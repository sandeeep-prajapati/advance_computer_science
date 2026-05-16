# Chapter 13: Performance & Caching

## Overview

Pimcore is powerful but can get heavy if not tuned right. Performance in Pimcore comes from stacking multiple cache layers, tuning the database, using search indexing for large datasets, and profiling to find bottlenecks.

---

## Cache Layer 1: PHP OPcache

PHP OPcache compiles PHP files to bytecode and keeps them in memory — essential for production.

**Config in `php.ini`:**
```ini
opcache.enable=1
opcache.memory_consumption=256
opcache.max_accelerated_files=20000
opcache.validate_timestamps=0   ; only set to 0 in production
opcache.revalidate_freq=0       ; production
opcache.preload=/app/var/cache/opcache_preload.php  ; optional preloading
```

---

## Cache Layer 2: Pimcore/Symfony Cache (Redis)

Pimcore uses Symfony Cache for metadata, class definitions, and object bricks. The default filesystem adapter is slow — **switch to Redis** in production.

**Install Redis extension:**
```bash
sudo apt install redis-server php-redis
```

**Config:** `config/packages/cache.yaml`

```yaml
framework:
  cache:
    app: cache.adapter.redis
    default_redis_provider: "redis://localhost:6379"
    pools:
      pimcore.cache.pool:
        adapter: cache.adapter.redis
        provider: "redis://localhost:6379"
```

**Verify Redis is running:**
```bash
redis-cli ping
# Expected: PONG
```

---

## Cache Layer 3: Full-Page Cache (Varnish / Nginx)

For CMS pages and API responses that don't change often, full-page caching prevents requests from hitting PHP at all.

### Nginx FastCGI Cache Setup

```nginx
# /etc/nginx/nginx.conf
fastcgi_cache_path /tmp/nginx_cache levels=1:2 keys_zone=PIMCORE:100m inactive=60m;
fastcgi_cache_key "$scheme$request_method$host$request_uri";

server {
    location ~ \.php$ {
        fastcgi_cache PIMCORE;
        fastcgi_cache_valid 200 60m;
        fastcgi_cache_use_stale error timeout updating http_500;
        add_header X-Cache-Status $upstream_cache_status;
    }
}
```

### Pimcore Cache Headers in Twig

```twig
{% do pimcore_cache(3600) %}  {# Cache this page for 1 hour #}
```

---

## Database Tuning (MySQL / MariaDB)

### Key Settings in `my.cnf`

```ini
[mysqld]
innodb_buffer_pool_size = 4G          ; 60-70% of available RAM
innodb_log_file_size = 256M
innodb_flush_log_at_trx_commit = 2    ; better write performance (lose max 1s on crash)
query_cache_type = 0                  ; disable in MySQL 8.x (deprecated)
max_connections = 200
```

### Add Indexes on Frequently Filtered Fields

```sql
-- Example: if you frequently filter by SKU or category
ALTER TABLE object_query_1 ADD INDEX idx_sku (sku);
ALTER TABLE object_query_1 ADD INDEX idx_category_id (category__id);
```

### Avoid Expensive Queries

```php
// BAD: loads all products with all relations eagerly
$listing = new Product\Listing();
$all = $listing->load();

// GOOD: paginate and only load what you need
$listing = new Product\Listing();
$listing->setLimit(20);
$listing->setOffset(0);
$listing->setOrderKey('name');
// Only call getRelations() when actually needed
```

---

## Enable MySQL Slow Query Log

```ini
[mysqld]
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 1   ; log queries taking > 1 second
```

---

## Elasticsearch / OpenSearch for Large Catalogs

For catalogs with more than ~10,000 products, MySQL fulltext search becomes slow. Use **Elasticsearch** for fast, faceted product search.

### Flow

```
Product saved in Pimcore
        ↓
Event subscriber fires (POST_UPDATE)
        ↓
Push document to Elasticsearch
        ↓
Frontend queries ES directly (fast)
```

### Example ES Document

```json
{
    "sku": "P001",
    "name": "Phone X",
    "price": 999,
    "categories": ["Smartphones", "Electronics"],
    "attributes": {
        "battery": "4000mAh",
        "display": "6.5 inch OLED"
    }
}
```

### Sync on Save (Event Subscriber)

```php
class ProductSearchIndexSubscriber implements EventSubscriberInterface
{
    public function __construct(private ElasticsearchClient $esClient) {}

    public static function getSubscribedEvents(): array
    {
        return [DataObjectEvents::POST_UPDATE => 'indexProduct'];
    }

    public function indexProduct(DataObjectEvent $event): void
    {
        $product = $event->getObject();
        if (!$product instanceof Product) return;

        $this->esClient->index([
            'index' => 'products',
            'id'    => $product->getId(),
            'body'  => [
                'sku'        => $product->getSku(),
                'name'       => $product->getName(),
                'price'      => $product->getPrice(),
                'category'   => $product->getCategory()?->getName(),
            ]
        ]);
    }
}
```

---

## Profiling & Identifying Bottlenecks

### 1. Symfony Profiler (Dev Mode)

The Symfony debug toolbar appears at the bottom of every page in dev mode.

Click it → see:
- Timeline (controller, Twig, events)
- All SQL queries with execution time
- Cache hits/misses
- Logs emitted during the request

**Enable if missing** in `config/packages/dev/web_profiler.yaml`:
```yaml
web_profiler:
  toolbar: true
  intercept_redirects: false
```

### 2. Blackfire.io (Production Profiling)

Blackfire generates flame graphs showing exactly where PHP time is spent.

```bash
# Install Blackfire agent
curl https://packages.blackfire.io/gpg.key | sudo apt-key add -
echo "deb http://packages.blackfire.io/debian any main" | sudo tee /etc/apt/sources.list.d/blackfire.list
sudo apt update && sudo apt install blackfire-php

# Profile a request
blackfire curl http://localhost/products
```

---

## Performance Checklist

| Item | Status |
|---|---|
| Enable OPcache in production | ✅ |
| Use Redis for Symfony/Pimcore cache | ✅ |
| Configure Nginx/Varnish full-page cache | ✅ |
| Tune MySQL buffer pool (60-70% RAM) | ✅ |
| Add DB indexes on filtered fields | ✅ |
| Enable MySQL slow query log | ✅ |
| Use Elasticsearch for large catalogs (>10k products) | ✅ |
| Profile with Symfony Profiler + Blackfire | ✅ |
| Paginate all listings (never load ALL records) | ✅ |
| Use lazy loading for relations | ✅ |
| Serve images via thumbnails + CDN | ✅ |

---

## Quick Command Reference

```bash
# Clear all caches
bin/console cache:clear
bin/console pimcore:cache:clear --all

# Warm up caches
bin/console cache:warmup --env=prod

# Check system requirements
bin/console pimcore:system:requirements:check

# Test database query
bin/console doctrine:query:sql "SELECT COUNT(*) FROM objects WHERE o_type='object'"

# Check Elasticsearch connection
curl http://localhost:9200/_cluster/health
```
