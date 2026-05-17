# Chapter 13: Informers & Watching

## Watch vs Informer

| | Raw Watch | SharedInformer |
|---|---|---|
| Local cache | No | Yes (synced) |
| Handles reconnection | Manual | Automatic |
| Multiple handlers | No | Yes |
| Event deduplication | No | Yes |
| Use case | Simple one-off watch | Controllers, operators |

---

## Raw Watch

```go
watcher, err := clientset.CoreV1().Pods("default").Watch(ctx, metav1.ListOptions{
    LabelSelector: "app=my-app",
})
if err != nil {
    log.Fatal(err)
}
defer watcher.Stop()

for event := range watcher.ResultChan() {
    pod, ok := event.Object.(*corev1.Pod)
    if !ok {
        continue
    }
    fmt.Printf("[%s] Pod: %s Phase: %s\n", event.Type, pod.Name, pod.Status.Phase)
    // event.Type: watch.Added | watch.Modified | watch.Deleted | watch.Error
}
```

---

## SharedInformerFactory (Recommended)

The `SharedInformerFactory` manages informers for all resource types and shares a single watch stream.

```go
import (
    "k8s.io/client-go/informers"
    "k8s.io/client-go/tools/cache"
)

// Create factory with 30-second resync period
factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

// Get a pod informer
podInformer := factory.Core().V1().Pods().Informer()

// Register event handlers
podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
    AddFunc: func(obj interface{}) {
        pod := obj.(*corev1.Pod)
        fmt.Printf("[ADDED] Pod: %s\n", pod.Name)
    },
    UpdateFunc: func(oldObj, newObj interface{}) {
        pod := newObj.(*corev1.Pod)
        fmt.Printf("[UPDATED] Pod: %s Phase: %s\n", pod.Name, pod.Status.Phase)
    },
    DeleteFunc: func(obj interface{}) {
        pod, ok := obj.(*corev1.Pod)
        if !ok {
            // Handle tombstone
            tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
            if !ok {
                return
            }
            pod = tombstone.Obj.(*corev1.Pod)
        }
        fmt.Printf("[DELETED] Pod: %s\n", pod.Name)
    },
})

// Start the factory (begins syncing)
stopCh := make(chan struct{})
defer close(stopCh)

factory.Start(stopCh)

// Wait until caches are synced (initial list complete)
if !cache.WaitForCacheSync(stopCh, podInformer.HasSynced) {
    log.Fatal("Timed out waiting for caches to sync")
}

fmt.Println("Cache synced — watching for events...")
<-stopCh
```

---

## Namespace-Scoped Informer

```go
// Only watch a specific namespace
factory := informers.NewSharedInformerFactoryWithOptions(
    clientset,
    30*time.Second,
    informers.WithNamespace("production"),
)
```

---

## Using a Lister (Read from Cache)

```go
// After the informer syncs, use Lister for fast local reads (no API calls)
podLister := factory.Core().V1().Pods().Lister()

// List all pods from cache
pods, err := podLister.List(labels.Everything())

// List pods in a namespace
pods, err := podLister.Pods("default").List(labels.Everything())

// Get a specific pod
pod, err := podLister.Pods("default").Get("my-pod")
```

---

## Informers for Different Resources

```go
factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

// All common informers
podInformer        := factory.Core().V1().Pods().Informer()
deployInformer     := factory.Apps().V1().Deployments().Informer()
svcInformer        := factory.Core().V1().Services().Informer()
cmInformer         := factory.Core().V1().ConfigMaps().Informer()
secretInformer     := factory.Core().V1().Secrets().Informer()
nsInformer         := factory.Core().V1().Namespaces().Informer()
nodeInformer       := factory.Core().V1().Nodes().Informer()
stsInformer        := factory.Apps().V1().StatefulSets().Informer()
dsInformer         := factory.Apps().V1().DaemonSets().Informer()
jobInformer        := factory.Batch().V1().Jobs().Informer()
cronJobInformer    := factory.Batch().V1().CronJobs().Informer()
ingressInformer    := factory.Networking().V1().Ingresses().Informer()
```

---

## WorkQueue — Controller Pattern

The standard pattern for production controllers: informer enqueues → worker dequeues → reconcile.

```go
import "k8s.io/client-go/util/workqueue"

type Controller struct {
    clientset   *kubernetes.Clientset
    podInformer cache.SharedIndexInformer
    podLister   listersv1.PodNamespaceLister
    queue       workqueue.RateLimitingInterface
}

func NewController(clientset *kubernetes.Clientset, factory informers.SharedInformerFactory) *Controller {
    podInformer := factory.Core().V1().Pods()
    queue       := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

    c := &Controller{
        clientset:   clientset,
        podInformer: podInformer.Informer(),
        podLister:   podInformer.Lister().Pods("default"),
        queue:       queue,
    }

    c.podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc:    func(obj interface{}) { c.enqueue(obj) },
        UpdateFunc: func(_, obj interface{}) { c.enqueue(obj) },
        DeleteFunc: func(obj interface{}) { c.enqueue(obj) },
    })

    return c
}

func (c *Controller) enqueue(obj interface{}) {
    key, err := cache.MetaNamespaceKeyFunc(obj)
    if err != nil {
        return
    }
    c.queue.Add(key)
}

func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
    defer c.queue.ShutDown()

    if !cache.WaitForCacheSync(stopCh, c.podInformer.HasSynced) {
        log.Fatal("Cache sync failed")
    }

    for i := 0; i < workers; i++ {
        go c.runWorker()
    }

    <-stopCh
}

func (c *Controller) runWorker() {
    for c.processNextItem() {}
}

func (c *Controller) processNextItem() bool {
    key, quit := c.queue.Get()
    if quit {
        return false
    }
    defer c.queue.Done(key)

    err := c.reconcile(key.(string))
    if err != nil {
        // Retry up to 5 times
        if c.queue.NumRequeues(key) < 5 {
            c.queue.AddRateLimited(key)
            return true
        }
        c.queue.Forget(key)
        log.Printf("Dropping key %s after max retries: %v", key, err)
        return true
    }

    c.queue.Forget(key)
    return true
}

func (c *Controller) reconcile(key string) error {
    namespace, name, err := cache.SplitMetaNamespaceKey(key)
    if err != nil {
        return err
    }

    pod, err := c.podLister.Get(name)
    if errors.IsNotFound(err) {
        fmt.Printf("Pod %s/%s was deleted\n", namespace, name)
        return nil
    }
    if err != nil {
        return err
    }

    fmt.Printf("Reconciling pod %s/%s phase=%s\n", namespace, name, pod.Status.Phase)
    // Your reconciliation logic here
    return nil
}
```

---

## Watching Events (kubectl get events equivalent)

```go
eventInformer := factory.Core().V1().Events().Informer()

eventInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
    AddFunc: func(obj interface{}) {
        event := obj.(*corev1.Event)
        if event.Type == corev1.EventTypeWarning {
            fmt.Printf("[WARNING] %s/%s — %s: %s\n",
                event.InvolvedObject.Kind,
                event.InvolvedObject.Name,
                event.Reason,
                event.Message,
            )
        }
    },
})
```

---

## Label Selector in Watch

```go
import "k8s.io/apimachinery/pkg/labels"

selector := labels.Set{"app": "nginx"}.String()

watcher, _ := clientset.CoreV1().Pods("default").Watch(ctx, metav1.ListOptions{
    LabelSelector: selector,
})
```

---

## Detecting Restart Storms

```go
podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
    UpdateFunc: func(oldObj, newObj interface{}) {
        pod    := newObj.(*corev1.Pod)
        oldPod := oldObj.(*corev1.Pod)

        for i, cs := range pod.Status.ContainerStatuses {
            if i < len(oldPod.Status.ContainerStatuses) {
                oldCS := oldPod.Status.ContainerStatuses[i]
                if cs.RestartCount > oldCS.RestartCount {
                    fmt.Printf("[RESTART] Pod=%s Container=%s restarts=%d\n",
                        pod.Name, cs.Name, cs.RestartCount)
                }
            }
            if cs.State.Waiting != nil && cs.State.Waiting.Reason == "CrashLoopBackOff" {
                fmt.Printf("[CRASHLOOP] Pod=%s Container=%s\n", pod.Name, cs.Name)
            }
        }
    },
})
```
