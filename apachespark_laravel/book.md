# Apache Spark + Laravel — Lessons

> Practical lessons for integrating a PHP/Laravel application with Apache Spark for batch processing, analytics, and data pipelines.

---

## Table of Contents

| Lesson | Title | Topics |
|---|---|---|
| 01 | [Introduction](01_introduction.md) | Where Spark fits, Laravel's role, integration options |
| 02 | [Architecture](02_architecture.md) | Laravel app, object storage, Spark cluster, result callbacks |
| 03 | [Local Setup](03_local_setup.md) | Docker services, Laravel config, Spark master/worker |
| 04 | [Data Exchange Patterns](04_data_exchange_patterns.md) | Database, CSV, Parquet, S3/MinIO, APIs |
| 05 | [Triggering Spark from Laravel](05_triggering_spark_from_laravel.md) | `spark-submit`, REST, queues, job records |
| 06 | [Using Apache Livy](06_using_apache_livy.md) | Submit Spark jobs through HTTP from Laravel |
| 07 | [Laravel Export Jobs](07_laravel_export_jobs.md) | Exporting source data for Spark |
| 08 | [Spark Job Examples](08_spark_job_examples.md) | PySpark jobs for reports and transformations |
| 09 | [Reading Spark Results in Laravel](09_reading_spark_results.md) | Status tables, result files, callbacks |
| 10 | [Scheduling](10_scheduling.md) | Laravel Scheduler vs Spark jobs vs Airflow |
| 11 | [Reliability](11_reliability.md) | Idempotency, retries, deterministic output, backfills |
| 12 | [Security](12_security.md) | Service tokens, storage credentials, secrets |
| 13 | [Monitoring and Debugging](13_monitoring_debugging.md) | Spark UI, Laravel logs, job IDs, metrics |
| 14 | [Deployment](14_deployment.md) | Standalone, Kubernetes, EMR/Dataproc, production checklist |
| 15 | [Capstone Project](15_capstone_project.md) | Customer revenue analytics pipeline |

---

## Main Idea

Laravel should not process huge datasets inside web requests. Laravel should create jobs, export or point to data, trigger Spark, and display results. Spark should do distributed computation.

---

# Lesson 1: Introduction

## What Problem Are We Solving?

Laravel is excellent for web apps, APIs, business rules, forms, dashboards, authentication, and transactional database workflows.

Apache Spark is excellent for processing large datasets:

- millions or billions of rows
- joins across large files
- analytics reports
- ETL pipelines
- machine learning features
- batch transformations

The integration goal is simple:

```text
Laravel starts and tracks work
Spark processes large data
Laravel shows the result
```

---

## Laravel Should Own

- users and permissions
- business actions
- database transactions
- API validation
- job request records
- report status pages
- notifications

---

## Spark Should Own

- large file processing
- distributed joins
- aggregations
- ETL
- report generation
- feature computation
- batch analytics

---

## Integration Options

| Option | Best For |
|---|---|
| Laravel runs `spark-submit` | Simple internal deployments |
| Laravel calls Apache Livy | HTTP-based Spark job submission |
| Laravel writes job row, worker submits Spark | Reliable async execution |
| Airflow schedules Spark and calls Laravel | Multi-step data workflows |
| Laravel publishes event, pipeline consumes | Event-driven analytics |

PHP does not need to run Spark directly. Laravel usually talks to Spark through process execution, HTTP APIs, queues, files, or orchestration tools.

---

# Lesson 2: Architecture

## Recommended Architecture

```text
                 ┌────────────────────┐
User request --> │ Laravel App         │
                 │ - creates job row   │
                 │ - validates input   │
                 └─────────┬──────────┘
                           │
                           │ submit job
                           v
                 ┌────────────────────┐
                 │ Spark Cluster       │
                 │ - PySpark/Scala job │
                 └─────────┬──────────┘
                           │
                           │ reads/writes
                           v
                 ┌────────────────────┐
                 │ Storage             │
                 │ S3/MinIO/HDFS/DB    │
                 └─────────┬──────────┘
                           │
                           │ result metadata
                           v
                 ┌────────────────────┐
                 │ Laravel Dashboard   │
                 │ status + download   │
                 └────────────────────┘
```

---

## Job Table

Laravel should track Spark work in a database table:

```php
Schema::create('spark_jobs', function (Blueprint $table) {
    $table->id();
    $table->string('type');
    $table->string('status')->default('pending');
    $table->json('parameters');
    $table->string('spark_application_id')->nullable();
    $table->string('output_path')->nullable();
    $table->text('error_message')->nullable();
    $table->timestamp('started_at')->nullable();
    $table->timestamp('finished_at')->nullable();
    $table->timestamps();
});
```

---

## Status Flow

```text
pending -> submitted -> running -> completed
                         |
                         v
                       failed
```

The user should never wait in a browser request while Spark processes data.

---

# Lesson 3: Local Setup

## Local Services

For development, use:

```text
laravel-app
mysql or postgres
spark-master
spark-worker
minio or local mounted data folder
```

Optional:

```text
livy
airflow
kafka
```

---

## Laravel Environment

```dotenv
SPARK_MASTER=spark://spark-master:7077
SPARK_SUBMIT=/opt/spark/bin/spark-submit
SPARK_JOBS_DIR=/opt/spark/jobs
SPARK_OUTPUT_DISK=s3
SPARK_OUTPUT_BUCKET=analytics
```

---

## Docker Network Rule

Inside containers, use service names:

```text
spark-master:7077
mysql:3306
minio:9000
laravel-app:80
```

Do not use `localhost` unless the service is inside the same container.

---

## Smoke Test Spark

From a Spark container:

```bash
spark-submit --master spark://spark-master:7077 /opt/spark/examples/src/main/python/pi.py 10
```

From Laravel:

```bash
php artisan spark:test
```

The Laravel command can submit a tiny Spark job and store the result in `spark_jobs`.

---

# Lesson 4: Data Exchange Patterns

## Pattern 1: Spark Reads Laravel Database

Spark can read MySQL/Postgres through JDBC.

Pros:

- simple source of truth
- no export step

Cons:

- can load production database heavily
- needs network and credentials
- large reads may affect app performance

Use read replicas when possible.

---

## Pattern 2: Laravel Exports Files

Laravel writes CSV or JSON files for Spark:

```text
storage/app/exports/orders-2026-05-30.csv
```

Spark reads the file and writes transformed output.

Pros:

- protects app database
- easy to replay
- good for medium-sized exports

Cons:

- export can be slow in PHP for huge datasets
- CSV needs careful typing

---

## Pattern 3: Object Storage

Best production pattern:

```text
Laravel -> S3/MinIO raw data
Spark -> S3/MinIO processed data
Laravel -> reads metadata/result path
```

Prefer Parquet for Spark outputs:

- compressed
- columnar
- typed
- efficient for analytics

---

## Pattern 4: API Calls

Spark can call Laravel APIs, but this is not ideal for large data.

Use APIs for:

- metadata
- small lookup tables
- job status callbacks

Avoid APIs for millions of rows.

---

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

---

# Lesson 6: Using Apache Livy

## What is Livy?

Apache Livy provides an HTTP API for submitting Spark jobs. This is often cleaner than running `spark-submit` directly from Laravel.

Laravel sends:

```text
POST /batches
```

Livy starts a Spark application and returns a batch ID.

---

## Laravel Livy Client

```php
use Illuminate\Support\Facades\Http;

final class LivyClient
{
    public function submitBatch(array $payload): array
    {
        $response = Http::baseUrl(config('services.livy.base_url'))
            ->post('/batches', $payload);

        $response->throw();

        return $response->json();
    }
}
```

---

## Submit PySpark Job

```php
$batch = app(LivyClient::class)->submitBatch([
    'file' => 's3://analytics/jobs/customer_revenue.py',
    'args' => [
        '--job-id', (string) $sparkJob->id,
        '--from', $fromDate,
        '--to', $toDate,
        '--output', $outputPath,
    ],
    'conf' => [
        'spark.sql.shuffle.partitions' => '64',
    ],
]);

$sparkJob->update([
    'status' => 'submitted',
    'spark_application_id' => (string) $batch['id'],
]);
```

---

## Poll Status

```php
$response = Http::baseUrl(config('services.livy.base_url'))
    ->get("/batches/{$batchId}");

$state = $response->json('state');
```

Map Livy states to Laravel statuses:

| Livy | Laravel |
|---|---|
| `starting` | `submitted` |
| `running` | `running` |
| `success` | `completed` |
| `dead` | `failed` |
| `killed` | `failed` |

---

# Lesson 7: Laravel Export Jobs

## Export Source Data

Laravel can prepare clean source data before Spark runs.

```php
final class ExportOrdersForSpark implements ShouldQueue
{
    public function handle(): void
    {
        $path = "spark/input/orders-{$this->sparkJobId}.csv";

        Storage::disk('s3')->put($path, $this->buildCsv());

        SparkJob::find($this->sparkJobId)->update([
            'status' => 'exported',
            'parameters->input_path' => $path,
        ]);
    }
}
```

---

## Chunk Large Exports

Never load millions of rows into memory:

```php
Order::query()
    ->whereBetween('created_at', [$from, $to])
    ->orderBy('id')
    ->chunkById(1000, function ($orders) use ($writer) {
        foreach ($orders as $order) {
            $writer->insertOne([
                $order->id,
                $order->customer_id,
                $order->total,
                $order->created_at->toIso8601String(),
            ]);
        }
    });
```

---

## Prefer Typed Data

CSV is okay for learning. For production analytics, prefer:

- Parquet
- Delta Lake
- database read replica through JDBC
- CDC into a data lake

PHP is not usually the best tool for writing huge Parquet datasets, so Laravel often prepares job metadata while Spark reads source systems directly.

---

# Lesson 8: Spark Job Examples

## Customer Revenue PySpark Job

```python
import argparse
from pyspark.sql import SparkSession
from pyspark.sql.functions import sum as spark_sum, count

parser = argparse.ArgumentParser()
parser.add_argument("--job-id", required=True)
parser.add_argument("--input", required=True)
parser.add_argument("--output", required=True)
args = parser.parse_args()

spark = SparkSession.builder.appName(f"customer-revenue-{args.job_id}").getOrCreate()

orders = spark.read.option("header", True).csv(args.input)

report = (
    orders
    .groupBy("customer_id")
    .agg(
        spark_sum("total").alias("total_revenue"),
        count("*").alias("order_count"),
    )
)

report.write.mode("overwrite").parquet(args.output)

spark.stop()
```

---

## Read from JDBC

```python
orders = (
    spark.read.format("jdbc")
    .option("url", "jdbc:mysql://mysql:3306/app")
    .option("dbtable", "orders")
    .option("user", "readonly")
    .option("password", "secret")
    .load()
)
```

Use a read replica and restrict credentials.

---

## Write Result Manifest

Spark can write a small JSON manifest for Laravel:

```json
{
  "job_id": 10,
  "status": "completed",
  "output_path": "s3://analytics/reports/customer-revenue/job-10/",
  "row_count": 5000
}
```

Laravel can read this manifest to update the job record.

---

# Lesson 9: Reading Spark Results in Laravel

## Result Metadata

Spark should write output to a known path:

```text
s3://analytics/reports/customer-revenue/job-123/
```

Laravel stores:

```php
$sparkJob->update([
    'status' => 'completed',
    'output_path' => 'reports/customer-revenue/job-123/',
    'finished_at' => now(),
]);
```

---

## User Downloads

Laravel can generate a temporary signed URL:

```php
return Storage::disk('s3')->temporaryUrl(
    $sparkJob->output_path.'/part-00000.csv',
    now()->addMinutes(10)
);
```

For multi-file Spark outputs, create a ZIP, a single CSV, or expose a report viewer.

---

## Callback Endpoint

Spark can call Laravel at the end of a job:

```text
POST /internal/spark/jobs/{job}/complete
```

Payload:

```json
{
  "output_path": "s3://analytics/reports/customer-revenue/job-123/",
  "row_count": 5000
}
```

Secure this endpoint with a service token.

---

## Polling

If callbacks are not possible, Laravel can poll Spark or Livy:

```text
php artisan spark:poll-status
```

Run it every minute through Laravel Scheduler.

---

# Lesson 10: Scheduling

## Laravel Scheduler

Use Laravel Scheduler for simple periodic Spark jobs:

```php
Schedule::command('spark:run-daily-revenue')
    ->dailyAt('01:00')
    ->withoutOverlapping();
```

Good for:

- one-step reports
- simple daily jobs
- app-owned workflows

---

## Airflow

Use Airflow when workflows become multi-step:

```text
extract -> transform -> quality checks -> publish -> notify Laravel
```

Good for:

- dependencies
- retries per task
- backfills
- operational UI
- cross-system pipelines

---

## Spark Native Scheduling

Spark itself runs applications. It does not replace a business scheduler.

You still need Laravel Scheduler, Airflow, cron, Kubernetes CronJob, or another orchestrator to decide when jobs start.

---

## Rule of Thumb

| Need | Tool |
|---|---|
| Simple app-owned daily report | Laravel Scheduler |
| User triggers report | Laravel controller + queue |
| Multi-step data pipeline | Airflow |
| Continuous stream processing | Spark Structured Streaming |

---

# Lesson 11: Reliability

## Idempotent Output

Use deterministic output paths:

```text
s3://analytics/reports/customer-revenue/job-123/
```

On retry, Spark overwrites the same path:

```python
report.write.mode("overwrite").parquet(output_path)
```

This prevents duplicate result folders.

---

## Laravel Job State

Update state carefully:

```text
pending -> exported -> submitted -> running -> completed
```

If submission fails, mark `failed` with the error message. Do not leave jobs stuck forever in `pending`.

---

## Retries

Retry transient failures:

- temporary storage outage
- Spark cluster busy
- network timeout

Do not retry forever:

- invalid input path
- broken job code
- missing credentials

---

## Backfills

Spark is excellent for historical backfills:

```bash
php artisan spark:run-customer-revenue --from=2026-01-01 --to=2026-01-31
```

Use date ranges and deterministic job keys so backfills are repeatable.

---

# Lesson 12: Security

## Separate Credentials

Use separate credentials for:

- Laravel application database user
- Spark read-only database user
- object storage access
- Livy API user
- callback API token

Do not give Spark the same database permissions as Laravel.

---

## Laravel Internal Endpoints

Secure Spark callback endpoints:

```php
Route::middleware('auth.spark')->prefix('internal/spark')->group(function () {
    Route::post('/jobs/{job}/complete', [SparkJobCallbackController::class, 'complete']);
    Route::post('/jobs/{job}/failed', [SparkJobCallbackController::class, 'failed']);
});
```

Middleware:

```php
abort_unless(
    hash_equals(config('services.spark.token'), request()->bearerToken() ?? ''),
    403
);
```

---

## Secrets

Store secrets in:

- environment variables
- Docker secrets
- Kubernetes secrets
- cloud secret managers
- Airflow/Livy connection stores when applicable

Never commit passwords into Spark job files.

---

## Data Protection

For sensitive data:

- encrypt object storage
- restrict bucket access
- avoid exporting unnecessary columns
- mask personal data before analytics when possible
- define retention for intermediate files

---

# Lesson 13: Monitoring and Debugging

## IDs to Track

Every job should have:

| ID | Purpose |
|---|---|
| Laravel `spark_jobs.id` | Business tracking |
| Spark application ID | Spark UI and logs |
| Output path | Result location |
| Correlation ID | Cross-service tracing |

Store these together in Laravel.

---

## Laravel Logs

Log:

- job ID
- job type
- submitted command or Livy batch ID
- Spark application ID
- status changes
- output path
- error message

---

## Spark UI

Use Spark UI to inspect:

- stages
- tasks
- shuffle reads/writes
- executor memory
- failed tasks
- skewed partitions

The Spark UI usually explains performance problems better than Laravel logs can.

---

## Common Failures

| Symptom | Likely Cause |
|---|---|
| Laravel cannot submit | wrong path, permissions, Spark master unreachable |
| Job stuck submitted | no workers, cluster busy |
| Out of memory | bad partitioning, huge shuffle, too much data cached |
| Empty output | wrong date filter or input path |
| Permission denied | storage credentials or bucket policy |

---

# Lesson 14: Deployment

## Deployment Options

| Option | Good For |
|---|---|
| Spark standalone | Simple self-managed clusters |
| Spark on Kubernetes | Container-native deployments |
| AWS EMR | Managed Spark on AWS |
| Google Dataproc | Managed Spark on GCP |
| Azure Synapse/HDInsight | Managed Spark on Azure |

---

## Laravel Deployment

Laravel needs:

- queue workers for submit jobs
- scheduler for polling or recurring jobs
- access to Spark submit endpoint
- access to object storage
- secure internal callback routes

---

## Spark Job Packaging

Keep Spark jobs versioned:

```text
spark-jobs/
  customer_revenue.py
  product_affinity.py
  requirements.txt
```

Deploy Spark job files with the application or publish them to object storage.

---

## Production Checklist

- [ ] Laravel creates a job row before submitting.
- [ ] Spark input and output paths are deterministic.
- [ ] Spark application ID is stored.
- [ ] Users can see job status.
- [ ] Failed jobs store error messages.
- [ ] Outputs are protected with correct permissions.
- [ ] Large reads use replicas or data lake files.
- [ ] Secrets are not committed.
- [ ] Spark UI/logs are accessible to operators.
- [ ] Backfill process is documented.

---

# Lesson 15: Capstone Project

## Build: Customer Revenue Analytics

Create a Laravel + Spark workflow:

```text
1. User requests a customer revenue report.
2. Laravel creates a spark_jobs row.
3. Laravel dispatches SubmitSparkJob.
4. Spark reads orders for the date range.
5. Spark groups revenue by customer.
6. Spark writes Parquet or CSV output.
7. Laravel marks the job completed.
8. User downloads the report from Laravel.
```

---

## Laravel Pieces

- `spark_jobs` table
- `POST /reports/customer-revenue`
- `SparkJob` model
- `SubmitSparkJob` queue job
- `SparkSubmitter` or `LivyClient`
- status page
- download route
- callback route or polling command

---

## Spark Pieces

- `customer_revenue.py`
- input argument parsing
- date filtering
- aggregation by customer
- deterministic output path
- manifest file

---

## Acceptance Criteria

- Web request returns quickly.
- Spark job runs outside the request.
- Laravel stores job status.
- Spark output path is deterministic.
- Retrying does not create duplicate reports.
- User can see failed/completed status.
- User can download the final report.
