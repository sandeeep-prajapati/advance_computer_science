# Chapter 11: Shards & Replicas

## What Are Shards?

A **shard** is a self-contained Lucene index — the fundamental unit of data storage and search in Elasticsearch. When you create an index, its data is split horizontally across multiple shards.

```
Index: "products" (3 primary shards, 1 replica each)

Node 1          Node 2          Node 3
───────         ───────         ───────
[P0] [R1]       [P1] [R2]       [P2] [R0]

P = Primary shard
R = Replica shard
```

---

## Primary vs Replica Shards

### Primary Shards

- Store the actual data
- Handle all **write** operations first
- Number is fixed at **index creation time** (cannot change without reindexing)
- Documents are routed to a specific primary shard by:

```
shard = hash(_id) % number_of_primary_shards
```

### Replica Shards

- Exact copies of primary shards
- Can handle **read** (search) requests — reduces primary load
- Number can be **changed at any time**
- Provide **high availability**: if a primary fails, a replica is promoted

```bash
# Change replica count
PUT /products/_settings
{
  "number_of_replicas": 2
}
```

---

## Choosing the Right Shard Count

This is one of the most important (and hardest) decisions in Elasticsearch.

### Common Rule of Thumb

- Target **20–50 GB per shard** for time-series data
- Target **5–40 GB per shard** for search-heavy indices
- A shard count too high wastes memory; too low limits parallelism

### Shard Sizing Formula

```
Number of primary shards ≈ total_expected_data_size / target_shard_size

Example:
  Expected index size: 200 GB
  Target shard size:   40 GB
  → 5 primary shards
```

### Over-Sharding Problems

Each shard:
- Has overhead: ~few MB of heap per shard even when empty
- Participates in every query (even empty ones)
- Creates coordination overhead

> Modern rule: **Keep shard count as low as needed.** One large shard is often better than ten small ones.

---

## Shard Allocation

The master node decides which shard lives on which node. Key settings:

```bash
# Maximum shards per node (default 1000)
PUT /_cluster/settings
{
  "persistent": {
    "cluster.max_shards_per_node": 2000
  }
}

# Force relocate shard manually
POST /_cluster/reroute
{
  "commands": [{
    "move": {
      "index": "products",
      "shard": 0,
      "from_node": "node-1",
      "to_node":   "node-3"
    }
  }]
}

# Allocate an unassigned shard
POST /_cluster/reroute
{
  "commands": [{
    "allocate_replica": {
      "index":   "products",
      "shard":   1,
      "node":    "node-2"
    }
  }]
}
```

---

## Shard Rebalancing

ES automatically rebalances shards when:
- Nodes join or leave the cluster
- Disk utilization becomes uneven
- Shard counts become uneven across nodes

Control rebalancing:

```bash
PUT /_cluster/settings
{
  "persistent": {
    "cluster.routing.rebalance.enable":           "all",     ← all|primaries|replicas|none
    "cluster.routing.allocation.enable":          "all",     ← control new allocations
    "indices.recovery.max_bytes_per_sec":         "40mb"     ← recovery speed throttle
  }
}
```

---

## Write Path: How a Document Gets Indexed

```
Client → Coordinating Node
           ↓
     Route to primary shard (based on _id hash)
           ↓
     Primary shard writes document
           ↓
     Replicate to all replica shards (in parallel)
           ↓
     Acknowledge to client (after quorum of replicas confirm)
```

### Write Consistency Settings

```bash
PUT /products/_doc/1?wait_for_active_shards=1
# 1      = only primary must confirm
# 2      = primary + 1 replica
# all    = all replicas
# quorum = majority (default for multi-replica)
```

---

## Read Path: How a Search is Executed

```
Client → Coordinating Node
           ↓
     Broadcast query to one copy of each shard
     (primary or any replica — round-robin)
           ↓
     Each shard executes query locally, returns top N results + scores
           ↓
     Coordinating node merges all results, re-ranks, returns top N
           ↓
     Fetch phase: retrieve _source for final top N documents
```

This two-phase approach (Query + Fetch) is called **scatter-gather**.

---

## Replica Promotion (Node Failure)

```
Normal:          P0(node1) → R0(node2)
                 P1(node2) → R1(node3)

node2 fails:     P0(node1) — healthy
                 P1(MISSING) — R1(node3) gets promoted to P1
                 R0(MISSING) — new replica allocated on node1/node3

After recovery:  P0(node1) → R0(node3)
                 P1(node3) → R1(node1)
```

During replica promotion, cluster is `yellow` (missing replicas). Once new replicas are assigned, it returns to `green`.

---

## Segment Merging

Each shard consists of multiple **Lucene segments**. As you index documents:

1. New segments are created on each refresh (default every 1 second)
2. Background merges combine small segments into larger ones
3. Deleted documents are only actually removed during a merge

Too many small segments = slow searches. Force merge on read-only indices:

```bash
# Merge down to 1 segment (best for frozen/read-only indices)
POST /logs-2024-01/_forcemerge?max_num_segments=1
```

---

## Monitoring Shard Health

```bash
# List all shards with status
GET /_cat/shards?v&s=index,shard,prirep

# Check why a shard is unassigned
GET /_cluster/allocation/explain

# Shard store info
GET /products/_shard_stores

# Recovery status
GET /_cat/recovery?v&active_only=true
```

### Shard States

| State | Meaning |
|---|---|
| `STARTED` | Shard is active and serving requests |
| `RELOCATING` | Being moved to another node |
| `INITIALIZING` | Being recovered/allocated |
| `UNASSIGNED` | Not assigned to any node (problem!) |

---

## Tips and Anti-Patterns

| Anti-Pattern | Problem | Solution |
|---|---|---|
| Too many shards | High memory overhead, slow queries | Merge/shrink, use ILM |
| One giant shard | No parallelism | Split across more shards |
| 0 replicas in prod | Single point of failure | At least 1 replica per primary |
| Changing shard count | Not possible on live index | Plan ahead; use reindex to resize |
| Replicas on same node | No HA benefit | Use zone awareness |
