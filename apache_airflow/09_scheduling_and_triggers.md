# Chapter 9: Scheduling & Triggers

## Cron-Based Scheduling

```python
with DAG(
    dag_id   = "my_dag",
    schedule = "0 6 * * *",   # daily at 6 AM UTC
    ...
)
```

### Cron Syntax

```
┌───────────── minute (0–59)
│ ┌─────────── hour (0–23)
│ │ ┌───────── day of month (1–31)
│ │ │ ┌─────── month (1–12)
│ │ │ │ ┌───── day of week (0=Sun–6=Sat)
│ │ │ │ │
* * * * *

0 6 * * *       Every day at 6 AM
0 */4 * * *     Every 4 hours
0 6 * * 1-5     Weekdays at 6 AM
30 2 1 * *      1st of every month at 2:30 AM
0 0 * * 0       Every Sunday midnight
*/15 * * * *    Every 15 minutes
0 0 1 1 *       Once a year (Jan 1st)
```

### Preset Schedules

| Preset | Equivalent Cron | When |
|---|---|---|
| `@once` | (run once) | Single run at start_date |
| `@hourly` | `0 * * * *` | Every hour |
| `@daily` | `0 0 * * *` | Midnight daily |
| `@weekly` | `0 0 * * 0` | Sunday midnight |
| `@monthly` | `0 0 1 * *` | 1st of month |
| `@yearly` | `0 0 1 1 *` | January 1st |
| `None` | — | Manual trigger only |

---

## timedelta Scheduling

```python
from datetime import timedelta

# Run every 6 hours
dag = DAG(schedule=timedelta(hours=6), ...)

# Run every 30 minutes
dag = DAG(schedule=timedelta(minutes=30), ...)
```

---

## Data Interval

In Airflow 2.2+, every scheduled run has a **data interval** — the period of data it processes:

```
Schedule: @daily

DAG Run triggered after 2024-01-03 00:00:00:
  data_interval_start = 2024-01-02 00:00:00
  data_interval_end   = 2024-01-03 00:00:00

The run processes data FROM Jan 2 TO Jan 3.
{{ ds }}            = "2024-01-02"   (data_interval_start date)
{{ data_interval_end | ds }} = "2024-01-03"
```

---

## Catchup

`catchup=True` (default): Airflow creates a DAG run for every missed interval since `start_date`.

```python
# Dangerous if start_date is far in the past!
dag = DAG(
    start_date = datetime(2020, 1, 1),
    schedule   = "@daily",
    catchup    = True,   # creates 4+ years of runs!
)

# Production best practice: disable catchup
dag = DAG(
    start_date = datetime(2024, 1, 1),
    schedule   = "@daily",
    catchup    = False,
)
```

---

## Backfill

Manually run a DAG for historical intervals:

```bash
# Backfill one week of data
airflow dags backfill \
  --start-date 2024-01-01 \
  --end-date   2024-01-07 \
  my_dag

# With dry run (don't actually run)
airflow dags backfill --dry-run \
  --start-date 2024-01-01 --end-date 2024-01-07 \
  my_dag

# Mark as successful (skip actual execution)
airflow dags backfill --mark-success \
  --start-date 2024-01-01 --end-date 2024-01-07 \
  my_dag
```

---

## Custom Timetables (Airflow 2.2+)

For non-standard schedules (e.g., business days only, every weekday at market open):

```python
from airflow.timetables.base import DagRunInfo, DataInterval, Timetable
from pendulum import DateTime, timezone
import pendulum

class WeekdayTimetable(Timetable):
    """Run once per weekday (Mon–Fri)."""

    description = "Runs on weekdays at 6 AM UTC"

    def infer_manual_data_interval(self, run_after: DateTime) -> DataInterval:
        start = run_after.start_of("day")
        return DataInterval(start=start, end=start.add(days=1))

    def next_dagrun_info(self, *, last_automated_data_interval, restriction):
        if last_automated_data_interval is None:
            next_start = restriction.earliest
        else:
            next_start = last_automated_data_interval.end

        # Skip weekends
        while next_start.day_of_week >= 5:  # Sat=5, Sun=6
            next_start = next_start.add(days=1)

        run_after = next_start.add(hours=6)  # 6 AM
        return DagRunInfo.interval(
            start=next_start,
            end=next_start.add(days=1),
        )
```

Register in a plugin and use:
```python
dag = DAG(schedule=WeekdayTimetable(), ...)
```

---

## Dataset-Driven Scheduling (Airflow 2.4+)

DAGs can be triggered when other DAGs produce datasets, enabling event-driven pipelines:

```python
from airflow import DAG, Dataset

# Define datasets (URIs, logical identifiers)
orders_dataset   = Dataset("s3://bucket/raw/orders/")
products_dataset = Dataset("s3://bucket/raw/products/")
output_dataset   = Dataset("s3://bucket/processed/sales/")

# Producer DAG: declares that it updates a dataset
with DAG("load_orders", schedule="@daily") as dag:
    @task(outlets=[orders_dataset])   # this task produces orders_dataset
    def load_orders():
        # ... load orders to S3 ...
        pass

# Consumer DAG: triggered when BOTH datasets are updated
with DAG(
    "process_sales",
    schedule = [orders_dataset, products_dataset],  # triggered by datasets
) as dag:
    @task(outlets=[output_dataset])
    def process(ds=None):
        # Runs only after both orders and products have been updated
        pass
```

---

## Manual Triggers

```bash
# Trigger via CLI
airflow dags trigger my_dag

# With configuration
airflow dags trigger my_dag --conf '{"env": "production", "date": "2024-01-15"}'

# With a specific run_id
airflow dags trigger my_dag --run-id "manual__2024-01-15T10:00:00"
```

```python
# Trigger via REST API
import requests

response = requests.post(
    "http://airflow:8080/api/v1/dags/my_dag/dagRuns",
    auth    = ("admin", "admin"),
    json    = {
        "conf":     {"env": "production", "date": "2024-01-15"},
        "run_id":   "api_trigger__2024-01-15",
    },
)
print(response.json())
```

---

## Schedule Summary

| Pattern | Use Case | Example |
|---|---|---|
| Cron string | Recurring schedule | `"0 2 * * *"` |
| `@daily`, `@hourly` | Common presets | `schedule="@daily"` |
| `timedelta` | Fixed intervals | `schedule=timedelta(hours=4)` |
| `None` | Manual only | `schedule=None` |
| Dataset | Event-driven | `schedule=[my_dataset]` |
| Custom timetable | Complex schedules | Business days, market hours |
