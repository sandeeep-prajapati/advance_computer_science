# Elasticsearch: The Complete Developer Book
> A comprehensive, single-file reference covering everything from your first index to production-grade clusters, the ELK stack, security, and advanced search.

---

## Quick Navigation

- [Table of Contents](#table-of-contents)
- [Chapter 1: Introduction to Elasticsearch](#chapter-1-introduction-to-elasticsearch)
- [Chapter 2: Setup & Installation](#chapter-2-setup--installation)
- [Chapter 3: Core Concepts](#chapter-3-core-concepts)
- [Chapter 4: CRUD Operations](#chapter-4-crud-operations)
- [Chapter 5: Search & Query DSL](#chapter-5-search--query-dsl)
- [Chapter 6: Aggregations](#chapter-6-aggregations)
- [Chapter 7: Text Analysis & Analyzers](#chapter-7-text-analysis--analyzers)
- [Chapter 8: Mappings & Data Types](#chapter-8-mappings--data-types)
- [Chapter 9: Index Management](#chapter-9-index-management)
- [Chapter 10: Cluster Architecture](#chapter-10-cluster-architecture)
- [Chapter 11: Shards & Replicas](#chapter-11-shards--replicas)
- [Chapter 12: Performance & Optimization](#chapter-12-performance--optimization)
- [Chapter 13: Security & Authentication](#chapter-13-security--authentication)
- [Chapter 14: Monitoring & Observability](#chapter-14-monitoring--observability)
- [Chapter 15: Elasticsearch with Kibana](#chapter-15-elasticsearch-with-kibana)
- [Chapter 16: The ELK Stack](#chapter-16-the-elk-stack)
- [Chapter 17: Snapshots & Backup](#chapter-17-snapshots--backup)
- [Chapter 18: Advanced Search](#chapter-18-advanced-search)
- [Chapter 19: Clients & Integration](#chapter-19-clients--integration)
- [Chapter 20: Learning Roadmap](#chapter-20-learning-roadmap)

---

# Elasticsearch Developer Book

A comprehensive, chapter-wise guide to mastering Elasticsearch — from first index to production-grade search, aggregations, the ELK stack, and cluster management.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_introduction.md) | **Introduction to Elasticsearch** | What is ES, why it exists, key concepts, use cases, ES vs SQL |
| [02](02_setup_and_installation.md) | **Setup & Installation** | Docker, local install, Elastic Cloud, Dev Tools, first index |
| [03](03_core_concepts.md) | **Core Concepts** | Index, Document, Field, Mapping, Cluster, Node, Shard, Replica |
| [04](04_crud_operations.md) | **CRUD Operations** | Index/Get/Update/Delete documents, Bulk API, upsert |
| [05](05_search_and_query_dsl.md) | **Search & Query DSL** | match, term, bool, range, exists, wildcard, multi_match |
| [06](06_aggregations.md) | **Aggregations** | Metric aggs, bucket aggs, pipeline aggs, nested aggs |
| [07](07_text_analysis.md) | **Text Analysis & Analyzers** | Analyzers, tokenizers, filters, custom analyzers, language analyzers |
| [08](08_mappings_and_data_types.md) | **Mappings & Data Types** | Dynamic vs explicit mapping, field types, nested, object, flattened |
| [09](09_index_management.md) | **Index Management** | Templates, aliases, ILM, rollover, reindex, shrink |
| [10](10_cluster_architecture.md) | **Cluster Architecture** | Master/data/ingest nodes, discovery, routing, split-brain |
| [11](11_shards_and_replicas.md) | **Shards & Replicas** | Primary vs replica, routing, over-sharding, rebalancing |
| [12](12_performance_and_optimization.md) | **Performance & Optimization** | Bulk indexing, refresh interval, force merge, caching, JVM tuning |
| [13](13_security.md) | **Security & Authentication** | TLS, X-Pack, roles, users, API keys, field/document-level security |
| [14](14_monitoring_and_observability.md) | **Monitoring & Observability** | Cluster health, cat APIs, Stack Monitoring, alerting, slow logs |
| [15](15_kibana.md) | **Elasticsearch with Kibana** | Discover, Dashboards, Lens, KQL, Dev Tools, Index Patterns |
| [16](16_elk_stack.md) | **The ELK Stack** | Logstash pipelines, Filebeat, Metricbeat, ingestion patterns |
| [17](17_snapshots_and_backup.md) | **Snapshots & Backup** | Repository setup, snapshot lifecycle, restore, cross-cluster |
| [18](18_advanced_search.md) | **Advanced Search** | Fuzzy, phonetic, geo, nested queries, percolate, more-like-this |
| [19](19_clients_and_integration.md) | **Clients & Integration** | Python, Node.js, Java clients, REST API patterns, pagination |
| [20](20_learning_roadmap.md) | **Learning Roadmap** | 5-week hands-on plan from zero to production-ready |

---

## How to Read This Book

**If you're brand new to Elasticsearch:** Read 1 → 2 → 3 → 4 → 5 in order. This gives you the core mental model.

**If you're building a search feature:** Read 5 → 6 → 7 → 8 → 18.

**If you're running the ELK stack:** Read 15 → 16 → 14.

**If you're going to production:** Read 9 → 10 → 11 → 12 → 13 → 17.

**If you need to tune performance:** Read 12 → 11 → 9.

**If something is slow or broken:** Jump to Chapter 14 (Monitoring) and Chapter 12 (Performance).

---

## Key Concepts at a Glance

```
Elasticsearch Cluster
├── Nodes
│   ├── Master Node         (cluster state, index management)
│   ├── Data Node           (stores shards, executes queries)
│   ├── Ingest Node         (pre-processes documents via pipelines)
│   ├── Coordinating Node   (routes requests, merges results)
│   └── ML Node             (machine learning tasks)
│
├── Indices
│   ├── Index               (like a database table)
│   ├── Mapping             (schema — field names and types)
│   ├── Alias               (virtual name pointing to one or more indices)
│   └── Template            (auto-apply settings/mappings to new indices)
│
├── Documents
│   ├── Document            (a JSON record — like a row)
│   ├── Field               (a key-value pair in a document)
│   └── _id                 (unique identifier per document)
│
├── Shards
│   ├── Primary Shard       (original copy of data)
│   └── Replica Shard       (copy of primary — for HA and read scaling)
│
├── Search
│   ├── Query DSL           (JSON-based query language)
│   ├── Relevance Score     (_score — how well a doc matches)
│   ├── Aggregations        (analytics, grouping, metrics)
│   └── Analyzer            (text processing pipeline for full-text search)
│
└── Operations
    ├── Index API           (create/update document)
    ├── Bulk API            (batch operations)
    ├── Search API          (query documents)
    ├── Reindex API         (copy/transform between indices)
    └── ILM                 (Index Lifecycle Management)
```


---

# Chapter 1: Introduction to Elasticsearch

## What is Elasticsearch?

**Elasticsearch** is an open-source, distributed **search and analytics engine** built on top of Apache Lucene. It is part of the **Elastic Stack** (formerly ELK Stack) and is designed to store, search, and analyze large volumes of data in near real-time.

> **Analogy:** Think of Elasticsearch as a **super-powered search engine with a database built in**:
> - You store JSON documents.
> - It indexes every field automatically.
> - You can search millions of documents in milliseconds.

---

## Why Elasticsearch Exists

Traditional relational databases are great for structured queries but struggle with:
- Full-text search across millions of records
- Fuzzy matching ("did you mean…?")
- Complex analytics over time-series data
- Log aggregation from hundreds of services

Elasticsearch was built specifically for these workloads.

| Problem | Elasticsearch Solution |
|---|---|
| Slow full-text search | Inverted index — search in O(1) |
| Can't scale reads | Replica shards for parallel queries |
| Rigid schema | Dynamic mapping, schema-free JSON |
| No analytics | Aggregations pipeline |
| Single point of failure | Distributed cluster with replication |

---

## Core Concepts (Quick Overview)

### Index
An **index** is the top-level container for your data — similar to a database table. Each index has a **mapping** (schema) that defines the fields and types.

### Document
A **document** is a single JSON record stored in an index — similar to a row in a table. Every document gets a unique `_id`.

```json
{
  "_index": "products",
  "_id": "1",
  "_source": {
    "name": "iPhone 15",
    "price": 999,
    "category": "electronics",
    "description": "Apple smartphone with A17 chip"
  }
}
```

### Shard
A **shard** is a horizontal partition of an index. Elasticsearch splits data across multiple shards to enable parallel processing and horizontal scaling.

### Cluster & Node
A **cluster** is a collection of nodes (servers) that together hold all your data. Each **node** holds some shards and participates in indexing and searching.

---

## Elasticsearch Architecture Overview

```
┌────────────────────────────────────────────────────────┐
│                    ES CLUSTER                           │
│                                                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │    Node 1    │  │    Node 2    │  │    Node 3    │  │
│  │  (master +   │  │  (data node) │  │  (data node) │  │
│  │   data)      │  │              │  │              │  │
│  │ [P0][P1][R2] │  │ [P2][R0][R1] │  │ [R0][R1][P2] │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│                                                         │
│  P0/P1/P2 = Primary Shards                             │
│  R0/R1/R2 = Replica Shards                             │
└────────────────────────────────────────────────────────┘
          ▲
          │ HTTP REST API
          │
     ┌────┴────┐
     │ Client  │  (Kibana, curl, Python, Node.js)
     └─────────┘
```

---

## How Elasticsearch Stores and Searches Data

### The Inverted Index

Unlike a relational database that scans rows, Elasticsearch builds an **inverted index** — a mapping from each unique word (token) to all documents containing it:

```
Word          →  Document IDs
──────────────────────────────
"apple"       →  [1, 5, 23]
"smartphone"  →  [1, 7]
"iphone"      →  [1, 5]
"chip"        →  [1, 3, 5]
```

When you search for "apple smartphone", Elasticsearch intersects those lists and scores results by relevance — in microseconds.

### Relevance Scoring

Every search result gets a `_score`. Higher score = better match. The default algorithm is **BM25** (Best Match 25), which factors in:
- **Term frequency** — how often the word appears in the doc
- **Inverse document frequency** — how rare the word is across all docs
- **Field length** — shorter fields score higher (title vs. full article)

---

## Elasticsearch vs SQL Databases

| Concept | SQL | Elasticsearch |
|---|---|---|
| Database | Database | Index |
| Table | Table | Index (1:1 in modern ES) |
| Row | Row | Document |
| Column | Column | Field |
| Schema | CREATE TABLE | Mapping |
| Primary key | Primary key | `_id` |
| Query | SQL SELECT | Query DSL (JSON) |
| Join | JOIN | Nested / Denormalization |
| Full-text | LIKE '%...%' (slow) | match query (fast) |
| Analytics | GROUP BY | Aggregations |

> **Key difference:** Elasticsearch is not a replacement for a relational database. It excels at search and analytics, not transactional operations (no ACID, no complex joins).

---

## Real-World Use Cases

### Log & Event Analytics (ELK Stack)
Ingest logs from hundreds of services via Filebeat/Logstash → store in Elasticsearch → visualize in Kibana. Used by Netflix, LinkedIn, and most DevOps teams.

### E-Commerce Search
Index products with titles, descriptions, tags. Support fuzzy search ("ipone" → "iphone"), faceted filters (brand, price range), autocomplete, and personalized ranking.

### Application Performance Monitoring (APM)
Elastic APM collects traces, metrics, and errors from your apps and stores them in Elasticsearch for unified observability.

### Security & SIEM
Elastic SIEM ingests security events, network flows, and endpoint data to detect threats, correlate events, and trigger alerts.

### Geospatial Search
Store location data (lat/lon, geo-shapes) and run geo-distance queries like "find all restaurants within 2km of me."

### Autocomplete & Suggestions
Use completion suggesters and edge n-gram tokenizers to power instant search-as-you-type features.

---

## Key Objects Summary

| Object | What It Does |
|---|---|
| `Index` | Container for documents (like a table) |
| `Document` | A single JSON record |
| `Field` | A key-value pair in a document |
| `Mapping` | Schema: defines field names and types |
| `Shard` | Horizontal partition of an index |
| `Replica` | Copy of a shard for HA and read scaling |
| `Alias` | Virtual name pointing to one or more indices |
| `Template` | Auto-applied settings/mappings for new indices |
| `Pipeline` | Pre-processing steps applied on ingest |
| `Analyzer` | Text processing pipeline (tokenize, filter, normalize) |
| `Aggregation` | Analytics operation (sum, avg, bucket grouping) |
| `Snapshot` | Backup of indices to a repository |

---

## Elasticsearch vs Competitors

| Feature | Elasticsearch | OpenSearch | Solr | Meilisearch |
|---|---|---|---|---|
| Full-text search | Excellent | Excellent | Good | Good |
| Analytics/Aggregations | Best in class | Good | Basic | Limited |
| ELK ecosystem | Native | Fork | No | No |
| Ease of setup | Medium | Medium | Complex | Easy |
| Managed cloud | Elastic Cloud | AWS OpenSearch | Various | Meilisearch Cloud |
| License | Elastic License 2.0 | Apache 2.0 | Apache 2.0 | MIT |
| Best for | Enterprise, logs, APM | AWS-native, cost | Legacy enterprise | Simple app search |


---

# Chapter 2: Setup & Installation

## Your Options

| Method | Best For |
|---|---|
| **Docker / Docker Compose** | Local development (recommended) |
| **Local install (tar/deb/rpm)** | Custom setups, bare metal |
| **Elastic Cloud** | Production, managed, no ops overhead |
| **AWS OpenSearch** | AWS-native, open-source fork |
| **Minikube / Helm** | Kubernetes-based development |

---

## Option A: Docker (Recommended for Development)

### Single Node (Quickest Start)

```bash
docker run -d \
  --name elasticsearch \
  -p 9200:9200 \
  -p 9300:9300 \
  -e "discovery.type=single-node" \
  -e "xpack.security.enabled=false" \
  docker.elastic.co/elasticsearch/elasticsearch:8.13.0
```

Verify it's running:

```bash
curl http://localhost:9200
```

Expected output:

```json
{
  "name" : "es-node-1",
  "cluster_name" : "docker-cluster",
  "version" : {
    "number" : "8.13.0"
  },
  "tagline" : "You Know, for Search"
}
```

---

### Docker Compose (ES + Kibana)

```yaml
# docker-compose.yml
version: '3.8'
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.13.0
    container_name: elasticsearch
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - ES_JAVA_OPTS=-Xms1g -Xmx1g
    ports:
      - "9200:9200"
    volumes:
      - esdata:/usr/share/elasticsearch/data

  kibana:
    image: docker.elastic.co/kibana/kibana:8.13.0
    container_name: kibana
    ports:
      - "5601:5601"
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    depends_on:
      - elasticsearch

volumes:
  esdata:
```

```bash
docker-compose up -d

# Check logs
docker-compose logs -f elasticsearch
```

Access Kibana at: `http://localhost:5601`

---

## Option B: Local Install (Linux)

```bash
# Import GPG key
wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo gpg --dearmor -o /usr/share/keyrings/elasticsearch-keyring.gpg

# Add repository
echo "deb [signed-by=/usr/share/keyrings/elasticsearch-keyring.gpg] https://artifacts.elastic.co/packages/8.x/apt stable main" | sudo tee /etc/apt/sources.list.d/elastic-8.x.list

# Install
sudo apt-get update && sudo apt-get install elasticsearch

# Start
sudo systemctl enable elasticsearch
sudo systemctl start elasticsearch

# Verify
curl -X GET "localhost:9200/"
```

---

## Option C: macOS (Homebrew)

```bash
brew tap elastic/tap
brew install elastic/tap/elasticsearch-full
brew services start elastic/tap/elasticsearch-full
```

---

## Elasticsearch Configuration

Key config file: `/etc/elasticsearch/elasticsearch.yml` (or `config/elasticsearch.yml`)

```yaml
# Cluster name
cluster.name: my-cluster

# Node name
node.name: node-1

# Data and log paths
path.data: /var/data/elasticsearch
path.logs: /var/log/elasticsearch

# Network
network.host: 0.0.0.0
http.port: 9200

# Discovery (for multi-node clusters)
discovery.seed_hosts: ["host1", "host2"]
cluster.initial_master_nodes: ["node-1"]

# Security (for production)
xpack.security.enabled: true
xpack.security.transport.ssl.enabled: true
```

---

## Using the Dev Tools Console (Kibana)

The **Dev Tools Console** in Kibana is the best way to interact with Elasticsearch interactively.

Access: `http://localhost:5601` → Dev Tools → Console

```
GET /
GET _cluster/health
GET _cat/indices?v
```

### REST API Conventions

All Elasticsearch operations use HTTP verbs:

| HTTP Method | Operation |
|---|---|
| `GET` | Retrieve a document or run a query |
| `POST` | Create a document (auto-generate `_id`) or run a query |
| `PUT` | Create/update a document with a specific `_id` |
| `DELETE` | Delete a document or index |
| `HEAD` | Check existence |

### URL Structure

```
http://localhost:9200/<index>/<endpoint>
                       ↑
                       index name

Examples:
GET  /products/_doc/1          ← get doc with id=1
PUT  /products/_doc/1          ← create/replace doc with id=1
POST /products/_doc             ← create doc with auto-generated id
POST /products/_search          ← search index
DELETE /products/_doc/1        ← delete doc
DELETE /products               ← delete entire index
```

---

## Your First Index and Document

### Create an Index

```bash
PUT /products
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  }
}
```

### Index a Document

```bash
PUT /products/_doc/1
{
  "name": "iPhone 15",
  "brand": "Apple",
  "price": 999,
  "in_stock": true,
  "description": "Latest Apple smartphone with A17 Pro chip"
}
```

### Search It

```bash
GET /products/_search
{
  "query": {
    "match": {
      "description": "apple smartphone"
    }
  }
}
```

---

## Useful `_cat` APIs

```bash
# List all indices with stats
GET _cat/indices?v

# Node health
GET _cat/nodes?v

# Shard allocation
GET _cat/shards?v

# Cluster health
GET _cluster/health

# Index settings
GET /products/_settings

# Index mapping
GET /products/_mapping
```

---

## Installing Kibana (Standalone)

```bash
sudo apt-get install kibana
sudo systemctl enable kibana
sudo systemctl start kibana
```

Kibana connects to Elasticsearch on `localhost:9200` by default. Edit `/etc/kibana/kibana.yml` to change:

```yaml
server.port: 5601
server.host: "0.0.0.0"
elasticsearch.hosts: ["http://localhost:9200"]
```

---

## Version Compatibility

| ES Version | Release | Notable Features |
|---|---|---|
| 5.x | 2016 | Unified stack, X-Pack |
| 6.x | 2017 | Single type per index |
| 7.x | 2019 | Removed types, vote-only master, ECS |
| 8.x | 2022 | Security on by default, kNN vector search, NLP |

> Always match your Elasticsearch and Kibana **major versions** exactly (e.g., ES 8.x with Kibana 8.x).


---

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


---

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


---

# Chapter 5: Search & Query DSL

## The Search API

All searches go through the `_search` endpoint:

```bash
GET /products/_search
{
  "query": { ... },
  "from": 0,
  "size": 10,
  "sort": [ ... ],
  "_source": ["name", "price"]
}
```

| Parameter | Default | Description |
|---|---|---|
| `query` | match_all | The query to filter/score documents |
| `from` | 0 | Offset for pagination |
| `size` | 10 | Number of results to return |
| `sort` | by `_score` | Sort criteria |
| `_source` | all fields | Which fields to return |
| `aggs` | none | Aggregation clauses |

---

## Two Types of Clauses

### Query Context (Scores Documents)
Used inside `query` → affects `_score`. "How well does this document match?"

### Filter Context (Binary Yes/No)
Used inside `filter` or `must_not` → no scoring, **cached**, faster.
"Does this document match? Yes/No."

> Use **filter** for exact matches (date range, status, keyword) and **query** for full-text relevance. Filters are cached and significantly faster.

---

## Match All (Default)

```bash
GET /products/_search
{
  "query": { "match_all": {} }
}
```

Returns all documents with equal score of 1.0.

---

## Full-Text Queries

### `match` — Standard Full-Text Search

Analyzes the search term and matches analyzed field content.

```bash
GET /products/_search
{
  "query": {
    "match": {
      "description": "apple smartphone"
    }
  }
}
```

- Default operator is `OR` (documents with "apple" OR "smartphone" match)
- Change to AND to require all terms:

```bash
"match": {
  "description": {
    "query": "apple smartphone",
    "operator": "and"
  }
}
```

### `match_phrase` — Exact Phrase in Order

```bash
GET /products/_search
{
  "query": {
    "match_phrase": {
      "description": "wireless noise cancelling"
    }
  }
}
```

All words must appear in the given order, adjacent (or within a `slop` distance).

### `multi_match` — Search Multiple Fields

```bash
GET /products/_search
{
  "query": {
    "multi_match": {
      "query": "laptop gaming",
      "fields": ["name^3", "description", "tags"],
      "type": "best_fields"
    }
  }
}
```

`^3` boosts the `name` field score by 3x.

**Types:**
- `best_fields` (default) — highest score from any single field
- `most_fields` — sum scores from all matching fields
- `cross_fields` — treats fields as one big field (good for first+last name)
- `phrase` — `match_phrase` across multiple fields

---

## Term-Level Queries (Exact Match)

These do NOT analyze the query string — used on `keyword`, `integer`, `boolean`, `date` fields.

### `term` — Exact Match

```bash
GET /products/_search
{
  "query": {
    "term": {
      "category": { "value": "electronics" }
    }
  }
}
```

### `terms` — Match Any of Multiple Values

```bash
GET /products/_search
{
  "query": {
    "terms": {
      "status": ["shipped", "delivered"]
    }
  }
}
```

### `range` — Numeric/Date Range

```bash
GET /products/_search
{
  "query": {
    "range": {
      "price": {
        "gte": 100,
        "lte": 500
      }
    }
  }
}

# Date range
{
  "query": {
    "range": {
      "created_at": {
        "gte": "2024-01-01",
        "lt": "2025-01-01",
        "format": "yyyy-MM-dd"
      }
    }
  }
}
```

### `exists` — Field Has a Value

```bash
GET /products/_search
{
  "query": {
    "exists": { "field": "discount_price" }
  }
}
```

### `ids` — Match Specific Document IDs

```bash
{
  "query": {
    "ids": { "values": ["1", "2", "3"] }
  }
}
```

### `prefix` — Starts With

```bash
{ "query": { "prefix": { "sku": { "value": "ELEC-" } } } }
```

### `wildcard` — Glob Pattern

```bash
{ "query": { "wildcard": { "email": { "value": "*@gmail.com" } } } }
```

> Wildcards are slow on large datasets — avoid leading wildcards (`*term`).

### `regexp` — Regular Expression

```bash
{ "query": { "regexp": { "sku": { "value": "ELEC-[0-9]{4}" } } } }
```

---

## Boolean Query — Combining Clauses

The `bool` query is the most powerful and commonly used query. It combines multiple clauses:

| Clause | Behavior |
|---|---|
| `must` | Document MUST match. Contributes to score. |
| `filter` | Document MUST match. Does NOT affect score. Cached. |
| `should` | Document SHOULD match. Increases score if it does. |
| `must_not` | Document MUST NOT match. Does not affect score. Cached. |

```bash
GET /products/_search
{
  "query": {
    "bool": {
      "must": [
        { "match": { "description": "wireless headphones" } }
      ],
      "filter": [
        { "term": { "in_stock": true } },
        { "range": { "price": { "lte": 300 } } }
      ],
      "should": [
        { "term": { "brand": "Sony" } },
        { "term": { "brand": "Bose" } }
      ],
      "must_not": [
        { "term": { "category": "refurbished" } }
      ]
    }
  }
}
```

### `minimum_should_match`

```bash
"bool": {
  "should": [
    { "term": { "tag": "sale" } },
    { "term": { "tag": "new" } },
    { "term": { "tag": "featured" } }
  ],
  "minimum_should_match": 2
}
```

---

## Sorting

```bash
GET /products/_search
{
  "query": { "match_all": {} },
  "sort": [
    { "price": { "order": "asc" } },
    { "created_at": { "order": "desc" } },
    "_score"
  ]
}
```

> You cannot sort by `text` fields. Use the `.keyword` sub-field for text sorting.

```bash
"sort": [{ "name.keyword": { "order": "asc" } }]
```

---

## Pagination

### From/Size (Simple Pagination)

```bash
GET /products/_search
{
  "from": 20,
  "size": 10
}
```

> Avoid deep pagination — `from=10000` requires ES to collect and discard 10,000 results. Max is controlled by `index.max_result_window` (default 10,000).

### Search After (Deep Pagination)

```bash
# First page
GET /products/_search
{
  "size": 10,
  "sort": [{ "price": "asc" }, { "_id": "asc" }]
}

# Next page — use last result's sort values
GET /products/_search
{
  "size": 10,
  "sort": [{ "price": "asc" }, { "_id": "asc" }],
  "search_after": [299.99, "abc123"]
}
```

### Scroll API (Bulk Export — Deprecated in favor of search_after)

Use `search_after` with a Point-In-Time (PIT) for modern bulk exports:

```bash
# Open a PIT
POST /products/_pit?keep_alive=1m

# Search with PIT
GET /_search
{
  "pit": { "id": "<pit_id>", "keep_alive": "1m" },
  "sort": [{ "_id": "asc" }],
  "search_after": ["last_id"]
}
```

---

## Source Filtering

Control what is returned in `_source`:

```bash
GET /products/_search
{
  "query": { "match_all": {} },
  "_source": {
    "includes": ["name", "price", "category"],
    "excludes": ["internal_*"]
  }
}
```

Or just list fields: `"_source": ["name", "price"]`

---

## Highlighting

Return highlighted snippets around matching terms:

```bash
GET /products/_search
{
  "query": { "match": { "description": "wireless" } },
  "highlight": {
    "fields": {
      "description": {
        "pre_tags": ["<em>"],
        "post_tags": ["</em>"],
        "number_of_fragments": 1,
        "fragment_size": 150
      }
    }
  }
}
```

Response includes `highlight.description` with `<em>wireless</em>` in context.

---

## Common Query Patterns

### Find products under $100 in stock

```bash
{
  "query": {
    "bool": {
      "filter": [
        { "range": { "price": { "lte": 100 } } },
        { "term":  { "in_stock": true } }
      ]
    }
  }
}
```

### Full-text search with category filter

```bash
{
  "query": {
    "bool": {
      "must": { "match": { "name": "gaming laptop" } },
      "filter": { "term": { "category": "electronics" } }
    }
  }
}
```

### Boost recent documents

```bash
{
  "query": {
    "function_score": {
      "query": { "match": { "title": "elasticsearch" } },
      "functions": [
        {
          "gauss": {
            "published_at": {
              "origin": "now",
              "scale": "30d",
              "decay": 0.5
            }
          }
        }
      ]
    }
  }
}
```


---

# Chapter 6: Aggregations

## What Are Aggregations?

Aggregations are Elasticsearch's analytics engine. They allow you to compute **summaries, statistics, and grouped analytics** over your data — like a powerful `GROUP BY` in SQL but capable of nesting, pipelines, and sub-aggregations.

```
Aggregation Types:
├── Metric Aggs     → Single numeric output (avg, sum, min, max, count)
├── Bucket Aggs     → Group documents into buckets (terms, range, date_histogram)
└── Pipeline Aggs   → Compute from other aggregation outputs (moving avg, derivative)
```

---

## Basic Structure

```bash
GET /orders/_search
{
  "size": 0,            ← don't return hits, just aggregation results
  "aggs": {
    "my_agg_name": {   ← arbitrary name you choose
      "agg_type": {    ← e.g., terms, avg, date_histogram
        "field": "category"
      }
    }
  }
}
```

---

## Metric Aggregations

### `avg`, `sum`, `min`, `max`

```bash
GET /orders/_search
{
  "size": 0,
  "aggs": {
    "avg_order_value":  { "avg":  { "field": "total" } },
    "total_revenue":    { "sum":  { "field": "total" } },
    "min_order":        { "min":  { "field": "total" } },
    "max_order":        { "max":  { "field": "total" } }
  }
}
```

Response:

```json
{
  "aggregations": {
    "avg_order_value":  { "value": 89.50 },
    "total_revenue":    { "value": 1243500.00 },
    "min_order":        { "value": 5.99 },
    "max_order":        { "value": 4999.00 }
  }
}
```

### `value_count` — Count Non-Null Values

```bash
"aggs": {
  "total_orders": { "value_count": { "field": "order_id" } }
}
```

### `stats` — All Basic Stats at Once

```bash
"aggs": {
  "price_stats": { "stats": { "field": "price" } }
}
```

Returns: `count`, `min`, `max`, `avg`, `sum`.

### `extended_stats` — Includes Variance and Std Deviation

```bash
"aggs": {
  "price_extended": { "extended_stats": { "field": "price" } }
}
```

### `cardinality` — Approximate Count of Distinct Values

```bash
"aggs": {
  "unique_customers": { "cardinality": { "field": "customer_id" } }
}
```

### `percentiles` — Distribution

```bash
"aggs": {
  "load_time_percentiles": {
    "percentiles": {
      "field": "response_time_ms",
      "percents": [50, 75, 90, 95, 99]
    }
  }
}
```

---

## Bucket Aggregations

Bucket aggs divide documents into groups (buckets). Each bucket can contain sub-aggregations.

### `terms` — Group by Field Value

```bash
GET /products/_search
{
  "size": 0,
  "aggs": {
    "by_category": {
      "terms": {
        "field": "category",
        "size": 10,
        "order": { "_count": "desc" }
      }
    }
  }
}
```

Response:

```json
{
  "aggregations": {
    "by_category": {
      "buckets": [
        { "key": "electronics", "doc_count": 1500 },
        { "key": "clothing",    "doc_count": 980 },
        { "key": "books",       "doc_count": 670 }
      ]
    }
  }
}
```

### `range` — Group by Numeric Ranges

```bash
"aggs": {
  "price_ranges": {
    "range": {
      "field": "price",
      "ranges": [
        { "to": 50 },
        { "from": 50, "to": 200 },
        { "from": 200 }
      ]
    }
  }
}
```

### `date_histogram` — Group by Time

```bash
GET /logs/_search
{
  "size": 0,
  "aggs": {
    "events_over_time": {
      "date_histogram": {
        "field":             "@timestamp",
        "calendar_interval": "day",
        "format":            "yyyy-MM-dd",
        "time_zone":         "America/New_York"
      }
    }
  }
}
```

Common intervals: `minute`, `hour`, `day`, `week`, `month`, `quarter`, `year`.

### `histogram` — Numeric Intervals

```bash
"aggs": {
  "price_distribution": {
    "histogram": {
      "field":    "price",
      "interval": 100
    }
  }
}
```

### `filter` — Single-Condition Bucket

```bash
"aggs": {
  "high_value_orders": {
    "filter": { "range": { "total": { "gte": 1000 } } },
    "aggs": {
      "avg_high_value": { "avg": { "field": "total" } }
    }
  }
}
```

### `filters` — Multiple Named Buckets

```bash
"aggs": {
  "by_status": {
    "filters": {
      "filters": {
        "pending":   { "term": { "status": "pending" } },
        "shipped":   { "term": { "status": "shipped" } },
        "delivered": { "term": { "status": "delivered" } }
      }
    }
  }
}
```

---

## Nested Aggregations (Sub-Aggregations)

The most powerful pattern: bucket aggs containing metric aggs.

### Average Revenue Per Category

```bash
GET /orders/_search
{
  "size": 0,
  "aggs": {
    "by_category": {
      "terms": { "field": "category", "size": 10 },
      "aggs": {
        "avg_revenue":  { "avg": { "field": "total" } },
        "total_revenue": { "sum": { "field": "total" } },
        "order_count":  { "value_count": { "field": "_id" } }
      }
    }
  }
}
```

### Monthly Sales with Top Products Per Month

```bash
GET /orders/_search
{
  "size": 0,
  "aggs": {
    "monthly": {
      "date_histogram": {
        "field": "created_at",
        "calendar_interval": "month"
      },
      "aggs": {
        "monthly_revenue": { "sum": { "field": "total" } },
        "top_products": {
          "terms": { "field": "product_id", "size": 5 }
        }
      }
    }
  }
}
```

---

## Pipeline Aggregations

Pipeline aggs operate on the output of other aggregations (not on documents directly).

### `avg_bucket` — Average of Bucket Values

```bash
GET /orders/_search
{
  "size": 0,
  "aggs": {
    "monthly": {
      "date_histogram": {
        "field": "created_at",
        "calendar_interval": "month"
      },
      "aggs": {
        "monthly_revenue": { "sum": { "field": "total" } }
      }
    },
    "avg_monthly_revenue": {
      "avg_bucket": {
        "buckets_path": "monthly>monthly_revenue"
      }
    }
  }
}
```

### `cumulative_sum` — Running Total

```bash
"aggs": {
  "cumulative_revenue": {
    "cumulative_sum": {
      "buckets_path": "monthly>monthly_revenue"
    }
  }
}
```

### `derivative` — Rate of Change

```bash
"aggs": {
  "revenue_change": {
    "derivative": {
      "buckets_path": "monthly>monthly_revenue"
    }
  }
}
```

### `moving_avg` — Smoothed Trend

```bash
"aggs": {
  "revenue_trend": {
    "moving_avg": {
      "buckets_path": "monthly>monthly_revenue",
      "window": 3
    }
  }
}
```

---

## Scoping Aggregations with a Query

Aggregations apply to the same document set as the query filter:

```bash
GET /orders/_search
{
  "query": {
    "bool": {
      "filter": [
        { "term": { "status": "delivered" } },
        { "range": { "created_at": { "gte": "2024-01-01" } } }
      ]
    }
  },
  "size": 0,
  "aggs": {
    "revenue_by_region": {
      "terms": { "field": "region" },
      "aggs": {
        "total_revenue": { "sum": { "field": "total" } }
      }
    }
  }
}
```

---

## Global Aggregation (Ignore Query Context)

```bash
GET /products/_search
{
  "query": { "term": { "category": "electronics" } },
  "size": 0,
  "aggs": {
    "electronics_avg": { "avg": { "field": "price" } },
    "all_products_avg": {
      "global": {},
      "aggs": {
        "avg_price": { "avg": { "field": "price" } }
      }
    }
  }
}
```

---

## Top Hits — Return Sample Documents Per Bucket

```bash
GET /orders/_search
{
  "size": 0,
  "aggs": {
    "by_status": {
      "terms": { "field": "status" },
      "aggs": {
        "sample_orders": {
          "top_hits": {
            "size": 2,
            "sort": [{ "total": { "order": "desc" } }],
            "_source": ["order_id", "total", "created_at"]
          }
        }
      }
    }
  }
}
```

---

## Performance Tips

| Tip | Why |
|---|---|
| Use `keyword` fields for `terms` aggs | `text` fields are not aggregatable by default |
| Set `size: 0` when you only need aggs | Skips fetching hit documents |
| Use `filter` aggs instead of filtered queries when mixing | Filter agg cached separately |
| Avoid `cardinality` on high-cardinality fields | Uses HyperLogLog++ — memory intensive |
| Limit `terms` `size` | Large sizes slow down coordination across shards |


---

# Chapter 7: Text Analysis & Analyzers

## What is Text Analysis?

When you index a `text` field, Elasticsearch doesn't store the raw string — it runs the string through an **analyzer** that transforms it into **tokens** (terms) stored in the inverted index.

The same analyzer is applied at **index time** and **search time** to ensure consistent matching.

```
Input:    "Quick Brown FOX jumped over the lazy Dog!"
          ↓ Analyzer
Tokens:   ["quick", "brown", "fox", "jump", "lazi", "dog"]
          ↓ Stored in inverted index
```

---

## Anatomy of an Analyzer

An analyzer is a pipeline with three components:

```
Input Text
    ↓
Character Filters   (strip HTML, replace chars)
    ↓
Tokenizer           (split into tokens — exactly ONE tokenizer)
    ↓
Token Filters       (lowercase, stop words, stemming — zero or more)
    ↓
Token Stream        (stored in inverted index)
```

---

## Built-In Analyzers

### `standard` (Default)

Splits on whitespace and punctuation, lowercases, removes common stop words.

```bash
POST /_analyze
{
  "analyzer": "standard",
  "text": "Quick Brown FOX jumped over the lazy Dog!"
}
```

Tokens: `["quick", "brown", "fox", "jumped", "over", "the", "lazy", "dog"]`

### `simple`

Splits on any non-letter character, lowercases. No stop words.

Tokens: `["quick", "brown", "fox", "jumped", "over", "the", "lazy", "dog"]`

### `whitespace`

Splits on whitespace only. Preserves case and punctuation.

Tokens: `["Quick", "Brown", "FOX!", "jumped", "over", "the", "lazy", "Dog!"]`

### `english` (and other language analyzers)

Applies stemming and language-specific stop words:

```bash
POST /_analyze
{
  "analyzer": "english",
  "text": "Running quickly and jumping over fences"
}
```

Tokens: `["run", "quick", "jump", "over", "fenc"]` — stemmed to root forms.

### `keyword`

No tokenization — entire string as one token. Used with `keyword` type fields.

---

## Testing Analyzers

Use the `_analyze` API to test how text is tokenized:

```bash
POST /_analyze
{
  "analyzer": "standard",
  "text": "My email is user@example.com"
}

# Test with a specific field's analyzer
POST /products/_analyze
{
  "field": "description",
  "text": "Wireless Bluetooth Headphones"
}

# Test a custom pipeline
POST /_analyze
{
  "tokenizer": "standard",
  "filter": ["lowercase", "stop", "snowball"],
  "text": "Running quickly through the forest"
}
```

---

## Custom Analyzers

Define custom analyzers in index settings:

```bash
PUT /my_index
{
  "settings": {
    "analysis": {
      "analyzer": {
        "my_custom_analyzer": {
          "type":            "custom",
          "char_filter":     ["html_strip"],
          "tokenizer":       "standard",
          "filter":          ["lowercase", "stop", "my_stemmer"]
        },
        "autocomplete_analyzer": {
          "type":      "custom",
          "tokenizer": "standard",
          "filter":    ["lowercase", "autocomplete_filter"]
        },
        "autocomplete_search": {
          "type":      "custom",
          "tokenizer": "standard",
          "filter":    ["lowercase"]
        }
      },
      "filter": {
        "my_stemmer": {
          "type":     "stemmer",
          "language": "english"
        },
        "autocomplete_filter": {
          "type":     "edge_ngram",
          "min_gram": 1,
          "max_gram": 20
        }
      }
    }
  }
}
```

---

## Common Tokenizers

| Tokenizer | Description | Example Input → Tokens |
|---|---|---|
| `standard` | Word boundaries, punctuation | `"hello world"` → `[hello, world]` |
| `whitespace` | Split on whitespace only | `"foo bar"` → `[foo, bar]` |
| `keyword` | Entire string as one token | `"foo bar"` → `[foo bar]` |
| `letter` | Split on non-letter chars | `"foo123bar"` → `[foo, bar]` |
| `ngram` | Generate n-grams | `"cat"` (n=2) → `[ca, at]` |
| `edge_ngram` | N-grams from start | `"cat"` → `[c, ca, cat]` |
| `pattern` | Split on regex | Split on comma: `"a,b,c"` → `[a, b, c]` |
| `path_hierarchy` | File path segments | `"/a/b/c"` → `[/a, /a/b, /a/b/c]` |

---

## Common Token Filters

| Filter | Description |
|---|---|
| `lowercase` | Converts tokens to lowercase |
| `stop` | Removes common stop words (the, is, are) |
| `stemmer` | Reduces words to root form (running→run) |
| `snowball` | Snowball stemmer algorithm |
| `synonym` | Replace tokens with synonyms |
| `asciifolding` | Convert accented chars (café→cafe) |
| `unique` | Remove duplicate tokens |
| `length` | Remove tokens outside min/max length |
| `edge_ngram` | Generate edge n-grams (for autocomplete) |
| `ngram` | Generate n-grams |
| `shingle` | Generate word n-grams (bigrams, trigrams) |

---

## Character Filters

| Filter | Description |
|---|---|
| `html_strip` | Remove HTML tags |
| `mapping` | Replace specific characters/strings |
| `pattern_replace` | Replace using regex |

```bash
# Remove HTML and replace & with and
{
  "char_filter": ["html_strip"],
  "text": "<p>Soup &amp; Salad</p>"
}
# → "Soup & Salad"
```

---

## Autocomplete with Edge N-Gram

A classic use case: "as-you-type" search.

```bash
PUT /products
{
  "settings": {
    "analysis": {
      "filter": {
        "autocomplete_filter": {
          "type":     "edge_ngram",
          "min_gram": 2,
          "max_gram": 20
        }
      },
      "analyzer": {
        "autocomplete": {
          "type":      "custom",
          "tokenizer": "standard",
          "filter":    ["lowercase", "autocomplete_filter"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "name": {
        "type": "text",
        "analyzer": "autocomplete",
        "search_analyzer": "standard"
      }
    }
  }
}
```

Index time: `"Laptop"` → `["la", "lap", "lapt", "lapto", "laptop"]`
Search time: `"lapt"` → `["lapt"]` — matched exactly.

---

## Synonym Filters

```bash
PUT /products
{
  "settings": {
    "analysis": {
      "filter": {
        "tech_synonyms": {
          "type": "synonym",
          "synonyms": [
            "phone, mobile, smartphone",
            "laptop => notebook, portable computer",
            "tv, television, smart tv"
          ]
        }
      },
      "analyzer": {
        "synonym_analyzer": {
          "tokenizer": "standard",
          "filter":    ["lowercase", "tech_synonyms"]
        }
      }
    }
  }
}
```

Now searching for "phone" also matches documents containing "mobile" and "smartphone".

---

## Field Mapping with Multiple Analyzers

You can use different analyzers for the same field:

```bash
"properties": {
  "title": {
    "type":            "text",
    "analyzer":        "english",         ← index-time analyzer
    "search_analyzer": "standard"         ← query-time analyzer
  }
}
```

---

## Multi-Fields: text + keyword

A very common pattern — store a field as both analyzed `text` and exact `keyword`:

```bash
"properties": {
  "name": {
    "type": "text",
    "analyzer": "standard",
    "fields": {
      "keyword": {
        "type": "keyword",
        "ignore_above": 256
      },
      "english": {
        "type": "text",
        "analyzer": "english"
      }
    }
  }
}
```

Usage:
- `name` → full-text search (analyzed)
- `name.keyword` → exact match, sorting, aggregations
- `name.english` → stemmed English search

---

## Update Analyzers on Existing Index

You **cannot change analyzers without reindexing**, but you can add new analyzers:

```bash
# 1. Close the index
POST /products/_close

# 2. Update settings
PUT /products/_settings
{
  "analysis": {
    "analyzer": {
      "new_analyzer": { ... }
    }
  }
}

# 3. Reopen
POST /products/_open

# 4. Reindex existing docs to apply new analyzer
POST /products/_update_by_query?conflicts=proceed
```


---

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


---

# Chapter 9: Index Management

## Index Settings

Every index has two categories of settings:

| Category | Changeable After Creation? | Examples |
|---|---|---|
| **Static** | No (requires close/reopen) | `number_of_shards`, `codec` |
| **Dynamic** | Yes (live update) | `number_of_replicas`, `refresh_interval` |

```bash
# Create index with explicit settings
PUT /logs-2024-01
{
  "settings": {
    "number_of_shards":   3,
    "number_of_replicas": 1,
    "refresh_interval":   "30s",
    "max_result_window":  50000
  }
}

# Update dynamic settings on existing index
PUT /logs-2024-01/_settings
{
  "number_of_replicas": 2,
  "refresh_interval":   "1s"
}
```

---

## Index Templates

Templates automatically apply settings and mappings to new indices that match a naming pattern.

### Component Templates (Reusable Pieces)

```bash
# Define a reusable mapping component
PUT /_component_template/timestamp_mapping
{
  "template": {
    "mappings": {
      "properties": {
        "@timestamp": { "type": "date" },
        "host":        { "type": "keyword" },
        "log_level":   { "type": "keyword" }
      }
    }
  }
}

# Define reusable settings component
PUT /_component_template/log_settings
{
  "template": {
    "settings": {
      "number_of_shards":   2,
      "number_of_replicas": 1,
      "refresh_interval":   "10s"
    }
  }
}
```

### Index Template (Compose Components)

```bash
PUT /_index_template/logs_template
{
  "index_patterns": ["logs-*", "events-*"],
  "priority":        100,
  "composed_of":     ["timestamp_mapping", "log_settings"],
  "template": {
    "mappings": {
      "properties": {
        "message": { "type": "text" },
        "service": { "type": "keyword" }
      }
    }
  }
}
```

Now any index created with a name matching `logs-*` automatically gets this mapping and settings.

---

## Index Aliases

An **alias** is a virtual name pointing to one or more real indices. Aliases allow zero-downtime reindexing and abstraction over index names.

```bash
# Create alias pointing to one index
POST /_aliases
{
  "actions": [
    { "add": { "index": "products_v1", "alias": "products" } }
  ]
}

# Query through alias (transparent to client)
GET /products/_search

# Switch alias atomically (zero downtime reindex cutover)
POST /_aliases
{
  "actions": [
    { "remove": { "index": "products_v1", "alias": "products" } },
    { "add":    { "index": "products_v2", "alias": "products" } }
  ]
}
```

### Write Index Alias

Point an alias to multiple indices but designate one as the write target:

```bash
POST /_aliases
{
  "actions": [
    { "add": { "index": "logs-2024-01", "alias": "logs", "is_write_index": false } },
    { "add": { "index": "logs-2024-02", "alias": "logs", "is_write_index": false } },
    { "add": { "index": "logs-2024-03", "alias": "logs", "is_write_index": true  } }
  ]
}
```

Now writes to `logs` go to `logs-2024-03`, but searches span all three.

### Filtered Alias

```bash
POST /_aliases
{
  "actions": [{
    "add": {
      "index": "orders",
      "alias": "eu_orders",
      "filter": { "term": { "region": "EU" } }
    }
  }]
}
```

`GET /eu_orders/_search` only returns EU orders.

---

## Index Lifecycle Management (ILM)

ILM automates the lifecycle of time-series indices (logs, metrics, events). Indices move through phases:

```
Hot  →  Warm  →  Cold  →  Frozen  →  Delete
(write)  (read) (rare)  (archive)   (expire)
```

### Create an ILM Policy

```bash
PUT /_ilm/policy/logs_policy
{
  "policy": {
    "phases": {
      "hot": {
        "actions": {
          "rollover": {
            "max_primary_shard_size": "50gb",
            "max_age":                "7d"
          },
          "set_priority": { "priority": 100 }
        }
      },
      "warm": {
        "min_age": "7d",
        "actions": {
          "shrink":      { "number_of_shards": 1 },
          "forcemerge":  { "max_num_segments": 1 },
          "set_priority": { "priority": 50 }
        }
      },
      "cold": {
        "min_age": "30d",
        "actions": {
          "freeze":       {},
          "set_priority": { "priority": 0 }
        }
      },
      "delete": {
        "min_age": "90d",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}
```

### Attach ILM Policy to an Index Template

```bash
PUT /_index_template/logs_template
{
  "index_patterns": ["logs-*"],
  "template": {
    "settings": {
      "index.lifecycle.name":            "logs_policy",
      "index.lifecycle.rollover_alias":  "logs"
    }
  }
}
```

### Bootstrap the First Index

```bash
PUT /logs-000001
{
  "aliases": {
    "logs": { "is_write_index": true }
  }
}
```

ILM will automatically create `logs-000002`, `logs-000003`, etc. on rollover.

---

## Rollover API

Manually trigger a rollover when conditions are met:

```bash
POST /logs/_rollover
{
  "conditions": {
    "max_age":                "1d",
    "max_primary_shard_size": "10gb",
    "max_docs":               1000000
  }
}
```

Creates a new index `logs-000002` and points the write alias to it.

---

## Reindex API

Copy documents from one index to another:

```bash
POST /_reindex
{
  "source": {
    "index": "products_old",
    "query": {
      "term": { "active": true }
    }
  },
  "dest": {
    "index": "products_new",
    "op_type": "create"    ← skip if already exists
  }
}
```

### Reindex with Transform (Script)

```bash
POST /_reindex
{
  "source": { "index": "products" },
  "dest":   { "index": "products_v2" },
  "script": {
    "source": """
      ctx._source.price_cents = (int)(ctx._source.price * 100);
      ctx._source.remove('price');
    """
  }
}
```

### Reindex from Remote Cluster

```bash
POST /_reindex
{
  "source": {
    "remote": {
      "host": "http://remote-cluster:9200"
    },
    "index": "products"
  },
  "dest": { "index": "products_local" }
}
```

---

## Shrink API

Reduce the number of primary shards (index must be read-only and all shards on one node):

```bash
# Mark source as read-only and move shards
PUT /logs-2024-01/_settings
{
  "settings": {
    "index.blocks.write":   true,
    "index.routing.allocation.require._name": "node-1"
  }
}

# Wait for relocation, then shrink
POST /logs-2024-01/_shrink/logs-2024-01-shrunk
{
  "settings": {
    "index.number_of_shards":   1,
    "index.blocks.write":       false
  }
}
```

---

## Force Merge

Merge Lucene segments to improve read performance (typically for read-only indices):

```bash
POST /logs-2024-01/_forcemerge?max_num_segments=1
```

> Only run on read-only (closed/frozen) indices in production — it's resource-intensive.

---

## Open & Close Index

Closing an index frees memory (no searches allowed):

```bash
POST /logs-2024-01/_close
POST /logs-2024-01/_open
```

---

## Index Stats

```bash
# All indices
GET /_cat/indices?v&s=index

# Specific index stats
GET /products/_stats

# Index segments info
GET /products/_segments
```

---

## Useful Index Operations Reference

| Operation | Command |
|---|---|
| Create index | `PUT /<index>` |
| Delete index | `DELETE /<index>` |
| Close index | `POST /<index>/_close` |
| Open index | `POST /<index>/_open` |
| Refresh | `POST /<index>/_refresh` |
| Flush | `POST /<index>/_flush` |
| Force merge | `POST /<index>/_forcemerge` |
| Reindex | `POST /_reindex` |
| Rollover | `POST /<alias>/_rollover` |
| Shrink | `POST /<index>/_shrink/<new>` |
| Get settings | `GET /<index>/_settings` |
| Get mapping | `GET /<index>/_mapping` |
| List aliases | `GET /_aliases` |


---

# Chapter 10: Cluster Architecture

## Overview

An Elasticsearch **cluster** is a group of one or more nodes that collectively hold all your data and provide indexing and search capabilities.

```
Cluster: "production-cluster"
├── Node 1  (master-eligible, data)
├── Node 2  (master-eligible, data)
├── Node 3  (master-eligible, data)
├── Node 4  (data)
├── Node 5  (data)
└── Node 6  (coordinating-only / ingest)
```

---

## Node Roles

Each node can serve one or more roles. In small clusters, a single node does everything. In production, roles are separated for isolation and performance.

### Master Node

Manages cluster-wide metadata:
- Index creation and deletion
- Shard allocation decisions
- Tracking which nodes are in the cluster

```yaml
# elasticsearch.yml
node.roles: [master]
```

> Recommendation: 3 dedicated master-eligible nodes (odd number to prevent split-brain).

### Data Node

Stores shards and executes index/search operations:

```yaml
node.roles: [data]
# Or tier-specific:
node.roles: [data_hot]    # Fast SSDs — active indices
node.roles: [data_warm]   # Cheaper storage — read-mostly indices
node.roles: [data_cold]   # Object storage — rare access
node.roles: [data_frozen] # Searchable snapshot — archive
```

### Ingest Node

Runs **ingest pipelines** to pre-process documents before indexing (parse, enrich, transform):

```yaml
node.roles: [ingest]
```

### Coordinating Node (Coordinating-Only)

Routes client requests to the correct data nodes, collects and merges results:

```yaml
node.roles: []    ← empty = coordinating only
```

A dedicated coordinating node offloads the coordination work from data and master nodes in large clusters.

### ML Node

Runs machine learning jobs (anomaly detection, NLP):

```yaml
node.roles: [ml]
```

---

## Master Election

Elasticsearch uses a **Raft-based consensus algorithm** to elect a master node.

1. Nodes with `master` role participate in elections
2. A node wins if it gets votes from a **quorum** (majority) of master-eligible nodes
3. The quorum prevents split-brain:

```
For 3 master-eligible nodes, quorum = 2
For 5 master-eligible nodes, quorum = 3
```

```yaml
# Configure who can be initial master candidates (set once at cluster bootstrap)
cluster.initial_master_nodes: ["node-1", "node-2", "node-3"]
```

> Never use an even number of master-eligible nodes. With 2 nodes, losing one means no quorum, and the cluster goes read-only.

---

## Cluster State

The **cluster state** is the metadata about the entire cluster maintained by the master node:
- Which indices exist
- Their mappings and settings
- Shard allocation (which shard on which node)

The master node broadcasts cluster state changes to all nodes.

---

## Node Discovery

How nodes find each other in the cluster:

```yaml
# Seed hosts (known nodes to contact for discovery)
discovery.seed_hosts: ["192.168.1.10", "192.168.1.11", "192.168.1.12"]

# Or use cloud-native discovery
discovery.ec2.tag.Environment: production
```

---

## Shard Allocation

The master decides which node gets each shard. Factors:
- **Disk usage** — avoid nodes with > 85% disk (watermark.low)
- **Shard count balance** — spread shards evenly
- **Rack/zone awareness** — keep primaries and replicas in different failure domains

### Zone Awareness (High Availability)

```yaml
# On each node, set its zone
node.attr.zone: us-east-1a

# In index settings
PUT /products/_settings
{
  "index.routing.allocation.awareness.attributes": "zone"
}

# Globally require one replica per zone
PUT /_cluster/settings
{
  "persistent": {
    "cluster.routing.allocation.awareness.attributes":       "zone",
    "cluster.routing.allocation.awareness.force.zone.values": "us-east-1a,us-east-1b,us-east-1c"
  }
}
```

### Controlling Shard Placement

```bash
# Pin an index to nodes with specific attributes
PUT /products/_settings
{
  "index.routing.allocation.require.box_type": "hot"
}

# Exclude a specific node from an index
PUT /products/_settings
{
  "index.routing.allocation.exclude._name": "node-3"
}
```

---

## Disk Watermarks

Elasticsearch automatically stops allocating new shards when nodes get full:

```bash
PUT /_cluster/settings
{
  "persistent": {
    "cluster.routing.allocation.disk.watermark.low":            "85%",   ← no new shards
    "cluster.routing.allocation.disk.watermark.high":           "90%",   ← start relocating
    "cluster.routing.allocation.disk.watermark.flood_stage":    "95%"    ← read-only
  }
}
```

---

## Cross-Cluster Search (CCS)

Query data across multiple independent clusters:

```yaml
# Configure remote cluster
PUT /_cluster/settings
{
  "persistent": {
    "cluster.remote.cluster_b.seeds": ["192.168.2.10:9300"]
  }
}
```

```bash
# Search across clusters
GET /cluster_b:products,local_products/_search
{
  "query": { "match_all": {} }
}
```

---

## Cross-Cluster Replication (CCR)

Replicate indices from a leader cluster to a follower cluster for geo-distribution or DR:

```bash
# On the follower cluster
PUT /replicated-products/_ccr/follow
{
  "remote_cluster":  "leader-cluster",
  "leader_index":    "products"
}
```

---

## Cluster Settings

```bash
# View all cluster settings
GET /_cluster/settings?include_defaults=true

# Update cluster settings
PUT /_cluster/settings
{
  "persistent": {
    "cluster.max_shards_per_node": 2000
  },
  "transient": {
    "indices.recovery.max_bytes_per_sec": "100mb"
  }
}
```

| Type | Persistence |
|---|---|
| `persistent` | Survives cluster restart |
| `transient` | Cleared on cluster restart |

---

## Architecture Patterns

### Single Node (Development)
```
[Node 1: master + data + ingest + coordinating]
```

### Small Production (3 nodes)
```
Node 1: master-eligible + data
Node 2: master-eligible + data
Node 3: master-eligible + data
```

### Medium Production (7 nodes)
```
Node 1-3: dedicated master (no data)
Node 4-7: dedicated data (hot tier)
```

### Large Production (tiered)
```
Node 1-3: dedicated master
Node 4-6: data_hot    (NVMe SSDs)
Node 7-9: data_warm   (SAS/SATA)
Node 10:  data_cold   (S3-backed)
Node 11:  ingest
Node 12:  coordinating
```

---

## Cluster Diagnostics

```bash
# Overall health
GET /_cluster/health

# Allocation explain (why is a shard unassigned?)
GET /_cluster/allocation/explain
{
  "index":   "products",
  "shard":   0,
  "primary": true
}

# Pending cluster tasks
GET /_cluster/pending_tasks

# Hot threads (what are nodes doing?)
GET /_nodes/hot_threads

# Node stats
GET /_nodes/stats

# Task management
GET /_tasks?detailed=true&actions=*search*
```


---

# Chapter 11: Shards & Replicas

## What Are Shards?

A **shard** is a self-contained Lucene index — the fundamental unit of data storage and search in Elasticsearch. When you create an index, its data is split horizontally across multiple shards.

```
Index: "products" (3 primary shards, 1 replica each)

Node 1          Node 2          Node 3
───────         ───────         ───────
[P0] [R1]       [P1] [R2]       [P2] [R0]

P = Primary shard
R = Replica shard
```

---

## Primary vs Replica Shards

### Primary Shards

- Store the actual data
- Handle all **write** operations first
- Number is fixed at **index creation time** (cannot change without reindexing)
- Documents are routed to a specific primary shard by:

```
shard = hash(_id) % number_of_primary_shards
```

### Replica Shards

- Exact copies of primary shards
- Can handle **read** (search) requests — reduces primary load
- Number can be **changed at any time**
- Provide **high availability**: if a primary fails, a replica is promoted

```bash
# Change replica count
PUT /products/_settings
{
  "number_of_replicas": 2
}
```

---

## Choosing the Right Shard Count

This is one of the most important (and hardest) decisions in Elasticsearch.

### Common Rule of Thumb

- Target **20–50 GB per shard** for time-series data
- Target **5–40 GB per shard** for search-heavy indices
- A shard count too high wastes memory; too low limits parallelism

### Shard Sizing Formula

```
Number of primary shards ≈ total_expected_data_size / target_shard_size

Example:
  Expected index size: 200 GB
  Target shard size:   40 GB
  → 5 primary shards
```

### Over-Sharding Problems

Each shard:
- Has overhead: ~few MB of heap per shard even when empty
- Participates in every query (even empty ones)
- Creates coordination overhead

> Modern rule: **Keep shard count as low as needed.** One large shard is often better than ten small ones.

---

## Shard Allocation

The master node decides which shard lives on which node. Key settings:

```bash
# Maximum shards per node (default 1000)
PUT /_cluster/settings
{
  "persistent": {
    "cluster.max_shards_per_node": 2000
  }
}

# Force relocate shard manually
POST /_cluster/reroute
{
  "commands": [{
    "move": {
      "index": "products",
      "shard": 0,
      "from_node": "node-1",
      "to_node":   "node-3"
    }
  }]
}

# Allocate an unassigned shard
POST /_cluster/reroute
{
  "commands": [{
    "allocate_replica": {
      "index":   "products",
      "shard":   1,
      "node":    "node-2"
    }
  }]
}
```

---

## Shard Rebalancing

ES automatically rebalances shards when:
- Nodes join or leave the cluster
- Disk utilization becomes uneven
- Shard counts become uneven across nodes

Control rebalancing:

```bash
PUT /_cluster/settings
{
  "persistent": {
    "cluster.routing.rebalance.enable":           "all",     ← all|primaries|replicas|none
    "cluster.routing.allocation.enable":          "all",     ← control new allocations
    "indices.recovery.max_bytes_per_sec":         "40mb"     ← recovery speed throttle
  }
}
```

---

## Write Path: How a Document Gets Indexed

```
Client → Coordinating Node
           ↓
     Route to primary shard (based on _id hash)
           ↓
     Primary shard writes document
           ↓
     Replicate to all replica shards (in parallel)
           ↓
     Acknowledge to client (after quorum of replicas confirm)
```

### Write Consistency Settings

```bash
PUT /products/_doc/1?wait_for_active_shards=1
# 1      = only primary must confirm
# 2      = primary + 1 replica
# all    = all replicas
# quorum = majority (default for multi-replica)
```

---

## Read Path: How a Search is Executed

```
Client → Coordinating Node
           ↓
     Broadcast query to one copy of each shard
     (primary or any replica — round-robin)
           ↓
     Each shard executes query locally, returns top N results + scores
           ↓
     Coordinating node merges all results, re-ranks, returns top N
           ↓
     Fetch phase: retrieve _source for final top N documents
```

This two-phase approach (Query + Fetch) is called **scatter-gather**.

---

## Replica Promotion (Node Failure)

```
Normal:          P0(node1) → R0(node2)
                 P1(node2) → R1(node3)

node2 fails:     P0(node1) — healthy
                 P1(MISSING) — R1(node3) gets promoted to P1
                 R0(MISSING) — new replica allocated on node1/node3

After recovery:  P0(node1) → R0(node3)
                 P1(node3) → R1(node1)
```

During replica promotion, cluster is `yellow` (missing replicas). Once new replicas are assigned, it returns to `green`.

---

## Segment Merging

Each shard consists of multiple **Lucene segments**. As you index documents:

1. New segments are created on each refresh (default every 1 second)
2. Background merges combine small segments into larger ones
3. Deleted documents are only actually removed during a merge

Too many small segments = slow searches. Force merge on read-only indices:

```bash
# Merge down to 1 segment (best for frozen/read-only indices)
POST /logs-2024-01/_forcemerge?max_num_segments=1
```

---

## Monitoring Shard Health

```bash
# List all shards with status
GET /_cat/shards?v&s=index,shard,prirep

# Check why a shard is unassigned
GET /_cluster/allocation/explain

# Shard store info
GET /products/_shard_stores

# Recovery status
GET /_cat/recovery?v&active_only=true
```

### Shard States

| State | Meaning |
|---|---|
| `STARTED` | Shard is active and serving requests |
| `RELOCATING` | Being moved to another node |
| `INITIALIZING` | Being recovered/allocated |
| `UNASSIGNED` | Not assigned to any node (problem!) |

---

## Tips and Anti-Patterns

| Anti-Pattern | Problem | Solution |
|---|---|---|
| Too many shards | High memory overhead, slow queries | Merge/shrink, use ILM |
| One giant shard | No parallelism | Split across more shards |
| 0 replicas in prod | Single point of failure | At least 1 replica per primary |
| Changing shard count | Not possible on live index | Plan ahead; use reindex to resize |
| Replicas on same node | No HA benefit | Use zone awareness |


---

# Chapter 12: Performance & Optimization

## The Performance Model

Elasticsearch performance has three dimensions to optimize:

```
Indexing Throughput  ←→  Search Latency  ←→  Storage Efficiency
```

Optimizing one often trades off against another. Understand your workload first.

---

## Bulk Indexing Optimization

### Use the Bulk API

Never index documents one at a time for large ingestion:

```bash
# Good: batch 5,000–10,000 documents per request
POST /_bulk
{ "index": { "_index": "logs" } }
{ "message": "event 1", "@timestamp": "2024-03-15T10:00:00Z" }
{ "index": { "_index": "logs" } }
{ "message": "event 2", "@timestamp": "2024-03-15T10:00:01Z" }
...
```

Rule of thumb: **5MB–15MB per bulk request** is a good starting point. Benchmark and tune.

### Disable Refresh During Bulk Ingest

```bash
# Disable auto-refresh before ingestion
PUT /logs/_settings
{ "refresh_interval": "-1" }

# Ingest all data...

# Re-enable and trigger a refresh
PUT /logs/_settings
{ "refresh_interval": "1s" }
POST /logs/_refresh
```

### Set replicas to 0 During Initial Load

```bash
PUT /logs/_settings
{ "number_of_replicas": 0 }

# Load data...

PUT /logs/_settings
{ "number_of_replicas": 1 }
```

### Disable `_source` for Metrics-Only Indices

If you never need to retrieve the original document:

```bash
"mappings": {
  "_source": { "enabled": false }
}
```

### Use `wait_for` Instead of Immediate Refresh

```bash
PUT /products/_doc/1?refresh=wait_for
```

Waits for the next natural refresh rather than forcing one.

---

## Search Performance

### Prefer `filter` Over `query` for Exact Matches

Filter context is cached by Elasticsearch — identical filters across requests hit the cache:

```bash
# Slow (scored, not cached)
"query": { "term": { "status": "active" } }

# Fast (not scored, cached)
"query": {
  "bool": {
    "filter": { "term": { "status": "active" } }
  }
}
```

### Use `keyword` for Aggregations and Sorting

Never aggregate or sort on `text` fields. Use `.keyword` sub-field:

```bash
"aggs": {
  "by_category": {
    "terms": { "field": "category.keyword" }   ← not "category"
  }
}
```

### Avoid `wildcard` and `regexp` on Large Indices

Leading wildcards (`*foo`) force a full term scan. Use `edge_ngram` tokenizer for prefix-style search instead.

### Limit `_source` Fields Returned

Return only what you need:

```bash
"_source": ["name", "price", "category"]
```

### Use `doc_values` for Range Queries on Numbers

`doc_values` (on-disk columnar store) is used automatically for `keyword`, numeric, and date fields. Don't disable it unless you never sort/aggregate that field.

### Query Caching

Elasticsearch caches the results of `filter` context queries in the **filter cache** (also called **node query cache**). Size it appropriately:

```yaml
# elasticsearch.yml
indices.queries.cache.size: 10%    ← % of heap (default 10%)
```

### Shard Request Cache

Aggregation results on entire shards are cached in the **shard request cache**:

```bash
GET /products/_search?request_cache=true
{
  "size": 0,
  "aggs": { ... }
}
```

---

## JVM Heap Tuning

### Set Heap to 50% of RAM (Max 32GB)

```yaml
# jvm.options or ES_JAVA_OPTS
-Xms16g
-Xmx16g
```

**Critical:** always set min and max to the same value to avoid GC pauses.

**Never exceed 32GB** — above this, the JVM loses compressed ordinary object pointers (OOPs), and heap efficiency drops sharply.

### GC Monitoring

```bash
GET /_nodes/stats/jvm

# Key metrics to watch:
# jvm.gc.collectors.old.collection_time_in_millis
# jvm.mem.heap_used_percent
```

Watch for:
- Heap > 85% consistently → add nodes or reduce shard count
- Old-gen GC > 1 minute → memory pressure, investigate large aggs or deep pagination

---

## Indexing Internals: Translog

Every index operation is first written to the **translog** (write-ahead log) for durability:

```yaml
# Async fsync (better write performance, slight data loss risk on crash)
index.translog.durability: async
index.translog.sync_interval: 5s

# Synchronous fsync (default — safe but slower)
index.translog.durability: request
```

---

## Codec: Compression

```bash
PUT /logs
{
  "settings": {
    "index.codec": "best_compression"   ← uses DEFLATE (slower, smaller)
    # default is "default" (LZ4 — faster, larger)
  }
}
```

Use `best_compression` for cold/warm indices where disk is constrained.

---

## Mappings for Performance

### Disable Norms for Non-Relevance Fields

Norms store per-document scoring info. Disable for fields where relevance doesn't matter:

```bash
"status": {
  "type":  "text",
  "norms": false
}
```

### Disable `_all` Field (ES < 6.0)

In older versions, disable the `_all` field to save space. Modern ES doesn't have it.

### Index Only What You Search

```bash
"internal_metadata": {
  "type":  "object",
  "enabled": false      ← stored in _source but not indexed
}
```

### Use `eager_global_ordinals` for High-Frequency Aggregation Fields

```bash
"category": {
  "type":                  "keyword",
  "eager_global_ordinals": true     ← precompute ordinals at refresh time
}
```

---

## Thread Pool Monitoring

```bash
GET /_cat/thread_pool?v&h=name,active,queue,rejected

# Key thread pools:
# write   → indexing and bulk
# search  → query execution
# get     → GET by ID
```

If `rejected` is non-zero, the queue is full — you're overloading the node. Options:
- Add more nodes
- Throttle ingest rate
- Reduce bulk request size

---

## Slow Log

Log queries/indexing that take too long:

```bash
PUT /products/_settings
{
  "index.search.slowlog.threshold.query.warn":  "10s",
  "index.search.slowlog.threshold.query.info":  "5s",
  "index.search.slowlog.threshold.fetch.warn":  "1s",
  "index.indexing.slowlog.threshold.index.warn": "10s"
}
```

Check logs in: `/var/log/elasticsearch/<cluster_name>_index_search_slowlog.log`

---

## Performance Checklist

| Area | Check |
|---|---|
| Indexing | Use Bulk API with 5MB–15MB batches |
| Indexing | Disable `refresh_interval` during initial load |
| Indexing | 0 replicas during initial load, then restore |
| Search | Use `filter` not `query` for exact/range filters |
| Search | Avoid deep pagination — use `search_after` |
| Search | Limit `_source` fields returned |
| Mapping | Use `keyword` for aggregations, not `text` |
| Mapping | Disable `norms` on non-text-relevance fields |
| JVM | Heap = 50% RAM, max 32GB, Xms=Xmx |
| Shards | 20–50GB per shard, minimize shard count |
| Hardware | SSDs for hot tier, spinning disks for warm/cold |


---

# Chapter 13: Security & Authentication

## Security Overview

Elasticsearch 8.x enables security by default. In ES 7.x and earlier, security (X-Pack) was optional.

Security layers:
```
1. Transport Layer Security (TLS)   ← encrypt node-to-node and client-to-node traffic
2. Authentication                   ← verify who the user is
3. Authorization (RBAC)             ← control what they can do
4. Audit Logging                    ← track who did what
5. Field/Document-Level Security    ← restrict which data users can see
```

---

## Enable Security (ES 8.x Default)

In ES 8.x, security is on by default when you start a new cluster. For development you may disable it:

```yaml
# elasticsearch.yml — development only
xpack.security.enabled: false
```

For production, configure properly:

```yaml
xpack.security.enabled: true
xpack.security.transport.ssl.enabled: true
xpack.security.transport.ssl.verification_mode: certificate
xpack.security.transport.ssl.keystore.path: elastic-certificates.p12

xpack.security.http.ssl.enabled: true
xpack.security.http.ssl.keystore.path: http.p12
```

---

## Generate TLS Certificates

```bash
# Generate CA
bin/elasticsearch-certutil ca

# Generate node certificate signed by that CA
bin/elasticsearch-certutil cert --ca elastic-stack-ca.p12

# For HTTP layer
bin/elasticsearch-certutil http

# Copy certificates to all nodes
cp elastic-certificates.p12 /etc/elasticsearch/
```

---

## Built-in Users

Elasticsearch ships with several reserved system users:

| User | Purpose |
|---|---|
| `elastic` | Superuser — full access |
| `kibana_system` | Used by Kibana to connect |
| `logstash_system` | Used by Logstash |
| `beats_system` | Used by Beats |
| `apm_system` | Used by APM Server |
| `remote_monitoring_user` | Monitoring |

```bash
# Set passwords for built-in users
bin/elasticsearch-setup-passwords interactive

# Or auto-generate
bin/elasticsearch-setup-passwords auto
```

---

## User Management

### Create a User

```bash
POST /_security/user/john_doe
{
  "password":  "SecurePass123!",
  "roles":     ["data_analyst"],
  "full_name": "John Doe",
  "email":     "john@example.com"
}
```

### Change Password

```bash
POST /_security/user/john_doe/_password
{
  "password": "NewSecurePass456!"
}
```

### Disable/Enable User

```bash
POST /_security/user/john_doe/_disable
POST /_security/user/john_doe/_enable
```

---

## Role-Based Access Control (RBAC)

### Built-In Roles

| Role | Access Level |
|---|---|
| `superuser` | Full access to everything |
| `kibana_admin` | Full Kibana administration |
| `viewer` | Read-only on all indices/spaces |
| `editor` | Read/write on all indices |
| `monitoring_user` | Read cluster monitoring data |
| `snapshot_user` | Create/restore snapshots |

### Create a Custom Role

```bash
PUT /_security/role/product_reader
{
  "cluster": ["monitor"],
  "indices": [
    {
      "names":      ["products*"],
      "privileges": ["read", "view_index_metadata"]
    }
  ]
}
```

### Full Admin Role for Specific Index

```bash
PUT /_security/role/orders_admin
{
  "cluster": ["monitor", "manage_index_templates"],
  "indices": [
    {
      "names":      ["orders*"],
      "privileges": ["all"]
    }
  ]
}
```

### Index Privileges Reference

| Privilege | Description |
|---|---|
| `read` | Search and GET |
| `write` | Index, update, delete documents |
| `create` | Index new documents only |
| `delete` | Delete documents |
| `index` | Index and update documents |
| `all` | All operations |
| `view_index_metadata` | View mappings, settings, aliases |
| `manage` | Manage index (settings, mappings, aliases) |
| `monitor` | Read index stats |

### Cluster Privileges

| Privilege | Description |
|---|---|
| `monitor` | Read cluster stats/health |
| `manage` | Manage cluster settings |
| `manage_index_templates` | Manage templates |
| `manage_ilm` | Manage ILM policies |
| `manage_security` | Manage users/roles |
| `all` | Everything |

---

## API Keys

For programmatic access without a username/password:

```bash
# Create an API key
POST /_security/api_key
{
  "name":       "my-app-key",
  "expiration": "30d",
  "role_descriptors": {
    "app_role": {
      "cluster": ["monitor"],
      "indices": [
        {
          "names":      ["products*"],
          "privileges": ["read"]
        }
      ]
    }
  }
}
```

Response:
```json
{
  "id":          "abc123",
  "name":        "my-app-key",
  "api_key":     "xyz789secret",
  "encoded":     "base64_of_id:api_key"
}
```

Use in requests:
```bash
curl -H "Authorization: ApiKey base64_of_id:api_key" \
     http://localhost:9200/products/_search
```

### Manage API Keys

```bash
# List your keys
GET /_security/api_key?mine=true

# Invalidate a key
DELETE /_security/api_key
{ "ids": ["abc123"] }
```

---

## Document-Level Security (DLS)

Restrict which documents a role can see:

```bash
PUT /_security/role/us_customer_service
{
  "indices": [
    {
      "names":      ["orders*"],
      "privileges": ["read"],
      "query": {
        "term": { "region": "US" }
      }
    }
  ]
}
```

Users with this role only see orders with `region: "US"`.

---

## Field-Level Security (FLS)

Restrict which fields a role can see:

```bash
PUT /_security/role/read_without_pii
{
  "indices": [
    {
      "names":      ["users*"],
      "privileges": ["read"],
      "field_security": {
        "grant":  ["*"],            ← allow all
        "except": ["ssn", "dob", "payment_info"]   ← except these
      }
    }
  ]
}
```

---

## Audit Logging

Track all access and changes:

```yaml
# elasticsearch.yml
xpack.security.audit.enabled: true
xpack.security.audit.logfile.events.include: [
  "authentication_success",
  "authentication_failed",
  "access_granted",
  "access_denied",
  "tampered_request",
  "connection_denied"
]
```

Audit logs go to `<cluster_name>_audit.json` in the logs directory.

---

## Security Best Practices

| Practice | Details |
|---|---|
| Always use TLS | For both transport and HTTP layers |
| Rotate API keys | Set expiration; rotate regularly |
| Use minimal privileges | Grant only what each service needs |
| Separate service accounts | One API key per service, not shared |
| Disable `elastic` superuser | Create named admin users instead |
| Enable audit logging | Required for compliance |
| Use DLS/FLS for multi-tenancy | Isolate tenants with row/column security |
| Network-level firewall | Don't expose port 9200 to the internet |


---

# Chapter 14: Monitoring & Observability

## What to Monitor

```
Cluster Health     →  Is the cluster green/yellow/red?
Node Resources     →  CPU, memory, disk, network
Indexing Rate      →  Documents/sec, bulk rejections
Search Rate        →  Queries/sec, query latency
JVM               →  Heap usage, GC frequency/duration
Shard Health       →  Unassigned shards, recovery progress
```

---

## Cluster Health API

```bash
GET /_cluster/health
```

```json
{
  "cluster_name":             "production",
  "status":                   "green",
  "number_of_nodes":          5,
  "number_of_data_nodes":     3,
  "active_primary_shards":    30,
  "active_shards":            60,
  "relocating_shards":        0,
  "initializing_shards":      0,
  "unassigned_shards":        0,
  "delayed_unassigned_shards": 0,
  "pending_tasks":            0
}
```

| Status | Action Required |
|---|---|
| `green` | All shards assigned and healthy |
| `yellow` | Missing replicas — investigate but not critical |
| `red` | Missing primaries — data may be unavailable — urgent |

```bash
# Health per index
GET /_cluster/health?level=indices

# Health per shard
GET /_cluster/health?level=shards
```

---

## `_cat` APIs (Human-Readable)

The `_cat` APIs output tabular data, great for quick checks:

```bash
# Nodes overview
GET /_cat/nodes?v&h=name,ip,role,heap.percent,cpu,load_1m,disk.avail

# Indices with stats
GET /_cat/indices?v&s=index&h=index,health,status,pri,rep,docs.count,store.size

# Shards
GET /_cat/shards?v&s=index,shard

# Pending recovery
GET /_cat/recovery?v&active_only=true

# Thread pool queues and rejections
GET /_cat/thread_pool?v&h=name,active,queue,rejected

# Pending master tasks
GET /_cat/pending_tasks?v

# Segment info
GET /_cat/segments?v
```

---

## Node Stats API

```bash
GET /_nodes/stats

# Specific subsections
GET /_nodes/stats/jvm,indices,os,process
```

Key sections:

### JVM Stats

```json
"jvm": {
  "mem": {
    "heap_used_in_bytes":    4294967296,
    "heap_used_percent":     53,
    "heap_max_in_bytes":     8589934592
  },
  "gc": {
    "collectors": {
      "young": { "collection_count": 120, "collection_time_in_millis": 3400 },
      "old":   { "collection_count": 2,   "collection_time_in_millis": 800  }
    }
  }
}
```

**Alert thresholds:**
- Heap > 85% → memory pressure
- Old GC time > 30s/hour → memory leak or undersized heap

### Indices Stats

```json
"indices": {
  "indexing": {
    "index_total":         5000000,
    "index_time_in_millis": 12000,
    "index_failed":        0
  },
  "search": {
    "query_total":          250000,
    "query_time_in_millis": 18000,
    "fetch_total":          245000
  }
}
```

---

## Index-Level Stats

```bash
GET /products/_stats

# Specific stats
GET /products/_stats/indexing,search,store
```

---

## Slow Logs

Identify slow queries and indexing:

```bash
PUT /products/_settings
{
  "index.search.slowlog.threshold.query.warn":  "5s",
  "index.search.slowlog.threshold.query.info":  "1s",
  "index.search.slowlog.threshold.query.debug": "500ms",
  "index.search.slowlog.threshold.fetch.warn":  "1s",
  "index.indexing.slowlog.threshold.index.warn": "5s",
  "index.indexing.slowlog.level":               "info"
}
```

Slow logs appear in: `<data-dir>/logs/<cluster-name>_index_search_slowlog.log`

---

## Task Management API

```bash
# List all running tasks
GET /_tasks?detailed=true

# Filter by type
GET /_tasks?actions=*search*&detailed=true

# Cancel a long-running task
POST /_tasks/<task-id>/_cancel
```

---

## Hot Threads

Find what CPU is actually doing:

```bash
GET /_nodes/hot_threads

GET /_nodes/node-1/hot_threads?threads=5&interval=1s
```

---

## Stack Monitoring (Kibana)

For production clusters, use Elastic's built-in monitoring UI:

1. In Kibana: **Stack Management → Stack Monitoring**
2. Or enable Metricbeat to collect ES metrics:

```yaml
# metricbeat.yml
metricbeat.modules:
  - module: elasticsearch
    metricsets:
      - node
      - node_stats
      - cluster_stats
      - index
      - index_recovery
      - shard
    period: 10s
    hosts: ["http://localhost:9200"]
```

---

## Alerting with Kibana Watcher

### Watch Example: Alert on Red Cluster Status

```bash
PUT /_watcher/watch/cluster_health_watch
{
  "trigger": {
    "schedule": { "interval": "1m" }
  },
  "input": {
    "http": {
      "request": {
        "host": "localhost",
        "port": 9200,
        "path": "/_cluster/health"
      }
    }
  },
  "condition": {
    "compare": { "ctx.payload.status": { "eq": "red" } }
  },
  "actions": {
    "send_email": {
      "email": {
        "to":      ["ops-team@company.com"],
        "subject": "ALERT: Elasticsearch cluster is RED",
        "body":    "Cluster {{ctx.payload.cluster_name}} status: {{ctx.payload.status}}"
      }
    }
  }
}
```

---

## Key Metrics to Dashboard

| Metric | Source | Alert Threshold |
|---|---|---|
| Cluster status | `_cluster/health` | != green |
| Heap used % | `_nodes/stats/jvm` | > 85% |
| Unassigned shards | `_cluster/health` | > 0 |
| GC old collection time | `_nodes/stats/jvm` | > 30s/min |
| Search query latency | `_nodes/stats/indices` | > 500ms p95 |
| Indexing failures | `_nodes/stats/indices` | > 0 |
| Disk watermark | `_cat/nodes` | > 80% |
| Thread pool rejections | `_cat/thread_pool` | > 0 |
| CPU usage | `_nodes/stats/os` | > 90% sustained |

---

## Diagnostic Runbook

### Yellow Status (Missing Replicas)

```bash
# 1. Find unassigned shards
GET /_cat/shards?v&h=index,shard,prirep,state,unassigned.reason

# 2. Explain why
GET /_cluster/allocation/explain

# 3. Common causes:
# - Not enough nodes (replicas > data nodes - 1)
# - Disk watermark exceeded
# - Shard filter rules preventing allocation
```

### Red Status (Missing Primary)

```bash
# Emergency: find which indices are red
GET /_cat/indices?v&health=red

# Last resort: force allocate a stale primary (data loss possible!)
POST /_cluster/reroute
{
  "commands": [{
    "allocate_stale_primary": {
      "index":           "products",
      "shard":           0,
      "node":            "node-1",
      "accept_data_loss": true
    }
  }]
}
```

### High Heap Usage

```bash
# 1. Check what's using heap
GET /_nodes/stats/jvm

# 2. Check for large aggregations or deep pagination
GET /_tasks?detailed=true&actions=*search*

# 3. Force GC (use with caution, causes pauses)
POST /_nodes/_all/_gc
```


---

# Chapter 15: Elasticsearch with Kibana

## What is Kibana?

**Kibana** is the visualization and management layer of the Elastic Stack. It provides:

- **Discover** — explore and search your data
- **Dashboards** — create visual analytics
- **Lens** — drag-and-drop chart builder
- **Dev Tools** — interactive REST API console
- **Stack Management** — index patterns, ILM, security, templates
- **Observability** — APM, logs, metrics, uptime
- **SIEM/Security** — threat detection and investigation

---

## Connecting Kibana to Elasticsearch

```yaml
# kibana.yml
server.port: 5601
server.host: "0.0.0.0"
elasticsearch.hosts: ["http://localhost:9200"]

# If security is enabled
elasticsearch.username: "kibana_system"
elasticsearch.password: "kibana_system_password"
```

---

## Dev Tools Console

The **Dev Tools Console** is the easiest way to run Elasticsearch queries interactively.

Access: `http://localhost:5601` → Dev Tools → Console

```
# Auto-complete is supported — press Ctrl+Space
GET _cluster/health

POST /products/_search
{
  "query": { "match": { "name": "laptop" } }
}
```

### Keyboard Shortcuts

| Shortcut | Action |
|---|---|
| `Ctrl+Enter` | Execute request |
| `Ctrl+/` | Toggle comment |
| `Ctrl+I` | Auto-indent |
| `Ctrl+Up/Down` | Move between requests |

---

## Index Patterns (Data Views in ES 8.x)

A **Data View** (formerly Index Pattern) tells Kibana which indices to search and which field is the time field.

### Create via UI
1. Stack Management → Data Views → Create Data View
2. Pattern: `logs-*` or `products`
3. Time field: `@timestamp` (or none for non-time-series)

### Create via API

```bash
POST /api/data_views/data_view
{
  "data_view": {
    "title":      "logs-*",
    "timeFieldName": "@timestamp",
    "name":       "Application Logs"
  }
}
```

---

## Discover

The **Discover** page lets you:
- Search documents using KQL or Lucene syntax
- Filter by time range
- Inspect field values
- Add columns to the document table

### Kibana Query Language (KQL)

KQL is Kibana's simplified query syntax for Discover and dashboard filters:

```
# Free-text search
error

# Field:value
status: "active"

# Range
price >= 100 and price <= 500

# Wildcard
email: *@gmail.com

# Negation
not status: "cancelled"

# Nested
order.items.product_id: "SKU-123"

# Boolean combinations
(status: "active" or status: "pending") and region: "EU"
```

---

## Dashboards

### Creating a Dashboard

1. Go to **Dashboards** → Create dashboard
2. Click **Create visualization** or **Add from library**
3. Use **Lens** (drag-and-drop) or other editors

### Lens Chart Types

| Type | Use Case |
|---|---|
| Bar / Horizontal Bar | Compare categories |
| Line | Trends over time |
| Area | Stacked trends |
| Pie / Donut | Proportions |
| Data table | Tabular aggregations |
| Metric | Single number KPI |
| Heatmap | Density/correlation |
| Gauge | Progress toward goal |
| Map | Geo data |

### Dashboard Controls

- **Time filter** — set global time range
- **Filters** — add field:value filters that apply to all panels
- **Controls** — dropdown/range sliders for interactive filtering

---

## Lens: Drag-and-Drop Visualizations

Lens auto-suggests chart types based on your data.

### Example: Average Price by Category Over Time

1. Open Lens
2. Set index: `products-*`
3. X-axis: `created_at` (date histogram)
4. Y-axis: `price` → Average
5. Break down by: `category.keyword`
6. Choose chart type: Line

### Example: Top 10 Products by Revenue

1. Horizontal bar chart
2. X-axis: `total` sum
3. Y-axis: `product_name.keyword` (terms, top 10)

---

## Canvas

Canvas provides pixel-perfect custom presentations:

- Use workpads to create presentations/reports
- Embed live Elasticsearch data via elements
- Write Elasticsearch SQL in expression editor

---

## Maps (Geo Data)

The **Maps** app visualizes geo data:

- Document layer (geo_point fields)
- Choropleth (fill regions by value)
- Heat map
- Point-to-point connections

---

## Saved Searches and Filters

In Discover, save a search and reuse it:

1. Search with filters active
2. Click **Save** → name your search
3. Use it in dashboards as a data source

---

## Kibana Alerting

Create alerts based on query thresholds:

1. **Stack Management → Rules → Create rule**
2. Rule types:
   - Elasticsearch query — run a query, alert if results hit threshold
   - Index threshold — alert when doc count crosses a level
   - Metrics threshold — CPU/memory based
3. Actions: email, Slack, PagerDuty, webhook

### Example: Alert on Error Rate

```
Trigger:    Count of log_level:error in last 5 minutes > 100
Action:     Send Slack message to #alerts
```

---

## Index Management via Kibana

Stack Management → Index Management:

- View all indices with health, size, doc count
- Create/delete indices
- Manage index templates
- Create ILM policies (visual policy editor)
- Manage rollover

---

## Kibana Spaces

Spaces isolate dashboards, data views, and saved objects per team:

- Each space has independent dashboards, visualizations, index patterns
- Users can be assigned to specific spaces
- Useful for multi-tenant environments

```bash
# Create a space via API
POST /api/spaces/space
{
  "id":          "marketing",
  "name":        "Marketing Team",
  "description": "Marketing analytics dashboards"
}
```

---

## Useful Kibana APIs

```bash
# List all saved objects
GET /api/saved_objects/_find?type=dashboard

# Export dashboard and all dependencies
POST /api/saved_objects/_export
{
  "type": "dashboard",
  "objects": [{ "type": "dashboard", "id": "<dashboard-id>" }],
  "includeReferencesDeep": true
}

# Import
POST /api/saved_objects/_import?overwrite=true
# (with ndjson file body)
```

---

## Kibana Tips

| Tip | Details |
|---|---|
| Use time range wisely | Narrow time ranges = faster Discover queries |
| Pin filters | Persistent filters survive page navigation |
| Use Inspect | Each panel has "Inspect" to see the raw ES query |
| Shared links | Share dashboards with `?_g` and `?_a` URL params for deep links |
| Dark mode | Stack Management → Advanced Settings → theme:darkMode |


---

# Chapter 16: The ELK Stack

## What is the ELK Stack?

ELK stands for three open-source projects:

```
E — Elasticsearch   (store, search, analyze)
L — Logstash        (collect, transform, enrich)
K — Kibana          (visualize, explore, alert)
```

Now more commonly called the **Elastic Stack**, which also includes **Beats** (lightweight shippers):

```
Data Sources
    ↓
Beats (Filebeat, Metricbeat, Packetbeat...)
    ↓ (optional)
Logstash (transform, enrich, parse)
    ↓
Elasticsearch (store + search)
    ↓
Kibana (visualize)
```

---

## Beats Family

Beats are **lightweight data shippers** — small agents installed on your servers. Each type specializes in one data source.

| Beat | Data Source |
|---|---|
| **Filebeat** | Log files (`/var/log/*.log`, Docker logs) |
| **Metricbeat** | System and service metrics (CPU, memory, MySQL) |
| **Packetbeat** | Network traffic analysis |
| **Heartbeat** | Uptime and availability monitoring |
| **Auditbeat** | Audit events (Linux audit framework) |
| **Winlogbeat** | Windows event logs |

### Filebeat Example

```yaml
# filebeat.yml
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - /var/log/app/*.log
      - /var/log/nginx/access.log
    fields:
      service: my-app
      environment: production

output.elasticsearch:
  hosts: ["http://elasticsearch:9200"]
  index:  "logs-%{+yyyy.MM.dd}"

# Or output to Logstash for processing
output.logstash:
  hosts: ["logstash:5044"]

setup.kibana:
  host: "http://kibana:5601"
```

```bash
# Start filebeat
sudo systemctl start filebeat

# Test config
filebeat test config -e

# Test output connection
filebeat test output -e
```

### Metricbeat Example

```yaml
# metricbeat.yml
metricbeat.modules:
  - module: system
    metricsets: [cpu, memory, disk, network, process]
    period: 10s

  - module: elasticsearch
    metricsets: [node, node_stats, cluster_stats]
    period: 30s
    hosts: ["http://elasticsearch:9200"]

  - module: nginx
    metricsets: [stubstatus]
    period: 10s
    hosts: ["http://localhost/server-status"]

output.elasticsearch:
  hosts: ["http://elasticsearch:9200"]
```

---

## Logstash

**Logstash** is a server-side data processing pipeline:

```
Input(s)  →  Filter(s)  →  Output(s)
```

### Logstash Pipeline Configuration

```bash
# /etc/logstash/conf.d/app-logs.conf

input {
  beats {
    port => 5044
  }
  file {
    path  => "/var/log/app/*.log"
    start_position => "beginning"
  }
}

filter {
  # Parse unstructured log lines
  grok {
    match => {
      "message" => '%{TIMESTAMP_ISO8601:timestamp} %{LOGLEVEL:level} %{GREEDYDATA:msg}'
    }
  }

  # Parse timestamp
  date {
    match => ["timestamp", "ISO8601"]
    target => "@timestamp"
    remove_field => ["timestamp"]
  }

  # Parse JSON message bodies
  json {
    source  => "msg"
    target  => "parsed"
    remove_field => ["msg"]
  }

  # Enrich with GeoIP
  geoip {
    source => "client_ip"
    target => "geoip"
  }

  # Add computed fields
  mutate {
    add_field  => { "environment" => "production" }
    rename     => { "host" => "server_hostname" }
    remove_field => ["agent", "ecs"]
  }

  # Drop health check noise
  if [request_path] == "/health" {
    drop {}
  }
}

output {
  elasticsearch {
    hosts  => ["http://elasticsearch:9200"]
    index  => "logs-%{+YYYY.MM.dd}"
    user   => "logstash_system"
    password => "password"
  }

  # Debug output during development
  stdout {
    codec => rubydebug
  }
}
```

---

## Logstash Grok Patterns

Grok is the most common Logstash filter for parsing unstructured text.

### Common Patterns

| Pattern | Matches |
|---|---|
| `%{IP:client_ip}` | IP address → `client_ip` field |
| `%{TIMESTAMP_ISO8601:timestamp}` | ISO 8601 timestamp |
| `%{LOGLEVEL:level}` | ERROR, WARN, INFO, DEBUG |
| `%{WORD:method}` | Single word (GET, POST) |
| `%{NUMBER:bytes}` | Numeric value |
| `%{URIPATH:path}` | URL path |
| `%{GREEDYDATA:message}` | Everything remaining |

### Apache Access Log Example

```
filter {
  grok {
    match => {
      "message" => '%{IPORHOST:client_ip} %{USER:ident} %{USER:auth} \[%{HTTPDATE:timestamp}\] "%{WORD:method} %{URIPATH:request_path}(?: HTTP/%{NUMBER:http_version})?" %{NUMBER:response_code} (?:%{NUMBER:bytes}|-)'
    }
  }
}
```

Test grok patterns at: [grokdebugger.com](https://grokdebugger.com) or Kibana Dev Tools → Grok Debugger.

---

## Ingest Pipelines (Alternative to Logstash)

For simpler transformations, use **Elasticsearch Ingest Pipelines** — no separate Logstash needed.

```bash
PUT /_ingest/pipeline/parse-apache-logs
{
  "processors": [
    {
      "grok": {
        "field": "message",
        "patterns": [
          "%{IPORHOST:client_ip} .* \\[%{HTTPDATE:timestamp}\\] \"%{WORD:method} %{URIPATH:path}\" %{NUMBER:status_code:int}"
        ]
      }
    },
    {
      "date": {
        "field":  "timestamp",
        "formats": ["dd/MMM/yyyy:HH:mm:ss Z"]
      }
    },
    {
      "geoip": {
        "field":  "client_ip",
        "target_field": "geo"
      }
    },
    {
      "remove": {
        "field": ["message", "timestamp"]
      }
    }
  ]
}

# Use the pipeline when indexing
POST /apache-logs/_doc?pipeline=parse-apache-logs
{
  "message": "192.168.1.1 - - [01/Jan/2024:12:00:00 +0000] \"GET /api/v1/products HTTP/1.1\" 200 1234"
}
```

---

## Common ELK Architectures

### Simple: Beats → ES → Kibana

```
App Servers
└── Filebeat → Elasticsearch → Kibana
```

Best for: small teams, low volume, simple parsing.

### Production: Beats → Logstash → ES → Kibana

```
App Servers
└── Filebeat ──┐
               ▼
            Logstash (parse, enrich, route)
               ↓
            Elasticsearch
               ↓
            Kibana
```

Best for: complex parsing, multiple sources, transformation.

### High-Volume: Beats → Kafka → Logstash → ES

```
App Servers
└── Filebeat → Kafka (buffer) → Logstash cluster → Elasticsearch
```

Best for: very high volume where Logstash can't keep up; Kafka buffers spikes.

---

## Elastic Common Schema (ECS)

ECS is a standardized field naming convention for all Elastic data:

| ECS Field | Description |
|---|---|
| `@timestamp` | Event timestamp |
| `event.action` | What happened (login, query) |
| `event.outcome` | success / failure |
| `source.ip` | Source IP address |
| `destination.ip` | Destination IP |
| `user.name` | Username |
| `host.name` | Hostname |
| `service.name` | Service name |
| `log.level` | Log severity |
| `http.request.method` | HTTP method |
| `http.response.status_code` | HTTP status |
| `url.path` | URL path |
| `error.message` | Error description |

Using ECS allows all Kibana dashboards, detections, and ML to work out of the box.


---

# Chapter 17: Snapshots & Backup

## What Are Snapshots?

A **snapshot** is an incremental backup of one or more Elasticsearch indices (or the entire cluster state) stored in a repository. Snapshots are:

- **Incremental** — only changed segments since last snapshot are copied
- **Non-destructive** — the cluster continues operating during a snapshot
- **Restorable** — indices can be restored to any cluster

---

## Snapshot Repositories

A **repository** is where snapshots are stored. Supported types:

| Type | Backend | Plugin Required |
|---|---|---|
| `fs` | Shared filesystem (NFS, etc.) | No |
| `s3` | Amazon S3 | `repository-s3` |
| `gcs` | Google Cloud Storage | `repository-gcs` |
| `azure` | Azure Blob Storage | `repository-azure` |
| `hdfs` | Hadoop HDFS | `repository-hdfs` |
| `url` | Read-only HTTP/S URL | No |

---

## Register a Repository

### Filesystem Repository (Development)

```bash
# First, add the path to elasticsearch.yml on ALL nodes:
# path.repo: ["/mnt/backups"]

PUT /_snapshot/my_backup_repo
{
  "type": "fs",
  "settings": {
    "location": "/mnt/backups/elasticsearch",
    "compress": true
  }
}
```

### S3 Repository

```bash
# Install plugin first: bin/elasticsearch-plugin install repository-s3

# Store credentials in keystore
bin/elasticsearch-keystore add s3.client.default.access_key
bin/elasticsearch-keystore add s3.client.default.secret_key

PUT /_snapshot/s3_backup
{
  "type": "s3",
  "settings": {
    "bucket":   "my-es-backups",
    "region":   "us-east-1",
    "base_path": "elasticsearch/snapshots",
    "compress": true
  }
}
```

### Verify Repository

```bash
POST /_snapshot/my_backup_repo/_verify
```

---

## Taking Snapshots

### Manual Snapshot

```bash
# Snapshot all indices
PUT /_snapshot/my_backup_repo/snapshot_2024_03_15
{
  "indices": "*",
  "ignore_unavailable": true,
  "include_global_state": false
}

# Snapshot specific indices
PUT /_snapshot/my_backup_repo/products_snapshot
{
  "indices": ["products", "products_v2"],
  "ignore_unavailable": true,
  "include_global_state": false,
  "metadata": {
    "taken_by":   "ops-team",
    "reason":     "pre-migration backup"
  }
}
```

### Check Snapshot Status

```bash
# Status of in-progress snapshot
GET /_snapshot/my_backup_repo/snapshot_2024_03_15/_status

# Details of completed snapshot
GET /_snapshot/my_backup_repo/snapshot_2024_03_15
```

### List All Snapshots

```bash
GET /_snapshot/my_backup_repo/*?verbose=false
```

---

## Snapshot Lifecycle Management (SLM)

Automate snapshot creation and retention:

```bash
PUT /_slm/policy/daily_snapshots
{
  "schedule":   "0 30 1 * * ?",    ← cron: 1:30 AM daily
  "name":       "<daily-snap-{now/d}>",
  "repository": "s3_backup",
  "config": {
    "indices":              ["*"],
    "ignore_unavailable":   true,
    "include_global_state": false
  },
  "retention": {
    "expire_after":   "30d",    ← delete snapshots older than 30 days
    "min_count":      5,        ← always keep at least 5
    "max_count":      50        ← never keep more than 50
  }
}
```

### Execute SLM Policy Immediately

```bash
POST /_slm/policy/daily_snapshots/_execute
```

### Check SLM Status

```bash
GET /_slm/policy/daily_snapshots

GET /_slm/stats
```

---

## Restoring Snapshots

### Restore All Indices

```bash
POST /_snapshot/my_backup_repo/snapshot_2024_03_15/_restore
{
  "ignore_unavailable":   true,
  "include_global_state": false
}
```

### Restore Specific Indices

```bash
POST /_snapshot/my_backup_repo/snapshot_2024_03_15/_restore
{
  "indices": ["products", "orders"],
  "ignore_unavailable": true
}
```

### Restore with Rename (to avoid conflicts with existing indices)

```bash
POST /_snapshot/my_backup_repo/snapshot_2024_03_15/_restore
{
  "indices":             ["products"],
  "rename_pattern":      "(.+)",
  "rename_replacement":  "restored_$1",
  "index_settings": {
    "index.number_of_replicas": 0
  }
}
```

This restores `products` as `restored_products`.

### Monitor Restore Progress

```bash
GET /_cat/recovery?v&active_only=true
GET /products/_recovery
```

---

## Searchable Snapshots

Searchable snapshots allow you to **search index data directly from the snapshot repository** without fully restoring it. Ideal for cold/frozen data.

```bash
POST /_snapshot/s3_backup/snapshot_2024_03_15/_mount
{
  "index":           "logs-2024-01",
  "renamed_index":   "logs-2024-01-frozen",
  "storage":         "full_copy"    ← or "shared_cache" for frozen tier
}
```

- `full_copy` — downloads all data to local storage (warm tier)
- `shared_cache` — fetches data on demand from repository (frozen tier)

---

## Cross-Cluster Restore

Restore from a snapshot taken on a different cluster:

```bash
# Register the remote cluster's repository
PUT /_snapshot/remote_backup
{
  "type": "s3",
  "settings": {
    "bucket":   "remote-cluster-backups",
    "readonly": true    ← read-only to prevent accidental writes
  }
}

# Restore
POST /_snapshot/remote_backup/snapshot_2024_03_15/_restore
{
  "indices": ["products"]
}
```

---

## Backup Strategy

### Recommended Production Backup Schedule

| Snapshot Type | Frequency | Retention |
|---|---|---|
| Daily | Every day at 2 AM | 30 days |
| Weekly | Every Sunday | 12 weeks |
| Monthly | 1st of each month | 12 months |

### What `include_global_state` Captures

When `true`, snapshots also include:
- Index templates
- Component templates
- ILM policies
- SLM policies
- Feature state (security users/roles, etc.)

Always set to `false` unless you're doing a full cluster migration.

---

## Delete Snapshots

```bash
# Delete a single snapshot
DELETE /_snapshot/my_backup_repo/snapshot_2024_03_15

# Delete multiple
DELETE /_snapshot/my_backup_repo/snapshot_2024_03_15,snapshot_2024_03_14

# Delete all in repo (careful!)
DELETE /_snapshot/my_backup_repo/*
```


---

# Chapter 18: Advanced Search

## Fuzzy Search

Fuzzy matching allows typos and misspellings. Elasticsearch uses **Levenshtein edit distance**.

```bash
GET /products/_search
{
  "query": {
    "fuzzy": {
      "name": {
        "value":       "labtop",     ← user typed "labtop"
        "fuzziness":   "AUTO",       ← AUTO: 0 for 1-2 chars, 1 for 3-5, 2 for 6+
        "prefix_length": 2,          ← first N chars must match exactly
        "max_expansions": 50
      }
    }
  }
}
```

### Fuzzy in `match` Query (Most Common)

```bash
{
  "query": {
    "match": {
      "description": {
        "query":      "wireles headpones",
        "fuzziness":  "AUTO",
        "operator":   "and"
      }
    }
  }
}
```

---

## Geo Queries

### `geo_distance` — Within a Radius

```bash
GET /restaurants/_search
{
  "query": {
    "bool": {
      "must": { "match_all": {} },
      "filter": {
        "geo_distance": {
          "distance":  "5km",
          "location":  { "lat": 40.7128, "lon": -74.0060 }
        }
      }
    }
  },
  "sort": [
    {
      "_geo_distance": {
        "location":     { "lat": 40.7128, "lon": -74.0060 },
        "order":        "asc",
        "unit":         "km",
        "distance_type": "plane"
      }
    }
  ]
}
```

### `geo_bounding_box` — Within a Rectangle

```bash
{
  "filter": {
    "geo_bounding_box": {
      "location": {
        "top_left":     { "lat": 42.0, "lon": -75.0 },
        "bottom_right": { "lat": 40.0, "lon": -70.0 }
      }
    }
  }
}
```

### `geo_shape` — Complex Shapes

```bash
{
  "filter": {
    "geo_shape": {
      "location": {
        "shape": {
          "type":        "polygon",
          "coordinates": [[
            [-74.1, 40.7], [-73.9, 40.7],
            [-73.9, 40.8], [-74.1, 40.8],
            [-74.1, 40.7]
          ]]
        },
        "relation": "within"   ← within|intersects|contains|disjoint
      }
    }
  }
}
```

---

## Nested Queries

For `nested` field type — query across fields of the same array element:

```bash
GET /orders/_search
{
  "query": {
    "nested": {
      "path":       "line_items",
      "query": {
        "bool": {
          "must": [
            { "term":  { "line_items.sku":      "SKU-123" } },
            { "range": { "line_items.quantity":  { "gte": 2 } } }
          ]
        }
      },
      "inner_hits": {         ← return matching nested elements
        "size": 5,
        "_source": ["line_items.sku", "line_items.quantity"]
      }
    }
  }
}
```

---

## Has Child / Has Parent (Join Field)

For parent-child relationships within the same index using a `join` field:

```bash
# Mapping with join field
PUT /store
{
  "mappings": {
    "properties": {
      "join_field": {
        "type":      "join",
        "relations": { "product": "review" }
      }
    }
  }
}

# Index parent
PUT /store/_doc/1
{ "name": "iPhone 15", "join_field": "product" }

# Index child (must go to same shard as parent using routing)
PUT /store/_doc/100?routing=1
{
  "rating": 5,
  "comment": "Great phone!",
  "join_field": { "name": "review", "parent": "1" }
}

# has_child query: find products that have a 5-star review
GET /store/_search
{
  "query": {
    "has_child": {
      "type":  "review",
      "query": { "term": { "rating": 5 } },
      "score_mode": "max",
      "inner_hits": {}
    }
  }
}

# has_parent query: find reviews whose product is an iPhone
GET /store/_search
{
  "query": {
    "has_parent": {
      "parent_type": "product",
      "query": { "match": { "name": "iphone" } }
    }
  }
}
```

---

## Percolate Query

**Percolation** is the reverse of a normal query — you store queries and test documents against them. Great for notification systems:

```bash
# Store a "query" document
PUT /alerts/_doc/1
{
  "query": {
    "bool": {
      "must": [
        { "term":  { "event_type": "login_failure" } },
        { "range": { "count": { "gte": 5 } } }
      ]
    }
  }
}

# Test if a new event matches any stored queries
GET /alerts/_search
{
  "query": {
    "percolate": {
      "field": "query",
      "document": {
        "event_type": "login_failure",
        "count":       7,
        "user":        "john@example.com"
      }
    }
  }
}
```

---

## More Like This (MLT)

Find documents similar to one or more given documents:

```bash
GET /articles/_search
{
  "query": {
    "more_like_this": {
      "fields":        ["title", "content"],
      "like": [
        { "_index": "articles", "_id": "42" }
      ],
      "min_term_freq":    1,
      "max_query_terms":  12,
      "min_doc_freq":     2
    }
  }
}

# Or based on a text snippet
{
  "query": {
    "more_like_this": {
      "fields": ["title", "body"],
      "like":   "Elasticsearch distributed search engine tutorial",
      "min_term_freq": 1,
      "min_doc_freq":  2
    }
  }
}
```

---

## Span Queries (Positional)

Control exact positions of terms:

```bash
# Match "quick brown" within 3 positions of each other
{
  "query": {
    "span_near": {
      "clauses": [
        { "span_term": { "description": "quick" } },
        { "span_term": { "description": "brown" } }
      ],
      "slop":     3,
      "in_order": true
    }
  }
}
```

---

## Function Score Query

Customize relevance scoring with functions:

```bash
GET /products/_search
{
  "query": {
    "function_score": {
      "query": { "match": { "name": "headphones" } },
      "functions": [
        {
          "filter": { "term": { "featured": true } },
          "weight": 3
        },
        {
          "field_value_factor": {
            "field":    "popularity_score",
            "factor":   1.5,
            "modifier": "sqrt",
            "missing":  1
          }
        },
        {
          "gauss": {
            "price": {
              "origin": "100",
              "scale":  "50",
              "offset": "10",
              "decay":  0.5
            }
          }
        }
      ],
      "score_mode":  "multiply",
      "boost_mode":  "multiply"
    }
  }
}
```

---

## Completion Suggester (Autocomplete)

Fast prefix autocomplete powered by an in-memory FST structure:

```bash
PUT /products
{
  "mappings": {
    "properties": {
      "suggest": {
        "type":     "completion",
        "analyzer": "simple"
      }
    }
  }
}

# Index with completion data
PUT /products/_doc/1
{
  "name":    "iPhone 15 Pro",
  "suggest": {
    "input":  ["iphone", "iphone 15", "iphone 15 pro", "apple phone"],
    "weight": 10
  }
}

# Query
GET /products/_search
{
  "suggest": {
    "product_suggestion": {
      "prefix": "ipho",
      "completion": {
        "field":        "suggest",
        "size":         5,
        "fuzzy":        { "fuzziness": 1 },
        "skip_duplicates": true
      }
    }
  }
}
```

---

## kNN Vector Search (Semantic Search)

Store and search dense vector embeddings for semantic similarity (ES 8.x):

```bash
PUT /docs
{
  "mappings": {
    "properties": {
      "title":     { "type": "text" },
      "embedding": {
        "type":       "dense_vector",
        "dims":       768,
        "index":      true,
        "similarity": "cosine"
      }
    }
  }
}

# Search by vector
GET /docs/_search
{
  "knn": {
    "field":          "embedding",
    "query_vector":   [0.12, -0.45, 0.78, ...],  ← 768 dims
    "k":              10,
    "num_candidates": 100
  },
  "query": {
    "match": { "title": "elasticsearch" }
  }
}
```

Use this to build semantic search, RAG systems, and image similarity engines.


---

# Chapter 19: Clients & Integration

## Official Elasticsearch Clients

| Language | Package |
|---|---|
| Python | `elasticsearch` (v8.x) |
| Node.js | `@elastic/elasticsearch` |
| Java | `elasticsearch-java` |
| Go | `go-elasticsearch` |
| Ruby | `elasticsearch-ruby` |
| PHP | `elasticsearch-php` |
| .NET (C#) | `Elastic.Clients.Elasticsearch` |
| Rust | `elasticsearch` |

All official clients:
- Auto-detect cluster topology (sniffing)
- Handle retries and failover
- Serialize/deserialize JSON
- Support async operations
- Follow the Elasticsearch REST API 1:1

---

## Python Client

### Installation

```bash
pip install elasticsearch==8.x.x
```

### Basic Connection

```python
from elasticsearch import Elasticsearch

# Without security (dev)
es = Elasticsearch("http://localhost:9200")

# With API key
es = Elasticsearch(
    "https://my-cluster.elastic.co",
    api_key=("key_id", "key_secret"),
)

# With username/password + TLS
es = Elasticsearch(
    "https://localhost:9200",
    basic_auth=("elastic", "password"),
    ca_certs="/path/to/http_ca.crt",
)

# Check connection
info = es.info()
print(info["version"]["number"])
```

### CRUD Operations

```python
# Index a document
response = es.index(
    index="products",
    id="1",
    document={
        "name":     "MacBook Pro",
        "price":    2499.00,
        "in_stock": True,
    }
)
print(response["result"])  # "created" or "updated"

# Get a document
doc = es.get(index="products", id="1")
print(doc["_source"]["name"])

# Check existence
exists = es.exists(index="products", id="1")

# Partial update
es.update(
    index="products",
    id="1",
    doc={"price": 2299.00}
)

# Delete
es.delete(index="products", id="1")
```

### Search

```python
# Simple match search
results = es.search(
    index="products",
    query={"match": {"name": "laptop"}},
    size=10,
    from_=0,
)

for hit in results["hits"]["hits"]:
    print(hit["_id"], hit["_source"]["name"], hit["_score"])

# Complex bool query
results = es.search(
    index="products",
    body={
        "query": {
            "bool": {
                "must": [{"match": {"name": "laptop"}}],
                "filter": [
                    {"term": {"in_stock": True}},
                    {"range": {"price": {"lte": 1500}}},
                ],
            }
        },
        "sort": [{"price": {"order": "asc"}}],
        "aggs": {
            "avg_price": {"avg": {"field": "price"}}
        },
    },
)

avg = results["aggregations"]["avg_price"]["value"]
total = results["hits"]["total"]["value"]
```

### Bulk Indexing

```python
from elasticsearch.helpers import bulk, streaming_bulk

# helpers.bulk — list of actions
actions = [
    {
        "_index":  "products",
        "_id":     str(product["id"]),
        "_source": product,
    }
    for product in products_list
]

success, failed = bulk(es, actions, chunk_size=1000)
print(f"Indexed: {success}, Failed: {len(failed)}")

# streaming_bulk — for memory efficiency with large datasets
for ok, info in streaming_bulk(es, actions, chunk_size=500):
    if not ok:
        print(f"Failed: {info}")
```

### Using Pandas with ES

```python
import pandas as pd
from elasticsearch.helpers import bulk

df = pd.read_csv("products.csv")

def df_to_actions(df, index):
    for _, row in df.iterrows():
        yield {
            "_index": index,
            "_source": row.to_dict(),
        }

bulk(es, df_to_actions(df, "products"))
```

---

## Node.js Client

### Installation

```bash
npm install @elastic/elasticsearch
```

### Basic Usage

```javascript
const { Client } = require('@elastic/elasticsearch');

const client = new Client({
  node: 'http://localhost:9200',
  auth: { apiKey: 'base64encodedkey' },
});

// Index a document
await client.index({
  index: 'products',
  id:    '1',
  document: {
    name:  'Wireless Keyboard',
    price: 79.99,
  },
});

// Search
const result = await client.search({
  index: 'products',
  query: {
    match: { name: 'keyboard' },
  },
});

result.hits.hits.forEach(hit => {
  console.log(hit._id, hit._source.name);
});

// Bulk indexing
const operations = products.flatMap(doc => [
  { index: { _index: 'products', _id: doc.id } },
  doc,
]);

const bulkResponse = await client.bulk({ operations });
if (bulkResponse.errors) {
  console.error('Bulk errors:', bulkResponse.items.filter(i => i.index?.error));
}
```

---

## REST API Patterns

### Pagination with Search After

```python
def paginate_all(es, index, query, batch_size=1000):
    """Retrieve all matching documents using search_after."""
    results = []
    sort = [{"_id": "asc"}]
    search_after = None

    while True:
        body = {
            "query":  query,
            "sort":   sort,
            "size":   batch_size,
        }
        if search_after:
            body["search_after"] = search_after

        response = es.search(index=index, body=body)
        hits = response["hits"]["hits"]
        if not hits:
            break

        results.extend(hits)
        search_after = hits[-1]["sort"]

    return results
```

### Point-in-Time Pagination

```python
# Open a PIT for consistent pagination
pit = es.open_point_in_time(index="products", keep_alive="1m")
pit_id = pit["id"]

try:
    search_after = None
    while True:
        body = {
            "pit":  {"id": pit_id, "keep_alive": "1m"},
            "query": {"match_all": {}},
            "sort": [{"_id": "asc"}],
            "size": 1000,
        }
        if search_after:
            body["search_after"] = search_after

        resp = es.search(body=body)
        hits = resp["hits"]["hits"]
        if not hits:
            break

        process(hits)
        search_after = hits[-1]["sort"]
        pit_id = resp["pit_id"]  # PIT ID can change
finally:
    es.close_point_in_time(id=pit_id)
```

---

## Using the REST API Directly (curl)

```bash
# Index with basic auth
curl -X PUT "http://localhost:9200/products/_doc/1" \
     -H "Content-Type: application/json" \
     -u elastic:password \
     -d '{"name": "MacBook", "price": 1999}'

# Search
curl -X GET "http://localhost:9200/products/_search" \
     -H "Content-Type: application/json" \
     -d '{"query": {"match": {"name": "macbook"}}}'

# With API key
curl -X GET "https://cluster.elastic.co:9243/products/_search" \
     -H "Authorization: ApiKey base64key==" \
     -H "Content-Type: application/json" \
     -d '{"query": {"match_all": {}}}'
```

---

## Common Integration Patterns

### Sync from a Database (CDC Pattern)

```python
import psycopg2
from elasticsearch import Elasticsearch
from elasticsearch.helpers import bulk

def sync_postgres_to_es():
    pg = psycopg2.connect("postgresql://...")
    es = Elasticsearch("http://localhost:9200")

    cursor = pg.cursor()
    cursor.execute("SELECT id, name, price, updated_at FROM products")

    def actions():
        for row in cursor:
            yield {
                "_index":   "products",
                "_id":       str(row[0]),
                "_source": {
                    "name":       row[1],
                    "price":      row[2],
                    "updated_at": row[3].isoformat(),
                },
            }

    success, failed = bulk(es, actions(), chunk_size=500)
    print(f"Synced {success} documents")
```

### Error Handling and Retries

```python
from elasticsearch import (
    Elasticsearch,
    ConnectionError,
    TransportError,
    NotFoundError,
    ConflictError,
)
import time

def safe_index(es, index, doc_id, document, retries=3):
    for attempt in range(retries):
        try:
            return es.index(index=index, id=doc_id, document=document)
        except ConflictError:
            raise  # Don't retry version conflicts
        except ConnectionError:
            if attempt < retries - 1:
                time.sleep(2 ** attempt)  # exponential backoff
            else:
                raise
        except TransportError as e:
            if e.status_code in (429, 503):  # throttled or unavailable
                time.sleep(2 ** attempt)
            else:
                raise
```

---

## Elasticsearch SQL

ES supports SQL queries via the SQL API (ES 6.3+):

```bash
POST /_sql
{
  "query": "SELECT name, price, category FROM products WHERE price < 500 ORDER BY price DESC LIMIT 10"
}

# Translate SQL to Query DSL
POST /_sql/translate
{
  "query": "SELECT * FROM products WHERE category = 'electronics'"
}
```

Python:

```python
result = es.sql.query(body={
    "query": "SELECT name, AVG(price) as avg_price FROM products GROUP BY category"
})
for row in result["rows"]:
    print(row)
```


---

# Chapter 20: Learning Roadmap

## 5-Week Hands-On Plan

This roadmap takes you from zero to production-ready Elasticsearch skills with a hands-on project each week.

---

## Week 1: Foundation

**Goal:** Run Elasticsearch, index data, run basic searches.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Install Docker, spin up ES + Kibana with docker-compose |
| 2 | Use Dev Tools: create an index, index 10 documents manually |
| 3 | Run all 5 Query DSL types: match, term, range, bool, exists |
| 4 | Use the Bulk API to load 1,000 sample records |
| 5 | Explore data in Kibana Discover |

### Practice Project: Product Catalog

Create a `products` index with 1,000 products (use a faker library or download a sample dataset).

```python
# Setup script
from faker import Faker
from elasticsearch.helpers import bulk
from elasticsearch import Elasticsearch

fake = Faker()
es   = Elasticsearch("http://localhost:9200")

def make_product(i):
    return {
        "name":        fake.catch_phrase(),
        "brand":       fake.company(),
        "price":       round(fake.random.uniform(5, 5000), 2),
        "category":    fake.random_element(["electronics", "clothing", "books", "home", "sports"]),
        "in_stock":    fake.boolean(chance_of_getting_true=75),
        "rating":      round(fake.random.uniform(1, 5), 1),
        "created_at":  fake.date_time_this_year().isoformat(),
        "description": fake.paragraph(nb_sentences=3),
    }

actions = [
    {"_index": "products", "_id": str(i), "_source": make_product(i)}
    for i in range(1, 1001)
]
bulk(es, actions)
```

### Week 1 Checklist
- [ ] Docker Compose with ES + Kibana running
- [ ] `GET _cluster/health` returns green
- [ ] Created index with explicit mapping
- [ ] Indexed 1,000 documents via Bulk API
- [ ] Run 5 different query types
- [ ] Filter results (bool query with filter context)
- [ ] Sort by price, paginate with from/size

---

## Week 2: Search Mastery

**Goal:** Master Query DSL, text analysis, and aggregations.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Implement multi_match, match_phrase, fuzzy queries |
| 2 | Build a custom analyzer with synonyms and stemming |
| 3 | Implement autocomplete using edge_ngram |
| 4 | Build aggregations: terms, date_histogram, avg, nested |
| 5 | Create a Kibana dashboard from your aggregations |

### Practice: Search Feature

Build a Python function that:
1. Full-text searches `name` and `description`
2. Filters by `category`, price range, and `in_stock`
3. Returns results sorted by relevance, with pagination
4. Shows facet counts (category breakdown)

```python
def search_products(query_text, category=None, min_price=None,
                    max_price=None, in_stock_only=False, page=1, size=10):
    must   = []
    filter = []

    if query_text:
        must.append({
            "multi_match": {
                "query":  query_text,
                "fields": ["name^3", "description", "brand^2"],
            }
        })
    else:
        must.append({"match_all": {}})

    if category:
        filter.append({"term": {"category": category}})
    if in_stock_only:
        filter.append({"term": {"in_stock": True}})
    if min_price or max_price:
        price_range = {}
        if min_price: price_range["gte"] = min_price
        if max_price: price_range["lte"] = max_price
        filter.append({"range": {"price": price_range}})

    return es.search(
        index="products",
        query={"bool": {"must": must, "filter": filter}},
        aggs={
            "categories":  {"terms": {"field": "category"}},
            "price_ranges": {
                "range": {
                    "field": "price",
                    "ranges": [
                        {"to": 50}, {"from": 50, "to": 200},
                        {"from": 200, "to": 1000}, {"from": 1000}
                    ]
                }
            },
            "avg_price": {"avg": {"field": "price"}},
        },
        from_=(page - 1) * size,
        size=size,
    )
```

### Week 2 Checklist
- [ ] Custom analyzer with edge_ngram for autocomplete
- [ ] Synonym filter configured
- [ ] multi_match with field boosting
- [ ] Fuzzy query with AUTO fuzziness
- [ ] terms, date_histogram, and nested aggregations
- [ ] Kibana dashboard with 4+ panels

---

## Week 3: Cluster Operations

**Goal:** Understand cluster architecture, mappings, and index management.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Create a 3-node cluster with docker-compose |
| 2 | Design and apply explicit mappings for a new use case |
| 3 | Set up ILM policy for log rotation |
| 4 | Practice reindexing with mapping changes |
| 5 | Configure index templates and aliases |

### Practice: Log Pipeline

Create a logs pipeline:
1. Template: `logs-*` with ILM policy attached
2. ILM: hot (7d rollover), warm (30d), delete (90d)
3. Load sample logs via Filebeat or Logstash
4. Verify rollover happens

### Week 3 Checklist
- [ ] 3-node cluster running (docker-compose)
- [ ] Explicit mapping with nested, keyword, text, date types
- [ ] ILM policy with hot/warm/delete phases
- [ ] Index template attached to ILM
- [ ] Reindex with field transformation script
- [ ] Zero-downtime alias swap after reindex

---

## Week 4: Production Readiness

**Goal:** Security, performance, monitoring, snapshots.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Enable security, create users and roles |
| 2 | Set up TLS for HTTP layer |
| 3 | Implement performance optimizations (bulk, refresh, codec) |
| 4 | Configure slow logs, monitor with `_cat` APIs |
| 5 | Set up S3/GCS snapshot repository and SLM policy |

### Practice: Secure Multi-Tenant Setup

Create two roles:
- `product_reader` — read products index
- `orders_admin` — full access to orders index

Create two users and test access restrictions.

### Week 4 Checklist
- [ ] Security enabled with custom users and roles
- [ ] API key created and used for authentication
- [ ] Bulk indexing benchmarked with refresh disabled
- [ ] Slow logs capturing queries > 1 second
- [ ] Snapshot repository registered and tested
- [ ] SLM daily snapshot policy running

---

## Week 5: Advanced Topics + Real Project

**Goal:** Apply everything to a real-world project.

### Capstone Project Options

#### Option A: E-Commerce Search Engine
- Product catalog with 100K+ products
- Autocomplete suggestions
- Faceted search (category, price, brand filters)
- Fuzzy search with typo tolerance
- Personalized ranking (boost purchased brands)
- Real-time stock updates via `_update_by_query`

#### Option B: Log Analytics Platform
- ELK Stack with Filebeat + Logstash
- Parse Apache/Nginx access logs
- Dashboard: RPS, error rate, top URLs, geo map
- Alert: error rate > 5% for 5 minutes
- ILM to auto-delete logs older than 90 days

#### Option C: Semantic Document Search
- Index company documents/PDFs
- Generate embeddings via Hugging Face / OpenAI API
- Store as `dense_vector` fields
- Build hybrid search: BM25 + kNN vector search
- Compare semantic results vs. keyword results

### Week 5 Checklist
- [ ] Capstone project deployed (Docker or Elastic Cloud)
- [ ] Explicit mapping with all relevant types
- [ ] Security configured
- [ ] ILM or manual index management
- [ ] Monitoring dashboard in Kibana
- [ ] Snapshot policy active
- [ ] Python or Node.js client integration

---

## Beyond the Roadmap

| Topic | Resources |
|---|---|
| **Elastic Certified Engineer** | Official certification exam |
| **Machine Learning** | Anomaly detection, NLP with ELSER |
| **Cross-Cluster Search** | Multi-datacenter search |
| **Fleet & Agent** | Unified data collection |
| **Elastic Observability** | APM, synthetics, profiling |
| **SIEM** | Elastic Security, detection rules |

---

## Key Commands Reference Card

```bash
# Health & Info
GET /                         ← cluster info
GET _cluster/health           ← cluster health
GET _cat/nodes?v              ← node list
GET _cat/indices?v            ← index list
GET _cat/shards?v             ← shard list

# Index Operations
PUT /<index>                  ← create index
DELETE /<index>               ← delete index
GET /<index>/_mapping         ← view mapping
GET /<index>/_settings        ← view settings
POST /<index>/_refresh        ← force refresh

# Documents
PUT /<index>/_doc/<id>        ← create/replace
POST /<index>/_doc            ← create (auto-id)
GET /<index>/_doc/<id>        ← get by id
POST /<index>/_update/<id>    ← partial update
DELETE /<index>/_doc/<id>     ← delete by id
POST /_bulk                   ← bulk operations

# Search
GET /<index>/_search          ← search
GET /<index>/_count           ← count matching docs
POST /_msearch                ← multi-search

# Analysis
POST /_analyze                ← test analyzer
POST /<index>/_analyze        ← test against index

# Snapshots
PUT /_snapshot/<repo>/<snap>  ← create snapshot
POST /_snapshot/<repo>/<snap>/_restore ← restore
GET /_snapshot/<repo>/*       ← list snapshots
```
