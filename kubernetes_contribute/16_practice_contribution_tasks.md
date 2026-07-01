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
