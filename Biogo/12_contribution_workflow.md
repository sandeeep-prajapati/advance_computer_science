# Chapter 12: Contribution Workflow

## Before Opening a PR

Read:

- `README.md`
- `CONTRIBUTING.md`
- package documentation
- related tests
- recent issues or PRs in the same area

Biogo welcomes bug reports and PRs. For API changes or new features, open a discussion first.

---

## Good First Contributions

Good first work:

- documentation fixes
- examples
- extra tests for existing behavior
- parser edge-case tests
- benchmark additions
- small bug fixes with reproduction

Avoid as a first contribution:

- broad rewrites
- API redesign
- changing exported type names
- large performance rewrites without benchmarks
- behavior changes without tests

---

## Branch Workflow

```bash
git clone https://github.com/<you>/biogo.git
cd biogo
git remote add upstream https://github.com/biogo/biogo.git
git checkout -b fix-small-thing
```

Make the smallest useful change.

Before commit:

```bash
gofmt -w changed.go
go test ./...
git diff
```

---

## PR Description Template

```text
Summary:
- What changed?

Reason:
- Why is this needed?

Tests:
- go test ./...
- go test ./specific/package

Notes:
- Any compatibility or performance concerns?
```

---

## API Change Caution

Library APIs are promises to users.

Before changing exported names or behavior:

- open an issue or discussion
- explain the use case
- describe compatibility impact
- include tests
- include examples if useful

Small code changes can have large downstream effects.
