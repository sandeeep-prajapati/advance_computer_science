# Chapter 8: XCom & Task Communication

## What is XCom?

XCom (Cross-Communication) allows tasks to exchange small amounts of data. Values are stored in the Airflow metadata database and retrieved by downstream tasks.

```
Task A (push) ──► Metadata DB ──► Task B (pull)
                  xcom table:
                  dag_id, run_id, task_id, key, value
```

**XCom limitations:**
- Default backend: stored in SQL DB, max ~64 KB
- Not designed for large data (DataFrames, files)
- For large data: store in S3/GCS, pass the path via XCom

---

## Traditional XCom Push/Pull

```python
from airflow import DAG
from airflow.operators.python import PythonOperator
from datetime import datetime

def extract(**context):
    # Push result to XCom
    context["ti"].xcom_push(key="row_count",    value=10000)
    context["ti"].xcom_push(key="output_path",  value="s3://bucket/out/2024-01-01/")
    # Returning a value also pushes to XCom with key="return_value"
    return {"status": "success", "rows": 10000}

def transform(**context):
    ti = context["ti"]

    # Pull by task_id and key
    row_count   = ti.xcom_pull(task_ids="extract", key="row_count")
    output_path = ti.xcom_pull(task_ids="extract", key="output_path")

    # Pull return_value (implicit push)
    result = ti.xcom_pull(task_ids="extract")   # key defaults to "return_value"

    print(f"Got {row_count} rows from {output_path}")
    return f"transformed_{result['rows']}_rows"

def load(**context):
    ti = context["ti"]
    # Pull from multiple tasks
    transform_result = ti.xcom_pull(task_ids="transform")
    extract_path     = ti.xcom_pull(task_ids="extract", key="output_path")
    print(f"Loading: {transform_result} from {extract_path}")
```

---

## TaskFlow XCom (Automatic)

The TaskFlow API handles XCom automatically — return values flow between tasks:

```python
from airflow.decorators import dag, task
from datetime import datetime

@dag(start_date=datetime(2024,1,1), schedule="@daily", catchup=False)
def taskflow_pipeline():

    @task()
    def extract(ds: str = None) -> dict:
        # Return value = XCom push (key="return_value")
        return {"path": f"s3://bucket/raw/{ds}/", "count": 10000}

    @task()
    def transform(data: dict) -> list:
        # data = auto-pulled from extract's return value
        print(f"Processing {data['count']} records from {data['path']}")
        return [{"id": i, "value": i * 2} for i in range(data["count"])]

    @task()
    def load(records: list) -> None:
        print(f"Loading {len(records)} records")
        # Insert to DB...

    @task()
    def notify(path: str, count: int) -> None:
        print(f"Processed {count} records, output at {path}")

    # Data flows naturally
    raw_data  = extract()
    clean     = transform(raw_data)
    load(clean)

    # Pass specific fields to different tasks
    notify(path=raw_data["path"], count=raw_data["count"])

taskflow_pipeline()
```

---

## Multiple Return Values (Airflow 2.3+)

```python
from airflow.decorators import task
from typing import Tuple

@task(multiple_outputs=True)
def extract() -> dict:
    # Each key becomes a separate XCom entry
    return {
        "path":  "s3://bucket/out/",
        "count": 10000,
        "status": "success",
    }

@dag(...)
def my_dag():
    data = extract()
    # Access individual keys
    process_data(path=data["path"])
    send_metrics(count=data["count"])
```

---

## XCom with Jinja Templates

```python
BashOperator(
    task_id      = "use_xcom",
    bash_command = """
        echo "Processing path: {{ ti.xcom_pull(task_ids='extract', key='path') }}"
        echo "Row count: {{ ti.xcom_pull(task_ids='extract', key='count') }}"
    """,
)
```

---

## Custom XCom Backend

For large data (DataFrames, models), use a custom backend that stores in S3/GCS:

```python
# plugins/custom_xcom_backend.py
from airflow.models.xcom import BaseXCom
import boto3
import json
import pickle
import uuid

class S3XComBackend(BaseXCom):
    BUCKET = "my-airflow-xcoms"

    @staticmethod
    def serialize_value(value, *, key=None, task_id=None, dag_id=None,
                        run_id=None, map_index=None):
        import pandas as pd
        if isinstance(value, pd.DataFrame):
            key_s3 = f"xcoms/{dag_id}/{run_id}/{task_id}/{uuid.uuid4()}.parquet"
            s3 = boto3.client("s3")
            s3.put_object(
                Bucket = S3XComBackend.BUCKET,
                Key    = key_s3,
                Body   = value.to_parquet(),
            )
            return BaseXCom.serialize_value({"__type": "s3_parquet", "key": key_s3})
        return BaseXCom.serialize_value(value)

    @staticmethod
    def deserialize_value(result):
        data = BaseXCom.deserialize_value(result)
        if isinstance(data, dict) and data.get("__type") == "s3_parquet":
            import pandas as pd
            s3 = boto3.client("s3")
            obj = s3.get_object(Bucket=S3XComBackend.BUCKET, Key=data["key"])
            return pd.read_parquet(obj["Body"])
        return data
```

```properties
# airflow.cfg
[core]
xcom_backend = plugins.custom_xcom_backend.S3XComBackend
```

---

## XCom Best Practices

```python
# 1. Use XCom only for small metadata, not bulk data
# BAD: Push 1M rows as a list
context["ti"].xcom_push("rows", [{"id": i} for i in range(1_000_000)])

# GOOD: Push S3 path, let downstream read from S3
s3_path = write_to_s3(rows)
context["ti"].xcom_push("output_path", s3_path)

# 2. Name your XCom keys clearly
ti.xcom_push("extract_row_count", 10000)       # clear
ti.xcom_push("result", 10000)                  # ambiguous

# 3. Use TaskFlow for clean data flow
@task()
def transform(data: dict) -> dict:  # type hints document the contract
    return process(data)

# 4. Clear XComs after use if storage is a concern
# In DAG:
dag = DAG(..., template_searchpath=["/path"])
# Or via:
ti.xcom_delete(key="large_temp_data")
```

---

## Passing Parameters Between DAGs

```python
# Trigger DAG with configuration
from airflow.operators.trigger_dagrun import TriggerDagRunOperator

trigger = TriggerDagRunOperator(
    task_id        = "trigger_downstream",
    trigger_dag_id = "downstream_pipeline",
    conf           = {
        "date":   "{{ ds }}",
        "source": "{{ dag.dag_id }}",
        "extra":  {"env": "production"},
    },
    wait_for_completion = True,
)

# In the downstream DAG, read the conf:
def process(**context):
    conf = context["dag_run"].conf
    date   = conf.get("date")
    source = conf.get("source")
    print(f"Triggered by {source} for date {date}")
```
