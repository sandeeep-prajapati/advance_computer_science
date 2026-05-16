# Chapter 11: StatefulSets & DaemonSets

## StatefulSets

### What is a StatefulSet?

A **StatefulSet** is like a Deployment but for **stateful applications** that need:
- **Stable, persistent pod names** (pod-0, pod-1, pod-2 — not random hashes)
- **Stable network identities** (each pod gets its own DNS hostname)
- **Ordered deployment and scaling** (pod-0 starts before pod-1)
- **Per-pod persistent volumes** (each pod gets its own PVC)

### When to Use StatefulSet vs Deployment

| Feature | Deployment | StatefulSet |
|---|---|---|
| Pod names | Random (nginx-abc123) | Stable (mysql-0, mysql-1) |
| Start order | Any order | Ordered (0, 1, 2...) |
| Storage | Shared PVC or emptyDir | Each pod has own PVC |
| Use case | Stateless apps (web, API) | Databases, queues, caches |

---

### StatefulSet Example: MySQL

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
  namespace: default
spec:
  serviceName: mysql-headless       # Must match headless service name
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
          ports:
            - containerPort: 3306
          env:
            - name: MYSQL_ROOT_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: mysql-secret
                  key: root-password
          volumeMounts:
            - name: mysql-data
              mountPath: /var/lib/mysql
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "2Gi"
              cpu: "1"

  # Each pod gets its own PVC
  volumeClaimTemplates:
    - metadata:
        name: mysql-data
      spec:
        accessModes: ["ReadWriteOnce"]
        storageClassName: fast-ssd
        resources:
          requests:
            storage: 20Gi
---
# Headless Service (required for stable DNS)
apiVersion: v1
kind: Service
metadata:
  name: mysql-headless
spec:
  clusterIP: None           # Headless — no virtual IP
  selector:
    app: mysql
  ports:
    - port: 3306
---
# Regular Service for app access
apiVersion: v1
kind: Service
metadata:
  name: mysql-service
spec:
  selector:
    app: mysql
  ports:
    - port: 3306
```

**DNS for each pod:**
```
mysql-0.mysql-headless.default.svc.cluster.local
mysql-1.mysql-headless.default.svc.cluster.local
```

---

### StatefulSet Scaling

```bash
# Scale up (pods added in order: mysql-2, then mysql-3)
kubectl scale statefulset mysql --replicas=3

# Scale down (pods removed in reverse order: mysql-2 first)
kubectl scale statefulset mysql --replicas=1

# Watch ordered creation
kubectl get pods -w -l app=mysql
```

---

### StatefulSet Update Strategies

```yaml
spec:
  updateStrategy:
    type: RollingUpdate    # Default — updates pods one at a time in reverse order
    rollingUpdate:
      partition: 0         # Only update pods with ordinal >= partition
```

```bash
# Canary update: only update mysql-1 and mysql-2 (not mysql-0)
kubectl patch statefulset mysql -p '{"spec":{"updateStrategy":{"rollingUpdate":{"partition":1}}}}'
kubectl set image statefulset/mysql mysql=mysql:8.1
# After validation:
kubectl patch statefulset mysql -p '{"spec":{"updateStrategy":{"rollingUpdate":{"partition":0}}}}'
```

---

## DaemonSets

### What is a DaemonSet?

A **DaemonSet** ensures that **one pod runs on every node** (or a subset of nodes). When a new node is added to the cluster, the pod is automatically scheduled on it.

**Use cases:**
- Log collection agents (Fluentd, Filebeat)
- Node monitoring agents (Prometheus Node Exporter, Datadog)
- Network plugins (CNI plugins)
- Security scanning agents

---

### DaemonSet Example: Log Collector

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd
  namespace: kube-system
  labels:
    app: fluentd
spec:
  selector:
    matchLabels:
      app: fluentd
  template:
    metadata:
      labels:
        app: fluentd
    spec:
      tolerations:
        # Also run on master/control-plane nodes
        - key: node-role.kubernetes.io/control-plane
          operator: Exists
          effect: NoSchedule

      containers:
        - name: fluentd
          image: fluent/fluentd-kubernetes-daemonset:v1-debian-elasticsearch
          env:
            - name: FLUENT_ELASTICSEARCH_HOST
              value: elasticsearch-service
            - name: FLUENT_ELASTICSEARCH_PORT
              value: "9200"
          resources:
            limits:
              memory: 200Mi
            requests:
              cpu: 100m
              memory: 200Mi
          volumeMounts:
            - name: varlog
              mountPath: /var/log
            - name: varlibdockercontainers
              mountPath: /var/lib/docker/containers
              readOnly: true

      volumes:
        - name: varlog
          hostPath:
            path: /var/log
        - name: varlibdockercontainers
          hostPath:
            path: /var/lib/docker/containers
```

---

### DaemonSet on Specific Nodes Only

```yaml
spec:
  template:
    spec:
      # Only run on nodes with this label
      nodeSelector:
        dedicated: monitoring

      # Or use affinity for more complex rules
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
              - matchExpressions:
                  - key: node-type
                    operator: In
                    values: [gpu, high-memory]
```

---

### DaemonSet Commands

```bash
# Create
kubectl apply -f daemonset.yaml

# List (shows desired vs ready per node)
kubectl get daemonsets
kubectl get daemonsets -n kube-system

# Describe
kubectl describe daemonset fluentd -n kube-system

# Check which nodes the DS is running on
kubectl get pods -l app=fluentd -o wide

# Rolling update (same as Deployment)
kubectl set image daemonset/fluentd fluentd=fluent/fluentd:v1.16
kubectl rollout status daemonset/fluentd

# Rollback
kubectl rollout undo daemonset/fluentd
```

---

## Comparison Summary

| Feature | Deployment | StatefulSet | DaemonSet |
|---|---|---|---|
| Pod count control | `replicas: N` | `replicas: N` | One per node |
| Pod name | Random hash | Stable (0, 1, 2) | Per node |
| Storage | Shared/none | Per-pod PVC | HostPath |
| Order | Arbitrary | Ordered | Per node |
| Use case | Stateless apps | DBs, caches, queues | Node agents |
| Scale up | Instant, any order | Ordered (0→1→2) | Auto when node joins |
| Scale down | Any order | Reverse order (2→1→0) | Auto when node leaves |
