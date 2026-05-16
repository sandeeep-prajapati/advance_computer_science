# Chapter 11: Dynamic DAGs & Task Mapping

## Dynamic Task Mapping (Airflow 2.3+)

Generate tasks at runtime based on upstream data — no need to know the count at DAG parse time.

```python
from airflow.decorators import dag, task
from datetime import datetime

@dag(start_date=datetime(2024, 1, 1), schedule="@daily", catchup=False)
def dynamic_pipeline():

    @task
    def get_files() -> list[str]:
        # Could query S3, DB, API, etc.
        return ["file_a.csv", "file_b.csv", "file_c.csv"]

    @task
    def process_file(filename: str) -> dict:
        print(f"Processing {filename}")
        return {"file": filename, "rows": 1000}

    @task
    def summarize(results: list[dict]) -> None:
        total = sum(r["rows"] for r in results)
        print(f"Total rows processed: {total}")

    files   = get_files()
    results = process_file.expand(filename=files)   # creates one task per file
    summarize(results)

dynamic_pipeline()
```

```
process_file[0]  → process_file[file_a.csv]
process_file[1]  → process_file[file_b.csv]
process_file[2]  → process_file[file_c.csv]
```

---

## expand() vs expand_kwargs()

### expand() — single argument

```python
@task
def transform(record: dict) -> dict:
    return {**record, "processed": True}

records = [{"id": 1}, {"id": 2}, {"id": 3}]
transform.expand(record=records)
```

### expand_kwargs() — multiple arguments together

```python
@task
def load_table(table: str, schema: str, mode: str) -> None:
    print(f"Loading {schema}.{table} in {mode} mode")

load_table.expand_kwargs([
    {"table": "orders",    "schema": "raw", "mode": "overwrite"},
    {"table": "customers", "schema": "raw", "mode": "append"},
    {"table": "products",  "schema": "raw", "mode": "overwrite"},
])
```

### partial() — fix some args, expand others

```python
@task
def upload_to_s3(filename: str, bucket: str, prefix: str) -> None:
    print(f"Uploading {filename} to s3://{bucket}/{prefix}/{filename}")

files = ["a.csv", "b.csv", "c.csv"]

# Fix bucket and prefix, expand only filename
upload_to_s3.partial(bucket="my-bucket", prefix="output/").expand(filename=files)
```

---

## Chained Dynamic Tasks

```python
@dag(start_date=datetime(2024, 1, 1), schedule="@daily", catchup=False)
def etl_pipeline():

    @task
    def get_sources() -> list[str]:
        return ["postgres", "mysql", "mongodb"]

    @task
    def extract(source: str) -> dict:
        return {"source": source, "data": f"data_from_{source}"}

    @task
    def transform(payload: dict) -> dict:
        return {**payload, "transformed": True}

    @task
    def load(payload: dict) -> None:
        print(f"Loading: {payload}")

    sources   = get_sources()
    extracted = extract.expand(source=sources)
    transformed = transform.expand(payload=extracted)  # still one task per item
    load.expand(payload=transformed)

etl_pipeline()
```

---

## Dynamic Mapping with Operators (non-TaskFlow)

```python
from airflow import DAG
from airflow.operators.python import PythonOperator
from datetime import datetime

def process(filename: str, **context) -> None:
    print(f"Processing {filename}")

with DAG("dynamic_operator_dag", start_date=datetime(2024, 1, 1), schedule=None) as dag:

    # expand() works on classic operators too
    process_task = PythonOperator.partial(
        task_id         = "process_file",
        python_callable = process,
    ).expand(op_kwargs=[
        {"filename": "a.csv"},
        {"filename": "b.csv"},
        {"filename": "c.csv"},
    ])
```

---

## Dynamic DAGs via Factory Functions

Generate multiple DAGs programmatically at parse time (pre-2.3 pattern, still valid):

```python
# dags/generated_dags.py
from airflow import DAG
from airflow.operators.python import PythonOperator
from datetime import datetime

CONFIGS = [
    {"dag_id": "etl_customers", "table": "customers", "schedule": "@daily"},
    {"dag_id": "etl_orders",    "table": "orders",    "schedule": "@hourly"},
    {"dag_id": "etl_products",  "table": "products",  "schedule": "@weekly"},
]

def make_dag(dag_id: str, table: str, schedule: str) -> DAG:
    with DAG(
        dag_id      = dag_id,
        start_date  = datetime(2024, 1, 1),
        schedule    = schedule,
        catchup     = False,
        tags        = ["generated", table],
    ) as dag:

        def extract(**ctx):
            print(f"Extracting {table} for {ctx['ds']}")

        def load(**ctx):
            print(f"Loading {table}")

        t1 = PythonOperator(task_id="extract", python_callable=extract)
        t2 = PythonOperator(task_id="load",    python_callable=load)
        t1 >> t2

    return dag

# Register each DAG in globals so Airflow picks them up
for cfg in CONFIGS:
    globals()[cfg["dag_id"]] = make_dag(**cfg)
```

---

## DAG Factory with YAML Config

```yaml
# config/pipelines.yaml
pipelines:
  - name: sales_etl
    source: postgres
    destination: bigquery
    schedule: "0 2 * * *"
    tables:
      - orders
      - order_items
      - customers

  - name: marketing_etl
    source: mysql
    destination: snowflake
    schedule: "0 4 * * *"
    tables:
      - campaigns
      - clicks
```

```python
# dags/yaml_factory.py
import yaml
import os
from airflow import DAG
from airflow.decorators import task
from datetime import datetime

config_path = os.path.join(os.path.dirname(__file__), "../config/pipelines.yaml")
with open(config_path) as f:
    config = yaml.safe_load(f)

def make_pipeline_dag(pipeline: dict) -> DAG:
    name        = pipeline["name"]
    source      = pipeline["source"]
    destination = pipeline["destination"]
    schedule    = pipeline["schedule"]
    tables      = pipeline["tables"]

    with DAG(
        dag_id     = name,
        start_date = datetime(2024, 1, 1),
        schedule   = schedule,
        catchup    = False,
    ) as dag:

        @task
        def extract_table(table: str) -> dict:
            print(f"Extracting {table} from {source}")
            return {"table": table, "rows": 1000}

        @task
        def load_table(payload: dict) -> None:
            print(f"Loading {payload['table']} to {destination}")

        results = extract_table.expand(table=tables)
        load_table.expand(payload=results)

    return dag

for pipeline in config["pipelines"]:
    globals()[pipeline["name"]] = make_pipeline_dag(pipeline)
```

---

## Dynamic Mapping — Limits and Controls

```python
# Limit max concurrent mapped tasks (Airflow 2.7+)
@task(max_active_tis_per_dag=10)
def process(item: str) -> None:
    print(item)

# Or per-dagrun
@task(max_active_tis_per_dagrun=5)
def process(item: str) -> None:
    print(item)

# Accessing index inside a mapped task
@task
def process(item: str, **context) -> None:
    map_index = context["task_instance"].map_index
    print(f"Processing item #{map_index}: {item}")
```

---

## Mixing Static and Dynamic Tasks

```python
@dag(start_date=datetime(2024, 1, 1), schedule="@daily", catchup=False)
def mixed_pipeline():

    @task
    def validate_config() -> None:
        print("Config valid")

    @task
    def get_items() -> list[str]:
        return ["x", "y", "z"]

    @task
    def process(item: str) -> str:
        return f"done_{item}"

    @task
    def aggregate(results: list[str]) -> None:
        print(f"All done: {results}")

    @task
    def notify() -> None:
        print("Pipeline complete")

    # Static → Dynamic → Static
    validate_config() >> get_items()
    items   = get_items()
    results = process.expand(item=items)
    aggregate(results) >> notify()

mixed_pipeline()
```

---

## Dynamic Task Mapping — Key Concepts

```
Static mapping (expand on a literal list):
  task.expand(x=[1, 2, 3])   ← count known at parse time

Dynamic mapping (expand on upstream XCom):
  items = get_items()
  task.expand(x=items)       ← count known only at runtime

Nested expand:
  task.expand(a=upstream_a, b=upstream_b)  ← cartesian product!
  Use expand_kwargs() instead to zip args.

Accessing mapped results downstream:
  Single mapped output → next expand() gets a list
  All complete → next task receives aggregated list
```

---

## Best Practices

```
1. Use dynamic task mapping (Airflow 2.3+) instead of DAG factories
   when the structure is the same, only data changes.

2. Use DAG factories when you need truly separate DAGs
   with different schedules, owners, or dependencies.

3. Avoid generating thousands of mapped tasks (> 1000) in one DAG run
   → Use chunking or batching instead.

4. Keep mapped task functions pure — they run in parallel.

5. Use max_active_tis_per_dag/dagrun to rate-limit mapped tasks
   against external systems.

6. YAML/JSON configs in a separate directory (not in dags/)
   prevent accidental re-parsing when config changes.
```
