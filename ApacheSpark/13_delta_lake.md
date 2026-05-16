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
