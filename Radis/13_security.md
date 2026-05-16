# Chapter 13: Security

## Security Layers

```
1. Network isolation      ← Redis only accessible from trusted hosts
2. Authentication (AUTH)  ← Require password or username+password
3. ACLs                   ← Per-user command and key restrictions
4. TLS                    ← Encrypt traffic in transit
5. Command restrictions   ← Disable dangerous commands
6. Protected mode         ← Default: reject external connections
```

---

## Authentication

### Simple Password (Redis < 6.0 style)

```properties
# redis.conf
requirepass yourStrongPassword123!
```

```bash
# Authenticate
AUTH yourStrongPassword123!
```

### ACL Users (Redis 6.0+, recommended)

```properties
# redis.conf — define users
# Disable default user (no password)
user default off nopass nocommands nokeys

# Create admin user
user admin on >adminPassword ~* &* +@all

# Create app user (read/write to app:* keys only)
user app-user on >appPassword ~app:* &* +@read +@write +@string +@hash +EXPIRE +DEL

# Create read-only user
user reader on >readPassword ~* &* +@read
```

---

## ACL System (Redis 6.0+)

### ACL Rule Components

```
user <username> [on|off]           ← enabled or disabled
                [>password]        ← set password (> prefix)
                [<password]        ← remove password (< prefix)
                [nopass]           ← no password required
                [~pattern]         ← allowed key patterns
                [%R~pattern]       ← read-only key pattern
                [%W~pattern]       ← write-only key pattern
                [&pattern]         ← allowed pub/sub channel patterns
                [+command]         ← allow command
                [-command]         ← deny command
                [+@category]       ← allow all commands in category
                [-@category]       ← deny all commands in category
                [allkeys]          ← all keys allowed (~*)
                [allchannels]      ← all channels allowed (&*)
                [allcommands]      ← all commands allowed (+@all)
                [nocommands]       ← no commands allowed (-@all)
                [resetkeys]        ← remove all key patterns
                [reset]            ← full reset to defaults
```

### ACL Commands

```bash
# List all users
ACL LIST

# Show current user
ACL WHOAMI

# Get specific user info
ACL GETUSER app-user

# Create/modify user
ACL SETUSER cache-user on >cachePass ~cache:* +GET +SET +DEL +EXPIRE +TTL

# Delete user
ACL DELUSER old-user

# Test permissions
ACL DRYRUN cache-user GET cache:product:1
ACL DRYRUN cache-user SET sensitive:data "value"   # returns error

# Log recent ACL failures
ACL LOG
ACL LOG RESET
ACL LOG COUNT
```

### ACL Categories

```bash
# List categories
ACL CAT

# List commands in a category
ACL CAT read
ACL CAT write
ACL CAT string
ACL CAT hash
ACL CAT pubsub
ACL CAT admin
ACL CAT dangerous
```

---

## TLS Configuration

### Generate Certificates

```bash
# Create CA and server certificates
mkdir -p /etc/redis/tls

# CA
openssl genrsa -out /etc/redis/tls/ca.key 4096
openssl req -x509 -new -nodes -key /etc/redis/tls/ca.key \
  -sha256 -days 1024 -out /etc/redis/tls/ca.crt \
  -subj "/CN=Redis-CA/O=MyOrg"

# Server cert
openssl genrsa -out /etc/redis/tls/redis.key 2048
openssl req -new -key /etc/redis/tls/redis.key \
  -out /etc/redis/tls/redis.csr \
  -subj "/CN=redis-server/O=MyOrg"
openssl x509 -req -in /etc/redis/tls/redis.csr \
  -CA /etc/redis/tls/ca.crt -CAkey /etc/redis/tls/ca.key \
  -CAcreateserial -out /etc/redis/tls/redis.crt -days 365 -sha256
```

### Redis TLS Config

```properties
# redis.conf
tls-port 6380
port 0                              # disable non-TLS port

tls-cert-file /etc/redis/tls/redis.crt
tls-key-file  /etc/redis/tls/redis.key
tls-ca-cert-file /etc/redis/tls/ca.crt
tls-auth-clients yes               # require client certificates
tls-protocols "TLSv1.3"
tls-ciphers DEFAULT:!MEDIUM
tls-prefer-server-ciphers yes

# Replication over TLS
tls-replication yes
```

### Client TLS (Python)

```python
import redis
import ssl

r = redis.Redis(
    host     = "redis-host",
    port     = 6380,
    ssl      = True,
    ssl_certfile  = "/etc/redis/tls/client.crt",
    ssl_keyfile   = "/etc/redis/tls/client.key",
    ssl_ca_certs  = "/etc/redis/tls/ca.crt",
    ssl_check_hostname = True,
)
r.ping()
```

---

## Disabling Dangerous Commands

```properties
# redis.conf — rename or disable dangerous commands
rename-command FLUSHALL ""          # disable completely
rename-command FLUSHDB  ""
rename-command CONFIG   "CONFIG-ADMIN-ONLY-abc123"
rename-command DEBUG    ""
rename-command KEYS     "SCAN-INSTEAD"
rename-command SHUTDOWN "SHUTDOWN-abc123"
```

---

## Network Security

```properties
# Bind to specific interfaces only
bind 127.0.0.1 10.0.0.1    # only localhost and internal network

# Protected mode (default: on)
# Rejects connections from external IPs if no auth is configured
protected-mode yes

# Max clients
maxclients 1000

# Connection timeout
timeout 300                  # close idle connections after 300s
tcp-keepalive 60
```

---

## Security Best Practices

```
1. Never expose Redis directly to the internet
   → Use firewall rules, VPC, security groups

2. Always use requirepass or ACL users
   → Rotate passwords periodically

3. Use TLS for all client connections in production
   → Especially across networks

4. Apply ACLs — least privilege per service
   → App user: only reads/writes to app:* keys
   → Monitoring user: INFO only, no writes

5. Rename or disable dangerous commands
   → FLUSHALL, FLUSHDB, DEBUG, CONFIG

6. Run Redis as a dedicated OS user (not root)
   → useradd -r -s /sbin/nologin redis

7. Use Redis Sentinel or Cluster behind a load balancer
   → Never expose individual node IPs to application code

8. Enable AOF persistence
   → Prevents data loss on Redis restart after attack

9. Monitor the ACL log for unauthorized access attempts
   → ACL LOG every 5 minutes

10. Keep Redis updated
    → Security patches released regularly
```
