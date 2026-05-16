# Chapter 12: Redis Cluster

## What is Redis Cluster?

Redis Cluster provides **horizontal scaling** and **high availability** simultaneously. Data is automatically sharded across multiple nodes using hash slots.

```
┌───────────────────────────────────────────────────────────┐
│                     Redis Cluster                         │
│                                                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │  Node 1      │  │  Node 2      │  │  Node 3      │   │
│  │  Primary     │  │  Primary     │  │  Primary     │   │
│  │  Slots 0-5460│  │Slots 5461-   │  │Slots 10923-  │   │
│  │              │  │     10922    │  │      16383   │   │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘   │
│         │                 │                 │            │
│  ┌──────▼───────┐  ┌──────▼───────┐  ┌──────▼───────┐   │
│  │  Node 1      │  │  Node 2      │  │  Node 3      │   │
│  │  Replica     │  │  Replica     │  │  Replica     │   │
│  └──────────────┘  └──────────────┘  └──────────────┘   │
└───────────────────────────────────────────────────────────┘
```

---

## Hash Slots

Redis Cluster uses **16,384 hash slots** (0–16383):

```
CRC16(key) % 16384 = slot number

"user:1000" → CRC16 = 42542 → slot 42542 % 16384 = 9774
             → handled by whatever node owns slot 9774
```

### Hash Tags

Force related keys to the same slot using `{}`:

```
{user:1000}:profile  → hash of "user:1000" → always same slot
{user:1000}:sessions → hash of "user:1000" → same slot as above
{user:1000}:orders   → hash of "user:1000" → same slot

Without tags:
  user:1000:profile  → different slot
  user:1000:sessions → different slot
  → Cross-slot MGET/transactions won't work!
```

---

## Create a Cluster

### Using redis-cli

```bash
# Start 6 Redis instances (3 primary + 3 replica)
redis-server --port 7001 --cluster-enabled yes --cluster-config-file nodes-7001.conf --cluster-node-timeout 5000 --daemonize yes
redis-server --port 7002 --cluster-enabled yes --cluster-config-file nodes-7002.conf --cluster-node-timeout 5000 --daemonize yes
redis-server --port 7003 --cluster-enabled yes --cluster-config-file nodes-7003.conf --cluster-node-timeout 5000 --daemonize yes
redis-server --port 7004 --cluster-enabled yes --cluster-config-file nodes-7004.conf --cluster-node-timeout 5000 --daemonize yes
redis-server --port 7005 --cluster-enabled yes --cluster-config-file nodes-7005.conf --cluster-node-timeout 5000 --daemonize yes
redis-server --port 7006 --cluster-enabled yes --cluster-config-file nodes-7006.conf --cluster-node-timeout 5000 --daemonize yes

# Create cluster (--cluster-replicas 1 = 1 replica per primary)
redis-cli --cluster create \
  127.0.0.1:7001 127.0.0.1:7002 127.0.0.1:7003 \
  127.0.0.1:7004 127.0.0.1:7005 127.0.0.1:7006 \
  --cluster-replicas 1
```

### Docker Compose Cluster

```yaml
version: "3.8"
services:
  redis-node-1:
    image: redis:7.2
    command: >
      redis-server
        --port 6379
        --cluster-enabled yes
        --cluster-config-file nodes.conf
        --cluster-node-timeout 5000
        --appendonly yes
    ports: ["7001:6379"]

  redis-node-2:
    image: redis:7.2
    command: >
      redis-server --port 6379 --cluster-enabled yes
        --cluster-config-file nodes.conf --cluster-node-timeout 5000
    ports: ["7002:6379"]
  # ... 4 more nodes ...
```

---

## Cluster CLI Commands

```bash
# Connect to cluster (with -c flag for auto-redirect)
redis-cli -c -p 7001

# Cluster info
CLUSTER INFO

# List all nodes
CLUSTER NODES

# Node for a specific key
CLUSTER KEYSLOT user:1000      # slot number

# Check cluster health
redis-cli --cluster check 127.0.0.1:7001

# Fix cluster issues
redis-cli --cluster fix 127.0.0.1:7001
```

---

## Python Client (redis-py cluster)

```python
from redis.cluster import RedisCluster, ClusterNode

# Connect to cluster
rc = RedisCluster(
    startup_nodes = [
        ClusterNode("localhost", 7001),
        ClusterNode("localhost", 7002),
        ClusterNode("localhost", 7003),
    ],
    decode_responses       = True,
    skip_full_coverage_check = True,
)

# Use exactly like regular redis client
rc.set("user:1000:name", "Alice")
rc.get("user:1000:name")

# Multi-key operations ONLY work with hash tags (same slot)
rc.mset({"{user:1000}:name": "Alice", "{user:1000}:age": "30"})
rc.mget("{user:1000}:name", "{user:1000}:age")

# Pipeline (within one slot)
pipe = rc.pipeline()
pipe.set("{order:abc}:status", "pending")
pipe.set("{order:abc}:amount", "99.99")
pipe.execute()
```

---

## Resharding (Rebalancing Slots)

When adding a new node, move slots to it:

```bash
# Add a new node
redis-cli --cluster add-node new-node:7007 existing-node:7001

# Reshard: move 1000 slots to new node
redis-cli --cluster reshard 127.0.0.1:7001 \
  --cluster-from all \
  --cluster-to <new-node-id> \
  --cluster-slots 1000 \
  --cluster-yes

# Remove a node (must first empty its slots)
redis-cli --cluster reshard 127.0.0.1:7001 \
  --cluster-from <node-to-remove-id> \
  --cluster-to <target-node-id> \
  --cluster-slots <all-slots-of-node>

redis-cli --cluster del-node 127.0.0.1:7001 <node-id-to-remove>

# Auto-rebalance slots evenly
redis-cli --cluster rebalance 127.0.0.1:7001
```

---

## Cluster Configuration

```properties
# redis.conf for each cluster node
cluster-enabled yes
cluster-config-file nodes.conf
cluster-node-timeout 15000           # failover timeout (ms)
cluster-require-full-coverage no     # allow partial operation when some slots unavailable
cluster-migration-barrier 1          # min replicas before migrating to orphaned primary
cluster-allow-reads-when-down no     # reject reads when cluster is down

# Announce IP (important in Docker/K8s)
cluster-announce-ip 10.0.0.1
cluster-announce-port 6379
cluster-announce-bus-port 16379
```

---

## Cluster Limitations

| Limitation | Workaround |
|---|---|
| Only database 0 | Use namespaces via key prefixes |
| Multi-key ops require same slot | Use hash tags `{tag}` |
| No cross-slot transactions | Design keys to use hash tags |
| SCAN iterates one node | Use SCAN on each node separately |
| Pub/Sub limited to node | Use cluster-aware pub/sub |
| Lua scripts: all keys must be same slot | Use hash tags |

---

## Cluster vs Sentinel: When to Choose

```
Use Sentinel when:
  - Data fits in one primary's RAM (< 100–200 GB)
  - You want simple setup
  - You need transparent failover without data sharding

Use Cluster when:
  - Data exceeds single node RAM
  - Need write throughput beyond single node
  - Horizontal scaling required
  - OK with application changes (hash tags, no cross-slot ops)
```
