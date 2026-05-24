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
