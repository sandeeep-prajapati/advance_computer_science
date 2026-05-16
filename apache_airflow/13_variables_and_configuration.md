# Chapter 13: Variables & Configuration

## Airflow Variables

Variables are key-value pairs stored in the Airflow metadata DB. Use them for runtime configuration that changes between environments.

### Managing Variables

**Via UI:** Admin → Variables → "+"

**Via CLI:**
```bash
# Set a variable
airflow variables set environment production
airflow variables set max_retries 3
airflow variables set config_json '{"batch_size": 500, "timeout": 30}'

# Get a variable
airflow variables get environment

# List all
airflow variables list

# Delete
airflow variables delete old_var

# Import from JSON file
airflow variables import variables.json

# Export to JSON file
airflow variables export variables.json
```

**Via Python:**
```python
from airflow.models import Variable

# Get a variable (returns string)
env = Variable.get("environment")

# With a default value
env = Variable.get("environment", default_var="staging")

# Deserialize JSON automatically
config = Variable.get("config_json", deserialize_json=True)
batch_size = config["batch_size"]   # int, not string

# Set a variable
Variable.set("last_run_date", "2024-01-15")
Variable.set("config_json", {"batch_size": 1000}, serialize_json=True)
```

---

## Using Variables in DAGs

```python
from airflow.decorators import dag, task
from airflow.models import Variable
from datetime import datetime

@dag(start_date=datetime(2024, 1, 1), schedule="@daily", catchup=False)
def configured_pipeline():

    @task
    def extract() -> dict:
        # Read variable at task execution time (not parse time!)
        config = Variable.get("etl_config", deserialize_json=True)
        source = Variable.get("source_db", default_var="prod-db")
        print(f"Extracting from {source} with batch_size={config['batch_size']}")
        return {"rows": config["batch_size"]}

    @task
    def load(data: dict) -> None:
        dest = Variable.get("destination", default_var="warehouse")
        print(f"Loading {data['rows']} rows to {dest}")

    extract() >> load(extract())

configured_pipeline()
```

**Important:** Read variables inside task functions, not at DAG parse time — reading at module level re-queries the DB on every scheduler heartbeat.

```python
# BAD: read at parse time (called every ~30 seconds by scheduler)
ENV = Variable.get("environment")  # ← module-level read

@dag(...)
def my_dag():
    ...

# GOOD: read at execution time
@task
def my_task():
    env = Variable.get("environment")  # ← only called when task runs
```

---

## Variables via Jinja Templating

```python
from airflow.operators.python import PythonOperator
from airflow.operators.bash import BashOperator

# In templated fields, use var.value or var.json
bash_task = BashOperator(
    task_id      = "use_variable",
    bash_command = "echo {{ var.value.environment }} && echo {{ var.json.etl_config.batch_size }}",
)
```

---

## DAG Params (Runtime Parameters)

Params allow passing values at trigger time. Unlike Variables, they are DAG-level and can differ per DAG run.

```python
from airflow.decorators import dag, task
from airflow.models.param import Param
from datetime import datetime

@dag(
    start_date = datetime(2024, 1, 1),
    schedule   = None,
    catchup    = False,
    params     = {
        "env":        Param("staging", type="string", enum=["staging", "production"]),
        "batch_size": Param(100, type="integer", minimum=1, maximum=10000),
        "dry_run":    Param(False, type="boolean"),
        "tables":     Param(["orders"], type="array"),
    },
)
def parameterized_pipeline():

    @task
    def run(**context) -> None:
        params     = context["params"]
        env        = params["env"]
        batch_size = params["batch_size"]
        dry_run    = params["dry_run"]
        tables     = params["tables"]
        print(f"Running on {env}, batch={batch_size}, dry_run={dry_run}, tables={tables}")

    run()

parameterized_pipeline()
```

```bash
# Trigger with custom params
airflow dags trigger parameterized_pipeline \
  --conf '{"env": "production", "batch_size": 500, "dry_run": false, "tables": ["orders", "products"]}'
```

```python
# Jinja access to params
BashOperator(
    task_id      = "show_params",
    bash_command = "echo {{ params.env }} {{ params.batch_size }}",
)
```

---

## Environment Variables in Airflow

### Setting Airflow Config via Environment

```bash
# Pattern: AIRFLOW__{SECTION}__{KEY} (uppercase, double underscores)
export AIRFLOW__CORE__EXECUTOR=CeleryExecutor
export AIRFLOW__DATABASE__SQL_ALCHEMY_CONN=postgresql+psycopg2://airflow:airflow@postgres/airflow
export AIRFLOW__CORE__FERNET_KEY=your_fernet_key_here
export AIRFLOW__WEBSERVER__SECRET_KEY=your_secret_key_here
export AIRFLOW__CELERY__BROKER_URL=redis://redis:6379/0
```

### Injecting Custom Env Vars

```python
from airflow.operators.python import PythonOperator
import os

def my_task(**context):
    # Access env vars at runtime
    api_key  = os.environ.get("MY_API_KEY")
    base_url = os.environ.get("BASE_URL", "https://api.example.com")
    print(f"Calling {base_url} with key {api_key[:4]}...")

task = PythonOperator(
    task_id         = "api_call",
    python_callable = my_task,
    env_vars        = {        # inject per-task env vars
        "TASK_MODE": "production",
    },
)
```

---

## Jinja Templating Reference

Templating is available in `template_fields` of operators. The template is rendered at execution time.

### Built-in Variables

```python
# In templates: {{ variable }}
# In Python: context["variable"]

# Date/time
{{ ds }}                           # 2024-01-15 (data_interval_start date)
{{ ds_nodash }}                    # 20240115
{{ ts }}                           # 2024-01-15T00:00:00+00:00
{{ ts_nodash }}                    # 20240115T000000

# Data interval
{{ data_interval_start | ds }}     # 2024-01-15
{{ data_interval_end   | ds }}     # 2024-01-16
{{ data_interval_start | ts }}     # 2024-01-15T00:00:00+00:00

# DAG and run info
{{ dag.dag_id }}
{{ dag_run.run_id }}
{{ dag_run.conf }}                  # dict from trigger --conf
{{ run_id }}
{{ task.task_id }}

# Previous/next run
{{ prev_ds }}                       # previous execution date
{{ next_ds }}

# Macros (callable functions)
{{ macros.ds_add(ds, 7) }}         # add 7 days to ds
{{ macros.ds_format(ds, "%Y-%m-%d", "%Y/%m/%d") }}
{{ macros.datetime.now() }}

# Variables and params
{{ var.value.my_var }}              # Variable.get("my_var")
{{ var.json.my_config.key }}        # Variable.get("my_config", deserialize_json=True)["key"]
{{ params.env }}                    # DAG param "env"

# Connections
{{ conn.my_postgres.host }}
{{ conn.my_postgres.login }}
```

### Jinja in Operators

```python
BashOperator(
    task_id      = "run_script",
    bash_command = """
        python /scripts/etl.py \
          --date {{ ds }} \
          --env  {{ var.value.environment }} \
          --dest {{ params.destination }}
    """,
)

PythonOperator(
    task_id             = "query_db",
    python_callable     = run_query,
    templates_dict      = {           # templated dict passed to op_kwargs
        "query_date": "{{ ds }}",
        "env":        "{{ var.value.environment }}",
    },
    op_kwargs           = {"database": "orders"},
)
```

---

## Custom Macros

```python
# In the DAG definition, register custom macros
def business_date(ds: str) -> str:
    """Return previous business day if ds is Monday."""
    from datetime import datetime, timedelta
    dt = datetime.strptime(ds, "%Y-%m-%d")
    if dt.weekday() == 0:   # Monday
        return (dt - timedelta(days=3)).strftime("%Y-%m-%d")
    return (dt - timedelta(days=1)).strftime("%Y-%m-%d")

with DAG(
    "my_dag",
    user_defined_macros = {"business_date": business_date},
    ...
) as dag:

    # Now available in templates
    task = BashOperator(
        task_id      = "run",
        bash_command = "python etl.py --date {{ business_date(ds) }}",
    )
```

---

## Variables File for Environment Management

```json
// variables/staging.json
{
  "environment":    "staging",
  "source_db_host": "staging-db.internal",
  "batch_size":     100,
  "etl_config": {
    "timeout":    30,
    "retries":    2,
    "batch_size": 100
  }
}
```

```json
// variables/production.json
{
  "environment":    "production",
  "source_db_host": "prod-db.internal",
  "batch_size":     5000,
  "etl_config": {
    "timeout":    60,
    "retries":    3,
    "batch_size": 5000
  }
}
```

```bash
# Load env-specific variables
airflow variables import variables/production.json
```

---

## Variables vs Params vs Env Vars

| Feature | Variable | Param | Env Var |
|---|---|---|---|
| Scope | Global | Per DAG run | Process/container |
| Set via | UI/CLI/API | Trigger conf | OS/docker |
| Access | `Variable.get()` | `context["params"]` | `os.environ` |
| Runtime change | Yes | Yes (at trigger) | No (at startup) |
| Sensitive data | Use secrets backend | No | Use secrets backend |
| Best for | Shared config | Per-run overrides | Infrastructure config |
