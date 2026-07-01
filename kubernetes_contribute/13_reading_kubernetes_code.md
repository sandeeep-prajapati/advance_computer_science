# Chapter 13: Reading Kubernetes Code

## Goal

Kubernetes is too large to read from top to bottom. A contributor reads it by following one behavior.

Good starting points:

- an exact error message
- a failing test
- a command name
- an API field
- a function mentioned in an issue
- a package changed by a recent PR

Bad starting point:

```text
I will understand the whole repo first.
```

That path is slow and frustrating.

---

## The Trace Method

Use this pattern:

```text
symptom
  |
  v
search exact text or symbol
  |
  v
read the nearest test
  |
  v
read the implementation
  |
  v
find the caller
  |
  v
make the smallest change
```

Example:

```bash
rg "forbidden" pkg staging test
rg "CrashLoopBackOff" pkg staging test
rg "NewCmdApply" cmd staging
```

The first useful file is often a test, not the production code.

---

## Follow Tests First

Tests show what the code is supposed to do.

Look for:

```bash
rg "func Test" pkg/kubelet
rg "func Test" pkg/scheduler
rg "func Test" staging/src/k8s.io/kubectl
```

When you find a test:

1. Read the test name.
2. Read the table cases.
3. Read expected output or expected error.
4. Only then read the implementation.

This prevents you from guessing intent.

---

## Common Go Patterns in Kubernetes

Kubernetes Go code often uses:

- table-driven tests
- interfaces for fake clients and test doubles
- informers and listers
- contexts and cancellation
- feature gates
- options/config structs
- generated clients and deep-copy code

Example table-driven test shape:

```go
tests := []struct {
    name    string
    input   string
    wantErr bool
}{
    {
        name:    "empty input is rejected",
        input:   "",
        wantErr: true,
    },
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test body
    })
}
```

If you add behavior, usually add one table case first.

---

## Read Ownership

Before editing a package, check ownership:

```bash
find . -name OWNERS -path "*pkg/kubelet*"
find . -name OWNERS -path "*staging/src/k8s.io/kubectl*"
```

`OWNERS` files tell you:

- who reviews changes
- which SIG owns the area
- whether the area has special rules

If a package is sensitive, start with tests or docs before changing behavior.

---

## Avoid Wandering

Stop after 20 to 30 minutes if you cannot explain:

- what package you are in
- what behavior you are tracing
- which test covers it
- what the smallest possible change is

Then return to the issue and search again with a more exact term.

Reading Kubernetes code is navigation work. The win is not reading more files. The win is finding the right file.
