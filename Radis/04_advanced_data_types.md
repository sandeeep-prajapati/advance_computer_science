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
