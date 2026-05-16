# Chapter 6: Spark SQL

## What is Spark SQL?

**Spark SQL** lets you query DataFrames using standard SQL syntax. It's fully integrated with the DataFrame API — you can mix SQL and DataFrame operations freely.

```python
# Create a temporary view
df.createOrReplaceTempView("users")

# Query with SQL
result = spark.sql("""
    SELECT city, COUNT(*) as total, AVG(score) as avg_score
    FROM users
    WHERE age > 18
    GROUP BY city
    ORDER BY total DESC
    LIMIT 10
""")
result.show()
```

SQL results are DataFrames — you can chain DataFrame operations after a SQL query.

---

## Creating Views

```python
# Temporary view — session-scoped (disappears when session ends)
df.createOrReplaceTempView("users")

# Global temporary view — available across sessions in same app
df.createOrReplaceGlobalTempView("global_users")
spark.sql("SELECT * FROM global_temp.global_users").show()

# Permanent table (persisted to metastore — requires Hive or Glue)
df.write.saveAsTable("default.users")
spark.sql("SELECT * FROM default.users").show()
```

---

## SQL Syntax Reference

### SELECT and Filtering

```sql
-- Basic select
SELECT name, age, score FROM users;

-- Aliases
SELECT name AS full_name, score * 100 AS score_pct FROM users;

-- Distinct
SELECT DISTINCT city FROM users;

-- Filtering
SELECT * FROM users WHERE age BETWEEN 25 AND 45;
SELECT * FROM users WHERE city IN ('London', 'Paris', 'Berlin');
SELECT * FROM users WHERE name LIKE 'A%';
SELECT * FROM users WHERE score IS NOT NULL;

-- Conditional
SELECT name,
       CASE WHEN score >= 90 THEN 'A'
            WHEN score >= 80 THEN 'B'
            ELSE 'C' END AS grade
FROM users;
```

### Aggregations

```sql
SELECT
    city,
    COUNT(*)                AS total_users,
    COUNT(DISTINCT user_id) AS unique_users,
    AVG(score)              AS avg_score,
    MAX(score)              AS max_score,
    MIN(age)                AS min_age,
    SUM(revenue)            AS total_revenue,
    PERCENTILE(score, 0.95) AS p95_score
FROM users
GROUP BY city
HAVING COUNT(*) > 100
ORDER BY total_revenue DESC;
```

### Joins

```sql
-- Inner join
SELECT u.name, o.order_id, o.total
FROM users u
INNER JOIN orders o ON u.user_id = o.user_id;

-- Left join with null handling
SELECT u.name, COALESCE(o.total, 0) AS total
FROM users u
LEFT JOIN orders o ON u.user_id = o.user_id;

-- Cross join
SELECT * FROM dim_dates CROSS JOIN dim_products;
```

---

## Window Functions

Window functions compute values across a sliding window of rows. They're one of Spark SQL's most powerful features.

```python
from pyspark.sql import Window
from pyspark.sql import functions as F

# Define window
w = Window.partitionBy("city").orderBy(F.col("score").desc())

# Rank functions
df.withColumn("rank",        F.rank().over(w))       \
  .withColumn("dense_rank",  F.dense_rank().over(w)) \
  .withColumn("row_number",  F.row_number().over(w)) \
  .withColumn("percent_rank",F.percent_rank().over(w)) \
  .show()

# Lag / Lead
w2 = Window.partitionBy("user_id").orderBy("event_date")
df.withColumn("prev_score", F.lag("score",  1).over(w2)) \
  .withColumn("next_score", F.lead("score", 1).over(w2)) \
  .show()

# Running total / moving average
w3 = Window.partitionBy("user_id").orderBy("event_date") \
           .rowsBetween(Window.unboundedPreceding, Window.currentRow)

df.withColumn("cumulative_sum", F.sum("revenue").over(w3)).show()

# Moving average (last 7 rows)
w4 = Window.partitionBy("user_id").orderBy("event_date") \
           .rowsBetween(-6, 0)

df.withColumn("moving_avg_7d", F.avg("score").over(w4)).show()

# ntile (quantile buckets)
w5 = Window.orderBy("score")
df.withColumn("quartile", F.ntile(4).over(w5)).show()
```

### In SQL

```sql
SELECT
    name,
    city,
    score,
    RANK()  OVER (PARTITION BY city ORDER BY score DESC) AS city_rank,
    SUM(revenue) OVER (PARTITION BY city ORDER BY event_date
                       ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS running_total,
    LAG(score, 1) OVER (PARTITION BY user_id ORDER BY event_date) AS prev_score
FROM events;
```

---

## CTEs (WITH Clause)

```sql
WITH active_users AS (
    SELECT user_id, name, city
    FROM users
    WHERE status = 'active'
),
user_orders AS (
    SELECT user_id, SUM(total) AS lifetime_value
    FROM orders
    GROUP BY user_id
)
SELECT u.name, u.city, COALESCE(uo.lifetime_value, 0) AS ltv
FROM active_users u
LEFT JOIN user_orders uo USING (user_id)
ORDER BY ltv DESC;
```

---

## Explain Plan

Always check the execution plan for complex queries:

```python
df.explain()              # physical plan only
df.explain("formatted")   # clean readable format
df.explain("cost")        # with cost estimates
df.explain("codegen")     # JVM code generation
```

Example output analysis:
```
== Physical Plan ==
AdaptiveSparkPlan
+- SortMergeJoin [user_id]         ← expensive shuffle join
   :- Sort [user_id]
   :  +- Exchange hashpartitioning  ← shuffle!
   :     +- Filter (age > 18)       ← pushed down ✓
   :        +- Scan parquet users
   +- Sort [user_id]
      +- Exchange hashpartitioning  ← shuffle!
         +- Scan parquet orders
```

Look for:
- `BroadcastHashJoin` → good (no shuffle for small table)
- `SortMergeJoin` → shuffle required (make sure it's necessary)
- `Exchange` → shuffle happening
- `Filter` pushed down before `Scan` → good (predicate pushdown)

---

## JDBC: Reading from Databases

```python
# Read from PostgreSQL
df = spark.read \
    .format("jdbc") \
    .option("url",      "jdbc:postgresql://host:5432/mydb") \
    .option("dbtable",  "public.orders") \
    .option("user",     "spark_user") \
    .option("password", "secret") \
    .option("driver",   "org.postgresql.Driver") \
    .load()

# Parallel read — partition by a numeric column
df = spark.read \
    .format("jdbc") \
    .option("url",           "jdbc:postgresql://host:5432/mydb") \
    .option("dbtable",       "public.orders") \
    .option("partitionColumn","order_id") \
    .option("lowerBound",    "1") \
    .option("upperBound",    "10000000") \
    .option("numPartitions", "50") \
    .load()
```

---

## The Spark Catalog

```python
# List databases
spark.catalog.listDatabases()

# List tables
spark.catalog.listTables("default")

# List columns
spark.catalog.listColumns("users", "default")

# Check if table exists
spark.catalog.tableExists("default.users")

# Refresh table metadata (after external changes)
spark.catalog.refreshTable("default.users")

# Clear cache
spark.catalog.clearCache()
```

---

## Built-In SQL Functions Reference

```python
import pyspark.sql.functions as F

# String
F.upper("col"), F.lower("col"), F.length("col")
F.trim("col"), F.ltrim("col"), F.rtrim("col")
F.split("col", ","), F.concat_ws(",", "a", "b")
F.regexp_replace("col", r"\d+", "")
F.regexp_extract("col", r"(\d{4})", 1)
F.substring("col", 1, 3)
F.lpad("col", 10, "0"), F.rpad("col", 10, " ")
F.translate("col", "abc", "ABC")
F.instr("col", "sub")        # position of substring

# Numeric
F.abs("col"), F.round("col", 2), F.ceil("col"), F.floor("col")
F.sqrt("col"), F.pow("col", 2), F.log("col"), F.log2("col")
F.greatest("a", "b", "c"),  F.least("a", "b", "c")
F.rand(), F.randn()          # random numbers

# Date / Timestamp
F.current_date(), F.current_timestamp()
F.to_date("col", "yyyy-MM-dd")
F.to_timestamp("col", "yyyy-MM-dd HH:mm:ss")
F.date_add("col", 7), F.date_sub("col", 7)
F.datediff("end_date", "start_date")
F.months_between("col1", "col2")
F.date_format("col", "yyyy-MM")
F.year("col"), F.month("col"), F.dayofmonth("col")
F.hour("col"), F.minute("col"), F.second("col")
F.dayofweek("col"), F.weekofyear("col")
F.trunc("col", "month"), F.date_trunc("hour", "col")

# Conditional
F.when(F.col("age") > 18, "adult").otherwise("minor")
F.coalesce("a", "b", F.lit("default"))
F.isnull("col"), F.isnan("col")
F.nullif("col", F.lit(0))

# Hash / Encoding
F.md5("col"), F.sha2("col", 256)
F.hash("col"), F.xxhash64("col")
F.base64("col"), F.unbase64("col")
```
