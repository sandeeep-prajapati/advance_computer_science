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
