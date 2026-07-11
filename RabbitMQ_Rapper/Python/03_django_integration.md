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

