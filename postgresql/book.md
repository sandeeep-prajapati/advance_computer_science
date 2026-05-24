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
# Chapter 2: Setup & Installation

## Installation Options

| Method | Best For |
|---|---|
| Package manager (apt/brew) | Local development on Linux/Mac |
| Docker | Isolated, reproducible, easy teardown |
| Installer | Windows |
| Cloud (RDS, Supabase, Neon) | Production, managed service |

---

## Linux (Ubuntu/Debian)

```bash
# Add official PostgreSQL repository
sudo apt install -y curl ca-certificates
sudo install -d /usr/share/postgresql-common/pgdg
sudo curl -o /usr/share/postgresql-common/pgdg/apt.postgresql.org.asc \
  --fail https://www.postgresql.org/media/keys/ACCC4CF8.asc
sudo sh -c 'echo "deb [signed-by=/usr/share/postgresql-common/pgdg/apt.postgresql.org.asc] \
  https://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > \
  /etc/apt/sources.list.d/pgdg.list'

# Install
sudo apt update
sudo apt install -y postgresql-16

# Start service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Check status
sudo systemctl status postgresql
```

---

## macOS (Homebrew)

```bash
brew install postgresql@16

# Start service
brew services start postgresql@16

# Add to PATH
echo 'export PATH="/opt/homebrew/opt/postgresql@16/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Verify
psql --version
```

Or use **Postgres.app** — a macOS menu bar app with one-click start/stop.

---

## Windows

Download the installer from `postgresql.org/download/windows` and run it. The installer includes:
- PostgreSQL server
- pgAdmin 4 (GUI)
- psql (CLI)
- Stack Builder (additional tools)

---

## Docker (Recommended for Development)

```bash
# Quick start — single container
docker run -d \
  --name postgres \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypassword \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql/data \
  postgres:16-alpine

# Connect
docker exec -it postgres psql -U myuser -d mydb
```

### Docker Compose (better for teams):
```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: myuser
      POSTGRES_PASSWORD: mypassword
      POSTGRES_DB: mydb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql   # auto-runs on first start
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U myuser -d mydb"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

```bash
docker-compose up -d
```

---

## psql — The Command-Line Client

`psql` is the primary CLI for PostgreSQL. You'll use it constantly.

### Connecting:
```bash
# Local connection
psql -U myuser -d mydb

# Remote connection
psql -h hostname -p 5432 -U myuser -d mydb

# From connection string
psql "postgresql://myuser:mypassword@localhost:5432/mydb"

# As postgres superuser (Linux)
sudo -u postgres psql
```

### Essential psql Commands:

```sql
-- Meta-commands (start with \)
\l                    -- list all databases
\c mydb               -- connect to database
\dt                   -- list tables in current schema
\dt *.*               -- list tables in all schemas
\d tablename          -- describe table (columns, types, constraints)
\di                   -- list indexes
\df                   -- list functions
\dv                   -- list views
\dn                   -- list schemas
\du                   -- list roles/users
\x                    -- toggle expanded display (useful for wide tables)
\timing               -- toggle query execution time display
\e                    -- open query in editor
\i file.sql           -- execute SQL from file
\o file.txt           -- send query output to file
\copy ...             -- client-side COPY (works without superuser)
\q                    -- quit
\?                    -- help on meta-commands
\h SELECT             -- help on SQL command
```

### psql Tips:
```bash
# Set prompt to show database and user
export PSQL_PROMPT='%n@%m:%>/%~%# '

# Save history
export PSQLHISTFILE=~/.psql_history

# .psqlrc configuration file
echo "\timing on" >> ~/.psqlrc
echo "\x auto" >> ~/.psqlrc
```

---

## pgAdmin 4 — GUI Client

pgAdmin 4 is the official web-based GUI for PostgreSQL.

```bash
# Install on Ubuntu
sudo apt install pgadmin4

# Or run via Docker
docker run -p 80:80 \
  -e 'PGADMIN_DEFAULT_EMAIL=admin@example.com' \
  -e 'PGADMIN_DEFAULT_PASSWORD=admin' \
  dpage/pgadmin4
```

Open `http://localhost` and register your PostgreSQL server.

---

## Connection String Format

```
postgresql://[user[:password]@][host][:port][/dbname][?param1=value1&...]

Examples:
  postgresql://localhost/mydb
  postgresql://myuser:mypass@localhost:5432/mydb
  postgresql://myuser:mypass@db.example.com/prod?sslmode=require
  postgresql:///mydb?host=/var/run/postgresql   # Unix socket
```

---

## First Steps After Installation

```sql
-- Connect as superuser
sudo -u postgres psql

-- Create a role (user)
CREATE USER myuser WITH PASSWORD 'mypassword';

-- Create a database owned by that user
CREATE DATABASE mydb OWNER myuser;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE mydb TO myuser;

-- Connect to new database
\c mydb

-- Create your first table
CREATE TABLE users (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Insert data
INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com');

-- Query
SELECT * FROM users;
```

---

## pg_hba.conf — Authentication Config

PostgreSQL's host-based authentication is configured in `pg_hba.conf`:

```
# TYPE  DATABASE  USER  ADDRESS        METHOD
local   all       all                  trust          # local Unix socket
host    all       all   127.0.0.1/32   scram-sha-256  # TCP from localhost
host    mydb      myuser 10.0.0.0/8   scram-sha-256  # TCP from private network
hostssl all       all   0.0.0.0/0      scram-sha-256  # SSL only from anywhere
```

After editing, reload: `sudo systemctl reload postgresql`

---

## Key Config File Locations

| File | Purpose |
|---|---|
| `postgresql.conf` | Server configuration (memory, connections, logging) |
| `pg_hba.conf` | Client authentication rules |
| `pg_ident.conf` | OS user → PostgreSQL user mapping |
| Data directory | `/var/lib/postgresql/16/main/` (Linux) |

```bash
# Find config file location
sudo -u postgres psql -c "SHOW config_file;"
sudo -u postgres psql -c "SHOW data_directory;"
```
# Chapter 3: Databases, Schemas & Tables

## Databases

A PostgreSQL **server** (cluster) can host multiple databases. Each database is isolated — you cannot query across databases in a single SQL statement.

```sql
-- List all databases
\l
SELECT datname FROM pg_database;

-- Create database
CREATE DATABASE myapp;
CREATE DATABASE myapp
    WITH OWNER = myuser
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TEMPLATE = template0;

-- Drop database (must disconnect first)
DROP DATABASE myapp;

-- Connect to database
\c myapp
```

---

## Schemas

A **schema** is a namespace within a database. The default schema is `public`. Use schemas to organize objects and control access.

```sql
-- List schemas
\dn
SELECT schema_name FROM information_schema.schemata;

-- Create schemas
CREATE SCHEMA app;
CREATE SCHEMA audit;
CREATE SCHEMA reporting;

-- Set search path (which schemas to look in, in order)
SET search_path TO app, public;

-- Permanent search path for a user
ALTER ROLE myuser SET search_path TO app, public;

-- Create a table in a specific schema
CREATE TABLE audit.activity_log (
    id BIGSERIAL PRIMARY KEY,
    action TEXT,
    occurred_at TIMESTAMPTZ DEFAULT NOW()
);

-- Reference cross-schema
SELECT * FROM audit.activity_log;
```

---

## Data Types

PostgreSQL has the richest set of data types of any RDBMS.

### Numeric:
```sql
SMALLINT          -- 2 bytes, -32768 to +32767
INTEGER / INT     -- 4 bytes, ~-2 billion to +2 billion
BIGINT            -- 8 bytes, ~-9.2 quintillion to +9.2 quintillion
SERIAL            -- auto-incrementing INTEGER (shorthand)
BIGSERIAL         -- auto-incrementing BIGINT
DECIMAL(p, s)     -- exact numeric, p total digits, s after decimal
NUMERIC(p, s)     -- same as DECIMAL
REAL              -- 4-byte floating point
DOUBLE PRECISION  -- 8-byte floating point
```

### Text:
```sql
CHAR(n)           -- fixed-length, padded with spaces
VARCHAR(n)        -- variable-length, max n characters
TEXT              -- unlimited variable-length (preferred in PostgreSQL)
```

### Date & Time:
```sql
DATE              -- date only (2024-01-15)
TIME              -- time only (14:30:00)
TIMESTAMP         -- date + time, no timezone
TIMESTAMPTZ       -- date + time WITH timezone (store in UTC) ← USE THIS
INTERVAL          -- time duration ('1 year 2 months 3 days')
```

### Boolean:
```sql
BOOLEAN           -- true, false, null
                  -- accepts: TRUE/FALSE, 't'/'f', 'yes'/'no', '1'/'0'
```

### JSON:
```sql
JSON              -- stores text, validates on input
JSONB             -- binary storage, indexed, faster queries ← PREFER THIS
```

### Arrays:
```sql
INTEGER[]         -- array of integers
TEXT[]            -- array of text
JSONB[]           -- array of JSONB
```

### UUID:
```sql
UUID              -- universally unique identifier
                  -- '550e8400-e29b-41d4-a716-446655440000'
```

### Other Useful Types:
```sql
BYTEA             -- binary data
INET              -- IPv4 or IPv6 address
CIDR              -- IP network
MACADDR           -- MAC address
POINT             -- geometric point (x, y)
TSRANGE           -- timestamp range
INT4RANGE         -- integer range
MONEY             -- currency amount (avoid — locale-dependent)
```

---

## Creating Tables

```sql
CREATE TABLE users (
    -- Identity
    id          BIGSERIAL PRIMARY KEY,

    -- Text
    name        VARCHAR(100)    NOT NULL,
    email       TEXT            NOT NULL,
    bio         TEXT,                       -- nullable

    -- Numbers
    age         SMALLINT        CHECK (age >= 0 AND age <= 150),
    balance     NUMERIC(12, 2)  DEFAULT 0.00,

    -- Boolean
    is_active   BOOLEAN         NOT NULL DEFAULT TRUE,

    -- Enum via CHECK or domain
    role        TEXT            NOT NULL DEFAULT 'user'
                                CHECK (role IN ('admin', 'user', 'guest')),

    -- Timestamps (always use TIMESTAMPTZ)
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT users_email_unique UNIQUE (email)
);

-- Table with foreign key
CREATE TABLE posts (
    id          BIGSERIAL       PRIMARY KEY,
    user_id     BIGINT          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       VARCHAR(255)    NOT NULL,
    body        TEXT            NOT NULL,
    status      TEXT            NOT NULL DEFAULT 'draft'
                                CHECK (status IN ('draft', 'published', 'archived')),
    published_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
```

---

## Constraints

```sql
-- NOT NULL
ALTER TABLE users ALTER COLUMN name SET NOT NULL;

-- UNIQUE (single column)
ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);

-- UNIQUE (composite)
ALTER TABLE team_members ADD CONSTRAINT one_role_per_user
    UNIQUE (team_id, user_id);

-- CHECK
ALTER TABLE users ADD CONSTRAINT valid_age CHECK (age BETWEEN 0 AND 150);

-- FOREIGN KEY with actions
ALTER TABLE posts ADD CONSTRAINT posts_user_fk
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE       -- delete posts when user deleted
    ON UPDATE CASCADE;      -- update FK when user id changes

-- ON DELETE options:
-- CASCADE     → delete child rows
-- SET NULL    → set FK to NULL
-- SET DEFAULT → set FK to default value
-- RESTRICT    → prevent parent deletion (same as NO ACTION for most purposes)
-- NO ACTION   → default; error if referenced row would be deleted
```

---

## Modifying Tables

```sql
-- Add column
ALTER TABLE users ADD COLUMN phone TEXT;
ALTER TABLE users ADD COLUMN avatar_url TEXT DEFAULT NULL;

-- Drop column
ALTER TABLE users DROP COLUMN phone;
ALTER TABLE users DROP COLUMN phone CASCADE;  -- also drops dependent objects

-- Rename column
ALTER TABLE users RENAME COLUMN bio TO biography;

-- Change data type
ALTER TABLE users ALTER COLUMN age TYPE INTEGER;

-- Set/drop NOT NULL
ALTER TABLE users ALTER COLUMN phone SET NOT NULL;
ALTER TABLE users ALTER COLUMN phone DROP NOT NULL;

-- Set/drop DEFAULT
ALTER TABLE users ALTER COLUMN is_active SET DEFAULT TRUE;
ALTER TABLE users ALTER COLUMN is_active DROP DEFAULT;

-- Rename table
ALTER TABLE users RENAME TO app_users;

-- Add constraint after creation
ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);
ALTER TABLE posts ADD CONSTRAINT valid_status
    CHECK (status IN ('draft', 'published', 'archived'));
```

---

## Generated Columns

```sql
CREATE TABLE orders (
    id          BIGSERIAL PRIMARY KEY,
    quantity    INTEGER NOT NULL,
    unit_price  NUMERIC(10, 2) NOT NULL,
    -- Automatically computed, stored on disk
    total       NUMERIC(10, 2) GENERATED ALWAYS AS (quantity * unit_price) STORED
);

-- Inserting doesn't include the generated column
INSERT INTO orders (quantity, unit_price) VALUES (3, 29.99);
SELECT * FROM orders;
-- total = 89.97 (computed automatically)
```

---

## Useful Table Queries

```sql
-- List all tables in current schema
\dt
SELECT table_name FROM information_schema.tables
WHERE table_schema = 'public';

-- Describe table structure
\d users
SELECT column_name, data_type, is_nullable, column_default
FROM information_schema.columns
WHERE table_name = 'users';

-- List constraints
SELECT constraint_name, constraint_type
FROM information_schema.table_constraints
WHERE table_name = 'users';

-- Table size
SELECT pg_size_pretty(pg_total_relation_size('users'));
SELECT pg_size_pretty(pg_relation_size('users'));          -- data only
SELECT pg_size_pretty(pg_indexes_size('users'));           -- indexes only
```

---

## Temporary Tables

```sql
-- Exists only for the current session
CREATE TEMP TABLE temp_results AS
SELECT id, name FROM users WHERE is_active = TRUE;

-- Gone automatically when session ends
```

---

## Table Inheritance

```sql
CREATE TABLE vehicles (
    id       BIGSERIAL PRIMARY KEY,
    make     TEXT NOT NULL,
    model    TEXT NOT NULL,
    year     INTEGER NOT NULL
);

CREATE TABLE cars (
    num_doors INTEGER NOT NULL DEFAULT 4
) INHERITS (vehicles);

CREATE TABLE trucks (
    payload_tons NUMERIC(6,2)
) INHERITS (vehicles);

-- Query parent includes all children
SELECT * FROM vehicles;       -- returns cars AND trucks

-- Query child only
SELECT * FROM ONLY vehicles;  -- returns only direct vehicles rows
```
# Chapter 4: CRUD Operations

## INSERT

```sql
-- Single row
INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com');

-- With all columns
INSERT INTO users (name, email, role, is_active)
VALUES ('Bob', 'bob@example.com', 'admin', TRUE);

-- Multiple rows (much faster than multiple INSERT statements)
INSERT INTO users (name, email) VALUES
    ('Charlie', 'charlie@example.com'),
    ('Diana', 'diana@example.com'),
    ('Eve', 'eve@example.com');

-- From a SELECT (copy data)
INSERT INTO users_backup
SELECT * FROM users WHERE is_active = FALSE;

-- RETURNING — get values back from inserted row
INSERT INTO users (name, email)
VALUES ('Frank', 'frank@example.com')
RETURNING id, created_at;

-- Get full inserted row
INSERT INTO users (name, email)
VALUES ('Grace', 'grace@example.com')
RETURNING *;
```

---

## SELECT

```sql
-- All rows, all columns
SELECT * FROM users;

-- Specific columns
SELECT id, name, email FROM users;

-- With WHERE filter
SELECT id, name FROM users WHERE is_active = TRUE;
SELECT id, name FROM users WHERE role = 'admin' AND is_active = TRUE;
SELECT id, name FROM users WHERE role IN ('admin', 'moderator');
SELECT id, name FROM users WHERE created_at > '2024-01-01';

-- Pattern matching
SELECT * FROM users WHERE name LIKE 'A%';       -- starts with A
SELECT * FROM users WHERE name ILIKE '%alice%'; -- case-insensitive
SELECT * FROM users WHERE email ~ '@gmail\.com$'; -- regex

-- NULL checks
SELECT * FROM users WHERE bio IS NULL;
SELECT * FROM users WHERE bio IS NOT NULL;

-- ORDER BY
SELECT id, name FROM users ORDER BY name ASC;
SELECT id, name FROM users ORDER BY created_at DESC, name ASC;

-- LIMIT and OFFSET (pagination)
SELECT id, name FROM users
ORDER BY created_at DESC
LIMIT 10 OFFSET 20;    -- page 3 (0-indexed), 10 per page

-- COUNT
SELECT COUNT(*) FROM users;
SELECT COUNT(*) FROM users WHERE is_active = TRUE;
SELECT COUNT(DISTINCT role) FROM users;

-- DISTINCT
SELECT DISTINCT role FROM users;
SELECT DISTINCT ON (role) id, role, name FROM users ORDER BY role, name;

-- Computed columns
SELECT
    id,
    name,
    UPPER(email) AS email_upper,
    EXTRACT(YEAR FROM created_at) AS year_joined,
    NOW() - created_at AS account_age
FROM users;
```

---

## Filtering with WHERE

```sql
-- Comparison operators
=   !=  <>  <  >  <=  >=

-- Logical operators
AND  OR  NOT

-- BETWEEN (inclusive)
SELECT * FROM users WHERE age BETWEEN 18 AND 65;

-- IN / NOT IN
SELECT * FROM users WHERE id IN (1, 2, 3);
SELECT * FROM users WHERE role NOT IN ('guest', 'banned');

-- ANY / ALL (with subquery)
SELECT * FROM orders WHERE amount > ANY (SELECT amount FROM orders WHERE user_id = 1);
SELECT * FROM orders WHERE amount > ALL (SELECT amount FROM orders WHERE user_id = 2);

-- EXISTS / NOT EXISTS
SELECT * FROM users u
WHERE EXISTS (
    SELECT 1 FROM orders o WHERE o.user_id = u.id
);

-- CASE expression
SELECT
    id,
    name,
    CASE
        WHEN age < 18 THEN 'minor'
        WHEN age BETWEEN 18 AND 65 THEN 'adult'
        ELSE 'senior'
    END AS age_group
FROM users;

-- COALESCE (first non-null)
SELECT id, COALESCE(phone, email, 'no contact') AS contact FROM users;

-- NULLIF (return null if values are equal)
SELECT NULLIF(score, 0) AS score FROM results;   -- turns 0 into NULL
```

---

## UPDATE

```sql
-- Update one column for all rows matching WHERE
UPDATE users SET is_active = FALSE WHERE last_login < NOW() - INTERVAL '1 year';

-- Update multiple columns
UPDATE users
SET
    name = 'Alice Smith',
    email = 'alice.smith@example.com',
    updated_at = NOW()
WHERE id = 1;

-- Update using RETURNING
UPDATE users
SET is_active = FALSE
WHERE last_login < NOW() - INTERVAL '1 year'
RETURNING id, name, email;

-- Update using a computed value
UPDATE products
SET price = price * 1.10    -- 10% price increase
WHERE category = 'electronics';

-- Update from another table (UPDATE ... FROM)
UPDATE orders o
SET status = 'cancelled'
FROM users u
WHERE o.user_id = u.id
AND u.is_active = FALSE
AND o.status = 'pending';
```

---

## DELETE

```sql
-- Delete with WHERE (always use WHERE unless you mean it!)
DELETE FROM users WHERE id = 42;
DELETE FROM sessions WHERE expires_at < NOW();

-- Delete with RETURNING
DELETE FROM notifications
WHERE read_at IS NOT NULL AND created_at < NOW() - INTERVAL '30 days'
RETURNING id, user_id;

-- Delete using a subquery
DELETE FROM posts
WHERE user_id IN (
    SELECT id FROM users WHERE is_active = FALSE
);

-- Delete with JOIN (using USING clause)
DELETE FROM orders o
USING users u
WHERE o.user_id = u.id
AND u.role = 'guest';

-- TRUNCATE — faster than DELETE for removing all rows
TRUNCATE TABLE sessions;
TRUNCATE TABLE sessions RESTART IDENTITY;   -- also reset sequences
TRUNCATE TABLE sessions CASCADE;            -- also truncate FK-referencing tables
```

---

## UPSERT — INSERT ON CONFLICT

UPSERT inserts a row, or updates it if it already exists (based on a unique constraint):

```sql
-- ON CONFLICT DO NOTHING (skip if duplicate)
INSERT INTO users (email, name)
VALUES ('alice@example.com', 'Alice')
ON CONFLICT (email) DO NOTHING;

-- ON CONFLICT DO UPDATE (upsert)
INSERT INTO users (email, name, role)
VALUES ('alice@example.com', 'Alice Updated', 'admin')
ON CONFLICT (email) DO UPDATE
SET
    name = EXCLUDED.name,          -- EXCLUDED refers to the rejected row
    role = EXCLUDED.role,
    updated_at = NOW();

-- Conditional upsert
INSERT INTO product_stock (product_id, quantity)
VALUES (42, 100)
ON CONFLICT (product_id) DO UPDATE
SET quantity = product_stock.quantity + EXCLUDED.quantity
WHERE product_stock.quantity < 1000;   -- only update if stock below threshold

-- RETURNING after upsert
INSERT INTO users (email, name)
VALUES ('alice@example.com', 'Alice')
ON CONFLICT (email) DO UPDATE
SET name = EXCLUDED.name
RETURNING id, name, (xmax = 0) AS was_inserted;
-- xmax = 0 means it was a fresh insert; xmax != 0 means it was updated
```

---

## Aggregate Functions

```sql
SELECT COUNT(*) FROM orders;
SELECT SUM(amount) FROM orders WHERE status = 'completed';
SELECT AVG(amount) FROM orders;
SELECT MIN(amount), MAX(amount) FROM orders;
SELECT ROUND(AVG(amount), 2) FROM orders;
SELECT STRING_AGG(name, ', ' ORDER BY name) FROM users;  -- concatenate strings
SELECT ARRAY_AGG(id ORDER BY id) FROM users;             -- collect into array

-- GROUP BY
SELECT
    status,
    COUNT(*) AS order_count,
    SUM(amount) AS total_revenue,
    AVG(amount) AS avg_order_value
FROM orders
GROUP BY status;

-- HAVING (filter after GROUP BY)
SELECT user_id, COUNT(*) AS order_count
FROM orders
GROUP BY user_id
HAVING COUNT(*) > 10;

-- GROUP BY with ROLLUP (subtotals)
SELECT
    category,
    status,
    SUM(amount) AS revenue
FROM orders
GROUP BY ROLLUP (category, status);
-- Adds subtotal rows for each category and a grand total
```

---

## Useful String Functions

```sql
SELECT UPPER('hello');                    -- 'HELLO'
SELECT LOWER('WORLD');                    -- 'world'
SELECT LENGTH('hello');                   -- 5
SELECT TRIM('  hello  ');                 -- 'hello'
SELECT LTRIM('  hello');                  -- 'hello'
SELECT SUBSTR('hello world', 7, 5);       -- 'world'
SELECT POSITION('world' IN 'hello world'); -- 7
SELECT REPLACE('hello world', 'world', 'postgres'); -- 'hello postgres'
SELECT SPLIT_PART('a,b,c', ',', 2);       -- 'b'
SELECT CONCAT(first_name, ' ', last_name) FROM users;
SELECT first_name || ' ' || last_name FROM users;   -- same as CONCAT
SELECT FORMAT('%s has %s posts', name, post_count) FROM ...;
SELECT LPAD('42', 5, '0');                -- '00042'
SELECT MD5('password');                   -- hash string
```

---

## Useful Date Functions

```sql
SELECT NOW();                             -- current timestamp with tz
SELECT CURRENT_DATE;                      -- current date
SELECT CURRENT_TIME;                      -- current time with tz
SELECT EXTRACT(YEAR FROM NOW());          -- 2024
SELECT EXTRACT(MONTH FROM NOW());         -- 1
SELECT DATE_TRUNC('month', NOW());        -- first moment of current month
SELECT NOW() + INTERVAL '7 days';        -- 7 days from now
SELECT NOW() - INTERVAL '1 year';        -- 1 year ago
SELECT AGE(NOW(), '1990-01-01');          -- age from date
SELECT TO_CHAR(NOW(), 'YYYY-MM-DD');      -- format timestamp as string
SELECT TO_TIMESTAMP('2024-01-15', 'YYYY-MM-DD');   -- string to timestamp
```
# Chapter 5: Advanced Queries

## JOINs

JOINs combine rows from two or more tables based on a related column.

```sql
-- INNER JOIN — only rows with matches in both tables
SELECT u.name, o.id AS order_id, o.amount
FROM users u
INNER JOIN orders o ON o.user_id = u.id;

-- LEFT JOIN — all rows from left table, nulls for non-matching right
SELECT u.name, COUNT(o.id) AS order_count
FROM users u
LEFT JOIN orders o ON o.user_id = u.id
GROUP BY u.id, u.name;

-- RIGHT JOIN — all rows from right table (less common, swap tables instead)
SELECT u.name, o.amount
FROM orders o
RIGHT JOIN users u ON o.user_id = u.id;

-- FULL OUTER JOIN — all rows from both tables
SELECT u.name, o.amount
FROM users u
FULL OUTER JOIN orders o ON o.user_id = u.id;

-- CROSS JOIN — cartesian product (every row with every row)
SELECT color, size FROM colors CROSS JOIN sizes;

-- SELF JOIN — join a table to itself
SELECT a.name AS employee, b.name AS manager
FROM employees a
JOIN employees b ON a.manager_id = b.id;
```

### Multiple JOINs:
```sql
SELECT
    u.name AS customer,
    o.id AS order_id,
    p.name AS product,
    oi.quantity,
    oi.unit_price
FROM orders o
JOIN users u ON u.id = o.user_id
JOIN order_items oi ON oi.order_id = o.id
JOIN products p ON p.id = oi.product_id
WHERE o.status = 'completed'
ORDER BY o.created_at DESC;
```

---

## Subqueries

A subquery is a query nested inside another query.

```sql
-- Subquery in WHERE (scalar subquery)
SELECT name FROM users
WHERE id = (SELECT user_id FROM orders WHERE amount = (SELECT MAX(amount) FROM orders));

-- Subquery in WHERE with IN
SELECT name FROM users
WHERE id IN (
    SELECT DISTINCT user_id FROM orders WHERE status = 'completed'
);

-- Correlated subquery (references outer query)
SELECT u.name,
    (SELECT COUNT(*) FROM orders WHERE user_id = u.id) AS order_count
FROM users u;

-- Subquery in FROM (derived table)
SELECT dept, AVG(salary) AS avg_salary
FROM (
    SELECT department AS dept, salary
    FROM employees
    WHERE hire_date > '2020-01-01'
) recent_hires
GROUP BY dept;

-- EXISTS with correlated subquery
SELECT * FROM users u
WHERE EXISTS (
    SELECT 1 FROM orders o
    WHERE o.user_id = u.id AND o.amount > 1000
);

-- NOT EXISTS
SELECT * FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM orders o WHERE o.user_id = u.id
);
```

---

## Common Table Expressions (CTEs)

CTEs (`WITH` clause) are named temporary result sets that make complex queries more readable.

```sql
-- Basic CTE
WITH active_users AS (
    SELECT id, name, email
    FROM users
    WHERE is_active = TRUE
),
high_value_orders AS (
    SELECT user_id, SUM(amount) AS total
    FROM orders
    WHERE status = 'completed'
    GROUP BY user_id
    HAVING SUM(amount) > 10000
)
SELECT u.name, u.email, hvo.total
FROM active_users u
JOIN high_value_orders hvo ON hvo.user_id = u.id
ORDER BY hvo.total DESC;

-- Multiple CTEs
WITH
    monthly_revenue AS (
        SELECT DATE_TRUNC('month', created_at) AS month, SUM(amount) AS revenue
        FROM orders
        WHERE status = 'completed'
        GROUP BY 1
    ),
    prev_month AS (
        SELECT LAG(revenue) OVER (ORDER BY month) AS prev_revenue, month, revenue
        FROM monthly_revenue
    )
SELECT
    month,
    revenue,
    prev_revenue,
    ROUND((revenue - prev_revenue) / prev_revenue * 100, 2) AS growth_pct
FROM prev_month
ORDER BY month;
```

### Recursive CTEs — for tree/graph structures:
```sql
-- Org chart traversal
WITH RECURSIVE org_chart AS (
    -- Base case: start with the CEO (no manager)
    SELECT id, name, manager_id, 1 AS level
    FROM employees
    WHERE manager_id IS NULL

    UNION ALL

    -- Recursive case: join with direct reports
    SELECT e.id, e.name, e.manager_id, oc.level + 1
    FROM employees e
    JOIN org_chart oc ON e.manager_id = oc.id
)
SELECT level, name FROM org_chart ORDER BY level, name;

-- Find all ancestors of a node
WITH RECURSIVE ancestors AS (
    SELECT id, parent_id, name
    FROM categories
    WHERE id = 42        -- starting node

    UNION ALL

    SELECT c.id, c.parent_id, c.name
    FROM categories c
    JOIN ancestors a ON c.id = a.parent_id
)
SELECT * FROM ancestors;
```

---

## Window Functions

Window functions perform calculations across a set of rows related to the current row, without collapsing them like GROUP BY.

```sql
-- ROW_NUMBER — unique sequential number within partition
SELECT
    id,
    name,
    department,
    salary,
    ROW_NUMBER() OVER (PARTITION BY department ORDER BY salary DESC) AS rank_in_dept
FROM employees;

-- RANK — like ROW_NUMBER but ties get same rank (with gaps)
-- DENSE_RANK — like RANK but no gaps

-- LAG / LEAD — access previous/next row
SELECT
    month,
    revenue,
    LAG(revenue, 1) OVER (ORDER BY month) AS prev_month_revenue,
    LEAD(revenue, 1) OVER (ORDER BY month) AS next_month_revenue,
    revenue - LAG(revenue) OVER (ORDER BY month) AS change
FROM monthly_revenue;

-- Running total
SELECT
    created_at::DATE AS date,
    amount,
    SUM(amount) OVER (ORDER BY created_at) AS running_total
FROM orders
WHERE status = 'completed';

-- Moving average (7-day)
SELECT
    date,
    revenue,
    AVG(revenue) OVER (
        ORDER BY date
        ROWS BETWEEN 6 PRECEDING AND CURRENT ROW
    ) AS moving_avg_7d
FROM daily_revenue;

-- FIRST_VALUE / LAST_VALUE
SELECT
    id,
    name,
    salary,
    FIRST_VALUE(salary) OVER (PARTITION BY department ORDER BY salary DESC) AS top_salary_in_dept
FROM employees;

-- NTILE — divide rows into N buckets
SELECT
    id,
    name,
    salary,
    NTILE(4) OVER (ORDER BY salary) AS quartile
FROM employees;

-- PERCENT_RANK / CUME_DIST
SELECT
    name,
    salary,
    ROUND(PERCENT_RANK() OVER (ORDER BY salary)::numeric * 100, 2) AS percentile
FROM employees;
```

---

## GROUPING SETS, ROLLUP, CUBE

```sql
-- ROLLUP — hierarchical subtotals
SELECT
    region,
    category,
    SUM(amount) AS revenue
FROM sales
GROUP BY ROLLUP (region, category);
-- Rows: (region, category), (region, NULL), (NULL, NULL)

-- CUBE — all possible combinations of subtotals
SELECT
    region,
    category,
    year,
    SUM(amount)
FROM sales
GROUP BY CUBE (region, category, year);

-- GROUPING SETS — explicit combinations
SELECT
    region,
    category,
    SUM(amount)
FROM sales
GROUP BY GROUPING SETS (
    (region, category),
    (region),
    (category),
    ()             -- grand total
);

-- GROUPING() function — tells you which columns were aggregated
SELECT
    GROUPING(region) AS region_grouped,
    region,
    SUM(amount)
FROM sales
GROUP BY ROLLUP (region);
```

---

## LATERAL Joins

LATERAL allows a subquery in FROM to reference columns from earlier in the FROM clause:

```sql
-- Get the 3 most recent orders for each user
SELECT u.name, recent.id, recent.amount, recent.created_at
FROM users u
CROSS JOIN LATERAL (
    SELECT id, amount, created_at
    FROM orders
    WHERE user_id = u.id
    ORDER BY created_at DESC
    LIMIT 3
) AS recent;

-- Get top product per category
SELECT c.name AS category, top_product.*
FROM categories c
CROSS JOIN LATERAL (
    SELECT p.name, p.price
    FROM products p
    WHERE p.category_id = c.id
    ORDER BY p.sales_count DESC
    LIMIT 1
) AS top_product;
```

---

## UNION, INTERSECT, EXCEPT

```sql
-- UNION — combine results, remove duplicates
SELECT name FROM customers
UNION
SELECT name FROM vendors;

-- UNION ALL — combine results, keep duplicates (faster)
SELECT name, 'customer' AS type FROM customers
UNION ALL
SELECT name, 'vendor' AS type FROM vendors;

-- INTERSECT — rows in both result sets
SELECT user_id FROM premium_users
INTERSECT
SELECT user_id FROM active_this_month;

-- EXCEPT — rows in first but not second
SELECT user_id FROM all_users
EXCEPT
SELECT user_id FROM cancelled_users;
```
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
# Chapter 7: Transactions & ACID

## What is a Transaction?

A **transaction** is a sequence of SQL operations that are treated as a single unit — either ALL succeed, or NONE do (atomicity). This prevents partial updates that leave your database in an inconsistent state.

```sql
-- Classic example: bank transfer
BEGIN;
    UPDATE accounts SET balance = balance - 100 WHERE id = 1;  -- debit
    UPDATE accounts SET balance = balance + 100 WHERE id = 2;  -- credit
COMMIT;
-- If any statement fails, ROLLBACK is triggered and neither update persists
```

---

## Transaction Commands

```sql
BEGIN;           -- or START TRANSACTION;
COMMIT;          -- save all changes
ROLLBACK;        -- discard all changes

-- Savepoints — partial rollback within a transaction
BEGIN;
    INSERT INTO orders (...) VALUES (...);
    SAVEPOINT order_saved;

    INSERT INTO order_items (...) VALUES (...);  -- this fails
    ROLLBACK TO order_saved;                     -- undo just the items insert

    -- order is still there, continue differently
    INSERT INTO order_items (...) VALUES (...);  -- different items
COMMIT;

-- Release savepoint (frees memory, doesn't commit)
RELEASE SAVEPOINT order_saved;
```

---

## ACID Properties in Practice

### Atomicity:
```sql
BEGIN;
    UPDATE inventory SET stock = stock - 5 WHERE product_id = 42;
    INSERT INTO order_items (order_id, product_id, quantity) VALUES (100, 42, 5);
    -- If INSERT fails (e.g., constraint), UPDATE is also rolled back
COMMIT;
```

### Consistency:
```sql
-- CHECK constraints, FOREIGN KEY constraints, NOT NULL — all enforced at commit
BEGIN;
    INSERT INTO users (email, age) VALUES ('test@test.com', -5);  -- violates CHECK(age >= 0)
    -- PostgreSQL raises error immediately
ROLLBACK;
```

### Isolation:
```sql
-- Other transactions don't see uncommitted changes
-- Transaction A:
BEGIN;
    UPDATE users SET name = 'Alice Updated' WHERE id = 1;
    -- (not committed yet)

-- Transaction B (concurrent):
SELECT name FROM users WHERE id = 1;
-- Returns: 'Alice' (the old value — B doesn't see A's uncommitted change)

-- Transaction A:
COMMIT;
-- Now Transaction B (or new queries) can see 'Alice Updated'
```

### Durability:
```sql
-- Once COMMIT returns, the change is safe even if the server crashes
-- PostgreSQL uses WAL (Write-Ahead Log) to guarantee this
COMMIT;  -- after this returns, the change survives a crash
```

---

## Isolation Levels

Isolation levels control how much a transaction can see of other concurrent transactions' changes.

| Level | Dirty Read | Non-Repeatable Read | Phantom Read |
|---|---|---|---|
| READ UNCOMMITTED | Possible (PG: same as RC) | Possible | Possible |
| READ COMMITTED | Not Possible | Possible | Possible |
| REPEATABLE READ | Not Possible | Not Possible | Not Possible (PG) |
| SERIALIZABLE | Not Possible | Not Possible | Not Possible |

### Read Committed (default in PostgreSQL):
```sql
-- Each statement sees the latest committed data at statement start
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;
BEGIN;
    -- This is the PostgreSQL default — good for most applications
    SELECT * FROM orders;   -- sees committed data as of this statement
    -- ... time passes, another transaction commits ...
    SELECT * FROM orders;   -- may see different data! (non-repeatable read)
COMMIT;
```

### Repeatable Read:
```sql
BEGIN TRANSACTION ISOLATION LEVEL REPEATABLE READ;
    SELECT * FROM orders;   -- sees snapshot at transaction start
    -- ... another transaction commits new orders ...
    SELECT * FROM orders;   -- SAME results (snapshot is frozen)
    -- No phantom reads in PostgreSQL's REPEATABLE READ (stronger than SQL standard)
COMMIT;
```

### Serializable:
```sql
BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;
    -- Strongest isolation — transactions appear to run one after another
    -- May fail with "could not serialize access" if concurrent conflict detected
    -- Retry the transaction on serialization failure
COMMIT;
```

---

## Locking

### Row-Level Locks:
```sql
-- SELECT FOR UPDATE — lock rows for update (prevents concurrent modification)
BEGIN;
    SELECT * FROM accounts WHERE id = 1 FOR UPDATE;
    -- Other transactions trying SELECT FOR UPDATE on the same row will wait
    UPDATE accounts SET balance = balance - 100 WHERE id = 1;
COMMIT;

-- FOR SHARE — lock rows for reading (others can read, not update)
SELECT * FROM products WHERE id = 42 FOR SHARE;

-- SKIP LOCKED — skip rows that are already locked
SELECT * FROM job_queue WHERE status = 'pending' LIMIT 1 FOR UPDATE SKIP LOCKED;
-- Great for building work queues — each worker gets a different job

-- NOWAIT — fail immediately instead of waiting for a lock
SELECT * FROM accounts WHERE id = 1 FOR UPDATE NOWAIT;
-- Raises error if row is locked (don't block)
```

### Table-Level Locks:
```sql
-- Lock entire table (usually done automatically by DDL)
LOCK TABLE users IN ACCESS EXCLUSIVE MODE;

-- Lock modes (least to most restrictive):
-- ACCESS SHARE         — SELECT
-- ROW SHARE            — SELECT FOR UPDATE
-- ROW EXCLUSIVE        — INSERT, UPDATE, DELETE
-- SHARE UPDATE EXCLUSIVE — VACUUM, CREATE INDEX CONCURRENTLY
-- SHARE                — CREATE INDEX
-- SHARE ROW EXCLUSIVE  — (rare)
-- EXCLUSIVE            — (rare)
-- ACCESS EXCLUSIVE     — DROP TABLE, ALTER TABLE, etc.
```

### Advisory Locks (application-level):
```sql
-- Lock by arbitrary integer (application assigns meaning)
SELECT pg_advisory_lock(12345);        -- session-level, blocks until acquired
SELECT pg_try_advisory_lock(12345);    -- returns true/false immediately
SELECT pg_advisory_unlock(12345);      -- release

-- Transaction-level advisory lock (released at commit/rollback)
SELECT pg_advisory_xact_lock(12345);
```

---

## Deadlocks

A deadlock occurs when two transactions each hold a lock the other needs:

```
Transaction A holds lock on accounts(1), wants lock on accounts(2)
Transaction B holds lock on accounts(2), wants lock on accounts(1)
→ Both wait forever → PostgreSQL detects and kills one
```

**Prevention:** Always acquire locks in the same order across transactions.

```sql
-- Bad: A does (1→2), B does (2→1) → deadlock possible
-- Good: Both always lock in order (1→2) → deadlock impossible

BEGIN;
    -- Always lock the lower ID first
    SELECT * FROM accounts WHERE id IN (1, 2) ORDER BY id FOR UPDATE;
    -- Now safe to update
COMMIT;
```

---

## Monitoring Transactions and Locks

```sql
-- See current active transactions
SELECT pid, state, wait_event_type, wait_event, query_start, query
FROM pg_stat_activity
WHERE state != 'idle';

-- See currently held locks
SELECT l.pid, l.relation::regclass, l.mode, l.granted, a.query
FROM pg_locks l
JOIN pg_stat_activity a ON a.pid = l.pid
WHERE l.relation IS NOT NULL;

-- See blocking queries
SELECT blocked.pid, blocked.query, blocking.pid AS blocking_pid, blocking.query AS blocking_query
FROM pg_stat_activity blocked
JOIN pg_stat_activity blocking ON blocking.pid = ANY(pg_blocking_pids(blocked.pid))
WHERE cardinality(pg_blocking_pids(blocked.pid)) > 0;

-- Kill a blocking query (replace 1234 with actual pid)
SELECT pg_cancel_backend(1234);   -- gentle: cancel current query
SELECT pg_terminate_backend(1234); -- forceful: disconnect entirely
```

---

## Long-Running Transactions

Long-running transactions block VACUUM, autovacuum, and cause table bloat.

```sql
-- Find transactions older than 5 minutes
SELECT pid, now() - pg_stat_activity.query_start AS duration, query, state
FROM pg_stat_activity
WHERE (now() - pg_stat_activity.query_start) > INTERVAL '5 minutes'
AND state != 'idle';

-- Set statement timeout (kills query if it runs too long)
SET statement_timeout = '30s';

-- Set idle in transaction timeout
SET idle_in_transaction_session_timeout = '60s';

-- In postgresql.conf (server-wide default)
statement_timeout = '30000'              # 30 seconds
idle_in_transaction_session_timeout = '60000'  # 60 seconds
```

---

## Two-Phase Commit (Distributed Transactions)

For distributed systems spanning multiple databases:

```sql
-- Phase 1: Prepare (vote to commit)
BEGIN;
    UPDATE accounts SET balance = balance - 100 WHERE id = 1;
PREPARE TRANSACTION 'transfer_123';   -- returns control, transaction suspended

-- Phase 2: Commit or rollback (after all participants agree)
COMMIT PREPARED 'transfer_123';
-- or
ROLLBACK PREPARED 'transfer_123';

-- List prepared transactions
SELECT * FROM pg_prepared_xacts;
```
# Chapter 8: Functions, Procedures & Triggers

## PL/pgSQL — PostgreSQL's Procedural Language

PL/pgSQL is the built-in procedural language for writing functions, procedures, and trigger logic in PostgreSQL.

---

## Functions

```sql
-- Basic function
CREATE OR REPLACE FUNCTION greet(name TEXT)
RETURNS TEXT AS $$
BEGIN
    RETURN 'Hello, ' || name || '!';
END;
$$ LANGUAGE plpgsql;

SELECT greet('World');    -- Hello, World!

-- Function with multiple parameters and computation
CREATE OR REPLACE FUNCTION calculate_tax(amount NUMERIC, rate NUMERIC DEFAULT 0.10)
RETURNS NUMERIC AS $$
BEGIN
    RETURN ROUND(amount * rate, 2);
END;
$$ LANGUAGE plpgsql;

SELECT calculate_tax(100.00);          -- 10.00
SELECT calculate_tax(100.00, 0.20);   -- 20.00
```

### SQL Functions (simpler for single expressions):
```sql
CREATE OR REPLACE FUNCTION full_name(first_name TEXT, last_name TEXT)
RETURNS TEXT AS $$
    SELECT first_name || ' ' || last_name;
$$ LANGUAGE sql IMMUTABLE;

-- IMMUTABLE: same inputs always return same output (enables function-based indexes)
-- STABLE: result is stable within a transaction (can use in index scans)
-- VOLATILE: result can change (default — cannot be used in index scans)
```

---

## Control Flow

```sql
CREATE OR REPLACE FUNCTION get_user_tier(user_id BIGINT)
RETURNS TEXT AS $$
DECLARE
    total_spending NUMERIC;
    tier TEXT;
BEGIN
    -- Declare variables, then assign
    SELECT COALESCE(SUM(amount), 0)
    INTO total_spending
    FROM orders
    WHERE user_id = get_user_tier.user_id
    AND status = 'completed';

    -- IF / ELSIF / ELSE
    IF total_spending >= 10000 THEN
        tier := 'platinum';
    ELSIF total_spending >= 1000 THEN
        tier := 'gold';
    ELSIF total_spending >= 100 THEN
        tier := 'silver';
    ELSE
        tier := 'bronze';
    END IF;

    RETURN tier;
END;
$$ LANGUAGE plpgsql;
```

### CASE expression:
```sql
CREATE OR REPLACE FUNCTION day_type(d DATE)
RETURNS TEXT AS $$
BEGIN
    RETURN CASE EXTRACT(DOW FROM d)
        WHEN 0, 6 THEN 'weekend'
        ELSE 'weekday'
    END;
END;
$$ LANGUAGE plpgsql;
```

---

## Loops

```sql
CREATE OR REPLACE FUNCTION sum_1_to_n(n INTEGER)
RETURNS BIGINT AS $$
DECLARE
    total BIGINT := 0;
    i INTEGER;
BEGIN
    -- Basic loop with EXIT
    i := 1;
    LOOP
        EXIT WHEN i > n;
        total := total + i;
        i := i + 1;
    END LOOP;
    RETURN total;
END;
$$ LANGUAGE plpgsql;

-- FOR loop (integer range)
CREATE OR REPLACE FUNCTION generate_invoice_numbers(year INTEGER, count INTEGER)
RETURNS VOID AS $$
DECLARE
    i INTEGER;
BEGIN
    FOR i IN 1..count LOOP
        INSERT INTO invoices (number, year)
        VALUES (LPAD(i::TEXT, 6, '0'), year);
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- FOR loop over query results
CREATE OR REPLACE FUNCTION deactivate_old_sessions()
RETURNS INTEGER AS $$
DECLARE
    session RECORD;
    deleted_count INTEGER := 0;
BEGIN
    FOR session IN
        SELECT id FROM sessions WHERE expires_at < NOW()
    LOOP
        DELETE FROM sessions WHERE id = session.id;
        deleted_count := deleted_count + 1;
    END LOOP;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;
```

---

## Exception Handling

```sql
CREATE OR REPLACE FUNCTION safe_divide(a NUMERIC, b NUMERIC)
RETURNS NUMERIC AS $$
BEGIN
    RETURN a / b;
EXCEPTION
    WHEN division_by_zero THEN
        RETURN NULL;
    WHEN OTHERS THEN
        RAISE NOTICE 'Unexpected error: %', SQLERRM;
        RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Raise errors from functions
CREATE OR REPLACE FUNCTION transfer_funds(
    from_id BIGINT, to_id BIGINT, amount NUMERIC
)
RETURNS VOID AS $$
DECLARE
    current_balance NUMERIC;
BEGIN
    SELECT balance INTO current_balance FROM accounts WHERE id = from_id FOR UPDATE;

    IF current_balance < amount THEN
        RAISE EXCEPTION 'Insufficient funds: balance=%, requested=%',
            current_balance, amount
            USING ERRCODE = 'insufficient_privilege';  -- or custom SQLSTATE
    END IF;

    UPDATE accounts SET balance = balance - amount WHERE id = from_id;
    UPDATE accounts SET balance = balance + amount WHERE id = to_id;
END;
$$ LANGUAGE plpgsql;
```

---

## Returning Tables and SETOF

```sql
-- Return a table
CREATE OR REPLACE FUNCTION get_top_users(n INTEGER)
RETURNS TABLE(user_id BIGINT, user_name TEXT, total_orders BIGINT) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, COUNT(o.id)
        FROM users u
        LEFT JOIN orders o ON o.user_id = u.id
        GROUP BY u.id, u.name
        ORDER BY COUNT(o.id) DESC
        LIMIT n;
END;
$$ LANGUAGE plpgsql;

SELECT * FROM get_top_users(10);

-- SETOF existing type
CREATE OR REPLACE FUNCTION get_active_users()
RETURNS SETOF users AS $$
BEGIN
    RETURN QUERY SELECT * FROM users WHERE is_active = TRUE;
END;
$$ LANGUAGE plpgsql;
```

---

## Stored Procedures (PostgreSQL 11+)

Procedures can manage transactions (COMMIT/ROLLBACK within them):

```sql
CREATE OR REPLACE PROCEDURE batch_update_prices(increase_pct NUMERIC)
LANGUAGE plpgsql AS $$
DECLARE
    batch RECORD;
    updated INTEGER := 0;
BEGIN
    FOR batch IN SELECT id FROM products WHERE is_active = TRUE LOOP
        UPDATE products
        SET price = price * (1 + increase_pct / 100)
        WHERE id = batch.id;

        updated := updated + 1;

        -- Commit every 1000 rows to avoid huge transactions
        IF updated % 1000 = 0 THEN
            COMMIT;
            RAISE NOTICE 'Committed % rows', updated;
        END IF;
    END LOOP;
    COMMIT;
    RAISE NOTICE 'Done. Updated % products.', updated;
END;
$$;

-- Call with CALL
CALL batch_update_prices(10.0);
```

---

## Triggers

A trigger automatically executes a function before or after a data-changing event.

### Trigger Function:
```sql
-- Trigger function must return TRIGGER type
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;    -- RETURN NEW for BEFORE triggers on INSERT/UPDATE
                   -- RETURN NULL to cancel the operation
END;
$$ LANGUAGE plpgsql;
```

### Attaching the Trigger:
```sql
-- Fire BEFORE each row INSERT or UPDATE
CREATE TRIGGER trigger_update_timestamp
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Now every update automatically sets updated_at
UPDATE users SET name = 'New Name' WHERE id = 1;
-- updated_at is automatically set
```

### Audit Log Trigger:
```sql
CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    table_name TEXT,
    operation TEXT,
    old_data JSONB,
    new_data JSONB,
    changed_by TEXT,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION audit_trigger_function()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO audit_log (table_name, operation, old_data, new_data, changed_by)
    VALUES (
        TG_TABLE_NAME,
        TG_OP,
        CASE WHEN TG_OP IN ('UPDATE', 'DELETE') THEN row_to_json(OLD)::JSONB END,
        CASE WHEN TG_OP IN ('INSERT', 'UPDATE') THEN row_to_json(NEW)::JSONB END,
        current_user
    );

    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- Apply to any table
CREATE TRIGGER audit_users
AFTER INSERT OR UPDATE OR DELETE ON users
FOR EACH ROW EXECUTE FUNCTION audit_trigger_function();
```

### Trigger Timing:

| Timing | Fired | Row Variables |
|---|---|---|
| `BEFORE` | Before the operation | NEW and/or OLD available; can modify NEW |
| `AFTER` | After the operation | NEW and/or OLD available; cannot modify |
| `INSTEAD OF` | Instead of the operation (on views) | Implement view DML |

### Statement-Level Triggers:
```sql
-- Fires once per statement, not per row
CREATE TRIGGER log_bulk_delete
AFTER DELETE ON orders
FOR EACH STATEMENT
EXECUTE FUNCTION log_bulk_operation();
```

---

## Managing Functions and Triggers

```sql
-- List functions
\df
SELECT routine_name, routine_type, data_type
FROM information_schema.routines
WHERE routine_schema = 'public';

-- List triggers
\dy
SELECT trigger_name, event_manipulation, event_object_table, action_timing
FROM information_schema.triggers
WHERE trigger_schema = 'public';

-- Drop function
DROP FUNCTION greet(TEXT);
DROP FUNCTION IF EXISTS greet(TEXT);

-- Drop trigger
DROP TRIGGER trigger_update_timestamp ON users;
DROP TRIGGER IF EXISTS trigger_update_timestamp ON users;

-- Enable/disable trigger
ALTER TABLE users DISABLE TRIGGER trigger_update_timestamp;
ALTER TABLE users ENABLE TRIGGER trigger_update_timestamp;
ALTER TABLE users DISABLE TRIGGER ALL;
```
# Chapter 9: Views & Materialized Views

## What is a View?

A **view** is a named SQL query stored in the database. It behaves like a table but doesn't store data — each access runs the underlying query. Views are used for:
- Simplifying complex queries
- Hiding sensitive columns
- Providing a stable API over changing table structures
- Enforcing row-level access patterns

---

## Creating and Using Views

```sql
-- Basic view
CREATE VIEW active_users AS
SELECT id, name, email, created_at
FROM users
WHERE is_active = TRUE;

-- Use like a table
SELECT * FROM active_users;
SELECT name FROM active_users WHERE name LIKE 'A%';

-- View joining multiple tables
CREATE OR REPLACE VIEW order_summary AS
SELECT
    o.id AS order_id,
    u.name AS customer_name,
    u.email AS customer_email,
    o.status,
    o.created_at,
    COUNT(oi.id) AS item_count,
    SUM(oi.quantity * oi.unit_price) AS total_amount
FROM orders o
JOIN users u ON u.id = o.user_id
LEFT JOIN order_items oi ON oi.order_id = o.id
GROUP BY o.id, u.name, u.email, o.status, o.created_at;

SELECT * FROM order_summary WHERE status = 'completed';
SELECT customer_name, SUM(total_amount) FROM order_summary GROUP BY customer_name;
```

---

## Updatable Views

Simple views (no aggregation, no DISTINCT, no JOIN) are automatically updatable:

```sql
CREATE VIEW active_products AS
SELECT id, name, price, stock
FROM products
WHERE is_active = TRUE;

-- These work on simple views
INSERT INTO active_products (name, price, stock) VALUES ('Widget', 9.99, 100);
UPDATE active_products SET price = 12.99 WHERE id = 42;
DELETE FROM active_products WHERE id = 42;
```

### WITH CHECK OPTION — enforce view conditions on INSERT/UPDATE:
```sql
CREATE VIEW active_products AS
SELECT id, name, price, is_active
FROM products
WHERE is_active = TRUE
WITH CHECK OPTION;

-- This will FAIL — would insert a row not visible through the view
INSERT INTO active_products (name, price, is_active) VALUES ('Hidden', 5.00, FALSE);
-- ERROR: new row violates check option for view "active_products"
```

---

## INSTEAD OF Triggers (for complex views)

For views that aren't automatically updatable:

```sql
CREATE VIEW user_profile AS
SELECT u.id, u.name, u.email, p.bio, p.avatar_url
FROM users u
LEFT JOIN profiles p ON p.user_id = u.id;

-- INSTEAD OF trigger makes this view updateable
CREATE OR REPLACE FUNCTION update_user_profile()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users SET name = NEW.name, email = NEW.email WHERE id = NEW.id;
    INSERT INTO profiles (user_id, bio, avatar_url)
    VALUES (NEW.id, NEW.bio, NEW.avatar_url)
    ON CONFLICT (user_id) DO UPDATE
    SET bio = EXCLUDED.bio, avatar_url = EXCLUDED.avatar_url;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_profile_trigger
INSTEAD OF UPDATE ON user_profile
FOR EACH ROW EXECUTE FUNCTION update_user_profile();

-- Now this works
UPDATE user_profile SET bio = 'New bio', name = 'New name' WHERE id = 1;
```

---

## Security Views — Column and Row Security

```sql
-- Hide sensitive columns
CREATE VIEW users_public AS
SELECT id, name, created_at  -- no email, no password_hash
FROM users
WHERE is_active = TRUE;

-- Row-level filtering per user
CREATE VIEW my_orders AS
SELECT * FROM orders
WHERE user_id = current_setting('app.current_user_id')::BIGINT;
-- Application sets: SET LOCAL app.current_user_id = '42';
```

---

## Materialized Views

A **materialized view** stores the query result physically on disk. Unlike regular views, queries against it read cached data — very fast for expensive aggregations.

```sql
-- Create materialized view
CREATE MATERIALIZED VIEW monthly_revenue AS
SELECT
    DATE_TRUNC('month', created_at) AS month,
    COUNT(*) AS order_count,
    SUM(amount) AS total_revenue,
    AVG(amount) AS avg_order_value
FROM orders
WHERE status = 'completed'
GROUP BY DATE_TRUNC('month', created_at)
ORDER BY month;

-- Index the materialized view
CREATE INDEX idx_monthly_revenue_month ON monthly_revenue(month);

-- Query it like a table (reads from cache — instant)
SELECT * FROM monthly_revenue WHERE month > NOW() - INTERVAL '1 year';
```

---

## Refreshing Materialized Views

The data goes stale — you must refresh it:

```sql
-- REFRESH — blocks reads while refreshing
REFRESH MATERIALIZED VIEW monthly_revenue;

-- REFRESH CONCURRENTLY — allows reads during refresh (requires unique index)
CREATE UNIQUE INDEX idx_monthly_revenue_month_unique ON monthly_revenue(month);
REFRESH MATERIALIZED VIEW CONCURRENTLY monthly_revenue;
```

### Automated refresh with pg_cron:
```sql
-- Install pg_cron extension
CREATE EXTENSION pg_cron;

-- Refresh every hour
SELECT cron.schedule('refresh-monthly-revenue', '0 * * * *',
    'REFRESH MATERIALIZED VIEW CONCURRENTLY monthly_revenue');

-- Refresh nightly at 2 AM
SELECT cron.schedule('refresh-product-stats', '0 2 * * *',
    'REFRESH MATERIALIZED VIEW CONCURRENTLY product_stats');
```

---

## Practical Materialized View Examples

### Product statistics:
```sql
CREATE MATERIALIZED VIEW product_stats AS
SELECT
    p.id,
    p.name,
    p.category,
    COUNT(oi.id) AS times_ordered,
    SUM(oi.quantity) AS total_units_sold,
    SUM(oi.quantity * oi.unit_price) AS total_revenue,
    AVG(r.rating) AS avg_rating,
    COUNT(r.id) AS review_count
FROM products p
LEFT JOIN order_items oi ON oi.product_id = p.id
LEFT JOIN reviews r ON r.product_id = p.id
GROUP BY p.id, p.name, p.category;

CREATE UNIQUE INDEX ON product_stats(id);
CREATE INDEX ON product_stats(category, total_revenue DESC);

-- Instant product rankings
SELECT name, total_revenue FROM product_stats
ORDER BY total_revenue DESC LIMIT 10;
```

### User activity summary:
```sql
CREATE MATERIALIZED VIEW user_activity_summary AS
SELECT
    u.id AS user_id,
    u.name,
    u.email,
    u.created_at AS joined_at,
    COUNT(DISTINCT o.id) AS total_orders,
    COALESCE(SUM(o.amount), 0) AS lifetime_value,
    MAX(o.created_at) AS last_order_at,
    COUNT(DISTINCT p.id) AS total_posts
FROM users u
LEFT JOIN orders o ON o.user_id = u.id AND o.status = 'completed'
LEFT JOIN posts p ON p.author_id = u.id
GROUP BY u.id, u.name, u.email, u.created_at;

CREATE UNIQUE INDEX ON user_activity_summary(user_id);
CREATE INDEX ON user_activity_summary(lifetime_value DESC);
```

---

## View vs Materialized View

| Feature | View | Materialized View |
|---|---|---|
| Data storage | None (runs query) | Yes (cached result) |
| Query performance | Depends on query | Fast (reads from disk) |
| Data freshness | Always current | Stale until refreshed |
| Index support | Uses base table indexes | Can have its own indexes |
| Storage cost | None | Yes |
| Use case | Simplicity, row security | Expensive aggregations, reports |
| DML support | Simple views only | No (read-only) |

---

## Listing and Managing Views

```sql
-- List views
\dv
SELECT table_name, view_definition
FROM information_schema.views
WHERE table_schema = 'public';

-- List materialized views
\dm
SELECT matviewname, definition, ispopulated
FROM pg_matviews
WHERE schemaname = 'public';

-- Drop
DROP VIEW active_users;
DROP VIEW IF EXISTS active_users CASCADE;   -- also drops dependent views

DROP MATERIALIZED VIEW monthly_revenue;
DROP MATERIALIZED VIEW IF EXISTS monthly_revenue;

-- Rename
ALTER VIEW active_users RENAME TO active_customers;
```
# Chapter 10: Full-Text Search

## What is Full-Text Search?

Full-text search (FTS) finds documents that match natural language queries, with ranking by relevance. It's smarter than `LIKE` — it understands language:
- **Stemming:** "running" matches "run", "runs", "runner"
- **Stop words:** ignores "the", "a", "is"
- **Ranking:** most relevant results first
- **Phrase search:** exact phrase matching
- **Dictionaries:** language-specific rules

---

## Core Types

```sql
-- tsvector: preprocessed document representation
SELECT to_tsvector('english', 'The quick brown fox jumps over the lazy dog');
-- 'brown':3 'dog':9 'fox':4 'jump':5 'lazi':8 'quick':2
-- (stop words "the", "over" removed; "jumps"→"jump", "lazy"→"lazi")

-- tsquery: search query
SELECT to_tsquery('english', 'quick & fox');         -- AND
SELECT to_tsquery('english', 'quick | fox');         -- OR
SELECT to_tsquery('english', 'quick & !cat');        -- NOT
SELECT plainto_tsquery('english', 'quick brown fox'); -- converts plain text
SELECT phraseto_tsquery('english', 'quick brown fox'); -- phrase search
SELECT websearch_to_tsquery('english', '"quick brown" -cat'); -- web-style
```

---

## Basic Full-Text Search

```sql
-- Simple search without index
SELECT title, body
FROM posts
WHERE to_tsvector('english', title || ' ' || body) @@ to_tsquery('english', 'postgresql');

-- Using plainto_tsquery (easier for user input)
SELECT title
FROM posts
WHERE to_tsvector('english', title || ' ' || body) @@ plainto_tsquery('english', 'learn postgresql tutorial');

-- Ranking results
SELECT
    title,
    ts_rank(to_tsvector('english', title || ' ' || body), query) AS rank
FROM posts,
     to_tsquery('english', 'postgresql') query
WHERE to_tsvector('english', title || ' ' || body) @@ query
ORDER BY rank DESC
LIMIT 10;
```

---

## Adding a tsvector Column (Recommended)

For performance, store the computed tsvector in a column:

```sql
-- Add tsvector column
ALTER TABLE posts ADD COLUMN search_vector TSVECTOR;

-- Populate it
UPDATE posts SET search_vector =
    setweight(to_tsvector('english', COALESCE(title, '')), 'A') ||
    setweight(to_tsvector('english', COALESCE(body, '')), 'B') ||
    setweight(to_tsvector('english', COALESCE(tags::TEXT, '')), 'C');
-- setweight assigns relevance: A > B > C > D
-- Matches in title rank higher than body

-- GIN index on the column
CREATE INDEX idx_posts_search ON posts USING GIN(search_vector);

-- Trigger to keep it updated automatically
CREATE OR REPLACE FUNCTION posts_search_trigger()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.body, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER posts_search_update
BEFORE INSERT OR UPDATE ON posts
FOR EACH ROW EXECUTE FUNCTION posts_search_trigger();
```

---

## Searching with the Index

```sql
-- Fast GIN index search
SELECT id, title,
    ts_rank(search_vector, query) AS rank
FROM posts,
     plainto_tsquery('english', 'learn postgresql') query
WHERE search_vector @@ query
ORDER BY rank DESC
LIMIT 20;

-- Web-style search (supports quotes and minus)
SELECT id, title, ts_rank(search_vector, query) AS rank
FROM posts,
     websearch_to_tsquery('english', '"full text" -mysql') query
WHERE search_vector @@ query
ORDER BY rank DESC;
```

---

## Ranking Functions

```sql
-- ts_rank: ranks based on term frequency
SELECT ts_rank(search_vector, query) AS rank

-- ts_rank_cd: ranks using cover density (considers proximity of terms)
SELECT ts_rank_cd(search_vector, query) AS rank

-- Normalize rank by document length (prevents long docs from always winning)
SELECT ts_rank(search_vector, query, 1) AS rank  -- 1 = divide by 1 + log(length)
SELECT ts_rank(search_vector, query, 32) AS rank -- 32 = divide by ndoc (# documents)
```

---

## Highlighting — ts_headline

```sql
SELECT
    id,
    title,
    ts_headline(
        'english',
        body,
        plainto_tsquery('english', 'postgresql index'),
        'MaxFragments=3, MaxWords=15, MinWords=5, StartSel=<b>, StopSel=</b>'
    ) AS excerpt
FROM posts
WHERE search_vector @@ plainto_tsquery('english', 'postgresql index')
ORDER BY ts_rank(search_vector, plainto_tsquery('english', 'postgresql index')) DESC;

-- Result: "...using a <b>PostgreSQL</b> <b>index</b> for full-text..."
```

---

## Phrase Search

```sql
-- Exact phrase search
SELECT * FROM posts
WHERE search_vector @@ phraseto_tsquery('english', 'getting started with postgresql');

-- Proximity search — within N positions of each other
SELECT * FROM posts
WHERE search_vector @@ to_tsquery('english', 'getting <-> started');  -- adjacent
WHERE search_vector @@ to_tsquery('english', 'getting <3> postgresql'); -- within 3 words
```

---

## Dictionaries and Languages

```sql
-- List available dictionaries
SELECT dictname FROM pg_ts_dict;

-- Available languages
SELECT cfgname FROM pg_ts_config;

-- Use different languages
SELECT to_tsvector('spanish', 'Los perros corren en el parque');
SELECT to_tsvector('german', 'Die Hunde laufen im Park');
SELECT to_tsvector('french', 'Les chiens courent dans le parc');

-- Simple dictionary (no stemming, just lowercase)
SELECT to_tsvector('simple', 'The Quick Brown Fox');
-- 'brown':3 'fox':4 'quick':2 'the':1  (no stemming, no stop words)
```

---

## Multi-language Search

```sql
CREATE TABLE articles (
    id BIGSERIAL PRIMARY KEY,
    title TEXT,
    body TEXT,
    language TEXT DEFAULT 'english',
    search_vector TSVECTOR
);

-- Trigger with dynamic language
CREATE OR REPLACE FUNCTION articles_search_trigger()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector(NEW.language::regconfig, COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector(NEW.language::regconfig, COALESCE(NEW.body, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Search in specific language
SELECT * FROM articles
WHERE search_vector @@ plainto_tsquery('spanish', 'perros')
AND language = 'spanish';
```

---

## Full-Text Search Performance Tips

```sql
-- Always use GIN index for full-text search (not B-tree, not GiST unless using phrase)
CREATE INDEX idx_fts ON posts USING GIN(search_vector);

-- Use GiST if you need phrase search (slower build, supports <->)
CREATE INDEX idx_fts_gist ON posts USING GIST(search_vector);

-- Check if index is being used
EXPLAIN ANALYZE
SELECT * FROM posts
WHERE search_vector @@ plainto_tsquery('english', 'postgresql')
LIMIT 10;
-- Should show: Bitmap Index Scan on idx_fts

-- Use LIMIT to avoid ranking millions of rows
SELECT id, title, ts_rank(search_vector, q) AS rank
FROM posts, plainto_tsquery('english', 'postgresql') q
WHERE search_vector @@ q
ORDER BY rank DESC
LIMIT 20;   -- rank only what you need
```

---

## When to Use PostgreSQL FTS vs Elasticsearch

| Feature | PostgreSQL FTS | Elasticsearch |
|---|---|---|
| Setup | None (built-in) | Separate cluster |
| Sync | Automatic (trigger) | ETL/CDC required |
| ACID | Yes | No |
| Performance | Good for millions | Better for billions |
| Aggregations | Basic | Advanced |
| Fuzzy search | Limited | Excellent |
| Relevance tuning | Limited | Extensive |
| Cost | Free (already running) | Additional infra |

**Use PostgreSQL FTS when:** you have < 10M documents, already use PostgreSQL, need transactional consistency.

**Use Elasticsearch when:** you have billions of documents, need advanced relevance tuning, or need complex aggregation analytics.
# Chapter 11: JSON & JSONB

## JSON vs JSONB

| Feature | JSON | JSONB |
|---|---|---|
| Storage | Text (preserves whitespace, key order) | Binary (no whitespace, sorted keys) |
| Write speed | Faster (no parsing) | Slightly slower |
| Read speed | Slower (re-parsed each time) | Faster |
| Indexing | Not indexable | GIN/btree indexable |
| Operators | Basic | Full set |
| Recommendation | Rarely use | **Always use JSONB** |

---

## Creating JSONB Columns

```sql
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    attributes JSONB,
    metadata JSONB DEFAULT '{}'::JSONB
);

INSERT INTO products (name, attributes) VALUES
('Laptop', '{"brand": "Dell", "cpu": "Intel i7", "ram_gb": 16, "colors": ["black", "silver"]}'),
('Phone', '{"brand": "Apple", "model": "iPhone 15", "storage_gb": 256, "5g": true}'),
('Headphones', '{"brand": "Sony", "wireless": true, "noise_cancelling": true, "colors": ["black"]}');
```

---

## Operators

```sql
-- -> returns JSON/JSONB element (keeps as JSON)
SELECT attributes -> 'brand' FROM products;           -- "Dell"
SELECT attributes -> 'colors' -> 0 FROM products;     -- "black"

-- ->> returns text (unwraps the value)
SELECT attributes ->> 'brand' FROM products;          -- Dell (no quotes)
SELECT attributes ->> 'ram_gb' FROM products;         -- 16 (as text)

-- #> path operator (returns JSON)
SELECT attributes #> '{colors, 0}' FROM products;    -- "black"

-- #>> path operator (returns text)
SELECT attributes #>> '{colors, 0}' FROM products;   -- black

-- @> contains (JSONB only)
SELECT * FROM products WHERE attributes @> '{"brand": "Apple"}';
SELECT * FROM products WHERE attributes @> '{"wireless": true}';

-- <@ is contained by
SELECT * FROM products WHERE '{"brand": "Sony"}' <@ attributes;

-- ? key exists
SELECT * FROM products WHERE attributes ? 'noise_cancelling';

-- ?| any of these keys exist
SELECT * FROM products WHERE attributes ?| ARRAY['wireless', '5g'];

-- ?& all of these keys exist
SELECT * FROM products WHERE attributes ?& ARRAY['brand', 'model'];
```

---

## JSONB Functions

```sql
-- jsonb_object_keys — list all keys
SELECT jsonb_object_keys(attributes) FROM products WHERE id = 1;

-- jsonb_each — expand to rows of (key, value)
SELECT key, value FROM products, jsonb_each(attributes) WHERE id = 1;

-- jsonb_each_text — expand to rows of (key, text_value)
SELECT key, value FROM products, jsonb_each_text(attributes) WHERE id = 1;

-- jsonb_extract_path / jsonb_extract_path_text
SELECT jsonb_extract_path(attributes, 'brand') FROM products;
SELECT jsonb_extract_path_text(attributes, 'brand') FROM products;

-- jsonb_set — modify a value
UPDATE products
SET attributes = jsonb_set(attributes, '{brand}', '"DELL"')
WHERE id = 1;

-- jsonb_set creating nested paths
UPDATE products
SET attributes = jsonb_set(attributes, '{specs, ssd}', 'true', true)
WHERE id = 1;
-- true parameter = create path if it doesn't exist

-- jsonb_insert — insert into array
UPDATE products
SET attributes = jsonb_insert(attributes, '{colors, -1}', '"gold"')
WHERE id = 1;
-- -1 position = append to end

-- Merge JSONB (|| operator)
UPDATE products
SET attributes = attributes || '{"weight_kg": 1.5, "warranty_years": 2}'::JSONB
WHERE id = 1;

-- Remove key (- operator)
UPDATE products
SET attributes = attributes - 'old_field'
WHERE id = 1;

-- Remove nested key (#- operator)
UPDATE products
SET attributes = attributes #- '{specs, old_spec}'
WHERE id = 1;

-- jsonb_array_elements — expand JSON array to rows
SELECT value FROM products, jsonb_array_elements(attributes -> 'colors') WHERE id = 1;

-- jsonb_array_length
SELECT jsonb_array_length(attributes -> 'colors') FROM products WHERE id = 1;

-- json_agg — aggregate rows into JSON array
SELECT json_agg(row_to_json(u)) FROM users u WHERE is_active = TRUE;

-- jsonb_agg
SELECT jsonb_agg(DISTINCT role) FROM users;

-- json_build_object
SELECT json_build_object('id', id, 'name', name, 'email', email) FROM users;

-- row_to_json
SELECT row_to_json(u) FROM users u WHERE id = 1;
```

---

## Indexing JSONB

```sql
-- GIN index — supports @>, ?, ?|, ?& operators
CREATE INDEX idx_products_attrs ON products USING GIN(attributes);

-- Expression index — for specific field queries
CREATE INDEX idx_products_brand ON products((attributes->>'brand'));

-- GIN with jsonb_path_ops (smaller, only supports @>)
CREATE INDEX idx_products_attrs_path ON products USING GIN(attributes jsonb_path_ops);
```

---

## JSONB Path Queries (PostgreSQL 12+)

```sql
-- jsonb_path_query — XPath-like queries
SELECT jsonb_path_query(attributes, '$.colors[*]') FROM products WHERE id = 1;

-- Find products where any color is "black"
SELECT * FROM products
WHERE jsonb_path_exists(attributes, '$.colors[*] ? (@ == "black")');

-- Find products with RAM >= 16
SELECT * FROM products
WHERE jsonb_path_exists(attributes, '$.ram_gb ? (@ >= 16)');

-- @? operator (path exists)
SELECT * FROM products WHERE attributes @? '$.colors[*] ? (@ == "black")';

-- @@ operator (path condition is true)
SELECT * FROM products WHERE attributes @@ '$.ram_gb >= 16';
```

---

## Practical Patterns

### EAV Alternative (Entity-Attribute-Value):
```sql
-- Instead of hundreds of nullable columns OR EAV tables:
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    preferences JSONB DEFAULT '{}',
    settings JSONB DEFAULT '{}'
);

-- Flexible schema per user
UPDATE users SET preferences = preferences || '{"theme": "dark", "language": "en"}' WHERE id = 1;
UPDATE users SET settings = settings || '{"notifications": {"email": true, "sms": false}}' WHERE id = 1;

SELECT preferences->>'theme' FROM users WHERE id = 1;
SELECT * FROM users WHERE preferences @> '{"theme": "dark"}';
```

### Storing Dynamic Forms:
```sql
CREATE TABLE form_submissions (
    id BIGSERIAL PRIMARY KEY,
    form_id BIGINT REFERENCES forms(id),
    user_id BIGINT REFERENCES users(id),
    data JSONB NOT NULL,
    submitted_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO form_submissions (form_id, user_id, data) VALUES
(1, 42, '{"name": "Alice", "age": 30, "interests": ["coding", "hiking"]}');

SELECT data->>'name', data->>'age'
FROM form_submissions
WHERE form_id = 1 AND data @> '{"interests": ["coding"]}';
```

---

## JSONB vs Relational — When to Use Which

**Use JSONB when:**
- Schema varies per row (product attributes, form responses)
- Data is accessed as a blob, not queried field-by-field
- You're prototyping and the schema is still evolving
- Storing external API responses as-is

**Use relational tables when:**
- You filter/sort/join on specific fields frequently
- You need strong integrity constraints (foreign keys, check constraints)
- Multiple rows reference the same data
- You need efficient updates to specific fields
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
# Chapter 13: Replication & High Availability

## Why Replication?

- **High availability** — automatic failover when primary fails
- **Read scaling** — distribute read queries to replicas
- **Disaster recovery** — geographic replication
- **Zero-downtime upgrades** — failover to replica, upgrade primary

---

## Streaming Replication (Physical)

Streaming replication sends WAL (Write-Ahead Log) data from the primary to standbys in near-real-time. Standbys are byte-for-byte copies of the primary.

### Primary Configuration (`postgresql.conf`):
```ini
# Replication settings
wal_level = replica          # or 'logical' for logical replication
max_wal_senders = 10         # max number of replication connections
wal_keep_size = 1GB          # keep WAL for replicas that fall behind
hot_standby = on             # allow read queries on standby

# Performance
synchronous_commit = on      # 'local', 'remote_write', 'remote_apply', or 'off'
```

### `pg_hba.conf` on Primary:
```
# Allow replication connections from standby
host replication replicator standby-ip/32 scram-sha-256
```

### Create Replication User:
```sql
CREATE ROLE replicator WITH REPLICATION LOGIN PASSWORD 'replpassword';
```

### Create Replication Slot (optional but recommended):
```sql
-- Prevents primary from removing WAL files the replica still needs
CREATE REPLICATION SLOT replica1 LOGICAL;    -- for logical replication
CREATE REPLICATION SLOT replica1 PHYSICAL;   -- for physical replication
SELECT slot_name, slot_type, active FROM pg_replication_slots;
```

### Standby Setup:
```bash
# Create base backup from primary
pg_basebackup -h primary-host -U replicator -D /var/lib/postgresql/16/main \
    -P -Xs -R -C -S replica1

# -P: show progress
# -Xs: use streaming
# -R: create standby.signal and recovery.conf
# -C: create replication slot
# -S replica1: slot name

# standby.signal is created automatically by pg_basebackup -R
# postgresql.auto.conf is created with primary_conninfo
```

### `postgresql.auto.conf` on Standby:
```ini
primary_conninfo = 'host=primary-host port=5432 user=replicator password=replpassword'
primary_slot_name = 'replica1'
```

---

## Monitoring Replication

```sql
-- On primary: view connected standbys
SELECT
    client_addr,
    state,
    sent_lsn,
    write_lsn,
    flush_lsn,
    replay_lsn,
    (sent_lsn - replay_lsn) AS replication_lag_bytes,
    sync_state
FROM pg_stat_replication;

-- On standby: view replication status
SELECT
    now() - pg_last_xact_replay_timestamp() AS replication_delay,
    pg_is_in_recovery() AS is_standby,
    pg_last_wal_receive_lsn() AS received_lsn,
    pg_last_wal_replay_lsn() AS replayed_lsn;
```

---

## Hot Standby — Read-Only Queries on Replica

```ini
# postgresql.conf on standby
hot_standby = on
hot_standby_feedback = on   # prevents primary from vacuuming rows needed by replica queries
```

```bash
# Connect to standby (read-only)
psql -h standby-host -U myuser -d mydb

-- Can run SELECT queries
SELECT * FROM users LIMIT 10;

-- DML is rejected
INSERT INTO users (name) VALUES ('test');
-- ERROR: cannot execute INSERT in a read-only transaction
```

---

## Failover

If the primary fails, promote a standby to primary:

```bash
# Manual failover — promote standby to primary
pg_ctl promote -D /var/lib/postgresql/16/main

# Or in psql
SELECT pg_promote();

# The standby.signal file is removed
# The standby begins accepting writes
```

After failover, update application connection strings to point to the new primary.

---

## Logical Replication

Logical replication replicates at the SQL level (individual rows/operations), not the physical WAL byte stream. This allows:
- Replication between different PostgreSQL versions
- Selective table replication
- Replication to different schemas
- Bidirectional replication (with care)

### Publisher (Source):
```sql
-- Enable logical replication
-- wal_level = logical (in postgresql.conf)

-- Create publication
CREATE PUBLICATION my_pub FOR TABLE users, orders;
CREATE PUBLICATION all_tables FOR ALL TABLES;

-- View publications
SELECT * FROM pg_publication;
SELECT * FROM pg_publication_tables;
```

### Subscriber (Target):
```sql
-- Create subscription
CREATE SUBSCRIPTION my_sub
    CONNECTION 'host=primary-host dbname=mydb user=replicator password=secret'
    PUBLICATION my_pub;

-- View subscriptions
SELECT subname, subenabled, subpublications FROM pg_subscription;

-- Subscription status
SELECT * FROM pg_stat_subscription;
```

---

## Synchronous Replication

By default, replication is asynchronous — the primary commits immediately, replication follows. For stronger durability:

```ini
# postgresql.conf on primary
synchronous_commit = remote_apply      # wait for replica to apply the WAL
synchronous_standby_names = 'replica1' # which standbys count
```

| synchronous_commit value | Durability | Latency |
|---|---|---|
| `off` | WAL not even written immediately | Lowest |
| `local` | WAL written to local disk | Low |
| `remote_write` | WAL sent to standby OS buffer | Medium |
| `remote_apply` | WAL applied on standby | Highest |

---

## Patroni — Automated Failover

Patroni is the industry standard for automated PostgreSQL HA:

```yaml
# patroni.yml
scope: my-postgres-cluster
name: node1

restapi:
  listen: 0.0.0.0:8008
  connect_address: 192.168.1.1:8008

etcd:
  hosts: etcd1:2379,etcd2:2379,etcd3:2379

bootstrap:
  dcs:
    ttl: 30
    loop_wait: 10
    retry_timeout: 10
    maximum_lag_on_failover: 1048576

postgresql:
  listen: 0.0.0.0:5432
  connect_address: 192.168.1.1:5432
  data_dir: /var/lib/postgresql/16/main
  parameters:
    wal_level: replica
    hot_standby: on
    max_wal_senders: 10
```

```bash
# Check cluster status
patronictl -c patroni.yml list

# Manual switchover (graceful)
patronictl -c patroni.yml switchover my-postgres-cluster

# Manual failover (emergency)
patronictl -c patroni.yml failover my-postgres-cluster
```

---

## Connection Routing with pgBouncer + Patroni

```
Applications
    │
    ├── Write connections → pgBouncer (primary) → Patroni primary
    │
    └── Read connections → pgBouncer (replica) → Patroni replica

# pgBouncer auto-detects primary via Patroni REST API
```

---

## Read Replica Routing in Application

```yaml
# Spring Boot
spring:
  datasource:
    write:
      url: jdbc:postgresql://primary:5432/mydb
    read:
      url: jdbc:postgresql://replica:5432/mydb
```

```java
@Service
public class UserService {

    @Transactional(readOnly = true)  // route to read replica
    public List<UserDto> getUsers() { ... }

    @Transactional                   // route to primary
    public UserDto createUser(CreateUserRequest req) { ... }
}
```
# Chapter 14: Security & User Management

## Roles and Users

In PostgreSQL, **roles** are the unified concept for both users and groups. A role with `LOGIN` privilege is effectively a "user".

```sql
-- Create a login role (user)
CREATE ROLE alice WITH LOGIN PASSWORD 'securepassword';
CREATE USER alice WITH PASSWORD 'securepassword';  -- same as above

-- Create a group role (no login)
CREATE ROLE developers;
CREATE ROLE readonly_users;

-- Add users to groups
GRANT developers TO alice;
GRANT readonly_users TO alice, bob;

-- Login as alice — inherits all privileges from developers role
```

---

## Role Attributes

```sql
CREATE ROLE myapp_admin WITH
    LOGIN                               -- can log in
    PASSWORD 'secure_password'
    SUPERUSER                           -- all privileges (dangerous!)
    CREATEDB                            -- can create databases
    CREATEROLE                          -- can create roles
    REPLICATION                         -- can initiate replication
    BYPASSRLS                           -- bypasses Row Level Security
    CONNECTION LIMIT 20                 -- max simultaneous connections
    VALID UNTIL '2025-01-01';           -- account expiry

-- Modify existing role
ALTER ROLE alice PASSWORD 'newpassword';
ALTER ROLE alice VALID UNTIL 'infinity';
ALTER ROLE alice CONNECTION LIMIT 50;
ALTER ROLE alice NOLOGIN;              -- disable login

-- Drop role
DROP ROLE alice;
DROP ROLE IF EXISTS alice;
```

---

## Privileges

```sql
-- Database privileges
GRANT CONNECT ON DATABASE mydb TO myapp_user;
GRANT CREATE ON DATABASE mydb TO developer;

-- Schema privileges
GRANT USAGE ON SCHEMA public TO myapp_user;
GRANT CREATE ON SCHEMA public TO developer;

-- Table privileges
GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly_users;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO myapp_user;
GRANT SELECT ON users TO reporting_user;

-- Specific columns
GRANT SELECT (id, name, email) ON users TO limited_user;  -- no password column!

-- Sequences (needed for INSERT with SERIAL/BIGSERIAL columns)
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO myapp_user;

-- Functions
GRANT EXECUTE ON FUNCTION calculate_tax(NUMERIC, NUMERIC) TO myapp_user;

-- Future tables (apply to new tables automatically)
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT ON TABLES TO readonly_users;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO myapp_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT USAGE, SELECT ON SEQUENCES TO myapp_user;
```

---

## Revoking Privileges

```sql
REVOKE ALL ON ALL TABLES IN SCHEMA public FROM old_user;
REVOKE DELETE ON orders FROM restricted_user;
REVOKE CREATE ON DATABASE mydb FROM developer;

-- Revoke from public (default: everyone has CONNECT)
REVOKE CONNECT ON DATABASE mydb FROM PUBLIC;
```

---

## Principle of Least Privilege — App Role Pattern

```sql
-- Create minimal-privilege application role
CREATE ROLE app_user WITH LOGIN PASSWORD 'app_password' CONNECTION LIMIT 50;

-- Grant only what the app needs
GRANT CONNECT ON DATABASE mydb TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;

-- Separate read-only role for reporting
CREATE ROLE reporting_user WITH LOGIN PASSWORD 'report_password';
GRANT CONNECT ON DATABASE mydb TO reporting_user;
GRANT USAGE ON SCHEMA public TO reporting_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO reporting_user;

-- Do NOT grant these to app_user:
-- SUPERUSER, CREATEDB, CREATEROLE
-- DROP TABLE, ALTER TABLE (handled by migrations with separate migration user)
```

---

## Row-Level Security (RLS)

RLS restricts which rows each user can see or modify — enforced at the database level.

```sql
-- Enable RLS on a table
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE orders FORCE ROW LEVEL SECURITY;  -- applies to table owner too

-- Policy: users can only see their own orders
CREATE POLICY orders_user_isolation ON orders
    FOR ALL
    TO app_user
    USING (user_id = current_setting('app.current_user_id')::BIGINT);

-- Separate read and write policies
CREATE POLICY orders_select ON orders
    FOR SELECT
    TO app_user
    USING (user_id = current_setting('app.current_user_id')::BIGINT);

CREATE POLICY orders_insert ON orders
    FOR INSERT
    TO app_user
    WITH CHECK (user_id = current_setting('app.current_user_id')::BIGINT);

-- Admin can see all orders
CREATE POLICY orders_admin ON orders
    FOR ALL
    TO admin_user
    USING (TRUE);
```

### Setting Context in Application:

```sql
-- Application sets user context before queries
SET LOCAL app.current_user_id = '42';
SELECT * FROM orders;  -- only sees orders for user 42

-- In a transaction
BEGIN;
    SET LOCAL app.current_user_id = '42';
    SELECT * FROM orders;  -- user 42's orders only
    INSERT INTO orders (amount) VALUES (100);  -- auto-sets user_id = 42
COMMIT;
```

### View RLS Policies:
```sql
SELECT tablename, policyname, cmd, roles, qual, with_check
FROM pg_policies
WHERE tablename = 'orders';
```

---

## Multi-Tenant RLS Pattern

```sql
CREATE TABLE tenants (id BIGSERIAL PRIMARY KEY, name TEXT);
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    name TEXT NOT NULL,
    price NUMERIC(10,2)
);

ALTER TABLE products ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON products
    USING (tenant_id = current_setting('app.tenant_id')::BIGINT);

-- Application sets tenant context
SET app.tenant_id = '100';
SELECT * FROM products;   -- only tenant 100's products
```

---

## pg_hba.conf — Client Authentication

```
# TYPE  DATABASE    USER         ADDRESS          METHOD
local   all         postgres                      peer         # local unix socket, postgres OS user
local   all         all                           md5          # local unix socket, password
host    mydb        app_user     10.0.0.0/8       scram-sha-256 # network, password
hostssl mydb        app_user     0.0.0.0/0        scram-sha-256 # SSL required
host    all         all          127.0.0.1/32     trust        # localhost, no password (dev only!)
reject  all         all          0.0.0.0/0                     # deny all others
```

Authentication methods:
- `trust` — no password (dev/local only!)
- `md5` — MD5-hashed password (deprecated, use scram)
- `scram-sha-256` — SCRAM authentication (current best practice)
- `peer` — OS username must match DB username (Unix sockets)
- `cert` — SSL certificate authentication
- `reject` — always deny

---

## SSL/TLS

```ini
# postgresql.conf
ssl = on
ssl_cert_file = 'server.crt'
ssl_key_file = 'server.key'
ssl_ca_file = 'root.crt'         # for client certificate verification
```

```
# Require SSL in pg_hba.conf
hostssl mydb app_user 0.0.0.0/0 scram-sha-256

# Connection string requiring SSL
postgresql://user:pass@host/db?sslmode=require
postgresql://user:pass@host/db?sslmode=verify-full&sslrootcert=/path/to/root.crt
```

---

## Auditing

```sql
-- Track all connections
log_connections = on
log_disconnections = on
log_hostname = on

-- Track slow queries
log_min_duration_statement = 1000   # log queries taking > 1 second

-- Track DDL changes
log_statement = 'ddl'   # 'none', 'ddl', 'mod', 'all'

-- pgaudit extension — detailed audit logging
CREATE EXTENSION pgaudit;
-- pgaudit.log = 'write, ddl, role'
```

---

## Secrets Management

```bash
# Never hardcode passwords in connection strings committed to git
# Use environment variables or a secrets manager

# .env file (never commit!)
DATABASE_URL=postgresql://app_user:$(vault read -field=password secret/db)@host/mydb

# Kubernetes secret
kubectl create secret generic db-secret \
    --from-literal=password=$(vault read -field=password secret/db)

# AWS Secrets Manager in application
aws secretsmanager get-secret-value --secret-id prod/myapp/db
```
# Chapter 15: Backup & Recovery

## Backup Strategies

| Method | Tool | Granularity | Restore Speed | Notes |
|---|---|---|---|---|
| Logical dump | `pg_dump` | Table level | Slow (reimport SQL) | Portable, cross-version |
| Physical backup | `pg_basebackup` | Cluster | Fast | Same PG version only |
| PITR | WAL archiving | Point-in-time | Medium | Recover to any moment |
| Continuous | pgBackRest / Barman | Incremental | Fast | Production standard |

---

## pg_dump — Logical Backups

`pg_dump` exports a database as SQL (or custom format) that can be imported anywhere.

```bash
# Plain SQL format
pg_dump -U postgres -d mydb -f mydb_backup.sql

# Custom format (compressed, faster restore, selective restore possible)
pg_dump -U postgres -d mydb -Fc -f mydb_backup.dump

# Directory format (parallel dump with -j)
pg_dump -U postgres -d mydb -Fd -j 4 -f mydb_backup/
# -j 4 = use 4 parallel workers

# With connection string
pg_dump "postgresql://user:pass@host:5432/mydb" -Fc -f backup.dump

# Exclude data (schema only)
pg_dump -U postgres -d mydb --schema-only -f schema_only.sql

# Data only (no schema)
pg_dump -U postgres -d mydb --data-only -f data_only.sql

# Specific tables only
pg_dump -U postgres -d mydb -t users -t orders -Fc -f partial.dump

# Exclude specific tables
pg_dump -U postgres -d mydb --exclude-table=logs --exclude-table=audit -Fc -f backup.dump

# Compressed SQL
pg_dump -U postgres -d mydb | gzip > backup_$(date +%Y%m%d).sql.gz
```

---

## pg_dumpall — Dump Entire Cluster

```bash
# Dump all databases + roles + tablespaces
pg_dumpall -U postgres -f full_cluster.sql

# Roles only (no data)
pg_dumpall -U postgres --roles-only -f roles.sql

# Globals only (roles + tablespaces, no data)
pg_dumpall -U postgres --globals-only -f globals.sql
```

---

## pg_restore — Restore from Custom Format

```bash
# Restore to existing (empty) database
pg_restore -U postgres -d mydb mydb_backup.dump

# Create database and restore
pg_restore -U postgres -d postgres --create mydb_backup.dump

# Parallel restore
pg_restore -U postgres -d mydb -j 4 mydb_backup.dump

# Restore specific tables only
pg_restore -U postgres -d mydb -t users -t orders mydb_backup.dump

# List contents of dump file
pg_restore --list mydb_backup.dump

# Restore with verbose output
pg_restore -U postgres -d mydb -v mydb_backup.dump

# From SQL format
psql -U postgres -d mydb -f mydb_backup.sql
# Or
psql "postgresql://user:pass@host/mydb" < mydb_backup.sql
```

---

## pg_basebackup — Physical Backup

Physical backup copies the entire data directory — same as what streaming replication uses.

```bash
# Basic physical backup
pg_basebackup -U replicator -h localhost -D /backup/pg_base -P -Xs

# With compression
pg_basebackup -U replicator -h localhost -D /backup/pg_base -P -Xs -z

# Tar format (good for S3 upload)
pg_basebackup -U replicator -h localhost -Ft -z -D - | \
    aws s3 cp - s3://my-bucket/pg-backup-$(date +%Y%m%d).tar.gz

# Options:
# -P: show progress
# -Xs: stream WAL during backup
# -z: gzip compress
# -R: write recovery config (for replica setup)
```

---

## WAL Archiving & Point-in-Time Recovery (PITR)

PITR lets you recover to ANY point in time, not just when a backup was taken.

### Enable WAL Archiving (`postgresql.conf`):
```ini
wal_level = replica
archive_mode = on
archive_command = 'cp %p /archive/%f'
# Or to S3:
archive_command = 'aws s3 cp %p s3://my-wal-bucket/%f'
archive_timeout = 300   # archive WAL every 5 minutes even if not full
```

### Verify Archive:
```sql
SELECT pg_walfile_name(pg_current_wal_lsn());  -- current WAL file name
SELECT * FROM pg_stat_archiver;                 -- archiver status
```

### PITR Recovery:

1. Restore a base backup
2. Create recovery configuration

```ini
# postgresql.conf (or recovery.conf for PG < 12)
restore_command = 'cp /archive/%f %p'
# Or from S3:
restore_command = 'aws s3 cp s3://my-wal-bucket/%f %p'

# Recover to specific time
recovery_target_time = '2024-02-15 14:30:00'

# Recover to specific transaction ID
recovery_target_xid = '12345678'

# Recover to specific LSN
recovery_target_lsn = '0/15000000'

# Allow connections after recovery
recovery_target_action = 'promote'   # or 'pause', 'shutdown'
```

```bash
# Create recovery signal file
touch /var/lib/postgresql/16/main/recovery.signal

# Start PostgreSQL — it replays WAL until target time
sudo systemctl start postgresql
```

---

## pgBackRest — Production Backup Tool

pgBackRest is the recommended tool for production. It handles incremental backups, parallel transfers, S3 storage, and more.

```ini
# /etc/pgbackrest/pgbackrest.conf
[global]
repo1-type=s3
repo1-s3-bucket=my-postgres-backups
repo1-s3-region=us-east-1
repo1-s3-key=AKIAIOSFODNN7EXAMPLE
repo1-s3-key-secret=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
repo1-retention-full=7        # keep 7 full backups
repo1-retention-diff=30       # keep 30 differential backups
compress-type=lz4
process-max=4

[mydb]
pg1-path=/var/lib/postgresql/16/main
pg1-user=postgres
```

```bash
# Initialize repository
pgbackrest --stanza=mydb stanza-create

# Full backup
pgbackrest --stanza=mydb backup --type=full

# Incremental backup (default)
pgbackrest --stanza=mydb backup

# Check backup info
pgbackrest --stanza=mydb info

# Restore latest backup
pgbackrest --stanza=mydb restore

# PITR restore
pgbackrest --stanza=mydb restore \
    --target='2024-02-15 14:30:00' \
    --target-action=promote

# Verify backups
pgbackrest --stanza=mydb check
```

---

## Automated Backup Schedule

### With cron:
```bash
# /etc/cron.d/postgresql-backup
# Full backup every Sunday at 1 AM
0 1 * * 0 postgres pgbackrest --stanza=mydb backup --type=full

# Incremental backup daily at 1 AM (Mon-Sat)
0 1 * * 1-6 postgres pgbackrest --stanza=mydb backup

# Verify backups
0 3 * * * postgres pgbackrest --stanza=mydb check
```

### pg_dump rotation script:
```bash
#!/bin/bash
BACKUP_DIR="/backups/postgresql"
DB_NAME="mydb"
RETENTION_DAYS=30

# Create backup
pg_dump -U postgres -Fc -d $DB_NAME \
    -f "$BACKUP_DIR/${DB_NAME}_$(date +%Y%m%d_%H%M%S).dump"

# Delete old backups
find $BACKUP_DIR -name "*.dump" -mtime +$RETENTION_DAYS -delete

echo "Backup complete. Files in $BACKUP_DIR:"
ls -lh $BACKUP_DIR | tail -5
```

---

## Testing Your Backups

A backup you've never tested is not a real backup.

```bash
# Monthly restore test procedure
# 1. Spin up test PostgreSQL instance
# 2. Restore backup
pg_restore -U postgres -d test_restore mydb_backup.dump

# 3. Verify data
psql -U postgres -d test_restore -c "SELECT COUNT(*) FROM users;"
psql -U postgres -d test_restore -c "SELECT MAX(created_at) FROM orders;"

# 4. Run application smoke test against restored DB

# 5. Verify pgBackRest backup integrity
pgbackrest --stanza=mydb check
```

---

## RTO and RPO

| Term | Definition | Strategy |
|---|---|---|
| **RTO** (Recovery Time Objective) | Max acceptable downtime | Standby servers, fast restore |
| **RPO** (Recovery Point Objective) | Max acceptable data loss | WAL archiving frequency, sync replication |

| Setup | RTO | RPO |
|---|---|---|
| pg_dump daily | Hours-days | Up to 24 hours |
| WAL archiving | Minutes-hours | Minutes |
| Streaming replication | Seconds-minutes | Seconds |
| Synchronous replication | Seconds | Near-zero |
# Chapter 16: PostgreSQL with Docker & Kubernetes

## Docker Quick Start

```bash
# Pull and run PostgreSQL
docker run -d \
  --name postgres \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypassword \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql/data \
  postgres:16-alpine

# Connect
docker exec -it postgres psql -U myuser -d mydb

# Stop/start
docker stop postgres
docker start postgres

# View logs
docker logs postgres
docker logs -f postgres   # follow
```

---

## Docker Compose for Development

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    container_name: myapp-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: myuser
      POSTGRES_PASSWORD: mypassword
      POSTGRES_DB: mydb
      POSTGRES_INITDB_ARGS: "--locale=en_US.UTF-8"
      TZ: UTC
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init:/docker-entrypoint-initdb.d    # init scripts run on first start
      - ./postgresql.conf:/etc/postgresql/postgresql.conf  # custom config
    command: postgres -c config_file=/etc/postgresql/postgresql.conf
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U myuser -d mydb"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s

  pgadmin:
    image: dpage/pgadmin4:latest
    container_name: myapp-pgadmin
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@example.com
      PGADMIN_DEFAULT_PASSWORD: admin
    ports:
      - "5050:80"
    depends_on:
      postgres:
        condition: service_healthy

volumes:
  postgres_data:
```

```bash
docker-compose up -d
docker-compose logs -f postgres
docker-compose exec postgres psql -U myuser -d mydb
docker-compose down              # stop but keep volumes
docker-compose down -v           # stop and remove volumes
```

---

## Init Scripts

Files in `/docker-entrypoint-initdb.d/` run automatically on the first start (when data directory is empty):

```
init/
├── 01_schema.sql      # create tables
├── 02_seed.sql        # insert test data
└── 03_extensions.sql  # create extensions
```

```sql
-- 01_schema.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

```bash
-- 02_seed.sql
INSERT INTO users (email, name) VALUES
    ('alice@example.com', 'Alice'),
    ('bob@example.com', 'Bob');
```

---

## Custom postgresql.conf

```ini
# postgresql.conf (mount into container)
# Memory
shared_buffers = 256MB
work_mem = 16MB
maintenance_work_mem = 64MB

# Connections
max_connections = 100

# Logging (for development)
log_statement = 'all'
log_duration = on
log_min_duration_statement = 0

# Performance
random_page_cost = 1.1    # SSD-optimized
effective_cache_size = 1GB
```

---

## Multiple PostgreSQL Versions

```yaml
version: '3.8'
services:
  pg14:
    image: postgres:14-alpine
    ports:
      - "5433:5432"     # different host port
    environment:
      POSTGRES_PASSWORD: postgres

  pg16:
    image: postgres:16-alpine
    ports:
      - "5434:5432"
    environment:
      POSTGRES_PASSWORD: postgres
```

---

## PostgreSQL on Kubernetes

### StatefulSet (basic setup):
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
spec:
  serviceName: postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:16-alpine
          env:
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef:
                  name: postgres-secret
                  key: username
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: postgres-secret
                  key: password
            - name: POSTGRES_DB
              value: mydb
          ports:
            - containerPort: 5432
          volumeMounts:
            - name: postgres-storage
              mountPath: /var/lib/postgresql/data
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "2Gi"
              cpu: "1000m"
          livenessProbe:
            exec:
              command: ["pg_isready", "-U", "myuser"]
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            exec:
              command: ["pg_isready", "-U", "myuser"]
            initialDelaySeconds: 5
            periodSeconds: 5
  volumeClaimTemplates:
    - metadata:
        name: postgres-storage
      spec:
        accessModes: ["ReadWriteOnce"]
        storageClassName: fast-ssd
        resources:
          requests:
            storage: 100Gi
```

### Service:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: postgres
spec:
  selector:
    app: postgres
  ports:
    - port: 5432
      targetPort: 5432
  type: ClusterIP
```

### Secret:
```bash
kubectl create secret generic postgres-secret \
  --from-literal=username=myuser \
  --from-literal=password=$(openssl rand -base64 32)
```

---

## CloudNativePG Operator (Production Recommendation)

CloudNativePG is the Kubernetes operator for PostgreSQL — manages HA, failover, backups automatically.

```bash
# Install operator
kubectl apply -f https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/release-1.22/releases/cnpg-1.22.0.yaml
```

```yaml
# postgresql-cluster.yaml
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: my-postgres-cluster
spec:
  instances: 3                    # 1 primary + 2 replicas

  postgresql:
    version: "16"
    parameters:
      shared_buffers: "512MB"
      max_connections: "200"

  storage:
    size: 100Gi
    storageClass: fast-ssd

  backup:
    retentionPolicy: "30d"
    barmanObjectStore:
      destinationPath: s3://my-postgres-backups/
      s3Credentials:
        accessKeyId:
          name: aws-creds
          key: ACCESS_KEY_ID
        secretAccessKey:
          name: aws-creds
          key: SECRET_ACCESS_KEY

  monitoring:
    enablePodMonitor: true   # Prometheus scraping

  resources:
    requests:
      memory: "1Gi"
      cpu: "500m"
    limits:
      memory: "4Gi"
      cpu: "2000m"
```

```bash
kubectl apply -f postgresql-cluster.yaml

# Check status
kubectl get clusters
kubectl describe cluster my-postgres-cluster

# Connect
kubectl exec -it my-postgres-cluster-1 -- psql -U app mydb

# Switchover (graceful failover)
kubectl cnpg promote my-postgres-cluster my-postgres-cluster-2

# Take a backup
kubectl cnpg backup my-postgres-cluster
```

---

## Persistence Checklist

- Use a **StorageClass** with `Retain` reclaim policy (not `Delete`) — prevents data loss on PVC deletion
- Enable **volume snapshots** for point-in-time storage snapshots
- Use **StatefulSet** not Deployment (stable pod names and volumes)
- Configure **resource limits** to prevent OOM kills
- Set **PodDisruptionBudget** to prevent all pods being killed during node drain
- Back up to **object storage** (S3/GCS), not just local disk

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: postgres-pdb
spec:
  minAvailable: 1   # always keep at least 1 instance running
  selector:
    matchLabels:
      app: postgres
```
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
# Chapter 18: Connection Pooling

## Why Connection Pooling?

PostgreSQL creates a new OS process per connection (~5–10 MB RAM each). At high concurrency, this becomes a bottleneck.

```
Without pooling:              With pooling:
App → PostgreSQL              App → PgBouncer → PostgreSQL
1000 connections = ~5 GB RAM  1000 clients → 20 server connections
```

---

## PgBouncer — The Standard Choice

PgBouncer is a lightweight connection pooler, the most widely used for PostgreSQL.

### Install:
```bash
# Ubuntu/Debian
sudo apt install pgbouncer

# Docker
docker run -d \
  --name pgbouncer \
  -e POSTGRESQL_HOST=postgres \
  -e POSTGRESQL_PORT=5432 \
  -e POSTGRESQL_DATABASE=mydb \
  -e POSTGRESQL_USERNAME=myuser \
  -e POSTGRESQL_PASSWORD=mypassword \
  -e PGBOUNCER_POOL_MODE=transaction \
  -e PGBOUNCER_MAX_CLIENT_CONN=1000 \
  -e PGBOUNCER_DEFAULT_POOL_SIZE=20 \
  -p 5432:5432 \
  bitnami/pgbouncer:latest
```

### Configuration (`/etc/pgbouncer/pgbouncer.ini`):
```ini
[databases]
mydb = host=127.0.0.1 port=5432 dbname=mydb

[pgbouncer]
listen_port = 6432
listen_addr = *
auth_type = md5
auth_file = /etc/pgbouncer/userlist.txt

; Pool mode — the most important setting
pool_mode = transaction         ; recommended for most apps

; Pool sizes
default_pool_size = 20          ; max server connections per database/user pair
min_pool_size = 5               ; keep alive minimum connections
reserve_pool_size = 5           ; extra connections for high-demand moments
max_client_conn = 1000          ; max total clients across all pools

; Timeouts
server_idle_timeout = 600       ; close server connection idle >10 min
client_idle_timeout = 0         ; keep client connection until disconnect
server_connect_timeout = 15
server_login_retry = 15

; Logging
logfile = /var/log/pgbouncer/pgbouncer.log
pidfile = /var/run/pgbouncer/pgbouncer.pid

; Admin
admin_users = postgres
stats_users = stats_user

; TLS
;server_tls_sslmode = require
;server_tls_ca_file = /etc/ssl/certs/ca-certificates.crt
```

### User authentication file (`/etc/pgbouncer/userlist.txt`):
```
"myuser" "md5abc123..."       ; md5 hash of password
"appuser" "md5xyz456..."
```

Generate the md5 hash:
```bash
echo -n "mypasswordmyuser" | md5sum
# Prefix with "md5"
```

Or use scram-sha-256 (preferred, PgBouncer 1.18+):
```bash
# Get the hash from PostgreSQL directly
psql -c "SELECT passwd FROM pg_shadow WHERE usename='myuser';"
```

---

## Pool Modes

| Mode | Description | Use Case |
|---|---|---|
| **session** | One server connection per client session | Simple, safest — but no multiplexing |
| **transaction** | Server connection released after each transaction | **Most common** — great for most web apps |
| **statement** | Released after each statement | Rarely used — breaks multi-statement transactions |

### Session vs Transaction mode:

```
Session mode (1:1 client:server):
Client A → Server 1  (connected the whole session)
Client B → Server 2  (connected the whole session)
Client C → waiting   (pool exhausted)

Transaction mode (multiplexed):
Client A: BEGIN...  → Server 1  (holds during transaction)
Client B: SELECT    → Server 2
Client A: COMMIT    → Server 1 released
Client C: SELECT    → Server 1 (reused!)
```

**Transaction mode limitations** — these features require session mode:
- Prepared statements (use protocol-level, not `PREPARE` SQL)
- `SET` commands that persist across transactions
- Advisory locks (`pg_advisory_lock` — use `pg_advisory_xact_lock` instead)
- `LISTEN`/`NOTIFY`
- `WITH HOLD` cursors

---

## PgBouncer Commands

```bash
# Connect to admin console
psql -h 127.0.0.1 -p 6432 -U postgres pgbouncer

# Show current state
SHOW POOLS;      -- pool stats
SHOW CLIENTS;    -- connected clients
SHOW SERVERS;    -- connected server connections
SHOW STATS;      -- aggregate statistics
SHOW CONFIG;     -- current configuration

# Reload config without restart
RELOAD;

# Pause/resume (for maintenance)
PAUSE mydb;
RESUME mydb;

# Graceful shutdown
SHUTDOWN;
```

```sql
-- SHOW POOLS output explained
 database | user   | cl_active | cl_waiting | sv_active | sv_idle | sv_used | maxwait
----------+--------+-----------+------------+-----------+---------+---------+--------
 mydb     | myuser |        45 |          0 |        20 |      15 |       0 |       0

-- cl_active: clients with active server connection
-- cl_waiting: clients waiting for a connection
-- sv_active: server connections in use
-- sv_idle: server connections idle and ready
-- maxwait: max seconds client waited (should be 0)
```

---

## PgBouncer with Docker Compose

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: myuser
      POSTGRES_PASSWORD: mypassword
      POSTGRES_DB: mydb
    volumes:
      - postgres_data:/var/lib/postgresql/data
    # Don't expose port — only pgbouncer should connect

  pgbouncer:
    image: edoburu/pgbouncer:latest
    environment:
      DATABASE_URL: "postgres://myuser:mypassword@postgres:5432/mydb"
      POOL_MODE: transaction
      MAX_CLIENT_CONN: 1000
      DEFAULT_POOL_SIZE: 25
    ports:
      - "5432:5432"    # Apps connect to pgbouncer, not postgres directly
    depends_on:
      - postgres

  app:
    build: .
    environment:
      DATABASE_URL: "postgresql://myuser:mypassword@pgbouncer:5432/mydb"

volumes:
  postgres_data:
```

---

## Application Configuration with PgBouncer

### Node.js (pg / node-postgres):
```javascript
import { Pool } from 'pg';

const pool = new Pool({
  host: 'pgbouncer',
  port: 5432,
  database: 'mydb',
  user: 'myuser',
  password: 'mypassword',
  // In transaction mode, keep pool small — pgbouncer multiplexes
  max: 10,                    // max connections from this app instance
  idleTimeoutMillis: 30000,
  connectionTimeoutMillis: 2000,
});

// Use pool, not single client (let pgbouncer do its job)
const result = await pool.query('SELECT * FROM users WHERE id = $1', [userId]);
```

### Java (Spring Boot / HikariCP):
```yaml
# application.yml — connect through pgbouncer
spring:
  datasource:
    url: jdbc:postgresql://pgbouncer:5432/mydb
    username: myuser
    password: mypassword
    hikari:
      maximum-pool-size: 10          # small — pgbouncer multiplexes
      minimum-idle: 2
      connection-timeout: 3000
      idle-timeout: 600000
      max-lifetime: 1800000
      # Disable prepared statement caching in transaction mode
      data-source-properties:
        prepareThreshold: 0          # disable server-side prepared statements
```

### Python (SQLAlchemy):
```python
from sqlalchemy import create_engine

engine = create_engine(
    "postgresql+psycopg2://myuser:mypassword@pgbouncer:5432/mydb",
    pool_size=10,
    max_overflow=5,
    pool_pre_ping=True,     # verify connection before using
    # Disable prepared statements for transaction mode
    connect_args={"options": "-c statement_timeout=30000"},
    execution_options={"no_parameters": True},
)
```

---

## Monitoring PgBouncer

```bash
# View stats via admin console
psql -h 127.0.0.1 -p 6432 -U postgres pgbouncer -c "SHOW STATS;"

# Key metrics to watch:
# avg_wait_time — average time clients wait for a connection (should be <1ms)
# avg_query_time — average query time
# total_xact_count — total transactions served
```

```sql
-- Prometheus exporter or manual monitoring
-- SHOW STATS_TOTALS has cumulative values
SHOW STATS_TOTALS;

-- Watch for connection wait spikes
SHOW CLIENTS WHERE state = 'cl_waiting';
```

---

## Pgpool-II — Alternative Pooler

Pgpool-II adds load balancing and connection pooling together.

```bash
# Install
sudo apt install pgpool2
```

```ini
# /etc/pgpool2/pgpool.conf
# Backend nodes
backend_hostname0 = 'primary-host'
backend_port0 = 5432
backend_weight0 = 1
backend_flag0 = 'ALLOW_TO_FAILOVER'

backend_hostname1 = 'replica-host'
backend_port1 = 5432
backend_weight1 = 1        # equal weight = 50/50 read load balance
backend_flag1 = 'ALLOW_TO_FAILOVER'

# Load balancing (reads go to replicas)
load_balance_mode = on

# Connection pooling
num_init_children = 32     # max concurrent connections
max_pool = 4               # connections per child

# Replication detection
sr_check_period = 10
health_check_period = 10
```

---

## Built-in Pooling: `max_connections` Tuning

Even without a pooler, configure PostgreSQL correctly:

```ini
# postgresql.conf
max_connections = 100       # lower is better — each costs ~5MB RAM

# Reserve connections for admin
reserved_connections = 3    # PG 16+  (or use superuser_reserved_connections)
```

```sql
-- Check current connections
SELECT count(*), state, wait_event_type, wait_event
FROM pg_stat_activity
WHERE pid <> pg_backend_pid()
GROUP BY state, wait_event_type, wait_event
ORDER BY count DESC;

-- Kill idle connections
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE state = 'idle'
  AND state_change < NOW() - INTERVAL '10 minutes'
  AND pid <> pg_backend_pid();
```

---

## Supabase Pooler (Supavisor)

Supabase runs Supavisor (an Elixir-based pooler) as a managed service:

```
Direct connection: postgresql://user:pass@db.project.supabase.co:5432/postgres
Pooled (transaction): postgresql://user:pass@db.project.supabase.co:6543/postgres
```

Use the pooled connection string for serverless functions (Vercel, AWS Lambda) — each invocation needs a fresh connection, so session-level state doesn't matter.

---

## Connection Pooling Quick Reference

| Tool | Mode | Best For |
|---|---|---|
| PgBouncer | Session / Transaction / Statement | Standard choice, lightweight |
| Pgpool-II | Session + load balancing | Read scaling + pooling |
| RDS Proxy | Session / Transaction | AWS RDS / Aurora |
| Supavisor | Session / Transaction | Supabase |
| HikariCP | Session (app-side) | JVM applications |
| pgBouncer in K8s | Transaction | Kubernetes sidecars |

### Rules of thumb:
- Start with `transaction` pool mode
- Set `default_pool_size` = (CPU cores × 2) to (CPU cores × 4)
- Keep app-side pool size small (5–20) when using pgbouncer
- Watch `cl_waiting` — any waiting means pool is undersized
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
