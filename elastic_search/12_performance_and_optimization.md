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
