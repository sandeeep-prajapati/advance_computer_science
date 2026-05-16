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
