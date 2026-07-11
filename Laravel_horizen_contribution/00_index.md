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

