# Chapter 5: Deployments & ReplicaSets

## Deployment Interface

```go
clientset.AppsV1().Deployments(namespace)
```

---

## Listing Deployments

```go
deployments, err := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
if err != nil {
    log.Fatal(err)
}

for _, d := range deployments.Items {
    fmt.Printf("%-40s  %d/%d ready\n",
        d.Name,
        d.Status.ReadyReplicas,
        *d.Spec.Replicas,
    )
}
```

---

## Getting a Single Deployment

```go
dep, err := clientset.AppsV1().Deployments("default").Get(ctx, "my-app", metav1.GetOptions{})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Desired:   %d\n", *dep.Spec.Replicas)
fmt.Printf("Ready:     %d\n", dep.Status.ReadyReplicas)
fmt.Printf("Updated:   %d\n", dep.Status.UpdatedReplicas)
fmt.Printf("Available: %d\n", dep.Status.AvailableReplicas)
fmt.Printf("Image:     %s\n", dep.Spec.Template.Spec.Containers[0].Image)
```

---

## Creating a Deployment

```go
import (
    appsv1  "k8s.io/api/apps/v1"
    corev1  "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
    "k8s.io/apimachinery/pkg/util/intstr"
)

func int32Ptr(i int32) *int32 { return &i }

deployment := &appsv1.Deployment{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "my-app",
        Namespace: "default",
        Labels:    map[string]string{"app": "my-app"},
    },
    Spec: appsv1.DeploymentSpec{
        Replicas: int32Ptr(3),
        Selector: &metav1.LabelSelector{
            MatchLabels: map[string]string{"app": "my-app"},
        },
        Strategy: appsv1.DeploymentStrategy{
            Type: appsv1.RollingUpdateDeploymentStrategyType,
            RollingUpdate: &appsv1.RollingUpdateDeployment{
                MaxUnavailable: &intstr.IntOrString{Type: intstr.Int, IntVal: 1},
                MaxSurge:       &intstr.IntOrString{Type: intstr.Int, IntVal: 1},
            },
        },
        Template: corev1.PodTemplateSpec{
            ObjectMeta: metav1.ObjectMeta{
                Labels: map[string]string{"app": "my-app"},
            },
            Spec: corev1.PodSpec{
                Containers: []corev1.Container{
                    {
                        Name:  "my-app",
                        Image: "myregistry/my-app:1.0.0",
                        Ports: []corev1.ContainerPort{{ContainerPort: 8080}},
                        Resources: corev1.ResourceRequirements{
                            Requests: corev1.ResourceList{
                                corev1.ResourceCPU:    resource.MustParse("100m"),
                                corev1.ResourceMemory: resource.MustParse("128Mi"),
                            },
                            Limits: corev1.ResourceList{
                                corev1.ResourceCPU:    resource.MustParse("500m"),
                                corev1.ResourceMemory: resource.MustParse("256Mi"),
                            },
                        },
                        ReadinessProbe: &corev1.Probe{
                            ProbeHandler: corev1.ProbeHandler{
                                HTTPGet: &corev1.HTTPGetAction{
                                    Path: "/health",
                                    Port: intstr.FromInt(8080),
                                },
                            },
                            InitialDelaySeconds: 5,
                            PeriodSeconds:       10,
                        },
                    },
                },
            },
        },
    },
}

created, err := clientset.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
if err != nil {
    log.Fatalf("Error creating deployment: %v", err)
}
fmt.Printf("Deployment created: %s\n", created.Name)
```

---

## Scaling a Deployment

```go
import "k8s.io/client-go/util/retry"

func scaleDeployment(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string, replicas int32) error {
    return retry.RetryOnConflict(retry.DefaultRetry, func() error {
        dep, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }
        dep.Spec.Replicas = &replicas
        _, err = clientset.AppsV1().Deployments(namespace).Update(ctx, dep, metav1.UpdateOptions{})
        return err
    })
}

scaleDeployment(ctx, clientset, "default", "my-app", 5)
```

---

## Update Image (Rolling Update)

```go
func updateDeploymentImage(ctx context.Context, clientset *kubernetes.Clientset, namespace, name, container, image string) error {
    return retry.RetryOnConflict(retry.DefaultRetry, func() error {
        dep, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }

        for i, c := range dep.Spec.Template.Spec.Containers {
            if c.Name == container {
                dep.Spec.Template.Spec.Containers[i].Image = image
                break
            }
        }

        _, err = clientset.AppsV1().Deployments(namespace).Update(ctx, dep, metav1.UpdateOptions{})
        return err
    })
}

updateDeploymentImage(ctx, clientset, "default", "my-app", "my-app", "myregistry/my-app:2.0.0")
```

---

## Patch a Deployment (Strategic Merge Patch)

```go
import (
    "encoding/json"
    "k8s.io/apimachinery/pkg/types"
)

type patchReplicas struct {
    Spec struct {
        Replicas int32 `json:"replicas"`
    } `json:"spec"`
}

patch := patchReplicas{}
patch.Spec.Replicas = 5

data, _ := json.Marshal(patch)
_, err = clientset.AppsV1().Deployments("default").Patch(
    ctx, "my-app",
    types.MergePatchType,
    data,
    metav1.PatchOptions{},
)
```

---

## Wait for Rollout Completion

```go
func waitForDeploymentRollout(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) error {
    for {
        dep, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }

        desired   := *dep.Spec.Replicas
        updated   := dep.Status.UpdatedReplicas
        ready     := dep.Status.ReadyReplicas
        available := dep.Status.AvailableReplicas

        fmt.Printf("updated=%d ready=%d available=%d desired=%d\n",
            updated, ready, available, desired)

        if updated == desired && ready == desired && available == desired {
            fmt.Println("Rollout complete!")
            return nil
        }

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(3 * time.Second):
        }
    }
}

timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
waitForDeploymentRollout(timeoutCtx, clientset, "default", "my-app")
```

---

## Rollout Restart (Force Pod Recreation)

```go
func rolloutRestart(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) error {
    return retry.RetryOnConflict(retry.DefaultRetry, func() error {
        dep, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }

        if dep.Spec.Template.Annotations == nil {
            dep.Spec.Template.Annotations = map[string]string{}
        }
        dep.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

        _, err = clientset.AppsV1().Deployments(namespace).Update(ctx, dep, metav1.UpdateOptions{})
        return err
    })
}
```

---

## Delete a Deployment

```go
// Delete deployment and all its pods (propagation)
propagation := metav1.DeletePropagationForeground
err := clientset.AppsV1().Deployments("default").Delete(ctx, "my-app", metav1.DeleteOptions{
    PropagationPolicy: &propagation,
})
```

---

## Working with ReplicaSets

```go
// List ReplicaSets for a deployment
rsets, err := clientset.AppsV1().ReplicaSets("default").List(ctx, metav1.ListOptions{
    LabelSelector: "app=my-app",
})

for _, rs := range rsets.Items {
    fmt.Printf("%-40s  desired=%d ready=%d\n",
        rs.Name, *rs.Spec.Replicas, rs.Status.ReadyReplicas)
}
```

---

## Deployment Status Conditions

```go
dep, _ := clientset.AppsV1().Deployments("default").Get(ctx, "my-app", metav1.GetOptions{})

for _, cond := range dep.Status.Conditions {
    fmt.Printf("%-25s %-8s %s\n", cond.Type, cond.Status, cond.Message)
}
// Types: Available, Progressing, ReplicaFailure
```
