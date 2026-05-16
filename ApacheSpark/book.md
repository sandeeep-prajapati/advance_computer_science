# Apache Spark: The Complete Developer Book
> A comprehensive, single-file reference covering everything from your first RDD to production-grade clusters, Delta Lake, MLlib, streaming, and cloud deployments.

---

## Quick Navigation

- [Table of Contents](#table-of-contents)
- [Chapter 1: Introduction to Apache Spark](#chapter-1-introduction-to-apache-spark)
- [Chapter 2: Setup & Installation](#chapter-2-setup--installation)
- [Chapter 3: Core Concepts](#chapter-3-core-concepts)
- [Chapter 4: RDDs — Resilient Distributed Datasets](#chapter-4-rdds--resilient-distributed-datasets)
- [Chapter 5: DataFrames & Datasets](#chapter-5-dataframes--datasets)
- [Chapter 6: Spark SQL](#chapter-6-spark-sql)
- [Chapter 7: Reading & Writing Data](#chapter-7-reading--writing-data)
- [Chapter 8: Transformations & Actions](#chapter-8-transformations--actions)
- [Chapter 9: Spark Structured Streaming](#chapter-9-spark-structured-streaming)
- [Chapter 10: Cluster Architecture](#chapter-10-cluster-architecture)
- [Chapter 11: Performance & Optimization](#chapter-11-performance--optimization)
- [Chapter 12: Spark MLlib](#chapter-12-spark-mllib)
- [Chapter 13: Delta Lake](#chapter-13-delta-lake)
- [Chapter 14: Spark on Kubernetes](#chapter-14-spark-on-kubernetes)
- [Chapter 15: Spark on YARN & EMR](#chapter-15-spark-on-yarn--emr)
- [Chapter 16: PySpark Advanced](#chapter-16-pyspark-advanced)
- [Chapter 17: Monitoring & Debugging](#chapter-17-monitoring--debugging)
- [Chapter 18: Security](#chapter-18-security)
- [Chapter 19: Spark with Airflow & Kafka](#chapter-19-spark-with-airflow--kafka)
- [Chapter 20: Learning Roadmap](#chapter-20-learning-roadmap)

---

# Apache Spark Developer Book

A comprehensive, chapter-wise guide to mastering Apache Spark — from your first RDD to production-grade pipelines, streaming, ML, and cluster tuning.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_introduction.md) | **Introduction to Apache Spark** | What is Spark, why it exists, Spark vs Hadoop MapReduce, components |
| [02](02_setup_and_installation.md) | **Setup & Installation** | Local mode, PySpark, Spark Shell, Docker, Databricks, spark-submit |
| [03](03_core_concepts.md) | **Core Concepts** | RDD, DataFrame, Dataset, DAG, lazy evaluation, partitions, transformations vs actions |
| [04](04_rdds.md) | **RDDs — Resilient Distributed Datasets** | Creating, transforming, and actioning RDDs, persistence, lineage |
| [05](05_dataframes_and_datasets.md) | **DataFrames & Datasets** | Schema, SparkSession, select, filter, groupBy, join, UDFs |
| [06](06_spark_sql.md) | **Spark SQL** | SQL queries, views, catalog, JDBC, window functions, explain plan |
| [07](07_reading_and_writing_data.md) | **Reading & Writing Data** | CSV, JSON, Parquet, ORC, Delta Lake, JDBC, Kafka, partitioned writes |
| [08](08_transformations_and_actions.md) | **Transformations & Actions** | map, filter, flatMap, groupBy, agg, join types, repartition, coalesce |
| [09](09_spark_streaming.md) | **Spark Structured Streaming** | Micro-batch vs continuous, triggers, watermarks, stateful ops, sinks |
| [10](10_cluster_architecture.md) | **Cluster Architecture** | Driver, Executors, Cluster Managers (YARN, K8s, Standalone), DAGScheduler |
| [11](11_performance_and_optimization.md) | **Performance & Optimization** | Shuffle, broadcast joins, caching, partitioning, AQE, skew handling |
| [12](12_spark_mllib.md) | **Spark MLlib** | Pipelines, feature engineering, classification, regression, clustering, tuning |
| [13](13_delta_lake.md) | **Delta Lake** | ACID transactions, time travel, schema enforcement, MERGE, OPTIMIZE |
| [14](14_spark_on_kubernetes.md) | **Spark on Kubernetes** | spark-submit on K8s, operator, dynamic allocation, resource management |
| [15](15_spark_on_yarn.md) | **Spark on YARN / EMR** | Client vs cluster mode, queue management, AWS EMR setup |
| [16](16_pyspark_advanced.md) | **PySpark Advanced** | Pandas UDFs (Arrow), pandas-on-Spark, type hints, koalas |
| [17](17_monitoring_and_debugging.md) | **Monitoring & Debugging** | Spark UI, event logs, Ganglia, metrics, GC tuning, task failures |
| [18](18_security.md) | **Security** | Kerberos, SSL, encryption at rest, ACLs, Spark on secure clusters |
| [19](19_spark_with_airflow_and_kafka.md) | **Spark with Airflow & Kafka** | SparkSubmitOperator, Kafka source/sink, end-to-end streaming pipeline |
| [20](20_learning_roadmap.md) | **Learning Roadmap** | 5-week hands-on plan from zero to production-ready |

---

## How to Read This Book

**If you're brand new to Spark:** Read 1 → 2 → 3 → 4 → 5 in order.

**If you're a Python/data engineer:** Read 2 → 5 → 6 → 7 → 8 → 16.

**If you're building streaming pipelines:** Read 9 → 7 → 19.

**If you're going to production:** Read 10 → 11 → 14 → 17.

**If you're doing ML:** Read 12 → 13 → 11.

**If something is slow:** Jump to Chapter 11 (Performance) and Chapter 17 (Debugging).

---

## Key Concepts at a Glance

```
Apache Spark Cluster
├── Driver Program
│   ├── SparkContext / SparkSession  (entry point)
│   ├── DAGScheduler                 (builds execution plan)
│   ├── TaskScheduler                (sends tasks to executors)
│   └── SparkUI                      (port 4040 — live monitoring)
│
├── Cluster Manager  (YARN / K8s / Standalone / Mesos)
│   └── Allocates executors on worker nodes
│
├── Executors  (JVM processes on worker nodes)
│   ├── Task Threads                 (run individual tasks)
│   ├── Block Manager                (caches RDDs/DataFrames)
│   └── Shuffle Service              (exchange data between stages)
│
├── Core Abstractions
│   ├── RDD          (low-level, distributed collection, fault-tolerant)
│   ├── DataFrame    (distributed table with schema, SQL-optimized)
│   ├── Dataset      (typed DataFrame — Scala/Java only)
│   └── DStream      (legacy streaming — use Structured Streaming)
│
├── Execution Model
│   ├── Job          (triggered by an action)
│   ├── Stage        (set of tasks separated by shuffle boundaries)
│   ├── Task         (unit of work — one partition on one executor)
│   └── Partition    (chunk of data processed in parallel)
│
└── Key Features
    ├── Lazy Evaluation    (transformations build a DAG, not executed until action)
    ├── In-Memory Caching  (cache() / persist() to avoid recomputation)
    ├── Fault Tolerance    (lineage-based recomputation on failure)
    └── Unified Engine     (batch + streaming + ML + graph in one framework)
```


---

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


---

# Chapter 2: Setup & Installation

## Your Options

| Method | Best For |
|---|---|
| **PySpark (pip)** | Local development, data science notebooks |
| **Docker** | Isolated local cluster, CI pipelines |
| **Databricks Community** | Managed cluster, free tier, notebooks |
| **Google Colab** | Zero-install, notebook-based learning |
| **spark-submit** | Running jobs on a real cluster |
| **AWS EMR / GCP Dataproc** | Managed cloud clusters |

---

## Option A: PySpark Local (Quickest Start)

### Requirements
- Python 3.8+
- Java 8 or 11 (must be on `PATH`)

```bash
# Install Java (Ubuntu/Debian)
sudo apt-get install -y openjdk-11-jdk
export JAVA_HOME=/usr/lib/jvm/java-11-openjdk-amd64

# Install PySpark
pip install pyspark==3.5.0

# Optional: Jupyter for interactive use
pip install jupyter
```

### Verify

```python
# test_spark.py
from pyspark.sql import SparkSession

spark = SparkSession.builder \
    .appName("test") \
    .master("local[*]") \
    .getOrCreate()

spark.range(10).show()
spark.stop()
```

```bash
python test_spark.py
```

---

## Option B: Docker Compose (Local Cluster)

```yaml
# docker-compose.yml
version: '3.8'
services:
  spark-master:
    image: bitnami/spark:3.5
    environment:
      - SPARK_MODE=master
      - SPARK_RPC_AUTHENTICATION_ENABLED=no
    ports:
      - "8080:8080"   # Spark UI
      - "7077:7077"   # Spark master port
    volumes:
      - ./jobs:/opt/spark/jobs

  spark-worker-1:
    image: bitnami/spark:3.5
    environment:
      - SPARK_MODE=worker
      - SPARK_MASTER_URL=spark://spark-master:7077
      - SPARK_WORKER_MEMORY=2G
      - SPARK_WORKER_CORES=2
    depends_on:
      - spark-master

  spark-worker-2:
    image: bitnami/spark:3.5
    environment:
      - SPARK_MODE=worker
      - SPARK_MASTER_URL=spark://spark-master:7077
      - SPARK_WORKER_MEMORY=2G
      - SPARK_WORKER_CORES=2
    depends_on:
      - spark-master
```

```bash
docker-compose up -d

# Submit a job to the cluster
docker exec spark-master spark-submit \
  --master spark://spark-master:7077 \
  /opt/spark/jobs/my_job.py
```

Access Spark UI at: `http://localhost:8080`

---

## Option C: Databricks Community Edition (Free)

1. Sign up at [community.cloud.databricks.com](https://community.cloud.databricks.com)
2. Create a cluster (Runtime: Databricks Runtime 13.x = Spark 3.4)
3. Create a notebook → choose Python
4. SparkSession is pre-created as `spark`

---

## The Spark Shell

### PySpark Shell

```bash
pyspark --master local[4]

# Inside the shell:
>>> spark.version
'3.5.0'
>>> sc = spark.sparkContext
>>> rdd = sc.parallelize([1, 2, 3, 4, 5])
>>> rdd.sum()
15
```

### Scala Spark Shell

```bash
spark-shell --master local[4]

// Inside the shell:
scala> spark.range(5).show()
```

---

## spark-submit: Running Jobs

`spark-submit` is how you run Spark applications on a cluster.

```bash
# Local mode
spark-submit \
  --master local[4] \
  --driver-memory 2g \
  my_job.py

# Standalone cluster
spark-submit \
  --master spark://master:7077 \
  --deploy-mode cluster \
  --executor-memory 4g \
  --executor-cores 2 \
  --num-executors 10 \
  my_job.py

# YARN
spark-submit \
  --master yarn \
  --deploy-mode cluster \
  --executor-memory 4g \
  --num-executors 20 \
  --conf spark.yarn.queue=default \
  hdfs:///jobs/my_job.py arg1 arg2

# Kubernetes
spark-submit \
  --master k8s://https://k8s-master:6443 \
  --deploy-mode cluster \
  --conf spark.kubernetes.container.image=my-spark:3.5 \
  --conf spark.kubernetes.namespace=spark-jobs \
  local:///opt/spark/jobs/my_job.py
```

---

## SparkSession Configuration

```python
from pyspark.sql import SparkSession

spark = SparkSession.builder \
    .appName("ProductionJob") \
    .master("local[*]") \
    \
    # Memory
    .config("spark.driver.memory",         "4g") \
    .config("spark.executor.memory",       "8g") \
    .config("spark.executor.cores",        "4") \
    .config("spark.memory.fraction",       "0.8") \
    \
    # Performance
    .config("spark.sql.shuffle.partitions","200") \
    .config("spark.default.parallelism",   "200") \
    .config("spark.sql.adaptive.enabled",  "true") \
    \
    # Delta Lake / Hive
    .config("spark.sql.extensions",        "io.delta.sql.DeltaSparkSessionExtension") \
    .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
    \
    # Logging
    .config("spark.eventLog.enabled",      "true") \
    .config("spark.eventLog.dir",          "/tmp/spark-events") \
    \
    .getOrCreate()

# Set log level
spark.sparkContext.setLogLevel("WARN")
```

---

## Project Structure

```
my_spark_project/
├── src/
│   ├── jobs/
│   │   ├── etl_job.py
│   │   └── streaming_job.py
│   └── utils/
│       ├── spark_session.py
│       └── schema.py
├── tests/
│   └── test_etl.py
├── data/
│   └── sample.csv
├── requirements.txt
└── submit.sh
```

### Reusable SparkSession Factory

```python
# src/utils/spark_session.py
from pyspark.sql import SparkSession

def get_spark(app_name: str, **configs) -> SparkSession:
    builder = SparkSession.builder.appName(app_name)
    for key, val in configs.items():
        builder = builder.config(key, val)
    return builder.getOrCreate()
```

---

## Jupyter + PySpark

```bash
pip install pyspark findspark jupyter

# Start Jupyter
jupyter notebook
```

```python
# In a notebook cell
import findspark
findspark.init()

from pyspark.sql import SparkSession
spark = SparkSession.builder.appName("notebook").master("local[*]").getOrCreate()
spark.range(5).show()
```

---

## Useful Environment Variables

```bash
export SPARK_HOME=/usr/local/spark
export PYSPARK_PYTHON=python3
export PYSPARK_DRIVER_PYTHON=jupyter
export PYSPARK_DRIVER_PYTHON_OPTS='notebook'
export JAVA_HOME=/usr/lib/jvm/java-11-openjdk-amd64
export PATH=$SPARK_HOME/bin:$JAVA_HOME/bin:$PATH
```


---

# Chapter 3: Core Concepts

## The Execution Model

Before writing any code, understand how Spark actually runs your program:

```
Your Code (Python/Scala/Java)
        ↓
  SparkSession / SparkContext
        ↓
  Transformations build a DAG (no execution yet)
        ↓
  Action called → Catalyst Optimizer rewrites plan
        ↓
  DAGScheduler splits into Stages
        ↓
  TaskScheduler sends Tasks to Executors
        ↓
  Executors run Tasks in parallel
        ↓
  Results sent back to Driver
```

---

## Lazy Evaluation

This is Spark's most important principle. **Nothing runs until you call an action.**

```python
# Step 1: Read — NO data is loaded yet
df = spark.read.csv("s3://bucket/data/*.csv", header=True)

# Step 2: Filter — adds to the plan, still NO execution
filtered = df.filter(df["age"] > 25)

# Step 3: Group — still just building the DAG
result = filtered.groupBy("country").count()

# Step 4: Action — NOW Spark executes everything at once, optimized
result.show()   ← triggers the entire pipeline
```

**Why it matters:** Spark can optimize the entire plan before running. If `country == "US"` is added as a later filter, Spark's Catalyst optimizer can push it down before the groupBy — skipping irrelevant data early.

---

## Transformations vs. Actions

### Transformations (Lazy — Build the DAG)

| Transformation | Description |
|---|---|
| `filter()` / `where()` | Keep rows matching condition |
| `select()` | Choose columns |
| `withColumn()` | Add or replace a column |
| `drop()` | Remove columns |
| `groupBy()` | Group rows by key |
| `agg()` | Apply aggregate functions |
| `join()` | Join two DataFrames |
| `union()` | Combine two DataFrames vertically |
| `distinct()` | Remove duplicate rows |
| `orderBy()` / `sort()` | Sort rows |
| `repartition()` | Redistribute data across partitions |
| `coalesce()` | Reduce number of partitions |
| `cache()` / `persist()` | Mark for in-memory caching |
| `map()` | Apply function to each element (RDD) |
| `flatMap()` | map + flatten (RDD) |

### Actions (Eager — Trigger Execution)

| Action | Description |
|---|---|
| `show()` | Print first N rows |
| `collect()` | Return all rows to driver (careful with large data!) |
| `count()` | Count rows |
| `first()` / `head()` | Return first row |
| `take(n)` | Return first n rows |
| `write` | Write DataFrame to storage |
| `foreach()` | Apply function to each row (side effects) |
| `reduce()` | Aggregate to a single value (RDD) |
| `saveAsTextFile()` | Write RDD to text files |
| `toPandas()` | Convert to Pandas DataFrame (driver only!) |

---

## Jobs, Stages, and Tasks

```
Action call
    ↓
Job  (one per action)
    ↓
Stages  (separated by shuffle boundaries)
    ├── Stage 0: read + filter + map
    │       ↓ shuffle (groupBy/join triggers data exchange)
    └── Stage 1: aggregate + write
            ↓
        Tasks  (one per partition per stage)
```

```
Example: 100 partitions, 2 stages
→ 100 tasks in Stage 0
→ 200 tasks in Stage 1 (after repartition(200))
```

Watch jobs/stages/tasks in the **Spark UI** at `http://localhost:4040`.

---

## Partitions

A **partition** is a chunk of data processed by a single task on a single executor core.

```python
# Check partition count
df.rdd.getNumPartitions()        # → e.g., 8

# Set during read
df = spark.read.option("numPartitions", 50).jdbc(...)

# Repartition (full shuffle — expensive)
df = df.repartition(200)

# Repartition by column (ensures same keys go to same partition)
df = df.repartition("country")

# Coalesce (reduce partitions WITHOUT shuffle)
df = df.coalesce(10)
```

### How Many Partitions?

| Rule | Value |
|---|---|
| Default for shuffle | `spark.sql.shuffle.partitions` (default: 200) |
| Ideal partition size | 128MB–256MB of data |
| Min partitions | Number of executor cores |
| For large joins/groupBy | 2–4× total executor cores |

---

## Shuffles

A **shuffle** occurs when data must be redistributed across partitions — it's the most expensive operation in Spark.

Triggered by: `groupBy`, `join`, `repartition`, `distinct`, `sort`, `reduceByKey`

```
Before shuffle:       After shuffle:
Partition 0: A,C,B    Partition 0: A,A,A  (all A's together)
Partition 1: A,B,C    Partition 1: B,B,B
Partition 2: B,A,C    Partition 2: C,C,C
```

**Cost:** Data written to disk → sent over the network → read from disk on the other side.

**Minimize shuffles** by:
- Filtering early (reduce data before the shuffle)
- Using `broadcast` joins for small tables
- Caching DataFrames before multiple shuffles
- Tuning `spark.sql.shuffle.partitions`

---

## Caching and Persistence

Cache a DataFrame when you'll reuse it multiple times:

```python
df = spark.read.parquet("s3://data/events/")
df.cache()          # stores in memory (MEMORY_AND_DISK by default for DataFrames)
df.count()          # first action materializes the cache

# Use df multiple times — second+ time reads from cache, not disk
df.filter(...).show()
df.groupBy(...).count().show()

# Release cache
df.unpersist()
```

### Persistence Levels (RDD API)

| Level | Memory | Disk | CPU | Use Case |
|---|---|---|---|---|
| `MEMORY_ONLY` | Yes | No | Low | Fits in memory, fast reuse |
| `MEMORY_AND_DISK` | Yes | Spill | Low | Doesn't fit in memory |
| `DISK_ONLY` | No | Yes | Low | Very large, infrequent reuse |
| `MEMORY_ONLY_SER` | Yes (serialized) | No | Higher | Memory-constrained |
| `OFF_HEAP` | Off-heap | No | Medium | Avoid GC pressure |

---

## The Catalyst Optimizer

Catalyst is Spark's query optimizer. It automatically rewrites your query for efficiency:

1. **Analysis** — resolve column names and types
2. **Logical optimization** — predicate pushdown, constant folding, column pruning
3. **Physical planning** — choose join strategy (broadcast vs sort-merge vs hash)
4. **Code generation** — compile query to JVM bytecode (Tungsten)

```python
# See the execution plan
df.explain()              # physical plan
df.explain(True)          # all plans (parsed, analyzed, optimized, physical)
df.explain("formatted")   # nicely formatted
```

---

## Narrow vs Wide Transformations

| Type | Description | Shuffle? | Example |
|---|---|---|---|
| **Narrow** | Each output partition depends on one input partition | No | `filter`, `map`, `select` |
| **Wide** | Output partitions depend on multiple input partitions | Yes | `groupBy`, `join`, `sort` |

Narrow transformations can be **pipelined** within a stage (no I/O). Wide transformations always create a new stage boundary.

---

## SparkContext vs SparkSession

| | SparkContext | SparkSession |
|---|---|---|
| Introduced | Spark 1.0 | Spark 2.0 |
| Use for | RDD API | DataFrame/SQL API (also wraps SparkContext) |
| Access | `sc = spark.sparkContext` | `spark = SparkSession.builder...getOrCreate()` |
| Status | Still valid | Preferred entry point for everything |


---

# Chapter 4: RDDs — Resilient Distributed Datasets

## What is an RDD?

An **RDD** (Resilient Distributed Dataset) is the foundational abstraction of Spark:
- **Resilient** — fault-tolerant via lineage (can recompute lost partitions)
- **Distributed** — data is split across multiple nodes in the cluster
- **Dataset** — a collection of records (any Python/Java/Scala object)

> In modern Spark (2.0+), **use DataFrames** for most work. RDDs are for low-level control, non-tabular data, or when you need Python objects rather than rows.

---

## Creating RDDs

### From a Collection

```python
sc = spark.sparkContext

# Parallelize a Python list
rdd = sc.parallelize([1, 2, 3, 4, 5])
rdd = sc.parallelize(range(1000), numSlices=10)  # 10 partitions

# From a list of tuples (key-value RDD)
kv_rdd = sc.parallelize([("a", 1), ("b", 2), ("a", 3)])
```

### From Files

```python
# Text file (each line = one element)
rdd = sc.textFile("data/logs/*.txt")
rdd = sc.textFile("hdfs:///data/events/", minPartitions=50)

# JSON file (each line = JSON string)
rdd = sc.textFile("data/events.json")

# Whole text files (filename, content) pairs
rdd = sc.wholeTextFiles("data/configs/")

# Binary files
rdd = sc.binaryFiles("data/images/")

# Sequence files (Hadoop format)
rdd = sc.sequenceFile("hdfs:///data/seqfiles/")
```

### From a DataFrame

```python
rdd = df.rdd           # Row objects
rdd = df.rdd.map(lambda row: row.asDict())  # dict per row
```

---

## RDD Transformations

### Element-wise

```python
rdd = sc.parallelize([1, 2, 3, 4, 5])

# map — apply function to each element
squared = rdd.map(lambda x: x ** 2)
# → [1, 4, 9, 16, 25]

# filter — keep matching elements
evens = rdd.filter(lambda x: x % 2 == 0)
# → [2, 4]

# flatMap — map + flatten
words = sc.parallelize(["hello world", "foo bar"])
word_list = words.flatMap(lambda s: s.split())
# → ["hello", "world", "foo", "bar"]

# mapPartitions — apply function to each partition (more efficient than map for setup)
def process_partition(iter):
    import json
    for item in iter:
        yield json.loads(item)

result = rdd.mapPartitions(process_partition)

# mapPartitionsWithIndex
def with_index(idx, iter):
    for item in iter:
        yield (idx, item)

result = rdd.mapPartitionsWithIndex(with_index)
```

### Set Operations

```python
rdd1 = sc.parallelize([1, 2, 3, 4])
rdd2 = sc.parallelize([3, 4, 5, 6])

rdd1.union(rdd2).collect()        # → [1,2,3,4,3,4,5,6]
rdd1.intersection(rdd2).collect() # → [3,4]
rdd1.subtract(rdd2).collect()     # → [1,2]
rdd1.cartesian(rdd2).collect()    # → all (a,b) pairs
```

### Key-Value (Pair RDD) Transformations

```python
kv = sc.parallelize([("a", 1), ("b", 2), ("a", 3), ("b", 4)])

# reduceByKey — combine values with same key
kv.reduceByKey(lambda a, b: a + b).collect()
# → [("a", 4), ("b", 6)]

# groupByKey — group all values by key (less efficient than reduceByKey)
kv.groupByKey().mapValues(list).collect()
# → [("a", [1, 3]), ("b", [2, 4])]

# sortByKey
kv.sortByKey().collect()
# → [("a", 1), ("a", 3), ("b", 2), ("b", 4)]

# mapValues — transform values, keep keys
kv.mapValues(lambda v: v * 10).collect()
# → [("a", 10), ("b", 20), ("a", 30), ("b", 40)]

# join
rdd_a = sc.parallelize([("a", 1), ("b", 2)])
rdd_b = sc.parallelize([("a", "x"), ("b", "y")])
rdd_a.join(rdd_b).collect()
# → [("a", (1, "x")), ("b", (2, "y"))]

# leftOuterJoin, rightOuterJoin, fullOuterJoin
rdd_a.leftOuterJoin(rdd_b).collect()

# aggregateByKey — flexible aggregation with different accumulator type
zero_value = (0, 0)  # (sum, count)
def seq_op(acc, v):   return (acc[0]+v, acc[1]+1)
def comb_op(a, b):    return (a[0]+b[0], a[1]+b[1])
kv.aggregateByKey(zero_value, seq_op, comb_op).collect()
# → [("a", (4, 2)), ("b", (6, 2))]  → avg = sum/count
```

---

## RDD Actions

```python
rdd = sc.parallelize([3, 1, 4, 1, 5, 9, 2, 6])

rdd.collect()          # → [3,1,4,1,5,9,2,6]  — returns all to driver
rdd.count()            # → 8
rdd.first()            # → 3
rdd.take(3)            # → [3,1,4]
rdd.top(3)             # → [9,6,5]  — largest first
rdd.min()              # → 1
rdd.max()              # → 9
rdd.sum()              # → 31
rdd.mean()             # → 3.875

rdd.reduce(lambda a, b: a + b)     # → 31
rdd.fold(0, lambda a, b: a + b)    # → 31 (fold with initial zero-value)

# Stats
rdd.stats()    # count, mean, stdev, min, max

# Save
rdd.saveAsTextFile("output/results/")
rdd.saveAsPickleFile("output/results.pkl")

# foreach — apply side effect (no return value)
rdd.foreach(lambda x: print(x))

# countByValue — frequency map
rdd.countByValue()   # → {3:1, 1:2, 4:1, 5:1, 9:1, 2:1, 6:1}
```

---

## Persistence / Caching

```python
from pyspark import StorageLevel

rdd = sc.textFile("hdfs:///data/large_file.txt")

# Default cache (MEMORY_AND_DISK)
rdd.cache()

# Explicit storage level
rdd.persist(StorageLevel.MEMORY_AND_DISK_SER)

# Trigger materialization
count = rdd.count()

# Release
rdd.unpersist()
```

---

## RDD Lineage and Fault Tolerance

Spark tracks the lineage of each RDD — the sequence of transformations that produced it. If a partition is lost (node failure), Spark recomputes just that partition from the lineage.

```python
# View the lineage
print(rdd.toDebugString().decode())
```

Example lineage:
```
(2) MapPartitionsRDD[3]
 |  FilterRDD[2]
 |  MapPartitionsRDD[1]
 |  HadoopRDD[0] at textFile  ← source
```

For long lineages (many transformations), call `checkpoint()` to break the lineage and save a copy to HDFS:

```python
sc.setCheckpointDir("hdfs:///spark-checkpoints/")
rdd.checkpoint()
rdd.count()  # materializes the checkpoint
```

---

## When to Use RDDs vs DataFrames

| Use RDDs when... | Use DataFrames when... |
|---|---|
| Processing non-tabular data (images, binary, custom objects) | Processing structured/semi-structured data |
| Need low-level control over partitioning | Want Catalyst query optimization |
| Python lambda functions with complex logic | SQL-style operations |
| Working with existing RDD-based code | Performance matters (DataFrames are faster) |
| Custom partitioners | Most production ETL and analytics |


---

# Chapter 5: DataFrames & Datasets

## What is a DataFrame?

A **DataFrame** is a distributed collection of data organized into named columns — like a table in a relational database or a Pandas DataFrame, but running across a cluster. It's the primary abstraction for data engineering in PySpark.

```
+-------+----+----------+-------+
|  name | age|   city   | score |
+-------+----+----------+-------+
| Alice |  30| New York |  92.5 |
| Bob   |  25| London   |  88.0 |
| Carol |  35| Tokyo    |  95.1 |
+-------+----+----------+-------+
```

DataFrames run through the **Catalyst optimizer** — always faster than equivalent RDD code.

---

## Creating DataFrames

### From a List of Dicts / Tuples

```python
from pyspark.sql import SparkSession
from pyspark.sql.types import *

spark = SparkSession.builder.master("local[*]").appName("demo").getOrCreate()

# From list of dicts (schema inferred)
data = [
    {"name": "Alice", "age": 30, "city": "New York"},
    {"name": "Bob",   "age": 25, "city": "London"},
]
df = spark.createDataFrame(data)

# From list of tuples with explicit schema
schema = StructType([
    StructField("name",  StringType(),  nullable=False),
    StructField("age",   IntegerType(), nullable=True),
    StructField("score", DoubleType(),  nullable=True),
])
rows = [("Alice", 30, 92.5), ("Bob", 25, 88.0)]
df = spark.createDataFrame(rows, schema=schema)
```

### From Files

```python
# CSV
df = spark.read.csv("data/users.csv", header=True, inferSchema=True)
df = spark.read.option("header", True).option("sep", ";").csv("data.csv")

# JSON
df = spark.read.json("data/events/*.json")

# Parquet (recommended for production)
df = spark.read.parquet("s3://bucket/data/")

# ORC
df = spark.read.orc("hdfs:///data/warehouse/")

# Delta Lake
df = spark.read.format("delta").load("s3://bucket/delta/products/")
```

### Schema Inspection

```python
df.printSchema()
# root
#  |-- name:  string (nullable = true)
#  |-- age:   integer (nullable = true)
#  |-- score: double  (nullable = true)

df.schema          # StructType object
df.columns         # ['name', 'age', 'score']
df.dtypes          # [('name', 'string'), ('age', 'int'), ...]
df.count()         # number of rows
df.show(5)         # print first 5 rows
df.show(truncate=False)  # don't truncate long strings
df.describe().show()    # stats for numeric columns
```

---

## Column Operations

```python
from pyspark.sql import functions as F
from pyspark.sql.types import IntegerType

# Select columns
df.select("name", "age").show()
df.select(df.name, df["age"]).show()
df.select(F.col("name"), F.col("age")).show()

# Add or replace a column
df = df.withColumn("score_rounded", F.round("score", 0))
df = df.withColumn("age_group", F.when(F.col("age") < 30, "young").otherwise("senior"))

# Rename
df = df.withColumnRenamed("score", "test_score")

# Cast type
df = df.withColumn("age", F.col("age").cast(IntegerType()))

# Drop columns
df = df.drop("internal_id", "temp_field")

# Select with expressions
df.selectExpr("name", "age * 2 as double_age", "upper(city) as city_upper").show()
```

---

## Filtering

```python
# filter / where (same function)
df.filter(F.col("age") > 25).show()
df.where("age > 25 AND city = 'London'").show()

# Multiple conditions
df.filter((F.col("age") > 20) & (F.col("score") >= 90)).show()
df.filter((F.col("city") == "London") | (F.col("city") == "Paris")).show()

# Null handling
df.filter(F.col("score").isNull()).show()
df.filter(F.col("score").isNotNull()).show()

# isin
df.filter(F.col("city").isin("London", "Paris", "Berlin")).show()
df.filter(~F.col("status").isin(["cancelled", "refunded"])).show()

# String match
df.filter(F.col("name").startswith("A")).show()
df.filter(F.col("email").endswith("@gmail.com")).show()
df.filter(F.col("description").contains("wireless")).show()
df.filter(F.col("name").rlike("^[A-Z].*")).show()  # regex
```

---

## GroupBy & Aggregations

```python
# Single aggregation
df.groupBy("city").count().show()

# Multiple aggregations
df.groupBy("city", "category").agg(
    F.count("*").alias("total"),
    F.avg("score").alias("avg_score"),
    F.max("score").alias("max_score"),
    F.min("age").alias("min_age"),
    F.sum("revenue").alias("total_revenue"),
    F.countDistinct("user_id").alias("unique_users"),
    F.collect_list("name").alias("names"),
    F.first("created_at").alias("first_seen"),
).show()

# Filter after groupBy
df.groupBy("city").count().filter(F.col("count") > 100).show()
```

---

## Joins

```python
users    = spark.read.parquet("data/users/")
orders   = spark.read.parquet("data/orders/")

# Inner join
result = users.join(orders, on="user_id", how="inner")

# Left join
result = users.join(orders, on="user_id", how="left")

# Right / outer
result = users.join(orders, on="user_id", how="right")
result = users.join(orders, on="user_id", how="outer")  # full outer

# Join on multiple columns
result = users.join(orders,
    on=["user_id", "account_type"],
    how="inner"
)

# Join with different column names
result = users.join(orders,
    users["id"] == orders["customer_id"],
    how="left"
)

# Cross join
result = df1.crossJoin(df2)

# Semi join (return left rows that have a match, no right columns)
result = users.join(orders, on="user_id", how="left_semi")

# Anti join (return left rows that have NO match)
unordered = users.join(orders, on="user_id", how="left_anti")
```

---

## Sorting

```python
df.orderBy("age").show()
df.orderBy(F.col("age").desc()).show()
df.orderBy(F.col("score").desc(), F.col("name").asc()).show()
df.sort("city", ascending=False).show()
```

---

## Deduplication

```python
df.distinct().show()                       # deduplicate all columns
df.dropDuplicates(["user_id"]).show()     # deduplicate by specific columns
df.dropDuplicates(["email", "date"]).show()
```

---

## Handling Nulls

```python
# Drop rows with any null
df.dropna().show()

# Drop rows where specific columns are null
df.dropna(subset=["score", "age"]).show()

# Fill nulls
df.fillna(0).show()                        # all numeric columns → 0
df.fillna({"score": 0.0, "city": "unknown"}).show()

# Replace with F.coalesce
df.withColumn("score", F.coalesce("score", F.lit(0.0))).show()

# Fill nulls with column computed value
mean_score = df.select(F.avg("score")).first()[0]
df = df.fillna({"score": mean_score})
```

---

## User-Defined Functions (UDFs)

```python
from pyspark.sql.functions import udf
from pyspark.sql.types import StringType

# Simple UDF
@udf(returnType=StringType())
def grade(score):
    if score is None:
        return "N/A"
    if score >= 90:  return "A"
    if score >= 80:  return "B"
    return "C"

df.withColumn("grade", grade("score")).show()

# Register for SQL use
spark.udf.register("grade_fn", grade)
spark.sql("SELECT name, grade_fn(score) FROM users").show()
```

> **Performance note:** Python UDFs are slow (serialize each row to Python, run, serialize back). Use Pandas UDFs (Chapter 16) for much better performance, or use built-in `F.*` functions wherever possible.

---

## Working with Arrays and Maps

```python
from pyspark.sql.functions import (
    array, array_contains, explode, explode_outer,
    size, sort_array, array_distinct, map_keys, map_values
)

# Arrays
df.withColumn("tags", array(F.lit("new"), F.lit("sale"))).show()
df.filter(array_contains("tags", "sale")).show()
df.withColumn("num_tags", size("tags")).show()
df.withColumn("tags_sorted", sort_array("tags")).show()

# Explode array → one row per element
df.withColumn("tag", explode("tags")).show()         # drops null arrays
df.withColumn("tag", explode_outer("tags")).show()   # keeps nulls as null

# Struct fields
df.select("address.city", "address.zip").show()

# Map type
df.withColumn("city", F.col("metadata")["city"]).show()
```

---

## Datasets (Scala / Java Only)

Datasets are strongly-typed DataFrames. In Python, everything is a DataFrame. In Scala:

```scala
case class User(name: String, age: Int, score: Double)
val ds: Dataset[User] = df.as[User]
ds.filter(_.age > 25).map(_.name).show()
```

Python/PySpark uses DataFrames exclusively.


---

# Chapter 6: Spark SQL

## What is Spark SQL?

**Spark SQL** lets you query DataFrames using standard SQL syntax. It's fully integrated with the DataFrame API — you can mix SQL and DataFrame operations freely.

```python
# Create a temporary view
df.createOrReplaceTempView("users")

# Query with SQL
result = spark.sql("""
    SELECT city, COUNT(*) as total, AVG(score) as avg_score
    FROM users
    WHERE age > 18
    GROUP BY city
    ORDER BY total DESC
    LIMIT 10
""")
result.show()
```

SQL results are DataFrames — you can chain DataFrame operations after a SQL query.

---

## Creating Views

```python
# Temporary view — session-scoped (disappears when session ends)
df.createOrReplaceTempView("users")

# Global temporary view — available across sessions in same app
df.createOrReplaceGlobalTempView("global_users")
spark.sql("SELECT * FROM global_temp.global_users").show()

# Permanent table (persisted to metastore — requires Hive or Glue)
df.write.saveAsTable("default.users")
spark.sql("SELECT * FROM default.users").show()
```

---

## SQL Syntax Reference

### SELECT and Filtering

```sql
-- Basic select
SELECT name, age, score FROM users;

-- Aliases
SELECT name AS full_name, score * 100 AS score_pct FROM users;

-- Distinct
SELECT DISTINCT city FROM users;

-- Filtering
SELECT * FROM users WHERE age BETWEEN 25 AND 45;
SELECT * FROM users WHERE city IN ('London', 'Paris', 'Berlin');
SELECT * FROM users WHERE name LIKE 'A%';
SELECT * FROM users WHERE score IS NOT NULL;

-- Conditional
SELECT name,
       CASE WHEN score >= 90 THEN 'A'
            WHEN score >= 80 THEN 'B'
            ELSE 'C' END AS grade
FROM users;
```

### Aggregations

```sql
SELECT
    city,
    COUNT(*)                AS total_users,
    COUNT(DISTINCT user_id) AS unique_users,
    AVG(score)              AS avg_score,
    MAX(score)              AS max_score,
    MIN(age)                AS min_age,
    SUM(revenue)            AS total_revenue,
    PERCENTILE(score, 0.95) AS p95_score
FROM users
GROUP BY city
HAVING COUNT(*) > 100
ORDER BY total_revenue DESC;
```

### Joins

```sql
-- Inner join
SELECT u.name, o.order_id, o.total
FROM users u
INNER JOIN orders o ON u.user_id = o.user_id;

-- Left join with null handling
SELECT u.name, COALESCE(o.total, 0) AS total
FROM users u
LEFT JOIN orders o ON u.user_id = o.user_id;

-- Cross join
SELECT * FROM dim_dates CROSS JOIN dim_products;
```

---

## Window Functions

Window functions compute values across a sliding window of rows. They're one of Spark SQL's most powerful features.

```python
from pyspark.sql import Window
from pyspark.sql import functions as F

# Define window
w = Window.partitionBy("city").orderBy(F.col("score").desc())

# Rank functions
df.withColumn("rank",        F.rank().over(w))       \
  .withColumn("dense_rank",  F.dense_rank().over(w)) \
  .withColumn("row_number",  F.row_number().over(w)) \
  .withColumn("percent_rank",F.percent_rank().over(w)) \
  .show()

# Lag / Lead
w2 = Window.partitionBy("user_id").orderBy("event_date")
df.withColumn("prev_score", F.lag("score",  1).over(w2)) \
  .withColumn("next_score", F.lead("score", 1).over(w2)) \
  .show()

# Running total / moving average
w3 = Window.partitionBy("user_id").orderBy("event_date") \
           .rowsBetween(Window.unboundedPreceding, Window.currentRow)

df.withColumn("cumulative_sum", F.sum("revenue").over(w3)).show()

# Moving average (last 7 rows)
w4 = Window.partitionBy("user_id").orderBy("event_date") \
           .rowsBetween(-6, 0)

df.withColumn("moving_avg_7d", F.avg("score").over(w4)).show()

# ntile (quantile buckets)
w5 = Window.orderBy("score")
df.withColumn("quartile", F.ntile(4).over(w5)).show()
```

### In SQL

```sql
SELECT
    name,
    city,
    score,
    RANK()  OVER (PARTITION BY city ORDER BY score DESC) AS city_rank,
    SUM(revenue) OVER (PARTITION BY city ORDER BY event_date
                       ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS running_total,
    LAG(score, 1) OVER (PARTITION BY user_id ORDER BY event_date) AS prev_score
FROM events;
```

---

## CTEs (WITH Clause)

```sql
WITH active_users AS (
    SELECT user_id, name, city
    FROM users
    WHERE status = 'active'
),
user_orders AS (
    SELECT user_id, SUM(total) AS lifetime_value
    FROM orders
    GROUP BY user_id
)
SELECT u.name, u.city, COALESCE(uo.lifetime_value, 0) AS ltv
FROM active_users u
LEFT JOIN user_orders uo USING (user_id)
ORDER BY ltv DESC;
```

---

## Explain Plan

Always check the execution plan for complex queries:

```python
df.explain()              # physical plan only
df.explain("formatted")   # clean readable format
df.explain("cost")        # with cost estimates
df.explain("codegen")     # JVM code generation
```

Example output analysis:
```
== Physical Plan ==
AdaptiveSparkPlan
+- SortMergeJoin [user_id]         ← expensive shuffle join
   :- Sort [user_id]
   :  +- Exchange hashpartitioning  ← shuffle!
   :     +- Filter (age > 18)       ← pushed down ✓
   :        +- Scan parquet users
   +- Sort [user_id]
      +- Exchange hashpartitioning  ← shuffle!
         +- Scan parquet orders
```

Look for:
- `BroadcastHashJoin` → good (no shuffle for small table)
- `SortMergeJoin` → shuffle required (make sure it's necessary)
- `Exchange` → shuffle happening
- `Filter` pushed down before `Scan` → good (predicate pushdown)

---

## JDBC: Reading from Databases

```python
# Read from PostgreSQL
df = spark.read \
    .format("jdbc") \
    .option("url",      "jdbc:postgresql://host:5432/mydb") \
    .option("dbtable",  "public.orders") \
    .option("user",     "spark_user") \
    .option("password", "secret") \
    .option("driver",   "org.postgresql.Driver") \
    .load()

# Parallel read — partition by a numeric column
df = spark.read \
    .format("jdbc") \
    .option("url",           "jdbc:postgresql://host:5432/mydb") \
    .option("dbtable",       "public.orders") \
    .option("partitionColumn","order_id") \
    .option("lowerBound",    "1") \
    .option("upperBound",    "10000000") \
    .option("numPartitions", "50") \
    .load()
```

---

## The Spark Catalog

```python
# List databases
spark.catalog.listDatabases()

# List tables
spark.catalog.listTables("default")

# List columns
spark.catalog.listColumns("users", "default")

# Check if table exists
spark.catalog.tableExists("default.users")

# Refresh table metadata (after external changes)
spark.catalog.refreshTable("default.users")

# Clear cache
spark.catalog.clearCache()
```

---

## Built-In SQL Functions Reference

```python
import pyspark.sql.functions as F

# String
F.upper("col"), F.lower("col"), F.length("col")
F.trim("col"), F.ltrim("col"), F.rtrim("col")
F.split("col", ","), F.concat_ws(",", "a", "b")
F.regexp_replace("col", r"\d+", "")
F.regexp_extract("col", r"(\d{4})", 1)
F.substring("col", 1, 3)
F.lpad("col", 10, "0"), F.rpad("col", 10, " ")
F.translate("col", "abc", "ABC")
F.instr("col", "sub")        # position of substring

# Numeric
F.abs("col"), F.round("col", 2), F.ceil("col"), F.floor("col")
F.sqrt("col"), F.pow("col", 2), F.log("col"), F.log2("col")
F.greatest("a", "b", "c"),  F.least("a", "b", "c")
F.rand(), F.randn()          # random numbers

# Date / Timestamp
F.current_date(), F.current_timestamp()
F.to_date("col", "yyyy-MM-dd")
F.to_timestamp("col", "yyyy-MM-dd HH:mm:ss")
F.date_add("col", 7), F.date_sub("col", 7)
F.datediff("end_date", "start_date")
F.months_between("col1", "col2")
F.date_format("col", "yyyy-MM")
F.year("col"), F.month("col"), F.dayofmonth("col")
F.hour("col"), F.minute("col"), F.second("col")
F.dayofweek("col"), F.weekofyear("col")
F.trunc("col", "month"), F.date_trunc("hour", "col")

# Conditional
F.when(F.col("age") > 18, "adult").otherwise("minor")
F.coalesce("a", "b", F.lit("default"))
F.isnull("col"), F.isnan("col")
F.nullif("col", F.lit(0))

# Hash / Encoding
F.md5("col"), F.sha2("col", 256)
F.hash("col"), F.xxhash64("col")
F.base64("col"), F.unbase64("col")
```


---

# Chapter 7: Reading & Writing Data

## The Unified Data Source API

Spark uses a unified API for all data sources:

```python
# Read
df = spark.read.format("parquet").option("key", "value").load("path")

# Write
df.write.format("parquet").option("key", "value").mode("overwrite").save("path")
```

Or use shorthand methods: `.read.parquet()`, `.read.csv()`, `.write.json()`, etc.

---

## CSV

### Reading

```python
df = spark.read.csv(
    path          = "data/sales/*.csv",
    header        = True,
    inferSchema   = True,
    sep           = ",",           # or ";" for European CSVs
    nullValue     = "N/A",
    nanValue      = "NaN",
    quote         = '"',
    escape        = '\\',
    multiLine     = True,          # allow newlines inside fields
    encoding      = "UTF-8",
    dateFormat    = "yyyy-MM-dd",
    timestampFormat = "yyyy-MM-dd HH:mm:ss",
)

# With explicit schema (always preferred in production)
from pyspark.sql.types import *
schema = StructType([
    StructField("id",         IntegerType(), True),
    StructField("name",       StringType(),  True),
    StructField("created_at", TimestampType(), True),
    StructField("revenue",    DoubleType(),  True),
])
df = spark.read.schema(schema).csv("data/sales/", header=True)
```

### Writing

```python
df.write.csv(
    path       = "output/sales/",
    header     = True,
    sep        = ",",
    mode       = "overwrite",   # overwrite|append|ignore|errorIfExists
    compression = "gzip",       # none|gzip|bz2|deflate|snappy
)
```

---

## JSON

### Reading

```python
# One JSON object per line (JSONL/NDJSON — default)
df = spark.read.json("data/events/*.json")

# Multiline JSON (one array/object spanning multiple lines)
df = spark.read.option("multiLine", True).json("data/config.json")

# With schema
df = spark.read.schema(schema).json("data/events/")
```

### Writing

```python
df.write.json("output/events/", mode="overwrite")
df.write.json("output/events/", compression="gzip")
```

---

## Parquet (Recommended Format)

Parquet is a **columnar, compressed** binary format. It's the default format for production Spark workloads.

**Why Parquet:**
- Columnar — read only the columns you need
- Compressed — typically 5–10× smaller than CSV
- Schema embedded — no need to specify schema on read
- Predicate pushdown — Spark can skip row groups that don't match filters

### Reading

```python
df = spark.read.parquet("s3://bucket/data/products/")

# Read specific partitioned path
df = spark.read.parquet("s3://bucket/data/events/date=2024-03-15/")
```

### Writing

```python
df.write.parquet(
    "output/products/",
    mode="overwrite",
    compression="snappy"    # snappy|gzip|lzo|zstd|none
)

# Partitioned write (creates subdirectories per partition key)
df.write \
    .partitionBy("year", "month") \
    .parquet("s3://bucket/data/events/", mode="overwrite")
```

Partitioned output:
```
s3://bucket/data/events/
├── year=2024/month=01/part-00000.parquet
├── year=2024/month=02/part-00000.parquet
└── year=2024/month=03/part-00000.parquet
```

Reading a partition:
```python
df = spark.read.parquet("s3://bucket/data/events/")
# Partition columns are automatically added
df.filter("year = 2024 AND month = 3").show()
```

---

## ORC

Similar to Parquet, optimized for Hive. Preferred in Hadoop/Hive ecosystems.

```python
df = spark.read.orc("hdfs:///warehouse/sales/")
df.write.orc("hdfs:///warehouse/output/", mode="overwrite")
```

---

## Delta Lake

Delta Lake adds ACID transactions, schema enforcement, and time travel on top of Parquet files stored in object storage.

```python
# Requires: pip install delta-spark

# Read
df = spark.read.format("delta").load("s3://bucket/delta/products/")

# Write (create or overwrite)
df.write.format("delta").mode("overwrite").save("s3://bucket/delta/products/")

# Append
df.write.format("delta").mode("append").save("s3://bucket/delta/products/")

# Merge (upsert)
from delta.tables import DeltaTable
delta_table = DeltaTable.forPath(spark, "s3://bucket/delta/products/")
delta_table.alias("target").merge(
    updates_df.alias("source"),
    "target.product_id = source.product_id"
).whenMatchedUpdateAll() \
 .whenNotMatchedInsertAll() \
 .execute()

# Time travel
df_yesterday = spark.read.format("delta") \
    .option("timestampAsOf", "2024-03-14") \
    .load("s3://bucket/delta/products/")

df_version_5 = spark.read.format("delta") \
    .option("versionAsOf", "5") \
    .load("s3://bucket/delta/products/")
```

---

## JDBC (Databases)

```python
# Read
df = spark.read.format("jdbc") \
    .option("url",      "jdbc:postgresql://host:5432/db") \
    .option("dbtable",  "(SELECT * FROM orders WHERE status='active') AS t") \
    .option("user",     "user") \
    .option("password", "pass") \
    .load()

# Write
df.write.format("jdbc") \
    .option("url",    "jdbc:postgresql://host:5432/db") \
    .option("dbtable","output_table") \
    .option("user",   "user") \
    .option("password","pass") \
    .mode("append") \
    .save()
```

---

## Avro

```python
# Requires: spark.jars.packages = org.apache.spark:spark-avro_2.12:3.5.0

df = spark.read.format("avro").load("data/events.avro")
df.write.format("avro").save("output/events/")
```

---

## Write Modes

| Mode | Behavior |
|---|---|
| `overwrite` | Delete existing data and write new data |
| `append` | Add new data to existing data |
| `ignore` | Do nothing if data already exists |
| `errorIfExists` | Throw error if data already exists (default) |

---

## Partitioned Writes — Tips

```python
# Control number of output files with coalesce/repartition
df.coalesce(1).write.parquet("output/small/")          # 1 file (for small data)
df.repartition(50).write.parquet("output/large/")      # 50 files

# Repartition by key before partitioned write
# (ensures each partition has all rows for that key)
df.repartition("country") \
  .write \
  .partitionBy("country") \
  .parquet("output/by_country/")

# Limit files per partition with maxRecordsPerFile
df.write \
    .option("maxRecordsPerFile", 1000000) \
    .partitionBy("date") \
    .parquet("output/")
```

---

## Reading Multiple Sources at Once

```python
# Union multiple paths
df = spark.read.parquet(
    "s3://bucket/data/year=2022/",
    "s3://bucket/data/year=2023/",
    "s3://bucket/data/year=2024/",
)

# Glob patterns
df = spark.read.parquet("s3://bucket/data/year=202*/")
df = spark.read.json("data/events/2024-{01,02,03}-*.json")
```

---

## Schema Evolution

```python
# Parquet: merge schemas across files with different columns
df = spark.read \
    .option("mergeSchema", "true") \
    .parquet("s3://bucket/data/")

# Delta: handle schema changes automatically
spark.conf.set("spark.databricks.delta.schema.autoMerge.enabled", "true")
df.write.format("delta").option("mergeSchema", "true").mode("append").save(path)
```


---

# Chapter 8: Transformations & Actions

## Common DataFrame Transformations Reference

### Column Functions

```python
from pyspark.sql import functions as F

# Arithmetic
df.withColumn("profit", F.col("revenue") - F.col("cost"))
df.withColumn("margin", (F.col("revenue") - F.col("cost")) / F.col("revenue"))

# String manipulation
df.withColumn("upper_name",  F.upper("name"))
df.withColumn("name_length", F.length("name"))
df.withColumn("first_name",  F.split("full_name", " ")[0])
df.withColumn("email_domain",F.regexp_extract("email", r"@(.+)$", 1))

# Date operations
df.withColumn("year",          F.year("created_at"))
df.withColumn("month",         F.month("created_at"))
df.withColumn("days_active",   F.datediff(F.current_date(), "signup_date"))
df.withColumn("week_start",    F.date_trunc("week", "event_time"))

# Struct / complex
df.withColumn("full_name",     F.struct("first_name", "last_name"))
df.withColumn("tags",          F.array(F.lit("new"), F.lit("sale")))
df.withColumn("metadata",      F.create_map(F.lit("env"), F.lit("prod")))
```

---

## Join Types Deep Dive

```
DataFrame A (users):          DataFrame B (orders):
user_id | name                user_id | total
--------|------               --------|------
1       | Alice               1       | 100
2       | Bob                 1       | 200
3       | Carol               4       | 50     ← no matching user

INNER  JOIN → {(1,Alice,100), (1,Alice,200)}      only matches
LEFT   JOIN → {(1,Alice,100), (1,Alice,200), (2,Bob,null), (3,Carol,null)}
RIGHT  JOIN → {(1,Alice,100), (1,Alice,200), (null,null,50)}
OUTER  JOIN → all of LEFT + RIGHT
LEFT SEMI  → {(1,Alice)}                          users WITH orders
LEFT ANTI  → {(2,Bob), (3,Carol)}                 users WITHOUT orders
```

```python
# Broadcast join — for small lookup tables
from pyspark.sql.functions import broadcast

result = large_df.join(
    broadcast(small_lookup_df),
    on="product_id",
    how="left"
)
```

> **Broadcast join threshold** (default: 10MB):
> ```python
> spark.conf.set("spark.sql.autoBroadcastJoinThreshold", 50 * 1024 * 1024)  # 50MB
> ```

---

## Repartition vs Coalesce

| | `repartition(n)` | `coalesce(n)` |
|---|---|---|
| Shuffle | Full shuffle (expensive) | Minimal shuffle |
| Direction | Increase or decrease | Only decrease |
| Output balance | Perfectly balanced | May be unbalanced |
| Use case | Before a wide transformation | Before writing output |

```python
# Increase partitions (before heavy shuffle operations)
df = df.repartition(400)

# Repartition by column (co-locate related data)
df = df.repartition("country")
df = df.repartition(200, "country", "category")

# Decrease partitions (before writing small output)
df = df.coalesce(10)

# Single file output
df.coalesce(1).write.parquet("output/")
```

---

## Pivot

Transform rows into columns:

```python
df = spark.createDataFrame([
    ("Alice", "Q1", 100), ("Alice", "Q2", 150),
    ("Bob",   "Q1", 200), ("Bob",   "Q2", 180),
], ["name", "quarter", "revenue"])

df.groupBy("name").pivot("quarter", ["Q1", "Q2", "Q3", "Q4"]) \
    .sum("revenue") \
    .fillna(0) \
    .show()

# Output:
# +-----+----+----+----+----+
# | name|  Q1|  Q2|  Q3|  Q4|
# +-----+----+----+----+----+
# |Alice| 100| 150|   0|   0|
# |  Bob| 200| 180|   0|   0|
# +-----+----+----+----+----+
```

---

## Union and Set Operations

```python
# Union (includes duplicates)
combined = df1.union(df2)
combined = df1.unionAll(df2)         # alias for union

# Union by name (handles different column order)
combined = df1.unionByName(df2)

# Union with missing columns filled as null
combined = df1.unionByName(df2, allowMissingColumns=True)

# Intersect (rows in both)
common = df1.intersect(df2)

# Except / Subtract (rows in df1 but not df2)
diff = df1.exceptAll(df2)
```

---

## Sampling

```python
# Random sample — fraction of rows
sample = df.sample(fraction=0.1, seed=42)
sample = df.sample(withReplacement=True, fraction=0.5, seed=42)

# Stratified sample — maintain class distribution
fractions = {"class_A": 0.1, "class_B": 0.2, "class_C": 0.05}
sample = df.sampleBy("label", fractions=fractions, seed=42)

# Limit rows
top_100 = df.limit(100)
```

---

## Apply Functions: `map`, `applyInPandas`, `mapInPandas`

```python
# Row-level transformation via map (RDD route)
result = df.rdd.map(lambda row: Row(
    name=row.name.upper(),
    score=row.score * 1.1
)).toDF()

# applyInPandas — apply pandas function to each group
def normalize_group(pdf):
    pdf["score"] = (pdf["score"] - pdf["score"].mean()) / pdf["score"].std()
    return pdf

result = df.groupBy("city").applyInPandas(
    normalize_group,
    schema=df.schema
)

# mapInPandas — apply pandas function to each partition
def process_batch(iterator):
    for pdf in iterator:
        pdf["processed"] = pdf["value"].apply(lambda x: x * 2)
        yield pdf

result = df.mapInPandas(process_batch, schema="value int, processed int")
```

---

## Actions Reference

```python
# Collect data to driver (use carefully — OOM risk for large DataFrames)
rows   = df.collect()                  # list of Row objects
df_pd  = df.toPandas()                 # Pandas DataFrame (driver memory)
first  = df.first()                    # first Row
sample = df.head(5)                    # list of first 5 Rows
n      = df.count()                    # row count

# Aggregation actions
df.agg({"score": "avg", "revenue": "sum"}).show()

# Write actions
df.write.parquet("output/")            # write to files
df.write.format("delta").save("path/") # write Delta table
df.createOrReplaceTempView("my_view")  # create SQL view (lazy)

# foreach — side effect per partition
def write_to_db(row):
    # send each row to an external system
    ...
df.foreach(write_to_db)

# foreachPartition — side effect per partition (one connection per partition)
def write_partition(rows):
    conn = create_db_connection()
    for row in rows:
        conn.insert(row)
    conn.commit()
df.foreachPartition(write_partition)
```

---

## Broadcast Variables and Accumulators

### Broadcast Variable (Read-Only Shared Data)

```python
# Send large lookup dictionary to all executors once
city_to_region = {"London": "EU", "Paris": "EU", "New York": "US"}
bc_lookup = spark.sparkContext.broadcast(city_to_region)

# Use in a UDF
@udf(returnType=StringType())
def get_region(city):
    return bc_lookup.value.get(city, "Unknown")

df.withColumn("region", get_region("city")).show()
```

### Accumulator (Distributed Counter)

```python
# Count bad records in a job
bad_records = spark.sparkContext.accumulator(0)

def process_row(row):
    if row["score"] < 0:
        bad_records.add(1)

df.foreach(process_row)
print(f"Bad records found: {bad_records.value}")
```

---

## Practical Patterns

### Normalize Column Names

```python
import re

def normalize_columns(df):
    return df.toDF(*[re.sub(r'[^a-zA-Z0-9_]', '_', c.lower()) for c in df.columns])

df = normalize_columns(df)
```

### Add Audit Columns

```python
from pyspark.sql import functions as F

def add_audit_cols(df, source: str):
    return df \
        .withColumn("_loaded_at",  F.current_timestamp()) \
        .withColumn("_source",     F.lit(source)) \
        .withColumn("_row_hash",   F.md5(F.concat_ws("|", *df.columns)))
```

### Deduplicate with Row Number

```python
from pyspark.sql import Window

w = Window.partitionBy("user_id").orderBy(F.col("updated_at").desc())

deduped = df \
    .withColumn("rn", F.row_number().over(w)) \
    .filter(F.col("rn") == 1) \
    .drop("rn")
```


---

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


---

# Chapter 10: Cluster Architecture

## The Big Picture

```
┌───────────────────────────────────────────────────────────────┐
│                        DRIVER                                  │
│                                                                │
│   SparkSession → SparkContext → DAGScheduler → TaskScheduler  │
│                       ↓                                        │
│                   Spark UI (port 4040)                         │
└───────────────────────────┬───────────────────────────────────┘
                            │ (request resources)
                            ▼
┌───────────────────────────────────────────────────────────────┐
│                    CLUSTER MANAGER                             │
│            (YARN / Kubernetes / Standalone / Mesos)            │
└────────┬──────────────────┬─────────────────┬─────────────────┘
         │                  │                 │
         ▼                  ▼                 ▼
   ┌──────────┐       ┌──────────┐       ┌──────────┐
   │Executor 1│       │Executor 2│       │Executor 3│
   │ 4 threads│       │ 4 threads│       │ 4 threads│
   │ 8GB RAM  │       │ 8GB RAM  │       │ 8GB RAM  │
   └──────────┘       └──────────┘       └──────────┘
```

---

## Driver

The **Driver** is the JVM process that runs your `main()` / Python script. It:

- Creates the `SparkSession` / `SparkContext`
- Parses your code and builds the logical plan
- Sends the execution plan to the **DAGScheduler**
- Receives results from executors
- Runs the Spark UI (port 4040)

### DAGScheduler

Converts the logical plan into a **DAG of Stages**:
- Finds where shuffle boundaries exist → splits into stages
- Submits stage by stage, in dependency order

### TaskScheduler

Schedules individual **tasks** within a stage:
- Assigns tasks to executors
- Handles task retries on failure
- Applies data locality preferences (PROCESS_LOCAL > NODE_LOCAL > RACK_LOCAL > ANY)

---

## Executors

Executors are JVM processes on worker nodes. Each executor:
- Has a fixed number of **CPU cores** and **memory**
- Runs one **task per core** in parallel
- Holds a **Block Manager** for caching (RDD/DataFrame partitions)
- Communicates shuffle data via the **Shuffle Service**

### Executor Memory Layout

```
executor-memory (e.g., 8g)
├── Reserved Memory              ~300MB (Spark internal use)
├── Unified Memory               spark.memory.fraction (default 60%)
│   ├── Execution Memory         (shuffle buffers, sort, aggregation)
│   └── Storage Memory           (cached RDDs/DataFrames)
└── User Memory                  remaining (your Python objects, UDF data)
```

```
spark.executor.memory:     8g
spark.memory.fraction:     0.6   → 6g for execution + storage
spark.memory.storageFraction: 0.5 → 3g reserved for cache
```

---

## Cluster Managers

### Standalone

Spark's built-in cluster manager. Simple setup, good for dedicated Spark clusters.

```bash
# Start master
$SPARK_HOME/sbin/start-master.sh

# Start workers
$SPARK_HOME/sbin/start-worker.sh spark://master:7077

# Submit
spark-submit --master spark://master:7077 my_job.py
```

### YARN (Hadoop Ecosystem)

Most common in enterprise environments (HDP, CDH, AWS EMR).

```bash
# Client mode: Driver runs on submitting machine
spark-submit --master yarn --deploy-mode client my_job.py

# Cluster mode: Driver runs on YARN worker node
spark-submit --master yarn --deploy-mode cluster my_job.py
```

| | Client Mode | Cluster Mode |
|---|---|---|
| Driver runs on | Submitting machine | YARN worker node |
| Good for | Interactive development | Production jobs |
| Logs | Local terminal | YARN log aggregation |

### Kubernetes

See Chapter 14 for full details.

```bash
spark-submit \
  --master k8s://https://cluster:6443 \
  --deploy-mode cluster \
  --conf spark.kubernetes.container.image=my-spark:3.5 \
  local:///opt/spark/my_job.py
```

---

## Deploy Modes

| Mode | Driver Location | Use Case |
|---|---|---|
| `local[N]` | Same process | Development, testing |
| `local[*]` | Same process | All cores on local machine |
| `client` | Machine running spark-submit | Interactive, debugging |
| `cluster` | Worker node chosen by cluster manager | Production |

```python
# Local with 4 cores
spark = SparkSession.builder.master("local[4]").getOrCreate()

# Local with all cores
spark = SparkSession.builder.master("local[*]").getOrCreate()
```

---

## Dynamic Resource Allocation

Spark automatically scales executors up/down based on workload:

```python
spark.conf.set("spark.dynamicAllocation.enabled",             "true")
spark.conf.set("spark.dynamicAllocation.minExecutors",        "2")
spark.conf.set("spark.dynamicAllocation.maxExecutors",        "50")
spark.conf.set("spark.dynamicAllocation.initialExecutors",    "5")
spark.conf.set("spark.dynamicAllocation.executorIdleTimeout", "60s")
spark.conf.set("spark.shuffle.service.enabled",               "true")
```

When enabled:
- Adds executors when tasks are queued
- Removes executors idle for > `executorIdleTimeout`
- Requires External Shuffle Service (so shuffle data isn't lost when executors are removed)

---

## Data Locality

Spark tries to run tasks on the node where the data lives to minimize network transfer:

| Locality Level | Meaning |
|---|---|
| `PROCESS_LOCAL` | Data in executor's JVM cache |
| `NODE_LOCAL` | Data on same node (HDFS block or disk) |
| `RACK_LOCAL` | Data on same rack |
| `NO_PREF` | No location preference |
| `ANY` | Any node (requires network transfer) |

```python
spark.conf.set("spark.locality.wait",         "3s")
spark.conf.set("spark.locality.wait.process", "0s")
spark.conf.set("spark.locality.wait.node",    "2s")
spark.conf.set("spark.locality.wait.rack",    "2s")
```

---

## Spark UI

The Spark UI (port 4040 by default) shows:

| Tab | Contents |
|---|---|
| **Jobs** | All jobs, duration, stages per job |
| **Stages** | Task counts, shuffle read/write, failures |
| **Storage** | Cached RDDs/DataFrames (size, partitions) |
| **Environment** | All Spark configs |
| **Executors** | Per-executor memory, tasks, shuffle |
| **SQL** | Query plans, execution details per query |
| **Streaming** | Micro-batch durations, input rates |

History Server (for completed jobs):
```bash
$SPARK_HOME/sbin/start-history-server.sh
# Access: http://master:18080
```

---

## Key Configurations Summary

```python
# Resources
spark.conf.set("spark.executor.instances",  "10")
spark.conf.set("spark.executor.cores",      "4")
spark.conf.set("spark.executor.memory",     "8g")
spark.conf.set("spark.driver.memory",       "4g")
spark.conf.set("spark.driver.maxResultSize","2g")

# Parallelism
spark.conf.set("spark.default.parallelism",          "400")
spark.conf.set("spark.sql.shuffle.partitions",        "400")

# Network
spark.conf.set("spark.network.timeout",               "800s")
spark.conf.set("spark.executor.heartbeatInterval",    "60s")

# Speculation (retry slow tasks)
spark.conf.set("spark.speculation",                   "true")
spark.conf.set("spark.speculation.multiplier",        "1.5")
spark.conf.set("spark.speculation.quantile",          "0.9")
```


---

# Chapter 11: Performance & Optimization

## The Performance Checklist

```
1. Read less data     → partition pruning, predicate pushdown, column pruning
2. Shuffle less       → broadcast joins, filter before groupBy
3. Cache wisely       → cache reused DataFrames, unpersist when done
4. Right-size         → correct number of partitions, executor sizing
5. Use AQE            → let Adaptive Query Execution handle the rest
```

---

## Adaptive Query Execution (AQE)

AQE (Spark 3.0+) automatically optimizes query plans at runtime. Enable it:

```python
spark.conf.set("spark.sql.adaptive.enabled", "true")   # default True in 3.2+
```

AQE automatically:
- **Coalesces shuffle partitions** — merges small post-shuffle partitions
- **Converts sort-merge joins to broadcast joins** — if one side shrinks after filtering
- **Skew join optimization** — splits large skewed partitions

```python
# Tune AQE
spark.conf.set("spark.sql.adaptive.coalescePartitions.enabled",    "true")
spark.conf.set("spark.sql.adaptive.coalescePartitions.minPartitionNum", "1")
spark.conf.set("spark.sql.adaptive.advisoryPartitionSizeInBytes",  "128m")
spark.conf.set("spark.sql.adaptive.skewJoin.enabled",              "true")
spark.conf.set("spark.sql.adaptive.skewJoin.skewedPartitionFactor","5")
```

---

## Broadcast Joins

Use broadcast joins when one DataFrame is small (fits in executor memory):

```python
from pyspark.sql.functions import broadcast

# Force broadcast join
result = large_orders.join(broadcast(small_lookup), on="product_id")

# Tune auto-broadcast threshold (default 10MB)
spark.conf.set("spark.sql.autoBroadcastJoinThreshold", 50 * 1024 * 1024)  # 50MB

# Disable auto-broadcast (forces sort-merge join)
spark.conf.set("spark.sql.autoBroadcastJoinThreshold", -1)
```

> **Rule of thumb:** If one side of the join is < 100MB, broadcast it.

---

## Partition Tuning

### Shuffle Partitions

Default `spark.sql.shuffle.partitions = 200` is rarely optimal:

```python
# Too many small partitions (200 for a 1GB dataset)
spark.conf.set("spark.sql.shuffle.partitions", "200")

# Better: 2-4 × total executor cores, or aim for 128-256MB per partition
total_cores = 10 * 4  # 10 executors × 4 cores
spark.conf.set("spark.sql.shuffle.partitions", str(total_cores * 3))

# With AQE enabled, set it high and let AQE coalesce
spark.conf.set("spark.sql.shuffle.partitions",             "2000")
spark.conf.set("spark.sql.adaptive.coalescePartitions.enabled", "true")
```

### Input Partitions

```python
# Check current partitions
print(df.rdd.getNumPartitions())

# Repartition if too few (before a heavy join or groupBy)
df = df.repartition(400)

# Coalesce if too many (before writing output)
df = df.coalesce(50)

# Repartition by key for better join performance
orders   = orders.repartition("user_id")
users    = users.repartition("user_id")
result   = orders.join(users, on="user_id")
```

---

## Caching

```python
# Cache at MEMORY_AND_DISK (default for DataFrames)
df.cache()
df.count()        # materialize the cache

# Explicit storage level
from pyspark import StorageLevel
df.persist(StorageLevel.MEMORY_AND_DISK_SER)

# Always unpersist when done
df.unpersist()

# Check cache storage in Spark UI → Storage tab
```

**When to cache:**
- DataFrame used 2+ times in the same job
- Iterative algorithms (ML training)
- DataFrame is the result of an expensive aggregation

**When NOT to cache:**
- Used only once
- So large it evicts other cached data (thrashing)
- At the start of a pipeline (cache after filtering, not before)

---

## Handling Data Skew

Skew occurs when some partitions have far more data than others — a few tasks take 10× longer.

### Detect Skew

```python
# Check partition sizes
df.groupBy(F.spark_partition_id()).count().orderBy("count", ascending=False).show()
```

### Solution 1: Salt the Key

```python
import pyspark.sql.functions as F

NUM_SALTS = 10

# Add salt to the large side
large = large_df.withColumn(
    "salted_key",
    F.concat_ws("_", F.col("skewed_key"), (F.rand() * NUM_SALTS).cast("int"))
)

# Explode the small side
from pyspark.sql.functions import array, explode, lit
small = small_df.withColumn("salt", array(*[lit(i) for i in range(NUM_SALTS)])) \
                .withColumn("salt_val", explode("salt")) \
                .withColumn("salted_key", F.concat_ws("_", F.col("key"), F.col("salt_val"))) \
                .drop("salt", "salt_val")

result = large.join(small, on="salted_key")
```

### Solution 2: AQE Skew Join (Automatic in Spark 3.0+)

```python
spark.conf.set("spark.sql.adaptive.skewJoin.enabled", "true")
spark.conf.set("spark.sql.adaptive.skewJoin.skewedPartitionFactor",       "5")
spark.conf.set("spark.sql.adaptive.skewJoin.skewedPartitionThresholdInBytes", "256m")
```

### Solution 3: Repartition by the Skewed Column

```python
df = df.repartition(200, "skewed_key")   # distributes by hash
```

---

## Filter and Projection Pushdown

Push filters as early as possible — reduce data before expensive operations:

```python
# Bad: full table scan then filter after join
result = large.join(other, "id").filter("country = 'US'")

# Good: filter before join
result = large.filter("country = 'US'").join(other, "id")

# Parquet predicate pushdown (automatic — Spark skips row groups)
df = spark.read.parquet("s3://data/")
df.filter("year = 2024 AND month = 3").explain()
# Look for: PushedFilters: [IsNotNull(year), EqualTo(year,2024), ...]
```

---

## Column Pruning

Select only the columns you need before joins or shuffles:

```python
# Bad: carry all columns through a heavy join
result = large_df.join(other_df, "id")

# Good: select needed columns first
result = large_df.select("id", "country", "revenue") \
                 .join(other_df.select("id", "name"), "id")
```

---

## Avoid UDFs — Use Built-in Functions

Python UDFs deserialize each row to Python, process it, and serialize back. This is 10–100× slower than built-in functions.

```python
# Slow: Python UDF
@udf(returnType=StringType())
def clean(s):
    return s.strip().lower() if s else None

df.withColumn("clean_name", clean("name"))

# Fast: built-in equivalent
df.withColumn("clean_name", F.lower(F.trim("name")))
```

Use **Pandas UDFs** (vectorized) when you need Python logic (see Chapter 16).

---

## Executor Sizing

**Fat executor vs many small executors:**

| Approach | Cores/Executor | Notes |
|---|---|---|
| Too thin | 1 core | Wastes memory (no data sharing between threads) |
| Optimal | 4–5 cores | Good balance of parallelism and memory efficiency |
| Too fat | 10+ cores | GC pressure, Python broadcasting slow |

```
10 worker nodes × 16 cores × 64GB RAM each

One executor per node (too fat):
→ 10 executors, 16 cores, 60GB each  ← GC issues

Better approach:
→ 3 executors per node = 30 executors, 5 cores, 19GB each
→ Total parallelism: 30 × 5 = 150 tasks

spark.conf.set("spark.executor.cores", "5")
spark.conf.set("spark.executor.memory", "19g")
spark.conf.set("spark.executor.instances", "30")
```

---

## Kryo Serialization

Switch from Java serialization to Kryo for 5–10× faster serialization:

```python
spark.conf.set("spark.serializer", "org.apache.spark.serializer.KryoSerializer")
spark.conf.set("spark.kryoserializer.buffer.max", "2000m")
```

---

## Off-Heap Memory

Avoid GC pressure by storing caches off-heap:

```python
spark.conf.set("spark.memory.offHeap.enabled",  "true")
spark.conf.set("spark.memory.offHeap.size",      "8g")
```

---

## Performance Anti-Patterns

| Anti-Pattern | Problem | Fix |
|---|---|---|
| `df.collect()` on large data | Driver OOM | Use `write` or `take(n)` |
| Unbounded `crossJoin` | Explodes to rows² | Add join condition |
| Nested UDFs | Extremely slow | Use built-in functions |
| Many small files input | Too many tasks, driver overhead | Compact files, use `maxFilesPerTrigger` |
| `repartition` before `coalesce` | Pointless shuffle | Use only `coalesce` to reduce |
| `shuffle.partitions=200` always | Wrong for data size | Tune per-job |
| Cache everything | Evicts useful caches | Cache only reused DataFrames |


---

# Chapter 12: Spark MLlib

## What is MLlib?

**MLlib** is Spark's scalable machine learning library. It runs on top of DataFrames and provides:

- **Feature engineering** (transformers)
- **ML algorithms** (estimators)
- **Pipelines** (chain transformers + estimators)
- **Model evaluation** and **hyperparameter tuning**
- **Model persistence** (save/load)

---

## The Pipeline API

MLlib follows the scikit-learn `Pipeline` pattern:

```
Transformer → transforms data (no fitting)
Estimator  → fits on data → produces a Transformer (the model)
Pipeline   → chains multiple Transformers and Estimators
```

```python
from pyspark.ml import Pipeline
from pyspark.ml.feature import (
    StringIndexer, OneHotEncoder, VectorAssembler,
    StandardScaler, Imputer
)
from pyspark.ml.classification import RandomForestClassifier
from pyspark.ml.evaluation import BinaryClassificationEvaluator

# 1. Define stages
indexer = StringIndexer(inputCol="category", outputCol="category_idx")
encoder = OneHotEncoder(inputCol="category_idx", outputCol="category_vec")
imputer = Imputer(inputCols=["age", "income"], outputCols=["age_imp", "income_imp"])

assembler = VectorAssembler(
    inputCols=["age_imp", "income_imp", "category_vec"],
    outputCol="features"
)

scaler = StandardScaler(inputCol="features", outputCol="scaled_features")
rf = RandomForestClassifier(featuresCol="scaled_features", labelCol="label", numTrees=100)

# 2. Build pipeline
pipeline = Pipeline(stages=[indexer, encoder, imputer, assembler, scaler, rf])

# 3. Train / Test split
train, test = df.randomSplit([0.8, 0.2], seed=42)

# 4. Fit
model = pipeline.fit(train)

# 5. Predict
predictions = model.transform(test)
predictions.select("label", "prediction", "probability").show()

# 6. Evaluate
evaluator = BinaryClassificationEvaluator(labelCol="label")
auc = evaluator.evaluate(predictions)
print(f"AUC: {auc:.4f}")

# 7. Save and load
model.save("models/rf_pipeline/")
from pyspark.ml import PipelineModel
loaded = PipelineModel.load("models/rf_pipeline/")
```

---

## Feature Engineering

### StringIndexer and OneHotEncoder

```python
# StringIndexer: "cat" → 0.0, "dog" → 1.0
indexer = StringIndexer(inputCol="animal", outputCol="animal_idx", handleInvalid="keep")
df = indexer.fit(df).transform(df)

# OneHotEncoder: 0.0 → [1,0], 1.0 → [0,1]
encoder = OneHotEncoder(inputCol="animal_idx", outputCol="animal_vec")
df = encoder.fit(df).transform(df)
```

### VectorAssembler

```python
assembler = VectorAssembler(
    inputCols=["feature1", "feature2", "feature3", "one_hot_vec"],
    outputCol="features",
    handleInvalid="skip"
)
df = assembler.transform(df)
```

### Scalers

```python
from pyspark.ml.feature import StandardScaler, MinMaxScaler, RobustScaler

# StandardScaler: zero mean, unit std
scaler = StandardScaler(inputCol="features", outputCol="scaled", withMean=True, withStd=True)
df = scaler.fit(df).transform(df)

# MinMaxScaler: scale to [0, 1]
scaler = MinMaxScaler(inputCol="features", outputCol="scaled", min=0.0, max=1.0)
```

### Text Features

```python
from pyspark.ml.feature import Tokenizer, HashingTF, IDF, Word2Vec, CountVectorizer

tokenizer = Tokenizer(inputCol="text", outputCol="words")

# TF-IDF
hashing   = HashingTF(inputCol="words", outputCol="tf", numFeatures=10000)
idf       = IDF(inputCol="tf", outputCol="tfidf")

# Word2Vec embeddings
w2v = Word2Vec(inputCol="words", outputCol="embedding", vectorSize=100, minCount=2)

# CountVectorizer (vocabulary-based)
cv = CountVectorizer(inputCol="words", outputCol="features", vocabSize=1000)
```

### Other Transformers

```python
from pyspark.ml.feature import (
    Bucketizer, QuantileDiscretizer,
    PCA, ChiSqSelector,
    SQLTransformer, Normalizer
)

# Bucketizer: bin continuous values
bucketizer = Bucketizer(
    splits=[0, 18, 25, 35, 60, float("inf")],
    inputCol="age", outputCol="age_bucket"
)

# PCA: dimensionality reduction
pca = PCA(k=50, inputCol="features", outputCol="pca_features")

# Feature selection with Chi-Square
selector = ChiSqSelector(numTopFeatures=20, featuresCol="features",
                          labelCol="label", outputCol="selected_features")

# SQL-based transformation
sql_t = SQLTransformer(statement="SELECT *, age * income AS age_income FROM __THIS__")
```

---

## Classification Algorithms

```python
from pyspark.ml.classification import (
    LogisticRegression,
    RandomForestClassifier,
    GBTClassifier,
    LinearSVC,
    NaiveBayes,
    DecisionTreeClassifier,
    MultilayerPerceptronClassifier,
)

# Logistic Regression
lr = LogisticRegression(
    featuresCol="features", labelCol="label",
    maxIter=100, regParam=0.01, elasticNetParam=0.8
)

# Random Forest
rf = RandomForestClassifier(
    featuresCol="features", labelCol="label",
    numTrees=200, maxDepth=10, featureSubsetStrategy="sqrt"
)

# Gradient Boosted Trees
gbt = GBTClassifier(
    featuresCol="features", labelCol="label",
    maxIter=50, maxDepth=5, stepSize=0.1
)

# Neural network (layers: input → hidden → output)
mlp = MultilayerPerceptronClassifier(
    featuresCol="features", labelCol="label",
    layers=[100, 64, 32, 2],    # input dim, hidden, output classes
    maxIter=100
)
```

---

## Regression Algorithms

```python
from pyspark.ml.regression import (
    LinearRegression,
    RandomForestRegressor,
    GBTRegressor,
    DecisionTreeRegressor,
    GeneralizedLinearRegression,
)

lr = LinearRegression(
    featuresCol="features", labelCol="price",
    regParam=0.1, elasticNetParam=0.5
)

rf_reg = RandomForestRegressor(
    featuresCol="features", labelCol="price",
    numTrees=100, maxDepth=8
)
```

---

## Clustering

```python
from pyspark.ml.clustering import KMeans, BisectingKMeans, GaussianMixture

# KMeans
kmeans = KMeans(featuresCol="features", k=5, seed=42, maxIter=20)
model  = kmeans.fit(df)
df     = model.transform(df)    # adds "prediction" column (cluster id)

# Evaluate (WSSSE — within-cluster sum of squared errors)
print(model.summary.trainingCost)

# Find optimal k with elbow method
costs = []
for k in range(2, 11):
    model = KMeans(k=k, seed=42).fit(df)
    costs.append((k, model.summary.trainingCost))
```

---

## Model Evaluation

```python
from pyspark.ml.evaluation import (
    BinaryClassificationEvaluator,
    MulticlassClassificationEvaluator,
    RegressionEvaluator,
    ClusteringEvaluator,
)

# Binary classification
bin_eval = BinaryClassificationEvaluator(
    labelCol="label", rawPredictionCol="rawPrediction",
    metricName="areaUnderROC"  # or "areaUnderPR"
)
auc = bin_eval.evaluate(predictions)

# Multiclass
mc_eval = MulticlassClassificationEvaluator(
    labelCol="label", predictionCol="prediction",
    metricName="accuracy"  # accuracy|f1|weightedPrecision|weightedRecall
)

# Regression
reg_eval = RegressionEvaluator(
    labelCol="price", predictionCol="prediction",
    metricName="rmse"   # rmse|mse|mae|r2
)
```

---

## Hyperparameter Tuning

```python
from pyspark.ml.tuning import ParamGridBuilder, CrossValidator, TrainValidationSplit

# Build grid
param_grid = ParamGridBuilder() \
    .addGrid(rf.numTrees, [50, 100, 200]) \
    .addGrid(rf.maxDepth, [5, 10, 15]) \
    .build()

# Cross-validation (k-fold)
cv = CrossValidator(
    estimator  = pipeline,
    estimatorParamMaps = param_grid,
    evaluator  = BinaryClassificationEvaluator(),
    numFolds   = 3,
    parallelism = 4,   # run N param sets in parallel
)
cv_model = cv.fit(train)
best_model = cv_model.bestModel

# Train-validation split (faster than CV)
tvs = TrainValidationSplit(
    estimator  = pipeline,
    estimatorParamMaps = param_grid,
    evaluator  = BinaryClassificationEvaluator(),
    trainRatio = 0.8,
)
tvs_model = tvs.fit(train)
```

---

## Feature Importance

```python
# After fitting a Random Forest model
rf_model = model.stages[-1]   # last stage in pipeline
importances = rf_model.featureImportances
feature_names = assembler.getInputCols()

for name, importance in sorted(zip(feature_names, importances), key=lambda x: -x[1]):
    print(f"{name}: {importance:.4f}")
```


---

# Chapter 13: Delta Lake

## What is Delta Lake?

**Delta Lake** is an open-source storage layer that brings ACID transactions and reliability to data lakes (S3, ADLS, GCS, HDFS). It sits on top of Parquet files and adds:

- **ACID transactions** — concurrent writes are safe
- **Schema enforcement** — rejects bad data at write time
- **Schema evolution** — add columns without rewriting
- **Time travel** — query historical versions of your data
- **MERGE** — efficient upsert (no need to rewrite entire table)
- **Streaming + batch unified** — same table for both

```bash
pip install delta-spark
```

```python
spark = SparkSession.builder \
    .config("spark.sql.extensions",
            "io.delta.sql.DeltaSparkSessionExtension") \
    .config("spark.sql.catalog.spark_catalog",
            "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
    .getOrCreate()
```

---

## Creating a Delta Table

```python
# Write as Delta
df.write.format("delta").mode("overwrite").save("s3://bucket/delta/products/")

# Or create with SQL
spark.sql("""
    CREATE TABLE products (
        product_id STRING,
        name       STRING,
        price      DOUBLE,
        category   STRING,
        updated_at TIMESTAMP
    )
    USING DELTA
    LOCATION 's3://bucket/delta/products/'
    PARTITIONED BY (category)
""")
```

---

## Reading Delta Tables

```python
# Read latest version
df = spark.read.format("delta").load("s3://bucket/delta/products/")

# Read as SQL table (if registered in metastore)
df = spark.table("products")
spark.sql("SELECT * FROM products WHERE category = 'electronics'").show()
```

---

## Time Travel

Every write to a Delta table creates a new **version**. You can query any past version:

```python
# Read a specific version
df_v3 = spark.read.format("delta") \
    .option("versionAsOf", 3) \
    .load("s3://bucket/delta/products/")

# Read as of a timestamp
df_old = spark.read.format("delta") \
    .option("timestampAsOf", "2024-03-01 00:00:00") \
    .load("s3://bucket/delta/products/")

# SQL time travel
spark.sql("""
    SELECT * FROM products VERSION AS OF 3
""")
spark.sql("""
    SELECT * FROM products TIMESTAMP AS OF '2024-03-01'
""")

# View history
from delta.tables import DeltaTable
dt = DeltaTable.forPath(spark, "s3://bucket/delta/products/")
dt.history().show(truncate=False)
```

---

## MERGE (Upsert)

The most powerful Delta feature — efficiently upsert, update, or delete rows:

```python
from delta.tables import DeltaTable

delta_table = DeltaTable.forPath(spark, "s3://bucket/delta/products/")

# Full upsert: update existing, insert new
delta_table.alias("target").merge(
    updates_df.alias("source"),
    "target.product_id = source.product_id"
) \
.whenMatchedUpdateAll() \
.whenNotMatchedInsertAll() \
.execute()

# Selective merge: only update specific columns
delta_table.alias("t").merge(
    source_df.alias("s"),
    "t.order_id = s.order_id"
) \
.whenMatchedUpdate(set={
    "t.status":     "s.status",
    "t.updated_at": "s.updated_at"
}) \
.whenNotMatchedInsert(values={
    "order_id":    "s.order_id",
    "status":      "s.status",
    "created_at":  "s.created_at",
    "updated_at":  "s.updated_at",
}) \
.execute()

# Merge with delete
delta_table.alias("t").merge(
    changes_df.alias("s"),
    "t.id = s.id"
) \
.whenMatchedDelete(condition="s.is_deleted = true") \
.whenMatchedUpdateAll(condition="s.is_deleted = false") \
.whenNotMatchedInsertAll(condition="s.is_deleted = false") \
.execute()
```

---

## Schema Enforcement and Evolution

### Schema Enforcement (Default On)

```python
# This will FAIL if df has different schema than the Delta table
df.write.format("delta").mode("append").save(path)
# → AnalysisException: A schema mismatch was detected
```

### Schema Evolution (Opt-in)

```python
# Add new columns automatically
df.write.format("delta") \
    .option("mergeSchema", "true") \
    .mode("append") \
    .save(path)

# Globally
spark.conf.set("spark.databricks.delta.schema.autoMerge.enabled", "true")
```

---

## Delta Table Operations

### UPDATE

```python
delta_table.update(
    condition="category = 'electronics'",
    set={"price": "price * 0.9"}
)

# SQL
spark.sql("""
    UPDATE products SET price = price * 0.9 WHERE category = 'electronics'
""")
```

### DELETE

```python
delta_table.delete("status = 'discontinued'")

# SQL
spark.sql("DELETE FROM products WHERE status = 'discontinued'")
```

### VACUUM (Remove Old Files)

Delta keeps old Parquet files for time travel. Vacuum removes files older than the retention period:

```python
delta_table.vacuum()              # default: 7 days
delta_table.vacuum(retentionHours=168)  # 7 days

# SQL
spark.sql("VACUUM products RETAIN 168 HOURS")

# Dry run (see what would be deleted)
spark.sql("VACUUM products RETAIN 168 HOURS DRY RUN")
```

### OPTIMIZE (Compact Small Files)

```python
# Compact all small files
delta_table.optimize().executeCompaction()

# Z-ORDER (colocate related data for better query pruning)
delta_table.optimize().executeZOrderBy("product_id", "category")

# SQL
spark.sql("OPTIMIZE products ZORDER BY (product_id, category)")
```

---

## Delta Streaming

Delta tables work seamlessly with Structured Streaming:

### Write Stream to Delta

```python
query = stream_df.writeStream \
    .format("delta") \
    .outputMode("append") \
    .option("checkpointLocation", "s3://bucket/checkpoints/") \
    .start("s3://bucket/delta/events/")
```

### Read Stream from Delta

```python
# Continuous stream of new rows added to a Delta table
stream_df = spark.readStream \
    .format("delta") \
    .option("maxFilesPerTrigger", 1) \
    .load("s3://bucket/delta/events/")
```

---

## Delta Lake Architecture

```
Delta Table = Parquet files + _delta_log/ directory

s3://bucket/delta/products/
├── _delta_log/
│   ├── 00000000000000000000.json   ← version 0 (CREATE)
│   ├── 00000000000000000001.json   ← version 1 (INSERT)
│   ├── 00000000000000000002.json   ← version 2 (UPDATE)
│   └── 00000000000000000010.checkpoint.parquet  ← checkpoint
├── part-00000.snappy.parquet
├── part-00001.snappy.parquet
└── ...
```

Each `.json` log entry records:
- Files added
- Files removed
- Schema changes
- Statistics (min/max per column per file — enables data skipping)

---

## Data Skipping and Statistics

Delta automatically collects min/max statistics per Parquet file per column. Queries with range filters skip entire files:

```python
# This query skips all files where max(price) < 100 or min(price) > 500
spark.sql("SELECT * FROM products WHERE price BETWEEN 100 AND 500")

# Force stats collection on existing table
delta_table.optimize().executeCompaction()
```

---

## Change Data Feed (CDF)

Track row-level changes (insert/update/delete) for CDC pipelines:

```python
# Enable CDF
spark.sql("""
    ALTER TABLE products
    SET TBLPROPERTIES (delta.enableChangeDataFeed = true)
""")

# Read changes between versions
changes = spark.read.format("delta") \
    .option("readChangeFeed", "true") \
    .option("startingVersion", 0) \
    .load("s3://bucket/delta/products/")

# Changes have extra column: _change_type = insert|update_preimage|update_postimage|delete
changes.filter("_change_type = 'update_postimage'").show()
```


---

# Chapter 14: Spark on Kubernetes

## Why Spark on Kubernetes?

Running Spark on Kubernetes gives you:
- **Unified infrastructure** — same cluster as the rest of your apps
- **Dynamic resource allocation** — no dedicated Spark cluster
- **Container isolation** — reproducible environments
- **Cloud-native** — works on EKS, GKE, AKS, on-premise

---

## How It Works

```
spark-submit → Driver Pod created → Requests Executor Pods
                      ↓
            K8s API Server allocates pods
                      ↓
            Executor Pods start, run tasks, terminate
                      ↓
            Driver Pod completes, terminates
```

No persistent cluster — pods are ephemeral. Each job gets its own pods.

---

## Basic spark-submit on Kubernetes

```bash
spark-submit \
  --master k8s://https://k8s-api-server:6443 \
  --deploy-mode cluster \
  --name my-spark-job \
  --conf spark.kubernetes.container.image=apache/spark:3.5.0 \
  --conf spark.kubernetes.namespace=spark-jobs \
  --conf spark.executor.instances=5 \
  --conf spark.executor.memory=4g \
  --conf spark.executor.cores=2 \
  --conf spark.driver.memory=2g \
  local:///opt/spark/examples/jars/spark-examples.jar \
  10000000
```

For a PySpark job:

```bash
spark-submit \
  --master k8s://https://k8s-api-server:6443 \
  --deploy-mode cluster \
  --conf spark.kubernetes.container.image=my-registry/pyspark:3.5 \
  --conf spark.kubernetes.namespace=spark-jobs \
  --conf spark.kubernetes.authenticate.driver.serviceAccountName=spark \
  --conf spark.executor.instances=10 \
  local:///opt/spark/jobs/my_job.py
```

---

## Custom Docker Image

```dockerfile
FROM apache/spark:3.5.0

# Install Python dependencies
USER root
COPY requirements.txt /opt/spark/requirements.txt
RUN pip3 install --no-cache-dir -r /opt/spark/requirements.txt

# Copy job files
COPY src/ /opt/spark/jobs/

USER spark
```

```bash
docker build -t my-registry/pyspark:3.5 .
docker push my-registry/pyspark:3.5
```

---

## RBAC Setup

Spark needs permission to create/delete pods:

```yaml
# spark-rbac.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: spark
  namespace: spark-jobs
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: spark-role
  namespace: spark-jobs
rules:
  - apiGroups: [""]
    resources: ["pods", "services", "configmaps", "persistentvolumeclaims"]
    verbs: ["create", "get", "watch", "list", "delete", "patch"]
  - apiGroups: [""]
    resources: ["pods/log"]
    verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: spark-rolebinding
  namespace: spark-jobs
subjects:
  - kind: ServiceAccount
    name: spark
    namespace: spark-jobs
roleRef:
  kind: Role
  name: spark-role
  apiGroup: rbac.authorization.k8s.io
```

```bash
kubectl apply -f spark-rbac.yaml
```

---

## Dynamic Resource Allocation on K8s

Scale executors automatically without pre-allocating:

```bash
spark-submit \
  --conf spark.dynamicAllocation.enabled=true \
  --conf spark.dynamicAllocation.shuffleTracking.enabled=true \
  --conf spark.dynamicAllocation.minExecutors=1 \
  --conf spark.dynamicAllocation.maxExecutors=50 \
  --conf spark.dynamicAllocation.initialExecutors=5 \
  --conf spark.dynamicAllocation.executorIdleTimeout=60s \
  ...
```

> Note: On Kubernetes, use `shuffleTracking` instead of External Shuffle Service.

---

## Spark Operator (Kubernetes-Native)

The **Spark Operator** (Google's open-source controller) manages Spark jobs as Kubernetes CRDs:

```bash
helm install spark-operator spark-operator/spark-operator \
  --namespace spark-operator \
  --create-namespace
```

```yaml
# spark-job.yaml
apiVersion: sparkoperator.k8s.io/v1beta2
kind: SparkApplication
metadata:
  name: my-spark-job
  namespace: spark-jobs
spec:
  type: Python
  pythonVersion: "3"
  mode: cluster
  image: "my-registry/pyspark:3.5"
  imagePullPolicy: Always
  mainApplicationFile: "local:///opt/spark/jobs/etl.py"
  sparkVersion: "3.5.0"

  driver:
    cores: 2
    memory: "4g"
    serviceAccount: spark
    labels:
      app: my-spark-job

  executor:
    cores: 4
    instances: 10
    memory: "8g"
    labels:
      app: my-spark-job

  restartPolicy:
    type: OnFailure
    onFailureRetries: 3
```

```bash
kubectl apply -f spark-job.yaml
kubectl get sparkapplication -n spark-jobs
kubectl describe sparkapplication my-spark-job -n spark-jobs
```

---

## Accessing Data from S3

```bash
# Add S3 credentials as env vars in the image, or use IAM roles for service accounts
spark-submit \
  --conf spark.hadoop.fs.s3a.access.key=YOUR_KEY \
  --conf spark.hadoop.fs.s3a.secret.key=YOUR_SECRET \
  --conf spark.hadoop.fs.s3a.impl=org.apache.hadoop.fs.s3a.S3AFileSystem \
  ...
```

For EKS with IAM Roles for Service Accounts (IRSA):

```yaml
metadata:
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789:role/spark-s3-role
```

---

## Spark History Server on K8s

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: spark-history-server
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: history-server
          image: apache/spark:3.5.0
          command: ["/opt/spark/bin/spark-class"]
          args: ["org.apache.spark.deploy.history.HistoryServer"]
          env:
            - name: SPARK_HISTORY_OPTS
              value: "-Dspark.history.fs.logDirectory=s3a://bucket/spark-logs/"
          ports:
            - containerPort: 18080
```

---

## Resource Management Tips

```bash
# Limit resource usage with resource quotas
--conf spark.kubernetes.driver.request.cores=1
--conf spark.kubernetes.driver.limit.cores=2
--conf spark.kubernetes.executor.request.cores=2
--conf spark.kubernetes.executor.limit.cores=4

# Node selectors (run on specific node pool)
--conf spark.kubernetes.node.selector.workload=spark

# Tolerations (run on tainted nodes)
--conf spark.kubernetes.driver.tolerations=[{"key":"spark","operator":"Equal","value":"true","effect":"NoSchedule"}]

# Pod priority
--conf spark.kubernetes.driver.priorityClassName=high-priority
--conf spark.kubernetes.executor.priorityClassName=medium-priority
```


---

# Chapter 15: Spark on YARN & EMR

## Spark on YARN

YARN (Yet Another Resource Negotiator) is the resource manager in the Hadoop ecosystem. It's the most common cluster manager for Spark in enterprise environments (Cloudera, Hortonworks, Amazon EMR).

```
YARN Architecture:
  ResourceManager (master)  ← allocates containers
  NodeManagers (workers)    ← run containers
  ApplicationMaster         ← Spark Driver (in cluster mode)
```

---

## Deploy Modes

### Client Mode

```bash
spark-submit \
  --master yarn \
  --deploy-mode client \
  --executor-memory 4g \
  --executor-cores 2 \
  --num-executors 10 \
  --driver-memory 4g \
  my_job.py
```

- **Driver** runs on the machine that ran `spark-submit`
- Good for: interactive use, debugging (logs to your terminal)
- Bad for: production (if your laptop disconnects, job dies)

### Cluster Mode

```bash
spark-submit \
  --master yarn \
  --deploy-mode cluster \
  --executor-memory 4g \
  --executor-cores 2 \
  --num-executors 20 \
  --driver-memory 4g \
  --conf spark.yarn.queue=production \
  hdfs:///jobs/my_job.py
```

- **Driver** runs on a YARN worker node (inside YARN)
- Good for: production (job survives disconnect)
- Logs go to YARN aggregated logs

---

## YARN-Specific Configuration

```bash
# Queue management
--conf spark.yarn.queue=high_priority

# AM memory overhead
--conf spark.yarn.am.memory=2g
--conf spark.yarn.am.memoryOverhead=512m

# Executor overhead (off-heap: PySpark, overhead for libraries)
--conf spark.yarn.executor.memoryOverhead=1024m

# Max app attempts
--conf spark.yarn.maxAppAttempts=2

# Archive Python deps to YARN
--conf spark.yarn.dist.archives=my_venv.tar.gz#env
--conf spark.pyspark.python=./env/bin/python

# Ship files to all nodes
--files log4j.properties
--jars external-lib.jar
```

---

## YARN UI and Log Access

```bash
# YARN ResourceManager UI
http://yarn-rm:8088

# Find app ID
yarn application -list
yarn application -status application_1234567890_0001

# View logs (after job completes)
yarn logs -applicationId application_1234567890_0001

# Kill a job
yarn application -kill application_1234567890_0001
```

---

## Dynamic Allocation on YARN

```python
spark.conf.set("spark.dynamicAllocation.enabled",                "true")
spark.conf.set("spark.shuffle.service.enabled",                  "true")   # required for YARN DRA
spark.conf.set("spark.dynamicAllocation.minExecutors",           "2")
spark.conf.set("spark.dynamicAllocation.maxExecutors",           "100")
spark.conf.set("spark.dynamicAllocation.executorIdleTimeout",    "60s")
spark.conf.set("spark.dynamicAllocation.schedulerBacklogTimeout","1s")
```

The **External Shuffle Service** (`spark.shuffle.service.enabled=true`) runs on each YARN NodeManager and holds shuffle data when executors are released.

---

## AWS EMR

**Amazon EMR** (Elastic MapReduce) is AWS's managed Hadoop/Spark service. It provisions clusters on EC2, runs your job, and optionally terminates.

### Launch EMR Cluster (AWS CLI)

```bash
aws emr create-cluster \
  --name "SparkProductionCluster" \
  --release-label emr-6.15.0 \
  --applications Name=Spark Name=Hadoop Name=Hive \
  --instance-type m5.2xlarge \
  --instance-count 5 \
  --ec2-attributes KeyName=my-key-pair \
  --use-default-roles \
  --log-uri s3://my-bucket/emr-logs/ \
  --configurations file://spark-config.json
```

### Spark Configuration on EMR (`spark-config.json`)

```json
[
  {
    "Classification": "spark-defaults",
    "Properties": {
      "spark.executor.memory":           "6g",
      "spark.executor.cores":            "4",
      "spark.sql.shuffle.partitions":    "400",
      "spark.sql.adaptive.enabled":      "true",
      "spark.dynamicAllocation.enabled": "true"
    }
  },
  {
    "Classification": "spark",
    "Properties": {
      "maximizeResourceAllocation": "true"
    }
  }
]
```

### Submit a Job to EMR

```bash
# Add a step to a running cluster
aws emr add-steps \
  --cluster-id j-XXXXXXXXXXXX \
  --steps Type=spark,\
Name=MyJob,\
Args=[--deploy-mode,cluster,\
      --master,yarn,\
      s3://my-bucket/jobs/my_job.py,\
      --input,s3://my-bucket/input/,\
      --output,s3://my-bucket/output/],\
ActionOnFailure=CONTINUE
```

### EMR with S3A Committer (Faster S3 Writes)

```json
{
  "Classification": "spark-defaults",
  "Properties": {
    "spark.hadoop.fs.s3a.committer.name":               "directory",
    "spark.sql.sources.commitProtocolClass":             "org.apache.spark.internal.io.cloud.PathOutputCommitProtocol",
    "spark.sql.parquet.output.committer.class":          "org.apache.spark.internal.io.cloud.BindingParquetOutputCommitter",
    "spark.hadoop.mapreduce.outputcommitter.factory.scheme.s3a": "org.apache.hadoop.fs.s3a.commit.S3ACommitterFactory"
  }
}
```

---

## EMR Serverless (No Cluster Management)

Submit Spark jobs without managing a cluster:

```bash
# Create an application
aws emr-serverless create-application \
  --name "MySparkApp" \
  --type SPARK \
  --release-label emr-6.15.0

# Submit a job
aws emr-serverless start-job-run \
  --application-id APP_ID \
  --execution-role-arn arn:aws:iam::123456789:role/EMRServerlessRole \
  --job-driver '{
    "sparkSubmit": {
      "entryPoint": "s3://bucket/jobs/my_job.py",
      "sparkSubmitParameters": "--conf spark.executor.memory=4g"
    }
  }' \
  --configuration-overrides '{
    "monitoringConfiguration": {
      "s3MonitoringConfiguration": {
        "logUri": "s3://bucket/emr-logs/"
      }
    }
  }'
```

---

## Databricks (Cloud-Managed Spark)

Databricks is the company behind Spark — their platform provides:
- Fully managed clusters (auto-scaling, auto-termination)
- Collaborative notebooks
- Delta Lake as first-class citizen
- MLflow for ML lifecycle
- Unity Catalog for governance

```python
# Databricks Runtime: SparkSession is pre-created as `spark`

# Read from Databricks Unity Catalog
df = spark.table("my_catalog.my_schema.my_table")

# Read from DBFS or cloud storage
df = spark.read.parquet("dbfs:/mnt/my-bucket/data/")
df = spark.read.parquet("s3://my-bucket/data/")

# Write as Delta (default in Databricks)
df.write.format("delta").saveAsTable("my_catalog.my_schema.output_table")
```


---

# Chapter 16: PySpark Advanced

## Pandas UDFs (Vectorized UDFs)

Regular Python UDFs serialize one row at a time between JVM and Python. **Pandas UDFs** use Apache Arrow to transfer entire batches (Series/DataFrames) — typically 10–100× faster.

```python
import pandas as pd
from pyspark.sql.functions import pandas_udf
from pyspark.sql.types import DoubleType, StringType
```

### Scalar Pandas UDF (Series → Series)

```python
@pandas_udf(returnType=DoubleType())
def normalize_score(scores: pd.Series) -> pd.Series:
    return (scores - scores.mean()) / scores.std()

df.withColumn("norm_score", normalize_score("score")).show()

# Multiple input Series
@pandas_udf(DoubleType())
def weighted_avg(values: pd.Series, weights: pd.Series) -> pd.Series:
    return values * weights / weights.sum()

df.withColumn("weighted", weighted_avg("revenue", "weight")).show()
```

### Iterator of Series (one-time setup cost)

```python
from typing import Iterator

@pandas_udf(StringType())
def classify_batch(iterator: Iterator[pd.Series]) -> Iterator[pd.Series]:
    # Load model once per partition (not per row!)
    import pickle
    model = pickle.load(open("model.pkl", "rb"))

    for series in iterator:
        yield pd.Series(model.predict(series.values.reshape(-1, 1)))

df.withColumn("category", classify_batch("features")).show()
```

### Grouped Map Pandas UDF (applyInPandas)

Apply a function that returns a DataFrame to each group:

```python
def normalize_within_group(pdf: pd.DataFrame) -> pd.DataFrame:
    pdf["norm_score"] = (pdf["score"] - pdf["score"].mean()) / pdf["score"].std()
    return pdf

result = df.groupBy("city").applyInPandas(
    func=normalize_within_group,
    schema="name STRING, city STRING, score DOUBLE, norm_score DOUBLE"
)
```

### Grouped Aggregate Pandas UDF (Series → Scalar)

```python
from pyspark.sql.functions import pandas_udf
from pyspark.sql import functions as F

@pandas_udf(DoubleType())
def median_agg(series: pd.Series) -> float:
    return series.median()

df.groupBy("city").agg(median_agg("score").alias("median_score")).show()
```

---

## pandas API on Spark (formerly Koalas)

Write **Pandas code that runs distributed** on Spark. Ideal for data scientists familiar with Pandas who need to scale.

```python
import pyspark.pandas as ps

# Create pyspark pandas DataFrame
psdf = ps.DataFrame({"a": [1, 2, 3], "b": [4, 5, 6]})

# Or from a Spark DataFrame
psdf = spark_df.to_pandas_on_spark()

# From CSV (distributed)
psdf = ps.read_csv("s3://bucket/data/*.csv")

# Use Pandas-like API
psdf["c"] = psdf["a"] + psdf["b"]
psdf.groupby("category")["revenue"].sum()
psdf.sort_values("score", ascending=False).head(10)
psdf[psdf["age"] > 25].describe()

# Convert back to Spark DataFrame
spark_df = psdf.to_spark()

# Convert to Pandas (driver memory — use for small data only)
pdf = psdf.to_pandas()
```

### Default Index Warning

Pandas-on-Spark uses a distributed index by default. For some operations this requires a shuffle. Use a natural unique column as index:

```python
psdf = psdf.set_index("user_id")
```

---

## Type Hints in PySpark

PySpark 3.1+ supports Python type hints for better IDE support and readability:

```python
from pyspark.sql import DataFrame
from pyspark.sql import functions as F
from pyspark.sql.types import StructType, StructField, StringType, DoubleType

def transform_users(df: DataFrame) -> DataFrame:
    return df \
        .filter(F.col("age") > 18) \
        .withColumn("score_pct", F.col("score") / 100.0)

def get_schema() -> StructType:
    return StructType([
        StructField("name",  StringType(),  nullable=False),
        StructField("score", DoubleType(),  nullable=True),
    ])
```

---

## Arrow Optimization

Apache Arrow speeds up data transfer between JVM (Spark) and Python (pandas):

```python
spark.conf.set("spark.sql.execution.arrow.pyspark.enabled",         "true")
spark.conf.set("spark.sql.execution.arrow.pyspark.fallback.enabled","true")

# This is now much faster with Arrow enabled:
pdf = spark_df.toPandas()          # Spark → Pandas
spark_df = spark.createDataFrame(pdf)  # Pandas → Spark
```

---

## UDF Best Practices

| Approach | Speed | Use Case |
|---|---|---|
| Built-in `F.*` functions | Fastest | Standard operations |
| SQL expressions | Fastest | String-based transformations |
| Pandas UDF (vectorized) | Fast | Custom Python logic, ML inference |
| `applyInPandas` | Medium | Per-group transformations |
| Python UDF (row-by-row) | Slowest | Simple fallback, avoid if possible |

---

## Broadcast Variables

Share large lookup tables across executors without re-sending per task:

```python
# Python dict broadcast
lookup = {"electronics": 0.1, "clothing": 0.2, "books": 0.05}
bc_tax = spark.sparkContext.broadcast(lookup)

@pandas_udf(DoubleType())
def apply_tax(categories: pd.Series, prices: pd.Series) -> pd.Series:
    tax_rates = categories.map(bc_tax.value).fillna(0.0)
    return prices * (1 + tax_rates)

df.withColumn("price_with_tax", apply_tax("category", "price")).show()

# Free memory when done
bc_tax.destroy()
```

---

## Accumulators for Custom Metrics

```python
from pyspark import AccumulatorParam

# Custom accumulator for a list of strings
class ListAccumulator(AccumulatorParam):
    def zero(self, v): return []
    def addInPlace(self, v1, v2): return v1 + v2

errors = spark.sparkContext.accumulator([], ListAccumulator())

def validate_row(row):
    if row["age"] < 0:
        errors.add([f"Invalid age for user {row['user_id']}"])

df.foreach(validate_row)
print(f"Errors: {errors.value[:10]}")
```

---

## Efficient Python Code Patterns

### Reuse SparkSession

```python
# Singleton pattern — don't create multiple sessions
def get_spark():
    return SparkSession.builder.appName("MyApp").getOrCreate()
```

### Avoid `collect()` on Large DataFrames

```python
# Bad: pulls all data to driver
rows = df.collect()
for row in rows:
    process(row)

# Good: process on executors
df.foreach(process)
df.foreachPartition(process_partition)
df.write.parquet("output/")
```

### Use `show()` for debugging, `write` for production

```python
# Debug
df.show(5)
df.printSchema()

# Production
df.write.mode("overwrite").parquet("output/")
```

### Chain Transformations

```python
# Build a clear transformation pipeline
result = (
    spark.read.parquet("input/")
    .filter(F.col("date") >= "2024-01-01")
    .withColumn("year",  F.year("date"))
    .withColumn("month", F.month("date"))
    .groupBy("year", "month", "category")
    .agg(
        F.sum("revenue").alias("total_revenue"),
        F.count("*").alias("order_count"),
    )
    .orderBy("year", "month")
)
result.write.partitionBy("year", "month").parquet("output/")
```


---

# Chapter 17: Monitoring & Debugging

## Spark UI (Port 4040)

The Spark UI is the most important debugging tool. Access it at `http://driver-host:4040` while a job runs.

| Tab | What to Look For |
|---|---|
| **Jobs** | Job duration, failed stages |
| **Stages** | Task counts, shuffle read/write, skewed tasks |
| **Storage** | Cached DataFrames (size, fraction cached) |
| **Environment** | All active Spark configs |
| **Executors** | Per-executor heap, tasks completed, GC time |
| **SQL** | Query plan, duration per operator |
| **Streaming** | Batch processing time, input rate, throughput |

---

## Reading the Jobs Tab

```
Job 0  Duration: 2.3 min  Stages: 3/3
 └── Stage 0: 200 tasks  45s   (read + filter)
 └── Stage 1: 200 tasks  1m    (groupBy shuffle)
 └── Stage 2: 50 tasks   15s   (write)
```

**Red flags:**
- Tasks at end of stage taking much longer → skew
- Stage retries → node failures or OOM
- Duration much longer than input scan → expensive shuffle

---

## Reading the Stages Tab

```
Stage 1: groupBy
  Input:  12.5 GB  (200 tasks)
  Output: 4.2 GB   (200 tasks)
  Shuffle Write: 4.2 GB   ← before shuffle
  Shuffle Read:  4.2 GB   ← after shuffle

  Task Metrics:
    Min:    2s
    Median: 18s
    Max:    4m 32s   ← massive skew! one task much slower
```

---

## Common Problems and Fixes

### Problem: OutOfMemoryError (OOM)

```
ERROR ExecutorLostFailure: ... Container killed by YARN for exceeding memory limits
```

**Fixes:**
```python
# Increase executor memory overhead
spark.conf.set("spark.yarn.executor.memoryOverhead", "2g")

# Reduce executor cores (less concurrent tasks per executor)
spark.conf.set("spark.executor.cores", "2")

# Reduce shuffle partitions (smaller buffers per task)
spark.conf.set("spark.sql.shuffle.partitions", "100")

# Enable AQE to coalesce small partitions
spark.conf.set("spark.sql.adaptive.enabled", "true")
```

### Problem: GC Overhead

Symptom: tasks are slow, executor GC time > 20% in Spark UI

```python
# Reduce heap fraction for execution (more for user objects)
spark.conf.set("spark.memory.fraction", "0.6")

# Use serialized caching (less GC pressure)
from pyspark import StorageLevel
df.persist(StorageLevel.MEMORY_AND_DISK_SER)

# Use off-heap
spark.conf.set("spark.memory.offHeap.enabled", "true")
spark.conf.set("spark.memory.offHeap.size", "4g")

# Use G1GC
spark.conf.set("spark.executor.extraJavaOptions",
               "-XX:+UseG1GC -XX:G1HeapRegionSize=16m -XX:+G1PrintRegionLivenessInfo")
```

### Problem: Task Skew

Symptom: 199 tasks finish in 10s, 1 task takes 10 minutes.

```python
# Enable AQE skew handling
spark.conf.set("spark.sql.adaptive.skewJoin.enabled", "true")
spark.conf.set("spark.sql.adaptive.skewJoin.skewedPartitionFactor", "5")

# Or repartition with salt (see Chapter 11)
```

### Problem: Too Many Small Files (Input)

Symptom: thousands of tasks, each reading tiny files.

```python
# Merge small files on read
df = spark.read \
    .option("mergeSchema", "true") \
    .parquet("s3://bucket/data/")

# Use spark.sql.files.maxPartitionBytes (default 128MB)
spark.conf.set("spark.sql.files.maxPartitionBytes", str(256 * 1024 * 1024))  # 256MB

# Use spark.sql.files.openCostInBytes to avoid combining too many tiny files
spark.conf.set("spark.sql.files.openCostInBytes", str(4 * 1024 * 1024))  # 4MB
```

### Problem: Job Hangs / Deadlock

```bash
# Check active tasks
curl http://driver:4040/api/v1/applications/<app-id>/stages

# Check executor thread dumps
GET /executors/<executor-id>/threads  (in Spark UI → Executors → Thread Dump)

# Kill a stage
spark.sparkContext.cancelJobGroup("group-id")
```

---

## Logging

### Set Log Level

```python
spark.sparkContext.setLogLevel("WARN")   # ERROR, WARN, INFO, DEBUG
```

### Custom log4j.properties

```properties
# log4j.properties
log4j.rootCategory=WARN, console
log4j.appender.console=org.apache.log4j.ConsoleAppender
log4j.appender.console.layout=org.apache.log4j.PatternLayout
log4j.appender.console.layout.ConversionPattern=%d{yy/MM/dd HH:mm:ss} %p %c{1}: %m%n

# Reduce chatty loggers
log4j.logger.org.apache.spark=WARN
log4j.logger.org.apache.hadoop=WARN
log4j.logger.org.spark_project.jetty=WARN

# Your app logs at INFO
log4j.logger.com.mycompany=INFO
```

```bash
spark-submit --files log4j.properties --conf "spark.driver.extraJavaOptions=-Dlog4j.configuration=log4j.properties"
```

### Python logging in PySpark

```python
import logging
logger = logging.getLogger(__name__)

def my_transform(df):
    count = df.count()
    logger.info(f"Processing {count} rows")
    return df.filter(...)
```

---

## Event Logs and History Server

```python
# Enable event logs
spark.conf.set("spark.eventLog.enabled", "true")
spark.conf.set("spark.eventLog.dir",     "hdfs:///spark-logs")
spark.conf.set("spark.eventLog.compress","true")
```

Start history server:

```bash
$SPARK_HOME/sbin/start-history-server.sh \
  --properties-file spark-history-server.conf
```

```properties
# spark-history-server.conf
spark.history.fs.logDirectory=hdfs:///spark-logs
spark.history.ui.port=18080
spark.history.retainedApplications=500
```

---

## Spark Metrics

Spark exposes metrics via JMX, Prometheus, or Graphite:

```properties
# metrics.properties
*.sink.graphite.class=org.apache.spark.metrics.sink.GraphiteSink
*.sink.graphite.host=graphite-host
*.sink.graphite.port=2003
*.sink.graphite.period=10

# Prometheus
*.sink.prometheusServlet.class=org.apache.spark.metrics.sink.PrometheusServlet
*.sink.prometheusServlet.path=/metrics/prometheus
```

Key metrics to watch:
- `executor.jvmGcTime` → GC overhead
- `executor.cpuTime` → actual CPU vs wallclock
- `shuffleRead/Write.bytes` → shuffle volume
- `executor.memoryUsed` → heap usage
- `BlockManager.memory.remainingMem_MB` → cache pressure

---

## Debugging Strategies

### Step 1: Read the Error Message

```
Job aborted due to stage failure:
Task 12 in stage 3.0 failed 4 times, most recent failure:
Lost task 12.3 in stage 3.0 (TID 1234, worker-2, executor 5):
java.lang.OutOfMemoryError: GC overhead limit exceeded
```

Key info: `stage 3.0`, `Task 12`, `executor 5`, `OutOfMemoryError`.

### Step 2: Use `explain()`

```python
df.explain("formatted")   # see the physical plan
```

### Step 3: Inspect Partition Sizes

```python
df.groupBy(F.spark_partition_id().alias("part_id")) \
  .count() \
  .orderBy("count", ascending=False) \
  .show(20)
```

### Step 4: Sample and Print

```python
df.sample(0.001).show(100, truncate=False)
df.filter("year IS NULL").count()
df.select(F.col("amount").cast("double")).filter(F.col("amount").isNull()).count()
```

### Step 5: Cache Intermediate Results

```python
# Cache at a checkpoint to see where slowness begins
df.cache()
print(df.count())   # materialize
df.explain()        # check plan reads from cache
```


---

# Chapter 18: Security

## Security Layers in Spark

```
1. Authentication   ← who is running this job?
2. Authorization    ← what are they allowed to do?
3. Encryption       ← is data protected in transit and at rest?
4. Audit            ← who did what and when?
```

---

## Authentication

### Kerberos (Enterprise Hadoop)

Kerberos is the standard authentication mechanism in Hadoop clusters:

```bash
# Obtain a Kerberos ticket
kinit spark-user@REALM.COMPANY.COM

# Submit with Kerberos
spark-submit \
  --master yarn \
  --principal spark-user@REALM.COMPANY.COM \
  --keytab /etc/security/keytabs/spark.keytab \
  my_job.py
```

```python
# In spark-defaults.conf or job config
spark.conf.set("spark.kerberos.principal", "spark-user@REALM.COM")
spark.conf.set("spark.kerberos.keytab",    "/etc/security/keytabs/spark.keytab")
```

### Shared Secret (Standalone / Local)

```python
spark.conf.set("spark.authenticate",       "true")
spark.conf.set("spark.authenticate.secret","my-secret-key")
```

---

## Encryption in Transit (SSL/TLS)

### RPC Encryption (Node-to-Node)

```python
spark.conf.set("spark.network.crypto.enabled",   "true")
spark.conf.set("spark.network.crypto.keyLength",  "128")
spark.conf.set("spark.authenticate",              "true")
```

### SSL for Spark UI and History Server

```properties
# spark-defaults.conf
spark.ssl.enabled=true
spark.ssl.keyStore=/path/to/keystore.jks
spark.ssl.keyStorePassword=keystore_pass
spark.ssl.trustStore=/path/to/truststore.jks
spark.ssl.trustStorePassword=truststore_pass
spark.ssl.protocol=TLSv1.2
```

```python
spark.conf.set("spark.ui.filters",        "org.apache.hadoop.security.authentication.server.AuthenticationFilter")
spark.conf.set("spark.ssl.ui.enabled",    "true")
```

---

## Encryption at Rest

Data at rest encryption is handled at the storage layer:

### HDFS Transparent Encryption

```bash
# Create encryption zone in HDFS
hdfs crypto -createZone -keyName myKey -path /user/spark/data

# Verify
hdfs crypto -listZones
```

### S3 Server-Side Encryption

```python
# SSE-S3 (AWS managed keys)
spark.conf.set("spark.hadoop.fs.s3a.server-side-encryption-algorithm", "AES256")

# SSE-KMS (customer managed keys)
spark.conf.set("spark.hadoop.fs.s3a.server-side-encryption-algorithm", "SSE-KMS")
spark.conf.set("spark.hadoop.fs.s3a.server-side-encryption-key",
               "arn:aws:kms:us-east-1:123456789:key/my-key-id")

# Client-side encryption
spark.conf.set("spark.hadoop.fs.s3a.encryption.algorithm", "CSE-KMS")
```

---

## Credential Management

### Hadoop Credential Provider

Store passwords in an encrypted credential store instead of plain-text config:

```bash
# Create credential store
hadoop credential create spark.password.db \
  -provider jceks://hdfs/user/spark/credentials.jceks \
  -value mySecretPassword

# Use it in config
spark.conf.set("spark.hadoop.hadoop.security.credential.provider.path",
               "jceks://hdfs/user/spark/credentials.jceks")
```

### Vault / Secrets Manager

```python
# AWS Secrets Manager
import boto3, json

def get_db_credentials():
    client = boto3.client("secretsmanager", region_name="us-east-1")
    secret = client.get_secret_value(SecretId="spark/db-credentials")
    return json.loads(secret["SecretString"])

creds = get_db_credentials()
df = spark.read.format("jdbc") \
    .option("user",     creds["username"]) \
    .option("password", creds["password"]) \
    .load()
```

---

## ACLs and Authorization

### Spark UI ACLs

Restrict who can view the Spark UI:

```python
spark.conf.set("spark.acls.enable",      "true")
spark.conf.set("spark.ui.view.acls",     "alice,bob")
spark.conf.set("spark.ui.view.acls.groups", "data-team")
spark.conf.set("spark.admin.acls",       "spark-admin")
spark.conf.set("spark.admin.acls.groups","ops-team")
```

### Modify ACLs (Cancel Jobs, Kill Tasks)

```python
spark.conf.set("spark.modify.acls",      "alice")
spark.conf.set("spark.modify.acls.groups","ops-team")
```

---

## Secure Spark on Kubernetes

```yaml
# Use a dedicated service account with minimal permissions
spark.kubernetes.authenticate.driver.serviceAccountName=spark

# Mount secrets as environment variables
spark.kubernetes.driverEnv.DB_PASSWORD=<SECRET_NAME>
```

```bash
# Create K8s secret for DB credentials
kubectl create secret generic spark-db-creds \
  --from-literal=username=spark_user \
  --from-literal=password=secret_pass \
  -n spark-jobs

# Reference in spark-submit
--conf spark.kubernetes.driver.secretKeyRef.DB_USER=spark-db-creds:username
--conf spark.kubernetes.driver.secretKeyRef.DB_PASS=spark-db-creds:password
```

---

## Data Masking / Anonymization

For PII compliance:

```python
from pyspark.sql import functions as F

def mask_pii(df):
    return df \
        .withColumn("email",  F.regexp_replace("email", r"^(.{2}).*(@.*)$", "$1***$2")) \
        .withColumn("phone",  F.regexp_replace("phone", r"\d{7}$", "*******")) \
        .withColumn("ssn",    F.lit("***-**-****")) \
        .withColumn("ip",     F.regexp_replace("ip", r"\d+$", "0"))   # anonymize last octet
```

### Tokenization

```python
# Replace PII with a stable hash for correlation without exposing real values
df = df.withColumn(
    "user_token",
    F.sha2(F.concat_ws("|", F.col("email"), F.lit("SALT_SECRET")), 256)
)
```

---

## Security Checklist

| Area | Action |
|---|---|
| Authentication | Enable Kerberos or shared secret |
| Encryption in transit | Enable `spark.network.crypto` or SSL |
| Encryption at rest | HDFS zones or S3 SSE |
| Credentials | Never hardcode passwords — use Vault or credential provider |
| Spark UI | Enable ACLs; don't expose port 4040 publicly |
| Kubernetes | Minimal RBAC; don't use `cluster-admin` for Spark SA |
| PII | Mask/tokenize before writing to data lake |
| Audit | Enable event logs; retain for compliance window |


---

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


---

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
