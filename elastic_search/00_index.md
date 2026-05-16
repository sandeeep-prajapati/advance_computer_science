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
