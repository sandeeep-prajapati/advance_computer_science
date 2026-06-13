# Chapter 9: Epigenetics and Regulation

## The Control Plane Above DNA

Your genome is mostly the same in every cell, but a neuron and a liver cell behave very differently. Epigenetics is the runtime configuration layer that decides which parts of the genome are readable, writable, or locked down.

```
DNA sequence        = source code
Chromatin state     = file permissions and compression
DNA methylation     = long-term silencing flag
Histone marks       = local access metadata
Transcription factor binding = process scheduler decisions
Enhancers/promoters = distributed config and entry points
```

---

## Core Regulatory Signals

```python
regulatory_layers = {
    "DNA methylation": "Usually represses promoters when high",
    "H3K4me3": "Active promoter mark",
    "H3K27ac": "Active enhancer/promoter mark",
    "H3K27me3": "Polycomb repression mark",
    "ATAC-seq": "Open chromatin accessibility",
    "ChIP-seq": "Protein or histone-mark binding enrichment",
    "Hi-C": "3D genome contacts and loops",
}

for layer, meaning in regulatory_layers.items():
    print(f"{layer:18} -> {meaning}")
```

Regulation is combinatorial. A gene may be silent because the promoter is methylated, because chromatin is closed, because an enhancer is inactive, or because the required transcription factor is absent.

---

## ChIP-seq and ATAC-seq as Signal Tracks

Both assays produce genomic intervals with signal strength:

```python
import pandas as pd

peaks = pd.DataFrame({
    "chrom": ["chr1", "chr1", "chr2"],
    "start": [100_000, 180_000, 55_000],
    "end": [100_450, 181_200, 55_700],
    "signal": [12.4, 34.8, 8.1],
    "annotation": ["promoter", "enhancer", "intergenic"],
})

promoter_peaks = peaks.query("annotation == 'promoter' and signal > 10")
enhancer_peaks = peaks.query("annotation == 'enhancer'").sort_values("signal", ascending=False)
```

Typical pipeline:

```
FASTQ -> QC -> alignment -> duplicate removal -> peak calling -> annotation
      -> motif enrichment -> differential accessibility/binding -> reports
```

Common tools:

```
FastQC / MultiQC        quality reports
bowtie2 / bwa / minimap2 alignment
samtools / sambamba     BAM processing
MACS2 / SEACR           peak calling
bedtools                interval operations
HOMER / MEME            motif discovery
deepTools               signal visualization
```

---

## Working with Genomic Intervals

Most regulatory genomics is interval algebra:

```python
def overlaps(a_start: int, a_end: int, b_start: int, b_end: int) -> bool:
    return a_start < b_end and b_start < a_end

def overlap_length(a_start: int, a_end: int, b_start: int, b_end: int) -> int:
    return max(0, min(a_end, b_end) - max(a_start, b_start))

gene_promoter = ("chr7", 55_086_724, 55_088_724)
peak = ("chr7", 55_087_100, 55_087_500)

same_chrom = gene_promoter[0] == peak[0]
is_hit = same_chrom and overlaps(gene_promoter[1], gene_promoter[2], peak[1], peak[2])
print(is_hit)
```

In production, use `bedtools`, `pybedtools`, `pyranges`, or `GenomicRanges` instead of hand-rolled loops.

---

## Motif Enrichment

A transcription factor motif is a position weight matrix. It scores whether a DNA window resembles a known binding pattern.

```python
import math

PWM = {
    "A": [0.8, 0.1, 0.1, 0.2],
    "C": [0.1, 0.7, 0.1, 0.2],
    "G": [0.05, 0.1, 0.7, 0.2],
    "T": [0.05, 0.1, 0.1, 0.4],
}

def score_motif(seq: str, pwm: dict[str, list[float]]) -> float:
    score = 0.0
    for i, base in enumerate(seq.upper()):
        score += math.log2(pwm.get(base, [0.25] * len(seq))[i] / 0.25)
    return score

print(score_motif("ACGT", PWM))
```

Motif enrichment answers: "Which transcription factors probably explain these peaks?"

---

## Enhancer-Gene Linking

Promoters are near genes; enhancers can be far away. Link enhancers to genes using several weak signals:

```python
signals = {
    "distance": "Enhancers often regulate nearby genes, but not always",
    "correlation": "Enhancer activity correlates with gene expression",
    "3d_contact": "Hi-C/pcHi-C shows physical contact",
    "perturbation": "CRISPRi/CRISPRa validates causality",
    "cell_type": "Regulation is highly context-specific",
}
```

For ML models, represent each candidate enhancer-gene pair as features:

```python
features = [
    "distance_to_tss",
    "enhancer_h3k27ac",
    "enhancer_atac",
    "promoter_h3k4me3",
    "hic_contact_score",
    "expression_correlation",
]
```

---

## Key Takeaways

- Epigenetics is not a change to DNA sequence; it is a change to how DNA is used.
- Regulatory genomics is mostly interval processing, signal tracks, statistics, and sequence modeling.
- A regulatory claim is strongest when multiple data layers agree: accessibility, histone marks, expression, 3D contact, and perturbation.

