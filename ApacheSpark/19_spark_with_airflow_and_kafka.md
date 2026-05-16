# Chapter 19: Spark with Airflow & Kafka

## Spark + Apache Airflow

Apache Airflow is the standard workflow orchestrator for data pipelines. It schedules and monitors Spark jobs.

### SparkSubmitOperator

```python
from airflow import DAG
from airflow.providers.apache.spark.operators.spark_submit import SparkSubmitOperator
from datetime import datetime, timedelta

default_args = {
    "owner":            "data-team",
    "retries":          2,
    "retry_delay":      timedelta(minutes=5),
    "email_on_failure": True,
    "email":            ["alerts@company.com"],
}

with DAG(
    dag_id        = "daily_etl",
    default_args  = default_args,
    start_date    = datetime(2024, 1, 1),
    schedule      = "0 2 * * *",    # 2 AM daily
    catchup       = False,
    tags          = ["spark", "etl"],
) as dag:

    ingest_orders = SparkSubmitOperator(
        task_id         = "ingest_orders",
        application     = "s3://bucket/jobs/ingest_orders.py",
        conn_id         = "spark_default",
        application_args= ["--date", "{{ ds }}", "--env", "production"],
        conf={
            "spark.executor.memory":          "8g",
            "spark.executor.cores":           "4",
            "spark.executor.instances":       "20",
            "spark.sql.shuffle.partitions":   "400",
            "spark.sql.adaptive.enabled":     "true",
        },
        jars            = "s3://bucket/jars/delta-core.jar",
        py_files        = "s3://bucket/packages/utils.zip",
        driver_memory   = "4g",
        verbose         = True,
    )

    transform_orders = SparkSubmitOperator(
        task_id     = "transform_orders",
        application = "s3://bucket/jobs/transform_orders.py",
        conn_id     = "spark_default",
        conf        = {"spark.executor.memory": "8g"},
    )

    load_to_warehouse = SparkSubmitOperator(
        task_id     = "load_to_warehouse",
        application = "s3://bucket/jobs/load_warehouse.py",
        conn_id     = "spark_default",
    )

    ingest_orders >> transform_orders >> load_to_warehouse
```

### Airflow Spark Connection

Configure in Airflow UI or via environment:

```bash
# Airflow connection URI format
AIRFLOW_CONN_SPARK_DEFAULT="spark://spark-master:7077?deploy-mode=cluster"

# For YARN
AIRFLOW_CONN_SPARK_DEFAULT="yarn://"
```

### Passing Dates and Parameters

```python
# In the Spark job: parse Airflow-passed args
import argparse
from pyspark.sql import SparkSession

def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("--date",        required=True)
    parser.add_argument("--env",         default="dev")
    parser.add_argument("--output-path", required=True)
    return parser.parse_args()

def main():
    args  = parse_args()
    spark = SparkSession.builder.appName(f"ETL_{args.date}").getOrCreate()

    df = spark.read.parquet(f"s3://bucket/raw/{args.date}/")
    # ... transform ...
    df.write.mode("overwrite").parquet(f"{args.output_path}/{args.date}/")

if __name__ == "__main__":
    main()
```

### Sensor: Wait for Data to Arrive

```python
from airflow.providers.amazon.aws.sensors.s3 import S3KeySensor

wait_for_data = S3KeySensor(
    task_id        = "wait_for_input",
    bucket_name    = "my-bucket",
    bucket_key     = "raw/{{ ds }}/_SUCCESS",
    aws_conn_id    = "aws_default",
    poke_interval  = 60,        # check every minute
    timeout        = 3600,      # fail after 1 hour
    mode           = "reschedule",
)

wait_for_data >> ingest_orders
```

---

## Spark + Apache Kafka

### Read from Kafka (Batch)

```python
spark = SparkSession.builder \
    .config("spark.jars.packages",
            "org.apache.spark:spark-sql-kafka-0-10_2.12:3.5.0") \
    .getOrCreate()

df = spark.read \
    .format("kafka") \
    .option("kafka.bootstrap.servers", "broker1:9092,broker2:9092") \
    .option("subscribe",               "orders") \
    .option("startingOffsets",         "earliest") \
    .option("endingOffsets",           "latest") \
    .load()

# Kafka schema: key, value, topic, partition, offset, timestamp, timestampType
from pyspark.sql.types import StructType, StructField, StringType, DoubleType
from pyspark.sql import functions as F

order_schema = StructType([
    StructField("order_id",   StringType()),
    StructField("user_id",    StringType()),
    StructField("total",      DoubleType()),
    StructField("created_at", StringType()),
])

orders = df.select(
    F.from_json(F.col("value").cast("string"), order_schema).alias("data")
).select("data.*")

orders.show()
```

### Structured Streaming from Kafka

```python
stream_df = spark.readStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers", "broker:9092") \
    .option("subscribe",               "events,clicks,impressions") \
    .option("startingOffsets",         "latest") \
    .option("maxOffsetsPerTrigger",    50000) \
    .option("kafka.group.id",          "spark-consumer-group") \
    .load()

event_schema = "user_id STRING, event STRING, ts TIMESTAMP, properties MAP<STRING, STRING>"

parsed = stream_df.select(
    F.from_json(F.col("value").cast("string"), event_schema).alias("e"),
    "topic", "partition", "offset", "timestamp"
).select("e.*", "topic", "offset", "timestamp")

# Windowed aggregation
result = parsed \
    .withWatermark("ts", "5 minutes") \
    .groupBy(
        F.window("ts", "10 minutes", "5 minutes"),
        "event"
    ).count()

query = result.writeStream \
    .format("console") \
    .outputMode("update") \
    .option("checkpointLocation", "/tmp/kafka-checkpoint/") \
    .trigger(processingTime="30 seconds") \
    .start()
```

### Write to Kafka

```python
# Batch write
output_df \
    .select(
        F.col("order_id").cast("string").alias("key"),
        F.to_json(F.struct("*")).alias("value"),
    ) \
    .write \
    .format("kafka") \
    .option("kafka.bootstrap.servers", "broker:9092") \
    .option("topic",                   "processed-orders") \
    .save()

# Streaming write
query = result \
    .select(
        F.col("user_id").cast("string").alias("key"),
        F.to_json(F.struct("*")).alias("value"),
    ) \
    .writeStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers", "broker:9092") \
    .option("topic",                   "aggregated-events") \
    .option("checkpointLocation",      "/tmp/kafka-write-checkpoint/") \
    .start()
```

---

## End-to-End Streaming Pipeline

```
Kafka (raw events)
       ↓
Spark Structured Streaming (parse, enrich, aggregate)
       ↓
Delta Lake (sink — ACID, time travel, batch queries)
       ↓
Airflow DAG (batch jobs trigger on schedule)
       ↓
Data Warehouse (Redshift / BigQuery / Snowflake)
```

```python
# streaming_pipeline.py
from pyspark.sql import SparkSession, functions as F

spark = SparkSession.builder \
    .appName("EventPipeline") \
    .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
    .getOrCreate()

EVENT_SCHEMA = "user_id STRING, event_type STRING, value DOUBLE, ts TIMESTAMP"

raw = spark.readStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers", "broker:9092") \
    .option("subscribe", "raw-events") \
    .load()

parsed = raw.select(
    F.from_json(F.col("value").cast("string"), EVENT_SCHEMA).alias("e")
).select("e.*")

enriched = parsed \
    .withWatermark("ts", "10 minutes") \
    .filter(F.col("value") > 0) \
    .withColumn("date", F.to_date("ts")) \
    .withColumn("hour", F.hour("ts"))

query = enriched.writeStream \
    .format("delta") \
    .outputMode("append") \
    .option("checkpointLocation", "s3://bucket/checkpoints/event-pipeline/") \
    .option("path",                "s3://bucket/delta/events/") \
    .partitionBy("date", "hour") \
    .trigger(processingTime="1 minute") \
    .start()

query.awaitTermination()
```
