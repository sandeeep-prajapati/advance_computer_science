# Chapter 2: Setup & Installation

## Requirements

Typical Laravel Kafka projects need:

- PHP and Laravel supported by your selected package version
- `ext-rdkafka`
- `librdkafka`
- Composer package `mateusjunges/laravel-kafka`
- A Kafka broker reachable from the Laravel runtime

Always check the package release notes for exact PHP and Laravel version support.

---

## Install System Dependencies

Ubuntu example:

```bash
sudo apt-get update
sudo apt-get install -y librdkafka-dev
sudo pecl install rdkafka
```

Enable the extension:

```ini
extension=rdkafka
```

Verify:

```bash
php -m | grep rdkafka
php --ri rdkafka
```

---

## Install the Laravel Package

```bash
composer require mateusjunges/laravel-kafka
```

Publish config if the package version provides a config file:

```bash
php artisan vendor:publish --provider="Junges\Kafka\KafkaServiceProvider"
```

---

## Local Kafka with Docker

Use a single-broker Kafka container for learning:

```yaml
services:
  kafka:
    image: bitnami/kafka:latest
    ports:
      - "9092:9092"
    environment:
      KAFKA_CFG_NODE_ID: 1
      KAFKA_CFG_PROCESS_ROLES: broker,controller
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: 1@kafka:9093
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE: "true"
```

For Laravel running inside another container, advertise an internal listener such as `kafka:9092`.

---

## Environment Variables

```dotenv
KAFKA_BROKERS=localhost:9092
KAFKA_CONSUMER_GROUP_ID=laravel-app
KAFKA_SECURITY_PROTOCOL=PLAINTEXT
```

For cloud Kafka, you will normally add SASL and TLS values.

---

## Smoke Test

Create a route or Artisan command that publishes one event:

```php
use Junges\Kafka\Facades\Kafka;
use Junges\Kafka\Message\Message;

Kafka::publish(config('kafka.brokers'))
    ->onTopic('dev.ping')
    ->withMessage(new Message(
        headers: ['source' => 'laravel'],
        body: ['message' => 'hello kafka']
    ))
    ->send();
```

If this fails, debug in this order:

1. Is `rdkafka` loaded by PHP CLI and FPM?
2. Can the Laravel container reach the broker host and port?
3. Is Kafka advertising a hostname Laravel can resolve?
4. Does your cloud provider require SASL/TLS?
