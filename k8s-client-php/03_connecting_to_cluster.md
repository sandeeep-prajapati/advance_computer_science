# Chapter 3: Connecting to a Cluster

## Connection Architecture

```
Your PHP App
     │
     │ HTTPS (port 6443)
     ▼
┌─────────────────────────┐
│   Kubernetes API Server │
│                         │
│  Authenticates via:     │
│  - Client certificate   │
│  - Bearer token         │
│  - Basic auth (legacy)  │
└─────────────────────────┘
```

---

## Local Cluster (Minikube / kind / k3s)

### Minikube

```php
use Maclof\Kubernetes\Client;

// Option A: Auto-detect from kubeconfig
$client = Client::fromKubeConfig();

// Option B: Explicit minikube certs
$client = new Client([
    'master'      => shell_exec('minikube ip') . ':8443',
    'ca_cert'     => $_SERVER['HOME'] . '/.minikube/ca.crt',
    'client_cert' => $_SERVER['HOME'] . '/.minikube/profiles/minikube/client.crt',
    'client_key'  => $_SERVER['HOME'] . '/.minikube/profiles/minikube/client.key',
]);
```

### kind (Kubernetes in Docker)

```bash
# Create cluster
kind create cluster --name dev

# Get config
kind get kubeconfig --name dev > /tmp/kind-config.yaml
```

```php
$client = Client::fromKubeConfig(
    kubeConfigPath: '/tmp/kind-config.yaml'
);
```

### k3s

```bash
# k3s writes kubeconfig here
sudo cat /etc/rancher/k3s/k3s.yaml
# Copy to user path
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $USER ~/.kube/config
```

```php
$client = Client::fromKubeConfig();  // reads ~/.kube/config automatically
```

---

## Remote / Cloud Cluster

### GKE (Google Kubernetes Engine)

```bash
# Authenticate and get credentials
gcloud container clusters get-credentials my-cluster --region us-central1

# kubeconfig is updated automatically
kubectl config current-context   # gke_project_region_cluster
```

```php
$client = Client::fromKubeConfig(
    context: 'gke_my-project_us-central1_my-cluster'
);
```

### EKS (AWS Elastic Kubernetes Service)

```bash
aws eks update-kubeconfig --region us-east-1 --name my-cluster
```

```php
// EKS uses AWS IAM token — use kubectl proxy for PHP simplicity
// OR use aws-iam-authenticator token

$token = shell_exec('aws eks get-token --cluster-name my-cluster --query "status.token" -o text');

$client = new Client([
    'master'  => 'https://ABCDEF.gr7.us-east-1.eks.amazonaws.com',
    'ca_cert' => '/path/to/eks-ca.crt',
    'token'   => trim($token),
]);
```

### AKS (Azure Kubernetes Service)

```bash
az aks get-credentials --resource-group myRG --name my-cluster
```

```php
$client = Client::fromKubeConfig(
    context: 'my-cluster'
);
```

---

## In-Cluster Connection (Pod Running Inside K8s)

When your PHP app runs as a Pod inside Kubernetes, Kubernetes automatically injects credentials:

```yaml
# Your Pod automatically gets:
# /var/run/secrets/kubernetes.io/serviceaccount/token
# /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
# /var/run/secrets/kubernetes.io/serviceaccount/namespace
```

```php
function createInClusterClient(): Client
{
    $tokenPath = '/var/run/secrets/kubernetes.io/serviceaccount/token';
    $caPath    = '/var/run/secrets/kubernetes.io/serviceaccount/ca.crt';
    $nsPath    = '/var/run/secrets/kubernetes.io/serviceaccount/namespace';

    return new Client([
        'master'    => 'https://kubernetes.default.svc',
        'ca_cert'   => $caPath,
        'token'     => file_get_contents($tokenPath),
        'namespace' => file_get_contents($nsPath),
    ]);
}
```

---

## Context Switching

Work with multiple clusters or namespaces in one PHP script:

```php
use Maclof\Kubernetes\Client;

// Load two different contexts
$dev  = Client::fromKubeConfig(context: 'dev-cluster');
$prod = Client::fromKubeConfig(context: 'prod-cluster');

// Query both
$devPods  = $dev->pods()->find();
$prodPods = $prod->pods()->find();

echo "Dev pods:  " . count($devPods)  . "\n";
echo "Prod pods: " . count($prodPods) . "\n";
```

---

## Namespace Switching at Runtime

```php
$client = Client::fromKubeConfig();

$namespaces = ['default', 'kube-system', 'production', 'staging'];

foreach ($namespaces as $ns) {
    $client->setNamespace($ns);
    $pods = $client->pods()->find();
    echo "[$ns] Pods: " . count($pods) . "\n";
}
```

---

## Using kubectl proxy (Simplest Dev Setup)

Avoid certificate complexity by proxying through kubectl:

```bash
# Start kubectl proxy (runs on localhost:8001)
kubectl proxy &
```

```php
// No auth needed — kubectl proxy handles it
$client = new Client([
    'master'    => 'http://127.0.0.1:8001',
    'verify'    => false,
]);

$pods = $client->pods()->find();
```

> Only use this for local development. Never expose `kubectl proxy` in production.

---

## Disabling TLS Verification (Dev Only)

```php
$client = new Client([
    'master' => 'https://...',
    'verify' => false,   // NEVER use in production
    'token'  => '...',
]);
```

---

## Connection Troubleshooting

| Error | Likely Cause | Fix |
|---|---|---|
| `Connection refused` | Wrong port or API server down | Check `kubectl cluster-info` |
| `SSL certificate error` | Wrong CA cert | Verify `ca_cert` path |
| `401 Unauthorized` | Invalid/expired token | Refresh token |
| `403 Forbidden` | Insufficient RBAC permissions | Grant role to ServiceAccount |
| `Timeout` | Network issue / firewall | Check cluster network access |

```php
// Test connection with verbose error reporting
try {
    $version = $client->getVersion();
    echo "K8s version: " . $version['gitVersion'] . "\n";
} catch (\GuzzleHttp\Exception\ConnectException $e) {
    echo "Cannot connect to cluster: " . $e->getMessage() . "\n";
} catch (\GuzzleHttp\Exception\ClientException $e) {
    echo "HTTP error " . $e->getResponse()->getStatusCode() . ": "
        . $e->getResponse()->getBody() . "\n";
}
```
