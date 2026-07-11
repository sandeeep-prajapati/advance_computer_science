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

