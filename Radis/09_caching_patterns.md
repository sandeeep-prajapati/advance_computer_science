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
