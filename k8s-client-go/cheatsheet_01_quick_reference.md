# k8s-client-go — Master Cheatsheet

## Install

```bash
go get k8s.io/client-go@v0.29.0
go get k8s.io/api@v0.29.0
go get k8s.io/apimachinery@v0.29.0
go mod tidy
```

---

## Connect

```go
import (
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

// From kubeconfig (dev)
config, _ := clientcmd.BuildConfigFromFlags("", filepath.Join(os.UserHomeDir(), ".kube", "config"))
clientset, _ := kubernetes.NewForConfig(config)

// In-cluster (Pod inside K8s)
config, _ = rest.InClusterConfig()
clientset, _ = kubernetes.NewForConfig(config)

// Auto-detect (try in-cluster, fall back to kubeconfig)
config, err := rest.InClusterConfig()
if err != nil {
    config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
}
clientset, _ = kubernetes.NewForConfig(config)

// Specific context
overrides := &clientcmd.ConfigOverrides{CurrentContext: "my-cluster"}
config, _ = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
    clientcmd.NewDefaultClientConfigLoadingRules(), overrides,
).ClientConfig()

// Tune QPS
config.QPS = 100
config.Burst = 200
```

---

## Pods

```go
// List
clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})                           // all namespaces
clientset.CoreV1().Pods("ns").List(ctx, metav1.ListOptions{LabelSelector: "app=nginx"})
clientset.CoreV1().Pods("ns").List(ctx, metav1.ListOptions{FieldSelector: "status.phase=Running"})

// Get
clientset.CoreV1().Pods("default").Get(ctx, "my-pod", metav1.GetOptions{})

// Create
clientset.CoreV1().Pods("default").Create(ctx, &corev1.Pod{...}, metav1.CreateOptions{})

// Delete
clientset.CoreV1().Pods("default").Delete(ctx, "my-pod", metav1.DeleteOptions{})
clientset.CoreV1().Pods("default").Delete(ctx, "my-pod", metav1.DeleteOptions{GracePeriodSeconds: int64Ptr(0)})

// Logs
req := clientset.CoreV1().Pods("default").GetLogs("my-pod", &corev1.PodLogOptions{})
req := clientset.CoreV1().Pods("default").GetLogs("my-pod", &corev1.PodLogOptions{TailLines: int64Ptr(100)})
req := clientset.CoreV1().Pods("default").GetLogs("my-pod", &corev1.PodLogOptions{Previous: true})
req := clientset.CoreV1().Pods("default").GetLogs("my-pod", &corev1.PodLogOptions{Container: "sidecar"})
stream, _ := req.Stream(ctx)
defer stream.Close()
io.Copy(os.Stdout, stream)

// Watch
watcher, _ := clientset.CoreV1().Pods("default").Watch(ctx, metav1.ListOptions{})
for event := range watcher.ResultChan() {
    pod := event.Object.(*corev1.Pod)
    fmt.Printf("[%s] %s\n", event.Type, pod.Name)
}

// Read pod fields
pod.Name
pod.Namespace
pod.Status.Phase           // corev1.PodRunning | PodPending | PodFailed | PodSucceeded
pod.Status.PodIP
pod.Spec.NodeName
pod.Status.ContainerStatuses[0].RestartCount
pod.Status.ContainerStatuses[0].State.Waiting.Reason   // CrashLoopBackOff, ImagePullBackOff
```

---

## Deployments

```go
// List / Get
clientset.AppsV1().Deployments("ns").List(ctx, metav1.ListOptions{})
clientset.AppsV1().Deployments("ns").Get(ctx, "name", metav1.GetOptions{})

// Create
clientset.AppsV1().Deployments("ns").Create(ctx, &appsv1.Deployment{...}, metav1.CreateOptions{})

// Scale (with retry)
retry.RetryOnConflict(retry.DefaultRetry, func() error {
    d, _ := clientset.AppsV1().Deployments("ns").Get(ctx, "name", metav1.GetOptions{})
    d.Spec.Replicas = int32Ptr(5)
    _, err := clientset.AppsV1().Deployments("ns").Update(ctx, d, metav1.UpdateOptions{})
    return err
})

// Update image
retry.RetryOnConflict(retry.DefaultRetry, func() error {
    d, _ := clientset.AppsV1().Deployments("ns").Get(ctx, "name", metav1.GetOptions{})
    d.Spec.Template.Spec.Containers[0].Image = "myapp:2.0"
    _, err := clientset.AppsV1().Deployments("ns").Update(ctx, d, metav1.UpdateOptions{})
    return err
})

// Rollout restart
retry.RetryOnConflict(retry.DefaultRetry, func() error {
    d, _ := clientset.AppsV1().Deployments("ns").Get(ctx, "name", metav1.GetOptions{})
    if d.Spec.Template.Annotations == nil { d.Spec.Template.Annotations = map[string]string{} }
    d.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
    _, err := clientset.AppsV1().Deployments("ns").Update(ctx, d, metav1.UpdateOptions{})
    return err
})

// Delete
clientset.AppsV1().Deployments("ns").Delete(ctx, "name", metav1.DeleteOptions{PropagationPolicy: &foreground})

// Read status
dep.Spec.Replicas
dep.Status.ReadyReplicas
dep.Status.UpdatedReplicas
dep.Status.AvailableReplicas
```

---

## Services

```go
clientset.CoreV1().Services("ns").List(ctx, metav1.ListOptions{})
clientset.CoreV1().Services("ns").Get(ctx, "name", metav1.GetOptions{})
clientset.CoreV1().Services("ns").Create(ctx, &corev1.Service{
    Spec: corev1.ServiceSpec{
        Type:     corev1.ServiceTypeClusterIP,    // NodePort | LoadBalancer | ExternalName
        Selector: map[string]string{"app": "x"},
        Ports:    []corev1.ServicePort{{Port: 80, TargetPort: intstr.FromInt(8080)}},
    },
}, metav1.CreateOptions{})
clientset.CoreV1().Services("ns").Delete(ctx, "name", metav1.DeleteOptions{})

svc.Spec.ClusterIP
svc.Spec.Type
svc.Status.LoadBalancer.Ingress[0].IP
```

---

## ConfigMaps

```go
clientset.CoreV1().ConfigMaps("ns").List(ctx, metav1.ListOptions{})
clientset.CoreV1().ConfigMaps("ns").Get(ctx, "name", metav1.GetOptions{})
cm.Data["KEY"]

clientset.CoreV1().ConfigMaps("ns").Create(ctx, &corev1.ConfigMap{
    Data: map[string]string{"KEY": "value"},
}, metav1.CreateOptions{})

retry.RetryOnConflict(retry.DefaultRetry, func() error {
    cm, _ := clientset.CoreV1().ConfigMaps("ns").Get(ctx, "name", metav1.GetOptions{})
    cm.Data["KEY"] = "new"
    _, err := clientset.CoreV1().ConfigMaps("ns").Update(ctx, cm, metav1.UpdateOptions{})
    return err
})

clientset.CoreV1().ConfigMaps("ns").Delete(ctx, "name", metav1.DeleteOptions{})
```

---

## Secrets

```go
clientset.CoreV1().Secrets("ns").Get(ctx, "name", metav1.GetOptions{})
string(secret.Data["password"])   // already decoded — no base64 needed

clientset.CoreV1().Secrets("ns").Create(ctx, &corev1.Secret{
    Type:       corev1.SecretTypeOpaque,
    StringData: map[string]string{"password": "secret"},  // auto-encoded
}, metav1.CreateOptions{})

clientset.CoreV1().Secrets("ns").Delete(ctx, "name", metav1.DeleteOptions{})
```

---

## Namespaces & RBAC

```go
clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
clientset.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "ns"}}, metav1.CreateOptions{})
clientset.CoreV1().Namespaces().Delete(ctx, "ns", metav1.DeleteOptions{})

clientset.RbacV1().Roles("ns").Create(ctx, &rbacv1.Role{Rules: []rbacv1.PolicyRule{{APIGroups: []string{""}, Resources: []string{"pods"}, Verbs: []string{"get","list"}}}}, metav1.CreateOptions{})
clientset.RbacV1().ClusterRoles().Create(ctx, &rbacv1.ClusterRole{...}, metav1.CreateOptions{})
clientset.RbacV1().RoleBindings("ns").Create(ctx, &rbacv1.RoleBinding{Subjects: [...], RoleRef: ...}, metav1.CreateOptions{})
clientset.RbacV1().ClusterRoleBindings().Create(ctx, &rbacv1.ClusterRoleBinding{...}, metav1.CreateOptions{})
clientset.CoreV1().ServiceAccounts("ns").Create(ctx, &corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: "sa"}}, metav1.CreateOptions{})
```

---

## PVC / Storage

```go
clientset.CoreV1().PersistentVolumeClaims("ns").List(ctx, metav1.ListOptions{})
clientset.CoreV1().PersistentVolumeClaims("ns").Create(ctx, &corev1.PersistentVolumeClaim{
    Spec: corev1.PersistentVolumeClaimSpec{
        AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
        StorageClassName: &storageClass,
        Resources:        corev1.ResourceRequirements{Requests: corev1.ResourceList{corev1.ResourceStorage: resource.MustParse("10Gi")}},
    },
}, metav1.CreateOptions{})
pvc.Status.Phase   // corev1.ClaimBound | ClaimPending | ClaimLost
clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
clientset.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
```

---

## Jobs & CronJobs

```go
clientset.BatchV1().Jobs("ns").List(ctx, metav1.ListOptions{})
clientset.BatchV1().Jobs("ns").Create(ctx, &batchv1.Job{...}, metav1.CreateOptions{})
job.Status.Succeeded
job.Status.Failed
job.Status.Active
clientset.BatchV1().Jobs("ns").Delete(ctx, "name", metav1.DeleteOptions{PropagationPolicy: &foreground})

clientset.BatchV1().CronJobs("ns").List(ctx, metav1.ListOptions{})
clientset.BatchV1().CronJobs("ns").Create(ctx, &batchv1.CronJob{Spec: batchv1.CronJobSpec{Schedule: "0 * * * *", ...}}, metav1.CreateOptions{})
cj.Spec.Schedule
cj.Status.LastScheduleTime
```

---

## Informers

```go
factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)
podInformer := factory.Core().V1().Pods().Informer()

podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
    AddFunc:    func(obj interface{}) { pod := obj.(*corev1.Pod); fmt.Println("ADDED:", pod.Name) },
    UpdateFunc: func(_, obj interface{}) { pod := obj.(*corev1.Pod); fmt.Println("UPDATED:", pod.Name) },
    DeleteFunc: func(obj interface{}) { /* handle tombstone */ },
})

factory.Start(stopCh)
cache.WaitForCacheSync(stopCh, podInformer.HasSynced)

// Lister (no API call — reads from cache)
lister := factory.Core().V1().Pods().Lister()
pods, _ := lister.Pods("default").List(labels.Everything())
pod, _ := lister.Pods("default").Get("my-pod")
```

---

## Dynamic Client (CRDs)

```go
dynClient, _ := dynamic.NewForConfig(config)

gvr := schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "certificates"}

dynClient.Resource(gvr).Namespace("default").List(ctx, metav1.ListOptions{})
dynClient.Resource(gvr).Namespace("default").Get(ctx, "name", metav1.GetOptions{})
dynClient.Resource(gvr).Namespace("default").Create(ctx, &unstructured.Unstructured{Object: map[string]interface{}{...}}, metav1.CreateOptions{})
dynClient.Resource(gvr).Namespace("default").Delete(ctx, "name", metav1.DeleteOptions{})

// Read fields
name, _, _ := unstructured.NestedString(obj.Object, "spec", "secretName")
unstructured.SetNestedField(obj.Object, "value", "spec", "field")
```

---

## Error Handling

```go
import "k8s.io/apimachinery/pkg/api/errors"

errors.IsNotFound(err)       // 404
errors.IsAlreadyExists(err)  // 409
errors.IsConflict(err)       // 409 (stale resourceVersion)
errors.IsForbidden(err)      // 403
errors.IsUnauthorized(err)   // 401
errors.IsInvalid(err)        // 422
errors.IsTooManyRequests(err)// 429

// Retry on conflict (always use for updates)
retry.RetryOnConflict(retry.DefaultRetry, func() error { ... })
```

---

## Helper Functions

```go
func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
func boolPtr(b bool) *bool    { return &b }
func stringPtr(s string) *string { return &s }

// Propagation
foreground := metav1.DeletePropagationForeground
background := metav1.DeletePropagationBackground
```

---

## Key Imports

```go
import (
    "context"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"

    appsv1  "k8s.io/api/apps/v1"
    batchv1 "k8s.io/api/batch/v1"
    corev1  "k8s.io/api/core/v1"
    netv1   "k8s.io/api/networking/v1"
    rbacv1  "k8s.io/api/rbac/v1"

    "k8s.io/apimachinery/pkg/api/errors"
    "k8s.io/apimachinery/pkg/api/resource"
    metav1  "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
    "k8s.io/apimachinery/pkg/runtime/schema"

    "k8s.io/client-go/dynamic"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/retry"
    "k8s.io/client-go/util/workqueue"
)
```
