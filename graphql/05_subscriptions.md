# Chapter 5: Subscriptions

## What are Subscriptions?

GraphQL Subscriptions enable **real-time data** from the server to the client. Unlike queries and mutations (request/response), subscriptions maintain a persistent connection and push data when events occur.

```
Client ──────── subscribe ──────────► Server
       ◄──── initial response ────────
       ◄──── event data ─────────────  (when something happens)
       ◄──── event data ─────────────  (again...)
       ◄──── event data ─────────────  (and again...)
       ──────── unsubscribe ──────────►
```

---

## Transport: WebSockets

Subscriptions typically use WebSockets. The `graphql-ws` protocol (modern) has replaced the older `subscriptions-transport-ws`:

```
Client → Server: { type: "connection_init" }
Server → Client: { type: "connection_ack" }
Client → Server: { type: "subscribe", id: "1", payload: { query: "subscription {...}" } }
Server → Client: { type: "next", id: "1", payload: { data: {...} } }
Server → Client: { type: "next", id: "1", payload: { data: {...} } }
Client → Server: { type: "complete", id: "1" }  (unsubscribe)
```

---

## Schema Definition

```graphql
type Subscription {
  # Post events
  postAdded: Post!
  postUpdated(id: ID!): Post!
  postDeleted: ID!                # returns deleted post ID

  # Chat
  messageSent(channelId: ID!): Message!
  userTyping(channelId: ID!): TypingEvent!
  userOnlineStatus(userId: ID!): UserStatus!

  # Notifications
  notificationReceived: Notification!

  # Order tracking
  orderStatusChanged(orderId: ID!): Order!
}

type TypingEvent {
  user: User!
  channelId: ID!
  isTyping: Boolean!
}

type UserStatus {
  user: User!
  online: Boolean!
  lastSeen: DateTime
}
```

---

## Server Implementation (Apollo Server + graphql-ws)

### Pub/Sub Setup:
```javascript
import { PubSub } from 'graphql-subscriptions';

// Simple in-memory pub/sub (use Redis for production multi-instance)
const pubsub = new PubSub();

// Event names as constants
const Events = {
  POST_ADDED: 'POST_ADDED',
  POST_UPDATED: 'POST_UPDATED',
  MESSAGE_SENT: 'MESSAGE_SENT',
  ORDER_STATUS_CHANGED: 'ORDER_STATUS_CHANGED',
};
```

### Subscription Resolvers:
```javascript
const resolvers = {
  Subscription: {
    postAdded: {
      subscribe: () => pubsub.asyncIterator([Events.POST_ADDED]),
      // resolver (optional — default passes through payload)
      resolve: (payload) => payload.postAdded,
    },

    postUpdated: {
      subscribe: (parent, { id }) => {
        // Only receive events for the specific post
        const topic = `${Events.POST_UPDATED}_${id}`;
        return pubsub.asyncIterator([topic]);
      },
    },

    messageSent: {
      subscribe: withFilter(
        () => pubsub.asyncIterator([Events.MESSAGE_SENT]),
        (payload, variables) => {
          // Only deliver if the message is in the requested channel
          return payload.messageSent.channelId === variables.channelId;
        }
      ),
    },
  },

  Mutation: {
    createPost: async (parent, { input }, context) => {
      const post = await context.db.posts.create(input);

      // Publish the event — all subscribers of POST_ADDED receive this
      await pubsub.publish(Events.POST_ADDED, { postAdded: post });

      return post;
    },

    sendMessage: async (parent, { input }, context) => {
      const message = await context.db.messages.create(input);
      await pubsub.publish(Events.MESSAGE_SENT, { messageSent: message });
      return message;
    },
  },
};
```

### Server Setup with graphql-ws:
```javascript
import { createServer } from 'http';
import { WebSocketServer } from 'ws';
import { useServer } from 'graphql-ws/lib/use/ws';
import { ApolloServer } from '@apollo/server';
import { makeExecutableSchema } from '@graphql-tools/schema';

const schema = makeExecutableSchema({ typeDefs, resolvers });
const httpServer = createServer(app);

// WebSocket server for subscriptions
const wsServer = new WebSocketServer({
  server: httpServer,
  path: '/graphql',
});

const serverCleanup = useServer(
  {
    schema,
    context: async (ctx) => {
      // Authentication for WebSocket connections
      const token = ctx.connectionParams?.authorization;
      const user = token ? await verifyToken(token) : null;
      return { user, pubsub, db };
    },
    onConnect: async (ctx) => {
      console.log('Client connected');
      return true;
    },
    onDisconnect: (ctx, code, reason) => {
      console.log('Client disconnected');
    },
  },
  wsServer
);

// Apollo Server setup
const server = new ApolloServer({
  schema,
  plugins: [
    ApolloServerPluginDrainHttpServer({ httpServer }),
    {
      async serverWillStart() {
        return {
          async drainServer() {
            await serverCleanup.dispose();
          },
        };
      },
    },
  ],
});

await server.start();
httpServer.listen(4000);
```

---

## Redis Pub/Sub for Production (Multi-Instance)

In-memory PubSub only works with a single server instance. Use Redis for horizontal scaling:

```javascript
import { RedisPubSub } from 'graphql-redis-subscriptions';
import Redis from 'ioredis';

const pubsub = new RedisPubSub({
  publisher: new Redis({ host: 'redis', port: 6379 }),
  subscriber: new Redis({ host: 'redis', port: 6379 }),
});
```

---

## withFilter — Filtering Events

```javascript
import { withFilter } from 'graphql-subscriptions';

const resolvers = {
  Subscription: {
    orderStatusChanged: {
      subscribe: withFilter(
        () => pubsub.asyncIterator([Events.ORDER_STATUS_CHANGED]),
        (payload, variables, context) => {
          // Only deliver if:
          // 1. The order ID matches
          // 2. The user owns this order
          return (
            payload.orderStatusChanged.id === variables.orderId &&
            payload.orderStatusChanged.userId === context.user.id
          );
        }
      ),
    },

    notificationReceived: {
      subscribe: withFilter(
        () => pubsub.asyncIterator([Events.NOTIFICATION]),
        (payload, variables, context) => {
          // Only deliver to the intended user
          return payload.notificationReceived.userId === context.user.id;
        }
      ),
    },
  },
};
```

---

## Client: Apollo Client Subscriptions

```javascript
import { ApolloClient, InMemoryCache, split, HttpLink } from '@apollo/client';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { createClient } from 'graphql-ws';
import { getMainDefinition } from '@apollo/client/utilities';

// HTTP link for queries/mutations
const httpLink = new HttpLink({ uri: 'http://localhost:4000/graphql' });

// WebSocket link for subscriptions
const wsLink = new GraphQLWsLink(
  createClient({
    url: 'ws://localhost:4000/graphql',
    connectionParams: {
      authorization: `Bearer ${localStorage.getItem('token')}`,
    },
  })
);

// Route to correct link based on operation type
const splitLink = split(
  ({ query }) => {
    const def = getMainDefinition(query);
    return def.kind === 'OperationDefinition' && def.operation === 'subscription';
  },
  wsLink,
  httpLink
);

const client = new ApolloClient({ link: splitLink, cache: new InMemoryCache() });
```

### Using useSubscription in React:
```javascript
import { useSubscription, gql } from '@apollo/client';

const POST_ADDED = gql`
  subscription OnPostAdded {
    postAdded {
      id
      title
      author { name avatar }
      createdAt
    }
  }
`;

function LiveFeed() {
  const [posts, setPosts] = useState([]);

  useSubscription(POST_ADDED, {
    onData: ({ data }) => {
      setPosts(prev => [data.data.postAdded, ...prev]);
    },
    onError: (error) => console.error('Subscription error:', error),
  });

  return (
    <ul>
      {posts.map(post => (
        <li key={post.id}>{post.title} by {post.author.name}</li>
      ))}
    </ul>
  );
}
```

---

## Subscription Lifecycle

```javascript
const resolvers = {
  Subscription: {
    userTyping: {
      subscribe: async function*(parent, { channelId }, context) {
        // Setup: register user as watching this channel
        await context.db.channelWatchers.add(channelId, context.user.id);

        try {
          const iterator = pubsub.asyncIterator([`TYPING_${channelId}`]);
          for await (const event of iterator) {
            yield event;
          }
        } finally {
          // Cleanup: runs when client disconnects
          await context.db.channelWatchers.remove(channelId, context.user.id);
        }
      },
    },
  },
};
```

---

## Authentication for Subscriptions

```javascript
useServer({
  schema,
  context: async (ctx, msg, args) => {
    // ctx.connectionParams contains the initial connection params
    const token = ctx.connectionParams?.authorization?.replace('Bearer ', '');

    if (!token) {
      throw new Error('Unauthorized: no token provided');
    }

    try {
      const user = jwt.verify(token, process.env.JWT_SECRET);
      return { user, pubsub, db };
    } catch {
      throw new Error('Unauthorized: invalid token');
    }
  },
}, wsServer);
```

---

## When to Use Subscriptions

| Use Case | Use | Alternative |
|---|---|---|
| Chat messages | Subscriptions | — |
| Live notifications | Subscriptions | SSE |
| Order tracking | Subscriptions | Polling |
| Dashboard metrics | Polling or SSE | Subscriptions |
| Collaborative editing | Subscriptions + CRDT | — |
| Real-time scores | Subscriptions | SSE |
| Periodic reports | Polling | Subscriptions (overkill) |

**SSE (Server-Sent Events)** is simpler for one-way server→client streams and works over HTTP/1.1 with no extra dependencies.
