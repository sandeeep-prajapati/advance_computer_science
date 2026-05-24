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
