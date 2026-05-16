# Chapter 14: Monitoring & Observability

## What to Monitor

```
Cluster Health     →  Is the cluster green/yellow/red?
Node Resources     →  CPU, memory, disk, network
Indexing Rate      →  Documents/sec, bulk rejections
Search Rate        →  Queries/sec, query latency
JVM               →  Heap usage, GC frequency/duration
Shard Health       →  Unassigned shards, recovery progress
```

---

## Cluster Health API

```bash
GET /_cluster/health
```

```json
{
  "cluster_name":             "production",
  "status":                   "green",
  "number_of_nodes":          5,
  "number_of_data_nodes":     3,
  "active_primary_shards":    30,
  "active_shards":            60,
  "relocating_shards":        0,
  "initializing_shards":      0,
  "unassigned_shards":        0,
  "delayed_unassigned_shards": 0,
  "pending_tasks":            0
}
```

| Status | Action Required |
|---|---|
| `green` | All shards assigned and healthy |
| `yellow` | Missing replicas — investigate but not critical |
| `red` | Missing primaries — data may be unavailable — urgent |

```bash
# Health per index
GET /_cluster/health?level=indices

# Health per shard
GET /_cluster/health?level=shards
```

---

## `_cat` APIs (Human-Readable)

The `_cat` APIs output tabular data, great for quick checks:

```bash
# Nodes overview
GET /_cat/nodes?v&h=name,ip,role,heap.percent,cpu,load_1m,disk.avail

# Indices with stats
GET /_cat/indices?v&s=index&h=index,health,status,pri,rep,docs.count,store.size

# Shards
GET /_cat/shards?v&s=index,shard

# Pending recovery
GET /_cat/recovery?v&active_only=true

# Thread pool queues and rejections
GET /_cat/thread_pool?v&h=name,active,queue,rejected

# Pending master tasks
GET /_cat/pending_tasks?v

# Segment info
GET /_cat/segments?v
```

---

## Node Stats API

```bash
GET /_nodes/stats

# Specific subsections
GET /_nodes/stats/jvm,indices,os,process
```

Key sections:

### JVM Stats

```json
"jvm": {
  "mem": {
    "heap_used_in_bytes":    4294967296,
    "heap_used_percent":     53,
    "heap_max_in_bytes":     8589934592
  },
  "gc": {
    "collectors": {
      "young": { "collection_count": 120, "collection_time_in_millis": 3400 },
      "old":   { "collection_count": 2,   "collection_time_in_millis": 800  }
    }
  }
}
```

**Alert thresholds:**
- Heap > 85% → memory pressure
- Old GC time > 30s/hour → memory leak or undersized heap

### Indices Stats

```json
"indices": {
  "indexing": {
    "index_total":         5000000,
    "index_time_in_millis": 12000,
    "index_failed":        0
  },
  "search": {
    "query_total":          250000,
    "query_time_in_millis": 18000,
    "fetch_total":          245000
  }
}
```

---

## Index-Level Stats

```bash
GET /products/_stats

# Specific stats
GET /products/_stats/indexing,search,store
```

---

## Slow Logs

Identify slow queries and indexing:

```bash
PUT /products/_settings
{
  "index.search.slowlog.threshold.query.warn":  "5s",
  "index.search.slowlog.threshold.query.info":  "1s",
  "index.search.slowlog.threshold.query.debug": "500ms",
  "index.search.slowlog.threshold.fetch.warn":  "1s",
  "index.indexing.slowlog.threshold.index.warn": "5s",
  "index.indexing.slowlog.level":               "info"
}
```

Slow logs appear in: `<data-dir>/logs/<cluster-name>_index_search_slowlog.log`

---

## Task Management API

```bash
# List all running tasks
GET /_tasks?detailed=true

# Filter by type
GET /_tasks?actions=*search*&detailed=true

# Cancel a long-running task
POST /_tasks/<task-id>/_cancel
```

---

## Hot Threads

Find what CPU is actually doing:

```bash
GET /_nodes/hot_threads

GET /_nodes/node-1/hot_threads?threads=5&interval=1s
```

---

## Stack Monitoring (Kibana)

For production clusters, use Elastic's built-in monitoring UI:

1. In Kibana: **Stack Management → Stack Monitoring**
2. Or enable Metricbeat to collect ES metrics:

```yaml
# metricbeat.yml
metricbeat.modules:
  - module: elasticsearch
    metricsets:
      - node
      - node_stats
      - cluster_stats
      - index
      - index_recovery
      - shard
    period: 10s
    hosts: ["http://localhost:9200"]
```

---

## Alerting with Kibana Watcher

### Watch Example: Alert on Red Cluster Status

```bash
PUT /_watcher/watch/cluster_health_watch
{
  "trigger": {
    "schedule": { "interval": "1m" }
  },
  "input": {
    "http": {
      "request": {
        "host": "localhost",
        "port": 9200,
        "path": "/_cluster/health"
      }
    }
  },
  "condition": {
    "compare": { "ctx.payload.status": { "eq": "red" } }
  },
  "actions": {
    "send_email": {
      "email": {
        "to":      ["ops-team@company.com"],
        "subject": "ALERT: Elasticsearch cluster is RED",
        "body":    "Cluster {{ctx.payload.cluster_name}} status: {{ctx.payload.status}}"
      }
    }
  }
}
```

---

## Key Metrics to Dashboard

| Metric | Source | Alert Threshold |
|---|---|---|
| Cluster status | `_cluster/health` | != green |
| Heap used % | `_nodes/stats/jvm` | > 85% |
| Unassigned shards | `_cluster/health` | > 0 |
| GC old collection time | `_nodes/stats/jvm` | > 30s/min |
| Search query latency | `_nodes/stats/indices` | > 500ms p95 |
| Indexing failures | `_nodes/stats/indices` | > 0 |
| Disk watermark | `_cat/nodes` | > 80% |
| Thread pool rejections | `_cat/thread_pool` | > 0 |
| CPU usage | `_nodes/stats/os` | > 90% sustained |

---

## Diagnostic Runbook

### Yellow Status (Missing Replicas)

```bash
# 1. Find unassigned shards
GET /_cat/shards?v&h=index,shard,prirep,state,unassigned.reason

# 2. Explain why
GET /_cluster/allocation/explain

# 3. Common causes:
# - Not enough nodes (replicas > data nodes - 1)
# - Disk watermark exceeded
# - Shard filter rules preventing allocation
```

### Red Status (Missing Primary)

```bash
# Emergency: find which indices are red
GET /_cat/indices?v&health=red

# Last resort: force allocate a stale primary (data loss possible!)
POST /_cluster/reroute
{
  "commands": [{
    "allocate_stale_primary": {
      "index":           "products",
      "shard":           0,
      "node":            "node-1",
      "accept_data_loss": true
    }
  }]
}
```

### High Heap Usage

```bash
# 1. Check what's using heap
GET /_nodes/stats/jvm

# 2. Check for large aggregations or deep pagination
GET /_tasks?detailed=true&actions=*search*

# 3. Force GC (use with caution, causes pauses)
POST /_nodes/_all/_gc
```
