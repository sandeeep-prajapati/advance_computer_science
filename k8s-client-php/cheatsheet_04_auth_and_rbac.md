# k8s-client-php — Auth & RBAC Cheatsheet

---

## Auth Methods at a Glance

```
Environment        Recommended Auth Method
─────────────────────────────────────────────────────
Local dev          kubeconfig  (Client::fromKubeConfig())
CI/CD pipeline     ServiceAccount token  (kubectl create token)
In-cluster app     Mounted SA token  (/var/run/secrets/...)
Cross-account AWS  IRSA + EKS token  (aws eks get-token)
GKE Workload ID    Workload Identity + SA token
```

---

## kubeconfig Auth

```php
// Default: uses ~/.kube/config, current-context
$client = Client::fromKubeConfig();

// Specific context
$client = Client::fromKubeConfig(context: 'prod-cluster');

// Specific kubeconfig file
$client = Client::fromKubeConfig(kubeConfigPath: '/tmp/custom.yaml');

// Get context names
// $ kubectl config get-contexts
// $ kubectl config current-context
```

---

## In-Cluster (Pod) Auth

```php
// Kubernetes mounts these automatically into every Pod:
$token = file_get_contents('/var/run/secrets/kubernetes.io/serviceaccount/token');
$ca    = '/var/run/secrets/kubernetes.io/serviceaccount/ca.crt';
$ns    = file_get_contents('/var/run/secrets/kubernetes.io/serviceaccount/namespace');

$client = new Client([
    'master'    => 'https://kubernetes.default.svc',
    'ca_cert'   => $ca,
    'token'     => $token,
    'namespace' => $ns,
]);
```

**Pod must have a ServiceAccount with the right RBAC permissions.**

---

## ServiceAccount Token (External)

```bash
# K8s 1.24+ — generate a short-lived token
kubectl create token my-sa -n default --duration=1h

# Or create a long-lived Secret-based token (older)
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: my-sa-token
  annotations:
    kubernetes.io/service-account.name: my-sa
type: kubernetes.io/service-account-token
EOF
kubectl get secret my-sa-token -o jsonpath='{.data.token}' | base64 -d
```

```php
$client = new Client([
    'master'  => 'https://k8s.example.com:6443',
    'ca_cert' => '/path/to/ca.crt',
    'token'   => 'eyJhbGci...',
]);
```

---

## RBAC — Common Roles

### Read-only access to pods and logs

```php
new Role([
    'metadata' => ['name' => 'pod-reader', 'namespace' => 'default'],
    'rules' => [[
        'apiGroups' => [''],
        'resources' => ['pods', 'pods/log', 'pods/status'],
        'verbs'     => ['get', 'list', 'watch'],
    ]],
])
```

### Deployment manager (create/update, no delete)

```php
new Role([
    'metadata' => ['name' => 'deployer', 'namespace' => 'production'],
    'rules' => [
        ['apiGroups' => ['apps'], 'resources' => ['deployments'], 'verbs' => ['get','list','watch','create','update','patch']],
        ['apiGroups' => [''],     'resources' => ['configmaps','secrets'], 'verbs' => ['get','list']],
    ],
])
```

### Full namespace admin

```php
new Role([
    'metadata' => ['name' => 'ns-admin', 'namespace' => 'team-alpha'],
    'rules' => [[
        'apiGroups' => ['', 'apps', 'batch', 'networking.k8s.io'],
        'resources' => ['*'],
        'verbs'     => ['*'],
    ]],
])
```

### Cluster-wide read-only (ClusterRole)

```php
new ClusterRole([
    'metadata' => ['name' => 'cluster-viewer'],
    'rules' => [
        ['apiGroups' => ['', 'apps', 'batch'], 'resources' => ['*'], 'verbs' => ['get', 'list', 'watch']],
        ['apiGroups' => [''],                  'resources' => ['nodes', 'namespaces', 'persistentvolumes'], 'verbs' => ['get', 'list']],
    ],
])
```

---

## RBAC — Common apiGroups + Resources

| Resource | apiGroup | Resource name |
|---|---|---|
| Pods, Services, ConfigMaps, Secrets | `""` (core) | `pods`, `services`, `configmaps`, `secrets` |
| Pod logs | `""` | `pods/log` |
| Pod exec | `""` | `pods/exec` |
| Deployments, ReplicaSets | `apps` | `deployments`, `replicasets` |
| StatefulSets, DaemonSets | `apps` | `statefulsets`, `daemonsets` |
| Jobs, CronJobs | `batch` | `jobs`, `cronjobs` |
| Ingresses, NetworkPolicies | `networking.k8s.io` | `ingresses`, `networkpolicies` |
| Roles, RoleBindings | `rbac.authorization.k8s.io` | `roles`, `rolebindings` |
| ClusterRoles, ClusterRoleBindings | `rbac.authorization.k8s.io` | `clusterroles`, `clusterrolebindings` |
| PVs, PVCs | `""` / `storage.k8s.io` | `persistentvolumes`, `persistentvolumeclaims` |
| StorageClasses | `storage.k8s.io` | `storageclasses` |
| Nodes | `""` | `nodes` |
| Namespaces | `""` | `namespaces` |
| CRDs | `apiextensions.k8s.io` | `customresourcedefinitions` |

---

## RBAC — All Verbs

| Verb | HTTP | Description |
|---|---|---|
| `get` | GET (single) | Read one resource |
| `list` | GET (collection) | List resources |
| `watch` | GET + `?watch=true` | Stream updates |
| `create` | POST | Create resource |
| `update` | PUT | Replace resource |
| `patch` | PATCH | Partial update |
| `delete` | DELETE | Delete resource |
| `deletecollection` | DELETE (collection) | Delete many |
| `*` | All | All verbs |

---

## Bind SA to Role — Full Example

```php
// 1. ServiceAccount
$client->serviceAccounts()->create(new ServiceAccount([
    'metadata' => ['name' => 'ci-sa', 'namespace' => 'default'],
]));

// 2. Role
$client->roles()->create(new Role([
    'metadata' => ['name' => 'ci-role', 'namespace' => 'default'],
    'rules'    => [
        ['apiGroups' => ['apps'], 'resources' => ['deployments'], 'verbs' => ['get','patch','update']],
        ['apiGroups' => [''],     'resources' => ['configmaps'],  'verbs' => ['get','create','update']],
    ],
]));

// 3. RoleBinding
$client->roleBindings()->create(new RoleBinding([
    'metadata' => ['name' => 'ci-binding', 'namespace' => 'default'],
    'subjects' => [['kind' => 'ServiceAccount', 'name' => 'ci-sa', 'namespace' => 'default']],
    'roleRef'  => ['apiGroup' => 'rbac.authorization.k8s.io', 'kind' => 'Role', 'name' => 'ci-role'],
]));

// 4. Get token for this SA
// $ kubectl create token ci-sa -n default
```

---

## Bind SA to ClusterRole (cluster-wide)

```php
$client->clusterRoleBindings()->create(new ClusterRoleBinding([
    'metadata' => ['name' => 'dashboard-viewer'],
    'subjects' => [['kind' => 'ServiceAccount', 'name' => 'dashboard-sa', 'namespace' => 'kube-system']],
    'roleRef'  => ['apiGroup' => 'rbac.authorization.k8s.io', 'kind' => 'ClusterRole', 'name' => 'cluster-viewer'],
]));
```

---

## Test RBAC Permissions

```bash
# Can this SA list pods?
kubectl auth can-i list pods --as=system:serviceaccount:default:ci-sa

# Can it create deployments?
kubectl auth can-i create deployments --as=system:serviceaccount:default:ci-sa -n production

# List all permissions for a SA
kubectl auth can-i --list --as=system:serviceaccount:default:ci-sa
```

---

## RBAC Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| `403 Forbidden` on pods | SA has no Role | Create Role + RoleBinding |
| `403 Forbidden` on nodes | Needs ClusterRole | Use ClusterRole + ClusterRoleBinding |
| `403 Forbidden` on another namespace | Role is namespaced | Create Role+RoleBinding in that namespace |
| Token expired | Short-lived token | `kubectl create token` again or use Secret-based token |
| `401 Unauthorized` | Wrong/missing token | Verify token with `kubectl auth whoami` |
