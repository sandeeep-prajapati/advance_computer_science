# Chapter 10: Bundles, Plugins & Extensions

## Bundles vs. Plugins

| | Bundles (Symfony) | Plugins (Pimcore) |
|---|---|---|
| **What they are** | Standard Symfony packages | Pimcore-specific UI/feature extensions |
| **Installed via** | Composer | Composer or manually under `/bundles` |
| **Can add** | Services, controllers, templates, commands | Custom admin menus, dashboard widgets, object field types, workflows |
| **Use when** | Business logic, APIs, integrations | Extending admin UI, editor features |

---

## Building a Custom Pimcore Bundle

We'll build a **DemoBundle** as a complete example.

### Bundle Directory Layout

```
project-root/
└── bundles/
    └── DemoBundle/
        ├── composer.json
        ├── DemoBundle.php
        └── src/
            ├── Controller/
            │   └── DefaultController.php
            ├── Service/
            │   └── GreetingService.php
            └── Resources/
                └── config/
                    └── services.yaml
```

> **Convention:** Custom bundles live in `/bundles/` to stay isolated from app logic in `/src/`.

---

### Step 1: composer.json (for the bundle)

```json
{
    "name": "mycompany/demo-bundle",
    "type": "pimcore-bundle",
    "description": "A demo Pimcore bundle",
    "autoload": {
        "psr-4": {
            "DemoBundle\\": "src/"
        }
    },
    "require": {
        "php": ">=8.1"
    }
}
```

---

### Step 2: The Bundle Class

```php
<?php
// bundles/DemoBundle/DemoBundle.php
namespace DemoBundle;

use Pimcore\Extension\Bundle\AbstractPimcoreBundle;

class DemoBundle extends AbstractPimcoreBundle
{
    public function getJsPaths(): array
    {
        // Register custom admin UI JavaScript files
        return [
            '/bundles/demobundle/js/admin.js',
        ];
    }

    public function getCssPaths(): array
    {
        return [];
    }
}
```

---

### Step 3: A Service

```php
<?php
// bundles/DemoBundle/src/Service/GreetingService.php
namespace DemoBundle\Service;

class GreetingService
{
    public function greet(string $name): string
    {
        return "Hello, " . $name . "!";
    }
}
```

---

### Step 4: A Controller

```php
<?php
// bundles/DemoBundle/src/Controller/DefaultController.php
namespace DemoBundle\Controller;

use DemoBundle\Service\GreetingService;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

class DefaultController extends AbstractController
{
    #[Route('/demo/hello/{name}', name: 'demo_hello')]
    public function hello(string $name, GreetingService $greeter): Response
    {
        return $this->json([
            'message' => $greeter->greet($name)
        ]);
    }
}
```

Visiting `/demo/hello/Sandeep` returns: `{ "message": "Hello, Sandeep!" }`

---

### Step 5: services.yaml

```yaml
# bundles/DemoBundle/src/Resources/config/services.yaml
services:
    DemoBundle\:
        resource: '../../*'
        exclude: '../../{Resources}'
```

---

### Step 6: Enable the Bundle

```bash
bin/console pimcore:bundle:enable DemoBundle
bin/console pimcore:bundle:install DemoBundle

# Verify
bin/console debug:bundle
```

---

## Extending the Admin UI (ExtJS)

Pimcore's admin is built with **ExtJS**. You can add custom menus, toolbar buttons, and panels.

### Add a Toolbar Button

```javascript
// bundles/DemoBundle/src/Resources/public/js/admin.js

pimcore.registerNS("pimcore.plugin.DemoBundle");

pimcore.plugin.DemoBundle = Class.create(pimcore.plugin.admin, {
    getClassName: function () {
        return "pimcore.plugin.DemoBundle";
    },

    initialize: function () {
        pimcore.plugin.broker.registerPlugin(this);
    },

    pimcoreReady: function (params, broker) {
        var toolbar = pimcore.globalmanager.get("layout_toolbar");

        toolbar.add({
            text: "Demo Page",
            iconCls: "pimcore_icon_overlay_add",
            handler: function () {
                Ext.Ajax.request({
                    url: "/admin/demo/hello",
                    success: function (response) {
                        var data = Ext.decode(response.responseText);
                        Ext.Msg.alert("Demo Bundle", data.message);
                    }
                });
            }
        });
    }
});

new pimcore.plugin.DemoBundle();
```

### Add a Button to the Product Grid

```javascript
// Inside admin.js — extend existing object grids
pimcore.object.grid.prototype.getToolbar =
pimcore.object.grid.prototype.getToolbar.createSequence(function () {
    if (this.class === "Product") {
        this.toolbar.add({
            text: "Custom Action",
            iconCls: "pimcore_icon_apply",
            handler: function () {
                var selected = this.getSelectionModel().getSelection();
                if (selected.length) {
                    Ext.Msg.alert("Custom Action", "Selected product ID: " + selected[0].data.id);
                }
            }.bind(this)
        });
    }
});
```

---

## Admin Backend Controller for Custom Pages

```php
<?php
// bundles/DemoBundle/src/Controller/AdminController.php
namespace DemoBundle\Controller;

use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/admin/demo', name: 'demo_admin_')]
class AdminController extends AbstractController
{
    #[Route('/hello', name: 'hello')]
    public function hello(): JsonResponse
    {
        return new JsonResponse(['message' => 'Hello from Demo Admin Controller!']);
    }
}
```

---

## How It All Fits Together

```
1. DemoBundle.php  →  getJsPaths() registers admin.js
2. admin.js        →  hooks into Pimcore UI lifecycle via pimcore.plugin.broker
3. AdminController →  provides JSON APIs for AJAX calls from admin.js
4. Result          →  editors see new UI buttons and panels seamlessly integrated
```

---

## Custom Package in a Pimcore Project (Alternative Layout)

If you want to put a package inside `src/`:

```
/src/MyPackage/
├── Controller/
│   └── AdminController.php
├── Resources/
│   └── config/
│       └── routes.yaml
└── MyPackageBundle.php
```

Register routes in `config/routes.yaml`:
```yaml
my_package:
    resource: '../src/MyPackage/Resources/config/routes.yaml'
```

Register as service in `config/services.yaml`:
```yaml
services:
    App\MyPackage\:
        resource: '../src/MyPackage/*'
        tags: ['controller.service_arguments']
```

---

## Pro Tips

- **Share publicly:** Push to GitHub/Packagist → install via `composer require mycompany/demo-bundle`
- **Bundles = Symfony power** (custom logic, APIs, integrations)
- **Plugins = Deep Pimcore integration** (UI, workflows, DAM, PIM features)
- Always run `bin/console cache:clear` after enabling/installing bundles
