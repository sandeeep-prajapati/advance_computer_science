# Chapter 4: Service Container & Dependency Injection

## What a Container Does

A container creates objects and injects dependencies.

```text
ProductController
└── CreateProduct
    ├── ProductRepository
    └── EventDispatcher
```

You ask for `ProductController`. The container builds the full object graph.

---

## Container API

```php
interface Container
{
    public function bind(string $id, callable|string $resolver): void;

    public function singleton(string $id, callable|string $resolver): void;

    public function get(string $id): object;
}
```

---

## Binding Interfaces

```php
$container->bind(ProductRepository::class, MysqlProductRepository::class);
$container->singleton(EventDispatcher::class, EventDispatcher::class);
```

Now classes can ask for interfaces:

```php
final class CreateProduct
{
    public function __construct(private ProductRepository $products) {}
}
```

---

## Service Provider

Laravel style:

```php
final class ProductServiceProvider implements ServiceProvider
{
    public function register(Container $container): void
    {
        $container->bind(ProductRepository::class, MysqlProductRepository::class);
    }

    public function boot(Container $container): void
    {
        $container->get(Router::class)
            ->get('/products', [ProductController::class, 'index']);
    }
}
```

`register()` binds services. `boot()` connects services after all bindings exist.

---

## Tagged Services

Symfony style:

```text
tag: console.command
tag: event.subscriber
tag: cache.warmer
tag: pimcore.bundle
```

Tags let a framework collect services automatically.

Example:

```php
foreach ($container->tagged('console.command') as $command) {
    $console->add($command);
}
```

---

## Rule

Use the container at framework edges. Do not call it everywhere inside business logic.

Good:

```php
final class ProductController
{
    public function __construct(private CreateProduct $createProduct) {}
}
```

Bad:

```php
$service = app()->get(CreateProduct::class);
```

inside every method.

