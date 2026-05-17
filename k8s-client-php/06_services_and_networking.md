# Chapter 6: Services & Networking

## Service Types

| Type | Access | Use Case |
|---|---|---|
| `ClusterIP` | Internal only (within cluster) | Microservice-to-microservice |
| `NodePort` | Node IP + port | Dev/testing, external access |
| `LoadBalancer` | Cloud load balancer IP | Production internet-facing |
| `ExternalName` | DNS alias to external service | Connect to external DB by name |

---

## Listing Services

```php
use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig();

$services = $client->services()->find();

foreach ($services as $svc) {
    $name      = $svc->getMetadata('name');
    $type      = $svc->getSpec()['type'] ?? 'ClusterIP';
    $clusterIP = $svc->getSpec()['clusterIP'] ?? 'None';
    $ports     = $svc->getSpec()['ports'] ?? [];

    $portStr = implode(',', array_map(
        fn($p) => ($p['port'] ?? '') . '/' . ($p['protocol'] ?? 'TCP'),
        $ports
    ));

    printf("%-35s %-15s %-16s %s\n", $name, $type, $clusterIP, $portStr);
}
```

---

## Creating a ClusterIP Service

```php
use Maclof\Kubernetes\Models\Service;

$service = new Service([
    'metadata' => [
        'name'      => 'my-app-svc',
        'namespace' => 'default',
        'labels'    => ['app' => 'my-app'],
    ],
    'spec' => [
        'type'     => 'ClusterIP',
        'selector' => ['app' => 'my-app'],
        'ports'    => [
            [
                'name'       => 'http',
                'port'       => 80,
                'targetPort' => 8080,
                'protocol'   => 'TCP',
            ],
        ],
    ],
]);

$client->services()->create($service);
echo "ClusterIP service created\n";
```

---

## Creating a NodePort Service

```php
$service = new Service([
    'metadata' => ['name' => 'my-app-nodeport'],
    'spec' => [
        'type'     => 'NodePort',
        'selector' => ['app' => 'my-app'],
        'ports'    => [
            [
                'port'       => 80,
                'targetPort' => 8080,
                'nodePort'   => 30080,   // 30000–32767
            ],
        ],
    ],
]);

$client->services()->create($service);
```

---

## Creating a LoadBalancer Service

```php
$service = new Service([
    'metadata' => [
        'name' => 'my-app-lb',
        'annotations' => [
            // GCP: Specify internal LB
            'cloud.google.com/load-balancer-type' => 'Internal',
            // AWS: Specify NLB type
            'service.beta.kubernetes.io/aws-load-balancer-type' => 'nlb',
        ],
    ],
    'spec' => [
        'type'     => 'LoadBalancer',
        'selector' => ['app' => 'my-app'],
        'ports'    => [
            ['port' => 80, 'targetPort' => 8080],
            ['port' => 443, 'targetPort' => 8443],
        ],
    ],
]);

$created = $client->services()->create($service);

// Wait for LoadBalancer IP to be assigned
$svc = waitForLoadBalancerIP($client, 'my-app-lb');
echo "External IP: " . $svc['status']['loadBalancer']['ingress'][0]['ip'] . "\n";
```

```php
function waitForLoadBalancerIP(Client $client, string $name, int $timeout = 120): array
{
    $start = time();
    while (time() - $start < $timeout) {
        $svc = $client->services()->findByName($name);
        $ingress = $svc->getStatus()['loadBalancer']['ingress'] ?? [];
        if (!empty($ingress)) {
            return $svc->toArray();
        }
        echo "Waiting for external IP...\n";
        sleep(5);
    }
    throw new RuntimeException("LoadBalancer IP not assigned within timeout");
}
```

---

## Creating an ExternalName Service

```php
$service = new Service([
    'metadata' => ['name' => 'external-db'],
    'spec' => [
        'type'         => 'ExternalName',
        'externalName' => 'my-database.rds.amazonaws.com',
    ],
]);

$client->services()->create($service);
// Now pods can connect to 'external-db' and it resolves to the RDS host
```

---

## Getting Service Details

```php
$svc = $client->services()->findByName('my-app-svc');

echo "Name:       " . $svc->getMetadata('name') . "\n";
echo "ClusterIP:  " . $svc->getSpec()['clusterIP'] . "\n";
echo "Type:       " . $svc->getSpec()['type'] . "\n";

foreach ($svc->getSpec()['ports'] as $port) {
    echo "Port: " . $port['port'] . " → " . $port['targetPort'] . "\n";
}

// External IPs for LoadBalancer
$ingress = $svc->getStatus()['loadBalancer']['ingress'] ?? [];
foreach ($ingress as $ing) {
    echo "External: " . ($ing['ip'] ?? $ing['hostname'] ?? 'pending') . "\n";
}
```

---

## Working with Endpoints

Endpoints show which Pod IPs a Service currently routes to:

```php
$endpoints = $client->endpoints()->findByName('my-app-svc');

$subsets = $endpoints->getAttribute('subsets') ?? [];
foreach ($subsets as $subset) {
    $port = $subset['ports'][0]['port'] ?? 'N/A';
    foreach ($subset['addresses'] ?? [] as $addr) {
        echo "Ready: " . $addr['ip'] . ":" . $port
            . " (Pod: " . ($addr['targetRef']['name'] ?? 'N/A') . ")\n";
    }
    foreach ($subset['notReadyAddresses'] ?? [] as $addr) {
        echo "NotReady: " . $addr['ip'] . "\n";
    }
}
```

---

## Deleting a Service

```php
$client->services()->deleteByName('my-app-svc');
```

---

## Headless Service (StatefulSet DNS)

```php
$service = new Service([
    'metadata' => ['name' => 'my-statefulset-headless'],
    'spec' => [
        'clusterIP' => 'None',   // headless — no VIP assigned
        'selector'  => ['app' => 'my-statefulset'],
        'ports'     => [['port' => 3306]],
    ],
]);

$client->services()->create($service);
// Each pod now has DNS: pod-0.my-statefulset-headless.default.svc.cluster.local
```

---

## Checking Service → Pod Connectivity

```php
function checkServiceHealth(Client $client, string $serviceName): array
{
    $svc       = $client->services()->findByName($serviceName);
    $endpoints = $client->endpoints()->findByName($serviceName);
    $selector  = $svc->getSpec()['selector'] ?? [];

    // Count pods matching selector
    $pods       = $client->pods()->setLabelSelector($selector)->find();
    $readyPods  = array_filter(iterator_to_array($pods), function ($pod) {
        $statuses = $pod->getStatus()['containerStatuses'] ?? [];
        return !empty(array_filter($statuses, fn($s) => $s['ready']));
    });

    $subsets       = $endpoints->getAttribute('subsets') ?? [];
    $endpointCount = array_sum(array_map(fn($s) => count($s['addresses'] ?? []), $subsets));

    return [
        'service'          => $serviceName,
        'type'             => $svc->getSpec()['type'],
        'clusterIP'        => $svc->getSpec()['clusterIP'],
        'total_pods'       => count(iterator_to_array($pods)),
        'ready_pods'       => count($readyPods),
        'endpoint_ips'     => $endpointCount,
    ];
}
```

---

## Network Policy Awareness

While PHP client cannot enforce network policies, you can read them:

```php
$netPolicies = $client->networkPolicies()->find();

foreach ($netPolicies as $np) {
    $name        = $np->getMetadata('name');
    $podSelector = json_encode($np->getSpec()['podSelector'] ?? []);
    echo "Policy: $name applies to: $podSelector\n";
}
```
