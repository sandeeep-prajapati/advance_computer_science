# Kubernetes Contribution Guide

A beginner-friendly, chapter-wise guide for learning enough Kubernetes to make real contributions to [kubernetes/kubernetes](https://github.com/kubernetes/kubernetes).

This guide assumes:

- You know a little Docker.
- You are new to Kubernetes.
- You want a practical path to your first useful pull request.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_what_is_kubernetes.md) | **What is Kubernetes?** | Docker vs Kubernetes, cluster mental model, why Kubernetes exists |
| [02](02_kubernetes_core_concepts.md) | **Core Concepts** | Pods, Deployments, Services, ConfigMaps, Secrets, Namespaces, Nodes |
| [03](03_setup_local_lab.md) | **Setup Local Lab** | Docker, kubectl, kind/minikube, first cluster, first workload |
| [04](04_hands_on_kubernetes_basics.md) | **Hands-on Basics** | Run containers, inspect pods, logs, exec, port-forward, YAML manifests |
| [05](05_kubernetes_architecture.md) | **Architecture** | API server, etcd, scheduler, controller-manager, kubelet, kube-proxy |
| [06](06_kubernetes_repo_map.md) | **Repository Map** | Official kubernetes/kubernetes folder structure, component paths, staging modules, what to avoid first |
| [07](07_contributor_setup.md) | **Contributor Setup** | Fork, clone, remotes, CLA, GitHub workflow, local build expectations |
| [08](08_find_first_issue.md) | **Find Your First Issue** | good first issue, help wanted, SIGs, asking for assignment, picking safely |
| [09](09_make_first_change.md) | **Make Your First Change** | Branch, edit, test, commit, push, PR description, release notes |
| [10](10_testing_kubernetes.md) | **Testing Kubernetes** | Unit tests, integration tests, e2e tests, focused test commands |
| [11](11_code_review_and_prow.md) | **Code Review & Prow** | Bots, labels, reviewers, approvers, CI commands, review etiquette |
| [12](12_learning_roadmap.md) | **Learning Roadmap** | 8-week plan from Docker basics to first Kubernetes PR |
| [13](13_reading_kubernetes_code.md) | **Reading Kubernetes Code** | Trace commands, errors, tests, packages, and call paths without getting lost |
| [14](14_api_server_internals.md) | **API Server Internals** | Request flow, authentication, admission, validation, storage, watch behavior |
| [15](15_scheduler_controller_kubelet.md) | **Scheduler, Controllers & Kubelet** | Scheduling loop, reconciliation loops, node agents, common source paths |
| [16](16_practice_contribution_tasks.md) | **Practice Contribution Tasks** | Safe exercises for docs, tests, kubectl, validation, and bug reproduction |

---

## How to Read This Guide

**If you know Docker but not Kubernetes:** Read chapters 1 to 5 first.

**If you want to contribute quickly:** Read chapters 1, 2, 7, 8, 9, and 11.

**If you want to understand the codebase:** Read chapters 5, 6, 10, and then choose a SIG area.

**If you want low-level code reading practice:** Read chapters 13, 14, and 15 after chapter 6.

**If you feel lost:** Follow chapter 12 week by week. Do not rush the first PR.

---

## Contribution Path at a Glance

```text
Docker knowledge
    |
    v
Kubernetes basics
    |
    v
Local cluster with kind
    |
    v
Understand kubernetes/kubernetes repo
    |
    v
Sign CLA + fork + clone
    |
    v
Pick good first issue / help wanted issue
    |
    v
Ask with /assign
    |
    v
Make small change + run focused tests
    |
    v
Open PR + respond to review
```

---

## Official Links

- Kubernetes repo: https://github.com/kubernetes/kubernetes
- Contributor guide: https://www.kubernetes.dev/docs/guide/
- First contribution: https://www.kubernetes.dev/docs/guide/first-contribution/
- GitHub workflow: https://www.kubernetes.dev/docs/guide/github-workflow/
- Good first issues: https://go.k8s.io/good-first-issue
- Help wanted issues: https://go.k8s.io/help-wanted
- Community repo: https://github.com/kubernetes/community
- Kubernetes Slack: https://slack.k8s.io/
