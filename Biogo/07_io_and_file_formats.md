# Chapter 7: IO & File Formats

## Why IO Matters

Bioinformatics is file-heavy.

Most workflows begin by reading:

- FASTA
- FASTQ
- feature files
- alignment-related files
- large text or binary datasets

Parser correctness is critical. A small parsing bug can corrupt a full analysis.

---

## IO Package

The `io` directory contains format-specific readers and writers.

Explore it:

```bash
find io -maxdepth 3 -type f | sort
go doc ./io/...
rg "FASTA|FASTQ|Reader|Writer" io
rg "func Test" io
```

---

## Reader Pattern

Many Go parsers follow this rough shape:

```text
create reader
  |
  v
read one record
  |
  v
handle EOF or parse error
  |
  v
repeat
```

A good bioinformatics reader should:

- stream data
- report useful errors
- handle valid edge cases
- reject malformed records
- avoid unnecessary full-file loading

---

## Parser Test Cases

When changing format behavior, test:

- empty file
- single record
- multiple records
- malformed header
- invalid sequence letters
- missing quality line
- long wrapped sequence
- EOF at awkward position

These tests are often more valuable than changing code immediately.

---

## Contributor Rule

Do not silently accept bad data unless the package already intentionally does that.

If behavior changes, document why.
