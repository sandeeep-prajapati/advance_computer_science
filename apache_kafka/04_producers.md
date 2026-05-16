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
