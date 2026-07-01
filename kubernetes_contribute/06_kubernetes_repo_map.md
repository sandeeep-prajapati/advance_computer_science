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
