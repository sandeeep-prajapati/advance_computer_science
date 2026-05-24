# GraphQL Developer Book

A comprehensive, chapter-wise guide to mastering GraphQL — from first query to production-grade APIs, federation, real-time subscriptions, and performance optimization.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_introduction.md) | **Introduction to GraphQL** | What is GraphQL, REST vs GraphQL, core concepts, when to use |
| [02](02_schema_definition_language.md) | **Schema Definition Language** | Types, fields, scalars, enums, interfaces, unions, input types |
| [03](03_queries.md) | **Queries** | Basic queries, nested fields, arguments, aliases, field selection |
| [04](04_mutations.md) | **Mutations** | Create/update/delete, input types, returning data, error handling |
| [05](05_subscriptions.md) | **Subscriptions** | Real-time data, WebSockets, pub/sub, filtering, lifecycle |
| [06](06_resolvers.md) | **Resolvers** | Resolver chain, context, info, default resolvers, resolver map |
| [07](07_type_system.md) | **Type System** | Scalars, objects, enums, interfaces, unions, custom scalars, nullability |
| [08](08_variables_and_fragments.md) | **Variables & Fragments** | Query variables, inline fragments, named fragments, directives |
| [09](09_directives.md) | **Directives** | @include, @skip, @deprecated, @specifiedBy, custom directives |
| [10](10_authentication_and_authorization.md) | **Authentication & Authorization** | JWT in context, field-level auth, schema directives, Shield |
| [11](11_n_plus_one_and_dataloader.md) | **N+1 Problem & DataLoader** | The N+1 problem, DataLoader batching, caching, per-request loaders |
| [12](12_error_handling.md) | **Error Handling** | Partial errors, error extensions, custom errors, error masking |
| [13](13_pagination.md) | **Pagination** | Offset, cursor-based (Relay spec), connection pattern, pageInfo |
| [14](14_file_uploads.md) | **File Uploads** | Multipart spec, upload scalar, handling on server and client |
| [15](15_graphql_with_react.md) | **GraphQL with React (Apollo Client)** | Apollo setup, useQuery, useMutation, useSubscription, caching |
| [16](16_server_implementations.md) | **Server Implementations** | Apollo Server, GraphQL Yoga, Pothos, NestJS, Spring GraphQL |
| [17](17_testing.md) | **Testing GraphQL APIs** | Unit testing resolvers, integration tests, supertest, mocking |
| [18](18_schema_design.md) | **Schema Design Best Practices** | Naming, nullability, versioning, evolution, domain modeling |
| [19](19_federation.md) | **Federation & Microservices** | Apollo Federation, subgraphs, @key, @external, supergraph |
| [20](20_learning_roadmap.md) | **Learning Roadmap** | 6-week hands-on plan from first query to production GraphQL API |

---

## How to Read This Book

**If you're brand new to GraphQL:** Read 1 → 2 → 3 → 4 → 6 in order. This gives you the core mental model.

**If you're building a client app:** Read 3 → 8 → 15 → 13.

**If you're building a server:** Read 2 → 6 → 16 → 10 → 11 → 12.

**If you have performance problems:** Read 11 → 13 → 19.

**If you're going real-time:** Read 5 → 15.

**If you're building microservices:** Read 19 → 6 → 10.

---

## Key Concepts at a Glance

```
GraphQL API
├── Schema (Contract between client and server)
│   ├── Root Types
│   │   ├── Query         (read operations)
│   │   ├── Mutation      (write operations)
│   │   └── Subscription  (real-time streams)
│   └── Type Definitions
│       ├── Object Types  (User, Post, Order)
│       ├── Input Types   (CreateUserInput)
│       ├── Enum Types    (UserStatus, Role)
│       ├── Interface     (Node, Timestamped)
│       ├── Union         (SearchResult = User | Post)
│       └── Scalars       (String, Int, Boolean, ID, custom)
│
├── Resolvers (Implementation)
│   ├── Query resolvers   → fetch data
│   ├── Mutation resolvers → mutate data
│   ├── Field resolvers   → nested field loading
│   └── Subscription resolvers → async iterables
│
├── Execution Engine
│   ├── Parse query → AST
│   ├── Validate against schema
│   ├── Execute resolvers (breadth-first)
│   └── Return JSON response
│
└── Client
    ├── Query document     → sent as POST /graphql
    ├── Variables          → separate from query string
    ├── Fragments          → reusable field sets
    └── Directives         → @include, @skip, @defer
```
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
# Chapter 2: Schema Definition Language (SDL)

## What is the SDL?

The **Schema Definition Language** is the language used to define a GraphQL schema. It describes the types in your API, the fields they contain, and what operations are possible. The SDL is language-agnostic — the same schema works whether your server is Node.js, Java, Python, or Go.

---

## Scalar Types (Leaf Values)

Scalars are the primitive types — they resolve to concrete values with no sub-fields.

### Built-in Scalars:
```graphql
type Example {
  id: ID!           # unique identifier (serialized as String)
  name: String!     # UTF-8 string
  age: Int!         # 32-bit signed integer
  price: Float!     # double-precision floating point
  active: Boolean!  # true or false
}
```

### Custom Scalars:
```graphql
scalar Date        # ISO 8601 date string
scalar DateTime    # ISO 8601 datetime string
scalar Email       # validated email format
scalar URL         # validated URL
scalar JSON        # arbitrary JSON value
scalar Upload      # file upload

type User {
  id: ID!
  createdAt: DateTime!
  birthDate: Date
  website: URL
}
```

---

## Object Types

Object types define the shape of data objects:

```graphql
type User {
  id: ID!
  name: String!
  email: String!
  bio: String          # nullable (no !)
  role: UserRole!
  posts: [Post!]!      # non-null list of non-null Posts
  createdAt: DateTime!
  updatedAt: DateTime!
}

type Post {
  id: ID!
  title: String!
  body: String!
  published: Boolean!
  author: User!        # relation — resolved by a resolver
  tags: [String!]!
  likeCount: Int!
}
```

### Nullability Rules:
```graphql
String      # nullable — can be null
String!     # non-null — never null
[String]    # nullable list of nullable strings
[String!]   # nullable list of non-null strings
[String]!   # non-null list of nullable strings
[String!]!  # non-null list of non-null strings
```

---

## Root Types

The three root types define entry points for operations:

```graphql
type Query {
  # Read operations
  user(id: ID!): User
  users(filter: UserFilter, page: Int, limit: Int): UserConnection!
  post(id: ID!): Post
  posts(authorId: ID): [Post!]!
  me: User               # current authenticated user
}

type Mutation {
  # Write operations
  createUser(input: CreateUserInput!): CreateUserPayload!
  updateUser(id: ID!, input: UpdateUserInput!): UpdateUserPayload!
  deleteUser(id: ID!): DeleteUserPayload!
  createPost(input: CreatePostInput!): Post!
  publishPost(id: ID!): Post!
}

type Subscription {
  # Real-time streams
  postAdded: Post!
  postUpdated(id: ID!): Post!
  userTyping(channelId: ID!): TypingEvent!
}
```

---

## Input Types

Input types are used for mutation arguments. They cannot be used as output types.

```graphql
input CreateUserInput {
  name: String!
  email: String!
  password: String!
  role: UserRole = USER        # default value
}

input UpdateUserInput {
  name: String                 # all fields optional for updates
  email: String
  bio: String
}

input UserFilter {
  role: UserRole
  search: String
  createdAfter: DateTime
}

# Usage in mutation
type Mutation {
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User!
  listUsers(filter: UserFilter): [User!]!
}
```

---

## Enum Types

```graphql
enum UserRole {
  ADMIN
  MODERATOR
  USER
  GUEST
}

enum PostStatus {
  DRAFT
  PUBLISHED
  ARCHIVED
}

enum SortOrder {
  ASC
  DESC
}

# Usage
type User {
  role: UserRole!
}

type Query {
  users(sortBy: String, order: SortOrder = DESC): [User!]!
}
```

---

## Interfaces

Interfaces define a contract — any type implementing it must include those fields:

```graphql
interface Node {
  id: ID!
}

interface Timestamped {
  createdAt: DateTime!
  updatedAt: DateTime!
}

type User implements Node & Timestamped {
  id: ID!
  name: String!
  email: String!
  createdAt: DateTime!
  updatedAt: DateTime!
}

type Post implements Node & Timestamped {
  id: ID!
  title: String!
  createdAt: DateTime!
  updatedAt: DateTime!
}

# Query any Node by ID
type Query {
  node(id: ID!): Node
}
```

---

## Union Types

Unions represent a field that can be one of several types (no shared fields required):

```graphql
union SearchResult = User | Post | Comment | Product

type Query {
  search(query: String!): [SearchResult!]!
}
```

Client uses inline fragments to handle each type:
```graphql
query Search {
  search(query: "Alice") {
    ... on User { id name email }
    ... on Post { id title author { name } }
    ... on Comment { id body }
    ... on Product { id name price }
  }
}
```

### Interface vs Union:
- **Interface** — types share common fields (use `... on` for unique fields)
- **Union** — types have NO common fields (must always use `... on`)

---

## Field Arguments

Any field can accept arguments:

```graphql
type Query {
  users(
    limit: Int = 20
    offset: Int = 0
    search: String
    role: UserRole
    sortBy: String = "createdAt"
    order: SortOrder = DESC
  ): [User!]!
}

type User {
  posts(
    status: PostStatus = PUBLISHED
    limit: Int = 10
  ): [Post!]!

  # Computed field with argument
  avatar(size: Int = 64): String
}
```

---

## Payload Types (Best Practice for Mutations)

Wrapping mutation results in a payload type allows returning either success data or errors:

```graphql
type CreateUserPayload {
  user: User
  errors: [UserError!]!
  success: Boolean!
}

type UserError {
  field: String      # which field caused the error (null for global errors)
  message: String!
}

type Mutation {
  createUser(input: CreateUserInput!): CreateUserPayload!
}
```

Client:
```graphql
mutation {
  createUser(input: { name: "Alice", email: "alice@example.com", password: "pass123" }) {
    success
    user { id name }
    errors { field message }
  }
}
```

---

## Descriptions (Documentation)

```graphql
"""
Represents a user account in the system.
"""
type User {
  "The unique identifier of the user."
  id: ID!

  "The user's display name."
  name: String!

  """
  The user's email address. Must be unique across all users.
  Used for authentication and notifications.
  """
  email: String!
}

"""
Input for creating a new user account.
"""
input CreateUserInput {
  "Must be between 2 and 100 characters."
  name: String!

  "Must be a valid email address."
  email: String!
}
```

Descriptions appear in GraphiQL, Apollo Studio, and generated documentation.

---

## Deprecation

```graphql
type User {
  id: ID!
  name: String!

  username: String @deprecated(reason: "Use `name` instead. Removing in v3.")
  phone: String @deprecated
}

type Query {
  users: [User!]!
  allUsers: [User!]! @deprecated(reason: "Use `users` instead.")
}
```

---

## Complete Schema Example

```graphql
scalar DateTime
scalar Upload

interface Node { id: ID! }

enum UserRole { ADMIN USER }
enum PostStatus { DRAFT PUBLISHED }

type User implements Node {
  id: ID!
  name: String!
  email: String!
  role: UserRole!
  posts(status: PostStatus): [Post!]!
  createdAt: DateTime!
}

type Post implements Node {
  id: ID!
  title: String!
  body: String!
  status: PostStatus!
  author: User!
  tags: [String!]!
  createdAt: DateTime!
}

input CreatePostInput {
  title: String!
  body: String!
  tags: [String!]! = []
}

type PostPayload {
  post: Post
  errors: [String!]!
}

type Query {
  me: User
  user(id: ID!): User
  posts(authorId: ID, status: PostStatus): [Post!]!
}

type Mutation {
  createPost(input: CreatePostInput!): PostPayload!
  publishPost(id: ID!): Post!
  deletePost(id: ID!): Boolean!
}

type Subscription {
  postPublished: Post!
}
```
# Chapter 3: Queries

## Basic Query Structure

A GraphQL query asks for specific fields on objects. Only the fields you request are returned.

```graphql
# Simplest query
query {
  me {
    id
    name
    email
  }
}

# Response
{
  "data": {
    "me": {
      "id": "1",
      "name": "Alice",
      "email": "alice@example.com"
    }
  }
}
```

---

## Query with Arguments

Arguments let you pass parameters to fields:

```graphql
query {
  user(id: "42") {
    name
    email
  }
}

query {
  posts(limit: 5, status: PUBLISHED) {
    id
    title
    createdAt
  }
}

query {
  users(search: "alice", role: ADMIN) {
    id
    name
  }
}
```

---

## Nested Fields

GraphQL excels at fetching related data in a single request:

```graphql
query {
  user(id: "1") {
    id
    name
    posts {
      id
      title
      tags
      comments {
        id
        body
        author {
          name
        }
      }
    }
    followers {
      id
      name
    }
  }
}
```

---

## Operation Name

Giving queries a name improves debugging and tooling:

```graphql
# Anonymous (valid but not recommended in production)
query {
  me { name }
}

# Named (recommended — shows up in logs, Apollo Studio, etc.)
query GetCurrentUser {
  me {
    id
    name
    email
    role
  }
}

query GetUserProfile {
  user(id: "42") {
    name
    bio
    avatar
  }
}
```

---

## Aliases

When you need the same field multiple times with different arguments:

```graphql
query GetUsers {
  alice: user(id: "1") {
    name
    email
  }
  bob: user(id: "2") {
    name
    email
  }
}

# Response
{
  "data": {
    "alice": { "name": "Alice", "email": "alice@example.com" },
    "bob": { "name": "Bob", "email": "bob@example.com" }
  }
}

# Also useful for renaming fields in response
query GetPosts {
  recentPosts: posts(limit: 5, order: DESC) {
    id
    heading: title     # rename 'title' to 'heading' in response
    preview: excerpt
  }
}
```

---

## Multiple Root Fields

A single query can request multiple root-level fields:

```graphql
query DashboardData {
  me {
    id
    name
  }
  recentPosts: posts(limit: 10) {
    id
    title
    createdAt
  }
  notifications(unreadOnly: true) {
    id
    message
  }
  siteStats {
    totalUsers
    totalPosts
  }
}
```

All these are resolved in parallel (depending on server configuration).

---

## Variables

Variables separate the query structure from runtime values — essential for dynamic queries:

```graphql
# Query with variable declarations
query GetUser($id: ID!, $includePrivate: Boolean = false) {
  user(id: $id) {
    id
    name
    email @include(if: $includePrivate)
    phone @include(if: $includePrivate)
  }
}
```

Variables are sent separately (not interpolated into the query string):

```json
{
  "query": "query GetUser($id: ID!, $includePrivate: Boolean = false) { ... }",
  "variables": {
    "id": "42",
    "includePrivate": true
  }
}
```

Variable syntax:
```graphql
query MyQuery(
  $id: ID!                    # required
  $limit: Int = 10            # optional with default
  $status: PostStatus         # optional, nullable
  $tags: [String!]!           # required, non-null list
) { ... }
```

---

## Fragments

Fragments are reusable field sets — avoid repeating the same fields in multiple queries:

```graphql
# Define a fragment
fragment UserBasic on User {
  id
  name
  email
  avatar
}

fragment PostPreview on Post {
  id
  title
  excerpt
  publishedAt
  author {
    ...UserBasic
  }
}

# Use in queries
query GetHomepage {
  featuredPosts: posts(featured: true, limit: 3) {
    ...PostPreview
  }
  recentPosts: posts(limit: 10) {
    ...PostPreview
  }
  me {
    ...UserBasic
    unreadNotificationsCount
  }
}
```

---

## Inline Fragments (on Union/Interface)

Required when a field can be multiple types:

```graphql
query SearchResults {
  search(query: "postgres") {
    __typename    # returns the type name: "Post", "User", etc.

    # Fields common to all (if using interface)
    ... on Node {
      id
    }

    ... on Post {
      title
      excerpt
      author { name }
    }

    ... on User {
      name
      bio
    }

    ... on Product {
      name
      price
    }
  }
}
```

---

## __typename — Meta Field

`__typename` is available on any type and returns the type name as a string:

```graphql
query {
  user(id: "1") {
    __typename   # "User"
    id
    name
  }
  posts {
    __typename   # "Post"
    id
    title
  }
}
```

Useful for:
- Distinguishing union/interface types
- Cache key generation in Apollo Client
- Debugging resolver issues

---

## @include and @skip Directives

Conditionally include or exclude fields:

```graphql
query GetUser($withPosts: Boolean!, $hideSensitive: Boolean!) {
  user(id: "1") {
    id
    name
    email @skip(if: $hideSensitive)
    phone @skip(if: $hideSensitive)

    posts @include(if: $withPosts) {
      id
      title
    }
  }
}
```

Variables:
```json
{
  "withPosts": true,
  "hideSensitive": false
}
```

---

## Pagination Queries

```graphql
# Offset-based
query GetUsers($offset: Int!, $limit: Int!) {
  users(offset: $offset, limit: $limit) {
    id
    name
  }
  usersCount
}

# Cursor-based (Relay-style)
query GetPosts($first: Int!, $after: String) {
  posts(first: $first, after: $after) {
    edges {
      cursor
      node {
        id
        title
        createdAt
      }
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
    totalCount
  }
}
```

---

## Query Best Practices

```graphql
# Bad: requesting everything
query {
  users {
    id name email phone bio avatar role
    posts { id title body tags status createdAt updatedAt }
    followers { id name email }
    following { id name email }
  }
}

# Good: request only what you render
query GetUserList {
  users(limit: 20) {
    id
    name
    avatar
    role
  }
}

# Good: separate queries for separate concerns
query GetUserProfile {
  user(id: $id) {
    id name bio avatar
  }
}

query GetUserPosts {
  posts(authorId: $id, limit: 10) {
    id title excerpt createdAt
  }
}
```

**Rules:**
- Name every query
- Use variables, never string interpolation
- Request only fields you render
- Use fragments for shared field sets
- Use pagination for list fields
# Chapter 4: Mutations

## What are Mutations?

Mutations are GraphQL operations that **modify** data on the server — create, update, delete. Unlike queries (which are typically safe to run in parallel), mutations run **sequentially** in the order listed.

```graphql
# Mutation keyword
mutation {
  createUser(input: { name: "Alice", email: "alice@example.com", password: "pass123" }) {
    id
    name
    email
  }
}
```

---

## Basic Mutation Structure

```graphql
# Schema definition
type Mutation {
  createPost(input: CreatePostInput!): Post!
  updatePost(id: ID!, input: UpdatePostInput!): Post!
  deletePost(id: ID!): Boolean!
  publishPost(id: ID!): Post!
}

input CreatePostInput {
  title: String!
  body: String!
  tags: [String!]! = []
}

input UpdatePostInput {
  title: String
  body: String
  tags: [String!]
}
```

Resolvers:
```javascript
const resolvers = {
  Mutation: {
    createPost: async (parent, { input }, context) => {
      if (!context.user) throw new AuthenticationError('Must be logged in');

      const post = await context.db.posts.create({
        ...input,
        authorId: context.user.id,
        status: 'DRAFT',
      });

      return post;
    },

    updatePost: async (parent, { id, input }, context) => {
      const post = await context.db.posts.findById(id);
      if (!post) throw new UserInputError('Post not found');
      if (post.authorId !== context.user.id) throw new ForbiddenError('Not your post');

      return context.db.posts.update(id, input);
    },

    deletePost: async (parent, { id }, context) => {
      await context.db.posts.delete(id);
      return true;
    },
  },
};
```

---

## Input Types for Mutations

Always use dedicated input types for mutations — don't reuse output types:

```graphql
# Good: dedicated input types
input CreateUserInput {
  name: String!
  email: String!
  password: String!
  role: UserRole = USER
}

input UpdateUserInput {
  name: String        # all fields optional for partial updates
  email: String
  bio: String
  avatarUrl: String
}

# Bad: passing fields directly (becomes unwieldy)
type Mutation {
  createUser(name: String!, email: String!, password: String!, role: UserRole): User!
}
```

---

## Mutation Return Types — Payload Pattern

Rather than returning the entity directly, return a payload type that can include both the result AND errors:

```graphql
type CreatePostPayload {
  post: Post              # null if mutation failed
  errors: [UserError!]!  # empty if success
}

type UserError {
  field: String          # which input field caused the error (null for global errors)
  message: String!
}

type Mutation {
  createPost(input: CreatePostInput!): CreatePostPayload!
  updatePost(id: ID!, input: UpdatePostInput!): UpdatePostPayload!
  deletePost(id: ID!): DeletePostPayload!
}
```

Client:
```graphql
mutation CreatePost($input: CreatePostInput!) {
  createPost(input: $input) {
    post {
      id
      title
    }
    errors {
      field
      message
    }
  }
}
```

This pattern avoids throwing HTTP errors for business-logic failures — the client always gets a 200 with errors inside the payload.

---

## Variables with Mutations

```graphql
mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    id
    name
    email
    createdAt
  }
}
```

Variables:
```json
{
  "input": {
    "name": "Alice",
    "email": "alice@example.com",
    "password": "SecurePass123!",
    "role": "USER"
  }
}
```

---

## Multiple Mutations in One Request

Mutations in a single request are executed **sequentially**, not in parallel:

```graphql
mutation BatchOperations {
  first: createPost(input: { title: "Post 1", body: "..." }) {
    id title
  }
  second: createPost(input: { title: "Post 2", body: "..." }) {
    id title
  }
  updateProfile: updateUser(id: "1", input: { bio: "New bio" }) {
    id bio
  }
}
```

`first` completes before `second` starts.

---

## Common Mutation Patterns

### Create:
```graphql
mutation CreateComment($input: CreateCommentInput!) {
  createComment(input: $input) {
    comment {
      id
      body
      author { name avatar }
      createdAt
    }
    errors { field message }
  }
}
```

### Update (Partial):
```graphql
mutation UpdateProfile($input: UpdateUserInput!) {
  updateUser(id: "me", input: $input) {
    user {
      id name bio email
    }
    errors { field message }
  }
}
```

### Delete:
```graphql
mutation DeletePost($id: ID!) {
  deletePost(id: $id) {
    success
    deletedId    # return ID so client can remove from cache
    errors { message }
  }
}
```

### Toggle (state change):
```graphql
mutation ToggleLike($postId: ID!) {
  toggleLike(postId: $postId) {
    post {
      id
      liked     # new state
      likeCount # updated count
    }
  }
}
```

---

## Optimistic Updates (Client Side)

```javascript
// Apollo Client — update UI before mutation completes
const [createPost] = useMutation(CREATE_POST, {
  optimisticResponse: {
    createPost: {
      post: {
        __typename: 'Post',
        id: 'temp-id',
        title: input.title,
        status: 'DRAFT',
        createdAt: new Date().toISOString(),
      },
      errors: [],
    },
  },
  update(cache, { data: { createPost } }) {
    // Update the Apollo cache with the new post
    cache.modify({
      fields: {
        posts(existingPosts = []) {
          const newPostRef = cache.writeFragment({
            data: createPost.post,
            fragment: POST_FRAGMENT,
          });
          return [newPostRef, ...existingPosts];
        },
      },
    });
  },
});
```

---

## Mutation Resolver Implementation Patterns

### With validation:
```javascript
createUser: async (parent, { input }, context) => {
  // Validate
  const errors = validateCreateUserInput(input);
  if (errors.length > 0) {
    return { user: null, errors };
  }

  // Check uniqueness
  const existing = await context.db.users.findByEmail(input.email);
  if (existing) {
    return {
      user: null,
      errors: [{ field: 'email', message: 'Email is already in use' }],
    };
  }

  // Hash password
  const passwordHash = await bcrypt.hash(input.password, 10);

  // Create
  const user = await context.db.users.create({
    ...input,
    password: passwordHash,
  });

  // Side effects
  await context.emailService.sendWelcome(user.email);
  context.pubsub.publish('USER_CREATED', { userCreated: user });

  return { user, errors: [] };
},
```

### With transactions:
```javascript
placeOrder: async (parent, { input }, context) => {
  return context.db.transaction(async (trx) => {
    const order = await context.db.orders.create(input, trx);
    await context.db.inventory.deduct(input.items, trx);
    await context.db.payments.charge(input.paymentMethod, order.total, trx);
    return { order, errors: [] };
  });
},
```

---

## Idempotency

For mutations that should be safe to retry:

```graphql
input CreateOrderInput {
  idempotencyKey: String!   # unique per attempt, same key = same result
  items: [OrderItemInput!]!
  paymentMethod: ID!
}
```

Server checks if the key was already used and returns the existing result.

---

## Mutation Best Practices

- Return the modified entity so clients can update their cache
- Use payload types for mutations — include both data and errors
- Validate input server-side — never trust client
- Use input types, not individual arguments for complex mutations
- Name mutations with verb + noun: `createUser`, `publishPost`, `cancelOrder`
- Publish events after successful mutations (for subscriptions)
- Use transactions for multi-step operations
- Return `deletedId` from delete mutations so clients can remove from cache
# Chapter 5: Subscriptions

## What are Subscriptions?

GraphQL Subscriptions enable **real-time data** from the server to the client. Unlike queries and mutations (request/response), subscriptions maintain a persistent connection and push data when events occur.

```
Client ──────── subscribe ──────────► Server
       ◄──── initial response ────────
       ◄──── event data ─────────────  (when something happens)
       ◄──── event data ─────────────  (again...)
       ◄──── event data ─────────────  (and again...)
       ──────── unsubscribe ──────────►
```

---

## Transport: WebSockets

Subscriptions typically use WebSockets. The `graphql-ws` protocol (modern) has replaced the older `subscriptions-transport-ws`:

```
Client → Server: { type: "connection_init" }
Server → Client: { type: "connection_ack" }
Client → Server: { type: "subscribe", id: "1", payload: { query: "subscription {...}" } }
Server → Client: { type: "next", id: "1", payload: { data: {...} } }
Server → Client: { type: "next", id: "1", payload: { data: {...} } }
Client → Server: { type: "complete", id: "1" }  (unsubscribe)
```

---

## Schema Definition

```graphql
type Subscription {
  # Post events
  postAdded: Post!
  postUpdated(id: ID!): Post!
  postDeleted: ID!                # returns deleted post ID

  # Chat
  messageSent(channelId: ID!): Message!
  userTyping(channelId: ID!): TypingEvent!
  userOnlineStatus(userId: ID!): UserStatus!

  # Notifications
  notificationReceived: Notification!

  # Order tracking
  orderStatusChanged(orderId: ID!): Order!
}

type TypingEvent {
  user: User!
  channelId: ID!
  isTyping: Boolean!
}

type UserStatus {
  user: User!
  online: Boolean!
  lastSeen: DateTime
}
```

---

## Server Implementation (Apollo Server + graphql-ws)

### Pub/Sub Setup:
```javascript
import { PubSub } from 'graphql-subscriptions';

// Simple in-memory pub/sub (use Redis for production multi-instance)
const pubsub = new PubSub();

// Event names as constants
const Events = {
  POST_ADDED: 'POST_ADDED',
  POST_UPDATED: 'POST_UPDATED',
  MESSAGE_SENT: 'MESSAGE_SENT',
  ORDER_STATUS_CHANGED: 'ORDER_STATUS_CHANGED',
};
```

### Subscription Resolvers:
```javascript
const resolvers = {
  Subscription: {
    postAdded: {
      subscribe: () => pubsub.asyncIterator([Events.POST_ADDED]),
      // resolver (optional — default passes through payload)
      resolve: (payload) => payload.postAdded,
    },

    postUpdated: {
      subscribe: (parent, { id }) => {
        // Only receive events for the specific post
        const topic = `${Events.POST_UPDATED}_${id}`;
        return pubsub.asyncIterator([topic]);
      },
    },

    messageSent: {
      subscribe: withFilter(
        () => pubsub.asyncIterator([Events.MESSAGE_SENT]),
        (payload, variables) => {
          // Only deliver if the message is in the requested channel
          return payload.messageSent.channelId === variables.channelId;
        }
      ),
    },
  },

  Mutation: {
    createPost: async (parent, { input }, context) => {
      const post = await context.db.posts.create(input);

      // Publish the event — all subscribers of POST_ADDED receive this
      await pubsub.publish(Events.POST_ADDED, { postAdded: post });

      return post;
    },

    sendMessage: async (parent, { input }, context) => {
      const message = await context.db.messages.create(input);
      await pubsub.publish(Events.MESSAGE_SENT, { messageSent: message });
      return message;
    },
  },
};
```

### Server Setup with graphql-ws:
```javascript
import { createServer } from 'http';
import { WebSocketServer } from 'ws';
import { useServer } from 'graphql-ws/lib/use/ws';
import { ApolloServer } from '@apollo/server';
import { makeExecutableSchema } from '@graphql-tools/schema';

const schema = makeExecutableSchema({ typeDefs, resolvers });
const httpServer = createServer(app);

// WebSocket server for subscriptions
const wsServer = new WebSocketServer({
  server: httpServer,
  path: '/graphql',
});

const serverCleanup = useServer(
  {
    schema,
    context: async (ctx) => {
      // Authentication for WebSocket connections
      const token = ctx.connectionParams?.authorization;
      const user = token ? await verifyToken(token) : null;
      return { user, pubsub, db };
    },
    onConnect: async (ctx) => {
      console.log('Client connected');
      return true;
    },
    onDisconnect: (ctx, code, reason) => {
      console.log('Client disconnected');
    },
  },
  wsServer
);

// Apollo Server setup
const server = new ApolloServer({
  schema,
  plugins: [
    ApolloServerPluginDrainHttpServer({ httpServer }),
    {
      async serverWillStart() {
        return {
          async drainServer() {
            await serverCleanup.dispose();
          },
        };
      },
    },
  ],
});

await server.start();
httpServer.listen(4000);
```

---

## Redis Pub/Sub for Production (Multi-Instance)

In-memory PubSub only works with a single server instance. Use Redis for horizontal scaling:

```javascript
import { RedisPubSub } from 'graphql-redis-subscriptions';
import Redis from 'ioredis';

const pubsub = new RedisPubSub({
  publisher: new Redis({ host: 'redis', port: 6379 }),
  subscriber: new Redis({ host: 'redis', port: 6379 }),
});
```

---

## withFilter — Filtering Events

```javascript
import { withFilter } from 'graphql-subscriptions';

const resolvers = {
  Subscription: {
    orderStatusChanged: {
      subscribe: withFilter(
        () => pubsub.asyncIterator([Events.ORDER_STATUS_CHANGED]),
        (payload, variables, context) => {
          // Only deliver if:
          // 1. The order ID matches
          // 2. The user owns this order
          return (
            payload.orderStatusChanged.id === variables.orderId &&
            payload.orderStatusChanged.userId === context.user.id
          );
        }
      ),
    },

    notificationReceived: {
      subscribe: withFilter(
        () => pubsub.asyncIterator([Events.NOTIFICATION]),
        (payload, variables, context) => {
          // Only deliver to the intended user
          return payload.notificationReceived.userId === context.user.id;
        }
      ),
    },
  },
};
```

---

## Client: Apollo Client Subscriptions

```javascript
import { ApolloClient, InMemoryCache, split, HttpLink } from '@apollo/client';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { createClient } from 'graphql-ws';
import { getMainDefinition } from '@apollo/client/utilities';

// HTTP link for queries/mutations
const httpLink = new HttpLink({ uri: 'http://localhost:4000/graphql' });

// WebSocket link for subscriptions
const wsLink = new GraphQLWsLink(
  createClient({
    url: 'ws://localhost:4000/graphql',
    connectionParams: {
      authorization: `Bearer ${localStorage.getItem('token')}`,
    },
  })
);

// Route to correct link based on operation type
const splitLink = split(
  ({ query }) => {
    const def = getMainDefinition(query);
    return def.kind === 'OperationDefinition' && def.operation === 'subscription';
  },
  wsLink,
  httpLink
);

const client = new ApolloClient({ link: splitLink, cache: new InMemoryCache() });
```

### Using useSubscription in React:
```javascript
import { useSubscription, gql } from '@apollo/client';

const POST_ADDED = gql`
  subscription OnPostAdded {
    postAdded {
      id
      title
      author { name avatar }
      createdAt
    }
  }
`;

function LiveFeed() {
  const [posts, setPosts] = useState([]);

  useSubscription(POST_ADDED, {
    onData: ({ data }) => {
      setPosts(prev => [data.data.postAdded, ...prev]);
    },
    onError: (error) => console.error('Subscription error:', error),
  });

  return (
    <ul>
      {posts.map(post => (
        <li key={post.id}>{post.title} by {post.author.name}</li>
      ))}
    </ul>
  );
}
```

---

## Subscription Lifecycle

```javascript
const resolvers = {
  Subscription: {
    userTyping: {
      subscribe: async function*(parent, { channelId }, context) {
        // Setup: register user as watching this channel
        await context.db.channelWatchers.add(channelId, context.user.id);

        try {
          const iterator = pubsub.asyncIterator([`TYPING_${channelId}`]);
          for await (const event of iterator) {
            yield event;
          }
        } finally {
          // Cleanup: runs when client disconnects
          await context.db.channelWatchers.remove(channelId, context.user.id);
        }
      },
    },
  },
};
```

---

## Authentication for Subscriptions

```javascript
useServer({
  schema,
  context: async (ctx, msg, args) => {
    // ctx.connectionParams contains the initial connection params
    const token = ctx.connectionParams?.authorization?.replace('Bearer ', '');

    if (!token) {
      throw new Error('Unauthorized: no token provided');
    }

    try {
      const user = jwt.verify(token, process.env.JWT_SECRET);
      return { user, pubsub, db };
    } catch {
      throw new Error('Unauthorized: invalid token');
    }
  },
}, wsServer);
```

---

## When to Use Subscriptions

| Use Case | Use | Alternative |
|---|---|---|
| Chat messages | Subscriptions | — |
| Live notifications | Subscriptions | SSE |
| Order tracking | Subscriptions | Polling |
| Dashboard metrics | Polling or SSE | Subscriptions |
| Collaborative editing | Subscriptions + CRDT | — |
| Real-time scores | Subscriptions | SSE |
| Periodic reports | Polling | Subscriptions (overkill) |

**SSE (Server-Sent Events)** is simpler for one-way server→client streams and works over HTTP/1.1 with no extra dependencies.
# Chapter 6: Resolvers

## What is a Resolver?

A **resolver** is a function that returns the value for a field in your GraphQL schema. Every field in GraphQL has a resolver — if you don't provide one, a default resolver returns the field from the parent object.

```
Schema field → Resolver function → Data (from DB, API, cache, etc.)
```

---

## Resolver Signature

Every resolver receives four arguments:

```javascript
fieldName(parent, args, context, info) {
  // parent  — result from parent resolver (null for root types)
  // args    — arguments passed to the field in the query
  // context — shared object across all resolvers (db, auth, dataloaders)
  // info    — AST info about the field being resolved (rarely used)
}
```

---

## Resolver Map

Resolvers are organized in a map matching the schema structure:

```javascript
const resolvers = {
  // Root Query type
  Query: {
    user: async (parent, { id }, context) => {
      return context.db.users.findById(id);
    },

    users: async (parent, { limit = 20, offset = 0 }, context) => {
      return context.db.users.findAll({ limit, offset });
    },

    me: async (parent, args, context) => {
      if (!context.user) return null;
      return context.db.users.findById(context.user.id);
    },
  },

  // Root Mutation type
  Mutation: {
    createUser: async (parent, { input }, context) => {
      return context.db.users.create(input);
    },
  },

  // Object type resolvers (field-level)
  User: {
    // posts is NOT in the DB user row — this resolver loads it
    posts: async (user, { limit = 10 }, context) => {
      return context.db.posts.findByAuthorId(user.id, { limit });
    },

    // fullName is a computed field
    fullName: (user) => `${user.firstName} ${user.lastName}`,

    // avatar with a fallback
    avatar: (user) => user.avatarUrl || `https://avatars.dicebear.com/api/identicon/${user.id}.svg`,
  },

  Post: {
    // author is stored as authorId in DB — resolver loads the User object
    author: async (post, args, context) => {
      return context.db.users.findById(post.authorId);
    },
  },
};
```

---

## Default Resolvers

If you don't define a resolver for a field, GraphQL uses the **default resolver**:

```javascript
// This is what the default resolver does:
(parent, args, context, info) => {
  return parent[info.fieldName];
}
```

So `user.name` works automatically if the DB object has a `name` field — no resolver needed.

---

## Context — Sharing State Across Resolvers

Context is created once per request and passed to every resolver:

```javascript
// Apollo Server
const server = new ApolloServer({
  typeDefs,
  resolvers,
  context: async ({ req }) => {
    // Extract token from header
    const token = req.headers.authorization?.replace('Bearer ', '');

    // Verify and decode
    let user = null;
    if (token) {
      try {
        user = jwt.verify(token, process.env.JWT_SECRET);
      } catch {
        // invalid token — user stays null
      }
    }

    return {
      user,                                    // current authenticated user
      db,                                      // database connection
      dataloaders: createDataLoaders(db),      // per-request DataLoader instances
      redis,                                   // cache
      logger,                                  // shared logger
    };
  },
});
```

---

## Resolver Chain (Execution Flow)

```graphql
query {
  user(id: "1") {      # Query.user resolver called
    name               # User.name (default resolver)
    posts {            # User.posts resolver called with user as parent
      title            # Post.title (default resolver)
      author {         # Post.author resolver called with post as parent
        name           # User.name (default resolver)
      }
    }
  }
}
```

Execution:
1. `Query.user(null, { id: "1" }, context)` → returns `{ id: 1, name: "Alice", ... }`
2. `User.posts({ id: 1, name: "Alice", ... }, {}, context)` → returns `[{ id: 10, title: "..." }]`
3. `Post.author({ id: 10, title: "...", authorId: 1 }, {}, context)` → returns `{ id: 1, name: "Alice" }`

---

## Async Resolvers

Most resolvers are async — they fetch data from databases or APIs:

```javascript
const resolvers = {
  Query: {
    // Can return a Promise
    user: (parent, { id }, context) => {
      return context.db.query('SELECT * FROM users WHERE id = $1', [id])
        .then(result => result.rows[0]);
    },

    // Or use async/await (equivalent, more readable)
    users: async (parent, { filter }, context) => {
      const users = await context.db.users.findAll(filter);
      return users;
    },
  },
};
```

---

## Error Handling in Resolvers

```javascript
const resolvers = {
  Query: {
    user: async (parent, { id }, context) => {
      const user = await context.db.users.findById(id);

      if (!user) {
        throw new GraphQLError('User not found', {
          extensions: {
            code: 'USER_NOT_FOUND',
            http: { status: 404 },
          },
        });
      }

      return user;
    },
  },

  Mutation: {
    createUser: async (parent, { input }, context) => {
      if (!context.user) {
        throw new GraphQLError('Authentication required', {
          extensions: { code: 'UNAUTHENTICATED' },
        });
      }

      if (context.user.role !== 'ADMIN') {
        throw new GraphQLError('Insufficient permissions', {
          extensions: { code: 'FORBIDDEN' },
        });
      }

      try {
        return await context.db.users.create(input);
      } catch (err) {
        if (err.code === '23505') {  // unique violation
          throw new GraphQLError('Email already in use', {
            extensions: { code: 'DUPLICATE_EMAIL', field: 'email' },
          });
        }
        throw err;
      }
    },
  },
};
```

---

## Resolver Best Practices: Thin Resolvers

Keep resolvers thin — delegate to service classes:

```javascript
// Bad: business logic in resolver
const resolvers = {
  Mutation: {
    createUser: async (parent, { input }, context) => {
      const hash = await bcrypt.hash(input.password, 12);
      const existing = await context.db.users.findByEmail(input.email);
      if (existing) throw new Error('Email taken');
      const user = await context.db.users.create({ ...input, password: hash });
      await sendWelcomeEmail(user.email);
      return user;
    },
  },
};

// Good: resolver delegates to service
const resolvers = {
  Mutation: {
    createUser: (parent, { input }, context) =>
      context.services.userService.create(input),
  },
};

// UserService handles the logic
class UserService {
  async create(input) {
    const hash = await bcrypt.hash(input.password, 12);
    const existing = await this.db.users.findByEmail(input.email);
    if (existing) throw new DuplicateEmailError(input.email);
    const user = await this.db.users.create({ ...input, password: hash });
    await this.emailService.sendWelcome(user.email);
    return user;
  }
}
```

---

## The `info` Argument

The `info` argument contains the full query AST — useful for advanced use cases:

```javascript
const resolvers = {
  Query: {
    user: async (parent, { id }, context, info) => {
      // Check which fields were requested — optimize DB query
      const requestedFields = info.fieldNodes[0].selectionSet.selections
        .map(s => s.name.value);

      if (!requestedFields.includes('posts')) {
        // Don't join posts if they weren't requested
        return context.db.users.findById(id, { includePosts: false });
      }

      return context.db.users.findById(id, { includePosts: true });
    },
  },
};
```

Libraries like `graphql-parse-resolve-info` and `join-monster` automate this pattern.

---

## Abstract Type Resolvers

For interfaces and unions, you need a `__resolveType` resolver:

```javascript
const resolvers = {
  // For union types
  SearchResult: {
    __resolveType(obj) {
      if (obj.title) return 'Post';
      if (obj.email) return 'User';
      if (obj.price) return 'Product';
      return null;
    },
  },

  // For interface types
  Node: {
    __resolveType(obj) {
      return obj.__typename || obj.constructor.name;
    },
  },
};
```
# Chapter 7: The Type System

## GraphQL's Type System

GraphQL is **strongly typed**. Every field, every argument, every return value has a type known at schema build time. This enables:
- Auto-generated documentation
- IDE autocomplete
- Client-side type safety (TypeScript generation)
- Automatic validation of queries before execution

---

## Named Types

### Scalar Types

Scalars are leaf values — they resolve to concrete data with no sub-fields.

**Built-in scalars:**
```graphql
String    # UTF-8 character sequence
Int       # signed 32-bit integer
Float     # double-precision floating-point
Boolean   # true or false
ID        # unique identifier (serialized as String, often an integer or UUID)
```

**Custom scalars:**
```graphql
scalar Date        # 'YYYY-MM-DD'
scalar DateTime    # ISO 8601 with timezone
scalar Email       # validated email
scalar URL         # validated URL
scalar JSON        # any JSON value
scalar UUID        # UUID string
scalar PositiveInt # integer > 0
scalar NonEmptyString # non-blank string

# Usage
type User {
  id: ID!
  email: Email!
  createdAt: DateTime!
  birthDate: Date
  website: URL
  metadata: JSON
}
```

Implementing a custom scalar (JavaScript):
```javascript
import { GraphQLScalarType, Kind } from 'graphql';

const DateScalar = new GraphQLScalarType({
  name: 'Date',
  description: 'Date custom scalar type (YYYY-MM-DD)',
  serialize(value) {
    // outgoing: JS Date → string
    return value instanceof Date
      ? value.toISOString().split('T')[0]
      : value;
  },
  parseValue(value) {
    // incoming variable: string → JS Date
    return new Date(value);
  },
  parseLiteral(ast) {
    // incoming inline: string literal → JS Date
    if (ast.kind === Kind.STRING) {
      return new Date(ast.value);
    }
    return null;
  },
});
```

---

### Object Types

```graphql
type User {
  id: ID!
  name: String!
  email: String!
  bio: String          # nullable field
  role: UserRole!
  posts: [Post!]!
  createdAt: DateTime!
}
```

Every object type field can have arguments:
```graphql
type User {
  id: ID!
  name: String!
  posts(limit: Int = 10, status: PostStatus): [Post!]!
  avatar(size: Int = 100, format: ImageFormat = WEBP): String
}
```

---

### Input Object Types

Used exclusively for mutation/query arguments — cannot be used as output:

```graphql
input CreatePostInput {
  title: String!
  body: String!
  tags: [String!]! = []
  status: PostStatus = DRAFT
}

input PostFilter {
  authorId: ID
  status: PostStatus
  search: String
  tags: [String!]
  createdAfter: DateTime
  createdBefore: DateTime
}

type Query {
  posts(filter: PostFilter, limit: Int = 20, offset: Int = 0): [Post!]!
}
```

Input types CAN reference other input types:
```graphql
input AddressInput {
  street: String!
  city: String!
  country: String!
}

input CreateUserInput {
  name: String!
  email: String!
  address: AddressInput    # nested input
}
```

---

### Enum Types

```graphql
enum UserRole {
  ADMIN
  MODERATOR
  USER
  GUEST
}

enum PostStatus {
  DRAFT
  PUBLISHED
  ARCHIVED
  DELETED
}

enum SortDirection {
  ASC
  DESC
}

enum ImageFormat {
  JPEG
  PNG
  WEBP
  AVIF
}
```

Enum values are all-caps by convention. Internally they can map to different values:

```javascript
const resolvers = {
  UserRole: {
    ADMIN: 'admin',       // GraphQL value → internal value
    MODERATOR: 'mod',
    USER: 'user',
    GUEST: 'guest',
  },
};
```

---

### Interface Types

Interfaces define a contract. Multiple types can implement one interface:

```graphql
interface Node {
  id: ID!
}

interface Timestamped {
  createdAt: DateTime!
  updatedAt: DateTime!
}

interface Searchable {
  searchTitle: String!
  searchBody: String!
}

type User implements Node & Timestamped {
  id: ID!
  name: String!
  email: String!
  createdAt: DateTime!
  updatedAt: DateTime!
}

type Post implements Node & Timestamped & Searchable {
  id: ID!
  title: String!
  body: String!
  createdAt: DateTime!
  updatedAt: DateTime!
  searchTitle: String!   # could resolve to title
  searchBody: String!    # could resolve to body
}
```

Using interfaces in queries:
```graphql
query {
  nodes(ids: ["user:1", "post:2"]) {
    id  # from Node interface — always available
    ... on User { name email }
    ... on Post { title body }
  }
}
```

---

### Union Types

Unions express "this can be one of these types" — no common fields required:

```graphql
union SearchResult = User | Post | Tag | Product

union MutationResult = CreatePostSuccess | CreatePostError

type CreatePostSuccess {
  post: Post!
}

type CreatePostError {
  code: String!
  message: String!
  field: String
}

type Mutation {
  createPost(input: CreatePostInput!): MutationResult!
}
```

Client:
```graphql
mutation {
  createPost(input: $input) {
    ... on CreatePostSuccess {
      post { id title }
    }
    ... on CreatePostError {
      code
      message
      field
    }
  }
}
```

---

## Nullability

Nullability is one of the most important design decisions in GraphQL.

```graphql
String      # nullable — field can return null
String!     # non-null — field always returns a value (never null)
[String]    # nullable list — the list itself can be null
[String!]   # non-null items, but list can be null
[String]!   # non-null list, but items can be null
[String!]!  # non-null list AND non-null items
```

**When to use `!` (non-null):**
- Always use `!` for IDs, required fields, reliable data
- Leave nullable for optional fields, fields that might not exist, fields that could fail

**Null propagation:** If a non-null field returns null, the error bubbles up to the nearest nullable ancestor:

```graphql
type User {
  id: ID!
  name: String!       # non-null
  address: Address    # nullable
}

type Address {
  city: String!       # non-null inside address
}

# If address.city resolver throws:
# - address becomes null (nullable — error stops here)
# - user still returns with address: null

# If user.name resolver throws:
# - name can't be null (non-null) → user becomes null
# - this propagates up... until it hits a nullable field or the root
```

**Recommendation for server-side schema:**
- Use `!` generously for fields the server controls
- Leave nullable for fields involving external systems, optional data, or that can meaningfully not exist

---

## Type Modifiers Summary

```graphql
User     # Object type — nullable
User!    # Object type — non-null
[User]   # List of nullable Users — list is nullable
[User!]  # List of non-null Users — list is nullable
[User]!  # List of nullable Users — list is non-null
[User!]! # List of non-null Users — list is non-null
```

---

## Generic / Parameterized Types (Pagination Example)

GraphQL doesn't have generics, but you can define typed connection types:

```graphql
# The Connection pattern (Relay spec)
type UserConnection {
  edges: [UserEdge!]!
  pageInfo: PageInfo!
  totalCount: Int!
}

type UserEdge {
  node: User!
  cursor: String!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}

type Query {
  users(first: Int, after: String, last: Int, before: String): UserConnection!
  posts(first: Int, after: String): PostConnection!
}
```

---

## TypeScript Integration

Generate types from schema automatically with GraphQL Code Generator:

```bash
npm install -D @graphql-codegen/cli @graphql-codegen/typescript
```

```yaml
# codegen.yml
generates:
  src/generated/types.ts:
    schema: schema.graphql
    plugins:
      - typescript
      - typescript-resolvers
```

```typescript
// Generated types used in resolvers
import { Resolvers } from './generated/types';

const resolvers: Resolvers = {
  Query: {
    user: async (parent, { id }, context): Promise<User> => {
      // TypeScript knows id: string, return type: User
      return context.db.users.findById(id);
    },
  },
};
```
# Chapter 8: Variables & Fragments

## Variables

Variables separate the query structure from the data values. They are sent alongside the query — never interpolated into the query string.

### Why Variables?
```graphql
# BAD: string interpolation (like SQL injection — never do this)
const query = `query { user(id: "${userId}") { name } }`;

# GOOD: variables
const query = `query GetUser($id: ID!) { user(id: $id) { name } }`;
const variables = { id: userId };
```

Benefits:
- **Safe** — no injection attacks
- **Cached** — same query structure, different variables = server can cache the parsed query
- **Type-checked** — variable types are validated against the schema
- **IDE support** — autocomplete for variable names and types

---

## Variable Declaration Syntax

```graphql
query OperationName($varName: Type = defaultValue) {
  field(arg: $varName)
}
```

```graphql
query GetUser(
  $id: ID!                      # required, no default
  $limit: Int = 10              # optional, with default
  $status: PostStatus           # optional, nullable
  $tags: [String!]! = []       # required non-null list, default to empty array
  $includePrivate: Boolean = false
) {
  user(id: $id) {
    id
    name
    email @include(if: $includePrivate)
    posts(limit: $limit, status: $status, tags: $tags) {
      id
      title
    }
  }
}
```

Variables JSON:
```json
{
  "id": "42",
  "limit": 5,
  "status": "PUBLISHED",
  "tags": ["javascript", "graphql"],
  "includePrivate": true
}
```

---

## Variable Types and Nullability

```graphql
query Examples(
  $a: String      # nullable String — can send null or omit
  $b: String!     # non-null String — must provide a value
  $c: [String]    # nullable list of nullable strings
  $d: [String!]   # nullable list of non-null strings
  $e: [String]!   # non-null list of nullable strings
  $f: [String!]!  # non-null list of non-null strings — must provide non-empty
) { ... }
```

---

## Named Fragments

Fragments are reusable pieces of a query. They prevent repeating the same fields in multiple places.

```graphql
# Define fragment
fragment UserBasicFields on User {
  id
  name
  email
  avatar(size: 64)
}

fragment PostPreviewFields on Post {
  id
  title
  excerpt
  publishedAt
  tags
  author {
    ...UserBasicFields    # compose fragments
  }
}

# Use in queries
query GetHomeFeed {
  me {
    ...UserBasicFields
    unreadCount
  }
  featuredPosts: posts(featured: true, limit: 3) {
    ...PostPreviewFields
  }
  recentPosts: posts(limit: 10) {
    ...PostPreviewFields
  }
}

query GetUserProfile($id: ID!) {
  user(id: $id) {
    ...UserBasicFields
    bio
    joinedAt: createdAt
    postCount
  }
}
```

Fragment rules:
- Must specify `on TypeName` — the type it applies to
- Only usable in queries on objects of that type or implementing types
- Cannot be self-referential (no circular fragments)
- Must be used — unused fragments cause validation errors

---

## Inline Fragments

Inline fragments are anonymous, one-use fragments. Required for union types and interfaces:

```graphql
query GetSearchResults($query: String!) {
  search(query: $query) {
    # Fields available on ALL types (if interface)
    __typename

    ... on User {
      id
      name
      email
      avatar
    }

    ... on Post {
      id
      title
      excerpt
      author { name }
    }

    ... on Product {
      id
      name
      price
      inStock
    }
  }
}
```

Inline fragments on interfaces (to access specific fields):
```graphql
query GetNodes($ids: [ID!]!) {
  nodes(ids: $ids) {
    id  # from Node interface — always available

    ... on User {
      name email
    }

    ... on Post {
      title status
    }
  }
}
```

---

## Fragment Variables (via Query Variables)

Fragments don't have their own variables — they share the enclosing operation's variables:

```graphql
fragment PostDetails on Post {
  id
  title
  body
  # Uses the operation's $showAuthor variable
  author @include(if: $showAuthor) {
    name
    avatar
  }
}

query GetPost($id: ID!, $showAuthor: Boolean = true) {
  post(id: $id) {
    ...PostDetails
  }
}
```

---

## Directives on Fragments

```graphql
fragment UserStats on User {
  postCount
  followerCount
  followingCount
}

query GetProfile($id: ID!, $withStats: Boolean!) {
  user(id: $id) {
    id
    name
    ...UserStats @include(if: $withStats)
  }
}
```

---

## @defer — Deferred Fragments (Experimental)

`@defer` allows a fragment to be streamed progressively — the main query returns first, then the deferred part arrives:

```graphql
query GetPost($id: ID!) {
  post(id: $id) {
    id
    title
    body         # arrives immediately

    ... on Post @defer {
      # arrives later (slow to load)
      comments { id body author { name } }
      relatedPosts { id title }
    }
  }
}
```

Useful for non-critical below-the-fold content.

---

## Fragment Colocation (Best Practice)

In component-based UIs, define fragments next to the components that use them:

```javascript
// UserCard.jsx
export const USER_CARD_FRAGMENT = gql`
  fragment UserCardFields on User {
    id
    name
    avatar(size: 48)
    role
  }
`;

function UserCard({ user }) {
  return (
    <div>
      <img src={user.avatar} />
      <h3>{user.name}</h3>
      <span>{user.role}</span>
    </div>
  );
}

// PostFeed.jsx — composes UserCard's fragment
import { USER_CARD_FRAGMENT } from './UserCard';

const GET_POSTS = gql`
  query GetPosts($limit: Int!) {
    posts(limit: $limit) {
      id
      title
      author {
        ...UserCardFields    # UserCard declares what it needs
      }
    }
  }
  ${USER_CARD_FRAGMENT}
`;
```

This pattern (popularized by Relay) ensures components own their data requirements — great for maintainability.

---

## Combining Variables and Fragments

```graphql
fragment OrderFields on Order {
  id
  status
  totalAmount
  items {
    id
    productName
    quantity
    price
  }
  createdAt
}

fragment PaginationFields on OrderConnection {
  pageInfo {
    hasNextPage
    endCursor
  }
  totalCount
}

query GetOrders(
  $userId: ID!
  $first: Int = 10
  $after: String
  $status: OrderStatus
) {
  user(id: $userId) {
    orders(first: $first, after: $after, status: $status) {
      ...PaginationFields
      edges {
        node {
          ...OrderFields
        }
      }
    }
  }
}
```

---

## Common Fragment Patterns

```graphql
# Error fragment for consistent error handling
fragment ErrorFields on UserError {
  field
  message
  code
}

# Payload fragment
fragment CreatePostPayloadFields on CreatePostPayload {
  post { id title status }
  errors { ...ErrorFields }
}

# Use across all post mutations
mutation CreatePost($input: CreatePostInput!) {
  createPost(input: $input) {
    ...CreatePostPayloadFields
  }
}

mutation UpdatePost($id: ID!, $input: UpdatePostInput!) {
  updatePost(id: $id, input: $input) {
    ...CreatePostPayloadFields
  }
}
```
# Chapter 9: Directives

## What are Directives?

Directives are annotations that modify the execution behavior of fields, types, or arguments. They start with `@`.

```graphql
query {
  user(id: "1") {
    name
    email @include(if: $showEmail)    # conditional field
    phone @skip(if: $hidePhone)       # skip this field
  }
}
```

---

## Built-in Directives

### @include(if: Boolean!)

Include a field only if the condition is true:

```graphql
query GetUser($withPosts: Boolean!) {
  user(id: "1") {
    id
    name
    posts @include(if: $withPosts) {
      id
      title
    }
  }
}

# Variables: { "withPosts": false }
# posts field is omitted from both query and response
```

### @skip(if: Boolean!)

Skip (exclude) a field if the condition is true:

```graphql
query GetProfile($hideSensitive: Boolean!) {
  me {
    id
    name
    email @skip(if: $hideSensitive)
    phone @skip(if: $hideSensitive)
    publicProfile {
      bio
      website
    }
  }
}
```

`@skip` and `@include` are logical opposites but both exist for clarity. You can use either.

### @deprecated(reason: String)

Marks a field or enum value as deprecated:

```graphql
type User {
  id: ID!
  name: String!
  username: String @deprecated(reason: "Use `name` instead. Will be removed in v3.")
  profileImage: String @deprecated(reason: "Use `avatar(size: Int)` instead.")
  avatar(size: Int = 100): String
}

type Query {
  allUsers: [User!]! @deprecated(reason: "Use `users(filter: ...)` instead.")
  users(filter: UserFilter): [User!]!
}

enum PostStatus {
  DRAFT
  PUBLISHED
  REMOVED @deprecated(reason: "Use ARCHIVED instead.")
  ARCHIVED
}
```

Deprecated fields still work but appear with strikethrough in GraphiQL and generate warnings in Apollo Client.

### @specifiedBy(url: String!)

Specifies documentation for a custom scalar:

```graphql
scalar DateTime @specifiedBy(url: "https://scalars.graphql.org/andimarek/date-time")
scalar UUID @specifiedBy(url: "https://tools.ietf.org/html/rfc4122")
```

---

## Custom Server-Side Directives

### Schema Directives (applied to type definitions):

```graphql
# Declare the directive
directive @auth(requires: Role = USER) on FIELD_DEFINITION | OBJECT

directive @rateLimit(
  max: Int!
  window: String!
) on FIELD_DEFINITION

directive @cacheControl(
  maxAge: Int
  scope: CacheScope
) on FIELD_DEFINITION | OBJECT

enum Role { ADMIN MODERATOR USER GUEST }
enum CacheScope { PUBLIC PRIVATE }

# Use it in schema
type Query {
  me: User @auth
  adminDashboard: Dashboard @auth(requires: ADMIN)
  publicPosts: [Post!]! @cacheControl(maxAge: 60, scope: PUBLIC)
  search(q: String!): [SearchResult!]! @rateLimit(max: 10, window: "1m")
}

type User @auth {
  id: ID!
  email: String!
  creditCard: String @auth(requires: ADMIN)  # field-level override
}
```

### Implementing Schema Directives (JavaScript with graphql-tools):

```javascript
import { mapSchema, getDirective, MapperKind } from '@graphql-tools/utils';
import { defaultFieldResolver } from 'graphql';

function authDirectiveTransformer(schema, directiveName) {
  return mapSchema(schema, {
    [MapperKind.OBJECT_FIELD]: (fieldConfig) => {
      const directive = getDirective(schema, fieldConfig, directiveName)?.[0];
      if (!directive) return fieldConfig;

      const { requires = 'USER' } = directive;
      const { resolve = defaultFieldResolver } = fieldConfig;

      return {
        ...fieldConfig,
        resolve(source, args, context, info) {
          if (!context.user) {
            throw new GraphQLError('Authentication required', {
              extensions: { code: 'UNAUTHENTICATED' },
            });
          }

          const roleHierarchy = ['GUEST', 'USER', 'MODERATOR', 'ADMIN'];
          const userLevel = roleHierarchy.indexOf(context.user.role);
          const requiredLevel = roleHierarchy.indexOf(requires);

          if (userLevel < requiredLevel) {
            throw new GraphQLError(`Requires role: ${requires}`, {
              extensions: { code: 'FORBIDDEN' },
            });
          }

          return resolve(source, args, context, info);
        },
      };
    },
  });
}

// Apply to schema
let schema = makeExecutableSchema({ typeDefs, resolvers });
schema = authDirectiveTransformer(schema, 'auth');
```

---

## @cacheControl Directive

Used with Apollo Server to set cache hints:

```graphql
directive @cacheControl(
  maxAge: Int
  scope: CacheControlScope
  inheritMaxAge: Boolean
) on FIELD_DEFINITION | OBJECT | INTERFACE | UNION

enum CacheControlScope { PUBLIC PRIVATE }

type Query {
  product(id: ID!): Product @cacheControl(maxAge: 300, scope: PUBLIC)  # 5 minutes
  me: User @cacheControl(maxAge: 0, scope: PRIVATE)  # never cache
  config: SiteConfig @cacheControl(maxAge: 3600, scope: PUBLIC)  # 1 hour
}
```

---

## @defer and @stream (Incremental Delivery)

These are experimental directives for streaming large responses:

```graphql
# @defer — return the rest of the query first, then stream the deferred part
query GetPost($id: ID!) {
  post(id: $id) {
    id
    title
    body

    # This arrives later
    ... @defer {
      comments { id body author { name } }
    }

    # This arrives even later (label for tracking)
    ... @defer(label: "related") {
      relatedPosts { id title }
    }
  }
}

# @stream — stream list items one by one
query GetPosts {
  posts @stream(initialCount: 3) {
    id
    title
  }
  # Returns first 3 immediately, then streams the rest
}
```

---

## @external, @key, @requires (Federation)

Used in Apollo Federation to compose subgraph schemas:

```graphql
# user-service subgraph
type User @key(fields: "id") {
  id: ID!
  name: String!
  email: String!
}

# review-service subgraph
type User @key(fields: "id") @extends {
  id: ID! @external          # field comes from user-service
  reviews: [Review!]!        # new field added by this service
}

type Review {
  id: ID!
  body: String!
  rating: Int!
  author: User!
}

type Query {
  reviews: [Review!]!
}
```

---

## Directive Locations

```graphql
# Where directives can be applied
directive @example on
  # Schema Definition Locations
  | SCHEMA
  | SCALAR
  | OBJECT
  | FIELD_DEFINITION
  | ARGUMENT_DEFINITION
  | INTERFACE
  | UNION
  | ENUM
  | ENUM_VALUE
  | INPUT_OBJECT
  | INPUT_FIELD_DEFINITION

  # Query Locations (client-side)
  | QUERY
  | MUTATION
  | SUBSCRIPTION
  | FIELD
  | FRAGMENT_DEFINITION
  | FRAGMENT_SPREAD
  | INLINE_FRAGMENT
  | VARIABLE_DEFINITION
```

---

## Directive Summary

| Directive | Location | Purpose |
|---|---|---|
| `@include(if: Boolean!)` | Field, fragment | Include field conditionally |
| `@skip(if: Boolean!)` | Field, fragment | Skip field conditionally |
| `@deprecated(reason: String)` | Field, enum value | Mark as deprecated |
| `@specifiedBy(url: String!)` | Scalar | Documentation URL |
| `@auth(requires: Role)` | Field, object | Authorization (custom) |
| `@cacheControl(maxAge: Int)` | Field, object | HTTP caching hints |
| `@defer` | Fragment | Stream deferred data |
| `@stream(initialCount: Int)` | Field (list) | Stream list items |
| `@key(fields: String!)` | Object | Federation entity key |
| `@external` | Field | Federation — from another subgraph |
# Chapter 10: Authentication & Authorization

## Authentication vs Authorization

- **Authentication** — Who are you? (verify identity via JWT, session, API key)
- **Authorization** — What can you do? (check permissions for a field or operation)

---

## Authentication via Context

The most common pattern: verify the JWT on every request and put the user in context.

```javascript
// Apollo Server context function
const server = new ApolloServer({
  typeDefs,
  resolvers,
  context: async ({ req }) => {
    const token = req.headers.authorization?.replace('Bearer ', '');

    let currentUser = null;
    if (token) {
      try {
        const payload = jwt.verify(token, process.env.JWT_SECRET);
        currentUser = await db.users.findById(payload.userId);
      } catch (err) {
        // invalid or expired token — user stays null
        // don't throw here — public endpoints need to work
      }
    }

    return {
      currentUser,
      db,
      dataloaders: createDataLoaders(db),
    };
  },
});
```

```javascript
// In resolvers
const resolvers = {
  Query: {
    me: (parent, args, context) => {
      if (!context.currentUser) {
        throw new GraphQLError('Not authenticated', {
          extensions: { code: 'UNAUTHENTICATED' },
        });
      }
      return context.currentUser;
    },
  },
};
```

---

## Login Mutation

```graphql
type AuthPayload {
  token: String!
  refreshToken: String!
  user: User!
}

type Mutation {
  login(email: String!, password: String!): AuthPayload!
  logout: Boolean!
  refreshToken(token: String!): AuthPayload!
  register(input: RegisterInput!): AuthPayload!
}
```

```javascript
const resolvers = {
  Mutation: {
    login: async (parent, { email, password }, context) => {
      const user = await context.db.users.findByEmail(email);
      if (!user || !await bcrypt.compare(password, user.passwordHash)) {
        throw new GraphQLError('Invalid credentials', {
          extensions: { code: 'UNAUTHENTICATED' },
        });
      }

      const token = jwt.sign(
        { userId: user.id, role: user.role },
        process.env.JWT_SECRET,
        { expiresIn: '15m' }
      );
      const refreshToken = jwt.sign(
        { userId: user.id },
        process.env.REFRESH_SECRET,
        { expiresIn: '7d' }
      );

      return { token, refreshToken, user };
    },
  },
};
```

---

## Authorization Patterns

### Pattern 1: Check in Each Resolver

```javascript
const resolvers = {
  Query: {
    adminStats: (parent, args, context) => {
      requireAuth(context);
      requireRole(context, 'ADMIN');
      return context.db.getAdminStats();
    },
  },
  Mutation: {
    deleteUser: async (parent, { id }, context) => {
      requireAuth(context);
      requireRole(context, 'ADMIN');
      return context.db.users.delete(id);
    },
    updatePost: async (parent, { id, input }, context) => {
      requireAuth(context);
      const post = await context.db.posts.findById(id);
      if (post.authorId !== context.currentUser.id) {
        throw new GraphQLError('Not your post', { extensions: { code: 'FORBIDDEN' } });
      }
      return context.db.posts.update(id, input);
    },
  },
};

// Helper functions
function requireAuth(context) {
  if (!context.currentUser) {
    throw new GraphQLError('Authentication required', {
      extensions: { code: 'UNAUTHENTICATED' },
    });
  }
}

function requireRole(context, role) {
  requireAuth(context);
  const roles = ['GUEST', 'USER', 'MODERATOR', 'ADMIN'];
  if (roles.indexOf(context.currentUser.role) < roles.indexOf(role)) {
    throw new GraphQLError(`Requires role: ${role}`, {
      extensions: { code: 'FORBIDDEN' },
    });
  }
}
```

### Pattern 2: Schema Directives (@auth)

```graphql
directive @auth(requires: Role = USER) on FIELD_DEFINITION | OBJECT

type Query {
  me: User @auth
  adminDashboard: Dashboard @auth(requires: ADMIN)
  publicFeed: [Post!]!   # no @auth — public
}

type User @auth {
  id: ID!
  name: String!
  email: String!
  secretField: String @auth(requires: ADMIN)
}
```

```javascript
// Implement the directive transformer
function authDirectiveTransformer(schema) {
  return mapSchema(schema, {
    [MapperKind.OBJECT_FIELD]: (fieldConfig, fieldName, typeName, schema) => {
      const authDirective = getDirective(schema, fieldConfig, 'auth')?.[0];
      if (!authDirective) return fieldConfig;

      const { requires = 'USER' } = authDirective;
      const { resolve = defaultFieldResolver } = fieldConfig;

      return {
        ...fieldConfig,
        async resolve(source, args, context, info) {
          checkAuth(context, requires);
          return resolve(source, args, context, info);
        },
      };
    },
  });
}
```

### Pattern 3: GraphQL Shield

```javascript
import { shield, rule, and, or, not } from 'graphql-shield';

const isAuthenticated = rule()(async (parent, args, context) => {
  return context.currentUser !== null;
});

const isAdmin = rule()(async (parent, args, context) => {
  return context.currentUser?.role === 'ADMIN';
});

const isPostAuthor = rule()(async (parent, { id }, context) => {
  const post = await context.db.posts.findById(id);
  return post.authorId === context.currentUser?.id;
});

const permissions = shield({
  Query: {
    me: isAuthenticated,
    adminStats: and(isAuthenticated, isAdmin),
    publicFeed: not(false),   // allow all
  },
  Mutation: {
    createPost: isAuthenticated,
    updatePost: and(isAuthenticated, isPostAuthor),
    deletePost: and(isAuthenticated, or(isAdmin, isPostAuthor)),
    deleteUser: and(isAuthenticated, isAdmin),
  },
});

// Apply to schema
const schema = applyMiddleware(makeExecutableSchema({ typeDefs, resolvers }), permissions);
```

---

## Field-Level Authorization

Hide fields based on the viewer's role:

```javascript
const resolvers = {
  User: {
    email: (user, args, context) => {
      // Only show email to the user themselves or admins
      if (context.currentUser?.id === user.id || context.currentUser?.role === 'ADMIN') {
        return user.email;
      }
      return null;  // or throw forbidden error
    },

    creditCard: (user, args, context) => {
      requireRole(context, 'ADMIN');
      return user.creditCard;
    },
  },
};
```

---

## Resource-Based Authorization

Check ownership of specific resources:

```javascript
const resolvers = {
  Mutation: {
    updatePost: async (parent, { id, input }, context) => {
      requireAuth(context);

      const post = await context.db.posts.findById(id);
      if (!post) {
        throw new GraphQLError('Post not found', { extensions: { code: 'NOT_FOUND' } });
      }

      // Check ownership
      const canEdit = post.authorId === context.currentUser.id
        || context.currentUser.role === 'ADMIN'
        || context.currentUser.role === 'MODERATOR';

      if (!canEdit) {
        throw new GraphQLError('Not authorized to edit this post', {
          extensions: { code: 'FORBIDDEN' },
        });
      }

      return context.db.posts.update(id, input);
    },
  },
};
```

---

## API Keys (for service-to-service auth)

```javascript
context: async ({ req }) => {
  // Check for API key (for programmatic access)
  const apiKey = req.headers['x-api-key'];
  if (apiKey) {
    const service = await db.apiKeys.findByKey(apiKey);
    if (service) {
      return { currentUser: service, isServiceAccount: true, db };
    }
  }

  // Fall back to JWT
  const token = req.headers.authorization?.replace('Bearer ', '');
  // ... JWT verification
}
```

---

## Security Best Practices

- Never expose internal IDs in errors (use codes like `USER_NOT_FOUND`)
- Always validate input in resolvers — never trust client
- Disable introspection in production (or require auth for it)
- Limit query depth and complexity to prevent DoS
- Use query timeouts

```javascript
import depthLimit from 'graphql-depth-limit';
import { createComplexityLimitRule } from 'graphql-validation-complexity';

const server = new ApolloServer({
  schema,
  validationRules: [
    depthLimit(7),                           // max query depth
    createComplexityLimitRule(1000),         // max query complexity score
  ],
  introspection: process.env.NODE_ENV === 'development',
});
```
# Chapter 11: N+1 Problem & DataLoader

## The N+1 Problem

The N+1 problem is the most common performance issue in GraphQL APIs. It occurs when loading a list of N items, then making 1 additional query per item.

```graphql
query {
  posts {        # 1 query → returns 10 posts
    id
    title
    author {     # 10 queries → 1 per post author!
      name
      avatar
    }
  }
}
```

Total: **11 queries** for 10 posts. With 100 posts: **101 queries**. This is the N+1 problem.

```javascript
// The problematic resolver
const resolvers = {
  Post: {
    author: async (post, args, context) => {
      // This runs once PER POST — terrible!
      return context.db.users.findById(post.authorId);
    },
  },
};
```

---

## DataLoader — The Solution

**DataLoader** batches and caches database calls within a single request.

```
Without DataLoader:     With DataLoader:
  findUser(1)             collect: [1, 2, 3, 4, 5]
  findUser(2)           ─→  findUsers([1,2,3,4,5])  (1 query!)
  findUser(3)
  findUser(4)
  findUser(5)
```

### Install:
```bash
npm install dataloader
```

### Basic DataLoader:
```javascript
import DataLoader from 'dataloader';

// Batch function: receives array of keys, must return array of values in same order
async function batchLoadUsers(userIds) {
  const users = await db.query(
    'SELECT * FROM users WHERE id = ANY($1)',
    [userIds]
  );

  // IMPORTANT: must return values in the same order as keys
  const userMap = new Map(users.rows.map(user => [user.id, user]));
  return userIds.map(id => userMap.get(id) || new Error(`User ${id} not found`));
}

const userLoader = new DataLoader(batchLoadUsers);
```

---

## Per-Request DataLoader Instances

DataLoaders should be **created per request** — not shared globally. This ensures the cache is per-request and doesn't leak between users.

```javascript
// Create fresh loaders for each request
function createDataLoaders(db) {
  return {
    users: new DataLoader(async (ids) => {
      const users = await db.users.findManyByIds(ids);
      const userMap = new Map(users.map(u => [u.id, u]));
      return ids.map(id => userMap.get(id) ?? null);
    }),

    postsbyAuthor: new DataLoader(async (authorIds) => {
      const posts = await db.posts.findByAuthorIds(authorIds);
      // Group posts by authorId
      const grouped = new Map();
      authorIds.forEach(id => grouped.set(id, []));
      posts.forEach(post => grouped.get(post.authorId)?.push(post));
      return authorIds.map(id => grouped.get(id) ?? []);
    }),

    commentsByPost: new DataLoader(async (postIds) => {
      const comments = await db.comments.findByPostIds(postIds);
      const grouped = new Map();
      postIds.forEach(id => grouped.set(id, []));
      comments.forEach(c => grouped.get(c.postId)?.push(c));
      return postIds.map(id => grouped.get(id) ?? []);
    }),
  };
}

// In Apollo Server context
const server = new ApolloServer({
  typeDefs,
  resolvers,
  context: ({ req }) => ({
    currentUser: getUser(req),
    db,
    loaders: createDataLoaders(db), // fresh instance per request
  }),
});
```

---

## Using DataLoader in Resolvers

```javascript
const resolvers = {
  Post: {
    // Uses DataLoader — batched automatically!
    author: (post, args, context) => {
      return context.loaders.users.load(post.authorId);
    },

    comments: (post, args, context) => {
      return context.loaders.commentsByPost.load(post.id);
    },
  },

  Comment: {
    author: (comment, args, context) => {
      return context.loaders.users.load(comment.authorId);
    },
  },

  User: {
    posts: (user, args, context) => {
      return context.loaders.postsByAuthor.load(user.id);
    },
  },
};
```

Now fetching 10 posts with authors:
- **Before DataLoader:** 1 + 10 = 11 queries
- **After DataLoader:** 1 + 1 = **2 queries**

---

## DataLoader with Arguments

When field arguments affect the query, include them in the key:

```javascript
// Problem: posts(limit: 5) vs posts(limit: 10) need different loaders
// Solution: include args in the key, or use separate loaders

const postsByAuthorLoader = new DataLoader(
  async (keys) => {
    // keys = [{ authorId: 1, limit: 5 }, { authorId: 2, limit: 10 }, ...]
    // Group by limit, batch per limit value
    const grouped = {};
    keys.forEach(key => {
      const k = key.limit;
      if (!grouped[k]) grouped[k] = [];
      grouped[k].push(key.authorId);
    });

    const results = {};
    await Promise.all(
      Object.entries(grouped).map(async ([limit, authorIds]) => {
        const posts = await db.posts.findByAuthorIds(authorIds, { limit: parseInt(limit) });
        // map back...
      })
    );

    return keys.map(key => results[key.authorId] ?? []);
  },
  {
    // Custom cache key function
    cacheKeyFn: (key) => `${key.authorId}:${key.limit}`,
  }
);

// Usage in resolver
const resolvers = {
  User: {
    posts: (user, { limit = 10 }, context) => {
      return context.loaders.postsByAuthor.load({ authorId: user.id, limit });
    },
  },
};
```

---

## DataLoader Options

```javascript
const loader = new DataLoader(batchFn, {
  // Maximum batch size (default: unlimited)
  maxBatchSize: 100,

  // Disable batching (useful for testing)
  batch: false,

  // Disable caching (for mutations or time-sensitive data)
  cache: false,

  // Custom cache key function (default: identity)
  cacheKeyFn: (key) => JSON.stringify(key),

  // Custom cache map (e.g., use a Map with TTL)
  cacheMap: new Map(),
});
```

---

## Priming the Cache

Pre-populate the DataLoader cache to avoid unnecessary loads:

```javascript
const resolvers = {
  Query: {
    posts: async (parent, { limit }, context) => {
      const posts = await context.db.posts.findAll({ limit });

      // Prime the author loader with the authors we already have joined data for
      posts.forEach(post => {
        if (post.author) {
          context.loaders.users.prime(post.author.id, post.author);
        }
      });

      return posts;
    },
  },
};
```

---

## Common DataLoader Patterns

### Many-to-many relationships:
```javascript
// Load tags for posts
const tagsByPostLoader = new DataLoader(async (postIds) => {
  const rows = await db.query(`
    SELECT pt.post_id, t.id, t.name
    FROM post_tags pt
    JOIN tags t ON t.id = pt.tag_id
    WHERE pt.post_id = ANY($1)
  `, [postIds]);

  const grouped = new Map(postIds.map(id => [id, []]));
  rows.forEach(row => grouped.get(row.post_id)?.push({ id: row.id, name: row.name }));
  return postIds.map(id => grouped.get(id) ?? []);
});
```

### Count loader:
```javascript
const commentCountLoader = new DataLoader(async (postIds) => {
  const rows = await db.query(`
    SELECT post_id, COUNT(*) AS count
    FROM comments
    WHERE post_id = ANY($1)
    GROUP BY post_id
  `, [postIds]);

  const counts = new Map(rows.map(r => [r.post_id, parseInt(r.count)]));
  return postIds.map(id => counts.get(id) ?? 0);
});

// Resolver
const resolvers = {
  Post: {
    commentCount: (post, args, context) =>
      context.loaders.commentCount.load(post.id),
  },
};
```

---

## Measuring the Improvement

```javascript
// Log all DB queries during a request
const db = {
  query: async (sql, params) => {
    console.log(`DB QUERY: ${sql}`);
    return pool.query(sql, params);
  },
};

// Without DataLoader:
// DB QUERY: SELECT * FROM posts LIMIT 10
// DB QUERY: SELECT * FROM users WHERE id = 1
// DB QUERY: SELECT * FROM users WHERE id = 2
// DB QUERY: SELECT * FROM users WHERE id = 1  (duplicate!)
// DB QUERY: SELECT * FROM users WHERE id = 3
// ... 10 more queries

// With DataLoader:
// DB QUERY: SELECT * FROM posts LIMIT 10
// DB QUERY: SELECT * FROM users WHERE id = ANY('{1,2,3,4,5,6,7,8,9,10}')
// Total: 2 queries (with deduplication — note id=1 queried once)
```
# Chapter 12: Error Handling

## GraphQL Error Structure

GraphQL responses always return HTTP 200, with errors in the `errors` array alongside partial `data`:

```json
{
  "data": {
    "user": null,
    "posts": [{ "id": "1", "title": "Hello" }]
  },
  "errors": [
    {
      "message": "User not found",
      "locations": [{ "line": 2, "column": 3 }],
      "path": ["user"],
      "extensions": {
        "code": "USER_NOT_FOUND",
        "http": { "status": 404 }
      }
    }
  ]
}
```

Key points:
- `errors` and `data` can coexist — **partial success** is valid
- `locations` — where in the query document the error occurred
- `path` — which field in the response failed
- `extensions` — custom metadata (error codes, HTTP status, etc.)

---

## GraphQLError

The `GraphQLError` class is the standard way to throw errors in resolvers:

```javascript
import { GraphQLError } from 'graphql';

const resolvers = {
  Query: {
    user: async (parent, { id }, context) => {
      const user = await context.db.users.findById(id);

      if (!user) {
        throw new GraphQLError(`User with id '${id}' not found`, {
          extensions: {
            code: 'USER_NOT_FOUND',
            http: { status: 404 },
            userId: id,
          },
        });
      }

      return user;
    },
  },
};
```

---

## Standard Error Codes

Use consistent error codes across your API:

```javascript
const ErrorCodes = {
  // Authentication
  UNAUTHENTICATED: 'UNAUTHENTICATED',   // not logged in
  FORBIDDEN: 'FORBIDDEN',               // logged in but not allowed

  // Input validation
  BAD_USER_INPUT: 'BAD_USER_INPUT',     // invalid input data
  VALIDATION_ERROR: 'VALIDATION_ERROR', // business rule violation

  // Resource errors
  NOT_FOUND: 'NOT_FOUND',
  CONFLICT: 'CONFLICT',                 // duplicate resource

  // Server errors
  INTERNAL_SERVER_ERROR: 'INTERNAL_SERVER_ERROR',
  SERVICE_UNAVAILABLE: 'SERVICE_UNAVAILABLE',
};

// Helper functions
const errors = {
  unauthenticated: () =>
    new GraphQLError('Authentication required', {
      extensions: { code: ErrorCodes.UNAUTHENTICATED, http: { status: 401 } },
    }),

  forbidden: (message = 'Access denied') =>
    new GraphQLError(message, {
      extensions: { code: ErrorCodes.FORBIDDEN, http: { status: 403 } },
    }),

  notFound: (resource, id) =>
    new GraphQLError(`${resource} '${id}' not found`, {
      extensions: { code: ErrorCodes.NOT_FOUND, http: { status: 404 }, resource, id },
    }),

  conflict: (message) =>
    new GraphQLError(message, {
      extensions: { code: ErrorCodes.CONFLICT, http: { status: 409 } },
    }),

  badInput: (field, message) =>
    new GraphQLError(message, {
      extensions: { code: ErrorCodes.BAD_USER_INPUT, field, http: { status: 400 } },
    }),
};
```

---

## Partial Errors (Field-Level)

When one field fails, the rest of the query can still succeed:

```graphql
query {
  user(id: "1") {    # succeeds
    id
    name
  }
  user(id: "999") {  # fails — user not found
    id
    name
  }
}
```

```json
{
  "data": {
    "firstUser": { "id": "1", "name": "Alice" },
    "missingUser": null
  },
  "errors": [
    {
      "message": "User '999' not found",
      "path": ["missingUser"],
      "extensions": { "code": "NOT_FOUND" }
    }
  ]
}
```

This is powerful — clients get as much data as possible even when parts fail.

---

## Null Propagation

When a non-null field throws an error, the error bubbles up to the nearest nullable ancestor:

```graphql
type Post {
  id: ID!
  title: String!    # non-null
  author: User!     # non-null — if this throws, post becomes null
}

type Query {
  posts: [Post!]!   # non-null list — but individual items can be null? No — Post! is non-null
}
```

```json
// If author resolver throws for post with id "1":
{
  "data": {
    "posts": [
      null,                                   // post became null due to error propagation
      { "id": "2", "title": "...", "author": { "name": "Bob" } }
    ]
  },
  "errors": [
    { "message": "...", "path": ["posts", 0, "author"] }
  ]
}
```

**Design tip:** Use nullable fields strategically to contain error propagation.

---

## Mutation Error Patterns

### Pattern 1: Throw GraphQLError (simple)
```javascript
createPost: async (parent, { input }, context) => {
  // Throws → goes to errors[] array, data.createPost = null
  if (!input.title) throw errors.badInput('title', 'Title is required');
  return db.posts.create(input);
},
```

### Pattern 2: Payload type with errors (recommended for mutations)
```graphql
type CreatePostPayload {
  post: Post
  errors: [UserError!]!
}

type UserError {
  field: String
  message: String!
  code: String!
}
```

```javascript
createPost: async (parent, { input }, context) => {
  const validationErrors = validate(input);
  if (validationErrors.length > 0) {
    // Return errors in payload — NOT as thrown error
    return { post: null, errors: validationErrors };
  }
  const post = await db.posts.create(input);
  return { post, errors: [] };
},
```

Payload pattern advantages:
- Business errors (validation) vs system errors (crashes) are clearly separated
- Client can handle them differently
- TypeScript can type-check the error handling

---

## Error Formatting (Server-Side)

Sanitize errors before sending to clients — don't leak internal details:

```javascript
const server = new ApolloServer({
  schema,
  formatError: (formattedError, error) => {
    // Log the original error for debugging
    console.error('GraphQL Error:', error);

    // Hide implementation details in production
    if (process.env.NODE_ENV === 'production') {
      // Keep expected errors (ones we threw ourselves)
      if (formattedError.extensions?.code) {
        return formattedError;
      }

      // Mask unexpected errors (DB errors, crashes, etc.)
      return {
        message: 'An unexpected error occurred',
        extensions: { code: 'INTERNAL_SERVER_ERROR' },
      };
    }

    // In development, show everything
    return formattedError;
  },
});
```

---

## Global Error Handler

```javascript
// Apollo Server plugin for error logging
const errorLoggingPlugin = {
  async requestDidStart() {
    return {
      async didEncounterErrors({ errors, request, contextValue }) {
        errors.forEach((error) => {
          // Only log unexpected errors — not user input errors
          if (!error.extensions?.code || error.extensions.code === 'INTERNAL_SERVER_ERROR') {
            logger.error({
              message: error.message,
              stack: error.originalError?.stack,
              path: error.path,
              query: request.query,
              userId: contextValue?.currentUser?.id,
            });
          }
        });
      },
    };
  },
};
```

---

## Client-Side Error Handling

```javascript
// Apollo Client
const { data, error, loading } = useQuery(GET_USER, { variables: { id } });

if (error) {
  // error.graphQLErrors — errors from the GraphQL server
  // error.networkError — HTTP/network failures
  // error.message — summary message

  const code = error.graphQLErrors[0]?.extensions?.code;

  if (code === 'UNAUTHENTICATED') redirect('/login');
  if (code === 'NOT_FOUND') return <NotFound />;
  return <ErrorMessage message={error.message} />;
}

// For mutations with payload pattern
const [createPost] = useMutation(CREATE_POST, {
  onCompleted: (data) => {
    if (data.createPost.errors.length > 0) {
      // Handle field errors
      data.createPost.errors.forEach(err => {
        setFieldError(err.field, err.message);
      });
    } else {
      // Success
      router.push(`/posts/${data.createPost.post.id}`);
    }
  },
  onError: (error) => {
    // Network or unexpected server error
    showToast({ type: 'error', message: 'Something went wrong' });
  },
});
```

---

## Error Masking Best Practices

| Error Type | Expose to Client | Log Server-Side |
|---|---|---|
| Input validation (bad user input) | Yes — full message | No (expected) |
| Business rules (insufficient funds) | Yes — full message | No (expected) |
| Authentication/authorization | Yes — generic message | No (expected) |
| Database errors | No — mask as internal error | Yes |
| Third-party API errors | No — mask | Yes |
| Unexpected exceptions | No — mask | Yes + alert |
# Chapter 13: Pagination

## Why Pagination?

Never return all records at once — always paginate. Without pagination:
- Responses can be megabytes of JSON
- Database queries become full table scans
- Memory usage grows unbounded
- UI becomes unusable

---

## Offset Pagination

The simplest approach — skip N records, take M.

### Schema:
```graphql
type Query {
  posts(limit: Int = 20, offset: Int = 0): PostList!
}

type PostList {
  items: [Post!]!
  total: Int!
  hasMore: Boolean!
}
```

### Resolver:
```javascript
const resolvers = {
  Query: {
    posts: async (parent, { limit = 20, offset = 0 }, context) => {
      const [items, total] = await Promise.all([
        context.db.posts.findAll({ limit, offset, orderBy: 'createdAt', order: 'DESC' }),
        context.db.posts.count(),
      ]);

      return {
        items,
        total,
        hasMore: offset + items.length < total,
      };
    },
  },
};
```

### Client:
```graphql
query GetPosts($limit: Int!, $offset: Int!) {
  posts(limit: $limit, offset: $offset) {
    items {
      id
      title
      createdAt
    }
    total
    hasMore
  }
}
```

### Problems with offset:
- **Skip problem:** If new items are added, items shift — page 2 might show items already on page 1
- **Performance:** `OFFSET 10000` requires DB to read and discard 10,000 rows
- **Consistency:** Results change between page loads if data changes

---

## Cursor-Based Pagination (Recommended)

Cursors point to a specific position in the result set. Stable even with inserts/deletes.

### Schema (Relay Connection Spec):
```graphql
type Query {
  posts(
    first: Int      # forward: take N items
    after: String   # forward: after this cursor
    last: Int       # backward: take N items from the end
    before: String  # backward: before this cursor
  ): PostConnection!
}

type PostConnection {
  edges: [PostEdge!]!
  pageInfo: PageInfo!
  totalCount: Int!
}

type PostEdge {
  node: Post!
  cursor: String!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}
```

### Cursor Implementation:

```javascript
// Encode/decode cursor (base64 of ID or timestamp)
const encodeCursor = (value) => Buffer.from(String(value)).toString('base64');
const decodeCursor = (cursor) => Buffer.from(cursor, 'base64').toString('utf-8');

const resolvers = {
  Query: {
    posts: async (parent, { first = 20, after, last, before }, context) => {
      let query = context.db.posts.query().orderBy('created_at', 'DESC');

      if (after) {
        const cursor = decodeCursor(after);
        query = query.where('created_at', '<', cursor); // assuming cursor is a timestamp
      }

      if (before) {
        const cursor = decodeCursor(before);
        query = query.where('created_at', '>', cursor);
      }

      const limit = first || last || 20;
      const posts = await query.limit(limit + 1); // fetch 1 extra to detect hasNextPage

      const hasNextPage = posts.length > limit;
      const hasPreviousPage = !!after;
      const items = hasNextPage ? posts.slice(0, limit) : posts;

      const edges = items.map(post => ({
        node: post,
        cursor: encodeCursor(post.createdAt.toISOString()),
      }));

      return {
        edges,
        totalCount: await context.db.posts.count(),
        pageInfo: {
          hasNextPage,
          hasPreviousPage,
          startCursor: edges[0]?.cursor ?? null,
          endCursor: edges[edges.length - 1]?.cursor ?? null,
        },
      };
    },
  },
};
```

### Client Usage:
```graphql
# First page
query GetFirstPage {
  posts(first: 10) {
    edges {
      cursor
      node { id title createdAt }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
  }
}

# Next page
query GetNextPage($after: String!) {
  posts(first: 10, after: $after) {
    edges {
      cursor
      node { id title createdAt }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

```javascript
// Apollo Client: fetchMore
const { data, fetchMore } = useQuery(GET_POSTS, { variables: { first: 10 } });

const loadMore = () => {
  fetchMore({
    variables: {
      first: 10,
      after: data.posts.pageInfo.endCursor,
    },
    updateQuery: (prev, { fetchMoreResult }) => ({
      posts: {
        ...fetchMoreResult.posts,
        edges: [...prev.posts.edges, ...fetchMoreResult.posts.edges],
      },
    }),
  });
};
```

---

## ID-Based Cursor (Simpler)

Use the row ID as the cursor — works well for tables with sequential IDs:

```javascript
const resolvers = {
  Query: {
    posts: async (parent, { first = 20, after }, context) => {
      const afterId = after ? parseInt(decodeCursor(after)) : null;

      const posts = await context.db.query(`
        SELECT * FROM posts
        WHERE ($1::bigint IS NULL OR id < $1)
        ORDER BY id DESC
        LIMIT $2
      `, [afterId, first + 1]);

      const hasNextPage = posts.length > first;
      const items = hasNextPage ? posts.slice(0, first) : posts;

      return {
        edges: items.map(p => ({ node: p, cursor: encodeCursor(p.id) })),
        pageInfo: {
          hasNextPage,
          hasPreviousPage: !!after,
          startCursor: items[0] ? encodeCursor(items[0].id) : null,
          endCursor: items[items.length-1] ? encodeCursor(items[items.length-1].id) : null,
        },
        totalCount: await context.db.posts.count(),
      };
    },
  },
};
```

---

## Simple Page Number Pagination

Sometimes page numbers are better for UX (especially for search results):

```graphql
type Query {
  searchPosts(query: String!, page: Int = 1, perPage: Int = 20): SearchResult!
}

type SearchResult {
  posts: [Post!]!
  currentPage: Int!
  totalPages: Int!
  totalCount: Int!
}
```

```javascript
const resolvers = {
  Query: {
    searchPosts: async (parent, { query, page = 1, perPage = 20 }, context) => {
      const offset = (page - 1) * perPage;
      const [posts, totalCount] = await Promise.all([
        context.db.posts.search(query, { limit: perPage, offset }),
        context.db.posts.searchCount(query),
      ]);

      return {
        posts,
        currentPage: page,
        totalPages: Math.ceil(totalCount / perPage),
        totalCount,
      };
    },
  },
};
```

---

## Pagination Comparison

| Method | Stable? | Performance | UX | Use When |
|---|---|---|---|---|
| Offset | No | Poor for large offsets | Simple | Admin lists, small datasets |
| Cursor | Yes | Excellent | Load more, infinite scroll | Social feeds, real-time data |
| Page number | No | Poor for large pages | Familiar | Search results, tables |

---

## Cursor Pagination Best Practices

- Always encode cursors (base64) — clients should treat them as opaque strings
- Don't expose implementation details in cursor values (ID, timestamp)
- Support both `first/after` (forward) and `last/before` (backward)
- Always return `pageInfo.hasNextPage` accurately
- Limit maximum `first` value to prevent abuse: `const safeLimit = Math.min(first ?? 20, 100)`
- Include `totalCount` carefully — it can be expensive; cache it or skip if not needed
# Chapter 14: File Uploads

## The GraphQL Multipart Spec

GraphQL doesn't natively support file uploads — files are binary data while GraphQL uses JSON. The **GraphQL multipart request spec** (by Jayden Seric) is the community standard for handling uploads.

The request is a `multipart/form-data` POST with three parts:
1. `operations` — the GraphQL query/mutation JSON
2. `map` — maps file variable paths to the multipart parts
3. File parts — the actual file data

---

## Server Setup (Apollo Server + graphql-upload)

```bash
npm install graphql-upload
```

```javascript
import { ApolloServer } from '@apollo/server';
import { graphqlUploadExpress } from 'graphql-upload';
import express from 'express';

const app = express();

// Must be before Apollo middleware
app.use(graphqlUploadExpress({
  maxFileSize: 10_000_000,   // 10 MB
  maxFiles: 5,
}));

// GraphQL Upload provides the Upload scalar
import { GraphQLUpload } from 'graphql-upload';

const resolvers = {
  Upload: GraphQLUpload,  // register the scalar
};
```

---

## Schema with Upload

```graphql
scalar Upload

type Mutation {
  uploadAvatar(file: Upload!): UploadResult!
  uploadFiles(files: [Upload!]!): [UploadResult!]!

  createPost(input: CreatePostWithImageInput!): Post!
}

input CreatePostWithImageInput {
  title: String!
  body: String!
  coverImage: Upload
}

type UploadResult {
  filename: String!
  mimetype: String!
  url: String!
  size: Int!
}
```

---

## Upload Resolver

The `Upload` scalar resolves to a Promise containing:
- `filename` — original file name
- `mimetype` — MIME type
- `encoding` — file encoding
- `createReadStream()` — function to create a readable stream

```javascript
import { createWriteStream, mkdir } from 'fs';
import { pipeline } from 'stream/promises';
import path from 'path';
import crypto from 'crypto';

const resolvers = {
  Mutation: {
    uploadAvatar: async (parent, { file }, context) => {
      requireAuth(context);

      const { filename, mimetype, createReadStream } = await file;

      // Validate file type
      const allowedTypes = ['image/jpeg', 'image/png', 'image/webp'];
      if (!allowedTypes.includes(mimetype)) {
        throw new GraphQLError('Invalid file type. Only JPEG, PNG, WebP allowed.', {
          extensions: { code: 'BAD_USER_INPUT' },
        });
      }

      // Generate unique filename
      const ext = path.extname(filename);
      const uniqueName = `${crypto.randomUUID()}${ext}`;
      const uploadDir = './uploads/avatars';
      const filePath = path.join(uploadDir, uniqueName);

      // Save to disk
      const stream = createReadStream();
      await pipeline(stream, createWriteStream(filePath));

      // Update user's avatar in DB
      const url = `/uploads/avatars/${uniqueName}`;
      await context.db.users.update(context.currentUser.id, { avatarUrl: url });

      return { filename: uniqueName, mimetype, url, size: 0 };
    },

    uploadFiles: async (parent, { files }, context) => {
      const uploads = await Promise.all(files.map(async (filePromise) => {
        const { filename, mimetype, createReadStream } = await filePromise;
        // ... save each file
        return { filename, mimetype, url: `/uploads/${filename}`, size: 0 };
      }));
      return uploads;
    },
  },
};
```

---

## Uploading to S3 (AWS)

In production, save to object storage — not local disk:

```javascript
import { S3Client, PutObjectCommand } from '@aws-sdk/client-s3';

const s3 = new S3Client({ region: process.env.AWS_REGION });

async function uploadToS3(stream, filename, mimetype) {
  const key = `uploads/${Date.now()}-${filename}`;

  await s3.send(new PutObjectCommand({
    Bucket: process.env.S3_BUCKET,
    Key: key,
    Body: stream,
    ContentType: mimetype,
    ACL: 'public-read',
  }));

  return `https://${process.env.S3_BUCKET}.s3.amazonaws.com/${key}`;
}

const resolvers = {
  Mutation: {
    uploadAvatar: async (parent, { file }, context) => {
      requireAuth(context);
      const { filename, mimetype, createReadStream } = await file;

      // Validate
      if (!['image/jpeg', 'image/png', 'image/webp'].includes(mimetype)) {
        throw new GraphQLError('Invalid file type');
      }

      const stream = createReadStream();
      const url = await uploadToS3(stream, filename, mimetype);

      await context.db.users.update(context.currentUser.id, { avatarUrl: url });

      return { filename, mimetype, url, size: 0 };
    },
  },
};
```

---

## Image Processing (Sharp)

```javascript
import sharp from 'sharp';
import { Readable } from 'stream';

async function processAndUploadImage(createReadStream, filename) {
  const chunks = [];
  for await (const chunk of createReadStream()) {
    chunks.push(chunk);
  }
  const buffer = Buffer.concat(chunks);

  // Resize and convert to WebP
  const processed = await sharp(buffer)
    .resize(800, 800, { fit: 'inside', withoutEnlargement: true })
    .webp({ quality: 85 })
    .toBuffer();

  const processedStream = Readable.from(processed);
  const newFilename = filename.replace(/\.[^.]+$/, '.webp');

  const url = await uploadToS3(processedStream, newFilename, 'image/webp');
  return url;
}
```

---

## Client-Side Upload

### With Apollo Client (React):

```javascript
import { useMutation, gql } from '@apollo/client';

const UPLOAD_AVATAR = gql`
  mutation UploadAvatar($file: Upload!) {
    uploadAvatar(file: $file) {
      url
      filename
    }
  }
`;

function AvatarUpload() {
  const [uploadAvatar, { loading, error }] = useMutation(UPLOAD_AVATAR);

  const handleChange = async (event) => {
    const file = event.target.files[0];
    if (!file) return;

    try {
      const { data } = await uploadAvatar({
        variables: { file },
      });
      console.log('Uploaded:', data.uploadAvatar.url);
    } catch (err) {
      console.error('Upload failed:', err);
    }
  };

  return (
    <div>
      <input type="file" accept="image/*" onChange={handleChange} disabled={loading} />
      {loading && <p>Uploading...</p>}
      {error && <p>Error: {error.message}</p>}
    </div>
  );
}
```

Apollo Client automatically formats the multipart request — you just pass the `File` object as a variable.

### Raw fetch (without Apollo):

```javascript
async function uploadFile(file, query, variables = {}) {
  const formData = new FormData();

  formData.append('operations', JSON.stringify({
    query,
    variables: { ...variables, file: null },
  }));

  formData.append('map', JSON.stringify({ '0': ['variables.file'] }));
  formData.append('0', file);

  const response = await fetch('/graphql', {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
    body: formData,
  });

  return response.json();
}

// Usage
uploadFile(
  fileInputElement.files[0],
  `mutation UploadAvatar($file: Upload!) { uploadAvatar(file: $file) { url } }`
);
```

---

## Alternative: Pre-signed URLs (Better for Large Files)

For large files, generate a pre-signed S3 URL and let the client upload directly:

```graphql
type Mutation {
  getUploadUrl(filename: String!, contentType: String!): UploadUrlPayload!
  confirmUpload(key: String!): Post!
}

type UploadUrlPayload {
  uploadUrl: String!   # pre-signed S3 PUT URL
  key: String!         # S3 key to confirm after upload
}
```

```javascript
import { getSignedUrl } from '@aws-sdk/s3-request-presigner';
import { PutObjectCommand } from '@aws-sdk/client-s3';

const resolvers = {
  Mutation: {
    getUploadUrl: async (parent, { filename, contentType }, context) => {
      requireAuth(context);
      const key = `uploads/${context.currentUser.id}/${Date.now()}-${filename}`;

      const url = await getSignedUrl(s3, new PutObjectCommand({
        Bucket: process.env.S3_BUCKET,
        Key: key,
        ContentType: contentType,
      }), { expiresIn: 300 }); // URL valid for 5 minutes

      return { uploadUrl: url, key };
    },
  },
};
```

Client: PUT file directly to S3, then call `confirmUpload(key)` to save the record.

---

## Security Considerations

- Always validate MIME type server-side (don't trust client)
- Scan uploaded files for malware (ClamAV, AWS Macie)
- Set maximum file size limits
- Never serve user-uploaded files from the same origin as your app (use S3/CDN)
- Set Content-Disposition and Content-Type headers when serving
- Use random filenames — never trust user-provided filenames
- Store files outside the web root if serving locally
# Chapter 15: GraphQL with React (Apollo Client)

## Apollo Client Setup

Apollo Client is the most popular GraphQL client for React.

```bash
npm install @apollo/client graphql
```

```javascript
// src/apollo.js
import { ApolloClient, InMemoryCache, ApolloProvider, createHttpLink } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';

const httpLink = createHttpLink({
  uri: 'http://localhost:4000/graphql',
});

// Add auth token to every request
const authLink = setContext((_, { headers }) => {
  const token = localStorage.getItem('token');
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : '',
    },
  };
});

export const client = new ApolloClient({
  link: authLink.concat(httpLink),
  cache: new InMemoryCache(),
  defaultOptions: {
    watchQuery: { errorPolicy: 'all' },
    query: { errorPolicy: 'all' },
  },
});

// Wrap your app
// src/index.js
root.render(
  <ApolloProvider client={client}>
    <App />
  </ApolloProvider>
);
```

---

## useQuery

```javascript
import { useQuery, gql } from '@apollo/client';

const GET_USERS = gql`
  query GetUsers($limit: Int!, $role: UserRole) {
    users(limit: $limit, role: $role) {
      id
      name
      email
      role
      avatar(size: 64)
    }
  }
`;

function UserList() {
  const { data, loading, error, refetch } = useQuery(GET_USERS, {
    variables: { limit: 20, role: 'USER' },
    fetchPolicy: 'cache-and-network',   // use cache but also refetch in background
    pollInterval: 30000,                // refetch every 30 seconds
    skip: !isLoggedIn,                  // skip query if not logged in
    onCompleted: (data) => console.log('Loaded', data.users.length, 'users'),
    onError: (error) => console.error('Failed to load users', error),
  });

  if (loading && !data) return <Spinner />;
  if (error) return <ErrorMessage error={error} />;

  return (
    <div>
      <button onClick={() => refetch()}>Refresh</button>
      {data.users.map(user => (
        <UserCard key={user.id} user={user} />
      ))}
    </div>
  );
}
```

### Fetch Policies:

| Policy | Cache | Network | When to Use |
|---|---|---|---|
| `cache-first` | ✓ always | Only on cache miss | Default — best performance |
| `cache-and-network` | ✓ immediately | Always refetches | Stale-while-revalidate |
| `network-only` | Writes to cache | Always fetches | Real-time accuracy needed |
| `no-cache` | Never | Always fetches | Sensitive data |
| `cache-only` | ✓ always | Never | Offline mode |

---

## useMutation

```javascript
import { useMutation, gql } from '@apollo/client';

const CREATE_POST = gql`
  mutation CreatePost($input: CreatePostInput!) {
    createPost(input: $input) {
      post {
        id
        title
        status
        createdAt
      }
      errors {
        field
        message
      }
    }
  }
`;

function CreatePostForm() {
  const [title, setTitle] = useState('');
  const [body, setBody] = useState('');
  const [fieldErrors, setFieldErrors] = useState({});

  const [createPost, { loading, error }] = useMutation(CREATE_POST, {
    // Update cache after mutation
    update(cache, { data: { createPost } }) {
      if (!createPost.post) return;
      cache.modify({
        fields: {
          posts(existing = []) {
            const newPostRef = cache.writeFragment({
              data: createPost.post,
              fragment: gql`fragment NewPost on Post { id title status createdAt }`,
            });
            return [newPostRef, ...existing];
          },
        },
      });
    },
    onCompleted: ({ createPost }) => {
      if (createPost.errors.length > 0) {
        const errors = {};
        createPost.errors.forEach(e => { errors[e.field] = e.message; });
        setFieldErrors(errors);
      } else {
        setTitle('');
        setBody('');
        alert('Post created!');
      }
    },
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    setFieldErrors({});
    createPost({ variables: { input: { title, body } } });
  };

  return (
    <form onSubmit={handleSubmit}>
      <input value={title} onChange={e => setTitle(e.target.value)} />
      {fieldErrors.title && <span>{fieldErrors.title}</span>}
      <textarea value={body} onChange={e => setBody(e.target.value)} />
      {fieldErrors.body && <span>{fieldErrors.body}</span>}
      <button type="submit" disabled={loading}>
        {loading ? 'Creating...' : 'Create Post'}
      </button>
      {error && <div>Network error: {error.message}</div>}
    </form>
  );
}
```

---

## Apollo Cache

Apollo Client's InMemoryCache is a normalized cache — it stores objects by type + ID.

```javascript
// Cache normalization
// Two queries returning the same user are merged in cache:
// Query 1: { user(id: "1") { id name } }
// Query 2: { me { id name email } }  (same user)
// Cache stores: User:1 → { id: "1", name: "Alice", email: "alice@..." }

const client = new ApolloClient({
  cache: new InMemoryCache({
    typePolicies: {
      Query: {
        fields: {
          // Custom field policy for pagination
          posts: {
            keyArgs: ['filter'],  // cache by filter only (not pagination args)
            merge(existing = [], incoming, { args }) {
              // Merge paginated results
              const merged = existing.slice();
              incoming.edges.forEach((edge, index) => {
                merged[args.offset + index] = edge;
              });
              return merged;
            },
          },
        },
      },
      User: {
        // Custom ID field (default is 'id')
        keyFields: ['email'],  // use email as unique identifier
      },
    },
  }),
});
```

---

## Cache Updates After Mutations

### Refetch queries (simplest):
```javascript
const [deletePost] = useMutation(DELETE_POST, {
  refetchQueries: ['GetPosts'],   // refetch by operation name
  // Or with variables:
  refetchQueries: [{ query: GET_POSTS, variables: { limit: 20 } }],
});
```

### Update cache manually (best performance):
```javascript
const [deletePost] = useMutation(DELETE_POST, {
  update(cache, { data }, { variables }) {
    cache.modify({
      fields: {
        posts(existingPosts, { readField }) {
          return existingPosts.filter(
            postRef => readField('id', postRef) !== variables.id
          );
        },
      },
    });
    // Also remove from normalized store
    cache.evict({ id: cache.identify({ __typename: 'Post', id: variables.id }) });
    cache.gc();
  },
});
```

---

## useSubscription

```javascript
import { useSubscription, gql } from '@apollo/client';

const POST_ADDED = gql`
  subscription OnPostAdded {
    postAdded {
      id
      title
      author { id name avatar }
      createdAt
    }
  }
`;

function LiveFeed() {
  const [posts, setPosts] = useState([]);

  // Load initial posts
  const { loading } = useQuery(GET_POSTS, {
    onCompleted: (data) => setPosts(data.posts),
  });

  // Subscribe to new posts
  useSubscription(POST_ADDED, {
    onData: ({ data }) => {
      setPosts(prev => [data.data.postAdded, ...prev]);
    },
  });

  if (loading) return <Spinner />;
  return <PostList posts={posts} />;
}
```

---

## Pagination with Apollo

```javascript
const GET_POSTS = gql`
  query GetPosts($first: Int!, $after: String) {
    posts(first: $first, after: $after) {
      edges { cursor node { id title excerpt } }
      pageInfo { hasNextPage endCursor }
    }
  }
`;

function InfinitePostList() {
  const { data, loading, fetchMore } = useQuery(GET_POSTS, {
    variables: { first: 10 },
  });

  const loadMore = () => {
    fetchMore({
      variables: {
        first: 10,
        after: data.posts.pageInfo.endCursor,
      },
      updateQuery: (prev, { fetchMoreResult }) => {
        if (!fetchMoreResult) return prev;
        return {
          posts: {
            ...fetchMoreResult.posts,
            edges: [...prev.posts.edges, ...fetchMoreResult.posts.edges],
          },
        };
      },
    });
  };

  if (loading && !data) return <Spinner />;

  return (
    <div>
      {data.posts.edges.map(({ node }) => (
        <PostCard key={node.id} post={node} />
      ))}
      {data.posts.pageInfo.hasNextPage && (
        <button onClick={loadMore} disabled={loading}>
          Load More
        </button>
      )}
    </div>
  );
}
```

---

## Code Generation with GraphQL Code Generator

Generate TypeScript types for your queries:

```bash
npm install -D @graphql-codegen/cli @graphql-codegen/client-preset
```

```yaml
# codegen.ts
import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  schema: 'http://localhost:4000/graphql',
  documents: ['src/**/*.tsx', 'src/**/*.ts'],
  generates: {
    './src/generated/': {
      preset: 'client',
      plugins: [],
    },
  },
};
export default config;
```

```bash
npx graphql-codegen --config codegen.ts
```

Use generated types:
```typescript
import { useQuery } from '@apollo/client';
import { graphql } from './generated';

const GET_USER = graphql(`
  query GetUser($id: ID!) {
    user(id: $id) {
      id
      name
      email
    }
  }
`);

function UserProfile({ id }: { id: string }) {
  const { data } = useQuery(GET_USER, { variables: { id } });
  // data is fully typed: data?.user?.name is string | undefined
  return <h1>{data?.user?.name}</h1>;
}
```
# Chapter 16: Server Implementations

## Overview

Multiple frameworks and languages support GraphQL. The best choice depends on your stack.

| Server | Language | Style | Best For |
|---|---|---|---|
| Apollo Server | Node.js | Feature-rich, enterprise | Large teams, rich ecosystem |
| GraphQL Yoga | Node.js | Lightweight, standards-based | Modern APIs, edge runtime |
| Pothos | Node.js | Code-first, TypeScript-native | Type safety without schema writing |
| NestJS + GraphQL | Node.js | Framework, decorators | Enterprise Node.js |
| Spring Boot GraphQL | Java/Kotlin | Annotation-based | Java ecosystem |
| Strawberry | Python | Decorator-based, type hints | Python, FastAPI-like |
| graphql-ruby | Ruby | Ruby-native | Rails apps |
| gqlgen | Go | Code-first, type-safe | Go microservices |

---

## Apollo Server (Node.js)

```bash
npm install @apollo/server graphql
```

```javascript
import { ApolloServer } from '@apollo/server';
import { startStandaloneServer } from '@apollo/server/standalone';

const typeDefs = `#graphql
  type Query {
    hello: String!
    user(id: ID!): User
  }

  type User {
    id: ID!
    name: String!
  }
`;

const resolvers = {
  Query: {
    hello: () => 'Hello, World!',
    user: (_, { id }, context) => context.db.users.findById(id),
  },
};

const server = new ApolloServer({
  typeDefs,
  resolvers,
  plugins: [
    ApolloServerPluginLandingPageLocalDefault(), // GraphQL sandbox in dev
  ],
  formatError: (formattedError, error) => {
    // Sanitize errors in production
    if (process.env.NODE_ENV === 'production' && !formattedError.extensions?.code) {
      return { message: 'Internal server error', extensions: { code: 'INTERNAL_SERVER_ERROR' } };
    }
    return formattedError;
  },
});

const { url } = await startStandaloneServer(server, {
  listen: { port: 4000 },
  context: async ({ req }) => ({
    user: await getUser(req.headers.authorization),
    db,
  }),
});

console.log(`Server running at ${url}`);
```

### With Express (more control):
```javascript
import { ApolloServer } from '@apollo/server';
import { expressMiddleware } from '@apollo/server/express4';
import express from 'express';
import cors from 'cors';

const app = express();
app.use(express.json());
app.use(cors());

const server = new ApolloServer({ typeDefs, resolvers });
await server.start();

app.use('/graphql', expressMiddleware(server, {
  context: async ({ req }) => ({
    user: await getUser(req.headers.authorization),
    db,
  }),
}));

app.listen(4000, () => console.log('Server on port 4000'));
```

---

## GraphQL Yoga (Node.js — Lightweight)

```bash
npm install graphql-yoga graphql
```

```javascript
import { createSchema, createYoga } from 'graphql-yoga';
import { createServer } from 'http';

const yoga = createYoga({
  schema: createSchema({
    typeDefs: `
      type Query { hello: String! }
    `,
    resolvers: {
      Query: { hello: () => 'Hello, Yoga!' },
    },
  }),
  context: ({ request }) => ({
    user: getUser(request.headers.get('authorization')),
    db,
  }),
});

createServer(yoga).listen(4000);
```

Yoga works on Node.js, Cloudflare Workers, Deno, Bun — highly portable.

---

## Pothos (Code-First TypeScript)

Pothos builds the schema in TypeScript — no SDL writing, full type safety:

```bash
npm install @pothos/core graphql
```

```typescript
import SchemaBuilder from '@pothos/core';

const builder = new SchemaBuilder<{
  Context: { user: User | null; db: DB };
}>({});

// Define types with code (no SDL)
const UserType = builder.objectType('User', {
  fields: (t) => ({
    id: t.exposeID('id'),
    name: t.exposeString('name'),
    email: t.exposeString('email'),
    posts: t.field({
      type: [PostType],
      resolve: (user, _, context) => context.db.posts.findByAuthor(user.id),
    }),
  }),
});

const PostType = builder.objectType('Post', {
  fields: (t) => ({
    id: t.exposeID('id'),
    title: t.exposeString('title'),
    author: t.field({
      type: UserType,
      resolve: (post, _, context) => context.db.users.findById(post.authorId),
    }),
  }),
});

builder.queryType({
  fields: (t) => ({
    user: t.field({
      type: UserType,
      nullable: true,
      args: { id: t.arg.id({ required: true }) },
      resolve: (root, { id }, context) => context.db.users.findById(id),
    }),
  }),
});

builder.mutationType({
  fields: (t) => ({
    createPost: t.field({
      type: PostType,
      args: { title: t.arg.string({ required: true }) },
      resolve: (root, { title }, context) => {
        if (!context.user) throw new Error('Unauthorized');
        return context.db.posts.create({ title, authorId: context.user.id });
      },
    }),
  }),
});

export const schema = builder.toSchema();
```

---

## NestJS with GraphQL

```bash
npm install @nestjs/graphql @nestjs/apollo @apollo/server graphql
```

```typescript
// app.module.ts
@Module({
  imports: [
    GraphQLModule.forRoot<ApolloDriverConfig>({
      driver: ApolloDriver,
      autoSchemaFile: 'schema.gql',  // code-first: auto-generates SDL
      sortSchema: true,
      context: ({ req }) => ({ req }),
    }),
    UsersModule,
    PostsModule,
  ],
})
export class AppModule {}

// user.resolver.ts
@Resolver(() => User)
export class UserResolver {
  constructor(
    private usersService: UsersService,
    private postsService: PostsService,
  ) {}

  @Query(() => [User])
  async users(): Promise<User[]> {
    return this.usersService.findAll();
  }

  @Query(() => User, { nullable: true })
  async user(@Args('id', { type: () => ID }) id: string): Promise<User | null> {
    return this.usersService.findById(id);
  }

  @ResolveField(() => [Post])
  async posts(@Parent() user: User): Promise<Post[]> {
    return this.postsService.findByAuthor(user.id);
  }

  @Mutation(() => User)
  @UseGuards(GqlAuthGuard)
  async createUser(
    @Args('input') input: CreateUserInput,
  ): Promise<User> {
    return this.usersService.create(input);
  }
}

// user.entity.ts
@ObjectType()
export class User {
  @Field(() => ID)
  id: string;

  @Field()
  name: string;

  @Field()
  email: string;

  @Field(() => [Post])
  posts: Post[];
}

// create-user.input.ts
@InputType()
export class CreateUserInput {
  @Field()
  @IsEmail()
  email: string;

  @Field()
  @MinLength(2)
  name: string;
}
```

---

## Spring Boot with GraphQL (Java)

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-graphql</artifactId>
</dependency>
```

```graphql
# src/main/resources/graphql/schema.graphqls
type Query {
  user(id: ID!): User
  users: [User!]!
}

type Mutation {
  createUser(input: CreateUserInput!): User!
}

type User {
  id: ID!
  name: String!
  posts: [Post!]!
}
```

```java
@Controller
public class UserController {

    private final UserService userService;

    @QueryMapping
    public Optional<User> user(@Argument Long id) {
        return userService.findById(id);
    }

    @QueryMapping
    public List<User> users() {
        return userService.findAll();
    }

    @MutationMapping
    public User createUser(@Argument @Valid CreateUserInput input) {
        return userService.create(input);
    }

    @SchemaMapping
    public List<Post> posts(User user) {
        return postService.findByAuthor(user.getId());
    }
}
```

---

## gqlgen (Go)

```bash
go get github.com/99designs/gqlgen
go run github.com/99designs/gqlgen generate
```

```go
// resolver.go
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
    return r.db.FindUser(ctx, id)
}

func (r *userResolver) Posts(ctx context.Context, obj *model.User) ([]*model.Post, error) {
    return r.db.FindPostsByAuthor(ctx, obj.ID)
}
```

Go's strong typing makes gqlgen generate type-safe resolver interfaces from your schema.

---

## Choosing Your Server

| Need | Recommendation |
|---|---|
| Node.js, flexibility | Apollo Server |
| Node.js, edge/serverless | GraphQL Yoga |
| TypeScript, type safety | Pothos + Yoga |
| Node.js enterprise, decorators | NestJS |
| Java/Kotlin | Spring Boot GraphQL |
| Python | Strawberry or Ariadne |
| Go | gqlgen |
| Ruby on Rails | graphql-ruby |
# Chapter 17: Testing GraphQL

## Testing Strategy

```
Unit Tests        → Resolver functions (business logic)
Integration Tests → Resolver with real database (via test DB or Testcontainers)
E2E Tests         → Full HTTP request through GraphQL server
Schema Tests      → Schema validity, breaking changes
```

---

## Unit Testing Resolvers

Resolvers are just functions — test them directly without a GraphQL server.

```typescript
// resolvers/user.resolver.ts
export const userResolvers = {
  Query: {
    user: async (_, { id }, { db, user: currentUser }) => {
      if (!currentUser) throw new AuthenticationError('Not authenticated');
      return db.users.findById(id);
    },
    users: async (_, { limit = 20, offset = 0 }, { db }) => {
      return db.users.findMany({ limit, offset });
    },
  },
  Mutation: {
    createUser: async (_, { input }, { db, user: currentUser }) => {
      if (!currentUser?.isAdmin) throw new ForbiddenError('Admins only');
      return db.users.create(input);
    },
  },
};
```

```typescript
// resolvers/user.resolver.test.ts
import { userResolvers } from './user.resolver';
import { AuthenticationError, ForbiddenError } from '../errors';

describe('user resolvers', () => {
  const mockDb = {
    users: {
      findById: jest.fn(),
      findMany: jest.fn(),
      create: jest.fn(),
    },
  };

  beforeEach(() => jest.clearAllMocks());

  describe('Query.user', () => {
    it('throws AuthenticationError when not authenticated', async () => {
      const context = { db: mockDb, user: null };
      await expect(
        userResolvers.Query.user({}, { id: '1' }, context)
      ).rejects.toThrow(AuthenticationError);
    });

    it('returns user when authenticated', async () => {
      const fakeUser = { id: '1', name: 'Alice', email: 'alice@example.com' };
      mockDb.users.findById.mockResolvedValue(fakeUser);

      const context = { db: mockDb, user: { id: '42' } };
      const result = await userResolvers.Query.user({}, { id: '1' }, context);

      expect(result).toEqual(fakeUser);
      expect(mockDb.users.findById).toHaveBeenCalledWith('1');
    });
  });

  describe('Mutation.createUser', () => {
    it('throws ForbiddenError when not admin', async () => {
      const context = { db: mockDb, user: { id: '1', isAdmin: false } };
      await expect(
        userResolvers.Mutation.createUser({}, { input: { name: 'Bob' } }, context)
      ).rejects.toThrow(ForbiddenError);
    });

    it('creates user when admin', async () => {
      const newUser = { id: '99', name: 'Bob', email: 'bob@example.com' };
      mockDb.users.create.mockResolvedValue(newUser);

      const context = { db: mockDb, user: { id: '1', isAdmin: true } };
      const result = await userResolvers.Mutation.createUser(
        {}, { input: { name: 'Bob', email: 'bob@example.com' } }, context
      );

      expect(result).toEqual(newUser);
    });
  });
});
```

---

## Integration Testing with Apollo Server

Test the full GraphQL pipeline including schema + resolvers.

```bash
npm install -D @apollo/server supertest jest @types/jest ts-jest
```

### Test server setup:

```typescript
// test/helpers/server.ts
import { ApolloServer } from '@apollo/server';
import { typeDefs } from '../../src/schema';
import { resolvers } from '../../src/resolvers';
import { createTestContext } from './context';

export function createTestServer() {
  return new ApolloServer({ typeDefs, resolvers });
}

export async function executeOperation(
  server: ApolloServer,
  query: string,
  variables?: Record<string, unknown>,
  contextOverrides?: Partial<Context>
) {
  const response = await server.executeOperation(
    { query, variables },
    { contextValue: createTestContext(contextOverrides) }
  );
  return response.body.kind === 'single' ? response.body.singleResult : null;
}
```

```typescript
// test/helpers/context.ts
export function createTestContext(overrides: Partial<Context> = {}): Context {
  return {
    user: null,
    db: testDb,     // real test database connection
    ...overrides,
  };
}
```

### Writing integration tests:

```typescript
// test/users.test.ts
import { ApolloServer } from '@apollo/server';
import { createTestServer, executeOperation } from './helpers/server';

describe('User GraphQL Operations', () => {
  let server: ApolloServer;

  beforeAll(async () => {
    server = createTestServer();
    await server.start();
    await seedTestDatabase();
  });

  afterAll(async () => {
    await server.stop();
    await cleanTestDatabase();
  });

  describe('query { users }', () => {
    const USERS_QUERY = `
      query GetUsers($limit: Int) {
        users(limit: $limit) {
          id
          name
          email
        }
      }
    `;

    it('returns list of users', async () => {
      const result = await executeOperation(server, USERS_QUERY, { limit: 5 });

      expect(result?.errors).toBeUndefined();
      expect(result?.data?.users).toHaveLength(5);
      expect(result?.data?.users[0]).toMatchObject({
        id: expect.any(String),
        name: expect.any(String),
        email: expect.any(String),
      });
    });

    it('requires authentication', async () => {
      const result = await executeOperation(
        server,
        USERS_QUERY,
        { limit: 5 },
        { user: null }   // unauthenticated
      );

      expect(result?.errors).toBeDefined();
      expect(result?.errors?.[0].extensions?.code).toBe('UNAUTHENTICATED');
    });
  });

  describe('mutation { createUser }', () => {
    const CREATE_USER = `
      mutation CreateUser($input: CreateUserInput!) {
        createUser(input: $input) {
          user {
            id
            name
            email
          }
          errors {
            field
            message
          }
        }
      }
    `;

    it('creates a user when authenticated as admin', async () => {
      const result = await executeOperation(
        server,
        CREATE_USER,
        { input: { name: 'New User', email: 'new@example.com' } },
        { user: { id: '1', isAdmin: true } }
      );

      expect(result?.errors).toBeUndefined();
      expect(result?.data?.createUser.user).toMatchObject({
        name: 'New User',
        email: 'new@example.com',
      });
      expect(result?.data?.createUser.errors).toHaveLength(0);
    });

    it('returns validation errors for invalid email', async () => {
      const result = await executeOperation(
        server,
        CREATE_USER,
        { input: { name: 'Bob', email: 'not-an-email' } },
        { user: { id: '1', isAdmin: true } }
      );

      expect(result?.data?.createUser.user).toBeNull();
      expect(result?.data?.createUser.errors).toContainEqual(
        expect.objectContaining({ field: 'email' })
      );
    });
  });
});
```

---

## HTTP-level Testing with Supertest

For testing the full HTTP layer (middleware, auth headers, CORS):

```typescript
// test/http.test.ts
import supertest from 'supertest';
import { buildApp } from '../src/app';   // your Express/Fastify app

describe('GraphQL HTTP endpoint', () => {
  let app: Express;

  beforeAll(async () => {
    app = await buildApp();
  });

  it('returns 200 for valid query', async () => {
    const response = await supertest(app)
      .post('/graphql')
      .set('Content-Type', 'application/json')
      .set('Authorization', `Bearer ${testToken}`)
      .send({
        query: `{ users { id name } }`,
      });

    expect(response.status).toBe(200);
    expect(response.body.errors).toBeUndefined();
    expect(response.body.data.users).toBeDefined();
  });

  it('handles introspection', async () => {
    const response = await supertest(app)
      .post('/graphql')
      .set('Content-Type', 'application/json')
      .send({ query: '{ __schema { types { name } } }' });

    expect(response.status).toBe(200);
    expect(response.body.data.__schema).toBeDefined();
  });

  it('returns 400 for malformed JSON', async () => {
    const response = await supertest(app)
      .post('/graphql')
      .set('Content-Type', 'application/json')
      .send('{ invalid json }');

    expect(response.status).toBe(400);
  });
});
```

---

## Snapshot Testing for Schema

Catch breaking changes to your schema:

```typescript
// test/schema.test.ts
import { printSchema, buildSchema } from 'graphql';
import { schema } from '../src/schema';

describe('Schema', () => {
  it('matches snapshot', () => {
    expect(printSchema(schema)).toMatchSnapshot();
  });

  it('is valid', () => {
    // buildSchema throws if schema has errors
    expect(() => buildSchema(printSchema(schema))).not.toThrow();
  });
});
```

```bash
# First run creates __snapshots__/schema.test.ts.snap
# Subsequent runs compare against snapshot
# Update snapshots when you intentionally change the schema:
jest --updateSnapshot
```

---

## Testing with Testcontainers (Real Database)

```bash
npm install -D testcontainers @testcontainers/postgresql
```

```typescript
// test/setup/database.ts
import { PostgreSqlContainer } from '@testcontainers/postgresql';
import { Pool } from 'pg';

let container: PostgreSqlContainer;
let pool: Pool;

export async function startTestDatabase() {
  container = await new PostgreSqlContainer('postgres:16-alpine')
    .withDatabase('testdb')
    .withUsername('test')
    .withPassword('test')
    .start();

  pool = new Pool({ connectionString: container.getConnectionUri() });
  
  // Run migrations
  await runMigrations(pool);
  
  return pool;
}

export async function stopTestDatabase() {
  await pool.end();
  await container.stop();
}

export async function clearDatabase() {
  await pool.query('TRUNCATE users, posts, comments CASCADE');
}
```

```typescript
// test/posts.integration.test.ts
import { startTestDatabase, stopTestDatabase, clearDatabase } from './setup/database';

describe('Posts integration tests', () => {
  let db: Pool;

  beforeAll(async () => { db = await startTestDatabase(); });
  afterAll(async () => { await stopTestDatabase(); });
  beforeEach(async () => { await clearDatabase(); });

  it('creates and retrieves a post', async () => {
    const user = await db.query(
      'INSERT INTO users (name, email) VALUES ($1, $2) RETURNING *',
      ['Alice', 'alice@example.com']
    );

    const server = createTestServer();
    await server.start();

    const result = await executeOperation(
      server,
      `mutation { createPost(input: { title: "Hello", authorId: "${user.rows[0].id}" }) { post { id title } } }`,
      {},
      { user: user.rows[0], db: createDbClient(pool) }
    );

    expect(result?.data?.createPost.post.title).toBe('Hello');
    await server.stop();
  });
});
```

---

## Testing Subscriptions

```typescript
// test/subscriptions.test.ts
import { createPubSub } from 'graphql-subscriptions';
import { postAddedResolver } from '../src/resolvers/subscriptions';

describe('Subscription resolvers', () => {
  it('publishes post added event', async () => {
    const pubsub = createPubSub();
    const context = { pubsub };

    // Get the async iterator
    const iterator = await postAddedResolver.subscribe({}, {}, context, {});

    const newPost = { id: '1', title: 'Hello', authorId: '42' };

    // Publish event
    pubsub.publish('POST_ADDED', { postAdded: newPost });

    // Read from iterator
    const { value } = await iterator.next();
    expect(value.postAdded).toEqual(newPost);

    // Cleanup
    iterator.return?.();
  });
});
```

---

## Mocking External Services

```typescript
// test/mocks/stripe.mock.ts
export const mockStripe = {
  charges: {
    create: jest.fn().mockResolvedValue({ id: 'ch_test_123', status: 'succeeded' }),
  },
  customers: {
    create: jest.fn().mockResolvedValue({ id: 'cus_test_456' }),
  },
};

// In test:
jest.mock('stripe', () => () => mockStripe);

it('processes payment', async () => {
  const result = await executeOperation(server, PROCESS_PAYMENT_MUTATION, {
    input: { amount: 1000, token: 'tok_test' },
  });

  expect(mockStripe.charges.create).toHaveBeenCalledWith({
    amount: 1000,
    currency: 'usd',
    source: 'tok_test',
  });
  expect(result?.data?.processPayment.success).toBe(true);
});
```

---

## Testing DataLoader

```typescript
// test/dataloader.test.ts
import { createDataLoaders } from '../src/loaders';

describe('DataLoader', () => {
  it('batches user loads', async () => {
    const findUsers = jest.fn().mockResolvedValue([
      { id: '1', name: 'Alice' },
      { id: '2', name: 'Bob' },
    ]);

    const loaders = createDataLoaders({ db: { users: { findByIds: findUsers } } });

    // Load multiple users — should batch into one call
    const [alice, bob] = await Promise.all([
      loaders.userLoader.load('1'),
      loaders.userLoader.load('2'),
    ]);

    expect(findUsers).toHaveBeenCalledTimes(1);
    expect(findUsers).toHaveBeenCalledWith(['1', '2']);
    expect(alice.name).toBe('Alice');
    expect(bob.name).toBe('Bob');
  });

  it('caches within request', async () => {
    const findUsers = jest.fn().mockResolvedValue([{ id: '1', name: 'Alice' }]);
    const loaders = createDataLoaders({ db: { users: { findByIds: findUsers } } });

    await loaders.userLoader.load('1');
    await loaders.userLoader.load('1'); // same ID

    expect(findUsers).toHaveBeenCalledTimes(1); // cached
  });
});
```

---

## Schema Change Detection

Prevent accidental breaking changes in CI:

```bash
npm install -D @graphql-inspector/cli
```

```bash
# Compare local schema against production
graphql-inspector diff \
  'http://api.production.com/graphql' \
  './schema.graphql'

# In CI pipeline:
graphql-inspector diff \
  'git:origin/main:schema.graphql' \
  './schema.graphql'
```

```yaml
# .github/workflows/schema-check.yml
- name: Check for breaking schema changes
  run: npx graphql-inspector diff "${{ secrets.PROD_GRAPHQL_URL }}" schema.graphql
  env:
    GRAPHQL_INSPECTOR_FAIL_ON_BREAKING: true
```

---

## Test Coverage Checklist

```
For each query:
  ✓ Returns correct data shape
  ✓ Applies correct auth (throws when unauthenticated)
  ✓ Handles not-found (returns null or throws)
  ✓ Pagination works correctly

For each mutation:
  ✓ Creates/updates/deletes the resource
  ✓ Returns validation errors for invalid input
  ✓ Is idempotent where expected
  ✓ Auth is enforced

General:
  ✓ Schema snapshot test
  ✓ No N+1 queries (check DataLoader)
  ✓ Error format matches expected shape
  ✓ Subscription events fire correctly
```
# Chapter 18: Schema Design Best Practices

## Schema Design Principles

Good GraphQL schema design is about thinking from the client's perspective, not the database schema.

```
❌ Database-centric thinking:
   Mirror your tables → "Here's what we store"

✓ Client-centric thinking:
   Model what clients need → "Here's what you can do"
```

---

## Naming Conventions

```graphql
# Types: PascalCase
type UserProfile { ... }
type BlogPost { ... }

# Fields: camelCase
type User {
  id: ID!
  firstName: String!
  createdAt: DateTime!
  isActive: Boolean!
}

# Enum values: SCREAMING_SNAKE_CASE
enum UserRole {
  ADMIN
  MODERATOR
  REGULAR_USER
}

# Mutations: verb + noun (imperative action)
type Mutation {
  createUser(input: CreateUserInput!): CreateUserPayload!
  updateUserProfile(input: UpdateUserProfileInput!): UpdateUserProfilePayload!
  deleteUser(id: ID!): DeleteUserPayload!
  publishPost(id: ID!): PublishPostPayload!  # domain action, not CRUD
  sendPasswordResetEmail(email: String!): SendPasswordResetEmailPayload!
}

# Queries: noun (what you're fetching)
type Query {
  user(id: ID!): User
  users(filter: UserFilter, pagination: PaginationInput): UserConnection!
  me: User                    # viewer pattern
  post(slug: String!): Post   # natural identifier
}
```

---

## Nullability Design

Nullability is your API contract. Think carefully before making fields nullable.

```graphql
type User {
  id: ID!           # always has an ID — non-null
  email: String!    # required for users — non-null
  bio: String       # optional profile bio — nullable
  avatar: String    # might not have one — nullable
  
  # Relationships: nullable vs non-null list
  posts: [Post!]!   # has posts field always, posts in list always Posts
  # vs
  posts: [Post!]    # posts field might be null (e.g., error loading)
  # vs
  posts: [Post]!    # always returns array, but items might be null (rare, avoid)
}
```

### The null propagation decision:
```graphql
# If user resolver throws, null propagates up until it hits a nullable field
type Query {
  user(id: ID!): User    # nullable — if user not found, returns null safely
  me: User!              # non-null — if auth fails, error returns, not null
}

# Partial errors: make subtypes nullable so one failure doesn't kill everything
type Order {
  id: ID!
  total: Float!
  # If payment service is down, don't fail the whole order
  paymentStatus: String    # nullable — can be null if service unavailable
  shipment: Shipment       # nullable — independent service
}
```

---

## Input Design

```graphql
# ✓ Always use Input types for mutations (not inline args)
mutation {
  createUser(input: CreateUserInput!): CreateUserPayload!
}

input CreateUserInput {
  name: String!
  email: String!
  role: UserRole = REGULAR_USER    # default value
  bio: String
}

# ✓ Separate input types per operation (don't reuse)
input UpdateUserInput {
  name: String     # all optional — only update what's provided
  bio: String
  avatar: String
}

# ❌ Don't pass ID in update input — use args instead
mutation updateUser($id: ID!, $input: UpdateUserInput!) {
  updateUser(id: $id, input: $input) { ... }
}
```

---

## Payload Pattern (Mutations)

Always return a payload type from mutations:

```graphql
# ✓ Good: rich payload, room to evolve
type Mutation {
  createPost(input: CreatePostInput!): CreatePostPayload!
}

type CreatePostPayload {
  post: Post             # null if validation failed
  errors: [UserError!]!  # always present (empty array = success)
}

type UserError {
  field: String          # which field caused the error
  message: String!
}

# ❌ Don't return the object directly — no room for errors or metadata
type Mutation {
  createPost(input: CreatePostInput!): Post!   # avoid
}

# More evolved payload (can add fields without breaking clients)
type PublishPostPayload {
  post: Post
  errors: [UserError!]!
  # Added later without breaking existing queries:
  notifiedFollowers: Int
  scheduledAt: DateTime
}
```

---

## ID Strategy

```graphql
# Always use ID type for identifiers
type User {
  id: ID!   # opaque — clients don't parse it
}

# ✓ Global IDs (Relay spec) — encodes type + database ID
# "dXNlcjox" = base64("User:1")
# Allows any object to be fetched with a single node query:
type Query {
  node(id: ID!): Node
}
interface Node {
  id: ID!
}
type User implements Node {
  id: ID!   # globally unique
  name: String!
}

# ❌ Don't expose database primary keys directly as meaningful values
# Clients shouldn't parse "1234" out of an ID field
```

---

## Pagination Design

```graphql
# ✓ Cursor-based pagination (Relay Connection spec)
type Query {
  users(
    first: Int
    after: String
    last: Int
    before: String
    filter: UserFilter
    orderBy: UserOrderBy
  ): UserConnection!
}

type UserConnection {
  edges: [UserEdge!]!
  pageInfo: PageInfo!
  totalCount: Int!          # optional but useful
}

type UserEdge {
  node: User!
  cursor: String!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}

# Filter and ordering as input types
input UserFilter {
  role: UserRole
  isActive: Boolean
  search: String
  createdAfter: DateTime
}

enum UserOrderBy {
  CREATED_AT_ASC
  CREATED_AT_DESC
  NAME_ASC
  NAME_DESC
}
```

---

## Versioning and Schema Evolution

GraphQL avoids versioning by evolving the schema:

```graphql
# ✓ Add new fields (backward compatible)
type User {
  id: ID!
  name: String!
  # Added in v2 — old clients ignore this:
  displayName: String
  lastSeenAt: DateTime
}

# ✓ Deprecate old fields instead of removing
type User {
  id: ID!
  name: String! @deprecated(reason: "Use displayName instead")
  displayName: String!
}

# ✓ Add new arguments with defaults
type Query {
  # Old: posts(limit: Int)
  # New: added status filter — backward compatible (default: all)
  posts(limit: Int, status: PostStatus): [Post!]!
}

# ❌ Never do these (breaking changes):
# - Remove a field
# - Change a field's type (String → Int)
# - Make a nullable field non-null
# - Rename a field
# - Remove an enum value
# - Change argument type
```

### Breaking change detection:
```bash
# graphql-inspector catches breaking changes
graphql-inspector diff old-schema.graphql new-schema.graphql

# Exit code 1 if breaking changes found — use in CI
```

---

## Domain Modeling vs Database Modeling

```graphql
# Database tables:
# users(id, email, name, stripe_customer_id, plan_id, ...)
# plans(id, name, price, features_json, ...)

# ✓ GraphQL: model the domain, not the DB
type User {
  id: ID!
  email: String!
  name: String!
  subscription: Subscription   # composed from users + plans
}

type Subscription {
  plan: Plan!
  status: SubscriptionStatus!
  renewsAt: DateTime
  canceledAt: DateTime
}

type Plan {
  id: ID!
  name: String!
  price: Money!
  features: [String!]!
}

type Money {
  amount: Int!      # in cents
  currency: String!
  formatted: String!  # "$29.99" — computed field
}

# Users never see plan_id or stripe_customer_id
```

---

## Designing for Performance

### Avoid over-fetching at the schema level:
```graphql
# ✓ Don't return all data when a summary suffices
type Query {
  userSummary(id: ID!): UserSummary    # lightweight
  user(id: ID!): User                  # full profile
}

type UserSummary {
  id: ID!
  name: String!
  avatar: String
}

# Separate expensive fields into their own query
type User {
  id: ID!
  name: String!
  # Expensive — triggers additional DB query
  # Clients only request if they need it
  analytics: UserAnalytics
  activityLog(days: Int = 30): [ActivityEvent!]!
}
```

### Use DataLoader-friendly patterns:
```graphql
# ✓ Resolve via ID — batchable with DataLoader
type Post {
  id: ID!
  title: String!
  author: User!          # resolved by author_id → DataLoader
  comments: [Comment!]!  # resolved by post_id → DataLoader
}

# ❌ Don't design resolvers that make queries per-item
# author { byEmail(email: $email) }  — not batchable
```

---

## Schema Documentation

Self-documenting GraphQL schema:

```graphql
"""
A registered user of the platform.
"""
type User {
  "The unique identifier of this user."
  id: ID!
  
  "The user's display name."
  name: String!
  
  "The user's email address. Only visible to the user themselves and admins."
  email: String!
  
  "All posts authored by this user, ordered by creation date (newest first)."
  posts(
    "Maximum number of posts to return (default: 20, max: 100)"
    limit: Int = 20
    "Return posts after this cursor for pagination"
    after: String
  ): PostConnection!
  
  "The role of this user, which determines their permissions."
  role: UserRole!
  
  """
  Whether this user's account is currently active.
  Inactive users cannot log in but their content remains visible.
  """
  isActive: Boolean!
  
  "When the user registered."
  createdAt: DateTime!
  
  "When the user last logged in. Null if never logged in."
  lastLoginAt: DateTime
}
```

---

## Common Schema Anti-Patterns

```graphql
# ❌ Generic "data" or "result" field
type UserPayload {
  data: JSON    # no type safety, terrible DX
  error: String
}

# ❌ RPC-style naming
type Query {
  getUserById(id: ID!): User        # redundant "ById"
  getAllUsers: [User]                # don't prefix with "getAll"
  fetchUserData(userId: String): User  # verb + "Data"
}

# ✓ Clean naming
type Query {
  user(id: ID!): User
  users(filter: UserFilter): UserConnection!
}

# ❌ Deeply nested mutations
type Mutation {
  user: UserMutation    # "namespace mutation" — breaks tooling
}
type UserMutation {
  create(input: CreateUserInput!): User
  update(id: ID!, input: UpdateUserInput!): User
}

# ✓ Flat mutations
type Mutation {
  createUser(input: CreateUserInput!): CreateUserPayload!
  updateUser(id: ID!, input: UpdateUserInput!): UpdateUserPayload!
}

# ❌ Using String for everything
type Event {
  startTime: String   # is this ISO 8601? Unix timestamp? "tomorrow"?
  coordinates: String # "lat,lng"?
}

# ✓ Custom scalars for domain types
type Event {
  startTime: DateTime!
  coordinates: GeoPoint
}
```

---

## Schema Design Checklist

```
Naming:
  ✓ Types in PascalCase
  ✓ Fields in camelCase
  ✓ Enum values in SCREAMING_SNAKE_CASE
  ✓ Mutations: verb + noun

Nullability:
  ✓ IDs, required fields are non-null
  ✓ Optional profile fields are nullable
  ✓ Considered null propagation impact

Mutations:
  ✓ All mutations accept Input type
  ✓ All mutations return a Payload type
  ✓ Payload includes errors: [UserError!]!

Pagination:
  ✓ List fields use Connection pattern
  ✓ Filter and orderBy as typed inputs

Evolution:
  ✓ @deprecated instead of removal
  ✓ New fields are nullable (or have defaults)
  ✓ Breaking change detection in CI

Documentation:
  ✓ Every type has a description
  ✓ Non-obvious fields are documented
  ✓ Enum values documented
```
# Chapter 19: Apollo Federation

## What is Federation?

Apollo Federation lets you split a large GraphQL API into independent subgraphs — each owned by a different team — that compose into a single unified supergraph.

```
Without Federation:          With Federation:
Single monolithic schema     Multiple subgraph services
     │                            │
  One server                 Router (Supergraph)
  All resolvers              ┌────┴─────┐
                         Users     Products    Orders
                         Service   Service     Service
```

---

## Core Concepts

| Term | Definition |
|---|---|
| **Subgraph** | An independent GraphQL service that owns part of the schema |
| **Supergraph** | The unified schema composed from all subgraphs |
| **Router** | The entry point that receives client queries and fans out to subgraphs |
| **Entity** | A type shared across subgraphs, identified by `@key` |
| **Reference resolver** | A resolver that can reconstruct an entity from its key fields |

---

## Setting Up a Subgraph

```bash
npm install @apollo/subgraph graphql
```

### Users Subgraph:

```typescript
// users-service/schema.ts
import { buildSubgraphSchema } from '@apollo/subgraph';
import { gql } from 'graphql-tag';

const typeDefs = gql`
  extend schema
    @link(url: "https://specs.apollo.dev/federation/v2.3", import: ["@key"])

  type Query {
    me: User
    user(id: ID!): User
    users: [User!]!
  }

  # @key declares what makes a User uniquely identifiable
  # Other subgraphs can reference User by id
  type User @key(fields: "id") {
    id: ID!
    name: String!
    email: String!
    joinedAt: DateTime!
  }
`;

const resolvers = {
  Query: {
    me: (_, __, { user }) => user,
    user: (_, { id }, { db }) => db.users.findById(id),
    users: (_, __, { db }) => db.users.findAll(),
  },
  User: {
    // Reference resolver: called when another subgraph references a User by id
    __resolveReference: async ({ id }, { db }) => {
      return db.users.findById(id);
    },
  },
};

export const schema = buildSubgraphSchema({ typeDefs, resolvers });
```

---

### Products Subgraph:

```typescript
// products-service/schema.ts
const typeDefs = gql`
  extend schema
    @link(url: "https://specs.apollo.dev/federation/v2.3",
          import: ["@key", "@external", "@requires", "@provides"])

  type Query {
    product(id: ID!): Product
    products(category: String): [Product!]!
    featuredProducts: [Product!]!
  }

  type Product @key(fields: "id") {
    id: ID!
    name: String!
    price: Float!
    category: String!
    inStock: Boolean!
    description: String
  }
`;

const resolvers = {
  Query: {
    product: (_, { id }, { db }) => db.products.findById(id),
    products: (_, { category }, { db }) => db.products.findByCategory(category),
    featuredProducts: (_, __, { db }) => db.products.findFeatured(),
  },
  Product: {
    __resolveReference: ({ id }, { db }) => db.products.findById(id),
  },
};
```

---

### Orders Subgraph — References Entities from Other Subgraphs:

```typescript
// orders-service/schema.ts
const typeDefs = gql`
  extend schema
    @link(url: "https://specs.apollo.dev/federation/v2.3",
          import: ["@key", "@external"])

  type Query {
    order(id: ID!): Order
    myOrders: [Order!]!
  }

  type Order @key(fields: "id") {
    id: ID!
    status: OrderStatus!
    total: Float!
    placedAt: DateTime!
    # Reference to User entity — owned by Users subgraph
    customer: User!
    # Reference to multiple Product entities
    items: [OrderItem!]!
  }

  type OrderItem {
    product: Product!    # Product from Products subgraph
    quantity: Int!
    unitPrice: Float!
  }

  # Stub: User is owned by Users subgraph, but referenced here
  type User @key(fields: "id", resolvable: false) {
    id: ID!
  }

  # Stub: Product is owned by Products subgraph
  type Product @key(fields: "id", resolvable: false) {
    id: ID!
  }

  enum OrderStatus {
    PENDING
    PROCESSING
    SHIPPED
    DELIVERED
    CANCELLED
  }
`;

const resolvers = {
  Query: {
    order: (_, { id }, { db }) => db.orders.findById(id),
    myOrders: (_, __, { db, user }) => db.orders.findByCustomer(user.id),
  },
  Order: {
    __resolveReference: ({ id }, { db }) => db.orders.findById(id),
    customer: (order) => ({ __typename: 'User', id: order.customerId }),
    items: (order, _, { db }) => db.orderItems.findByOrder(order.id),
  },
  OrderItem: {
    product: (item) => ({ __typename: 'Product', id: item.productId }),
  },
};
```

---

## The Router (Apollo Router)

The router composes all subgraphs and serves clients a unified API.

```bash
# Install Apollo Router
curl -sSL https://router.apollo.dev/download/nix/latest | sh

# Or via npm (for development)
npm install @apollo/gateway
```

### Using Apollo Gateway (Node.js, simpler):

```typescript
// gateway/index.ts
import { ApolloServer } from '@apollo/server';
import { startStandaloneServer } from '@apollo/server/standalone';
import { ApolloGateway, IntrospectAndCompose } from '@apollo/gateway';

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs: [
      { name: 'users', url: 'http://users-service:4001/graphql' },
      { name: 'products', url: 'http://products-service:4002/graphql' },
      { name: 'orders', url: 'http://orders-service:4003/graphql' },
    ],
  }),
  // Pass headers from client to subgraphs
  buildService({ url }) {
    return new RemoteGraphQLDataSource({
      url,
      willSendRequest({ request, context }) {
        request.http?.headers.set('authorization', context.token || '');
        request.http?.headers.set('x-user-id', context.userId || '');
      },
    });
  },
});

const server = new ApolloServer({ gateway });

const { url } = await startStandaloneServer(server, {
  listen: { port: 4000 },
  context: async ({ req }) => ({
    token: req.headers.authorization,
    userId: getUserIdFromToken(req.headers.authorization),
  }),
});
```

### Apollo Router (production — binary):

```yaml
# router.yaml
supergraph:
  listen: 0.0.0.0:4000

subgraphs:
  users:
    routing_url: http://users-service:4001/graphql
  products:
    routing_url: http://products-service:4002/graphql
  orders:
    routing_url: http://orders-service:4003/graphql

# Pass JWT to all subgraphs
headers:
  all:
    request:
      - propagate:
          named: authorization
```

```bash
./router --config router.yaml --supergraph supergraph.graphql
```

---

## Composing the Supergraph Schema

```bash
# Install Rover CLI
npm install -g @apollo/rover

# Compose locally
rover supergraph compose --config supergraph.yaml

# supergraph.yaml:
federation_version: =2.3.0
subgraphs:
  users:
    routing_url: http://localhost:4001/graphql
    schema:
      file: ./users-service/schema.graphql
  products:
    routing_url: http://localhost:4002/graphql
    schema:
      file: ./products-service/schema.graphql
  orders:
    routing_url: http://localhost:4003/graphql
    schema:
      file: ./orders-service/schema.graphql
```

---

## Federation Directives

```graphql
# @key — declare entity, specify primary key fields
type Product @key(fields: "id") { ... }
type OrderItem @key(fields: "orderId productId") { ... }  # composite key

# @external — field owned by another subgraph
type User @key(fields: "id") {
  id: ID! @external
  email: String! @external   # owned by Users subgraph
  orders: [Order!]!          # owned by Orders subgraph
}

# @requires — need fields from another subgraph to resolve this field
type Product @key(fields: "id") {
  id: ID!
  price: Float!
  weight: Float! @external          # owned by another subgraph
  shippingCost: Float! @requires(fields: "weight")  # needs weight
}

# Resolver for shippingCost receives weight from Users subgraph:
Product: {
  shippingCost: (product) => calculateShipping(product.weight),
}

# @provides — tells router this subgraph can provide fields normally from another
type Order @key(fields: "id") {
  id: ID!
  product: Product @provides(fields: "name price")  # avoids extra roundtrip
}

type Product @key(fields: "id") {
  id: ID! @external
  name: String! @external
  price: Float! @external
}

# @shareable — multiple subgraphs can resolve this type/field
type Coordinates @shareable {
  lat: Float!
  lng: Float!
}

# @inaccessible — hide from client-facing schema (internal use)
type InternalAuditLog @inaccessible {
  ... 
}

# @override — migrate a field from one subgraph to another
type User @key(fields: "id") {
  id: ID!
  shippingAddress: Address @override(from: "orders")  # migrating this field here
}
```

---

## Query Planning

When a client requests data spanning multiple subgraphs, the router creates a query plan:

```graphql
# Client query spanning 3 subgraphs
query {
  myOrders {         # → Orders subgraph
    id
    status
    customer {       # → Users subgraph (entity lookup by id)
      name
      email
    }
    items {
      product {      # → Products subgraph (entity lookup by id)
        name
        price
        inStock
      }
      quantity
    }
  }
}
```

```
Query Plan:
1. Fetch myOrders from Orders subgraph
   → returns [{id, status, customerId, items: [{productId, quantity}]}]

2. Parallel fetch:
   a. Fetch users by ids=[customerId1, customerId2] from Users subgraph
   b. Fetch products by ids=[prod1, prod2, ...] from Products subgraph

3. Merge results and return to client
```

---

## Managed Federation with Apollo Studio

```bash
# Publish subgraph schema to Apollo Studio
rover subgraph publish my-graph@production \
  --name users \
  --schema ./users-service/schema.graphql \
  --routing-url https://users.internal/graphql

# Check for composition errors before deploying
rover subgraph check my-graph@production \
  --name users \
  --schema ./users-service/schema.graphql
```

Studio validates composition and detects breaking changes before you push.

---

## Docker Compose for Local Federation Development

```yaml
version: '3.8'
services:
  users-service:
    build: ./users-service
    ports:
      - "4001:4001"
    environment:
      DATABASE_URL: postgresql://postgres:postgres@postgres:5432/users
      PORT: 4001

  products-service:
    build: ./products-service
    ports:
      - "4002:4002"
    environment:
      DATABASE_URL: postgresql://postgres:postgres@postgres:5432/products
      PORT: 4002

  orders-service:
    build: ./orders-service
    ports:
      - "4003:4003"
    environment:
      DATABASE_URL: postgresql://postgres:postgres@postgres:5432/orders
      PORT: 4003

  router:
    image: ghcr.io/apollographql/router:v1.35.0
    ports:
      - "4000:4000"
    volumes:
      - ./router.yaml:/dist/config/router.yaml:ro
      - ./supergraph.graphql:/dist/config/supergraph.graphql:ro
    command: ["--config", "/dist/config/router.yaml", "--supergraph", "/dist/config/supergraph.graphql"]
    depends_on:
      - users-service
      - products-service
      - orders-service

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_PASSWORD: postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

---

## When to Use Federation

| Situation | Recommendation |
|---|---|
| Single team, one codebase | Don't use federation — monolith is simpler |
| 2–5 teams, growing API | Consider federation |
| Multiple teams, clear ownership | Federation is the right choice |
| Different tech stacks per service | Federation works across languages |
| Independent deploy cadences | Federation enables team autonomy |

### Federation is not for:
- Small projects with a single team
- When subgraph boundaries aren't clear
- When you want to avoid operational complexity

---

## Common Federation Pitfalls

```
❌ Circular @requires (A requires B, B requires A)
❌ Missing __resolveReference — entity can't be fetched
❌ Subgraph returning wrong __typename in references
❌ Large composite @key fields that are hard to look up
❌ Putting too many fields in @requires — triggers extra roundtrips
✓ Keep @key fields simple (usually just id)
✓ Return { __typename: 'User', id: user.id } from reference fields
✓ Use managed federation for production schema validation
```
# Chapter 20: Learning Roadmap

## Overview

This roadmap takes you from GraphQL beginner to production-ready API developer in 6 weeks with hands-on projects.

```
Week 1: Fundamentals        → Write queries, mutations, understand schemas
Week 2: Server Side         → Build a GraphQL server with real resolvers
Week 3: Data & Performance  → DataLoader, N+1, caching, pagination
Week 4: Client Side         → Apollo Client in React, subscriptions
Week 5: Production          → Auth, error handling, testing, monitoring
Week 6: Advanced Patterns   → Federation, schema design, real-world architecture
```

---

## Week 1: GraphQL Fundamentals

**Goal:** Be comfortable reading and writing GraphQL queries, mutations, and schemas.

### Topics:
- REST vs GraphQL — what problems GraphQL solves
- Schema Definition Language: types, scalars, nullability
- Queries: basic, nested, aliases, fragments, variables
- Mutations: input types, payload pattern
- Introspection: `__schema`, `__type`
- Using GraphiQL / Sandbox to explore APIs

### Exercises:

```graphql
# Practice with the GitHub GraphQL API (api.github.com/graphql)

# 1. Fetch your profile
query {
  viewer {
    login
    name
    bio
    followers { totalCount }
    following { totalCount }
  }
}

# 2. Fetch your repositories
query MyRepos($first: Int = 10) {
  viewer {
    repositories(first: $first, orderBy: { field: UPDATED_AT, direction: DESC }) {
      nodes {
        name
        description
        stargazerCount
        primaryLanguage { name }
        updatedAt
      }
      pageInfo { hasNextPage endCursor }
    }
  }
}

# 3. Create an issue (mutation)
mutation CreateIssue($repoId: ID!, $title: String!) {
  createIssue(input: { repositoryId: $repoId, title: $title }) {
    issue {
      id
      number
      url
    }
  }
}

# 4. Use a fragment to avoid repeating fields
fragment RepoDetails on Repository {
  id
  name
  stargazerCount
  description
}

query {
  viewer {
    starredRepositories(first: 5) {
      nodes { ...RepoDetails }
    }
  }
}
```

### Milestone check:
- [ ] Can read and write any GraphQL query or mutation
- [ ] Understand the SDL type system (nullability, interfaces, unions)
- [ ] Can navigate a schema via introspection

---

## Week 2: Building a GraphQL Server

**Goal:** Build a working GraphQL API from scratch with real database integration.

### Tech stack:
- Node.js + TypeScript
- Apollo Server or GraphQL Yoga
- PostgreSQL (from the PostgreSQL book)
- Prisma or pg library

### Project: Blog API

```typescript
// Build a blog API with these types:
// User, Post, Comment, Tag

// Required operations:
// Query: user(id), users, post(id), posts(filter, pagination), tags
// Mutation: createPost, updatePost, deletePost, addComment, createTag
// Types: all with correct nullability

// Step-by-step:
// 1. Write the schema in SDL
// 2. Implement stub resolvers that return mock data
// 3. Connect to PostgreSQL and implement real resolvers
// 4. Add context with db connection
// 5. Test in GraphiQL/Sandbox
```

```graphql
# Schema to implement:
type Query {
  me: User
  user(id: ID!): User
  post(id: ID!): Post
  posts(
    filter: PostFilter
    first: Int = 20
    after: String
  ): PostConnection!
}

type Mutation {
  createPost(input: CreatePostInput!): CreatePostPayload!
  updatePost(id: ID!, input: UpdatePostInput!): UpdatePostPayload!
  deletePost(id: ID!): DeletePostPayload!
  addComment(postId: ID!, body: String!): AddCommentPayload!
}

type Post @key(fields: "id") {
  id: ID!
  title: String!
  body: String!
  author: User!
  comments: [Comment!]!
  tags: [Tag!]!
  publishedAt: DateTime
  createdAt: DateTime!
}
```

### Milestone check:
- [ ] Server starts and accepts GraphQL requests
- [ ] At least 3 queries and 2 mutations work end-to-end
- [ ] Errors are handled gracefully

---

## Week 3: Data Layer and Performance

**Goal:** Implement DataLoader, cursor pagination, and understand N+1.

### Topics:
- N+1 problem — demonstrate and measure
- DataLoader: batching and caching
- Cursor-based pagination (Relay Connection spec)
- JSONB fields, custom scalars (DateTime, JSON)
- Query complexity limits
- Persisted queries

### Exercises:

```typescript
// 1. Demonstrate N+1: log all DB queries for a posts query
//    that fetches authors — count the DB calls

// 2. Fix with DataLoader
const userLoader = new DataLoader<string, User>(async (ids) => {
  const users = await db.query(
    'SELECT * FROM users WHERE id = ANY($1)', [ids]
  );
  // MUST return results in same order as ids
  const userMap = new Map(users.rows.map(u => [u.id, u]));
  return ids.map(id => userMap.get(id) ?? new Error(`User ${id} not found`));
});

// 3. Implement cursor pagination for posts
// 4. Add query depth limit (prevent deeply nested attacks)
// 5. Add query complexity scoring
```

### Milestone check:
- [ ] DataLoader eliminates N+1 queries (verified by query count logs)
- [ ] Pagination works correctly (hasNextPage, endCursor, fetchMore)
- [ ] Query depth/complexity limits prevent abuse

---

## Week 4: GraphQL Client (React + Apollo)

**Goal:** Build a React frontend that consumes your GraphQL API.

### Project: Blog Frontend

```typescript
// Build a React app with:
// 1. Post list with infinite scroll (fetchMore)
// 2. Post detail with comments
// 3. Create post form (useMutation + cache update)
// 4. Authentication (login mutation → store JWT)
// 5. Real-time: new comment notification via subscription

// Required hooks to master:
// - useQuery (fetchPolicy: 'cache-and-network')
// - useMutation (with cache.modify or refetchQueries)
// - useSubscription
// - useLazyQuery (trigger query on user action)

// Required cache operations:
// - Read from cache: client.readQuery / readFragment
// - Write to cache: cache.modify, cache.writeFragment
// - Evict from cache: cache.evict + cache.gc
```

```typescript
// Key pattern: optimistic updates
const [likePost] = useMutation(LIKE_POST, {
  optimisticResponse: {
    likePost: {
      __typename: 'LikePostPayload',
      post: {
        __typename: 'Post',
        id: postId,
        likeCount: currentLikeCount + 1,
        isLikedByMe: true,
      },
      errors: [],
    },
  },
});
```

### Milestone check:
- [ ] List renders with loading/error/empty states
- [ ] Mutations update the cache without refetching
- [ ] Subscriptions update the UI in real time

---

## Week 5: Production Readiness

**Goal:** Harden the API with authentication, error handling, testing, and monitoring.

### Topics:
- JWT authentication via context
- graphql-shield for authorization rules
- Error formatting and sanitization
- Comprehensive testing: unit, integration, E2E
- Logging and monitoring (Apollo Studio, custom plugins)
- Rate limiting and introspection disabling in production

### Hardening checklist:

```typescript
// 1. Disable introspection in production
const server = new ApolloServer({
  introspection: process.env.NODE_ENV !== 'production',
});

// 2. Sanitize errors in production
formatError: (formattedError) => {
  if (process.env.NODE_ENV === 'production') {
    if (formattedError.extensions?.code === 'INTERNAL_SERVER_ERROR') {
      return { message: 'Internal server error', extensions: { code: 'INTERNAL_SERVER_ERROR' } };
    }
  }
  return formattedError;
},

// 3. Add query depth and complexity limits
import depthLimit from 'graphql-depth-limit';
import { createComplexityLimitRule } from 'graphql-validation-complexity';

const server = new ApolloServer({
  validationRules: [
    depthLimit(10),
    createComplexityLimitRule(5000),
  ],
});

// 4. Write tests for every resolver
// Unit: mock DB, test business logic
// Integration: real test DB, test full operation

// 5. Monitor with Apollo Studio or custom plugin
const loggingPlugin: ApolloServerPlugin = {
  requestDidStart: async () => ({
    willSendResponse: async ({ response, contextValue }) => {
      logger.info({ 
        userId: contextValue.user?.id,
        errors: response.body.singleResult?.errors?.length,
      });
    },
  }),
};
```

### Milestone check:
- [ ] Unauthenticated requests to protected fields are rejected
- [ ] Production errors don't leak implementation details
- [ ] Test suite covers happy path + error cases for all operations
- [ ] Slow operations are logged and traced

---

## Week 6: Advanced Patterns and Architecture

**Goal:** Understand federation, schema design patterns, and real-world GraphQL architecture.

### Topics:
- Apollo Federation: subgraphs, @key, router
- Schema design review: nullability, naming, versioning
- Schema registry and change management
- GraphQL at scale: caching (CDN, APQ), persisted queries
- Alternatives: GraphQL Yoga, Pothos, URQL, TanStack Query + GraphQL

### Project: Federated E-commerce API

```
Split your blog API into federated subgraphs:

Users Subgraph (port 4001):
  - User type with @key(fields: "id")
  - Auth: login, signup mutations
  - Profile management

Content Subgraph (port 4002):
  - Post type with @key(fields: "id")
  - References User entity
  - Full-text search

Comments Subgraph (port 4003):
  - Comment type with @key(fields: "id")
  - References User and Post entities
  - Real-time via subscriptions

Router (port 4000):
  - Composes all subgraphs
  - Handles auth header propagation
```

### Milestone check:
- [ ] Three subgraphs running independently
- [ ] Cross-subgraph query works (e.g., post with author from users subgraph)
- [ ] Can explain query plan to a colleague

---

## Capstone Project: Production GraphQL API

Build a complete, production-ready GraphQL API:

```
Requirements:
1. Federated architecture (2+ subgraphs)
2. JWT authentication with graphql-shield authorization
3. DataLoader for all related data
4. Cursor-based pagination for all lists
5. File uploads (profile pictures)
6. Real-time subscriptions (at least one feature)
7. Full test suite (unit + integration)
8. Apollo Studio connected for monitoring
9. Docker Compose for local development
10. Breaking change detection in CI/CD pipeline

Suggested domain: Project Management Tool
  - Users with roles (admin, member, viewer)
  - Projects with members
  - Tasks with assignments, due dates, attachments
  - Comments with real-time notifications
  - Activity log
```

---

## Essential Resources

| Resource | Type | Focus |
|---|---|---|
| graphql.org/learn | Official docs | Language fundamentals |
| apollographql.com/docs | Official docs | Apollo ecosystem |
| the-guild.dev/blog | Blog | Framework-agnostic best practices |
| graphqlconf.org | Conference talks | Real-world case studies |
| "Production Ready GraphQL" — Marc-André Giroux | Book | API design patterns |
| graphql-spec.github.io | Spec | Language specification |
| Principled GraphQL (apollographql.com) | Guide | Design principles at scale |

---

## Skill Level Benchmarks

| Level | Can Do |
|---|---|
| **Beginner** | Write queries/mutations, use GraphiQL, understand the type system |
| **Intermediate** | Build a server, implement resolvers, use Apollo Client in React |
| **Advanced** | DataLoader, cursor pagination, testing, authentication, error handling |
| **Expert** | Federation, schema design at scale, performance optimization, custom directives |

---

## Interview Topics

```
Concepts interviewers test:
1. "Explain N+1 and how DataLoader solves it"
2. "When would you use GraphQL vs REST?"
3. "How do you handle authentication in GraphQL?"
4. "Explain cursor vs offset pagination — tradeoffs"
5. "What is Apollo Federation and when would you use it?"
6. "How does the GraphQL cache work in Apollo Client?"
7. "How do you prevent malicious queries (depth, complexity)?"
8. "What are the tradeoffs of nullable vs non-null fields?"
```

---

## Next Steps Beyond This Course

- **URQL** — lighter alternative to Apollo Client
- **TanStack Query + graphql-request** — for React Query fans
- **Relay** — Facebook's opinionated GraphQL client (strict Connection spec)
- **GraphQL Code Generator** — generate full TypeScript types from schema
- **gqlgen** (Go) or **Strawberry** (Python) — server implementations in other languages
- **GraphQL over HTTP** — the standard spec for transport
- **Persisted Queries / APQ** — reduce request size, enable CDN caching
- **Defer / Stream** — incremental delivery for large responses
