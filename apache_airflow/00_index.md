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
