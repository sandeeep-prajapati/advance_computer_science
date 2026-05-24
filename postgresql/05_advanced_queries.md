# Chapter 5: Advanced Queries

## JOINs

JOINs combine rows from two or more tables based on a related column.

```sql
-- INNER JOIN — only rows with matches in both tables
SELECT u.name, o.id AS order_id, o.amount
FROM users u
INNER JOIN orders o ON o.user_id = u.id;

-- LEFT JOIN — all rows from left table, nulls for non-matching right
SELECT u.name, COUNT(o.id) AS order_count
FROM users u
LEFT JOIN orders o ON o.user_id = u.id
GROUP BY u.id, u.name;

-- RIGHT JOIN — all rows from right table (less common, swap tables instead)
SELECT u.name, o.amount
FROM orders o
RIGHT JOIN users u ON o.user_id = u.id;

-- FULL OUTER JOIN — all rows from both tables
SELECT u.name, o.amount
FROM users u
FULL OUTER JOIN orders o ON o.user_id = u.id;

-- CROSS JOIN — cartesian product (every row with every row)
SELECT color, size FROM colors CROSS JOIN sizes;

-- SELF JOIN — join a table to itself
SELECT a.name AS employee, b.name AS manager
FROM employees a
JOIN employees b ON a.manager_id = b.id;
```

### Multiple JOINs:
```sql
SELECT
    u.name AS customer,
    o.id AS order_id,
    p.name AS product,
    oi.quantity,
    oi.unit_price
FROM orders o
JOIN users u ON u.id = o.user_id
JOIN order_items oi ON oi.order_id = o.id
JOIN products p ON p.id = oi.product_id
WHERE o.status = 'completed'
ORDER BY o.created_at DESC;
```

---

## Subqueries

A subquery is a query nested inside another query.

```sql
-- Subquery in WHERE (scalar subquery)
SELECT name FROM users
WHERE id = (SELECT user_id FROM orders WHERE amount = (SELECT MAX(amount) FROM orders));

-- Subquery in WHERE with IN
SELECT name FROM users
WHERE id IN (
    SELECT DISTINCT user_id FROM orders WHERE status = 'completed'
);

-- Correlated subquery (references outer query)
SELECT u.name,
    (SELECT COUNT(*) FROM orders WHERE user_id = u.id) AS order_count
FROM users u;

-- Subquery in FROM (derived table)
SELECT dept, AVG(salary) AS avg_salary
FROM (
    SELECT department AS dept, salary
    FROM employees
    WHERE hire_date > '2020-01-01'
) recent_hires
GROUP BY dept;

-- EXISTS with correlated subquery
SELECT * FROM users u
WHERE EXISTS (
    SELECT 1 FROM orders o
    WHERE o.user_id = u.id AND o.amount > 1000
);

-- NOT EXISTS
SELECT * FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM orders o WHERE o.user_id = u.id
);
```

---

## Common Table Expressions (CTEs)

CTEs (`WITH` clause) are named temporary result sets that make complex queries more readable.

```sql
-- Basic CTE
WITH active_users AS (
    SELECT id, name, email
    FROM users
    WHERE is_active = TRUE
),
high_value_orders AS (
    SELECT user_id, SUM(amount) AS total
    FROM orders
    WHERE status = 'completed'
    GROUP BY user_id
    HAVING SUM(amount) > 10000
)
SELECT u.name, u.email, hvo.total
FROM active_users u
JOIN high_value_orders hvo ON hvo.user_id = u.id
ORDER BY hvo.total DESC;

-- Multiple CTEs
WITH
    monthly_revenue AS (
        SELECT DATE_TRUNC('month', created_at) AS month, SUM(amount) AS revenue
        FROM orders
        WHERE status = 'completed'
        GROUP BY 1
    ),
    prev_month AS (
        SELECT LAG(revenue) OVER (ORDER BY month) AS prev_revenue, month, revenue
        FROM monthly_revenue
    )
SELECT
    month,
    revenue,
    prev_revenue,
    ROUND((revenue - prev_revenue) / prev_revenue * 100, 2) AS growth_pct
FROM prev_month
ORDER BY month;
```

### Recursive CTEs — for tree/graph structures:
```sql
-- Org chart traversal
WITH RECURSIVE org_chart AS (
    -- Base case: start with the CEO (no manager)
    SELECT id, name, manager_id, 1 AS level
    FROM employees
    WHERE manager_id IS NULL

    UNION ALL

    -- Recursive case: join with direct reports
    SELECT e.id, e.name, e.manager_id, oc.level + 1
    FROM employees e
    JOIN org_chart oc ON e.manager_id = oc.id
)
SELECT level, name FROM org_chart ORDER BY level, name;

-- Find all ancestors of a node
WITH RECURSIVE ancestors AS (
    SELECT id, parent_id, name
    FROM categories
    WHERE id = 42        -- starting node

    UNION ALL

    SELECT c.id, c.parent_id, c.name
    FROM categories c
    JOIN ancestors a ON c.id = a.parent_id
)
SELECT * FROM ancestors;
```

---

## Window Functions

Window functions perform calculations across a set of rows related to the current row, without collapsing them like GROUP BY.

```sql
-- ROW_NUMBER — unique sequential number within partition
SELECT
    id,
    name,
    department,
    salary,
    ROW_NUMBER() OVER (PARTITION BY department ORDER BY salary DESC) AS rank_in_dept
FROM employees;

-- RANK — like ROW_NUMBER but ties get same rank (with gaps)
-- DENSE_RANK — like RANK but no gaps

-- LAG / LEAD — access previous/next row
SELECT
    month,
    revenue,
    LAG(revenue, 1) OVER (ORDER BY month) AS prev_month_revenue,
    LEAD(revenue, 1) OVER (ORDER BY month) AS next_month_revenue,
    revenue - LAG(revenue) OVER (ORDER BY month) AS change
FROM monthly_revenue;

-- Running total
SELECT
    created_at::DATE AS date,
    amount,
    SUM(amount) OVER (ORDER BY created_at) AS running_total
FROM orders
WHERE status = 'completed';

-- Moving average (7-day)
SELECT
    date,
    revenue,
    AVG(revenue) OVER (
        ORDER BY date
        ROWS BETWEEN 6 PRECEDING AND CURRENT ROW
    ) AS moving_avg_7d
FROM daily_revenue;

-- FIRST_VALUE / LAST_VALUE
SELECT
    id,
    name,
    salary,
    FIRST_VALUE(salary) OVER (PARTITION BY department ORDER BY salary DESC) AS top_salary_in_dept
FROM employees;

-- NTILE — divide rows into N buckets
SELECT
    id,
    name,
    salary,
    NTILE(4) OVER (ORDER BY salary) AS quartile
FROM employees;

-- PERCENT_RANK / CUME_DIST
SELECT
    name,
    salary,
    ROUND(PERCENT_RANK() OVER (ORDER BY salary)::numeric * 100, 2) AS percentile
FROM employees;
```

---

## GROUPING SETS, ROLLUP, CUBE

```sql
-- ROLLUP — hierarchical subtotals
SELECT
    region,
    category,
    SUM(amount) AS revenue
FROM sales
GROUP BY ROLLUP (region, category);
-- Rows: (region, category), (region, NULL), (NULL, NULL)

-- CUBE — all possible combinations of subtotals
SELECT
    region,
    category,
    year,
    SUM(amount)
FROM sales
GROUP BY CUBE (region, category, year);

-- GROUPING SETS — explicit combinations
SELECT
    region,
    category,
    SUM(amount)
FROM sales
GROUP BY GROUPING SETS (
    (region, category),
    (region),
    (category),
    ()             -- grand total
);

-- GROUPING() function — tells you which columns were aggregated
SELECT
    GROUPING(region) AS region_grouped,
    region,
    SUM(amount)
FROM sales
GROUP BY ROLLUP (region);
```

---

## LATERAL Joins

LATERAL allows a subquery in FROM to reference columns from earlier in the FROM clause:

```sql
-- Get the 3 most recent orders for each user
SELECT u.name, recent.id, recent.amount, recent.created_at
FROM users u
CROSS JOIN LATERAL (
    SELECT id, amount, created_at
    FROM orders
    WHERE user_id = u.id
    ORDER BY created_at DESC
    LIMIT 3
) AS recent;

-- Get top product per category
SELECT c.name AS category, top_product.*
FROM categories c
CROSS JOIN LATERAL (
    SELECT p.name, p.price
    FROM products p
    WHERE p.category_id = c.id
    ORDER BY p.sales_count DESC
    LIMIT 1
) AS top_product;
```

---

## UNION, INTERSECT, EXCEPT

```sql
-- UNION — combine results, remove duplicates
SELECT name FROM customers
UNION
SELECT name FROM vendors;

-- UNION ALL — combine results, keep duplicates (faster)
SELECT name, 'customer' AS type FROM customers
UNION ALL
SELECT name, 'vendor' AS type FROM vendors;

-- INTERSECT — rows in both result sets
SELECT user_id FROM premium_users
INTERSECT
SELECT user_id FROM active_this_month;

-- EXCEPT — rows in first but not second
SELECT user_id FROM all_users
EXCEPT
SELECT user_id FROM cancelled_users;
```
