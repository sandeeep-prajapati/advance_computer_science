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
