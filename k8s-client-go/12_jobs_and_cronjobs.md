# Chapter 12: Jobs & CronJobs

## Interfaces

```go
clientset.BatchV1().Jobs(namespace)
clientset.BatchV1().CronJobs(namespace)
```

---

## Working with Jobs

### List Jobs

```go
jobs, err := clientset.BatchV1().Jobs("default").List(ctx, metav1.ListOptions{})
for _, job := range jobs.Items {
    active    := job.Status.Active
    succeeded := job.Status.Succeeded
    failed    := job.Status.Failed
    fmt.Printf("%-30s active=%-3d succeeded=%-3d failed=%d\n",
        job.Name, active, succeeded, failed)
}
```

### Create a Job

```go
import batchv1 "k8s.io/api/batch/v1"

backoffLimit    := int32(3)
ttlSeconds      := int32(3600)
completions     := int32(1)
parallelism     := int32(1)

job := &batchv1.Job{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "db-migration",
        Namespace: "default",
        Labels:    map[string]string{"type": "migration"},
    },
    Spec: batchv1.JobSpec{
        Completions:             &completions,
        Parallelism:             &parallelism,
        BackoffLimit:            &backoffLimit,
        TTLSecondsAfterFinished: &ttlSeconds,
        Template: corev1.PodTemplateSpec{
            ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"job": "db-migration"}},
            Spec: corev1.PodSpec{
                RestartPolicy: corev1.RestartPolicyOnFailure,
                Containers: []corev1.Container{
                    {
                        Name:    "migrator",
                        Image:   "myapp:latest",
                        Command: []string{"go", "run", "migrate.go"},
                        Env: []corev1.EnvVar{
                            {Name: "DB_HOST", Value: "postgres"},
                            {
                                Name: "DB_PASSWORD",
                                ValueFrom: &corev1.EnvVarSource{
                                    SecretKeyRef: &corev1.SecretKeySelector{
                                        LocalObjectReference: corev1.LocalObjectReference{Name: "db-secret"},
                                        Key: "password",
                                    },
                                },
                            },
                        },
                        Resources: corev1.ResourceRequirements{
                            Requests: corev1.ResourceList{
                                corev1.ResourceCPU:    resource.MustParse("100m"),
                                corev1.ResourceMemory: resource.MustParse("128Mi"),
                            },
                        },
                    },
                },
            },
        },
    },
}

created, err := clientset.BatchV1().Jobs("default").Create(ctx, job, metav1.CreateOptions{})
fmt.Printf("Job created: %s\n", created.Name)
```

### Wait for Job Completion

```go
func waitForJob(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) (string, error) {
    for {
        job, err := clientset.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return "", err
        }

        if job.Status.Succeeded > 0 {
            return "succeeded", nil
        }

        backoff := *job.Spec.BackoffLimit
        if job.Status.Failed > backoff {
            return "failed", fmt.Errorf("job %s failed after %d attempts", name, job.Status.Failed)
        }

        fmt.Printf("Job: active=%d succeeded=%d failed=%d\n",
            job.Status.Active, job.Status.Succeeded, job.Status.Failed)

        select {
        case <-ctx.Done():
            return "timeout", ctx.Err()
        case <-time.After(5 * time.Second):
        }
    }
}

timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

status, err := waitForJob(timeoutCtx, clientset, "default", "db-migration")
fmt.Printf("Job status: %s\n", status)
```

### Parallel Job (Batch Processing)

```go
completions  := int32(10)
parallelism  := int32(3)
backoffLimit := int32(2)

job := &batchv1.Job{
    ObjectMeta: metav1.ObjectMeta{Name: "batch-processor"},
    Spec: batchv1.JobSpec{
        Completions:  &completions,
        Parallelism:  &parallelism,
        BackoffLimit: &backoffLimit,
        Template: corev1.PodTemplateSpec{
            Spec: corev1.PodSpec{
                RestartPolicy: corev1.RestartPolicyNever,
                Containers: []corev1.Container{{
                    Name:    "worker",
                    Image:   "myapp:batch",
                    Command: []string{"/app/worker"},
                }},
            },
        },
    },
}
```

### Delete a Job

```go
propagation := metav1.DeletePropagationForeground

err := clientset.BatchV1().Jobs("default").Delete(ctx, "db-migration", metav1.DeleteOptions{
    PropagationPolicy: &propagation,
})
```

---

## Working with CronJobs

### List CronJobs

```go
cjs, err := clientset.BatchV1().CronJobs("default").List(ctx, metav1.ListOptions{})
for _, cj := range cjs.Items {
    lastSchedule := "never"
    if cj.Status.LastScheduleTime != nil {
        lastSchedule = cj.Status.LastScheduleTime.String()
    }
    suspended := ""
    if cj.Spec.Suspend != nil && *cj.Spec.Suspend {
        suspended = "SUSPENDED"
    }
    fmt.Printf("%-30s %-15s %-25s active=%-2d %s\n",
        cj.Name, cj.Spec.Schedule, lastSchedule, len(cj.Status.Active), suspended)
}
```

### Create a CronJob

```go
successLimit := int32(3)
failLimit    := int32(1)
suspend      := false
startDeadline := int64(300)

cronJob := &batchv1.CronJob{
    ObjectMeta: metav1.ObjectMeta{Name: "daily-report", Namespace: "default"},
    Spec: batchv1.CronJobSpec{
        Schedule:                   "0 6 * * *",
        ConcurrencyPolicy:          batchv1.ForbidConcurrent,
        SuccessfulJobsHistoryLimit: &successLimit,
        FailedJobsHistoryLimit:     &failLimit,
        Suspend:                    &suspend,
        StartingDeadlineSeconds:    &startDeadline,
        JobTemplate: batchv1.JobTemplateSpec{
            Spec: batchv1.JobSpec{
                BackoffLimit: int32Ptr(2),
                Template: corev1.PodTemplateSpec{
                    Spec: corev1.PodSpec{
                        RestartPolicy: corev1.RestartPolicyOnFailure,
                        Containers: []corev1.Container{{
                            Name:    "reporter",
                            Image:   "myapp:latest",
                            Command: []string{"go", "run", "report.go"},
                            Resources: corev1.ResourceRequirements{
                                Requests: corev1.ResourceList{
                                    corev1.ResourceCPU:    resource.MustParse("200m"),
                                    corev1.ResourceMemory: resource.MustParse("256Mi"),
                                },
                            },
                        }},
                    },
                },
            },
        },
    },
}

clientset.BatchV1().CronJobs("default").Create(ctx, cronJob, metav1.CreateOptions{})
```

### Suspend / Resume a CronJob

```go
import (
    "encoding/json"
    "k8s.io/apimachinery/pkg/types"
)

type suspendPatch struct {
    Spec struct {
        Suspend bool `json:"suspend"`
    } `json:"spec"`
}

func setCronJobSuspend(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string, suspend bool) error {
    p := suspendPatch{}
    p.Spec.Suspend = suspend
    data, err := json.Marshal(p)
    if err != nil {
        return err
    }
    _, err = clientset.BatchV1().CronJobs(namespace).Patch(
        ctx, name, types.MergePatchType, data, metav1.PatchOptions{})
    return err
}

setCronJobSuspend(ctx, clientset, "default", "daily-report", true)  // pause
setCronJobSuspend(ctx, clientset, "default", "daily-report", false) // resume
```

### Manually Trigger a CronJob

```go
func triggerCronJob(ctx context.Context, clientset *kubernetes.Clientset, namespace, cronJobName string) (string, error) {
    cj, err := clientset.BatchV1().CronJobs(namespace).Get(ctx, cronJobName, metav1.GetOptions{})
    if err != nil {
        return "", err
    }

    jobName := fmt.Sprintf("%s-manual-%d", cronJobName, time.Now().Unix())
    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name:        jobName,
            Namespace:   namespace,
            Annotations: map[string]string{"triggered-by": "manual"},
        },
        Spec: cj.Spec.JobTemplate.Spec,
    }

    created, err := clientset.BatchV1().Jobs(namespace).Create(ctx, job, metav1.CreateOptions{})
    if err != nil {
        return "", err
    }

    fmt.Printf("Triggered job: %s\n", created.Name)
    return created.Name, nil
}
```

---

## ConcurrencyPolicy Reference

```go
batchv1.AllowConcurrent    // multiple instances allowed (default)
batchv1.ForbidConcurrent   // skip if previous still running
batchv1.ReplaceConcurrent  // cancel previous, start new
```
