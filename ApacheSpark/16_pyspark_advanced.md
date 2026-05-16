# Chapter 16: PySpark Advanced

## Pandas UDFs (Vectorized UDFs)

Regular Python UDFs serialize one row at a time between JVM and Python. **Pandas UDFs** use Apache Arrow to transfer entire batches (Series/DataFrames) — typically 10–100× faster.

```python
import pandas as pd
from pyspark.sql.functions import pandas_udf
from pyspark.sql.types import DoubleType, StringType
```

### Scalar Pandas UDF (Series → Series)

```python
@pandas_udf(returnType=DoubleType())
def normalize_score(scores: pd.Series) -> pd.Series:
    return (scores - scores.mean()) / scores.std()

df.withColumn("norm_score", normalize_score("score")).show()

# Multiple input Series
@pandas_udf(DoubleType())
def weighted_avg(values: pd.Series, weights: pd.Series) -> pd.Series:
    return values * weights / weights.sum()

df.withColumn("weighted", weighted_avg("revenue", "weight")).show()
```

### Iterator of Series (one-time setup cost)

```python
from typing import Iterator

@pandas_udf(StringType())
def classify_batch(iterator: Iterator[pd.Series]) -> Iterator[pd.Series]:
    # Load model once per partition (not per row!)
    import pickle
    model = pickle.load(open("model.pkl", "rb"))

    for series in iterator:
        yield pd.Series(model.predict(series.values.reshape(-1, 1)))

df.withColumn("category", classify_batch("features")).show()
```

### Grouped Map Pandas UDF (applyInPandas)

Apply a function that returns a DataFrame to each group:

```python
def normalize_within_group(pdf: pd.DataFrame) -> pd.DataFrame:
    pdf["norm_score"] = (pdf["score"] - pdf["score"].mean()) / pdf["score"].std()
    return pdf

result = df.groupBy("city").applyInPandas(
    func=normalize_within_group,
    schema="name STRING, city STRING, score DOUBLE, norm_score DOUBLE"
)
```

### Grouped Aggregate Pandas UDF (Series → Scalar)

```python
from pyspark.sql.functions import pandas_udf
from pyspark.sql import functions as F

@pandas_udf(DoubleType())
def median_agg(series: pd.Series) -> float:
    return series.median()

df.groupBy("city").agg(median_agg("score").alias("median_score")).show()
```

---

## pandas API on Spark (formerly Koalas)

Write **Pandas code that runs distributed** on Spark. Ideal for data scientists familiar with Pandas who need to scale.

```python
import pyspark.pandas as ps

# Create pyspark pandas DataFrame
psdf = ps.DataFrame({"a": [1, 2, 3], "b": [4, 5, 6]})

# Or from a Spark DataFrame
psdf = spark_df.to_pandas_on_spark()

# From CSV (distributed)
psdf = ps.read_csv("s3://bucket/data/*.csv")

# Use Pandas-like API
psdf["c"] = psdf["a"] + psdf["b"]
psdf.groupby("category")["revenue"].sum()
psdf.sort_values("score", ascending=False).head(10)
psdf[psdf["age"] > 25].describe()

# Convert back to Spark DataFrame
spark_df = psdf.to_spark()

# Convert to Pandas (driver memory — use for small data only)
pdf = psdf.to_pandas()
```

### Default Index Warning

Pandas-on-Spark uses a distributed index by default. For some operations this requires a shuffle. Use a natural unique column as index:

```python
psdf = psdf.set_index("user_id")
```

---

## Type Hints in PySpark

PySpark 3.1+ supports Python type hints for better IDE support and readability:

```python
from pyspark.sql import DataFrame
from pyspark.sql import functions as F
from pyspark.sql.types import StructType, StructField, StringType, DoubleType

def transform_users(df: DataFrame) -> DataFrame:
    return df \
        .filter(F.col("age") > 18) \
        .withColumn("score_pct", F.col("score") / 100.0)

def get_schema() -> StructType:
    return StructType([
        StructField("name",  StringType(),  nullable=False),
        StructField("score", DoubleType(),  nullable=True),
    ])
```

---

## Arrow Optimization

Apache Arrow speeds up data transfer between JVM (Spark) and Python (pandas):

```python
spark.conf.set("spark.sql.execution.arrow.pyspark.enabled",         "true")
spark.conf.set("spark.sql.execution.arrow.pyspark.fallback.enabled","true")

# This is now much faster with Arrow enabled:
pdf = spark_df.toPandas()          # Spark → Pandas
spark_df = spark.createDataFrame(pdf)  # Pandas → Spark
```

---

## UDF Best Practices

| Approach | Speed | Use Case |
|---|---|---|
| Built-in `F.*` functions | Fastest | Standard operations |
| SQL expressions | Fastest | String-based transformations |
| Pandas UDF (vectorized) | Fast | Custom Python logic, ML inference |
| `applyInPandas` | Medium | Per-group transformations |
| Python UDF (row-by-row) | Slowest | Simple fallback, avoid if possible |

---

## Broadcast Variables

Share large lookup tables across executors without re-sending per task:

```python
# Python dict broadcast
lookup = {"electronics": 0.1, "clothing": 0.2, "books": 0.05}
bc_tax = spark.sparkContext.broadcast(lookup)

@pandas_udf(DoubleType())
def apply_tax(categories: pd.Series, prices: pd.Series) -> pd.Series:
    tax_rates = categories.map(bc_tax.value).fillna(0.0)
    return prices * (1 + tax_rates)

df.withColumn("price_with_tax", apply_tax("category", "price")).show()

# Free memory when done
bc_tax.destroy()
```

---

## Accumulators for Custom Metrics

```python
from pyspark import AccumulatorParam

# Custom accumulator for a list of strings
class ListAccumulator(AccumulatorParam):
    def zero(self, v): return []
    def addInPlace(self, v1, v2): return v1 + v2

errors = spark.sparkContext.accumulator([], ListAccumulator())

def validate_row(row):
    if row["age"] < 0:
        errors.add([f"Invalid age for user {row['user_id']}"])

df.foreach(validate_row)
print(f"Errors: {errors.value[:10]}")
```

---

## Efficient Python Code Patterns

### Reuse SparkSession

```python
# Singleton pattern — don't create multiple sessions
def get_spark():
    return SparkSession.builder.appName("MyApp").getOrCreate()
```

### Avoid `collect()` on Large DataFrames

```python
# Bad: pulls all data to driver
rows = df.collect()
for row in rows:
    process(row)

# Good: process on executors
df.foreach(process)
df.foreachPartition(process_partition)
df.write.parquet("output/")
```

### Use `show()` for debugging, `write` for production

```python
# Debug
df.show(5)
df.printSchema()

# Production
df.write.mode("overwrite").parquet("output/")
```

### Chain Transformations

```python
# Build a clear transformation pipeline
result = (
    spark.read.parquet("input/")
    .filter(F.col("date") >= "2024-01-01")
    .withColumn("year",  F.year("date"))
    .withColumn("month", F.month("date"))
    .groupBy("year", "month", "category")
    .agg(
        F.sum("revenue").alias("total_revenue"),
        F.count("*").alias("order_count"),
    )
    .orderBy("year", "month")
)
result.write.partitionBy("year", "month").parquet("output/")
```
