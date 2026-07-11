# Laravel Horizon Contribution Guide: Complete Book

> A practical single-file guide for learning Laravel Horizon internals and choosing useful contribution areas.

---

# Laravel Horizon Contribution Guide

A practical guide for learning Laravel Horizon deeply enough to contribute useful pull requests.

Laravel Horizon is the official dashboard and worker supervisor for Laravel Redis queues. Contributing to it means understanding **Laravel queues**, **Redis**, **worker processes**, **balancing/autoscaling**, **metrics**, **failed jobs**, **Vue dashboard UI**, and Laravel's contribution style.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_what_is_horizon.md) | **What Horizon Is** | Purpose, architecture, Redis queues, dashboard, supervisors |
| [02](02_repo_setup_and_workflow.md) | **Repo Setup & Workflow** | Fork, clone, branches, Composer, npm, tests, PR flow |
| [03](03_architecture_map.md) | **Architecture Map** | Source folders, service provider, commands, supervisors, repositories |
| [04](04_core_concepts_to_learn.md) | **Core Concepts to Learn** | Queues, Redis, processes, balancing, metrics, tags, failed jobs |
| [05](05_backend_contribution_areas.md) | **Backend Contribution Areas** | Autoscaling, balancing, supervisors, Redis storage, tests |
| [06](06_frontend_contribution_areas.md) | **Frontend Contribution Areas** | Vue dashboard, assets, UX, accessibility, polling states |
| [07](07_testing_strategy.md) | **Testing Strategy** | Unit tests, feature tests, process tests, regression tests |
| [08](08_first_pr_roadmap.md) | **First PR Roadmap** | How to pick issues, reproduce bugs, write small PRs |
| [09](09_good_work_areas_for_you.md) | **Where You Can Work** | Beginner, intermediate, advanced, and long-term areas |
| [10](10_sources_and_tracking.md) | **Sources & Tracking** | Official links, issue tracking checklist, reading plan |
| [book](book.md) | **Combined Book** | All chapters in one file |

---

## Fast Start

Read in this order:

```text
01 -> 02 -> 03 -> 04 -> 09 -> 08
```

Then choose one area:

```text
Small docs/test PR
-> frontend dashboard fix
-> backend regression test
-> queue balancing/autoscaling bug
```

---

## Contribution Mindset

Good Horizon PRs are usually:

- small
- reproducible
- tested
- aligned with Laravel style
- focused on one bug or improvement
- careful about backward compatibility

Avoid starting with a huge feature. Start by learning the codebase through tests and small fixes.


---

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


---

# Chapter 2: Repo Setup & Workflow

## Official Repository

Repository:

```text
https://github.com/laravel/horizon
```

The current active branch observed for the repo is:

```text
5.x
```

Always check the repo before starting because Laravel branches change with framework versions.

---

## Contribution Flow

```text
fork laravel/horizon
clone your fork
add upstream remote
create feature branch
write failing test or reproduction
fix code
run tests
open PR to correct branch
respond to review
```

---

## Clone

```bash
git clone git@github.com:YOUR_USERNAME/horizon.git
cd horizon
git remote add upstream git@github.com:laravel/horizon.git
```

Sync later:

```bash
git fetch upstream
git checkout 5.x
git merge upstream/5.x
```

---

## Install Dependencies

PHP:

```bash
composer install
```

Frontend:

```bash
npm install
```

Build frontend assets:

```bash
npm run dev
npm run build
```

Exact scripts can change, so always inspect:

```bash
cat composer.json
cat package.json
```

---

## Branch Naming

Use focused branch names:

```text
fix-autoscaler-terminating-process-count
test-failed-jobs-filtering
improve-dashboard-empty-state
docs-supervisor-balancing-option
```

Avoid:

```text
my-changes
fix-all
big-update
```

---

## Commit Style

Keep commits small and clear.

Good:

```text
Add regression test for terminating process autoscaling
Fix process count used during autoscale decisions
```

Weak:

```text
changes
final fix
update files
```

---

## Before Opening PR

Checklist:

- issue reproduced or behavior explained
- test added or updated
- code follows existing style
- no unrelated formatting changes
- frontend assets built only if repo expects it
- branch targets correct Horizon branch
- PR body explains problem and solution

---

## PR Body Template

```md
## Problem

Explain what is broken or missing.

## Cause

Explain why it happens.

## Fix

Explain what changed.

## Tests

List commands run and what the test covers.
```


---

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


---

# Chapter 4: Core Concepts to Learn

## Laravel Queues

You must understand:

- jobs
- dispatching
- queue connections
- queue names
- `queue:work`
- retries
- timeouts
- failed jobs
- middleware
- tags

Practice:

```php
dispatch(new ImportProductsJob($fileId))->onQueue('imports');
```

Then run:

```bash
php artisan queue:work redis --queue=imports
```

---

## Redis

Horizon is Redis-based.

Learn:

- lists
- sorted sets
- hashes
- key expiration
- Lua scripts basics
- Redis connection names
- Redis cluster caveats

You do not need to become a Redis expert first, but you must know enough to understand Horizon's repositories.

---

## Processes

Horizon manages OS processes.

Learn:

- process start/stop
- signals
- graceful shutdown
- long-running PHP behavior
- memory limits
- process exit codes

Important idea:

```text
HTTP request lifecycle is short.
Horizon worker lifecycle is long.
```

Long-running code has different bugs.

---

## Supervisor Concepts

Horizon supervisor config controls:

- connection
- queues
- balance strategy
- min/max processes
- memory
- timeout
- tries
- nice value
- max time
- max jobs

You should know how config maps to runtime process behavior.

---

## Balancing and Autoscaling

This is one of Horizon's hardest areas.

Concepts:

- current process count
- desired process count
- queue size
- runtime average
- scale up
- scale down
- terminating processes
- min/max process limits
- balancing between multiple queues

Good contribution area after you understand tests.

---

## Metrics

Horizon dashboard depends on metrics.

Learn:

- throughput
- runtime
- job counts
- failed counts
- snapshots
- trimming old metrics

Contribution areas:

- inaccurate numbers
- missing edge cases
- dashboard display mismatch
- cleanup behavior

---

## Tags

Horizon can show jobs by tags.

Job tags help track:

- user IDs
- order IDs
- model IDs
- tenant IDs

Contribution areas:

- tag display
- tag search
- tag cleanup
- handling special characters


---

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


---

# Chapter 6: Frontend Contribution Areas

## Horizon Dashboard Stack

Horizon's dashboard uses frontend assets in the package, commonly involving:

- Vue components
- JavaScript API calls
- SCSS/CSS
- Blade views
- compiled assets

Exact versions can change, so inspect:

```bash
cat package.json
ls resources/js
ls resources/sass
ls resources/views
```

---

## Good Frontend Contribution Areas

| Area | Difficulty | Examples |
|---|---|---|
| empty states | beginner | better text when no failed jobs |
| loading states | beginner | spinners, disabled controls |
| copy clarity | beginner | clearer labels |
| accessibility | intermediate | keyboard focus, aria labels |
| responsive layout | intermediate | mobile dashboard issues |
| API error handling | intermediate | failed polling request state |
| dashboard performance | advanced | polling and rendering large lists |

---

## Dashboard API Flow

Frontend usually follows:

```text
Vue component
└── HTTP request to Horizon route
    └── controller/API endpoint
        └── repository reads Redis
            └── JSON response
```

If dashboard data is wrong, check both:

- backend JSON response
- frontend rendering logic

---

## Frontend PR Checklist

- reproduce the UI issue
- inspect API response in browser devtools
- identify whether bug is backend or frontend
- make minimal UI change
- build assets if repo expects committed assets
- include screenshot or explanation in PR

---

## Useful UI Improvements

Good areas for first dashboard PR:

- improve empty dashboard states
- add clearer labels for paused supervisors
- make failed-job error display easier to read
- fix small responsive layout problem
- improve keyboard focus on dashboard controls
- improve loading behavior when polling fails

---

## Avoid at First

Avoid starting with:

- total dashboard redesign
- changing build system
- replacing Vue version
- changing routes/API shape
- adding large new feature

These are high-friction PRs.


---

# Chapter 7: Testing Strategy

## Why Tests Matter

Horizon controls production workers. A small bug can affect job processing.

A good contribution usually needs a test.

---

## Test Types

| Test Type | Covers |
|---|---|
| unit test | one class, one behavior |
| feature test | Laravel app behavior |
| console test | Artisan command behavior |
| repository test | Redis-backed state |
| regression test | prevents a known bug from returning |

---

## Regression Test Pattern

Use this structure:

```text
test name explains bug
arrange exact bad state
run the method/command
assert fixed behavior
```

Example:

```php
public function test_terminating_processes_do_not_block_autoscaling()
{
    // Arrange supervisor and process pool state.
    // Run autoscaler.
    // Assert expected process count.
}
```

---

## What Makes a Good Test

Good:

- small
- focused
- deterministic
- names the behavior
- fails before the fix
- passes after the fix

Bad:

- tests many things at once
- depends on timing too much
- passes even without the fix
- only checks implementation details

---

## Testing Redis Logic

Repository tests should cover:

- empty Redis data
- existing records
- deleted records
- malformed payloads if relevant
- trimming behavior
- expiration behavior

---

## Testing Process Logic

Process and autoscaling tests should cover:

- min process limits
- max process limits
- scale up
- scale down
- terminating processes
- multiple queues
- long-running jobs
- balancing strategy

These tests are harder but valuable.

---

## Commands to Discover

Inspect available scripts:

```bash
cat composer.json
cat phpunit.xml.dist
```

Likely commands:

```bash
vendor/bin/phpunit
composer test
```

Use the repo's actual scripts.


---

# Chapter 8: First PR Roadmap

## Best First PR Types

Choose one:

1. Add missing test for existing behavior.
2. Improve small dashboard empty/loading state.
3. Fix typo or unclear docs in comments/readme.
4. Add regression test for a confirmed bug.
5. Fix small command output bug.

Do not start with autoscaling rewrite.

---

## Step 1: Read Before Coding

Spend one session reading:

- `README.md`
- `composer.json`
- `package.json`
- `src/HorizonServiceProvider.php`
- console commands
- `AutoScaler`
- `ProcessPool`
- tests folder

Write notes:

```text
What does this class do?
Who calls it?
What tests cover it?
What Redis keys does it touch?
```

---

## Step 2: Reproduce

For a bug:

```text
Can I explain exact input/state?
Can I make a failing test?
Can I reproduce manually?
```

If you cannot reproduce, do not patch yet.

---

## Step 3: Small Fix

Fix one behavior.

Good PR scope:

```text
One bug
One test file
One or two source files
Clear PR explanation
```

Bad PR scope:

```text
Bug fix + refactor + UI redesign + formatting + dependency upgrade
```

---

## Step 4: Run Tests

Run targeted test first:

```bash
vendor/bin/phpunit --filter test_name
```

Then full suite if possible:

```bash
vendor/bin/phpunit
```

For frontend:

```bash
npm run build
```

Use actual scripts from `package.json`.

---

## Step 5: PR Explanation

Write:

```md
This fixes an autoscaling edge case where ...

The problem happens when ...

The fix changes ...

I added a regression test covering ...
```

Maintainers should understand the PR without reverse-engineering it.

---

## First 30-Day Plan

### Week 1

- Read docs and repo structure.
- Run tests locally.
- Read 5 existing merged PRs.
- Read 5 open/closed issues.

### Week 2

- Add one small test.
- Try one dashboard UI improvement locally.

### Week 3

- Study `AutoScaler`, `Supervisor`, `ProcessPool`.
- Reproduce one bug in a test.

### Week 4

- Open a small PR.
- Respond calmly to review.
- Keep learning from maintainer feedback.


---

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


---

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

