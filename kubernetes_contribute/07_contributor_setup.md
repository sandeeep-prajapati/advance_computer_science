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

