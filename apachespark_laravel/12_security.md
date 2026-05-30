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
