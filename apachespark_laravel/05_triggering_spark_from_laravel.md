# Lesson 5: Triggering Spark from Laravel

## Simple `spark-submit`

Laravel can submit a Spark job through Symfony Process:

```php
use Symfony\Component\Process\Process;

$process = new Process([
    config('spark.submit_path'),
    '--master', config('spark.master'),
    config('spark.jobs_dir').'/customer_revenue.py',
    '--job-id', (string) $sparkJob->id,
    '--from', $fromDate,
    '--to', $toDate,
    '--output', $outputPath,
]);

$process->setTimeout(null);
$process->start();
```

Do this from a queue worker or command, not directly inside a web request.

---

## Queue-Based Submitter

Controller:

```php
$job = SparkJob::create([
    'type' => 'customer_revenue',
    'parameters' => compact('fromDate', 'toDate'),
]);

SubmitSparkJob::dispatch($job->id);
```

Queued job:

```php
final class SubmitSparkJob implements ShouldQueue
{
    public function handle(): void
    {
        $job = SparkJob::findOrFail($this->sparkJobId);
        app(SparkSubmitter::class)->submit($job);
    }
}
```

---

## Capture Application ID

Spark prints an application ID such as:

```text
application_1710000000000_0001
```

Store it:

```php
$sparkJob->update([
    'status' => 'submitted',
    'spark_application_id' => $applicationId,
]);
```

That ID helps debug the job in Spark UI.
