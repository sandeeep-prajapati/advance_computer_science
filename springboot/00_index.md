# Spring Boot Developer Book

A comprehensive, chapter-wise guide to mastering Spring Boot — from first application to production-grade microservices, security, monitoring, and cloud deployment.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_introduction.md) | **Introduction to Spring Boot** | What is Spring Boot, why it exists, Spring vs Spring Boot, key features, auto-configuration |
| [02](02_setup_and_installation.md) | **Setup & Installation** | Spring Initializr, IDE setup, Maven vs Gradle, project structure, first REST endpoint |
| [03](03_core_concepts.md) | **Core Concepts** | IoC container, Dependency Injection, Beans, Application Context, Component scan |
| [04](04_rest_apis.md) | **REST APIs** | @RestController, @RequestMapping, HTTP methods, request/response, DTOs, validation |
| [05](05_spring_data_jpa.md) | **Spring Data JPA & Databases** | JPA entities, repositories, CRUD, relationships, JPQL, pagination, transactions |
| [06](06_spring_security.md) | **Spring Security** | Authentication, authorization, JWT, OAuth2, form login, method-level security |
| [07](07_configuration_and_profiles.md) | **Configuration & Profiles** | application.properties/yaml, profiles, externalized config, @ConfigurationProperties |
| [08](08_testing.md) | **Testing Spring Boot** | Unit tests, @SpringBootTest, MockMvc, @DataJpaTest, Testcontainers, test slices |
| [09](09_actuator_and_monitoring.md) | **Actuator & Monitoring** | Health checks, metrics, Prometheus, Grafana, custom endpoints, info endpoint |
| [10](10_exception_handling_and_validation.md) | **Exception Handling & Validation** | @ExceptionHandler, @ControllerAdvice, Bean Validation, custom validators, error responses |
| [11](11_spring_boot_with_docker.md) | **Spring Boot with Docker** | Dockerfile, multi-stage builds, docker-compose, buildpacks, environment variables |
| [12](12_spring_boot_with_kubernetes.md) | **Spring Boot with Kubernetes** | Deploying to K8s, ConfigMaps, Secrets, health probes, Helm charts, scaling |
| [13](13_microservices.md) | **Spring Boot Microservices** | Microservices architecture, service decomposition, inter-service communication, API Gateway |
| [14](14_spring_cloud.md) | **Spring Cloud** | Eureka, Config Server, Feign, Resilience4j, Circuit Breaker, distributed tracing |
| [15](15_caching.md) | **Caching** | @Cacheable, @CacheEvict, Redis cache, Caffeine, cache configuration, TTL |
| [16](16_messaging.md) | **Messaging — Kafka & RabbitMQ** | Spring Kafka, @KafkaListener, RabbitMQ, AMQP, producers/consumers, message routing |
| [17](17_scheduling_and_async.md) | **Scheduling & Async** | @Scheduled, @Async, thread pools, @EnableAsync, virtual threads (Java 21) |
| [18](18_graphql.md) | **Spring Boot with GraphQL** | GraphQL schema, @QueryMapping, @MutationMapping, subscriptions, N+1 problem |
| [19](19_performance_and_production.md) | **Performance & Production** | Connection pooling, lazy loading, query optimization, GraalVM native image, health tuning |
| [20](20_learning_roadmap.md) | **Learning Roadmap** | 6-week hands-on plan from zero to production-ready microservices |

---

## How to Read This Book

**If you're brand new to Spring Boot:** Read 1 → 2 → 3 → 4 → 5 in order. This gives you the core mental model.

**If you're building a REST API:** Read 2 → 3 → 4 → 10 → 7.

**If you're securing an application:** Read 6 → 7 → 10.

**If you're going to production:** Read 8 → 9 → 11 → 12 → 19.

**If you're building microservices:** Read 13 → 14 → 16 → 15.

**If you need to deploy on Kubernetes:** Read 11 → 12.

---

## Key Concepts at a Glance

```
Spring Boot Application
├── Auto-Configuration
│   ├── @SpringBootApplication      (entry point)
│   ├── @EnableAutoConfiguration    (magic starter)
│   └── spring.factories / imports  (auto-config registry)
│
├── IoC Container
│   ├── @Component / @Service / @Repository / @Controller
│   ├── @Bean (manual bean definition)
│   └── @Autowired / Constructor Injection
│
├── Web Layer
│   ├── @RestController             (REST endpoints)
│   ├── @RequestMapping / @GetMapping etc.
│   ├── @RequestBody / @PathVariable / @RequestParam
│   └── ResponseEntity<T>
│
├── Data Layer
│   ├── @Entity / @Table            (JPA mapping)
│   ├── JpaRepository<T, ID>        (CRUD + paging)
│   ├── @Transactional              (transaction boundary)
│   └── Spring Data JDBC / R2DBC    (reactive / simple)
│
├── Security
│   ├── SecurityFilterChain         (request filter chain)
│   ├── UserDetailsService          (custom user loading)
│   ├── JWT / OAuth2                (token-based auth)
│   └── @PreAuthorize               (method security)
│
├── Config & Profiles
│   ├── application.yml             (base config)
│   ├── application-{profile}.yml   (profile-specific)
│   └── @ConfigurationProperties    (typed config binding)
│
└── Production Tooling
    ├── Actuator                    (health, metrics, info)
    ├── Micrometer + Prometheus     (metrics export)
    └── Testcontainers              (real DB in tests)
```
