# Chapter 13: Security

## Security Layers

```
1. Encryption in Transit   ← TLS between clients and brokers
2. Authentication          ← Who is this client? (SASL/TLS)
3. Authorization           ← What can this client do? (ACLs)
4. Encryption at Rest      ← Disk encryption (OS/cloud level)
5. Audit Logging           ← Who did what and when
```

---

## TLS (Encryption in Transit)

### Generate Certificates

```bash
# 1. Create CA
openssl req -new -x509 -keyout ca.key -out ca.crt -days 365 \
  -subj "/CN=Kafka-CA/O=MyOrg/C=US" -passout pass:ca-secret

# 2. Create broker keystore
keytool -genkeypair -alias broker1 \
  -keyalg RSA -keysize 2048 \
  -keystore broker1.keystore.jks \
  -dname "CN=broker1.example.com,O=MyOrg,C=US" \
  -storepass broker-secret -keypass broker-secret

# 3. Create CSR
keytool -certreq -alias broker1 \
  -keystore broker1.keystore.jks \
  -file broker1.csr \
  -storepass broker-secret

# 4. Sign with CA
openssl x509 -req -CA ca.crt -CAkey ca.key -in broker1.csr \
  -out broker1.crt -days 365 -CAcreateserial -passin pass:ca-secret

# 5. Import CA cert + signed cert into keystore
keytool -importcert -alias ca-cert -file ca.crt \
  -keystore broker1.keystore.jks -storepass broker-secret -noprompt
keytool -importcert -alias broker1 -file broker1.crt \
  -keystore broker1.keystore.jks -storepass broker-secret

# 6. Create truststore (contains CA cert)
keytool -importcert -alias ca-cert -file ca.crt \
  -keystore kafka.truststore.jks -storepass truststore-secret -noprompt
```

### Broker TLS Config (server.properties)

```properties
listeners=PLAINTEXT://0.0.0.0:9092,SSL://0.0.0.0:9093
advertised.listeners=PLAINTEXT://broker1:9092,SSL://broker1:9093
listener.security.protocol.map=PLAINTEXT:PLAINTEXT,SSL:SSL

ssl.keystore.location=/etc/kafka/ssl/broker1.keystore.jks
ssl.keystore.password=broker-secret
ssl.key.password=broker-secret
ssl.truststore.location=/etc/kafka/ssl/kafka.truststore.jks
ssl.truststore.password=truststore-secret
ssl.client.auth=required          # require mutual TLS
ssl.protocol=TLSv1.3
ssl.enabled.protocols=TLSv1.3,TLSv1.2
```

### Client TLS Config

```python
producer = Producer({
    "bootstrap.servers":       "broker1:9093",
    "security.protocol":       "SSL",
    "ssl.ca.location":         "/etc/kafka/ssl/ca.crt",
    "ssl.certificate.location":"/etc/kafka/ssl/client.crt",
    "ssl.key.location":        "/etc/kafka/ssl/client.key",
    "ssl.key.password":        "client-secret",
})
```

---

## SASL Authentication

### SASL/PLAIN (Username/Password — Use Over TLS)

```properties
# server.properties
listeners=SASL_SSL://0.0.0.0:9093
security.inter.broker.protocol=SASL_SSL
sasl.mechanism.inter.broker.protocol=PLAIN
sasl.enabled.mechanisms=PLAIN

# JAAS config (kafka_server_jaas.conf)
KafkaServer {
  org.apache.kafka.common.security.plain.PlainLoginModule required
  username="admin"
  password="admin-secret"
  user_admin="admin-secret"
  user_producer="producer-secret"
  user_consumer="consumer-secret";
};
```

```bash
export KAFKA_OPTS="-Djava.security.auth.login.config=/etc/kafka/kafka_server_jaas.conf"
```

```python
# Client config
producer = Producer({
    "bootstrap.servers":             "broker1:9093",
    "security.protocol":             "SASL_SSL",
    "sasl.mechanism":                "PLAIN",
    "sasl.username":                 "producer",
    "sasl.password":                 "producer-secret",
    "ssl.ca.location":               "/etc/kafka/ssl/ca.crt",
})
```

### SASL/SCRAM (Stored in ZooKeeper/KRaft)

```bash
# Create user
kafka-configs --bootstrap-server localhost:9092 \
  --alter --entity-type users --entity-name alice \
  --add-config 'SCRAM-SHA-256=[iterations=8192,password=alice-secret]'

# List users
kafka-configs --bootstrap-server localhost:9092 \
  --describe --entity-type users
```

```properties
# server.properties
sasl.enabled.mechanisms=SCRAM-SHA-256
sasl.mechanism.inter.broker.protocol=SCRAM-SHA-256
```

### SASL/GSSAPI (Kerberos)

```properties
# server.properties
sasl.enabled.mechanisms=GSSAPI
sasl.kerberos.service.name=kafka

# JAAS
KafkaServer {
  com.sun.security.auth.module.Krb5LoginModule required
  useKeyTab=true
  storeKey=true
  keyTab="/etc/security/keytabs/kafka.keytab"
  principal="kafka/broker1.example.com@REALM.COM";
};
```

---

## ACLs (Authorization)

### Enable ACLs

```properties
# server.properties
authorizer.class.name=org.apache.kafka.metadata.authorizer.StandardAuthorizer
allow.everyone.if.no.acl.found=false   # deny by default
super.users=User:admin
```

### Manage ACLs

```bash
# Allow producer to write to topic
kafka-acls --bootstrap-server localhost:9092 \
  --add \
  --allow-principal User:producer-service \
  --operation Write \
  --operation Describe \
  --topic orders

# Allow consumer group to read from topic
kafka-acls --bootstrap-server localhost:9092 \
  --add \
  --allow-principal User:consumer-service \
  --operation Read \
  --operation Describe \
  --topic orders

# Allow consumer to use a consumer group
kafka-acls --bootstrap-server localhost:9092 \
  --add \
  --allow-principal User:consumer-service \
  --operation Read \
  --group order-processor

# Allow all operations on a topic (admin)
kafka-acls --bootstrap-server localhost:9092 \
  --add \
  --allow-principal User:admin \
  --operation All \
  --topic '*'

# List ACLs
kafka-acls --bootstrap-server localhost:9092 --list --topic orders

# Delete ACL
kafka-acls --bootstrap-server localhost:9092 \
  --remove \
  --allow-principal User:old-service \
  --operation Write \
  --topic orders
```

---

## ACL Operations Reference

| Resource Type | Operations |
|---|---|
| Topic | Read, Write, Create, Delete, Alter, Describe, DescribeConfigs, AlterConfigs |
| Group | Read, Delete, Describe |
| Cluster | Create, Alter, Describe, AlterConfigs, IdempotentWrite, ClusterAction |
| TransactionalId | Write, Describe |

---

## Encryption at Rest

Kafka does not encrypt data at rest natively. Use:

```bash
# Linux LUKS disk encryption
cryptsetup luksFormat /dev/sdb
cryptsetup luksOpen /dev/sdb kafka-data
mkfs.xfs /dev/mapper/kafka-data
mount /dev/mapper/kafka-data /var/kafka/data

# Cloud provider solutions:
# AWS: EBS encryption (KMS), S3 SSE for tiered storage
# GCP: GCE persistent disk encryption (default)
# Azure: Azure Disk Encryption
```

---

## Security Best Practices

```
1. Always use TLS in production (SASL_SSL, not SASL_PLAINTEXT)
2. Use SASL/SCRAM or GSSAPI (not PLAIN) for authentication
3. Enable ACLs with deny-by-default
4. Give each service its own credentials (principle of least privilege)
5. Rotate credentials periodically
6. Never use super.users in application code
7. Use mutual TLS for inter-broker communication
8. Enable audit logging for compliance
9. Separate networks: brokers on internal network, expose only via load balancer
10. Encrypt at rest using cloud provider disk encryption
```

---

## Audit Logging

```properties
# Enable log4j audit (add to log4j.properties)
log4j.logger.kafka.authorizer.logger=INFO, authorizerAppender
log4j.additivity.kafka.authorizer.logger=false
log4j.appender.authorizerAppender=org.apache.log4j.DailyRollingFileAppender
log4j.appender.authorizerAppender.DatePattern='.'yyyy-MM-dd-HH
log4j.appender.authorizerAppender.File=${kafka.logs.dir}/kafka-authorizer.log
log4j.appender.authorizerAppender.layout=org.apache.log4j.PatternLayout
log4j.appender.authorizerAppender.layout.ConversionPattern=[%d] %p %m (%c)%n
```

Audit log records `ALLOWED` and `DENIED` authorization events with principal, operation, and resource.
