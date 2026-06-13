# Chapter 20: Learning Roadmap

## How to Use This Course

You do not need to memorize all of biology first. Learn enough biology to ask better computational questions, then build small projects that force the concepts to become practical.

```
read -> implement -> analyze public data -> explain results -> repeat
```

---

## Six-Week Plan

### Week 1: Biology and Sequences

Read:

```
01_biology_for_cs.md
02_genomics_and_sequencing.md
03_sequence_analysis.md
```

Build:

```python
projects = [
    "FASTA/FASTQ parser",
    "reverse-complement and translation utilities",
    "Needleman-Wunsch or Smith-Waterman alignment",
    "k-mer counter and simple sequence classifier",
]
```

Goal: DNA/RNA/protein should feel like data structures you can manipulate.

### Week 2: Genomics Workflows

Read:

```
10_population_genetics_gwas.md
16_bioinformatics_pipelines.md
```

Build:

```python
projects = [
    "VCF parser and variant summary report",
    "small Snakemake or Nextflow QC pipeline",
    "Manhattan plot from public GWAS summary statistics",
]
```

Goal: understand the file formats and reproducibility discipline.

### Week 3: Expression and Single Cell

Read:

```
06_gene_expression_transcriptomics.md
15_single_cell_multiomics.md
```

Build:

```python
projects = [
    "differential expression notebook",
    "Scanpy clustering workflow on a public scRNA-seq dataset",
    "marker-gene based cell type annotation report",
]
```

Goal: turn count matrices into biological interpretation.

### Week 4: Proteins and Drug Discovery

Read:

```
05_protein_structure_alphafold.md
08_drug_discovery_ai.md
12_protein_engineering_ai.md
```

Build:

```python
projects = [
    "protein embedding similarity search",
    "small QSAR model using RDKit fingerprints",
    "variant ranking workflow with structure confidence notes",
]
```

Goal: connect sequence, structure, function, and molecules.

### Week 5: Editing and Synthetic Biology

Read:

```
04_crispr.md
07_synthetic_biology.md
14_crispr_screening.md
18_ethics_safety_regulation.md
```

Build:

```python
projects = [
    "guide RNA ranking toy model",
    "CRISPR screen count-analysis notebook",
    "codon optimization simulator",
    "risk checklist for a hypothetical therapeutic project",
]
```

Goal: understand perturbation design, screening statistics, and safety constraints.

### Week 6: Integration Project

Read:

```
09_epigenetics_regulation.md
11_cancer_genomics.md
13_metagenomics_microbiome.md
17_ai_for_genomics.md
19_solving_real_problems.md
```

Pick one capstone:

```python
capstones = [
    "Cancer variant annotation dashboard",
    "Microbiome classifier with leakage-safe validation",
    "Enhancer prediction model from DNA sequence",
    "End-to-end RNA-seq mini pipeline with report",
    "Protein property prediction model with uncertainty estimates",
]
```

Goal: finish one project that looks like a real portfolio artifact.

---

## Suggested Project Standards

Every project should include:

```
README with biological question
environment file
small reproducible dataset or download instructions
notebook or script
figures
limitations section
```

For ML projects, also include:

```
train/validation/test split explanation
baseline model
main metric
failure modes
data leakage check
```

---

## Career Paths

```python
paths = {
    "bioinformatics engineer": "pipelines, data processing, cloud/HPC",
    "computational biologist": "analysis, statistics, biological interpretation",
    "ML scientist in biotech": "models for sequence, structure, drugs, cells",
    "genomics data engineer": "large-scale storage, metadata, cohort systems",
    "clinical genomics analyst": "variant interpretation and reporting",
    "synthetic biology engineer": "design-build-test-learn systems",
}
```

The fastest path from CS into biotech is usually: Python + statistics + sequencing formats + one strong biological domain.

---

## Reference Stack

```bash
# Core Python
pip install numpy pandas scipy scikit-learn matplotlib seaborn jupyter

# Bioinformatics
pip install biopython pysam pybedtools pyranges anndata scanpy

# ML
pip install torch transformers lightning

# Chemistry / proteins
conda install -c conda-forge rdkit
pip install fair-esm biotite

# Pipelines
conda install -c bioconda snakemake nextflow
```

---

## Final Mental Model

```
Biology gives you the question.
Sequencing gives you the measurement.
Statistics tells you what is reliable.
Machine learning helps search huge spaces.
Engineering makes the work reproducible.
Ethics decides what should be built and how.
```

If you can move comfortably across those layers, your CS and AI background becomes a serious advantage in genetic engineering.

