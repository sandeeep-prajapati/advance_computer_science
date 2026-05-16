# Chapter 10: Executors

## What is an Executor?

The executor determines **how and where tasks run**. It receives task instances from the scheduler and dispatches them to workers.

```
Scheduler ──► Executor ──► Worker(s) ──► Run task
```

| Executor | Workers | Use Case |
|---|---|---|
| `SequentialExecutor` | None (in-process) | Testing only |
| `LocalExecutor` | Subprocesses (same machine) | Single-node, small workloads |
| `CeleryExecutor` | Remote worker pool | Multi-node, production |
| `KubernetesExecutor` | Kubernetes pods | Cloud-native, strong isolation |
| `CeleryKubernetesExecutor` | Both Celery + K8s | Hybrid workloads |
| `LocalKubernetesExecutor` | Local + K8s | Local light + K8s heavy |

---

## SequentialExecutor

Runs one task at a time in the scheduler process. Only for local development.

```properties
AIRFLOW__CORE__EXECUTOR=SequentialExecutor
```

---

## LocalExecutor

Spawns subprocesses on the same machine as the scheduler. Simple and fast for single-node setups.

```properties
AIRFLOW__CORE__EXECUTOR=LocalExecutor
AIRFLOW__DATABASE__SQL_ALCHEMY_CONN=postgresql+psycopg2://airflow:airflow@postgres/airflow
```

```yaml
# docker-compose (simplified)
scheduler:
  image: apache/airflow:2.9.0
  command: scheduler
  environment:
    AIRFLOW__CORE__EXECUTOR: LocalExecutor
```

**Max parallelism:** controlled by `AIRFLOW__CORE__PARALLELISM` (default 32).

---

## CeleryExecutor

Distributes tasks across a pool of worker nodes using Celery. Requires a message broker (Redis or RabbitMQ).

```
Scheduler ──► Redis/RabbitMQ ──► Celery Workers (N nodes)
                                       │
                                 run task instance
```

```properties
AIRFLOW__CORE__EXECUTOR=CeleryExecutor
AIRFLOW__CELERY__BROKER_URL=redis://:password@redis:6379/0
AIRFLOW__CELERY__RESULT_BACKEND=db+postgresql://airflow:airflow@postgres/airflow
AIRFLOW__CELERY__WORKER_CONCURRENCY=16  # tasks per worker
```

```yaml
# docker-compose snippet
worker:
  image: apache/airflow:2.9.0
  command: celery worker
  environment:
    AIRFLOW__CORE__EXECUTOR: CeleryExecutor
    AIRFLOW__CELERY__BROKER_URL: redis://redis:6379/0
  deploy:
    replicas: 3  # scale workers
```

### Worker Queues

Route specific tasks to specific workers:

```python
# Task uses "gpu" queue
train_model = PythonOperator(
    task_id = "train_model",
    queue   = "gpu",           # route to GPU worker
    ...
)
```

```bash
# Start a worker that handles the "gpu" queue
airflow celery worker -q gpu,default
```

---

## KubernetesExecutor

Creates a new Kubernetes pod for every task instance. Full isolation, no persistent workers.

```
Scheduler ──► Kubernetes API ──► Pod (task runs) ──► Pod deleted
```

```properties
AIRFLOW__CORE__EXECUTOR=KubernetesExecutor
AIRFLOW__KUBERNETES__NAMESPACE=airflow
AIRFLOW__KUBERNETES__WORKER_CONTAINER_REPOSITORY=apache/airflow
AIRFLOW__KUBERNETES__WORKER_CONTAINER_TAG=2.9.0
AIRFLOW__KUBERNETES__IN_CLUSTER=True
```

### Pod Configuration Per Task

```python
from airflow.providers.cncf.kubernetes.operators.pod import KubernetesPodOperator
from airflow.kubernetes.secret import Secret
from kubernetes.client import models as k8s

# Per-task pod override
task = PythonOperator(
    task_id = "heavy_task",
    executor_config = {
        "pod_override": k8s.V1Pod(
            spec=k8s.V1PodSpec(
                containers=[k8s.V1Container(
                    name      = "base",
                    resources = k8s.V1ResourceRequirements(
                        requests = {"cpu": "4", "memory": "16Gi"},
                        limits   = {"cpu": "8", "memory": "32Gi"},
                    ),
                    env = [
                        k8s.V1EnvVar(name="EXTRA_VAR", value="hello"),
                    ],
                )],
                node_selector = {"cloud.google.com/gke-nodepool": "high-memory"},
            )
        )
    },
)
```

---

## Concurrency Settings

```properties
# Total concurrent tasks across all DAGs
AIRFLOW__CORE__PARALLELISM=32

# Concurrent tasks per DAG
AIRFLOW__CORE__MAX_ACTIVE_TASKS_PER_DAG=16

# Concurrent DAG runs per DAG
AIRFLOW__CORE__MAX_ACTIVE_RUNS_PER_DAG=3

# CeleryExecutor worker concurrency
AIRFLOW__CELERY__WORKER_CONCURRENCY=16
```

```python
# Per-DAG overrides
with DAG(
    dag_id              = "my_dag",
    max_active_runs     = 2,
    max_active_tasks    = 8,
    concurrency         = 8,   # legacy: max tasks across all runs of this DAG
) as dag:
    pass
```

---

## Task Queues and Pools

### Pools

Pools limit the number of concurrently running tasks that use a shared resource (e.g., a slow DB that can't handle more than 5 connections):

```bash
# Create a pool
airflow pools set db_pool 5 "Limit concurrent DB tasks"
airflow pools list
```

```python
task = PythonOperator(
    task_id = "query_db",
    pool    = "db_pool",    # limited to 5 concurrent tasks
    pool_slots = 1,         # each task takes 1 slot (default)
)

# Heavy tasks take more slots
heavy_task = PythonOperator(
    task_id    = "full_table_scan",
    pool       = "db_pool",
    pool_slots = 3,   # this task alone uses 3 of 5 db_pool slots
)
```

### Priority Weight

```python
# Higher priority_weight = runs first when slots are contested
urgent_task = PythonOperator(
    task_id         = "urgent_load",
    priority_weight = 100,    # default is 1
    weight_rule     = "downstream",  # or "upstream", "absolute"
)
```

---

## Choosing an Executor

```
SequentialExecutor:
  ✓ Testing DAG logic
  ✗ Never in production

LocalExecutor:
  ✓ Single machine, moderate workloads
  ✓ Simpler setup
  ✗ No horizontal scaling
  ✗ Worker failure = scheduler failure

CeleryExecutor:
  ✓ Multi-node, horizontal scaling
  ✓ Worker pools + queues
  ✗ Requires Redis/RabbitMQ
  ✗ Workers share dependencies

KubernetesExecutor:
  ✓ Full isolation (each task = new pod)
  ✓ Custom resources per task
  ✓ Scales to zero (no idle workers)
  ✗ Pod startup overhead (~10-30s per task)
  ✗ Complex setup
  ✓ Best for: cloud-native, varying resource needs
```
