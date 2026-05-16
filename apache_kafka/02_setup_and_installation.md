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
