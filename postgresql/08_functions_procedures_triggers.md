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
