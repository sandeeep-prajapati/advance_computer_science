# Chapter 8: Mappings & Data Types

## What is a Mapping?

A **mapping** is the schema for an Elasticsearch index. It defines:
- Field names
- Field data types
- How fields should be indexed and searched
- Which fields are stored, analyzed, or excluded

```bash
# View current mapping
GET /products/_mapping

# View specific field
GET /products/_mapping/field/price
```

---

## Dynamic vs. Explicit Mapping

### Dynamic Mapping (Auto-Detection)

Elasticsearch infers types from the first document:

| JSON Value | Detected Type |
|---|---|
| `"hello"` | `text` + `keyword` (multi-field) |
| `42` | `long` |
| `3.14` | `float` |
| `true/false` | `boolean` |
| `"2024-01-01"` | `date` (if matches date format) |
| `{ "city": "..." }` | `object` with nested properties |
| `[1, 2, 3]` | Same type as element (e.g., `long`) |

### Dynamic Mapping Modes

```bash
PUT /products
{
  "mappings": {
    "dynamic": "strict"    ← reject documents with unknown fields
  }
}
```

| Mode | Behavior |
|---|---|
| `true` (default) | Auto-detect and add new fields |
| `false` | Ignore unknown fields (not indexed, not searchable) |
| `strict` | Reject documents with unmapped fields (throw error) |

### Explicit Mapping (Recommended for Production)

```bash
PUT /orders
{
  "mappings": {
    "dynamic": "strict",
    "properties": {
      "order_id":    { "type": "keyword" },
      "customer_id": { "type": "keyword" },
      "total":       { "type": "float" },
      "created_at":  { "type": "date" },
      "status":      { "type": "keyword" },
      "notes":       { "type": "text" }
    }
  }
}
```

---

## Core Data Types

### String Types

#### `text` — Full-Text Search

```bash
"description": {
  "type":     "text",
  "analyzer": "english"
}
```

- Analyzed at index time (tokenized, lowercased, stemmed)
- Used for `match`, `match_phrase` queries
- Cannot be used for exact match, sort, or aggregations directly

#### `keyword` — Exact Values

```bash
"status":  { "type": "keyword" }
"sku":     { "type": "keyword" }
"user_id": { "type": "keyword" }
```

- Not analyzed — stored as-is
- Used for `term`, `terms` queries
- Can be sorted and aggregated
- `ignore_above: 256` silently ignores strings longer than 256 chars

---

### Numeric Types

| Type | Range | Storage |
|---|---|---|
| `byte` | -128 to 127 | 1 byte |
| `short` | -32,768 to 32,767 | 2 bytes |
| `integer` | ±2.1 billion | 4 bytes |
| `long` | ±9.2 × 10¹⁸ | 8 bytes |
| `float` | ~7 decimal digits | 4 bytes |
| `double` | ~15 decimal digits | 8 bytes |
| `half_float` | Lower precision | 2 bytes |
| `scaled_float` | Store as long with scaling factor | — |

```bash
"price": {
  "type":           "scaled_float",
  "scaling_factor": 100       ← store $9.99 as 999
}
```

---

### Date Type

```bash
"created_at": {
  "type":   "date",
  "format": "strict_date_optional_time||epoch_millis"
}
```

Stored internally as UTC epoch milliseconds.

Common formats:
- `strict_date_optional_time` — ISO 8601 (`2024-03-15T10:30:00Z`)
- `epoch_millis` — milliseconds since epoch
- `yyyy-MM-dd` — date only
- `dd/MM/yyyy HH:mm:ss`

You can accept multiple formats: `"format": "yyyy-MM-dd||yyyy/MM/dd||epoch_millis"`

---

### Boolean

```bash
"in_stock": { "type": "boolean" }
```

Accepts: `true`, `false`, `"true"`, `"false"`, `1`, `0`.

---

### Binary

Stores base64-encoded binary data. Not searchable by default.

```bash
"thumbnail": { "type": "binary", "store": true }
```

---

### Range Types

```bash
"age_group":      { "type": "integer_range" }
"event_window":   { "type": "date_range" }
"price_bracket":  { "type": "float_range" }
```

```json
{
  "age_group": { "gte": 18, "lte": 35 }
}
```

---

## Object and Nested Types

### `object` — Nested JSON

Automatically created for JSON sub-objects:

```bash
"address": {
  "type": "object",
  "properties": {
    "city":    { "type": "keyword" },
    "country": { "type": "keyword" },
    "zip":     { "type": "keyword" }
  }
}
```

Query: `"term": { "address.city": "Berlin" }`

> **Limitation:** When stored in an array, object relationships are flattened, which can cause incorrect query results.

### `nested` — Arrays of Objects (Preserves Relationships)

```bash
"line_items": {
  "type": "nested",
  "properties": {
    "product_id": { "type": "keyword" },
    "quantity":   { "type": "integer" },
    "price":      { "type": "float" }
  }
}
```

Query (requires `nested` query type):

```bash
{
  "query": {
    "nested": {
      "path": "line_items",
      "query": {
        "bool": {
          "must": [
            { "term":  { "line_items.product_id": "SKU-123" } },
            { "range": { "line_items.quantity": { "gte": 2 } } }
          ]
        }
      }
    }
  }
}
```

> Use `nested` when you need to query across fields of the same array element together. Otherwise, use `object`.

---

## Geo Types

### `geo_point`

```bash
"location": { "type": "geo_point" }
```

Accepted formats:
```json
{ "lat": 40.7128, "lon": -74.0060 }
"40.7128,-74.0060"
[-74.0060, 40.7128]    ← GeoJSON format (lon, lat)
```

Geo queries:
```bash
{
  "query": {
    "geo_distance": {
      "distance":  "5km",
      "location":  { "lat": 40.7128, "lon": -74.0060 }
    }
  }
}
```

### `geo_shape`

Stores polygons, lines, and complex shapes. Used for spatial relationships (within, intersects, contains).

---

## Specialized Types

### `ip` — IP Addresses

```bash
"client_ip": { "type": "ip" }
```

Supports CIDR range queries: `"term": { "client_ip": "192.168.0.0/16" }`

### `completion` — Autocomplete Suggester

```bash
"suggest": {
  "type": "completion",
  "analyzer": "simple"
}
```

### `flattened` — Dynamic Object as Keywords

All subfields treated as `keyword` regardless of type. Good for unknown/arbitrary JSON.

```bash
"labels": { "type": "flattened" }
```

### `dense_vector` — ML/kNN Search

```bash
"embedding": {
  "type":  "dense_vector",
  "dims":  768,
  "index": true,
  "similarity": "cosine"
}
```

---

## Mapping Parameters

| Parameter | Applies To | Description |
|---|---|---|
| `index` | All | Set to `false` to not index the field |
| `store` | All | Store separately for retrieval (rarely needed) |
| `doc_values` | keyword, numeric | On-disk data structure for sorting/aggs |
| `norms` | text | Scoring norms (set to `false` to save space) |
| `ignore_above` | keyword | Ignore strings longer than N chars |
| `null_value` | keyword, numeric | Value to use when field is null |
| `copy_to` | text, keyword | Copy field value to another field |
| `fields` | text | Multi-fields (index same value multiple ways) |
| `enabled` | object | `false` = store but don't index sub-fields |

---

## Adding Fields to Existing Mapping

You can add new fields at any time:

```bash
PUT /products/_mapping
{
  "properties": {
    "supplier_id": { "type": "keyword" },
    "weight_kg":   { "type": "float" }
  }
}
```

You **cannot** change the type of an existing field. To do that, you must **reindex**.

---

## Reindexing to Change Mapping

```bash
# 1. Create new index with updated mapping
PUT /products_v2
{
  "mappings": {
    "properties": {
      "price": { "type": "scaled_float", "scaling_factor": 100 }
    }
  }
}

# 2. Copy data from old index
POST /_reindex
{
  "source": { "index": "products" },
  "dest":   { "index": "products_v2" }
}

# 3. Switch alias
POST /_aliases
{
  "actions": [
    { "remove": { "index": "products",    "alias": "products_alias" } },
    { "add":    { "index": "products_v2", "alias": "products_alias" } }
  ]
}
```
