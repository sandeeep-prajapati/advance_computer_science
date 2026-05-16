# Chapter 5: CMS Documents

## What are Documents?

**Documents** are the CMS layer of Pimcore. They represent web pages, reusable content blocks, and email templates. Stored in the **Documents** tree in the Admin UI, they form your website's page hierarchy.

### Three Types of Documents

| Type | Use Case | Example |
|---|---|---|
| **Page** | Full CMS web page | Homepage, About Us, Product Detail |
| **Snippet** | Reusable content block | Footer, Promo Banner, Newsletter Signup |
| **Email** | Email templates | Order confirmation, Marketing newsletter |

---

## Step 1: Create a New Page

1. In Admin Panel → **Documents** (left sidebar)
2. Right-click on root (or a folder) → **Add Document** → **Page**
3. Enter details:
   - Name: `about-us`
   - Title: "About Us"
   - Template: `default.html.twig`
4. Save & Publish

> The page is now accessible at `/about-us`

---

## Step 2: Define Editable Areas in Twig

Pimcore's power is the **editable area** system — editors can change content directly in the Admin UI without touching code.

### Example: `templates/default.html.twig`

```twig
{% extends "layout.html.twig" %}

{% block content %}
    <h1>{% pimcore_input('headline', {'width': 600}) %}</h1>

    <div>
        {% pimcore_wysiwyg('content', {'height': 300}) %}
    </div>

    <div>
        {% pimcore_image('mainImage', {'thumbnail': 'content'}) %}
    </div>
{% endblock %}
```

### Available Editable Types

| Tag | Type | Use Case |
|---|---|---|
| `pimcore_input` | Single-line text | Headlines, titles |
| `pimcore_wysiwyg` | Rich text editor | Body content |
| `pimcore_image` | Image from DAM | Hero images, product photos |
| `pimcore_link` | Link element | CTAs, navigation links |
| `pimcore_snippet` | Include a Snippet | Reusable blocks (footer, banner) |
| `pimcore_select` | Dropdown selection | Layout options |
| `pimcore_areablock` | Drag-and-drop area blocks | Complex page builders |

---

## Step 3: Edit Content in Admin

1. Go to **Documents** → About Us → **Edit**
2. You'll see the WYSIWYG editor and editable fields (Headline, Content, Image)
3. Enter text in the headline field
4. Click into the WYSIWYG area to type content
5. Drag & drop an image from Assets into the image area
6. Click **Save & Publish**

> Non-technical editors can now update text & images without touching code.

---

## Step 4: Organize Content Pages

### Folder Structure Best Practice

```
/documents
  /en
    /home
    /about-us
    /products
      /laptop-x1
      /wireless-mouse
  /de
    /startseite
    /ueber-uns
    /produkte
  /fr
    /accueil
    /a-propos
```

### Use Snippets for Reusable Blocks

- Create a **Snippet** document for your newsletter signup box
- Include it in any Page template using `{% pimcore_snippet('newsletter') %}`
- Update the Snippet once → all pages using it update automatically

---

## Step 5: Connecting Documents to Data Objects

The real power of Pimcore CMS is combining **Documents** (pages) with **Data Objects** (structured data).

### Example: Dynamic Product Detail Page

**Template:** `templates/product/detail.html.twig`
```twig
{% extends "layout.html.twig" %}

{% block content %}
    <h1>{{ product.name }}</h1>
    <p>Price: ${{ product.price }}</p>
    <div>{{ product.description|raw }}</div>

    {% if product.image %}
        <img src="{{ product.image.getThumbnail('galleryThumbnail') }}" alt="{{ product.name }}">
    {% endif %}
{% endblock %}
```

**Controller Route:**
```php
// In Admin: Documents → product-detail page → Settings → Controller: ProductController::detail
```

This way, the same template renders different product data depending on which product is loaded.

---

## CMS Page vs. Headless Approach

| Approach | When to Use |
|---|---|
| **Traditional CMS** (Twig templates) | Classic websites, blog, marketing pages |
| **Headless CMS** (API-first) | React/Vue/Next.js frontends, mobile apps, omnichannel delivery |

Pimcore supports both. You can start with traditional CMS and add headless APIs later — without changing your data models.

---

## Summary

- **Documents** = the CMS layer (pages, snippets, emails)
- **Editable areas** in Twig = let editors manage content without code
- **Folder structure** = organize by language and content type
- **Snippets** = reuse common blocks across pages
- **Data Object integration** = documents can load and display dynamic structured data
