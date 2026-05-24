# Chapter 7: Configuration & Profiles

## Configuration Sources (Priority Order)

Spring Boot loads configuration from multiple sources. Higher items override lower items:

```
1. Command-line arguments              (--server.port=9090)
2. OS environment variables            (SERVER_PORT=9090)
3. application-{profile}.properties   (profile-specific)
4. application.properties / yml       (base config)
5. @PropertySource annotations
6. Default values in @Value
```

---

## application.properties vs application.yml

Both are equivalent. Pick one style and be consistent.

**application.properties:**
```properties
server.port=8080
server.servlet.context-path=/api

spring.datasource.url=jdbc:postgresql://localhost:5432/mydb
spring.datasource.username=myuser
spring.datasource.password=${DB_PASSWORD}

spring.jpa.show-sql=false
spring.jpa.hibernate.ddl-auto=validate

logging.level.com.example=DEBUG
logging.level.org.hibernate.SQL=DEBUG
```

**application.yml (equivalent):**
```yaml
server:
  port: 8080
  servlet:
    context-path: /api

spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/mydb
    username: myuser
    password: ${DB_PASSWORD}
  jpa:
    show-sql: false
    hibernate:
      ddl-auto: validate

logging:
  level:
    com.example: DEBUG
    org.hibernate.SQL: DEBUG
```

---

## Profiles

Profiles let you have different configurations for different environments.

```
application.yml           ← base config (loaded always)
application-dev.yml       ← dev overrides
application-test.yml      ← test overrides
application-prod.yml      ← production overrides
```

### application-dev.yml:
```yaml
spring:
  datasource:
    url: jdbc:h2:mem:devdb
  jpa:
    show-sql: true
    hibernate:
      ddl-auto: create-drop

logging:
  level:
    com.example: DEBUG
```

### application-prod.yml:
```yaml
spring:
  datasource:
    url: ${DATABASE_URL}
    username: ${DATABASE_USER}
    password: ${DATABASE_PASS}
    hikari:
      maximum-pool-size: 20
  jpa:
    show-sql: false
    hibernate:
      ddl-auto: validate

logging:
  level:
    root: WARN
    com.example: INFO
```

### Activating a Profile

```bash
# Command line
java -jar app.jar --spring.profiles.active=prod

# Environment variable
export SPRING_PROFILES_ACTIVE=prod
java -jar app.jar

# In application.yml
spring:
  profiles:
    active: dev   # default profile
```

### In tests:
```java
@SpringBootTest
@ActiveProfiles("test")
class UserServiceTest { ... }
```

---

## @ConfigurationProperties — Typed Configuration

Instead of scattered `@Value` annotations, bind entire config sections to a typed class.

```yaml
# application.yml
app:
  mail:
    host: smtp.gmail.com
    port: 587
    username: noreply@example.com
    password: ${MAIL_PASSWORD}
    from-name: My App
  jwt:
    secret: ${JWT_SECRET}
    expiration-ms: 3600000
  feature:
    signup-enabled: true
    max-users-per-org: 100
```

```java
@ConfigurationProperties(prefix = "app.mail")
@Validated   // enables Bean Validation on config fields
public record MailProperties(
    @NotBlank String host,
    @Min(1) @Max(65535) int port,
    @Email String username,
    @NotBlank String password,
    String fromName
) {}

@ConfigurationProperties(prefix = "app.jwt")
public record JwtProperties(
    @NotBlank String secret,
    @Positive long expirationMs
) {}
```

Enable with `@EnableConfigurationProperties` or just add `@Component`:
```java
@Configuration
@EnableConfigurationProperties({MailProperties.class, JwtProperties.class})
public class AppConfig {}
```

Use them anywhere via injection:
```java
@Service
public class EmailService {
    private final MailProperties mail;

    public EmailService(MailProperties mail) {
        this.mail = mail;
    }
}
```

---

## @Value Annotation

For simple single-value injection:

```java
@Service
public class MyService {

    @Value("${app.name:My Default App}")    // with default value
    private String appName;

    @Value("${server.port}")
    private int port;

    @Value("${feature.flags}")              // comma-separated → List
    private List<String> featureFlags;

    @Value("#{systemProperties['user.home']}") // Spring Expression Language
    private String userHome;
}
```

Prefer `@ConfigurationProperties` over `@Value` for groups of related config.

---

## Externalized Configuration (12-Factor App)

In production, configuration should come from environment variables — never hardcoded or committed to git.

```yaml
# application-prod.yml — references env vars
spring:
  datasource:
    url: ${DATABASE_URL}
    username: ${DATABASE_USER}
    password: ${DATABASE_PASS}

app:
  jwt:
    secret: ${JWT_SECRET}
  stripe:
    api-key: ${STRIPE_API_KEY}
```

### In Docker / Kubernetes:
```bash
# Docker
docker run -e DATABASE_URL=jdbc:postgresql://db/mydb \
           -e DATABASE_USER=admin \
           -e DATABASE_PASS=secret \
           my-app.jar

# Kubernetes Secret + env ref
env:
  - name: DATABASE_PASS
    valueFrom:
      secretKeyRef:
        name: db-secret
        key: password
```

---

## Relaxed Binding

Spring Boot accepts config keys in multiple formats — they all map to the same property:

```
# All equivalent:
spring.datasource.driver-class-name=...
spring.datasource.driverClassName=...
SPRING_DATASOURCE_DRIVER_CLASS_NAME=...       ← env var
--spring.datasource.driver-class-name=...     ← CLI arg
```

---

## Configuration Metadata (IDE Support)

Create `META-INF/spring-configuration-metadata.json` or use the annotation processor:

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-configuration-processor</artifactId>
    <optional>true</optional>
</dependency>
```

This generates metadata for `@ConfigurationProperties` classes, enabling IDE autocomplete in `application.yml`.

---

## Logging Configuration

```yaml
logging:
  level:
    root: INFO
    com.example.myapp: DEBUG
    org.springframework.security: DEBUG
    org.hibernate.SQL: DEBUG
    org.hibernate.orm.jdbc.bind: TRACE    # logs parameter values
  file:
    name: logs/application.log
  pattern:
    console: "%d{HH:mm:ss} [%thread] %-5level %logger{36} - %msg%n"
```

### Custom Logback config (`logback-spring.xml`):
```xml
<configuration>
    <springProfile name="prod">
        <appender name="JSON" class="ch.qos.logback.core.ConsoleAppender">
            <encoder class="net.logstash.logback.encoder.LogstashEncoder"/>
        </appender>
        <root level="INFO">
            <appender-ref ref="JSON"/>
        </root>
    </springProfile>
</configuration>
```

---

## Profile-Specific Beans

```java
@Configuration
public class StorageConfig {

    @Bean
    @Profile("dev")
    public StorageService localStorageService() {
        return new LocalFileStorageService("/tmp/uploads");
    }

    @Bean
    @Profile("prod")
    public StorageService s3StorageService(S3Client s3Client) {
        return new S3StorageService(s3Client, "my-s3-bucket");
    }
}
```

---

## Config Server (Centralized Configuration)

For microservices, use **Spring Cloud Config Server** to serve config from a Git repository:

```yaml
# bootstrap.yml (on each microservice)
spring:
  config:
    import: configserver:http://config-server:8888
  application:
    name: user-service
  profiles:
    active: prod
```

The config server fetches `user-service-prod.yml` from Git and serves it to the app.
