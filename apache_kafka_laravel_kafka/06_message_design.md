# Chapter 6: Message Design

## Event Envelope

Use a consistent envelope:

```json
{
  "event_id": "01J...",
  "event_type": "orders.created",
  "event_version": 1,
  "occurred_at": "2026-05-30T10:00:00+05:30",
  "producer": "checkout-api",
  "data": {
    "order_id": 123,
    "user_id": 50,
    "total": 1999,
    "currency": "INR"
  }
}
```

This makes consumers easier to build and debug.

---

## Required Metadata

| Field | Why |
|---|---|
| `event_id` | Idempotency and tracing |
| `event_type` | Dispatching and logging |
| `event_version` | Schema evolution |
| `occurred_at` | Business time |
| `producer` | Ownership |
| `correlation_id` | Request tracing across services |

---

## Compatibility Rules

Safe changes:

- Add optional fields.
- Add new event types.
- Add nullable fields with defaults.

Dangerous changes:

- Rename fields.
- Change numeric strings to numbers without coordination.
- Remove fields that consumers may still read.
- Reuse an event name for a different meaning.

---

## Versioning

Prefer versioning in the payload or header:

```php
headers: [
    'event' => 'orders.created',
    'event_version' => '2',
]
```

Avoid creating `orders.created.v1`, `orders.created.v2`, and `orders.created.v3` topics unless the old and new streams truly need different retention, permissions, or consumers.
