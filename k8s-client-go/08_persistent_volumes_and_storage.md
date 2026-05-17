# Chapter 8: Persistent Volumes & Storage

## Storage Interfaces

```go
clientset.CoreV1().PersistentVolumes()                      // cluster-scoped
clientset.CoreV1().PersistentVolumeClaims(namespace)        // namespaced
clientset.StorageV1().StorageClasses()                      // cluster-scoped
```

---

## Working with StorageClasses

### List StorageClasses

```go
scs, err := clientset.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
for _, sc := range scs.Items {
    isDefault := sc.Annotations["storageclass.kubernetes.io/is-default-class"] == "true"
    defaultStr := ""
    if isDefault {
        defaultStr = "(default)"
    }
    fmt.Printf("%-30s %-40s %-10s %s\n",
        sc.Name, sc.Provisioner, *sc.ReclaimPolicy, defaultStr)
}
```

### Create a StorageClass

```go
import (
    storagev1 "k8s.io/api/storage/v1"
    corev1    "k8s.io/api/core/v1"
)

reclaimPolicy    := corev1.PersistentVolumeReclaimRetain
bindingMode      := storagev1.VolumeBindingWaitForFirstConsumer

sc := &storagev1.StorageClass{
    ObjectMeta: metav1.ObjectMeta{
        Name: "fast-ssd",
        Annotations: map[string]string{
            "storageclass.kubernetes.io/is-default-class": "false",
        },
    },
    Provisioner:          "kubernetes.io/gce-pd",
    ReclaimPolicy:        &reclaimPolicy,
    VolumeBindingMode:    &bindingMode,
    Parameters: map[string]string{
        "type":             "pd-ssd",
        "replication-type": "regional-pd",
    },
}

clientset.StorageV1().StorageClasses().Create(ctx, sc, metav1.CreateOptions{})
```

---

## Working with PersistentVolumeClaims

### List PVCs

```go
pvcs, err := clientset.CoreV1().PersistentVolumeClaims("default").List(ctx, metav1.ListOptions{})
for _, pvc := range pvcs.Items {
    capacity := pvc.Status.Capacity.Storage()
    fmt.Printf("%-30s %-10s %-10s %-20s %s\n",
        pvc.Name,
        pvc.Status.Phase,
        capacity.String(),
        pvc.Spec.StorageClassName,
        strings.Join(pvc.Spec.AccessModes, ","),
    )
}
```

### Create a PVC

```go
import (
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
)

storageClass := "standard"

pvc := &corev1.PersistentVolumeClaim{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "my-data-pvc",
        Namespace: "default",
        Labels:    map[string]string{"app": "my-app"},
    },
    Spec: corev1.PersistentVolumeClaimSpec{
        AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
        StorageClassName: &storageClass,
        Resources: corev1.ResourceRequirements{
            Requests: corev1.ResourceList{
                corev1.ResourceStorage: resource.MustParse("10Gi"),
            },
        },
    },
}

created, err := clientset.CoreV1().PersistentVolumeClaims("default").Create(ctx, pvc, metav1.CreateOptions{})
fmt.Printf("PVC created: %s\n", created.Name)
```

### Wait for PVC to be Bound

```go
func waitForPVCBound(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) error {
    for {
        pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }

        if pvc.Status.Phase == corev1.ClaimBound {
            fmt.Printf("PVC %s bound to PV %s\n", name, pvc.Spec.VolumeName)
            return nil
        }

        fmt.Printf("PVC phase: %s — waiting...\n", pvc.Status.Phase)
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(3 * time.Second):
        }
    }
}
```

### PVC with ReadWriteMany (Shared Storage)

```go
storageClass := "nfs-client"

pvc := &corev1.PersistentVolumeClaim{
    ObjectMeta: metav1.ObjectMeta{Name: "shared-data", Namespace: "default"},
    Spec: corev1.PersistentVolumeClaimSpec{
        AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
        StorageClassName: &storageClass,
        Resources: corev1.ResourceRequirements{
            Requests: corev1.ResourceList{
                corev1.ResourceStorage: resource.MustParse("50Gi"),
            },
        },
    },
}
```

---

## Working with PersistentVolumes

### List PVs

```go
pvs, err := clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
for _, pv := range pvs.Items {
    capacity := pv.Spec.Capacity.Storage()
    claim := ""
    if pv.Spec.ClaimRef != nil {
        claim = pv.Spec.ClaimRef.Namespace + "/" + pv.Spec.ClaimRef.Name
    }
    fmt.Printf("%-30s %-8s %-10s %-10s %s\n",
        pv.Name, capacity.String(), pv.Status.Phase, pv.Spec.PersistentVolumeReclaimPolicy, claim)
}
```

### Create a PV (NFS Static Provisioning)

```go
pv := &corev1.PersistentVolume{
    ObjectMeta: metav1.ObjectMeta{
        Name:   "my-nfs-pv",
        Labels: map[string]string{"type": "nfs"},
    },
    Spec: corev1.PersistentVolumeSpec{
        Capacity: corev1.ResourceList{
            corev1.ResourceStorage: resource.MustParse("100Gi"),
        },
        AccessModes:                   []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
        PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
        StorageClassName:              "nfs",
        PersistentVolumeSource: corev1.PersistentVolumeSource{
            NFS: &corev1.NFSVolumeSource{
                Server: "10.0.0.5",
                Path:   "/exported/data",
            },
        },
    },
}

clientset.CoreV1().PersistentVolumes().Create(ctx, pv, metav1.CreateOptions{})
```

---

## Mount PVC in a Deployment

```go
deployment := &appsv1.Deployment{
    ObjectMeta: metav1.ObjectMeta{Name: "app-with-storage", Namespace: "default"},
    Spec: appsv1.DeploymentSpec{
        Replicas: int32Ptr(1),
        Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "storage-app"}},
        Template: corev1.PodTemplateSpec{
            ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "storage-app"}},
            Spec: corev1.PodSpec{
                Containers: []corev1.Container{
                    {
                        Name:  "app",
                        Image: "myapp:1.0",
                        VolumeMounts: []corev1.VolumeMount{
                            {Name: "data-volume", MountPath: "/var/data"},
                        },
                    },
                },
                Volumes: []corev1.Volume{
                    {
                        Name: "data-volume",
                        VolumeSource: corev1.VolumeSource{
                            PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
                                ClaimName: "my-data-pvc",
                            },
                        },
                    },
                },
            },
        },
    },
}
```

---

## Delete PVC

```go
// PVC must not be in use by any Pod
err := clientset.CoreV1().PersistentVolumeClaims("default").Delete(ctx, "my-data-pvc", metav1.DeleteOptions{})

// Delete PV manually if reclaimPolicy is Retain
err = clientset.CoreV1().PersistentVolumes().Delete(ctx, "my-nfs-pv", metav1.DeleteOptions{})
```

---

## Storage Capacity Summary

```go
func printStorageSummary(ctx context.Context, clientset *kubernetes.Clientset) {
    pvcs, _ := clientset.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})

    var totalGi int64
    for _, pvc := range pvcs.Items {
        req := pvc.Spec.Resources.Requests[corev1.ResourceStorage]
        totalGi += req.Value() / (1024 * 1024 * 1024)
        fmt.Printf("%-30s %-15s %-10s %s\n",
            pvc.Name, pvc.Namespace, pvc.Status.Phase, req.String())
    }
    fmt.Printf("\nTotal PVC storage: %dGi across %d PVCs\n", totalGi, len(pvcs.Items))
}
```
