# Chapter 10: Sources & Tracking

## Official Sources

Use these as your source of truth:

- Laravel Horizon repository: `https://github.com/laravel/horizon`
- Laravel contribution guide: `https://laravel.com/docs/contributions`
- Laravel queues documentation: `https://laravel.com/docs/queues`
- Laravel Horizon documentation: `https://laravel.com/docs/horizon`

Always check the branch before contributing. Horizon currently uses the `5.x` branch in the public repository context checked for this guide.

---

## What I Verified for This Guide

Current public context checked:

- repository: `laravel/horizon`
- active branch context: `5.x`
- open PR example around autoscaling and terminating process counts
- Horizon is a Laravel Redis queue dashboard/worker supervisor
- contribution areas include backend process balancing, Redis repositories, failed jobs, dashboard UI, tests, and docs

---

## Personal Tracking Template

Use this when choosing work:

```md
## Candidate Issue / PR

Link:

Area:

Difficulty:

What is broken:

How to reproduce:

Files likely involved:

Tests likely needed:

Risk:

Decision:
```

---

## Learning Checklist

Before first PR:

- [ ] I can explain what Horizon does.
- [ ] I can run the test suite.
- [ ] I can find service provider, commands, supervisor code, and tests.
- [ ] I understand Laravel queues basics.
- [ ] I understand Redis basics used by Horizon.
- [ ] I can write a failing test for a small bug.
- [ ] I can explain my PR in less than 10 lines.

---

## Your Contribution Ladder

```text
Level 1: docs and small tests
Level 2: dashboard UI states
Level 3: failed jobs and Redis repository tests
Level 4: dashboard API accuracy
Level 5: autoscaling regression tests
Level 6: autoscaling/process lifecycle fixes
```

Move one level at a time.

