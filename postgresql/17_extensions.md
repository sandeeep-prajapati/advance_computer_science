# Chapter 17: Extensions

## What are Extensions?

PostgreSQL extensions add new functionality to the database: data types, functions, operators, and more. There are hundreds of extensions — both built-in and community-contributed.

```sql
-- List available extensions
SELECT name, default_version, comment FROM pg_available_extensions ORDER BY name;

-- List installed extensions
SELECT extname, extversion FROM pg_extension;

-- Install an extension (requires superuser or appropriate permissions)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Remove extension
DROP EXTENSION IF EXISTS "uuid-ossp";
```

---

## pg_stat_statements — Query Analytics

Tracks execution statistics for all SQL statements. Essential for performance tuning.

```sql
-- Enable in postgresql.conf
-- shared_preload_libraries = 'pg_stat_statements'

CREATE EXTENSION pg_stat_statements;

-- Top 10 slowest queries
SELECT
    round(mean_exec_time::numeric, 2) AS avg_ms,
    calls,
    round(total_exec_time::numeric, 2) AS total_ms,
    rows,
    substring(query, 1, 60) AS short_query
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

-- Most called queries
SELECT calls, round(mean_exec_time::numeric, 2) AS avg_ms, query
FROM pg_stat_statements
ORDER BY calls DESC
LIMIT 10;

-- Queries with most I/O
SELECT query, shared_blks_read + shared_blks_written AS total_io
FROM pg_stat_statements
ORDER BY total_io DESC
LIMIT 10;

-- Reset statistics
SELECT pg_stat_statements_reset();
```

---

## uuid-ossp — UUID Generation

```sql
CREATE EXTENSION "uuid-ossp";

-- Generate v4 UUID (random)
SELECT uuid_generate_v4();    -- '550e8400-e29b-41d4-a716-446655440000'

-- Use as default ID
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Or use built-in (PostgreSQL 13+, no extension needed)
SELECT gen_random_uuid();
```

---

## pgcrypto — Cryptographic Functions

```sql
CREATE EXTENSION pgcrypto;

-- Hash password (bcrypt)
SELECT crypt('mypassword', gen_salt('bf', 12));

-- Verify password
SELECT crypt('mypassword', stored_hash) = stored_hash AS is_valid;

-- Generate random bytes
SELECT encode(gen_random_bytes(32), 'hex');    -- 64-char hex string
SELECT encode(gen_random_bytes(16), 'base64'); -- base64 string

-- Encrypt/decrypt data
SELECT pgp_sym_encrypt('secret data', 'encryption_key');
SELECT pgp_sym_decrypt(encrypted_data, 'encryption_key');

-- SHA hashes
SELECT encode(digest('hello', 'sha256'), 'hex');
SELECT encode(digest('hello', 'sha512'), 'hex');
SELECT encode(hmac('message', 'key', 'sha256'), 'hex');
```

---

## PostGIS — Geospatial Data

PostGIS adds geographic objects, allowing PostgreSQL to be used as a spatial database.

```sql
CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;

-- Store geographic coordinates
CREATE TABLE locations (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    coordinates GEOGRAPHY(POINT, 4326)  -- 4326 = WGS84 (GPS coordinates)
);

-- Insert a point (longitude, latitude)
INSERT INTO locations (name, coordinates)
VALUES
    ('Eiffel Tower', ST_MakePoint(2.2945, 48.8584)),
    ('Statue of Liberty', ST_MakePoint(-74.0445, 40.6892));

-- Find locations within 10km of a point
SELECT name,
    ST_Distance(coordinates, ST_MakePoint(2.3522, 48.8566)::GEOGRAPHY) / 1000 AS km
FROM locations
WHERE ST_DWithin(
    coordinates,
    ST_MakePoint(2.3522, 48.8566)::GEOGRAPHY,
    10000   -- meters
)
ORDER BY km;

-- GiST index for spatial queries
CREATE INDEX idx_locations_coords ON locations USING GIST(coordinates);

-- Distance between two points
SELECT ST_Distance(
    ST_MakePoint(2.2945, 48.8584)::GEOGRAPHY,
    ST_MakePoint(-74.0445, 40.6892)::GEOGRAPHY
) / 1000 AS km;
-- ~5837 km

-- Area of a polygon
SELECT ST_Area(ST_GeomFromText('POLYGON((0 0, 1 0, 1 1, 0 1, 0 0))'));

-- Buffer (area around a point)
SELECT ST_Buffer(ST_MakePoint(2.2945, 48.8584)::GEOGRAPHY, 500);  -- 500m radius
```

---

## pg_cron — Scheduled Jobs Inside PostgreSQL

```sql
-- shared_preload_libraries = 'pg_cron' (in postgresql.conf)
CREATE EXTENSION pg_cron;

-- Schedule a job: refresh materialized view every hour
SELECT cron.schedule('refresh-stats', '0 * * * *',
    'REFRESH MATERIALIZED VIEW CONCURRENTLY product_stats');

-- Run vacuum every night at 2 AM
SELECT cron.schedule('nightly-vacuum', '0 2 * * *',
    'VACUUM ANALYZE');

-- Delete old sessions every 5 minutes
SELECT cron.schedule('cleanup-sessions', '*/5 * * * *',
    'DELETE FROM sessions WHERE expires_at < NOW()');

-- List scheduled jobs
SELECT jobid, schedule, command, active FROM cron.job;

-- Unschedule a job
SELECT cron.unschedule(jobid) FROM cron.job WHERE jobname = 'nightly-vacuum';

-- View job run history
SELECT * FROM cron.job_run_details ORDER BY start_time DESC LIMIT 20;
```

---

## TimescaleDB — Time-Series Data

TimescaleDB makes PostgreSQL excel at time-series workloads (IoT, metrics, logs):

```sql
CREATE EXTENSION timescaledb;

-- Create a hypertable (partitioned automatically by time)
CREATE TABLE metrics (
    time TIMESTAMPTZ NOT NULL,
    device_id INTEGER,
    cpu_usage DOUBLE PRECISION,
    memory_usage DOUBLE PRECISION
);

SELECT create_hypertable('metrics', 'time');

-- Insert works normally
INSERT INTO metrics VALUES (NOW(), 1, 72.5, 45.2);

-- Time-series queries are automatically optimized
SELECT time_bucket('1 hour', time) AS hour,
       device_id,
       AVG(cpu_usage) AS avg_cpu
FROM metrics
WHERE time > NOW() - INTERVAL '24 hours'
GROUP BY hour, device_id
ORDER BY hour;

-- Continuous aggregates (like materialized views, auto-refreshed)
CREATE MATERIALIZED VIEW metrics_hourly
WITH (timescaledb.continuous) AS
SELECT time_bucket('1 hour', time) AS hour,
       device_id,
       AVG(cpu_usage) AS avg_cpu,
       MAX(cpu_usage) AS max_cpu
FROM metrics
GROUP BY hour, device_id;

-- Automatic data retention (drop data older than 90 days)
SELECT add_retention_policy('metrics', INTERVAL '90 days');
```

---

## hstore — Key-Value Pairs

```sql
CREATE EXTENSION hstore;

CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    attributes HSTORE
);

INSERT INTO products (name, attributes)
VALUES ('Laptop', 'brand=>Dell, cpu=>i7, ram=>16GB, color=>black');

-- Query
SELECT attributes->'brand' FROM products WHERE id = 1;       -- 'Dell'
SELECT * FROM products WHERE attributes->'brand' = 'Dell';
SELECT * FROM products WHERE attributes ? 'cpu';              -- key exists

-- Update
UPDATE products SET attributes = attributes || 'ssd=>512GB'::hstore WHERE id = 1;

-- GIN index
CREATE INDEX idx_products_attrs ON products USING GIN(attributes);
```

Note: JSONB is generally preferred over hstore now — more flexible and feature-rich.

---

## pg_trgm — Fuzzy String Matching

```sql
CREATE EXTENSION pg_trgm;

-- Similarity between strings (0 to 1)
SELECT similarity('postgresql', 'postgress');  -- ~0.7

-- Find similar strings
SELECT name FROM products WHERE name % 'laptpp';  -- typo-tolerant search
-- % operator: returns true if similarity > pg_trgm.similarity_threshold (default 0.3)

-- GIN or GIST index for fast similarity search
CREATE INDEX idx_products_name_trgm ON products USING GIN(name gin_trgm_ops);

-- LIKE with trigram index (normally LIKE doesn't use indexes for %word%)
SELECT name FROM products WHERE name ILIKE '%dell laptop%';
-- With GIN trgm index, this is fast!

-- Combined FTS + fuzzy matching
SELECT name,
    similarity(name, 'laptpp') AS score
FROM products
WHERE name % 'laptpp'
ORDER BY score DESC;
```

---

## Extension Quick Reference

| Extension | Purpose | Install |
|---|---|---|
| `pg_stat_statements` | Query analytics | Requires `shared_preload_libraries` |
| `uuid-ossp` | UUID generation | `CREATE EXTENSION "uuid-ossp"` |
| `pgcrypto` | Encryption, hashing | `CREATE EXTENSION pgcrypto` |
| `postgis` | Geospatial | Separate install + `CREATE EXTENSION` |
| `pg_cron` | Scheduled jobs | Requires `shared_preload_libraries` |
| `timescaledb` | Time-series | Separate install + `shared_preload_libraries` |
| `hstore` | Key-value pairs | `CREATE EXTENSION hstore` |
| `pg_trgm` | Fuzzy search | `CREATE EXTENSION pg_trgm` |
| `btree_gin` | GIN index for btree ops | `CREATE EXTENSION btree_gin` |
| `intarray` | Integer array operators | `CREATE EXTENSION intarray` |
| `earthdistance` | Earth-distance calculations | `CREATE EXTENSION earthdistance` |
