# Python RabbitMQ Wrapper: Complete Book

> Detailed Python and Django RabbitMQ wrapper notes.

---

# Python RabbitMQ Wrapper Book

Detailed notes for building RabbitMQ wrappers in Python and Django.

This folder covers:

- plain Python worker architecture
- `pika` publisher and consumer
- message envelope and handler registry
- Django settings integration
- Django management command workers
- retries, DLQ, idempotency, and operations

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_python_package_architecture.md) | **Python Package Architecture** | package layout, dataclasses, interfaces, settings |
| [02](02_plain_python_pika_worker.md) | **Plain Python pika Worker** | connection, publisher, consumer, handler registry |
| [03](03_django_integration.md) | **Django Integration** | settings, services, management command, views publishing messages |
| [04](04_python_retries_dlq_operations.md) | **Python Retries, DLQ & Operations** | retry policy, DLQ, idempotency, logging, process managers |
| [book](book.md) | **Combined Python Book** | All Python chapters in one file |

---

## Python Wrapper Goal

Publishing should be simple:

```python
publisher.publish("product.updated", {"product_id": 501})
```

Consuming should be a command:

```bash
python worker.py product.search_index.queue
python manage.py consume_rabbitmq product.search_index.queue
```

Business handlers should receive a clean `MessageEnvelope`, not raw `pika` delivery objects.


---

# Chapter 1: Python Package Architecture

## Package Structure

```text
python-rabbit-app/
├── pyproject.toml
├── messaging/
│   ├── __init__.py
│   ├── envelope.py
│   ├── serializer.py
│   ├── publisher.py
│   ├── consumer.py
│   ├── registry.py
│   ├── topology.py
│   └── exceptions.py
├── handlers/
│   ├── __init__.py
│   └── product_updated.py
└── worker.py
```

Install:

```bash
pip install pika
```

For Django:

```bash
pip install pika django
```

---

## Message Envelope

```python
from dataclasses import dataclass, field
from datetime import datetime, timezone
from typing import Any

@dataclass(frozen=True)
class MessageEnvelope:
    id: str
    type: str
    version: int
    occurred_at: str
    payload: dict[str, Any]
    meta: dict[str, Any] = field(default_factory=dict)

    @staticmethod
    def now_iso() -> str:
        return datetime.now(timezone.utc).isoformat()
```

---

## Handler Protocol

```python
from typing import Protocol

class MessageHandler(Protocol):
    def handle(self, message: MessageEnvelope) -> None:
        ...
```

Using a protocol keeps handlers easy to test.

---

## Serializer

```python
import json

class JsonMessageSerializer:
    def encode(self, message: MessageEnvelope) -> bytes:
        return json.dumps({
            "id": message.id,
            "type": message.type,
            "version": message.version,
            "occurred_at": message.occurred_at,
            "payload": message.payload,
            "meta": message.meta,
        }).encode("utf-8")

    def decode(self, body: bytes) -> MessageEnvelope:
        data = json.loads(body.decode("utf-8"))

        return MessageEnvelope(
            id=data["id"],
            type=data["type"],
            version=int(data["version"]),
            occurred_at=data["occurred_at"],
            payload=data.get("payload", {}),
            meta=data.get("meta", {}),
        )
```

---

## Exceptions

```python
class TemporaryMessageError(Exception):
    pass

class PermanentMessageError(Exception):
    pass
```

Consumer behavior:

```text
TemporaryMessageError -> retry
PermanentMessageError -> DLQ
unknown Exception -> limited retry then DLQ
```

---

## Design Rule

Keep the `messaging` package independent from Django.

Good:

```text
messaging package -> pure Python
django app -> imports messaging and wires settings/commands
```

Bad:

```text
messaging package imports Django models/settings everywhere
```


---

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


---

# Chapter 3: Django Integration

## Django Architecture

```text
django project
├── settings.py
├── messaging/
│   ├── publisher.py
│   ├── consumer.py
│   ├── registry.py
│   ├── handlers.py
│   └── management/
│       └── commands/
│           └── consume_rabbitmq.py
└── products/
    └── views.py
```

Views publish messages. Management commands consume messages.

---

## settings.py

```python
RABBITMQ = {
    "HOST": "127.0.0.1",
    "PORT": 5672,
    "USER": "guest",
    "PASSWORD": "guest",
    "VHOST": "/",
    "EXCHANGE": "app.events",
    "HANDLERS": {
        "product.updated": "messaging.handlers.ProductUpdatedHandler",
    },
}
```

---

## Connection Factory

```python
from django.conf import settings

def rabbitmq_config() -> dict:
    cfg = settings.RABBITMQ
    return {
        "host": cfg["HOST"],
        "port": cfg["PORT"],
        "user": cfg["USER"],
        "password": cfg["PASSWORD"],
        "vhost": cfg.get("VHOST", "/"),
    }
```

---

## Django Handler Registry

```python
from django.conf import settings
from django.utils.module_loading import import_string

class DjangoHandlerRegistry:
    def get(self, message_type: str):
        path = settings.RABBITMQ["HANDLERS"].get(message_type)

        if not path:
            raise PermanentMessageError(f"No handler for {message_type}")

        handler_class = import_string(path)
        return handler_class()
```

For bigger apps, inject services into handlers through factories.

---

## Management Command

`messaging/management/commands/consume_rabbitmq.py`:

```python
from django.core.management.base import BaseCommand

class Command(BaseCommand):
    help = "Consume RabbitMQ messages"

    def add_arguments(self, parser):
        parser.add_argument("queue")

    def handle(self, *args, **options):
        connection = make_connection(rabbitmq_config())
        channel = connection.channel()
        declare_topology(channel)

        consumer = RabbitConsumer(
            channel=channel,
            serializer=JsonMessageSerializer(),
            registry=DjangoHandlerRegistry(),
        )

        consumer.consume(options["queue"])
```

Run:

```bash
python manage.py consume_rabbitmq product.search_index.queue
```

---

## Django View Publishing

```python
from django.http import JsonResponse

def update_product(request, product_id: int):
    connection = make_connection(rabbitmq_config())
    channel = connection.channel()
    declare_topology(channel)

    publisher = RabbitPublisher(
        channel=channel,
        serializer=JsonMessageSerializer(),
        exchange="app.events",
    )

    publisher.publish(
        "product.updated",
        {"product_id": product_id},
        {
            "source": "django",
            "correlation_id": request.headers.get("X-Correlation-ID"),
        },
    )

    connection.close()

    return JsonResponse({"status": "queued"}, status=202)
```

For production, avoid creating a new RabbitMQ connection for every request if request volume is high. Use a managed connection layer or publish through a service.

---

## Django Handler

```python
class ProductUpdatedHandler:
    def handle(self, message: MessageEnvelope) -> None:
        product_id = message.payload["product_id"]

        product = Product.objects.get(id=product_id)
        index_product(product)
```

---

## Django Rule

Do not run long RabbitMQ consumers inside views, middleware, or app startup.

Use:

```bash
python manage.py consume_rabbitmq queue.name
```

as a separate worker process.


---

# Chapter 4: Python Retries, DLQ & Operations

## Retry Policy

```python
class RetryPolicy:
    def __init__(self, max_retries: int = 3):
        self.max_retries = max_retries

    def retry_count(self, message: MessageEnvelope) -> int:
        return int(message.meta.get("retry_count", 0))

    def should_retry(self, message: MessageEnvelope) -> bool:
        return self.retry_count(message) < self.max_retries
```

---

## Idempotency in Django

Model:

```python
class ProcessedMessage(models.Model):
    message_id = models.CharField(max_length=100, unique=True)
    processed_at = models.DateTimeField(auto_now_add=True)
```

Usage:

```python
if ProcessedMessage.objects.filter(message_id=message.id).exists():
    return

process_business_action()

ProcessedMessage.objects.create(message_id=message.id)
```

Use database transactions for important side effects.

---

## Logging

```python
logger.info(
    "rabbitmq.message_processed",
    extra={
        "message_id": message.id,
        "type": message.type,
        "correlation_id": message.meta.get("correlation_id"),
    },
)
```

Prefer structured logs in production.

---

## systemd Service

```ini
[Unit]
Description=Django RabbitMQ Worker
After=network.target

[Service]
WorkingDirectory=/var/www/app
ExecStart=/var/www/app/.venv/bin/python manage.py consume_rabbitmq product.search_index.queue
Restart=always
User=www-data
Environment=DJANGO_SETTINGS_MODULE=config.settings.production

[Install]
WantedBy=multi-user.target
```

---

## Supervisor Config

```ini
[program:django-rabbitmq-worker]
command=/var/www/app/.venv/bin/python manage.py consume_rabbitmq product.search_index.queue
directory=/var/www/app
autostart=true
autorestart=true
redirect_stderr=true
stdout_logfile=/var/log/django-rabbitmq-worker.log
```

---

## Python Production Checklist

- use durable queues and exchanges
- use persistent messages
- use manual ack
- configure prefetch
- add DLQ
- add retry policy
- make handlers idempotent
- use management command workers
- supervise worker processes
- log message id/type/correlation id
- close connections cleanly
- handle SIGTERM in container deployments

