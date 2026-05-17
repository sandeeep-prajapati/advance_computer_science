# k8s-client-php — Master Cheatsheet

## Install

```bash
composer require maclof/kubernetes-client
# Laravel alternative:
composer require renoki-co/php-k8s
```

---

## Connect

```php
use Maclof\Kubernetes\Client;

// From kubeconfig (dev)
$client = Client::fromKubeConfig();
$client = Client::fromKubeConfig(context: 'my-cluster');

// In-cluster (Pod inside K8s)
$client = new Client([
    'master'    => 'https://kubernetes.default.svc',
    'ca_cert'   => '/var/run/secrets/kubernetes.io/serviceaccount/ca.crt',
    'token'     => file_get_contents('/var/run/secrets/kubernetes.io/serviceaccount/token'),
    'namespace' => file_get_contents('/var/run/secrets/kubernetes.io/serviceaccount/namespace'),
]);

// Client cert
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

// Switch namespace
$client->setNamespace('production');
$client->setNamespace('');  // all namespaces
```

---

## Pods

```php
// List
$client->pods()->find()
$client->pods()->setLabelSelector(['app' => 'nginx'])->find()
$client->pods()->setFieldSelector(['status.phase' => 'Running'])->find()

// Get
$client->pods()->findByName('my-pod')

// Logs
$client->pods()->logs('my-pod')
$client->pods()->logs('my-pod', ['tailLines' => 100])
$client->pods()->logs('my-pod', ['sinceSeconds' => 3600])
$client->pods()->logs('my-pod', ['previous' => true])          // crashed container
$client->pods()->logs('my-pod', ['container' => 'sidecar'])    // multi-container

// Create
$client->pods()->create(new Pod([
    'metadata' => ['name' => 'my-pod', 'labels' => ['app' => 'nginx']],
    'spec'     => ['containers' => [['name' => 'nginx', 'image' => 'nginx:alpine']]],
]))

// Delete
$client->pods()->deleteByName('my-pod')
$client->pods()->delete($pod, ['gracePeriodSeconds' => 0])     // force

// Read pod data
$pod->getMetadata('name')
$pod->getMetadata('namespace')
$pod->getSpec()['nodeName']
$pod->getStatus()['phase']              // Pending | Running | Succeeded | Failed
$pod->getStatus()['podIP']
$pod->getStatus()['containerStatuses']  // array of container states
```

---

## Deployments

```php
// List / Get
$client->deployments()->find()
$client->deployments()->findByName('my-app')

// Create
$client->deployments()->create(new Deployment([...]))

// Scale
$client->deployments()->patch($dep, ['spec' => ['replicas' => 5]])

// Update image (rolling update)
$client->deployments()->patch($dep, [
    'spec' => ['template' => ['spec' => ['containers' => [
        ['name' => 'app', 'image' => 'myapp:2.0'],
    ]]]],
])

// Delete
$client->deployments()->deleteByName('my-app')

// Read status
$dep->getSpec()['replicas']
$dep->getStatus()['readyReplicas']
$dep->getStatus()['updatedReplicas']
$dep->getStatus()['availableReplicas']
```

---

## Services

```php
// List / Get
$client->services()->find()
$client->services()->findByName('my-svc')

// Create ClusterIP
$client->services()->create(new Service([
    'metadata' => ['name' => 'my-svc'],
    'spec'     => [
        'type'     => 'ClusterIP',
        'selector' => ['app' => 'my-app'],
        'ports'    => [['port' => 80, 'targetPort' => 8080]],
    ],
]))

// Create LoadBalancer
// (same but 'type' => 'LoadBalancer')

// Read
$svc->getSpec()['clusterIP']
$svc->getSpec()['type']
$svc->getStatus()['loadBalancer']['ingress']   // external IPs

// Delete
$client->services()->deleteByName('my-svc')
```

---

## ConfigMaps

```php
// List / Get
$client->configMaps()->find()
$client->configMaps()->findByName('app-config')

// Read values
$cm->getAttribute('data')['MY_KEY']

// Create
$client->configMaps()->create(new ConfigMap([
    'metadata' => ['name' => 'app-config'],
    'data'     => ['KEY' => 'value', 'OTHER' => 'val'],
]))

// Update
$cm = $client->configMaps()->findByName('app-config');
$data = $cm->getAttribute('data');
$data['KEY'] = 'new-value';
$cm->setAttribute('data', $data);
$client->configMaps()->update($cm);

// Delete
$client->configMaps()->deleteByName('app-config')
```

---

## Secrets

```php
// List / Get  (never print values!)
$client->secrets()->find()
$client->secrets()->findByName('db-creds')

// Read (always base64_decode)
$val = base64_decode($secret->getAttribute('data')['password'])

// Create (always base64_encode values)
$client->secrets()->create(new Secret([
    'metadata' => ['name' => 'db-creds'],
    'type'     => 'Opaque',
    'data'     => [
        'username' => base64_encode('admin'),
        'password' => base64_encode('secret'),
    ],
]))

// Docker registry secret
'type' => 'kubernetes.io/dockerconfigjson',
'data' => ['.dockerconfigjson' => base64_encode(json_encode(['auths' => [...]]))],

// TLS secret
'type' => 'kubernetes.io/tls',
'data' => ['tls.crt' => base64_encode($cert), 'tls.key' => base64_encode($key)],

// Delete
$client->secrets()->deleteByName('db-creds')
```

---

## Namespaces

```php
$client->namespaces()->find()
$client->namespaces()->findByName('production')

$client->namespaces()->create(new NamespaceModel([
    'metadata' => ['name' => 'team-alpha', 'labels' => ['team' => 'alpha']],
]))

$client->namespaces()->deleteByName('team-alpha')  // deletes ALL resources inside!
```

---

## PVC / Storage

```php
$client->persistentVolumeClaims()->find()
$client->persistentVolumeClaims()->findByName('my-pvc')

$client->persistentVolumeClaims()->create(new PersistentVolumeClaim([
    'metadata' => ['name' => 'my-pvc'],
    'spec'     => [
        'accessModes'      => ['ReadWriteOnce'],
        'storageClassName' => 'standard',
        'resources'        => ['requests' => ['storage' => '10Gi']],
    ],
]))

// Status
$pvc->getStatus()['phase']          // Pending | Bound | Lost
$pvc->getSpec()['volumeName']       // which PV it's bound to

$client->persistentVolumeClaims()->deleteByName('my-pvc')
$client->persistentVolumes()->find()
$client->storageClasses()->find()
```

---

## Jobs & CronJobs

```php
// Jobs
$client->jobs()->find()
$client->jobs()->findByName('my-job')
$client->jobs()->create(new Job([
    'metadata' => ['name' => 'my-job'],
    'spec' => [
        'backoffLimit' => 3,
        'ttlSecondsAfterFinished' => 3600,
        'template' => ['spec' => [
            'restartPolicy' => 'OnFailure',
            'containers'    => [['name' => 'worker', 'image' => 'myapp:latest', 'command' => ['php', 'task.php']]],
        ]],
    ],
]))
$job->getStatus()['succeeded']
$job->getStatus()['failed']
$job->getStatus()['active']
$client->jobs()->deleteByName('my-job', ['propagationPolicy' => 'Foreground'])

// CronJobs
$client->cronJobs()->find()
$client->cronJobs()->patch($cj, ['spec' => ['suspend' => true]])   // pause
$client->cronJobs()->patch($cj, ['spec' => ['suspend' => false]])  // resume
$cj->getSpec()['schedule']
$cj->getStatus()['lastScheduleTime']
```

---

## StatefulSets & DaemonSets

```php
$client->statefulSets()->find()
$client->statefulSets()->findByName('mysql')
$client->statefulSets()->patch($sts, ['spec' => ['replicas' => 5]])
$sts->getStatus()['readyReplicas']

$client->daemonSets()->find()
$ds->getStatus()['desiredNumberScheduled']
$ds->getStatus()['numberReady']
```

---

## Ingress

```php
$client->ingresses()->find()
$client->ingresses()->findByName('my-ingress')
$client->ingresses()->create(new Ingress([
    'metadata' => ['name' => 'my-ingress', 'annotations' => [
        'nginx.ingress.kubernetes.io/rewrite-target' => '/',
    ]],
    'spec' => [
        'ingressClassName' => 'nginx',
        'rules' => [[
            'host' => 'myapp.example.com',
            'http' => ['paths' => [[
                'path'     => '/',
                'pathType' => 'Prefix',
                'backend'  => ['service' => ['name' => 'my-svc', 'port' => ['number' => 80]]],
            ]]],
        ]],
        'tls' => [['hosts' => ['myapp.example.com'], 'secretName' => 'my-tls']],
    ],
]))
$client->ingresses()->deleteByName('my-ingress')
```

---

## RBAC

```php
// Roles
$client->roles()->create(new Role([
    'metadata' => ['name' => 'pod-reader', 'namespace' => 'default'],
    'rules'    => [['apiGroups' => [''], 'resources' => ['pods'], 'verbs' => ['get', 'list', 'watch']]],
]))

// ClusterRoles
$client->clusterRoles()->create(new ClusterRole([
    'metadata' => ['name' => 'cluster-reader'],
    'rules'    => [['apiGroups' => ['*'], 'resources' => ['*'], 'verbs' => ['get', 'list', 'watch']]],
]))

// RoleBindings
$client->roleBindings()->create(new RoleBinding([
    'metadata' => ['name' => 'ci-binding', 'namespace' => 'default'],
    'subjects' => [['kind' => 'ServiceAccount', 'name' => 'ci-sa', 'namespace' => 'default']],
    'roleRef'  => ['apiGroup' => 'rbac.authorization.k8s.io', 'kind' => 'Role', 'name' => 'pod-reader'],
]))

// ServiceAccounts
$client->serviceAccounts()->create(new ServiceAccount(['metadata' => ['name' => 'ci-sa']]))
$client->serviceAccounts()->find()
```

---

## Watch (Real-time Events)

```php
$client->pods()->watch(function ($type, $pod) {
    // $type = ADDED | MODIFIED | DELETED | ERROR
    echo "[$type] " . $pod->getMetadata('name') . "\n";
    return true;   // return false to stop watching
});

$client->deployments()->watch(function ($type, $dep) {
    return true;
});

$client->events()->watch(function ($type, $event) {
    if ($event->getAttribute('type') === 'Warning') {
        echo $event->getAttribute('message') . "\n";
    }
    return true;
});
```

---

## Nodes

```php
$client->nodes()->find()
$client->nodes()->findByName('node-1')

$node->getStatus()['conditions']         // Ready, MemoryPressure, DiskPressure
$node->getStatus()['capacity']           // cpu, memory, pods
$node->getStatus()['allocatable']        // available capacity
$node->getMetadata('labels')             // node labels
```

---

## Error Handling

```php
use GuzzleHttp\Exception\ClientException;
use GuzzleHttp\Exception\ConnectException;

try {
    $client->pods()->find();
} catch (ClientException $e) {
    $code = $e->getResponse()->getStatusCode();
    // 401 = bad token | 403 = no RBAC | 404 = not found | 409 = conflict | 422 = invalid spec
} catch (ConnectException $e) {
    // API server unreachable
}
```

---

## Model Imports

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
