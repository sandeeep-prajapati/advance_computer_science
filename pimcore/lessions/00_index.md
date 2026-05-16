# Pimcore Developer Book

A comprehensive, chapter-wise guide to building production-grade applications with Pimcore — from first install to Kubernetes deployment and open-source contribution.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_introduction.md) | **Introduction to Pimcore** | What is Pimcore, modules (PIM/DAM/CMS/MDM/eCommerce), core building blocks, architecture |
| [02](02_setup_and_installation.md) | **Setup & Installation** | System requirements, Composer install, Docker setup, database setup, common pitfalls |
| [03](03_data_objects_and_classes.md) | **Data Objects & Classes** | Classes, Data Objects, Object Bricks, Field Collections, Localized Fields, Relations, PHP API |
| [04](04_assets_and_dam.md) | **Assets & Digital Asset Management** | Upload, metadata, tagging, image thumbnails, connecting assets to objects, CDN strategy |
| [05](05_cms_documents.md) | **CMS Documents** | Pages, Snippets, Emails, editable areas in Twig, organizing content, headless vs. traditional |
| [06](06_controllers_and_routing.md) | **Controllers & Routing** | Symfony controllers in Pimcore, product list/detail/category/search, pagination, card grids |
| [07](07_rest_api_and_graphql.md) | **REST API & GraphQL** | Custom REST endpoints, pagination & filtering, DataHub GraphQL, comparing approaches |
| [08](08_workflows_and_localization.md) | **Workflows & Localization** | Workflow states/transitions/roles, YAML config, multilingual fields, language trees |
| [09](09_ecommerce_and_product_catalog.md) | **eCommerce & Product Catalog** | eCommerce framework, product index, catalog pages, related products, cart & checkout |
| [10](10_bundles_and_extensions.md) | **Bundles & Extensions** | Custom bundles, bundle structure, services & controllers, Admin UI extensions (ExtJS) |
| [11](11_event_subscribers_and_hooks.md) | **Event Subscribers & Hooks** | Pimcore events, auto-slug generation, audit logging, external sync on save |
| [12](12_data_import_export.md) | **Data Import & Export** | CSV import command, image downloading, export to CSV, dry-run mode, error handling |
| [13](13_performance_and_caching.md) | **Performance & Caching** | OPcache, Redis, Varnish/Nginx full-page cache, MySQL tuning, Elasticsearch, profiling |
| [14](14_security_and_permissions.md) | **Security & Permissions** | Users & roles, object/asset permissions, Symfony firewall, JWT, auditing |
| [15](15_debugging_and_troubleshooting.md) | **Debugging & Troubleshooting** | Log files, Symfony profiler, SQL debugging, common errors, Xdebug, debug commands |
| [16](16_testing.md) | **Testing** | PHPUnit setup, unit tests, functional controller tests, CI/CD with GitHub Actions |
| [17](17_deployment_and_cicd.md) | **Deployment & CI/CD** | Env vars, class definitions, migration sequence, zero-downtime deploy, GitHub Actions |
| [18](18_kubernetes_setup_and_debugging.md) | **Kubernetes Setup & Debugging** | K8s architecture, Dockerfiles, manifests (Deployments/Services/PVC/Ingress), HPA, debugging pods |
| [19](19_contributing_to_pimcore.md) | **Contributing to Pimcore** | Finding issues, fork & branch, PR process, coding style |
| [20](20_learning_roadmap.md) | **Learning Roadmap** | 5-week hands-on plan from zero to PR contribution |

---

## How to Read This Book

**If you're brand new to Pimcore:** Read chapters 1 → 2 → 3 → 4 → 5 in order. This gives you the core mental model.

**If you're building a product catalog:** Read 3 → 4 → 6 → 9 → 7.

**If you're adding APIs:** Read 7 → 8 → 11.

**If you're going to production:** Read 13 → 14 → 15 → 17 → 18.

**If you want to extend Pimcore:** Read 10 → 11 → 12.

---

## Key Concepts at a Glance

```
Pimcore
├── Data Layer
│   ├── Classes (blueprints/schemas)
│   ├── Data Objects (records/instances)
│   ├── Assets (images, videos, PDFs)
│   └── Documents (CMS pages, snippets, emails)
│
├── Application Layer (Symfony)
│   ├── Controllers (routing + responses)
│   ├── Services (business logic)
│   ├── Bundles (reusable feature packages)
│   └── Event Subscribers (hooks)
│
├── API Layer
│   ├── Custom REST (Symfony controllers)
│   └── GraphQL (DataHub bundle)
│
└── Infrastructure
    ├── MySQL (data store)
    ├── Redis (cache + sessions)
    ├── Elasticsearch (search + product index)
    └── Kubernetes (container orchestration)
```
