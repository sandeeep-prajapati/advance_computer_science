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
