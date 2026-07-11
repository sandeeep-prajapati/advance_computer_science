# Chapter 1: RabbitMQ Core Concepts

## RabbitMQ in One Picture

```text
Producer
└── publishes message to Exchange
    └── Exchange checks bindings
        └── routes message to Queue
            └── Consumer receives message
                ├── ack success
                └── nack/reject failure
```

RabbitMQ is not just a queue. It is a broker with routing rules.

---

## Producer

A producer creates a message.

Examples:

- Product service publishes `product.created`.
- Order service publishes `order.placed`.
- Pimcore publishes `asset.updated`.
- Laravel API publishes `email.send`.
- Gin API publishes `payment.captured`.
- Django admin publishes `report.generate`.

Producer code should be short:

```text
create message envelope
publish to exchange with routing key
return response to user
```

Do not run slow work in the HTTP request when it can be processed asynchronously.

---

## Message

A message has two parts:

```text
body       JSON payload
properties metadata controlled by AMQP client
```

Body example:

```json
{
  "id": "msg-1001",
  "type": "product.created",
  "version": 1,
  "occurred_at": "2026-07-11T10:15:00Z",
  "payload": {
    "product_id": 501
  },
  "meta": {
    "correlation_id": "req-abc",
    "source": "admin"
  }
}
```

Properties can include:

- `content_type`
- `delivery_mode`
- `correlation_id`
- `message_id`
- `timestamp`
- custom headers

---

## Exchange

Producer publishes to an exchange.

| Exchange Type | Routing Style | Example |
|---|---|---|
| direct | exact routing key | `email.send` |
| topic | wildcard routing key | `product.*`, `order.#` |
| fanout | broadcast to all queues | cache clear broadcast |
| headers | route by headers | advanced filtering |

Most app wrappers should start with a **topic exchange**.

---

## Queue

A queue stores messages until consumers process them.

Good queue names describe the job:

```text
product.search_index.queue
email.order_confirmation.queue
erp.product_sync.queue
report.pdf_generation.queue
```

Weak queue names:

```text
events
queue1
products
```

The queue belongs to the consumer purpose.

---

## Binding and Routing Key

Binding connects an exchange to a queue.

```text
exchange: app.events
queue: product.search_index.queue
routing key pattern: product.*
```

Messages:

```text
product.created  -> delivered
product.updated  -> delivered
order.placed     -> not delivered
```

Topic wildcard rules:

| Pattern | Meaning |
|---|---|
| `*` | one word |
| `#` | zero or more words |
| `product.*` | product.created, product.updated |
| `order.#` | order.placed, order.payment.captured |

---

## Ack, Nack, Reject

Manual ack is critical.

| Action | Meaning |
|---|---|
| `ack` | processed successfully, remove from queue |
| `nack requeue=true` | failed, put back into queue |
| `nack requeue=false` | failed, reject or dead-letter |
| `reject` | reject one message |

Rule:

```text
ack only after successful processing
```

If you ack before processing and your worker crashes, the message is lost.

---

## Prefetch

Prefetch controls how many unacked messages a worker can hold.

```text
prefetch = 1
worker gets one message at a time
good for slow jobs

prefetch = 20
worker can receive 20 messages before acking
good for fast jobs
```

Start with `prefetch = 1` while learning.

---

## Durability

For production:

- exchange durable: true
- queue durable: true
- message persistent: true
- manual ack: true

Durable queues survive broker restart. Persistent messages are saved to disk by RabbitMQ.

---

## Competing Consumers

Multiple workers can consume the same queue.

```text
product.search_index.queue
├── worker 1
├── worker 2
└── worker 3
```

RabbitMQ distributes messages among them. This is how you scale workers horizontally.

