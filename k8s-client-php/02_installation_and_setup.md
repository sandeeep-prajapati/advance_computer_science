# Chapter 2: Installation & Setup

## Requirements

- PHP >= 7.4 (8.x recommended)
- Composer
- A running Kubernetes cluster (Minikube, k3s, kind, or cloud-managed)
- kubectl configured (for getting credentials)

---

## Installing maclof/kubernetes-client

```bash
composer require maclof/kubernetes-client
```

```bash
# Verify installation
composer show maclof/kubernetes-client
```

**Key dependencies it installs:**
- `guzzlehttp/guzzle` — HTTP client
- `symfony/yaml` — YAML parsing

---

## Installing renoki-co/php-k8s (Laravel alternative)

```bash
composer require renoki-co/php-k8s
```

For Laravel service provider:
```bash
php artisan vendor:publish --provider="RenokiCo\PhpK8s\PhpK8sServiceProvider"
```

---

## Authentication Methods

### Method 1: kubeconfig File (Local Development)

The easiest method — reads the same config `kubectl` uses.

```php
use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig(
    kubeConfigPath: '/home/user/.kube/config',
    context: 'my-cluster-context'   // optional, uses current-context if omitted
);
```

**Get your kubeconfig context names:**
```bash
kubectl config get-contexts
kubectl config current-context
```

---

### Method 2: Service Account Token (In-Cluster / Production)

When your PHP app runs **inside** the cluster (as a Pod), Kubernetes mounts credentials automatically:

```
/var/run/secrets/kubernetes.io/serviceaccount/token    ← Bearer token
/var/run/secrets/kubernetes.io/serviceaccount/ca.crt   ← CA certificate
```

```php
use Maclof\Kubernetes\Client;

$client = new Client([
    'master'          => 'https://kubernetes.default.svc',
    'ca_cert'         => '/var/run/secrets/kubernetes.io/serviceaccount/ca.crt',
    'token'           => file_get_contents('/var/run/secrets/kubernetes.io/serviceaccount/token'),
    'namespace'       => file_get_contents('/var/run/secrets/kubernetes.io/serviceaccount/namespace'),
]);
```

This is the **recommended approach for production** — no credentials are baked into your image.

---

### Method 3: Client Certificate + Key

```php
$client = new Client([
    'master'      => 'https://k8s.example.com:6443',
    'ca_cert'     => '/path/to/ca.crt',
    'client_cert' => '/path/to/client.crt',
    'client_key'  => '/path/to/client.key',
]);
```

**Extract certs from kubeconfig:**
```bash
# Get CA cert
kubectl config view --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}' | base64 -d > ca.crt

# Get client cert
kubectl config view --raw -o jsonpath='{.users[0].user.client-certificate-data}' | base64 -d > client.crt

# Get client key
kubectl config view --raw -o jsonpath='{.users[0].user.client-key-data}' | base64 -d > client.key
```

---

### Method 4: Bearer Token

```php
$client = new Client([
    'master' => 'https://k8s.example.com:6443',
    'ca_cert' => '/path/to/ca.crt',
    'token'  => 'eyJhbGci...your-token...',
]);
```

**Get a ServiceAccount token:**
```bash
# Kubernetes 1.24+
kubectl create token my-service-account -n default

# Or create a long-lived secret (older approach)
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: my-sa-token
  annotations:
    kubernetes.io/service-account.name: my-service-account
type: kubernetes.io/service-account-token
EOF

kubectl get secret my-sa-token -o jsonpath='{.data.token}' | base64 -d
```

---

## Environment-Based Configuration

Keep credentials out of code — load from environment:

```php
$client = new Client([
    'master'  => getenv('K8S_MASTER'),
    'ca_cert' => getenv('K8S_CA_CERT_PATH'),
    'token'   => getenv('K8S_TOKEN'),
]);
```

**.env file (development only — never commit):**
```
K8S_MASTER=https://127.0.0.1:6443
K8S_CA_CERT_PATH=/home/user/.minikube/ca.crt
K8S_TOKEN=eyJhbGci...
```

---

## Namespace Configuration

By default, the client targets the `default` namespace. You can set it globally or per-request:

```php
// Global namespace
$client = new Client([
    'master'    => '...',
    'namespace' => 'production',
    // ...
]);

// Switch namespace at runtime
$client->setNamespace('staging');
$pods = $client->pods()->find();

$client->setNamespace('production');
$deployments = $client->deployments()->find();
```

---

## Verifying the Connection

```php
<?php
require 'vendor/autoload.php';

use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig();

try {
    $namespaces = $client->namespaces()->find();
    echo "Connected! Namespaces found: " . count($namespaces) . "\n";
    foreach ($namespaces as $ns) {
        echo " - " . $ns->getMetadata('name') . "\n";
    }
} catch (\Exception $e) {
    echo "Connection failed: " . $e->getMessage() . "\n";
}
```

---

## Minikube Quick Start

```bash
# Start minikube
minikube start

# Get API endpoint and credentials
minikube ip                      # e.g. 192.168.49.2
kubectl cluster-info             # Shows API server URL

# Copy minikube certs
cp ~/.minikube/ca.crt /tmp/k8s-ca.crt
cp ~/.minikube/profiles/minikube/client.crt /tmp/k8s-client.crt
cp ~/.minikube/profiles/minikube/client.key /tmp/k8s-client.key
```

```php
$client = new Client([
    'master'      => 'https://192.168.49.2:8443',
    'ca_cert'     => '/tmp/k8s-ca.crt',
    'client_cert' => '/tmp/k8s-client.crt',
    'client_key'  => '/tmp/k8s-client.key',
]);
```

---

## Client Configuration Options Reference

| Option | Type | Description |
|---|---|---|
| `master` | string | API server URL |
| `ca_cert` | string | Path to CA certificate file |
| `client_cert` | string | Path to client certificate |
| `client_key` | string | Path to client private key |
| `token` | string | Bearer token string |
| `namespace` | string | Default namespace (default: `default`) |
| `timeout` | int | HTTP request timeout in seconds |
| `verify` | bool | Verify TLS (disable only for dev) |
