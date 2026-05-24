# Chapter 4: Mutations

## What are Mutations?

Mutations are GraphQL operations that **modify** data on the server — create, update, delete. Unlike queries (which are typically safe to run in parallel), mutations run **sequentially** in the order listed.

```graphql
# Mutation keyword
mutation {
  createUser(input: { name: "Alice", email: "alice@example.com", password: "pass123" }) {
    id
    name
    email
  }
}
```

---

## Basic Mutation Structure

```graphql
# Schema definition
type Mutation {
  createPost(input: CreatePostInput!): Post!
  updatePost(id: ID!, input: UpdatePostInput!): Post!
  deletePost(id: ID!): Boolean!
  publishPost(id: ID!): Post!
}

input CreatePostInput {
  title: String!
  body: String!
  tags: [String!]! = []
}

input UpdatePostInput {
  title: String
  body: String
  tags: [String!]
}
```

Resolvers:
```javascript
const resolvers = {
  Mutation: {
    createPost: async (parent, { input }, context) => {
      if (!context.user) throw new AuthenticationError('Must be logged in');

      const post = await context.db.posts.create({
        ...input,
        authorId: context.user.id,
        status: 'DRAFT',
      });

      return post;
    },

    updatePost: async (parent, { id, input }, context) => {
      const post = await context.db.posts.findById(id);
      if (!post) throw new UserInputError('Post not found');
      if (post.authorId !== context.user.id) throw new ForbiddenError('Not your post');

      return context.db.posts.update(id, input);
    },

    deletePost: async (parent, { id }, context) => {
      await context.db.posts.delete(id);
      return true;
    },
  },
};
```

---

## Input Types for Mutations

Always use dedicated input types for mutations — don't reuse output types:

```graphql
# Good: dedicated input types
input CreateUserInput {
  name: String!
  email: String!
  password: String!
  role: UserRole = USER
}

input UpdateUserInput {
  name: String        # all fields optional for partial updates
  email: String
  bio: String
  avatarUrl: String
}

# Bad: passing fields directly (becomes unwieldy)
type Mutation {
  createUser(name: String!, email: String!, password: String!, role: UserRole): User!
}
```

---

## Mutation Return Types — Payload Pattern

Rather than returning the entity directly, return a payload type that can include both the result AND errors:

```graphql
type CreatePostPayload {
  post: Post              # null if mutation failed
  errors: [UserError!]!  # empty if success
}

type UserError {
  field: String          # which input field caused the error (null for global errors)
  message: String!
}

type Mutation {
  createPost(input: CreatePostInput!): CreatePostPayload!
  updatePost(id: ID!, input: UpdatePostInput!): UpdatePostPayload!
  deletePost(id: ID!): DeletePostPayload!
}
```

Client:
```graphql
mutation CreatePost($input: CreatePostInput!) {
  createPost(input: $input) {
    post {
      id
      title
    }
    errors {
      field
      message
    }
  }
}
```

This pattern avoids throwing HTTP errors for business-logic failures — the client always gets a 200 with errors inside the payload.

---

## Variables with Mutations

```graphql
mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    id
    name
    email
    createdAt
  }
}
```

Variables:
```json
{
  "input": {
    "name": "Alice",
    "email": "alice@example.com",
    "password": "SecurePass123!",
    "role": "USER"
  }
}
```

---

## Multiple Mutations in One Request

Mutations in a single request are executed **sequentially**, not in parallel:

```graphql
mutation BatchOperations {
  first: createPost(input: { title: "Post 1", body: "..." }) {
    id title
  }
  second: createPost(input: { title: "Post 2", body: "..." }) {
    id title
  }
  updateProfile: updateUser(id: "1", input: { bio: "New bio" }) {
    id bio
  }
}
```

`first` completes before `second` starts.

---

## Common Mutation Patterns

### Create:
```graphql
mutation CreateComment($input: CreateCommentInput!) {
  createComment(input: $input) {
    comment {
      id
      body
      author { name avatar }
      createdAt
    }
    errors { field message }
  }
}
```

### Update (Partial):
```graphql
mutation UpdateProfile($input: UpdateUserInput!) {
  updateUser(id: "me", input: $input) {
    user {
      id name bio email
    }
    errors { field message }
  }
}
```

### Delete:
```graphql
mutation DeletePost($id: ID!) {
  deletePost(id: $id) {
    success
    deletedId    # return ID so client can remove from cache
    errors { message }
  }
}
```

### Toggle (state change):
```graphql
mutation ToggleLike($postId: ID!) {
  toggleLike(postId: $postId) {
    post {
      id
      liked     # new state
      likeCount # updated count
    }
  }
}
```

---

## Optimistic Updates (Client Side)

```javascript
// Apollo Client — update UI before mutation completes
const [createPost] = useMutation(CREATE_POST, {
  optimisticResponse: {
    createPost: {
      post: {
        __typename: 'Post',
        id: 'temp-id',
        title: input.title,
        status: 'DRAFT',
        createdAt: new Date().toISOString(),
      },
      errors: [],
    },
  },
  update(cache, { data: { createPost } }) {
    // Update the Apollo cache with the new post
    cache.modify({
      fields: {
        posts(existingPosts = []) {
          const newPostRef = cache.writeFragment({
            data: createPost.post,
            fragment: POST_FRAGMENT,
          });
          return [newPostRef, ...existingPosts];
        },
      },
    });
  },
});
```

---

## Mutation Resolver Implementation Patterns

### With validation:
```javascript
createUser: async (parent, { input }, context) => {
  // Validate
  const errors = validateCreateUserInput(input);
  if (errors.length > 0) {
    return { user: null, errors };
  }

  // Check uniqueness
  const existing = await context.db.users.findByEmail(input.email);
  if (existing) {
    return {
      user: null,
      errors: [{ field: 'email', message: 'Email is already in use' }],
    };
  }

  // Hash password
  const passwordHash = await bcrypt.hash(input.password, 10);

  // Create
  const user = await context.db.users.create({
    ...input,
    password: passwordHash,
  });

  // Side effects
  await context.emailService.sendWelcome(user.email);
  context.pubsub.publish('USER_CREATED', { userCreated: user });

  return { user, errors: [] };
},
```

### With transactions:
```javascript
placeOrder: async (parent, { input }, context) => {
  return context.db.transaction(async (trx) => {
    const order = await context.db.orders.create(input, trx);
    await context.db.inventory.deduct(input.items, trx);
    await context.db.payments.charge(input.paymentMethod, order.total, trx);
    return { order, errors: [] };
  });
},
```

---

## Idempotency

For mutations that should be safe to retry:

```graphql
input CreateOrderInput {
  idempotencyKey: String!   # unique per attempt, same key = same result
  items: [OrderItemInput!]!
  paymentMethod: ID!
}
```

Server checks if the key was already used and returns the existing result.

---

## Mutation Best Practices

- Return the modified entity so clients can update their cache
- Use payload types for mutations — include both data and errors
- Validate input server-side — never trust client
- Use input types, not individual arguments for complex mutations
- Name mutations with verb + noun: `createUser`, `publishPost`, `cancelOrder`
- Publish events after successful mutations (for subscriptions)
- Use transactions for multi-step operations
- Return `deletedId` from delete mutations so clients can remove from cache
