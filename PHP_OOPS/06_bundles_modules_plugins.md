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

