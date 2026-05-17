# Chapter 13: Watching Events & Real-Time Updates

## The Watch API

The Kubernetes Watch API streams resource change events in real time over a long-lived HTTP connection. Each event has a type:

| Event Type | Meaning |
|---|---|
| `ADDED` | Resource was created |
| `MODIFIED` | Resource was updated |
| `DELETED` | Resource was deleted |
| `ERROR` | Watch error (restart the watch) |

---

## Watching Pods

```php
use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig();

// Watch all pods in the current namespace
$client->pods()->watch(function (string $type, $pod) {
    $name  = $pod->getMetadata('name');
    $phase = $pod->getStatus()['phase'] ?? 'Unknown';

    switch ($type) {
        case 'ADDED':
            echo "[+] Pod added:    $name ($phase)\n";
            break;
        case 'MODIFIED':
            echo "[~] Pod changed:  $name ($phase)\n";
            break;
        case 'DELETED':
            echo "[-] Pod deleted:  $name\n";
            break;
        case 'ERROR':
            echo "[!] Watch error for $name\n";
            break;
    }

    return true;  // return false to stop the watch
});
```

---

## Watching Pods with Label Filter

```php
$client->pods()->setLabelSelector(['app' => 'my-app'])->watch(function ($type, $pod) {
    $name   = $pod->getMetadata('name');
    $reason = $pod->getStatus()['conditions'][0]['reason'] ?? '-';
    echo "[$type] $name — $reason\n";
    return true;
});
```

---

## Watching Deployments (Detect Rollouts)

```php
$client->deployments()->watch(function ($type, $deployment) {
    if ($type !== 'MODIFIED') return true;

    $name      = $deployment->getMetadata('name');
    $desired   = $deployment->getSpec()['replicas'] ?? 0;
    $ready     = $deployment->getStatus()['readyReplicas'] ?? 0;
    $updated   = $deployment->getStatus()['updatedReplicas'] ?? 0;

    echo "Deployment '$name': updated=$updated ready=$ready desired=$desired\n";

    if ($updated === $desired && $ready === $desired) {
        echo "✓ Deployment '$name' rollout complete!\n";
    }

    return true;
});
```

---

## Watching Events (kubectl get events equivalent)

```php
use Maclof\Kubernetes\Models\Event;

$client->events()->watch(function ($type, $event) {
    if ($type !== 'ADDED') return true;

    $name      = $event->getMetadata('name');
    $reason    = $event->getAttribute('reason') ?? '-';
    $message   = $event->getAttribute('message') ?? '-';
    $objName   = $event->getAttribute('involvedObject')['name'] ?? '-';
    $objKind   = $event->getAttribute('involvedObject')['kind'] ?? '-';
    $eventType = $event->getAttribute('type') ?? 'Normal';  // Normal or Warning

    $prefix = $eventType === 'Warning' ? '[WARN]' : '[INFO]';
    echo "$prefix $objKind/$objName — $reason: $message\n";

    return true;
});
```

### Filter Warning Events Only

```php
$client->events()->watch(function ($type, $event) {
    $evtType = $event->getAttribute('type') ?? 'Normal';

    if ($evtType !== 'Warning') return true;  // skip non-warnings

    $reason  = $event->getAttribute('reason') ?? '-';
    $message = $event->getAttribute('message') ?? '-';
    $obj     = $event->getAttribute('involvedObject')['name'] ?? '-';
    $count   = $event->getAttribute('count') ?? 1;

    echo "[WARNING x$count] $obj — $reason: $message\n";

    return true;
});
```

---

## Detecting Pod Restarts

```php
function watchForPodRestarts(Client $client, string $namespace = 'default'): void
{
    $client->setNamespace($namespace);

    $knownRestarts = [];

    $client->pods()->watch(function ($type, $pod) use (&$knownRestarts) {
        if ($type === 'DELETED') {
            unset($knownRestarts[$pod->getMetadata('name')]);
            return true;
        }

        $name     = $pod->getMetadata('name');
        $statuses = $pod->getStatus()['containerStatuses'] ?? [];

        foreach ($statuses as $cs) {
            $container     = $cs['name'];
            $restartCount  = $cs['restartCount'] ?? 0;
            $key           = "$name/$container";

            $prev = $knownRestarts[$key] ?? 0;

            if ($restartCount > $prev) {
                $diff = $restartCount - $prev;
                echo "[RESTART] Pod=$name Container=$container restarts=$restartCount (+$diff)\n";

                // Optionally fetch logs
                if ($restartCount > 5) {
                    echo "[CRITICAL] High restart count! Fetching previous logs...\n";
                    // $logs = $client->pods()->logs($name, ['previous' => true, 'container' => $container]);
                }
            }

            $knownRestarts[$key] = $restartCount;
        }

        return true;
    });
}

watchForPodRestarts($client, 'production');
```

---

## Watch with Timeout (Non-blocking)

```php
function watchWithTimeout(Client $client, int $seconds): void
{
    $deadline = time() + $seconds;

    $client->pods()->watch(function ($type, $pod) use ($deadline) {
        $name = $pod->getMetadata('name');
        echo "[$type] $name\n";

        // Stop watching after timeout
        return time() < $deadline;
    });

    echo "Watch ended after $seconds seconds\n";
}

watchWithTimeout($client, 30);
```

---

## Polling (Alternative to Watching)

For PHP-FPM environments where long-lived connections aren't suitable, use polling:

```php
function pollDeploymentStatus(Client $client, string $name, int $pollEvery = 5): void
{
    echo "Polling deployment '$name' every {$pollEvery}s...\n";

    while (true) {
        $dep    = $client->deployments()->findByName($name);
        $status = $dep->getStatus();

        $desired   = $dep->getSpec()['replicas'] ?? 0;
        $ready     = $status['readyReplicas'] ?? 0;
        $updated   = $status['updatedReplicas'] ?? 0;

        printf("[%s] desired=%d updated=%d ready=%d\n",
            date('H:i:s'), $desired, $updated, $ready);

        if ($updated === $desired && $ready === $desired) {
            echo "Rollout complete!\n";
            break;
        }

        sleep($pollEvery);
    }
}

pollDeploymentStatus($client, 'my-app');
```

---

## Watching Multiple Resource Types

```php
function startClusterMonitor(Client $client): void
{
    // Note: PHP is single-threaded — use separate processes or queues for parallel watching
    // This example watches pods sequentially for demonstration

    $resources = ['pods', 'deployments', 'services'];

    foreach ($resources as $resource) {
        echo "Starting watch for: $resource\n";
        $client->$resource()->watch(function ($type, $obj) use ($resource) {
            echo "[$resource] $type: " . $obj->getMetadata('name') . "\n";
            return true;
        });
    }
}
```

---

## Resource Version and Watch Resumption

The Watch API uses `resourceVersion` to resume from where it left off:

```php
// Get current resource version from a list
$pods = $client->pods()->find();
$resourceVersion = $pods->getResourceVersion();  // e.g. "12345"

// Watch only events that happened AFTER this version
$client->pods()->watch(function ($type, $pod) {
    echo "[$type] " . $pod->getMetadata('name') . "\n";
    return true;
}, ['resourceVersion' => $resourceVersion]);
```

---

## Kubernetes Events vs Watch Events

| | `$client->events()->watch()` | `$client->pods()->watch()` |
|---|---|---|
| What it shows | Human-readable cluster events (like `kubectl get events`) | Raw API resource change stream |
| CrashLoopBackOff | `BackOff` Warning event | MODIFIED pod with `CrashLoopBackOff` status |
| Scaling event | `ScalingReplicaSet` event | ADDED/DELETED pod events |
| Best for | Alerting, audit trail | Automation, operators |
