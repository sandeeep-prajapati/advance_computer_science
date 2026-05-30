# Lesson 11: Reliability

## Idempotent Output

Use deterministic output paths:

```text
s3://analytics/reports/customer-revenue/job-123/
```

On retry, Spark overwrites the same path:

```python
report.write.mode("overwrite").parquet(output_path)
```

This prevents duplicate result folders.

---

## Laravel Job State

Update state carefully:

```text
pending -> exported -> submitted -> running -> completed
```

If submission fails, mark `failed` with the error message. Do not leave jobs stuck forever in `pending`.

---

## Retries

Retry transient failures:

- temporary storage outage
- Spark cluster busy
- network timeout

Do not retry forever:

- invalid input path
- broken job code
- missing credentials

---

## Backfills

Spark is excellent for historical backfills:

```bash
php artisan spark:run-customer-revenue --from=2026-01-01 --to=2026-01-31
```

Use date ranges and deterministic job keys so backfills are repeatable.
