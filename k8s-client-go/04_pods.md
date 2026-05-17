# Chapter 4: Working with Pods

## Pod Interface

```go
clientset.CoreV1().Pods(namespace)
```

All methods accept a `context.Context` as the first argument — always pass one with a timeout.

---

## Listing Pods

```go
import (
    "context"
    "fmt"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

ctx := context.TODO()

// All pods in a namespace
pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
if err != nil {
    log.Fatalf("Error listing pods: %v", err)
}

for _, pod := range pods.Items {
    fmt.Printf("%-40s %-10s %s\n", pod.Name, pod.Status.Phase, pod.Spec.NodeName)
}
```

### Filter by Label Selector

```go
pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{
    LabelSelector: "app=nginx,env=production",
})

// Programmatic label selector
import "k8s.io/apimachinery/pkg/labels"

selector := labels.Set{"app": "nginx", "env": "production"}.String()
pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{
    LabelSelector: selector,
})
```

### Filter by Field Selector

```go
pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{
    FieldSelector: "status.phase=Running",
})

// Multiple field selectors
pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{
    FieldSelector: "status.phase=Running,spec.nodeName=node-1",
})
```

### All Namespaces

```go
pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
// "" = all namespaces (like kubectl get pods -A)
```

---

## Getting a Single Pod

```go
pod, err := clientset.CoreV1().Pods("default").Get(ctx, "my-nginx", metav1.GetOptions{})
if err != nil {
    if errors.IsNotFound(err) {
        fmt.Println("Pod not found")
        return
    }
    log.Fatalf("Error getting pod: %v", err)
}

fmt.Printf("Name:      %s\n", pod.Name)
fmt.Printf("Namespace: %s\n", pod.Namespace)
fmt.Printf("Phase:     %s\n", pod.Status.Phase)
fmt.Printf("Node:      %s\n", pod.Spec.NodeName)
fmt.Printf("PodIP:     %s\n", pod.Status.PodIP)
```

---

## Creating a Pod

```go
import (
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
)

pod := &corev1.Pod{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "my-nginx",
        Namespace: "default",
        Labels:    map[string]string{"app": "nginx", "env": "dev"},
    },
    Spec: corev1.PodSpec{
        Containers: []corev1.Container{
            {
                Name:  "nginx",
                Image: "nginx:1.25-alpine",
                Ports: []corev1.ContainerPort{
                    {ContainerPort: 80},
                },
                Resources: corev1.ResourceRequirements{
                    Requests: corev1.ResourceList{
                        corev1.ResourceCPU:    resource.MustParse("100m"),
                        corev1.ResourceMemory: resource.MustParse("64Mi"),
                    },
                    Limits: corev1.ResourceList{
                        corev1.ResourceCPU:    resource.MustParse("200m"),
                        corev1.ResourceMemory: resource.MustParse("128Mi"),
                    },
                },
                Env: []corev1.EnvVar{
                    {Name: "ENV", Value: "production"},
                },
            },
        },
        RestartPolicy: corev1.RestartPolicyAlways,
    },
}

created, err := clientset.CoreV1().Pods("default").Create(ctx, pod, metav1.CreateOptions{})
if err != nil {
    log.Fatalf("Error creating pod: %v", err)
}
fmt.Printf("Pod created: %s\n", created.Name)
```

---

## Pod Logs

```go
import (
    "io"
    corev1 "k8s.io/api/core/v1"
)

// Fetch logs
req := clientset.CoreV1().Pods("default").GetLogs("my-nginx", &corev1.PodLogOptions{})
stream, err := req.Stream(ctx)
if err != nil {
    log.Fatalf("Error getting logs: %v", err)
}
defer stream.Close()

logs, _ := io.ReadAll(stream)
fmt.Println(string(logs))
```

### Log Options

```go
tailLines := int64(100)
sinceSeconds := int64(3600)

// Last 100 lines
req := clientset.CoreV1().Pods("default").GetLogs("my-nginx", &corev1.PodLogOptions{
    TailLines: &tailLines,
})

// Last hour
req := clientset.CoreV1().Pods("default").GetLogs("my-nginx", &corev1.PodLogOptions{
    SinceSeconds: &sinceSeconds,
})

// Previous crashed container
req := clientset.CoreV1().Pods("default").GetLogs("my-nginx", &corev1.PodLogOptions{
    Previous: true,
})

// Specific container (multi-container pod)
req := clientset.CoreV1().Pods("default").GetLogs("my-nginx", &corev1.PodLogOptions{
    Container: "sidecar",
})

// Follow (streaming)
req := clientset.CoreV1().Pods("default").GetLogs("my-nginx", &corev1.PodLogOptions{
    Follow: true,
})
stream, _ := req.Stream(ctx)
io.Copy(os.Stdout, stream)
```

---

## Deleting a Pod

```go
// Delete by name
err := clientset.CoreV1().Pods("default").Delete(ctx, "my-nginx", metav1.DeleteOptions{})

// Force delete (grace period = 0)
gracePeriod := int64(0)
err := clientset.CoreV1().Pods("default").Delete(ctx, "my-nginx", metav1.DeleteOptions{
    GracePeriodSeconds: &gracePeriod,
})
```

---

## Exec into a Pod

```go
import (
    "k8s.io/client-go/tools/remotecommand"
    "k8s.io/client-go/kubernetes/scheme"
)

func execInPod(clientset *kubernetes.Clientset, config *rest.Config, namespace, podName, container string, command []string) (string, string, error) {
    req := clientset.CoreV1().RESTClient().Post().
        Resource("pods").
        Name(podName).
        Namespace(namespace).
        SubResource("exec").
        VersionedParams(&corev1.PodExecOptions{
            Container: container,
            Command:   command,
            Stdin:     false,
            Stdout:    true,
            Stderr:    true,
            TTY:       false,
        }, scheme.ParameterCodec)

    exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
    if err != nil {
        return "", "", err
    }

    var stdout, stderr bytes.Buffer
    err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
        Stdout: &stdout,
        Stderr: &stderr,
    })

    return stdout.String(), stderr.String(), err
}

// Usage
stdout, stderr, err := execInPod(clientset, config, "default", "my-pod", "app", []string{"ls", "/var/data"})
```

---

## Watch Pods

```go
watcher, err := clientset.CoreV1().Pods("default").Watch(ctx, metav1.ListOptions{})
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
}
```

---

## Wait for Pod Running

```go
import (
    "k8s.io/apimachinery/pkg/watch"
    "k8s.io/apimachinery/pkg/fields"
)

func waitForPodRunning(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) error {
    watcher, err := clientset.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{
        FieldSelector: fields.OneTermEqualSelector("metadata.name", name).String(),
    })
    if err != nil {
        return err
    }
    defer watcher.Stop()

    for event := range watcher.ResultChan() {
        pod, ok := event.Object.(*corev1.Pod)
        if !ok {
            continue
        }
        if pod.Status.Phase == corev1.PodRunning {
            return nil
        }
        if pod.Status.Phase == corev1.PodFailed {
            return fmt.Errorf("pod %s failed", name)
        }
    }
    return fmt.Errorf("watch channel closed")
}
```

---

## Container Status Check

```go
pod, _ := clientset.CoreV1().Pods("default").Get(ctx, "my-pod", metav1.GetOptions{})

for _, cs := range pod.Status.ContainerStatuses {
    fmt.Printf("Container: %-20s Ready: %-5v Restarts: %d\n",
        cs.Name, cs.Ready, cs.RestartCount)

    if cs.State.Waiting != nil {
        fmt.Printf("  Waiting: %s — %s\n", cs.State.Waiting.Reason, cs.State.Waiting.Message)
    }
    if cs.State.Terminated != nil {
        fmt.Printf("  Exit code: %d\n", cs.State.Terminated.ExitCode)
    }
}
```
