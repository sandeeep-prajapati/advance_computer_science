# Lesson 6: Triggering Airflow from Laravel

## Use Airflow REST API

Laravel can trigger a DAG run when a user starts a workflow.

Endpoint shape:

```text
POST /api/v1/dags/{dag_id}/dagRuns
```

Example Laravel service:

```php
use Illuminate\Support\Facades\Http;

final class AirflowClient
{
    public function triggerDag(string $dagId, array $conf = []): array
    {
        $response = Http::withBasicAuth(
                config('services.airflow.username'),
                config('services.airflow.password')
            )
            ->baseUrl(config('services.airflow.base_url'))
            ->post("/api/v1/dags/{$dagId}/dagRuns", [
                'conf' => $conf,
            ]);

        $response->throw();

        return $response->json();
    }
}
```

---

## Trigger with Business Context

```php
app(AirflowClient::class)->triggerDag('generate_order_report', [
    'report_id' => $report->id,
    'from' => $request->date('from')->toDateString(),
    'to' => $request->date('to')->toDateString(),
    'requested_by' => auth()->id(),
    'correlation_id' => (string) str()->uuid(),
]);
```

Airflow reads this as `dag_run.conf`.

---

## Airflow DAG Reads Conf

```python
from airflow.decorators import dag, task
from pendulum import datetime

@dag(start_date=datetime(2026, 1, 1), schedule=None, catchup=False)
def generate_order_report():
    @task
    def read_conf(**context):
        conf = context["dag_run"].conf or {}
        return {
            "report_id": conf["report_id"],
            "from": conf["from"],
            "to": conf["to"],
            "correlation_id": conf["correlation_id"],
        }

    read_conf()

generate_order_report()
```

---

## Laravel Status Model

Store the Airflow DAG run ID in Laravel:

```php
$dagRun = app(AirflowClient::class)->triggerDag(...);

$report->update([
    'status' => 'processing',
    'airflow_dag_run_id' => $dagRun['dag_run_id'],
]);
```
