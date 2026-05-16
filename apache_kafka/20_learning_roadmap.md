# Chapter 20: Learning Roadmap

## 5-Week Hands-On Plan

This roadmap takes you from zero to production-ready Apache Kafka skills with a real project each week.

---

## Week 1: Foundation

**Goal:** Run Kafka locally, send and receive messages, understand the core model.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Start Kafka with Docker Compose (KRaft mode), verify with CLI tools |
| 2 | Create topics, produce messages from CLI, consume from beginning |
| 3 | Write Python producer (confluent-kafka) and consumer |
| 4 | Explore consumer groups: run 3 consumers on a 6-partition topic |
| 5 | Trigger rebalances by killing/starting consumers, observe in Kafka UI |

### Practice Project: User Activity Tracker

```python
# Producer: simulate user activity events
from confluent_kafka import Producer
import json, time, random

producer = Producer({"bootstrap.servers": "localhost:9092"})

EVENTS = ["page_view", "click", "add_to_cart", "purchase", "logout"]
PAGES  = ["/home", "/product/123", "/cart", "/checkout", "/profile"]

while True:
    event = {
        "user_id":    f"user_{random.randint(1, 1000):04d}",
        "event_type": random.choice(EVENTS),
        "page":       random.choice(PAGES),
        "ts":         int(time.time() * 1000),
        "session_id": f"sess_{random.randint(1, 100):03d}",
    }
    producer.produce(
        "user-events",
        key   = event["user_id"].encode(),
        value = json.dumps(event).encode(),
    )
    producer.poll(0)
    time.sleep(0.05)
```

```python
# Consumer: count events per type
from confluent_kafka import Consumer
import json
from collections import defaultdict

consumer = Consumer({
    "bootstrap.servers": "localhost:9092",
    "group.id":          "activity-counter",
    "auto.offset.reset": "earliest",
})
consumer.subscribe(["user-events"])

counts = defaultdict(int)
while True:
    msg = consumer.poll(1.0)
    if msg is None or msg.error(): continue
    event = json.loads(msg.value())
    counts[event["event_type"]] += 1
    if sum(counts.values()) % 100 == 0:
        print(dict(counts))
    consumer.commit(asynchronous=True)
```

### Week 1 Checklist
- [ ] Kafka running in Docker
- [ ] Topic created with 6 partitions, RF=1
- [ ] CLI producer + consumer verified
- [ ] Python producer sending events
- [ ] Python consumer with manual commit
- [ ] Consumer group with 2+ instances observed in Kafka UI

---

## Week 2: Deep Dive — Reliability & Schema

**Goal:** Understand delivery guarantees, implement idempotency, add Schema Registry.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Configure `acks=all`, `enable.idempotence=true`, test with broker restart |
| 2 | Add Schema Registry, define Avro schema, produce/consume with serialization |
| 3 | Evolve schema (add nullable field), verify BACKWARD compatibility |
| 4 | Implement at-least-once consumer with idempotent processing |
| 5 | Implement dead letter queue for poison pill messages |

### Practice: Order Processing with Avro

```python
# order_schema.avsc
order_schema_str = """
{
  "type": "record",
  "name": "Order",
  "namespace": "com.example",
  "fields": [
    {"name": "order_id",   "type": "string"},
    {"name": "user_id",    "type": "string"},
    {"name": "amount",     "type": "double"},
    {"name": "status",     "type": {"type": "enum", "name": "Status",
                                   "symbols": ["PENDING", "CONFIRMED", "SHIPPED"]}},
    {"name": "created_at", "type": "long"}
  ]
}
"""
```

### Week 2 Checklist
- [ ] Producer with `acks=all` + idempotence enabled
- [ ] Avro schema registered in Schema Registry
- [ ] Producer serializing with Avro
- [ ] Consumer deserializing with Avro
- [ ] Schema evolution tested (add field with default)
- [ ] Dead letter queue implemented

---

## Week 3: Streaming + Connect

**Goal:** Build a real-time pipeline with Kafka Connect and stream processing.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Set up Kafka Connect (distributed mode) with Docker |
| 2 | Configure Debezium source to capture PostgreSQL CDC |
| 3 | Configure Elasticsearch sink to index changed rows |
| 4 | Write Faust (or raw consumer) stream processor with aggregation |
| 5 | Add windowed session tracking (count events per user per 5 min) |

### Practice: CDC Pipeline

```yaml
# Docker Compose addition
postgres:
  image: postgres:15
  environment:
    POSTGRES_USER: user
    POSTGRES_PASSWORD: secret
    POSTGRES_DB: shop
  command: ["postgres", "-c", "wal_level=logical"]
```

```sql
-- In PostgreSQL
CREATE TABLE orders (
    order_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    VARCHAR(100),
    amount     NUMERIC(10, 2),
    status     VARCHAR(20) DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Enable replication
ALTER TABLE orders REPLICA IDENTITY FULL;
```

### Week 3 Checklist
- [ ] Kafka Connect running in distributed mode
- [ ] Debezium capturing inserts/updates/deletes from PostgreSQL
- [ ] Changes flowing to Kafka topic
- [ ] Elasticsearch sink indexing events
- [ ] Python stream processor aggregating metrics
- [ ] Window-based aggregation working

---

## Week 4: Performance & Operations

**Goal:** Tune Kafka for throughput, monitor lag, handle failures.

### Daily Tasks

| Day | Task |
|---|---|
| 1 | Run producer perf test, tune batch.size + linger.ms + compression |
| 2 | Simulate broker failure, observe leader election, verify no data loss |
| 3 | Set up Prometheus + Grafana with Kafka JMX metrics |
| 4 | Introduce artificial consumer lag, alert when lag > threshold |
| 5 | Configure SASL/PLAIN auth + ACLs for least-privilege access |

### Week 4 Checklist
- [ ] Perf test baseline measured (messages/sec, MB/s)
- [ ] Compression (lz4) enabled and measured improvement
- [ ] Broker failure simulated, cluster recovered automatically
- [ ] Grafana dashboard with: broker health, consumer lag, throughput
- [ ] Alert rule for consumer lag > 10,000
- [ ] SASL auth configured, consumer with wrong password rejected

---

## Week 5: Capstone Project

**Goal:** Build an end-to-end production-grade streaming system.

### Capstone Option A: Real-Time Analytics Platform

```
User App → Kafka (events) → Spark Streaming → Delta Lake → BI Dashboard
                         → Faust aggregator → Kafka (alerts) → Slack/PagerDuty
```

Features:
- High-throughput event ingestion (simulate 10,000 events/sec)
- 1-minute, 5-minute, 1-hour windowed KPIs
- Delta Lake as analytical store (time travel queries)
- Real-time alerting on anomalies (revenue drops, error spikes)
- Consumer lag dashboard in Grafana

### Capstone Option B: CDC-Driven Data Sync

```
PostgreSQL → Debezium → Kafka → (Elasticsearch + ClickHouse + Redis)
                              → Schema Registry (Avro)
                              → MirrorMaker 2 → DR Cluster
```

Features:
- Full-table initial snapshot + ongoing CDC
- Avro schema with Schema Registry
- Elasticsearch for search, ClickHouse for analytics, Redis for cache
- Schema evolution: add new columns without downtime
- Failover test to DR cluster

### Capstone Option C: Exactly-Once Payment Processor

```
Payment API → Kafka (raw-payments) → EOS Processor → Kafka (settled) → Ledger DB
                                   → Idempotent consumer → PostgreSQL
```

Features:
- Exactly-once semantics with transactions
- Idempotent consumer with database upsert
- Dead letter queue for invalid payments
- Audit log (every payment event stored forever with compaction)
- Load test at 1,000 TPS

---

## Key Commands Reference Card

```bash
# Topics
kafka-topics --bootstrap-server localhost:9092 --list
kafka-topics --bootstrap-server localhost:9092 --describe --topic orders
kafka-topics --bootstrap-server localhost:9092 --create --topic test \
  --partitions 6 --replication-factor 1

# Consumer groups
kafka-consumer-groups --bootstrap-server localhost:9092 --list
kafka-consumer-groups --bootstrap-server localhost:9092 --describe --group my-app
kafka-consumer-groups --bootstrap-server localhost:9092 \
  --reset-offsets --to-earliest --group my-app --topic orders --execute

# Config
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics --entity-name orders --describe
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics --entity-name orders \
  --alter --add-config retention.ms=86400000

# Performance
kafka-producer-perf-test --topic perf-test --num-records 1000000 \
  --record-size 1000 --throughput -1 \
  --producer-props bootstrap.servers=localhost:9092

kafka-consumer-perf-test --bootstrap-server localhost:9092 \
  --topic perf-test --messages 1000000

# Cluster health
kafka-topics --bootstrap-server localhost:9092 --describe --under-replicated-partitions
kafka-topics --bootstrap-server localhost:9092 --describe --unavailable-partitions
kafka-leader-election --bootstrap-server localhost:9092 \
  --election-type preferred --all-topic-partitions
```

---

## Beyond the Roadmap

| Topic | Notes |
|---|---|
| **ksqlDB** | SQL-based stream processing on Kafka |
| **Kafka Tiered Storage** | Offload old segments to S3 (KIP-405) |
| **Kafka Raft (KRaft)** | ZooKeeper-free Kafka — future default |
| **Confluent Platform** | Enterprise Kafka with Schema Registry, KSQL, Control Center |
| **Apache Pulsar** | Alternative to Kafka — multi-tenancy, geo-replication |
| **Redpanda** | Drop-in Kafka replacement in C++ (no JVM) |
| **Apache Iceberg + Kafka** | Stream directly into Iceberg tables |
| **Kafka MCP** | OpenTelemetry integration for distributed tracing |
| **Debezium** | Advanced CDC patterns (filtering, routing, transforms) |
| **Apache NiFi** | Visual data flow with Kafka integration |
