# Chapter 3: Data Objects & Classes

## Overview

Data Objects and Classes are the backbone of Pimcore's data modeling system. Think of it like a database but with a visual editor:

- **Class** = the schema/blueprint (like a table definition)
- **Data Object** = the actual record (like a row in that table)

---

## Core Concepts

### 1. Classes
A **Class** defines the structure (fields and layout) of a Data Object type.

- Created in Admin UI → Settings → Data Objects → Classes
- Each class generates a PHP model automatically (e.g., `Pimcore\Model\DataObject\Product`)
- Example: A `Product` class with fields: `name`, `sku`, `price`, `description`, `image`, `categories`

### 2. Data Objects
**Data Objects** are instances created from a Class — the actual content/records.

- Stored in the Data Objects tree (left sidebar in Admin)
- Example: A `Product` object = "iPhone 16 Pro" with SKU `IP16P-256`, price `₹120,000`

### 3. Object Bricks
**Object Bricks** are reusable sets of optional fields that can be attached to Classes dynamically.

- Use when: optional or type-specific attributes that not every object has
- Example: An `ElectronicsSpecs` brick with `batteryLife`, `screenSize` — only attached to electronic products, not clothing

### 4. Field Collections
**Field Collections** are groups of fields you can repeat multiple times within one object.

- Use when: you need repeatable structured data
- Example: A `TechnicalSpecs` field collection with `label` + `value`:
  - (label=Battery, value=24h)
  - (label=Weight, value=180g)

### 5. Localized Fields
**Localized Fields** are a special container for multi-language content.

- Use when: you need fields translated per language
- Example: `name[en]` = "iPhone 16", `name[de]` = "iPhone 16", `name[fr]` = "iPhone 16"

### 6. Relations
**Relations** link objects, assets, or documents to each other.

| Type | Use Case |
|---|---|
| Object Relation (Many-to-One) | Product → Category |
| Asset Relation | Product → Image |
| Many-to-Many Object Relation | Product → multiple Tags |
| Many-to-Many Asset Relation | Product → multiple Images |

---

## Creating a Product Class (Admin UI — Click by Click)

### Step 1: Open Class Editor
1. Log in → Settings (gear icon) → Data Objects → Classes
2. Click **Add** → Choose **Class** → Enter Name: `Product` → Click **Create**

### Step 2: Add Fields

Inside the class editor, add these fields:

| Field | Type | Notes |
|---|---|---|
| `sku` | Input | Mark as unique |
| `name` | Input | Product name |
| `price` | Numeric | Set min/max optional |
| `description` | WYSIWYG or Textarea | Rich text |
| `images` | Many-to-Many Asset Relation | Allow `image` type assets |
| `categories` | Many-to-Many Object Relation | Allow `Category` class |

> Tip: You can organize fields into **Panels**, **Fieldsets**, and **Tabs** for better Admin UI layout.

### Step 3: Save & Publish the Class
Click **Save & Publish** — Pimcore generates the PHP model and updates the database.

---

## Creating Product Objects

### Via Admin UI
1. Go to **Data Objects** in the sidebar
2. Right-click → **Add Object** → `Product`
3. Fill in fields (Name, SKU, Price, etc.)
4. Upload an image from Assets and link it
5. Click **Save & Publish**

### Via PHP (Programmatic)

```php
<?php
use Pimcore\Model\DataObject;
use Pimcore\Model\Asset;
use Pimcore\Model\Element\Service as ElementService;

// 1) Ensure parent folder exists
$parent = DataObject\Folder::getByPath('/products');
if (!$parent) {
    $parent = new DataObject\Folder();
    $parent->setKey('products');
    $parent->setParentId(1); // root
    $parent->save();
}

// 2) Create product object
$product = new DataObject\Product();
$product->setKey(ElementService::getValidKey('iPhone-16-pro', 'object'));
$product->setParentId($parent->getId());
$product->setPublished(true);

// 3) Set fields (setters follow field name in camelCase)
$product->setSku('IP16P-256');
$product->setName('iPhone 16 Pro');
$product->setPrice(120000);
$product->setDescription('Flagship phone — example description.');

// 4) Attach images (Many-to-Many asset relation = array of Assets)
$imageIds = [123, 124];
$assets = array_filter(array_map(fn($id) => Asset::getById($id), $imageIds));
$product->setImages(array_values($assets));

// 5) Attach categories (Many-to-Many object relation = array of Objects)
$categoryIds = [10, 11];
$cats = array_filter(array_map(fn($id) => DataObject\Category::getById($id), $categoryIds));
$product->setCategories(array_values($cats));

// 6) Save
$product->save(['versionNote' => 'Created via script']);
echo "Saved product ID: " . $product->getId();
```

---

## Deploying Classes Between Environments

The recommended workflow for moving class definitions between dev → staging → production:

```bash
# Export class definition to JSON (from Admin UI: Classes → select → Export JSON)
# OR via console:
php bin/console pimcore:definition:import:class /path/to/Product.json

# Rebuild classes & update DB schema
php bin/console pimcore:deployment:classes-rebuild --create-classes

# Clear cache
php bin/console cache:clear
```

> Always commit class JSON files to Git (`config/pimcore/classes/Product.json`).

---

## Common Pitfalls

| Problem | Fix |
|---|---|
| `Class not found` when using `DataObject\Product` | Run `pimcore:deployment:classes-rebuild --create-classes` |
| Setter not working | Setter names follow field name exactly in camelCase — `sku` → `setSku()` |
| Object saved in wrong location | Check `setParentId()` — `1` is the root |
| Single image vs. array | `Image` field type = single `Asset` instance; Many-to-Many = array of Assets |
| Field renamed, data lost | Renaming fields can lose stored data — be careful |

---

## Data Modeling Walkthrough: Products with Variants

A complete example showing all building blocks together:

1. **Create Class:** `Product` → fields: `sku`, `name`, `price`
2. **Add Localized Fields:** for `name` and `description` (multi-language)
3. **Add Object Relation:** to `Category` class
4. **Add Asset Relation:** for `images`
5. **Add Object Brick:** `ElectronicsSpecs` → `batteryLife`, `screenSize`
6. **Add Field Collection:** `TechnicalSpecs` → `label` + `value` pairs
7. **Create Data Objects:**
   - "iPhone 16" → has images + `ElectronicsSpecs`
   - "iPhone 16 Pro" → variant with higher price + extra specs

> **Classes define the structure. Data Objects are the actual records. Object Bricks, Field Collections, and Localized Fields make it flexible. Relations tie everything together.**
