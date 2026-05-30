# Apache Airflow + Laravel + Kafka Integration — Lessons

> Practical lessons for orchestrating Laravel/PHP workflows with Apache Airflow and Apache Kafka.

---

## Table of Contents

| Lesson | Title | Topics |
|---|---|---|
| 01 | [Introduction](01_introduction.md) | What Airflow should do, what Laravel should do, where Kafka fits |
| 02 | [Architecture](02_architecture.md) | Event-driven batch orchestration, services, topics, DAGs |
| 03 | [Local Setup](03_local_setup.md) | Docker layout, Laravel app, Kafka, Airflow, networking |
| 04 | [Laravel Kafka Producer](04_laravel_kafka_producer.md) | Publishing events from Laravel |
| 05 | [Laravel Kafka Consumer](05_laravel_kafka_consumer.md) | Consuming Kafka events in Laravel |
| 06 | [Triggering Airflow from Laravel](06_trigger_airflow_from_laravel.md) | REST API, auth, `dag_run.conf`, PHP HTTP client |
| 07 | [Airflow Calling Laravel](07_airflow_calling_laravel.md) | `HttpOperator`, internal APIs, Artisan commands |
| 08 | [Airflow and Kafka](08_airflow_and_kafka.md) | Producing/consuming Kafka from DAG tasks |
| 09 | [Data Pipeline Pattern](09_data_pipeline_pattern.md) | Export from Laravel, process in Airflow, publish result |
| 10 | [Reliability](10_reliability.md) | Idempotency, retries, DLQ, backfills, exactly-once thinking |
| 11 | [Security](11_security.md) | API tokens, network boundaries, SASL/TLS, secrets |
| 12 | [Monitoring](12_monitoring.md) | Airflow logs, Laravel logs, Kafka lag, correlation IDs |
| 13 | [Deployment](13_deployment.md) | Docker, Supervisor, Kubernetes, production checklist |
| 14 | [Capstone Project](14_capstone_project.md) | Order analytics pipeline |

---

## Integration Rule

Use each tool for its natural job:

- **Laravel** owns web requests, business rules, database writes, and user-facing APIs.
- **Kafka** carries durable events between systems.
- **Airflow** orchestrates scheduled, retryable, observable workflows.

Do not turn Airflow into your web app, and do not turn Laravel queue workers into a full workflow scheduler.
