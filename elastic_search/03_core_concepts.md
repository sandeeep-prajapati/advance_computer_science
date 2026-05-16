# Chapter 3: Core Concepts

## The Mental Model

Before writing any queries, you need to understand how Elasticsearch organizes and thinks about data.

```
Cluster
└── Node(s)
    └── Index (e.g., "products")
        ├── Settings  (shards, replicas, refresh interval)
        ├── Mapping   (field types / schema)
        └── Shards
            └── Documents (JSON records)
```

---

## Index

An **index** is the primary container for data in Elasticsearch. Think of it as a table in SQL, but more flexible.

- Each index has a **unique name** (lowercase, no spaces)
- Each index has **settings** (shard count, replica count, etc.)
- Each index has a **mapping** (schema — field names and types)

```bash
# Create an index with explicit settings
PUT /orders
{
  "settings": {
    "number_of_shards": 3,
    "number_of_replicas": 1
  },
  "mappings": {
    "properties": {
      "order_id":    { "type": "keyword" },
      "customer":    { "type": "text" },
      "total":       { "type": "float" },
      "created_at":  { "type": "date" },
      "status":      { "type": "keyword" }
    }
  }
}
```

---

## Document

A **document** is a single JSON object stored in an index — analogous to a row in a database.

### Document Metadata Fields

| Field | Description |
|---|---|
| `_index` | The index this document belongs to |
| `_id` | Unique identifier within the index |
| `_version` | Version number (increments on each update) |
| `_seq_no` | Sequence number for optimistic concurrency |
| `_primary_term` | Primary term for optimistic concurrency |
| `_source` | The original JSON document |
| `_score` | Relevance score (only in search results) |

### Example Document

```json
{
  "_index": "orders",
  "_id": "100",
  "_version": 1,
  "_source": {
    "order_id": "ORD-100",
    "customer": "Alice Smith",
    "total": 299.99,
    "created_at": "2024-03-15T10:30:00Z",
    "status": "shipped"
  }
}
```

---

## Field

A **field** is a key-value pair inside a document. Elasticsearch needs to know the **type** of each field to index it correctly.

### Most Important Field Types

| Type | Use Case | Example |
|---|---|---|
| `text` | Full-text search (analyzed) | blog post, product description |
| `keyword` | Exact match, sorting, aggregations | email, status, tag |
| `integer` / `long` | Whole numbers | age, count, quantity |
| `float` / `double` | Decimal numbers | price, rating |
| `boolean` | True/false | in_stock, active |
| `date` | Timestamps and dates | created_at, updated_at |
| `object` | Nested JSON object | address: {city, zip} |
| `nested` | Array of objects with independent fields | line_items |
| `geo_point` | Latitude/longitude | store location |
| `ip` | IP addresses | client_ip |
| `completion` | Autocomplete | search suggestions |

> **Critical distinction:** `text` vs `keyword`
> - `text` → analyzed (tokenized, lowercased) → full-text search
> - `keyword` → exact string, not analyzed → filters, sorting, aggregations

---

## Mapping

A **mapping** is the schema of an index. It defines field names and their types.

### Dynamic Mapping (Auto-detection)

By default, Elasticsearch infers types when you first index a document:

```json
{ "name": "Alice", "age": 30, "active": true }
```

Elasticsearch creates:
- `name` → `text` + `keyword` (multi-field)
- `age` → `long`
- `active` → `boolean`

### Explicit Mapping (Recommended for Production)

```bash
PUT /users
{
  "mappings": {
    "properties": {
      "name": {
        "type": "text",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
      "email":      { "type": "keyword" },
      "age":        { "type": "integer" },
      "created_at": { "type": "date", "format": "strict_date_optional_time" },
      "address": {
        "type": "object",
        "properties": {
          "city":    { "type": "keyword" },
          "country": { "type": "keyword" }
        }
      }
    }
  }
}
```

> **Mappings are mostly immutable.** Once a field type is set, you cannot change it without reindexing. Plan your mapping carefully upfront.

---

## Cluster

A **cluster** is one or more Elasticsearch nodes working together. Key properties:

- Has a unique **cluster name** (e.g., `production-cluster`)
- Stores the complete dataset distributed across nodes
- Automatically handles node failures and data redistribution
- Coordinates distributed search across all nodes

```bash
# Check cluster health
GET _cluster/health

# Response
{
  "cluster_name": "my-cluster",
  "status": "green",          ← green/yellow/red
  "number_of_nodes": 3,
  "number_of_data_nodes": 3,
  "active_primary_shards": 10,
  "active_shards": 20,        ← primaries + replicas
  "unassigned_shards": 0
}
```

### Cluster Health Colors

| Color | Meaning |
|---|---|
| `green` | All primary and replica shards are assigned |
| `yellow` | All primary shards assigned, but some replicas are not |
| `red` | Some primary shards are unassigned — data loss risk |

---

## Node

A **node** is a single server in the cluster. Each node has specific **roles**:

| Role | Responsibility |
|---|---|
| `master` | Manages cluster state (index creation, shard allocation) |
| `data` | Stores shards, executes indexing and search |
| `ingest` | Runs ingest pipelines to transform documents before indexing |
| `coordinating` | Routes requests, merges shard results (all nodes do this by default) |
| `ml` | Runs machine learning jobs |
| `remote_cluster_client` | For cross-cluster search |

In development (single node), one node does everything. In production, roles are separated.

```bash
# List all nodes
GET _cat/nodes?v&h=name,ip,role,heap.percent,cpu
```

---

## Shards

A **shard** is a self-contained Lucene index — the actual unit of data storage and search.

### Primary Shards
- Determined at **index creation time** — cannot be changed afterward
- Data is distributed across primary shards by routing (default: hash of `_id`)
- More primary shards = more parallelism for indexing

### Replica Shards
- **Copies** of primary shards
- Can be changed at any time (`number_of_replicas`)
- Provide **high availability** (if primary fails, replica is promoted)
- Provide **read scaling** (searches can hit replicas)

```
Index with 3 primary shards, 1 replica:
Total shards = 3 primaries + 3 replicas = 6

Node 1:  [P0] [R1]
Node 2:  [P1] [R2]
Node 3:  [P2] [R0]
```

### Routing

When you index a document, ES determines which shard gets it:

```
shard = hash(_id) % number_of_primary_shards
```

Custom routing is possible:

```bash
PUT /products/_doc/1?routing=electronics
```

---

## The Refresh Cycle

Elasticsearch is **near real-time**, not real-time. Here's why:

1. Document is indexed → stored in an in-memory buffer
2. Every **1 second** (default), the buffer is written to a new **Lucene segment** (refresh)
3. The segment becomes visible to search
4. Periodically, segments are **merged** to improve performance

```
Index operation → buffer → (1s refresh) → segment → searchable
```

```bash
# Force a refresh (make all buffered docs immediately searchable)
POST /products/_refresh

# Disable auto-refresh for bulk indexing (then re-enable after)
PUT /products/_settings
{ "refresh_interval": "-1" }
```

---

## Summary: Analogies Table

| Elasticsearch | SQL Database | MongoDB |
|---|---|---|
| Index | Table | Collection |
| Document | Row | Document |
| Field | Column | Field |
| Mapping | Schema | Schema (optional) |
| Shard | Partition | Shard |
| `_id` | Primary Key | `_id` |
| Cluster | Database Server | Replica Set |
| Node | Server Instance | `mongod` process |
