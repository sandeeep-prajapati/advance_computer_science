# Lesson 8: Airflow and Kafka

## Airflow Publishing to Kafka

Install a Kafka Python client in the Airflow image, such as `confluent-kafka`.

```python
from confluent_kafka import Producer
import json
import os

def publish_report_generated(report_id: int, file_path: str):
    producer = Producer({
        "bootstrap.servers": os.environ["KAFKA_BOOTSTRAP_SERVERS"],
        "acks": "all",
    })

    producer.produce(
        "reports.generated",
        key=f"report-{report_id}",
        value=json.dumps({
            "event_type": "reports.generated",
            "event_version": 1,
            "report_id": report_id,
            "file_path": file_path,
        }),
    )
    producer.flush()
```

---

## Airflow Task

```python
from airflow.decorators import task

@task
def notify_laravel(report_id: int, file_path: str):
    publish_report_generated(report_id, file_path)
```

Laravel then consumes `reports.generated`.

---

## Should Airflow Consume Kafka Continuously?

Airflow is best for scheduled or finite workflows. A never-ending Kafka consumer is usually better as:

- Laravel consumer command
- standalone Python service
- Kafka Connect
- Flink or Spark streaming job

Airflow can read a bounded batch from Kafka, but it should not normally run as a permanent stream processor.

---

## Bounded Kafka Batch Pattern

```text
Airflow DAG every 15 minutes
  -> read offsets from tracking table
  -> consume up to N messages
  -> process batch
  -> save latest offsets
```

Use this for batch analytics, not low-latency event processing.
