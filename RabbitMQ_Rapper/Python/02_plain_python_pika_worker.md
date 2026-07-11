# Chapter 2: Plain Python pika Worker

## Connection

```python
import pika

def make_connection(config: dict) -> pika.BlockingConnection:
    credentials = pika.PlainCredentials(config["user"], config["password"])

    params = pika.ConnectionParameters(
        host=config["host"],
        port=config["port"],
        virtual_host=config.get("vhost", "/"),
        credentials=credentials,
        heartbeat=60,
        blocked_connection_timeout=30,
    )

    return pika.BlockingConnection(params)
```

---

## Topology

```python
def declare_topology(channel) -> None:
    channel.exchange_declare(
        exchange="app.events",
        exchange_type="topic",
        durable=True,
    )

    channel.exchange_declare(
        exchange="app.dlx",
        exchange_type="topic",
        durable=True,
    )

    channel.queue_declare(
        queue="product.search_index.queue",
        durable=True,
        arguments={
            "x-dead-letter-exchange": "app.dlx",
        },
    )

    channel.queue_bind(
        queue="product.search_index.queue",
        exchange="app.events",
        routing_key="product.*",
    )

    channel.queue_declare(queue="product.search_index.dlq", durable=True)
    channel.queue_bind(
        queue="product.search_index.dlq",
        exchange="app.dlx",
        routing_key="#",
    )
```

---

## Publisher

```python
import uuid
import pika

class RabbitPublisher:
    def __init__(self, channel, serializer, exchange: str = "app.events"):
        self.channel = channel
        self.serializer = serializer
        self.exchange = exchange

    def publish(self, routing_key: str, payload: dict, meta: dict | None = None) -> None:
        message = MessageEnvelope(
            id=str(uuid.uuid4()),
            type=routing_key,
            version=1,
            occurred_at=MessageEnvelope.now_iso(),
            payload=payload,
            meta=meta or {},
        )

        self.channel.basic_publish(
            exchange=self.exchange,
            routing_key=routing_key,
            body=self.serializer.encode(message),
            properties=pika.BasicProperties(
                content_type="application/json",
                delivery_mode=pika.DeliveryMode.Persistent,
                message_id=message.id,
                correlation_id=message.meta.get("correlation_id"),
            ),
        )
```

---

## Handler Registry

```python
class HandlerRegistry:
    def __init__(self):
        self.handlers = {}

    def register(self, message_type: str, handler) -> None:
        self.handlers[message_type] = handler

    def get(self, message_type: str):
        try:
            return self.handlers[message_type]
        except KeyError:
            raise PermanentMessageError(f"No handler for {message_type}")
```

---

## Consumer

```python
import logging

logger = logging.getLogger(__name__)

class RabbitConsumer:
    def __init__(self, channel, serializer, registry):
        self.channel = channel
        self.serializer = serializer
        self.registry = registry

    def consume(self, queue: str) -> None:
        self.channel.basic_qos(prefetch_count=1)
        self.channel.basic_consume(
            queue=queue,
            on_message_callback=self.process,
            auto_ack=False,
        )
        self.channel.start_consuming()

    def process(self, channel, method, properties, body):
        try:
            message = self.serializer.decode(body)
            handler = self.registry.get(message.type)
            handler.handle(message)

            channel.basic_ack(delivery_tag=method.delivery_tag)

            logger.info("rabbitmq.message_processed", extra={
                "message_id": message.id,
                "type": message.type,
            })
        except PermanentMessageError as exc:
            logger.exception("rabbitmq.permanent_failure")
            channel.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
        except Exception as exc:
            logger.exception("rabbitmq.unknown_failure")
            channel.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
```

---

## Handler

```python
class ProductUpdatedHandler:
    def handle(self, message: MessageEnvelope) -> None:
        product_id = message.payload.get("product_id")

        if not product_id:
            raise PermanentMessageError("product_id is required")

        print(f"Index product {product_id}")
```

---

## Worker Script

```python
import sys

config = {
    "host": "127.0.0.1",
    "port": 5672,
    "user": "guest",
    "password": "guest",
}

connection = make_connection(config)
channel = connection.channel()
declare_topology(channel)

registry = HandlerRegistry()
registry.register("product.updated", ProductUpdatedHandler())

consumer = RabbitConsumer(
    channel=channel,
    serializer=JsonMessageSerializer(),
    registry=registry,
)

consumer.consume(sys.argv[1] if len(sys.argv) > 1 else "product.search_index.queue")
```

Run:

```bash
python worker.py product.search_index.queue
```

