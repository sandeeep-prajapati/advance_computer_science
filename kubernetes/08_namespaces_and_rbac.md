# Chapter 8: Namespaces & RBAC

## Namespaces

### What is a Namespace?

A **Namespace** is a logical partition within a Kubernetes cluster. It lets you:
- Isolate teams (frontend team, backend team)
- Separate environments (dev, staging, production)
- Apply resource quotas and limits per team
- Scope RBAC permissions

> Resources within a namespace can communicate. Cross-namespace requires full DNS: `service.namespace.svc.cluster.local`

### Built-in Namespaces

| Namespace | Purpose |
|---|---|
| `default` | Where resources go if you don't specify a namespace |
| `kube-system` | Kubernetes internal components (CoreDNS, kube-proxy) |
| `kube-public` | Publicly readable data |
| `kube-node-lease` | Node heartbeat leases |

### Creating and Using Namespaces

```bash
# Create a namespace
kubectl create namespace production
kubectl create namespace staging
kubectl create namespace development

# View namespaces
kubectl get namespaces

# Run commands in a specific namespace
kubectl get pods -n production
kubectl apply -f deployment.yaml -n production

# Set default namespace for your session
kubectl config set-context --current --namespace=production

# View all resources across all namespaces
kubectl get pods --all-namespaces
kubectl get pods -A
```

**Namespace YAML:**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    env: production
    team: platform
```

---

## Resource Quotas

Limit the total resources a namespace can use:

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: production-quota
  namespace: production
spec:
  hard:
    requests.cpu: "10"            # Total CPU requests limit
    requests.memory: "20Gi"       # Total memory requests limit
    limits.cpu: "20"
    limits.memory: "40Gi"
    pods: "50"                    # Max 50 pods
    services: "20"
    persistentvolumeclaims: "10"
    secrets: "30"
    configmaps: "30"
```

```bash
kubectl apply -f quota.yaml
kubectl describe resourcequota production-quota -n production
```

---

## LimitRange

Set default resource requests/limits for containers in a namespace:

```yaml
apiVersion: v1
kind: LimitRange
metadata:
  name: default-limits
  namespace: development
spec:
  limits:
    - type: Container
      default:               # Default limits (if not set in pod)
        cpu: "500m"
        memory: "256Mi"
      defaultRequest:        # Default requests (if not set in pod)
        cpu: "100m"
        memory: "128Mi"
      max:                   # Maximum allowed
        cpu: "2"
        memory: "2Gi"
      min:                   # Minimum allowed
        cpu: "50m"
        memory: "64Mi"
```

---

## RBAC (Role-Based Access Control)

RBAC controls **who can do what** in your cluster.

### RBAC Objects

| Object | Scope | Purpose |
|---|---|---|
| `Role` | Namespace | Permissions within one namespace |
| `ClusterRole` | Cluster-wide | Permissions across all namespaces or for cluster-level resources |
| `RoleBinding` | Namespace | Assigns a Role to a user/group/ServiceAccount |
| `ClusterRoleBinding` | Cluster-wide | Assigns a ClusterRole to a user/group/ServiceAccount |

### API Verbs (Permissions)

| Verb | Action |
|---|---|
| `get` | Read a single resource |
| `list` | List resources |
| `watch` | Watch for changes |
| `create` | Create resources |
| `update` | Update resources |
| `patch` | Partially update resources |
| `delete` | Delete resources |
| `deletecollection` | Delete all matching resources |

---

## Role (Namespace-scoped)

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pod-reader
  namespace: production
rules:
  - apiGroups: [""]              # "" = core API group
    resources: ["pods", "pods/log"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get", "list"]
```

---

## ClusterRole (Cluster-scoped)

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: node-reader
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["namespaces"]
    verbs: ["get", "list"]
```

---

## RoleBinding

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: alice-pod-reader
  namespace: production
subjects:
  - kind: User
    name: alice                     # Kubernetes user
    apiGroup: rbac.authorization.k8s.io
  - kind: Group
    name: developers
    apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```

---

## ServiceAccount

A **ServiceAccount** is an identity for Pods (not humans). Use it to give Pods permission to call the Kubernetes API.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-app-sa
  namespace: production
```

```yaml
# Bind ServiceAccount to a Role
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: my-app-sa-binding
  namespace: production
subjects:
  - kind: ServiceAccount
    name: my-app-sa
    namespace: production
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```

```yaml
# Use ServiceAccount in a Deployment
spec:
  template:
    spec:
      serviceAccountName: my-app-sa
      containers:
        - name: app
          image: my-app:latest
```

---

## Common RBAC Patterns

### Read-Only Access to a Namespace

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: namespace-viewer
  namespace: staging
rules:
  - apiGroups: ["", "apps", "batch"]
    resources: ["*"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: developers-viewer
  namespace: staging
subjects:
  - kind: Group
    name: developers
roleRef:
  kind: Role
  name: namespace-viewer
  apiGroup: rbac.authorization.k8s.io
```

### Deploy-Only Role (CI/CD)

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: deployer
  namespace: production
rules:
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get", "list", "update", "patch"]
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["get", "list", "create", "update"]
```

---

## RBAC Commands

```bash
# Check your permissions
kubectl auth can-i create pods
kubectl auth can-i get secrets -n production
kubectl auth can-i "*" "*"   # Am I cluster-admin?

# Check another user's permissions
kubectl auth can-i get pods --as alice -n production
kubectl auth can-i get secrets --as system:serviceaccount:production:my-app-sa

# List roles
kubectl get roles -A
kubectl get clusterroles | grep -v system:

# List bindings
kubectl get rolebindings -A
kubectl get clusterrolebindings

# Describe a role
kubectl describe role pod-reader -n production
```

---

## Best Practices

- **Least privilege:** Grant only the minimum required permissions
- **Avoid `cluster-admin` bindings** for application service accounts
- **Use Namespaces** to scope RBAC and resource quotas per team/environment
- **Audit RBAC regularly** using `kubectl auth can-i` and `kubectl get rolebindings`
- **Don't use the `default` ServiceAccount** for application pods — create dedicated ones
- **Separate concerns:** One ServiceAccount per application/service
