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

