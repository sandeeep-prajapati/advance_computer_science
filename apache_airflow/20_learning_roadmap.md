# Chapter 20: Learning Roadmap

## 5-Week Study Plan

### Week 1: Foundations

| Day | Topic | Practice |
|---|---|---|
| Mon | Ch 01–02: What is Airflow, Setup | Run `docker-compose up`, open UI |
| Tue | Ch 03: Core concepts | Draw the DAG → TaskInstance lifecycle |
| Wed | Ch 04: Writing DAGs | Write a 3-task pipeline (extract → transform → load) |
| Thu | Ch 05: Operators | Replace PythonOperator with BashOperator + SparkSubmitOperator |
| Fri | Ch 06–07: Sensors + Hooks | Add an S3KeySensor; use PostgresHook directly |
| Sat | Review | Refactor your pipeline to use TaskFlow API |
| Sun | Rest / catch-up | — |

**Week 1 Checkpoint:** You have a working DAG with sensors, operators, and a hook, deployed via docker-compose.

---

### Week 2: Data Flow & Scheduling

| Day | Topic | Practice |
|---|---|---|
| Mon | Ch 08: XCom & TaskFlow | Pass data between 3 tasks using @task decorators |
| Tue | Ch 09: Scheduling | Change schedule: cron → timedelta → Dataset-driven |
| Wed | Ch 09: Backfill | Run a backfill for 7 days; observe DAG runs in UI |
| Thu | Ch 10: Executors | Switch LocalExecutor → CeleryExecutor with Redis |
| Fri | Ch 11: Dynamic DAGs | Convert static 3-task DAG to dynamic mapping on a list |
| Sat | Ch 13: Variables + Params | Add env variable; trigger DAG with --conf params |
| Sun | Review | Write 5 unit tests for your task callables |

**Week 2 Checkpoint:** You understand scheduling semantics (data_interval, catchup) and have a dynamically-mapped pipeline.

---

### Week 3: Reliability & Testing

| Day | Topic | Practice |
|---|---|---|
| Mon | Ch 12: Testing DAGs | Add DAG validation + task logic tests; run in CI |
| Tue | Ch 14: Monitoring | Add Slack on_failure_callback to your DAG |
| Wed | Ch 14: SLAs | Add an SLA (1 hour) and trigger a miss |
| Thu | Ch 15: Security | Set Fernet key; use env var instead of UI connection |
| Fri | Ch 16: Performance | Measure DAG parse time; profile a slow import |
| Sat | Review | Create a pool; add priority_weight to critical tasks |
| Sun | Rest / catch-up | — |

**Week 3 Checkpoint:** Your pipeline has tests, alerting, and doesn't have import errors or slow parse time.

---

### Week 4: Cloud & Custom Extensions

| Day | Topic | Practice |
|---|---|---|
| Mon | Ch 17: AWS operators | Replace PythonOperator with S3KeySensor + S3Hook |
| Tue | Ch 17: BigQuery/GCP | Write a BigQueryInsertJobOperator task |
| Wed | Ch 18: Custom operators | Build S3ToPostgresOperator; use it in a DAG |
| Thu | Ch 18: Plugins | Wrap operator in an AirflowPlugin |
| Fri | Ch 10: KubernetesExecutor | Deploy to local K8s (minikube); configure pod_override |
| Sat | Ch 19: Production setup | Deploy with Helm on minikube; configure GitSync |
| Sun | Review | Add monitoring: Prometheus + Grafana dashboard |

**Week 4 Checkpoint:** You can write custom operators, deploy to Kubernetes, and have cloud provider integration.

---

### Week 5: Capstone Projects

Pick one capstone to build end-to-end:

---

#### Capstone A: E-Commerce Data Platform

```
Pipeline architecture:

[S3: raw orders]
      │
      ▼ S3KeySensor (deferrable)
[extract_orders] → PythonOperator: validate + dedupe
      │
      ▼
[transform] → PythonOperator: compute revenue, AOV
      │
      ├──► [load_to_postgres] → PostgresHook
      │
      └──► [load_to_redshift] → RedshiftSQLOperator
                │
                ▼
      [notify_slack] → SlackWebhookHook (on_success_callback)

Features to implement:
  ☐ Dynamic task mapping per product category
  ☐ ExternalTaskSensor waiting for inventory DAG
  ☐ Pool limiting Redshift concurrency
  ☐ SLA of 2 hours from midnight UTC
  ☐ Retry with exponential backoff
  ☐ Backfill last 30 days on first deploy
```

---

#### Capstone B: ML Training Pipeline

```
Pipeline architecture:

[Dataset: feature_store updated]
      │  (Dataset-driven schedule)
      ▼
[validate_features] → SqlSensor: count > 1000
      │
      ▼
[prepare_training_data] → dynamic mapping per model_type
      │                    (model_type: ["lr", "xgb", "nn"])
      ▼
[train_model[0/1/2]] → KubernetesPodOperator per model
      │                  (GPU node pool, 8 CPU, 32 GB RAM each)
      ▼
[evaluate_models] → compare metrics, pick best
      │
      ▼
[register_best_model] → push to MLflow / SageMaker
      │
      ▼
[notify_team] → Slack + email

Features to implement:
  ☐ Dataset-driven trigger from feature engineering DAG
  ☐ Dynamic mapping for parallel model training
  ☐ KubernetesPodOperator with GPU node selector
  ☐ Custom S3XComBackend for model artifacts
  ☐ Branching: only deploy if accuracy > 0.90
```

---

#### Capstone C: Real-Time Ingestion Monitor

```
Pipeline architecture (hourly schedule):

[extract_kafka_lag] → PythonOperator (confluent-kafka AdminClient)
      │
      ▼
[check_lag_thresholds] → BranchPythonOperator
      │
      ├── [lag_ok]      → EmptyOperator
      │
      └── [lag_critical] → TriggerDagRunOperator → scale_consumers_dag
                │
                ▼
      [alert_oncall] → PagerDuty / Slack

Features to implement:
  ☐ Hourly schedule with catchup=False
  ☐ Branching on consumer lag threshold
  ☐ TriggerDagRunOperator with conf payload
  ☐ Variable-driven thresholds (lag_warning_threshold)
  ☐ Audit trail: write task outcomes to PostgreSQL
```

---

## Commands Reference Card

```bash
# Start Airflow
docker-compose up -d

# DAG management
airflow dags list
airflow dags trigger my_dag --conf '{"key": "value"}'
airflow dags pause / unpause my_dag
airflow dags delete my_dag --yes

# Task management
airflow tasks list my_dag
airflow tasks test my_dag task_id 2024-01-15         # test without DB record
airflow tasks run  my_dag task_id 2024-01-15         # run and record in DB
airflow tasks clear my_dag --task-ids task_id --yes

# Scheduler
airflow scheduler
airflow jobs check --job-type SchedulerJob

# Celery
airflow celery worker -q default,heavy
airflow celery flower                                 # web UI for workers

# Database
airflow db init
airflow db upgrade
airflow db clean --clean-before-timestamp 2024-01-01 --yes

# Variables & Connections
airflow variables set key value
airflow variables get key
airflow connections add --conn-id id --conn-type type ...
airflow connections list

# Users
airflow users create --role Admin --username admin --email a@b.com --password pass
airflow users list
```

---

## Beyond Airflow

| What you'll learn next | Technology |
|---|---|
| Streaming pipelines | Apache Kafka + Flink |
| Lakehouse patterns | Apache Iceberg, Delta Lake |
| Feature engineering | Feast, Tecton |
| ML orchestration | MLflow, Kubeflow Pipelines |
| Modern orchestration | Prefect, Dagster |
| Data contracts | dbt, Great Expectations |
| Infrastructure as code | Terraform + Helm |

---

## Key Takeaways

```
1. DAGs are code — treat them like production software
   (tests, CI/CD, code review, linting)

2. Use TaskFlow API (@dag, @task) for new pipelines
   — cleaner, auto-XCom, less boilerplate

3. Executors define WHERE tasks run:
   Local = simplicity, Celery = scale out, K8s = isolation

4. Scheduling semantics matter:
   data_interval_start is what you process, not when you run

5. Sensors should use mode="reschedule" or deferrable=True
   for long waits

6. Never store large data in XCom — store in S3/GCS, pass the path

7. Use Secrets Backends for credentials — never hardcode or
   use plain Airflow Variables for sensitive data

8. Measure DAG parse time — module-level code runs constantly

9. The scheduler is a single point of failure in most setups —
   monitor its heartbeat aggressively

10. Production Airflow = PostgreSQL + TLS + Fernet + CI/CD + Alerting
```
