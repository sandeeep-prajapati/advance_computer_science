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
