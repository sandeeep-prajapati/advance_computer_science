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
