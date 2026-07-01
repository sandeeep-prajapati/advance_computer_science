# Chapter 5: Kubernetes Architecture

## Big Picture

```text
kubectl
  |
  v
API Server
  |
  +--> etcd
  +--> Scheduler
  +--> Controller Manager
  +--> Kubelet on each node
  +--> Kube Proxy / networking
```

Kubernetes is API-driven. Most components watch the API server and react to changes.

---

## Control Plane Components

| Component | Job |
|---|---|
| API server | Validates and serves Kubernetes API requests |
| etcd | Stores cluster state |
| scheduler | Chooses which node should run a Pod |
| controller-manager | Runs controllers that reconcile desired state |

---

## Node Components

| Component | Job |
|---|---|
| kubelet | Ensures containers run on a node |
| container runtime | Runs containers, commonly containerd |
| kube-proxy | Handles Service networking rules |

---

## Controllers

A controller watches resources and tries to make reality match the desired state.

Example:

```text
Desired state:
Deployment says replicas=3

Actual state:
Only 2 Pods exist

Controller action:
Create 1 more Pod
```

This pattern appears everywhere in Kubernetes.

---

## API Objects

Most Kubernetes features are API objects:

- Pod
- Deployment
- Service
- ConfigMap
- Secret
- Node
- PersistentVolume
- Job

The API server validates and stores these objects. Controllers and kubelets act on them.

---

## Contributor Mental Map

When reading a bug or code path, classify it:

| Area | Likely Component |
|---|---|
| API validation/defaulting | API server / API machinery |
| Pod placement | scheduler |
| Pod lifecycle on node | kubelet |
| Replica management | controller-manager |
| Service networking | kube-proxy / networking |
| CLI behavior | kubectl / staging modules |
| Tests failing | test framework or component behavior |

