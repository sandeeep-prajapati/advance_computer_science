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
