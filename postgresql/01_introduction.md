# Chapter 1: Introduction to PostgreSQL

## What is PostgreSQL?

**PostgreSQL** (often called "Postgres") is a powerful, open-source **object-relational database management system (ORDBMS)** with over 35 years of active development. It is known for reliability, feature robustness, and standards compliance.

> **Analogy:** If databases were cities, SQLite is a small town, MySQL is a mid-size city, and PostgreSQL is a full metropolis with every infrastructure system you could need — and it's still growing.

---

## Brief History

| Year | Milestone |
|---|---|
| 1986 | POSTGRES project begins at UC Berkeley (Michael Stonebraker) |
| 1994 | SQL support added → Postgres95 |
| 1996 | Renamed to **PostgreSQL**, open-source release |
| 2005 | 8.0 — native Windows support, point-in-time recovery |
| 2010 | 9.0 — streaming replication, hot standby |
| 2016 | 9.6 — parallel queries |
| 2017 | 10 — declarative partitioning, logical replication |
| 2019 | 12 — generated columns, improved partitioning |
| 2021 | 14 — range types, JSON subscripting |
| 2023 | 16 — logical replication from standby, pg_stat improvements |

---

## Why PostgreSQL?

| Feature | PostgreSQL |
|---|---|
| **ACID compliance** | Full — every transaction is safe |
| **SQL standards** | Most compliant of any RDBMS |
| **Data types** | 100+ built-in + custom types |
| **JSON support** | First-class JSONB — better than most NoSQL |
| **Full-text search** | Built-in, no Elasticsearch needed for many use cases |
| **Extensibility** | Write extensions in C, PL/pgSQL, Python, Rust |
| **Scalability** | Handles terabytes, petabytes with partitioning |
| **License** | PostgreSQL License (very permissive — free for commercial use) |
| **Community** | Active, large, no single corporate owner |

---

## PostgreSQL vs Alternatives

| Feature | PostgreSQL | MySQL | SQLite | MongoDB |
|---|---|---|---|---|
| Type | ORDBMS | RDBMS | Embedded RDBMS | Document DB |
| ACID | Full | Mostly (InnoDB) | Yes | Yes (4.0+) |
| JSON | JSONB (excellent) | JSON (limited) | As text | Native documents |
| Full-text search | Built-in | Basic | No | Yes |
| Replication | Streaming + logical | Binary + GTID | No | Replica sets |
| Partitioning | Declarative | Limited | No | Sharding |
| Extensions | Rich ecosystem | Limited | FTS only | Plugins |
| Geospatial | PostGIS | Basic | No | Basic |
| Window functions | Full | Partial | Partial | No |

---

## Core RDBMS Concepts

### ACID Properties
Every transaction in PostgreSQL guarantees:

- **Atomicity** — all steps in a transaction succeed, or none do
- **Consistency** — data always moves from one valid state to another
- **Isolation** — concurrent transactions don't interfere with each other
- **Durability** — committed data survives system crashes

### Relational Model
Data is stored in **tables** (relations) with **rows** (tuples) and **columns** (attributes). Tables relate to each other via **foreign keys**.

```
users                    orders
┌────┬──────┬──────┐    ┌────┬─────────┬──────┐
│ id │ name │ email│    │ id │ user_id │total │
├────┼──────┼──────┤    ├────┼─────────┼──────┤
│  1 │ Alice│ a@.. │    │  1 │    1    │ 99.9 │
│  2 │ Bob  │ b@.. │    │  2 │    1    │ 45.0 │
└────┴──────┴──────┘    │  3 │    2    │120.0 │
                        └────┴─────────┴──────┘
                              ↑ foreign key → users.id
```

---

## PostgreSQL Architecture

```
Client (psql / app)
        │ TCP/IP or Unix socket
        ▼
┌───────────────────┐
│   postmaster      │  ← supervisor process; forks backend per connection
└────────┬──────────┘
         │ forks
┌────────▼──────────┐
│   backend process │  ← one per client connection; runs queries
│   (query engine)  │
│   ┌─────────────┐ │
│   │ Parser      │ │  parse SQL → AST
│   │ Rewriter    │ │  apply view rules
│   │ Planner     │ │  choose optimal plan (EXPLAIN output)
│   │ Executor    │ │  run the plan
│   └─────────────┘ │
└────────┬──────────┘
         │
┌────────▼──────────────────────────────────────────┐
│   Shared Memory                                    │
│   ┌─────────────┐  ┌──────────────┐  ┌──────────┐ │
│   │ shared_buffers│ │ WAL buffers  │  │ lock table│ │
│   │ (page cache) │  │              │  │           │ │
│   └─────────────┘  └──────────────┘  └──────────┘ │
└────────────────────────────────────────────────────┘
         │
┌────────▼──────────┐
│   Disk Storage     │
│   ├── base/        │  table & index data files
│   ├── pg_wal/      │  Write-Ahead Log
│   └── pg_stat/     │  statistics files
└────────────────────┘
```

---

## MVCC — Multi-Version Concurrency Control

PostgreSQL uses MVCC to allow readers and writers to operate concurrently without blocking each other.

- Every row has **xmin** (transaction that created it) and **xmax** (transaction that deleted it)
- A query sees rows where `xmin <= snapshot_xid < xmax`
- Writers create new row versions; old versions become "dead tuples"
- **VACUUM** reclaims dead tuple space

```sql
-- Readers never block writers
-- Writers never block readers
-- Writers only block other writers (on the same row)
```

---

## Real-World Use Cases

- **Web applications** — the most popular choice alongside MySQL
- **Financial systems** — strict ACID guarantees, audit trails
- **Geospatial applications** — PostGIS extension
- **Time-series data** — TimescaleDB extension
- **Analytics** — window functions, CTEs, parallel query
- **Document storage** — JSONB as a MongoDB alternative
- **Full-text search** — replace or supplement Elasticsearch
- **Multi-tenant SaaS** — row-level security per tenant

---

## Who Uses PostgreSQL?

Apple, Instagram, Reddit, Skype, NASA, Twitch, Spotify, Uber (historically), GitHub, GitLab, Heroku, Supabase, and hundreds of thousands of companies worldwide.
