# Chapter 18: Custom Operators & Plugins

## Custom Operator

Extend `BaseOperator` to encapsulate reusable task logic:

```python
# plugins/operators/s3_to_postgres_operator.py
from airflow.models import BaseOperator
from airflow.providers.amazon.aws.hooks.s3 import S3Hook
from airflow.providers.postgres.hooks.postgres import PostgresHook
import pandas as pd
import io


class S3ToPostgresOperator(BaseOperator):
    """
    Load a Parquet file from S3 into a PostgreSQL table.
    Supports REPLACE, APPEND, and UPSERT modes.
    """

    template_fields = ("s3_key", "table")   # Jinja-templated fields

    def __init__(
        self,
        s3_conn_id:       str,
        s3_bucket:        str,
        s3_key:           str,
        postgres_conn_id: str,
        table:            str,
        schema:           str = "public",
        mode:             str = "replace",   # replace | append | upsert
        upsert_keys:      list[str] = None,
        **kwargs,
    ):
        super().__init__(**kwargs)
        self.s3_conn_id       = s3_conn_id
        self.s3_bucket        = s3_bucket
        self.s3_key           = s3_key
        self.postgres_conn_id = postgres_conn_id
        self.table            = table
        self.schema           = schema
        self.mode             = mode
        self.upsert_keys      = upsert_keys or []

    def execute(self, context):
        self.log.info("Reading s3://%s/%s", self.s3_bucket, self.s3_key)

        s3_hook = S3Hook(aws_conn_id=self.s3_conn_id)
        obj     = s3_hook.get_key(self.s3_key, bucket_name=self.s3_bucket)
        df      = pd.read_parquet(io.BytesIO(obj.get()["Body"].read()))

        self.log.info("Loaded %d rows from S3", len(df))

        pg_hook = PostgresHook(postgres_conn_id=self.postgres_conn_id)
        engine  = pg_hook.get_sqlalchemy_engine()

        if self.mode == "replace":
            df.to_sql(self.table, engine, schema=self.schema,
                      if_exists="replace", index=False)
        elif self.mode == "append":
            df.to_sql(self.table, engine, schema=self.schema,
                      if_exists="append", index=False)
        elif self.mode == "upsert":
            self._upsert(df, pg_hook)
        else:
            raise ValueError(f"Unknown mode: {self.mode}")

        self.log.info("Loaded %d rows into %s.%s", len(df), self.schema, self.table)
        return len(df)

    def _upsert(self, df: pd.DataFrame, hook: PostgresHook) -> None:
        staging = f"_staging_{self.table}"
        engine  = hook.get_sqlalchemy_engine()
        df.to_sql(staging, engine, schema=self.schema, if_exists="replace", index=False)

        key_clause    = " AND ".join(f"t.{k} = s.{k}" for k in self.upsert_keys)
        update_clause = ", ".join(
            f"{c} = s.{c}" for c in df.columns if c not in self.upsert_keys
        )
        insert_cols = ", ".join(df.columns)
        insert_vals = ", ".join(f"s.{c}" for c in df.columns)

        sql = f"""
            INSERT INTO {self.schema}.{self.table} ({insert_cols})
            SELECT {insert_vals} FROM {self.schema}.{staging} s
            ON CONFLICT ({", ".join(self.upsert_keys)})
            DO UPDATE SET {update_clause}
        """
        hook.run(sql)
        hook.run(f"DROP TABLE IF EXISTS {self.schema}.{staging}")
```

### Using the Custom Operator

```python
from plugins.operators.s3_to_postgres_operator import S3ToPostgresOperator

load_orders = S3ToPostgresOperator(
    task_id          = "load_orders",
    s3_conn_id       = "aws_default",
    s3_bucket        = "my-data-bucket",
    s3_key           = "processed/{{ ds }}/orders.parquet",
    postgres_conn_id = "postgres_default",
    table            = "orders",
    schema           = "staging",
    mode             = "upsert",
    upsert_keys      = ["order_id"],
)
```

---

## Custom Sensor

```python
# plugins/sensors/api_data_sensor.py
from airflow.sensors.base import BaseSensorOperator
import requests


class APIDataReadySensor(BaseSensorOperator):
    """Polls an internal API until data is available for the given date."""

    template_fields = ("endpoint", "date")

    def __init__(self, base_url: str, endpoint: str, date: str, min_records: int = 1, **kwargs):
        super().__init__(**kwargs)
        self.base_url    = base_url
        self.endpoint    = endpoint
        self.date        = date
        self.min_records = min_records

    def poke(self, context) -> bool:
        url = f"{self.base_url}/availability/{self.endpoint}"
        try:
            resp = requests.get(url, params={"date": self.date}, timeout=10)
            resp.raise_for_status()
            count = resp.json().get("record_count", 0)
            self.log.info("Records available: %d / %d needed", count, self.min_records)
            return count >= self.min_records
        except requests.RequestException as e:
            self.log.warning("API check failed: %s", e)
            return False
```

---

## Custom Hook

```python
# plugins/hooks/internal_api_hook.py
from airflow.hooks.base import BaseHook
import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry


class InternalAPIHook(BaseHook):
    """Hook for the internal data API with built-in retries."""

    conn_name_attr    = "api_conn_id"
    default_conn_name = "internal_api_default"

    def __init__(self, api_conn_id: str = "internal_api_default"):
        super().__init__()
        self.api_conn_id = api_conn_id
        self._session    = None

    def get_conn(self) -> requests.Session:
        if self._session is None:
            conn       = self.get_connection(self.api_conn_id)
            self._base_url = f"https://{conn.host}:{conn.port or 443}"

            session = requests.Session()
            retry   = Retry(total=3, backoff_factor=1, status_forcelist=[500, 502, 503])
            adapter = HTTPAdapter(max_retries=retry)
            session.mount("https://", adapter)
            session.headers.update({
                "Authorization": f"Bearer {conn.password}",
                "Content-Type":  "application/json",
            })
            self._session = session
        return self._session

    def get_data(self, path: str, params: dict = None) -> dict:
        resp = self.get_conn().get(f"{self._base_url}{path}", params=params)
        resp.raise_for_status()
        return resp.json()

    def post_data(self, path: str, payload: dict) -> dict:
        resp = self.get_conn().post(f"{self._base_url}{path}", json=payload)
        resp.raise_for_status()
        return resp.json()
```

---

## Plugins

Plugins let you register operators, hooks, sensors, macros, blueprints, and UI menu items:

```python
# plugins/my_plugin.py
from airflow.plugins_manager import AirflowPlugin
from plugins.operators.s3_to_postgres_operator import S3ToPostgresOperator
from plugins.sensors.api_data_sensor import APIDataReadySensor
from plugins.hooks.internal_api_hook import InternalAPIHook


def greeting_macro(name: str) -> str:
    return f"Hello, {name}!"


class MyCompanyPlugin(AirflowPlugin):
    name     = "my_company_plugin"

    operators = [S3ToPostgresOperator]
    sensors   = [APIDataReadySensor]
    hooks     = [InternalAPIHook]

    macros    = [greeting_macro]   # available as {{ greeting_macro("world") }}
```

Place the file in the `plugins/` directory and Airflow auto-discovers it.

---

## Decorating an Existing Operator

```python
# Add pre/post hooks to any operator by subclassing
from airflow.providers.postgres.operators.postgres import PostgresOperator
from datetime import datetime


class AuditedPostgresOperator(PostgresOperator):
    """PostgresOperator that logs execution to an audit table."""

    def execute(self, context):
        start = datetime.utcnow()
        try:
            result = super().execute(context)
            self._write_audit(context, "success", start)
            return result
        except Exception as e:
            self._write_audit(context, "failed", start, str(e))
            raise

    def _write_audit(self, context, status: str, start: datetime, error: str = None):
        from airflow.providers.postgres.hooks.postgres import PostgresHook
        hook = PostgresHook(postgres_conn_id=self.postgres_conn_id)
        hook.run(
            """
            INSERT INTO audit.task_runs (dag_id, task_id, ds, status, started_at, error)
            VALUES (%s, %s, %s, %s, %s, %s)
            """,
            parameters=(
                context["dag"].dag_id,
                context["task"].task_id,
                context["ds"],
                status,
                start,
                error,
            ),
        )
```

---

## Custom XCom Backend as a Plugin

```python
# plugins/xcoms/gcs_xcom_backend.py
from airflow.models.xcom import BaseXCom
import json, uuid
from google.cloud import storage


class GCSXComBackend(BaseXCom):
    BUCKET = "my-airflow-xcoms"

    @staticmethod
    def serialize_value(value, *, key=None, task_id=None, dag_id=None,
                        run_id=None, map_index=None):
        import pandas as pd
        if isinstance(value, pd.DataFrame):
            gcs_key = f"xcoms/{dag_id}/{run_id}/{task_id}/{uuid.uuid4()}.parquet"
            client  = storage.Client()
            bucket  = client.bucket(GCSXComBackend.BUCKET)
            bucket.blob(gcs_key).upload_from_string(value.to_parquet())
            return BaseXCom.serialize_value({"__type": "gcs_parquet", "key": gcs_key})
        return BaseXCom.serialize_value(value)

    @staticmethod
    def deserialize_value(result):
        import pandas as pd
        data = BaseXCom.deserialize_value(result)
        if isinstance(data, dict) and data.get("__type") == "gcs_parquet":
            client = storage.Client()
            bucket = client.bucket(GCSXComBackend.BUCKET)
            return pd.read_parquet(bucket.blob(data["key"]).download_as_bytes())
        return data
```

```properties
# airflow.cfg
[core]
xcom_backend = plugins.xcoms.gcs_xcom_backend.GCSXComBackend
```

---

## Plugin Directory Layout

```
project/
├── dags/
│   └── my_dag.py
├── plugins/
│   ├── __init__.py
│   ├── my_plugin.py              ← AirflowPlugin registration
│   ├── operators/
│   │   ├── __init__.py
│   │   └── s3_to_postgres_operator.py
│   ├── sensors/
│   │   ├── __init__.py
│   │   └── api_data_sensor.py
│   ├── hooks/
│   │   ├── __init__.py
│   │   └── internal_api_hook.py
│   └── xcoms/
│       ├── __init__.py
│       └── gcs_xcom_backend.py
└── tests/
    └── test_operators.py
```

---

## Testing Custom Operators

```python
# tests/test_s3_to_postgres_operator.py
from unittest.mock import MagicMock, patch
import pandas as pd
from plugins.operators.s3_to_postgres_operator import S3ToPostgresOperator

@patch("plugins.operators.s3_to_postgres_operator.S3Hook")
@patch("plugins.operators.s3_to_postgres_operator.PostgresHook")
def test_execute_append_mode(mock_pg_hook_cls, mock_s3_hook_cls):
    # Arrange
    sample_df   = pd.DataFrame({"id": [1, 2], "value": ["a", "b"]})
    mock_s3     = MagicMock()
    mock_s3.get_key.return_value.get.return_value = {"Body": sample_df.to_parquet()}
    mock_s3_hook_cls.return_value = mock_s3

    mock_pg     = MagicMock()
    mock_engine = MagicMock()
    mock_pg.get_sqlalchemy_engine.return_value = mock_engine
    mock_pg_hook_cls.return_value = mock_pg

    op = S3ToPostgresOperator(
        task_id          = "test",
        s3_conn_id       = "aws",
        s3_bucket        = "bucket",
        s3_key           = "key.parquet",
        postgres_conn_id = "pg",
        table            = "orders",
        mode             = "append",
    )

    # Act
    result = op.execute({})

    # Assert
    assert result == 2
    mock_s3.get_key.assert_called_once_with("key.parquet", bucket_name="bucket")
```
