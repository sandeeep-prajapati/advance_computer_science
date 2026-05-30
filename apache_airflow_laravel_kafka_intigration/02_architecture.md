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
