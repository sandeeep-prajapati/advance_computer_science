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
