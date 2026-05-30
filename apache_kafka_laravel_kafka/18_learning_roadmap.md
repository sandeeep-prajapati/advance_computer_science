# Chapter 18: Learning Roadmap

## Week 1: Kafka and Laravel Basics

- Understand topics, partitions, offsets, and consumer groups.
- Run Kafka locally with Docker.
- Install `ext-rdkafka` and the Laravel package.
- Publish a message from an Artisan command.
- Consume a message with a simple handler.

---

## Week 2: Application Design

- Create a domain event.
- Publish after database commit.
- Add a standard event envelope.
- Add correlation IDs.
- Test producers with Kafka fake.

---

## Week 3: Reliability

- Implement idempotent consumers.
- Add retry strategy.
- Add DLQ topic.
- Learn offset behavior.
- Practice replaying a small DLQ batch.

---

## Week 4: Production Readiness

- Configure SASL/TLS for managed Kafka.
- Run consumers under Supervisor, Docker, or Kubernetes.
- Add structured logs and metrics.
- Monitor consumer lag.
- Document topic ownership and schemas.

---

## Practice Checklist

- [ ] Produce `orders.created`
- [ ] Consume with a named group
- [ ] Use stable Kafka keys
- [ ] Add event IDs
- [ ] Make handler idempotent
- [ ] Test with Kafka fake
- [ ] Run a consumer as a managed process
- [ ] Add DLQ handling
- [ ] Configure SASL/TLS
- [ ] Monitor lag
