# Chapter 1: Introduction to Pimcore

## What is Pimcore?

Pimcore is an **open-source Digital Experience Platform (DXP)** built on the **Symfony PHP framework**. It is a single, unified platform that helps businesses manage, enrich, and deliver digital data across multiple channels — websites, mobile apps, marketplaces, print, and IoT devices.

Instead of running five separate tools (one for product data, one for images, one for content, one for eCommerce), Pimcore brings everything under one roof.

> **Analogy:** Think of Pimcore as a "digital warehouse + factory":
> - **Warehouse shelves** (Data Objects & Assets) store your products, images, and content.
> - **Blueprints** (Classes) define how things are structured.
> - **Factory machines** (Admin UI + APIs) let you organize, enrich, and distribute data anywhere.

---

## Main Modules

### 1. PIM — Product Information Management
- Central hub for all product-related data.
- Manages attributes (size, color, material, SKU), variants, categories, and relationships.
- Ensures consistent product data across websites, apps, and marketplaces.

### 2. DAM — Digital Asset Management
- Central repository for all media assets (images, videos, PDFs, etc.).
- Provides metadata management, version control, tagging, and rights management.
- Assets can be directly linked to products, categories, or content pages.

### 3. CMS — Content Management System
- Manages website content (pages, blogs, landing pages).
- WYSIWYG editing and drag-and-drop layouts.
- Supports both traditional CMS and **headless CMS** approaches.

### 4. MDM — Master Data Management
- Extends PIM to handle all types of data (customers, suppliers, vendors, locations).
- Ensures a single source of truth for critical business data.

### 5. eCommerce Framework
- Flexible toolkit to create product catalogs, carts, and checkout processes.
- Can integrate with ERP, CRM, and payment systems.
- Supports personalized shopping experiences.

---

## Core Building Blocks

| Building Block | What It Is | Example |
|---|---|---|
| **Data Objects** | Structured business entities | A `Product` with name, SKU, price |
| **Classes** | Schema/blueprint that defines Object structure | The `Product` class with field definitions |
| **Assets** | Raw files (images, videos, PDFs) | `laptop.jpg`, `manual.pdf` |
| **Documents** | CMS pages, snippets, email templates | `/about-us` page, product detail page |
| **Admin UI** | ExtJS-powered backend for editors and admins | The `/admin` dashboard |

> **Rule of thumb:**
> - Objects = structured data
> - Assets = media files
> - Documents = content pages

---

## Problems Pimcore Solves

| Problem | Solution |
|---|---|
| Data chaos | Unifies scattered product, asset, and content data |
| Content silos | Single-source publishing across websites, apps, marketplaces, and print |
| Integration pain | API-first and Symfony-based, plugs into modern stacks |
| Scalability | Handles millions of products, files, or pages efficiently |

---

## Architecture Overview

```
+-----------------------------+
|       Admin UI (ExtJS)      |
|  (Editors / Data Managers)  |
+-------------+---------------+
              |
              v
+-------------+---------------+
|    Application Layer         |
| (Symfony Controllers,        |
| Services, Bundles, API)      |
+---+----------+-----------+--+
    |          |           |
    v          v           v
+-------+  +-------+  +----------+
| MySQL |  | Files |  | Elastic  |
| / DB  |  | / S3  |  | Search   |
+-------+  +-------+  +----------+
              |
              v
+-----------------------------+
|    Cache (Redis / Varnish)  |
+-----------------------------+
              |
              v
+-----------------------------+
|  APIs & Frontends           |
|  (REST, GraphQL, Twig,      |
|   React, Vue, Next.js)      |
+-----------------------------+
```

---

## Real-World Use Cases

- **PIM:** Manage product catalogs, specs, attributes, multi-language data.
- **MDM:** Single source of truth for enterprise data.
- **DAM:** Organize, transform, and deliver digital assets.
- **CMS / Headless CMS:** Manage websites or deliver structured content to apps via APIs.
- **eCommerce:** Extend with product management, carts, pricing, personalization.

---

## Key Technologies Under the Hood

- **PHP 8.1+** — runtime
- **Symfony** — framework (routing, DI, events, security)
- **MySQL / MariaDB** — primary data store
- **Redis** — caching and sessions (recommended)
- **Elasticsearch / OpenSearch** — search and product indexing (optional but critical for large catalogs)
- **Twig** — templating engine
- **ExtJS** — Admin UI framework
- **Composer** — dependency management
