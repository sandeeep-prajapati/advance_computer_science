# Chapter 9: Grover's Search Algorithm

## Problem Statement

```
PROBLEM: Unstructured Search
  Input:  N items in a list, one (or k) are "marked"
  Goal:   Find a marked item
  Oracle: O(x) = 1 if x is marked, 0 otherwise

CLASSICAL:
  Expected: N/2 queries (check randomly)
  Worst:    N queries
  Lower bound: Ω(N) — must check all items

QUANTUM (Grover, 1996):
  O(√N) queries — provably optimal quantum algorithm
  Quadratic speedup

PRACTICAL SIZES:
  N = 1,000,000 → Classical: 500,000 queries, Grover: ~785 queries
  N = 10^18 (RSA keyspace) → Classical: impossible, Grover: 10^9 queries
```

---

## The Algorithm

### Setup:
```
n qubits to encode N = 2ⁿ items
Oracle Oφ: phase oracle that marks the target(s)

Initial state: |s⟩ = H⊗n|0⟩ⁿ = (1/√N) Σₓ |x⟩
```

### Grover Operator G = DOs:
```
1. Apply phase oracle Oφ: |x⟩ → (-1)^f(x)|x⟩
   (flips sign of marked items' amplitudes)

2. Apply Grover diffusion D = 2|s⟩⟨s| - I
   (inversion about the mean amplitude)
   D|ψ⟩: each amplitude aₓ → 2⟨a⟩ - aₓ

Combined effect: REFLECTS in marked-state subspace then in initial-state subspace.
Net rotation by 2θ in 2D space.
```

### Geometric interpretation:
```
View state as 2D vector in {|good⟩, |bad⟩} space:

Initial: sin(θ)|good⟩ + cos(θ)|bad⟩    where sin(θ) = √(k/N)
After G: sin(3θ)|good⟩ + cos(3θ)|bad⟩
After G²: sin(5θ)|good⟩ + cos(5θ)|bad⟩
After m steps: sin((2m+1)θ)|good⟩ + ...

Maximum when (2m+1)θ = π/2  →  m = π/(4θ) - 1/2 ≈ (π/4)√(N/k)

For k=1: m ≈ (π/4)√N iterations
```

### Circuit:
```
|0⟩^⊗n ─H⊗n─[─Oφ─DiffuserD─]^m─ Measure
              └─── m rounds ─┘
```

---

## Step-by-Step Example (N=4, n=2, target |11⟩)

```
N=4, target x* = 11, initial θ = π/6 (sin(π/6) = 1/2 = √(1/4))
Optimal iterations: m = 1

Initial: |s⟩ = (|00⟩+|01⟩+|10⟩+|11⟩)/2 = [1/2, 1/2, 1/2, 1/2]

After Oracle (flip |11⟩):
  [1/2, 1/2, 1/2, -1/2]

Mean amplitude = (1/2 + 1/2 + 1/2 - 1/2)/4 = 1/4

After Diffuser (invert about mean 1/4):
  each amp aₓ → 2(1/4) - aₓ = 1/2 - aₓ

  1/2 → 0
  1/2 → 0
  1/2 → 0
  -1/2 → 1

State = |11⟩ (amplitude 1!) → measure → get 11 with certainty.
```

---

## Implementation in Qiskit

```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator
import numpy as np

def grover_oracle(n: int, target: str):
    """Phase oracle: flips phase of |target⟩"""
    qc = QuantumCircuit(n)
    # Flip qubits where target has 0 (to use multi-controlled X on all-1s)
    for i, bit in enumerate(reversed(target)):
        if bit == '0':
            qc.x(i)
    # Multi-controlled Z (phase flip for |1...1⟩)
    qc.h(n - 1)
    qc.mcx(list(range(n - 1)), n - 1)  # multi-controlled X
    qc.h(n - 1)
    # Flip back
    for i, bit in enumerate(reversed(target)):
        if bit == '0':
            qc.x(i)
    return qc

def grover_diffuser(n: int):
    """Grover diffusion operator: 2|s><s| - I"""
    qc = QuantumCircuit(n)
    qc.h(range(n))       # H⊗n
    qc.x(range(n))       # X⊗n
    qc.h(n - 1)
    qc.mcx(list(range(n - 1)), n - 1)  # multi-controlled X
    qc.h(n - 1)
    qc.x(range(n))       # X⊗n
    qc.h(range(n))       # H⊗n
    return qc

def grover_circuit(n: int, target: str, num_iterations: int = None):
    """Complete Grover's search circuit."""
    if num_iterations is None:
        num_iterations = int(np.pi / 4 * np.sqrt(2**n))
    
    qc = QuantumCircuit(n, n)
    
    # Initialize superposition
    qc.h(range(n))
    
    oracle = grover_oracle(n, target)
    diffuser = grover_diffuser(n)
    
    # Grover iterations
    for _ in range(num_iterations):
        qc.compose(oracle, inplace=True)
        qc.compose(diffuser, inplace=True)
    
    # Measure
    qc.measure(range(n), range(n))
    
    return qc

# Search for '101' in 8 items (n=3)
target = '101'
n = 3
qc = grover_circuit(n, target)
print(qc.draw())

sim = AerSimulator()
counts = sim.run(qc, shots=1000).result().get_counts()
print(f"Target: {target}")
print(f"Results: {sorted(counts.items(), key=lambda x: -x[1])[:3]}")
# '101' should have highest count
```

---

## Multiple Solutions (k marked items)

```
For k marked items out of N:
  θ = arcsin(√(k/N))
  Optimal iterations: m ≈ (π/4)√(N/k)
  Success probability: ≈ 1

Unknown k → use quantum counting first, or use BBHT algorithm
(Boyer-Brassard-Høyer-Tapp: finds marked item with O(√(N/k)) queries even unknown k)
```

---

## Quantum Counting

How many marked items are there?

```
Uses Quantum Phase Estimation to estimate k:
  Eigenvalue of Grover operator G: e^(±2iθ) where sin(θ) = √(k/N)
  QPE estimates θ → compute k = N·sin²(θ)

Circuit: Phase estimation with G as the unitary
Precision: uses t ancilla qubits → estimates k to t bits
Queries: O(√N)
```

---

## Grover's Algorithm Complexity Analysis

```
Query complexity:     O(√N)   — provably optimal
Quantum time:         O(n·√N)  — n-bit operations per query
Classical lower bound: Ω(N)

Improvement over classical: √N factor

Note: For k marked items → O(√(N/k))
  If k = N/4 (quarter marked): O(1) iterations!
  If k = N (all marked): 0 iterations needed

Limitations:
  - Only quadratic (not exponential) speedup
  - Need quantum oracle — hard to implement for many real problems
  - Practical speedup limited by decoherence and gate fidelity
```

---

## Real-World Applications

```
Cryptography:
  Break symmetric encryption by searching keyspace
  128-bit AES: 2^128 keys → √(2^128) = 2^64 Grover queries
  Mitigation: Double key lengths (256-bit AES)

Optimization:
  Combinatorial optimization: find minimum in unstructured space
  Speedup: O(√N) → helps but not a silver bullet

SAT solving:
  Apply Grover to find satisfying assignment
  Speedup over brute force, but NP-complete remains hard

Database search:
  Only useful without classical structure
  Indexed databases are already O(log N) classically — Grover doesn't help
```

---

## Problem-Solving with Grover

**Template:**
```python
def grover_solve(n, is_solution):
    """
    n: number of qubits
    is_solution: function that checks if a bitstring is a solution
    """
    target_count = sum(1 for x in range(2**n) if is_solution(format(x, f'0{n}b')))
    iterations = int(np.pi / 4 * np.sqrt(2**n / max(target_count, 1)))
    
    # Build phase oracle from is_solution (classically)
    # For real hardware: implement is_solution as quantum circuit
    
    # ... (build and run Grover's circuit)
    pass

# Example: 3-SAT oracle
# Find assignment satisfying: (x0 OR x1) AND (NOT x1 OR x2) AND (x0 OR NOT x2)
def satisfies_3sat(assignment: str) -> bool:
    x0, x1, x2 = [b == '1' for b in assignment]
    return (x0 or x1) and (not x1 or x2) and (x0 or not x2)

# Count solutions
n = 3
solutions = [x for x in range(8) if satisfies_3sat(format(x, '03b'))]
print(f"Solutions: {[format(x,'03b') for x in solutions]}")
# Grover would find one in O(√8) ≈ 2.8 iterations
```

---

## Grover's Optimality

```
THEOREM (Bennett et al., 1997):
  Any quantum algorithm needs Ω(√N) queries to find a marked item
  in an unstructured database of N items.

PROOF IDEA (hybrid argument):
  Consider a quantum algorithm with T < cπ√N queries.
  By hybrid argument, it cannot distinguish a database with one
  marked item from one with no marked item with probability > 1/2.
  Therefore, T must be Ω(√N).

This means Grover is OPTIMAL — no quantum algorithm can do better for unstructured search.
```
