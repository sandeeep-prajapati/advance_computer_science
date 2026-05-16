# Chapter 8: Kafka Connect

## What is Kafka Connect?

Kafka Connect is a **framework for streaming data between Kafka and external systems** using reusable, configurable connectors. You declare what to connect — Kafka handles the rest.

```
External System ──► Source Connector ──► Kafka ──► Sink Connector ──► External System

Examples:
  PostgreSQL CDC ──► Debezium ──► Kafka ──► Elasticsearch Sink ──► Elasticsearch
  S3 Files       ──► S3 Source ──► Kafka ──► JDBC Sink           ──► MySQL
```

**Key properties:**
- Configuration-driven (JSON/REST API)
- Distributed mode: scales with multiple workers
- Exactly-once or at-least-once delivery
- Automatic offset management
- Transformations (SMTs) without custom code
- 200+ connectors in the Confluent Hub ecosystem

---

## Standalone vs Distributed Mode

| Feature | Standalone | Distributed |
|---|---|---|
| Workers | Single process | Multiple processes |
| Fault tolerance | None | Yes (tasks redistributed) |
| Scalability | No | Yes |
| Use case | Dev / testing | Production |
| Config | Local file | REST API |

---

## Start Connect (Distributed)

```properties
# connect-distributed.properties
bootstrap.servers=localhost:9092
group.id=connect-cluster

key.converter=org.apache.kafka.connect.json.JsonConverter
value.converter=org.apache.kafka.connect.json.JsonConverter
key.converter.schemas.enable=false
value.converter.schemas.enable=false

# For Avro with Schema Registry:
# key.converter=io.confluent.connect.avro.AvroConverter
# key.converter.schema.registry.url=http://localhost:8081
# value.converter=io.confluent.connect.avro.AvroConverter
# value.converter.schema.registry.url=http://localhost:8081

offset.storage.topic=connect-offsets
config.storage.topic=connect-configs
status.storage.topic=connect-status

offset.flush.interval.ms=10000
```

```bash
# Start Connect worker
connect-distributed connect-distributed.properties
```

---

## REST API

```bash
# List connectors
curl http://localhost:8083/connectors

# Create connector
curl -X POST http://localhost:8083/connectors \
  -H "Content-Type: application/json" \
  -d @connector-config.json

# Connector status
curl http://localhost:8083/connectors/my-connector/status

# Pause / resume
curl -X PUT http://localhost:8083/connectors/my-connector/pause
curl -X PUT http://localhost:8083/connectors/my-connector/resume

# Delete
curl -X DELETE http://localhost:8083/connectors/my-connector

# Available plugins
curl http://localhost:8083/connector-plugins
```

---

## Source Connectors

### Debezium PostgreSQL CDC

Capture every INSERT/UPDATE/DELETE from Postgres as Kafka events:

```json
{
  "name": "postgres-cdc",
  "config": {
    "connector.class":            "io.debezium.connector.postgresql.PostgresConnector",
    "database.hostname":          "postgres",
    "database.port":              "5432",
    "database.user":              "replicator",
    "database.password":          "secret",
    "database.dbname":            "shop",
    "database.server.name":       "shop",
    "table.include.list":         "public.orders,public.users",
    "plugin.name":                "pgoutput",
    "topic.prefix":               "cdc",
    "snapshot.mode":              "initial"
  }
}
```

Events published to: `cdc.public.orders`, `cdc.public.users`

### File Source (Simple)

```json
{
  "name": "file-source",
  "config": {
    "connector.class": "org.apache.kafka.connect.file.FileStreamSourceConnector",
    "file":            "/var/log/app.log",
    "topic":           "app-logs"
  }
}
```

---

## Sink Connectors

### Elasticsearch Sink

```json
{
  "name": "elasticsearch-sink",
  "config": {
    "connector.class":               "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
    "tasks.max":                     "3",
    "topics":                        "orders,events",
    "connection.url":                "http://elasticsearch:9200",
    "type.name":                     "_doc",
    "key.ignore":                    "false",
    "schema.ignore":                 "true",
    "behavior.on.null.values":       "DELETE",
    "batch.size":                    "2000",
    "max.retries":                   "5",
    "retry.backoff.ms":              "100"
  }
}
```

### JDBC Sink (Write to any SQL DB)

```json
{
  "name": "postgres-sink",
  "config": {
    "connector.class":        "io.confluent.connect.jdbc.JdbcSinkConnector",
    "tasks.max":              "3",
    "topics":                 "processed-orders",
    "connection.url":         "jdbc:postgresql://postgres:5432/warehouse",
    "connection.user":        "writer",
    "connection.password":    "secret",
    "auto.create":            "true",
    "auto.evolve":            "true",
    "insert.mode":            "upsert",
    "pk.mode":                "record_key",
    "pk.fields":              "order_id"
  }
}
```

### S3 Sink

```json
{
  "name": "s3-sink",
  "config": {
    "connector.class":         "io.confluent.connect.s3.S3SinkConnector",
    "tasks.max":               "4",
    "topics":                  "user-events",
    "s3.region":               "us-east-1",
    "s3.bucket.name":          "my-kafka-data",
    "s3.part.size":            "67108864",
    "storage.class":           "io.confluent.connect.s3.storage.S3Storage",
    "format.class":            "io.confluent.connect.s3.format.parquet.ParquetFormat",
    "flush.size":              "10000",
    "rotate.interval.ms":      "60000",
    "locale":                  "en_US",
    "timezone":                "UTC",
    "timestamp.extractor":     "RecordField",
    "timestamp.field":         "created_at",
    "partitioner.class":       "io.confluent.connect.storage.partitioner.TimeBasedPartitioner",
    "path.format":             "'year'=YYYY/'month'=MM/'day'=dd/'hour'=HH",
    "partition.duration.ms":   "3600000"
  }
}
```

---

## Single Message Transforms (SMTs)

SMTs modify messages in-flight without custom code:

```json
{
  "transforms":                   "extractKey,addTimestamp,dropField",

  "transforms.extractKey.type":   "org.apache.kafka.connect.transforms.ExtractField$Key",
  "transforms.extractKey.field":  "order_id",

  "transforms.addTimestamp.type": "org.apache.kafka.connect.transforms.InsertField$Value",
  "transforms.addTimestamp.timestamp.field": "kafka_timestamp",

  "transforms.dropField.type":    "org.apache.kafka.connect.transforms.ReplaceField$Value",
  "transforms.dropField.blacklist": "internal_field,debug_info"
}
```

### Common SMTs

| SMT | Purpose |
|---|---|
| `ReplaceField` | Rename, include, or exclude fields |
| `ExtractField` | Pull a field out to become the key or value |
| `InsertField` | Add fields (offset, partition, timestamp, static) |
| `MaskField` | Replace field value with null or a constant |
| `Filter` | Drop messages that match a condition |
| `TimestampConverter` | Convert timestamp formats |
| `RegexRouter` | Route to different topics by regex |
| `Flatten` | Flatten nested structs into dot-notation |

---

## Connector Monitoring

```bash
# Check task status
curl http://localhost:8083/connectors/postgres-cdc/status | jq .

# Check offset lag (source)
curl http://localhost:8083/connectors/postgres-cdc/offsets

# Restart a failed task
curl -X POST http://localhost:8083/connectors/postgres-cdc/tasks/0/restart

# Worker metrics via JMX / REST
curl http://localhost:8083/ | jq .version
```

Key JMX metrics:
- `kafka.connect:type=connector-task-metrics` → task state
- `kafka.connect:type=source-task-metrics` → poll rate, offset commit rate
- `kafka.connect:type=sink-task-metrics` → put rate, offset commit rate
