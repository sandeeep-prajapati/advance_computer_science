# Chapter 7: Laravel Events, Jobs, and Outbox

## Domain Event to Kafka

Keep your business event separate from Kafka infrastructure:

```php
final class OrderCreated
{
    public function __construct(public readonly int $orderId) {}
}
```

Listener:

```php
final class PublishOrderCreatedToKafka
{
    public bool $afterCommit = true;

    public function handle(OrderCreated $event): void
    {
        $order = Order::findOrFail($event->orderId);

        app(OrderEventPublisher::class)->orderCreated($order);
    }
}
```

`$afterCommit = true` prevents publishing if the database transaction rolls back.

---

## Queue the Publish

For request latency, publish from a queued listener:

```php
final class PublishOrderCreatedToKafka implements ShouldQueue
{
    public bool $afterCommit = true;

    public function handle(OrderCreated $event): void
    {
        app(OrderEventPublisher::class)->orderCreated(
            Order::findOrFail($event->orderId)
        );
    }
}
```

This improves user-facing latency, but it does not fully solve atomicity.

---

## Transactional Outbox Pattern

For high-value workflows:

1. Save business data and an outbox row in the same database transaction.
2. A separate worker reads unpublished outbox rows.
3. The worker publishes to Kafka.
4. The worker marks the outbox row as published.

Schema:

```php
Schema::create('outbox_messages', function (Blueprint $table) {
    $table->id();
    $table->uuid('event_id')->unique();
    $table->string('topic');
    $table->string('key')->nullable();
    $table->json('headers');
    $table->json('payload');
    $table->timestamp('published_at')->nullable();
    $table->timestamps();
});
```

The outbox is the most reliable bridge between a relational database transaction and Kafka.
