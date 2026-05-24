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
