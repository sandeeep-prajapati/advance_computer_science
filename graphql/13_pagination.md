# Chapter 13: Pagination

## Why Pagination?

Never return all records at once — always paginate. Without pagination:
- Responses can be megabytes of JSON
- Database queries become full table scans
- Memory usage grows unbounded
- UI becomes unusable

---

## Offset Pagination

The simplest approach — skip N records, take M.

### Schema:
```graphql
type Query {
  posts(limit: Int = 20, offset: Int = 0): PostList!
}

type PostList {
  items: [Post!]!
  total: Int!
  hasMore: Boolean!
}
```

### Resolver:
```javascript
const resolvers = {
  Query: {
    posts: async (parent, { limit = 20, offset = 0 }, context) => {
      const [items, total] = await Promise.all([
        context.db.posts.findAll({ limit, offset, orderBy: 'createdAt', order: 'DESC' }),
        context.db.posts.count(),
      ]);

      return {
        items,
        total,
        hasMore: offset + items.length < total,
      };
    },
  },
};
```

### Client:
```graphql
query GetPosts($limit: Int!, $offset: Int!) {
  posts(limit: $limit, offset: $offset) {
    items {
      id
      title
      createdAt
    }
    total
    hasMore
  }
}
```

### Problems with offset:
- **Skip problem:** If new items are added, items shift — page 2 might show items already on page 1
- **Performance:** `OFFSET 10000` requires DB to read and discard 10,000 rows
- **Consistency:** Results change between page loads if data changes

---

## Cursor-Based Pagination (Recommended)

Cursors point to a specific position in the result set. Stable even with inserts/deletes.

### Schema (Relay Connection Spec):
```graphql
type Query {
  posts(
    first: Int      # forward: take N items
    after: String   # forward: after this cursor
    last: Int       # backward: take N items from the end
    before: String  # backward: before this cursor
  ): PostConnection!
}

type PostConnection {
  edges: [PostEdge!]!
  pageInfo: PageInfo!
  totalCount: Int!
}

type PostEdge {
  node: Post!
  cursor: String!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}
```

### Cursor Implementation:

```javascript
// Encode/decode cursor (base64 of ID or timestamp)
const encodeCursor = (value) => Buffer.from(String(value)).toString('base64');
const decodeCursor = (cursor) => Buffer.from(cursor, 'base64').toString('utf-8');

const resolvers = {
  Query: {
    posts: async (parent, { first = 20, after, last, before }, context) => {
      let query = context.db.posts.query().orderBy('created_at', 'DESC');

      if (after) {
        const cursor = decodeCursor(after);
        query = query.where('created_at', '<', cursor); // assuming cursor is a timestamp
      }

      if (before) {
        const cursor = decodeCursor(before);
        query = query.where('created_at', '>', cursor);
      }

      const limit = first || last || 20;
      const posts = await query.limit(limit + 1); // fetch 1 extra to detect hasNextPage

      const hasNextPage = posts.length > limit;
      const hasPreviousPage = !!after;
      const items = hasNextPage ? posts.slice(0, limit) : posts;

      const edges = items.map(post => ({
        node: post,
        cursor: encodeCursor(post.createdAt.toISOString()),
      }));

      return {
        edges,
        totalCount: await context.db.posts.count(),
        pageInfo: {
          hasNextPage,
          hasPreviousPage,
          startCursor: edges[0]?.cursor ?? null,
          endCursor: edges[edges.length - 1]?.cursor ?? null,
        },
      };
    },
  },
};
```

### Client Usage:
```graphql
# First page
query GetFirstPage {
  posts(first: 10) {
    edges {
      cursor
      node { id title createdAt }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
  }
}

# Next page
query GetNextPage($after: String!) {
  posts(first: 10, after: $after) {
    edges {
      cursor
      node { id title createdAt }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

```javascript
// Apollo Client: fetchMore
const { data, fetchMore } = useQuery(GET_POSTS, { variables: { first: 10 } });

const loadMore = () => {
  fetchMore({
    variables: {
      first: 10,
      after: data.posts.pageInfo.endCursor,
    },
    updateQuery: (prev, { fetchMoreResult }) => ({
      posts: {
        ...fetchMoreResult.posts,
        edges: [...prev.posts.edges, ...fetchMoreResult.posts.edges],
      },
    }),
  });
};
```

---

## ID-Based Cursor (Simpler)

Use the row ID as the cursor — works well for tables with sequential IDs:

```javascript
const resolvers = {
  Query: {
    posts: async (parent, { first = 20, after }, context) => {
      const afterId = after ? parseInt(decodeCursor(after)) : null;

      const posts = await context.db.query(`
        SELECT * FROM posts
        WHERE ($1::bigint IS NULL OR id < $1)
        ORDER BY id DESC
        LIMIT $2
      `, [afterId, first + 1]);

      const hasNextPage = posts.length > first;
      const items = hasNextPage ? posts.slice(0, first) : posts;

      return {
        edges: items.map(p => ({ node: p, cursor: encodeCursor(p.id) })),
        pageInfo: {
          hasNextPage,
          hasPreviousPage: !!after,
          startCursor: items[0] ? encodeCursor(items[0].id) : null,
          endCursor: items[items.length-1] ? encodeCursor(items[items.length-1].id) : null,
        },
        totalCount: await context.db.posts.count(),
      };
    },
  },
};
```

---

## Simple Page Number Pagination

Sometimes page numbers are better for UX (especially for search results):

```graphql
type Query {
  searchPosts(query: String!, page: Int = 1, perPage: Int = 20): SearchResult!
}

type SearchResult {
  posts: [Post!]!
  currentPage: Int!
  totalPages: Int!
  totalCount: Int!
}
```

```javascript
const resolvers = {
  Query: {
    searchPosts: async (parent, { query, page = 1, perPage = 20 }, context) => {
      const offset = (page - 1) * perPage;
      const [posts, totalCount] = await Promise.all([
        context.db.posts.search(query, { limit: perPage, offset }),
        context.db.posts.searchCount(query),
      ]);

      return {
        posts,
        currentPage: page,
        totalPages: Math.ceil(totalCount / perPage),
        totalCount,
      };
    },
  },
};
```

---

## Pagination Comparison

| Method | Stable? | Performance | UX | Use When |
|---|---|---|---|---|
| Offset | No | Poor for large offsets | Simple | Admin lists, small datasets |
| Cursor | Yes | Excellent | Load more, infinite scroll | Social feeds, real-time data |
| Page number | No | Poor for large pages | Familiar | Search results, tables |

---

## Cursor Pagination Best Practices

- Always encode cursors (base64) — clients should treat them as opaque strings
- Don't expose implementation details in cursor values (ID, timestamp)
- Support both `first/after` (forward) and `last/before` (backward)
- Always return `pageInfo.hasNextPage` accurately
- Limit maximum `first` value to prevent abuse: `const safeLimit = Math.min(first ?? 20, 100)`
- Include `totalCount` carefully — it can be expensive; cache it or skip if not needed
