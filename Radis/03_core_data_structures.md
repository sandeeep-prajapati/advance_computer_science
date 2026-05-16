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
