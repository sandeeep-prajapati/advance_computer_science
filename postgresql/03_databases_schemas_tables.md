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
