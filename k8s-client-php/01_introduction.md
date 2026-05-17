# Chapter 1: Introduction to k8s-client-php

## What is k8s-client-php?

**k8s-client-php** refers to PHP libraries that wrap the Kubernetes REST API, letting you manage cluster resources (Pods, Deployments, Services, etc.) directly from PHP code — without shelling out to `kubectl`.

The two most popular options are:

| Library | Style | Best For |
|---|---|---|
| `maclof/kubernetes-client` | Procedural, fluent | General PHP apps |
| `renoki-co/php-k8s` | Eloquent-like, Laravel | Laravel/Lumen projects |

This book uses `maclof/kubernetes-client` as the primary example, with `php-k8s` patterns shown where relevant.

---

## Why Use a PHP Client Instead of kubectl?

| Scenario | Use kubectl | Use PHP Client |
|---|---|---|
| One-off cluster debugging | Yes | No |
| CI/CD pipeline shell script | Yes | Maybe |
| Web dashboard managing resources | No | Yes |
| SaaS platform provisioning namespaces per customer | No | Yes |
| Kubernetes Operator written in PHP | No | Yes |
| Automated reconciliation loop | No | Yes |

> **Rule:** If a human is running the command → use `kubectl`. If your application is doing it programmatically → use a PHP client.

---

## How the Kubernetes API Works

The Kubernetes API is a **RESTful HTTP API** with a well-defined schema:

```
GET    /api/v1/namespaces/{namespace}/pods          → list pods
GET    /api/v1/namespaces/{namespace}/pods/{name}   → get pod
POST   /api/v1/namespaces/{namespace}/pods          → create pod
PATCH  /api/v1/namespaces/{namespace}/pods/{name}   → update pod
DELETE /api/v1/namespaces/{namespace}/pods/{name}   → delete pod
```

All responses are JSON. The PHP client libraries are thin wrappers that:
1. Build the correct URL and HTTP method
2. Serialize PHP arrays/objects to JSON
3. Handle authentication headers
4. Parse responses back into PHP objects

---

## API Groups and Versions

Kubernetes has multiple API groups:

| API Group | Path | Contains |
|---|---|---|
| Core | `/api/v1` | Pods, Services, ConfigMaps, Secrets, PVs |
| Apps | `/apis/apps/v1` | Deployments, StatefulSets, DaemonSets, ReplicaSets |
| Batch | `/apis/batch/v1` | Jobs, CronJobs |
| Networking | `/apis/networking.k8s.io/v1` | Ingress, NetworkPolicy |
| RBAC | `/apis/rbac.authorization.k8s.io/v1` | Roles, ClusterRoles, Bindings |
| Custom | `/apis/{group}/{version}` | CRDs |

The PHP client handles routing to the correct API group automatically.

---

## Architecture Overview

```
Your PHP App
     │
     │  HTTP/TLS
     ▼
Kubernetes API Server  (port 6443)
     │
     ├── etcd          (stores all state)
     ├── Scheduler
     └── Controllers
```

Your PHP code:
1. Authenticates with the API Server (cert, token, or kubeconfig)
2. Sends REST requests to create/read/update/delete resources
3. The API Server validates and stores the changes
4. Controllers reconcile toward the new desired state

---

## What You Can Do with the PHP Client

- **Read resources:** List pods, get deployment status, check service endpoints
- **Write resources:** Create namespaces, deploy applications, create secrets
- **Watch resources:** Stream real-time events (pod restarts, deployment rollouts)
- **Execute into pods:** Run commands inside running containers
- **Manage RBAC:** Create roles, bindings, service accounts programmatically
- **Handle CRDs:** Interact with custom resource types

---

## What You Cannot Do (Limitations)

- Cannot replace the full power of `kubectl` for interactive debugging
- No built-in leader election (you implement it yourself)
- CRD support varies by library — raw HTTP is sometimes needed
- Watching streams (long-polling) requires careful connection management in PHP-FPM environments

---

## PHP vs Other Language Clients

| Feature | PHP Client | Go client-go | Python kubernetes |
|---|---|---|---|
| Official support | Community | Official (CNCF) | Official (CNCF) |
| Operator frameworks | None | controller-runtime | kopf |
| Stability | Good | Excellent | Excellent |
| Ecosystem | Laravel-friendly | Best for operators | Strong ML/data |
| Async support | Limited | Yes (goroutines) | Yes (asyncio) |

> For building production operators or controllers, Go is the standard. PHP clients are ideal for web dashboards, SaaS provisioning, and CI/CD integrations within existing PHP stacks.
