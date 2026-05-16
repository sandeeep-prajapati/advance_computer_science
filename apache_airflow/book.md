# Mastering Apache Airflow

> A comprehensive guide to Apache Airflow — from core concepts and DAG authoring to production deployment, cloud integrations, and custom extensions.

---

## Quick Navigation

- [Index & Overview](#index-overview)
- [Introduction to Apache Airflow](#introduction-to-apache-airflow)
- [Setup & Installation](#setup-installation)
- [Core Concepts](#core-concepts)
- [Writing DAGs](#writing-dags)
- [Operators](#operators)
- [Sensors](#sensors)
- [Hooks & Connections](#hooks-connections)
- [XCom & Task Communication](#xcom-task-communication)
- [Scheduling & Triggers](#scheduling-triggers)
- [Executors](#executors)
- [Dynamic DAGs & Task Mapping](#dynamic-dags-task-mapping)
- [Testing DAGs](#testing-dags)
- [Variables & Configuration](#variables-configuration)
- [Monitoring & Alerting](#monitoring-alerting)
- [Security](#security)
- [Performance & Scaling](#performance-scaling)
- [Airflow with Cloud Providers](#airflow-with-cloud-providers)
- [Custom Operators & Plugins](#custom-operators-plugins)
- [Airflow in Production](#airflow-in-production)
- [Learning Roadmap](#learning-roadmap)

---

# Apache Airflow — Complete Learning Guide

> A comprehensive, hands-on guide to Apache Airflow: from core concepts to production-grade data pipeline orchestration.

---

## Table of Contents

| Chapter | Title | Topics |
|---|---|---|
| 01 | [Introduction](01_introduction.md) | What is Airflow, DAGs, use cases, architecture |
| 02 | [Setup & Installation](02_setup_and_installation.md) | Docker Compose, local install, Astro CLI |
| 03 | [Core Concepts](03_core_concepts.md) | DAG, Task, Operator, Sensor, Hook, XCom |
| 04 | [Writing DAGs](04_writing_dags.md) | DAG definition, task flow, TaskFlow API |
| 05 | [Operators](05_operators.md) | BashOperator, PythonOperator, built-in operators |
| 06 | [Sensors](06_sensors.md) | FileSensor, S3Sensor, HttpSensor, custom sensors |
| 07 | [Hooks & Connections](07_hooks_and_connections.md) | Connections UI, hooks, secret backends |
| 08 | [XCom & Task Communication](08_xcom_and_task_communication.md) | XCom push/pull, TaskFlow, custom backends |
| 09 | [Scheduling & Triggers](09_scheduling_and_triggers.md) | Cron, timetables, datasets, event-driven |
| 10 | [Executors](10_executors.md) | Sequential, Local, Celery, Kubernetes |
| 11 | [Dynamic DAGs](11_dynamic_dags.md) | Dynamic task mapping, factory patterns |
| 12 | [Testing DAGs](12_testing_dags.md) | Unit tests, integration tests, DAG validation |
| 13 | [Variables & Configuration](13_variables_and_configuration.md) | Variables, Params, env vars, Jinja templating |
| 14 | [Monitoring & Alerting](14_monitoring_and_alerting.md) | Webserver UI, logs, callbacks, SLAs |
| 15 | [Security](15_security.md) | RBAC, authentication, secrets, Fernet encryption |
| 16 | [Performance & Scaling](16_performance_and_scaling.md) | Concurrency, pools, priority, tuning |
| 17 | [Airflow with Cloud Providers](17_airflow_with_cloud_providers.md) | AWS, GCP, Azure operators and hooks |
| 18 | [Custom Operators & Plugins](18_custom_operators_and_plugins.md) | Building operators, hooks, macros, plugins |
| 19 | [Airflow in Production](19_airflow_in_production.md) | MWAA, Cloud Composer, Astronomer, best practices |
| 20 | [Learning Roadmap](20_learning_roadmap.md) | 5-week plan, capstone projects, beyond Airflow |

---

## How to Read This Guide

- **Beginners**: Chapters 1–6 cover fundamentals. Build and run a DAG before moving on.
- **Data engineers**: Focus on Chapters 4–9, 11, 13 for building production pipelines.
- **Platform engineers**: Chapters 10, 15–16, 19 cover deployment and scaling.
- **Power users**: Chapters 11–12, 18 cover dynamic DAGs, testing, and custom components.
- **All readers**: Chapter 20 provides a structured 5-week hands-on roadmap.

---

## Apache Airflow Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                      Apache Airflow                              │
│                                                                  │
│  ┌─────────────┐   ┌─────────────┐   ┌──────────────────────┐   │
│  │  Web Server │   │  Scheduler  │   │  Metadata Database   │   │
│  │  (Flask/    │   │  (DAG       │   │  (PostgreSQL/MySQL)  │   │
│  │   Gunicorn) │   │   parsing + │   │  Tasks, runs, logs,  │   │
│  │  UI + REST  │   │   trigger)  │   │  connections, vars   │   │
│  └─────────────┘   └──────┬──────┘   └──────────────────────┘   │
│                            │                                     │
│                     ┌──────▼──────┐                              │
│                     │   Executor  │                              │
│                     │  (where     │                              │
│                     │  tasks run) │                              │
│                     └──────┬──────┘                              │
│                            │                                     │
│          ┌─────────────────┼─────────────────┐                  │
│          │                 │                 │                   │
│    ┌─────▼─────┐    ┌──────▼─────┐   ┌──────▼─────┐            │
│    │  Worker 1 │    │  Worker 2  │   │  Worker 3  │            │
│    └───────────┘    └────────────┘   └────────────┘            │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    DAGs Folder                          │    │
│  │  my_etl_dag.py   reporting_dag.py   ml_pipeline.py     │    │
│  └─────────────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────────────┘
```

---

## Core Abstractions

| Concept | Description |
|---|---|
| **DAG** | Directed Acyclic Graph — a workflow definition |
| **Task** | A single unit of work within a DAG |
| **Operator** | Template for a task (defines what to do) |
| **Sensor** | Operator that waits for a condition to be true |
| **Hook** | Interface to an external service (DB, API, cloud) |
| **Connection** | Stored credentials for hooks |
| **XCom** | Cross-task communication mechanism |
| **Variable** | Global key-value store for DAG configuration |
| **DAG Run** | A specific execution of a DAG |
| **Task Instance** | A specific execution of a task in a DAG Run |
| **Executor** | Component that runs task instances |


---

# Chapter 1: Introduction to Apache Airflow

## What is Apache Airflow?

Apache Airflow is an open-source **workflow orchestration platform** for programmatically authoring, scheduling, and monitoring data pipelines. Created at Airbnb in 2014 and donated to the Apache Software Foundation, it is now the industry standard for managing complex data workflows.

The core idea: **pipelines as code**. You define your workflows in Python, which gives you version control, testing, reuse, and the full power of a programming language.

```
Without Airflow:                With Airflow:
  Cron jobs everywhere            DAGs as Python code
  No visibility                   Full visibility via UI
  No retry logic                  Automatic retries
  No dependencies                 Task dependencies
  No alerting                     Email/Slack on failure
  Manual backfill                 Automatic backfill
```

---

## Why Airflow?

| Need | Airflow Solution |
|---|---|
| Schedule a pipeline daily | Cron expression or custom timetable |
| Run task B only after task A | Upstream dependencies |
| Retry failed tasks automatically | `retries` and `retry_delay` |
| Monitor all pipelines in one place | Web UI with DAG/task status |
| Alert on failure | Email/Slack callbacks |
| Re-run historical data | Backfill |
| Pass data between tasks | XCom |
| Parameterize a run | DAG Params |
| Wait for a file before processing | Sensors |

---

## DAGs: The Core Concept

A **DAG** (Directed Acyclic Graph) is a collection of tasks with dependencies between them. It defines:
- **What** to run (tasks)
- **When** to run (schedule)
- **In what order** to run (dependencies)

```
Example DAG: Daily ETL Pipeline

extract_data ──► transform_data ──► load_to_warehouse ──► send_report
                       │
                       ▼
                  validate_data
```

```python
from airflow import DAG
from airflow.operators.python import PythonOperator
from datetime import datetime

with DAG(
    dag_id      = "daily_etl",
    start_date  = datetime(2024, 1, 1),
    schedule    = "@daily",
) as dag:

    extract   = PythonOperator(task_id="extract",   python_callable=extract_fn)
    transform = PythonOperator(task_id="transform", python_callable=transform_fn)
    load      = PythonOperator(task_id="load",      python_callable=load_fn)

    extract >> transform >> load   # dependency chain
```

---

## Key Properties

| Property | Description |
|---|---|
| **Pipelines as code** | Python DAGs — version-controlled, testable |
| **Dynamic** | DAGs can generate tasks programmatically |
| **Extensible** | Hundreds of operators for every service |
| **Scalable** | Celery/Kubernetes executors for thousands of tasks |
| **Backfill** | Re-run historical DAG runs |
| **Retry** | Automatic retry with exponential backoff |
| **SLAs** | Miss SLA → alert |
| **UI** | Full web dashboard: DAG/task status, logs, Gantt |

---

## Airflow vs Alternatives

| Feature | Airflow | Prefect | Dagster | Luigi |
|---|---|---|---|---|
| Paradigm | DAG-centric | Flow-centric | Asset-centric | Pipeline |
| Language | Python | Python | Python | Python |
| UI | Web UI | Cloud/OSS UI | Web UI | Basic |
| Dynamic tasks | Yes (v2.3+) | Yes | Yes | Limited |
| Backfill | Yes | Limited | Yes | Manual |
| Cloud managed | MWAA, Composer | Prefect Cloud | Dagster Cloud | No |
| Ecosystem | Huge (providers) | Growing | Growing | Small |
| Learning curve | Medium | Low | Medium | Low |
| Best for | Complex ETL | Microservices | Data assets | Simple pipelines |

---

## Airflow Architecture

### Components

```
┌──────────────────────────────────────────────────────────────┐
│  Web Server        → UI + REST API (port 8080)               │
│  Scheduler         → Parses DAGs, triggers task instances    │
│  Executor          → Dispatches tasks to workers             │
│  Worker(s)         → Executes task instances                 │
│  Metadata Database → State of all DAGs, tasks, logs         │
│  DAGs Folder       → Python files defining workflows         │
│  Message Broker    → (Celery only) Redis or RabbitMQ         │
└──────────────────────────────────────────────────────────────┘
```

### Request Flow

```
1. User writes dag.py → places in ~/airflow/dags/
2. Scheduler parses dag.py every ~30 seconds
3. Scheduler checks: "is it time to trigger this DAG?"
4. Scheduler creates a DagRun, queues TaskInstances
5. Executor picks up TaskInstances → sends to workers
6. Workers run the task code
7. Task writes result to metadata DB + logs
8. Web server displays status in UI
```

---

## Real-World Use Cases

### 1. Daily ETL Pipeline
```
S3 (raw data) → Spark job → Data warehouse → dbt models → Tableau refresh
```

### 2. ML Model Retraining
```
Feature extraction → Train model → Evaluate → If accuracy > threshold → Deploy → Notify team
```

### 3. Data Quality Checks
```
Load yesterday's data → Run Great Expectations → If fail → alert → Else → publish
```

### 4. Multi-Cloud Data Sync
```
BigQuery export → GCS → Transfer to S3 → Load to Redshift → Update dashboard
```

### 5. Event-Driven Pipeline
```
Wait for file in S3 → Parse file → Enrich with API → Load to DB → Archive original
```

---

## Airflow Version History

| Version | Key Features |
|---|---|
| 1.x | Original release, Python 2/3, ZooKeeper scheduling |
| 2.0 (2020) | TaskFlow API, revamped scheduler, Python 3 only |
| 2.2 (2021) | Custom timetables, deferred operators |
| 2.3 (2022) | Dynamic task mapping |
| 2.4 (2022) | Dataset-driven scheduling |
| 2.5 (2022) | Dynamic task mapping improvements |
| 2.7 (2023) | Airflow Object Storage, UI improvements |
| 2.9 (2024) | Better dynamic task mapping, new UI |
| 3.0 (2025) | DAG versioning, Task Execution Interface, new UI |

This guide focuses on **Airflow 2.x** patterns, which are the production standard today.


---

# Chapter 2: Setup & Installation

## Quick Start with Docker Compose (Recommended)

The official Docker Compose setup runs the full Airflow stack with one command.

```bash
# Download the official docker-compose.yml
curl -LfO 'https://airflow.apache.org/docs/apache-airflow/stable/docker-compose.yaml'

# Create required directories
mkdir -p ./dags ./logs ./plugins ./config

# Set the UID (Linux/Mac)
echo -e "AIRFLOW_UID=$(id -u)" > .env

# Initialize the database
docker compose up airflow-init

# Start all services
docker compose up -d

# Access UI: http://localhost:8080
# Default credentials: airflow / airflow
```

### Services Started

| Service | Description | Port |
|---|---|---|
| `airflow-webserver` | Web UI + REST API | 8080 |
| `airflow-scheduler` | DAG parsing + triggering | — |
| `airflow-worker` | Task execution (Celery) | — |
| `airflow-triggerer` | Deferred operators | — |
| `postgres` | Metadata database | 5432 |
| `redis` | Celery broker | 6379 |

---

## Minimal Docker Compose (LocalExecutor)

For development without Celery:

```yaml
# docker-compose-local.yml
version: "3.8"

x-airflow-common: &airflow-common
  image: apache/airflow:2.9.0
  environment:
    AIRFLOW__CORE__EXECUTOR: LocalExecutor
    AIRFLOW__DATABASE__SQL_ALCHEMY_CONN: postgresql+psycopg2://airflow:airflow@postgres/airflow
    AIRFLOW__CORE__FERNET_KEY: ""
    AIRFLOW__CORE__DAGS_ARE_PAUSED_AT_CREATION: "true"
    AIRFLOW__CORE__LOAD_EXAMPLES: "false"
    AIRFLOW__API__AUTH_BACKENDS: "airflow.api.auth.backend.basic_auth"
    _PIP_ADDITIONAL_REQUIREMENTS: "pandas sqlalchemy"
  volumes:
    - ./dags:/opt/airflow/dags
    - ./logs:/opt/airflow/logs
    - ./plugins:/opt/airflow/plugins
  depends_on:
    postgres:
      condition: service_healthy

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: airflow
      POSTGRES_PASSWORD: airflow
      POSTGRES_DB: airflow
    healthcheck:
      test: ["CMD", "pg_isready", "-U", "airflow"]
      interval: 5s
      retries: 5

  airflow-webserver:
    <<: *airflow-common
    command: webserver
    ports: ["8080:8080"]
    healthcheck:
      test: ["CMD", "curl", "--fail", "http://localhost:8080/health"]
      interval: 10s

  airflow-scheduler:
    <<: *airflow-common
    command: scheduler

  airflow-init:
    <<: *airflow-common
    command: >
      bash -c "
        airflow db migrate &&
        airflow users create --username admin --password admin
          --firstname Admin --lastname User --role Admin --email admin@example.com
      "
```

```bash
docker compose -f docker-compose-local.yml up -d
```

---

## Astro CLI (Astronomer)

The Astronomer Astro CLI is the easiest local development experience:

```bash
# Install Astro CLI
curl -sSL install.astronomer.io | sudo bash -s

# Initialize a new Airflow project
astro dev init

# Start Airflow locally
astro dev start
# UI: http://localhost:8080 (admin/admin)

# Stop
astro dev stop

# Restart
astro dev restart

# Run a DAG
astro dev run dags trigger my_dag

# Deploy to Astronomer Cloud
astro deploy
```

Project structure created by `astro dev init`:
```
my_project/
  dags/           ← your DAG files
  plugins/        ← custom plugins
  include/        ← helper modules
  tests/          ← DAG tests
  requirements.txt ← Python packages
  Dockerfile      ← custom image
  .env            ← local env vars
```

---

## Local Installation (pip)

```bash
# Python 3.8+ required
python -m venv airflow-env
source airflow-env/bin/activate

# Install Airflow with constraints (prevents dependency conflicts)
AIRFLOW_VERSION=2.9.0
PYTHON_VERSION="$(python --version | cut -d " " -f 2 | cut -d "." -f 1-2)"
CONSTRAINT_URL="https://raw.githubusercontent.com/apache/airflow/constraints-${AIRFLOW_VERSION}/constraints-${PYTHON_VERSION}.txt"

pip install "apache-airflow==${AIRFLOW_VERSION}" --constraint "${CONSTRAINT_URL}"

# Install providers you need
pip install apache-airflow-providers-postgres
pip install apache-airflow-providers-amazon
pip install apache-airflow-providers-google
pip install apache-airflow-providers-slack

# Initialize database (SQLite by default)
export AIRFLOW_HOME=~/airflow
airflow db init

# Create admin user
airflow users create \
  --username admin \
  --firstname Admin \
  --lastname User \
  --role Admin \
  --email admin@example.com \
  --password admin

# Start services (2 terminals)
airflow webserver --port 8080
airflow scheduler
```

---

## Configuration (airflow.cfg)

Key settings in `~/airflow/airflow.cfg` or via environment variables:

```properties
# Environment variable format: AIRFLOW__{SECTION}__{KEY}

# Core
AIRFLOW__CORE__EXECUTOR=LocalExecutor          # or CeleryExecutor, KubernetesExecutor
AIRFLOW__CORE__DAGS_FOLDER=/opt/airflow/dags
AIRFLOW__CORE__LOAD_EXAMPLES=False
AIRFLOW__CORE__MAX_ACTIVE_RUNS_PER_DAG=3
AIRFLOW__CORE__MAX_ACTIVE_TASKS_PER_DAG=10
AIRFLOW__CORE__PARALLELISM=32                  # total concurrent tasks cluster-wide
AIRFLOW__CORE__DAGBAG_IMPORT_TIMEOUT=30        # seconds to import a DAG file

# Database
AIRFLOW__DATABASE__SQL_ALCHEMY_CONN=postgresql+psycopg2://user:pass@host/db

# Scheduler
AIRFLOW__SCHEDULER__MIN_FILE_PROCESS_INTERVAL=30   # DAG file re-parse interval
AIRFLOW__SCHEDULER__DAG_DIR_LIST_INTERVAL=30       # scan dags folder interval
AIRFLOW__SCHEDULER__MAX_TIS_PER_QUERY=512

# Webserver
AIRFLOW__WEBSERVER__BASE_URL=http://localhost:8080
AIRFLOW__WEBSERVER__SECRET_KEY=your-secret-key
AIRFLOW__WEBSERVER__WORKERS=4

# Email
AIRFLOW__SMTP__SMTP_HOST=smtp.gmail.com
AIRFLOW__SMTP__SMTP_PORT=587
AIRFLOW__SMTP__SMTP_USER=your@email.com
AIRFLOW__SMTP__SMTP_PASSWORD=app_password
AIRFLOW__SMTP__SMTP_MAIL_FROM=airflow@yourco.com
```

---

## First DAG: Verify Setup

```python
# dags/hello_world.py
from airflow import DAG
from airflow.operators.python import PythonOperator
from airflow.operators.bash import BashOperator
from datetime import datetime

def say_hello():
    print("Hello from Airflow!")
    return "success"

with DAG(
    dag_id      = "hello_world",
    start_date  = datetime(2024, 1, 1),
    schedule    = None,             # manual trigger only
    catchup     = False,
    tags        = ["example"],
) as dag:

    greet = PythonOperator(
        task_id          = "say_hello",
        python_callable  = say_hello,
    )

    echo = BashOperator(
        task_id      = "echo_date",
        bash_command = "echo 'Current date: $(date)'",
    )

    greet >> echo
```

```bash
# Trigger the DAG
airflow dags trigger hello_world

# Check status
airflow dags list
airflow tasks list hello_world
airflow dags state hello_world $(date +%Y-%m-%dT%H:%M:%S)
```


---

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


---

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


---

# Chapter 5: Operators

## PythonOperator

Runs any Python callable. The most flexible operator.

```python
from airflow.operators.python import PythonOperator

def process_data(ds: str, table: str, **context) -> dict:
    print(f"Processing {table} for {ds}")
    # ... business logic ...
    return {"rows_processed": 1000}

task = PythonOperator(
    task_id          = "process",
    python_callable  = process_data,
    op_kwargs        = {"table": "orders"},  # extra keyword args
    op_args          = [],                   # positional args
    provide_context  = True,                 # deprecated in 2.x (always provided)
)
```

---

## BashOperator

Runs shell commands. Supports Jinja templating.

```python
from airflow.operators.bash import BashOperator

run_script = BashOperator(
    task_id      = "run_etl",
    bash_command = """
        set -e
        python /opt/scripts/etl.py \
          --date {{ ds }} \
          --env {{ var.value.environment }}
        echo "Exit code: $?"
    """,
    env          = {"PYTHONPATH": "/opt/scripts"},
    cwd          = "/opt/scripts",
)

# Multi-line with heredoc
run_spark = BashOperator(
    task_id      = "spark_job",
    bash_command = """
        spark-submit \
          --master yarn \
          --deploy-mode cluster \
          --conf spark.executor.memory=4g \
          s3://bucket/jobs/transform.py \
          --date {{ ds }}
    """,
)
```

---

## EmailOperator

```python
from airflow.operators.email import EmailOperator

send_report = EmailOperator(
    task_id = "send_report",
    to      = ["team@company.com", "manager@company.com"],
    subject = "Daily ETL Complete: {{ ds }}",
    html_content = """
        <h2>Pipeline Status</h2>
        <p>The daily ETL for <b>{{ ds }}</b> completed successfully.</p>
        <ul>
            <li>Records processed: {{ ti.xcom_pull(task_ids='load') }}</li>
        </ul>
    """,
    files   = ["/tmp/report_{{ ds }}.csv"],
)
```

---

## EmptyOperator (DummyOperator)

No-op task — useful for grouping, branching joins, and marking pipeline stages.

```python
from airflow.operators.empty import EmptyOperator  # Airflow 2.4+
# or:
from airflow.operators.dummy import DummyOperator  # legacy

start = EmptyOperator(task_id="start")
end   = EmptyOperator(task_id="end", trigger_rule="all_done")

start >> [task_a, task_b, task_c] >> end
```

---

## SparkSubmitOperator

```python
from airflow.providers.apache.spark.operators.spark_submit import SparkSubmitOperator

spark_job = SparkSubmitOperator(
    task_id     = "spark_transform",
    application = "s3://bucket/jobs/transform.py",
    conn_id     = "spark_default",
    application_args = ["--date", "{{ ds }}", "--env", "prod"],
    conf = {
        "spark.executor.memory":        "8g",
        "spark.executor.cores":         "4",
        "spark.executor.instances":     "20",
        "spark.sql.shuffle.partitions": "400",
        "spark.sql.adaptive.enabled":   "true",
    },
    jars         = "s3://bucket/jars/delta-core.jar",
    py_files     = "s3://bucket/packages/utils.zip",
    driver_memory = "4g",
    num_executors = 20,
    verbose      = True,
)
```

---

## BigQueryOperator

```python
from airflow.providers.google.cloud.operators.bigquery import (
    BigQueryInsertJobOperator,
    BigQueryCreateEmptyTableOperator,
)

run_query = BigQueryInsertJobOperator(
    task_id    = "run_bq_query",
    gcp_conn_id = "google_cloud_default",
    configuration = {
        "query": {
            "query": """
                INSERT INTO `project.dataset.output_table`
                SELECT
                    user_id,
                    SUM(amount) AS total_spent,
                    COUNT(*)    AS order_count
                FROM `project.dataset.orders`
                WHERE DATE(created_at) = '{{ ds }}'
                GROUP BY user_id
            """,
            "useLegacySql": False,
            "destinationTable": {
                "projectId": "my-project",
                "datasetId": "analytics",
                "tableId":   "daily_summary",
            },
            "writeDisposition": "WRITE_APPEND",
        }
    },
)
```

---

## PostgresOperator

```python
from airflow.providers.postgres.operators.postgres import PostgresOperator

insert_data = PostgresOperator(
    task_id     = "insert_summary",
    postgres_conn_id = "my_postgres",
    sql         = """
        INSERT INTO daily_summary (date, total_orders, revenue)
        SELECT
            '{{ ds }}'::date,
            COUNT(*),
            SUM(amount)
        FROM orders
        WHERE order_date = '{{ ds }}'
        ON CONFLICT (date) DO UPDATE
            SET total_orders = EXCLUDED.total_orders,
                revenue      = EXCLUDED.revenue;
    """,
)
```

---

## KubernetesPodOperator

Run a Docker container in Kubernetes — isolates dependencies per task:

```python
from airflow.providers.cncf.kubernetes.operators.pod import KubernetesPodOperator
from kubernetes.client import models as k8s

run_ml_task = KubernetesPodOperator(
    task_id        = "train_model",
    name           = "train-model-{{ ds_nodash }}",
    namespace      = "airflow",
    image          = "my-registry/ml-pipeline:latest",
    cmds           = ["python", "train.py"],
    arguments      = ["--date", "{{ ds }}", "--model", "xgboost"],
    env_vars       = [
        k8s.V1EnvVar(name="MLFLOW_TRACKING_URI", value="http://mlflow:5000"),
    ],
    container_resources = k8s.V1ResourceRequirements(
        requests = {"cpu": "2", "memory": "8Gi"},
        limits   = {"cpu": "4", "memory": "16Gi"},
    ),
    is_delete_operator_pod = True,   # clean up pod after run
    get_logs       = True,
    log_events_on_failure = True,
    in_cluster     = True,
)
```

---

## HttpOperator

```python
from airflow.providers.http.operators.http import SimpleHttpOperator
import json

trigger_api = SimpleHttpOperator(
    task_id     = "trigger_refresh",
    http_conn_id = "my_api",
    endpoint    = "/api/v1/refresh",
    method      = "POST",
    data        = json.dumps({"date": "{{ ds }}", "force": True}),
    headers     = {"Content-Type": "application/json", "Authorization": "Bearer {{ var.value.api_token }}"},
    response_check = lambda response: response.json()["status"] == "accepted",
    log_response = True,
)
```

---

## TriggerDagRunOperator

```python
from airflow.operators.trigger_dagrun import TriggerDagRunOperator

trigger_downstream = TriggerDagRunOperator(
    task_id            = "trigger_reporting_dag",
    trigger_dag_id     = "daily_reporting",
    conf               = {"source_date": "{{ ds }}", "table": "orders"},
    wait_for_completion = True,
    execution_date     = "{{ data_interval_end }}",
    reset_dag_run      = True,
)
```

---

## Operator Execution Modes

| Mode | How | Notes |
|---|---|---|
| Default | Runs in worker process | Same Python environment |
| `KubernetesPodOperator` | New K8s pod | Full isolation |
| `DockerOperator` | New Docker container | Isolation without K8s |
| `ExternalTaskSensor` | Waits, not runs | Coordinates across DAGs |
| `BashOperator` subprocesses | Shell process | Can run any executable |


---

# Chapter 6: Sensors

## What are Sensors?

Sensors are operators that **poll** an external condition and block the pipeline until the condition is met. They are commonly used to:
- Wait for data files to arrive
- Wait for another pipeline to complete
- Wait for an API endpoint to be available
- Wait for a database row or result to appear

```
Sensor Execution Modes:

poke mode (default):
  Worker holds a slot while polling every poke_interval
  Uses a worker slot the entire time it's waiting
  Good for: short waits (< 5 minutes)

reschedule mode (recommended for long waits):
  Worker releases slot between checks
  Scheduler reschedules it at poke_interval
  Good for: long waits (hours)
  AIRFLOW__SCHEDULER__MIN_FILE_PROCESS_INTERVAL must be < poke_interval
```

---

## FileSensor

Wait for a file or directory to appear on disk:

```python
from airflow.sensors.filesystem import FileSensor

wait_for_csv = FileSensor(
    task_id       = "wait_for_input",
    filepath      = "/data/input/{{ ds }}/orders.csv",
    poke_interval = 60,      # check every 60 seconds
    timeout       = 7200,    # fail after 2 hours
    mode          = "reschedule",
    soft_fail     = False,   # True = skip instead of fail on timeout
)
```

---

## S3KeySensor

Wait for an object to appear in Amazon S3:

```python
from airflow.providers.amazon.aws.sensors.s3 import S3KeySensor

wait_for_data = S3KeySensor(
    task_id        = "wait_for_s3_data",
    bucket_name    = "my-data-bucket",
    bucket_key     = "raw/{{ ds }}/orders/*.parquet",
    wildcard_match = True,
    aws_conn_id    = "aws_default",
    poke_interval  = 120,
    timeout        = 3600,
    mode           = "reschedule",
)

# Check multiple keys
wait_for_all = S3KeySensor(
    task_id    = "wait_for_all_files",
    bucket_name = "my-bucket",
    bucket_key  = [
        "raw/{{ ds }}/orders.parquet",
        "raw/{{ ds }}/customers.parquet",
        "raw/{{ ds }}/products.parquet",
    ],
    aws_conn_id = "aws_default",
    mode        = "reschedule",
)
```

---

## HttpSensor

Wait for an HTTP endpoint to return a success response:

```python
from airflow.providers.http.sensors.http import HttpSensor

wait_for_api = HttpSensor(
    task_id         = "wait_for_api",
    http_conn_id    = "my_api",
    endpoint        = "/health",
    request_params  = {},
    response_check  = lambda response: response.json()["status"] == "healthy",
    poke_interval   = 30,
    timeout         = 300,
    mode            = "reschedule",
)

# Wait for a data endpoint
wait_for_report = HttpSensor(
    task_id         = "wait_for_daily_report",
    http_conn_id    = "reporting_api",
    endpoint        = "/api/reports/{{ ds }}",
    response_check  = lambda r: r.status_code == 200 and r.json()["ready"] is True,
    poke_interval   = 60,
    timeout         = 3600,
)
```

---

## SqlSensor

Wait for a SQL query to return rows:

```python
from airflow.providers.common.sql.sensors.sql import SqlSensor

wait_for_data = SqlSensor(
    task_id      = "wait_for_loaded_data",
    conn_id      = "my_postgres",
    sql          = """
        SELECT COUNT(*) FROM orders
        WHERE order_date = '{{ ds }}'
        HAVING COUNT(*) > 0
    """,
    poke_interval = 60,
    timeout       = 3600,
    mode          = "reschedule",
    success       = lambda rows: rows and rows[0][0] > 0,
)
```

---

## ExternalTaskSensor

Wait for a task in another DAG to complete:

```python
from airflow.sensors.external_task import ExternalTaskSensor

wait_for_upstream = ExternalTaskSensor(
    task_id             = "wait_for_source_dag",
    external_dag_id     = "source_data_pipeline",
    external_task_id    = "load_complete",   # None = wait for whole DAG
    execution_date_fn   = lambda dt: dt,     # map this DAG's run to source DAG's run
    mode                = "reschedule",
    timeout             = 3600,
    poke_interval       = 60,
    allowed_states      = ["success"],
    failed_states       = ["failed", "skipped"],
)

# Wait for yesterday's run of upstream DAG
from datetime import timedelta

wait_for_yesterday = ExternalTaskSensor(
    task_id          = "wait_for_yesterday",
    external_dag_id  = "daily_etl",
    execution_date_fn = lambda dt: dt - timedelta(days=1),
    mode             = "reschedule",
)
```

---

## PythonSensor

Wait for any custom Python condition:

```python
from airflow.sensors.python import PythonSensor

def check_kafka_topic(**context) -> bool:
    from confluent_kafka.admin import AdminClient
    admin = AdminClient({"bootstrap.servers": "kafka:9092"})
    metadata = admin.list_topics(timeout=5)
    topic = f"orders-{context['ds_nodash']}"
    return topic in metadata.topics

wait_for_kafka = PythonSensor(
    task_id         = "wait_for_kafka_topic",
    python_callable = check_kafka_topic,
    poke_interval   = 30,
    timeout         = 600,
    mode            = "reschedule",
)
```

---

## Deferred Sensors (Deferrable Operators)

Airflow 2.2+ supports async deferral — the sensor completely frees the worker while waiting, using a Triggerer process:

```python
from airflow.providers.amazon.aws.sensors.s3 import S3KeySensor

# Just add deferrable=True
wait_for_data = S3KeySensor(
    task_id     = "wait_for_s3",
    bucket_name = "my-bucket",
    bucket_key  = "data/{{ ds }}/input.csv",
    deferrable  = True,          # uses Triggerer, frees worker slot
    poke_interval = timedelta(seconds=60),
    timeout     = timedelta(hours=6),
)
```

**Benefits of deferrable:**
- Worker slot freed while waiting (much more scalable)
- Triggerer uses async I/O — handles thousands of deferred sensors with minimal resources
- Available for most built-in sensors (S3, GCS, BigQuery, etc.)

---

## Custom Sensor

```python
from airflow.sensors.base import BaseSensorOperator
from airflow.utils.decorators import apply_defaults

class DataAvailableSensor(BaseSensorOperator):
    """Waits until data is available in our internal API."""

    template_fields = ("endpoint",)

    def __init__(self, endpoint: str, min_records: int = 1, **kwargs):
        super().__init__(**kwargs)
        self.endpoint    = endpoint
        self.min_records = min_records

    def poke(self, context) -> bool:
        import requests
        response = requests.get(
            f"http://data-api/availability/{self.endpoint}",
            params={"date": context["ds"]},
            timeout=10,
        )
        if response.status_code != 200:
            return False
        data = response.json()
        available = data.get("record_count", 0) >= self.min_records
        self.log.info(f"Records available: {data.get('record_count', 0)}, needed: {self.min_records}")
        return available

# Usage
wait_for_orders = DataAvailableSensor(
    task_id     = "wait_for_orders",
    endpoint    = "orders",
    min_records = 100,
    poke_interval = 120,
    timeout     = 3600,
    mode        = "reschedule",
)
```

---

## Sensor Best Practices

```
1. Always use mode="reschedule" for waits > 5 minutes
   → Prevents worker slot exhaustion

2. Set reasonable timeouts
   → soft_fail=True: skip downstream instead of fail

3. Use deferrable=True when available (Airflow 2.2+)
   → Most scalable option for long-running sensors

4. Avoid tight poke_interval (< 30s)
   → Creates unnecessary load on external systems

5. Set failed_states in ExternalTaskSensor
   → Otherwise waits forever if upstream fails
```


---

# Chapter 7: Hooks & Connections

## Connections

Connections store credentials for external systems. They are the single source of truth for how Airflow connects to databases, cloud services, and APIs.

### Managing Connections

**Via UI:** Admin → Connections → "+"

**Via CLI:**
```bash
# Add a PostgreSQL connection
airflow connections add \
  --conn-id       my_postgres \
  --conn-type     postgres \
  --conn-host     localhost \
  --conn-port     5432 \
  --conn-login    myuser \
  --conn-password mypassword \
  --conn-schema   mydb \
  --conn-extra    '{"sslmode": "require"}'

# Add an S3 connection (AWS)
airflow connections add \
  --conn-id   aws_default \
  --conn-type aws \
  --conn-extra '{"aws_access_key_id": "AKIA...", "aws_secret_access_key": "..."}'

# List all connections
airflow connections list

# Delete
airflow connections delete my_old_conn
```

**Via environment variables** (best for containers):
```bash
# URI format
export AIRFLOW_CONN_MY_POSTGRES="postgresql://myuser:mypassword@localhost:5432/mydb"

# JSON format
export AIRFLOW_CONN_MY_POSTGRES='{
  "conn_type": "postgres",
  "host": "localhost",
  "login": "myuser",
  "password": "mypassword",
  "schema": "mydb",
  "port": 5432
}'
```

---

## Common Connection Types

| Conn Type | Provider Package | Use Case |
|---|---|---|
| `postgres` | apache-airflow-providers-postgres | PostgreSQL |
| `mysql` | apache-airflow-providers-mysql | MySQL |
| `aws` | apache-airflow-providers-amazon | AWS (S3, Redshift, etc.) |
| `google_cloud_platform` | apache-airflow-providers-google | GCP (BigQuery, GCS, etc.) |
| `azure` | apache-airflow-providers-microsoft-azure | Azure services |
| `spark` | apache-airflow-providers-apache-spark | Spark cluster |
| `kafka` | apache-airflow-providers-apache-kafka | Kafka |
| `http` | apache-airflow-providers-http | HTTP/REST APIs |
| `ssh` | apache-airflow-providers-ssh | SSH connections |
| `slack` | apache-airflow-providers-slack | Slack webhooks |
| `snowflake` | apache-airflow-providers-snowflake | Snowflake DWH |
| `databricks` | apache-airflow-providers-databricks | Databricks |

---

## Hooks

Hooks are the Python interface to external services. They read a `conn_id` and manage the connection lifecycle.

### PostgresHook

```python
from airflow.providers.postgres.hooks.postgres import PostgresHook
import pandas as pd

hook = PostgresHook(postgres_conn_id="my_postgres")

# Run SQL
hook.run("INSERT INTO logs VALUES (%s, %s)", parameters=["2024-01-01", "success"])

# Get records as list of tuples
records = hook.get_records("SELECT id, name FROM users LIMIT 100")

# Get first record
first = hook.get_first("SELECT MAX(created_at) FROM orders")

# Get Pandas DataFrame
df = hook.get_pandas_df("SELECT * FROM orders WHERE date = '{{ ds }}'")

# Get SQLAlchemy engine
engine = hook.get_sqlalchemy_engine()
df.to_sql("staging", engine, if_exists="append", index=False)

# Execute in PythonOperator
def load_to_postgres(**context):
    hook = PostgresHook(postgres_conn_id="my_postgres")
    df = pd.DataFrame(...)
    hook.insert_rows(
        table  = "daily_summary",
        rows   = df.to_records(index=False).tolist(),
        target_fields = df.columns.tolist(),
        replace = True,
        replace_index = ["date"],
    )
```

### S3Hook

```python
from airflow.providers.amazon.aws.hooks.s3 import S3Hook

hook = S3Hook(aws_conn_id="aws_default")

# Check if key exists
exists = hook.check_for_key("raw/2024-01-01/data.csv", "my-bucket")

# Download file
hook.download_file(
    key      = "raw/{{ ds }}/data.csv",
    bucket_name = "my-bucket",
    local_path = "/tmp/",
)

# Upload file
hook.load_file(
    filename    = "/tmp/output.parquet",
    key         = "processed/{{ ds }}/output.parquet",
    bucket_name = "my-bucket",
    replace     = True,
)

# List keys
keys = hook.list_keys(bucket_name="my-bucket", prefix="raw/2024-01-01/")

# Read file content
content = hook.read_key("config/settings.json", "my-bucket")
```

### HttpHook

```python
from airflow.providers.http.hooks.http import HttpHook

hook = HttpHook(http_conn_id="my_api", method="POST")

response = hook.run(
    endpoint = "/api/v1/process",
    data     = {"date": "2024-01-01"},
    headers  = {"Content-Type": "application/json"},
)
data = response.json()
```

### SlackHook

```python
from airflow.providers.slack.hooks.slack_webhook import SlackWebhookHook

def notify_slack(message: str) -> None:
    hook = SlackWebhookHook(slack_webhook_conn_id="slack_webhook")
    hook.send(text=message, icon_emoji=":airflow:")

# In a callback
def on_failure_callback(context):
    dag_id  = context["dag"].dag_id
    task_id = context["task"].task_id
    ds      = context["ds"]
    notify_slack(f":red_circle: Task failed: `{dag_id}.{task_id}` for `{ds}`")
```

---

## Secrets Backends

Instead of storing credentials in the Airflow metadata database, use an external secrets manager.

### AWS Secrets Manager

```properties
# airflow.cfg
[secrets]
backend = airflow.providers.amazon.aws.secrets.secrets_manager.SecretsManagerBackend
backend_kwargs = {
    "connections_prefix": "airflow/connections",
    "variables_prefix": "airflow/variables",
    "config_prefix": "airflow/config"
}
```

```bash
# Store connection in AWS Secrets Manager
aws secretsmanager create-secret \
  --name "airflow/connections/my_postgres" \
  --secret-string "postgresql://user:pass@host:5432/db"
```

### HashiCorp Vault

```properties
[secrets]
backend = airflow.providers.hashicorp.secrets.vault.VaultBackend
backend_kwargs = {
    "connections_path": "connections",
    "variables_path": "variables",
    "mount_point": "airflow",
    "url": "http://vault:8200",
    "token": "vault-token"
}
```

---

## Custom Hook

```python
from airflow.hooks.base import BaseHook
import requests

class MyAPIHook(BaseHook):
    conn_name_attr  = "my_api_conn_id"
    default_conn_name = "my_api_default"

    def __init__(self, my_api_conn_id: str = "my_api_default"):
        super().__init__()
        self.my_api_conn_id = my_api_conn_id
        self._session = None

    def get_conn(self):
        if self._session is None:
            conn = self.get_connection(self.my_api_conn_id)
            self._session = requests.Session()
            self._session.headers.update({
                "Authorization": f"Bearer {conn.password}",
                "Content-Type":  "application/json",
            })
            self._base_url = f"http://{conn.host}:{conn.port}"
        return self._session

    def get_data(self, endpoint: str, params: dict = None) -> dict:
        session = self.get_conn()
        response = session.get(f"{self._base_url}{endpoint}", params=params)
        response.raise_for_status()
        return response.json()

    def post_data(self, endpoint: str, payload: dict) -> dict:
        session = self.get_conn()
        response = session.post(f"{self._base_url}{endpoint}", json=payload)
        response.raise_for_status()
        return response.json()
```


---

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


---

# Chapter 9: Scheduling & Triggers

## Cron-Based Scheduling

```python
with DAG(
    dag_id   = "my_dag",
    schedule = "0 6 * * *",   # daily at 6 AM UTC
    ...
)
```

### Cron Syntax

```
┌───────────── minute (0–59)
│ ┌─────────── hour (0–23)
│ │ ┌───────── day of month (1–31)
│ │ │ ┌─────── month (1–12)
│ │ │ │ ┌───── day of week (0=Sun–6=Sat)
│ │ │ │ │
* * * * *

0 6 * * *       Every day at 6 AM
0 */4 * * *     Every 4 hours
0 6 * * 1-5     Weekdays at 6 AM
30 2 1 * *      1st of every month at 2:30 AM
0 0 * * 0       Every Sunday midnight
*/15 * * * *    Every 15 minutes
0 0 1 1 *       Once a year (Jan 1st)
```

### Preset Schedules

| Preset | Equivalent Cron | When |
|---|---|---|
| `@once` | (run once) | Single run at start_date |
| `@hourly` | `0 * * * *` | Every hour |
| `@daily` | `0 0 * * *` | Midnight daily |
| `@weekly` | `0 0 * * 0` | Sunday midnight |
| `@monthly` | `0 0 1 * *` | 1st of month |
| `@yearly` | `0 0 1 1 *` | January 1st |
| `None` | — | Manual trigger only |

---

## timedelta Scheduling

```python
from datetime import timedelta

# Run every 6 hours
dag = DAG(schedule=timedelta(hours=6), ...)

# Run every 30 minutes
dag = DAG(schedule=timedelta(minutes=30), ...)
```

---

## Data Interval

In Airflow 2.2+, every scheduled run has a **data interval** — the period of data it processes:

```
Schedule: @daily

DAG Run triggered after 2024-01-03 00:00:00:
  data_interval_start = 2024-01-02 00:00:00
  data_interval_end   = 2024-01-03 00:00:00

The run processes data FROM Jan 2 TO Jan 3.
{{ ds }}            = "2024-01-02"   (data_interval_start date)
{{ data_interval_end | ds }} = "2024-01-03"
```

---

## Catchup

`catchup=True` (default): Airflow creates a DAG run for every missed interval since `start_date`.

```python
# Dangerous if start_date is far in the past!
dag = DAG(
    start_date = datetime(2020, 1, 1),
    schedule   = "@daily",
    catchup    = True,   # creates 4+ years of runs!
)

# Production best practice: disable catchup
dag = DAG(
    start_date = datetime(2024, 1, 1),
    schedule   = "@daily",
    catchup    = False,
)
```

---

## Backfill

Manually run a DAG for historical intervals:

```bash
# Backfill one week of data
airflow dags backfill \
  --start-date 2024-01-01 \
  --end-date   2024-01-07 \
  my_dag

# With dry run (don't actually run)
airflow dags backfill --dry-run \
  --start-date 2024-01-01 --end-date 2024-01-07 \
  my_dag

# Mark as successful (skip actual execution)
airflow dags backfill --mark-success \
  --start-date 2024-01-01 --end-date 2024-01-07 \
  my_dag
```

---

## Custom Timetables (Airflow 2.2+)

For non-standard schedules (e.g., business days only, every weekday at market open):

```python
from airflow.timetables.base import DagRunInfo, DataInterval, Timetable
from pendulum import DateTime, timezone
import pendulum

class WeekdayTimetable(Timetable):
    """Run once per weekday (Mon–Fri)."""

    description = "Runs on weekdays at 6 AM UTC"

    def infer_manual_data_interval(self, run_after: DateTime) -> DataInterval:
        start = run_after.start_of("day")
        return DataInterval(start=start, end=start.add(days=1))

    def next_dagrun_info(self, *, last_automated_data_interval, restriction):
        if last_automated_data_interval is None:
            next_start = restriction.earliest
        else:
            next_start = last_automated_data_interval.end

        # Skip weekends
        while next_start.day_of_week >= 5:  # Sat=5, Sun=6
            next_start = next_start.add(days=1)

        run_after = next_start.add(hours=6)  # 6 AM
        return DagRunInfo.interval(
            start=next_start,
            end=next_start.add(days=1),
        )
```

Register in a plugin and use:
```python
dag = DAG(schedule=WeekdayTimetable(), ...)
```

---

## Dataset-Driven Scheduling (Airflow 2.4+)

DAGs can be triggered when other DAGs produce datasets, enabling event-driven pipelines:

```python
from airflow import DAG, Dataset

# Define datasets (URIs, logical identifiers)
orders_dataset   = Dataset("s3://bucket/raw/orders/")
products_dataset = Dataset("s3://bucket/raw/products/")
output_dataset   = Dataset("s3://bucket/processed/sales/")

# Producer DAG: declares that it updates a dataset
with DAG("load_orders", schedule="@daily") as dag:
    @task(outlets=[orders_dataset])   # this task produces orders_dataset
    def load_orders():
        # ... load orders to S3 ...
        pass

# Consumer DAG: triggered when BOTH datasets are updated
with DAG(
    "process_sales",
    schedule = [orders_dataset, products_dataset],  # triggered by datasets
) as dag:
    @task(outlets=[output_dataset])
    def process(ds=None):
        # Runs only after both orders and products have been updated
        pass
```

---

## Manual Triggers

```bash
# Trigger via CLI
airflow dags trigger my_dag

# With configuration
airflow dags trigger my_dag --conf '{"env": "production", "date": "2024-01-15"}'

# With a specific run_id
airflow dags trigger my_dag --run-id "manual__2024-01-15T10:00:00"
```

```python
# Trigger via REST API
import requests

response = requests.post(
    "http://airflow:8080/api/v1/dags/my_dag/dagRuns",
    auth    = ("admin", "admin"),
    json    = {
        "conf":     {"env": "production", "date": "2024-01-15"},
        "run_id":   "api_trigger__2024-01-15",
    },
)
print(response.json())
```

---

## Schedule Summary

| Pattern | Use Case | Example |
|---|---|---|
| Cron string | Recurring schedule | `"0 2 * * *"` |
| `@daily`, `@hourly` | Common presets | `schedule="@daily"` |
| `timedelta` | Fixed intervals | `schedule=timedelta(hours=4)` |
| `None` | Manual only | `schedule=None` |
| Dataset | Event-driven | `schedule=[my_dataset]` |
| Custom timetable | Complex schedules | Business days, market hours |


---

# Chapter 10: Executors

## What is an Executor?

The executor determines **how and where tasks run**. It receives task instances from the scheduler and dispatches them to workers.

```
Scheduler ──► Executor ──► Worker(s) ──► Run task
```

| Executor | Workers | Use Case |
|---|---|---|
| `SequentialExecutor` | None (in-process) | Testing only |
| `LocalExecutor` | Subprocesses (same machine) | Single-node, small workloads |
| `CeleryExecutor` | Remote worker pool | Multi-node, production |
| `KubernetesExecutor` | Kubernetes pods | Cloud-native, strong isolation |
| `CeleryKubernetesExecutor` | Both Celery + K8s | Hybrid workloads |
| `LocalKubernetesExecutor` | Local + K8s | Local light + K8s heavy |

---

## SequentialExecutor

Runs one task at a time in the scheduler process. Only for local development.

```properties
AIRFLOW__CORE__EXECUTOR=SequentialExecutor
```

---

## LocalExecutor

Spawns subprocesses on the same machine as the scheduler. Simple and fast for single-node setups.

```properties
AIRFLOW__CORE__EXECUTOR=LocalExecutor
AIRFLOW__DATABASE__SQL_ALCHEMY_CONN=postgresql+psycopg2://airflow:airflow@postgres/airflow
```

```yaml
# docker-compose (simplified)
scheduler:
  image: apache/airflow:2.9.0
  command: scheduler
  environment:
    AIRFLOW__CORE__EXECUTOR: LocalExecutor
```

**Max parallelism:** controlled by `AIRFLOW__CORE__PARALLELISM` (default 32).

---

## CeleryExecutor

Distributes tasks across a pool of worker nodes using Celery. Requires a message broker (Redis or RabbitMQ).

```
Scheduler ──► Redis/RabbitMQ ──► Celery Workers (N nodes)
                                       │
                                 run task instance
```

```properties
AIRFLOW__CORE__EXECUTOR=CeleryExecutor
AIRFLOW__CELERY__BROKER_URL=redis://:password@redis:6379/0
AIRFLOW__CELERY__RESULT_BACKEND=db+postgresql://airflow:airflow@postgres/airflow
AIRFLOW__CELERY__WORKER_CONCURRENCY=16  # tasks per worker
```

```yaml
# docker-compose snippet
worker:
  image: apache/airflow:2.9.0
  command: celery worker
  environment:
    AIRFLOW__CORE__EXECUTOR: CeleryExecutor
    AIRFLOW__CELERY__BROKER_URL: redis://redis:6379/0
  deploy:
    replicas: 3  # scale workers
```

### Worker Queues

Route specific tasks to specific workers:

```python
# Task uses "gpu" queue
train_model = PythonOperator(
    task_id = "train_model",
    queue   = "gpu",           # route to GPU worker
    ...
)
```

```bash
# Start a worker that handles the "gpu" queue
airflow celery worker -q gpu,default
```

---

## KubernetesExecutor

Creates a new Kubernetes pod for every task instance. Full isolation, no persistent workers.

```
Scheduler ──► Kubernetes API ──► Pod (task runs) ──► Pod deleted
```

```properties
AIRFLOW__CORE__EXECUTOR=KubernetesExecutor
AIRFLOW__KUBERNETES__NAMESPACE=airflow
AIRFLOW__KUBERNETES__WORKER_CONTAINER_REPOSITORY=apache/airflow
AIRFLOW__KUBERNETES__WORKER_CONTAINER_TAG=2.9.0
AIRFLOW__KUBERNETES__IN_CLUSTER=True
```

### Pod Configuration Per Task

```python
from airflow.providers.cncf.kubernetes.operators.pod import KubernetesPodOperator
from airflow.kubernetes.secret import Secret
from kubernetes.client import models as k8s

# Per-task pod override
task = PythonOperator(
    task_id = "heavy_task",
    executor_config = {
        "pod_override": k8s.V1Pod(
            spec=k8s.V1PodSpec(
                containers=[k8s.V1Container(
                    name      = "base",
                    resources = k8s.V1ResourceRequirements(
                        requests = {"cpu": "4", "memory": "16Gi"},
                        limits   = {"cpu": "8", "memory": "32Gi"},
                    ),
                    env = [
                        k8s.V1EnvVar(name="EXTRA_VAR", value="hello"),
                    ],
                )],
                node_selector = {"cloud.google.com/gke-nodepool": "high-memory"},
            )
        )
    },
)
```

---

## Concurrency Settings

```properties
# Total concurrent tasks across all DAGs
AIRFLOW__CORE__PARALLELISM=32

# Concurrent tasks per DAG
AIRFLOW__CORE__MAX_ACTIVE_TASKS_PER_DAG=16

# Concurrent DAG runs per DAG
AIRFLOW__CORE__MAX_ACTIVE_RUNS_PER_DAG=3

# CeleryExecutor worker concurrency
AIRFLOW__CELERY__WORKER_CONCURRENCY=16
```

```python
# Per-DAG overrides
with DAG(
    dag_id              = "my_dag",
    max_active_runs     = 2,
    max_active_tasks    = 8,
    concurrency         = 8,   # legacy: max tasks across all runs of this DAG
) as dag:
    pass
```

---

## Task Queues and Pools

### Pools

Pools limit the number of concurrently running tasks that use a shared resource (e.g., a slow DB that can't handle more than 5 connections):

```bash
# Create a pool
airflow pools set db_pool 5 "Limit concurrent DB tasks"
airflow pools list
```

```python
task = PythonOperator(
    task_id = "query_db",
    pool    = "db_pool",    # limited to 5 concurrent tasks
    pool_slots = 1,         # each task takes 1 slot (default)
)

# Heavy tasks take more slots
heavy_task = PythonOperator(
    task_id    = "full_table_scan",
    pool       = "db_pool",
    pool_slots = 3,   # this task alone uses 3 of 5 db_pool slots
)
```

### Priority Weight

```python
# Higher priority_weight = runs first when slots are contested
urgent_task = PythonOperator(
    task_id         = "urgent_load",
    priority_weight = 100,    # default is 1
    weight_rule     = "downstream",  # or "upstream", "absolute"
)
```

---

## Choosing an Executor

```
SequentialExecutor:
  ✓ Testing DAG logic
  ✗ Never in production

LocalExecutor:
  ✓ Single machine, moderate workloads
  ✓ Simpler setup
  ✗ No horizontal scaling
  ✗ Worker failure = scheduler failure

CeleryExecutor:
  ✓ Multi-node, horizontal scaling
  ✓ Worker pools + queues
  ✗ Requires Redis/RabbitMQ
  ✗ Workers share dependencies

KubernetesExecutor:
  ✓ Full isolation (each task = new pod)
  ✓ Custom resources per task
  ✓ Scales to zero (no idle workers)
  ✗ Pod startup overhead (~10-30s per task)
  ✗ Complex setup
  ✓ Best for: cloud-native, varying resource needs
```


---

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


---

# Chapter 12: Testing DAGs

## Why Test DAGs?

```
Parse errors         → DAG won't load at all
Import errors        → Silent failures in production
Wrong dependencies   → Tasks run in wrong order
Broken logic         → Incorrect data silently written
Missing connections  → Task fails at runtime, not at review time
```

---

## DAG Validation Tests

The fastest, cheapest tests — just import the DAG file:

```python
# tests/test_dag_validation.py
import os
import pytest
from airflow.models import DagBag

DAG_FOLDER = os.path.join(os.path.dirname(__file__), "../dags")

@pytest.fixture(scope="module")
def dagbag():
    return DagBag(dag_folder=DAG_FOLDER, include_examples=False)

def test_no_import_errors(dagbag):
    assert dagbag.import_errors == {}, f"DAG import errors: {dagbag.import_errors}"

def test_all_dags_loaded(dagbag):
    assert len(dagbag.dags) > 0, "No DAGs found"

def test_dag_ids_are_unique(dagbag):
    dag_ids = list(dagbag.dags.keys())
    assert len(dag_ids) == len(set(dag_ids))

@pytest.mark.parametrize("dag_id", ["my_etl_dag", "daily_report_dag"])
def test_specific_dag_exists(dagbag, dag_id):
    assert dag_id in dagbag.dags, f"DAG '{dag_id}' not found"
```

---

## Task Structure Tests

```python
# tests/test_my_etl_dag.py
import pytest
from airflow.models import DagBag

@pytest.fixture(scope="module")
def dag():
    dagbag = DagBag(dag_folder="dags/", include_examples=False)
    return dagbag.get_dag("my_etl_dag")

def test_dag_loaded(dag):
    assert dag is not None

def test_task_count(dag):
    assert len(dag.tasks) == 4

def test_task_ids(dag):
    task_ids = {t.task_id for t in dag.tasks}
    assert task_ids == {"extract", "transform", "load", "notify"}

def test_dependencies(dag):
    extract   = dag.get_task("extract")
    transform = dag.get_task("transform")
    load      = dag.get_task("load")
    notify    = dag.get_task("notify")

    assert transform in extract.downstream_list
    assert load      in transform.downstream_list
    assert notify    in load.downstream_list

def test_no_catchup(dag):
    assert dag.catchup is False

def test_retry_config(dag):
    extract = dag.get_task("extract")
    assert extract.retries >= 1
```

---

## Unit Testing Task Logic

Test the Python callable directly — no Airflow context needed:

```python
# dags/my_etl_dag.py
def extract_records(date: str) -> list[dict]:
    """Extract orders from DB for a given date."""
    # In production this hits a real DB
    return []  # stub

def transform_records(records: list[dict]) -> list[dict]:
    return [
        {**r, "amount_usd": r["amount"] / 100}
        for r in records
        if r.get("status") == "completed"
    ]

def load_records(records: list[dict], table: str) -> int:
    return len(records)
```

```python
# tests/test_task_logic.py
from dags.my_etl_dag import transform_records, load_records

def test_transform_filters_non_completed():
    records = [
        {"id": 1, "amount": 1000, "status": "completed"},
        {"id": 2, "amount": 500,  "status": "pending"},
        {"id": 3, "amount": 200,  "status": "completed"},
    ]
    result = transform_records(records)
    assert len(result) == 2
    assert all(r["status"] == "completed" for r in result)

def test_transform_converts_amount():
    records = [{"id": 1, "amount": 1000, "status": "completed"}]
    result = transform_records(records)
    assert result[0]["amount_usd"] == 10.0

def test_transform_empty_input():
    assert transform_records([]) == []

def test_load_returns_count():
    records = [{"id": i} for i in range(5)]
    assert load_records(records, "orders") == 5
```

---

## Testing with Mocked Context

```python
# dags/report_dag.py
from airflow.operators.python import PythonOperator

def send_report(**context):
    ds         = context["ds"]
    dag_run    = context["dag_run"]
    task_inst  = context["ti"]

    path = task_inst.xcom_pull(task_ids="generate", key="report_path")
    # send report at path for ds
    return f"Sent report for {ds}"
```

```python
# tests/test_report_task.py
from unittest.mock import MagicMock, patch
from dags.report_dag import send_report

def test_send_report_uses_correct_ds():
    ti = MagicMock()
    ti.xcom_pull.return_value = "s3://bucket/reports/2024-01-01/report.pdf"

    context = {
        "ds":      "2024-01-01",
        "dag_run": MagicMock(),
        "ti":      ti,
    }

    result = send_report(**context)

    ti.xcom_pull.assert_called_once_with(task_ids="generate", key="report_path")
    assert "2024-01-01" in result
```

---

## Testing with pytest + airflow DB

For tests that need the metadata DB (XCom, task state):

```python
# conftest.py
import pytest
from airflow.models import DagBag, TaskInstance
from airflow.utils.state import State
from airflow.utils.types import DagRunType
from datetime import datetime

@pytest.fixture(scope="session")
def airflow_db():
    """Initialize Airflow DB for tests."""
    from airflow.utils.db import initdb
    initdb()

@pytest.fixture
def create_dag_run(airflow_db):
    from airflow.models import DagRun
    def _create(dag, execution_date=None):
        execution_date = execution_date or datetime(2024, 1, 1)
        dag_run = dag.create_dagrun(
            run_id          = f"test__{execution_date}",
            execution_date  = execution_date,
            state           = State.RUNNING,
            run_type        = DagRunType.MANUAL,
        )
        return dag_run
    return _create
```

```python
# tests/test_xcom_flow.py
from airflow.models import DagBag, TaskInstance
from datetime import datetime

def test_extract_pushes_xcom(create_dag_run):
    dagbag   = DagBag(dag_folder="dags/", include_examples=False)
    dag      = dagbag.get_dag("xcom_test_dag")
    dag_run  = create_dag_run(dag)

    ti = TaskInstance(task=dag.get_task("extract"), run_id=dag_run.run_id)
    ti.run(ignore_ti_state=True)

    value = ti.xcom_pull(task_ids="extract", key="row_count")
    assert value == 100
```

---

## Integration Testing with Docker

```python
# tests/test_integration.py
import subprocess
import pytest

@pytest.fixture(scope="session")
def airflow_services():
    """Start Airflow via docker-compose for integration tests."""
    subprocess.run(
        ["docker-compose", "-f", "docker-compose.test.yml", "up", "-d"],
        check=True,
    )
    # Wait for webserver to be ready
    import time; time.sleep(30)
    yield
    subprocess.run(
        ["docker-compose", "-f", "docker-compose.test.yml", "down", "-v"],
    )

def test_dag_trigger_via_api(airflow_services):
    import requests
    response = requests.post(
        "http://localhost:8080/api/v1/dags/my_etl_dag/dagRuns",
        auth = ("admin", "admin"),
        json = {"conf": {}, "logical_date": "2024-01-01T00:00:00Z"},
    )
    assert response.status_code == 200
    dag_run_id = response.json()["dag_run_id"]
    assert dag_run_id is not None
```

---

## Testing Dynamic Task Mapping

```python
# tests/test_dynamic_dag.py
from airflow.models import DagBag

def test_dynamic_mapping_tasks():
    dagbag = DagBag(dag_folder="dags/", include_examples=False)
    dag    = dagbag.get_dag("dynamic_etl_dag")

    # Verify the mapped task exists
    task_ids = {t.task_id for t in dag.tasks}
    assert "process_file" in task_ids

    # Verify it's a mapped operator
    process_task = dag.get_task("process_file")
    assert hasattr(process_task, "expand_input")
```

---

## CI/CD Pipeline for DAGs

```yaml
# .github/workflows/test-dags.yml
name: Test DAGs

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    env:
      AIRFLOW_HOME: /tmp/airflow
      AIRFLOW__CORE__LOAD_EXAMPLES: "False"
      AIRFLOW__DATABASE__SQL_ALCHEMY_CONN: "sqlite:////tmp/airflow/airflow.db"

    steps:
      - uses: actions/checkout@v3

      - name: Set up Python
        uses: actions/setup-python@v4
        with:
          python-version: "3.11"

      - name: Install dependencies
        run: |
          pip install apache-airflow==2.9.0 pytest pytest-cov

      - name: Initialize Airflow DB
        run: airflow db init

      - name: Run tests
        run: |
          pytest tests/ -v --cov=dags --cov-report=xml

      - name: Upload coverage
        uses: codecov/codecov-action@v3
```

---

## Test Project Structure

```
project/
├── dags/
│   ├── my_etl_dag.py
│   ├── report_dag.py
│   └── dynamic_dag.py
├── tests/
│   ├── conftest.py
│   ├── test_dag_validation.py    # Import and structure checks
│   ├── test_task_logic.py        # Pure function unit tests
│   ├── test_my_etl_dag.py        # DAG-level assertions
│   └── test_integration.py       # End-to-end (optional)
├── pytest.ini
└── requirements-test.txt
```

```ini
# pytest.ini
[pytest]
testpaths = tests
log_cli = true
log_level = INFO
```

---

## Testing Checklist

```
DAG validation:
  ☐ No import errors (DagBag.import_errors is empty)
  ☐ All expected DAGs are present
  ☐ No duplicate task IDs
  ☐ Correct dependencies (upstream/downstream)
  ☐ catchup=False on production DAGs

Task logic:
  ☐ Happy path tested
  ☐ Empty input handled
  ☐ Edge cases covered
  ☐ Side effects mocked (DB, S3, HTTP)

Configuration:
  ☐ Retries configured
  ☐ Timeouts configured
  ☐ Correct schedule
  ☐ Correct start_date
```


---

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


---

# Chapter 14: Monitoring & Alerting

## Airflow Web UI Overview

```
Main views:
  DAGs list     → All DAGs with status, schedule, last run, next run
  Grid view     → Task run history per DAG (rows=tasks, cols=runs)
  Graph view    → DAG dependency visualization
  Gantt chart   → Task duration visualization per run
  Calendar view → Run success/failure per day
  Task Logs     → Per-task stdout + airflow logs
  Code view     → DAG source code snapshot
  Audit Log     → Who changed what
```

---

## Callbacks

Callbacks are Python functions invoked at specific lifecycle events.

### DAG-Level Callbacks

```python
from airflow import DAG
from datetime import datetime

def on_dag_success(context):
    dag_id  = context["dag"].dag_id
    run_id  = context["run_id"]
    print(f"DAG {dag_id} succeeded: {run_id}")

def on_dag_failure(context):
    dag_id     = context["dag"].dag_id
    exception  = context.get("exception")
    print(f"DAG {dag_id} failed: {exception}")

with DAG(
    dag_id               = "monitored_dag",
    start_date           = datetime(2024, 1, 1),
    schedule             = "@daily",
    catchup              = False,
    on_success_callback  = on_dag_success,
    on_failure_callback  = on_dag_failure,
) as dag:
    pass
```

### Task-Level Callbacks

```python
from airflow.operators.python import PythonOperator

def on_task_failure(context):
    task_id  = context["task"].task_id
    dag_id   = context["dag"].dag_id
    ds       = context["ds"]
    exc      = context.get("exception")
    print(f"Task {dag_id}.{task_id} failed on {ds}: {exc}")

def on_task_retry(context):
    ti        = context["task_instance"]
    attempt   = ti.try_number
    print(f"Retry #{attempt} for {ti.task_id}")

def on_task_success(context):
    print("Task succeeded!")

task = PythonOperator(
    task_id             = "my_task",
    python_callable     = lambda: None,
    on_failure_callback = on_task_failure,
    on_retry_callback   = on_task_retry,
    on_success_callback = on_task_success,
)
```

---

## Slack Alerting

```python
from airflow.providers.slack.hooks.slack_webhook import SlackWebhookHook

def slack_alert(context, color: str = "#ff0000", title: str = "Task Failed"):
    ti      = context["task_instance"]
    dag_id  = context["dag"].dag_id
    task_id = ti.task_id
    ds      = context["ds"]
    log_url = ti.log_url
    exc     = context.get("exception", "Unknown error")

    hook = SlackWebhookHook(slack_webhook_conn_id="slack_webhook")
    hook.send(
        attachments=[{
            "color":  color,
            "title":  f"{title} | {dag_id}.{task_id}",
            "fields": [
                {"title": "Date",    "value": ds,      "short": True},
                {"title": "Error",   "value": str(exc), "short": False},
            ],
            "actions": [{
                "type": "button",
                "text": "View Logs",
                "url":  log_url,
            }],
        }]
    )

def slack_failure_callback(context):
    slack_alert(context, color="#ff0000", title=":red_circle: Task Failed")

def slack_success_callback(context):
    slack_alert(context, color="#36a64f", title=":white_check_mark: Task Succeeded")

# Apply as default to all tasks in a DAG
default_args = {
    "on_failure_callback": slack_failure_callback,
    "retries":             2,
}

with DAG("my_dag", default_args=default_args, ...) as dag:
    pass
```

---

## Email Alerting

```python
# airflow.cfg
# [smtp]
# smtp_host = smtp.gmail.com
# smtp_starttls = True
# smtp_ssl = False
# smtp_user = airflow@example.com
# smtp_password = yourpassword
# smtp_port = 587
# smtp_mail_from = airflow@example.com

with DAG(
    dag_id       = "email_alerts_dag",
    schedule     = "@daily",
    start_date   = datetime(2024, 1, 1),
    catchup      = False,
    default_args = {
        "email":            ["team@example.com"],
        "email_on_failure": True,
        "email_on_retry":   False,
        "retries":          2,
    },
) as dag:
    pass
```

```python
# Custom email content
from airflow.utils.email import send_email

def send_custom_alert(context):
    task_id = context["task_instance"].task_id
    dag_id  = context["dag"].dag_id
    ds      = context["ds"]
    log_url = context["task_instance"].log_url

    send_email(
        to      = ["oncall@example.com"],
        subject = f"[AIRFLOW] {dag_id}.{task_id} failed on {ds}",
        html_content = f"""
            <h3>Task Failure Alert</h3>
            <p>DAG: <b>{dag_id}</b></p>
            <p>Task: <b>{task_id}</b></p>
            <p>Date: <b>{ds}</b></p>
            <p><a href="{log_url}">View Logs</a></p>
        """,
    )
```

---

## SLAs (Service Level Agreements)

SLAs trigger a callback if a task hasn't completed within the specified duration:

```python
from datetime import timedelta

def sla_miss_callback(dag, task_list, blocking_task_list, slas, blocking_tis):
    missed = [str(s) for s in slas]
    print(f"SLA missed for tasks: {missed}")
    # Send alert here

with DAG(
    dag_id           = "sla_dag",
    schedule         = "@daily",
    start_date       = datetime(2024, 1, 1),
    catchup          = False,
    sla_miss_callback = sla_miss_callback,
) as dag:

    # Task must complete within 1 hour of the scheduled run time
    extract = PythonOperator(
        task_id         = "extract",
        python_callable = lambda: None,
        sla             = timedelta(hours=1),
    )

    # Task must complete within 2 hours
    load = PythonOperator(
        task_id         = "load",
        python_callable = lambda: None,
        sla             = timedelta(hours=2),
    )
```

---

## Task Instance Timeouts

```python
from datetime import timedelta

task = PythonOperator(
    task_id             = "long_running_task",
    python_callable     = run_heavy_job,
    execution_timeout   = timedelta(hours=2),   # kill task after 2 hours
    retries             = 3,
    retry_delay         = timedelta(minutes=5),
    retry_exponential_backoff = True,           # delay doubles each retry
    max_retry_delay     = timedelta(minutes=60),
)
```

---

## Logs

### Accessing Logs

```bash
# CLI
airflow tasks logs my_dag my_task 2024-01-15

# List log files
airflow tasks logs --full-content my_dag my_task 2024-01-15
```

### Log Configuration

```properties
# airflow.cfg
[logging]
base_log_folder = /opt/airflow/logs
logging_level   = INFO

# Remote logging to S3
remote_logging        = True
remote_base_log_folder = s3://my-airflow-logs/logs
remote_log_conn_id    = aws_default

# Remote logging to GCS
remote_base_log_folder = gs://my-airflow-logs/logs
remote_log_conn_id    = google_cloud_default
```

### Structured Logging in Tasks

```python
import logging

log = logging.getLogger(__name__)

def my_task(**context):
    log.info("Starting task for date: %s", context["ds"])
    log.warning("Potential issue: %s", "low memory")
    log.error("Something went wrong: %s", "connection refused")

    # airflow.task logger is accessible via context
    context["task"].log.info("Also works via task logger")
```

---

## Metrics & Prometheus

Airflow emits StatsD metrics by default; use the Prometheus exporter for Grafana:

```properties
# airflow.cfg
[metrics]
statsd_on   = True
statsd_host = localhost
statsd_port = 8125
statsd_prefix = airflow
```

```yaml
# docker-compose — add Prometheus exporter
prometheus:
  image: prom/prometheus
  volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml

airflow-statsd-exporter:
  image: prom/statsd-exporter
  ports:
    - "9102:9102"
  command:
    - "--statsd.listen-udp=:8125"
    - "--web.listen-address=:9102"
```

### Key Metrics

| Metric | What It Measures |
|---|---|
| `airflow.dag.loading_duration` | Time to load DAG files |
| `airflow.scheduler.heartbeat` | Scheduler is alive |
| `airflow.task.duration` | Task execution time |
| `airflow.task.start_delta` | Delay from scheduled to actual start |
| `airflow.dagrun.duration.success` | DAG run duration (success) |
| `airflow.dagrun.duration.failed` | DAG run duration (failed) |
| `airflow.executor.open_slots` | Available worker slots |
| `airflow.executor.queued_tasks` | Tasks waiting to run |
| `airflow.executor.running_tasks` | Tasks currently running |

---

## Monitoring DAG Health via REST API

```python
import requests

AIRFLOW_URL = "http://localhost:8080"
AUTH        = ("admin", "admin")

def get_dag_runs(dag_id: str, state: str = None) -> list:
    params = {"state": state} if state else {}
    resp   = requests.get(
        f"{AIRFLOW_URL}/api/v1/dags/{dag_id}/dagRuns",
        auth   = AUTH,
        params = params,
    )
    return resp.json()["dag_runs"]

def get_failed_tasks(dag_id: str, run_id: str) -> list:
    resp = requests.get(
        f"{AIRFLOW_URL}/api/v1/dags/{dag_id}/dagRuns/{run_id}/taskInstances",
        auth   = AUTH,
        params = {"state": "failed"},
    )
    return resp.json()["task_instances"]

def clear_task(dag_id: str, run_id: str, task_id: str) -> None:
    requests.post(
        f"{AIRFLOW_URL}/api/v1/dags/{dag_id}/dagRuns/{run_id}/taskInstances/clear",
        auth = AUTH,
        json = {"task_ids": [task_id], "include_downstream": True},
    )
```

---

## Alerting Best Practices

```
1. Always set on_failure_callback for production DAGs
   → Prefer Slack/PagerDuty over email for on-call alerting

2. Set execution_timeout on every task
   → Prevents stuck tasks from blocking the worker pool

3. Use SLAs for data freshness SLOs
   → SLA miss ≠ failure — task may still succeed late

4. Set retry + retry_delay for transient failures
   → retry_exponential_backoff for external API calls

5. Emit custom metrics from tasks
   → airflow.stats.timing("my_dag.extract.rows", row_count)
   → Enables data quality dashboards

6. Log structured data (JSON lines) when possible
   → Enables log aggregation tools (Datadog, CloudWatch, ELK)
```


---

# Chapter 15: Security

## Authentication

### Default Setup (Username/Password)

```python
# webserver_config.py
from flask_appbuilder.security.manager import AUTH_DB

AUTH_TYPE = AUTH_DB

# Create admin user via CLI
# airflow users create --role Admin --username admin --firstname Admin \
#   --lastname User --email admin@example.com --password admin123
```

### OAuth2 / OIDC (Google, GitHub, Okta)

```python
# webserver_config.py
from flask_appbuilder.security.manager import AUTH_OAUTH

AUTH_TYPE       = AUTH_OAUTH
AUTH_USER_REGISTRATION      = True
AUTH_USER_REGISTRATION_ROLE = "Viewer"

OAUTH_PROVIDERS = [
    {
        "name":        "google",
        "token_key":   "access_token",
        "icon":        "fa-google",
        "remote_app":  {
            "client_id":     "your-client-id.apps.googleusercontent.com",
            "client_secret": "your-client-secret",
            "api_base_url":  "https://www.googleapis.com/oauth2/v2/",
            "client_kwargs": {"scope": "email profile"},
            "request_token_url": None,
            "access_token_url":  "https://accounts.google.com/o/oauth2/token",
            "authorize_url":     "https://accounts.google.com/o/oauth2/auth",
        },
    }
]
```

### LDAP Authentication

```python
# webserver_config.py
from flask_appbuilder.security.manager import AUTH_LDAP

AUTH_TYPE         = AUTH_LDAP
AUTH_LDAP_SERVER  = "ldap://ldap.example.com:389"
AUTH_LDAP_USE_TLS = True
AUTH_LDAP_SEARCH   = "dc=example,dc=com"
AUTH_LDAP_BIND_USER     = "cn=service,dc=example,dc=com"
AUTH_LDAP_BIND_PASSWORD = "bind-password"
AUTH_LDAP_UID_FIELD     = "uid"
```

---

## RBAC (Role-Based Access Control)

Airflow ships with five built-in roles:

| Role | Permissions |
|---|---|
| `Admin` | Full access including security management |
| `Op` | Create/edit DAGs, connections, variables |
| `User` | Trigger DAGs, view runs and logs |
| `Viewer` | Read-only access |
| `Public` | No access (anonymous users) |

### Managing Users and Roles

```bash
# Create a user
airflow users create \
  --role   User \
  --username john.doe \
  --firstname John \
  --lastname  Doe \
  --email john.doe@example.com \
  --password  SecurePass123

# List users
airflow users list

# Add role to user
airflow users add-role --username john.doe --role Op

# Remove role
airflow users remove-role --username john.doe --role Op

# List roles
airflow roles list

# Create custom role
airflow roles create data_engineer
```

### Custom Roles via Python

```python
# webserver_config.py
from airflow.www.security import AirflowSecurityManager

class MySecurityManager(AirflowSecurityManager):

    def get_user_roles(self, user=None):
        roles = super().get_user_roles(user)
        # Add custom logic: e.g., assign role based on email domain
        if user and user.email.endswith("@data-team.example.com"):
            roles.append(self.find_role("Op"))
        return roles

SECURITY_MANAGER_CLASS = MySecurityManager
```

### DAG-Level Access Control

```python
# In the DAG definition, restrict which roles can view/trigger
with DAG(
    dag_id       = "finance_pipeline",
    access_control = {
        "finance_team": {"can_read", "can_edit", "can_dag_run"},
        "Admin":        {"can_read", "can_edit", "can_dag_run", "can_delete"},
    },
    ...
) as dag:
    pass
```

---

## Fernet Encryption

Airflow encrypts connection passwords and variable values in the metadata DB using Fernet symmetric encryption.

```bash
# Generate a Fernet key
python -c "from cryptography.fernet import Fernet; print(Fernet.generate_key().decode())"
```

```properties
# airflow.cfg
[core]
fernet_key = your_fernet_key_here

# Or via environment variable
AIRFLOW__CORE__FERNET_KEY=your_fernet_key_here
```

**Important:** Keep the Fernet key secret and back it up. Losing it means you cannot decrypt stored passwords.

### Key Rotation

```bash
# After setting a new fernet_key alongside the old one:
airflow rotate-fernet-key
```

---

## Secrets Backends

Instead of storing secrets in the Airflow metadata DB, delegate to a dedicated secrets manager.

### AWS Secrets Manager

```properties
# airflow.cfg
[secrets]
backend = airflow.providers.amazon.aws.secrets.secrets_manager.SecretsManagerBackend
backend_kwargs = {
    "connections_prefix": "airflow/connections",
    "variables_prefix":   "airflow/variables",
    "config_prefix":      "airflow/config"
}
```

```bash
# Store a connection as a plain URI
aws secretsmanager create-secret \
  --name "airflow/connections/my_postgres" \
  --secret-string "postgresql://user:pass@host:5432/db"

# Store a variable
aws secretsmanager create-secret \
  --name "airflow/variables/api_key" \
  --secret-string "super-secret-key"
```

### HashiCorp Vault

```properties
[secrets]
backend = airflow.providers.hashicorp.secrets.vault.VaultBackend
backend_kwargs = {
    "connections_path": "connections",
    "variables_path":   "variables",
    "mount_point":      "airflow",
    "url":              "http://vault:8200",
    "auth_type":        "token",
    "token":            "vault-root-token"
}
```

```bash
# Store in Vault
vault kv put airflow/connections/my_postgres \
  conn_uri="postgresql://user:pass@host:5432/db"

vault kv put airflow/variables/batch_size value=500
```

### GCP Secret Manager

```properties
[secrets]
backend = airflow.providers.google.cloud.secrets.secret_manager.CloudSecretManagerBackend
backend_kwargs = {
    "connections_prefix": "airflow-connections",
    "variables_prefix":   "airflow-variables",
    "gcp_key_path":       "/opt/airflow/gcp_service_account.json"
}
```

---

## Network Security

### TLS for Webserver

```properties
# airflow.cfg
[webserver]
web_server_ssl_cert = /opt/airflow/certs/airflow.crt
web_server_ssl_key  = /opt/airflow/certs/airflow.key
```

### Webserver Behind a Reverse Proxy (Nginx)

```nginx
# /etc/nginx/sites-available/airflow
server {
    listen 443 ssl;
    server_name airflow.example.com;

    ssl_certificate     /etc/nginx/certs/airflow.crt;
    ssl_certificate_key /etc/nginx/certs/airflow.key;

    location / {
        proxy_pass         http://localhost:8080;
        proxy_set_header   Host $host;
        proxy_set_header   X-Real-IP $remote_addr;
        proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
    }
}
```

```properties
# airflow.cfg — tell Airflow it's behind a proxy
[webserver]
base_url              = https://airflow.example.com
enable_proxy_fix      = True
proxy_fix_x_for       = 1
proxy_fix_x_proto     = 1
proxy_fix_x_host      = 1
proxy_fix_x_prefix    = 1
```

---

## API Security

```properties
# airflow.cfg
[api]
# Require authentication for REST API (Airflow 2.0+)
auth_backends = airflow.api.auth.backend.basic_auth

# Or JWT-based
auth_backends = airflow.api.auth.backend.kerberos_auth
```

```bash
# API with basic auth
curl -X GET "http://localhost:8080/api/v1/dags" \
  -u admin:admin \
  -H "Content-Type: application/json"

# API with JWT token
TOKEN=$(curl -X POST "http://localhost:8080/auth/token" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' | jq -r '.access_token')

curl -X GET "http://localhost:8080/api/v1/dags" \
  -H "Authorization: Bearer $TOKEN"
```

---

## Security Checklist

```
Authentication:
  ☐ Disable default admin/admin credentials
  ☐ Enable OAuth or LDAP for enterprise deployments
  ☐ Enforce strong password policy

Encryption:
  ☐ Set a strong FERNET_KEY
  ☐ Rotate Fernet key periodically
  ☐ Enable TLS on webserver
  ☐ Enable TLS between Airflow components and DB

Secrets management:
  ☐ Use a secrets backend (Vault, AWS SM, GCP SM)
  ☐ Never hardcode credentials in DAG files
  ☐ Never store secrets in Airflow Variables (plaintext in UI)

Access control:
  ☐ Apply least-privilege roles to all users
  ☐ Use DAG-level access_control for sensitive pipelines
  ☐ Audit user access regularly

Network:
  ☐ Webserver behind reverse proxy with TLS
  ☐ Restrict metadata DB access to Airflow services only
  ☐ Enable API authentication
  ☐ Put Airflow in a private network, expose only webserver

Code security:
  ☐ Review DAGs before merging (no shell injection)
  ☐ Validate BashOperator commands
  ☐ Avoid dynamic SQL without parameterization
```


---

# Chapter 16: Performance & Scaling

## Understanding the Bottlenecks

```
Parsing bottleneck:
  Scheduler scans DAG files every min_file_process_interval seconds
  Slow imports or heavy module-level code slows ALL DAGs

Scheduling bottleneck:
  Scheduler runs the critical section loop every scheduler_heartbeat_sec
  Thousands of tasks in the ready queue slow scheduling

Executor bottleneck:
  Limited worker slots → tasks queue up
  CeleryExecutor: add workers
  KubernetesExecutor: add node pool capacity

Database bottleneck:
  Too many tasks writing metadata simultaneously
  Use connection pooling (PgBouncer for PostgreSQL)
```

---

## Scheduler Tuning

```properties
# airflow.cfg

[scheduler]
# Seconds between scheduling loops (default 5)
scheduler_heartbeat_sec = 5

# How often the scheduler scans for new/changed DAG files (default 30)
min_file_process_interval = 30

# Number of parallel threads for DAG parsing (default 2)
parsing_processes = 4

# Max DAG runs created per scheduler loop iteration (default 100)
max_dagruns_per_loop_to_schedule = 100

# Max TIs that can be scheduled per scheduler loop (default 512)
max_tis_per_query = 512

# Seconds to wait for task to heartbeat before marking zombie (default 300)
scheduler_zombie_task_threshold = 300

# Run health check every N seconds
scheduler_health_check_threshold = 30
```

---

## DAG Parsing Performance

```python
# BAD: heavy import at module level — runs on every scheduler heartbeat
import pandas as pd
import tensorflow as tf                # ← slow import
from my_utils import load_big_config   # ← reads a file

GLOBAL_CONFIG = load_big_config()      # ← executed at parse time

@dag(...)
def my_dag():
    @task
    def train():
        import tensorflow as tf        # ← import inside task instead
        ...
```

```python
# GOOD: lazy imports, no work at parse time
from airflow.decorators import dag, task
from datetime import datetime

@dag(start_date=datetime(2024, 1, 1), schedule="@daily", catchup=False)
def fast_dag():

    @task
    def extract() -> list:
        import pandas as pd        # lazy import
        return []

    @task
    def transform(data: list) -> list:
        from my_utils import process   # lazy import
        return process(data)

    transform(extract())

fast_dag()
```

### Measuring Parse Time

```bash
# Time how long it takes to parse all DAGs
airflow dags report

# Parse a specific file and print time
time python /opt/airflow/dags/my_dag.py

# Profile DAG parsing
python -c "
import cProfile
import airflow.models.dagbag
cProfile.run('airflow.models.dagbag.DagBag()')
"
```

---

## Concurrency Configuration

```properties
# airflow.cfg

[core]
# Total tasks running simultaneously across ALL DAGs (default 32)
parallelism = 64

# Max tasks per single DAG across all active runs (default 16)
max_active_tasks_per_dag = 32

# Max concurrent DAG runs per DAG (default 16)
max_active_runs_per_dag = 5

[celery]
# Tasks each Celery worker can run simultaneously (default 16)
worker_concurrency = 32
```

```python
# Per-DAG concurrency limits
with DAG(
    dag_id          = "my_dag",
    max_active_runs = 2,       # no more than 2 runs at the same time
    max_active_tasks = 16,     # no more than 16 tasks running across all runs
    concurrency     = 16,      # legacy alias for max_active_tasks
) as dag:
    pass
```

---

## Pools for Resource Management

```bash
# Create pools
airflow pools set default_pool      128  "Default pool"
airflow pools set db_pool             5  "Limit concurrent DB queries"
airflow pools set api_pool           10  "Limit concurrent API calls"
airflow pools set heavy_compute_pool  3  "GPU/heavy jobs"
```

```python
from airflow.operators.python import PythonOperator

# Light task — no pool needed (uses default_pool)
light_task = PythonOperator(
    task_id         = "light",
    python_callable = lambda: None,
)

# DB-heavy task — limit concurrent access
db_task = PythonOperator(
    task_id         = "query_warehouse",
    python_callable = run_query,
    pool            = "db_pool",
    pool_slots      = 1,            # uses 1 of 5 db_pool slots
)

# Very heavy task — takes 3 slots
heavy_task = PythonOperator(
    task_id         = "full_table_scan",
    python_callable = run_heavy,
    pool            = "db_pool",
    pool_slots      = 3,            # blocks 3 other tasks in db_pool
)
```

---

## Priority Weights

```python
# Higher priority_weight = scheduled first when slots are contested
from airflow.utils.weight_rule import WeightRule

urgent = PythonOperator(
    task_id         = "urgent_report",
    python_callable = run_report,
    priority_weight = 100,          # default is 1
    weight_rule     = WeightRule.ABSOLUTE,  # use exactly this weight
)

# weight_rule options:
# DOWNSTREAM (default): weight = own + sum of all downstream weights
# UPSTREAM:             weight = own + sum of all upstream weights
# ABSOLUTE:             weight = own weight only
```

---

## CeleryExecutor Scaling

```yaml
# docker-compose.yml — scale workers horizontally
worker:
  image: apache/airflow:2.9.0
  command: celery worker
  environment:
    AIRFLOW__CORE__EXECUTOR: CeleryExecutor
    AIRFLOW__CELERY__WORKER_CONCURRENCY: 16
  deploy:
    replicas: 4   # 4 workers × 16 concurrency = 64 total slots
```

```bash
# Scale workers at runtime (docker swarm / kubernetes)
docker service scale airflow_worker=8

# Graceful worker shutdown (finish current tasks, then stop)
airflow celery worker --autoscale=16,4  # max 16, min 4 concurrent tasks

# Worker queues for heterogeneous workloads
airflow celery worker -q default,heavy_compute
airflow celery worker -q gpu_tasks -c 2    # dedicated GPU worker, low concurrency
```

---

## KubernetesExecutor Scaling

```yaml
# The KubernetesExecutor scales to zero — no idle workers
# Configure the pod template for all tasks

# airflow.cfg
[kubernetes]
worker_pods_creation_batch_size = 16  # create up to 16 pods per scheduler loop
delete_worker_pods               = True
delete_worker_pods_on_failure    = False   # keep failed pods for debugging
worker_pods_pending_timeout      = 300     # fail task if pod stays pending > 5min
```

```python
# Per-task pod resource requests
from kubernetes.client import models as k8s

light_task = PythonOperator(
    task_id         = "light",
    python_callable = run_light,
    executor_config = {
        "pod_override": k8s.V1Pod(spec=k8s.V1PodSpec(containers=[
            k8s.V1Container(
                name="base",
                resources=k8s.V1ResourceRequirements(
                    requests={"cpu": "0.5", "memory": "512Mi"},
                    limits  ={"cpu": "1",   "memory": "1Gi"},
                ),
            )
        ]))
    },
)

heavy_task = PythonOperator(
    task_id         = "heavy",
    python_callable = run_heavy,
    executor_config = {
        "pod_override": k8s.V1Pod(spec=k8s.V1PodSpec(containers=[
            k8s.V1Container(
                name="base",
                resources=k8s.V1ResourceRequirements(
                    requests={"cpu": "4",   "memory": "16Gi"},
                    limits  ={"cpu": "8",   "memory": "32Gi"},
                ),
            )
        ], node_selector={"workload": "high-memory"}))
    },
)
```

---

## Database Connection Pooling

For large deployments, use PgBouncer between Airflow and PostgreSQL:

```yaml
# docker-compose.yml
pgbouncer:
  image: edoburu/pgbouncer
  environment:
    DATABASE_URL: postgresql://airflow:airflow@postgres:5432/airflow
    POOL_MODE:    transaction
    MAX_CLIENT_CONN: 200
    DEFAULT_POOL_SIZE: 20
  ports:
    - "5432:5432"
```

```properties
# airflow.cfg — point to PgBouncer instead of Postgres
[database]
sql_alchemy_conn = postgresql+psycopg2://airflow:airflow@pgbouncer:5432/airflow
sql_alchemy_pool_size = 5
sql_alchemy_max_overflow = 10
sql_alchemy_pool_recycle = 1800
```

---

## Performance Checklist

```
DAG parsing:
  ☐ No heavy imports at module level
  ☐ No network/DB calls at parse time
  ☐ No Variable.get() at module level
  ☐ DAG files < 1s parse time (measure with airflow dags report)

Scheduling:
  ☐ parallelism matches executor capacity
  ☐ parsing_processes = number of CPU cores / 2
  ☐ max_active_runs_per_dag prevents thundering herd

Worker capacity:
  ☐ worker_concurrency matched to workload (IO vs CPU)
  ☐ Pools configured for shared resources (DB, API)
  ☐ Priority weights set for critical pipelines

Database:
  ☐ PostgreSQL (not SQLite or MySQL for production)
  ☐ PgBouncer for high-concurrency deployments
  ☐ Regular VACUUM / ANALYZE on metadata DB
  ☐ Prune old dag runs: airflow db clean
```

---

## Pruning Metadata DB

```bash
# Clean old task instances and dag runs (keeps last 30 days)
airflow db clean --clean-before-timestamp "$(date -d '-30 days' --iso-8601=seconds)"

# Clean specific table
airflow db clean --tables task_instance --clean-before-timestamp 2024-01-01

# Dry run first
airflow db clean --dry-run --clean-before-timestamp 2024-01-01
```


---

# Chapter 17: Airflow with Cloud Providers

## Provider Packages

Cloud integrations are installed as separate provider packages:

```bash
pip install apache-airflow-providers-amazon          # AWS
pip install apache-airflow-providers-google          # GCP
pip install apache-airflow-providers-microsoft-azure # Azure
pip install apache-airflow-providers-databricks      # Databricks
pip install apache-airflow-providers-snowflake       # Snowflake
```

---

## AWS

### Connections

```bash
# Set up AWS connection
airflow connections add aws_default \
  --conn-type aws \
  --conn-extra '{"aws_access_key_id": "AKIA...", "aws_secret_access_key": "...", "region_name": "us-east-1"}'

# Or via environment variable
export AIRFLOW_CONN_AWS_DEFAULT='{"conn_type": "aws", "extra": {"region_name": "us-east-1"}}'
# → Uses IAM role or instance profile automatically when running on AWS
```

### S3 Operations

```python
from airflow.providers.amazon.aws.operators.s3 import S3CopyObjectOperator, S3DeleteObjectsOperator
from airflow.providers.amazon.aws.sensors.s3 import S3KeySensor
from airflow.providers.amazon.aws.hooks.s3 import S3Hook

# Sensor: wait for file
wait = S3KeySensor(
    task_id     = "wait_for_file",
    bucket_name = "my-data-bucket",
    bucket_key  = "raw/{{ ds }}/orders.parquet",
    aws_conn_id = "aws_default",
    mode        = "reschedule",
    timeout     = 3600,
    deferrable  = True,
)

# Copy object
copy = S3CopyObjectOperator(
    task_id             = "copy_to_processed",
    source_bucket_name  = "my-data-bucket",
    source_bucket_key   = "raw/{{ ds }}/orders.parquet",
    dest_bucket_name    = "my-data-bucket",
    dest_bucket_key     = "processed/{{ ds }}/orders.parquet",
    aws_conn_id         = "aws_default",
)

# Use S3Hook directly in a task
from airflow.decorators import task

@task
def upload_results(results: list, ds: str = None) -> str:
    import json, io
    hook = S3Hook(aws_conn_id="aws_default")
    key  = f"output/{ds}/results.json"
    hook.load_string(
        string_data = json.dumps(results),
        key         = key,
        bucket_name = "my-data-bucket",
        replace     = True,
    )
    return f"s3://my-data-bucket/{key}"
```

### Redshift

```python
from airflow.providers.amazon.aws.operators.redshift_sql import RedshiftSQLOperator
from airflow.providers.amazon.aws.transfers.s3_to_redshift import S3ToRedshiftOperator

# Run SQL on Redshift
query = RedshiftSQLOperator(
    task_id      = "transform_redshift",
    redshift_conn_id = "redshift_default",
    sql          = """
        INSERT INTO analytics.daily_orders
        SELECT date, COUNT(*) as order_count, SUM(amount) as total
        FROM staging.orders
        WHERE date = '{{ ds }}'
        GROUP BY date
    """,
)

# Load S3 → Redshift via COPY
load = S3ToRedshiftOperator(
    task_id             = "load_to_redshift",
    schema              = "staging",
    table               = "orders",
    s3_bucket           = "my-data-bucket",
    s3_key              = "raw/{{ ds }}/orders.parquet",
    redshift_conn_id    = "redshift_default",
    aws_conn_id         = "aws_default",
    copy_options        = ["FORMAT AS PARQUET"],
    method              = "REPLACE",
)
```

### EMR (Spark on AWS)

```python
from airflow.providers.amazon.aws.operators.emr import (
    EmrCreateJobFlowOperator,
    EmrAddStepsOperator,
    EmrTerminateJobFlowOperator,
)
from airflow.providers.amazon.aws.sensors.emr import EmrStepSensor

create_cluster = EmrCreateJobFlowOperator(
    task_id         = "create_emr_cluster",
    aws_conn_id     = "aws_default",
    job_flow_overrides = {
        "Name": "airflow-spark-cluster",
        "ReleaseLabel": "emr-7.0.0",
        "Instances": {
            "InstanceGroups": [
                {"Name": "Master", "InstanceRole": "MASTER", "InstanceType": "m5.xlarge", "InstanceCount": 1},
                {"Name": "Workers","InstanceRole": "CORE",   "InstanceType": "m5.xlarge", "InstanceCount": 3},
            ],
            "KeepJobFlowAliveWhenNoSteps": True,
        },
    },
)

add_step = EmrAddStepsOperator(
    task_id       = "run_spark_job",
    job_flow_id   = "{{ task_instance.xcom_pull('create_emr_cluster', key='return_value') }}",
    aws_conn_id   = "aws_default",
    steps         = [{
        "Name": "Transform",
        "ActionOnFailure": "TERMINATE_CLUSTER",
        "HadoopJarStep": {
            "Jar": "command-runner.jar",
            "Args": [
                "spark-submit", "--deploy-mode", "cluster",
                "s3://my-bucket/scripts/transform.py",
                "--date", "{{ ds }}",
            ],
        },
    }],
)

wait_for_step = EmrStepSensor(
    task_id     = "wait_for_spark",
    job_flow_id = "{{ task_instance.xcom_pull('create_emr_cluster', key='return_value') }}",
    step_id     = "{{ task_instance.xcom_pull('run_spark_job', key='return_value')[0] }}",
    aws_conn_id = "aws_default",
)

terminate = EmrTerminateJobFlowOperator(
    task_id    = "terminate_cluster",
    job_flow_id = "{{ task_instance.xcom_pull('create_emr_cluster', key='return_value') }}",
    aws_conn_id = "aws_default",
)

create_cluster >> add_step >> wait_for_step >> terminate
```

---

## GCP (Google Cloud Platform)

### Connection

```bash
airflow connections add google_cloud_default \
  --conn-type google_cloud_platform \
  --conn-extra '{"key_path": "/opt/airflow/gcp_key.json", "project": "my-gcp-project"}'
```

### BigQuery

```python
from airflow.providers.google.cloud.operators.bigquery import (
    BigQueryInsertJobOperator,
    BigQueryCreateEmptyTableOperator,
    BigQueryDeleteTableOperator,
    BigQueryCheckOperator,
)
from airflow.providers.google.cloud.sensors.bigquery import BigQueryTableExistenceSensor

# Run a BQ query
run_query = BigQueryInsertJobOperator(
    task_id    = "run_bq_query",
    gcp_conn_id = "google_cloud_default",
    configuration = {
        "query": {
            "query": """
                SELECT date, COUNT(*) as orders
                FROM `my-project.raw.orders`
                WHERE DATE(created_at) = '{{ ds }}'
                GROUP BY date
            """,
            "useLegacySql":     False,
            "destinationTable": {
                "projectId": "my-project",
                "datasetId": "analytics",
                "tableId":   "daily_orders${{ ds_nodash }}",
            },
            "writeDisposition": "WRITE_TRUNCATE",
        }
    },
)

# Check row count > 0
check = BigQueryCheckOperator(
    task_id    = "check_data",
    sql        = "SELECT COUNT(*) FROM `my-project.analytics.daily_orders` WHERE date = '{{ ds }}'",
    gcp_conn_id = "google_cloud_default",
    use_legacy_sql = False,
)
```

### GCS Operations

```python
from airflow.providers.google.cloud.operators.gcs import GCSDeleteObjectsOperator
from airflow.providers.google.cloud.transfers.gcs_to_bigquery import GCSToBigQueryOperator
from airflow.providers.google.cloud.sensors.gcs import GCSObjectExistenceSensor

# Sensor: wait for file in GCS
wait = GCSObjectExistenceSensor(
    task_id    = "wait_for_gcs",
    bucket     = "my-gcs-bucket",
    object     = "raw/{{ ds }}/orders.parquet",
    gcp_conn_id = "google_cloud_default",
    mode       = "reschedule",
)

# Load GCS → BigQuery
load = GCSToBigQueryOperator(
    task_id              = "gcs_to_bq",
    bucket               = "my-gcs-bucket",
    source_objects       = ["raw/{{ ds }}/orders.parquet"],
    destination_project_dataset_table = "my-project.staging.orders",
    source_format        = "PARQUET",
    write_disposition    = "WRITE_TRUNCATE",
    gcp_conn_id          = "google_cloud_default",
)
```

### Dataproc (Spark on GCP)

```python
from airflow.providers.google.cloud.operators.dataproc import (
    DataprocCreateClusterOperator,
    DataprocSubmitJobOperator,
    DataprocDeleteClusterOperator,
)

create = DataprocCreateClusterOperator(
    task_id      = "create_cluster",
    project_id   = "my-project",
    cluster_name = "airflow-spark-{{ ds_nodash }}",
    region       = "us-central1",
    cluster_config = {
        "master_config": {"num_instances": 1, "machine_type_uri": "n1-standard-4"},
        "worker_config": {"num_instances": 2, "machine_type_uri": "n1-standard-4"},
    },
    gcp_conn_id  = "google_cloud_default",
)

submit_job = DataprocSubmitJobOperator(
    task_id    = "run_spark",
    project_id = "my-project",
    region     = "us-central1",
    job        = {
        "placement":  {"cluster_name": "airflow-spark-{{ ds_nodash }}"},
        "pyspark_job": {
            "main_python_file_uri": "gs://my-bucket/scripts/transform.py",
            "args": ["--date", "{{ ds }}"],
        },
    },
    gcp_conn_id = "google_cloud_default",
)

delete = DataprocDeleteClusterOperator(
    task_id      = "delete_cluster",
    project_id   = "my-project",
    cluster_name = "airflow-spark-{{ ds_nodash }}",
    region       = "us-central1",
    gcp_conn_id  = "google_cloud_default",
    trigger_rule = "all_done",   # delete even if spark job failed
)

create >> submit_job >> delete
```

---

## Azure

### Azure Blob Storage

```python
from airflow.providers.microsoft.azure.operators.wasb_delete_blob import WasbDeleteBlobOperator
from airflow.providers.microsoft.azure.sensors.wasb import WasbBlobSensor
from airflow.providers.microsoft.azure.hooks.wasb import WasbHook

# Sensor: wait for blob
wait = WasbBlobSensor(
    task_id       = "wait_for_blob",
    container_name = "raw-data",
    blob_name     = "orders/{{ ds }}/data.parquet",
    wasb_conn_id  = "wasb_default",
    mode          = "reschedule",
)

# Use hook directly
@task
def upload_to_blob(data: bytes, ds: str = None) -> str:
    hook = WasbHook(wasb_conn_id="wasb_default")
    blob_name = f"output/{ds}/results.parquet"
    hook.load_bytes(data, container_name="output", blob_name=blob_name, overwrite=True)
    return blob_name
```

### Azure Synapse / SQL

```python
from airflow.providers.microsoft.azure.operators.synapse import AzureSynapseRunSparkBatchOperator

run_spark = AzureSynapseRunSparkBatchOperator(
    task_id     = "run_synapse_spark",
    azure_synapse_conn_id = "synapse_default",
    spark_pool  = "my-spark-pool",
    payload     = {
        "name":           "transform_orders",
        "file":           "abfss://scripts@mystorageaccount.dfs.core.windows.net/transform.py",
        "args":           ["--date", "{{ ds }}"],
        "driverMemory":   "4g",
        "executorMemory": "4g",
        "numExecutors":   4,
    },
)
```

---

## Databricks

```python
from airflow.providers.databricks.operators.databricks import DatabricksRunNowOperator, DatabricksSubmitRunOperator

# Trigger an existing Databricks job
run_job = DatabricksRunNowOperator(
    task_id          = "run_databricks_job",
    databricks_conn_id = "databricks_default",
    job_id           = 12345,
    notebook_params  = {"date": "{{ ds }}", "env": "production"},
)

# Submit a one-time run
submit = DatabricksSubmitRunOperator(
    task_id          = "submit_notebook",
    databricks_conn_id = "databricks_default",
    new_cluster      = {
        "spark_version": "13.3.x-scala2.12",
        "node_type_id":  "i3.xlarge",
        "num_workers":   4,
    },
    notebook_task    = {
        "notebook_path": "/Shared/ETL/transform_orders",
        "base_parameters": {"date": "{{ ds }}"},
    },
)
```

---

## Snowflake

```python
from airflow.providers.snowflake.operators.snowflake import SnowflakeOperator
from airflow.providers.snowflake.transfers.s3_to_snowflake import S3ToSnowflakeOperator

# Run SQL
query = SnowflakeOperator(
    task_id        = "transform_snowflake",
    snowflake_conn_id = "snowflake_default",
    sql            = """
        MERGE INTO ANALYTICS.ORDERS t
        USING STAGING.ORDERS s ON t.ORDER_ID = s.ORDER_ID
        WHEN MATCHED THEN UPDATE SET t.STATUS = s.STATUS
        WHEN NOT MATCHED THEN INSERT VALUES (s.ORDER_ID, s.STATUS, s.AMOUNT)
    """,
    warehouse      = "COMPUTE_WH",
    database       = "MY_DB",
    schema         = "ANALYTICS",
)

# Load from S3
load = S3ToSnowflakeOperator(
    task_id           = "s3_to_snowflake",
    s3_keys           = ["raw/{{ ds }}/orders.parquet"],
    table             = "ORDERS",
    schema            = "STAGING",
    stage             = "MY_S3_STAGE",
    file_format       = "PARQUET_FORMAT",
    snowflake_conn_id = "snowflake_default",
)
```


---

# Chapter 18: Custom Operators & Plugins

## Custom Operator

Extend `BaseOperator` to encapsulate reusable task logic:

```python
# plugins/operators/s3_to_postgres_operator.py
from airflow.models import BaseOperator
from airflow.providers.amazon.aws.hooks.s3 import S3Hook
from airflow.providers.postgres.hooks.postgres import PostgresHook
import pandas as pd
import io


class S3ToPostgresOperator(BaseOperator):
    """
    Load a Parquet file from S3 into a PostgreSQL table.
    Supports REPLACE, APPEND, and UPSERT modes.
    """

    template_fields = ("s3_key", "table")   # Jinja-templated fields

    def __init__(
        self,
        s3_conn_id:       str,
        s3_bucket:        str,
        s3_key:           str,
        postgres_conn_id: str,
        table:            str,
        schema:           str = "public",
        mode:             str = "replace",   # replace | append | upsert
        upsert_keys:      list[str] = None,
        **kwargs,
    ):
        super().__init__(**kwargs)
        self.s3_conn_id       = s3_conn_id
        self.s3_bucket        = s3_bucket
        self.s3_key           = s3_key
        self.postgres_conn_id = postgres_conn_id
        self.table            = table
        self.schema           = schema
        self.mode             = mode
        self.upsert_keys      = upsert_keys or []

    def execute(self, context):
        self.log.info("Reading s3://%s/%s", self.s3_bucket, self.s3_key)

        s3_hook = S3Hook(aws_conn_id=self.s3_conn_id)
        obj     = s3_hook.get_key(self.s3_key, bucket_name=self.s3_bucket)
        df      = pd.read_parquet(io.BytesIO(obj.get()["Body"].read()))

        self.log.info("Loaded %d rows from S3", len(df))

        pg_hook = PostgresHook(postgres_conn_id=self.postgres_conn_id)
        engine  = pg_hook.get_sqlalchemy_engine()

        if self.mode == "replace":
            df.to_sql(self.table, engine, schema=self.schema,
                      if_exists="replace", index=False)
        elif self.mode == "append":
            df.to_sql(self.table, engine, schema=self.schema,
                      if_exists="append", index=False)
        elif self.mode == "upsert":
            self._upsert(df, pg_hook)
        else:
            raise ValueError(f"Unknown mode: {self.mode}")

        self.log.info("Loaded %d rows into %s.%s", len(df), self.schema, self.table)
        return len(df)

    def _upsert(self, df: pd.DataFrame, hook: PostgresHook) -> None:
        staging = f"_staging_{self.table}"
        engine  = hook.get_sqlalchemy_engine()
        df.to_sql(staging, engine, schema=self.schema, if_exists="replace", index=False)

        key_clause    = " AND ".join(f"t.{k} = s.{k}" for k in self.upsert_keys)
        update_clause = ", ".join(
            f"{c} = s.{c}" for c in df.columns if c not in self.upsert_keys
        )
        insert_cols = ", ".join(df.columns)
        insert_vals = ", ".join(f"s.{c}" for c in df.columns)

        sql = f"""
            INSERT INTO {self.schema}.{self.table} ({insert_cols})
            SELECT {insert_vals} FROM {self.schema}.{staging} s
            ON CONFLICT ({", ".join(self.upsert_keys)})
            DO UPDATE SET {update_clause}
        """
        hook.run(sql)
        hook.run(f"DROP TABLE IF EXISTS {self.schema}.{staging}")
```

### Using the Custom Operator

```python
from plugins.operators.s3_to_postgres_operator import S3ToPostgresOperator

load_orders = S3ToPostgresOperator(
    task_id          = "load_orders",
    s3_conn_id       = "aws_default",
    s3_bucket        = "my-data-bucket",
    s3_key           = "processed/{{ ds }}/orders.parquet",
    postgres_conn_id = "postgres_default",
    table            = "orders",
    schema           = "staging",
    mode             = "upsert",
    upsert_keys      = ["order_id"],
)
```

---

## Custom Sensor

```python
# plugins/sensors/api_data_sensor.py
from airflow.sensors.base import BaseSensorOperator
import requests


class APIDataReadySensor(BaseSensorOperator):
    """Polls an internal API until data is available for the given date."""

    template_fields = ("endpoint", "date")

    def __init__(self, base_url: str, endpoint: str, date: str, min_records: int = 1, **kwargs):
        super().__init__(**kwargs)
        self.base_url    = base_url
        self.endpoint    = endpoint
        self.date        = date
        self.min_records = min_records

    def poke(self, context) -> bool:
        url = f"{self.base_url}/availability/{self.endpoint}"
        try:
            resp = requests.get(url, params={"date": self.date}, timeout=10)
            resp.raise_for_status()
            count = resp.json().get("record_count", 0)
            self.log.info("Records available: %d / %d needed", count, self.min_records)
            return count >= self.min_records
        except requests.RequestException as e:
            self.log.warning("API check failed: %s", e)
            return False
```

---

## Custom Hook

```python
# plugins/hooks/internal_api_hook.py
from airflow.hooks.base import BaseHook
import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry


class InternalAPIHook(BaseHook):
    """Hook for the internal data API with built-in retries."""

    conn_name_attr    = "api_conn_id"
    default_conn_name = "internal_api_default"

    def __init__(self, api_conn_id: str = "internal_api_default"):
        super().__init__()
        self.api_conn_id = api_conn_id
        self._session    = None

    def get_conn(self) -> requests.Session:
        if self._session is None:
            conn       = self.get_connection(self.api_conn_id)
            self._base_url = f"https://{conn.host}:{conn.port or 443}"

            session = requests.Session()
            retry   = Retry(total=3, backoff_factor=1, status_forcelist=[500, 502, 503])
            adapter = HTTPAdapter(max_retries=retry)
            session.mount("https://", adapter)
            session.headers.update({
                "Authorization": f"Bearer {conn.password}",
                "Content-Type":  "application/json",
            })
            self._session = session
        return self._session

    def get_data(self, path: str, params: dict = None) -> dict:
        resp = self.get_conn().get(f"{self._base_url}{path}", params=params)
        resp.raise_for_status()
        return resp.json()

    def post_data(self, path: str, payload: dict) -> dict:
        resp = self.get_conn().post(f"{self._base_url}{path}", json=payload)
        resp.raise_for_status()
        return resp.json()
```

---

## Plugins

Plugins let you register operators, hooks, sensors, macros, blueprints, and UI menu items:

```python
# plugins/my_plugin.py
from airflow.plugins_manager import AirflowPlugin
from plugins.operators.s3_to_postgres_operator import S3ToPostgresOperator
from plugins.sensors.api_data_sensor import APIDataReadySensor
from plugins.hooks.internal_api_hook import InternalAPIHook


def greeting_macro(name: str) -> str:
    return f"Hello, {name}!"


class MyCompanyPlugin(AirflowPlugin):
    name     = "my_company_plugin"

    operators = [S3ToPostgresOperator]
    sensors   = [APIDataReadySensor]
    hooks     = [InternalAPIHook]

    macros    = [greeting_macro]   # available as {{ greeting_macro("world") }}
```

Place the file in the `plugins/` directory and Airflow auto-discovers it.

---

## Decorating an Existing Operator

```python
# Add pre/post hooks to any operator by subclassing
from airflow.providers.postgres.operators.postgres import PostgresOperator
from datetime import datetime


class AuditedPostgresOperator(PostgresOperator):
    """PostgresOperator that logs execution to an audit table."""

    def execute(self, context):
        start = datetime.utcnow()
        try:
            result = super().execute(context)
            self._write_audit(context, "success", start)
            return result
        except Exception as e:
            self._write_audit(context, "failed", start, str(e))
            raise

    def _write_audit(self, context, status: str, start: datetime, error: str = None):
        from airflow.providers.postgres.hooks.postgres import PostgresHook
        hook = PostgresHook(postgres_conn_id=self.postgres_conn_id)
        hook.run(
            """
            INSERT INTO audit.task_runs (dag_id, task_id, ds, status, started_at, error)
            VALUES (%s, %s, %s, %s, %s, %s)
            """,
            parameters=(
                context["dag"].dag_id,
                context["task"].task_id,
                context["ds"],
                status,
                start,
                error,
            ),
        )
```

---

## Custom XCom Backend as a Plugin

```python
# plugins/xcoms/gcs_xcom_backend.py
from airflow.models.xcom import BaseXCom
import json, uuid
from google.cloud import storage


class GCSXComBackend(BaseXCom):
    BUCKET = "my-airflow-xcoms"

    @staticmethod
    def serialize_value(value, *, key=None, task_id=None, dag_id=None,
                        run_id=None, map_index=None):
        import pandas as pd
        if isinstance(value, pd.DataFrame):
            gcs_key = f"xcoms/{dag_id}/{run_id}/{task_id}/{uuid.uuid4()}.parquet"
            client  = storage.Client()
            bucket  = client.bucket(GCSXComBackend.BUCKET)
            bucket.blob(gcs_key).upload_from_string(value.to_parquet())
            return BaseXCom.serialize_value({"__type": "gcs_parquet", "key": gcs_key})
        return BaseXCom.serialize_value(value)

    @staticmethod
    def deserialize_value(result):
        import pandas as pd
        data = BaseXCom.deserialize_value(result)
        if isinstance(data, dict) and data.get("__type") == "gcs_parquet":
            client = storage.Client()
            bucket = client.bucket(GCSXComBackend.BUCKET)
            return pd.read_parquet(bucket.blob(data["key"]).download_as_bytes())
        return data
```

```properties
# airflow.cfg
[core]
xcom_backend = plugins.xcoms.gcs_xcom_backend.GCSXComBackend
```

---

## Plugin Directory Layout

```
project/
├── dags/
│   └── my_dag.py
├── plugins/
│   ├── __init__.py
│   ├── my_plugin.py              ← AirflowPlugin registration
│   ├── operators/
│   │   ├── __init__.py
│   │   └── s3_to_postgres_operator.py
│   ├── sensors/
│   │   ├── __init__.py
│   │   └── api_data_sensor.py
│   ├── hooks/
│   │   ├── __init__.py
│   │   └── internal_api_hook.py
│   └── xcoms/
│       ├── __init__.py
│       └── gcs_xcom_backend.py
└── tests/
    └── test_operators.py
```

---

## Testing Custom Operators

```python
# tests/test_s3_to_postgres_operator.py
from unittest.mock import MagicMock, patch
import pandas as pd
from plugins.operators.s3_to_postgres_operator import S3ToPostgresOperator

@patch("plugins.operators.s3_to_postgres_operator.S3Hook")
@patch("plugins.operators.s3_to_postgres_operator.PostgresHook")
def test_execute_append_mode(mock_pg_hook_cls, mock_s3_hook_cls):
    # Arrange
    sample_df   = pd.DataFrame({"id": [1, 2], "value": ["a", "b"]})
    mock_s3     = MagicMock()
    mock_s3.get_key.return_value.get.return_value = {"Body": sample_df.to_parquet()}
    mock_s3_hook_cls.return_value = mock_s3

    mock_pg     = MagicMock()
    mock_engine = MagicMock()
    mock_pg.get_sqlalchemy_engine.return_value = mock_engine
    mock_pg_hook_cls.return_value = mock_pg

    op = S3ToPostgresOperator(
        task_id          = "test",
        s3_conn_id       = "aws",
        s3_bucket        = "bucket",
        s3_key           = "key.parquet",
        postgres_conn_id = "pg",
        table            = "orders",
        mode             = "append",
    )

    # Act
    result = op.execute({})

    # Assert
    assert result == 2
    mock_s3.get_key.assert_called_once_with("key.parquet", bucket_name="bucket")
```


---

# Chapter 19: Airflow in Production

## Managed Airflow Options

| Service | Provider | Executor | Best For |
|---|---|---|---|
| MWAA | AWS | CeleryExecutor | AWS-native teams |
| Cloud Composer | GCP | KubernetesExecutor | GCP-native, K8s teams |
| Astronomer | Multi-cloud | CeleryExecutor or K8s | Enterprise, support |
| Astro (SaaS) | Astronomer Cloud | KubernetesExecutor | Fully managed |
| Self-hosted K8s | Any | KubernetesExecutor | Full control |

---

## AWS MWAA (Managed Workflows for Apache Airflow)

```python
# Create MWAA environment via CloudFormation / Terraform
# Resource: aws_mwaa_environment

resource "aws_mwaa_environment" "airflow" {
  name              = "production-airflow"
  airflow_version   = "2.9.1"
  environment_class = "mw1.medium"   # mw1.small | mw1.medium | mw1.large

  source_bucket_arn = aws_s3_bucket.airflow.arn
  dag_s3_path       = "dags/"
  requirements_s3_path = "requirements.txt"
  plugins_s3_path   = "plugins.zip"

  airflow_configuration_options = {
    "core.default_timezone"          = "UTC"
    "webserver.dag_default_view"     = "grid"
    "scheduler.min_file_process_interval" = "60"
  }

  network_configuration {
    security_group_ids = [aws_security_group.airflow.id]
    subnet_ids         = aws_subnet.private[*].id
  }

  logging_configuration {
    dag_processing_logs { enabled = true; log_level = "WARNING" }
    scheduler_logs      { enabled = true; log_level = "WARNING" }
    task_logs           { enabled = true; log_level = "INFO" }
    webserver_logs      { enabled = true; log_level = "WARNING" }
    worker_logs         { enabled = true; log_level = "WARNING" }
  }
}
```

```bash
# Deploy DAGs to MWAA (sync to S3)
aws s3 sync ./dags s3://my-mwaa-bucket/dags/ --delete
aws s3 cp requirements.txt s3://my-mwaa-bucket/requirements.txt
aws s3 cp plugins.zip s3://my-mwaa-bucket/plugins.zip
```

---

## GCP Cloud Composer

```bash
# Create a Cloud Composer environment
gcloud composer environments create production-airflow \
  --location us-central1 \
  --image-version composer-3-airflow-2.9.1 \
  --machine-type n1-standard-2 \
  --node-count 3 \
  --disk-size 50 \
  --python-version 3

# Deploy DAGs
gcloud composer environments storage dags import \
  --environment production-airflow \
  --location us-central1 \
  --source ./dags/

# Install Python packages
gcloud composer environments update production-airflow \
  --location us-central1 \
  --update-pypi-packages-from-file requirements.txt

# Set Airflow variables
gcloud composer environments run production-airflow \
  --location us-central1 variables -- \
  --set environment production
```

---

## Self-Hosted with Helm (Kubernetes)

```bash
# Add the official Helm chart repo
helm repo add apache-airflow https://airflow.apache.org
helm repo update

# Create values.yaml and install
helm upgrade --install airflow apache-airflow/airflow \
  --namespace airflow \
  --create-namespace \
  -f values.yaml
```

```yaml
# values.yaml (key sections)
airflowVersion: "2.9.1"
executor: KubernetesExecutor

images:
  airflow:
    repository: my-registry/airflow
    tag: "2.9.1-custom"

config:
  core:
    parallelism: "64"
    max_active_tasks_per_dag: "32"
  webserver:
    base_url: "https://airflow.mycompany.com"

dags:
  gitSync:
    enabled: true
    repo: https://github.com/mycompany/airflow-dags
    branch: main
    subPath: dags/
    credentialsSecret: airflow-git-credentials

workers:
  replicas: 0   # KubernetesExecutor creates pods on demand

webserver:
  replicas: 2
  service:
    type: ClusterIP

scheduler:
  replicas: 1

postgresql:
  enabled: false  # use external PostgreSQL

externalDatabase:
  type: postgres
  host: my-postgres.example.com
  port: 5432
  database: airflow
  user: airflow
  passwordSecret: airflow-db-credentials
  passwordSecretKey: password
```

---

## Production Architecture

```
                           ┌─────────────────────────────────┐
                           │           Kubernetes              │
                           │                                   │
  Users ──► ALB ──► ┌──────────┐   ┌──────────┐              │
                    │Webserver │   │Scheduler │              │
                    └──────────┘   └─────┬────┘              │
                         │               │ creates pods        │
                         │         ┌─────▼──────────────┐     │
                         │         │  Worker Pods (K8s)  │     │
                         │         └────────────────────-┘     │
                         │                                     │
                    ┌────▼─────────────────────────────┐      │
                    │         PostgreSQL (RDS)          │      │
                    └──────────────────────────────────┘      │
                                                               │
  Secrets ◄──── AWS Secrets Manager / Vault                   │
  DAGs    ◄──── Git Sync (github.com/org/dags)                │
  Logs    ────► S3 / GCS                                       │
  Metrics ────► Prometheus / Grafana                          │
                           └─────────────────────────────────┘
```

---

## DAG Deployment Strategy

### Git Sync (Recommended)

```yaml
# Kubernetes sidecar that continuously syncs DAGs from Git
dags:
  gitSync:
    enabled:           true
    repo:              https://github.com/mycompany/airflow-dags
    branch:            main
    rev:               HEAD
    depth:             1
    maxFailures:       3
    subPath:           dags/
    wait:              60   # sync every 60 seconds
    containerName:     git-sync
    uid:               65533
    credentialsSecret: airflow-git-credentials   # K8s secret with SSH key or token
```

### CI/CD Pipeline

```yaml
# .github/workflows/deploy-dags.yml
name: Deploy DAGs

on:
  push:
    branches: [main]
    paths: [dags/**, plugins/**, requirements.txt]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: pip install apache-airflow==2.9.0 pytest
      - run: airflow db init && pytest tests/

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      # MWAA
      - name: Sync to S3
        run: aws s3 sync ./dags s3://my-airflow-bucket/dags/ --delete

      # OR: Cloud Composer
      - name: Deploy to Composer
        run: |
          gcloud composer environments storage dags import \
            --environment production --location us-central1 \
            --source ./dags/
```

---

## High Availability

```yaml
# Multiple schedulers (Airflow 2.0+ supports HA scheduler)
scheduler:
  replicas: 2   # Both schedulers use optimistic row locking — safe

# Multiple webservers behind a load balancer
webserver:
  replicas: 2

# Shared AIRFLOW__WEBSERVER__SECRET_KEY — must match across replicas
# Shared Fernet key — must match across all components
# Shared metadata DB — PostgreSQL with connection pooling (PgBouncer)
# Shared DAG storage — Git sync or shared volume
```

---

## Production Runbook

### Health Checks

```bash
# Check scheduler is alive
airflow jobs check --job-type SchedulerJob --limit 1

# Check DAG parsing errors
airflow dags report | grep -i error

# Check stuck tasks (running > 2 hours)
airflow tasks states-for-dag-run my_dag 2024-01-15T00:00:00

# Check executor queue depth
airflow celery flower   # Celery
kubectl get pods -n airflow | grep "airflow-worker"   # K8s
```

### Incident Responses

```bash
# Clear and re-run a failed task
airflow tasks clear my_dag --task-ids extract --start-date 2024-01-15 --yes

# Mark a task as success (skip it)
airflow tasks state my_dag extract 2024-01-15
airflow tasks mark-success my_dag extract 2024-01-15

# Unpause a DAG
airflow dags unpause my_dag

# Kill a stuck running task instance
airflow tasks failed-deps my_dag extract 2024-01-15

# Force-fail a running DAG run
airflow dags state my_dag 2024-01-15T00:00:00

# Purge old metadata (keep 90 days)
airflow db clean --clean-before-timestamp "$(date -d '-90 days' --iso-8601=seconds)"
```

---

## Production Checklist

```
Infrastructure:
  ☐ PostgreSQL (not SQLite/MySQL)
  ☐ PgBouncer for connection pooling
  ☐ Remote log storage (S3, GCS)
  ☐ External secrets backend (Vault, AWS SM)
  ☐ TLS on webserver
  ☐ Webserver behind reverse proxy

DAGs:
  ☐ catchup=False on all production DAGs
  ☐ execution_timeout on every task
  ☐ retries + retry_delay configured
  ☐ on_failure_callback with alerting
  ☐ No module-level Variable.get() calls
  ☐ SLAs set for critical pipelines

Security:
  ☐ Fernet key set and backed up
  ☐ Default admin credentials changed
  ☐ RBAC configured
  ☐ API authentication enabled
  ☐ DAG-level access_control for sensitive DAGs

Observability:
  ☐ Prometheus metrics scraping
  ☐ Grafana dashboard
  ☐ Alerting on scheduler_heartbeat absence
  ☐ Alerting on task failure rate
  ☐ DAG SLA miss notifications

Operations:
  ☐ DB cleanup cron job (airflow db clean)
  ☐ Backup PostgreSQL metadata DB
  ☐ DAG deployment via CI/CD (not manual copy)
  ☐ Staging environment mirrors production
```


---

# Chapter 20: Learning Roadmap

## 5-Week Study Plan

### Week 1: Foundations

| Day | Topic | Practice |
|---|---|---|
| Mon | Ch 01–02: What is Airflow, Setup | Run `docker-compose up`, open UI |
| Tue | Ch 03: Core concepts | Draw the DAG → TaskInstance lifecycle |
| Wed | Ch 04: Writing DAGs | Write a 3-task pipeline (extract → transform → load) |
| Thu | Ch 05: Operators | Replace PythonOperator with BashOperator + SparkSubmitOperator |
| Fri | Ch 06–07: Sensors + Hooks | Add an S3KeySensor; use PostgresHook directly |
| Sat | Review | Refactor your pipeline to use TaskFlow API |
| Sun | Rest / catch-up | — |

**Week 1 Checkpoint:** You have a working DAG with sensors, operators, and a hook, deployed via docker-compose.

---

### Week 2: Data Flow & Scheduling

| Day | Topic | Practice |
|---|---|---|
| Mon | Ch 08: XCom & TaskFlow | Pass data between 3 tasks using @task decorators |
| Tue | Ch 09: Scheduling | Change schedule: cron → timedelta → Dataset-driven |
| Wed | Ch 09: Backfill | Run a backfill for 7 days; observe DAG runs in UI |
| Thu | Ch 10: Executors | Switch LocalExecutor → CeleryExecutor with Redis |
| Fri | Ch 11: Dynamic DAGs | Convert static 3-task DAG to dynamic mapping on a list |
| Sat | Ch 13: Variables + Params | Add env variable; trigger DAG with --conf params |
| Sun | Review | Write 5 unit tests for your task callables |

**Week 2 Checkpoint:** You understand scheduling semantics (data_interval, catchup) and have a dynamically-mapped pipeline.

---

### Week 3: Reliability & Testing

| Day | Topic | Practice |
|---|---|---|
| Mon | Ch 12: Testing DAGs | Add DAG validation + task logic tests; run in CI |
| Tue | Ch 14: Monitoring | Add Slack on_failure_callback to your DAG |
| Wed | Ch 14: SLAs | Add an SLA (1 hour) and trigger a miss |
| Thu | Ch 15: Security | Set Fernet key; use env var instead of UI connection |
| Fri | Ch 16: Performance | Measure DAG parse time; profile a slow import |
| Sat | Review | Create a pool; add priority_weight to critical tasks |
| Sun | Rest / catch-up | — |

**Week 3 Checkpoint:** Your pipeline has tests, alerting, and doesn't have import errors or slow parse time.

---

### Week 4: Cloud & Custom Extensions

| Day | Topic | Practice |
|---|---|---|
| Mon | Ch 17: AWS operators | Replace PythonOperator with S3KeySensor + S3Hook |
| Tue | Ch 17: BigQuery/GCP | Write a BigQueryInsertJobOperator task |
| Wed | Ch 18: Custom operators | Build S3ToPostgresOperator; use it in a DAG |
| Thu | Ch 18: Plugins | Wrap operator in an AirflowPlugin |
| Fri | Ch 10: KubernetesExecutor | Deploy to local K8s (minikube); configure pod_override |
| Sat | Ch 19: Production setup | Deploy with Helm on minikube; configure GitSync |
| Sun | Review | Add monitoring: Prometheus + Grafana dashboard |

**Week 4 Checkpoint:** You can write custom operators, deploy to Kubernetes, and have cloud provider integration.

---

### Week 5: Capstone Projects

Pick one capstone to build end-to-end:

---

#### Capstone A: E-Commerce Data Platform

```
Pipeline architecture:

[S3: raw orders]
      │
      ▼ S3KeySensor (deferrable)
[extract_orders] → PythonOperator: validate + dedupe
      │
      ▼
[transform] → PythonOperator: compute revenue, AOV
      │
      ├──► [load_to_postgres] → PostgresHook
      │
      └──► [load_to_redshift] → RedshiftSQLOperator
                │
                ▼
      [notify_slack] → SlackWebhookHook (on_success_callback)

Features to implement:
  ☐ Dynamic task mapping per product category
  ☐ ExternalTaskSensor waiting for inventory DAG
  ☐ Pool limiting Redshift concurrency
  ☐ SLA of 2 hours from midnight UTC
  ☐ Retry with exponential backoff
  ☐ Backfill last 30 days on first deploy
```

---

#### Capstone B: ML Training Pipeline

```
Pipeline architecture:

[Dataset: feature_store updated]
      │  (Dataset-driven schedule)
      ▼
[validate_features] → SqlSensor: count > 1000
      │
      ▼
[prepare_training_data] → dynamic mapping per model_type
      │                    (model_type: ["lr", "xgb", "nn"])
      ▼
[train_model[0/1/2]] → KubernetesPodOperator per model
      │                  (GPU node pool, 8 CPU, 32 GB RAM each)
      ▼
[evaluate_models] → compare metrics, pick best
      │
      ▼
[register_best_model] → push to MLflow / SageMaker
      │
      ▼
[notify_team] → Slack + email

Features to implement:
  ☐ Dataset-driven trigger from feature engineering DAG
  ☐ Dynamic mapping for parallel model training
  ☐ KubernetesPodOperator with GPU node selector
  ☐ Custom S3XComBackend for model artifacts
  ☐ Branching: only deploy if accuracy > 0.90
```

---

#### Capstone C: Real-Time Ingestion Monitor

```
Pipeline architecture (hourly schedule):

[extract_kafka_lag] → PythonOperator (confluent-kafka AdminClient)
      │
      ▼
[check_lag_thresholds] → BranchPythonOperator
      │
      ├── [lag_ok]      → EmptyOperator
      │
      └── [lag_critical] → TriggerDagRunOperator → scale_consumers_dag
                │
                ▼
      [alert_oncall] → PagerDuty / Slack

Features to implement:
  ☐ Hourly schedule with catchup=False
  ☐ Branching on consumer lag threshold
  ☐ TriggerDagRunOperator with conf payload
  ☐ Variable-driven thresholds (lag_warning_threshold)
  ☐ Audit trail: write task outcomes to PostgreSQL
```

---

## Commands Reference Card

```bash
# Start Airflow
docker-compose up -d

# DAG management
airflow dags list
airflow dags trigger my_dag --conf '{"key": "value"}'
airflow dags pause / unpause my_dag
airflow dags delete my_dag --yes

# Task management
airflow tasks list my_dag
airflow tasks test my_dag task_id 2024-01-15         # test without DB record
airflow tasks run  my_dag task_id 2024-01-15         # run and record in DB
airflow tasks clear my_dag --task-ids task_id --yes

# Scheduler
airflow scheduler
airflow jobs check --job-type SchedulerJob

# Celery
airflow celery worker -q default,heavy
airflow celery flower                                 # web UI for workers

# Database
airflow db init
airflow db upgrade
airflow db clean --clean-before-timestamp 2024-01-01 --yes

# Variables & Connections
airflow variables set key value
airflow variables get key
airflow connections add --conn-id id --conn-type type ...
airflow connections list

# Users
airflow users create --role Admin --username admin --email a@b.com --password pass
airflow users list
```

---

## Beyond Airflow

| What you'll learn next | Technology |
|---|---|
| Streaming pipelines | Apache Kafka + Flink |
| Lakehouse patterns | Apache Iceberg, Delta Lake |
| Feature engineering | Feast, Tecton |
| ML orchestration | MLflow, Kubeflow Pipelines |
| Modern orchestration | Prefect, Dagster |
| Data contracts | dbt, Great Expectations |
| Infrastructure as code | Terraform + Helm |

---

## Key Takeaways

```
1. DAGs are code — treat them like production software
   (tests, CI/CD, code review, linting)

2. Use TaskFlow API (@dag, @task) for new pipelines
   — cleaner, auto-XCom, less boilerplate

3. Executors define WHERE tasks run:
   Local = simplicity, Celery = scale out, K8s = isolation

4. Scheduling semantics matter:
   data_interval_start is what you process, not when you run

5. Sensors should use mode="reschedule" or deferrable=True
   for long waits

6. Never store large data in XCom — store in S3/GCS, pass the path

7. Use Secrets Backends for credentials — never hardcode or
   use plain Airflow Variables for sensitive data

8. Measure DAG parse time — module-level code runs constantly

9. The scheduler is a single point of failure in most setups —
   monitor its heartbeat aggressively

10. Production Airflow = PostgreSQL + TLS + Fernet + CI/CD + Alerting
```
