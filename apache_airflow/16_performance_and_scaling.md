# Chapter 16: Performance & Scaling

## Understanding the Bottlenecks

```
Parsing bottleneck:
  Scheduler scans DAG files every min_file_process_interval seconds
  Slow imports or heavy module-level code slows ALL DAGs

Scheduling bottleneck:
  Scheduler runs the critical section loop every scheduler_heartbeat_sec
  Thousands of tasks in the ready queue slow scheduling

Executor bottleneck:
  Limited worker slots → tasks queue up
  CeleryExecutor: add workers
  KubernetesExecutor: add node pool capacity

Database bottleneck:
  Too many tasks writing metadata simultaneously
  Use connection pooling (PgBouncer for PostgreSQL)
```

---

## Scheduler Tuning

```properties
# airflow.cfg

[scheduler]
# Seconds between scheduling loops (default 5)
scheduler_heartbeat_sec = 5

# How often the scheduler scans for new/changed DAG files (default 30)
min_file_process_interval = 30

# Number of parallel threads for DAG parsing (default 2)
parsing_processes = 4

# Max DAG runs created per scheduler loop iteration (default 100)
max_dagruns_per_loop_to_schedule = 100

# Max TIs that can be scheduled per scheduler loop (default 512)
max_tis_per_query = 512

# Seconds to wait for task to heartbeat before marking zombie (default 300)
scheduler_zombie_task_threshold = 300

# Run health check every N seconds
scheduler_health_check_threshold = 30
```

---

## DAG Parsing Performance

```python
# BAD: heavy import at module level — runs on every scheduler heartbeat
import pandas as pd
import tensorflow as tf                # ← slow import
from my_utils import load_big_config   # ← reads a file

GLOBAL_CONFIG = load_big_config()      # ← executed at parse time

@dag(...)
def my_dag():
    @task
    def train():
        import tensorflow as tf        # ← import inside task instead
        ...
```

```python
# GOOD: lazy imports, no work at parse time
from airflow.decorators import dag, task
from datetime import datetime

@dag(start_date=datetime(2024, 1, 1), schedule="@daily", catchup=False)
def fast_dag():

    @task
    def extract() -> list:
        import pandas as pd        # lazy import
        return []

    @task
    def transform(data: list) -> list:
        from my_utils import process   # lazy import
        return process(data)

    transform(extract())

fast_dag()
```

### Measuring Parse Time

```bash
# Time how long it takes to parse all DAGs
airflow dags report

# Parse a specific file and print time
time python /opt/airflow/dags/my_dag.py

# Profile DAG parsing
python -c "
import cProfile
import airflow.models.dagbag
cProfile.run('airflow.models.dagbag.DagBag()')
"
```

---

## Concurrency Configuration

```properties
# airflow.cfg

[core]
# Total tasks running simultaneously across ALL DAGs (default 32)
parallelism = 64

# Max tasks per single DAG across all active runs (default 16)
max_active_tasks_per_dag = 32

# Max concurrent DAG runs per DAG (default 16)
max_active_runs_per_dag = 5

[celery]
# Tasks each Celery worker can run simultaneously (default 16)
worker_concurrency = 32
```

```python
# Per-DAG concurrency limits
with DAG(
    dag_id          = "my_dag",
    max_active_runs = 2,       # no more than 2 runs at the same time
    max_active_tasks = 16,     # no more than 16 tasks running across all runs
    concurrency     = 16,      # legacy alias for max_active_tasks
) as dag:
    pass
```

---

## Pools for Resource Management

```bash
# Create pools
airflow pools set default_pool      128  "Default pool"
airflow pools set db_pool             5  "Limit concurrent DB queries"
airflow pools set api_pool           10  "Limit concurrent API calls"
airflow pools set heavy_compute_pool  3  "GPU/heavy jobs"
```

```python
from airflow.operators.python import PythonOperator

# Light task — no pool needed (uses default_pool)
light_task = PythonOperator(
    task_id         = "light",
    python_callable = lambda: None,
)

# DB-heavy task — limit concurrent access
db_task = PythonOperator(
    task_id         = "query_warehouse",
    python_callable = run_query,
    pool            = "db_pool",
    pool_slots      = 1,            # uses 1 of 5 db_pool slots
)

# Very heavy task — takes 3 slots
heavy_task = PythonOperator(
    task_id         = "full_table_scan",
    python_callable = run_heavy,
    pool            = "db_pool",
    pool_slots      = 3,            # blocks 3 other tasks in db_pool
)
```

---

## Priority Weights

```python
# Higher priority_weight = scheduled first when slots are contested
from airflow.utils.weight_rule import WeightRule

urgent = PythonOperator(
    task_id         = "urgent_report",
    python_callable = run_report,
    priority_weight = 100,          # default is 1
    weight_rule     = WeightRule.ABSOLUTE,  # use exactly this weight
)

# weight_rule options:
# DOWNSTREAM (default): weight = own + sum of all downstream weights
# UPSTREAM:             weight = own + sum of all upstream weights
# ABSOLUTE:             weight = own weight only
```

---

## CeleryExecutor Scaling

```yaml
# docker-compose.yml — scale workers horizontally
worker:
  image: apache/airflow:2.9.0
  command: celery worker
  environment:
    AIRFLOW__CORE__EXECUTOR: CeleryExecutor
    AIRFLOW__CELERY__WORKER_CONCURRENCY: 16
  deploy:
    replicas: 4   # 4 workers × 16 concurrency = 64 total slots
```

```bash
# Scale workers at runtime (docker swarm / kubernetes)
docker service scale airflow_worker=8

# Graceful worker shutdown (finish current tasks, then stop)
airflow celery worker --autoscale=16,4  # max 16, min 4 concurrent tasks

# Worker queues for heterogeneous workloads
airflow celery worker -q default,heavy_compute
airflow celery worker -q gpu_tasks -c 2    # dedicated GPU worker, low concurrency
```

---

## KubernetesExecutor Scaling

```yaml
# The KubernetesExecutor scales to zero — no idle workers
# Configure the pod template for all tasks

# airflow.cfg
[kubernetes]
worker_pods_creation_batch_size = 16  # create up to 16 pods per scheduler loop
delete_worker_pods               = True
delete_worker_pods_on_failure    = False   # keep failed pods for debugging
worker_pods_pending_timeout      = 300     # fail task if pod stays pending > 5min
```

```python
# Per-task pod resource requests
from kubernetes.client import models as k8s

light_task = PythonOperator(
    task_id         = "light",
    python_callable = run_light,
    executor_config = {
        "pod_override": k8s.V1Pod(spec=k8s.V1PodSpec(containers=[
            k8s.V1Container(
                name="base",
                resources=k8s.V1ResourceRequirements(
                    requests={"cpu": "0.5", "memory": "512Mi"},
                    limits  ={"cpu": "1",   "memory": "1Gi"},
                ),
            )
        ]))
    },
)

heavy_task = PythonOperator(
    task_id         = "heavy",
    python_callable = run_heavy,
    executor_config = {
        "pod_override": k8s.V1Pod(spec=k8s.V1PodSpec(containers=[
            k8s.V1Container(
                name="base",
                resources=k8s.V1ResourceRequirements(
                    requests={"cpu": "4",   "memory": "16Gi"},
                    limits  ={"cpu": "8",   "memory": "32Gi"},
                ),
            )
        ], node_selector={"workload": "high-memory"}))
    },
)
```

---

## Database Connection Pooling

For large deployments, use PgBouncer between Airflow and PostgreSQL:

```yaml
# docker-compose.yml
pgbouncer:
  image: edoburu/pgbouncer
  environment:
    DATABASE_URL: postgresql://airflow:airflow@postgres:5432/airflow
    POOL_MODE:    transaction
    MAX_CLIENT_CONN: 200
    DEFAULT_POOL_SIZE: 20
  ports:
    - "5432:5432"
```

```properties
# airflow.cfg — point to PgBouncer instead of Postgres
[database]
sql_alchemy_conn = postgresql+psycopg2://airflow:airflow@pgbouncer:5432/airflow
sql_alchemy_pool_size = 5
sql_alchemy_max_overflow = 10
sql_alchemy_pool_recycle = 1800
```

---

## Performance Checklist

```
DAG parsing:
  ☐ No heavy imports at module level
  ☐ No network/DB calls at parse time
  ☐ No Variable.get() at module level
  ☐ DAG files < 1s parse time (measure with airflow dags report)

Scheduling:
  ☐ parallelism matches executor capacity
  ☐ parsing_processes = number of CPU cores / 2
  ☐ max_active_runs_per_dag prevents thundering herd

Worker capacity:
  ☐ worker_concurrency matched to workload (IO vs CPU)
  ☐ Pools configured for shared resources (DB, API)
  ☐ Priority weights set for critical pipelines

Database:
  ☐ PostgreSQL (not SQLite or MySQL for production)
  ☐ PgBouncer for high-concurrency deployments
  ☐ Regular VACUUM / ANALYZE on metadata DB
  ☐ Prune old dag runs: airflow db clean
```

---

## Pruning Metadata DB

```bash
# Clean old task instances and dag runs (keeps last 30 days)
airflow db clean --clean-before-timestamp "$(date -d '-30 days' --iso-8601=seconds)"

# Clean specific table
airflow db clean --tables task_instance --clean-before-timestamp 2024-01-01

# Dry run first
airflow db clean --dry-run --clean-before-timestamp 2024-01-01
```
