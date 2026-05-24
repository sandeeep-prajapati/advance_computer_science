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
