# Chapter 11: Serialization and Schema Registry

## JSON First

JSON is a good first choice for Laravel teams:

```php
body: [
    'event_id' => (string) str()->uuid(),
    'event_type' => 'orders.created',
    'event_version' => 1,
    'data' => [
        'order_id' => $order->id,
    ],
]
```

Benefits:

- Easy to inspect
- Native PHP arrays map naturally
- No schema infrastructure required

Costs:

- Runtime validation is your responsibility
- Type drift can happen silently
- Payloads are larger than binary formats

---

## Custom Serializers

Use custom serializers when:

- every message must share a standard envelope
- dates and money need consistent formatting
- you need encryption or compression at the application layer
- you are integrating with Avro, Protobuf, or JSON Schema

---

## Schema Registry

Schema Registry stores schemas and enforces compatibility. It is common with Avro, Protobuf, and JSON Schema.

Compatibility modes:

| Mode | Meaning |
|---|---|
| Backward | New consumers can read old data |
| Forward | Old consumers can read new data |
| Full | Both backward and forward |
| None | No compatibility checks |

For multi-service teams, schema compatibility prevents accidental breaking changes.

---

## Laravel Validation at Boundaries

Validate consumed payloads before business logic:

```php
Validator::make($payload, [
    'event_id' => ['required', 'uuid'],
    'event_type' => ['required', 'string'],
    'event_version' => ['required', 'integer'],
    'data.order_id' => ['required', 'integer'],
])->validate();
```

Invalid messages should go to a DLQ with enough context to debug.
