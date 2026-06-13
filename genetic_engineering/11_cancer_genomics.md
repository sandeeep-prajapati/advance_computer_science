# Chapter 11: Cancer Genomics

## Cancer as Evolution in a Tissue

Cancer is not one disease. It is an evolving population of cells that accumulate mutations, compete for growth, evade immune control, and adapt under treatment pressure.

```
Normal cell -> driver mutation -> clonal expansion -> more mutations
            -> subclones -> treatment -> resistant clone selection
```

CS analogy:

```
Genome instability = high mutation rate
Driver mutation    = bug that increases process priority
Tumor clone        = forked process lineage
Therapy            = adversarial selection pressure
Resistance         = surviving branch with escape condition
```

---

## Somatic vs Germline Variants

```python
variant_context = {
    "germline": "Inherited; present in normal and tumor cells",
    "somatic": "Acquired in tumor; present in tumor but absent from matched normal",
}
```

Tumor-normal comparison:

```python
def classify_variant(tumor_af: float, normal_af: float) -> str:
    if normal_af > 0.25 and tumor_af > 0.25:
        return "likely_germline"
    if tumor_af > 0.05 and normal_af < 0.02:
        return "likely_somatic"
    return "uncertain"
```

Real analysis must account for tumor purity, ploidy, local copy number, mapping artifacts, and sequencing depth.

---

## Driver and Passenger Mutations

```python
drivers = {
    "oncogene": "Gain-of-function; one activated copy can matter",
    "tumor_suppressor": "Loss-of-function; often requires both copies hit",
    "DNA_repair": "Raises mutation rate and creates therapy vulnerabilities",
}

examples = {
    "TP53": "tumor suppressor, genome stress response",
    "KRAS": "oncogene, signaling",
    "EGFR": "oncogene, targetable in some lung cancers",
    "BRCA1/2": "DNA repair, PARP inhibitor sensitivity",
}
```

Passenger mutations hitchhike with tumor growth but do not cause the growth advantage.

---

## Variant Calling Workflow

```
Tumor FASTQ + Normal FASTQ
  -> QC
  -> alignment
  -> duplicate marking
  -> base recalibration
  -> somatic SNV/indel calling
  -> copy-number calling
  -> structural-variant calling
  -> annotation
  -> clinical interpretation
```

Common tools:

```
Mutect2 / Strelka2     somatic SNV/indel calling
FACETS / CNVkit        copy-number and purity
Manta / GRIDSS         structural variants
VEP / OncoKB / CIViC   annotation and actionability
SigProfiler            mutational signatures
maftools               cancer cohort analysis in R
```

---

## Mutational Signatures

Mutation patterns can reveal biological processes:

```python
import pandas as pd

signature_clues = pd.DataFrame([
    {"pattern": "C>T at CpG", "possible_process": "spontaneous deamination / aging"},
    {"pattern": "C>A", "possible_process": "tobacco-associated damage"},
    {"pattern": "small indels at repeats", "possible_process": "mismatch repair deficiency"},
    {"pattern": "homologous recombination deficiency", "possible_process": "BRCA pathway defect"},
])
```

This is like profiling logs to infer which subsystem caused the error.

---

## Tumor Heterogeneity

Variant allele frequency is not equal to clone frequency unless copy number and purity are simple.

```python
def rough_cancer_cell_fraction(vaf: float, purity: float) -> float:
    """
    Extremely simplified diploid assumption.
    Real CCF estimation needs copy number and multiplicity.
    """
    return min(1.0, (2 * vaf) / purity)
```

Single-cell sequencing can resolve subclones more directly, but introduces sparse noisy measurements.

---

## Liquid Biopsy

Circulating tumor DNA enables monitoring from blood:

```
Use cases:
  minimal residual disease
  recurrence monitoring
  resistance mutation detection
  therapy response tracking

Challenges:
  very low allele fractions
  clonal hematopoiesis confounding
  assay sensitivity and false positives
```

---

## Key Takeaways

- Cancer genomics is evolutionary analysis under noisy measurement.
- Matched normal samples make somatic calling much stronger.
- Clinical interpretation requires evidence tiers, not just variant detection.
- Copy number, purity, and heterogeneity are central to making sense of tumor data.

