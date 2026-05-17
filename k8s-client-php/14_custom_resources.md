# Chapter 14: Custom Resources (CRDs)

## What Are CRDs?

**CustomResourceDefinitions (CRDs)** let you extend the Kubernetes API with your own resource types. Operators (like cert-manager, ArgoCD, Prometheus Operator) define CRDs to manage their own concepts.

```
Standard K8s API:  Pod, Deployment, Service
CRD examples:      Certificate, Application, PrometheusRule, IngressRoute
```

---

## Listing CRDs in the Cluster

```php
use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig();

// CRDs are cluster-scoped (not namespaced)
$crds = $client->customResourceDefinitions()->find();

foreach ($crds as $crd) {
    $name    = $crd->getMetadata('name');    // e.g. certificates.cert-manager.io
    $group   = $crd->getSpec()['group'];     // e.g. cert-manager.io
    $kinds   = array_column($crd->getSpec()['names'] ?? [], 'kind');
    $scope   = $crd->getSpec()['scope'];     // Namespaced or Cluster

    printf("%-55s %-30s %-12s\n", $name, $group, $scope);
}
```

---

## Reading Custom Resources via Raw HTTP

The `maclof/kubernetes-client` library handles standard resources natively. For CRDs, use the raw API or the `customResources` method:

```php
// Generic custom resource access
// Format: /apis/{group}/{version}/namespaces/{namespace}/{plural}

$response = $client->sendRequest(
    'GET',
    '/apis/cert-manager.io/v1/namespaces/default/certificates'
);

$body  = json_decode($response->getBody()->getContents(), true);
$certs = $body['items'] ?? [];

foreach ($certs as $cert) {
    $name   = $cert['metadata']['name'];
    $ready  = collect($cert['status']['conditions'] ?? [])
        ->firstWhere('type', 'Ready')['status'] ?? 'Unknown';
    $expiry = $cert['status']['notAfter'] ?? 'N/A';

    echo "$name — ready=$ready expires=$expiry\n";
}
```

---

## Using php-k8s for CRD Support

The `renoki-co/php-k8s` library has better CRD support via its `K8sResource` mechanism:

```php
use RenokiCo\PhpK8s\KubernetesCluster;
use RenokiCo\PhpK8s\Kinds\K8sResource;

$cluster = KubernetesCluster::fromKubeConfigFile('~/.kube/config');

// Define a custom resource
class K8sCertificate extends K8sResource
{
    protected static $kind        = 'Certificate';
    protected static $apiVersion  = 'cert-manager.io/v1';
    protected static $pluralName  = 'certificates';
}

// List all certificates
$certs = K8sCertificate::fromCluster($cluster)->all();

foreach ($certs as $cert) {
    echo $cert->getName() . "\n";
}
```

---

## Creating a Custom Resource

```php
// Create a cert-manager Certificate via raw HTTP
$certBody = [
    'apiVersion' => 'cert-manager.io/v1',
    'kind'       => 'Certificate',
    'metadata'   => [
        'name'      => 'myapp-tls',
        'namespace' => 'default',
    ],
    'spec' => [
        'secretName' => 'myapp-tls-secret',
        'issuerRef'  => [
            'name' => 'letsencrypt-prod',
            'kind' => 'ClusterIssuer',
        ],
        'dnsNames' => [
            'myapp.example.com',
            'www.myapp.example.com',
        ],
        'duration'    => '2160h',   // 90 days
        'renewBefore' => '360h',    // renew 15 days before expiry
    ],
];

$client->sendRequest(
    'POST',
    '/apis/cert-manager.io/v1/namespaces/default/certificates',
    ['json' => $certBody]
);

echo "Certificate resource created\n";
```

---

## Working with ArgoCD Applications

```php
function createArgoCDApp(Client $client, string $name, string $repoUrl, string $path): void
{
    $app = [
        'apiVersion' => 'argoproj.io/v1alpha1',
        'kind'       => 'Application',
        'metadata'   => [
            'name'      => $name,
            'namespace' => 'argocd',
            'finalizers' => ['resources-finalizer.argocd.argoproj.io'],
        ],
        'spec' => [
            'project' => 'default',
            'source'  => [
                'repoURL'        => $repoUrl,
                'targetRevision' => 'HEAD',
                'path'           => $path,
            ],
            'destination' => [
                'server'    => 'https://kubernetes.default.svc',
                'namespace' => 'default',
            ],
            'syncPolicy' => [
                'automated' => [
                    'prune'    => true,
                    'selfHeal' => true,
                ],
                'syncOptions' => ['CreateNamespace=true'],
            ],
        ],
    ];

    $client->sendRequest(
        'POST',
        '/apis/argoproj.io/v1alpha1/namespaces/argocd/applications',
        ['json' => $app]
    );

    echo "ArgoCD Application '$name' created\n";
}
```

---

## Watching Custom Resources

```php
// Watch cert-manager certificates for status changes
$response = $client->sendRequest(
    'GET',
    '/apis/cert-manager.io/v1/namespaces/default/certificates',
    ['query' => ['watch' => 'true', 'timeoutSeconds' => '60']]
);

$body = $response->getBody();
while (!$body->eof()) {
    $line = $body->read(4096);
    if (empty(trim($line))) continue;

    $event = json_decode($line, true);
    $type  = $event['type'] ?? 'UNKNOWN';
    $name  = $event['object']['metadata']['name'] ?? '-';
    $ready = 'Unknown';

    foreach ($event['object']['status']['conditions'] ?? [] as $cond) {
        if ($cond['type'] === 'Ready') {
            $ready = $cond['status'];
            break;
        }
    }

    echo "[$type] Certificate: $name — Ready: $ready\n";
}
```

---

## Defining Your Own CRD via PHP

```php
function createCRD(Client $client, string $group, string $plural, string $kind): void
{
    $crd = [
        'apiVersion' => 'apiextensions.k8s.io/v1',
        'kind'       => 'CustomResourceDefinition',
        'metadata'   => ['name' => "$plural.$group"],
        'spec'       => [
            'group' => $group,
            'scope' => 'Namespaced',
            'names' => [
                'plural'   => $plural,
                'singular' => strtolower($kind),
                'kind'     => $kind,
            ],
            'versions' => [[
                'name'    => 'v1',
                'served'  => true,
                'storage' => true,
                'schema'  => [
                    'openAPIV3Schema' => [
                        'type'       => 'object',
                        'properties' => [
                            'spec' => [
                                'type'       => 'object',
                                'properties' => [
                                    'replicas' => ['type' => 'integer'],
                                    'image'    => ['type' => 'string'],
                                ],
                            ],
                        ],
                    ],
                ],
            ]],
        ],
    ];

    $client->sendRequest(
        'POST',
        '/apis/apiextensions.k8s.io/v1/customresourcedefinitions',
        ['json' => $crd]
    );

    echo "CRD '$plural.$group' created\n";
}

createCRD($client, 'mycompany.io', 'myapps', 'MyApp');
```

---

## List All Custom Resources of a Type

```php
function listCustomResources(Client $client, string $group, string $version, string $plural, string $namespace = 'default'): array
{
    $path = "/apis/$group/$version/namespaces/$namespace/$plural";

    $response = $client->sendRequest('GET', $path);
    $body     = json_decode($response->getBody()->getContents(), true);

    return $body['items'] ?? [];
}

$myApps = listCustomResources($client, 'mycompany.io', 'v1', 'myapps');
foreach ($myApps as $app) {
    echo $app['metadata']['name'] . "\n";
}
```

---

## CRD Reference Table

| CRD (name.group) | Kind | Used By |
|---|---|---|
| `certificates.cert-manager.io` | Certificate | cert-manager |
| `clusterissuers.cert-manager.io` | ClusterIssuer | cert-manager |
| `applications.argoproj.io` | Application | ArgoCD |
| `prometheusrules.monitoring.coreos.com` | PrometheusRule | Prometheus Operator |
| `ingressroutes.traefik.containo.us` | IngressRoute | Traefik |
| `helmreleases.helm.toolkit.fluxcd.io` | HelmRelease | FluxCD |
