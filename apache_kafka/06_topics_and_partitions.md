# Chapter 6: Topics & Partitions

## Topic Configuration

Topics have two kinds of settings:
- **Static** (set at creation, require recreation to change): `replication.factor`, `num.partitions`
- **Dynamic** (alterable at runtime): retention, compression, cleanup policy

### Create with Custom Config

```bash
kafka-topics --bootstrap-server localhost:9092 \
  --create \
  --topic user-events \
  --partitions 12 \
  --replication-factor 3 \
  --config retention.ms=604800000 \
  --config compression.type=lz4 \
  --config cleanup.policy=delete
```

### Alter Topic Config

```bash
# Change retention to 30 days
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics \
  --entity-name user-events \
  --alter \
  --add-config retention.ms=2592000000

# Describe current config
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics \
  --entity-name user-events \
  --describe
```

---

## Key Topic Configurations

| Config | Default | Description |
|---|---|---|
| `retention.ms` | `604800000` (7 days) | How long to keep messages |
| `retention.bytes` | `-1` (unlimited) | Max size per partition |
| `cleanup.policy` | `delete` | `delete` or `compact` or `delete,compact` |
| `compression.type` | `producer` | `none`, `gzip`, `snappy`, `lz4`, `zstd` |
| `replication.factor` | Broker default | Copies of each partition |
| `min.insync.replicas` | `1` | Min ISR for producer `acks=all` |
| `segment.ms` | `604800000` (7 days) | When to roll a new segment |
| `segment.bytes` | `1073741824` (1 GB) | Max segment file size |
| `max.message.bytes` | `1048588` (~1 MB) | Max message size |
| `delete.retention.ms` | `86400000` (1 day) | How long to keep tombstones (compact) |

---

## Retention Policies

### Time-Based Retention (delete)

```bash
# Keep messages for 3 days
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics --entity-name logs \
  --alter --add-config "retention.ms=259200000"
```

### Size-Based Retention

```bash
# Keep at most 10 GB per partition
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics --entity-name logs \
  --alter --add-config "retention.bytes=10737418240"
```

### Combined (Whichever Triggers First)

```bash
--add-config "retention.ms=86400000,retention.bytes=5368709120"
```

---

## Log Compaction

Log compaction keeps only the **latest value per key**. Perfect for changelog/state topics.

```
Before compaction:
  [key=user1: v1][key=user2: v1][key=user1: v2][key=user1: v3]

After compaction:
  [key=user2: v1][key=user1: v3]

Tombstone: key=user1, value=null → delete key from compacted log
```

```bash
# Enable compaction
kafka-configs --bootstrap-server localhost:9092 \
  --entity-type topics --entity-name user-profiles \
  --alter --add-config "cleanup.policy=compact"

# Compact + delete (compaction on recent data, delete old)
--add-config "cleanup.policy=compact,delete,retention.ms=86400000"
```

**Compaction settings:**

| Config | Default | Description |
|---|---|---|
| `min.cleanable.dirty.ratio` | `0.5` | Compact when dirty/total ratio > threshold |
| `min.compaction.lag.ms` | `0` | Don't compact messages newer than this |
| `max.compaction.lag.ms` | `Long.MAX_VALUE` | Force compaction after this time |
| `delete.retention.ms` | `86400000` | Keep tombstones this long |

---

## Partition Count: How to Choose

```
Rule of thumb:
  Partitions = max(throughput_MB_s / 10, desired_consumer_parallelism)

Example:
  Target throughput: 100 MB/s
  Max consumer instances: 20
  → partitions = max(10, 20) = 20 partitions
```

**More partitions = more parallelism but:**
- More file handles on brokers
- Longer leader election time on failure
- More memory in producer and consumer

```bash
# Increase partitions (can only go up, never down)
kafka-topics --bootstrap-server localhost:9092 \
  --alter --topic orders --partitions 24
```

---

## Replication Factor

```
replication.factor=1: No redundancy. Single point of failure.
replication.factor=2: One copy. Can survive 1 broker failure.
replication.factor=3: Two copies. Can survive 2 broker failures.
                       Standard for production.
```

**Relationship with min.insync.replicas:**

```
replication.factor=3, min.insync.replicas=2, acks=all:
  - 3 replicas, need 2 to acknowledge each write
  - Can survive 1 broker being down
  - If 2 brokers fail → write unavailable (but no data loss)

replication.factor=3, min.insync.replicas=1, acks=all:
  - Only 1 replica needed → leader only
  - Same as acks=1 effectively
  - Don't do this
```

---

## Partition Assignment: Leaders and Followers

```bash
# View partition leadership
kafka-topics --bootstrap-server localhost:9092 \
  --describe --topic orders

# Output:
# Topic: orders   PartitionCount: 6  ReplicationFactor: 3
# Topic: orders   Partition: 0  Leader: 1  Replicas: 1,2,3  Isr: 1,2,3
# Topic: orders   Partition: 1  Leader: 2  Replicas: 2,3,1  Isr: 2,3,1
# Topic: orders   Partition: 2  Leader: 3  Replicas: 3,1,2  Isr: 3,1,2
```

### Rebalance Preferred Leaders

After a broker restart, it may not be the preferred leader:

```bash
# Reassign partitions back to preferred leaders
kafka-leader-election --bootstrap-server localhost:9092 \
  --election-type preferred \
  --all-topic-partitions
```

---

## Topic Best Practices

```
1. Naming conventions:
   <domain>.<entity>.<event-type>
   Example: payments.orders.created
            analytics.pageviews.raw

2. Partition count:
   Start with 6 for moderate-traffic topics
   24+ for high-throughput topics
   Use multiples of consumer instance count

3. Replication:
   Always 3 in production
   min.insync.replicas = replication.factor - 1

4. Retention:
   Event topics: 7 days (replay window)
   Log aggregation: 3 days or size-limited
   Compaction topics: forever (with tombstones)

5. Message size:
   Target < 1MB per message
   For large payloads: store in S3, send reference in Kafka

6. Schema:
   Always use Schema Registry for value schema
   Key should be a simple string or int (user_id, order_id)
```

---

## Monitoring Topic Health

```bash
# Under-replicated partitions (ISR < replication factor = problem!)
kafka-topics --bootstrap-server localhost:9092 \
  --describe --under-replicated-partitions

# Unavailable partitions (leader = none = outage!)
kafka-topics --bootstrap-server localhost:9092 \
  --describe --unavailable-partitions

# Topics with no leader
kafka-topics --bootstrap-server localhost:9092 \
  --describe --under-min-isr-partitions
```
