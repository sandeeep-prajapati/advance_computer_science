# Chapter 8: Indexing & K-mers

## What is a K-mer?

A k-mer is a substring of length `k`.

For sequence:

```text
ACGTAC
```

The 3-mers are:

```text
ACG
CGT
GTA
TAC
```

K-mers are used for fast sequence search and comparison.

---

## K-mer Indexing

The `index/kmerindex` package is for indexing sequence fragments.

Explore it:

```bash
find index -maxdepth 3 -type f | sort
go doc ./index/kmerindex
rg "kmer|Kmer|Index" index
rg "func Test" index
```

---

## Why Indexing is Tricky

Indexing code must care about:

- memory usage
- repeated sequences
- ambiguous bases
- lookup speed
- large input data
- deterministic behavior

Small changes can affect performance significantly.

---

## Reading Strategy

When reading k-mer code:

1. Find how `k` is represented.
2. Find how sequence windows are generated.
3. Find how keys are stored.
4. Find lookup behavior.
5. Find tests for edge cases.

Good search terms:

```bash
rg "Kmer|kmer|KmerIndex|Index" index/kmerindex
rg "ambiguous|N" index/kmerindex
```

---

## Good First Work

Safer work:

- add examples
- add table-driven tests
- add benchmarks for existing behavior
- clarify documentation

Riskier work:

- changing key encoding
- changing matching semantics
- rewriting storage
- optimizing before measuring
