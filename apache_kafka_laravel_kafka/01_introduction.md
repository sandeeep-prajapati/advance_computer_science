# Chapter 1: Introduction

## What is Laravel Kafka?

Laravel Kafka usually refers to using Apache Kafka from a Laravel application through the `mateusjunges/laravel-kafka` package. The package wraps the lower-level PHP Kafka client experience in Laravel-friendly APIs for producing, consuming, testing, and configuring messages.

Kafka remains the broker. Laravel remains your application framework. The package is the bridge.

---

## Why Use Kafka with Laravel?

Laravel applications often start with synchronous workflows:

```text
Checkout request
  -> create order
  -> charge payment
  -> update inventory
  -> send email
  -> notify analytics
```

That works at small scale, but it couples unrelated systems to the request lifecycle. Kafka lets the checkout app publish an event once and allows many independent services to react.

```text
Checkout service
  -> Kafka topic: orders.created
      -> payment service
      -> inventory service
      -> email service
      -> analytics service
```

---

## Kafka vs Laravel Queues

| Feature | Laravel Queue | Kafka |
|---|---|---|
| Main purpose | Background jobs | Durable event stream |
| Consumption | One worker handles a job | Each consumer group can read every event |
| Replay | Not the default model | First-class offset-based replay |
| Ordering | Queue dependent | Guaranteed within a partition |
| Retention | Jobs are usually deleted | Events remain by retention policy |
| Best for | Commands and background work | Integration events and streaming data |

Use Laravel queues for internal work such as sending emails or resizing images. Use Kafka when other services need a durable stream of business events.

---

## Good Laravel Kafka Use Cases

- Publish `orders.created`, `payments.captured`, or `users.registered` events.
- Stream audit logs from many Laravel apps into one analytics system.
- Fan out domain events to independent services.
- Consume events from upstream systems and update Laravel read models.
- Build event-driven integrations without coupling services by HTTP calls.

---

## When Kafka is Too Much

Kafka adds operational and design complexity. Prefer simpler tools when:

- Only one Laravel app needs the data.
- A normal queue job is enough.
- You need request-response behavior.
- You cannot operate Kafka or use a managed Kafka provider.
- You do not need replay, partitioning, or many independent consumers.

---

## Mental Model

Think in events, not jobs:

```text
Job: "SendWelcomeEmail"
Event: "UserRegistered"
```

An event says something already happened. Consumers decide what to do with it.
