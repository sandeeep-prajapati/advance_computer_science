# Chapter 7: Synthetic Biology

## Engineering Biology Like Software

Synthetic biology treats biological systems as programmable:

```
Software Engineering  →  Synthetic Biology
──────────────────────────────────────────
Function/Class        →  Gene / Protein
API / Interface       →  Promoter / RBS
Library               →  BioBrick Part
Package               →  Genetic Circuit
Compiler              →  Codon Optimizer
Debugger              →  Flow Cytometer
Version Control       →  Registry of Parts
Unit Test             →  Cell-Free Assay
Deployment            →  Cell Transformation
Production Server     →  Fermenter/Bioreactor
```

---

## Genetic Parts: The Building Blocks

```python
from dataclasses import dataclass
from typing import Optional

@dataclass
class GeneticPart:
    """A reusable biological component (BioBrick standard)."""
    part_id: str          # e.g., "BBa_B0034"
    part_type: str        # promoter, RBS, CDS, terminator
    sequence: str         # DNA sequence
    strength: float       # relative to standard reference
    organism: str         # E. coli, S. cerevisiae, etc.
    description: str

# Standard BioBrick prefix/suffix for assembly
BIOBRICK_PREFIX = "GAATTCGCGGCCGCTTCTAGAG"  # EcoRI + NotI + XbaI
BIOBRICK_SUFFIX = "CTGCAGCGGCCGCTACTAGTAGC"  # PstI + NotI + SpeI

def assemble_biobricks(parts: list) -> str:
    """Simulate BioBrick assembly (3A ligation)."""
    return BIOBRICK_PREFIX + "".join(p.sequence for p in parts) + BIOBRICK_SUFFIX

# Common parts database
STANDARD_PARTS = {
    'BBa_J23100': GeneticPart('BBa_J23100', 'promoter', 'TTGACAGCTAGCTCAGTCCTAGGTATAATGCTAGC', 1.0, 'E. coli', 'Standard constitutive promoter (strongest)'),
    'BBa_J23106': GeneticPart('BBa_J23106', 'promoter', 'TTGACAGCTAGCTCAGTCCTAGGTATTATGCTAGC', 0.47, 'E. coli', 'Constitutive promoter (medium)'),
    'BBa_B0034':  GeneticPart('BBa_B0034',  'RBS',      'AAAGAGGAGAAA', 1.0, 'E. coli', 'Strong RBS'),
    'BBa_B0015':  GeneticPart('BBa_B0015',  'terminator','CCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATA', 1.0, 'E. coli', 'Double terminator'),
}
```

---

## Genetic Circuits: Boolean Logic in Cells

```python
# Genetic circuits implement computational logic
# Gate types implemented in vivo:

class GeneticGate:
    """Simulate genetic logic gate behavior."""
    
    @staticmethod
    def NOT(input_signal: float, max_output: float = 1.0,
            hill_n: float = 2.0, k: float = 0.5) -> float:
        """NOT gate: repressor controlled by input.
        
        High input → low output (repression)
        Uses Hill function for transcription kinetics.
        """
        return max_output / (1 + (input_signal / k) ** hill_n)
    
    @staticmethod
    def AND(input1: float, input2: float) -> float:
        """AND gate: both TFs needed for transcription."""
        return input1 * input2  # simplified; real gate uses cooperative binding
    
    @staticmethod
    def OR(input1: float, input2: float) -> float:
        """OR gate: either TF drives transcription."""
        return 1 - (1 - input1) * (1 - input2)
    
    @staticmethod
    def NAND(input1: float, input2: float) -> float:
        return GeneticGate.NOT(GeneticGate.AND(input1, input2))

# Toggle Switch (Gardner et al. 2000 — Nature)
def toggle_switch_ode(state: list, t: float,
                       a1: float = 156.25, a2: float = 15.6,
                       beta: float = 2.5, gamma: float = 1.0) -> list:
    """
    Two-gene bistable switch ODE model.
    state = [u, v] = [protein1_level, protein2_level]
    
    du/dt = a1 / (1 + v^beta) - u
    dv/dt = a2 / (1 + u^gamma) - v
    """
    u, v = state
    dudt = a1 / (1 + v**beta) - u
    dvdt = a2 / (1 + u**gamma) - v
    return [dudt, dvdt]

def simulate_toggle_switch(initial_state: list = [0.1, 10.0],
                            t_max: float = 50.0) -> tuple:
    """Simulate bistable toggle switch dynamics."""
    from scipy.integrate import odeint
    import numpy as np
    
    t = np.linspace(0, t_max, 1000)
    trajectory = odeint(toggle_switch_ode, initial_state, t)
    return t, trajectory

# Repressilator (Elowitz & Leibler 2000 — oscillator)
def repressilator_ode(state: list, t: float,
                       alpha: float = 216, beta: float = 5.0,
                       n: float = 2.0) -> list:
    """
    Three-gene oscillator: lacI → tetR → cI → lacI (cycle).
    state = [m1, m2, m3, p1, p2, p3]  (mRNA + protein for each gene)
    """
    m1, m2, m3, p1, p2, p3 = state
    dm1 = -m1 + alpha / (1 + p3**n)
    dm2 = -m2 + alpha / (1 + p1**n)
    dm3 = -m3 + alpha / (1 + p2**n)
    dp1 = -beta * (p1 - m1)
    dp2 = -beta * (p2 - m2)
    dp3 = -beta * (p3 - m3)
    return [dm1, dm2, dm3, dp1, dp2, dp3]
```

---

## Codon Optimization

Optimizing DNA sequences for expression in target organism:

```python
from collections import defaultdict
import numpy as np

# Codon usage table for E. coli (frequency per thousand)
ECOLI_CODON_USAGE = {
    'TTT': 0.58, 'TTC': 0.42, 'TTA': 0.14, 'TTG': 0.13,
    'CTT': 0.12, 'CTC': 0.10, 'CTA': 0.04, 'CTG': 0.47,
    'ATT': 0.49, 'ATC': 0.39, 'ATA': 0.11, 'ATG': 1.00,
    'GTT': 0.28, 'GTC': 0.20, 'GTA': 0.17, 'GTG': 0.35,
    'TCT': 0.17, 'TCC': 0.15, 'TCA': 0.14, 'TCG': 0.14,
    'CCT': 0.18, 'CCC': 0.13, 'CCA': 0.20, 'CCG': 0.49,
    'ACT': 0.19, 'ACC': 0.40, 'ACA': 0.17, 'ACG': 0.25,
    'GCT': 0.18, 'GCC': 0.26, 'GCA': 0.23, 'GCG': 0.33,
    # ... (full table has all 64 codons)
}

# Reverse: amino acid → codon frequency table
AA_TO_CODONS = defaultdict(list)
from ..genetic_engineering.genetics_utils import CODON_TABLE  # from ch01
# Build reverse mapping
_AA_CODON_MAP = defaultdict(list)
for codon, aa in {
    'TTT': 'F', 'TTC': 'F', 'TTA': 'L', 'TTG': 'L',
    'CTT': 'L', 'CTC': 'L', 'CTA': 'L', 'CTG': 'L',
    'ATT': 'I', 'ATC': 'I', 'ATA': 'I', 'ATG': 'M',
    'GTT': 'V', 'GTC': 'V', 'GTA': 'V', 'GTG': 'V',
    'TAA': '*', 'TAG': '*', 'TGA': '*',
}.items():
    _AA_CODON_MAP[aa].append(codon)

def codon_optimize(protein_sequence: str,
                   codon_usage: dict = None,
                   strategy: str = 'most_frequent') -> str:
    """
    Optimize codons for expression in target organism.
    
    strategy:
      'most_frequent' — always use most common codon (may cause RNA structure issues)
      'harmonized'    — match source organism's codon usage pattern
      'random_weighted' — sample proportionally to usage frequency
    """
    if codon_usage is None:
        codon_usage = ECOLI_CODON_USAGE
    
    # Group codons by amino acid
    aa_to_codons = defaultdict(list)
    for codon, freq in codon_usage.items():
        if len(codon) == 3 and codon not in ('TAA', 'TAG', 'TGA'):
            # Map codon → amino acid using genetic code
            aa = _get_aa(codon)
            if aa and aa != '*':
                aa_to_codons[aa].append((codon, freq))
    
    optimized = []
    for aa in protein_sequence:
        if aa == '*':
            optimized.append('TAA')  # preferred stop codon
            break
        
        candidates = aa_to_codons.get(aa, [])
        if not candidates:
            continue
        
        if strategy == 'most_frequent':
            best_codon = max(candidates, key=lambda x: x[1])[0]
        elif strategy == 'random_weighted':
            codons, freqs = zip(*candidates)
            freqs = np.array(freqs)
            freqs /= freqs.sum()
            best_codon = np.random.choice(codons, p=freqs)
        else:
            best_codon = max(candidates, key=lambda x: x[1])[0]
        
        optimized.append(best_codon)
    
    return ''.join(optimized)

def _get_aa(codon: str) -> str:
    """Return amino acid for codon."""
    table = {
        'TTT': 'F', 'TTC': 'F', 'TTA': 'L', 'TTG': 'L',
        'CTT': 'L', 'CTC': 'L', 'CTA': 'L', 'CTG': 'L',
        'ATT': 'I', 'ATC': 'I', 'ATA': 'I', 'ATG': 'M',
        'GTT': 'V', 'GTC': 'V', 'GTA': 'V', 'GTG': 'V',
        'TCT': 'S', 'TCC': 'S', 'TCA': 'S', 'TCG': 'S',
        'CCT': 'P', 'CCC': 'P', 'CCA': 'P', 'CCG': 'P',
        'ACT': 'T', 'ACC': 'T', 'ACA': 'T', 'ACG': 'T',
        'GCT': 'A', 'GCC': 'A', 'GCA': 'A', 'GCG': 'A',
        'TAT': 'Y', 'TAC': 'Y', 'TAA': '*', 'TAG': '*',
        'CAT': 'H', 'CAC': 'H', 'CAA': 'Q', 'CAG': 'Q',
        'AAT': 'N', 'AAC': 'N', 'AAA': 'K', 'AAG': 'K',
        'GAT': 'D', 'GAC': 'D', 'GAA': 'E', 'GAG': 'E',
        'TGT': 'C', 'TGC': 'C', 'TGA': '*', 'TGG': 'W',
        'CGT': 'R', 'CGC': 'R', 'CGA': 'R', 'CGG': 'R',
        'AGT': 'S', 'AGC': 'S', 'AGA': 'R', 'AGG': 'R',
        'GGT': 'G', 'GGC': 'G', 'GGA': 'G', 'GGG': 'G',
    }
    return table.get(codon, '?')

def calculate_cai(dna_sequence: str, codon_usage: dict) -> float:
    """
    Codon Adaptation Index — measure of codon optimization.
    Range: 0 (bad) to 1.0 (perfectly optimized for organism).
    """
    # Find max frequency codon per amino acid
    aa_max_freq = defaultdict(float)
    for codon, freq in codon_usage.items():
        aa = _get_aa(codon)
        if aa and aa != '*':
            aa_max_freq[aa] = max(aa_max_freq[aa], freq)
    
    log_cai = 0
    n = 0
    for i in range(0, len(dna_sequence) - 2, 3):
        codon = dna_sequence[i:i+3]
        aa = _get_aa(codon)
        if aa and aa != '*' and aa_max_freq[aa] > 0:
            freq = codon_usage.get(codon, 0.01)
            rscu = freq / aa_max_freq[aa]  # Relative Synonymous Codon Usage
            log_cai += np.log(rscu + 1e-10)
            n += 1
    
    return np.exp(log_cai / n) if n > 0 else 0
```

---

## Metabolic Engineering

Engineering organisms to produce valuable compounds:

```python
# Metabolic Flux Analysis — how much product can we make?

def parse_stoichiometry(reactions: list) -> tuple:
    """
    Parse metabolic network reactions into stoichiometry matrix S.
    
    reactions: list of strings like "A + 2B -> C + D"
    Returns: (S_matrix, metabolites, reaction_names)
    """
    import re
    from collections import defaultdict
    
    metabolites = set()
    stoich = []
    
    for rxn in reactions:
        left, right = rxn.split('->')
        coeffs = defaultdict(float)
        
        for side, sign in [(left, -1), (right, 1)]:
            for term in side.strip().split('+'):
                term = term.strip()
                match = re.match(r'(\d+\.?\d*)\s*(\w+)', term)
                if match:
                    coeff, met = float(match.group(1)), match.group(2)
                else:
                    coeff, met = 1.0, term
                coeffs[met] += sign * coeff
                metabolites.add(met)
        stoich.append(coeffs)
    
    met_list = sorted(metabolites)
    S = np.zeros((len(met_list), len(reactions)))
    for j, rxn_coeffs in enumerate(stoich):
        for met, coeff in rxn_coeffs.items():
            S[met_list.index(met), j] = coeff
    
    return S, met_list

def flux_balance_analysis(S: np.ndarray, objective_idx: int,
                           lb: np.ndarray, ub: np.ndarray) -> dict:
    """
    Flux Balance Analysis (FBA): maximize product flux.
    Solves: max c^T v  s.t. Sv=0, lb≤v≤ub
    """
    from scipy.optimize import linprog
    
    n_rxns = S.shape[1]
    c = np.zeros(n_rxns)
    c[objective_idx] = -1  # maximize (linprog minimizes)
    
    result = linprog(c, A_eq=S, b_eq=np.zeros(S.shape[0]),
                     bounds=list(zip(lb, ub)), method='highs')
    
    if result.success:
        return {'fluxes': result.x, 'max_flux': -result.fun, 'status': 'optimal'}
    return {'status': result.message}
```

---

## Cell-Free Systems and PURE System

```python
# Cell-free protein synthesis: transcription + translation in a test tube
# No living cells needed — ideal for prototyping genetic circuits

def design_cell_free_circuit(circuit: dict) -> str:
    """
    Design a genetic circuit for cell-free expression.
    
    circuit: {
        'promoter': 'T7',     # T7 RNA polymerase promoter for in vitro
        'gene': 'GFP',
        'reporter': 'fluorescence',
        'input': 'IPTG',      # inducer molecule
    }
    """
    sequence_parts = []
    
    if circuit['promoter'] == 'T7':
        # T7 promoter: strong, works in cell-free with T7 RNAP
        sequence_parts.append('TAATACGACTCACTATA')  # T7 promoter sequence
    
    sequence_parts.append('AAGGAGATATACAT')  # Shine-Dalgarno (RBS)
    sequence_parts.append('ATG')              # Start codon
    sequence_parts.append(f"[{circuit['gene']}_CDS_sequence]")
    sequence_parts.append('CTGCAGCGGCCGC')   # terminator
    
    return ''.join(sequence_parts)
```

---

## iGEM and Standardized Parts

```python
# iGEM (International Genetically Engineered Machine)
# Registry: parts.igem.org — open source biology

# Key standardized promoter family: Anderson Promoters
ANDERSON_PROMOTERS = {
    'BBa_J23100': 1.000,   # strongest
    'BBa_J23101': 0.700,
    'BBa_J23102': 0.860,
    'BBa_J23103': 0.010,   # very weak
    'BBa_J23104': 0.720,
    'BBa_J23105': 0.240,
    'BBa_J23106': 0.470,
    'BBa_J23107': 0.360,
    'BBa_J23108': 0.510,
    'BBa_J23109': 0.040,
    'BBa_J23110': 0.330,
    'BBa_J23111': 0.580,
    'BBa_J23112': 0.010,
    'BBa_J23113': 0.010,
    'BBa_J23114': 0.100,
    'BBa_J23115': 0.150,
    'BBa_J23116': 0.160,
    'BBa_J23117': 0.060,
    'BBa_J23118': 0.560,
    'BBa_J23119': 0.900,
}

def design_expression_system(target_expression_level: float,
                              protein_sequence: str,
                              organism: str = 'E. coli') -> dict:
    """Design a complete expression cassette for target protein."""
    # Select promoter closest to desired level
    best_promoter = min(ANDERSON_PROMOTERS.items(),
                        key=lambda x: abs(x[1] - target_expression_level))
    
    # Optimize codons for organism
    optimized_cds = codon_optimize(protein_sequence)
    
    return {
        'promoter': best_promoter[0],
        'promoter_strength': best_promoter[1],
        'rbs': 'BBa_B0034',           # strong RBS
        'cds': optimized_cds,
        'terminator': 'BBa_B0015',
        'vector': 'pUC19' if organism == 'E. coli' else 'pRS416',
        'antibiotic_resistance': 'AmpR',
    }
```

---

## Biosafety and Containment

```python
# Kill switches: containment of GMOs
KILL_SWITCH_DESIGNS = {
    'essential_nutrient_dependency': {
        'description': 'Delete gene for essential amino acid, supplement in lab only',
        'example': 'Delete pyrF (uracil synthesis); grow only with uracil supplementation',
        'escape_frequency': 1e-6,
    },
    'toxin_antitoxin': {
        'description': 'Express antitoxin continuously; stop = cell death',
        'example': 'MazE antitoxin (need IPTG), MazF toxin (mRNA interferase)',
        'escape_frequency': 1e-7,
    },
    'orthogonal_amino_acid_dependency': {
        'description': 'Use synthetic amino acid (not in nature) for essential protein',
        'example': 'Unnatural amino acid at amber stop codon in essential gene',
        'escape_frequency': 1e-10,
    },
    'crispr_based': {
        'description': 'CRISPR continuously targets chromosome if trigger absent',
        'example': 'dCas9-repressor on essential gene requires specific inducer to release',
        'escape_frequency': 1e-8,
    },
}
```
