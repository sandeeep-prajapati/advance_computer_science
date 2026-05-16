# Chapter 17: Snapshots & Backup

## What Are Snapshots?

A **snapshot** is an incremental backup of one or more Elasticsearch indices (or the entire cluster state) stored in a repository. Snapshots are:

- **Incremental** — only changed segments since last snapshot are copied
- **Non-destructive** — the cluster continues operating during a snapshot
- **Restorable** — indices can be restored to any cluster

---

## Snapshot Repositories

A **repository** is where snapshots are stored. Supported types:

| Type | Backend | Plugin Required |
|---|---|---|
| `fs` | Shared filesystem (NFS, etc.) | No |
| `s3` | Amazon S3 | `repository-s3` |
| `gcs` | Google Cloud Storage | `repository-gcs` |
| `azure` | Azure Blob Storage | `repository-azure` |
| `hdfs` | Hadoop HDFS | `repository-hdfs` |
| `url` | Read-only HTTP/S URL | No |

---

## Register a Repository

### Filesystem Repository (Development)

```bash
# First, add the path to elasticsearch.yml on ALL nodes:
# path.repo: ["/mnt/backups"]

PUT /_snapshot/my_backup_repo
{
  "type": "fs",
  "settings": {
    "location": "/mnt/backups/elasticsearch",
    "compress": true
  }
}
```

### S3 Repository

```bash
# Install plugin first: bin/elasticsearch-plugin install repository-s3

# Store credentials in keystore
bin/elasticsearch-keystore add s3.client.default.access_key
bin/elasticsearch-keystore add s3.client.default.secret_key

PUT /_snapshot/s3_backup
{
  "type": "s3",
  "settings": {
    "bucket":   "my-es-backups",
    "region":   "us-east-1",
    "base_path": "elasticsearch/snapshots",
    "compress": true
  }
}
```

### Verify Repository

```bash
POST /_snapshot/my_backup_repo/_verify
```

---

## Taking Snapshots

### Manual Snapshot

```bash
# Snapshot all indices
PUT /_snapshot/my_backup_repo/snapshot_2024_03_15
{
  "indices": "*",
  "ignore_unavailable": true,
  "include_global_state": false
}

# Snapshot specific indices
PUT /_snapshot/my_backup_repo/products_snapshot
{
  "indices": ["products", "products_v2"],
  "ignore_unavailable": true,
  "include_global_state": false,
  "metadata": {
    "taken_by":   "ops-team",
    "reason":     "pre-migration backup"
  }
}
```

### Check Snapshot Status

```bash
# Status of in-progress snapshot
GET /_snapshot/my_backup_repo/snapshot_2024_03_15/_status

# Details of completed snapshot
GET /_snapshot/my_backup_repo/snapshot_2024_03_15
```

### List All Snapshots

```bash
GET /_snapshot/my_backup_repo/*?verbose=false
```

---

## Snapshot Lifecycle Management (SLM)

Automate snapshot creation and retention:

```bash
PUT /_slm/policy/daily_snapshots
{
  "schedule":   "0 30 1 * * ?",    ← cron: 1:30 AM daily
  "name":       "<daily-snap-{now/d}>",
  "repository": "s3_backup",
  "config": {
    "indices":              ["*"],
    "ignore_unavailable":   true,
    "include_global_state": false
  },
  "retention": {
    "expire_after":   "30d",    ← delete snapshots older than 30 days
    "min_count":      5,        ← always keep at least 5
    "max_count":      50        ← never keep more than 50
  }
}
```

### Execute SLM Policy Immediately

```bash
POST /_slm/policy/daily_snapshots/_execute
```

### Check SLM Status

```bash
GET /_slm/policy/daily_snapshots

GET /_slm/stats
```

---

## Restoring Snapshots

### Restore All Indices

```bash
POST /_snapshot/my_backup_repo/snapshot_2024_03_15/_restore
{
  "ignore_unavailable":   true,
  "include_global_state": false
}
```

### Restore Specific Indices

```bash
POST /_snapshot/my_backup_repo/snapshot_2024_03_15/_restore
{
  "indices": ["products", "orders"],
  "ignore_unavailable": true
}
```

### Restore with Rename (to avoid conflicts with existing indices)

```bash
POST /_snapshot/my_backup_repo/snapshot_2024_03_15/_restore
{
  "indices":             ["products"],
  "rename_pattern":      "(.+)",
  "rename_replacement":  "restored_$1",
  "index_settings": {
    "index.number_of_replicas": 0
  }
}
```

This restores `products` as `restored_products`.

### Monitor Restore Progress

```bash
GET /_cat/recovery?v&active_only=true
GET /products/_recovery
```

---

## Searchable Snapshots

Searchable snapshots allow you to **search index data directly from the snapshot repository** without fully restoring it. Ideal for cold/frozen data.

```bash
POST /_snapshot/s3_backup/snapshot_2024_03_15/_mount
{
  "index":           "logs-2024-01",
  "renamed_index":   "logs-2024-01-frozen",
  "storage":         "full_copy"    ← or "shared_cache" for frozen tier
}
```

- `full_copy` — downloads all data to local storage (warm tier)
- `shared_cache` — fetches data on demand from repository (frozen tier)

---

## Cross-Cluster Restore

Restore from a snapshot taken on a different cluster:

```bash
# Register the remote cluster's repository
PUT /_snapshot/remote_backup
{
  "type": "s3",
  "settings": {
    "bucket":   "remote-cluster-backups",
    "readonly": true    ← read-only to prevent accidental writes
  }
}

# Restore
POST /_snapshot/remote_backup/snapshot_2024_03_15/_restore
{
  "indices": ["products"]
}
```

---

## Backup Strategy

### Recommended Production Backup Schedule

| Snapshot Type | Frequency | Retention |
|---|---|---|
| Daily | Every day at 2 AM | 30 days |
| Weekly | Every Sunday | 12 weeks |
| Monthly | 1st of each month | 12 months |

### What `include_global_state` Captures

When `true`, snapshots also include:
- Index templates
- Component templates
- ILM policies
- SLM policies
- Feature state (security users/roles, etc.)

Always set to `false` unless you're doing a full cluster migration.

---

## Delete Snapshots

```bash
# Delete a single snapshot
DELETE /_snapshot/my_backup_repo/snapshot_2024_03_15

# Delete multiple
DELETE /_snapshot/my_backup_repo/snapshot_2024_03_15,snapshot_2024_03_14

# Delete all in repo (careful!)
DELETE /_snapshot/my_backup_repo/*
```
