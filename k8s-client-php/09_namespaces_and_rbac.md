# Chapter 9: Namespaces & RBAC

## Why Namespaces & RBAC?

| Problem | Solution |
|---|---|
| Isolate teams (dev/staging/prod) | Separate Namespaces |
| Prevent team A from deleting team B's resources | RBAC Roles |
| Give CI pipeline only create/read pod access | ServiceAccount + RoleBinding |
| Allow dashboard to read cluster-wide resources | ClusterRole + ClusterRoleBinding |

---

## Working with Namespaces

### List Namespaces

```php
use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig();
$client->setNamespace('');  // target all namespaces for cluster-scoped resources

$namespaces = $client->namespaces()->find();

foreach ($namespaces as $ns) {
    $name   = $ns->getMetadata('name');
    $status = $ns->getStatus()['phase'] ?? 'Unknown';
    echo "[$status] $name\n";
}
```

### Create a Namespace

```php
use Maclof\Kubernetes\Models\NamespaceModel;

$namespace = new NamespaceModel([
    'metadata' => [
        'name' => 'team-alpha',
        'labels' => [
            'team'        => 'alpha',
            'environment' => 'development',
        ],
        'annotations' => [
            'owner' => 'alice@example.com',
        ],
    ],
]);

$client->namespaces()->create($namespace);
echo "Namespace 'team-alpha' created\n";
```

### Delete a Namespace

```php
// Deleting a namespace deletes ALL resources inside it
$client->namespaces()->deleteByName('team-alpha');
// Termination takes time — pods are cleaned up first
```

### Watch Namespace Status (Terminating)

```php
function waitForNamespaceGone(Client $client, string $name, int $timeout = 120): void
{
    $start = time();
    while (time() - $start < $timeout) {
        $ns = $client->namespaces()->findByName($name);
        if ($ns === null) {
            echo "Namespace '$name' fully deleted\n";
            return;
        }
        $phase = $ns->getStatus()['phase'] ?? 'Active';
        echo "Namespace phase: $phase — waiting...\n";
        sleep(3);
    }
    throw new RuntimeException("Namespace '$name' still terminating after $timeout seconds");
}
```

---

## Working with ServiceAccounts

```php
use Maclof\Kubernetes\Models\ServiceAccount;

// Create a ServiceAccount for a CI pipeline
$sa = new ServiceAccount([
    'metadata' => [
        'name'      => 'ci-pipeline',
        'namespace' => 'default',
        'labels'    => ['purpose' => 'ci'],
        'annotations' => [
            'description' => 'ServiceAccount for GitHub Actions CI pipeline',
        ],
    ],
]);

$client->serviceAccounts()->create($sa);

// Generate a token for it (K8s 1.24+)
// Use kubectl: kubectl create token ci-pipeline
```

### List ServiceAccounts

```php
$serviceAccounts = $client->serviceAccounts()->find();

foreach ($serviceAccounts as $sa) {
    echo $sa->getMetadata('name') . "\n";
}
```

---

## Working with Roles (Namespaced)

A **Role** grants permissions within a single namespace.

```php
use Maclof\Kubernetes\Models\Role;

$role = new Role([
    'metadata' => [
        'name'      => 'pod-reader',
        'namespace' => 'default',
    ],
    'rules' => [
        [
            'apiGroups' => [''],           // '' = core API group
            'resources' => ['pods', 'pods/log'],
            'verbs'     => ['get', 'list', 'watch'],
        ],
    ],
]);

$client->roles()->create($role);
```

### Deployment Manager Role

```php
$role = new Role([
    'metadata' => ['name' => 'deployment-manager', 'namespace' => 'production'],
    'rules' => [
        [
            'apiGroups' => ['apps'],
            'resources' => ['deployments', 'replicasets'],
            'verbs'     => ['get', 'list', 'watch', 'create', 'update', 'patch'],
        ],
        [
            'apiGroups' => [''],
            'resources' => ['pods', 'services', 'configmaps'],
            'verbs'     => ['get', 'list', 'watch'],
        ],
    ],
]);

$client->setNamespace('production');
$client->roles()->create($role);
```

---

## Working with ClusterRoles (Cluster-Wide)

```php
use Maclof\Kubernetes\Models\ClusterRole;

$clusterRole = new ClusterRole([
    'metadata' => ['name' => 'cluster-reader'],
    'rules' => [
        [
            'apiGroups' => ['', 'apps', 'batch'],
            'resources' => ['*'],
            'verbs'     => ['get', 'list', 'watch'],
        ],
        [
            'apiGroups' => [''],
            'resources' => ['nodes', 'namespaces', 'persistentvolumes'],
            'verbs'     => ['get', 'list'],
        ],
    ],
]);

$client->clusterRoles()->create($clusterRole);
```

---

## Working with RoleBindings

Bind a Role to a user or ServiceAccount within a namespace:

```php
use Maclof\Kubernetes\Models\RoleBinding;

$roleBinding = new RoleBinding([
    'metadata' => [
        'name'      => 'ci-pipeline-binding',
        'namespace' => 'default',
    ],
    'subjects' => [
        [
            'kind'      => 'ServiceAccount',
            'name'      => 'ci-pipeline',
            'namespace' => 'default',
        ],
    ],
    'roleRef' => [
        'apiGroup' => 'rbac.authorization.k8s.io',
        'kind'     => 'Role',
        'name'     => 'deployment-manager',
    ],
]);

$client->roleBindings()->create($roleBinding);
echo "RoleBinding created: ci-pipeline → deployment-manager\n";
```

### Bind to a User (not ServiceAccount)

```php
$roleBinding = new RoleBinding([
    'metadata' => ['name' => 'alice-pod-reader', 'namespace' => 'team-alpha'],
    'subjects' => [
        [
            'kind'     => 'User',
            'name'     => 'alice@example.com',
            'apiGroup' => 'rbac.authorization.k8s.io',
        ],
    ],
    'roleRef' => [
        'apiGroup' => 'rbac.authorization.k8s.io',
        'kind'     => 'Role',
        'name'     => 'pod-reader',
    ],
]);

$client->roleBindings()->create($roleBinding);
```

---

## Working with ClusterRoleBindings

```php
use Maclof\Kubernetes\Models\ClusterRoleBinding;

$crb = new ClusterRoleBinding([
    'metadata' => ['name' => 'dashboard-cluster-reader'],
    'subjects' => [
        [
            'kind'      => 'ServiceAccount',
            'name'      => 'dashboard',
            'namespace' => 'kube-system',
        ],
    ],
    'roleRef' => [
        'apiGroup' => 'rbac.authorization.k8s.io',
        'kind'     => 'ClusterRole',
        'name'     => 'cluster-reader',
    ],
]);

$client->clusterRoleBindings()->create($crb);
```

---

## Multi-Tenant Namespace Provisioning

Automate creating everything needed for a new team:

```php
function provisionTenantNamespace(Client $client, string $teamName, string $ownerEmail): void
{
    $ns = $teamName . '-workspace';

    // 1. Create namespace
    $client->namespaces()->create(new NamespaceModel([
        'metadata' => [
            'name'        => $ns,
            'labels'      => ['team' => $teamName],
            'annotations' => ['owner' => $ownerEmail],
        ],
    ]));

    $client->setNamespace($ns);

    // 2. Create ServiceAccount
    $client->serviceAccounts()->create(new ServiceAccount([
        'metadata' => ['name' => "$teamName-sa", 'namespace' => $ns],
    ]));

    // 3. Create Role
    $client->roles()->create(new Role([
        'metadata' => ['name' => "$teamName-role", 'namespace' => $ns],
        'rules' => [
            ['apiGroups' => ['', 'apps', 'batch'], 'resources' => ['*'], 'verbs' => ['*']],
        ],
    ]));

    // 4. Bind SA to Role
    $client->roleBindings()->create(new RoleBinding([
        'metadata' => ['name' => "$teamName-binding", 'namespace' => $ns],
        'subjects' => [['kind' => 'ServiceAccount', 'name' => "$teamName-sa", 'namespace' => $ns]],
        'roleRef'  => ['apiGroup' => 'rbac.authorization.k8s.io', 'kind' => 'Role', 'name' => "$teamName-role"],
    ]));

    echo "Namespace '$ns' provisioned for team '$teamName'\n";
}
```

---

## Listing Roles and Bindings

```php
// List all roles in current namespace
$roles = $client->roles()->find();
foreach ($roles as $role) {
    echo "Role: " . $role->getMetadata('name') . "\n";
    foreach ($role->getAttribute('rules') as $rule) {
        echo "  Resources: " . implode(', ', $rule['resources'])
            . " | Verbs: " . implode(', ', $rule['verbs']) . "\n";
    }
}

// List all role bindings
$bindings = $client->roleBindings()->find();
foreach ($bindings as $binding) {
    $subjects = array_map(fn($s) => $s['name'], $binding->getAttribute('subjects'));
    echo $binding->getMetadata('name') . " → " . implode(', ', $subjects) . "\n";
}
```
