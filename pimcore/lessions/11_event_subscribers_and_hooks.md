# Chapter 11: Event Subscribers & Hooks

## Overview

Pimcore builds on Symfony's event system and adds its own events for Data Objects, Assets, and Documents. Event subscribers let you **react to things happening** in Pimcore without modifying core code.

Common uses:
- Auto-generate a slug when a product is saved
- Send a notification when content is published
- Sync data to external systems after save
- Log admin actions for auditing

---

## Pimcore Data Object Events

```php
use Pimcore\Event\DataObjectEvents;

DataObjectEvents::PRE_UPDATE   // before saving an existing object
DataObjectEvents::POST_UPDATE  // after saving
DataObjectEvents::PRE_ADD      // before creating a new object
DataObjectEvents::POST_ADD     // after creating
DataObjectEvents::PRE_DELETE   // before deleting
DataObjectEvents::POST_DELETE  // after deleting
```

---

## Example 1: Auto-Generate Slug on Product Save

### Step 1: Add `slug` Field to Product Class

In Admin UI → Settings → Data Objects → Classes → Product:
- Add a field: `slug` (type: Input, optional: read-only)

### Step 2: Create the Event Subscriber

**File:** `src/EventSubscriber/ProductSlugSubscriber.php`

```php
<?php
namespace App\EventSubscriber;

use Pimcore\Event\DataObjectEvents;
use Pimcore\Event\Model\DataObjectEvent;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use Pimcore\Model\DataObject\Product;

class ProductSlugSubscriber implements EventSubscriberInterface
{
    public static function getSubscribedEvents(): array
    {
        return [
            DataObjectEvents::PRE_UPDATE => 'onProductSave',
            DataObjectEvents::PRE_ADD    => 'onProductSave',
        ];
    }

    public function onProductSave(DataObjectEvent $event): void
    {
        $object = $event->getObject();

        // Only handle Product objects
        if (!$object instanceof Product) {
            return;
        }

        $name = $object->getName();
        if ($name) {
            // Slugify: lowercase + replace non-alphanumeric with dashes
            $slug = strtolower(trim(preg_replace('/[^A-Za-z0-9-]+/', '-', $name), '-'));
            $object->setSlug($slug);
        }
    }
}
```

### Step 3: Register in services.yaml

```yaml
# config/services.yaml
services:
    App\EventSubscriber\ProductSlugSubscriber:
        tags:
            - { name: kernel.event_subscriber }
```

### What Happens

1. You save or create a Product: "iPhone 16 Pro"
2. Event fires → `onProductSave()` runs
3. Slug is auto-populated: `iphone-16-pro`
4. Object is persisted with the slug already set

---

## Example 2: Admin Audit Log (Login Tracking)

**File:** `src/EventSubscriber/AdminAuditSubscriber.php`

```php
<?php
namespace App\EventSubscriber;

use Pimcore\Event\AdminEvents;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;

class AdminAuditSubscriber implements EventSubscriberInterface
{
    public static function getSubscribedEvents(): array
    {
        return [
            AdminEvents::ADMIN_LOGIN => 'onAdminLogin',
        ];
    }

    public function onAdminLogin($event): void
    {
        $request = $event->getRequest();
        $user = $event->getUser();

        \Pimcore\Log\Simple::log(
            'admin_audit',
            "User {$user->getName()} logged in from {$request->getClientIp()}"
        );
    }
}
```

---

## Example 3: Sync to External System After Product Save

```php
<?php
namespace App\EventSubscriber;

use Pimcore\Event\DataObjectEvents;
use Pimcore\Event\Model\DataObjectEvent;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use Pimcore\Model\DataObject\Product;
use App\Service\ExternalSyncService;

class ProductSyncSubscriber implements EventSubscriberInterface
{
    public function __construct(private ExternalSyncService $syncService)
    {}

    public static function getSubscribedEvents(): array
    {
        return [
            DataObjectEvents::POST_UPDATE => 'onProductSaved',
        ];
    }

    public function onProductSaved(DataObjectEvent $event): void
    {
        $object = $event->getObject();

        if (!$object instanceof Product) {
            return;
        }

        // Sync to ERP, Elasticsearch, external API, etc.
        $this->syncService->syncProduct($object);
    }
}
```

---

## Available Events

### Data Object Events
| Event | When |
|---|---|
| `PRE_ADD` | Before a new object is created |
| `POST_ADD` | After a new object is created |
| `PRE_UPDATE` | Before an existing object is saved |
| `POST_UPDATE` | After an existing object is saved |
| `PRE_DELETE` | Before an object is deleted |
| `POST_DELETE` | After an object is deleted |
| `PRE_PUBLISH` | Before an object is published |
| `POST_PUBLISH` | After an object is published |

### Asset Events
```php
use Pimcore\Event\AssetEvents;

AssetEvents::PRE_UPDATE
AssetEvents::POST_UPDATE
AssetEvents::PRE_ADD
AssetEvents::POST_ADD
```

### Document Events
```php
use Pimcore\Event\DocumentEvents;

DocumentEvents::PRE_UPDATE
DocumentEvents::POST_UPDATE
```

### Admin Events
```php
use Pimcore\Event\AdminEvents;

AdminEvents::ADMIN_LOGIN
AdminEvents::ADMIN_LOGOUT
```

---

## Tips & Advanced Usage

### Unique Slugs

If you need unique slugs, check for existing ones and append a counter:

```php
$baseSlug = $slug;
$counter = 1;
while ($this->slugExists($slug, $object->getId())) {
    $slug = $baseSlug . '-' . $counter++;
}
$object->setSlug($slug);
```

### Use Symfony's SluggerInterface

```php
use Symfony\Component\String\Slugger\SluggerInterface;

class ProductSlugSubscriber implements EventSubscriberInterface
{
    public function __construct(private SluggerInterface $slugger)
    {}

    public function onProductSave(DataObjectEvent $event): void
    {
        $object = $event->getObject();
        if (!$object instanceof Product) return;

        $slug = strtolower($this->slugger->slug($object->getName()));
        $object->setSlug($slug);
    }
}
```

### preUpdate vs postUpdate

| Use | When |
|---|---|
| `PRE_UPDATE` | Modify the object **before** it's saved (slug generation, validation) |
| `POST_UPDATE` | Act **after** save (sync to Elasticsearch, send webhook, update cache) |

> Never save the object again inside a `PRE_UPDATE` subscriber — it causes infinite loops.
