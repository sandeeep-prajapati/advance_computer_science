# k8s-client-php — Code Snippets Cheatsheet

Common copy-paste patterns for real-world tasks.

---

## Apply (Create or Update)

```php
function apply(Client $client, string $resource, string $name, array $spec): string
{
    try {
        $existing = $client->$resource()->findByName($name);
        if ($existing) {
            $client->$resource()->patch($existing, $spec);
            return 'updated';
        }
    } catch (\Exception $e) {}

    $modelClass = getModelClass($resource);
    $client->$resource()->create(new $modelClass($spec));
    return 'created';
}
```

---

## Wait for Pod Running

```php
function waitForPodRunning(Client $client, string $name, int $timeout = 60): bool
{
    $start = time();
    while (time() - $start < $timeout) {
        $pod   = $client->pods()->findByName($name);
        $phase = $pod?->getStatus()['phase'] ?? 'Pending';
        if ($phase === 'Running') return true;
        sleep(2);
    }
    return false;
}
```

---

## Wait for Deployment Ready

```php
function waitForDeployment(Client $client, string $name, int $timeout = 120): bool
{
    $start = time();
    while (time() - $start < $timeout) {
        $dep     = $client->deployments()->findByName($name);
        $desired = $dep->getSpec()['replicas'] ?? 0;
        $ready   = $dep->getStatus()['readyReplicas'] ?? 0;
        $updated = $dep->getStatus()['updatedReplicas'] ?? 0;
        if ($updated === $desired && $ready === $desired) return true;
        sleep(3);
    }
    return false;
}
```

---

## Wait for Job Done

```php
function waitForJob(Client $client, string $name, int $timeout = 300): string
{
    $start = time();
    while (time() - $start < $timeout) {
        $job          = $client->jobs()->findByName($name);
        $backoffLimit = $job->getSpec()['backoffLimit'] ?? 3;
        $succeeded    = $job->getStatus()['succeeded'] ?? 0;
        $failed       = $job->getStatus()['failed'] ?? 0;
        if ($succeeded > 0) return 'succeeded';
        if ($failed > $backoffLimit) return 'failed';
        sleep(5);
    }
    return 'timeout';
}
```

---

## Wait for PVC Bound

```php
function waitForPVC(Client $client, string $name, int $timeout = 60): bool
{
    $start = time();
    while (time() - $start < $timeout) {
        $pvc   = $client->persistentVolumeClaims()->findByName($name);
        $phase = $pvc?->getStatus()['phase'] ?? 'Pending';
        if ($phase === 'Bound') return true;
        sleep(3);
    }
    return false;
}
```

---

## Wait for LoadBalancer IP

```php
function waitForLBIP(Client $client, string $name, int $timeout = 120): ?string
{
    $start = time();
    while (time() - $start < $timeout) {
        $svc     = $client->services()->findByName($name);
        $ingress = $svc->getStatus()['loadBalancer']['ingress'] ?? [];
        if (!empty($ingress)) return $ingress[0]['ip'] ?? $ingress[0]['hostname'] ?? null;
        sleep(5);
    }
    return null;
}
```

---

## Scale Deployment

```php
function scale(Client $client, string $name, int $replicas): void
{
    $dep = $client->deployments()->findByName($name);
    $client->deployments()->patch($dep, ['spec' => ['replicas' => $replicas]]);
}

scale($client, 'my-app', 0);   // scale to zero
scale($client, 'my-app', 5);   // scale up
```

---

## Update Image (Trigger Rolling Update)

```php
function updateImage(Client $client, string $deployment, string $container, string $image): void
{
    $dep = $client->deployments()->findByName($deployment);
    $client->deployments()->patch($dep, [
        'spec' => ['template' => ['spec' => ['containers' => [
            ['name' => $container, 'image' => $image],
        ]]]],
    ]);
}

updateImage($client, 'my-app', 'app', 'myregistry/myapp:v2.0.0');
```

---

## Get Pods for a Deployment

```php
function getDeploymentPods(Client $client, string $deploymentName): array
{
    $dep      = $client->deployments()->findByName($deploymentName);
    $labels   = $dep->getSpec()['selector']['matchLabels'] ?? [];
    $pods     = $client->pods()->setLabelSelector($labels)->find();
    return iterator_to_array($pods);
}
```

---

## Get Pod Logs for a Deployment

```php
function getDeploymentLogs(Client $client, string $deploymentName, int $tailLines = 100): array
{
    $logs = [];
    foreach (getDeploymentPods($client, $deploymentName) as $pod) {
        $name         = $pod->getMetadata('name');
        $logs[$name]  = $client->pods()->logs($name, ['tailLines' => $tailLines]);
    }
    return $logs;
}
```

---

## Restart a Deployment (Rollout Restart)

```php
function restartDeployment(Client $client, string $name): void
{
    $dep = $client->deployments()->findByName($name);
    // Add restart annotation to force pod recreation
    $client->deployments()->patch($dep, [
        'spec' => ['template' => ['metadata' => ['annotations' => [
            'kubectl.kubernetes.io/restartedAt' => date('c'),
        ]]]],
    ]);
}
```

---

## Sync ConfigMap from PHP Array

```php
function syncConfigMap(Client $client, string $name, array $config): void
{
    $data = array_map('strval', $config);
    try {
        $cm = $client->configMaps()->findByName($name);
        $cm->setAttribute('data', $data);
        $client->configMaps()->update($cm);
    } catch (\Exception $e) {
        $client->configMaps()->create(new ConfigMap([
            'metadata' => ['name' => $name],
            'data'     => $data,
        ]));
    }
}

syncConfigMap($client, 'feature-flags', [
    'FEATURE_NEW_UI'    => 'true',
    'FEATURE_DARK_MODE' => 'false',
]);
```

---

## Rotate a Secret

```php
function rotateSecret(Client $client, string $name, string $key, string $newValue): void
{
    $secret = $client->secrets()->findByName($name);
    $data   = $secret->getAttribute('data') ?? [];
    $data[$key . '_previous'] = $data[$key] ?? '';
    $data[$key] = base64_encode($newValue);
    $secret->setAttribute('data', $data);
    $client->secrets()->update($secret);
}
```

---

## Provision Tenant Namespace

```php
function provisionTenant(Client $client, string $team): void
{
    $ns = "$team-workspace";

    // Namespace
    $client->namespaces()->create(new NamespaceModel([
        'metadata' => ['name' => $ns, 'labels' => ['team' => $team]],
    ]));

    $client->setNamespace($ns);

    // ServiceAccount
    $client->serviceAccounts()->create(new ServiceAccount([
        'metadata' => ['name' => "$team-sa", 'namespace' => $ns],
    ]));

    // Role (full access in namespace)
    $client->roles()->create(new Role([
        'metadata' => ['name' => "$team-role", 'namespace' => $ns],
        'rules'    => [['apiGroups' => ['', 'apps', 'batch'], 'resources' => ['*'], 'verbs' => ['*']]],
    ]));

    // Bind SA to Role
    $client->roleBindings()->create(new RoleBinding([
        'metadata' => ['name' => "$team-binding", 'namespace' => $ns],
        'subjects' => [['kind' => 'ServiceAccount', 'name' => "$team-sa", 'namespace' => $ns]],
        'roleRef'  => ['apiGroup' => 'rbac.authorization.k8s.io', 'kind' => 'Role', 'name' => "$team-role"],
    ]));
}
```

---

## Cluster-wide Pod Status Summary

```php
function podSummary(Client $client): void
{
    $client->setNamespace('');
    $phases = ['Running' => 0, 'Pending' => 0, 'Failed' => 0, 'Succeeded' => 0, 'Unknown' => 0];
    foreach ($client->pods()->find() as $pod) {
        $phase = $pod->getStatus()['phase'] ?? 'Unknown';
        $phases[$phase] = ($phases[$phase] ?? 0) + 1;
    }
    foreach ($phases as $phase => $count) {
        printf("%-12s %d\n", $phase, $count);
    }
}
```

---

## Delete All Pods with Label (Force Restart)

```php
function forceRestartPods(Client $client, array $labels): void
{
    $pods = $client->pods()->setLabelSelector($labels)->find();
    foreach ($pods as $pod) {
        $name = $pod->getMetadata('name');
        $client->pods()->delete($pod, ['gracePeriodSeconds' => 0]);
        echo "Deleted $name (will be recreated by ReplicaSet)\n";
    }
}

forceRestartPods($client, ['app' => 'my-app', 'env' => 'staging']);
```

---

## Retry with Exponential Backoff

```php
function withRetry(callable $fn, int $attempts = 3, float $delay = 1.0): mixed
{
    for ($i = 1; $i <= $attempts; $i++) {
        try {
            return $fn();
        } catch (\GuzzleHttp\Exception\ServerException $e) {
            if ($i === $attempts) throw $e;
            sleep((int) ($delay * 2 ** ($i - 1)));
        } catch (\GuzzleHttp\Exception\ClientException $e) {
            if ($e->getResponse()->getStatusCode() === 429) {
                sleep((int) ($e->getResponse()->getHeader('Retry-After')[0] ?? 5));
                continue;
            }
            throw $e;
        }
    }
}

$pods = withRetry(fn() => $client->pods()->find());
```

---

## Run a One-off Job and Collect Output

```php
function runJobAndGetLogs(Client $client, string $name, array $command, string $image): string
{
    $jobSpec = [
        'metadata' => ['name' => $name],
        'spec'     => [
            'backoffLimit' => 0,
            'ttlSecondsAfterFinished' => 600,
            'template' => ['spec' => [
                'restartPolicy' => 'Never',
                'containers'    => [['name' => 'runner', 'image' => $image, 'command' => $command]],
            ]],
        ],
    ];

    $client->jobs()->create(new Job($jobSpec));
    $result = waitForJob($client, $name, 120);

    $pods = $client->pods()->setLabelSelector(['job-name' => $name])->find();
    foreach ($pods as $pod) {
        return $client->pods()->logs($pod->getMetadata('name'));
    }

    return '';
}

$output = runJobAndGetLogs($client, 'php-task-' . time(), ['php', '-r', 'echo phpversion();'], 'php:8.2');
echo "Output: $output\n";
```

---

## Suspend / Resume All CronJobs

```php
function setCronJobsSuspended(Client $client, bool $suspend): void
{
    foreach ($client->cronJobs()->find() as $cj) {
        $client->cronJobs()->patch($cj, ['spec' => ['suspend' => $suspend]]);
        echo ($suspend ? 'Suspended' : 'Resumed') . ': ' . $cj->getMetadata('name') . "\n";
    }
}

setCronJobsSuspended($client, true);   // pause all before a deploy
setCronJobsSuspended($client, false);  // resume after
```

---

## Watch Pod Restarts (Alert on CrashLoop)

```php
function alertOnCrashLoop(Client $client, callable $alertFn): void
{
    $seen = [];

    $client->pods()->watch(function ($type, $pod) use (&$seen, $alertFn) {
        $name     = $pod->getMetadata('name');
        $statuses = $pod->getStatus()['containerStatuses'] ?? [];

        foreach ($statuses as $cs) {
            $key     = "$name/" . $cs['name'];
            $reason  = $cs['state']['waiting']['reason'] ?? null;
            $count   = $cs['restartCount'] ?? 0;

            if ($reason === 'CrashLoopBackOff' && !in_array($key, $seen)) {
                $alertFn("CrashLoopBackOff: Pod=$name restarts=$count");
                $seen[] = $key;
            }
        }
        return true;
    });
}

alertOnCrashLoop($client, function ($msg) {
    // Send to Slack, PagerDuty, email, etc.
    error_log("[ALERT] $msg");
});
```
