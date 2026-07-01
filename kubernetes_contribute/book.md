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

# Chapter 2: Kubernetes Core Concepts

## Pod

A **Pod** is the smallest unit Kubernetes runs.

Usually:

- 1 Pod = 1 application container
- Sometimes 1 Pod = app container + helper sidecar container

Example:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
    - name: nginx
      image: nginx:1.27
      ports:
        - containerPort: 80
```

You rarely manage Pods directly in production. You usually create a Deployment.

---

## Deployment

A **Deployment** manages replicated Pods and rollouts.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.27
          ports:
            - containerPort: 80
```

If one Pod dies, the Deployment controller creates another.

---

## Service

Pods can be replaced, so their IPs are not stable. A **Service** gives a stable network identity.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx
spec:
  selector:
    app: nginx
  ports:
    - port: 80
      targetPort: 80
```

Think:

```text
Deployment = keep app running
Service    = let others reach app
```

---

## ConfigMap and Secret

Use **ConfigMaps** for non-sensitive configuration.

Use **Secrets** for sensitive values, but remember: default Kubernetes Secrets are encoded, not magically secure. Real security depends on cluster configuration.

---

## Namespace

A **Namespace** groups resources.

```bash
kubectl create namespace learning
kubectl get pods -n learning
```

Namespaces are useful for separating environments, teams, and experiments.

---

## Labels and Selectors

Labels are key-value metadata.

Selectors find resources by labels.

```text
Pod label:
  app=nginx

Service selector:
  app=nginx

Result:
  Service sends traffic to Pods with app=nginx
```

Many Kubernetes systems depend on labels.

---

## The Contributor Angle

When reading Kubernetes code, always ask:

- What object is this code working with?
- Which controller owns the behavior?
- Is this desired state or actual state?
- Is this validation, scheduling, reconciliation, networking, storage, or node behavior?

# Chapter 3: Setup Local Lab

## Tools You Need

Install:

- Docker
- Git
- Go
- `kubectl`
- `kind` or `minikube`

For contribution work, `kind` is a good default because it runs Kubernetes clusters inside Docker containers.

---

## Check Docker

```bash
docker version
docker ps
```

If Docker works, you already have the base runtime needed for a local Kubernetes lab.

---

## Install kubectl

`kubectl` is the Kubernetes command-line tool.

Check:

```bash
kubectl version --client
```

---

## Install kind

`kind` means **Kubernetes in Docker**.

Check:

```bash
kind version
```

Create a cluster:

```bash
kind create cluster --name learning
```

Check nodes:

```bash
kubectl get nodes
```

Delete when done:

```bash
kind delete cluster --name learning
```

---

## First Workload

```bash
kubectl create deployment hello --image=nginx:1.27
kubectl get deployments
kubectl get pods
kubectl expose deployment hello --port=80 --type=ClusterIP
kubectl port-forward service/hello 8080:80
```

Open:

```text
http://localhost:8080
```

Stop `port-forward` with `Ctrl+C`.

---

## Cleanup

```bash
kubectl delete service hello
kubectl delete deployment hello
```

---

## What You Should Understand Before Moving On

- `kind` creates a real Kubernetes cluster locally.
- `kubectl` talks to the API server.
- A Deployment creates Pods.
- A Service gives stable access to Pods.
- You can inspect everything with `kubectl get` and `kubectl describe`.

# Chapter 4: Hands-on Kubernetes Basics

## Important kubectl Commands

```bash
kubectl get pods
kubectl get deployments
kubectl get services
kubectl describe pod <pod-name>
kubectl logs <pod-name>
kubectl exec -it <pod-name> -- sh
kubectl delete pod <pod-name>
```

When you delete a Pod owned by a Deployment, Kubernetes creates a replacement. That is reconciliation in action.

---

## Create from YAML

Create `nginx.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.27
```

Apply:

```bash
kubectl apply -f nginx.yaml
kubectl get pods
```

Delete:

```bash
kubectl delete -f nginx.yaml
```

---

## Imperative vs Declarative

| Style | Example | Best For |
|---|---|---|
| Imperative | `kubectl create deployment ...` | Quick learning |
| Declarative | `kubectl apply -f app.yaml` | Real work |

Kubernetes is built around declarative APIs. You describe what you want; controllers make it happen.

---

## Debugging Loop

Use this loop when something fails:

```bash
kubectl get pods
kubectl describe pod <pod>
kubectl logs <pod>
kubectl get events --sort-by=.lastTimestamp
```

Common beginner issues:

- image name typo
- container crashes immediately
- port mismatch
- labels do not match Service selector
- resource is created in a different namespace

---

## Skill Goal

Before contributing to Kubernetes code, you should be comfortable with:

- creating a Deployment
- exposing it with a Service
- reading Pod logs
- describing failed Pods
- understanding labels/selectors
- deleting and recreating resources

# Chapter 5: Kubernetes Architecture

## Big Picture

```text
kubectl
  |
  v
API Server
  |
  +--> etcd
  +--> Scheduler
  +--> Controller Manager
  +--> Kubelet on each node
  +--> Kube Proxy / networking
```

Kubernetes is API-driven. Most components watch the API server and react to changes.

---

## Control Plane Components

| Component | Job |
|---|---|
| API server | Validates and serves Kubernetes API requests |
| etcd | Stores cluster state |
| scheduler | Chooses which node should run a Pod |
| controller-manager | Runs controllers that reconcile desired state |

---

## Node Components

| Component | Job |
|---|---|
| kubelet | Ensures containers run on a node |
| container runtime | Runs containers, commonly containerd |
| kube-proxy | Handles Service networking rules |

---

## Controllers

A controller watches resources and tries to make reality match the desired state.

Example:

```text
Desired state:
Deployment says replicas=3

Actual state:
Only 2 Pods exist

Controller action:
Create 1 more Pod
```

This pattern appears everywhere in Kubernetes.

---

## API Objects

Most Kubernetes features are API objects:

- Pod
- Deployment
- Service
- ConfigMap
- Secret
- Node
- PersistentVolume
- Job

The API server validates and stores these objects. Controllers and kubelets act on them.

---

## Contributor Mental Map

When reading a bug or code path, classify it:

| Area | Likely Component |
|---|---|
| API validation/defaulting | API server / API machinery |
| Pod placement | scheduler |
| Pod lifecycle on node | kubelet |
| Replica management | controller-manager |
| Service networking | kube-proxy / networking |
| CLI behavior | kubectl / staging modules |
| Tests failing | test framework or component behavior |

# Chapter 6: kubernetes/kubernetes Repository Map

## The Repository is Large

`kubernetes/kubernetes` is one of the largest Go projects in open source. Do not try to understand everything first.

Your goal is:

1. Learn the top-level map.
2. Pick a small area.
3. Follow one issue deeply.

---

## Official Top-Level Folder Structure

At the root of the official repo you will see a structure like this:

```text
kubernetes/
|-- .github/          # GitHub workflows, issue templates, PR templates
|-- CHANGELOG/        # Release changelogs
|-- LICENSES/         # License files for the project and dependencies
|-- api/              # API definitions, generated OpenAPI material, API checks
|-- build/            # Build tooling, containers, and release build support
|-- cluster/          # Cluster bootstrap scripts and older cluster tooling
|-- cmd/              # Main binary entry points
|-- docs/             # Repo-local developer and design documentation
|-- hack/             # Developer scripts for build, test, verify, and codegen
|-- logo/             # Kubernetes logo assets
|-- pkg/              # Main Kubernetes implementation packages
|-- plugin/           # Plugin-related code
|-- staging/          # Code developed here but published as k8s.io/* modules
|-- test/             # Unit helpers, integration tests, e2e tests, test images
|-- third_party/      # Third-party helper files
`-- vendor/           # Vendored Go dependencies
```

For contribution work, the folders you will open most often are:

| Path | Meaning |
|---|---|
| `cmd/` | Entry points for Kubernetes binaries |
| `pkg/` | Core implementation packages |
| `staging/` | Code published as separate Go modules under `k8s.io/*` |
| `test/` | e2e, integration, and test utilities |
| `hack/` | Developer scripts for build, test, verify, and generated code |
| `build/` | Build tooling and images |
| `cluster/` | Cluster scripts and legacy tooling |
| `api/` | API-related definitions and OpenAPI material |
| `docs/` | Repo-local docs |

---

## Component Entry Points

| Binary | Path |
|---|---|
| kube-apiserver | `cmd/kube-apiserver` |
| kube-controller-manager | `cmd/kube-controller-manager` |
| kube-scheduler | `cmd/kube-scheduler` |
| kubelet | `cmd/kubelet` |
| kube-proxy | `cmd/kube-proxy` |
| kubectl | `cmd/kubectl` and staged CLI modules |

The `cmd/` folder usually contains startup wiring. After the binary starts, most real behavior lives in `pkg/` or `staging/`.

---

## Where Core Logic Usually Lives

| Area | Common Paths |
|---|---|
| API server behavior | `cmd/kube-apiserver/`, `pkg/controlplane/`, `pkg/registry/`, `staging/src/k8s.io/apiserver/` |
| API types | `pkg/apis/`, `staging/src/k8s.io/api/`, `staging/src/k8s.io/apimachinery/` |
| Scheduler | `cmd/kube-scheduler/`, `pkg/scheduler/` |
| Controllers | `cmd/kube-controller-manager/`, `pkg/controller/` |
| Kubelet | `cmd/kubelet/`, `pkg/kubelet/` |
| Kube-proxy | `cmd/kube-proxy/`, `pkg/proxy/` |
| kubectl | `cmd/kubectl/`, `staging/src/k8s.io/kubectl/` |
| Client libraries | `staging/src/k8s.io/client-go/` |
| Tests | `test/integration/`, `test/e2e/`, package-level `*_test.go` files |

Think of the repo like this:

```text
cmd/ starts a component
  |
  v
pkg/ contains much of the component behavior
  |
  v
staging/src/k8s.io/* contains shared libraries used by many components
  |
  v
test/ and *_test.go files prove the behavior
```

---

## Staging Folder

`staging/` is important. It contains code that is developed inside `kubernetes/kubernetes` but published as separate modules, such as:

- `k8s.io/api`
- `k8s.io/apimachinery`
- `k8s.io/client-go`
- `k8s.io/apiserver`
- `k8s.io/kubectl`

If you see an import like:

```go
import "k8s.io/apimachinery/pkg/apis/meta/v1"
```

the source may live under `staging/src/k8s.io/apimachinery/...`.

---

## How to Read a Code Path

Do not start by reading the whole repo. Start with one symptom, one command, or one issue.

Example for a kubelet bug:

```text
issue text
  |
  v
rg "exact error message"
  |
  v
pkg/kubelet/...
  |
  v
nearby *_test.go
  |
  v
focused unit test
```

Example for a kubectl bug:

```text
kubectl command
  |
  v
cmd/kubectl/
  |
  v
staging/src/k8s.io/kubectl/
  |
  v
printers, builders, command options, or tests
```

Example for an API behavior bug:

```text
API kind / field name
  |
  v
pkg/apis/ or staging/src/k8s.io/api/
  |
  v
validation, defaulting, conversion, registry
  |
  v
unit or integration test
```

---

## Good First Areas

Good beginner areas often include:

- small docs fixes inside the repo
- tests around existing behavior
- error message improvements requested by maintainers
- focused bug fixes with clear reproduction
- small `kubectl` behavior fixes

Hard first areas:

- scheduler design changes
- API changes
- storage behavior
- networking behavior
- large refactors
- performance-sensitive kubelet changes
- generated API machinery changes without maintainer guidance

---

## Reading Strategy

Use search:

```bash
rg "error message text"
rg "FunctionName"
rg "type Pod"
rg "TODO"
```

Then trace:

```text
test -> function -> package -> caller -> command/component
```

Do not read randomly. Start from an issue, failing test, or exact symbol.

---

## Useful Commands

After cloning the official repo:

```bash
git clone https://github.com/kubernetes/kubernetes.git
cd kubernetes
find . -maxdepth 2 -type d | sort
```

Search fast:

```bash
rg "CrashLoopBackOff"
rg "func .*Schedule" pkg/scheduler
rg "type Pod" pkg/apis staging/src/k8s.io/api
rg "NewCmd" staging/src/k8s.io/kubectl cmd/kubectl
```

Check changed files before committing:

```bash
git status
git diff
```
# Chapter 7: Contributor Setup

## Official Prerequisites

Before contributing:

- Create a GitHub account.
- Read the Kubernetes Code of Conduct.
- Sign the Kubernetes Contributor License Agreement.
- Read the Contributor Guide.
- Set up Git correctly.

The upstream `CONTRIBUTING.md` points contributors to the Kubernetes Contributor Guide and says the CLA is required.

Links:

- https://github.com/kubernetes/kubernetes/blob/master/CONTRIBUTING.md
- https://www.kubernetes.dev/docs/guide/
- https://git.k8s.io/community/CLA.md

---

## Fork and Clone

Fork:

```text
https://github.com/kubernetes/kubernetes
```

Clone your fork:

```bash
mkdir -p ~/src/k8s.io
cd ~/src/k8s.io
git clone https://github.com/<your-github-username>/kubernetes.git
cd kubernetes
git remote add upstream https://github.com/kubernetes/kubernetes.git
git remote set-url --push upstream no_push
git remote -v
```

`upstream` points to the real Kubernetes repo. You fetch from it, but you do not push to it.

---

## Keep Your Branch Updated

```bash
git fetch upstream
git checkout master
git rebase upstream/master
```

Create a branch:

```bash
git checkout -b fix-small-thing
```

Kubernetes docs recommend `fetch` + `rebase` instead of merge commits.

---

## Go Setup

Kubernetes is written mostly in Go.

Check:

```bash
go version
```

Use the Go version expected by the repo. After cloning, check:

```bash
cat go.mod
```

Look for the `go` version line and any toolchain line.

---

## Build Smoke Test

From the Kubernetes repo:

```bash
make help
make all WHAT=cmd/kubectl
```

For a first contribution, do not try to build every possible target unless you need it. Use focused builds and tests.

---

## Before You Edit

Run:

```bash
git status
git branch
```

Make sure:

- you are on your own branch
- your branch is based on current upstream master
- the issue is assigned to you or maintainers said your approach is okay

# Chapter 8: Find Your First Issue

## Best Labels

Use these:

- `good first issue`
- `help wanted`

Official shortcuts:

- https://go.k8s.io/good-first-issue
- https://go.k8s.io/help-wanted

The Kubernetes contributor docs say `good first issue` means Kubernetes members have committed to giving extra help to new contributors.

---

## Good First Contribution Types

Start with:

- documentation improvement
- small test coverage improvement
- clear error message improvement
- small bug with reproduction steps
- issue triage
- broken link or outdated reference

Avoid as first PR:

- huge refactor
- formatting-only PR across many files
- API design change
- behavior change without an issue
- drive-by linter warning cleanup

Kubernetes maintainers explicitly warn that every PR creates review and CI cost, so make the change useful.

---

## How to Ask for Assignment

On an issue, comment:

```text
/assign
```

or:

```text
/assign @your-github-username
```

The Kubernetes bot can assign it to you.

---

## How to Ask a Good Question

Bad:

```text
I am new. Please guide me.
```

Better:

```text
I would like to work on this. I found the likely code path in pkg/...
My plan is to add/update a test for ... and then change ...
Does this approach sound acceptable?
/assign
```

Maintainers are busy. A specific question gets a better answer.

---

## SIGs

Kubernetes is organized into **SIGs**: Special Interest Groups.

Examples:

- SIG API Machinery
- SIG Apps
- SIG CLI
- SIG Node
- SIG Network
- SIG Scheduling
- SIG Storage
- SIG Testing

If you know Docker, SIG Node may feel interesting later because kubelet and container runtime behavior live close to your existing knowledge. For a first PR, choose the issue, not the SIG prestige.

---

## First Issue Checklist

Pick an issue only if:

- it has `good first issue` or `help wanted`
- it has recent maintainer activity
- the expected change is small
- you can reproduce or understand the problem
- nobody else is actively working on it
- tests are mentioned or you can identify a focused test

# Chapter 9: Make Your First Change

## Workflow

```bash
git fetch upstream
git checkout master
git rebase upstream/master
git checkout -b my-first-k8s-fix
```

Edit the smallest set of files needed.

---

## Run Focused Tests

Examples:

```bash
go test ./pkg/some/package
go test ./staging/src/k8s.io/some/module/...
```

If the issue or package has special test commands, use those.

Do not claim you ran tests that you did not run.

---

## Commit

```bash
git status
git diff
git add <files>
git commit
```

Kubernetes requires the CLA check to pass. If a maintainer or repository template asks for a signed-off commit, use `git commit -s`; otherwise, a normal commit is fine.

Good commit message:

```text
Fix kubectl error message for missing foo

Explain the behavior change in one or two lines if needed.
```

---

## Push

```bash
git push origin my-first-k8s-fix
```

Open a PR from your fork to:

```text
kubernetes/kubernetes:master
```

---

## PR Description

Use the template. Include:

- what changed
- why it changed
- issue link
- tests run
- release note block if required

Example:

```text
Fixes #12345

Tests:
- go test ./pkg/example
```

For release notes, follow the PR template. If no user-visible change:

```text
NONE
```

---

## PR Size Rule

Your first PR should be boringly small.

Good:

- one bug
- one package
- one focused test

Risky:

- multiple unrelated files
- broad cleanup
- style changes mixed with behavior changes
- no test when behavior changed
# Chapter 10: Testing Kubernetes

## Test Types

| Type | Meaning | Beginner Use |
|---|---|---|
| Unit test | Tests one package/function | Best first choice |
| Integration test | Tests components together | Useful but heavier |
| e2e test | Tests real cluster behavior | Powerful but expensive |
| Conformance | Required behavior across clusters | Not a first stop |

---

## Start With Unit Tests

Run package tests:

```bash
go test ./pkg/example
```

Run one test:

```bash
go test ./pkg/example -run TestName
```

Verbose:

```bash
go test ./pkg/example -run TestName -v
```

---

## Find Existing Tests

```bash
rg "func Test" pkg/example
rg "TestName" .
```

Add your test near similar tests. Follow the existing style.

---

## Table-Driven Tests

Kubernetes Go tests often use table-driven structure:

```go
func TestThing(t *testing.T) {
    tests := []struct {
        name string
        input string
        want string
    }{
        {
            name: "simple case",
            input: "abc",
            want: "abc",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Thing(tt.input)
            if got != tt.want {
                t.Fatalf("Thing() = %q, want %q", got, tt.want)
            }
        })
    }
}
```

---

## When Tests Are Hard

If you cannot run the full suite locally:

- run the smallest package test
- mention exactly what you ran in the PR
- ask maintainers if there is a better focused test
- let CI run the larger jobs

Do not run huge test commands blindly. Kubernetes CI is complex and local machine time matters.

---

## What Maintainers Like

Maintainers like tests that:

- fail before the fix
- pass after the fix
- are focused
- do not sleep unnecessarily
- do not depend on test order
- explain the behavior being protected

# Chapter 11: Code Review and Prow

## Kubernetes Uses Bots

Kubernetes uses automation heavily. You will see bot comments, labels, test results, and commands.

Common examples:

```text
/assign
/sig cli
/kind bug
/retest
/test pull-kubernetes-unit
```

Do not be scared by bot noise. Read the failing job name and logs carefully.

---

## Review Roles

| Role | Meaning |
|---|---|
| Contributor | Opens issues/PRs |
| Reviewer | Reviews code quality and correctness |
| Approver | Can approve changes for an owned area |
| Maintainer | Trusted project participant with broader responsibility |

Kubernetes uses `OWNERS` files to route reviews. If you change files in a folder, reviewers/approvers are usually inferred from nearby ownership metadata.

---

## Responding to Review

Good response:

```text
Thanks, updated the test case and renamed the variable as suggested.
```

If you disagree:

```text
I may be missing context. I chose this because ...
Would you prefer ... instead?
```

Stay specific and calm. Review is part of the contribution, not a final exam.

---

## Updating Your PR

Make changes on the same branch:

```bash
git add <files>
git commit --amend
git push --force-with-lease
```

For multiple meaningful commits, keep them clean. For small first PRs, one commit is usually easiest.

---

## CI Failures

When CI fails:

1. Open the failed job.
2. Find the first real error.
3. Check whether the error relates to your change.
4. If unrelated/flaky, wait or ask before retesting.
5. If related, fix locally and push again.

Avoid repeatedly typing `/retest` without understanding the failure.

---

## Merge

After review and approval, Kubernetes automation handles merging when requirements pass. Your job is to keep the PR clear, responsive, and small enough to review.

# Chapter 12: Learning Roadmap

## Goal

In 8 weeks, go from basic Docker knowledge to a realistic first Kubernetes contribution.

Move faster if you already understand a topic. Move slower if commands feel magical.

---

## Week 1: Kubernetes Basics

Learn:

- Pod
- Deployment
- Service
- Namespace
- labels/selectors

Practice:

```bash
kind create cluster --name learning
kubectl create deployment nginx --image=nginx:1.27
kubectl expose deployment nginx --port=80
kubectl get pods,svc
kubectl port-forward service/nginx 8080:80
```

Outcome: you can run and inspect a simple app.

---

## Week 2: YAML and Debugging

Learn:

- declarative manifests
- `kubectl apply`
- `kubectl describe`
- `kubectl logs`
- events

Outcome: you can diagnose common Pod failures.

---

## Week 3: Architecture

Learn:

- API server
- etcd
- scheduler
- controllers
- kubelet

Outcome: you can classify bugs by component.

---

## Week 4: Go Basics for Kubernetes

Learn:

- Go modules
- packages
- interfaces
- table-driven tests
- `go test`

Outcome: you can read and run a small Go package test.

---

## Week 5: Repository Orientation

Do:

```bash
git clone https://github.com/<you>/kubernetes.git
cd kubernetes
git remote add upstream https://github.com/kubernetes/kubernetes.git
rg "func Test" pkg
```

Learn:

- `cmd/`
- `pkg/`
- `staging/`
- `test/`
- `hack/`
- `OWNERS`

Outcome: the repo feels large but not impossible.

---

## Week 6: Community Setup

Do:

- sign the CLA
- join Kubernetes Slack
- read the Contributor Guide
- read the GitHub workflow
- browse `good first issue`

Outcome: you understand the contribution process.

---

## Week 7: Pick One Issue

Find:

- `good first issue`
- `help wanted`
- recent maintainer comments
- small scope

Comment:

```text
I would like to work on this. My plan is ...
/assign
```

Outcome: one issue is assigned or you have a maintainer-approved direction.

---

## Week 8: First PR

Do:

- create a branch
- make the smallest useful change
- add/update a test if needed
- run focused tests
- open PR
- respond to review

Outcome: your first Kubernetes PR is open.

---

## After First PR

Good next steps:

- stay in the same SIG area for a while
- review simple PRs after you gain context
- help triage issues
- improve tests
- attend a SIG meeting

Consistency beats intensity. One small useful contribution every week teaches more than one giant stressful PR.

# Chapter 13: Reading Kubernetes Code

## Goal

Kubernetes is too large to read from top to bottom. A contributor reads it by following one behavior.

Good starting points:

- an exact error message
- a failing test
- a command name
- an API field
- a function mentioned in an issue
- a package changed by a recent PR

Bad starting point:

```text
I will understand the whole repo first.
```

That path is slow and frustrating.

---

## The Trace Method

Use this pattern:

```text
symptom
  |
  v
search exact text or symbol
  |
  v
read the nearest test
  |
  v
read the implementation
  |
  v
find the caller
  |
  v
make the smallest change
```

Example:

```bash
rg "forbidden" pkg staging test
rg "CrashLoopBackOff" pkg staging test
rg "NewCmdApply" cmd staging
```

The first useful file is often a test, not the production code.

---

## Follow Tests First

Tests show what the code is supposed to do.

Look for:

```bash
rg "func Test" pkg/kubelet
rg "func Test" pkg/scheduler
rg "func Test" staging/src/k8s.io/kubectl
```

When you find a test:

1. Read the test name.
2. Read the table cases.
3. Read expected output or expected error.
4. Only then read the implementation.

This prevents you from guessing intent.

---

## Common Go Patterns in Kubernetes

Kubernetes Go code often uses:

- table-driven tests
- interfaces for fake clients and test doubles
- informers and listers
- contexts and cancellation
- feature gates
- options/config structs
- generated clients and deep-copy code

Example table-driven test shape:

```go
tests := []struct {
    name    string
    input   string
    wantErr bool
}{
    {
        name:    "empty input is rejected",
        input:   "",
        wantErr: true,
    },
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test body
    })
}
```

If you add behavior, usually add one table case first.

---

## Read Ownership

Before editing a package, check ownership:

```bash
find . -name OWNERS -path "*pkg/kubelet*"
find . -name OWNERS -path "*staging/src/k8s.io/kubectl*"
```

`OWNERS` files tell you:

- who reviews changes
- which SIG owns the area
- whether the area has special rules

If a package is sensitive, start with tests or docs before changing behavior.

---

## Avoid Wandering

Stop after 20 to 30 minutes if you cannot explain:

- what package you are in
- what behavior you are tracing
- which test covers it
- what the smallest possible change is

Then return to the issue and search again with a more exact term.

Reading Kubernetes code is navigation work. The win is not reading more files. The win is finding the right file.
# Chapter 14: API Server Internals

## What the API Server Does

The API server is the front door of Kubernetes.

Every `kubectl` request and almost every component request goes through it.

It handles:

- authentication
- authorization
- admission
- defaulting
- validation
- conversion between API versions
- storage in etcd
- watch/list behavior

---

## Request Flow

A simplified request looks like this:

```text
client request
  |
  v
authentication
  |
  v
authorization
  |
  v
admission
  |
  v
defaulting and validation
  |
  v
storage strategy
  |
  v
etcd
  |
  v
response or watch event
```

When debugging API behavior, classify the failure:

| Symptom | Likely Area |
|---|---|
| `Unauthorized` | authentication |
| `Forbidden` | authorization / RBAC |
| object rejected with field error | validation or admission |
| field silently changed | defaulting or admission mutation |
| object stored differently | conversion or storage |
| watch/list odd behavior | API machinery or storage watch cache |

---

## Important Source Paths

| Area | Common Paths |
|---|---|
| API server entry point | `cmd/kube-apiserver/` |
| Kubernetes API server wiring | `pkg/controlplane/` |
| REST storage and registry logic | `pkg/registry/` |
| Internal API types | `pkg/apis/` |
| External API types | `staging/src/k8s.io/api/` |
| API machinery | `staging/src/k8s.io/apimachinery/` |
| Generic API server | `staging/src/k8s.io/apiserver/` |

The exact file depends on the resource. For a Pod behavior, search for Pod-related registry, validation, or strategy code.

---

## Validation, Defaulting, Conversion

These words appear constantly in API work:

| Term | Meaning |
|---|---|
| Defaulting | Fill missing fields with default values |
| Validation | Reject invalid objects |
| Conversion | Convert between API versions or internal/external types |
| Admission | Plugin phase that can mutate or reject requests |
| Strategy | Resource-specific create/update/delete rules |

Useful searches:

```bash
rg "Validate.*Pod" pkg staging
rg "SetDefaults" pkg staging
rg "Convert_.*Pod" pkg staging
rg "Strategy" pkg/registry
```

---

## Watches Matter

Kubernetes components do not constantly poll everything manually. They watch the API server.

Example:

```text
Deployment object changes
  |
  v
API server stores new state
  |
  v
controller receives watch event
  |
  v
controller reconciles ReplicaSets and Pods
```

If a controller did not react, the bug may be:

- object was not stored
- watch event was not delivered
- informer cache did not sync
- controller logic ignored the event

---

## Beginner-Safe API Contributions

Safer first contributions:

- improve validation test coverage
- improve error messages when maintainers ask
- add missing unit tests for existing API behavior
- update documentation comments
- fix small kubectl/API interaction bugs

Risky first contributions:

- adding new API fields
- changing default behavior
- changing storage format
- changing conversion logic
- changing admission order

API compatibility is serious. For API changes, wait for maintainer guidance.
# Chapter 15: Scheduler, Controllers & Kubelet

## Three Different Jobs

These components are easy to mix up:

| Component | Main Question |
|---|---|
| Scheduler | Which node should run this Pod? |
| Controller manager | Does actual cluster state match desired state? |
| Kubelet | Are the containers running correctly on this node? |

They all watch the API server, but they do different work.

---

## Scheduler Mental Model

The scheduler handles Pods that do not have a node assigned yet.

Simplified flow:

```text
unscheduled Pod
  |
  v
filter nodes that cannot run it
  |
  v
score remaining nodes
  |
  v
choose best node
  |
  v
bind Pod to node
```

Common source paths:

```text
cmd/kube-scheduler/
pkg/scheduler/
pkg/scheduler/framework/
```

Useful searches:

```bash
rg "framework.Plugin" pkg/scheduler
rg "Score" pkg/scheduler
rg "Filter" pkg/scheduler
```

Scheduler changes can affect many workloads. For a first PR, prefer tests, docs, or small bug fixes with clear maintainer direction.

---

## Controller Mental Model

A controller is a reconciliation loop.

```text
watch resources
  |
  v
enqueue key
  |
  v
sync desired state vs actual state
  |
  v
create, update, or delete resources
  |
  v
retry on error
```

Example:

```text
Deployment wants 3 replicas
  |
  v
ReplicaSet has 2 Pods
  |
  v
controller creates 1 more Pod
```

Common source paths:

```text
cmd/kube-controller-manager/
pkg/controller/
```

Useful searches:

```bash
rg "sync.*" pkg/controller
rg "enqueue" pkg/controller
rg "workqueue" pkg/controller
```

When reading controller code, find the `sync` function. It usually contains the heart of the behavior.

---

## Kubelet Mental Model

The kubelet runs on every node. It talks to the API server and the container runtime.

Simplified flow:

```text
Pod assigned to this node
  |
  v
kubelet observes desired Pod
  |
  v
create Pod sandbox
  |
  v
start containers through runtime
  |
  v
report Pod status back to API server
```

Common source paths:

```text
cmd/kubelet/
pkg/kubelet/
pkg/kubelet/kuberuntime/
pkg/kubelet/status/
```

Useful searches:

```bash
rg "SyncPod" pkg/kubelet
rg "PodStatus" pkg/kubelet
rg "kuberuntime" pkg/kubelet
```

Kubelet code touches real node behavior. Keep changes small and test-focused.

---

## How Components Work Together

For a Deployment:

```text
kubectl apply deployment.yaml
  |
  v
API server stores Deployment
  |
  v
Deployment controller creates ReplicaSet
  |
  v
ReplicaSet controller creates Pods
  |
  v
Scheduler assigns Pods to nodes
  |
  v
Kubelet starts containers
  |
  v
Kubelet reports status
```

When something fails, ask:

| Failure | Check |
|---|---|
| Deployment exists but no Pods | controller path |
| Pod exists but no node | scheduler path |
| Pod assigned but not running | kubelet/runtime path |
| Service cannot reach Pod | networking/kube-proxy path |
| kubectl output looks wrong | kubectl/printer path |

This classification saves hours.
# Chapter 16: Practice Contribution Tasks

## Goal

Before taking a real issue, practice the moves contributors use every day.

These exercises are not meant to become PRs immediately. They build confidence with search, tests, and small edits.

---

## Exercise 1: Trace a kubectl Command

Pick a command:

```bash
kubectl get pods
kubectl describe pod nginx
kubectl apply -f pod.yaml
```

Search:

```bash
rg "NewCmdGet" cmd staging/src/k8s.io/kubectl
rg "NewCmdDescribe" cmd staging/src/k8s.io/kubectl
rg "NewCmdApply" cmd staging/src/k8s.io/kubectl
```

Write down:

- command entry point
- options struct
- validation function
- run function
- nearest tests

Outcome: you understand how one CLI command is wired.

---

## Exercise 2: Find a Validation Test

Pick an API object:

- Pod
- Service
- Deployment
- ConfigMap
- Secret

Search:

```bash
rg "Validate.*Pod" pkg staging
rg "Validate.*Service" pkg staging
rg "Validate.*Deployment" pkg staging
```

Then find tests:

```bash
rg "Validate.*Pod" pkg staging -g "*_test.go"
```

Outcome: you can connect an API rule to its test.

---

## Exercise 3: Reproduce a Small Failure

Create a bad manifest in a local cluster:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: bad-pod
spec:
  containers:
    - name: app
      image: nginx:1.27
      resources:
        requests:
          cpu: "not-a-number"
```

Apply it:

```bash
kubectl apply -f bad-pod.yaml
```

Then search the repo for the error text.

Outcome: you learn how user-facing validation errors map back to source code.

---

## Exercise 4: Read One Controller Sync Function

Pick one controller folder:

```text
pkg/controller/deployment/
pkg/controller/replicaset/
pkg/controller/job/
```

Search:

```bash
rg "sync" pkg/controller/deployment
rg "workqueue" pkg/controller/deployment
```

Write a short note:

```text
This controller watches:
It creates/updates/deletes:
The main sync function is:
The nearest tests are:
```

Outcome: reconciliation starts to feel concrete.

---

## Exercise 5: Make a Local Test-Only Change

Do this locally first:

1. Find a table-driven test.
2. Add one extra test case for existing behavior.
3. Run the focused package test.
4. Revert the change if it was only for practice.

Example:

```bash
go test ./pkg/scheduler/... -run TestName
go test ./staging/src/k8s.io/kubectl/... -run TestName
```

Outcome: you learn the edit-test loop without risking a messy PR.

---

## Exercise 6: Prepare a Real Issue Note

Before asking for assignment, write this in your own notes:

```text
Issue:
Component:
Files likely involved:
Tests likely involved:
Risk level:
My planned first step:
Question for maintainers:
```

Then comment only if your plan is specific.

Good comment:

```text
I would like to work on this. I found the behavior in pkg/... and the related test in pkg/... .
My plan is to add a focused regression test first, then make the smallest code change needed.
/assign
```

Avoid vague comments:

```text
Can I work on this?
```

Maintainers trust contributors who show they have already looked.
