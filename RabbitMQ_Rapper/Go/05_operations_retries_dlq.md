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

