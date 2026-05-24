# Chapter 11: N+1 Problem & DataLoader

## The N+1 Problem

The N+1 problem is the most common performance issue in GraphQL APIs. It occurs when loading a list of N items, then making 1 additional query per item.

```graphql
query {
  posts {        # 1 query → returns 10 posts
    id
    title
    author {     # 10 queries → 1 per post author!
      name
      avatar
    }
  }
}
```

Total: **11 queries** for 10 posts. With 100 posts: **101 queries**. This is the N+1 problem.

```javascript
// The problematic resolver
const resolvers = {
  Post: {
    author: async (post, args, context) => {
      // This runs once PER POST — terrible!
      return context.db.users.findById(post.authorId);
    },
  },
};
```

---

## DataLoader — The Solution

**DataLoader** batches and caches database calls within a single request.

```
Without DataLoader:     With DataLoader:
  findUser(1)             collect: [1, 2, 3, 4, 5]
  findUser(2)           ─→  findUsers([1,2,3,4,5])  (1 query!)
  findUser(3)
  findUser(4)
  findUser(5)
```

### Install:
```bash
npm install dataloader
```

### Basic DataLoader:
```javascript
import DataLoader from 'dataloader';

// Batch function: receives array of keys, must return array of values in same order
async function batchLoadUsers(userIds) {
  const users = await db.query(
    'SELECT * FROM users WHERE id = ANY($1)',
    [userIds]
  );

  // IMPORTANT: must return values in the same order as keys
  const userMap = new Map(users.rows.map(user => [user.id, user]));
  return userIds.map(id => userMap.get(id) || new Error(`User ${id} not found`));
}

const userLoader = new DataLoader(batchLoadUsers);
```

---

## Per-Request DataLoader Instances

DataLoaders should be **created per request** — not shared globally. This ensures the cache is per-request and doesn't leak between users.

```javascript
// Create fresh loaders for each request
function createDataLoaders(db) {
  return {
    users: new DataLoader(async (ids) => {
      const users = await db.users.findManyByIds(ids);
      const userMap = new Map(users.map(u => [u.id, u]));
      return ids.map(id => userMap.get(id) ?? null);
    }),

    postsbyAuthor: new DataLoader(async (authorIds) => {
      const posts = await db.posts.findByAuthorIds(authorIds);
      // Group posts by authorId
      const grouped = new Map();
      authorIds.forEach(id => grouped.set(id, []));
      posts.forEach(post => grouped.get(post.authorId)?.push(post));
      return authorIds.map(id => grouped.get(id) ?? []);
    }),

    commentsByPost: new DataLoader(async (postIds) => {
      const comments = await db.comments.findByPostIds(postIds);
      const grouped = new Map();
      postIds.forEach(id => grouped.set(id, []));
      comments.forEach(c => grouped.get(c.postId)?.push(c));
      return postIds.map(id => grouped.get(id) ?? []);
    }),
  };
}

// In Apollo Server context
const server = new ApolloServer({
  typeDefs,
  resolvers,
  context: ({ req }) => ({
    currentUser: getUser(req),
    db,
    loaders: createDataLoaders(db), // fresh instance per request
  }),
});
```

---

## Using DataLoader in Resolvers

```javascript
const resolvers = {
  Post: {
    // Uses DataLoader — batched automatically!
    author: (post, args, context) => {
      return context.loaders.users.load(post.authorId);
    },

    comments: (post, args, context) => {
      return context.loaders.commentsByPost.load(post.id);
    },
  },

  Comment: {
    author: (comment, args, context) => {
      return context.loaders.users.load(comment.authorId);
    },
  },

  User: {
    posts: (user, args, context) => {
      return context.loaders.postsByAuthor.load(user.id);
    },
  },
};
```

Now fetching 10 posts with authors:
- **Before DataLoader:** 1 + 10 = 11 queries
- **After DataLoader:** 1 + 1 = **2 queries**

---

## DataLoader with Arguments

When field arguments affect the query, include them in the key:

```javascript
// Problem: posts(limit: 5) vs posts(limit: 10) need different loaders
// Solution: include args in the key, or use separate loaders

const postsByAuthorLoader = new DataLoader(
  async (keys) => {
    // keys = [{ authorId: 1, limit: 5 }, { authorId: 2, limit: 10 }, ...]
    // Group by limit, batch per limit value
    const grouped = {};
    keys.forEach(key => {
      const k = key.limit;
      if (!grouped[k]) grouped[k] = [];
      grouped[k].push(key.authorId);
    });

    const results = {};
    await Promise.all(
      Object.entries(grouped).map(async ([limit, authorIds]) => {
        const posts = await db.posts.findByAuthorIds(authorIds, { limit: parseInt(limit) });
        // map back...
      })
    );

    return keys.map(key => results[key.authorId] ?? []);
  },
  {
    // Custom cache key function
    cacheKeyFn: (key) => `${key.authorId}:${key.limit}`,
  }
);

// Usage in resolver
const resolvers = {
  User: {
    posts: (user, { limit = 10 }, context) => {
      return context.loaders.postsByAuthor.load({ authorId: user.id, limit });
    },
  },
};
```

---

## DataLoader Options

```javascript
const loader = new DataLoader(batchFn, {
  // Maximum batch size (default: unlimited)
  maxBatchSize: 100,

  // Disable batching (useful for testing)
  batch: false,

  // Disable caching (for mutations or time-sensitive data)
  cache: false,

  // Custom cache key function (default: identity)
  cacheKeyFn: (key) => JSON.stringify(key),

  // Custom cache map (e.g., use a Map with TTL)
  cacheMap: new Map(),
});
```

---

## Priming the Cache

Pre-populate the DataLoader cache to avoid unnecessary loads:

```javascript
const resolvers = {
  Query: {
    posts: async (parent, { limit }, context) => {
      const posts = await context.db.posts.findAll({ limit });

      // Prime the author loader with the authors we already have joined data for
      posts.forEach(post => {
        if (post.author) {
          context.loaders.users.prime(post.author.id, post.author);
        }
      });

      return posts;
    },
  },
};
```

---

## Common DataLoader Patterns

### Many-to-many relationships:
```javascript
// Load tags for posts
const tagsByPostLoader = new DataLoader(async (postIds) => {
  const rows = await db.query(`
    SELECT pt.post_id, t.id, t.name
    FROM post_tags pt
    JOIN tags t ON t.id = pt.tag_id
    WHERE pt.post_id = ANY($1)
  `, [postIds]);

  const grouped = new Map(postIds.map(id => [id, []]));
  rows.forEach(row => grouped.get(row.post_id)?.push({ id: row.id, name: row.name }));
  return postIds.map(id => grouped.get(id) ?? []);
});
```

### Count loader:
```javascript
const commentCountLoader = new DataLoader(async (postIds) => {
  const rows = await db.query(`
    SELECT post_id, COUNT(*) AS count
    FROM comments
    WHERE post_id = ANY($1)
    GROUP BY post_id
  `, [postIds]);

  const counts = new Map(rows.map(r => [r.post_id, parseInt(r.count)]));
  return postIds.map(id => counts.get(id) ?? 0);
});

// Resolver
const resolvers = {
  Post: {
    commentCount: (post, args, context) =>
      context.loaders.commentCount.load(post.id),
  },
};
```

---

## Measuring the Improvement

```javascript
// Log all DB queries during a request
const db = {
  query: async (sql, params) => {
    console.log(`DB QUERY: ${sql}`);
    return pool.query(sql, params);
  },
};

// Without DataLoader:
// DB QUERY: SELECT * FROM posts LIMIT 10
// DB QUERY: SELECT * FROM users WHERE id = 1
// DB QUERY: SELECT * FROM users WHERE id = 2
// DB QUERY: SELECT * FROM users WHERE id = 1  (duplicate!)
// DB QUERY: SELECT * FROM users WHERE id = 3
// ... 10 more queries

// With DataLoader:
// DB QUERY: SELECT * FROM posts LIMIT 10
// DB QUERY: SELECT * FROM users WHERE id = ANY('{1,2,3,4,5,6,7,8,9,10}')
// Total: 2 queries (with deduplication — note id=1 queried once)
```
