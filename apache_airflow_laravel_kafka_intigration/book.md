# Apache Airflow + Laravel + Kafka Integration — Lessons

> Practical lessons for orchestrating Laravel/PHP workflows with Apache Airflow and Apache Kafka.

---

## Table of Contents

| Lesson | Title | Topics |
|---|---|---|
| 01 | [Introduction](01_introduction.md) | What Airflow should do, what Laravel should do, where Kafka fits |
| 02 | [Architecture](02_architecture.md) | Event-driven batch orchestration, services, topics, DAGs |
| 03 | [Local Setup](03_local_setup.md) | Docker layout, Laravel app, Kafka, Airflow, networking |
| 04 | [Laravel Kafka Producer](04_laravel_kafka_producer.md) | Publishing events from Laravel |
| 05 | [Laravel Kafka Consumer](05_laravel_kafka_consumer.md) | Consuming Kafka events in Laravel |
| 06 | [Triggering Airflow from Laravel](06_trigger_airflow_from_laravel.md) | REST API, auth, `dag_run.conf`, PHP HTTP client |
| 07 | [Airflow Calling Laravel](07_airflow_calling_laravel.md) | `HttpOperator`, internal APIs, Artisan commands |
| 08 | [Airflow and Kafka](08_airflow_and_kafka.md) | Producing/consuming Kafka from DAG tasks |
| 09 | [Data Pipeline Pattern](09_data_pipeline_pattern.md) | Export from Laravel, process in Airflow, publish result |
| 10 | [Reliability](10_reliability.md) | Idempotency, retries, DLQ, backfills, exactly-once thinking |
| 11 | [Security](11_security.md) | API tokens, network boundaries, SASL/TLS, secrets |
| 12 | [Monitoring](12_monitoring.md) | Airflow logs, Laravel logs, Kafka lag, correlation IDs |
| 13 | [Deployment](13_deployment.md) | Docker, Supervisor, Kubernetes, production checklist |
| 14 | [Capstone Project](14_capstone_project.md) | Order analytics pipeline |

---

## Integration Rule

Use each tool for its natural job:

- **Laravel** owns web requests, business rules, database writes, and user-facing APIs.
- **Kafka** carries durable events between systems.
- **Airflow** orchestrates scheduled, retryable, observable workflows.

Do not turn Airflow into your web app, and do not turn Laravel queue workers into a full workflow scheduler.

---

# Lesson 1: Introduction

## Goal

This guide teaches how to integrate Apache Airflow with a PHP/Laravel application that uses Kafka.

The common pattern is:

```text
Laravel publishes business events
Kafka stores and distributes those events
Airflow runs scheduled or triggered workflows
Laravel consumes results or receives callbacks
```

---

## When This Integration Makes Sense

Use Airflow with Laravel and Kafka when you need:

- scheduled exports from a Laravel database
- nightly reports or analytics pipelines
- long-running data processing outside HTTP requests
- retryable workflows with dependencies
- backfills for historical data
- coordination between Laravel, Kafka, storage, and data warehouses

---

## When Not to Use Airflow

Airflow is not a replacement for:

- Laravel controllers
- Laravel queues
- real-time request handling
- simple cron jobs
- short background jobs

If a task is small and internal to Laravel, use Laravel queues or Laravel Scheduler. Use Airflow when workflows are multi-step, observable, retryable, and operationally important.

---

## Main Integration Styles

| Style | Direction | Example |
|---|---|---|
| Laravel triggers Airflow | Laravel -> Airflow API | User uploads CSV and starts processing DAG |
| Airflow calls Laravel | Airflow -> Laravel API | DAG asks Laravel for pending orders |
| Laravel publishes Kafka | Laravel -> Kafka | `orders.created` event |
| Airflow reads Kafka | Kafka -> Airflow | DAG processes event batches |
| Airflow publishes Kafka | Airflow -> Kafka | `reports.generated` event |
| Laravel consumes Kafka | Kafka -> Laravel | Laravel updates report status |

---

# Lesson 2: Architecture

## Recommended Architecture

```text
                 ┌──────────────────┐
HTTP Requests -> │ Laravel App       │
                 │ - business logic  │
                 │ - MySQL/Postgres  │
                 └───────┬──────────┘
                         │ publish events
                         v
                 ┌──────────────────┐
                 │ Kafka             │
                 │ orders.created    │
                 │ reports.ready     │
                 └───────┬──────────┘
                         │
          ┌──────────────┴──────────────┐
          v                             v
┌──────────────────┐           ┌──────────────────┐
│ Laravel Consumer │           │ Airflow DAGs      │
│ updates app DB   │           │ batch workflows   │
└──────────────────┘           └──────────────────┘
```

---

## Clear Ownership

Laravel should own:

- users
- permissions
- business entities
- API validation
- transactional writes
- user-facing status pages

Airflow should own:

- schedules
- workflow dependencies
- retries
- backfills
- data movement
- operational visibility

Kafka should own:

- event storage
- fan-out
- replay
- decoupling between Laravel and Airflow

---

## Example Workflow

```text
1. Laravel creates an order.
2. Laravel publishes orders.created to Kafka.
3. Airflow scheduled DAG reads new order data every hour.
4. Airflow transforms data into a report.
5. Airflow publishes reports.generated to Kafka.
6. Laravel consumer updates report status in the Laravel database.
```

This keeps the user-facing app fast while heavy workflows run outside Laravel.

---

# Lesson 3: Local Setup

## Services

A local development stack usually has:

```text
laravel-app
mysql or postgres
kafka
airflow-webserver
airflow-scheduler
airflow-worker
```

Keep all services on one Docker network so they can call each other by service name.

---

## Environment Variables

Laravel:

```dotenv
KAFKA_BROKERS=kafka:9092
AIRFLOW_BASE_URL=http://airflow-webserver:8080
AIRFLOW_USERNAME=airflow
AIRFLOW_PASSWORD=airflow
```

Airflow:

```dotenv
LARAVEL_BASE_URL=http://laravel-app
KAFKA_BOOTSTRAP_SERVERS=kafka:9092
```

---

## Network Rule

Inside Docker, do not use `localhost` to reach another container.

Use:

```text
kafka:9092
airflow-webserver:8080
laravel-app:80
```

`localhost` means "this same container", not your whole Docker project.

---

## First Smoke Tests

From Laravel:

```bash
php artisan route:list
php artisan kafka:publish-test
```

From Airflow container:

```bash
airflow dags list
curl http://laravel-app/api/health
```

From Kafka tooling:

```bash
kafka-topics.sh --bootstrap-server kafka:9092 --list
```

---

# Lesson 4: Laravel Kafka Producer

## Install Laravel Kafka Package

```bash
composer require mateusjunges/laravel-kafka
```

Make sure PHP has `rdkafka` installed and enabled.

---

## Publish an Event

```php
use Junges\Kafka\Facades\Kafka;
use Junges\Kafka\Message\Message;

Kafka::publish(config('kafka.brokers'))
    ->onTopic('orders.created')
    ->withMessage(new Message(
        headers: [
            'event' => 'orders.created',
            'event_version' => '1',
            'correlation_id' => (string) str()->uuid(),
        ],
        body: [
            'order_id' => $order->id,
            'user_id' => $order->user_id,
            'total' => $order->total,
            'created_at' => $order->created_at->toIso8601String(),
        ],
        key: "order-{$order->id}",
    ))
    ->send();
```

---

## Publish After Commit

If the event depends on database data, publish after the database transaction commits.

```php
final class PublishOrderCreatedToKafka implements ShouldQueue
{
    public bool $afterCommit = true;

    public function handle(OrderCreated $event): void
    {
        app(OrderKafkaPublisher::class)->publishCreated($event->orderId);
    }
}
```

For critical workflows, use an outbox table and a dedicated publisher worker.

---

## Topic Design

Good topic names:

```text
orders.created
orders.cancelled
reports.generated
airflow.pipeline.failed
```

Use stable keys when order matters:

```text
order-123
report-2026-05-30
customer-50
```

---

# Lesson 5: Laravel Kafka Consumer

## Consume Airflow Results

Airflow can publish workflow result events to Kafka. Laravel consumes them and updates application state.

```php
use Junges\Kafka\Facades\Kafka;
use Junges\Kafka\Contracts\ConsumerMessage;
use Junges\Kafka\Contracts\MessageConsumer;

Kafka::consumer(['reports.generated'], 'laravel-report-updater', config('kafka.brokers'))
    ->withHandler(function (ConsumerMessage $message, MessageConsumer $consumer): void {
        $payload = $message->getBody();

        Report::where('id', $payload['report_id'])->update([
            'status' => 'ready',
            'file_path' => $payload['file_path'],
            'generated_at' => $payload['generated_at'],
        ]);
    })
    ->build()
    ->consume();
```

---

## Run as Artisan Command

```php
final class ConsumeReportsCommand extends Command
{
    protected $signature = 'kafka:consume-reports';

    public function handle(): int
    {
        Kafka::consumer(['reports.generated'], 'laravel-report-updater', config('kafka.brokers'))
            ->withHandler(app(ReportGeneratedHandler::class))
            ->build()
            ->consume();

        return self::SUCCESS;
    }
}
```

Run it with Supervisor, systemd, Docker, or Kubernetes.

---

## Idempotency

Kafka consumers may receive duplicates. Store processed event IDs:

```php
ProcessedKafkaMessage::firstOrCreate([
    'event_id' => $payload['event_id'],
    'consumer_group' => 'laravel-report-updater',
]);
```

Make database updates safe to repeat.

---

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

---

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

---

# Lesson 8: Airflow and Kafka

## Airflow Publishing to Kafka

Install a Kafka Python client in the Airflow image, such as `confluent-kafka`.

```python
from confluent_kafka import Producer
import json
import os

def publish_report_generated(report_id: int, file_path: str):
    producer = Producer({
        "bootstrap.servers": os.environ["KAFKA_BOOTSTRAP_SERVERS"],
        "acks": "all",
    })

    producer.produce(
        "reports.generated",
        key=f"report-{report_id}",
        value=json.dumps({
            "event_type": "reports.generated",
            "event_version": 1,
            "report_id": report_id,
            "file_path": file_path,
        }),
    )
    producer.flush()
```

---

## Airflow Task

```python
from airflow.decorators import task

@task
def notify_laravel(report_id: int, file_path: str):
    publish_report_generated(report_id, file_path)
```

Laravel then consumes `reports.generated`.

---

## Should Airflow Consume Kafka Continuously?

Airflow is best for scheduled or finite workflows. A never-ending Kafka consumer is usually better as:

- Laravel consumer command
- standalone Python service
- Kafka Connect
- Flink or Spark streaming job

Airflow can read a bounded batch from Kafka, but it should not normally run as a permanent stream processor.

---

## Bounded Kafka Batch Pattern

```text
Airflow DAG every 15 minutes
  -> read offsets from tracking table
  -> consume up to N messages
  -> process batch
  -> save latest offsets
```

Use this for batch analytics, not low-latency event processing.

---

# Lesson 9: Data Pipeline Pattern

## Scenario

Laravel stores orders. Airflow generates daily revenue reports. Kafka tells Laravel when reports are ready.

```text
Laravel DB -> Airflow DAG -> CSV/Parquet/S3 -> Kafka event -> Laravel UI status
```

---

## DAG Shape

```text
extract_orders
  -> transform_orders
  -> build_report
  -> publish_reports_generated
```

---

## Extract from Laravel API

```python
import requests

def fetch_orders(base_url, token, start_date, end_date):
    response = requests.get(
        f"{base_url}/internal/airflow/orders",
        headers={"Authorization": f"Bearer {token}"},
        params={"from": start_date, "to": end_date},
        timeout=60,
    )
    response.raise_for_status()
    return response.json()
```

---

## Publish Result Event

```json
{
  "event_id": "01J...",
  "event_type": "reports.generated",
  "event_version": 1,
  "report_id": 55,
  "file_path": "s3://bucket/reports/orders-2026-05-30.csv",
  "generated_at": "2026-05-30T12:00:00Z"
}
```

Laravel consumes this event and updates the report row.

---

## Why This Pattern Works

- Laravel remains the source of business truth.
- Airflow owns heavy data workflow execution.
- Kafka decouples report generation from Laravel status updates.
- The UI can show progress without running the workflow in a web request.

---

# Lesson 10: Reliability

## Idempotency Everywhere

Airflow tasks can retry. Kafka messages can be duplicated. HTTP calls can time out after succeeding.

Design every step to be safe to repeat.

Laravel examples:

```php
Report::updateOrCreate(
    ['id' => $payload['report_id']],
    ['status' => 'ready', 'file_path' => $payload['file_path']]
);
```

Airflow examples:

```python
# write output to deterministic path
file_path = f"s3://reports/orders/{report_id}.csv"
```

---

## Retries

Airflow retry config:

```python
@dag(
    default_args={
        "retries": 3,
        "retry_delay": timedelta(minutes=5),
    },
)
def report_pipeline():
    ...
```

Laravel Kafka consumer retries should have limits and DLQ behavior.

---

## Dead Letter Topics

Use DLQ topics for messages Laravel cannot process:

```text
reports.generated.dlq
orders.created.dlq
```

Include:

- original topic
- original key
- payload
- error message
- failed timestamp

---

## Backfills

Airflow is excellent for backfills:

```text
Run daily_report for 2026-05-01 through 2026-05-30
```

Make sure Laravel APIs can accept date ranges and deterministic report IDs so backfills do not create duplicate business records.

---

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

---

# Lesson 12: Monitoring

## Correlation IDs

Generate one correlation ID in Laravel and pass it through:

```text
Laravel request
  -> Kafka headers
  -> Airflow dag_run.conf
  -> Airflow logs
  -> reports.generated event
  -> Laravel consumer logs
```

This makes one workflow traceable across all systems.

---

## Laravel Logs

Log:

- event type
- event ID
- topic
- key
- correlation ID
- Airflow DAG ID
- Airflow DAG run ID

---

## Airflow Monitoring

Watch:

- DAG failures
- task duration
- retry count
- queued task count
- scheduler health
- worker health

Airflow should alert when a business-critical DAG fails.

---

## Kafka Monitoring

Watch:

- consumer lag
- broker health
- produce failures
- DLQ message count
- topic size and retention

Consumer lag connects directly to business freshness. If lag grows, users may see stale report status or delayed updates.

---

# Lesson 13: Deployment

## Process Layout

Production commonly has:

```text
Laravel web containers
Laravel queue workers
Laravel Kafka consumer workers
Airflow webserver
Airflow scheduler
Airflow workers
Kafka or managed Kafka
Database
Object storage
```

---

## Laravel Consumers

Run Kafka consumers as managed long-running processes:

```ini
[program:kafka-report-consumer]
command=php /var/www/html/artisan kafka:consume-reports
autostart=true
autorestart=true
stopwaitsecs=30
```

Restart consumers during deploys so they load new code.

---

## Airflow DAG Deployment

Keep DAG files in source control. Deploy them the same way every environment receives application code.

Use environment-specific config through:

- Airflow Connections
- Airflow Variables
- environment variables
- secret manager

---

## Production Checklist

- [ ] Laravel can reach Airflow API.
- [ ] Airflow can reach Laravel internal API.
- [ ] Laravel and Airflow can reach Kafka.
- [ ] Kafka uses SASL/TLS in production.
- [ ] Internal APIs require service authentication.
- [ ] Consumers are idempotent.
- [ ] Airflow tasks are idempotent.
- [ ] DAG failures alert the team.
- [ ] Kafka lag is monitored.
- [ ] DLQ replay process is documented.

---

# Lesson 14: Capstone Project

## Build: Order Analytics Pipeline

Create a Laravel + Kafka + Airflow workflow:

```text
1. User requests an order report in Laravel.
2. Laravel creates report row with status=pending.
3. Laravel triggers Airflow DAG with report_id and date range.
4. Airflow calls Laravel internal API to fetch orders.
5. Airflow builds CSV report.
6. Airflow publishes reports.generated to Kafka.
7. Laravel consumes reports.generated.
8. Laravel marks report status=ready.
```

---

## Laravel Pieces

- `reports` table
- `POST /reports`
- `GET /internal/airflow/orders`
- `POST /internal/airflow/reports/{report}/failed`
- `AirflowClient`
- `kafka:consume-reports`
- `ReportGeneratedHandler`

---

## Airflow Pieces

- DAG `generate_order_report`
- task `read_conf`
- task `fetch_orders_from_laravel`
- task `build_csv`
- task `publish_reports_generated`
- failure callback that notifies Laravel

---

## Kafka Topics

```text
reports.generated
reports.generated.dlq
airflow.pipeline.failed
```

---

## Acceptance Criteria

- Laravel can trigger a DAG.
- Airflow can read `dag_run.conf`.
- Airflow can call Laravel internal API.
- Airflow publishes a Kafka result event.
- Laravel consumes the result and updates report status.
- Duplicate Kafka messages do not duplicate report rows.
- Failed workflows are visible in Laravel and Airflow.
