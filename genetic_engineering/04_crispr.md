# Chapter 4: CRISPR-Cas9 Gene Editing

## What is CRISPR?

CRISPR-Cas9 is a molecular "find-and-replace" tool for DNA:

```
Find:    Guide RNA (gRNA) directs Cas9 to a specific DNA sequence
Cut:     Cas9 protein creates a double-strand break (DSB)
Replace: Cell's repair machinery either:
           - NHEJ: introduces insertions/deletions (disrupts gene)
           - HDR: inserts provided DNA template (precise edit)

Your CS analogy:
  gRNA = search query (matches ~20 bp in genome)
  Cas9 = restriction enzyme (cuts at match)
  NHEJ = lossy compression (garbles the sequence → gene knockout)
  HDR  = find-replace with template (precise correction)
```

---

## The Mechanism

```
1. gRNA design:
   Target sequence:     5'- ATGCTAGCTAGCTAGCTAGC -3'   (20 nt protospacer)
   PAM sequence (NGG):  5'-                   NGG -3'   (Cas9 requirement)

2. Cas9-gRNA searches genome for: [20-nt match] + NGG
   When found: Cas9 unwinds DNA, gRNA hybridizes, Cas9 cuts

3. After cut: DNA repair
   NHEJ (error-prone): 
     |---20bp---|--NGG--|
          ↓
     |---20bp--XXXX--|  (insertion/deletion → frameshift → gene knockout)
   
   HDR (template repair):
     Provide repair template → insert exact sequence
     More precise, but less efficient (~1-10% efficiency)
```

---

## Guide RNA Design

The most critical computational component:

```python
import re
import numpy as np

def find_crispr_sites(sequence: str, pam: str = 'NGG') -> list:
    """Find all possible CRISPR target sites in a sequence."""
    sites = []
    # Convert PAM pattern to regex
    pam_regex = pam.replace('N', '[ATGC]')
    
    # Search both strands
    for strand, seq in [('+', sequence), ('-', reverse_complement(sequence))]:
        for match in re.finditer(f'(?=(.{{20}}{pam_regex}))', seq):
            start = match.start()
            protospacer = match.group(1)[:20]
            pam_seq = match.group(1)[20:]
            
            # Skip sequences with poly-T (affects gRNA expression)
            if 'TTTT' in protospacer:
                continue
            
            sites.append({
                'strand': strand,
                'position': start,
                'protospacer': protospacer,
                'pam': pam_seq,
                'gc_content': gc_content(protospacer),
                'on_target_score': estimate_on_target(protospacer),
            })
    
    return sites

def gc_content(seq: str) -> float:
    return (seq.count('G') + seq.count('C')) / len(seq)

def estimate_on_target_doench(protospacer: str) -> float:
    """
    Simplified Doench 2016 on-target efficiency score.
    Full model uses Random Forest trained on screen data.
    """
    score = 0.597636
    
    # Nucleotide preferences at each position (simplified)
    # Real model: 30 features including dinucleotides
    gc = gc_content(protospacer)
    
    # GC content penalty (too low or too high is bad)
    if gc < 0.25 or gc > 0.75:
        score -= 0.2
    
    # Position-specific weights (very simplified)
    # Position 20 (cleavage site adjacent): G preferred
    if protospacer[-1] == 'G':
        score += 0.05
    
    # No poly-G
    if 'GGGG' in protospacer:
        score -= 0.3
    
    return max(0, min(1, score))

def reverse_complement(seq: str) -> str:
    complement = {'A': 'T', 'T': 'A', 'G': 'C', 'C': 'G', 'N': 'N'}
    return ''.join(complement[b] for b in reversed(seq))
```

---

## Off-Target Prediction

The key safety challenge in CRISPR:

```python
def count_mismatches(seq1: str, seq2: str) -> int:
    """Count mismatches between two equal-length sequences."""
    return sum(a != b for a, b in zip(seq1, seq2))

def find_off_targets(guide: str, genome: str, max_mismatches: int = 3) -> list:
    """
    Find potential off-target sites (simplified).
    Real tools: Cas-OFFinder, CRISPOR, CHOPCHOP
    Note: this naive implementation is O(genome_size) — real tools use indexing.
    """
    off_targets = []
    n = len(guide)
    
    for strand, seq in [('+', genome), ('-', reverse_complement(genome))]:
        for i in range(len(seq) - n - 2):
            window = seq[i:i+n]
            pam = seq[i+n:i+n+3]
            
            if not re.match(r'[ATGC]GG', pam):
                continue
            
            mismatches = count_mismatches(guide, window)
            if mismatches <= max_mismatches:
                # Weight mismatches by position (seed region is more important)
                # Seed region: last 12 bp adjacent to PAM
                seed_mismatches = count_mismatches(guide[-12:], window[-12:])
                off_targets.append({
                    'sequence': window,
                    'position': i,
                    'strand': strand,
                    'total_mismatches': mismatches,
                    'seed_mismatches': seed_mismatches,
                    'pam': pam,
                })
    
    # Sort by fewest mismatches (highest risk first)
    return sorted(off_targets, key=lambda x: (x['seed_mismatches'], x['total_mismatches']))

# ML-based off-target prediction (using trained models)
def predict_off_target_activity(guide: str, target: str, model=None) -> float:
    """
    Predict off-target cleavage activity using trained neural network.
    Real model: CRISPR-ML, DeepCRISPR, Azimuth
    """
    mismatches = count_mismatches(guide, target)
    mismatch_positions = [i for i, (a, b) in enumerate(zip(guide, target)) if a != b]
    
    # Simple heuristic (use trained DL model in practice)
    if mismatches == 0:
        return 1.0
    
    # Penalize mismatches in seed region (positions 8-20 from PAM)
    seed_mismatches = sum(1 for pos in mismatch_positions if pos >= 8)
    non_seed_mismatches = mismatches - seed_mismatches
    
    activity = np.exp(-1.5 * seed_mismatches - 0.5 * non_seed_mismatches)
    return min(1.0, max(0.0, activity))
```

---

## CRISPR Variants

Beyond Cas9:

```
Cas9 (SpCas9):     Standard, large (4.2 kb), requires NGG PAM
SaCas9:            Smaller (3.3 kb), NNGRRT PAM, better for AAV delivery
Cas12a (Cpf1):     Staggered cut, T-rich PAM (TTTN), better specificity
Cas13:             Targets RNA (not DNA), useful for RNA knockdown
CasRx:             Compact Cas13, RNA knockdown
BE3/ABE:           Base editors (C→T or A→G, no double-strand break)
PE (Prime Editor):  "Search-and-replace" without double-strand break

Newer: CRISPRi, CRISPRa (no cut, just regulate expression):
  dCas9-KRAB: transcriptional repressor
  dCas9-VP64: transcriptional activator
```

---

## Base Editing

More precise than standard CRISPR — no double-strand break:

```python
# Cytosine Base Editor (CBE): C→T (or G→A on complementary strand)
# Adenine Base Editor (ABE): A→G (or T→C)
# No double-strand break → much safer, fewer indels

def design_base_editor(mutation: str) -> dict:
    """
    Given a target SNP correction, design base editor strategy.
    mutation: e.g., 'C>T' (change C to T)
    """
    strategies = {
        'C>T': {'editor': 'CBE (BE4max)', 'window': '4-8 from PAM', 'notes': 'Edit window 4-8 nt from PAM'},
        'G>A': {'editor': 'CBE (on minus strand)', 'window': '4-8', 'notes': 'Target complementary strand'},
        'A>G': {'editor': 'ABE (ABE8e)', 'window': '4-8', 'notes': 'A→G in edit window'},
        'T>C': {'editor': 'ABE (on minus strand)', 'window': '4-8', 'notes': 'Target complementary strand'},
    }
    return strategies.get(mutation, {'editor': 'Prime Editor', 'window': 'N/A', 'notes': 'Use PE for complex edits'})
```

---

## gRNA Selection with ML

```python
# Training a guide RNA efficiency predictor
import torch
import torch.nn as nn

class GuideRNAPredictor(nn.Module):
    """Predict CRISPR guide RNA on-target efficiency."""
    
    def __init__(self, sequence_length: int = 30):
        super().__init__()
        # One-hot encode sequence: 4 bases × length
        self.conv = nn.Sequential(
            nn.Conv1d(4, 32, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.Conv1d(32, 64, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.AdaptiveAvgPool1d(1),
        )
        self.fc = nn.Sequential(
            nn.Linear(64, 32),
            nn.ReLU(),
            nn.Dropout(0.3),
            nn.Linear(32, 1),
            nn.Sigmoid()
        )
    
    def encode_sequence(self, sequence: str) -> torch.Tensor:
        """One-hot encode DNA sequence."""
        base_map = {'A': 0, 'T': 1, 'G': 2, 'C': 3}
        encoding = torch.zeros(4, len(sequence))
        for i, base in enumerate(sequence):
            if base in base_map:
                encoding[base_map[base], i] = 1
        return encoding
    
    def forward(self, x: torch.Tensor) -> torch.Tensor:
        features = self.conv(x).squeeze(-1)
        return self.fc(features).squeeze(-1)

def train_guide_predictor(sequences: list, efficiencies: list,
                          epochs: int = 50) -> GuideRNAPredictor:
    """Train guide RNA efficiency predictor."""
    model = GuideRNAPredictor()
    optimizer = torch.optim.Adam(model.parameters(), lr=1e-3)
    criterion = nn.MSELoss()
    
    # Encode sequences
    X = torch.stack([model.encode_sequence(s[:30]) for s in sequences])
    y = torch.tensor(efficiencies, dtype=torch.float32)
    
    for epoch in range(epochs):
        optimizer.zero_grad()
        predictions = model(X)
        loss = criterion(predictions, y)
        loss.backward()
        optimizer.step()
        if epoch % 10 == 0:
            print(f"Epoch {epoch}: Loss = {loss.item():.4f}")
    
    return model
```

---

## CRISPR Delivery Methods

Getting CRISPR into cells is the major engineering challenge:

```
Method          Cargo         Cell types      Efficiency  Safety
──────────────────────────────────────────────────────────────────
Plasmid DNA     Cas9+gRNA     Dividing cells  Medium      Low (integrates?)
Electroporation RNP or RNA    Most cells      High        Better (transient)
Viral (AAV)     Small Cas9+g  In vivo         High        Good, limited size
Viral (LV)      Cas9+gRNA     Dividing cells  High        Integrates (risk)
Lipid NPs (LNP) RNA           Liver, lung     High in vivo  Good
Nanoparticles   Various       Various         Developing  Promising

In vivo delivery choices:
  Liver:  LNP (similar to mRNA vaccines)
  Eye:    Subretinal injection + AAV
  Muscle: AAV intramuscular injection
  Brain:  AAV stereotaxic injection (invasive)
  Blood:  Ex vivo: edit cells, reinfuse
```

---

## Real Applications

```python
# Case study: Sickle Cell Disease
# Disease: mutation in HBB gene: GAG→GTG (E6V) → abnormal hemoglobin
# CRISPR approach: reactivate fetal hemoglobin (HbF)
# Target: BCL11A enhancer (suppressor of HbF expression)
# Method: Delete 58-bp segment in BCL11A enhancer → de-repress HbF
# Result: HbF levels rise → functional hemoglobin replaces sickle HbF

# FDA approved: Casgevy (2023) — first approved CRISPR therapy

# Algorithm for sickle cell:
def design_sickle_cell_therapy():
    return {
        'target_gene': 'BCL11A',
        'edit_type': 'NHEJ deletion in enhancer',
        'guides': ['GTAACGGCAGACTTGACCTGG', 'AAGCCTGGCTTTAAAGCTCAC'],
        'rationale': 'BCL11A represses HbF; deletion reactivates fetal hemoglobin',
        'delivery': 'Ex vivo HSC editing + autologous transplant',
        'clinical_status': 'FDA approved (Casgevy, 2023)',
    }
```

---

## CRISPR Tools and APIs

```python
# CRISPOR: comprehensive guide design tool
# Usage: crispor.tefor.net or local install

# Python API for CHOPCHOP
import requests

def chopchop_design(target_sequence: str, organism: str = 'human') -> list:
    """Use CHOPCHOP web API for guide design."""
    # (Demo — actual API call would vary)
    url = "http://chopchop.cbu.uib.no/cgi-bin/chopchop.cgi"
    # ... parse response
    return []

# Azimuth: Doench 2016 on-target scoring
# pip install azimuth
try:
    import azimuth.model_comparison as az
    def doench_score(sequences_30mer: list) -> np.ndarray:
        """Score 30-mer sequences (4+20+3+3 = 30)."""
        return az.predict(np.array(sequences_30mer))
except ImportError:
    pass
```
