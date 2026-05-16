# Chapter 5: Operators

## PythonOperator

Runs any Python callable. The most flexible operator.

```python
from airflow.operators.python import PythonOperator

def process_data(ds: str, table: str, **context) -> dict:
    print(f"Processing {table} for {ds}")
    # ... business logic ...
    return {"rows_processed": 1000}

task = PythonOperator(
    task_id          = "process",
    python_callable  = process_data,
    op_kwargs        = {"table": "orders"},  # extra keyword args
    op_args          = [],                   # positional args
    provide_context  = True,                 # deprecated in 2.x (always provided)
)
```

---

## BashOperator

Runs shell commands. Supports Jinja templating.

```python
from airflow.operators.bash import BashOperator

run_script = BashOperator(
    task_id      = "run_etl",
    bash_command = """
        set -e
        python /opt/scripts/etl.py \
          --date {{ ds }} \
          --env {{ var.value.environment }}
        echo "Exit code: $?"
    """,
    env          = {"PYTHONPATH": "/opt/scripts"},
    cwd          = "/opt/scripts",
)

# Multi-line with heredoc
run_spark = BashOperator(
    task_id      = "spark_job",
    bash_command = """
        spark-submit \
          --master yarn \
          --deploy-mode cluster \
          --conf spark.executor.memory=4g \
          s3://bucket/jobs/transform.py \
          --date {{ ds }}
    """,
)
```

---

## EmailOperator

```python
from airflow.operators.email import EmailOperator

send_report = EmailOperator(
    task_id = "send_report",
    to      = ["team@company.com", "manager@company.com"],
    subject = "Daily ETL Complete: {{ ds }}",
    html_content = """
        <h2>Pipeline Status</h2>
        <p>The daily ETL for <b>{{ ds }}</b> completed successfully.</p>
        <ul>
            <li>Records processed: {{ ti.xcom_pull(task_ids='load') }}</li>
        </ul>
    """,
    files   = ["/tmp/report_{{ ds }}.csv"],
)
```

---

## EmptyOperator (DummyOperator)

No-op task — useful for grouping, branching joins, and marking pipeline stages.

```python
from airflow.operators.empty import EmptyOperator  # Airflow 2.4+
# or:
from airflow.operators.dummy import DummyOperator  # legacy

start = EmptyOperator(task_id="start")
end   = EmptyOperator(task_id="end", trigger_rule="all_done")

start >> [task_a, task_b, task_c] >> end
```

---

## SparkSubmitOperator

```python
from airflow.providers.apache.spark.operators.spark_submit import SparkSubmitOperator

spark_job = SparkSubmitOperator(
    task_id     = "spark_transform",
    application = "s3://bucket/jobs/transform.py",
    conn_id     = "spark_default",
    application_args = ["--date", "{{ ds }}", "--env", "prod"],
    conf = {
        "spark.executor.memory":        "8g",
        "spark.executor.cores":         "4",
        "spark.executor.instances":     "20",
        "spark.sql.shuffle.partitions": "400",
        "spark.sql.adaptive.enabled":   "true",
    },
    jars         = "s3://bucket/jars/delta-core.jar",
    py_files     = "s3://bucket/packages/utils.zip",
    driver_memory = "4g",
    num_executors = 20,
    verbose      = True,
)
```

---

## BigQueryOperator

```python
from airflow.providers.google.cloud.operators.bigquery import (
    BigQueryInsertJobOperator,
    BigQueryCreateEmptyTableOperator,
)

run_query = BigQueryInsertJobOperator(
    task_id    = "run_bq_query",
    gcp_conn_id = "google_cloud_default",
    configuration = {
        "query": {
            "query": """
                INSERT INTO `project.dataset.output_table`
                SELECT
                    user_id,
                    SUM(amount) AS total_spent,
                    COUNT(*)    AS order_count
                FROM `project.dataset.orders`
                WHERE DATE(created_at) = '{{ ds }}'
                GROUP BY user_id
            """,
            "useLegacySql": False,
            "destinationTable": {
                "projectId": "my-project",
                "datasetId": "analytics",
                "tableId":   "daily_summary",
            },
            "writeDisposition": "WRITE_APPEND",
        }
    },
)
```

---

## PostgresOperator

```python
from airflow.providers.postgres.operators.postgres import PostgresOperator

insert_data = PostgresOperator(
    task_id     = "insert_summary",
    postgres_conn_id = "my_postgres",
    sql         = """
        INSERT INTO daily_summary (date, total_orders, revenue)
        SELECT
            '{{ ds }}'::date,
            COUNT(*),
            SUM(amount)
        FROM orders
        WHERE order_date = '{{ ds }}'
        ON CONFLICT (date) DO UPDATE
            SET total_orders = EXCLUDED.total_orders,
                revenue      = EXCLUDED.revenue;
    """,
)
```

---

## KubernetesPodOperator

Run a Docker container in Kubernetes — isolates dependencies per task:

```python
from airflow.providers.cncf.kubernetes.operators.pod import KubernetesPodOperator
from kubernetes.client import models as k8s

run_ml_task = KubernetesPodOperator(
    task_id        = "train_model",
    name           = "train-model-{{ ds_nodash }}",
    namespace      = "airflow",
    image          = "my-registry/ml-pipeline:latest",
    cmds           = ["python", "train.py"],
    arguments      = ["--date", "{{ ds }}", "--model", "xgboost"],
    env_vars       = [
        k8s.V1EnvVar(name="MLFLOW_TRACKING_URI", value="http://mlflow:5000"),
    ],
    container_resources = k8s.V1ResourceRequirements(
        requests = {"cpu": "2", "memory": "8Gi"},
        limits   = {"cpu": "4", "memory": "16Gi"},
    ),
    is_delete_operator_pod = True,   # clean up pod after run
    get_logs       = True,
    log_events_on_failure = True,
    in_cluster     = True,
)
```

---

## HttpOperator

```python
from airflow.providers.http.operators.http import SimpleHttpOperator
import json

trigger_api = SimpleHttpOperator(
    task_id     = "trigger_refresh",
    http_conn_id = "my_api",
    endpoint    = "/api/v1/refresh",
    method      = "POST",
    data        = json.dumps({"date": "{{ ds }}", "force": True}),
    headers     = {"Content-Type": "application/json", "Authorization": "Bearer {{ var.value.api_token }}"},
    response_check = lambda response: response.json()["status"] == "accepted",
    log_response = True,
)
```

---

## TriggerDagRunOperator

```python
from airflow.operators.trigger_dagrun import TriggerDagRunOperator

trigger_downstream = TriggerDagRunOperator(
    task_id            = "trigger_reporting_dag",
    trigger_dag_id     = "daily_reporting",
    conf               = {"source_date": "{{ ds }}", "table": "orders"},
    wait_for_completion = True,
    execution_date     = "{{ data_interval_end }}",
    reset_dag_run      = True,
)
```

---

## Operator Execution Modes

| Mode | How | Notes |
|---|---|---|
| Default | Runs in worker process | Same Python environment |
| `KubernetesPodOperator` | New K8s pod | Full isolation |
| `DockerOperator` | New Docker container | Isolation without K8s |
| `ExternalTaskSensor` | Waits, not runs | Coordinates across DAGs |
| `BashOperator` subprocesses | Shell process | Can run any executable |
