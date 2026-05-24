# Chapter 18: Spring Boot with GraphQL

## What is GraphQL?

**GraphQL** is an API query language that lets clients request exactly the data they need — no more, no less. Unlike REST where the server decides what's in the response, GraphQL gives clients control.

```
REST:                           GraphQL:
GET /users/1                    query {
GET /users/1/orders               user(id: 1) {
GET /users/1/orders/5/items         id
                                    name
= 3 requests, possible overfetch    orders {
                                      id
                                      totalAmount
                                      items { name, price }
                                    }
                                  }
                                }
                                = 1 request, exactly what's needed
```

---

## Dependencies

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-graphql</artifactId>
</dependency>
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-web</artifactId>
</dependency>
```

---

## GraphQL Schema

Define your schema in `src/main/resources/graphql/schema.graphqls`:

```graphql
type Query {
    user(id: ID!): User
    users(page: Int = 0, size: Int = 10): UserPage!
    searchUsers(name: String!): [User!]!
}

type Mutation {
    createUser(input: CreateUserInput!): User!
    updateUser(id: ID!, input: UpdateUserInput!): User!
    deleteUser(id: ID!): Boolean!
}

type Subscription {
    orderStatusChanged(orderId: ID!): Order!
}

type User {
    id: ID!
    name: String!
    email: String!
    createdAt: String!
    orders: [Order!]!
    orderCount: Int!
}

type Order {
    id: ID!
    totalAmount: Float!
    status: OrderStatus!
    items: [OrderItem!]!
    createdAt: String!
}

type OrderItem {
    id: ID!
    productName: String!
    quantity: Int!
    price: Float!
}

type UserPage {
    content: [User!]!
    totalElements: Int!
    totalPages: Int!
    currentPage: Int!
}

enum OrderStatus {
    PENDING
    CONFIRMED
    SHIPPED
    DELIVERED
    CANCELLED
}

input CreateUserInput {
    name: String!
    email: String!
    password: String!
}

input UpdateUserInput {
    name: String
    email: String
}
```

---

## Query Controllers

```java
@Controller
public class UserQueryController {

    private final UserService userService;

    // Maps to Query.user
    @QueryMapping
    public UserDto user(@Argument Long id) {
        return userService.findById(id)
            .orElseThrow(() -> new ResourceNotFoundException("User", id));
    }

    // Maps to Query.users
    @QueryMapping
    public Page<UserDto> users(@Argument int page, @Argument int size) {
        return userService.findAll(PageRequest.of(page, size));
    }

    // Maps to Query.searchUsers
    @QueryMapping
    public List<UserDto> searchUsers(@Argument String name) {
        return userService.search(name);
    }
}
```

---

## Mutation Controllers

```java
@Controller
public class UserMutationController {

    private final UserService userService;

    // Maps to Mutation.createUser
    @MutationMapping
    public UserDto createUser(@Argument @Valid CreateUserInput input) {
        return userService.create(input);
    }

    // Maps to Mutation.updateUser
    @MutationMapping
    public UserDto updateUser(@Argument Long id, @Argument UpdateUserInput input) {
        return userService.update(id, input);
    }

    // Maps to Mutation.deleteUser
    @MutationMapping
    public boolean deleteUser(@Argument Long id) {
        return userService.delete(id);
    }
}
```

---

## Resolving Nested Fields

When a field requires a separate query (like loading `orders` for each `User`), use `@SchemaMapping`:

```java
@Controller
public class UserFieldController {

    private final OrderService orderService;

    // Maps to User.orders — called for each User in the response
    @SchemaMapping(typeName = "User", field = "orders")
    public List<OrderDto> orders(UserDto user) {
        return orderService.findByUserId(user.id());
    }

    // Maps to User.orderCount
    @SchemaMapping(typeName = "User", field = "orderCount")
    public int orderCount(UserDto user) {
        return orderService.countByUserId(user.id());
    }
}
```

---

## The N+1 Problem and @BatchMapping

Without batching, loading orders for 10 users = 10 separate database queries (N+1 problem).

`@BatchMapping` solves this by loading all at once:

```java
@Controller
public class UserFieldController {

    private final OrderService orderService;

    // Called ONCE with ALL users in the response — returns a Map<User, List<Order>>
    @BatchMapping(typeName = "User", field = "orders")
    public Map<UserDto, List<OrderDto>> orders(List<UserDto> users) {
        List<Long> userIds = users.stream().map(UserDto::id).toList();
        Map<Long, List<OrderDto>> ordersByUser = orderService.findByUserIds(userIds);

        return users.stream().collect(Collectors.toMap(
            user -> user,
            user -> ordersByUser.getOrDefault(user.id(), List.of())
        ));
    }
}
```

This turns N+1 queries into a single query.

---

## Subscriptions (Real-Time Updates)

```xml
<!-- Additional dependency for WebSocket subscriptions -->
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-websocket</artifactId>
</dependency>
```

```java
@Controller
public class OrderSubscriptionController {

    private final Sinks.Many<OrderDto> orderSink =
        Sinks.many().multicast().onBackpressureBuffer();

    // Maps to Subscription.orderStatusChanged
    @SubscriptionMapping
    public Flux<OrderDto> orderStatusChanged(@Argument Long orderId) {
        return orderSink.asFlux()
            .filter(order -> order.id().equals(orderId));
    }

    // Called when order status changes
    public void publishOrderUpdate(OrderDto order) {
        orderSink.tryEmitNext(order);
    }
}
```

---

## Configuration

```yaml
spring:
  graphql:
    graphiql:
      enabled: true        # enables GraphiQL browser UI at /graphiql
      path: /graphiql
    schema:
      printer:
        enabled: true
    websocket:
      path: /graphql-ws    # for subscriptions
  mvc:
    async:
      request-timeout: 60s
```

---

## Error Handling

```java
@Component
public class GraphQLExceptionHandler implements DataFetcherExceptionResolver {

    @Override
    public Mono<List<GraphQLError>> resolveException(
            Throwable exception,
            DataFetchingEnvironment environment) {

        if (exception instanceof ResourceNotFoundException ex) {
            return Mono.just(List.of(GraphqlErrorBuilder.newError(environment)
                .message(ex.getMessage())
                .errorType(ErrorType.NOT_FOUND)
                .build()));
        }

        if (exception instanceof AccessDeniedException) {
            return Mono.just(List.of(GraphqlErrorBuilder.newError(environment)
                .message("Access denied")
                .errorType(ErrorType.FORBIDDEN)
                .build()));
        }

        return Mono.empty();   // unhandled — Spring returns default error
    }
}
```

---

## Security with GraphQL

```java
@Controller
public class OrderQueryController {

    @QueryMapping
    @PreAuthorize("isAuthenticated()")
    public List<OrderDto> myOrders(@AuthenticationPrincipal UserDetails user) {
        return orderService.findByUserEmail(user.getUsername());
    }

    @MutationMapping
    @PreAuthorize("hasRole('ADMIN')")
    public boolean cancelOrder(@Argument Long id) {
        return orderService.cancel(id);
    }
}
```

---

## Testing GraphQL

```java
@SpringBootTest
@AutoConfigureHttpGraphQlTester
class UserGraphQLTest {

    @Autowired
    private HttpGraphQlTester graphQlTester;

    @Test
    void query_user_shouldReturnUser() {
        graphQlTester.document("""
            query {
                user(id: 1) {
                    id
                    name
                    email
                }
            }
        """)
        .execute()
        .path("user.name").entity(String.class).isEqualTo("Alice")
        .path("user.email").entity(String.class).isEqualTo("alice@example.com");
    }

    @Test
    void mutation_createUser_shouldReturnCreatedUser() {
        graphQlTester.document("""
            mutation {
                createUser(input: { name: "Bob", email: "bob@example.com", password: "pass1234" }) {
                    id
                    name
                }
            }
        """)
        .execute()
        .path("createUser.name").entity(String.class).isEqualTo("Bob");
    }
}
```

---

## REST vs GraphQL — When to Use

| Use REST when | Use GraphQL when |
|---|---|
| Simple CRUD resources | Complex, nested data relationships |
| Public API for 3rd parties | Multiple clients (mobile, web) with different data needs |
| File uploads | Rapid UI development (client controls data shape) |
| HTTP caching is important | Reducing network round trips is critical |
| Team is already familiar with REST | Real-time subscriptions are needed |
