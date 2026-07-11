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

