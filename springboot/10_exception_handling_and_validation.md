# Chapter 10: Exception Handling & Validation

## Why Centralized Exception Handling?

Without centralized handling, Spring Boot returns a generic error response. With `@ControllerAdvice`, you get consistent, informative error responses across all controllers.

**Default Spring Boot error response:**
```json
{
  "timestamp": "2024-01-15T10:30:00.000+00:00",
  "status": 500,
  "error": "Internal Server Error",
  "path": "/api/users"
}
```

**Custom error response:**
```json
{
  "status": 404,
  "error": "User Not Found",
  "message": "No user exists with id: 42",
  "path": "/api/users/42",
  "timestamp": "2024-01-15T10:30:00.000+00:00"
}
```

---

## Standard Error Response DTO

```java
public record ErrorResponse(
    int status,
    String error,
    String message,
    String path,
    LocalDateTime timestamp,
    Map<String, String> fieldErrors  // for validation errors
) {
    public static ErrorResponse of(int status, String error, String message, String path) {
        return new ErrorResponse(status, error, message, path, LocalDateTime.now(), null);
    }
}
```

---

## Global Exception Handler with @ControllerAdvice

```java
@RestControllerAdvice
public class GlobalExceptionHandler {

    // Business exceptions
    @ExceptionHandler(ResourceNotFoundException.class)
    @ResponseStatus(HttpStatus.NOT_FOUND)
    public ErrorResponse handleNotFound(ResourceNotFoundException ex, HttpServletRequest request) {
        return ErrorResponse.of(404, "Not Found", ex.getMessage(), request.getRequestURI());
    }

    @ExceptionHandler(DuplicateResourceException.class)
    @ResponseStatus(HttpStatus.CONFLICT)
    public ErrorResponse handleConflict(DuplicateResourceException ex, HttpServletRequest request) {
        return ErrorResponse.of(409, "Conflict", ex.getMessage(), request.getRequestURI());
    }

    @ExceptionHandler(AccessDeniedException.class)
    @ResponseStatus(HttpStatus.FORBIDDEN)
    public ErrorResponse handleForbidden(AccessDeniedException ex, HttpServletRequest request) {
        return ErrorResponse.of(403, "Forbidden", "Access denied", request.getRequestURI());
    }

    // Validation exception — @Valid failures
    @ExceptionHandler(MethodArgumentNotValidException.class)
    @ResponseStatus(HttpStatus.BAD_REQUEST)
    public ErrorResponse handleValidation(MethodArgumentNotValidException ex,
                                          HttpServletRequest request) {
        Map<String, String> fieldErrors = new LinkedHashMap<>();
        ex.getBindingResult().getFieldErrors().forEach(error ->
            fieldErrors.put(error.getField(), error.getDefaultMessage())
        );
        return new ErrorResponse(
            400, "Validation Failed", "Request validation failed",
            request.getRequestURI(), LocalDateTime.now(), fieldErrors
        );
    }

    // @PathVariable / @RequestParam type mismatch
    @ExceptionHandler(MethodArgumentTypeMismatchException.class)
    @ResponseStatus(HttpStatus.BAD_REQUEST)
    public ErrorResponse handleTypeMismatch(MethodArgumentTypeMismatchException ex,
                                            HttpServletRequest request) {
        String message = "Invalid value '" + ex.getValue() + "' for parameter '" + ex.getName() + "'";
        return ErrorResponse.of(400, "Bad Request", message, request.getRequestURI());
    }

    // Fallback for unexpected exceptions
    @ExceptionHandler(Exception.class)
    @ResponseStatus(HttpStatus.INTERNAL_SERVER_ERROR)
    public ErrorResponse handleGeneral(Exception ex, HttpServletRequest request) {
        log.error("Unexpected error", ex);   // log full stack trace
        return ErrorResponse.of(500, "Internal Server Error",
            "An unexpected error occurred", request.getRequestURI());
    }
}
```

---

## Custom Business Exceptions

```java
// Base exception
public abstract class AppException extends RuntimeException {
    private final HttpStatus status;

    protected AppException(String message, HttpStatus status) {
        super(message);
        this.status = status;
    }

    public HttpStatus getStatus() { return status; }
}

// Specific exceptions
public class ResourceNotFoundException extends AppException {
    public ResourceNotFoundException(String resourceName, Long id) {
        super(resourceName + " not found with id: " + id, HttpStatus.NOT_FOUND);
    }
}

public class DuplicateResourceException extends AppException {
    public DuplicateResourceException(String message) {
        super(message, HttpStatus.CONFLICT);
    }
}

public class BusinessRuleException extends AppException {
    public BusinessRuleException(String message) {
        super(message, HttpStatus.UNPROCESSABLE_ENTITY);
    }
}
```

---

## Bean Validation Annotations

```java
public record CreateUserRequest(

    @NotNull(message = "Name cannot be null")
    @NotBlank(message = "Name cannot be blank")
    @Size(min = 2, max = 100, message = "Name must be between 2 and 100 characters")
    String name,

    @NotBlank
    @Email(message = "Must be a valid email address")
    String email,

    @NotBlank
    @Size(min = 8, message = "Password must be at least 8 characters")
    @Pattern(
        regexp = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d).*$",
        message = "Password must contain uppercase, lowercase, and a digit"
    )
    String password,

    @Min(value = 18, message = "Must be at least 18 years old")
    @Max(value = 120)
    Integer age,

    @Past(message = "Birth date must be in the past")
    LocalDate birthDate,

    @Future(message = "Appointment must be in the future")
    LocalDateTime appointment,

    @Positive(message = "Price must be positive")
    BigDecimal price,

    @NotEmpty(message = "At least one tag is required")
    @Size(max = 10, message = "Too many tags")
    List<@NotBlank String> tags
) {}
```

---

## Custom Validators

```java
// 1. Define the annotation
@Documented
@Constraint(validatedBy = UniqueEmailValidator.class)
@Target({ElementType.FIELD})
@Retention(RetentionPolicy.RUNTIME)
public @interface UniqueEmail {
    String message() default "Email is already in use";
    Class<?>[] groups() default {};
    Class<? extends Payload>[] payload() default {};
}

// 2. Implement the validator
@Component
public class UniqueEmailValidator implements ConstraintValidator<UniqueEmail, String> {

    private final UserRepository userRepository;

    public UniqueEmailValidator(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    @Override
    public boolean isValid(String email, ConstraintValidatorContext context) {
        if (email == null) return true;   // let @NotNull handle null
        return !userRepository.existsByEmail(email);
    }
}

// 3. Use it
public record CreateUserRequest(
    @UniqueEmail
    @Email
    String email
) {}
```

### Class-Level Validation (cross-field):
```java
@PasswordMatch(message = "Passwords do not match")
public record ChangePasswordRequest(
    @NotBlank String newPassword,
    @NotBlank String confirmPassword
) {}

@Constraint(validatedBy = PasswordMatchValidator.class)
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
public @interface PasswordMatch {
    String message() default "Passwords do not match";
    Class<?>[] groups() default {};
    Class<? extends Payload>[] payload() default {};
}

public class PasswordMatchValidator implements ConstraintValidator<PasswordMatch, ChangePasswordRequest> {
    @Override
    public boolean isValid(ChangePasswordRequest req, ConstraintValidatorContext ctx) {
        return req.newPassword().equals(req.confirmPassword());
    }
}
```

---

## Validation Groups

Apply different validation rules for create vs update:

```java
public interface OnCreate {}
public interface OnUpdate {}

public record UserRequest(
    @NotNull(groups = OnCreate.class)   // only on create
    @Null(groups = OnUpdate.class)      // must be null on update
    String email,

    @NotBlank
    String name
) {}

// In controller:
@PostMapping
public UserDto create(@Validated(OnCreate.class) @RequestBody UserRequest req) { ... }

@PutMapping("/{id}")
public UserDto update(@PathVariable Long id,
                      @Validated(OnUpdate.class) @RequestBody UserRequest req) { ... }
```

---

## Handling Constraint Violations from Service Layer

```java
@Service
@Validated   // enables method-level validation
public class UserService {

    public UserDto findById(@NotNull @Positive Long id) {
        return userRepository.findById(id)
            .map(this::toDto)
            .orElseThrow(() -> new ResourceNotFoundException("User", id));
    }
}
```

`ConstraintViolationException` is thrown for method-level validation failures — handle it in `@ControllerAdvice`.

---

## Problem Details (RFC 7807) — Spring Boot 3+

Spring Boot 3.x supports RFC 7807 Problem Details out of the box:

```yaml
spring:
  mvc:
    problemdetails:
      enabled: true
```

This produces structured error responses automatically for standard exceptions:
```json
{
  "type": "about:blank",
  "title": "Bad Request",
  "status": 400,
  "detail": "Failed to read request",
  "instance": "/api/users"
}
```

---

## Validation Error Response Example

```
POST /api/users
{
  "name": "",
  "email": "not-an-email",
  "age": 10
}

Response 400:
{
  "status": 400,
  "error": "Validation Failed",
  "message": "Request validation failed",
  "path": "/api/users",
  "timestamp": "2024-01-15T10:30:00",
  "fieldErrors": {
    "name": "Name cannot be blank",
    "email": "Must be a valid email address",
    "age": "Must be at least 18 years old"
  }
}
```
