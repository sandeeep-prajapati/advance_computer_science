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
