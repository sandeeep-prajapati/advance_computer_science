# Chapter 20: Learning Roadmap

## Overview

This roadmap takes you from zero to production-ready PostgreSQL skills in 6 weeks, with hands-on projects at each stage.

```
Week 1: Foundations         → Can write SQL and design schemas
Week 2: Advanced Queries    → Can analyze data with complex SQL
Week 3: Indexing & Tuning   → Can diagnose and fix slow queries
Week 4: PL/pgSQL & Security → Can write functions and secure databases
Week 5: Reliability         → Can set up replication and backups
Week 6: Production Systems  → Can run PostgreSQL in Docker/K8s at scale
```

---

## Week 1: Foundations

**Goal:** Be comfortable with basic SQL, schema design, and the psql CLI.

### Topics:
- Install PostgreSQL and explore psql CLI (`\l`, `\dt`, `\d tablename`, `\x`, `\timing`)
- Data types: TEXT, INTEGER, BIGINT, NUMERIC, BOOLEAN, TIMESTAMPTZ, UUID
- CREATE TABLE with PRIMARY KEY, NOT NULL, UNIQUE, DEFAULT, CHECK
- CRUD: SELECT, INSERT, UPDATE, DELETE, RETURNING
- Aggregate functions: COUNT, SUM, AVG, MIN, MAX, GROUP BY, HAVING
- JOINs: INNER, LEFT, RIGHT, FULL OUTER
- ORDER BY, LIMIT, OFFSET

### Project: Library Database
```sql
-- Design and implement a library system:
CREATE TABLE authors (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    birth_year INTEGER,
    nationality TEXT
);

CREATE TABLE books (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author_id BIGINT REFERENCES authors(id),
    isbn TEXT UNIQUE,
    published_year INTEGER,
    genre TEXT,
    available_copies INTEGER NOT NULL DEFAULT 1 CHECK (available_copies >= 0)
);

CREATE TABLE members (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE loans (
    id BIGSERIAL PRIMARY KEY,
    book_id BIGINT NOT NULL REFERENCES books(id),
    member_id BIGINT NOT NULL REFERENCES members(id),
    loaned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    due_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '14 days',
    returned_at TIMESTAMPTZ
);

-- Queries to write:
-- 1. All books currently on loan (not yet returned)
-- 2. Most popular books (most loans)
-- 3. Members with overdue books
-- 4. Books by genre with available copy count
-- 5. Author with the most books in the library
```

### Milestone check:
- [ ] Can write SELECT with JOINs across 3+ tables
- [ ] Can design a normalized schema from requirements
- [ ] Comfortable with psql CLI meta-commands

---

## Week 2: Advanced Queries and JSONB

**Goal:** Use CTEs, window functions, and JSONB for real analytical queries.

### Topics:
- Subqueries (scalar, correlated, EXISTS)
- CTEs (WITH) — regular and recursive
- Window functions: ROW_NUMBER, RANK, LAG, LEAD, SUM OVER
- JSONB operators and functions
- Full-text search with tsvector/tsquery
- UPSERT (ON CONFLICT DO UPDATE)
- COPY for bulk data import

### Project: E-commerce Analytics
```sql
-- Build on Week 1 or start fresh with:
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT,
    price NUMERIC(10,2),
    metadata JSONB DEFAULT '{}'
);

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT,
    placed_at TIMESTAMPTZ DEFAULT NOW(),
    total NUMERIC(10,2),
    status TEXT
);

CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT REFERENCES orders(id),
    product_id BIGINT REFERENCES products(id),
    quantity INTEGER,
    unit_price NUMERIC(10,2)
);

-- Queries to write:
-- 1. Rolling 7-day revenue using window functions
-- 2. Customer ranking by total spend (RANK, DENSE_RANK)
-- 3. Products bought together (self-join on order_items)
-- 4. Filter products by JSONB attributes (brand, color)
-- 5. UPSERT product inventory from CSV using COPY + ON CONFLICT
-- 6. Recursive CTE: product category hierarchy
```

### Milestone check:
- [ ] Can write recursive CTEs (e.g., org chart traversal)
- [ ] Can use window functions for rankings and running totals
- [ ] Can query JSONB data with GIN index

---

## Week 3: Indexing and Query Optimization

**Goal:** Diagnose slow queries and fix them with the right indexes.

### Topics:
- `EXPLAIN` and `EXPLAIN (ANALYZE, BUFFERS)`
- B-tree, GIN, GiST, BRIN indexes
- Partial indexes, expression indexes, covering indexes (INCLUDE)
- Common slow query patterns and fixes
- `pg_stat_statements` for query analytics
- `autovacuum`, table bloat, `pg_stat_user_tables`

### Project: Slow Query Clinic
```sql
-- Setup: Load large dataset
INSERT INTO orders (customer_id, placed_at, total, status)
SELECT
    (random() * 10000)::INT,
    NOW() - (random() * 365 * 3)::INT * INTERVAL '1 day',
    (random() * 500)::NUMERIC(10,2),
    (ARRAY['pending','processing','shipped','delivered','cancelled'])[(random()*4)::INT+1]
FROM generate_series(1, 1000000);

-- Diagnose and fix each slow query:

-- 1. Slow: no index on status
SELECT * FROM orders WHERE status = 'pending';
-- Fix: CREATE INDEX ...

-- 2. Slow: function on indexed column
SELECT * FROM orders WHERE date_trunc('day', placed_at) = '2024-01-15';
-- Fix: range condition or expression index

-- 3. Slow: COUNT(*) on large table
SELECT COUNT(*) FROM orders;
-- Fix: pg_class.reltuples estimate

-- 4. N+1: correlated subquery per row
SELECT id, (SELECT COUNT(*) FROM order_items WHERE order_id = orders.id) FROM orders;
-- Fix: aggregating JOIN or window function

-- 5. Verify with EXPLAIN that your indexes are being used
```

### Milestone check:
- [ ] Can read EXPLAIN output and identify seq scans, hash joins, sort spills
- [ ] Can add the right index for any given slow query
- [ ] Understand when NOT to add an index (low-selectivity, write-heavy)

---

## Week 4: PL/pgSQL, Transactions, and Security

**Goal:** Write stored procedures, understand transactions, and secure a database.

### Topics:
- PL/pgSQL syntax: variables, IF/ELSE, loops, cursors, exceptions
- RETURNS TABLE, SETOF, polymorphic functions
- Stored procedures with COMMIT
- Triggers: BEFORE/AFTER, per-row/per-statement
- Transaction isolation levels and locking
- Roles, GRANT, REVOKE, Row-Level Security (RLS)

### Project: Banking System
```sql
-- Account system with proper transaction handling

CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL,
    balance NUMERIC(15,2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    from_account_id BIGINT REFERENCES accounts(id),
    to_account_id BIGINT REFERENCES accounts(id),
    amount NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Write these:
-- 1. Function: transfer_funds(from_id, to_id, amount) with proper locking
CREATE OR REPLACE FUNCTION transfer_funds(
    from_id BIGINT, to_id BIGINT, amount NUMERIC
) RETURNS VOID AS $$
DECLARE
    from_balance NUMERIC;
BEGIN
    -- Lock both accounts in consistent order to prevent deadlock
    PERFORM id FROM accounts
    WHERE id IN (from_id, to_id)
    ORDER BY id
    FOR UPDATE;

    SELECT balance INTO from_balance FROM accounts WHERE id = from_id;
    IF from_balance < amount THEN
        RAISE EXCEPTION 'Insufficient funds';
    END IF;

    UPDATE accounts SET balance = balance - amount WHERE id = from_id;
    UPDATE accounts SET balance = balance + amount WHERE id = to_id;
    INSERT INTO transactions (from_account_id, to_account_id, amount)
    VALUES (from_id, to_id, amount);
END;
$$ LANGUAGE plpgsql;

-- 2. Trigger: prevent balance going negative (redundant CHECK but shows triggers)
-- 3. RLS: owners can only see their own accounts
ALTER TABLE accounts ENABLE ROW LEVEL SECURITY;
CREATE POLICY accounts_owner ON accounts
    USING (owner_id = current_setting('app.current_user_id')::BIGINT);

-- 4. Audit log trigger on all account changes
```

### Milestone check:
- [ ] Can write PL/pgSQL functions with exception handling
- [ ] Understand difference between SERIALIZABLE and READ COMMITTED
- [ ] Can set up RLS for a multi-tenant table

---

## Week 5: Reliability — Backup, Replication, and Monitoring

**Goal:** Set up production-grade reliability mechanisms.

### Topics:
- `pg_dump` / `pg_restore` — all formats
- WAL archiving and Point-in-Time Recovery (PITR)
- pgBackRest setup and scheduling
- Streaming replication: primary → standby
- Logical replication for zero-downtime migrations
- `pg_stat_replication`, `pg_stat_archiver`

### Project: HA PostgreSQL Cluster
```bash
# Set up using Docker Compose:
# Primary: accepts reads and writes
# Replica: streams WAL, accepts reads

# docker-compose.yml:
# primary:
#   image: postgres:16
#   environment: POSTGRES_PASSWORD=secret
#   command: >
#     postgres
#     -c wal_level=replica
#     -c max_wal_senders=3
#     -c hot_standby=on
#
# replica:
#   image: postgres:16
#   depends_on: [primary]
#   command: >
#     bash -c "
#     pg_basebackup -h primary -U replicator -D $PGDATA -P -Xs -R
#     postgres
#     "

# Tasks:
# 1. Verify replication lag with pg_stat_replication
# 2. Simulate primary failure, promote replica
# 3. Take pg_dump backup and restore to a third instance
# 4. Configure WAL archiving to a local directory
# 5. Practice PITR: delete data, restore to before deletion
```

### Milestone check:
- [ ] Can restore a pg_dump backup to a fresh instance
- [ ] Can set up a streaming replica from scratch
- [ ] Understand RPO/RTO tradeoffs for different backup strategies

---

## Week 6: Production Systems — Docker, Kubernetes, and Tuning

**Goal:** Deploy PostgreSQL in containers, tune performance, and operate at scale.

### Topics:
- Docker + Docker Compose with health checks, init scripts
- Kubernetes StatefulSet, CloudNativePG operator
- PgBouncer setup in transaction mode
- postgresql.conf tuning (shared_buffers, work_mem, autovacuum)
- TimescaleDB for time-series data
- PostGIS for geospatial queries

### Project: Production-Ready Deployment
```yaml
# Full docker-compose.yml with:
# - postgres:16 with custom postgresql.conf
# - pgbouncer in transaction mode
# - pgadmin for management
# - promtail + prometheus + grafana for monitoring
# - Automated daily backups to a volume

# Tuning tasks:
# 1. Load 5M rows, measure query time before tuning
# 2. Tune shared_buffers, work_mem, random_page_cost
# 3. Verify cache hit ratio > 99%
# 4. Set up PgBouncer and measure connection overhead reduction
# 5. Create a Kubernetes StatefulSet with CloudNativePG
# 6. Configure HPA and PodDisruptionBudget
```

### Milestone check:
- [ ] Can deploy PostgreSQL in Docker/K8s with proper persistence
- [ ] Can tune postgresql.conf for a given workload
- [ ] Can use PgBouncer and measure its impact

---

## Capstone Project: Build a Complete Application Backend

Using all skills from weeks 1–6, build a multi-tenant SaaS backend database:

```sql
-- Requirements:
-- 1. Multi-tenant with Row-Level Security (one DB, many tenants)
-- 2. Audit logging on all critical tables via triggers
-- 3. Full-text search on articles/content
-- 4. JSONB for flexible metadata per tenant
-- 5. Partitioned orders table by month
-- 6. Materialized view for dashboard statistics (auto-refresh via pg_cron)
-- 7. Connection pooling via PgBouncer
-- 8. Automated backups with pgBackRest
-- 9. Read replica for analytics queries
-- 10. Monitoring with pg_stat_statements + Grafana
```

---

## Essential Resources

| Resource | Type | Focus |
|---|---|---|
| PostgreSQL Official Docs | Reference | Everything — the gold standard |
| "The Art of PostgreSQL" — Dimitri Fontaine | Book | SQL mastery, advanced features |
| "PostgreSQL: Up and Running" — Regina Obe | Book | Practical overview |
| explain.dalibo.com | Tool | Visual EXPLAIN ANALYZE |
| pganalyze.com | SaaS tool | Production query monitoring |
| Citus Data Blog | Blog | Scaling, advanced features |
| Supabase Blog | Blog | Practical PostgreSQL patterns |
| postgres.fm | Podcast | PostgreSQL news and deep dives |

---

## Skill Level Benchmarks

| Level | Can Do |
|---|---|
| **Beginner** | Write SELECT/INSERT/UPDATE, design simple schemas, use psql |
| **Intermediate** | Write CTEs, window functions, create indexes, read EXPLAIN |
| **Advanced** | PL/pgSQL, triggers, RLS, replication setup, performance tuning |
| **Expert** | Architecture decisions, custom extensions, disaster recovery, K8s operations |

---

## Common Certifications and Interviews

```
PostgreSQL does not have an official certification.
Practical skills that interviewers test:

1. Schema design: normalize a given scenario
2. Write a query with window functions (TOP N per group, running total)
3. EXPLAIN ANALYZE: "what's wrong with this plan?"
4. Index selection: "how would you optimize this query?"
5. Transactions: "explain MVCC and isolation levels"
6. Security: "how do you implement multi-tenant data isolation?"
7. Backup: "what's your backup strategy? What's your RPO/RTO?"
```

---

## Next Steps Beyond This Course

- **PostGIS**: geospatial data, map-based applications
- **TimescaleDB**: time-series data at scale (IoT, metrics)
- **Citus**: distributed PostgreSQL, horizontal sharding
- **Logical Replication**: zero-downtime major version upgrades
- **PL/Python, PL/Rust**: write stored functions in Python or Rust
- **pgvector**: vector similarity search for AI applications
- **Contribute to PostgreSQL**: the project is community-driven and welcoming
