# Chapter 5: Deployments & ReplicaSets

## What Deployments Give You

| Feature | How |
|---|---|
| Desired pod count | `replicas: 3` — always maintained |
| Rolling updates | New pods before old ones are removed |
| Rollback | `kubectl rollout undo` — or via API |
| Pause / resume | Pause rollout mid-way for validation |
| History | Last N revisions stored in etcd |

---

## Listing Deployments

```php
use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig();

$deployments = $client->deployments()->find();

foreach ($deployments as $dep) {
    $name      = $dep->getMetadata('name');
    $desired   = $dep->getSpec()['replicas'] ?? 0;
    $ready     = $dep->getStatus()['readyReplicas'] ?? 0;
    $available = $dep->getStatus()['availableReplicas'] ?? 0;

    printf("%-40s %d/%d  available=%d\n", $name, $ready, $desired, $available);
}
```

---

## Getting a Single Deployment

```php
$dep = $client->deployments()->findByName('my-app');

echo "Image: " . $dep->getSpec()['template']['spec']['containers'][0]['image'] . "\n";
echo "Replicas desired:   " . $dep->getSpec()['replicas'] . "\n";
echo "Replicas ready:     " . ($dep->getStatus()['readyReplicas'] ?? 0) . "\n";
echo "Replicas available: " . ($dep->getStatus()['availableReplicas'] ?? 0) . "\n";
```

---

## Creating a Deployment

```php
use Maclof\Kubernetes\Models\Deployment;

$deployment = new Deployment([
    'metadata' => [
        'name'      => 'my-app',
        'namespace' => 'default',
        'labels'    => ['app' => 'my-app'],
    ],
    'spec' => [
        'replicas' => 3,
        'selector' => [
            'matchLabels' => ['app' => 'my-app'],
        ],
        'strategy' => [
            'type' => 'RollingUpdate',
            'rollingUpdate' => [
                'maxUnavailable' => 1,
                'maxSurge'       => 1,
            ],
        ],
        'template' => [
            'metadata' => [
                'labels' => ['app' => 'my-app'],
            ],
            'spec' => [
                'containers' => [
                    [
                        'name'  => 'my-app',
                        'image' => 'my-registry/my-app:1.0.0',
                        'ports' => [['containerPort' => 8080]],
                        'resources' => [
                            'requests' => ['memory' => '128Mi', 'cpu' => '100m'],
                            'limits'   => ['memory' => '256Mi', 'cpu' => '500m'],
                        ],
                        'readinessProbe' => [
                            'httpGet' => ['path' => '/health', 'port' => 8080],
                            'initialDelaySeconds' => 5,
                            'periodSeconds'       => 10,
                        ],
                        'livenessProbe' => [
                            'httpGet' => ['path' => '/health', 'port' => 8080],
                            'initialDelaySeconds' => 15,
                            'periodSeconds'       => 20,
                        ],
                    ],
                ],
            ],
        ],
    ],
]);

$client->deployments()->create($deployment);
echo "Deployment created\n";
```

---

## Scaling a Deployment

```php
// Read current deployment
$dep = $client->deployments()->findByName('my-app');

// Modify replica count
$spec = $dep->getSpec();
$spec['replicas'] = 5;

// Patch the deployment
$client->deployments()->patch($dep, [
    'spec' => ['replicas' => 5],
]);

echo "Scaled to 5 replicas\n";
```

### Scale via Dedicated Scale Subresource

```php
use Maclof\Kubernetes\Models\Scale;

$scale = new Scale([
    'metadata' => ['name' => 'my-app', 'namespace' => 'default'],
    'spec'     => ['replicas' => 10],
]);

$client->deployments()->scale('my-app', $scale);
```

---

## Updating a Deployment (Rolling Update)

```php
// Update image to trigger rolling update
$dep = $client->deployments()->findByName('my-app');

$dep->setAttribute('spec', array_merge_recursive($dep->getSpec(), [
    'template' => [
        'spec' => [
            'containers' => [
                [
                    'name'  => 'my-app',
                    'image' => 'my-registry/my-app:2.0.0',  // new version
                ],
            ],
        ],
    ],
]));

$client->deployments()->update($dep);
echo "Rolling update triggered to v2.0.0\n";
```

### Patch Method (Safer)

```php
$client->deployments()->patch($dep, [
    'spec' => [
        'template' => [
            'spec' => [
                'containers' => [[
                    'name'  => 'my-app',
                    'image' => 'my-registry/my-app:2.0.0',
                ]],
            ],
        ],
    ],
]);
```

---

## Waiting for Rollout Completion

```php
function waitForDeployment(Client $client, string $name, int $timeout = 120): bool
{
    $start = time();
    while (time() - $start < $timeout) {
        $dep = $client->deployments()->findByName($name);

        $desired   = $dep->getSpec()['replicas'] ?? 0;
        $ready     = $dep->getStatus()['readyReplicas'] ?? 0;
        $updated   = $dep->getStatus()['updatedReplicas'] ?? 0;
        $available = $dep->getStatus()['availableReplicas'] ?? 0;

        echo "Rollout: updated=$updated ready=$ready available=$available desired=$desired\n";

        if ($updated === $desired && $ready === $desired && $available === $desired) {
            return true;
        }

        sleep(3);
    }
    return false;
}

// Usage
$client->deployments()->patch($dep, ['spec' => ['template' => ['spec' => ['containers' => [['name' => 'app', 'image' => 'app:2.0']]]]]]);

if (waitForDeployment($client, 'my-app')) {
    echo "Deployment rolled out successfully\n";
} else {
    echo "Deployment timed out — check pod status\n";
}
```

---

## Deleting a Deployment

```php
// Delete by name (also deletes associated ReplicaSets and Pods)
$client->deployments()->deleteByName('my-app');

// With propagation policy
$dep = $client->deployments()->findByName('my-app');
$client->deployments()->delete($dep, [
    'propagationPolicy' => 'Foreground',  // wait for pods to terminate
]);
```

---

## Working with ReplicaSets

ReplicaSets are usually managed by Deployments, but you can interact with them directly:

```php
// List ReplicaSets owned by a deployment
$rsets = $client->replicaSets()->setLabelSelector(['app' => 'my-app'])->find();

foreach ($rsets as $rs) {
    $name    = $rs->getMetadata('name');
    $desired = $rs->getSpec()['replicas'] ?? 0;
    $ready   = $rs->getStatus()['readyReplicas'] ?? 0;

    // Check if this is the active RS (has ready replicas)
    $active = $ready > 0 ? 'ACTIVE' : 'old';
    echo "[$active] $name — $ready/$desired ready\n";
}
```

---

## Deployment Status Conditions

```php
$dep = $client->deployments()->findByName('my-app');
$conditions = $dep->getStatus()['conditions'] ?? [];

foreach ($conditions as $cond) {
    printf(
        "%-25s %-10s %s\n",
        $cond['type'],
        $cond['status'],
        $cond['message'] ?? ''
    );
}
```

| Condition Type | Meaning |
|---|---|
| `Available` | Minimum replicas are ready |
| `Progressing` | Rollout is in progress or complete |
| `ReplicaFailure` | Replica creation failed |

---

## Deployment Strategy Comparison

```php
// RollingUpdate (default) — no downtime
'strategy' => [
    'type' => 'RollingUpdate',
    'rollingUpdate' => [
        'maxUnavailable' => '25%',
        'maxSurge'       => '25%',
    ],
],

// Recreate — downtime between old and new
'strategy' => [
    'type' => 'Recreate',
],
```
