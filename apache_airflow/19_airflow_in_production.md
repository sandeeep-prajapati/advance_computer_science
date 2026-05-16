# Chapter 19: Airflow in Production

## Managed Airflow Options

| Service | Provider | Executor | Best For |
|---|---|---|---|
| MWAA | AWS | CeleryExecutor | AWS-native teams |
| Cloud Composer | GCP | KubernetesExecutor | GCP-native, K8s teams |
| Astronomer | Multi-cloud | CeleryExecutor or K8s | Enterprise, support |
| Astro (SaaS) | Astronomer Cloud | KubernetesExecutor | Fully managed |
| Self-hosted K8s | Any | KubernetesExecutor | Full control |

---

## AWS MWAA (Managed Workflows for Apache Airflow)

```python
# Create MWAA environment via CloudFormation / Terraform
# Resource: aws_mwaa_environment

resource "aws_mwaa_environment" "airflow" {
  name              = "production-airflow"
  airflow_version   = "2.9.1"
  environment_class = "mw1.medium"   # mw1.small | mw1.medium | mw1.large

  source_bucket_arn = aws_s3_bucket.airflow.arn
  dag_s3_path       = "dags/"
  requirements_s3_path = "requirements.txt"
  plugins_s3_path   = "plugins.zip"

  airflow_configuration_options = {
    "core.default_timezone"          = "UTC"
    "webserver.dag_default_view"     = "grid"
    "scheduler.min_file_process_interval" = "60"
  }

  network_configuration {
    security_group_ids = [aws_security_group.airflow.id]
    subnet_ids         = aws_subnet.private[*].id
  }

  logging_configuration {
    dag_processing_logs { enabled = true; log_level = "WARNING" }
    scheduler_logs      { enabled = true; log_level = "WARNING" }
    task_logs           { enabled = true; log_level = "INFO" }
    webserver_logs      { enabled = true; log_level = "WARNING" }
    worker_logs         { enabled = true; log_level = "WARNING" }
  }
}
```

```bash
# Deploy DAGs to MWAA (sync to S3)
aws s3 sync ./dags s3://my-mwaa-bucket/dags/ --delete
aws s3 cp requirements.txt s3://my-mwaa-bucket/requirements.txt
aws s3 cp plugins.zip s3://my-mwaa-bucket/plugins.zip
```

---

## GCP Cloud Composer

```bash
# Create a Cloud Composer environment
gcloud composer environments create production-airflow \
  --location us-central1 \
  --image-version composer-3-airflow-2.9.1 \
  --machine-type n1-standard-2 \
  --node-count 3 \
  --disk-size 50 \
  --python-version 3

# Deploy DAGs
gcloud composer environments storage dags import \
  --environment production-airflow \
  --location us-central1 \
  --source ./dags/

# Install Python packages
gcloud composer environments update production-airflow \
  --location us-central1 \
  --update-pypi-packages-from-file requirements.txt

# Set Airflow variables
gcloud composer environments run production-airflow \
  --location us-central1 variables -- \
  --set environment production
```

---

## Self-Hosted with Helm (Kubernetes)

```bash
# Add the official Helm chart repo
helm repo add apache-airflow https://airflow.apache.org
helm repo update

# Create values.yaml and install
helm upgrade --install airflow apache-airflow/airflow \
  --namespace airflow \
  --create-namespace \
  -f values.yaml
```

```yaml
# values.yaml (key sections)
airflowVersion: "2.9.1"
executor: KubernetesExecutor

images:
  airflow:
    repository: my-registry/airflow
    tag: "2.9.1-custom"

config:
  core:
    parallelism: "64"
    max_active_tasks_per_dag: "32"
  webserver:
    base_url: "https://airflow.mycompany.com"

dags:
  gitSync:
    enabled: true
    repo: https://github.com/mycompany/airflow-dags
    branch: main
    subPath: dags/
    credentialsSecret: airflow-git-credentials

workers:
  replicas: 0   # KubernetesExecutor creates pods on demand

webserver:
  replicas: 2
  service:
    type: ClusterIP

scheduler:
  replicas: 1

postgresql:
  enabled: false  # use external PostgreSQL

externalDatabase:
  type: postgres
  host: my-postgres.example.com
  port: 5432
  database: airflow
  user: airflow
  passwordSecret: airflow-db-credentials
  passwordSecretKey: password
```

---

## Production Architecture

```
                           ┌─────────────────────────────────┐
                           │           Kubernetes              │
                           │                                   │
  Users ──► ALB ──► ┌──────────┐   ┌──────────┐              │
                    │Webserver │   │Scheduler │              │
                    └──────────┘   └─────┬────┘              │
                         │               │ creates pods        │
                         │         ┌─────▼──────────────┐     │
                         │         │  Worker Pods (K8s)  │     │
                         │         └────────────────────-┘     │
                         │                                     │
                    ┌────▼─────────────────────────────┐      │
                    │         PostgreSQL (RDS)          │      │
                    └──────────────────────────────────┘      │
                                                               │
  Secrets ◄──── AWS Secrets Manager / Vault                   │
  DAGs    ◄──── Git Sync (github.com/org/dags)                │
  Logs    ────► S3 / GCS                                       │
  Metrics ────► Prometheus / Grafana                          │
                           └─────────────────────────────────┘
```

---

## DAG Deployment Strategy

### Git Sync (Recommended)

```yaml
# Kubernetes sidecar that continuously syncs DAGs from Git
dags:
  gitSync:
    enabled:           true
    repo:              https://github.com/mycompany/airflow-dags
    branch:            main
    rev:               HEAD
    depth:             1
    maxFailures:       3
    subPath:           dags/
    wait:              60   # sync every 60 seconds
    containerName:     git-sync
    uid:               65533
    credentialsSecret: airflow-git-credentials   # K8s secret with SSH key or token
```

### CI/CD Pipeline

```yaml
# .github/workflows/deploy-dags.yml
name: Deploy DAGs

on:
  push:
    branches: [main]
    paths: [dags/**, plugins/**, requirements.txt]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: pip install apache-airflow==2.9.0 pytest
      - run: airflow db init && pytest tests/

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      # MWAA
      - name: Sync to S3
        run: aws s3 sync ./dags s3://my-airflow-bucket/dags/ --delete

      # OR: Cloud Composer
      - name: Deploy to Composer
        run: |
          gcloud composer environments storage dags import \
            --environment production --location us-central1 \
            --source ./dags/
```

---

## High Availability

```yaml
# Multiple schedulers (Airflow 2.0+ supports HA scheduler)
scheduler:
  replicas: 2   # Both schedulers use optimistic row locking — safe

# Multiple webservers behind a load balancer
webserver:
  replicas: 2

# Shared AIRFLOW__WEBSERVER__SECRET_KEY — must match across replicas
# Shared Fernet key — must match across all components
# Shared metadata DB — PostgreSQL with connection pooling (PgBouncer)
# Shared DAG storage — Git sync or shared volume
```

---

## Production Runbook

### Health Checks

```bash
# Check scheduler is alive
airflow jobs check --job-type SchedulerJob --limit 1

# Check DAG parsing errors
airflow dags report | grep -i error

# Check stuck tasks (running > 2 hours)
airflow tasks states-for-dag-run my_dag 2024-01-15T00:00:00

# Check executor queue depth
airflow celery flower   # Celery
kubectl get pods -n airflow | grep "airflow-worker"   # K8s
```

### Incident Responses

```bash
# Clear and re-run a failed task
airflow tasks clear my_dag --task-ids extract --start-date 2024-01-15 --yes

# Mark a task as success (skip it)
airflow tasks state my_dag extract 2024-01-15
airflow tasks mark-success my_dag extract 2024-01-15

# Unpause a DAG
airflow dags unpause my_dag

# Kill a stuck running task instance
airflow tasks failed-deps my_dag extract 2024-01-15

# Force-fail a running DAG run
airflow dags state my_dag 2024-01-15T00:00:00

# Purge old metadata (keep 90 days)
airflow db clean --clean-before-timestamp "$(date -d '-90 days' --iso-8601=seconds)"
```

---

## Production Checklist

```
Infrastructure:
  ☐ PostgreSQL (not SQLite/MySQL)
  ☐ PgBouncer for connection pooling
  ☐ Remote log storage (S3, GCS)
  ☐ External secrets backend (Vault, AWS SM)
  ☐ TLS on webserver
  ☐ Webserver behind reverse proxy

DAGs:
  ☐ catchup=False on all production DAGs
  ☐ execution_timeout on every task
  ☐ retries + retry_delay configured
  ☐ on_failure_callback with alerting
  ☐ No module-level Variable.get() calls
  ☐ SLAs set for critical pipelines

Security:
  ☐ Fernet key set and backed up
  ☐ Default admin credentials changed
  ☐ RBAC configured
  ☐ API authentication enabled
  ☐ DAG-level access_control for sensitive DAGs

Observability:
  ☐ Prometheus metrics scraping
  ☐ Grafana dashboard
  ☐ Alerting on scheduler_heartbeat absence
  ☐ Alerting on task failure rate
  ☐ DAG SLA miss notifications

Operations:
  ☐ DB cleanup cron job (airflow db clean)
  ☐ Backup PostgreSQL metadata DB
  ☐ DAG deployment via CI/CD (not manual copy)
  ☐ Staging environment mirrors production
```
