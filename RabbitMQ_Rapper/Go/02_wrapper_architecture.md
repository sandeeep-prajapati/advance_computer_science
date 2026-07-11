# Chapter 2: Wrapper Architecture

## Why Build a Wrapper?

Raw RabbitMQ client code is low level:

- channels
- exchanges
- queue declarations
- delivery tags
- ack/nack
- headers
- reconnects
- JSON decoding
- retries

If every handler manages these details, your code becomes hard to test and unsafe.

A wrapper gives one clean API:

```text
publisher.publish("product.created", payload)
consumer.consume("product.search_index.queue")
```

---

## Clean Architecture

```text
Infrastructure Layer
└── RabbitMQ client, channel, exchange, queue, ack, retry

Messaging Layer
└── Publisher, Consumer, Envelope, Serializer, Handler Registry

Application Layer
└── ProductIndexer, SendEmail, SyncToERP

Domain Layer
└── Product, Order, Customer, business rules
```

RabbitMQ belongs to infrastructure. Your application services should not depend directly on it.

---

## Message Envelope

Use the same envelope shape in every language.

```json
{
  "id": "msg-01",
  "type": "product.updated",
  "version": 1,
  "occurred_at": "2026-07-11T10:30:00Z",
  "payload": {
    "product_id": 501
  },
  "meta": {
    "correlation_id": "req-777",
    "source": "pimcore",
    "retry_count": 0
  }
}
```

Fields:

| Field | Purpose |
|---|---|
| `id` | idempotency and tracing |
| `type` | maps to handler |
| `version` | schema evolution |
| `occurred_at` | event time |
| `payload` | business data |
| `meta` | correlation, source, retries |

---

## Core Components

```text
MessageEnvelope
MessageSerializer
MessagePublisher
MessageConsumer
MessageHandler
HandlerRegistry
RetryPolicy
ErrorClassifier
```

### Publisher

Responsible for:

- creating AMQP message
- JSON serialization
- setting headers
- selecting exchange
- selecting routing key

### Consumer

Responsible for:

- reading from queue
- decoding message
- finding handler
- calling handler
- ack/nack
- retry/DLQ

### Handler

Responsible for business process only:

```text
ProductUpdatedHandler -> ProductIndexer.index(product_id)
```

---

## Handler Registry

The registry maps message type to handler.

```text
product.created -> ProductCreatedHandler
product.updated -> ProductUpdatedHandler
order.placed    -> SendOrderConfirmationHandler
```

If a type has no handler:

- log error
- reject message
- send to DLQ

Do not silently ack unknown messages.

---

## Error Classification

Not all failures are equal.

| Failure | Example | Action |
|---|---|---|
| temporary | database timeout | retry |
| permanent | missing required `product_id` | DLQ |
| unknown | unexpected exception | limited retry then DLQ |

Your wrapper should classify errors so poison messages do not retry forever.

---

## Consumer Algorithm

```text
receive delivery
try decode envelope
try validate envelope
try find handler
try run handler
if success:
    ack
if temporary failure:
    send to retry queue
    ack original
if permanent failure:
    send to DLQ
    ack original
if wrapper cannot republish:
    nack original safely
```

This is safer than blindly `nack requeue=true`, which can create infinite fast loops.

---

## Folder Design

Use this idea in every stack:

```text
Messaging/
├── Envelope
├── Serializer
├── Publisher
├── Consumer
├── Handler
├── HandlerRegistry
├── RetryPolicy
└── RabbitMQ
    ├── ConnectionFactory
    ├── TopologySetup
    └── RabbitPublisher
```

Framework-specific code should live outside this core:

```text
Laravel Artisan Command
Symfony Console Command
Django Management Command
Gin HTTP Handler
```

