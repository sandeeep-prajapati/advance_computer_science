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
