# Kubernetes PHP Client — Complete Developer Book

A comprehensive guide to managing Kubernetes resources programmatically using PHP. Covers everything from first connection to production-grade operators, canary deployments, and multi-cluster management.

---

## Quick Start

```bash
composer require maclof/kubernetes-client
```

```php
use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig();  // reads ~/.kube/config

// List pods
foreach ($client->pods()->find() as $pod) {
    echo $pod->getMetadata('name') . "\n";
}
```

---

## Full Table of Contents

### Part 1: Foundations
- [Chapter 1: Introduction](01_introduction.md) — What is k8s-client-php, when to use PHP vs kubectl, API architecture
- [Chapter 2: Installation & Setup](02_installation_and_setup.md) — Composer, authentication methods, kubeconfig, env config
- [Chapter 3: Connecting to a Cluster](03_connecting_to_cluster.md) — Local (minikube/kind/k3s), remote (GKE/EKS/AKS), in-cluster

### Part 2: Core Resources
- [Chapter 4: Pods](04_pods.md) — List, get, create, delete, logs, exec, label selectors
- [Chapter 5: Deployments & ReplicaSets](05_deployments_and_replicasets.md) — Create, scale, rolling update, rollback, wait for rollout
- [Chapter 6: Services & Networking](06_services_and_networking.md) — ClusterIP, NodePort, LoadBalancer, ExternalName, endpoints
- [Chapter 7: ConfigMaps & Secrets](07_configmaps_and_secrets.md) — Create, read, update config and secrets, rotation patterns
- [Chapter 8: Persistent Volumes & Storage](08_persistent_volumes_and_storage.md) — PV, PVC, StorageClass, dynamic provisioning
- [Chapter 9: Namespaces & RBAC](09_namespaces_and_rbac.md) — Namespaces, Roles, ClusterRoles, Bindings, ServiceAccounts

### Part 3: Advanced Resources
- [Chapter 10: Ingress](10_ingress.md) — Path routing, TLS, nginx annotations, multi-host
- [Chapter 11: StatefulSets & DaemonSets](11_statefulsets_and_daemonsets.md) — Databases, log agents, stable DNS
- [Chapter 12: Jobs & CronJobs](12_jobs_and_cronjobs.md) — One-off tasks, scheduled jobs, parallel batch processing
- [Chapter 13: Watching Events](13_watching_events.md) — Real-time Watch API, pod restarts, deployment rollouts
- [Chapter 14: Custom Resources (CRDs)](14_custom_resources.md) — cert-manager, ArgoCD, defining your own CRDs

### Part 4: Production
- [Chapter 15: Error Handling & Retry Logic](15_error_handling.md) — HTTP error codes, retries, backoff, idempotent patterns
- [Chapter 16: Production Patterns](16_production_patterns.md) — Operator pattern, multi-cluster, canary deployments, health dashboards

---

## API Method Reference

### Pods
```php
$client->pods()->find()                          // list
$client->pods()->setLabelSelector(['k'=>'v'])->find()
$client->pods()->findByName('name')
$client->pods()->create($model)
$client->pods()->delete($model)
$client->pods()->deleteByName('name')
$client->pods()->logs('name', ['tailLines'=>100])
```

### Deployments
```php
$client->deployments()->find()
$client->deployments()->findByName('name')
$client->deployments()->create($model)
$client->deployments()->update($model)
$client->deployments()->patch($model, $partialSpec)
$client->deployments()->delete($model)
$client->deployments()->deleteByName('name')
```

### Services
```php
$client->services()->find()
$client->services()->findByName('name')
$client->services()->create($model)
$client->services()->update($model)
$client->services()->delete($model)
```

### ConfigMaps & Secrets
```php
$client->configMaps()->find()
$client->configMaps()->findByName('name')
$client->configMaps()->create($model)
$client->configMaps()->update($model)
$client->secrets()->find()
$client->secrets()->create($model)
```

### Other Resources
```php
$client->namespaces()->find()
$client->namespaces()->create($model)
$client->persistentVolumeClaims()->find()
$client->persistentVolumes()->find()
$client->storageClasses()->find()
$client->statefulSets()->find()
$client->daemonSets()->find()
$client->jobs()->find()
$client->cronJobs()->find()
$client->ingresses()->find()
$client->roles()->find()
$client->roleBindings()->find()
$client->clusterRoles()->find()
$client->clusterRoleBindings()->find()
$client->serviceAccounts()->find()
$client->nodes()->find()
$client->events()->find()
```

---

## Authentication Quick Reference

```php
// From kubeconfig (local dev)
$client = Client::fromKubeConfig();
$client = Client::fromKubeConfig(context: 'my-cluster');

// In-cluster (running as a Pod)
$client = new Client([
    'master'    => 'https://kubernetes.default.svc',
    'ca_cert'   => '/var/run/secrets/kubernetes.io/serviceaccount/ca.crt',
    'token'     => file_get_contents('/var/run/secrets/kubernetes.io/serviceaccount/token'),
    'namespace' => file_get_contents('/var/run/secrets/kubernetes.io/serviceaccount/namespace'),
]);

// Client certificates
$client = new Client([
    'master'      => 'https://k8s.example.com:6443',
    'ca_cert'     => '/path/to/ca.crt',
    'client_cert' => '/path/to/client.crt',
    'client_key'  => '/path/to/client.key',
]);

// Bearer token
$client = new Client([
    'master'  => 'https://k8s.example.com:6443',
    'ca_cert' => '/path/to/ca.crt',
    'token'   => 'eyJhbGci...',
]);
```

---

## Model Import Reference

```php
use Maclof\Kubernetes\Models\Pod;
use Maclof\Kubernetes\Models\Deployment;
use Maclof\Kubernetes\Models\Service;
use Maclof\Kubernetes\Models\ConfigMap;
use Maclof\Kubernetes\Models\Secret;
use Maclof\Kubernetes\Models\PersistentVolume;
use Maclof\Kubernetes\Models\PersistentVolumeClaim;
use Maclof\Kubernetes\Models\StorageClass;
use Maclof\Kubernetes\Models\StatefulSet;
use Maclof\Kubernetes\Models\DaemonSet;
use Maclof\Kubernetes\Models\Job;
use Maclof\Kubernetes\Models\CronJob;
use Maclof\Kubernetes\Models\Ingress;
use Maclof\Kubernetes\Models\NamespaceModel;
use Maclof\Kubernetes\Models\ServiceAccount;
use Maclof\Kubernetes\Models\Role;
use Maclof\Kubernetes\Models\ClusterRole;
use Maclof\Kubernetes\Models\RoleBinding;
use Maclof\Kubernetes\Models\ClusterRoleBinding;
```
