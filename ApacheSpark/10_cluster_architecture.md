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
