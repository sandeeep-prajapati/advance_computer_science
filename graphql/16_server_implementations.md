# Chapter 16: Server Implementations

## Overview

Multiple frameworks and languages support GraphQL. The best choice depends on your stack.

| Server | Language | Style | Best For |
|---|---|---|---|
| Apollo Server | Node.js | Feature-rich, enterprise | Large teams, rich ecosystem |
| GraphQL Yoga | Node.js | Lightweight, standards-based | Modern APIs, edge runtime |
| Pothos | Node.js | Code-first, TypeScript-native | Type safety without schema writing |
| NestJS + GraphQL | Node.js | Framework, decorators | Enterprise Node.js |
| Spring Boot GraphQL | Java/Kotlin | Annotation-based | Java ecosystem |
| Strawberry | Python | Decorator-based, type hints | Python, FastAPI-like |
| graphql-ruby | Ruby | Ruby-native | Rails apps |
| gqlgen | Go | Code-first, type-safe | Go microservices |

---

## Apollo Server (Node.js)

```bash
npm install @apollo/server graphql
```

```javascript
import { ApolloServer } from '@apollo/server';
import { startStandaloneServer } from '@apollo/server/standalone';

const typeDefs = `#graphql
  type Query {
    hello: String!
    user(id: ID!): User
  }

  type User {
    id: ID!
    name: String!
  }
`;

const resolvers = {
  Query: {
    hello: () => 'Hello, World!',
    user: (_, { id }, context) => context.db.users.findById(id),
  },
};

const server = new ApolloServer({
  typeDefs,
  resolvers,
  plugins: [
    ApolloServerPluginLandingPageLocalDefault(), // GraphQL sandbox in dev
  ],
  formatError: (formattedError, error) => {
    // Sanitize errors in production
    if (process.env.NODE_ENV === 'production' && !formattedError.extensions?.code) {
      return { message: 'Internal server error', extensions: { code: 'INTERNAL_SERVER_ERROR' } };
    }
    return formattedError;
  },
});

const { url } = await startStandaloneServer(server, {
  listen: { port: 4000 },
  context: async ({ req }) => ({
    user: await getUser(req.headers.authorization),
    db,
  }),
});

console.log(`Server running at ${url}`);
```

### With Express (more control):
```javascript
import { ApolloServer } from '@apollo/server';
import { expressMiddleware } from '@apollo/server/express4';
import express from 'express';
import cors from 'cors';

const app = express();
app.use(express.json());
app.use(cors());

const server = new ApolloServer({ typeDefs, resolvers });
await server.start();

app.use('/graphql', expressMiddleware(server, {
  context: async ({ req }) => ({
    user: await getUser(req.headers.authorization),
    db,
  }),
}));

app.listen(4000, () => console.log('Server on port 4000'));
```

---

## GraphQL Yoga (Node.js — Lightweight)

```bash
npm install graphql-yoga graphql
```

```javascript
import { createSchema, createYoga } from 'graphql-yoga';
import { createServer } from 'http';

const yoga = createYoga({
  schema: createSchema({
    typeDefs: `
      type Query { hello: String! }
    `,
    resolvers: {
      Query: { hello: () => 'Hello, Yoga!' },
    },
  }),
  context: ({ request }) => ({
    user: getUser(request.headers.get('authorization')),
    db,
  }),
});

createServer(yoga).listen(4000);
```

Yoga works on Node.js, Cloudflare Workers, Deno, Bun — highly portable.

---

## Pothos (Code-First TypeScript)

Pothos builds the schema in TypeScript — no SDL writing, full type safety:

```bash
npm install @pothos/core graphql
```

```typescript
import SchemaBuilder from '@pothos/core';

const builder = new SchemaBuilder<{
  Context: { user: User | null; db: DB };
}>({});

// Define types with code (no SDL)
const UserType = builder.objectType('User', {
  fields: (t) => ({
    id: t.exposeID('id'),
    name: t.exposeString('name'),
    email: t.exposeString('email'),
    posts: t.field({
      type: [PostType],
      resolve: (user, _, context) => context.db.posts.findByAuthor(user.id),
    }),
  }),
});

const PostType = builder.objectType('Post', {
  fields: (t) => ({
    id: t.exposeID('id'),
    title: t.exposeString('title'),
    author: t.field({
      type: UserType,
      resolve: (post, _, context) => context.db.users.findById(post.authorId),
    }),
  }),
});

builder.queryType({
  fields: (t) => ({
    user: t.field({
      type: UserType,
      nullable: true,
      args: { id: t.arg.id({ required: true }) },
      resolve: (root, { id }, context) => context.db.users.findById(id),
    }),
  }),
});

builder.mutationType({
  fields: (t) => ({
    createPost: t.field({
      type: PostType,
      args: { title: t.arg.string({ required: true }) },
      resolve: (root, { title }, context) => {
        if (!context.user) throw new Error('Unauthorized');
        return context.db.posts.create({ title, authorId: context.user.id });
      },
    }),
  }),
});

export const schema = builder.toSchema();
```

---

## NestJS with GraphQL

```bash
npm install @nestjs/graphql @nestjs/apollo @apollo/server graphql
```

```typescript
// app.module.ts
@Module({
  imports: [
    GraphQLModule.forRoot<ApolloDriverConfig>({
      driver: ApolloDriver,
      autoSchemaFile: 'schema.gql',  // code-first: auto-generates SDL
      sortSchema: true,
      context: ({ req }) => ({ req }),
    }),
    UsersModule,
    PostsModule,
  ],
})
export class AppModule {}

// user.resolver.ts
@Resolver(() => User)
export class UserResolver {
  constructor(
    private usersService: UsersService,
    private postsService: PostsService,
  ) {}

  @Query(() => [User])
  async users(): Promise<User[]> {
    return this.usersService.findAll();
  }

  @Query(() => User, { nullable: true })
  async user(@Args('id', { type: () => ID }) id: string): Promise<User | null> {
    return this.usersService.findById(id);
  }

  @ResolveField(() => [Post])
  async posts(@Parent() user: User): Promise<Post[]> {
    return this.postsService.findByAuthor(user.id);
  }

  @Mutation(() => User)
  @UseGuards(GqlAuthGuard)
  async createUser(
    @Args('input') input: CreateUserInput,
  ): Promise<User> {
    return this.usersService.create(input);
  }
}

// user.entity.ts
@ObjectType()
export class User {
  @Field(() => ID)
  id: string;

  @Field()
  name: string;

  @Field()
  email: string;

  @Field(() => [Post])
  posts: Post[];
}

// create-user.input.ts
@InputType()
export class CreateUserInput {
  @Field()
  @IsEmail()
  email: string;

  @Field()
  @MinLength(2)
  name: string;
}
```

---

## Spring Boot with GraphQL (Java)

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-graphql</artifactId>
</dependency>
```

```graphql
# src/main/resources/graphql/schema.graphqls
type Query {
  user(id: ID!): User
  users: [User!]!
}

type Mutation {
  createUser(input: CreateUserInput!): User!
}

type User {
  id: ID!
  name: String!
  posts: [Post!]!
}
```

```java
@Controller
public class UserController {

    private final UserService userService;

    @QueryMapping
    public Optional<User> user(@Argument Long id) {
        return userService.findById(id);
    }

    @QueryMapping
    public List<User> users() {
        return userService.findAll();
    }

    @MutationMapping
    public User createUser(@Argument @Valid CreateUserInput input) {
        return userService.create(input);
    }

    @SchemaMapping
    public List<Post> posts(User user) {
        return postService.findByAuthor(user.getId());
    }
}
```

---

## gqlgen (Go)

```bash
go get github.com/99designs/gqlgen
go run github.com/99designs/gqlgen generate
```

```go
// resolver.go
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
    return r.db.FindUser(ctx, id)
}

func (r *userResolver) Posts(ctx context.Context, obj *model.User) ([]*model.Post, error) {
    return r.db.FindPostsByAuthor(ctx, obj.ID)
}
```

Go's strong typing makes gqlgen generate type-safe resolver interfaces from your schema.

---

## Choosing Your Server

| Need | Recommendation |
|---|---|
| Node.js, flexibility | Apollo Server |
| Node.js, edge/serverless | GraphQL Yoga |
| TypeScript, type safety | Pothos + Yoga |
| Node.js enterprise, decorators | NestJS |
| Java/Kotlin | Spring Boot GraphQL |
| Python | Strawberry or Ariadne |
| Go | gqlgen |
| Ruby on Rails | graphql-ruby |
