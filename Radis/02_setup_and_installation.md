# Chapter 2: Setup & Installation

## Quick Start with Docker

```bash
# Single Redis instance
docker run -d \
  --name redis \
  -p 6379:6379 \
  redis:7.2

# With persistence
docker run -d \
  --name redis \
  -p 6379:6379 \
  -v redis-data:/data \
  redis:7.2 redis-server --appendonly yes

# Connect
docker exec -it redis redis-cli ping
# PONG
```

---

## Docker Compose (Redis + RedisInsight)

```yaml
# docker-compose.yml
version: "3.8"
services:
  redis:
    image: redis:7.2
    container_name: redis
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
      - ./redis.conf:/usr/local/etc/redis/redis.conf
    command: redis-server /usr/local/etc/redis/redis.conf
    restart: unless-stopped

  redis-insight:
    image: redis/redisinsight:latest
    container_name: redis-insight
    ports:
      - "5540:5540"
    depends_on:
      - redis

volumes:
  redis-data:
```

```bash
docker compose up -d
# Redis UI: http://localhost:5540
```

---

## Redis Configuration File (redis.conf)

```properties
# Network
bind 127.0.0.1
port 6379
protected-mode yes

# Authentication
requirepass yourStrongPassword123

# Persistence
save 900 1          # Save if 1 key changed in 900 seconds
save 300 10         # Save if 10 keys changed in 300 seconds
save 60 10000       # Save if 10000 keys changed in 60 seconds
dir /data
dbfilename dump.rdb

# AOF
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec

# Memory
maxmemory 256mb
maxmemory-policy allkeys-lru

# Logging
loglevel notice
logfile "/var/log/redis/redis-server.log"

# Connection limits
maxclients 10000
timeout 300

# Slow log
slowlog-log-slower-than 10000   # microseconds
slowlog-max-len 128
```

---

## Local Installation

### Ubuntu / Debian
```bash
# Official Redis repo
curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" \
  | sudo tee /etc/apt/sources.list.d/redis.list
sudo apt-get update
sudo apt-get install -y redis

sudo systemctl enable redis-server
sudo systemctl start redis-server
redis-cli ping
```

### macOS
```bash
brew install redis
brew services start redis
redis-cli ping
```

### Build from Source
```bash
wget https://download.redis.io/redis-stable.tar.gz
tar xvzf redis-stable.tar.gz
cd redis-stable
make
sudo make install
redis-server redis.conf
```

---

## redis-cli Basics

```bash
# Connect
redis-cli                           # localhost:6379
redis-cli -h myhost -p 6380         # remote
redis-cli -a password               # with auth
redis-cli -n 2                      # database 2

# Basic commands
redis-cli ping
redis-cli info server
redis-cli dbsize                    # number of keys

# Interactive mode
redis-cli
127.0.0.1:6379> SET name "Alice"
OK
127.0.0.1:6379> GET name
"Alice"
127.0.0.1:6379> DEL name
(integer) 1

# Execute single command
redis-cli SET counter 0
redis-cli INCR counter
redis-cli GET counter

# Pipeline from file
redis-cli --pipe < commands.txt

# Monitor (watch all commands in real time)
redis-cli MONITOR

# Benchmark
redis-benchmark -n 100000 -c 50 -q

# Scan all keys (safe, no KEYS * in production)
redis-cli --scan --pattern "user:*"
```

---

## Redis Insight (GUI)

Redis Insight is the official web-based GUI for Redis.

**Features:**
- Browse and edit keys visually
- Run CLI commands
- Analyze memory usage per key
- View slow log
- Profiler (live command stream)
- Pub/Sub explorer

```bash
# Docker
docker run -d -p 5540:5540 redis/redisinsight:latest

# Access: http://localhost:5540
# Add connection: host=localhost, port=6379
```

---

## Verify Installation: First Commands

```bash
redis-cli

# Test connectivity
PING
# PONG

# Store and retrieve
SET greeting "Hello, Redis!"
GET greeting
# "Hello, Redis!"

# With expiry (TTL in seconds)
SETEX session:abc 3600 "user-data"
TTL session:abc
# 3600

# Atomic counter
SET visits 0
INCR visits
INCR visits
INCR visits
GET visits
# "3"

# List of recent items
LPUSH recent:items "item3"
LPUSH recent:items "item2"
LPUSH recent:items "item1"
LRANGE recent:items 0 -1
# 1) "item1"
# 2) "item2"
# 3) "item3"

# Check server info
INFO server
INFO memory
INFO stats
```

---

## Multiple Databases

Redis has 16 logical databases (0–15) by default, all in the same instance:

```bash
# Switch databases
SELECT 0    # default
SELECT 1
SELECT 15

# Copy key across databases
COPY mykey 1    # copy to DB 1

# Move key to another database
MOVE mykey 2

# Flush current database (destructive!)
FLUSHDB

# Flush all databases (very destructive!)
FLUSHALL
```

**Note:** In Redis Cluster mode, only database 0 is available.
