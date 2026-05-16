# Chapter 6: Persistence

## Persistence Overview

Redis is in-memory but offers two persistence mechanisms:

```
RDB (Redis Database):          AOF (Append Only File):
  Point-in-time snapshots        Log of every write command
  Compact binary file            Fully reproducible
  Fast restart                   Near-zero data loss
  Some data loss possible        Larger files
  Best for backups               Best for durability
```

**Recommendation:** Use **both** (hybrid) in production.

---

## RDB (Snapshotting)

RDB saves a binary snapshot of all data at intervals.

### Configuration

```properties
# redis.conf

# Auto-save rules: save after N seconds if M keys changed
save 900 1          # 1 key changed in 15 minutes
save 300 10         # 10 keys changed in 5 minutes
save 60 10000       # 10000 keys changed in 1 minute

# Disable auto-save
# save ""

# RDB file settings
dir /var/lib/redis
dbfilename dump.rdb
rdbcompression yes       # compress with LZF
rdbchecksum yes          # CRC64 checksum

# Stop accepting writes if RDB save fails (data safety guard)
stop-writes-on-bgsave-error yes
```

### Manual Save

```bash
# Non-blocking background save
BGSAVE

# Blocking save (blocks all clients — avoid in production)
SAVE

# Check last save time
LASTSAVE    # Unix timestamp

# Save status
INFO persistence
# rdb_changes_since_last_save: 0
# rdb_last_bgsave_status: ok
# rdb_last_bgsave_time_sec: 2
```

### How BGSAVE Works

```
Redis (parent process)
       │
       │── fork() ──► Child process
       │                   │
       │ (continues         │ Writes snapshot to dump.rdb.tmp
       │  serving clients)  │
       │                   └─► Atomic rename to dump.rdb
       │ (Copy-on-write)
```

The parent uses Copy-on-Write semantics — only modified pages are duplicated. Very efficient for mostly-read workloads.

---

## AOF (Append Only File)

AOF logs every write command. On restart, Redis replays the file to reconstruct state.

### Configuration

```properties
# redis.conf
appendonly yes
appendfilename "appendonly.aof"
dir /var/lib/redis

# Fsync policy
appendfsync always      # safest (every write), slowest
appendfsync everysec    # recommended (1 second window risk)
appendfsync no          # fastest, OS decides (up to 30s risk)

# Prevent AOF rewrite during BGSAVE
no-appendfsync-on-rewrite yes

# Auto-rewrite AOF when it grows by 100% and is > 64MB
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb

# Handle truncated AOF on restart (power failure)
aof-use-rdb-preamble yes   # hybrid format (see below)
aof-load-truncated yes
```

### AOF Rewrite

The AOF file grows over time (e.g., 1M INCR commands = 1M lines). Rewrite compacts it to the current state (just 1 SET command):

```bash
# Manual rewrite
BGREWRITEAOF

# Check status
INFO persistence
# aof_rewrite_in_progress: 0
# aof_last_rewrite_time_sec: 1
# aof_current_size: 52428800
# aof_base_size: 10485760
```

---

## Hybrid Persistence (Recommended)

Redis 4.0+ supports a hybrid AOF format: RDB as the base + AOF for changes since last RDB.

```properties
# redis.conf
appendonly yes
aof-use-rdb-preamble yes    # hybrid = fast load + near-zero data loss
```

**Benefits:**
- Fast restart (RDB is binary, fast to load)
- Near-zero data loss (AOF for recent changes)
- Best of both worlds

---

## Restart and Recovery

```
Startup sequence:
  1. If AOF enabled → load AOF (more recent, better durability)
  2. If only RDB   → load dump.rdb
  3. If neither    → start empty

Recovery after crash:
  RDB: lose up to last-save-interval of data
  AOF everysec: lose up to ~1 second of data
  AOF always: theoretically 0 data loss (1 write per fsync)
```

```bash
# Verify AOF file integrity
redis-check-aof --fix appendonly.aof

# Verify RDB file integrity
redis-check-rdb dump.rdb
```

---

## Backup Strategy

```bash
# 1. Copy RDB file (safe — it's a complete point-in-time snapshot)
cp /var/lib/redis/dump.rdb /backup/redis-$(date +%Y%m%d-%H%M%S).rdb

# 2. Script with BGSAVE
redis-cli BGSAVE
sleep 5
cp /var/lib/redis/dump.rdb /backup/redis-$(date +%Y%m%d).rdb

# 3. Use DUMP/RESTORE for individual keys
redis-cli DUMP user:1000      # serialize key to binary
redis-cli RESTORE user:1000 0 "<binary>"  # restore

# 4. Use redis-cli --rdb for remote backup
redis-cli --rdb /backup/dump.rdb   # triggers SYNC protocol

# 5. Replicate to a replica and take RDB from replica
# (avoids fork() overhead on primary)
```

---

## Persistence Tradeoffs

| Feature | No Persistence | RDB Only | AOF Only | RDB + AOF |
|---|---|---|---|---|
| Data loss on crash | All | Since last save | ~1 second | ~1 second |
| Restart speed | Instant | Fast | Slow (replay) | Fast (hybrid) |
| File size | None | Small (binary) | Large | Medium |
| I/O overhead | None | Periodic | Continuous | Moderate |
| Use case | Cache only | Snapshots | Critical data | Production |

---

## Persistence Monitoring

```bash
INFO persistence

# Key metrics:
# rdb_last_bgsave_status: ok       ← should be "ok"
# rdb_last_bgsave_time_sec: 2      ← how long last RDB took
# aof_enabled: 1
# aof_last_write_status: ok        ← should be "ok"
# aof_current_size: 52428800       ← bytes
# aof_rewrite_in_progress: 0

# Alert if:
# rdb_last_bgsave_status: err       → RDB save failed
# aof_last_write_status: err        → AOF write failed (data loss risk)
```
