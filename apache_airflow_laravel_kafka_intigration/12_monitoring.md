# Lesson 12: Monitoring

## Correlation IDs

Generate one correlation ID in Laravel and pass it through:

```text
Laravel request
  -> Kafka headers
  -> Airflow dag_run.conf
  -> Airflow logs
  -> reports.generated event
  -> Laravel consumer logs
```

This makes one workflow traceable across all systems.

---

## Laravel Logs

Log:

- event type
- event ID
- topic
- key
- correlation ID
- Airflow DAG ID
- Airflow DAG run ID

---

## Airflow Monitoring

Watch:

- DAG failures
- task duration
- retry count
- queued task count
- scheduler health
- worker health

Airflow should alert when a business-critical DAG fails.

---

## Kafka Monitoring

Watch:

- consumer lag
- broker health
- produce failures
- DLQ message count
- topic size and retention

Consumer lag connects directly to business freshness. If lag grows, users may see stale report status or delayed updates.
