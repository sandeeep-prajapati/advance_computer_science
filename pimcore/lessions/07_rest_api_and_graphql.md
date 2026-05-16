# Chapter 7: REST API & GraphQL

## Overview

Pimcore supports two major API approaches:
1. **Custom REST APIs** — built with Symfony controllers (full control)
2. **DataHub GraphQL** — official Pimcore bundle for flexible, frontend-friendly queries

---

## Option A: Custom REST API (Symfony Controller)

### Basic Product API

**File:** `src/Controller/Api/ProductApiController.php`

```php
<?php
namespace App\Controller\Api;

use Pimcore\Model\DataObject\Product;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Routing\Annotation\Route;

class ProductApiController extends AbstractController
{
    #[Route('/api/products', name: 'api_products', methods: ['GET'])]
    public function list(Request $request): JsonResponse
    {
        $page  = max(1, (int) $request->query->get('page', 1));
        $limit = min(50, (int) $request->query->get('limit', 10)); // cap for safety
        $name  = $request->query->get('name');

        $listing = new Product\Listing();
        $listing->setLimit($limit);
        $listing->setOffset(($page - 1) * $limit);

        if ($name) {
            $listing->addConditionParam("name LIKE ?", ["%$name%"]);
        }

        $items = [];
        foreach ($listing as $product) {
            $items[] = [
                'id'          => $product->getId(),
                'sku'         => $product->getSku(),
                'name'        => $product->getName(),
                'price'       => $product->getPrice(),
                'description' => $product->getDescription(),
                'category'    => $product->getCategory()?->getName(),
                'image'       => $product->getImage()
                                   ? $product->getImage()->getFullPath()
                                   : null,
            ];
        }

        return $this->json([
            'page'     => $page,
            'limit'    => $limit,
            'count'    => $listing->getTotalCount(),
            'products' => $items,
        ]);
    }
}
```

### Test with curl

```bash
# Fetch first page, 5 products
curl "http://localhost/api/products?page=1&limit=5"

# Filter by name containing 'iPhone'
curl "http://localhost/api/products?page=1&limit=10&name=iPhone"
```

### Sample JSON Response

```json
{
    "page": 1,
    "limit": 5,
    "count": 27,
    "products": [
        {
            "id": 123,
            "sku": "IP16-256GB",
            "name": "iPhone 16 Pro",
            "price": 119999,
            "description": "Latest flagship phone",
            "category": "Smartphones",
            "image": "/var/assets/products/iphone16.jpg"
        }
    ]
}
```

---

## Option B: Pimcore DataHub (GraphQL)

### Installation

```bash
composer require pimcore/data-hub
bin/console pimcore:bundle:enable PimcoreDataHubBundle
bin/console pimcore:bundle:install PimcoreDataHubBundle
```

### Create a GraphQL Endpoint

1. In Admin UI → **Settings → DataHub**
2. Click **Add Configuration**
3. Select **GraphQL** as type
4. Choose the Data Object class to expose (e.g., `Product`)
5. Configure which fields/relations are accessible
6. Save → Pimcore auto-generates a GraphQL endpoint

**Default endpoint URL:**
```
/graphql/{configName}
# Example:
http://localhost/graphql/products
```

### Example GraphQL Query

```graphql
query {
  getProductListing(first: 5) {
    edges {
      node {
        id
        name
        sku
        price
        description
        images {
          edges {
            node {
              path
              thumbnails {
                product_preview {
                  path
                }
              }
            }
          }
        }
        categories {
          edges {
            node {
              id
              name
            }
          }
        }
      }
    }
  }
}
```

### Fetching a Single Product

```graphql
query {
  getProduct(id: 123) {
    name
    sku
    price
    description
  }
}
```

### Response:
```json
{
    "data": {
        "getProduct": {
            "name": "Laptop X1",
            "sku": "LAP123",
            "price": 1200,
            "description": "A lightweight business laptop"
        }
    }
}
```

### Fetching a Document Page

```graphql
query {
  getDocument(id: 45) {
    fullPath
    title
    properties {
      key
      value
    }
  }
}
```

---

## Comparison: REST vs GraphQL

| Approach | Pros | Cons |
|---|---|---|
| **DataHub (GraphQL)** | Frontend devs love it; flexible field selection; configurable in Admin UI | Extra bundle; some setup overhead |
| **Custom REST (Symfony)** | Full control; lightweight; no extra dependencies | Fixed response shapes; less frontend flexibility |
| **API Platform** | Auto-generated REST + GraphQL APIs; powerful filters | Extra setup; can feel heavy for simple use cases |

### When to Use Which?
- **GraphQL (DataHub):** Frontend teams using React, Next.js, Vue that need flexible queries
- **Custom REST:** Backend-controlled APIs for performance/security
- **API Platform:** Teams that want auto-generated, standards-compliant APIs

---

## Securing Your APIs

```yaml
# config/packages/security.yaml
security:
  firewalls:
    api:
      pattern: ^/api
      stateless: true
      json_login:
        check_path: /api/login
  access_control:
    - { path: ^/api, roles: ROLE_API_USER }
```

### Best Practices
- Use **JWT or API keys** for machine-to-machine communication
- Never expose `/admin/*` endpoints publicly
- **Rate-limit** APIs at the proxy/load balancer level
- In DataHub: configure per-endpoint permissions (read-only, field restrictions)
- Always **cap the `limit` parameter** in REST APIs (max 50) to prevent large DB queries

---

## Real-World Use Cases

| Use Case | Solution |
|---|---|
| Headless website | Pimcore manages data → React/Vue fetches via GraphQL |
| Mobile app | iOS/Android fetches product catalog via REST |
| Marketplace | Syndicate product data to Amazon, Shopify via REST |
| Omnichannel | Same content powers website, app, in-store displays |
