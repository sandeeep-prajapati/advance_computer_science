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
