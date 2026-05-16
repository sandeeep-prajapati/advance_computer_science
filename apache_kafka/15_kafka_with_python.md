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
