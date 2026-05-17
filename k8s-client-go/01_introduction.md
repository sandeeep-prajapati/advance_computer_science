# Chapter 1: Introduction to client-go

## What is client-go?

**client-go** is the official Go client library for the Kubernetes API, maintained by the Kubernetes project itself. It is the foundation for:
- `kubectl` (the Kubernetes CLI)
- `kube-controller-manager` (built-in controllers)
- Every major Kubernetes operator and controller

```
client-go is to Kubernetes what the AWS SDK is to AWS.
```

---

## Why client-go Over Other Options?

| Option | Pros | Cons |
|---|---|---|
| `client-go` | Official, full-featured, informers, workqueue | Verbose, steep learning curve |
| `controller-runtime` | Higher-level, operator-friendly | Abstraction hides internals |
| Raw HTTP | Total control | No type safety, lots of boilerplate |
| `kubectl` CLI | Simple | No programmatic control |

> **Rule:** Use `client-go` directly when building low-level controllers or need full control. Use `controller-runtime` (via kubebuilder) for production operators.

---

## client-go Architecture

```
Your Go Program
     │
     ├── kubernetes.Clientset       ← typed API access
     │       ├── CoreV1()           ← Pods, Services, ConfigMaps
     │       ├── AppsV1()           ← Deployments, StatefulSets
     │       ├── BatchV1()          ← Jobs, CronJobs
     │       ├── NetworkingV1()     ← Ingress, NetworkPolicy
     │       └── RbacV1()           ← Roles, Bindings
     │
     ├── dynamic.Interface          ← unstructured / CRD access
     │
     ├── cache.SharedInformer       ← watch + local cache
     │       └── workqueue          ← rate-limited event queue
     │
     └── rest.Config                ← connection & auth
```

---

## API Groups and Versions in Go

| API Group | Go Client Method | Resources |
|---|---|---|
| Core (`v1`) | `clientset.CoreV1()` | Pod, Service, ConfigMap, Secret, PV, PVC, Namespace, Node |
| Apps (`apps/v1`) | `clientset.AppsV1()` | Deployment, StatefulSet, DaemonSet, ReplicaSet |
| Batch (`batch/v1`) | `clientset.BatchV1()` | Job, CronJob |
| Networking | `clientset.NetworkingV1()` | Ingress, NetworkPolicy |
| RBAC | `clientset.RbacV1()` | Role, ClusterRole, RoleBinding, ClusterRoleBinding |
| Storage | `clientset.StorageV1()` | StorageClass, VolumeAttachment |

---

## The Desired State Model (Reconciliation)

client-go is built around the Kubernetes control-loop model:

```
┌─────────────────────────────────────────────┐
│              Reconciliation Loop             │
│                                              │
│  1. Observe  →  Watch API for events         │
│  2. Diff     →  Desired state vs actual      │
│  3. Act      →  Create / Update / Delete     │
│  4. Repeat                                   │
└─────────────────────────────────────────────┘
```

Your controller:
1. Uses an **Informer** to watch resources (cached locally)
2. Puts events into a **WorkQueue**
3. A **worker goroutine** dequeues and calls your **Reconcile** function
4. Reconcile compares desired vs actual and acts

---

## What You Can Do with client-go

- **CRUD on any resource** — Pods, Deployments, Services, Secrets, CRDs
- **Watch resources** — real-time event streaming with Informers
- **Leader election** — only one controller instance acts at a time
- **Custom controllers** — reconcile any resource type
- **Exec into pods** — run commands programmatically
- **Port forwarding** — forward cluster ports to local
- **Certificate operations** — create CSRs, approve certificates

---

## Go vs Other Language Clients

| Feature | Go (client-go) | Python | PHP |
|---|---|---|---|
| Official support | Yes (core team) | Yes | Community |
| Informer / cache | Full support | Partial | No |
| Operator frameworks | controller-runtime, kubebuilder | kopf | None |
| Concurrency | Goroutines (excellent) | asyncio | Limited |
| Type safety | Full generated types | Yes | No |
| Performance | Best | Good | Good |
| Use in production operators | Standard | Used | Rare |

---

## Real-World Use Cases

- **kubectl itself** — built on client-go
- **cert-manager** — manages TLS certificates as a controller
- **ArgoCD** — GitOps controller built with client-go
- **Prometheus Operator** — manages Prometheus instances
- **Custom CI/CD** — deploy apps, watch rollout status
- **Multi-tenant platforms** — provision namespaces, RBAC, quotas per customer
- **Cluster maintenance tools** — drain nodes, evict pods, rotate secrets
