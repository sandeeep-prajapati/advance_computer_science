# Chapter 5: Services & Networking

## Why Services Exist

Pods are **ephemeral** — they get new IP addresses every time they restart. A **Service** provides a stable, permanent endpoint to reach a set of pods regardless of their IPs.

```
Without Service:
  Client → Pod IP (10.244.1.5) → Pod dies → new IP (10.244.1.9) → client broken

With Service:
  Client → Service IP (10.96.45.200) → Service routes to healthy pods → always works
```

---

## Service Types

| Type | Accessible From | Use Case |
|---|---|---|
| `ClusterIP` (default) | Inside cluster only | Internal service-to-service communication |
| `NodePort` | Outside cluster via Node IP + port | Dev/testing, simple external access |
| `LoadBalancer` | Public internet | Production on cloud providers |
| `ExternalName` | Maps to external DNS | Route to external DB/API |

---

## ClusterIP (Default)

Used for internal communication between services. Kubernetes assigns a virtual IP that only exists inside the cluster.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app-service
  namespace: default
spec:
  type: ClusterIP       # Default — can omit this line
  selector:
    app: my-app         # Routes to pods with this label
  ports:
    - protocol: TCP
      port: 80          # Service port (what clients connect to)
      targetPort: 8080  # Container port (what your app listens on)
```

```bash
kubectl apply -f service.yaml
kubectl get services
# my-app-service   ClusterIP   10.96.45.200   <none>   80/TCP

# Test from inside cluster
kubectl run test --image=busybox --rm -it --restart=Never -- \
  wget -qO- http://my-app-service:80
```

---

## NodePort

Exposes the service on a port on every node (`30000–32767`). Accessible from outside via `<NodeIP>:<NodePort>`.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app-nodeport
spec:
  type: NodePort
  selector:
    app: my-app
  ports:
    - port: 80
      targetPort: 8080
      nodePort: 30080   # Optional — if omitted, K8s assigns one
```

```bash
# Minikube: get the URL
minikube service my-app-nodeport --url

# Manual access
kubectl get nodes -o wide   # Get node external IP
curl http://<node-ip>:30080
```

---

## LoadBalancer

Creates an external load balancer in the cloud provider (AWS ELB, GCP LB, Azure LB) and assigns a public IP.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app-lb
spec:
  type: LoadBalancer
  selector:
    app: my-app
  ports:
    - port: 80
      targetPort: 8080
```

```bash
kubectl get services
# NAME         TYPE           CLUSTER-IP     EXTERNAL-IP      PORT(S)
# my-app-lb    LoadBalancer   10.96.45.200   34.125.78.90     80:31234/TCP

curl http://34.125.78.90
```

> On Minikube: run `minikube tunnel` in a separate terminal to simulate LoadBalancer.

---

## ExternalName

Maps a service to an external DNS name (not pods).

```yaml
apiVersion: v1
kind: Service
metadata:
  name: external-db
spec:
  type: ExternalName
  externalName: my-database.rds.amazonaws.com
```

Pods can now connect to `external-db` and K8s resolves it to the RDS hostname.

---

## Headless Services (StatefulSets)

Set `clusterIP: None` to get direct DNS to each pod (used by StatefulSets).

```yaml
spec:
  clusterIP: None
  selector:
    app: mysql
```

Each pod gets its own DNS: `mysql-0.mysql-service.default.svc.cluster.local`

---

## DNS in Kubernetes

Kubernetes runs CoreDNS to provide automatic DNS resolution:

```
<service-name>.<namespace>.svc.cluster.local

# Examples:
my-app-service.default.svc.cluster.local
mysql-service.production.svc.cluster.local
```

From within the same namespace, you can use just the service name:
```bash
curl http://my-app-service/
```

From a different namespace:
```bash
curl http://my-app-service.default.svc.cluster.local/
```

---

## Service Discovery

```bash
# Inside a pod: env vars are auto-injected
env | grep MY_APP_SERVICE
# MY_APP_SERVICE_SERVICE_HOST=10.96.45.200
# MY_APP_SERVICE_SERVICE_PORT=80

# DNS lookup from inside a pod
kubectl exec -it my-pod -- nslookup my-app-service
```

---

## Network Policies

**NetworkPolicy** acts as a firewall between pods. By default, all pods can talk to each other.

### Deny All Ingress to a Namespace

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all
  namespace: production
spec:
  podSelector: {}          # Applies to all pods
  policyTypes:
    - Ingress
  ingress: []              # No ingress allowed
```

### Allow Only Frontend to Backend

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-frontend-to-backend
spec:
  podSelector:
    matchLabels:
      app: backend          # Applies to backend pods
  policyTypes:
    - Ingress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: frontend  # Only allow from frontend pods
      ports:
        - protocol: TCP
          port: 8080
```

### Allow Specific Namespace

```yaml
ingress:
  - from:
      - namespaceSelector:
          matchLabels:
            name: monitoring  # Allow from monitoring namespace
```

---

## Port Forwarding (Dev/Debug)

```bash
# Forward service port to localhost
kubectl port-forward svc/my-app-service 8080:80

# Forward pod port directly
kubectl port-forward pod/my-app-pod-abc123 8080:8080

# Forward in background
kubectl port-forward svc/my-app-service 8080:80 &
```

---

## Multi-Port Services

```yaml
spec:
  ports:
    - name: http
      port: 80
      targetPort: 8080
    - name: https
      port: 443
      targetPort: 8443
    - name: metrics
      port: 9090
      targetPort: 9090
```

---

## Service Networking Cheatsheet

```bash
# List all services
kubectl get services
kubectl get svc -A           # All namespaces

# Describe a service
kubectl describe svc my-app-service

# Check endpoints (which pod IPs are behind the service)
kubectl get endpoints my-app-service

# Delete service
kubectl delete svc my-app-service

# Expose a deployment quickly
kubectl expose deployment my-app --port=80 --target-port=8080 --type=ClusterIP
```
