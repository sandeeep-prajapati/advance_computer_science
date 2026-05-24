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
