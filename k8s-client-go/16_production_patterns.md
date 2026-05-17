# Chapter 16: Production Patterns

## 1. Controller Pattern (Full Example)

The standard Go controller loop: watch → enqueue → reconcile.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    corev1  "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    "k8s.io/apimachinery/pkg/util/wait"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    listersv1 "k8s.io/client-go/listers/core/v1"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/util/workqueue"
)

type PodController struct {
    clientset *kubernetes.Clientset
    lister    listersv1.PodLister
    synced    cache.InformerSynced
    queue     workqueue.RateLimitingInterface
}

func NewPodController(clientset *kubernetes.Clientset, factory informers.SharedInformerFactory) *PodController {
    podInformer := factory.Core().V1().Pods()
    c := &PodController{
        clientset: clientset,
        lister:    podInformer.Lister(),
        synced:    podInformer.Informer().HasSynced,
        queue:     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "pods"),
    }

    podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc:    func(obj interface{}) { c.enqueue(obj) },
        UpdateFunc: func(_, obj interface{}) { c.enqueue(obj) },
        DeleteFunc: func(obj interface{}) { c.enqueue(obj) },
    })

    return c
}

func (c *PodController) enqueue(obj interface{}) {
    key, err := cache.MetaNamespaceKeyFunc(obj)
    if err != nil {
        return
    }
    c.queue.Add(key)
}

func (c *PodController) Run(ctx context.Context, workers int) {
    defer c.queue.ShutDown()
    log.Println("Starting PodController")

    if !cache.WaitForCacheSync(ctx.Done(), c.synced) {
        log.Fatal("Timed out waiting for caches to sync")
    }
    log.Println("Cache synced")

    for i := 0; i < workers; i++ {
        go wait.UntilWithContext(ctx, c.runWorker, time.Second)
    }

    <-ctx.Done()
    log.Println("Shutting down PodController")
}

func (c *PodController) runWorker(ctx context.Context) {
    for c.processNext(ctx) {}
}

func (c *PodController) processNext(ctx context.Context) bool {
    key, quit := c.queue.Get()
    if quit {
        return false
    }
    defer c.queue.Done(key)

    if err := c.reconcile(ctx, key.(string)); err != nil {
        if c.queue.NumRequeues(key) < 5 {
            log.Printf("Retrying %s: %v", key, err)
            c.queue.AddRateLimited(key)
            return true
        }
        log.Printf("Dropping %s after max retries: %v", key, err)
    }

    c.queue.Forget(key)
    return true
}

func (c *PodController) reconcile(ctx context.Context, key string) error {
    namespace, name, _ := cache.SplitMetaNamespaceKey(key)

    pod, err := c.lister.Pods(namespace).Get(name)
    if errors.IsNotFound(err) {
        log.Printf("Pod %s/%s deleted", namespace, name)
        return nil
    }
    if err != nil {
        return err
    }

    log.Printf("Reconciling pod %s/%s phase=%s", namespace, name, pod.Status.Phase)
    // Your business logic here
    return nil
}

func main() {
    config, _ := rest.InClusterConfig()
    clientset, _ := kubernetes.NewForConfig(config)

    factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)
    ctrl := NewPodController(clientset, factory)

    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    factory.Start(ctx.Done())
    ctrl.Run(ctx, 2) // 2 worker goroutines
}
```

---

## 2. Leader Election

Only one instance of your controller runs at a time across multiple replicas:

```go
import (
    "k8s.io/client-go/tools/leaderelection"
    "k8s.io/client-go/tools/leaderelection/resourcelock"
)

func runWithLeaderElection(ctx context.Context, clientset *kubernetes.Clientset, id string, run func(ctx context.Context)) {
    lock := &resourcelock.LeaseLock{
        LeaseMeta: metav1.ObjectMeta{
            Name:      "my-controller-lock",
            Namespace: "default",
        },
        Client: clientset.CoordinationV1(),
        LockConfig: resourcelock.ResourceLockConfig{
            Identity: id,
        },
    }

    leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
        Lock:            lock,
        ReleaseOnCancel: true,
        LeaseDuration:   15 * time.Second,
        RenewDeadline:   10 * time.Second,
        RetryPeriod:     2 * time.Second,
        Callbacks: leaderelection.LeaderCallbacks{
            OnStartedLeading: func(ctx context.Context) {
                fmt.Printf("I am now the leader: %s\n", id)
                run(ctx)
            },
            OnStoppedLeading: func() {
                fmt.Printf("Lost leadership: %s\n", id)
            },
            OnNewLeader: func(identity string) {
                if identity == id {
                    return
                }
                fmt.Printf("New leader elected: %s\n", identity)
            },
        },
    })
}

// Usage — each pod gets a unique ID
podName := os.Getenv("POD_NAME")
runWithLeaderElection(ctx, clientset, podName, func(ctx context.Context) {
    ctrl.Run(ctx, 2)
})
```

---

## 3. Deployment Pipeline

```go
type Deployer struct {
    clientset *kubernetes.Clientset
}

func (d *Deployer) Deploy(ctx context.Context, namespace, deployment, container, image string) error {
    fmt.Printf("Deploying %s to %s/%s...\n", image, namespace, deployment)

    if err := d.updateImage(ctx, namespace, deployment, container, image); err != nil {
        return fmt.Errorf("image update failed: %w", err)
    }

    rolloutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
    defer cancel()

    if err := d.waitForRollout(rolloutCtx, namespace, deployment); err != nil {
        fmt.Println("Rollout failed — rolling back")
        d.rollback(ctx, namespace, deployment)
        return fmt.Errorf("rollout failed: %w", err)
    }

    fmt.Printf("Deploy of %s succeeded\n", image)
    return nil
}

func (d *Deployer) updateImage(ctx context.Context, ns, name, container, image string) error {
    return retry.RetryOnConflict(retry.DefaultRetry, func() error {
        dep, err := d.clientset.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }
        for i, c := range dep.Spec.Template.Spec.Containers {
            if c.Name == container {
                dep.Spec.Template.Spec.Containers[i].Image = image
                break
            }
        }
        _, err = d.clientset.AppsV1().Deployments(ns).Update(ctx, dep, metav1.UpdateOptions{})
        return err
    })
}

func (d *Deployer) waitForRollout(ctx context.Context, ns, name string) error {
    for {
        dep, err := d.clientset.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }
        desired := *dep.Spec.Replicas
        if dep.Status.UpdatedReplicas == desired &&
            dep.Status.ReadyReplicas == desired &&
            dep.Status.AvailableReplicas == desired {
            return nil
        }
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(3 * time.Second):
        }
    }
}

func (d *Deployer) rollback(ctx context.Context, ns, name string) {
    retry.RetryOnConflict(retry.DefaultRetry, func() error {
        dep, err := d.clientset.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }
        if dep.Spec.Template.Annotations == nil {
            dep.Spec.Template.Annotations = map[string]string{}
        }
        dep.Spec.Template.Annotations["rollback"] = time.Now().Format(time.RFC3339)
        _, err = d.clientset.AppsV1().Deployments(ns).Update(ctx, dep, metav1.UpdateOptions{})
        return err
    })
}
```

---

## 4. Cluster Health Dashboard

```go
type ClusterHealth struct {
    Nodes       NodeHealth       `json:"nodes"`
    Pods        PodHealth        `json:"pods"`
    Deployments DeploymentHealth `json:"deployments"`
    Timestamp   string           `json:"timestamp"`
}

func clusterHealthReport(ctx context.Context, clientset *kubernetes.Clientset) (*ClusterHealth, error) {
    health := &ClusterHealth{Timestamp: time.Now().Format(time.RFC3339)}

    // Nodes
    nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    health.Nodes.Total = len(nodes.Items)
    for _, node := range nodes.Items {
        for _, cond := range node.Status.Conditions {
            if cond.Type == corev1.NodeReady && cond.Status == corev1.ConditionTrue {
                health.Nodes.Ready++
            }
        }
    }

    // Pods (all namespaces)
    pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    health.Pods.Total = len(pods.Items)
    for _, pod := range pods.Items {
        switch pod.Status.Phase {
        case corev1.PodRunning:
            health.Pods.Running++
        case corev1.PodFailed:
            health.Pods.Failed++
        case corev1.PodPending:
            health.Pods.Pending++
        }
    }

    // Deployments
    deps, err := clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    health.Deployments.Total = len(deps.Items)
    for _, dep := range deps.Items {
        if dep.Status.ReadyReplicas == *dep.Spec.Replicas {
            health.Deployments.Healthy++
        }
    }

    return health, nil
}
```

---

## 5. Graceful Shutdown

```go
import (
    "os"
    "os/signal"
    "syscall"
)

func main() {
    config, _ := rest.InClusterConfig()
    clientset, _ := kubernetes.NewForConfig(config)
    factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)
    ctrl := NewPodController(clientset, factory)

    // Context cancelled on SIGINT / SIGTERM
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    factory.Start(ctx.Done())
    ctrl.Run(ctx, 2)

    log.Println("Controller stopped cleanly")
}
```

---

## Production Checklist

| Concern | Solution |
|---|---|
| Authentication | `rest.InClusterConfig()` in production pods |
| Least privilege | ServiceAccount with minimal RBAC |
| Conflict retry | Always use `retry.RetryOnConflict` for updates |
| Context timeouts | Every API call gets a context with deadline |
| Rate limiting | Set `config.QPS` and `config.Burst` appropriately |
| Leader election | Use `leaderelection` for multi-replica controllers |
| Informer cache | Use informers + listers — avoid direct API calls in hot paths |
| Graceful shutdown | `signal.NotifyContext` + `factory.Shutdown()` |
| Error logging | Log with structured fields (namespace, name, reason) |
| Health probes | Expose `/healthz` and `/readyz` endpoints |
| Metrics | Use `k8s.io/component-base/metrics` for controller metrics |
