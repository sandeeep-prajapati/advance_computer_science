# Chapter 14: Security & Permissions

## Overview

Pimcore has a granular permission model because editors, developers, and API consumers all interact with the system differently. Security is handled at multiple layers:

1. **Users & Roles** — who can access the admin
2. **Object/Document/Asset permissions** — what content they can touch
3. **API security** — Symfony firewall + JWT/API keys
4. **Auditing** — log who did what

---

## 1. Users & Roles

- **Users** = actual accounts (human editors or system/API accounts)
- **Roles** = reusable permission sets assigned to multiple users
- Users can belong to **multiple roles** — permissions are merged (most permissive wins)

**Best Practice:** Define roles first, then assign users to roles.

**Example roles:**
- `content_editor` — can create/edit CMS pages
- `product_manager` — can manage Products and Categories
- `api_client` — API access only (no admin UI)

**Admin UI:** Settings → Users / Roles

---

## 2. Object / Document / Asset Permissions

Pimcore's tree-based structure makes it natural to apply permissions hierarchically.

### Available Permission Types

| Permission | Description |
|---|---|
| **List** | Can see the element in the tree |
| **View** | Can open and read the element |
| **Create** | Can create new elements |
| **Delete** | Can delete elements |
| **Publish** | Can publish/unpublish |
| **Write** | Can edit field values |

### Applying Permissions (Best Practice)

1. Create a folder (e.g., `/Products`, `/MarketingDocs`)
2. Right-click the folder → **Security** tab
3. Add the role (e.g., `product_manager`)
4. Grant permissions (List, View, Create, Delete, Publish, Write)
5. Children **inherit** from the parent folder

> Keep the `admin` role only for superusers. Give editors the **minimum necessary permissions** (least privilege principle).

### Example Permission Structure

```
/data-objects
  /products         → product_manager: Read + Write + Publish
    /electronics    → (inherits from /products)
    /clothing       → (inherits from /products)
  /customers        → customer_manager: Read only (GDPR)
  /orders           → order_manager: Read + Write

/assets
  /product-images   → product_manager: Upload + Edit
  /legal-docs       → legal_team only

/documents
  /en               → content_editor: Create + Edit + Publish
  /de               → german_editor: Create + Edit + Publish
```

---

## 3. Securing APIs

### Symfony Firewall Configuration

```yaml
# config/packages/security.yaml
security:
  providers:
    pimcore_admin:
      id: Pimcore\Security\User\UserProvider

  firewalls:
    # Protect admin UI
    pimcore_admin:
      pattern: ^/admin
      anonymous: false
      provider: pimcore_admin

    # API endpoints — stateless with token auth
    api:
      pattern: ^/api
      stateless: true
      json_login:
        check_path: /api/login
        username_path: username
        password_path: password
      jwt: true  # if using lexik/jwt-authentication-bundle

  access_control:
    - { path: ^/admin, roles: ROLE_PIMCORE_ADMIN }
    - { path: ^/api,   roles: ROLE_API_USER }
```

### JWT Authentication Setup

```bash
composer require lexik/jwt-authentication-bundle
```

```yaml
# config/packages/lexik_jwt_authentication.yaml
lexik_jwt_authentication:
    secret_key: '%env(JWT_SECRET_KEY)%'
    public_key: '%env(JWT_PUBLIC_KEY)%'
    pass_phrase: '%env(JWT_PASSPHRASE)%'
    token_ttl: 3600  # 1 hour
```

```bash
# Generate keys
php bin/console lexik:jwt:generate-keypair
```

**Get a token:**
```bash
curl -X POST http://localhost/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "api_user", "password": "secret"}'
```

**Use the token:**
```bash
curl http://localhost/api/products \
  -H "Authorization: Bearer eyJ0eXAiOiJ..."
```

---

## 4. DataHub (GraphQL) Security

```yaml
# In DataHub configuration (Admin UI):
# - Assign permissions per endpoint
# - Example: Product API allows only queries (no mutations)
# - Control which fields are exposed (hide sensitive fields)
# - Set allowed IP ranges for API access
```

---

## 5. Auditing Admin Actions

### Built-in Logs

Pimcore logs user actions in `var/log/admin-YYYY-MM-DD.log`:

```
2025-09-26 12:01:23 User "editor1" published object ID 123
2025-09-26 12:05:44 User "manager1" deleted asset ID 456
```

### Versioning

Enable **versioning** on Data Objects → full history of changes is stored:
1. Open Settings → Data Objects → Classes → `Product`
2. Enable "Save Versions" in class settings

### Custom Audit Subscriber

```php
<?php
namespace App\EventSubscriber;

use Pimcore\Event\DataObjectEvents;
use Pimcore\Event\Model\DataObjectEvent;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use Pimcore\Model\DataObject\Product;

class ProductAuditSubscriber implements EventSubscriberInterface
{
    public static function getSubscribedEvents(): array
    {
        return [
            DataObjectEvents::POST_DELETE => 'onProductDeleted',
            DataObjectEvents::POST_UPDATE => 'onProductUpdated',
        ];
    }

    public function onProductDeleted(DataObjectEvent $event): void
    {
        $product = $event->getObject();
        if (!$product instanceof Product) return;

        \Pimcore\Log\Simple::log(
            'audit',
            sprintf("Product deleted: ID=%d, SKU=%s", $product->getId(), $product->getSku())
        );
    }

    public function onProductUpdated(DataObjectEvent $event): void
    {
        $product = $event->getObject();
        if (!$product instanceof Product) return;

        \Pimcore\Log\Simple::log(
            'audit',
            sprintf("Product updated: ID=%d, Name=%s", $product->getId(), $product->getName())
        );
    }
}
```

---

## 6. API Key Authentication (Simpler Alternative to JWT)

For server-to-server or simple API clients:

```php
// src/EventSubscriber/ApiKeyAuthenticator.php
namespace App\Security;

use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Security\Core\Exception\AuthenticationException;
use Symfony\Component\Security\Http\Authenticator\AbstractAuthenticator;
use Symfony\Component\Security\Http\Authenticator\Passport\Passport;
use Symfony\Component\Security\Http\Authenticator\Passport\SelfValidatingPassport;
use Symfony\Component\Security\Http\Authenticator\Passport\Badge\UserBadge;

class ApiKeyAuthenticator extends AbstractAuthenticator
{
    public function supports(Request $request): ?bool
    {
        return $request->headers->has('X-API-KEY');
    }

    public function authenticate(Request $request): Passport
    {
        $apiKey = $request->headers->get('X-API-KEY');
        return new SelfValidatingPassport(new UserBadge($apiKey));
    }

    public function onAuthenticationFailure(Request $request, AuthenticationException $exception): JsonResponse
    {
        return new JsonResponse(['error' => 'Invalid API Key'], 401);
    }

    public function onAuthenticationSuccess(Request $request, $token, string $firewallName): ?Response
    {
        return null; // continue request
    }
}
```

---

## Security Best Practices Checklist

| Practice | Description |
|---|---|
| Use roles, not per-user permissions | Easier to manage at scale |
| Least privilege | Give only what's needed |
| Secure folder trees | Apply at folder root, let children inherit |
| Stateless APIs | Use JWT or API keys — never session-based |
| Never expose `/admin` publicly | Use IP whitelisting or VPN |
| Rate-limit APIs | At proxy/load balancer level |
| Enable versioning | Full change history for Data Objects |
| Audit logs | Keep `admin-*.log` files and custom event log |
| Separate system accounts | API clients get their own user (not human editor accounts) |
| Review DataHub configs | Only expose fields that are safe for public consumption |
