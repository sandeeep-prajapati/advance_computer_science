# Chapter 2: Linear Algebra for Quantum Computing

## Why Linear Algebra?

Quantum mechanics IS linear algebra. Every quantum state is a vector, every operation is a matrix, and measurement is a projection. Master this and everything else follows.

---

## Vectors and Dirac Notation

### Ket (column vector):
```
|0⟩ = [1]    |1⟩ = [0]    |ψ⟩ = [α] = α|0⟩ + β|1⟩
      [0]          [1]           [β]
```

### Bra (row vector — conjugate transpose):
```
⟨0| = [1  0]    ⟨ψ| = [α*  β*]    (* denotes complex conjugate)
```

### Inner product (bra × ket → scalar):
```
⟨φ|ψ⟩ = φ†ψ    (dot product of complex vectors)

⟨0|0⟩ = 1    ⟨0|1⟩ = 0    ⟨1|1⟩ = 1    (orthonormal basis)

For |ψ⟩ = α|0⟩ + β|1⟩:
⟨ψ|ψ⟩ = |α|² + |β|² = 1    (normalization)
```

### Outer product (ket × bra → matrix):
```
|0⟩⟨0| = [1][1 0] = [1 0]    (projection onto |0⟩)
          [0]       [0 0]

|ψ⟩⟨ψ| = density matrix for pure state ψ
```

---

## Complex Numbers Review

Quantum amplitudes are complex numbers: `z = a + bi`

```
Conjugate:       z* = a - bi
Magnitude:       |z| = √(a² + b²)
Phase:           z = |z|·e^(iθ) = |z|(cos θ + i sin θ)    (Euler's formula)

Key: e^(iπ) = -1    (Euler's identity — phase flips sign)
```

In quantum mechanics, the *global phase* is unobservable:
```
|ψ⟩ and e^(iθ)|ψ⟩ represent the same physical state.
But *relative phase* between |0⟩ and |1⟩ components is crucial!

|+⟩ = (|0⟩ + |1⟩)/√2    ← relative phase = 0
|-⟩ = (|0⟩ - |1⟩)/√2    ← relative phase = π (very different state!)
```

---

## Matrices and Operations

### Matrix multiplication:
```
Gate U applied to state |ψ⟩:
U|ψ⟩ = (matrix) × (column vector)

Example — X gate (bit flip):
X = [0 1] applied to |0⟩ = [1]:
    [1 0]                  [0]

X|0⟩ = [0 1][1] = [0] = |1⟩  ✓
        [1 0][0]   [1]
```

### Conjugate transpose (Hermitian adjoint):
```
For matrix M, M† = (M*)ᵀ    (conjugate of each entry, then transpose)

X† = [0 1]† = [0 1] = X    (X is its own inverse — Hermitian and unitary)
     [1 0]    [1 0]
```

### Unitary matrices (all quantum gates are unitary):
```
U is unitary iff U†U = UU† = I    (preserves normalization)

This means:
  - All quantum gates are reversible
  - Inner products are preserved: ⟨Uφ|Uψ⟩ = ⟨φ|ψ⟩
  - Column (and row) vectors form orthonormal sets
```

### Hermitian matrices (observables):
```
H is Hermitian iff H† = H

Properties:
  - All eigenvalues are real (physical measurement values)
  - Eigenvectors form orthonormal basis
  
Example: Pauli Z matrix Z = [1  0 ] has eigenvalues +1 (for |0⟩) and -1 (for |1⟩)
                             [0 -1]
```

---

## Tensor Products (Multi-Qubit Systems)

The tensor product ⊗ combines systems:

```
|00⟩ = |0⟩⊗|0⟩ = [1] ⊗ [1] = [1·[1]] = [1]
                   [0]   [0]   [1·[0]]   [0]
                               [0·[1]]   [0]
                               [0·[0]]   [0]

So |00⟩ = [1,0,0,0]ᵀ, |01⟩ = [0,1,0,0]ᵀ, |10⟩ = [0,0,1,0]ᵀ, |11⟩ = [0,0,0,1]ᵀ
```

For gates on multi-qubit systems:
```
Apply H to qubit 0, I (identity) to qubit 1:
H⊗I = [1  1] ⊗ [1 0] = [H  0] = [1  0  1  0]
      [1 -1]   [0 1]   [0  H]   [0  1  0  1]
      ——/√2——                    [1  0 -1  0]
                                 [0  1  0 -1]  (all divided by √2)
```

### General rule:
```
n-qubit system has state space of dimension 2ⁿ
n-qubit gate is a 2ⁿ × 2ⁿ unitary matrix

This exponential growth is why quantum simulation is hard classically,
and why quantum computers can explore exponential state spaces.
```

---

## Eigenvalues and Eigenvectors

```
M|v⟩ = λ|v⟩    (|v⟩ is eigenvector, λ is eigenvalue)

For Pauli Z = [1  0]:
              [0 -1]

Z|0⟩ = +1·|0⟩    (|0⟩ is eigenstate with eigenvalue +1)
Z|1⟩ = -1·|1⟩    (|1⟩ is eigenstate with eigenvalue -1)

For Hadamard H = [1  1]/√2:
                 [1 -1]

H|+⟩ = +1·|+⟩    where |+⟩ = (|0⟩+|1⟩)/√2
H|-⟩ = -1·|-⟩    where |-⟩ = (|0⟩-|1⟩)/√2
```

---

## Key Bases

```
Computational (Z) basis:    |0⟩, |1⟩
Hadamard (X) basis:         |+⟩ = (|0⟩+|1⟩)/√2,  |-⟩ = (|0⟩-|1⟩)/√2
Y basis:                    |+i⟩= (|0⟩+i|1⟩)/√2, |-i⟩= (|0⟩-i|1⟩)/√2

Change of basis:
  H transforms between Z and X bases:
  H|0⟩ = |+⟩,    H|1⟩ = |-⟩
  H|+⟩ = |0⟩,    H|-⟩ = |1⟩
```

---

## The Pauli Matrices

The fundamental building blocks of quantum operations:

```
I = [1 0]    X = [0 1]    Y = [0 -i]    Z = [1  0]
    [0 1]        [1 0]        [i  0]        [0 -1]

Identity    Bit flip    Bit+Phase flip    Phase flip

Key properties:
  X² = Y² = Z² = I
  XY = iZ,  YZ = iX,  ZX = iY
  Pauli matrices are both unitary and Hermitian
  They form a basis for all 2×2 Hermitian matrices
```

---

## Important Linear Algebra Facts for QC

### Spectral decomposition:
```
Any Hermitian H = Σᵢ λᵢ|eᵢ⟩⟨eᵢ|    (sum over eigenvectors weighted by eigenvalues)

Measurement in a basis {|eᵢ⟩} gives outcome i with probability |⟨eᵢ|ψ⟩|²
```

### Matrix exponential (used for rotations):
```
e^(iθX/2) = cos(θ/2)I + i·sin(θ/2)X = Rx(θ) rotation gate

e^(iθZ/2) = [e^(iθ/2)    0      ] = Rz(θ) rotation gate
             [0        e^(-iθ/2)]
```

### Trace:
```
tr(M) = sum of diagonal elements = sum of eigenvalues

tr(|ψ⟩⟨ψ|) = 1    (pure state density matrix has trace 1)
```

---

## Problem-Solving with Linear Algebra

**Problem type 1: Is this state normalized?**
```python
import numpy as np

psi = np.array([1/np.sqrt(2), 1/np.sqrt(2)])  # |+⟩
print(np.dot(psi.conj(), psi))  # Should be 1.0
```

**Problem type 2: Apply a gate to a state**
```python
# Apply Hadamard to |0⟩
H = np.array([[1, 1], [1, -1]]) / np.sqrt(2)
zero = np.array([1, 0])  # |0⟩
result = H @ zero
print(result)  # [0.707, 0.707] = |+⟩
```

**Problem type 3: Compute measurement probabilities**
```python
# State: (|0⟩ + i|1⟩) / √2
psi = np.array([1/np.sqrt(2), 1j/np.sqrt(2)])
prob_0 = abs(psi[0])**2  # = 0.5
prob_1 = abs(psi[1])**2  # = 0.5
```

**Problem type 4: Is this gate unitary?**
```python
H = np.array([[1, 1], [1, -1]]) / np.sqrt(2)
print(np.allclose(H @ H.conj().T, np.eye(2)))  # True if unitary
```

**Problem type 5: Tensor product of two states**
```python
psi1 = np.array([1, 0])      # |0⟩
psi2 = np.array([1, 0])      # |0⟩
joint = np.kron(psi1, psi2)  # |00⟩ = [1,0,0,0]
```

---

## Quick Reference

| Operation | Math | Python (numpy) |
|---|---|---|
| Inner product | `⟨φ\|ψ⟩` | `phi.conj() @ psi` |
| Outer product | `\|ψ⟩⟨φ\|` | `np.outer(psi, phi.conj())` |
| Gate application | `U\|ψ⟩` | `U @ psi` |
| Tensor product | `A⊗B` | `np.kron(A, B)` |
| Conjugate transpose | `U†` | `U.conj().T` |
| Probability | `\|⟨i\|ψ⟩\|²` | `abs(psi[i])**2` |
| Norm | `\|\|ψ\|\|` | `np.linalg.norm(psi)` |
| Eigendecomposition | `Mv=λv` | `np.linalg.eigh(M)` |
