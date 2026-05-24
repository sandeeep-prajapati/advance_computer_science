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
