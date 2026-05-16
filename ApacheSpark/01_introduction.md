# Chapter 1: Introduction to Apache Spark

## What is Apache Spark?

**Apache Spark** is an open-source, distributed computing framework designed for fast, large-scale data processing. Originally developed at UC Berkeley's AMPLab in 2009 and donated to the Apache Software Foundation in 2013, it has become the industry standard for big data engineering and data science.

> **Analogy:** Think of Spark as a **massive parallel assembly line** for data:
> - You split data across hundreds of machines.
> - Each machine processes its chunk simultaneously.
> - Results are combined and returned in seconds — not hours.

---

## Why Spark Exists

Before Spark, the dominant big data framework was **Hadoop MapReduce**. MapReduce had fundamental limitations:

| Problem with MapReduce | Spark Solution |
|---|---|
| Writes intermediate results to disk after every step | Keeps data in-memory across steps |
| Slow for iterative algorithms (ML, graph) | In-memory caching — 10–100× faster |
| Only supports batch processing | Supports batch + streaming + ML + SQL |
| Java-only, verbose API | Python, Scala, Java, R — concise API |
| No interactive querying | Interactive notebooks (Jupyter, Databricks) |

---

## The Spark Ecosystem

```
┌──────────────────────────────────────────────────────┐
│                  SPARK APPLICATIONS                   │
│                                                       │
│  Spark SQL    Spark Streaming   MLlib      GraphX    │
│  (SQL/DF)    (real-time data)  (ML/AI)   (graphs)   │
│                                                       │
│             Spark Core (RDD, Scheduling, Memory)      │
└──────────────────────────────────────────────────────┘
         │               │              │
      Standalone       YARN           Kubernetes
      Cluster         (Hadoop)
```

| Component | Purpose |
|---|---|
| **Spark Core** | Task scheduling, memory management, fault recovery |
| **Spark SQL** | SQL queries, DataFrames, Hive integration |
| **Structured Streaming** | Real-time stream processing on DataFrames |
| **MLlib** | Scalable machine learning pipelines |
| **GraphX** | Graph computation (PageRank, connected components) |

---

## How Spark Executes Code

Spark uses **lazy evaluation** — your code builds a plan, not actual execution:

```
1. You write transformations (map, filter, groupBy)
   → Spark records them as a DAG (Directed Acyclic Graph)
   → Nothing runs yet

2. You call an action (collect, count, write)
   → Spark optimizes the DAG (Catalyst optimizer)
   → Splits into Stages separated by shuffle boundaries
   → Sends Tasks to Executors in parallel
   → Returns result to Driver
```

```python
# This builds a DAG — no computation yet
df = spark.read.csv("data.csv")
filtered = df.filter(df.age > 18)
grouped  = filtered.groupBy("city").count()

# THIS triggers execution
grouped.show()   ← action
```

---

## Key Concepts

### RDD (Resilient Distributed Dataset)
The original Spark abstraction. An immutable, distributed collection of objects partitioned across the cluster. Fault-tolerant via lineage — if a partition is lost, Spark recomputes it from the source.

### DataFrame
A distributed table with a schema — like a database table or a Pandas DataFrame, but distributed across the cluster. Built on top of RDDs but with the **Catalyst query optimizer** for dramatic speed improvements.

### Dataset (Scala/Java only)
A typed DataFrame with compile-time type safety. Not available in Python (PySpark uses DataFrames for everything).

### SparkSession
The entry point for all Spark functionality. Replaced `SparkContext`, `SQLContext`, and `HiveContext` in Spark 2.0+.

```python
from pyspark.sql import SparkSession

spark = SparkSession.builder \
    .appName("MyApp") \
    .config("spark.executor.memory", "4g") \
    .getOrCreate()
```

### DAG (Directed Acyclic Graph)
Spark's representation of your computation. Each node is an RDD/DataFrame operation; edges are dependencies. Spark optimizes and executes the DAG when an action is called.

### Partition
The smallest unit of parallelism. Each partition is processed by one task on one executor. More partitions = more parallelism (up to the number of CPU cores available).

---

## Spark vs. Alternatives

| Feature | Spark | Hadoop MapReduce | Flink | Dask |
|---|---|---|---|---|
| Speed | Very fast (in-memory) | Slow (disk I/O) | Very fast | Fast (single machine scale) |
| Streaming | Micro-batch + continuous | No | True streaming | Limited |
| ML support | MLlib | Mahout (limited) | FlinkML | scikit-learn integration |
| SQL support | Spark SQL | Hive (slow) | Flink SQL | Limited |
| Scale | 1 to 10,000+ nodes | 1 to 10,000+ nodes | 1 to 1,000+ nodes | Single machine (mostly) |
| Python API | PySpark (excellent) | Streaming (poor) | PyFlink | Native Python |
| Ease of use | Medium | Complex | Medium | Easy |
| Best for | Batch + streaming + ML | Legacy batch | Low-latency streaming | Data science / medium data |

---

## Real-World Use Cases

### ETL Pipelines
Read from S3/HDFS → clean and transform → write to data warehouse (Redshift, BigQuery, Snowflake). Spark replaces complex chains of SQL scripts with a unified, versioned Python codebase.

### Log Processing
Process billions of log lines per day. Parse, aggregate, and enrich in parallel across hundreds of nodes. Write results to Elasticsearch or a data lake.

### Machine Learning at Scale
Train ML models on datasets too large to fit on one machine. Distributed feature engineering with Spark SQL + MLlib pipelines for model training and hyperparameter tuning.

### Real-Time Streaming
Spark Structured Streaming ingests from Kafka, processes events with windowed aggregations, and writes results to a database — with exactly-once semantics.

### Data Lake Analytics
Query petabyte-scale Parquet/Delta Lake tables stored on S3 using Spark SQL. Power BI, Tableau, or Jupyter connects to Spark as the compute engine.

---

## Spark Versions

| Version | Year | Key Changes |
|---|---|---|
| 1.x | 2014 | DataFrames introduced |
| 2.0 | 2016 | SparkSession, Structured Streaming, Dataset API |
| 2.4 | 2018 | Barrier execution mode, Python 3 default |
| 3.0 | 2020 | Adaptive Query Execution (AQE), Pandas UDFs v2 |
| 3.1 | 2021 | K8s improvements, PySpark type hints |
| 3.2 | 2021 | pandas API on Spark, ANSI SQL |
| 3.3 | 2022 | Python 3.x default, better error messages |
| 3.4 | 2023 | Connect (remote mode), Spark Connect protocol |
| 3.5 | 2023 | Spark Connect GA, improved Python UDFs |
