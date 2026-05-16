# Chapter 9: Spark Structured Streaming

## What is Structured Streaming?

**Structured Streaming** is Spark's streaming engine built on top of the DataFrame/SQL API. It treats a live data stream as an **unbounded table** that grows continuously. Your batch DataFrame code works on streaming data with minimal changes.

```
Kafka / Files / Socket
        ↓
Streaming DataFrame (unbounded table)
        ↓
Transformations (filter, groupBy, join)
        ↓
Output Sink (Kafka, files, console, database)
```

---

## Micro-Batch vs Continuous Processing

| Mode | Latency | Fault Tolerance | Use Case |
|---|---|---|---|
| **Micro-batch** (default) | ~100ms–seconds | Exactly-once | Most streaming use cases |
| **Continuous** | ~1ms | At-least-once | Ultra-low latency |

---

## Your First Streaming Job

```python
from pyspark.sql import SparkSession
from pyspark.sql import functions as F

spark = SparkSession.builder.appName("StreamingApp").getOrCreate()

# Read a stream (from files appearing in a directory)
stream_df = spark.readStream \
    .format("json") \
    .schema("user_id STRING, event STRING, ts TIMESTAMP") \
    .option("path", "data/stream_input/") \
    .load()

# Apply transformations (same as batch!)
result = stream_df \
    .filter(F.col("event") != "heartbeat") \
    .groupBy("event") \
    .count()

# Start the streaming query and write output
query = result.writeStream \
    .format("console") \
    .outputMode("complete") \
    .option("checkpointLocation", "/tmp/checkpoint/") \
    .trigger(processingTime="10 seconds") \
    .start()

query.awaitTermination()
```

---

## Input Sources

### Kafka

```python
# Requires: spark.jars.packages = org.apache.spark:spark-sql-kafka-0-10_2.12:3.5.0

df = spark.readStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers", "broker:9092") \
    .option("subscribe",               "events,clicks") \
    .option("startingOffsets",         "earliest")    \
    .option("maxOffsetsPerTrigger",    100000) \
    .load()

# Schema: kafka always gives (key, value, topic, partition, offset, timestamp)
# value is binary — you need to cast and parse it
parsed = df.select(
    F.col("key").cast("string").alias("key"),
    F.from_json(
        F.col("value").cast("string"),
        schema="user_id STRING, event STRING, ts TIMESTAMP"
    ).alias("data"),
    "topic", "partition", "offset", "timestamp"
).select("key", "data.*", "topic", "offset", "timestamp")
```

### Files (Auto-Ingest New Files)

```python
df = spark.readStream \
    .format("parquet") \
    .schema(my_schema) \
    .option("path", "s3://bucket/landing/") \
    .option("maxFilesPerTrigger", 100) \
    .load()
```

### Rate Source (Testing)

```python
# Generates N rows per second for testing
df = spark.readStream.format("rate") \
    .option("rowsPerSecond", 1000) \
    .load()
# Columns: timestamp, value (0,1,2,...)
```

### Socket (Debugging Only)

```python
df = spark.readStream.format("socket") \
    .option("host", "localhost").option("port", 9999).load()
```

---

## Output Modes

| Mode | Meaning | Requires |
|---|---|---|
| `append` | Only new rows added since last trigger | No aggregation (or with watermark) |
| `complete` | Full result table every trigger | Aggregation |
| `update` | Only rows that changed since last trigger | Aggregation or stateless |

```python
query = result.writeStream.outputMode("append").format("parquet").start()
```

---

## Output Sinks

### Console (Debug)

```python
query = df.writeStream.format("console").outputMode("append").start()
```

### Files (Parquet / Delta)

```python
query = df.writeStream \
    .format("parquet") \
    .option("path",               "s3://bucket/stream_output/") \
    .option("checkpointLocation", "s3://bucket/checkpoints/stream1/") \
    .outputMode("append") \
    .partitionBy("date") \
    .start()
```

### Kafka

```python
result.selectExpr("CAST(user_id AS STRING) AS key",
                  "to_json(struct(*)) AS value") \
      .writeStream \
      .format("kafka") \
      .option("kafka.bootstrap.servers", "broker:9092") \
      .option("topic",                   "output-events") \
      .option("checkpointLocation",      "/tmp/kafka-checkpoint/") \
      .start()
```

### foreach / foreachBatch

```python
# foreachBatch — most flexible: treat each micro-batch as a DataFrame
def process_batch(batch_df, batch_id):
    batch_df.write \
        .format("jdbc") \
        .option("url", "jdbc:postgresql://...") \
        .option("dbtable", "events") \
        .mode("append") \
        .save()
    print(f"Processed batch {batch_id}: {batch_df.count()} rows")

query = df.writeStream \
    .foreachBatch(process_batch) \
    .option("checkpointLocation", "/tmp/checkpoint/") \
    .start()
```

---

## Triggers

```python
# Micro-batch every 10 seconds
.trigger(processingTime="10 seconds")

# Micro-batch every 1 minute
.trigger(processingTime="1 minute")

# One-time: process all available data and stop
.trigger(once=True)

# Available-now: process all data and stop (replacement for once=True in 3.3+)
.trigger(availableNow=True)

# Continuous (low latency — at-least-once)
.trigger(continuous="1 second")

# Default: as fast as possible (0 delay)
.trigger(processingTime="0 seconds")
```

---

## Watermarks (Handling Late Data)

A **watermark** tells Spark how long to wait for late-arriving data before closing a time window.

```python
df_with_watermark = parsed \
    .withWatermark("event_time", "10 minutes")   ← tolerate up to 10min late data

result = df_with_watermark \
    .groupBy(
        F.window("event_time", "5 minutes", "1 minute"),  # 5-min sliding window
        "event_type"
    ) \
    .count()
```

### Tumbling Window (Non-overlapping)

```python
result = df \
    .withWatermark("ts", "5 minutes") \
    .groupBy(F.window("ts", "1 hour"), "country") \
    .agg(F.count("*").alias("events"), F.sum("revenue").alias("revenue"))
```

### Sliding Window (Overlapping)

```python
# 1-hour window, slide every 15 minutes
result = df \
    .withWatermark("ts", "10 minutes") \
    .groupBy(F.window("ts", "1 hour", "15 minutes")) \
    .count()
```

---

## Stateful Operations

### Stream-Static Join

```python
# Join stream with a static reference table
products = spark.read.parquet("s3://bucket/products/")  # static

enriched = stream_df.join(
    broadcast(products),
    on="product_id",
    how="left"
)
```

### Stream-Stream Join (with Watermark)

```python
impressions = spark.readStream.format("kafka")...  # ad impressions
clicks       = spark.readStream.format("kafka")...  # ad clicks

impressions_wm = impressions.withWatermark("impression_time", "2 hours")
clicks_wm      = clicks.withWatermark("click_time",       "3 hours")

result = impressions_wm.join(
    clicks_wm,
    expr("""
        impression_id = click_id AND
        click_time BETWEEN impression_time AND impression_time + INTERVAL 1 HOUR
    """),
    how="leftOuter"
)
```

---

## Checkpointing

Checkpointing saves streaming state to durable storage — enables fault recovery:

```python
query = df.writeStream \
    .option("checkpointLocation", "s3://bucket/checkpoints/my-stream/") \
    .format("delta") \
    .start()
```

- Saves: offsets consumed, aggregation state, metadata
- On restart, Spark resumes from exactly where it left off

---

## Managing Streaming Queries

```python
# Start query
query = df.writeStream.format("console").start()

# Status
print(query.status)
print(query.lastProgress)
print(query.recentProgress)

# Stop
query.stop()

# Wait for termination (blocking)
query.awaitTermination()
query.awaitTermination(timeout=300)  # seconds

# List all active queries
for q in spark.streams.active:
    print(q.name, q.status)
```
