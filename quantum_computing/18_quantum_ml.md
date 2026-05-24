# Chapter 18: Quantum Machine Learning

## Overview of QML

Quantum Machine Learning (QML) applies quantum computing to machine learning:

```
Three paradigms:
1. CC → Classical data, Classical ML (traditional)
2. QC → Quantum data, Classical ML (quantum sensor → classical analysis)
3. CQ → Classical data, Quantum ML (most common QML research)
4. QQ → Quantum data, Quantum ML (quantum network learning)

Current focus: CQ — use quantum circuits to process classical data

Central question: Does quantum ML offer genuine speedup?
Honest answer: Unclear. Active research. Many claims are overstated.
```

---

## Variational Quantum Classifier

The quantum analog of a classical neural network:

```
Input: classical data x ∈ ℝⁿ
Encoding: U(x)|0⟩ → quantum state (feature map)
Model: V(θ)|ψ(x)⟩ → parameterized circuit (ansatz)
Output: ⟨O⟩ = ⟨ψ(x,θ)|O|ψ(x,θ)⟩ → real value → binary/multiclass label
Training: minimize loss L(θ) w.r.t. θ using gradient descent
```

```python
import pennylane as qml
import numpy as np
from sklearn.datasets import make_classification
from sklearn.preprocessing import normalize

# Create quantum device
dev = qml.device('default.qubit', wires=4)
n_qubits = 4

# Feature map: encode data as angles
def feature_map(x):
    qml.AngleEmbedding(x, wires=range(n_qubits))
    for i in range(n_qubits - 1):
        qml.CNOT(wires=[i, i+1])

# Variational ansatz
def ansatz(weights):
    qml.BasicEntanglerLayers(weights, wires=range(n_qubits))

@qml.qnode(dev)
def quantum_classifier(x, weights):
    feature_map(x)
    ansatz(weights)
    return qml.expval(qml.PauliZ(0))  # binary: >0 → class 1, <0 → class 0

# Training
def loss(weights, X, y):
    preds = np.array([quantum_classifier(x, weights) for x in X])
    return np.mean((preds - y)**2)  # MSE loss

# Dataset
X, y = make_classification(n_samples=100, n_features=4, random_state=42)
X = normalize(X) * np.pi  # scale to [-π, π]
y = 2 * y - 1  # convert {0,1} to {-1,+1}

# Initialize weights
weight_shape = qml.BasicEntanglerLayers.shape(n_layers=2, n_wires=n_qubits)
weights = np.random.uniform(-np.pi, np.pi, weight_shape, requires_grad=True)

# Optimize
opt = qml.AdamOptimizer(stepsize=0.1)
for i in range(50):
    weights, l = opt.step_and_cost(lambda w: loss(w, X[:20], y[:20]), weights)
    if i % 10 == 0:
        preds = np.sign([quantum_classifier(x, weights) for x in X])
        acc = np.mean(preds == y)
        print(f"Step {i}: loss={l:.4f}, acc={acc:.3f}")
```

---

## Quantum Kernel Methods

A mathematically rigorous approach to QML:

```
Classical kernel: k(x, x') = φ(x)·φ(x')  (implicit feature map)
Quantum kernel: k_Q(x, x') = |⟨φ_Q(x)|φ_Q(x')⟩|²

where |φ_Q(x)⟩ = U(x)|0⟩ is a quantum feature state

Advantage: Quantum feature spaces may be exponentially larger than classical.
           Hard to compute classically → potential quantum advantage.

Use with classical SVM:
  Compute kernel matrix K[i,j] = k_Q(xᵢ, xⱼ) on quantum computer
  Train SVM with this kernel on classical computer
```

```python
import pennylane as qml
import numpy as np
from sklearn.svm import SVC
from sklearn.metrics import accuracy_score

dev = qml.device('default.qubit', wires=4)

@qml.qnode(dev)
def kernel_circuit(x1, x2):
    """Compute |⟨φ(x1)|φ(x2)⟩|² quantum kernel."""
    # Encode x1
    qml.AngleEmbedding(x1, wires=range(4))
    # Encode x2 inverse
    qml.AngleEmbedding(x2, wires=range(4)).inv()
    return qml.probs(wires=range(4))

def quantum_kernel(X1, X2):
    """Compute quantum kernel matrix."""
    K = np.zeros((len(X1), len(X2)))
    for i, x1 in enumerate(X1):
        for j, x2 in enumerate(X2):
            # |⟨0|U†(x2)U(x1)|0⟩|² = probability of measuring |0...0⟩
            probs = kernel_circuit(x1, x2)
            K[i, j] = probs[0]  # |0...0⟩ probability
    return K

# Generate data
from sklearn.datasets import make_classification
X, y = make_classification(n_samples=60, n_features=4, random_state=42)
X_train, X_test = X[:40], X[40:]
y_train, y_test = y[:40], y[40:]

# Compute kernel matrix
print("Computing quantum kernel matrix...")
K_train = quantum_kernel(X_train, X_train)
K_test = quantum_kernel(X_test, X_train)

# Train SVM with quantum kernel
svm = SVC(kernel='precomputed')
svm.fit(K_train, y_train)
y_pred = svm.predict(K_test)
print(f"Quantum SVM accuracy: {accuracy_score(y_test, y_pred):.3f}")
```

---

## Quantum Neural Networks (QNN)

```python
import pennylane as qml
import torch
import numpy as np

# PennyLane + PyTorch integration
n_qubits = 4
n_layers = 3
dev = qml.device('default.qubit', wires=n_qubits)

@qml.qnode(dev, interface='torch')
def quantum_neural_net(inputs, weights):
    """Quantum neural network as a PyTorch layer."""
    # Angle embedding
    qml.AngleEmbedding(inputs, wires=range(n_qubits))
    
    # Parameterized layers
    for layer_weights in weights:
        for i in range(n_qubits):
            qml.RY(layer_weights[i, 0], wires=i)
            qml.RZ(layer_weights[i, 1], wires=i)
        for i in range(n_qubits - 1):
            qml.CNOT(wires=[i, i+1])
    
    return [qml.expval(qml.PauliZ(i)) for i in range(n_qubits)]

# Define as PyTorch module
class QuantumLayer(torch.nn.Module):
    def __init__(self, n_qubits, n_layers):
        super().__init__()
        weight_shapes = (n_layers, n_qubits, 2)
        self.weights = torch.nn.Parameter(
            torch.randn(*weight_shapes)
        )
    
    def forward(self, x):
        return torch.stack([
            torch.tensor(quantum_neural_net(xi, self.weights), dtype=torch.float32)
            for xi in x
        ])

class HybridClassifier(torch.nn.Module):
    def __init__(self):
        super().__init__()
        self.quantum = QuantumLayer(n_qubits, n_layers)
        self.classical = torch.nn.Linear(n_qubits, 2)  # 2-class output
    
    def forward(self, x):
        q_out = self.quantum(x)
        return self.classical(q_out)

# Train
model = HybridClassifier()
optimizer = torch.optim.Adam(model.parameters(), lr=0.01)
criterion = torch.nn.CrossEntropyLoss()

# Example training loop
X = torch.randn(50, n_qubits)  # 50 samples, 4 features
y = torch.randint(0, 2, (50,))  # binary labels

for epoch in range(10):
    optimizer.zero_grad()
    out = model(X[:10])  # small batch due to quantum simulation cost
    loss = criterion(out, y[:10])
    loss.backward()
    optimizer.step()
    print(f"Epoch {epoch}: loss = {loss.item():.4f}")
```

---

## Quantum Transfer Learning

Use a pre-trained classical model + quantum circuit:

```python
import torch
import torchvision
import pennylane as qml

# Classical pre-trained feature extractor
classical_model = torchvision.models.resnet18(pretrained=True)
classical_model.fc = torch.nn.Identity()  # remove final layer → 512-dim features

n_qubits = 4
dev = qml.device('default.qubit', wires=n_qubits)

@qml.qnode(dev, interface='torch')
def quantum_layer(inputs, weights):
    qml.AngleEmbedding(inputs[:n_qubits], wires=range(n_qubits))
    qml.BasicEntanglerLayers(weights, wires=range(n_qubits))
    return [qml.expval(qml.PauliZ(i)) for i in range(n_qubits)]

class QuantumTransferModel(torch.nn.Module):
    def __init__(self):
        super().__init__()
        self.classical = classical_model
        self.reducer = torch.nn.Linear(512, n_qubits)
        self.q_weights = torch.nn.Parameter(
            torch.randn(2, n_qubits)
        )
        self.classifier = torch.nn.Linear(n_qubits, 2)
    
    def forward(self, x):
        features = self.classical(x)          # 512-dim
        reduced = torch.tanh(self.reducer(features)) * np.pi  # n_qubits-dim, scaled
        q_out = torch.stack([
            torch.tensor(quantum_layer(r, self.q_weights))
            for r in reduced
        ])
        return self.classifier(q_out)
```

---

## Quantum Generative Models

Quantum circuit born machines (QCBMs) as generative models:

```python
import pennylane as qml
import numpy as np

dev = qml.device('default.qubit', wires=4)

@qml.qnode(dev)
def generative_circuit(params):
    """Parameterized circuit that generates probability distribution."""
    for layer in range(3):
        for q in range(4):
            qml.RY(params[layer, q, 0], wires=q)
            qml.RZ(params[layer, q, 1], wires=q)
        for q in range(3):
            qml.CNOT(wires=[q, q+1])
    
    return qml.probs(wires=range(4))  # probability distribution over 16 states

# Target distribution (e.g., from data)
target_dist = np.array([0.1, 0.2, 0.05, 0.15, 0.0, 0.3, 0.1, 0.05,
                        0.02, 0.01, 0.01, 0.01, 0.0, 0.0, 0.0, 0.0])
target_dist /= target_dist.sum()  # normalize

# Train to match distribution using KL divergence
params = np.random.uniform(0, np.pi, (3, 4, 2), requires_grad=True)

def kl_divergence(params):
    gen_dist = generative_circuit(params)
    return np.sum(target_dist * np.log(target_dist / (gen_dist + 1e-10)))

opt = qml.GradientDescentOptimizer(0.1)
for step in range(50):
    params, cost = opt.step_and_cost(kl_divergence, params)
    if step % 10 == 0:
        print(f"Step {step}: KL divergence = {cost:.4f}")
```

---

## Quantum Convolutional Neural Networks (QCNN)

```python
import pennylane as qml
import numpy as np

dev = qml.device('default.qubit', wires=8)

def conv_layer(params, wires):
    """Quantum convolutional layer: 2-qubit gate applied to pairs."""
    for i, (w1, w2) in enumerate(zip(wires[::2], wires[1::2])):
        qml.RY(params[i, 0], wires=w1)
        qml.RY(params[i, 1], wires=w2)
        qml.CNOT(wires=[w1, w2])

def pool_layer(wires):
    """Quantum pooling: measure half the qubits, condition on result."""
    for w1, w2 in zip(wires[::2], wires[1::2]):
        # Simplified pooling: CNOT and discard first qubit
        qml.CNOT(wires=[w1, w2])

@qml.qnode(dev)
def qcnn(x, params):
    # Encode 8-dimensional data
    qml.AngleEmbedding(x, wires=range(8))
    
    # Conv layer 1: acts on all 8 qubits in pairs
    conv_layer(params[0], wires=range(8))
    
    # Pool: 8 → 4 qubits
    pool_layer(wires=range(8))
    
    # Conv layer 2: acts on remaining 4 qubits
    conv_layer(params[1], wires=range(1, 8, 2))
    
    # Pool: 4 → 2 qubits
    pool_layer(wires=range(1, 8, 2))
    
    # Final measurement
    return qml.expval(qml.PauliZ(3))
```

---

## QML Reality Check

```
Claims about QML speedup — what's real?

POTENTIALLY REAL:
  Quantum kernel methods: If data has quantum structure, kernel
  computation is classically hard → genuine speedup.
  
  Quantum simulation for chemistry/materials: Yes, real speedup.

PROBABLY NOT REAL (yet):
  Quantum speedup for classical data: Most QML algorithms have
  "dequantization" — equally fast classical algorithms exist.
  
  Liu et al. (2021): Classification speedup claims are based on
  assumptions that hide classical algorithms of similar complexity.

THE HONEST PICTURE:
  QML is a nascent field. Most results are proofs-of-concept.
  True quantum advantage for ML requires:
    1. Fault-tolerant quantum computers
    2. Problems with genuine quantum structure
    3. Better classical lower bounds

BEST BET: Quantum simulation for drug discovery.
          Processing data that is itself quantum (quantum sensors, quantum networks).
```

---

## Practical QML Workflow

```
1. Classical preprocessing:
   - Feature engineering
   - Dimensionality reduction (must fit in n qubits = 2ⁿ amplitudes)
   - Normalization

2. Feature encoding choice:
   - Angle encoding: n features → n qubits (simpler, more practical)
   - Amplitude encoding: 2ⁿ features → n qubits (expressive but hard to prepare)
   - Basis encoding: binary data → computational basis

3. Ansatz choice:
   - Hardware-efficient: short depth, hardware-friendly
   - Problem-inspired: better expressibility, harder to run

4. Training:
   - Parameter shift rule: exact gradients
   - SPSA: gradient-free, works with shot noise
   - Watch for barren plateaus

5. Evaluation:
   - Compare to classical baselines (SVM, neural nets)
   - Honest about shot noise and circuit depth requirements
```
