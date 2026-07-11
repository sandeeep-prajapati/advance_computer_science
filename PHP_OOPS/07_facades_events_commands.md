# Chapter 7: Facades, Events & Commands

## Facades

Laravel-style facades give static-looking access to services.

```php
Cache::put('theme', 'dark', 3600);
Route::get('/admin', [AdminController::class, 'index']);
```

Internally, a facade resolves a real service from the container.

```php
abstract class Facade
{
    protected static Container $container;

    abstract protected static function accessor(): string;

    public static function __callStatic(string $method, array $arguments): mixed
    {
        $service = static::$container->get(static::accessor());

        return $service->{$method}(...$arguments);
    }
}
```

Use facades for developer experience. Use dependency injection for core business logic.

---

## Events

Events decouple actions from side effects.

```php
final readonly class ProductCreated
{
    public function __construct(public int $productId) {}
}
```

Dispatcher:

```php
$events->dispatch(new ProductCreated($id));
```

Listeners:

- send email
- update search index
- clear cache
- publish RabbitMQ message

---

## Console Commands

Frameworks need CLI tools.

```php
interface ConsoleCommand
{
    public function name(): string;

    public function handle(array $arguments): int;
}
```

Examples:

```bash
php console make:controller ProductController
php console route:list
php console cache:warm
php console migrate
php console assets:publish
```

Commands are just objects with dependencies.

