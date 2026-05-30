# Lesson 2: Architecture

## Recommended Architecture

```text
                 ┌────────────────────┐
User request --> │ Laravel App         │
                 │ - creates job row   │
                 │ - validates input   │
                 └─────────┬──────────┘
                           │
                           │ submit job
                           v
                 ┌────────────────────┐
                 │ Spark Cluster       │
                 │ - PySpark/Scala job │
                 └─────────┬──────────┘
                           │
                           │ reads/writes
                           v
                 ┌────────────────────┐
                 │ Storage             │
                 │ S3/MinIO/HDFS/DB    │
                 └─────────┬──────────┘
                           │
                           │ result metadata
                           v
                 ┌────────────────────┐
                 │ Laravel Dashboard   │
                 │ status + download   │
                 └────────────────────┘
```

---

## Job Table

Laravel should track Spark work in a database table:

```php
Schema::create('spark_jobs', function (Blueprint $table) {
    $table->id();
    $table->string('type');
    $table->string('status')->default('pending');
    $table->json('parameters');
    $table->string('spark_application_id')->nullable();
    $table->string('output_path')->nullable();
    $table->text('error_message')->nullable();
    $table->timestamp('started_at')->nullable();
    $table->timestamp('finished_at')->nullable();
    $table->timestamps();
});
```

---

## Status Flow

```text
pending -> submitted -> running -> completed
                         |
                         v
                       failed
```

The user should never wait in a browser request while Spark processes data.
