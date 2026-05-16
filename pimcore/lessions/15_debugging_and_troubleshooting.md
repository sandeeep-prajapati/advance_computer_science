# Chapter 15: Debugging & Troubleshooting

## Overview

Pimcore is Symfony under the hood, so you use Symfony tools plus Pimcore-specific logs. The key tools are: **log files**, **Symfony profiler**, **SQL debugging**, and a set of **console commands** for common problems.

---

## 1. Log Files

Pimcore writes all logs to `var/log/` (per environment).

| Log File | Contents |
|---|---|
| `var/log/dev.log` | All runtime events in dev mode |
| `var/log/prod.log` | Runtime events in prod (less verbose) |
| `var/log/system.log` | Pimcore core system events |
| `var/log/admin-YYYY-MM-DD.log` | Admin actions (who saved/published what) |
| `var/log/php_error.log` | PHP errors (depending on PHP config) |

### Tail logs in real-time

```bash
# Watch dev log
tail -f var/log/dev.log

# Watch all logs
tail -f var/log/*.log

# Watch logs inside Docker
docker compose logs -f php
```

### Configure Log Level

```yaml
# config/packages/monolog.yaml
monolog:
  handlers:
    main:
      type: stream
      path: "%kernel.logs_dir%/%kernel.environment%.log"
      level: debug   # use 'info' in prod, 'debug' in dev
```

---

## 2. Symfony Profiler (Dev Mode)

The Symfony debug toolbar appears at the **bottom of every page** in dev mode.

### Enable if Missing

```yaml
# config/packages/dev/web_profiler.yaml
web_profiler:
  toolbar: true
  intercept_redirects: false
```

### What You Can See

- **Timeline** — controller execution, Twig rendering, event dispatching
- **SQL queries** — each query with time + EXPLAIN link
- **Cache** — hits and misses
- **Logs** — all log messages during the request
- **Security** — current user, roles, firewall

**URL:** `http://localhost/_profiler`

---

## 3. Debugging SQL Queries

### Enable Doctrine Query Logging

```bash
# In .env.local
APP_DEBUG=1
```

Or in `config/packages/doctrine.yaml`:
```yaml
doctrine:
  dbal:
    logging: true
```

### Run SQL Manually

```bash
# Query the object table
bin/console doctrine:query:sql "SELECT * FROM object_query_1 WHERE o_key='ProductX'"

# Check how many objects exist
bin/console doctrine:query:sql "SELECT COUNT(*) FROM objects WHERE o_type='object'"
```

### MySQL Slow Query Log

```ini
# /etc/mysql/my.cnf
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 1    ; log queries taking > 1 second
```

---

## 4. Common Errors & Fixes

### Class Not Found (`DataObject\Product`)

**Symptom:** `Class 'Pimcore\Model\DataObject\Product' not found`

**Fix:**
```bash
bin/console pimcore:deployment:classes-rebuild --create-classes
composer dump-autoload --optimize
bin/console cache:clear
```

### 404 After Deployment

**Symptom:** Routes return 404 after deploying to a new server.

**Checklist:**
```bash
# 1. Rebuild class definitions
bin/console pimcore:deployment:classes-rebuild -c -d

# 2. Clear all caches
bin/console cache:clear
bin/console pimcore:cache:clear --all

# 3. Check document root points to /public
# Nginx config: root /var/www/pimcore/public;

# 4. Fix file permissions
sudo chown -R www-data:www-data var/ public/var/
sudo chmod -R 775 var/ public/var/
```

### Database Connection Error (SQLSTATE[HY000] [2002])

**Fix:**
```bash
# Check if DB container is running
docker compose ps

# Check DB logs
docker compose logs db

# Verify .env DATABASE_URL
cat .env.local | grep DATABASE
# Should be: mysql://user:pass@db:3306/pimcore (not localhost!)
```

### Cache Issues (Stale Data)

```bash
# Clear everything
bin/console cache:clear
bin/console pimcore:cache:clear --all

# For Redis
redis-cli FLUSHDB

# Rebuild assets
bin/console assets:install public --symlink
```

### File Permission Errors

```bash
# Fix write permissions
sudo chown -R www-data:www-data var/ public/var/
sudo chmod -R 775 var/ public/var/ config/

# Or inside Docker
docker compose exec php bash -c "chown -R www-data:www-data var public/var"
```

### Missing Migrations

```bash
bin/console doctrine:migrations:migrate --no-interaction --env=prod
```

---

## 5. Debugging Permissions

**Symptom:** Editor cannot see objects in the tree.

**Fix:**
1. Check roles & users: Settings → Users / Roles
2. Verify the folder has **"List children"** permission enabled for the role
3. Permissions don't automatically cascade if "list children" is unchecked

```bash
# List all users
bin/console pimcore:user:list

# Check security configuration
bin/console security:check
```

---

## 6. Debug Console Commands

```bash
# List all routes
bin/console debug:router

# Show all registered services
bin/console debug:container

# List all events
bin/console debug:event-dispatcher

# List enabled bundles
bin/console debug:bundle

# Check class definitions
bin/console pimcore:deployment:classes-rebuild --dump-classes

# Test your API
curl -v http://localhost/api/products

# Check system requirements
bin/console pimcore:system:requirements:check
```

---

## 7. Debugging in Docker

```bash
# Get a shell inside PHP container
docker compose exec php bash

# Run Symfony console
docker compose exec php bin/console cache:clear

# Check PHP configuration
docker compose exec php php -i | grep -i opcache

# Check Redis connection
docker compose exec php php -r "new Redis(); $r = new Redis(); $r->connect('redis', 6379); echo $r->ping();"

# Check running processes
docker compose ps

# Full log output
docker compose logs --tail=100 php

# Restart a container
docker compose restart php
```

---

## 8. Xdebug Setup (IDE Debugging)

**Install Xdebug:**
```bash
sudo apt install php8.3-xdebug
```

**Config in `php.ini`:**
```ini
zend_extension=xdebug.so
xdebug.mode=debug
xdebug.start_with_request=yes
xdebug.client_host=host.docker.internal   ; from inside Docker
xdebug.client_port=9003
xdebug.idekey=PHPSTORM
```

**In VS Code** (PHP Debug extension):
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Listen for Xdebug",
            "type": "php",
            "request": "launch",
            "port": 9003,
            "pathMappings": {
                "/var/www/html": "${workspaceFolder}"
            }
        }
    ]
}
```

---

## Quick Debugging Checklist

When Pimcore misbehaves, follow this sequence:

```
1. Check logs:       tail -f var/log/dev.log
2. Clear cache:      bin/console cache:clear && bin/console pimcore:cache:clear
3. Rebuild classes:  bin/console pimcore:deployment:classes-rebuild -c
4. Check routes:     bin/console debug:router | grep your-route
5. Check DB:         docker compose ps && docker compose logs db
6. Check Symfony profiler for slow queries
7. Check file permissions:  ls -la var/ public/var/
8. Run requirements check:  bin/console pimcore:system:requirements:check
```
