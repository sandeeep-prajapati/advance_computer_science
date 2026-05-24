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
