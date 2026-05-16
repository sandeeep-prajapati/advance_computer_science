# Chapter 2: Setup & Installation

## Your Options

| Method | Best For |
|---|---|
| **PySpark (pip)** | Local development, data science notebooks |
| **Docker** | Isolated local cluster, CI pipelines |
| **Databricks Community** | Managed cluster, free tier, notebooks |
| **Google Colab** | Zero-install, notebook-based learning |
| **spark-submit** | Running jobs on a real cluster |
| **AWS EMR / GCP Dataproc** | Managed cloud clusters |

---

## Option A: PySpark Local (Quickest Start)

### Requirements
- Python 3.8+
- Java 8 or 11 (must be on `PATH`)

```bash
# Install Java (Ubuntu/Debian)
sudo apt-get install -y openjdk-11-jdk
export JAVA_HOME=/usr/lib/jvm/java-11-openjdk-amd64

# Install PySpark
pip install pyspark==3.5.0

# Optional: Jupyter for interactive use
pip install jupyter
```

### Verify

```python
# test_spark.py
from pyspark.sql import SparkSession

spark = SparkSession.builder \
    .appName("test") \
    .master("local[*]") \
    .getOrCreate()

spark.range(10).show()
spark.stop()
```

```bash
python test_spark.py
```

---

## Option B: Docker Compose (Local Cluster)

```yaml
# docker-compose.yml
version: '3.8'
services:
  spark-master:
    image: bitnami/spark:3.5
    environment:
      - SPARK_MODE=master
      - SPARK_RPC_AUTHENTICATION_ENABLED=no
    ports:
      - "8080:8080"   # Spark UI
      - "7077:7077"   # Spark master port
    volumes:
      - ./jobs:/opt/spark/jobs

  spark-worker-1:
    image: bitnami/spark:3.5
    environment:
      - SPARK_MODE=worker
      - SPARK_MASTER_URL=spark://spark-master:7077
      - SPARK_WORKER_MEMORY=2G
      - SPARK_WORKER_CORES=2
    depends_on:
      - spark-master

  spark-worker-2:
    image: bitnami/spark:3.5
    environment:
      - SPARK_MODE=worker
      - SPARK_MASTER_URL=spark://spark-master:7077
      - SPARK_WORKER_MEMORY=2G
      - SPARK_WORKER_CORES=2
    depends_on:
      - spark-master
```

```bash
docker-compose up -d

# Submit a job to the cluster
docker exec spark-master spark-submit \
  --master spark://spark-master:7077 \
  /opt/spark/jobs/my_job.py
```

Access Spark UI at: `http://localhost:8080`

---

## Option C: Databricks Community Edition (Free)

1. Sign up at [community.cloud.databricks.com](https://community.cloud.databricks.com)
2. Create a cluster (Runtime: Databricks Runtime 13.x = Spark 3.4)
3. Create a notebook → choose Python
4. SparkSession is pre-created as `spark`

---

## The Spark Shell

### PySpark Shell

```bash
pyspark --master local[4]

# Inside the shell:
>>> spark.version
'3.5.0'
>>> sc = spark.sparkContext
>>> rdd = sc.parallelize([1, 2, 3, 4, 5])
>>> rdd.sum()
15
```

### Scala Spark Shell

```bash
spark-shell --master local[4]

// Inside the shell:
scala> spark.range(5).show()
```

---

## spark-submit: Running Jobs

`spark-submit` is how you run Spark applications on a cluster.

```bash
# Local mode
spark-submit \
  --master local[4] \
  --driver-memory 2g \
  my_job.py

# Standalone cluster
spark-submit \
  --master spark://master:7077 \
  --deploy-mode cluster \
  --executor-memory 4g \
  --executor-cores 2 \
  --num-executors 10 \
  my_job.py

# YARN
spark-submit \
  --master yarn \
  --deploy-mode cluster \
  --executor-memory 4g \
  --num-executors 20 \
  --conf spark.yarn.queue=default \
  hdfs:///jobs/my_job.py arg1 arg2

# Kubernetes
spark-submit \
  --master k8s://https://k8s-master:6443 \
  --deploy-mode cluster \
  --conf spark.kubernetes.container.image=my-spark:3.5 \
  --conf spark.kubernetes.namespace=spark-jobs \
  local:///opt/spark/jobs/my_job.py
```

---

## SparkSession Configuration

```python
from pyspark.sql import SparkSession

spark = SparkSession.builder \
    .appName("ProductionJob") \
    .master("local[*]") \
    \
    # Memory
    .config("spark.driver.memory",         "4g") \
    .config("spark.executor.memory",       "8g") \
    .config("spark.executor.cores",        "4") \
    .config("spark.memory.fraction",       "0.8") \
    \
    # Performance
    .config("spark.sql.shuffle.partitions","200") \
    .config("spark.default.parallelism",   "200") \
    .config("spark.sql.adaptive.enabled",  "true") \
    \
    # Delta Lake / Hive
    .config("spark.sql.extensions",        "io.delta.sql.DeltaSparkSessionExtension") \
    .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
    \
    # Logging
    .config("spark.eventLog.enabled",      "true") \
    .config("spark.eventLog.dir",          "/tmp/spark-events") \
    \
    .getOrCreate()

# Set log level
spark.sparkContext.setLogLevel("WARN")
```

---

## Project Structure

```
my_spark_project/
├── src/
│   ├── jobs/
│   │   ├── etl_job.py
│   │   └── streaming_job.py
│   └── utils/
│       ├── spark_session.py
│       └── schema.py
├── tests/
│   └── test_etl.py
├── data/
│   └── sample.csv
├── requirements.txt
└── submit.sh
```

### Reusable SparkSession Factory

```python
# src/utils/spark_session.py
from pyspark.sql import SparkSession

def get_spark(app_name: str, **configs) -> SparkSession:
    builder = SparkSession.builder.appName(app_name)
    for key, val in configs.items():
        builder = builder.config(key, val)
    return builder.getOrCreate()
```

---

## Jupyter + PySpark

```bash
pip install pyspark findspark jupyter

# Start Jupyter
jupyter notebook
```

```python
# In a notebook cell
import findspark
findspark.init()

from pyspark.sql import SparkSession
spark = SparkSession.builder.appName("notebook").master("local[*]").getOrCreate()
spark.range(5).show()
```

---

## Useful Environment Variables

```bash
export SPARK_HOME=/usr/local/spark
export PYSPARK_PYTHON=python3
export PYSPARK_DRIVER_PYTHON=jupyter
export PYSPARK_DRIVER_PYTHON_OPTS='notebook'
export JAVA_HOME=/usr/lib/jvm/java-11-openjdk-amd64
export PATH=$SPARK_HOME/bin:$JAVA_HOME/bin:$PATH
```
