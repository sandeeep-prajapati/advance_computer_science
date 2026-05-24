# Mastering Quantum Computing
## Complete Lecture Notes — Problem-Solving Approach

---

## Table of Contents

| Chapter | Title | Key Concepts |
|---|---|---|
| 01 | Introduction to Quantum Computing | Qubits, superposition, entanglement, interference, use cases |
| 02 | Linear Algebra for Quantum Computing | Vectors, matrices, inner products, tensor products, Dirac notation |
| 03 | Qubits and Quantum States | State vectors, Bloch sphere, multi-qubit systems, density matrices |
| 04 | Quantum Gates and Circuits | Single-qubit gates, multi-qubit gates, universal gate sets, circuit diagrams |
| 05 | Quantum Entanglement | Bell states, EPR paradox, quantum teleportation, dense coding |
| 06 | Quantum Measurement | Projective measurement, Born rule, partial measurement, POVM |
| 07 | Quantum Algorithm Design | Algorithm paradigms, oracle model, amplitude amplification, speedup types |
| 08 | Deutsch-Jozsa and Simon's Algorithm | Problem setup, circuit, complexity advantage, Simon's algorithm |
| 09 | Grover's Search Algorithm | Unstructured search, oracle, diffusion operator, quadratic speedup |
| 10 | Quantum Fourier Transform | DFT vs QFT, circuit, period finding |
| 11 | Shor's Factoring Algorithm | RSA threat, classical reduction, quantum period finding, circuit |
| 12 | Quantum Phase Estimation | Eigenvalue problem, QPE circuit, precision, applications |
| 13 | Variational Quantum Algorithms | VQE, QAOA, parameterized circuits, hybrid classical-quantum |
| 14 | Quantum Error Correction | Decoherence, bit flip code, Shor code, stabilizer formalism, surface codes |
| 15 | Quantum Complexity Theory | BQP, QMA, quantum supremacy, complexity landscape |
| 16 | Quantum Programming with Qiskit | IBM Q, circuits, simulation, running on real hardware |
| 17 | Quantum Programming with Cirq and PennyLane | Google's Cirq, PennyLane for ML, noise models |
| 18 | Quantum Machine Learning | QML algorithms, quantum kernels, quantum neural networks, QSVM |
| 19 | Quantum Hardware | Superconducting, trapped ion, photonic, topological, fidelity metrics |
| 20 | Learning Roadmap | 6-week plan, problem-solving patterns, competitive programming approach |

---

## Architecture Overview

```
Quantum Computing Stack
═══════════════════════════════════════════════════════════════

Application Layer
┌─────────────────────────────────────────────────────────────┐
│  Algorithms: Grover  Shor  VQE  QAOA  QML  QPE  Simon's    │
└─────────────────────────────────────────────────────────────┘
                           │
Algorithm Primitives
┌─────────────────────────────────────────────────────────────┐
│  QFT  Phase Estimation  Amplitude Amplification  Oracle     │
└─────────────────────────────────────────────────────────────┘
                           │
Circuit Layer
┌─────────────────────────────────────────────────────────────┐
│  Single-qubit gates      Multi-qubit gates    Measurements  │
│  H  X  Y  Z  S  T  Rx   CNOT  CZ  Toffoli   M_basis       │
└─────────────────────────────────────────────────────────────┘
                           │
Physics Layer
┌─────────────────────────────────────────────────────────────┐
│  Superposition   Entanglement   Interference   Decoherence  │
│  |ψ⟩ = α|0⟩ + β|1⟩    Bell states    Phase kickback        │
└─────────────────────────────────────────────────────────────┘
                           │
Hardware Layer
┌─────────────────────────────────────────────────────────────┐
│  Superconducting    Trapped Ion    Photonic    Topological  │
│  (IBM, Google)      (IonQ, Quantinuum)        (Microsoft)  │
└─────────────────────────────────────────────────────────────┘


Problem-Solving Paradigm
═══════════════════════════════════════════════════════════════

Classical Problem
      │
      ▼
Is there quantum advantage?
  │             │
 YES            NO → Use classical algorithm
  │
  ▼
Identify Structure:
 Periodicity? → Shor / QFT
 Search?      → Grover
 Optimization?→ QAOA / VQE
 Sampling?    → Quantum Supremacy circuits
  │
  ▼
Design Oracle / Ansatz
  │
  ▼
Build Quantum Circuit
  │
  ▼
Classical Post-processing
  │
  ▼
Solution
```

---

## Key Mathematical Notation

| Symbol | Meaning |
|---|---|
| `\|0⟩, \|1⟩` | Computational basis states |
| `\|ψ⟩ = α\|0⟩ + β\|1⟩` | General qubit state (superposition) |
| `α, β ∈ ℂ, \|α\|² + \|β\|² = 1` | Normalization constraint |
| `⟨ψ\|` | Bra (dual/conjugate of ket) |
| `⟨φ\|ψ⟩` | Inner product |
| `\|φ⟩⊗\|ψ⟩ = \|φψ⟩` | Tensor product (combining systems) |
| `H` | Hadamard gate (superposition) |
| `CNOT` | Controlled-NOT gate (entanglement) |
| `U†` | Conjugate transpose (inverse of unitary U) |
| `tr(ρ)` | Trace of density matrix |

---

## Prerequisites

- **Mathematics:** Linear algebra (vectors, matrices, eigenvalues), complex numbers, probability
- **Computer Science:** Classical algorithms, complexity theory (P vs NP), Boolean logic
- **Physics:** Helpful but not required — key physics explained from first principles

---

## Software Tools

| Tool | Language | Use Case |
|---|---|---|
| Qiskit | Python | IBM quantum computers, simulation |
| Cirq | Python | Google quantum computers |
| PennyLane | Python | Quantum machine learning |
| Braket SDK | Python | AWS quantum cloud |
| Q# | Q# | Microsoft Azure Quantum |
| Quirk | Browser | Visual circuit builder |
| IBM Quantum Lab | Browser | Free cloud quantum access |
# Chapter 1: Introduction to Quantum Computing

## What is Quantum Computing?

Quantum computing harnesses quantum mechanical phenomena — **superposition**, **entanglement**, and **interference** — to perform computations that are exponentially harder for classical computers.

```
Classical Bit:       Qubit:
    0 or 1           α|0⟩ + β|1⟩
                     (both simultaneously until measured)
```

A classical computer with n bits has exactly one of 2ⁿ states at any moment.
A quantum computer with n qubits can exist in a superposition of all 2ⁿ states simultaneously.

---

## The Three Quantum Phenomena

### 1. Superposition
A qubit can be in a combination of |0⟩ and |1⟩ at the same time:
```
|ψ⟩ = α|0⟩ + β|1⟩

where α, β are complex amplitudes
and |α|² + |β|² = 1 (probabilities sum to 1)

After measurement:
  P(outcome = 0) = |α|²
  P(outcome = 1) = |β|²
  And the state collapses to the measured value
```

### 2. Entanglement
Two or more qubits can be correlated in a way that measuring one instantly determines the state of the other, regardless of distance:
```
Bell state: |Φ+⟩ = (|00⟩ + |11⟩) / √2

Measuring qubit 1 as |0⟩ → qubit 2 instantly collapses to |0⟩
Measuring qubit 1 as |1⟩ → qubit 2 instantly collapses to |1⟩
```

### 3. Interference
Quantum amplitudes (unlike probabilities) can be negative or complex. They can interfere constructively (amplify correct answers) or destructively (cancel wrong answers):
```
Correct paths: interfere constructively → higher probability
Wrong paths:   interfere destructively  → lower probability

This is how quantum algorithms achieve speedup.
```

---

## Quantum vs Classical Computing

| Aspect | Classical | Quantum |
|---|---|---|
| Basic unit | Bit (0 or 1) | Qubit (superposition) |
| Parallelism | Execute one path | Explore all paths simultaneously |
| Randomness | Pseudorandom | Inherently probabilistic |
| Copying | Easy | Impossible (no-cloning theorem) |
| Error | Rare, correctable | Decoherence, needs QEC |
| Best at | Sequential logic, exact computation | Factoring, search, simulation, optimization |

---

## What Quantum Computers Are Good At

| Problem Type | Quantum Speedup | Example Algorithm |
|---|---|---|
| Integer factoring | Exponential | Shor's algorithm |
| Unstructured search | Quadratic | Grover's algorithm |
| Quantum simulation | Exponential | VQE (chemistry) |
| Linear systems | Polynomial | HHL algorithm |
| Machine learning | Potential | QSVM, QNN |
| Optimization | Heuristic | QAOA |
| Sampling | Exponential | Boson Sampling |

### What quantum computers are NOT good at:
- General-purpose computing (web servers, databases)
- Problems without mathematical structure
- Tasks where classical algorithms are already optimal
- Real-time interaction (high latency currently)

---

## Current State of Quantum Computing (2024–2025)

```
NISQ Era (Noisy Intermediate-Scale Quantum)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Qubits available:  50–1000+ physical qubits
Error rates:       ~0.1–1% per gate
Coherence time:    microseconds to milliseconds
Connectivity:      limited (not all-to-all)
Error correction:  not yet fault-tolerant

Key players:
  IBM:      Eagle (127Q), Osprey (433Q), Condor (1121Q)
  Google:   Sycamore (53Q), Willow (105Q)
  IonQ:     Trapped ion, ~30 algorithmic qubits
  Quantinuum: H-series trapped ion
  Microsoft: Topological qubits (in development)
  AWS:      Braket cloud access to multiple providers
```

---

## The Quantum Advantage Problem

**Quantum advantage** (or supremacy) means a quantum computer solves a problem faster than any classical computer could in practice.

```
Google's 2019 Sycamore claim:
  Task: Sample from random quantum circuit distribution
  Quantum: 200 seconds
  Classical (estimated): 10,000 years

But: The classical simulation was later improved to 300 seconds.
This shows how hard it is to demonstrate true quantum advantage.

For practical advantage (useful problems):
  Need: Fault-tolerant quantum computers
  Timeline: ~2030s (optimistic estimate)
```

---

## Real-World Applications

### Cryptography
- **Threat:** Shor's algorithm breaks RSA and elliptic curve cryptography
- **Defense:** Post-quantum cryptography (CRYSTALS-Kyber, CRYSTALS-Dilithium)
- **NIST PQC standards released in 2024**

### Drug Discovery / Chemistry
- Simulate molecular interactions at the quantum level
- Classical computers can't accurately simulate molecules with >50 electrons
- VQE algorithm running on quantum hardware can find ground state energies

### Optimization
- Logistics: route optimization, supply chain
- Finance: portfolio optimization, risk analysis
- QAOA algorithm provides heuristic solutions

### Machine Learning
- Quantum kernels for SVM classification
- Quantum neural networks (parametrized circuits)
- Still theoretical — unclear if practical advantage exists

---

## The Quantum Computing Mindset

```
Classical problem-solving:
  if/else → sequential decisions
  for loop → iterate through possibilities
  recursion → break into subproblems

Quantum problem-solving:
  Superposition  → prepare all inputs simultaneously
  Oracle         → mark correct answers
  Interference   → amplify correct, cancel wrong
  Measurement    → extract classical answer
```

**Key insight:** Quantum algorithms don't just "run in parallel" — they use *interference* to make correct answers more likely when measured. Simply preparing a superposition doesn't give you the answer — you need the algorithm to reshape the probability distribution.

---

## Quantum Computing Timeline

```
1981  Richard Feynman proposes quantum computers to simulate physics
1985  David Deutsch: first quantum Turing machine
1994  Shor's algorithm: exponential speedup for factoring
1996  Grover's algorithm: quadratic speedup for search
2001  IBM: 7-qubit NMR demonstration of Shor's algorithm
2009  Google: first superconducting qubit experiments
2016  IBM puts 5-qubit quantum computer on the cloud
2019  Google claims quantum supremacy (Sycamore)
2021  IBM: 127-qubit Eagle processor
2023  IBM: 1121-qubit Condor processor
2024  Google: Willow chip — exponential error reduction with scale
2030s Expected: fault-tolerant quantum computers
```

---

## Setting Up Your Learning Environment

```bash
# Install Qiskit (Python)
pip install qiskit qiskit-aer qiskit-ibm-runtime

# Install visualization
pip install matplotlib pylatexenc

# Install PennyLane for quantum ML
pip install pennylane

# Verify installation
python3 -c "import qiskit; print(qiskit.__version__)"
```

### Your first quantum circuit:
```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator

# Create a 2-qubit circuit
qc = QuantumCircuit(2, 2)

# Put qubit 0 in superposition
qc.h(0)

# Entangle qubit 0 and qubit 1
qc.cx(0, 1)

# Measure
qc.measure([0, 1], [0, 1])

# Simulate
simulator = AerSimulator()
job = simulator.run(qc, shots=1000)
counts = job.result().get_counts()
print(counts)
# Expected: {'00': ~500, '11': ~500} — Bell state!
```

---

## Problem-Solving Framework

Throughout this course, we use this framework for every quantum algorithm:

```
1. PROBLEM
   What is the classical problem? What's the input/output?

2. CLASSICAL COMPLEXITY
   Best known classical algorithm and its time complexity.

3. QUANTUM STRUCTURE
   What mathematical structure enables speedup?
   (Periodicity, symmetry, oracle access, amplitude amplification)

4. QUANTUM CIRCUIT
   How is the algorithm implemented as quantum gates?

5. COMPLEXITY ANALYSIS
   Quantum time/query/space complexity. What's the speedup?

6. LIMITATIONS
   When does this NOT help? What are the caveats?

7. CODE
   Qiskit implementation.
```

We'll apply this framework to every algorithm in this course.
# Chapter 2: Linear Algebra for Quantum Computing

## Why Linear Algebra?

Quantum mechanics IS linear algebra. Every quantum state is a vector, every operation is a matrix, and measurement is a projection. Master this and everything else follows.

---

## Vectors and Dirac Notation

### Ket (column vector):
```
|0⟩ = [1]    |1⟩ = [0]    |ψ⟩ = [α] = α|0⟩ + β|1⟩
      [0]          [1]           [β]
```

### Bra (row vector — conjugate transpose):
```
⟨0| = [1  0]    ⟨ψ| = [α*  β*]    (* denotes complex conjugate)
```

### Inner product (bra × ket → scalar):
```
⟨φ|ψ⟩ = φ†ψ    (dot product of complex vectors)

⟨0|0⟩ = 1    ⟨0|1⟩ = 0    ⟨1|1⟩ = 1    (orthonormal basis)

For |ψ⟩ = α|0⟩ + β|1⟩:
⟨ψ|ψ⟩ = |α|² + |β|² = 1    (normalization)
```

### Outer product (ket × bra → matrix):
```
|0⟩⟨0| = [1][1 0] = [1 0]    (projection onto |0⟩)
          [0]       [0 0]

|ψ⟩⟨ψ| = density matrix for pure state ψ
```

---

## Complex Numbers Review

Quantum amplitudes are complex numbers: `z = a + bi`

```
Conjugate:       z* = a - bi
Magnitude:       |z| = √(a² + b²)
Phase:           z = |z|·e^(iθ) = |z|(cos θ + i sin θ)    (Euler's formula)

Key: e^(iπ) = -1    (Euler's identity — phase flips sign)
```

In quantum mechanics, the *global phase* is unobservable:
```
|ψ⟩ and e^(iθ)|ψ⟩ represent the same physical state.
But *relative phase* between |0⟩ and |1⟩ components is crucial!

|+⟩ = (|0⟩ + |1⟩)/√2    ← relative phase = 0
|-⟩ = (|0⟩ - |1⟩)/√2    ← relative phase = π (very different state!)
```

---

## Matrices and Operations

### Matrix multiplication:
```
Gate U applied to state |ψ⟩:
U|ψ⟩ = (matrix) × (column vector)

Example — X gate (bit flip):
X = [0 1] applied to |0⟩ = [1]:
    [1 0]                  [0]

X|0⟩ = [0 1][1] = [0] = |1⟩  ✓
        [1 0][0]   [1]
```

### Conjugate transpose (Hermitian adjoint):
```
For matrix M, M† = (M*)ᵀ    (conjugate of each entry, then transpose)

X† = [0 1]† = [0 1] = X    (X is its own inverse — Hermitian and unitary)
     [1 0]    [1 0]
```

### Unitary matrices (all quantum gates are unitary):
```
U is unitary iff U†U = UU† = I    (preserves normalization)

This means:
  - All quantum gates are reversible
  - Inner products are preserved: ⟨Uφ|Uψ⟩ = ⟨φ|ψ⟩
  - Column (and row) vectors form orthonormal sets
```

### Hermitian matrices (observables):
```
H is Hermitian iff H† = H

Properties:
  - All eigenvalues are real (physical measurement values)
  - Eigenvectors form orthonormal basis
  
Example: Pauli Z matrix Z = [1  0 ] has eigenvalues +1 (for |0⟩) and -1 (for |1⟩)
                             [0 -1]
```

---

## Tensor Products (Multi-Qubit Systems)

The tensor product ⊗ combines systems:

```
|00⟩ = |0⟩⊗|0⟩ = [1] ⊗ [1] = [1·[1]] = [1]
                   [0]   [0]   [1·[0]]   [0]
                               [0·[1]]   [0]
                               [0·[0]]   [0]

So |00⟩ = [1,0,0,0]ᵀ, |01⟩ = [0,1,0,0]ᵀ, |10⟩ = [0,0,1,0]ᵀ, |11⟩ = [0,0,0,1]ᵀ
```

For gates on multi-qubit systems:
```
Apply H to qubit 0, I (identity) to qubit 1:
H⊗I = [1  1] ⊗ [1 0] = [H  0] = [1  0  1  0]
      [1 -1]   [0 1]   [0  H]   [0  1  0  1]
      ——/√2——                    [1  0 -1  0]
                                 [0  1  0 -1]  (all divided by √2)
```

### General rule:
```
n-qubit system has state space of dimension 2ⁿ
n-qubit gate is a 2ⁿ × 2ⁿ unitary matrix

This exponential growth is why quantum simulation is hard classically,
and why quantum computers can explore exponential state spaces.
```

---

## Eigenvalues and Eigenvectors

```
M|v⟩ = λ|v⟩    (|v⟩ is eigenvector, λ is eigenvalue)

For Pauli Z = [1  0]:
              [0 -1]

Z|0⟩ = +1·|0⟩    (|0⟩ is eigenstate with eigenvalue +1)
Z|1⟩ = -1·|1⟩    (|1⟩ is eigenstate with eigenvalue -1)

For Hadamard H = [1  1]/√2:
                 [1 -1]

H|+⟩ = +1·|+⟩    where |+⟩ = (|0⟩+|1⟩)/√2
H|-⟩ = -1·|-⟩    where |-⟩ = (|0⟩-|1⟩)/√2
```

---

## Key Bases

```
Computational (Z) basis:    |0⟩, |1⟩
Hadamard (X) basis:         |+⟩ = (|0⟩+|1⟩)/√2,  |-⟩ = (|0⟩-|1⟩)/√2
Y basis:                    |+i⟩= (|0⟩+i|1⟩)/√2, |-i⟩= (|0⟩-i|1⟩)/√2

Change of basis:
  H transforms between Z and X bases:
  H|0⟩ = |+⟩,    H|1⟩ = |-⟩
  H|+⟩ = |0⟩,    H|-⟩ = |1⟩
```

---

## The Pauli Matrices

The fundamental building blocks of quantum operations:

```
I = [1 0]    X = [0 1]    Y = [0 -i]    Z = [1  0]
    [0 1]        [1 0]        [i  0]        [0 -1]

Identity    Bit flip    Bit+Phase flip    Phase flip

Key properties:
  X² = Y² = Z² = I
  XY = iZ,  YZ = iX,  ZX = iY
  Pauli matrices are both unitary and Hermitian
  They form a basis for all 2×2 Hermitian matrices
```

---

## Important Linear Algebra Facts for QC

### Spectral decomposition:
```
Any Hermitian H = Σᵢ λᵢ|eᵢ⟩⟨eᵢ|    (sum over eigenvectors weighted by eigenvalues)

Measurement in a basis {|eᵢ⟩} gives outcome i with probability |⟨eᵢ|ψ⟩|²
```

### Matrix exponential (used for rotations):
```
e^(iθX/2) = cos(θ/2)I + i·sin(θ/2)X = Rx(θ) rotation gate

e^(iθZ/2) = [e^(iθ/2)    0      ] = Rz(θ) rotation gate
             [0        e^(-iθ/2)]
```

### Trace:
```
tr(M) = sum of diagonal elements = sum of eigenvalues

tr(|ψ⟩⟨ψ|) = 1    (pure state density matrix has trace 1)
```

---

## Problem-Solving with Linear Algebra

**Problem type 1: Is this state normalized?**
```python
import numpy as np

psi = np.array([1/np.sqrt(2), 1/np.sqrt(2)])  # |+⟩
print(np.dot(psi.conj(), psi))  # Should be 1.0
```

**Problem type 2: Apply a gate to a state**
```python
# Apply Hadamard to |0⟩
H = np.array([[1, 1], [1, -1]]) / np.sqrt(2)
zero = np.array([1, 0])  # |0⟩
result = H @ zero
print(result)  # [0.707, 0.707] = |+⟩
```

**Problem type 3: Compute measurement probabilities**
```python
# State: (|0⟩ + i|1⟩) / √2
psi = np.array([1/np.sqrt(2), 1j/np.sqrt(2)])
prob_0 = abs(psi[0])**2  # = 0.5
prob_1 = abs(psi[1])**2  # = 0.5
```

**Problem type 4: Is this gate unitary?**
```python
H = np.array([[1, 1], [1, -1]]) / np.sqrt(2)
print(np.allclose(H @ H.conj().T, np.eye(2)))  # True if unitary
```

**Problem type 5: Tensor product of two states**
```python
psi1 = np.array([1, 0])      # |0⟩
psi2 = np.array([1, 0])      # |0⟩
joint = np.kron(psi1, psi2)  # |00⟩ = [1,0,0,0]
```

---

## Quick Reference

| Operation | Math | Python (numpy) |
|---|---|---|
| Inner product | `⟨φ\|ψ⟩` | `phi.conj() @ psi` |
| Outer product | `\|ψ⟩⟨φ\|` | `np.outer(psi, phi.conj())` |
| Gate application | `U\|ψ⟩` | `U @ psi` |
| Tensor product | `A⊗B` | `np.kron(A, B)` |
| Conjugate transpose | `U†` | `U.conj().T` |
| Probability | `\|⟨i\|ψ⟩\|²` | `abs(psi[i])**2` |
| Norm | `\|\|ψ\|\|` | `np.linalg.norm(psi)` |
| Eigendecomposition | `Mv=λv` | `np.linalg.eigh(M)` |
# Chapter 3: Qubits and Quantum States

## The Qubit

A qubit is the fundamental unit of quantum information. Unlike a classical bit (exactly 0 or 1), a qubit exists in a superposition:

```
|ψ⟩ = α|0⟩ + β|1⟩

where:
  α, β ∈ ℂ    (complex amplitudes)
  |α|² + |β|² = 1    (normalization — probabilities sum to 1)

Measurement outcomes:
  P(0) = |α|²    (probability of measuring 0)
  P(1) = |β|²    (probability of measuring 1)
  After measurement → state collapses to |0⟩ or |1⟩
```

Physical implementations of qubits:
- Superconducting: two energy levels of a Josephson junction (IBM, Google)
- Trapped ion: two electronic states of an ion (IonQ, Quantinuum)
- Photon: polarization (horizontal/vertical)
- Nitrogen-vacancy center: electron spin in diamond
- Quantum dot: electron spin

---

## The Bloch Sphere

Any single-qubit state can be visualized as a point on the Bloch sphere:

```
          |0⟩ (North Pole)
           │
           │  ← |ψ⟩ = cos(θ/2)|0⟩ + e^(iφ)sin(θ/2)|1⟩
    |+i⟩ ──┼── |+⟩
           │
           │
          |1⟩ (South Pole)

θ ∈ [0, π]:    polar angle (from north pole)
φ ∈ [0, 2π):   azimuthal angle

Key points:
  North pole (θ=0):         |0⟩ = [1, 0]
  South pole (θ=π):         |1⟩ = [0, 1]
  +X axis (θ=π/2, φ=0):     |+⟩ = (|0⟩+|1⟩)/√2
  -X axis (θ=π/2, φ=π):     |-⟩ = (|0⟩-|1⟩)/√2
  +Y axis (θ=π/2, φ=π/2):   |+i⟩= (|0⟩+i|1⟩)/√2
  -Y axis (θ=π/2, φ=3π/2):  |-i⟩= (|0⟩-i|1⟩)/√2
```

Quantum gates = rotations on the Bloch sphere:
- X gate: 180° rotation around X axis (|0⟩ ↔ |1⟩)
- Z gate: 180° rotation around Z axis (phase flip)
- H gate: 90° rotation that maps poles to equator

---

## Multi-Qubit Systems

For n qubits, the state space is 2ⁿ dimensional:

```
2 qubits: {|00⟩, |01⟩, |10⟩, |11⟩}  — 4 basis states

|ψ⟩ = α₀₀|00⟩ + α₀₁|01⟩ + α₁₀|10⟩ + α₁₁|11⟩

Normalization: |α₀₀|² + |α₀₁|² + |α₁₀|² + |α₁₁|² = 1

Tensor product (independent qubits):
  |ψ₁⟩⊗|ψ₂⟩ = (α|0⟩+β|1⟩)⊗(γ|0⟩+δ|1⟩)
              = αγ|00⟩ + αδ|01⟩ + βγ|10⟩ + βδ|11⟩
```

### Entangled states (cannot be written as tensor product):
```
Bell state: |Φ+⟩ = (|00⟩ + |11⟩)/√2

Can we write this as |ψ₁⟩⊗|ψ₂⟩?
(α|0⟩+β|1⟩)⊗(γ|0⟩+δ|1⟩) = αγ|00⟩ + αδ|01⟩ + βγ|10⟩ + βδ|11⟩

For (|00⟩ + |11⟩)/√2:
  αγ = 1/√2  →  α ≠ 0 and γ ≠ 0
  βδ = 1/√2  →  β ≠ 0 and δ ≠ 0
  αδ = 0     →  α = 0 or δ = 0  CONTRADICTION!

Therefore |Φ+⟩ is entangled — no product state representation.
```

---

## Separable vs Entangled States

```
Separable (product state):      |ψ⟩ = |ψ₁⟩ ⊗ |ψ₂⟩
  Example: |+⟩⊗|0⟩ = (|00⟩+|10⟩)/√2

Entangled:                      |ψ⟩ ≠ |ψ₁⟩ ⊗ |ψ₂⟩
  Example: (|00⟩+|11⟩)/√2

How to check: Schmidt decomposition
  If only one Schmidt coefficient is non-zero → separable
  If more than one → entangled
```

---

## The Four Bell States

The maximally entangled 2-qubit states:

```
|Φ+⟩ = (|00⟩ + |11⟩)/√2    ← most common entangled state
|Φ-⟩ = (|00⟩ - |11⟩)/√2
|Ψ+⟩ = (|01⟩ + |10⟩)/√2
|Ψ-⟩ = (|01⟩ - |10⟩)/√2

Properties:
  - Orthonormal: ⟨Φ+|Φ-⟩ = 0, etc.
  - Maximally entangled: measuring one qubit gives random result
  - Used in: quantum teleportation, dense coding, QKD
```

---

## Mixed States and Density Matrices

When we don't know the exact quantum state (due to noise or entanglement with environment), we use a **density matrix** ρ:

```
Pure state |ψ⟩:   ρ = |ψ⟩⟨ψ|    (rank 1, tr(ρ²) = 1)

Mixed state:       ρ = Σᵢ pᵢ|ψᵢ⟩⟨ψᵢ|    (probabilistic mixture)
                   where pᵢ ≥ 0 and Σpᵢ = 1

Properties of density matrix:
  - Hermitian: ρ† = ρ
  - Positive semidefinite: ⟨v|ρ|v⟩ ≥ 0 for all |v⟩
  - Unit trace: tr(ρ) = 1
  - Pure state: tr(ρ²) = 1
  - Mixed state: tr(ρ²) < 1
```

### Example: maximally mixed state
```
If a qubit is in |0⟩ with probability 1/2 and |1⟩ with probability 1/2:

ρ = (1/2)|0⟩⟨0| + (1/2)|1⟩⟨1| = [1/2   0 ]
                                    [0    1/2]
= I/2

tr(ρ²) = tr(I/4) = 1/2 < 1  → mixed state

This is different from superposition! (|0⟩+|1⟩)/√2 has ρ=[1/2  1/2; 1/2  1/2]
```

---

## Partial Trace (Tracing Out a Subsystem)

For a 2-qubit system, the reduced density matrix of qubit 1:

```
ρ₁ = tr₂(ρ)    (trace out qubit 2)

Example: Bell state |Φ+⟩ = (|00⟩+|11⟩)/√2

ρ_AB = |Φ+⟩⟨Φ+| = (1/2)(|00⟩⟨00| + |00⟩⟨11| + |11⟩⟨00| + |11⟩⟨11|)

ρ_A = tr_B(ρ_AB) = (1/2)(|0⟩⟨0| + |1⟩⟨1|) = I/2

Qubit A is in a maximally mixed state!
This is the signature of maximal entanglement.
```

---

## Qubit State Examples

```python
import numpy as np
from qiskit.quantum_info import Statevector, DensityMatrix

# |0⟩ state
zero = Statevector([1, 0])

# |+⟩ state (superposition)
plus = Statevector([1/np.sqrt(2), 1/np.sqrt(2)])

# |i⟩ state (complex amplitude)
plus_i = Statevector([1/np.sqrt(2), 1j/np.sqrt(2)])

# Measure probabilities
print(plus.probabilities_dict())   # {'0': 0.5, '1': 0.5}
print(plus_i.probabilities_dict()) # {'0': 0.5, '1': 0.5} — same probs, different phase

# Check if pure
dm = DensityMatrix(plus)
print(dm.purity())  # 1.0 → pure state

# Bell state
bell = Statevector([1/np.sqrt(2), 0, 0, 1/np.sqrt(2)])  # |Φ+⟩
print(bell.probabilities_dict())  # {'00': 0.5, '11': 0.5}

# Partial trace: reduced state of qubit 0
rho_A = DensityMatrix(bell).partial_trace([1])  # trace out qubit 1
print(rho_A.purity())  # 0.5 → maximally mixed! (entangled)
```

---

## No-Cloning Theorem

**Statement:** It is impossible to copy an unknown quantum state.

**Proof sketch:**
```
Suppose cloning unitary U exists:
  U(|ψ⟩|0⟩) = |ψ⟩|ψ⟩

For |ψ⟩ = |0⟩:  U(|0⟩|0⟩) = |00⟩
For |ψ⟩ = |1⟩:  U(|1⟩|0⟩) = |11⟩

Now try |ψ⟩ = |+⟩ = (|0⟩+|1⟩)/√2:

U(|+⟩|0⟩) = U((|0⟩+|1⟩)/√2 ⊗ |0⟩)
            = (U|00⟩ + U|10⟩)/√2
            = (|00⟩ + |11⟩)/√2    (by linearity)
            = |Φ+⟩    (Bell state!)

But we wanted |+⟩|+⟩ = (|00⟩+|01⟩+|10⟩+|11⟩)/2 ≠ |Φ+⟩

Contradiction! No cloning unitary exists.
```

**Implications:**
- Cannot copy quantum information (unlike classical bits)
- Enables quantum key distribution (eavesdropping is detectable)
- Quantum error correction must use entanglement, not copying

---

## Decoherence

Real qubits interact with their environment, causing loss of quantum information:

```
Ideal: |ψ⟩ = α|0⟩ + β|1⟩  (coherent superposition)

After decoherence:
  ρ → [|α|²    0  ]    (off-diagonal elements decay → classical mixture)
      [0     |β|²]

Types of noise:
  T₁ (energy relaxation): |1⟩ → |0⟩ decay
  T₂ (dephasing):         phase information lost
  T₂* (inhomogeneous dephasing): faster T₂ from static noise

Typical values (superconducting qubits):
  T₁ ≈ 100–500 μs
  T₂ ≈ 50–200 μs
  Gate time ≈ 10–100 ns

Available gate operations ≈ T₂ / gate_time ≈ 1000–10000
```

This is why error correction (Chapter 14) is critical for large computations.
# Chapter 4: Quantum Gates and Circuits

## Quantum Gates Fundamentals

Quantum gates are unitary operations on qubits:
- **Unitary:** U†U = I → reversible, preserves normalization
- **Linear:** gate on superposition = superposition of gate on basis states
- **No-cloning:** cannot copy a qubit state

---

## Single-Qubit Gates

### Pauli Gates (π rotations):
```
X = [0 1]  (bit flip: |0⟩↔|1⟩, NOT gate analog)
    [1 0]

Y = [0 -i]  (bit+phase flip)
    [i  0]

Z = [1  0]  (phase flip: |0⟩→|0⟩, |1⟩→-|1⟩)
    [0 -1]
```

### Hadamard Gate (superposition):
```
H = 1/√2 [1  1]
          [1 -1]

H|0⟩ = |+⟩ = (|0⟩+|1⟩)/√2
H|1⟩ = |-⟩ = (|0⟩-|1⟩)/√2
H|+⟩ = |0⟩    H|-⟩ = |1⟩    (H is its own inverse!)
```

### Phase Gates:
```
S = [1 0]    T = [1    0  ]    Rz(θ) = [1    0    ]
    [0 i]        [0  e^iπ/4]           [0  e^iθ  ]

S = Z^(1/2)    T = Z^(1/4)    (quarter-turn, eighth-turn phase)

P(λ) = [1     0  ]    (general phase gate — S=P(π/2), T=P(π/4), Z=P(π))
        [0  e^iλ ]
```

### Rotation Gates (Bloch sphere rotations):
```
Rx(θ) = cos(θ/2)I - i·sin(θ/2)X = [cos(θ/2)    -i·sin(θ/2)]
                                    [-i·sin(θ/2)   cos(θ/2) ]

Ry(θ) = cos(θ/2)I - i·sin(θ/2)Y = [cos(θ/2)   -sin(θ/2)]
                                    [sin(θ/2)    cos(θ/2) ]

Rz(θ) = cos(θ/2)I - i·sin(θ/2)Z = [e^(-iθ/2)    0     ]
                                    [0          e^(iθ/2)]
```

### Universal single-qubit gate:
```
U(θ,φ,λ) = [cos(θ/2)           -e^(iλ)·sin(θ/2)]
            [e^(iφ)·sin(θ/2)    e^(i(φ+λ))cos(θ/2)]

Every single-qubit gate = U(θ,φ,λ) for some θ,φ,λ
```

---

## Two-Qubit Gates

### CNOT (Controlled-NOT):
```
CNOT = [1 0 0 0]    Control qubit 0, target qubit 1:
       [0 1 0 0]
       [0 0 0 1]    |00⟩ → |00⟩
       [0 0 1 0]    |01⟩ → |01⟩
                    |10⟩ → |11⟩    (flip target when control is |1⟩)
                    |11⟩ → |10⟩

Circuit notation:
  q0: ─●─       (filled dot = control)
  q1: ─⊕─       (⊕ = target, XOR)
```

### CZ (Controlled-Z):
```
CZ = diag(1, 1, 1, -1)

|11⟩ → -|11⟩    (phase flip only when both are |1⟩)

CZ is symmetric — either qubit can be "control"

Circuit:
  q0: ─●─
  q1: ─●─
```

### SWAP:
```
SWAP = [1 0 0 0]
       [0 0 1 0]
       [0 1 0 0]
       [0 0 0 1]

Swaps the states of two qubits.
SWAP = CNOT₀₁ · CNOT₁₀ · CNOT₀₁    (3 CNOTs)
```

### iSWAP, fSWAP (used in some hardware):
```
iSWAP = [1  0  0  0]    (SWAP with i phase on swapped terms)
        [0  0  i  0]
        [0  i  0  0]
        [0  0  0  1]
```

---

## Three-Qubit Gates

### Toffoli (CCNOT — doubly controlled NOT):
```
Toffoli: Flips target (q2) only when both controls (q0, q1) are |1⟩

|110⟩ → |111⟩
|111⟩ → |110⟩
All other basis states unchanged

Classical NAND gate → Toffoli is universal for classical computation
Circuit:
  q0: ─●─
  q1: ─●─
  q2: ─⊕─
```

### Fredkin (CSWAP):
```
Controlled-SWAP: swaps q1 and q2 when q0 = |1⟩
```

---

## Quantum Circuit Diagrams

Circuits are read left to right (time flows left to right):

```
Example: Creating a Bell state

q0: ─H─●─     H gate on q0, then CNOT with q0 control, q1 target
q1: ───⊕─

Step 1: |00⟩
Step 2: After H on q0: (|0⟩+|1⟩)/√2 ⊗ |0⟩ = (|00⟩+|10⟩)/√2
Step 3: After CNOT: (|00⟩+|11⟩)/√2 = |Φ+⟩ ✓

Example: Quantum teleportation circuit

q0: ─H─●─────M─────────────
q1: ───⊕─M───│─────────────
q2: ─────────│─X─Z─── output
             │ │ │
             └─┘ └── (classical control)
```

---

## Circuit Depth and Width

```
Width  = number of qubits
Depth  = longest path through gates (critical path)
T-count = number of T gates (expensive for fault-tolerant QC)
CNOT-count = number of 2-qubit gates (often the bottleneck)

Goal: minimize depth (speed) and CNOT-count (fidelity)
```

---

## Universal Gate Sets

A gate set is **universal** if any unitary can be approximated to arbitrary precision.

```
Common universal sets:
  {H, T, CNOT}     — most common theoretically
  {Rx, Ry, CNOT}   — hardware-friendly (continuous rotations)
  {H, S, CNOT, T}  — Clifford+T basis

Solovay-Kitaev theorem:
  Any single-qubit unitary can be approximated to ε precision
  using O(log^c(1/ε)) gates from a universal set (c ≈ 2).
```

### Clifford gates (stabilizer operations):
```
{H, S, CNOT}

Properties:
  - Map Pauli operators to Pauli operators
  - Efficiently simulable classically (Gottesman-Knill theorem!)
  - Cannot achieve quantum advantage alone
  - Adding T gate makes set universal

This means: all quantum advantage requires non-Clifford operations.
```

---

## Building Circuits in Qiskit

```python
from qiskit import QuantumCircuit, QuantumRegister, ClassicalRegister
from qiskit.circuit.library import QFT

# Basic circuit
qc = QuantumCircuit(3, 3)  # 3 qubits, 3 classical bits

# Single-qubit gates
qc.h(0)        # Hadamard on qubit 0
qc.x(1)        # X (NOT) on qubit 1
qc.z(2)        # Z (phase flip) on qubit 2
qc.s(0)        # S gate
qc.t(1)        # T gate
qc.rx(0.5, 0)  # Rx(0.5 rad) on qubit 0
qc.ry(1.2, 1)  # Ry
qc.rz(0.3, 2)  # Rz

# Two-qubit gates
qc.cx(0, 1)    # CNOT: control=0, target=1
qc.cz(0, 2)    # CZ
qc.swap(1, 2)  # SWAP

# Three-qubit gates
qc.ccx(0, 1, 2)  # Toffoli
qc.cswap(0, 1, 2)  # Fredkin

# Barriers (for visualization separation)
qc.barrier()

# Measurement
qc.measure([0, 1, 2], [0, 1, 2])

# Draw the circuit
print(qc.draw())

# Draw with matplotlib
qc.draw('mpl')
```

### Composing circuits:
```python
from qiskit import QuantumCircuit

# Bell state circuit
def bell_circuit():
    qc = QuantumCircuit(2)
    qc.h(0)
    qc.cx(0, 1)
    return qc

# GHZ state: (|000⟩ + |111⟩)/√2
def ghz_circuit(n):
    qc = QuantumCircuit(n)
    qc.h(0)
    for i in range(n - 1):
        qc.cx(i, i + 1)
    return qc

# Combine circuits
main = QuantumCircuit(3)
bell = bell_circuit()
main.compose(bell, qubits=[0, 1], inplace=True)
main.cx(1, 2)  # extend entanglement
print(main.draw())
```

---

## Circuit Optimization

```python
from qiskit.compiler import transpile
from qiskit_aer import AerSimulator

qc = QuantumCircuit(3)
# ... build circuit

# Transpile for a specific backend (optimizes gate count)
backend = AerSimulator()
optimized = transpile(qc, backend, optimization_level=3)

print(f"Original depth: {qc.depth()}")
print(f"Optimized depth: {optimized.depth()}")
print(f"CNOT count: {optimized.count_ops().get('cx', 0)}")
```

---

## Common Circuit Patterns

### Phase kickback:
```
Used in almost every quantum algorithm:

q0: ─H─●─H─    q0 is the "control qubit"
q1: ─X─⊕──    q1 is in eigenstate of U

Result: phase of q1's eigenvalue kicks back to q0

If U|ψ⟩ = e^(iφ)|ψ⟩, then:
H|0⟩⊗|ψ⟩ → (1/√2)(|0⟩ + e^(iφ)|1⟩)⊗|ψ⟩

The eigenvalue φ is now encoded in q0's phase.
This is the core of Shor's, phase estimation, and many others.
```

### Quantum parallelism:
```
Prepare superposition of all inputs, apply function once:

q0..n: ─H⊗n─[Uf]─
q_out: ──────[Uf]─

Uf(|x⟩|0⟩) = |x⟩|f(x)⟩

After H⊗n: (1/√N) Σₓ |x⟩
After Uf: (1/√N) Σₓ |x⟩|f(x)⟩   (all f(x) computed at once)

But measurement collapses to ONE x — need interference to extract useful info.
```

---

## Gate Equivalences

```
X = HZH    (X in the Z basis = Z in the X basis)
Z = HXH
CNOT = (I⊗H)(CZ)(I⊗H)
Toffoli = uses 6 CNOTs + single-qubit gates

Key decompositions for hardware:
Most hardware only natively supports {Rx, Ry, Rz, CNOT/CZ}
Compiler (transpiler) converts all other gates to native set
```

---

## Problem-Solving: Circuit Construction

**Pattern:** Build circuits bottom-up from primitives.

```python
# Problem: Implement a circuit that puts n qubits in equal superposition
# Solution: Apply H to all qubits

def equal_superposition(n):
    qc = QuantumCircuit(n)
    qc.h(range(n))  # H on all qubits
    return qc
    # State: (1/√2ⁿ) Σₓ |x⟩   for all x in {0,...,2ⁿ-1}

# Problem: Create W state: (|100⟩ + |010⟩ + |001⟩)/√3
# (useful for distributed quantum computing)
def w_state():
    qc = QuantumCircuit(3)
    # Start with |100⟩
    qc.x(0)
    # Ry rotation to split amplitude
    qc.ry(2 * np.arccos(1/np.sqrt(3)), 0)
    # ... (full construction requires careful angle computation)
    return qc
```
# Chapter 5: Quantum Entanglement

## What is Entanglement?

Entanglement is a quantum correlation between qubits that has no classical analog. Entangled qubits share a joint quantum state that cannot be described independently.

```
Classical correlation (coin pair):
  One shows H → other shows T
  But each coin has a definite state before you look.

Quantum entanglement (Bell state):
  |Φ+⟩ = (|00⟩ + |11⟩)/√2
  Neither qubit has a definite state until measured.
  Measuring one instantly determines the other.
  This holds even if they're light-years apart.
```

---

## Creating Entanglement

```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator

# Create Bell state |Φ+⟩
qc = QuantumCircuit(2, 2)
qc.h(0)        # q0 → (|0⟩+|1⟩)/√2
qc.cx(0, 1)   # CNOT → (|00⟩+|11⟩)/√2

qc.measure([0, 1], [0, 1])
result = AerSimulator().run(qc, shots=1000).result()
print(result.get_counts())
# {'00': ~500, '11': ~500}  — never '01' or '10'!
```

All four Bell states:
```python
def bell_state(phi_plus=True, phi_minus=False, psi_plus=False, psi_minus=False):
    qc = QuantumCircuit(2)
    qc.h(0)
    qc.cx(0, 1)
    
    if phi_minus:
        qc.z(0)     # |Φ-⟩ = (|00⟩ - |11⟩)/√2
    elif psi_plus:
        qc.x(0)     # |Ψ+⟩ = (|01⟩ + |10⟩)/√2
    elif psi_minus:
        qc.x(0); qc.z(0)  # |Ψ-⟩ = (|01⟩ - |10⟩)/√2
    
    return qc
```

---

## GHZ States (Multi-Qubit Entanglement)

Greenberger-Horne-Zeilinger states are the n-qubit generalization of Bell states:

```
|GHZ⟩ = (|00...0⟩ + |11...1⟩)/√2

n=3: (|000⟩ + |111⟩)/√2
```

```python
def ghz_state(n):
    qc = QuantumCircuit(n)
    qc.h(0)
    for i in range(n - 1):
        qc.cx(i, i + 1)
    return qc
```

GHZ states violate Bell inequalities maximally and are used in quantum error correction.

---

## Quantum Teleportation

Quantum teleportation transmits a quantum state using an entangled pair + 2 classical bits.

**Problem:** Alice wants to send quantum state |ψ⟩ = α|0⟩+β|1⟩ to Bob.

**Resources:**
- Alice and Bob share an entangled pair |Φ+⟩
- Alice has a qubit in state |ψ⟩
- Classical communication channel

**Protocol:**
```
                Alice                    Bob
q0 (Alice):  |ψ⟩ = α|0⟩+β|1⟩
q1 (Alice):  ─┐─ (shared Bell pair)
q2 (Bob):    ─┘─

Step 1: Alice applies CNOT(q0→q1) then H(q0)

Full 3-qubit state before measurement:
(α|0⟩+β|1⟩)(|00⟩+|11⟩)/√2

After CNOT(q0,q1) and H(q0):
= (1/2)[|00⟩(α|0⟩+β|1⟩) + |01⟩(α|1⟩+β|0⟩)
       +|10⟩(α|0⟩-β|1⟩) + |11⟩(α|1⟩-β|0⟩)]

Step 2: Alice measures q0, q1 → gets 2 classical bits (m0, m1)

Step 3: Alice sends m0, m1 to Bob (2 classical bits)

Step 4: Bob applies corrections based on (m0, m1):
  00 → I    (nothing)       → Bob has α|0⟩+β|1⟩ ✓
  01 → X    (bit flip)      → Bob has α|0⟩+β|1⟩ ✓
  10 → Z    (phase flip)    → Bob has α|0⟩+β|1⟩ ✓
  11 → XZ   (both)          → Bob has α|0⟩+β|1⟩ ✓
```

```python
from qiskit import QuantumCircuit, ClassicalRegister, QuantumRegister

def teleportation_circuit():
    qr = QuantumRegister(3, 'q')
    cr = ClassicalRegister(2, 'c')
    qc = QuantumCircuit(qr, cr)
    
    # Prepare state to teleport (example: |+⟩)
    qc.h(0)  # q0 will be teleported
    
    qc.barrier()
    
    # Create Bell pair between q1 and q2
    qc.h(1)
    qc.cx(1, 2)
    
    qc.barrier()
    
    # Alice's operations
    qc.cx(0, 1)
    qc.h(0)
    
    qc.barrier()
    
    # Alice measures
    qc.measure(0, 0)
    qc.measure(1, 1)
    
    qc.barrier()
    
    # Bob applies corrections (classically controlled)
    with qc.if_test((cr[1], 1)):
        qc.x(2)
    with qc.if_test((cr[0], 1)):
        qc.z(2)
    
    return qc
```

**Key insight:** No information travels faster than light — Alice must send 2 classical bits to Bob for the protocol to work. The quantum state's content is not transmitted superluminally.

---

## Superdense Coding

The reverse of teleportation: send 2 classical bits using 1 qubit + 1 entangled pair.

```
Alice and Bob share |Φ+⟩ = (|00⟩+|11⟩)/√2
Alice holds qubit A, Bob holds qubit B.

Alice encodes 2 bits by applying one of 4 operations to her qubit A:

  00 → I      → (|00⟩+|11⟩)/√2 = |Φ+⟩
  01 → X      → (|10⟩+|01⟩)/√2 = |Ψ+⟩
  10 → Z      → (|00⟩-|11⟩)/√2 = |Φ-⟩
  11 → XZ     → (|10⟩-|01⟩)/√2 = |Ψ-⟩

Alice sends her qubit A to Bob.
Bob receives both qubits, applies CNOT(A,B) then H(A), then measures.
Bob can distinguish all 4 Bell states → recovers 2 bits.
```

```python
def superdense_coding(bits: str):
    """bits: '00', '01', '10', or '11'"""
    qc = QuantumCircuit(2, 2)
    
    # Create Bell pair
    qc.h(0); qc.cx(0, 1)
    
    # Alice encodes (on qubit 0)
    if bits[1] == '1': qc.x(0)
    if bits[0] == '1': qc.z(0)
    
    # Bob decodes
    qc.cx(0, 1)
    qc.h(0)
    qc.measure([0, 1], [0, 1])
    
    return qc
```

---

## Bell Inequality and EPR Paradox

**Einstein, Podolsky, Rosen (1935):** argued quantum mechanics was incomplete — there must be "hidden variables" that predetermine measurement outcomes.

**Bell's theorem (1964):** ANY hidden variable theory that reproduces QM predictions must be non-local. Tested by Bell inequalities.

**CHSH inequality (measurable form):**
```
Classical (local hidden variables): |S| ≤ 2

S = E(a,b) - E(a,b') + E(a',b) + E(a',b')

where E(a,b) = correlation of measurements along angles a and b.

Quantum mechanics: |S| ≤ 2√2 ≈ 2.828

Experiments consistently show |S| ≈ 2.7 → quantum mechanics is right,
and local hidden variables are ruled out. Quantum nonlocality is real.
```

```python
# CHSH inequality test (conceptual)
import numpy as np

def chsh_quantum(theta_a, theta_b):
    """Quantum correlation for angles theta_a, theta_b"""
    return -np.cos(theta_a - theta_b)

# Optimal angles for maximum CHSH violation
a = 0;     a_prime = np.pi/2
b = np.pi/4;  b_prime = -np.pi/4

S = (chsh_quantum(a, b) - chsh_quantum(a, b_prime)
   + chsh_quantum(a_prime, b) + chsh_quantum(a_prime, b_prime))
print(f"S = {S:.3f}")  # ≈ 2.828 = 2√2 (violates classical bound of 2)
```

---

## Entanglement Entropy

A measure of how entangled a bipartite system is:

```
Von Neumann entropy: S(ρ_A) = -tr(ρ_A log ρ_A)

For Bell state |Φ+⟩:
  ρ_A = I/2    → S = -tr(I/2 · log(I/2)) = log(2) = 1 ebit (maximally entangled)

For product state |ψ₁⟩|ψ₂⟩:
  ρ_A = |ψ₁⟩⟨ψ₁|    → S = 0    (no entanglement)

S ranges from 0 (no entanglement) to log(d) (maximally entangled, d = min(dim A, dim B))
```

```python
from qiskit.quantum_info import entropy, partial_trace, Statevector

# Bell state
bell = Statevector([1/np.sqrt(2), 0, 0, 1/np.sqrt(2)])
rho = bell.to_density_matrix()
rho_A = partial_trace(rho, [1])  # trace out qubit 1
print(f"Entanglement entropy: {entropy(rho_A):.4f}")  # ≈ 1.0 (in bits)
```

---

## Quantum Key Distribution (QKD)

Entanglement enables unconditionally secure key distribution.

**BB84 Protocol (no entanglement needed, but illustrative):**
```
Alice sends qubits in randomly chosen bases: Z ({|0⟩,|1⟩}) or X ({|+⟩,|-⟩})
Bob measures in randomly chosen bases.
They compare bases (publicly) — keep results where bases matched.

Security: Eavesdropper (Eve) measuring qubits collapses superpositions,
introducing detectable errors.

E91 Protocol (Ekert, 1991): Uses Bell state pairs.
Alice and Bob measure their respective qubits along random angles.
Bell inequality violation confirms no eavesdropping.
```

---

## Entanglement in Quantum Algorithms

Entanglement is not just philosophically interesting — it's computationally essential:

```
Grover's algorithm:   Oracle creates entanglement between query and answer
Shor's algorithm:     QFT creates entanglement across register
VQE:                  Ansatz with entangling gates explores correlated states
Quantum teleportation: Entanglement as a communication resource
One-way QC:           Entangled cluster states as computation substrate

Rule of thumb:
  More entanglement ≠ always better
  Algorithms need the RIGHT entanglement structure
  Classically simulable circuits: limited entanglement (low Schmidt rank)
  Quantum advantage: high-entanglement states + interference
```

---

## Problem-Solving: Entanglement Detection

```python
from qiskit.quantum_info import Statevector, entropy, partial_trace

def is_entangled(state_vector, threshold=1e-10):
    """Check if a 2-qubit state is entangled."""
    sv = Statevector(state_vector)
    rho = sv.to_density_matrix()
    rho_A = partial_trace(rho, [1])
    ent_entropy = entropy(rho_A)
    return ent_entropy > threshold, ent_entropy

# Test
bell = [1/np.sqrt(2), 0, 0, 1/np.sqrt(2)]
product = [1/np.sqrt(2), 1/np.sqrt(2), 0, 0]  # |+⟩|0⟩

print(is_entangled(bell))     # (True, 1.0)
print(is_entangled(product))  # (False, 0.0)
```
# Chapter 6: Quantum Measurement

## The Measurement Postulate

Measurement is the act of extracting classical information from a quantum system. It is irreversible and probabilistic.

```
State before measurement: |ψ⟩ = Σᵢ αᵢ|i⟩

Measure in basis {|i⟩}:
  Outcome i occurs with probability P(i) = |αᵢ|² = |⟨i|ψ⟩|²
  After outcome i: state collapses to |i⟩ (normalized)

This is the Born rule — the fundamental postulate of quantum mechanics.
```

---

## Projective Measurement

A projective measurement is described by a set of projection operators {Pᵢ}:

```
Pᵢ = |i⟩⟨i|    (project onto basis state |i⟩)

Properties:
  ΣᵢPᵢ = I        (completeness)
  PᵢPⱼ = δᵢⱼPᵢ   (orthogonality)
  Pᵢ² = Pᵢ        (idempotent)

Probability of outcome i: P(i) = ⟨ψ|Pᵢ|ψ⟩ = |⟨i|ψ⟩|²
Post-measurement state: Pᵢ|ψ⟩ / √P(i)
```

---

## Measuring in Different Bases

Measurement outcome depends on the basis chosen:

```
State |+⟩ = (|0⟩+|1⟩)/√2

Measure in Z basis ({|0⟩,|1⟩}):
  P(0) = |⟨0|+⟩|² = |1/√2|² = 1/2
  P(1) = |⟨1|+⟩|² = 1/2
  Random outcome!

Measure in X basis ({|+⟩,|-⟩}):
  P(+) = |⟨+|+⟩|² = 1
  P(-) = |⟨-|+⟩|² = 0
  Deterministic outcome: always get |+⟩
```

### To measure in X basis:
```
Apply H before measuring in Z basis:
H|+⟩ = |0⟩ → always measure 0

In general: to measure in basis {H|0⟩, H|1⟩}, apply H† = H before standard measurement.
```

```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator

# Measure |+⟩ in different bases
qc_z = QuantumCircuit(1, 1)
qc_z.h(0)              # prepare |+⟩
qc_z.measure(0, 0)     # measure in Z basis → random

qc_x = QuantumCircuit(1, 1)
qc_x.h(0)              # prepare |+⟩
qc_x.h(0)              # change to X basis (apply H before measuring)
qc_x.measure(0, 0)     # now measures in X basis → deterministic 0

sim = AerSimulator()
print(sim.run(qc_z, shots=1000).result().get_counts())  # {'0': ~500, '1': ~500}
print(sim.run(qc_x, shots=1000).result().get_counts())  # {'0': 1000}
```

---

## Expectation Values

The expected value of an observable O in state |ψ⟩:
```
⟨O⟩ = ⟨ψ|O|ψ⟩

For O = Z (Pauli-Z), measuring |0⟩ gives +1, |1⟩ gives -1:
  ⟨Z⟩ for |0⟩ = ⟨0|Z|0⟩ = +1
  ⟨Z⟩ for |1⟩ = ⟨1|Z|1⟩ = -1
  ⟨Z⟩ for |+⟩ = ⟨+|Z|+⟩ = 0    (equal superposition → average 0)

Multiple measurements needed to estimate expectation value:
  ⟨O⟩ ≈ (1/N) Σᵢ oᵢ    (average of N measurement outcomes)
  
Error scales as 1/√N — need many shots for precision.
```

```python
from qiskit.quantum_info import SparsePauliOp, Statevector

psi = Statevector.from_label('+')  # |+⟩
Z = SparsePauliOp('Z')
H_op = SparsePauliOp('X')  # Pauli X as observable

print(psi.expectation_value(Z))   # ≈ 0 + 0i
print(psi.expectation_value(H_op)) # ≈ 1 + 0i  (|+⟩ is eigenstate of X with +1)
```

---

## Partial Measurement

Measuring only some qubits of a multi-qubit system:

```
2-qubit state: |ψ⟩ = (1/2)(|00⟩ + |01⟩ + |10⟩ + |11⟩) = |+⟩|+⟩

Measure qubit 0 only:
  P(0) = |⟨0|_0 |ψ⟩|² = (|00⟩+|01⟩)/√2 → norm² = 1/2
  P(1) = 1/2

If outcome is 0: post-measurement state = (|00⟩+|01⟩)/√2 = |0⟩|+⟩

For entangled state |Φ+⟩ = (|00⟩+|11⟩)/√2:
  Measure qubit 0:
    P(0) = 1/2, post-state = |00⟩    (qubit 1 collapses to |0⟩)
    P(1) = 1/2, post-state = |11⟩    (qubit 1 collapses to |1⟩)
```

```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator

# Partial measurement of entangled state
qc = QuantumCircuit(2, 1)
qc.h(0)
qc.cx(0, 1)   # Bell state

qc.measure(0, 0)  # Measure only qubit 0

# After measurement, qubit 1 is in a definite state
# (demonstrated by measuring qubit 1 conditioned on result)
result = AerSimulator().run(qc, shots=1000).result()
print(result.get_counts())  # {'0': ~500, '1': ~500}
```

---

## Mid-Circuit Measurement and Reset

Modern quantum computers support measuring qubits mid-circuit (not just at the end):

```python
from qiskit import QuantumCircuit, ClassicalRegister, QuantumRegister

qr = QuantumRegister(3, 'q')
cr = ClassicalRegister(2, 'c')
qc = QuantumCircuit(qr, cr)

# Prepare something
qc.h(0)
qc.cx(0, 1)

# Mid-circuit measurement of qubit 0
qc.measure(0, 0)

# Reset qubit 0 (back to |0⟩)
qc.reset(0)

# Continue using qubit 0
qc.h(0)

# Classical control: apply gate conditioned on measurement outcome
with qc.if_test((cr[0], 1)):
    qc.x(1)

qc.measure([1, 2], [0, 1])
```

---

## POVM (Positive Operator-Valued Measures)

POVMs generalize projective measurement — used when measurement is not perfectly projective (e.g., noisy hardware, ancilla-mediated measurement):

```
A POVM is a set {Mᵢ} satisfying:
  Mᵢ ≥ 0    (positive semidefinite)
  ΣᵢMᵢ = I  (completeness)

Probability of outcome i: P(i) = ⟨ψ|Mᵢ|ψ⟩ = tr(Mᵢ|ψ⟩⟨ψ|)

Projective measurement is special case where Mᵢ = Pᵢ = |i⟩⟨i|

POVM example: Unambiguous state discrimination
Distinguish |0⟩ and |+⟩ with 3 outcomes:
  M₀ = (|1⟩⟨1|) / (1 + 1/√2)      → outcome 0 means "definitely was |+⟩"
  M₊ = |-⟩⟨-| / (1 + 1/√2)        → outcome 1 means "definitely was |0⟩"
  M? = I - M₀ - M₊                  → outcome 2 means "inconclusive"
```

---

## Quantum State Tomography

Reconstruct an unknown quantum state by measuring many copies in different bases:

```python
from qiskit.quantum_info import state_fidelity, Statevector
from qiskit_aer import AerSimulator
from qiskit import QuantumCircuit
import numpy as np

# For a 1-qubit state, measure in X, Y, Z bases
# Expectation values:
#   ⟨X⟩ = 2 Re(α*β)     (off-diagonal real)
#   ⟨Y⟩ = 2 Im(α*β)     (off-diagonal imaginary)
#   ⟨Z⟩ = |α|² - |β|²   (diagonal difference)

# Reconstruct density matrix:
# ρ = (I + ⟨X⟩X + ⟨Y⟩Y + ⟨Z⟩Z) / 2

def tomography_1qubit(counts_x, counts_y, counts_z, shots):
    """Reconstruct 1-qubit state from Pauli measurements."""
    exp_x = (counts_x.get('0', 0) - counts_x.get('1', 0)) / shots
    exp_y = (counts_y.get('0', 0) - counts_y.get('1', 0)) / shots
    exp_z = (counts_z.get('0', 0) - counts_z.get('1', 0)) / shots
    
    I = np.eye(2)
    X = np.array([[0,1],[1,0]])
    Y = np.array([[0,-1j],[1j,0]])
    Z = np.array([[1,0],[0,-1]])
    
    rho = (I + exp_x * X + exp_y * Y + exp_z * Z) / 2
    return rho
```

---

## Measurement in Algorithms

```
Grover's Search:
  Measure at the end → marked state appears with high probability

Shor's Algorithm:
  Measure quantum register → get a multiple of 2ⁿ/r
  Classical post-processing extracts period r

Quantum Phase Estimation:
  Measure ancilla register → encodes eigenvalue as binary

VQE:
  Measure Hamiltonian terms one by one
  Estimate energy as sum of expectation values

Key principle: Design the algorithm so measurement gives useful information.
Constructive interference → correct answer has high probability before measuring.
```

---

## The Measurement Problem

Why does measurement collapse the wave function? This is an unsolved foundational question:

```
Copenhagen interpretation: Collapse is fundamental. The wave function describes
what we know, not objective reality.

Many-worlds interpretation: No collapse — universe branches at each measurement.
Each outcome occurs in a different branch.

Decoherence: Interaction with environment effectively collapses the state.
Superposition still exists but becomes practically unobservable.

For quantum computing purposes:
  We use the Born rule operationally — it correctly predicts outcomes.
  The interpretation doesn't affect algorithm design or hardware.
```

---

## Problem-Solving: Measurement Design

**Pattern:** Decide what to measure and when to extract maximum information.

```python
# Problem: Determine if a 2-qubit state is |00⟩ or |11⟩ (Bell state component)
# Without destroying the state entirely

from qiskit import QuantumCircuit

def measure_parity():
    """Measures X⊗X parity without collapsing into |00⟩ or |11⟩."""
    qc = QuantumCircuit(3, 1)  # 2 data + 1 ancilla
    
    # Prepare Bell state
    qc.h(0); qc.cx(0, 1)
    
    # Parity measurement: CNOT from both data qubits to ancilla
    qc.cx(0, 2)
    qc.cx(1, 2)
    
    # Measure ancilla only
    qc.measure(2, 0)
    
    # If ancilla = 0: parity is even (|00⟩ or |11⟩ component)
    # If ancilla = 1: parity is odd  (|01⟩ or |10⟩ component)
    # Data qubits not disturbed within their parity sector
    return qc
```
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
# Chapter 8: Deutsch-Jozsa and Simon's Algorithm

## Deutsch's Problem (1985)

The first quantum algorithm to show exponential speedup (for a promise problem).

### Problem:
```
Given: Black-box function f:{0,1} → {0,1}
Promise: f is either:
  - Constant: f(0) = f(1)   (both same output)
  - Balanced: f(0) ≠ f(1)   (different outputs)
Determine which case.

Classical: 2 queries (must check both inputs)
Quantum:   1 query
```

### Solution (Deutsch's Algorithm):
```
Initial: |0⟩|1⟩
Step 1:  Apply H to both: |+⟩|-⟩ = (|0⟩+|1⟩)/√2 · (|0⟩-|1⟩)/√2
Step 2:  Apply Uf: |x⟩|y⟩ → |x⟩|y⊕f(x)⟩

After Uf (using phase kickback with |-⟩):
  = ((-1)^f(0)|0⟩ + (-1)^f(1)|1⟩)/√2 · |−⟩

Step 3: Apply H to first qubit:
  If f constant: (±|0⟩) → measurement gives 0
  If f balanced: (±|1⟩) → measurement gives 1

Step 4: Measure first qubit:
  0 → constant
  1 → balanced
```

```python
from qiskit import QuantumCircuit, QuantumRegister, ClassicalRegister
from qiskit_aer import AerSimulator

def deutsch_circuit(f_type: str):
    """
    f_type: 'constant_0', 'constant_1', 'balanced_id', 'balanced_not'
    """
    qr = QuantumRegister(2, 'q')
    cr = ClassicalRegister(1, 'c')
    qc = QuantumCircuit(qr, cr)
    
    # Initialize: |01⟩
    qc.x(1)
    
    # Apply H to both
    qc.h(0)
    qc.h(1)
    
    # Oracle
    if f_type == 'constant_0':
        pass                    # f(x) = 0: identity
    elif f_type == 'constant_1':
        qc.x(1)                 # f(x) = 1: flip ancilla (but phase cancels → still constant)
    elif f_type == 'balanced_id':
        qc.cx(0, 1)             # f(x) = x: CNOT
    elif f_type == 'balanced_not':
        qc.cx(0, 1)             # f(x) = NOT x
        qc.x(1)
    
    # Apply H to first qubit
    qc.h(0)
    
    # Measure
    qc.measure(0, 0)
    
    return qc

sim = AerSimulator()
for f_type in ['constant_0', 'constant_1', 'balanced_id', 'balanced_not']:
    qc = deutsch_circuit(f_type)
    counts = sim.run(qc, shots=100).result().get_counts()
    verdict = 'constant' if '0' in counts else 'balanced'
    print(f"{f_type}: {counts} → {verdict}")
# constant_0: {'0': 100} → constant
# constant_1: {'0': 100} → constant
# balanced_id: {'1': 100} → balanced
# balanced_not: {'1': 100} → balanced
```

---

## Deutsch-Jozsa Algorithm (1992)

Generalizes Deutsch's problem to n-bit functions.

### Problem:
```
Given: Black-box f:{0,1}ⁿ → {0,1}
Promise: f is either:
  - Constant: f(x) = c for all x
  - Balanced: f(x)=0 for exactly 2^(n-1) inputs, f(x)=1 for the other half
Determine which case.

Classical (deterministic): up to 2^(n-1)+1 queries (exponential worst case)
Classical (randomized): O(1) queries with bounded error — but Deutsch-Jozsa is exact
Quantum: 1 query (exact)
```

### Circuit:
```
|0⟩^⊗n ─H⊗n─[  ]─H⊗n─ Measure
|1⟩    ─H───[ Uf]─────

Step 1: Apply H to all n+1 qubits
  → (1/√2ⁿ) Σₓ |x⟩ ⊗ |−⟩

Step 2: Apply phase oracle Uf
  → (1/√2ⁿ) Σₓ (-1)^f(x)|x⟩ ⊗ |−⟩

Step 3: Apply H⊗n to first n qubits
  → Σₖ [(1/2ⁿ) Σₓ (-1)^(f(x)+x·k)] |k⟩

Step 4: Measure first n qubits:
  If f is constant: all amplitude at |00...0⟩ → measure 0...0
  If f is balanced: zero amplitude at |00...0⟩ → NEVER measure 0...0
```

**Mathematical proof:**
```
Amplitude at |0...0⟩ after H⊗n:
  = (1/2ⁿ) Σₓ (-1)^f(x)

If constant:    (1/2ⁿ) · 2ⁿ · (-1)^c = ±1   (large amplitude)
If balanced:    (1/2ⁿ) · 0 = 0               (exactly zero)

So measuring |0...0⟩ → constant
   measuring anything else → balanced
```

```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator
import numpy as np

def deutsch_jozsa(n: int, oracle_type: str = 'balanced'):
    """
    n: number of input qubits
    oracle_type: 'constant_0', 'constant_1', or 'balanced'
    """
    qc = QuantumCircuit(n + 1, n)
    
    # Initialize last qubit as |1⟩
    qc.x(n)
    
    # Apply H to all qubits
    qc.h(range(n + 1))
    
    qc.barrier()
    
    # Oracle
    if oracle_type == 'constant_0':
        pass
    elif oracle_type == 'constant_1':
        qc.x(n)
    elif oracle_type == 'balanced':
        # Simple balanced oracle: CNOT from each input qubit to output
        for i in range(n):
            qc.cx(i, n)
    
    qc.barrier()
    
    # Apply H to input qubits
    qc.h(range(n))
    
    # Measure input qubits
    qc.measure(range(n), range(n))
    
    return qc

# Test with n=4
sim = AerSimulator()
for oracle in ['constant_0', 'constant_1', 'balanced']:
    qc = deutsch_jozsa(4, oracle)
    counts = sim.run(qc, shots=100).result().get_counts()
    # Constant → only '0000'; Balanced → never '0000'
    verdict = 'constant' if '0000' in counts else 'balanced'
    print(f"{oracle}: {verdict} | counts={counts}")
```

---

## Simon's Algorithm (1994)

Simon's problem: the first problem with proven exponential quantum speedup.

### Problem:
```
Given: f:{0,1}ⁿ → {0,1}ⁿ satisfying: f(x) = f(x⊕s) for some hidden s
Promise: either s = 0ⁿ (one-to-one) or s ≠ 0ⁿ (two-to-one)
Find s.

Classical: Ω(2^(n/2)) queries (birthday paradox)
Quantum: O(n) queries
Exponential separation!
```

### Algorithm:
```
Repeat O(n) times:
  1. Prepare (1/√2ⁿ) Σₓ |x⟩|0⟩
  2. Apply Uf: (1/√2ⁿ) Σₓ |x⟩|f(x)⟩
  3. Measure second register → gets some value y = f(x₀)
     First register collapses to: (|x₀⟩ + |x₀⊕s⟩)/√2
  4. Apply H⊗n to first register:
     → (1/√2ⁿ) Σₖ (-1)^(k·x₀) [1 + (-1)^(k·s)] |k⟩
     Amplitude is zero unless k·s = 0 (mod 2)!
  5. Measure first register → get k with k·s = 0 (mod 2)

Collect n-1 linearly independent equations k·s = 0
Solve the system (classical Gaussian elimination) → find s
```

```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator
import numpy as np

def simon_circuit(n: int, s: str):
    """
    Simon's algorithm circuit
    s: the hidden period as binary string (length n)
    """
    qc = QuantumCircuit(2 * n, n)
    
    # Apply H to first n qubits
    qc.h(range(n))
    
    # Oracle: f(x) = f(x XOR s)
    # Simple construction: copy x to output, XOR with s when needed
    for i in range(n):
        qc.cx(i, i + n)
    
    # Apply XOR with s (for non-trivial period)
    if int(s, 2) != 0:
        # This creates a period — simplified oracle
        for i, bit in enumerate(s):
            if bit == '1':
                qc.cx(0, n + i)  # XOR with s when first bit is 1
    
    # Apply H to first n qubits
    qc.h(range(n))
    
    # Measure first n qubits
    qc.measure(range(n), range(n))
    
    return qc

def solve_simon_system(equations: list[str]) -> str:
    """Classical post-processing: solve k·s=0 system via Gaussian elimination."""
    n = len(equations[0])
    # Convert to matrix over GF(2)
    matrix = [[int(b) for b in eq] for eq in equations]
    
    # Gaussian elimination over GF(2)
    pivot_col = {}
    row = 0
    for col in range(n):
        pivot = None
        for r in range(row, len(matrix)):
            if matrix[r][col] == 1:
                pivot = r
                break
        if pivot is None:
            continue
        matrix[row], matrix[pivot] = matrix[pivot], matrix[row]
        pivot_col[col] = row
        for r in range(len(matrix)):
            if r != row and matrix[r][col] == 1:
                matrix[r] = [(matrix[r][i] ^ matrix[row][i]) for i in range(n)]
        row += 1
    
    # Back-substitute to find s
    s = [0] * n
    free_vars = [c for c in range(n) if c not in pivot_col]
    if free_vars:
        s[free_vars[0]] = 1  # set free variable to 1
        for col, row_idx in pivot_col.items():
            s[col] = matrix[row_idx][free_vars[0]]
    
    return ''.join(map(str, s))

# Run Simon's algorithm for s = '110'
n = 3
s = '110'
sim = AerSimulator()
equations = []
for _ in range(n * 3):  # collect more equations than needed
    qc = simon_circuit(n, s)
    result = sim.run(qc, shots=1).result()
    k = list(result.get_counts().keys())[0]
    if k != '0' * n:  # skip trivial equation
        equations.append(k)
    if len(equations) >= n - 1:
        break

print(f"Equations: {equations}")
recovered_s = solve_simon_system(equations)
print(f"Recovered s = {recovered_s}")  # Should be '110'
```

---

## Why Simon's Algorithm Matters

Simon's algorithm inspired Shor's algorithm. The structure is almost identical:

```
Simon's:
  Hidden period under XOR: f(x) = f(x⊕s)
  Quantum gives linear equations k·s = 0

Shor's:
  Hidden period under multiplication: f(x) = aˣ mod N
  QFT gives equations → period r
  Period r → factors of N

Simon's → Shor's: Replace XOR (cyclic group Z₂ⁿ) with multiplication (cyclic group Zᴺ*)
```

---

## Bernstein-Vazirani Algorithm

A special case of Deutsch-Jozsa — find a hidden string using 1 query:

```
Problem: f(x) = s·x (mod 2) for hidden string s
Classical: n queries (one bit of s at a time)
Quantum: 1 query

Circuit: Same as Deutsch-Jozsa!
After H⊗n, Oracle, H⊗n:
Measurement directly gives s

This shows quantum algorithms can learn hidden linear structure in 1 shot.
```

```python
def bernstein_vazirani(n: int, s: str):
    """Find hidden string s where f(x) = s·x mod 2"""
    qc = QuantumCircuit(n + 1, n)
    qc.x(n)
    qc.h(range(n + 1))
    
    # Oracle: f(x) = s·x → CNOT for each bit of s
    for i, bit in enumerate(s):
        if bit == '1':
            qc.cx(i, n)
    
    qc.h(range(n))
    qc.measure(range(n), range(n))
    return qc

sim = AerSimulator()
s = '10110'
qc = bernstein_vazirani(len(s), s)
counts = sim.run(qc, shots=1).result().get_counts()
print(f"Hidden s = {s}, recovered = {list(counts.keys())[0]}")  # Same!
```

---

## Summary: First Generation Quantum Algorithms

| Algorithm | Problem | Classical | Quantum | Speedup |
|---|---|---|---|---|
| Deutsch | 1-bit function constant/balanced | 2 queries | 1 query | 2× |
| Deutsch-Jozsa | n-bit function constant/balanced | 2^(n-1)+1 | 1 query | Exponential |
| Bernstein-Vazirani | Find hidden string s in s·x | n queries | 1 query | n× |
| Simon's | Find hidden XOR period | Ω(2^(n/2)) | O(n) | Exponential |

These algorithms demonstrate quantum advantage for artificial problems, paving the way for Shor's and Grover's algorithms on real-world problems.
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
# Chapter 10: Quantum Fourier Transform

## Discrete Fourier Transform (Classical Review)

```
DFT: x → x̂   where x̂ₖ = (1/√N) Σⱼ xⱼ · ωₙʲᵏ

ωₙ = e^(2πi/N)    (N-th root of unity)

Classical FFT: O(N log N)
```

---

## Quantum Fourier Transform

The QFT performs the DFT on quantum amplitudes — exponentially faster:

```
QFT: |j⟩ → (1/√N) Σₖ e^(2πijk/N) |k⟩

For N = 2ⁿ:
QFT|j⟩ = (1/√2ⁿ) Σₖ e^(2πijk/2ⁿ) |k⟩

Classical FFT: O(N log N) = O(2ⁿ · n)
Quantum QFT:  O(n²)  (exponentially faster!)
```

---

## QFT Product Formula

For n-bit number j = j₁j₂...jₙ (binary), the QFT has a beautiful product structure:

```
QFT|j₁j₂...jₙ⟩ = (1/√2)(|0⟩ + e^(2πi·0.jₙ)|1⟩)
                ⊗ (1/√2)(|0⟩ + e^(2πi·0.jₙ₋₁jₙ)|1⟩)
                ⊗ ...
                ⊗ (1/√2)(|0⟩ + e^(2πi·0.j₁j₂...jₙ)|1⟩)

where 0.j₁j₂...jₙ = j₁/2 + j₂/4 + ... + jₙ/2ⁿ  (binary fraction)

This factored form gives us the efficient circuit!
```

---

## QFT Circuit

```
For 3 qubits (j = j₁j₂j₃):

j₁: ─H─P(π/2)─P(π/4)─────────────── SWAP ──── output bit 0
j₂: ───────────────H──P(π/2)──────── SWAP ──── output bit 1
j₃: ──────────────────────────H ─────────────── output bit 2
     ↑           ↑            ↑
     H    R₂ gate  R₃ gate   H

Phase gate Rₖ = [1    0        ]    (phase = 2π/2ᵏ)
               [0    e^(2πi/2ᵏ)]
```

```python
from qiskit import QuantumCircuit
import numpy as np

def qft_circuit(n: int, inverse: bool = False):
    """Build QFT circuit for n qubits."""
    qc = QuantumCircuit(n)
    
    def qft_rotation(qc, n):
        if n == 0:
            return
        n -= 1
        qc.h(n)
        for qubit in range(n):
            qc.cp(np.pi / 2**(n - qubit), qubit, n)
        qft_rotation(qc, n)
    
    qft_rotation(qc, n)
    
    # Swap to get correct bit order
    for qubit in range(n // 2):
        qc.swap(qubit, n - qubit - 1)
    
    if inverse:
        return qc.inverse()
    return qc

# Build 3-qubit QFT
qft = qft_circuit(3)
print(qft.draw())
print(f"Gate count: {qft.count_ops()}")
# H: 3, CP: 3, SWAP: 1 — only O(n²) gates!
```

---

## Verifying QFT

```python
from qiskit.quantum_info import Statevector, Operator
import numpy as np

n = 3
N = 2**n

# Classical DFT matrix
omega = np.exp(2j * np.pi / N)
DFT = np.array([[omega**(j*k) for k in range(N)] for j in range(N)]) / np.sqrt(N)

# QFT unitary
qft = qft_circuit(n)
QFT_matrix = Operator(qft).data

# Check they match (up to qubit ordering conventions)
print(np.allclose(np.abs(DFT), np.abs(QFT_matrix)))  # Should be True

# Apply QFT to |5⟩ = |101⟩
state_5 = Statevector.from_label('101')
qft_state = state_5.evolve(qft)
print("QFT|5⟩ amplitudes (magnitudes):")
print(np.abs(qft_state.data))  # All equal = 1/√8 (uniform superposition phase-encoded)
```

---

## Inverse QFT

The inverse QFT (IQFT) is the complex conjugate of QFT:

```
IQFT = QFT†

Used in: Phase estimation (measuring eigenphases)
         Shor's algorithm (period extraction)

Circuit: Run QFT circuit in reverse with Rₖ replaced by Rₖ†
```

```python
iqft = qft_circuit(n, inverse=True)
```

---

## Period Finding with QFT

The key application: find period r of a function f(x) = f(x+r).

```
If f has period r, then Σₓ f(x)|x⟩ has support on multiples of r.
QFT maps this to: support on multiples of N/r

Measure → get j·N/r for some random j
→ use continued fractions to extract r from j·N/r

This is the core of Shor's algorithm!
```

```
State after computing f: (1/√N) Σₓ |x⟩|f(x)⟩
Measure second register → get f(x₀)
First register collapses to: (1/√(N/r)) Σⱼ |x₀ + j·r⟩

Apply QFT: gets sharp peaks at multiples of N/r
Measure → j·N/r for random integer j
Classical: gcd(j·N/r, N) computations find r
```

---

## QFT Gate Complexity

```
n-qubit QFT:
  H gates: n
  Controlled phase gates: n(n-1)/2
  SWAP gates: ⌊n/2⌋
  Total: O(n²) gates

Compared to classical FFT: O(N log N) = O(2ⁿ · n)

Speedup: exponential!

But: QFT computes the transform on the AMPLITUDES, not on classical data.
You can't directly read off all N transform coefficients (measuring collapses state).
The usefulness comes from USING the QFT state in further computation.
```

---

## Phase Estimation Preview

QFT is the key ingredient in quantum phase estimation:

```
Given: Unitary U with eigenstate |u⟩: U|u⟩ = e^(2πiφ)|u⟩
Goal: Estimate φ (the phase)

Circuit:
  Ancilla: ─H⊗t─[cU^1 cU^2 ... cU^2^(t-1)]─IQFT─ Measure
  |u⟩:    ──────────────────────────────────────────────────

After IQFT, ancilla is in state |φ̃⟩ (binary representation of φ)
Measurement gives t-bit approximation of φ.
```

---

## Common QFT Applications

| Application | How QFT is Used |
|---|---|
| Period finding (Shor's) | Transform periodic superposition to peaks at multiples of N/r |
| Phase estimation (QPE) | IQFT to read off eigenphase as binary fraction |
| Quantum simulation | Time evolution as phase kicks, QFT to read energy |
| Hidden subgroup problem | Generalized QFT over groups |
| Quantum machine learning | Quantum signal processing |

---

## Problem-Solving: QFT-based Algorithms

**Pattern:** Whenever you see "find the period/frequency of something," think QFT.

```python
# Demonstrate QFT detects periodicity
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator
from qiskit.quantum_info import Statevector
import numpy as np

# Create state with period r=2 in 4-qubit register
# Periodic state: |0⟩+|2⟩+|4⟩+|6⟩+|8⟩+|10⟩+|12⟩+|14⟩ (all even → period 2 in 4 bits)
n = 4
periodic_state = np.zeros(2**n)
r = 2
for k in range(2**n // r):
    periodic_state[k * r] = 1
periodic_state /= np.linalg.norm(periodic_state)

sv = Statevector(periodic_state)
qft = qft_circuit(n)
qft_sv = sv.evolve(qft)

# Peaks should appear at multiples of N/r = 16/2 = 8
probs = np.abs(qft_sv.data)**2
print("QFT probabilities (should peak at 0 and 8):")
for i, p in enumerate(probs):
    if p > 0.01:
        print(f"  |{i}⟩: {p:.3f}")
# |0⟩: 0.5, |8⟩: 0.5 → period r = N/8 = 2 ✓
```
# Chapter 11: Shor's Factoring Algorithm

## Why Shor's Algorithm Matters

```
RSA encryption (used in HTTPS, banking, email):
  Security relies on: factoring n = p·q (two large primes) is HARD

Best classical: General Number Field Sieve
  Time: O(exp(n^(1/3) · (log n)^(2/3)))  — superpolynomial in bit length

Shor's algorithm (1994):
  Time: O(n³) — polynomial!  (n = number of bits in the number)

Impact: A large enough quantum computer would break RSA.
Current cryptography: migrating to post-quantum standards (NIST 2024).
```

---

## Problem Reduction: Factoring → Period Finding

Shor's insight: factoring reduces to **order finding**.

```
STEP 1: Choose random a with 1 < a < N and gcd(a, N) = 1
        (If gcd(a,N) ≠ 1, you've found a factor already — done!)

STEP 2: Find the ORDER r of a modulo N:
        r = smallest positive integer such that aʳ ≡ 1 (mod N)

STEP 3: If r is odd or aʳ/² ≡ -1 (mod N): restart with new a

STEP 4: Otherwise:
        gcd(aʳ/² - 1, N) and gcd(aʳ/² + 1, N) are factors of N!

WHY? Because aʳ ≡ 1 (mod N) means:
  (aʳ/² - 1)(aʳ/² + 1) ≡ 0 (mod N)
  So N divides this product → factors appear in the GCDs.

QUANTUM PART: Finding r is the hard part — use quantum period finding.
CLASSICAL PART: Steps 1, 3, 4 are fast classical operations.
```

---

## Classical Example

```
Factor N = 15, choose a = 7

Order finding: 7ˢ mod 15 for s = 1, 2, 3, ...
  7¹ = 7
  7² = 49 = 3·15 + 4 → 4
  7³ = 7·4 = 28 = 15 + 13 → 13
  7⁴ = 7·13 = 91 = 6·15 + 1 → 1  ← period r = 4

r = 4 (even, and 7² = 4 ≠ -1 ≡ 14 mod 15) ✓

gcd(7² - 1, 15) = gcd(48, 15) = gcd(3·15+3, 15) = gcd(3, 15) = 3  ✓
gcd(7² + 1, 15) = gcd(50, 15) = gcd(3·15+5, 15) = gcd(5, 15) = 5  ✓

Factors: 3 and 5 (15 = 3 × 5) ✓
```

---

## Quantum Period Finding

The quantum algorithm finds r by:

```
1. Prepare superposition: (1/√2ⁿ) Σₓ |x⟩|0⟩

2. Compute modular exponentiation: |x⟩|0⟩ → |x⟩|aˣ mod N⟩

3. Measure second register → collapse to some value aˣ⁰ mod N
   First register collapses to: (1/√(2ⁿ/r)) Σⱼ |x₀ + j·r⟩  (periodic with period r)

4. Apply QFT to first register:
   → peaks at multiples of 2ⁿ/r

5. Measure first register → get j·(2ⁿ/r) for some integer j

6. Classical: continued fractions algorithm on j/2ⁿ gives r
```

---

## Circuit Architecture

```
n qubits for x register (n = 2·log₂(N), need ~2n input qubits)
n qubits for result register

q₀,...,q_{2n-1}: ─H⊗2n──[  modular  ]──QFT⁻¹── Measure
q_{2n},...,q_{3n-1}: ─────[exponent  ]────────── Measure

Controlled unitary: U_a,j: |x⟩|y⟩ → |x⟩|y·aˣ mod N⟩

Key expense: Modular exponentiation circuit (O(n³) gates using repeated squaring)
```

---

## Implementation for Small N

```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator
from qiskit.circuit.library import QFT
from fractions import Fraction
import numpy as np

def quantum_period_finding(N: int, a: int):
    """
    Find period r of f(x) = a^x mod N using quantum circuit.
    For small N only (demonstration purposes).
    """
    n = int(np.ceil(np.log2(N)))
    t = 2 * n  # ancilla bits for phase estimation
    
    qc = QuantumCircuit(t + n, t)
    
    # Initialize target register to |1⟩
    qc.x(t)  # first qubit of target register
    
    # Hadamard on counting register
    qc.h(range(t))
    
    # Controlled modular multiplication: a^(2^j) mod N
    for j in range(t):
        # Implement controlled-U^(2^j) where U|y⟩ = |a*y mod N⟩
        # For real implementation, use modular exponentiation circuit
        # Here we use a classical simulation trick
        _controlled_modular_mult(qc, a, j, N, t, n)
    
    # Inverse QFT on counting register
    qc.compose(QFT(t, inverse=True), qubits=range(t), inplace=True)
    
    # Measure counting register
    qc.measure(range(t), range(t))
    
    return qc

def shor_classical(N: int, measurement: int, t: int) -> int:
    """Classical post-processing: extract period from measurement."""
    # Use continued fractions to find r/N approximation
    phase = measurement / (2**t)
    frac = Fraction(phase).limit_denominator(N)
    r = frac.denominator
    return r

def factor_with_shor(N: int) -> tuple:
    """Full Shor's algorithm (classical simulation for small N)."""
    import random
    from math import gcd
    
    for _ in range(100):  # retry with different a
        a = random.randint(2, N - 1)
        g = gcd(a, N)
        if g != 1:
            return g, N // g  # Lucky! gcd gives factor directly
        
        # Find order r of a modulo N (classically for demo)
        r = 1
        val = a
        while val != 1:
            val = (val * a) % N
            r += 1
        
        if r % 2 == 1:
            continue  # odd period, retry
        
        half_power = pow(a, r // 2, N)
        if half_power == N - 1:  # aʳ/² ≡ -1 mod N
            continue
        
        p = gcd(half_power - 1, N)
        q = gcd(half_power + 1, N)
        
        if 1 < p < N:
            return p, N // p
        if 1 < q < N:
            return q, N // q
    
    return None, None

# Demonstrate Shor's classical reduction
print(factor_with_shor(15))   # (3, 5) or (5, 3)
print(factor_with_shor(21))   # (3, 7) or (7, 3)
print(factor_with_shor(35))   # (5, 7)
```

---

## Simplified 4-Qubit Demo (N=15)

IBM's team demonstrated Shor's on 5 qubits for N=15 (2001). A simplified version:

```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator

def shors_15_a7():
    """Simplified Shor's for N=15, a=7 (period r=4)."""
    # 4 counting qubits + 4 for modular arithmetic
    qc = QuantumCircuit(8, 4)
    
    # Initialize |1⟩ in work register
    qc.x(4)
    
    # Hadamard on counting register
    qc.h([0, 1, 2, 3])
    
    # Simplified modular exponentiation (7^x mod 15)
    # For real circuits, this needs many more gates
    qc.cx(3, 6)
    qc.cx(3, 7)
    qc.h(3)
    qc.swap(5, 7)
    qc.swap(4, 6)
    # ... (full circuit is more complex)
    
    # QFT inverse on counting register
    # ... (apply IQFT)
    
    qc.measure([0, 1, 2, 3], [0, 1, 2, 3])
    return qc
```

---

## Continued Fractions Algorithm

Key classical post-processing:

```python
from fractions import Fraction
import math

def find_period(measurement: int, t: int, N: int) -> int:
    """
    measurement: result from quantum register
    t: number of counting qubits
    N: the number to factor
    """
    if measurement == 0:
        return None  # trivial result, retry
    
    # Phase = measurement / 2^t
    phase = measurement / (2**t)
    
    # Find best rational approximation p/q with q ≤ N
    frac = Fraction(phase).limit_denominator(N)
    r = frac.denominator
    
    # Verify: a^r ≡ 1 mod N
    # (this check is done by the caller)
    return r

# Example: N=15, t=8, a=7, period r=4
# QFT peaks at j·(256/4) = j·64 for j = 0,1,2,3
for j in range(4):
    measurement = j * 64
    r = find_period(measurement, 8, 15)
    if r:
        print(f"measurement={measurement}, phase={measurement/256:.3f}, period={r}")
# measurement=64, phase=0.250, period=4 ✓
# measurement=128, phase=0.500, period=2 (r=2 → need r=lcm(2,...))
# measurement=192, phase=0.750, period=4 ✓
```

---

## Shor's Algorithm Complexity

```
Classical integer factoring (GNFS): O(exp(n^(1/3) · (log n)^(2/3)))
Shor's quantum:                     O(n³) gate operations, O(n) qubits

This is an EXPONENTIAL speedup.

For 2048-bit RSA key (n=2048):
  Classical: O(exp(2048^(1/3))) — astronomical
  Quantum:   O(2048³) = O(10^10) — feasible with fault-tolerant QC

Required quantum resources for 2048-bit RSA:
  Logical qubits: ~4000
  Physical qubits (with error correction): ~20 million
  Gate operations: ~10^12
  
Current NISQ computers: 50-1000 noisy qubits → far from breaking RSA.
Timeline: 2030s-2040s for fault-tolerant systems capable of this.
```

---

## Post-Quantum Cryptography

In response to Shor's threat:
```
NIST PQC Standards (2024):
  CRYSTALS-Kyber → Key encapsulation (replaces RSA key exchange)
  CRYSTALS-Dilithium → Digital signatures
  FALCON → Digital signatures (compact)
  SPHINCS+ → Hash-based signatures (conservative)

These are based on hard problems that quantum computers DON'T speed up significantly:
  Lattice problems (Learning With Errors)
  Hash functions (not significantly affected by Grover)
```
# Chapter 12: Quantum Phase Estimation

## The Phase Estimation Problem

```
PROBLEM: Quantum Phase Estimation (QPE)
  Given:  Unitary U and one of its eigenstates |u⟩
          U|u⟩ = e^(2πiφ)|u⟩   (φ is unknown, 0 ≤ φ < 1)
  Goal:   Estimate φ to t bits of precision

WHY IT MATTERS:
  - Finding energy levels (chemistry, materials science)
  - Shor's algorithm (order finding is phase estimation)
  - Quantum simulation (Hamiltonian eigenvalues)
  - HHL algorithm (linear systems)
  - QMA verification
```

---

## Circuit

```
Ancilla register (t qubits):
|0⟩: ─H─────────────────────────────── IQFT ─ Measure
|0⟩: ─H───────────────────────────────         Measure
|0⟩: ─H─────────────────────────────────────── Measure
...
|0⟩: ─H─●─────────────────────────────         Measure

Eigenstate register:
|u⟩:    ─U^1─U^2─U^4─U^8─ ... ─U^(2^(t-1))──────────────

Each ancilla qubit j controls U^(2^j)

Phase kickback: if U|u⟩ = e^(2πiφ)|u⟩
  controlled-U^(2^j) puts phase e^(2πi·2^j·φ) on ancilla j
  
After all controlled unitaries:
  ancilla in state: (1/√2ᵗ) Σⱼ e^(2πi·j·φ)|j⟩

After IQFT:
  ancilla in state: |φ̃⟩  (φ to t bits precision)
```

---

## Mathematical Derivation

```
Initial: |0⟩^t ⊗ |u⟩

After H^⊗t on ancilla:
  (1/√2ᵗ) Σⱼ₌₀^(2ᵗ-1) |j⟩ ⊗ |u⟩

Apply controlled-U^j (j controlled by j-th ancilla collectively):
  = (1/√2ᵗ) Σⱼ |j⟩ ⊗ U^j|u⟩
  = (1/√2ᵗ) Σⱼ |j⟩ ⊗ e^(2πiφj)|u⟩     (eigenvalue kicks back)
  = (1/√2ᵗ) Σⱼ e^(2πiφj)|j⟩ ⊗ |u⟩     (phase on ancilla)

This is the QFT of |φ̃⟩ applied to the ancilla!
After IQFT:
  |φ̃⟩ ⊗ |u⟩   where φ̃ is the t-bit binary fraction representation of φ
```

---

## Implementation

```python
from qiskit import QuantumCircuit
from qiskit.circuit.library import QFT, PhaseEstimation
from qiskit_aer import AerSimulator
import numpy as np

def qpe_circuit(t: int, unitary: QuantumCircuit) -> QuantumCircuit:
    """
    Quantum Phase Estimation
    t: number of ancilla qubits (bits of precision)
    unitary: circuit implementing U acting on n qubits
    """
    n = unitary.num_qubits
    qc = QuantumCircuit(t + n, t)
    
    # Initialize eigenstate (|1⟩ for simple T gate example)
    # qc.x(t)  # depends on the specific unitary
    
    # Hadamard on ancilla
    qc.h(range(t))
    
    # Controlled-U^(2^j) for each ancilla qubit
    for j in range(t):
        # Apply U^(2^j) controlled by ancilla qubit j
        controlled_u = unitary.repeat(2**j).control(1)
        qc.compose(controlled_u, qubits=[j] + list(range(t, t + n)), inplace=True)
    
    # Inverse QFT on ancilla
    iqft = QFT(t, inverse=True)
    qc.compose(iqft, qubits=range(t), inplace=True)
    
    # Measure ancilla
    qc.measure(range(t), range(t))
    
    return qc

# Example: Estimate phase of T gate
# T|1⟩ = e^(iπ/4)|1⟩, so φ = 1/8
def estimate_T_phase():
    t = 4  # 4 bits precision → can represent 1/8 = 0.001 in binary
    
    # T gate: phase = π/4, so φ = 1/8
    T_gate = QuantumCircuit(1)
    T_gate.t(0)
    
    qc = QuantumCircuit(t + 1, t)
    qc.x(t)    # initialize eigenstate |1⟩ of T gate
    qc.h(range(t))
    
    # Controlled T^(2^j)
    for j in range(t):
        angle = np.pi * (2**j) / 4  # phase of T^(2^j)
        qc.cp(angle, j, t)          # controlled phase gate
    
    # Inverse QFT
    from qiskit.circuit.library import QFT
    qc.compose(QFT(t, inverse=True), qubits=range(t), inplace=True)
    
    qc.measure(range(t), range(t))
    
    return qc

sim = AerSimulator()
qc = estimate_T_phase()
counts = sim.run(qc, shots=1000).result().get_counts()
print("Counts:", counts)
# Expected: '0010' with high probability (binary 0010 = 2/16 = 1/8 ✓)

# Convert measurement to phase
for outcome, count in sorted(counts.items(), key=lambda x: -x[1])[:3]:
    phase = int(outcome, 2) / (2**4)
    print(f"  outcome {outcome} → phase φ = {phase:.4f} (counts: {count})")
```

---

## Precision and Error

```
With t ancilla qubits:
  Precision: 2^(-t)  (can estimate φ to t binary digits)
  
  If φ = 0.001 (binary) = 1/8:
    t=1: can't resolve (0 or 1 → 0 or 1/2)
    t=2: can't resolve (1/4)
    t=3: can't resolve (1/8 exactly at boundary)
    t=4: '0010' → 2/16 = 1/8 ✓

Success probability: ≥ 4/π² ≈ 0.405 for t such that 2ᵗ ≥ 2/ε
  (for error ε in the phase estimate)

To get probability ≥ 1-δ:
  Use t_total = t + ⌈log₂(1 + 1/(2δ))⌉ ancilla qubits
```

---

## QPE for Chemistry: Energy Estimation

QPE is used to estimate eigenvalues of molecular Hamiltonians:

```
Molecule H₂ (hydrogen):
  Hamiltonian H_mol has eigenvalues E₀, E₁, ... (energy levels)
  Ground state energy E₀ = lowest eigenvalue

QPE approach:
  U = e^(-iH_mol·t)   (time evolution operator)
  Eigenvalues of U: e^(-iE_k·t)
  Phase φₖ = E_k·t/(2π)
  
  QPE estimates φₖ → compute E_k = 2π·φₖ/t

Challenge:
  Preparing eigenstate |u⟩ is hard (that's VQE's job)
  Implementing e^(-iH·t) efficiently (quantum simulation)
```

---

## QPE in Shor's Algorithm

```
Order finding: find r such that aʳ ≡ 1 (mod N)

Define U_a: |x⟩ → |ax mod N⟩

Eigenstates of U_a: |u_s⟩ = (1/√r) Σⱼ e^(-2πijs/r) |aʲ mod N⟩

Eigenvalues: e^(2πis/r)  for s = 0, 1, ..., r-1

Phase: φ_s = s/r

QPE on |u_s⟩ (or uniform superposition = Σ_s|u_s⟩):
  Measures s/r for random s
  Continued fractions → extracts r from s/r

This IS quantum period finding!
```

---

## Iterative QPE (Kitaev's Algorithm)

For resource-constrained hardware, estimate one bit of φ at a time:

```python
def iterative_qpe_bit(unitary: QuantumCircuit, eigenstate: QuantumCircuit,
                      k: int, known_bits: int) -> int:
    """
    Estimate k-th bit of φ given known lower bits.
    Uses only 1 ancilla qubit (instead of t).
    """
    # Requires: eigenstate preparation + controlled-U^(2^k)
    qc = QuantumCircuit(1 + unitary.num_qubits, 1)
    
    # Prepare eigenstate
    qc.compose(eigenstate, qubits=range(1, 1 + unitary.num_qubits), inplace=True)
    
    # H on ancilla
    qc.h(0)
    
    # Phase correction for known bits
    omega = 2 * np.pi * known_bits / (2**k)
    qc.p(-omega, 0)
    
    # Controlled-U^(2^k)
    controlled_u = unitary.repeat(2**k).control(1)
    qc.compose(controlled_u, qubits=[0] + list(range(1, 1 + unitary.num_qubits)), inplace=True)
    
    # Hadamard
    qc.h(0)
    qc.measure(0, 0)
    
    # Majority vote over multiple runs
    sim = AerSimulator()
    counts = sim.run(qc, shots=100).result().get_counts()
    return int(max(counts, key=counts.get))
```

---

## Summary

| Aspect | QPE |
|---|---|
| Input | Unitary U, eigenstate |u⟩ |
| Output | t-bit estimate of phase φ |
| Circuit depth | O(t) controlled-U applications |
| Qubits needed | t ancilla + size of |u⟩ |
| Success probability | ≥ 4/π² ≈ 40% (boosted with more ancilla) |
| Key ingredient | Inverse QFT |
| Applications | Energy levels, Shor's, simulation, HHL |
# Chapter 13: Variational Quantum Algorithms (VQE & QAOA)

## The Variational Principle

Hybrid classical-quantum algorithms designed for NISQ hardware:

```
Idea: Use a parameterized quantum circuit (ansatz) with classical optimizer

quantum circuit (parameterized) → measure expectation value → classical optimizer
         ↑___________________________________________|

Theorem (Variational Principle):
  For Hamiltonian H with ground state energy E₀:
  ⟨ψ(θ)|H|ψ(θ)⟩ ≥ E₀   for any state |ψ(θ)⟩

So minimizing ⟨H⟩ over parameters θ gives upper bound on E₀.
```

---

## VQE — Variational Quantum Eigensolver

VQE finds the ground state energy of a Hamiltonian:

```
Application: Quantum chemistry
  Find lowest energy molecular configuration
  Simulate H₂, LiH, BeH₂, etc.

Classical methods:
  Exact diagonalization: O(4ⁿ) — scales exponentially
  Density Functional Theory (DFT): O(n³) — approximate
  CCSD(T): gold standard, O(n⁷) — too slow for large molecules

VQE: polynomial quantum resources (for good ansatz)
```

### VQE Algorithm:
```
1. Encode Hamiltonian H as weighted sum of Pauli operators:
   H = Σᵢ cᵢ Pᵢ    where Pᵢ ∈ {I,X,Y,Z}^⊗n

2. Choose ansatz |ψ(θ)⟩ = U(θ)|0⟩ⁿ

3. Quantum computer: measure ⟨Pᵢ⟩ for each Pauli term
   E(θ) = Σᵢ cᵢ ⟨ψ(θ)|Pᵢ|ψ(θ)⟩

4. Classical computer: update θ to minimize E(θ)
   θ* = argmin E(θ)

5. Return: E(θ*) ≈ E₀ (ground state energy)
```

---

## Implementing VQE

```python
from qiskit import QuantumCircuit
from qiskit.quantum_info import SparsePauliOp, Statevector
from qiskit_aer import AerSimulator
from scipy.optimize import minimize
import numpy as np

# H₂ Hamiltonian (simplified, 2-qubit, at fixed bond length)
H2_hamiltonian = SparsePauliOp.from_list([
    ('II', -1.0523732),
    ('IZ',  0.3979374),
    ('ZI', -0.3979374),
    ('ZZ', -0.0112801),
    ('XX',  0.1809312),
])

def ansatz(theta: list) -> QuantumCircuit:
    """Hardware-efficient ansatz for 2 qubits."""
    qc = QuantumCircuit(2)
    qc.ry(theta[0], 0)
    qc.ry(theta[1], 1)
    qc.cx(0, 1)
    qc.ry(theta[2], 0)
    qc.ry(theta[3], 1)
    return qc

def compute_energy(theta: np.ndarray) -> float:
    """Measure expectation value of H₂ Hamiltonian."""
    qc = ansatz(theta)
    sv = Statevector(qc)
    return sv.expectation_value(H2_hamiltonian).real

# Classical optimization
theta_init = np.random.uniform(0, 2*np.pi, 4)
result = minimize(compute_energy, theta_init, method='COBYLA',
                  options={'maxiter': 1000, 'rhobeg': 0.5})

print(f"VQE ground state energy: {result.fun:.6f} Ha")
print(f"Classical exact:         -1.857275 Ha")
print(f"Error: {abs(result.fun - (-1.857275)):.6f} Ha")

# Visualize optimal state
optimal_theta = result.x
qc = ansatz(optimal_theta)
print(f"Optimal ansatz circuit:")
print(qc.draw())
```

---

## Ansatz Design

The ansatz is critical — it determines what states VQE can explore:

### Hardware-Efficient Ansatz:
```python
def hardware_efficient_ansatz(n: int, depth: int, theta: np.ndarray) -> QuantumCircuit:
    qc = QuantumCircuit(n)
    idx = 0
    for _ in range(depth):
        for q in range(n):
            qc.ry(theta[idx], q); idx += 1
            qc.rz(theta[idx], q); idx += 1
        for q in range(n - 1):
            qc.cx(q, q + 1)  # linear entanglement
    return qc
# Pros: short depth, hardware-friendly
# Cons: may not reach ground state (expressibility vs trainability tradeoff)
```

### Chemistry-Inspired (UCCSD):
```python
# Unitary Coupled Cluster Single Double (UCCSD)
# Based on quantum chemistry theory
from qiskit_nature.second_q.circuit.library import UCCSD
# Better for chemistry but deeper circuit
```

---

## Barren Plateaus Problem

A major challenge for VQA training:

```
For deep circuits with many qubits:
  Gradients of cost function become exponentially small
  ∂E/∂θᵢ → 0 exponentially in n

Cause: random initialization → state is uniformly distributed
       → cost function landscape is flat almost everywhere

Solutions:
  1. Local cost functions (instead of global)
  2. Layer-by-layer training
  3. Problem-specific initialization
  4. Warm starting from classical solution
  5. Rotosolve / parameter shift rules
```

### Parameter shift rule (exact gradient):
```python
def parameter_shift_gradient(theta: np.ndarray, i: int, shift: float = np.pi/2) -> float:
    """Compute exact gradient of E w.r.t. theta[i] using parameter shift."""
    theta_plus = theta.copy(); theta_plus[i] += shift
    theta_minus = theta.copy(); theta_minus[i] -= shift
    return (compute_energy(theta_plus) - compute_energy(theta_minus)) / (2 * np.sin(shift))
```

---

## QAOA — Quantum Approximate Optimization Algorithm

QAOA finds approximate solutions to combinatorial optimization problems.

### Target: MaxCut Problem
```
PROBLEM: MaxCut
  Given: Graph G = (V, E)
  Find: Partition vertices into two sets S, S̄
  Maximize: Number of edges between S and S̄

Classical: NP-hard for general graphs
Best approximation: Goemans-Williamson (0.878 approximation)
QAOA: Parameterizable approximation, performance improves with circuit depth p
```

### QAOA Structure:
```
Encode MaxCut as Ising Hamiltonian:
  H_C = Σ_{(i,j)∈E} (1 - ZᵢZⱼ)/2

QAOA state: |β,γ⟩ = e^(-iβ_p H_B) e^(-iγ_p H_C) ... e^(-iβ₁ H_B) e^(-iγ₁ H_C) |+⟩^n

H_B = Σᵢ Xᵢ  (mixing Hamiltonian — transverse field)
H_C = cost Hamiltonian

Optimize angles β, γ to maximize ⟨H_C⟩
```

```python
from qiskit import QuantumCircuit
from qiskit.quantum_info import SparsePauliOp, Statevector
import numpy as np
from scipy.optimize import minimize

# MaxCut on a simple 4-node cycle graph
# Edges: (0,1), (1,2), (2,3), (3,0)

edges = [(0, 1), (1, 2), (2, 3), (3, 0)]
n = 4

def cost_hamiltonian(n, edges) -> SparsePauliOp:
    """H_C = Σ (I - ZᵢZⱼ)/2 for each edge."""
    terms = []
    for (i, j) in edges:
        # ZᵢZⱼ term
        pauli_str = ['I'] * n
        pauli_str[i] = 'Z'
        pauli_str[j] = 'Z'
        terms.append((''.join(reversed(pauli_str)), -0.5))
        terms.append(('I' * n, 0.5))
    return SparsePauliOp.from_list(terms)

def qaoa_circuit(n: int, edges: list, gamma: np.ndarray, beta: np.ndarray) -> QuantumCircuit:
    p = len(gamma)  # QAOA depth
    qc = QuantumCircuit(n)
    
    # Initialize |+⟩^n
    qc.h(range(n))
    
    for layer in range(p):
        # Phase separation (cost Hamiltonian)
        for (i, j) in edges:
            qc.cx(i, j)
            qc.rz(2 * gamma[layer], j)
            qc.cx(i, j)
        
        # Mixing (X Hamiltonian)
        for q in range(n):
            qc.rx(2 * beta[layer], q)
    
    return qc

def qaoa_cost(params: np.ndarray) -> float:
    p = len(params) // 2
    gamma = params[:p]
    beta = params[p:]
    
    qc = qaoa_circuit(n, edges, gamma, beta)
    sv = Statevector(qc)
    H_C = cost_hamiltonian(n, edges)
    return -sv.expectation_value(H_C).real  # negate for minimization

# Optimize QAOA with p=1
p = 1
params0 = np.random.uniform(0, np.pi, 2 * p)
result = minimize(qaoa_cost, params0, method='COBYLA')

gamma_opt = result.x[:p]
beta_opt = result.x[p:]
qc_opt = qaoa_circuit(n, edges, gamma_opt, beta_opt)
sv = Statevector(qc_opt)

probs = sv.probabilities_dict()
print("Top solutions (MaxCut on 4-cycle):")
for state, prob in sorted(probs.items(), key=lambda x: -x[1])[:4]:
    cut = sum(1 for (i,j) in edges if state[n-1-i] != state[n-1-j])
    print(f"  |{state}⟩: prob={prob:.3f}, cut={cut}")
# Best cut for 4-cycle is 4 (alternate 0101 or 1010)
```

---

## QAOA Performance

```
QAOA depth p:
  p=1: may outperform random, but not optimal for hard instances
  p→∞: approaches optimal solution (adiabatic theorem)

MaxCut approximation ratio:
  p=1 QAOA on 3-regular graphs: ≥ 0.6924
  Best classical (Goemans-Williamson): 0.878
  QAOA with large p: potentially better for specific problems

Current NISQ limitation:
  Decoherence limits useful circuit depth
  p=1-5 is practical on current hardware
```

---

## Other VQA Variants

| Algorithm | Problem | Ansatz |
|---|---|---|
| VQE | Ground state energy | Chemistry-inspired |
| QAOA | Combinatorial optimization | Problem-specific phases |
| VQLS | Linear systems Ax=b | Generic |
| QNN | Machine learning | Parameterized circuit |
| QSVM | Classification | Quantum kernel |
| ADAPT-VQE | Chemistry (adaptive) | Grows greedily |

---

## Classical Simulation of VQAs

```python
# PennyLane: best library for VQA with autodiff
import pennylane as qml
import numpy as np

dev = qml.device('default.qubit', wires=2)

@qml.qnode(dev)
def circuit(theta):
    qml.RY(theta[0], wires=0)
    qml.RY(theta[1], wires=1)
    qml.CNOT(wires=[0, 1])
    qml.RY(theta[2], wires=0)
    return qml.expval(qml.PauliZ(0) @ qml.PauliZ(1))

# Automatic differentiation — exact gradients!
theta = np.array([0.1, 0.2, 0.3], requires_grad=True)
grad_fn = qml.grad(circuit)
print(f"Energy: {circuit(theta):.4f}")
print(f"Gradient: {grad_fn(theta)}")

# Gradient descent
opt = qml.GradientDescentOptimizer(stepsize=0.1)
for i in range(100):
    theta, cost = opt.step_and_cost(circuit, theta)
    if i % 10 == 0:
        print(f"Step {i}: cost = {cost:.4f}")
```
# Chapter 14: Quantum Error Correction

## Why Quantum Errors Are Different

```
Classical error correction:
  Copy a bit: 0 → 000, 1 → 111
  Majority vote: if 001 → correct to 000
  
Quantum: no-cloning theorem! Cannot copy |ψ⟩.
Quantum errors are continuous: small rotations, not just bit flips.
Measuring to check destroys superposition.

SOLUTION: Encode information in entangled states.
Measure only the error (syndrome) without touching the data.
```

---

## Types of Quantum Errors

```
Single-qubit errors (Pauli channel):
  X error:  bit flip  |0⟩↔|1⟩
  Z error:  phase flip |+⟩↔|-⟩
  Y error:  both X and Z
  I:        no error

General error: any unitary close to I = aI + bX + cY + dZ
If we can correct X, Z, Y independently → can correct all single-qubit errors
(digitization of errors)

Multi-qubit errors:
  Depolarizing: random Pauli error on each qubit with probability p
  Erasure: qubit is lost (known location)
  Leakage: qubit leaves the 2D subspace
```

---

## 3-Qubit Bit Flip Code

The simplest quantum error correcting code:

```
Encoding:
  |0⟩ → |000⟩
  |1⟩ → |111⟩
  |ψ⟩ = α|0⟩ + β|1⟩ → α|000⟩ + β|111⟩

Error: X on qubit 1 → α|010⟩ + β|101⟩

Syndrome measurement (check parity WITHOUT measuring data):
  M₁ = Z₀Z₁: compare qubits 0 and 1
  M₂ = Z₁Z₂: compare qubits 1 and 2

  No error: M₁=+1, M₂=+1 → syndrome (0,0)
  Error on q0: M₁=-1, M₂=+1 → syndrome (1,0) → apply X to q0
  Error on q1: M₁=-1, M₂=-1 → syndrome (1,1) → apply X to q1
  Error on q2: M₁=+1, M₂=-1 → syndrome (0,1) → apply X to q2
```

```python
from qiskit import QuantumCircuit, QuantumRegister, ClassicalRegister
from qiskit_aer import AerSimulator
from qiskit_aer.noise import NoiseModel, pauli_error

def bit_flip_encode():
    """Encode 1 logical qubit into 3 physical qubits."""
    qc = QuantumCircuit(3)
    qc.cx(0, 1)  # CNOT qubit 0 → qubit 1
    qc.cx(0, 2)  # CNOT qubit 0 → qubit 2
    return qc

def bit_flip_syndrome():
    """Syndrome measurement for bit flip code."""
    qr = QuantumRegister(3, 'data')
    ar = QuantumRegister(2, 'ancilla')
    cr = ClassicalRegister(2, 'syndrome')
    qc = QuantumCircuit(qr, ar, cr)
    
    # M₁ = Z₀Z₁
    qc.cx(qr[0], ar[0])
    qc.cx(qr[1], ar[0])
    
    # M₂ = Z₁Z₂
    qc.cx(qr[1], ar[1])
    qc.cx(qr[2], ar[1])
    
    qc.measure(ar, cr)
    return qc

def bit_flip_correct():
    """Apply correction based on syndrome."""
    qc = QuantumCircuit(3, 2)
    # Syndrome 01 → flip q0
    # Syndrome 11 → flip q1
    # Syndrome 10 → flip q2
    # (implemented as classical control in real circuit)
    return qc
```

**Limitation:** Only corrects bit flip (X) errors, not phase flip (Z) errors.

---

## 3-Qubit Phase Flip Code

```
Encoding (in X basis):
  |+⟩ → |+++⟩
  |-⟩ → |---⟩

Equivalent to: apply H to all qubits of bit flip code!

|0⟩_L = |+⟩|+⟩|+⟩ = H⊗3|000⟩
|1⟩_L = |-⟩|-⟩|-⟩ = H⊗3|111⟩

Syndrome measurement:
  M₁ = X₀X₁:  compare phase of qubits 0 and 1
  M₂ = X₁X₂:  compare phase of qubits 1 and 2
```

---

## Shor's 9-Qubit Code

Combines bit flip and phase flip codes — the first code to correct ALL single-qubit errors:

```
Logical 0: |0⟩_L = (|000⟩+|111⟩)^⊗3 / 2√2
Logical 1: |1⟩_L = (|000⟩-|111⟩)^⊗3 / 2√2

Structure:
  Outer code (3 qubits): phase flip code
  Inner code (each of 3 → 3 qubits): bit flip code

Encodes 1 logical qubit into 9 physical qubits
Corrects any single-qubit X, Z, Y error
```

```python
def shor_9_encode(state: QuantumCircuit) -> QuantumCircuit:
    """Encode 1 logical qubit into 9 physical qubits using Shor's code."""
    qc = QuantumCircuit(9)
    
    # First, encode in phase flip code:
    # |0⟩ → (|000⟩+|111⟩)/√2 structure across blocks
    qc.cx(0, 3)  # q0 → q3 (outer code)
    qc.cx(0, 6)  # q0 → q6
    
    # Apply H to start each block
    qc.h([0, 3, 6])
    
    # Within each block, apply bit flip code
    for start in [0, 3, 6]:
        qc.cx(start, start + 1)
        qc.cx(start, start + 2)
    
    return qc
```

---

## Stabilizer Formalism

A powerful framework for quantum error correcting codes:

```
Stabilizer group S: set of commuting Pauli operators {Sᵢ}
Codespace: all states |ψ⟩ such that Sᵢ|ψ⟩ = |ψ⟩ for all i

3-qubit bit flip code stabilizers:
  S₁ = Z₀Z₁    (qubits 0,1 have same value)
  S₂ = Z₁Z₂    (qubits 1,2 have same value)

Syndrome: which stabilizers anticommute with the error?
  X on q0: {X₀,Z₀Z₁} = 0 → S₁ anticommutes → measure -1 for S₁
           {X₀,Z₁Z₂} = I → S₂ commutes → measure +1 for S₂
  Syndrome: (S₁=-1, S₂=+1) → error on q0
```

---

## Surface Codes (Topological QEC)

The most promising code for practical quantum computers:

```
Arrangement: Qubits on 2D grid
  Data qubits: on vertices (or edges)
  Ancilla qubits: on faces and edges (for syndrome measurement)

Logical qubit: encoded in topology of the lattice
Syndrome: local parity measurements (weight-4 Pauli operators)

Key advantage:
  Only requires nearest-neighbor gates
  High threshold: ~1% physical error rate (vs ~0.01% for other codes)
  Scalable: increase grid size → increase protection
```

```
Surface code layout (d=3):
  
  q─s─q─s─q
  │       │
  s   q   s
  │       │
  q─s─q─s─q
  
  q = data qubit
  s = syndrome (ancilla) qubit
  
  Distance d=3: can correct t = ⌊(d-1)/2⌋ = 1 error
  Physical qubits: d² + (d-1)² = 13 for d=3

Resource overhead:
  1 logical qubit = d² physical qubits
  d=3:  9 physical → 1 logical
  d=7: 49 physical → 1 logical
  d=23: 529 physical → 1 logical (needed for Shor's on 2048-bit RSA)
  
  For Shor's 2048-bit RSA: ~4000 logical × ~529 physical = ~2M physical qubits
```

---

## Fault-Tolerant Quantum Computing

Error correction alone isn't enough — the correction operations must also be fault-tolerant:

```
Fault-tolerant gate set for surface code: Clifford + T gate
  Clifford gates: H, S, CNOT — implemented transversally (bit-by-bit)
  T gate: expensive — requires magic state distillation

Magic state distillation:
  Prepare many noisy T states, distill down to fewer high-fidelity ones
  Resources: ~O(15) noisy states per distilled state (15-to-1 protocol)

Full fault-tolerant computation:
  Each logical gate → many physical operations
  Overhead factor: thousands to millions vs bare quantum computation
```

---

## Quantum Error Threshold

```
THRESHOLD THEOREM:
  If physical error rate p < p_threshold:
    Can perform arbitrarily long computation with overhead polylog(n)
    where n = computation depth

  p_threshold for surface code: ≈ 1%
  Current superconducting qubit errors: ≈ 0.1-1%

  We're close to the threshold! This is why surface codes are the main focus.

Timeline:
  NISQ era (now): 50-1000 noisy qubits, no error correction
  FTQC era (2030s?): millions of physical → thousands of logical qubits
```

---

## Simulating Quantum Error Correction

```python
from qiskit_aer.noise import NoiseModel, pauli_error, depolarizing_error
from qiskit_aer import AerSimulator
from qiskit import QuantumCircuit

# Build noise model
noise_model = NoiseModel()
p = 0.01  # 1% error rate

# Add depolarizing error after each gate
error_1q = depolarizing_error(p, 1)
error_2q = depolarizing_error(p, 2)
noise_model.add_all_qubit_quantum_error(error_1q, ['h', 'rx', 'ry', 'rz'])
noise_model.add_all_qubit_quantum_error(error_2q, ['cx', 'cz'])

# Simulate with noise
noisy_sim = AerSimulator(noise_model=noise_model)

# Test: prepare |+⟩, measure many times
qc = QuantumCircuit(1, 1)
qc.h(0)
qc.measure(0, 0)

ideal_result = AerSimulator().run(qc, shots=1000).result().get_counts()
noisy_result = noisy_sim.run(qc, shots=1000).result().get_counts()
print(f"Ideal:  {ideal_result}")
print(f"Noisy:  {noisy_result}")
```

---

## Quantum Error Mitigation (Near-Term Alternative)

Without full QEC, use classical post-processing to reduce errors:

```python
# Zero-noise extrapolation (ZNE)
def zne_result(circuit, noise_factors=[1, 2, 3]):
    """Estimate noiseless result by extrapolating to zero noise."""
    results = []
    for factor in noise_factors:
        # Scale noise by inserting identity pairs
        noisy_circuit = scale_noise(circuit, factor)
        result = run_circuit(noisy_circuit)
        results.append(result)
    
    # Fit polynomial and extrapolate to factor=0
    from numpy.polynomial import polynomial as P
    coeffs = P.polyfit(noise_factors, results, deg=2)
    return P.polyval(0, coeffs)  # extrapolate to zero noise
```
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
# Chapter 16: Quantum Programming with Qiskit

## Qiskit Ecosystem

Qiskit (IBM's quantum SDK) is the most widely used quantum programming framework:

```
qiskit              → Core: circuits, gates, transpilation
qiskit-aer          → Simulators (statevector, QASM, noise)
qiskit-ibm-runtime  → Access to IBM quantum hardware
qiskit-nature       → Quantum chemistry
qiskit-ml           → Quantum machine learning
qiskit-optimization → Combinatorial optimization
```

```bash
pip install qiskit qiskit-aer qiskit-ibm-runtime
pip install matplotlib pylatexenc  # visualization
```

---

## Building Circuits

```python
from qiskit import QuantumCircuit, QuantumRegister, ClassicalRegister

# Method 1: Simple
qc = QuantumCircuit(3, 2)  # 3 qubits, 2 classical bits

# Method 2: Named registers
qr = QuantumRegister(2, 'data')
ar = QuantumRegister(1, 'ancilla')
cr = ClassicalRegister(2, 'output')
qc = QuantumCircuit(qr, ar, cr)

# All gates
qc.h(qr[0])
qc.x(qr[1])
qc.y(qr[0])
qc.z(qr[1])
qc.s(qr[0])       # S gate = Z^(1/2)
qc.sdg(qr[0])     # S† gate
qc.t(qr[0])       # T gate = Z^(1/4)
qc.tdg(qr[0])     # T† gate
qc.rx(0.5, qr[0])
qc.ry(1.2, qr[1])
qc.rz(0.3, qr[0])
qc.p(0.7, qr[0])  # phase gate
qc.u(0.5, 1.2, 0.3, qr[0])  # general U gate

# Multi-qubit gates
qc.cx(qr[0], qr[1])      # CNOT
qc.cz(qr[0], qr[1])      # CZ
qc.cy(qr[0], qr[1])      # CY
qc.swap(qr[0], qr[1])    # SWAP
qc.ccx(qr[0], qr[1], ar[0])  # Toffoli
qc.cswap(ar[0], qr[0], qr[1])  # Fredkin

# Controlled versions
qc.ch(qr[0], qr[1])      # Controlled H
qc.crx(0.5, qr[0], qr[1])

# Measurement
qc.measure(qr, cr)         # all at once
qc.measure(qr[0], cr[0])  # individual

# Reset
qc.reset(qr[0])   # reset to |0⟩ (mid-circuit)

# Barrier (visual separator, no effect on computation)
qc.barrier()

# Draw
print(qc.draw())           # text
qc.draw('mpl')             # matplotlib figure
```

---

## Simulating Circuits

```python
from qiskit_aer import AerSimulator, StatevectorSimulator
from qiskit.quantum_info import Statevector

# 1. Statevector simulation (exact, exponential memory)
qc = QuantumCircuit(3)
qc.h(0)
qc.cx(0, 1)
qc.cx(0, 2)  # GHZ state

sv = Statevector(qc)
print(sv)                           # full state vector
print(sv.probabilities_dict())      # {'000': 0.5, '111': 0.5}
sv.draw('bloch')                    # Bloch sphere (single qubits only)
sv.draw('qsphere')                  # Q-sphere visualization

# 2. Sampling simulation (QASM)
qc_m = qc.copy()
qc_m.measure_all()
sim = AerSimulator()
job = sim.run(qc_m, shots=8192)
counts = job.result().get_counts()
print(counts)  # {'000': 4103, '111': 4089}

# 3. Density matrix simulation
from qiskit_aer import DensityMatrixSimulator
from qiskit.quantum_info import DensityMatrix
dm = DensityMatrix(qc)
print(dm.purity())  # 1.0 = pure state

# 4. Noisy simulation
from qiskit_aer.noise import NoiseModel, depolarizing_error
noise_model = NoiseModel()
noise_model.add_all_qubit_quantum_error(depolarizing_error(0.01, 1), ['h', 'rx'])
noise_model.add_all_qubit_quantum_error(depolarizing_error(0.05, 2), ['cx'])
noisy_sim = AerSimulator(noise_model=noise_model)
noisy_counts = noisy_sim.run(qc_m, shots=8192).result().get_counts()
print(noisy_counts)  # Imperfect, may see '001', '010', etc.
```

---

## Transpilation

Before running on hardware, circuits are transpiled to the target's native gate set and connectivity:

```python
from qiskit.compiler import transpile
from qiskit_aer import AerSimulator

# Simulate a backend with limited connectivity
from qiskit.providers.fake_provider import FakeManilaV2
backend = FakeManilaV2()  # simulated 5-qubit IBM backend

qc = QuantumCircuit(3)
qc.h(0)
qc.cx(0, 2)  # non-adjacent qubits
qc.cx(1, 2)

# Transpile for backend
t_qc = transpile(qc, backend=backend, optimization_level=3)
print(f"Original: {qc.count_ops()}")
print(f"Transpiled: {t_qc.count_ops()}")
print(f"Added SWAP gates: {t_qc.count_ops().get('swap', 0)}")
```

---

## Running on IBM Quantum Hardware

```python
from qiskit_ibm_runtime import QiskitRuntimeService, SamplerV2 as Sampler

# Setup (one-time): save credentials
QiskitRuntimeService.save_account(
    channel="ibm_quantum",
    token="YOUR_IBM_API_TOKEN",  # get from quantum.ibm.com
    set_as_default=True,
    overwrite=True
)

# Load account
service = QiskitRuntimeService()

# List available backends
backends = service.backends()
for b in backends:
    print(f"{b.name}: {b.num_qubits} qubits, pending: {b.status().pending_jobs}")

# Select a backend
backend = service.least_busy(min_num_qubits=5)

# Build circuit
qc = QuantumCircuit(2, 2)
qc.h(0)
qc.cx(0, 1)
qc.measure_all()

# Transpile
from qiskit.compiler import transpile
t_qc = transpile(qc, backend=backend, optimization_level=3)

# Run with Sampler (returns probability distribution)
sampler = Sampler(backend)
job = sampler.run([t_qc], shots=1024)
result = job.result()
counts = result[0].data.meas.get_counts()
print(f"Results from {backend.name}:", counts)
```

---

## Quantum Information Tools

```python
from qiskit.quantum_info import (
    Statevector, DensityMatrix, Operator,
    state_fidelity, average_gate_fidelity,
    entropy, partial_trace, random_statevector,
    SparsePauliOp, Pauli
)

# State fidelity
psi = Statevector.from_label('+')   # |+⟩
phi = Statevector([1/2**0.5, 0])    # |0⟩ + 0|1⟩ — wait, this is |0⟩/√... 
phi = Statevector([1, 0])           # |0⟩
print(state_fidelity(psi, phi))     # 0.5

# Operator fidelity
U = Operator.from_label('H')   # Hadamard
V = Operator.from_label('X')   # X gate
print(average_gate_fidelity(U, V))

# Entanglement entropy
bell = Statevector([1/2**0.5, 0, 0, 1/2**0.5])
rho = DensityMatrix(bell)
rho_A = partial_trace(rho, [1])
print(f"Entanglement entropy: {entropy(rho_A, base=2):.4f}")  # 1.0 ebit

# Pauli operator expectation
qc = QuantumCircuit(2)
qc.h(0); qc.cx(0, 1)
sv = Statevector(qc)
ZZ = SparsePauliOp('ZZ')
print(f"⟨ZZ⟩ = {sv.expectation_value(ZZ).real:.4f}")  # 1.0 for Bell state
```

---

## Quantum Circuit Library

Qiskit includes many pre-built circuits:

```python
from qiskit.circuit.library import (
    QFT,                    # Quantum Fourier Transform
    PhaseEstimation,        # QPE circuit
    GroverOperator,         # Grover's diffuser
    TwoLocal,               # Hardware-efficient ansatz
    EfficientSU2,           # Another VQA ansatz
    PauliEvolutionGate,     # exp(-iHt) for Pauli H
    ZZFeatureMap,           # Quantum kernel for ML
    RealAmplitudes,         # Parameterized circuit
)
import numpy as np

# QFT
qft = QFT(4)
print(qft.draw())

# Grover operator
from qiskit.circuit.library import GroverOperator
oracle = QuantumCircuit(3)
oracle.ccz(0, 1, 2)  # marks |111⟩
grover_op = GroverOperator(oracle)
print(grover_op.draw())

# Ansatz for VQE
ansatz = EfficientSU2(num_qubits=4, reps=2, entanglement='linear')
print(f"Parameters: {len(ansatz.parameters)}")
print(ansatz.draw())

# Bind parameters
params = {p: np.random.uniform(0, 2*np.pi) for p in ansatz.parameters}
bound = ansatz.assign_parameters(params)
```

---

## Advanced: Custom Gate Definitions

```python
from qiskit.circuit import Gate
from qiskit import QuantumCircuit
import numpy as np

# Define a custom gate
class iSWAP(Gate):
    def __init__(self):
        super().__init__('iswap', 2, [])
    
    def __array__(self, dtype=None):
        return np.array([[1, 0, 0, 0],
                         [0, 0, 1j, 0],
                         [0, 1j, 0, 0],
                         [0, 0, 0, 1]], dtype=dtype)

# Use in circuit
qc = QuantumCircuit(2)
qc.append(iSWAP(), [0, 1])
print(qc.draw())

# Decompose to native gates for simulation
from qiskit.quantum_info import Operator
print(Operator(qc))
```

---

## Primitives API (Modern Qiskit)

```python
from qiskit_ibm_runtime import SamplerV2, EstimatorV2

# Sampler: get measurement statistics
sampler = SamplerV2(backend)
pub = (circuit, )  # PUB = Primitive Unified Bloc
job = sampler.run([pub])
result = job.result()
counts = result[0].data.meas.get_counts()

# Estimator: compute expectation values
from qiskit.quantum_info import SparsePauliOp
estimator = EstimatorV2(backend)
observable = SparsePauliOp('ZZ')
pub = (circuit, observable)
job = estimator.run([pub])
exp_val = job.result()[0].data.evs
print(f"⟨ZZ⟩ = {exp_val}")
```

---

## Complete Example: Variational Circuit on Real Hardware

```python
from qiskit import QuantumCircuit
from qiskit.quantum_info import SparsePauliOp
from qiskit_ibm_runtime import QiskitRuntimeService, EstimatorV2
from qiskit.compiler import transpile
from scipy.optimize import minimize
import numpy as np

# Hamiltonian: Z⊗Z + X⊗I
H = SparsePauliOp.from_list([('ZZ', 1.0), ('XI', 0.5)])

def parameterized_circuit(theta):
    qc = QuantumCircuit(2)
    qc.ry(theta[0], 0)
    qc.ry(theta[1], 1)
    qc.cx(0, 1)
    qc.ry(theta[2], 0)
    return qc

# Load backend
service = QiskitRuntimeService()
backend = service.least_busy(min_num_qubits=2)
estimator = EstimatorV2(backend)

# Cost function (calls quantum hardware)
def cost(theta):
    qc = parameterized_circuit(theta)
    t_qc = transpile(qc, backend=backend)
    pub = (t_qc, H)
    result = estimator.run([pub]).result()
    return result[0].data.evs

# Optimize
theta0 = np.random.uniform(0, 2*np.pi, 3)
result = minimize(cost, theta0, method='COBYLA', options={'maxiter': 100})
print(f"Ground state energy: {result.fun:.4f}")
```

---

## Qiskit Cheat Sheet

```python
# Create         qc = QuantumCircuit(n, m)
# Gates          qc.h/x/y/z/s/t/cx/cz/swap/ccx(qubit)
# Rotation       qc.rx/ry/rz(angle, qubit)
# Measure        qc.measure(qubits, clbits)
# Draw           qc.draw() or qc.draw('mpl')
# Depth          qc.depth()
# Gate count     qc.count_ops()
# Statevector    Statevector(qc)
# Density matrix DensityMatrix(qc)
# Simulate       AerSimulator().run(qc, shots=1024).result().get_counts()
# Transpile      transpile(qc, backend, optimization_level=3)
# Operator       Operator(qc)  or SparsePauliOp('ZZ')
# Expectation    sv.expectation_value(operator)
# Fidelity       state_fidelity(sv1, sv2)
# Compose        qc1.compose(qc2, inplace=True)
# Inverse        qc.inverse()
# Repeat         qc.repeat(k)
# Control        qc.control(num_ctrl_qubits)
```
# Chapter 17: Cirq and PennyLane

## Cirq (Google)

Cirq is Google's quantum computing framework, designed for NISQ devices:

```bash
pip install cirq cirq-google
```

### Core concepts:
```python
import cirq
import numpy as np

# Qubits
q0, q1, q2 = cirq.LineQubit.range(3)       # linear array
a = cirq.GridQubit(0, 0); b = cirq.GridQubit(0, 1)  # grid qubits (Google hardware)
q = cirq.NamedQubit('alice')                # named qubit

# Gates and moments
circuit = cirq.Circuit([
    cirq.H(q0),
    cirq.CNOT(q0, q1),
    cirq.measure(q0, q1, key='result')
])
print(circuit)

# Parametrized gates
theta = cirq.Parameter('theta')
qc = cirq.Circuit(cirq.Ry(rads=theta)(q0))
```

### Bell state in Cirq:
```python
import cirq

q0, q1 = cirq.LineQubit.range(2)
circuit = cirq.Circuit([
    cirq.H(q0),
    cirq.CNOT(q0, q1),
    cirq.measure(q0, q1, key='bell'),
])
print(circuit)

# Simulate
simulator = cirq.Simulator()
result = simulator.run(circuit, repetitions=1000)
print(result.histogram(key='bell'))
# {0: ~500 (=|00⟩), 3: ~500 (=|11⟩)}
```

### Statevector simulation:
```python
simulator = cirq.Simulator()
result = simulator.simulate(circuit)
print(result.final_state_vector)  # complex array

# Density matrix simulation
dm_sim = cirq.DensityMatrixSimulator()
dm_result = dm_sim.simulate(circuit)
print(dm_result.final_density_matrix)
```

### Cirq noise models:
```python
# Depolarizing noise
noise = cirq.depolarize(p=0.01)
noisy_circuit = circuit.with_noise(noise)

result = cirq.DensityMatrixSimulator(noise=noise).run(circuit, repetitions=1000)

# Channel as gate
noisy_q = cirq.Circuit(
    cirq.H(q0),
    cirq.depolarize(0.01)(q0),  # apply noise channel
    cirq.measure(q0, key='m')
)
```

### Grover's algorithm in Cirq:
```python
import cirq
import numpy as np

def oracle_circuit(qubits, target_state: list):
    """Phase oracle that marks target_state with a phase flip."""
    n = len(qubits)
    circuit = cirq.Circuit()
    
    # Flip qubits where target is 0
    flips = [qubits[i] for i, b in enumerate(target_state) if b == 0]
    if flips:
        circuit.append(cirq.X.on_each(*flips))
    
    # Multi-controlled Z (phase flip for |11...1⟩)
    circuit.append(cirq.Z.controlled(n-1).on(*qubits))
    
    # Undo flips
    if flips:
        circuit.append(cirq.X.on_each(*flips))
    
    return circuit

def diffuser(qubits):
    """Grover diffusion operator."""
    n = len(qubits)
    circuit = cirq.Circuit()
    circuit.append(cirq.H.on_each(*qubits))
    circuit.append(cirq.X.on_each(*qubits))
    circuit.append(cirq.Z.controlled(n-1).on(*qubits))
    circuit.append(cirq.X.on_each(*qubits))
    circuit.append(cirq.H.on_each(*qubits))
    return circuit

def grover(n, target_state, num_iterations=None):
    qubits = cirq.LineQubit.range(n)
    if num_iterations is None:
        num_iterations = int(np.pi / 4 * np.sqrt(2**n))
    
    circuit = cirq.Circuit()
    circuit.append(cirq.H.on_each(*qubits))  # superposition
    
    for _ in range(num_iterations):
        circuit += oracle_circuit(qubits, target_state)
        circuit += diffuser(qubits)
    
    circuit.append(cirq.measure(*qubits, key='result'))
    return circuit

# Search for |101⟩
qc = grover(3, [1, 0, 1])
result = cirq.Simulator().run(qc, repetitions=1000)
print(result.histogram(key='result'))
```

---

## PennyLane

PennyLane is the best framework for quantum machine learning — supports automatic differentiation and multiple backends:

```bash
pip install pennylane pennylane-qiskit pennylane-cirq torch
```

### Core concepts:
```python
import pennylane as qml
import numpy as np

# Create a device (simulator or hardware)
dev = qml.device('default.qubit', wires=3)      # exact statevector
dev_noise = qml.device('qiskit.aer', wires=3)   # use Qiskit's Aer

# QNode: quantum function decorator
@qml.qnode(dev)
def circuit(params):
    qml.RY(params[0], wires=0)
    qml.RY(params[1], wires=1)
    qml.CNOT(wires=[0, 1])
    return qml.expval(qml.PauliZ(0))  # return expectation value

# Call like a function
params = np.array([0.5, 1.2])
print(circuit(params))  # float
```

### Automatic differentiation:
```python
import pennylane as qml
import numpy as np

dev = qml.device('default.qubit', wires=2)

@qml.qnode(dev)
def circuit(theta):
    qml.RY(theta, wires=0)
    qml.CNOT(wires=[0, 1])
    return qml.expval(qml.PauliZ(0))

# Gradient via parameter shift rule (exact)
theta = np.array(0.5, requires_grad=True)
grad = qml.grad(circuit)(theta)
print(f"dE/dtheta = {grad}")

# Interface with PyTorch
import torch
dev_torch = qml.device('default.qubit', wires=2)

@qml.qnode(dev_torch, interface='torch')
def torch_circuit(theta):
    qml.RY(theta, wires=0)
    return qml.expval(qml.PauliZ(0))

theta_torch = torch.tensor(0.5, requires_grad=True)
loss = torch_circuit(theta_torch)
loss.backward()
print(f"Gradient: {theta_torch.grad}")
```

### VQE with PennyLane:
```python
import pennylane as qml
from pennylane import numpy as np

dev = qml.device('default.qubit', wires=2)

# H₂ Hamiltonian
coeffs = [-1.0523732, 0.3979374, -0.3979374, -0.0112801, 0.1809312]
obs = [
    qml.Identity(0) @ qml.Identity(1),
    qml.PauliZ(0) @ qml.Identity(1),
    qml.Identity(0) @ qml.PauliZ(1),
    qml.PauliZ(0) @ qml.PauliZ(1),
    qml.PauliX(0) @ qml.PauliX(1),
]
H = qml.Hamiltonian(coeffs, obs)

@qml.qnode(dev)
def vqe_circuit(theta):
    qml.RY(theta[0], wires=0)
    qml.RY(theta[1], wires=1)
    qml.CNOT(wires=[0, 1])
    qml.RY(theta[2], wires=0)
    qml.RY(theta[3], wires=1)
    return qml.expval(H)

# Optimize
theta = np.random.uniform(0, 2*np.pi, 4, requires_grad=True)
opt = qml.AdamOptimizer(stepsize=0.1)
for i in range(200):
    theta, cost = opt.step_and_cost(vqe_circuit, theta)
    if i % 20 == 0:
        print(f"Step {i}: E = {cost:.6f} Ha")

print(f"\nFinal energy: {vqe_circuit(theta):.6f} Ha")
print(f"Exact (FCI):  -1.857275 Ha")
```

### QAOA with PennyLane:
```python
import pennylane as qml
from pennylane import numpy as np
import networkx as nx

# Create graph
G = nx.cycle_graph(4)  # 4-cycle
n_wires = G.number_of_nodes()

dev = qml.device('default.qubit', wires=n_wires)

def cost_layer(gamma):
    for edge in G.edges():
        qml.CNOT(wires=list(edge))
        qml.RZ(2 * gamma, wires=edge[1])
        qml.CNOT(wires=list(edge))

def mixer_layer(beta):
    for q in range(n_wires):
        qml.RX(2 * beta, wires=q)

@qml.qnode(dev)
def qaoa_circuit(params, depth=1):
    gammas = params[:depth]
    betas = params[depth:]
    
    # Initialize
    for q in range(n_wires):
        qml.Hadamard(wires=q)
    
    # QAOA layers
    for d in range(depth):
        cost_layer(gammas[d])
        mixer_layer(betas[d])
    
    # Cost function: negative MaxCut
    cost_obs = []
    for edge in G.edges():
        cost_obs.append(qml.PauliZ(edge[0]) @ qml.PauliZ(edge[1]))
    
    return qml.expval(sum(cost_obs))

# Optimize with p=1
depth = 1
params = np.random.uniform(0, np.pi, 2*depth, requires_grad=True)
opt = qml.GradientDescentOptimizer(stepsize=0.5)

for i in range(100):
    params, cost = opt.step_and_cost(lambda p: qaoa_circuit(p, depth), params)

print(f"MaxCut approximation: {-qaoa_circuit(params, depth):.3f}")
print(f"(True MaxCut for 4-cycle = 4)")
```

---

## PennyLane Quantum ML Templates

```python
import pennylane as qml
import numpy as np

dev = qml.device('default.qubit', wires=4)

# Angle embedding: encode classical data as rotation angles
@qml.qnode(dev)
def angle_embedding_circuit(x, weights):
    qml.AngleEmbedding(x, wires=range(4), rotation='Y')
    qml.BasicEntanglerLayers(weights, wires=range(4))
    return [qml.expval(qml.PauliZ(i)) for i in range(4)]

x = np.array([0.1, 0.5, 0.3, 0.9])  # input data
weight_shape = qml.BasicEntanglerLayers.shape(n_layers=2, n_wires=4)
weights = np.random.uniform(0, 2*np.pi, weight_shape)
print(angle_embedding_circuit(x, weights))

# Amplitude embedding: encode 2^n data points in n-qubit state
@qml.qnode(dev)
def amplitude_embedding_circuit(x):
    qml.AmplitudeEmbedding(x, wires=range(4), normalize=True)
    return qml.expval(qml.PauliZ(0))

x_data = np.array([0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8,
                   0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8])  # 2^4 = 16 values
print(amplitude_embedding_circuit(x_data))
```

---

## Comparing Frameworks

| Feature | Qiskit | Cirq | PennyLane |
|---|---|---|---|
| Hardware access | IBM Q | Google | Multiple backends |
| Best for | IBM hardware, general QC | Google hardware, NISQ | QML, autodiff |
| Circuit style | Object-oriented | Gate-on-qubit | Functional/decorator |
| Autodiff | No (manual) | No | Yes (parameter shift, adjoint) |
| ML integration | Limited | Limited | PyTorch, TF, JAX, NumPy |
| Noise models | Excellent (Aer) | Good | Via Qiskit/Cirq backend |
| Community | Largest | Large | Growing (ML focus) |
| Visualization | Good | Good | Good |
| Learning curve | Medium | Medium | Lower for ML devs |

### When to use which:
```
Qiskit: IBM quantum hardware, comprehensive QC library, production
Cirq:   Google Sycamore hardware, NISQ algorithm research
PennyLane: Quantum ML, hybrid optimization, autodiff gradients

Can mix: PennyLane can use Qiskit/Cirq as backend!
qml.device('qiskit.ibm', ...)  # PennyLane → IBM hardware via Qiskit
```

---

## Problem-Solving: Choosing a Framework

```python
# Decision guide
def choose_framework(task: str) -> str:
    if "IBM hardware" in task:
        return "Qiskit"
    elif "Google hardware" in task:
        return "Cirq"
    elif "machine learning" in task or "gradient" in task or "VQE" in task:
        return "PennyLane"
    elif "research" in task or "noise" in task:
        return "Qiskit (Aer noise models)"
    elif "production" in task:
        return "Qiskit (largest ecosystem)"
    else:
        return "PennyLane (most versatile)"
```
# Chapter 18: Quantum Machine Learning

## Overview of QML

Quantum Machine Learning (QML) applies quantum computing to machine learning:

```
Three paradigms:
1. CC → Classical data, Classical ML (traditional)
2. QC → Quantum data, Classical ML (quantum sensor → classical analysis)
3. CQ → Classical data, Quantum ML (most common QML research)
4. QQ → Quantum data, Quantum ML (quantum network learning)

Current focus: CQ — use quantum circuits to process classical data

Central question: Does quantum ML offer genuine speedup?
Honest answer: Unclear. Active research. Many claims are overstated.
```

---

## Variational Quantum Classifier

The quantum analog of a classical neural network:

```
Input: classical data x ∈ ℝⁿ
Encoding: U(x)|0⟩ → quantum state (feature map)
Model: V(θ)|ψ(x)⟩ → parameterized circuit (ansatz)
Output: ⟨O⟩ = ⟨ψ(x,θ)|O|ψ(x,θ)⟩ → real value → binary/multiclass label
Training: minimize loss L(θ) w.r.t. θ using gradient descent
```

```python
import pennylane as qml
import numpy as np
from sklearn.datasets import make_classification
from sklearn.preprocessing import normalize

# Create quantum device
dev = qml.device('default.qubit', wires=4)
n_qubits = 4

# Feature map: encode data as angles
def feature_map(x):
    qml.AngleEmbedding(x, wires=range(n_qubits))
    for i in range(n_qubits - 1):
        qml.CNOT(wires=[i, i+1])

# Variational ansatz
def ansatz(weights):
    qml.BasicEntanglerLayers(weights, wires=range(n_qubits))

@qml.qnode(dev)
def quantum_classifier(x, weights):
    feature_map(x)
    ansatz(weights)
    return qml.expval(qml.PauliZ(0))  # binary: >0 → class 1, <0 → class 0

# Training
def loss(weights, X, y):
    preds = np.array([quantum_classifier(x, weights) for x in X])
    return np.mean((preds - y)**2)  # MSE loss

# Dataset
X, y = make_classification(n_samples=100, n_features=4, random_state=42)
X = normalize(X) * np.pi  # scale to [-π, π]
y = 2 * y - 1  # convert {0,1} to {-1,+1}

# Initialize weights
weight_shape = qml.BasicEntanglerLayers.shape(n_layers=2, n_wires=n_qubits)
weights = np.random.uniform(-np.pi, np.pi, weight_shape, requires_grad=True)

# Optimize
opt = qml.AdamOptimizer(stepsize=0.1)
for i in range(50):
    weights, l = opt.step_and_cost(lambda w: loss(w, X[:20], y[:20]), weights)
    if i % 10 == 0:
        preds = np.sign([quantum_classifier(x, weights) for x in X])
        acc = np.mean(preds == y)
        print(f"Step {i}: loss={l:.4f}, acc={acc:.3f}")
```

---

## Quantum Kernel Methods

A mathematically rigorous approach to QML:

```
Classical kernel: k(x, x') = φ(x)·φ(x')  (implicit feature map)
Quantum kernel: k_Q(x, x') = |⟨φ_Q(x)|φ_Q(x')⟩|²

where |φ_Q(x)⟩ = U(x)|0⟩ is a quantum feature state

Advantage: Quantum feature spaces may be exponentially larger than classical.
           Hard to compute classically → potential quantum advantage.

Use with classical SVM:
  Compute kernel matrix K[i,j] = k_Q(xᵢ, xⱼ) on quantum computer
  Train SVM with this kernel on classical computer
```

```python
import pennylane as qml
import numpy as np
from sklearn.svm import SVC
from sklearn.metrics import accuracy_score

dev = qml.device('default.qubit', wires=4)

@qml.qnode(dev)
def kernel_circuit(x1, x2):
    """Compute |⟨φ(x1)|φ(x2)⟩|² quantum kernel."""
    # Encode x1
    qml.AngleEmbedding(x1, wires=range(4))
    # Encode x2 inverse
    qml.AngleEmbedding(x2, wires=range(4)).inv()
    return qml.probs(wires=range(4))

def quantum_kernel(X1, X2):
    """Compute quantum kernel matrix."""
    K = np.zeros((len(X1), len(X2)))
    for i, x1 in enumerate(X1):
        for j, x2 in enumerate(X2):
            # |⟨0|U†(x2)U(x1)|0⟩|² = probability of measuring |0...0⟩
            probs = kernel_circuit(x1, x2)
            K[i, j] = probs[0]  # |0...0⟩ probability
    return K

# Generate data
from sklearn.datasets import make_classification
X, y = make_classification(n_samples=60, n_features=4, random_state=42)
X_train, X_test = X[:40], X[40:]
y_train, y_test = y[:40], y[40:]

# Compute kernel matrix
print("Computing quantum kernel matrix...")
K_train = quantum_kernel(X_train, X_train)
K_test = quantum_kernel(X_test, X_train)

# Train SVM with quantum kernel
svm = SVC(kernel='precomputed')
svm.fit(K_train, y_train)
y_pred = svm.predict(K_test)
print(f"Quantum SVM accuracy: {accuracy_score(y_test, y_pred):.3f}")
```

---

## Quantum Neural Networks (QNN)

```python
import pennylane as qml
import torch
import numpy as np

# PennyLane + PyTorch integration
n_qubits = 4
n_layers = 3
dev = qml.device('default.qubit', wires=n_qubits)

@qml.qnode(dev, interface='torch')
def quantum_neural_net(inputs, weights):
    """Quantum neural network as a PyTorch layer."""
    # Angle embedding
    qml.AngleEmbedding(inputs, wires=range(n_qubits))
    
    # Parameterized layers
    for layer_weights in weights:
        for i in range(n_qubits):
            qml.RY(layer_weights[i, 0], wires=i)
            qml.RZ(layer_weights[i, 1], wires=i)
        for i in range(n_qubits - 1):
            qml.CNOT(wires=[i, i+1])
    
    return [qml.expval(qml.PauliZ(i)) for i in range(n_qubits)]

# Define as PyTorch module
class QuantumLayer(torch.nn.Module):
    def __init__(self, n_qubits, n_layers):
        super().__init__()
        weight_shapes = (n_layers, n_qubits, 2)
        self.weights = torch.nn.Parameter(
            torch.randn(*weight_shapes)
        )
    
    def forward(self, x):
        return torch.stack([
            torch.tensor(quantum_neural_net(xi, self.weights), dtype=torch.float32)
            for xi in x
        ])

class HybridClassifier(torch.nn.Module):
    def __init__(self):
        super().__init__()
        self.quantum = QuantumLayer(n_qubits, n_layers)
        self.classical = torch.nn.Linear(n_qubits, 2)  # 2-class output
    
    def forward(self, x):
        q_out = self.quantum(x)
        return self.classical(q_out)

# Train
model = HybridClassifier()
optimizer = torch.optim.Adam(model.parameters(), lr=0.01)
criterion = torch.nn.CrossEntropyLoss()

# Example training loop
X = torch.randn(50, n_qubits)  # 50 samples, 4 features
y = torch.randint(0, 2, (50,))  # binary labels

for epoch in range(10):
    optimizer.zero_grad()
    out = model(X[:10])  # small batch due to quantum simulation cost
    loss = criterion(out, y[:10])
    loss.backward()
    optimizer.step()
    print(f"Epoch {epoch}: loss = {loss.item():.4f}")
```

---

## Quantum Transfer Learning

Use a pre-trained classical model + quantum circuit:

```python
import torch
import torchvision
import pennylane as qml

# Classical pre-trained feature extractor
classical_model = torchvision.models.resnet18(pretrained=True)
classical_model.fc = torch.nn.Identity()  # remove final layer → 512-dim features

n_qubits = 4
dev = qml.device('default.qubit', wires=n_qubits)

@qml.qnode(dev, interface='torch')
def quantum_layer(inputs, weights):
    qml.AngleEmbedding(inputs[:n_qubits], wires=range(n_qubits))
    qml.BasicEntanglerLayers(weights, wires=range(n_qubits))
    return [qml.expval(qml.PauliZ(i)) for i in range(n_qubits)]

class QuantumTransferModel(torch.nn.Module):
    def __init__(self):
        super().__init__()
        self.classical = classical_model
        self.reducer = torch.nn.Linear(512, n_qubits)
        self.q_weights = torch.nn.Parameter(
            torch.randn(2, n_qubits)
        )
        self.classifier = torch.nn.Linear(n_qubits, 2)
    
    def forward(self, x):
        features = self.classical(x)          # 512-dim
        reduced = torch.tanh(self.reducer(features)) * np.pi  # n_qubits-dim, scaled
        q_out = torch.stack([
            torch.tensor(quantum_layer(r, self.q_weights))
            for r in reduced
        ])
        return self.classifier(q_out)
```

---

## Quantum Generative Models

Quantum circuit born machines (QCBMs) as generative models:

```python
import pennylane as qml
import numpy as np

dev = qml.device('default.qubit', wires=4)

@qml.qnode(dev)
def generative_circuit(params):
    """Parameterized circuit that generates probability distribution."""
    for layer in range(3):
        for q in range(4):
            qml.RY(params[layer, q, 0], wires=q)
            qml.RZ(params[layer, q, 1], wires=q)
        for q in range(3):
            qml.CNOT(wires=[q, q+1])
    
    return qml.probs(wires=range(4))  # probability distribution over 16 states

# Target distribution (e.g., from data)
target_dist = np.array([0.1, 0.2, 0.05, 0.15, 0.0, 0.3, 0.1, 0.05,
                        0.02, 0.01, 0.01, 0.01, 0.0, 0.0, 0.0, 0.0])
target_dist /= target_dist.sum()  # normalize

# Train to match distribution using KL divergence
params = np.random.uniform(0, np.pi, (3, 4, 2), requires_grad=True)

def kl_divergence(params):
    gen_dist = generative_circuit(params)
    return np.sum(target_dist * np.log(target_dist / (gen_dist + 1e-10)))

opt = qml.GradientDescentOptimizer(0.1)
for step in range(50):
    params, cost = opt.step_and_cost(kl_divergence, params)
    if step % 10 == 0:
        print(f"Step {step}: KL divergence = {cost:.4f}")
```

---

## Quantum Convolutional Neural Networks (QCNN)

```python
import pennylane as qml
import numpy as np

dev = qml.device('default.qubit', wires=8)

def conv_layer(params, wires):
    """Quantum convolutional layer: 2-qubit gate applied to pairs."""
    for i, (w1, w2) in enumerate(zip(wires[::2], wires[1::2])):
        qml.RY(params[i, 0], wires=w1)
        qml.RY(params[i, 1], wires=w2)
        qml.CNOT(wires=[w1, w2])

def pool_layer(wires):
    """Quantum pooling: measure half the qubits, condition on result."""
    for w1, w2 in zip(wires[::2], wires[1::2]):
        # Simplified pooling: CNOT and discard first qubit
        qml.CNOT(wires=[w1, w2])

@qml.qnode(dev)
def qcnn(x, params):
    # Encode 8-dimensional data
    qml.AngleEmbedding(x, wires=range(8))
    
    # Conv layer 1: acts on all 8 qubits in pairs
    conv_layer(params[0], wires=range(8))
    
    # Pool: 8 → 4 qubits
    pool_layer(wires=range(8))
    
    # Conv layer 2: acts on remaining 4 qubits
    conv_layer(params[1], wires=range(1, 8, 2))
    
    # Pool: 4 → 2 qubits
    pool_layer(wires=range(1, 8, 2))
    
    # Final measurement
    return qml.expval(qml.PauliZ(3))
```

---

## QML Reality Check

```
Claims about QML speedup — what's real?

POTENTIALLY REAL:
  Quantum kernel methods: If data has quantum structure, kernel
  computation is classically hard → genuine speedup.
  
  Quantum simulation for chemistry/materials: Yes, real speedup.

PROBABLY NOT REAL (yet):
  Quantum speedup for classical data: Most QML algorithms have
  "dequantization" — equally fast classical algorithms exist.
  
  Liu et al. (2021): Classification speedup claims are based on
  assumptions that hide classical algorithms of similar complexity.

THE HONEST PICTURE:
  QML is a nascent field. Most results are proofs-of-concept.
  True quantum advantage for ML requires:
    1. Fault-tolerant quantum computers
    2. Problems with genuine quantum structure
    3. Better classical lower bounds

BEST BET: Quantum simulation for drug discovery.
          Processing data that is itself quantum (quantum sensors, quantum networks).
```

---

## Practical QML Workflow

```
1. Classical preprocessing:
   - Feature engineering
   - Dimensionality reduction (must fit in n qubits = 2ⁿ amplitudes)
   - Normalization

2. Feature encoding choice:
   - Angle encoding: n features → n qubits (simpler, more practical)
   - Amplitude encoding: 2ⁿ features → n qubits (expressive but hard to prepare)
   - Basis encoding: binary data → computational basis

3. Ansatz choice:
   - Hardware-efficient: short depth, hardware-friendly
   - Problem-inspired: better expressibility, harder to run

4. Training:
   - Parameter shift rule: exact gradients
   - SPSA: gradient-free, works with shot noise
   - Watch for barren plateaus

5. Evaluation:
   - Compare to classical baselines (SVM, neural nets)
   - Honest about shot noise and circuit depth requirements
```
# Chapter 19: Quantum Hardware

## Physical Implementations of Qubits

Different technologies encode quantum information in different physical systems:

| Technology | Company | Qubit type | Gate fidelity | T₂ coherence |
|---|---|---|---|---|
| Superconducting | IBM, Google | Energy levels of Josephson junction | 99.5% (1Q), 99% (2Q) | ~100 μs |
| Trapped ion | IonQ, Quantinuum | Electronic states of ions | 99.9% (1Q), 99.5% (2Q) | >10 minutes |
| Photonic | PsiQuantum, Xanadu | Photon polarization/path | Variable | Infinite (no decoherence) |
| Neutral atom | QuEra, Pasqal | Rydberg states of neutral atoms | 99.5% (1Q), 99% (2Q) | ~5 s |
| Silicon spin | Intel | Electron spin in silicon | ~99% | ms |
| Topological | Microsoft | Non-Abelian anyons | Theoretical | Inherently protected |

---

## Superconducting Qubits (IBM, Google)

The most advanced commercial technology:

```
Physics:
  Josephson junction: superconductor-insulator-superconductor
  At ~15 mK (colder than outer space), becomes quantum harmonic oscillator
  Anharmonicity makes it a 2-level system (qubit)
  
  Qubit types: Transmon (most common), Flux qubit, Charge qubit
  
Gate operations:
  Single-qubit: microwave pulses (~10-50 ns)
  Two-qubit: cross-resonance, CZ via flux tuning (~100-500 ns)
  
Cooling requirements: dilution refrigerator (~15 mK)
Scale: 50-1000+ physical qubits on chip

IBM Quantum roadmap:
  2021: 127 qubits (Eagle)
  2022: 433 qubits (Osprey)
  2023: 1121 qubits (Condor)
  2025: 100,000+ qubits (modular systems)
  
Connectivity:
  Heavy-hex topology: each qubit connects to 2-3 others
  Not all-to-all → SWAP gates needed for non-adjacent operations
```

```
IBM Heavy-Hex Topology (Eagle, 127 qubits):
  ○─○─○─○─○─○─○─○─○─○─○─○
  │   │   │   │   │   │  
  ○   ○   ○   ○   ○   ○  
  │   │   │   │   │   │  
  ○─○─○─○─○─○─○─○─○─○─○─○
```

---

## Trapped Ion Qubits (IonQ, Quantinuum)

```
Physics:
  Individual ions (e.g., Yb⁺, Ca⁺) trapped by electromagnetic fields
  Qubit: two electronic energy levels of the ion
  Gates: laser pulses (optical, ~1-100 μs) or microwave pulses
  
Advantages:
  Very long coherence times (>10 minutes vs ~100 μs for superconducting)
  High gate fidelity (>99.9% single-qubit, >99.5% two-qubit)
  All-to-all connectivity (any ion can interact with any other)
  No need for dilution refrigerator (room temperature trap)
  
Disadvantages:
  Slow gate speeds (μs vs ns for superconducting)
  Difficult to scale to thousands of qubits
  Photonic interconnects needed for modular systems
  
Key companies:
  IonQ: Ytterbium ions, publicly traded, AWS/Azure integration
  Quantinuum (Honeywell): Barium ions, H-series, highest fidelity
  Oxford Ionics: Erbium, chip-scale
  
Scale: 20-30 "algorithmic qubits" currently
```

---

## Photonic Quantum Computing

```
Physics:
  Qubits: photon polarization (H/V) or dual-rail (path)
  Gates: beam splitters, phase shifters, photon detectors
  
Types:
  Gaussian Boson Sampling: Xanadu's Borealis
  Linear optical QC (LOQC): probabilistic gates
  Fusion-based QC: PsiQuantum's approach
  
Advantages:
  Room temperature operation (no cooling)
  Natural for quantum communication/networks
  Photons are perfect flying qubits
  
Disadvantages:
  Photon loss (each optical element absorbs some photons)
  Probabilistic gates (entangling gates succeed with probability <1)
  Measurement is destructive
  
PsiQuantum approach:
  Fault-tolerant silicon photonics
  Manufacturing at GlobalFoundries (standard semiconductor fab)
  Target: 1 million+ photonic qubits
```

---

## Neutral Atom Qubits (QuEra, Pasqal)

```
Physics:
  Neutral atoms (Rb, Cs) held in optical tweezers (focused laser beams)
  Ground/excited states or hyperfine levels as qubit states
  Long-range interactions via Rydberg states (highly excited atoms)
  
Advantages:
  Large-scale connectivity via Rydberg interactions
  Reconfigurable: move atoms to create different topologies
  Long coherence times
  Hundreds to thousands of qubits
  
Disadvantages:
  Moderate gate fidelities (~99.5%)
  Atom loading is probabilistic
  Vacuum system required
  
QuEra:
  Aquila system: 256 neutral atom qubits (publicly available via AWS Braket)
  Programmable Hamiltonian simulation
  
Notable result: Harvard/QuEra (2023) demonstrated 48 logical qubits via error correction
  First demonstration of fault-tolerant algorithms at this scale
```

---

## Hardware Metrics

### Gate fidelity:
```
Average Gate Fidelity F = (1 - ε)
  ε = error per gate

Single-qubit fidelity: typically 99.5–99.9%
Two-qubit fidelity: typically 99.0–99.7%
State preparation and measurement (SPAM) error: ~1-5%

Error per gate needed for fault-tolerance: < ~1%
Current superconducting: ~0.1-1% (at the threshold)
Current trapped ion: ~0.01-0.1% (better!)
```

### Quantum Volume:
```
IBM's benchmark for processor capability:
  QV = max m such that m-qubit random circuit has heavy output > 2/3

Factors: qubit count, connectivity, gate fidelity, measurement error

QV progression:
  2017: QV=4
  2019: QV=16
  2020: QV=64
  2022: QV=512
  2023: QV=8192 (IBM Torino)

Higher QV = more complex circuits can run successfully
```

### Circuit Layer Operations per Second (CLOPS):
```
Speed metric: how fast can a backend run parametric circuit layers?
  Includes compilation, execution, and result retrieval

IBM Eagle: ~5,000 CLOPS
Important for VQA optimization (many circuit evaluations needed)
```

---

## Quantum Cloud Access

```
IBM Quantum:
  Free tier: Falcon (5Q), Heron (133Q), Eagle (127Q)
  Access via qiskit-ibm-runtime
  Jobs queue (wait time varies)

Google Quantum AI:
  Research collaboration required
  Access via cirq and Google Cloud

AWS Braket:
  Multiple providers: IonQ, Rigetti, OQC, QuEra, IQM
  Pay per shot/task
  
Azure Quantum:
  IonQ, Quantinuum, Rigetti
  Quantum credits for researchers
  
IonQ:
  Harmony, Aria, Forte systems
  Available on AWS/Azure/GCP

Quantinuum:
  H1-1, H2 systems (highest fidelity available)
  Commercial and research access
```

```python
# Connect to IBM Quantum (free tier)
from qiskit_ibm_runtime import QiskitRuntimeService

service = QiskitRuntimeService(channel="ibm_quantum", token="YOUR_TOKEN")
backends = service.backends(operational=True, min_num_qubits=5)
for b in backends[:5]:
    print(f"{b.name}: {b.num_qubits}Q, queue={b.status().pending_jobs}")

# Connect to AWS Braket
# pip install amazon-braket-sdk
from braket.aws import AwsDevice
from braket.circuits import Circuit

device = AwsDevice("arn:aws:braket:::device/quantum-simulator/amazon/sv1")
circuit = Circuit().h(0).cnot(0, 1)
task = device.run(circuit, shots=1000)
result = task.result()
print(result.measurement_counts)
```

---

## Noise Characterization

Understanding hardware noise is critical for algorithm design:

```python
from qiskit_ibm_runtime import QiskitRuntimeService
from qiskit.providers import BackendV2

# Get backend properties
service = QiskitRuntimeService()
backend = service.backend('ibm_brisbane')

# View coupling map
print(backend.coupling_map)

# View error rates (via calibration data)
from qiskit.transpiler import CouplingMap
cm = backend.coupling_map

# Randomized benchmarking (estimate gate errors)
# (requires qiskit-experiments)
from qiskit_experiments.library import StandardRB

qubits = [0]
rb_exp = StandardRB(qubits, lengths=range(1, 200, 20), backend=backend)
rb_data = rb_exp.run(backend).block_for_results()
result_data = rb_data.analysis_results("EPC")
print(f"1Q gate error: {result_data.value:.4f}")
```

---

## Roadmap to Fault-Tolerant QC

```
NISQ era (now, 2024):
  Physical qubits: 50–1000
  Error rate: 0.1–1%
  No error correction
  Useful for: NISQ algorithms (VQE, QAOA), quantum simulation demo

Early FTQC era (2027–2030?):
  Physical qubits: 10,000–100,000
  Logical qubits: 10–100 (with surface codes)
  Surface code threshold crossed
  Useful for: Small-scale error-corrected computations

Full FTQC era (2030s+):
  Physical qubits: millions
  Logical qubits: thousands
  Useful for: Shor's algorithm, large-scale quantum simulation, quantum chemistry

Key milestones to watch:
  1. Sustained error below surface code threshold (~1%)
  2. First logical qubit outperforming physical qubits
  3. Fault-tolerant T gate via magic state distillation at scale
  4. Break-even: fault-tolerant computation faster than classical
```

---

## Comparing Hardware Technologies

```
For algorithm developers:
  Short-term advantage: Trapped ion (high fidelity, all-to-all)
  Long-term scale: Superconducting (faster fabrication scaling)
  New paradigm: Neutral atoms (reconfigurable, large-scale)
  
For networking/communication: Photonic
For topological protection: Microsoft's topological (not yet realized)

"Horses for courses":
  Optimization (short depth): Superconducting (speed)
  Deep circuits (high fidelity): Trapped ion
  Quantum simulation: Neutral atoms (programmable Hamiltonians)
  Communication: Photonic
```
# Chapter 20: Learning Roadmap & Problem-Solving Approach

## Overview

A structured 6-week plan to master quantum computing from fundamentals to expert-level problem solving, with a problem-first approach inspired by competitive programming.

```
Week 1: Mathematical Foundations
Week 2: Quantum Circuits and Simulation
Week 3: Core Algorithms (Deutsch-Jozsa, Simon, Grover)
Week 4: Advanced Algorithms (QFT, Shor, QPE)
Week 5: Near-Term Algorithms (VQE, QAOA, QML)
Week 6: Hardware, Error Correction, and Research Frontier
```

---

## Week 1: Mathematical Foundations

**Goal:** Be fluent in the math language of quantum computing.

### Study:
- Complex numbers: polar form, Euler's formula e^(iθ)
- Linear algebra: vectors, matrices, inner product, tensor product
- Dirac notation: kets, bras, bra-kets
- Unitary matrices and their properties
- Bloch sphere representation

### Practice Problems:
```python
import numpy as np

# Problem 1: Normalize this state vector
v = np.array([1+1j, 2-1j])
# Solution: divide by norm
norm = np.linalg.norm(v)
normalized = v / norm
assert abs(np.dot(normalized.conj(), normalized) - 1) < 1e-10

# Problem 2: Apply Hadamard to |1⟩, verify result
H = np.array([[1, 1], [1, -1]]) / np.sqrt(2)
one = np.array([0, 1])
result = H @ one
# Should be (|0⟩-|1⟩)/√2 = |-⟩
assert np.allclose(result, np.array([1/np.sqrt(2), -1/np.sqrt(2)]))

# Problem 3: Compute tensor product |+⟩⊗|0⟩
plus = np.array([1, 1]) / np.sqrt(2)
zero = np.array([1, 0])
product = np.kron(plus, zero)
# Should be [1/√2, 0, 1/√2, 0]

# Problem 4: Apply CNOT to |+⟩|0⟩ → Bell state?
CNOT = np.array([[1,0,0,0],[0,1,0,0],[0,0,0,1],[0,0,1,0]])
bell = CNOT @ product
# Should be (|00⟩+|11⟩)/√2

# Problem 5: Compute ⟨Ψ|Z⊗Z|Ψ⟩ for Bell state
Z = np.array([[1, 0], [0, -1]])
ZZ = np.kron(Z, Z)
exp_ZZ = bell.conj() @ ZZ @ bell
print(f"⟨ZZ⟩ = {exp_ZZ.real}")  # Should be +1
```

### Milestone:
- [ ] Can represent any single/multi-qubit state as vector
- [ ] Can apply gates by matrix multiplication
- [ ] Understand why unitary = reversible = norm-preserving

---

## Week 2: Quantum Circuits and Programming

**Goal:** Build and simulate arbitrary quantum circuits in Qiskit.

### Study:
- All standard gates (H, X, Y, Z, S, T, CNOT, CZ, Toffoli)
- Circuit composition and decomposition
- Quantum parallelism and phase kickback
- Running simulations (statevector and sampling)

### Project: Circuit Builder Challenge
```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator
from qiskit.quantum_info import Statevector

# Challenge 1: Build a 4-qubit GHZ state
# (|0000⟩ + |1111⟩)/√2
def ghz_4():
    qc = QuantumCircuit(4)
    # YOUR CODE HERE
    qc.h(0)
    for i in range(3):
        qc.cx(i, i+1)
    return qc

# Challenge 2: Implement quantum teleportation
# Given qubit in |+⟩, teleport from q0 to q2

# Challenge 3: Implement superdense coding
# Send '10' using entangled pair

# Challenge 4: Build the Bernstein-Vazirani circuit for s='1011'

# Challenge 5: Build Grover's oracle for |101⟩
# Verify: after 1 Grover iteration, probability of |101⟩ > 0.9

# Verify your circuits:
for name, qc in [("GHZ-4", ghz_4())]:
    sv = Statevector(qc)
    probs = sv.probabilities_dict()
    print(f"{name}: {sorted(probs.items(), key=lambda x: -x[1])[:3]}")
```

### Milestone:
- [ ] Can implement Bell states, GHZ, teleportation from memory
- [ ] Comfortable with Qiskit API
- [ ] Understand phase kickback mechanism

---

## Week 3: First Quantum Algorithms

**Goal:** Implement and understand the first generation of quantum algorithms.

### Algorithms to implement:
1. Deutsch's algorithm
2. Deutsch-Jozsa algorithm (n-qubit)
3. Bernstein-Vazirani
4. Simon's algorithm (with classical post-processing)
5. Grover's search (with diffuser)

### Problem-Solving Approach:
```
For each algorithm, answer:
1. What problem does it solve?
2. What's the classical complexity?
3. What's the quantum complexity?
4. What mathematical structure does the oracle encode?
5. How does interference help?
6. How is the answer extracted?
```

### Challenge Problems:
```python
# Challenge A: Multi-solution Grover
# Given oracle marking 3 items in a 32-item database
# How many iterations? Implement it.
N = 32; k = 3
import numpy as np
iterations = int(np.pi / 4 * np.sqrt(N/k))
print(f"Optimal iterations: {iterations}")  # 2

# Challenge B: Solve a 3-bit Simon's problem
# f has period s='110', verify your algorithm finds it

# Challenge C: Use Grover to solve 3-SAT on 3 variables
# (x0 OR x1) AND (NOT x1 OR x2) AND (x0 OR NOT x2)
# Build oracle for the satisfying assignment
```

---

## Week 4: Core Algorithm Techniques

**Goal:** Master QFT, phase estimation, and Shor's algorithm.

### Topics:
- Quantum Fourier Transform: circuit and verification
- Quantum Phase Estimation: eigenvalue finding
- Shor's factoring algorithm: classical reduction

### Problems:
```python
from qiskit.circuit.library import QFT
from qiskit.quantum_info import Statevector
import numpy as np

# Problem 1: QFT of |5⟩ in 4 qubits
qc = QuantumCircuit(4)
qc.x(0); qc.x(2)  # |5⟩ = |0101⟩
qft = QFT(4)
qc.compose(qft, inplace=True)
sv = Statevector(qc)
print("QFT|5⟩ amplitudes:", np.round(np.abs(sv.data), 3))
# All equal (1/4) — uniform distribution

# Problem 2: Use QFT to detect period 4 in 16-element register
# Create state |0⟩+|4⟩+|8⟩+|12⟩ → apply QFT → verify peaks at 0,4,8,12

# Problem 3: QPE for T gate
# T|1⟩ = e^(iπ/4)|1⟩ → phase = 1/8 = 0.001 in binary
# Use 4-bit QPE to estimate → should measure '0010'

# Problem 4: Factor 15 using Shor's
# Try a=7: verify period r=4, then gcd(7^2-1, 15) and gcd(7^2+1, 15)
from math import gcd
a, N = 7, 15
r = 4  # (found by quantum period finding)
p = gcd(pow(a, r//2) - 1, N)
q = gcd(pow(a, r//2) + 1, N)
print(f"Factors of {N}: {p} and {q}")  # 3 and 5
```

---

## Week 5: Near-Term Algorithms and Hardware

**Goal:** Implement and tune VQE, QAOA, and understand quantum hardware.

### Projects:

**VQE Project: H₂ ground state**
```python
# Goal: Find ground state energy of H₂ molecule
# Use different ansätze, compare results
# Measure: accuracy vs circuit depth vs number of parameters

# Ansätze to compare:
# 1. Hardware-efficient (Ry+CNOT layers)
# 2. UCCSD (chemistry-inspired)
# 3. Problem-inspired
```

**QAOA Project: MaxCut optimization**
```python
# Goal: Find MaxCut on a 6-node graph
# 1. Implement QAOA for p=1,2,3
# 2. Compare to classical greedy algorithm
# 3. Analyze: does more depth help?
import networkx as nx
G = nx.random_regular_graph(3, 6, seed=42)
# Expected MaxCut = ? (find classically first for validation)
```

**Hardware Project: Run on real IBM quantum device**
```python
# 1. Build a simple Bell state circuit
# 2. Transpile for IBM backend
# 3. Run with 1024 shots
# 4. Compare ideal vs noisy results
# 5. Apply readout error mitigation
```

---

## Week 6: Error Correction and Research Frontier

**Goal:** Understand quantum error correction and identify research problems.

### Study:
- Bit flip code (3-qubit)
- Shor's 9-qubit code
- Surface code basics
- Noise models and mitigation

### Research Reading:
```
Key papers to read:
1. Grover (1996): "A fast quantum mechanical algorithm for database search"
2. Shor (1994): "Algorithms for quantum computation: discrete logarithms and factoring"
3. Preskill (2018): "Quantum Computing in the NISQ Era and Beyond"
4. Arute et al. (2019): "Quantum supremacy using a programmable superconducting processor"
5. Bharti et al. (2022): "Noisy intermediate-scale quantum algorithms" (review)
6. Cerezo et al. (2021): "Variational quantum algorithms" (review)
```

---

## The Problem-Solving Framework

**Apply this to every quantum algorithm problem:**

```
STEP 1: UNDERSTAND THE PROBLEM
  - What is the input? Output? Size n?
  - What does "solving" mean classically?
  - Classical complexity?

STEP 2: FIND THE STRUCTURE
  Ask: Is there periodicity? → QFT-based
  Ask: Is there a checkable oracle? → Grover/amplitude amplification
  Ask: Is it an eigenvalue problem? → QPE
  Ask: Is it optimization? → VQE/QAOA
  Ask: Is it search with structure? → Look for quantum walk or algebraic structure

STEP 3: DESIGN THE ORACLE/UNITARY
  - What does the oracle compute?
  - How to implement it reversibly?
  - How to uncompute ancilla?

STEP 4: ANALYZE COMPLEXITY
  - How many oracle calls?
  - Total gates? Circuit depth?
  - Speedup over classical?

STEP 5: IMPLEMENT
  - Start with small n in Qiskit simulation
  - Verify classical test cases
  - Check interference is happening correctly

STEP 6: VERIFY
  - Does the correct answer have high probability?
  - Check with known cases
  - Compare to classical result
```

---

## Common Problem Patterns

```
Pattern 1: "Find x such that f(x) = y"
→ Grover's search if unstructured
→ If f has period: Simon's / Shor-like

Pattern 2: "What is the period/order of f?"
→ QFT / Shor-like quantum period finding

Pattern 3: "Find ground state of Hamiltonian H"
→ VQE (short-depth NISQ) or QPE (fault-tolerant)

Pattern 4: "Solve optimization problem with binary variables"
→ QAOA (combinatorial), VQE (with problem Hamiltonian)

Pattern 5: "Classify data points"
→ Quantum kernel SVM, VQC, QCNN

Pattern 6: "Solve linear system Ax = b"
→ HHL algorithm (exponential speedup under conditions)

Pattern 7: "Sample from a complex probability distribution"
→ Quantum Gibbs sampling, boson sampling
```

---

## Competitive Programming Mindset Applied to Quantum

```
Classical CP skill → Quantum analog

Brute force (O(2ⁿ))    → Superposition + oracle (O(√2ⁿ) via Grover)
DFT / FFT               → QFT (O(n²) instead of O(N·n))
GCD / number theory     → Quantum period finding + continued fractions
Graph algorithms        → Quantum walks (speedup for some problems)
Dynamic programming     → Phase kickback accumulation (amplitude encoding)

Think of quantum algorithm = classical algorithm + interference
Interference = your "loop optimization" trick in QC
```

---

## Essential Resources

| Resource | Type | Why |
|---|---|---|
| Nielsen & Chuang "Quantum Computation and Quantum Information" | Textbook | The Bible of QC — start here |
| Preskill's lecture notes (theory.caltech.edu/~preskill/ph219) | Notes | Free, comprehensive |
| Qiskit Textbook (learning.quantum.ibm.com) | Interactive | Hands-on with Qiskit |
| Quanta Magazine (quantum topics) | Articles | Accessible research coverage |
| arXiv quant-ph | Papers | Latest research |
| IBM Quantum Lab | Cloud | Free hardware access |
| Quirk (algassert.com/quirk) | Browser tool | Visual circuit building |
| PennyLane documentation | Docs | QML and VQA |
| Quantum Computing Stack Exchange | Q&A | Problem solving help |

---

## Career Paths

```
1. Quantum Algorithm Research
   Skills: QC theory, complexity theory, algorithm design
   Jobs: IBM Research, Google Quantum AI, Microsoft Research, universities
   
2. Quantum Software Engineering
   Skills: Qiskit, PennyLane, quantum compilation, error mitigation
   Jobs: IBM, PennyLane/Xanadu, IonQ, startups
   
3. Quantum Machine Learning
   Skills: ML + quantum circuits, PennyLane, autodiff
   Jobs: QC startups, pharmaceutical companies, finance
   
4. Quantum Hardware
   Skills: Physics, cryogenic engineering, microwave engineering
   Jobs: IBM, Google, IonQ, Quantinuum
   
5. Post-Quantum Cryptography
   Skills: Cryptography, mathematics
   Jobs: NIST, NSA, cloud providers (AWS, Azure, GCP), banks

6. Quantum-Safe Security Consulting
   Skills: Security + basic QC knowledge
   Jobs: Consulting firms, financial sector
```

---

## 30-Day Quick Start (If You're in a Hurry)

```
Days 1-5:   Read Ch 1-4 of Nielsen & Chuang. Run Qiskit "Hello World."
Days 6-10:  Implement Bell states, GHZ, teleportation in Qiskit.
Days 11-15: Implement and run Deutsch-Jozsa and Grover's from scratch.
Days 16-20: Implement QFT. Read Shor's paper. Simulate period finding.
Days 21-25: Implement VQE for H₂. Learn PennyLane basics.
Days 26-30: Run circuits on IBM Quantum hardware. Read Preskill's NISQ paper.

After 30 days: You can read most quantum computing papers with comprehension.
```
