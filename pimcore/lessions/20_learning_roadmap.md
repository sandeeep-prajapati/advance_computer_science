# Chapter 20: Learning Roadmap — 5 Weeks to Pimcore Proficiency

## Overview

This is a structured, hands-on 5-week plan that takes you from zero → a real open-source contribution, with concrete daily tasks and clear success criteria.

---

## Week 1: Install & Basics

**Goal:** Have Pimcore running locally, know the Admin UI, commit first project code.

| Day | Task |
|---|---|
| Day 1–2 | Install Pimcore with Docker + Composer. Access Admin UI at `/admin`. Explore Data Objects, Assets, Documents sections. |
| Day 3–4 | Create a simple Product class (name, price, sku). Add a few objects manually. |
| Day 5 | Explore Assets — upload images and PDFs. Try linking assets to your product. |
| Day 6–7 | Play with Documents — create a CMS page, add text/images. Read Pimcore docs on Data Objects & Assets. |

**Success Criteria:**
- Pimcore runs locally
- Can create/update a Product object with an image via Admin UI
- First Git commit: project bootstrapped

---

## Week 2: Classes & Assets

**Goal:** Model richer data and handle media programmatically.

| Day | Task |
|---|---|
| Day 1–2 | Add fields: description, categories (relation), images (asset relation). Add localized fields for multi-language. |
| Day 3 | Write a PHP script to create/update Product objects programmatically. |
| Day 4 | Import 5 products from a small CSV using `fgetcsv` or `League\Csv`. |
| Day 5 | Define a thumbnail preset in Admin UI (Settings → Thumbnails). |
| Day 6 | Create a thumbnail variant programmatically via PHP. |
| Day 7 | Document your Product data model (draw a diagram + notes). |

**Success Criteria:**
- Product model with relations + localized fields exists
- CSV import script works
- Thumbnail variant displays correctly

---

## Week 3: APIs & Import/Export

**Goal:** Expose and consume Pimcore data via APIs.

| Day | Task |
|---|---|
| Day 1–2 | Create a Symfony controller returning Products as JSON. Add pagination + filtering (`?page=2&category=shoes`). |
| Day 3 | Test your API with curl and Postman. |
| Day 4–5 | Enable Pimcore DataHub/GraphQL bundle. Write a GraphQL query to fetch products with images/categories. |
| Day 6 | Write an export command: Products → CSV. |
| Day 7 | Combine: fetch via API → save to CSV → re-import. |

**Success Criteria:**
- REST endpoint works with pagination & filters
- GraphQL query returns product + image data
- Export + re-import cycle proven

---

## Week 4: Admin Extensions & Testing

**Goal:** Extend Admin UI and write tests.

| Day | Task |
|---|---|
| Day 1–2 | Add a custom button in Product grid (e.g., "Export CSV"). |
| Day 3–4 | Build a small bundle with a custom admin page (e.g., product stats). |
| Day 5 | Set up PHPUnit config. Write a unit test for a simple service. |
| Day 6–7 | Write a functional test for your Product API endpoint. |

**Success Criteria:**
- Custom button appears in Product grid
- Bundle loads a custom admin page
- PHPUnit runs with passing unit + functional tests

---

## Week 5: Kubernetes & Contribution

**Goal:** Deploy to Kubernetes and submit your first Pimcore PR.

| Day | Task |
|---|---|
| Day 1 | Set up a local Kubernetes cluster (Minikube or k3s). Deploy Pimcore using the manifests from Chapter 18. |
| Day 2 | Practice `kubectl logs`, `kubectl exec`, `kubectl describe pod` for debugging. |
| Day 3 | Read Pimcore's Contributing Guidelines. Find a "good first issue" on GitHub. |
| Day 4–5 | Fork Pimcore repo. Fix the issue. Run the full test suite. |
| Day 6 | Open a PR against Pimcore's main dev branch. Follow the PR template. |
| Day 7 | Respond to maintainer feedback if any. |

**Success Criteria:**
- Pimcore running in Kubernetes locally
- Can debug pods and view logs
- PR open in Pimcore core repo with all CI tests passing

---

## Summary Timeline

| Week | Focus | Deliverable |
|---|---|---|
| Week 1 | Setup & Admin UI basics | Pimcore running + Product object via Admin UI |
| Week 2 | Data modeling & Assets | CSV import + thumbnails working |
| Week 3 | APIs (REST + GraphQL) | Live API + export/import cycle |
| Week 4 | Admin Extensions + Tests | Custom UI + PHPUnit passing |
| Week 5 | Kubernetes + Contribution | K8s deployment + PR submitted |

---

## Quick Reference: Most-Used Commands

```bash
# Development
docker compose up -d
docker compose exec php bin/console cache:clear
docker compose exec php bin/console pimcore:deployment:classes-rebuild -c
docker compose exec php bin/phpunit

# Kubernetes
kubectl get pods -n pimcore
kubectl logs -f <pod> -n pimcore -c php-fpm
kubectl exec -it <pod> -n pimcore -c php-fpm -- bash
kubectl rollout restart deployment/pimcore-app -n pimcore

# Pimcore Console
bin/console debug:router
bin/console debug:container
bin/console pimcore:system:requirements:check
bin/console app:import-products

# Git workflow
git checkout -b feature/my-feature
git fetch upstream && git rebase upstream/10.x
```

---

## Recommended Reading

| Resource | URL |
|---|---|
| Pimcore Docs | `docs.pimcore.com` |
| Pimcore GitHub | `github.com/pimcore/pimcore` |
| Symfony Docs | `symfony.com/doc` |
| Pimcore Community | `pimcore.com/community` |
| Docker Docs | `docs.docker.com` |
| Kubernetes Docs | `kubernetes.io/docs` |
