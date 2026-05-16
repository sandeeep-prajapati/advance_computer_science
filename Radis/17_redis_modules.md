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
