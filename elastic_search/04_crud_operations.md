# Chapter 4: CRUD Operations

## Overview

Elasticsearch uses a REST API for all operations. Every document operation maps to an HTTP method:

| Operation | HTTP | Endpoint |
|---|---|---|
| Create (specific ID) | `PUT` | `/<index>/_doc/<id>` |
| Create (auto ID) | `POST` | `/<index>/_doc` |
| Read | `GET` | `/<index>/_doc/<id>` |
| Update (partial) | `POST` | `/<index>/_update/<id>` |
| Replace (full) | `PUT` | `/<index>/_doc/<id>` |
| Delete | `DELETE` | `/<index>/_doc/<id>` |
| Bulk | `POST` | `/_bulk` |

---

## Create: Indexing a Document

### PUT — Specific ID

```bash
PUT /products/_doc/1
{
  "name": "Laptop Pro 15",
  "brand": "Dell",
  "price": 1299.99,
  "category": "electronics",
  "in_stock": true,
  "tags": ["laptop", "work", "portable"],
  "created_at": "2024-03-15"
}
```

Response:

```json
{
  "_index": "products",
  "_id": "1",
  "_version": 1,
  "result": "created",
  "_shards": { "total": 2, "successful": 1 }
}
```

### POST — Auto-Generated ID

```bash
POST /products/_doc
{
  "name": "Wireless Mouse",
  "brand": "Logitech",
  "price": 49.99
}
```

Response includes auto-generated `_id`:
```json
{ "_id": "a1B2c3D4e5F6" }
```

### Create Only (Fail if Exists)

Use `_create` endpoint or `op_type=create`:

```bash
PUT /products/_create/1
{ "name": "Laptop Pro 15" }
```

Returns `409 Conflict` if document already exists.

---

## Read: Getting a Document

### Get by ID

```bash
GET /products/_doc/1
```

Response:

```json
{
  "_index": "products",
  "_id": "1",
  "_version": 1,
  "found": true,
  "_source": {
    "name": "Laptop Pro 15",
    "price": 1299.99
  }
}
```

### Check Existence (no body returned)

```bash
HEAD /products/_doc/1
# 200 OK = exists, 404 = not found
```

### Retrieve Specific Fields Only

```bash
GET /products/_doc/1?_source_includes=name,price
```

### Get Multiple Documents

```bash
GET /products/_mget
{
  "ids": ["1", "2", "3"]
}
```

---

## Update: Modifying a Document

### Partial Update (most common)

```bash
POST /products/_update/1
{
  "doc": {
    "price": 1199.99,
    "in_stock": false
  }
}
```

Only specified fields are updated. Unmentioned fields are preserved.

### Update with Script (Groovy/Painless)

```bash
POST /products/_update/1
{
  "script": {
    "source": "ctx._source.price -= params.discount",
    "params": {
      "discount": 100
    }
  }
}
```

Common script operations:
```bash
# Increment a counter
"source": "ctx._source.views += 1"

# Add to an array
"source": "ctx._source.tags.add(params.tag)",
"params": { "tag": "bestseller" }

# Remove from an array
"source": "ctx._source.tags.remove(ctx._source.tags.indexOf(params.tag))",
"params": { "tag": "outdated" }

# Delete a field
"source": "ctx._source.remove('temp_field')"

# Conditionally delete document
"source": "if (ctx._source.stock == 0) { ctx.op = 'delete' }"
```

### Upsert (Update or Insert)

```bash
POST /products/_update/99
{
  "doc": {
    "name": "New Product",
    "price": 49.99
  },
  "doc_as_upsert": true
}
```

Creates the document if it doesn't exist; updates it if it does.

---

## Delete: Removing Documents

### Delete by ID

```bash
DELETE /products/_doc/1
```

Response:

```json
{
  "_id": "1",
  "result": "deleted",
  "_version": 2
}
```

### Delete by Query

```bash
POST /products/_delete_by_query
{
  "query": {
    "term": {
      "category": "discontinued"
    }
  }
}
```

### Delete an Entire Index

```bash
DELETE /products
```

> This is irreversible. Always double-check the index name.

---

## Bulk API

The **Bulk API** is the most efficient way to perform many operations at once. It reduces HTTP overhead dramatically.

### Format

Each operation takes **two lines**:
1. Action line (metadata)
2. Document line (body) — not needed for `delete`

```bash
POST /_bulk
{ "index": { "_index": "products", "_id": "1" } }
{ "name": "iPhone 15", "price": 999 }
{ "index": { "_index": "products", "_id": "2" } }
{ "name": "Samsung S24", "price": 849 }
{ "update": { "_index": "products", "_id": "1" } }
{ "doc": { "price": 949 } }
{ "delete": { "_index": "products", "_id": "3" } }
```

### Bulk to a Specific Index

```bash
POST /products/_bulk
{ "index": { "_id": "10" } }
{ "name": "AirPods Pro", "price": 249 }
{ "index": { "_id": "11" } }
{ "name": "Apple Watch", "price": 399 }
```

### Bulk Response

```json
{
  "took": 23,
  "errors": false,
  "items": [
    { "index": { "_id": "10", "result": "created", "status": 201 } },
    { "index": { "_id": "11", "result": "created", "status": 201 } }
  ]
}
```

If `"errors": true`, check each item's `status` and `error` fields.

---

## Optimistic Concurrency Control

Elasticsearch uses `_seq_no` and `_primary_term` to prevent lost updates (no locking needed).

```bash
# Read the current seq_no
GET /products/_doc/1
# → "_seq_no": 5, "_primary_term": 1

# Update only if version matches
POST /products/_update/1?if_seq_no=5&if_primary_term=1
{
  "doc": { "price": 999 }
}
```

If another process updated the document between your read and write, you get `409 Conflict`.

---

## Versioning

Every document update increments `_version`. Deleted documents stay at their last version briefly, then are cleaned up.

```bash
GET /products/_doc/1
# "_version": 3
```

---

## Update by Query

Update multiple documents matching a query:

```bash
POST /products/_update_by_query
{
  "script": {
    "source": "ctx._source.price = ctx._source.price * 0.9"
  },
  "query": {
    "term": {
      "category": "electronics"
    }
  }
}
```

---

## Practical Tips

| Scenario | Recommended Approach |
|---|---|
| Ingest 1M documents | Bulk API with batches of 5,000–10,000 |
| Update one field | `_update` with `doc` partial update |
| Compute a new value | `_update` with Painless script |
| Delete old records | `_delete_by_query` with date range filter |
| Retry on conflict | Add `retry_on_conflict=3` to update |
| Avoid overwrites | Use `_create` endpoint or `op_type=create` |
