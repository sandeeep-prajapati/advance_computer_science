# Chapter 12: Table Partitioning

## What is Partitioning?

Partitioning divides a large table into smaller physical pieces (partitions) while appearing as one logical table. Benefits:
- **Query performance** — queries scan only relevant partitions (partition pruning)
- **Maintenance** — VACUUM, index rebuilds operate on smaller tables
- **Data lifecycle** — drop old partitions instead of slow DELETE
- **Storage tiering** — move old partitions to slower/cheaper storage

---

## Partition Types

| Type | Split by | Example |
|---|---|---|
| **RANGE** | Continuous value ranges | Date ranges, ID ranges |
| **LIST** | Discrete values | Country, status, region |
| **HASH** | Hash function | Distribute evenly by key |

---

## Range Partitioning (Most Common — by Date)

```sql
-- Create partitioned parent table (no data stored here)
CREATE TABLE orders (
    id          BIGSERIAL,
    user_id     BIGINT NOT NULL,
    amount      NUMERIC(10, 2) NOT NULL,
    status      TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Create partitions (one per month)
CREATE TABLE orders_2024_01 PARTITION OF orders
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE orders_2024_02 PARTITION OF orders
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');

CREATE TABLE orders_2024_03 PARTITION OF orders
    FOR VALUES FROM ('2024-03-01') TO ('2024-04-01');

-- Default partition (catches values not matching any partition)
CREATE TABLE orders_default PARTITION OF orders DEFAULT;

-- Insert works normally
INSERT INTO orders (user_id, amount, status, created_at)
VALUES (1, 99.99, 'completed', '2024-02-15 10:00:00');
-- Goes into orders_2024_02 automatically

-- Query works on parent table
SELECT * FROM orders WHERE created_at >= '2024-02-01' AND created_at < '2024-03-01';
-- PostgreSQL only scans orders_2024_02 (partition pruning!)
```

---

## Primary Keys and Indexes

```sql
-- PRIMARY KEY must include the partition key
ALTER TABLE orders ADD PRIMARY KEY (id, created_at);

-- Or use a sequence separately
CREATE SEQUENCE orders_id_seq;
ALTER TABLE orders ALTER COLUMN id SET DEFAULT nextval('orders_id_seq');

-- Indexes on the parent are automatically created on all partitions
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status, created_at);
-- These create indexes on ALL partitions automatically

-- Or index individual partitions
CREATE INDEX idx_orders_2024_01_user ON orders_2024_01(user_id);
```

---

## List Partitioning (by Discrete Values)

```sql
-- Partition by region
CREATE TABLE sales (
    id         BIGSERIAL,
    region     TEXT NOT NULL,
    amount     NUMERIC(10, 2),
    created_at TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY LIST (region);

CREATE TABLE sales_us     PARTITION OF sales FOR VALUES IN ('us-east', 'us-west', 'us-central');
CREATE TABLE sales_europe PARTITION OF sales FOR VALUES IN ('de', 'fr', 'uk', 'es', 'it');
CREATE TABLE sales_apac   PARTITION OF sales FOR VALUES IN ('au', 'jp', 'sg', 'in');
CREATE TABLE sales_other  PARTITION OF sales DEFAULT;

-- Query only hits us partition
SELECT * FROM sales WHERE region IN ('us-east', 'us-west') AND amount > 1000;
```

---

## Hash Partitioning (Even Distribution)

```sql
-- Distribute users evenly across 4 partitions
CREATE TABLE users (
    id    BIGSERIAL,
    email TEXT NOT NULL,
    name  TEXT
) PARTITION BY HASH (id);

CREATE TABLE users_p0 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 0);
CREATE TABLE users_p1 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 1);
CREATE TABLE users_p2 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 2);
CREATE TABLE users_p3 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 3);
```

---

## Automating Partition Creation (pg_partman)

Manually creating monthly partitions is tedious. `pg_partman` extension automates it:

```sql
CREATE EXTENSION pg_partman;

-- Create the partitioned table
CREATE TABLE events (
    id BIGSERIAL,
    event_type TEXT NOT NULL,
    payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Let pg_partman manage it
SELECT partman.create_parent(
    p_parent_table => 'public.events',
    p_control      => 'created_at',
    p_interval     => 'monthly',    -- or 'daily', 'weekly', 'yearly'
    p_premake      => 3             -- pre-create 3 future partitions
);

-- Schedule maintenance (run from cron or pg_cron)
SELECT partman.run_maintenance();
-- Creates new partitions, optionally drops old ones
```

---

## Dropping Old Partitions

One of the biggest advantages of partitioning — dropping old data is instant:

```sql
-- View all partitions
SELECT
    child.relname AS partition_name,
    pg_size_pretty(pg_relation_size(child.oid)) AS size
FROM pg_inherits
JOIN pg_class parent ON pg_inherits.inhparent = parent.oid
JOIN pg_class child ON pg_inherits.inhrelid = child.oid
WHERE parent.relname = 'orders';

-- Detach partition first (makes it independent)
ALTER TABLE orders DETACH PARTITION orders_2022_01;

-- Then drop it
DROP TABLE orders_2022_01;

-- Or detach and archive
ALTER TABLE orders DETACH PARTITION orders_2022_01 CONCURRENTLY;
-- Now orders_2022_01 is a standalone table you can move to cold storage
```

---

## Partition Pruning

PostgreSQL automatically skips partitions that can't contain matching rows:

```sql
-- This only scans orders_2024_02
EXPLAIN SELECT * FROM orders
WHERE created_at >= '2024-02-01' AND created_at < '2024-03-01';
-- Output: Seq Scan on orders_2024_02

-- Enable runtime partition pruning (default: on in PostgreSQL 11+)
SET enable_partition_pruning = on;

-- Pruning also works with parameters
PREPARE order_query(TIMESTAMPTZ) AS
    SELECT * FROM orders WHERE created_at = $1;
-- PostgreSQL prunes at execution time based on the actual $1 value
```

---

## Sub-Partitioning

```sql
-- Partition by year, then by month within each year
CREATE TABLE events (
    id         BIGSERIAL,
    created_at TIMESTAMPTZ NOT NULL
) PARTITION BY RANGE (created_at);

CREATE TABLE events_2024 PARTITION OF events
    FOR VALUES FROM ('2024-01-01') TO ('2025-01-01')
    PARTITION BY RANGE (created_at);    -- sub-partition

CREATE TABLE events_2024_01 PARTITION OF events_2024
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE events_2024_02 PARTITION OF events_2024
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
```

---

## When to Partition

**Good candidates:**
- Tables > 100 GB
- Tables where queries always filter on the partition key (e.g., always by date)
- Tables with clear data lifecycle (old data archived or deleted regularly)
- Time-series data (logs, events, metrics)

**Not needed for:**
- Tables < 10 GB (overhead rarely worth it)
- Tables queried without the partition key (no pruning benefit)
- Tables with complex cross-partition JOINs (slower)

---

## Foreign Keys Limitation

PostgreSQL does NOT allow foreign keys to reference partitioned tables (only individual partitions). Design around this:

```sql
-- This fails:
-- FOREIGN KEY (order_id) REFERENCES orders(id)  -- orders is partitioned

-- Solutions:
-- 1. Reference the specific partition (not ideal)
-- 2. Use application-level integrity checks
-- 3. Use triggers to enforce integrity
-- 4. Denormalize the data (store user_id in the referencing table)
```
