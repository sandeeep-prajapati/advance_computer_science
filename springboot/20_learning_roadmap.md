# Chapter 20: Learning Roadmap

## 6-Week Hands-On Plan: Zero to Production-Ready

This roadmap takes you from installing Java to deploying a production-grade Spring Boot microservice. Each week has a project milestone.

---

## Week 1: Foundation

**Goal:** Build and run your first Spring Boot REST API.

### Topics:
- Install JDK 21 and set up IntelliJ IDEA
- Use Spring Initializr to create a project
- Understand `@SpringBootApplication`, application context, DI
- Build a CRUD REST API for a simple resource (e.g., Books)
- Use H2 in-memory database with Spring Data JPA
- Run and test with curl or Postman

### Milestone Project: **Library API**
```
POST /api/books         { title, author, isbn, year }
GET  /api/books         (all books, paginated)
GET  /api/books/{id}
PUT  /api/books/{id}
DELETE /api/books/{id}
GET  /api/books/search?author=...
```

### Exercises:
- [ ] Create a `Book` entity with `@Entity`, `@Id`, `@Column`
- [ ] Create `BookRepository extends JpaRepository<Book, Long>`
- [ ] Create `BookService` with business logic
- [ ] Create `BookController` with all HTTP methods
- [ ] Add `@Valid` validation to create/update requests
- [ ] Handle `ResourceNotFoundException` with `@ControllerAdvice`
- [ ] Write unit tests for `BookService` with Mockito
- [ ] Write `@WebMvcTest` for `BookController`

---

## Week 2: Data, Config & Security

**Goal:** Add a real database, profiles, and JWT authentication.

### Topics:
- Connect to PostgreSQL (use Docker Compose)
- Database migrations with Flyway
- Configuration profiles (dev/prod)
- `@ConfigurationProperties` for typed config
- Spring Security + JWT authentication

### Milestone Project: **Secure Library API**
```
POST /api/auth/register
POST /api/auth/login      → returns JWT token
GET  /api/books           → requires authentication
POST /api/books           → requires LIBRARIAN role
DELETE /api/books/{id}    → requires ADMIN role
```

### Exercises:
- [ ] Add PostgreSQL + Flyway migrations
- [ ] Create `application-dev.yml` and `application-prod.yml`
- [ ] Create `User` entity with email and password
- [ ] Implement `UserDetailsService`
- [ ] Implement JWT generation and validation
- [ ] Add `JwtAuthFilter`
- [ ] Configure `SecurityFilterChain`
- [ ] Write integration tests with `@SpringBootTest` + Testcontainers

---

## Week 3: Production Readiness

**Goal:** Make the app observable, containerized, and deployable.

### Topics:
- Spring Boot Actuator — health, metrics, info
- Micrometer + Prometheus metrics
- Docker with multi-stage build
- Docker Compose for full local stack
- Graceful shutdown

### Milestone Project: **Library API in Docker**

```bash
# Full stack: app + postgres + redis + prometheus + grafana
docker-compose up -d

# API serves on :8080
# Prometheus on :9090
# Grafana on :3000 (JVM + HTTP dashboards)
```

### Exercises:
- [ ] Add Actuator with health, info, metrics endpoints
- [ ] Add custom `HealthIndicator` for external dependency
- [ ] Add Micrometer Prometheus registry
- [ ] Write a multi-stage Dockerfile
- [ ] Create `docker-compose.yml` with all services
- [ ] Add Redis for session caching
- [ ] Configure Prometheus scraping
- [ ] Import JVM Grafana dashboard (ID: 4701)
- [ ] Configure graceful shutdown

---

## Week 4: Caching, Messaging & Advanced Features

**Goal:** Add caching, async processing, and event-driven patterns.

### Topics:
- Spring Cache with Redis
- `@Async` and thread pools
- `@Scheduled` jobs
- Apache Kafka basics
- Exception handling patterns

### Milestone Project: **Library API with Notifications**

```
Book loan event → Kafka → Notification Service → Email user
                                               → Update stats (async)
Cache expensive catalog stats with Redis
Daily report generated via @Scheduled job
```

### Exercises:
- [ ] Add Redis cache for `findById` with `@Cacheable`
- [ ] Add `@CacheEvict` on update/delete
- [ ] Create a `BookLoanService` that publishes Kafka events
- [ ] Create a separate `NotificationService` (can be a different module)
- [ ] Add `@KafkaListener` for loan events
- [ ] Add `@Async` email sending
- [ ] Add `@Scheduled` daily overdue check job
- [ ] Configure proper thread pools for async work

---

## Week 5: Microservices & Spring Cloud

**Goal:** Split the monolith into microservices with service discovery.

### Topics:
- Microservices architecture patterns
- Eureka service discovery
- OpenFeign declarative clients
- Spring Cloud Gateway
- Resilience4j circuit breaker
- Distributed tracing with Zipkin

### Milestone Project: **Library Platform (3 services)**

```
API Gateway (:8080)
├── User Service (:8081)    → users DB
├── Book Service (:8082)    → books DB
└── Loan Service (:8083)    → loans DB

Eureka Server (:8761)
Zipkin (:9411)
```

### Exercises:
- [ ] Create Eureka Server project
- [ ] Split Library API into 3 separate services
- [ ] Register all services with Eureka
- [ ] Add Feign client in Loan Service to call User Service
- [ ] Add Spring Cloud Gateway with routing config
- [ ] Add authentication filter in Gateway
- [ ] Add Resilience4j circuit breaker to Feign client
- [ ] Configure Zipkin distributed tracing
- [ ] Verify trace propagates across services in Zipkin UI

---

## Week 6: Kubernetes Deployment

**Goal:** Deploy all services to Kubernetes.

### Topics:
- Kubernetes Deployments, Services, Ingress
- ConfigMaps and Secrets
- Health probes with Actuator
- Helm charts
- Horizontal Pod Autoscaler

### Milestone Project: **Library Platform on K8s**

```bash
# Local K8s (kind or minikube)
kubectl apply -f k8s/

# All services running in K8s
# Ingress routes traffic to correct service
# HPA scales services under load
# Secrets managed securely
```

### Exercises:
- [ ] Write Deployment YAML for each service
- [ ] Configure liveness/readiness probes
- [ ] Create K8s Secrets for DB passwords and JWT secret
- [ ] Create ConfigMaps for non-sensitive config
- [ ] Create Services and Ingress
- [ ] Package as Helm chart with `values.yaml`
- [ ] Configure HPA for Book Service
- [ ] Deploy with `helm install`
- [ ] Test rolling update: `helm upgrade --set image.tag=v2`

---

## Learning Path Visualization

```
Week 1: [Basics] → REST API, JPA, Unit Tests
    ↓
Week 2: [Security] → JWT, PostgreSQL, Profiles, Flyway
    ↓
Week 3: [Production] → Actuator, Docker, Prometheus, Grafana
    ↓
Week 4: [Async] → Redis Cache, Kafka, @Async, @Scheduled
    ↓
Week 5: [Microservices] → Eureka, Feign, Gateway, Resilience4j
    ↓
Week 6: [Kubernetes] → K8s, Helm, HPA, Probes
    ↓
[Production-Ready Spring Boot Microservice Engineer]
```

---

## Key Concepts Cheat Sheet

| Concept | Annotation/Class | Chapter |
|---|---|---|
| REST Controller | `@RestController`, `@GetMapping` | 4 |
| JPA Entity | `@Entity`, `@Id`, `@Column` | 5 |
| Repository | `JpaRepository`, `@Query` | 5 |
| Transactions | `@Transactional` | 5 |
| Security | `SecurityFilterChain`, `@PreAuthorize` | 6 |
| Configuration | `@ConfigurationProperties`, `@Value` | 7 |
| Testing | `@WebMvcTest`, `@DataJpaTest`, Testcontainers | 8 |
| Health check | `HealthIndicator`, `/actuator/health` | 9 |
| Error handling | `@ControllerAdvice`, `@ExceptionHandler` | 10 |
| Validation | `@Valid`, `@NotBlank`, `@Email` | 10 |
| Docker | Dockerfile, `docker-compose.yml` | 11 |
| Kubernetes | Deployment, Service, Ingress, Helm | 12 |
| Microservices | Service decomposition, Saga pattern | 13 |
| Service discovery | Eureka, `@LoadBalanced` | 14 |
| Circuit breaker | `@CircuitBreaker`, Resilience4j | 14 |
| Caching | `@Cacheable`, `@CacheEvict`, Redis | 15 |
| Kafka | `KafkaTemplate`, `@KafkaListener` | 16 |
| RabbitMQ | `RabbitTemplate`, `@RabbitListener` | 16 |
| Async | `@Async`, `CompletableFuture` | 17 |
| Scheduling | `@Scheduled`, cron expressions | 17 |
| GraphQL | `@QueryMapping`, `@MutationMapping` | 18 |
| Native image | GraalVM, `native:compile` | 19 |

---

## Recommended Resources

- **Official Docs:** docs.spring.io — always start here
- **Spring Guides:** spring.io/guides — step-by-step tutorials
- **Baeldung:** baeldung.com — practical Spring Boot examples
- **Spring Initializr:** start.spring.io — project generator
- **Books:** "Spring Boot in Action" (Craig Walls), "Cloud Native Spring in Action" (Thomas Vitale)
- **GitHub:** github.com/spring-projects/spring-boot/tree/main/spring-boot-tests — official test examples

---

## Common Mistakes to Avoid

1. **Field injection** (`@Autowired` on field) — use constructor injection
2. **Exposing @Entity in API** — always use DTOs
3. **Missing `@Transactional`** — service layer mutations need it
4. **N+1 queries** — use JOIN FETCH or @BatchMapping
5. **`ddl-auto=update` in production** — use Flyway
6. **Secrets in application.properties** — use environment variables
7. **Single thread for `@Scheduled`** — configure a thread pool
8. **`@Async` self-call** — won't work without proxy
9. **Infinite recursion in JSON** — use `@JsonIgnore` or DTOs for bidirectional relationships
10. **No timeout on external calls** — always set connection and read timeouts
