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
