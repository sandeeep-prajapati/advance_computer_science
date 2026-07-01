# Chapter 11: Testing Biogo

## Test First

For library work, tests are your safety net.

Run all tests:

```bash
go test ./...
```

Run one package:

```bash
go test ./seq
go test ./alphabet
go test ./io/...
```

Run one test:

```bash
go test ./seq -run TestName
```

---

## Table-Driven Tests

Go libraries often use table-driven tests:

```go
tests := []struct {
    name string
    in   string
    want string
}{
    {name: "simple", in: "ACGT", want: "ACGT"},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // check behavior
    })
}
```

Add a new test case before changing code.

---

## Parser Tests

For parsers, include:

- valid minimal input
- valid multi-record input
- malformed input
- partial input
- unexpected EOF
- invalid characters
- wrapped lines

Parser behavior should be explicit.

---

## Benchmark Tests

Benchmarks help performance-sensitive packages:

```go
func BenchmarkSomething(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // operation
    }
}
```

Run:

```bash
go test -bench . -benchmem ./path/to/package
```

---

## Before a PR

Check:

```bash
gofmt -w changed.go
go test ./...
go vet ./...
git diff
```

If you cannot run all tests, say exactly what you ran.
