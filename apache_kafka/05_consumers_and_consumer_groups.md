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
