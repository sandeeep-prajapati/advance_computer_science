# Chapter 10: Cluster Architecture

## Overview

An Elasticsearch **cluster** is a group of one or more nodes that collectively hold all your data and provide indexing and search capabilities.

```
Cluster: "production-cluster"
├── Node 1  (master-eligible, data)
├── Node 2  (master-eligible, data)
├── Node 3  (master-eligible, data)
├── Node 4  (data)
├── Node 5  (data)
└── Node 6  (coordinating-only / ingest)
```

---

## Node Roles

Each node can serve one or more roles. In small clusters, a single node does everything. In production, roles are separated for isolation and performance.

### Master Node

Manages cluster-wide metadata:
- Index creation and deletion
- Shard allocation decisions
- Tracking which nodes are in the cluster

```yaml
# elasticsearch.yml
node.roles: [master]
```

> Recommendation: 3 dedicated master-eligible nodes (odd number to prevent split-brain).

### Data Node

Stores shards and executes index/search operations:

```yaml
node.roles: [data]
# Or tier-specific:
node.roles: [data_hot]    # Fast SSDs — active indices
node.roles: [data_warm]   # Cheaper storage — read-mostly indices
node.roles: [data_cold]   # Object storage — rare access
node.roles: [data_frozen] # Searchable snapshot — archive
```

### Ingest Node

Runs **ingest pipelines** to pre-process documents before indexing (parse, enrich, transform):

```yaml
node.roles: [ingest]
```

### Coordinating Node (Coordinating-Only)

Routes client requests to the correct data nodes, collects and merges results:

```yaml
node.roles: []    ← empty = coordinating only
```

A dedicated coordinating node offloads the coordination work from data and master nodes in large clusters.

### ML Node

Runs machine learning jobs (anomaly detection, NLP):

```yaml
node.roles: [ml]
```

---

## Master Election

Elasticsearch uses a **Raft-based consensus algorithm** to elect a master node.

1. Nodes with `master` role participate in elections
2. A node wins if it gets votes from a **quorum** (majority) of master-eligible nodes
3. The quorum prevents split-brain:

```
For 3 master-eligible nodes, quorum = 2
For 5 master-eligible nodes, quorum = 3
```

```yaml
# Configure who can be initial master candidates (set once at cluster bootstrap)
cluster.initial_master_nodes: ["node-1", "node-2", "node-3"]
```

> Never use an even number of master-eligible nodes. With 2 nodes, losing one means no quorum, and the cluster goes read-only.

---

## Cluster State

The **cluster state** is the metadata about the entire cluster maintained by the master node:
- Which indices exist
- Their mappings and settings
- Shard allocation (which shard on which node)

The master node broadcasts cluster state changes to all nodes.

---

## Node Discovery

How nodes find each other in the cluster:

```yaml
# Seed hosts (known nodes to contact for discovery)
discovery.seed_hosts: ["192.168.1.10", "192.168.1.11", "192.168.1.12"]

# Or use cloud-native discovery
discovery.ec2.tag.Environment: production
```

---

## Shard Allocation

The master decides which node gets each shard. Factors:
- **Disk usage** — avoid nodes with > 85% disk (watermark.low)
- **Shard count balance** — spread shards evenly
- **Rack/zone awareness** — keep primaries and replicas in different failure domains

### Zone Awareness (High Availability)

```yaml
# On each node, set its zone
node.attr.zone: us-east-1a

# In index settings
PUT /products/_settings
{
  "index.routing.allocation.awareness.attributes": "zone"
}

# Globally require one replica per zone
PUT /_cluster/settings
{
  "persistent": {
    "cluster.routing.allocation.awareness.attributes":       "zone",
    "cluster.routing.allocation.awareness.force.zone.values": "us-east-1a,us-east-1b,us-east-1c"
  }
}
```

### Controlling Shard Placement

```bash
# Pin an index to nodes with specific attributes
PUT /products/_settings
{
  "index.routing.allocation.require.box_type": "hot"
}

# Exclude a specific node from an index
PUT /products/_settings
{
  "index.routing.allocation.exclude._name": "node-3"
}
```

---

## Disk Watermarks

Elasticsearch automatically stops allocating new shards when nodes get full:

```bash
PUT /_cluster/settings
{
  "persistent": {
    "cluster.routing.allocation.disk.watermark.low":            "85%",   ← no new shards
    "cluster.routing.allocation.disk.watermark.high":           "90%",   ← start relocating
    "cluster.routing.allocation.disk.watermark.flood_stage":    "95%"    ← read-only
  }
}
```

---

## Cross-Cluster Search (CCS)

Query data across multiple independent clusters:

```yaml
# Configure remote cluster
PUT /_cluster/settings
{
  "persistent": {
    "cluster.remote.cluster_b.seeds": ["192.168.2.10:9300"]
  }
}
```

```bash
# Search across clusters
GET /cluster_b:products,local_products/_search
{
  "query": { "match_all": {} }
}
```

---

## Cross-Cluster Replication (CCR)

Replicate indices from a leader cluster to a follower cluster for geo-distribution or DR:

```bash
# On the follower cluster
PUT /replicated-products/_ccr/follow
{
  "remote_cluster":  "leader-cluster",
  "leader_index":    "products"
}
```

---

## Cluster Settings

```bash
# View all cluster settings
GET /_cluster/settings?include_defaults=true

# Update cluster settings
PUT /_cluster/settings
{
  "persistent": {
    "cluster.max_shards_per_node": 2000
  },
  "transient": {
    "indices.recovery.max_bytes_per_sec": "100mb"
  }
}
```

| Type | Persistence |
|---|---|
| `persistent` | Survives cluster restart |
| `transient` | Cleared on cluster restart |

---

## Architecture Patterns

### Single Node (Development)
```
[Node 1: master + data + ingest + coordinating]
```

### Small Production (3 nodes)
```
Node 1: master-eligible + data
Node 2: master-eligible + data
Node 3: master-eligible + data
```

### Medium Production (7 nodes)
```
Node 1-3: dedicated master (no data)
Node 4-7: dedicated data (hot tier)
```

### Large Production (tiered)
```
Node 1-3: dedicated master
Node 4-6: data_hot    (NVMe SSDs)
Node 7-9: data_warm   (SAS/SATA)
Node 10:  data_cold   (S3-backed)
Node 11:  ingest
Node 12:  coordinating
```

---

## Cluster Diagnostics

```bash
# Overall health
GET /_cluster/health

# Allocation explain (why is a shard unassigned?)
GET /_cluster/allocation/explain
{
  "index":   "products",
  "shard":   0,
  "primary": true
}

# Pending cluster tasks
GET /_cluster/pending_tasks

# Hot threads (what are nodes doing?)
GET /_nodes/hot_threads

# Node stats
GET /_nodes/stats

# Task management
GET /_tasks?detailed=true&actions=*search*
```
