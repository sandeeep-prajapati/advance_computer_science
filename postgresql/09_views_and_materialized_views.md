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
