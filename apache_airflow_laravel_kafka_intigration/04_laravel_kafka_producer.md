# Lesson 4: Laravel Kafka Producer

## Install Laravel Kafka Package

```bash
composer require mateusjunges/laravel-kafka
```

Make sure PHP has `rdkafka` installed and enabled.

---

## Publish an Event

```php
use Junges\Kafka\Facades\Kafka;
use Junges\Kafka\Message\Message;

Kafka::publish(config('kafka.brokers'))
    ->onTopic('orders.created')
    ->withMessage(new Message(
        headers: [
            'event' => 'orders.created',
            'event_version' => '1',
            'correlation_id' => (string) str()->uuid(),
        ],
        body: [
            'order_id' => $order->id,
            'user_id' => $order->user_id,
            'total' => $order->total,
            'created_at' => $order->created_at->toIso8601String(),
        ],
        key: "order-{$order->id}",
    ))
    ->send();
```

---

## Publish After Commit

If the event depends on database data, publish after the database transaction commits.

```php
final class PublishOrderCreatedToKafka implements ShouldQueue
{
    public bool $afterCommit = true;

    public function handle(OrderCreated $event): void
    {
        app(OrderKafkaPublisher::class)->publishCreated($event->orderId);
    }
}
```

For critical workflows, use an outbox table and a dedicated publisher worker.

---

## Topic Design

Good topic names:

```text
orders.created
orders.cancelled
reports.generated
airflow.pipeline.failed
```

Use stable keys when order matters:

```text
order-123
report-2026-05-30
customer-50
```
