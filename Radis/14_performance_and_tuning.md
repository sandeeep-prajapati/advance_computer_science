# Chapter 14: Performance & Tuning

## Pipelining

Pipelining sends multiple commands in a single network round-trip, dramatically reducing latency for bulk operations.

```python
import redis

r = redis.Redis()

# Without pipeline: N round trips
for i in range(1000):
    r.set(f"key:{i}", f"value:{i}")

# With pipeline: 1 round trip
pipe = r.pipeline()
for i in range(1000):
    pipe.set(f"key:{i}", f"value:{i}")
pipe.execute()   # all 1000 commands sent at once

# Pipeline with transactions
pipe = r.pipeline(transaction=True)   # wraps in MULTI/EXEC
pipe.incr("counter")
pipe.get("counter")
result = pipe.execute()

# Benchmark difference
import time

start = time.time()
for i in range(10000):
    r.set(f"k:{i}", "v")
print(f"No pipeline: {time.time()-start:.2f}s")

start = time.time()
pipe = r.pipeline()
for i in range(10000):
    pipe.set(f"k:{i}", "v")
pipe.execute()
print(f"Pipeline: {time.time()-start:.2f}s")
# Pipeline is typically 5-10× faster for bulk ops
```

---

## Connection Pooling

Each Redis connection consumes ~64KB server-side. Reuse connections with a pool:

```python
import redis

# Connection pool (default in redis-py for Redis())
pool = redis.ConnectionPool(
    host          = "localhost",
    port          = 6379,
    max_connections = 50,
    socket_connect_timeout = 2,
    socket_timeout         = 2,
    retry_on_timeout       = True,
)

r = redis.Redis(connection_pool=pool)

# For async code
import redis.asyncio as aioredis

async def main():
    pool = aioredis.ConnectionPool.from_url(
        "redis://localhost:6379",
        max_connections = 20,
    )
    r = aioredis.Redis(connection_pool=pool)
    await r.set("key", "value")
```

---

## Benchmarking

### redis-benchmark

```bash
# Basic benchmark (100k ops, 50 clients, default commands)
redis-benchmark -n 100000 -c 50 -q

# Specific commands
redis-benchmark -n 100000 -t set,get,lpush,lpop -q

# With pipelining
redis-benchmark -n 1000000 -c 100 -P 16 -q

# Large values
redis-benchmark -n 50000 -d 1024 -t set,get -q

# Output example:
# SET: 245098.04 requests per second
# GET: 312500.00 requests per second
# INCR: 338983.05 requests per second
# LPUSH: 256410.25 requests per second
```

---

## Slow Log

Find commands taking too long:

```bash
# Commands slower than 10ms (default: 10000 microseconds)
CONFIG SET slowlog-log-slower-than 10000

# Max entries to keep
CONFIG SET slowlog-max-len 128

# View slow log
SLOWLOG GET 10          # last 10 entries
SLOWLOG GET             # all entries

# Entry format:
# 1) (integer) 14          ← log entry ID
# 2) (integer) 1704067200  ← timestamp
# 3) (integer) 15000       ← duration in microseconds
# 4) 1) "KEYS"             ← command + args
#    2) "user:*"

# Reset slow log
SLOWLOG RESET
SLOWLOG LEN             # current count
```

---

## Memory Optimization

### Key Size

```python
# BAD: long verbose keys
"user_profile_for_customer_id_1000"      # 36 bytes

# GOOD: short but readable
"u:1000:prof"                            # 11 bytes

# Save significant memory at scale:
# 1M keys × 25 bytes saved = 25 MB
```

### Value Encoding

```bash
# Check encoding
OBJECT ENCODING mykey

# Keep small structures in efficient encodings:
# Hash: use HSET for structured data (not JSON string)
# List: keep under 128 items for listpack encoding
# Set: keep under 128 items for listpack/intset encoding

# Configure thresholds
CONFIG SET hash-max-listpack-entries 128
CONFIG SET hash-max-listpack-value 64
CONFIG SET list-max-listpack-size 128
CONFIG SET set-max-intset-entries 512
CONFIG SET zset-max-listpack-entries 128
CONFIG SET zset-max-listpack-value 64
```

### Memory Analysis

```bash
# Total memory usage
INFO memory

# Memory per key (Redis 4.0+)
MEMORY USAGE user:1000:profile   # bytes including overhead

# Find memory hogs
redis-cli --memkeys --memkeys-samples 100

# Memory doctor
MEMORY DOCTOR
# Output: suggestions for reducing memory

# Active defragmentation
CONFIG SET activedefrag yes
CONFIG SET active-defrag-ignore-bytes 100mb
CONFIG SET active-defrag-threshold-lower 10
CONFIG SET active-defrag-threshold-upper 100
```

---

## OS-Level Tuning

```bash
# Disable transparent huge pages (major Redis performance killer)
echo never > /sys/kernel/mm/transparent_hugepage/enabled
echo never > /sys/kernel/mm/transparent_hugepage/defrag

# Add to /etc/rc.local for persistence

# Increase max connections
sysctl -w net.core.somaxconn=65535
sysctl -w net.ipv4.tcp_max_syn_backlog=65535

# TCP keepalive
sysctl -w net.ipv4.tcp_keepalive_time=60

# Disable swap (or set vm.swappiness very low)
sysctl -w vm.swappiness=0

# Increase max open file descriptors
ulimit -n 65535
```

```properties
# redis.conf
tcp-backlog 511       # accept queue size
hz 10                 # timer frequency (10 = default, 100 = faster expiry)
aof-rewrite-incremental-fsync yes
rdb-save-incremental-fsync yes
```

---

## Performance Anti-Patterns

| Anti-Pattern | Problem | Fix |
|---|---|---|
| `KEYS *` in production | Blocks server | Use `SCAN` with MATCH |
| Large values (>1MB) | Slow serialization, memory pressure | Store in S3, reference in Redis |
| Many connections, no pool | Connection overhead | Use connection pool |
| Single-threaded bottleneck | CPU bound | Redis 6+ I/O threads |
| No pipelining for bulk ops | High round-trip latency | Use pipeline() |
| Storing all data forever | Memory grows unbounded | Always set TTL on cache keys |
| SELECT 0..15 for isolation | Logical but not physical | Use separate instances |
| SMEMBERS on huge sets | Returns all members at once | Use SSCAN |
| LRANGE 0 -1 on long lists | Returns whole list | Use pagination |
| DEBUG SLEEP/RELOAD in prod | Blocks all clients | Disable debug commands |

---

## Key Metrics for Performance

```bash
INFO stats

# Important metrics:
# instantaneous_ops_per_sec   ← current throughput
# total_commands_processed    ← total lifetime
# rejected_connections        ← client limit hit (increase maxclients)
# expired_keys                ← TTL expirations
# evicted_keys                ← evictions (increase maxmemory or change policy)
# keyspace_hits               ← cache hits
# keyspace_misses             ← cache misses
# hit_rate = hits / (hits + misses)   ← should be > 95% for a cache
```
