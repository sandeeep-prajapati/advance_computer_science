# Chapter 15: Quantum Complexity Theory

## Classical Complexity Review

```
P:    Problems solvable in polynomial time (deterministic)
NP:   Problems verifiable in polynomial time
BPP:  Problems solvable in polynomial time with bounded error (randomized)
PSPACE: Problems solvable in polynomial space

Known: P ⊆ NP ⊆ PSPACE,  P ⊆ BPP
Open: P = NP? BPP = P?
```

---

## Quantum Complexity Classes

### BQP (Bounded-error Quantum Polynomial time):
```
Decision problems solvable by quantum computer in polynomial time
with error probability ≤ 1/3.

Formally: L ∈ BQP if there's a uniform family of polynomial-size
quantum circuits {Cₙ} such that for x of length n:
  x ∈ L → P[Cₙ(x) = 1] ≥ 2/3
  x ∉ L → P[Cₙ(x) = 1] ≤ 1/3

BQP is the quantum analog of BPP.
```

### Known inclusions:
```
P ⊆ BPP ⊆ BQP ⊆ PSPACE

P ⊆ BQP: Classical computation is special case
BQP ⊆ PSPACE: Quantum computation can be simulated in polynomial space

OPEN: BPP = BQP? (Is quantum strictly more powerful than classical randomized?)
      NP ⊆ BQP? (Can quantum solve NP-hard problems efficiently? Probably not.)
```

### QMA (Quantum Merlin-Arthur):
```
Quantum analog of NP.
L ∈ QMA: there exists a quantum proof |ψ⟩ that a verifier circuit accepts with ≥2/3 probability.

QMA-complete problems:
  k-local Hamiltonian: Is ground energy of H below threshold?
    (Quantum analogy of SAT)
  Consistency of density matrices
  Quantum circuit satisfiability

Known: NP ⊆ QMA (classical NP proofs are quantum proofs too)
       BQP ⊆ QMA (any BQP language has trivial quantum proof)
Open:  QMA = NP?  QMA = PSPACE?
```

---

## The Quantum Complexity Landscape

```
        PSPACE
           │
          QMA
         / │ \
       NP  │  BQP
       │   │   │
       │  BPP  │
        \ / \ /
          P

Arrows = containment (A → B means A ⊆ B)
Dashed lines = unknown relationship

Key question for quantum advantage:
  BQP vs NP:
    Factoring ∈ BQP but factoring ∈ NP ∩ coNP (not known NP-complete)
    BQP likely doesn't contain NP-complete problems
    (Quantum cannot solve Traveling Salesman, SAT, etc. exponentially faster — probably)
```

---

## Query Complexity

Studying problems via the number of oracle queries:

```
Classical randomized: R(f) queries
Quantum: Q(f) queries

Known separations:
  Deutsch-Jozsa:  R(f) = Ω(N)  vs  Q(f) = 1   (exponential)
  Simon's:        R(f) = Ω(2^(n/2)) vs Q(f) = O(n)  (exponential)
  Grover's:       R(f) = Ω(N)  vs  Q(f) = Θ(√N)  (quadratic)

Known limits:
  Polynomial method: Q(f) ≥ Ω(√R(f))  (quantum gives at most polynomial speedup in query model)
  This is why we can't get exponential speedup for NP-complete search
```

### Polynomial method (key technique):
```
For any quantum algorithm making T queries:
  The acceptance probability is a degree-2T polynomial in the input bits.
  
  Lower bound on queries = deg(f)/2
  
  Example: OR function requires Ω(√N) quantum queries (Grover is optimal).
  
  PARITY: deg(PARITY) = n, so Q(PARITY) ≥ n/2 (linear — no quantum speedup)
```

---

## Quantum Supremacy

Demonstrating that a quantum computer performs a task no classical computer can feasibly do:

### Boson Sampling:
```
Photonic quantum computer with n photons, m modes:
  Sample from distribution of boson interference patterns
  
  Computing this classically: #P-hard (as hard as counting problems)
  Quantum: natural by design
  
  First demonstrations: ~50 photons in 2020 (Jiuzhang, USTC)
```

### Random Circuit Sampling (Google 2019):
```
Sycamore chip: 53 qubits, depth-20 random circuit
  Quantum: ~200 seconds
  Classical estimate: 10,000 years
  (Later improved to ~300 seconds classically — supremacy claims disputed)
```

### Limitations of supremacy claims:
```
1. Classical simulation algorithms are rapidly improving
2. Supremacy tasks are not useful computationally
3. Error mitigation challenges make results hard to verify
4. "Supremacy" ≠ practical advantage

True criterion: Useful quantum advantage on industry-relevant problem.
```

---

## Computational Hardness of Classical Simulation

```
Simulating n-qubit quantum computer:
  Naive: O(2ⁿ) memory for state vector
  
  Clifford circuits: efficient (Gottesman-Knill)
  Low-entanglement circuits: efficient (tensor networks, MPS)
  
  Hard to simulate:
    High-entanglement circuits beyond ~50 qubits
    Random circuits
    Quantum error correction circuits at large scale

Why BQP ⊆ PSPACE:
  Can simulate by computing each amplitude separately
  Each amplitude: exponential time but polynomial space
  
Why PSPACE (probably) ⊄ BQP:
  Quantum can't solve PSPACE-complete problems efficiently (believed)
  IP = PSPACE but no evidence quantum collapses this
```

---

## Oracle Separations

Using query complexity to prove theoretical separations:

```
Bernstein-Vazirani (1993):
  Relative to a random oracle: BQP ⊄ BPP
  (Quantum is more powerful than classical — relative to oracle)

Simon (1994):
  There exists oracle A: BQP^A ⊄ BPP^A
  (Exponential separation with oracle)

But: Oracles don't resolve real complexity questions!
  Example: P^Oracle = NP^Oracle doesn't mean P = NP.
```

---

## Quantum PCP and the NLTS Theorem

Recent breakthrough (2022):

```
Classical PCP Theorem: NP-hard problems have short proofs with local checks.

Quantum PCP Conjecture: QMA-hard problems have quantum proofs checkable locally.
  Status: OPEN (one of most important open problems in quantum complexity)

NLTS Theorem (Anshu, Breuckmann, Nirkhe 2022):
  There exist local Hamiltonians whose ground states require linear circuit depth.
  
  Implication: Some quantum states cannot be efficiently prepared.
  Progress toward proving quantum PCP.
```

---

## Practical Complexity Insights

```
When designing quantum algorithms, these complexity facts matter:

1. Don't expect exponential speedup for NP-complete problems.
   (TSP, SAT, coloring — quantum probably doesn't help exponentially)

2. Quadratic speedup (Grover) is provably optimal for unstructured search.
   Can't do better than O(√N) for unstructured oracles.

3. Exponential speedup is possible for:
   - Problems with hidden algebraic structure (Shor's, Simon's)
   - Quantum simulation (Feynman's original vision)
   - Phase estimation for eigenvalue problems

4. QMA-complete problems: Need quantum proof — classical verification insufficient.
   These may be HARDER than NP for quantum too!

5. NISQ advantages: unclear, ongoing research
   Many claimed quantum ML advantages don't hold up under scrutiny.
```

---

## Complexity Table for Key Problems

| Problem | Classical | Quantum | Speedup |
|---|---|---|---|
| Unstructured search | O(N) | O(√N) | Quadratic |
| Integer factoring | O(exp(n^(1/3))) | O(n³) | Exponential |
| Discrete logarithm | O(exp(n^(1/2))) | O(n³) | Exponential |
| Parity | O(N) queries | O(N) queries | None |
| Graph isomorphism | exp? | BQP? | Unknown |
| NP-complete (SAT) | O(2^n) | O(2^(n/2)) max | Quadratic (at best) |
| Quantum simulation | O(exp(n)) | O(poly(n)) | Exponential |
| Linear systems (HHL) | O(N) | O(polylog N) | Exponential* |

*HHL has caveats: exponential speedup only for well-conditioned systems, quantum output

---

## Problem-Solving with Complexity Theory

```
Before designing a quantum algorithm, ask:

1. Is the problem in P classically? 
   → Quantum won't give useful speedup (exponential speedup would collapse P≠BQP)

2. Is the problem NP-hard?
   → Quantum probably gives at most quadratic speedup (via Grover)
   → Unless there's hidden algebraic structure

3. Is the problem in QMA?
   → May need quantum proof — potentially no classical verification

4. Is there periodic/algebraic structure?
   → QFT-based algorithms may give exponential speedup

5. Is it a search problem?
   → Grover gives √N speedup (sometimes worth it, sometimes not vs classical heuristics)
```
