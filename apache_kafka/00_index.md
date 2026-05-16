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
