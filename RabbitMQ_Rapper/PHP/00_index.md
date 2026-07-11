# PHP RabbitMQ Wrapper Book

Detailed notes for building RabbitMQ wrappers in PHP.

This folder covers:

- Core PHP with `php-amqplib/php-amqplib`
- Laravel service provider + Artisan worker
- Symfony bundle + console command
- Pimcore bundle/event integration
- Composer package design
- Retry, DLQ, idempotency, and deployment for PHP workers

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_php_package_architecture.md) | **PHP Package Architecture** | Composer, PSR-4, contracts, package layout |
| [02](02_core_php_consumer_publisher.md) | **Core PHP Consumer & Publisher** | php-amqplib, connection, publisher, consumer, CLI worker |
| [03](03_laravel_integration.md) | **Laravel Integration** | config, service provider, Artisan command, handlers through container |
| [04](04_symfony_integration.md) | **Symfony Integration** | services.yaml, bundle, console command, tagged handlers |
| [05](05_pimcore_integration.md) | **Pimcore Integration** | Pimcore events, DataObject publishing, bundle design |
| [06](06_php_retries_dlq_operations.md) | **PHP Retries, DLQ & Operations** | retry policy, DLQ, Supervisor, memory limits, graceful restart |
| [book](book.md) | **Combined PHP Book** | All PHP chapters in one file |

---

## PHP Wrapper Goal

Your PHP application should use RabbitMQ like this:

```php
$publisher->publish('product.updated', [
    'product_id' => 501,
]);
```

And your worker should look like this:

```bash
php bin/consume product.search_index.queue
php artisan rabbitmq:consume product.search_index.queue
php bin/console app:rabbitmq:consume product.search_index.queue
```

The business handler should not know about AMQP channels, delivery tags, or queue declarations.

