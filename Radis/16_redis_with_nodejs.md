# Chapter 16: Redis with Node.js

## Libraries

| Library | Notes |
|---|---|
| `ioredis` | Feature-rich, Cluster + Sentinel, TypeScript support — recommended |
| `node-redis` | Official Redis client (v4+), modern async API |

```bash
npm install ioredis
npm install redis      # node-redis v4
```

---

## ioredis: Connecting

```typescript
import Redis from "ioredis";

// Single instance
const redis = new Redis({
  host:            "localhost",
  port:            6379,
  password:        "secret",
  db:              0,
  retryStrategy:   (times) => Math.min(times * 50, 2000),
  lazyConnect:     false,
});

// From URL
const redis = new Redis("redis://:password@localhost:6379/0");
const redis = new Redis("rediss://localhost:6380/0");  // TLS

// Sentinel
const redis = new Redis({
  sentinels: [
    { host: "sentinel1", port: 26379 },
    { host: "sentinel2", port: 26380 },
  ],
  name:         "myredis",
  password:     "secret",
  sentinelPassword: "sentinel-password",
});

// Cluster
const cluster = new Redis.Cluster([
  { host: "node1", port: 7001 },
  { host: "node2", port: 7002 },
  { host: "node3", port: 7003 },
]);
```

---

## All Data Types (ioredis)

```typescript
const r = new Redis({ host: "localhost" });

// Strings
await r.set("name", "Alice");
await r.get("name");                              // "Alice"
await r.setex("session:abc", 3600, "user-data"); // with TTL
await r.incr("counter");
await r.mset("k1", "v1", "k2", "v2");
await r.mget("k1", "k2");                        // ["v1", "v2"]

// Lists
await r.lpush("queue", "a", "b", "c");
await r.lrange("queue", 0, -1);                  // ["c","b","a"]
await r.brpop("queue", 5);                       // blocking pop

// Sets
await r.sadd("tags", "python", "redis");
await r.smembers("tags");
await r.sismember("tags", "python");             // 1

// Hashes
await r.hset("user:1000", "name", "Alice", "age", "30");
await r.hget("user:1000", "name");               // "Alice"
await r.hgetall("user:1000");                    // { name: "Alice", age: "30" }
await r.hmget("user:1000", ["name", "age"]);

// Sorted Sets
await r.zadd("leaderboard", 9850, "alice", 7200, "bob");
await r.zrevrange("leaderboard", 0, 9, "WITHSCORES");
await r.zrank("leaderboard", "alice");
await r.zscore("leaderboard", "alice");          // "9850"

// Expiry
await r.expire("key", 300);
await r.ttl("key");
await r.persist("key");
```

---

## Pipelining

```typescript
// Non-transactional pipeline
const pipe = r.pipeline();
for (let i = 0; i < 1000; i++) {
  pipe.set(`key:${i}`, `value:${i}`);
  pipe.expire(`key:${i}`, 3600);
}
const results = await pipe.exec();
// results: [[null, "OK"], [null, 1], [null, "OK"], ...]

// Transactional pipeline (MULTI/EXEC)
const [err, results2] = await r
  .multi()
  .incr("counter")
  .get("counter")
  .exec();
```

---

## Transactions with WATCH

```typescript
async function transferBalance(
  fromKey: string,
  toKey: string,
  amount: number
): Promise<boolean> {
  const result = await r.watch(fromKey, toKey, async (err, multi) => {
    const from = parseFloat((await r.get(fromKey)) ?? "0");
    const to   = parseFloat((await r.get(toKey))   ?? "0");

    if (from < amount) return null;

    return multi
      .set(fromKey, (from - amount).toString())
      .set(toKey,   (to   + amount).toString())
      .exec();
  });

  return result !== null;
}
```

---

## Lua Scripts

```typescript
// Define and call a Lua script
const rateLimitScript = r.defineCommand("rateLimit", {
  numberOfKeys: 1,
  lua: `
    local current = tonumber(redis.call('INCR', KEYS[1]))
    if current == 1 then
      redis.call('EXPIRE', KEYS[1], ARGV[2])
    end
    return current <= tonumber(ARGV[1]) and 1 or 0
  `,
});

// Type-safe call
const allowed = await (r as any).rateLimit("rate:user-1", 100, 60);

// Or use EVAL directly
const result = await r.eval(
  `return redis.call('GET', KEYS[1])`,
  1,
  "my-key"
);
```

---

## Pub/Sub

```typescript
import Redis from "ioredis";

// Use separate connections for pub and sub
const pub = new Redis();
const sub = new Redis();

// Subscribe
sub.subscribe("notifications", "alerts", (err, count) => {
  console.log(`Subscribed to ${count} channels`);
});

sub.on("message", (channel: string, message: string) => {
  console.log(`[${channel}] ${message}`);
});

// Pattern subscribe
sub.psubscribe("user:*:events");
sub.on("pmessage", (pattern, channel, message) => {
  console.log(`Pattern ${pattern} → [${channel}] ${message}`);
});

// Publish
await pub.publish("notifications", JSON.stringify({ type: "alert", msg: "Server down" }));

// Unsubscribe
sub.unsubscribe("notifications");
```

---

## node-redis (Official v4)

```typescript
import { createClient } from "redis";

const client = createClient({
  url:      "redis://localhost:6379",
  password: "secret",
});

client.on("error", (err) => console.error("Redis error:", err));
await client.connect();

// Commands
await client.set("key", "value");
const val = await client.get("key");

// Hash
await client.hSet("user:1", { name: "Alice", age: "30" });
const user = await client.hGetAll("user:1");

// Sorted Set
await client.zAdd("lb", [{ score: 9850, value: "alice" }]);
const top = await client.zRangeWithScores("lb", 0, 9, { REV: true });

// Pipeline
const multi = client.multi();
multi.set("a", "1");
multi.incr("a");
const results = await multi.exec();

await client.disconnect();
```

---

## Real-World Patterns

### Express.js Session Middleware

```typescript
import express from "express";
import session from "express-session";
import { createClient } from "redis";
import connectRedis from "connect-redis";

const app     = express();
const client  = createClient({ url: "redis://localhost:6379" });
await client.connect();

const RedisStore = connectRedis(session);

app.use(session({
  store:  new RedisStore({ client }),
  secret: "your-session-secret",
  resave: false,
  saveUninitialized: false,
  cookie: { secure: true, maxAge: 3600000 },
}));

app.get("/login", (req, res) => {
  (req.session as any).userId = "user-123";
  res.json({ status: "logged in" });
});
```

### Cache Middleware

```typescript
import { Request, Response, NextFunction } from "express";

function cacheMiddleware(ttl: number) {
  return async (req: Request, res: Response, next: NextFunction) => {
    const key = `cache:${req.path}:${JSON.stringify(req.query)}`;
    const cached = await r.get(key);

    if (cached) {
      return res.json(JSON.parse(cached));
    }

    const originalJson = res.json.bind(res);
    res.json = (data: any) => {
      r.setex(key, ttl, JSON.stringify(data));
      return originalJson(data);
    };

    next();
  };
}

app.get("/api/products", cacheMiddleware(300), async (req, res) => {
  const products = await db.getProducts();
  res.json(products);
});
```
