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
