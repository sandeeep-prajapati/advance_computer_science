# Chapter 7: Pub/Sub & Messaging

## Pub/Sub

Redis Pub/Sub implements a simple publish-subscribe messaging model. Publishers send messages to channels; subscribers receive all messages on channels they subscribe to.

```
Publisher ──PUBLISH──► Channel "news:tech" ──► Subscriber A
                                           ──► Subscriber B
                                           ──► Subscriber C
```

**Key properties:**
- Fire-and-forget: messages not stored, not replayable
- No message acknowledgment
- Subscribers must be connected to receive messages
- Scales well for real-time notifications, chat, broadcasts

---

## Basic Pub/Sub Commands

### Subscriber

```bash
# Subscribe to a channel
SUBSCRIBE news:tech

# Subscribe to multiple channels
SUBSCRIBE news:tech news:sports alerts

# Pattern subscribe (glob patterns)
PSUBSCRIBE news:*       # all news channels
PSUBSCRIBE user:*:notifications
PSUBSCRIBE "h?llo"      # hello, hallo, hxllo

# Unsubscribe
UNSUBSCRIBE news:tech
PUNSUBSCRIBE news:*
```

### Publisher

```bash
# Publish a message
PUBLISH news:tech "Redis 8.0 released!"
# Returns: number of subscribers that received the message

PUBLISH user:1000:notifications '{"type":"mention","from":"alice"}'
```

### Inspect Channels

```bash
# List active channels (with subscribers)
PUBSUB CHANNELS              # all active channels
PUBSUB CHANNELS "news:*"     # matching pattern

# Count subscribers per channel
PUBSUB NUMSUB news:tech news:sports

# Count pattern subscribers
PUBSUB NUMPAT
```

---

## Pub/Sub in Python

```python
import redis
import threading
import json

r = redis.Redis(host="localhost", port=6379)


# ─── Subscriber ───
def message_handler(message):
    if message["type"] == "message":
        channel = message["channel"].decode()
        data    = json.loads(message["data"])
        print(f"[{channel}] {data}")

pubsub = r.pubsub()
pubsub.subscribe(**{
    "notifications:user-1": message_handler,
    "broadcast:all":        message_handler,
})

# Run listener in background thread
thread = pubsub.run_in_thread(sleep_time=0.01, daemon=True)


# ─── Publisher ───
def notify_user(user_id: str, event: dict) -> None:
    r.publish(f"notifications:{user_id}", json.dumps(event))

notify_user("user-1", {"type": "message", "from": "alice", "text": "Hello!"})

# Pattern subscription
pubsub2 = r.pubsub()
pubsub2.psubscribe(**{"notifications:*": message_handler})
```

---

## Pub/Sub Limitations

| Limitation | Workaround |
|---|---|
| No message persistence | Use Redis Streams (Chapter 4 / below) |
| No acknowledgment | Use Redis Streams with consumer groups |
| No replay | Use Redis Streams |
| Subscriber must be online | Use Streams or List queues |
| No filtering | Pattern subscribe or app-level filtering |

---

## Redis Streams as a Message Queue

Redis Streams (XADD/XREADGROUP) are the production-grade alternative to Pub/Sub when you need persistence, consumer groups, and acknowledgment.

### Producer Pattern

```python
import redis
import json
import time

r = redis.Redis()

def publish_event(stream: str, event: dict) -> str:
    return r.xadd(stream, event, maxlen=100000, approximate=True)

# Produce events
msg_id = publish_event("orders", {
    "order_id": "ORD-123",
    "user_id":  "user-42",
    "amount":   "99.99",
})
print(f"Published: {msg_id}")
```

### Consumer Group Pattern

```python
import redis
import json

r = redis.Redis()

STREAM  = "orders"
GROUP   = "order-processor"
CONSUMER = "worker-1"

# Create group (only once)
try:
    r.xgroup_create(STREAM, GROUP, id="0", mkstream=True)
except redis.exceptions.ResponseError:
    pass  # group already exists

def process_orders():
    while True:
        # Read up to 10 undelivered messages
        results = r.xreadgroup(
            groupname  = GROUP,
            consumername = CONSUMER,
            streams    = {STREAM: ">"},
            count      = 10,
            block      = 2000,   # block 2 seconds if no messages
        )

        if not results:
            continue

        for stream_name, messages in results:
            for msg_id, fields in messages:
                try:
                    order_id = fields[b"order_id"].decode()
                    amount   = float(fields[b"amount"])
                    print(f"Processing order {order_id}: ${amount}")
                    process_order(order_id, amount)

                    # Acknowledge: remove from pending
                    r.xack(STREAM, GROUP, msg_id)

                except Exception as e:
                    print(f"Failed: {e} — will retry")
                    # Don't ack — will be redelivered or claimed

process_orders()
```

### Handle Pending (Unacked) Messages

```python
# Check pending messages
pending = r.xpending(STREAM, GROUP, "-", "+", 10)
for p in pending:
    msg_id, consumer, idle_ms, delivery_count = p
    if idle_ms > 30000:  # stale for 30+ seconds
        # Claim and reprocess
        claimed = r.xclaim(STREAM, GROUP, CONSUMER, 30000, [msg_id])
        for mid, fields in claimed:
            reprocess(mid, fields)
```

---

## List as a Simple Queue

For simple task queues without consumer groups:

```python
QUEUE = "tasks:email"

# Producer
def enqueue(task: dict) -> None:
    r.lpush(QUEUE, json.dumps(task))

# Consumer (blocking)
def worker():
    while True:
        result = r.brpop(QUEUE, timeout=5)
        if result:
            _, data = result
            task = json.loads(data)
            process_task(task)

# Priority queue with multiple lists
HIGH_QUEUE = "tasks:high"
LOW_QUEUE  = "tasks:low"

# BRPOP waits on multiple lists in order (high priority first)
result = r.brpop([HIGH_QUEUE, LOW_QUEUE], timeout=10)
```

---

## Choosing Between Pub/Sub, Streams, and Lists

| Feature | Pub/Sub | Streams | List Queue |
|---|---|---|---|
| Persistence | No | Yes | Yes |
| Replay | No | Yes | No |
| Multiple consumers | All get each msg | Groups = divided | Round-robin |
| Acknowledgment | No | Yes (XACK) | Manual |
| Consumer groups | No | Yes | No |
| Message ordering | FIFO (channel) | FIFO (stream) | FIFO |
| Use case | Live notifications | Event sourcing, durable queue | Simple task queue |
