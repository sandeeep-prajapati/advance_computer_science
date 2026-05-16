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
