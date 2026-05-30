# Chapter 17: Capstone Project

## Goal

Build a Laravel order event pipeline:

```text
POST /orders
  -> creates order in MySQL
  -> publishes orders.created
  -> payment consumer reads event
  -> creates payment row
  -> publishes payments.authorized
```

---

## Topics

```text
orders.created
payments.authorized
orders.created.dlq
```

---

## Steps

1. Install `mateusjunges/laravel-kafka` and `ext-rdkafka`.
2. Add Kafka config and topic names.
3. Create `OrderCreated` domain event.
4. Publish Kafka event after database commit.
5. Create `kafka:consume-orders` command.
6. Build `OrdersCreatedHandler`.
7. Add idempotency table using `event_id`.
8. Add tests with Kafka fake.
9. Add Supervisor or Docker process for the consumer.
10. Add logs for topic, partition, offset, key, and event ID.

---

## Acceptance Criteria

- Creating an order publishes one Kafka message.
- Payment consumer processes the message once even after duplicate delivery.
- Invalid messages go to DLQ.
- Tests assert topic and payload.
- Consumer can be restarted without losing events.
- Lag and handler failures are visible in logs or metrics.
