# Chapter 16: PostgreSQL with Docker & Kubernetes

## Docker Quick Start

```bash
# Pull and run PostgreSQL
docker run -d \
  --name postgres \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypassword \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql/data \
  postgres:16-alpine

# Connect
docker exec -it postgres psql -U myuser -d mydb

# Stop/start
docker stop postgres
docker start postgres

# View logs
docker logs postgres
docker logs -f postgres   # follow
```

---

## Docker Compose for Development

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    container_name: myapp-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: myuser
      POSTGRES_PASSWORD: mypassword
      POSTGRES_DB: mydb
      POSTGRES_INITDB_ARGS: "--locale=en_US.UTF-8"
      TZ: UTC
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init:/docker-entrypoint-initdb.d    # init scripts run on first start
      - ./postgresql.conf:/etc/postgresql/postgresql.conf  # custom config
    command: postgres -c config_file=/etc/postgresql/postgresql.conf
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U myuser -d mydb"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s

  pgadmin:
    image: dpage/pgadmin4:latest
    container_name: myapp-pgadmin
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@example.com
      PGADMIN_DEFAULT_PASSWORD: admin
    ports:
      - "5050:80"
    depends_on:
      postgres:
        condition: service_healthy

volumes:
  postgres_data:
```

```bash
docker-compose up -d
docker-compose logs -f postgres
docker-compose exec postgres psql -U myuser -d mydb
docker-compose down              # stop but keep volumes
docker-compose down -v           # stop and remove volumes
```

---

## Init Scripts

Files in `/docker-entrypoint-initdb.d/` run automatically on the first start (when data directory is empty):

```
init/
├── 01_schema.sql      # create tables
├── 02_seed.sql        # insert test data
└── 03_extensions.sql  # create extensions
```

```sql
-- 01_schema.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

```bash
-- 02_seed.sql
INSERT INTO users (email, name) VALUES
    ('alice@example.com', 'Alice'),
    ('bob@example.com', 'Bob');
```

---

## Custom postgresql.conf

```ini
# postgresql.conf (mount into container)
# Memory
shared_buffers = 256MB
work_mem = 16MB
maintenance_work_mem = 64MB

# Connections
max_connections = 100

# Logging (for development)
log_statement = 'all'
log_duration = on
log_min_duration_statement = 0

# Performance
random_page_cost = 1.1    # SSD-optimized
effective_cache_size = 1GB
```

---

## Multiple PostgreSQL Versions

```yaml
version: '3.8'
services:
  pg14:
    image: postgres:14-alpine
    ports:
      - "5433:5432"     # different host port
    environment:
      POSTGRES_PASSWORD: postgres

  pg16:
    image: postgres:16-alpine
    ports:
      - "5434:5432"
    environment:
      POSTGRES_PASSWORD: postgres
```

---

## PostgreSQL on Kubernetes

### StatefulSet (basic setup):
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
spec:
  serviceName: postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:16-alpine
          env:
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef:
                  name: postgres-secret
                  key: username
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: postgres-secret
                  key: password
            - name: POSTGRES_DB
              value: mydb
          ports:
            - containerPort: 5432
          volumeMounts:
            - name: postgres-storage
              mountPath: /var/lib/postgresql/data
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "2Gi"
              cpu: "1000m"
          livenessProbe:
            exec:
              command: ["pg_isready", "-U", "myuser"]
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            exec:
              command: ["pg_isready", "-U", "myuser"]
            initialDelaySeconds: 5
            periodSeconds: 5
  volumeClaimTemplates:
    - metadata:
        name: postgres-storage
      spec:
        accessModes: ["ReadWriteOnce"]
        storageClassName: fast-ssd
        resources:
          requests:
            storage: 100Gi
```

### Service:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: postgres
spec:
  selector:
    app: postgres
  ports:
    - port: 5432
      targetPort: 5432
  type: ClusterIP
```

### Secret:
```bash
kubectl create secret generic postgres-secret \
  --from-literal=username=myuser \
  --from-literal=password=$(openssl rand -base64 32)
```

---

## CloudNativePG Operator (Production Recommendation)

CloudNativePG is the Kubernetes operator for PostgreSQL — manages HA, failover, backups automatically.

```bash
# Install operator
kubectl apply -f https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/release-1.22/releases/cnpg-1.22.0.yaml
```

```yaml
# postgresql-cluster.yaml
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: my-postgres-cluster
spec:
  instances: 3                    # 1 primary + 2 replicas

  postgresql:
    version: "16"
    parameters:
      shared_buffers: "512MB"
      max_connections: "200"

  storage:
    size: 100Gi
    storageClass: fast-ssd

  backup:
    retentionPolicy: "30d"
    barmanObjectStore:
      destinationPath: s3://my-postgres-backups/
      s3Credentials:
        accessKeyId:
          name: aws-creds
          key: ACCESS_KEY_ID
        secretAccessKey:
          name: aws-creds
          key: SECRET_ACCESS_KEY

  monitoring:
    enablePodMonitor: true   # Prometheus scraping

  resources:
    requests:
      memory: "1Gi"
      cpu: "500m"
    limits:
      memory: "4Gi"
      cpu: "2000m"
```

```bash
kubectl apply -f postgresql-cluster.yaml

# Check status
kubectl get clusters
kubectl describe cluster my-postgres-cluster

# Connect
kubectl exec -it my-postgres-cluster-1 -- psql -U app mydb

# Switchover (graceful failover)
kubectl cnpg promote my-postgres-cluster my-postgres-cluster-2

# Take a backup
kubectl cnpg backup my-postgres-cluster
```

---

## Persistence Checklist

- Use a **StorageClass** with `Retain` reclaim policy (not `Delete`) — prevents data loss on PVC deletion
- Enable **volume snapshots** for point-in-time storage snapshots
- Use **StatefulSet** not Deployment (stable pod names and volumes)
- Configure **resource limits** to prevent OOM kills
- Set **PodDisruptionBudget** to prevent all pods being killed during node drain
- Back up to **object storage** (S3/GCS), not just local disk

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: postgres-pdb
spec:
  minAvailable: 1   # always keep at least 1 instance running
  selector:
    matchLabels:
      app: postgres
```
