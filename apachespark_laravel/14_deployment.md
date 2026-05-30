# Lesson 14: Deployment

## Deployment Options

| Option | Good For |
|---|---|
| Spark standalone | Simple self-managed clusters |
| Spark on Kubernetes | Container-native deployments |
| AWS EMR | Managed Spark on AWS |
| Google Dataproc | Managed Spark on GCP |
| Azure Synapse/HDInsight | Managed Spark on Azure |

---

## Laravel Deployment

Laravel needs:

- queue workers for submit jobs
- scheduler for polling or recurring jobs
- access to Spark submit endpoint
- access to object storage
- secure internal callback routes

---

## Spark Job Packaging

Keep Spark jobs versioned:

```text
spark-jobs/
  customer_revenue.py
  product_affinity.py
  requirements.txt
```

Deploy Spark job files with the application or publish them to object storage.

---

## Production Checklist

- [ ] Laravel creates a job row before submitting.
- [ ] Spark input and output paths are deterministic.
- [ ] Spark application ID is stored.
- [ ] Users can see job status.
- [ ] Failed jobs store error messages.
- [ ] Outputs are protected with correct permissions.
- [ ] Large reads use replicas or data lake files.
- [ ] Secrets are not committed.
- [ ] Spark UI/logs are accessible to operators.
- [ ] Backfill process is documented.
