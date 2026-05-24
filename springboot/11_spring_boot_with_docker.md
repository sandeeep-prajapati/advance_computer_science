# Chapter 11: Spring Boot with Docker

## Why Docker?

Docker packages your application with all its dependencies into a container — ensuring it runs identically on any machine: developer laptop, CI server, production.

---

## Option A: Simple Dockerfile

```dockerfile
# Stage 1: Build
FROM eclipse-temurin:21-jdk-alpine AS builder
WORKDIR /app
COPY . .
RUN ./mvnw clean package -DskipTests

# Stage 2: Runtime (smaller image)
FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY --from=builder /app/target/*.jar app.jar
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]
```

```bash
# Build image
docker build -t my-spring-app:latest .

# Run container
docker run -p 8080:8080 -e SPRING_PROFILES_ACTIVE=prod my-spring-app:latest
```

---

## Option B: Layered JAR (Optimized for Docker)

Spring Boot's layered JAR puts dependencies, Spring Boot libs, and app code in separate layers. Docker caches unchanged layers, making rebuilds fast.

```dockerfile
FROM eclipse-temurin:21-jdk-alpine AS builder
WORKDIR /app
COPY . .
RUN ./mvnw clean package -DskipTests
RUN java -Djarmode=layertools -jar target/*.jar extract --destination target/extracted

FROM eclipse-temurin:21-jre-alpine
WORKDIR /app

# Copy layers in order from least to most frequently changed
COPY --from=builder /app/target/extracted/dependencies/ ./
COPY --from=builder /app/target/extracted/spring-boot-loader/ ./
COPY --from=builder /app/target/extracted/snapshot-dependencies/ ./
COPY --from=builder /app/target/extracted/application/ ./

EXPOSE 8080
ENTRYPOINT ["java", "org.springframework.boot.loader.launch.JarLauncher"]
```

---

## Option C: Cloud Native Buildpacks (No Dockerfile)

Spring Boot 2.3+ can build optimized Docker images without a Dockerfile:

```bash
# Build image directly (requires Docker running)
./mvnw spring-boot:build-image -Dspring-boot.build-image.imageName=my-spring-app:latest

# With Gradle
./gradlew bootBuildImage --imageName=my-spring-app:latest
```

Configure in `pom.xml`:
```xml
<plugin>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-maven-plugin</artifactId>
    <configuration>
        <image>
            <name>my-spring-app:${project.version}</name>
            <builder>paketobuildpacks/builder-jammy-base</builder>
        </image>
        <env>
            <BP_JVM_VERSION>21</BP_JVM_VERSION>
        </env>
    </configuration>
</plugin>
```

---

## Docker Compose for Local Development

```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SPRING_PROFILES_ACTIVE=dev
      - DATABASE_URL=jdbc:postgresql://db:5432/mydb
      - DATABASE_USER=myuser
      - DATABASE_PASS=mypassword
      - REDIS_HOST=redis
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_started
    restart: unless-stopped

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: mydb
      POSTGRES_USER: myuser
      POSTGRES_PASSWORD: mypassword
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U myuser -d mydb"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  mailhog:
    image: mailhog/mailhog
    ports:
      - "1025:1025"   # SMTP
      - "8025:8025"   # Web UI

volumes:
  postgres_data:
```

```bash
# Start all services
docker-compose up -d

# Start only dependencies (run app locally)
docker-compose up -d db redis

# View logs
docker-compose logs -f app

# Rebuild app image
docker-compose up -d --build app

# Tear down
docker-compose down -v    # -v removes volumes
```

---

## Environment Variables in Spring Boot + Docker

```yaml
# application-prod.yml
spring:
  datasource:
    url: ${DATABASE_URL}
    username: ${DATABASE_USER}
    password: ${DATABASE_PASS}
```

Pass via Docker:
```bash
# Inline
docker run -e DATABASE_URL=jdbc:postgresql://host/db my-app

# From .env file
docker run --env-file .env my-app
```

**.env file (never commit to git):**
```env
DATABASE_URL=jdbc:postgresql://db:5432/mydb
DATABASE_USER=myuser
DATABASE_PASS=secret
JWT_SECRET=your-very-long-secret-key
```

---

## JVM Tuning for Containers

Containers have memory limits — JVM must respect them:

```dockerfile
ENTRYPOINT ["java",
  "-XX:MaxRAMPercentage=75.0",
  "-XX:InitialRAMPercentage=50.0",
  "-XX:+UseContainerSupport",
  "-jar", "app.jar"]
```

`-XX:+UseContainerSupport` (default in JDK 11+) makes the JVM read cgroup limits instead of host memory.

For memory-constrained environments:
```bash
java -Xms256m -Xmx512m -jar app.jar
```

---

## Health Checks in Docker

```dockerfile
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
  CMD wget -qO- http://localhost:8080/actuator/health || exit 1
```

Or using curl:
```dockerfile
HEALTHCHECK --interval=30s --timeout=10s \
  CMD curl -f http://localhost:8080/actuator/health || exit 1
```

---

## Multi-Stage Build Optimization Tips

```dockerfile
# Cache Maven/Gradle dependencies as a separate layer
FROM eclipse-temurin:21-jdk-alpine AS deps
WORKDIR /app
COPY pom.xml .
COPY .mvn .mvn
COPY mvnw .
RUN ./mvnw dependency:go-offline -B    # downloads all deps

FROM deps AS builder
COPY src ./src
RUN ./mvnw package -DskipTests

FROM eclipse-temurin:21-jre-alpine
# ... rest of runtime stage
```

Now `dependency:go-offline` layer is cached — only rebuilt when `pom.xml` changes.

---

## Pushing to a Container Registry

```bash
# Docker Hub
docker tag my-spring-app:latest myusername/my-spring-app:latest
docker push myusername/my-spring-app:latest

# Google Container Registry
docker tag my-spring-app:latest gcr.io/my-project/my-spring-app:latest
docker push gcr.io/my-project/my-spring-app:latest

# AWS ECR
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin 123456789.dkr.ecr.us-east-1.amazonaws.com
docker push 123456789.dkr.ecr.us-east-1.amazonaws.com/my-spring-app:latest
```

---

## Docker Image Size Comparison

| Base Image | Size | Notes |
|---|---|---|
| `eclipse-temurin:21-jdk` | ~430 MB | Full JDK — use for building |
| `eclipse-temurin:21-jre` | ~215 MB | JRE only — use for running |
| `eclipse-temurin:21-jre-alpine` | ~75 MB | Alpine-based — smallest |
| `gcr.io/distroless/java21` | ~60 MB | No shell — most secure |
| Spring Boot Buildpacks | ~90 MB | Optimized layers, no Dockerfile |

---

## .dockerignore

```
.git
target/
*.md
.env
.gitignore
Dockerfile
docker-compose*.yml
**/*.log
```
