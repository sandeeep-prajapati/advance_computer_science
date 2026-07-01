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

