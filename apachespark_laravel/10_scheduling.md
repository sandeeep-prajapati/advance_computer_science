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
