# Chapter 10: Ingress

## What is Ingress?

**Ingress** is the HTTP/HTTPS layer router for Kubernetes. Instead of creating a LoadBalancer per service, one Ingress controller handles all external traffic and routes by hostname or path.

```
Internet
   │
   ▼
Ingress Controller (nginx, traefik, etc.)
   ├── /api/*       → api-service:80
   ├── /app/*       → frontend-service:80
   └── admin.host   → admin-service:8080
```

---

## Prerequisites

```bash
# Install nginx ingress controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/cloud/deploy.yaml

# For minikube
minikube addons enable ingress
```

---

## Listing Ingresses

```php
use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig();

$ingresses = $client->ingresses()->find();

foreach ($ingresses as $ing) {
    $name  = $ing->getMetadata('name');
    $rules = $ing->getSpec()['rules'] ?? [];

    foreach ($rules as $rule) {
        $host  = $rule['host'] ?? '*';
        $paths = $rule['http']['paths'] ?? [];
        foreach ($paths as $path) {
            $backend = $path['backend']['service']['name'] ?? 'N/A';
            $port    = $path['backend']['service']['port']['number'] ?? 'N/A';
            $p       = $path['path'] ?? '/';
            printf("%-20s %-25s %-10s → %s:%s\n", $name, $host, $p, $backend, $port);
        }
    }
}
```

---

## Creating a Basic Ingress (Path-based Routing)

```php
use Maclof\Kubernetes\Models\Ingress;

$ingress = new Ingress([
    'metadata' => [
        'name'      => 'my-app-ingress',
        'namespace' => 'default',
        'annotations' => [
            'nginx.ingress.kubernetes.io/rewrite-target' => '/',
        ],
    ],
    'spec' => [
        'ingressClassName' => 'nginx',
        'rules' => [
            [
                'host' => 'myapp.example.com',
                'http' => [
                    'paths' => [
                        [
                            'path'     => '/api',
                            'pathType' => 'Prefix',
                            'backend'  => [
                                'service' => [
                                    'name' => 'api-service',
                                    'port' => ['number' => 80],
                                ],
                            ],
                        ],
                        [
                            'path'     => '/',
                            'pathType' => 'Prefix',
                            'backend'  => [
                                'service' => [
                                    'name' => 'frontend-service',
                                    'port' => ['number' => 80],
                                ],
                            ],
                        ],
                    ],
                ],
            ],
        ],
    ],
]);

$client->ingresses()->create($ingress);
echo "Ingress created\n";
```

---

## Creating an Ingress with TLS

```php
$ingress = new Ingress([
    'metadata' => [
        'name'      => 'my-app-tls-ingress',
        'namespace' => 'default',
        'annotations' => [
            'cert-manager.io/cluster-issuer'            => 'letsencrypt-prod',
            'nginx.ingress.kubernetes.io/ssl-redirect'  => 'true',
        ],
    ],
    'spec' => [
        'ingressClassName' => 'nginx',
        'tls' => [
            [
                'hosts'      => ['myapp.example.com', 'api.example.com'],
                'secretName' => 'myapp-tls-secret',
            ],
        ],
        'rules' => [
            [
                'host' => 'myapp.example.com',
                'http' => [
                    'paths' => [[
                        'path'     => '/',
                        'pathType' => 'Prefix',
                        'backend'  => ['service' => ['name' => 'frontend', 'port' => ['number' => 80]]],
                    ]],
                ],
            ],
            [
                'host' => 'api.example.com',
                'http' => [
                    'paths' => [[
                        'path'     => '/',
                        'pathType' => 'Prefix',
                        'backend'  => ['service' => ['name' => 'api', 'port' => ['number' => 80]]],
                    ]],
                ],
            ],
        ],
    ],
]);

$client->ingresses()->create($ingress);
```

---

## Multi-Host Ingress

```php
function createMultiHostIngress(Client $client, array $routes): void
{
    // $routes format: ['host.example.com' => ['path' => 'service-name:port']]
    $rules = [];

    foreach ($routes as $host => $paths) {
        $httpPaths = [];
        foreach ($paths as $path => $backend) {
            [$svcName, $port] = explode(':', $backend);
            $httpPaths[] = [
                'path'     => $path,
                'pathType' => 'Prefix',
                'backend'  => ['service' => ['name' => $svcName, 'port' => ['number' => (int)$port]]],
            ];
        }
        $rules[] = ['host' => $host, 'http' => ['paths' => $httpPaths]];
    }

    $ingress = new Ingress([
        'metadata' => ['name' => 'multi-host-ingress'],
        'spec'     => ['ingressClassName' => 'nginx', 'rules' => $rules],
    ]);

    $client->ingresses()->create($ingress);
}

createMultiHostIngress($client, [
    'app.example.com'    => ['/' => 'frontend:80', '/api' => 'api:8080'],
    'admin.example.com'  => ['/' => 'admin-panel:3000'],
    'metrics.example.com' => ['/' => 'grafana:3000'],
]);
```

---

## Adding Annotations (nginx-specific)

```php
$annotations = [
    // Rate limiting
    'nginx.ingress.kubernetes.io/limit-rps'       => '10',
    'nginx.ingress.kubernetes.io/limit-connections' => '20',

    // Auth
    'nginx.ingress.kubernetes.io/auth-type'       => 'basic',
    'nginx.ingress.kubernetes.io/auth-secret'     => 'basic-auth',

    // Proxy settings
    'nginx.ingress.kubernetes.io/proxy-body-size'    => '50m',
    'nginx.ingress.kubernetes.io/proxy-read-timeout' => '300',

    // CORS
    'nginx.ingress.kubernetes.io/enable-cors'     => 'true',
    'nginx.ingress.kubernetes.io/cors-allow-origin' => 'https://myapp.com',

    // Redirect HTTP to HTTPS
    'nginx.ingress.kubernetes.io/ssl-redirect'    => 'true',
    'nginx.ingress.kubernetes.io/force-ssl-redirect' => 'true',
];

$ingress->setAttribute('metadata', array_merge(
    $ingress->getMetadata(),
    ['annotations' => $annotations]
));
```

---

## Getting Ingress Address

```php
$ing = $client->ingresses()->findByName('my-app-ingress');

$lbStatus = $ing->getStatus()['loadBalancer']['ingress'] ?? [];
foreach ($lbStatus as $entry) {
    echo "Address: " . ($entry['ip'] ?? $entry['hostname'] ?? 'pending') . "\n";
}

// Rules summary
foreach ($ing->getSpec()['rules'] ?? [] as $rule) {
    echo "Host: " . ($rule['host'] ?? '*') . "\n";
}
```

---

## Updating an Ingress

```php
$ing = $client->ingresses()->findByName('my-app-ingress');

$spec = $ing->getSpec();
$spec['rules'][] = [
    'host' => 'new-service.example.com',
    'http' => ['paths' => [[
        'path'     => '/',
        'pathType' => 'Prefix',
        'backend'  => ['service' => ['name' => 'new-service', 'port' => ['number' => 80]]],
    ]]],
];

$ing->setAttribute('spec', $spec);
$client->ingresses()->update($ing);
echo "Ingress updated with new host\n";
```

---

## Deleting an Ingress

```php
$client->ingresses()->deleteByName('my-app-ingress');
```

---

## Path Type Reference

| PathType | Behavior |
|---|---|
| `Exact` | Matches the exact path only |
| `Prefix` | Matches path prefix |
| `ImplementationSpecific` | Controller-specific behavior |
