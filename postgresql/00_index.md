# PostgreSQL Developer Book

A comprehensive, chapter-wise guide to mastering PostgreSQL — from first database to production-grade performance tuning, replication, and cloud deployment.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_introduction.md) | **Introduction to PostgreSQL** | What is PostgreSQL, history, RDBMS concepts, strengths, use cases |
| [02](02_setup_and_installation.md) | **Setup & Installation** | Install on Linux/Mac/Windows, psql CLI, pgAdmin, Docker setup |
| [03](03_databases_schemas_tables.md) | **Databases, Schemas & Tables** | CREATE DATABASE, schemas, data types, constraints, ALTER/DROP |
| [04](04_crud_operations.md) | **CRUD Operations** | SELECT, INSERT, UPDATE, DELETE, RETURNING, UPSERT (ON CONFLICT) |
| [05](05_advanced_queries.md) | **Advanced Queries** | JOINs, subqueries, CTEs, window functions, GROUPING SETS |
| [06](06_indexes.md) | **Indexes & Performance** | B-tree, Hash, GIN, GiST, BRIN, partial, composite indexes, EXPLAIN |
| [07](07_transactions.md) | **Transactions & ACID** | BEGIN/COMMIT/ROLLBACK, isolation levels, savepoints, locking |
| [08](08_functions_procedures_triggers.md) | **Functions, Procedures & Triggers** | PL/pgSQL, stored functions, procedures, triggers, event triggers |
| [09](09_views_and_materialized_views.md) | **Views & Materialized Views** | CREATE VIEW, updatable views, MATERIALIZED VIEW, REFRESH |
| [10](10_full_text_search.md) | **Full-Text Search** | tsvector, tsquery, GIN index, ranking, dictionaries, phrase search |
| [11](11_json_and_jsonb.md) | **JSON & JSONB** | JSON vs JSONB, operators, functions, indexing, document patterns |
| [12](12_partitioning.md) | **Table Partitioning** | Range, list, hash partitioning, partition pruning, sub-partitions |
| [13](13_replication.md) | **Replication & High Availability** | Streaming replication, logical replication, failover, Patroni |
| [14](14_security.md) | **Security & User Management** | Roles, privileges, row-level security, pg_hba.conf, SSL |
| [15](15_backup_and_recovery.md) | **Backup & Recovery** | pg_dump, pg_basebackup, PITR, WAL archiving, pgBackRest |
| [16](16_postgresql_with_docker.md) | **PostgreSQL with Docker & K8s** | Docker Compose, init scripts, volumes, CloudNativePG operator |
| [17](17_extensions.md) | **Extensions** | pg_stat_statements, PostGIS, pgcrypto, uuid-ossp, TimescaleDB |
| [18](18_connection_pooling.md) | **Connection Pooling** | PgBouncer, pg_pool, transaction vs session pooling, Supabase |
| [19](19_performance_tuning.md) | **Performance Tuning & Monitoring** | postgresql.conf, VACUUM, ANALYZE, pg_stat_*, query optimization |
| [20](20_learning_roadmap.md) | **Learning Roadmap** | 6-week hands-on plan from first query to production DBA skills |

---

## How to Read This Book

**If you're brand new to PostgreSQL:** Read 1 → 2 → 3 → 4 → 5 in order. This gives you the core mental model.

**If you're writing application queries:** Read 4 → 5 → 6 → 7.

**If you're handling JSON/documents:** Read 11 → 5 → 6.

**If you're optimizing performance:** Read 6 → 19 → 7 → 12.

**If you're going to production:** Read 13 → 14 → 15 → 16 → 18.

**If you need search:** Jump to Chapter 10 (Full-Text Search).

---

## Key Concepts at a Glance

```
PostgreSQL Server
├── Cluster (one data directory, multiple databases)
│   ├── Database
│   │   ├── Schema (public, app, audit, ...)
│   │   │   ├── Tables
│   │   │   ├── Views / Materialized Views
│   │   │   ├── Functions / Procedures
│   │   │   ├── Sequences
│   │   │   ├── Indexes
│   │   │   └── Triggers
│   │   └── Extensions (pg_stat_statements, PostGIS, ...)
│   └── pg_catalog (system catalog — metadata about everything)
│
├── Process Architecture
│   ├── postmaster          (supervisor, accepts connections)
│   ├── backend process     (one per client connection)
│   ├── WAL writer          (writes WAL to disk)
│   ├── checkpointer        (flushes dirty pages to disk)
│   ├── autovacuum          (reclaims dead tuples)
│   └── background worker   (extensions, logical replication)
│
├── Storage
│   ├── Heap files          (table data in 8KB pages)
│   ├── WAL files           (Write-Ahead Log — crash recovery)
│   ├── Index files         (B-tree, GIN, GiST, etc.)
│   └── TOAST              (large values stored separately)
│
└── MVCC (Multi-Version Concurrency Control)
    ├── Each row has xmin/xmax (transaction visibility)
    ├── Readers don't block writers
    └── Dead tuples reclaimed by VACUUM
```
