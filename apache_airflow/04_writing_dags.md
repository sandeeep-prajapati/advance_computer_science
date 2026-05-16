# Chapter 4: Writing DAGs

## Traditional DAG Style

```python
# dags/etl_pipeline.py
from airflow import DAG
from airflow.operators.python import PythonOperator
from airflow.operators.bash import BashOperator
from datetime import datetime, timedelta

default_args = {
    "owner":             "data-team",
    "retries":           2,
    "retry_delay":       timedelta(minutes=5),
    "email_on_failure":  True,
    "email":             ["alerts@company.com"],
    "execution_timeout": timedelta(hours=2),
}

def extract(**context):
    ds = context["ds"]   # data interval start date as string YYYY-MM-DD
    print(f"Extracting data for {ds}")
    return f"s3://bucket/raw/{ds}/data.parquet"

def transform(**context):
    path = context["ti"].xcom_pull(task_ids="extract")
    print(f"Transforming data from {path}")

def load(**context):
    print("Loading data to warehouse")

with DAG(
    dag_id       = "etl_pipeline",
    default_args = default_args,
    start_date   = datetime(2024, 1, 1),
    schedule     = "0 2 * * *",    # 2 AM daily
    catchup      = False,
    max_active_runs = 1,
    tags         = ["etl", "daily"],
) as dag:

    extract_task   = PythonOperator(task_id="extract",   python_callable=extract)
    transform_task = PythonOperator(task_id="transform", python_callable=transform)
    load_task      = PythonOperator(task_id="load",      python_callable=load)
    notify_task    = BashOperator(
        task_id      = "notify",
        bash_command = "echo 'Pipeline complete for {{ ds }}'",
    )

    extract_task >> transform_task >> load_task >> notify_task
```

---

## TaskFlow API (Airflow 2.0+)

The TaskFlow API uses decorators for cleaner, more Pythonic DAG definitions. XCom is handled automatically.

```python
# dags/taskflow_etl.py
from airflow.decorators import dag, task
from datetime import datetime

@dag(
    dag_id      = "taskflow_etl",
    start_date  = datetime(2024, 1, 1),
    schedule    = "@daily",
    catchup     = False,
    tags        = ["etl"],
)
def etl_pipeline():

    @task()
    def extract(ds: str = None) -> dict:
        """Extract raw data from source."""
        print(f"Extracting for {ds}")
        return {"path": f"s3://bucket/raw/{ds}/", "records": 10000}

    @task()
    def transform(extracted: dict) -> dict:
        """Clean and transform data."""
        print(f"Transforming {extracted['records']} records")
        return {"records_transformed": extracted["records"], "status": "ok"}

    @task()
    def load(transformed: dict) -> None:
        """Load data to warehouse."""
        print(f"Loading {transformed['records_transformed']} records")

    # Build pipeline by calling tasks
    raw    = extract()
    clean  = transform(raw)
    load(clean)

# Instantiate the DAG
etl_pipeline()
```

**TaskFlow advantages:**
- Return values automatically become XComs
- Parameters flow naturally between tasks
- Less boilerplate
- Type hints improve readability

---

## Dependency Patterns

### Linear Chain
```python
t1 >> t2 >> t3 >> t4
# or equivalently
t1.set_downstream(t2)
t3.set_upstream(t2)
```

### Fan-Out (one to many)
```python
extract >> [transform_a, transform_b, transform_c]
```

### Fan-In (many to one)
```python
[transform_a, transform_b, transform_c] >> load
```

### Diamond Pattern
```python
extract >> [transform_a, transform_b]
[transform_a, transform_b] >> load
load >> notify
```

### Complex Pipeline
```python
from airflow.operators.empty import EmptyOperator

start   = EmptyOperator(task_id="start")
end     = EmptyOperator(task_id="end")

start >> [task_a, task_b]
task_a >> [task_c, task_d]
task_b >> task_d
[task_c, task_d] >> end
```

---

## Trigger Rules

By default, a task runs only when all upstream tasks succeed. Trigger rules change this:

```python
from airflow.utils.trigger_rule import TriggerRule

cleanup = PythonOperator(
    task_id      = "cleanup",
    python_callable = cleanup_fn,
    trigger_rule = TriggerRule.ALL_DONE,  # run regardless of upstream state
)
```

| Rule | When Task Runs |
|---|---|
| `all_success` | All parents succeeded (default) |
| `all_failed` | All parents failed |
| `all_done` | All parents finished (any state) |
| `one_success` | At least one parent succeeded |
| `one_failed` | At least one parent failed |
| `none_failed` | No parents failed (success or skipped) |
| `none_skipped` | No parents were skipped |
| `always` | Regardless of parent states |

---

## Branching

```python
from airflow.operators.python import BranchPythonOperator
from airflow.operators.empty import EmptyOperator

def choose_branch(**context):
    hour = context["logical_date"].hour
    if hour < 12:
        return "morning_task"
    else:
        return "afternoon_task"   # task_id to execute next

branch = BranchPythonOperator(
    task_id         = "branch",
    python_callable = choose_branch,
)

morning_task   = EmptyOperator(task_id="morning_task")
afternoon_task = EmptyOperator(task_id="afternoon_task")
join           = EmptyOperator(task_id="join", trigger_rule="none_failed_min_one_success")

branch >> [morning_task, afternoon_task] >> join
```

---

## ShortCircuitOperator

Skip all downstream tasks if a condition is False:

```python
from airflow.operators.python import ShortCircuitOperator

def is_weekday(**context):
    return context["logical_date"].weekday() < 5   # Mon-Fri

check_weekday = ShortCircuitOperator(
    task_id         = "check_weekday",
    python_callable = is_weekday,
)

check_weekday >> process_data >> send_report
# On weekends: process_data and send_report are SKIPPED
```

---

## Jinja Templating

Airflow templates fields of operators with Jinja2:

```python
BashOperator(
    task_id      = "process",
    bash_command = """
        python process.py \
          --date {{ ds }} \
          --start {{ data_interval_start }} \
          --end {{ data_interval_end }} \
          --run-id {{ run_id }}
    """,
)
```

### Template Variables

| Variable | Value |
|---|---|
| `{{ ds }}` | Data interval start date: `2024-01-15` |
| `{{ ds_nodash }}` | `20240115` |
| `{{ ts }}` | ISO timestamp: `2024-01-15T06:00:00+00:00` |
| `{{ data_interval_start }}` | `2024-01-15 06:00:00+00:00` |
| `{{ data_interval_end }}` | `2024-01-16 06:00:00+00:00` |
| `{{ run_id }}` | `scheduled__2024-01-15T06:00:00+00:00` |
| `{{ dag.dag_id }}` | `my_dag` |
| `{{ task.task_id }}` | `my_task` |
| `{{ params.my_param }}` | Value from DAG params |
| `{{ var.value.my_var }}` | Value from Airflow Variable |
| `{{ conn.my_conn.host }}` | Connection field |
| `{{ macros.ds_add(ds, 7) }}` | Date arithmetic |

---

## DAG Documentation

Add documentation visible in the Airflow UI:

```python
dag.doc_md = """
## Daily ETL Pipeline

Extracts from PostgreSQL, transforms with dbt, loads to BigQuery.

**Schedule:** Daily at 2 AM UTC  
**Owner:** Data Engineering Team  
**SLA:** Must complete by 4 AM UTC  

### Dependencies
- PostgreSQL source must be available
- dbt models must be up to date
"""

task.doc_md = """
Runs dbt models in the `staging` layer.
Fails if any model has > 1% error rate.
"""
```
