# Chapter 10: Security and Cloud Kafka

## Local vs Production

Local Kafka often uses plaintext:

```dotenv
KAFKA_SECURITY_PROTOCOL=PLAINTEXT
```

Managed Kafka usually uses SASL and TLS:

```dotenv
KAFKA_SECURITY_PROTOCOL=SASL_SSL
KAFKA_SASL_MECHANISM=PLAIN
KAFKA_SASL_USERNAME=...
KAFKA_SASL_PASSWORD=...
```

---

## SASL Options

```php
->withConfigOptions([
    'security.protocol' => env('KAFKA_SECURITY_PROTOCOL', 'SASL_SSL'),
    'sasl.mechanisms' => env('KAFKA_SASL_MECHANISM', 'PLAIN'),
    'sasl.username' => env('KAFKA_SASL_USERNAME'),
    'sasl.password' => env('KAFKA_SASL_PASSWORD'),
])
```

Use the same security options for producers and consumers.

---

## Secret Handling

- Store credentials in environment variables or a secret manager.
- Do not commit `.env` files.
- Rotate cloud Kafka API keys.
- Give each service its own Kafka principal.
- Restrict ACLs by topic and consumer group.

---

## ACL Shape

For a producer service:

```text
WRITE on topic orders.created
DESCRIBE on topic orders.created
```

For a consumer service:

```text
READ on topic orders.created
DESCRIBE on topic orders.created
READ on group payment-service
```

Least privilege makes mistakes much easier to contain.
