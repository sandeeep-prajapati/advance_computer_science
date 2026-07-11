# Chapter 1: PHP Package Architecture

## Composer Package Structure

Build your wrapper as a reusable Composer package or internal library.

```text
rabbitmq-wrapper/
├── composer.json
├── src/
│   ├── Contracts/
│   │   ├── MessagePublisher.php
│   │   ├── MessageConsumer.php
│   │   └── MessageHandler.php
│   ├── Envelope/
│   │   └── MessageEnvelope.php
│   ├── Serialization/
│   │   └── JsonMessageSerializer.php
│   ├── Registry/
│   │   └── HandlerRegistry.php
│   ├── RabbitMQ/
│   │   ├── RabbitConnectionFactory.php
│   │   ├── RabbitPublisher.php
│   │   ├── RabbitConsumer.php
│   │   └── Topology.php
│   └── Retry/
│       └── RetryPolicy.php
└── tests/
```

---

## composer.json

```json
{
  "name": "app/rabbitmq-wrapper",
  "type": "library",
  "require": {
    "php": "^8.2",
    "php-amqplib/php-amqplib": "^3.6",
    "psr/log": "^3.0"
  },
  "autoload": {
    "psr-4": {
      "App\\RabbitMQ\\": "src/"
    }
  },
  "autoload-dev": {
    "psr-4": {
      "App\\RabbitMQ\\Tests\\": "tests/"
    }
  }
}
```

Run:

```bash
composer install
composer dump-autoload
```

---

## Contracts

```php
namespace App\RabbitMQ\Contracts;

use App\RabbitMQ\Envelope\MessageEnvelope;

interface MessagePublisher
{
    public function publish(string $routingKey, array $payload, array $meta = []): void;
}

interface MessageHandler
{
    public function handle(MessageEnvelope $message): void;
}
```

Framework code depends on contracts. RabbitMQ implementation stays replaceable.

---

## Message Envelope

```php
namespace App\RabbitMQ\Envelope;

final readonly class MessageEnvelope
{
    public function __construct(
        public string $id,
        public string $type,
        public int $version,
        public string $occurredAt,
        public array $payload,
        public array $meta = []
    ) {}
}
```

The same envelope can be used in Core PHP, Laravel, Symfony, and Pimcore.

---

## Serializer

```php
namespace App\RabbitMQ\Serialization;

use App\RabbitMQ\Envelope\MessageEnvelope;

final class JsonMessageSerializer
{
    public function encode(MessageEnvelope $message): string
    {
        return json_encode([
            'id' => $message->id,
            'type' => $message->type,
            'version' => $message->version,
            'occurred_at' => $message->occurredAt,
            'payload' => $message->payload,
            'meta' => $message->meta,
        ], JSON_THROW_ON_ERROR);
    }

    public function decode(string $body): MessageEnvelope
    {
        $data = json_decode($body, true, flags: JSON_THROW_ON_ERROR);

        return new MessageEnvelope(
            $data['id'],
            $data['type'],
            (int) $data['version'],
            $data['occurred_at'],
            $data['payload'] ?? [],
            $data['meta'] ?? []
        );
    }
}
```

---

## Design Rule

Keep package core framework-neutral.

Framework-specific integrations should be adapters:

```text
Core package
├── envelope
├── publisher
├── consumer
└── handler contract

Laravel adapter
Symfony adapter
Pimcore adapter
```

