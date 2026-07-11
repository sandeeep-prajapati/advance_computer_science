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

