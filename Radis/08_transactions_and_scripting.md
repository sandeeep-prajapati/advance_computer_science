# Chapter 8: Transactions & Scripting

## MULTI / EXEC Transactions

Redis transactions queue commands and execute them atomically as a single batch. Either all commands execute or none do (on EXEC failure).

```bash
# Basic transaction
MULTI                   # start transaction
SET balance 100
INCR counter
ZADD leaderboard 9850 "alice"
EXEC                    # execute all atomically

# Discard transaction
MULTI
SET key1 "value1"
DISCARD                 # cancel, queue is cleared
```

**Important nuances:**
- Redis transactions are NOT rollback-capable — if a command fails after EXEC, other commands still run
- Command errors detected at EXEC time (e.g., wrong type) don't abort the transaction
- Commands are queued, not executed immediately

```bash
MULTI
SET counter 0
INCR counter         # queued as QUEUED
LPUSH counter "x"    # will error (wrong type) but others still run
EXEC
# 1) OK
# 2) (integer) 1
# 3) (error) WRONGTYPE ...  ← only this failed
```

---

## WATCH (Optimistic Locking)

WATCH implements optimistic concurrency control. If any watched key changes before EXEC, the transaction is aborted.

```bash
# Pattern: check-and-set
WATCH balance
current = GET balance          # read current value
MULTI
SET balance (current + 10)     # set new value
EXEC                           # returns nil if balance changed
                               # returns [OK] if no conflict
```

```python
import redis

r = redis.Redis()

def transfer(from_key: str, to_key: str, amount: float) -> bool:
    with r.pipeline() as pipe:
        while True:
            try:
                pipe.watch(from_key, to_key)     # watch both keys

                from_bal = float(pipe.get(from_key) or 0)
                to_bal   = float(pipe.get(to_key) or 0)

                if from_bal < amount:
                    pipe.unwatch()
                    return False  # insufficient funds

                pipe.multi()                      # start transaction
                pipe.set(from_key, from_bal - amount)
                pipe.set(to_key,   to_bal   + amount)
                pipe.execute()                    # atomic execution
                return True

            except redis.WatchError:
                # Another client modified the keys — retry
                continue

transfer("account:alice", "account:bob", 50.0)
```

**When WATCH returns nil (transaction aborted):**
- The watched key was modified by another client between WATCH and EXEC
- Retry the entire operation from WATCH

---

## Lua Scripting

Lua scripts run on the Redis server atomically. They can read and write any key within the script.

```bash
# Run inline script
EVAL "return 'hello from lua'" 0

# Script with keys and arguments
EVAL "
  local current = redis.call('GET', KEYS[1])
  if not current then
    redis.call('SET', KEYS[1], ARGV[1])
    return 1
  end
  return 0
" 1 my_key default_value

# KEYS[1], KEYS[2]... = key arguments (cluster-routing aware)
# ARGV[1], ARGV[2]... = value arguments
```

### Redis-Callable Functions

```lua
-- redis.call: raises error on failure
local val = redis.call("GET", KEYS[1])

-- redis.pcall: returns error as a Lua table (no exception)
local ok, err = pcall(function()
    return redis.call("INCR", KEYS[1])
end)

-- redis.status_reply: return status
return redis.status_reply("OK")

-- redis.error_reply: return error
return redis.error_reply("ERR something went wrong")
```

### Practical Lua Scripts

#### Atomic Increment with Limit (Rate Limiter)

```lua
-- rate_limit.lua
-- KEYS[1]: rate limit key
-- ARGV[1]: max requests
-- ARGV[2]: window seconds
local current = tonumber(redis.call("INCR", KEYS[1]))
if current == 1 then
    redis.call("EXPIRE", KEYS[1], ARGV[2])
end
if current > tonumber(ARGV[1]) then
    return 0  -- rate limited
end
return 1  -- allowed
```

```python
rate_script = r.register_script(open("rate_limit.lua").read())
allowed = rate_script(keys=["rate:user-1"], args=[100, 60])
```

#### Get and Delete (Atomic)

```lua
-- get_and_delete.lua
local val = redis.call("GET", KEYS[1])
if val then
    redis.call("DEL", KEYS[1])
end
return val
```

#### Conditional Set

```lua
-- set_if_less.lua: set value only if current value is less
local current = tonumber(redis.call("GET", KEYS[1])) or 0
local new_val = tonumber(ARGV[1])
if new_val < current or current == 0 then
    redis.call("SET", KEYS[1], ARGV[1])
    return 1
end
return 0
```

---

## Script Caching (EVALSHA)

Large scripts can be cached on the server to avoid resending on every call:

```bash
# Load script and get SHA1
SCRIPT LOAD "return redis.call('GET', KEYS[1])"
# "2067d915024a3e1657c4169c84f809f8ec75b9a7"

# Execute by SHA1 (faster — no script transmission)
EVALSHA 2067d915024a3e1657c4169c84f809f8ec75b9a7 1 mykey

# Check if script is cached
SCRIPT EXISTS sha1_hash

# Flush all cached scripts
SCRIPT FLUSH
```

```python
import redis

r = redis.Redis()

# register_script caches and uses EVALSHA automatically
get_and_delete = r.register_script("""
local val = redis.call("GET", KEYS[1])
if val then redis.call("DEL", KEYS[1]) end
return val
""")

value = get_and_delete(keys=["otp:user-1"])
```

---

## Redis Functions (Redis 7.0+)

Functions are named, persistent Lua scripts loaded into the server:

```bash
# Load a function library
FUNCTION LOAD "#!lua name=mylib\n
redis.register_function('getset', function(keys, args)
    local old = redis.call('GET', keys[1])
    redis.call('SET', keys[1], args[1])
    return old
end)"

# Call the function
FCALL getset 1 mykey new-value

# List loaded functions
FUNCTION LIST

# Delete a library
FUNCTION DELETE mylib
```

---

## Transactions vs Lua: When to Use Each

| Feature | MULTI/EXEC | Lua Script |
|---|---|---|
| Atomicity | Yes (all-or-nothing) | Yes (single thread) |
| Read-then-write | Only with WATCH | Yes (native) |
| Conditional logic | No | Yes |
| Error handling | Limited | Full (pcall) |
| Performance | Batched round-trips | Single round-trip |
| Complexity | Low | Higher |
| Use case | Simple atomic batches | Complex atomic logic |

**Rule of thumb:**
- Simple batched operations → MULTI/EXEC
- Conditional logic, read-modify-write → Lua
- Reusable server-side logic → Redis Functions
