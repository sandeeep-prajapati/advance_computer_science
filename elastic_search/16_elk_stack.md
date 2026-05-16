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
