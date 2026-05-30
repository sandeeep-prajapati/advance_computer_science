# Lesson 1: Introduction

## What Problem Are We Solving?

Laravel is excellent for web apps, APIs, business rules, forms, dashboards, authentication, and transactional database workflows.

Apache Spark is excellent for processing large datasets:

- millions or billions of rows
- joins across large files
- analytics reports
- ETL pipelines
- machine learning features
- batch transformations

The integration goal is simple:

```text
Laravel starts and tracks work
Spark processes large data
Laravel shows the result
```

---

## Laravel Should Own

- users and permissions
- business actions
- database transactions
- API validation
- job request records
- report status pages
- notifications

---

## Spark Should Own

- large file processing
- distributed joins
- aggregations
- ETL
- report generation
- feature computation
- batch analytics

---

## Integration Options

| Option | Best For |
|---|---|
| Laravel runs `spark-submit` | Simple internal deployments |
| Laravel calls Apache Livy | HTTP-based Spark job submission |
| Laravel writes job row, worker submits Spark | Reliable async execution |
| Airflow schedules Spark and calls Laravel | Multi-step data workflows |
| Laravel publishes event, pipeline consumes | Event-driven analytics |

PHP does not need to run Spark directly. Laravel usually talks to Spark through process execution, HTTP APIs, queues, files, or orchestration tools.
