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
