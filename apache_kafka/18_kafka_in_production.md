# Chapter 18: Kafka in Production

## Production Cluster Sizing

### Broker Hardware

```
Throughput target: 100 MB/s write, 7-day retention, RF=3

Storage per broker:
  = (100 MB/s × 86400 sec/day × 7 days × 3 replicas) / num_brokers
  For 6 brokers: = 100 × 86400 × 7 × 3 / 6 = 3.024 TB per broker

Recommended hardware (per broker):
  CPU:    8–16 cores
  RAM:    32–64 GB (leave most for OS page cache)
  Disk:   4–8 TB NVMe SSDs (or JBOD HDDs for archival)
  Network: 10 Gbps
  OS:     Linux (CentOS, Ubuntu, Amazon Linux)
```

### Partition Count

```
Rule: partitions ≥ max(producer_parallelism, consumer_parallelism)

For 100 MB/s write from 10 producer instances:
  Each instance writes 10 MB/s
  Kafka partition can handle ~10 MB/s write
  → Need at least 10 partitions

For consumer parallelism of 20:
  → Need at least 20 partitions
  → Use 24 (nice multiple)
```

---

## Multi-Datacenter Setup

### Active-Active (MirrorMaker 2)

```
DC1: Kafka Cluster A           DC2: Kafka Cluster B
  Topic: orders       ←──MM2──►   Topic: A.orders
  Topic: A.orders ←──────────      Topic: orders
```

MirrorMaker 2 provides:
- Bidirectional replication
- Consumer offset translation (so consumers can failover)
- Topic renaming convention (`{source_cluster}.{topic}`)

### MirrorMaker 2 Config

```properties
# mm2.properties
clusters=us-east, eu-west

us-east.bootstrap.servers=kafka-us:9092
eu-west.bootstrap.servers=kafka-eu:9092

# Replication flows
us-east->eu-west.enabled=true
eu-west->us-east.enabled=true

# Which topics to replicate
us-east->eu-west.topics=orders,payments,events
eu-west->us-east.topics=orders,payments,events

# Replication factor in target cluster
replication.factor=3

# Replicate consumer group offsets (enables consumer failover)
us-east->eu-west.sync.group.offsets.enabled=true
eu-west->us-east.sync.group.offsets.enabled=true
```

```bash
# Start MirrorMaker 2
kafka-mirror-maker --mm2-config mm2.properties
# OR with Kafka Connect
connect-mirror-maker mm2.properties
```

---

## Disaster Recovery

### RTO and RPO Targets

```
RPO (Recovery Point Objective) = how much data can you lose?
  MirrorMaker 2 replicates near-real-time → RPO ~ seconds to minutes

RTO (Recovery Time Objective) = how fast can you recover?
  Depends on consumer failover speed
  With offset replication enabled → consumers resume from same position in minutes
```

### Failover Procedure

```bash
# 1. Stop MirrorMaker replication (prevent split-brain)
# Stop MM2 process

# 2. Update consumer bootstrap servers to DR cluster
# Deploy consumers pointing to eu-west brokers

# 3. Consumers use translated offsets (auto-translated by MM2)
# Consumer group offsets automatically available in DR cluster

# 4. Once DC1 is restored, replicate back and switchover
```

---

## Kafka on Kubernetes

### Strimzi Operator

```yaml
# kafka-cluster.yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: Kafka
metadata:
  name: production-cluster
  namespace: kafka
spec:
  kafka:
    version: 3.7.0
    replicas: 3
    listeners:
      - name: plain
        port: 9092
        type: internal
        tls: false
      - name: tls
        port: 9093
        type: internal
        tls: true
      - name: external
        port: 9094
        type: loadbalancer
        tls: true
    config:
      default.replication.factor: 3
      min.insync.replicas: 2
      offsets.topic.replication.factor: 3
      transaction.state.log.replication.factor: 3
      transaction.state.log.min.isr: 2
      auto.create.topics.enable: "false"
    storage:
      type: jbod
      volumes:
        - id: 0
          type: persistent-claim
          size: 2000Gi
          class: fast-ssd
          deleteClaim: false
    resources:
      requests:
        memory: 16Gi
        cpu: "4"
      limits:
        memory: 16Gi
        cpu: "8"
    rack:
      topologyKey: topology.kubernetes.io/zone
  zookeeper:
    replicas: 3
    storage:
      type: persistent-claim
      size: 100Gi
      class: fast-ssd
```

---

## Confluent Cloud / MSK

### Confluent Cloud Python Client

```python
producer = Producer({
    "bootstrap.servers":    "pkc-xxxxx.us-east-1.aws.confluent.cloud:9092",
    "security.protocol":    "SASL_SSL",
    "sasl.mechanism":       "PLAIN",
    "sasl.username":        "API_KEY",
    "sasl.password":        "API_SECRET",
    "acks":                 "all",
    "enable.idempotence":   True,
})
```

### AWS MSK (Managed Streaming for Kafka)

```python
# IAM authentication (MSK Serverless / MSK with IAM)
producer = Producer({
    "bootstrap.servers":                         "b-1.msk-cluster.us-east-1.amazonaws.com:9098",
    "security.protocol":                          "SASL_SSL",
    "sasl.mechanism":                             "AWS_MSK_IAM",
    "sasl.jaas.config":                           "software.amazon.msk.auth.iam.IAMLoginModule required;",
    "sasl.client.callback.handler.class":         "software.amazon.msk.auth.iam.IAMClientCallbackHandler",
})
```

---

## Capacity Planning

```
Daily growth estimate:
  events_per_second × avg_event_bytes × 86400 × replication_factor

Example: 10,000 events/s × 500 bytes × 86400 × 3 = 1.296 TB/day

Disk capacity needed (7 days retention):
  1.296 TB/day × 7 = ~9 TB raw

With compression (snappy ~50% ratio):
  ~4.5 TB compressed

Add 20% overhead buffer: ~5.4 TB across brokers
```

---

## Production Checklist

| Area | Setting | Value |
|---|---|---|
| Replication | `replication.factor` | 3 |
| Replication | `min.insync.replicas` | 2 |
| Reliability | `acks` (producer) | `all` |
| Reliability | `enable.idempotence` | `true` |
| Safety | `unclean.leader.election.enable` | `false` |
| Topics | `auto.create.topics.enable` | `false` |
| Security | TLS + SASL | Enabled |
| Security | ACLs | Least privilege |
| Monitoring | Consumer lag alert | Configurable threshold |
| Monitoring | Under-replicated partitions alert | 0 expected |
| OS | `vm.swappiness` | `1` |
| OS | Max open files | `100000+` |
| JVM | Heap | `6g max, G1GC` |
| Storage | XFS with `noatime` | Enabled |
| Retention | per-topic | Based on use case |

---

## Capacity Monitoring and Alerts

```bash
# Prometheus alert rules
- alert: KafkaUnderReplicatedPartitions
  expr: kafka_server_replicamanager_underreplicatedpartitions > 0
  for: 1m
  annotations:
    summary: "Under-replicated partitions detected"

- alert: KafkaConsumerLagHigh
  expr: kafka_consumergroup_lag > 100000
  for: 5m
  annotations:
    summary: "Consumer group {{ $labels.consumergroup }} is lagging"

- alert: KafkaBrokerDown
  expr: count(kafka_server_kafkaserver_brokerstate == 3) < 3
  annotations:
    summary: "Fewer than 3 healthy Kafka brokers"

- alert: KafkaDiskUsageHigh
  expr: (node_filesystem_size_bytes{mountpoint="/var/kafka"} -
         node_filesystem_free_bytes{mountpoint="/var/kafka"}) /
        node_filesystem_size_bytes{mountpoint="/var/kafka"} > 0.80
  annotations:
    summary: "Kafka disk >80% full"
```
