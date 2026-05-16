# Chapter 18: Security

## Security Layers in Spark

```
1. Authentication   ← who is running this job?
2. Authorization    ← what are they allowed to do?
3. Encryption       ← is data protected in transit and at rest?
4. Audit            ← who did what and when?
```

---

## Authentication

### Kerberos (Enterprise Hadoop)

Kerberos is the standard authentication mechanism in Hadoop clusters:

```bash
# Obtain a Kerberos ticket
kinit spark-user@REALM.COMPANY.COM

# Submit with Kerberos
spark-submit \
  --master yarn \
  --principal spark-user@REALM.COMPANY.COM \
  --keytab /etc/security/keytabs/spark.keytab \
  my_job.py
```

```python
# In spark-defaults.conf or job config
spark.conf.set("spark.kerberos.principal", "spark-user@REALM.COM")
spark.conf.set("spark.kerberos.keytab",    "/etc/security/keytabs/spark.keytab")
```

### Shared Secret (Standalone / Local)

```python
spark.conf.set("spark.authenticate",       "true")
spark.conf.set("spark.authenticate.secret","my-secret-key")
```

---

## Encryption in Transit (SSL/TLS)

### RPC Encryption (Node-to-Node)

```python
spark.conf.set("spark.network.crypto.enabled",   "true")
spark.conf.set("spark.network.crypto.keyLength",  "128")
spark.conf.set("spark.authenticate",              "true")
```

### SSL for Spark UI and History Server

```properties
# spark-defaults.conf
spark.ssl.enabled=true
spark.ssl.keyStore=/path/to/keystore.jks
spark.ssl.keyStorePassword=keystore_pass
spark.ssl.trustStore=/path/to/truststore.jks
spark.ssl.trustStorePassword=truststore_pass
spark.ssl.protocol=TLSv1.2
```

```python
spark.conf.set("spark.ui.filters",        "org.apache.hadoop.security.authentication.server.AuthenticationFilter")
spark.conf.set("spark.ssl.ui.enabled",    "true")
```

---

## Encryption at Rest

Data at rest encryption is handled at the storage layer:

### HDFS Transparent Encryption

```bash
# Create encryption zone in HDFS
hdfs crypto -createZone -keyName myKey -path /user/spark/data

# Verify
hdfs crypto -listZones
```

### S3 Server-Side Encryption

```python
# SSE-S3 (AWS managed keys)
spark.conf.set("spark.hadoop.fs.s3a.server-side-encryption-algorithm", "AES256")

# SSE-KMS (customer managed keys)
spark.conf.set("spark.hadoop.fs.s3a.server-side-encryption-algorithm", "SSE-KMS")
spark.conf.set("spark.hadoop.fs.s3a.server-side-encryption-key",
               "arn:aws:kms:us-east-1:123456789:key/my-key-id")

# Client-side encryption
spark.conf.set("spark.hadoop.fs.s3a.encryption.algorithm", "CSE-KMS")
```

---

## Credential Management

### Hadoop Credential Provider

Store passwords in an encrypted credential store instead of plain-text config:

```bash
# Create credential store
hadoop credential create spark.password.db \
  -provider jceks://hdfs/user/spark/credentials.jceks \
  -value mySecretPassword

# Use it in config
spark.conf.set("spark.hadoop.hadoop.security.credential.provider.path",
               "jceks://hdfs/user/spark/credentials.jceks")
```

### Vault / Secrets Manager

```python
# AWS Secrets Manager
import boto3, json

def get_db_credentials():
    client = boto3.client("secretsmanager", region_name="us-east-1")
    secret = client.get_secret_value(SecretId="spark/db-credentials")
    return json.loads(secret["SecretString"])

creds = get_db_credentials()
df = spark.read.format("jdbc") \
    .option("user",     creds["username"]) \
    .option("password", creds["password"]) \
    .load()
```

---

## ACLs and Authorization

### Spark UI ACLs

Restrict who can view the Spark UI:

```python
spark.conf.set("spark.acls.enable",      "true")
spark.conf.set("spark.ui.view.acls",     "alice,bob")
spark.conf.set("spark.ui.view.acls.groups", "data-team")
spark.conf.set("spark.admin.acls",       "spark-admin")
spark.conf.set("spark.admin.acls.groups","ops-team")
```

### Modify ACLs (Cancel Jobs, Kill Tasks)

```python
spark.conf.set("spark.modify.acls",      "alice")
spark.conf.set("spark.modify.acls.groups","ops-team")
```

---

## Secure Spark on Kubernetes

```yaml
# Use a dedicated service account with minimal permissions
spark.kubernetes.authenticate.driver.serviceAccountName=spark

# Mount secrets as environment variables
spark.kubernetes.driverEnv.DB_PASSWORD=<SECRET_NAME>
```

```bash
# Create K8s secret for DB credentials
kubectl create secret generic spark-db-creds \
  --from-literal=username=spark_user \
  --from-literal=password=secret_pass \
  -n spark-jobs

# Reference in spark-submit
--conf spark.kubernetes.driver.secretKeyRef.DB_USER=spark-db-creds:username
--conf spark.kubernetes.driver.secretKeyRef.DB_PASS=spark-db-creds:password
```

---

## Data Masking / Anonymization

For PII compliance:

```python
from pyspark.sql import functions as F

def mask_pii(df):
    return df \
        .withColumn("email",  F.regexp_replace("email", r"^(.{2}).*(@.*)$", "$1***$2")) \
        .withColumn("phone",  F.regexp_replace("phone", r"\d{7}$", "*******")) \
        .withColumn("ssn",    F.lit("***-**-****")) \
        .withColumn("ip",     F.regexp_replace("ip", r"\d+$", "0"))   # anonymize last octet
```

### Tokenization

```python
# Replace PII with a stable hash for correlation without exposing real values
df = df.withColumn(
    "user_token",
    F.sha2(F.concat_ws("|", F.col("email"), F.lit("SALT_SECRET")), 256)
)
```

---

## Security Checklist

| Area | Action |
|---|---|
| Authentication | Enable Kerberos or shared secret |
| Encryption in transit | Enable `spark.network.crypto` or SSL |
| Encryption at rest | HDFS zones or S3 SSE |
| Credentials | Never hardcode passwords — use Vault or credential provider |
| Spark UI | Enable ACLs; don't expose port 4040 publicly |
| Kubernetes | Minimal RBAC; don't use `cluster-admin` for Spark SA |
| PII | Mask/tokenize before writing to data lake |
| Audit | Enable event logs; retain for compliance window |
