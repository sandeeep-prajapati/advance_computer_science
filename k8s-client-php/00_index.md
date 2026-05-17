# Kubernetes PHP Client (k8s-client-php) — Developer Book

A comprehensive, chapter-wise guide to interacting with the Kubernetes API using PHP — from first connection to production-grade cluster management, event watching, RBAC, and real-world automation.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_introduction.md) | **Introduction** | What is k8s-client-php, when to use it, architecture overview, PHP vs kubectl |
| [02](02_installation_and_setup.md) | **Installation & Setup** | Composer install, authentication methods, KubeConfig, in-cluster config |
| [03](03_connecting_to_cluster.md) | **Connecting to a Cluster** | Local cluster, remote cluster, service accounts, TLS, context switching |
| [04](04_pods.md) | **Working with Pods** | List, get, create, delete pods, exec into pods, logs, pod status |
| [05](05_deployments_and_replicasets.md) | **Deployments & ReplicaSets** | Create/update deployments, rolling updates, scale, rollback |
| [06](06_services_and_networking.md) | **Services & Networking** | ClusterIP, NodePort, LoadBalancer services, endpoints |
| [07](07_configmaps_and_secrets.md) | **ConfigMaps & Secrets** | Create, read, update ConfigMaps and Secrets via PHP |
| [08](08_persistent_volumes_and_storage.md) | **Persistent Volumes & Storage** | PV, PVC, StorageClass management through the API |
| [09](09_namespaces_and_rbac.md) | **Namespaces & RBAC** | Create namespaces, Roles, ClusterRoles, RoleBindings, ServiceAccounts |
| [10](10_ingress.md) | **Ingress** | Create and manage Ingress rules, TLS, annotations |
| [11](11_statefulsets_and_daemonsets.md) | **StatefulSets & DaemonSets** | Managing stateful workloads and node-level agents |
| [12](12_jobs_and_cronjobs.md) | **Jobs & CronJobs** | One-off tasks and scheduled jobs via PHP |
| [13](13_watching_events.md) | **Watching Events & Real-Time Updates** | Watch API, streaming events, detecting pod restarts |
| [14](14_custom_resources.md) | **Custom Resources (CRDs)** | Listing and managing CRDs, dynamic client usage |
| [15](15_error_handling.md) | **Error Handling & Retry Logic** | API exceptions, status codes, retries, rate limiting |
| [16](16_production_patterns.md) | **Production Patterns** | Operator pattern, leader election, reconciliation loops, multi-cluster |

---

## How to Read This Book

**If you're brand new to k8s-client-php:** Read 01 → 02 → 03 → 04 in order. This gives you the core mental model.

**If you're building a deployment tool:** Read 05 → 06 → 07 → 10.

**If you're managing cluster resources programmatically:** Read 09 → 08 → 12 → 13.

**If you're writing a Kubernetes Operator in PHP:** Read 13 → 14 → 16.

**If something broke:** Jump to Chapter 15 (Error Handling & Retry Logic).

---

## Quick Reference

```php
use Maclof\Kubernetes\Client;

$client = new Client([
    'master'      => 'https://k8s.example.com:6443',
    'ca_cert'     => '/path/to/ca.crt',
    'client_cert' => '/path/to/client.crt',
    'client_key'  => '/path/to/client.key',
]);

// List all pods in default namespace
$pods = $client->pods()->find();

// Get a specific pod
$pod = $client->pods()->findByName('my-pod');

// Delete a pod
$client->pods()->deleteByName('my-pod');
```

---

## Key Libraries

| Library | GitHub | Notes |
|---|---|---|
| `maclof/kubernetes-client` | maclof/kubernetes-client | Most widely used PHP K8s client |
| `php-k8s/php-k8s` | renoki-co/php-k8s | Laravel-friendly, fluent API |
| Raw HTTP via Guzzle | — | For advanced/custom resources |
