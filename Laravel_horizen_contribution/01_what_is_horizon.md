# Chapter 1: What Horizon Is

## The Simple Explanation

Laravel Horizon is a queue worker manager and dashboard for Laravel applications that use Redis queues.

It helps you:

- run queue workers
- monitor queued jobs
- see throughput and runtime
- inspect failed jobs
- pause and continue queues
- manage supervisors
- balance workers across queues
- monitor job tags and metrics

---

## Mental Model

```text
Laravel App
├── dispatches jobs
│
Redis Queue
├── stores pending jobs
├── stores reserved jobs
├── stores failed job metadata
└── stores Horizon metrics
│
Horizon Master Supervisor
├── starts supervisors
├── monitors worker processes
├── scales workers
└── handles termination/restart
│
Horizon Dashboard
├── displays jobs
├── displays queues
├── displays supervisors
├── displays metrics
└── controls pause/continue
```

---

## What Horizon Adds Over Normal Queue Workers

Laravel already has:

```bash
php artisan queue:work
```

Horizon adds:

```bash
php artisan horizon
```

With Horizon you get:

- a dashboard
- supervisor configuration
- process balancing
- queue metrics
- job metrics
- tags
- failed job visibility
- production-friendly worker management

---

## Main Runtime Pieces

| Piece | Responsibility |
|---|---|
| Service provider | registers Horizon into Laravel |
| Artisan commands | `horizon`, `horizon:terminate`, `horizon:pause`, etc. |
| Master supervisor | manages supervisor lifecycle |
| Supervisor | manages process pools |
| Process pool | manages worker processes |
| Auto scaler | decides process counts |
| Repositories | read/write Horizon state in Redis |
| Dashboard assets | Vue/JS/CSS interface |

---

## Why Contributing Is Valuable

Horizon sits at a serious production boundary:

```text
background jobs
Redis
worker processes
process balancing
dashboard observability
deployment restarts
failed job handling
```

Small bugs can affect production workloads. Good tests and clear fixes matter.

---

## Important Warning

Do not treat Horizon as only a UI package.

The dashboard is only one part. The harder part is process orchestration:

- workers
- supervisors
- scaling
- Redis state
- restarts
- queue balancing
- metrics accuracy

