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
