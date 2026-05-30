# Lesson 4: Data Exchange Patterns

## Pattern 1: Spark Reads Laravel Database

Spark can read MySQL/Postgres through JDBC.

Pros:

- simple source of truth
- no export step

Cons:

- can load production database heavily
- needs network and credentials
- large reads may affect app performance

Use read replicas when possible.

---

## Pattern 2: Laravel Exports Files

Laravel writes CSV or JSON files for Spark:

```text
storage/app/exports/orders-2026-05-30.csv
```

Spark reads the file and writes transformed output.

Pros:

- protects app database
- easy to replay
- good for medium-sized exports

Cons:

- export can be slow in PHP for huge datasets
- CSV needs careful typing

---

## Pattern 3: Object Storage

Best production pattern:

```text
Laravel -> S3/MinIO raw data
Spark -> S3/MinIO processed data
Laravel -> reads metadata/result path
```

Prefer Parquet for Spark outputs:

- compressed
- columnar
- typed
- efficient for analytics

---

## Pattern 4: API Calls

Spark can call Laravel APIs, but this is not ideal for large data.

Use APIs for:

- metadata
- small lookup tables
- job status callbacks

Avoid APIs for millions of rows.
