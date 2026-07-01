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

