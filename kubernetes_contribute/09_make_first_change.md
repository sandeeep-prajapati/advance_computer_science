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
