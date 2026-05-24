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
