# Chapter 6: Sensors

## What are Sensors?

Sensors are operators that **poll** an external condition and block the pipeline until the condition is met. They are commonly used to:
- Wait for data files to arrive
- Wait for another pipeline to complete
- Wait for an API endpoint to be available
- Wait for a database row or result to appear

```
Sensor Execution Modes:

poke mode (default):
  Worker holds a slot while polling every poke_interval
  Uses a worker slot the entire time it's waiting
  Good for: short waits (< 5 minutes)

reschedule mode (recommended for long waits):
  Worker releases slot between checks
  Scheduler reschedules it at poke_interval
  Good for: long waits (hours)
  AIRFLOW__SCHEDULER__MIN_FILE_PROCESS_INTERVAL must be < poke_interval
```

---

## FileSensor

Wait for a file or directory to appear on disk:

```python
from airflow.sensors.filesystem import FileSensor

wait_for_csv = FileSensor(
    task_id       = "wait_for_input",
    filepath      = "/data/input/{{ ds }}/orders.csv",
    poke_interval = 60,      # check every 60 seconds
    timeout       = 7200,    # fail after 2 hours
    mode          = "reschedule",
    soft_fail     = False,   # True = skip instead of fail on timeout
)
```

---

## S3KeySensor

Wait for an object to appear in Amazon S3:

```python
from airflow.providers.amazon.aws.sensors.s3 import S3KeySensor

wait_for_data = S3KeySensor(
    task_id        = "wait_for_s3_data",
    bucket_name    = "my-data-bucket",
    bucket_key     = "raw/{{ ds }}/orders/*.parquet",
    wildcard_match = True,
    aws_conn_id    = "aws_default",
    poke_interval  = 120,
    timeout        = 3600,
    mode           = "reschedule",
)

# Check multiple keys
wait_for_all = S3KeySensor(
    task_id    = "wait_for_all_files",
    bucket_name = "my-bucket",
    bucket_key  = [
        "raw/{{ ds }}/orders.parquet",
        "raw/{{ ds }}/customers.parquet",
        "raw/{{ ds }}/products.parquet",
    ],
    aws_conn_id = "aws_default",
    mode        = "reschedule",
)
```

---

## HttpSensor

Wait for an HTTP endpoint to return a success response:

```python
from airflow.providers.http.sensors.http import HttpSensor

wait_for_api = HttpSensor(
    task_id         = "wait_for_api",
    http_conn_id    = "my_api",
    endpoint        = "/health",
    request_params  = {},
    response_check  = lambda response: response.json()["status"] == "healthy",
    poke_interval   = 30,
    timeout         = 300,
    mode            = "reschedule",
)

# Wait for a data endpoint
wait_for_report = HttpSensor(
    task_id         = "wait_for_daily_report",
    http_conn_id    = "reporting_api",
    endpoint        = "/api/reports/{{ ds }}",
    response_check  = lambda r: r.status_code == 200 and r.json()["ready"] is True,
    poke_interval   = 60,
    timeout         = 3600,
)
```

---

## SqlSensor

Wait for a SQL query to return rows:

```python
from airflow.providers.common.sql.sensors.sql import SqlSensor

wait_for_data = SqlSensor(
    task_id      = "wait_for_loaded_data",
    conn_id      = "my_postgres",
    sql          = """
        SELECT COUNT(*) FROM orders
        WHERE order_date = '{{ ds }}'
        HAVING COUNT(*) > 0
    """,
    poke_interval = 60,
    timeout       = 3600,
    mode          = "reschedule",
    success       = lambda rows: rows and rows[0][0] > 0,
)
```

---

## ExternalTaskSensor

Wait for a task in another DAG to complete:

```python
from airflow.sensors.external_task import ExternalTaskSensor

wait_for_upstream = ExternalTaskSensor(
    task_id             = "wait_for_source_dag",
    external_dag_id     = "source_data_pipeline",
    external_task_id    = "load_complete",   # None = wait for whole DAG
    execution_date_fn   = lambda dt: dt,     # map this DAG's run to source DAG's run
    mode                = "reschedule",
    timeout             = 3600,
    poke_interval       = 60,
    allowed_states      = ["success"],
    failed_states       = ["failed", "skipped"],
)

# Wait for yesterday's run of upstream DAG
from datetime import timedelta

wait_for_yesterday = ExternalTaskSensor(
    task_id          = "wait_for_yesterday",
    external_dag_id  = "daily_etl",
    execution_date_fn = lambda dt: dt - timedelta(days=1),
    mode             = "reschedule",
)
```

---

## PythonSensor

Wait for any custom Python condition:

```python
from airflow.sensors.python import PythonSensor

def check_kafka_topic(**context) -> bool:
    from confluent_kafka.admin import AdminClient
    admin = AdminClient({"bootstrap.servers": "kafka:9092"})
    metadata = admin.list_topics(timeout=5)
    topic = f"orders-{context['ds_nodash']}"
    return topic in metadata.topics

wait_for_kafka = PythonSensor(
    task_id         = "wait_for_kafka_topic",
    python_callable = check_kafka_topic,
    poke_interval   = 30,
    timeout         = 600,
    mode            = "reschedule",
)
```

---

## Deferred Sensors (Deferrable Operators)

Airflow 2.2+ supports async deferral — the sensor completely frees the worker while waiting, using a Triggerer process:

```python
from airflow.providers.amazon.aws.sensors.s3 import S3KeySensor

# Just add deferrable=True
wait_for_data = S3KeySensor(
    task_id     = "wait_for_s3",
    bucket_name = "my-bucket",
    bucket_key  = "data/{{ ds }}/input.csv",
    deferrable  = True,          # uses Triggerer, frees worker slot
    poke_interval = timedelta(seconds=60),
    timeout     = timedelta(hours=6),
)
```

**Benefits of deferrable:**
- Worker slot freed while waiting (much more scalable)
- Triggerer uses async I/O — handles thousands of deferred sensors with minimal resources
- Available for most built-in sensors (S3, GCS, BigQuery, etc.)

---

## Custom Sensor

```python
from airflow.sensors.base import BaseSensorOperator
from airflow.utils.decorators import apply_defaults

class DataAvailableSensor(BaseSensorOperator):
    """Waits until data is available in our internal API."""

    template_fields = ("endpoint",)

    def __init__(self, endpoint: str, min_records: int = 1, **kwargs):
        super().__init__(**kwargs)
        self.endpoint    = endpoint
        self.min_records = min_records

    def poke(self, context) -> bool:
        import requests
        response = requests.get(
            f"http://data-api/availability/{self.endpoint}",
            params={"date": context["ds"]},
            timeout=10,
        )
        if response.status_code != 200:
            return False
        data = response.json()
        available = data.get("record_count", 0) >= self.min_records
        self.log.info(f"Records available: {data.get('record_count', 0)}, needed: {self.min_records}")
        return available

# Usage
wait_for_orders = DataAvailableSensor(
    task_id     = "wait_for_orders",
    endpoint    = "orders",
    min_records = 100,
    poke_interval = 120,
    timeout     = 3600,
    mode        = "reschedule",
)
```

---

## Sensor Best Practices

```
1. Always use mode="reschedule" for waits > 5 minutes
   → Prevents worker slot exhaustion

2. Set reasonable timeouts
   → soft_fail=True: skip downstream instead of fail

3. Use deferrable=True when available (Airflow 2.2+)
   → Most scalable option for long-running sensors

4. Avoid tight poke_interval (< 30s)
   → Creates unnecessary load on external systems

5. Set failed_states in ExternalTaskSensor
   → Otherwise waits forever if upstream fails
```
