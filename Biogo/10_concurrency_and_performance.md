# Chapter 10: Concurrency & Performance

## Why Performance Matters

Biological datasets can be large.

Performance issues often appear in:

- parsers
- alignment loops
- k-mer indexing
- sequence scanning
- sorting large intermediate data
- repeated allocations

Go makes concurrency convenient, but concurrency still needs careful design.

---

## Concurrent Package

Explore:

```bash
go doc ./concurrent
find concurrent -maxdepth 2 -type f | sort
rg "goroutine|chan|WaitGroup|parallel" concurrent
rg "func Test" concurrent
```

Look for:

- worker patterns
- channels
- cancellation behavior
- ordering guarantees
- error handling

---

## Morass Package

The `morass` package is related to handling large sortable data.

Explore:

```bash
go doc ./morass
find morass -maxdepth 2 -type f | sort
rg "sort|temp|file|merge" morass
rg "func Test|Benchmark" morass
```

---

## Benchmark Before Optimizing

Use benchmarks:

```bash
go test -bench . ./...
go test -bench . -benchmem ./path/to/package
```

Useful benchmark output includes:

- ns/op
- B/op
- allocs/op

Do not optimize from feeling. Measure first.

---

## Performance-Safe Contribution Pattern

1. Add or find a benchmark.
2. Run baseline.
3. Make one change.
4. Run benchmark again.
5. Explain tradeoffs in the PR.

If a change improves speed but makes code much harder to understand, discuss it first.
