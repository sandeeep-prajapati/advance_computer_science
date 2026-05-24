# Chapter 17: Cirq and PennyLane

## Cirq (Google)

Cirq is Google's quantum computing framework, designed for NISQ devices:

```bash
pip install cirq cirq-google
```

### Core concepts:
```python
import cirq
import numpy as np

# Qubits
q0, q1, q2 = cirq.LineQubit.range(3)       # linear array
a = cirq.GridQubit(0, 0); b = cirq.GridQubit(0, 1)  # grid qubits (Google hardware)
q = cirq.NamedQubit('alice')                # named qubit

# Gates and moments
circuit = cirq.Circuit([
    cirq.H(q0),
    cirq.CNOT(q0, q1),
    cirq.measure(q0, q1, key='result')
])
print(circuit)

# Parametrized gates
theta = cirq.Parameter('theta')
qc = cirq.Circuit(cirq.Ry(rads=theta)(q0))
```

### Bell state in Cirq:
```python
import cirq

q0, q1 = cirq.LineQubit.range(2)
circuit = cirq.Circuit([
    cirq.H(q0),
    cirq.CNOT(q0, q1),
    cirq.measure(q0, q1, key='bell'),
])
print(circuit)

# Simulate
simulator = cirq.Simulator()
result = simulator.run(circuit, repetitions=1000)
print(result.histogram(key='bell'))
# {0: ~500 (=|00⟩), 3: ~500 (=|11⟩)}
```

### Statevector simulation:
```python
simulator = cirq.Simulator()
result = simulator.simulate(circuit)
print(result.final_state_vector)  # complex array

# Density matrix simulation
dm_sim = cirq.DensityMatrixSimulator()
dm_result = dm_sim.simulate(circuit)
print(dm_result.final_density_matrix)
```

### Cirq noise models:
```python
# Depolarizing noise
noise = cirq.depolarize(p=0.01)
noisy_circuit = circuit.with_noise(noise)

result = cirq.DensityMatrixSimulator(noise=noise).run(circuit, repetitions=1000)

# Channel as gate
noisy_q = cirq.Circuit(
    cirq.H(q0),
    cirq.depolarize(0.01)(q0),  # apply noise channel
    cirq.measure(q0, key='m')
)
```

### Grover's algorithm in Cirq:
```python
import cirq
import numpy as np

def oracle_circuit(qubits, target_state: list):
    """Phase oracle that marks target_state with a phase flip."""
    n = len(qubits)
    circuit = cirq.Circuit()
    
    # Flip qubits where target is 0
    flips = [qubits[i] for i, b in enumerate(target_state) if b == 0]
    if flips:
        circuit.append(cirq.X.on_each(*flips))
    
    # Multi-controlled Z (phase flip for |11...1⟩)
    circuit.append(cirq.Z.controlled(n-1).on(*qubits))
    
    # Undo flips
    if flips:
        circuit.append(cirq.X.on_each(*flips))
    
    return circuit

def diffuser(qubits):
    """Grover diffusion operator."""
    n = len(qubits)
    circuit = cirq.Circuit()
    circuit.append(cirq.H.on_each(*qubits))
    circuit.append(cirq.X.on_each(*qubits))
    circuit.append(cirq.Z.controlled(n-1).on(*qubits))
    circuit.append(cirq.X.on_each(*qubits))
    circuit.append(cirq.H.on_each(*qubits))
    return circuit

def grover(n, target_state, num_iterations=None):
    qubits = cirq.LineQubit.range(n)
    if num_iterations is None:
        num_iterations = int(np.pi / 4 * np.sqrt(2**n))
    
    circuit = cirq.Circuit()
    circuit.append(cirq.H.on_each(*qubits))  # superposition
    
    for _ in range(num_iterations):
        circuit += oracle_circuit(qubits, target_state)
        circuit += diffuser(qubits)
    
    circuit.append(cirq.measure(*qubits, key='result'))
    return circuit

# Search for |101⟩
qc = grover(3, [1, 0, 1])
result = cirq.Simulator().run(qc, repetitions=1000)
print(result.histogram(key='result'))
```

---

## PennyLane

PennyLane is the best framework for quantum machine learning — supports automatic differentiation and multiple backends:

```bash
pip install pennylane pennylane-qiskit pennylane-cirq torch
```

### Core concepts:
```python
import pennylane as qml
import numpy as np

# Create a device (simulator or hardware)
dev = qml.device('default.qubit', wires=3)      # exact statevector
dev_noise = qml.device('qiskit.aer', wires=3)   # use Qiskit's Aer

# QNode: quantum function decorator
@qml.qnode(dev)
def circuit(params):
    qml.RY(params[0], wires=0)
    qml.RY(params[1], wires=1)
    qml.CNOT(wires=[0, 1])
    return qml.expval(qml.PauliZ(0))  # return expectation value

# Call like a function
params = np.array([0.5, 1.2])
print(circuit(params))  # float
```

### Automatic differentiation:
```python
import pennylane as qml
import numpy as np

dev = qml.device('default.qubit', wires=2)

@qml.qnode(dev)
def circuit(theta):
    qml.RY(theta, wires=0)
    qml.CNOT(wires=[0, 1])
    return qml.expval(qml.PauliZ(0))

# Gradient via parameter shift rule (exact)
theta = np.array(0.5, requires_grad=True)
grad = qml.grad(circuit)(theta)
print(f"dE/dtheta = {grad}")

# Interface with PyTorch
import torch
dev_torch = qml.device('default.qubit', wires=2)

@qml.qnode(dev_torch, interface='torch')
def torch_circuit(theta):
    qml.RY(theta, wires=0)
    return qml.expval(qml.PauliZ(0))

theta_torch = torch.tensor(0.5, requires_grad=True)
loss = torch_circuit(theta_torch)
loss.backward()
print(f"Gradient: {theta_torch.grad}")
```

### VQE with PennyLane:
```python
import pennylane as qml
from pennylane import numpy as np

dev = qml.device('default.qubit', wires=2)

# H₂ Hamiltonian
coeffs = [-1.0523732, 0.3979374, -0.3979374, -0.0112801, 0.1809312]
obs = [
    qml.Identity(0) @ qml.Identity(1),
    qml.PauliZ(0) @ qml.Identity(1),
    qml.Identity(0) @ qml.PauliZ(1),
    qml.PauliZ(0) @ qml.PauliZ(1),
    qml.PauliX(0) @ qml.PauliX(1),
]
H = qml.Hamiltonian(coeffs, obs)

@qml.qnode(dev)
def vqe_circuit(theta):
    qml.RY(theta[0], wires=0)
    qml.RY(theta[1], wires=1)
    qml.CNOT(wires=[0, 1])
    qml.RY(theta[2], wires=0)
    qml.RY(theta[3], wires=1)
    return qml.expval(H)

# Optimize
theta = np.random.uniform(0, 2*np.pi, 4, requires_grad=True)
opt = qml.AdamOptimizer(stepsize=0.1)
for i in range(200):
    theta, cost = opt.step_and_cost(vqe_circuit, theta)
    if i % 20 == 0:
        print(f"Step {i}: E = {cost:.6f} Ha")

print(f"\nFinal energy: {vqe_circuit(theta):.6f} Ha")
print(f"Exact (FCI):  -1.857275 Ha")
```

### QAOA with PennyLane:
```python
import pennylane as qml
from pennylane import numpy as np
import networkx as nx

# Create graph
G = nx.cycle_graph(4)  # 4-cycle
n_wires = G.number_of_nodes()

dev = qml.device('default.qubit', wires=n_wires)

def cost_layer(gamma):
    for edge in G.edges():
        qml.CNOT(wires=list(edge))
        qml.RZ(2 * gamma, wires=edge[1])
        qml.CNOT(wires=list(edge))

def mixer_layer(beta):
    for q in range(n_wires):
        qml.RX(2 * beta, wires=q)

@qml.qnode(dev)
def qaoa_circuit(params, depth=1):
    gammas = params[:depth]
    betas = params[depth:]
    
    # Initialize
    for q in range(n_wires):
        qml.Hadamard(wires=q)
    
    # QAOA layers
    for d in range(depth):
        cost_layer(gammas[d])
        mixer_layer(betas[d])
    
    # Cost function: negative MaxCut
    cost_obs = []
    for edge in G.edges():
        cost_obs.append(qml.PauliZ(edge[0]) @ qml.PauliZ(edge[1]))
    
    return qml.expval(sum(cost_obs))

# Optimize with p=1
depth = 1
params = np.random.uniform(0, np.pi, 2*depth, requires_grad=True)
opt = qml.GradientDescentOptimizer(stepsize=0.5)

for i in range(100):
    params, cost = opt.step_and_cost(lambda p: qaoa_circuit(p, depth), params)

print(f"MaxCut approximation: {-qaoa_circuit(params, depth):.3f}")
print(f"(True MaxCut for 4-cycle = 4)")
```

---

## PennyLane Quantum ML Templates

```python
import pennylane as qml
import numpy as np

dev = qml.device('default.qubit', wires=4)

# Angle embedding: encode classical data as rotation angles
@qml.qnode(dev)
def angle_embedding_circuit(x, weights):
    qml.AngleEmbedding(x, wires=range(4), rotation='Y')
    qml.BasicEntanglerLayers(weights, wires=range(4))
    return [qml.expval(qml.PauliZ(i)) for i in range(4)]

x = np.array([0.1, 0.5, 0.3, 0.9])  # input data
weight_shape = qml.BasicEntanglerLayers.shape(n_layers=2, n_wires=4)
weights = np.random.uniform(0, 2*np.pi, weight_shape)
print(angle_embedding_circuit(x, weights))

# Amplitude embedding: encode 2^n data points in n-qubit state
@qml.qnode(dev)
def amplitude_embedding_circuit(x):
    qml.AmplitudeEmbedding(x, wires=range(4), normalize=True)
    return qml.expval(qml.PauliZ(0))

x_data = np.array([0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8,
                   0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8])  # 2^4 = 16 values
print(amplitude_embedding_circuit(x_data))
```

---

## Comparing Frameworks

| Feature | Qiskit | Cirq | PennyLane |
|---|---|---|---|
| Hardware access | IBM Q | Google | Multiple backends |
| Best for | IBM hardware, general QC | Google hardware, NISQ | QML, autodiff |
| Circuit style | Object-oriented | Gate-on-qubit | Functional/decorator |
| Autodiff | No (manual) | No | Yes (parameter shift, adjoint) |
| ML integration | Limited | Limited | PyTorch, TF, JAX, NumPy |
| Noise models | Excellent (Aer) | Good | Via Qiskit/Cirq backend |
| Community | Largest | Large | Growing (ML focus) |
| Visualization | Good | Good | Good |
| Learning curve | Medium | Medium | Lower for ML devs |

### When to use which:
```
Qiskit: IBM quantum hardware, comprehensive QC library, production
Cirq:   Google Sycamore hardware, NISQ algorithm research
PennyLane: Quantum ML, hybrid optimization, autodiff gradients

Can mix: PennyLane can use Qiskit/Cirq as backend!
qml.device('qiskit.ibm', ...)  # PennyLane → IBM hardware via Qiskit
```

---

## Problem-Solving: Choosing a Framework

```python
# Decision guide
def choose_framework(task: str) -> str:
    if "IBM hardware" in task:
        return "Qiskit"
    elif "Google hardware" in task:
        return "Cirq"
    elif "machine learning" in task or "gradient" in task or "VQE" in task:
        return "PennyLane"
    elif "research" in task or "noise" in task:
        return "Qiskit (Aer noise models)"
    elif "production" in task:
        return "Qiskit (largest ecosystem)"
    else:
        return "PennyLane (most versatile)"
```
