# Chapter 14: CRISPR Screening

## Perturb Thousands of Genes at Once

CRISPR screens connect genotype to phenotype at scale. Instead of editing one target, you perturb many genes across a cell population and read out which perturbations change survival, expression, morphology, or drug response.

```
guide library -> cell population -> selection/phenotype -> sequencing
              -> guide counts -> gene-level effects
```

This chapter focuses on computational analysis and experimental design concepts, not wet-lab execution.

---

## Screen Types

```python
screen_types = {
    "knockout": "Cas9 creates loss-of-function mutations",
    "CRISPRi": "dCas9-KRAB represses transcription",
    "CRISPRa": "dCas9 activator increases transcription",
    "base_editor": "Creates targeted base substitutions",
    "prime_editor": "Creates more flexible small edits",
}
```

Readout types:

```
pooled viability      guide abundance changes after selection
FACS sorting          guides enriched in high/low phenotype bins
single-cell CRISPR    perturbation plus transcriptome per cell
imaging screen        morphology or reporter phenotypes
```

---

## Count Matrix

The key input is guide counts per sample:

```python
import pandas as pd

counts = pd.DataFrame({
    "guide": ["g1", "g2", "g3", "g4"],
    "gene": ["TP53", "TP53", "KRAS", "KRAS"],
    "control_t0": [1000, 980, 1020, 1100],
    "treated_t1": [2200, 2100, 200, 180],
})

counts["log2fc"] = ((counts["treated_t1"] + 1) / (counts["control_t0"] + 1)).apply(lambda x: __import__("math").log2(x))
```

Guide-level changes are aggregated to gene-level scores.

---

## Gene-Level Scoring

```python
def summarize_gene_effects(counts: pd.DataFrame) -> pd.DataFrame:
    return (
        counts.groupby("gene")
        .agg(
            mean_log2fc=("log2fc", "mean"),
            median_log2fc=("log2fc", "median"),
            num_guides=("guide", "count"),
        )
        .sort_values("mean_log2fc")
    )
```

Negative score: guides drop out, suggesting the gene may be essential under that condition.

Positive score: guides enrich, suggesting loss of the gene provides advantage or resistance.

---

## MAGeCK Concept

MAGeCK is a standard method for pooled CRISPR screen analysis. It models guide counts and uses robust rank aggregation to identify genes whose guides are consistently enriched or depleted.

```
Input:  count table + sample design
Output: gene rankings, p-values, FDR, guide-level evidence
```

CLI shape:

```bash
mageck count -l library.txt -n screen --sample-label t0,t1 --fastq t0.fastq t1.fastq
mageck test -k screen.count.txt -t t1 -c t0 -n screen_results
```

---

## Controls and QC

```python
qc_items = {
    "library_coverage": "Enough cells/reads per guide",
    "guide_distribution": "No severe bottlenecking",
    "non_targeting_controls": "Estimate null distribution",
    "positive_controls": "Known essential or pathway genes should score",
    "replicate_correlation": "Biological replicates should agree",
    "copy_number_effects": "Amplified regions can create false positives",
}
```

Computational screens fail quietly when controls are weak. Always plot guide count distributions, replicate correlations, and control gene behavior before trusting hits.

---

## Single-Cell CRISPR Screens

Single-cell perturb-seq style screens measure each cell's expression profile and perturbation label.

```python
workflow = [
    "assign guide identities to cells",
    "filter low-quality cells and ambiguous guide calls",
    "normalize expression",
    "compare perturbed vs control cells",
    "cluster perturbation signatures",
    "infer pathways and gene programs",
]
```

Tools:

```
Seurat / Scanpy       single-cell analysis
scMAGeCK              single-cell CRISPR screen analysis
Mixscape              perturbation signature classification
Cell Ranger Feature Barcode guide capture preprocessing
```

---

## Key Takeaways

- CRISPR screens are high-throughput causal perturbation systems.
- The computational core is count modeling, ranking, and replicate-aware statistics.
- Controls are not optional; they define whether a hit is believable.
- Pooled screens find candidates; follow-up validation establishes mechanism.

