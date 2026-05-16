# Chapter 17: Stateful Workloads

## Running Databases on Kubernetes

Running databases in Kubernetes requires careful planning. Key concerns:
- **Persistent storage** — data must survive pod restarts
- **Stable network identity** — replicas need to address each other
- **Ordered operations** — primary must start before replicas
- **Backup strategy** — automated, tested backups

> For production, prefer **managed cloud databases** (RDS, Cloud SQL) over self-hosted on K8s unless you have the operational expertise. Self-hosted is great for dev/staging.

---

## MySQL on Kubernetes

### Single Instance (Dev/Staging)

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mysql-secret
  namespace: default
type: Opaque
stringData:
  MYSQL_ROOT_PASSWORD: "rootpassword"
  MYSQL_DATABASE: "myapp"
  MYSQL_USER: "appuser"
  MYSQL_PASSWORD: "apppassword"
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pvc
spec:
  accessModes: [ReadWriteOnce]
  storageClassName: standard
  resources:
    requests:
      storage: 20Gi
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  serviceName: mysql-headless
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
        - name: mysql
          image: mysql:8.0
          envFrom:
            - secretRef:
                name: mysql-secret
          ports:
            - containerPort: 3306
          volumeMounts:
            - name: data
              mountPath: /var/lib/mysql
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "2Gi"
              cpu: "1"
          livenessProbe:
            exec:
              command: ["mysqladmin", "ping", "-h", "localhost"]
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            exec:
              command: ["mysql", "-uroot", "-p$(MYSQL_ROOT_PASSWORD)", "-e", "SELECT 1"]
            initialDelaySeconds: 10
            periodSeconds: 5
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: mysql-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: mysql-headless
spec:
  clusterIP: None
  selector:
    app: mysql
  ports:
    - port: 3306
---
apiVersion: v1
kind: Service
metadata:
  name: mysql-service
spec:
  selector:
    app: mysql
  ports:
    - port: 3306
      targetPort: 3306
```

### MySQL Backup CronJob

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: mysql-backup
spec:
  schedule: "0 3 * * *"   # 3 AM daily
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: backup
              image: mysql:8.0
              command:
                - /bin/sh
                - -c
                - |
                  mysqldump -h mysql-service -u root -p$MYSQL_ROOT_PASSWORD --all-databases | \
                  gzip > /backup/full-$(date +%Y%m%d-%H%M%S).sql.gz
                  echo "Backup done"
              env:
                - name: MYSQL_ROOT_PASSWORD
                  valueFrom:
                    secretKeyRef:
                      name: mysql-secret
                      key: MYSQL_ROOT_PASSWORD
              volumeMounts:
                - name: backup
                  mountPath: /backup
          volumes:
            - name: backup
              persistentVolumeClaim:
                claimName: mysql-backup-pvc
```

---

## Redis on Kubernetes

### Redis Single Instance

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
        - name: redis
          image: redis:7-alpine
          command: ["redis-server", "--appendonly", "yes", "--requirepass", "$(REDIS_PASSWORD)"]
          env:
            - name: REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: redis-secret
                  key: password
          ports:
            - containerPort: 6379
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "1Gi"
              cpu: "500m"
          volumeMounts:
            - name: data
              mountPath: /data
          livenessProbe:
            exec:
              command: ["redis-cli", "ping"]
            periodSeconds: 10
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: redis-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: redis-service
spec:
  selector:
    app: redis
  ports:
    - port: 6379
```

### Redis with Helm (Production-ready)

```bash
helm install redis bitnami/redis \
  --set auth.password=secretpassword \
  --set master.persistence.size=10Gi \
  --set replica.replicaCount=2 \
  -n production
```

---

## Elasticsearch on Kubernetes

```bash
# Install ECK (Elastic Cloud on Kubernetes) operator
kubectl create -f https://download.elastic.co/downloads/eck/2.11.1/crds.yaml
kubectl apply -f https://download.elastic.co/downloads/eck/2.11.1/operator.yaml
```

```yaml
apiVersion: elasticsearch.k8s.elastic.co/v1
kind: Elasticsearch
metadata:
  name: my-elasticsearch
  namespace: elastic
spec:
  version: 8.12.0
  nodeSets:
    - name: default
      count: 3
      config:
        node.store.allow_mmap: false
      podTemplate:
        spec:
          containers:
            - name: elasticsearch
              resources:
                requests:
                  memory: 2Gi
                  cpu: 1
                limits:
                  memory: 4Gi
                  cpu: 2
      volumeClaimTemplates:
        - metadata:
            name: elasticsearch-data
          spec:
            accessModes: [ReadWriteOnce]
            storageClassName: fast-ssd
            resources:
              requests:
                storage: 50Gi
---
# Kibana
apiVersion: kibana.k8s.elastic.co/v1
kind: Kibana
metadata:
  name: my-kibana
  namespace: elastic
spec:
  version: 8.12.0
  count: 1
  elasticsearchRef:
    name: my-elasticsearch
```

---

## Kubernetes Operators

An **Operator** extends Kubernetes with domain knowledge for a specific application. It manages complex stateful apps (like databases) using custom controllers and CRDs.

| Operator | Manages |
|---|---|
| ECK | Elasticsearch + Kibana |
| MySQL Operator | MySQL InnoDB Cluster |
| Postgres Operator (Zalando) | PostgreSQL + replication |
| Redis Operator (OpsTree) | Redis cluster |
| Strimzi | Apache Kafka |
| Prometheus Operator | Prometheus + Alertmanager |

### Install PostgreSQL with Zalando Operator

```bash
helm install postgres-operator postgres-operator-charts/postgres-operator \
  -n postgres-operator --create-namespace
```

```yaml
apiVersion: acid.zalan.do/v1
kind: postgresql
metadata:
  name: my-postgres
  namespace: default
spec:
  teamId: "myteam"
  volume:
    size: 20Gi
  numberOfInstances: 2    # 1 primary + 1 replica
  users:
    appuser:
      - superuser
      - createdb
  databases:
    myapp: appuser
  postgresql:
    version: "16"
  resources:
    requests:
      cpu: 100m
      memory: 100Mi
    limits:
      cpu: 500m
      memory: 500Mi
```

---

## Production Recommendations

| Workload | Dev/Staging | Production |
|---|---|---|
| MySQL | StatefulSet in K8s | AWS RDS / Google Cloud SQL |
| PostgreSQL | Zalando Operator | AWS RDS / Cloud SQL |
| Redis | Deployment in K8s | AWS ElastiCache / Redis Cloud |
| Elasticsearch | ECK operator | AWS OpenSearch / Elastic Cloud |
| Kafka | Strimzi operator | AWS MSK / Confluent Cloud |

**Why managed services for production?**
- Automated backups
- Automated failover
- Security patches
- No operational overhead
- SLA guarantees

---

## Shared Volumes for Application Files

For apps needing shared file storage across pods (e.g., CMS uploads):

```bash
# AWS: Install EFS CSI Driver
helm install aws-efs-csi-driver \
  aws-efs-csi-driver/aws-efs-csi-driver \
  -n kube-system
```

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: efs-sc
provisioner: efs.csi.aws.com
parameters:
  provisioningMode: efs-ap
  fileSystemId: fs-xxxxxxxxx
  directoryPerms: "700"
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: shared-uploads
spec:
  storageClassName: efs-sc
  accessModes: [ReadWriteMany]   # Multiple pods can write
  resources:
    requests:
      storage: 100Gi
```
