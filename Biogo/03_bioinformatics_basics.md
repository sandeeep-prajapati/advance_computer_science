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
