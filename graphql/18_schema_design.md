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
