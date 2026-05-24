# Chapter 15: GraphQL with React (Apollo Client)

## Apollo Client Setup

Apollo Client is the most popular GraphQL client for React.

```bash
npm install @apollo/client graphql
```

```javascript
// src/apollo.js
import { ApolloClient, InMemoryCache, ApolloProvider, createHttpLink } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';

const httpLink = createHttpLink({
  uri: 'http://localhost:4000/graphql',
});

// Add auth token to every request
const authLink = setContext((_, { headers }) => {
  const token = localStorage.getItem('token');
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : '',
    },
  };
});

export const client = new ApolloClient({
  link: authLink.concat(httpLink),
  cache: new InMemoryCache(),
  defaultOptions: {
    watchQuery: { errorPolicy: 'all' },
    query: { errorPolicy: 'all' },
  },
});

// Wrap your app
// src/index.js
root.render(
  <ApolloProvider client={client}>
    <App />
  </ApolloProvider>
);
```

---

## useQuery

```javascript
import { useQuery, gql } from '@apollo/client';

const GET_USERS = gql`
  query GetUsers($limit: Int!, $role: UserRole) {
    users(limit: $limit, role: $role) {
      id
      name
      email
      role
      avatar(size: 64)
    }
  }
`;

function UserList() {
  const { data, loading, error, refetch } = useQuery(GET_USERS, {
    variables: { limit: 20, role: 'USER' },
    fetchPolicy: 'cache-and-network',   // use cache but also refetch in background
    pollInterval: 30000,                // refetch every 30 seconds
    skip: !isLoggedIn,                  // skip query if not logged in
    onCompleted: (data) => console.log('Loaded', data.users.length, 'users'),
    onError: (error) => console.error('Failed to load users', error),
  });

  if (loading && !data) return <Spinner />;
  if (error) return <ErrorMessage error={error} />;

  return (
    <div>
      <button onClick={() => refetch()}>Refresh</button>
      {data.users.map(user => (
        <UserCard key={user.id} user={user} />
      ))}
    </div>
  );
}
```

### Fetch Policies:

| Policy | Cache | Network | When to Use |
|---|---|---|---|
| `cache-first` | ✓ always | Only on cache miss | Default — best performance |
| `cache-and-network` | ✓ immediately | Always refetches | Stale-while-revalidate |
| `network-only` | Writes to cache | Always fetches | Real-time accuracy needed |
| `no-cache` | Never | Always fetches | Sensitive data |
| `cache-only` | ✓ always | Never | Offline mode |

---

## useMutation

```javascript
import { useMutation, gql } from '@apollo/client';

const CREATE_POST = gql`
  mutation CreatePost($input: CreatePostInput!) {
    createPost(input: $input) {
      post {
        id
        title
        status
        createdAt
      }
      errors {
        field
        message
      }
    }
  }
`;

function CreatePostForm() {
  const [title, setTitle] = useState('');
  const [body, setBody] = useState('');
  const [fieldErrors, setFieldErrors] = useState({});

  const [createPost, { loading, error }] = useMutation(CREATE_POST, {
    // Update cache after mutation
    update(cache, { data: { createPost } }) {
      if (!createPost.post) return;
      cache.modify({
        fields: {
          posts(existing = []) {
            const newPostRef = cache.writeFragment({
              data: createPost.post,
              fragment: gql`fragment NewPost on Post { id title status createdAt }`,
            });
            return [newPostRef, ...existing];
          },
        },
      });
    },
    onCompleted: ({ createPost }) => {
      if (createPost.errors.length > 0) {
        const errors = {};
        createPost.errors.forEach(e => { errors[e.field] = e.message; });
        setFieldErrors(errors);
      } else {
        setTitle('');
        setBody('');
        alert('Post created!');
      }
    },
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    setFieldErrors({});
    createPost({ variables: { input: { title, body } } });
  };

  return (
    <form onSubmit={handleSubmit}>
      <input value={title} onChange={e => setTitle(e.target.value)} />
      {fieldErrors.title && <span>{fieldErrors.title}</span>}
      <textarea value={body} onChange={e => setBody(e.target.value)} />
      {fieldErrors.body && <span>{fieldErrors.body}</span>}
      <button type="submit" disabled={loading}>
        {loading ? 'Creating...' : 'Create Post'}
      </button>
      {error && <div>Network error: {error.message}</div>}
    </form>
  );
}
```

---

## Apollo Cache

Apollo Client's InMemoryCache is a normalized cache — it stores objects by type + ID.

```javascript
// Cache normalization
// Two queries returning the same user are merged in cache:
// Query 1: { user(id: "1") { id name } }
// Query 2: { me { id name email } }  (same user)
// Cache stores: User:1 → { id: "1", name: "Alice", email: "alice@..." }

const client = new ApolloClient({
  cache: new InMemoryCache({
    typePolicies: {
      Query: {
        fields: {
          // Custom field policy for pagination
          posts: {
            keyArgs: ['filter'],  // cache by filter only (not pagination args)
            merge(existing = [], incoming, { args }) {
              // Merge paginated results
              const merged = existing.slice();
              incoming.edges.forEach((edge, index) => {
                merged[args.offset + index] = edge;
              });
              return merged;
            },
          },
        },
      },
      User: {
        // Custom ID field (default is 'id')
        keyFields: ['email'],  // use email as unique identifier
      },
    },
  }),
});
```

---

## Cache Updates After Mutations

### Refetch queries (simplest):
```javascript
const [deletePost] = useMutation(DELETE_POST, {
  refetchQueries: ['GetPosts'],   // refetch by operation name
  // Or with variables:
  refetchQueries: [{ query: GET_POSTS, variables: { limit: 20 } }],
});
```

### Update cache manually (best performance):
```javascript
const [deletePost] = useMutation(DELETE_POST, {
  update(cache, { data }, { variables }) {
    cache.modify({
      fields: {
        posts(existingPosts, { readField }) {
          return existingPosts.filter(
            postRef => readField('id', postRef) !== variables.id
          );
        },
      },
    });
    // Also remove from normalized store
    cache.evict({ id: cache.identify({ __typename: 'Post', id: variables.id }) });
    cache.gc();
  },
});
```

---

## useSubscription

```javascript
import { useSubscription, gql } from '@apollo/client';

const POST_ADDED = gql`
  subscription OnPostAdded {
    postAdded {
      id
      title
      author { id name avatar }
      createdAt
    }
  }
`;

function LiveFeed() {
  const [posts, setPosts] = useState([]);

  // Load initial posts
  const { loading } = useQuery(GET_POSTS, {
    onCompleted: (data) => setPosts(data.posts),
  });

  // Subscribe to new posts
  useSubscription(POST_ADDED, {
    onData: ({ data }) => {
      setPosts(prev => [data.data.postAdded, ...prev]);
    },
  });

  if (loading) return <Spinner />;
  return <PostList posts={posts} />;
}
```

---

## Pagination with Apollo

```javascript
const GET_POSTS = gql`
  query GetPosts($first: Int!, $after: String) {
    posts(first: $first, after: $after) {
      edges { cursor node { id title excerpt } }
      pageInfo { hasNextPage endCursor }
    }
  }
`;

function InfinitePostList() {
  const { data, loading, fetchMore } = useQuery(GET_POSTS, {
    variables: { first: 10 },
  });

  const loadMore = () => {
    fetchMore({
      variables: {
        first: 10,
        after: data.posts.pageInfo.endCursor,
      },
      updateQuery: (prev, { fetchMoreResult }) => {
        if (!fetchMoreResult) return prev;
        return {
          posts: {
            ...fetchMoreResult.posts,
            edges: [...prev.posts.edges, ...fetchMoreResult.posts.edges],
          },
        };
      },
    });
  };

  if (loading && !data) return <Spinner />;

  return (
    <div>
      {data.posts.edges.map(({ node }) => (
        <PostCard key={node.id} post={node} />
      ))}
      {data.posts.pageInfo.hasNextPage && (
        <button onClick={loadMore} disabled={loading}>
          Load More
        </button>
      )}
    </div>
  );
}
```

---

## Code Generation with GraphQL Code Generator

Generate TypeScript types for your queries:

```bash
npm install -D @graphql-codegen/cli @graphql-codegen/client-preset
```

```yaml
# codegen.ts
import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  schema: 'http://localhost:4000/graphql',
  documents: ['src/**/*.tsx', 'src/**/*.ts'],
  generates: {
    './src/generated/': {
      preset: 'client',
      plugins: [],
    },
  },
};
export default config;
```

```bash
npx graphql-codegen --config codegen.ts
```

Use generated types:
```typescript
import { useQuery } from '@apollo/client';
import { graphql } from './generated';

const GET_USER = graphql(`
  query GetUser($id: ID!) {
    user(id: $id) {
      id
      name
      email
    }
  }
`);

function UserProfile({ id }: { id: string }) {
  const { data } = useQuery(GET_USER, { variables: { id } });
  // data is fully typed: data?.user?.name is string | undefined
  return <h1>{data?.user?.name}</h1>;
}
```
