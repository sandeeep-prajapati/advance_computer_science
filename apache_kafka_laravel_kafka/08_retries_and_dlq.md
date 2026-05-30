# Chapter 8: Retries and Dead Letter Topics

## Failure Types

| Failure | Example | Retry? |
|---|---|---|
| Transient | HTTP timeout, DB deadlock | Yes |
| Recoverable later | Missing dependent record | Yes, delayed |
| Permanent | Invalid payload shape | No, send to DLQ |
| Bug | Handler throws every time | DLQ after limited attempts |

---

## Idempotent Consumer First

Retries create duplicates. Make handlers safe:

```php
ProcessedKafkaMessage::firstOrCreate([
    'event_id' => $payload['event_id'],
]);
```

Or enforce natural uniqueness:

```php
Payment::updateOrCreate(
    ['order_id' => $payload['data']['order_id']],
    ['status' => 'authorized']
);
```

---

## Dead Letter Topic Naming

```text
orders.created.dlq
payments.captured.dlq
```

DLQ payload should include:

- original topic
- original partition and offset
- original key
- original headers
- original body
- exception class and message
- failed at timestamp

---

## Replay Strategy

Do not replay DLQ blindly. Use a command:

```bash
php artisan kafka:replay-dlq orders.created.dlq --limit=100
```

Before replay:

1. Fix the bug or data issue.
2. Pick a small batch.
3. Monitor consumer lag and error logs.
4. Keep the original DLQ event for audit if required.
