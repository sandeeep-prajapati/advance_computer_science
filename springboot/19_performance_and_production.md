# Chapter 19: Performance & Production Best Practices

## Performance Mindset

Production performance problems fall into a few categories:
1. **Database** — N+1 queries, missing indexes, full table scans
2. **Memory** — memory leaks, excessive GC, too-large caches
3. **Threading** — thread pool exhaustion, blocking I/O
4. **Network** — chatty microservices, large payloads
5. **Cold start** — slow JVM startup (relevant for serverless/K8s scale-to-zero)

---

## Database Performance

### 1. Connection Pool Tuning (HikariCP)

HikariCP is Spring Boot's default connection pool.

```yaml
spring:
  datasource:
    hikari:
      maximum-pool-size: 10      # don't blindly set high — DB can't handle unlimited
      minimum-idle: 2
      idle-timeout: 300000       # remove idle connections after 5 min
      connection-timeout: 30000  # wait max 30s for a connection
      max-lifetime: 1800000      # recycle connections every 30 min
      validation-timeout: 5000
      leak-detection-threshold: 60000  # warn if connection held > 60s
```

**Rule of thumb:** `pool-size = (number of CPU cores * 2) + effective_spindle_count`
For most apps: 10 connections is plenty. Don't go above 50.

### 2. Query Optimization

```java
// Bad: N+1 problem
List<User> users = userRepository.findAll();
users.forEach(user -> {
    List<Order> orders = orderRepository.findByUserId(user.getId());  // N queries!
});

// Good: JOIN FETCH
@Query("SELECT u FROM User u LEFT JOIN FETCH u.orders WHERE u.status = :status")
List<User> findUsersWithOrders(UserStatus status);

// Good: DTO projection (faster than loading entities + mapping)
@Query("SELECT new com.example.UserSummary(u.id, u.name, COUNT(o.id)) " +
       "FROM User u LEFT JOIN u.orders o GROUP BY u.id, u.name")
List<UserSummary> findUserSummaries();
```

### 3. Indexes

Add indexes for columns used in `WHERE`, `JOIN`, `ORDER BY`:

```sql
-- In Flyway migration
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status_created ON orders(status, created_at DESC);

-- Composite index for pagination queries
CREATE INDEX idx_products_category_price ON products(category_id, price);
```

### 4. Pagination — Always Paginate Large Queries

```java
// Never load all records
List<User> allUsers = userRepository.findAll();  // could be millions!

// Always use pagination
Page<User> users = userRepository.findAll(PageRequest.of(0, 100, Sort.by("createdAt").descending()));
```

### 5. Read-Only Transactions

```java
@Transactional(readOnly = true)   // skips dirty checking, faster flush
public List<UserDto> getUsers() {
    return userRepository.findAll().stream().map(this::toDto).toList();
}
```

---

## Caching Strategy

```java
// Cache expensive aggregation queries
@Cacheable(value = "dashboardStats", key = "#orgId")
@Transactional(readOnly = true)
public DashboardStats getDashboardStats(Long orgId) {
    // Heavy query with multiple joins and aggregations
    return statsRepository.computeStats(orgId);
}

// Evict on data change
@CacheEvict(value = "dashboardStats", key = "#order.orgId")
public Order placeOrder(OrderRequest request, Long orgId) { ... }
```

---

## JVM Tuning

```bash
# Production JVM flags (tune for your heap size)
java \
  -Xms512m \                        # initial heap
  -Xmx2g \                          # max heap
  -XX:MaxMetaspaceSize=256m \       # metaspace limit
  -XX:+UseG1GC \                    # G1 GC (good all-rounder, default JDK 9+)
  -XX:MaxGCPauseMillis=200 \        # target max GC pause
  -XX:+HeapDumpOnOutOfMemoryError \ # dump heap on OOM for analysis
  -XX:HeapDumpPath=/tmp/heapdump.hprof \
  -XX:+UseContainerSupport \        # respect container memory limits
  -XX:MaxRAMPercentage=75.0 \       # use 75% of container memory for heap
  -jar app.jar
```

### GC Options:
| GC | Flag | Best For |
|---|---|---|
| G1GC | `-XX:+UseG1GC` | Balanced latency/throughput (default Java 9+) |
| ZGC | `-XX:+UseZGC` | Ultra-low pause (<1ms), large heaps |
| Shenandoah | `-XX:+UseShenandoahGC` | Low latency (Red Hat JDK) |
| SerialGC | `-XX:+UseSerialGC` | Small footprint, single-core |

---

## GraalVM Native Image

Compile Spring Boot app to a native binary — millisecond startup, much lower memory.

```xml
<plugin>
    <groupId>org.graalvm.buildtools</groupId>
    <artifactId>native-maven-plugin</artifactId>
</plugin>
```

```bash
# Build native image (requires GraalVM JDK)
./mvnw native:compile -Pnative

# Run native binary
./target/my-app
# Starts in ~100ms vs 5-10s for JVM!
```

**Trade-offs:**
- Build time: 3-10 minutes (vs seconds for JVM)
- No dynamic class loading, reflection requires configuration hints
- Profiling/debugging harder
- Best for: serverless, sidecar containers, CLI tools

---

## HTTP Performance

### Enable HTTP/2:
```yaml
server:
  http2:
    enabled: true
```

### Enable Compression:
```yaml
server:
  compression:
    enabled: true
    mime-types: application/json,application/xml,text/plain,text/html
    min-response-size: 1024   # only compress if response > 1KB
```

### Async Request Processing:
```java
// Return Callable — releases Tomcat thread while processing
@GetMapping("/reports/{id}")
public Callable<ReportDto> getReport(@PathVariable Long id) {
    return () -> reportService.generate(id);   // executed in separate thread
}
```

---

## Lazy Loading and Initialization

```yaml
spring:
  main:
    lazy-initialization: true   # beans initialized on first use, not at startup
                                # reduces startup time; test thoroughly
```

```java
@Bean
@Lazy   // this specific bean is lazy
public SlowToInitializeService slowService() { ... }
```

### CDS (Class Data Sharing) — Spring Boot 3.3+:

```bash
# Training run — capture loaded classes
java -XX:ArchiveClassesAtExit=app.jsa -jar app.jar -Dspring.context.exit=onRefresh

# Production run — use cached class data
java -XX:SharedArchiveFile=app.jsa -jar app.jar
# 30-40% faster startup!
```

---

## Security Performance

```java
// Cache BCrypt result — BCrypt is intentionally slow
@Cacheable(value = "passwordChecks", key = "#rawPassword + #encodedPassword")
public boolean checkPassword(String rawPassword, String encodedPassword) {
    return passwordEncoder.matches(rawPassword, encodedPassword);
}
// Note: Only appropriate for rate-limited auth flows — don't cache blindly
```

---

## Observability and Profiling

### Enable JFR (Java Flight Recorder) — zero-overhead profiling:
```bash
java -XX:StartFlightRecording=duration=60s,filename=recording.jfr -jar app.jar
```

### Heap analysis:
```bash
# Trigger heap dump
curl -X POST http://localhost:8080/actuator/heapdump > heap.hprof
# Analyze with Eclipse MAT or IntelliJ
```

---

## Production Checklist

### Application:
- [ ] `server.shutdown=graceful` configured
- [ ] Health probes configured (liveness, readiness)
- [ ] Actuator on separate port, secured
- [ ] All secrets in environment variables or vault
- [ ] Structured logging (JSON) for log aggregation
- [ ] Correlation IDs / tracing enabled
- [ ] Connection pool properly sized
- [ ] Pagination on all list endpoints
- [ ] Input validation on all endpoints

### JVM:
- [ ] Heap size explicitly set (`-Xmx`)
- [ ] `+UseContainerSupport` if running in containers
- [ ] GC algorithm chosen for workload
- [ ] OOM heap dump enabled
- [ ] JFR or async profiler available

### Database:
- [ ] All query-time columns indexed
- [ ] N+1 queries eliminated
- [ ] `readOnly=true` on read transactions
- [ ] Connection pool sized appropriately
- [ ] Database migrations with Flyway/Liquibase
- [ ] Slow query logging enabled in DB

### Infrastructure:
- [ ] Resource requests and limits set in K8s
- [ ] HPA configured
- [ ] PodDisruptionBudget to prevent full outage during deploys
- [ ] Rolling update strategy (`maxUnavailable=0`)
- [ ] Circuit breakers for all external calls
- [ ] Retry with exponential backoff

---

## Quick Reference: Common Performance Issues

| Symptom | Likely Cause | Fix |
|---|---|---|
| Slow queries | Missing index, N+1 | Add index, use JOIN FETCH |
| High memory | Memory leak, large heap | Heap dump → MAT analysis |
| Request timeouts | Thread pool exhaustion | Increase pool, use async |
| Slow startup | Eager init, heavy beans | Lazy init, CDS |
| 503 under load | Too few pods, small pool | HPA, increase pool |
| DB connection errors | Pool exhausted | Increase pool, fix long transactions |
