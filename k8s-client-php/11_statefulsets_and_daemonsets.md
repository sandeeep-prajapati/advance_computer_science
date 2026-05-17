# Chapter 11: StatefulSets & DaemonSets

## StatefulSets vs Deployments

| Feature | Deployment | StatefulSet |
|---|---|---|
| Pod names | `app-7d6f8-xyz` (random) | `app-0`, `app-1`, `app-2` (stable) |
| Pod DNS | Random | `pod-0.headless-svc.ns.svc.cluster.local` |
| Scaling order | Parallel | Ordered (0 → 1 → 2) |
| Storage | Shared or ephemeral | Each pod gets its own PVC (volumeClaimTemplates) |
| Use case | Stateless apps | Databases, queues, distributed systems |

---

## Working with StatefulSets

### List StatefulSets

```php
use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig();

$statefulSets = $client->statefulSets()->find();

foreach ($statefulSets as $sts) {
    $name     = $sts->getMetadata('name');
    $desired  = $sts->getSpec()['replicas'] ?? 0;
    $ready    = $sts->getStatus()['readyReplicas'] ?? 0;
    $current  = $sts->getStatus()['currentReplicas'] ?? 0;

    printf("%-30s %d/%d  current=%d\n", $name, $ready, $desired, $current);
}
```

### Create a StatefulSet (MySQL example)

```php
use Maclof\Kubernetes\Models\StatefulSet;
use Maclof\Kubernetes\Models\Service;

// First, create the headless service for stable DNS
$headlessService = new Service([
    'metadata' => ['name' => 'mysql-headless', 'namespace' => 'default'],
    'spec' => [
        'clusterIP' => 'None',
        'selector'  => ['app' => 'mysql'],
        'ports'     => [['port' => 3306, 'name' => 'mysql']],
    ],
]);
$client->services()->create($headlessService);

// Then create the StatefulSet
$statefulSet = new StatefulSet([
    'metadata' => [
        'name'      => 'mysql',
        'namespace' => 'default',
    ],
    'spec' => [
        'serviceName' => 'mysql-headless',
        'replicas'    => 3,
        'selector'    => ['matchLabels' => ['app' => 'mysql']],
        'template'    => [
            'metadata' => ['labels' => ['app' => 'mysql']],
            'spec'     => [
                'containers' => [[
                    'name'  => 'mysql',
                    'image' => 'mysql:8.0',
                    'env'   => [
                        ['name' => 'MYSQL_ROOT_PASSWORD', 'valueFrom' => [
                            'secretKeyRef' => ['name' => 'mysql-secret', 'key' => 'root-password'],
                        ]],
                    ],
                    'ports' => [['containerPort' => 3306, 'name' => 'mysql']],
                    'volumeMounts' => [
                        ['name' => 'data', 'mountPath' => '/var/lib/mysql'],
                    ],
                    'resources' => [
                        'requests' => ['memory' => '512Mi', 'cpu' => '250m'],
                        'limits'   => ['memory' => '1Gi', 'cpu' => '500m'],
                    ],
                    'readinessProbe' => [
                        'exec' => ['command' => ['bash', '-c', 'mysqladmin ping -h 127.0.0.1']],
                        'initialDelaySeconds' => 20,
                        'periodSeconds'       => 10,
                    ],
                ]],
            ],
        ],
        'volumeClaimTemplates' => [
            [
                'metadata' => ['name' => 'data'],
                'spec'     => [
                    'accessModes'      => ['ReadWriteOnce'],
                    'storageClassName' => 'standard',
                    'resources'        => ['requests' => ['storage' => '20Gi']],
                ],
            ],
        ],
    ],
]);

$client->statefulSets()->create($statefulSet);
echo "MySQL StatefulSet created with 3 replicas\n";
// Pods: mysql-0, mysql-1, mysql-2
// DNS:  mysql-0.mysql-headless.default.svc.cluster.local
```

### Scaling a StatefulSet

```php
$client->statefulSets()->patch(
    $client->statefulSets()->findByName('mysql'),
    ['spec' => ['replicas' => 5]]
);
// Scales up ordered: mysql-3, then mysql-4
```

### StatefulSet Pod DNS Names

```php
function getStatefulSetPodDNS(string $podName, string $serviceName, string $namespace = 'default'): string
{
    return "$podName.$serviceName.$namespace.svc.cluster.local";
}

// MySQL cluster primary/replica discovery
$sts    = $client->statefulSets()->findByName('mysql');
$count  = $sts->getSpec()['replicas'] ?? 0;
$svcName = $sts->getSpec()['serviceName'];

for ($i = 0; $i < $count; $i++) {
    $dns = getStatefulSetPodDNS("mysql-$i", $svcName);
    echo "Replica $i: $dns\n";
}
```

---

## Update Strategy for StatefulSets

```php
// RollingUpdate (default): pods updated one at a time, from highest to lowest index
// OnDelete: pods only updated when manually deleted

$sts = $client->statefulSets()->findByName('mysql');
$client->statefulSets()->patch($sts, [
    'spec' => [
        'updateStrategy' => [
            'type' => 'RollingUpdate',
            'rollingUpdate' => [
                'partition' => 2,  // only update pods with index >= 2 (canary)
            ],
        ],
        'template' => [
            'spec' => [
                'containers' => [['name' => 'mysql', 'image' => 'mysql:8.1']],
            ],
        ],
    ],
]);
```

---

## Working with DaemonSets

A **DaemonSet** runs exactly one pod per node. Use for:
- Log collectors (Fluentd, Filebeat)
- Monitoring agents (Prometheus node-exporter, Datadog)
- Network plugins (CNI agents)
- Security scanners

### List DaemonSets

```php
$daemonSets = $client->daemonSets()->find();

foreach ($daemonSets as $ds) {
    $name      = $ds->getMetadata('name');
    $desired   = $ds->getStatus()['desiredNumberScheduled'] ?? 0;
    $ready     = $ds->getStatus()['numberReady'] ?? 0;
    $available = $ds->getStatus()['numberAvailable'] ?? 0;

    printf("%-35s desired=%d  ready=%d  available=%d\n",
        $name, $desired, $ready, $available);
}
```

### Create a DaemonSet (log collector)

```php
use Maclof\Kubernetes\Models\DaemonSet;

$daemonSet = new DaemonSet([
    'metadata' => [
        'name'      => 'fluentd',
        'namespace' => 'kube-system',
        'labels'    => ['app' => 'fluentd'],
    ],
    'spec' => [
        'selector' => ['matchLabels' => ['app' => 'fluentd']],
        'updateStrategy' => [
            'type' => 'RollingUpdate',
            'rollingUpdate' => ['maxUnavailable' => 1],
        ],
        'template' => [
            'metadata' => ['labels' => ['app' => 'fluentd']],
            'spec' => [
                'tolerations' => [
                    ['key' => 'node-role.kubernetes.io/control-plane', 'effect' => 'NoSchedule'],
                ],
                'containers' => [[
                    'name'  => 'fluentd',
                    'image' => 'fluent/fluentd-kubernetes-daemonset:v1.16',
                    'resources' => [
                        'limits'   => ['memory' => '200Mi'],
                        'requests' => ['cpu' => '100m', 'memory' => '200Mi'],
                    ],
                    'volumeMounts' => [
                        ['name' => 'varlog',        'mountPath' => '/var/log'],
                        ['name' => 'varlibdockercontainers', 'mountPath' => '/var/lib/docker/containers', 'readOnly' => true],
                    ],
                ]],
                'volumes' => [
                    ['name' => 'varlog', 'hostPath' => ['path' => '/var/log']],
                    ['name' => 'varlibdockercontainers', 'hostPath' => ['path' => '/var/lib/docker/containers']],
                ],
            ],
        ],
    ],
]);

$client->daemonSets()->create($daemonSet);
echo "Fluentd DaemonSet deployed to all nodes\n";
```

### DaemonSet with Node Selector (Specific Nodes Only)

```php
$daemonSet = new DaemonSet([
    'metadata' => ['name' => 'gpu-monitor'],
    'spec' => [
        'selector'  => ['matchLabels' => ['app' => 'gpu-monitor']],
        'template'  => [
            'metadata' => ['labels' => ['app' => 'gpu-monitor']],
            'spec'     => [
                'nodeSelector' => ['gpu' => 'true'],   // Only nodes with this label
                'containers'   => [[
                    'name'  => 'gpu-monitor',
                    'image' => 'gpu-metrics:latest',
                ]],
            ],
        ],
    ],
]);

$client->daemonSets()->create($daemonSet);
```

### Check DaemonSet Health

```php
function isDaemonSetHealthy(Client $client, string $name): bool
{
    $ds = $client->daemonSets()->findByName($name);
    $status = $ds->getStatus();

    $desired   = $status['desiredNumberScheduled'] ?? 0;
    $ready     = $status['numberReady'] ?? 0;

    return $desired > 0 && $ready === $desired;
}
```

---

## Delete StatefulSet / DaemonSet

```php
// Delete StatefulSet (does NOT delete associated PVCs — they persist)
$client->statefulSets()->deleteByName('mysql');

// Delete PVCs manually after
$pvcs = $client->persistentVolumeClaims()->setLabelSelector(['app' => 'mysql'])->find();
foreach ($pvcs as $pvc) {
    $client->persistentVolumeClaims()->deleteByName($pvc->getMetadata('name'));
}

// Delete DaemonSet
$client->daemonSets()->deleteByName('fluentd');
```
