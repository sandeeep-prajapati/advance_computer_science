# Chapter 14: Spring Cloud

## What is Spring Cloud?

**Spring Cloud** provides tools for building distributed systems — service discovery, centralized config, circuit breakers, distributed tracing, and more.

```xml
<!-- BOM — manages all Spring Cloud versions -->
<dependencyManagement>
    <dependencies>
        <dependency>
            <groupId>org.springframework.cloud</groupId>
            <artifactId>spring-cloud-dependencies</artifactId>
            <version>2023.0.1</version>
            <type>pom</type>
            <scope>import</scope>
        </dependency>
    </dependencies>
</dependencyManagement>
```

---

## Service Discovery with Eureka

**Eureka** lets services find each other by name instead of hardcoded host:port.

### Eureka Server:
```xml
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-starter-netflix-eureka-server</artifactId>
</dependency>
```

```java
@SpringBootApplication
@EnableEurekaServer
public class DiscoveryServerApplication { ... }
```

```yaml
# application.yml (Eureka Server)
server:
  port: 8761
eureka:
  client:
    register-with-eureka: false
    fetch-registry: false
```

### Eureka Client (each microservice):
```xml
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-starter-netflix-eureka-client</artifactId>
</dependency>
```

```yaml
# application.yml (each service)
spring:
  application:
    name: user-service    # registered name in Eureka
eureka:
  client:
    service-url:
      defaultZone: http://eureka-server:8761/eureka/
```

```java
// RestClient with load balancing via Eureka
@Bean
@LoadBalanced    // resolves service names via Eureka
public RestClient.Builder restClientBuilder() {
    return RestClient.builder();
}

// Usage: "http://user-service" resolves to actual instance
restClient.get().uri("http://user-service/api/users/{id}", id).retrieve().body(UserDto.class);
```

---

## OpenFeign — Declarative REST Clients

Feign generates REST client implementations from interfaces:

```xml
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-starter-openfeign</artifactId>
</dependency>
```

```java
@SpringBootApplication
@EnableFeignClients
public class OrderServiceApplication { ... }

@FeignClient(name = "user-service", path = "/api/users")
public interface UserServiceClient {

    @GetMapping("/{id}")
    UserDto getUserById(@PathVariable Long id);

    @GetMapping
    Page<UserDto> getAllUsers(@RequestParam int page, @RequestParam int size);

    @PostMapping
    UserDto createUser(@RequestBody CreateUserRequest request);
}

@Service
public class OrderService {
    private final UserServiceClient userClient;

    public OrderDto placeOrder(Long userId, OrderRequest request) {
        UserDto user = userClient.getUserById(userId);  // looks like a local call!
        // ...
    }
}
```

---

## Spring Cloud Config Server

Centralize configuration for all services in a Git repository.

### Config Server:
```xml
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-config-server</artifactId>
</dependency>
```

```java
@SpringBootApplication
@EnableConfigServer
public class ConfigServerApplication { ... }
```

```yaml
# Config Server application.yml
spring:
  cloud:
    config:
      server:
        git:
          uri: https://github.com/myorg/config-repo
          default-label: main
          search-paths: '{application}'
```

Git repo structure:
```
config-repo/
├── application.yml           # shared by all services
├── user-service.yml          # user-service specific
├── user-service-prod.yml     # user-service in prod
└── order-service.yml
```

### Config Client (each service):
```xml
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-starter-config</artifactId>
</dependency>
```

```yaml
# bootstrap.yml (loaded before application.yml)
spring:
  application:
    name: user-service
  config:
    import: configserver:http://config-server:8888
  profiles:
    active: prod
```

**Refresh config without restart:**
```java
@RefreshScope   // re-creates bean when /actuator/refresh is called
@Service
public class FeatureFlagService {
    @Value("${feature.new-checkout: false}")
    private boolean newCheckoutEnabled;
}
```

```bash
# Trigger refresh after pushing new config to Git
curl -X POST http://user-service:8080/actuator/refresh
```

---

## Resilience4j — Circuit Breaker

Circuit Breaker prevents cascade failures when a downstream service is down.

```xml
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-starter-circuitbreaker-resilience4j</artifactId>
</dependency>
```

### States:
```
CLOSED (normal) → 50% failures → OPEN (rejects all) → wait 60s → HALF_OPEN → success? → CLOSED
```

```java
@Service
public class OrderService {

    private final CircuitBreakerFactory circuitBreakerFactory;

    public UserDto getUser(Long userId) {
        CircuitBreaker cb = circuitBreakerFactory.create("user-service");
        return cb.run(
            () -> userServiceClient.getUserById(userId),   // primary
            throwable -> getDefaultUser(userId)            // fallback
        );
    }

    private UserDto getDefaultUser(Long userId) {
        return new UserDto(userId, "Unknown", "unknown@example.com");
    }
}
```

### Annotation approach:
```java
@CircuitBreaker(name = "user-service", fallbackMethod = "getDefaultUser")
@Retry(name = "user-service")
@RateLimiter(name = "user-service")
@TimeLimiter(name = "user-service")
public UserDto getUser(Long userId) {
    return userServiceClient.getUserById(userId);
}

public UserDto getDefaultUser(Long userId, Exception e) {
    log.warn("Fallback for userId={}, reason={}", userId, e.getMessage());
    return new UserDto(userId, "Unknown", "");
}
```

### Configuration:
```yaml
resilience4j:
  circuitbreaker:
    instances:
      user-service:
        sliding-window-size: 10
        minimum-number-of-calls: 5
        failure-rate-threshold: 50        # open circuit if 50% fail
        wait-duration-in-open-state: 60s
        permitted-number-of-calls-in-half-open-state: 3
  retry:
    instances:
      user-service:
        max-attempts: 3
        wait-duration: 500ms
        retry-exceptions:
          - java.net.ConnectException
  timelimiter:
    instances:
      user-service:
        timeout-duration: 3s
```

---

## Distributed Tracing

### Micrometer Tracing (Spring Boot 3+):
```xml
<dependency>
    <groupId>io.micrometer</groupId>
    <artifactId>micrometer-tracing-bridge-brave</artifactId>
</dependency>
<dependency>
    <groupId>io.zipkin.reporter2</groupId>
    <artifactId>zipkin-reporter-brave</artifactId>
</dependency>
```

```yaml
management:
  tracing:
    sampling:
      probability: 0.1   # sample 10% in production
  zipkin:
    tracing:
      endpoint: http://zipkin:9411/api/v2/spans
```

Trace IDs automatically propagate through:
- RestClient / RestTemplate / WebClient
- Feign clients
- Kafka messages
- Scheduled tasks

Logs automatically include trace/span IDs:
```
2024-01-15 10:30:00 [user-service,6532d5d8e83e3aa2,6532d5d8e83e3aa2] INFO  - Processing request
```

---

## Spring Cloud Gateway Advanced

```java
@Bean
public RouteLocator customRouteLocator(RouteLocatorBuilder builder) {
    return builder.routes()
        .route("user-service", r -> r
            .path("/api/users/**")
            .filters(f -> f
                .addRequestHeader("X-Request-Source", "gateway")
                .circuitBreaker(c -> c
                    .setName("user-service")
                    .setFallbackUri("forward:/fallback/users"))
            )
            .uri("lb://user-service"))
        .build();
}

// Fallback controller
@RestController
@RequestMapping("/fallback")
public class FallbackController {
    @GetMapping("/users")
    public ResponseEntity<String> usersFallback() {
        return ResponseEntity.serviceUnavailable()
            .body("User service is temporarily unavailable");
    }
}
```

---

## Spring Cloud Components Summary

| Component | Purpose | Implementation |
|---|---|---|
| **Eureka** | Service registry/discovery | `spring-cloud-starter-netflix-eureka` |
| **Config Server** | Centralized configuration | `spring-cloud-config-server` |
| **Gateway** | API Gateway, routing | `spring-cloud-starter-gateway` |
| **OpenFeign** | Declarative REST clients | `spring-cloud-starter-openfeign` |
| **Resilience4j** | Circuit breaker, retry | `spring-cloud-starter-circuitbreaker-resilience4j` |
| **Micrometer Tracing** | Distributed tracing | `micrometer-tracing-bridge-brave` |
| **Bus** | Config refresh propagation | `spring-cloud-starter-bus-kafka` |
| **Kubernetes** | K8s-native service discovery | `spring-cloud-starter-kubernetes` |

---

## Alternative: Kubernetes-Native Discovery

Instead of Eureka, use K8s Services for service discovery:

```xml
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-starter-kubernetes-client-all</artifactId>
</dependency>
```

Services discover each other via K8s DNS:
- `http://user-service` resolves to the K8s Service
- No Eureka server needed
- Works natively in K8s environments
