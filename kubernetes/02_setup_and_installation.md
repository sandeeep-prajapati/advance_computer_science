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
