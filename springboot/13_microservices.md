# Chapter 13: Spring Boot Microservices

## What are Microservices?

**Microservices** is an architectural style where an application is composed of small, independently deployable services, each responsible for a specific business capability.

```
Monolith:                          Microservices:
┌──────────────────────────┐       ┌──────────┐  ┌──────────┐  ┌──────────┐
│   User Management        │       │  User    │  │  Order   │  │ Payment  │
│   Order Management       │  →    │  Service │  │  Service │  │  Service │
│   Payment Processing     │       └──────────┘  └──────────┘  └──────────┘
│   Inventory              │             │             │             │
│   Notifications          │       ┌──────────┐  ┌──────────┐       │
└──────────────────────────┘       │Inventory │  │Notif.    │       │
                                   │  Service │  │  Service │       │
                                   └──────────┘  └──────────┘       │
```

---

## When to Use Microservices

**Use microservices when:**
- Different parts of the app scale differently (payment vs catalog)
- Different teams own different domains
- Different release cycles for different features
- Technology diversity is needed (Python for ML, Java for core)

**Don't use microservices when:**
- Team is small (< 10 engineers)
- Domain boundaries aren't clear yet
- You don't have strong DevOps capabilities
- Starting a new product (start as a monolith, extract later)

---

## Decomposition Strategies

| Strategy | How | Example |
|---|---|---|
| By business capability | Split by what the business does | Orders, Inventory, Billing |
| By subdomain (DDD) | Domain-Driven Design bounded contexts | Catalogue, Fulfillment, Customer |
| By team | Conway's Law — structure matches org | Team A → User Service, Team B → Order Service |
| Strangler Fig | Gradually extract from monolith | Start with notifications, then orders |

---

## Sample Microservices Architecture

```
Client (Browser/Mobile)
        │
        ▼
┌───────────────┐
│  API Gateway  │  (Spring Cloud Gateway / Kong / Nginx)
└───────┬───────┘
        │ routes requests
   ┌────┴──────────────────────────┐
   │           │                  │
   ▼           ▼                  ▼
┌──────┐  ┌──────────┐  ┌──────────────┐
│ User │  │  Order   │  │   Payment    │
│ Svc  │  │  Service │  │   Service    │
└──┬───┘  └────┬─────┘  └──────┬───────┘
   │            │               │
   ▼            ▼               ▼
┌──────┐  ┌──────────┐  ┌──────────────┐
│User  │  │ Orders   │  │  Payments    │
│  DB  │  │    DB    │  │     DB       │
└──────┘  └──────────┘  └──────────────┘
                │
           (publishes events)
                │
          ┌─────▼──────┐
          │    Kafka    │   (async messaging)
          └─────┬───────┘
                │
         ┌──────▼───────┐
         │ Notification │
         │   Service    │
         └──────────────┘
```

---

## Inter-Service Communication

### Synchronous (REST / gRPC)
Use when you need an immediate response.

```java
// RestClient (Spring Boot 3.2+)
@Service
public class OrderService {

    private final RestClient restClient;

    public OrderService(RestClient.Builder builder) {
        this.restClient = builder.baseUrl("http://user-service").build();
    }

    public UserDto getUser(Long userId) {
        return restClient.get()
            .uri("/api/users/{id}", userId)
            .retrieve()
            .body(UserDto.class);
    }
}
```

### Asynchronous (Kafka / RabbitMQ)
Use when you don't need an immediate response. Decouples services.

```java
// Order Service publishes event
@Service
public class OrderService {
    private final KafkaTemplate<String, OrderCreatedEvent> kafkaTemplate;

    public Order placeOrder(OrderRequest request) {
        Order order = orderRepository.save(new Order(request));
        kafkaTemplate.send("orders", new OrderCreatedEvent(order.getId(), order.getUserId()));
        return order;
    }
}

// Notification Service consumes event
@Component
public class OrderEventListener {
    @KafkaListener(topics = "orders", groupId = "notification-service")
    public void handleOrderCreated(OrderCreatedEvent event) {
        notificationService.sendOrderConfirmation(event.getUserId(), event.getOrderId());
    }
}
```

---

## API Gateway

The API Gateway is the single entry point for all clients. It handles:
- Routing to correct service
- Authentication (verify JWT once, not in every service)
- Rate limiting
- Request/response transformation
- Load balancing

### Spring Cloud Gateway:
```xml
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-starter-gateway</artifactId>
</dependency>
```

```yaml
spring:
  cloud:
    gateway:
      routes:
        - id: user-service
          uri: lb://user-service    # lb = load-balanced via Eureka
          predicates:
            - Path=/api/users/**
          filters:
            - StripPrefix=0
            - AddRequestHeader=X-Gateway-Source, api-gateway

        - id: order-service
          uri: lb://order-service
          predicates:
            - Path=/api/orders/**

        - id: rate-limit
          uri: lb://user-service
          predicates:
            - Path=/api/auth/**
          filters:
            - name: RequestRateLimiter
              args:
                redis-rate-limiter.replenishRate: 10
                redis-rate-limiter.burstCapacity: 20
```

---

## Service-to-Service Authentication

Each service should verify incoming requests are from trusted sources:

### Option 1: Pass JWT through
```java
// Gateway passes user's JWT to downstream services
// Downstream service validates the same JWT
```

### Option 2: Internal API Key
```java
// Each service has an internal API key
// Gateway adds X-Internal-Key header before forwarding
http.authorizeHttpRequests(auth -> auth
    .requestMatchers("/internal/**").hasRole("INTERNAL_SERVICE")
    .anyRequest().authenticated()
);
```

### Option 3: mTLS (production)
Services authenticate to each other using certificates — handled by service mesh (Istio).

---

## Database Per Service

Each service owns its own database — other services cannot query it directly.

```
User Service    → PostgreSQL (users DB)
Order Service   → PostgreSQL (orders DB)
Payment Service → PostgreSQL (payments DB)
Catalog Service → MongoDB (products DB)
Session Service → Redis
```

To share data across services:
- Expose via API (synchronous)
- Publish events to message bus (asynchronous)
- Replicate read-only data (CQRS)

---

## Saga Pattern — Distributed Transactions

Traditional ACID transactions don't work across microservices. The Saga pattern coordinates multi-step transactions:

**Choreography Saga (event-driven):**
```
Order Service   → publishes OrderCreated
                ← Inventory Service reserves stock → publishes StockReserved
                ← Payment Service charges card → publishes PaymentCompleted
                ← Fulfillment Service ships order → publishes OrderFulfilled

On failure:
Payment fails → publishes PaymentFailed
             ← Inventory Service releases stock (compensating transaction)
             ← Order Service marks order as failed
```

**Orchestration Saga (central coordinator):**
```
Order Saga Orchestrator:
1. Reserve inventory (call Inventory Service)
2. Charge payment (call Payment Service)
3. Ship order (call Fulfillment Service)
On any failure: call compensating transactions in reverse
```

---

## Shared Libraries vs Code Duplication

**Don't** share domain entities across services — this creates coupling.

**Do** share:
- Common utilities (logging format, error response model)
- OpenAPI-generated client stubs
- Security utilities

Create a shared library:
```xml
<!-- shared-lib/pom.xml -->
<artifactId>my-shared-lib</artifactId>

<!-- Each service's pom.xml -->
<dependency>
    <groupId>com.example</groupId>
    <artifactId>my-shared-lib</artifactId>
    <version>1.0.0</version>
</dependency>
```

---

## Health of the Whole System

Add readiness checks for dependencies:
```java
@Component
public class OrderServiceReadinessIndicator extends AbstractHealthIndicator {
    private final InventoryServiceClient inventoryClient;
    private final PaymentServiceClient paymentClient;

    @Override
    protected void doHealthCheck(Health.Builder builder) {
        boolean inventoryUp = inventoryClient.isHealthy();
        boolean paymentUp = paymentClient.isHealthy();

        if (inventoryUp && paymentUp) {
            builder.up().withDetail("dependencies", "all healthy");
        } else {
            builder.down()
                .withDetail("inventory", inventoryUp ? "UP" : "DOWN")
                .withDetail("payment", paymentUp ? "UP" : "DOWN");
        }
    }
}
```

---

## Key Design Principles

| Principle | What it means |
|---|---|
| Single Responsibility | Each service does one thing well |
| Loose Coupling | Services don't know each other's internals |
| High Cohesion | Related logic stays together |
| Design for failure | Assume services will fail; use timeouts, retries, circuit breakers |
| API-first | Design APIs before implementation |
| Automate everything | CI/CD, infra-as-code, automated testing |
