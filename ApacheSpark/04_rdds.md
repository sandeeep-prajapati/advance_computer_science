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
