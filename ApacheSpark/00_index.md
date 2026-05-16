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
