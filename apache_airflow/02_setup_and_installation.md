# Chapter 2: Setup & Installation

## Quick Start with Docker Compose (Recommended)

The official Docker Compose setup runs the full Airflow stack with one command.

```bash
# Download the official docker-compose.yml
curl -LfO 'https://airflow.apache.org/docs/apache-airflow/stable/docker-compose.yaml'

# Create required directories
mkdir -p ./dags ./logs ./plugins ./config

# Set the UID (Linux/Mac)
echo -e "AIRFLOW_UID=$(id -u)" > .env

# Initialize the database
docker compose up airflow-init

# Start all services
docker compose up -d

# Access UI: http://localhost:8080
# Default credentials: airflow / airflow
```

### Services Started

| Service | Description | Port |
|---|---|---|
| `airflow-webserver` | Web UI + REST API | 8080 |
| `airflow-scheduler` | DAG parsing + triggering | — |
| `airflow-worker` | Task execution (Celery) | — |
| `airflow-triggerer` | Deferred operators | — |
| `postgres` | Metadata database | 5432 |
| `redis` | Celery broker | 6379 |

---

## Minimal Docker Compose (LocalExecutor)

For development without Celery:

```yaml
# docker-compose-local.yml
version: "3.8"

x-airflow-common: &airflow-common
  image: apache/airflow:2.9.0
  environment:
    AIRFLOW__CORE__EXECUTOR: LocalExecutor
    AIRFLOW__DATABASE__SQL_ALCHEMY_CONN: postgresql+psycopg2://airflow:airflow@postgres/airflow
    AIRFLOW__CORE__FERNET_KEY: ""
    AIRFLOW__CORE__DAGS_ARE_PAUSED_AT_CREATION: "true"
    AIRFLOW__CORE__LOAD_EXAMPLES: "false"
    AIRFLOW__API__AUTH_BACKENDS: "airflow.api.auth.backend.basic_auth"
    _PIP_ADDITIONAL_REQUIREMENTS: "pandas sqlalchemy"
  volumes:
    - ./dags:/opt/airflow/dags
    - ./logs:/opt/airflow/logs
    - ./plugins:/opt/airflow/plugins
  depends_on:
    postgres:
      condition: service_healthy

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: airflow
      POSTGRES_PASSWORD: airflow
      POSTGRES_DB: airflow
    healthcheck:
      test: ["CMD", "pg_isready", "-U", "airflow"]
      interval: 5s
      retries: 5

  airflow-webserver:
    <<: *airflow-common
    command: webserver
    ports: ["8080:8080"]
    healthcheck:
      test: ["CMD", "curl", "--fail", "http://localhost:8080/health"]
      interval: 10s

  airflow-scheduler:
    <<: *airflow-common
    command: scheduler

  airflow-init:
    <<: *airflow-common
    command: >
      bash -c "
        airflow db migrate &&
        airflow users create --username admin --password admin
          --firstname Admin --lastname User --role Admin --email admin@example.com
      "
```

```bash
docker compose -f docker-compose-local.yml up -d
```

---

## Astro CLI (Astronomer)

The Astronomer Astro CLI is the easiest local development experience:

```bash
# Install Astro CLI
curl -sSL install.astronomer.io | sudo bash -s

# Initialize a new Airflow project
astro dev init

# Start Airflow locally
astro dev start
# UI: http://localhost:8080 (admin/admin)

# Stop
astro dev stop

# Restart
astro dev restart

# Run a DAG
astro dev run dags trigger my_dag

# Deploy to Astronomer Cloud
astro deploy
```

Project structure created by `astro dev init`:
```
my_project/
  dags/           ← your DAG files
  plugins/        ← custom plugins
  include/        ← helper modules
  tests/          ← DAG tests
  requirements.txt ← Python packages
  Dockerfile      ← custom image
  .env            ← local env vars
```

---

## Local Installation (pip)

```bash
# Python 3.8+ required
python -m venv airflow-env
source airflow-env/bin/activate

# Install Airflow with constraints (prevents dependency conflicts)
AIRFLOW_VERSION=2.9.0
PYTHON_VERSION="$(python --version | cut -d " " -f 2 | cut -d "." -f 1-2)"
CONSTRAINT_URL="https://raw.githubusercontent.com/apache/airflow/constraints-${AIRFLOW_VERSION}/constraints-${PYTHON_VERSION}.txt"

pip install "apache-airflow==${AIRFLOW_VERSION}" --constraint "${CONSTRAINT_URL}"

# Install providers you need
pip install apache-airflow-providers-postgres
pip install apache-airflow-providers-amazon
pip install apache-airflow-providers-google
pip install apache-airflow-providers-slack

# Initialize database (SQLite by default)
export AIRFLOW_HOME=~/airflow
airflow db init

# Create admin user
airflow users create \
  --username admin \
  --firstname Admin \
  --lastname User \
  --role Admin \
  --email admin@example.com \
  --password admin

# Start services (2 terminals)
airflow webserver --port 8080
airflow scheduler
```

---

## Configuration (airflow.cfg)

Key settings in `~/airflow/airflow.cfg` or via environment variables:

```properties
# Environment variable format: AIRFLOW__{SECTION}__{KEY}

# Core
AIRFLOW__CORE__EXECUTOR=LocalExecutor          # or CeleryExecutor, KubernetesExecutor
AIRFLOW__CORE__DAGS_FOLDER=/opt/airflow/dags
AIRFLOW__CORE__LOAD_EXAMPLES=False
AIRFLOW__CORE__MAX_ACTIVE_RUNS_PER_DAG=3
AIRFLOW__CORE__MAX_ACTIVE_TASKS_PER_DAG=10
AIRFLOW__CORE__PARALLELISM=32                  # total concurrent tasks cluster-wide
AIRFLOW__CORE__DAGBAG_IMPORT_TIMEOUT=30        # seconds to import a DAG file

# Database
AIRFLOW__DATABASE__SQL_ALCHEMY_CONN=postgresql+psycopg2://user:pass@host/db

# Scheduler
AIRFLOW__SCHEDULER__MIN_FILE_PROCESS_INTERVAL=30   # DAG file re-parse interval
AIRFLOW__SCHEDULER__DAG_DIR_LIST_INTERVAL=30       # scan dags folder interval
AIRFLOW__SCHEDULER__MAX_TIS_PER_QUERY=512

# Webserver
AIRFLOW__WEBSERVER__BASE_URL=http://localhost:8080
AIRFLOW__WEBSERVER__SECRET_KEY=your-secret-key
AIRFLOW__WEBSERVER__WORKERS=4

# Email
AIRFLOW__SMTP__SMTP_HOST=smtp.gmail.com
AIRFLOW__SMTP__SMTP_PORT=587
AIRFLOW__SMTP__SMTP_USER=your@email.com
AIRFLOW__SMTP__SMTP_PASSWORD=app_password
AIRFLOW__SMTP__SMTP_MAIL_FROM=airflow@yourco.com
```

---

## First DAG: Verify Setup

```python
# dags/hello_world.py
from airflow import DAG
from airflow.operators.python import PythonOperator
from airflow.operators.bash import BashOperator
from datetime import datetime

def say_hello():
    print("Hello from Airflow!")
    return "success"

with DAG(
    dag_id      = "hello_world",
    start_date  = datetime(2024, 1, 1),
    schedule    = None,             # manual trigger only
    catchup     = False,
    tags        = ["example"],
) as dag:

    greet = PythonOperator(
        task_id          = "say_hello",
        python_callable  = say_hello,
    )

    echo = BashOperator(
        task_id      = "echo_date",
        bash_command = "echo 'Current date: $(date)'",
    )

    greet >> echo
```

```bash
# Trigger the DAG
airflow dags trigger hello_world

# Check status
airflow dags list
airflow tasks list hello_world
airflow dags state hello_world $(date +%Y-%m-%dT%H:%M:%S)
```
