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
