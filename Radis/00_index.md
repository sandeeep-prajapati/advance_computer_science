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
