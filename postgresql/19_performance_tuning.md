# Chapter 19: Performance Tuning

## Performance Tuning Stack

```
1. Schema design    → right data types, normalized schema, partitioning
2. Indexes          → B-tree, GIN, partial, covering
3. Query tuning     → EXPLAIN ANALYZE, rewrite slow queries
4. postgresql.conf  → memory, checkpoints, I/O
5. Hardware         → SSD, RAM, CPU cores
6. Connection layer → PgBouncer, connection limits
7. Autovacuum       → prevent table bloat
```

---

## postgresql.conf — Key Parameters

Tune these in `/etc/postgresql/16/main/postgresql.conf` or mount into Docker:

### Memory:
```ini
# shared_buffers: PostgreSQL's own cache (~25% of total RAM)
shared_buffers = 4GB          # for 16 GB server

# work_mem: per sort/hash operation — can multiply!
# (max_connections × parallel workers × work_mem = peak usage)
work_mem = 64MB               # conservative; increase for heavy analytics

# maintenance_work_mem: for VACUUM, CREATE INDEX, etc.
maintenance_work_mem = 512MB

# effective_cache_size: estimate of OS + PG cache (not actual allocation)
# Used by query planner to decide index vs seq scan
effective_cache_size = 12GB   # ~75% of total RAM
```

### Checkpoints:
```ini
# checkpoint_completion_target: spread checkpoint writes over this fraction of interval
checkpoint_completion_target = 0.9   # default 0.5 — higher reduces I/O spikes

# max_wal_size: allow WAL to grow this large between checkpoints
max_wal_size = 4GB            # increase if you see frequent checkpoints

# min_wal_size: keep at least this much WAL allocated
min_wal_size = 1GB
```

### Planner:
```ini
# random_page_cost: cost of non-sequential disk read
# Lower for SSD (default is 4.0 for HDD)
random_page_cost = 1.1        # SSD
# random_page_cost = 4.0      # HDD

# effective_io_concurrency: parallel I/O (SSD can handle many)
effective_io_concurrency = 200    # SSD
# effective_io_concurrency = 2    # HDD

# parallel query
max_parallel_workers_per_gather = 4   # use up to 4 workers for parallel seq scans
max_parallel_workers = 8
max_parallel_maintenance_workers = 4  # for CREATE INDEX, VACUUM
```

### Logging (for tuning):
```ini
# Log slow queries
log_min_duration_statement = 1000    # log queries taking >1 second (ms)
log_checkpoints = on                  # log checkpoint activity
log_lock_waits = on                   # log when waiting for locks
log_temp_files = 0                    # log all temp file creation (sorts that spill)
```

### Connections:
```ini
max_connections = 100         # keep low; use PgBouncer for high concurrency
```

Apply changes:
```sql
-- Reload config without restart (for most params)
SELECT pg_reload_conf();

-- Some params require restart (shared_buffers, max_connections)
-- check pg_settings
SELECT name, setting, pending_restart FROM pg_settings WHERE pending_restart;
```

---

## EXPLAIN ANALYZE Deep Dive

```sql
-- Full plan with timing and buffer stats
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)
SELECT u.name, COUNT(o.id)
FROM users u
JOIN orders o ON u.id = o.user_id
WHERE u.created_at > NOW() - INTERVAL '30 days'
GROUP BY u.id, u.name;
```

### Reading the output:
```
Finalize GroupAggregate  (cost=1234.56..5678.90 rows=100 width=40)
                          (actual time=45.678..50.123 rows=87 loops=1)
  Buffers: shared hit=890 read=234
  ->  Sort  ...
        Sort Method: quicksort  Memory: 256kB
        ->  Hash Join  (cost=...) (actual ... rows=15234 loops=1)
              Hash Cond: (o.user_id = u.id)
              Buffers: shared hit=450 read=234
              ->  Seq Scan on orders o
              ->  Hash  (cost=...) (actual rows=89 loops=1)
                    ->  Index Scan on users u
                          Index Cond: (created_at > ...)
                          Rows Removed by Filter: 12

Planning Time: 2.345 ms
Execution Time: 51.456 ms
```

Key things to look for:

| Signal | Meaning | Fix |
|---|---|---|
| `Seq Scan` on large table | No index used | Add appropriate index |
| `cost` rows >> `actual` rows | Bad row estimates | Run ANALYZE or increase statistics |
| `Rows Removed by Filter` >> rows returned | Ineffective index | Add more selective index |
| `Hash Batches: 8` | Hash join spilled to disk | Increase work_mem |
| `Sort Method: external` | Sort spilled to disk | Increase work_mem |
| High `shared read` | Cache miss | Add RAM / shared_buffers |
| `loops=10000` | Nested loop is problematic | Rewrite query or add index |

---

## Query Optimization Techniques

### Use covering indexes to avoid heap access:
```sql
-- Query: get user names and emails filtered by plan
SELECT name, email FROM users WHERE plan = 'pro';

-- Without covering index: Index Scan → heap fetch for name, email
CREATE INDEX idx_users_plan ON users(plan);

-- With covering index: Index Only Scan (no heap access)
CREATE INDEX idx_users_plan_covering ON users(plan) INCLUDE (name, email);
```

### Avoid functions on indexed columns:
```sql
-- BAD: function prevents index use
WHERE date_trunc('day', created_at) = '2024-01-15'
WHERE lower(email) = 'alice@example.com'

-- GOOD: use range condition
WHERE created_at >= '2024-01-15' AND created_at < '2024-01-16'

-- Or expression index
CREATE INDEX idx_users_email_lower ON users(lower(email));
WHERE lower(email) = 'alice@example.com'  -- now uses index
```

### Rewrite correlated subqueries as JOINs:
```sql
-- SLOW: correlated subquery (runs once per row)
SELECT name,
    (SELECT COUNT(*) FROM orders WHERE user_id = users.id) AS order_count
FROM users;

-- FAST: aggregating JOIN
SELECT u.name, COALESCE(o.order_count, 0)
FROM users u
LEFT JOIN (
    SELECT user_id, COUNT(*) AS order_count FROM orders GROUP BY user_id
) o ON u.id = o.user_id;
```

### Use EXISTS instead of IN for large sets:
```sql
-- Potentially slow with large subquery result
WHERE user_id IN (SELECT id FROM users WHERE plan = 'pro')

-- Usually faster
WHERE EXISTS (SELECT 1 FROM users WHERE id = orders.user_id AND plan = 'pro')
```

### Limit rows early:
```sql
-- Push filters into CTEs / subqueries
WITH recent_orders AS (
    SELECT * FROM orders
    WHERE created_at > NOW() - INTERVAL '7 days'  -- filter first
    LIMIT 1000
)
SELECT u.name, o.total
FROM recent_orders o
JOIN users u ON u.id = o.user_id;
```

---

## Statistics and the Planner

```sql
-- Check table statistics
SELECT schemaname, tablename, n_live_tup, n_dead_tup,
       last_analyze, last_autoanalyze
FROM pg_stat_user_tables
ORDER BY n_live_tup DESC;

-- Increase statistics for columns with skewed data
ALTER TABLE orders ALTER COLUMN status SET STATISTICS 500;  -- default 100
ANALYZE orders;

-- Check column statistics used by planner
SELECT * FROM pg_stats
WHERE tablename = 'orders' AND attname = 'status';

-- Check if estimates are off
EXPLAIN (ANALYZE, FORMAT JSON)
SELECT * FROM orders WHERE status = 'pending';
-- Compare rows estimate vs actual rows
```

---

## VACUUM and Table Bloat

PostgreSQL MVCC leaves dead tuples after UPDATE/DELETE. VACUUM reclaims them.

```sql
-- Check for table bloat
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS total_size,
    n_dead_tup,
    n_live_tup,
    round(n_dead_tup::numeric / NULLIF(n_live_tup + n_dead_tup, 0) * 100, 1) AS dead_pct,
    last_vacuum,
    last_autovacuum
FROM pg_stat_user_tables
ORDER BY n_dead_tup DESC;

-- Manual VACUUM (doesn't require table lock)
VACUUM users;

-- VACUUM ANALYZE: reclaim + update statistics
VACUUM ANALYZE users;

-- VACUUM FULL: rewrites table (reclaims disk, takes lock — use sparingly)
VACUUM FULL users;    -- locks table!
```

### Autovacuum tuning:
```ini
# postgresql.conf
autovacuum = on
autovacuum_max_workers = 5          # more workers for busy databases
autovacuum_vacuum_scale_factor = 0.05   # vacuum when 5% dead (default 0.2 = 20%)
autovacuum_analyze_scale_factor = 0.02  # analyze when 2% changed

# For high-write tables, per-table override:
```

```sql
-- Per-table autovacuum tuning (for high-churn tables)
ALTER TABLE orders SET (
    autovacuum_vacuum_scale_factor = 0.01,   -- 1% threshold
    autovacuum_vacuum_cost_delay = 2          -- less throttling
);
```

---

## Index Bloat

```sql
-- Find bloated indexes (uses pgstattuple if available)
SELECT
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY pg_relation_size(indexrelid) DESC;

-- Rebuild an index online (no lock)
REINDEX INDEX CONCURRENTLY idx_orders_user_id;

-- Or rebuild all indexes on a table
REINDEX TABLE CONCURRENTLY orders;
```

---

## pg_stat_* Views

```sql
-- Active connections and their queries
SELECT pid, state, wait_event_type, wait_event,
       now() - query_start AS duration,
       substring(query, 1, 80) AS query
FROM pg_stat_activity
WHERE state != 'idle'
ORDER BY duration DESC;

-- Lock waits
SELECT
    blocked.pid AS blocked_pid,
    blocked.query AS blocked_query,
    blocking.pid AS blocking_pid,
    blocking.query AS blocking_query,
    now() - blocked.query_start AS wait_time
FROM pg_stat_activity blocked
JOIN pg_stat_activity blocking
    ON blocking.pid = ANY(pg_blocking_pids(blocked.pid))
WHERE blocked.cardinality(pg_blocking_pids(blocked.pid)) > 0;

-- Table access patterns (sequential vs index scans)
SELECT
    relname,
    seq_scan,
    idx_scan,
    seq_tup_read,
    idx_tup_fetch,
    n_tup_ins + n_tup_upd + n_tup_del AS write_ops
FROM pg_stat_user_tables
ORDER BY seq_scan DESC;

-- Unused indexes (waste of disk and write overhead)
SELECT
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE idx_scan = 0
  AND NOT EXISTS (
      SELECT 1 FROM pg_constraint c
      WHERE c.conindid = indexrelid
  )
ORDER BY pg_relation_size(indexrelid) DESC;

-- Cache hit ratio (should be >99%)
SELECT
    sum(heap_blks_read) AS heap_read,
    sum(heap_blks_hit) AS heap_hit,
    round(sum(heap_blks_hit)::numeric / NULLIF(sum(heap_blks_hit + heap_blks_read), 0) * 100, 2) AS cache_hit_ratio
FROM pg_statio_user_tables;
```

---

## Slow Query Identification

```sql
-- Requires pg_stat_statements extension
CREATE EXTENSION pg_stat_statements;

-- Top 10 slowest by total time
SELECT
    round(total_exec_time::numeric, 2) AS total_ms,
    calls,
    round(mean_exec_time::numeric, 2) AS avg_ms,
    round(stddev_exec_time::numeric, 2) AS stddev_ms,
    rows / NULLIF(calls, 0) AS rows_per_call,
    substring(query, 1, 100) AS query
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 10;

-- Top 10 by I/O
SELECT
    shared_blks_read + shared_blks_written AS total_io,
    calls,
    round(mean_exec_time::numeric, 2) AS avg_ms,
    substring(query, 1, 100) AS query
FROM pg_stat_statements
ORDER BY total_io DESC
LIMIT 10;

-- Reset stats
SELECT pg_stat_statements_reset();
```

---

## Parallel Query

```sql
-- Enable parallel queries (usually on by default)
SET max_parallel_workers_per_gather = 4;

-- Force parallel plan for testing
SET parallel_setup_cost = 0;
SET parallel_tuple_cost = 0;
SET min_parallel_table_scan_size = 0;

-- Check if parallel is being used
EXPLAIN SELECT avg(price) FROM orders WHERE status = 'completed';
-- Look for: Gather or Gather Merge node at the top

-- Disable parallel for a session if causing issues
SET max_parallel_workers_per_gather = 0;
```

---

## Partitioning for Performance

Already covered in Chapter 12, but key performance note:

```sql
-- Partition pruning: PostgreSQL only scans relevant partition
EXPLAIN SELECT * FROM orders
WHERE created_at >= '2024-01-01' AND created_at < '2024-02-01';

-- Should show: Seq Scan on orders_2024_01 (just one partition)
-- NOT: Seq Scan on orders (entire table)
```

---

## Performance Tuning Checklist

```
Schema:
  ✓ Appropriate data types (int4 not int8 where possible)
  ✓ NOT NULL wherever applicable (enables optimizations)
  ✓ Partition large tables (>100M rows)

Indexes:
  ✓ FK columns indexed
  ✓ WHERE clause columns indexed
  ✓ Covering indexes for read-heavy queries
  ✓ Removed unused indexes

Memory:
  ✓ shared_buffers = 25% RAM
  ✓ effective_cache_size = 75% RAM
  ✓ work_mem tuned per query type
  ✓ Cache hit ratio > 99%

Maintenance:
  ✓ Autovacuum enabled and healthy
  ✓ No tables with >20% dead tuples
  ✓ Table/index bloat monitored

Queries:
  ✓ Slow query log enabled
  ✓ pg_stat_statements installed
  ✓ Top 10 slowest queries reviewed and optimized

Connections:
  ✓ PgBouncer in front (transaction mode)
  ✓ max_connections set conservatively

Hardware:
  ✓ random_page_cost = 1.1 for SSD
  ✓ effective_io_concurrency = 200 for SSD
```
