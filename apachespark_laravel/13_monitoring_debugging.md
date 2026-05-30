# Lesson 13: Monitoring and Debugging

## IDs to Track

Every job should have:

| ID | Purpose |
|---|---|
| Laravel `spark_jobs.id` | Business tracking |
| Spark application ID | Spark UI and logs |
| Output path | Result location |
| Correlation ID | Cross-service tracing |

Store these together in Laravel.

---

## Laravel Logs

Log:

- job ID
- job type
- submitted command or Livy batch ID
- Spark application ID
- status changes
- output path
- error message

---

## Spark UI

Use Spark UI to inspect:

- stages
- tasks
- shuffle reads/writes
- executor memory
- failed tasks
- skewed partitions

The Spark UI usually explains performance problems better than Laravel logs can.

---

## Common Failures

| Symptom | Likely Cause |
|---|---|
| Laravel cannot submit | wrong path, permissions, Spark master unreachable |
| Job stuck submitted | no workers, cluster busy |
| Out of memory | bad partitioning, huge shuffle, too much data cached |
| Empty output | wrong date filter or input path |
| Permission denied | storage credentials or bucket policy |
