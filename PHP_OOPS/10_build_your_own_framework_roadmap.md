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

