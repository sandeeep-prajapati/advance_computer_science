# Chapter 8: Deutsch-Jozsa and Simon's Algorithm

## Deutsch's Problem (1985)

The first quantum algorithm to show exponential speedup (for a promise problem).

### Problem:
```
Given: Black-box function f:{0,1} → {0,1}
Promise: f is either:
  - Constant: f(0) = f(1)   (both same output)
  - Balanced: f(0) ≠ f(1)   (different outputs)
Determine which case.

Classical: 2 queries (must check both inputs)
Quantum:   1 query
```

### Solution (Deutsch's Algorithm):
```
Initial: |0⟩|1⟩
Step 1:  Apply H to both: |+⟩|-⟩ = (|0⟩+|1⟩)/√2 · (|0⟩-|1⟩)/√2
Step 2:  Apply Uf: |x⟩|y⟩ → |x⟩|y⊕f(x)⟩

After Uf (using phase kickback with |-⟩):
  = ((-1)^f(0)|0⟩ + (-1)^f(1)|1⟩)/√2 · |−⟩

Step 3: Apply H to first qubit:
  If f constant: (±|0⟩) → measurement gives 0
  If f balanced: (±|1⟩) → measurement gives 1

Step 4: Measure first qubit:
  0 → constant
  1 → balanced
```

```python
from qiskit import QuantumCircuit, QuantumRegister, ClassicalRegister
from qiskit_aer import AerSimulator

def deutsch_circuit(f_type: str):
    """
    f_type: 'constant_0', 'constant_1', 'balanced_id', 'balanced_not'
    """
    qr = QuantumRegister(2, 'q')
    cr = ClassicalRegister(1, 'c')
    qc = QuantumCircuit(qr, cr)
    
    # Initialize: |01⟩
    qc.x(1)
    
    # Apply H to both
    qc.h(0)
    qc.h(1)
    
    # Oracle
    if f_type == 'constant_0':
        pass                    # f(x) = 0: identity
    elif f_type == 'constant_1':
        qc.x(1)                 # f(x) = 1: flip ancilla (but phase cancels → still constant)
    elif f_type == 'balanced_id':
        qc.cx(0, 1)             # f(x) = x: CNOT
    elif f_type == 'balanced_not':
        qc.cx(0, 1)             # f(x) = NOT x
        qc.x(1)
    
    # Apply H to first qubit
    qc.h(0)
    
    # Measure
    qc.measure(0, 0)
    
    return qc

sim = AerSimulator()
for f_type in ['constant_0', 'constant_1', 'balanced_id', 'balanced_not']:
    qc = deutsch_circuit(f_type)
    counts = sim.run(qc, shots=100).result().get_counts()
    verdict = 'constant' if '0' in counts else 'balanced'
    print(f"{f_type}: {counts} → {verdict}")
# constant_0: {'0': 100} → constant
# constant_1: {'0': 100} → constant
# balanced_id: {'1': 100} → balanced
# balanced_not: {'1': 100} → balanced
```

---

## Deutsch-Jozsa Algorithm (1992)

Generalizes Deutsch's problem to n-bit functions.

### Problem:
```
Given: Black-box f:{0,1}ⁿ → {0,1}
Promise: f is either:
  - Constant: f(x) = c for all x
  - Balanced: f(x)=0 for exactly 2^(n-1) inputs, f(x)=1 for the other half
Determine which case.

Classical (deterministic): up to 2^(n-1)+1 queries (exponential worst case)
Classical (randomized): O(1) queries with bounded error — but Deutsch-Jozsa is exact
Quantum: 1 query (exact)
```

### Circuit:
```
|0⟩^⊗n ─H⊗n─[  ]─H⊗n─ Measure
|1⟩    ─H───[ Uf]─────

Step 1: Apply H to all n+1 qubits
  → (1/√2ⁿ) Σₓ |x⟩ ⊗ |−⟩

Step 2: Apply phase oracle Uf
  → (1/√2ⁿ) Σₓ (-1)^f(x)|x⟩ ⊗ |−⟩

Step 3: Apply H⊗n to first n qubits
  → Σₖ [(1/2ⁿ) Σₓ (-1)^(f(x)+x·k)] |k⟩

Step 4: Measure first n qubits:
  If f is constant: all amplitude at |00...0⟩ → measure 0...0
  If f is balanced: zero amplitude at |00...0⟩ → NEVER measure 0...0
```

**Mathematical proof:**
```
Amplitude at |0...0⟩ after H⊗n:
  = (1/2ⁿ) Σₓ (-1)^f(x)

If constant:    (1/2ⁿ) · 2ⁿ · (-1)^c = ±1   (large amplitude)
If balanced:    (1/2ⁿ) · 0 = 0               (exactly zero)

So measuring |0...0⟩ → constant
   measuring anything else → balanced
```

```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator
import numpy as np

def deutsch_jozsa(n: int, oracle_type: str = 'balanced'):
    """
    n: number of input qubits
    oracle_type: 'constant_0', 'constant_1', or 'balanced'
    """
    qc = QuantumCircuit(n + 1, n)
    
    # Initialize last qubit as |1⟩
    qc.x(n)
    
    # Apply H to all qubits
    qc.h(range(n + 1))
    
    qc.barrier()
    
    # Oracle
    if oracle_type == 'constant_0':
        pass
    elif oracle_type == 'constant_1':
        qc.x(n)
    elif oracle_type == 'balanced':
        # Simple balanced oracle: CNOT from each input qubit to output
        for i in range(n):
            qc.cx(i, n)
    
    qc.barrier()
    
    # Apply H to input qubits
    qc.h(range(n))
    
    # Measure input qubits
    qc.measure(range(n), range(n))
    
    return qc

# Test with n=4
sim = AerSimulator()
for oracle in ['constant_0', 'constant_1', 'balanced']:
    qc = deutsch_jozsa(4, oracle)
    counts = sim.run(qc, shots=100).result().get_counts()
    # Constant → only '0000'; Balanced → never '0000'
    verdict = 'constant' if '0000' in counts else 'balanced'
    print(f"{oracle}: {verdict} | counts={counts}")
```

---

## Simon's Algorithm (1994)

Simon's problem: the first problem with proven exponential quantum speedup.

### Problem:
```
Given: f:{0,1}ⁿ → {0,1}ⁿ satisfying: f(x) = f(x⊕s) for some hidden s
Promise: either s = 0ⁿ (one-to-one) or s ≠ 0ⁿ (two-to-one)
Find s.

Classical: Ω(2^(n/2)) queries (birthday paradox)
Quantum: O(n) queries
Exponential separation!
```

### Algorithm:
```
Repeat O(n) times:
  1. Prepare (1/√2ⁿ) Σₓ |x⟩|0⟩
  2. Apply Uf: (1/√2ⁿ) Σₓ |x⟩|f(x)⟩
  3. Measure second register → gets some value y = f(x₀)
     First register collapses to: (|x₀⟩ + |x₀⊕s⟩)/√2
  4. Apply H⊗n to first register:
     → (1/√2ⁿ) Σₖ (-1)^(k·x₀) [1 + (-1)^(k·s)] |k⟩
     Amplitude is zero unless k·s = 0 (mod 2)!
  5. Measure first register → get k with k·s = 0 (mod 2)

Collect n-1 linearly independent equations k·s = 0
Solve the system (classical Gaussian elimination) → find s
```

```python
from qiskit import QuantumCircuit
from qiskit_aer import AerSimulator
import numpy as np

def simon_circuit(n: int, s: str):
    """
    Simon's algorithm circuit
    s: the hidden period as binary string (length n)
    """
    qc = QuantumCircuit(2 * n, n)
    
    # Apply H to first n qubits
    qc.h(range(n))
    
    # Oracle: f(x) = f(x XOR s)
    # Simple construction: copy x to output, XOR with s when needed
    for i in range(n):
        qc.cx(i, i + n)
    
    # Apply XOR with s (for non-trivial period)
    if int(s, 2) != 0:
        # This creates a period — simplified oracle
        for i, bit in enumerate(s):
            if bit == '1':
                qc.cx(0, n + i)  # XOR with s when first bit is 1
    
    # Apply H to first n qubits
    qc.h(range(n))
    
    # Measure first n qubits
    qc.measure(range(n), range(n))
    
    return qc

def solve_simon_system(equations: list[str]) -> str:
    """Classical post-processing: solve k·s=0 system via Gaussian elimination."""
    n = len(equations[0])
    # Convert to matrix over GF(2)
    matrix = [[int(b) for b in eq] for eq in equations]
    
    # Gaussian elimination over GF(2)
    pivot_col = {}
    row = 0
    for col in range(n):
        pivot = None
        for r in range(row, len(matrix)):
            if matrix[r][col] == 1:
                pivot = r
                break
        if pivot is None:
            continue
        matrix[row], matrix[pivot] = matrix[pivot], matrix[row]
        pivot_col[col] = row
        for r in range(len(matrix)):
            if r != row and matrix[r][col] == 1:
                matrix[r] = [(matrix[r][i] ^ matrix[row][i]) for i in range(n)]
        row += 1
    
    # Back-substitute to find s
    s = [0] * n
    free_vars = [c for c in range(n) if c not in pivot_col]
    if free_vars:
        s[free_vars[0]] = 1  # set free variable to 1
        for col, row_idx in pivot_col.items():
            s[col] = matrix[row_idx][free_vars[0]]
    
    return ''.join(map(str, s))

# Run Simon's algorithm for s = '110'
n = 3
s = '110'
sim = AerSimulator()
equations = []
for _ in range(n * 3):  # collect more equations than needed
    qc = simon_circuit(n, s)
    result = sim.run(qc, shots=1).result()
    k = list(result.get_counts().keys())[0]
    if k != '0' * n:  # skip trivial equation
        equations.append(k)
    if len(equations) >= n - 1:
        break

print(f"Equations: {equations}")
recovered_s = solve_simon_system(equations)
print(f"Recovered s = {recovered_s}")  # Should be '110'
```

---

## Why Simon's Algorithm Matters

Simon's algorithm inspired Shor's algorithm. The structure is almost identical:

```
Simon's:
  Hidden period under XOR: f(x) = f(x⊕s)
  Quantum gives linear equations k·s = 0

Shor's:
  Hidden period under multiplication: f(x) = aˣ mod N
  QFT gives equations → period r
  Period r → factors of N

Simon's → Shor's: Replace XOR (cyclic group Z₂ⁿ) with multiplication (cyclic group Zᴺ*)
```

---

## Bernstein-Vazirani Algorithm

A special case of Deutsch-Jozsa — find a hidden string using 1 query:

```
Problem: f(x) = s·x (mod 2) for hidden string s
Classical: n queries (one bit of s at a time)
Quantum: 1 query

Circuit: Same as Deutsch-Jozsa!
After H⊗n, Oracle, H⊗n:
Measurement directly gives s

This shows quantum algorithms can learn hidden linear structure in 1 shot.
```

```python
def bernstein_vazirani(n: int, s: str):
    """Find hidden string s where f(x) = s·x mod 2"""
    qc = QuantumCircuit(n + 1, n)
    qc.x(n)
    qc.h(range(n + 1))
    
    # Oracle: f(x) = s·x → CNOT for each bit of s
    for i, bit in enumerate(s):
        if bit == '1':
            qc.cx(i, n)
    
    qc.h(range(n))
    qc.measure(range(n), range(n))
    return qc

sim = AerSimulator()
s = '10110'
qc = bernstein_vazirani(len(s), s)
counts = sim.run(qc, shots=1).result().get_counts()
print(f"Hidden s = {s}, recovered = {list(counts.keys())[0]}")  # Same!
```

---

## Summary: First Generation Quantum Algorithms

| Algorithm | Problem | Classical | Quantum | Speedup |
|---|---|---|---|---|
| Deutsch | 1-bit function constant/balanced | 2 queries | 1 query | 2× |
| Deutsch-Jozsa | n-bit function constant/balanced | 2^(n-1)+1 | 1 query | Exponential |
| Bernstein-Vazirani | Find hidden string s in s·x | n queries | 1 query | n× |
| Simon's | Find hidden XOR period | Ω(2^(n/2)) | O(n) | Exponential |

These algorithms demonstrate quantum advantage for artificial problems, paving the way for Shor's and Grover's algorithms on real-world problems.
