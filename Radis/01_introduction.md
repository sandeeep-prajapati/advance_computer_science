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
