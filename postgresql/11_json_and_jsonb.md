# Chapter 11: JSON & JSONB

## JSON vs JSONB

| Feature | JSON | JSONB |
|---|---|---|
| Storage | Text (preserves whitespace, key order) | Binary (no whitespace, sorted keys) |
| Write speed | Faster (no parsing) | Slightly slower |
| Read speed | Slower (re-parsed each time) | Faster |
| Indexing | Not indexable | GIN/btree indexable |
| Operators | Basic | Full set |
| Recommendation | Rarely use | **Always use JSONB** |

---

## Creating JSONB Columns

```sql
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    attributes JSONB,
    metadata JSONB DEFAULT '{}'::JSONB
);

INSERT INTO products (name, attributes) VALUES
('Laptop', '{"brand": "Dell", "cpu": "Intel i7", "ram_gb": 16, "colors": ["black", "silver"]}'),
('Phone', '{"brand": "Apple", "model": "iPhone 15", "storage_gb": 256, "5g": true}'),
('Headphones', '{"brand": "Sony", "wireless": true, "noise_cancelling": true, "colors": ["black"]}');
```

---

## Operators

```sql
-- -> returns JSON/JSONB element (keeps as JSON)
SELECT attributes -> 'brand' FROM products;           -- "Dell"
SELECT attributes -> 'colors' -> 0 FROM products;     -- "black"

-- ->> returns text (unwraps the value)
SELECT attributes ->> 'brand' FROM products;          -- Dell (no quotes)
SELECT attributes ->> 'ram_gb' FROM products;         -- 16 (as text)

-- #> path operator (returns JSON)
SELECT attributes #> '{colors, 0}' FROM products;    -- "black"

-- #>> path operator (returns text)
SELECT attributes #>> '{colors, 0}' FROM products;   -- black

-- @> contains (JSONB only)
SELECT * FROM products WHERE attributes @> '{"brand": "Apple"}';
SELECT * FROM products WHERE attributes @> '{"wireless": true}';

-- <@ is contained by
SELECT * FROM products WHERE '{"brand": "Sony"}' <@ attributes;

-- ? key exists
SELECT * FROM products WHERE attributes ? 'noise_cancelling';

-- ?| any of these keys exist
SELECT * FROM products WHERE attributes ?| ARRAY['wireless', '5g'];

-- ?& all of these keys exist
SELECT * FROM products WHERE attributes ?& ARRAY['brand', 'model'];
```

---

## JSONB Functions

```sql
-- jsonb_object_keys — list all keys
SELECT jsonb_object_keys(attributes) FROM products WHERE id = 1;

-- jsonb_each — expand to rows of (key, value)
SELECT key, value FROM products, jsonb_each(attributes) WHERE id = 1;

-- jsonb_each_text — expand to rows of (key, text_value)
SELECT key, value FROM products, jsonb_each_text(attributes) WHERE id = 1;

-- jsonb_extract_path / jsonb_extract_path_text
SELECT jsonb_extract_path(attributes, 'brand') FROM products;
SELECT jsonb_extract_path_text(attributes, 'brand') FROM products;

-- jsonb_set — modify a value
UPDATE products
SET attributes = jsonb_set(attributes, '{brand}', '"DELL"')
WHERE id = 1;

-- jsonb_set creating nested paths
UPDATE products
SET attributes = jsonb_set(attributes, '{specs, ssd}', 'true', true)
WHERE id = 1;
-- true parameter = create path if it doesn't exist

-- jsonb_insert — insert into array
UPDATE products
SET attributes = jsonb_insert(attributes, '{colors, -1}', '"gold"')
WHERE id = 1;
-- -1 position = append to end

-- Merge JSONB (|| operator)
UPDATE products
SET attributes = attributes || '{"weight_kg": 1.5, "warranty_years": 2}'::JSONB
WHERE id = 1;

-- Remove key (- operator)
UPDATE products
SET attributes = attributes - 'old_field'
WHERE id = 1;

-- Remove nested key (#- operator)
UPDATE products
SET attributes = attributes #- '{specs, old_spec}'
WHERE id = 1;

-- jsonb_array_elements — expand JSON array to rows
SELECT value FROM products, jsonb_array_elements(attributes -> 'colors') WHERE id = 1;

-- jsonb_array_length
SELECT jsonb_array_length(attributes -> 'colors') FROM products WHERE id = 1;

-- json_agg — aggregate rows into JSON array
SELECT json_agg(row_to_json(u)) FROM users u WHERE is_active = TRUE;

-- jsonb_agg
SELECT jsonb_agg(DISTINCT role) FROM users;

-- json_build_object
SELECT json_build_object('id', id, 'name', name, 'email', email) FROM users;

-- row_to_json
SELECT row_to_json(u) FROM users u WHERE id = 1;
```

---

## Indexing JSONB

```sql
-- GIN index — supports @>, ?, ?|, ?& operators
CREATE INDEX idx_products_attrs ON products USING GIN(attributes);

-- Expression index — for specific field queries
CREATE INDEX idx_products_brand ON products((attributes->>'brand'));

-- GIN with jsonb_path_ops (smaller, only supports @>)
CREATE INDEX idx_products_attrs_path ON products USING GIN(attributes jsonb_path_ops);
```

---

## JSONB Path Queries (PostgreSQL 12+)

```sql
-- jsonb_path_query — XPath-like queries
SELECT jsonb_path_query(attributes, '$.colors[*]') FROM products WHERE id = 1;

-- Find products where any color is "black"
SELECT * FROM products
WHERE jsonb_path_exists(attributes, '$.colors[*] ? (@ == "black")');

-- Find products with RAM >= 16
SELECT * FROM products
WHERE jsonb_path_exists(attributes, '$.ram_gb ? (@ >= 16)');

-- @? operator (path exists)
SELECT * FROM products WHERE attributes @? '$.colors[*] ? (@ == "black")';

-- @@ operator (path condition is true)
SELECT * FROM products WHERE attributes @@ '$.ram_gb >= 16';
```

---

## Practical Patterns

### EAV Alternative (Entity-Attribute-Value):
```sql
-- Instead of hundreds of nullable columns OR EAV tables:
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    preferences JSONB DEFAULT '{}',
    settings JSONB DEFAULT '{}'
);

-- Flexible schema per user
UPDATE users SET preferences = preferences || '{"theme": "dark", "language": "en"}' WHERE id = 1;
UPDATE users SET settings = settings || '{"notifications": {"email": true, "sms": false}}' WHERE id = 1;

SELECT preferences->>'theme' FROM users WHERE id = 1;
SELECT * FROM users WHERE preferences @> '{"theme": "dark"}';
```

### Storing Dynamic Forms:
```sql
CREATE TABLE form_submissions (
    id BIGSERIAL PRIMARY KEY,
    form_id BIGINT REFERENCES forms(id),
    user_id BIGINT REFERENCES users(id),
    data JSONB NOT NULL,
    submitted_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO form_submissions (form_id, user_id, data) VALUES
(1, 42, '{"name": "Alice", "age": 30, "interests": ["coding", "hiking"]}');

SELECT data->>'name', data->>'age'
FROM form_submissions
WHERE form_id = 1 AND data @> '{"interests": ["coding"]}';
```

---

## JSONB vs Relational — When to Use Which

**Use JSONB when:**
- Schema varies per row (product attributes, form responses)
- Data is accessed as a blob, not queried field-by-field
- You're prototyping and the schema is still evolving
- Storing external API responses as-is

**Use relational tables when:**
- You filter/sort/join on specific fields frequently
- You need strong integrity constraints (foreign keys, check constraints)
- Multiple rows reference the same data
- You need efficient updates to specific fields
