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

