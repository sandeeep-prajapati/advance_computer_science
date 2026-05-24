# Chapter 6: Gene Expression & Transcriptomics

## Gene Expression as Logging/Monitoring

Every cell runs the same "source code" (genome) but has completely different expression profiles:

```
Neuron:   high NEFL, SYP, CHAT        (neuronal proteins)
Liver:    high ALB, CYP3A4, APOB      (albumin, metabolism)
Muscle:   high MYH7, ACTN2, TNNT2    (contractile proteins)

Same genome → different program execution based on:
  - Transcription factors (config files)
  - Epigenetic marks (environment variables)
  - Signal transduction (runtime messages)
```

---

## RNA-seq Pipeline: End to End

```python
import subprocess
import pandas as pd
import numpy as np

def full_rnaseq_pipeline(sample_sheet: str, reference_genome: str,
                          gtf_annotation: str, output_dir: str) -> None:
    """
    Complete RNA-seq pipeline from FASTQ to count matrix.
    Tools: STAR (alignment), featureCounts (quantification)
    """
    import os
    samples = pd.read_csv(sample_sheet)  # columns: sample_name, fastq_r1, fastq_r2
    
    # Step 1: Build STAR index (one-time)
    star_index = f"{output_dir}/star_index"
    os.makedirs(star_index, exist_ok=True)
    subprocess.run([
        "STAR", "--runMode", "genomeGenerate",
        "--genomeDir", star_index,
        "--genomeFastaFiles", reference_genome,
        "--sjdbGTFfile", gtf_annotation,
        "--runThreadN", "8"
    ], check=True)
    
    bam_files = []
    for _, row in samples.iterrows():
        sample_out = f"{output_dir}/{row['sample_name']}"
        os.makedirs(sample_out, exist_ok=True)
        
        # Step 2: Align
        subprocess.run([
            "STAR", "--runMode", "alignReads",
            "--genomeDir", star_index,
            "--readFilesIn", row['fastq_r1'], row['fastq_r2'],
            "--readFilesCommand", "zcat",
            "--outSAMtype", "BAM", "SortedByCoordinate",
            "--outFileNamePrefix", f"{sample_out}/",
            "--runThreadN", "8",
        ], check=True)
        bam_files.append(f"{sample_out}/Aligned.sortedByCoord.out.bam")
    
    # Step 3: Count reads per gene
    subprocess.run([
        "featureCounts", "-T", "8", "-p",  # paired-end
        "-a", gtf_annotation,
        "-o", f"{output_dir}/counts.txt",
    ] + bam_files, check=True)
```

---

## Working with Count Matrices

```python
import pandas as pd
import numpy as np

def load_featurecounts(counts_file: str) -> pd.DataFrame:
    """Parse featureCounts output into clean count matrix."""
    df = pd.read_csv(counts_file, sep='\t', comment='#', index_col=0)
    # featureCounts adds metadata columns: Chr, Start, End, Strand, Length
    count_cols = [c for c in df.columns if c.endswith('.bam')]
    counts = df[count_cols].copy()
    counts.columns = [c.split('/')[-1].replace('.bam', '') for c in count_cols]
    return counts

def filter_low_expression(counts: pd.DataFrame,
                           min_count: int = 10,
                           min_samples: int = 2) -> pd.DataFrame:
    """Remove genes with very low counts across samples."""
    keep = (counts >= min_count).sum(axis=1) >= min_samples
    print(f"Keeping {keep.sum()} / {len(counts)} genes after filtering")
    return counts[keep]

# CPM and TPM normalization
def normalize_cpm(counts: pd.DataFrame) -> pd.DataFrame:
    """Counts Per Million — simple library size correction."""
    return counts / counts.sum() * 1e6

def normalize_tpm(counts: pd.DataFrame, gene_lengths: pd.Series) -> pd.DataFrame:
    """Transcripts Per Million — length and depth corrected."""
    rpk = counts.div(gene_lengths / 1000, axis=0)
    return rpk / rpk.sum() * 1e6
```

---

## Differential Expression with PyDESeq2

DESeq2 is the gold standard for RNA-seq differential expression:

```python
from pydeseq2.dds import DeseqDataSet
from pydeseq2.ds import DeseqStats
import pandas as pd

def run_deseq2(counts: pd.DataFrame, metadata: pd.DataFrame,
               design_factor: str = 'condition',
               reference_level: str = 'normal') -> pd.DataFrame:
    """
    Run DESeq2 differential expression analysis.
    
    counts:   genes × samples (integer counts, no normalization)
    metadata: samples × factors (must include design_factor column)
    """
    # PyDESeq2 expects samples × genes
    counts_T = counts.T
    
    # Build DESeq dataset
    dds = DeseqDataSet(
        counts=counts_T,
        metadata=metadata,
        design_factors=design_factor,
        ref_level=[design_factor, reference_level],
        n_cpus=4,
    )
    
    # Estimate size factors + dispersions + fit model
    dds.deseq2()
    
    # Statistical test
    stat_res = DeseqStats(dds, contrast=[design_factor, 'tumor', reference_level])
    stat_res.summary()
    
    results = stat_res.results_df
    results = results.sort_values('padj')
    return results

def volcano_data(results: pd.DataFrame, fc_threshold: float = 1.0,
                 padj_threshold: float = 0.05) -> pd.DataFrame:
    """Annotate results for volcano plot."""
    results = results.copy()
    results['neg_log10_padj'] = -np.log10(results['padj'].clip(lower=1e-300))
    results['category'] = 'Not significant'
    up = (results['log2FoldChange'] >= fc_threshold) & (results['padj'] < padj_threshold)
    down = (results['log2FoldChange'] <= -fc_threshold) & (results['padj'] < padj_threshold)
    results.loc[up, 'category'] = 'Upregulated'
    results.loc[down, 'category'] = 'Downregulated'
    return results

# Pathway enrichment analysis
def gene_set_enrichment(gene_list: list, gene_sets: dict,
                        background_size: int = 20000) -> pd.DataFrame:
    """
    Fisher's exact test for pathway enrichment.
    gene_sets: {'pathway_name': ['GENE1', 'GENE2', ...], ...}
    """
    from scipy.stats import fisher_exact
    
    results = []
    query_set = set(gene_list)
    
    for pathway, genes in gene_sets.items():
        pathway_set = set(genes)
        overlap = len(query_set & pathway_set)
        
        # 2x2 contingency table
        a = overlap
        b = len(query_set) - overlap
        c = len(pathway_set) - overlap
        d = background_size - len(query_set) - len(pathway_set) + overlap
        
        _, pval = fisher_exact([[a, b], [c, d]], alternative='greater')
        results.append({
            'pathway': pathway,
            'overlap': overlap,
            'pathway_size': len(pathway_set),
            'pvalue': pval,
            'gene_ratio': overlap / len(query_set),
        })
    
    df = pd.DataFrame(results)
    from scipy.stats import false_discovery_control
    df['padj'] = false_discovery_control(df['pvalue'])
    return df.sort_values('padj')
```

---

## Single-Cell RNA-seq (scRNA-seq)

scRNA-seq measures expression in individual cells (not bulk averages):

```python
import scanpy as sc
import anndata as ad
import numpy as np
import pandas as pd

def process_scrna_seq(h5_file: str, min_genes: int = 200,
                       min_cells: int = 3) -> sc.AnnData:
    """
    Process 10x Genomics scRNA-seq data with Scanpy.
    
    Returns AnnData object: cells × genes matrix + metadata
    """
    # Load 10x output
    adata = sc.read_10x_h5(h5_file)
    adata.var_names_make_unique()
    
    # ── Quality Control ──────────────────────────────
    # Mitochondrial gene fraction = proxy for dead/damaged cells
    adata.var['mt'] = adata.var_names.str.startswith('MT-')
    sc.pp.calculate_qc_metrics(adata, qc_vars=['mt'], percent_top=None, inplace=True)
    
    # Filter cells: too few genes = empty/dead, too many = doublets
    sc.pp.filter_cells(adata, min_genes=min_genes)
    sc.pp.filter_genes(adata, min_cells=min_cells)
    adata = adata[adata.obs['pct_counts_mt'] < 20].copy()  # <20% mitochondrial
    
    print(f"After QC: {adata.n_obs} cells × {adata.n_vars} genes")
    
    # ── Normalization ────────────────────────────────
    sc.pp.normalize_total(adata, target_sum=1e4)  # 10k counts per cell
    sc.pp.log1p(adata)                             # log(1 + x)
    
    # ── Feature selection ────────────────────────────
    sc.pp.highly_variable_genes(adata, min_mean=0.0125, max_mean=3,
                                 min_disp=0.5, n_top_genes=2000)
    adata.raw = adata  # save normalized counts
    adata = adata[:, adata.var.highly_variable]
    
    # ── Dimensionality reduction ─────────────────────
    sc.pp.scale(adata, max_value=10)
    sc.tl.pca(adata, svd_solver='arpack', n_comps=50)
    
    # ── Neighborhood graph ───────────────────────────
    sc.pp.neighbors(adata, n_neighbors=15, n_pcs=40)
    
    # ── Clustering ───────────────────────────────────
    sc.tl.leiden(adata, resolution=0.5)  # Leiden community detection
    
    # ── Visualization ────────────────────────────────
    sc.tl.umap(adata)
    
    return adata

def find_marker_genes(adata: sc.AnnData) -> pd.DataFrame:
    """Find genes that define each cluster."""
    sc.tl.rank_genes_groups(adata, 'leiden', method='wilcoxon')
    markers = sc.get.rank_genes_groups_df(adata, group=None)
    return markers

def annotate_cell_types(adata: sc.AnnData, marker_genes: dict) -> sc.AnnData:
    """
    Score each cell against known marker gene lists.
    marker_genes: {'T_cell': ['CD3D', 'CD3E', 'TRAC'], 'B_cell': ['MS4A1', ...]}
    """
    for cell_type, markers in marker_genes.items():
        sc.tl.score_genes(adata, gene_list=markers, score_name=f'{cell_type}_score')
    
    # Assign most likely cell type
    score_cols = [f'{ct}_score' for ct in marker_genes]
    adata.obs['cell_type'] = adata.obs[score_cols].idxmax(axis=1).str.replace('_score', '')
    return adata
```

---

## Pseudotime and Trajectory Analysis

Track developmental progression within scRNA-seq data:

```python
import scvelo as scv  # pip install scvelo

def rna_velocity_analysis(adata: sc.AnnData, loom_file: str) -> sc.AnnData:
    """
    RNA velocity: predict future cell state from spliced/unspliced RNA.
    Requires loom file with spliced/unspliced counts (from velocyto).
    """
    # Load velocity data
    velocity_adata = scv.read(loom_file, cache=True)
    adata = scv.utils.merge(adata, velocity_adata)
    
    # Preprocessing
    scv.pp.filter_and_normalize(adata, min_shared_counts=20, n_top_genes=2000)
    scv.pp.moments(adata, n_pcs=30, n_neighbors=30)
    
    # Fit velocity model
    scv.tl.velocity(adata, mode='stochastic')  # or 'dynamical' for more accuracy
    scv.tl.velocity_graph(adata)
    
    return adata

# PAGA for trajectory inference (no velocity needed)
def trajectory_analysis(adata: sc.AnnData, root_cluster: str) -> sc.AnnData:
    """PAGA + diffusion pseudotime for developmental trajectory."""
    sc.tl.paga(adata, groups='leiden')
    sc.pl.paga(adata)
    
    # Diffusion pseudotime from root
    adata.uns['iroot'] = np.flatnonzero(adata.obs['leiden'] == root_cluster)[0]
    sc.tl.dpt(adata)
    return adata
```

---

## Gene Regulatory Networks

Model which transcription factors drive which genes:

```python
import pyscenic  # pip install pyscenic
from arboreto.algo import grnboost2  # GRN inference

def infer_grn(expression_matrix: pd.DataFrame,
              tf_list: list) -> pd.DataFrame:
    """
    Infer Gene Regulatory Network with GRNBoost2 (gradient boosting).
    
    expression_matrix: cells × genes
    tf_list:          list of transcription factor gene names
    Returns:          DataFrame of TF → target gene regulatory links
    """
    network = grnboost2(
        expression_data=expression_matrix,
        gene_names=expression_matrix.columns.tolist(),
        tf_names=tf_list,
        verbose=True
    )
    return network.sort_values('importance', ascending=False)

def scenic_regulon_analysis(expression_matrix: pd.DataFrame,
                             network: pd.DataFrame,
                             motif_annotations: str) -> dict:
    """
    SCENIC: from GRN to transcription factor regulons.
    Validates TF targets by checking motif enrichment.
    """
    import json
    from pyscenic.prune import prune2df
    from pyscenic.aucell import aucell
    
    # Step 1: Prune indirect links using motif databases
    regulons = prune2df(network, motif_annotations, num_workers=4)
    
    # Step 2: Score each cell for each regulon activity (AUCell)
    auc_matrix = aucell(expression_matrix, regulons)
    
    return {'regulons': regulons, 'auc_scores': auc_matrix}
```

---

## Bulk vs Single-Cell: When to Use Each

```
Bulk RNA-seq:
  ✓ Cheap ($50-100/sample)
  ✓ Well-established analysis (DESeq2, edgeR)
  ✓ Great for comparing conditions (treated vs control)
  ✗ Averages across all cells (hides heterogeneity)
  Use: well-defined cell populations, large cohorts

Single-cell RNA-seq:
  ✓ Resolves cell-to-cell variability
  ✓ Discovers new cell types/states
  ✓ Trajectory/developmental analysis
  ✗ Expensive ($1000-5000/sample)
  ✗ Higher technical noise
  Use: development, tumor heterogeneity, immune profiling

Spatial transcriptomics (Visium, Slide-seq):
  ✓ Expression + spatial location in tissue
  ✓ Understand cell-cell interactions in tissue context
  ✗ Lower throughput, lower sensitivity
  Use: tissue architecture questions
```

---

## AI/ML Applications in Transcriptomics

```python
import torch
import torch.nn as nn
from sklearn.preprocessing import LabelEncoder

class ExpressionClassifier(nn.Module):
    """Classify tissue type or disease state from RNA-seq profile."""
    
    def __init__(self, n_genes: int, n_classes: int, hidden: int = 512):
        super().__init__()
        self.net = nn.Sequential(
            nn.Linear(n_genes, hidden),
            nn.BatchNorm1d(hidden),
            nn.ReLU(),
            nn.Dropout(0.4),
            nn.Linear(hidden, hidden // 2),
            nn.BatchNorm1d(hidden // 2),
            nn.ReLU(),
            nn.Dropout(0.3),
            nn.Linear(hidden // 2, n_classes)
        )
    
    def forward(self, x: torch.Tensor) -> torch.Tensor:
        return self.net(x)

# Variational Autoencoder for single-cell data
class scVAE(nn.Module):
    """VAE for learning continuous latent space of cell states."""
    
    def __init__(self, n_genes: int, latent_dim: int = 10):
        super().__init__()
        self.encoder = nn.Sequential(
            nn.Linear(n_genes, 128), nn.ReLU(),
            nn.Linear(128, 64), nn.ReLU(),
        )
        self.fc_mu = nn.Linear(64, latent_dim)
        self.fc_logvar = nn.Linear(64, latent_dim)
        self.decoder = nn.Sequential(
            nn.Linear(latent_dim, 64), nn.ReLU(),
            nn.Linear(64, 128), nn.ReLU(),
            nn.Linear(128, n_genes), nn.Softplus()  # non-negative counts
        )
    
    def encode(self, x):
        h = self.encoder(x)
        return self.fc_mu(h), self.fc_logvar(h)
    
    def reparameterize(self, mu, logvar):
        std = torch.exp(0.5 * logvar)
        return mu + torch.randn_like(std) * std
    
    def forward(self, x):
        mu, logvar = self.encode(x)
        z = self.reparameterize(mu, logvar)
        return self.decoder(z), mu, logvar

def vae_loss(recon_x, x, mu, logvar):
    """ELBO loss: reconstruction + KL divergence."""
    recon = nn.functional.mse_loss(recon_x, x, reduction='sum')
    kld = -0.5 * torch.sum(1 + logvar - mu.pow(2) - logvar.exp())
    return recon + kld
```

---

## Key Tools Summary

```
Bulk RNA-seq:
  Alignment:   STAR, HISAT2, Salmon (pseudo-alignment)
  Counting:    featureCounts, HTSeq
  DE analysis: DESeq2 (R), PyDESeq2 (Python), edgeR, limma
  Pathway:     clusterProfiler, GSEA, gseapy

Single-cell:
  Processing:  Scanpy (Python), Seurat (R), Cell Ranger (10x)
  Clustering:  Leiden, Louvain
  Trajectory:  scVelo (velocity), Monocle3, PAGA
  Integration: Harmony, scVI (VAE-based), Scanorama
  Cell typing: CellTypist, SingleR

Databases:
  Expression:  GEO (Gene Expression Omnibus), GTEx, ENCODE
  Pathways:    KEGG, Reactome, MSigDB, Gene Ontology
  Single-cell: Human Cell Atlas, cellxgene
```
