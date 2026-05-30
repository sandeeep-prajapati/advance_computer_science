# Chapter 15: Production Patterns

## Idempotency

Kafka provides at-least-once consumption in common Laravel setups. Expect duplicates.

Use an event log table:

```php
Schema::create('processed_kafka_messages', function (Blueprint $table) {
    $table->id();
    $table->uuid('event_id')->unique();
    $table->string('consumer_group');
    $table->timestamp('processed_at');
});
```

Before processing, check whether the event was already handled.

---

## Ordering

Kafka ordering is per partition. If order matters, use a stable key:

```php
key: "order-{$order->id}"
```

All events for that order should use the same key.

---

## Sagas

For workflows across services:

```text
orders.created
  -> payment authorized
  -> inventory reserved
  -> shipment requested
```

Each service owns its local transaction and publishes the next event. Compensating events handle failure:

```text
inventory.reserve_failed
  -> payment refund requested
```

---

## Contract Ownership

Every topic needs an owner:

- Who can publish?
- Who maintains the schema?
- What is the retention policy?
- What is the compatibility rule?
- Who responds when consumers break?

Kafka works best when event contracts are treated as product interfaces.
