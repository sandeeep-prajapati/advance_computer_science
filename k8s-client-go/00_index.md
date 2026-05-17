# Kubernetes Go Client (client-go) — Developer Book

A comprehensive, chapter-wise guide to interacting with the Kubernetes API using Go's official `client-go` library — from first connection to production-grade controllers, informers, and custom operators.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_introduction.md) | **Introduction** | What is client-go, architecture, when to use it, Go vs other clients |
| [02](02_installation_and_setup.md) | **Installation & Setup** | go mod, client-go packages, kubeconfig, code generation |
| [03](03_connecting_to_cluster.md) | **Connecting to a Cluster** | InClusterConfig, kubeconfig, multi-cluster, REST config |
| [04](04_pods.md) | **Working with Pods** | List, get, create, delete, logs, exec, watch pods |
| [05](05_deployments_and_replicasets.md) | **Deployments & ReplicaSets** | Create, scale, rolling update, rollback, status |
| [06](06_services_and_networking.md) | **Services & Networking** | ClusterIP, NodePort, LoadBalancer, endpoints |
| [07](07_configmaps_and_secrets.md) | **ConfigMaps & Secrets** | Create, read, update, delete config and secrets |
| [08](08_persistent_volumes_and_storage.md) | **Persistent Volumes & Storage** | PV, PVC, StorageClass management |
| [09](09_namespaces_and_rbac.md) | **Namespaces & RBAC** | Namespaces, Roles, ClusterRoles, Bindings, ServiceAccounts |
| [10](10_ingress.md) | **Ingress** | Create and manage Ingress rules, TLS, annotations |
| [11](11_statefulsets_and_daemonsets.md) | **StatefulSets & DaemonSets** | Stateful workloads, node-level agents |
| [12](12_jobs_and_cronjobs.md) | **Jobs & CronJobs** | One-off tasks, scheduled jobs, parallel batch |
| [13](13_informers_and_watching.md) | **Informers & Watching** | SharedInformer, event handlers, workqueues, cache |
| [14](14_custom_resources.md) | **Custom Resources (CRDs)** | Dynamic client, typed CRD client, code-gen |
| [15](15_error_handling.md) | **Error Handling & Retry Logic** | k8s errors, retry, backoff, conflict resolution |
| [16](16_production_patterns.md) | **Production Patterns** | Controller pattern, leader election, operator, multi-cluster |

---

## How to Read This Book

**If you're brand new to client-go:** Read 01 → 02 → 03 → 04 in order.

**If you're building a deployment tool:** Read 05 → 06 → 07 → 10.

**If you're writing a Kubernetes controller/operator:** Read 13 → 14 → 15 → 16.

**If you're managing RBAC programmatically:** Read 09 → 15.

**If something broke:** Jump to Chapter 15 (Error Handling).

---

## Quick Reference

```go
import (
    "context"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Build client from kubeconfig
config, _ := clientcmd.BuildConfigFromFlags("", "/home/user/.kube/config")
clientset, _ := kubernetes.NewForConfig(config)

// List pods
pods, _ := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
for _, p := range pods.Items {
    fmt.Println(p.Name)
}
```

---

## Key Packages

| Package | Import Path | Purpose |
|---|---|---|
| `kubernetes` | `k8s.io/client-go/kubernetes` | Main typed clientset |
| `clientcmd` | `k8s.io/client-go/tools/clientcmd` | Load kubeconfig |
| `rest` | `k8s.io/client-go/rest` | REST config, in-cluster |
| `cache` | `k8s.io/client-go/tools/cache` | Informers, listers |
| `workqueue` | `k8s.io/client-go/util/workqueue` | Rate-limiting queue |
| `retry` | `k8s.io/client-go/util/retry` | RetryOnConflict |
| `dynamic` | `k8s.io/client-go/dynamic` | Unstructured / CRD client |
| `corev1` | `k8s.io/api/core/v1` | Pod, Service, ConfigMap types |
| `appsv1` | `k8s.io/api/apps/v1` | Deployment, StatefulSet types |
| `metav1` | `k8s.io/apimachinery/pkg/apis/meta/v1` | ObjectMeta, ListOptions |
