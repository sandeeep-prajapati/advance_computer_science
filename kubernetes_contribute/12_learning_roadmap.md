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

