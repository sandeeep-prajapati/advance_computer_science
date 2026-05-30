# Lesson 7: Airflow Calling Laravel

## Why Airflow Calls Laravel

Airflow may need data or commands from Laravel:

- get pending records
- mark report as processing
- fetch signed export URL
- notify Laravel that a workflow finished

Prefer HTTP APIs for cross-service calls.

---

## Laravel Internal API

```php
Route::middleware('auth.airflow')->prefix('internal/airflow')->group(function () {
    Route::get('/reports/{report}', [AirflowReportController::class, 'show']);
    Route::post('/reports/{report}/complete', [AirflowReportController::class, 'complete']);
});
```

Use a dedicated token or service account.

---

## Airflow SimpleHttpOperator

```python
from airflow.providers.http.operators.http import SimpleHttpOperator

mark_complete = SimpleHttpOperator(
    task_id="mark_report_complete",
    http_conn_id="laravel_api",
    endpoint="/internal/airflow/reports/{{ dag_run.conf['report_id'] }}/complete",
    method="POST",
    headers={"Authorization": "Bearer {{ var.value.LARAVEL_AIRFLOW_TOKEN }}"},
    json={
        "file_path": "{{ ti.xcom_pull(task_ids='build_report') }}",
    },
)
```

---

## Calling Artisan from Airflow

If Airflow and Laravel share infrastructure, Airflow can run Artisan commands:

```python
from airflow.operators.bash import BashOperator

export_orders = BashOperator(
    task_id="export_orders",
    bash_command="cd /var/www/html && php artisan reports:export-orders --report-id {{ dag_run.conf['report_id'] }}",
)
```

Use this carefully. HTTP APIs are usually cleaner across services.
