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
