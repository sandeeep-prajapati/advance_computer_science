# Chapter 3: Core Concepts — IoC, DI, Beans, and Application Context

## Inversion of Control (IoC)

**IoC** is a design principle where the framework controls the flow of the program, not your code. Instead of you creating objects, the Spring container creates and manages them.

```
Traditional:                      With IoC:
  YourCode → creates → Objects      Spring Container → creates → Objects
  YourCode → wires   → Dependencies Spring Container → injects → Dependencies
  YourCode → manages → Lifecycle    Spring Container → manages → Lifecycle
```

---

## Dependency Injection (DI)

DI is the mechanism Spring uses to implement IoC. Instead of a class creating its own dependencies, they are **injected** from outside.

### Without DI (tightly coupled):
```java
public class OrderService {
    private PaymentService paymentService = new PaymentService(); // bad: hard-coded dependency

    public void placeOrder() {
        paymentService.charge();
    }
}
```

### With Constructor Injection (recommended):
```java
@Service
public class OrderService {
    private final PaymentService paymentService;

    public OrderService(PaymentService paymentService) { // Spring injects this
        this.paymentService = paymentService;
    }

    public void placeOrder() {
        paymentService.charge();
    }
}
```

### Three Types of Injection

| Type | Example | Recommendation |
|---|---|---|
| **Constructor injection** | `public Service(Dep dep)` | **Preferred** — enforces required deps, immutable |
| **Setter injection** | `@Autowired public void setDep(Dep dep)` | For optional dependencies |
| **Field injection** | `@Autowired private Dep dep;` | Avoid — hard to test, hides dependencies |

---

## Spring Beans

A **Bean** is an object that Spring creates, configures, and manages. Every object managed by the Spring IoC container is a bean.

### Declaring Beans

**Method 1: Stereotype Annotations (auto-detected via component scan):**

```java
@Component          // Generic bean
@Service            // Business logic layer
@Repository         // Data access layer (also translates DB exceptions)
@Controller         // Web layer (returns views)
@RestController     // Web layer (returns data — @Controller + @ResponseBody)
```

```java
@Service
public class UserService {
    // Spring creates and manages this bean
}
```

**Method 2: @Bean in @Configuration class (explicit):**

```java
@Configuration
public class AppConfig {

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    @Bean
    public RestTemplate restTemplate() {
        return new RestTemplate();
    }
}
```

Use `@Bean` when:
- You don't own the class (third-party library)
- You need fine-grained control over initialization

---

## Application Context

The **ApplicationContext** is Spring's IoC container. It holds all beans and their relationships.

```java
// Programmatic access (rarely needed in Spring Boot apps)
ApplicationContext ctx = SpringApplication.run(MyApp.class, args);
UserService userService = ctx.getBean(UserService.class);
```

In Spring Boot, the ApplicationContext is created automatically when the app starts. You rarely interact with it directly — just use `@Autowired` or constructor injection.

### Context Lifecycle:

```
Application starts
    ↓
Spring scans @Component, @Service, @Repository, @Controller
    ↓
Spring creates beans in dependency order
    ↓
Spring injects dependencies
    ↓
@PostConstruct methods run
    ↓
Application is ready → serves requests
    ↓
Application shuts down
    ↓
@PreDestroy methods run
    ↓
Beans destroyed
```

---

## Component Scan

Spring Boot automatically scans the package containing `@SpringBootApplication` and all its sub-packages.

```
com.example.myapp
├── MyAppApplication.java     ← @SpringBootApplication here
├── controller/
│   └── UserController.java   ← auto-detected
├── service/
│   └── UserService.java      ← auto-detected
└── repository/
    └── UserRepository.java   ← auto-detected
```

If you put a class outside this package, you must add it explicitly:
```java
@SpringBootApplication(scanBasePackages = {"com.example.myapp", "com.other.package"})
```

---

## Bean Scopes

| Scope | Behavior | When to Use |
|---|---|---|
| `singleton` (default) | One instance per container | Stateless services (most beans) |
| `prototype` | New instance every injection | Stateful objects |
| `request` | One per HTTP request | Web-layer beans |
| `session` | One per HTTP session | User session data |

```java
@Component
@Scope("prototype")
public class ReportGenerator {
    // new instance created every time it's injected
}
```

---

## @PostConstruct and @PreDestroy

```java
@Service
public class CacheService {

    @PostConstruct
    public void init() {
        // Runs after bean is created and dependencies are injected
        loadCacheFromDatabase();
    }

    @PreDestroy
    public void cleanup() {
        // Runs before bean is destroyed (app shutdown)
        flushCacheToDisk();
    }
}
```

---

## Qualifiers — Resolving Multiple Beans

When two beans of the same type exist, Spring doesn't know which to inject:

```java
@Service
public class EmailNotifier implements Notifier { ... }

@Service
public class SmsNotifier implements Notifier { ... }
```

**Solution 1: @Qualifier**
```java
@Service
public class OrderService {
    public OrderService(@Qualifier("emailNotifier") Notifier notifier) {
        this.notifier = notifier;
    }
}
```

**Solution 2: @Primary**
```java
@Service
@Primary   // this one is injected by default
public class EmailNotifier implements Notifier { ... }
```

---

## @Conditional Beans

Create beans only under certain conditions:

```java
@Configuration
public class DataSourceConfig {

    @Bean
    @Profile("dev")
    public DataSource h2DataSource() {
        return new EmbeddedDatabaseBuilder().setType(H2).build();
    }

    @Bean
    @Profile("prod")
    public DataSource postgresDataSource() {
        return DataSourceBuilder.create().url("jdbc:postgresql://...").build();
    }
}
```

Other conditions:
```java
@ConditionalOnProperty(name = "feature.cache.enabled", havingValue = "true")
@ConditionalOnClass(name = "redis.clients.jedis.Jedis")
@ConditionalOnMissingBean(CacheManager.class)
```

---

## Aspect-Oriented Programming (AOP)

AOP lets you add cross-cutting concerns (logging, transactions, security) without modifying business code.

```java
@Aspect
@Component
public class LoggingAspect {

    @Around("execution(* com.example.myapp.service.*.*(..))")
    public Object logExecutionTime(ProceedingJoinPoint joinPoint) throws Throwable {
        long start = System.currentTimeMillis();
        Object result = joinPoint.proceed();
        long duration = System.currentTimeMillis() - start;
        log.info("{} executed in {}ms", joinPoint.getSignature(), duration);
        return result;
    }
}
```

Spring uses AOP internally for:
- `@Transactional` — wraps methods in database transactions
- `@Cacheable` — intercepts methods to return cached results
- `@PreAuthorize` — checks security before method execution

---

## Summary: Spring Core Mental Model

```
@SpringBootApplication
    ↓ triggers
Component Scan + Auto-Configuration
    ↓ creates
ApplicationContext (IoC Container)
    ↓ manages
Beans (@Component / @Service / @Repository / @Controller / @Bean)
    ↓ wired via
Dependency Injection (Constructor / Setter / Field)
    ↓ enhanced by
AOP (@Transactional, @Cacheable, @PreAuthorize)
```
