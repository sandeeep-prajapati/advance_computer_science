# Chapter 14: Monitoring & Operations

## Key Metrics to Watch

### Broker Health

| Metric | JMX Path | Normal | Alert |
|---|---|---|---|
| Under-replicated partitions | `kafka.server:ReplicaManager,UnderReplicatedPartitions` | 0 | > 0 |
| Active controller count | `kafka.controller:KafkaController,ActiveControllerCount` | 1 (per cluster) | ≠ 1 |
| Offline partitions | `kafka.controller:KafkaController,OfflinePartitionsCount` | 0 | > 0 |
| ISR shrinks/sec | `kafka.server:ReplicaManager,IsrShrinksPerSec` | ~0 | Spike |
| Request handler idle % | `kafka.server:KafkaRequestHandlerPool,RequestHandlerAvgIdlePercent` | > 0.3 | < 0.2 |
| Network processor idle % | `kafka.network:SocketServer,NetworkProcessorAvgIdlePercent` | > 0.3 | < 0.2 |

### Producer Metrics

| Metric | Normal | Alert |
|---|---|---|
| `record-error-rate` | 0 | > 0 |
| `record-send-rate` | Expected | Drop |
| `request-latency-avg` | < 100ms | > 500ms |
| `buffer-available-bytes` | > 10% total | < 10% |

### Consumer Metrics

| Metric | JMX Path | Alert |
|---|---|---|
| Consumer lag | `kafka.consumer:ConsumerFetchManager-Metrics,records-lag-max` | > threshold |
| Fetch rate | `kafka.consumer:ConsumerFetchManager-Metrics,fetch-rate` | Drop |
| Commit rate | `kafka.consumer:ConsumerCoordinator-Metrics,commit-rate` | Drop |

---

## Prometheus + Grafana Setup

### JMX Exporter

```yaml
# kafka-jmx-exporter.yml
startDelaySeconds: 0
ssl: false
lowercaseOutputName: false
lowercaseOutputLabelNames: false
whitelistObjectNames:
  - "kafka.server:type=BrokerTopicMetrics,name=*"
  - "kafka.server:type=ReplicaManager,name=*"
  - "kafka.controller:type=KafkaController,name=*"
  - "kafka.network:type=RequestMetrics,name=*,request=*"
  - "kafka.log:type=LogFlushStats,name=*"
```

```bash
# In KAFKA_OPTS
-javaagent:/opt/jmx_exporter/jmx_prometheus_javaagent.jar=7071:/opt/jmx_exporter/kafka.yml
```

### docker-compose with monitoring

```yaml
prometheus:
  image: prom/prometheus
  volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml
  ports: ["9090:9090"]

grafana:
  image: grafana/grafana
  ports: ["3000:3000"]
  environment:
    GF_SECURITY_ADMIN_PASSWORD: admin
```

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'kafka-broker-1'
    static_configs:
      - targets: ['broker1:7071']
  - job_name: 'kafka-broker-2'
    static_configs:
      - targets: ['broker2:7071']
```

---

## Consumer Lag Monitoring

Consumer lag is the most important consumer-side metric:

```bash
# CLI lag check
kafka-consumer-groups --bootstrap-server localhost:9092 \
  --describe --group my-app

# Output:
# TOPIC    PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG  CONSUMER-ID
# orders   0          1234            1260            26   consumer-1
# orders   1          890             905             15   consumer-2
# orders   2          2100            2100            0    consumer-3
```

### Programmatic Lag Check (Python)

```python
from confluent_kafka.admin import AdminClient
from confluent_kafka import Consumer, TopicPartition

admin = AdminClient({"bootstrap.servers": "localhost:9092"})
consumer = Consumer({"bootstrap.servers": "localhost:9092", "group.id": "lag-monitor"})

def get_consumer_lag(group_id: str, topic: str, num_partitions: int) -> dict:
    committed = consumer.committed([
        TopicPartition(topic, p) for p in range(num_partitions)
    ])

    end_offsets = consumer.get_watermark_offsets(TopicPartition(topic, 0))

    # For all partitions:
    lag = {}
    for tp in committed:
        _, high = consumer.get_watermark_offsets(tp, timeout=5)
        current = tp.offset if tp.offset >= 0 else 0
        lag[tp.partition] = high - current

    return lag

print(get_consumer_lag("order-processor", "orders", 6))
```

---

## Operations Runbook

### Broker Is Down

```bash
# 1. Check which partitions are affected
kafka-topics --bootstrap-server localhost:9092 \
  --describe --under-replicated-partitions

# 2. Check if leader election happened
kafka-topics --bootstrap-server localhost:9092 \
  --describe --topic orders | grep "Isr:"

# 3. Restart broker
systemctl restart kafka

# 4. Monitor ISR expansion
watch -n 5 "kafka-topics --bootstrap-server localhost:9092 \
  --describe --under-replicated-partitions"

# 5. After restart, trigger preferred leader election
kafka-leader-election --bootstrap-server localhost:9092 \
  --election-type preferred --all-topic-partitions
```

### High Consumer Lag

```bash
# 1. Check lag per partition
kafka-consumer-groups --bootstrap-server localhost:9092 \
  --describe --group my-app

# 2. Scale up consumers (if partitions available)
# Deploy more consumer instances

# 3. Check if consumer is crashing
kafka-consumer-groups --bootstrap-server localhost:9092 \
  --describe --group my-app | grep -v "STABLE"

# 4. Temporarily increase consumers beyond partition count (won't help)
# → Add partitions instead
kafka-topics --bootstrap-server localhost:9092 \
  --alter --topic orders --partitions 12
```

### Topic Deletion

```bash
# Delete topic
kafka-topics --bootstrap-server localhost:9092 \
  --delete --topic old-topic

# If deletion is not completing, check:
# delete.topic.enable=true must be set in server.properties
```

### Purge a Topic (Reset Retention)

```bash
# Set retention to 1 second
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics --entity-name orders \
  --alter --add-config "retention.ms=1000"

# Wait ~30 seconds for deletion
sleep 30

# Restore original retention
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics --entity-name orders \
  --alter --delete-config "retention.ms"
```

---

## Disk Usage

```bash
# Log sizes per topic-partition
du -sh /var/kafka/data/*/

# Topic-level disk usage
kafka-log-dirs --bootstrap-server localhost:9092 \
  --topic-list orders \
  --describe | jq .

# Find large segments
find /var/kafka/data -name "*.log" -size +1G | sort -k5 -n
```

---

## Kafka Admin API (Python)

```python
from confluent_kafka.admin import AdminClient, NewTopic, ConfigResource

admin = AdminClient({"bootstrap.servers": "localhost:9092"})

# Create topics
admin.create_topics([
    NewTopic("events", num_partitions=12, replication_factor=3,
             config={"retention.ms": "604800000", "compression.type": "lz4"}),
])

# List topics
metadata = admin.list_topics(timeout=10)
for topic in metadata.topics.values():
    print(f"{topic.topic}: {len(topic.partitions)} partitions")

# Describe config
resource = ConfigResource(ConfigResource.Type.TOPIC, "orders")
configs = admin.describe_configs([resource])
for key, val in configs[resource].result().items():
    print(f"{key} = {val.value}")

# Alter config
from confluent_kafka.admin import ConfigEntry
admin.alter_configs([
    (resource, {"retention.ms": ConfigEntry("retention.ms", "86400000")})
])
```
