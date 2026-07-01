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
