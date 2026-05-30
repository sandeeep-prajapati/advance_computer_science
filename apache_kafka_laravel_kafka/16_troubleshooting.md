# Chapter 16: Troubleshooting

## `Class RdKafka\... not found`

The PHP extension is not loaded.

```bash
php -m | grep rdkafka
php --ri rdkafka
```

Remember PHP CLI and PHP-FPM can use different ini files.

---

## Cannot Connect to Broker

Common causes:

- wrong `KAFKA_BROKERS`
- Docker container using `localhost` instead of service name
- Kafka `advertised.listeners` points to a host Laravel cannot resolve
- firewall or security group blocks port
- cloud Kafka requires SASL/TLS

Inside the Laravel container, test DNS and TCP connectivity.

---

## Consumer Reads Nothing

Check:

- correct topic name
- topic has messages
- correct consumer group
- committed offset is already at end
- `auto.offset.reset` only applies when there is no committed offset
- handler is actually attached before `consume()`

Use a new temporary group ID to test from the beginning.

---

## Duplicates

Duplicates are normal after crashes and retries. Fix with:

- idempotent handlers
- unique event IDs
- database constraints
- offset commit after successful processing

---

## Lag Keeps Growing

Possible causes:

- handler too slow
- downstream database or API slow
- not enough partitions
- not enough consumers
- poison message repeatedly failing
- consumer process restarting

Start with logs, lag by partition, and handler duration metrics.
