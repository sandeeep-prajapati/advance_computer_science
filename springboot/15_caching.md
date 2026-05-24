# Chapter 15: Caching

## Why Caching?

Caching stores the results of expensive operations (database queries, API calls, computations) in fast storage so future requests get the cached result instead of repeating the expensive work.

```
Without cache:   Request → Service → Database (50ms) → Response
With cache hit:  Request → Service → Cache (1ms) → Response
With cache miss: Request → Service → Cache miss → Database (50ms) → Cache.store → Response
```

---

## Spring Cache Abstraction

Spring provides a cache abstraction that works with multiple backends: in-memory (Caffeine), Redis, Hazelcast, JCache, etc.

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-cache</artifactId>
</dependency>
```

```java
@SpringBootApplication
@EnableCaching    // activate the cache infrastructure
public class MyApplication { ... }
```

---

## Core Cache Annotations

```java
@Service
public class ProductService {

    // Cache the result. Key defaults to method args.
    @Cacheable(value = "products", key = "#id")
    public ProductDto getProduct(Long id) {
        log.info("Loading product {} from database", id);   // only logs on cache miss
        return productRepository.findById(id)
            .map(this::toDto)
            .orElseThrow(() -> new ResourceNotFoundException("Product", id));
    }

    // Evict a specific entry when a product is updated
    @CacheEvict(value = "products", key = "#id")
    public ProductDto updateProduct(Long id, UpdateProductRequest request) {
        Product product = productRepository.findById(id).orElseThrow();
        product.setName(request.name());
        product.setPrice(request.price());
        return toDto(productRepository.save(product));
    }

    // Evict ALL entries in the cache when bulk-updating
    @CacheEvict(value = "products", allEntries = true)
    public void clearProductCache() { }

    // Update cache entry after method completes (always executes the method)
    @CachePut(value = "products", key = "#result.id")
    public ProductDto createProduct(CreateProductRequest request) {
        Product product = productRepository.save(new Product(request));
        return toDto(product);
    }

    // Multiple cache operations
    @Caching(
        evict = {
            @CacheEvict(value = "products", key = "#id"),
            @CacheEvict(value = "productsByCategory", allEntries = true)
        }
    )
    public void deleteProduct(Long id) {
        productRepository.deleteById(id);
    }
}
```

---

## Cache Keys

```java
// Default: all method parameters as key
@Cacheable("users")
public UserDto getUser(Long id) { ... }   // key = 1, 2, 3...

// Custom SpEL expression
@Cacheable(value = "users", key = "#email.toLowerCase()")
public UserDto getUserByEmail(String email) { ... }

// Complex key
@Cacheable(value = "search", key = "T(java.util.Objects).hash(#query, #page, #size)")
public Page<ProductDto> search(String query, int page, int size) { ... }

// Conditional caching
@Cacheable(value = "products", condition = "#id > 0", unless = "#result == null")
public ProductDto getProduct(Long id) { ... }
```

---

## Option A: Simple In-Memory Cache (ConcurrentHashMap)

Works out of the box, no extra dependencies. For single-instance apps only.

```yaml
spring:
  cache:
    type: simple
```

Not suitable for production with multiple instances — each instance has its own cache.

---

## Option B: Caffeine (Local Cache with TTL)

High-performance, feature-rich in-memory cache. Best for single-instance or when Redis is overkill.

```xml
<dependency>
    <groupId>com.github.ben-manes.caffeine</groupId>
    <artifactId>caffeine</artifactId>
</dependency>
```

```yaml
spring:
  cache:
    type: caffeine
    caffeine:
      spec: maximumSize=1000,expireAfterWrite=5m
```

```java
@Configuration
@EnableCaching
public class CacheConfig {

    @Bean
    public CacheManager cacheManager() {
        CaffeineCacheManager manager = new CaffeineCacheManager();
        manager.setCaffeine(Caffeine.newBuilder()
            .maximumSize(1_000)
            .expireAfterWrite(Duration.ofMinutes(10))
            .recordStats());   // enable cache statistics
        return manager;
    }

    // Different TTL per cache
    @Bean
    public CacheManager cacheManager() {
        SimpleCacheManager manager = new SimpleCacheManager();
        manager.setCaches(List.of(
            buildCache("products", 1000, Duration.ofMinutes(10)),
            buildCache("users", 500, Duration.ofMinutes(30)),
            buildCache("config", 100, Duration.ofHours(1))
        ));
        return manager;
    }

    private CaffeineCache buildCache(String name, int maxSize, Duration ttl) {
        return new CaffeineCache(name, Caffeine.newBuilder()
            .maximumSize(maxSize)
            .expireAfterWrite(ttl)
            .build());
    }
}
```

---

## Option C: Redis (Distributed Cache)

For multi-instance deployments — all instances share the same cache.

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-data-redis</artifactId>
</dependency>
```

```yaml
spring:
  cache:
    type: redis
    redis:
      time-to-live: 10m        # default TTL for all caches
      key-prefix: myapp::      # prefix all keys
      use-key-prefix: true
  data:
    redis:
      host: localhost
      port: 6379
      password: ${REDIS_PASSWORD}
      timeout: 2000ms
      lettuce:
        pool:
          max-active: 8
          max-idle: 8
          min-idle: 1
```

```java
@Configuration
@EnableCaching
public class RedisConfig {

    @Bean
    public RedisCacheManager cacheManager(RedisConnectionFactory factory) {
        RedisCacheConfiguration defaultConfig = RedisCacheConfiguration.defaultCacheConfig()
            .entryTtl(Duration.ofMinutes(10))
            .serializeValuesWith(RedisSerializationContext.SerializationPair.fromSerializer(
                new GenericJackson2JsonRedisSerializer()
            ))
            .disableCachingNullValues();

        Map<String, RedisCacheConfiguration> cacheConfigs = Map.of(
            "products", defaultConfig.entryTtl(Duration.ofMinutes(30)),
            "users", defaultConfig.entryTtl(Duration.ofHours(1)),
            "sessions", defaultConfig.entryTtl(Duration.ofDays(7))
        );

        return RedisCacheManager.builder(factory)
            .cacheDefaults(defaultConfig)
            .withInitialCacheConfigurations(cacheConfigs)
            .build();
    }
}
```

---

## Redis as a Data Store (Beyond Cache)

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-data-redis</artifactId>
</dependency>
```

```java
@Service
public class SessionService {

    private final StringRedisTemplate redisTemplate;

    // Store session
    public void createSession(String sessionId, UserSession session) {
        String key = "session:" + sessionId;
        redisTemplate.opsForValue().set(key, serialize(session), Duration.ofHours(24));
    }

    // Get session
    public Optional<UserSession> getSession(String sessionId) {
        String value = redisTemplate.opsForValue().get("session:" + sessionId);
        return Optional.ofNullable(value).map(this::deserialize);
    }

    // Extend TTL
    public void extendSession(String sessionId) {
        redisTemplate.expire("session:" + sessionId, Duration.ofHours(24));
    }

    // Delete session
    public void deleteSession(String sessionId) {
        redisTemplate.delete("session:" + sessionId);
    }
}
```

---

## Cache-Aside vs Write-Through vs Write-Behind

| Pattern | How | Use When |
|---|---|---|
| **Cache-aside** (default) | App reads cache; on miss, loads from DB and stores in cache | Most read-heavy workloads |
| **Write-through** | Write to cache AND DB together | Strong consistency needed |
| **Write-behind** | Write to cache, async write to DB | High write throughput |
| **Read-through** | Cache fetches from DB on miss automatically | Simplicity |

Spring's `@Cacheable` implements cache-aside.

---

## Cache Eviction Strategies

| Strategy | Description |
|---|---|
| `expireAfterWrite` | TTL from time of write |
| `expireAfterAccess` | TTL from time of last access |
| `maximumSize` | LRU eviction when limit reached |
| Manual eviction | `@CacheEvict` or `cache.evict(key)` |

---

## Cache Best Practices

- Cache **read-heavy, write-rarely** data (product catalog, user profiles)
- **Don't cache** frequently changing data or user-specific sensitive data
- Set appropriate **TTL** — shorter for volatile data, longer for stable data
- **Monitor cache hit rate** — Caffeine.recordStats() or Redis INFO stats
- Use **cache warming** on startup for critical data
- Plan for **cache stampede** — use probabilistic early expiration or locking

```java
// Cache statistics with Caffeine
@Bean
public CacheManager cacheManager(MeterRegistry registry) {
    CaffeineCacheManager manager = new CaffeineCacheManager();
    manager.setCaffeine(Caffeine.newBuilder()
        .recordStats());
    // Micrometer will export cache stats automatically
    return manager;
}
```
