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
