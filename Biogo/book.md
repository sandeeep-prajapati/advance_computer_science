# Biogo Contribution Guide

A practical, chapter-wise guide for learning the [biogo/biogo](https://github.com/biogo/biogo) codebase and preparing useful contributions.

This guide assumes:

- You know basic Go.
- You know basic biology terms like DNA, RNA, protein, FASTA, and FASTQ.
- You want to understand the repo structure before changing code.

---

## Table of Contents

| # | Chapter | Topics Covered |
|---|---|---|
| [01](01_what_is_biogo.md) | **What is Biogo?** | Purpose, project scope, why Go for bioinformatics |
| [02](02_setup_and_installation.md) | **Setup & Installation** | Clone, Go tooling, running tests, local workspace |
| [03](03_bioinformatics_basics.md) | **Bioinformatics Basics** | Sequences, alphabets, quality scores, FASTA, FASTQ |
| [04](04_repo_structure.md) | **Repository Structure** | Official biogo/biogo folder map and package responsibilities |
| [05](05_sequences_and_alphabets.md) | **Sequences & Alphabets** | `seq`, `alphabet`, quality scores, sequence representation |
| [06](06_alignment_package.md) | **Alignment Package** | Alignment concepts, scoring, package reading strategy |
| [07](07_io_and_file_formats.md) | **IO & File Formats** | Bioinformatics readers/writers, streaming, parser tests |
| [08](08_indexing_and_kmers.md) | **Indexing & K-mers** | K-mer indexing, sequence search, package map |
| [09](09_complexity_pwm_features.md) | **Complexity, PWM & Features** | Complexity measures, position weight matrices, genomic features |
| [10](10_concurrency_and_performance.md) | **Concurrency & Performance** | `concurrent`, large datasets, benchmarks, memory care |
| [11](11_testing_biogo.md) | **Testing Biogo** | Unit tests, fixtures, golden files, benchmarks |
| [12](12_contribution_workflow.md) | **Contribution Workflow** | Issues, discussions, PRs, API change caution |
| [13](13_reading_biogo_code.md) | **Reading Biogo Code** | Trace APIs, examples, tests, package-level docs |
| [14](14_practice_tasks.md) | **Practice Tasks** | Safe exercises for docs, tests, parsers, benchmarks |

---

## How to Read This Guide

**If you only want to use Biogo:** Read chapters 1, 2, 3, 5, and 7.

**If you want to contribute:** Read chapters 1, 2, 4, 11, 12, and 13.

**If you want low-level code understanding:** Read chapters 4 through 10, then practice with chapter 14.

---

## Project Links

- Repository: https://github.com/biogo/biogo
- Issues: https://github.com/biogo/biogo/issues
- Go package docs: https://pkg.go.dev/github.com/biogo/biogo
- User group: https://groups.google.com/forum/#!forum/biogo-user
- Developer group: https://groups.google.com/forum/#!forum/biogo-dev
# Chapter 1: What is Biogo?

## Simple Definition

Biogo is a bioinformatics library written in Go.

It provides packages for working with:

- biological sequences
- alphabets such as DNA, RNA, and protein
- sequence alignment
- sequence file formats
- sequence complexity
- k-mer indexing
- genomic features
- concurrent processing patterns

The project is a library, not a single command-line application.

---

## Why Bioinformatics Needs Libraries

Bioinformatics code often has to process large files:

- genomes
- transcriptomes
- metagenomic reads
- FASTA files
- FASTQ files
- alignment inputs
- feature annotations

Good libraries help with:

- correct parsing
- memory-conscious processing
- reusable sequence types
- tested algorithms
- readable research software

---

## Why Go Helps

Go is useful for bioinformatics because it gives:

- compiled performance
- simple syntax
- strong typing
- fast tests
- easy concurrency
- simple distribution

Biogo uses Go packages to make bioinformatics code easier to read, test, and reuse.

---

## Contributor Mindset

When contributing to Biogo, think like a library maintainer:

- avoid breaking public APIs
- keep behavior well tested
- document exported names
- prefer simple Go idioms
- care about memory use on large data
- preserve biological correctness

Small, focused improvements are better than large rewrites.
# Chapter 2: Setup & Installation

## Prerequisites

Install:

- Go
- Git
- a code editor
- `rg` for searching code quickly

Check Go:

```bash
go version
```

---

## Clone the Repo

```bash
git clone https://github.com/biogo/biogo.git
cd biogo
```

Check files:

```bash
find . -maxdepth 2 -type d | sort
```

---

## Install or Fetch Packages

For a normal user project:

```bash
go get github.com/biogo/biogo/...
```

For contributing, clone the repo and run tests locally.

---

## Run Tests

From the repo root:

```bash
go test ./...
```

Run one package:

```bash
go test ./seq/...
go test ./io/...
go test ./align/...
```

Run verbose tests:

```bash
go test -v ./...
```

---

## Useful Developer Commands

```bash
go test ./...
go test -run TestName ./path/to/package
go test -bench . ./path/to/package
go test -race ./...
go vet ./...
gofmt -w file.go
```

Use `gofmt` before committing Go code.

---

## First Local Check

After cloning, do:

```bash
git status
go test ./...
```

If tests pass before your change, later failures are easier to understand.
# Chapter 3: Bioinformatics Basics

## Sequence

A biological sequence is an ordered list of letters.

Examples:

```text
DNA:     A C G T
RNA:     A C G U
Protein: A R N D C Q E G H I L K M F P S T W Y V
```

In code, a sequence is not just a string. It often needs:

- alphabet validation
- metadata
- location
- quality scores
- efficient slicing
- biological operations

---

## Alphabet

An alphabet defines which letters are valid.

Examples:

| Alphabet | Letters |
|---|---|
| DNA | A, C, G, T |
| RNA | A, C, G, U |
| Protein | amino acid symbols |

Ambiguous biological letters may also exist, such as `N` for unknown DNA base.

---

## Quality Scores

Sequencing reads often include quality scores.

FASTQ stores:

- sequence identifier
- sequence bases
- separator
- quality string

Conceptually:

```text
@read1
ACGTACGT
+
IIIIIIII
```

Quality scores estimate confidence in each base call.

---

## FASTA vs FASTQ

| Format | Stores | Common Use |
|---|---|---|
| FASTA | sequence name and sequence | genomes, proteins, references |
| FASTQ | sequence plus quality scores | sequencing reads |

FASTA example:

```text
>seq1
ACGTACGT
```

FASTQ example:

```text
@read1
ACGTACGT
+
IIIIIIII
```

---

## Contributor Warning

Bioinformatics parsers must be strict enough to catch bad data but flexible enough to handle real files.

When changing parser behavior:

- add tests
- include malformed input cases
- include edge cases
- avoid loading huge files into memory unnecessarily
# Chapter 4: Repository Structure

## Official Top-Level Folder Structure

The official `biogo/biogo` repo has a root structure like this:

```text
biogo/
|-- .github/                 # GitHub project files
|-- THIRD_PARTY_LICENSES/    # Third-party license material
|-- align/                   # Sequence alignment code
|-- alphabet/                # Biological alphabets and quality scores
|-- complexity/              # Sequence complexity calculations
|-- concurrent/              # Concurrency helpers
|-- errors/                  # Project error helpers
|-- feat/                    # Feature interfaces and genomic feature concepts
|-- index/kmerindex/         # K-mer indexing
|-- io/                      # Bioinformatics file format IO
|-- morass/                  # External sorting / large-data support
|-- paper/                   # Paper-related material
|-- pwm/                     # Position weight matrices
|-- seq/                     # Sequence types and operations
|-- util/                    # Utility packages
|-- doc.go                   # Package documentation
|-- go.mod                   # Go module definition
`-- version.go               # Version information
```

---

## Package Map

| Package | Main Responsibility |
|---|---|
| `alphabet` | letters, alphabets, quality scores |
| `seq` | sequence representations and sequence operations |
| `io` | reading and writing biological formats |
| `align` | alignment algorithms and scoring support |
| `index/kmerindex` | k-mer based indexing/search |
| `complexity` | measuring sequence complexity |
| `pwm` | position weight matrix logic |
| `feat` | genomic feature abstractions |
| `concurrent` | concurrency utilities |
| `morass` | handling large sortable data |
| `util` | smaller helpers |

---

## Where to Start Reading

Good first packages:

- `alphabet`
- `seq`
- `io`
- `complexity`

Harder first packages:

- `align`
- `index/kmerindex`
- `morass`

Why harder? They may involve algorithms, performance, or large-data behavior.

---

## Reading Strategy

Start with package docs:

```bash
go doc ./seq
go doc ./alphabet
go doc ./io/...
```

Then search:

```bash
rg "type .*Seq" seq alphabet
rg "func Test" seq alphabet io
rg "FASTA|FASTQ" io
```

Read tests before changing behavior.
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
# Chapter 9: Complexity, PWM & Features

## Complexity Package

Sequence complexity measures how repetitive or information-rich a sequence is.

Low complexity examples:

```text
AAAAAAAAAAAA
ATATATATATAT
```

Higher complexity example:

```text
ACGTTGCAAGTC
```

Explore:

```bash
go doc ./complexity
rg "complexity|Entropy|Dust" complexity
rg "func Test" complexity
```

---

## PWM Package

PWM means Position Weight Matrix.

It is often used for motif-like scoring:

```text
position 1: A likely
position 2: C likely
position 3: G likely
position 4: T likely
```

Explore:

```bash
go doc ./pwm
rg "PWM|Matrix|Score" pwm
rg "func Test" pwm
```

---

## Feat Package

The `feat` package relates to biological features.

Features can represent regions such as:

- genes
- exons
- regulatory regions
- sequence intervals

Explore:

```bash
go doc ./feat
rg "type .*Feature|Start|End|Location" feat
rg "func Test" feat
```

---

## Contribution Caution

These packages may encode biological meaning.

Before changing behavior:

- understand the biological model
- inspect tests
- add examples
- avoid renaming exported APIs casually
- benchmark scoring or scanning changes

Documentation improvements are very useful here because the concepts are domain-heavy.
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
# Chapter 12: Contribution Workflow

## Before Opening a PR

Read:

- `README.md`
- `CONTRIBUTING.md`
- package documentation
- related tests
- recent issues or PRs in the same area

Biogo welcomes bug reports and PRs. For API changes or new features, open a discussion first.

---

## Good First Contributions

Good first work:

- documentation fixes
- examples
- extra tests for existing behavior
- parser edge-case tests
- benchmark additions
- small bug fixes with reproduction

Avoid as a first contribution:

- broad rewrites
- API redesign
- changing exported type names
- large performance rewrites without benchmarks
- behavior changes without tests

---

## Branch Workflow

```bash
git clone https://github.com/<you>/biogo.git
cd biogo
git remote add upstream https://github.com/biogo/biogo.git
git checkout -b fix-small-thing
```

Make the smallest useful change.

Before commit:

```bash
gofmt -w changed.go
go test ./...
git diff
```

---

## PR Description Template

```text
Summary:
- What changed?

Reason:
- Why is this needed?

Tests:
- go test ./...
- go test ./specific/package

Notes:
- Any compatibility or performance concerns?
```

---

## API Change Caution

Library APIs are promises to users.

Before changing exported names or behavior:

- open an issue or discussion
- explain the use case
- describe compatibility impact
- include tests
- include examples if useful

Small code changes can have large downstream effects.
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
