# Chapter 6: Resolvers

## What is a Resolver?

A **resolver** is a function that returns the value for a field in your GraphQL schema. Every field in GraphQL has a resolver — if you don't provide one, a default resolver returns the field from the parent object.

```
Schema field → Resolver function → Data (from DB, API, cache, etc.)
```

---

## Resolver Signature

Every resolver receives four arguments:

```javascript
fieldName(parent, args, context, info) {
  // parent  — result from parent resolver (null for root types)
  // args    — arguments passed to the field in the query
  // context — shared object across all resolvers (db, auth, dataloaders)
  // info    — AST info about the field being resolved (rarely used)
}
```

---

## Resolver Map

Resolvers are organized in a map matching the schema structure:

```javascript
const resolvers = {
  // Root Query type
  Query: {
    user: async (parent, { id }, context) => {
      return context.db.users.findById(id);
    },

    users: async (parent, { limit = 20, offset = 0 }, context) => {
      return context.db.users.findAll({ limit, offset });
    },

    me: async (parent, args, context) => {
      if (!context.user) return null;
      return context.db.users.findById(context.user.id);
    },
  },

  // Root Mutation type
  Mutation: {
    createUser: async (parent, { input }, context) => {
      return context.db.users.create(input);
    },
  },

  // Object type resolvers (field-level)
  User: {
    // posts is NOT in the DB user row — this resolver loads it
    posts: async (user, { limit = 10 }, context) => {
      return context.db.posts.findByAuthorId(user.id, { limit });
    },

    // fullName is a computed field
    fullName: (user) => `${user.firstName} ${user.lastName}`,

    // avatar with a fallback
    avatar: (user) => user.avatarUrl || `https://avatars.dicebear.com/api/identicon/${user.id}.svg`,
  },

  Post: {
    // author is stored as authorId in DB — resolver loads the User object
    author: async (post, args, context) => {
      return context.db.users.findById(post.authorId);
    },
  },
};
```

---

## Default Resolvers

If you don't define a resolver for a field, GraphQL uses the **default resolver**:

```javascript
// This is what the default resolver does:
(parent, args, context, info) => {
  return parent[info.fieldName];
}
```

So `user.name` works automatically if the DB object has a `name` field — no resolver needed.

---

## Context — Sharing State Across Resolvers

Context is created once per request and passed to every resolver:

```javascript
// Apollo Server
const server = new ApolloServer({
  typeDefs,
  resolvers,
  context: async ({ req }) => {
    // Extract token from header
    const token = req.headers.authorization?.replace('Bearer ', '');

    // Verify and decode
    let user = null;
    if (token) {
      try {
        user = jwt.verify(token, process.env.JWT_SECRET);
      } catch {
        // invalid token — user stays null
      }
    }

    return {
      user,                                    // current authenticated user
      db,                                      // database connection
      dataloaders: createDataLoaders(db),      // per-request DataLoader instances
      redis,                                   // cache
      logger,                                  // shared logger
    };
  },
});
```

---

## Resolver Chain (Execution Flow)

```graphql
query {
  user(id: "1") {      # Query.user resolver called
    name               # User.name (default resolver)
    posts {            # User.posts resolver called with user as parent
      title            # Post.title (default resolver)
      author {         # Post.author resolver called with post as parent
        name           # User.name (default resolver)
      }
    }
  }
}
```

Execution:
1. `Query.user(null, { id: "1" }, context)` → returns `{ id: 1, name: "Alice", ... }`
2. `User.posts({ id: 1, name: "Alice", ... }, {}, context)` → returns `[{ id: 10, title: "..." }]`
3. `Post.author({ id: 10, title: "...", authorId: 1 }, {}, context)` → returns `{ id: 1, name: "Alice" }`

---

## Async Resolvers

Most resolvers are async — they fetch data from databases or APIs:

```javascript
const resolvers = {
  Query: {
    // Can return a Promise
    user: (parent, { id }, context) => {
      return context.db.query('SELECT * FROM users WHERE id = $1', [id])
        .then(result => result.rows[0]);
    },

    // Or use async/await (equivalent, more readable)
    users: async (parent, { filter }, context) => {
      const users = await context.db.users.findAll(filter);
      return users;
    },
  },
};
```

---

## Error Handling in Resolvers

```javascript
const resolvers = {
  Query: {
    user: async (parent, { id }, context) => {
      const user = await context.db.users.findById(id);

      if (!user) {
        throw new GraphQLError('User not found', {
          extensions: {
            code: 'USER_NOT_FOUND',
            http: { status: 404 },
          },
        });
      }

      return user;
    },
  },

  Mutation: {
    createUser: async (parent, { input }, context) => {
      if (!context.user) {
        throw new GraphQLError('Authentication required', {
          extensions: { code: 'UNAUTHENTICATED' },
        });
      }

      if (context.user.role !== 'ADMIN') {
        throw new GraphQLError('Insufficient permissions', {
          extensions: { code: 'FORBIDDEN' },
        });
      }

      try {
        return await context.db.users.create(input);
      } catch (err) {
        if (err.code === '23505') {  // unique violation
          throw new GraphQLError('Email already in use', {
            extensions: { code: 'DUPLICATE_EMAIL', field: 'email' },
          });
        }
        throw err;
      }
    },
  },
};
```

---

## Resolver Best Practices: Thin Resolvers

Keep resolvers thin — delegate to service classes:

```javascript
// Bad: business logic in resolver
const resolvers = {
  Mutation: {
    createUser: async (parent, { input }, context) => {
      const hash = await bcrypt.hash(input.password, 12);
      const existing = await context.db.users.findByEmail(input.email);
      if (existing) throw new Error('Email taken');
      const user = await context.db.users.create({ ...input, password: hash });
      await sendWelcomeEmail(user.email);
      return user;
    },
  },
};

// Good: resolver delegates to service
const resolvers = {
  Mutation: {
    createUser: (parent, { input }, context) =>
      context.services.userService.create(input),
  },
};

// UserService handles the logic
class UserService {
  async create(input) {
    const hash = await bcrypt.hash(input.password, 12);
    const existing = await this.db.users.findByEmail(input.email);
    if (existing) throw new DuplicateEmailError(input.email);
    const user = await this.db.users.create({ ...input, password: hash });
    await this.emailService.sendWelcome(user.email);
    return user;
  }
}
```

---

## The `info` Argument

The `info` argument contains the full query AST — useful for advanced use cases:

```javascript
const resolvers = {
  Query: {
    user: async (parent, { id }, context, info) => {
      // Check which fields were requested — optimize DB query
      const requestedFields = info.fieldNodes[0].selectionSet.selections
        .map(s => s.name.value);

      if (!requestedFields.includes('posts')) {
        // Don't join posts if they weren't requested
        return context.db.users.findById(id, { includePosts: false });
      }

      return context.db.users.findById(id, { includePosts: true });
    },
  },
};
```

Libraries like `graphql-parse-resolve-info` and `join-monster` automate this pattern.

---

## Abstract Type Resolvers

For interfaces and unions, you need a `__resolveType` resolver:

```javascript
const resolvers = {
  // For union types
  SearchResult: {
    __resolveType(obj) {
      if (obj.title) return 'Post';
      if (obj.email) return 'User';
      if (obj.price) return 'Product';
      return null;
    },
  },

  // For interface types
  Node: {
    __resolveType(obj) {
      return obj.__typename || obj.constructor.name;
    },
  },
};
```
