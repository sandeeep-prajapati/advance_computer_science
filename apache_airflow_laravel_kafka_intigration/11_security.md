# Lesson 11: Security

## Service Authentication

Use a dedicated Airflow service token for Laravel internal APIs:

```dotenv
AIRFLOW_SERVICE_TOKEN=...
```

Laravel middleware:

```php
public function handle($request, Closure $next)
{
    abort_unless(
        hash_equals(config('services.airflow.token'), $request->bearerToken() ?? ''),
        403
    );

    return $next($request);
}
```

---

## Airflow Secrets

Store secrets in:

- Airflow Connections
- Airflow Variables for non-sensitive config
- environment variables
- cloud secret manager

Do not hardcode Laravel tokens or Kafka passwords in DAG files.

---

## Kafka Security

Production Kafka should use:

```text
SASL_SSL
ACLs per service
separate producer and consumer credentials
topic-level permissions
consumer-group permissions
```

Laravel producer permission:

```text
WRITE orders.created
DESCRIBE orders.created
```

Laravel consumer permission:

```text
READ reports.generated
READ group laravel-report-updater
```

Airflow producer permission:

```text
WRITE reports.generated
DESCRIBE reports.generated
```
