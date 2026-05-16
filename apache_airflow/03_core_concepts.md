# Chapter 3: Core Concepts

## DAG (Directed Acyclic Graph)

A DAG is a Python object that defines:
- The tasks to run
- The order and dependencies between tasks
- When to run (schedule)
- How to behave on failure

```
"Directed" = dependencies have direction (A → B means B depends on A)
"Acyclic"  = no cycles (A → B → A is not allowed)
"Graph"    = tasks are nodes, dependencies are edges
```

```python
from airflow import DAG
from datetime import datetime, timedelta

with DAG(
    dag_id            = "my_pipeline",
    description       = "Daily data pipeline",
    start_date        = datetime(2024, 1, 1),
    schedule          = "0 6 * * *",    # daily at 6 AM
    catchup           = False,
    max_active_runs   = 1,
    default_args      = {
        "owner":           "data-team",
        "retries":         2,
        "retry_delay":     timedelta(minutes=5),
        "email_on_failure": True,
        "email":           ["alerts@company.com"],
    },
    tags              = ["etl", "daily"],
) as dag:
    pass  # tasks go here
```

### Key DAG Parameters

| Parameter | Description |
|---|---|
| `dag_id` | Unique identifier |
| `start_date` | First date the DAG can run |
| `schedule` | Cron, timedelta, timetable, or `None` |
| `catchup` | Run missed intervals since start_date (default: True) |
| `max_active_runs` | Max concurrent DAG runs |
| `default_args` | Default task parameters |
| `tags` | Filter DAGs in the UI |
| `params` | Runtime parameters (Airflow 2.2+) |

---

## Task

A Task is a node in the DAG graph. Each task is an instance of an Operator.

```python
from airflow.operators.python import PythonOperator

task = PythonOperator(
    task_id          = "process_data",   # unique within the DAG
    python_callable  = process_fn,
    dag              = dag,              # or use context manager (with DAG)

    # Task-level overrides
    retries          = 3,
    retry_delay      = timedelta(minutes=2),
    execution_timeout= timedelta(hours=1),
    depends_on_past  = False,
    trigger_rule     = "all_success",    # default
    priority_weight  = 10,
    pool             = "default_pool",
    queue            = "default",
)
```

### Task States

```
none     → scheduled → queued → running → success
                                        → failed → up_for_retry → queued (retry)
                                        → upstream_failed
                                        → skipped
```

---

## Operators

Operators are templates for tasks. Airflow ships with hundreds:

| Operator | Use Case |
|---|---|
| `PythonOperator` | Run any Python function |
| `BashOperator` | Run shell commands |
| `EmailOperator` | Send email |
| `DummyOperator` / `EmptyOperator` | No-op (branching, grouping) |
| `BranchPythonOperator` | Conditional branching |
| `TriggerDagRunOperator` | Trigger another DAG |
| `SparkSubmitOperator` | Submit a Spark job |
| `BigQueryOperator` | Run BigQuery SQL |
| `S3ToRedshiftOperator` | Copy S3 → Redshift |
| `KubernetesPodOperator` | Run a container in Kubernetes |
| `DockerOperator` | Run a Docker container |
| `HttpOperator` | Make HTTP requests |
| `SlackWebhookOperator` | Send Slack message |

---

## Sensors

Sensors are special operators that **wait** for a condition to become True before allowing downstream tasks to proceed.

```python
from airflow.sensors.filesystem import FileSensor

wait_for_file = FileSensor(
    task_id         = "wait_for_data",
    filepath        = "/data/input/{{ ds }}/data.csv",
    poke_interval   = 60,    # check every 60 seconds
    timeout         = 3600,  # fail after 1 hour
    mode            = "reschedule",  # release slot while waiting
)

wait_for_file >> process_data
```

| Sensor | Waits For |
|---|---|
| `FileSensor` | File to appear on disk |
| `S3KeySensor` | Object to appear in S3 |
| `HttpSensor` | HTTP endpoint to return success |
| `SqlSensor` | SQL query to return rows |
| `ExternalTaskSensor` | Another DAG/task to complete |
| `PythonSensor` | Custom Python condition |
| `TimeDeltaSensor` | Time to pass |

---

## Hooks

Hooks provide a clean interface to external systems. They handle connection management, authentication, and retries. Operators use hooks internally.

```python
from airflow.providers.postgres.hooks.postgres import PostgresHook

hook = PostgresHook(postgres_conn_id="my_postgres")

# Use directly
records = hook.get_records("SELECT * FROM orders WHERE date = '{{ ds }}'")
hook.run("INSERT INTO staging SELECT * FROM raw WHERE dt = %s", parameters=["2024-01-01"])

# Get a SQLAlchemy engine
engine = hook.get_sqlalchemy_engine()
df = pd.read_sql("SELECT ...", engine)
```

---

## Connections

Connections store credentials for hooks. Managed in:
- Airflow UI: Admin → Connections
- Environment variables: `AIRFLOW_CONN_<CONN_ID>=uri`
- Secrets backends (Vault, AWS Secrets Manager)

```bash
# CLI: add a connection
airflow connections add \
  --conn-id my_postgres \
  --conn-type postgres \
  --conn-host localhost \
  --conn-port 5432 \
  --conn-login myuser \
  --conn-password mypassword \
  --conn-schema mydb

# Environment variable format
export AIRFLOW_CONN_MY_POSTGRES="postgresql://myuser:mypassword@localhost:5432/mydb"
```

---

## XCom (Cross-Communication)

XComs let tasks share small pieces of data with each other via the metadata database.

```python
def push_fn(**context):
    context["ti"].xcom_push(key="result", value={"count": 42, "status": "ok"})

def pull_fn(**context):
    result = context["ti"].xcom_pull(task_ids="push_task", key="result")
    print(f"Got: {result}")   # {"count": 42, "status": "ok"}

push_task = PythonOperator(task_id="push_task", python_callable=push_fn)
pull_task = PythonOperator(task_id="pull_task", python_callable=pull_fn)
push_task >> pull_task
```

**XCom size limit:** The default metadata database stores XComs up to ~64 KB. For larger data, use a custom backend (S3, GCS) or pass file paths via XCom.

---

## DAG Run vs Task Instance

```
DAG: daily_etl  (definition)
  │
  ├─ DAG Run: daily_etl @ 2024-01-01  (one execution)
  │    ├─ Task Instance: extract @ 2024-01-01   (state: success)
  │    ├─ Task Instance: transform @ 2024-01-01 (state: success)
  │    └─ Task Instance: load @ 2024-01-01      (state: running)
  │
  └─ DAG Run: daily_etl @ 2024-01-02  (another execution)
       ├─ Task Instance: extract @ 2024-01-02   (state: success)
       ├─ Task Instance: transform @ 2024-01-02 (state: failed)
       └─ Task Instance: load @ 2024-01-02      (state: upstream_failed)
```

### Execution Date vs Data Interval

In Airflow 2.2+, each DAG run has a **data interval** instead of just an execution date:

```
Schedule: @daily
DAG Run for 2024-01-02:
  data_interval_start = 2024-01-02 00:00:00
  data_interval_end   = 2024-01-03 00:00:00
  run_id              = scheduled__2024-01-02T00:00:00
  logical_date        = 2024-01-02 (same as execution_date)
```

Templates use `{{ ds }}` (data interval start date), `{{ data_interval_end }}`, etc.
