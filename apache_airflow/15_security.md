# Chapter 15: Security

## Authentication

### Default Setup (Username/Password)

```python
# webserver_config.py
from flask_appbuilder.security.manager import AUTH_DB

AUTH_TYPE = AUTH_DB

# Create admin user via CLI
# airflow users create --role Admin --username admin --firstname Admin \
#   --lastname User --email admin@example.com --password admin123
```

### OAuth2 / OIDC (Google, GitHub, Okta)

```python
# webserver_config.py
from flask_appbuilder.security.manager import AUTH_OAUTH

AUTH_TYPE       = AUTH_OAUTH
AUTH_USER_REGISTRATION      = True
AUTH_USER_REGISTRATION_ROLE = "Viewer"

OAUTH_PROVIDERS = [
    {
        "name":        "google",
        "token_key":   "access_token",
        "icon":        "fa-google",
        "remote_app":  {
            "client_id":     "your-client-id.apps.googleusercontent.com",
            "client_secret": "your-client-secret",
            "api_base_url":  "https://www.googleapis.com/oauth2/v2/",
            "client_kwargs": {"scope": "email profile"},
            "request_token_url": None,
            "access_token_url":  "https://accounts.google.com/o/oauth2/token",
            "authorize_url":     "https://accounts.google.com/o/oauth2/auth",
        },
    }
]
```

### LDAP Authentication

```python
# webserver_config.py
from flask_appbuilder.security.manager import AUTH_LDAP

AUTH_TYPE         = AUTH_LDAP
AUTH_LDAP_SERVER  = "ldap://ldap.example.com:389"
AUTH_LDAP_USE_TLS = True
AUTH_LDAP_SEARCH   = "dc=example,dc=com"
AUTH_LDAP_BIND_USER     = "cn=service,dc=example,dc=com"
AUTH_LDAP_BIND_PASSWORD = "bind-password"
AUTH_LDAP_UID_FIELD     = "uid"
```

---

## RBAC (Role-Based Access Control)

Airflow ships with five built-in roles:

| Role | Permissions |
|---|---|
| `Admin` | Full access including security management |
| `Op` | Create/edit DAGs, connections, variables |
| `User` | Trigger DAGs, view runs and logs |
| `Viewer` | Read-only access |
| `Public` | No access (anonymous users) |

### Managing Users and Roles

```bash
# Create a user
airflow users create \
  --role   User \
  --username john.doe \
  --firstname John \
  --lastname  Doe \
  --email john.doe@example.com \
  --password  SecurePass123

# List users
airflow users list

# Add role to user
airflow users add-role --username john.doe --role Op

# Remove role
airflow users remove-role --username john.doe --role Op

# List roles
airflow roles list

# Create custom role
airflow roles create data_engineer
```

### Custom Roles via Python

```python
# webserver_config.py
from airflow.www.security import AirflowSecurityManager

class MySecurityManager(AirflowSecurityManager):

    def get_user_roles(self, user=None):
        roles = super().get_user_roles(user)
        # Add custom logic: e.g., assign role based on email domain
        if user and user.email.endswith("@data-team.example.com"):
            roles.append(self.find_role("Op"))
        return roles

SECURITY_MANAGER_CLASS = MySecurityManager
```

### DAG-Level Access Control

```python
# In the DAG definition, restrict which roles can view/trigger
with DAG(
    dag_id       = "finance_pipeline",
    access_control = {
        "finance_team": {"can_read", "can_edit", "can_dag_run"},
        "Admin":        {"can_read", "can_edit", "can_dag_run", "can_delete"},
    },
    ...
) as dag:
    pass
```

---

## Fernet Encryption

Airflow encrypts connection passwords and variable values in the metadata DB using Fernet symmetric encryption.

```bash
# Generate a Fernet key
python -c "from cryptography.fernet import Fernet; print(Fernet.generate_key().decode())"
```

```properties
# airflow.cfg
[core]
fernet_key = your_fernet_key_here

# Or via environment variable
AIRFLOW__CORE__FERNET_KEY=your_fernet_key_here
```

**Important:** Keep the Fernet key secret and back it up. Losing it means you cannot decrypt stored passwords.

### Key Rotation

```bash
# After setting a new fernet_key alongside the old one:
airflow rotate-fernet-key
```

---

## Secrets Backends

Instead of storing secrets in the Airflow metadata DB, delegate to a dedicated secrets manager.

### AWS Secrets Manager

```properties
# airflow.cfg
[secrets]
backend = airflow.providers.amazon.aws.secrets.secrets_manager.SecretsManagerBackend
backend_kwargs = {
    "connections_prefix": "airflow/connections",
    "variables_prefix":   "airflow/variables",
    "config_prefix":      "airflow/config"
}
```

```bash
# Store a connection as a plain URI
aws secretsmanager create-secret \
  --name "airflow/connections/my_postgres" \
  --secret-string "postgresql://user:pass@host:5432/db"

# Store a variable
aws secretsmanager create-secret \
  --name "airflow/variables/api_key" \
  --secret-string "super-secret-key"
```

### HashiCorp Vault

```properties
[secrets]
backend = airflow.providers.hashicorp.secrets.vault.VaultBackend
backend_kwargs = {
    "connections_path": "connections",
    "variables_path":   "variables",
    "mount_point":      "airflow",
    "url":              "http://vault:8200",
    "auth_type":        "token",
    "token":            "vault-root-token"
}
```

```bash
# Store in Vault
vault kv put airflow/connections/my_postgres \
  conn_uri="postgresql://user:pass@host:5432/db"

vault kv put airflow/variables/batch_size value=500
```

### GCP Secret Manager

```properties
[secrets]
backend = airflow.providers.google.cloud.secrets.secret_manager.CloudSecretManagerBackend
backend_kwargs = {
    "connections_prefix": "airflow-connections",
    "variables_prefix":   "airflow-variables",
    "gcp_key_path":       "/opt/airflow/gcp_service_account.json"
}
```

---

## Network Security

### TLS for Webserver

```properties
# airflow.cfg
[webserver]
web_server_ssl_cert = /opt/airflow/certs/airflow.crt
web_server_ssl_key  = /opt/airflow/certs/airflow.key
```

### Webserver Behind a Reverse Proxy (Nginx)

```nginx
# /etc/nginx/sites-available/airflow
server {
    listen 443 ssl;
    server_name airflow.example.com;

    ssl_certificate     /etc/nginx/certs/airflow.crt;
    ssl_certificate_key /etc/nginx/certs/airflow.key;

    location / {
        proxy_pass         http://localhost:8080;
        proxy_set_header   Host $host;
        proxy_set_header   X-Real-IP $remote_addr;
        proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
    }
}
```

```properties
# airflow.cfg — tell Airflow it's behind a proxy
[webserver]
base_url              = https://airflow.example.com
enable_proxy_fix      = True
proxy_fix_x_for       = 1
proxy_fix_x_proto     = 1
proxy_fix_x_host      = 1
proxy_fix_x_prefix    = 1
```

---

## API Security

```properties
# airflow.cfg
[api]
# Require authentication for REST API (Airflow 2.0+)
auth_backends = airflow.api.auth.backend.basic_auth

# Or JWT-based
auth_backends = airflow.api.auth.backend.kerberos_auth
```

```bash
# API with basic auth
curl -X GET "http://localhost:8080/api/v1/dags" \
  -u admin:admin \
  -H "Content-Type: application/json"

# API with JWT token
TOKEN=$(curl -X POST "http://localhost:8080/auth/token" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' | jq -r '.access_token')

curl -X GET "http://localhost:8080/api/v1/dags" \
  -H "Authorization: Bearer $TOKEN"
```

---

## Security Checklist

```
Authentication:
  ☐ Disable default admin/admin credentials
  ☐ Enable OAuth or LDAP for enterprise deployments
  ☐ Enforce strong password policy

Encryption:
  ☐ Set a strong FERNET_KEY
  ☐ Rotate Fernet key periodically
  ☐ Enable TLS on webserver
  ☐ Enable TLS between Airflow components and DB

Secrets management:
  ☐ Use a secrets backend (Vault, AWS SM, GCP SM)
  ☐ Never hardcode credentials in DAG files
  ☐ Never store secrets in Airflow Variables (plaintext in UI)

Access control:
  ☐ Apply least-privilege roles to all users
  ☐ Use DAG-level access_control for sensitive pipelines
  ☐ Audit user access regularly

Network:
  ☐ Webserver behind reverse proxy with TLS
  ☐ Restrict metadata DB access to Airflow services only
  ☐ Enable API authentication
  ☐ Put Airflow in a private network, expose only webserver

Code security:
  ☐ Review DAGs before merging (no shell injection)
  ☐ Validate BashOperator commands
  ☐ Avoid dynamic SQL without parameterization
```
