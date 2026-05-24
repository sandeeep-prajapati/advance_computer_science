# Chapter 13: Variational Quantum Algorithms (VQE & QAOA)

## The Variational Principle

Hybrid classical-quantum algorithms designed for NISQ hardware:

```
Idea: Use a parameterized quantum circuit (ansatz) with classical optimizer

quantum circuit (parameterized) → measure expectation value → classical optimizer
         ↑___________________________________________|

Theorem (Variational Principle):
  For Hamiltonian H with ground state energy E₀:
  ⟨ψ(θ)|H|ψ(θ)⟩ ≥ E₀   for any state |ψ(θ)⟩

So minimizing ⟨H⟩ over parameters θ gives upper bound on E₀.
```

---

## VQE — Variational Quantum Eigensolver

VQE finds the ground state energy of a Hamiltonian:

```
Application: Quantum chemistry
  Find lowest energy molecular configuration
  Simulate H₂, LiH, BeH₂, etc.

Classical methods:
  Exact diagonalization: O(4ⁿ) — scales exponentially
  Density Functional Theory (DFT): O(n³) — approximate
  CCSD(T): gold standard, O(n⁷) — too slow for large molecules

VQE: polynomial quantum resources (for good ansatz)
```

### VQE Algorithm:
```
1. Encode Hamiltonian H as weighted sum of Pauli operators:
   H = Σᵢ cᵢ Pᵢ    where Pᵢ ∈ {I,X,Y,Z}^⊗n

2. Choose ansatz |ψ(θ)⟩ = U(θ)|0⟩ⁿ

3. Quantum computer: measure ⟨Pᵢ⟩ for each Pauli term
   E(θ) = Σᵢ cᵢ ⟨ψ(θ)|Pᵢ|ψ(θ)⟩

4. Classical computer: update θ to minimize E(θ)
   θ* = argmin E(θ)

5. Return: E(θ*) ≈ E₀ (ground state energy)
```

---

## Implementing VQE

```python
from qiskit import QuantumCircuit
from qiskit.quantum_info import SparsePauliOp, Statevector
from qiskit_aer import AerSimulator
from scipy.optimize import minimize
import numpy as np

# H₂ Hamiltonian (simplified, 2-qubit, at fixed bond length)
H2_hamiltonian = SparsePauliOp.from_list([
    ('II', -1.0523732),
    ('IZ',  0.3979374),
    ('ZI', -0.3979374),
    ('ZZ', -0.0112801),
    ('XX',  0.1809312),
])

def ansatz(theta: list) -> QuantumCircuit:
    """Hardware-efficient ansatz for 2 qubits."""
    qc = QuantumCircuit(2)
    qc.ry(theta[0], 0)
    qc.ry(theta[1], 1)
    qc.cx(0, 1)
    qc.ry(theta[2], 0)
    qc.ry(theta[3], 1)
    return qc

def compute_energy(theta: np.ndarray) -> float:
    """Measure expectation value of H₂ Hamiltonian."""
    qc = ansatz(theta)
    sv = Statevector(qc)
    return sv.expectation_value(H2_hamiltonian).real

# Classical optimization
theta_init = np.random.uniform(0, 2*np.pi, 4)
result = minimize(compute_energy, theta_init, method='COBYLA',
                  options={'maxiter': 1000, 'rhobeg': 0.5})

print(f"VQE ground state energy: {result.fun:.6f} Ha")
print(f"Classical exact:         -1.857275 Ha")
print(f"Error: {abs(result.fun - (-1.857275)):.6f} Ha")

# Visualize optimal state
optimal_theta = result.x
qc = ansatz(optimal_theta)
print(f"Optimal ansatz circuit:")
print(qc.draw())
```

---

## Ansatz Design

The ansatz is critical — it determines what states VQE can explore:

### Hardware-Efficient Ansatz:
```python
def hardware_efficient_ansatz(n: int, depth: int, theta: np.ndarray) -> QuantumCircuit:
    qc = QuantumCircuit(n)
    idx = 0
    for _ in range(depth):
        for q in range(n):
            qc.ry(theta[idx], q); idx += 1
            qc.rz(theta[idx], q); idx += 1
        for q in range(n - 1):
            qc.cx(q, q + 1)  # linear entanglement
    return qc
# Pros: short depth, hardware-friendly
# Cons: may not reach ground state (expressibility vs trainability tradeoff)
```

### Chemistry-Inspired (UCCSD):
```python
# Unitary Coupled Cluster Single Double (UCCSD)
# Based on quantum chemistry theory
from qiskit_nature.second_q.circuit.library import UCCSD
# Better for chemistry but deeper circuit
```

---

## Barren Plateaus Problem

A major challenge for VQA training:

```
For deep circuits with many qubits:
  Gradients of cost function become exponentially small
  ∂E/∂θᵢ → 0 exponentially in n

Cause: random initialization → state is uniformly distributed
       → cost function landscape is flat almost everywhere

Solutions:
  1. Local cost functions (instead of global)
  2. Layer-by-layer training
  3. Problem-specific initialization
  4. Warm starting from classical solution
  5. Rotosolve / parameter shift rules
```

### Parameter shift rule (exact gradient):
```python
def parameter_shift_gradient(theta: np.ndarray, i: int, shift: float = np.pi/2) -> float:
    """Compute exact gradient of E w.r.t. theta[i] using parameter shift."""
    theta_plus = theta.copy(); theta_plus[i] += shift
    theta_minus = theta.copy(); theta_minus[i] -= shift
    return (compute_energy(theta_plus) - compute_energy(theta_minus)) / (2 * np.sin(shift))
```

---

## QAOA — Quantum Approximate Optimization Algorithm

QAOA finds approximate solutions to combinatorial optimization problems.

### Target: MaxCut Problem
```
PROBLEM: MaxCut
  Given: Graph G = (V, E)
  Find: Partition vertices into two sets S, S̄
  Maximize: Number of edges between S and S̄

Classical: NP-hard for general graphs
Best approximation: Goemans-Williamson (0.878 approximation)
QAOA: Parameterizable approximation, performance improves with circuit depth p
```

### QAOA Structure:
```
Encode MaxCut as Ising Hamiltonian:
  H_C = Σ_{(i,j)∈E} (1 - ZᵢZⱼ)/2

QAOA state: |β,γ⟩ = e^(-iβ_p H_B) e^(-iγ_p H_C) ... e^(-iβ₁ H_B) e^(-iγ₁ H_C) |+⟩^n

H_B = Σᵢ Xᵢ  (mixing Hamiltonian — transverse field)
H_C = cost Hamiltonian

Optimize angles β, γ to maximize ⟨H_C⟩
```

```python
from qiskit import QuantumCircuit
from qiskit.quantum_info import SparsePauliOp, Statevector
import numpy as np
from scipy.optimize import minimize

# MaxCut on a simple 4-node cycle graph
# Edges: (0,1), (1,2), (2,3), (3,0)

edges = [(0, 1), (1, 2), (2, 3), (3, 0)]
n = 4

def cost_hamiltonian(n, edges) -> SparsePauliOp:
    """H_C = Σ (I - ZᵢZⱼ)/2 for each edge."""
    terms = []
    for (i, j) in edges:
        # ZᵢZⱼ term
        pauli_str = ['I'] * n
        pauli_str[i] = 'Z'
        pauli_str[j] = 'Z'
        terms.append((''.join(reversed(pauli_str)), -0.5))
        terms.append(('I' * n, 0.5))
    return SparsePauliOp.from_list(terms)

def qaoa_circuit(n: int, edges: list, gamma: np.ndarray, beta: np.ndarray) -> QuantumCircuit:
    p = len(gamma)  # QAOA depth
    qc = QuantumCircuit(n)
    
    # Initialize |+⟩^n
    qc.h(range(n))
    
    for layer in range(p):
        # Phase separation (cost Hamiltonian)
        for (i, j) in edges:
            qc.cx(i, j)
            qc.rz(2 * gamma[layer], j)
            qc.cx(i, j)
        
        # Mixing (X Hamiltonian)
        for q in range(n):
            qc.rx(2 * beta[layer], q)
    
    return qc

def qaoa_cost(params: np.ndarray) -> float:
    p = len(params) // 2
    gamma = params[:p]
    beta = params[p:]
    
    qc = qaoa_circuit(n, edges, gamma, beta)
    sv = Statevector(qc)
    H_C = cost_hamiltonian(n, edges)
    return -sv.expectation_value(H_C).real  # negate for minimization

# Optimize QAOA with p=1
p = 1
params0 = np.random.uniform(0, np.pi, 2 * p)
result = minimize(qaoa_cost, params0, method='COBYLA')

gamma_opt = result.x[:p]
beta_opt = result.x[p:]
qc_opt = qaoa_circuit(n, edges, gamma_opt, beta_opt)
sv = Statevector(qc_opt)

probs = sv.probabilities_dict()
print("Top solutions (MaxCut on 4-cycle):")
for state, prob in sorted(probs.items(), key=lambda x: -x[1])[:4]:
    cut = sum(1 for (i,j) in edges if state[n-1-i] != state[n-1-j])
    print(f"  |{state}⟩: prob={prob:.3f}, cut={cut}")
# Best cut for 4-cycle is 4 (alternate 0101 or 1010)
```

---

## QAOA Performance

```
QAOA depth p:
  p=1: may outperform random, but not optimal for hard instances
  p→∞: approaches optimal solution (adiabatic theorem)

MaxCut approximation ratio:
  p=1 QAOA on 3-regular graphs: ≥ 0.6924
  Best classical (Goemans-Williamson): 0.878
  QAOA with large p: potentially better for specific problems

Current NISQ limitation:
  Decoherence limits useful circuit depth
  p=1-5 is practical on current hardware
```

---

## Other VQA Variants

| Algorithm | Problem | Ansatz |
|---|---|---|
| VQE | Ground state energy | Chemistry-inspired |
| QAOA | Combinatorial optimization | Problem-specific phases |
| VQLS | Linear systems Ax=b | Generic |
| QNN | Machine learning | Parameterized circuit |
| QSVM | Classification | Quantum kernel |
| ADAPT-VQE | Chemistry (adaptive) | Grows greedily |

---

## Classical Simulation of VQAs

```python
# PennyLane: best library for VQA with autodiff
import pennylane as qml
import numpy as np

dev = qml.device('default.qubit', wires=2)

@qml.qnode(dev)
def circuit(theta):
    qml.RY(theta[0], wires=0)
    qml.RY(theta[1], wires=1)
    qml.CNOT(wires=[0, 1])
    qml.RY(theta[2], wires=0)
    return qml.expval(qml.PauliZ(0) @ qml.PauliZ(1))

# Automatic differentiation — exact gradients!
theta = np.array([0.1, 0.2, 0.3], requires_grad=True)
grad_fn = qml.grad(circuit)
print(f"Energy: {circuit(theta):.4f}")
print(f"Gradient: {grad_fn(theta)}")

# Gradient descent
opt = qml.GradientDescentOptimizer(stepsize=0.1)
for i in range(100):
    theta, cost = opt.step_and_cost(circuit, theta)
    if i % 10 == 0:
        print(f"Step {i}: cost = {cost:.4f}")
```
