# Chapter 12: Jobs & CronJobs

## Jobs vs CronJobs vs Deployments

| Type | Purpose | Restarts on Completion |
|---|---|---|
| `Deployment` | Long-running service | Yes (Always) |
| `Job` | One-off task to completion | Only on failure |
| `CronJob` | Scheduled recurring task | On schedule |

---

## Working with Jobs

### List Jobs

```php
use Maclof\Kubernetes\Client;

$client = Client::fromKubeConfig();

$jobs = $client->jobs()->find();

foreach ($jobs as $job) {
    $name      = $job->getMetadata('name');
    $active    = $job->getStatus()['active'] ?? 0;
    $succeeded = $job->getStatus()['succeeded'] ?? 0;
    $failed    = $job->getStatus()['failed'] ?? 0;
    $startTime = $job->getStatus()['startTime'] ?? '-';

    printf("%-30s active=%-3d succeeded=%-3d failed=%-3d started=%s\n",
        $name, $active, $succeeded, $failed, $startTime);
}
```

### Create a Job

```php
use Maclof\Kubernetes\Models\Job;

$job = new Job([
    'metadata' => [
        'name'      => 'db-migration',
        'namespace' => 'default',
        'labels'    => ['type' => 'migration'],
    ],
    'spec' => [
        'completions'  => 1,         // number of successful pod completions needed
        'parallelism'  => 1,         // max pods running in parallel
        'backoffLimit' => 3,         // retry on failure (max 3 times)
        'ttlSecondsAfterFinished' => 3600,  // auto-delete 1 hour after completion
        'template' => [
            'metadata' => ['labels' => ['job' => 'db-migration']],
            'spec' => [
                'restartPolicy' => 'OnFailure',  // or 'Never'
                'containers' => [[
                    'name'    => 'migrator',
                    'image'   => 'myapp:latest',
                    'command' => ['php', 'artisan', 'migrate', '--force'],
                    'env' => [
                        ['name' => 'DB_HOST', 'value' => 'postgres'],
                        ['name' => 'DB_PASSWORD', 'valueFrom' => [
                            'secretKeyRef' => ['name' => 'db-secret', 'key' => 'password'],
                        ]],
                    ],
                    'resources' => [
                        'requests' => ['memory' => '128Mi', 'cpu' => '100m'],
                        'limits'   => ['memory' => '256Mi', 'cpu' => '500m'],
                    ],
                ]],
            ],
        ],
    ],
]);

$created = $client->jobs()->create($job);
echo "Job created: " . $created->getMetadata('name') . "\n";
```

### Wait for Job Completion

```php
function waitForJob(Client $client, string $jobName, int $timeout = 300): string
{
    $start = time();
    while (time() - $start < $timeout) {
        $job       = $client->jobs()->findByName($jobName);
        $status    = $job->getStatus();
        $succeeded = $status['succeeded'] ?? 0;
        $failed    = $status['failed'] ?? 0;
        $active    = $status['active'] ?? 0;

        if ($succeeded > 0) {
            echo "Job '$jobName' completed successfully\n";
            return 'succeeded';
        }

        $backoffLimit = $job->getSpec()['backoffLimit'] ?? 3;
        if ($failed > $backoffLimit) {
            echo "Job '$jobName' failed after $failed attempts\n";
            return 'failed';
        }

        echo "Job status: active=$active succeeded=$succeeded failed=$failed\n";
        sleep(5);
    }

    return 'timeout';
}

$client->jobs()->create($job);
$result = waitForJob($client, 'db-migration');

if ($result !== 'succeeded') {
    // Fetch logs from the failed pod
    $pods = $client->pods()->setLabelSelector(['job' => 'db-migration'])->find();
    foreach ($pods as $pod) {
        echo "=== Logs for " . $pod->getMetadata('name') . " ===\n";
        echo $client->pods()->logs($pod->getMetadata('name'), ['previous' => false]);
    }
}
```

### Parallel Job (Batch Processing)

```php
$job = new Job([
    'metadata' => ['name' => 'batch-processor'],
    'spec' => [
        'completions' => 10,    // process 10 items total
        'parallelism' => 3,     // 3 workers at a time
        'backoffLimit' => 2,
        'template' => [
            'spec' => [
                'restartPolicy' => 'Never',
                'containers' => [[
                    'name'    => 'worker',
                    'image'   => 'myapp:batch',
                    'command' => ['php', 'batch.php'],
                    'env'     => [
                        ['name' => 'JOB_INDEX', 'valueFrom' => [
                            'fieldRef' => ['fieldPath' => 'metadata.annotations[\'batch.kubernetes.io/job-completion-index\']'],
                        ]],
                    ],
                ]],
            ],
        ],
    ],
]);

$client->jobs()->create($job);
```

### Delete a Job

```php
// Delete job and its pods
$client->jobs()->deleteByName('db-migration', [
    'propagationPolicy' => 'Foreground',
]);
```

---

## Working with CronJobs

### List CronJobs

```php
$cronJobs = $client->cronJobs()->find();

foreach ($cronJobs as $cj) {
    $name          = $cj->getMetadata('name');
    $schedule      = $cj->getSpec()['schedule'];
    $lastSchedule  = $cj->getStatus()['lastScheduleTime'] ?? 'never';
    $active        = count($cj->getStatus()['active'] ?? []);
    $suspended     = $cj->getSpec()['suspend'] ? 'SUSPENDED' : 'active';

    printf("%-30s %-15s %-20s active=%-2d %s\n",
        $name, $schedule, $lastSchedule, $active, $suspended);
}
```

### Create a CronJob

```php
use Maclof\Kubernetes\Models\CronJob;

$cronJob = new CronJob([
    'metadata' => [
        'name'      => 'daily-report',
        'namespace' => 'default',
    ],
    'spec' => [
        'schedule'                   => '0 6 * * *',   // every day at 6am UTC
        'timezone'                   => 'America/New_York',
        'concurrencyPolicy'          => 'Forbid',       // don't run if previous still running
        'successfulJobsHistoryLimit' => 3,
        'failedJobsHistoryLimit'     => 1,
        'suspend'                    => false,
        'startingDeadlineSeconds'    => 300,            // skip if > 5 min late
        'jobTemplate' => [
            'spec' => [
                'backoffLimit' => 2,
                'ttlSecondsAfterFinished' => 86400,
                'template' => [
                    'spec' => [
                        'restartPolicy' => 'OnFailure',
                        'containers' => [[
                            'name'    => 'reporter',
                            'image'   => 'myapp:latest',
                            'command' => ['php', 'artisan', 'report:daily'],
                            'resources' => [
                                'requests' => ['memory' => '256Mi', 'cpu' => '200m'],
                                'limits'   => ['memory' => '512Mi', 'cpu' => '1'],
                            ],
                        ]],
                    ],
                ],
            ],
        ],
    ],
]);

$client->cronJobs()->create($cronJob);
echo "CronJob 'daily-report' scheduled: 0 6 * * *\n";
```

### CronJob Schedule Reference

```
# ┌────── minute (0-59)
# │ ┌──── hour (0-23)
# │ │ ┌── day of month (1-31)
# │ │ │ ┌ month (1-12)
# │ │ │ │ ┌ day of week (0-6, Sun=0)
# │ │ │ │ │
  * * * * *

Examples:
  0 * * * *     = every hour
  0 6 * * *     = every day at 6am
  */15 * * * *  = every 15 minutes
  0 6 * * 1     = every Monday at 6am
  0 0 1 * *     = first of every month at midnight
```

### Suspend / Resume a CronJob

```php
$cj = $client->cronJobs()->findByName('daily-report');

// Suspend
$client->cronJobs()->patch($cj, ['spec' => ['suspend' => true]]);
echo "CronJob suspended\n";

// Resume
$client->cronJobs()->patch($cj, ['spec' => ['suspend' => false]]);
echo "CronJob resumed\n";
```

### Manually Trigger a CronJob

```php
function triggerCronJob(Client $client, string $cronJobName): string
{
    $cj = $client->cronJobs()->findByName($cronJobName);

    $jobSpec     = $cj->getSpec()['jobTemplate']['spec'];
    $manualJobName = $cronJobName . '-manual-' . time();

    $job = new Job([
        'metadata' => [
            'name'            => $manualJobName,
            'namespace'       => 'default',
            'ownerReferences' => [],
            'annotations'     => ['triggered-by' => 'manual'],
        ],
        'spec' => $jobSpec,
    ]);

    $client->jobs()->create($job);
    echo "Manual trigger: $manualJobName\n";

    return $manualJobName;
}

$jobName = triggerCronJob($client, 'daily-report');
$result  = waitForJob($client, $jobName);
echo "Manual run result: $result\n";
```

### List Recent Job Runs from a CronJob

```php
$cj = $client->cronJobs()->findByName('daily-report');

// Active jobs
$activeJobs = $cj->getStatus()['active'] ?? [];
echo "Currently running: " . count($activeJobs) . "\n";

// Find historical jobs by label
$historicalJobs = $client->jobs()->setLabelSelector([
    'app' => 'daily-report',
])->find();

foreach ($historicalJobs as $job) {
    $name      = $job->getMetadata('name');
    $startTime = $job->getStatus()['startTime'] ?? '-';
    $endTime   = $job->getStatus()['completionTime'] ?? 'running';
    $succeeded = $job->getStatus()['succeeded'] ?? 0;

    echo "$name  start=$startTime  end=$endTime  success=$succeeded\n";
}
```

---

## Concurrency Policy Options

| Policy | Behavior |
|---|---|
| `Allow` (default) | Multiple instances can run simultaneously |
| `Forbid` | Skip if previous run still running |
| `Replace` | Cancel previous run, start new one |
