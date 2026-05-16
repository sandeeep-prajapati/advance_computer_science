# Chapter 17: Deployment & CI/CD

## Overview

Pimcore deployments require careful orchestration: code, class definitions, database migrations, caches, and assets all need to be handled correctly. This chapter covers a safe, repeatable deployment workflow.

---

## 1. Environment Variables

Never hardcode credentials. Store environment-specific values in `.env.local` (never commit to Git).

**`.env` (committed — safe defaults):**
```dotenv
APP_ENV=prod
APP_DEBUG=0
DATABASE_URL="mysql://user:pass@mysql:3306/pimcore_prod"
PIMCORE_ENVIRONMENT=prod
```

**`.env.local` (NOT committed — overrides per server):**
```dotenv
DATABASE_URL="mysql://real_user:real_pass@mysql:3306/pimcore_prod"
PIMCORE_IMAGE_PROXY_SECRET=xxxx
REDIS_URL=redis://redis:6379
```

> **In Docker/Kubernetes:** Inject env vars through **Secrets**, not `.env` files.

---

## 2. Database Changes

Pimcore stores data structure in **Class Definitions** (JSON files in `var/classes` or `config/pimcore/classes`).

**Rules:**
1. Always commit class JSON files to Git
2. Use `classes-rebuild` to sync DB schema after deploying class changes
3. Use Doctrine migrations for custom tables/entities

```bash
# After deploying new class definitions
bin/console pimcore:deployment:classes-rebuild --create-classes --delete-removed-classes

# Run Doctrine migrations for custom tables
bin/console doctrine:migrations:migrate --no-interaction --env=prod
```

> **Golden rule:** No manual DB changes. Everything through migrations or class definitions.

---

## 3. Assets (Files, Images, Documents)

- Assets live in `/var/assets` (or external storage)
- In CI/CD: **don't overwrite!** Use shared storage (NFS, S3, Azure Blob)
- For large deployments: serve assets via CDN

```bash
# If using S3:
# Configure in config/packages/flysystem.yaml
# Use pimcore/cloud-storage bundle
```

---

## 4. Full Deployment Sequence

Run these commands in order on every deploy:

```bash
#!/bin/bash
set -e

echo "=== Step 1: Install Composer dependencies ==="
composer install --no-dev --optimize-autoloader --prefer-dist

echo "=== Step 2: Sync class definitions to DB ==="
bin/console pimcore:deployment:classes-rebuild --create-classes --delete-removed-classes

echo "=== Step 3: Run database migrations ==="
bin/console doctrine:migrations:migrate --no-interaction --env=prod

echo "=== Step 4: Clear and warm up caches ==="
bin/console cache:clear --env=prod
bin/console cache:warmup --env=prod
bin/console pimcore:cache:clear --all

echo "=== Step 5: Install bundle assets (if custom bundles added) ==="
bin/console assets:install public --symlink --env=prod

echo "=== Done ==="
```

---

## 5. Git Branching Strategy

```
main          →  production (protected, deployment by tag/release)
develop       →  staging (auto-deploys on push)
feature/xxx   →  development (merge into develop via PR)
hotfix/xxx    →  emergency fix (merge into main AND develop)
```

---

## 6. CI/CD Pipeline (GitHub Actions)

**.github/workflows/deploy.yml**

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [develop, main]
  pull_request:
    branches: [develop]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:8
        env:
          MYSQL_DATABASE: pimcore_test
          MYSQL_USER: root
          MYSQL_ALLOW_EMPTY_PASSWORD: yes
        ports:
          - 3306:3306

    steps:
      - uses: actions/checkout@v4

      - name: Setup PHP
        uses: shivammathur/setup-php@v2
        with:
          php-version: '8.3'
          extensions: mbstring, intl, pdo_mysql

      - name: Install dependencies
        run: composer install --prefer-dist --no-progress

      - name: Run tests
        run: bin/phpunit --testdox

  deploy-staging:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/develop'
    steps:
      - uses: actions/checkout@v4

      - name: Deploy to staging
        uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.STAGING_HOST }}
          username: deploy
          key: ${{ secrets.STAGING_SSH_KEY }}
          script: |
            cd /var/www/pimcore-staging
            git pull origin develop
            composer install --no-dev --optimize-autoloader
            bin/console pimcore:deployment:classes-rebuild -c -d
            bin/console doctrine:migrations:migrate --no-interaction --env=prod
            bin/console cache:clear --env=prod
            bin/console pimcore:cache:clear --all

  deploy-production:
    needs: test
    runs-on: ubuntu-latest
    if: startsWith(github.ref, 'refs/tags/v')
    steps:
      - uses: actions/checkout@v4

      - name: Deploy to production
        uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.PROD_HOST }}
          username: deploy
          key: ${{ secrets.PROD_SSH_KEY }}
          script: |
            cd /var/www/pimcore-prod

            # Backup DB before migration
            mysqldump -u pimcore -p$DB_PASS pimcore_prod > /backups/pimcore_$(date +%Y%m%d_%H%M%S).sql

            git pull origin main
            composer install --no-dev --optimize-autoloader
            bin/console pimcore:deployment:classes-rebuild -c -d
            bin/console doctrine:migrations:migrate --no-interaction --env=prod
            bin/console cache:clear --env=prod
            bin/console cache:warmup --env=prod
            bin/console pimcore:cache:clear --all
```

---

## 7. Zero-Downtime Deployment with Symlinks

For production deployments that must have zero downtime:

```bash
#!/bin/bash
# deploy-zero-downtime.sh

DEPLOY_DIR="/var/www/releases/$(date +%Y%m%d_%H%M%S)"
CURRENT_LINK="/var/www/pimcore"
SHARED_DIR="/var/www/shared"

# Create new release directory
mkdir -p "$DEPLOY_DIR"
git clone --depth 1 git@github.com:company/pimcore.git "$DEPLOY_DIR"

cd "$DEPLOY_DIR"

# Symlink shared directories (assets, var/assets persist across deploys)
ln -nfs "$SHARED_DIR/var/assets" "$DEPLOY_DIR/var/assets"
ln -nfs "$SHARED_DIR/.env.local" "$DEPLOY_DIR/.env.local"

# Install dependencies
composer install --no-dev --optimize-autoloader

# Run migrations (before switching symlink)
bin/console pimcore:deployment:classes-rebuild -c -d
bin/console doctrine:migrations:migrate --no-interaction --env=prod
bin/console cache:clear --env=prod
bin/console cache:warmup --env=prod

# Switch symlink atomically (near-zero downtime)
ln -nfs "$DEPLOY_DIR" "$CURRENT_LINK"

# Reload PHP-FPM
sudo systemctl reload php8.3-fpm

echo "Deployed to $DEPLOY_DIR"
```

---

## 8. Extra Safety Measures

| Practice | Command / Config |
|---|---|
| DB backup before migration | `mysqldump` before `doctrine:migrations:migrate` |
| Verify migration before running | `bin/console doctrine:migrations:status` |
| Feature flags | Environment variables or DB-toggle for risky features |
| Health check endpoint | Custom route returning 200 + app version |
| Rollback plan | Keep last 3 releases, symlink back if deploy fails |
| Audit logs | Keep Pimcore admin logs enabled |

---

## Deployment in One Sentence

> Commit code + class JSONs → push to develop → CI runs tests → deploy to staging → QA review → merge/tag release → deploy to production with composer install + classes-rebuild + migrations + cache clear, keeping assets in shared storage.
