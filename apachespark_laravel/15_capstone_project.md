# Lesson 15: Capstone Project

## Build: Customer Revenue Analytics

Create a Laravel + Spark workflow:

```text
1. User requests a customer revenue report.
2. Laravel creates a spark_jobs row.
3. Laravel dispatches SubmitSparkJob.
4. Spark reads orders for the date range.
5. Spark groups revenue by customer.
6. Spark writes Parquet or CSV output.
7. Laravel marks the job completed.
8. User downloads the report from Laravel.
```

---

## Laravel Pieces

- `spark_jobs` table
- `POST /reports/customer-revenue`
- `SparkJob` model
- `SubmitSparkJob` queue job
- `SparkSubmitter` or `LivyClient`
- status page
- download route
- callback route or polling command

---

## Spark Pieces

- `customer_revenue.py`
- input argument parsing
- date filtering
- aggregation by customer
- deterministic output path
- manifest file

---

## Acceptance Criteria

- Web request returns quickly.
- Spark job runs outside the request.
- Laravel stores job status.
- Spark output path is deterministic.
- Retrying does not create duplicate reports.
- User can see failed/completed status.
- User can download the final report.
