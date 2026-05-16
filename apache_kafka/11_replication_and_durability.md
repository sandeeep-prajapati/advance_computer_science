# Chapter 11: Replication & Durability

## The Replication Protocol

Kafka replication is **leader-based**: all reads and writes go through the partition leader. Followers pull data from the leader in parallel.

```
Write path:
  Producer ──► Broker 1 (Leader, P0)
                    │
                    ├──► Broker 2 (Follower, P0) ← replicate
                    └──► Broker 3 (Follower, P0) ← replicate
                    │
                    ACK to producer (when ISR confirms)

Read path:
  Consumer ──► Broker 1 (Leader, P0)
                    │
                    Return records
                    (followers also serve reads in Kafka 2.4+ with replica.selector.class)
```

---

## Replication Factor

```bash
# Create topic with RF=3
kafka-topics --bootstrap-server localhost:9092 \
  --create --topic payments \
  --partitions 6 \
  --replication-factor 3

# Describe to see replicas
kafka-topics --bootstrap-server localhost:9092 --describe --topic payments

# Output shows Replicas and Isr (In-Sync Replicas) per partition:
# Partition: 0  Leader: 2  Replicas: 2,1,3  Isr: 2,1,3
# Partition: 1  Leader: 1  Replicas: 1,3,2  Isr: 1,3,2
```

---

## Producer Acks and Durability

The `acks` setting controls the durability/latency tradeoff:

```
acks=0 (fire and forget)
───────────────────────
Producer ──(write)──► Broker
          ←(no ack)─
Risk: Network or broker failure = lost message
Use: Metrics, non-critical logging

acks=1 (leader only)
────────────────────
Producer ──(write)──► Leader ──(ack)──► Producer
                        │
                        Followers might not have it yet
Risk: Leader crashes before replication = lost message
Use: Moderate throughput, acceptable loss

acks=all (full ISR)
────────────────────
Producer ──(write)──► Leader ──(replicate)──► All ISR followers
                        │◄──────────────────────
                        ACK after all ISR confirm
Risk: None (as long as min.insync.replicas respected)
Use: Production, financial data, audit logs
```

---

## min.insync.replicas

This setting determines the minimum number of replicas that must acknowledge a write:

```
Configuration:
  replication.factor=3
  min.insync.replicas=2
  acks=all

Scenario 1: ISR = [B1, B2, B3] → writes succeed
Scenario 2: ISR = [B1, B2]     → writes succeed (ISR ≥ min.insync)
Scenario 3: ISR = [B1]         → writes FAIL
  → Producer gets: NotEnoughReplicasException
  → Cluster is write-unavailable for this partition

Why this is correct:
  If ISR=[B1] and we allowed writes, data would be on only 1 broker.
  If B1 then crashes → DATA LOSS despite acks=all.
  Better to reject writes than accept unacknowledged data.
```

**Recommended production settings:**
```properties
# server.properties (applied cluster-wide as defaults)
default.replication.factor=3
min.insync.replicas=2

# Per-topic override
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics --entity-name payments \
  --alter --add-config min.insync.replicas=2
```

---

## Follower Replication

Followers continuously fetch batches from the leader:

```properties
# How long a follower can lag before removed from ISR
replica.lag.time.max.ms=30000       # default: 30 seconds

# Max bytes per fetch request
replica.fetch.max.bytes=1048576     # 1 MB

# How fast followers replicate
replica.fetch.min.bytes=1
replica.fetch.wait.max.ms=500
```

---

## Unclean Leader Election

When all ISR replicas fail simultaneously and only out-of-sync replicas remain:

```properties
# Allow out-of-sync replica to become leader (risk: data loss)
unclean.leader.election.enable=false   # default: false for most topics

# Use true only for high-availability over data durability
# e.g., real-time metrics where losing some data is acceptable
```

---

## Log Segments and Durability

Kafka stores data in segment files:

```
/kafka-logs/orders-0/
  00000000000000000000.log         ← active segment (writable)
  00000000000000012345.log         ← closed segment
  00000000000000012345.index       ← offset index
  00000000000000012345.timeindex   ← timestamp index
```

**Log flushing (fsync) settings:**

```properties
# Flush data to disk after N messages (default: Long.MAX_VALUE = rely on OS)
log.flush.interval.messages=Long.MAX_VALUE

# Flush at least every N ms (default: Long.MAX_VALUE)
log.flush.interval.ms=Long.MAX_VALUE

# For strongest durability (but high latency):
log.flush.interval.messages=1      # fsync every message — VERY slow
```

**Recommendation:** Leave at defaults. Rely on replication (ISR) for durability, not fsync. OS page cache provides good protection.

---

## Replication Factor Change (Partition Reassignment)

You cannot change replication factor via `kafka-topics --alter`. Use `kafka-reassign-partitions`:

```bash
# 1. Create reassignment JSON
cat > reassignment.json << 'EOF'
{
  "version": 1,
  "partitions": [
    {"topic": "orders", "partition": 0, "replicas": [1,2,3], "log_dirs": ["any","any","any"]},
    {"topic": "orders", "partition": 1, "replicas": [2,3,1], "log_dirs": ["any","any","any"]}
  ]
}
EOF

# 2. Execute reassignment
kafka-reassign-partitions --bootstrap-server localhost:9092 \
  --reassignment-json-file reassignment.json \
  --execute \
  --throttle 50000000   # 50 MB/s throttle (avoid saturating network)

# 3. Verify progress
kafka-reassign-partitions --bootstrap-server localhost:9092 \
  --reassignment-json-file reassignment.json \
  --verify

# 4. Remove throttle after completion
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type brokers --entity-default \
  --alter --delete-config leader.replication.throttled.rate
```

---

## Durability Checklist

| Setting | Production Value | Why |
|---|---|---|
| `replication.factor` | 3 | Survive 2 simultaneous broker failures |
| `min.insync.replicas` | 2 | No data accepted with only 1 copy |
| `acks` (producer) | `all` | All ISR must confirm |
| `unclean.leader.election.enable` | `false` | No data loss on leader election |
| `enable.idempotence` (producer) | `true` | No duplicates on retry |
| `replica.lag.time.max.ms` | `30000` | Remove lagging replicas from ISR promptly |

---

## Monitoring Replication Health

```bash
# Under-replicated partitions — ISR smaller than replication factor
kafka-topics --bootstrap-server localhost:9092 \
  --describe --under-replicated-partitions

# Non-preferred leaders (performance degraded)
kafka-topics --bootstrap-server localhost:9092 \
  --describe

# JMX metrics
# kafka.server:type=ReplicaManager,name=UnderReplicatedPartitions → should be 0
# kafka.server:type=ReplicaManager,name=IsrShrinksPerSec → should be 0 normally
# kafka.server:type=ReplicaManager,name=IsrExpandsPerSec → normal during recoveries
```
