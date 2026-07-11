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

