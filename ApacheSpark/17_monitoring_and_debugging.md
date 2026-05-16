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
