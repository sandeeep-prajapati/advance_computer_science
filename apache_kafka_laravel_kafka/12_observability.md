# Chapter 12: Observability

## What to Log

Log structured context:

```php
Log::info('Kafka message consumed', [
    'topic' => $message->getTopicName(),
    'partition' => $message->getPartition(),
    'offset' => $message->getOffset(),
    'key' => $message->getKey(),
    'event_id' => data_get($message->getBody(), 'event_id'),
    'event_type' => data_get($message->getBody(), 'event_type'),
]);
```

Avoid logging full payloads if they contain personal data.

---

## Metrics

Track:

- messages produced per minute
- produce failures
- messages consumed per minute
- handler duration
- handler failures
- consumer lag
- DLQ count
- retry count

Consumer lag is the most important operational signal. Growing lag means consumers are slower than producers or are failing.

---

## Tracing

Propagate correlation IDs:

```php
headers: [
    'correlation_id' => request()->header('X-Correlation-ID', (string) str()->uuid()),
]
```

Each consumer should include the same ID in logs and outbound calls.

---

## Health Checks

A Laravel HTTP health check should not only say the web app is alive. Add separate checks for:

- broker connectivity
- producer publish ability to a health topic
- consumer process running
- consumer lag below threshold

Process health and stream health are different things.
