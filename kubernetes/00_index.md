# Kubernetes Developer Book

A comprehensive, chapter-wise guide to mastering Kubernetes — from first cluster to production-grade deployments, autoscaling, security, and cloud providers.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_introduction.md) | **Introduction to Kubernetes** | What is K8s, why it exists, key concepts, architecture, control plane vs worker nodes |
| [02](02_setup_and_installation.md) | **Setup & Installation** | Minikube, k3s, kind, kubectl, Docker Desktop, cloud cluster setup |
| [03](03_pods_and_nodes.md) | **Pods & Nodes** | What is a Pod, multi-container pods, lifecycle, node scheduling, resource requests |
| [04](04_deployments_and_replicasets.md) | **Deployments & ReplicaSets** | Creating deployments, rolling updates, rollback, scaling, strategies |
| [05](05_services_and_networking.md) | **Services & Networking** | ClusterIP, NodePort, LoadBalancer, DNS, service discovery, network policies |
| [06](06_configmaps_and_secrets.md) | **ConfigMaps & Secrets** | Externalizing config, mounting as env vars or files, secret management |
| [07](07_persistent_volumes_and_storage.md) | **Persistent Volumes & Storage** | PV, PVC, StorageClass, dynamic provisioning, ReadWriteMany, StatefulSets storage |
| [08](08_namespaces_and_rbac.md) | **Namespaces & RBAC** | Organizing resources, roles, ClusterRoles, service accounts, least-privilege |
| [09](09_ingress_and_load_balancing.md) | **Ingress & Load Balancing** | Ingress controllers, Nginx ingress, TLS/cert-manager, path-based routing |
| [10](10_helm.md) | **Helm — Package Manager** | Charts, values.yaml, install/upgrade/rollback, writing your own chart |
| [11](11_statefulsets_and_daemonsets.md) | **StatefulSets & DaemonSets** | When to use each, stable identity, ordered updates, node-level agents |
| [12](12_jobs_and_cronjobs.md) | **Jobs & CronJobs** | One-off tasks, scheduled tasks, parallelism, failure handling |
| [13](13_resource_management_and_autoscaling.md) | **Resource Management & Autoscaling** | Requests/limits, LimitRange, ResourceQuota, HPA, VPA, KEDA |
| [14](14_monitoring_and_logging.md) | **Monitoring & Logging** | Prometheus, Grafana, Loki, kubectl logs, metrics-server, alerting |
| [15](15_security_and_best_practices.md) | **Security & Best Practices** | SecurityContext, PodSecurity, NetworkPolicy, image scanning, secrets management |
| [16](16_cicd_with_kubernetes.md) | **CI/CD with Kubernetes** | GitHub Actions, ArgoCD, GitOps pattern, rolling deployments, blue-green |
| [17](17_stateful_workloads.md) | **Stateful Workloads** | Databases on K8s, operators, MySQL/Redis/Elasticsearch deployments |
| [18](18_debugging_and_troubleshooting.md) | **Debugging & Troubleshooting** | kubectl debug commands, CrashLoopBackOff, OOMKilled, ImagePullBackOff, events |
| [19](19_cloud_providers.md) | **Cloud Providers** | GKE, EKS, AKS — managed clusters, node pools, cloud-native storage, IAM |
| [20](20_learning_roadmap.md) | **Learning Roadmap** | 5-week hands-on plan from zero to production-ready |

---

## How to Read This Book

**If you're brand new to Kubernetes:** Read 1 → 2 → 3 → 4 → 5 in order. This gives you the core mental model.

**If you're deploying an application:** Read 4 → 5 → 6 → 7 → 9.

**If you're securing a cluster:** Read 8 → 15 → 6 → 13.

**If you're going to production:** Read 13 → 14 → 15 → 16 → 19.

**If you're running stateful services:** Read 7 → 11 → 17.

**If something broke:** Jump to Chapter 18 (Debugging & Troubleshooting).

---

## Key Concepts at a Glance

```
Kubernetes Cluster
├── Control Plane
│   ├── API Server          (all kubectl commands go here)
│   ├── etcd                (cluster state database)
│   ├── Scheduler           (decides which node runs a pod)
│   └── Controller Manager  (enforces desired state)
│
├── Worker Nodes
│   ├── kubelet             (runs pods on the node)
│   ├── kube-proxy          (network rules)
│   └── Container Runtime   (Docker/containerd)
│
├── Workloads
│   ├── Pod                 (smallest deployable unit)
│   ├── Deployment          (stateless apps, rolling updates)
│   ├── StatefulSet         (stateful apps, stable identity)
│   ├── DaemonSet           (one pod per node)
│   ├── Job / CronJob       (one-off or scheduled tasks)
│   └── ReplicaSet          (desired pod count enforcement)
│
├── Networking
│   ├── Service             (stable endpoint for pods)
│   ├── Ingress             (HTTP routing, TLS)
│   └── NetworkPolicy       (firewall rules between pods)
│
├── Config & Storage
│   ├── ConfigMap           (non-sensitive config)
│   ├── Secret              (sensitive data)
│   ├── PersistentVolume    (storage resource)
│   └── PersistentVolumeClaim (storage request)
│
└── Access Control
    ├── Namespace           (logical isolation)
    ├── ServiceAccount      (pod identity)
    ├── Role / ClusterRole  (permission definitions)
    └── RoleBinding         (bind role to account)
```
