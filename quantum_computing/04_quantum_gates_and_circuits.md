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
