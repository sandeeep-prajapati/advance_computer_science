# Lesson 8: Spark Job Examples

## Customer Revenue PySpark Job

```python
import argparse
from pyspark.sql import SparkSession
from pyspark.sql.functions import sum as spark_sum, count

parser = argparse.ArgumentParser()
parser.add_argument("--job-id", required=True)
parser.add_argument("--input", required=True)
parser.add_argument("--output", required=True)
args = parser.parse_args()

spark = SparkSession.builder.appName(f"customer-revenue-{args.job_id}").getOrCreate()

orders = spark.read.option("header", True).csv(args.input)

report = (
    orders
    .groupBy("customer_id")
    .agg(
        spark_sum("total").alias("total_revenue"),
        count("*").alias("order_count"),
    )
)

report.write.mode("overwrite").parquet(args.output)

spark.stop()
```

---

## Read from JDBC

```python
orders = (
    spark.read.format("jdbc")
    .option("url", "jdbc:mysql://mysql:3306/app")
    .option("dbtable", "orders")
    .option("user", "readonly")
    .option("password", "secret")
    .load()
)
```

Use a read replica and restrict credentials.

---

## Write Result Manifest

Spark can write a small JSON manifest for Laravel:

```json
{
  "job_id": 10,
  "status": "completed",
  "output_path": "s3://analytics/reports/customer-revenue/job-10/",
  "row_count": 5000
}
```

Laravel can read this manifest to update the job record.
