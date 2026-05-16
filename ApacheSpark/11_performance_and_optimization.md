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
