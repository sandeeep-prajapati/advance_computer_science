# Chapter 6: PHP Retries, DLQ & Operations

## PHP Worker Risks

Long-running PHP workers need care:

- memory can grow over time
- database connections can go stale
- service container state can become old
- deploys need worker restarts
- fatal errors can kill worker

Frameworks hide some of this, but your RabbitMQ wrapper should still be defensive.

---

## Retry Policy

```php
final class RetryPolicy
{
    public function __construct(
        private int $maxRetries = 3
    ) {}

    public function shouldRetry(MessageEnvelope $message): bool
    {
        return ($message->meta['retry_count'] ?? 0) < $this->maxRetries;
    }

    public function nextRetryCount(MessageEnvelope $message): int
    {
        return ((int) ($message->meta['retry_count'] ?? 0)) + 1;
    }
}
```

---

## Error Classes

```php
class TemporaryMessageException extends RuntimeException {}

class PermanentMessageException extends RuntimeException {}
```

Consumer behavior:

```text
TemporaryMessageException -> retry
PermanentMessageException -> DLQ
Throwable -> retry until limit, then DLQ
```

---

## Memory Limit Restart

```php
if (memory_get_usage(true) > 256 * 1024 * 1024) {
    $logger->warning('worker.memory_limit_reached');
    exit(0);
}
```

Supervisor restarts the worker.

---

## Supervisor Config

```ini
[program:php-rabbitmq-worker]
command=php /var/www/app/bin/console app:rabbitmq:consume product.search_index.queue
directory=/var/www/app
autostart=true
autorestart=true
stopwaitsecs=60
user=www-data
redirect_stderr=true
stdout_logfile=/var/log/php-rabbitmq-worker.log
```

---

## Graceful Shutdown

Use signal handling when available:

```php
pcntl_signal(SIGTERM, function () use (&$running) {
    $running = false;
});
```

The worker should finish the current message, then exit.

---

## PHP Production Checklist

- `composer install --no-dev --optimize-autoloader`
- durable exchange and queue
- persistent messages
- manual ack
- prefetch configured
- structured logs
- retry queue
- DLQ
- idempotency table
- Supervisor/systemd/Docker restart
- graceful shutdown
- memory limit restart

