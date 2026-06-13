# Chapter 13: Metagenomics and Microbiome

## Sequencing Communities Instead of Individuals

Metagenomics studies mixtures of organisms from an environment: gut, soil, ocean, skin, wastewater, or hospital surfaces.

```
Single genome sequencing: one sample -> one organism
Metagenomics:             one sample -> many organisms + abundance estimates
```

Two common approaches:

```python
approaches = {
    "16S rRNA amplicon": "Cheap bacterial/archaeal profiling; genus-level often",
    "shotgun metagenomics": "All DNA; species/strain and functional potential",
}
```

---

## 16S Amplicon Workflow

```
FASTQ -> QC -> denoise -> ASVs/OTUs -> taxonomy assignment
      -> abundance table -> diversity analysis -> differential abundance
```

Tools:

```
QIIME 2       end-to-end microbiome analysis
DADA2         denoising and ASV inference
SILVA/GTDB    reference taxonomies
phyloseq      R ecosystem for microbiome statistics
```

ASVs are exact inferred sequences. OTUs are clustered groups, often at 97% identity.

---

## Shotgun Metagenomics Workflow

```
FASTQ -> host read removal -> taxonomic profiling -> assembly/binning
      -> gene prediction -> functional annotation -> pathway analysis
```

Tools:

```
Kraken2 / Bracken      taxonomic classification and abundance
MetaPhlAn              marker-gene profiling
MEGAHIT / metaSPAdes   assembly
MetaBAT / CONCOCT      binning genomes from assemblies
Prokka / Bakta         gene annotation
HUMAnN                 pathway and functional profiling
```

---

## Diversity Metrics

Alpha diversity: diversity within one sample.

Beta diversity: difference between samples.

```python
import numpy as np

def shannon(counts: np.ndarray) -> float:
    counts = counts[counts > 0]
    p = counts / counts.sum()
    return float(-(p * np.log(p)).sum())

def bray_curtis(a: np.ndarray, b: np.ndarray) -> float:
    return float(np.abs(a - b).sum() / (a + b).sum())
```

Interpret carefully: sequencing depth, contamination, diet, geography, medication, and batch effects can dominate biological signal.

---

## Compositional Data Problem

Microbiome abundance data is compositional: counts are relative, not absolute. If one species rises, others appear to fall even if their absolute amount is unchanged.

```python
def centered_log_ratio(counts: np.ndarray, pseudocount: float = 1.0) -> np.ndarray:
    x = counts + pseudocount
    log_x = np.log(x)
    return log_x - log_x.mean()
```

Use compositional-aware methods when possible:

```
ALDEx2, ANCOM-BC, songbird, q2-composition
```

---

## Functional Profiling

Taxonomy says who is there. Functional profiling asks what they can do.

```python
functional_questions = [
    "Which metabolic pathways are enriched?",
    "Are antibiotic resistance genes present?",
    "Which carbohydrate-active enzymes exist?",
    "Can the community produce short-chain fatty acids?",
]
```

Functional potential is not the same as active expression. For activity, integrate metatranscriptomics, metabolomics, or proteomics.

---

## ML for Microbiome Data

```python
from sklearn.ensemble import RandomForestClassifier
from sklearn.model_selection import cross_val_score

def train_microbiome_classifier(X, y):
    """
    X: samples x taxa/pathways after compositional transform
    y: phenotype labels
    """
    model = RandomForestClassifier(n_estimators=500, random_state=42)
    auc = cross_val_score(model, X, y, cv=5, scoring="roc_auc")
    model.fit(X, y)
    return model, auc.mean()
```

Guard against leakage: samples from the same subject, extraction batch, or sequencing run can create fake performance.

---

## Key Takeaways

- 16S is cheaper and narrower; shotgun is richer and more expensive.
- Microbiome data is sparse, noisy, compositional, and batch-sensitive.
- Taxonomy, function, and activity are different layers.
- Good metadata often matters as much as the sequencing itself.

