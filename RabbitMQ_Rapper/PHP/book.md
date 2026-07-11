# PHP RabbitMQ Wrapper: Complete Book

> Detailed PHP RabbitMQ wrapper notes for Core PHP, Laravel, Symfony, and Pimcore.

---

# PHP RabbitMQ Wrapper Book

Detailed notes for building RabbitMQ wrappers in PHP.

This folder covers:

- Core PHP with `php-amqplib/php-amqplib`
- Laravel service provider + Artisan worker
- Symfony bundle + console command
- Pimcore bundle/event integration
- Composer package design
- Retry, DLQ, idempotency, and deployment for PHP workers

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_php_package_architecture.md) | **PHP Package Architecture** | Composer, PSR-4, contracts, package layout |
| [02](02_core_php_consumer_publisher.md) | **Core PHP Consumer & Publisher** | php-amqplib, connection, publisher, consumer, CLI worker |
| [03](03_laravel_integration.md) | **Laravel Integration** | config, service provider, Artisan command, handlers through container |
| [04](04_symfony_integration.md) | **Symfony Integration** | services.yaml, bundle, console command, tagged handlers |
| [05](05_pimcore_integration.md) | **Pimcore Integration** | Pimcore events, DataObject publishing, bundle design |
| [06](06_php_retries_dlq_operations.md) | **PHP Retries, DLQ & Operations** | retry policy, DLQ, Supervisor, memory limits, graceful restart |
| [book](book.md) | **Combined PHP Book** | All PHP chapters in one file |

---

## PHP Wrapper Goal

Your PHP application should use RabbitMQ like this:

```php
$publisher->publish('product.updated', [
    'product_id' => 501,
]);
```

And your worker should look like this:

```bash
php bin/consume product.search_index.queue
php artisan rabbitmq:consume product.search_index.queue
php bin/console app:rabbitmq:consume product.search_index.queue
```

The business handler should not know about AMQP channels, delivery tags, or queue declarations.


---

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


---

# Chapter 2: Core PHP Consumer & Publisher

## Connection Factory

```php
namespace App\RabbitMQ\RabbitMQ;

use PhpAmqpLib\Connection\AMQPStreamConnection;

final class RabbitConnectionFactory
{
    public function __construct(private array $config) {}

    public function connect(): AMQPStreamConnection
    {
        return new AMQPStreamConnection(
            $this->config['host'],
            $this->config['port'],
            $this->config['user'],
            $this->config['password'],
            $this->config['vhost'] ?? '/'
        );
    }
}
```

---

## Topology

Declare exchanges, queues, and bindings before consuming.

```php
use PhpAmqpLib\Channel\AMQPChannel;

final class Topology
{
    public function declare(AMQPChannel $channel): void
    {
        $channel->exchange_declare('app.events', 'topic', false, true, false);
        $channel->exchange_declare('app.dlx', 'topic', false, true, false);

        $channel->queue_declare(
            'product.search_index.queue',
            false,
            true,
            false,
            false,
            false,
            ['x-dead-letter-exchange' => ['S', 'app.dlx']]
        );

        $channel->queue_bind('product.search_index.queue', 'app.events', 'product.*');

        $channel->queue_declare('product.search_index.dlq', false, true, false, false);
        $channel->queue_bind('product.search_index.dlq', 'app.dlx', '#');
    }
}
```

---

## Publisher

```php
use App\RabbitMQ\Envelope\MessageEnvelope;
use App\RabbitMQ\Serialization\JsonMessageSerializer;
use PhpAmqpLib\Channel\AMQPChannel;
use PhpAmqpLib\Message\AMQPMessage;

final class RabbitPublisher
{
    public function __construct(
        private AMQPChannel $channel,
        private JsonMessageSerializer $serializer,
        private string $exchange = 'app.events'
    ) {}

    public function publish(string $routingKey, array $payload, array $meta = []): void
    {
        $message = new MessageEnvelope(
            id: bin2hex(random_bytes(16)),
            type: $routingKey,
            version: 1,
            occurredAt: gmdate('c'),
            payload: $payload,
            meta: $meta
        );

        $body = $this->serializer->encode($message);

        $amqpMessage = new AMQPMessage($body, [
            'content_type' => 'application/json',
            'delivery_mode' => AMQPMessage::DELIVERY_MODE_PERSISTENT,
            'message_id' => $message->id,
            'correlation_id' => $meta['correlation_id'] ?? null,
        ]);

        $this->channel->basic_publish($amqpMessage, $this->exchange, $routingKey);
    }
}
```

---

## Handler Registry

```php
use App\RabbitMQ\Contracts\MessageHandler;

final class HandlerRegistry
{
    public function __construct(private array $handlers) {}

    public function get(string $type): MessageHandler
    {
        return $this->handlers[$type]
            ?? throw new RuntimeException("No handler registered for {$type}");
    }
}
```

---

## Consumer

```php
use PhpAmqpLib\Channel\AMQPChannel;
use PhpAmqpLib\Message\AMQPMessage;
use Psr\Log\LoggerInterface;

final class RabbitConsumer
{
    public function __construct(
        private AMQPChannel $channel,
        private JsonMessageSerializer $serializer,
        private HandlerRegistry $registry,
        private LoggerInterface $logger
    ) {}

    public function consume(string $queue): void
    {
        $this->channel->basic_qos(null, 1, null);

        $this->channel->basic_consume(
            $queue,
            '',
            false,
            false,
            false,
            false,
            fn (AMQPMessage $message) => $this->process($message, $queue)
        );

        while ($this->channel->is_consuming()) {
            $this->channel->wait();
        }
    }

    private function process(AMQPMessage $raw, string $queue): void
    {
        try {
            $message = $this->serializer->decode($raw->getBody());
            $handler = $this->registry->get($message->type);
            $handler->handle($message);

            $raw->ack();

            $this->logger->info('rabbitmq.message_processed', [
                'message_id' => $message->id,
                'type' => $message->type,
                'queue' => $queue,
            ]);
        } catch (Throwable $exception) {
            $this->logger->error('rabbitmq.message_failed', [
                'queue' => $queue,
                'error' => $exception->getMessage(),
            ]);

            $raw->nack(false, false);
        }
    }
}
```

---

## CLI Worker

```php
require dirname(__DIR__) . '/vendor/autoload.php';

$config = [
    'host' => '127.0.0.1',
    'port' => 5672,
    'user' => 'guest',
    'password' => 'guest',
];

$connection = (new RabbitConnectionFactory($config))->connect();
$channel = $connection->channel();

(new Topology())->declare($channel);

$consumer = new RabbitConsumer(
    $channel,
    new JsonMessageSerializer(),
    new HandlerRegistry([
        'product.updated' => new ProductUpdatedHandler(),
    ]),
    new SimpleLogger()
);

$consumer->consume($argv[1] ?? 'product.search_index.queue');
```

Run:

```bash
php bin/consume.php product.search_index.queue
```


---

# Chapter 3: Laravel Integration

## Laravel Architecture

```text
config/rabbitmq.php
└── RabbitMqServiceProvider
    ├── binds AMQP connection
    ├── binds publisher
    ├── binds consumer
    ├── registers handler map
    └── exposes Artisan command
```

---

## Config

```php
return [
    'host' => env('RABBITMQ_HOST', '127.0.0.1'),
    'port' => (int) env('RABBITMQ_PORT', 5672),
    'user' => env('RABBITMQ_USER', 'guest'),
    'password' => env('RABBITMQ_PASSWORD', 'guest'),
    'vhost' => env('RABBITMQ_VHOST', '/'),
    'exchange' => env('RABBITMQ_EXCHANGE', 'app.events'),
    'handlers' => [
        'product.updated' => App\Messaging\Handlers\ProductUpdatedHandler::class,
    ],
];
```

---

## Service Provider

```php
namespace App\Providers;

use Illuminate\Support\ServiceProvider;
use PhpAmqpLib\Connection\AMQPStreamConnection;

final class RabbitMqServiceProvider extends ServiceProvider
{
    public function register(): void
    {
        $this->app->singleton(AMQPStreamConnection::class, function () {
            return new AMQPStreamConnection(
                config('rabbitmq.host'),
                config('rabbitmq.port'),
                config('rabbitmq.user'),
                config('rabbitmq.password'),
                config('rabbitmq.vhost')
            );
        });

        $this->app->singleton(RabbitPublisher::class);
        $this->app->singleton(RabbitConsumer::class);
        $this->app->singleton(HandlerRegistry::class);
    }

    public function boot(): void
    {
        $this->commands([
            ConsumeRabbitMqCommand::class,
        ]);
    }
}
```

---

## Laravel Handler Registry

```php
use Illuminate\Contracts\Container\Container;

final class HandlerRegistry
{
    public function __construct(private Container $container) {}

    public function get(string $type): MessageHandler
    {
        $class = config("rabbitmq.handlers.{$type}");

        if (! $class) {
            throw new RuntimeException("No RabbitMQ handler for {$type}");
        }

        return $this->container->make($class);
    }
}
```

---

## Handler

```php
final class ProductUpdatedHandler implements MessageHandler
{
    public function __construct(private ProductIndexer $indexer) {}

    public function handle(MessageEnvelope $message): void
    {
        $productId = (int) $message->payload['product_id'];

        $this->indexer->index($productId);
    }
}
```

Laravel injects `ProductIndexer` automatically.

---

## Artisan Command

```php
use Illuminate\Console\Command;

final class ConsumeRabbitMqCommand extends Command
{
    protected $signature = 'rabbitmq:consume {queue}';

    protected $description = 'Consume RabbitMQ messages';

    public function handle(RabbitConsumer $consumer): int
    {
        $consumer->consume($this->argument('queue'));

        return self::SUCCESS;
    }
}
```

Run:

```bash
php artisan rabbitmq:consume product.search_index.queue
```

---

## Laravel Publishing Example

```php
final class ProductController
{
    public function update(Request $request, RabbitPublisher $publisher): JsonResponse
    {
        $productId = 501;

        $publisher->publish('product.updated', [
            'product_id' => $productId,
        ], [
            'correlation_id' => $request->headers->get('X-Correlation-ID'),
            'source' => 'laravel',
        ]);

        return response()->json(['status' => 'queued'], 202);
    }
}
```

---

## Laravel Worker Deployment

Use Supervisor:

```ini
[program:laravel-rabbitmq-product-index]
command=php /var/www/app/artisan rabbitmq:consume product.search_index.queue
directory=/var/www/app
autostart=true
autorestart=true
user=www-data
redirect_stderr=true
stdout_logfile=/var/log/rabbitmq-product-index.log
```


---

# Chapter 4: Symfony Integration

## Symfony Architecture

```text
RabbitMqBundle
├── DependencyInjection
├── services.yaml
├── Console Command
├── tagged handlers
└── publisher service
```

Symfony is strong for this because it has autowiring, autoconfiguration, service tags, and console commands.

---

## services.yaml

```yaml
services:
  App\RabbitMQ\:
    resource: '../src/RabbitMQ/'
    autowire: true
    autoconfigure: true

  PhpAmqpLib\Connection\AMQPStreamConnection:
    arguments:
      $host: '%env(RABBITMQ_HOST)%'
      $port: '%env(int:RABBITMQ_PORT)%'
      $user: '%env(RABBITMQ_USER)%'
      $password: '%env(RABBITMQ_PASSWORD)%'
      $vhost: '%env(default::RABBITMQ_VHOST)%'
```

---

## Tagged Handler Attribute

```php
#[Attribute(Attribute::TARGET_CLASS)]
final class AsRabbitHandler
{
    public function __construct(public string $type) {}
}
```

Handler:

```php
#[AsRabbitHandler('product.updated')]
final class ProductUpdatedHandler implements MessageHandler
{
    public function __construct(private ProductIndexer $indexer) {}

    public function handle(MessageEnvelope $message): void
    {
        $this->indexer->index((int) $message->payload['product_id']);
    }
}
```

---

## Console Command

```php
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(name: 'app:rabbitmq:consume')]
final class ConsumeRabbitMqCommand extends Command
{
    public function __construct(private RabbitConsumer $consumer)
    {
        parent::__construct();
    }

    protected function configure(): void
    {
        $this->addArgument('queue', InputArgument::REQUIRED);
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $this->consumer->consume($input->getArgument('queue'));

        return Command::SUCCESS;
    }
}
```

Run:

```bash
php bin/console app:rabbitmq:consume product.search_index.queue
```

---

## Symfony Publishing Example

```php
final class ProductController
{
    public function update(RabbitPublisher $publisher): JsonResponse
    {
        $publisher->publish('product.updated', [
            'product_id' => 501,
        ], [
            'source' => 'symfony',
        ]);

        return new JsonResponse(['status' => 'queued'], 202);
    }
}
```

---

## Symfony Best Practices

- Use service tags for handlers.
- Use console commands for workers.
- Use Monolog for structured logs.
- Keep message handlers as services.
- Run workers through Supervisor/systemd/Docker.


---

# Chapter 5: Pimcore Integration

## Pimcore Mental Model

Pimcore is Symfony-based. Use Symfony services, events, console commands, and bundles.

RabbitMQ is useful in Pimcore for:

- product indexing
- asset processing
- ERP sync
- PIM data export
- image transformation
- cache invalidation
- webhook delivery

---

## Pimcore Bundle Structure

```text
src/RabbitMqBundle/
├── RabbitMqBundle.php
├── Resources/
│   └── config/
│       └── services.yaml
├── EventSubscriber/
│   └── ProductChangedSubscriber.php
├── Command/
│   └── ConsumeRabbitMqCommand.php
└── Handler/
    └── ProductUpdatedHandler.php
```

---

## Publish When DataObject Changes

```php
use Pimcore\Event\DataObjectEvent;
use Pimcore\Event\DataObjectEvents;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;

final class ProductChangedSubscriber implements EventSubscriberInterface
{
    public function __construct(private RabbitPublisher $publisher) {}

    public static function getSubscribedEvents(): array
    {
        return [
            DataObjectEvents::POST_UPDATE => 'onPostUpdate',
        ];
    }

    public function onPostUpdate(DataObjectEvent $event): void
    {
        $object = $event->getObject();

        if (! $object instanceof Product) {
            return;
        }

        $this->publisher->publish('product.updated', [
            'product_id' => $object->getId(),
        ], [
            'source' => 'pimcore',
            'class' => Product::class,
        ]);
    }
}
```

---

## Pimcore Handler Example

```php
final class ProductUpdatedHandler implements MessageHandler
{
    public function __construct(private SearchIndexer $indexer) {}

    public function handle(MessageEnvelope $message): void
    {
        $product = Product::getById((int) $message->payload['product_id']);

        if (! $product instanceof Product) {
            throw new RuntimeException('Product not found.');
        }

        $this->indexer->index($product);
    }
}
```

---

## Pimcore CLI Worker

Because Pimcore has Symfony console:

```bash
bin/console app:rabbitmq:consume product.search_index.queue
```

Run it under Supervisor or Docker.

---

## Pimcore Design Advice

- Publish small IDs, not full object dumps.
- Load fresh Pimcore objects inside handlers.
- Do heavy sync/index/export work in workers.
- Use DLQ for failed product/asset sync messages.
- Include object class and ID in metadata.
- Do not block admin save actions with long external API calls.


---

# Chapter 6: PHP Retries, DLQ & Operations

## PHP Worker Risks

Long-running PHP workers need care:

- memory can grow over time
- database connections can go stale
- service container state can become old
- deploys need worker restarts
- fatal errors can kill worker

Frameworks hide some of this, but your RabbitMQ wrapper should still be defensive.

---

## Retry Policy

```php
final class RetryPolicy
{
    public function __construct(
        private int $maxRetries = 3
    ) {}

    public function shouldRetry(MessageEnvelope $message): bool
    {
        return ($message->meta['retry_count'] ?? 0) < $this->maxRetries;
    }

    public function nextRetryCount(MessageEnvelope $message): int
    {
        return ((int) ($message->meta['retry_count'] ?? 0)) + 1;
    }
}
```

---

## Error Classes

```php
class TemporaryMessageException extends RuntimeException {}

class PermanentMessageException extends RuntimeException {}
```

Consumer behavior:

```text
TemporaryMessageException -> retry
PermanentMessageException -> DLQ
Throwable -> retry until limit, then DLQ
```

---

## Memory Limit Restart

```php
if (memory_get_usage(true) > 256 * 1024 * 1024) {
    $logger->warning('worker.memory_limit_reached');
    exit(0);
}
```

Supervisor restarts the worker.

---

## Supervisor Config

```ini
[program:php-rabbitmq-worker]
command=php /var/www/app/bin/console app:rabbitmq:consume product.search_index.queue
directory=/var/www/app
autostart=true
autorestart=true
stopwaitsecs=60
user=www-data
redirect_stderr=true
stdout_logfile=/var/log/php-rabbitmq-worker.log
```

---

## Graceful Shutdown

Use signal handling when available:

```php
pcntl_signal(SIGTERM, function () use (&$running) {
    $running = false;
});
```

The worker should finish the current message, then exit.

---

## PHP Production Checklist

- `composer install --no-dev --optimize-autoloader`
- durable exchange and queue
- persistent messages
- manual ack
- prefetch configured
- structured logs
- retry queue
- DLQ
- idempotency table
- Supervisor/systemd/Docker restart
- graceful shutdown
- memory limit restart

