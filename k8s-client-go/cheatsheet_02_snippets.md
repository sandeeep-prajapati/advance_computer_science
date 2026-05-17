# k8s-client-go — Code Snippets Cheatsheet

Copy-paste patterns for real-world tasks.

---

## Auto-detect Config (In-cluster or kubeconfig)

```go
func buildConfig() (*rest.Config, error) {
    if config, err := rest.InClusterConfig(); err == nil {
        return config, nil
    }
    home, _ := os.UserHomeDir()
    return clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
}
```

---

## Create or Update (Apply)

```go
func applyDeployment(ctx context.Context, cs *kubernetes.Clientset, dep *appsv1.Deployment) (string, error) {
    _, err := cs.AppsV1().Deployments(dep.Namespace).Get(ctx, dep.Name, metav1.GetOptions{})
    if errors.IsNotFound(err) {
        _, err = cs.AppsV1().Deployments(dep.Namespace).Create(ctx, dep, metav1.CreateOptions{})
        return "created", err
    }
    if err != nil {
        return "", err
    }
    err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
        existing, err := cs.AppsV1().Deployments(dep.Namespace).Get(ctx, dep.Name, metav1.GetOptions{})
        if err != nil {
            return err
        }
        dep.ResourceVersion = existing.ResourceVersion
        _, err = cs.AppsV1().Deployments(dep.Namespace).Update(ctx, dep, metav1.UpdateOptions{})
        return err
    })
    return "updated", err
}
```

---

## Delete If Exists

```go
func deleteIfExists(ctx context.Context, cs *kubernetes.Clientset, ns, name string) error {
    err := cs.AppsV1().Deployments(ns).Delete(ctx, name, metav1.DeleteOptions{})
    if errors.IsNotFound(err) {
        return nil
    }
    return err
}
```

---

## Wait for Pod Running

```go
func waitForPodRunning(ctx context.Context, cs *kubernetes.Clientset, ns, name string) error {
    for {
        pod, err := cs.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }
        if pod.Status.Phase == corev1.PodRunning {
            return nil
        }
        if pod.Status.Phase == corev1.PodFailed {
            return fmt.Errorf("pod %s failed", name)
        }
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(2 * time.Second):
        }
    }
}
```

---

## Wait for Deployment Rollout

```go
func waitForDeployment(ctx context.Context, cs *kubernetes.Clientset, ns, name string) error {
    for {
        dep, err := cs.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }
        d := *dep.Spec.Replicas
        if dep.Status.UpdatedReplicas == d && dep.Status.ReadyReplicas == d && dep.Status.AvailableReplicas == d {
            return nil
        }
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(3 * time.Second):
        }
    }
}
```

---

## Wait for Job Completion

```go
func waitForJob(ctx context.Context, cs *kubernetes.Clientset, ns, name string) (string, error) {
    for {
        job, err := cs.BatchV1().Jobs(ns).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return "", err
        }
        if job.Status.Succeeded > 0 {
            return "succeeded", nil
        }
        if job.Status.Failed > *job.Spec.BackoffLimit {
            return "failed", fmt.Errorf("job failed after %d attempts", job.Status.Failed)
        }
        select {
        case <-ctx.Done():
            return "timeout", ctx.Err()
        case <-time.After(5 * time.Second):
        }
    }
}
```

---

## Wait for PVC Bound

```go
func waitForPVC(ctx context.Context, cs *kubernetes.Clientset, ns, name string) error {
    for {
        pvc, err := cs.CoreV1().PersistentVolumeClaims(ns).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }
        if pvc.Status.Phase == corev1.ClaimBound {
            return nil
        }
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(3 * time.Second):
        }
    }
}
```

---

## Wait for LoadBalancer IP

```go
func waitForLBIP(ctx context.Context, cs *kubernetes.Clientset, ns, name string) (string, error) {
    for {
        svc, err := cs.CoreV1().Services(ns).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return "", err
        }
        for _, ing := range svc.Status.LoadBalancer.Ingress {
            if ing.IP != "" {
                return ing.IP, nil
            }
            if ing.Hostname != "" {
                return ing.Hostname, nil
            }
        }
        select {
        case <-ctx.Done():
            return "", ctx.Err()
        case <-time.After(5 * time.Second):
        }
    }
}
```

---

## Scale Deployment

```go
func scale(ctx context.Context, cs *kubernetes.Clientset, ns, name string, replicas int32) error {
    return retry.RetryOnConflict(retry.DefaultRetry, func() error {
        dep, err := cs.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }
        dep.Spec.Replicas = &replicas
        _, err = cs.AppsV1().Deployments(ns).Update(ctx, dep, metav1.UpdateOptions{})
        return err
    })
}

scale(ctx, clientset, "default", "my-app", 0)  // scale to zero
scale(ctx, clientset, "default", "my-app", 5)  // scale up
```

---

## Update Image

```go
func updateImage(ctx context.Context, cs *kubernetes.Clientset, ns, dep, container, image string) error {
    return retry.RetryOnConflict(retry.DefaultRetry, func() error {
        d, err := cs.AppsV1().Deployments(ns).Get(ctx, dep, metav1.GetOptions{})
        if err != nil {
            return err
        }
        for i, c := range d.Spec.Template.Spec.Containers {
            if c.Name == container {
                d.Spec.Template.Spec.Containers[i].Image = image
            }
        }
        _, err = cs.AppsV1().Deployments(ns).Update(ctx, d, metav1.UpdateOptions{})
        return err
    })
}
```

---

## Rollout Restart

```go
func rolloutRestart(ctx context.Context, cs *kubernetes.Clientset, ns, name string) error {
    return retry.RetryOnConflict(retry.DefaultRetry, func() error {
        dep, err := cs.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }
        if dep.Spec.Template.Annotations == nil {
            dep.Spec.Template.Annotations = map[string]string{}
        }
        dep.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
        _, err = cs.AppsV1().Deployments(ns).Update(ctx, dep, metav1.UpdateOptions{})
        return err
    })
}
```

---

## Get Pod Logs as String

```go
func getPodLogs(ctx context.Context, cs *kubernetes.Clientset, ns, name string, tailLines int64) (string, error) {
    req := cs.CoreV1().Pods(ns).GetLogs(name, &corev1.PodLogOptions{TailLines: &tailLines})
    stream, err := req.Stream(ctx)
    if err != nil {
        return "", err
    }
    defer stream.Close()
    data, err := io.ReadAll(stream)
    return string(data), err
}
```

---

## Get Pods for a Deployment

```go
func getDeploymentPods(ctx context.Context, cs *kubernetes.Clientset, ns, name string) ([]corev1.Pod, error) {
    dep, err := cs.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
    if err != nil {
        return nil, err
    }
    selector := labels.Set(dep.Spec.Selector.MatchLabels).String()
    pods, err := cs.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{LabelSelector: selector})
    if err != nil {
        return nil, err
    }
    return pods.Items, nil
}
```

---

## Provision Tenant Namespace

```go
func provisionTenant(ctx context.Context, cs *kubernetes.Clientset, team string) error {
    ns := team + "-workspace"
    cs.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
        ObjectMeta: metav1.ObjectMeta{Name: ns, Labels: map[string]string{"team": team}},
    }, metav1.CreateOptions{})

    cs.CoreV1().ServiceAccounts(ns).Create(ctx, &corev1.ServiceAccount{
        ObjectMeta: metav1.ObjectMeta{Name: team + "-sa"},
    }, metav1.CreateOptions{})

    cs.RbacV1().Roles(ns).Create(ctx, &rbacv1.Role{
        ObjectMeta: metav1.ObjectMeta{Name: team + "-role"},
        Rules: []rbacv1.PolicyRule{{
            APIGroups: []string{"", "apps", "batch"}, Resources: []string{"*"}, Verbs: []string{"*"},
        }},
    }, metav1.CreateOptions{})

    cs.RbacV1().RoleBindings(ns).Create(ctx, &rbacv1.RoleBinding{
        ObjectMeta: metav1.ObjectMeta{Name: team + "-binding"},
        Subjects:   []rbacv1.Subject{{Kind: "ServiceAccount", Name: team + "-sa", Namespace: ns}},
        RoleRef:    rbacv1.RoleRef{APIGroup: "rbac.authorization.k8s.io", Kind: "Role", Name: team + "-role"},
    }, metav1.CreateOptions{})

    return nil
}
```

---

## Retry with Exponential Backoff

```go
func withRetry(ctx context.Context, maxSteps int, fn func() error) error {
    backoff := wait.Backoff{
        Duration: 500 * time.Millisecond,
        Factor:   2.0,
        Jitter:   0.1,
        Steps:    maxSteps,
    }
    return wait.ExponentialBackoffWithContext(ctx, backoff, func(ctx context.Context) (bool, error) {
        err := fn()
        if err == nil {
            return true, nil
        }
        if errors.IsServerTimeout(err) || errors.IsServiceUnavailable(err) || errors.IsTooManyRequests(err) {
            return false, nil  // retry
        }
        return false, err  // permanent error
    })
}
```

---

## Cluster Pod Summary

```go
func podSummary(ctx context.Context, cs *kubernetes.Clientset) {
    pods, _ := cs.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
    phases := map[corev1.PodPhase]int{}
    for _, p := range pods.Items {
        phases[p.Status.Phase]++
    }
    for phase, count := range phases {
        fmt.Printf("%-12s %d\n", phase, count)
    }
}
```

---

## Watch for CrashLoopBackOff

```go
func watchCrashLoops(ctx context.Context, cs *kubernetes.Clientset) {
    factory := informers.NewSharedInformerFactory(cs, 30*time.Second)
    factory.Core().V1().Pods().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
        UpdateFunc: func(_, obj interface{}) {
            pod := obj.(*corev1.Pod)
            for _, cs := range pod.Status.ContainerStatuses {
                if cs.State.Waiting != nil && cs.State.Waiting.Reason == "CrashLoopBackOff" {
                    fmt.Printf("[ALERT] CrashLoopBackOff: pod=%s container=%s restarts=%d\n",
                        pod.Name, cs.Name, cs.RestartCount)
                }
            }
        },
    })
    factory.Start(ctx.Done())
    <-ctx.Done()
}
```
