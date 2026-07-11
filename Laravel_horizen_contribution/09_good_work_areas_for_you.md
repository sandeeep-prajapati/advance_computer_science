# Chapter 9: Where You Can Work

## Beginner-Friendly Areas

Start here if you are new to open source.

### 1. Tests Around Existing Behavior

Why:

- low risk
- teaches internals
- maintainers like stronger coverage

Work examples:

- command tests
- repository edge cases
- dashboard API response tests
- failed job tests

### 2. Dashboard UI States

Why:

- visible
- small changes possible
- easier to manually verify

Work examples:

- empty state text
- loading state
- error display
- button disabled states
- small responsive fixes

### 3. Documentation and Examples

Why:

- easiest entry
- forces you to learn behavior

Work examples:

- clarify config option explanation
- add example for supervisor balancing
- improve command description

---

## Intermediate Areas

Move here after you understand tests.

### 1. Redis Repository Edge Cases

Work examples:

- stale data cleanup
- missing key behavior
- trimming old metrics
- tag lookup behavior

### 2. Failed Job Flow

Work examples:

- failed job display
- failed job retry/clear behavior
- failed job payload parsing
- better tests around failed jobs

### 3. Dashboard API Accuracy

Work examples:

- ensure JSON response matches Redis state
- fix wrong counts
- improve polling error handling

---

## Advanced Areas

These are high-value but harder.

### 1. Autoscaling and Balancing

Work examples:

- process count edge cases
- terminating process handling
- long-running job scenarios
- min/max process constraints
- multi-queue balancing

This is a strong area if you want to become serious in Horizon.

### 2. Process Lifecycle

Work examples:

- graceful shutdown
- memory limit behavior
- worker restart behavior
- deploy termination behavior

### 3. Performance

Work examples:

- dashboard query performance
- Redis lookup optimization
- reducing unnecessary polling cost
- large queue display behavior

---

## Recommended Area for You

Based on your interest in queues, RabbitMQ wrappers, framework internals, Laravel, Pimcore, and backend architecture, your best path is:

```text
1. Tests for queue/process behavior
2. Failed job and Redis repository edge cases
3. Dashboard API accuracy
4. Autoscaling/balancing regression tests
5. Autoscaling bug fixes
```

This path builds from safe to advanced.

---

## Concrete First Issue Hunt

Search issues/PRs for:

```text
autoscale
balancing
terminating process
failed jobs
Redis
metrics
dashboard count
pause supervisor
queue wait
memory limit
```

For each candidate:

```text
Can I reproduce?
Can I write a test?
Is the fix small?
Is it still relevant on 5.x?
```

If yes, it may be a good contribution.

