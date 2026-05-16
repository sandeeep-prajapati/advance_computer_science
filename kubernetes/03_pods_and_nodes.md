# Chapter 3: Pods & Nodes

## What is a Pod?

A **Pod** is the smallest deployable unit in Kubernetes. Every container you run lives inside a Pod. A Pod:
- Wraps one or more containers
- Gives them a **shared IP address**
- Gives them **shared storage volumes**
- Lives and dies together — all containers start and stop as a unit

> **Rule:** In 99% of cases, one Pod = one container. Multi-container Pods are for tightly-coupled helpers (sidecars).

---

## Pod Lifecycle

```
Pending → Running → Succeeded (or Failed)
              ↓
           (if crash)
           CrashLoopBackOff
```

| Phase | Meaning |
|---|---|
| `Pending` | Pod accepted, waiting for a node/image pull |
| `Running` | At least one container is running |
| `Succeeded` | All containers exited with code 0 (Jobs) |
| `Failed` | At least one container exited with non-zero |
| `Unknown` | Node communication lost |
| `CrashLoopBackOff` | Container keeps crashing and restarting |

---

## Writing Your First Pod YAML

**File:** `pod.yaml`

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-nginx
  namespace: default
  labels:
    app: nginx
    env: dev
spec:
  containers:
    - name: nginx
      image: nginx:1.25-alpine
      ports:
        - containerPort: 80
      resources:
        requests:
          memory: "64Mi"
          cpu: "100m"
        limits:
          memory: "128Mi"
          cpu: "200m"
```

```bash
# Apply the pod
kubectl apply -f pod.yaml

# Check status
kubectl get pods

# Describe the pod (see events, IP, node)
kubectl describe pod my-nginx

# Get logs
kubectl logs my-nginx

# Exec into the container
kubectl exec -it my-nginx -- sh

# Delete the pod
kubectl delete pod my-nginx
```

---

## Multi-Container Pods

Use multi-container pods when containers are tightly coupled. Common patterns:

### Sidecar Pattern

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-with-log-sidecar
spec:
  containers:
    # Main application
    - name: app
      image: my-app:latest
      volumeMounts:
        - name: shared-logs
          mountPath: /var/log/app

    # Sidecar: ships logs to external system
    - name: log-shipper
      image: fluent/fluent-bit:latest
      volumeMounts:
        - name: shared-logs
          mountPath: /var/log/app
          readOnly: true

  volumes:
    - name: shared-logs
      emptyDir: {}
```

### Init Container Pattern

Init containers run **before** the main container starts — useful for setup tasks (wait for DB, run migrations).

```yaml
spec:
  initContainers:
    - name: wait-for-db
      image: busybox
      command: ['sh', '-c', 'until nc -z mysql-service 3306; do echo waiting; sleep 2; done']

  containers:
    - name: app
      image: my-app:latest
```

---

## Resource Requests and Limits

```yaml
resources:
  requests:
    memory: "256Mi"   # Guaranteed minimum
    cpu: "250m"       # 250 millicores = 0.25 CPU
  limits:
    memory: "512Mi"   # Hard cap — OOMKill if exceeded
    cpu: "500m"       # Throttled if exceeded (not killed)
```

| Type | Behavior |
|---|---|
| `requests` | Used by scheduler to find a node with enough capacity |
| `limits.memory` | Container is OOMKilled if it exceeds this |
| `limits.cpu` | Container is CPU-throttled (not killed) |

**CPU units:**
- `1` = 1 full CPU core
- `500m` = half a core
- `100m` = 10% of a core

---

## Pod Labels & Selectors

Labels are key-value pairs attached to objects. Services, Deployments, and network policies use **selectors** to target pods by label.

```yaml
metadata:
  labels:
    app: my-app
    version: v2
    env: production
```

```bash
# List pods with label selectors
kubectl get pods -l app=my-app
kubectl get pods -l env=production,version=v2

# Add a label to an existing pod
kubectl label pod my-nginx color=blue

# Remove a label
kubectl label pod my-nginx color-
```

---

## Pod Scheduling: Node Selection

### nodeSelector (simple)

```yaml
spec:
  nodeSelector:
    disktype: ssd
```

Label the node first:
```bash
kubectl label node my-node disktype=ssd
```

### nodeAffinity (advanced)

```yaml
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
          - matchExpressions:
              - key: kubernetes.io/arch
                operator: In
                values: [amd64]
      preferredDuringSchedulingIgnoredDuringExecution:
        - weight: 1
          preference:
            matchExpressions:
              - key: disktype
                operator: In
                values: [ssd]
```

### Taints & Tolerations

**Taints** repel pods from nodes. **Tolerations** allow specific pods to land on tainted nodes.

```bash
# Taint a node (repels all pods without a toleration)
kubectl taint nodes gpu-node dedicated=gpu:NoSchedule
```

```yaml
# Pod that can tolerate the taint
spec:
  tolerations:
    - key: "dedicated"
      operator: "Equal"
      value: "gpu"
      effect: "NoSchedule"
```

---

## Nodes

### View Nodes

```bash
# List nodes with status
kubectl get nodes

# Detailed node info (capacity, conditions, images)
kubectl describe node my-node

# Node resource usage
kubectl top nodes
```

### Node Conditions

| Condition | Meaning |
|---|---|
| `Ready` | Node is healthy and accepting pods |
| `MemoryPressure` | Node is running low on memory |
| `DiskPressure` | Node is running low on disk |
| `PIDPressure` | Too many processes running on the node |
| `NetworkUnavailable` | Network not configured correctly |

### Drain and Cordon

```bash
# Cordon: mark node as unschedulable (no new pods, existing stay)
kubectl cordon my-node

# Drain: evict all pods + cordon (for maintenance)
kubectl drain my-node --ignore-daemonsets --delete-emptydir-data

# Uncordon: mark as schedulable again
kubectl uncordon my-node
```

---

## Useful Pod Commands Cheatsheet

```bash
# Get pods in all namespaces
kubectl get pods --all-namespaces
kubectl get pods -A

# Get pods on a specific node
kubectl get pods --field-selector spec.nodeName=my-node

# Watch pods in real time
kubectl get pods -w

# Get pod YAML definition
kubectl get pod my-pod -o yaml

# Copy a file from/to a pod
kubectl cp my-pod:/etc/nginx/nginx.conf ./nginx.conf
kubectl cp ./nginx.conf my-pod:/tmp/nginx.conf

# Run a temporary debug pod
kubectl run debug --image=busybox --rm -it --restart=Never -- sh
```
