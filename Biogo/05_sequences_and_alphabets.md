# Chapter 5: Sequences & Alphabets

## Why This Layer Matters

Most bioinformatics code begins with sequence data.

Before alignment, indexing, or feature analysis, the library needs to know:

- what letters are valid
- how sequence data is stored
- whether quality scores exist
- how sequence operations preserve metadata

---

## Alphabet Package

The `alphabet` package is where biological letters and scoring concepts live.

Read it when working with:

- DNA symbols
- RNA symbols
- protein symbols
- ambiguous letters
- Phred quality scores
- Solexa quality scores

Useful commands:

```bash
go doc ./alphabet
rg "type .*Alphabet" alphabet
rg "Phred|Solexa" alphabet
rg "func Test" alphabet
```

---

## Seq Package

The `seq` package holds sequence types and operations.

Read it when working with:

- sequence representation
- reverse/complement logic
- slicing
- metadata on sequences
- quality-aware sequences

Useful commands:

```bash
go doc ./seq
rg "type .*Seq" seq
rg "Reverse|Complement" seq
rg "func Test" seq
```

---

## Common Mistakes

Avoid treating biological sequences as plain strings too early.

Plain strings lose:

- alphabet information
- quality scores
- validation behavior
- biological intent

Use the package types where the library expects them.

---

## Contribution Ideas

Safe beginner work:

- improve package documentation
- add missing tests for edge cases
- clarify examples
- add tests around invalid alphabet input
- improve error messages without changing API behavior

Risky work:

- changing exported sequence types
- changing validation semantics
- changing quality score conversion behavior
