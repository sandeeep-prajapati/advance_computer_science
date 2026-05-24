# Chapter 7: Quantum Algorithm Design

## The Quantum Algorithm Design Framework

Designing quantum algorithms requires a completely different mindset from classical algorithms. The goal is not to "run things in parallel" — it's to engineer interference so the right answer becomes highly probable.

```
Classical Algorithm Design:
  1. Break problem into subproblems
  2. Solve each subproblem
  3. Combine solutions

Quantum Algorithm Design:
  1. Find mathematical structure in the problem
  2. Encode inputs as quantum superposition
  3. Design oracle/unitary that marks correct answers
  4. Use interference to amplify correct answer's amplitude
  5. Measure to extract classical solution
```

---

## Query Complexity Model

Most quantum algorithms are analyzed in the **oracle model**: the algorithm queries a black-box function f(x) and we count the number of queries.

```
Classical algorithm: Ω(N) queries to find marked item in list of N
Grover's algorithm: O(√N) queries  → quadratic speedup

Classical algorithm: Ω(N) queries to determine if f is constant or balanced
Deutsch-Jozsa: O(1) queries → exponential speedup

Why oracle model?
  - Separates quantum from classical advantage cleanly
  - Lower bounds in oracle model apply to real algorithms
  - Makes algorithms hardware-independent
```

---

## Quantum Parallelism and its Limits

```
Prepare superposition: H⊗n|0...0⟩ = (1/√N) Σₓ |x⟩

Apply oracle: Uf: (1/√N) Σₓ |x⟩|0⟩ → (1/√N) Σₓ |x⟩|f(x)⟩

This evaluates f(x) for all x simultaneously!

BUT: Measuring now gives ONE random x — useless.

The trick: use interference BEFORE measuring.
Design unitary that transforms the state so the desired answer has amplitude ≈ 1.
```

---

## The Oracle (Black Box)

An oracle is a quantum gate that encodes a classical function f: {0,1}ⁿ → {0,1}:

### Phase oracle (most useful for algorithms):
```
Uφ|x⟩ = (-1)^f(x)|x⟩

If f(x) = 0: |x⟩ → |x⟩       (no change)
If f(x) = 1: |x⟩ → -|x⟩      (phase flip)

Construction from standard oracle:
  Start: |x⟩|−⟩ = |x⟩(|0⟩-|1⟩)/√2
  Apply Uf: |x⟩|f(x)⊕0⟩ - |x⟩|f(x)⊕1⟩)/√2
           = (-1)^f(x)|x⟩|−⟩

The |−⟩ ancilla picks up the phase — this is "phase kickback".
```

### Standard oracle:
```
Uf|x⟩|b⟩ = |x⟩|b ⊕ f(x)⟩

|x⟩|0⟩ → |x⟩|f(x)⟩    (compute f into ancilla)
|x⟩|1⟩ → |x⟩|1⊕f(x)⟩  (flip ancilla if f(x)=1)
```

```python
from qiskit import QuantumCircuit

def phase_oracle_example(n):
    """Oracle that marks the all-ones string 111...1"""
    qc = QuantumCircuit(n)
    qc.h(n-1)  # convert ancilla to |-⟩ (or use phase oracle directly)
    qc.mcx(list(range(n-1)), n-1)  # multi-controlled X
    qc.h(n-1)  # convert back
    return qc
```

---

## Amplitude Amplification

The core mechanism behind Grover's and many other algorithms:

```
Given state: |ψ⟩ = α|good⟩ + β|bad⟩

Grover operator G = -HWH · Oracle
  Oracle: flip phase of |good⟩
  -HWH:   reflect about average amplitude

Each G application rotates state in 2D subspace {|good⟩, |bad⟩}
Rotation angle: 2θ where sin(θ) = α = √(k/N)  (k = number of marked states)

After O(√(N/k)) applications:
  |ψ⟩ ≈ |good⟩   (amplitude of good state ≈ 1)
```

Generalization: Given any quantum algorithm A that prepares state with success probability p:
```
Amplitude Amplification gives success probability ≈ 1 with O(1/√p) applications of A.
```

---

## Quantum Fourier Transform (Preview)

The QFT is to quantum algorithms what FFT is to classical signal processing:

```
Classical DFT: f̂(k) = Σⱼ f(j) e^(2πijk/N)    — O(N log N)

QFT: |j⟩ → (1/√N) Σₖ e^(2πijk/N)|k⟩            — O(log²N)

Exponential speedup!

Used in: Shor's (period finding), Phase Estimation, quantum simulation
```

---

## Algorithm Paradigms

| Paradigm | Key Technique | Example Algorithms |
|---|---|---|
| **Amplitude Amplification** | Grover's operator | Grover's search, BBHT |
| **Phase Estimation** | QFT + controlled-U | Shor's, QPE, HHL |
| **Variable Time Amplitude Amp.** | Stopping condition | VTAA algorithms |
| **Quantum Walks** | Walk on graph | Element distinctness, MNRS |
| **Hybrid VQA** | Parameterized circuit + optimizer | VQE, QAOA |
| **Linear Algebra** | HHL algorithm | Solving Ax=b |
| **Adiabatic** | Slow Hamiltonian evolution | Adiabatic optimization |

---

## Speedup Types

| Speedup | Description | Example |
|---|---|---|
| **Exponential** | exp(n) → poly(n) | Shor's: factoring O(exp(n)) → O(n³) |
| **Polynomial** | n^a → n^b (b < a) | HHL: O(N) → O(polylog N) (under conditions) |
| **Quadratic** | N → √N | Grover's: O(N) → O(√N) |
| **Constant factor** | k·N → N | Some arithmetic speedups |
| **No speedup** | Same complexity | Most classical algorithms |

**Important:** Most problems do NOT have quantum speedup.
Quantum advantage requires specific mathematical structure.

---

## Quantum Complexity Classes

```
P   ⊆ BPP  ⊆ BQP  ⊆ PSPACE

P:     Classical polynomial time (deterministic)
BPP:   Classical probabilistic polynomial time
BQP:   Quantum polynomial time (with bounded error)
PSPACE: Classical polynomial space
NP:    Nondeterministic polynomial (classical)

BQP vs NP:
  - BQP ⊄ NP (likely) — quantum doesn't solve NP-hard problems generically
  - NP ⊄ BQP (likely) — quantum doesn't solve all of NP
  - Quantum gives superpolynomial speedup for specific problems

Key point: Shor's factors in BQP but factoring is not known to be in P.
Grover's searches in BQP but unstructured search lower bound is Ω(√N) for quantum too.
```

---

## Problem-Solving Heuristics

When you see a new problem, ask:

```
1. Is there periodicity?
   → Quantum Fourier Transform → period finding → Shor-like algorithms

2. Is there a search with a clear oracle?
   → Grover's / Amplitude Amplification → quadratic speedup

3. Is it an optimization problem with constraints?
   → QAOA / VQE → variational approach

4. Is it a graph problem?
   → Quantum walks → potentially polynomial speedup

5. Does it involve eigenvalues/eigenvectors?
   → Quantum Phase Estimation → exponential speedup if unitary can be simulated

6. Is it a linear system Ax = b?
   → HHL algorithm → exponential speedup (but output is quantum state)

7. Is it a machine learning kernel method?
   → Quantum kernel methods → potentially exponential speedup

8. None of the above?
   → Quantum probably doesn't help. Use classical algorithm.
```

---

## Designing an Oracle

The hardest part of quantum algorithm design is often the oracle:

```python
from qiskit import QuantumCircuit

# Example: Oracle for "is x a multiple of 3?" (x is 3-bit number)
# Multiples of 3 in {0,...,7}: 0, 3, 6

def multiple_of_3_oracle():
    """Phase oracle: flips phase of |x⟩ if x mod 3 == 0"""
    qc = QuantumCircuit(3)
    # Mark 000 (0), 011 (3), 110 (6)
    
    # Mark |000⟩
    qc.x([0, 1, 2])  # flip all
    qc.ccx(0, 1, 2)  # ... this approach needs careful construction
    qc.x([0, 1, 2])  # flip back
    
    # In practice: use reversible arithmetic
    # Compute x mod 3 into ancilla, phase flip if 0, uncompute
    return qc
```

**Oracle construction principle:**
1. Compute f(x) reversibly into an ancilla
2. Apply phase flip conditioned on correct answer
3. Uncompute (run step 1 in reverse) to restore ancilla

Uncomputation is essential to prevent ancilla qubits from entangling with data.

---

## Classical vs Quantum Algorithm Comparison

```
Problem: Find a marked item in a list of N items (one marked)

Classical:
  - Query items one by one: O(N) expected
  - No structure to exploit

Quantum (Grover):
  - O(√N) queries — provably optimal for quantum
  - Speedup: √N (quadratic)

Problem: Factor N-bit integer n

Classical:
  - Best: Number Field Sieve: O(exp(n^(1/3))) — superpolynomial

Quantum (Shor):
  - O(n³) — polynomial
  - Speedup: exponential

Why different speedups?
  Factoring has algebraic period structure → QFT exploits it.
  Unstructured search has no structure → only quadratic speedup possible.
```

---

## Algorithm Design Practice Problems

```
Problem 1: You have N items, k of which are "good." 
How many Grover queries do you need to find a good item?
Answer: O(√(N/k))

Problem 2: You can query f:{0,1}ⁿ→{0,1}. f is either:
  - Constant: f(x) = 0 for all x, or f(x) = 1 for all x
  - Balanced: f(x) = 0 for exactly N/2 inputs
Determine which case with minimum queries.
Classical: Ω(N/2+1). Quantum: 1 query (Deutsch-Jozsa).

Problem 3: f:{0,1}ⁿ→{0,1}ⁿ has a period s: f(x) = f(x⊕s) for all x.
Find s.
Classical: O(2^(n/2)) queries (birthday paradox). Quantum: O(n) queries (Simon's).

These are in the next chapters!
```
