# Chapter 14: Practice Tasks

## Goal

Practice useful contribution movements without starting with a risky PR.

---

## Task 1: Map One Package

Pick one:

- `alphabet`
- `seq`
- `io`
- `complexity`

Write notes:

```text
Package:
Main exported types:
Main exported functions:
Important tests:
One example:
One thing I do not understand yet:
```

Commands:

```bash
go doc ./seq
rg "^type [A-Z]|^func [A-Z]" seq
rg "func Test|func Example" seq
```

---

## Task 2: Add a Test Case Locally

Find a table-driven test.

Add one extra case for existing behavior.

Run:

```bash
go test ./path/to/package -run TestName
```

This teaches the edit-test loop.

---

## Task 3: Parser Edge Case

In an IO package, find a parser test.

Think of one malformed input:

- empty record
- missing header
- invalid sequence character
- unexpected EOF
- quality length mismatch

Add a local test. Run the focused test.

---

## Task 4: Benchmark Reading

Find an existing benchmark:

```bash
rg "func Benchmark" .
```

Run:

```bash
go test -bench . -benchmem ./path/to/package
```

Write down:

```text
Benchmark:
ns/op:
B/op:
allocs/op:
What operation is being measured:
```

---

## Task 5: Prepare a Real Issue Note

Before opening or commenting on an issue:

```text
Problem:
Package:
Reproduction:
Expected behavior:
Actual behavior:
Likely files:
Test plan:
Compatibility risk:
```

Good issue reports save maintainers time.

---

## Task 6: Documentation Improvement

Find an exported type without an example.

Add a small example locally:

```go
func ExampleName() {
    // minimal example
}
```

Run:

```bash
go test ./path/to/package
```

Examples are valuable because users read them first.
