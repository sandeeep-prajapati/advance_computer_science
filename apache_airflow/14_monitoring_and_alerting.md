# Chapter 14: Monitoring & Alerting

## Airflow Web UI Overview

```
Main views:
  DAGs list     → All DAGs with status, schedule, last run, next run
  Grid view     → Task run history per DAG (rows=tasks, cols=runs)
  Graph view    → DAG dependency visualization
  Gantt chart   → Task duration visualization per run
  Calendar view → Run success/failure per day
  Task Logs     → Per-task stdout + airflow logs
  Code view     → DAG source code snapshot
  Audit Log     → Who changed what
```

---

## Callbacks

Callbacks are Python functions invoked at specific lifecycle events.

### DAG-Level Callbacks

```python
from airflow import DAG
from datetime import datetime

def on_dag_success(context):
    dag_id  = context["dag"].dag_id
    run_id  = context["run_id"]
    print(f"DAG {dag_id} succeeded: {run_id}")

def on_dag_failure(context):
    dag_id     = context["dag"].dag_id
    exception  = context.get("exception")
    print(f"DAG {dag_id} failed: {exception}")

with DAG(
    dag_id               = "monitored_dag",
    start_date           = datetime(2024, 1, 1),
    schedule             = "@daily",
    catchup              = False,
    on_success_callback  = on_dag_success,
    on_failure_callback  = on_dag_failure,
) as dag:
    pass
```

### Task-Level Callbacks

```python
from airflow.operators.python import PythonOperator

def on_task_failure(context):
    task_id  = context["task"].task_id
    dag_id   = context["dag"].dag_id
    ds       = context["ds"]
    exc      = context.get("exception")
    print(f"Task {dag_id}.{task_id} failed on {ds}: {exc}")

def on_task_retry(context):
    ti        = context["task_instance"]
    attempt   = ti.try_number
    print(f"Retry #{attempt} for {ti.task_id}")

def on_task_success(context):
    print("Task succeeded!")

task = PythonOperator(
    task_id             = "my_task",
    python_callable     = lambda: None,
    on_failure_callback = on_task_failure,
    on_retry_callback   = on_task_retry,
    on_success_callback = on_task_success,
)
```

---

## Slack Alerting

```python
from airflow.providers.slack.hooks.slack_webhook import SlackWebhookHook

def slack_alert(context, color: str = "#ff0000", title: str = "Task Failed"):
    ti      = context["task_instance"]
    dag_id  = context["dag"].dag_id
    task_id = ti.task_id
    ds      = context["ds"]
    log_url = ti.log_url
    exc     = context.get("exception", "Unknown error")

    hook = SlackWebhookHook(slack_webhook_conn_id="slack_webhook")
    hook.send(
        attachments=[{
            "color":  color,
            "title":  f"{title} | {dag_id}.{task_id}",
            "fields": [
                {"title": "Date",    "value": ds,      "short": True},
                {"title": "Error",   "value": str(exc), "short": False},
            ],
            "actions": [{
                "type": "button",
                "text": "View Logs",
                "url":  log_url,
            }],
        }]
    )

def slack_failure_callback(context):
    slack_alert(context, color="#ff0000", title=":red_circle: Task Failed")

def slack_success_callback(context):
    slack_alert(context, color="#36a64f", title=":white_check_mark: Task Succeeded")

# Apply as default to all tasks in a DAG
default_args = {
    "on_failure_callback": slack_failure_callback,
    "retries":             2,
}

with DAG("my_dag", default_args=default_args, ...) as dag:
    pass
```

---

## Email Alerting

```python
# airflow.cfg
# [smtp]
# smtp_host = smtp.gmail.com
# smtp_starttls = True
# smtp_ssl = False
# smtp_user = airflow@example.com
# smtp_password = yourpassword
# smtp_port = 587
# smtp_mail_from = airflow@example.com

with DAG(
    dag_id       = "email_alerts_dag",
    schedule     = "@daily",
    start_date   = datetime(2024, 1, 1),
    catchup      = False,
    default_args = {
        "email":            ["team@example.com"],
        "email_on_failure": True,
        "email_on_retry":   False,
        "retries":          2,
    },
) as dag:
    pass
```

```python
# Custom email content
from airflow.utils.email import send_email

def send_custom_alert(context):
    task_id = context["task_instance"].task_id
    dag_id  = context["dag"].dag_id
    ds      = context["ds"]
    log_url = context["task_instance"].log_url

    send_email(
        to      = ["oncall@example.com"],
        subject = f"[AIRFLOW] {dag_id}.{task_id} failed on {ds}",
        html_content = f"""
            <h3>Task Failure Alert</h3>
            <p>DAG: <b>{dag_id}</b></p>
            <p>Task: <b>{task_id}</b></p>
            <p>Date: <b>{ds}</b></p>
            <p><a href="{log_url}">View Logs</a></p>
        """,
    )
```

---

## SLAs (Service Level Agreements)

SLAs trigger a callback if a task hasn't completed within the specified duration:

```python
from datetime import timedelta

def sla_miss_callback(dag, task_list, blocking_task_list, slas, blocking_tis):
    missed = [str(s) for s in slas]
    print(f"SLA missed for tasks: {missed}")
    # Send alert here

with DAG(
    dag_id           = "sla_dag",
    schedule         = "@daily",
    start_date       = datetime(2024, 1, 1),
    catchup          = False,
    sla_miss_callback = sla_miss_callback,
) as dag:

    # Task must complete within 1 hour of the scheduled run time
    extract = PythonOperator(
        task_id         = "extract",
        python_callable = lambda: None,
        sla             = timedelta(hours=1),
    )

    # Task must complete within 2 hours
    load = PythonOperator(
        task_id         = "load",
        python_callable = lambda: None,
        sla             = timedelta(hours=2),
    )
```

---

## Task Instance Timeouts

```python
from datetime import timedelta

task = PythonOperator(
    task_id             = "long_running_task",
    python_callable     = run_heavy_job,
    execution_timeout   = timedelta(hours=2),   # kill task after 2 hours
    retries             = 3,
    retry_delay         = timedelta(minutes=5),
    retry_exponential_backoff = True,           # delay doubles each retry
    max_retry_delay     = timedelta(minutes=60),
)
```

---

## Logs

### Accessing Logs

```bash
# CLI
airflow tasks logs my_dag my_task 2024-01-15

# List log files
airflow tasks logs --full-content my_dag my_task 2024-01-15
```

### Log Configuration

```properties
# airflow.cfg
[logging]
base_log_folder = /opt/airflow/logs
logging_level   = INFO

# Remote logging to S3
remote_logging        = True
remote_base_log_folder = s3://my-airflow-logs/logs
remote_log_conn_id    = aws_default

# Remote logging to GCS
remote_base_log_folder = gs://my-airflow-logs/logs
remote_log_conn_id    = google_cloud_default
```

### Structured Logging in Tasks

```python
import logging

log = logging.getLogger(__name__)

def my_task(**context):
    log.info("Starting task for date: %s", context["ds"])
    log.warning("Potential issue: %s", "low memory")
    log.error("Something went wrong: %s", "connection refused")

    # airflow.task logger is accessible via context
    context["task"].log.info("Also works via task logger")
```

---

## Metrics & Prometheus

Airflow emits StatsD metrics by default; use the Prometheus exporter for Grafana:

```properties
# airflow.cfg
[metrics]
statsd_on   = True
statsd_host = localhost
statsd_port = 8125
statsd_prefix = airflow
```

```yaml
# docker-compose — add Prometheus exporter
prometheus:
  image: prom/prometheus
  volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml

airflow-statsd-exporter:
  image: prom/statsd-exporter
  ports:
    - "9102:9102"
  command:
    - "--statsd.listen-udp=:8125"
    - "--web.listen-address=:9102"
```

### Key Metrics

| Metric | What It Measures |
|---|---|
| `airflow.dag.loading_duration` | Time to load DAG files |
| `airflow.scheduler.heartbeat` | Scheduler is alive |
| `airflow.task.duration` | Task execution time |
| `airflow.task.start_delta` | Delay from scheduled to actual start |
| `airflow.dagrun.duration.success` | DAG run duration (success) |
| `airflow.dagrun.duration.failed` | DAG run duration (failed) |
| `airflow.executor.open_slots` | Available worker slots |
| `airflow.executor.queued_tasks` | Tasks waiting to run |
| `airflow.executor.running_tasks` | Tasks currently running |

---

## Monitoring DAG Health via REST API

```python
import requests

AIRFLOW_URL = "http://localhost:8080"
AUTH        = ("admin", "admin")

def get_dag_runs(dag_id: str, state: str = None) -> list:
    params = {"state": state} if state else {}
    resp   = requests.get(
        f"{AIRFLOW_URL}/api/v1/dags/{dag_id}/dagRuns",
        auth   = AUTH,
        params = params,
    )
    return resp.json()["dag_runs"]

def get_failed_tasks(dag_id: str, run_id: str) -> list:
    resp = requests.get(
        f"{AIRFLOW_URL}/api/v1/dags/{dag_id}/dagRuns/{run_id}/taskInstances",
        auth   = AUTH,
        params = {"state": "failed"},
    )
    return resp.json()["task_instances"]

def clear_task(dag_id: str, run_id: str, task_id: str) -> None:
    requests.post(
        f"{AIRFLOW_URL}/api/v1/dags/{dag_id}/dagRuns/{run_id}/taskInstances/clear",
        auth = AUTH,
        json = {"task_ids": [task_id], "include_downstream": True},
    )
```

---

## Alerting Best Practices

```
1. Always set on_failure_callback for production DAGs
   → Prefer Slack/PagerDuty over email for on-call alerting

2. Set execution_timeout on every task
   → Prevents stuck tasks from blocking the worker pool

3. Use SLAs for data freshness SLOs
   → SLA miss ≠ failure — task may still succeed late

4. Set retry + retry_delay for transient failures
   → retry_exponential_backoff for external API calls

5. Emit custom metrics from tasks
   → airflow.stats.timing("my_dag.extract.rows", row_count)
   → Enables data quality dashboards

6. Log structured data (JSON lines) when possible
   → Enables log aggregation tools (Datadog, CloudWatch, ELK)
```
