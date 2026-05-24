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
