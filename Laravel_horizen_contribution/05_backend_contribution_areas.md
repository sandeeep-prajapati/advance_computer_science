# Chapter 5: Backend Contribution Areas

## Best Backend Areas to Target

Start with one of these:

| Area | Difficulty | Why Good |
|---|---|---|
| tests for existing behavior | beginner | teaches repo safely |
| command output edge cases | beginner | small surface area |
| failed job cleanup | intermediate | useful and testable |
| Redis repository edge cases | intermediate | important backend layer |
| dashboard API response bugs | intermediate | connects backend/frontend |
| autoscaling/balancing | advanced | high impact, tricky |
| process lifecycle | advanced | production critical |

---

## Area 1: Regression Tests

Look for old bug reports or open PRs and write tests proving the behavior.

Good test PR:

```text
Given queue A has long-running jobs
And queue B receives many short jobs
When autoscaler runs
Then terminating processes should not block scaling decisions
```

Why this is good:

- maintainers value tests
- small PRs are easier to review
- you learn internals safely

---

## Area 2: Autoscaling and Balancing

This is a meaningful advanced area.

Current open PR context shows autoscaling edge cases around terminating processes and worker counts. This is a real example of the kind of bug you can study.

What to learn:

- `AutoScaler`
- `ProcessPool`
- supervisor options
- min/max process counts
- terminating process handling
- tests for scaling decisions

Possible work:

- add missing regression tests
- improve naming around process counts
- document balance behavior
- test edge cases with long-running jobs

---

## Area 3: Supervisor Lifecycle

Useful questions:

- What happens on `horizon:terminate`?
- What happens when a worker hits memory limit?
- What happens when a process is terminating but still running a job?
- What happens during deploy?

Possible work:

- graceful shutdown tests
- clearer command output
- edge case bug fixes
- documentation examples

---

## Area 4: Redis Repository Behavior

Horizon stores state in Redis.

Look for:

- stale keys
- trimming bugs
- missing expiration
- inconsistent counts
- inefficient lookups
- edge cases with empty data

Good first backend task:

```text
Find repository method with weak test coverage.
Add tests for empty Redis state, missing keys, and malformed payloads.
```

---

## Area 5: Failed Jobs

Failed-job flows are practical and visible.

Possible work:

- failed job filtering
- retry behavior
- clearing failed jobs
- display accuracy
- tests for failed job payload shape

---

## Area 6: Compatibility

Laravel packages often need updates for:

- new Laravel versions
- new Symfony components
- new PHP versions
- Redis client differences
- frontend build tooling changes

This is useful, but be careful. Compatibility PRs must be precise and tested.

