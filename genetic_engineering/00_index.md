# Mastering Genetic Engineering with CS & AI Skills
## A Computer Science / AI Practitioner's Guide to Computational Biology

---

## Why This Course Exists

You have deep CS and AI skills. Genetic engineering is experiencing a **software moment**: the tools, data, and problems are increasingly computational. This course bridges your existing expertise with the biology you need to solve real-world problems in biotech, drug discovery, and synthetic biology.

```
Your CS/AI Skills → Genetic Engineering Problems

Sequence Analysis     → DNA/protein sequences as strings (NLP analogy)
Graph Algorithms      → Metabolic networks, gene regulatory networks
Machine Learning      → Protein structure prediction, drug discovery
Optimization          → DNA codon optimization, primer design
Database Systems      → Genomic databases, BLAST, variant calling
Statistics            → GWAS, differential expression, QTL mapping
Computer Vision       → Cell microscopy analysis, pathology
NLP/Transformers      → AlphaFold2, ESMFold, DNABERT, Evo
Reinforcement Learning → Protein design, CRISPR guide RNA optimization
Software Engineering  → Bioinformatics pipelines, workflow orchestration
```

---

## Table of Contents

| Chapter | Title | Key Concepts |
|---|---|---|
| 01 | [Biology for CS Programmers](01_biology_for_cs.md) | DNA, RNA, proteins, cells — from a data perspective |
| 02 | [Genomics and Sequencing](02_genomics_and_sequencing.md) | NGS, FASTQ, BAM, variant calling, WGS/WES/RNA-seq |
| 03 | [Sequence Analysis & Bioinformatics](03_sequence_analysis.md) | BLAST, alignment, Smith-Waterman, CLUSTAL |
| 04 | [CRISPR-Cas9 Gene Editing](04_crispr.md) | Mechanism, guide RNA design, off-targets, delivery |
| 05 | [Protein Structure and AlphaFold](05_protein_structure_alphafold.md) | Folding problem, AlphaFold2, ESMFold, structure-function |
| 06 | [Gene Expression & Transcriptomics](06_gene_expression_transcriptomics.md) | RNA-seq, differential expression, DESeq2, scRNA-seq |
| 07 | [Synthetic Biology](07_synthetic_biology.md) | Parts, circuits, BioBricks, metabolic engineering |
| 08 | [Drug Discovery with AI](08_drug_discovery_ai.md) | Virtual screening, QSAR, generative models, AlphaFold |
| 09 | [Epigenetics and Regulation](09_epigenetics_regulation.md) | DNA methylation, ChIP-seq, regulatory networks |
| 10 | [Population Genetics & GWAS](10_population_genetics_gwas.md) | SNPs, GWAS, polygenic scores, ancestry |
| 11 | [Cancer Genomics](11_cancer_genomics.md) | Mutations, tumor evolution, liquid biopsy, immunotherapy |
| 12 | [Protein Engineering with AI](12_protein_engineering_ai.md) | Directed evolution, protein language models, RFdiffusion |
| 13 | [Metagenomics & Microbiome](13_metagenomics_microbiome.md) | 16S rRNA, shotgun sequencing, microbiome analysis |
| 14 | [CRISPR Screening](14_crispr_screening.md) | Pooled screens, MAGeCK algorithm, genome-wide screens |
| 15 | [Single-Cell Multi-omics](15_single_cell_multiomics.md) | scRNA-seq, ATAC-seq, spatial transcriptomics, Seurat |
| 16 | [Bioinformatics Pipelines](16_bioinformatics_pipelines.md) | Nextflow, Snakemake, WDL, cloud genomics (AWS/GCP) |
| 17 | [AI for Genomics](17_ai_for_genomics.md) | Deep learning, graph neural nets, transformers for DNA/protein |
| 18 | [Ethics, Safety, and Regulation](18_ethics_safety_regulation.md) | Biosafety, dual-use research, FDA/EMA, gene therapy ethics |
| 19 | [Solving Real Problems](19_solving_real_problems.md) | Case studies: COVID-19, sickle cell, CAR-T, insulin |
| 20 | [Learning Roadmap](20_learning_roadmap.md) | 6-week plan, projects, career paths in computational biology |

---

## Architecture: How Genetic Engineering Works Computationally

```
Biological Central Dogma
═══════════════════════════════════════════════════════════
DNA → RNA → Protein → Function

As data structures:
  DNA:     String of {A, T, G, C}    length: ~3 billion for human
  RNA:     String of {A, U, G, C}    (subset of genes, ~20,000 types)
  Protein: String of 20 amino acids  (folded 3D structure = function)


Genetic Engineering Pipeline (CS/AI view)
═══════════════════════════════════════════════════════════

Problem Definition              Analysis                  Solution
┌─────────────────┐       ┌─────────────────┐      ┌─────────────────┐
│ What gene/protein│  →   │ Sequence,        │  →  │ Edit genome:     │
│ to target?       │      │ structure, model  │     │ CRISPR, gene     │
│ What disease?    │      │ (AI/ML/align.)   │     │ therapy, drug    │
│ What function?   │      └─────────────────┘      └─────────────────┘
└─────────────────┘
           │
           ▼
  AI-Powered Tools:
  AlphaFold2    → Predict protein structure from sequence
  DNABERT/Evo  → Understand DNA regulatory grammar
  ProteinMPNN   → Design protein sequences
  RFdiffusion   → Generate novel protein backbones
  ChatNT        → Genomic sequence foundation model
  ESM-2         → Protein language model (650M params)


Data Scales
═══════════════════════════════════════════════════════════
Human genome:   3.2 billion base pairs × 2 = 6.4 GB raw
Proteome:       ~20,000 proteins
Global sequencing: ~10 exabytes/year generated
PDB:            230,000+ protein structures
UniProt:        250 million+ protein sequences
```

---

## Your CS/AI Advantage

```
Domain                 CS/AI Approach
─────────────────────────────────────────────────────────
Protein structures    Computer vision (3D CNN, GNN)
DNA sequences         NLP (BERT-style transformers, LLMs)
Gene networks         Graph neural networks
Drug design           Generative models (VAE, GAN, diffusion)
Cell images           Computer vision (U-Net segmentation)
Clinical trials       Causal inference, survival analysis
Genome-wide assoc.    Statistical ML, polygenic models
Protein design        RL + protein language models
CRISPR guide RNA      Sequence-to-function ML models
Pipeline automation   Workflow engines, cloud computing
```

---

## Prerequisites

- **Required:** Python programming (intermediate level)
- **Helpful:** Basic statistics, probability
- **No biology required** — Chapter 1 covers everything you need
- **Tools:** Python 3.10+, BioPython, pandas, scikit-learn, PyTorch/JAX

```bash
# Setup your environment
conda create -n bioml python=3.10
conda activate bioml
pip install biopython pandas numpy scikit-learn torch
pip install scanpy anndata seurat-py
pip install nextflow  # pipeline tool
pip install pydeseq2  # differential expression
pip install fair-esm  # Meta's ESM protein models
```
