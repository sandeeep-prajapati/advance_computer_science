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
