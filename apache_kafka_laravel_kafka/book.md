# Laravel Kafka — Complete Learning Guide

> A practical, hands-on guide to using Apache Kafka from Laravel with the `mateusjunges/laravel-kafka` package.

---

## Table of Contents

| Chapter | Title | Topics |
|---|---|---|
| 01 | [Introduction](01_introduction.md) | Laravel + Kafka use cases, event-driven apps, when Kafka fits |
| 02 | [Setup & Installation](02_setup_and_installation.md) | PHP extension, Composer package, Docker Kafka, env config |
| 03 | [Configuration](03_configuration.md) | Brokers, topics, consumer groups, config publishing, options |
| 04 | [Producing Messages](04_producing_messages.md) | Facade API, keys, headers, payloads, reliability options |
| 05 | [Consuming Messages](05_consuming_messages.md) | Consumers, handlers, groups, offsets, long-running workers |
| 06 | [Message Design](06_message_design.md) | Event names, payload contracts, correlation IDs, versioning |
| 07 | [Laravel Events, Jobs, and Outbox](07_events_jobs_outbox.md) | Domain events, queued publishing, transactional outbox pattern |
| 08 | [Retries and Dead Letter Topics](08_retries_and_dlq.md) | Backoff, poison messages, DLQ naming, replay strategy |
| 09 | [Testing](09_testing.md) | Kafka fake, assertions, handler tests, integration tests |
| 10 | [Security and Cloud Kafka](10_security_and_cloud.md) | SASL, TLS, Confluent Cloud, MSK, secrets |
| 11 | [Serialization and Schema Registry](11_serialization_schema_registry.md) | JSON, custom serializers, Avro/Protobuf concepts |
| 12 | [Observability](12_observability.md) | Logs, metrics, lag, traces, health checks |
| 13 | [Performance and Scaling](13_performance_scaling.md) | Batching, async publishing, partitions, Horizon/Supervisor |
| 14 | [Deployment](14_deployment.md) | Containers, graceful shutdown, process management |
| 15 | [Production Patterns](15_production_patterns.md) | Idempotency, ordering, sagas, CDC, multi-service contracts |
| 16 | [Troubleshooting](16_troubleshooting.md) | Common errors, extension issues, connectivity, offsets |
| 17 | [Capstone Project](17_capstone_project.md) | Order event pipeline in Laravel |
| 18 | [Learning Roadmap](18_learning_roadmap.md) | 4-week plan and practice checklist |

---

## Architecture

```text
Laravel App
  HTTP Controller / Command / Domain Service
        |
        | publish event
        v
  Junges\Kafka Producer  --->  Kafka Topic: orders.created
                                      |
                                      | consumed by group payment-service
                                      v
                              Laravel Kafka Consumer
                                      |
                                      v
                              Handler / Job / Application Service
```

---

## Recommended Reading Path

- Start with Chapters 1-5 if you are new to Kafka in Laravel.
- Use Chapters 7-9 before wiring Kafka into real business workflows.
- Use Chapters 10-16 before production.
- Build Chapter 17 after reading the guide once.

---

## References

- Laravel Kafka package: https://github.com/mateusjunges/laravel-kafka
- Laravel Kafka docs: https://laravelkafka.com
- Apache Kafka docs: https://kafka.apache.org/documentation/

---

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

---

# Chapter 2: Setup & Installation

## Requirements

Typical Laravel Kafka projects need:

- PHP and Laravel supported by your selected package version
- `ext-rdkafka`
- `librdkafka`
- Composer package `mateusjunges/laravel-kafka`
- A Kafka broker reachable from the Laravel runtime

Always check the package release notes for exact PHP and Laravel version support.

---

## Install System Dependencies

Ubuntu example:

```bash
sudo apt-get update
sudo apt-get install -y librdkafka-dev
sudo pecl install rdkafka
```

Enable the extension:

```ini
extension=rdkafka
```

Verify:

```bash
php -m | grep rdkafka
php --ri rdkafka
```

---

## Install the Laravel Package

```bash
composer require mateusjunges/laravel-kafka
```

Publish config if the package version provides a config file:

```bash
php artisan vendor:publish --provider="Junges\Kafka\KafkaServiceProvider"
```

---

## Local Kafka with Docker

Use a single-broker Kafka container for learning:

```yaml
services:
  kafka:
    image: bitnami/kafka:latest
    ports:
      - "9092:9092"
    environment:
      KAFKA_CFG_NODE_ID: 1
      KAFKA_CFG_PROCESS_ROLES: broker,controller
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: 1@kafka:9093
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE: "true"
```

For Laravel running inside another container, advertise an internal listener such as `kafka:9092`.

---

## Environment Variables

```dotenv
KAFKA_BROKERS=localhost:9092
KAFKA_CONSUMER_GROUP_ID=laravel-app
KAFKA_SECURITY_PROTOCOL=PLAINTEXT
```

For cloud Kafka, you will normally add SASL and TLS values.

---

## Smoke Test

Create a route or Artisan command that publishes one event:

```php
use Junges\Kafka\Facades\Kafka;
use Junges\Kafka\Message\Message;

Kafka::publish(config('kafka.brokers'))
    ->onTopic('dev.ping')
    ->withMessage(new Message(
        headers: ['source' => 'laravel'],
        body: ['message' => 'hello kafka']
    ))
    ->send();
```

If this fails, debug in this order:

1. Is `rdkafka` loaded by PHP CLI and FPM?
2. Can the Laravel container reach the broker host and port?
3. Is Kafka advertising a hostname Laravel can resolve?
4. Does your cloud provider require SASL/TLS?

---

# Chapter 3: Configuration

## Core Settings

Most Laravel Kafka apps need a small set of shared settings:

```php
return [
    'brokers' => env('KAFKA_BROKERS', 'localhost:9092'),

    'consumer_group_id' => env('KAFKA_CONSUMER_GROUP_ID', 'laravel-app'),

    'security' => [
        'protocol' => env('KAFKA_SECURITY_PROTOCOL', 'PLAINTEXT'),
        'sasl_mechanism' => env('KAFKA_SASL_MECHANISM'),
        'sasl_username' => env('KAFKA_SASL_USERNAME'),
        'sasl_password' => env('KAFKA_SASL_PASSWORD'),
    ],
];
```

Keep topic names in config too:

```php
'topics' => [
    'orders_created' => env('KAFKA_TOPIC_ORDERS_CREATED', 'orders.created'),
    'orders_failed' => env('KAFKA_TOPIC_ORDERS_FAILED', 'orders.failed'),
],
```

---

## Producer Options

Producer options are passed to the underlying Kafka client:

```php
Kafka::publish(config('kafka.brokers'))
    ->onTopic(config('kafka.topics.orders_created'))
    ->withConfigOptions([
        'acks' => 'all',
        'enable.idempotence' => 'true',
        'compression.type' => 'lz4',
        'message.timeout.ms' => '30000',
    ]);
```

Recommended defaults for important business events:

| Option | Suggested value | Why |
|---|---:|---|
| `acks` | `all` | Wait for in-sync replicas |
| `enable.idempotence` | `true` | Avoid duplicates caused by producer retry |
| `compression.type` | `lz4` or `zstd` | Reduce network and broker storage |
| `message.timeout.ms` | `30000` | Fail clearly instead of hanging forever |

---

## Consumer Options

```php
$consumer = Kafka::consumer(
        [config('kafka.topics.orders_created')],
        config('kafka.consumer_group_id'),
        config('kafka.brokers')
    )
    ->withOptions([
        'auto.offset.reset' => 'earliest',
        'enable.auto.commit' => 'false',
        'max.poll.interval.ms' => '300000',
    ]);
```

Common settings:

| Option | Purpose |
|---|---|
| `group.id` | Identifies the independent consumer group |
| `auto.offset.reset` | Where to start when no committed offset exists |
| `enable.auto.commit` | Whether offsets are committed automatically |
| `max.poll.interval.ms` | Max processing time before the broker considers consumer stuck |
| `session.timeout.ms` | Failure detection window |

---

## Naming Conventions

Use topic names that communicate ownership and event type:

```text
orders.created
orders.cancelled
payments.captured
inventory.stock-reserved
```

For environments:

```text
local.orders.created
staging.orders.created
prod.orders.created
```

Avoid putting implementation details into topic names. `orders.created` is better than `laravel_order_controller_created`.

---

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

---

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

---

# Chapter 6: Message Design

## Event Envelope

Use a consistent envelope:

```json
{
  "event_id": "01J...",
  "event_type": "orders.created",
  "event_version": 1,
  "occurred_at": "2026-05-30T10:00:00+05:30",
  "producer": "checkout-api",
  "data": {
    "order_id": 123,
    "user_id": 50,
    "total": 1999,
    "currency": "INR"
  }
}
```

This makes consumers easier to build and debug.

---

## Required Metadata

| Field | Why |
|---|---|
| `event_id` | Idempotency and tracing |
| `event_type` | Dispatching and logging |
| `event_version` | Schema evolution |
| `occurred_at` | Business time |
| `producer` | Ownership |
| `correlation_id` | Request tracing across services |

---

## Compatibility Rules

Safe changes:

- Add optional fields.
- Add new event types.
- Add nullable fields with defaults.

Dangerous changes:

- Rename fields.
- Change numeric strings to numbers without coordination.
- Remove fields that consumers may still read.
- Reuse an event name for a different meaning.

---

## Versioning

Prefer versioning in the payload or header:

```php
headers: [
    'event' => 'orders.created',
    'event_version' => '2',
]
```

Avoid creating `orders.created.v1`, `orders.created.v2`, and `orders.created.v3` topics unless the old and new streams truly need different retention, permissions, or consumers.

---

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

---

# Chapter 8: Retries and Dead Letter Topics

## Failure Types

| Failure | Example | Retry? |
|---|---|---|
| Transient | HTTP timeout, DB deadlock | Yes |
| Recoverable later | Missing dependent record | Yes, delayed |
| Permanent | Invalid payload shape | No, send to DLQ |
| Bug | Handler throws every time | DLQ after limited attempts |

---

## Idempotent Consumer First

Retries create duplicates. Make handlers safe:

```php
ProcessedKafkaMessage::firstOrCreate([
    'event_id' => $payload['event_id'],
]);
```

Or enforce natural uniqueness:

```php
Payment::updateOrCreate(
    ['order_id' => $payload['data']['order_id']],
    ['status' => 'authorized']
);
```

---

## Dead Letter Topic Naming

```text
orders.created.dlq
payments.captured.dlq
```

DLQ payload should include:

- original topic
- original partition and offset
- original key
- original headers
- original body
- exception class and message
- failed at timestamp

---

## Replay Strategy

Do not replay DLQ blindly. Use a command:

```bash
php artisan kafka:replay-dlq orders.created.dlq --limit=100
```

Before replay:

1. Fix the bug or data issue.
2. Pick a small batch.
3. Monitor consumer lag and error logs.
4. Keep the original DLQ event for audit if required.

---

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

---

# Chapter 10: Security and Cloud Kafka

## Local vs Production

Local Kafka often uses plaintext:

```dotenv
KAFKA_SECURITY_PROTOCOL=PLAINTEXT
```

Managed Kafka usually uses SASL and TLS:

```dotenv
KAFKA_SECURITY_PROTOCOL=SASL_SSL
KAFKA_SASL_MECHANISM=PLAIN
KAFKA_SASL_USERNAME=...
KAFKA_SASL_PASSWORD=...
```

---

## SASL Options

```php
->withConfigOptions([
    'security.protocol' => env('KAFKA_SECURITY_PROTOCOL', 'SASL_SSL'),
    'sasl.mechanisms' => env('KAFKA_SASL_MECHANISM', 'PLAIN'),
    'sasl.username' => env('KAFKA_SASL_USERNAME'),
    'sasl.password' => env('KAFKA_SASL_PASSWORD'),
])
```

Use the same security options for producers and consumers.

---

## Secret Handling

- Store credentials in environment variables or a secret manager.
- Do not commit `.env` files.
- Rotate cloud Kafka API keys.
- Give each service its own Kafka principal.
- Restrict ACLs by topic and consumer group.

---

## ACL Shape

For a producer service:

```text
WRITE on topic orders.created
DESCRIBE on topic orders.created
```

For a consumer service:

```text
READ on topic orders.created
DESCRIBE on topic orders.created
READ on group payment-service
```

Least privilege makes mistakes much easier to contain.

---

# Chapter 11: Serialization and Schema Registry

## JSON First

JSON is a good first choice for Laravel teams:

```php
body: [
    'event_id' => (string) str()->uuid(),
    'event_type' => 'orders.created',
    'event_version' => 1,
    'data' => [
        'order_id' => $order->id,
    ],
]
```

Benefits:

- Easy to inspect
- Native PHP arrays map naturally
- No schema infrastructure required

Costs:

- Runtime validation is your responsibility
- Type drift can happen silently
- Payloads are larger than binary formats

---

## Custom Serializers

Use custom serializers when:

- every message must share a standard envelope
- dates and money need consistent formatting
- you need encryption or compression at the application layer
- you are integrating with Avro, Protobuf, or JSON Schema

---

## Schema Registry

Schema Registry stores schemas and enforces compatibility. It is common with Avro, Protobuf, and JSON Schema.

Compatibility modes:

| Mode | Meaning |
|---|---|
| Backward | New consumers can read old data |
| Forward | Old consumers can read new data |
| Full | Both backward and forward |
| None | No compatibility checks |

For multi-service teams, schema compatibility prevents accidental breaking changes.

---

## Laravel Validation at Boundaries

Validate consumed payloads before business logic:

```php
Validator::make($payload, [
    'event_id' => ['required', 'uuid'],
    'event_type' => ['required', 'string'],
    'event_version' => ['required', 'integer'],
    'data.order_id' => ['required', 'integer'],
])->validate();
```

Invalid messages should go to a DLQ with enough context to debug.

---

# Chapter 12: Observability

## What to Log

Log structured context:

```php
Log::info('Kafka message consumed', [
    'topic' => $message->getTopicName(),
    'partition' => $message->getPartition(),
    'offset' => $message->getOffset(),
    'key' => $message->getKey(),
    'event_id' => data_get($message->getBody(), 'event_id'),
    'event_type' => data_get($message->getBody(), 'event_type'),
]);
```

Avoid logging full payloads if they contain personal data.

---

## Metrics

Track:

- messages produced per minute
- produce failures
- messages consumed per minute
- handler duration
- handler failures
- consumer lag
- DLQ count
- retry count

Consumer lag is the most important operational signal. Growing lag means consumers are slower than producers or are failing.

---

## Tracing

Propagate correlation IDs:

```php
headers: [
    'correlation_id' => request()->header('X-Correlation-ID', (string) str()->uuid()),
]
```

Each consumer should include the same ID in logs and outbound calls.

---

## Health Checks

A Laravel HTTP health check should not only say the web app is alive. Add separate checks for:

- broker connectivity
- producer publish ability to a health topic
- consumer process running
- consumer lag below threshold

Process health and stream health are different things.

---

# Chapter 13: Performance and Scaling

## Producer Performance

For low throughput, the normal `send()` API is simple. For high throughput, avoid flushing after every message and use the package's async or batch publishing support when available in your installed version.

Useful producer options:

```php
->withConfigOptions([
    'compression.type' => 'lz4',
    'linger.ms' => '10',
    'batch.num.messages' => '1000',
    'queue.buffering.max.messages' => '100000',
])
```

---

## Consumer Scaling

Consumer parallelism is bounded by partitions.

```text
Topic has 6 partitions
Consumer group has 3 instances
Each instance receives about 2 partitions
```

If a topic has 3 partitions, running 10 consumers in the same group will still process with only 3 active consumers.

---

## Handler Performance

Keep handlers lean:

- validate payload
- call application service
- persist result
- commit offset

Move slow side effects to Laravel jobs when the processing contract allows it.

---

## Backpressure

If lag grows:

1. Check handler error rate.
2. Check database or HTTP dependency latency.
3. Increase partitions if ordering allows.
4. Add consumer instances.
5. Batch work where safe.
6. Split heavy consumers into separate topics or services.

---

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

---

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

---

# Chapter 16: Troubleshooting

## `Class RdKafka\... not found`

The PHP extension is not loaded.

```bash
php -m | grep rdkafka
php --ri rdkafka
```

Remember PHP CLI and PHP-FPM can use different ini files.

---

## Cannot Connect to Broker

Common causes:

- wrong `KAFKA_BROKERS`
- Docker container using `localhost` instead of service name
- Kafka `advertised.listeners` points to a host Laravel cannot resolve
- firewall or security group blocks port
- cloud Kafka requires SASL/TLS

Inside the Laravel container, test DNS and TCP connectivity.

---

## Consumer Reads Nothing

Check:

- correct topic name
- topic has messages
- correct consumer group
- committed offset is already at end
- `auto.offset.reset` only applies when there is no committed offset
- handler is actually attached before `consume()`

Use a new temporary group ID to test from the beginning.

---

## Duplicates

Duplicates are normal after crashes and retries. Fix with:

- idempotent handlers
- unique event IDs
- database constraints
- offset commit after successful processing

---

## Lag Keeps Growing

Possible causes:

- handler too slow
- downstream database or API slow
- not enough partitions
- not enough consumers
- poison message repeatedly failing
- consumer process restarting

Start with logs, lag by partition, and handler duration metrics.

---

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

---

# Chapter 18: Learning Roadmap

## Week 1: Kafka and Laravel Basics

- Understand topics, partitions, offsets, and consumer groups.
- Run Kafka locally with Docker.
- Install `ext-rdkafka` and the Laravel package.
- Publish a message from an Artisan command.
- Consume a message with a simple handler.

---

## Week 2: Application Design

- Create a domain event.
- Publish after database commit.
- Add a standard event envelope.
- Add correlation IDs.
- Test producers with Kafka fake.

---

## Week 3: Reliability

- Implement idempotent consumers.
- Add retry strategy.
- Add DLQ topic.
- Learn offset behavior.
- Practice replaying a small DLQ batch.

---

## Week 4: Production Readiness

- Configure SASL/TLS for managed Kafka.
- Run consumers under Supervisor, Docker, or Kubernetes.
- Add structured logs and metrics.
- Monitor consumer lag.
- Document topic ownership and schemas.

---

## Practice Checklist

- [ ] Produce `orders.created`
- [ ] Consume with a named group
- [ ] Use stable Kafka keys
- [ ] Add event IDs
- [ ] Make handler idempotent
- [ ] Test with Kafka fake
- [ ] Run a consumer as a managed process
- [ ] Add DLQ handling
- [ ] Configure SASL/TLS
- [ ] Monitor lag
