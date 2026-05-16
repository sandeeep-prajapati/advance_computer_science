# Chapter 7: Kafka Streams

## What is Kafka Streams?

Kafka Streams is a **Java/Kotlin library** for building real-time stream processing applications on top of Kafka. Unlike Spark or Flink, it runs inside your application process — no separate cluster needed.

```
Input Topic ──► Kafka Streams App ──► Output Topic
                  (your process)
```

**Key properties:**
- Library (not a framework) — embedded in your app
- No external cluster — just Kafka
- Fault-tolerant with local RocksDB state stores
- Exactly-once semantics support
- Scales by adding instances

---

## Core Abstractions

| Abstraction | Description |
|---|---|
| **KStream** | Unbounded stream of records (insert semantics) |
| **KTable** | Changelog stream — latest value per key (update semantics) |
| **GlobalKTable** | Fully replicated KTable on every instance |
| **StreamsBuilder** | Entry point to build the topology |
| **Topology** | The processing graph (DAG of operations) |
| **State Store** | Local RocksDB for aggregations and joins |

---

## Basic Kafka Streams App (Java)

```java
import org.apache.kafka.streams.*;
import org.apache.kafka.streams.kstream.*;
import java.util.Properties;

public class WordCount {
    public static void main(String[] args) throws InterruptedException {
        Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG,    "wordcount-app");
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG,   Serdes.String().getClass());
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, Serdes.String().getClass());

        StreamsBuilder builder = new StreamsBuilder();

        KStream<String, String> textLines = builder.stream("text-input");

        KTable<String, Long> wordCounts = textLines
            .flatMapValues(line -> Arrays.asList(line.toLowerCase().split("\\W+")))
            .groupBy((key, word) -> word)
            .count();

        wordCounts.toStream().to("wordcount-output",
            Produced.with(Serdes.String(), Serdes.Long()));

        KafkaStreams streams = new KafkaStreams(builder.build(), props);
        streams.start();

        Runtime.getRuntime().addShutdownHook(new Thread(streams::close));
    }
}
```

---

## KStream Operations

```java
StreamsBuilder builder = new StreamsBuilder();
KStream<String, Order> orders = builder.stream("orders",
    Consumed.with(Serdes.String(), orderSerde));

// Filter
KStream<String, Order> highValue = orders
    .filter((key, order) -> order.getAmount() > 100.0);

// Map values
KStream<String, OrderEnriched> enriched = orders
    .mapValues(order -> enrich(order));

// FlatMap (one-to-many)
KStream<String, LineItem> lineItems = orders
    .flatMapValues(order -> order.getLineItems());

// Branch (split stream)
Map<String, KStream<String, Order>> branches = orders
    .split(Named.as("branch-"))
    .branch((k, v) -> v.getAmount() > 1000, Branched.as("premium"))
    .branch((k, v) -> v.getAmount() > 100,  Branched.as("standard"))
    .defaultBranch(Branched.as("small"));

KStream<String, Order> premium  = branches.get("branch-premium");
KStream<String, Order> standard = branches.get("branch-standard");

// Repartition (change key)
KStream<String, Order> byProduct = orders
    .selectKey((k, order) -> order.getProductId());

// Write to topic
highValue.to("high-value-orders", Produced.with(Serdes.String(), orderSerde));
```

---

## Aggregations

```java
// Count events per user per hour
KTable<Windowed<String>, Long> clicksPerUser = orders
    .groupByKey()
    .windowedBy(TimeWindows.ofSizeWithNoGrace(Duration.ofHours(1)))
    .count(Materialized.as("clicks-per-user-store"));

clicksPerUser
    .toStream()
    .map((windowedKey, count) -> KeyValue.pair(
        windowedKey.key(),
        windowedKey.window().start() + ":" + count
    ))
    .to("clicks-per-user-output");

// Sum revenue per category
KTable<String, Double> revenueByCategory = orders
    .groupBy((key, order) -> KeyValue.pair(order.getCategory(), order))
    .aggregate(
        () -> 0.0,
        (key, order, total) -> total + order.getAmount(),
        Materialized.<String, Double, KeyValueStore<Bytes, byte[]>>as("revenue-store")
            .withValueSerde(Serdes.Double())
    );
```

---

## Windowing

```java
// Tumbling window (non-overlapping, fixed size)
TimeWindows tumbling = TimeWindows.ofSizeWithNoGrace(Duration.ofMinutes(5));

// Sliding window (overlapping by step)
SlidingWindows sliding = SlidingWindows.ofTimeDifferenceAndGrace(
    Duration.ofMinutes(10),
    Duration.ofMinutes(1)
);

// Session window (gap-based)
SessionWindows session = SessionWindows.ofInactivityGapWithNoGrace(Duration.ofMinutes(30));

// Usage
KTable<Windowed<String>, Long> result = stream
    .groupByKey()
    .windowedBy(tumbling)
    .count(Materialized.as("windowed-count-store"));
```

---

## KStream-KTable Join

```java
KStream<String, Order>   orders   = builder.stream("orders");
KTable<String, Customer> customers = builder.table("customers");

// Enrich stream with table lookup
KStream<String, EnrichedOrder> enriched = orders
    .join(customers,
        (order, customer) -> new EnrichedOrder(order, customer),
        Joined.with(Serdes.String(), orderSerde, customerSerde)
    );

enriched.to("enriched-orders");
```

---

## KStream-KStream Join

```java
KStream<String, Order>   orders   = builder.stream("orders");
KStream<String, Payment> payments = builder.stream("payments");

// Join two streams within a 5-minute window
KStream<String, OrderWithPayment> combined = orders
    .join(payments,
        (order, payment) -> new OrderWithPayment(order, payment),
        JoinWindows.ofTimeDifferenceWithNoGrace(Duration.ofMinutes(5)),
        StreamJoined.with(Serdes.String(), orderSerde, paymentSerde)
    );
```

---

## Stateful Processing with State Stores

```java
// Custom state store for deduplication
StoreBuilder<KeyValueStore<String, Boolean>> storeBuilder =
    Stores.keyValueStoreBuilder(
        Stores.persistentKeyValueStore("dedup-store"),
        Serdes.String(),
        Serdes.Boolean()
    );

builder.addStateStore(storeBuilder);

KStream<String, Order> deduplicated = orders
    .transform(() -> new Transformer<String, Order, KeyValue<String, Order>>() {
        KeyValueStore<String, Boolean> store;

        @Override
        public void init(ProcessorContext context) {
            store = context.getStateStore("dedup-store");
        }

        @Override
        public KeyValue<String, Order> transform(String key, Order value) {
            if (store.get(value.getOrderId()) != null) {
                return null;  // duplicate, drop
            }
            store.put(value.getOrderId(), true);
            return KeyValue.pair(key, value);
        }
    }, "dedup-store");
```

---

## Interactive Queries

Query local state stores from outside the app:

```java
KafkaStreams streams = new KafkaStreams(builder.build(), props);
streams.start();

// Query local state store directly
ReadOnlyKeyValueStore<String, Long> store =
    streams.store(StoreQueryParameters.fromNameAndType(
        "wordcount-store",
        QueryableStoreTypes.keyValueStore()
    ));

Long count = store.get("kafka");
```

---

## Kafka Streams Configuration

```java
Properties props = new Properties();
props.put(StreamsConfig.APPLICATION_ID_CONFIG,    "my-streams-app");
props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");

// Exactly-once semantics
props.put(StreamsConfig.PROCESSING_GUARANTEE_CONFIG, StreamsConfig.EXACTLY_ONCE_V2);

// Parallelism
props.put(StreamsConfig.NUM_STREAM_THREADS_CONFIG, 4);

// State directory (RocksDB)
props.put(StreamsConfig.STATE_DIR_CONFIG, "/tmp/kafka-streams");

// Commit interval
props.put(StreamsConfig.COMMIT_INTERVAL_MS_CONFIG, 100);

// Replication factor for internal topics
props.put(StreamsConfig.REPLICATION_FACTOR_CONFIG, 3);
```

---

## Kafka Streams vs Other Frameworks

| Feature | Kafka Streams | Apache Flink | Apache Spark Streaming |
|---|---|---|---|
| Deployment | Embedded library | Separate cluster | Separate cluster |
| Latency | Low (ms) | Very low (ms) | Moderate (seconds) |
| State backend | RocksDB | RocksDB / heap | In-memory / HDFS |
| Exactly-once | Yes | Yes | Yes |
| Windowing | Yes | Rich | Yes |
| Language | Java / Kotlin | Java / Scala / Python | Python / Java / Scala |
| SQL interface | No | Flink SQL | Spark SQL |
| Learning curve | Moderate | Steep | Moderate |
