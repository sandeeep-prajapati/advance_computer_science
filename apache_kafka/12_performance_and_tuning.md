# Chapter 12: Performance & Tuning

## Understanding Kafka Performance

Kafka achieves high throughput through:

1. **Sequential disk I/O** — append-only writes are fast (disk can do 500 MB/s sequential vs ~random IOPS)
2. **OS page cache** — Linux caches recently read disk pages in RAM
3. **Zero-copy** (`sendfile` syscall) — kernel sends data directly from page cache to network socket
4. **Batching** — producers accumulate messages, send in one network call
5. **Compression** — reduce network and disk bandwidth

```
Without zero-copy:
  Disk → Kernel page cache → User space → Kernel socket buffer → Network

With zero-copy (sendfile):
  Disk → Kernel page cache ──────────────────────────────► Network
  (2 copies vs 4 copies)
```

---

## Producer Tuning

### Maximize Throughput

```python
producer = Producer({
    "bootstrap.servers": "broker1:9092,broker2:9092,broker3:9092",

    # Batching — let messages accumulate
    "linger.ms":        20,         # wait 20ms to fill a batch
    "batch.size":       65536,      # 64KB batch size
    "buffer.memory":    67108864,   # 64MB total send buffer

    # Compression reduces network bandwidth significantly
    "compression.type": "lz4",      # lz4: fast, good ratio
                                    # zstd: best ratio, slower CPU
                                    # snappy: balanced

    # Throughput-optimized acks
    "acks":             "1",        # leader only (accept some risk)
})
```

### Maximize Reliability

```python
producer = Producer({
    "bootstrap.servers":         "localhost:9092",
    "acks":                      "all",
    "enable.idempotence":        True,
    "retries":                   10,
    "retry.backoff.ms":          100,
    "delivery.timeout.ms":       120000,
    "max.in.flight.requests.per.connection": 5,
    "linger.ms":                 5,      # still batch a bit
    "compression.type":         "snappy",
})
```

---

## Consumer Tuning

```python
consumer = Consumer({
    "bootstrap.servers":      "localhost:9092",
    "group.id":               "fast-consumer",
    "max.poll.records":       2000,     # fetch more records per poll
    "fetch.min.bytes":        1048576,  # 1MB — wait for enough data
    "fetch.max.wait.ms":      500,      # max wait for fetch.min.bytes
    "fetch.max.bytes":        52428800, # 50MB max per fetch
    "max.partition.fetch.bytes": 1048576,

    # Avoid rebalancing from slow consumers
    "max.poll.interval.ms":   600000,   # 10 min max processing time
    "session.timeout.ms":     30000,
    "heartbeat.interval.ms":  3000,
})
```

---

## Broker-Level Tuning

### OS Configuration

```bash
# Increase file descriptors (Kafka needs many open files)
echo "* soft nofile 100000" >> /etc/security/limits.conf
echo "* hard nofile 100000" >> /etc/security/limits.conf

# Increase network buffers
sysctl -w net.core.rmem_max=16777216
sysctl -w net.core.wmem_max=16777216
sysctl -w net.ipv4.tcp_rmem="4096 87380 16777216"
sysctl -w net.ipv4.tcp_wmem="4096 65536 16777216"

# Disable swapping for Kafka heap
sysctl -w vm.swappiness=1

# Dirty page flushing (don't flush too aggressively)
sysctl -w vm.dirty_ratio=80
sysctl -w vm.dirty_background_ratio=5

# Use deadline or noop I/O scheduler (not cfq)
echo deadline > /sys/block/sda/queue/scheduler
```

### JVM Tuning (kafka-server-start.sh)

```bash
export KAFKA_HEAP_OPTS="-Xmx6g -Xms6g"   # 6GB heap, pre-allocate
export KAFKA_JVM_PERFORMANCE_OPTS="
  -server
  -XX:+UseG1GC
  -XX:MaxGCPauseMillis=20
  -XX:InitiatingHeapOccupancyPercent=35
  -XX:+ExplicitGCInvokesConcurrent
  -Djava.awt.headless=true"
```

**Heap size guidance:**
- Don't give Kafka more than 6–8 GB heap
- Kafka relies heavily on OS page cache (RAM not in heap)
- More RAM → more page cache → better read performance

### Storage

```bash
# Use XFS or ext4 (XFS preferred for Kafka)
mkfs.xfs -f /dev/sdb
mount -t xfs -o noatime,nodiratime /dev/sdb /var/kafka/data

# RAID: use JBOD (Just a Bunch of Disks) not RAID5/6
# Multiple log.dirs = Kafka stripes across disks automatically
log.dirs=/kafka/disk1,/kafka/disk2,/kafka/disk3
```

---

## Compression

| Codec | CPU | Ratio | Latency | Best For |
|---|---|---|---|---|
| `none` | None | 1.0× | Lowest | Already-compressed data |
| `gzip` | High | 3-5× | High | Archival, max compression |
| `snappy` | Low | 2-3× | Low | General purpose |
| `lz4` | Low | 2-3× | Very low | High throughput |
| `zstd` | Medium | 3-6× | Low | Best ratio/speed tradeoff |

**Producer-side vs broker-side compression:**
- `compression.type=producer` (default) — broker stores as-is
- Set on topic to recompress: risky, usually avoid

---

## Partitions and Parallelism

```
Consumer throughput ≈ consumer_instances × per_consumer_MB_s

To process 100 MB/s:
  If each consumer handles 10 MB/s → need 10 consumers → need ≥10 partitions

Increasing partitions:
  - Increases producer parallelism
  - Increases consumer parallelism
  - Increases broker overhead (file handles, leader election time)

Don't over-partition:
  - 1,000 partitions × 3 replicas = 3,000 partition replicas per cluster
  - Kafka recommends max ~4,000 partitions per broker
```

---

## Benchmarking

### Producer Perf Test

```bash
kafka-producer-perf-test \
  --topic benchmark \
  --num-records 10000000 \
  --record-size 1000 \
  --throughput -1 \
  --producer-props \
    bootstrap.servers=localhost:9092 \
    acks=all \
    linger.ms=5 \
    batch.size=65536 \
    compression.type=lz4
```

### Consumer Perf Test

```bash
kafka-consumer-perf-test \
  --bootstrap-server localhost:9092 \
  --topic benchmark \
  --messages 10000000 \
  --group perf-test-group \
  --threads 4
```

---

## Performance Anti-Patterns

| Anti-Pattern | Problem | Fix |
|---|---|---|
| `linger.ms=0` + small messages | Tiny batches, high overhead | Set `linger.ms=5-20` |
| No compression | Wastes network + disk | Enable `lz4` or `snappy` |
| Too many small partitions | Leader election overhead | Consolidate topics |
| `auto.offset.reset=latest` during replay | Miss events | Use `earliest` for replay |
| Long `max.poll.interval.ms` | Slow rebalance detection | Tune to actual processing time |
| Heap > 8 GB | GC pauses | Cap heap, add RAM for page cache |
| Sync commit per message | Low throughput | Batch commits or async |
| RAID 5/6 | Write amplification | Use JBOD |
