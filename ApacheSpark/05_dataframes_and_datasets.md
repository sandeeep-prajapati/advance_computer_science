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
