# Chapter 12: Jobs & CronJobs

## Jobs

### What is a Job?

A **Job** creates one or more Pods to run a task to **completion**. Unlike a Deployment (which keeps pods running forever), a Job:
- Runs a task once
- Marks itself as complete when the pod exits with code `0`
- Retries if the pod fails

**Use cases:**
- Database migrations
- Data processing batch jobs
- Report generation
- One-time setup tasks
- Image processing pipelines

---

### Basic Job

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: db-migration
spec:
  completions: 1          # How many successful pod completions needed
  parallelism: 1          # How many pods run at the same time
  backoffLimit: 3         # Retry up to 3 times before marking failed
  activeDeadlineSeconds: 600  # Kill job if takes > 10 minutes
  ttlSecondsAfterFinished: 3600  # Auto-delete job 1hr after completion

  template:
    spec:
      restartPolicy: OnFailure    # Required for Jobs (Never or OnFailure)
      containers:
        - name: migration
          image: my-app:latest
          command: ["bin/console", "doctrine:migrations:migrate", "--no-interaction"]
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: app-secrets
                  key: DATABASE_URL
          resources:
            requests:
              memory: "256Mi"
              cpu: "100m"
```

```bash
kubectl apply -f job.yaml
kubectl get jobs
kubectl get pods -l job-name=db-migration
kubectl logs job/db-migration
```

---

### Parallel Jobs

Run multiple pods at the same time to speed up processing:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: image-processor
spec:
  completions: 10      # Process 10 batches total
  parallelism: 3       # Run 3 pods simultaneously
  backoffLimit: 2
  template:
    spec:
      restartPolicy: OnFailure
      containers:
        - name: processor
          image: image-processor:latest
          env:
            - name: BATCH_SIZE
              value: "100"
```

---

### Job Patterns

| Pattern | completions | parallelism | Use Case |
|---|---|---|---|
| Single task | 1 | 1 | One-shot migration, setup |
| Parallel fixed count | N | N | Process N batches |
| Work queue | unset | N | Pull from a queue until empty |

---

### Job Lifecycle & Cleanup

```bash
# View all jobs
kubectl get jobs

# View job details (see success/failure counts)
kubectl describe job db-migration

# View pods created by the job
kubectl get pods --selector=job-name=db-migration

# Follow logs
kubectl logs -f job/db-migration

# Delete job and its pods
kubectl delete job db-migration

# Auto-cleanup: add this to job spec
spec:
  ttlSecondsAfterFinished: 300   # Delete 5 minutes after completion
```

---

## CronJobs

### What is a CronJob?

A **CronJob** creates Jobs on a schedule — like Linux `cron` but in Kubernetes.

**Use cases:**
- Nightly database backups
- Daily report generation
- Hourly cache warming
- Periodic cleanup tasks

---

### CronJob Syntax

```
┌─────── minute (0-59)
│ ┌───── hour (0-23)
│ │ ┌─── day of month (1-31)
│ │ │ ┌─ month (1-12)
│ │ │ │ └── day of week (0-7, 0=Sun)
│ │ │ │
* * * * *
```

| Schedule | Meaning |
|---|---|
| `0 2 * * *` | Every day at 2:00 AM |
| `*/15 * * * *` | Every 15 minutes |
| `0 0 * * 0` | Every Sunday at midnight |
| `0 9 * * 1-5` | Every weekday at 9 AM |
| `@daily` | Once a day (= `0 0 * * *`) |
| `@hourly` | Once an hour |
| `@weekly` | Once a week |

---

### CronJob Example: Nightly Backup

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: nightly-backup
  namespace: production
spec:
  schedule: "0 2 * * *"       # 2 AM every day
  timeZone: "Asia/Kolkata"     # Optional timezone (K8s 1.27+)

  concurrencyPolicy: Forbid    # Don't run if previous job still running
  successfulJobsHistoryLimit: 3  # Keep last 3 successful jobs
  failedJobsHistoryLimit: 1      # Keep last 1 failed job
  startingDeadlineSeconds: 300   # If missed, skip (don't catch up after 5min)

  jobTemplate:
    spec:
      backoffLimit: 2
      ttlSecondsAfterFinished: 86400   # Delete job after 24hrs
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: backup
              image: mysql:8.0
              command:
                - /bin/sh
                - -c
                - |
                  mysqldump -h mysql-service -u root -p$MYSQL_ROOT_PASSWORD mydb | \
                  gzip > /backup/mydb-$(date +%Y%m%d-%H%M%S).sql.gz
                  echo "Backup completed"
              env:
                - name: MYSQL_ROOT_PASSWORD
                  valueFrom:
                    secretKeyRef:
                      name: mysql-secret
                      key: root-password
              volumeMounts:
                - name: backup-storage
                  mountPath: /backup
          volumes:
            - name: backup-storage
              persistentVolumeClaim:
                claimName: backup-pvc
```

---

### CronJob Concurrency Policies

| Policy | Behavior |
|---|---|
| `Allow` (default) | Multiple jobs can run at the same time |
| `Forbid` | Skip new job if previous is still running |
| `Replace` | Cancel current job and start a new one |

---

### CronJob Commands

```bash
# View cronjobs
kubectl get cronjobs
kubectl get cj   # short form

# Describe (see last schedule, active jobs)
kubectl describe cronjob nightly-backup

# Manually trigger a cronjob right now
kubectl create job manual-backup-now \
  --from=cronjob/nightly-backup

# View jobs created by cronjob
kubectl get jobs -l app=nightly-backup

# Watch job execution
kubectl get jobs -w

# Suspend a cronjob (pause scheduling)
kubectl patch cronjob nightly-backup -p '{"spec":{"suspend":true}}'

# Resume
kubectl patch cronjob nightly-backup -p '{"spec":{"suspend":false}}'

# Delete
kubectl delete cronjob nightly-backup
```

---

### Common Job Failures & Fixes

| Issue | Cause | Fix |
|---|---|---|
| Job stuck in `Pending` | Resource constraints | Check node capacity, `kubectl describe pod` |
| Pod `Error` immediately | Container command fails | Check `kubectl logs pod/...` |
| Job never completes | `restartPolicy: Always` set | Change to `OnFailure` or `Never` |
| CronJob misses schedule | `startingDeadlineSeconds` too low | Increase or remove the deadline |
| Duplicate jobs running | `concurrencyPolicy: Allow` with slow job | Change to `Forbid` |
| Old jobs pile up | No `ttlSecondsAfterFinished` | Add TTL or history limits |
