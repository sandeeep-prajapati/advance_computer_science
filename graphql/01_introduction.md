# Chapter 1: Introduction to GraphQL

## What is GraphQL?

**GraphQL** is a query language for APIs and a runtime for executing those queries. Developed by Facebook in 2012 and open-sourced in 2015, it gives clients the power to ask for exactly the data they need — nothing more, nothing less.

> **Analogy:** REST is like a restaurant with a fixed menu — you order "dish 5" and get everything that comes with it. GraphQL is like a build-your-own meal where you specify exactly which ingredients you want.

---

## The Problem GraphQL Solves

### Over-fetching (REST):
```
GET /users/1
Response: { id, name, email, phone, address, bio, avatar, preferences, ... }
You only needed: { id, name }
→ Wasted bandwidth and parsing
```

### Under-fetching (REST):
```
GET /users/1           → get user
GET /users/1/posts     → get their posts
GET /users/1/followers → get followers
→ 3 round trips for one screen
```

### GraphQL solves both:
```graphql
query {
  user(id: 1) {
    id
    name
    posts { title }
    followers { count }
  }
}
→ 1 request, exactly the fields you need
```

---

## GraphQL vs REST vs gRPC

| Feature | REST | GraphQL | gRPC |
|---|---|---|---|
| Protocol | HTTP | HTTP | HTTP/2 |
| Data format | JSON/XML | JSON | Protocol Buffers |
| Schema | OpenAPI (optional) | Strongly typed (required) | .proto files |
| Versioning | URL versions (`/v2/`) | Schema evolution | .proto versioning |
| Over-fetching | Common | Eliminated | N/A |
| Real-time | SSE/Webhooks | Subscriptions (built-in) | Server streaming |
| Learning curve | Low | Medium | Medium-High |
| Browser support | Native | Native | Needs proxy (gRPC-Web) |
| Best for | Simple CRUD, public APIs | Complex UI data, mobile | Internal microservices |

---

## Core Concepts

### 1. Schema
The schema is the contract between client and server — it defines what data exists and what operations are possible.

```graphql
type Query {
  user(id: ID!): User
  posts: [Post!]!
}

type User {
  id: ID!
  name: String!
  email: String!
  posts: [Post!]!
}

type Post {
  id: ID!
  title: String!
  body: String!
  author: User!
}
```

### 2. Queries (Read)
```graphql
query GetUser {
  user(id: "1") {
    name
    posts {
      title
    }
  }
}
```

### 3. Mutations (Write)
```graphql
mutation CreatePost {
  createPost(input: { title: "Hello", body: "World" }) {
    id
    title
  }
}
```

### 4. Subscriptions (Real-time)
```graphql
subscription {
  postAdded {
    id
    title
    author { name }
  }
}
```

### 5. Resolvers
Resolvers are functions that return data for each field in the schema. They connect the schema to your data sources.

```javascript
const resolvers = {
  Query: {
    user: (parent, args, context) => context.db.users.findById(args.id),
  },
  User: {
    posts: (user, args, context) => context.db.posts.findByAuthor(user.id),
  },
};
```

---

## GraphQL Execution Model

```
Client sends query
    ↓
Parse query string → Abstract Syntax Tree (AST)
    ↓
Validate AST against schema
    ↓
Execute: walk the AST, call resolvers
    ├── Query.user resolver → returns user object
    └── User.posts resolver → returns posts array
         └── Post.title → returns string (no resolver needed)
    ↓
Collect all results
    ↓
Return JSON response
```

---

## The Single Endpoint

Unlike REST (many URLs), GraphQL has **one endpoint**:

```
REST:
  GET    /users
  POST   /users
  GET    /users/:id
  PUT    /users/:id
  DELETE /users/:id
  GET    /users/:id/posts
  ...

GraphQL:
  POST /graphql    ← everything goes here
```

The **operation type** (query/mutation/subscription) and **fields** determine what happens.

---

## When to Use GraphQL

**GraphQL shines when:**
- Multiple clients (web, mobile, TV) need different data shapes
- UI screens require data from multiple resources
- Rapid frontend development — frontend can evolve independently
- API is consumed by many different teams
- You need built-in introspection and self-documentation

**Stick with REST when:**
- Simple CRUD with uniform responses
- Public API consumed by third parties (REST is more familiar)
- File uploads are the primary operation
- HTTP caching is critical (REST caches better by default)
- Team has no GraphQL experience and simple API suffices

---

## GraphQL Request Format

GraphQL runs over HTTP, typically as a POST:

```http
POST /graphql
Content-Type: application/json

{
  "query": "query GetUser($id: ID!) { user(id: $id) { name email } }",
  "variables": { "id": "42" },
  "operationName": "GetUser"
}
```

Response:
```json
{
  "data": {
    "user": {
      "name": "Alice",
      "email": "alice@example.com"
    }
  },
  "errors": null
}
```

---

## Introspection

GraphQL schemas are self-documenting. Clients can query the schema itself:

```graphql
query {
  __schema {
    types {
      name
      fields { name type { name } }
    }
  }
}

query {
  __type(name: "User") {
    fields {
      name
      type { name kind }
      description
    }
  }
}
```

This powers tools like **GraphiQL**, **Apollo Studio**, and **Insomnia**.

---

## Key Terminology

| Term | Meaning |
|---|---|
| **Schema** | Type definitions and available operations |
| **Type** | A named set of fields (like a class) |
| **Field** | A named value on a type |
| **Scalar** | Leaf types: String, Int, Float, Boolean, ID |
| **Resolver** | Function that returns a field's value |
| **Query** | Read operation |
| **Mutation** | Write operation |
| **Subscription** | Real-time stream |
| **Fragment** | Reusable set of fields |
| **Variable** | Runtime parameter for a query |
| **Directive** | Annotation that affects execution (`@include`, `@skip`) |
| **Context** | Shared data across all resolvers (auth, DB connection) |
| **DataLoader** | Batching utility to solve N+1 problem |
| **Introspection** | Querying the schema itself |
