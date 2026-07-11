# PHP OOP Framework Design Book

A practical guide to PHP OOP concepts that teach you how to design frameworks like **Symfony**, **Laravel**, and **Pimcore**, with **Composer packages**, **bundles**, **service containers**, **facades**, and frontend integration using **JavaScript**, **React**, **Next.js**, and **TypeScript**.

This book is not only about PHP syntax. It is about learning how framework internals are designed.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_oop_foundation_for_frameworks.md) | **OOP Foundation for Frameworks** | Classes, objects, encapsulation, composition, lifecycle objects |
| [02](02_contracts_interfaces_and_solid.md) | **Contracts, Interfaces & SOLID** | Interfaces, abstract classes, dependency inversion, framework boundaries |
| [03](03_composer_autoloading_and_packages.md) | **Composer, Autoloading & Packages** | PSR-4, package structure, scripts, service discovery, reusable libraries |
| [04](04_service_container_and_dependency_injection.md) | **Service Container & Dependency Injection** | Binding, singleton, auto-wiring, service providers, tagged services |
| [05](05_kernel_router_middleware_controller.md) | **Kernel, Router, Middleware & Controller** | Front controller, HTTP kernel, route resolver, controller dispatcher |
| [06](06_bundles_modules_plugins.md) | **Bundles, Modules & Plugins** | Symfony bundles, Pimcore bundles, Laravel packages, extension points |
| [07](07_facades_events_commands.md) | **Facades, Events & Commands** | Laravel facades, event dispatcher, console commands, jobs |
| [08](08_ui_asset_integration.md) | **UI & Asset Integration** | Blade/Twig, Vite, npm, React widgets, Pimcore admin assets |
| [09](09_react_next_typescript_frontend.md) | **React, Next.js & TypeScript Frontends** | API-first backend, independent frontend, auth, DTOs, generated clients |
| [10](10_build_your_own_framework_roadmap.md) | **Build Your Own Framework Roadmap** | Step-by-step mini-framework plan |

---

## How to Read

**If you are new to OOP:** Read 1 → 2 first.

**If you want Symfony/Pimcore style:** Read 3 → 4 → 5 → 6 → 8.

**If you want Laravel style:** Read 3 → 4 → 5 → 7 → 8.

**If you want independent frontend architecture:** Read 8 → 9.

**If you want to build your own framework:** Read all chapters, then follow Chapter 10.

---

## Big Picture

```text
PHP Framework
├── public/index.php          front controller
├── Kernel                    bootstraps app and handles lifecycle
├── Service Container         creates and stores services
├── Router                    maps URL to controller
├── Middleware                wraps request/response flow
├── Controller                translates HTTP to application action
├── Events                    decouple side effects
├── Console Commands          CLI tools like migrate/cache/make
├── Bundles/Packages          reusable modules
├── Composer                  autoloading and dependency management
└── UI Integration
    ├── Blade/Twig templates
    ├── Vite/npm assets
    ├── React widgets inside PHP pages
    └── separate Next.js/TypeScript frontend
```

---

## Core Design Rule

Keep your **business logic independent**.

```text
Good:
Controller -> Application Service -> Domain Object -> Repository Interface

Bad:
Controller -> SQL + RabbitMQ + HTML + API calls + business rules in one method
```

Frameworks are powerful when they separate concerns cleanly.

