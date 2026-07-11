# Chapter 1: OOP Foundation for Frameworks

## Why OOP Matters in Framework Design

A framework is a group of collaborating objects.

```text
Request
└── Kernel
    ├── Router
    ├── Middleware
    ├── ControllerResolver
    ├── Controller
    ├── Service
    └── Response
```

OOP helps you give each object a clear responsibility.

---

## Class and Object

```php
final class Request
{
    public function __construct(
        public readonly string $method,
        public readonly string $path,
        public readonly array $query = [],
        public readonly array $body = []
    ) {}
}

$request = new Request('GET', '/products');
```

The class is the design. The object is the runtime value.

---

## Encapsulation

Encapsulation hides internal details.

```php
final class Config
{
    public function __construct(private array $items) {}

    public function get(string $key, mixed $default = null): mixed
    {
        return $this->items[$key] ?? $default;
    }
}
```

Other code does not need to know if config came from PHP files, YAML, `.env`, database, or cache.

---

## Composition Over Inheritance

Bad:

```php
class ProductController extends DatabaseConnection
{
}
```

A controller is not a database connection.

Better:

```php
final class ProductController
{
    public function __construct(private ProductRepository $products) {}
}
```

The controller has a repository. This is composition.

---

## Framework Object Example

```php
final class Application
{
    public function __construct(
        private Router $router,
        private ControllerDispatcher $dispatcher
    ) {}

    public function handle(Request $request): Response
    {
        $route = $this->router->match($request);

        return $this->dispatcher->dispatch($route, $request);
    }
}
```

This is small, but it already looks like a framework.

---

## Value Objects

Value objects represent concepts safely.

```php
final readonly class RouteName
{
    public function __construct(public string $value)
    {
        if ($value === '') {
            throw new InvalidArgumentException('Route name cannot be empty.');
        }
    }
}
```

Use value objects for IDs, money, slugs, route names, emails, date ranges, and configuration keys.

---

## Practice

Build these:

- `Request`
- `Response`
- `Route`
- `Router`
- `Application`

Goal:

```php
$app->handle(new Request('GET', '/'));
```

