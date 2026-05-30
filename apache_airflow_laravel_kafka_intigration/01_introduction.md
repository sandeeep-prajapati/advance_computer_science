# Lesson 1: Introduction

## Goal

This guide teaches how to integrate Apache Airflow with a PHP/Laravel application that uses Kafka.

The common pattern is:

```text
Laravel publishes business events
Kafka stores and distributes those events
Airflow runs scheduled or triggered workflows
Laravel consumes results or receives callbacks
```

---

## When This Integration Makes Sense

Use Airflow with Laravel and Kafka when you need:

- scheduled exports from a Laravel database
- nightly reports or analytics pipelines
- long-running data processing outside HTTP requests
- retryable workflows with dependencies
- backfills for historical data
- coordination between Laravel, Kafka, storage, and data warehouses

---

## When Not to Use Airflow

Airflow is not a replacement for:

- Laravel controllers
- Laravel queues
- real-time request handling
- simple cron jobs
- short background jobs

If a task is small and internal to Laravel, use Laravel queues or Laravel Scheduler. Use Airflow when workflows are multi-step, observable, retryable, and operationally important.

---

## Main Integration Styles

| Style | Direction | Example |
|---|---|---|
| Laravel triggers Airflow | Laravel -> Airflow API | User uploads CSV and starts processing DAG |
| Airflow calls Laravel | Airflow -> Laravel API | DAG asks Laravel for pending orders |
| Laravel publishes Kafka | Laravel -> Kafka | `orders.created` event |
| Airflow reads Kafka | Kafka -> Airflow | DAG processes event batches |
| Airflow publishes Kafka | Airflow -> Kafka | `reports.generated` event |
| Laravel consumes Kafka | Kafka -> Laravel | Laravel updates report status |
