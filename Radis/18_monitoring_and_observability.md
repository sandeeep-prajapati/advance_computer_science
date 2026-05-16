# Chapter 18: Monitoring & Observability

## INFO Command

The `INFO` command is the primary window into Redis's health and performance:

```bash
# All sections
INFO

# Specific sections
INFO server       # Redis version, OS, uptime
INFO clients      # connected clients, blocked clients
INFO memory       # memory usage, fragmentation
INFO persistence  # RDB/AOF status
INFO stats        # commands/sec, hits/misses, evictions
INFO replication  # role, replicas, replication offset
INFO cpu          # CPU usage
INFO keyspace     # keys per database
INFO latencystats # command latency percentiles
INFO commandstats # per-command statistics
```

### Key Metrics Reference

```bash
INFO all

# Memory
used_memory:104857600             # bytes used
used_memory_human:100.00M
used_memory_peak:209715200        # peak usage
mem_fragmentation_ratio:1.21      # > 1.5 = fragmentation issue
maxmemory:536870912               # configured limit

# Stats
total_commands_processed:1234567  # lifetime commands
instantaneous_ops_per_sec:12345   # current ops/sec
keyspace_hits:900000              # cache hits
keyspace_misses:100000            # cache misses
# hit_rate = 900000/(900000+100000) = 90%
evicted_keys:500                  # keys removed by eviction policy
expired_keys:10000                # keys expired by TTL
rejected_connections:0            # 0 = good, > 0 = too many clients

# Clients
connected_clients:250
blocked_clients:3                 # BRPOP/BLPOP/WAIT
client_recent_max_input_buffer:20480

# Persistence
rdb_last_bgsave_status:ok         # ok = good
aof_last_write_status:ok

# Replication
role:master
connected_slaves:2
master_repl_offset:123456
```

---

## MONITOR

Stream every command executed on the server (use only for debugging — impacts performance):

```bash
redis-cli MONITOR

# Output:
# 1704067200.123456 [0 127.0.0.1:54321] "SET" "key" "value"
# 1704067200.234567 [0 127.0.0.1:54321] "GET" "key"

# Filter commands by client
redis-cli MONITOR | grep "192.168.1.100"
```

---

## Slow Log

```bash
CONFIG SET slowlog-log-slower-than 1000    # log commands > 1ms (1000 microseconds)
CONFIG SET slowlog-max-len 1000

SLOWLOG GET 25        # last 25 slow commands
SLOWLOG LEN           # total entries
SLOWLOG RESET

# Each entry:
# 1) ID
# 2) timestamp
# 3) duration (microseconds)
# 4) command + args
# 5) client IP:port
# 6) client name
```

---

## LATENCY Framework

Redis tracks latency spikes in specific event categories:

```bash
# Enable latency monitoring
CONFIG SET latency-tracking yes
CONFIG SET latency-tracking-info-percentiles "50 99 99.9"
CONFIG SET latency-monitor-threshold 1   # ms

# View latency history
LATENCY HISTORY event
LATENCY LATEST              # latest latency per event type
LATENCY RESET               # clear history

# Common events:
# command    ← slow commands
# fast-command ← fast commands (< 1ms)
# aof-stat  ← AOF fsync
# rdb-unlink-temp-file ← RDB cleanup
```

---

## MEMORY Commands

```bash
# Per-key memory usage
MEMORY USAGE user:1000:profile        # bytes (includes overhead)
MEMORY USAGE user:1000:profile SAMPLES 0  # exact (slower)

# Server-wide memory stats
MEMORY STATS

# Memory breakdown
MEMORY DOCTOR

# Purge freed memory back to OS
MEMORY PURGE

# Malloc stats
MEMORY MALLOC-STATS
```

---

## Prometheus + Grafana

### redis_exporter

```bash
# Run redis_exporter
docker run -d \
  --name redis-exporter \
  -p 9121:9121 \
  oliver006/redis_exporter \
  --redis.addr=redis://localhost:6379 \
  --redis.password=yourPassword
```

### Prometheus Config

```yaml
# prometheus.yml
scrape_configs:
  - job_name: "redis"
    static_configs:
      - targets: ["redis-exporter:9121"]
```

### Key Prometheus Metrics

| Metric | Alert Condition |
|---|---|
| `redis_up` | == 0 → Redis is down |
| `redis_memory_used_bytes` | > 80% of maxmemory |
| `redis_connected_clients` | > maxclients × 0.8 |
| `redis_keyspace_hits_total` | Hit rate < 90% |
| `redis_evicted_keys_total` | Rate > 0 (unexpected eviction) |
| `redis_replication_offset` | Primary–replica delta > threshold |
| `redis_commands_duration_seconds` | p99 > 10ms |
| `redis_aof_last_write_status` | != ok |
| `redis_rdb_last_save_timestamp` | > 24h ago |

### Sample Alert Rules

```yaml
# prometheus-alerts.yml
groups:
  - name: redis
    rules:
      - alert: RedisDown
        expr: redis_up == 0
        for: 1m
        annotations:
          summary: "Redis instance is down"

      - alert: RedisHighMemory
        expr: redis_memory_used_bytes / redis_config_maxmemory > 0.85
        for: 5m
        annotations:
          summary: "Redis memory usage > 85%"

      - alert: RedisLowHitRate
        expr: |
          redis_keyspace_hits_total /
          (redis_keyspace_hits_total + redis_keyspace_misses_total) < 0.9
        for: 10m
        annotations:
          summary: "Redis cache hit rate < 90%"

      - alert: RedisUnexpectedEvictions
        expr: rate(redis_evicted_keys_total[5m]) > 0
        for: 5m
        annotations:
          summary: "Redis is evicting keys — increase maxmemory?"
```

---

## Keyspace Analysis

```bash
# Count keys matching pattern (safe — uses SCAN)
redis-cli --scan --pattern "session:*" | wc -l

# Count all keys
redis-cli DBSIZE

# Key type distribution
redis-cli INFO keyspace
# db0:keys=10000,expires=5000,avg_ttl=3600

# Find large keys
redis-cli --bigkeys

# Memory usage per key type
redis-cli --memkeys --memkeys-samples 100

# Top-N largest keys
redis-cli --hotkeys    # requires maxmemory-policy = *lfu*
```

---

## DEBUG and Diagnostic Commands

```bash
# Object info
OBJECT ENCODING key           # internal encoding
OBJECT IDLETIME key           # seconds since last access
OBJECT FREQ key               # LFU frequency counter (requires LFU policy)
OBJECT HELP

# Debug sleep (simulate slow server — DEV ONLY)
DEBUG SLEEP 0.1               # sleep 100ms

# Reload all data from RDB/AOF
DEBUG RELOAD

# Cluster-specific diagnostics
CLUSTER INFO
CLUSTER NODES
```
