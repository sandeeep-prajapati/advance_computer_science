# Chapter 15: Single-Cell Multi-omics

## Measuring Biology Cell by Cell

Bulk sequencing averages many cells together. Single-cell methods preserve cell-to-cell variation.

```
Bulk RNA-seq:        one expression vector per sample
scRNA-seq:           one expression vector per cell
single-cell multiome: RNA + chromatin accessibility in the same cell
spatial omics:       molecular data plus tissue coordinates
```

---

## Data Model

Single-cell data is usually a sparse matrix:

```python
import numpy as np
from scipy import sparse

# cells x genes
X = sparse.csr_matrix(np.array([
    [0, 3, 0, 1],
    [2, 0, 0, 4],
    [0, 0, 5, 1],
]))
```

In Python, `AnnData` is the common container:

```python
import anndata as ad

adata = ad.AnnData(X)
adata.obs["sample"] = ["A", "A", "B"]
adata.var["gene_symbol"] = ["GATA1", "MS4A1", "COL1A1", "ACTB"]
```

---

## scRNA-seq Pipeline

```python
import scanpy as sc

def basic_scanpy_pipeline(adata):
    sc.pp.filter_cells(adata, min_genes=200)
    sc.pp.filter_genes(adata, min_cells=3)
    sc.pp.normalize_total(adata, target_sum=1e4)
    sc.pp.log1p(adata)
    sc.pp.highly_variable_genes(adata, n_top_genes=2000)
    adata = adata[:, adata.var.highly_variable].copy()
    sc.pp.scale(adata, max_value=10)
    sc.tl.pca(adata)
    sc.pp.neighbors(adata)
    sc.tl.umap(adata)
    sc.tl.leiden(adata)
    return adata
```

Typical outputs:

```
clusters
marker genes
cell type annotations
differential expression
trajectories
batch-corrected embeddings
```

---

## Single-Cell ATAC-seq

scATAC-seq measures open chromatin in each cell.

```
cells x peaks matrix
peak calling
TF motif enrichment
gene activity scores
integration with scRNA-seq
```

Tools:

```
Signac        scATAC analysis in R
ArchR         scalable scATAC and multiome
muon          multimodal omics in Python
SnapATAC      large-scale scATAC
```

---

## Multiome Integration

Multiome combines modalities:

```python
modalities = {
    "RNA": "what genes are expressed",
    "ATAC": "which regulatory regions are open",
    "protein": "surface markers or selected proteins",
    "spatial": "where cells are in tissue",
}
```

Integration goal:

```
Find a shared representation where cells with similar biological state cluster together,
while modality-specific features remain interpretable.
```

Common methods:

```
weighted nearest neighbors
canonical correlation analysis
matrix factorization
variational autoencoders
totalVI / multiVI / scVI-tools
```

---

## Batch Effects

Single-cell data is especially vulnerable to batch effects.

```python
batch_sources = [
    "donor",
    "sample handling",
    "library kit",
    "sequencing lane",
    "chemistry version",
    "ambient RNA",
    "doublets",
]
```

Good correction removes technical variation while preserving real biology. Overcorrection can erase meaningful differences.

---

## Spatial Transcriptomics

Spatial methods preserve tissue location:

```
spot/cell coordinate + expression vector
```

Use cases:

```
tumor microenvironment mapping
developmental patterning
brain region analysis
cell-cell interaction inference
pathology augmentation
```

---

## Key Takeaways

- Single-cell analysis turns biology into sparse high-dimensional matrix work.
- QC decisions strongly shape downstream interpretation.
- Multiome links regulation to expression in the same biological state.
- Spatial context adds the missing "where" to molecular measurements.

