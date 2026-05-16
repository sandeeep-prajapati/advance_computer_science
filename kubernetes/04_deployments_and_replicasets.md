# Chapter 4: Deployments & ReplicaSets

## What is a Deployment?

A **Deployment** is the standard way to run stateless applications in Kubernetes. It manages:
- **Desired replica count** — always keeps N pods running
- **Rolling updates** — replace old pods with new ones with zero downtime
- **Rollback** — instantly go back to a previous version

Under the hood, a Deployment creates and manages a **ReplicaSet**, which manages the actual Pods.

```
Deployment
  └── ReplicaSet (v1: nginx:1.24)
        ├── Pod 1
        ├── Pod 2
        └── Pod 3
```

When you update the image:
```
Deployment
  ├── ReplicaSet (v1: nginx:1.24)  → scales down to 0
  └── ReplicaSet (v2: nginx:1.25)  → scales up to 3
```

---

## Creating a Deployment

**File:** `deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: default
  labels:
    app: my-app
spec:
  replicas: 3

  selector:
    matchLabels:
      app: my-app        # Must match template labels below

  template:
    metadata:
      labels:
        app: my-app      # Pods will get this label
    spec:
      containers:
        - name: my-app
          image: nginx:1.25-alpine
          ports:
            - containerPort: 80
          resources:
            requests:
              memory: "64Mi"
              cpu: "100m"
            limits:
              memory: "256Mi"
              cpu: "500m"
          readinessProbe:
            httpGet:
              path: /
              port: 80
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /
              port: 80
            initialDelaySeconds: 15
            periodSeconds: 20
```

```bash
kubectl apply -f deployment.yaml

kubectl get deployments
kubectl get pods -l app=my-app
```

---

## Scaling

```bash
# Scale to 5 replicas
kubectl scale deployment my-app --replicas=5

# Or edit the YAML and apply
kubectl apply -f deployment.yaml

# Watch pods scale up
kubectl get pods -w
```

---

## Rolling Updates (Zero Downtime)

### Update Strategy

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1           # Max extra pods during update
      maxUnavailable: 0     # Never take pods below desired count
```

| Strategy | Behavior |
|---|---|
| `RollingUpdate` (default) | Gradually replaces old pods with new ones |
| `Recreate` | Kills all old pods first, then creates new ones (causes downtime) |

### Trigger a Rolling Update

```bash
# Update the image
kubectl set image deployment/my-app my-app=nginx:1.26-alpine

# Watch the rollout
kubectl rollout status deployment/my-app

# See rollout history
kubectl rollout history deployment/my-app
```

### Annotate Updates (Best Practice)

```bash
kubectl set image deployment/my-app my-app=nginx:1.26-alpine \
  --record  # deprecated but still works

# Better: use --annotation
kubectl annotate deployment my-app \
  kubernetes.io/change-cause="Upgraded to nginx 1.26"
```

---

## Rollback

```bash
# Rollback to previous version
kubectl rollout undo deployment/my-app

# Rollback to a specific revision
kubectl rollout history deployment/my-app  # find revision number
kubectl rollout undo deployment/my-app --to-revision=2

# Check rollback status
kubectl rollout status deployment/my-app
```

---

## Readiness and Liveness Probes

**Liveness probe** — is the container alive? If it fails, Kubernetes **restarts** the container.

**Readiness probe** — is the container ready to receive traffic? If it fails, the pod is **removed from Service endpoints** (traffic stops routing to it).

```yaml
containers:
  - name: app
    image: my-app:latest

    # Readiness: HTTP check on /health
    readinessProbe:
      httpGet:
        path: /health
        port: 8080
      initialDelaySeconds: 5    # Wait 5s before first check
      periodSeconds: 10         # Check every 10s
      failureThreshold: 3       # 3 failures = not ready

    # Liveness: TCP socket check
    livenessProbe:
      tcpSocket:
        port: 8080
      initialDelaySeconds: 15
      periodSeconds: 20

    # Startup probe: gives slow-starting apps time to start
    startupProbe:
      httpGet:
        path: /health
        port: 8080
      failureThreshold: 30    # Wait up to 5 minutes (30 * 10s)
      periodSeconds: 10
```

---

## Environment Variables

```yaml
containers:
  - name: app
    image: my-app:latest
    env:
      # Direct value
      - name: APP_ENV
        value: "production"

      # From a ConfigMap
      - name: DB_HOST
        valueFrom:
          configMapKeyRef:
            name: app-config
            key: db_host

      # From a Secret
      - name: DB_PASSWORD
        valueFrom:
          secretKeyRef:
            name: app-secrets
            key: db_password

    # Load all keys from a ConfigMap as env vars
    envFrom:
      - configMapRef:
          name: app-config
      - secretRef:
          name: app-secrets
```

---

## Deployment Commands Reference

```bash
# Create
kubectl apply -f deployment.yaml

# View
kubectl get deployments
kubectl describe deployment my-app
kubectl get deployment my-app -o yaml

# Scale
kubectl scale deployment my-app --replicas=5

# Update image
kubectl set image deployment/my-app my-app=nginx:1.26

# Rollout status
kubectl rollout status deployment/my-app

# Rollout history
kubectl rollout history deployment/my-app

# Rollback
kubectl rollout undo deployment/my-app

# Pause a rollout (to make multiple changes before resuming)
kubectl rollout pause deployment/my-app
kubectl set image deployment/my-app my-app=nginx:1.26
kubectl set resources deployment/my-app -c=my-app --limits=cpu=200m,memory=512Mi
kubectl rollout resume deployment/my-app

# Delete
kubectl delete deployment my-app
```

---

## ReplicaSet (Direct Use — Rare)

You rarely create ReplicaSets directly (use Deployments instead). But it's useful to understand:

```bash
# View ReplicaSets (Deployment creates and manages these)
kubectl get replicasets
kubectl describe replicaset my-app-7d6b9f8c4
```

A Deployment keeps both old and new ReplicaSets around. Old ones are scaled to 0 (kept for rollback). The number of old ReplicaSets kept is controlled by `.spec.revisionHistoryLimit` (default: 10).

---

## Complete Example: App + Service + Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web-app
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: web-app
    spec:
      containers:
        - name: web-app
          image: nginx:1.25-alpine
          ports:
            - containerPort: 80
          resources:
            requests:
              memory: "64Mi"
              cpu: "100m"
            limits:
              memory: "256Mi"
              cpu: "500m"
          readinessProbe:
            httpGet:
              path: /
              port: 80
            initialDelaySeconds: 5
            periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: web-app-service
spec:
  selector:
    app: web-app
  ports:
    - port: 80
      targetPort: 80
  type: ClusterIP
```

```bash
kubectl apply -f deployment.yaml
kubectl get pods,deployments,services
kubectl port-forward svc/web-app-service 8080:80
```
