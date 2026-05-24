# Chapter 5: Spring Data JPA & Databases

## What is JPA?

**JPA (Jakarta Persistence API)** is a specification for mapping Java objects to relational database tables (ORM). **Hibernate** is the most popular JPA implementation. **Spring Data JPA** adds a repository abstraction on top of JPA that eliminates boilerplate CRUD code.

```
Your Code → Spring Data JPA → JPA (Hibernate) → JDBC → Database
```

---

## Dependencies

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-data-jpa</artifactId>
</dependency>
<dependency>
    <groupId>org.postgresql</groupId>
    <artifactId>postgresql</artifactId>
    <scope>runtime</scope>
</dependency>
```

---

## Database Configuration

```yaml
# application.yml
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/mydb
    username: myuser
    password: mypassword
    driver-class-name: org.postgresql.Driver
    hikari:
      maximum-pool-size: 10
      minimum-idle: 2

  jpa:
    hibernate:
      ddl-auto: validate        # use 'create-drop' in dev, 'validate' in prod
    show-sql: true
    properties:
      hibernate:
        format_sql: true
        dialect: org.hibernate.dialect.PostgreSQLDialect
```

**ddl-auto values:**
| Value | Behavior |
|---|---|
| `none` | No action |
| `validate` | Validates schema matches entities (prod) |
| `update` | Alters tables to match entities (dev) |
| `create` | Drops and recreates on every start |
| `create-drop` | Creates on start, drops on stop (tests) |

---

## Defining Entities

```java
@Entity
@Table(name = "users")
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false, length = 100)
    private String name;

    @Column(unique = true, nullable = false)
    private String email;

    @Enumerated(EnumType.STRING)
    private UserStatus status;

    @CreationTimestamp
    @Column(updatable = false)
    private LocalDateTime createdAt;

    @UpdateTimestamp
    private LocalDateTime updatedAt;

    // getters, setters (or use Lombok)
}
```

### Use Lombok to reduce boilerplate:
```java
@Entity
@Table(name = "users")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String name;
    private String email;
}
```

---

## Spring Data Repositories

Spring Data JPA generates implementations automatically at runtime.

```java
// Basic CRUD
public interface UserRepository extends JpaRepository<User, Long> {
    // JpaRepository gives you: save, findById, findAll, delete, count, etc.
}
```

### Derived Query Methods

Spring Data parses method names and generates SQL:

```java
public interface UserRepository extends JpaRepository<User, Long> {

    Optional<User> findByEmail(String email);

    List<User> findByNameContainingIgnoreCase(String name);

    List<User> findByStatusAndCreatedAtAfter(UserStatus status, LocalDateTime date);

    boolean existsByEmail(String email);

    long countByStatus(UserStatus status);

    void deleteByEmail(String email);

    List<User> findTop10ByOrderByCreatedAtDesc();
}
```

**Keyword mapping:**
| Keyword | SQL | Example |
|---|---|---|
| `findBy` | `WHERE` | `findByName` |
| `And` | `AND` | `findByNameAndEmail` |
| `Or` | `OR` | `findByNameOrEmail` |
| `Containing` | `LIKE %x%` | `findByNameContaining` |
| `StartingWith` | `LIKE x%` | `findByNameStartingWith` |
| `Between` | `BETWEEN` | `findByAgeBetween` |
| `LessThan` | `<` | `findByAgeLessThan` |
| `GreaterThan` | `>` | `findByAgeGreaterThan` |
| `IsNull` | `IS NULL` | `findByEmailIsNull` |
| `OrderBy` | `ORDER BY` | `findByStatusOrderByNameAsc` |
| `Top` / `First` | `LIMIT` | `findTop5ByStatus` |

---

## Custom JPQL Queries

```java
public interface UserRepository extends JpaRepository<User, Long> {

    // JPQL — uses entity names, not table names
    @Query("SELECT u FROM User u WHERE u.email = :email AND u.status = :status")
    Optional<User> findActiveByEmail(@Param("email") String email,
                                     @Param("status") UserStatus status);

    // Native SQL — when JPQL isn't expressive enough
    @Query(value = "SELECT * FROM users WHERE created_at > NOW() - INTERVAL '7 days'",
           nativeQuery = true)
    List<User> findRecentUsers();

    // Modifying queries
    @Modifying
    @Transactional
    @Query("UPDATE User u SET u.status = :status WHERE u.id = :id")
    int updateStatus(@Param("id") Long id, @Param("status") UserStatus status);
}
```

---

## Pagination and Sorting

```java
// Repository
public interface UserRepository extends JpaRepository<User, Long> {
    Page<User> findByStatus(UserStatus status, Pageable pageable);
}

// Service
@Service
public class UserService {
    public Page<UserDto> getUsers(int page, int size, String sortField) {
        Pageable pageable = PageRequest.of(page, size, Sort.by(sortField).ascending());
        return userRepository.findByStatus(UserStatus.ACTIVE, pageable)
            .map(this::toDto);
    }
}

// Controller
@GetMapping
public Page<UserDto> getUsers(
        @RequestParam(defaultValue = "0") int page,
        @RequestParam(defaultValue = "10") int size,
        @RequestParam(defaultValue = "name") String sort) {
    return userService.getUsers(page, size, sort);
}
```

Page response includes: `content`, `totalElements`, `totalPages`, `number`, `size`.

---

## Entity Relationships

### One-to-Many / Many-to-One

```java
@Entity
public class Post {
    @Id @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private String title;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "author_id")
    private User author;

    @OneToMany(mappedBy = "post", cascade = CascadeType.ALL, orphanRemoval = true)
    private List<Comment> comments = new ArrayList<>();
}

@Entity
public class Comment {
    @Id @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private String content;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "post_id")
    private Post post;
}
```

### Many-to-Many

```java
@Entity
public class User {
    @ManyToMany(fetch = FetchType.LAZY)
    @JoinTable(
        name = "user_roles",
        joinColumns = @JoinColumn(name = "user_id"),
        inverseJoinColumns = @JoinColumn(name = "role_id")
    )
    private Set<Role> roles = new HashSet<>();
}
```

---

## Fetch Types and N+1 Problem

`LAZY` = load on demand (default for collections)
`EAGER` = load immediately (default for @ManyToOne)

**N+1 Problem:** Loading 100 posts then querying for each author separately = 101 queries.

**Fix with JOIN FETCH:**
```java
@Query("SELECT p FROM Post p JOIN FETCH p.author WHERE p.status = 'PUBLISHED'")
List<Post> findPublishedPostsWithAuthors();
```

**Fix with @EntityGraph:**
```java
@EntityGraph(attributePaths = {"author", "comments"})
List<Post> findByStatus(String status);
```

---

## @Transactional

```java
@Service
public class OrderService {

    @Transactional
    public Order placeOrder(Long userId, List<Long> itemIds) {
        User user = userRepository.findById(userId).orElseThrow();
        List<Item> items = itemRepository.findAllById(itemIds);

        Order order = new Order(user, items);
        orderRepository.save(order);

        inventoryService.deduct(items);   // if this throws, order is rolled back

        return order;
    }

    @Transactional(readOnly = true)   // optimization: no dirty check, read replica possible
    public List<Order> getOrders(Long userId) {
        return orderRepository.findByUserId(userId);
    }
}
```

| Property | Options |
|---|---|
| `propagation` | `REQUIRED` (default), `REQUIRES_NEW`, `SUPPORTS`, `NEVER` |
| `isolation` | `DEFAULT`, `READ_COMMITTED`, `REPEATABLE_READ`, `SERIALIZABLE` |
| `rollbackFor` | Exception classes that trigger rollback |
| `readOnly` | `true` for read-only transactions |

---

## Auditing

Automatically track who created/updated entities:

```java
@Configuration
@EnableJpaAuditing
public class JpaConfig {
    @Bean
    public AuditorAware<String> auditorProvider() {
        return () -> Optional.ofNullable(SecurityContextHolder.getContext())
            .map(ctx -> ctx.getAuthentication().getName());
    }
}

@Entity
@EntityListeners(AuditingEntityListener.class)
public class User {
    @CreatedDate
    @Column(updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    private LocalDateTime updatedAt;

    @CreatedBy
    @Column(updatable = false)
    private String createdBy;

    @LastModifiedBy
    private String updatedBy;
}
```

---

## Database Migration with Flyway

Flyway manages database schema changes with versioned SQL scripts. Much better than `ddl-auto=update` in production.

```xml
<dependency>
    <groupId>org.flywaydb</groupId>
    <artifactId>flyway-core</artifactId>
</dependency>
```

```
src/main/resources/db/migration/
├── V1__create_users_table.sql
├── V2__add_email_index.sql
└── V3__add_roles_table.sql
```

```sql
-- V1__create_users_table.sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

Flyway runs missing migrations on startup automatically.

---

## Spring Data with Multiple Databases

```yaml
spring:
  datasource:
    primary:
      url: jdbc:postgresql://primary-host/mydb
    secondary:
      url: jdbc:postgresql://readonly-host/mydb
```

```java
@Configuration
@EnableJpaRepositories(
    basePackages = "com.example.repository.primary",
    entityManagerFactoryRef = "primaryEntityManagerFactory"
)
public class PrimaryDataSourceConfig { ... }
```
