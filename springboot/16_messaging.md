# Chapter 16: Messaging — Kafka & RabbitMQ

## Why Messaging?

Messaging decouples services — the sender doesn't need to know if the receiver is running, how many receivers exist, or how fast they process. This enables:
- **Asynchronous processing** — respond to the user immediately, process in background
- **Load leveling** — messages queue up during spikes, consumers work at their pace
- **Event-driven architecture** — services react to events, not direct calls

---

## Apache Kafka

Kafka is a distributed, fault-tolerant event streaming platform. Best for high-throughput, event replay, and ordered processing.

### Dependency:
```xml
<dependency>
    <groupId>org.springframework.kafka</groupId>
    <artifactId>spring-kafka</artifactId>
</dependency>
```

### Configuration:
```yaml
spring:
  kafka:
    bootstrap-servers: localhost:9092
    producer:
      key-serializer: org.apache.kafka.common.serialization.StringSerializer
      value-serializer: org.springframework.kafka.support.serializer.JsonSerializer
      acks: all          # wait for all replicas to acknowledge
      retries: 3
      properties:
        enable.idempotence: true   # exactly-once semantics
    consumer:
      group-id: my-service
      key-deserializer: org.apache.kafka.common.serialization.StringDeserializer
      value-deserializer: org.springframework.kafka.support.serializer.JsonDeserializer
      auto-offset-reset: earliest  # start from beginning for new consumer groups
      enable-auto-commit: false    # manual acknowledgment for reliability
      properties:
        spring.json.trusted.packages: "com.example.events"
    listener:
      ack-mode: manual   # manually commit offset after processing
```

---

## Kafka Producer

```java
// Event class
public record OrderCreatedEvent(
    Long orderId,
    Long userId,
    BigDecimal totalAmount,
    List<Long> itemIds,
    LocalDateTime createdAt
) {}

@Service
public class OrderService {

    private final KafkaTemplate<String, OrderCreatedEvent> kafkaTemplate;

    @Transactional
    public Order placeOrder(OrderRequest request) {
        Order order = orderRepository.save(new Order(request));

        OrderCreatedEvent event = new OrderCreatedEvent(
            order.getId(), order.getUserId(), order.getTotal(),
            order.getItemIds(), order.getCreatedAt()
        );

        // Send async — returns CompletableFuture
        kafkaTemplate.send("orders.created", String.valueOf(order.getId()), event)
            .whenComplete((result, ex) -> {
                if (ex != null) {
                    log.error("Failed to send order event: orderId={}", order.getId(), ex);
                } else {
                    log.info("Order event sent: partition={}, offset={}",
                        result.getRecordMetadata().partition(),
                        result.getRecordMetadata().offset());
                }
            });

        return order;
    }
}
```

---

## Kafka Consumer

```java
@Component
public class NotificationConsumer {

    private final EmailService emailService;

    @KafkaListener(
        topics = "orders.created",
        groupId = "notification-service",
        containerFactory = "kafkaListenerContainerFactory"
    )
    public void handleOrderCreated(
            ConsumerRecord<String, OrderCreatedEvent> record,
            Acknowledgment acknowledgment) {
        try {
            OrderCreatedEvent event = record.value();
            log.info("Received order event: orderId={}", event.orderId());

            emailService.sendOrderConfirmation(event.userId(), event.orderId());

            acknowledgment.acknowledge();   // commit offset only after successful processing

        } catch (Exception e) {
            log.error("Failed to process order event: orderId={}", record.key(), e);
            // Don't acknowledge — message will be redelivered
            throw e;
        }
    }

    // Listen to multiple topics
    @KafkaListener(topics = {"orders.created", "orders.updated"}, groupId = "audit-service")
    public void auditOrderEvents(OrderCreatedEvent event) {
        auditRepository.save(new AuditLog("ORDER_EVENT", event.orderId()));
    }

    // Batch listener
    @KafkaListener(topics = "analytics.events", groupId = "analytics")
    public void processBatch(List<ConsumerRecord<String, AnalyticsEvent>> records,
                             Acknowledgment ack) {
        analyticsRepository.saveAll(records.stream().map(r -> r.value()).toList());
        ack.acknowledge();
    }
}
```

---

## Kafka Configuration Bean

```java
@Configuration
public class KafkaConfig {

    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, Object> kafkaListenerContainerFactory(
            ConsumerFactory<String, Object> consumerFactory) {
        ConcurrentKafkaListenerContainerFactory<String, Object> factory =
            new ConcurrentKafkaListenerContainerFactory<>();
        factory.setConsumerFactory(consumerFactory);
        factory.setConcurrency(3);   // 3 threads per listener (up to partition count)
        factory.getContainerProperties().setAckMode(ContainerProperties.AckMode.MANUAL);

        // Dead letter queue for failed messages
        DefaultErrorHandler errorHandler = new DefaultErrorHandler(
            new DeadLetterPublishingRecoverer(kafkaTemplate),
            new FixedBackOff(1000L, 3)   // retry 3 times with 1s delay
        );
        factory.setCommonErrorHandler(errorHandler);

        return factory;
    }

    @Bean
    public NewTopic orderCreatedTopic() {
        return TopicBuilder.name("orders.created")
            .partitions(6)
            .replicas(3)
            .build();
    }
}
```

---

## RabbitMQ

RabbitMQ is a message broker using the AMQP protocol. Best for task queues, request-reply patterns, and complex routing.

### Dependency:
```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-amqp</artifactId>
</dependency>
```

### Configuration:
```yaml
spring:
  rabbitmq:
    host: localhost
    port: 5672
    username: guest
    password: guest
    virtual-host: /
    listener:
      simple:
        acknowledge-mode: manual
        retry:
          enabled: true
          max-attempts: 3
          initial-interval: 1000ms
```

### Setup Queues, Exchanges, Bindings:

```java
@Configuration
public class RabbitMQConfig {

    public static final String ORDER_QUEUE = "orders.queue";
    public static final String ORDER_EXCHANGE = "orders.exchange";
    public static final String ORDER_ROUTING_KEY = "order.created";
    public static final String DLQ = "orders.dlq";
    public static final String DLX = "orders.dlx";

    // Dead Letter Queue
    @Bean
    public Queue deadLetterQueue() {
        return QueueBuilder.durable(DLQ).build();
    }

    @Bean
    public DirectExchange deadLetterExchange() {
        return new DirectExchange(DLX);
    }

    @Bean
    public Binding dlqBinding() {
        return BindingBuilder.bind(deadLetterQueue()).to(deadLetterExchange()).with(DLQ);
    }

    // Main queue with DLQ configured
    @Bean
    public Queue orderQueue() {
        return QueueBuilder.durable(ORDER_QUEUE)
            .withArgument("x-dead-letter-exchange", DLX)
            .withArgument("x-dead-letter-routing-key", DLQ)
            .withArgument("x-message-ttl", 300000)   // 5 min TTL
            .build();
    }

    @Bean
    public TopicExchange orderExchange() {
        return new TopicExchange(ORDER_EXCHANGE);
    }

    @Bean
    public Binding orderBinding() {
        return BindingBuilder.bind(orderQueue()).to(orderExchange()).with(ORDER_ROUTING_KEY);
    }

    @Bean
    public MessageConverter jsonMessageConverter() {
        return new Jackson2JsonMessageConverter();
    }

    @Bean
    public RabbitTemplate rabbitTemplate(ConnectionFactory factory) {
        RabbitTemplate template = new RabbitTemplate(factory);
        template.setMessageConverter(jsonMessageConverter());
        template.setConfirmCallback((correlationData, ack, cause) -> {
            if (!ack) log.error("Message not confirmed: {}", cause);
        });
        return template;
    }
}
```

### RabbitMQ Producer:
```java
@Service
public class OrderService {

    private final RabbitTemplate rabbitTemplate;

    public Order placeOrder(OrderRequest request) {
        Order order = orderRepository.save(new Order(request));

        rabbitTemplate.convertAndSend(
            RabbitMQConfig.ORDER_EXCHANGE,
            RabbitMQConfig.ORDER_ROUTING_KEY,
            new OrderCreatedEvent(order.getId(), order.getUserId())
        );

        return order;
    }
}
```

### RabbitMQ Consumer:
```java
@Component
public class NotificationConsumer {

    @RabbitListener(queues = RabbitMQConfig.ORDER_QUEUE)
    public void handleOrderCreated(
            OrderCreatedEvent event,
            Channel channel,
            @Header(AmqpHeaders.DELIVERY_TAG) long tag) throws IOException {
        try {
            notificationService.sendOrderConfirmation(event.getUserId(), event.getOrderId());
            channel.basicAck(tag, false);          // acknowledge
        } catch (Exception e) {
            channel.basicNack(tag, false, false);  // reject, send to DLQ
            log.error("Failed to process order: {}", event.getOrderId(), e);
        }
    }
}
```

---

## Kafka vs RabbitMQ Comparison

| Feature | Kafka | RabbitMQ |
|---|---|---|
| Model | Event log (pull) | Message queue (push) |
| Message retention | Configurable (days/weeks) | Until consumed |
| Replay | Yes — rewind offset | No — consumed messages gone |
| Throughput | Very high (millions/sec) | High (tens of thousands/sec) |
| Ordering | Per partition | Per queue |
| Routing | Topics/partitions | Exchanges/routing keys |
| Consumer groups | Multiple independent groups | Competing consumers |
| Best for | Event sourcing, analytics, audit | Task queues, RPC, complex routing |

---

## Transactional Messaging (Outbox Pattern)

Ensures messages are sent **if and only if** the database transaction commits:

```java
@Entity
public class OutboxEvent {
    @Id @GeneratedValue
    private Long id;
    private String topic;
    private String key;
    private String payload;
    private boolean sent = false;
    private LocalDateTime createdAt;
}

@Service
public class OrderService {
    @Transactional
    public Order placeOrder(OrderRequest request) {
        Order order = orderRepository.save(new Order(request));

        // Save event atomically with the order
        OutboxEvent outbox = new OutboxEvent("orders.created",
            String.valueOf(order.getId()),
            objectMapper.writeValueAsString(new OrderCreatedEvent(order)));
        outboxRepository.save(outbox);

        return order;
    }
}

// Separate scheduler reads and sends outbox events
@Scheduled(fixedDelay = 1000)
public void processOutbox() {
    outboxRepository.findBySentFalse().forEach(event -> {
        kafkaTemplate.send(event.getTopic(), event.getKey(), event.getPayload());
        event.setSent(true);
        outboxRepository.save(event);
    });
}
```
