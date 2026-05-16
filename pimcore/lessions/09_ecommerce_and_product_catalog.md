# Chapter 9: eCommerce & Product Catalog

## What is Pimcore's eCommerce Framework?

Pimcore's eCommerce Framework is **not** a ready-made shop like Shopify. Instead, it's a **flexible toolkit** to build custom e-commerce solutions:

- Model products in Data Objects (SKU, price, variants, categories)
- Organize them into a searchable, filterable catalog
- Add pricing, availability, cart & checkout logic
- Integrate with ERP, CRM, and payment systems

---

## Prerequisites

Before building a catalog, ensure you have:
- A `Product` class with fields: `name`, `sku`, `price`, `category`, `image`
- A `Category` class for organizing products
- At least 2–3 Product objects in the Data Objects tree

---

## Step 1: Install the eCommerce Bundle

```bash
composer require pimcore/pimcore-ecommerce-framework
```

Enable in `config/bundles.php` if not auto-registered.

---

## Step 2: Create Category Structure

1. Go to **Data Objects** → Categories
2. Add categories:
   - Electronics
   - Accessories
   - Home Appliances

## Step 3: Assign Products to Categories

In the Product objects, set the `category` relation:
- Laptop X1 → Category = **Electronics**
- Wireless Mouse → Category = **Accessories**

---

## Step 4: Configure the Product Index

Pimcore uses an **index service** (similar to Elasticsearch but also works with MySQL) to make products searchable and filterable.

**File:** `config/packages/pimcore_ecommerce.yaml`

```yaml
pimcore_ecommerce_framework:
  index_services:
    mysql:
      type: MySql
      config:
        tablename: ecommerce_index
  workers:
    default:
      type: default
```

Bootstrap the index:
```bash
bin/console ecommerce:indexservice:bootstrap
```

---

## Step 5: Add Products to Index

In Pimcore Admin:
1. Go to **Marketing → E-Commerce → Index**
2. Select **Add Products to Index** → choose your Product objects

---

## Step 6: Create a Catalog Page

1. Go to **Documents** → Add Page → name it `catalog`

**Template:** `templates/catalog.html.twig`

```twig
{% set products = pimcore_ecommerce_index('default').getProducts() %}

<h1>Product Catalog</h1>
<div class="product-grid">
    {% for product in products %}
        <div class="product-card">
            <h2>{{ product.name }}</h2>
            <p>SKU: {{ product.sku }}</p>
            <p>Price: {{ product.price }} €</p>
            <img src="{{ product.image|pimcore_image({'thumbnail': 'catalog'}) }}" alt="{{ product.name }}">
            <a href="{{ path('product_detail', {id: product.id}) }}">View Details</a>
        </div>
    {% endfor %}
</div>
```

---

## Product Relationships

### Related Products (Upsells, Bundles)

Add to the `Product` class:
- `relatedProducts` → Many-to-Many Object Relation (to `Product`)
- `alternativeProducts` → Many-to-Many Object Relation (to `Product`)

Example in Admin:
- **Laptop X1** → Related Products: Wireless Mouse, Laptop Bag
- **iPhone 16** → Alternative Products: iPhone 16 Pro, iPhone 16 Plus

### Display in Templates

```twig
{% if product.relatedProducts %}
    <h3>Related Products</h3>
    {% for related in product.relatedProducts %}
        <a href="{{ path('product_detail', {id: related.id}) }}">
            {{ related.name }} — ${{ related.price }}
        </a>
    {% endfor %}
{% endif %}
```

---

## Advanced: Filter by Category in Catalog

```php
// Controller: filter by category
#[Route('/shop/{categorySlug}', name: 'catalog_by_category')]
public function catalog(string $categorySlug): Response
{
    // Find category by slug
    $categories = new \Pimcore\Model\DataObject\Category\Listing();
    $categories->setCondition('slug = ?', [$categorySlug]);
    $category = $categories->current();

    // Fetch products in this category
    $products = new \Pimcore\Model\DataObject\Product\Listing();
    $products->setCondition('category__id = ?', [$category->getId()]);

    return $this->render('catalog.html.twig', [
        'products' => $products,
        'category' => $category,
    ]);
}
```

---

## Cart & Checkout (Extending the Framework)

For a full e-commerce flow, you can extend with:

```php
// CartManager (Pimcore eCommerce)
$cartManager = $this->get('pimcore_ecommerce_framework.cart_manager');
$cart = $cartManager->getOrCreateCartByName('default');
$cart->addItem($product, 1);
$cart->save();
```

The framework provides:
- **CartManager** — manage sessions/persistent carts
- **PriceSystem** — flexible pricing rules (customer groups, promotions)
- **AvailabilitySystem** — stock management
- **Checkout** — multi-step order processing

---

## Summary

| Feature | How to Implement |
|---|---|
| Product listing | `Product\Listing()` + controller |
| Category filtering | `setCondition('category__id = ?')` |
| Product search | `setCondition('name LIKE ?')` |
| Product index | `pimcore_ecommerce_framework` config + bootstrap |
| Related products | Many-to-Many relation in class definition |
| Cart | `CartManager` from eCommerce Framework |
| Checkout | Custom Symfony controllers + CartManager |
| Integrations | ERP/CRM via Symfony services / events |

> **Next Steps:** Add faceted filtering (category, price range) using the Index Service, or integrate with a payment gateway via Symfony event subscribers.
