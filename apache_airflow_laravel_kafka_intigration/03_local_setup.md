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
