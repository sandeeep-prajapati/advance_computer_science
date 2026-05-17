# Chapter 11: StatefulSets & DaemonSets

## Interfaces

```go
clientset.AppsV1().StatefulSets(namespace)
clientset.AppsV1().DaemonSets(namespace)
```

---

## Working with StatefulSets

### List StatefulSets

```go
stsList, err := clientset.AppsV1().StatefulSets("default").List(ctx, metav1.ListOptions{})
for _, sts := range stsList.Items {
    fmt.Printf("%-30s  %d/%d ready\n", sts.Name, sts.Status.ReadyReplicas, *sts.Spec.Replicas)
}
```

### Create a StatefulSet (MySQL)

```go
import (
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
)

storageClass := "standard"

sts := &appsv1.StatefulSet{
    ObjectMeta: metav1.ObjectMeta{Name: "mysql", Namespace: "default"},
    Spec: appsv1.StatefulSetSpec{
        ServiceName: "mysql-headless",
        Replicas:    int32Ptr(3),
        Selector:    &metav1.LabelSelector{MatchLabels: map[string]string{"app": "mysql"}},
        Template: corev1.PodTemplateSpec{
            ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "mysql"}},
            Spec: corev1.PodSpec{
                Containers: []corev1.Container{
                    {
                        Name:  "mysql",
                        Image: "mysql:8.0",
                        Env: []corev1.EnvVar{
                            {
                                Name: "MYSQL_ROOT_PASSWORD",
                                ValueFrom: &corev1.EnvVarSource{
                                    SecretKeyRef: &corev1.SecretKeySelector{
                                        LocalObjectReference: corev1.LocalObjectReference{Name: "mysql-secret"},
                                        Key: "root-password",
                                    },
                                },
                            },
                        },
                        Ports: []corev1.ContainerPort{{ContainerPort: 3306, Name: "mysql"}},
                        VolumeMounts: []corev1.VolumeMount{
                            {Name: "data", MountPath: "/var/lib/mysql"},
                        },
                        Resources: corev1.ResourceRequirements{
                            Requests: corev1.ResourceList{
                                corev1.ResourceCPU:    resource.MustParse("250m"),
                                corev1.ResourceMemory: resource.MustParse("512Mi"),
                            },
                        },
                        ReadinessProbe: &corev1.Probe{
                            ProbeHandler: corev1.ProbeHandler{
                                Exec: &corev1.ExecAction{
                                    Command: []string{"bash", "-c", "mysqladmin ping -h 127.0.0.1"},
                                },
                            },
                            InitialDelaySeconds: 20,
                            PeriodSeconds:       10,
                        },
                    },
                },
            },
        },
        VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
            {
                ObjectMeta: metav1.ObjectMeta{Name: "data"},
                Spec: corev1.PersistentVolumeClaimSpec{
                    AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
                    StorageClassName: &storageClass,
                    Resources: corev1.ResourceRequirements{
                        Requests: corev1.ResourceList{
                            corev1.ResourceStorage: resource.MustParse("20Gi"),
                        },
                    },
                },
            },
        },
    },
}

clientset.AppsV1().StatefulSets("default").Create(ctx, sts, metav1.CreateOptions{})
// Pods: mysql-0, mysql-1, mysql-2
// DNS:  mysql-0.mysql-headless.default.svc.cluster.local
```

### Scale a StatefulSet

```go
err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
    sts, err := clientset.AppsV1().StatefulSets("default").Get(ctx, "mysql", metav1.GetOptions{})
    if err != nil {
        return err
    }
    sts.Spec.Replicas = int32Ptr(5)
    _, err = clientset.AppsV1().StatefulSets("default").Update(ctx, sts, metav1.UpdateOptions{})
    return err
})
```

### StatefulSet Pod DNS

```go
func statefulSetPodDNS(podIndex int, stsName, serviceName, namespace string) string {
    return fmt.Sprintf("%s-%d.%s.%s.svc.cluster.local", stsName, podIndex, serviceName, namespace)
}

sts, _ := clientset.AppsV1().StatefulSets("default").Get(ctx, "mysql", metav1.GetOptions{})
for i := int32(0); i < *sts.Spec.Replicas; i++ {
    dns := statefulSetPodDNS(int(i), "mysql", sts.Spec.ServiceName, "default")
    fmt.Printf("Replica %d: %s\n", i, dns)
}
```

### Update Strategy (Partition / Canary)

```go
partition := int32(2) // only update pods with index >= 2

err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
    sts, err := clientset.AppsV1().StatefulSets("default").Get(ctx, "mysql", metav1.GetOptions{})
    if err != nil {
        return err
    }
    sts.Spec.UpdateStrategy = appsv1.StatefulSetUpdateStrategy{
        Type: appsv1.RollingUpdateStatefulSetStrategyType,
        RollingUpdate: &appsv1.RollingUpdateStatefulSetStrategy{
            Partition: &partition,
        },
    }
    sts.Spec.Template.Spec.Containers[0].Image = "mysql:8.1"
    _, err = clientset.AppsV1().StatefulSets("default").Update(ctx, sts, metav1.UpdateOptions{})
    return err
})
```

---

## Working with DaemonSets

### List DaemonSets

```go
dsList, err := clientset.AppsV1().DaemonSets("kube-system").List(ctx, metav1.ListOptions{})
for _, ds := range dsList.Items {
    fmt.Printf("%-35s desired=%-3d ready=%-3d available=%d\n",
        ds.Name,
        ds.Status.DesiredNumberScheduled,
        ds.Status.NumberReady,
        ds.Status.NumberAvailable,
    )
}
```

### Create a DaemonSet (Fluentd Log Collector)

```go
ds := &appsv1.DaemonSet{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "fluentd",
        Namespace: "kube-system",
        Labels:    map[string]string{"app": "fluentd"},
    },
    Spec: appsv1.DaemonSetSpec{
        Selector: &metav1.LabelSelector{
            MatchLabels: map[string]string{"app": "fluentd"},
        },
        UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
            Type: appsv1.RollingUpdateDaemonSetStrategyType,
            RollingUpdate: &appsv1.RollingUpdateDaemonSet{
                MaxUnavailable: &intstr.IntOrString{Type: intstr.Int, IntVal: 1},
            },
        },
        Template: corev1.PodTemplateSpec{
            ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "fluentd"}},
            Spec: corev1.PodSpec{
                Tolerations: []corev1.Toleration{
                    {
                        Key:    "node-role.kubernetes.io/control-plane",
                        Effect: corev1.TaintEffectNoSchedule,
                    },
                },
                Containers: []corev1.Container{
                    {
                        Name:  "fluentd",
                        Image: "fluent/fluentd-kubernetes-daemonset:v1.16",
                        Resources: corev1.ResourceRequirements{
                            Limits:   corev1.ResourceList{corev1.ResourceMemory: resource.MustParse("200Mi")},
                            Requests: corev1.ResourceList{
                                corev1.ResourceCPU:    resource.MustParse("100m"),
                                corev1.ResourceMemory: resource.MustParse("200Mi"),
                            },
                        },
                        VolumeMounts: []corev1.VolumeMount{
                            {Name: "varlog", MountPath: "/var/log"},
                            {Name: "varlibdockercontainers", MountPath: "/var/lib/docker/containers", ReadOnly: true},
                        },
                    },
                },
                Volumes: []corev1.Volume{
                    {Name: "varlog", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/log"}}},
                    {Name: "varlibdockercontainers", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/lib/docker/containers"}}},
                },
            },
        },
    },
}

clientset.AppsV1().DaemonSets("kube-system").Create(ctx, ds, metav1.CreateOptions{})
```

### DaemonSet on Specific Nodes Only

```go
ds.Spec.Template.Spec.NodeSelector = map[string]string{
    "gpu": "true",   // only nodes labelled gpu=true
}
```

### DaemonSet Health Check

```go
func isDaemonSetHealthy(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) bool {
    ds, err := clientset.AppsV1().DaemonSets(namespace).Get(ctx, name, metav1.GetOptions{})
    if err != nil {
        return false
    }
    return ds.Status.NumberReady == ds.Status.DesiredNumberScheduled && ds.Status.DesiredNumberScheduled > 0
}
```

---

## Delete StatefulSet and DaemonSet

```go
propagation := metav1.DeletePropagationForeground

// Delete StatefulSet (PVCs are NOT deleted automatically)
clientset.AppsV1().StatefulSets("default").Delete(ctx, "mysql", metav1.DeleteOptions{
    PropagationPolicy: &propagation,
})

// Manually delete PVCs
pvcs, _ := clientset.CoreV1().PersistentVolumeClaims("default").List(ctx, metav1.ListOptions{
    LabelSelector: "app=mysql",
})
for _, pvc := range pvcs.Items {
    clientset.CoreV1().PersistentVolumeClaims("default").Delete(ctx, pvc.Name, metav1.DeleteOptions{})
}

// Delete DaemonSet
clientset.AppsV1().DaemonSets("kube-system").Delete(ctx, "fluentd", metav1.DeleteOptions{})
```
