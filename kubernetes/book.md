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


---

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


---

# Chapter 2: Setup & Installation

## What You Need

| Tool | Purpose |
|---|---|
| `kubectl` | CLI to interact with any Kubernetes cluster |
| `minikube` | Local single-node cluster (best for beginners) |
| `kind` | Kubernetes in Docker (fast, CI-friendly) |
| `k3s` | Lightweight K8s (great for Raspberry Pi, edge) |
| `Docker Desktop` | Includes K8s option on Mac/Windows |

---

## Step 1: Install kubectl

**kubectl** is the command-line tool for every Kubernetes operation.

### Linux

```bash
# Download latest stable version
curl -LO "https://dl.k8s.io/release/$(curl -Ls https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"

# Install
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Verify
kubectl version --client
```

### macOS (Homebrew)

```bash
brew install kubectl
kubectl version --client
```

### Windows

```powershell
# Using winget
winget install -e --id Kubernetes.kubectl
```

### Enable Shell Autocompletion

```bash
# Bash
echo 'source <(kubectl completion bash)' >> ~/.bashrc
source ~/.bashrc

# Zsh
echo 'source <(kubectl completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

---

## Option A: Minikube (Recommended for beginners)

Minikube runs a single-node Kubernetes cluster locally inside a VM or container.

### Install Minikube

```bash
# Linux
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# macOS
brew install minikube
```

### Start a Cluster

```bash
# Start with Docker driver (recommended)
minikube start --driver=docker

# Start with more resources
minikube start --driver=docker --cpus=4 --memory=8192

# Check status
minikube status

# View dashboard
minikube dashboard
```

### Minikube Useful Commands

```bash
# Stop cluster (preserves state)
minikube stop

# Delete cluster (clean slate)
minikube delete

# SSH into the node
minikube ssh

# Load a local Docker image into minikube
minikube image load my-app:latest

# Enable addons
minikube addons enable ingress
minikube addons enable metrics-server
minikube addons list
```

---

## Option B: kind (Kubernetes in Docker)

**kind** is fast, lightweight, and great for CI pipelines.

### Install kind

```bash
# Linux
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.23.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# macOS
brew install kind
```

### Create a Cluster

```bash
# Single node
kind create cluster --name my-cluster

# Multi-node (1 control plane + 2 workers)
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
  - role: worker
  - role: worker
EOF

# Delete cluster
kind delete cluster --name my-cluster
```

### Load Local Images into kind

```bash
# Build your image
docker build -t my-app:latest .

# Load into kind (required — kind doesn't see your local Docker daemon)
kind load docker-image my-app:latest --name my-cluster
```

---

## Option C: k3s (Lightweight, Production-like)

**k3s** is a minimal Kubernetes distribution (single ~70MB binary). Great for local Linux or Raspberry Pi.

```bash
# Install k3s (single command)
curl -sfL https://get.k3s.io | sh -

# Check status
sudo k3s kubectl get nodes

# Configure kubectl to use k3s
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $(id -u):$(id -g) ~/.kube/config
kubectl get nodes

# Uninstall
/usr/local/bin/k3s-uninstall.sh
```

---

## Verify Your Cluster

After any setup, run these to confirm everything works:

```bash
# Check cluster info
kubectl cluster-info

# List nodes (should show Ready)
kubectl get nodes

# List all pods in all namespaces
kubectl get pods --all-namespaces

# Check component health
kubectl get componentstatuses
```

---

## kubectl Configuration (kubeconfig)

kubectl uses `~/.kube/config` to know which cluster to talk to.

```bash
# View current context (which cluster you're targeting)
kubectl config current-context

# List all configured contexts
kubectl config get-contexts

# Switch to a different cluster
kubectl config use-context my-other-cluster

# View full kubeconfig
kubectl config view
```

---

## Deploy Your First App

```bash
# Deploy nginx
kubectl create deployment nginx --image=nginx:latest

# Check pods
kubectl get pods

# Expose as NodePort service
kubectl expose deployment nginx --port=80 --type=NodePort

# Get the URL (minikube)
minikube service nginx --url

# OR port-forward to localhost
kubectl port-forward deployment/nginx 8080:80
# Now open: http://localhost:8080
```

---

## Cleanup

```bash
kubectl delete deployment nginx
kubectl delete service nginx
```

---

## Quick Reference: Which Tool to Use?

| Scenario | Tool |
|---|---|
| Learning at home | Minikube |
| CI/CD pipelines | kind |
| Low-resource Linux / Raspberry Pi | k3s |
| Mac or Windows with GUI | Docker Desktop (enable K8s) |
| Team-shared dev cluster | Cloud: GKE/EKS/AKS free tier |
| Production | Managed cloud K8s (GKE, EKS, AKS) |


---

# Chapter 3: Pods & Nodes

## What is a Pod?

A **Pod** is the smallest deployable unit in Kubernetes. Every container you run lives inside a Pod. A Pod:
- Wraps one or more containers
- Gives them a **shared IP address**
- Gives them **shared storage volumes**
- Lives and dies together — all containers start and stop as a unit

> **Rule:** In 99% of cases, one Pod = one container. Multi-container Pods are for tightly-coupled helpers (sidecars).

---

## Pod Lifecycle

```
Pending → Running → Succeeded (or Failed)
              ↓
           (if crash)
           CrashLoopBackOff
```

| Phase | Meaning |
|---|---|
| `Pending` | Pod accepted, waiting for a node/image pull |
| `Running` | At least one container is running |
| `Succeeded` | All containers exited with code 0 (Jobs) |
| `Failed` | At least one container exited with non-zero |
| `Unknown` | Node communication lost |
| `CrashLoopBackOff` | Container keeps crashing and restarting |

---

## Writing Your First Pod YAML

**File:** `pod.yaml`

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-nginx
  namespace: default
  labels:
    app: nginx
    env: dev
spec:
  containers:
    - name: nginx
      image: nginx:1.25-alpine
      ports:
        - containerPort: 80
      resources:
        requests:
          memory: "64Mi"
          cpu: "100m"
        limits:
          memory: "128Mi"
          cpu: "200m"
```

```bash
# Apply the pod
kubectl apply -f pod.yaml

# Check status
kubectl get pods

# Describe the pod (see events, IP, node)
kubectl describe pod my-nginx

# Get logs
kubectl logs my-nginx

# Exec into the container
kubectl exec -it my-nginx -- sh

# Delete the pod
kubectl delete pod my-nginx
```

---

## Multi-Container Pods

Use multi-container pods when containers are tightly coupled. Common patterns:

### Sidecar Pattern

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-with-log-sidecar
spec:
  containers:
    # Main application
    - name: app
      image: my-app:latest
      volumeMounts:
        - name: shared-logs
          mountPath: /var/log/app

    # Sidecar: ships logs to external system
    - name: log-shipper
      image: fluent/fluent-bit:latest
      volumeMounts:
        - name: shared-logs
          mountPath: /var/log/app
          readOnly: true

  volumes:
    - name: shared-logs
      emptyDir: {}
```

### Init Container Pattern

Init containers run **before** the main container starts — useful for setup tasks (wait for DB, run migrations).

```yaml
spec:
  initContainers:
    - name: wait-for-db
      image: busybox
      command: ['sh', '-c', 'until nc -z mysql-service 3306; do echo waiting; sleep 2; done']

  containers:
    - name: app
      image: my-app:latest
```

---

## Resource Requests and Limits

```yaml
resources:
  requests:
    memory: "256Mi"   # Guaranteed minimum
    cpu: "250m"       # 250 millicores = 0.25 CPU
  limits:
    memory: "512Mi"   # Hard cap — OOMKill if exceeded
    cpu: "500m"       # Throttled if exceeded (not killed)
```

| Type | Behavior |
|---|---|
| `requests` | Used by scheduler to find a node with enough capacity |
| `limits.memory` | Container is OOMKilled if it exceeds this |
| `limits.cpu` | Container is CPU-throttled (not killed) |

**CPU units:**
- `1` = 1 full CPU core
- `500m` = half a core
- `100m` = 10% of a core

---

## Pod Labels & Selectors

Labels are key-value pairs attached to objects. Services, Deployments, and network policies use **selectors** to target pods by label.

```yaml
metadata:
  labels:
    app: my-app
    version: v2
    env: production
```

```bash
# List pods with label selectors
kubectl get pods -l app=my-app
kubectl get pods -l env=production,version=v2

# Add a label to an existing pod
kubectl label pod my-nginx color=blue

# Remove a label
kubectl label pod my-nginx color-
```

---

## Pod Scheduling: Node Selection

### nodeSelector (simple)

```yaml
spec:
  nodeSelector:
    disktype: ssd
```

Label the node first:
```bash
kubectl label node my-node disktype=ssd
```

### nodeAffinity (advanced)

```yaml
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
          - matchExpressions:
              - key: kubernetes.io/arch
                operator: In
                values: [amd64]
      preferredDuringSchedulingIgnoredDuringExecution:
        - weight: 1
          preference:
            matchExpressions:
              - key: disktype
                operator: In
                values: [ssd]
```

### Taints & Tolerations

**Taints** repel pods from nodes. **Tolerations** allow specific pods to land on tainted nodes.

```bash
# Taint a node (repels all pods without a toleration)
kubectl taint nodes gpu-node dedicated=gpu:NoSchedule
```

```yaml
# Pod that can tolerate the taint
spec:
  tolerations:
    - key: "dedicated"
      operator: "Equal"
      value: "gpu"
      effect: "NoSchedule"
```

---

## Nodes

### View Nodes

```bash
# List nodes with status
kubectl get nodes

# Detailed node info (capacity, conditions, images)
kubectl describe node my-node

# Node resource usage
kubectl top nodes
```

### Node Conditions

| Condition | Meaning |
|---|---|
| `Ready` | Node is healthy and accepting pods |
| `MemoryPressure` | Node is running low on memory |
| `DiskPressure` | Node is running low on disk |
| `PIDPressure` | Too many processes running on the node |
| `NetworkUnavailable` | Network not configured correctly |

### Drain and Cordon

```bash
# Cordon: mark node as unschedulable (no new pods, existing stay)
kubectl cordon my-node

# Drain: evict all pods + cordon (for maintenance)
kubectl drain my-node --ignore-daemonsets --delete-emptydir-data

# Uncordon: mark as schedulable again
kubectl uncordon my-node
```

---

## Useful Pod Commands Cheatsheet

```bash
# Get pods in all namespaces
kubectl get pods --all-namespaces
kubectl get pods -A

# Get pods on a specific node
kubectl get pods --field-selector spec.nodeName=my-node

# Watch pods in real time
kubectl get pods -w

# Get pod YAML definition
kubectl get pod my-pod -o yaml

# Copy a file from/to a pod
kubectl cp my-pod:/etc/nginx/nginx.conf ./nginx.conf
kubectl cp ./nginx.conf my-pod:/tmp/nginx.conf

# Run a temporary debug pod
kubectl run debug --image=busybox --rm -it --restart=Never -- sh
```


---

# Chapter 4: Deployments & ReplicaSets

## What is a Deployment?

A **Deployment** is the standard way to run stateless applications in Kubernetes. It manages:
- **Desired replica count** — always keeps N pods running
- **Rolling updates** — replace old pods with new ones with zero downtime
- **Rollback** — instantly go back to a previous version

Under the hood, a Deployment creates and manages a **ReplicaSet**, which manages the actual Pods.

```
Deployment
  └── ReplicaSet (v1: nginx:1.24)
        ├── Pod 1
        ├── Pod 2
        └── Pod 3
```

When you update the image:
```
Deployment
  ├── ReplicaSet (v1: nginx:1.24)  → scales down to 0
  └── ReplicaSet (v2: nginx:1.25)  → scales up to 3
```

---

## Creating a Deployment

**File:** `deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: default
  labels:
    app: my-app
spec:
  replicas: 3

  selector:
    matchLabels:
      app: my-app        # Must match template labels below

  template:
    metadata:
      labels:
        app: my-app      # Pods will get this label
    spec:
      containers:
        - name: my-app
          image: nginx:1.25-alpine
          ports:
            - containerPort: 80
          resources:
            requests:
              memory: "64Mi"
              cpu: "100m"
            limits:
              memory: "256Mi"
              cpu: "500m"
          readinessProbe:
            httpGet:
              path: /
              port: 80
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /
              port: 80
            initialDelaySeconds: 15
            periodSeconds: 20
```

```bash
kubectl apply -f deployment.yaml

kubectl get deployments
kubectl get pods -l app=my-app
```

---

## Scaling

```bash
# Scale to 5 replicas
kubectl scale deployment my-app --replicas=5

# Or edit the YAML and apply
kubectl apply -f deployment.yaml

# Watch pods scale up
kubectl get pods -w
```

---

## Rolling Updates (Zero Downtime)

### Update Strategy

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1           # Max extra pods during update
      maxUnavailable: 0     # Never take pods below desired count
```

| Strategy | Behavior |
|---|---|
| `RollingUpdate` (default) | Gradually replaces old pods with new ones |
| `Recreate` | Kills all old pods first, then creates new ones (causes downtime) |

### Trigger a Rolling Update

```bash
# Update the image
kubectl set image deployment/my-app my-app=nginx:1.26-alpine

# Watch the rollout
kubectl rollout status deployment/my-app

# See rollout history
kubectl rollout history deployment/my-app
```

### Annotate Updates (Best Practice)

```bash
kubectl set image deployment/my-app my-app=nginx:1.26-alpine \
  --record  # deprecated but still works

# Better: use --annotation
kubectl annotate deployment my-app \
  kubernetes.io/change-cause="Upgraded to nginx 1.26"
```

---

## Rollback

```bash
# Rollback to previous version
kubectl rollout undo deployment/my-app

# Rollback to a specific revision
kubectl rollout history deployment/my-app  # find revision number
kubectl rollout undo deployment/my-app --to-revision=2

# Check rollback status
kubectl rollout status deployment/my-app
```

---

## Readiness and Liveness Probes

**Liveness probe** — is the container alive? If it fails, Kubernetes **restarts** the container.

**Readiness probe** — is the container ready to receive traffic? If it fails, the pod is **removed from Service endpoints** (traffic stops routing to it).

```yaml
containers:
  - name: app
    image: my-app:latest

    # Readiness: HTTP check on /health
    readinessProbe:
      httpGet:
        path: /health
        port: 8080
      initialDelaySeconds: 5    # Wait 5s before first check
      periodSeconds: 10         # Check every 10s
      failureThreshold: 3       # 3 failures = not ready

    # Liveness: TCP socket check
    livenessProbe:
      tcpSocket:
        port: 8080
      initialDelaySeconds: 15
      periodSeconds: 20

    # Startup probe: gives slow-starting apps time to start
    startupProbe:
      httpGet:
        path: /health
        port: 8080
      failureThreshold: 30    # Wait up to 5 minutes (30 * 10s)
      periodSeconds: 10
```

---

## Environment Variables

```yaml
containers:
  - name: app
    image: my-app:latest
    env:
      # Direct value
      - name: APP_ENV
        value: "production"

      # From a ConfigMap
      - name: DB_HOST
        valueFrom:
          configMapKeyRef:
            name: app-config
            key: db_host

      # From a Secret
      - name: DB_PASSWORD
        valueFrom:
          secretKeyRef:
            name: app-secrets
            key: db_password

    # Load all keys from a ConfigMap as env vars
    envFrom:
      - configMapRef:
          name: app-config
      - secretRef:
          name: app-secrets
```

---

## Deployment Commands Reference

```bash
# Create
kubectl apply -f deployment.yaml

# View
kubectl get deployments
kubectl describe deployment my-app
kubectl get deployment my-app -o yaml

# Scale
kubectl scale deployment my-app --replicas=5

# Update image
kubectl set image deployment/my-app my-app=nginx:1.26

# Rollout status
kubectl rollout status deployment/my-app

# Rollout history
kubectl rollout history deployment/my-app

# Rollback
kubectl rollout undo deployment/my-app

# Pause a rollout (to make multiple changes before resuming)
kubectl rollout pause deployment/my-app
kubectl set image deployment/my-app my-app=nginx:1.26
kubectl set resources deployment/my-app -c=my-app --limits=cpu=200m,memory=512Mi
kubectl rollout resume deployment/my-app

# Delete
kubectl delete deployment my-app
```

---

## ReplicaSet (Direct Use — Rare)

You rarely create ReplicaSets directly (use Deployments instead). But it's useful to understand:

```bash
# View ReplicaSets (Deployment creates and manages these)
kubectl get replicasets
kubectl describe replicaset my-app-7d6b9f8c4
```

A Deployment keeps both old and new ReplicaSets around. Old ones are scaled to 0 (kept for rollback). The number of old ReplicaSets kept is controlled by `.spec.revisionHistoryLimit` (default: 10).

---

## Complete Example: App + Service + Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web-app
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: web-app
    spec:
      containers:
        - name: web-app
          image: nginx:1.25-alpine
          ports:
            - containerPort: 80
          resources:
            requests:
              memory: "64Mi"
              cpu: "100m"
            limits:
              memory: "256Mi"
              cpu: "500m"
          readinessProbe:
            httpGet:
              path: /
              port: 80
            initialDelaySeconds: 5
            periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: web-app-service
spec:
  selector:
    app: web-app
  ports:
    - port: 80
      targetPort: 80
  type: ClusterIP
```

```bash
kubectl apply -f deployment.yaml
kubectl get pods,deployments,services
kubectl port-forward svc/web-app-service 8080:80
```


---

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


---

# Chapter 6: ConfigMaps & Secrets

## Why Externalize Configuration?

Hard-coding config values in container images is an anti-pattern:
- You'd need different images per environment
- Secrets in images = security risk
- Config changes would require rebuilding images

Kubernetes solves this with **ConfigMaps** (non-sensitive config) and **Secrets** (sensitive data).

---

## ConfigMaps

### Create a ConfigMap

**From YAML:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: default
data:
  APP_ENV: "production"
  APP_DEBUG: "false"
  LOG_LEVEL: "info"
  DB_HOST: "mysql-service"
  DB_PORT: "3306"
  db.conf: |
    [client]
    host=mysql-service
    port=3306
    database=myapp
```

```bash
kubectl apply -f configmap.yaml
```

**From literal values:**
```bash
kubectl create configmap app-config \
  --from-literal=APP_ENV=production \
  --from-literal=LOG_LEVEL=info
```

**From a file:**
```bash
kubectl create configmap nginx-conf --from-file=nginx.conf
kubectl create configmap app-props --from-file=config/  # all files in dir
```

### Use ConfigMap as Environment Variables

```yaml
spec:
  containers:
    - name: app
      image: my-app:latest

      # Load all keys as env vars
      envFrom:
        - configMapRef:
            name: app-config

      # Load specific key
      env:
        - name: APP_ENV
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: APP_ENV
```

### Mount ConfigMap as Files

```yaml
spec:
  containers:
    - name: app
      image: my-app:latest
      volumeMounts:
        - name: config-volume
          mountPath: /etc/app/config   # Files appear here

  volumes:
    - name: config-volume
      configMap:
        name: app-config
        items:
          - key: db.conf
            path: database.conf        # Mounted as /etc/app/config/database.conf
```

---

## Secrets

Secrets store sensitive data like passwords, tokens, and keys. They are base64-encoded (not encrypted by default — use encryption at rest for production).

### Create a Secret

**From YAML (base64 values):**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
  namespace: default
type: Opaque
data:
  DB_PASSWORD: cGFzc3dvcmQxMjM=    # base64 of "password123"
  API_KEY: c2VjcmV0a2V5MTIz          # base64 of "secretkey123"
```

> Encode: `echo -n "password123" | base64`
> Decode: `echo "cGFzc3dvcmQxMjM=" | base64 -d`

**Using `stringData` (plain text — K8s encodes it):**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
stringData:
  DB_PASSWORD: "password123"
  API_KEY: "secretkey123"
  DATABASE_URL: "mysql://user:password123@mysql-service:3306/myapp"
```

**From command line:**
```bash
kubectl create secret generic app-secrets \
  --from-literal=DB_PASSWORD=password123 \
  --from-literal=API_KEY=secretkey123

# From a file (e.g., TLS cert)
kubectl create secret tls my-tls \
  --cert=tls.crt \
  --key=tls.key

# From a Docker registry credential
kubectl create secret docker-registry regcred \
  --docker-server=registry.example.com \
  --docker-username=myuser \
  --docker-password=mypassword
```

### Use Secrets as Environment Variables

```yaml
spec:
  containers:
    - name: app
      image: my-app:latest
      envFrom:
        - secretRef:
            name: app-secrets
      env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: DB_PASSWORD
```

### Mount Secret as Files

```yaml
spec:
  containers:
    - name: app
      volumeMounts:
        - name: secret-volume
          mountPath: /etc/secrets
          readOnly: true

  volumes:
    - name: secret-volume
      secret:
        secretName: app-secrets
```

### Pull Private Images with Image Pull Secrets

```yaml
spec:
  imagePullSecrets:
    - name: regcred
  containers:
    - name: app
      image: registry.example.com/my-private-app:latest
```

---

## Combined ConfigMap + Secret Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: web-config
data:
  APP_ENV: "production"
  REDIS_HOST: "redis-service"
  REDIS_PORT: "6379"
---
apiVersion: v1
kind: Secret
metadata:
  name: web-secrets
type: Opaque
stringData:
  DB_PASSWORD: "super-secret"
  JWT_SECRET: "my-jwt-signing-key"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: web-app
  template:
    metadata:
      labels:
        app: web-app
    spec:
      containers:
        - name: web-app
          image: my-app:latest
          envFrom:
            - configMapRef:
                name: web-config
            - secretRef:
                name: web-secrets
```

---

## Managing ConfigMaps and Secrets

```bash
# View ConfigMap
kubectl get configmaps
kubectl describe configmap app-config
kubectl get configmap app-config -o yaml

# Edit ConfigMap
kubectl edit configmap app-config

# View Secret (values are base64 encoded)
kubectl get secrets
kubectl describe secret app-secrets
kubectl get secret app-secrets -o yaml

# Decode a secret value
kubectl get secret app-secrets -o jsonpath='{.data.DB_PASSWORD}' | base64 -d

# Delete
kubectl delete configmap app-config
kubectl delete secret app-secrets
```

---

## Production Secret Management

> Default Kubernetes Secrets are only base64-encoded, not encrypted. For production, use:

| Solution | Description |
|---|---|
| **Kubernetes encryption at rest** | Encrypt etcd with `EncryptionConfiguration` |
| **HashiCorp Vault** | External secret store, injects secrets at pod startup |
| **AWS Secrets Manager** | With `external-secrets` operator |
| **Sealed Secrets** | Encrypt secrets in Git safely (by Bitnami) |
| **SOPS** | File-based secret encryption for GitOps |

### Sealed Secrets (Quick Example)

```bash
# Install kubeseal CLI
kubeseal --fetch-cert > pub-sealed-secrets.pem

# Seal a secret (safe to commit to Git)
kubectl create secret generic my-secret --from-literal=password=abc123 --dry-run=client -o yaml | \
  kubeseal --cert pub-sealed-secrets.pem -o yaml > my-sealed-secret.yaml

# Apply (controller decrypts it in cluster)
kubectl apply -f my-sealed-secret.yaml
```

---

## Best Practices

- Never hardcode passwords in Dockerfiles or application code
- Never commit raw Secrets YAML to Git (use Sealed Secrets or Vault)
- Use `readOnly: true` when mounting secrets as volumes
- Rotate secrets regularly and update the Secret object
- Use RBAC to restrict who can `get`/`list` secrets in each namespace
- Consider using `external-secrets` operator to sync from AWS/GCP/Azure secret managers


---

# Chapter 7: Persistent Volumes & Storage

## Why Persistent Storage?

Containers are **ephemeral** — when a pod dies, all data inside it is lost. For databases, file uploads, logs, and any stateful data, you need persistent storage that outlives the pod.

```
Pod dies → Container filesystem gone
Pod with PVC → Data survives in Persistent Volume
```

---

## Storage Hierarchy

```
StorageClass
    ↓ (defines provisioner + parameters)
PersistentVolume (PV)
    ↓ (actual storage resource)
PersistentVolumeClaim (PVC)
    ↓ (request for storage)
Pod → mounts the PVC
```

---

## PersistentVolume (PV)

A **PersistentVolume** is a piece of storage in the cluster provisioned by an admin (or automatically by a StorageClass). It has:
- A **capacity** (size)
- An **access mode**
- A **reclaim policy**

### Access Modes

| Mode | Abbreviation | Meaning |
|---|---|---|
| `ReadWriteOnce` | RWO | One node can read/write |
| `ReadOnlyMany` | ROX | Many nodes can read |
| `ReadWriteMany` | RWX | Many nodes can read/write |
| `ReadWriteOncePod` | RWOP | Only one pod can read/write |

### Static PV Example

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: my-pv
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain  # Keep data after PVC deleted
  storageClassName: manual
  hostPath:                               # Only for single-node dev!
    path: /mnt/data
```

---

## PersistentVolumeClaim (PVC)

A **PVC** is a request for storage by a user. Kubernetes finds a matching PV and binds them.

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-pvc
  namespace: default
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: manual   # Match the PV's storageClassName
  resources:
    requests:
      storage: 5Gi           # Request 5Gi (must be ≤ PV capacity)
```

### Mount PVC in a Pod

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
        - name: mysql
          image: mysql:8.0
          env:
            - name: MYSQL_ROOT_PASSWORD
              value: "rootpassword"
          volumeMounts:
            - name: mysql-data
              mountPath: /var/lib/mysql    # MySQL data directory
      volumes:
        - name: mysql-data
          persistentVolumeClaim:
            claimName: my-pvc             # Use the PVC
```

---

## StorageClass & Dynamic Provisioning

**StorageClass** enables **dynamic provisioning** — PVs are created automatically when a PVC is created, without admin intervention.

### Built-in StorageClasses

```bash
kubectl get storageclasses

# Output (example on cloud):
# NAME            PROVISIONER              RECLAIMPOLICY
# standard        kubernetes.io/gce-pd     Delete
# premium-rwo     pd.csi.storage.gke.io    Delete
```

### Custom StorageClass

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: pd.csi.storage.gke.io   # GCP provisioner
parameters:
  type: pd-ssd
reclaimPolicy: Retain
allowVolumeExpansion: true
```

### Dynamic PVC (No Manual PV Needed)

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: data-pvc
spec:
  storageClassName: fast-ssd     # Triggers dynamic provisioning
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 20Gi
```

---

## Reclaim Policies

| Policy | What Happens When PVC is Deleted |
|---|---|
| `Retain` | PV stays — data preserved, must be manually reclaimed |
| `Delete` | PV and underlying storage are deleted |
| `Recycle` | Deprecated — don't use |

> Use `Retain` for databases. Use `Delete` for temporary/ephemeral data.

---

## Storage for Multi-Pod Sharing (ReadWriteMany)

Some apps need multiple pods to share the same storage (e.g., Pimcore assets, shared file uploads):

| Provider | ReadWriteMany Solution |
|---|---|
| AWS | EFS with `aws-efs-csi-driver` |
| GCP | Filestore with `filestore.csi.storage.gke.io` |
| Azure | Azure Files |
| On-premise | NFS server with `nfs-subdir-external-provisioner` |
| Any | Rook-Ceph distributed storage |

```yaml
# PVC for shared file storage
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: shared-assets-pvc
spec:
  storageClassName: efs-sc       # Must support RWX
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 50Gi
```

---

## emptyDir (Temporary Shared Volume)

`emptyDir` is a temporary volume that exists as long as the pod lives. Data is lost when the pod is deleted. Useful for sharing data between containers in the same pod.

```yaml
spec:
  containers:
    - name: app
      image: my-app
      volumeMounts:
        - name: tmp-data
          mountPath: /tmp/cache

    - name: sidecar
      image: my-sidecar
      volumeMounts:
        - name: tmp-data
          mountPath: /shared

  volumes:
    - name: tmp-data
      emptyDir: {}
```

---

## hostPath (Single Node Dev Only)

Mounts a file or directory from the host node. **Never use in production multi-node clusters.**

```yaml
volumes:
  - name: host-data
    hostPath:
      path: /data/myapp
      type: DirectoryOrCreate
```

---

## Volume Expansion

```yaml
# PVC must have allowVolumeExpansion: true on its StorageClass
spec:
  resources:
    requests:
      storage: 50Gi   # Increase from 20Gi to 50Gi
```

```bash
kubectl patch pvc data-pvc -p '{"spec":{"resources":{"requests":{"storage":"50Gi"}}}}'
kubectl get pvc data-pvc  # watch for Resizing → FileSystemResizePending → Bound
```

---

## Storage Commands Cheatsheet

```bash
# View PVs
kubectl get pv
kubectl describe pv my-pv

# View PVCs
kubectl get pvc
kubectl get pvc -A   # all namespaces
kubectl describe pvc my-pvc

# View StorageClasses
kubectl get storageclass

# Check which pod uses which PVC
kubectl get pods -o json | jq '.items[].spec.volumes[].persistentVolumeClaim'

# Delete PVC (data fate depends on reclaimPolicy)
kubectl delete pvc my-pvc

# Force delete stuck PVC
kubectl patch pvc my-pvc -p '{"metadata":{"finalizers":null}}'
```

---

## Common Issues

| Issue | Cause | Fix |
|---|---|---|
| PVC stuck in `Pending` | No matching PV or StorageClass | Check `kubectl describe pvc` for events |
| PV in `Released` state | PVC deleted but PV not reclaimed | Delete and recreate PV or patch claimRef |
| Data lost after pod restart | Used emptyDir instead of PVC | Switch to PVC |
| Multi-pod write errors | Using RWO with multiple pods | Switch to RWX storage |
| Volume not expandable | StorageClass missing `allowVolumeExpansion: true` | Update StorageClass |


---

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


---

# Chapter 9: Ingress & Load Balancing

## What is Ingress?

A **Service** of type LoadBalancer gives you one IP per service — expensive on cloud providers. **Ingress** solves this by providing a single entry point that routes HTTP/HTTPS traffic to multiple services based on hostname or URL path.

```
Internet
    ↓
Ingress (single LoadBalancer IP)
    ├── /api    → api-service
    ├── /admin  → admin-service
    └── /       → frontend-service
```

---

## Ingress Controller

An **Ingress** resource alone does nothing — you need an **Ingress Controller** to implement it.

| Controller | When to Use |
|---|---|
| **Nginx Ingress Controller** | Most common, feature-rich |
| **Traefik** | Auto-discovers services, good for Kubernetes-native |
| **AWS ALB Ingress** | AWS-native, integrates with ALB |
| **GCE Ingress** | GKE-native, integrates with Google Cloud LB |
| **HAProxy** | High-performance edge |

---

## Install Nginx Ingress Controller

### On Minikube

```bash
minikube addons enable ingress
kubectl get pods -n ingress-nginx
```

### On Any Cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.10.1/deploy/static/provider/cloud/deploy.yaml

# Wait for controller to be ready
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s
```

---

## Basic Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-ingress
  namespace: default
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
    - host: myapp.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend-service
                port:
                  number: 80
```

---

## Path-Based Routing

Route different paths to different services on the same domain:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
  annotations:
    nginx.ingress.kubernetes.io/use-regex: "true"
spec:
  ingressClassName: nginx
  rules:
    - host: myapp.example.com
      http:
        paths:
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: api-service
                port:
                  number: 8080

          - path: /admin
            pathType: Prefix
            backend:
              service:
                name: admin-service
                port:
                  number: 3000

          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend-service
                port:
                  number: 80
```

---

## Host-Based Routing

Route different hostnames to different services:

```yaml
spec:
  rules:
    - host: api.myapp.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api-service
                port:
                  number: 8080

    - host: admin.myapp.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: admin-service
                port:
                  number: 3000

    - host: myapp.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend-service
                port:
                  number: 80
```

---

## TLS / HTTPS with cert-manager

### Install cert-manager

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.5/cert-manager.yaml

kubectl get pods -n cert-manager
```

### Create a Let's Encrypt ClusterIssuer

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
      - http01:
          ingress:
            class: nginx
```

### Ingress with TLS

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: secure-ingress
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - myapp.example.com
        - api.myapp.example.com
      secretName: myapp-tls        # cert-manager creates this Secret
  rules:
    - host: myapp.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend-service
                port:
                  number: 80
    - host: api.myapp.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api-service
                port:
                  number: 8080
```

---

## Common Nginx Ingress Annotations

```yaml
annotations:
  # Redirect HTTP to HTTPS
  nginx.ingress.kubernetes.io/ssl-redirect: "true"

  # Increase upload size limit
  nginx.ingress.kubernetes.io/proxy-body-size: "100m"

  # Increase timeouts
  nginx.ingress.kubernetes.io/proxy-read-timeout: "300"
  nginx.ingress.kubernetes.io/proxy-send-timeout: "300"

  # Rate limiting
  nginx.ingress.kubernetes.io/limit-rps: "10"

  # CORS
  nginx.ingress.kubernetes.io/enable-cors: "true"
  nginx.ingress.kubernetes.io/cors-allow-origin: "https://myapp.com"

  # Basic auth
  nginx.ingress.kubernetes.io/auth-type: basic
  nginx.ingress.kubernetes.io/auth-secret: basic-auth
  nginx.ingress.kubernetes.io/auth-realm: "Authentication Required"

  # Rewrite URL
  nginx.ingress.kubernetes.io/rewrite-target: /$2
```

---

## Ingress on Minikube (Local Testing)

```bash
# Enable ingress
minikube addons enable ingress

# Get Minikube IP
minikube ip   # e.g., 192.168.49.2

# Add to /etc/hosts for local DNS
echo "192.168.49.2 myapp.local" | sudo tee -a /etc/hosts

# Apply ingress pointing to myapp.local
# Now curl http://myapp.local works
```

---

## Ingress Commands

```bash
# List ingresses
kubectl get ingress
kubectl get ingress -A

# Describe (see rules and backend status)
kubectl describe ingress my-ingress

# Check ingress controller logs
kubectl logs -n ingress-nginx \
  deployment/ingress-nginx-controller

# Check cert-manager certificate status
kubectl get certificate
kubectl describe certificate myapp-tls
```

---

## Full Working Example

```yaml
# All in one: Deployment + Service + Ingress
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: hello-app
  template:
    metadata:
      labels:
        app: hello-app
    spec:
      containers:
        - name: hello-app
          image: hashicorp/http-echo
          args:
            - "-text=Hello from Kubernetes!"
          ports:
            - containerPort: 5678
---
apiVersion: v1
kind: Service
metadata:
  name: hello-service
spec:
  selector:
    app: hello-app
  ports:
    - port: 80
      targetPort: 5678
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: hello-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
    - host: hello.local
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: hello-service
                port:
                  number: 80
```


---

# Chapter 10: Helm — Package Manager for Kubernetes

## What is Helm?

**Helm** is the package manager for Kubernetes. Instead of maintaining many separate YAML files for each environment, Helm bundles them into a **Chart** — a reusable, versioned, configurable package.

> **Analogy:** Helm is to Kubernetes what `apt` or `brew` is to Linux/macOS.

---

## Why Use Helm?

| Without Helm | With Helm |
|---|---|
| 10+ YAML files per app | One chart |
| Manual sed/envsubst for different environments | `values.yaml` per environment |
| No versioning | Chart versioning + release history |
| Manual rollback = error-prone | `helm rollback` |
| Copy-paste config per cluster | `helm install -f values.prod.yaml` |

---

## Install Helm

```bash
# Linux/macOS (recommended)
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# macOS via Homebrew
brew install helm

# Verify
helm version
```

---

## Key Concepts

| Concept | Description |
|---|---|
| **Chart** | Package of Kubernetes YAML templates + metadata |
| **Release** | A deployed instance of a chart |
| **Repository** | Storage for charts (like Docker Hub for images) |
| **Values** | Config variables that customize a chart |

---

## Adding Chart Repositories

```bash
# Add popular repos
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add stable https://charts.helm.sh/stable
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx

# Update repos
helm repo update

# Search for a chart
helm search repo nginx
helm search repo mysql
helm search hub wordpress   # Search Artifact Hub
```

---

## Installing Charts

```bash
# Install with default values
helm install my-nginx ingress-nginx/ingress-nginx

# Install with custom values
helm install my-mysql bitnami/mysql \
  --set auth.rootPassword=secret123 \
  --set primary.persistence.size=20Gi

# Install from a values file
helm install my-mysql bitnami/mysql -f values.yaml

# Install in a specific namespace
helm install my-mysql bitnami/mysql -n database --create-namespace

# Install and wait for readiness
helm install my-mysql bitnami/mysql --wait --timeout=300s
```

---

## Viewing Releases

```bash
# List all releases
helm list
helm list -A   # all namespaces

# Get details of a release
helm status my-mysql

# Get the values used in a release
helm get values my-mysql

# Get all rendered manifests
helm get manifest my-mysql
```

---

## Upgrading and Rolling Back

```bash
# Upgrade a release (new image, new values)
helm upgrade my-mysql bitnami/mysql \
  --set auth.rootPassword=newpassword

# Upgrade from values file
helm upgrade my-mysql bitnami/mysql -f values.prod.yaml

# Upgrade or install if not exists
helm upgrade --install my-mysql bitnami/mysql -f values.yaml

# View history
helm history my-mysql

# Rollback to previous version
helm rollback my-mysql

# Rollback to specific revision
helm rollback my-mysql 2
```

---

## Uninstalling

```bash
helm uninstall my-mysql
helm uninstall my-mysql -n database
```

---

## Creating Your Own Chart

```bash
# Generate scaffold
helm create my-app

# Structure:
# my-app/
# ├── Chart.yaml          # Chart metadata
# ├── values.yaml         # Default values
# ├── templates/
# │   ├── deployment.yaml
# │   ├── service.yaml
# │   ├── ingress.yaml
# │   ├── _helpers.tpl    # Template helper functions
# │   ├── hpa.yaml
# │   └── serviceaccount.yaml
# └── charts/             # Sub-charts (dependencies)
```

### Chart.yaml

```yaml
apiVersion: v2
name: my-app
description: A Helm chart for my application
type: application
version: 0.1.0        # Chart version
appVersion: "1.0.0"   # App version (image tag)
```

### values.yaml

```yaml
replicaCount: 3

image:
  repository: my-registry/my-app
  tag: "1.0.0"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: true
  className: nginx
  host: myapp.example.com
  tls: true

resources:
  requests:
    cpu: "100m"
    memory: "128Mi"
  limits:
    cpu: "500m"
    memory: "512Mi"

env:
  APP_ENV: production
  LOG_LEVEL: info

secrets:
  DB_PASSWORD: ""
  API_KEY: ""
```

### templates/deployment.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "my-app.fullname" . }}
  labels:
    {{- include "my-app.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "my-app.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "my-app.selectorLabels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - containerPort: {{ .Values.service.port }}
          env:
            {{- range $key, $val := .Values.env }}
            - name: {{ $key }}
              value: {{ $val | quote }}
            {{- end }}
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
```

---

## Multiple Environments with Values Files

```bash
# values.yaml       — defaults
# values.dev.yaml   — dev overrides
# values.prod.yaml  — production overrides
```

**values.prod.yaml:**
```yaml
replicaCount: 5
image:
  tag: "2.1.0"
resources:
  requests:
    cpu: "500m"
    memory: "512Mi"
  limits:
    cpu: "2"
    memory: "2Gi"
ingress:
  host: myapp.com
```

```bash
# Deploy to dev
helm upgrade --install my-app ./my-app -f values.dev.yaml -n dev

# Deploy to production
helm upgrade --install my-app ./my-app -f values.prod.yaml -n production
```

---

## Lint and Debug

```bash
# Lint the chart for errors
helm lint ./my-app

# Render templates without deploying (dry run)
helm template my-app ./my-app -f values.prod.yaml

# Dry run against a real cluster
helm install my-app ./my-app --dry-run --debug

# Show computed values
helm install my-app ./my-app --debug 2>&1 | head -50
```

---

## Helm Cheatsheet

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm search repo nginx
helm install my-release bitnami/nginx
helm list
helm upgrade my-release bitnami/nginx --set replicaCount=3
helm rollback my-release
helm history my-release
helm uninstall my-release
helm create my-chart
helm lint ./my-chart
helm template ./my-chart
helm package ./my-chart
```


---

# Chapter 11: StatefulSets & DaemonSets

## StatefulSets

### What is a StatefulSet?

A **StatefulSet** is like a Deployment but for **stateful applications** that need:
- **Stable, persistent pod names** (pod-0, pod-1, pod-2 — not random hashes)
- **Stable network identities** (each pod gets its own DNS hostname)
- **Ordered deployment and scaling** (pod-0 starts before pod-1)
- **Per-pod persistent volumes** (each pod gets its own PVC)

### When to Use StatefulSet vs Deployment

| Feature | Deployment | StatefulSet |
|---|---|---|
| Pod names | Random (nginx-abc123) | Stable (mysql-0, mysql-1) |
| Start order | Any order | Ordered (0, 1, 2...) |
| Storage | Shared PVC or emptyDir | Each pod has own PVC |
| Use case | Stateless apps (web, API) | Databases, queues, caches |

---

### StatefulSet Example: MySQL

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
  namespace: default
spec:
  serviceName: mysql-headless       # Must match headless service name
  replicas: 1
  selector:
    matchLabels:
      app: mysql

  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
        - name: mysql
          image: mysql:8.0
          ports:
            - containerPort: 3306
          env:
            - name: MYSQL_ROOT_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: mysql-secret
                  key: root-password
          volumeMounts:
            - name: mysql-data
              mountPath: /var/lib/mysql
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "2Gi"
              cpu: "1"

  # Each pod gets its own PVC
  volumeClaimTemplates:
    - metadata:
        name: mysql-data
      spec:
        accessModes: ["ReadWriteOnce"]
        storageClassName: fast-ssd
        resources:
          requests:
            storage: 20Gi
---
# Headless Service (required for stable DNS)
apiVersion: v1
kind: Service
metadata:
  name: mysql-headless
spec:
  clusterIP: None           # Headless — no virtual IP
  selector:
    app: mysql
  ports:
    - port: 3306
---
# Regular Service for app access
apiVersion: v1
kind: Service
metadata:
  name: mysql-service
spec:
  selector:
    app: mysql
  ports:
    - port: 3306
```

**DNS for each pod:**
```
mysql-0.mysql-headless.default.svc.cluster.local
mysql-1.mysql-headless.default.svc.cluster.local
```

---

### StatefulSet Scaling

```bash
# Scale up (pods added in order: mysql-2, then mysql-3)
kubectl scale statefulset mysql --replicas=3

# Scale down (pods removed in reverse order: mysql-2 first)
kubectl scale statefulset mysql --replicas=1

# Watch ordered creation
kubectl get pods -w -l app=mysql
```

---

### StatefulSet Update Strategies

```yaml
spec:
  updateStrategy:
    type: RollingUpdate    # Default — updates pods one at a time in reverse order
    rollingUpdate:
      partition: 0         # Only update pods with ordinal >= partition
```

```bash
# Canary update: only update mysql-1 and mysql-2 (not mysql-0)
kubectl patch statefulset mysql -p '{"spec":{"updateStrategy":{"rollingUpdate":{"partition":1}}}}'
kubectl set image statefulset/mysql mysql=mysql:8.1
# After validation:
kubectl patch statefulset mysql -p '{"spec":{"updateStrategy":{"rollingUpdate":{"partition":0}}}}'
```

---

## DaemonSets

### What is a DaemonSet?

A **DaemonSet** ensures that **one pod runs on every node** (or a subset of nodes). When a new node is added to the cluster, the pod is automatically scheduled on it.

**Use cases:**
- Log collection agents (Fluentd, Filebeat)
- Node monitoring agents (Prometheus Node Exporter, Datadog)
- Network plugins (CNI plugins)
- Security scanning agents

---

### DaemonSet Example: Log Collector

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd
  namespace: kube-system
  labels:
    app: fluentd
spec:
  selector:
    matchLabels:
      app: fluentd
  template:
    metadata:
      labels:
        app: fluentd
    spec:
      tolerations:
        # Also run on master/control-plane nodes
        - key: node-role.kubernetes.io/control-plane
          operator: Exists
          effect: NoSchedule

      containers:
        - name: fluentd
          image: fluent/fluentd-kubernetes-daemonset:v1-debian-elasticsearch
          env:
            - name: FLUENT_ELASTICSEARCH_HOST
              value: elasticsearch-service
            - name: FLUENT_ELASTICSEARCH_PORT
              value: "9200"
          resources:
            limits:
              memory: 200Mi
            requests:
              cpu: 100m
              memory: 200Mi
          volumeMounts:
            - name: varlog
              mountPath: /var/log
            - name: varlibdockercontainers
              mountPath: /var/lib/docker/containers
              readOnly: true

      volumes:
        - name: varlog
          hostPath:
            path: /var/log
        - name: varlibdockercontainers
          hostPath:
            path: /var/lib/docker/containers
```

---

### DaemonSet on Specific Nodes Only

```yaml
spec:
  template:
    spec:
      # Only run on nodes with this label
      nodeSelector:
        dedicated: monitoring

      # Or use affinity for more complex rules
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
              - matchExpressions:
                  - key: node-type
                    operator: In
                    values: [gpu, high-memory]
```

---

### DaemonSet Commands

```bash
# Create
kubectl apply -f daemonset.yaml

# List (shows desired vs ready per node)
kubectl get daemonsets
kubectl get daemonsets -n kube-system

# Describe
kubectl describe daemonset fluentd -n kube-system

# Check which nodes the DS is running on
kubectl get pods -l app=fluentd -o wide

# Rolling update (same as Deployment)
kubectl set image daemonset/fluentd fluentd=fluent/fluentd:v1.16
kubectl rollout status daemonset/fluentd

# Rollback
kubectl rollout undo daemonset/fluentd
```

---

## Comparison Summary

| Feature | Deployment | StatefulSet | DaemonSet |
|---|---|---|---|
| Pod count control | `replicas: N` | `replicas: N` | One per node |
| Pod name | Random hash | Stable (0, 1, 2) | Per node |
| Storage | Shared/none | Per-pod PVC | HostPath |
| Order | Arbitrary | Ordered | Per node |
| Use case | Stateless apps | DBs, caches, queues | Node agents |
| Scale up | Instant, any order | Ordered (0→1→2) | Auto when node joins |
| Scale down | Any order | Reverse order (2→1→0) | Auto when node leaves |


---

# Chapter 12: Jobs & CronJobs

## Jobs

### What is a Job?

A **Job** creates one or more Pods to run a task to **completion**. Unlike a Deployment (which keeps pods running forever), a Job:
- Runs a task once
- Marks itself as complete when the pod exits with code `0`
- Retries if the pod fails

**Use cases:**
- Database migrations
- Data processing batch jobs
- Report generation
- One-time setup tasks
- Image processing pipelines

---

### Basic Job

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: db-migration
spec:
  completions: 1          # How many successful pod completions needed
  parallelism: 1          # How many pods run at the same time
  backoffLimit: 3         # Retry up to 3 times before marking failed
  activeDeadlineSeconds: 600  # Kill job if takes > 10 minutes
  ttlSecondsAfterFinished: 3600  # Auto-delete job 1hr after completion

  template:
    spec:
      restartPolicy: OnFailure    # Required for Jobs (Never or OnFailure)
      containers:
        - name: migration
          image: my-app:latest
          command: ["bin/console", "doctrine:migrations:migrate", "--no-interaction"]
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: app-secrets
                  key: DATABASE_URL
          resources:
            requests:
              memory: "256Mi"
              cpu: "100m"
```

```bash
kubectl apply -f job.yaml
kubectl get jobs
kubectl get pods -l job-name=db-migration
kubectl logs job/db-migration
```

---

### Parallel Jobs

Run multiple pods at the same time to speed up processing:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: image-processor
spec:
  completions: 10      # Process 10 batches total
  parallelism: 3       # Run 3 pods simultaneously
  backoffLimit: 2
  template:
    spec:
      restartPolicy: OnFailure
      containers:
        - name: processor
          image: image-processor:latest
          env:
            - name: BATCH_SIZE
              value: "100"
```

---

### Job Patterns

| Pattern | completions | parallelism | Use Case |
|---|---|---|---|
| Single task | 1 | 1 | One-shot migration, setup |
| Parallel fixed count | N | N | Process N batches |
| Work queue | unset | N | Pull from a queue until empty |

---

### Job Lifecycle & Cleanup

```bash
# View all jobs
kubectl get jobs

# View job details (see success/failure counts)
kubectl describe job db-migration

# View pods created by the job
kubectl get pods --selector=job-name=db-migration

# Follow logs
kubectl logs -f job/db-migration

# Delete job and its pods
kubectl delete job db-migration

# Auto-cleanup: add this to job spec
spec:
  ttlSecondsAfterFinished: 300   # Delete 5 minutes after completion
```

---

## CronJobs

### What is a CronJob?

A **CronJob** creates Jobs on a schedule — like Linux `cron` but in Kubernetes.

**Use cases:**
- Nightly database backups
- Daily report generation
- Hourly cache warming
- Periodic cleanup tasks

---

### CronJob Syntax

```
┌─────── minute (0-59)
│ ┌───── hour (0-23)
│ │ ┌─── day of month (1-31)
│ │ │ ┌─ month (1-12)
│ │ │ │ └── day of week (0-7, 0=Sun)
│ │ │ │
* * * * *
```

| Schedule | Meaning |
|---|---|
| `0 2 * * *` | Every day at 2:00 AM |
| `*/15 * * * *` | Every 15 minutes |
| `0 0 * * 0` | Every Sunday at midnight |
| `0 9 * * 1-5` | Every weekday at 9 AM |
| `@daily` | Once a day (= `0 0 * * *`) |
| `@hourly` | Once an hour |
| `@weekly` | Once a week |

---

### CronJob Example: Nightly Backup

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: nightly-backup
  namespace: production
spec:
  schedule: "0 2 * * *"       # 2 AM every day
  timeZone: "Asia/Kolkata"     # Optional timezone (K8s 1.27+)

  concurrencyPolicy: Forbid    # Don't run if previous job still running
  successfulJobsHistoryLimit: 3  # Keep last 3 successful jobs
  failedJobsHistoryLimit: 1      # Keep last 1 failed job
  startingDeadlineSeconds: 300   # If missed, skip (don't catch up after 5min)

  jobTemplate:
    spec:
      backoffLimit: 2
      ttlSecondsAfterFinished: 86400   # Delete job after 24hrs
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: backup
              image: mysql:8.0
              command:
                - /bin/sh
                - -c
                - |
                  mysqldump -h mysql-service -u root -p$MYSQL_ROOT_PASSWORD mydb | \
                  gzip > /backup/mydb-$(date +%Y%m%d-%H%M%S).sql.gz
                  echo "Backup completed"
              env:
                - name: MYSQL_ROOT_PASSWORD
                  valueFrom:
                    secretKeyRef:
                      name: mysql-secret
                      key: root-password
              volumeMounts:
                - name: backup-storage
                  mountPath: /backup
          volumes:
            - name: backup-storage
              persistentVolumeClaim:
                claimName: backup-pvc
```

---

### CronJob Concurrency Policies

| Policy | Behavior |
|---|---|
| `Allow` (default) | Multiple jobs can run at the same time |
| `Forbid` | Skip new job if previous is still running |
| `Replace` | Cancel current job and start a new one |

---

### CronJob Commands

```bash
# View cronjobs
kubectl get cronjobs
kubectl get cj   # short form

# Describe (see last schedule, active jobs)
kubectl describe cronjob nightly-backup

# Manually trigger a cronjob right now
kubectl create job manual-backup-now \
  --from=cronjob/nightly-backup

# View jobs created by cronjob
kubectl get jobs -l app=nightly-backup

# Watch job execution
kubectl get jobs -w

# Suspend a cronjob (pause scheduling)
kubectl patch cronjob nightly-backup -p '{"spec":{"suspend":true}}'

# Resume
kubectl patch cronjob nightly-backup -p '{"spec":{"suspend":false}}'

# Delete
kubectl delete cronjob nightly-backup
```

---

### Common Job Failures & Fixes

| Issue | Cause | Fix |
|---|---|---|
| Job stuck in `Pending` | Resource constraints | Check node capacity, `kubectl describe pod` |
| Pod `Error` immediately | Container command fails | Check `kubectl logs pod/...` |
| Job never completes | `restartPolicy: Always` set | Change to `OnFailure` or `Never` |
| CronJob misses schedule | `startingDeadlineSeconds` too low | Increase or remove the deadline |
| Duplicate jobs running | `concurrencyPolicy: Allow` with slow job | Change to `Forbid` |
| Old jobs pile up | No `ttlSecondsAfterFinished` | Add TTL or history limits |


---

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


---

# Chapter 14: Monitoring & Logging

## Overview

A production Kubernetes cluster needs:
- **Metrics** — CPU, memory, request rates, error rates
- **Logs** — container stdout/stderr, application logs
- **Alerts** — notify when something goes wrong
- **Dashboards** — visualize everything at a glance

**Standard Stack:**
```
Prometheus (collect metrics)
    + Grafana (visualize)
    + Alertmanager (alerting)
    + Loki (log aggregation)
    + Promtail/Fluentd (log shipping)
```

---

## kubectl Logs

### Basic Log Commands

```bash
# Logs from a pod (latest container)
kubectl logs my-pod

# Follow logs in real time
kubectl logs -f my-pod

# Last 100 lines
kubectl logs my-pod --tail=100

# Logs since a time
kubectl logs my-pod --since=1h
kubectl logs my-pod --since-time=2025-01-01T10:00:00Z

# Logs from a specific container in a multi-container pod
kubectl logs my-pod -c nginx
kubectl logs my-pod -c php-fpm

# Previous container (after crash)
kubectl logs my-pod --previous

# Logs from all pods matching a label
kubectl logs -l app=my-app --all-containers=true

# Logs from a deployment
kubectl logs deployment/my-app
kubectl logs deployment/my-app -c my-app -f
```

---

## metrics-server

**metrics-server** enables `kubectl top` commands for basic CPU/memory metrics.

```bash
# Install
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# On Minikube
minikube addons enable metrics-server

# Use
kubectl top nodes
kubectl top pods -A
kubectl top pods -n production --sort-by=cpu
```

---

## Prometheus + Grafana Stack

### Install with Helm (Recommended)

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Install kube-prometheus-stack (includes Prometheus, Grafana, Alertmanager)
helm install monitoring prometheus-community/kube-prometheus-stack \
  -n monitoring --create-namespace \
  --set grafana.adminPassword=admin123 \
  --set prometheus.prometheusSpec.retention=30d

# Check pods
kubectl get pods -n monitoring
```

### Access Grafana

```bash
kubectl port-forward -n monitoring svc/monitoring-grafana 3000:80
# Open: http://localhost:3000
# Login: admin / admin123
```

### Access Prometheus

```bash
kubectl port-forward -n monitoring svc/monitoring-kube-prometheus-prometheus 9090:9090
# Open: http://localhost:9090
```

---

## Prometheus Basics

### Key Concepts

| Concept | Description |
|---|---|
| **Scrape** | Prometheus pulls metrics from targets at intervals |
| **Metric** | Named time-series data point |
| **Label** | Key-value tags on metrics |
| **PromQL** | Query language for metrics |
| **Exporter** | App that exposes metrics at `/metrics` |

### Common PromQL Queries

```promql
# CPU usage by pod
sum(rate(container_cpu_usage_seconds_total[5m])) by (pod)

# Memory usage by namespace
sum(container_memory_usage_bytes) by (namespace)

# HTTP error rate
rate(http_requests_total{status=~"5.."}[5m])

# Pod restart count
kube_pod_container_status_restarts_total

# Node disk pressure
node_filesystem_avail_bytes / node_filesystem_size_bytes < 0.10

# Request latency p99
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))
```

### ServiceMonitor (Auto-scrape your apps)

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: my-app-monitor
  namespace: production
spec:
  selector:
    matchLabels:
      app: my-app
  endpoints:
    - port: http
      path: /metrics
      interval: 30s
```

Your app must expose a `/metrics` endpoint in Prometheus format (use a Prometheus client library).

---

## Alerting with Alertmanager

### PrometheusRule Example

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: my-app-alerts
  namespace: production
spec:
  groups:
    - name: my-app
      rules:
        - alert: HighErrorRate
          expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
          for: 2m
          labels:
            severity: critical
          annotations:
            summary: "High error rate on {{ $labels.pod }}"
            description: "Error rate is {{ $value | humanizePercentage }}"

        - alert: PodCrashLooping
          expr: rate(kube_pod_container_status_restarts_total[15m]) > 0
          for: 5m
          labels:
            severity: warning
          annotations:
            summary: "Pod {{ $labels.pod }} is crash looping"

        - alert: HighMemoryUsage
          expr: |
            container_memory_usage_bytes / container_spec_memory_limit_bytes > 0.9
          for: 5m
          labels:
            severity: warning
          annotations:
            summary: "Container using >90% memory limit"
```

---

## Loki — Log Aggregation

**Loki** is Grafana's log aggregation system — lightweight, stores logs by labels, not full-text indexing.

### Install Loki Stack

```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm install loki grafana/loki-stack \
  -n monitoring \
  --set grafana.enabled=false \
  --set promtail.enabled=true
```

**Promtail** auto-discovers pods and ships their logs to Loki.

### Query Logs in Grafana (LogQL)

```logql
# All logs from my-app in production
{app="my-app", namespace="production"}

# Errors only
{app="my-app"} |= "ERROR"

# Count errors per minute
rate({app="my-app"} |= "ERROR" [1m])

# Parse JSON logs
{app="my-app"} | json | level="error"
```

---

## Pre-built Grafana Dashboards

Import these from grafana.com/dashboards:

| Dashboard | ID | What It Shows |
|---|---|---|
| Kubernetes cluster overview | 315 | Node CPU/memory, pod counts |
| Node Exporter | 1860 | Per-node CPU, memory, disk, network |
| Pod monitoring | 6417 | Pod CPU/memory/restarts |
| Nginx Ingress | 9614 | Request rates, errors, latency |
| MySQL | 7362 | Queries, connections, slow queries |
| Redis | 11835 | Operations, memory, hit rate |

```bash
# In Grafana UI:
# Dashboards → Import → Enter ID → Load → Select Prometheus data source → Import
```

---

## Kubernetes Events

Events are lightweight logs of what happened to K8s objects:

```bash
# All events in a namespace
kubectl get events -n production

# Sort by time
kubectl get events -n production --sort-by='.lastTimestamp'

# Events for a specific pod
kubectl describe pod my-pod -n production  # Shows Events at bottom

# Watch events in real time
kubectl get events -w -n production
```

---

## Monitoring Checklist for Production

- [ ] metrics-server installed (`kubectl top` works)
- [ ] Prometheus scraping cluster metrics
- [ ] Grafana dashboards for cluster + applications
- [ ] Alertmanager configured (Slack/PagerDuty/email)
- [ ] Loki/EFK for log aggregation
- [ ] Pod restart alerts set up
- [ ] Node disk/memory pressure alerts
- [ ] Application-level metrics exposed at `/metrics`
- [ ] Log retention policy defined
- [ ] On-call runbook linked in alert annotations


---

# Chapter 15: Security & Best Practices

## Security Layers in Kubernetes

```
1. Cluster Infrastructure  — secure nodes, network, etcd encryption
2. Authentication          — who can connect to the API server
3. Authorization (RBAC)    — what they can do
4. Admission Control       — validate/mutate resources before creation
5. Network Policies        — firewall between pods
6. Pod Security            — what containers can do on the node
7. Image Security          — trusted, scanned images only
8. Secrets Management      — proper handling of sensitive data
```

---

## SecurityContext

**SecurityContext** controls what a container can do on the host.

### Pod-level SecurityContext

```yaml
spec:
  securityContext:
    runAsNonRoot: true        # Container must not run as root
    runAsUser: 1000           # Run as user ID 1000
    runAsGroup: 3000
    fsGroup: 2000             # All volumes owned by this group
    seccompProfile:
      type: RuntimeDefault    # Apply default seccomp filter
```

### Container-level SecurityContext

```yaml
containers:
  - name: app
    securityContext:
      allowPrivilegeEscalation: false   # Can't gain more privileges
      readOnlyRootFilesystem: true      # Filesystem is read-only
      capabilities:
        drop:
          - ALL                         # Drop all Linux capabilities
        add:
          - NET_BIND_SERVICE            # Only add what's needed
```

### Recommended Secure Container Template

```yaml
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    fsGroup: 2000
    seccompProfile:
      type: RuntimeDefault
  containers:
    - name: app
      image: my-app:latest
      securityContext:
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
        capabilities:
          drop: ["ALL"]
      volumeMounts:
        - name: tmp
          mountPath: /tmp          # Writable tmp directory
        - name: cache
          mountPath: /app/cache
  volumes:
    - name: tmp
      emptyDir: {}
    - name: cache
      emptyDir: {}
```

---

## Pod Security Standards (PSS)

Kubernetes enforces security policies via **Pod Security Admission** (replaced PodSecurityPolicy in K8s 1.25+).

Three levels:
| Level | Behavior |
|---|---|
| `privileged` | No restrictions |
| `baseline` | Minimal restrictions, prevents known privilege escalations |
| `restricted` | Heavily restricted, follows current pod hardening best practices |

Apply at the namespace level:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

---

## Network Policies

By default, all pods can communicate with all other pods. Use NetworkPolicy to restrict this.

### Zero-Trust: Deny All, Then Allow Selectively

```yaml
# 1. Deny all ingress and egress
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all
  namespace: production
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress
---
# 2. Allow frontend to reach backend
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-frontend-to-backend
  namespace: production
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
    - Ingress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: frontend
      ports:
        - port: 8080
---
# 3. Allow all pods to reach DNS (required for service discovery)
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-dns
  namespace: production
spec:
  podSelector: {}
  policyTypes:
    - Egress
  egress:
    - ports:
        - port: 53
          protocol: UDP
        - port: 53
          protocol: TCP
```

---

## Image Security

### Use Specific Tags (Never `latest` in Production)

```yaml
# Bad
image: nginx:latest

# Good
image: nginx:1.25.3
# Or use digest (immutable)
image: nginx@sha256:abc123...
```

### Image Scanning

```bash
# Trivy: open-source vulnerability scanner
trivy image nginx:1.25.3
trivy image my-app:latest

# In CI/CD pipeline
trivy image --exit-code 1 --severity HIGH,CRITICAL my-app:latest
```

### Limit Image Registries (OPA/Kyverno)

```yaml
# Kyverno policy: only allow images from trusted registry
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: allowed-registries
spec:
  validationFailureAction: Enforce
  rules:
    - name: check-registry
      match:
        resources:
          kinds: [Pod]
      validate:
        message: "Image must come from registry.company.com"
        pattern:
          spec:
            containers:
              - image: "registry.company.com/*"
```

---

## etcd Encryption at Rest

```yaml
# /etc/kubernetes/encryption-config.yaml
apiVersion: apiserver.config.k8s.io/v1
kind: EncryptionConfiguration
resources:
  - resources:
      - secrets
    providers:
      - aescbc:
          keys:
            - name: key1
              secret: <base64-encoded-32-byte-key>
      - identity: {}
```

Add to kube-apiserver: `--encryption-provider-config=/etc/kubernetes/encryption-config.yaml`

---

## TLS Everywhere

```bash
# All internal K8s communication is TLS by default
# For your apps: use cert-manager (Chapter 9) to manage certs
# For ingress: always enable TLS termination at the ingress layer

# Rotate kubelet certificates
kubeadm certs renew all
```

---

## Admission Controllers

Admission controllers intercept requests before they're saved to etcd.

**OPA Gatekeeper** and **Kyverno** are popular policy engines:

```bash
# Install Kyverno
helm install kyverno kyverno/kyverno -n kyverno --create-namespace
```

```yaml
# Require all pods to have resource limits
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-resource-limits
spec:
  validationFailureAction: Enforce
  rules:
    - name: check-limits
      match:
        resources:
          kinds: [Pod]
      validate:
        message: "Resource limits are required"
        pattern:
          spec:
            containers:
              - resources:
                  limits:
                    memory: "?*"
                    cpu: "?*"
```

---

## Security Audit Commands

```bash
# Check all privileged pods
kubectl get pods -A -o json | jq '.items[] | select(.spec.containers[].securityContext.privileged == true) | .metadata.name'

# Check pods running as root
kubectl get pods -A -o json | jq '.items[] | select(.spec.securityContext.runAsUser == 0 or .spec.securityContext.runAsNonRoot != true) | .metadata.name'

# Who has cluster-admin?
kubectl get clusterrolebindings -o json | jq '.items[] | select(.roleRef.name == "cluster-admin") | .subjects'

# Check exposed secrets
kubectl get secrets -A

# Audit API server access logs
kubectl logs -n kube-system kube-apiserver-<node> | grep "authentication"

# kube-bench: CIS Kubernetes benchmark
kubectl apply -f https://raw.githubusercontent.com/aquasecurity/kube-bench/main/job.yaml
kubectl logs job/kube-bench
```

---

## Security Checklist

- [ ] No containers running as root (`runAsNonRoot: true`)
- [ ] All containers have `allowPrivilegeEscalation: false`
- [ ] `readOnlyRootFilesystem: true` where possible
- [ ] All capabilities dropped (`drop: [ALL]`)
- [ ] NetworkPolicies restrict east-west traffic
- [ ] RBAC follows least privilege
- [ ] Secrets not in Git (use Sealed Secrets / Vault)
- [ ] etcd encryption at rest enabled
- [ ] Images scanned for CVEs before deployment
- [ ] Pod Security Standards enforced (`restricted` for prod)
- [ ] Ingress has TLS termination
- [ ] Admission controller (Kyverno/Gatekeeper) enforcing policies
- [ ] Regular CIS benchmark audits (kube-bench)


---

# Chapter 16: CI/CD with Kubernetes

## CI/CD Pipeline Overview

```
Developer pushes code
    ↓
CI (GitHub Actions / GitLab CI)
    ├── Run tests
    ├── Build Docker image
    ├── Push to registry
    └── Deploy to Kubernetes
         ├── kubectl apply (imperative)
         └── ArgoCD sync (GitOps)
```

---

## GitHub Actions: Build & Deploy

### Complete Workflow

**File:** `.github/workflows/deploy.yml`

```yaml
name: Build and Deploy

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: test
          MYSQL_DATABASE: testdb
        ports:
          - 3306:3306
        options: --health-cmd="mysqladmin ping" --health-interval=10s

    steps:
      - uses: actions/checkout@v4

      - name: Run tests
        run: |
          composer install
          vendor/bin/phpunit

  build:
    needs: test
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    outputs:
      image-tag: ${{ steps.meta.outputs.tags }}

    steps:
      - uses: actions/checkout@v4

      - name: Log in to registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=sha,prefix=sha-
            type=ref,event=branch
            type=semver,pattern={{version}}

      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy-staging:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: staging

    steps:
      - uses: actions/checkout@v4

      - name: Configure kubectl
        uses: azure/k8s-set-context@v3
        with:
          method: kubeconfig
          kubeconfig: ${{ secrets.KUBE_CONFIG_STAGING }}

      - name: Deploy to staging
        run: |
          IMAGE_TAG="sha-${{ github.sha }}"
          kubectl set image deployment/my-app \
            my-app=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:$IMAGE_TAG \
            -n staging
          kubectl rollout status deployment/my-app -n staging --timeout=5m

  deploy-production:
    needs: deploy-staging
    runs-on: ubuntu-latest
    environment: production   # Requires manual approval

    steps:
      - uses: actions/checkout@v4

      - name: Configure kubectl
        uses: azure/k8s-set-context@v3
        with:
          method: kubeconfig
          kubeconfig: ${{ secrets.KUBE_CONFIG_PRODUCTION }}

      - name: Deploy to production
        run: |
          IMAGE_TAG="sha-${{ github.sha }}"
          kubectl set image deployment/my-app \
            my-app=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:$IMAGE_TAG \
            -n production
          kubectl rollout status deployment/my-app -n production --timeout=10m
```

---

## ArgoCD — GitOps Continuous Delivery

**ArgoCD** watches a Git repository and automatically syncs the cluster to match what's in Git.

```
Git repo (source of truth)
    ↓ ArgoCD watches for changes
    ↓ Detects drift between Git state and cluster state
    ↓ Auto-syncs cluster to match Git
Kubernetes cluster
```

### Install ArgoCD

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f \
  https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for pods
kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=argocd-server -n argocd

# Get admin password
kubectl get secret argocd-initial-admin-secret -n argocd \
  -o jsonpath='{.data.password}' | base64 -d

# Access UI
kubectl port-forward svc/argocd-server 8080:443 -n argocd
# Open: https://localhost:8080 (admin / <password above>)
```

### Create an ArgoCD Application

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: my-app
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/my-org/my-app-config
    targetRevision: HEAD
    path: k8s/production         # Path to K8s manifests in the repo

  destination:
    server: https://kubernetes.default.svc
    namespace: production

  syncPolicy:
    automated:
      prune: true                # Delete resources removed from Git
      selfHeal: true             # Revert manual changes to cluster
    syncOptions:
      - CreateNamespace=true
```

### GitOps Workflow

```
Developer pushes YAML change to Git
    ↓
ArgoCD detects difference
    ↓
ArgoCD applies change to cluster
    ↓
ArgoCD shows sync status: Synced ✓
```

### ArgoCD CLI

```bash
# Install ArgoCD CLI
curl -sSL -o argocd https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
sudo install argocd /usr/local/bin/argocd

# Login
argocd login localhost:8080

# List apps
argocd app list

# Sync an app manually
argocd app sync my-app

# Check sync status
argocd app get my-app

# Rollback to previous version
argocd app rollback my-app
```

---

## Blue-Green Deployment

Two identical environments: only one receives traffic at a time.

```yaml
# Blue deployment (current live)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app-blue
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
      version: blue
  template:
    metadata:
      labels:
        app: my-app
        version: blue
    spec:
      containers:
        - name: my-app
          image: my-app:1.0

---
# Green deployment (new version)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app-green
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
      version: green
  template:
    metadata:
      labels:
        app: my-app
        version: green
    spec:
      containers:
        - name: my-app
          image: my-app:2.0

---
# Service: switch traffic from blue to green
apiVersion: v1
kind: Service
metadata:
  name: my-app-service
spec:
  selector:
    app: my-app
    version: blue    # Change to 'green' to switch traffic
  ports:
    - port: 80
```

```bash
# Switch to green
kubectl patch service my-app-service \
  -p '{"spec":{"selector":{"version":"green"}}}'

# Rollback: switch back to blue
kubectl patch service my-app-service \
  -p '{"spec":{"selector":{"version":"blue"}}}'
```

---

## Canary Deployment

Send a small percentage of traffic to the new version:

```yaml
# Current: 3 replicas of v1 (75% traffic)
# Canary:  1 replica of v2 (25% traffic)

apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app-stable
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
        - name: my-app
          image: my-app:1.0
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app-canary
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
        - name: my-app
          image: my-app:2.0   # New version
---
apiVersion: v1
kind: Service
metadata:
  name: my-app-service
spec:
  selector:
    app: my-app     # Routes to BOTH deployments (both share this label)
  ports:
    - port: 80
```

---

## Docker Image Best Practices for CI/CD

```dockerfile
# Multi-stage build — keep final image small
FROM composer:2.7 AS composer
WORKDIR /app
COPY composer.json composer.lock ./
RUN composer install --no-dev --no-scripts --prefer-dist

FROM php:8.3-fpm-alpine AS production
WORKDIR /var/www/html

# Install runtime dependencies only
RUN apk add --no-cache nginx

COPY --from=composer /app/vendor ./vendor
COPY . .

# Non-root user
RUN adduser -D -u 1000 appuser
USER appuser

EXPOSE 9000
CMD ["php-fpm"]
```

---

## CI/CD Best Practices

- **Build once, deploy many:** Same image to staging then production (just change env)
- **Tag images with Git SHA:** `my-app:sha-abc1234` — immutable, traceable
- **Never use `latest` in production:** Immutable tags only
- **Rollback = re-deploy old tag:** Not a special process — just set the old image tag
- **Store Kubernetes manifests in Git:** Enable ArgoCD GitOps
- **Separate CI (test/build) from CD (deploy):** Different workflows, different secrets
- **Use environments + approvals:** Staging auto-deploys, production requires approval


---

# Chapter 17: Stateful Workloads

## Running Databases on Kubernetes

Running databases in Kubernetes requires careful planning. Key concerns:
- **Persistent storage** — data must survive pod restarts
- **Stable network identity** — replicas need to address each other
- **Ordered operations** — primary must start before replicas
- **Backup strategy** — automated, tested backups

> For production, prefer **managed cloud databases** (RDS, Cloud SQL) over self-hosted on K8s unless you have the operational expertise. Self-hosted is great for dev/staging.

---

## MySQL on Kubernetes

### Single Instance (Dev/Staging)

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mysql-secret
  namespace: default
type: Opaque
stringData:
  MYSQL_ROOT_PASSWORD: "rootpassword"
  MYSQL_DATABASE: "myapp"
  MYSQL_USER: "appuser"
  MYSQL_PASSWORD: "apppassword"
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pvc
spec:
  accessModes: [ReadWriteOnce]
  storageClassName: standard
  resources:
    requests:
      storage: 20Gi
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  serviceName: mysql-headless
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
        - name: mysql
          image: mysql:8.0
          envFrom:
            - secretRef:
                name: mysql-secret
          ports:
            - containerPort: 3306
          volumeMounts:
            - name: data
              mountPath: /var/lib/mysql
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "2Gi"
              cpu: "1"
          livenessProbe:
            exec:
              command: ["mysqladmin", "ping", "-h", "localhost"]
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            exec:
              command: ["mysql", "-uroot", "-p$(MYSQL_ROOT_PASSWORD)", "-e", "SELECT 1"]
            initialDelaySeconds: 10
            periodSeconds: 5
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: mysql-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: mysql-headless
spec:
  clusterIP: None
  selector:
    app: mysql
  ports:
    - port: 3306
---
apiVersion: v1
kind: Service
metadata:
  name: mysql-service
spec:
  selector:
    app: mysql
  ports:
    - port: 3306
      targetPort: 3306
```

### MySQL Backup CronJob

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: mysql-backup
spec:
  schedule: "0 3 * * *"   # 3 AM daily
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: backup
              image: mysql:8.0
              command:
                - /bin/sh
                - -c
                - |
                  mysqldump -h mysql-service -u root -p$MYSQL_ROOT_PASSWORD --all-databases | \
                  gzip > /backup/full-$(date +%Y%m%d-%H%M%S).sql.gz
                  echo "Backup done"
              env:
                - name: MYSQL_ROOT_PASSWORD
                  valueFrom:
                    secretKeyRef:
                      name: mysql-secret
                      key: MYSQL_ROOT_PASSWORD
              volumeMounts:
                - name: backup
                  mountPath: /backup
          volumes:
            - name: backup
              persistentVolumeClaim:
                claimName: mysql-backup-pvc
```

---

## Redis on Kubernetes

### Redis Single Instance

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
        - name: redis
          image: redis:7-alpine
          command: ["redis-server", "--appendonly", "yes", "--requirepass", "$(REDIS_PASSWORD)"]
          env:
            - name: REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: redis-secret
                  key: password
          ports:
            - containerPort: 6379
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "1Gi"
              cpu: "500m"
          volumeMounts:
            - name: data
              mountPath: /data
          livenessProbe:
            exec:
              command: ["redis-cli", "ping"]
            periodSeconds: 10
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: redis-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: redis-service
spec:
  selector:
    app: redis
  ports:
    - port: 6379
```

### Redis with Helm (Production-ready)

```bash
helm install redis bitnami/redis \
  --set auth.password=secretpassword \
  --set master.persistence.size=10Gi \
  --set replica.replicaCount=2 \
  -n production
```

---

## Elasticsearch on Kubernetes

```bash
# Install ECK (Elastic Cloud on Kubernetes) operator
kubectl create -f https://download.elastic.co/downloads/eck/2.11.1/crds.yaml
kubectl apply -f https://download.elastic.co/downloads/eck/2.11.1/operator.yaml
```

```yaml
apiVersion: elasticsearch.k8s.elastic.co/v1
kind: Elasticsearch
metadata:
  name: my-elasticsearch
  namespace: elastic
spec:
  version: 8.12.0
  nodeSets:
    - name: default
      count: 3
      config:
        node.store.allow_mmap: false
      podTemplate:
        spec:
          containers:
            - name: elasticsearch
              resources:
                requests:
                  memory: 2Gi
                  cpu: 1
                limits:
                  memory: 4Gi
                  cpu: 2
      volumeClaimTemplates:
        - metadata:
            name: elasticsearch-data
          spec:
            accessModes: [ReadWriteOnce]
            storageClassName: fast-ssd
            resources:
              requests:
                storage: 50Gi
---
# Kibana
apiVersion: kibana.k8s.elastic.co/v1
kind: Kibana
metadata:
  name: my-kibana
  namespace: elastic
spec:
  version: 8.12.0
  count: 1
  elasticsearchRef:
    name: my-elasticsearch
```

---

## Kubernetes Operators

An **Operator** extends Kubernetes with domain knowledge for a specific application. It manages complex stateful apps (like databases) using custom controllers and CRDs.

| Operator | Manages |
|---|---|
| ECK | Elasticsearch + Kibana |
| MySQL Operator | MySQL InnoDB Cluster |
| Postgres Operator (Zalando) | PostgreSQL + replication |
| Redis Operator (OpsTree) | Redis cluster |
| Strimzi | Apache Kafka |
| Prometheus Operator | Prometheus + Alertmanager |

### Install PostgreSQL with Zalando Operator

```bash
helm install postgres-operator postgres-operator-charts/postgres-operator \
  -n postgres-operator --create-namespace
```

```yaml
apiVersion: acid.zalan.do/v1
kind: postgresql
metadata:
  name: my-postgres
  namespace: default
spec:
  teamId: "myteam"
  volume:
    size: 20Gi
  numberOfInstances: 2    # 1 primary + 1 replica
  users:
    appuser:
      - superuser
      - createdb
  databases:
    myapp: appuser
  postgresql:
    version: "16"
  resources:
    requests:
      cpu: 100m
      memory: 100Mi
    limits:
      cpu: 500m
      memory: 500Mi
```

---

## Production Recommendations

| Workload | Dev/Staging | Production |
|---|---|---|
| MySQL | StatefulSet in K8s | AWS RDS / Google Cloud SQL |
| PostgreSQL | Zalando Operator | AWS RDS / Cloud SQL |
| Redis | Deployment in K8s | AWS ElastiCache / Redis Cloud |
| Elasticsearch | ECK operator | AWS OpenSearch / Elastic Cloud |
| Kafka | Strimzi operator | AWS MSK / Confluent Cloud |

**Why managed services for production?**
- Automated backups
- Automated failover
- Security patches
- No operational overhead
- SLA guarantees

---

## Shared Volumes for Application Files

For apps needing shared file storage across pods (e.g., CMS uploads):

```bash
# AWS: Install EFS CSI Driver
helm install aws-efs-csi-driver \
  aws-efs-csi-driver/aws-efs-csi-driver \
  -n kube-system
```

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: efs-sc
provisioner: efs.csi.aws.com
parameters:
  provisioningMode: efs-ap
  fileSystemId: fs-xxxxxxxxx
  directoryPerms: "700"
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: shared-uploads
spec:
  storageClassName: efs-sc
  accessModes: [ReadWriteMany]   # Multiple pods can write
  resources:
    requests:
      storage: 100Gi
```


---

# Chapter 18: Debugging & Troubleshooting

## The Debugging Mindset

When something breaks in Kubernetes, work top-down:

```
1. What is the symptom?   (app down, slow, 502, pod crash)
2. Check the pod status   (kubectl get pods)
3. Describe the pod       (kubectl describe pod)
4. Check logs             (kubectl logs)
5. Check events           (kubectl get events)
6. Check service/ingress  (kubectl describe svc/ingress)
7. Check resource limits  (kubectl top pods/nodes)
8. Exec into the pod      (kubectl exec -it ... -- sh)
```

---

## Common Pod States

| State | Cause | Debug Step |
|---|---|---|
| `Pending` | No node available, PVC not bound | `kubectl describe pod` → Events |
| `CrashLoopBackOff` | Container keeps crashing | `kubectl logs pod --previous` |
| `OOMKilled` | Container exceeded memory limit | Increase limits or fix memory leak |
| `ImagePullBackOff` | Can't pull the image | Check image name, registry auth |
| `ErrImagePull` | First image pull failure | Same as above |
| `CreateContainerConfigError` | Bad config (missing ConfigMap/Secret) | `kubectl describe pod` → Events |
| `Error` | Container exited with non-zero | `kubectl logs` to see error |
| `Terminating` (stuck) | Finalizers blocking deletion | Remove finalizers |

---

## CrashLoopBackOff

```bash
# Step 1: Check current logs
kubectl logs my-pod

# Step 2: Check logs from the PREVIOUS (crashed) container
kubectl logs my-pod --previous

# Step 3: Describe to see last state + exit code
kubectl describe pod my-pod
# Look for:
# Last State: Terminated
#   Reason: Error / OOMKilled
#   Exit Code: 1 (or 137 for OOM)

# Step 4: Test your container command locally
docker run --rm -e DATABASE_URL=... my-app:latest
```

---

## ImagePullBackOff / ErrImagePull

```bash
kubectl describe pod my-pod
# Look for Events section:
# Failed to pull image "registry/my-app:latest": ...

# Common causes and fixes:
# 1. Wrong image name/tag
kubectl set image deployment/my-app my-app=registry/my-app:v1.2.3

# 2. Private registry — need imagePullSecrets
kubectl create secret docker-registry regcred \
  --docker-server=registry.example.com \
  --docker-username=myuser \
  --docker-password=mypassword

# Add to deployment:
spec:
  imagePullSecrets:
    - name: regcred

# 3. For kind: load image manually
kind load docker-image my-app:latest --name my-cluster

# 4. For minikube: load image manually
minikube image load my-app:latest
```

---

## OOMKilled

Container exceeded its memory limit — kernel kills it with exit code 137.

```bash
kubectl describe pod my-pod
# Look for:
# Last State: Terminated
#   Reason: OOMKilled
#   Exit Code: 137

# Fix: increase memory limits
kubectl set resources deployment/my-app \
  --limits=memory=1Gi

# Or edit the deployment
kubectl edit deployment my-app
# Change limits.memory to a higher value

# Find actual usage to set appropriate limits
kubectl top pod my-pod
```

---

## Pending Pods

```bash
kubectl describe pod my-pod
# Look at Events section

# Common causes:
# 1. Insufficient resources
kubectl describe nodes | grep -A5 "Allocated resources"

# 2. PVC not bound
kubectl get pvc
kubectl describe pvc my-pvc

# 3. Node selector doesn't match
kubectl get nodes --show-labels
kubectl describe pod my-pod | grep -A5 "Node-Selectors"

# 4. Taints blocking scheduling
kubectl describe nodes | grep Taints

# 5. No nodes in the cluster
kubectl get nodes
```

---

## Service Not Routing Traffic

```bash
# Step 1: Check if pods are Running
kubectl get pods -l app=my-app

# Step 2: Check service
kubectl get svc my-service
kubectl describe svc my-service

# Step 3: Check endpoints (are pods being selected?)
kubectl get endpoints my-service
# If ENDPOINTS is <none>: label selector doesn't match any pods

# Step 4: Verify labels match
kubectl get pods --show-labels
# Compare with service selector

# Step 5: Test service DNS from inside cluster
kubectl run test --image=busybox --rm -it --restart=Never -- \
  nslookup my-service
kubectl run test --image=busybox --rm -it --restart=Never -- \
  wget -qO- http://my-service

# Step 6: Test directly on pod port
kubectl port-forward pod/my-pod 8080:8080
curl http://localhost:8080
```

---

## Ingress Not Working

```bash
# Check ingress resource
kubectl describe ingress my-ingress

# Check if ingress controller is running
kubectl get pods -n ingress-nginx

# Check ingress controller logs
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller

# Verify backend service exists
kubectl get svc my-backend-service

# Check if backend pods are healthy
kubectl get endpoints my-backend-service

# Test path routing
curl -H "Host: myapp.example.com" http://<ingress-ip>/api
```

---

## Persistent Volume Issues

```bash
# PVC stuck in Pending
kubectl describe pvc my-pvc
# Events: no matching StorageClass, no PV available

# List storage classes
kubectl get sc

# Check if PV is Available
kubectl get pv

# Force delete stuck PVC (remove finalizers)
kubectl patch pvc my-pvc -p '{"metadata":{"finalizers":null}}'

# Check mounted volume inside pod
kubectl exec -it my-pod -- ls -la /data
kubectl exec -it my-pod -- df -h
```

---

## Debug Container (kubectl debug)

```bash
# Add a debug container to a running pod (K8s 1.23+)
kubectl debug -it my-pod --image=busybox --target=my-container

# Create a copy of a pod with debug tools
kubectl debug my-pod -it --copy-to=debug-pod --image=ubuntu:22.04

# Debug a node (runs privileged pod on the node)
kubectl debug node/my-node -it --image=ubuntu:22.04
```

---

## Common Application Debugging

### Test DNS Resolution

```bash
kubectl run dns-test --image=busybox --rm -it --restart=Never -- \
  nslookup kubernetes.default.svc.cluster.local
```

### Test HTTP Connectivity

```bash
kubectl run curl-test --image=curlimages/curl --rm -it --restart=Never -- \
  curl -v http://my-service:80/health
```

### Check Environment Variables

```bash
kubectl exec my-pod -- env | sort
```

### Inspect Mounted Files

```bash
kubectl exec my-pod -- ls -la /etc/config/
kubectl exec my-pod -- cat /etc/config/app.conf
```

---

## Node Debugging

```bash
# List nodes
kubectl get nodes

# Node details
kubectl describe node my-node

# Node resource usage
kubectl top node my-node

# Check node conditions
kubectl get nodes -o json | jq '.items[].status.conditions'

# SSH into a node (if you have access)
ssh user@<node-ip>

# Or use kubectl debug node
kubectl debug node/my-node -it --image=ubuntu:22.04
```

---

## Stuck Terminating Pods/Namespaces

```bash
# Pod stuck Terminating
kubectl delete pod my-pod --grace-period=0 --force

# Namespace stuck Terminating — remove finalizers
kubectl get namespace my-namespace -o json \
  | jq '.spec.finalizers = []' \
  | kubectl replace --raw "/api/v1/namespaces/my-namespace/finalize" -f -
```

---

## Debugging Cheatsheet

```bash
# Pod basics
kubectl get pods -A                          # All pods, all namespaces
kubectl get pods -o wide                     # Show node and IP
kubectl describe pod <pod>                   # Full details + events
kubectl logs <pod>                           # Logs
kubectl logs <pod> --previous                # Previous container logs
kubectl logs <pod> -c <container>            # Specific container
kubectl exec -it <pod> -- sh                 # Shell into container
kubectl top pod <pod>                        # CPU/memory usage

# Service/network
kubectl get svc,ep -n <namespace>            # Services + endpoints
kubectl describe svc <service>               # Service details
kubectl port-forward svc/<svc> 8080:80       # Forward to localhost

# Events
kubectl get events -n <namespace> --sort-by='.lastTimestamp'

# Nodes
kubectl get nodes                            # Node status
kubectl describe node <node>                 # Node details + events
kubectl top nodes                            # Node resource usage
kubectl drain <node> --ignore-daemonsets     # Drain node for maintenance

# YAML output
kubectl get pod <pod> -o yaml                # Full YAML
kubectl get all -n <namespace> -o yaml       # All resources as YAML
```


---

# Chapter 19: Cloud Providers — GKE, EKS, AKS

## Why Use Managed Kubernetes?

Running your own Kubernetes control plane is complex. Managed Kubernetes services offload control plane management:

| Feature | Self-managed | Managed (GKE/EKS/AKS) |
|---|---|---|
| Control plane setup | Manual | Automatic |
| Control plane upgrades | Manual | 1-click or automatic |
| etcd backups | Manual | Automatic |
| HA control plane | Complex setup | Built-in |
| Cost | Node + ops cost | Node cost only |
| SLA | Your responsibility | 99.95%+ SLA |

---

## Google Kubernetes Engine (GKE)

### Create a GKE Cluster

```bash
# Install gcloud CLI
# https://cloud.google.com/sdk/docs/install

# Authenticate
gcloud auth login
gcloud config set project my-project-id

# Create cluster (Autopilot — Google manages nodes too)
gcloud container clusters create-auto my-cluster \
  --region=us-central1

# Create standard cluster
gcloud container clusters create my-cluster \
  --region=us-central1 \
  --num-nodes=3 \
  --machine-type=e2-standard-2 \
  --disk-size=50GB \
  --enable-autoscaling \
  --min-nodes=1 \
  --max-nodes=10

# Get credentials for kubectl
gcloud container clusters get-credentials my-cluster --region=us-central1
kubectl get nodes
```

### Node Pools

```bash
# Add a GPU node pool
gcloud container node-pools create gpu-pool \
  --cluster=my-cluster \
  --region=us-central1 \
  --machine-type=n1-standard-4 \
  --accelerator=type=nvidia-tesla-t4,count=1 \
  --num-nodes=1

# Add a high-memory pool
gcloud container node-pools create high-mem \
  --cluster=my-cluster \
  --region=us-central1 \
  --machine-type=n2-highmem-8 \
  --num-nodes=2
```

### GKE Storage

```yaml
# GKE standard persistent disk
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: pd.csi.storage.gke.io
parameters:
  type: pd-ssd
reclaimPolicy: Retain
allowVolumeExpansion: true
---
# GKE Filestore for ReadWriteMany
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: filestore-rwx
provisioner: filestore.csi.storage.gke.io
parameters:
  tier: standard
  network: default
volumeBindingMode: WaitForFirstConsumer
```

### GKE Useful Commands

```bash
# List clusters
gcloud container clusters list

# Upgrade cluster to new K8s version
gcloud container clusters upgrade my-cluster --master --cluster-version=1.28

# Upgrade node pool
gcloud container clusters upgrade my-cluster \
  --node-pool=default-pool \
  --cluster-version=1.28

# Delete cluster
gcloud container clusters delete my-cluster --region=us-central1
```

---

## Amazon Elastic Kubernetes Service (EKS)

### Create an EKS Cluster

```bash
# Install eksctl
curl --silent --location "https://github.com/eksfargate/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin

# Create cluster
eksctl create cluster \
  --name my-cluster \
  --region us-east-1 \
  --nodegroup-name standard-workers \
  --node-type t3.medium \
  --nodes 3 \
  --nodes-min 1 \
  --nodes-max 10 \
  --managed

# Get credentials
aws eks update-kubeconfig --name my-cluster --region us-east-1
kubectl get nodes
```

### EKS with eksctl Config File

```yaml
# cluster.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: my-cluster
  region: us-east-1
  version: "1.28"

managedNodeGroups:
  - name: general
    instanceType: t3.medium
    minSize: 2
    maxSize: 10
    desiredCapacity: 3
    volumeSize: 50

  - name: spot
    instanceTypes: [t3.medium, t3.large, t3a.medium]
    spot: true
    minSize: 0
    maxSize: 5
```

```bash
eksctl create cluster -f cluster.yaml
```

### EKS Storage (EBS + EFS)

```bash
# Install EBS CSI Driver
eksctl create addon --name aws-ebs-csi-driver --cluster my-cluster

# Install EFS CSI Driver (for ReadWriteMany)
helm install aws-efs-csi-driver aws-efs-csi-driver/aws-efs-csi-driver \
  -n kube-system
```

```yaml
# EBS StorageClass for RWO
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: ebs-gp3
provisioner: ebs.csi.aws.com
parameters:
  type: gp3
  iopsPerGB: "3000"
  throughput: "125"
reclaimPolicy: Retain
allowVolumeExpansion: true
---
# EFS StorageClass for RWX
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: efs-sc
provisioner: efs.csi.aws.com
parameters:
  provisioningMode: efs-ap
  fileSystemId: fs-xxxxxxxxx
  directoryPerms: "700"
```

### IAM Roles for Service Accounts (IRSA)

Give pods access to AWS services (S3, SQS, etc.) without hardcoding credentials:

```bash
# Create OIDC provider
eksctl utils associate-iam-oidc-provider --cluster=my-cluster --approve

# Create IAM service account
eksctl create iamserviceaccount \
  --name=s3-access-sa \
  --namespace=production \
  --cluster=my-cluster \
  --attach-policy-arn=arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess \
  --approve
```

---

## Azure Kubernetes Service (AKS)

### Create an AKS Cluster

```bash
# Install Azure CLI
# https://docs.microsoft.com/en-us/cli/azure/install-azure-cli

az login
az account set --subscription my-subscription-id

# Create resource group
az group create --name my-rg --location eastus

# Create AKS cluster
az aks create \
  --resource-group my-rg \
  --name my-cluster \
  --node-count 3 \
  --node-vm-size Standard_D2s_v3 \
  --enable-cluster-autoscaler \
  --min-count 1 \
  --max-count 10 \
  --generate-ssh-keys

# Get credentials
az aks get-credentials --resource-group my-rg --name my-cluster
kubectl get nodes
```

### AKS Node Pools

```bash
# Add node pool
az aks nodepool add \
  --resource-group my-rg \
  --cluster-name my-cluster \
  --name highperformance \
  --node-count 2 \
  --node-vm-size Standard_D8s_v3 \
  --enable-cluster-autoscaler \
  --min-count 1 \
  --max-count 5
```

### AKS Storage (Azure Disk + Azure Files)

```yaml
# Azure Disk (RWO)
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: azure-ssd
provisioner: disk.csi.azure.com
parameters:
  skuName: Premium_LRS
allowVolumeExpansion: true
---
# Azure Files (RWX)
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: azure-files-rwx
provisioner: file.csi.azure.com
parameters:
  skuName: Standard_LRS
mountOptions:
  - dir_mode=0777
  - file_mode=0777
```

---

## Cross-Provider Comparison

| Feature | GKE | EKS | AKS |
|---|---|---|---|
| Managed control plane | Yes | Yes (free) | Yes |
| Autopilot mode | Yes (GKE Autopilot) | Fargate | Virtual Nodes |
| Default storage | pd-standard / pd-ssd | gp2 / gp3 (EBS) | azuredisk |
| RWX storage | Filestore | EFS | Azure Files |
| IAM integration | Workload Identity | IRSA | AAD Pod Identity |
| Cluster upgrade | Auto-upgrade available | Manual (eksctl) | Auto-upgrade |
| Networking | VPC-native | VPC CNI | Azure CNI |
| Free control plane | Yes (standard) | Yes | Yes |
| Container registry | Artifact Registry | ECR | ACR |

---

## Multi-Cloud Best Practices

```bash
# Use kubeconfig contexts to manage multiple clusters
kubectl config get-contexts
kubectl config use-context gke_my-project_us-central1_my-cluster
kubectl config use-context arn:aws:eks:us-east-1:123:cluster/my-cluster
kubectl config use-context my-cluster

# Give contexts friendly names
kubectl config rename-context gke_my-project_us-central1_my-cluster gke-prod
kubectl config rename-context arn:aws:eks:us-east-1:123:cluster/my-cluster eks-staging
```

---

## Cost Optimization Tips

```bash
# Use spot/preemptible nodes for batch workloads (60-90% cheaper)
# GKE: add spot node pool
gcloud container node-pools create spot-pool \
  --cluster=my-cluster \
  --spot \
  --machine-type=e2-standard-4

# EKS
eksctl create nodegroup --spot --instance-types=m5.large,m5.xlarge

# Right-size nodes using VPA recommendations
kubectl describe vpa my-app-vpa

# Use Cluster Autoscaler to scale down idle nodes

# Use namespace resource quotas to prevent runaway usage

# Set pod disruption budgets to enable safe node drain
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: my-app-pdb
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: my-app
```


---

# Chapter 20: Learning Roadmap — 5 Weeks to Kubernetes Proficiency

## Overview

A structured, hands-on 5-week plan that takes you from zero → running production-grade workloads, with concrete daily tasks and clear success criteria.

---

## Week 1: Foundations & First Deployment

**Goal:** Run a cluster locally, understand pods and services, deploy your first app.

| Day | Task |
|---|---|
| Day 1 | Install kubectl + Minikube. Run `minikube start`. Deploy nginx with `kubectl create deployment`. |
| Day 2 | Write your first Pod YAML. Apply it. Get logs, exec into it, delete it. |
| Day 3 | Write a Deployment YAML (3 replicas). Scale it. Update the image. Watch rolling update. |
| Day 4 | Create a ClusterIP Service. Connect to it from another pod. Try NodePort. |
| Day 5 | Create a ConfigMap and a Secret. Load them as env vars in a Deployment. |
| Day 6–7 | Deploy a full stack: Nginx + MySQL (each with a Service). Connect one to the other. |

**Success Criteria:**
- Cluster running locally
- Can apply, scale, update, and rollback a Deployment
- Understands Pod lifecycle and pod status debugging
- ConfigMap/Secret mounted in a running pod

---

## Week 2: Networking & Storage

**Goal:** Route external traffic, persist data, understand namespace isolation.

| Day | Task |
|---|---|
| Day 1–2 | Install ingress-nginx. Create an Ingress for your app. Test path-based and host-based routing. |
| Day 3 | Install cert-manager. Create a ClusterIssuer. Add TLS to your Ingress. |
| Day 4 | Create a PersistentVolumeClaim. Mount it in a MySQL StatefulSet. Verify data persists after pod restart. |
| Day 5 | Create `dev` and `staging` namespaces. Deploy app in each. Verify cross-namespace DNS. |
| Day 6 | Add ResourceQuota and LimitRange to each namespace. |
| Day 7 | Implement a NetworkPolicy: deny all → allow only frontend→backend. |

**Success Criteria:**
- External traffic reaches app via Ingress with TLS
- Database data survives pod restart
- Namespaces with quotas in place
- NetworkPolicy restricting pod communication

---

## Week 3: RBAC, Helm, and Autoscaling

**Goal:** Secure the cluster, package your app, automate scaling.

| Day | Task |
|---|---|
| Day 1–2 | Create a ServiceAccount. Bind a read-only Role to it. Test permissions with `kubectl auth can-i`. |
| Day 3 | Install Helm. Deploy MySQL with `helm install bitnami/mysql`. Upgrade it with new values. |
| Day 4 | Run `helm create my-app`. Convert your Deployment into a Helm chart with values.yaml. |
| Day 5 | Install metrics-server. Set up HPA for your app (target 70% CPU). |
| Day 6 | Generate load with a busybox loop. Watch HPA scale pods up. Stop load, watch scale down. |
| Day 7 | Set up RBAC: one ServiceAccount per namespace with least-privilege permissions. |

**Success Criteria:**
- RBAC with least-privilege ServiceAccounts
- App packaged as a Helm chart
- HPA scaling in response to load
- `kubectl auth can-i` returns correct results

---

## Week 4: Monitoring, Security & CI/CD

**Goal:** Observe what's happening, harden security, automate deployments.

| Day | Task |
|---|---|
| Day 1–2 | Install kube-prometheus-stack. Open Grafana. Import K8s cluster dashboard. |
| Day 3 | Create a PrometheusRule: alert when pod restarts > 0. Configure Alertmanager to print to a log. |
| Day 4 | Apply PodSecurity labels on namespaces (`restricted` mode). Fix any policy violations. |
| Day 5 | Add SecurityContext to your app: `runAsNonRoot`, `readOnlyRootFilesystem`, drop all capabilities. |
| Day 6 | Write a GitHub Actions workflow: test → build Docker image → push to GHCR → deploy to cluster. |
| Day 7 | Install ArgoCD. Create an Application pointing to your Helm chart repo. Watch GitOps auto-sync. |

**Success Criteria:**
- Grafana dashboards showing cluster health
- Alert fires on pod restart
- Pods passing `restricted` Pod Security Standard
- GitHub Actions CI/CD pipeline deploys to cluster

---

## Week 5: Production Readiness & Cloud

**Goal:** Production hardening, cloud deployment, advanced patterns.

| Day | Task |
|---|---|
| Day 1 | Create a cloud cluster (GKE free tier or EKS). Connect kubectl to it. Migrate your app. |
| Day 2 | Set up cloud-native storage (GKE persistent disk or EBS). Test persistence on cloud. |
| Day 3 | Implement blue-green deployment for your app. Test traffic switching. |
| Day 4 | Implement a CronJob for nightly database backup. Trigger it manually and verify backup file. |
| Day 5 | Install Loki + Grafana. Ship container logs to Loki. Query logs in Grafana. |
| Day 6 | Run kube-bench for CIS compliance checks. Fix at least 5 findings. |
| Day 7 | Full DR drill: delete your deployment, restore from backup, verify app recovers cleanly. |

**Success Criteria:**
- App running on real cloud cluster
- Blue-green deployment working
- Automated database backups running via CronJob
- Logs queryable in Grafana/Loki
- kube-bench score improved

---

## Summary Timeline

| Week | Focus | Deliverable |
|---|---|---|
| Week 1 | Pods, Deployments, Services | App running locally |
| Week 2 | Networking, Storage, Namespaces | Full stack with TLS + persistent DB |
| Week 3 | RBAC, Helm, Autoscaling | App as Helm chart + HPA working |
| Week 4 | Monitoring, Security, CI/CD | GitOps pipeline + Grafana dashboards |
| Week 5 | Cloud + Production hardening | Production-grade cloud deployment |

---

## Quick Reference: Most-Used Commands

```bash
# ---- Cluster ----
kubectl cluster-info
kubectl get nodes -o wide
kubectl top nodes

# ---- Pods ----
kubectl get pods -A
kubectl describe pod <pod>
kubectl logs <pod> -f
kubectl logs <pod> --previous
kubectl exec -it <pod> -- sh
kubectl top pod <pod>

# ---- Deployments ----
kubectl get deployments
kubectl scale deployment <name> --replicas=5
kubectl set image deployment/<name> <container>=<image>:<tag>
kubectl rollout status deployment/<name>
kubectl rollout history deployment/<name>
kubectl rollout undo deployment/<name>

# ---- Services ----
kubectl get svc,ep
kubectl port-forward svc/<name> 8080:80

# ---- Debugging ----
kubectl get events --sort-by='.lastTimestamp'
kubectl describe <resource> <name>
kubectl get all -n <namespace>
kubectl run debug --image=busybox --rm -it --restart=Never -- sh

# ---- Helm ----
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install <release> <chart> -f values.yaml
helm upgrade <release> <chart> -f values.yaml
helm rollback <release>
helm list

# ---- Config/Secrets ----
kubectl create configmap <name> --from-literal=key=value
kubectl create secret generic <name> --from-literal=key=value
kubectl get secret <name> -o jsonpath='{.data.key}' | base64 -d

# ---- RBAC ----
kubectl auth can-i <verb> <resource>
kubectl auth can-i get pods --as=user@example.com

# ---- Storage ----
kubectl get pv,pvc
kubectl describe pvc <name>
```

---

## Recommended Resources

| Resource | URL |
|---|---|
| Official K8s Docs | `kubernetes.io/docs` |
| Kubernetes Playground | `labs.play-with-k8s.com` |
| Helm Charts | `artifacthub.io` |
| CNCF Landscape | `landscape.cncf.io` |
| kube.academy (free) | `kube.academy` |
| Killer.sh CKA Practice | `killer.sh` |
| The Kubernetes Book | `thekubernetesbook.com` |

---

## Certifications Path

| Certification | Level | Focus |
|---|---|---|
| **KCNA** | Beginner | Cloud Native fundamentals |
| **CKA** | Intermediate | Cluster administration |
| **CKAD** | Intermediate | Application development |
| **CKS** | Advanced | Security specialist |

**Study tips for CKA/CKAD:**
- All exams are hands-on (open book, real cluster)
- Master `kubectl` — speed matters
- Practice on `killer.sh` (exam simulator)
- Read `kubernetes.io/docs` — that's all you get during the exam
