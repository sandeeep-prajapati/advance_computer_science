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

