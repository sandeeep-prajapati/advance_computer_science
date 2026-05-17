# Chapter 7: ConfigMaps & Secrets

## ConfigMaps vs Secrets

| Feature | ConfigMap | Secret |
|---|---|---|
| Data type | Plain text | Base64-encoded |
| Stored in etcd | Plaintext | Encrypted (if etcd encryption enabled) |
| Use case | App config, feature flags | Passwords, tokens, TLS certs |
| Size limit | 1 MiB | 1 MiB |
| Mount as | File, env var | File, env var |

---

## Working with ConfigMaps

### List ConfigMaps

```php
use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig();

$configMaps = $client->configMaps()->find();

foreach ($configMaps as $cm) {
    $name = $cm->getMetadata('name');
    $keys = array_keys($cm->getAttribute('data') ?? []);
    echo "$name: " . implode(', ', $keys) . "\n";
}
```

### Get ConfigMap Data

```php
$cm = $client->configMaps()->findByName('app-config');

$data = $cm->getAttribute('data') ?? [];

echo "DB_HOST: " . ($data['DB_HOST'] ?? 'not set') . "\n";
echo "APP_ENV: " . ($data['APP_ENV'] ?? 'not set') . "\n";

// Read as file content (multi-line config)
$nginxConf = $data['nginx.conf'] ?? '';
file_put_contents('/tmp/nginx.conf', $nginxConf);
```

### Create a ConfigMap

```php
use Maclof\Kubernetes\Models\ConfigMap;

$configMap = new ConfigMap([
    'metadata' => [
        'name'      => 'app-config',
        'namespace' => 'default',
        'labels'    => ['app' => 'my-app'],
    ],
    'data' => [
        'APP_ENV'       => 'production',
        'DB_HOST'       => 'postgres.default.svc.cluster.local',
        'DB_PORT'       => '5432',
        'CACHE_TTL'     => '3600',
        'nginx.conf'    => "server {\n    listen 80;\n    server_name _;\n    location / { proxy_pass http://app:8080; }\n}\n",
    ],
]);

$client->configMaps()->create($configMap);
echo "ConfigMap created\n";
```

### Update a ConfigMap

```php
$cm = $client->configMaps()->findByName('app-config');

$data = $cm->getAttribute('data') ?? [];
$data['CACHE_TTL'] = '7200';   // update a value
$data['NEW_FLAG']  = 'true';   // add a new key

$cm->setAttribute('data', $data);
$client->configMaps()->update($cm);
echo "ConfigMap updated\n";
```

### Delete a ConfigMap

```php
$client->configMaps()->deleteByName('app-config');
```

---

## Working with Secrets

> **Security Rule:** Never log, echo, or expose Secret values. Always handle them with care.

### List Secrets (Names Only — Never Print Values)

```php
$secrets = $client->secrets()->find();

foreach ($secrets as $secret) {
    $name = $secret->getMetadata('name');
    $type = $secret->getAttribute('type') ?? 'Opaque';
    $keys = array_keys($secret->getAttribute('data') ?? []);
    echo "$name ($type): " . implode(', ', $keys) . "\n";
}
```

### Create an Opaque Secret

Kubernetes expects Secret data to be **base64-encoded**:

```php
use Maclof\Kubernetes\Models\Secret;

$secret = new Secret([
    'metadata' => [
        'name'      => 'db-credentials',
        'namespace' => 'default',
        'labels'    => ['app' => 'my-app'],
    ],
    'type' => 'Opaque',
    'data' => [
        'username' => base64_encode('admin'),
        'password' => base64_encode('super-secret-password'),
        'db_url'   => base64_encode('postgres://admin:pass@postgres:5432/mydb'),
    ],
]);

$client->secrets()->create($secret);
echo "Secret created\n";
```

### Create a Docker Registry Secret

```php
$dockerConfig = json_encode([
    'auths' => [
        'registry.example.com' => [
            'username' => 'myuser',
            'password' => 'mypass',
            'email'    => 'ci@example.com',
            'auth'     => base64_encode('myuser:mypass'),
        ],
    ],
]);

$secret = new Secret([
    'metadata' => ['name' => 'registry-credentials'],
    'type'     => 'kubernetes.io/dockerconfigjson',
    'data'     => [
        '.dockerconfigjson' => base64_encode($dockerConfig),
    ],
]);

$client->secrets()->create($secret);
```

### Create a TLS Secret

```php
$tlsSecret = new Secret([
    'metadata' => ['name' => 'my-tls-secret'],
    'type'     => 'kubernetes.io/tls',
    'data'     => [
        'tls.crt' => base64_encode(file_get_contents('/path/to/tls.crt')),
        'tls.key' => base64_encode(file_get_contents('/path/to/tls.key')),
    ],
]);

$client->secrets()->create($tlsSecret);
```

### Read a Secret Value

```php
$secret = $client->secrets()->findByName('db-credentials');

$data = $secret->getAttribute('data') ?? [];

// Values are base64-encoded in the API response — decode them
$username = base64_decode($data['username'] ?? '');
$password = base64_decode($data['password'] ?? '');

// Use, but don't log the password
$dsn = "pgsql:host=postgres;dbname=mydb;user=$username;password=$password";
$pdo = new PDO($dsn);
```

### Update a Secret

```php
$secret = $client->secrets()->findByName('db-credentials');

$data = $secret->getAttribute('data') ?? [];
$data['password'] = base64_encode('new-rotated-password');

$secret->setAttribute('data', $data);
$client->secrets()->update($secret);
echo "Secret rotated\n";
```

### Delete a Secret

```php
$client->secrets()->deleteByName('db-credentials');
```

---

## Secret Rotation Pattern

```php
function rotateSecret(Client $client, string $secretName, string $key, string $newValue): void
{
    $secret = $client->secrets()->findByName($secretName);
    $data   = $secret->getAttribute('data') ?? [];

    // Store previous value for rollback
    if (isset($data[$key])) {
        $data[$key . '_previous'] = $data[$key];
    }

    $data[$key] = base64_encode($newValue);
    $secret->setAttribute('data', $data);
    $client->secrets()->update($secret);

    echo "Secret '$secretName' key '$key' rotated successfully\n";
}

rotateSecret($client, 'db-credentials', 'password', 'new-secure-password-v2');
```

---

## Bulk ConfigMap from PHP Array

```php
function syncConfigMap(Client $client, string $name, array $config): void
{
    $data = array_map('strval', $config);  // ensure all values are strings

    try {
        $cm = $client->configMaps()->findByName($name);
        $cm->setAttribute('data', $data);
        $client->configMaps()->update($cm);
        echo "ConfigMap '$name' updated\n";
    } catch (\Exception $e) {
        $cm = new \Maclof\Kubernetes\Models\ConfigMap([
            'metadata' => ['name' => $name],
            'data'     => $data,
        ]);
        $client->configMaps()->create($cm);
        echo "ConfigMap '$name' created\n";
    }
}

syncConfigMap($client, 'feature-flags', [
    'FEATURE_NEW_UI'    => 'true',
    'FEATURE_DARK_MODE' => 'false',
    'MAX_RETRY'         => '3',
]);
```

---

## ConfigMap from .env File

```php
function configMapFromEnvFile(Client $client, string $name, string $envFilePath): void
{
    $lines = file($envFilePath, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);
    $data  = [];

    foreach ($lines as $line) {
        if (str_starts_with(trim($line), '#')) continue;
        [$key, $value] = explode('=', $line, 2);
        $data[trim($key)] = trim($value, '"\'');
    }

    syncConfigMap($client, $name, $data);
}

configMapFromEnvFile($client, 'app-env-config', '/path/to/.env.production');
```
