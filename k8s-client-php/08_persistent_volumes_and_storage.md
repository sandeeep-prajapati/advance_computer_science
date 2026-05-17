# Chapter 8: Persistent Volumes & Storage

## Storage Architecture

```
Developer defines:    PersistentVolumeClaim (PVC)
Cluster provides:     PersistentVolume (PV)
Cloud provisions:     Disk (EBS, GCE PD, Azure Disk) via StorageClass
```

| Object | Who Creates It | Purpose |
|---|---|---|
| `StorageClass` | Admin | Defines provisioner and disk type |
| `PersistentVolume` | Admin or Dynamically | The actual storage resource |
| `PersistentVolumeClaim` | Developer | Request for storage |

---

## Working with StorageClasses

```php
use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig();

// List available StorageClasses
$storageClasses = $client->storageClasses()->find();

foreach ($storageClasses as $sc) {
    $name        = $sc->getMetadata('name');
    $provisioner = $sc->getAttribute('provisioner');
    $isDefault   = ($sc->getMetadata('annotations')['storageclass.kubernetes.io/is-default-class'] ?? 'false') === 'true';
    $reclaimPolicy = $sc->getAttribute('reclaimPolicy') ?? 'Delete';

    printf("%-30s %-40s %-10s %s\n",
        $name,
        $provisioner,
        $reclaimPolicy,
        $isDefault ? '(default)' : ''
    );
}
```

### Create a StorageClass

```php
use Maclof\Kubernetes\Models\StorageClass;

$sc = new StorageClass([
    'metadata' => [
        'name'        => 'fast-ssd',
        'annotations' => [
            'storageclass.kubernetes.io/is-default-class' => 'false',
        ],
    ],
    'provisioner'       => 'kubernetes.io/gce-pd',
    'reclaimPolicy'     => 'Retain',
    'volumeBindingMode' => 'WaitForFirstConsumer',
    'parameters'        => [
        'type'              => 'pd-ssd',
        'replication-type'  => 'regional-pd',
        'zones'             => 'us-central1-a,us-central1-b',
    ],
]);

$client->storageClasses()->create($sc);
```

---

## Working with PersistentVolumeClaims (PVC)

### List PVCs

```php
$pvcs = $client->persistentVolumeClaims()->find();

foreach ($pvcs as $pvc) {
    $name         = $pvc->getMetadata('name');
    $status       = $pvc->getStatus()['phase'] ?? 'Unknown';
    $capacity     = $pvc->getStatus()['capacity']['storage'] ?? 'N/A';
    $storageClass = $pvc->getSpec()['storageClassName'] ?? '-';
    $accessModes  = implode(',', $pvc->getSpec()['accessModes'] ?? []);

    printf("%-30s %-10s %-10s %-20s %s\n",
        $name, $status, $capacity, $storageClass, $accessModes);
}
```

### Create a PVC

```php
use Maclof\Kubernetes\Models\PersistentVolumeClaim;

$pvc = new PersistentVolumeClaim([
    'metadata' => [
        'name'      => 'my-data-pvc',
        'namespace' => 'default',
        'labels'    => ['app' => 'my-app'],
    ],
    'spec' => [
        'accessModes'      => ['ReadWriteOnce'],
        'storageClassName' => 'standard',
        'resources'        => [
            'requests' => ['storage' => '10Gi'],
        ],
    ],
]);

$created = $client->persistentVolumeClaims()->create($pvc);
echo "PVC created: " . $created->getMetadata('name') . "\n";

// Wait for PVC to be Bound
$pvc = waitForPVCBound($client, 'my-data-pvc');
echo "PVC bound to PV: " . $pvc->getSpec()['volumeName'] . "\n";
```

```php
function waitForPVCBound(Client $client, string $pvcName, int $timeout = 60): mixed
{
    $start = time();
    while (time() - $start < $timeout) {
        $pvc    = $client->persistentVolumeClaims()->findByName($pvcName);
        $phase  = $pvc->getStatus()['phase'] ?? 'Pending';

        if ($phase === 'Bound') return $pvc;

        echo "PVC phase: $phase — waiting...\n";
        sleep(3);
    }
    throw new RuntimeException("PVC not bound within $timeout seconds");
}
```

### Create a PVC with Specific Access Mode

```php
// ReadWriteOnce — single node read/write (most disks)
// ReadOnlyMany  — multiple nodes read-only (NFS, etc.)
// ReadWriteMany — multiple nodes read/write (NFS, CephFS)

$sharedPvc = new PersistentVolumeClaim([
    'metadata' => ['name' => 'shared-data'],
    'spec' => [
        'accessModes'      => ['ReadWriteMany'],
        'storageClassName' => 'nfs-client',
        'resources'        => ['requests' => ['storage' => '50Gi']],
    ],
]);

$client->persistentVolumeClaims()->create($sharedPvc);
```

---

## Working with PersistentVolumes (PV)

### List PVs

```php
$pvs = $client->persistentVolumes()->find();

foreach ($pvs as $pv) {
    $name         = $pv->getMetadata('name');
    $capacity     = $pv->getSpec()['capacity']['storage'] ?? 'N/A';
    $accessModes  = implode(',', $pv->getSpec()['accessModes'] ?? []);
    $phase        = $pv->getStatus()['phase'] ?? 'Unknown';
    $claim        = $pv->getSpec()['claimRef']['name'] ?? '-';

    printf("%-30s %-8s %-20s %-10s %s\n",
        $name, $capacity, $accessModes, $phase, $claim);
}
```

### Create a PV (Static Provisioning)

```php
use Maclof\Kubernetes\Models\PersistentVolume;

$pv = new PersistentVolume([
    'metadata' => [
        'name'   => 'my-nfs-pv',
        'labels' => ['type' => 'nfs'],
    ],
    'spec' => [
        'capacity'        => ['storage' => '100Gi'],
        'accessModes'     => ['ReadWriteMany'],
        'reclaimPolicy'   => 'Retain',
        'storageClassName' => 'nfs',
        'nfs' => [
            'server' => '10.0.0.5',
            'path'   => '/exported/data',
        ],
    ],
]);

$client->persistentVolumes()->create($pv);
```

---

## Mounting a PVC in a Deployment

```php
use Maclof\Kubernetes\Models\Deployment;

$deployment = new Deployment([
    'metadata' => ['name' => 'app-with-storage'],
    'spec' => [
        'replicas' => 1,
        'selector' => ['matchLabels' => ['app' => 'app-with-storage']],
        'template' => [
            'metadata' => ['labels' => ['app' => 'app-with-storage']],
            'spec' => [
                'containers' => [[
                    'name'  => 'app',
                    'image' => 'myapp:1.0',
                    'volumeMounts' => [
                        [
                            'name'      => 'data-volume',
                            'mountPath' => '/var/data',
                        ],
                    ],
                ]],
                'volumes' => [
                    [
                        'name' => 'data-volume',
                        'persistentVolumeClaim' => [
                            'claimName' => 'my-data-pvc',
                        ],
                    ],
                ],
            ],
        ],
    ],
]);

$client->deployments()->create($deployment);
```

---

## Storage Capacity Check

```php
function listStorageUsage(Client $client): void
{
    $pvcs = $client->persistentVolumeClaims()->find();

    $total = 0;
    foreach ($pvcs as $pvc) {
        $name     = $pvc->getMetadata('name');
        $ns       = $pvc->getMetadata('namespace');
        $phase    = $pvc->getStatus()['phase'] ?? '-';
        $size     = $pvc->getSpec()['resources']['requests']['storage'] ?? '0';

        printf("%-30s %-15s %-10s %s\n", $name, $ns, $phase, $size);

        // Parse sizes (basic)
        if (str_ends_with($size, 'Gi')) {
            $total += (int) $size;
        }
    }

    echo "\nTotal storage requested: {$total}Gi\n";
}
```

---

## Deleting a PVC

```php
// PVC must not be in use by any pod
$client->persistentVolumeClaims()->deleteByName('my-data-pvc');

// The PV behavior depends on reclaimPolicy:
// - Delete: PV and underlying disk are deleted
// - Retain: PV remains, must be manually reclaimed
// - Recycle: PV is scrubbed (deprecated)
```
