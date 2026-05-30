# Chapter 5: Consuming Messages

## Basic Consumer

The package docs show consumers being created with `Kafka::consumer()`, optionally passing topics, group ID, and broker.

```php
use Junges\Kafka\Facades\Kafka;
use Junges\Kafka\Contracts\ConsumerMessage;
use Junges\Kafka\Contracts\MessageConsumer;

Kafka::consumer(
        [config('kafka.topics.orders_created')],
        'payment-service',
        config('kafka.brokers')
    )
    ->withHandler(function (ConsumerMessage $message, MessageConsumer $consumer): void {
        $body = $message->getBody();

        app(ProcessCreatedOrder::class)->handle(
            orderId: $body['order_id'],
            userId: $body['user_id'],
            total: $body['total'],
        );
    })
    ->build()
    ->consume();
```

---

## Consumer Groups

Every independent service should have its own group:

```text
Topic: orders.created

Group: payment-service
  reads every order event for payment authorization

Group: email-service
  reads every order event for email notifications

Group: analytics-service
  reads every order event for reporting
```

Adding more instances to the same group scales processing across partitions.

---

## Handler Classes

Prefer handler classes once logic grows:

```php
use Junges\Kafka\Contracts\ConsumerMessage;
use Junges\Kafka\Contracts\Handler;
use Junges\Kafka\Contracts\MessageConsumer;

final class OrdersCreatedHandler implements Handler
{
    public function __invoke(ConsumerMessage $message, MessageConsumer $consumer): void
    {
        $payload = $message->getBody();

        app(ProcessCreatedOrder::class)->handle($payload);
    }
}
```

Then wire it:

```php
Kafka::consumer([...], 'payment-service', config('kafka.brokers'))
    ->withHandler(new OrdersCreatedHandler())
    ->build()
    ->consume();
```

---

## Long-Running Consumers

Kafka consumers are long-running processes. That means:

- Do not assume fresh config or code after deploy until process restart.
- Avoid memory leaks in handlers.
- Close database connections if they become stale.
- Use Supervisor, systemd, Docker, or Kubernetes to restart failed consumers.
- Handle SIGTERM for graceful deploys.

---

## Offset Rule

Commit an offset only after processing succeeds.

```text
process message
  -> database update succeeds
  -> external call succeeds
  -> commit offset
```

If the process crashes before commit, Kafka can redeliver the message. Your handler must be idempotent.
