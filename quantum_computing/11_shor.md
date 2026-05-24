# Chapter 11: Shor's Factoring Algorithm

## Why Shor's Algorithm Matters

```
RSA encryption (used in HTTPS, banking, email):
  Security relies on: factoring n = p·q (two large primes) is HARD

Best classical: General Number Field Sieve
  Time: O(exp(n^(1/3) · (log n)^(2/3)))  — superpolynomial in bit length

Shor's algorithm (1994):
  Time: O(n³) — polynomial!  (n = number of bits in the number)

Impact: A large enough quantum computer would break RSA.
Current cryptography: migrating to post-quantum standards (NIST 2024).
```

---

## Problem Reduction: Factoring → Period Finding

Shor's insight: factoring reduces to **order finding**.

```
STEP 1: Choose random a with 1 < a < N and gcd(a, N) = 1
        (If gcd(a,N) ≠ 1, you've found a factor already — done!)

STEP 2: Find the ORDER r of a modulo N:
        r = smallest positive integer such that aʳ ≡ 1 (mod N)

STEP 3: If r is odd or aʳ/² ≡ -1 (mod N): restart with new a

STEP 4: Otherwise:
        gcd(aʳ/² - 1, N) and gcd(aʳ/² + 1, N) are factors of N!

WHY? Because aʳ ≡ 1 (mod N) means:
  (aʳ/² - 1)(aʳ/² + 1) ≡ 0 (mod N)
  So N divides this product → factors appear in the GCDs.

QUANTUM PART: Finding r is the hard part — use quantum period finding.
CLASSICAL PART: Steps 1, 3, 4 are fast classical operations.
```

---

## Classical Example

```
Factor N = 15, choose a = 7

Order finding: 7ˢ mod 15 for s = 1, 2, 3, ...
  7¹ = 7
  7² = 49 = 3·15 + 4 → 4
  7³ = 7·4 = 28 = 15 + 13 → 13
  7⁴ = 7·13 = 91 = 6·15 + 1 → 1  ← period r = 4

r = 4 (even, and 7² = 4 ≠ -1 ≡ 14 mod 15) ✓

gcd(7² - 1, 15) = gcd(48, 15) = gcd(3·15+3, 15) = gcd(3, 15) = 3  ✓
gcd(7² + 1, 15) = gcd(50, 15) = gcd(3·15+5, 15) = gcd(5, 15) = 5  ✓

Factors: 3 and 5 (15 = 3 × 5) ✓
```

---

## Quantum Period Finding

The quantum algorithm finds r by:

```
1. Prepare superposition: (1/√2ⁿ) Σₓ |x⟩|0⟩

2. Compute modular exponentiation: |x⟩|0⟩ → |x⟩|aˣ mod N⟩

3. Measure second register → collapse to some value aˣ⁰ mod N
   First register collapses to: (1/√(2ⁿ/r)) Σⱼ |x₀ + j·r⟩  (periodic with period r)

4. Apply QFT to first register:
   → peaks at multiples of 2ⁿ/r

5. Measure first register → get j·(2ⁿ/r) for some integer j

6. Classical: continued fractions algorithm on j/2ⁿ gives r
```

---

## Circuit Architecture

```
n qubits for x register (n = 2·log₂(N), need ~2n input qubits)
n qubits for result register

q₀,...,q_{2n-1}: ─H⊗2n──[  modular  ]──QFT⁻¹── Measure
q_{2n},...,q_{3n-1}: ─────[exponent  ]────────── Measure

Controlled unitary: U_a,j: |x⟩|y⟩ → |x⟩|y·aˣ mod N⟩

Key expense: Modular exponentiation circuit (O(n³) gates using repeated squaring)
```

---

## Implementation for Small N

```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator
from qiskit.circuit.library import QFT
from fractions import Fraction
import numpy as np

def quantum_period_finding(N: int, a: int):
    """
    Find period r of f(x) = a^x mod N using quantum circuit.
    For small N only (demonstration purposes).
    """
    n = int(np.ceil(np.log2(N)))
    t = 2 * n  # ancilla bits for phase estimation
    
    qc = QuantumCircuit(t + n, t)
    
    # Initialize target register to |1⟩
    qc.x(t)  # first qubit of target register
    
    # Hadamard on counting register
    qc.h(range(t))
    
    # Controlled modular multiplication: a^(2^j) mod N
    for j in range(t):
        # Implement controlled-U^(2^j) where U|y⟩ = |a*y mod N⟩
        # For real implementation, use modular exponentiation circuit
        # Here we use a classical simulation trick
        _controlled_modular_mult(qc, a, j, N, t, n)
    
    # Inverse QFT on counting register
    qc.compose(QFT(t, inverse=True), qubits=range(t), inplace=True)
    
    # Measure counting register
    qc.measure(range(t), range(t))
    
    return qc

def shor_classical(N: int, measurement: int, t: int) -> int:
    """Classical post-processing: extract period from measurement."""
    # Use continued fractions to find r/N approximation
    phase = measurement / (2**t)
    frac = Fraction(phase).limit_denominator(N)
    r = frac.denominator
    return r

def factor_with_shor(N: int) -> tuple:
    """Full Shor's algorithm (classical simulation for small N)."""
    import random
    from math import gcd
    
    for _ in range(100):  # retry with different a
        a = random.randint(2, N - 1)
        g = gcd(a, N)
        if g != 1:
            return g, N // g  # Lucky! gcd gives factor directly
        
        # Find order r of a modulo N (classically for demo)
        r = 1
        val = a
        while val != 1:
            val = (val * a) % N
            r += 1
        
        if r % 2 == 1:
            continue  # odd period, retry
        
        half_power = pow(a, r // 2, N)
        if half_power == N - 1:  # aʳ/² ≡ -1 mod N
            continue
        
        p = gcd(half_power - 1, N)
        q = gcd(half_power + 1, N)
        
        if 1 < p < N:
            return p, N // p
        if 1 < q < N:
            return q, N // q
    
    return None, None

# Demonstrate Shor's classical reduction
print(factor_with_shor(15))   # (3, 5) or (5, 3)
print(factor_with_shor(21))   # (3, 7) or (7, 3)
print(factor_with_shor(35))   # (5, 7)
```

---

## Simplified 4-Qubit Demo (N=15)

IBM's team demonstrated Shor's on 5 qubits for N=15 (2001). A simplified version:

```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator

def shors_15_a7():
    """Simplified Shor's for N=15, a=7 (period r=4)."""
    # 4 counting qubits + 4 for modular arithmetic
    qc = QuantumCircuit(8, 4)
    
    # Initialize |1⟩ in work register
    qc.x(4)
    
    # Hadamard on counting register
    qc.h([0, 1, 2, 3])
    
    # Simplified modular exponentiation (7^x mod 15)
    # For real circuits, this needs many more gates
    qc.cx(3, 6)
    qc.cx(3, 7)
    qc.h(3)
    qc.swap(5, 7)
    qc.swap(4, 6)
    # ... (full circuit is more complex)
    
    # QFT inverse on counting register
    # ... (apply IQFT)
    
    qc.measure([0, 1, 2, 3], [0, 1, 2, 3])
    return qc
```

---

## Continued Fractions Algorithm

Key classical post-processing:

```python
from fractions import Fraction
import math

def find_period(measurement: int, t: int, N: int) -> int:
    """
    measurement: result from quantum register
    t: number of counting qubits
    N: the number to factor
    """
    if measurement == 0:
        return None  # trivial result, retry
    
    # Phase = measurement / 2^t
    phase = measurement / (2**t)
    
    # Find best rational approximation p/q with q ≤ N
    frac = Fraction(phase).limit_denominator(N)
    r = frac.denominator
    
    # Verify: a^r ≡ 1 mod N
    # (this check is done by the caller)
    return r

# Example: N=15, t=8, a=7, period r=4
# QFT peaks at j·(256/4) = j·64 for j = 0,1,2,3
for j in range(4):
    measurement = j * 64
    r = find_period(measurement, 8, 15)
    if r:
        print(f"measurement={measurement}, phase={measurement/256:.3f}, period={r}")
# measurement=64, phase=0.250, period=4 ✓
# measurement=128, phase=0.500, period=2 (r=2 → need r=lcm(2,...))
# measurement=192, phase=0.750, period=4 ✓
```

---

## Shor's Algorithm Complexity

```
Classical integer factoring (GNFS): O(exp(n^(1/3) · (log n)^(2/3)))
Shor's quantum:                     O(n³) gate operations, O(n) qubits

This is an EXPONENTIAL speedup.

For 2048-bit RSA key (n=2048):
  Classical: O(exp(2048^(1/3))) — astronomical
  Quantum:   O(2048³) = O(10^10) — feasible with fault-tolerant QC

Required quantum resources for 2048-bit RSA:
  Logical qubits: ~4000
  Physical qubits (with error correction): ~20 million
  Gate operations: ~10^12
  
Current NISQ computers: 50-1000 noisy qubits → far from breaking RSA.
Timeline: 2030s-2040s for fault-tolerant systems capable of this.
```

---

## Post-Quantum Cryptography

In response to Shor's threat:
```
NIST PQC Standards (2024):
  CRYSTALS-Kyber → Key encapsulation (replaces RSA key exchange)
  CRYSTALS-Dilithium → Digital signatures
  FALCON → Digital signatures (compact)
  SPHINCS+ → Hash-based signatures (conservative)

These are based on hard problems that quantum computers DON'T speed up significantly:
  Lattice problems (Learning With Errors)
  Hash functions (not significantly affected by Grover)
```
