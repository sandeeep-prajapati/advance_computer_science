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
