# Redis — Complete Guide

> All 20 chapters in one file. From fundamentals to production-grade caching and distributed data systems.

---

## Quick Navigation

- [Index & Architecture Overview](#index-architecture-overview)
- [Introduction to Redis](#introduction-to-redis)
- [Setup & Installation](#setup-installation)
- [Core Data Structures](#core-data-structures)
- [Advanced Data Types](#advanced-data-types)
- [Expiry & Eviction](#expiry-eviction)
- [Persistence](#persistence)
- [Pub/Sub & Messaging](#pubsub-messaging)
- [Transactions & Scripting](#transactions-scripting)
- [Caching Patterns](#caching-patterns)
- [Replication](#replication)
- [Redis Sentinel](#redis-sentinel)
- [Redis Cluster](#redis-cluster)
- [Security](#security)
- [Performance & Tuning](#performance-tuning)
- [Redis with Python](#redis-with-python)
- [Redis with Node.js](#redis-with-nodejs)
- [Redis Modules](#redis-modules)
- [Monitoring & Observability](#monitoring-observability)
- [Redis in Production](#redis-in-production)
- [Learning Roadmap](#learning-roadmap)

---

# Redis — Complete Learning Guide

> A comprehensive, hands-on guide to Redis: from core data structures to production-grade caching, streaming, and distributed systems.

---

## Table of Contents

| Chapter | Title | Topics |
|---|---|---|
| 01 | [Introduction](01_introduction.md) | What is Redis, why it exists, use cases, architecture |
| 02 | [Setup & Installation](02_setup_and_installation.md) | Docker, local install, redis-cli, Redis Insight |
| 03 | [Core Data Structures](03_core_data_structures.md) | String, List, Set, Hash, Sorted Set |
| 04 | [Advanced Data Types](04_advanced_data_types.md) | Streams, HyperLogLog, Bitmaps, Geospatial |
| 05 | [Expiry & Eviction](05_expiry_and_eviction.md) | TTL, eviction policies, keyspace notifications |
| 06 | [Persistence](06_persistence.md) | RDB snapshots, AOF, hybrid persistence |
| 07 | [Pub/Sub & Messaging](07_pub_sub_and_messaging.md) | Channels, patterns, Redis Streams |
| 08 | [Transactions & Scripting](08_transactions_and_scripting.md) | MULTI/EXEC, WATCH, Lua scripting |
| 09 | [Caching Patterns](09_caching_patterns.md) | Cache-aside, write-through, TTL, cache stampede |
| 10 | [Replication](10_replication.md) | Leader-replica, replication lag, partial sync |
| 11 | [Redis Sentinel](11_redis_sentinel.md) | High availability, automatic failover |
| 12 | [Redis Cluster](12_redis_cluster.md) | Horizontal scaling, hash slots, resharding |
| 13 | [Security](13_security.md) | AUTH, ACL, TLS, network security |
| 14 | [Performance & Tuning](14_performance_and_tuning.md) | Memory optimization, pipelining, benchmarking |
| 15 | [Redis with Python](15_redis_with_python.md) | redis-py, connection pools, pipelining |
| 16 | [Redis with Node.js](16_redis_with_nodejs.md) | ioredis, node-redis, patterns |
| 17 | [Redis Modules](17_redis_modules.md) | RediSearch, RedisJSON, RedisTimeSeries, RedisBloom |
| 18 | [Monitoring & Observability](18_monitoring_and_observability.md) | INFO, MONITOR, Prometheus, slow log |
| 19 | [Redis in Production](19_redis_in_production.md) | Cloud options, capacity planning, disaster recovery |
| 20 | [Learning Roadmap](20_learning_roadmap.md) | 5-week plan, practice projects, beyond Redis |

---

## How to Read This Guide

- **Beginners**: Chapters 1–5 cover fundamentals. Practice each data structure before moving on.
- **Application developers**: Focus on Chapters 3–5, 7–9, 15–16 for common patterns.
- **Platform engineers**: Chapters 6, 10–14, 18–19 cover persistence, HA, and production operations.
- **Power users**: Chapters 4, 7–8, 17 cover advanced capabilities and Redis modules.
- **All readers**: Chapter 20 provides a structured 5-week hands-on roadmap.

---

## Redis Architecture

```
                    ┌─────────────────────────────────────────┐
                    │              Redis Server                │
                    │                                         │
                    │  ┌────────────────────────────────────┐ │
                    │  │         In-Memory Data Store        │ │
                    │  │                                     │ │
                    │  │  Strings  Lists  Sets  Hashes       │ │
                    │  │  Sorted Sets  Streams  HyperLogLog  │ │
                    │  └────────────────────────────────────┘ │
                    │                                         │
                    │  ┌──────────────┐  ┌─────────────────┐  │
                    │  │  Persistence │  │  Pub/Sub Engine │  │
                    │  │  (RDB / AOF) │  │  (Channels)     │  │
                    │  └──────────────┘  └─────────────────┘  │
                    │                                         │
                    │  ┌──────────────┐  ┌─────────────────┐  │
                    │  │  Replication │  │   Lua Scripting │  │
                    │  │  (Leader→    │  │   & Transactions│  │
                    │  │   Replica)   │  │                 │  │
                    │  └──────────────┘  └─────────────────┘  │
                    └─────────────────────────────────────────┘
                              ▲                  ▲
                    TCP :6379            Cluster Bus :16379
                    Clients              (node-to-node)
```

---

## Core Abstractions

| Concept | Description |
|---|---|
| **Key** | Unique identifier for every value (binary-safe string) |
| **Value** | One of the supported data types |
| **TTL** | Optional expiry time per key |
| **Database** | Logical namespace (0–15, default 0) |
| **Command** | Single-operation unit (atomic) |
| **Pipeline** | Batch of commands in one network round-trip |
| **Transaction** | MULTI/EXEC block — all-or-nothing execution |
| **Lua Script** | Server-side script evaluated atomically |
| **Keyspace** | All keys in a database |


---

# Chapter 1: Introduction to Redis

## What is Redis?

Redis (Remote Dictionary Server) is an **open-source, in-memory data structure store** that can be used as a database, cache, message broker, and streaming engine. Created by Salvatore Sanfilippo in 2009, it is consistently one of the most popular databases in the world.

The key insight: **data lives in RAM**. This makes Redis orders of magnitude faster than disk-based databases for the right workloads.

```
Traditional DB:             Redis:
  App ──► DB ──► Disk       App ──► Redis ──► RAM
  ~1–10ms per query         ~0.1–1ms per command
  Disk I/O bottleneck       Memory bandwidth = speed
```

---

## Why Redis?

| Need | Redis Solution |
|---|---|
| Fast key-value lookup | String + Hash data types |
| Session storage | String with TTL |
| Rate limiting | Atomic INCR + TTL |
| Leaderboards | Sorted Set (ZADD/ZRANGE) |
| Task queues | List (LPUSH/BRPOP) |
| Pub/Sub messaging | Channels |
| Event streaming | Redis Streams |
| Full-text search | RediSearch module |
| Time series | RedisTimeSeries module |
| Bloom filter | RedisBloom module |

---

## Redis vs Other Stores

| Feature | Redis | Memcached | PostgreSQL | MongoDB |
|---|---|---|---|---|
| Data model | Multiple structures | Key-value only | Relational | Document |
| Persistence | Yes (RDB/AOF) | No | Yes | Yes |
| Replication | Yes | No | Yes | Yes |
| Clustering | Yes (native) | No | Extensions | Yes |
| Pub/Sub | Yes | No | Limited | No |
| Transactions | Yes (MULTI/EXEC) | No | Full ACID | Multi-doc |
| Scripting | Lua | No | PL/pgSQL | JavaScript |
| Latency | Sub-ms | Sub-ms | 1–10ms | 1–10ms |
| RAM-bound | Yes | Yes | No | No |

---

## How Redis Works

### Single-Threaded Event Loop

Redis uses a single-threaded I/O event loop (like Node.js):

```
Network requests ──► Event loop ──► Command queue ──► Execute ──► Reply
                                                        │
                                                  Single thread
                                                  No locks needed
                                                  → Atomic commands
```

This means:
- Every command is atomic by design
- No concurrency bugs from parallel mutations
- Very low latency — no thread context switching
- I/O-bound: single CPU core is often sufficient

**Redis 6.0+**: Added multithreading for I/O handling (reads/writes), but command execution remains single-threaded.

---

## Real-World Use Cases

### 1. Caching
```
User requests ──► App ──► Redis cache?
                              │ YES → return instantly
                              │ NO  → query DB → store in Redis → return
```

### 2. Session Storage
```python
# Login: store session
redis.setex(f"session:{token}", 3600, json.dumps(user_data))

# Request: validate session
data = redis.get(f"session:{token}")
if data: return json.loads(data)
else: return 401
```

### 3. Rate Limiting
```python
key = f"rate:{user_id}:{minute}"
count = redis.incr(key)
redis.expire(key, 60)
if count > 100: return 429  # Too Many Requests
```

### 4. Real-Time Leaderboard
```python
# Update score
redis.zadd("game:leaderboard", {"player:42": 9850})

# Top 10
redis.zrevrange("game:leaderboard", 0, 9, withscores=True)
```

### 5. Distributed Lock
```python
# Acquire lock (NX = only if not exists)
acquired = redis.set("lock:resource", worker_id, nx=True, ex=30)
if acquired:
    try:
        do_critical_work()
    finally:
        if redis.get("lock:resource") == worker_id:
            redis.delete("lock:resource")
```

### 6. Message Queue
```python
# Producer
redis.lpush("queue:emails", json.dumps(email_payload))

# Consumer (blocking pop — waits for new items)
while True:
    _, data = redis.brpop("queue:emails", timeout=0)
    send_email(json.loads(data))
```

---

## Redis Ecosystem

```
┌─────────────────────────────────────────────────────────┐
│                    Redis Ecosystem                       │
│                                                         │
│  ┌───────────┐  ┌───────────┐  ┌──────────────────┐   │
│  │RediSearch │  │ RedisJSON │  │ RedisTimeSeries   │   │
│  │(full-text)│  │(JSON docs)│  │(metrics/IoT)      │   │
│  └───────────┘  └───────────┘  └──────────────────┘   │
│  ┌───────────┐  ┌───────────┐  ┌──────────────────┐   │
│  │RedisBloom │  │ RedisGraph│  │ RedisAI           │   │
│  │(prob. DS) │  │(graph DB) │  │(ML models)        │   │
│  └───────────┘  └───────────┘  └──────────────────┘   │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │                   Redis Core                    │   │
│  │  Data Structures, Persistence, Pub/Sub, Streams │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  ┌──────────┐  ┌───────────┐  ┌──────────────────┐    │
│  │ Sentinel │  │  Cluster  │  │ Redis Insight UI  │    │
│  │ (HA)     │  │ (Scaling) │  │ (GUI management)  │    │
│  └──────────┘  └───────────┘  └──────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

---

## Redis Data Model

Everything in Redis is a **key → value** mapping:

```
Key:   "user:1000:profile"   (any binary string, : is convention)
Value: One of:
  - String:      "Alice"
  - List:        ["item1", "item2", "item3"]
  - Set:         {"apple", "banana", "cherry"}
  - Hash:        {name: "Alice", age: "30", city: "NYC"}
  - Sorted Set:  {("Alice", score=9850), ("Bob", score=7200)}
  - Stream:      time-ordered entries with fields
  - HyperLogLog: approximate unique count
  - Bitmap:      bit array
  - Geospatial:  lat/lon indexed set
```

**Key naming conventions:**
```
object-type:id:field
user:1000:profile
session:abc123def
rate:user-42:2024-01-15-14:00
cache:products:category:electronics
```


---

# Chapter 2: Setup & Installation

## Quick Start with Docker

```bash
# Single Redis instance
docker run -d \
  --name redis \
  -p 6379:6379 \
  redis:7.2

# With persistence
docker run -d \
  --name redis \
  -p 6379:6379 \
  -v redis-data:/data \
  redis:7.2 redis-server --appendonly yes

# Connect
docker exec -it redis redis-cli ping
# PONG
```

---

## Docker Compose (Redis + RedisInsight)

```yaml
# docker-compose.yml
version: "3.8"
services:
  redis:
    image: redis:7.2
    container_name: redis
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
      - ./redis.conf:/usr/local/etc/redis/redis.conf
    command: redis-server /usr/local/etc/redis/redis.conf
    restart: unless-stopped

  redis-insight:
    image: redis/redisinsight:latest
    container_name: redis-insight
    ports:
      - "5540:5540"
    depends_on:
      - redis

volumes:
  redis-data:
```

```bash
docker compose up -d
# Redis UI: http://localhost:5540
```

---

## Redis Configuration File (redis.conf)

```properties
# Network
bind 127.0.0.1
port 6379
protected-mode yes

# Authentication
requirepass yourStrongPassword123

# Persistence
save 900 1          # Save if 1 key changed in 900 seconds
save 300 10         # Save if 10 keys changed in 300 seconds
save 60 10000       # Save if 10000 keys changed in 60 seconds
dir /data
dbfilename dump.rdb

# AOF
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec

# Memory
maxmemory 256mb
maxmemory-policy allkeys-lru

# Logging
loglevel notice
logfile "/var/log/redis/redis-server.log"

# Connection limits
maxclients 10000
timeout 300

# Slow log
slowlog-log-slower-than 10000   # microseconds
slowlog-max-len 128
```

---

## Local Installation

### Ubuntu / Debian
```bash
# Official Redis repo
curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" \
  | sudo tee /etc/apt/sources.list.d/redis.list
sudo apt-get update
sudo apt-get install -y redis

sudo systemctl enable redis-server
sudo systemctl start redis-server
redis-cli ping
```

### macOS
```bash
brew install redis
brew services start redis
redis-cli ping
```

### Build from Source
```bash
wget https://download.redis.io/redis-stable.tar.gz
tar xvzf redis-stable.tar.gz
cd redis-stable
make
sudo make install
redis-server redis.conf
```

---

## redis-cli Basics

```bash
# Connect
redis-cli                           # localhost:6379
redis-cli -h myhost -p 6380         # remote
redis-cli -a password               # with auth
redis-cli -n 2                      # database 2

# Basic commands
redis-cli ping
redis-cli info server
redis-cli dbsize                    # number of keys

# Interactive mode
redis-cli
127.0.0.1:6379> SET name "Alice"
OK
127.0.0.1:6379> GET name
"Alice"
127.0.0.1:6379> DEL name
(integer) 1

# Execute single command
redis-cli SET counter 0
redis-cli INCR counter
redis-cli GET counter

# Pipeline from file
redis-cli --pipe < commands.txt

# Monitor (watch all commands in real time)
redis-cli MONITOR

# Benchmark
redis-benchmark -n 100000 -c 50 -q

# Scan all keys (safe, no KEYS * in production)
redis-cli --scan --pattern "user:*"
```

---

## Redis Insight (GUI)

Redis Insight is the official web-based GUI for Redis.

**Features:**
- Browse and edit keys visually
- Run CLI commands
- Analyze memory usage per key
- View slow log
- Profiler (live command stream)
- Pub/Sub explorer

```bash
# Docker
docker run -d -p 5540:5540 redis/redisinsight:latest

# Access: http://localhost:5540
# Add connection: host=localhost, port=6379
```

---

## Verify Installation: First Commands

```bash
redis-cli

# Test connectivity
PING
# PONG

# Store and retrieve
SET greeting "Hello, Redis!"
GET greeting
# "Hello, Redis!"

# With expiry (TTL in seconds)
SETEX session:abc 3600 "user-data"
TTL session:abc
# 3600

# Atomic counter
SET visits 0
INCR visits
INCR visits
INCR visits
GET visits
# "3"

# List of recent items
LPUSH recent:items "item3"
LPUSH recent:items "item2"
LPUSH recent:items "item1"
LRANGE recent:items 0 -1
# 1) "item1"
# 2) "item2"
# 3) "item3"

# Check server info
INFO server
INFO memory
INFO stats
```

---

## Multiple Databases

Redis has 16 logical databases (0–15) by default, all in the same instance:

```bash
# Switch databases
SELECT 0    # default
SELECT 1
SELECT 15

# Copy key across databases
COPY mykey 1    # copy to DB 1

# Move key to another database
MOVE mykey 2

# Flush current database (destructive!)
FLUSHDB

# Flush all databases (very destructive!)
FLUSHALL
```

**Note:** In Redis Cluster mode, only database 0 is available.


---

# Chapter 3: Core Data Structures

## Strings

The most fundamental Redis type. Despite the name, a Redis string can store any binary data: text, integers, serialized JSON, or even binary blobs up to 512 MB.

```bash
# Basic set/get
SET name "Alice"
GET name             # "Alice"

# Set with expiry
SETEX session:abc 3600 "user-data"
SET session:xyz "data" EX 3600   # same thing
SET session:xyz "data" PX 3600000   # milliseconds
SET key val EX 60 NX             # set ONLY if not exists (atomic lock)
SET key val EX 60 XX             # set ONLY if exists

# Integer operations (atomic)
SET counter 10
INCR counter         # 11
INCRBY counter 5     # 16
DECR counter         # 15
DECRBY counter 3     # 12
INCRBYFLOAT price 1.5

# Get and set atomically
GETSET counter 0     # returns old value, sets new
GETDEL key           # get and delete
GETEX key EX 300     # get and reset TTL

# Multiple keys at once
MSET user:1:name "Alice" user:2:name "Bob"
MGET user:1:name user:2:name

# String length
STRLEN name          # 5

# Append
APPEND log "entry1\n"
APPEND log "entry2\n"
```

**Use cases:** counters, session tokens, cached HTML, feature flags, distributed IDs.

---

## Lists

Ordered collections of strings, implemented as a doubly-linked list. Fast O(1) push/pop at both ends.

```bash
# Push elements
LPUSH tasks "task3"     # push to left (head)
LPUSH tasks "task2"
LPUSH tasks "task1"
RPUSH tasks "task4"     # push to right (tail)

# View elements
LRANGE tasks 0 -1       # all elements
# 1) "task1"
# 2) "task2"
# 3) "task3"
# 4) "task4"

LRANGE tasks 0 1        # first two
LLEN tasks              # length = 4

# Pop elements
LPOP tasks              # "task1" (removes from head)
RPOP tasks              # "task4" (removes from tail)
LPOP tasks 2            # pop 2 from left

# Blocking pop (for queue consumers)
BRPOP queue:jobs 0      # block forever until item available
BLPOP queue:jobs 30     # block up to 30 seconds

# Insert and trim
LINSERT tasks BEFORE "task3" "task2.5"
LTRIM tasks 0 99        # keep only first 100 elements

# Move between lists (atomic)
LMOVE source dest LEFT RIGHT

# Index access
LINDEX tasks 0          # get element at index 0
LSET tasks 0 "new-task" # set element at index 0
```

**Use cases:** message queues, activity feeds, recent items (LPUSH + LTRIM), task lists.

---

## Sets

Unordered collections of unique strings. O(1) add, remove, and lookup.

```bash
# Add members
SADD tags "python" "redis" "backend"

# Check membership
SISMEMBER tags "redis"    # 1 (true)
SISMEMBER tags "java"     # 0 (false)

# Count and list
SCARD tags                # 3
SMEMBERS tags             # {"python", "redis", "backend"}

# Remove
SREM tags "backend"

# Random sample (useful for random selection)
SRANDMEMBER tags          # random member
SRANDMEMBER tags 2        # 2 random members
SPOP tags                 # random member + remove

# Set operations
SADD set1 "a" "b" "c"
SADD set2 "b" "c" "d"

SINTER set1 set2          # {"b", "c"} — intersection
SUNION set1 set2          # {"a","b","c","d"} — union
SDIFF  set1 set2          # {"a"} — in set1 but not set2

# Store results in new key
SINTERSTORE result set1 set2
SUNIONSTORE result set1 set2
```

**Use cases:** unique visitors, tags, friend lists, "users who liked this", deduplication.

---

## Hashes

Maps of field → value pairs. Like a Python dict or a database row. More memory-efficient than storing each field as a separate key.

```bash
# Set fields
HSET user:1000 name "Alice" age 30 city "NYC" email "alice@example.com"

# Get field(s)
HGET  user:1000 name          # "Alice"
HMGET user:1000 name city     # "Alice", "NYC"
HGETALL user:1000             # all fields and values

# Update single field
HSET user:1000 age 31

# Delete field
HDEL user:1000 email

# Check field exists
HEXISTS user:1000 name        # 1
HEXISTS user:1000 phone       # 0

# Get just keys or values
HKEYS user:1000               # name, age, city
HVALS user:1000               # "Alice", 31, "NYC"
HLEN  user:1000               # 3

# Atomic increment
HINCRBY  user:1000 age 1      # 32
HINCRBYFLOAT user:1000 balance 10.50

# Set only if not exists
HSETNX user:1000 phone "555-1234"
```

**Use cases:** user profiles, product objects, configuration, game state, session data.

---

## Sorted Sets (ZSets)

Members with associated float scores, kept in order. Every operation maintains sorted order. One of Redis's most powerful structures.

```bash
# Add members with scores
ZADD leaderboard 9850 "alice"
ZADD leaderboard 7200 "bob"
ZADD leaderboard 8100 "charlie"
ZADD leaderboard 9850 "diana"   # same score as alice

# Range by rank (0 = lowest score)
ZRANGE leaderboard 0 -1              # all, ascending
ZREVRANGE leaderboard 0 2            # top 3, descending
ZRANGE leaderboard 0 -1 WITHSCORES  # with scores
ZRANGE leaderboard 0 -1 REV         # descending (Redis 6.2+)

# Range by score
ZRANGEBYSCORE leaderboard 8000 9999          # scores 8000–9999
ZRANGEBYSCORE leaderboard -inf +inf          # all
ZRANGEBYSCORE leaderboard "(8000" 9999       # exclusive lower bound

# Rank of a member (0-based)
ZRANK    leaderboard "alice"         # 3 (0 = lowest)
ZREVRANK leaderboard "alice"         # 0 (0 = highest)

# Score of a member
ZSCORE leaderboard "alice"           # 9850

# Count members in score range
ZCOUNT leaderboard 8000 9999         # 2

# Remove members
ZREM leaderboard "bob"
ZPOPMIN leaderboard                  # remove lowest
ZPOPMAX leaderboard 3                # remove top 3

# Increment score
ZINCRBY leaderboard 100 "charlie"    # 8200

# Set operations
ZUNIONSTORE dest 2 set1 set2 WEIGHTS 1 0.5
ZINTERSTORE dest 2 set1 set2
```

**Use cases:** leaderboards, priority queues, rate limiting with sliding window, time-sorted events, autocomplete.

---

## Choosing the Right Data Structure

```
Question                                    → Use
─────────────────────────────────────────────────────────
Store a single value or counter             → String
Cache a full object (JSON)                  → String (serialized)
Store an object with individual fields      → Hash
Ordered list, queue, or stack               → List
Unique items, tags, or set math             → Set
Ranked list, priority queue, time window    → Sorted Set
Event stream with consumer groups           → Stream (Chapter 4)
Count distinct items (approx)              → HyperLogLog (Chapter 4)
Track which items were seen (bit array)    → Bitmap (Chapter 4)
Geographic proximity search                → Geospatial (Chapter 4)
```

---

## Key Commands (All Types)

```bash
# Check if key exists
EXISTS user:1000             # 1 or 0
EXISTS k1 k2 k3              # count of existing keys

# Get type of a key
TYPE user:1000               # "hash"
TYPE counter                 # "string"
TYPE tasks                   # "list"

# Delete keys
DEL key1 key2 key3           # synchronous delete
UNLINK key1 key2             # async delete (non-blocking)

# Find keys (avoid in production — use SCAN)
KEYS "user:*"                # DANGEROUS on large keyspaces
SCAN 0 MATCH "user:*" COUNT 100  # safe iteration

# Rename
RENAME oldkey newkey
RENAMENX oldkey newkey       # only if newkey doesn't exist

# Object encoding (memory optimization info)
OBJECT ENCODING mykey        # ziplist, listpack, hashtable, etc.
OBJECT IDLETIME mykey        # seconds since last access
OBJECT FREQ mykey            # LFU frequency counter
```


---

# Chapter 4: Advanced Data Types

## Redis Streams

A persistent, append-only log of messages — like Apache Kafka but built into Redis. Each entry has an auto-generated time-based ID and a set of fields.

```bash
# Add entries to a stream
XADD events * user_id user-1 event_type page_view page /home
XADD events * user_id user-2 event_type purchase amount 99.99
XADD events 1704067200000-0 user_id user-1 event_type click  # custom ID

# Read entries
XLEN events                         # count
XRANGE events - +                   # all entries (- = min, + = max)
XRANGE events - + COUNT 10          # first 10
XREVRANGE events + - COUNT 5        # last 5

# Read from a specific ID onward
XRANGE events 1704067200000-0 +

# Trim to max length
XADD events MAXLEN ~ 10000 * field value  # ~ = approximate trim
XTRIM events MAXLEN 5000
```

### Consumer Groups

Consumer groups allow multiple consumers to process different entries:

```bash
# Create consumer group ($ = only new entries)
XGROUP CREATE events mygroup $ MKSTREAM
XGROUP CREATE events mygroup 0      # from beginning

# Read as a consumer
XREADGROUP GROUP mygroup worker-1 COUNT 10 STREAMS events >
# > means: give me undelivered entries

# Acknowledge processed entries
XACK events mygroup 1704067200000-0

# Check pending (unacknowledged) entries
XPENDING events mygroup - + 10

# Claim stale entries from a crashed consumer
XCLAIM events mygroup worker-2 30000 1704067200000-0
# or auto-claim
XAUTOCLAIM events mygroup worker-2 30000 0-0 COUNT 50
```

**Use cases:** event sourcing, task queues with consumer groups, activity logs.

---

## HyperLogLog

A probabilistic data structure for **approximate cardinality estimation** (counting unique items). Uses only ~12 KB of memory regardless of set size, with ~0.81% error.

```bash
# Add elements
PFADD unique-visitors "user-1" "user-2" "user-3"
PFADD unique-visitors "user-1"   # duplicates ignored

# Get approximate count
PFCOUNT unique-visitors           # ~3

# Per-day visitors
PFADD visitors:2024-01-01 "user-1" "user-2" "user-100"
PFADD visitors:2024-01-02 "user-2" "user-3" "user-101"

# Merge and count unique across multiple keys
PFMERGE visitors:week visitors:2024-01-01 visitors:2024-01-02
PFCOUNT visitors:week            # ~5 unique users across 2 days
```

**Use cases:** unique visitor counts, unique search queries, approximate distinct counts where exactness isn't critical.

---

## Bitmaps

Bit arrays stored efficiently in Redis strings. Set or get individual bits. 512 MB string = ~4 billion bits.

```bash
# Set bits (e.g., user activity: 1 = active on that day)
SETBIT activity:user-1:2024-01 0  1   # day 1 = active
SETBIT activity:user-1:2024-01 4  1   # day 5 = active
SETBIT activity:user-1:2024-01 29 1   # day 30 = active

# Get bit
GETBIT activity:user-1:2024-01 0      # 1
GETBIT activity:user-1:2024-01 1      # 0

# Count set bits
BITCOUNT activity:user-1:2024-01       # 3 active days
BITCOUNT key 0 0                       # bits in first byte

# Find first set or clear bit
BITPOS activity:user-1:2024-01 1       # first active day offset
BITPOS activity:user-1:2024-01 0       # first inactive day

# Bitwise operations on multiple bitmaps
BITOP AND result active:jan active:feb  # active both months
BITOP OR  result active:jan active:feb  # active either month
BITOP XOR result active:jan active:feb  # active one but not both

# How many users were active on day 5?
BITCOUNT active:2024-01-05             # count all users active that day
```

**Use cases:** daily active user tracking, feature flag rollouts, visited items, compact boolean arrays.

---

## Geospatial

Store geographic coordinates and query by proximity. Internally uses a Sorted Set with geohash scores.

```bash
# Add locations (longitude, latitude, name)
GEOADD restaurants -73.9857 40.7484 "Le Bernardin"
GEOADD restaurants -73.9712 40.7614 "Per Se"
GEOADD restaurants -73.9442 40.6782 "Peter Luger"

# Distance between two members
GEODIST restaurants "Le Bernardin" "Per Se" km  # distance in km

# Find nearby (within radius)
GEOSEARCH restaurants FROMMEMBER "Le Bernardin" BYRADIUS 2 km ASC COUNT 5
GEOSEARCH restaurants FROMLONLAT -73.9857 40.7484 BYRADIUS 5 km ASC

# Search with box
GEOSEARCH restaurants FROMLONLAT -73.9857 40.7484 BYBOX 10 10 km ASC

# Get coordinates
GEOPOS restaurants "Le Bernardin"    # (-73.9857, 40.7484)

# Get geohash string
GEOHASH restaurants "Le Bernardin"   # "dr5ruecge40"
```

**Use cases:** store locator, "find nearby" features, delivery radius checks, geographic filtering.

---

## Sorted Set Advanced Patterns

### Sliding Window Rate Limiter

```python
import time
import redis

r = redis.Redis()

def check_rate_limit(user_id: str, max_requests: int, window_seconds: int) -> bool:
    key   = f"rate:{user_id}"
    now   = time.time()
    window_start = now - window_seconds

    pipe = r.pipeline()
    pipe.zremrangebyscore(key, 0, window_start)   # remove old entries
    pipe.zadd(key, {str(now): now})               # add current request
    pipe.zcard(key)                                # count in window
    pipe.expire(key, window_seconds)
    _, _, count, _ = pipe.execute()

    return count <= max_requests
```

### Time-Series with Sorted Set

```python
# Store events by timestamp
score = int(time.time() * 1000)   # milliseconds
r.zadd("events:stream", {f"event:{uuid}": score})

# Get events in time range
start_ms = int((time.time() - 3600) * 1000)  # last hour
end_ms   = int(time.time() * 1000)
r.zrangebyscore("events:stream", start_ms, end_ms)
```

---

## Data Type Memory Usage

| Type | Internal Encoding | Threshold |
|---|---|---|
| String | `int` | Integer values |
| String | `embstr` | ≤ 44 bytes |
| String | `raw` | > 44 bytes |
| List | `listpack` | ≤ 128 elements, values ≤ 64 bytes |
| List | `quicklist` | > threshold |
| Hash | `listpack` | ≤ 128 fields, values ≤ 64 bytes |
| Hash | `hashtable` | > threshold |
| Set | `listpack` | ≤ 128 members, values ≤ 64 bytes |
| Set | `intset` | All integers, ≤ 512 members |
| Set | `hashtable` | > threshold |
| Sorted Set | `listpack` | ≤ 128 members, values ≤ 64 bytes |
| Sorted Set | `skiplist` | > threshold |

```bash
# Check encoding
OBJECT ENCODING mykey
```

Redis automatically upgrades encoding when thresholds are exceeded. Keep values small to stay in memory-efficient encodings.


---

# Chapter 5: Expiry & Eviction

## Key Expiry (TTL)

Every Redis key can have an optional Time-To-Live. When the TTL expires, Redis automatically deletes the key.

```bash
# Set with expiry at creation
SET session:abc "data" EX 3600     # seconds
SET session:xyz "data" PX 3600000  # milliseconds

# Set expiry on existing key
EXPIRE key 300               # expire in 300 seconds
PEXPIRE key 300000           # expire in 300,000 milliseconds
EXPIREAT key 1735689600      # expire at Unix timestamp
PEXPIREAT key 1735689600000  # expire at Unix timestamp (ms)

# Check TTL
TTL key      # seconds remaining (-1 = no expiry, -2 = doesn't exist)
PTTL key     # milliseconds remaining

# Remove expiry (make persistent)
PERSIST key

# EXPIRETIME (Redis 7.0+)
EXPIRETIME key    # Unix timestamp when key expires

# Conditional expiry (Redis 7.0+)
EXPIRE key 300 NX   # set expiry only if key has no expiry
EXPIRE key 300 XX   # set expiry only if key already has expiry
EXPIRE key 300 GT   # set expiry only if new TTL > current TTL
EXPIRE key 300 LT   # set expiry only if new TTL < current TTL
```

---

## How Redis Expires Keys

Redis uses two mechanisms:

### 1. Lazy Expiration
When a key is accessed, Redis checks if it's expired and deletes it on access.

### 2. Active Expiration
Every 100ms, Redis randomly samples keys with TTLs and deletes expired ones. Redis 7+ uses a more efficient approach.

```
Active expiry loop:
  1. Sample 20 random keys from the set of keys with expiry
  2. Delete all expired keys found
  3. If > 25% were expired, repeat immediately
  4. Otherwise, wait 100ms
```

**Implication:** A key may live slightly past its TTL if not accessed and not sampled. For most use cases this is fine.

---

## Eviction Policies

When Redis reaches `maxmemory`, it uses an eviction policy to free space.

```properties
# redis.conf
maxmemory 512mb
maxmemory-policy allkeys-lru
```

### Policy Reference

| Policy | Description | Use Case |
|---|---|---|
| `noeviction` | Return error when full (default) | Databases, never lose data |
| `allkeys-lru` | Evict least recently used key from all keys | General cache |
| `allkeys-lfu` | Evict least frequently used from all keys | Cache with skewed access |
| `allkeys-random` | Evict random key from all keys | Uniform access patterns |
| `volatile-lru` | Evict LRU key among keys with TTL set | Mixed persistent + cache |
| `volatile-lfu` | Evict LFU key among keys with TTL set | Mixed with frequency bias |
| `volatile-random` | Random eviction among keys with TTL | Mixed workloads |
| `volatile-ttl` | Evict key with shortest TTL first | Time-sensitive cache |

```bash
# Set policy at runtime
CONFIG SET maxmemory-policy allkeys-lru
CONFIG SET maxmemory 1gb

# Check current memory usage
INFO memory
# used_memory: 104857600
# used_memory_human: 100.00M
# maxmemory: 536870912
# maxmemory_human: 512.00M
# mem_fragmentation_ratio: 1.05
```

---

## LRU vs LFU

### LRU (Least Recently Used)
Evicts the key that has not been accessed for the longest time.

```
Key A: last accessed 10 min ago  ← evicted first
Key B: last accessed 5 min ago
Key C: last accessed 1 min ago   ← kept
```

**Best for:** Working set cache — recently used data is likely to be used again.

### LFU (Least Frequently Used)
Evicts the key that has been accessed the fewest times (with time decay).

```
Key A: accessed 2 times total    ← evicted first
Key B: accessed 50 times total
Key C: accessed 100 times total  ← kept
```

**Best for:** Long-term caches where popular items should survive even if not recently accessed.

```properties
# LFU tuning
lfu-log-factor 10    # higher = less saturation of counter
lfu-decay-time 1     # minutes per decay (counter halves)
```

---

## Keyspace Notifications

Redis can publish notifications when keys expire or are modified:

```properties
# redis.conf — enable notifications
# K = keyspace events, E = keyevent events
# x = expired events, g = generic commands, s = set commands
# l = list commands, z = sorted set commands
notify-keyspace-events "KExa"
```

```bash
# Subscribe to expiry events on database 0
redis-cli SUBSCRIBE "__keyevent@0__:expired"

# Subscribe to all set operations
redis-cli PSUBSCRIBE "__keyevent@0__:set"

# Subscribe to keyspace events for a specific key
redis-cli SUBSCRIBE "__keyspace@0__:mykey"
```

```python
# Python: listen for expired keys
import redis

r = redis.Redis()
p = r.pubsub()
p.psubscribe("__keyevent@0__:expired")

for msg in p.listen():
    if msg["type"] == "pmessage":
        expired_key = msg["data"].decode()
        print(f"Key expired: {expired_key}")
        # Take action: refresh cache, notify service, etc.
```

**Use cases:** cache refresh triggers, session timeout notifications, scheduled job cleanup.

---

## TTL Best Practices

```python
# 1. Always set a TTL on cache keys
redis.set("cache:product:123", json.dumps(product), ex=3600)

# 2. Add jitter to prevent cache stampede
import random
ttl = 3600 + random.randint(-300, 300)   # ±5 min jitter
redis.set(key, value, ex=ttl)

# 3. Use PERSIST for permanent data, TTL for temporary
# Permanent: user profiles, product catalog
# Temporary: sessions, cache, rate limits, OTPs

# 4. Don't rely on exact expiry timing
# A key may live slightly past TTL (lazy expiry)
# Never use Redis TTL for precise scheduling

# 5. Refresh TTL on access (sliding expiry)
value = redis.get(key)
if value:
    redis.expire(key, 3600)   # reset TTL on hit
```

---

## Memory Optimization

```bash
# Check memory usage of a specific key
MEMORY USAGE key            # bytes

# Check memory breakdown
INFO memory

# Analyze large keys (Redis 4.0+)
redis-cli --memkeys --memkeys-samples 100

# Run memory doctor
MEMORY DOCTOR

# Defragment memory (active defrag)
CONFIG SET activedefrag yes
CONFIG SET active-defrag-ignore-bytes 100mb
CONFIG SET active-defrag-threshold-lower 10
```


---

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


---

# Chapter 7: Pub/Sub & Messaging

## Pub/Sub

Redis Pub/Sub implements a simple publish-subscribe messaging model. Publishers send messages to channels; subscribers receive all messages on channels they subscribe to.

```
Publisher ──PUBLISH──► Channel "news:tech" ──► Subscriber A
                                           ──► Subscriber B
                                           ──► Subscriber C
```

**Key properties:**
- Fire-and-forget: messages not stored, not replayable
- No message acknowledgment
- Subscribers must be connected to receive messages
- Scales well for real-time notifications, chat, broadcasts

---

## Basic Pub/Sub Commands

### Subscriber

```bash
# Subscribe to a channel
SUBSCRIBE news:tech

# Subscribe to multiple channels
SUBSCRIBE news:tech news:sports alerts

# Pattern subscribe (glob patterns)
PSUBSCRIBE news:*       # all news channels
PSUBSCRIBE user:*:notifications
PSUBSCRIBE "h?llo"      # hello, hallo, hxllo

# Unsubscribe
UNSUBSCRIBE news:tech
PUNSUBSCRIBE news:*
```

### Publisher

```bash
# Publish a message
PUBLISH news:tech "Redis 8.0 released!"
# Returns: number of subscribers that received the message

PUBLISH user:1000:notifications '{"type":"mention","from":"alice"}'
```

### Inspect Channels

```bash
# List active channels (with subscribers)
PUBSUB CHANNELS              # all active channels
PUBSUB CHANNELS "news:*"     # matching pattern

# Count subscribers per channel
PUBSUB NUMSUB news:tech news:sports

# Count pattern subscribers
PUBSUB NUMPAT
```

---

## Pub/Sub in Python

```python
import redis
import threading
import json

r = redis.Redis(host="localhost", port=6379)


# ─── Subscriber ───
def message_handler(message):
    if message["type"] == "message":
        channel = message["channel"].decode()
        data    = json.loads(message["data"])
        print(f"[{channel}] {data}")

pubsub = r.pubsub()
pubsub.subscribe(**{
    "notifications:user-1": message_handler,
    "broadcast:all":        message_handler,
})

# Run listener in background thread
thread = pubsub.run_in_thread(sleep_time=0.01, daemon=True)


# ─── Publisher ───
def notify_user(user_id: str, event: dict) -> None:
    r.publish(f"notifications:{user_id}", json.dumps(event))

notify_user("user-1", {"type": "message", "from": "alice", "text": "Hello!"})

# Pattern subscription
pubsub2 = r.pubsub()
pubsub2.psubscribe(**{"notifications:*": message_handler})
```

---

## Pub/Sub Limitations

| Limitation | Workaround |
|---|---|
| No message persistence | Use Redis Streams (Chapter 4 / below) |
| No acknowledgment | Use Redis Streams with consumer groups |
| No replay | Use Redis Streams |
| Subscriber must be online | Use Streams or List queues |
| No filtering | Pattern subscribe or app-level filtering |

---

## Redis Streams as a Message Queue

Redis Streams (XADD/XREADGROUP) are the production-grade alternative to Pub/Sub when you need persistence, consumer groups, and acknowledgment.

### Producer Pattern

```python
import redis
import json
import time

r = redis.Redis()

def publish_event(stream: str, event: dict) -> str:
    return r.xadd(stream, event, maxlen=100000, approximate=True)

# Produce events
msg_id = publish_event("orders", {
    "order_id": "ORD-123",
    "user_id":  "user-42",
    "amount":   "99.99",
})
print(f"Published: {msg_id}")
```

### Consumer Group Pattern

```python
import redis
import json

r = redis.Redis()

STREAM  = "orders"
GROUP   = "order-processor"
CONSUMER = "worker-1"

# Create group (only once)
try:
    r.xgroup_create(STREAM, GROUP, id="0", mkstream=True)
except redis.exceptions.ResponseError:
    pass  # group already exists

def process_orders():
    while True:
        # Read up to 10 undelivered messages
        results = r.xreadgroup(
            groupname  = GROUP,
            consumername = CONSUMER,
            streams    = {STREAM: ">"},
            count      = 10,
            block      = 2000,   # block 2 seconds if no messages
        )

        if not results:
            continue

        for stream_name, messages in results:
            for msg_id, fields in messages:
                try:
                    order_id = fields[b"order_id"].decode()
                    amount   = float(fields[b"amount"])
                    print(f"Processing order {order_id}: ${amount}")
                    process_order(order_id, amount)

                    # Acknowledge: remove from pending
                    r.xack(STREAM, GROUP, msg_id)

                except Exception as e:
                    print(f"Failed: {e} — will retry")
                    # Don't ack — will be redelivered or claimed

process_orders()
```

### Handle Pending (Unacked) Messages

```python
# Check pending messages
pending = r.xpending(STREAM, GROUP, "-", "+", 10)
for p in pending:
    msg_id, consumer, idle_ms, delivery_count = p
    if idle_ms > 30000:  # stale for 30+ seconds
        # Claim and reprocess
        claimed = r.xclaim(STREAM, GROUP, CONSUMER, 30000, [msg_id])
        for mid, fields in claimed:
            reprocess(mid, fields)
```

---

## List as a Simple Queue

For simple task queues without consumer groups:

```python
QUEUE = "tasks:email"

# Producer
def enqueue(task: dict) -> None:
    r.lpush(QUEUE, json.dumps(task))

# Consumer (blocking)
def worker():
    while True:
        result = r.brpop(QUEUE, timeout=5)
        if result:
            _, data = result
            task = json.loads(data)
            process_task(task)

# Priority queue with multiple lists
HIGH_QUEUE = "tasks:high"
LOW_QUEUE  = "tasks:low"

# BRPOP waits on multiple lists in order (high priority first)
result = r.brpop([HIGH_QUEUE, LOW_QUEUE], timeout=10)
```

---

## Choosing Between Pub/Sub, Streams, and Lists

| Feature | Pub/Sub | Streams | List Queue |
|---|---|---|---|
| Persistence | No | Yes | Yes |
| Replay | No | Yes | No |
| Multiple consumers | All get each msg | Groups = divided | Round-robin |
| Acknowledgment | No | Yes (XACK) | Manual |
| Consumer groups | No | Yes | No |
| Message ordering | FIFO (channel) | FIFO (stream) | FIFO |
| Use case | Live notifications | Event sourcing, durable queue | Simple task queue |


---

# Chapter 8: Transactions & Scripting

## MULTI / EXEC Transactions

Redis transactions queue commands and execute them atomically as a single batch. Either all commands execute or none do (on EXEC failure).

```bash
# Basic transaction
MULTI                   # start transaction
SET balance 100
INCR counter
ZADD leaderboard 9850 "alice"
EXEC                    # execute all atomically

# Discard transaction
MULTI
SET key1 "value1"
DISCARD                 # cancel, queue is cleared
```

**Important nuances:**
- Redis transactions are NOT rollback-capable — if a command fails after EXEC, other commands still run
- Command errors detected at EXEC time (e.g., wrong type) don't abort the transaction
- Commands are queued, not executed immediately

```bash
MULTI
SET counter 0
INCR counter         # queued as QUEUED
LPUSH counter "x"    # will error (wrong type) but others still run
EXEC
# 1) OK
# 2) (integer) 1
# 3) (error) WRONGTYPE ...  ← only this failed
```

---

## WATCH (Optimistic Locking)

WATCH implements optimistic concurrency control. If any watched key changes before EXEC, the transaction is aborted.

```bash
# Pattern: check-and-set
WATCH balance
current = GET balance          # read current value
MULTI
SET balance (current + 10)     # set new value
EXEC                           # returns nil if balance changed
                               # returns [OK] if no conflict
```

```python
import redis

r = redis.Redis()

def transfer(from_key: str, to_key: str, amount: float) -> bool:
    with r.pipeline() as pipe:
        while True:
            try:
                pipe.watch(from_key, to_key)     # watch both keys

                from_bal = float(pipe.get(from_key) or 0)
                to_bal   = float(pipe.get(to_key) or 0)

                if from_bal < amount:
                    pipe.unwatch()
                    return False  # insufficient funds

                pipe.multi()                      # start transaction
                pipe.set(from_key, from_bal - amount)
                pipe.set(to_key,   to_bal   + amount)
                pipe.execute()                    # atomic execution
                return True

            except redis.WatchError:
                # Another client modified the keys — retry
                continue

transfer("account:alice", "account:bob", 50.0)
```

**When WATCH returns nil (transaction aborted):**
- The watched key was modified by another client between WATCH and EXEC
- Retry the entire operation from WATCH

---

## Lua Scripting

Lua scripts run on the Redis server atomically. They can read and write any key within the script.

```bash
# Run inline script
EVAL "return 'hello from lua'" 0

# Script with keys and arguments
EVAL "
  local current = redis.call('GET', KEYS[1])
  if not current then
    redis.call('SET', KEYS[1], ARGV[1])
    return 1
  end
  return 0
" 1 my_key default_value

# KEYS[1], KEYS[2]... = key arguments (cluster-routing aware)
# ARGV[1], ARGV[2]... = value arguments
```

### Redis-Callable Functions

```lua
-- redis.call: raises error on failure
local val = redis.call("GET", KEYS[1])

-- redis.pcall: returns error as a Lua table (no exception)
local ok, err = pcall(function()
    return redis.call("INCR", KEYS[1])
end)

-- redis.status_reply: return status
return redis.status_reply("OK")

-- redis.error_reply: return error
return redis.error_reply("ERR something went wrong")
```

### Practical Lua Scripts

#### Atomic Increment with Limit (Rate Limiter)

```lua
-- rate_limit.lua
-- KEYS[1]: rate limit key
-- ARGV[1]: max requests
-- ARGV[2]: window seconds
local current = tonumber(redis.call("INCR", KEYS[1]))
if current == 1 then
    redis.call("EXPIRE", KEYS[1], ARGV[2])
end
if current > tonumber(ARGV[1]) then
    return 0  -- rate limited
end
return 1  -- allowed
```

```python
rate_script = r.register_script(open("rate_limit.lua").read())
allowed = rate_script(keys=["rate:user-1"], args=[100, 60])
```

#### Get and Delete (Atomic)

```lua
-- get_and_delete.lua
local val = redis.call("GET", KEYS[1])
if val then
    redis.call("DEL", KEYS[1])
end
return val
```

#### Conditional Set

```lua
-- set_if_less.lua: set value only if current value is less
local current = tonumber(redis.call("GET", KEYS[1])) or 0
local new_val = tonumber(ARGV[1])
if new_val < current or current == 0 then
    redis.call("SET", KEYS[1], ARGV[1])
    return 1
end
return 0
```

---

## Script Caching (EVALSHA)

Large scripts can be cached on the server to avoid resending on every call:

```bash
# Load script and get SHA1
SCRIPT LOAD "return redis.call('GET', KEYS[1])"
# "2067d915024a3e1657c4169c84f809f8ec75b9a7"

# Execute by SHA1 (faster — no script transmission)
EVALSHA 2067d915024a3e1657c4169c84f809f8ec75b9a7 1 mykey

# Check if script is cached
SCRIPT EXISTS sha1_hash

# Flush all cached scripts
SCRIPT FLUSH
```

```python
import redis

r = redis.Redis()

# register_script caches and uses EVALSHA automatically
get_and_delete = r.register_script("""
local val = redis.call("GET", KEYS[1])
if val then redis.call("DEL", KEYS[1]) end
return val
""")

value = get_and_delete(keys=["otp:user-1"])
```

---

## Redis Functions (Redis 7.0+)

Functions are named, persistent Lua scripts loaded into the server:

```bash
# Load a function library
FUNCTION LOAD "#!lua name=mylib\n
redis.register_function('getset', function(keys, args)
    local old = redis.call('GET', keys[1])
    redis.call('SET', keys[1], args[1])
    return old
end)"

# Call the function
FCALL getset 1 mykey new-value

# List loaded functions
FUNCTION LIST

# Delete a library
FUNCTION DELETE mylib
```

---

## Transactions vs Lua: When to Use Each

| Feature | MULTI/EXEC | Lua Script |
|---|---|---|
| Atomicity | Yes (all-or-nothing) | Yes (single thread) |
| Read-then-write | Only with WATCH | Yes (native) |
| Conditional logic | No | Yes |
| Error handling | Limited | Full (pcall) |
| Performance | Batched round-trips | Single round-trip |
| Complexity | Low | Higher |
| Use case | Simple atomic batches | Complex atomic logic |

**Rule of thumb:**
- Simple batched operations → MULTI/EXEC
- Conditional logic, read-modify-write → Lua
- Reusable server-side logic → Redis Functions


---

# Chapter 9: Caching Patterns

## Why Cache?

```
Without cache:
  Every request → Database → Disk I/O → 5–50ms
  100 req/s = 100 DB queries/s → DB overload

With cache:
  Cache hit  → Redis in-memory → 0.1ms
  Cache miss → DB query → Store in Redis → 5–50ms
  99% cache hit rate → 100× fewer DB queries
```

---

## Cache-Aside (Lazy Loading)

The application controls the cache. On a cache miss, the app fetches from DB and populates the cache.

```python
import redis
import json
from typing import Optional

r = redis.Redis()

def get_product(product_id: str) -> Optional[dict]:
    cache_key = f"cache:product:{product_id}"

    # 1. Try cache first
    cached = r.get(cache_key)
    if cached:
        return json.loads(cached)

    # 2. Cache miss — fetch from database
    product = db.query("SELECT * FROM products WHERE id = %s", product_id)
    if not product:
        return None

    # 3. Store in cache with TTL
    r.setex(cache_key, 3600, json.dumps(product))

    return product
```

**Pros:**
- Only cache what's actually requested
- Cache failure doesn't break the app
- Data in cache is always fresh (at most TTL-old)

**Cons:**
- First request always misses (cold cache)
- Possible stale data for up to TTL

---

## Write-Through

Update cache synchronously on every write:

```python
def update_product(product_id: str, data: dict) -> None:
    # 1. Update database
    db.execute("UPDATE products SET ... WHERE id = %s", product_id)

    # 2. Update cache immediately
    cache_key = f"cache:product:{product_id}"
    r.setex(cache_key, 3600, json.dumps(data))
```

**Pros:** Cache always consistent with DB. No stale reads.
**Cons:** Every write hits both DB and cache. Wasted space for items never read.

---

## Write-Behind (Write-Back)

Write to cache first, persist to DB asynchronously:

```python
import threading

def update_product_fast(product_id: str, data: dict) -> None:
    cache_key  = f"cache:product:{product_id}"
    dirty_key  = f"dirty:product:{product_id}"

    # 1. Write to cache immediately (fast response)
    pipe = r.pipeline()
    pipe.set(cache_key, json.dumps(data))
    pipe.sadd("dirty:products", product_id)   # track dirty keys
    pipe.execute()

# Background worker flushes dirty keys to DB
def flush_dirty_products():
    while True:
        product_ids = r.smembers("dirty:products")
        for pid in product_ids:
            data = r.get(f"cache:product:{pid.decode()}")
            if data:
                db.execute("UPDATE products ...", json.loads(data))
                r.srem("dirty:products", pid)
        time.sleep(1)
```

**Pros:** Extremely fast writes. DB write batching.
**Cons:** Risk of data loss if Redis crashes before flush. Complex.

---

## Cache Invalidation

```python
# 1. Delete on write (simplest)
def update_product(product_id: str, data: dict) -> None:
    db.execute("UPDATE products ...", data)
    r.delete(f"cache:product:{product_id}")
    # Next read will miss and repopulate

# 2. Update on write
def update_product(product_id: str, data: dict) -> None:
    db.execute("UPDATE products ...", data)
    r.setex(f"cache:product:{product_id}", 3600, json.dumps(data))

# 3. Tag-based invalidation
def invalidate_category(category: str) -> None:
    # Track all product keys for a category
    product_ids = r.smembers(f"cache:tag:category:{category}")
    if product_ids:
        keys = [f"cache:product:{pid.decode()}" for pid in product_ids]
        r.delete(*keys)
        r.delete(f"cache:tag:category:{category}")
```

---

## Cache Stampede (Thundering Herd)

When a popular cache key expires, many requests simultaneously miss the cache and flood the DB.

```
t=0: Key expires
t=1: 1000 requests arrive → all miss cache
t=1: 1000 DB queries fire simultaneously → DB overload
```

### Solution 1: Mutex Lock

```python
import redis
import json

r = redis.Redis()

def get_with_lock(key: str, ttl: int, fetch_fn) -> dict:
    # Try cache
    cached = r.get(key)
    if cached:
        return json.loads(cached)

    lock_key = f"lock:{key}"

    # Try to acquire lock
    acquired = r.set(lock_key, "1", nx=True, ex=10)
    if acquired:
        try:
            data = fetch_fn()
            r.setex(key, ttl, json.dumps(data))
            return data
        finally:
            r.delete(lock_key)
    else:
        # Wait and retry (another worker is fetching)
        import time
        time.sleep(0.1)
        return get_with_lock(key, ttl, fetch_fn)
```

### Solution 2: Probabilistic Early Refresh

```python
import math, random, time

def get_with_early_refresh(key: str, ttl: int, fetch_fn, beta: float = 1.0) -> dict:
    data  = r.get(key)
    expiry = r.expiretime(key)

    if data:
        # Probabilistic early expiry check
        gap = expiry - time.time()
        if gap > 0 and -beta * math.log(random.random()) >= gap:
            # Refresh early in background
            threading.Thread(target=lambda: refresh(key, ttl, fetch_fn)).start()
        return json.loads(data)

    return refresh(key, ttl, fetch_fn)
```

### Solution 3: TTL Jitter

```python
import random

def cache_set(key: str, value: dict, base_ttl: int) -> None:
    jitter = random.randint(-base_ttl // 10, base_ttl // 10)
    r.setex(key, base_ttl + jitter, json.dumps(value))
```

---

## Session Storage Pattern

```python
import uuid
import json
import redis

r = redis.Redis()
SESSION_TTL = 3600  # 1 hour

def create_session(user_data: dict) -> str:
    token = str(uuid.uuid4())
    r.setex(f"session:{token}", SESSION_TTL, json.dumps(user_data))
    return token

def get_session(token: str) -> dict | None:
    data = r.get(f"session:{token}")
    if data:
        r.expire(f"session:{token}", SESSION_TTL)  # sliding expiry
        return json.loads(data)
    return None

def invalidate_session(token: str) -> None:
    r.delete(f"session:{token}")

def invalidate_all_user_sessions(user_id: str) -> None:
    # Track all sessions per user
    session_keys = r.smembers(f"user:{user_id}:sessions")
    for key in session_keys:
        r.delete(f"session:{key.decode()}")
    r.delete(f"user:{user_id}:sessions")
```

---

## Rate Limiting Patterns

### Fixed Window

```python
def is_allowed_fixed(user_id: str, max_req: int, window: int) -> bool:
    key   = f"rate:{user_id}:{int(time.time()) // window}"
    count = r.incr(key)
    if count == 1:
        r.expire(key, window)
    return count <= max_req
```

### Sliding Window with Sorted Set

```python
def is_allowed_sliding(user_id: str, max_req: int, window: int) -> bool:
    key = f"rate:sliding:{user_id}"
    now = time.time()

    pipe = r.pipeline()
    pipe.zremrangebyscore(key, 0, now - window)
    pipe.zadd(key, {str(now): now})
    pipe.zcard(key)
    pipe.expire(key, window)
    _, _, count, _ = pipe.execute()

    return count <= max_req
```

---

## Distributed Lock

```python
import uuid, time

class RedisLock:
    def __init__(self, r: redis.Redis, key: str, ttl: int = 30):
        self.r   = r
        self.key = f"lock:{key}"
        self.ttl = ttl
        self.token = str(uuid.uuid4())

    def acquire(self, timeout: float = 10.0) -> bool:
        deadline = time.time() + timeout
        while time.time() < deadline:
            if self.r.set(self.key, self.token, nx=True, ex=self.ttl):
                return True
            time.sleep(0.05)
        return False

    def release(self) -> None:
        # Lua: only delete if token matches (atomic)
        script = """
        if redis.call('GET', KEYS[1]) == ARGV[1] then
            return redis.call('DEL', KEYS[1])
        else
            return 0
        end
        """
        self.r.eval(script, 1, self.key, self.token)

    def __enter__(self):
        if not self.acquire():
            raise TimeoutError(f"Could not acquire lock: {self.key}")
        return self

    def __exit__(self, *args):
        self.release()

# Usage
with RedisLock(r, "invoice:generator"):
    generate_invoice()
```


---

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


---

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


---

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


---

# Chapter 13: Security

## Security Layers

```
1. Network isolation      ← Redis only accessible from trusted hosts
2. Authentication (AUTH)  ← Require password or username+password
3. ACLs                   ← Per-user command and key restrictions
4. TLS                    ← Encrypt traffic in transit
5. Command restrictions   ← Disable dangerous commands
6. Protected mode         ← Default: reject external connections
```

---

## Authentication

### Simple Password (Redis < 6.0 style)

```properties
# redis.conf
requirepass yourStrongPassword123!
```

```bash
# Authenticate
AUTH yourStrongPassword123!
```

### ACL Users (Redis 6.0+, recommended)

```properties
# redis.conf — define users
# Disable default user (no password)
user default off nopass nocommands nokeys

# Create admin user
user admin on >adminPassword ~* &* +@all

# Create app user (read/write to app:* keys only)
user app-user on >appPassword ~app:* &* +@read +@write +@string +@hash +EXPIRE +DEL

# Create read-only user
user reader on >readPassword ~* &* +@read
```

---

## ACL System (Redis 6.0+)

### ACL Rule Components

```
user <username> [on|off]           ← enabled or disabled
                [>password]        ← set password (> prefix)
                [<password]        ← remove password (< prefix)
                [nopass]           ← no password required
                [~pattern]         ← allowed key patterns
                [%R~pattern]       ← read-only key pattern
                [%W~pattern]       ← write-only key pattern
                [&pattern]         ← allowed pub/sub channel patterns
                [+command]         ← allow command
                [-command]         ← deny command
                [+@category]       ← allow all commands in category
                [-@category]       ← deny all commands in category
                [allkeys]          ← all keys allowed (~*)
                [allchannels]      ← all channels allowed (&*)
                [allcommands]      ← all commands allowed (+@all)
                [nocommands]       ← no commands allowed (-@all)
                [resetkeys]        ← remove all key patterns
                [reset]            ← full reset to defaults
```

### ACL Commands

```bash
# List all users
ACL LIST

# Show current user
ACL WHOAMI

# Get specific user info
ACL GETUSER app-user

# Create/modify user
ACL SETUSER cache-user on >cachePass ~cache:* +GET +SET +DEL +EXPIRE +TTL

# Delete user
ACL DELUSER old-user

# Test permissions
ACL DRYRUN cache-user GET cache:product:1
ACL DRYRUN cache-user SET sensitive:data "value"   # returns error

# Log recent ACL failures
ACL LOG
ACL LOG RESET
ACL LOG COUNT
```

### ACL Categories

```bash
# List categories
ACL CAT

# List commands in a category
ACL CAT read
ACL CAT write
ACL CAT string
ACL CAT hash
ACL CAT pubsub
ACL CAT admin
ACL CAT dangerous
```

---

## TLS Configuration

### Generate Certificates

```bash
# Create CA and server certificates
mkdir -p /etc/redis/tls

# CA
openssl genrsa -out /etc/redis/tls/ca.key 4096
openssl req -x509 -new -nodes -key /etc/redis/tls/ca.key \
  -sha256 -days 1024 -out /etc/redis/tls/ca.crt \
  -subj "/CN=Redis-CA/O=MyOrg"

# Server cert
openssl genrsa -out /etc/redis/tls/redis.key 2048
openssl req -new -key /etc/redis/tls/redis.key \
  -out /etc/redis/tls/redis.csr \
  -subj "/CN=redis-server/O=MyOrg"
openssl x509 -req -in /etc/redis/tls/redis.csr \
  -CA /etc/redis/tls/ca.crt -CAkey /etc/redis/tls/ca.key \
  -CAcreateserial -out /etc/redis/tls/redis.crt -days 365 -sha256
```

### Redis TLS Config

```properties
# redis.conf
tls-port 6380
port 0                              # disable non-TLS port

tls-cert-file /etc/redis/tls/redis.crt
tls-key-file  /etc/redis/tls/redis.key
tls-ca-cert-file /etc/redis/tls/ca.crt
tls-auth-clients yes               # require client certificates
tls-protocols "TLSv1.3"
tls-ciphers DEFAULT:!MEDIUM
tls-prefer-server-ciphers yes

# Replication over TLS
tls-replication yes
```

### Client TLS (Python)

```python
import redis
import ssl

r = redis.Redis(
    host     = "redis-host",
    port     = 6380,
    ssl      = True,
    ssl_certfile  = "/etc/redis/tls/client.crt",
    ssl_keyfile   = "/etc/redis/tls/client.key",
    ssl_ca_certs  = "/etc/redis/tls/ca.crt",
    ssl_check_hostname = True,
)
r.ping()
```

---

## Disabling Dangerous Commands

```properties
# redis.conf — rename or disable dangerous commands
rename-command FLUSHALL ""          # disable completely
rename-command FLUSHDB  ""
rename-command CONFIG   "CONFIG-ADMIN-ONLY-abc123"
rename-command DEBUG    ""
rename-command KEYS     "SCAN-INSTEAD"
rename-command SHUTDOWN "SHUTDOWN-abc123"
```

---

## Network Security

```properties
# Bind to specific interfaces only
bind 127.0.0.1 10.0.0.1    # only localhost and internal network

# Protected mode (default: on)
# Rejects connections from external IPs if no auth is configured
protected-mode yes

# Max clients
maxclients 1000

# Connection timeout
timeout 300                  # close idle connections after 300s
tcp-keepalive 60
```

---

## Security Best Practices

```
1. Never expose Redis directly to the internet
   → Use firewall rules, VPC, security groups

2. Always use requirepass or ACL users
   → Rotate passwords periodically

3. Use TLS for all client connections in production
   → Especially across networks

4. Apply ACLs — least privilege per service
   → App user: only reads/writes to app:* keys
   → Monitoring user: INFO only, no writes

5. Rename or disable dangerous commands
   → FLUSHALL, FLUSHDB, DEBUG, CONFIG

6. Run Redis as a dedicated OS user (not root)
   → useradd -r -s /sbin/nologin redis

7. Use Redis Sentinel or Cluster behind a load balancer
   → Never expose individual node IPs to application code

8. Enable AOF persistence
   → Prevents data loss on Redis restart after attack

9. Monitor the ACL log for unauthorized access attempts
   → ACL LOG every 5 minutes

10. Keep Redis updated
    → Security patches released regularly
```


---

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


---

# Chapter 15: Redis with Python

## Installation

```bash
pip install redis          # synchronous
pip install redis[hiredis] # + C parser (faster)
pip install redis[asyncio] # + async support
```

---

## Connecting

```python
import redis

# Basic connection
r = redis.Redis(host="localhost", port=6379, db=0)
r = redis.Redis.from_url("redis://localhost:6379/0")
r = redis.Redis.from_url("redis://:password@localhost:6379/0")
r = redis.Redis.from_url("rediss://localhost:6380/0")  # TLS

# With connection pool (recommended for production)
pool = redis.ConnectionPool.from_url(
    "redis://localhost:6379/0",
    max_connections     = 20,
    socket_timeout      = 2.0,
    socket_connect_timeout = 2.0,
    retry_on_timeout    = True,
    health_check_interval = 30,
    decode_responses    = True,   # auto-decode bytes to str
)
r = redis.Redis(connection_pool=pool)

# Sentinel connection
from redis.sentinel import Sentinel
sentinel = Sentinel([("sentinel1", 26379), ("sentinel2", 26380)],
                    socket_timeout=0.1)
primary = sentinel.master_for("myredis", password="secret", socket_timeout=0.1)
replica = sentinel.slave_for("myredis",  password="secret", socket_timeout=0.1)
```

---

## All Data Types

```python
r = redis.Redis(decode_responses=True)

# Strings
r.set("name", "Alice")
r.get("name")                       # "Alice"
r.setex("otp", 300, "123456")       # with TTL
r.incr("counter")
r.incrby("counter", 5)
r.mset({"k1": "v1", "k2": "v2"})
r.mget(["k1", "k2"])                # ["v1", "v2"]

# Lists
r.lpush("queue", "a", "b", "c")
r.rpush("queue", "d")
r.lrange("queue", 0, -1)            # all elements
r.llen("queue")
r.lpop("queue")
r.brpop("queue", timeout=5)         # blocking

# Sets
r.sadd("tags", "python", "redis", "backend")
r.smembers("tags")
r.sismember("tags", "python")
r.scard("tags")
r.srandmember("tags", 2)
r.sinter("tags", "other_tags")

# Hashes
r.hset("user:1000", mapping={"name": "Alice", "age": "30"})
r.hget("user:1000", "name")
r.hmget("user:1000", ["name", "age"])
r.hgetall("user:1000")
r.hincrby("user:1000", "age", 1)
r.hdel("user:1000", "temp_field")

# Sorted Sets
r.zadd("leaderboard", {"alice": 9850, "bob": 7200})
r.zrevrange("leaderboard", 0, 9, withscores=True)
r.zrank("leaderboard", "alice")
r.zscore("leaderboard", "alice")
r.zincrby("leaderboard", 100, "bob")

# Streams
r.xadd("events", {"user": "alice", "action": "login"})
r.xrange("events", "-", "+")
r.xlen("events")
```

---

## Pipelining

```python
# Non-transactional pipeline (fire-and-forget batch)
pipe = r.pipeline(transaction=False)
for i in range(10000):
    pipe.set(f"key:{i}", f"val:{i}")
    pipe.expire(f"key:{i}", 3600)
results = pipe.execute()

# Transactional pipeline (MULTI/EXEC)
pipe = r.pipeline(transaction=True)
pipe.incr("balance")
pipe.get("balance")
ok, balance = pipe.execute()
```

---

## WATCH and Optimistic Locking

```python
def safe_increment(key: str, amount: float) -> float:
    with r.pipeline() as pipe:
        while True:
            try:
                pipe.watch(key)
                current = float(pipe.get(key) or 0)
                pipe.multi()
                pipe.set(key, current + amount)
                pipe.execute()
                return current + amount
            except redis.WatchError:
                continue  # retry on conflict
```

---

## Lua Scripts

```python
# Inline script
result = r.eval("""
    local current = redis.call('GET', KEYS[1])
    if not current then return 0 end
    return tonumber(current)
""", 1, "my_counter")

# Registered script (cached via EVALSHA)
rate_limit_script = r.register_script("""
local key    = KEYS[1]
local limit  = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = tonumber(redis.call('INCR', key))
if current == 1 then redis.call('EXPIRE', key, window) end
return current <= limit and 1 or 0
""")

allowed = rate_limit_script(keys=["rate:user-1"], args=[100, 60])
```

---

## Pub/Sub

```python
import threading

r = redis.Redis(decode_responses=True)

# Subscriber
pubsub = r.pubsub()
pubsub.subscribe("notifications")

def message_loop():
    for message in pubsub.listen():
        if message["type"] == "message":
            print(f"Received: {message['data']}")

thread = threading.Thread(target=message_loop, daemon=True)
thread.start()

# Publisher
r.publish("notifications", "Hello, subscribers!")

# Pattern subscribe
pubsub.psubscribe("user:*:events")
```

---

## Async Redis (asyncio)

```python
import asyncio
import redis.asyncio as aioredis

async def main():
    r = aioredis.Redis.from_url("redis://localhost:6379", decode_responses=True)

    # Basic operations
    await r.set("key", "value")
    val = await r.get("key")
    print(val)

    # Pipeline
    async with r.pipeline() as pipe:
        pipe.set("a", 1)
        pipe.set("b", 2)
        pipe.mget(["a", "b"])
        results = await pipe.execute()

    # Pub/Sub
    pubsub = r.pubsub()
    await pubsub.subscribe("channel")
    async for message in pubsub.listen():
        if message["type"] == "message":
            print(message["data"])

    await r.aclose()

asyncio.run(main())
```

---

## Real-World Patterns

### Cache Decorator

```python
import json
import functools
from typing import Callable

def cached(ttl: int = 3600, prefix: str = "cache"):
    def decorator(fn: Callable):
        @functools.wraps(fn)
        def wrapper(*args, **kwargs):
            key = f"{prefix}:{fn.__name__}:{args}:{sorted(kwargs.items())}"
            cached_val = r.get(key)
            if cached_val:
                return json.loads(cached_val)
            result = fn(*args, **kwargs)
            r.setex(key, ttl, json.dumps(result))
            return result
        return wrapper
    return decorator

@cached(ttl=3600, prefix="products")
def get_product(product_id: str) -> dict:
    return db.query("SELECT * FROM products WHERE id = %s", product_id)
```

### Leaderboard Service

```python
class Leaderboard:
    def __init__(self, r: redis.Redis, name: str):
        self.r    = r
        self.key  = f"leaderboard:{name}"

    def update_score(self, player: str, score: float) -> None:
        self.r.zadd(self.key, {player: score})

    def increment_score(self, player: str, amount: float) -> float:
        return self.r.zincrby(self.key, amount, player)

    def get_top(self, n: int) -> list[tuple[str, float]]:
        return self.r.zrevrange(self.key, 0, n - 1, withscores=True)

    def get_rank(self, player: str) -> int | None:
        rank = self.r.zrevrank(self.key, player)
        return rank + 1 if rank is not None else None

    def get_score(self, player: str) -> float | None:
        return self.r.zscore(self.key, player)

lb = Leaderboard(r, "weekly")
lb.update_score("alice", 9850)
lb.update_score("bob", 7200)
print(lb.get_top(5))         # [("alice", 9850), ("bob", 7200)]
print(lb.get_rank("alice"))  # 1
```


---

# Chapter 16: Redis with Node.js

## Libraries

| Library | Notes |
|---|---|
| `ioredis` | Feature-rich, Cluster + Sentinel, TypeScript support — recommended |
| `node-redis` | Official Redis client (v4+), modern async API |

```bash
npm install ioredis
npm install redis      # node-redis v4
```

---

## ioredis: Connecting

```typescript
import Redis from "ioredis";

// Single instance
const redis = new Redis({
  host:            "localhost",
  port:            6379,
  password:        "secret",
  db:              0,
  retryStrategy:   (times) => Math.min(times * 50, 2000),
  lazyConnect:     false,
});

// From URL
const redis = new Redis("redis://:password@localhost:6379/0");
const redis = new Redis("rediss://localhost:6380/0");  // TLS

// Sentinel
const redis = new Redis({
  sentinels: [
    { host: "sentinel1", port: 26379 },
    { host: "sentinel2", port: 26380 },
  ],
  name:         "myredis",
  password:     "secret",
  sentinelPassword: "sentinel-password",
});

// Cluster
const cluster = new Redis.Cluster([
  { host: "node1", port: 7001 },
  { host: "node2", port: 7002 },
  { host: "node3", port: 7003 },
]);
```

---

## All Data Types (ioredis)

```typescript
const r = new Redis({ host: "localhost" });

// Strings
await r.set("name", "Alice");
await r.get("name");                              // "Alice"
await r.setex("session:abc", 3600, "user-data"); // with TTL
await r.incr("counter");
await r.mset("k1", "v1", "k2", "v2");
await r.mget("k1", "k2");                        // ["v1", "v2"]

// Lists
await r.lpush("queue", "a", "b", "c");
await r.lrange("queue", 0, -1);                  // ["c","b","a"]
await r.brpop("queue", 5);                       // blocking pop

// Sets
await r.sadd("tags", "python", "redis");
await r.smembers("tags");
await r.sismember("tags", "python");             // 1

// Hashes
await r.hset("user:1000", "name", "Alice", "age", "30");
await r.hget("user:1000", "name");               // "Alice"
await r.hgetall("user:1000");                    // { name: "Alice", age: "30" }
await r.hmget("user:1000", ["name", "age"]);

// Sorted Sets
await r.zadd("leaderboard", 9850, "alice", 7200, "bob");
await r.zrevrange("leaderboard", 0, 9, "WITHSCORES");
await r.zrank("leaderboard", "alice");
await r.zscore("leaderboard", "alice");          // "9850"

// Expiry
await r.expire("key", 300);
await r.ttl("key");
await r.persist("key");
```

---

## Pipelining

```typescript
// Non-transactional pipeline
const pipe = r.pipeline();
for (let i = 0; i < 1000; i++) {
  pipe.set(`key:${i}`, `value:${i}`);
  pipe.expire(`key:${i}`, 3600);
}
const results = await pipe.exec();
// results: [[null, "OK"], [null, 1], [null, "OK"], ...]

// Transactional pipeline (MULTI/EXEC)
const [err, results2] = await r
  .multi()
  .incr("counter")
  .get("counter")
  .exec();
```

---

## Transactions with WATCH

```typescript
async function transferBalance(
  fromKey: string,
  toKey: string,
  amount: number
): Promise<boolean> {
  const result = await r.watch(fromKey, toKey, async (err, multi) => {
    const from = parseFloat((await r.get(fromKey)) ?? "0");
    const to   = parseFloat((await r.get(toKey))   ?? "0");

    if (from < amount) return null;

    return multi
      .set(fromKey, (from - amount).toString())
      .set(toKey,   (to   + amount).toString())
      .exec();
  });

  return result !== null;
}
```

---

## Lua Scripts

```typescript
// Define and call a Lua script
const rateLimitScript = r.defineCommand("rateLimit", {
  numberOfKeys: 1,
  lua: `
    local current = tonumber(redis.call('INCR', KEYS[1]))
    if current == 1 then
      redis.call('EXPIRE', KEYS[1], ARGV[2])
    end
    return current <= tonumber(ARGV[1]) and 1 or 0
  `,
});

// Type-safe call
const allowed = await (r as any).rateLimit("rate:user-1", 100, 60);

// Or use EVAL directly
const result = await r.eval(
  `return redis.call('GET', KEYS[1])`,
  1,
  "my-key"
);
```

---

## Pub/Sub

```typescript
import Redis from "ioredis";

// Use separate connections for pub and sub
const pub = new Redis();
const sub = new Redis();

// Subscribe
sub.subscribe("notifications", "alerts", (err, count) => {
  console.log(`Subscribed to ${count} channels`);
});

sub.on("message", (channel: string, message: string) => {
  console.log(`[${channel}] ${message}`);
});

// Pattern subscribe
sub.psubscribe("user:*:events");
sub.on("pmessage", (pattern, channel, message) => {
  console.log(`Pattern ${pattern} → [${channel}] ${message}`);
});

// Publish
await pub.publish("notifications", JSON.stringify({ type: "alert", msg: "Server down" }));

// Unsubscribe
sub.unsubscribe("notifications");
```

---

## node-redis (Official v4)

```typescript
import { createClient } from "redis";

const client = createClient({
  url:      "redis://localhost:6379",
  password: "secret",
});

client.on("error", (err) => console.error("Redis error:", err));
await client.connect();

// Commands
await client.set("key", "value");
const val = await client.get("key");

// Hash
await client.hSet("user:1", { name: "Alice", age: "30" });
const user = await client.hGetAll("user:1");

// Sorted Set
await client.zAdd("lb", [{ score: 9850, value: "alice" }]);
const top = await client.zRangeWithScores("lb", 0, 9, { REV: true });

// Pipeline
const multi = client.multi();
multi.set("a", "1");
multi.incr("a");
const results = await multi.exec();

await client.disconnect();
```

---

## Real-World Patterns

### Express.js Session Middleware

```typescript
import express from "express";
import session from "express-session";
import { createClient } from "redis";
import connectRedis from "connect-redis";

const app     = express();
const client  = createClient({ url: "redis://localhost:6379" });
await client.connect();

const RedisStore = connectRedis(session);

app.use(session({
  store:  new RedisStore({ client }),
  secret: "your-session-secret",
  resave: false,
  saveUninitialized: false,
  cookie: { secure: true, maxAge: 3600000 },
}));

app.get("/login", (req, res) => {
  (req.session as any).userId = "user-123";
  res.json({ status: "logged in" });
});
```

### Cache Middleware

```typescript
import { Request, Response, NextFunction } from "express";

function cacheMiddleware(ttl: number) {
  return async (req: Request, res: Response, next: NextFunction) => {
    const key = `cache:${req.path}:${JSON.stringify(req.query)}`;
    const cached = await r.get(key);

    if (cached) {
      return res.json(JSON.parse(cached));
    }

    const originalJson = res.json.bind(res);
    res.json = (data: any) => {
      r.setex(key, ttl, JSON.stringify(data));
      return originalJson(data);
    };

    next();
  };
}

app.get("/api/products", cacheMiddleware(300), async (req, res) => {
  const products = await db.getProducts();
  res.json(products);
});
```


---

# Chapter 17: Redis Modules

## What are Redis Modules?

Redis modules extend Redis with new data types and commands. Loaded as shared libraries, they run in-process with Redis for maximum performance.

```bash
# Load a module in redis.conf
loadmodule /usr/lib/redis/modules/redisearch.so
loadmodule /usr/lib/redis/modules/redisjson.so
loadmodule /usr/lib/redis/modules/redistimeseries.so
loadmodule /usr/lib/redis/modules/redisbloom.so

# Or at runtime
MODULE LOAD /path/to/module.so
MODULE LIST
MODULE UNLOAD module-name
```

**Redis Stack** (Docker) includes all major modules:
```bash
docker run -d -p 6379:6379 redis/redis-stack-server:latest
```

---

## RediSearch (Full-Text Search)

RediSearch adds full-text search, secondary indexing, and aggregations to Redis.

### Create an Index

```bash
# Index a Hash-based schema
FT.CREATE products:idx
  ON HASH
  PREFIX 1 product:
  SCHEMA
    name        TEXT    WEIGHT 5.0
    description TEXT
    category    TAG     SEPARATOR ","
    price       NUMERIC SORTABLE
    rating      NUMERIC SORTABLE
    in_stock    TAG

# Index a JSON schema (with RedisJSON)
FT.CREATE orders:idx
  ON JSON
  PREFIX 1 order:
  SCHEMA
    $.user_id   AS user_id   TAG
    $.amount    AS amount    NUMERIC SORTABLE
    $.status    AS status    TAG
    $.item_name AS item_name TEXT
```

### Basic Search

```bash
# Simple text search
FT.SEARCH products:idx "laptop"

# Filter by tag and range
FT.SEARCH products:idx "@category:{electronics} @price:[100 500]"

# Full-text + filters
FT.SEARCH products:idx "gaming @category:{computers} @in_stock:{true}"

# Sort and paginate
FT.SEARCH products:idx "phone" SORTBY price ASC LIMIT 0 10

# With highlighting
FT.SEARCH products:idx "redis" HIGHLIGHT FIELDS 1 description

# Return specific fields
FT.SEARCH products:idx "*" RETURN 2 name price

# Count only
FT.SEARCH products:idx "laptop" LIMIT 0 0
```

### Aggregations

```bash
# Average price per category
FT.AGGREGATE products:idx "*"
  GROUPBY 1 @category
  REDUCE AVG 1 @price AS avg_price
  REDUCE COUNT 0 AS count
  SORTBY 2 @avg_price DESC

# Revenue by status
FT.AGGREGATE orders:idx "*"
  GROUPBY 1 @status
  REDUCE SUM 1 @amount AS total_revenue
  REDUCE COUNT 0 AS order_count
```

### Python with RediSearch

```python
from redis.commands.search.field import TextField, TagField, NumericField
from redis.commands.search.indexDefinition import IndexDefinition, IndexType
from redis.commands.search.query import Query

# Create index
r.ft("products:idx").create_index(
    [
        TextField("name",        weight=5.0),
        TextField("description"),
        TagField("category"),
        NumericField("price",    sortable=True),
    ],
    definition = IndexDefinition(prefix=["product:"], index_type=IndexType.HASH)
)

# Add documents
r.hset("product:1", mapping={"name": "Redis Laptop", "category": "electronics", "price": "999"})

# Search
results = r.ft("products:idx").search(
    Query("laptop @price:[500 1500]")
    .sort_by("price", asc=True)
    .paging(0, 10)
)
for doc in results.docs:
    print(doc.name, doc.price)
```

---

## RedisJSON

Store and query JSON documents natively in Redis. Uses JSONPath syntax.

```bash
# Set a JSON document
JSON.SET user:1000 $ '{"name":"Alice","age":30,"address":{"city":"NYC","zip":"10001"},"tags":["vip","premium"]}'

# Get full document
JSON.GET user:1000 $

# Get nested field
JSON.GET user:1000 $.name
JSON.GET user:1000 $.address.city
JSON.GET user:1000 $.tags[0]

# Update a field
JSON.SET user:1000 $.age 31
JSON.SET user:1000 $.address.city '"Los Angeles"'

# Increment numeric field
JSON.NUMINCRBY user:1000 $.age 1

# Array operations
JSON.ARRAPPEND user:1000 $.tags '"elite"'
JSON.ARRLEN    user:1000 $.tags      # 3
JSON.ARRINDEX  user:1000 $.tags '"vip"'  # 0

# Delete a field
JSON.DEL user:1000 $.address.zip

# Type of a field
JSON.TYPE user:1000 $.age      # integer
JSON.TYPE user:1000 $.tags     # array
```

```python
import json

# Set
r.json().set("user:1000", "$", {"name": "Alice", "age": 30})

# Get
doc = r.json().get("user:1000", "$")     # returns list
doc = r.json().get("user:1000", "$.name") # ["Alice"]

# Update
r.json().set("user:1000", "$.age", 31)
r.json().numincrby("user:1000", "$.age", 1)

# Array append
r.json().arrappend("user:1000", "$.tags", "premium")
```

---

## RedisTimeSeries

Purpose-built time series database inside Redis. Perfect for metrics, IoT, and monitoring data.

```bash
# Create a time series
TS.CREATE sensor:temp:room1
  RETENTION 86400000        # 1 day in ms
  LABELS location room1 type temperature

# Add data points
TS.ADD sensor:temp:room1 * 22.5         # * = current timestamp
TS.ADD sensor:temp:room1 1704067200000 23.1
TS.MADD sensor:temp:room1 1704067260000 22.8
       sensor:temp:room2 1704067260000 21.4

# Query range
TS.RANGE sensor:temp:room1 - +          # all data
TS.RANGE sensor:temp:room1 1704067200000 1704153600000  # last 24h

# Aggregation (average per 1 minute)
TS.RANGE sensor:temp:room1 - + AGGREGATION avg 60000

# Latest value
TS.GET sensor:temp:room1

# Multi-query with label filter
TS.MRANGE - + FILTER type=temperature AGGREGATION avg 60000
```

```python
r.ts().create("sensor:temp", retention_msecs=86400000, labels={"room": "1"})
r.ts().add("sensor:temp", "*", 22.5)
r.ts().range("sensor:temp", "-", "+")
r.ts().range("sensor:temp", "-", "+", aggregation_type="avg", bucket_size_msec=60000)
```

---

## RedisBloom

Probabilistic data structures: Bloom Filter, Cuckoo Filter, Count-Min Sketch, Top-K.

### Bloom Filter

Answers "has this item been seen before?" with configurable false-positive rate. Uses ~10 bits per element.

```bash
# Create filter (1M items, 0.1% false positive rate)
BF.RESERVE emails:seen 0.001 1000000

# Add items
BF.ADD    emails:seen "user@example.com"
BF.MADD   emails:seen "a@b.com" "c@d.com" "e@f.com"

# Check membership
BF.EXISTS  emails:seen "user@example.com"    # 1 (probably yes)
BF.EXISTS  emails:seen "new@user.com"        # 0 (definitely no)
BF.MEXISTS emails:seen "a@b.com" "new@b.com"

# Info
BF.INFO emails:seen
```

```python
r.bf().reserve("urls:crawled", 0.001, 1000000)
r.bf().add("urls:crawled", "https://example.com")
seen = r.bf().exists("urls:crawled", "https://example.com")  # True
```

### Count-Min Sketch

Approximate frequency counter — count occurrences with small memory:

```bash
CMS.INITBYPROB searches:freq 0.001 0.01
CMS.INCRBY searches:freq "redis" 1 "python" 3
CMS.QUERY  searches:freq "redis"     # approximate count
```

### Top-K

Track the K most frequent items:

```bash
TOPK.RESERVE trending:hashtags 10 50 5 0.9
TOPK.ADD    trending:hashtags "redis" "python" "kafka" "redis" "redis"
TOPK.LIST   trending:hashtags WITHCOUNT
TOPK.QUERY  trending:hashtags "redis"
```

---

## Module Availability

| Module | Redis OSS | Redis Stack | Redis Enterprise |
|---|---|---|---|
| RediSearch | No | Yes | Yes |
| RedisJSON | No | Yes | Yes |
| RedisTimeSeries | No | Yes | Yes |
| RedisBloom | No | Yes | Yes |
| RedisGraph | Deprecated | No | Yes |
| RedisAI | No | No | Yes |

**For development:** use `redis/redis-stack-server` Docker image — all modules included.


---

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


---

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


---

# Chapter 20: Learning Roadmap

## 5-Week Hands-On Plan

This roadmap takes you from zero to production-ready Redis skills with a real project each week.

---

## Week 1: Foundation

**Goal:** Run Redis locally, master all 5 core data structures, understand TTL and eviction.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Install Redis with Docker, use redis-cli, run PING/SET/GET |
| 2 | Practice all String commands (INCR, MSET, GETEX, SETEX) |
| 3 | Practice List, Set, and Hash commands with real use cases |
| 4 | Practice Sorted Set (leaderboard), understand ZRANK vs ZSCORE |
| 5 | Set TTLs, observe expiry, configure eviction policies |

### Practice Project: Social App Data Layer

```python
import redis
import json
import time

r = redis.Redis(decode_responses=True)

# 1. User profiles (Hash)
def create_user(user_id: str, name: str, email: str) -> None:
    r.hset(f"user:{user_id}", mapping={
        "name":       name,
        "email":      email,
        "created_at": int(time.time()),
        "post_count": 0,
    })

# 2. Activity feed (List — newest first)
def add_activity(user_id: str, activity: str) -> None:
    r.lpush(f"feed:{user_id}", json.dumps({
        "text": activity,
        "ts":   int(time.time()),
    }))
    r.ltrim(f"feed:{user_id}", 0, 99)   # keep last 100

def get_feed(user_id: str, page: int = 0) -> list:
    start = page * 20
    items = r.lrange(f"feed:{user_id}", start, start + 19)
    return [json.loads(i) for i in items]

# 3. Follower graph (Sets)
def follow(user_id: str, target_id: str) -> None:
    r.sadd(f"following:{user_id}", target_id)
    r.sadd(f"followers:{target_id}", user_id)

def mutual_friends(user_a: str, user_b: str) -> set:
    return r.sinter(f"following:{user_a}", f"following:{user_b}")

# 4. Leaderboard (Sorted Set)
def record_score(user_id: str, points: int) -> None:
    r.zincrby("leaderboard:weekly", points, user_id)

def get_top_10() -> list:
    return r.zrevrange("leaderboard:weekly", 0, 9, withscores=True)

# 5. Session (String with TTL)
def create_session(user_id: str) -> str:
    import uuid
    token = str(uuid.uuid4())
    r.setex(f"session:{token}", 3600, user_id)
    return token
```

### Week 1 Checklist
- [ ] Redis running in Docker
- [ ] All 5 core data types used in code
- [ ] TTL set and verified with TTL command
- [ ] Eviction policy configured
- [ ] 5+ real use cases implemented

---

## Week 2: Caching & Patterns

**Goal:** Build a production-quality caching layer with cache-aside, stampede prevention, and rate limiting.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Implement cache-aside for a database query |
| 2 | Add TTL jitter and stampede prevention (mutex lock) |
| 3 | Build a sliding-window rate limiter with Sorted Set |
| 4 | Implement a distributed lock (Lua-based) |
| 5 | Add keyspace notifications to trigger cache refresh |

### Practice: API Gateway Cache + Rate Limiter

```python
import time, redis, json, uuid

r = redis.Redis(decode_responses=True)

# Cache with stampede protection
def get_product_cached(product_id: str) -> dict:
    key      = f"cache:product:{product_id}"
    lock_key = f"lock:product:{product_id}"

    data = r.get(key)
    if data:
        return json.loads(data)

    # Try to acquire lock
    if r.set(lock_key, "1", nx=True, ex=10):
        try:
            product = db_fetch_product(product_id)
            r.setex(key, 3600, json.dumps(product))
            return product
        finally:
            r.delete(lock_key)
    else:
        time.sleep(0.05)
        return get_product_cached(product_id)  # retry

# Sliding window rate limiter
def rate_limit(user_id: str, limit: int = 100, window: int = 60) -> bool:
    key = f"rate:{user_id}"
    now = time.time()
    pipe = r.pipeline()
    pipe.zremrangebyscore(key, 0, now - window)
    pipe.zadd(key, {str(now): now})
    pipe.zcard(key)
    pipe.expire(key, window)
    _, _, count, _ = pipe.execute()
    return count <= limit
```

### Week 2 Checklist
- [ ] Cache-aside pattern implemented and tested
- [ ] Cache stampede prevented with mutex lock
- [ ] Rate limiter working (fixed and sliding window)
- [ ] Distributed lock implemented and tested
- [ ] Keyspace notification triggers refresh

---

## Week 3: Pub/Sub, Streams, & Scripting

**Goal:** Build a real-time notification system and event pipeline.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Implement Pub/Sub notification system |
| 2 | Switch to Redis Streams for durable events |
| 3 | Add consumer groups with acknowledgment |
| 4 | Write Lua scripts for complex atomic operations |
| 5 | Build a full event pipeline with DLQ |

### Practice: Real-Time Order Pipeline

```python
# Producer: place orders
def place_order(order_data: dict) -> str:
    msg_id = r.xadd("orders:stream", order_data, maxlen=50000, approximate=True)
    r.publish("orders:notifications", json.dumps({"id": msg_id, **order_data}))
    return msg_id

# Consumer with consumer group
r.xgroup_create("orders:stream", "payment-processor", "$", mkstream=True)

def process_orders():
    while True:
        msgs = r.xreadgroup("payment-processor", "worker-1",
                            {"orders:stream": ">"}, count=10, block=2000)
        for _, messages in (msgs or []):
            for msg_id, fields in messages:
                try:
                    charge_customer(fields)
                    r.xack("orders:stream", "payment-processor", msg_id)
                except Exception:
                    r.xadd("orders:dlq", {"original_id": msg_id, **fields})
                    r.xack("orders:stream", "payment-processor", msg_id)
```

### Week 3 Checklist
- [ ] Pub/Sub fan-out to 3+ subscribers
- [ ] Redis Streams producer + consumer
- [ ] Consumer group with multiple workers
- [ ] Message acknowledgment and DLQ
- [ ] Lua atomic script for complex operation

---

## Week 4: HA & Performance

**Goal:** Set up Sentinel, tune performance, add monitoring.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Set up 1 primary + 2 replicas with Docker Compose |
| 2 | Add 3 Sentinel nodes, test automatic failover |
| 3 | Benchmark with redis-benchmark, add pipelining |
| 4 | Set up Prometheus + Grafana with redis_exporter |
| 5 | Configure ACLs, disable dangerous commands, add TLS |

### Week 4 Checklist
- [ ] Replication working (1 primary, 2 replicas)
- [ ] Sentinel automatic failover tested
- [ ] 2× throughput improvement from pipelining
- [ ] Grafana dashboard showing memory, ops/sec, hit rate
- [ ] ACL users configured (admin, app-user, read-only)
- [ ] Dangerous commands disabled

---

## Week 5: Capstone Project

**Goal:** Build a complete Redis-powered production system.

### Capstone Option A: High-Performance API Backend

```
Frontend → API Gateway (rate limiter) → App → Redis Cache → PostgreSQL
                                            → Redis Sessions
                                            → Redis Pub/Sub (real-time)
                                            → Redis Leaderboards
```

Features:
- Session management with sliding expiry
- Cache-aside with stampede protection
- Rate limiting (100 req/min per user)
- Real-time leaderboard with daily/weekly/all-time buckets
- Pub/Sub for live score updates to WebSocket clients

### Capstone Option B: Real-Time Event Processing

```
Event Producers → Redis Streams → Consumer Groups → Multiple Sinks
                                                  → PostgreSQL (persist)
                                                  → Elasticsearch (search)
                                                  → Redis PubSub (notify)
```

Features:
- High-throughput event ingestion (Streams + MAXLEN)
- Multiple consumer groups (payment, analytics, notification)
- Dead letter queue for failed events
- Consumer lag monitoring dashboard

### Capstone Option C: Distributed Task System

```
Job API → Redis List (queue) → Workers → Redis Hash (status) → Result Store
       → Redis Sorted Set (priority queue)
       → Redis Pub/Sub (progress updates)
       → Redis Streams (audit log)
```

Features:
- Priority queue (Sorted Set, score = priority)
- Job status tracking (Hash per job)
- Distributed locking (prevent duplicate processing)
- Progress streaming via Pub/Sub
- Retry logic with backoff (Sorted Set as delayed queue)

---

## Key Commands Reference Card

```bash
# Connection
redis-cli -h host -p 6379 -a password -n 0

# Diagnostics
PING
INFO memory
INFO stats
INFO replication
DBSIZE
CONFIG GET maxmemory
CONFIG GET maxmemory-policy
SLOWLOG GET 10

# Keyspace navigation
SCAN 0 MATCH "user:*" COUNT 100
TYPE key
TTL key
MEMORY USAGE key
OBJECT ENCODING key

# Cluster / Sentinel
redis-cli -c -p 7001          # cluster mode
redis-cli -p 26379 SENTINEL masters
redis-cli -p 26379 SENTINEL get-master-addr-by-name myredis
redis-cli -p 26379 SENTINEL failover myredis

# Benchmarking
redis-benchmark -n 100000 -c 50 -t set,get -q
redis-benchmark -n 1000000 -c 100 -P 16 -t set -q
```

---

## Beyond the Roadmap

| Topic | Notes |
|---|---|
| **Redisson** | Feature-rich Java Redis client (distributed locks, collections) |
| **RedisOM** | Object mapper for Redis with search (Python/Java/Node) |
| **Celery + Redis** | Distributed task queue in Python |
| **Bull/BullMQ** | Redis-backed job queues for Node.js |
| **Redis Gears** | Server-side serverless functions (deprecated → Functions) |
| **Redis Functions** | Persistent Lua function libraries (Redis 7.0+) |
| **Keydb** | Multi-threaded Redis fork (higher throughput) |
| **Dragonfly** | Modern Redis alternative with better multi-core utilization |
| **RESP3 Protocol** | New protocol with richer types (Redis 6+) |
| **Client-side caching** | Track-and-invalidate with server-assisted invalidation |
