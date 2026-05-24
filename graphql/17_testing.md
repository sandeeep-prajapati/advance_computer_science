# Chapter 17: Testing GraphQL

## Testing Strategy

```
Unit Tests        → Resolver functions (business logic)
Integration Tests → Resolver with real database (via test DB or Testcontainers)
E2E Tests         → Full HTTP request through GraphQL server
Schema Tests      → Schema validity, breaking changes
```

---

## Unit Testing Resolvers

Resolvers are just functions — test them directly without a GraphQL server.

```typescript
// resolvers/user.resolver.ts
export const userResolvers = {
  Query: {
    user: async (_, { id }, { db, user: currentUser }) => {
      if (!currentUser) throw new AuthenticationError('Not authenticated');
      return db.users.findById(id);
    },
    users: async (_, { limit = 20, offset = 0 }, { db }) => {
      return db.users.findMany({ limit, offset });
    },
  },
  Mutation: {
    createUser: async (_, { input }, { db, user: currentUser }) => {
      if (!currentUser?.isAdmin) throw new ForbiddenError('Admins only');
      return db.users.create(input);
    },
  },
};
```

```typescript
// resolvers/user.resolver.test.ts
import { userResolvers } from './user.resolver';
import { AuthenticationError, ForbiddenError } from '../errors';

describe('user resolvers', () => {
  const mockDb = {
    users: {
      findById: jest.fn(),
      findMany: jest.fn(),
      create: jest.fn(),
    },
  };

  beforeEach(() => jest.clearAllMocks());

  describe('Query.user', () => {
    it('throws AuthenticationError when not authenticated', async () => {
      const context = { db: mockDb, user: null };
      await expect(
        userResolvers.Query.user({}, { id: '1' }, context)
      ).rejects.toThrow(AuthenticationError);
    });

    it('returns user when authenticated', async () => {
      const fakeUser = { id: '1', name: 'Alice', email: 'alice@example.com' };
      mockDb.users.findById.mockResolvedValue(fakeUser);

      const context = { db: mockDb, user: { id: '42' } };
      const result = await userResolvers.Query.user({}, { id: '1' }, context);

      expect(result).toEqual(fakeUser);
      expect(mockDb.users.findById).toHaveBeenCalledWith('1');
    });
  });

  describe('Mutation.createUser', () => {
    it('throws ForbiddenError when not admin', async () => {
      const context = { db: mockDb, user: { id: '1', isAdmin: false } };
      await expect(
        userResolvers.Mutation.createUser({}, { input: { name: 'Bob' } }, context)
      ).rejects.toThrow(ForbiddenError);
    });

    it('creates user when admin', async () => {
      const newUser = { id: '99', name: 'Bob', email: 'bob@example.com' };
      mockDb.users.create.mockResolvedValue(newUser);

      const context = { db: mockDb, user: { id: '1', isAdmin: true } };
      const result = await userResolvers.Mutation.createUser(
        {}, { input: { name: 'Bob', email: 'bob@example.com' } }, context
      );

      expect(result).toEqual(newUser);
    });
  });
});
```

---

## Integration Testing with Apollo Server

Test the full GraphQL pipeline including schema + resolvers.

```bash
npm install -D @apollo/server supertest jest @types/jest ts-jest
```

### Test server setup:

```typescript
// test/helpers/server.ts
import { ApolloServer } from '@apollo/server';
import { typeDefs } from '../../src/schema';
import { resolvers } from '../../src/resolvers';
import { createTestContext } from './context';

export function createTestServer() {
  return new ApolloServer({ typeDefs, resolvers });
}

export async function executeOperation(
  server: ApolloServer,
  query: string,
  variables?: Record<string, unknown>,
  contextOverrides?: Partial<Context>
) {
  const response = await server.executeOperation(
    { query, variables },
    { contextValue: createTestContext(contextOverrides) }
  );
  return response.body.kind === 'single' ? response.body.singleResult : null;
}
```

```typescript
// test/helpers/context.ts
export function createTestContext(overrides: Partial<Context> = {}): Context {
  return {
    user: null,
    db: testDb,     // real test database connection
    ...overrides,
  };
}
```

### Writing integration tests:

```typescript
// test/users.test.ts
import { ApolloServer } from '@apollo/server';
import { createTestServer, executeOperation } from './helpers/server';

describe('User GraphQL Operations', () => {
  let server: ApolloServer;

  beforeAll(async () => {
    server = createTestServer();
    await server.start();
    await seedTestDatabase();
  });

  afterAll(async () => {
    await server.stop();
    await cleanTestDatabase();
  });

  describe('query { users }', () => {
    const USERS_QUERY = `
      query GetUsers($limit: Int) {
        users(limit: $limit) {
          id
          name
          email
        }
      }
    `;

    it('returns list of users', async () => {
      const result = await executeOperation(server, USERS_QUERY, { limit: 5 });

      expect(result?.errors).toBeUndefined();
      expect(result?.data?.users).toHaveLength(5);
      expect(result?.data?.users[0]).toMatchObject({
        id: expect.any(String),
        name: expect.any(String),
        email: expect.any(String),
      });
    });

    it('requires authentication', async () => {
      const result = await executeOperation(
        server,
        USERS_QUERY,
        { limit: 5 },
        { user: null }   // unauthenticated
      );

      expect(result?.errors).toBeDefined();
      expect(result?.errors?.[0].extensions?.code).toBe('UNAUTHENTICATED');
    });
  });

  describe('mutation { createUser }', () => {
    const CREATE_USER = `
      mutation CreateUser($input: CreateUserInput!) {
        createUser(input: $input) {
          user {
            id
            name
            email
          }
          errors {
            field
            message
          }
        }
      }
    `;

    it('creates a user when authenticated as admin', async () => {
      const result = await executeOperation(
        server,
        CREATE_USER,
        { input: { name: 'New User', email: 'new@example.com' } },
        { user: { id: '1', isAdmin: true } }
      );

      expect(result?.errors).toBeUndefined();
      expect(result?.data?.createUser.user).toMatchObject({
        name: 'New User',
        email: 'new@example.com',
      });
      expect(result?.data?.createUser.errors).toHaveLength(0);
    });

    it('returns validation errors for invalid email', async () => {
      const result = await executeOperation(
        server,
        CREATE_USER,
        { input: { name: 'Bob', email: 'not-an-email' } },
        { user: { id: '1', isAdmin: true } }
      );

      expect(result?.data?.createUser.user).toBeNull();
      expect(result?.data?.createUser.errors).toContainEqual(
        expect.objectContaining({ field: 'email' })
      );
    });
  });
});
```

---

## HTTP-level Testing with Supertest

For testing the full HTTP layer (middleware, auth headers, CORS):

```typescript
// test/http.test.ts
import supertest from 'supertest';
import { buildApp } from '../src/app';   // your Express/Fastify app

describe('GraphQL HTTP endpoint', () => {
  let app: Express;

  beforeAll(async () => {
    app = await buildApp();
  });

  it('returns 200 for valid query', async () => {
    const response = await supertest(app)
      .post('/graphql')
      .set('Content-Type', 'application/json')
      .set('Authorization', `Bearer ${testToken}`)
      .send({
        query: `{ users { id name } }`,
      });

    expect(response.status).toBe(200);
    expect(response.body.errors).toBeUndefined();
    expect(response.body.data.users).toBeDefined();
  });

  it('handles introspection', async () => {
    const response = await supertest(app)
      .post('/graphql')
      .set('Content-Type', 'application/json')
      .send({ query: '{ __schema { types { name } } }' });

    expect(response.status).toBe(200);
    expect(response.body.data.__schema).toBeDefined();
  });

  it('returns 400 for malformed JSON', async () => {
    const response = await supertest(app)
      .post('/graphql')
      .set('Content-Type', 'application/json')
      .send('{ invalid json }');

    expect(response.status).toBe(400);
  });
});
```

---

## Snapshot Testing for Schema

Catch breaking changes to your schema:

```typescript
// test/schema.test.ts
import { printSchema, buildSchema } from 'graphql';
import { schema } from '../src/schema';

describe('Schema', () => {
  it('matches snapshot', () => {
    expect(printSchema(schema)).toMatchSnapshot();
  });

  it('is valid', () => {
    // buildSchema throws if schema has errors
    expect(() => buildSchema(printSchema(schema))).not.toThrow();
  });
});
```

```bash
# First run creates __snapshots__/schema.test.ts.snap
# Subsequent runs compare against snapshot
# Update snapshots when you intentionally change the schema:
jest --updateSnapshot
```

---

## Testing with Testcontainers (Real Database)

```bash
npm install -D testcontainers @testcontainers/postgresql
```

```typescript
// test/setup/database.ts
import { PostgreSqlContainer } from '@testcontainers/postgresql';
import { Pool } from 'pg';

let container: PostgreSqlContainer;
let pool: Pool;

export async function startTestDatabase() {
  container = await new PostgreSqlContainer('postgres:16-alpine')
    .withDatabase('testdb')
    .withUsername('test')
    .withPassword('test')
    .start();

  pool = new Pool({ connectionString: container.getConnectionUri() });
  
  // Run migrations
  await runMigrations(pool);
  
  return pool;
}

export async function stopTestDatabase() {
  await pool.end();
  await container.stop();
}

export async function clearDatabase() {
  await pool.query('TRUNCATE users, posts, comments CASCADE');
}
```

```typescript
// test/posts.integration.test.ts
import { startTestDatabase, stopTestDatabase, clearDatabase } from './setup/database';

describe('Posts integration tests', () => {
  let db: Pool;

  beforeAll(async () => { db = await startTestDatabase(); });
  afterAll(async () => { await stopTestDatabase(); });
  beforeEach(async () => { await clearDatabase(); });

  it('creates and retrieves a post', async () => {
    const user = await db.query(
      'INSERT INTO users (name, email) VALUES ($1, $2) RETURNING *',
      ['Alice', 'alice@example.com']
    );

    const server = createTestServer();
    await server.start();

    const result = await executeOperation(
      server,
      `mutation { createPost(input: { title: "Hello", authorId: "${user.rows[0].id}" }) { post { id title } } }`,
      {},
      { user: user.rows[0], db: createDbClient(pool) }
    );

    expect(result?.data?.createPost.post.title).toBe('Hello');
    await server.stop();
  });
});
```

---

## Testing Subscriptions

```typescript
// test/subscriptions.test.ts
import { createPubSub } from 'graphql-subscriptions';
import { postAddedResolver } from '../src/resolvers/subscriptions';

describe('Subscription resolvers', () => {
  it('publishes post added event', async () => {
    const pubsub = createPubSub();
    const context = { pubsub };

    // Get the async iterator
    const iterator = await postAddedResolver.subscribe({}, {}, context, {});

    const newPost = { id: '1', title: 'Hello', authorId: '42' };

    // Publish event
    pubsub.publish('POST_ADDED', { postAdded: newPost });

    // Read from iterator
    const { value } = await iterator.next();
    expect(value.postAdded).toEqual(newPost);

    // Cleanup
    iterator.return?.();
  });
});
```

---

## Mocking External Services

```typescript
// test/mocks/stripe.mock.ts
export const mockStripe = {
  charges: {
    create: jest.fn().mockResolvedValue({ id: 'ch_test_123', status: 'succeeded' }),
  },
  customers: {
    create: jest.fn().mockResolvedValue({ id: 'cus_test_456' }),
  },
};

// In test:
jest.mock('stripe', () => () => mockStripe);

it('processes payment', async () => {
  const result = await executeOperation(server, PROCESS_PAYMENT_MUTATION, {
    input: { amount: 1000, token: 'tok_test' },
  });

  expect(mockStripe.charges.create).toHaveBeenCalledWith({
    amount: 1000,
    currency: 'usd',
    source: 'tok_test',
  });
  expect(result?.data?.processPayment.success).toBe(true);
});
```

---

## Testing DataLoader

```typescript
// test/dataloader.test.ts
import { createDataLoaders } from '../src/loaders';

describe('DataLoader', () => {
  it('batches user loads', async () => {
    const findUsers = jest.fn().mockResolvedValue([
      { id: '1', name: 'Alice' },
      { id: '2', name: 'Bob' },
    ]);

    const loaders = createDataLoaders({ db: { users: { findByIds: findUsers } } });

    // Load multiple users — should batch into one call
    const [alice, bob] = await Promise.all([
      loaders.userLoader.load('1'),
      loaders.userLoader.load('2'),
    ]);

    expect(findUsers).toHaveBeenCalledTimes(1);
    expect(findUsers).toHaveBeenCalledWith(['1', '2']);
    expect(alice.name).toBe('Alice');
    expect(bob.name).toBe('Bob');
  });

  it('caches within request', async () => {
    const findUsers = jest.fn().mockResolvedValue([{ id: '1', name: 'Alice' }]);
    const loaders = createDataLoaders({ db: { users: { findByIds: findUsers } } });

    await loaders.userLoader.load('1');
    await loaders.userLoader.load('1'); // same ID

    expect(findUsers).toHaveBeenCalledTimes(1); // cached
  });
});
```

---

## Schema Change Detection

Prevent accidental breaking changes in CI:

```bash
npm install -D @graphql-inspector/cli
```

```bash
# Compare local schema against production
graphql-inspector diff \
  'http://api.production.com/graphql' \
  './schema.graphql'

# In CI pipeline:
graphql-inspector diff \
  'git:origin/main:schema.graphql' \
  './schema.graphql'
```

```yaml
# .github/workflows/schema-check.yml
- name: Check for breaking schema changes
  run: npx graphql-inspector diff "${{ secrets.PROD_GRAPHQL_URL }}" schema.graphql
  env:
    GRAPHQL_INSPECTOR_FAIL_ON_BREAKING: true
```

---

## Test Coverage Checklist

```
For each query:
  ✓ Returns correct data shape
  ✓ Applies correct auth (throws when unauthenticated)
  ✓ Handles not-found (returns null or throws)
  ✓ Pagination works correctly

For each mutation:
  ✓ Creates/updates/deletes the resource
  ✓ Returns validation errors for invalid input
  ✓ Is idempotent where expected
  ✓ Auth is enforced

General:
  ✓ Schema snapshot test
  ✓ No N+1 queries (check DataLoader)
  ✓ Error format matches expected shape
  ✓ Subscription events fire correctly
```
