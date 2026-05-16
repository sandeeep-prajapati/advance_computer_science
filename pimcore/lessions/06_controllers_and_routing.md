# Chapter 6: Controllers & Routing

## Overview

Pimcore is **100% Symfony under the hood**, so controllers work exactly like Symfony controllers. You define routes with PHP attributes or YAML, inject services via Dependency Injection, and return responses (HTML via Twig or JSON).

---

## Project Structure

```
src/
├── Controller/
│   ├── ProductController.php       # Product list, detail, category, search
│   └── Api/
│       └── ProductApiController.php # JSON API endpoints
├── Service/
│   └── ProductService.php
└── Repository/
    └── ProductRepository.php
```

---

## 1. Product List Controller

**File:** `src/Controller/ProductController.php`

```php
<?php
namespace App\Controller;

use Pimcore\Model\DataObject\Product;
use Pimcore\Model\DataObject\Product\Listing;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

class ProductController extends AbstractController
{
    #[Route('/products', name: 'product_list')]
    public function list(): Response
    {
        $products = new Listing();
        $products->setOrderKey('name');
        $products->setOrder('asc');

        return $this->render('product/list.html.twig', [
            'products' => $products,
        ]);
    }
}
```

**Template:** `templates/product/list.html.twig`

```twig
{% extends "base.html.twig" %}

{% block title %}Product Catalog{% endblock %}

{% block content %}
<h1>Product Catalog</h1>
<ul>
    {% for product in products %}
        <li>
            <a href="{{ path('product_detail', {id: product.id}) }}">
                <strong>{{ product.name }}</strong>
            </a>
            - ${{ product.price }}
        </li>
    {% else %}
        <li>No products found.</li>
    {% endfor %}
</ul>
{% endblock %}
```

---

## 2. Product Detail Controller

```php
#[Route('/product/{id}', name: 'product_detail')]
public function detail(int $id): Response
{
    $product = Product::getById($id);

    if (!$product) {
        throw $this->createNotFoundException("Product with ID $id not found.");
    }

    return $this->render('product/detail.html.twig', [
        'product' => $product,
        'category' => $product->getCategory(),
    ]);
}
```

**Template:** `templates/product/detail.html.twig`

```twig
{% extends "base.html.twig" %}

{% block content %}
<h1>{{ product.name }}</h1>
<p><strong>Price:</strong> ${{ product.price }}</p>
<div><strong>Description:</strong> {{ product.description|raw }}</div>

{% if product.image %}
    <img src="{{ product.image.getThumbnail('galleryThumbnail') }}" alt="{{ product.name }}">
{% endif %}

{% if category %}
    <hr>
    <h2>Category: {{ category.name }}</h2>
{% endif %}

<a href="{{ path('product_list') }}">← Back to catalog</a>
{% endblock %}
```

---

## 3. Category Filtering

```php
#[Route('/category/{id}', name: 'products_by_category')]
public function byCategory(int $id): Response
{
    $category = \Pimcore\Model\DataObject\Category::getById($id);

    if (!$category) {
        throw $this->createNotFoundException("Category not found.");
    }

    $products = new \Pimcore\Model\DataObject\Product\Listing();
    $products->setCondition('category__id = ?', [$category->getId()]);
    $products->setOrderKey('name');
    $products->setOrder('asc');

    return $this->render('product/by_category.html.twig', [
        'category' => $category,
        'products' => $products,
    ]);
}
```

---

## 4. Search Controller

```php
#[Route('/products/search', name: 'product_search')]
public function search(Request $request): Response
{
    $query = $request->query->get('q');
    $products = [];

    if ($query) {
        $products = new Listing();
        $products->setCondition('name LIKE ?', ["%$query%"]);
        $products->setOrderKey('name');
        $products->setOrder('asc');
    }

    return $this->render('product/search.html.twig', [
        'products' => $products,
        'query'    => $query,
    ]);
}
```

**Template search form:**
```twig
<form method="get" action="{{ path('product_search') }}">
    <input type="text" name="q" value="{{ query }}" placeholder="Search by name">
    <button type="submit">Search</button>
</form>
```

---

## 5. Pagination

```php
#[Route('/products', name: 'product_list')]
public function list(Request $request): Response
{
    $page    = max(1, (int) $request->query->get('page', 1));
    $perPage = 10;
    $offset  = ($page - 1) * $perPage;

    $products = new Listing();
    $products->setLimit($perPage);
    $products->setOffset($offset);
    $products->setOrderKey('name');
    $products->setOrder('asc');

    $totalPages = ceil($products->getTotalCount() / $perPage);

    return $this->render('product/list.html.twig', [
        'products'   => $products,
        'page'       => $page,
        'totalPages' => $totalPages,
    ]);
}
```

**Pagination in Twig:**
```twig
<div class="pagination">
    {% if page > 1 %}
        <a href="{{ path('product_list', {page: page - 1}) }}">← Previous</a>
    {% endif %}

    <span>Page {{ page }} of {{ totalPages }}</span>

    {% if page < totalPages %}
        <a href="{{ path('product_list', {page: page + 1}) }}">Next →</a>
    {% endif %}
</div>
```

---

## 6. Card Grid Layout (TailwindCSS)

**Template:** `templates/product/list.html.twig`

```twig
<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
    {% for product in products %}
        <div class="bg-white rounded-2xl shadow-md overflow-hidden hover:shadow-lg transition">
            {% if product.image %}
                <img src="{{ product.image.getThumbnail('galleryThumbnail') }}"
                     alt="{{ product.name }}"
                     class="w-full h-48 object-cover">
            {% else %}
                <div class="w-full h-48 bg-gray-200 flex items-center justify-center text-gray-500">
                    No Image
                </div>
            {% endif %}

            <div class="p-4">
                <h2 class="text-lg font-semibold">{{ product.name }}</h2>
                <p class="text-gray-600">${{ product.price }}</p>

                {% if product.category %}
                    <p class="text-sm text-gray-500 mt-1">
                        Category: <span class="font-medium">{{ product.category.name }}</span>
                    </p>
                {% endif %}

                <a href="{{ path('product_detail', {id: product.id}) }}"
                   class="mt-3 inline-block bg-blue-600 text-white px-3 py-1 rounded-lg text-sm hover:bg-blue-700">
                    View Details
                </a>
            </div>
        </div>
    {% endfor %}
</div>
```

---

## 7. Injecting a Service

```php
// src/Service/ProductService.php
class ProductService
{
    public function getFeaturedProducts(int $limit = 5): array
    {
        $listing = new \Pimcore\Model\DataObject\Product\Listing();
        $listing->setLimit($limit);
        $listing->setOrderKey('price');
        $listing->setOrder('DESC');
        return $listing->load();
    }
}

// In Controller:
public function featured(ProductService $productService): Response
{
    return $this->json($productService->getFeaturedProducts());
}
```

---

## Routes Summary

| URL | Route Name | Description |
|---|---|---|
| `/products` | `product_list` | All products (paginated) |
| `/product/{id}` | `product_detail` | Single product detail |
| `/category/{id}` | `products_by_category` | Products filtered by category |
| `/products/search?q=...` | `product_search` | Search by product name |
| `/products/sort/{field}` | `product_sort` | Sort products by field |
| `/api/products` | `api_products` | JSON API (see Chapter 8) |

---

## Connecting a Route to a Pimcore Document

To make a controller power a Pimcore Document page:

1. Go to **Documents** → Create Page → Name it `products`
2. Open the page → **Settings** tab → set:
   - **Controller:** `App\Controller\ProductController::list`
   - **Template:** `product/list.html.twig`
3. Publish

Now `/products` loads through both the Pimcore Document tree and your Symfony controller.
