# Chapter 10: Testing Kubernetes

## Test Types

| Type | Meaning | Beginner Use |
|---|---|---|
| Unit test | Tests one package/function | Best first choice |
| Integration test | Tests components together | Useful but heavier |
| e2e test | Tests real cluster behavior | Powerful but expensive |
| Conformance | Required behavior across clusters | Not a first stop |

---

## Start With Unit Tests

Run package tests:

```bash
go test ./pkg/example
```

Run one test:

```bash
go test ./pkg/example -run TestName
```

Verbose:

```bash
go test ./pkg/example -run TestName -v
```

---

## Find Existing Tests

```bash
rg "func Test" pkg/example
rg "TestName" .
```

Add your test near similar tests. Follow the existing style.

---

## Table-Driven Tests

Kubernetes Go tests often use table-driven structure:

```go
func TestThing(t *testing.T) {
    tests := []struct {
        name string
        input string
        want string
    }{
        {
            name: "simple case",
            input: "abc",
            want: "abc",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Thing(tt.input)
            if got != tt.want {
                t.Fatalf("Thing() = %q, want %q", got, tt.want)
            }
        })
    }
}
```

---

## When Tests Are Hard

If you cannot run the full suite locally:

- run the smallest package test
- mention exactly what you ran in the PR
- ask maintainers if there is a better focused test
- let CI run the larger jobs

Do not run huge test commands blindly. Kubernetes CI is complex and local machine time matters.

---

## What Maintainers Like

Maintainers like tests that:

- fail before the fix
- pass after the fix
- are focused
- do not sleep unnecessarily
- do not depend on test order
- explain the behavior being protected

