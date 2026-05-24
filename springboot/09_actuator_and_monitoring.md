# Chapter 9: Actuator & Monitoring

## What is Spring Boot Actuator?

**Spring Boot Actuator** adds production-ready features to your application — health checks, metrics, configuration info, thread dumps, and more — accessible via HTTP endpoints or JMX.

---

## Dependency

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-actuator</artifactId>
</dependency>
```

---

## Built-In Endpoints

| Endpoint | URL | Description |
|---|---|---|
| `health` | `/actuator/health` | Application and component health |
| `info` | `/actuator/info` | App metadata (version, git, etc.) |
| `metrics` | `/actuator/metrics` | JVM, HTTP, custom metrics |
| `env` | `/actuator/env` | All configuration properties |
| `beans` | `/actuator/beans` | All Spring beans |
| `mappings` | `/actuator/mappings` | All @RequestMapping routes |
| `loggers` | `/actuator/loggers` | View/change log levels at runtime |
| `threaddump` | `/actuator/threaddump` | JVM thread dump |
| `heapdump` | `/actuator/heapdump` | JVM heap dump (binary) |
| `shutdown` | `/actuator/shutdown` | Graceful shutdown (disabled by default) |
| `httptrace` | `/actuator/httptrace` | Recent HTTP requests |

---

## Exposing Endpoints

By default, only `health` is exposed over HTTP. Configure exposure:

```yaml
management:
  endpoints:
    web:
      exposure:
        include: health, info, metrics, loggers, env
        # or include: "*" to expose all (only in dev)
      base-path: /actuator     # default
  endpoint:
    health:
      show-details: always     # show component details
    shutdown:
      enabled: true            # enable graceful shutdown
  server:
    port: 8081                 # separate port for actuator (production best practice)
```

---

## Health Endpoint

```json
// GET /actuator/health
{
  "status": "UP",
  "components": {
    "db": {
      "status": "UP",
      "details": { "database": "PostgreSQL", "validationQuery": "isValid()" }
    },
    "redis": {
      "status": "UP",
      "details": { "version": "7.0.11" }
    },
    "diskSpace": {
      "status": "UP",
      "details": { "total": 499963174912, "free": 291275534336, "threshold": 10485760 }
    }
  }
}
```

### Custom Health Indicator:

```java
@Component
public class ExternalApiHealthIndicator implements HealthIndicator {

    private final ExternalApiClient client;

    @Override
    public Health health() {
        try {
            client.ping();
            return Health.up().withDetail("responseTime", "12ms").build();
        } catch (Exception e) {
            return Health.down()
                .withDetail("error", e.getMessage())
                .build();
        }
    }
}
```

---

## Info Endpoint

```yaml
management:
  info:
    env:
      enabled: true
    git:
      enabled: true
      mode: full
    build:
      enabled: true

info:
  app:
    name: My Application
    version: '@project.version@'    # pulled from pom.xml
    description: '@project.description@'
  contact:
    team: backend
    email: backend@example.com
```

For git info, add the plugin:
```xml
<plugin>
    <groupId>io.github.git-commit-id</groupId>
    <artifactId>git-commit-id-maven-plugin</artifactId>
</plugin>
```

---

## Metrics

Spring Boot uses **Micrometer** — a metrics facade that supports many backends.

```json
// GET /actuator/metrics
{
  "names": [
    "http.server.requests",
    "jvm.memory.used",
    "jvm.gc.pause",
    "process.cpu.usage",
    "spring.data.repository.invocations",
    "hikaricp.connections.active"
  ]
}

// GET /actuator/metrics/http.server.requests
{
  "name": "http.server.requests",
  "measurements": [
    { "statistic": "COUNT", "value": 1042 },
    { "statistic": "TOTAL_TIME", "value": 5.231 },
    { "statistic": "MAX", "value": 0.152 }
  ]
}
```

### Custom Metrics:

```java
@Service
public class OrderService {

    private final Counter orderCreatedCounter;
    private final Timer orderProcessingTimer;

    public OrderService(MeterRegistry registry) {
        this.orderCreatedCounter = registry.counter("orders.created", "status", "success");
        this.orderProcessingTimer = registry.timer("orders.processing.duration");
    }

    public Order placeOrder(OrderRequest request) {
        return orderProcessingTimer.record(() -> {
            Order order = processOrder(request);
            orderCreatedCounter.increment();
            return order;
        });
    }
}
```

Annotations approach:
```java
@Timed(value = "user.creation.time", description = "Time to create a user")
public UserDto createUser(CreateUserRequest request) { ... }
```

---

## Prometheus + Grafana Integration

### Add Micrometer Prometheus registry:
```xml
<dependency>
    <groupId>io.micrometer</groupId>
    <artifactId>micrometer-registry-prometheus</artifactId>
</dependency>
```

```yaml
management:
  endpoints:
    web:
      exposure:
        include: health, prometheus
  metrics:
    export:
      prometheus:
        enabled: true
    tags:
      application: ${spring.application.name}
      environment: ${spring.profiles.active}
```

Prometheus scrapes `GET /actuator/prometheus` in Prometheus format.

### Prometheus config (`prometheus.yml`):
```yaml
scrape_configs:
  - job_name: 'spring-boot-app'
    metrics_path: '/actuator/prometheus'
    static_configs:
      - targets: ['app:8080']
    scrape_interval: 15s
```

### Docker Compose for full monitoring stack:
```yaml
version: '3.8'
services:
  app:
    image: my-spring-app
    ports:
      - "8080:8080"

  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

### Key Grafana Dashboards for Spring Boot:
- Dashboard ID **4701** — JVM (Micrometer)
- Dashboard ID **12900** — Spring Boot 3.x Statistics
- Dashboard ID **11378** — Spring Boot Hikari Connection Pool

---

## Distributed Tracing with Micrometer Tracing

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
      probability: 1.0      # sample 100% (use 0.1 in prod)
  zipkin:
    tracing:
      endpoint: http://zipkin:9411/api/v2/spans
```

Trace IDs are automatically propagated across HTTP calls. Logs include `[traceId, spanId]`.

---

## Changing Log Levels at Runtime

```bash
# View current level
curl http://localhost:8080/actuator/loggers/com.example

# Change to DEBUG at runtime (no restart!)
curl -X POST http://localhost:8080/actuator/loggers/com.example \
  -H "Content-Type: application/json" \
  -d '{"configuredLevel": "DEBUG"}'
```

---

## Graceful Shutdown

```yaml
server:
  shutdown: graceful        # wait for active requests to complete

spring:
  lifecycle:
    timeout-per-shutdown-phase: 30s   # max wait time

management:
  endpoint:
    shutdown:
      enabled: true
```

```bash
# Trigger graceful shutdown
curl -X POST http://localhost:8080/actuator/shutdown
```

---

## Security for Actuator

```java
@Bean
public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
    http.authorizeHttpRequests(auth -> auth
        .requestMatchers("/actuator/health", "/actuator/info").permitAll()
        .requestMatchers("/actuator/**").hasRole("ADMIN")
        .anyRequest().authenticated()
    );
    return http.build();
}
```

Best practice in production: run actuator on a different port (not exposed to internet):
```yaml
management:
  server:
    port: 8081
```

Then firewall port 8081 from external access.
