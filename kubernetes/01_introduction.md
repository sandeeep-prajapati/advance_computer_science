# Chapter 1: Introduction to Kubernetes

## What is Kubernetes?

**Kubernetes** (K8s) is an open-source **container orchestration platform** originally developed by Google and now maintained by the Cloud Native Computing Foundation (CNCF). It automates the deployment, scaling, and management of containerized applications.

> **Analogy:** Think of Kubernetes as an **intelligent operations team** for your containers:
> - You tell it *what* you want running (desired state).
> - It figures out *how* to make it happen and keeps it that way — even when servers fail.

---

## Why Kubernetes Exists

Before Kubernetes, deploying containers at scale meant manually managing:
- Which server runs which container?
- What happens when a container crashes?
- How do you scale up during traffic spikes?
- How do you deploy new versions without downtime?

Kubernetes solves all of this automatically.

| Problem | Kubernetes Solution |
|---|---|
| Container crashes | Automatically restarts it |
| Server dies | Reschedules containers on healthy nodes |
| Traffic spike | Scales pods up automatically (HPA) |
| New deployment | Rolling update with zero downtime |
| Multiple environments | Namespaces for dev/staging/prod isolation |

---

## Core Concepts

### Cluster
A **cluster** is the entire Kubernetes environment — a set of machines (nodes) that run your workloads.

```
Kubernetes Cluster
├── Control Plane (brain)
└── Worker Nodes  (muscle)
```

### Node
A **node** is a physical or virtual machine that runs your containers. There are two types:
- **Control Plane Node** — manages the cluster (API server, scheduler, etcd)
- **Worker Node** — runs your application containers

### Pod
A **Pod** is the smallest deployable unit in Kubernetes. It wraps one or more containers that share:
- The same network namespace (same IP)
- The same storage volumes

> Pods are **ephemeral** — they can be killed and replaced at any time.

### Deployment
A **Deployment** manages a set of identical Pods. You tell it "I want 3 replicas of my app" and it ensures 3 Pods are always running.

### Service
A **Service** gives your Pods a stable network address. Since Pod IPs change when Pods restart, you always access them through a Service.

### Namespace
A **Namespace** is a logical partition of the cluster for isolating teams, environments, or projects.

---

## Kubernetes Architecture

```
┌─────────────────────────────────────────────────────┐
│                  CONTROL PLANE                       │
│                                                      │
│  ┌──────────────┐  ┌──────┐  ┌────────────────────┐ │
│  │  API Server  │  │ etcd │  │ Controller Manager  │ │
│  │  (kube-      │  │(state│  │ (reconciliation     │ │
│  │   apiserver) │  │  DB) │  │  loop)              │ │
│  └──────────────┘  └──────┘  └────────────────────┘ │
│         │                    ┌──────────────────────┐│
│         │                    │      Scheduler       ││
│         │                    │  (pod placement)     ││
│         │                    └──────────────────────┘│
└─────────┼───────────────────────────────────────────┘
          │ kubectl / API calls
          │
┌─────────▼──────────────────────────────────────────┐
│                  WORKER NODES                        │
│                                                      │
│  Node 1                Node 2              Node 3   │
│  ┌──────────┐          ┌──────────┐        ┌──────┐  │
│  │ kubelet  │          │ kubelet  │        │kubel.│  │
│  │kube-proxy│          │kube-proxy│        │k-p   │  │
│  │[Pod][Pod]│          │[Pod][Pod]│        │[Pod] │  │
│  └──────────┘          └──────────┘        └──────┘  │
└──────────────────────────────────────────────────────┘
```

### Control Plane Components

| Component | Role |
|---|---|
| **API Server** | The front door — all `kubectl` commands go here |
| **etcd** | Key-value store holding the entire cluster state |
| **Scheduler** | Decides which node each new Pod runs on |
| **Controller Manager** | Watches cluster state and reconciles toward desired state |
| **Cloud Controller** | Integrates with cloud providers (AWS, GCP, Azure) |

### Worker Node Components

| Component | Role |
|---|---|
| **kubelet** | Agent on each node — ensures containers run as specified |
| **kube-proxy** | Manages network rules for Service traffic routing |
| **Container Runtime** | Runs containers (containerd, Docker, CRI-O) |

---

## The Desired State Model

This is the most important concept in Kubernetes:

1. **You declare** what you want: "Run 3 replicas of my app"
2. **Kubernetes stores** this desired state in etcd
3. **Controllers watch** the actual state continuously
4. **When drift happens** (a pod crashes), controllers reconcile — they create a new pod to match the desired state

```
You declare:        replicas: 3
Current state:      replicas: 2  (one crashed)
Controller action:  Create 1 new pod  →  replicas: 3 ✓
```

---

## Kubernetes vs Docker Compose vs Traditional VMs

| Feature | Traditional VMs | Docker Compose | Kubernetes |
|---|---|---|---|
| Scaling | Manual | Manual | Automatic (HPA) |
| Self-healing | No | No | Yes |
| Rolling updates | Manual/downtime | Manual | Automatic |
| Load balancing | External setup | Manual | Built-in Services |
| Multi-host | Complex | Limited | Native |
| Production-ready | Yes (ops-heavy) | No | Yes |

---

## Key Kubernetes Objects Summary

| Object | What It Does |
|---|---|
| `Pod` | Runs containers |
| `Deployment` | Manages stateless pod replicas |
| `StatefulSet` | Manages stateful pods with stable identity |
| `DaemonSet` | Runs one pod per node |
| `Job` | Runs a task to completion |
| `CronJob` | Scheduled tasks |
| `Service` | Stable network endpoint for pods |
| `Ingress` | HTTP routing and TLS termination |
| `ConfigMap` | Non-sensitive configuration |
| `Secret` | Sensitive data (passwords, tokens) |
| `PersistentVolumeClaim` | Request for storage |
| `Namespace` | Logical isolation |
| `ServiceAccount` | Identity for pods |
| `Role` / `ClusterRole` | Permission definitions |

---

## Real-World Use Cases

- **Web applications:** Run multiple replicas, auto-scale under load, zero-downtime deploys
- **Microservices:** Each service in its own Deployment, communicating via Services
- **Batch processing:** Jobs for data pipelines, CronJobs for scheduled reports
- **Machine learning:** GPU node pools for training, CPU pools for serving
- **Multi-tenant SaaS:** Namespaces + RBAC to isolate customers
- **CI/CD pipelines:** Ephemeral pods for test runners
