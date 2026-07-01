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
