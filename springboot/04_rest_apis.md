# Chapter 4: REST APIs with Spring Boot

## Building REST APIs

Spring Boot's `spring-boot-starter-web` gives you everything to build RESTful APIs:
- **Spring MVC** for routing and request handling
- **Jackson** for JSON serialization/deserialization
- **Embedded Tomcat** as the server

---

## Core Annotations

| Annotation | Purpose |
|---|---|
| `@RestController` | Marks a class as a REST controller (combines @Controller + @ResponseBody) |
| `@RequestMapping` | Maps URL path to a class or method |
| `@GetMapping` | HTTP GET handler |
| `@PostMapping` | HTTP POST handler |
| `@PutMapping` | HTTP PUT handler |
| `@PatchMapping` | HTTP PATCH handler |
| `@DeleteMapping` | HTTP DELETE handler |
| `@PathVariable` | Extracts value from URL path |
| `@RequestParam` | Extracts query parameter |
| `@RequestBody` | Deserializes request body from JSON |
| `@ResponseStatus` | Sets HTTP status code on method/exception |

---

## A Complete CRUD Controller

```java
@RestController
@RequestMapping("/api/users")
public class UserController {

    private final UserService userService;

    public UserController(UserService userService) {
        this.userService = userService;
    }

    // GET /api/users
    @GetMapping
    public List<UserDto> getAllUsers() {
        return userService.findAll();
    }

    // GET /api/users/42
    @GetMapping("/{id}")
    public ResponseEntity<UserDto> getUserById(@PathVariable Long id) {
        return userService.findById(id)
            .map(ResponseEntity::ok)
            .orElse(ResponseEntity.notFound().build());
    }

    // GET /api/users?page=0&size=10&sort=name
    @GetMapping("/search")
    public Page<UserDto> searchUsers(
            @RequestParam(defaultValue = "") String name,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size) {
        return userService.search(name, PageRequest.of(page, size));
    }

    // POST /api/users
    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public UserDto createUser(@Valid @RequestBody CreateUserRequest request) {
        return userService.create(request);
    }

    // PUT /api/users/42
    @PutMapping("/{id}")
    public ResponseEntity<UserDto> updateUser(
            @PathVariable Long id,
            @Valid @RequestBody UpdateUserRequest request) {
        return userService.update(id, request)
            .map(ResponseEntity::ok)
            .orElse(ResponseEntity.notFound().build());
    }

    // DELETE /api/users/42
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteUser(@PathVariable Long id) {
        if (userService.delete(id)) {
            return ResponseEntity.noContent().build();
        }
        return ResponseEntity.notFound().build();
    }
}
```

---

## DTOs (Data Transfer Objects)

Never expose your `@Entity` classes directly in API responses. Use DTOs to:
- Control what fields are exposed
- Decouple API contract from database schema
- Add validation annotations

```java
// Request DTO
public record CreateUserRequest(
    @NotBlank @Size(max = 100) String name,
    @Email @NotBlank String email,
    @NotBlank @Size(min = 8) String password
) {}

// Response DTO
public record UserDto(
    Long id,
    String name,
    String email,
    LocalDateTime createdAt
) {}
```

**Mapping entity to DTO:**
```java
@Service
public class UserService {
    private UserDto toDto(User user) {
        return new UserDto(user.getId(), user.getName(), user.getEmail(), user.getCreatedAt());
    }
}
```

For larger projects, consider **MapStruct** to auto-generate mappers at compile time.

---

## ResponseEntity — Full Control Over Responses

`ResponseEntity<T>` lets you control status code, headers, and body:

```java
@GetMapping("/{id}")
public ResponseEntity<UserDto> getUser(@PathVariable Long id) {
    return userService.findById(id)
        .map(user -> ResponseEntity
            .ok()
            .header("X-User-Found", "true")
            .body(user))
        .orElse(ResponseEntity.notFound().build());
}

@PostMapping
public ResponseEntity<UserDto> createUser(@RequestBody CreateUserRequest req) {
    UserDto created = userService.create(req);
    URI location = URI.create("/api/users/" + created.id());
    return ResponseEntity.created(location).body(created);
}
```

---

## HTTP Status Codes Cheat Sheet

| Status | Meaning | When to Use |
|---|---|---|
| `200 OK` | Success | GET, PUT, PATCH success |
| `201 Created` | Resource created | POST success |
| `204 No Content` | Success, no body | DELETE success |
| `400 Bad Request` | Invalid input | Validation failure |
| `401 Unauthorized` | Not authenticated | Missing/invalid token |
| `403 Forbidden` | Not authorized | Authenticated but no permission |
| `404 Not Found` | Resource missing | Entity not found |
| `409 Conflict` | State conflict | Duplicate resource |
| `422 Unprocessable Entity` | Semantic validation error | Business rule violation |
| `500 Internal Server Error` | Server error | Unexpected exception |

---

## Request Validation with Bean Validation

Add `spring-boot-starter-validation` to enable Bean Validation (Jakarta Validation API):

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-validation</artifactId>
</dependency>
```

```java
public record CreateUserRequest(
    @NotBlank(message = "Name is required")
    @Size(max = 100, message = "Name must be under 100 characters")
    String name,

    @Email(message = "Invalid email format")
    @NotBlank
    String email,

    @Min(value = 0, message = "Age must be positive")
    @Max(value = 150)
    Integer age,

    @Pattern(regexp = "^[A-Z]{2}$", message = "Country code must be 2 uppercase letters")
    String countryCode
) {}
```

Use `@Valid` on the controller parameter to trigger validation:
```java
@PostMapping
public UserDto create(@Valid @RequestBody CreateUserRequest request) { ... }
```

Validation failures throw `MethodArgumentNotValidException` — handle it in `@ControllerAdvice` (see Chapter 10).

---

## Path Variables vs Request Params

```
GET /api/users/42           → @PathVariable Long id = 42
GET /api/users?role=ADMIN   → @RequestParam String role = "ADMIN"
GET /api/users/42/posts     → nested resource
```

**Rule of thumb:**
- Use **path variables** for resource identity (`/users/42`)
- Use **request params** for filtering, sorting, pagination (`?sort=name&page=2`)

---

## Consuming External REST APIs (RestTemplate / RestClient)

### RestClient (Spring Boot 3.2+ — recommended):
```java
@Service
public class GithubService {

    private final RestClient restClient;

    public GithubService(RestClient.Builder builder) {
        this.restClient = builder
            .baseUrl("https://api.github.com")
            .defaultHeader("Accept", "application/vnd.github.v3+json")
            .build();
    }

    public GithubUser getUser(String username) {
        return restClient.get()
            .uri("/users/{username}", username)
            .retrieve()
            .body(GithubUser.class);
    }
}
```

### WebClient (reactive / non-blocking):
```java
WebClient webClient = WebClient.builder()
    .baseUrl("https://api.github.com")
    .build();

Mono<GithubUser> user = webClient.get()
    .uri("/users/{username}", username)
    .retrieve()
    .bodyToMono(GithubUser.class);
```

---

## Content Negotiation

Spring Boot auto-configures JSON (Jackson). For XML support:
```xml
<dependency>
    <groupId>com.fasterxml.jackson.dataformat</groupId>
    <artifactId>jackson-dataformat-xml</artifactId>
</dependency>
```

Clients request format via `Accept` header:
```
Accept: application/json   → JSON response
Accept: application/xml    → XML response
```

---

## CORS Configuration

```java
@Configuration
public class CorsConfig {

    @Bean
    public CorsConfigurationSource corsConfigurationSource() {
        CorsConfiguration config = new CorsConfiguration();
        config.setAllowedOrigins(List.of("https://myfrontend.com"));
        config.setAllowedMethods(List.of("GET", "POST", "PUT", "DELETE", "OPTIONS"));
        config.setAllowedHeaders(List.of("*"));
        config.setAllowCredentials(true);

        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
        source.registerCorsConfiguration("/api/**", config);
        return source;
    }
}
```

Or per-controller:
```java
@RestController
@CrossOrigin(origins = "https://myfrontend.com")
public class UserController { ... }
```

---

## API Versioning Strategies

| Strategy | Example | Pros/Cons |
|---|---|---|
| **URI versioning** | `/api/v1/users` | Simple, visible, most common |
| **Header versioning** | `Accept: application/vnd.myapi.v1+json` | Clean URLs, harder to test |
| **Query param** | `/api/users?version=1` | Easy but pollutes URL |

```java
// URI versioning example
@RestController
@RequestMapping("/api/v1/users")
public class UserControllerV1 { ... }

@RestController
@RequestMapping("/api/v2/users")
public class UserControllerV2 { ... }
```

---

## Useful Jackson Annotations

```java
public class UserDto {
    @JsonProperty("user_id")       // rename field in JSON
    private Long id;

    @JsonIgnore                    // exclude from serialization
    private String password;

    @JsonFormat(pattern = "yyyy-MM-dd HH:mm:ss")
    private LocalDateTime createdAt;

    @JsonInclude(JsonInclude.Include.NON_NULL)  // omit null fields
    private String middleName;
}
```

Global Jackson config in `application.yml`:
```yaml
spring:
  jackson:
    default-property-inclusion: non_null
    serialization:
      write-dates-as-timestamps: false
    deserialization:
      fail-on-unknown-properties: false
```
