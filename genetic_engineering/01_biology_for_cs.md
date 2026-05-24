# Chapter 1: Biology for CS Programmers

## The Central Dogma (Data Flow View)

Biology's central dogma is really a **data transformation pipeline**:

```
Storage    →    Transcription    →    Translation    →    Function
 DNA               mRNA                Protein
 HDD/SSD          RAM                  Running Process

DNA: blueprint (read-mostly, very stable)
RNA: temporary working copy (single-stranded, short-lived)
Protein: the actual worker (enzymes, structural, signaling)
```

---

## DNA: The String Data Structure

DNA is a 4-character alphabet string:
```python
# DNA alphabet
BASES = {'A', 'T', 'G', 'C'}

# A-T pair, G-C pair (complementary base pairing)
complement = {'A': 'T', 'T': 'A', 'G': 'C', 'C': 'G'}

# Double helix: two antiparallel strands
dna_strand_1 = "ATGCGTACGTAGC"
dna_strand_2 = "".join(complement[b] for b in reversed(dna_strand_1))
# dna_strand_2 = "GCTACGTACGCAT" (reverse complement)

# Human genome: ~3.2 billion base pairs
# Each base: 2 bits to encode (A=00, T=01, G=10, C=11)
# Human genome: ~800 MB uncompressed, ~150 MB compressed

# Chromosome: one long DNA molecule
# Humans: 46 chromosomes (23 pairs)
# Each chromosome: millions to hundreds of millions of base pairs
```

---

## Genes: Functions in the Codebase

```python
# Gene: a segment of DNA encoding a protein (or functional RNA)
# Human genes: ~20,000 (only ~2% of genome is protein-coding!)
# The other 98% = regulatory elements, introns, repeats, "dark matter"

class Gene:
    def __init__(self, name: str, chromosome: str, start: int, end: int, strand: str):
        self.name = name          # e.g., "BRCA1"
        self.chromosome = chromosome  # e.g., "chr17"
        self.start = start        # genomic position
        self.end = end
        self.strand = strand      # '+' or '-' (which strand it's on)
    
    @property
    def length(self) -> int:
        return self.end - self.start

# BRCA1 gene example
brca1 = Gene("BRCA1", "chr17", 43044295, 43170245, "-")
print(f"BRCA1 length: {brca1.length:,} bp")  # ~125,951 bp
```

---

## The Genetic Code: A 3-letter Encoding

```python
# Codon table: 3 DNA bases → 1 amino acid
# 4³ = 64 codons → 20 amino acids + 3 stop codons (redundant code)

CODON_TABLE = {
    'ATA': 'I', 'ATC': 'I', 'ATT': 'I', 'ATG': 'M',  # Start codon = M
    'ACA': 'T', 'ACC': 'T', 'ACG': 'T', 'ACT': 'T',
    'AAC': 'N', 'AAT': 'N', 'AAA': 'K', 'AAG': 'K',
    'AGC': 'S', 'AGT': 'S', 'AGA': 'R', 'AGG': 'R',
    'CTA': 'L', 'CTC': 'L', 'CTG': 'L', 'CTT': 'L',
    'CCA': 'P', 'CCC': 'P', 'CCG': 'P', 'CCT': 'P',
    'CAC': 'H', 'CAT': 'H', 'CAA': 'Q', 'CAG': 'Q',
    'CGA': 'R', 'CGC': 'R', 'CGG': 'R', 'CGT': 'R',
    'GTA': 'V', 'GTC': 'V', 'GTG': 'V', 'GTT': 'V',
    'GCA': 'A', 'GCC': 'A', 'GCG': 'A', 'GCT': 'A',
    'GAC': 'D', 'GAT': 'D', 'GAA': 'E', 'GAG': 'E',
    'GGA': 'G', 'GGC': 'G', 'GGG': 'G', 'GGT': 'G',
    'TCA': 'S', 'TCC': 'S', 'TCG': 'S', 'TCT': 'S',
    'TTC': 'F', 'TTT': 'F', 'TTA': 'L', 'TTG': 'L',
    'TAC': 'Y', 'TAT': 'Y', 'TAA': '*', 'TAG': '*',  # Stop codons
    'TGC': 'C', 'TGT': 'C', 'TGA': '*', 'TGG': 'W',
}

def translate(dna: str) -> str:
    """Translate DNA sequence to amino acid sequence."""
    protein = []
    for i in range(0, len(dna) - 2, 3):
        codon = dna[i:i+3]
        aa = CODON_TABLE.get(codon, '?')
        if aa == '*':
            break  # stop codon
        protein.append(aa)
    return ''.join(protein)

# Transcription + Translation pipeline
def central_dogma(dna: str) -> str:
    rna = dna.replace('T', 'U')          # DNA → RNA (T→U)
    coding_dna = rna.replace('U', 'T')   # work with DNA codons
    return translate(coding_dna)

# Example: insulin gene first exon
insulin_dna = "ATGGCCCTGTGGATGCGCCTCCTGCCCCTGCTGGCGCTGCTGGCCCTCTGGGGACCT"
print(translate(insulin_dna))  # MACWMRLLLPLLALLALWGP (insulin signal peptide)
```

---

## Gene Expression: Process Scheduling

```python
# Gene expression = how much of a protein is being made
# Controlled by: promoters, enhancers, transcription factors

# Think of it like:
# - Gene = process definition (like a class/function)
# - Gene expression = number of running instances
# - Transcription factor = program scheduler
# - Promoter = main() entry point

# Expression levels measured by RNA-seq:
# Low expression  = 0-10 TPM (transcripts per million)
# Medium          = 10-100 TPM
# High            = >100 TPM
# Housekeeping genes: expressed everywhere (e.g., ACTB, GAPDH)
# Tissue-specific:    expressed only in specific cell types

# Each cell has the same DNA but different expression profiles
# → Same source code, different runtime configurations
```

---

## Cells: The Runtime Environment

```
Eukaryotic cell components:
  Nucleus       → Storage (DNA)
  Ribosome      → CPU (translates mRNA to protein)
  Mitochondria  → Power supply (ATP generation)
  Endoplasmic   → Assembly line (protein folding, modification)
  Reticulum
  Golgi         → Shipping center (protein sorting and secretion)
  Cell membrane → API boundary (signals in/out)
  Cytoplasm     → RAM (active molecules)
```

---

## Mutations: Bugs in the Code

```python
# Types of DNA mutations (software bugs analogy):

# 1. SNP (Single Nucleotide Polymorphism) — single character substitution
dna_normal = "ATGCATGCATGC"
dna_snp    = "ATGCATGCATGC".replace("G", "A", 1)  # G→A at position 4
# Synonymous: same amino acid (no effect)
# Missense: different amino acid (may alter function)
# Nonsense: creates stop codon (truncates protein — often deleterious)

# 2. Insertion/Deletion (Indel) — insertion or deletion bug
dna_insertion = "ATGC" + "T" + "ATGCATGC"  # extra T inserted
dna_deletion  = "ATG" + "ATGCATGC"         # C deleted

# Frameshift: indel shifts reading frame (like corrupting all downstream bytes)
def show_frameshift(dna: str, indel_pos: int) -> None:
    normal   = translate(dna)
    with_del = translate(dna[:indel_pos] + dna[indel_pos+1:])  # delete 1 base
    print(f"Normal:    {normal}")
    print(f"Frameshift:{with_del}")  # completely different downstream

# 3. Copy Number Variation (CNV) — duplicated/deleted gene segments
# Gene amplification: BRCA1 duplicated 10× in some cancers

# 4. Chromosomal rearrangement — "swap" between files
# Philadelphia chromosome: BCR-ABL fusion → chronic myeloid leukemia
# Treated by imatinib (Gleevec) — discovered by computational drug design!
```

---

## The Protein: The Deployed Application

```python
from Bio import SeqIO
from Bio.Seq import Seq

# Amino acids: 20 characters
AA_ALPHABET = "ACDEFGHIKLMNPQRSTVWY"

# Protein structure hierarchy:
# Primary structure: amino acid sequence (1D string)
# Secondary structure: local folding patterns (α-helices, β-sheets)
# Tertiary structure: complete 3D fold (determines function)
# Quaternary structure: multiple protein chains assembled together

# Hemoglobin example: 4 chains (2α + 2β) — quaternary structure

# Key insight: SEQUENCE → STRUCTURE → FUNCTION
# AlphaFold2 predicts structure from sequence (with >90% accuracy!)

# Working with proteins in Python
from Bio.SeqUtils.ProtParam import ProteinAnalysis

insulin_protein = "MALWMRLLPLLALLALWGPDPAAAFVNQHLCGSHLVEALYLVCGERGFFYTPKTRREAEDLQVGQVELGGGPGAGSLQPLALEGSLQKRGIVEQCCTSICSLYQLENYCN"

analysis = ProteinAnalysis(insulin_protein)
print(f"Length: {len(insulin_protein)} aa")
print(f"Molecular weight: {analysis.molecular_weight():.0f} Da")
print(f"Isoelectric point: {analysis.isoelectric_point():.2f}")
print(f"Instability index: {analysis.instability_index():.2f}")
```

---

## Key Databases (Think: APIs and Datasets)

```python
# NCBI GenBank: DNA sequences
# UniProt: protein sequences and functions
# PDB: protein 3D structures
# Ensembl: annotated genomes
# KEGG: metabolic pathways
# OMIM: genetic diseases

# Accessing databases via BioPython
from Bio import Entrez, SeqIO

Entrez.email = "your.email@example.com"

# Fetch insulin gene sequence from NCBI
handle = Entrez.efetch(db="nucleotide", id="NM_000207",
                       rettype="gb", retmode="text")
record = SeqIO.read(handle, "genbank")
handle.close()

print(f"Gene: {record.name}")
print(f"Description: {record.description}")
print(f"Sequence length: {len(record.seq)} bp")
print(f"First 60 bp: {record.seq[:60]}")

# Fetch protein structure IDs
handle = Entrez.esearch(db="structure", term="BRCA1[gene]")
result = Entrez.read(handle)
print(f"Structures for BRCA1: {result['Count']}")
```

---

## Quick Reference: Biology for CS Developers

| Biology Term | CS Analog |
|---|---|
| Genome | Source code repository |
| Gene | Function/class definition |
| Chromosome | Source file |
| DNA | Hard disk (storage) |
| mRNA | RAM copy of code |
| Protein | Running process |
| Mutation | Bug / code change |
| Transcription factor | Scheduler/config |
| Promoter | Main entry point |
| Enhancer | Feature flag |
| Codon | 3-character opcode |
| Amino acid | Machine instruction |
| Protein folding | Compilation + linking |
| Cell | Operating system instance |
| Gene expression | Process priority |
| SNP | Variable substitution |
| Frameshift | Corrupted binary |
| CRISPR | Find-and-replace tool |
| Gene therapy | Patch/hotfix |
| RNA-seq | Profiling/logging |
| AlphaFold | Compile-time structure inference |
