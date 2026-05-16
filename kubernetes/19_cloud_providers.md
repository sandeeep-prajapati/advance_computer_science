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
