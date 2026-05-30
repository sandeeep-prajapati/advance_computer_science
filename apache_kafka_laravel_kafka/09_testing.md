# Chapter 9: Testing

## Producer Tests

The package provides a Kafka fake for tests, so application tests do not need a real broker.

```php
use Junges\Kafka\Facades\Kafka;

public function test_it_publishes_order_created_event(): void
{
    Kafka::fake();

    $order = Order::factory()->create();

    app(OrderEventPublisher::class)->orderCreated($order);

    Kafka::assertPublishedOn(config('kafka.topics.orders_created'));
}
```

---

## Assert Payload Shape

```php
Kafka::assertPublishedOn(
    config('kafka.topics.orders_created'),
    function ($message) use ($order) {
        return $message->getBody()['order_id'] === $order->id;
    }
);
```

Exact assertion methods vary by package version, so keep tests close to the installed docs.

---

## Handler Unit Test

Test handler behavior without Kafka:

```php
public function test_handler_processes_order_created(): void
{
    $handler = app(OrdersCreatedHandler::class);

    $message = new FakeConsumerMessage([
        'event_id' => (string) str()->uuid(),
        'data' => ['order_id' => 123],
    ]);

    $handler($message, new FakeMessageConsumer());

    $this->assertDatabaseHas('payments', [
        'order_id' => 123,
    ]);
}
```

---

## Integration Tests

Use a real Kafka broker for a small number of integration tests:

- producer can publish to test topic
- consumer can read one message
- SASL/TLS configuration works
- serialization/deserialization works

Keep integration tests separate from normal unit tests because they are slower and require infrastructure.
