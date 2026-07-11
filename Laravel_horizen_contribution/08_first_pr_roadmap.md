# Chapter 8: First PR Roadmap

## Best First PR Types

Choose one:

1. Add missing test for existing behavior.
2. Improve small dashboard empty/loading state.
3. Fix typo or unclear docs in comments/readme.
4. Add regression test for a confirmed bug.
5. Fix small command output bug.

Do not start with autoscaling rewrite.

---

## Step 1: Read Before Coding

Spend one session reading:

- `README.md`
- `composer.json`
- `package.json`
- `src/HorizonServiceProvider.php`
- console commands
- `AutoScaler`
- `ProcessPool`
- tests folder

Write notes:

```text
What does this class do?
Who calls it?
What tests cover it?
What Redis keys does it touch?
```

---

## Step 2: Reproduce

For a bug:

```text
Can I explain exact input/state?
Can I make a failing test?
Can I reproduce manually?
```

If you cannot reproduce, do not patch yet.

---

## Step 3: Small Fix

Fix one behavior.

Good PR scope:

```text
One bug
One test file
One or two source files
Clear PR explanation
```

Bad PR scope:

```text
Bug fix + refactor + UI redesign + formatting + dependency upgrade
```

---

## Step 4: Run Tests

Run targeted test first:

```bash
vendor/bin/phpunit --filter test_name
```

Then full suite if possible:

```bash
vendor/bin/phpunit
```

For frontend:

```bash
npm run build
```

Use actual scripts from `package.json`.

---

## Step 5: PR Explanation

Write:

```md
This fixes an autoscaling edge case where ...

The problem happens when ...

The fix changes ...

I added a regression test covering ...
```

Maintainers should understand the PR without reverse-engineering it.

---

## First 30-Day Plan

### Week 1

- Read docs and repo structure.
- Run tests locally.
- Read 5 existing merged PRs.
- Read 5 open/closed issues.

### Week 2

- Add one small test.
- Try one dashboard UI improvement locally.

### Week 3

- Study `AutoScaler`, `Supervisor`, `ProcessPool`.
- Reproduce one bug in a test.

### Week 4

- Open a small PR.
- Respond calmly to review.
- Keep learning from maintainer feedback.

