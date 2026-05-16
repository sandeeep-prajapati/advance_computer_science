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
