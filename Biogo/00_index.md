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
