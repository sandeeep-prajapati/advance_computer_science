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

