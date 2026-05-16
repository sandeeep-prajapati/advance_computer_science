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
