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
