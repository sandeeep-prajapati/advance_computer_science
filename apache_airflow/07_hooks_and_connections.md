# Chapter 7: Hooks & Connections

## Connections

Connections store credentials for external systems. They are the single source of truth for how Airflow connects to databases, cloud services, and APIs.

### Managing Connections

**Via UI:** Admin → Connections → "+"

**Via CLI:**
```bash
# Add a PostgreSQL connection
airflow connections add \
  --conn-id       my_postgres \
  --conn-type     postgres \
  --conn-host     localhost \
  --conn-port     5432 \
  --conn-login    myuser \
  --conn-password mypassword \
  --conn-schema   mydb \
  --conn-extra    '{"sslmode": "require"}'

# Add an S3 connection (AWS)
airflow connections add \
  --conn-id   aws_default \
  --conn-type aws \
  --conn-extra '{"aws_access_key_id": "AKIA...", "aws_secret_access_key": "..."}'

# List all connections
airflow connections list

# Delete
airflow connections delete my_old_conn
```

**Via environment variables** (best for containers):
```bash
# URI format
export AIRFLOW_CONN_MY_POSTGRES="postgresql://myuser:mypassword@localhost:5432/mydb"

# JSON format
export AIRFLOW_CONN_MY_POSTGRES='{
  "conn_type": "postgres",
  "host": "localhost",
  "login": "myuser",
  "password": "mypassword",
  "schema": "mydb",
  "port": 5432
}'
```

---

## Common Connection Types

| Conn Type | Provider Package | Use Case |
|---|---|---|
| `postgres` | apache-airflow-providers-postgres | PostgreSQL |
| `mysql` | apache-airflow-providers-mysql | MySQL |
| `aws` | apache-airflow-providers-amazon | AWS (S3, Redshift, etc.) |
| `google_cloud_platform` | apache-airflow-providers-google | GCP (BigQuery, GCS, etc.) |
| `azure` | apache-airflow-providers-microsoft-azure | Azure services |
| `spark` | apache-airflow-providers-apache-spark | Spark cluster |
| `kafka` | apache-airflow-providers-apache-kafka | Kafka |
| `http` | apache-airflow-providers-http | HTTP/REST APIs |
| `ssh` | apache-airflow-providers-ssh | SSH connections |
| `slack` | apache-airflow-providers-slack | Slack webhooks |
| `snowflake` | apache-airflow-providers-snowflake | Snowflake DWH |
| `databricks` | apache-airflow-providers-databricks | Databricks |

---

## Hooks

Hooks are the Python interface to external services. They read a `conn_id` and manage the connection lifecycle.

### PostgresHook

```python
from airflow.providers.postgres.hooks.postgres import PostgresHook
import pandas as pd

hook = PostgresHook(postgres_conn_id="my_postgres")

# Run SQL
hook.run("INSERT INTO logs VALUES (%s, %s)", parameters=["2024-01-01", "success"])

# Get records as list of tuples
records = hook.get_records("SELECT id, name FROM users LIMIT 100")

# Get first record
first = hook.get_first("SELECT MAX(created_at) FROM orders")

# Get Pandas DataFrame
df = hook.get_pandas_df("SELECT * FROM orders WHERE date = '{{ ds }}'")

# Get SQLAlchemy engine
engine = hook.get_sqlalchemy_engine()
df.to_sql("staging", engine, if_exists="append", index=False)

# Execute in PythonOperator
def load_to_postgres(**context):
    hook = PostgresHook(postgres_conn_id="my_postgres")
    df = pd.DataFrame(...)
    hook.insert_rows(
        table  = "daily_summary",
        rows   = df.to_records(index=False).tolist(),
        target_fields = df.columns.tolist(),
        replace = True,
        replace_index = ["date"],
    )
```

### S3Hook

```python
from airflow.providers.amazon.aws.hooks.s3 import S3Hook

hook = S3Hook(aws_conn_id="aws_default")

# Check if key exists
exists = hook.check_for_key("raw/2024-01-01/data.csv", "my-bucket")

# Download file
hook.download_file(
    key      = "raw/{{ ds }}/data.csv",
    bucket_name = "my-bucket",
    local_path = "/tmp/",
)

# Upload file
hook.load_file(
    filename    = "/tmp/output.parquet",
    key         = "processed/{{ ds }}/output.parquet",
    bucket_name = "my-bucket",
    replace     = True,
)

# List keys
keys = hook.list_keys(bucket_name="my-bucket", prefix="raw/2024-01-01/")

# Read file content
content = hook.read_key("config/settings.json", "my-bucket")
```

### HttpHook

```python
from airflow.providers.http.hooks.http import HttpHook

hook = HttpHook(http_conn_id="my_api", method="POST")

response = hook.run(
    endpoint = "/api/v1/process",
    data     = {"date": "2024-01-01"},
    headers  = {"Content-Type": "application/json"},
)
data = response.json()
```

### SlackHook

```python
from airflow.providers.slack.hooks.slack_webhook import SlackWebhookHook

def notify_slack(message: str) -> None:
    hook = SlackWebhookHook(slack_webhook_conn_id="slack_webhook")
    hook.send(text=message, icon_emoji=":airflow:")

# In a callback
def on_failure_callback(context):
    dag_id  = context["dag"].dag_id
    task_id = context["task"].task_id
    ds      = context["ds"]
    notify_slack(f":red_circle: Task failed: `{dag_id}.{task_id}` for `{ds}`")
```

---

## Secrets Backends

Instead of storing credentials in the Airflow metadata database, use an external secrets manager.

### AWS Secrets Manager

```properties
# airflow.cfg
[secrets]
backend = airflow.providers.amazon.aws.secrets.secrets_manager.SecretsManagerBackend
backend_kwargs = {
    "connections_prefix": "airflow/connections",
    "variables_prefix": "airflow/variables",
    "config_prefix": "airflow/config"
}
```

```bash
# Store connection in AWS Secrets Manager
aws secretsmanager create-secret \
  --name "airflow/connections/my_postgres" \
  --secret-string "postgresql://user:pass@host:5432/db"
```

### HashiCorp Vault

```properties
[secrets]
backend = airflow.providers.hashicorp.secrets.vault.VaultBackend
backend_kwargs = {
    "connections_path": "connections",
    "variables_path": "variables",
    "mount_point": "airflow",
    "url": "http://vault:8200",
    "token": "vault-token"
}
```

---

## Custom Hook

```python
from airflow.hooks.base import BaseHook
import requests

class MyAPIHook(BaseHook):
    conn_name_attr  = "my_api_conn_id"
    default_conn_name = "my_api_default"

    def __init__(self, my_api_conn_id: str = "my_api_default"):
        super().__init__()
        self.my_api_conn_id = my_api_conn_id
        self._session = None

    def get_conn(self):
        if self._session is None:
            conn = self.get_connection(self.my_api_conn_id)
            self._session = requests.Session()
            self._session.headers.update({
                "Authorization": f"Bearer {conn.password}",
                "Content-Type":  "application/json",
            })
            self._base_url = f"http://{conn.host}:{conn.port}"
        return self._session

    def get_data(self, endpoint: str, params: dict = None) -> dict:
        session = self.get_conn()
        response = session.get(f"{self._base_url}{endpoint}", params=params)
        response.raise_for_status()
        return response.json()

    def post_data(self, endpoint: str, payload: dict) -> dict:
        session = self.get_conn()
        response = session.post(f"{self._base_url}{endpoint}", json=payload)
        response.raise_for_status()
        return response.json()
```
