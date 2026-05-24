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
