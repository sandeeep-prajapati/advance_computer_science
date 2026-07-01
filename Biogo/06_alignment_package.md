# Chapter 6: Alignment Package

## What Alignment Means

Sequence alignment compares biological sequences to find similarity.

It is used for:

- comparing genes
- comparing proteins
- finding mutations
- detecting conserved regions
- building analysis pipelines

Simple example:

```text
ACGT-A
ACGCTA
```

The gap allows two related sequences to line up.

---

## What to Look For in `align`

The `align` package may include:

- scoring
- gap handling
- pairwise alignment logic
- local or global alignment concepts
- tests for expected alignments

Start with:

```bash
go doc ./align
find align -maxdepth 2 -type f | sort
rg "func Test" align
rg "score|gap|align" align
```

---

## Alignment Vocabulary

| Term | Meaning |
|---|---|
| Match | Same or compatible letters align |
| Mismatch | Different letters align |
| Gap | Missing letter introduced to improve alignment |
| Score | Numeric value for alignment quality |
| Local alignment | Best matching region |
| Global alignment | Whole sequence alignment |

---

## How to Read Algorithmic Code

For alignment code:

1. Find the public function.
2. Read the test that calls it.
3. Identify input sequences.
4. Identify scoring rules.
5. Identify expected output.
6. Only then inspect internal loops.

Do not start inside the deepest loop. It will feel worse than it is.

---

## Contribution Caution

Alignment changes can alter scientific results.

Before changing behavior:

- add regression tests
- compare old and new output
- document the biological meaning
- benchmark if performance may change
- keep public API compatibility in mind
