# Chapter 17: Airflow with Cloud Providers

## Provider Packages

Cloud integrations are installed as separate provider packages:

```bash
pip install apache-airflow-providers-amazon          # AWS
pip install apache-airflow-providers-google          # GCP
pip install apache-airflow-providers-microsoft-azure # Azure
pip install apache-airflow-providers-databricks      # Databricks
pip install apache-airflow-providers-snowflake       # Snowflake
```

---

## AWS

### Connections

```bash
# Set up AWS connection
airflow connections add aws_default \
  --conn-type aws \
  --conn-extra '{"aws_access_key_id": "AKIA...", "aws_secret_access_key": "...", "region_name": "us-east-1"}'

# Or via environment variable
export AIRFLOW_CONN_AWS_DEFAULT='{"conn_type": "aws", "extra": {"region_name": "us-east-1"}}'
# → Uses IAM role or instance profile automatically when running on AWS
```

### S3 Operations

```python
from airflow.providers.amazon.aws.operators.s3 import S3CopyObjectOperator, S3DeleteObjectsOperator
from airflow.providers.amazon.aws.sensors.s3 import S3KeySensor
from airflow.providers.amazon.aws.hooks.s3 import S3Hook

# Sensor: wait for file
wait = S3KeySensor(
    task_id     = "wait_for_file",
    bucket_name = "my-data-bucket",
    bucket_key  = "raw/{{ ds }}/orders.parquet",
    aws_conn_id = "aws_default",
    mode        = "reschedule",
    timeout     = 3600,
    deferrable  = True,
)

# Copy object
copy = S3CopyObjectOperator(
    task_id             = "copy_to_processed",
    source_bucket_name  = "my-data-bucket",
    source_bucket_key   = "raw/{{ ds }}/orders.parquet",
    dest_bucket_name    = "my-data-bucket",
    dest_bucket_key     = "processed/{{ ds }}/orders.parquet",
    aws_conn_id         = "aws_default",
)

# Use S3Hook directly in a task
from airflow.decorators import task

@task
def upload_results(results: list, ds: str = None) -> str:
    import json, io
    hook = S3Hook(aws_conn_id="aws_default")
    key  = f"output/{ds}/results.json"
    hook.load_string(
        string_data = json.dumps(results),
        key         = key,
        bucket_name = "my-data-bucket",
        replace     = True,
    )
    return f"s3://my-data-bucket/{key}"
```

### Redshift

```python
from airflow.providers.amazon.aws.operators.redshift_sql import RedshiftSQLOperator
from airflow.providers.amazon.aws.transfers.s3_to_redshift import S3ToRedshiftOperator

# Run SQL on Redshift
query = RedshiftSQLOperator(
    task_id      = "transform_redshift",
    redshift_conn_id = "redshift_default",
    sql          = """
        INSERT INTO analytics.daily_orders
        SELECT date, COUNT(*) as order_count, SUM(amount) as total
        FROM staging.orders
        WHERE date = '{{ ds }}'
        GROUP BY date
    """,
)

# Load S3 → Redshift via COPY
load = S3ToRedshiftOperator(
    task_id             = "load_to_redshift",
    schema              = "staging",
    table               = "orders",
    s3_bucket           = "my-data-bucket",
    s3_key              = "raw/{{ ds }}/orders.parquet",
    redshift_conn_id    = "redshift_default",
    aws_conn_id         = "aws_default",
    copy_options        = ["FORMAT AS PARQUET"],
    method              = "REPLACE",
)
```

### EMR (Spark on AWS)

```python
from airflow.providers.amazon.aws.operators.emr import (
    EmrCreateJobFlowOperator,
    EmrAddStepsOperator,
    EmrTerminateJobFlowOperator,
)
from airflow.providers.amazon.aws.sensors.emr import EmrStepSensor

create_cluster = EmrCreateJobFlowOperator(
    task_id         = "create_emr_cluster",
    aws_conn_id     = "aws_default",
    job_flow_overrides = {
        "Name": "airflow-spark-cluster",
        "ReleaseLabel": "emr-7.0.0",
        "Instances": {
            "InstanceGroups": [
                {"Name": "Master", "InstanceRole": "MASTER", "InstanceType": "m5.xlarge", "InstanceCount": 1},
                {"Name": "Workers","InstanceRole": "CORE",   "InstanceType": "m5.xlarge", "InstanceCount": 3},
            ],
            "KeepJobFlowAliveWhenNoSteps": True,
        },
    },
)

add_step = EmrAddStepsOperator(
    task_id       = "run_spark_job",
    job_flow_id   = "{{ task_instance.xcom_pull('create_emr_cluster', key='return_value') }}",
    aws_conn_id   = "aws_default",
    steps         = [{
        "Name": "Transform",
        "ActionOnFailure": "TERMINATE_CLUSTER",
        "HadoopJarStep": {
            "Jar": "command-runner.jar",
            "Args": [
                "spark-submit", "--deploy-mode", "cluster",
                "s3://my-bucket/scripts/transform.py",
                "--date", "{{ ds }}",
            ],
        },
    }],
)

wait_for_step = EmrStepSensor(
    task_id     = "wait_for_spark",
    job_flow_id = "{{ task_instance.xcom_pull('create_emr_cluster', key='return_value') }}",
    step_id     = "{{ task_instance.xcom_pull('run_spark_job', key='return_value')[0] }}",
    aws_conn_id = "aws_default",
)

terminate = EmrTerminateJobFlowOperator(
    task_id    = "terminate_cluster",
    job_flow_id = "{{ task_instance.xcom_pull('create_emr_cluster', key='return_value') }}",
    aws_conn_id = "aws_default",
)

create_cluster >> add_step >> wait_for_step >> terminate
```

---

## GCP (Google Cloud Platform)

### Connection

```bash
airflow connections add google_cloud_default \
  --conn-type google_cloud_platform \
  --conn-extra '{"key_path": "/opt/airflow/gcp_key.json", "project": "my-gcp-project"}'
```

### BigQuery

```python
from airflow.providers.google.cloud.operators.bigquery import (
    BigQueryInsertJobOperator,
    BigQueryCreateEmptyTableOperator,
    BigQueryDeleteTableOperator,
    BigQueryCheckOperator,
)
from airflow.providers.google.cloud.sensors.bigquery import BigQueryTableExistenceSensor

# Run a BQ query
run_query = BigQueryInsertJobOperator(
    task_id    = "run_bq_query",
    gcp_conn_id = "google_cloud_default",
    configuration = {
        "query": {
            "query": """
                SELECT date, COUNT(*) as orders
                FROM `my-project.raw.orders`
                WHERE DATE(created_at) = '{{ ds }}'
                GROUP BY date
            """,
            "useLegacySql":     False,
            "destinationTable": {
                "projectId": "my-project",
                "datasetId": "analytics",
                "tableId":   "daily_orders${{ ds_nodash }}",
            },
            "writeDisposition": "WRITE_TRUNCATE",
        }
    },
)

# Check row count > 0
check = BigQueryCheckOperator(
    task_id    = "check_data",
    sql        = "SELECT COUNT(*) FROM `my-project.analytics.daily_orders` WHERE date = '{{ ds }}'",
    gcp_conn_id = "google_cloud_default",
    use_legacy_sql = False,
)
```

### GCS Operations

```python
from airflow.providers.google.cloud.operators.gcs import GCSDeleteObjectsOperator
from airflow.providers.google.cloud.transfers.gcs_to_bigquery import GCSToBigQueryOperator
from airflow.providers.google.cloud.sensors.gcs import GCSObjectExistenceSensor

# Sensor: wait for file in GCS
wait = GCSObjectExistenceSensor(
    task_id    = "wait_for_gcs",
    bucket     = "my-gcs-bucket",
    object     = "raw/{{ ds }}/orders.parquet",
    gcp_conn_id = "google_cloud_default",
    mode       = "reschedule",
)

# Load GCS → BigQuery
load = GCSToBigQueryOperator(
    task_id              = "gcs_to_bq",
    bucket               = "my-gcs-bucket",
    source_objects       = ["raw/{{ ds }}/orders.parquet"],
    destination_project_dataset_table = "my-project.staging.orders",
    source_format        = "PARQUET",
    write_disposition    = "WRITE_TRUNCATE",
    gcp_conn_id          = "google_cloud_default",
)
```

### Dataproc (Spark on GCP)

```python
from airflow.providers.google.cloud.operators.dataproc import (
    DataprocCreateClusterOperator,
    DataprocSubmitJobOperator,
    DataprocDeleteClusterOperator,
)

create = DataprocCreateClusterOperator(
    task_id      = "create_cluster",
    project_id   = "my-project",
    cluster_name = "airflow-spark-{{ ds_nodash }}",
    region       = "us-central1",
    cluster_config = {
        "master_config": {"num_instances": 1, "machine_type_uri": "n1-standard-4"},
        "worker_config": {"num_instances": 2, "machine_type_uri": "n1-standard-4"},
    },
    gcp_conn_id  = "google_cloud_default",
)

submit_job = DataprocSubmitJobOperator(
    task_id    = "run_spark",
    project_id = "my-project",
    region     = "us-central1",
    job        = {
        "placement":  {"cluster_name": "airflow-spark-{{ ds_nodash }}"},
        "pyspark_job": {
            "main_python_file_uri": "gs://my-bucket/scripts/transform.py",
            "args": ["--date", "{{ ds }}"],
        },
    },
    gcp_conn_id = "google_cloud_default",
)

delete = DataprocDeleteClusterOperator(
    task_id      = "delete_cluster",
    project_id   = "my-project",
    cluster_name = "airflow-spark-{{ ds_nodash }}",
    region       = "us-central1",
    gcp_conn_id  = "google_cloud_default",
    trigger_rule = "all_done",   # delete even if spark job failed
)

create >> submit_job >> delete
```

---

## Azure

### Azure Blob Storage

```python
from airflow.providers.microsoft.azure.operators.wasb_delete_blob import WasbDeleteBlobOperator
from airflow.providers.microsoft.azure.sensors.wasb import WasbBlobSensor
from airflow.providers.microsoft.azure.hooks.wasb import WasbHook

# Sensor: wait for blob
wait = WasbBlobSensor(
    task_id       = "wait_for_blob",
    container_name = "raw-data",
    blob_name     = "orders/{{ ds }}/data.parquet",
    wasb_conn_id  = "wasb_default",
    mode          = "reschedule",
)

# Use hook directly
@task
def upload_to_blob(data: bytes, ds: str = None) -> str:
    hook = WasbHook(wasb_conn_id="wasb_default")
    blob_name = f"output/{ds}/results.parquet"
    hook.load_bytes(data, container_name="output", blob_name=blob_name, overwrite=True)
    return blob_name
```

### Azure Synapse / SQL

```python
from airflow.providers.microsoft.azure.operators.synapse import AzureSynapseRunSparkBatchOperator

run_spark = AzureSynapseRunSparkBatchOperator(
    task_id     = "run_synapse_spark",
    azure_synapse_conn_id = "synapse_default",
    spark_pool  = "my-spark-pool",
    payload     = {
        "name":           "transform_orders",
        "file":           "abfss://scripts@mystorageaccount.dfs.core.windows.net/transform.py",
        "args":           ["--date", "{{ ds }}"],
        "driverMemory":   "4g",
        "executorMemory": "4g",
        "numExecutors":   4,
    },
)
```

---

## Databricks

```python
from airflow.providers.databricks.operators.databricks import DatabricksRunNowOperator, DatabricksSubmitRunOperator

# Trigger an existing Databricks job
run_job = DatabricksRunNowOperator(
    task_id          = "run_databricks_job",
    databricks_conn_id = "databricks_default",
    job_id           = 12345,
    notebook_params  = {"date": "{{ ds }}", "env": "production"},
)

# Submit a one-time run
submit = DatabricksSubmitRunOperator(
    task_id          = "submit_notebook",
    databricks_conn_id = "databricks_default",
    new_cluster      = {
        "spark_version": "13.3.x-scala2.12",
        "node_type_id":  "i3.xlarge",
        "num_workers":   4,
    },
    notebook_task    = {
        "notebook_path": "/Shared/ETL/transform_orders",
        "base_parameters": {"date": "{{ ds }}"},
    },
)
```

---

## Snowflake

```python
from airflow.providers.snowflake.operators.snowflake import SnowflakeOperator
from airflow.providers.snowflake.transfers.s3_to_snowflake import S3ToSnowflakeOperator

# Run SQL
query = SnowflakeOperator(
    task_id        = "transform_snowflake",
    snowflake_conn_id = "snowflake_default",
    sql            = """
        MERGE INTO ANALYTICS.ORDERS t
        USING STAGING.ORDERS s ON t.ORDER_ID = s.ORDER_ID
        WHEN MATCHED THEN UPDATE SET t.STATUS = s.STATUS
        WHEN NOT MATCHED THEN INSERT VALUES (s.ORDER_ID, s.STATUS, s.AMOUNT)
    """,
    warehouse      = "COMPUTE_WH",
    database       = "MY_DB",
    schema         = "ANALYTICS",
)

# Load from S3
load = S3ToSnowflakeOperator(
    task_id           = "s3_to_snowflake",
    s3_keys           = ["raw/{{ ds }}/orders.parquet"],
    table             = "ORDERS",
    schema            = "STAGING",
    stage             = "MY_S3_STAGE",
    file_format       = "PARQUET_FORMAT",
    snowflake_conn_id = "snowflake_default",
)
```
