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
