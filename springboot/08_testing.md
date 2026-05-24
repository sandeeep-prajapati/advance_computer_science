# Chapter 8: Testing Spring Boot Applications

## Testing Philosophy

Spring Boot supports multiple testing approaches at different speeds and scope:

```
Unit Tests          → fast, no Spring context, mock dependencies
Slice Tests         → load partial context (@WebMvcTest, @DataJpaTest)
Integration Tests   → full Spring context (@SpringBootTest)
End-to-End Tests    → full stack with real DB (Testcontainers)
```

---

## Dependencies

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-test</artifactId>
    <scope>test</scope>
</dependency>
```

This includes: JUnit 5, Mockito, MockMvc, AssertJ, Hamcrest, Testcontainers support.

---

## Unit Tests (No Spring Context)

Test pure business logic without starting Spring:

```java
class UserServiceTest {

    @Mock
    private UserRepository userRepository;

    @Mock
    private EmailService emailService;

    @InjectMocks
    private UserService userService;

    @BeforeEach
    void setUp() {
        MockitoAnnotations.openMocks(this);
    }

    @Test
    void createUser_ShouldReturnDto_WhenEmailIsUnique() {
        // Given
        CreateUserRequest request = new CreateUserRequest("Alice", "alice@example.com");
        User savedUser = new User(1L, "Alice", "alice@example.com");
        when(userRepository.existsByEmail("alice@example.com")).thenReturn(false);
        when(userRepository.save(any(User.class))).thenReturn(savedUser);

        // When
        UserDto result = userService.create(request);

        // Then
        assertThat(result.name()).isEqualTo("Alice");
        assertThat(result.email()).isEqualTo("alice@example.com");
        verify(emailService).sendWelcomeEmail("alice@example.com");
    }

    @Test
    void createUser_ShouldThrow_WhenEmailAlreadyExists() {
        when(userRepository.existsByEmail("alice@example.com")).thenReturn(true);
        CreateUserRequest request = new CreateUserRequest("Alice", "alice@example.com");

        assertThatThrownBy(() -> userService.create(request))
            .isInstanceOf(DuplicateEmailException.class)
            .hasMessageContaining("alice@example.com");
    }
}
```

### Use @ExtendWith(MockitoExtension.class) to remove boilerplate:

```java
@ExtendWith(MockitoExtension.class)
class UserServiceTest {

    @Mock
    private UserRepository userRepository;

    @InjectMocks
    private UserService userService;

    // no setUp needed
}
```

---

## @WebMvcTest — Controller Layer Tests

Loads only the web layer (controllers, filters, exception handlers). No real database.

```java
@WebMvcTest(UserController.class)
class UserControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean    // creates a Mockito mock and adds it to Spring context
    private UserService userService;

    @Autowired
    private ObjectMapper objectMapper;

    @Test
    void getUser_ShouldReturn200_WhenUserExists() throws Exception {
        UserDto user = new UserDto(1L, "Alice", "alice@example.com", LocalDateTime.now());
        when(userService.findById(1L)).thenReturn(Optional.of(user));

        mockMvc.perform(get("/api/users/1")
                .contentType(MediaType.APPLICATION_JSON))
            .andExpect(status().isOk())
            .andExpect(jsonPath("$.name").value("Alice"))
            .andExpect(jsonPath("$.email").value("alice@example.com"));
    }

    @Test
    void createUser_ShouldReturn400_WhenEmailIsInvalid() throws Exception {
        CreateUserRequest request = new CreateUserRequest("Alice", "not-an-email");

        mockMvc.perform(post("/api/users")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(request)))
            .andExpect(status().isBadRequest());
    }

    @Test
    void createUser_ShouldReturn201_WithLocationHeader() throws Exception {
        CreateUserRequest request = new CreateUserRequest("Alice", "alice@example.com");
        UserDto created = new UserDto(1L, "Alice", "alice@example.com", LocalDateTime.now());
        when(userService.create(any())).thenReturn(created);

        mockMvc.perform(post("/api/users")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(request)))
            .andExpect(status().isCreated())
            .andExpect(header().string("Location", containsString("/api/users/1")));
    }
}
```

---

## @DataJpaTest — Repository Layer Tests

Loads only JPA components with an in-memory H2 database:

```java
@DataJpaTest
class UserRepositoryTest {

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private TestEntityManager entityManager;

    @Test
    void findByEmail_ShouldReturnUser_WhenExists() {
        // Setup test data
        User user = new User();
        user.setName("Alice");
        user.setEmail("alice@example.com");
        entityManager.persistAndFlush(user);

        // Test
        Optional<User> found = userRepository.findByEmail("alice@example.com");
        assertThat(found).isPresent();
        assertThat(found.get().getName()).isEqualTo("Alice");
    }

    @Test
    void findByStatus_WithPagination_ShouldReturnPage() {
        // create 15 users
        IntStream.range(1, 16).forEach(i -> {
            User u = new User("User" + i, "user" + i + "@example.com", UserStatus.ACTIVE);
            entityManager.persist(u);
        });
        entityManager.flush();

        Page<User> page = userRepository.findByStatus(
            UserStatus.ACTIVE,
            PageRequest.of(0, 10)
        );

        assertThat(page.getTotalElements()).isEqualTo(15);
        assertThat(page.getContent()).hasSize(10);
    }
}
```

---

## @SpringBootTest — Full Integration Tests

Loads the entire Spring context:

```java
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@ActiveProfiles("test")
class UserIntegrationTest {

    @Autowired
    private TestRestTemplate restTemplate;

    @Autowired
    private UserRepository userRepository;

    @BeforeEach
    void setUp() {
        userRepository.deleteAll();
    }

    @Test
    void fullUserLifecycle() {
        // Create
        CreateUserRequest createReq = new CreateUserRequest("Alice", "alice@example.com");
        ResponseEntity<UserDto> createResponse = restTemplate.postForEntity(
            "/api/users", createReq, UserDto.class);
        assertThat(createResponse.getStatusCode()).isEqualTo(HttpStatus.CREATED);

        Long userId = createResponse.getBody().id();

        // Read
        ResponseEntity<UserDto> getResponse = restTemplate.getForEntity(
            "/api/users/" + userId, UserDto.class);
        assertThat(getResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        assertThat(getResponse.getBody().name()).isEqualTo("Alice");

        // Delete
        restTemplate.delete("/api/users/" + userId);
        assertThat(userRepository.existsById(userId)).isFalse();
    }
}
```

---

## Testcontainers — Real Database in Tests

Use real PostgreSQL/MySQL/Redis in tests via Docker containers:

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-testcontainers</artifactId>
    <scope>test</scope>
</dependency>
<dependency>
    <groupId>org.testcontainers</groupId>
    <artifactId>postgresql</artifactId>
    <scope>test</scope>
</dependency>
```

```java
@SpringBootTest
@Testcontainers
class UserRepositoryIntegrationTest {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:16-alpine")
        .withDatabaseName("testdb")
        .withUsername("test")
        .withPassword("test");

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
    }

    @Autowired
    private UserRepository userRepository;

    @Test
    void shouldPersistUser() {
        User user = new User("Alice", "alice@example.com");
        userRepository.save(user);
        assertThat(userRepository.findByEmail("alice@example.com")).isPresent();
    }
}
```

Spring Boot 3.1+ supports `@ServiceConnection` for simpler setup:
```java
@Container
@ServiceConnection
static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:16-alpine");
// No @DynamicPropertySource needed!
```

---

## Testing with Security

```java
@WebMvcTest(UserController.class)
class SecuredControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @Test
    @WithMockUser(username = "alice", roles = "USER")
    void getProfile_ShouldReturn200_ForAuthenticatedUser() throws Exception {
        mockMvc.perform(get("/api/profile"))
            .andExpect(status().isOk());
    }

    @Test
    void getProfile_ShouldReturn401_WhenUnauthenticated() throws Exception {
        mockMvc.perform(get("/api/profile"))
            .andExpect(status().isUnauthorized());
    }

    @Test
    @WithMockUser(roles = "ADMIN")
    void deleteUser_ShouldReturn204_ForAdmin() throws Exception {
        mockMvc.perform(delete("/api/users/1"))
            .andExpect(status().isNoContent());
    }
}
```

---

## application-test.yml

```yaml
spring:
  datasource:
    url: jdbc:h2:mem:testdb;DB_CLOSE_DELAY=-1
    driver-class-name: org.h2.Driver
  jpa:
    hibernate:
      ddl-auto: create-drop
    show-sql: true

logging:
  level:
    root: WARN
    com.example: INFO
```

---

## AssertJ Cheat Sheet

```java
// Basic assertions
assertThat(value).isEqualTo(expected);
assertThat(value).isNotNull();
assertThat(value).isNull();
assertThat(list).hasSize(3);
assertThat(list).isEmpty();
assertThat(list).contains("element");
assertThat(list).containsExactly("a", "b", "c");
assertThat(str).startsWith("Hello");
assertThat(str).contains("world");

// Exceptions
assertThatThrownBy(() -> service.method())
    .isInstanceOf(MyException.class)
    .hasMessage("Expected message");

// Optional
assertThat(optional).isPresent();
assertThat(optional).isEmpty();
assertThat(optional).hasValue("expected");
```

---

## Test Slice Summary

| Annotation | What it loads | Use for |
|---|---|---|
| `@WebMvcTest` | Controllers, filters, security | HTTP layer, request/response |
| `@DataJpaTest` | JPA repositories, H2 | Database queries |
| `@DataMongoTest` | MongoDB repositories | MongoDB queries |
| `@JsonTest` | Jackson serialization | JSON mapping |
| `@RestClientTest` | RestTemplate/RestClient | External HTTP calls |
| `@SpringBootTest` | Full context | Integration / E2E |
