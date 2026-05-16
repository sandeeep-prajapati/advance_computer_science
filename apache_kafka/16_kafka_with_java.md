# Chapter 16: Kafka with Java

## Maven Dependency

```xml
<dependency>
    <groupId>org.apache.kafka</groupId>
    <artifactId>kafka-clients</artifactId>
    <version>3.7.0</version>
</dependency>

<!-- Spring Kafka -->
<dependency>
    <groupId>org.springframework.kafka</groupId>
    <artifactId>spring-kafka</artifactId>
    <version>3.1.0</version>
</dependency>
```

---

## Java Producer

```java
import org.apache.kafka.clients.producer.*;
import org.apache.kafka.common.serialization.StringSerializer;
import java.util.Properties;
import java.util.concurrent.ExecutionException;

public class OrderProducer {

    public static void main(String[] args) throws ExecutionException, InterruptedException {

        Properties props = new Properties();
        props.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG,  "localhost:9092");
        props.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG,   StringSerializer.class.getName());
        props.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        props.put(ProducerConfig.ACKS_CONFIG,               "all");
        props.put(ProducerConfig.ENABLE_IDEMPOTENCE_CONFIG, true);
        props.put(ProducerConfig.LINGER_MS_CONFIG,          5);
        props.put(ProducerConfig.BATCH_SIZE_CONFIG,         65536);
        props.put(ProducerConfig.COMPRESSION_TYPE_CONFIG,   "lz4");

        try (KafkaProducer<String, String> producer = new KafkaProducer<>(props)) {

            // Async send with callback
            ProducerRecord<String, String> record = new ProducerRecord<>(
                "orders",
                "user-123",
                "{\"order_id\":\"A1\",\"amount\":99.99}"
            );

            producer.send(record, (metadata, exception) -> {
                if (exception != null) {
                    System.err.println("Send failed: " + exception.getMessage());
                } else {
                    System.out.printf("Sent to %s[%d]@%d%n",
                        metadata.topic(), metadata.partition(), metadata.offset());
                }
            });

            // Synchronous send
            RecordMetadata meta = producer.send(record).get();
            System.out.printf("Sync send: partition=%d offset=%d%n",
                meta.partition(), meta.offset());

            producer.flush();
        }
    }
}
```

---

## Java Consumer

```java
import org.apache.kafka.clients.consumer.*;
import org.apache.kafka.common.serialization.StringDeserializer;
import java.time.Duration;
import java.util.*;

public class OrderConsumer {

    public static void main(String[] args) {
        Properties props = new Properties();
        props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG,   "localhost:9092");
        props.put(ConsumerConfig.GROUP_ID_CONFIG,            "order-processor");
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG,   StringDeserializer.class.getName());
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG,   "earliest");
        props.put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG,  false);
        props.put(ConsumerConfig.MAX_POLL_RECORDS_CONFIG,    500);

        try (KafkaConsumer<String, String> consumer = new KafkaConsumer<>(props)) {
            consumer.subscribe(List.of("orders"));

            Runtime.getRuntime().addShutdownHook(new Thread(consumer::wakeup));

            try {
                while (true) {
                    ConsumerRecords<String, String> records =
                        consumer.poll(Duration.ofMillis(1000));

                    for (ConsumerRecord<String, String> record : records) {
                        System.out.printf("topic=%s partition=%d offset=%d key=%s value=%s%n",
                            record.topic(), record.partition(), record.offset(),
                            record.key(), record.value());

                        processOrder(record.key(), record.value());
                    }

                    // Commit after processing batch
                    consumer.commitSync();
                }
            } catch (WakeupException e) {
                // Shutdown signal, exit cleanly
            }
        }
    }

    static void processOrder(String key, String value) {
        // Business logic
    }
}
```

---

## Transactions in Java

```java
Properties props = new Properties();
props.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
props.put(ProducerConfig.TRANSACTIONAL_ID_CONFIG,  "order-enricher-1");
props.put(ProducerConfig.ENABLE_IDEMPOTENCE_CONFIG, true);
props.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG,   StringSerializer.class.getName());
props.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());

KafkaProducer<String, String> producer = new KafkaProducer<>(props);
producer.initTransactions();

Properties consumerProps = new Properties();
consumerProps.put(ConsumerConfig.ISOLATION_LEVEL_CONFIG, "read_committed");
// ... other consumer config ...
KafkaConsumer<String, String> consumer = new KafkaConsumer<>(consumerProps);
consumer.subscribe(List.of("raw-orders"));

while (true) {
    ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(100));
    if (records.isEmpty()) continue;

    try {
        producer.beginTransaction();

        for (ConsumerRecord<String, String> record : records) {
            String enriched = enrich(record.value());
            producer.send(new ProducerRecord<>("enriched-orders", record.key(), enriched));
        }

        // Commit offsets atomically with transaction
        Map<TopicPartition, OffsetAndMetadata> offsets = new HashMap<>();
        for (TopicPartition tp : consumer.assignment()) {
            offsets.put(tp, new OffsetAndMetadata(consumer.position(tp)));
        }
        producer.sendOffsetsToTransaction(offsets, consumer.groupMetadata());
        producer.commitTransaction();

    } catch (Exception e) {
        producer.abortTransaction();
        consumer.seekToBeginning(consumer.assignment());
    }
}
```

---

## Spring Kafka

### Configuration

```java
@Configuration
@EnableKafka
public class KafkaConfig {

    @Bean
    public ProducerFactory<String, Object> producerFactory() {
        Map<String, Object> config = new HashMap<>();
        config.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
        config.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG,   StringSerializer.class);
        config.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, JsonSerializer.class);
        config.put(ProducerConfig.ACKS_CONFIG, "all");
        config.put(ProducerConfig.ENABLE_IDEMPOTENCE_CONFIG, true);
        return new DefaultKafkaProducerFactory<>(config);
    }

    @Bean
    public KafkaTemplate<String, Object> kafkaTemplate() {
        return new KafkaTemplate<>(producerFactory());
    }

    @Bean
    public ConsumerFactory<String, Object> consumerFactory() {
        Map<String, Object> config = new HashMap<>();
        config.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
        config.put(ConsumerConfig.GROUP_ID_CONFIG,          "spring-group");
        config.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG,   StringDeserializer.class);
        config.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, JsonDeserializer.class);
        config.put(JsonDeserializer.TRUSTED_PACKAGES, "*");
        return new DefaultKafkaConsumerFactory<>(config);
    }

    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, Object> kafkaListenerContainerFactory() {
        ConcurrentKafkaListenerContainerFactory<String, Object> factory =
            new ConcurrentKafkaListenerContainerFactory<>();
        factory.setConsumerFactory(consumerFactory());
        factory.setConcurrency(3);  // 3 threads (one per partition)
        factory.getContainerProperties().setAckMode(AckMode.MANUAL_IMMEDIATE);
        return factory;
    }
}
```

### Producer Service

```java
@Service
public class OrderProducerService {

    @Autowired
    private KafkaTemplate<String, Object> kafkaTemplate;

    public void sendOrder(Order order) {
        kafkaTemplate.send("orders", order.getOrderId(), order)
            .whenComplete((result, ex) -> {
                if (ex != null) {
                    log.error("Failed to send order {}: {}", order.getOrderId(), ex.getMessage());
                } else {
                    log.info("Order sent: partition={} offset={}",
                        result.getRecordMetadata().partition(),
                        result.getRecordMetadata().offset());
                }
            });
    }
}
```

### Consumer (Listener)

```java
@Component
public class OrderConsumer {

    @KafkaListener(
        topics   = "orders",
        groupId  = "order-processor",
        concurrency = "3"   // 3 consumer threads
    )
    public void consume(
        @Payload Order order,
        @Header(KafkaHeaders.RECEIVED_PARTITION) int partition,
        @Header(KafkaHeaders.OFFSET) long offset,
        Acknowledgment ack
    ) {
        try {
            processOrder(order);
            ack.acknowledge();   // manual commit
        } catch (Exception e) {
            log.error("Failed to process order at offset {}: {}", offset, e.getMessage());
            // Don't ack — will be redelivered
        }
    }

    // Batch listener
    @KafkaListener(topics = "events", containerFactory = "batchFactory")
    public void consumeBatch(List<Order> orders, Acknowledgment ack) {
        processBatch(orders);
        ack.acknowledge();
    }
}
```

---

## Error Handling and Retry in Spring Kafka

```java
@Bean
public DefaultErrorHandler errorHandler(KafkaOperations<?, ?> template) {
    // Fixed backoff: 3 retries, 1 second apart
    FixedBackOff backOff = new FixedBackOff(1000L, 3);

    DeadLetterPublishingRecoverer recoverer =
        new DeadLetterPublishingRecoverer(template,
            (rec, ex) -> new TopicPartition(rec.topic() + ".dlq", rec.partition()));

    DefaultErrorHandler handler = new DefaultErrorHandler(recoverer, backOff);

    // Don't retry on these exceptions
    handler.addNotRetryableExceptions(
        DeserializationException.class,
        NullPointerException.class
    );

    return handler;
}
```
