# Chapter 5: Kernel, Router, Middleware & Controller

## Front Controller

`public/index.php` is the single HTTP entry point.

```php
require dirname(__DIR__) . '/vendor/autoload.php';

$kernel = require dirname(__DIR__) . '/bootstrap/app.php';

$request = Request::capture();
$response = $kernel->handle($request);
$response->send();
$kernel->terminate($request, $response);
```

---

## Kernel

The kernel controls lifecycle.

```php
interface HttpKernel
{
    public function bootstrap(): void;

    public function handle(Request $request): Response;

    public function terminate(Request $request, Response $response): void;
}
```

Responsibilities:

- load config
- register providers
- boot bundles
- run middleware
- dispatch controller
- terminate after response

---

## Router

```php
final class Router
{
    private array $routes = [];

    public function get(string $path, mixed $handler): void
    {
        $this->routes['GET'][$path] = $handler;
    }

    public function match(Request $request): mixed
    {
        return $this->routes[$request->method][$request->path]
            ?? throw new RuntimeException('Route not found.');
    }
}
```

---

## Middleware

Middleware wraps request handling.

```php
interface Middleware
{
    public function handle(Request $request, callable $next): Response;
}
```

Examples:

- auth
- CSRF
- session
- logging
- CORS
- rate limiting

---

## Controller

Controllers translate HTTP to application actions.

```php
final class ProductController
{
    public function __construct(private CreateProduct $createProduct) {}

    public function store(Request $request): Response
    {
        $id = $this->createProduct->handle($request->body);

        return JsonResponse::created(['id' => $id]);
    }
}
```

Keep business logic in services, not controllers.

