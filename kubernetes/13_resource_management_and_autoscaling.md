# Chapter 13: Resource Management & Autoscaling

## Resource Requests and Limits

Every container should declare `requests` and `limits` for CPU and memory.

```yaml
resources:
  requests:
    memory: "256Mi"    # Scheduler uses this to pick a node
    cpu: "250m"        # Guaranteed minimum
  limits:
    memory: "512Mi"    # OOMKilled if exceeded
    cpu: "1000m"       # Throttled (not killed) if exceeded
```

**Units:**
- CPU: `1` = 1 core, `500m` = 0.5 core, `100m` = 0.1 core
- Memory: `128Mi`, `1Gi`, `2Gi` (use Mi/Gi, not MB/GB)

---

## Quality of Service (QoS) Classes

Kubernetes assigns a QoS class based on resource settings:

| QoS Class | Condition | Eviction Priority |
|---|---|---|
| `Guaranteed` | limits == requests for all containers | Last to be evicted |
| `Burstable` | At least one container has requests < limits | Middle |
| `BestEffort` | No requests or limits set | First to be evicted |

```yaml
# Guaranteed (limits == requests)
resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

---

## LimitRange (Namespace Defaults)

```yaml
apiVersion: v1
kind: LimitRange
metadata:
  name: default-limits
  namespace: staging
spec:
  limits:
    - type: Container
      default:
        cpu: "500m"
        memory: "256Mi"
      defaultRequest:
        cpu: "100m"
        memory: "128Mi"
      max:
        cpu: "2"
        memory: "4Gi"
      min:
        cpu: "50m"
        memory: "64Mi"
```

---

## ResourceQuota (Namespace Caps)

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: team-quota
  namespace: team-a
spec:
  hard:
    requests.cpu: "4"
    requests.memory: "8Gi"
    limits.cpu: "8"
    limits.memory: "16Gi"
    pods: "20"
    services: "10"
    persistentvolumeclaims: "5"
    count/deployments.apps: "10"
```

```bash
kubectl describe resourcequota team-quota -n team-a
```

---

## Horizontal Pod Autoscaler (HPA)

HPA automatically scales the number of pod replicas based on CPU/memory usage or custom metrics.

### HPA based on CPU

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: my-app-hpa
  namespace: production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  minReplicas: 2
  maxReplicas: 20
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70   # Scale up when avg CPU > 70%
```

### HPA based on CPU + Memory

```yaml
metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### HPA Behavior (Tuning Scale Up/Down)

```yaml
spec:
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60   # Wait 60s after scale-up before another
      policies:
        - type: Pods
          value: 4                     # Add at most 4 pods per period
          periodSeconds: 60
    scaleDown:
      stabilizationWindowSeconds: 300  # Wait 5 min before scaling down
      policies:
        - type: Percent
          value: 25                    # Remove at most 25% per period
          periodSeconds: 60
```

```bash
# View HPA status
kubectl get hpa
kubectl describe hpa my-app-hpa

# Watch HPA in real time
kubectl get hpa -w

# Generate load to trigger scaling
kubectl run load-test --image=busybox --rm -it -- \
  /bin/sh -c "while true; do wget -q -O- http://my-app-service; done"
```

---

## Vertical Pod Autoscaler (VPA)

VPA automatically adjusts CPU and memory **requests** for containers (doesn't change replica count).

```bash
# Install VPA
kubectl apply -f https://github.com/kubernetes/autoscaler/releases/latest/download/vertical-pod-autoscaler.yaml
```

```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: my-app-vpa
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  updatePolicy:
    updateMode: "Auto"    # Auto, Off (recommend only), Initial
  resourcePolicy:
    containerPolicies:
      - containerName: my-app
        minAllowed:
          cpu: 50m
          memory: 64Mi
        maxAllowed:
          cpu: 4
          memory: 4Gi
```

```bash
# Check VPA recommendations
kubectl describe vpa my-app-vpa
```

> **Note:** Don't use HPA and VPA together on CPU/memory for the same deployment — they conflict. Use HPA on CPU and VPA on memory, or use KEDA for custom metrics.

---

## KEDA — Event-Driven Autoscaling

**KEDA** (Kubernetes Event-Driven Autoscaling) scales pods based on external events: queue length, HTTP request rate, database row count, etc.

```bash
# Install KEDA
helm repo add kedacore https://kedacore.github.io/charts
helm install keda kedacore/keda -n keda --create-namespace
```

```yaml
# Scale based on RabbitMQ queue depth
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: queue-consumer-scaler
spec:
  scaleTargetRef:
    name: queue-consumer
  minReplicaCount: 0    # Scale to zero when queue empty!
  maxReplicaCount: 50
  triggers:
    - type: rabbitmq
      metadata:
        host: amqp://rabbitmq:5672/
        queueName: my-tasks
        queueLength: "10"   # One pod per 10 messages
```

---

## Cluster Autoscaler

The **Cluster Autoscaler** adds/removes **nodes** when pods can't be scheduled due to resource constraints.

- Cloud providers: GKE, EKS, AKS have built-in node pool autoscaling
- Configure min/max node counts on your node pool

```bash
# GKE: enable autoscaling
gcloud container clusters update my-cluster \
  --enable-autoscaling \
  --min-nodes=1 \
  --max-nodes=10 \
  --node-pool=default-pool
```

---

## Resource Monitoring Commands

```bash
# Node resource usage
kubectl top nodes

# Pod resource usage (across all namespaces)
kubectl top pods -A

# Pod resource usage in a namespace
kubectl top pods -n production

# Sort by memory
kubectl top pods --sort-by=memory -A

# Describe a pod to see resource spec + events
kubectl describe pod my-pod -n production
```

---

## Autoscaling Cheatsheet

```bash
# Create HPA imperatively
kubectl autoscale deployment my-app --cpu-percent=70 --min=2 --max=10

# View HPA
kubectl get hpa
kubectl describe hpa my-app

# Delete HPA
kubectl delete hpa my-app

# View resource quotas
kubectl get resourcequota -A
kubectl describe resourcequota team-quota -n team-a

# Check limitrange
kubectl describe limitrange default-limits -n staging
```
