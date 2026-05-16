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
