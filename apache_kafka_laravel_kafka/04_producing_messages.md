# Chapter 4: Producing Messages

## Basic Producer

```php
use Junges\Kafka\Facades\Kafka;
use Junges\Kafka\Message\Message;

$message = new Message(
    headers: [
        'event' => 'orders.created',
        'correlation_id' => (string) str()->uuid(),
    ],
    body: [
        'order_id' => 123,
        'user_id' => 50,
        'total' => 1999,
        'currency' => 'INR',
    ],
    key: 'order-123'
);

Kafka::publish(config('kafka.brokers'))
    ->onTopic(config('kafka.topics.orders_created'))
    ->withMessage($message)
    ->send();
```

The official package docs show the same producer shape: create a producer with `Kafka::publish('broker')`, select a topic with `onTopic(...)`, add options such as key and headers, then call `send()`.

---

## Choosing a Kafka Key

The key controls partitioning. Events with the same key usually go to the same partition, preserving order for that key.

Good keys:

| Event | Key |
|---|---|
| `orders.created` | `order-{id}` |
| `payments.captured` | `payment-{id}` or `order-{id}` |
| `users.updated` | `user-{id}` |

Use `order_id` as the key when all order lifecycle events must stay ordered.

---

## Headers

Headers are useful for metadata that should not live in the business payload:

```php
[
    'event' => 'orders.created',
    'event_version' => '1',
    'correlation_id' => $request->header('X-Correlation-ID', (string) str()->uuid()),
    'producer' => config('app.name'),
    'occurred_at' => now()->toIso8601String(),
]
```

---

## Publishing from an Application Service

```php
final class OrderEventPublisher
{
    public function orderCreated(Order $order): void
    {
        Kafka::publish(config('kafka.brokers'))
            ->onTopic(config('kafka.topics.orders_created'))
            ->withMessage(new Message(
                headers: [
                    'event' => 'orders.created',
                    'event_version' => '1',
                    'correlation_id' => request()->header('X-Correlation-ID', (string) str()->uuid()),
                ],
                body: [
                    'order_id' => $order->id,
                    'user_id' => $order->user_id,
                    'total' => $order->total,
                    'currency' => $order->currency,
                    'created_at' => $order->created_at->toIso8601String(),
                ],
                key: "order-{$order->id}",
            ))
            ->send();
    }
}
```

---

## Avoid Publishing Inside Open Transactions

This is risky:

```php
DB::transaction(function () use ($order) {
    $order->save();
    Kafka::publish(...)->send();
});
```

If Kafka succeeds and the database rolls back, consumers see an event for data that does not exist. Prefer:

- publish after commit for simple cases
- use Laravel's `afterCommit`
- use the transactional outbox pattern for important workflows
