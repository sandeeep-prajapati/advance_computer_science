# Chapter 20: Learning Roadmap

## 5-Week Hands-On Plan

This roadmap takes you from zero to production-ready Apache Spark skills with a real project each week.

---

## Week 1: Foundation

**Goal:** Run PySpark locally, understand the execution model, work with DataFrames.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Install PySpark (`pip install pyspark`), run `spark.range(100).show()` |
| 2 | Create DataFrames from CSV/JSON, inspect schema, `show()`, `printSchema()` |
| 3 | Practice: `filter`, `select`, `withColumn`, `groupBy`, `agg` |
| 4 | Load a real dataset (NYC Taxi, Airlines, Kaggle), run 10 transformations |
| 5 | Use `explain()`, observe DAG in Spark UI (port 4040) |

### Practice Project: NYC Taxi Analytics

```python
# Download: https://www.nyc.gov/site/tlc/about/tlc-trip-record-data.page

df = spark.read.parquet("yellow_tripdata_2024-01.parquet")
df.printSchema()
df.describe("trip_distance", "fare_amount", "tip_amount").show()

# Top 10 pickup locations
df.groupBy("PULocationID") \
  .agg(F.count("*").alias("trips"), F.avg("fare_amount").alias("avg_fare")) \
  .orderBy("trips", ascending=False) \
  .show(10)

# Peak hours
df.withColumn("hour", F.hour("tpep_pickup_datetime")) \
  .groupBy("hour").count() \
  .orderBy("hour").show(24)
```

### Week 1 Checklist
- [ ] SparkSession created locally
- [ ] Read CSV and Parquet files
- [ ] Applied 5+ transformations (filter, select, groupBy, withColumn, join)
- [ ] Called actions: show, count, collect, write
- [ ] Examined DAG in Spark UI
- [ ] Understood lazy evaluation with a 3-step example

---

## Week 2: Deep Dive — SQL, Windows, Optimization

**Goal:** Master Spark SQL, window functions, joins, and basic performance tuning.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Register DataFrames as views, run 10 SQL queries |
| 2 | Implement all window functions: rank, dense_rank, lag, lead, running_total |
| 3 | Practice all join types; test broadcast join vs sort-merge join |
| 4 | Read and write Parquet with partitioning, observe partition pruning |
| 5 | Profile a slow query: use `explain()`, tune shuffle partitions, enable AQE |

### Practice: E-Commerce Analytics

```python
# Simulate e-commerce data
orders   = spark.read.parquet("data/orders/")
products = spark.read.parquet("data/products/")
users    = spark.read.parquet("data/users/")

# Revenue by month and category
monthly = orders \
    .join(broadcast(products), "product_id") \
    .groupBy(F.date_trunc("month", "order_date").alias("month"), "category") \
    .agg(F.sum("revenue").alias("revenue")) \
    .orderBy("month", "category")

# Rank customers by lifetime value
from pyspark.sql import Window
w = Window.orderBy(F.col("ltv").desc())
top_customers = orders \
    .groupBy("user_id").agg(F.sum("revenue").alias("ltv")) \
    .join(users, "user_id") \
    .withColumn("rank", F.rank().over(w))
```

### Week 2 Checklist
- [ ] 10 SQL queries using CTEs and window functions
- [ ] All join types implemented and tested
- [ ] Broadcast join forced with `broadcast()` hint
- [ ] Partitioned Parquet write and read with partition pruning
- [ ] `explain()` read and understood (look for shuffle, filter pushdown)
- [ ] AQE enabled and shuffle partitions tuned

---

## Week 3: Streaming + Delta Lake

**Goal:** Build a streaming pipeline, store in Delta, run batch and streaming queries.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Set up Kafka with Docker, produce 1,000 test messages |
| 2 | Read Kafka stream, parse JSON, write to console |
| 3 | Add watermark and 5-minute window aggregation |
| 4 | Write stream to Delta Lake, read back with time travel |
| 5 | Implement MERGE (upsert) from a batch job |

### Practice: Event Stream Pipeline

```python
# Docker Kafka setup
# docker run -p 9092:9092 confluentinc/cp-kafka:latest

# Producer script (Python)
from kafka import KafkaProducer
import json, time, random

producer = KafkaProducer(bootstrap_servers="localhost:9092",
                         value_serializer=lambda v: json.dumps(v).encode())

while True:
    producer.send("events", {
        "user_id":     f"user_{random.randint(1, 1000)}",
        "event_type":  random.choice(["page_view", "click", "purchase"]),
        "value":       round(random.uniform(0, 1000), 2),
        "ts":          time.strftime("%Y-%m-%dT%H:%M:%SZ"),
    })
    time.sleep(0.1)
```

### Week 3 Checklist
- [ ] Kafka running in Docker
- [ ] Structured Streaming job reading from Kafka
- [ ] Watermark and window aggregation working
- [ ] Stream writing to Delta Lake
- [ ] Time travel query on Delta table
- [ ] MERGE implemented for upsert pattern

---

## Week 4: Cluster & Production

**Goal:** Deploy on a real cluster (Docker Compose or Kubernetes), tune performance.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Launch 3-node Spark cluster with Docker Compose |
| 2 | Submit jobs with spark-submit, observe Spark UI on master |
| 3 | Tune executor memory, cores, shuffle partitions |
| 4 | Trigger and fix a data skew problem (use AQE or salt) |
| 5 | Set up History Server, review completed job event logs |

### Week 4 Checklist
- [ ] 3-node standalone cluster running
- [ ] Job submitted via spark-submit (not local mode)
- [ ] Executor sizing tuned (cores × memory)
- [ ] Skew identified in Spark UI and mitigated
- [ ] Dynamic allocation enabled
- [ ] History Server showing event logs

---

## Week 5: Capstone Project

**Goal:** Build an end-to-end production data pipeline.

### Capstone Option A: Streaming Analytics Platform

```
Kafka (user events) → Spark Streaming → Delta Lake (events) → Airflow (batch aggregation) → Parquet/DB (dashboard)
```

Features:
- Real-time event ingestion from Kafka
- 5-min, 1-hour, and daily windowed aggregations
- Delta Lake as unified storage (streaming writes, batch reads)
- Airflow DAG for daily batch jobs
- Spark UI monitoring

### Capstone Option B: Large-Scale ETL Pipeline

```
S3 (raw data, 100GB) → Spark ETL → Delta Lake → Spark SQL analytics → Output reports
```

Features:
- Read 100GB+ of CSV/JSON
- Clean: null handling, deduplication, type normalization
- Transform: join 5+ tables, compute derived metrics
- Partition output by date/country
- Benchmark before/after optimization (AQE, broadcast, caching)

### Capstone Option C: ML Pipeline at Scale

```
Raw data (1M rows) → Spark feature engineering → MLlib training → Model evaluation → Predictions
```

Features:
- Feature pipeline with StringIndexer, VectorAssembler, StandardScaler
- Train RandomForest + GBT, compare AUC
- Cross-validation with ParamGridBuilder
- Save model, load and run batch predictions
- Use Pandas UDF for real-time inference on stream

---

## Key Commands Reference Card

```bash
# Start PySpark shell
pyspark --master local[4]

# Submit a job
spark-submit --master local[4] my_job.py
spark-submit --master yarn --deploy-mode cluster my_job.py
spark-submit --master spark://master:7077 my_job.py
spark-submit --master k8s://https://k8s:6443 --deploy-mode cluster my_job.py

# Key Spark UI ports
http://localhost:4040    # Active job UI
http://localhost:18080   # History Server
http://localhost:8080    # Standalone master

# YARN commands
yarn application -list
yarn logs -applicationId application_XXX
yarn application -kill application_XXX
```

```python
# Essential DataFrame operations
df.show(5, truncate=False)
df.printSchema()
df.explain("formatted")
df.count()
df.rdd.getNumPartitions()
df.groupBy(F.spark_partition_id()).count().show()

# Performance
df.cache()
df.unpersist()
spark.catalog.clearCache()
spark.conf.set("spark.sql.shuffle.partitions", "200")
spark.conf.set("spark.sql.adaptive.enabled",   "true")
```

---

## Beyond the Roadmap

| Topic | Notes |
|---|---|
| **Apache Iceberg** | Delta Lake alternative — multi-engine support |
| **Apache Hudi** | Upserts on data lakes, CDC patterns |
| **Spark Connect** | Thin client mode — Spark 3.4+ |
| **Photon Engine** | Databricks' vectorized C++ engine |
| **dbt + Spark** | SQL transformations on Spark via dbt-spark |
| **Great Expectations** | Data quality testing integrated with Spark |
| **MLflow** | Track ML experiments, deploy models from Spark |
| **Ray on Spark** | Python ML/RL workloads on Spark clusters |
