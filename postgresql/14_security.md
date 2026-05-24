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
