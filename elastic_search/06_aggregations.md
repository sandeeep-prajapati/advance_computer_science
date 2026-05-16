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
