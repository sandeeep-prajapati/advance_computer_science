# Chapter 19: Redis in Production

## Choosing Your Deployment

| Option | Setup | Use Case |
|---|---|---|
| Single instance | Simplest | Dev, small apps with persistence |
| Replication + Sentinel | Medium | Production HA, < 100GB |
| Redis Cluster | Complex | Large datasets, high write throughput |
| Managed (ElastiCache, MemoryDB, Upstash) | Easiest | Cloud-native, ops-free |

---

## Cloud Options

### AWS ElastiCache for Redis

```python
# ElastiCache connection (Cluster mode disabled)
import redis
r = redis.Redis(
    host     = "mycluster.abc123.ng.0001.use1.cache.amazonaws.com",
    port     = 6379,
    ssl      = True,
    ssl_cert_reqs = None,   # ElastiCache uses self-signed cert
)

# ElastiCache Cluster mode
from redis.cluster import RedisCluster, ClusterNode
rc = RedisCluster(
    startup_nodes = [ClusterNode("mycluster.abc.clustercfg.use1.cache.amazonaws.com", 6379)],
    ssl = True,
)
```

### AWS MemoryDB for Redis

MemoryDB persists every transaction to a Multi-AZ log before acknowledging — stronger durability than ElastiCache:

```python
r = redis.Redis(
    host = "clustercfg.my-memorydb.abc.memorydb.us-east-1.amazonaws.com",
    port = 6379,
    ssl  = True,
    password = "AUTH_TOKEN",
)
```

### Upstash (Serverless Redis)

```python
r = redis.Redis(
    host     = "us1-abc.upstash.io",
    port     = 6380,
    password = "your-token",
    ssl      = True,
)
```

---

## Production Checklist

| Category | Setting | Recommended Value |
|---|---|---|
| Auth | `requirepass` | Strong password (or ACL users) |
| Memory | `maxmemory` | 70–80% of available RAM |
| Eviction | `maxmemory-policy` | `allkeys-lru` (cache) / `noeviction` (DB) |
| Persistence | `appendonly` | `yes` |
| Persistence | `aof-use-rdb-preamble` | `yes` (hybrid) |
| Network | `bind` | Internal IP only |
| Security | `protected-mode` | `yes` |
| Danger | `rename-command FLUSHALL` | `""` |
| Danger | `rename-command DEBUG` | `""` |
| Replication | `min-replicas-to-write` | `1` |
| Replication | `min-replicas-max-lag` | `10` |
| OS | THP | Disabled |
| OS | `vm.swappiness` | `0` or `1` |
| Logging | `slowlog-log-slower-than` | `10000` (10ms) |
| Monitoring | Prometheus exporter | Running |
| Backup | RDB + AOF | Both enabled |

---

## Capacity Planning

```
Memory needed = (avg_key_size_bytes + avg_value_size_bytes + overhead_per_key)
                × expected_key_count × replication_factor

Key overhead: ~55–70 bytes per key
Value overhead: varies by type

Example: 1M sessions of 500 bytes each, RF=1
  = (20 + 500 + 65) × 1,000,000 = 585 MB

Add 30% overhead buffer: ~760 MB → provision 1 GB maxmemory

Rule of thumb:
  Small values (<100 bytes): ~100–150 bytes per key
  Medium values (~1KB):      ~1.1–1.2 KB per key
  Large values (>10KB):      close to value size + ~100 bytes
```

---

## Disaster Recovery

### Backup Procedure

```bash
#!/bin/bash
# backup-redis.sh

DATE=$(date +%Y%m%d-%H%M%S)
BACKUP_DIR="/backup/redis"
REDIS_DATA="/var/lib/redis"

mkdir -p $BACKUP_DIR

# Trigger a background save
redis-cli BGSAVE

# Wait for save to complete
while [ "$(redis-cli LASTSAVE)" == "$LAST_SAVE" ]; do
    sleep 1
done

# Copy files
cp $REDIS_DATA/dump.rdb     $BACKUP_DIR/dump-$DATE.rdb
cp $REDIS_DATA/appendonly.aof $BACKUP_DIR/aof-$DATE.aof

# Upload to S3
aws s3 cp $BACKUP_DIR/dump-$DATE.rdb s3://my-backups/redis/

# Cleanup old backups (keep 30 days)
find $BACKUP_DIR -mtime +30 -delete

echo "Backup completed: $DATE"
```

### Point-in-Time Restore

```bash
# Stop Redis
systemctl stop redis

# Replace data files
cp /backup/redis/dump-20240101.rdb /var/lib/redis/dump.rdb

# If using AOF:
cp /backup/redis/aof-20240101.aof /var/lib/redis/appendonly.aof

# Start Redis
systemctl start redis

# Verify data
redis-cli DBSIZE
redis-cli INFO keyspace
```

---

## High Availability Deployment

### Sentinel: 3-Node Setup

```
Node 1: Primary + Sentinel (zone a)
Node 2: Replica + Sentinel (zone b)
Node 3: Replica + Sentinel (zone c)

failover quorum = 2
min-replicas-to-write = 1
```

### Cluster: 6-Node Setup

```
Node 1a: Primary (slots 0–5460)    | Zone A
Node 1b: Replica of Node 1a        | Zone B
Node 2a: Primary (slots 5461–10922)| Zone B
Node 2b: Replica of Node 2a        | Zone C
Node 3a: Primary (slots 10923–16383)| Zone C
Node 3b: Replica of Node 3a        | Zone A
```

---

## Troubleshooting Guide

| Symptom | Diagnosis | Fix |
|---|---|---|
| High memory, no evictions | Objects too large / never expiring | Add TTL, reduce value size |
| Low hit rate | Wrong eviction policy or TTL too short | Increase TTL, switch to LFU |
| Slow commands | Missing pipeline / large value / KEYS | Add pipeline, use SCAN |
| OOM crashes | maxmemory too low | Increase memory or reduce data |
| Replication lag | Slow replica or network | Check replica resources, bandwidth |
| Disconnected replica | Backlog overflow | Increase repl-backlog-size |
| Client timeout | Connection pool exhausted | Increase max_connections, add pool |
| Fragmentation > 1.5 | Many write/delete cycles | Enable activedefrag or restart |

---

## Operational Runbooks

### Rolling Upgrade

```bash
# 1. Upgrade replica first (no downtime)
systemctl stop redis-replica
apt-get install redis-server=7.2.x
systemctl start redis-replica

# 2. Verify replica synced
redis-cli -p 6380 INFO replication | grep master_link_status

# 3. Failover: promote replica to primary
redis-cli -p 26379 SENTINEL failover myredis

# 4. Verify new primary
redis-cli -p 26379 SENTINEL get-master-addr-by-name myredis

# 5. Upgrade old primary (now a replica)
apt-get install redis-server=7.2.x
systemctl restart redis
```

### Emergency Memory Recovery

```bash
# 1. Check what's consuming memory
redis-cli --bigkeys
redis-cli MEMORY STATS

# 2. Immediately free memory
redis-cli MEMORY PURGE

# 3. Reduce TTL on large keys
redis-cli EXPIRE large:key:name 300

# 4. Delete unused keys
redis-cli --scan --pattern "tmp:*" | xargs redis-cli DEL

# 5. If critical: temporarily change eviction policy
redis-cli CONFIG SET maxmemory-policy allkeys-lru
```
