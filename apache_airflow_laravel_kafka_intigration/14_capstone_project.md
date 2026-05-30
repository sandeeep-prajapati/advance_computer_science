# Lesson 14: Capstone Project

## Build: Order Analytics Pipeline

Create a Laravel + Kafka + Airflow workflow:

```text
1. User requests an order report in Laravel.
2. Laravel creates report row with status=pending.
3. Laravel triggers Airflow DAG with report_id and date range.
4. Airflow calls Laravel internal API to fetch orders.
5. Airflow builds CSV report.
6. Airflow publishes reports.generated to Kafka.
7. Laravel consumes reports.generated.
8. Laravel marks report status=ready.
```

---

## Laravel Pieces

- `reports` table
- `POST /reports`
- `GET /internal/airflow/orders`
- `POST /internal/airflow/reports/{report}/failed`
- `AirflowClient`
- `kafka:consume-reports`
- `ReportGeneratedHandler`

---

## Airflow Pieces

- DAG `generate_order_report`
- task `read_conf`
- task `fetch_orders_from_laravel`
- task `build_csv`
- task `publish_reports_generated`
- failure callback that notifies Laravel

---

## Kafka Topics

```text
reports.generated
reports.generated.dlq
airflow.pipeline.failed
```

---

## Acceptance Criteria

- Laravel can trigger a DAG.
- Airflow can read `dag_run.conf`.
- Airflow can call Laravel internal API.
- Airflow publishes a Kafka result event.
- Laravel consumes the result and updates report status.
- Duplicate Kafka messages do not duplicate report rows.
- Failed workflows are visible in Laravel and Airflow.
