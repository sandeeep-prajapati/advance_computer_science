# Chapter 7: Testing Strategy

## Why Tests Matter

Horizon controls production workers. A small bug can affect job processing.

A good contribution usually needs a test.

---

## Test Types

| Test Type | Covers |
|---|---|
| unit test | one class, one behavior |
| feature test | Laravel app behavior |
| console test | Artisan command behavior |
| repository test | Redis-backed state |
| regression test | prevents a known bug from returning |

---

## Regression Test Pattern

Use this structure:

```text
test name explains bug
arrange exact bad state
run the method/command
assert fixed behavior
```

Example:

```php
public function test_terminating_processes_do_not_block_autoscaling()
{
    // Arrange supervisor and process pool state.
    // Run autoscaler.
    // Assert expected process count.
}
```

---

## What Makes a Good Test

Good:

- small
- focused
- deterministic
- names the behavior
- fails before the fix
- passes after the fix

Bad:

- tests many things at once
- depends on timing too much
- passes even without the fix
- only checks implementation details

---

## Testing Redis Logic

Repository tests should cover:

- empty Redis data
- existing records
- deleted records
- malformed payloads if relevant
- trimming behavior
- expiration behavior

---

## Testing Process Logic

Process and autoscaling tests should cover:

- min process limits
- max process limits
- scale up
- scale down
- terminating processes
- multiple queues
- long-running jobs
- balancing strategy

These tests are harder but valuable.

---

## Commands to Discover

Inspect available scripts:

```bash
cat composer.json
cat phpunit.xml.dist
```

Likely commands:

```bash
vendor/bin/phpunit
composer test
```

Use the repo's actual scripts.

