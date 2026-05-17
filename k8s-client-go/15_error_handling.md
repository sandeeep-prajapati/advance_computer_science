# Chapter 15: Error Handling & Retry Logic

## Kubernetes Error Package

```go
import "k8s.io/apimachinery/pkg/api/errors"
```

| Function | Checks For |
|---|---|
| `errors.IsNotFound(err)` | 404 — resource doesn't exist |
| `errors.IsAlreadyExists(err)` | 409 — resource already exists |
| `errors.IsConflict(err)` | 409 — optimistic lock (stale resourceVersion) |
| `errors.IsForbidden(err)` | 403 — no RBAC permission |
| `errors.IsUnauthorized(err)` | 401 — bad/expired credentials |
| `errors.IsInvalid(err)` | 422 — spec validation failed |
| `errors.IsTooManyRequests(err)` | 429 — rate limited |
| `errors.IsServerTimeout(err)` | 503/504 — server timeout |
| `errors.IsTimeout(err)` | Client-side timeout |
| `errors.IsServiceUnavailable(err)` | 503 — server unavailable |

---

## Basic Error Handling

```go
pod, err := clientset.CoreV1().Pods("default").Get(ctx, "my-pod", metav1.GetOptions{})
if err != nil {
    switch {
    case errors.IsNotFound(err):
        fmt.Println("Pod not found")
    case errors.IsForbidden(err):
        fmt.Println("Permission denied — check RBAC")
    case errors.IsUnauthorized(err):
        fmt.Println("Invalid credentials — refresh token")
    default:
        log.Fatalf("Unexpected error: %v", err)
    }
    return
}
```

---

## RetryOnConflict

The most important retry pattern — handles 409 Conflict from stale resourceVersion:

```go
import "k8s.io/client-go/util/retry"

err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
    // Always GET first to get the latest resourceVersion
    dep, err := clientset.AppsV1().Deployments("default").Get(ctx, "my-app", metav1.GetOptions{})
    if err != nil {
        return err
    }

    // Modify
    dep.Spec.Replicas = int32Ptr(5)

    // Update — may fail with 409 if another writer updated concurrently
    _, err = clientset.AppsV1().Deployments("default").Update(ctx, dep, metav1.UpdateOptions{})
    return err
})
```

### DefaultRetry Settings

```go
retry.DefaultRetry = wait.Backoff{
    Duration: 10 * time.Millisecond,
    Factor:   1.0,
    Jitter:   0.1,
    Steps:    5,
}

// More aggressive retry
retry.RetryOnConflict(wait.Backoff{
    Duration: 100 * time.Millisecond,
    Factor:   2.0,
    Jitter:   0.1,
    Steps:    10,
}, func() error {
    // ...
    return nil
})
```

---

## Custom Retry with Backoff

```go
import (
    "k8s.io/apimachinery/pkg/util/wait"
)

func withRetry(ctx context.Context, maxAttempts int, fn func() error) error {
    backoff := wait.Backoff{
        Duration: 500 * time.Millisecond,
        Factor:   2.0,
        Jitter:   0.1,
        Steps:    maxAttempts,
    }

    return wait.ExponentialBackoffWithContext(ctx, backoff, func(ctx context.Context) (bool, error) {
        err := fn()
        if err == nil {
            return true, nil
        }

        // Retry on transient errors
        if errors.IsServerTimeout(err) || errors.IsServiceUnavailable(err) || errors.IsTooManyRequests(err) {
            return false, nil  // retry
        }

        // Don't retry on permanent errors
        return false, err
    })
}

// Usage
err := withRetry(ctx, 5, func() error {
    _, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
    return err
})
```

---

## Apply Pattern (Create or Update)

```go
func applyDeployment(ctx context.Context, clientset *kubernetes.Clientset, dep *appsv1.Deployment) error {
    _, err := clientset.AppsV1().Deployments(dep.Namespace).Get(ctx, dep.Name, metav1.GetOptions{})

    if errors.IsNotFound(err) {
        _, err = clientset.AppsV1().Deployments(dep.Namespace).Create(ctx, dep, metav1.CreateOptions{})
        return err
    }
    if err != nil {
        return err
    }

    return retry.RetryOnConflict(retry.DefaultRetry, func() error {
        existing, err := clientset.AppsV1().Deployments(dep.Namespace).Get(ctx, dep.Name, metav1.GetOptions{})
        if err != nil {
            return err
        }
        dep.ResourceVersion = existing.ResourceVersion
        _, err = clientset.AppsV1().Deployments(dep.Namespace).Update(ctx, dep, metav1.UpdateOptions{})
        return err
    })
}
```

---

## Delete If Exists

```go
func deleteIfExists(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) error {
    err := clientset.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
    if errors.IsNotFound(err) {
        return nil  // already gone — that's fine
    }
    return err
}
```

---

## Context Timeout

```go
// Always use context timeouts for API calls
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
if err != nil {
    if ctx.Err() != nil {
        fmt.Println("API call timed out")
    } else {
        fmt.Printf("API error: %v\n", err)
    }
}
```

---

## Structured Error Logging

```go
import (
    "k8s.io/apimachinery/pkg/api/errors"
    "k8s.io/apimachinery/pkg/runtime/schema"
)

func logK8sError(operation, resource, name string, err error) {
    if statusErr, ok := err.(*errors.StatusError); ok {
        fmt.Printf("K8s error: op=%s resource=%s name=%s status=%d reason=%s message=%s\n",
            operation, resource, name,
            statusErr.Status().Code,
            statusErr.Status().Reason,
            statusErr.Status().Message,
        )
    } else {
        fmt.Printf("Error: op=%s resource=%s name=%s err=%v\n", operation, resource, name, err)
    }
}
```

---

## Handling NotFound in Controllers

```go
func (c *Controller) reconcile(key string) error {
    namespace, name, err := cache.SplitMetaNamespaceKey(key)
    if err != nil {
        return err
    }

    pod, err := c.podLister.Pods(namespace).Get(name)
    if errors.IsNotFound(err) {
        // Pod was deleted — clean up any dependent resources
        fmt.Printf("Pod %s/%s deleted, cleaning up\n", namespace, name)
        return nil  // not an error — expected in deletion flow
    }
    if err != nil {
        return err  // will be retried
    }

    // Reconcile the existing pod
    return c.syncPod(pod)
}
```

---

## Server-Side Apply (SSA)

```go
import (
    "encoding/json"
    "k8s.io/apimachinery/pkg/types"
)

// Server-Side Apply — the safest way to manage resources
// (no read-modify-write, no conflict on resourceVersion)
patchData := map[string]interface{}{
    "apiVersion": "v1",
    "kind":       "ConfigMap",
    "metadata":   map[string]interface{}{"name": "app-config", "namespace": "default"},
    "data":       map[string]interface{}{"KEY": "value"},
}

data, _ := json.Marshal(patchData)
_, err = clientset.CoreV1().ConfigMaps("default").Patch(
    ctx,
    "app-config",
    types.ApplyPatchType,
    data,
    metav1.PatchOptions{
        FieldManager: "my-controller",
        Force:        boolPtr(true),
    },
)
```

---

## Common Error Patterns

```go
// Pattern: Create only if not exists
_, err := clientset.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
if err != nil && !errors.IsAlreadyExists(err) {
    return err
}

// Pattern: Update with conflict retry
err = retry.RetryOnConflict(retry.DefaultRetry, func() error { ... })

// Pattern: Delete idempotently
err = clientset.AppsV1().Deployments(ns).Delete(ctx, name, metav1.DeleteOptions{})
if err != nil && !errors.IsNotFound(err) {
    return err
}

// Pattern: Check RBAC before proceeding
_, err = clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
if errors.IsForbidden(err) {
    log.Fatal("Need permission: pods/list in default namespace")
}
```
