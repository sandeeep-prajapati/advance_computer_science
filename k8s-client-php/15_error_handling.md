# Chapter 15: Error Handling & Retry Logic

## Kubernetes API Error Codes

| HTTP Status | Meaning | Common Cause |
|---|---|---|
| `400 Bad Request` | Invalid resource spec | Malformed YAML/JSON |
| `401 Unauthorized` | Not authenticated | Expired/missing token |
| `403 Forbidden` | No permission | RBAC not configured |
| `404 Not Found` | Resource doesn't exist | Wrong name or namespace |
| `409 Conflict` | Resource already exists | Creating duplicate |
| `422 Unprocessable` | Validation failed | Invalid field values |
| `429 Too Many Requests` | Rate limited | Too many API calls |
| `500 Internal Server Error` | K8s API bug | Rare |
| `503 Service Unavailable` | API server overloaded | High load |

---

## Basic Exception Handling

```php
use Maclof\Kubernetes\Client;
use GuzzleHttp\Exception\ClientException;
use GuzzleHttp\Exception\ConnectException;
use GuzzleHttp\Exception\ServerException;

$client = Client::fromKubeConfig();

function createDeploymentSafe(Client $client, array $spec): bool
{
    try {
        $deployment = new \Maclof\Kubernetes\Models\Deployment($spec);
        $client->deployments()->create($deployment);
        return true;

    } catch (ClientException $e) {
        $code    = $e->getResponse()->getStatusCode();
        $body    = json_decode($e->getResponse()->getBody()->getContents(), true);
        $message = $body['message'] ?? $e->getMessage();

        switch ($code) {
            case 401:
                error_log("Authentication failed. Check your token/certificate.");
                break;
            case 403:
                error_log("Permission denied. Check RBAC for this ServiceAccount.");
                break;
            case 409:
                error_log("Resource already exists: " . ($body['details']['name'] ?? 'unknown'));
                break;
            case 422:
                error_log("Validation failed: $message");
                break;
            default:
                error_log("API error $code: $message");
        }
        return false;

    } catch (ConnectException $e) {
        error_log("Cannot connect to Kubernetes API: " . $e->getMessage());
        return false;

    } catch (ServerException $e) {
        error_log("Kubernetes API server error: " . $e->getMessage());
        return false;
    }
}
```

---

## Create-or-Update (Apply) Pattern

```php
function applyDeployment(Client $client, array $spec): string
{
    $name = $spec['metadata']['name'];

    try {
        $existing = $client->deployments()->findByName($name);

        if ($existing === null) {
            // Create
            $client->deployments()->create(new \Maclof\Kubernetes\Models\Deployment($spec));
            return 'created';
        } else {
            // Update
            $client->deployments()->patch($existing, $spec);
            return 'updated';
        }

    } catch (ClientException $e) {
        if ($e->getResponse()->getStatusCode() === 404) {
            $client->deployments()->create(new \Maclof\Kubernetes\Models\Deployment($spec));
            return 'created';
        }
        throw $e;
    }
}

$action = applyDeployment($client, $myDeploymentSpec);
echo "Deployment $action\n";
```

---

## Retry Logic with Exponential Backoff

```php
function withRetry(callable $operation, int $maxAttempts = 3, float $baseDelay = 1.0): mixed
{
    $attempt = 0;

    while (true) {
        $attempt++;
        try {
            return $operation();

        } catch (ServerException $e) {
            // 5xx errors are transient — retry
            if ($attempt >= $maxAttempts) throw $e;

            $delay = $baseDelay * (2 ** ($attempt - 1));  // 1s, 2s, 4s
            error_log("Server error (attempt $attempt/$maxAttempts), retrying in {$delay}s...");
            sleep((int) $delay);

        } catch (ClientException $e) {
            $code = $e->getResponse()->getStatusCode();

            if ($code === 429) {
                // Rate limited — respect Retry-After header
                $retryAfter = (int) ($e->getResponse()->getHeader('Retry-After')[0] ?? 5);
                error_log("Rate limited, waiting {$retryAfter}s...");
                sleep($retryAfter);
                continue;
            }

            // 4xx errors (except 429) are not retryable
            throw $e;

        } catch (ConnectException $e) {
            if ($attempt >= $maxAttempts) throw $e;
            $delay = $baseDelay * (2 ** ($attempt - 1));
            error_log("Connection failed (attempt $attempt), retrying in {$delay}s...");
            sleep((int) $delay);
        }
    }
}

// Usage
$pods = withRetry(fn() => $client->pods()->find(), maxAttempts: 3);
```

---

## Idempotent Resource Management

```php
class K8sResourceManager
{
    public function __construct(private Client $client) {}

    public function ensureNamespace(string $name): void
    {
        try {
            $this->client->namespaces()->findByName($name);
        } catch (\Exception $e) {
            $this->client->namespaces()->create(
                new \Maclof\Kubernetes\Models\NamespaceModel([
                    'metadata' => ['name' => $name],
                ])
            );
        }
    }

    public function ensureConfigMap(string $name, string $namespace, array $data): void
    {
        $this->client->setNamespace($namespace);

        try {
            $cm = $this->client->configMaps()->findByName($name);
            if ($cm !== null) {
                $cm->setAttribute('data', array_map('strval', $data));
                $this->client->configMaps()->update($cm);
                return;
            }
        } catch (\Exception $e) {
            // Not found — fall through to create
        }

        $this->client->configMaps()->create(
            new \Maclof\Kubernetes\Models\ConfigMap([
                'metadata' => ['name' => $name, 'namespace' => $namespace],
                'data'     => array_map('strval', $data),
            ])
        );
    }

    public function deleteIfExists(string $resourceType, string $name): bool
    {
        try {
            $this->client->$resourceType()->deleteByName($name);
            return true;
        } catch (ClientException $e) {
            if ($e->getResponse()->getStatusCode() === 404) {
                return false;  // already gone
            }
            throw $e;
        }
    }
}
```

---

## Validation Before Creating Resources

```php
function validateDeploymentSpec(array $spec): array
{
    $errors = [];

    if (empty($spec['metadata']['name'])) {
        $errors[] = 'metadata.name is required';
    }

    $replicas = $spec['spec']['replicas'] ?? null;
    if ($replicas !== null && (!is_int($replicas) || $replicas < 0)) {
        $errors[] = 'spec.replicas must be a non-negative integer';
    }

    $containers = $spec['spec']['template']['spec']['containers'] ?? [];
    if (empty($containers)) {
        $errors[] = 'spec.template.spec.containers must not be empty';
    }

    foreach ($containers as $i => $c) {
        if (empty($c['name'])) $errors[] = "containers[$i].name is required";
        if (empty($c['image'])) $errors[] = "containers[$i].image is required";
    }

    return $errors;
}

$errors = validateDeploymentSpec($mySpec);
if (!empty($errors)) {
    foreach ($errors as $err) {
        echo "Validation error: $err\n";
    }
    exit(1);
}

$client->deployments()->create(new Deployment($mySpec));
```

---

## Structured Error Logging

```php
function logK8sError(string $operation, string $resource, \Exception $e): void
{
    $context = [
        'operation'  => $operation,
        'resource'   => $resource,
        'error_type' => get_class($e),
        'message'    => $e->getMessage(),
    ];

    if ($e instanceof ClientException) {
        $context['http_status'] = $e->getResponse()->getStatusCode();
        $body = json_decode($e->getResponse()->getBody()->getContents(), true);
        $context['k8s_reason']  = $body['reason'] ?? null;
        $context['k8s_message'] = $body['message'] ?? null;
    }

    error_log(json_encode($context));
}

try {
    $client->deployments()->create($deployment);
} catch (\Exception $e) {
    logK8sError('create', 'Deployment/my-app', $e);
}
```

---

## Timeout Configuration

```php
$client = new Client([
    'master'  => 'https://k8s.example.com:6443',
    'token'   => '...',
    'ca_cert' => '...',
    'timeout' => 30,    // 30 second timeout per request
]);

// For long-running operations (watch, exec), increase timeout
$client = new Client([
    // ...
    'timeout'         => 0,      // no timeout for watches
    'connect_timeout' => 10,     // but still timeout on initial connect
]);
```

---

## Common Error Patterns and Fixes

```php
// Pattern: Resource already exists (409)
try {
    $client->secrets()->create($secret);
} catch (ClientException $e) {
    if ($e->getResponse()->getStatusCode() === 409) {
        // Update instead
        $existing = $client->secrets()->findByName($secret->getMetadata('name'));
        $client->secrets()->update($existing);
    }
}

// Pattern: Namespace not found
try {
    $client->setNamespace('team-xyz');
    $pods = $client->pods()->find();
} catch (ClientException $e) {
    if ($e->getResponse()->getStatusCode() === 404) {
        echo "Namespace 'team-xyz' does not exist\n";
    }
}

// Pattern: Forbidden — wrong ServiceAccount permissions
try {
    $client->nodes()->find();  // requires cluster-level access
} catch (ClientException $e) {
    if ($e->getResponse()->getStatusCode() === 403) {
        echo "Missing ClusterRole permission to list nodes\n";
    }
}
```
