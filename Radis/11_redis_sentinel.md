# Chapter 11: Redis Sentinel

## What is Redis Sentinel?

Redis Sentinel provides **high availability** for Redis without data sharding. It monitors Redis instances, detects failures, and performs **automatic failover** — promoting a replica to primary without human intervention.

```
┌─────────────────────────────────────────────────────────┐
│                   Sentinel Cluster                       │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐            │
│  │Sentinel 1│   │Sentinel 2│   │Sentinel 3│            │
│  └──────────┘   └──────────┘   └──────────┘            │
│       │               │               │                 │
│       └───────────────┼───────────────┘                 │
│                       │ monitor + vote                  │
└───────────────────────┼─────────────────────────────────┘
                        │
         ┌──────────────┼──────────────┐
         │              │              │
    ┌────▼────┐    ┌────▼────┐    ┌────▼────┐
    │ Primary │    │Replica 1│    │Replica 2│
    │ :6379   │    │  :6380  │    │  :6381  │
    └─────────┘    └─────────┘    └─────────┘
```

---

## Sentinel Configuration

```properties
# sentinel.conf (same file for all sentinel instances)

# Monitor a primary (name, host, port, quorum)
sentinel monitor myredis 127.0.0.1 6379 2
# quorum=2: need 2 sentinels to agree on failure

# Authentication
sentinel auth-pass myredis yourPassword

# Failure detection
sentinel down-after-milliseconds myredis 5000    # ms before "subjectively down"
sentinel failover-timeout myredis 60000          # max failover time

# How many replicas to reconfigure simultaneously during failover
sentinel parallel-syncs myredis 1

# Sentinel port
port 26379
daemonize yes
logfile /var/log/redis/sentinel.log
```

---

## Docker Compose: Full HA Setup

```yaml
version: "3.8"
services:
  redis-primary:
    image: redis:7.2
    command: redis-server --requirepass secret
    ports: ["6379:6379"]

  redis-replica1:
    image: redis:7.2
    command: redis-server --replicaof redis-primary 6379 --masterauth secret
    ports: ["6380:6380"]

  redis-replica2:
    image: redis:7.2
    command: redis-server --replicaof redis-primary 6379 --masterauth secret
    ports: ["6381:6381"]

  sentinel1:
    image: redis:7.2
    command: >
      sh -c 'echo "sentinel monitor myredis redis-primary 6379 2
      sentinel auth-pass myredis secret
      sentinel down-after-milliseconds myredis 5000
      sentinel failover-timeout myredis 60000
      sentinel parallel-syncs myredis 1" > /sentinel.conf &&
      redis-sentinel /sentinel.conf'
    ports: ["26379:26379"]

  sentinel2:
    image: redis:7.2
    command: >
      sh -c 'echo "sentinel monitor myredis redis-primary 6379 2
      sentinel auth-pass myredis secret
      sentinel down-after-milliseconds myredis 5000
      sentinel failover-timeout myredis 60000
      sentinel parallel-syncs myredis 1" > /sentinel.conf &&
      redis-sentinel /sentinel.conf'
    ports: ["26380:26380"]

  sentinel3:
    image: redis:7.2
    command: >
      sh -c 'echo "sentinel monitor myredis redis-primary 6379 2
      sentinel auth-pass myredis secret
      sentinel down-after-milliseconds myredis 5000
      sentinel failover-timeout myredis 60000
      sentinel parallel-syncs myredis 1" > /sentinel.conf &&
      redis-sentinel /sentinel.conf'
    ports: ["26381:26381"]
```

---

## Connecting via Sentinel (Python)

```python
from redis.sentinel import Sentinel

sentinel = Sentinel(
    [("sentinel1", 26379), ("sentinel2", 26380), ("sentinel3", 26381)],
    socket_timeout = 0.1,
    password       = "secret",
)

# Get primary (write)
primary = sentinel.master_for("myredis", socket_timeout=0.1, password="secret")

# Get replica (read)
replica = sentinel.slave_for("myredis", socket_timeout=0.1, password="secret")

# Use normally
primary.set("key", "value")
value = replica.get("key")
print(value)
```

If a failover happens, the Sentinel client automatically discovers the new primary on the next operation.

---

## Failover Process

### Step-by-Step

```
1. Primary stops responding
2. After down-after-milliseconds, each sentinel marks it "subjectively down" (SDOWN)
3. When quorum sentinels agree → "objectively down" (ODOWN)
4. Sentinels elect a leader sentinel (via Raft-like voting)
5. Leader sentinel selects best replica to promote:
   - Lowest replica-priority
   - Most recent replication offset (least lag)
   - Lowest run ID (tie-breaker)
6. Leader sends REPLICAOF NO ONE to chosen replica
7. New primary is announced
8. Other replicas point to new primary
9. Old primary reconfigured as replica when it comes back
```

### Sentinel CLI Commands

```bash
# Connect to sentinel
redis-cli -p 26379

# Info about monitored primaries
SENTINEL masters
SENTINEL master myredis

# Info about replicas
SENTINEL replicas myredis

# Info about other sentinels
SENTINEL sentinels myredis

# Get current primary address
SENTINEL get-master-addr-by-name myredis
# Returns: 127.0.0.1, 6379

# Force a failover (test)
SENTINEL failover myredis

# Reset a primary (clear history)
SENTINEL reset myredis
```

---

## Sentinel Notifications and Events

```bash
# Subscribe to sentinel events
redis-cli -p 26379 SUBSCRIBE +switch-master

# Event types:
# +switch-master: failover completed
# +odown:         primary objectively down
# +sdown:         instance subjectively down
# -sdown:         instance recovered
# +reboot:        instance restarted
# +slave:         new replica detected
```

```python
# Python: react to failover events
import redis

sentinel_r = redis.Redis(host="sentinel1", port=26379)
p = sentinel_r.pubsub()
p.subscribe("+switch-master")

for msg in p.listen():
    if msg["type"] == "message":
        print(f"Failover: {msg['data']}")
        # Notify ops team, clear connection pools, etc.
```

---

## Sentinel vs Cluster

| Feature | Sentinel | Cluster |
|---|---|---|
| Purpose | HA for single primary | HA + horizontal sharding |
| Data distribution | All data on primary | Data split across nodes |
| Max memory | Primary RAM limit | Sum of all node RAM |
| Failover | Automatic | Automatic |
| Client complexity | Sentinel-aware client | Cluster-aware client |
| Setup complexity | Medium | High |
| Use case | < 100 GB, one primary | > 100 GB, need sharding |
