# Chapter 9: Index Management

## Index Settings

Every index has two categories of settings:

| Category | Changeable After Creation? | Examples |
|---|---|---|
| **Static** | No (requires close/reopen) | `number_of_shards`, `codec` |
| **Dynamic** | Yes (live update) | `number_of_replicas`, `refresh_interval` |

```bash
# Create index with explicit settings
PUT /logs-2024-01
{
  "settings": {
    "number_of_shards":   3,
    "number_of_replicas": 1,
    "refresh_interval":   "30s",
    "max_result_window":  50000
  }
}

# Update dynamic settings on existing index
PUT /logs-2024-01/_settings
{
  "number_of_replicas": 2,
  "refresh_interval":   "1s"
}
```

---

## Index Templates

Templates automatically apply settings and mappings to new indices that match a naming pattern.

### Component Templates (Reusable Pieces)

```bash
# Define a reusable mapping component
PUT /_component_template/timestamp_mapping
{
  "template": {
    "mappings": {
      "properties": {
        "@timestamp": { "type": "date" },
        "host":        { "type": "keyword" },
        "log_level":   { "type": "keyword" }
      }
    }
  }
}

# Define reusable settings component
PUT /_component_template/log_settings
{
  "template": {
    "settings": {
      "number_of_shards":   2,
      "number_of_replicas": 1,
      "refresh_interval":   "10s"
    }
  }
}
```

### Index Template (Compose Components)

```bash
PUT /_index_template/logs_template
{
  "index_patterns": ["logs-*", "events-*"],
  "priority":        100,
  "composed_of":     ["timestamp_mapping", "log_settings"],
  "template": {
    "mappings": {
      "properties": {
        "message": { "type": "text" },
        "service": { "type": "keyword" }
      }
    }
  }
}
```

Now any index created with a name matching `logs-*` automatically gets this mapping and settings.

---

## Index Aliases

An **alias** is a virtual name pointing to one or more real indices. Aliases allow zero-downtime reindexing and abstraction over index names.

```bash
# Create alias pointing to one index
POST /_aliases
{
  "actions": [
    { "add": { "index": "products_v1", "alias": "products" } }
  ]
}

# Query through alias (transparent to client)
GET /products/_search

# Switch alias atomically (zero downtime reindex cutover)
POST /_aliases
{
  "actions": [
    { "remove": { "index": "products_v1", "alias": "products" } },
    { "add":    { "index": "products_v2", "alias": "products" } }
  ]
}
```

### Write Index Alias

Point an alias to multiple indices but designate one as the write target:

```bash
POST /_aliases
{
  "actions": [
    { "add": { "index": "logs-2024-01", "alias": "logs", "is_write_index": false } },
    { "add": { "index": "logs-2024-02", "alias": "logs", "is_write_index": false } },
    { "add": { "index": "logs-2024-03", "alias": "logs", "is_write_index": true  } }
  ]
}
```

Now writes to `logs` go to `logs-2024-03`, but searches span all three.

### Filtered Alias

```bash
POST /_aliases
{
  "actions": [{
    "add": {
      "index": "orders",
      "alias": "eu_orders",
      "filter": { "term": { "region": "EU" } }
    }
  }]
}
```

`GET /eu_orders/_search` only returns EU orders.

---

## Index Lifecycle Management (ILM)

ILM automates the lifecycle of time-series indices (logs, metrics, events). Indices move through phases:

```
Hot  →  Warm  →  Cold  →  Frozen  →  Delete
(write)  (read) (rare)  (archive)   (expire)
```

### Create an ILM Policy

```bash
PUT /_ilm/policy/logs_policy
{
  "policy": {
    "phases": {
      "hot": {
        "actions": {
          "rollover": {
            "max_primary_shard_size": "50gb",
            "max_age":                "7d"
          },
          "set_priority": { "priority": 100 }
        }
      },
      "warm": {
        "min_age": "7d",
        "actions": {
          "shrink":      { "number_of_shards": 1 },
          "forcemerge":  { "max_num_segments": 1 },
          "set_priority": { "priority": 50 }
        }
      },
      "cold": {
        "min_age": "30d",
        "actions": {
          "freeze":       {},
          "set_priority": { "priority": 0 }
        }
      },
      "delete": {
        "min_age": "90d",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}
```

### Attach ILM Policy to an Index Template

```bash
PUT /_index_template/logs_template
{
  "index_patterns": ["logs-*"],
  "template": {
    "settings": {
      "index.lifecycle.name":            "logs_policy",
      "index.lifecycle.rollover_alias":  "logs"
    }
  }
}
```

### Bootstrap the First Index

```bash
PUT /logs-000001
{
  "aliases": {
    "logs": { "is_write_index": true }
  }
}
```

ILM will automatically create `logs-000002`, `logs-000003`, etc. on rollover.

---

## Rollover API

Manually trigger a rollover when conditions are met:

```bash
POST /logs/_rollover
{
  "conditions": {
    "max_age":                "1d",
    "max_primary_shard_size": "10gb",
    "max_docs":               1000000
  }
}
```

Creates a new index `logs-000002` and points the write alias to it.

---

## Reindex API

Copy documents from one index to another:

```bash
POST /_reindex
{
  "source": {
    "index": "products_old",
    "query": {
      "term": { "active": true }
    }
  },
  "dest": {
    "index": "products_new",
    "op_type": "create"    ← skip if already exists
  }
}
```

### Reindex with Transform (Script)

```bash
POST /_reindex
{
  "source": { "index": "products" },
  "dest":   { "index": "products_v2" },
  "script": {
    "source": """
      ctx._source.price_cents = (int)(ctx._source.price * 100);
      ctx._source.remove('price');
    """
  }
}
```

### Reindex from Remote Cluster

```bash
POST /_reindex
{
  "source": {
    "remote": {
      "host": "http://remote-cluster:9200"
    },
    "index": "products"
  },
  "dest": { "index": "products_local" }
}
```

---

## Shrink API

Reduce the number of primary shards (index must be read-only and all shards on one node):

```bash
# Mark source as read-only and move shards
PUT /logs-2024-01/_settings
{
  "settings": {
    "index.blocks.write":   true,
    "index.routing.allocation.require._name": "node-1"
  }
}

# Wait for relocation, then shrink
POST /logs-2024-01/_shrink/logs-2024-01-shrunk
{
  "settings": {
    "index.number_of_shards":   1,
    "index.blocks.write":       false
  }
}
```

---

## Force Merge

Merge Lucene segments to improve read performance (typically for read-only indices):

```bash
POST /logs-2024-01/_forcemerge?max_num_segments=1
```

> Only run on read-only (closed/frozen) indices in production — it's resource-intensive.

---

## Open & Close Index

Closing an index frees memory (no searches allowed):

```bash
POST /logs-2024-01/_close
POST /logs-2024-01/_open
```

---

## Index Stats

```bash
# All indices
GET /_cat/indices?v&s=index

# Specific index stats
GET /products/_stats

# Index segments info
GET /products/_segments
```

---

## Useful Index Operations Reference

| Operation | Command |
|---|---|
| Create index | `PUT /<index>` |
| Delete index | `DELETE /<index>` |
| Close index | `POST /<index>/_close` |
| Open index | `POST /<index>/_open` |
| Refresh | `POST /<index>/_refresh` |
| Flush | `POST /<index>/_flush` |
| Force merge | `POST /<index>/_forcemerge` |
| Reindex | `POST /_reindex` |
| Rollover | `POST /<alias>/_rollover` |
| Shrink | `POST /<index>/_shrink/<new>` |
| Get settings | `GET /<index>/_settings` |
| Get mapping | `GET /<index>/_mapping` |
| List aliases | `GET /_aliases` |
