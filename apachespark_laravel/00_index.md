# Apache Spark + Laravel — Lessons

> Practical lessons for integrating a PHP/Laravel application with Apache Spark for batch processing, analytics, and data pipelines.

---

## Table of Contents

| Lesson | Title | Topics |
|---|---|---|
| 01 | [Introduction](01_introduction.md) | Where Spark fits, Laravel's role, integration options |
| 02 | [Architecture](02_architecture.md) | Laravel app, object storage, Spark cluster, result callbacks |
| 03 | [Local Setup](03_local_setup.md) | Docker services, Laravel config, Spark master/worker |
| 04 | [Data Exchange Patterns](04_data_exchange_patterns.md) | Database, CSV, Parquet, S3/MinIO, APIs |
| 05 | [Triggering Spark from Laravel](05_triggering_spark_from_laravel.md) | `spark-submit`, REST, queues, job records |
| 06 | [Using Apache Livy](06_using_apache_livy.md) | Submit Spark jobs through HTTP from Laravel |
| 07 | [Laravel Export Jobs](07_laravel_export_jobs.md) | Exporting source data for Spark |
| 08 | [Spark Job Examples](08_spark_job_examples.md) | PySpark jobs for reports and transformations |
| 09 | [Reading Spark Results in Laravel](09_reading_spark_results.md) | Status tables, result files, callbacks |
| 10 | [Scheduling](10_scheduling.md) | Laravel Scheduler vs Spark jobs vs Airflow |
| 11 | [Reliability](11_reliability.md) | Idempotency, retries, deterministic output, backfills |
| 12 | [Security](12_security.md) | Service tokens, storage credentials, secrets |
| 13 | [Monitoring and Debugging](13_monitoring_debugging.md) | Spark UI, Laravel logs, job IDs, metrics |
| 14 | [Deployment](14_deployment.md) | Standalone, Kubernetes, EMR/Dataproc, production checklist |
| 15 | [Capstone Project](15_capstone_project.md) | Customer revenue analytics pipeline |

---

## Main Idea

Laravel should not process huge datasets inside web requests. Laravel should create jobs, export or point to data, trigger Spark, and display results. Spark should do distributed computation.
