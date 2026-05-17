# Chapter 16: Production Patterns

## 1. Operator Pattern in PHP

A Kubernetes Operator is a controller that watches custom resources and reconciles the actual state toward the desired state. While Go is the standard, PHP can implement basic operators for internal SaaS platforms.

```
Watch CRD events
     │
     ▼
Reconcile loop:
  - Read desired state from CRD
  - Read actual state from cluster
  - Create/update/delete resources to match
     │
     ▼
Repeat forever
```

```php
use Maclof\Kubernetes\Client;

class TenantOperator
{
    public function __construct(private Client $client) {}

    public function run(): void
    {
        echo "Tenant operator started\n";

        // Watch Tenant CRDs (assuming CRD is defined)
        while (true) {
            try {
                $this->client->sendRequest('GET', '/apis/mycompany.io/v1/tenants', [
                    'query' => ['watch' => 'true', 'timeoutSeconds' => '3600'],
                    'stream' => true,
                ])->then(function ($response) {
                    $body = $response->getBody();
                    while (!$body->eof()) {
                        $line  = $body->readline();
                        $event = json_decode($line, true);
                        if ($event) {
                            $this->reconcile($event['type'], $event['object']);
                        }
                    }
                });
            } catch (\Exception $e) {
                error_log("Watch error: " . $e->getMessage());
                sleep(5);  // backoff before reconnecting
            }
        }
    }

    private function reconcile(string $type, array $tenant): void
    {
        $name  = $tenant['metadata']['name'];
        $ns    = $tenant['spec']['namespace'] ?? $name . '-workspace';
        $owner = $tenant['spec']['owner'] ?? '';

        echo "[$type] Reconciling tenant: $name → namespace: $ns\n";

        if ($type === 'DELETED') {
            $this->deprovisionTenant($ns);
            return;
        }

        $this->provisionTenant($ns, $name, $owner);
    }

    private function provisionTenant(string $ns, string $team, string $owner): void
    {
        // Ensure namespace exists
        $this->ensureNamespace($ns, $team, $owner);

        // Ensure resource quota
        $this->ensureResourceQuota($ns);

        // Ensure limit range
        $this->ensureLimitRange($ns);

        echo "Tenant '$team' provisioned in namespace '$ns'\n";
    }

    private function deprovisionTenant(string $ns): void
    {
        try {
            $this->client->namespaces()->deleteByName($ns);
            echo "Namespace '$ns' deleted\n";
        } catch (\Exception $e) {
            error_log("Failed to delete namespace '$ns': " . $e->getMessage());
        }
    }

    private function ensureNamespace(string $name, string $team, string $owner): void
    {
        try {
            $this->client->namespaces()->findByName($name);
        } catch (\Exception $e) {
            $this->client->namespaces()->create(new \Maclof\Kubernetes\Models\NamespaceModel([
                'metadata' => [
                    'name'        => $name,
                    'labels'      => ['team' => $team, 'managed-by' => 'tenant-operator'],
                    'annotations' => ['owner' => $owner],
                ],
            ]));
        }
    }

    private function ensureResourceQuota(string $ns): void
    {
        $this->client->setNamespace($ns);

        // Basic quota: 4 CPU, 8Gi RAM, 10 pods
        $quota = [
            'metadata' => ['name' => 'tenant-quota', 'namespace' => $ns],
            'spec' => ['hard' => [
                'requests.cpu'    => '4',
                'requests.memory' => '8Gi',
                'limits.cpu'      => '8',
                'limits.memory'   => '16Gi',
                'pods'            => '20',
            ]],
        ];

        try {
            $existing = $this->client->resourceQuotas()->findByName('tenant-quota');
            $existing->setAttribute('spec', $quota['spec']);
            $this->client->resourceQuotas()->update($existing);
        } catch (\Exception $e) {
            $this->client->resourceQuotas()->create(
                new \Maclof\Kubernetes\Models\ResourceQuota($quota)
            );
        }
    }

    private function ensureLimitRange(string $ns): void
    {
        $this->client->setNamespace($ns);

        $limitRange = [
            'metadata' => ['name' => 'tenant-limits', 'namespace' => $ns],
            'spec' => ['limits' => [[
                'type'           => 'Container',
                'defaultRequest' => ['cpu' => '100m', 'memory' => '128Mi'],
                'default'        => ['cpu' => '500m', 'memory' => '256Mi'],
            ]]],
        ];

        // create or update...
    }
}

// Entry point
$client   = Client::fromKubeConfig();
$operator = new TenantOperator($client);
$operator->run();
```

---

## 2. Multi-Cluster Management

```php
class MultiClusterClient
{
    private array $clients = [];

    public function addCluster(string $name, Client $client): void
    {
        $this->clients[$name] = $client;
    }

    public function forCluster(string $name): Client
    {
        return $this->clients[$name] ?? throw new InvalidArgumentException("Unknown cluster: $name");
    }

    public function allClusters(): array
    {
        return array_keys($this->clients);
    }

    public function aggregatePods(string $namespace = 'default'): array
    {
        $result = [];
        foreach ($this->clients as $clusterName => $client) {
            $client->setNamespace($namespace);
            foreach ($client->pods()->find() as $pod) {
                $result[] = [
                    'cluster' => $clusterName,
                    'name'    => $pod->getMetadata('name'),
                    'phase'   => $pod->getStatus()['phase'] ?? 'Unknown',
                ];
            }
        }
        return $result;
    }
}

$manager = new MultiClusterClient();
$manager->addCluster('us-east', Client::fromKubeConfig(context: 'eks-us-east'));
$manager->addCluster('eu-west', Client::fromKubeConfig(context: 'gke-eu-west'));
$manager->addCluster('ap-south', Client::fromKubeConfig(context: 'aks-ap-south'));

$allPods = $manager->aggregatePods('production');
echo "Total pods across all clusters: " . count($allPods) . "\n";
```

---

## 3. Deployment Pipeline Integration

```php
class KubernetesDeployer
{
    public function __construct(private Client $client) {}

    public function deploy(string $image, string $deploymentName, string $containerName): array
    {
        $startTime = time();

        // 1. Update image
        $dep = $this->client->deployments()->findByName($deploymentName);
        $this->client->deployments()->patch($dep, [
            'spec' => ['template' => ['spec' => ['containers' => [
                ['name' => $containerName, 'image' => $image],
            ]]]],
        ]);

        echo "Rolling out $image to $deploymentName...\n";

        // 2. Wait for rollout
        $success = $this->waitForRollout($deploymentName);
        $elapsed = time() - $startTime;

        return [
            'success'  => $success,
            'image'    => $image,
            'duration' => $elapsed,
        ];
    }

    private function waitForRollout(string $name, int $timeout = 300): bool
    {
        $start = time();
        while (time() - $start < $timeout) {
            $dep       = $this->client->deployments()->findByName($name);
            $status    = $dep->getStatus();
            $desired   = $dep->getSpec()['replicas'] ?? 0;
            $updated   = $status['updatedReplicas'] ?? 0;
            $ready     = $status['readyReplicas'] ?? 0;
            $available = $status['availableReplicas'] ?? 0;

            printf("  updated=%-3d ready=%-3d available=%-3d desired=%-3d\n",
                $updated, $ready, $available, $desired);

            if ($updated === $desired && $ready === $desired && $available === $desired) {
                return true;
            }

            sleep(3);
        }
        return false;
    }

    public function rollback(string $deploymentName): void
    {
        // Trigger rollback by removing the latest ReplicaSet
        // In practice: patch with previous image or use revision annotation
        $dep = $this->client->deployments()->findByName($deploymentName);
        $this->client->deployments()->patch($dep, [
            'metadata' => ['annotations' => ['deployment.kubernetes.io/revision' => null]],
        ]);
        echo "Rollback triggered for $deploymentName\n";
    }
}

// CI/CD usage
$deployer = new KubernetesDeployer($client);
$result   = $deployer->deploy('myapp:2.5.0', 'my-app', 'app');

if (!$result['success']) {
    echo "Deploy failed after {$result['duration']}s — rolling back\n";
    $deployer->rollback('my-app');
    exit(1);
}

echo "Deploy succeeded in {$result['duration']}s\n";
```

---

## 4. Cluster Health Dashboard

```php
function clusterHealthReport(Client $client): array
{
    $client->setNamespace('');  // all namespaces

    // Nodes
    $nodes       = $client->nodes()->find();
    $readyNodes  = 0;
    foreach ($nodes as $node) {
        $conditions = $node->getStatus()['conditions'] ?? [];
        foreach ($conditions as $cond) {
            if ($cond['type'] === 'Ready' && $cond['status'] === 'True') {
                $readyNodes++;
                break;
            }
        }
    }

    // Pods
    $pods        = $client->pods()->find();
    $runningPods = 0;
    $failedPods  = 0;
    foreach ($pods as $pod) {
        $phase = $pod->getStatus()['phase'] ?? 'Unknown';
        if ($phase === 'Running') $runningPods++;
        if ($phase === 'Failed')  $failedPods++;
    }

    // Deployments
    $deployments       = $client->deployments()->find();
    $healthyDeployments = 0;
    foreach ($deployments as $dep) {
        $desired = $dep->getSpec()['replicas'] ?? 0;
        $ready   = $dep->getStatus()['readyReplicas'] ?? 0;
        if ($ready === $desired) $healthyDeployments++;
    }

    return [
        'nodes'        => ['total' => count($nodes), 'ready' => $readyNodes],
        'pods'         => ['total' => count($pods), 'running' => $runningPods, 'failed' => $failedPods],
        'deployments'  => ['total' => count($deployments), 'healthy' => $healthyDeployments],
        'timestamp'    => date('c'),
    ];
}

$health = clusterHealthReport($client);
echo json_encode($health, JSON_PRETTY_PRINT);
```

---

## 5. Canary Deployment Pattern

```php
function canaryDeploy(Client $client, string $app, string $newImage, int $canaryWeight = 10): void
{
    // Scale down canary deployment to ~10% of traffic weight
    // Requires a traffic management solution (Istio, Argo Rollouts, etc.)
    // This example manages it with two Deployments + Services

    $totalReplicas  = 10;
    $canaryReplicas = (int) ceil($totalReplicas * ($canaryWeight / 100));
    $stableReplicas = $totalReplicas - $canaryReplicas;

    // Patch stable deployment
    $stable = $client->deployments()->findByName($app);
    $client->deployments()->patch($stable, ['spec' => ['replicas' => $stableReplicas]]);

    // Create or update canary deployment
    $canaryName = $app . '-canary';
    try {
        $canary = $client->deployments()->findByName($canaryName);
        $client->deployments()->patch($canary, [
            'spec' => [
                'replicas' => $canaryReplicas,
                'template' => ['spec' => ['containers' => [['name' => $app, 'image' => $newImage]]]],
            ],
        ]);
    } catch (\Exception $e) {
        // Create canary from stable spec
        $stableSpec         = $stable->getSpec();
        $stableSpec['replicas'] = $canaryReplicas;
        $stableSpec['template']['spec']['containers'][0]['image'] = $newImage;
        $stableSpec['selector']['matchLabels']['track'] = 'canary';
        $stableSpec['template']['metadata']['labels']['track'] = 'canary';

        $client->deployments()->create(new \Maclof\Kubernetes\Models\Deployment([
            'metadata' => ['name' => $canaryName, 'labels' => ['app' => $app, 'track' => 'canary']],
            'spec'     => $stableSpec,
        ]));
    }

    echo "Canary deployed: $canaryReplicas/{$totalReplicas} pods running $newImage\n";
}

// Promote canary → full rollout
function promoteCanary(Client $client, string $app, string $newImage): void
{
    $totalReplicas = 10;

    $stable = $client->deployments()->findByName($app);
    $client->deployments()->patch($stable, [
        'spec' => [
            'replicas' => $totalReplicas,
            'template' => ['spec' => ['containers' => [['name' => $app, 'image' => $newImage]]]],
        ],
    ]);

    $client->deployments()->deleteByName($app . '-canary', ['propagationPolicy' => 'Foreground']);
    echo "Canary promoted — full rollout complete\n";
}
```

---

## 6. Automated Cleanup Job

```php
function cleanupCompletedJobs(Client $client, int $keepSuccessful = 3, int $keepFailed = 1): void
{
    $jobs = $client->jobs()->find();
    $succeeded = $failed = [];

    foreach ($jobs as $job) {
        $status = $job->getStatus();
        if (($status['succeeded'] ?? 0) > 0) {
            $succeeded[] = $job;
        } elseif (($status['failed'] ?? 0) > 0) {
            $failed[] = $job;
        }
    }

    // Sort oldest first
    usort($succeeded, fn($a, $b) =>
        strtotime($a->getStatus()['completionTime'] ?? 0)
        <=> strtotime($b->getStatus()['completionTime'] ?? 0)
    );

    $toDelete = array_slice($succeeded, 0, max(0, count($succeeded) - $keepSuccessful));
    foreach ($toDelete as $job) {
        $client->jobs()->delete($job, ['propagationPolicy' => 'Foreground']);
        echo "Deleted completed job: " . $job->getMetadata('name') . "\n";
    }

    $toDeleteFailed = array_slice($failed, 0, max(0, count($failed) - $keepFailed));
    foreach ($toDeleteFailed as $job) {
        $client->jobs()->delete($job, ['propagationPolicy' => 'Foreground']);
        echo "Deleted failed job: " . $job->getMetadata('name') . "\n";
    }
}

cleanupCompletedJobs($client, keepSuccessful: 5, keepFailed: 2);
```

---

## Production Checklist

| Concern | Solution |
|---|---|
| Authentication | Use in-cluster ServiceAccount tokens, not hardcoded certs |
| Least privilege | ServiceAccount with only the Roles it needs |
| Retry logic | Exponential backoff for 5xx and 429 errors |
| Timeout | Always set connection + request timeouts |
| Logging | Structured logs with cluster, namespace, resource context |
| Secret handling | Never log Secret values; use `valueFrom.secretKeyRef` in pod specs |
| Multi-cluster | Separate Client instances per cluster, aggregate carefully |
| Rate limiting | Batch API calls; avoid listing all resources on every request |
| Watch reconnection | Auto-restart watch on EOF or timeout |
| idempotency | Use create-or-update pattern for all resource management |
