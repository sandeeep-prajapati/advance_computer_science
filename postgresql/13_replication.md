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
