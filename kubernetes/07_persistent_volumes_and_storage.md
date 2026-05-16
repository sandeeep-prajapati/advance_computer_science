# Chapter 7: Persistent Volumes & Storage

## Why Persistent Storage?

Containers are **ephemeral** — when a pod dies, all data inside it is lost. For databases, file uploads, logs, and any stateful data, you need persistent storage that outlives the pod.

```
Pod dies → Container filesystem gone
Pod with PVC → Data survives in Persistent Volume
```

---

## Storage Hierarchy

```
StorageClass
    ↓ (defines provisioner + parameters)
PersistentVolume (PV)
    ↓ (actual storage resource)
PersistentVolumeClaim (PVC)
    ↓ (request for storage)
Pod → mounts the PVC
```

---

## PersistentVolume (PV)

A **PersistentVolume** is a piece of storage in the cluster provisioned by an admin (or automatically by a StorageClass). It has:
- A **capacity** (size)
- An **access mode**
- A **reclaim policy**

### Access Modes

| Mode | Abbreviation | Meaning |
|---|---|---|
| `ReadWriteOnce` | RWO | One node can read/write |
| `ReadOnlyMany` | ROX | Many nodes can read |
| `ReadWriteMany` | RWX | Many nodes can read/write |
| `ReadWriteOncePod` | RWOP | Only one pod can read/write |

### Static PV Example

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: my-pv
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain  # Keep data after PVC deleted
  storageClassName: manual
  hostPath:                               # Only for single-node dev!
    path: /mnt/data
```

---

## PersistentVolumeClaim (PVC)

A **PVC** is a request for storage by a user. Kubernetes finds a matching PV and binds them.

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-pvc
  namespace: default
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: manual   # Match the PV's storageClassName
  resources:
    requests:
      storage: 5Gi           # Request 5Gi (must be ≤ PV capacity)
```

### Mount PVC in a Pod

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
        - name: mysql
          image: mysql:8.0
          env:
            - name: MYSQL_ROOT_PASSWORD
              value: "rootpassword"
          volumeMounts:
            - name: mysql-data
              mountPath: /var/lib/mysql    # MySQL data directory
      volumes:
        - name: mysql-data
          persistentVolumeClaim:
            claimName: my-pvc             # Use the PVC
```

---

## StorageClass & Dynamic Provisioning

**StorageClass** enables **dynamic provisioning** — PVs are created automatically when a PVC is created, without admin intervention.

### Built-in StorageClasses

```bash
kubectl get storageclasses

# Output (example on cloud):
# NAME            PROVISIONER              RECLAIMPOLICY
# standard        kubernetes.io/gce-pd     Delete
# premium-rwo     pd.csi.storage.gke.io    Delete
```

### Custom StorageClass

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: pd.csi.storage.gke.io   # GCP provisioner
parameters:
  type: pd-ssd
reclaimPolicy: Retain
allowVolumeExpansion: true
```

### Dynamic PVC (No Manual PV Needed)

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: data-pvc
spec:
  storageClassName: fast-ssd     # Triggers dynamic provisioning
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 20Gi
```

---

## Reclaim Policies

| Policy | What Happens When PVC is Deleted |
|---|---|
| `Retain` | PV stays — data preserved, must be manually reclaimed |
| `Delete` | PV and underlying storage are deleted |
| `Recycle` | Deprecated — don't use |

> Use `Retain` for databases. Use `Delete` for temporary/ephemeral data.

---

## Storage for Multi-Pod Sharing (ReadWriteMany)

Some apps need multiple pods to share the same storage (e.g., Pimcore assets, shared file uploads):

| Provider | ReadWriteMany Solution |
|---|---|
| AWS | EFS with `aws-efs-csi-driver` |
| GCP | Filestore with `filestore.csi.storage.gke.io` |
| Azure | Azure Files |
| On-premise | NFS server with `nfs-subdir-external-provisioner` |
| Any | Rook-Ceph distributed storage |

```yaml
# PVC for shared file storage
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: shared-assets-pvc
spec:
  storageClassName: efs-sc       # Must support RWX
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 50Gi
```

---

## emptyDir (Temporary Shared Volume)

`emptyDir` is a temporary volume that exists as long as the pod lives. Data is lost when the pod is deleted. Useful for sharing data between containers in the same pod.

```yaml
spec:
  containers:
    - name: app
      image: my-app
      volumeMounts:
        - name: tmp-data
          mountPath: /tmp/cache

    - name: sidecar
      image: my-sidecar
      volumeMounts:
        - name: tmp-data
          mountPath: /shared

  volumes:
    - name: tmp-data
      emptyDir: {}
```

---

## hostPath (Single Node Dev Only)

Mounts a file or directory from the host node. **Never use in production multi-node clusters.**

```yaml
volumes:
  - name: host-data
    hostPath:
      path: /data/myapp
      type: DirectoryOrCreate
```

---

## Volume Expansion

```yaml
# PVC must have allowVolumeExpansion: true on its StorageClass
spec:
  resources:
    requests:
      storage: 50Gi   # Increase from 20Gi to 50Gi
```

```bash
kubectl patch pvc data-pvc -p '{"spec":{"resources":{"requests":{"storage":"50Gi"}}}}'
kubectl get pvc data-pvc  # watch for Resizing → FileSystemResizePending → Bound
```

---

## Storage Commands Cheatsheet

```bash
# View PVs
kubectl get pv
kubectl describe pv my-pv

# View PVCs
kubectl get pvc
kubectl get pvc -A   # all namespaces
kubectl describe pvc my-pvc

# View StorageClasses
kubectl get storageclass

# Check which pod uses which PVC
kubectl get pods -o json | jq '.items[].spec.volumes[].persistentVolumeClaim'

# Delete PVC (data fate depends on reclaimPolicy)
kubectl delete pvc my-pvc

# Force delete stuck PVC
kubectl patch pvc my-pvc -p '{"metadata":{"finalizers":null}}'
```

---

## Common Issues

| Issue | Cause | Fix |
|---|---|---|
| PVC stuck in `Pending` | No matching PV or StorageClass | Check `kubectl describe pvc` for events |
| PV in `Released` state | PVC deleted but PV not reclaimed | Delete and recreate PV or patch claimRef |
| Data lost after pod restart | Used emptyDir instead of PVC | Switch to PVC |
| Multi-pod write errors | Using RWO with multiple pods | Switch to RWX storage |
| Volume not expandable | StorageClass missing `allowVolumeExpansion: true` | Update StorageClass |
