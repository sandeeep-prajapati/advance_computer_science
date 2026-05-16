# Chapter 15: Spark on YARN & EMR

## Spark on YARN

YARN (Yet Another Resource Negotiator) is the resource manager in the Hadoop ecosystem. It's the most common cluster manager for Spark in enterprise environments (Cloudera, Hortonworks, Amazon EMR).

```
YARN Architecture:
  ResourceManager (master)  ← allocates containers
  NodeManagers (workers)    ← run containers
  ApplicationMaster         ← Spark Driver (in cluster mode)
```

---

## Deploy Modes

### Client Mode

```bash
spark-submit \
  --master yarn \
  --deploy-mode client \
  --executor-memory 4g \
  --executor-cores 2 \
  --num-executors 10 \
  --driver-memory 4g \
  my_job.py
```

- **Driver** runs on the machine that ran `spark-submit`
- Good for: interactive use, debugging (logs to your terminal)
- Bad for: production (if your laptop disconnects, job dies)

### Cluster Mode

```bash
spark-submit \
  --master yarn \
  --deploy-mode cluster \
  --executor-memory 4g \
  --executor-cores 2 \
  --num-executors 20 \
  --driver-memory 4g \
  --conf spark.yarn.queue=production \
  hdfs:///jobs/my_job.py
```

- **Driver** runs on a YARN worker node (inside YARN)
- Good for: production (job survives disconnect)
- Logs go to YARN aggregated logs

---

## YARN-Specific Configuration

```bash
# Queue management
--conf spark.yarn.queue=high_priority

# AM memory overhead
--conf spark.yarn.am.memory=2g
--conf spark.yarn.am.memoryOverhead=512m

# Executor overhead (off-heap: PySpark, overhead for libraries)
--conf spark.yarn.executor.memoryOverhead=1024m

# Max app attempts
--conf spark.yarn.maxAppAttempts=2

# Archive Python deps to YARN
--conf spark.yarn.dist.archives=my_venv.tar.gz#env
--conf spark.pyspark.python=./env/bin/python

# Ship files to all nodes
--files log4j.properties
--jars external-lib.jar
```

---

## YARN UI and Log Access

```bash
# YARN ResourceManager UI
http://yarn-rm:8088

# Find app ID
yarn application -list
yarn application -status application_1234567890_0001

# View logs (after job completes)
yarn logs -applicationId application_1234567890_0001

# Kill a job
yarn application -kill application_1234567890_0001
```

---

## Dynamic Allocation on YARN

```python
spark.conf.set("spark.dynamicAllocation.enabled",                "true")
spark.conf.set("spark.shuffle.service.enabled",                  "true")   # required for YARN DRA
spark.conf.set("spark.dynamicAllocation.minExecutors",           "2")
spark.conf.set("spark.dynamicAllocation.maxExecutors",           "100")
spark.conf.set("spark.dynamicAllocation.executorIdleTimeout",    "60s")
spark.conf.set("spark.dynamicAllocation.schedulerBacklogTimeout","1s")
```

The **External Shuffle Service** (`spark.shuffle.service.enabled=true`) runs on each YARN NodeManager and holds shuffle data when executors are released.

---

## AWS EMR

**Amazon EMR** (Elastic MapReduce) is AWS's managed Hadoop/Spark service. It provisions clusters on EC2, runs your job, and optionally terminates.

### Launch EMR Cluster (AWS CLI)

```bash
aws emr create-cluster \
  --name "SparkProductionCluster" \
  --release-label emr-6.15.0 \
  --applications Name=Spark Name=Hadoop Name=Hive \
  --instance-type m5.2xlarge \
  --instance-count 5 \
  --ec2-attributes KeyName=my-key-pair \
  --use-default-roles \
  --log-uri s3://my-bucket/emr-logs/ \
  --configurations file://spark-config.json
```

### Spark Configuration on EMR (`spark-config.json`)

```json
[
  {
    "Classification": "spark-defaults",
    "Properties": {
      "spark.executor.memory":           "6g",
      "spark.executor.cores":            "4",
      "spark.sql.shuffle.partitions":    "400",
      "spark.sql.adaptive.enabled":      "true",
      "spark.dynamicAllocation.enabled": "true"
    }
  },
  {
    "Classification": "spark",
    "Properties": {
      "maximizeResourceAllocation": "true"
    }
  }
]
```

### Submit a Job to EMR

```bash
# Add a step to a running cluster
aws emr add-steps \
  --cluster-id j-XXXXXXXXXXXX \
  --steps Type=spark,\
Name=MyJob,\
Args=[--deploy-mode,cluster,\
      --master,yarn,\
      s3://my-bucket/jobs/my_job.py,\
      --input,s3://my-bucket/input/,\
      --output,s3://my-bucket/output/],\
ActionOnFailure=CONTINUE
```

### EMR with S3A Committer (Faster S3 Writes)

```json
{
  "Classification": "spark-defaults",
  "Properties": {
    "spark.hadoop.fs.s3a.committer.name":               "directory",
    "spark.sql.sources.commitProtocolClass":             "org.apache.spark.internal.io.cloud.PathOutputCommitProtocol",
    "spark.sql.parquet.output.committer.class":          "org.apache.spark.internal.io.cloud.BindingParquetOutputCommitter",
    "spark.hadoop.mapreduce.outputcommitter.factory.scheme.s3a": "org.apache.hadoop.fs.s3a.commit.S3ACommitterFactory"
  }
}
```

---

## EMR Serverless (No Cluster Management)

Submit Spark jobs without managing a cluster:

```bash
# Create an application
aws emr-serverless create-application \
  --name "MySparkApp" \
  --type SPARK \
  --release-label emr-6.15.0

# Submit a job
aws emr-serverless start-job-run \
  --application-id APP_ID \
  --execution-role-arn arn:aws:iam::123456789:role/EMRServerlessRole \
  --job-driver '{
    "sparkSubmit": {
      "entryPoint": "s3://bucket/jobs/my_job.py",
      "sparkSubmitParameters": "--conf spark.executor.memory=4g"
    }
  }' \
  --configuration-overrides '{
    "monitoringConfiguration": {
      "s3MonitoringConfiguration": {
        "logUri": "s3://bucket/emr-logs/"
      }
    }
  }'
```

---

## Databricks (Cloud-Managed Spark)

Databricks is the company behind Spark — their platform provides:
- Fully managed clusters (auto-scaling, auto-termination)
- Collaborative notebooks
- Delta Lake as first-class citizen
- MLflow for ML lifecycle
- Unity Catalog for governance

```python
# Databricks Runtime: SparkSession is pre-created as `spark`

# Read from Databricks Unity Catalog
df = spark.table("my_catalog.my_schema.my_table")

# Read from DBFS or cloud storage
df = spark.read.parquet("dbfs:/mnt/my-bucket/data/")
df = spark.read.parquet("s3://my-bucket/data/")

# Write as Delta (default in Databricks)
df.write.format("delta").saveAsTable("my_catalog.my_schema.output_table")
```
