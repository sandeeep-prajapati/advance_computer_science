# Chapter 4: Working with Pods

## Pod Resource Methods

| Method | Description |
|---|---|
| `$client->pods()->find()` | List all pods in current namespace |
| `$client->pods()->findByName('name')` | Get a single pod |
| `$client->pods()->create($model)` | Create a new pod |
| `$client->pods()->delete($model)` | Delete a pod |
| `$client->pods()->deleteByName('name')` | Delete pod by name |
| `$client->pods()->find(['label' => 'value'])` | Filter pods by label |
| `$client->pods()->logs('name')` | Fetch pod logs |

---

## Listing Pods

```php
use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig();

// List all pods in current namespace
$pods = $client->pods()->find();

foreach ($pods as $pod) {
    $name      = $pod->getMetadata('name');
    $namespace = $pod->getMetadata('namespace');
    $phase     = $pod->getStatus()['phase'] ?? 'Unknown';
    $nodeName  = $pod->getSpec()['nodeName'] ?? 'unscheduled';

    printf("%-40s %-15s %-10s %s\n", $name, $namespace, $phase, $nodeName);
}
```

### List Pods Across All Namespaces

```php
$client->setNamespace('');  // empty string = all namespaces

$allPods = $client->pods()->find();
echo "Total pods cluster-wide: " . count($allPods) . "\n";
```

---

## Filtering Pods by Label

```php
// Pods with label app=nginx
$nginxPods = $client->pods()->setLabelSelector([
    'app' => 'nginx',
])->find();

// Multiple labels (AND condition)
$prodNginx = $client->pods()->setLabelSelector([
    'app' => 'nginx',
    'env' => 'production',
])->find();

// Field selector
$runningPods = $client->pods()->setFieldSelector([
    'status.phase' => 'Running',
])->find();
```

---

## Getting a Single Pod

```php
$pod = $client->pods()->findByName('my-nginx-pod');

if ($pod === null) {
    echo "Pod not found\n";
    exit(1);
}

// Access metadata
echo "Name:      " . $pod->getMetadata('name') . "\n";
echo "Namespace: " . $pod->getMetadata('namespace') . "\n";
echo "Node:      " . ($pod->getSpec()['nodeName'] ?? 'N/A') . "\n";
echo "Phase:     " . ($pod->getStatus()['phase'] ?? 'N/A') . "\n";

// Container statuses
$containerStatuses = $pod->getStatus()['containerStatuses'] ?? [];
foreach ($containerStatuses as $cs) {
    echo "Container: {$cs['name']}  Ready: " . ($cs['ready'] ? 'yes' : 'no') . "\n";
}
```

---

## Creating a Pod

```php
use Maclof\Kubernetes\Models\Pod;

$pod = new Pod([
    'metadata' => [
        'name'      => 'my-nginx',
        'namespace' => 'default',
        'labels'    => [
            'app' => 'nginx',
            'env' => 'dev',
        ],
    ],
    'spec' => [
        'containers' => [
            [
                'name'  => 'nginx',
                'image' => 'nginx:1.25-alpine',
                'ports' => [
                    ['containerPort' => 80],
                ],
                'resources' => [
                    'requests' => ['memory' => '64Mi', 'cpu' => '100m'],
                    'limits'   => ['memory' => '128Mi', 'cpu' => '200m'],
                ],
                'env' => [
                    ['name' => 'ENV',  'value' => 'production'],
                    ['name' => 'PORT', 'value' => '8080'],
                ],
            ],
        ],
        'restartPolicy' => 'Always',
    ],
]);

$created = $client->pods()->create($pod);
echo "Pod created: " . $created->getMetadata('name') . "\n";
```

---

## Creating a Pod with Multiple Containers

```php
$pod = new Pod([
    'metadata' => ['name' => 'sidecar-example'],
    'spec' => [
        'containers' => [
            [
                'name'  => 'app',
                'image' => 'myapp:1.0',
                'ports' => [['containerPort' => 8080]],
            ],
            [
                'name'  => 'log-shipper',
                'image' => 'fluent/fluent-bit:latest',
                'volumeMounts' => [
                    ['name' => 'app-logs', 'mountPath' => '/var/log/app'],
                ],
            ],
        ],
        'volumes' => [
            ['name' => 'app-logs', 'emptyDir' => []],
        ],
    ],
]);

$client->pods()->create($pod);
```

---

## Pod Logs

```php
// Get current logs
$logs = $client->pods()->logs('my-nginx');
echo $logs;

// Get logs for a specific container (multi-container pod)
$logs = $client->pods()->logs('sidecar-pod', ['container' => 'log-shipper']);

// Last N lines
$logs = $client->pods()->logs('my-nginx', ['tailLines' => 100]);

// Logs since a duration
$logs = $client->pods()->logs('my-nginx', ['sinceSeconds' => 3600]);

// Previous container instance (for CrashLoopBackOff debugging)
$logs = $client->pods()->logs('my-nginx', ['previous' => true]);
```

---

## Checking Pod Status

```php
function getPodStatus(Client $client, string $podName): string
{
    $pod = $client->pods()->findByName($podName);

    if ($pod === null) {
        return 'NOT_FOUND';
    }

    $phase = $pod->getStatus()['phase'] ?? 'Unknown';

    $containerStatuses = $pod->getStatus()['containerStatuses'] ?? [];
    foreach ($containerStatuses as $cs) {
        if (isset($cs['state']['waiting']['reason'])) {
            return $cs['state']['waiting']['reason'];  // CrashLoopBackOff, ImagePullBackOff, etc.
        }
    }

    return $phase;
}

// Wait for a pod to be Running
function waitForPodRunning(Client $client, string $podName, int $timeoutSec = 60): bool
{
    $start = time();
    while (time() - $start < $timeoutSec) {
        $status = getPodStatus($client, $podName);
        if ($status === 'Running') {
            return true;
        }
        echo "Pod status: $status — waiting...\n";
        sleep(2);
    }
    return false;
}
```

---

## Deleting a Pod

```php
// Delete by name
$client->pods()->deleteByName('my-nginx');
echo "Pod deleted\n";

// Delete with grace period (seconds)
$pod = $client->pods()->findByName('my-nginx');
$client->pods()->delete($pod, [
    'gracePeriodSeconds' => 0,  // force delete
]);

// Delete all pods with a label
$oldPods = $client->pods()->setLabelSelector(['app' => 'old-app'])->find();
foreach ($oldPods as $pod) {
    $client->pods()->delete($pod);
    echo "Deleted: " . $pod->getMetadata('name') . "\n";
}
```

---

## Pod IP and Networking Info

```php
$pod = $client->pods()->findByName('my-nginx');

$status = $pod->getStatus();
echo "Pod IP:   " . ($status['podIP'] ?? 'N/A') . "\n";
echo "Host IP:  " . ($status['hostIP'] ?? 'N/A') . "\n";

// All pod IPs (for multi-homed pods)
foreach ($status['podIPs'] ?? [] as $ipEntry) {
    echo "IP: " . $ipEntry['ip'] . "\n";
}
```

---

## Resource Summary

```php
// Print a kubectl get pods-like summary
function printPodTable(Client $client): void
{
    $pods = $client->pods()->find();

    printf("%-45s %-10s %-5s %-12s %-6s\n", 'NAME', 'STATUS', 'READY', 'NODE', 'AGE');
    printf("%s\n", str_repeat('-', 85));

    foreach ($pods as $pod) {
        $name     = $pod->getMetadata('name');
        $phase    = $pod->getStatus()['phase'] ?? 'Unknown';
        $node     = $pod->getSpec()['nodeName'] ?? '-';
        $statuses = $pod->getStatus()['containerStatuses'] ?? [];
        $ready    = count(array_filter($statuses, fn($s) => $s['ready']));
        $total    = count($statuses);

        $start = $pod->getStatus()['startTime'] ?? null;
        $age   = $start ? human_age(new DateTime($start)) : '-';

        printf("%-45s %-10s %-5s %-12s %-6s\n", $name, $phase, "$ready/$total", $node, $age);
    }
}

function human_age(DateTime $start): string
{
    $diff = (new DateTime())->diff($start);
    if ($diff->days > 0) return $diff->days . 'd';
    if ($diff->h > 0)    return $diff->h . 'h';
    return $diff->i . 'm';
}
```
