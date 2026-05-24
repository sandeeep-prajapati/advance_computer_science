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
