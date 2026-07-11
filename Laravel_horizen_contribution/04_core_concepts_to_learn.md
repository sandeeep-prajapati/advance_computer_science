# Chapter 4: Core Concepts to Learn

## Laravel Queues

You must understand:

- jobs
- dispatching
- queue connections
- queue names
- `queue:work`
- retries
- timeouts
- failed jobs
- middleware
- tags

Practice:

```php
dispatch(new ImportProductsJob($fileId))->onQueue('imports');
```

Then run:

```bash
php artisan queue:work redis --queue=imports
```

---

## Redis

Horizon is Redis-based.

Learn:

- lists
- sorted sets
- hashes
- key expiration
- Lua scripts basics
- Redis connection names
- Redis cluster caveats

You do not need to become a Redis expert first, but you must know enough to understand Horizon's repositories.

---

## Processes

Horizon manages OS processes.

Learn:

- process start/stop
- signals
- graceful shutdown
- long-running PHP behavior
- memory limits
- process exit codes

Important idea:

```text
HTTP request lifecycle is short.
Horizon worker lifecycle is long.
```

Long-running code has different bugs.

---

## Supervisor Concepts

Horizon supervisor config controls:

- connection
- queues
- balance strategy
- min/max processes
- memory
- timeout
- tries
- nice value
- max time
- max jobs

You should know how config maps to runtime process behavior.

---

## Balancing and Autoscaling

This is one of Horizon's hardest areas.

Concepts:

- current process count
- desired process count
- queue size
- runtime average
- scale up
- scale down
- terminating processes
- min/max process limits
- balancing between multiple queues

Good contribution area after you understand tests.

---

## Metrics

Horizon dashboard depends on metrics.

Learn:

- throughput
- runtime
- job counts
- failed counts
- snapshots
- trimming old metrics

Contribution areas:

- inaccurate numbers
- missing edge cases
- dashboard display mismatch
- cleanup behavior

---

## Tags

Horizon can show jobs by tags.

Job tags help track:

- user IDs
- order IDs
- model IDs
- tenant IDs

Contribution areas:

- tag display
- tag search
- tag cleanup
- handling special characters

