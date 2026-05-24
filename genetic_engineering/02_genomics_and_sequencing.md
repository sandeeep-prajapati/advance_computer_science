# Chapter 2: Genomics and Sequencing

## DNA Sequencing: Reading the Code

Sequencing is the process of determining the exact order of bases (A, T, G, C) in a DNA molecule.

```
First-generation (Sanger, 1977): sequence ~1 kb, one fragment at a time
Second-gen (NGS, 2007+):       sequence billions of fragments in parallel
Third-generation (Oxford Nanopore, PacBio): long reads, real-time
```

---

## Next-Generation Sequencing (NGS)

```
NGS Pipeline:
  DNA/RNA sample
      ↓
  Library preparation (fragment DNA, add adapters)
      ↓
  Sequencing (read millions of short fragments simultaneously)
      ↓
  FASTQ files (raw reads + quality scores)
      ↓
  Alignment to reference genome → BAM/CRAM files
      ↓
  Variant calling → VCF files
      ↓
  Analysis (differential expression, GWAS, etc.)
```

---

## FASTQ Format

The raw output of sequencing machines:

```
@SRR1551047.1 1 length=101
TGGCTTTACACTTTATGCTTCGGCTCGATAGTGTTTATTTCCGAGCCAGGTCTCGTTCGTTATCGGAATTTCGTATTGCAGCACAAACAATTTCAGTCAT
+
CCCFFFFFHHHHHJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJIJJJJJJJJJJJJJJJJJJJJJJJJJJIJJJJJJJJJHHHGHGFFFFFCCC

Line 1: @READ_ID
Line 2: DNA sequence
Line 3: + (separator)
Line 4: Quality scores (ASCII encoding: Q = ord(char) - 33)

Quality score Q = -10 log₁₀(P_error)
  Q10 = 10% error probability
  Q20 = 1% error probability
  Q30 = 0.1% error probability  ← acceptable threshold
```

```python
from Bio import SeqIO
from Bio.SeqRecord import SeqRecord
import numpy as np

def parse_fastq(filename: str):
    """Parse FASTQ file and return reads with quality scores."""
    reads = []
    for record in SeqIO.parse(filename, "fastq"):
        seq = str(record.seq)
        qualities = record.letter_annotations["phred_quality"]
        reads.append({
            'id': record.id,
            'sequence': seq,
            'quality': qualities,
            'avg_quality': np.mean(qualities),
            'length': len(seq)
        })
    return reads

def quality_filter(reads: list, min_avg_quality: float = 20.0,
                   min_length: int = 50) -> list:
    """Filter reads by quality."""
    return [r for r in reads
            if r['avg_quality'] >= min_avg_quality
            and r['length'] >= min_length]

# Quality summary
def fastq_summary(filename: str) -> dict:
    reads = parse_fastq(filename)
    qualities = [r['avg_quality'] for r in reads]
    lengths = [r['length'] for r in reads]
    return {
        'total_reads': len(reads),
        'mean_quality': np.mean(qualities),
        'median_quality': np.median(qualities),
        'mean_length': np.mean(lengths),
        'total_bases': sum(lengths),
    }
```

---

## Sequence Alignment

Aligning short reads to a reference genome:

```
Problem: given millions of reads (~100 bp each) and a genome (~3 GB),
find where each read originated.

Tools: BWA-MEM, STAR (RNA-seq), HISAT2, Bowtie2

Algorithm insight:
  Naive: compare each read to every position → O(N·M) → too slow
  BWA:   Burrows-Wheeler Aligner using FM-index (compressed suffix array)
  STAR:  Maximal mappable prefix approach + splice-aware
```

```python
import subprocess
import os

# Run BWA alignment (command-line tool)
def align_reads(reference: str, reads_r1: str, reads_r2: str,
                output_bam: str, threads: int = 4) -> None:
    """
    Align paired-end reads to reference genome using BWA-MEM.
    Requires: bwa, samtools installed
    """
    # Index reference (one-time)
    if not os.path.exists(f"{reference}.bwt"):
        subprocess.run(["bwa", "index", reference], check=True)
    
    # Align and sort
    bwa_cmd = ["bwa", "mem", "-t", str(threads), reference, reads_r1, reads_r2]
    sort_cmd = ["samtools", "sort", "-o", output_bam]
    
    bwa = subprocess.Popen(bwa_cmd, stdout=subprocess.PIPE)
    subprocess.run(sort_cmd, stdin=bwa.stdout, check=True)
    bwa.wait()
    
    # Index BAM
    subprocess.run(["samtools", "index", output_bam], check=True)

# Parse BAM file
import pysam  # pip install pysam

def count_reads_per_gene(bam_file: str, gene_bed: str) -> dict:
    """Count reads overlapping gene regions."""
    bam = pysam.AlignmentFile(bam_file, "rb")
    gene_counts = {}
    
    with open(gene_bed) as f:
        for line in f:
            chrom, start, end, gene_name = line.strip().split('\t')[:4]
            count = bam.count(chrom, int(start), int(end))
            gene_counts[gene_name] = count
    
    bam.close()
    return gene_counts
```

---

## VCF: Variant Call Format

The standard format for genetic variants:

```
##fileformat=VCFv4.2
##reference=GRCh38
#CHROM  POS     ID       REF  ALT  QUAL  FILTER  INFO              FORMAT  SAMPLE1
chr17   43094077 rs80357713  G    A    999  PASS  AC=1;AF=0.5;DP=45  GT:DP:GQ  0/1:45:99
chr17   43094077 .           GT   G    200  PASS  INDEL;DP=30        GT        0/1

CHROM: chromosome
POS:   1-based genomic position
REF:   reference allele
ALT:   alternate allele
QUAL:  variant quality score
GT:    genotype (0/0=homozygous ref, 0/1=heterozygous, 1/1=homozygous alt)
```

```python
import vcfpy  # pip install vcfpy

def parse_variants(vcf_file: str) -> list:
    """Parse VCF and extract variant information."""
    variants = []
    reader = vcfpy.Reader.from_path(vcf_file)
    
    for record in reader:
        for alt in record.ALT:
            variant = {
                'chrom': record.CHROM,
                'pos': record.POS,
                'ref': record.REF,
                'alt': str(alt.value),
                'qual': record.QUAL,
                'filter': list(record.FILTER),
                'variant_type': classify_variant(record.REF, str(alt.value)),
            }
            variants.append(variant)
    
    return variants

def classify_variant(ref: str, alt: str) -> str:
    if len(ref) == len(alt) == 1:
        return 'SNV'
    elif len(ref) > len(alt):
        return 'DEL'
    elif len(ref) < len(alt):
        return 'INS'
    else:
        return 'MNV'

# Variant annotation with ANNOVAR or SnpEff
def annotate_variants(vcf_file: str, output_dir: str) -> None:
    """Annotate variants with gene, consequence, clinical significance."""
    subprocess.run([
        "snpEff", "hg38", "-v", vcf_file,
        "-s", f"{output_dir}/snpeff_summary.html"
    ], stdout=open(f"{output_dir}/annotated.vcf", 'w'), check=True)
```

---

## RNA-seq: Expression Profiling

```
Goal: Measure gene expression levels in a sample.

Pipeline:
  1. Extract RNA from cells
  2. Convert to cDNA (complementary DNA)
  3. Sequence (NGS)
  4. Align to genome or transcriptome
  5. Count reads per gene
  6. Normalize counts (TPM, FPKM, TMM)
  7. Differential expression analysis
```

```python
import pandas as pd
import numpy as np
from scipy import stats

# Raw counts matrix: rows=genes, cols=samples
counts = pd.DataFrame({
    'gene': ['BRCA1', 'TP53', 'VEGFA', 'ACTB'],
    'tumor_1': [120, 800, 2400, 15000],
    'tumor_2': [90, 920, 2100, 16000],
    'normal_1': [250, 750, 400, 14000],
    'normal_2': [280, 700, 380, 15500],
}).set_index('gene')

def normalize_tpm(counts: pd.DataFrame, gene_lengths: dict) -> pd.DataFrame:
    """
    TPM normalization (Transcripts Per Million).
    Accounts for gene length and sequencing depth.
    """
    # Reads per kilobase (RPK)
    rpk = counts.div([gene_lengths.get(g, 1000) for g in counts.index], axis=0) * 1000
    
    # Per million scaling
    tpm = rpk / rpk.sum() * 1e6
    return tpm

# Simple differential expression (t-test for demo; use DESeq2/limma for real work)
def differential_expression(counts: pd.DataFrame,
                             group1_cols: list, group2_cols: list) -> pd.DataFrame:
    """Simple t-test based DE (use PyDESeq2 in practice)."""
    results = []
    for gene in counts.index:
        g1 = np.log2(counts.loc[gene, group1_cols] + 1)
        g2 = np.log2(counts.loc[gene, group2_cols] + 1)
        log2fc = g1.mean() - g2.mean()
        _, pval = stats.ttest_ind(g1, g2)
        results.append({'gene': gene, 'log2FC': log2fc, 'pvalue': pval})
    
    df = pd.DataFrame(results)
    df['padj'] = stats.false_discovery_control(df['pvalue'].fillna(1))
    return df.sort_values('padj')

de_results = differential_expression(counts,
                                      ['tumor_1', 'tumor_2'],
                                      ['normal_1', 'normal_2'])
print(de_results[de_results['padj'] < 0.05])
```

---

## Whole Genome vs Whole Exome vs Targeted

| Sequencing Type | What's Covered | Cost | Use Case |
|---|---|---|---|
| WGS (Whole Genome) | 100% of genome | ~$300 | Structural variants, non-coding, research |
| WES (Whole Exome) | ~1% (all genes) | ~$100 | Clinical diagnosis, rare disease |
| Targeted panel | Specific genes | ~$50 | Cancer panels, pharmacogenomics |
| RNA-seq | Expressed genes | ~$100 | Differential expression, splicing |
| ATAC-seq | Open chromatin | ~$150 | Regulatory elements |
| ChIP-seq | Protein-DNA binding | ~$200 | Transcription factor binding |

---

## Long-Read Sequencing (Third Generation)

```python
# Oxford Nanopore: reads >50 kb, real-time, portable
# PacBio HiFi: reads ~15-20 kb, >99.9% accuracy

# Advantages over short reads:
# - Resolve repetitive regions (centromeres, telomeres)
# - Detect structural variants directly
# - Phased sequencing (distinguish maternal/paternal haplotypes)
# - Direct methylation detection (no bisulfite conversion needed)

# Use cases:
# - De novo genome assembly
# - Long-range phasing
# - Repeat expansion diseases (Huntington's, fragile X)
# - Direct RNA sequencing

# Quality metrics for sequencing runs:
def sequencing_qc(reads: list) -> dict:
    lengths = [r['length'] for r in reads]
    return {
        'total_reads': len(reads),
        'total_bases': sum(lengths),
        'n50': calculate_n50(lengths),       # read length where >50% bases in longer reads
        'mean_read_length': np.mean(lengths),
        'q20_bases_fraction': sum(1 for r in reads if r['avg_quality'] >= 20) / len(reads),
    }

def calculate_n50(lengths: list) -> int:
    sorted_lens = sorted(lengths, reverse=True)
    cumsum = 0
    total = sum(lengths)
    for l in sorted_lens:
        cumsum += l
        if cumsum >= total / 2:
            return l
    return 0
```

---

## Bioinformatics Tool Ecosystem

```
Alignment:    BWA, STAR, HISAT2, Minimap2 (long reads)
QC:           FastQC, MultiQC, Trimmomatic, fastp
Variant:      GATK, DeepVariant (Google's DL-based), FreeBayes
Expression:   featureCounts, HTSeq, Kallisto, Salmon
DE analysis:  DESeq2 (R), PyDESeq2 (Python), edgeR, limma
Assembly:     SPAdes, Flye (long read), Hifiasm
Annotation:   ANNOVAR, SnpEff, VEP (Ensembl)
Databases:    NCBI, Ensembl, UCSC Genome Browser, gnomAD

Google's AI contribution:
  DeepVariant: CNN-based variant caller, outperforms classical methods
  DeepConsensus: PacBio read correction
  AlphaFold: structure prediction (see Chapter 5)
  
Your AI skills apply here:
  Train CNN on read images (DeepVariant approach)
  Transformer on sequences (DNABERT, Nucleotide Transformer)
  Graph NN on regulatory networks
```
