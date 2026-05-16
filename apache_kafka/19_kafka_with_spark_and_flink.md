# Chapter 19: Kafka with Spark & Flink

## Kafka + Apache Spark Structured Streaming

### Read from Kafka

```python
from pyspark.sql import SparkSession
from pyspark.sql import functions as F
from pyspark.sql.types import StructType, StructField, StringType, DoubleType, TimestampType

spark = SparkSession.builder \
    .appName("KafkaStreamProcessor") \
    .config("spark.jars.packages",
            "org.apache.spark:spark-sql-kafka-0-10_2.12:3.5.0") \
    .getOrCreate()

# Read stream from Kafka
raw = spark.readStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers",  "broker:9092") \
    .option("subscribe",                "orders,events") \
    .option("startingOffsets",          "latest") \
    .option("maxOffsetsPerTrigger",     50000) \
    .option("kafka.group.id",           "spark-consumer") \
    .option("kafka.security.protocol",  "SASL_SSL") \
    .option("kafka.sasl.mechanism",     "PLAIN") \
    .load()

# Kafka schema: key, value, topic, partition, offset, timestamp, timestampType
# Both key and value are binary — cast to string/parse
```

### Parse JSON Messages

```python
order_schema = StructType([
    StructField("order_id",   StringType()),
    StructField("user_id",    StringType()),
    StructField("amount",     DoubleType()),
    StructField("currency",   StringType()),
    StructField("created_at", TimestampType()),
])

orders = raw \
    .filter(F.col("topic") == "orders") \
    .select(
        F.from_json(F.col("value").cast("string"), order_schema).alias("data"),
        "topic", "partition", "offset", "timestamp"
    ) \
    .select("data.*", "topic", "partition", "offset", "timestamp")
```

### Windowed Aggregation

```python
# Revenue per user per 10-minute window
revenue = orders \
    .withWatermark("created_at", "5 minutes") \
    .groupBy(
        F.window("created_at", "10 minutes", "5 minutes"),
        "user_id"
    ) \
    .agg(
        F.sum("amount").alias("total_revenue"),
        F.count("*").alias("order_count")
    )

query = revenue.writeStream \
    .format("console") \
    .outputMode("update") \
    .option("checkpointLocation", "/tmp/checkpoints/revenue/") \
    .trigger(processingTime="30 seconds") \
    .start()

query.awaitTermination()
```

### Write Back to Kafka

```python
# Transform and write to another Kafka topic
enriched = orders \
    .withColumn("tier", F.when(F.col("amount") > 1000, "premium").otherwise("standard")) \
    .withColumn("processed_at", F.current_timestamp())

query = enriched \
    .select(
        F.col("order_id").cast("string").alias("key"),
        F.to_json(F.struct("*")).alias("value"),
    ) \
    .writeStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers",   "broker:9092") \
    .option("topic",                     "enriched-orders") \
    .option("checkpointLocation",        "/tmp/checkpoints/enriched/") \
    .start()
```

### foreachBatch for Complex Sinks

```python
def process_batch(batch_df, batch_id):
    # Write to multiple destinations
    batch_df.write \
        .format("delta") \
        .mode("append") \
        .option("path", "s3://my-bucket/delta/orders/") \
        .save()

    # Write aggregated metrics to PostgreSQL
    batch_df.groupBy("user_id") \
        .agg(F.sum("amount").alias("total")) \
        .write \
        .format("jdbc") \
        .option("url", "jdbc:postgresql://db:5432/analytics") \
        .option("dbtable", "user_metrics") \
        .mode("append") \
        .save()

query = orders.writeStream \
    .foreachBatch(process_batch) \
    .option("checkpointLocation", "/tmp/checkpoints/batch/") \
    .trigger(processingTime="1 minute") \
    .start()
```

---

## Kafka + Apache Flink

### Maven Dependency

```xml
<dependency>
    <groupId>org.apache.flink</groupId>
    <artifactId>flink-connector-kafka</artifactId>
    <version>3.2.0-1.19</version>
</dependency>
```

### Flink Kafka Source

```java
import org.apache.flink.api.common.serialization.SimpleStringSchema;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;

KafkaSource<String> source = KafkaSource.<String>builder()
    .setBootstrapServers("localhost:9092")
    .setTopics("orders")
    .setGroupId("flink-consumer")
    .setStartingOffsets(OffsetsInitializer.earliest())
    .setValueOnlyDeserializer(new SimpleStringSchema())
    .build();

StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

DataStream<String> stream = env.fromSource(
    source,
    WatermarkStrategy.forMonotonousTimestamps(),
    "Kafka Source"
);
```

### Flink Kafka Sink

```java
import org.apache.flink.connector.kafka.sink.KafkaRecordSerializationSchema;
import org.apache.flink.connector.kafka.sink.KafkaSink;

KafkaSink<String> sink = KafkaSink.<String>builder()
    .setBootstrapServers("localhost:9092")
    .setRecordSerializer(KafkaRecordSerializationSchema.builder()
        .setTopic("processed-orders")
        .setValueSerializationSchema(new SimpleStringSchema())
        .build())
    .setDeliveryGuarantee(DeliveryGuarantee.EXACTLY_ONCE)
    .setTransactionalIdPrefix("flink-tx-")
    .setKafkaProducerConfig(producerProps)
    .build();

stream.sinkTo(sink);
```

### Flink Windowing

```java
import org.apache.flink.api.common.eventtime.*;
import org.apache.flink.streaming.api.windowing.assigners.TumblingEventTimeWindows;
import java.time.Duration;

DataStream<Order> orders = rawStream
    .map(json -> objectMapper.readValue(json, Order.class))
    .assignTimestampsAndWatermarks(
        WatermarkStrategy
            .<Order>forBoundedOutOfOrderness(Duration.ofMinutes(1))
            .withTimestampAssigner((event, ts) -> event.getCreatedAt().toEpochMilli())
    );

// Tumbling window: revenue per user per 5 minutes
DataStream<UserRevenue> result = orders
    .keyBy(Order::getUserId)
    .window(TumblingEventTimeWindows.of(Duration.ofMinutes(5)))
    .aggregate(new RevenueAggregator());
```

### Flink Exactly-Once with Kafka

```java
// Enable checkpointing
env.enableCheckpointing(60000, CheckpointingMode.EXACTLY_ONCE);
env.getCheckpointConfig().setCheckpointStorage("s3://my-bucket/flink-checkpoints/");
env.getCheckpointConfig().setMinPauseBetweenCheckpoints(30000);
env.getCheckpointConfig().setCheckpointTimeout(120000);

// Kafka source with exactly-once
KafkaSource<Order> source = KafkaSource.<Order>builder()
    .setBootstrapServers("localhost:9092")
    .setTopics("orders")
    .setGroupId("flink-eos")
    .setStartingOffsets(OffsetsInitializer.committedOffsets())
    .setDeserializer(KafkaRecordDeserializationSchema.valueOnly(new OrderDeserializationSchema()))
    .build();

// Kafka sink with exactly-once transactions
KafkaSink<String> sink = KafkaSink.<String>builder()
    .setBootstrapServers("localhost:9092")
    .setDeliveryGuarantee(DeliveryGuarantee.EXACTLY_ONCE)
    .setTransactionalIdPrefix("flink-eos-")
    // ...
    .build();
```

---

## Spark vs Flink for Kafka Consumers

| Feature | Spark Structured Streaming | Apache Flink |
|---|---|---|
| Processing model | Micro-batch (default) / continuous | True streaming |
| Latency | 100ms – seconds | < 10ms |
| State management | In-memory / HDFS | RocksDB / heap |
| Windowing | Tumbling, sliding, session | Rich (all types) |
| Exactly-once | Yes (foreachBatch) | Yes (native) |
| Language | Python, Java, Scala, R | Java, Scala, Python, SQL |
| SQL support | Spark SQL | Flink SQL (powerful) |
| Learning curve | Easier | Steeper |
| Ecosystem | Spark ecosystem | Flink ecosystem |
| Cloud managed | Databricks, EMR | Confluent, Kinesis Data Analytics |

---

## End-to-End Pipeline

```
Kafka (raw-events)
       │
       ├──► Flink (real-time, < 1s latency)
       │      └──► Kafka (alerts) ──► PagerDuty
       │
       └──► Spark Structured Streaming (batch, 1-min windows)
              ├──► Delta Lake (for BI queries)
              └──► Kafka (enriched-events)
                       └──► Elasticsearch (Kibana dashboards)
```
