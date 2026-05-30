# Laravel Kafka — Complete Learning Guide

> A practical, hands-on guide to using Apache Kafka from Laravel with the `mateusjunges/laravel-kafka` package.

---

## Table of Contents

| Chapter | Title | Topics |
|---|---|---|
| 01 | [Introduction](01_introduction.md) | Laravel + Kafka use cases, event-driven apps, when Kafka fits |
| 02 | [Setup & Installation](02_setup_and_installation.md) | PHP extension, Composer package, Docker Kafka, env config |
| 03 | [Configuration](03_configuration.md) | Brokers, topics, consumer groups, config publishing, options |
| 04 | [Producing Messages](04_producing_messages.md) | Facade API, keys, headers, payloads, reliability options |
| 05 | [Consuming Messages](05_consuming_messages.md) | Consumers, handlers, groups, offsets, long-running workers |
| 06 | [Message Design](06_message_design.md) | Event names, payload contracts, correlation IDs, versioning |
| 07 | [Laravel Events, Jobs, and Outbox](07_events_jobs_outbox.md) | Domain events, queued publishing, transactional outbox pattern |
| 08 | [Retries and Dead Letter Topics](08_retries_and_dlq.md) | Backoff, poison messages, DLQ naming, replay strategy |
| 09 | [Testing](09_testing.md) | Kafka fake, assertions, handler tests, integration tests |
| 10 | [Security and Cloud Kafka](10_security_and_cloud.md) | SASL, TLS, Confluent Cloud, MSK, secrets |
| 11 | [Serialization and Schema Registry](11_serialization_schema_registry.md) | JSON, custom serializers, Avro/Protobuf concepts |
| 12 | [Observability](12_observability.md) | Logs, metrics, lag, traces, health checks |
| 13 | [Performance and Scaling](13_performance_scaling.md) | Batching, async publishing, partitions, Horizon/Supervisor |
| 14 | [Deployment](14_deployment.md) | Containers, graceful shutdown, process management |
| 15 | [Production Patterns](15_production_patterns.md) | Idempotency, ordering, sagas, CDC, multi-service contracts |
| 16 | [Troubleshooting](16_troubleshooting.md) | Common errors, extension issues, connectivity, offsets |
| 17 | [Capstone Project](17_capstone_project.md) | Order event pipeline in Laravel |
| 18 | [Learning Roadmap](18_learning_roadmap.md) | 4-week plan and practice checklist |

---

## Architecture

```text
Laravel App
  HTTP Controller / Command / Domain Service
        |
        | publish event
        v
  Junges\Kafka Producer  --->  Kafka Topic: orders.created
                                      |
                                      | consumed by group payment-service
                                      v
                              Laravel Kafka Consumer
                                      |
                                      v
                              Handler / Job / Application Service
```

---

## Recommended Reading Path

- Start with Chapters 1-5 if you are new to Kafka in Laravel.
- Use Chapters 7-9 before wiring Kafka into real business workflows.
- Use Chapters 10-16 before production.
- Build Chapter 17 after reading the guide once.

---

## References

- Laravel Kafka package: https://github.com/mateusjunges/laravel-kafka
- Laravel Kafka docs: https://laravelkafka.com
- Apache Kafka docs: https://kafka.apache.org/documentation/
