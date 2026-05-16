# Chapter 12: Testing DAGs

## Why Test DAGs?

```
Parse errors         → DAG won't load at all
Import errors        → Silent failures in production
Wrong dependencies   → Tasks run in wrong order
Broken logic         → Incorrect data silently written
Missing connections  → Task fails at runtime, not at review time
```

---

## DAG Validation Tests

The fastest, cheapest tests — just import the DAG file:

```python
# tests/test_dag_validation.py
import os
import pytest
from airflow.models import DagBag

DAG_FOLDER = os.path.join(os.path.dirname(__file__), "../dags")

@pytest.fixture(scope="module")
def dagbag():
    return DagBag(dag_folder=DAG_FOLDER, include_examples=False)

def test_no_import_errors(dagbag):
    assert dagbag.import_errors == {}, f"DAG import errors: {dagbag.import_errors}"

def test_all_dags_loaded(dagbag):
    assert len(dagbag.dags) > 0, "No DAGs found"

def test_dag_ids_are_unique(dagbag):
    dag_ids = list(dagbag.dags.keys())
    assert len(dag_ids) == len(set(dag_ids))

@pytest.mark.parametrize("dag_id", ["my_etl_dag", "daily_report_dag"])
def test_specific_dag_exists(dagbag, dag_id):
    assert dag_id in dagbag.dags, f"DAG '{dag_id}' not found"
```

---

## Task Structure Tests

```python
# tests/test_my_etl_dag.py
import pytest
from airflow.models import DagBag

@pytest.fixture(scope="module")
def dag():
    dagbag = DagBag(dag_folder="dags/", include_examples=False)
    return dagbag.get_dag("my_etl_dag")

def test_dag_loaded(dag):
    assert dag is not None

def test_task_count(dag):
    assert len(dag.tasks) == 4

def test_task_ids(dag):
    task_ids = {t.task_id for t in dag.tasks}
    assert task_ids == {"extract", "transform", "load", "notify"}

def test_dependencies(dag):
    extract   = dag.get_task("extract")
    transform = dag.get_task("transform")
    load      = dag.get_task("load")
    notify    = dag.get_task("notify")

    assert transform in extract.downstream_list
    assert load      in transform.downstream_list
    assert notify    in load.downstream_list

def test_no_catchup(dag):
    assert dag.catchup is False

def test_retry_config(dag):
    extract = dag.get_task("extract")
    assert extract.retries >= 1
```

---

## Unit Testing Task Logic

Test the Python callable directly — no Airflow context needed:

```python
# dags/my_etl_dag.py
def extract_records(date: str) -> list[dict]:
    """Extract orders from DB for a given date."""
    # In production this hits a real DB
    return []  # stub

def transform_records(records: list[dict]) -> list[dict]:
    return [
        {**r, "amount_usd": r["amount"] / 100}
        for r in records
        if r.get("status") == "completed"
    ]

def load_records(records: list[dict], table: str) -> int:
    return len(records)
```

```python
# tests/test_task_logic.py
from dags.my_etl_dag import transform_records, load_records

def test_transform_filters_non_completed():
    records = [
        {"id": 1, "amount": 1000, "status": "completed"},
        {"id": 2, "amount": 500,  "status": "pending"},
        {"id": 3, "amount": 200,  "status": "completed"},
    ]
    result = transform_records(records)
    assert len(result) == 2
    assert all(r["status"] == "completed" for r in result)

def test_transform_converts_amount():
    records = [{"id": 1, "amount": 1000, "status": "completed"}]
    result = transform_records(records)
    assert result[0]["amount_usd"] == 10.0

def test_transform_empty_input():
    assert transform_records([]) == []

def test_load_returns_count():
    records = [{"id": i} for i in range(5)]
    assert load_records(records, "orders") == 5
```

---

## Testing with Mocked Context

```python
# dags/report_dag.py
from airflow.operators.python import PythonOperator

def send_report(**context):
    ds         = context["ds"]
    dag_run    = context["dag_run"]
    task_inst  = context["ti"]

    path = task_inst.xcom_pull(task_ids="generate", key="report_path")
    # send report at path for ds
    return f"Sent report for {ds}"
```

```python
# tests/test_report_task.py
from unittest.mock import MagicMock, patch
from dags.report_dag import send_report

def test_send_report_uses_correct_ds():
    ti = MagicMock()
    ti.xcom_pull.return_value = "s3://bucket/reports/2024-01-01/report.pdf"

    context = {
        "ds":      "2024-01-01",
        "dag_run": MagicMock(),
        "ti":      ti,
    }

    result = send_report(**context)

    ti.xcom_pull.assert_called_once_with(task_ids="generate", key="report_path")
    assert "2024-01-01" in result
```

---

## Testing with pytest + airflow DB

For tests that need the metadata DB (XCom, task state):

```python
# conftest.py
import pytest
from airflow.models import DagBag, TaskInstance
from airflow.utils.state import State
from airflow.utils.types import DagRunType
from datetime import datetime

@pytest.fixture(scope="session")
def airflow_db():
    """Initialize Airflow DB for tests."""
    from airflow.utils.db import initdb
    initdb()

@pytest.fixture
def create_dag_run(airflow_db):
    from airflow.models import DagRun
    def _create(dag, execution_date=None):
        execution_date = execution_date or datetime(2024, 1, 1)
        dag_run = dag.create_dagrun(
            run_id          = f"test__{execution_date}",
            execution_date  = execution_date,
            state           = State.RUNNING,
            run_type        = DagRunType.MANUAL,
        )
        return dag_run
    return _create
```

```python
# tests/test_xcom_flow.py
from airflow.models import DagBag, TaskInstance
from datetime import datetime

def test_extract_pushes_xcom(create_dag_run):
    dagbag   = DagBag(dag_folder="dags/", include_examples=False)
    dag      = dagbag.get_dag("xcom_test_dag")
    dag_run  = create_dag_run(dag)

    ti = TaskInstance(task=dag.get_task("extract"), run_id=dag_run.run_id)
    ti.run(ignore_ti_state=True)

    value = ti.xcom_pull(task_ids="extract", key="row_count")
    assert value == 100
```

---

## Integration Testing with Docker

```python
# tests/test_integration.py
import subprocess
import pytest

@pytest.fixture(scope="session")
def airflow_services():
    """Start Airflow via docker-compose for integration tests."""
    subprocess.run(
        ["docker-compose", "-f", "docker-compose.test.yml", "up", "-d"],
        check=True,
    )
    # Wait for webserver to be ready
    import time; time.sleep(30)
    yield
    subprocess.run(
        ["docker-compose", "-f", "docker-compose.test.yml", "down", "-v"],
    )

def test_dag_trigger_via_api(airflow_services):
    import requests
    response = requests.post(
        "http://localhost:8080/api/v1/dags/my_etl_dag/dagRuns",
        auth = ("admin", "admin"),
        json = {"conf": {}, "logical_date": "2024-01-01T00:00:00Z"},
    )
    assert response.status_code == 200
    dag_run_id = response.json()["dag_run_id"]
    assert dag_run_id is not None
```

---

## Testing Dynamic Task Mapping

```python
# tests/test_dynamic_dag.py
from airflow.models import DagBag

def test_dynamic_mapping_tasks():
    dagbag = DagBag(dag_folder="dags/", include_examples=False)
    dag    = dagbag.get_dag("dynamic_etl_dag")

    # Verify the mapped task exists
    task_ids = {t.task_id for t in dag.tasks}
    assert "process_file" in task_ids

    # Verify it's a mapped operator
    process_task = dag.get_task("process_file")
    assert hasattr(process_task, "expand_input")
```

---

## CI/CD Pipeline for DAGs

```yaml
# .github/workflows/test-dags.yml
name: Test DAGs

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    env:
      AIRFLOW_HOME: /tmp/airflow
      AIRFLOW__CORE__LOAD_EXAMPLES: "False"
      AIRFLOW__DATABASE__SQL_ALCHEMY_CONN: "sqlite:////tmp/airflow/airflow.db"

    steps:
      - uses: actions/checkout@v3

      - name: Set up Python
        uses: actions/setup-python@v4
        with:
          python-version: "3.11"

      - name: Install dependencies
        run: |
          pip install apache-airflow==2.9.0 pytest pytest-cov

      - name: Initialize Airflow DB
        run: airflow db init

      - name: Run tests
        run: |
          pytest tests/ -v --cov=dags --cov-report=xml

      - name: Upload coverage
        uses: codecov/codecov-action@v3
```

---

## Test Project Structure

```
project/
├── dags/
│   ├── my_etl_dag.py
│   ├── report_dag.py
│   └── dynamic_dag.py
├── tests/
│   ├── conftest.py
│   ├── test_dag_validation.py    # Import and structure checks
│   ├── test_task_logic.py        # Pure function unit tests
│   ├── test_my_etl_dag.py        # DAG-level assertions
│   └── test_integration.py       # End-to-end (optional)
├── pytest.ini
└── requirements-test.txt
```

```ini
# pytest.ini
[pytest]
testpaths = tests
log_cli = true
log_level = INFO
```

---

## Testing Checklist

```
DAG validation:
  ☐ No import errors (DagBag.import_errors is empty)
  ☐ All expected DAGs are present
  ☐ No duplicate task IDs
  ☐ Correct dependencies (upstream/downstream)
  ☐ catchup=False on production DAGs

Task logic:
  ☐ Happy path tested
  ☐ Empty input handled
  ☐ Edge cases covered
  ☐ Side effects mocked (DB, S3, HTTP)

Configuration:
  ☐ Retries configured
  ☐ Timeouts configured
  ☐ Correct schedule
  ☐ Correct start_date
```
