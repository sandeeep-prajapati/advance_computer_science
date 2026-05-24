# Chapter 19: Apollo Federation

## What is Federation?

Apollo Federation lets you split a large GraphQL API into independent subgraphs — each owned by a different team — that compose into a single unified supergraph.

```
Without Federation:          With Federation:
Single monolithic schema     Multiple subgraph services
     │                            │
  One server                 Router (Supergraph)
  All resolvers              ┌────┴─────┐
                         Users     Products    Orders
                         Service   Service     Service
```

---

## Core Concepts

| Term | Definition |
|---|---|
| **Subgraph** | An independent GraphQL service that owns part of the schema |
| **Supergraph** | The unified schema composed from all subgraphs |
| **Router** | The entry point that receives client queries and fans out to subgraphs |
| **Entity** | A type shared across subgraphs, identified by `@key` |
| **Reference resolver** | A resolver that can reconstruct an entity from its key fields |

---

## Setting Up a Subgraph

```bash
npm install @apollo/subgraph graphql
```

### Users Subgraph:

```typescript
// users-service/schema.ts
import { buildSubgraphSchema } from '@apollo/subgraph';
import { gql } from 'graphql-tag';

const typeDefs = gql`
  extend schema
    @link(url: "https://specs.apollo.dev/federation/v2.3", import: ["@key"])

  type Query {
    me: User
    user(id: ID!): User
    users: [User!]!
  }

  # @key declares what makes a User uniquely identifiable
  # Other subgraphs can reference User by id
  type User @key(fields: "id") {
    id: ID!
    name: String!
    email: String!
    joinedAt: DateTime!
  }
`;

const resolvers = {
  Query: {
    me: (_, __, { user }) => user,
    user: (_, { id }, { db }) => db.users.findById(id),
    users: (_, __, { db }) => db.users.findAll(),
  },
  User: {
    // Reference resolver: called when another subgraph references a User by id
    __resolveReference: async ({ id }, { db }) => {
      return db.users.findById(id);
    },
  },
};

export const schema = buildSubgraphSchema({ typeDefs, resolvers });
```

---

### Products Subgraph:

```typescript
// products-service/schema.ts
const typeDefs = gql`
  extend schema
    @link(url: "https://specs.apollo.dev/federation/v2.3",
          import: ["@key", "@external", "@requires", "@provides"])

  type Query {
    product(id: ID!): Product
    products(category: String): [Product!]!
    featuredProducts: [Product!]!
  }

  type Product @key(fields: "id") {
    id: ID!
    name: String!
    price: Float!
    category: String!
    inStock: Boolean!
    description: String
  }
`;

const resolvers = {
  Query: {
    product: (_, { id }, { db }) => db.products.findById(id),
    products: (_, { category }, { db }) => db.products.findByCategory(category),
    featuredProducts: (_, __, { db }) => db.products.findFeatured(),
  },
  Product: {
    __resolveReference: ({ id }, { db }) => db.products.findById(id),
  },
};
```

---

### Orders Subgraph — References Entities from Other Subgraphs:

```typescript
// orders-service/schema.ts
const typeDefs = gql`
  extend schema
    @link(url: "https://specs.apollo.dev/federation/v2.3",
          import: ["@key", "@external"])

  type Query {
    order(id: ID!): Order
    myOrders: [Order!]!
  }

  type Order @key(fields: "id") {
    id: ID!
    status: OrderStatus!
    total: Float!
    placedAt: DateTime!
    # Reference to User entity — owned by Users subgraph
    customer: User!
    # Reference to multiple Product entities
    items: [OrderItem!]!
  }

  type OrderItem {
    product: Product!    # Product from Products subgraph
    quantity: Int!
    unitPrice: Float!
  }

  # Stub: User is owned by Users subgraph, but referenced here
  type User @key(fields: "id", resolvable: false) {
    id: ID!
  }

  # Stub: Product is owned by Products subgraph
  type Product @key(fields: "id", resolvable: false) {
    id: ID!
  }

  enum OrderStatus {
    PENDING
    PROCESSING
    SHIPPED
    DELIVERED
    CANCELLED
  }
`;

const resolvers = {
  Query: {
    order: (_, { id }, { db }) => db.orders.findById(id),
    myOrders: (_, __, { db, user }) => db.orders.findByCustomer(user.id),
  },
  Order: {
    __resolveReference: ({ id }, { db }) => db.orders.findById(id),
    customer: (order) => ({ __typename: 'User', id: order.customerId }),
    items: (order, _, { db }) => db.orderItems.findByOrder(order.id),
  },
  OrderItem: {
    product: (item) => ({ __typename: 'Product', id: item.productId }),
  },
};
```

---

## The Router (Apollo Router)

The router composes all subgraphs and serves clients a unified API.

```bash
# Install Apollo Router
curl -sSL https://router.apollo.dev/download/nix/latest | sh

# Or via npm (for development)
npm install @apollo/gateway
```

### Using Apollo Gateway (Node.js, simpler):

```typescript
// gateway/index.ts
import { ApolloServer } from '@apollo/server';
import { startStandaloneServer } from '@apollo/server/standalone';
import { ApolloGateway, IntrospectAndCompose } from '@apollo/gateway';

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs: [
      { name: 'users', url: 'http://users-service:4001/graphql' },
      { name: 'products', url: 'http://products-service:4002/graphql' },
      { name: 'orders', url: 'http://orders-service:4003/graphql' },
    ],
  }),
  // Pass headers from client to subgraphs
  buildService({ url }) {
    return new RemoteGraphQLDataSource({
      url,
      willSendRequest({ request, context }) {
        request.http?.headers.set('authorization', context.token || '');
        request.http?.headers.set('x-user-id', context.userId || '');
      },
    });
  },
});

const server = new ApolloServer({ gateway });

const { url } = await startStandaloneServer(server, {
  listen: { port: 4000 },
  context: async ({ req }) => ({
    token: req.headers.authorization,
    userId: getUserIdFromToken(req.headers.authorization),
  }),
});
```

### Apollo Router (production — binary):

```yaml
# router.yaml
supergraph:
  listen: 0.0.0.0:4000

subgraphs:
  users:
    routing_url: http://users-service:4001/graphql
  products:
    routing_url: http://products-service:4002/graphql
  orders:
    routing_url: http://orders-service:4003/graphql

# Pass JWT to all subgraphs
headers:
  all:
    request:
      - propagate:
          named: authorization
```

```bash
./router --config router.yaml --supergraph supergraph.graphql
```

---

## Composing the Supergraph Schema

```bash
# Install Rover CLI
npm install -g @apollo/rover

# Compose locally
rover supergraph compose --config supergraph.yaml

# supergraph.yaml:
federation_version: =2.3.0
subgraphs:
  users:
    routing_url: http://localhost:4001/graphql
    schema:
      file: ./users-service/schema.graphql
  products:
    routing_url: http://localhost:4002/graphql
    schema:
      file: ./products-service/schema.graphql
  orders:
    routing_url: http://localhost:4003/graphql
    schema:
      file: ./orders-service/schema.graphql
```

---

## Federation Directives

```graphql
# @key — declare entity, specify primary key fields
type Product @key(fields: "id") { ... }
type OrderItem @key(fields: "orderId productId") { ... }  # composite key

# @external — field owned by another subgraph
type User @key(fields: "id") {
  id: ID! @external
  email: String! @external   # owned by Users subgraph
  orders: [Order!]!          # owned by Orders subgraph
}

# @requires — need fields from another subgraph to resolve this field
type Product @key(fields: "id") {
  id: ID!
  price: Float!
  weight: Float! @external          # owned by another subgraph
  shippingCost: Float! @requires(fields: "weight")  # needs weight
}

# Resolver for shippingCost receives weight from Users subgraph:
Product: {
  shippingCost: (product) => calculateShipping(product.weight),
}

# @provides — tells router this subgraph can provide fields normally from another
type Order @key(fields: "id") {
  id: ID!
  product: Product @provides(fields: "name price")  # avoids extra roundtrip
}

type Product @key(fields: "id") {
  id: ID! @external
  name: String! @external
  price: Float! @external
}

# @shareable — multiple subgraphs can resolve this type/field
type Coordinates @shareable {
  lat: Float!
  lng: Float!
}

# @inaccessible — hide from client-facing schema (internal use)
type InternalAuditLog @inaccessible {
  ... 
}

# @override — migrate a field from one subgraph to another
type User @key(fields: "id") {
  id: ID!
  shippingAddress: Address @override(from: "orders")  # migrating this field here
}
```

---

## Query Planning

When a client requests data spanning multiple subgraphs, the router creates a query plan:

```graphql
# Client query spanning 3 subgraphs
query {
  myOrders {         # → Orders subgraph
    id
    status
    customer {       # → Users subgraph (entity lookup by id)
      name
      email
    }
    items {
      product {      # → Products subgraph (entity lookup by id)
        name
        price
        inStock
      }
      quantity
    }
  }
}
```

```
Query Plan:
1. Fetch myOrders from Orders subgraph
   → returns [{id, status, customerId, items: [{productId, quantity}]}]

2. Parallel fetch:
   a. Fetch users by ids=[customerId1, customerId2] from Users subgraph
   b. Fetch products by ids=[prod1, prod2, ...] from Products subgraph

3. Merge results and return to client
```

---

## Managed Federation with Apollo Studio

```bash
# Publish subgraph schema to Apollo Studio
rover subgraph publish my-graph@production \
  --name users \
  --schema ./users-service/schema.graphql \
  --routing-url https://users.internal/graphql

# Check for composition errors before deploying
rover subgraph check my-graph@production \
  --name users \
  --schema ./users-service/schema.graphql
```

Studio validates composition and detects breaking changes before you push.

---

## Docker Compose for Local Federation Development

```yaml
version: '3.8'
services:
  users-service:
    build: ./users-service
    ports:
      - "4001:4001"
    environment:
      DATABASE_URL: postgresql://postgres:postgres@postgres:5432/users
      PORT: 4001

  products-service:
    build: ./products-service
    ports:
      - "4002:4002"
    environment:
      DATABASE_URL: postgresql://postgres:postgres@postgres:5432/products
      PORT: 4002

  orders-service:
    build: ./orders-service
    ports:
      - "4003:4003"
    environment:
      DATABASE_URL: postgresql://postgres:postgres@postgres:5432/orders
      PORT: 4003

  router:
    image: ghcr.io/apollographql/router:v1.35.0
    ports:
      - "4000:4000"
    volumes:
      - ./router.yaml:/dist/config/router.yaml:ro
      - ./supergraph.graphql:/dist/config/supergraph.graphql:ro
    command: ["--config", "/dist/config/router.yaml", "--supergraph", "/dist/config/supergraph.graphql"]
    depends_on:
      - users-service
      - products-service
      - orders-service

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_PASSWORD: postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

---

## When to Use Federation

| Situation | Recommendation |
|---|---|
| Single team, one codebase | Don't use federation — monolith is simpler |
| 2–5 teams, growing API | Consider federation |
| Multiple teams, clear ownership | Federation is the right choice |
| Different tech stacks per service | Federation works across languages |
| Independent deploy cadences | Federation enables team autonomy |

### Federation is not for:
- Small projects with a single team
- When subgraph boundaries aren't clear
- When you want to avoid operational complexity

---

## Common Federation Pitfalls

```
❌ Circular @requires (A requires B, B requires A)
❌ Missing __resolveReference — entity can't be fetched
❌ Subgraph returning wrong __typename in references
❌ Large composite @key fields that are hard to look up
❌ Putting too many fields in @requires — triggers extra roundtrips
✓ Keep @key fields simple (usually just id)
✓ Return { __typename: 'User', id: user.id } from reference fields
✓ Use managed federation for production schema validation
```
