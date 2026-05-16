# Chapter 13: Security & Authentication

## Security Overview

Elasticsearch 8.x enables security by default. In ES 7.x and earlier, security (X-Pack) was optional.

Security layers:
```
1. Transport Layer Security (TLS)   ← encrypt node-to-node and client-to-node traffic
2. Authentication                   ← verify who the user is
3. Authorization (RBAC)             ← control what they can do
4. Audit Logging                    ← track who did what
5. Field/Document-Level Security    ← restrict which data users can see
```

---

## Enable Security (ES 8.x Default)

In ES 8.x, security is on by default when you start a new cluster. For development you may disable it:

```yaml
# elasticsearch.yml — development only
xpack.security.enabled: false
```

For production, configure properly:

```yaml
xpack.security.enabled: true
xpack.security.transport.ssl.enabled: true
xpack.security.transport.ssl.verification_mode: certificate
xpack.security.transport.ssl.keystore.path: elastic-certificates.p12

xpack.security.http.ssl.enabled: true
xpack.security.http.ssl.keystore.path: http.p12
```

---

## Generate TLS Certificates

```bash
# Generate CA
bin/elasticsearch-certutil ca

# Generate node certificate signed by that CA
bin/elasticsearch-certutil cert --ca elastic-stack-ca.p12

# For HTTP layer
bin/elasticsearch-certutil http

# Copy certificates to all nodes
cp elastic-certificates.p12 /etc/elasticsearch/
```

---

## Built-in Users

Elasticsearch ships with several reserved system users:

| User | Purpose |
|---|---|
| `elastic` | Superuser — full access |
| `kibana_system` | Used by Kibana to connect |
| `logstash_system` | Used by Logstash |
| `beats_system` | Used by Beats |
| `apm_system` | Used by APM Server |
| `remote_monitoring_user` | Monitoring |

```bash
# Set passwords for built-in users
bin/elasticsearch-setup-passwords interactive

# Or auto-generate
bin/elasticsearch-setup-passwords auto
```

---

## User Management

### Create a User

```bash
POST /_security/user/john_doe
{
  "password":  "SecurePass123!",
  "roles":     ["data_analyst"],
  "full_name": "John Doe",
  "email":     "john@example.com"
}
```

### Change Password

```bash
POST /_security/user/john_doe/_password
{
  "password": "NewSecurePass456!"
}
```

### Disable/Enable User

```bash
POST /_security/user/john_doe/_disable
POST /_security/user/john_doe/_enable
```

---

## Role-Based Access Control (RBAC)

### Built-In Roles

| Role | Access Level |
|---|---|
| `superuser` | Full access to everything |
| `kibana_admin` | Full Kibana administration |
| `viewer` | Read-only on all indices/spaces |
| `editor` | Read/write on all indices |
| `monitoring_user` | Read cluster monitoring data |
| `snapshot_user` | Create/restore snapshots |

### Create a Custom Role

```bash
PUT /_security/role/product_reader
{
  "cluster": ["monitor"],
  "indices": [
    {
      "names":      ["products*"],
      "privileges": ["read", "view_index_metadata"]
    }
  ]
}
```

### Full Admin Role for Specific Index

```bash
PUT /_security/role/orders_admin
{
  "cluster": ["monitor", "manage_index_templates"],
  "indices": [
    {
      "names":      ["orders*"],
      "privileges": ["all"]
    }
  ]
}
```

### Index Privileges Reference

| Privilege | Description |
|---|---|
| `read` | Search and GET |
| `write` | Index, update, delete documents |
| `create` | Index new documents only |
| `delete` | Delete documents |
| `index` | Index and update documents |
| `all` | All operations |
| `view_index_metadata` | View mappings, settings, aliases |
| `manage` | Manage index (settings, mappings, aliases) |
| `monitor` | Read index stats |

### Cluster Privileges

| Privilege | Description |
|---|---|
| `monitor` | Read cluster stats/health |
| `manage` | Manage cluster settings |
| `manage_index_templates` | Manage templates |
| `manage_ilm` | Manage ILM policies |
| `manage_security` | Manage users/roles |
| `all` | Everything |

---

## API Keys

For programmatic access without a username/password:

```bash
# Create an API key
POST /_security/api_key
{
  "name":       "my-app-key",
  "expiration": "30d",
  "role_descriptors": {
    "app_role": {
      "cluster": ["monitor"],
      "indices": [
        {
          "names":      ["products*"],
          "privileges": ["read"]
        }
      ]
    }
  }
}
```

Response:
```json
{
  "id":          "abc123",
  "name":        "my-app-key",
  "api_key":     "xyz789secret",
  "encoded":     "base64_of_id:api_key"
}
```

Use in requests:
```bash
curl -H "Authorization: ApiKey base64_of_id:api_key" \
     http://localhost:9200/products/_search
```

### Manage API Keys

```bash
# List your keys
GET /_security/api_key?mine=true

# Invalidate a key
DELETE /_security/api_key
{ "ids": ["abc123"] }
```

---

## Document-Level Security (DLS)

Restrict which documents a role can see:

```bash
PUT /_security/role/us_customer_service
{
  "indices": [
    {
      "names":      ["orders*"],
      "privileges": ["read"],
      "query": {
        "term": { "region": "US" }
      }
    }
  ]
}
```

Users with this role only see orders with `region: "US"`.

---

## Field-Level Security (FLS)

Restrict which fields a role can see:

```bash
PUT /_security/role/read_without_pii
{
  "indices": [
    {
      "names":      ["users*"],
      "privileges": ["read"],
      "field_security": {
        "grant":  ["*"],            ← allow all
        "except": ["ssn", "dob", "payment_info"]   ← except these
      }
    }
  ]
}
```

---

## Audit Logging

Track all access and changes:

```yaml
# elasticsearch.yml
xpack.security.audit.enabled: true
xpack.security.audit.logfile.events.include: [
  "authentication_success",
  "authentication_failed",
  "access_granted",
  "access_denied",
  "tampered_request",
  "connection_denied"
]
```

Audit logs go to `<cluster_name>_audit.json` in the logs directory.

---

## Security Best Practices

| Practice | Details |
|---|---|
| Always use TLS | For both transport and HTTP layers |
| Rotate API keys | Set expiration; rotate regularly |
| Use minimal privileges | Grant only what each service needs |
| Separate service accounts | One API key per service, not shared |
| Disable `elastic` superuser | Create named admin users instead |
| Enable audit logging | Required for compliance |
| Use DLS/FLS for multi-tenancy | Isolate tenants with row/column security |
| Network-level firewall | Don't expose port 9200 to the internet |
