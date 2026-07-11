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

