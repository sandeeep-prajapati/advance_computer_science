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

