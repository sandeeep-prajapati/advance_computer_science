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
