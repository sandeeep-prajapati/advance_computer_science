# Chapter 3: Sequence Analysis & Bioinformatics Algorithms

## Sequence Analysis as String Algorithms

DNA and protein sequences ARE strings. Classic CS string algorithms apply directly.

```python
from Bio import SeqIO, pairwise2
from Bio.Seq import Seq
from Bio.SeqUtils import gc_fraction
import numpy as np

# Basic sequence operations
dna = Seq("ATGGCTATCGAAGGTCTGAAAGATGCAGGCTGGATCAAGCA")
print(f"Length: {len(dna)}")
print(f"GC content: {gc_fraction(dna)*100:.1f}%")
print(f"Complement: {dna.complement()}")
print(f"Reverse complement: {dna.reverse_complement()}")
print(f"Translation: {dna.translate()}")
```

---

## Pairwise Sequence Alignment

The **Smith-Waterman** (local) and **Needleman-Wunsch** (global) algorithms are classic dynamic programming.

### Needleman-Wunsch (Global Alignment):
```python
import numpy as np

def needleman_wunsch(seq1: str, seq2: str,
                     match: int = 2, mismatch: int = -1, gap: int = -2):
    """Global pairwise alignment via dynamic programming."""
    n, m = len(seq1), len(seq2)
    
    # Initialize DP matrix
    dp = np.zeros((n + 1, m + 1), dtype=int)
    dp[0, :] = [gap * j for j in range(m + 1)]
    dp[:, 0] = [gap * i for i in range(n + 1)]
    
    # Fill matrix
    for i in range(1, n + 1):
        for j in range(1, m + 1):
            diag = dp[i-1, j-1] + (match if seq1[i-1] == seq2[j-1] else mismatch)
            delete = dp[i-1, j] + gap
            insert = dp[i, j-1] + gap
            dp[i, j] = max(diag, delete, insert)
    
    # Traceback
    aligned1, aligned2 = [], []
    i, j = n, m
    while i > 0 or j > 0:
        if i > 0 and j > 0 and dp[i, j] == dp[i-1, j-1] + (match if seq1[i-1] == seq2[j-1] else mismatch):
            aligned1.append(seq1[i-1])
            aligned2.append(seq2[j-1])
            i -= 1; j -= 1
        elif i > 0 and dp[i, j] == dp[i-1, j] + gap:
            aligned1.append(seq1[i-1]); aligned2.append('-'); i -= 1
        else:
            aligned1.append('-'); aligned2.append(seq2[j-1]); j -= 1
    
    return ''.join(reversed(aligned1)), ''.join(reversed(aligned2)), dp[n, m]

seq1 = "ATCGCATGA"
seq2 = "ATCGATG"
a1, a2, score = needleman_wunsch(seq1, seq2)
print(f"Seq1: {a1}")
print(f"Seq2: {a2}")
print(f"Score: {score}")
```

### Smith-Waterman (Local Alignment):
```python
def smith_waterman(seq1: str, seq2: str,
                   match: int = 2, mismatch: int = -1, gap: int = -2):
    """Local pairwise alignment — finds best local match."""
    n, m = len(seq1), len(seq2)
    dp = np.zeros((n + 1, m + 1), dtype=int)
    max_score = 0
    max_pos = (0, 0)
    
    for i in range(1, n + 1):
        for j in range(1, m + 1):
            diag = dp[i-1, j-1] + (match if seq1[i-1] == seq2[j-1] else mismatch)
            dp[i, j] = max(0, diag, dp[i-1, j] + gap, dp[i, j-1] + gap)
            if dp[i, j] > max_score:
                max_score = dp[i, j]
                max_pos = (i, j)
    
    return max_score, max_pos
```

### Using BioPython for alignment:
```python
from Bio import pairwise2
from Bio.pairwise2 import format_alignment

seq1, seq2 = "GCATGCG", "GATTACA"
alignments = pairwise2.align.globalms(seq1, seq2, 2, -1, -2, -2)
print(format_alignment(*alignments[0]))
```

---

## BLAST: Heuristic Database Search

**BLAST** (Basic Local Alignment Search Tool) finds similar sequences in databases:

```python
from Bio.Blast import NCBIWWW, NCBIXML

def blast_sequence(query: str, database: str = "nt",
                   program: str = "blastn") -> list:
    """Search NCBI databases via BLAST."""
    result_handle = NCBIWWW.qblast(program, database, query)
    blast_records = NCBIXML.parse(result_handle)
    
    hits = []
    for blast_record in blast_records:
        for alignment in blast_record.alignments[:5]:  # top 5 hits
            for hsp in alignment.hsps:
                hits.append({
                    'title': alignment.title[:60],
                    'length': alignment.length,
                    'score': hsp.score,
                    'e_value': hsp.expect,
                    'identity': hsp.identities / hsp.align_length,
                    'query_seq': hsp.query[:50],
                    'subject_seq': hsp.sbjct[:50],
                })
    return hits

# BLAST algorithm (simplified — real BLAST uses k-mer lookup + gapped extension)
def blast_algorithm_overview():
    """
    BLAST steps:
    1. Hash query into k-mers (words, e.g., k=11 for DNA, k=3 for protein)
    2. Find word matches in database (lookup table)
    3. Extend matches in both directions (gapped extension)
    4. Score with substitution matrices (BLOSUM62 for protein, identity for DNA)
    5. Report hits with E-value (expected number of such hits by chance)
    
    Time complexity: O(query_length × database_size) with hash speedup
    Space complexity: O(k-mer_count) for hash table
    """
    pass
```

---

## Multiple Sequence Alignment

Aligning 3+ sequences simultaneously:

```python
from Bio import AlignIO, SeqIO
from Bio.Seq import Seq
from Bio.SeqRecord import SeqRecord
import subprocess
import tempfile

def multiple_alignment_mafft(sequences: dict) -> str:
    """Perform multiple sequence alignment using MAFFT."""
    # Write sequences to temp file
    with tempfile.NamedTemporaryFile(mode='w', suffix='.fasta', delete=False) as f:
        for name, seq in sequences.items():
            f.write(f">{name}\n{seq}\n")
        temp_input = f.name
    
    # Run MAFFT
    result = subprocess.run(
        ["mafft", "--auto", temp_input],
        capture_output=True, text=True
    )
    return result.stdout  # aligned FASTA

# Parse and analyze alignment
def analyze_alignment(aligned_fasta: str) -> dict:
    from io import StringIO
    alignment = AlignIO.read(StringIO(aligned_fasta), "fasta")
    
    n_seqs = len(alignment)
    align_length = alignment.get_alignment_length()
    
    # Conservation per position
    conservation = []
    for i in range(align_length):
        column = alignment[:, i]
        bases = [b for b in column if b != '-']
        if bases:
            most_common = max(set(bases), key=bases.count)
            cons_score = bases.count(most_common) / len(bases)
        else:
            cons_score = 0
        conservation.append(cons_score)
    
    return {
        'n_sequences': n_seqs,
        'alignment_length': align_length,
        'mean_conservation': np.mean(conservation),
        'conserved_positions': sum(1 for c in conservation if c > 0.9),
        'conservation_array': conservation,
    }
```

---

## Substitution Matrices

PAM and BLOSUM matrices encode evolutionary probabilities between amino acids:

```python
import numpy as np

# BLOSUM62 matrix (simplified — full version is 20×20)
BLOSUM62_PARTIAL = {
    ('A', 'A'): 4,  ('A', 'R'): -1, ('A', 'N'): -2, ('A', 'D'): -2,
    ('R', 'A'): -1, ('R', 'R'): 5,  ('R', 'N'): 0,  ('R', 'D'): -2,
    # ... (full 20×20 matrix)
}

def blosum62_score(seq1: str, seq2: str) -> int:
    """Score protein alignment using BLOSUM62."""
    from Bio.Align import substitution_matrices
    blosum62 = substitution_matrices.load("BLOSUM62")
    
    score = 0
    for a, b in zip(seq1, seq2):
        if a != '-' and b != '-':
            try:
                score += blosum62[a][b]
            except KeyError:
                score -= 4  # penalty for unknown
    return score

# Gap penalty models
def affine_gap_penalty(gap_length: int, open_penalty: float = 10,
                       extend_penalty: float = 0.5) -> float:
    """Affine gap penalty: open + extend × length."""
    if gap_length == 0:
        return 0
    return open_penalty + extend_penalty * gap_length
```

---

## K-mer Analysis

K-mers (substrings of length k) are the workhorse of modern bioinformatics:

```python
from collections import Counter
import hashlib

def kmer_spectrum(sequence: str, k: int = 21) -> Counter:
    """Count all k-mers in a sequence."""
    return Counter(sequence[i:i+k] for i in range(len(sequence) - k + 1))

def kmer_distance(seq1: str, seq2: str, k: int = 4) -> float:
    """K-mer based sequence distance (alignment-free)."""
    km1 = kmer_spectrum(seq1, k)
    km2 = kmer_spectrum(seq2, k)
    
    all_kmers = set(km1) | set(km2)
    
    # Cosine distance
    v1 = np.array([km1.get(kmer, 0) for kmer in all_kmers])
    v2 = np.array([km2.get(kmer, 0) for kmer in all_kmers])
    
    cosine_sim = np.dot(v1, v2) / (np.linalg.norm(v1) * np.linalg.norm(v2) + 1e-10)
    return 1 - cosine_sim

# MinHash for fast sequence similarity
def minhash_similarity(seq1: str, seq2: str, k: int = 21, n_hashes: int = 128) -> float:
    """Approximate Jaccard similarity using MinHash."""
    def get_minhashes(sequence, k, n):
        kmers = [sequence[i:i+k] for i in range(len(sequence)-k+1)]
        hashes = [int(hashlib.md5(km.encode()).hexdigest(), 16) for km in kmers]
        return [min(h % (2**32) for h in hashes)]  # simplified
    
    mh1 = set(get_minhashes(seq1, k, n_hashes))
    mh2 = set(get_minhashes(seq2, k, n_hashes))
    return len(mh1 & mh2) / len(mh1 | mh2)
```

---

## De Bruijn Graphs (Genome Assembly)

Modern genome assembly uses de Bruijn graphs:

```python
from collections import defaultdict

def build_debruijn_graph(reads: list, k: int = 31) -> dict:
    """Build de Bruijn graph from reads."""
    graph = defaultdict(list)
    
    for read in reads:
        for i in range(len(read) - k + 1):
            kmer = read[i:i+k]
            prefix = kmer[:-1]  # (k-1)-mer
            suffix = kmer[1:]   # (k-1)-mer
            graph[prefix].append(suffix)
    
    return dict(graph)

def find_eulerian_path(graph: dict) -> list:
    """Find Eulerian path through de Bruijn graph (genome sequence)."""
    # Check for Eulerian path conditions
    in_degree = defaultdict(int)
    out_degree = defaultdict(int)
    
    for node, neighbors in graph.items():
        out_degree[node] += len(neighbors)
        for neighbor in neighbors:
            in_degree[neighbor] += 1
    
    # Find start node
    start = None
    for node in set(list(graph.keys()) + [n for neighbors in graph.values() for n in neighbors]):
        if out_degree[node] - in_degree[node] == 1:
            start = node
            break
    
    if start is None:
        start = next(iter(graph))
    
    # Hierholzer's algorithm for Eulerian path
    path = []
    stack = [start]
    local_graph = {k: list(v) for k, v in graph.items()}
    
    while stack:
        v = stack[-1]
        if local_graph.get(v):
            u = local_graph[v].pop()
            stack.append(u)
        else:
            path.append(stack.pop())
    
    return path[::-1]

def reconstruct_sequence(path: list) -> str:
    """Reconstruct sequence from Eulerian path."""
    return path[0] + ''.join(node[-1] for node in path[1:])
```

---

## Motif Finding

Finding regulatory sequence motifs (like transcription factor binding sites):

```python
import numpy as np
from scipy.stats import entropy

def build_pwm(sequences: list) -> np.ndarray:
    """Build Position Weight Matrix from aligned sequences."""
    n = len(sequences[0])
    bases = 'ACGT'
    
    # Count matrix
    counts = np.zeros((4, n))
    for seq in sequences:
        for i, base in enumerate(seq):
            if base in bases:
                counts[bases.index(base), i] += 1
    
    # Add pseudocounts and normalize
    pwm = (counts + 0.25) / (len(sequences) + 1)  # Laplace smoothing
    return pwm

def score_sequence(pwm: np.ndarray, sequence: str) -> float:
    """Score how well a sequence matches a PWM motif."""
    bases = 'ACGT'
    score = 0
    for i, base in enumerate(sequence):
        if base in bases and i < pwm.shape[1]:
            score += np.log2(pwm[bases.index(base), i] / 0.25)  # log-odds
    return score

def scan_for_motif(genome: str, pwm: np.ndarray,
                   threshold: float = 8.0) -> list:
    """Scan genome for motif occurrences."""
    motif_len = pwm.shape[1]
    hits = []
    for i in range(len(genome) - motif_len + 1):
        window = genome[i:i+motif_len]
        score = score_sequence(pwm, window)
        if score >= threshold:
            hits.append({'position': i, 'sequence': window, 'score': score})
    return hits

def motif_information_content(pwm: np.ndarray) -> np.ndarray:
    """Calculate information content per position."""
    background = np.array([0.25, 0.25, 0.25, 0.25])
    ic = np.log2(4) - np.array([entropy(col, background, base=2) for col in pwm.T])
    return ic  # bits per position
```

---

## Problem-Solving: Sequence Analysis

**Common Problem Patterns:**

```python
# Pattern 1: "Find all occurrences of pattern P in sequence S"
def find_all_occurrences(sequence: str, pattern: str) -> list:
    return [i for i in range(len(sequence) - len(pattern) + 1)
            if sequence[i:i+len(pattern)] == pattern]
# → Use KMP or Boyer-Moore for efficiency

# Pattern 2: "What is the evolutionary distance between two sequences?"
def evolutionary_distance(seq1: str, seq2: str) -> float:
    a1, a2, _ = needleman_wunsch(seq1, seq2)
    matches = sum(1 for a, b in zip(a1, a2) if a == b and a != '-')
    return 1 - matches / len(a1)

# Pattern 3: "Find all open reading frames (ORFs)"
def find_orfs(sequence: str, min_length: int = 100) -> list:
    orfs = []
    for frame in range(3):
        for strand, seq in [(1, sequence), (-1, str(Seq(sequence).reverse_complement()))]:
            for i in range(frame, len(seq) - 2, 3):
                codon = seq[i:i+3]
                if codon == 'ATG':  # start codon
                    for j in range(i+3, len(seq) - 2, 3):
                        stop = seq[j:j+3]
                        if stop in ['TAA', 'TAG', 'TGA']:
                            if j - i >= min_length:
                                orfs.append({'start': i, 'end': j+3, 'strand': strand, 'length': j-i})
                            break
    return orfs
```
