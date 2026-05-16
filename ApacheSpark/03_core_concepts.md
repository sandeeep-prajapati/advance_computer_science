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
