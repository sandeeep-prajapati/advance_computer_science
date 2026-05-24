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
