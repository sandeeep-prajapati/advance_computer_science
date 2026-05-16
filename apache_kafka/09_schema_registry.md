# Chapter 9: Schema Registry

## Why Schema Registry?

Without schema management, Kafka is just raw bytes. Problems arise when:
- A producer changes the message format, breaking consumers
- Fields are renamed or removed without coordination
- Types change (string → integer)
- Multiple teams produce to the same topic with different schemas

**Schema Registry solves this:**
1. Centralized schema storage with versioning
2. Compatibility enforcement (BACKWARD, FORWARD, FULL)
3. Automatic schema evolution with safety guarantees
4. Schema IDs embedded in messages — consumers always know the schema

---

## How It Works

```
Producer                  Schema Registry         Consumer
   │                            │                    │
   │── register schema ────────►│                    │
   │◄─ schema_id=5 ────────────┤                    │
   │                            │                    │
   │── [magic byte][schema_id=5][serialized payload] ──►│
   │                            │                    │
   │                            │◄─ fetch schema 5 ──┤
   │                            │─ return schema ───►│
   │                            │                    │
   │                            │                 deserialize
```

Every Confluent-format message starts with:
- 1 magic byte (`0x00`)
- 4-byte schema ID
- serialized payload (Avro, Protobuf, or JSON Schema)

---

## Avro Schema Example

### Define Schema

```json
// order.avsc
{
  "type": "record",
  "name": "Order",
  "namespace": "com.example.shop",
  "fields": [
    {"name": "order_id",   "type": "string"},
    {"name": "user_id",    "type": "string"},
    {"name": "amount",     "type": "double"},
    {"name": "currency",   "type": {"type": "enum", "name": "Currency",
                                   "symbols": ["USD", "EUR", "GBP"]}},
    {"name": "created_at", "type": "long",   "logicalType": "timestamp-millis"},
    {"name": "metadata",   "type": ["null", {"type": "map", "values": "string"}],
                           "default": null}
  ]
}
```

### Producer with Avro

```python
from confluent_kafka import Producer
from confluent_kafka.schema_registry import SchemaRegistryClient
from confluent_kafka.schema_registry.avro import AvroSerializer
from confluent_kafka.serialization import SerializationContext, MessageField

schema_registry = SchemaRegistryClient({"url": "http://localhost:8081"})

order_schema = open("order.avsc").read()

avro_serializer = AvroSerializer(schema_registry, order_schema)

producer = Producer({"bootstrap.servers": "localhost:9092"})

order = {
    "order_id":   "ORD-123",
    "user_id":    "user-456",
    "amount":     99.99,
    "currency":   "USD",
    "created_at": 1704067200000,
    "metadata":   None,
}

producer.produce(
    topic = "orders",
    key   = order["order_id"],
    value = avro_serializer(order, SerializationContext("orders", MessageField.VALUE)),
)
producer.flush()
```

### Consumer with Avro

```python
from confluent_kafka import Consumer
from confluent_kafka.schema_registry import SchemaRegistryClient
from confluent_kafka.schema_registry.avro import AvroDeserializer
from confluent_kafka.serialization import SerializationContext, MessageField

schema_registry = SchemaRegistryClient({"url": "http://localhost:8081"})
avro_deserializer = AvroDeserializer(schema_registry)

consumer = Consumer({
    "bootstrap.servers": "localhost:9092",
    "group.id":          "order-consumer",
    "auto.offset.reset": "earliest",
})
consumer.subscribe(["orders"])

while True:
    msg = consumer.poll(1.0)
    if msg is None or msg.error(): continue

    order = avro_deserializer(msg.value(),
                              SerializationContext(msg.topic(), MessageField.VALUE))
    print(f"Order: {order}")
```

---

## Protobuf Schema Example

```protobuf
// order.proto
syntax = "proto3";
package com.example.shop;

message Order {
  string order_id   = 1;
  string user_id    = 2;
  double amount     = 3;
  string currency   = 4;
  int64  created_at = 5;
  map<string, string> metadata = 6;
}
```

```python
from confluent_kafka.schema_registry.protobuf import ProtobufSerializer, ProtobufDeserializer
from order_pb2 import Order  # generated from proto file

protobuf_serializer = ProtobufSerializer(Order, schema_registry)

order = Order(order_id="ORD-123", user_id="user-456", amount=99.99, currency="USD")
producer.produce(
    "orders",
    value = protobuf_serializer(order, SerializationContext("orders", MessageField.VALUE))
)
```

---

## JSON Schema

```python
json_schema_str = """
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title":   "Order",
  "type":    "object",
  "properties": {
    "order_id": {"type": "string"},
    "amount":   {"type": "number", "minimum": 0}
  },
  "required": ["order_id", "amount"]
}
"""

from confluent_kafka.schema_registry.json_schema import JSONSerializer
json_serializer = JSONSerializer(json_schema_str, schema_registry)
```

---

## Compatibility Modes

| Mode | What It Allows |
|---|---|
| `BACKWARD` | New schema can read data written with old schema (consumers upgrade first) |
| `BACKWARD_TRANSITIVE` | New schema can read all previous versions |
| `FORWARD` | Old schema can read data written with new schema (producers upgrade first) |
| `FORWARD_TRANSITIVE` | Any old schema can read data written with new schema |
| `FULL` | Both BACKWARD and FORWARD |
| `FULL_TRANSITIVE` | Both, transitively |
| `NONE` | No compatibility checks |

**Rule of thumb: use `FULL_TRANSITIVE` for production.**

### Safe Schema Changes (FULL compatibility)

| Change | Safe? |
|---|---|
| Add optional field with default | Yes |
| Remove optional field | Yes |
| Change field type | No |
| Remove required field | No |
| Rename field | No (use aliases) |
| Add required field without default | No |

---

## Schema Registry REST API

```bash
# Register a schema
curl -X POST http://localhost:8081/subjects/orders-value/versions \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  -d '{"schema": "{\"type\":\"record\",\"name\":\"Order\",...}"}'

# List subjects
curl http://localhost:8081/subjects

# Get all versions
curl http://localhost:8081/subjects/orders-value/versions

# Get specific version
curl http://localhost:8081/subjects/orders-value/versions/1

# Get latest schema
curl http://localhost:8081/subjects/orders-value/versions/latest

# Check compatibility
curl -X POST http://localhost:8081/compatibility/subjects/orders-value/versions/latest \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  -d '{"schema": "..."}'

# Set compatibility
curl -X PUT http://localhost:8081/config/orders-value \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  -d '{"compatibility": "FULL_TRANSITIVE"}'

# Delete a subject (careful!)
curl -X DELETE http://localhost:8081/subjects/orders-value
```

---

## Schema Evolution Example

```
Version 1: {order_id, user_id, amount}
Version 2: {order_id, user_id, amount, currency="USD"}  ← add with default: SAFE
Version 3: {order_id, user_id, amount, currency="USD", metadata=null} ← SAFE

NOT SAFE:
  {order_id, user_id, total}   ← renamed "amount" to "total"
  {order_id, user_id}          ← removed required field
```

---

## Using Confluent Cloud Schema Registry

```python
schema_registry = SchemaRegistryClient({
    "url":              "https://your-sr.confluent.cloud",
    "basic.auth.user.info": "API_KEY:API_SECRET",
})
```
