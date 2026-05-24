# Chapter 2: Setup & Installation

## Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| JDK | 17 or 21 (LTS) | Spring Boot 3.x requires Java 17+ |
| Maven or Gradle | Latest | Build tool |
| IDE | IntelliJ IDEA / VS Code / Eclipse | Development environment |
| Docker | Latest | For containers and Testcontainers |

---

## Option A: Spring Initializr (Recommended)

The fastest way to create a Spring Boot project is at **start.spring.io**.

### Steps:
1. Go to [start.spring.io](https://start.spring.io)
2. Choose:
   - **Project:** Maven or Gradle
   - **Language:** Java (or Kotlin/Groovy)
   - **Spring Boot version:** Latest stable (3.x)
   - **Group:** `com.example`
   - **Artifact:** `my-app`
   - **Java version:** 17 or 21
3. Add Dependencies: `Spring Web`, `Spring Data JPA`, `H2 Database`
4. Click **Generate** → download ZIP → extract → open in IDE

### Via curl (terminal):
```bash
curl https://start.spring.io/starter.zip \
  -d type=maven-project \
  -d language=java \
  -d bootVersion=3.3.0 \
  -d baseDir=my-app \
  -d groupId=com.example \
  -d artifactId=my-app \
  -d dependencies=web,data-jpa,h2 \
  -o my-app.zip && unzip my-app.zip
```

---

## Option B: IntelliJ IDEA (Ultimate)

IntelliJ IDEA Ultimate has Spring Initializr built in:

1. **File → New Project → Spring**
2. Fill in project metadata
3. Select dependencies
4. Click **Create**

---

## Project Structure

```
my-app/
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/example/myapp/
│   │   │       ├── MyAppApplication.java     ← entry point
│   │   │       ├── controller/               ← REST controllers
│   │   │       ├── service/                  ← business logic
│   │   │       ├── repository/               ← data access
│   │   │       └── model/                    ← entities/DTOs
│   │   └── resources/
│   │       ├── application.properties        ← config
│   │       ├── application.yml               ← alternative config (YAML)
│   │       └── static/                       ← static files (HTML, CSS, JS)
│   └── test/
│       └── java/
│           └── com/example/myapp/
│               └── MyAppApplicationTests.java
├── pom.xml                                   ← Maven build file
└── mvnw / gradlew                            ← Maven/Gradle wrapper
```

---

## pom.xml Explained

```xml
<?xml version="1.0" encoding="UTF-8"?>
<project>
    <modelVersion>4.0.0</modelVersion>

    <!-- Spring Boot parent manages all dependency versions -->
    <parent>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-parent</artifactId>
        <version>3.3.0</version>
    </parent>

    <groupId>com.example</groupId>
    <artifactId>my-app</artifactId>
    <version>0.0.1-SNAPSHOT</version>

    <properties>
        <java.version>17</java.version>
    </properties>

    <dependencies>
        <!-- Web: REST API support + embedded Tomcat -->
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>

        <!-- JPA: Hibernate + Spring Data -->
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-data-jpa</artifactId>
        </dependency>

        <!-- H2: In-memory database (dev only) -->
        <dependency>
            <groupId>com.h2database</groupId>
            <artifactId>h2</artifactId>
            <scope>runtime</scope>
        </dependency>

        <!-- Testing: JUnit 5 + Mockito + MockMvc -->
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-test</artifactId>
            <scope>test</scope>
        </dependency>
    </dependencies>

    <build>
        <plugins>
            <!-- Creates executable fat JAR -->
            <plugin>
                <groupId>org.springframework.boot</groupId>
                <artifactId>spring-boot-maven-plugin</artifactId>
            </plugin>
        </plugins>
    </build>
</project>
```

---

## Maven vs Gradle

| Feature | Maven | Gradle |
|---|---|---|
| Config format | XML (pom.xml) | Groovy/Kotlin DSL (build.gradle) |
| Learning curve | Lower | Slightly higher |
| Build speed | Slower | Faster (incremental, parallel) |
| Flexibility | Convention-heavy | Highly customizable |
| Industry adoption | Very common | Growing fast (Android, large projects) |

### Equivalent Gradle (Kotlin DSL):

```kotlin
// build.gradle.kts
plugins {
    id("org.springframework.boot") version "3.3.0"
    id("io.spring.dependency-management") version "1.1.5"
    kotlin("jvm") version "1.9.24"
    kotlin("plugin.spring") version "1.9.24"
}

dependencies {
    implementation("org.springframework.boot:spring-boot-starter-web")
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")
    runtimeOnly("com.h2database:h2")
    testImplementation("org.springframework.boot:spring-boot-starter-test")
}
```

---

## Running the Application

### With Maven wrapper:
```bash
./mvnw spring-boot:run
```

### With Gradle wrapper:
```bash
./gradlew bootRun
```

### Build and run JAR:
```bash
./mvnw clean package
java -jar target/my-app-0.0.1-SNAPSHOT.jar
```

### With IntelliJ IDEA:
- Click the green play button next to `MyAppApplication.main()`

---

## Your First REST Endpoint

Edit `MyAppApplication.java` or create `HelloController.java`:

```java
package com.example.myapp.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class HelloController {

    @GetMapping("/hello")
    public String hello() {
        return "Hello, Spring Boot!";
    }
}
```

Run the app and test:
```bash
curl http://localhost:8080/hello
# Response: Hello, Spring Boot!
```

---

## application.properties vs application.yml

Both are equivalent — use whichever you prefer.

**application.properties:**
```properties
server.port=8080
spring.datasource.url=jdbc:h2:mem:testdb
spring.datasource.driver-class-name=org.h2.Driver
spring.jpa.show-sql=true
spring.jpa.hibernate.ddl-auto=update
```

**application.yml (equivalent):**
```yaml
server:
  port: 8080

spring:
  datasource:
    url: jdbc:h2:mem:testdb
    driver-class-name: org.h2.Driver
  jpa:
    show-sql: true
    hibernate:
      ddl-auto: update
```

YAML is preferred for complex nested config. Properties is simpler for flat config.

---

## Hot Reload with DevTools

Add to `pom.xml`:
```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-devtools</artifactId>
    <scope>runtime</scope>
    <optional>true</optional>
</dependency>
```

DevTools provides:
- **Automatic restart** when classpath files change
- **LiveReload** browser refresh for static resources
- **Relaxed property binding** in dev mode

In IntelliJ: Enable "Build project automatically" and "Allow auto-make to start even if developed application is currently running".

---

## Install Java (if needed)

### Ubuntu/Debian:
```bash
sudo apt update
sudo apt install openjdk-21-jdk
java -version
```

### macOS (SDKMAN — recommended):
```bash
curl -s "https://get.sdkman.io" | bash
sdk install java 21.0.3-tem
sdk use java 21.0.3-tem
java -version
```

### Windows:
Download from [adoptium.net](https://adoptium.net) and run the installer.
