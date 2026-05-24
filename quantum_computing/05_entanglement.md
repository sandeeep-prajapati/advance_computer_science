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
