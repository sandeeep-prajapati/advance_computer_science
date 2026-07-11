# Chapter 3: Architecture Map

## Common Repo Areas

Laravel Horizon is a package, so expect a structure like:

```text
horizon/
├── src/
├── resources/
│   ├── js/
│   ├── sass/
│   └── views/
├── routes/
├── tests/
├── config/
├── composer.json
├── package.json
└── phpunit.xml.dist
```

Exact files may change, but this map gives you the mental model.

---

## Backend Source

`src/` usually contains the core package code:

```text
src/
├── HorizonServiceProvider
├── Console commands
├── MasterSupervisor
├── Supervisor
├── ProcessPool
├── AutoScaler
├── Repositories
├── Job payload helpers
├── Metrics
└── HTTP controllers/API resources
```

Read these areas slowly. Horizon is not only normal Laravel app code; it manages long-running processes.

---

## Service Provider

The service provider usually does package bootstrapping:

- merge config
- publish assets/config
- register routes
- register commands
- bind services
- load views

Contribution areas:

- registration bugs
- config defaults
- route/dashboard behavior
- compatibility with new Laravel versions

---

## Console Commands

Commands are important because Horizon is operated from CLI:

```bash
php artisan horizon
php artisan horizon:terminate
php artisan horizon:pause
php artisan horizon:continue
php artisan horizon:status
php artisan horizon:clear
```

Contribution areas:

- command output clarity
- termination behavior
- status edge cases
- safer failure messages
- tests around command behavior

---

## Supervisors and Process Pools

This is the heart of Horizon.

```text
MasterSupervisor
└── Supervisor
    └── ProcessPool
        ├── active worker processes
        └── terminating worker processes
```

Contribution areas:

- autoscaling
- process count calculation
- worker termination
- balancing between queues
- memory/time limits
- graceful shutdown

This is an advanced area. Start by reading tests before editing code.

---

## Redis Repositories

Horizon stores runtime state in Redis.

Repository classes may track:

- masters
- supervisors
- jobs
- failed jobs
- metrics
- tags
- recent jobs
- pending jobs

Contribution areas:

- cleanup behavior
- missing edge cases
- Redis key consistency
- performance of lookups
- test coverage for repositories

---

## Dashboard

Dashboard areas:

```text
resources/js
resources/sass
resources/views
routes
HTTP controllers
```

Contribution areas:

- loading states
- empty states
- filters
- failed job display
- supervisor/queue cards
- accessibility
- responsive layout
- API error handling

