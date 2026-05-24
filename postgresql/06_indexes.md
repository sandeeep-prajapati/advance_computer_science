# Chapter 6: Indexes & Performance

## What is an Index?

An index is a data structure that allows PostgreSQL to find rows quickly without scanning the entire table. Without an index, every query scans every row (sequential scan). With an index, the database jumps directly to matching rows.

```
Without index: Scan all 10,000,000 rows to find email = 'alice@example.com'
With index:    Jump to the exact row in ~3 operations (B-tree depth)
```

---

## Index Types

### B-tree (Default — for most use cases)
```sql
-- Created automatically for PRIMARY KEY and UNIQUE
-- Supports: =, <, >, <=, >=, BETWEEN, IN, LIKE 'prefix%', IS NULL

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_status ON orders(status);
```

### Hash (for equality only)
```sql
-- Faster than B-tree for equality checks, but no range queries
CREATE INDEX idx_sessions_token ON sessions USING HASH (token);
```

### GIN (Generalized Inverted Index — for full-text, arrays, JSONB)
```sql
-- Full-text search
CREATE INDEX idx_posts_search ON posts USING GIN (to_tsvector('english', title || ' ' || body));

-- JSONB
CREATE INDEX idx_products_attrs ON products USING GIN (attributes);

-- Arrays
CREATE INDEX idx_posts_tags ON posts USING GIN (tags);
```

### GiST (Generalized Search Tree — for geometric, range, full-text)
```sql
-- Range types
CREATE INDEX idx_reservations_period ON reservations USING GIST (during);

-- PostGIS geometry
CREATE INDEX idx_locations_geom ON locations USING GIST (coordinates);

-- Full-text (GiST supports phrase search, GIN is faster for most FTS)
CREATE INDEX idx_docs_search ON documents USING GIST (to_tsvector('english', content));
```

### BRIN (Block Range INdex — for large ordered tables)
```sql
-- Very small size, good for time-series data in insertion order
-- Works on physical page ranges — tables must be physically sorted
CREATE INDEX idx_events_time ON events USING BRIN (created_at);
```

---

## Composite Indexes

```sql
-- Covers queries filtering on both columns
CREATE INDEX idx_orders_user_status ON orders(user_id, status);

-- Covers: WHERE user_id = 1 AND status = 'pending'
-- Covers: WHERE user_id = 1 (leftmost prefix)
-- Does NOT cover: WHERE status = 'pending' alone (not leftmost)

-- For ORDER BY + WHERE
CREATE INDEX idx_orders_user_created ON orders(user_id, created_at DESC);
-- Covers: WHERE user_id = 1 ORDER BY created_at DESC LIMIT 10
```

**Rule:** Order columns from most selective to least in composite index. Put equality columns before range columns.

---

## Partial Indexes

Index only a subset of rows — smaller, faster:

```sql
-- Only index active users (don't waste space on deleted ones)
CREATE INDEX idx_active_users_email ON users(email) WHERE is_active = TRUE;

-- Only index pending orders (completed orders rarely queried by status)
CREATE INDEX idx_pending_orders ON orders(user_id, created_at)
WHERE status = 'pending';

-- Only index unread notifications
CREATE INDEX idx_unread_notifications ON notifications(user_id, created_at)
WHERE read_at IS NULL;
```

---

## Expression Indexes (Functional Indexes)

```sql
-- Case-insensitive email lookup
CREATE INDEX idx_users_email_lower ON users(LOWER(email));
-- Query must also use LOWER: WHERE LOWER(email) = 'alice@example.com'

-- Index on JSON field
CREATE INDEX idx_users_prefs_theme ON users((preferences->>'theme'));

-- Index on computed expression
CREATE INDEX idx_orders_year ON orders(EXTRACT(YEAR FROM created_at));
```

---

## Covering Indexes (INCLUDE)

Include extra columns in the index so queries can be satisfied from the index alone (index-only scan):

```sql
-- Query: SELECT email, name FROM users WHERE id = 42
-- Without INCLUDE: index lookup for id, then heap fetch for email, name
-- With INCLUDE: everything in the index

CREATE INDEX idx_users_id_cover ON users(id) INCLUDE (email, name);
CREATE INDEX idx_orders_user_cover ON orders(user_id, status) INCLUDE (amount, created_at);
```

---

## Unique Indexes

```sql
-- Enforces uniqueness (also creates a constraint)
CREATE UNIQUE INDEX idx_users_email ON users(email);

-- Partial unique index (only one active record per user)
CREATE UNIQUE INDEX idx_one_active_subscription ON subscriptions(user_id)
WHERE status = 'active';

-- Unique on expression
CREATE UNIQUE INDEX idx_users_email_unique ON users(LOWER(email));
```

---

## EXPLAIN — Reading Query Plans

`EXPLAIN` shows the query plan without executing. `EXPLAIN ANALYZE` executes and shows actual timings.

```sql
EXPLAIN SELECT * FROM users WHERE email = 'alice@example.com';
-- If index exists: shows Index Scan
-- If no index: shows Seq Scan

EXPLAIN ANALYZE SELECT * FROM orders WHERE user_id = 1 ORDER BY created_at DESC LIMIT 10;
```

### Reading EXPLAIN Output:
```
Index Scan using idx_orders_user_created on orders  (cost=0.43..8.90 rows=10 width=96) (actual time=0.078..0.099 rows=10 loops=1)
  Index Cond: (user_id = 1)
Planning Time: 0.234 ms
Execution Time: 0.127 ms
```

| Term | Meaning |
|---|---|
| `Seq Scan` | Full table scan (usually bad for large tables) |
| `Index Scan` | Uses index, fetches rows from heap |
| `Index Only Scan` | Uses index only (fastest — all data in index) |
| `Bitmap Index Scan` | Uses index to build bitmap, then heap fetch |
| `cost=X..Y` | Estimated startup cost .. total cost |
| `rows=N` | Estimated row count |
| `actual time=X..Y` | Real timing in milliseconds |
| `loops=N` | How many times this node executed |

### EXPLAIN options:
```sql
EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) SELECT ...;
-- ANALYZE: execute and show actual times
-- BUFFERS: show buffer cache hits/misses
-- FORMAT JSON: machine-readable output (use with explain.dalibo.com)
```

---

## Finding Slow Queries

```sql
-- Enable pg_stat_statements extension
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Top 10 slowest queries
SELECT
    round(total_exec_time::numeric, 2) AS total_ms,
    calls,
    round(mean_exec_time::numeric, 2) AS avg_ms,
    round((100 * total_exec_time / sum(total_exec_time) OVER ())::numeric, 2) AS pct,
    substring(query, 1, 80) AS query
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 10;

-- Queries with sequential scans on large tables
SELECT schemaname, tablename, seq_scan, seq_tup_read, idx_scan
FROM pg_stat_user_tables
WHERE seq_scan > 0
ORDER BY seq_tup_read DESC;
```

---

## Index Management

```sql
-- List all indexes
\di
SELECT tablename, indexname, indexdef FROM pg_indexes WHERE schemaname = 'public';

-- Index sizes
SELECT
    indexname,
    pg_size_pretty(pg_relation_size(indexname::regclass)) AS size
FROM pg_indexes
WHERE tablename = 'users';

-- Index usage stats
SELECT
    schemaname, tablename, indexname,
    idx_scan AS times_used,
    idx_tup_read, idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;

-- Find unused indexes (idx_scan = 0 after reasonable uptime)
SELECT indexname, tablename, idx_scan
FROM pg_stat_user_indexes
WHERE idx_scan = 0;

-- Drop index
DROP INDEX idx_users_email;
DROP INDEX CONCURRENTLY idx_users_email;   -- non-blocking (preferred in production)

-- Rebuild index
REINDEX INDEX idx_users_email;
REINDEX TABLE users;
REINDEX CONCURRENTLY INDEX idx_users_email;   -- non-blocking (PostgreSQL 12+)

-- Create index without blocking writes (takes longer but safe in production)
CREATE INDEX CONCURRENTLY idx_users_name ON users(name);
```

---

## When to Add an Index

**Add an index when:**
- Column appears in WHERE, JOIN ON, or ORDER BY frequently
- Table has more than ~1,000 rows
- The query is slow and EXPLAIN shows a Seq Scan

**Don't add an index when:**
- Table is small (Seq Scan is faster for small tables)
- Column has very low cardinality (e.g., boolean with 50/50 split)
- Table gets far more writes than reads (indexes slow down writes)
- You already have a covering composite index that includes this column

**Index Overhead:**
- Each index adds write overhead (UPDATE, INSERT, DELETE must update all indexes)
- Indexes use disk space
- Too many indexes confuse the planner
- Rule of thumb: start with no indexes, add based on query analysis
