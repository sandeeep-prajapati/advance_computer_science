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
