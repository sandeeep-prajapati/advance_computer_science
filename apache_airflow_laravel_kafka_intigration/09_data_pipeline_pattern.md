# Lesson 9: Data Pipeline Pattern

## Scenario

Laravel stores orders. Airflow generates daily revenue reports. Kafka tells Laravel when reports are ready.

```text
Laravel DB -> Airflow DAG -> CSV/Parquet/S3 -> Kafka event -> Laravel UI status
```

---

## DAG Shape

```text
extract_orders
  -> transform_orders
  -> build_report
  -> publish_reports_generated
```

---

## Extract from Laravel API

```python
import requests

def fetch_orders(base_url, token, start_date, end_date):
    response = requests.get(
        f"{base_url}/internal/airflow/orders",
        headers={"Authorization": f"Bearer {token}"},
        params={"from": start_date, "to": end_date},
        timeout=60,
    )
    response.raise_for_status()
    return response.json()
```

---

## Publish Result Event

```json
{
  "event_id": "01J...",
  "event_type": "reports.generated",
  "event_version": 1,
  "report_id": 55,
  "file_path": "s3://bucket/reports/orders-2026-05-30.csv",
  "generated_at": "2026-05-30T12:00:00Z"
}
```

Laravel consumes this event and updates the report row.

---

## Why This Pattern Works

- Laravel remains the source of business truth.
- Airflow owns heavy data workflow execution.
- Kafka decouples report generation from Laravel status updates.
- The UI can show progress without running the workflow in a web request.
