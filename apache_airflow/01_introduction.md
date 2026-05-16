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
