# Apache Kafka — Complete Guide

> All 20 chapters in one file. From fundamentals to production-grade distributed streaming.

---

## Quick Navigation

- [Index & Architecture Overview](#index-architecture-overview)
- [Introduction to Apache Kafka](#introduction-to-apache-kafka)
- [Setup & Installation](#setup-installation)
- [Core Concepts](#core-concepts)
- [Producers](#producers)
- [Consumers & Consumer Groups](#consumers-consumer-groups)
- [Topics & Partitions](#topics-partitions)
- [Kafka Streams](#kafka-streams)
- [Kafka Connect](#kafka-connect)
- [Schema Registry](#schema-registry)
- [Cluster Architecture](#cluster-architecture)
- [Replication & Durability](#replication-durability)
- [Performance & Tuning](#performance-tuning)
- [Security](#security)
- [Monitoring & Operations](#monitoring-operations)
- [Kafka with Python](#kafka-with-python)
- [Kafka with Java](#kafka-with-java)
- [Exactly-Once Semantics](#exactly-once-semantics)
- [Kafka in Production](#kafka-in-production)
- [Kafka with Spark & Flink](#kafka-with-spark-flink)
- [Learning Roadmap](#learning-roadmap)

---

# Apache Kafka — Complete Learning Guide

> A comprehensive, hands-on guide to Apache Kafka: from core concepts to production-grade distributed streaming systems.

---

## Table of Contents

| Chapter | Title | Topics |
|---|---|---|
| 01 | [Introduction](01_introduction.md) | What is Kafka, why it exists, pub/sub model, use cases |
| 02 | [Setup & Installation](02_setup_and_installation.md) | Docker, KRaft mode, multi-broker, CLI tools |
| 03 | [Core Concepts](03_core_concepts.md) | Topics, partitions, offsets, brokers, producers, consumers |
| 04 | [Producers](04_producers.md) | Producer API, partitioning, acks, retries, idempotency |
| 05 | [Consumers & Consumer Groups](05_consumers_and_consumer_groups.md) | Consumer API, group rebalance, offset management, lag |
| 06 | [Topics & Partitions](06_topics_and_partitions.md) | Topic config, compaction, retention, replication |
| 07 | [Kafka Streams](07_kafka_streams.md) | KStream, KTable, aggregations, joins, windowing |
| 08 | [Kafka Connect](08_kafka_connect.md) | Source/sink connectors, distributed mode, SMTs |
| 09 | [Schema Registry](09_schema_registry.md) | Avro, Protobuf, JSON Schema, compatibility modes |
| 10 | [Cluster Architecture](10_cluster_architecture.md) | ZooKeeper vs KRaft, ISR, leader election, controller |
| 11 | [Replication & Durability](11_replication_and_durability.md) | Replication factor, acks, ISR, min.insync.replicas |
| 12 | [Performance & Tuning](12_performance_and_tuning.md) | Batching, compression, page cache, OS tuning |
| 13 | [Security](13_security.md) | TLS, SASL, ACLs, encryption at rest |
| 14 | [Monitoring & Operations](14_monitoring_and_operations.md) | JMX, Prometheus, consumer lag, key metrics |
| 15 | [Kafka with Python](15_kafka_with_python.md) | confluent-kafka, kafka-python, Faust streaming |
| 16 | [Kafka with Java](16_kafka_with_java.md) | Java Producer/Consumer API, Spring Kafka |
| 17 | [Exactly-Once Semantics](17_exactly_once_semantics.md) | Idempotent producer, transactions, EOS patterns |
| 18 | [Kafka in Production](18_kafka_in_production.md) | Multi-DC, MirrorMaker 2, capacity planning, disaster recovery |
| 19 | [Kafka with Spark & Flink](19_kafka_with_spark_and_flink.md) | Spark Structured Streaming, Flink Kafka connector |
| 20 | [Learning Roadmap](20_learning_roadmap.md) | 5-week plan, capstone projects, beyond Kafka |

---

## How to Read This Guide

- **Beginners**: Chapters 1–6 cover fundamentals. Complete the practice exercises before moving on.
- **Application developers**: Focus on Chapters 4–5, 9, 15–16 for producer/consumer patterns.
- **Platform engineers**: Chapters 10–14, 18 cover cluster operations and production readiness.
- **Stream processors**: Chapters 7–8, 17, 19 cover stateful processing and exactly-once delivery.
- **All readers**: Chapter 20 provides a structured 5-week hands-on roadmap.

---

## Apache Kafka Architecture

```
                         ┌─────────────────────────────────────────────┐
                         │              Apache Kafka Cluster            │
                         │                                             │
  Producer A ──────────► │  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
  Producer B ──────────► │  │ Broker 1 │  │ Broker 2 │  │ Broker 3 │ │
  Producer C ──────────► │  │(Leader:  │  │(Leader:  │  │(Leader:  │ │
                         │  │ P0,P3)   │  │ P1,P4)   │  │ P2,P5)   │ │
                         │  └──────────┘  └──────────┘  └──────────┘ │
                         │        │              │              │      │
                         │  ┌─────┴──────────────┴──────────────┴───┐ │
                         │  │           ZooKeeper / KRaft            │ │
                         │  │        (Cluster Metadata & Leader      │ │
                         │  │              Election)                 │ │
                         │  └────────────────────────────────────────┘ │
                         └─────────────────────────────────────────────┘
                                  ▲                  ▼
                         Consumer Group A      Consumer Group B
                         (App: Payments)       (App: Analytics)
                         ├── Consumer 1        ├── Consumer X
                         └── Consumer 2        └── Consumer Y
```

---

## Core Abstractions

| Concept | Description |
|---|---|
| **Topic** | Named, append-only log of events |
| **Partition** | Ordered sub-log within a topic; unit of parallelism |
| **Offset** | Position of a message within a partition |
| **Broker** | Kafka server that stores and serves messages |
| **Producer** | Client that writes messages to topics |
| **Consumer** | Client that reads messages from topics |
| **Consumer Group** | Set of consumers sharing partitions for parallel reads |
| **Replication** | Copies of partitions across brokers for fault tolerance |
| **ISR** | In-Sync Replicas — replicas caught up with the leader |


---

# Chapter 1: Introduction to Apache Kafka

## What is Apache Kafka?

Apache Kafka is a **distributed event streaming platform** built for high-throughput, fault-tolerant, real-time data pipelines. Originally developed at LinkedIn in 2011 and open-sourced through the Apache Software Foundation, Kafka is now the backbone of data infrastructure at thousands of companies worldwide.

At its core, Kafka is a **distributed commit log**: an ordered, immutable sequence of events that can be read by many consumers without removing the data.

---

## The Problem Kafka Solves

Before Kafka, companies faced a **point-to-point integration explosion**:

```
Without Kafka:
  Service A ──────► Service B
  Service A ──────► Service C
  Service A ──────► Service D
  Service B ──────► Service C
  Service B ──────► Service E
  ... N×(N-1) connections for N services

With Kafka:
  Service A ──► Kafka ◄── Service B
  Service C ──►        ◄── Service D
               ▼
         Any consumer
         subscribes once
```

**Problems in a pre-Kafka world:**
- Every new integration requires custom code on both sides
- No replay — if a consumer is down, events are lost
- Different systems have different protocols (HTTP, TCP, queues)
- No historical data — can't re-process past events
- Tight coupling — producer must know every consumer

---

## Kafka's Key Properties

| Property | Description |
|---|---|
| **High throughput** | Millions of messages per second per broker |
| **Low latency** | Sub-10ms end-to-end at scale |
| **Durable** | Messages persisted to disk with replication |
| **Fault-tolerant** | Survives broker failures with no data loss |
| **Scalable** | Horizontal scale by adding brokers and partitions |
| **Replayable** | Consumers can re-read from any point in history |
| **Decoupled** | Producers and consumers are fully independent |

---

## The Publish-Subscribe Model

Kafka uses a **publish-subscribe** model:

```
PUBLISHER (Producer)                SUBSCRIBER (Consumer)
     │                                       │
     │  write("orders", key="user-1",        │
     │         value={"amount": 99.99})      │
     ▼                                       │
┌─────────────────────────────┐              │
│  Topic: orders              │              │
│  Partition 0: [msg0][msg1]  │◄─────────────┤ read from offset 42
│  Partition 1: [msg2][msg3]  │◄─────────────┘
│  Partition 2: [msg4][msg5]  │
└─────────────────────────────┘
```

Key insight: **multiple consumer groups** can each independently read the same topic. Adding a new consumer group never impacts producers or other consumers.

---

## Kafka vs Traditional Message Queues

| Feature | Traditional MQ (RabbitMQ, SQS) | Apache Kafka |
|---|---|---|
| Message deletion | Deleted after consumption | Retained by time/size policy |
| Multiple consumers | Competing consumers (one gets each) | Each group gets all messages |
| Replay | Not possible | Yes — seek to any offset |
| Ordering | Per-queue (limited) | Per-partition (guaranteed) |
| Throughput | Moderate | Extremely high (millions/s) |
| Persistence | Short-lived | Long-term (days to forever) |
| Use case | Task queues, RPC | Event streaming, CDC, pipelines |

---

## Core Use Cases

### 1. Real-Time Event Streaming
```
User clicks ──► Kafka ──► Recommendation engine (real-time)
                      ──► Analytics pipeline (batch)
                      ──► Fraud detection (real-time ML)
```

### 2. Change Data Capture (CDC)
```
PostgreSQL WAL ──► Debezium Connector ──► Kafka ──► Data warehouse
                                                ──► Search index (ES)
                                                ──► Cache invalidation
```

### 3. Log Aggregation
```
App Server 1 ──►
App Server 2 ──► Kafka (logs topic) ──► Elasticsearch
App Server 3 ──►                   ──► S3 / HDFS
```

### 4. Metrics & Monitoring
```
IoT sensors ──► Kafka ──► Spark Streaming ──► Dashboards
                      ──► Alerting system
```

### 5. Event Sourcing
```
Every state change = an event in Kafka
Current state = replay of all past events
Audit log = free, built into Kafka retention
```

---

## Kafka Ecosystem

```
┌─────────────────────────────────────────────────────────┐
│                    Kafka Ecosystem                       │
│                                                         │
│  ┌───────────┐  ┌─────────────┐  ┌──────────────────┐  │
│  │  Kafka    │  │   Kafka     │  │  Schema Registry  │  │
│  │  Streams  │  │   Connect   │  │  (Avro/Protobuf)  │  │
│  └───────────┘  └─────────────┘  └──────────────────┘  │
│                                                         │
│  ┌─────────────────────────────────────────────────┐    │
│  │               Apache Kafka Core                 │    │
│  │  (Topics, Partitions, Brokers, Replication)     │    │
│  └─────────────────────────────────────────────────┘    │
│                                                         │
│  ┌──────────┐  ┌──────────┐  ┌────────┐  ┌─────────┐   │
│  │  Flink   │  │  Spark   │  │  ksqlDB│  │ Debezium│   │
│  └──────────┘  └──────────┘  └────────┘  └─────────┘   │
└─────────────────────────────────────────────────────────┘
```

| Component | Purpose |
|---|---|
| **Kafka Core** | Broker, producer, consumer, topic management |
| **Kafka Streams** | Stateful stream processing in Java/Kotlin |
| **Kafka Connect** | Plug-and-play source/sink connectors |
| **Schema Registry** | Enforce schemas, prevent breaking changes |
| **ksqlDB** | SQL interface for stream processing |
| **Flink / Spark** | Advanced stream processing frameworks |
| **Debezium** | CDC from databases (Postgres, MySQL, MongoDB) |

---

## Kafka vs Competitors

| Feature | Kafka | Pulsar | Kinesis | RabbitMQ |
|---|---|---|---|---|
| Throughput | Very high | Very high | High | Moderate |
| Multi-tenancy | Limited | Native | N/A | Virtual hosts |
| Replay | Configurable | Configurable | 7 days max | No |
| Protocol | Custom binary | Custom + multi | Proprietary | AMQP |
| Managed cloud | Confluent, MSK | StreamNative | Kinesis | CloudAMQP |
| Tiered storage | Yes (KIP-405) | Native | Yes | No |
| Streams API | Kafka Streams | Pulsar Functions | Lambda | N/A |

---

## How Kafka Stores Data

Kafka stores data in **topic-partitions** on disk as segment files:

```
/kafka-logs/
  orders-0/         ← Partition 0 of "orders" topic
    00000000000000000000.log       ← messages
    00000000000000000000.index     ← offset → byte position
    00000000000000000000.timeindex ← timestamp → offset
    00000000000000012345.log       ← next segment
  orders-1/
  orders-2/
```

**Key storage properties:**
- Sequential disk writes (fast, no random I/O)
- OS page cache used extensively for reads
- Messages are immutable — never updated, only appended
- Compaction: can keep only latest value per key
- Retention: configurable by time (e.g., 7 days) or size (e.g., 100 GB)


---

# Chapter 2: Setup & Installation

## Quick Start with Docker (Recommended)

### Single-Node KRaft Mode (No ZooKeeper)

```yaml
# docker-compose.yml
version: "3.8"
services:
  kafka:
    image: confluentinc/cp-kafka:7.6.0
    hostname: kafka
    container_name: kafka
    ports:
      - "9092:9092"
      - "9101:9101"
    environment:
      KAFKA_NODE_ID: 1
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_PROCESS_ROLES: broker,controller
      KAFKA_CONTROLLER_QUORUM_VOTERS: 1@kafka:29093
      KAFKA_LISTENERS: PLAINTEXT://kafka:29092,CONTROLLER://kafka:29093,PLAINTEXT_HOST://0.0.0.0:9092
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_LOG_DIRS: /tmp/kraft-combined-logs
      CLUSTER_ID: MkU3OEVBNTcwNTJENDM2Qk

  schema-registry:
    image: confluentinc/cp-schema-registry:7.6.0
    ports:
      - "8081:8081"
    environment:
      SCHEMA_REGISTRY_HOST_NAME: schema-registry
      SCHEMA_REGISTRY_KAFKASTORE_BOOTSTRAP_SERVERS: kafka:29092

  kafka-ui:
    image: provectuslabs/kafka-ui:latest
    ports:
      - "8080:8080"
    environment:
      KAFKA_CLUSTERS_0_NAME: local
      KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:29092
      KAFKA_CLUSTERS_0_SCHEMAREGISTRY: http://schema-registry:8081
```

```bash
docker compose up -d
```

---

### Full Cluster with ZooKeeper (3 Brokers)

```yaml
# docker-compose-cluster.yml
version: "3.8"

services:
  zookeeper:
    image: confluentinc/cp-zookeeper:7.6.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000

  broker1:
    image: confluentinc/cp-kafka:7.6.0
    depends_on: [zookeeper]
    ports: ["9092:9092"]
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://broker1:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_DEFAULT_REPLICATION_FACTOR: 3
      KAFKA_MIN_INSYNC_REPLICAS: 2
      KAFKA_LOG_RETENTION_HOURS: 168

  broker2:
    image: confluentinc/cp-kafka:7.6.0
    depends_on: [zookeeper]
    ports: ["9093:9093"]
    environment:
      KAFKA_BROKER_ID: 2
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://broker2:29092,PLAINTEXT_HOST://localhost:9093
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT

  broker3:
    image: confluentinc/cp-kafka:7.6.0
    depends_on: [zookeeper]
    ports: ["9094:9094"]
    environment:
      KAFKA_BROKER_ID: 3
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://broker3:29092,PLAINTEXT_HOST://localhost:9094
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
```

---

## CLI Tools

### Topic Management

```bash
# Create topic
kafka-topics --bootstrap-server localhost:9092 \
  --create --topic orders \
  --partitions 6 \
  --replication-factor 1

# List topics
kafka-topics --bootstrap-server localhost:9092 --list

# Describe topic
kafka-topics --bootstrap-server localhost:9092 --describe --topic orders

# Delete topic
kafka-topics --bootstrap-server localhost:9092 --delete --topic orders

# Alter topic (increase partitions only)
kafka-topics --bootstrap-server localhost:9092 \
  --alter --topic orders --partitions 12
```

### Producer CLI

```bash
# Simple producer (type messages interactively)
kafka-console-producer --bootstrap-server localhost:9092 --topic orders

# Producer with keys
kafka-console-producer --bootstrap-server localhost:9092 \
  --topic orders \
  --property "key.separator=:" \
  --property "parse.key=true"

# Type: user-123:{"amount": 99.99}
```

### Consumer CLI

```bash
# Read from beginning
kafka-console-consumer --bootstrap-server localhost:9092 \
  --topic orders \
  --from-beginning

# Consumer group
kafka-console-consumer --bootstrap-server localhost:9092 \
  --topic orders \
  --group my-app \
  --from-beginning

# Show keys
kafka-console-consumer --bootstrap-server localhost:9092 \
  --topic orders \
  --from-beginning \
  --property print.key=true \
  --property key.separator=" | "

# Read with Avro (Schema Registry)
kafka-avro-console-consumer --bootstrap-server localhost:9092 \
  --topic orders \
  --property schema.registry.url=http://localhost:8081
```

### Consumer Group Management

```bash
# List consumer groups
kafka-consumer-groups --bootstrap-server localhost:9092 --list

# Describe group (shows lag!)
kafka-consumer-groups --bootstrap-server localhost:9092 \
  --describe --group my-app

# Reset offsets to beginning
kafka-consumer-groups --bootstrap-server localhost:9092 \
  --group my-app \
  --topic orders \
  --reset-offsets --to-earliest --execute

# Reset offsets to specific offset
kafka-consumer-groups --bootstrap-server localhost:9092 \
  --group my-app \
  --topic orders:2 \
  --reset-offsets --to-offset 100 --execute

# Reset offsets by time
kafka-consumer-groups --bootstrap-server localhost:9092 \
  --group my-app --topic orders \
  --reset-offsets --to-datetime 2024-01-01T00:00:00.000 --execute
```

---

## kafka.properties (server.properties Key Settings)

```properties
# Broker identity
broker.id=1
node.id=1                           # KRaft mode

# Networking
listeners=PLAINTEXT://0.0.0.0:9092
advertised.listeners=PLAINTEXT://kafka-host:9092

# Storage
log.dirs=/var/kafka/logs
log.retention.hours=168             # 7 days
log.retention.bytes=107374182400    # 100 GB per partition
log.segment.bytes=1073741824        # 1 GB per segment

# Replication
default.replication.factor=3
min.insync.replicas=2

# Performance
num.network.threads=3
num.io.threads=8
socket.send.buffer.bytes=102400
socket.receive.buffer.bytes=102400
num.partitions=6

# Log compaction
log.cleanup.policy=delete           # or compact
log.compaction.lag.ms=3600000
```

---

## First Steps: Verify Installation

```bash
# 1. Create a test topic
kafka-topics --bootstrap-server localhost:9092 \
  --create --topic test --partitions 1 --replication-factor 1

# 2. Produce a message (in one terminal)
echo "hello kafka" | kafka-console-producer \
  --bootstrap-server localhost:9092 --topic test

# 3. Consume the message (in another terminal)
kafka-console-consumer --bootstrap-server localhost:9092 \
  --topic test --from-beginning

# 4. Check broker metadata
kafka-broker-api-versions --bootstrap-server localhost:9092

# 5. Check cluster info
kafka-metadata-quorum --bootstrap-server localhost:9092 --describe --status
```

---

## Kafka UI Tools

| Tool | Type | Notes |
|---|---|---|
| **Kafka UI** (provectuslabs) | Web UI | Open source, Docker-friendly |
| **Conduktor** | Desktop/Web | Commercial, powerful |
| **AKHQ** | Web UI | Open source, feature-rich |
| **Kafka Magic** | Desktop | Simple GUI |
| **Confluent Control Center** | Web UI | Enterprise, part of Confluent Platform |
| **kcat (kafkacat)** | CLI | Power-user CLI tool |

```bash
# kcat examples
# Produce
echo '{"user": "alice"}' | kcat -P -b localhost:9092 -t events -k user-1

# Consume (JSON pretty print)
kcat -C -b localhost:9092 -t events -o beginning | jq .

# Metadata
kcat -b localhost:9092 -L
```


---

# Chapter 3: Core Concepts

## Topics

A **topic** is a named, ordered, append-only log. Think of it like a database table — but immutable and distributed.

```
Topic: "orders"

Partition 0: [offset 0][offset 1][offset 2][offset 3] ──► append
Partition 1: [offset 0][offset 1][offset 2]           ──► append
Partition 2: [offset 0][offset 1][offset 2][offset 3] ──► append
```

**Key properties:**
- Topics are split into partitions for parallelism
- Within a partition, messages are strictly ordered
- Across partitions, there is no ordering guarantee
- Retention is configurable (time-based or size-based)

---

## Partitions

**Partitions** are the unit of parallelism and scalability in Kafka.

```
┌────────────────────────────────────────────┐
│ Topic: user-events   (4 partitions)        │
│                                            │
│ P0: [e0][e1][e2][e3][e4]                  │
│ P1: [e0][e1][e2]                          │
│ P2: [e0][e1][e2][e3]                      │
│ P3: [e0]                                  │
└────────────────────────────────────────────┘

Each partition is an independent ordered log.
Partitions can be distributed across brokers.
```

**Rules:**
- You can increase partitions but NEVER decrease them
- More partitions = more parallelism (up to a point)
- Partition count should equal max consumer instances in a group

---

## Offsets

An **offset** is a monotonically increasing integer that identifies the position of a message within a partition.

```
Partition 0:
  Offset 0: {"order_id": "A1", "amount": 50.00}
  Offset 1: {"order_id": "A2", "amount": 75.00}
  Offset 2: {"order_id": "A3", "amount": 22.50}  ← consumer at offset 2
  Offset 3: {"order_id": "A4", "amount": 99.99}  ← new, unread
```

**Consumer offset tracking:**
- Each consumer group tracks its own offset per partition
- Stored in internal topic `__consumer_offsets`
- Enables independent consumption by multiple groups
- Replay = reset offset to an earlier position

---

## Messages (Records)

Every Kafka message has these fields:

| Field | Type | Description |
|---|---|---|
| **Key** | bytes (nullable) | Routing key; determines partition |
| **Value** | bytes (nullable) | The message payload |
| **Topic** | string | Destination topic name |
| **Partition** | int | Which partition (set by broker/producer) |
| **Offset** | long | Position within the partition |
| **Timestamp** | long | Create time or log-append time |
| **Headers** | Map<String, bytes> | Optional metadata (trace ID, source, etc.) |

```python
# A Kafka record in Python
record = {
    "key":       b"user-123",
    "value":     b'{"amount": 99.99, "currency": "USD"}',
    "headers":   [("source", b"checkout-service"), ("version", b"2")],
    "timestamp": 1704067200000,  # ms since epoch
}
```

---

## Brokers

A **broker** is a Kafka server that:
- Stores partitions on disk
- Handles read/write requests from clients
- Replicates partitions to other brokers
- Participates in leader election

```
Cluster: 3 brokers

Broker 1 (leader for P0, P3):
  /kafka-logs/orders-0/
  /kafka-logs/orders-3/
  /kafka-logs/events-1/  ← replica

Broker 2 (leader for P1, P4):
  /kafka-logs/orders-1/
  /kafka-logs/orders-4/

Broker 3 (leader for P2, P5):
  /kafka-logs/orders-2/
  /kafka-logs/orders-5/
```

---

## Producers

A **producer** writes messages to Kafka topics.

```python
from confluent_kafka import Producer

p = Producer({"bootstrap.servers": "localhost:9092"})

# Async produce (non-blocking)
p.produce("orders", key="user-1", value='{"amount": 50}')
p.flush()  # wait for all in-flight messages to be delivered
```

**Partitioning logic:**
1. If a partition is specified explicitly → use it
2. If a key is provided → `hash(key) % num_partitions`
3. No key → round-robin across partitions (sticky by batch)

---

## Consumers

A **consumer** reads messages from Kafka topics.

```python
from confluent_kafka import Consumer

c = Consumer({
    "bootstrap.servers": "localhost:9092",
    "group.id":          "my-app",
    "auto.offset.reset": "earliest",
})
c.subscribe(["orders"])

while True:
    msg = c.poll(timeout=1.0)
    if msg is None: continue
    if msg.error(): print(f"Error: {msg.error()}"); continue

    print(f"Received: key={msg.key()}, value={msg.value()}, offset={msg.offset()}")

c.close()
```

---

## Consumer Groups

**Consumer groups** enable parallel consumption:

```
Topic: orders (6 partitions)

Consumer Group A (3 consumers):
  Consumer A1 ──► P0, P1
  Consumer A2 ──► P2, P3
  Consumer A3 ──► P4, P5

Consumer Group B (2 consumers):
  Consumer B1 ──► P0, P1, P2
  Consumer B2 ──► P3, P4, P5

Both groups get ALL messages independently.
```

**Rules:**
- Each partition is consumed by exactly one consumer per group
- If consumers > partitions, some consumers are idle
- If consumers < partitions, some consumers handle multiple partitions
- Adding/removing consumers triggers a **rebalance**

---

## Replication

Kafka replicates partitions for fault tolerance:

```
Topic: orders, Partition 0, Replication Factor: 3

Broker 1: [LEADER]   ← handles all reads & writes
Broker 2: [FOLLOWER] ← replicates from leader (ISR)
Broker 3: [FOLLOWER] ← replicates from leader (ISR)

If Broker 1 fails:
  Kafka elects new leader (Broker 2 or 3)
  Zero data loss (if ISR was up to date)
```

**ISR (In-Sync Replicas):** The set of replicas fully caught up with the leader. Only ISR members can become the new leader.

---

## The Commit Log Model

```
Time flows right ──────────────────────────────────►

P0: [A][B][C][D][E][F][G][H]
                 ▲           ▲
           Oldest msg    New writes appended here
           (may be
           compacted
           or deleted
           by retention)
```

This model gives Kafka its power:
- **Immutability**: messages never change
- **Replayability**: seek to any offset and re-read history
- **Decoupling**: readers don't affect writers
- **Speed**: sequential disk I/O is extremely fast

---

## ZooKeeper vs KRaft

| Feature | ZooKeeper Mode | KRaft Mode (Kafka 3.3+) |
|---|---|---|
| Metadata storage | ZooKeeper | Kafka's own Raft log |
| Dependency | External ZK cluster | Built in |
| Partition limit | ~200,000 | Millions |
| Controller | ZooKeeper-based | Self-managed |
| Recovery time | Seconds-minutes | Milliseconds |
| Production ready | Yes (legacy) | Yes (Kafka 3.3+) |

KRaft mode is the future — ZooKeeper support will be removed in a future Kafka version.


---

# Chapter 4: Producers

## Producer Overview

The Kafka producer is responsible for publishing records to Kafka topics. Under the hood, the producer:

1. Serializes the key and value to bytes
2. Determines the target partition (via partitioner)
3. Batches records into a `RecordAccumulator`
4. Sends batches to the broker in background threads
5. Handles retries and error callbacks

```
Producer Application
       │
  ┌────▼──────────────────────────────────────────┐
  │                 Producer Client               │
  │                                               │
  │  serialize() ──► partition() ──► accumulate() │
  │                                      │        │
  │                               ┌──────▼──────┐ │
  │                               │  Batch P0   │ │
  │                               │  Batch P1   │ │ ──► Broker
  │                               │  Batch P2   │ │
  │                               └─────────────┘ │
  └───────────────────────────────────────────────┘
```

---

## Basic Producer

```python
from confluent_kafka import Producer

producer = Producer({
    "bootstrap.servers": "localhost:9092",
})

def delivery_report(err, msg):
    if err:
        print(f"Delivery failed: {err}")
    else:
        print(f"Delivered to {msg.topic()} [{msg.partition()}] @ offset {msg.offset()}")

# Async produce
producer.produce(
    topic    = "orders",
    key      = "user-123",
    value    = '{"order_id": "A1", "amount": 99.99}',
    headers  = {"source": "checkout"},
    callback = delivery_report,
)

# Flush all pending messages
producer.flush()
```

---

## Producer Configuration

### Reliability Settings

```python
producer = Producer({
    "bootstrap.servers": "broker1:9092,broker2:9092,broker3:9092",

    # Acknowledgement mode
    "acks":              "all",    # 0=fire-forget, 1=leader-only, all=full ISR
    "enable.idempotence": True,    # exactly-once per partition (requires acks=all)

    # Retry behavior
    "retries":              10,
    "retry.backoff.ms":     100,
    "delivery.timeout.ms":  120000,   # total retry window

    # In-flight requests (keep at 1 for strict ordering without idempotence)
    "max.in.flight.requests.per.connection": 5,

    # Compression
    "compression.type":  "snappy",   # none, gzip, snappy, lz4, zstd
})
```

### Throughput Tuning

```python
producer = Producer({
    "bootstrap.servers": "localhost:9092",

    # Batching
    "batch.size":           65536,   # 64KB batch size in bytes
    "linger.ms":            10,      # wait up to 10ms to fill a batch
    "buffer.memory":        33554432, # 32MB total buffer

    # Compression (reduces network + storage)
    "compression.type":     "lz4",
})
```

---

## Partitioning Strategies

### Default: Key Hash

```python
# Same key always goes to the same partition
producer.produce("orders", key="user-123", value="order data")
# hash("user-123") % num_partitions = consistent partition
```

### Custom Partitioner (Python)

```python
def my_partitioner(key: bytes, all_partitions: list, available: list) -> int:
    if key is None:
        return all_partitions[0]
    # Route VIP orders to partition 0
    if key.startswith(b"vip-"):
        return 0
    return hash(key) % len(all_partitions)

producer = Producer({
    "bootstrap.servers": "localhost:9092",
    "partitioner":        my_partitioner,
})
```

### Round-Robin (No Key)

```python
# Without a key, Kafka uses sticky partitioning per batch
producer.produce("events", value="no key message")
```

---

## Serialization

Kafka transmits raw bytes — you handle serialization:

```python
import json
from confluent_kafka.serialization import StringSerializer, JSONSerializer
from confluent_kafka.schema_registry import SchemaRegistryClient
from confluent_kafka.schema_registry.json_schema import JSONSerializer as SchemaJSONSerializer

# Simple JSON
def json_serialize(obj: dict) -> bytes:
    return json.dumps(obj).encode("utf-8")

producer.produce(
    "orders",
    key   = "user-1".encode(),
    value = json_serialize({"amount": 99.99}),
)

# Using Avro with Schema Registry (see Chapter 9)
from confluent_kafka.schema_registry.avro import AvroSerializer

schema_registry = SchemaRegistryClient({"url": "http://localhost:8081"})
avro_serializer = AvroSerializer(schema_registry, schema_str)

producer.produce(
    "orders",
    value = avro_serializer({"amount": 99.99}, SerializationContext("orders", MessageField.VALUE))
)
```

---

## Acks Explained

```
acks=0  (fire and forget)
  Producer ──► Broker
  No confirmation. Fastest. Data loss possible.

acks=1  (leader ack)
  Producer ──► Broker Leader ──► ACK
  Leader confirms write. Follower lag can cause loss on crash.

acks=all (strongest)
  Producer ──► Broker Leader ──► ACK after all ISR replicas confirm
  Safest. Highest latency. No data loss if ISR ≥ min.insync.replicas.
```

---

## Idempotent Producer

Prevents duplicate messages on retry:

```python
producer = Producer({
    "bootstrap.servers": "localhost:9092",
    "enable.idempotence": True,
    # These are automatically set when idempotence is enabled:
    # acks = all
    # max.in.flight.requests.per.connection = 5
    # retries = INT_MAX
})
```

How it works:
- Each producer gets a **PID** (Producer ID) from the broker
- Each message carries a **sequence number**
- Broker deduplicates by (PID, partition, sequence number)
- Guarantees: **at least once** → **exactly once per partition**

---

## Transactions (Atomic Multi-Topic Writes)

```python
from confluent_kafka import Producer

producer = Producer({
    "bootstrap.servers":     "localhost:9092",
    "transactional.id":      "my-producer-1",
    "enable.idempotence":     True,
})

producer.init_transactions()

try:
    producer.begin_transaction()

    # Write to multiple topics atomically
    producer.produce("orders",   key="order-1",   value="created")
    producer.produce("inventory", key="product-A", value="decrement")
    producer.produce("payments",  key="order-1",   value="charged")

    producer.commit_transaction()

except Exception as e:
    producer.abort_transaction()
    raise
```

Transactions guarantee **all-or-nothing** across multiple topics and partitions.

---

## Producer Metrics to Monitor

| Metric | What it Tells You |
|---|---|
| `record-send-rate` | Messages/second being produced |
| `record-error-rate` | Failed deliveries per second |
| `request-latency-avg` | Average time to get a broker ACK |
| `batch-size-avg` | Average bytes per batch (higher = more efficient) |
| `record-queue-time-avg` | Time records wait in accumulator before sending |
| `compression-rate-avg` | Compression effectiveness |
| `buffer-available-bytes` | Free space in send buffer |

---

## Producer Best Practices

```python
# 1. Always flush() before shutdown
import atexit
atexit.register(producer.flush)

# 2. Use callbacks for error handling
def on_delivery(err, msg):
    if err:
        # Retry, dead-letter queue, or alert
        logger.error(f"Produce failed: {err}")

# 3. Batch small messages
producer = Producer({
    "batch.size":   65536,
    "linger.ms":    5,
    "compression.type": "lz4",
})

# 4. Monitor buffer memory
# If buffer fills up, produce() blocks or raises BufferError
producer = Producer({
    "queue.buffering.max.messages": 100000,
    "queue.buffering.max.kbytes":   1048576,  # 1 GB
})
```


---

# Chapter 5: Consumers & Consumer Groups

## Consumer Basics

A consumer reads records from Kafka topics and tracks its position (offset) per partition.

```python
from confluent_kafka import Consumer, KafkaError

consumer = Consumer({
    "bootstrap.servers":  "localhost:9092",
    "group.id":           "order-processor",
    "auto.offset.reset":  "earliest",   # start from beginning if no committed offset
    "enable.auto.commit": False,        # manual offset control (recommended)
})

consumer.subscribe(["orders"])

try:
    while True:
        msg = consumer.poll(timeout=1.0)

        if msg is None:
            continue
        if msg.error():
            if msg.error().code() == KafkaError._PARTITION_EOF:
                print(f"Reached end of partition {msg.partition()}")
            else:
                raise Exception(msg.error())
            continue

        # Process message
        print(f"[{msg.topic()}|P{msg.partition()}|@{msg.offset()}] "
              f"key={msg.key()} value={msg.value()}")

        # Commit offset AFTER successful processing
        consumer.commit(asynchronous=False)

finally:
    consumer.close()
```

---

## Consumer Groups

**Consumer groups** provide scalable, fault-tolerant consumption.

```
Topic: orders (6 partitions: P0–P5)

Group "payment-service" (3 instances):
  instance-1: P0, P1
  instance-2: P2, P3
  instance-3: P4, P5

Group "analytics" (1 instance):
  instance-1: P0, P1, P2, P3, P4, P5

Both groups receive ALL messages independently.
```

Each consumer group maintains its own offsets in `__consumer_offsets`.

---

## Rebalancing

A **rebalance** happens when the group membership changes:
- New consumer joins
- Consumer crashes or leaves
- Consumer fails to heartbeat in time
- New partitions added to a topic

```
Before rebalance (3 consumers, 6 partitions):
  C1: P0, P1
  C2: P2, P3
  C3: P4, P5

C3 crashes → rebalance triggered:
  C1: P0, P1, P4
  C2: P2, P3, P5
```

**Rebalance strategies:**

| Strategy | Behavior | Best For |
|---|---|---|
| `range` | Sequential ranges per consumer | Low churn groups |
| `roundrobin` | Alternating assignment | Balanced load |
| `sticky` | Minimize partition movement | Stateful consumers |
| `cooperative-sticky` | Incremental, no stop-the-world | Production (recommended) |

```python
consumer = Consumer({
    "bootstrap.servers":      "localhost:9092",
    "group.id":               "my-group",
    "partition.assignment.strategy": "cooperative-sticky",
})
```

---

## Offset Management

### Auto Commit (Simple but Risky)

```python
consumer = Consumer({
    "bootstrap.servers":   "localhost:9092",
    "group.id":            "simple-group",
    "enable.auto.commit":  True,
    "auto.commit.interval.ms": 5000,  # commit every 5 seconds
    # Risk: can commit offsets before processing completes → message loss on crash
})
```

### Manual Commit (Recommended)

```python
consumer = Consumer({
    "bootstrap.servers":  "localhost:9092",
    "group.id":           "reliable-group",
    "enable.auto.commit": False,
    "auto.offset.reset":  "earliest",
})

consumer.subscribe(["orders"])

while True:
    msg = consumer.poll(1.0)
    if msg is None or msg.error(): continue

    try:
        process_order(msg.value())
        # Synchronous: guaranteed before next poll
        consumer.commit(asynchronous=False)
        # Async: higher throughput, small risk of duplicate on crash
        # consumer.commit(asynchronous=True)
    except Exception:
        # Don't commit — message will be redelivered
        logger.error("Processing failed, not committing")
```

### Commit Specific Offsets

```python
from confluent_kafka import TopicPartition

# Commit offset for a specific partition
consumer.commit(offsets=[
    TopicPartition("orders", partition=0, offset=msg.offset() + 1)
])
```

---

## Seeking and Replaying

```python
from confluent_kafka import TopicPartition
import datetime

# Seek to beginning
consumer.subscribe(["orders"])
consumer.poll(0)   # trigger assignment
partitions = consumer.assignment()
consumer.seek_to_beginning(*partitions)

# Seek to specific offset
consumer.seek(TopicPartition("orders", partition=0, offset=42))

# Seek to timestamp
target_ts = int(datetime.datetime(2024, 1, 1).timestamp() * 1000)
partitions = [TopicPartition("orders", p, target_ts) for p in range(6)]
result = consumer.offsets_for_times(partitions)
for tp in result:
    consumer.seek(tp)
```

---

## Handling Rebalance Events

```python
from confluent_kafka import Consumer

def on_assign(consumer, partitions):
    print(f"Partitions assigned: {[p.partition for p in partitions]}")
    # Restore state or seek to custom position here

def on_revoke(consumer, partitions):
    print(f"Partitions revoked: {[p.partition for p in partitions]}")
    # Flush any pending state or offsets here
    consumer.commit(asynchronous=False)

def on_lost(consumer, partitions):
    print(f"Partitions lost (crash): {[p.partition for p in partitions]}")

consumer.subscribe(
    ["orders"],
    on_assign = on_assign,
    on_revoke = on_revoke,
    on_lost   = on_lost,
)
```

---

## Consumer Lag

**Consumer lag** = how many messages behind the latest offset a consumer is.

```bash
# CLI: describe group and see LAG column
kafka-consumer-groups --bootstrap-server localhost:9092 \
  --describe --group order-processor

# Output:
# GROUP           TOPIC   PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG
# order-processor orders  0          1234            1260            26
# order-processor orders  1          890             905             15
```

**High lag means:**
- Consumer is too slow (CPU-bound processing)
- Not enough consumer instances
- Rebalancing is happening frequently
- Network or broker issues

---

## Consumer Configuration Reference

| Config | Default | Description |
|---|---|---|
| `auto.offset.reset` | `latest` | `earliest` or `latest` when no committed offset |
| `enable.auto.commit` | `true` | Auto-commit every interval |
| `auto.commit.interval.ms` | `5000` | How often to auto-commit |
| `max.poll.interval.ms` | `300000` | Max time between polls before kicked out of group |
| `session.timeout.ms` | `45000` | Heartbeat timeout before consumer is considered dead |
| `heartbeat.interval.ms` | `3000` | How often to send heartbeats |
| `max.poll.records` | `500` | Max records per poll() call |
| `fetch.min.bytes` | `1` | Min bytes before broker responds |
| `fetch.max.wait.ms` | `500` | Max wait if fetch.min.bytes not met |

---

## At-Least-Once vs At-Most-Once vs Exactly-Once

| Semantics | How | Risk |
|---|---|---|
| **At-most-once** | Commit before processing | Message loss on crash |
| **At-least-once** | Commit after processing | Duplicate processing |
| **Exactly-once** | Transactions + idempotent writes | Complex, see Chapter 17 |

For most applications, **at-least-once + idempotent processing** is the right choice:

```python
def process_order(order_id: str, data: dict):
    # Make your processing idempotent
    if db.exists(order_id):
        return   # already processed, skip
    db.upsert(order_id, data)
```

---

## Batch Processing Pattern

```python
consumer = Consumer({
    "bootstrap.servers":  "localhost:9092",
    "group.id":           "batch-processor",
    "enable.auto.commit": False,
    "max.poll.records":   1000,
})

consumer.subscribe(["orders"])

while True:
    # Collect a batch
    batch = []
    while len(batch) < 1000:
        msg = consumer.poll(timeout=0.1)
        if msg is None or msg.error(): break
        batch.append(msg)

    if batch:
        # Bulk process the whole batch
        process_batch([msg.value() for msg in batch])
        # Commit the highest offset per partition
        consumer.commit(asynchronous=False)
```


---

# Chapter 6: Topics & Partitions

## Topic Configuration

Topics have two kinds of settings:
- **Static** (set at creation, require recreation to change): `replication.factor`, `num.partitions`
- **Dynamic** (alterable at runtime): retention, compression, cleanup policy

### Create with Custom Config

```bash
kafka-topics --bootstrap-server localhost:9092 \
  --create \
  --topic user-events \
  --partitions 12 \
  --replication-factor 3 \
  --config retention.ms=604800000 \
  --config compression.type=lz4 \
  --config cleanup.policy=delete
```

### Alter Topic Config

```bash
# Change retention to 30 days
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics \
  --entity-name user-events \
  --alter \
  --add-config retention.ms=2592000000

# Describe current config
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics \
  --entity-name user-events \
  --describe
```

---

## Key Topic Configurations

| Config | Default | Description |
|---|---|---|
| `retention.ms` | `604800000` (7 days) | How long to keep messages |
| `retention.bytes` | `-1` (unlimited) | Max size per partition |
| `cleanup.policy` | `delete` | `delete` or `compact` or `delete,compact` |
| `compression.type` | `producer` | `none`, `gzip`, `snappy`, `lz4`, `zstd` |
| `replication.factor` | Broker default | Copies of each partition |
| `min.insync.replicas` | `1` | Min ISR for producer `acks=all` |
| `segment.ms` | `604800000` (7 days) | When to roll a new segment |
| `segment.bytes` | `1073741824` (1 GB) | Max segment file size |
| `max.message.bytes` | `1048588` (~1 MB) | Max message size |
| `delete.retention.ms` | `86400000` (1 day) | How long to keep tombstones (compact) |

---

## Retention Policies

### Time-Based Retention (delete)

```bash
# Keep messages for 3 days
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics --entity-name logs \
  --alter --add-config "retention.ms=259200000"
```

### Size-Based Retention

```bash
# Keep at most 10 GB per partition
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics --entity-name logs \
  --alter --add-config "retention.bytes=10737418240"
```

### Combined (Whichever Triggers First)

```bash
--add-config "retention.ms=86400000,retention.bytes=5368709120"
```

---

## Log Compaction

Log compaction keeps only the **latest value per key**. Perfect for changelog/state topics.

```
Before compaction:
  [key=user1: v1][key=user2: v1][key=user1: v2][key=user1: v3]

After compaction:
  [key=user2: v1][key=user1: v3]

Tombstone: key=user1, value=null → delete key from compacted log
```

```bash
# Enable compaction
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics --entity-name user-profiles \
  --alter --add-config "cleanup.policy=compact"

# Compact + delete (compaction on recent data, delete old)
--add-config "cleanup.policy=compact,delete,retention.ms=86400000"
```

**Compaction settings:**

| Config | Default | Description |
|---|---|---|
| `min.cleanable.dirty.ratio` | `0.5` | Compact when dirty/total ratio > threshold |
| `min.compaction.lag.ms` | `0` | Don't compact messages newer than this |
| `max.compaction.lag.ms` | `Long.MAX_VALUE` | Force compaction after this time |
| `delete.retention.ms` | `86400000` | Keep tombstones this long |

---

## Partition Count: How to Choose

```
Rule of thumb:
  Partitions = max(throughput_MB_s / 10, desired_consumer_parallelism)

Example:
  Target throughput: 100 MB/s
  Max consumer instances: 20
  → partitions = max(10, 20) = 20 partitions
```

**More partitions = more parallelism but:**
- More file handles on brokers
- Longer leader election time on failure
- More memory in producer and consumer

```bash
# Increase partitions (can only go up, never down)
kafka-topics --bootstrap-server localhost:9092 \
  --alter --topic orders --partitions 24
```

---

## Replication Factor

```
replication.factor=1: No redundancy. Single point of failure.
replication.factor=2: One copy. Can survive 1 broker failure.
replication.factor=3: Two copies. Can survive 2 broker failures.
                       Standard for production.
```

**Relationship with min.insync.replicas:**

```
replication.factor=3, min.insync.replicas=2, acks=all:
  - 3 replicas, need 2 to acknowledge each write
  - Can survive 1 broker being down
  - If 2 brokers fail → write unavailable (but no data loss)

replication.factor=3, min.insync.replicas=1, acks=all:
  - Only 1 replica needed → leader only
  - Same as acks=1 effectively
  - Don't do this
```

---

## Partition Assignment: Leaders and Followers

```bash
# View partition leadership
kafka-topics --bootstrap-server localhost:9092 \
  --describe --topic orders

# Output:
# Topic: orders   PartitionCount: 6  ReplicationFactor: 3
# Topic: orders   Partition: 0  Leader: 1  Replicas: 1,2,3  Isr: 1,2,3
# Topic: orders   Partition: 1  Leader: 2  Replicas: 2,3,1  Isr: 2,3,1
# Topic: orders   Partition: 2  Leader: 3  Replicas: 3,1,2  Isr: 3,1,2
```

### Rebalance Preferred Leaders

After a broker restart, it may not be the preferred leader:

```bash
# Reassign partitions back to preferred leaders
kafka-leader-election --bootstrap-server localhost:9092 \
  --election-type preferred \
  --all-topic-partitions
```

---

## Topic Best Practices

```
1. Naming conventions:
   <domain>.<entity>.<event-type>
   Example: payments.orders.created
            analytics.pageviews.raw

2. Partition count:
   Start with 6 for moderate-traffic topics
   24+ for high-throughput topics
   Use multiples of consumer instance count

3. Replication:
   Always 3 in production
   min.insync.replicas = replication.factor - 1

4. Retention:
   Event topics: 7 days (replay window)
   Log aggregation: 3 days or size-limited
   Compaction topics: forever (with tombstones)

5. Message size:
   Target < 1MB per message
   For large payloads: store in S3, send reference in Kafka

6. Schema:
   Always use Schema Registry for value schema
   Key should be a simple string or int (user_id, order_id)
```

---

## Monitoring Topic Health

```bash
# Under-replicated partitions (ISR < replication factor = problem!)
kafka-topics --bootstrap-server localhost:9092 \
  --describe --under-replicated-partitions

# Unavailable partitions (leader = none = outage!)
kafka-topics --bootstrap-server localhost:9092 \
  --describe --unavailable-partitions

# Topics with no leader
kafka-topics --bootstrap-server localhost:9092 \
  --describe --under-min-isr-partitions
```


---

# Chapter 7: Kafka Streams

## What is Kafka Streams?

Kafka Streams is a **Java/Kotlin library** for building real-time stream processing applications on top of Kafka. Unlike Spark or Flink, it runs inside your application process — no separate cluster needed.

```
Input Topic ──► Kafka Streams App ──► Output Topic
                  (your process)
```

**Key properties:**
- Library (not a framework) — embedded in your app
- No external cluster — just Kafka
- Fault-tolerant with local RocksDB state stores
- Exactly-once semantics support
- Scales by adding instances

---

## Core Abstractions

| Abstraction | Description |
|---|---|
| **KStream** | Unbounded stream of records (insert semantics) |
| **KTable** | Changelog stream — latest value per key (update semantics) |
| **GlobalKTable** | Fully replicated KTable on every instance |
| **StreamsBuilder** | Entry point to build the topology |
| **Topology** | The processing graph (DAG of operations) |
| **State Store** | Local RocksDB for aggregations and joins |

---

## Basic Kafka Streams App (Java)

```java
import org.apache.kafka.streams.*;
import org.apache.kafka.streams.kstream.*;
import java.util.Properties;

public class WordCount {
    public static void main(String[] args) throws InterruptedException {
        Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG,    "wordcount-app");
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG,   Serdes.String().getClass());
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, Serdes.String().getClass());

        StreamsBuilder builder = new StreamsBuilder();

        KStream<String, String> textLines = builder.stream("text-input");

        KTable<String, Long> wordCounts = textLines
            .flatMapValues(line -> Arrays.asList(line.toLowerCase().split("\\W+")))
            .groupBy((key, word) -> word)
            .count();

        wordCounts.toStream().to("wordcount-output",
            Produced.with(Serdes.String(), Serdes.Long()));

        KafkaStreams streams = new KafkaStreams(builder.build(), props);
        streams.start();

        Runtime.getRuntime().addShutdownHook(new Thread(streams::close));
    }
}
```

---

## KStream Operations

```java
StreamsBuilder builder = new StreamsBuilder();
KStream<String, Order> orders = builder.stream("orders",
    Consumed.with(Serdes.String(), orderSerde));

// Filter
KStream<String, Order> highValue = orders
    .filter((key, order) -> order.getAmount() > 100.0);

// Map values
KStream<String, OrderEnriched> enriched = orders
    .mapValues(order -> enrich(order));

// FlatMap (one-to-many)
KStream<String, LineItem> lineItems = orders
    .flatMapValues(order -> order.getLineItems());

// Branch (split stream)
Map<String, KStream<String, Order>> branches = orders
    .split(Named.as("branch-"))
    .branch((k, v) -> v.getAmount() > 1000, Branched.as("premium"))
    .branch((k, v) -> v.getAmount() > 100,  Branched.as("standard"))
    .defaultBranch(Branched.as("small"));

KStream<String, Order> premium  = branches.get("branch-premium");
KStream<String, Order> standard = branches.get("branch-standard");

// Repartition (change key)
KStream<String, Order> byProduct = orders
    .selectKey((k, order) -> order.getProductId());

// Write to topic
highValue.to("high-value-orders", Produced.with(Serdes.String(), orderSerde));
```

---

## Aggregations

```java
// Count events per user per hour
KTable<Windowed<String>, Long> clicksPerUser = orders
    .groupByKey()
    .windowedBy(TimeWindows.ofSizeWithNoGrace(Duration.ofHours(1)))
    .count(Materialized.as("clicks-per-user-store"));

clicksPerUser
    .toStream()
    .map((windowedKey, count) -> KeyValue.pair(
        windowedKey.key(),
        windowedKey.window().start() + ":" + count
    ))
    .to("clicks-per-user-output");

// Sum revenue per category
KTable<String, Double> revenueByCategory = orders
    .groupBy((key, order) -> KeyValue.pair(order.getCategory(), order))
    .aggregate(
        () -> 0.0,
        (key, order, total) -> total + order.getAmount(),
        Materialized.<String, Double, KeyValueStore<Bytes, byte[]>>as("revenue-store")
            .withValueSerde(Serdes.Double())
    );
```

---

## Windowing

```java
// Tumbling window (non-overlapping, fixed size)
TimeWindows tumbling = TimeWindows.ofSizeWithNoGrace(Duration.ofMinutes(5));

// Sliding window (overlapping by step)
SlidingWindows sliding = SlidingWindows.ofTimeDifferenceAndGrace(
    Duration.ofMinutes(10),
    Duration.ofMinutes(1)
);

// Session window (gap-based)
SessionWindows session = SessionWindows.ofInactivityGapWithNoGrace(Duration.ofMinutes(30));

// Usage
KTable<Windowed<String>, Long> result = stream
    .groupByKey()
    .windowedBy(tumbling)
    .count(Materialized.as("windowed-count-store"));
```

---

## KStream-KTable Join

```java
KStream<String, Order>   orders   = builder.stream("orders");
KTable<String, Customer> customers = builder.table("customers");

// Enrich stream with table lookup
KStream<String, EnrichedOrder> enriched = orders
    .join(customers,
        (order, customer) -> new EnrichedOrder(order, customer),
        Joined.with(Serdes.String(), orderSerde, customerSerde)
    );

enriched.to("enriched-orders");
```

---

## KStream-KStream Join

```java
KStream<String, Order>   orders   = builder.stream("orders");
KStream<String, Payment> payments = builder.stream("payments");

// Join two streams within a 5-minute window
KStream<String, OrderWithPayment> combined = orders
    .join(payments,
        (order, payment) -> new OrderWithPayment(order, payment),
        JoinWindows.ofTimeDifferenceWithNoGrace(Duration.ofMinutes(5)),
        StreamJoined.with(Serdes.String(), orderSerde, paymentSerde)
    );
```

---

## Stateful Processing with State Stores

```java
// Custom state store for deduplication
StoreBuilder<KeyValueStore<String, Boolean>> storeBuilder =
    Stores.keyValueStoreBuilder(
        Stores.persistentKeyValueStore("dedup-store"),
        Serdes.String(),
        Serdes.Boolean()
    );

builder.addStateStore(storeBuilder);

KStream<String, Order> deduplicated = orders
    .transform(() -> new Transformer<String, Order, KeyValue<String, Order>>() {
        KeyValueStore<String, Boolean> store;

        @Override
        public void init(ProcessorContext context) {
            store = context.getStateStore("dedup-store");
        }

        @Override
        public KeyValue<String, Order> transform(String key, Order value) {
            if (store.get(value.getOrderId()) != null) {
                return null;  // duplicate, drop
            }
            store.put(value.getOrderId(), true);
            return KeyValue.pair(key, value);
        }
    }, "dedup-store");
```

---

## Interactive Queries

Query local state stores from outside the app:

```java
KafkaStreams streams = new KafkaStreams(builder.build(), props);
streams.start();

// Query local state store directly
ReadOnlyKeyValueStore<String, Long> store =
    streams.store(StoreQueryParameters.fromNameAndType(
        "wordcount-store",
        QueryableStoreTypes.keyValueStore()
    ));

Long count = store.get("kafka");
```

---

## Kafka Streams Configuration

```java
Properties props = new Properties();
props.put(StreamsConfig.APPLICATION_ID_CONFIG,    "my-streams-app");
props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");

// Exactly-once semantics
props.put(StreamsConfig.PROCESSING_GUARANTEE_CONFIG, StreamsConfig.EXACTLY_ONCE_V2);

// Parallelism
props.put(StreamsConfig.NUM_STREAM_THREADS_CONFIG, 4);

// State directory (RocksDB)
props.put(StreamsConfig.STATE_DIR_CONFIG, "/tmp/kafka-streams");

// Commit interval
props.put(StreamsConfig.COMMIT_INTERVAL_MS_CONFIG, 100);

// Replication factor for internal topics
props.put(StreamsConfig.REPLICATION_FACTOR_CONFIG, 3);
```

---

## Kafka Streams vs Other Frameworks

| Feature | Kafka Streams | Apache Flink | Apache Spark Streaming |
|---|---|---|---|
| Deployment | Embedded library | Separate cluster | Separate cluster |
| Latency | Low (ms) | Very low (ms) | Moderate (seconds) |
| State backend | RocksDB | RocksDB / heap | In-memory / HDFS |
| Exactly-once | Yes | Yes | Yes |
| Windowing | Yes | Rich | Yes |
| Language | Java / Kotlin | Java / Scala / Python | Python / Java / Scala |
| SQL interface | No | Flink SQL | Spark SQL |
| Learning curve | Moderate | Steep | Moderate |


---

# Chapter 8: Kafka Connect

## What is Kafka Connect?

Kafka Connect is a **framework for streaming data between Kafka and external systems** using reusable, configurable connectors. You declare what to connect — Kafka handles the rest.

```
External System ──► Source Connector ──► Kafka ──► Sink Connector ──► External System

Examples:
  PostgreSQL CDC ──► Debezium ──► Kafka ──► Elasticsearch Sink ──► Elasticsearch
  S3 Files       ──► S3 Source ──► Kafka ──► JDBC Sink           ──► MySQL
```

**Key properties:**
- Configuration-driven (JSON/REST API)
- Distributed mode: scales with multiple workers
- Exactly-once or at-least-once delivery
- Automatic offset management
- Transformations (SMTs) without custom code
- 200+ connectors in the Confluent Hub ecosystem

---

## Standalone vs Distributed Mode

| Feature | Standalone | Distributed |
|---|---|---|
| Workers | Single process | Multiple processes |
| Fault tolerance | None | Yes (tasks redistributed) |
| Scalability | No | Yes |
| Use case | Dev / testing | Production |
| Config | Local file | REST API |

---

## Start Connect (Distributed)

```properties
# connect-distributed.properties
bootstrap.servers=localhost:9092
group.id=connect-cluster

key.converter=org.apache.kafka.connect.json.JsonConverter
value.converter=org.apache.kafka.connect.json.JsonConverter
key.converter.schemas.enable=false
value.converter.schemas.enable=false

# For Avro with Schema Registry:
# key.converter=io.confluent.connect.avro.AvroConverter
# key.converter.schema.registry.url=http://localhost:8081
# value.converter=io.confluent.connect.avro.AvroConverter
# value.converter.schema.registry.url=http://localhost:8081

offset.storage.topic=connect-offsets
config.storage.topic=connect-configs
status.storage.topic=connect-status

offset.flush.interval.ms=10000
```

```bash
# Start Connect worker
connect-distributed connect-distributed.properties
```

---

## REST API

```bash
# List connectors
curl http://localhost:8083/connectors

# Create connector
curl -X POST http://localhost:8083/connectors \
  -H "Content-Type: application/json" \
  -d @connector-config.json

# Connector status
curl http://localhost:8083/connectors/my-connector/status

# Pause / resume
curl -X PUT http://localhost:8083/connectors/my-connector/pause
curl -X PUT http://localhost:8083/connectors/my-connector/resume

# Delete
curl -X DELETE http://localhost:8083/connectors/my-connector

# Available plugins
curl http://localhost:8083/connector-plugins
```

---

## Source Connectors

### Debezium PostgreSQL CDC

Capture every INSERT/UPDATE/DELETE from Postgres as Kafka events:

```json
{
  "name": "postgres-cdc",
  "config": {
    "connector.class":            "io.debezium.connector.postgresql.PostgresConnector",
    "database.hostname":          "postgres",
    "database.port":              "5432",
    "database.user":              "replicator",
    "database.password":          "secret",
    "database.dbname":            "shop",
    "database.server.name":       "shop",
    "table.include.list":         "public.orders,public.users",
    "plugin.name":                "pgoutput",
    "topic.prefix":               "cdc",
    "snapshot.mode":              "initial"
  }
}
```

Events published to: `cdc.public.orders`, `cdc.public.users`

### File Source (Simple)

```json
{
  "name": "file-source",
  "config": {
    "connector.class": "org.apache.kafka.connect.file.FileStreamSourceConnector",
    "file":            "/var/log/app.log",
    "topic":           "app-logs"
  }
}
```

---

## Sink Connectors

### Elasticsearch Sink

```json
{
  "name": "elasticsearch-sink",
  "config": {
    "connector.class":               "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
    "tasks.max":                     "3",
    "topics":                        "orders,events",
    "connection.url":                "http://elasticsearch:9200",
    "type.name":                     "_doc",
    "key.ignore":                    "false",
    "schema.ignore":                 "true",
    "behavior.on.null.values":       "DELETE",
    "batch.size":                    "2000",
    "max.retries":                   "5",
    "retry.backoff.ms":              "100"
  }
}
```

### JDBC Sink (Write to any SQL DB)

```json
{
  "name": "postgres-sink",
  "config": {
    "connector.class":        "io.confluent.connect.jdbc.JdbcSinkConnector",
    "tasks.max":              "3",
    "topics":                 "processed-orders",
    "connection.url":         "jdbc:postgresql://postgres:5432/warehouse",
    "connection.user":        "writer",
    "connection.password":    "secret",
    "auto.create":            "true",
    "auto.evolve":            "true",
    "insert.mode":            "upsert",
    "pk.mode":                "record_key",
    "pk.fields":              "order_id"
  }
}
```

### S3 Sink

```json
{
  "name": "s3-sink",
  "config": {
    "connector.class":         "io.confluent.connect.s3.S3SinkConnector",
    "tasks.max":               "4",
    "topics":                  "user-events",
    "s3.region":               "us-east-1",
    "s3.bucket.name":          "my-kafka-data",
    "s3.part.size":            "67108864",
    "storage.class":           "io.confluent.connect.s3.storage.S3Storage",
    "format.class":            "io.confluent.connect.s3.format.parquet.ParquetFormat",
    "flush.size":              "10000",
    "rotate.interval.ms":      "60000",
    "locale":                  "en_US",
    "timezone":                "UTC",
    "timestamp.extractor":     "RecordField",
    "timestamp.field":         "created_at",
    "partitioner.class":       "io.confluent.connect.storage.partitioner.TimeBasedPartitioner",
    "path.format":             "'year'=YYYY/'month'=MM/'day'=dd/'hour'=HH",
    "partition.duration.ms":   "3600000"
  }
}
```

---

## Single Message Transforms (SMTs)

SMTs modify messages in-flight without custom code:

```json
{
  "transforms":                   "extractKey,addTimestamp,dropField",

  "transforms.extractKey.type":   "org.apache.kafka.connect.transforms.ExtractField$Key",
  "transforms.extractKey.field":  "order_id",

  "transforms.addTimestamp.type": "org.apache.kafka.connect.transforms.InsertField$Value",
  "transforms.addTimestamp.timestamp.field": "kafka_timestamp",

  "transforms.dropField.type":    "org.apache.kafka.connect.transforms.ReplaceField$Value",
  "transforms.dropField.blacklist": "internal_field,debug_info"
}
```

### Common SMTs

| SMT | Purpose |
|---|---|
| `ReplaceField` | Rename, include, or exclude fields |
| `ExtractField` | Pull a field out to become the key or value |
| `InsertField` | Add fields (offset, partition, timestamp, static) |
| `MaskField` | Replace field value with null or a constant |
| `Filter` | Drop messages that match a condition |
| `TimestampConverter` | Convert timestamp formats |
| `RegexRouter` | Route to different topics by regex |
| `Flatten` | Flatten nested structs into dot-notation |

---

## Connector Monitoring

```bash
# Check task status
curl http://localhost:8083/connectors/postgres-cdc/status | jq .

# Check offset lag (source)
curl http://localhost:8083/connectors/postgres-cdc/offsets

# Restart a failed task
curl -X POST http://localhost:8083/connectors/postgres-cdc/tasks/0/restart

# Worker metrics via JMX / REST
curl http://localhost:8083/ | jq .version
```

Key JMX metrics:
- `kafka.connect:type=connector-task-metrics` → task state
- `kafka.connect:type=source-task-metrics` → poll rate, offset commit rate
- `kafka.connect:type=sink-task-metrics` → put rate, offset commit rate


---

# Chapter 9: Schema Registry

## Why Schema Registry?

Without schema management, Kafka is just raw bytes. Problems arise when:
- A producer changes the message format, breaking consumers
- Fields are renamed or removed without coordination
- Types change (string → integer)
- Multiple teams produce to the same topic with different schemas

**Schema Registry solves this:**
1. Centralized schema storage with versioning
2. Compatibility enforcement (BACKWARD, FORWARD, FULL)
3. Automatic schema evolution with safety guarantees
4. Schema IDs embedded in messages — consumers always know the schema

---

## How It Works

```
Producer                  Schema Registry         Consumer
   │                            │                    │
   │── register schema ────────►│                    │
   │◄─ schema_id=5 ────────────┤                    │
   │                            │                    │
   │── [magic byte][schema_id=5][serialized payload] ──►│
   │                            │                    │
   │                            │◄─ fetch schema 5 ──┤
   │                            │─ return schema ───►│
   │                            │                    │
   │                            │                 deserialize
```

Every Confluent-format message starts with:
- 1 magic byte (`0x00`)
- 4-byte schema ID
- serialized payload (Avro, Protobuf, or JSON Schema)

---

## Avro Schema Example

### Define Schema

```json
// order.avsc
{
  "type": "record",
  "name": "Order",
  "namespace": "com.example.shop",
  "fields": [
    {"name": "order_id",   "type": "string"},
    {"name": "user_id",    "type": "string"},
    {"name": "amount",     "type": "double"},
    {"name": "currency",   "type": {"type": "enum", "name": "Currency",
                                   "symbols": ["USD", "EUR", "GBP"]}},
    {"name": "created_at", "type": "long",   "logicalType": "timestamp-millis"},
    {"name": "metadata",   "type": ["null", {"type": "map", "values": "string"}],
                           "default": null}
  ]
}
```

### Producer with Avro

```python
from confluent_kafka import Producer
from confluent_kafka.schema_registry import SchemaRegistryClient
from confluent_kafka.schema_registry.avro import AvroSerializer
from confluent_kafka.serialization import SerializationContext, MessageField

schema_registry = SchemaRegistryClient({"url": "http://localhost:8081"})

order_schema = open("order.avsc").read()

avro_serializer = AvroSerializer(schema_registry, order_schema)

producer = Producer({"bootstrap.servers": "localhost:9092"})

order = {
    "order_id":   "ORD-123",
    "user_id":    "user-456",
    "amount":     99.99,
    "currency":   "USD",
    "created_at": 1704067200000,
    "metadata":   None,
}

producer.produce(
    topic = "orders",
    key   = order["order_id"],
    value = avro_serializer(order, SerializationContext("orders", MessageField.VALUE)),
)
producer.flush()
```

### Consumer with Avro

```python
from confluent_kafka import Consumer
from confluent_kafka.schema_registry import SchemaRegistryClient
from confluent_kafka.schema_registry.avro import AvroDeserializer
from confluent_kafka.serialization import SerializationContext, MessageField

schema_registry = SchemaRegistryClient({"url": "http://localhost:8081"})
avro_deserializer = AvroDeserializer(schema_registry)

consumer = Consumer({
    "bootstrap.servers": "localhost:9092",
    "group.id":          "order-consumer",
    "auto.offset.reset": "earliest",
})
consumer.subscribe(["orders"])

while True:
    msg = consumer.poll(1.0)
    if msg is None or msg.error(): continue

    order = avro_deserializer(msg.value(),
                              SerializationContext(msg.topic(), MessageField.VALUE))
    print(f"Order: {order}")
```

---

## Protobuf Schema Example

```protobuf
// order.proto
syntax = "proto3";
package com.example.shop;

message Order {
  string order_id   = 1;
  string user_id    = 2;
  double amount     = 3;
  string currency   = 4;
  int64  created_at = 5;
  map<string, string> metadata = 6;
}
```

```python
from confluent_kafka.schema_registry.protobuf import ProtobufSerializer, ProtobufDeserializer
from order_pb2 import Order  # generated from proto file

protobuf_serializer = ProtobufSerializer(Order, schema_registry)

order = Order(order_id="ORD-123", user_id="user-456", amount=99.99, currency="USD")
producer.produce(
    "orders",
    value = protobuf_serializer(order, SerializationContext("orders", MessageField.VALUE))
)
```

---

## JSON Schema

```python
json_schema_str = """
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title":   "Order",
  "type":    "object",
  "properties": {
    "order_id": {"type": "string"},
    "amount":   {"type": "number", "minimum": 0}
  },
  "required": ["order_id", "amount"]
}
"""

from confluent_kafka.schema_registry.json_schema import JSONSerializer
json_serializer = JSONSerializer(json_schema_str, schema_registry)
```

---

## Compatibility Modes

| Mode | What It Allows |
|---|---|
| `BACKWARD` | New schema can read data written with old schema (consumers upgrade first) |
| `BACKWARD_TRANSITIVE` | New schema can read all previous versions |
| `FORWARD` | Old schema can read data written with new schema (producers upgrade first) |
| `FORWARD_TRANSITIVE` | Any old schema can read data written with new schema |
| `FULL` | Both BACKWARD and FORWARD |
| `FULL_TRANSITIVE` | Both, transitively |
| `NONE` | No compatibility checks |

**Rule of thumb: use `FULL_TRANSITIVE` for production.**

### Safe Schema Changes (FULL compatibility)

| Change | Safe? |
|---|---|
| Add optional field with default | Yes |
| Remove optional field | Yes |
| Change field type | No |
| Remove required field | No |
| Rename field | No (use aliases) |
| Add required field without default | No |

---

## Schema Registry REST API

```bash
# Register a schema
curl -X POST http://localhost:8081/subjects/orders-value/versions \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  -d '{"schema": "{\"type\":\"record\",\"name\":\"Order\",...}"}'

# List subjects
curl http://localhost:8081/subjects

# Get all versions
curl http://localhost:8081/subjects/orders-value/versions

# Get specific version
curl http://localhost:8081/subjects/orders-value/versions/1

# Get latest schema
curl http://localhost:8081/subjects/orders-value/versions/latest

# Check compatibility
curl -X POST http://localhost:8081/compatibility/subjects/orders-value/versions/latest \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  -d '{"schema": "..."}'

# Set compatibility
curl -X PUT http://localhost:8081/config/orders-value \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  -d '{"compatibility": "FULL_TRANSITIVE"}'

# Delete a subject (careful!)
curl -X DELETE http://localhost:8081/subjects/orders-value
```

---

## Schema Evolution Example

```
Version 1: {order_id, user_id, amount}
Version 2: {order_id, user_id, amount, currency="USD"}  ← add with default: SAFE
Version 3: {order_id, user_id, amount, currency="USD", metadata=null} ← SAFE

NOT SAFE:
  {order_id, user_id, total}   ← renamed "amount" to "total"
  {order_id, user_id}          ← removed required field
```

---

## Using Confluent Cloud Schema Registry

```python
schema_registry = SchemaRegistryClient({
    "url":              "https://your-sr.confluent.cloud",
    "basic.auth.user.info": "API_KEY:API_SECRET",
})
```


---

# Chapter 10: Cluster Architecture

## Cluster Components

```
┌──────────────────────────────────────────────────────────────────┐
│                        Kafka Cluster                             │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │                Controller (KRaft / ZooKeeper)              │  │
│  │  - Manages partition leadership                            │  │
│  │  - Maintains cluster metadata                             │  │
│  │  - Handles broker failures                                │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌───────────┐        ┌───────────┐        ┌───────────┐        │
│  │  Broker 1 │        │  Broker 2 │        │  Broker 3 │        │
│  │           │        │           │        │           │        │
│  │ Leader:   │        │ Leader:   │        │ Leader:   │        │
│  │ P0, P3    │        │ P1, P4    │        │ P2, P5    │        │
│  │           │        │           │        │           │        │
│  │ Replica:  │        │ Replica:  │        │ Replica:  │        │
│  │ P1, P4    │        │ P2, P5    │        │ P0, P3    │        │
│  └───────────┘        └───────────┘        └───────────┘        │
└──────────────────────────────────────────────────────────────────┘
          ▲                    ▲                    ▲
   Producers/Consumers connect to leader broker for each partition
```

---

## ZooKeeper Mode (Legacy)

ZooKeeper manages:
- Broker registry (who is alive)
- Controller election
- Topic configuration
- ACLs

```
ZooKeeper Ensemble (3 or 5 nodes for quorum)
  zk1 ◄──► zk2 ◄──► zk3

Kafka Brokers connect to ZooKeeper at startup:
  broker.id=1
  zookeeper.connect=zk1:2181,zk2:2181,zk3:2181
```

**Limitations of ZooKeeper mode:**
- Hard limit ~200,000 partitions per cluster
- Long recovery time after controller election
- Extra operational component to manage
- Being removed in Kafka 4.0+

---

## KRaft Mode (Kafka 3.3+)

KRaft uses Kafka's own Raft consensus protocol to replace ZooKeeper:

```
KRaft Controllers (3 nodes, Raft quorum)
  controller-1 ◄──► controller-2 ◄──► controller-3

Combined mode (broker + controller, OK for small clusters):
  node.id=1
  process.roles=broker,controller
  controller.quorum.voters=1@host1:9093,2@host2:9093,3@host3:9093

Dedicated mode (production):
  # 3 controller-only nodes + N broker-only nodes
  # controller: process.roles=controller
  # broker:     process.roles=broker
```

---

## Leader Election

Every partition has exactly one **leader** and zero or more **followers**:

```
Topic: orders, Partition 0
  Leader: Broker 1    ← all reads and writes go here
  Follower: Broker 2  ← replicates from leader (ISR)
  Follower: Broker 3  ← replicates from leader (ISR)

Broker 1 fails:
  1. Controller detects via heartbeat timeout
  2. Controller picks new leader from ISR (e.g., Broker 2)
  3. Updates cluster metadata
  4. Producers/consumers reconnect to new leader
  5. Broker 3 now replicates from Broker 2
```

---

## ISR (In-Sync Replicas)

The ISR is the set of replicas that are fully caught up with the leader:

```
Scenario: replication.factor=3, min.insync.replicas=2

ISR = [Broker1, Broker2, Broker3]   → writes succeed (acks=all)
ISR = [Broker1, Broker2]            → writes succeed (ISR ≥ 2)
ISR = [Broker1]                     → writes fail (ISR < 2)
                                       NotEnoughReplicasException
```

A follower is **removed from ISR** when it falls behind by:
```
replica.lag.time.max.ms=30000  (default: 30 seconds)
```

---

## Cluster Metadata

All brokers cache cluster metadata:
- List of all brokers
- Topics, partition counts, replication factors
- Leader for each partition
- Consumer group state

Clients refresh metadata periodically:
```
metadata.max.age.ms=300000   (default: 5 minutes)
```

---

## Multi-Rack / Multi-Zone Awareness

Distribute partition replicas across availability zones:

```properties
# broker.rack in server.properties for each broker
# AZ-1 brokers:
broker.id=1
broker.rack=us-east-1a

# AZ-2 brokers:
broker.id=2
broker.rack=us-east-1b

# AZ-3 brokers:
broker.id=3
broker.rack=us-east-1c
```

Kafka then ensures replicas spread across racks:
```
Partition 0: Leader (az1), Follower (az2), Follower (az3)
If az1 fails → az2 or az3 takes over immediately
```

---

## Broker Configuration Reference

```properties
# Identity
broker.id=1
node.id=1                           # KRaft

# Network
listeners=PLAINTEXT://0.0.0.0:9092
advertised.listeners=PLAINTEXT://broker1.example.com:9092
num.network.threads=3
socket.send.buffer.bytes=102400
socket.receive.buffer.bytes=102400
socket.request.max.bytes=104857600

# IO
num.io.threads=8
num.recovery.threads.per.data.dir=1
log.dirs=/var/kafka/data1,/var/kafka/data2  # multiple disks

# Replication
default.replication.factor=3
min.insync.replicas=2
offsets.topic.replication.factor=3
transaction.state.log.replication.factor=3
transaction.state.log.min.isr=2

# Retention
log.retention.hours=168
log.segment.bytes=1073741824
log.retention.check.interval.ms=300000

# Performance
num.partitions=6
message.max.bytes=1048576
replica.fetch.max.bytes=1048576
replica.lag.time.max.ms=30000
```

---

## Cluster Sizing

```
Throughput Formula:
  Brokers needed = ceil(
    (producer_MB_s × replication_factor) / per_broker_throughput_MB_s
  )

Example:
  Producer throughput: 100 MB/s
  Replication factor: 3
  Per-broker disk I/O limit: 150 MB/s
  → Brokers = ceil(100 × 3 / 150) = ceil(2) = 3 brokers minimum

Storage Formula:
  Storage per broker = (daily_ingest_GB × retention_days × replication_factor) / brokers

Partition Formula:
  Partitions = max(
    throughput_MB_s / 10,           ← producer throughput estimate
    expected_consumer_parallelism   ← max consumer instances
  )
```

---

## Internal Topics

| Topic | Purpose |
|---|---|
| `__consumer_offsets` | Consumer group offset storage |
| `__transaction_state` | Exactly-once transaction state |
| `_schemas` | Schema Registry schema storage |
| `connect-offsets` | Kafka Connect source offsets |
| `connect-configs` | Kafka Connect connector configs |
| `connect-status` | Kafka Connect task status |

These are created automatically. Don't delete them.


---

# Chapter 11: Replication & Durability

## The Replication Protocol

Kafka replication is **leader-based**: all reads and writes go through the partition leader. Followers pull data from the leader in parallel.

```
Write path:
  Producer ──► Broker 1 (Leader, P0)
                    │
                    ├──► Broker 2 (Follower, P0) ← replicate
                    └──► Broker 3 (Follower, P0) ← replicate
                    │
                    ACK to producer (when ISR confirms)

Read path:
  Consumer ──► Broker 1 (Leader, P0)
                    │
                    Return records
                    (followers also serve reads in Kafka 2.4+ with replica.selector.class)
```

---

## Replication Factor

```bash
# Create topic with RF=3
kafka-topics --bootstrap-server localhost:9092 \
  --create --topic payments \
  --partitions 6 \
  --replication-factor 3

# Describe to see replicas
kafka-topics --bootstrap-server localhost:9092 --describe --topic payments

# Output shows Replicas and Isr (In-Sync Replicas) per partition:
# Partition: 0  Leader: 2  Replicas: 2,1,3  Isr: 2,1,3
# Partition: 1  Leader: 1  Replicas: 1,3,2  Isr: 1,3,2
```

---

## Producer Acks and Durability

The `acks` setting controls the durability/latency tradeoff:

```
acks=0 (fire and forget)
───────────────────────
Producer ──(write)──► Broker
          ←(no ack)─
Risk: Network or broker failure = lost message
Use: Metrics, non-critical logging

acks=1 (leader only)
────────────────────
Producer ──(write)──► Leader ──(ack)──► Producer
                        │
                        Followers might not have it yet
Risk: Leader crashes before replication = lost message
Use: Moderate throughput, acceptable loss

acks=all (full ISR)
────────────────────
Producer ──(write)──► Leader ──(replicate)──► All ISR followers
                        │◄──────────────────────
                        ACK after all ISR confirm
Risk: None (as long as min.insync.replicas respected)
Use: Production, financial data, audit logs
```

---

## min.insync.replicas

This setting determines the minimum number of replicas that must acknowledge a write:

```
Configuration:
  replication.factor=3
  min.insync.replicas=2
  acks=all

Scenario 1: ISR = [B1, B2, B3] → writes succeed
Scenario 2: ISR = [B1, B2]     → writes succeed (ISR ≥ min.insync)
Scenario 3: ISR = [B1]         → writes FAIL
  → Producer gets: NotEnoughReplicasException
  → Cluster is write-unavailable for this partition

Why this is correct:
  If ISR=[B1] and we allowed writes, data would be on only 1 broker.
  If B1 then crashes → DATA LOSS despite acks=all.
  Better to reject writes than accept unacknowledged data.
```

**Recommended production settings:**
```properties
# server.properties (applied cluster-wide as defaults)
default.replication.factor=3
min.insync.replicas=2

# Per-topic override
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics --entity-name payments \
  --alter --add-config min.insync.replicas=2
```

---

## Follower Replication

Followers continuously fetch batches from the leader:

```properties
# How long a follower can lag before removed from ISR
replica.lag.time.max.ms=30000       # default: 30 seconds

# Max bytes per fetch request
replica.fetch.max.bytes=1048576     # 1 MB

# How fast followers replicate
replica.fetch.min.bytes=1
replica.fetch.wait.max.ms=500
```

---

## Unclean Leader Election

When all ISR replicas fail simultaneously and only out-of-sync replicas remain:

```properties
# Allow out-of-sync replica to become leader (risk: data loss)
unclean.leader.election.enable=false   # default: false for most topics

# Use true only for high-availability over data durability
# e.g., real-time metrics where losing some data is acceptable
```

---

## Log Segments and Durability

Kafka stores data in segment files:

```
/kafka-logs/orders-0/
  00000000000000000000.log         ← active segment (writable)
  00000000000000012345.log         ← closed segment
  00000000000000012345.index       ← offset index
  00000000000000012345.timeindex   ← timestamp index
```

**Log flushing (fsync) settings:**

```properties
# Flush data to disk after N messages (default: Long.MAX_VALUE = rely on OS)
log.flush.interval.messages=Long.MAX_VALUE

# Flush at least every N ms (default: Long.MAX_VALUE)
log.flush.interval.ms=Long.MAX_VALUE

# For strongest durability (but high latency):
log.flush.interval.messages=1      # fsync every message — VERY slow
```

**Recommendation:** Leave at defaults. Rely on replication (ISR) for durability, not fsync. OS page cache provides good protection.

---

## Replication Factor Change (Partition Reassignment)

You cannot change replication factor via `kafka-topics --alter`. Use `kafka-reassign-partitions`:

```bash
# 1. Create reassignment JSON
cat > reassignment.json << 'EOF'
{
  "version": 1,
  "partitions": [
    {"topic": "orders", "partition": 0, "replicas": [1,2,3], "log_dirs": ["any","any","any"]},
    {"topic": "orders", "partition": 1, "replicas": [2,3,1], "log_dirs": ["any","any","any"]}
  ]
}
EOF

# 2. Execute reassignment
kafka-reassign-partitions --bootstrap-server localhost:9092 \
  --reassignment-json-file reassignment.json \
  --execute \
  --throttle 50000000   # 50 MB/s throttle (avoid saturating network)

# 3. Verify progress
kafka-reassign-partitions --bootstrap-server localhost:9092 \
  --reassignment-json-file reassignment.json \
  --verify

# 4. Remove throttle after completion
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type brokers --entity-default \
  --alter --delete-config leader.replication.throttled.rate
```

---

## Durability Checklist

| Setting | Production Value | Why |
|---|---|---|
| `replication.factor` | 3 | Survive 2 simultaneous broker failures |
| `min.insync.replicas` | 2 | No data accepted with only 1 copy |
| `acks` (producer) | `all` | All ISR must confirm |
| `unclean.leader.election.enable` | `false` | No data loss on leader election |
| `enable.idempotence` (producer) | `true` | No duplicates on retry |
| `replica.lag.time.max.ms` | `30000` | Remove lagging replicas from ISR promptly |

---

## Monitoring Replication Health

```bash
# Under-replicated partitions — ISR smaller than replication factor
kafka-topics --bootstrap-server localhost:9092 \
  --describe --under-replicated-partitions

# Non-preferred leaders (performance degraded)
kafka-topics --bootstrap-server localhost:9092 \
  --describe

# JMX metrics
# kafka.server:type=ReplicaManager,name=UnderReplicatedPartitions → should be 0
# kafka.server:type=ReplicaManager,name=IsrShrinksPerSec → should be 0 normally
# kafka.server:type=ReplicaManager,name=IsrExpandsPerSec → normal during recoveries
```


---

# Chapter 12: Performance & Tuning

## Understanding Kafka Performance

Kafka achieves high throughput through:

1. **Sequential disk I/O** — append-only writes are fast (disk can do 500 MB/s sequential vs ~random IOPS)
2. **OS page cache** — Linux caches recently read disk pages in RAM
3. **Zero-copy** (`sendfile` syscall) — kernel sends data directly from page cache to network socket
4. **Batching** — producers accumulate messages, send in one network call
5. **Compression** — reduce network and disk bandwidth

```
Without zero-copy:
  Disk → Kernel page cache → User space → Kernel socket buffer → Network

With zero-copy (sendfile):
  Disk → Kernel page cache ──────────────────────────────► Network
  (2 copies vs 4 copies)
```

---

## Producer Tuning

### Maximize Throughput

```python
producer = Producer({
    "bootstrap.servers": "broker1:9092,broker2:9092,broker3:9092",

    # Batching — let messages accumulate
    "linger.ms":        20,         # wait 20ms to fill a batch
    "batch.size":       65536,      # 64KB batch size
    "buffer.memory":    67108864,   # 64MB total send buffer

    # Compression reduces network bandwidth significantly
    "compression.type": "lz4",      # lz4: fast, good ratio
                                    # zstd: best ratio, slower CPU
                                    # snappy: balanced

    # Throughput-optimized acks
    "acks":             "1",        # leader only (accept some risk)
})
```

### Maximize Reliability

```python
producer = Producer({
    "bootstrap.servers":         "localhost:9092",
    "acks":                      "all",
    "enable.idempotence":        True,
    "retries":                   10,
    "retry.backoff.ms":          100,
    "delivery.timeout.ms":       120000,
    "max.in.flight.requests.per.connection": 5,
    "linger.ms":                 5,      # still batch a bit
    "compression.type":         "snappy",
})
```

---

## Consumer Tuning

```python
consumer = Consumer({
    "bootstrap.servers":      "localhost:9092",
    "group.id":               "fast-consumer",
    "max.poll.records":       2000,     # fetch more records per poll
    "fetch.min.bytes":        1048576,  # 1MB — wait for enough data
    "fetch.max.wait.ms":      500,      # max wait for fetch.min.bytes
    "fetch.max.bytes":        52428800, # 50MB max per fetch
    "max.partition.fetch.bytes": 1048576,

    # Avoid rebalancing from slow consumers
    "max.poll.interval.ms":   600000,   # 10 min max processing time
    "session.timeout.ms":     30000,
    "heartbeat.interval.ms":  3000,
})
```

---

## Broker-Level Tuning

### OS Configuration

```bash
# Increase file descriptors (Kafka needs many open files)
echo "* soft nofile 100000" >> /etc/security/limits.conf
echo "* hard nofile 100000" >> /etc/security/limits.conf

# Increase network buffers
sysctl -w net.core.rmem_max=16777216
sysctl -w net.core.wmem_max=16777216
sysctl -w net.ipv4.tcp_rmem="4096 87380 16777216"
sysctl -w net.ipv4.tcp_wmem="4096 65536 16777216"

# Disable swapping for Kafka heap
sysctl -w vm.swappiness=1

# Dirty page flushing (don't flush too aggressively)
sysctl -w vm.dirty_ratio=80
sysctl -w vm.dirty_background_ratio=5

# Use deadline or noop I/O scheduler (not cfq)
echo deadline > /sys/block/sda/queue/scheduler
```

### JVM Tuning (kafka-server-start.sh)

```bash
export KAFKA_HEAP_OPTS="-Xmx6g -Xms6g"   # 6GB heap, pre-allocate
export KAFKA_JVM_PERFORMANCE_OPTS="
  -server
  -XX:+UseG1GC
  -XX:MaxGCPauseMillis=20
  -XX:InitiatingHeapOccupancyPercent=35
  -XX:+ExplicitGCInvokesConcurrent
  -Djava.awt.headless=true"
```

**Heap size guidance:**
- Don't give Kafka more than 6–8 GB heap
- Kafka relies heavily on OS page cache (RAM not in heap)
- More RAM → more page cache → better read performance

### Storage

```bash
# Use XFS or ext4 (XFS preferred for Kafka)
mkfs.xfs -f /dev/sdb
mount -t xfs -o noatime,nodiratime /dev/sdb /var/kafka/data

# RAID: use JBOD (Just a Bunch of Disks) not RAID5/6
# Multiple log.dirs = Kafka stripes across disks automatically
log.dirs=/kafka/disk1,/kafka/disk2,/kafka/disk3
```

---

## Compression

| Codec | CPU | Ratio | Latency | Best For |
|---|---|---|---|---|
| `none` | None | 1.0× | Lowest | Already-compressed data |
| `gzip` | High | 3-5× | High | Archival, max compression |
| `snappy` | Low | 2-3× | Low | General purpose |
| `lz4` | Low | 2-3× | Very low | High throughput |
| `zstd` | Medium | 3-6× | Low | Best ratio/speed tradeoff |

**Producer-side vs broker-side compression:**
- `compression.type=producer` (default) — broker stores as-is
- Set on topic to recompress: risky, usually avoid

---

## Partitions and Parallelism

```
Consumer throughput ≈ consumer_instances × per_consumer_MB_s

To process 100 MB/s:
  If each consumer handles 10 MB/s → need 10 consumers → need ≥10 partitions

Increasing partitions:
  - Increases producer parallelism
  - Increases consumer parallelism
  - Increases broker overhead (file handles, leader election time)

Don't over-partition:
  - 1,000 partitions × 3 replicas = 3,000 partition replicas per cluster
  - Kafka recommends max ~4,000 partitions per broker
```

---

## Benchmarking

### Producer Perf Test

```bash
kafka-producer-perf-test \
  --topic benchmark \
  --num-records 10000000 \
  --record-size 1000 \
  --throughput -1 \
  --producer-props \
    bootstrap.servers=localhost:9092 \
    acks=all \
    linger.ms=5 \
    batch.size=65536 \
    compression.type=lz4
```

### Consumer Perf Test

```bash
kafka-consumer-perf-test \
  --bootstrap-server localhost:9092 \
  --topic benchmark \
  --messages 10000000 \
  --group perf-test-group \
  --threads 4
```

---

## Performance Anti-Patterns

| Anti-Pattern | Problem | Fix |
|---|---|---|
| `linger.ms=0` + small messages | Tiny batches, high overhead | Set `linger.ms=5-20` |
| No compression | Wastes network + disk | Enable `lz4` or `snappy` |
| Too many small partitions | Leader election overhead | Consolidate topics |
| `auto.offset.reset=latest` during replay | Miss events | Use `earliest` for replay |
| Long `max.poll.interval.ms` | Slow rebalance detection | Tune to actual processing time |
| Heap > 8 GB | GC pauses | Cap heap, add RAM for page cache |
| Sync commit per message | Low throughput | Batch commits or async |
| RAID 5/6 | Write amplification | Use JBOD |


---

# Chapter 13: Security

## Security Layers

```
1. Encryption in Transit   ← TLS between clients and brokers
2. Authentication          ← Who is this client? (SASL/TLS)
3. Authorization           ← What can this client do? (ACLs)
4. Encryption at Rest      ← Disk encryption (OS/cloud level)
5. Audit Logging           ← Who did what and when
```

---

## TLS (Encryption in Transit)

### Generate Certificates

```bash
# 1. Create CA
openssl req -new -x509 -keyout ca.key -out ca.crt -days 365 \
  -subj "/CN=Kafka-CA/O=MyOrg/C=US" -passout pass:ca-secret

# 2. Create broker keystore
keytool -genkeypair -alias broker1 \
  -keyalg RSA -keysize 2048 \
  -keystore broker1.keystore.jks \
  -dname "CN=broker1.example.com,O=MyOrg,C=US" \
  -storepass broker-secret -keypass broker-secret

# 3. Create CSR
keytool -certreq -alias broker1 \
  -keystore broker1.keystore.jks \
  -file broker1.csr \
  -storepass broker-secret

# 4. Sign with CA
openssl x509 -req -CA ca.crt -CAkey ca.key -in broker1.csr \
  -out broker1.crt -days 365 -CAcreateserial -passin pass:ca-secret

# 5. Import CA cert + signed cert into keystore
keytool -importcert -alias ca-cert -file ca.crt \
  -keystore broker1.keystore.jks -storepass broker-secret -noprompt
keytool -importcert -alias broker1 -file broker1.crt \
  -keystore broker1.keystore.jks -storepass broker-secret

# 6. Create truststore (contains CA cert)
keytool -importcert -alias ca-cert -file ca.crt \
  -keystore kafka.truststore.jks -storepass truststore-secret -noprompt
```

### Broker TLS Config (server.properties)

```properties
listeners=PLAINTEXT://0.0.0.0:9092,SSL://0.0.0.0:9093
advertised.listeners=PLAINTEXT://broker1:9092,SSL://broker1:9093
listener.security.protocol.map=PLAINTEXT:PLAINTEXT,SSL:SSL

ssl.keystore.location=/etc/kafka/ssl/broker1.keystore.jks
ssl.keystore.password=broker-secret
ssl.key.password=broker-secret
ssl.truststore.location=/etc/kafka/ssl/kafka.truststore.jks
ssl.truststore.password=truststore-secret
ssl.client.auth=required          # require mutual TLS
ssl.protocol=TLSv1.3
ssl.enabled.protocols=TLSv1.3,TLSv1.2
```

### Client TLS Config

```python
producer = Producer({
    "bootstrap.servers":       "broker1:9093",
    "security.protocol":       "SSL",
    "ssl.ca.location":         "/etc/kafka/ssl/ca.crt",
    "ssl.certificate.location":"/etc/kafka/ssl/client.crt",
    "ssl.key.location":        "/etc/kafka/ssl/client.key",
    "ssl.key.password":        "client-secret",
})
```

---

## SASL Authentication

### SASL/PLAIN (Username/Password — Use Over TLS)

```properties
# server.properties
listeners=SASL_SSL://0.0.0.0:9093
security.inter.broker.protocol=SASL_SSL
sasl.mechanism.inter.broker.protocol=PLAIN
sasl.enabled.mechanisms=PLAIN

# JAAS config (kafka_server_jaas.conf)
KafkaServer {
  org.apache.kafka.common.security.plain.PlainLoginModule required
  username="admin"
  password="admin-secret"
  user_admin="admin-secret"
  user_producer="producer-secret"
  user_consumer="consumer-secret";
};
```

```bash
export KAFKA_OPTS="-Djava.security.auth.login.config=/etc/kafka/kafka_server_jaas.conf"
```

```python
# Client config
producer = Producer({
    "bootstrap.servers":             "broker1:9093",
    "security.protocol":             "SASL_SSL",
    "sasl.mechanism":                "PLAIN",
    "sasl.username":                 "producer",
    "sasl.password":                 "producer-secret",
    "ssl.ca.location":               "/etc/kafka/ssl/ca.crt",
})
```

### SASL/SCRAM (Stored in ZooKeeper/KRaft)

```bash
# Create user
kafka-configs --bootstrap-server localhost:9092 \
  --alter --entity-type users --entity-name alice \
  --add-config 'SCRAM-SHA-256=[iterations=8192,password=alice-secret]'

# List users
kafka-configs --bootstrap-server localhost:9092 \
  --describe --entity-type users
```

```properties
# server.properties
sasl.enabled.mechanisms=SCRAM-SHA-256
sasl.mechanism.inter.broker.protocol=SCRAM-SHA-256
```

### SASL/GSSAPI (Kerberos)

```properties
# server.properties
sasl.enabled.mechanisms=GSSAPI
sasl.kerberos.service.name=kafka

# JAAS
KafkaServer {
  com.sun.security.auth.module.Krb5LoginModule required
  useKeyTab=true
  storeKey=true
  keyTab="/etc/security/keytabs/kafka.keytab"
  principal="kafka/broker1.example.com@REALM.COM";
};
```

---

## ACLs (Authorization)

### Enable ACLs

```properties
# server.properties
authorizer.class.name=org.apache.kafka.metadata.authorizer.StandardAuthorizer
allow.everyone.if.no.acl.found=false   # deny by default
super.users=User:admin
```

### Manage ACLs

```bash
# Allow producer to write to topic
kafka-acls --bootstrap-server localhost:9092 \
  --add \
  --allow-principal User:producer-service \
  --operation Write \
  --operation Describe \
  --topic orders

# Allow consumer group to read from topic
kafka-acls --bootstrap-server localhost:9092 \
  --add \
  --allow-principal User:consumer-service \
  --operation Read \
  --operation Describe \
  --topic orders

# Allow consumer to use a consumer group
kafka-acls --bootstrap-server localhost:9092 \
  --add \
  --allow-principal User:consumer-service \
  --operation Read \
  --group order-processor

# Allow all operations on a topic (admin)
kafka-acls --bootstrap-server localhost:9092 \
  --add \
  --allow-principal User:admin \
  --operation All \
  --topic '*'

# List ACLs
kafka-acls --bootstrap-server localhost:9092 --list --topic orders

# Delete ACL
kafka-acls --bootstrap-server localhost:9092 \
  --remove \
  --allow-principal User:old-service \
  --operation Write \
  --topic orders
```

---

## ACL Operations Reference

| Resource Type | Operations |
|---|---|
| Topic | Read, Write, Create, Delete, Alter, Describe, DescribeConfigs, AlterConfigs |
| Group | Read, Delete, Describe |
| Cluster | Create, Alter, Describe, AlterConfigs, IdempotentWrite, ClusterAction |
| TransactionalId | Write, Describe |

---

## Encryption at Rest

Kafka does not encrypt data at rest natively. Use:

```bash
# Linux LUKS disk encryption
cryptsetup luksFormat /dev/sdb
cryptsetup luksOpen /dev/sdb kafka-data
mkfs.xfs /dev/mapper/kafka-data
mount /dev/mapper/kafka-data /var/kafka/data

# Cloud provider solutions:
# AWS: EBS encryption (KMS), S3 SSE for tiered storage
# GCP: GCE persistent disk encryption (default)
# Azure: Azure Disk Encryption
```

---

## Security Best Practices

```
1. Always use TLS in production (SASL_SSL, not SASL_PLAINTEXT)
2. Use SASL/SCRAM or GSSAPI (not PLAIN) for authentication
3. Enable ACLs with deny-by-default
4. Give each service its own credentials (principle of least privilege)
5. Rotate credentials periodically
6. Never use super.users in application code
7. Use mutual TLS for inter-broker communication
8. Enable audit logging for compliance
9. Separate networks: brokers on internal network, expose only via load balancer
10. Encrypt at rest using cloud provider disk encryption
```

---

## Audit Logging

```properties
# Enable log4j audit (add to log4j.properties)
log4j.logger.kafka.authorizer.logger=INFO, authorizerAppender
log4j.additivity.kafka.authorizer.logger=false
log4j.appender.authorizerAppender=org.apache.log4j.DailyRollingFileAppender
log4j.appender.authorizerAppender.DatePattern='.'yyyy-MM-dd-HH
log4j.appender.authorizerAppender.File=${kafka.logs.dir}/kafka-authorizer.log
log4j.appender.authorizerAppender.layout=org.apache.log4j.PatternLayout
log4j.appender.authorizerAppender.layout.ConversionPattern=[%d] %p %m (%c)%n
```

Audit log records `ALLOWED` and `DENIED` authorization events with principal, operation, and resource.


---

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


---

# Chapter 15: Kafka with Python

## Libraries

| Library | Type | Notes |
|---|---|---|
| `confluent-kafka` | Production-grade | C librdkafka bindings, fastest, recommended |
| `kafka-python` | Pure Python | Easier to install, lower throughput |
| `faust` | Stream processing | Kafka Streams-like in Python |
| `aiokafka` | Async/await | asyncio-based consumer/producer |

```bash
pip install confluent-kafka
pip install kafka-python
pip install faust
```

---

## confluent-kafka: Production Producer

```python
import json
import logging
from typing import Optional
from confluent_kafka import Producer, KafkaException
from confluent_kafka.admin import AdminClient, NewTopic

logger = logging.getLogger(__name__)


class KafkaProducerService:
    def __init__(self, bootstrap_servers: str, topic: str):
        self.topic = topic
        self._producer = Producer({
            "bootstrap.servers":     bootstrap_servers,
            "acks":                  "all",
            "enable.idempotence":    True,
            "retries":               5,
            "retry.backoff.ms":      100,
            "linger.ms":             5,
            "batch.size":            65536,
            "compression.type":      "lz4",
            "delivery.timeout.ms":   30000,
        })

    def send(self, key: str, value: dict, headers: Optional[dict] = None) -> None:
        kafka_headers = [(k, v.encode()) for k, v in (headers or {}).items()]
        self._producer.produce(
            topic    = self.topic,
            key      = key.encode("utf-8"),
            value    = json.dumps(value).encode("utf-8"),
            headers  = kafka_headers,
            callback = self._on_delivery,
        )
        # Poll for callbacks to prevent buffer overflow
        self._producer.poll(0)

    def flush(self, timeout: float = 10.0) -> None:
        self._producer.flush(timeout)

    @staticmethod
    def _on_delivery(err, msg) -> None:
        if err:
            logger.error(f"Delivery failed: topic={msg.topic()} err={err}")
        else:
            logger.debug(f"Delivered: topic={msg.topic()} partition={msg.partition()} offset={msg.offset()}")


# Usage
producer = KafkaProducerService("localhost:9092", "orders")
producer.send("user-123", {"order_id": "A1", "amount": 99.99})
producer.flush()
```

---

## confluent-kafka: Production Consumer

```python
import json
import signal
import logging
from typing import Callable
from confluent_kafka import Consumer, KafkaError, KafkaException

logger = logging.getLogger(__name__)


class KafkaConsumerService:
    def __init__(self, bootstrap_servers: str, group_id: str, topics: list[str]):
        self._running = True
        self._consumer = Consumer({
            "bootstrap.servers":      bootstrap_servers,
            "group.id":               group_id,
            "auto.offset.reset":      "earliest",
            "enable.auto.commit":     False,
            "max.poll.records":       500,
            "session.timeout.ms":     30000,
            "heartbeat.interval.ms":  3000,
            "max.poll.interval.ms":   300000,
        })
        self._consumer.subscribe(topics)

        # Graceful shutdown
        signal.signal(signal.SIGTERM, lambda *_: self.stop())
        signal.signal(signal.SIGINT,  lambda *_: self.stop())

    def run(self, handler: Callable[[str, dict, dict], None]) -> None:
        try:
            while self._running:
                msg = self._consumer.poll(timeout=1.0)

                if msg is None:
                    continue
                if msg.error():
                    if msg.error().code() == KafkaError._PARTITION_EOF:
                        continue
                    raise KafkaException(msg.error())

                key     = msg.key().decode("utf-8")   if msg.key()   else None
                value   = json.loads(msg.value())      if msg.value() else {}
                headers = {k: v.decode() for k, v in (msg.headers() or [])}

                try:
                    handler(key, value, headers)
                    self._consumer.commit(asynchronous=False)
                except Exception:
                    logger.exception(f"Handler failed for offset {msg.offset()}")
                    # Don't commit — will be reprocessed
        finally:
            self._consumer.close()

    def stop(self) -> None:
        self._running = False


# Usage
def handle_order(key: str, value: dict, headers: dict) -> None:
    print(f"Processing order {key}: {value}")
    # ... business logic ...

consumer = KafkaConsumerService("localhost:9092", "order-processor", ["orders"])
consumer.run(handle_order)
```

---

## kafka-python

```python
from kafka import KafkaProducer, KafkaConsumer
import json

# Producer
producer = KafkaProducer(
    bootstrap_servers  = "localhost:9092",
    value_serializer   = lambda v: json.dumps(v).encode("utf-8"),
    key_serializer     = lambda k: k.encode("utf-8"),
    acks               = "all",
    retries            = 3,
    compression_type   = "lz4",
)

producer.send("orders", key="user-1", value={"amount": 99.99})
producer.flush()
producer.close()

# Consumer
consumer = KafkaConsumer(
    "orders",
    bootstrap_servers  = "localhost:9092",
    group_id           = "my-app",
    auto_offset_reset  = "earliest",
    enable_auto_commit = False,
    value_deserializer = lambda v: json.loads(v.decode("utf-8")),
    key_deserializer   = lambda k: k.decode("utf-8") if k else None,
)

for msg in consumer:
    print(f"Received: {msg.key} → {msg.value}")
    consumer.commit()
```

---

## Faust: Stream Processing in Python

```python
import faust

app = faust.App(
    "order-processor",
    broker       = "kafka://localhost:9092",
    value_serializer = "json",
)

# Define message types
class Order(faust.Record):
    order_id: str
    user_id:  str
    amount:   float

# Input/output topics
orders_topic   = app.topic("orders",   value_type=Order)
enriched_topic = app.topic("enriched-orders", value_type=dict)

# Stateful table (KTable equivalent)
order_counts = app.Table("order-counts", default=int)

# Stream processor
@app.agent(orders_topic)
async def process_orders(orders):
    async for order in orders:
        # Stateful: count per user
        order_counts[order.user_id] += 1

        # Filter and route
        if order.amount > 1000:
            await enriched_topic.send(
                key   = order.order_id,
                value = {**order.asdict(), "tier": "premium",
                         "user_order_count": order_counts[order.user_id]}
            )

# Windowed aggregation
@app.agent(orders_topic)
async def window_count(orders):
    async for order in orders.group_by(Order.user_id):
        pass  # grouped by user_id

# Cron job
@app.timer(interval=60.0)
async def log_counts():
    print(f"Order counts: {dict(order_counts.items())}")

if __name__ == "__main__":
    app.main()
```

```bash
faust -A app worker -l info
```

---

## aiokafka (Async)

```python
import asyncio
import json
from aiokafka import AIOKafkaProducer, AIOKafkaConsumer

async def produce():
    producer = AIOKafkaProducer(
        bootstrap_servers = "localhost:9092",
        value_serializer  = lambda v: json.dumps(v).encode(),
    )
    await producer.start()
    try:
        await producer.send_and_wait("orders", value={"amount": 99.99})
    finally:
        await producer.stop()

async def consume():
    consumer = AIOKafkaConsumer(
        "orders",
        bootstrap_servers = "localhost:9092",
        group_id          = "async-group",
        auto_offset_reset = "earliest",
    )
    await consumer.start()
    try:
        async for msg in consumer:
            data = json.loads(msg.value)
            print(f"Received: {data}")
    finally:
        await consumer.stop()

asyncio.run(consume())
```

---

## Dead Letter Queue Pattern

```python
from confluent_kafka import Producer, Consumer
import json, logging

logger = logging.getLogger(__name__)

dlq_producer = Producer({"bootstrap.servers": "localhost:9092"})

def process_with_dlq(msg):
    try:
        data = json.loads(msg.value())
        process_order(data)
    except Exception as e:
        logger.error(f"Failed to process message, sending to DLQ: {e}")
        dlq_producer.produce(
            topic   = f"{msg.topic()}.dlq",
            key     = msg.key(),
            value   = msg.value(),
            headers = {
                "error":           str(e),
                "original-topic":  msg.topic(),
                "original-offset": str(msg.offset()),
                "original-partition": str(msg.partition()),
            },
        )
        dlq_producer.flush()
```

---

## Transactional Exactly-Once

```python
from confluent_kafka import Producer, Consumer, KafkaError

producer = Producer({
    "bootstrap.servers":  "localhost:9092",
    "transactional.id":   "transform-pipeline-1",
    "enable.idempotence": True,
})
producer.init_transactions()

consumer = Consumer({
    "bootstrap.servers":              "localhost:9092",
    "group.id":                       "transform-group",
    "isolation.level":                "read_committed",  # skip uncommitted
    "enable.auto.commit":             False,
})
consumer.subscribe(["raw-orders"])

while True:
    msg = consumer.poll(1.0)
    if msg is None or msg.error(): continue

    try:
        producer.begin_transaction()
        transformed = transform(json.loads(msg.value()))
        producer.produce("processed-orders", value=json.dumps(transformed))

        # Commit offsets atomically with the transaction
        producer.send_offsets_to_transaction(
            consumer.position(consumer.assignment()),
            consumer.consumer_group_metadata()
        )
        producer.commit_transaction()
    except Exception:
        producer.abort_transaction()
```


---

# Chapter 16: Kafka with Java

## Maven Dependency

```xml
<dependency>
    <groupId>org.apache.kafka</groupId>
    <artifactId>kafka-clients</artifactId>
    <version>3.7.0</version>
</dependency>

<!-- Spring Kafka -->
<dependency>
    <groupId>org.springframework.kafka</groupId>
    <artifactId>spring-kafka</artifactId>
    <version>3.1.0</version>
</dependency>
```

---

## Java Producer

```java
import org.apache.kafka.clients.producer.*;
import org.apache.kafka.common.serialization.StringSerializer;
import java.util.Properties;
import java.util.concurrent.ExecutionException;

public class OrderProducer {

    public static void main(String[] args) throws ExecutionException, InterruptedException {

        Properties props = new Properties();
        props.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG,  "localhost:9092");
        props.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG,   StringSerializer.class.getName());
        props.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        props.put(ProducerConfig.ACKS_CONFIG,               "all");
        props.put(ProducerConfig.ENABLE_IDEMPOTENCE_CONFIG, true);
        props.put(ProducerConfig.LINGER_MS_CONFIG,          5);
        props.put(ProducerConfig.BATCH_SIZE_CONFIG,         65536);
        props.put(ProducerConfig.COMPRESSION_TYPE_CONFIG,   "lz4");

        try (KafkaProducer<String, String> producer = new KafkaProducer<>(props)) {

            // Async send with callback
            ProducerRecord<String, String> record = new ProducerRecord<>(
                "orders",
                "user-123",
                "{\"order_id\":\"A1\",\"amount\":99.99}"
            );

            producer.send(record, (metadata, exception) -> {
                if (exception != null) {
                    System.err.println("Send failed: " + exception.getMessage());
                } else {
                    System.out.printf("Sent to %s[%d]@%d%n",
                        metadata.topic(), metadata.partition(), metadata.offset());
                }
            });

            // Synchronous send
            RecordMetadata meta = producer.send(record).get();
            System.out.printf("Sync send: partition=%d offset=%d%n",
                meta.partition(), meta.offset());

            producer.flush();
        }
    }
}
```

---

## Java Consumer

```java
import org.apache.kafka.clients.consumer.*;
import org.apache.kafka.common.serialization.StringDeserializer;
import java.time.Duration;
import java.util.*;

public class OrderConsumer {

    public static void main(String[] args) {
        Properties props = new Properties();
        props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG,   "localhost:9092");
        props.put(ConsumerConfig.GROUP_ID_CONFIG,            "order-processor");
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG,   StringDeserializer.class.getName());
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG,   "earliest");
        props.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG,  false);
        props.put(ConsumerConfig.MAX_POLL_RECORDS_CONFIG,    500);

        try (KafkaConsumer<String, String> consumer = new KafkaConsumer<>(props)) {
            consumer.subscribe(List.of("orders"));

            Runtime.getRuntime().addShutdownHook(new Thread(consumer::wakeup));

            try {
                while (true) {
                    ConsumerRecords<String, String> records =
                        consumer.poll(Duration.ofMillis(1000));

                    for (ConsumerRecord<String, String> record : records) {
                        System.out.printf("topic=%s partition=%d offset=%d key=%s value=%s%n",
                            record.topic(), record.partition(), record.offset(),
                            record.key(), record.value());

                        processOrder(record.key(), record.value());
                    }

                    // Commit after processing batch
                    consumer.commitSync();
                }
            } catch (WakeupException e) {
                // Shutdown signal, exit cleanly
            }
        }
    }

    static void processOrder(String key, String value) {
        // Business logic
    }
}
```

---

## Transactions in Java

```java
Properties props = new Properties();
props.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
props.put(ProducerConfig.TRANSACTIONAL_ID_CONFIG,  "order-enricher-1");
props.put(ProducerConfig.ENABLE_IDEMPOTENCE_CONFIG, true);
props.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG,   StringSerializer.class.getName());
props.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());

KafkaProducer<String, String> producer = new KafkaProducer<>(props);
producer.initTransactions();

Properties consumerProps = new Properties();
consumerProps.put(ConsumerConfig.ISOLATION_LEVEL_CONFIG, "read_committed");
// ... other consumer config ...
KafkaConsumer<String, String> consumer = new KafkaConsumer<>(consumerProps);
consumer.subscribe(List.of("raw-orders"));

while (true) {
    ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(100));
    if (records.isEmpty()) continue;

    try {
        producer.beginTransaction();

        for (ConsumerRecord<String, String> record : records) {
            String enriched = enrich(record.value());
            producer.send(new ProducerRecord<>("enriched-orders", record.key(), enriched));
        }

        // Commit offsets atomically with transaction
        Map<TopicPartition, OffsetAndMetadata> offsets = new HashMap<>();
        for (TopicPartition tp : consumer.assignment()) {
            offsets.put(tp, new OffsetAndMetadata(consumer.position(tp)));
        }
        producer.sendOffsetsToTransaction(offsets, consumer.groupMetadata());
        producer.commitTransaction();

    } catch (Exception e) {
        producer.abortTransaction();
        consumer.seekToBeginning(consumer.assignment());
    }
}
```

---

## Spring Kafka

### Configuration

```java
@Configuration
@EnableKafka
public class KafkaConfig {

    @Bean
    public ProducerFactory<String, Object> producerFactory() {
        Map<String, Object> config = new HashMap<>();
        config.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
        config.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG,   StringSerializer.class);
        config.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, JsonSerializer.class);
        config.put(ProducerConfig.ACKS_CONFIG, "all");
        config.put(ProducerConfig.ENABLE_IDEMPOTENCE_CONFIG, true);
        return new DefaultKafkaProducerFactory<>(config);
    }

    @Bean
    public KafkaTemplate<String, Object> kafkaTemplate() {
        return new KafkaTemplate<>(producerFactory());
    }

    @Bean
    public ConsumerFactory<String, Object> consumerFactory() {
        Map<String, Object> config = new HashMap<>();
        config.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
        config.put(ConsumerConfig.GROUP_ID_CONFIG,          "spring-group");
        config.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG,   StringDeserializer.class);
        config.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, JsonDeserializer.class);
        config.put(JsonDeserializer.TRUSTED_PACKAGES, "*");
        return new DefaultKafkaConsumerFactory<>(config);
    }

    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, Object> kafkaListenerContainerFactory() {
        ConcurrentKafkaListenerContainerFactory<String, Object> factory =
            new ConcurrentKafkaListenerContainerFactory<>();
        factory.setConsumerFactory(consumerFactory());
        factory.setConcurrency(3);  // 3 threads (one per partition)
        factory.getContainerProperties().setAckMode(AckMode.MANUAL_IMMEDIATE);
        return factory;
    }
}
```

### Producer Service

```java
@Service
public class OrderProducerService {

    @Autowired
    private KafkaTemplate<String, Object> kafkaTemplate;

    public void sendOrder(Order order) {
        kafkaTemplate.send("orders", order.getOrderId(), order)
            .whenComplete((result, ex) -> {
                if (ex != null) {
                    log.error("Failed to send order {}: {}", order.getOrderId(), ex.getMessage());
                } else {
                    log.info("Order sent: partition={} offset={}",
                        result.getRecordMetadata().partition(),
                        result.getRecordMetadata().offset());
                }
            });
    }
}
```

### Consumer (Listener)

```java
@Component
public class OrderConsumer {

    @KafkaListener(
        topics   = "orders",
        groupId  = "order-processor",
        concurrency = "3"   // 3 consumer threads
    )
    public void consume(
        @Payload Order order,
        @Header(KafkaHeaders.RECEIVED_PARTITION) int partition,
        @Header(KafkaHeaders.OFFSET) long offset,
        Acknowledgment ack
    ) {
        try {
            processOrder(order);
            ack.acknowledge();   // manual commit
        } catch (Exception e) {
            log.error("Failed to process order at offset {}: {}", offset, e.getMessage());
            // Don't ack — will be redelivered
        }
    }

    // Batch listener
    @KafkaListener(topics = "events", containerFactory = "batchFactory")
    public void consumeBatch(List<Order> orders, Acknowledgment ack) {
        processBatch(orders);
        ack.acknowledge();
    }
}
```

---

## Error Handling and Retry in Spring Kafka

```java
@Bean
public DefaultErrorHandler errorHandler(KafkaOperations<?, ?> template) {
    // Fixed backoff: 3 retries, 1 second apart
    FixedBackOff backOff = new FixedBackOff(1000L, 3);

    DeadLetterPublishingRecoverer recoverer =
        new DeadLetterPublishingRecoverer(template,
            (rec, ex) -> new TopicPartition(rec.topic() + ".dlq", rec.partition()));

    DefaultErrorHandler handler = new DefaultErrorHandler(recoverer, backOff);

    // Don't retry on these exceptions
    handler.addNotRetryableExceptions(
        DeserializationException.class,
        NullPointerException.class
    );

    return handler;
}
```


---

# Chapter 17: Exactly-Once Semantics

## Delivery Guarantees

```
At-most-once:   Message may be lost, never duplicated
                Commit offset BEFORE processing
                Use: non-critical metrics, fire-and-forget logs

At-least-once:  Message never lost, may be duplicated
                Commit offset AFTER processing
                Use: most business applications (with idempotent consumers)

Exactly-once:   Message delivered and processed exactly once
                Requires: idempotent producer + transactions + read_committed
                Use: financial transactions, inventory, billing
```

---

## Building Blocks of EOS

### 1. Idempotent Producer

Prevents duplicates caused by producer retries:

```
Without idempotence:
  Producer sends msg → network failure → retry
  Broker received original + retry = DUPLICATE

With idempotence (PID + sequence number):
  Producer sends msg (PID=1, seq=42) → network failure → retry
  Broker sees (PID=1, seq=42) again → deduplicates → EXACTLY ONCE per partition
```

```python
producer = Producer({
    "bootstrap.servers":  "localhost:9092",
    "enable.idempotence": True,
    # Automatically set:
    # acks = all
    # max.in.flight.requests.per.connection = 5
    # retries = INT_MAX
})
```

**Scope:** Idempotent producer prevents duplicates within a single producer session. PID resets on producer restart.

---

### 2. Transactions

Transactions enable **atomic writes across multiple partitions and topics**:

```python
from confluent_kafka import Producer

producer = Producer({
    "bootstrap.servers":  "localhost:9092",
    "transactional.id":   "payment-processor-1",
    "enable.idempotence": True,
})

producer.init_transactions()

# Transaction 1: charge customer + update inventory atomically
producer.begin_transaction()
producer.produce("payment-events",   key="order-1", value="charged")
producer.produce("inventory-updates", key="prod-A",  value="decrement")
producer.commit_transaction()

# If any produce fails:
try:
    producer.begin_transaction()
    producer.produce("payment-events", key="order-2", value="charged")
    raise Exception("something went wrong")
    producer.commit_transaction()
except Exception:
    producer.abort_transaction()
    # Both messages are rolled back — neither visible to read_committed consumers
```

**Transaction coordinator:** A special broker partition manages transaction state. Identified by `transactional.id`.

---

### 3. Consumer read_committed

Consumers must use `isolation.level=read_committed` to skip uncommitted or aborted transactions:

```python
consumer = Consumer({
    "bootstrap.servers": "localhost:9092",
    "group.id":          "payment-auditor",
    "isolation.level":   "read_committed",  # skip uncommitted + aborted
    # Default is "read_uncommitted" — reads everything including aborted
})
```

```
Without read_committed:
  Consumer reads all messages including those in open transactions
  → sees partial state of multi-topic atomic writes

With read_committed:
  Consumer only reads messages from committed transactions
  → consistent view of multi-topic writes
```

---

## End-to-End Exactly-Once: Consume-Transform-Produce

The canonical EOS pattern — read from input, transform, write to output, commit offsets **atomically**:

```python
from confluent_kafka import Producer, Consumer, KafkaError

# Producer with transactional.id
producer = Producer({
    "bootstrap.servers": "localhost:9092",
    "transactional.id":  "transformer-instance-1",
    "enable.idempotence": True,
})
producer.init_transactions()

# Consumer reading committed only
consumer = Consumer({
    "bootstrap.servers":   "localhost:9092",
    "group.id":            "transformer-group",
    "isolation.level":     "read_committed",
    "enable.auto.commit":  False,
    "auto.offset.reset":   "earliest",
})
consumer.subscribe(["raw-events"])

import json

while True:
    messages = []
    while len(messages) < 100:
        msg = consumer.poll(0.1)
        if msg is None or msg.error(): break
        messages.append(msg)

    if not messages:
        continue

    try:
        producer.begin_transaction()

        for msg in messages:
            event = json.loads(msg.value())
            enriched = enrich(event)
            producer.produce("enriched-events",
                             key=msg.key(),
                             value=json.dumps(enriched).encode())

        # Commit consumer offsets INSIDE the transaction
        producer.send_offsets_to_transaction(
            {tp: offset for tp, offset in
             [(msg.topic_partition(), msg.offset() + 1) for msg in messages]},
            consumer.consumer_group_metadata()
        )

        producer.commit_transaction()

    except Exception as e:
        print(f"Transaction failed: {e}, aborting")
        producer.abort_transaction()
        # Seek back to uncommitted offsets — messages will be reprocessed
        consumer.seek_to_committed(consumer.assignment())
```

---

## Transaction Timeout and Fencing

### Zombie Fencing

Prevents "zombie" producer instances from interfering:

```
Scenario: Producer with transactional.id="proc-1" crashes mid-transaction
  New instance starts with same transactional.id="proc-1"
  Broker fences the old instance (bumps epoch)
  Old instance's subsequent requests fail with ProducerFencedException
  New instance's transactions proceed normally
```

```properties
# Maximum time a transaction can stay open
transaction.timeout.ms=60000   (default: 60 seconds)
# Longer timeout → transactions can hold up consumers (read_committed must wait)
```

---

## EOS Performance Impact

```
Overhead of transactions:
1. Extra round trip: begin_transaction, commit_transaction
2. Transaction coordinator load
3. Consumers with read_committed have higher latency (wait for commit markers)

Typical overhead: 5–15% throughput reduction vs at-least-once
```

---

## When to Use EOS

| Use Case | Semantics | Why |
|---|---|---|
| Financial ledger, billing | Exactly-once | No duplicate charges |
| Inventory management | Exactly-once | Prevent over-decrement |
| Real-time analytics | At-least-once | Duplicates skew by tiny % |
| Log aggregation | At-least-once | Duplicate logs acceptable |
| Metrics, heartbeats | At-most-once | Loss OK, duplicates worse |
| CDC (Debezium) | At-least-once + idempotent consumer | Upsert = natural idempotency |

---

## Idempotent Consumer Pattern

Even with at-least-once delivery, make processing idempotent:

```python
import redis

r = redis.Redis(host="localhost", port=6379)

def process_order(order_id: str, order_data: dict) -> None:
    # Use Redis SET NX (set if not exists) as dedup check
    already_processed = not r.set(f"processed:{order_id}", "1",
                                   nx=True, ex=86400)
    if already_processed:
        return  # skip duplicate

    # Process the order
    db.insert_order(order_data)
    r.delete(f"processed:{order_id}")  # optional cleanup
```

Or use a database upsert:

```sql
INSERT INTO orders (order_id, user_id, amount, processed_at)
VALUES (%s, %s, %s, NOW())
ON CONFLICT (order_id) DO NOTHING;   -- idempotent: safe to retry
```


---

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


---

# Chapter 19: Kafka with Spark & Flink

## Kafka + Apache Spark Structured Streaming

### Read from Kafka

```python
from pyspark.sql import SparkSession
from pyspark.sql import functions as F
from pyspark.sql.types import StructType, StructField, StringType, DoubleType, TimestampType

spark = SparkSession.builder \
    .appName("KafkaStreamProcessor") \
    .config("spark.jars.packages",
            "org.apache.spark:spark-sql-kafka-0-10_2.12:3.5.0") \
    .getOrCreate()

# Read stream from Kafka
raw = spark.readStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers",  "broker:9092") \
    .option("subscribe",                "orders,events") \
    .option("startingOffsets",          "latest") \
    .option("maxOffsetsPerTrigger",     50000) \
    .option("kafka.group.id",           "spark-consumer") \
    .option("kafka.security.protocol",  "SASL_SSL") \
    .option("kafka.sasl.mechanism",     "PLAIN") \
    .load()

# Kafka schema: key, value, topic, partition, offset, timestamp, timestampType
# Both key and value are binary — cast to string/parse
```

### Parse JSON Messages

```python
order_schema = StructType([
    StructField("order_id",   StringType()),
    StructField("user_id",    StringType()),
    StructField("amount",     DoubleType()),
    StructField("currency",   StringType()),
    StructField("created_at", TimestampType()),
])

orders = raw \
    .filter(F.col("topic") == "orders") \
    .select(
        F.from_json(F.col("value").cast("string"), order_schema).alias("data"),
        "topic", "partition", "offset", "timestamp"
    ) \
    .select("data.*", "topic", "partition", "offset", "timestamp")
```

### Windowed Aggregation

```python
# Revenue per user per 10-minute window
revenue = orders \
    .withWatermark("created_at", "5 minutes") \
    .groupBy(
        F.window("created_at", "10 minutes", "5 minutes"),
        "user_id"
    ) \
    .agg(
        F.sum("amount").alias("total_revenue"),
        F.count("*").alias("order_count")
    )

query = revenue.writeStream \
    .format("console") \
    .outputMode("update") \
    .option("checkpointLocation", "/tmp/checkpoints/revenue/") \
    .trigger(processingTime="30 seconds") \
    .start()

query.awaitTermination()
```

### Write Back to Kafka

```python
# Transform and write to another Kafka topic
enriched = orders \
    .withColumn("tier", F.when(F.col("amount") > 1000, "premium").otherwise("standard")) \
    .withColumn("processed_at", F.current_timestamp())

query = enriched \
    .select(
        F.col("order_id").cast("string").alias("key"),
        F.to_json(F.struct("*")).alias("value"),
    ) \
    .writeStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers",   "broker:9092") \
    .option("topic",                     "enriched-orders") \
    .option("checkpointLocation",        "/tmp/checkpoints/enriched/") \
    .start()
```

### foreachBatch for Complex Sinks

```python
def process_batch(batch_df, batch_id):
    # Write to multiple destinations
    batch_df.write \
        .format("delta") \
        .mode("append") \
        .option("path", "s3://my-bucket/delta/orders/") \
        .save()

    # Write aggregated metrics to PostgreSQL
    batch_df.groupBy("user_id") \
        .agg(F.sum("amount").alias("total")) \
        .write \
        .format("jdbc") \
        .option("url", "jdbc:postgresql://db:5432/analytics") \
        .option("dbtable", "user_metrics") \
        .mode("append") \
        .save()

query = orders.writeStream \
    .foreachBatch(process_batch) \
    .option("checkpointLocation", "/tmp/checkpoints/batch/") \
    .trigger(processingTime="1 minute") \
    .start()
```

---

## Kafka + Apache Flink

### Maven Dependency

```xml
<dependency>
    <groupId>org.apache.flink</groupId>
    <artifactId>flink-connector-kafka</artifactId>
    <version>3.2.0-1.19</version>
</dependency>
```

### Flink Kafka Source

```java
import org.apache.flink.api.common.serialization.SimpleStringSchema;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;

KafkaSource<String> source = KafkaSource.<String>builder()
    .setBootstrapServers("localhost:9092")
    .setTopics("orders")
    .setGroupId("flink-consumer")
    .setStartingOffsets(OffsetsInitializer.earliest())
    .setValueOnlyDeserializer(new SimpleStringSchema())
    .build();

StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

DataStream<String> stream = env.fromSource(
    source,
    WatermarkStrategy.forMonotonousTimestamps(),
    "Kafka Source"
);
```

### Flink Kafka Sink

```java
import org.apache.flink.connector.kafka.sink.KafkaRecordSerializationSchema;
import org.apache.flink.connector.kafka.sink.KafkaSink;

KafkaSink<String> sink = KafkaSink.<String>builder()
    .setBootstrapServers("localhost:9092")
    .setRecordSerializer(KafkaRecordSerializationSchema.builder()
        .setTopic("processed-orders")
        .setValueSerializationSchema(new SimpleStringSchema())
        .build())
    .setDeliveryGuarantee(DeliveryGuarantee.EXACTLY_ONCE)
    .setTransactionalIdPrefix("flink-tx-")
    .setKafkaProducerConfig(producerProps)
    .build();

stream.sinkTo(sink);
```

### Flink Windowing

```java
import org.apache.flink.api.common.eventtime.*;
import org.apache.flink.streaming.api.windowing.assigners.TumblingEventTimeWindows;
import java.time.Duration;

DataStream<Order> orders = rawStream
    .map(json -> objectMapper.readValue(json, Order.class))
    .assignTimestampsAndWatermarks(
        WatermarkStrategy
            .<Order>forBoundedOutOfOrderness(Duration.ofMinutes(1))
            .withTimestampAssigner((event, ts) -> event.getCreatedAt().toEpochMilli())
    );

// Tumbling window: revenue per user per 5 minutes
DataStream<UserRevenue> result = orders
    .keyBy(Order::getUserId)
    .window(TumblingEventTimeWindows.of(Duration.ofMinutes(5)))
    .aggregate(new RevenueAggregator());
```

### Flink Exactly-Once with Kafka

```java
// Enable checkpointing
env.enableCheckpointing(60000, CheckpointingMode.EXACTLY_ONCE);
env.getCheckpointConfig().setCheckpointStorage("s3://my-bucket/flink-checkpoints/");
env.getCheckpointConfig().setMinPauseBetweenCheckpoints(30000);
env.getCheckpointConfig().setCheckpointTimeout(120000);

// Kafka source with exactly-once
KafkaSource<Order> source = KafkaSource.<Order>builder()
    .setBootstrapServers("localhost:9092")
    .setTopics("orders")
    .setGroupId("flink-eos")
    .setStartingOffsets(OffsetsInitializer.committedOffsets())
    .setDeserializer(KafkaRecordDeserializationSchema.valueOnly(new OrderDeserializationSchema()))
    .build();

// Kafka sink with exactly-once transactions
KafkaSink<String> sink = KafkaSink.<String>builder()
    .setBootstrapServers("localhost:9092")
    .setDeliveryGuarantee(DeliveryGuarantee.EXACTLY_ONCE)
    .setTransactionalIdPrefix("flink-eos-")
    // ...
    .build();
```

---

## Spark vs Flink for Kafka Consumers

| Feature | Spark Structured Streaming | Apache Flink |
|---|---|---|
| Processing model | Micro-batch (default) / continuous | True streaming |
| Latency | 100ms – seconds | < 10ms |
| State management | In-memory / HDFS | RocksDB / heap |
| Windowing | Tumbling, sliding, session | Rich (all types) |
| Exactly-once | Yes (foreachBatch) | Yes (native) |
| Language | Python, Java, Scala, R | Java, Scala, Python, SQL |
| SQL support | Spark SQL | Flink SQL (powerful) |
| Learning curve | Easier | Steeper |
| Ecosystem | Spark ecosystem | Flink ecosystem |
| Cloud managed | Databricks, EMR | Confluent, Kinesis Data Analytics |

---

## End-to-End Pipeline

```
Kafka (raw-events)
       │
       ├──► Flink (real-time, < 1s latency)
       │      └──► Kafka (alerts) ──► PagerDuty
       │
       └──► Spark Structured Streaming (batch, 1-min windows)
              ├──► Delta Lake (for BI queries)
              └──► Kafka (enriched-events)
                       └──► Elasticsearch (Kibana dashboards)
```


---

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
