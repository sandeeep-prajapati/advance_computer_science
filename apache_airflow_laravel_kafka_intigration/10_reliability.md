# Lesson 10: Reliability

## Idempotency Everywhere

Airflow tasks can retry. Kafka messages can be duplicated. HTTP calls can time out after succeeding.

Design every step to be safe to repeat.

Laravel examples:

```php
Report::updateOrCreate(
    ['id' => $payload['report_id']],
    ['status' => 'ready', 'file_path' => $payload['file_path']]
);
```

Airflow examples:

```python
# write output to deterministic path
file_path = f"s3://reports/orders/{report_id}.csv"
```

---

## Retries

Airflow retry config:

```python
@dag(
    default_args={
        "retries": 3,
        "retry_delay": timedelta(minutes=5),
    },
)
def report_pipeline():
    ...
```

Laravel Kafka consumer retries should have limits and DLQ behavior.

---

## Dead Letter Topics

Use DLQ topics for messages Laravel cannot process:

```text
reports.generated.dlq
orders.created.dlq
```

Include:

- original topic
- original key
- payload
- error message
- failed timestamp

---

## Backfills

Airflow is excellent for backfills:

```text
Run daily_report for 2026-05-01 through 2026-05-30
```

Make sure Laravel APIs can accept date ranges and deterministic report IDs so backfills do not create duplicate business records.
