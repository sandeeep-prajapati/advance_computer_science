# Chapter 14: Deployment

## Running Consumers

Run consumers as managed processes, not from web requests.

Supervisor example:

```ini
[program:kafka-orders-consumer]
command=php /var/www/html/artisan kafka:consume-orders
directory=/var/www/html
autostart=true
autorestart=true
stopwaitsecs=30
stdout_logfile=/var/log/laravel/kafka-orders.log
stderr_logfile=/var/log/laravel/kafka-orders-error.log
```

---

## Artisan Command Wrapper

```php
final class ConsumeOrdersCommand extends Command
{
    protected $signature = 'kafka:consume-orders';

    public function handle(): int
    {
        Kafka::consumer(
                [config('kafka.topics.orders_created')],
                'payment-service',
                config('kafka.brokers')
            )
            ->withHandler(app(OrdersCreatedHandler::class))
            ->build()
            ->consume();

        return self::SUCCESS;
    }
}
```

---

## Graceful Shutdown

On deploy:

1. Stop accepting new work.
2. Let the current message finish.
3. Commit the offset if processing succeeded.
4. Close the consumer.
5. Start the new version.

Abrupt shutdown is acceptable only if handlers are idempotent and duplicate processing is safe.

---

## Kubernetes Notes

- Run consumers as Deployments.
- Use one container process per consumer command.
- Set termination grace period long enough for handlers.
- Expose lag metrics for autoscaling.
- Do not scale consumers beyond partition count unless you plan for idle replicas.
