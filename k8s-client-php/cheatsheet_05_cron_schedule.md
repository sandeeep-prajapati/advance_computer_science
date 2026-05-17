# k8s-client-php — CronJob Schedule Cheatsheet

---

## Cron Syntax

```
┌─────────── minute        (0–59)
│ ┌───────── hour          (0–23)
│ │ ┌─────── day of month  (1–31)
│ │ │ ┌───── month         (1–12 or JAN–DEC)
│ │ │ │ ┌─── day of week   (0–6 or SUN–SAT, 0=Sunday)
│ │ │ │ │
* * * * *
```

---

## Common Schedules

| Schedule | Meaning | Use Case |
|---|---|---|
| `* * * * *` | Every minute | Testing only |
| `*/5 * * * *` | Every 5 minutes | Health checks, polling |
| `*/15 * * * *` | Every 15 minutes | Metrics collection |
| `0 * * * *` | Every hour (on the hour) | Hourly reports |
| `0 */2 * * *` | Every 2 hours | Cache warming |
| `0 6 * * *` | Every day at 6am UTC | Daily jobs |
| `0 6 * * 1-5` | Mon–Fri at 6am | Business hours jobs |
| `0 0 * * 0` | Every Sunday midnight | Weekly cleanup |
| `0 0 1 * *` | 1st of every month | Monthly billing |
| `0 0 1 1 *` | 1st Jan every year | Annual archival |
| `30 3 * * *` | Every day at 3:30am | Backups (off-peak) |
| `0 9,17 * * 1-5` | 9am & 5pm weekdays | Twice-daily sync |

---

## PHP CronJob Creation

```php
use Maclof\Kubernetes\Client;
use Maclof\Kubernetes\Models\CronJob;

$client = Client::fromKubeConfig();

function createCronJob(Client $client, string $name, string $schedule, array $command, string $image): void
{
    $cronJob = new CronJob([
        'metadata' => ['name' => $name, 'namespace' => 'default'],
        'spec' => [
            'schedule'                   => $schedule,
            'concurrencyPolicy'          => 'Forbid',
            'successfulJobsHistoryLimit' => 3,
            'failedJobsHistoryLimit'     => 1,
            'startingDeadlineSeconds'    => 300,
            'jobTemplate' => ['spec' => [
                'backoffLimit'            => 2,
                'ttlSecondsAfterFinished' => 86400,
                'template' => ['spec' => [
                    'restartPolicy' => 'OnFailure',
                    'containers'    => [['name' => 'job', 'image' => $image, 'command' => $command]],
                ]],
            ]],
        ],
    ]);
    $client->cronJobs()->create($cronJob);
    echo "CronJob '$name' created: $schedule\n";
}

// Examples
createCronJob($client, 'daily-report',   '0 6 * * *',     ['php', 'artisan', 'report:daily'],   'myapp:latest');
createCronJob($client, 'hourly-sync',    '0 * * * *',     ['php', 'artisan', 'sync:data'],       'myapp:latest');
createCronJob($client, 'weekly-cleanup', '0 0 * * 0',     ['php', 'artisan', 'cleanup:old'],     'myapp:latest');
createCronJob($client, 'every-5min',     '*/5 * * * *',   ['php', 'artisan', 'queue:heartbeat'], 'myapp:latest');
```

---

## ConcurrencyPolicy Options

```php
'concurrencyPolicy' => 'Allow',    // Run overlapping jobs (default)
'concurrencyPolicy' => 'Forbid',   // Skip if previous still running
'concurrencyPolicy' => 'Replace',  // Kill previous, start new
```

---

## Manage CronJobs

```php
// List all
$client->cronJobs()->find()

// Get one
$cj = $client->cronJobs()->findByName('daily-report')

// Check schedule
$cj->getSpec()['schedule']
$cj->getStatus()['lastScheduleTime']
$cj->getSpec()['suspend']

// Suspend (pause)
$client->cronJobs()->patch($cj, ['spec' => ['suspend' => true]])

// Resume
$client->cronJobs()->patch($cj, ['spec' => ['suspend' => false]])

// Change schedule
$client->cronJobs()->patch($cj, ['spec' => ['schedule' => '0 8 * * *']])

// Delete
$client->cronJobs()->deleteByName('daily-report')

// Suspend all
foreach ($client->cronJobs()->find() as $cj) {
    $client->cronJobs()->patch($cj, ['spec' => ['suspend' => true]]);
}
```

---

## Manually Trigger a CronJob

```php
function triggerCronJob(Client $client, string $cronJobName): string
{
    $cj      = $client->cronJobs()->findByName($cronJobName);
    $jobName = $cronJobName . '-manual-' . time();

    $job = new \Maclof\Kubernetes\Models\Job([
        'metadata' => ['name' => $jobName, 'namespace' => 'default'],
        'spec'     => $cj->getSpec()['jobTemplate']['spec'],
    ]);

    $client->jobs()->create($job);
    echo "Triggered: $jobName\n";
    return $jobName;
}

$jobName = triggerCronJob($client, 'daily-report');
// Then monitor: waitForJob($client, $jobName)
```

---

## CronJob vs Job Quick Pick

```
Need to run something...
  ├── Once, right now             → Job
  ├── Once, at a specific time    → CronJob (one-time schedule) or Job with delay
  ├── Every N minutes             → CronJob  */N * * * *
  ├── Daily                       → CronJob  0 H * * *
  ├── Weekly                      → CronJob  0 H * * DOW
  └── Monthly                     → CronJob  0 H 1 * *
```

---

## Timezone Reference

```php
'timezone' => 'UTC',                    // Default
'timezone' => 'America/New_York',       // EST/EDT
'timezone' => 'America/Los_Angeles',    // PST/PDT
'timezone' => 'Europe/London',          // GMT/BST
'timezone' => 'Europe/Paris',           // CET/CEST
'timezone' => 'Asia/Tokyo',             // JST
'timezone' => 'Asia/Kolkata',           // IST
'timezone' => 'Australia/Sydney',       // AEST/AEDT
```

> **Requires** Kubernetes 1.27+ and `CronJobTimeZone` feature gate enabled.

---

## List Jobs Created by a CronJob

```php
// Jobs created by a CronJob have a label: job-name={cronjob-name}-{timestamp}
// Or find by ownerReferences

$jobs = $client->jobs()->find();
foreach ($jobs as $job) {
    $owners = $job->getMetadata()['ownerReferences'] ?? [];
    foreach ($owners as $owner) {
        if ($owner['kind'] === 'CronJob' && $owner['name'] === 'daily-report') {
            $name      = $job->getMetadata('name');
            $startTime = $job->getStatus()['startTime'] ?? '-';
            $succeeded = $job->getStatus()['succeeded'] ?? 0;
            echo "$name  started=$startTime  succeeded=$succeeded\n";
        }
    }
}
```
