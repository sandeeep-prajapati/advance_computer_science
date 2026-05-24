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
