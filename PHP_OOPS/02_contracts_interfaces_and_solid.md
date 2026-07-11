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

