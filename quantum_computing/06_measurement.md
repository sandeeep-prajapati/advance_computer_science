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
