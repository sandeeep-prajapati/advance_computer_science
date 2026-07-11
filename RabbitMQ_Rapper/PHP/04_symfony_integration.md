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

