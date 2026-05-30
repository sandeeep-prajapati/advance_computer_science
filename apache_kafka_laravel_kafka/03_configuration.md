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
