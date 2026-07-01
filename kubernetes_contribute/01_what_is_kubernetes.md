# Chapter 1: What is Kubernetes?

## Simple Definition

**Kubernetes** is a system for running containers across many machines.

If Docker answers:

> How do I package and run one container?

Kubernetes answers:

> How do I run many containers reliably in production?

---

## Docker vs Kubernetes

| Topic | Docker | Kubernetes |
|---|---|---|
| Main job | Build and run containers | Orchestrate containers |
| Scale | Usually one machine | Many machines |
| Restart failed app | Manual or Docker restart policy | Automatic |
| Networking | Container/network setup | Service discovery and stable networking |
| Deployment | `docker run` / Compose | Declarative YAML through API server |
| Production operations | Limited | Scheduling, rollout, scaling, healing |

Your Docker knowledge is useful because Kubernetes still runs container images. Kubernetes just adds the production control system around those containers.

---

## The Mental Model

```text
You give Kubernetes desired state:

"Run 3 copies of nginx"

Kubernetes continuously works to make actual state match desired state:

- choose machines
- start containers
- restart failures
- expose networking
- roll out updates
```

This is the most important idea in Kubernetes: **desired state reconciliation**.

---

## Why Kubernetes Exists

Without Kubernetes, teams must manually solve:

- Which machine should run this container?
- What happens if the container crashes?
- How do other services find it?
- How do we roll out a new version safely?
- How do we scale up during traffic?
- How do we attach config and secrets?

Kubernetes gives standard APIs for these problems.

---

## First Vocabulary

| Word | Meaning |
|---|---|
| Cluster | A group of machines managed by Kubernetes |
| Node | One machine in the cluster |
| Pod | Smallest runnable unit; usually one app container |
| Deployment | Keeps a desired number of Pods running |
| Service | Stable network name/IP for Pods |
| Namespace | Logical grouping inside a cluster |
| Controller | Code that watches state and fixes drift |
| API server | Front door of the Kubernetes control plane |

---

## What Contributors Work On

The `kubernetes/kubernetes` repo contains the core Kubernetes code. Contributors work on:

- API behavior
- controllers
- scheduler logic
- kubelet behavior
- CLI behavior related to core components
- tests
- documentation inside the repo
- bug fixes
- cleanup that maintainers actually want

For a first contribution, avoid huge features. Start small, useful, and reviewable.

