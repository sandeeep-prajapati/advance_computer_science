# PHP OOP Framework Design: The Complete Book

> A practical single-file guide to PHP OOP concepts for designing Symfony, Laravel, and Pimcore-style frameworks with Composer packages, UI asset integration, and independent React/Next/TypeScript frontends.

---

# PHP OOP Framework Design Book

A practical guide to PHP OOP concepts that teach you how to design frameworks like **Symfony**, **Laravel**, and **Pimcore**, with **Composer packages**, **bundles**, **service containers**, **facades**, and frontend integration using **JavaScript**, **React**, **Next.js**, and **TypeScript**.

This book is not only about PHP syntax. It is about learning how framework internals are designed.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_oop_foundation_for_frameworks.md) | **OOP Foundation for Frameworks** | Classes, objects, encapsulation, composition, lifecycle objects |
| [02](02_contracts_interfaces_and_solid.md) | **Contracts, Interfaces & SOLID** | Interfaces, abstract classes, dependency inversion, framework boundaries |
| [03](03_composer_autoloading_and_packages.md) | **Composer, Autoloading & Packages** | PSR-4, package structure, scripts, service discovery, reusable libraries |
| [04](04_service_container_and_dependency_injection.md) | **Service Container & Dependency Injection** | Binding, singleton, auto-wiring, service providers, tagged services |
| [05](05_kernel_router_middleware_controller.md) | **Kernel, Router, Middleware & Controller** | Front controller, HTTP kernel, route resolver, controller dispatcher |
| [06](06_bundles_modules_plugins.md) | **Bundles, Modules & Plugins** | Symfony bundles, Pimcore bundles, Laravel packages, extension points |
| [07](07_facades_events_commands.md) | **Facades, Events & Commands** | Laravel facades, event dispatcher, console commands, jobs |
| [08](08_ui_asset_integration.md) | **UI & Asset Integration** | Blade/Twig, Vite, npm, React widgets, Pimcore admin assets |
| [09](09_react_next_typescript_frontend.md) | **React, Next.js & TypeScript Frontends** | API-first backend, independent frontend, auth, DTOs, generated clients |
| [10](10_build_your_own_framework_roadmap.md) | **Build Your Own Framework Roadmap** | Step-by-step mini-framework plan |

---

## How to Read

**If you are new to OOP:** Read 1 → 2 first.

**If you want Symfony/Pimcore style:** Read 3 → 4 → 5 → 6 → 8.

**If you want Laravel style:** Read 3 → 4 → 5 → 7 → 8.

**If you want independent frontend architecture:** Read 8 → 9.

**If you want to build your own framework:** Read all chapters, then follow Chapter 10.

---

## Big Picture

```text
PHP Framework
├── public/index.php          front controller
├── Kernel                    bootstraps app and handles lifecycle
├── Service Container         creates and stores services
├── Router                    maps URL to controller
├── Middleware                wraps request/response flow
├── Controller                translates HTTP to application action
├── Events                    decouple side effects
├── Console Commands          CLI tools like migrate/cache/make
├── Bundles/Packages          reusable modules
├── Composer                  autoloading and dependency management
└── UI Integration
    ├── Blade/Twig templates
    ├── Vite/npm assets
    ├── React widgets inside PHP pages
    └── separate Next.js/TypeScript frontend
```

---

## Core Design Rule

Keep your **business logic independent**.

```text
Good:
Controller -> Application Service -> Domain Object -> Repository Interface

Bad:
Controller -> SQL + RabbitMQ + HTML + API calls + business rules in one method
```

Frameworks are powerful when they separate concerns cleanly.


---

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


---

# Chapter 2: Contracts, Interfaces & SOLID

## Frameworks Depend on Contracts

Laravel, Symfony, and Pimcore are flexible because high-level code depends on interfaces.

```php
interface CacheStore
{
    public function get(string $key, mixed $default = null): mixed;

    public function put(string $key, mixed $value, int $seconds): void;
}
```

Implementations can change:

- array cache
- file cache
- Redis cache
- database cache

The application code stays stable.

---

## Interface vs Abstract Class

| Tool | Use When |
|---|---|
| interface | you need a contract |
| abstract class | you need shared base behavior |
| trait | you need reusable methods, carefully |

Example:

```php
interface ServiceProvider
{
    public function register(Container $container): void;

    public function boot(Container $container): void;
}
```

---

## SOLID in Framework Design

### Single Responsibility

One object, one reason to change.

```text
Router      matches routes
Container   resolves services
Kernel      controls lifecycle
Controller  handles HTTP boundary
Repository  hides persistence
```

### Open/Closed

Add new behavior through interfaces, bundles, events, and middleware.

### Dependency Inversion

```php
final class SendInvoice
{
    public function __construct(private Mailer $mailer) {}
}
```

It depends on `Mailer`, not `SmtpMailer`.

---

## Boundary Thinking

| Layer | Should Know | Should Not Know |
|---|---|---|
| Controller | Request, response, validation | SQL details |
| Service | Business use case | HTTP framework |
| Repository | Storage | HTML rendering |
| Kernel | Boot lifecycle | Product rules |
| Frontend | API contract | database tables |

This is how you design a framework that does not become tangled.


---

# Chapter 3: Composer, Autoloading & Packages

## Composer Mental Model

Composer handles PHP dependencies and autoloading.

```text
composer.json
└── composer install
    ├── vendor/
    ├── vendor/autoload.php
    └── composer.lock
```

Almost every PHP framework starts with:

```php
require dirname(__DIR__) . '/vendor/autoload.php';
```

---

## PSR-4 Autoloading

```json
{
  "autoload": {
    "psr-4": {
      "App\\": "app/",
      "Framework\\": "framework/",
      "Acme\\BlogBundle\\": "packages/blog-bundle/src/"
    }
  }
}
```

Composer maps namespaces to folders.

```text
Framework\Http\Request
→ framework/Http/Request.php
```

---

## Package Structure

Reusable framework package:

```text
packages/blog-bundle/
├── composer.json
├── src/
│   ├── BlogBundle.php
│   ├── BlogServiceProvider.php
│   ├── Controller/
│   ├── Command/
│   └── Resources/
│       ├── views/
│       └── assets/
└── tests/
```

---

## Composer Package Metadata

```json
{
  "name": "acme/blog-bundle",
  "type": "framework-bundle",
  "autoload": {
    "psr-4": {
      "Acme\\BlogBundle\\": "src/"
    }
  },
  "extra": {
    "my-framework": {
      "bundles": [
        "Acme\\BlogBundle\\BlogBundle"
      ],
      "providers": [
        "Acme\\BlogBundle\\BlogServiceProvider"
      ]
    }
  }
}
```

Frameworks use `extra` metadata for package discovery.

Laravel uses providers. Symfony uses bundles. Pimcore uses bundles on top of Symfony.

---

## Composer Scripts

```json
{
  "scripts": {
    "serve": "php -S 127.0.0.1:8000 -t public",
    "test": "phpunit",
    "analyse": "phpstan analyse app framework",
    "format": "pint",
    "check": [
      "@format",
      "@analyse",
      "@test"
    ]
  }
}
```

Run:

```bash
composer serve
composer test
composer check
```

---

## Install vs Update

| Command | Meaning |
|---|---|
| `composer install` | install exact versions from lock file |
| `composer update` | resolve new versions and rewrite lock file |
| `composer require` | add package |
| `composer dump-autoload` | rebuild autoload files |

Use `install` normally. Use `update` intentionally.


---

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


---

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


---

# Chapter 6: Bundles, Modules & Plugins

## Why Bundles Exist

Bundles/package modules let you add features without editing framework core.

```text
BlogBundle
├── routes
├── controllers
├── services
├── commands
├── templates
├── assets
├── migrations
└── event subscribers
```

Symfony and Pimcore use bundle architecture. Laravel packages use service providers.

---

## Bundle Contract

```php
interface Bundle
{
    public function name(): string;

    public function register(Container $container): void;

    public function boot(Container $container): void;
}
```

---

## Example Bundle

```php
final class BlogBundle implements Bundle
{
    public function name(): string
    {
        return 'blog';
    }

    public function register(Container $container): void
    {
        $container->bind(PostRepository::class, MysqlPostRepository::class);
    }

    public function boot(Container $container): void
    {
        $container->get(Router::class)
            ->get('/blog', [PostController::class, 'index']);
    }
}
```

---

## Extension Points

Good frameworks provide clear extension points:

- service registration
- route registration
- event subscribers
- console commands
- template paths
- asset publishing
- migrations
- admin UI extensions
- custom field types

Pimcore-style systems need strong extension points because bundles often add admin features and data-object behavior.

---

## Laravel vs Symfony/Pimcore

| Concept | Laravel | Symfony/Pimcore |
|---|---|---|
| Package boot | Service Provider | Bundle |
| Config publish | `vendor:publish` | bundle config/packages |
| Commands | Artisan | Console command |
| Events | Event/listener | Event subscriber |
| Templates | Blade | Twig |
| Assets | Vite | Webpack Encore/Vite/custom assets |


---

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


---

# Chapter 8: UI & Asset Integration

## Two UI Styles

PHP frameworks usually support two frontend styles.

```text
Server-rendered UI
├── Laravel Blade
├── Symfony/Pimcore Twig
└── React widgets mounted inside pages

Independent frontend
├── React SPA
├── Next.js app
├── TypeScript client
└── PHP backend as API
```

---

## npm and Vite

Composer handles PHP. npm handles frontend assets.

```text
composer install
→ PHP packages and autoloading

npm install
→ JS/CSS dependencies

npm run build
→ compiled assets in public/build
```

Example `package.json`:

```json
{
  "private": true,
  "scripts": {
    "dev": "vite",
    "build": "vite build"
  },
  "devDependencies": {
    "vite": "^5.0.0",
    "typescript": "^5.0.0",
    "react": "^18.0.0",
    "react-dom": "^18.0.0"
  }
}
```

---

## React Widget in PHP Page

Twig/Blade renders a mount point:

```html
<div id="product-editor" data-product-id="1001"></div>
<script type="module" src="/build/product-editor.js"></script>
```

React mounts there:

```tsx
import { createRoot } from 'react-dom/client';

const el = document.getElementById('product-editor');

if (el) {
  createRoot(el).render(
    <ProductEditor productId={el.dataset.productId!} />
  );
}
```

This is useful for Pimcore admin extensions and Laravel/Symfony admin pages.

---

## Asset Publishing

Packages/bundles may contain assets.

```text
BlogBundle/
└── Resources/assets/
    ├── blog-admin.tsx
    └── blog-admin.css
```

Command:

```bash
php console assets:publish
```

Output:

```text
public/bundles/blog/
├── blog-admin.js
└── blog-admin.css
```

---

## Backend-to-Frontend Contract

Do not let frontend depend on database tables directly.

Good:

```text
React -> HTTP API -> Application Service -> Repository
```

Bad:

```text
React assumes exact database columns and hidden backend rules
```


---

# Chapter 9: React, Next.js & TypeScript Frontends

## Independent Frontend Architecture

When using Next.js or React as a separate app, PHP becomes an API backend.

```text
Next.js / React / TypeScript
└── calls REST or GraphQL API
    └── PHP Framework
        ├── Controller/API Resource
        ├── Application Service
        ├── Domain Logic
        └── Repository
```

This is framework-independent frontend design.

---

## API DTO

Backend response:

```json
{
  "id": 1001,
  "name": "Blue Shirt",
  "price": {
    "amount": 1999,
    "currency": "INR"
  }
}
```

TypeScript type:

```ts
export type ProductDto = {
  id: number;
  name: string;
  price: {
    amount: number;
    currency: string;
  };
};
```

---

## API Client

```ts
export async function getProduct(id: number): Promise<ProductDto> {
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/products/${id}`);

  if (!response.ok) {
    throw new Error('Failed to load product');
  }

  return response.json();
}
```

---

## PHP API Controller

```php
final class ProductApiController
{
    public function __construct(private ProductQuery $products) {}

    public function show(int $id): JsonResponse
    {
        return new JsonResponse($this->products->findDto($id));
    }
}
```

The controller returns DTOs, not raw database entities.

---

## Auth Options

| Style | Use Case |
|---|---|
| session cookie | same-domain app |
| Sanctum-style cookie auth | Laravel SPA |
| JWT | separate frontend/backend |
| OAuth2/OIDC | enterprise login |

Pick auth based on deployment shape.

---

## Frontend Independence Rule

React/Next should know API contracts, not PHP framework internals.

That means the frontend can change from React to Vue, or backend can change from Laravel to Symfony, if the API contract stays stable.


---

# Chapter 10: Build Your Own Framework Roadmap

## Goal

Build a small framework that teaches the internals behind Laravel, Symfony, and Pimcore.

---

## Step 1: Composer Project

Create:

```text
composer.json
public/index.php
bootstrap/app.php
framework/
app/
```

Add PSR-4:

```json
{
  "autoload": {
    "psr-4": {
      "App\\": "app/",
      "Framework\\": "framework/"
    }
  }
}
```

---

## Step 2: HTTP Core

Build:

- `Request`
- `Response`
- `Router`
- `Route`
- `ControllerDispatcher`

Goal:

```php
$router->get('/', fn () => new Response('Hello'));
```

---

## Step 3: Container

Build:

- `bind`
- `singleton`
- `get`
- constructor auto-wiring
- interface binding

Goal:

```php
$controller = $container->get(HomeController::class);
```

---

## Step 4: Kernel

Build:

- `HttpKernel`
- provider registration
- middleware pipeline
- terminate lifecycle

Goal:

```php
$response = $kernel->handle(Request::capture());
```

---

## Step 5: Packages and Bundles

Build:

- `Bundle` contract
- service provider contract
- package discovery through Composer `extra`
- route/service/command extension points

Goal:

```php
$kernel->registerBundle(new BlogBundle());
```

---

## Step 6: UI Integration

Build:

- template renderer
- Vite manifest reader
- asset publishing command
- React widget mount example

Goal:

```php
vite('resources/js/product-editor.tsx');
```

---

## Step 7: Independent Frontend

Build:

- `/api/products`
- DTO classes
- TypeScript API client
- Next.js page consuming PHP API

Goal:

```text
Next.js frontend -> PHP API -> Application Service
```

---

## Final Checklist

You understand framework design when you can explain:

- how Composer autoloading starts the app
- how the kernel handles lifecycle
- how the container resolves services
- how routes reach controllers
- how middleware wraps requests
- how bundles extend the framework
- how facades proxy to services
- how events decouple side effects
- how npm/Vite assets enter PHP pages
- how React/Next can live independently from PHP

