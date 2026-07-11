# RabbitMQ Wrapper: Root and Go Book

> Shared RabbitMQ architecture plus Go/Gin wrapper notes. PHP and Python have separate books in their folders.

---

# RabbitMQ Wrapper Developer Book

A detailed guide to designing RabbitMQ wrappers for real applications.

This folder is organized by language/framework:

```text
RabbitMQ_Rapper/
├── 00_index.md                         shared concepts and reading path
├── 01_rabbitmq_core_concepts.md         broker, exchange, queue, ack, prefetch
├── 02_wrapper_architecture.md           clean wrapper design used by every stack
├── 03_go_rabbitmq_wrapper.md            Go worker wrapper
├── 04_go_gin_integration.md             Gin API producer + background worker
├── 05_operations_retries_dlq.md         retry, DLQ, idempotency, monitoring
├── book.md                              combined root book
├── PHP/                                 Core PHP, Laravel, Symfony, Pimcore
└── Python/                              Python and Django
```

---

## What You Will Learn

- How RabbitMQ routes messages through exchanges, queues, bindings, and routing keys.
- How to write a clean wrapper around RabbitMQ instead of scattering raw AMQP code.
- How to design producers, consumers, handlers, serializers, retry policies, and DLQs.
- How to write Go RabbitMQ workers and connect them with Gin APIs.
- How to implement the same architecture in PHP frameworks.
- How to implement the same architecture in Python/Django.
- How to write production-safe message code with logging, metrics, idempotency, and graceful shutdown.

---

## Reading Path

**Start here:**

1. [RabbitMQ Core Concepts](01_rabbitmq_core_concepts.md)
2. [Wrapper Architecture](02_wrapper_architecture.md)
3. [Operations, Retries & DLQ](05_operations_retries_dlq.md)

**For Go:**

1. [Go RabbitMQ Wrapper](03_go_rabbitmq_wrapper.md)
2. [Go Gin Integration](04_go_gin_integration.md)

**For PHP:**

1. [PHP RabbitMQ Notes](PHP/00_index.md)

**For Python:**

1. [Python RabbitMQ Notes](Python/00_index.md)

---

## Wrapper Mental Model

```text
Producer Application
└── Publisher Wrapper
    ├── validates message
    ├── creates envelope
    ├── serializes JSON
    ├── sets headers
    └── publishes to exchange with routing key

RabbitMQ Broker
├── exchange
├── binding
├── queue
├── retry queue
└── dead-letter queue

Consumer Worker
└── Consumer Wrapper
    ├── receives raw delivery
    ├── decodes envelope
    ├── finds handler
    ├── runs business service
    ├── ack on success
    ├── retry on temporary failure
    └── reject/DLQ on permanent failure
```

---

## Main Design Rule

Keep AMQP details at the boundary.

```text
Good:
RabbitMQ delivery -> wrapper -> message envelope -> handler -> application service

Bad:
business service directly creates channels, acks messages, parses headers, and handles DLQ rules
```

Your handler should process a clean message. The wrapper should handle RabbitMQ mechanics.


---

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


---

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


---

# Chapter 3: Go RabbitMQ Wrapper

## Project Structure

```text
go-rabbit-app/
├── go.mod
├── cmd/
│   └── worker/
│       └── main.go
└── internal/
    ├── messaging/
    │   ├── envelope.go
    │   ├── publisher.go
    │   ├── consumer.go
    │   ├── registry.go
    │   └── topology.go
    └── handlers/
        └── product_updated.go
```

Install:

```bash
go get github.com/rabbitmq/amqp091-go
go get github.com/google/uuid
```

---

## Envelope

```go
package messaging

import "time"

type Envelope struct {
    ID         string                 `json:"id"`
    Type       string                 `json:"type"`
    Version    int                    `json:"version"`
    OccurredAt time.Time              `json:"occurred_at"`
    Payload    map[string]interface{} `json:"payload"`
    Meta       map[string]interface{} `json:"meta"`
}
```

---

## Handler Interface

```go
package messaging

import "context"

type Handler interface {
    Handle(ctx context.Context, message Envelope) error
}
```

Handlers return errors. The consumer decides ack/retry/DLQ.

---

## Registry

```go
package messaging

import "fmt"

type Registry struct {
    handlers map[string]Handler
}

func NewRegistry() *Registry {
    return &Registry{handlers: map[string]Handler{}}
}

func (r *Registry) Register(messageType string, handler Handler) {
    r.handlers[messageType] = handler
}

func (r *Registry) Get(messageType string) (Handler, error) {
    handler, ok := r.handlers[messageType]
    if !ok {
        return nil, fmt.Errorf("no handler registered for %s", messageType)
    }

    return handler, nil
}
```

---

## Topology Setup

```go
package messaging

import amqp "github.com/rabbitmq/amqp091-go"

func DeclareTopology(ch *amqp.Channel) error {
    if err := ch.ExchangeDeclare(
        "app.events",
        "topic",
        true,
        false,
        false,
        false,
        nil,
    ); err != nil {
        return err
    }

    _, err := ch.QueueDeclare(
        "product.search_index.queue",
        true,
        false,
        false,
        false,
        amqp.Table{
            "x-dead-letter-exchange": "app.dlx",
        },
    )
    if err != nil {
        return err
    }

    return ch.QueueBind(
        "product.search_index.queue",
        "product.*",
        "app.events",
        false,
        nil,
    )
}
```

---

## Publisher

```go
package messaging

import (
    "context"
    "encoding/json"

    amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
    channel  *amqp.Channel
    exchange string
}

func NewPublisher(channel *amqp.Channel, exchange string) *Publisher {
    return &Publisher{channel: channel, exchange: exchange}
}

func (p *Publisher) Publish(ctx context.Context, routingKey string, message Envelope) error {
    body, err := json.Marshal(message)
    if err != nil {
        return err
    }

    return p.channel.PublishWithContext(
        ctx,
        p.exchange,
        routingKey,
        false,
        false,
        amqp.Publishing{
            ContentType:  "application/json",
            DeliveryMode: amqp.Persistent,
            MessageId:    message.ID,
            Body:         body,
        },
    )
}
```

---

## Consumer

```go
package messaging

import (
    "context"
    "encoding/json"
    "log"

    amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
    channel  *amqp.Channel
    registry *Registry
}

func NewConsumer(channel *amqp.Channel, registry *Registry) *Consumer {
    return &Consumer{channel: channel, registry: registry}
}

func (c *Consumer) Consume(ctx context.Context, queue string) error {
    if err := c.channel.Qos(1, 0, false); err != nil {
        return err
    }

    deliveries, err := c.channel.Consume(queue, "", false, false, false, false, nil)
    if err != nil {
        return err
    }

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case delivery, ok := <-deliveries:
            if !ok {
                return nil
            }
            c.process(ctx, delivery)
        }
    }
}

func (c *Consumer) process(ctx context.Context, delivery amqp.Delivery) {
    var message Envelope

    if err := json.Unmarshal(delivery.Body, &message); err != nil {
        log.Printf("invalid message json: %v", err)
        _ = delivery.Nack(false, false)
        return
    }

    handler, err := c.registry.Get(message.Type)
    if err != nil {
        log.Printf("handler missing: %v", err)
        _ = delivery.Nack(false, false)
        return
    }

    if err := handler.Handle(ctx, message); err != nil {
        log.Printf("handler failed message_id=%s type=%s error=%v", message.ID, message.Type, err)
        _ = delivery.Nack(false, false)
        return
    }

    _ = delivery.Ack(false)
}
```

---

## Example Handler

```go
package handlers

import (
    "context"
    "fmt"

    "go-rabbit-app/internal/messaging"
)

type ProductUpdatedHandler struct{}

func (h *ProductUpdatedHandler) Handle(ctx context.Context, message messaging.Envelope) error {
    productID, ok := message.Payload["product_id"]
    if !ok {
        return fmt.Errorf("product_id is required")
    }

    fmt.Printf("index product %v\n", productID)
    return nil
}
```

---

## Worker Main

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "go-rabbit-app/internal/handlers"
    "go-rabbit-app/internal/messaging"

    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatal(err)
    }
    defer ch.Close()

    if err := messaging.DeclareTopology(ch); err != nil {
        log.Fatal(err)
    }

    registry := messaging.NewRegistry()
    registry.Register("product.updated", &handlers.ProductUpdatedHandler{})

    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    consumer := messaging.NewConsumer(ch, registry)
    if err := consumer.Consume(ctx, "product.search_index.queue"); err != nil {
        log.Println(err)
    }
}
```

---

## Go Best Practices

- Use `context.Context` everywhere.
- Return errors from handlers.
- Keep Gin HTTP server and worker process separate.
- Use manual ack.
- Start with prefetch `1`.
- Add reconnect logic before production.
- Add structured logs and metrics.


---

# Chapter 4: Go Gin Integration

## Architecture

Gin should publish messages from HTTP handlers. A separate worker should consume messages.

```text
Gin API process
└── POST /products
    ├── save product
    └── publish product.created

Worker process
└── consume product.search_index.queue
    └── update search index
```

Do not run the RabbitMQ consumer loop inside the Gin request process unless you intentionally build a combined service.

---

## Gin Project Structure

```text
go-rabbit-app/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
└── internal/
    ├── http/
    │   └── product_controller.go
    ├── messaging/
    └── services/
```

---

## Product Controller

```go
package http

import (
    "net/http"
    "time"

    "go-rabbit-app/internal/messaging"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type ProductController struct {
    publisher *messaging.Publisher
}

func NewProductController(publisher *messaging.Publisher) *ProductController {
    return &ProductController{publisher: publisher}
}

func (pc *ProductController) Create(c *gin.Context) {
    productID := uuid.NewString()

    message := messaging.Envelope{
        ID:         uuid.NewString(),
        Type:       "product.created",
        Version:    1,
        OccurredAt: time.Now().UTC(),
        Payload: map[string]interface{}{
            "product_id": productID,
        },
        Meta: map[string]interface{}{
            "source":         "gin-api",
            "correlation_id": c.GetHeader("X-Correlation-ID"),
        },
    }

    if err := pc.publisher.Publish(c.Request.Context(), "product.created", message); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "message publish failed"})
        return
    }

    c.JSON(http.StatusAccepted, gin.H{
        "product_id": productID,
        "status":     "queued",
    })
}
```

---

## API Main

```go
package main

import (
    "log"

    apphttp "go-rabbit-app/internal/http"
    "go-rabbit-app/internal/messaging"

    "github.com/gin-gonic/gin"
    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatal(err)
    }
    defer ch.Close()

    if err := messaging.DeclareTopology(ch); err != nil {
        log.Fatal(err)
    }

    publisher := messaging.NewPublisher(ch, "app.events")
    controller := apphttp.NewProductController(publisher)

    router := gin.Default()
    router.POST("/products", controller.Create)

    if err := router.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
```

---

## Correlation ID Middleware

```go
func CorrelationID() gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.GetHeader("X-Correlation-ID")
        if id == "" {
            id = uuid.NewString()
        }

        c.Header("X-Correlation-ID", id)
        c.Set("correlation_id", id)
        c.Next()
    }
}
```

Use the same correlation ID in message metadata.

---

## HTTP Response Choice

When publishing async work, return `202 Accepted`.

```json
{
  "status": "queued",
  "product_id": "abc"
}
```

Do not return `200 OK` claiming the slow background work is complete.

---

## Gin Checklist

- API publishes messages.
- Worker consumes messages.
- Use correlation IDs.
- Return `202 Accepted` for async work.
- Do not block HTTP request on slow worker job.
- Use separate Docker/systemd processes for API and worker.


---

# Chapter 5: Operations, Retries & DLQ

## Production Failure Model

RabbitMQ workers fail for normal reasons:

- database timeout
- external API slow
- invalid JSON
- missing payload field
- duplicate message
- worker deployment restart
- RabbitMQ reconnect
- handler bug

Your wrapper exists to make these failures controlled.

---

## Retry Queue Pattern

Do not retry immediately forever.

Use delayed retry queues with TTL.

```text
main queue
└── worker failure
    └── publish to retry queue
        └── retry queue waits using TTL
            └── dead-letters back to main exchange
```

Example queues:

```text
product.search_index.queue
product.search_index.retry.10s
product.search_index.retry.1m
product.search_index.retry.5m
product.search_index.dlq
```

---

## Dead Letter Queue

DLQ stores messages that cannot be processed.

Send to DLQ when:

- JSON cannot be decoded
- required payload is missing
- handler is not registered
- retry limit exceeded
- permanent business rule failure

DLQ should be monitored. A DLQ with messages is a signal to investigate.

---

## Idempotency

RabbitMQ can deliver a message more than once. Your handler must be safe.

```text
message id = msg-123
processed_messages table already has msg-123
handler skips work
consumer ack
```

Use idempotency for:

- payment capture
- email sending
- stock update
- ERP sync
- search indexing

---

## Structured Logging

Log fields, not only strings.

```json
{
  "level": "info",
  "message": "rabbitmq_message_processed",
  "message_id": "msg-123",
  "type": "product.updated",
  "queue": "product.search_index.queue",
  "correlation_id": "req-999",
  "duration_ms": 51
}
```

Every log should include:

- message id
- message type
- queue
- correlation id
- status

---

## Metrics

Track:

- consumed count
- success count
- failure count
- retry count
- DLQ count
- processing duration
- queue depth
- consumer count

Alerts:

- DLQ has messages
- queue depth increasing
- no consumers
- processing duration too high
- retry rate too high

---

## Graceful Shutdown

Workers should handle SIGTERM.

```text
receive shutdown signal
stop accepting new work
finish current message
ack/nack current message
close channel
close connection
exit
```

This matters for Docker, Kubernetes, Supervisor, and systemd.

---

## Deployment

Run workers as separate processes.

```text
web process      handles HTTP
worker process   consumes RabbitMQ
scheduler        publishes periodic jobs
```

Never rely on a browser request to keep a worker running.

