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
