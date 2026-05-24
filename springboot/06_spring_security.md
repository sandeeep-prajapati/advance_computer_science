# Chapter 6: Spring Security

## What is Spring Security?

**Spring Security** is a powerful, highly customizable authentication and authorization framework for Spring applications. It handles:
- **Authentication** — who are you? (login, JWT, OAuth2)
- **Authorization** — what can you do? (roles, permissions)
- **Protection** — CSRF, session fixation, clickjacking, etc.

---

## Dependency

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-security</artifactId>
</dependency>
```

Adding this dependency alone immediately:
- Blocks all requests by default (requires authentication)
- Generates a random password printed to the console at startup
- Enables HTTP Basic authentication

---

## Security Filter Chain

Spring Security is a chain of filters that process every HTTP request. You configure it with a `SecurityFilterChain` bean:

```java
@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http
            .csrf(AbstractHttpConfigurer::disable)          // disable for REST APIs
            .sessionManagement(session -> session
                .sessionCreationPolicy(SessionCreationPolicy.STATELESS)) // no sessions for JWT
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/api/auth/**").permitAll()              // public endpoints
                .requestMatchers("/api/admin/**").hasRole("ADMIN")        // admin only
                .requestMatchers(HttpMethod.GET, "/api/products/**").permitAll()
                .anyRequest().authenticated()                             // all else need auth
            )
            .addFilterBefore(jwtAuthFilter, UsernamePasswordAuthenticationFilter.class);

        return http.build();
    }
}
```

---

## JWT Authentication (Most Common for REST APIs)

### Flow:
```
Client → POST /api/auth/login {username, password}
       ← 200 OK { token: "eyJhbGci..." }

Client → GET /api/users (Authorization: Bearer eyJhbGci...)
       → JwtAuthFilter validates token
       → SecurityContext populated with user
       ← 200 OK [users...]
```

### Dependencies:
```xml
<dependency>
    <groupId>io.jsonwebtoken</groupId>
    <artifactId>jjwt-api</artifactId>
    <version>0.12.5</version>
</dependency>
<dependency>
    <groupId>io.jsonwebtoken</groupId>
    <artifactId>jjwt-impl</artifactId>
    <version>0.12.5</version>
    <scope>runtime</scope>
</dependency>
<dependency>
    <groupId>io.jsonwebtoken</groupId>
    <artifactId>jjwt-jackson</artifactId>
    <version>0.12.5</version>
    <scope>runtime</scope>
</dependency>
```

### JWT Service:
```java
@Service
public class JwtService {

    @Value("${app.jwt.secret}")
    private String secretKey;

    @Value("${app.jwt.expiration-ms:3600000}")
    private long expirationMs;

    public String generateToken(UserDetails userDetails) {
        return Jwts.builder()
            .subject(userDetails.getUsername())
            .issuedAt(new Date())
            .expiration(new Date(System.currentTimeMillis() + expirationMs))
            .signWith(getSigningKey())
            .compact();
    }

    public String extractUsername(String token) {
        return getClaims(token).getSubject();
    }

    public boolean isTokenValid(String token, UserDetails userDetails) {
        return extractUsername(token).equals(userDetails.getUsername())
            && !isTokenExpired(token);
    }

    private boolean isTokenExpired(String token) {
        return getClaims(token).getExpiration().before(new Date());
    }

    private Claims getClaims(String token) {
        return Jwts.parser().verifyWith(getSigningKey()).build().parseSignedClaims(token).getPayload();
    }

    private SecretKey getSigningKey() {
        return Keys.hmacShaKeyFor(Decoders.BASE64.decode(secretKey));
    }
}
```

### JWT Filter:
```java
@Component
public class JwtAuthFilter extends OncePerRequestFilter {

    private final JwtService jwtService;
    private final UserDetailsService userDetailsService;

    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain chain) throws ServletException, IOException {
        String authHeader = request.getHeader("Authorization");

        if (authHeader == null || !authHeader.startsWith("Bearer ")) {
            chain.doFilter(request, response);
            return;
        }

        String token = authHeader.substring(7);
        String username = jwtService.extractUsername(token);

        if (username != null && SecurityContextHolder.getContext().getAuthentication() == null) {
            UserDetails userDetails = userDetailsService.loadUserByUsername(username);
            if (jwtService.isTokenValid(token, userDetails)) {
                UsernamePasswordAuthenticationToken authToken =
                    new UsernamePasswordAuthenticationToken(
                        userDetails, null, userDetails.getAuthorities());
                authToken.setDetails(new WebAuthenticationDetailsSource().buildDetails(request));
                SecurityContextHolder.getContext().setAuthentication(authToken);
            }
        }

        chain.doFilter(request, response);
    }
}
```

### Auth Controller:
```java
@RestController
@RequestMapping("/api/auth")
public class AuthController {

    private final AuthenticationManager authenticationManager;
    private final JwtService jwtService;
    private final UserDetailsService userDetailsService;

    @PostMapping("/login")
    public ResponseEntity<AuthResponse> login(@RequestBody AuthRequest request) {
        authenticationManager.authenticate(
            new UsernamePasswordAuthenticationToken(request.username(), request.password())
        );
        UserDetails userDetails = userDetailsService.loadUserByUsername(request.username());
        String token = jwtService.generateToken(userDetails);
        return ResponseEntity.ok(new AuthResponse(token));
    }
}
```

---

## Custom UserDetailsService

```java
@Service
public class CustomUserDetailsService implements UserDetailsService {

    private final UserRepository userRepository;

    @Override
    public UserDetails loadUserByUsername(String email) throws UsernameNotFoundException {
        return userRepository.findByEmail(email)
            .map(user -> org.springframework.security.core.userdetails.User.builder()
                .username(user.getEmail())
                .password(user.getPassword())
                .roles(user.getRole().name())
                .build())
            .orElseThrow(() -> new UsernameNotFoundException("User not found: " + email));
    }
}
```

---

## Password Encoding

**Never store plain-text passwords.**

```java
@Bean
public PasswordEncoder passwordEncoder() {
    return new BCryptPasswordEncoder();  // uses bcrypt with cost factor 10
}

// In registration:
user.setPassword(passwordEncoder.encode(rawPassword));

// Spring Security compares automatically during authentication
```

---

## Role-Based Authorization

```java
// In SecurityFilterChain:
.authorizeHttpRequests(auth -> auth
    .requestMatchers("/api/admin/**").hasRole("ADMIN")
    .requestMatchers("/api/manager/**").hasAnyRole("ADMIN", "MANAGER")
    .anyRequest().authenticated()
)
```

### Method-Level Security

```java
@Configuration
@EnableMethodSecurity   // must enable this
public class SecurityConfig { ... }

@Service
public class OrderService {

    @PreAuthorize("hasRole('ADMIN') or #userId == authentication.name")
    public List<Order> getOrders(String userId) { ... }

    @PreAuthorize("hasRole('ADMIN')")
    public void deleteOrder(Long id) { ... }

    @PostAuthorize("returnObject.userId == authentication.name")
    public Order getOrder(Long id) { ... }
}
```

---

## OAuth2 Login (Social Login)

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-oauth2-client</artifactId>
</dependency>
```

```yaml
spring:
  security:
    oauth2:
      client:
        registration:
          google:
            client-id: ${GOOGLE_CLIENT_ID}
            client-secret: ${GOOGLE_CLIENT_SECRET}
            scope: email, profile
          github:
            client-id: ${GITHUB_CLIENT_ID}
            client-secret: ${GITHUB_CLIENT_SECRET}
```

```java
http.oauth2Login(oauth2 -> oauth2
    .loginPage("/login")
    .defaultSuccessUrl("/dashboard")
);
```

---

## Securing Actuator Endpoints

```yaml
management:
  endpoints:
    web:
      exposure:
        include: health, info, metrics
```

```java
http.authorizeHttpRequests(auth -> auth
    .requestMatchers("/actuator/health").permitAll()
    .requestMatchers("/actuator/**").hasRole("ADMIN")
    .anyRequest().authenticated()
);
```

---

## CORS + Security

Spring Security overrides any CORS config you set elsewhere. Configure CORS inside the SecurityFilterChain:

```java
http.cors(cors -> cors.configurationSource(corsConfigurationSource()));
```

---

## Common Security Configurations

```java
// Permit static resources
.requestMatchers("/static/**", "/css/**", "/js/**", "/images/**").permitAll()

// H2 console (dev only)
.requestMatchers("/h2-console/**").permitAll()
http.headers(headers -> headers.frameOptions(FrameOptionsConfig::disable));

// Remember-me
http.rememberMe(remember -> remember.key("my-secret-key").tokenValiditySeconds(86400));
```

---

## Security Best Practices

- Always use HTTPS in production
- Store JWT secret in environment variables, not in code
- Use short-lived access tokens (15 min) + refresh tokens (7 days)
- Never log raw passwords or tokens
- Use `@Transactional(readOnly = true)` for `loadUserByUsername`
- Rate-limit login endpoints to prevent brute force
- Use `BCryptPasswordEncoder` — never MD5/SHA1
