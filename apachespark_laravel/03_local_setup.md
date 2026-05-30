# Lesson 3: Local Setup

## Local Services

For development, use:

```text
laravel-app
mysql or postgres
spark-master
spark-worker
minio or local mounted data folder
```

Optional:

```text
livy
airflow
kafka
```

---

## Laravel Environment

```dotenv
SPARK_MASTER=spark://spark-master:7077
SPARK_SUBMIT=/opt/spark/bin/spark-submit
SPARK_JOBS_DIR=/opt/spark/jobs
SPARK_OUTPUT_DISK=s3
SPARK_OUTPUT_BUCKET=analytics
```

---

## Docker Network Rule

Inside containers, use service names:

```text
spark-master:7077
mysql:3306
minio:9000
laravel-app:80
```

Do not use `localhost` unless the service is inside the same container.

---

## Smoke Test Spark

From a Spark container:

```bash
spark-submit --master spark://spark-master:7077 /opt/spark/examples/src/main/python/pi.py 10
```

From Laravel:

```bash
php artisan spark:test
```

The Laravel command can submit a tiny Spark job and store the result in `spark_jobs`.
