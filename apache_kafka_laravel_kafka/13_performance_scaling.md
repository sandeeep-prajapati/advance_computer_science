# Chapter 13: Performance and Scaling

## Producer Performance

For low throughput, the normal `send()` API is simple. For high throughput, avoid flushing after every message and use the package's async or batch publishing support when available in your installed version.

Useful producer options:

```php
->withConfigOptions([
    'compression.type' => 'lz4',
    'linger.ms' => '10',
    'batch.num.messages' => '1000',
    'queue.buffering.max.messages' => '100000',
])
```

---

## Consumer Scaling

Consumer parallelism is bounded by partitions.

```text
Topic has 6 partitions
Consumer group has 3 instances
Each instance receives about 2 partitions
```

If a topic has 3 partitions, running 10 consumers in the same group will still process with only 3 active consumers.

---

## Handler Performance

Keep handlers lean:

- validate payload
- call application service
- persist result
- commit offset

Move slow side effects to Laravel jobs when the processing contract allows it.

---

## Backpressure

If lag grows:

1. Check handler error rate.
2. Check database or HTTP dependency latency.
3. Increase partitions if ordering allows.
4. Add consumer instances.
5. Batch work where safe.
6. Split heavy consumers into separate topics or services.
