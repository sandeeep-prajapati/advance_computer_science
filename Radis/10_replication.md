# Chapter 10: Replication

## Replication Overview

Redis replication is **asynchronous** and **leader-based**. One primary (leader) handles all writes; one or more replicas copy the data stream.

```
Primary (Read/Write)
    │
    ├──► Replica 1 (Read-only)
    ├──► Replica 2 (Read-only)
    └──► Replica 3 (Read-only)
              │
              └──► Replica 4 (cascading replica)
```

**Key properties:**
- Replicas are read-only by default
- Replication is asynchronous (minimal latency on primary)
- Partial resync after short disconnections
- Full resync (RDB transfer) on first connection or long disconnect

---

## Setting Up Replication

### Primary (no special config needed)

```properties
# primary/redis.conf
bind 0.0.0.0
port 6379
requirepass primarypassword
```

### Replica

```properties
# replica/redis.conf
port 6380
replicaof primary-host 6379          # or: slaveof (legacy)
masterauth primarypassword           # auth with primary
replica-read-only yes                # default: replicas reject writes

# Announce itself to Sentinel
replica-announce-ip replica-host
replica-announce-port 6380
```

### Docker Compose Example

```yaml
version: "3.8"
services:
  primary:
    image: redis:7.2
    command: redis-server --requirepass secret
    ports: ["6379:6379"]

  replica1:
    image: redis:7.2
    command: redis-server --replicaof primary 6379 --masterauth secret
    ports: ["6380:6380"]
    depends_on: [primary]

  replica2:
    image: redis:7.2
    command: redis-server --replicaof primary 6379 --masterauth secret
    ports: ["6381:6381"]
    depends_on: [primary]
```

---

## CLI Replication Management

```bash
# Set up replication at runtime
redis-cli -p 6380 REPLICAOF primary-host 6379
redis-cli -p 6380 REPLICAOF NO ONE      # promote to primary

# Check replication status
redis-cli INFO replication

# Primary output:
# role:master
# connected_slaves:2
# slave0:ip=192.168.1.2,port=6380,state=online,offset=12345,lag=0
# slave1:ip=192.168.1.3,port=6381,state=online,offset=12340,lag=1
# master_replid:abc123...
# master_repl_offset:12345
# repl_backlog_active:1
# repl_backlog_size:1048576

# Replica output:
# role:slave
# master_host:primary
# master_port:6379
# master_link_status:up
# master_last_io_seconds_ago:1
# master_sync_in_progress:0
# slave_repl_offset:12345
# slave_priority:100
# slave_read_only:1
```

---

## Replication Process

### First-Time Connection (Full Sync)

```
1. Replica connects to primary
2. Primary forks and creates RDB snapshot (BGSAVE)
3. Primary sends RDB file to replica
4. During transfer: new commands buffered in replication backlog
5. Replica loads RDB
6. Primary sends buffered commands
7. Replica now in sync — receives live command stream
```

### Reconnection (Partial Sync)

```
Replica disconnects briefly:
1. Primary stores recent commands in replication backlog (default 1MB)
2. Replica reconnects with its offset + replid
3. Primary checks: "is offset in my backlog?"
   - YES: send only missing commands (partial sync, fast)
   - NO:  full sync required (expensive)
```

```properties
# Increase backlog for high-write environments
repl-backlog-size 64mb
repl-backlog-ttl 3600    # keep backlog 1 hour after last replica disconnects
```

---

## Read Scaling with Replicas

Route reads to replicas to reduce primary load:

```python
import redis

# Connection pool for primary (writes)
primary = redis.Redis(host="primary-host", port=6379)

# Connection pools for replicas (reads)
replicas = [
    redis.Redis(host="replica1-host", port=6380),
    redis.Redis(host="replica2-host", port=6381),
]

import random

def get_replica() -> redis.Redis:
    return random.choice(replicas)

# Write to primary
primary.set("user:1000:name", "Alice")

# Read from replica
name = get_replica().get("user:1000:name")
```

**Warning:** Replicas are asynchronously updated. There's a small window where reads from replicas may return stale data (replication lag).

---

## Replication Lag

```bash
# Check lag per replica (on primary)
INFO replication
# slave0:ip=...,port=6380,state=online,offset=12345,lag=0
# lag = 0: in sync
# lag = 1: 1 second behind (normal)
# lag > 5: investigate (slow network, overloaded replica)
```

```properties
# Control lag tolerance
min-replicas-to-write 1      # primary requires 1 replica in sync before accepting writes
min-replicas-max-lag 10      # replica must have lag < 10 seconds
```

This prevents data loss: if all replicas are lagging badly, the primary stops accepting writes.

---

## WAIT Command (Sync Confirmation)

Force the primary to wait until replicas acknowledge writes:

```bash
# Wait for 2 replicas to sync, timeout 100ms
WAIT 2 100
# Returns: number of replicas that acknowledged

# In Python
primary.set("critical:key", "value")
synced = primary.wait(2, 100)
if synced < 2:
    raise Exception("Not enough replicas confirmed write")
```

---

## Diskless Replication

Send RDB directly over the network without writing to disk first:

```properties
# Useful when disk I/O is slow but network is fast
repl-diskless-sync yes
repl-diskless-sync-delay 5      # wait 5s for multiple replicas to connect
repl-diskless-sync-max-replicas 0  # 0 = no limit
```

---

## Cascading Replicas

Replicas can themselves have replicas, reducing primary load:

```
Primary ──► Replica 1 ──► Sub-Replica A
                      ──► Sub-Replica B
         ──► Replica 2
```

```properties
# On sub-replica:
replicaof replica1-host 6380
```

Useful for: geographically distributed setups, offloading RDB generation from the primary.
