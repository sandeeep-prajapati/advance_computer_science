# Chapter 13: Reading Biogo Code

## Start from Package Documentation

For Go libraries, package docs are the front door.

Use:

```bash
go doc ./seq
go doc ./alphabet
go doc ./io/...
go doc ./align
```

Then open `doc.go` files and package comments.

---

## Trace Method

Use this flow:

```text
package docs
  |
  v
exported type or function
  |
  v
example or test
  |
  v
implementation
  |
  v
edge-case tests
```

Do not begin with private helper functions.

---

## Search Patterns

Find exported types:

```bash
rg "^type [A-Z]" seq alphabet io align
```

Find exported functions:

```bash
rg "^func [A-Z]" seq alphabet io align
```

Find tests:

```bash
rg "func Test" seq alphabet io align
```

Find examples:

```bash
rg "func Example" .
```

---

## Read Errors Carefully

Bioinformatics libraries should return helpful errors.

Search:

```bash
rg "error|Errorf|New" errors io seq alphabet
```

When changing errors:

- preserve useful context
- avoid breaking tests unnecessarily
- add malformed input tests
- do not hide parse positions if they exist

---

## Stop Conditions

Pause and re-search if you cannot answer:

- Which package am I changing?
- Is this exported API?
- What tests cover it?
- What biological concept does it represent?
- What input file or sequence triggers it?

Good contributors navigate carefully.
