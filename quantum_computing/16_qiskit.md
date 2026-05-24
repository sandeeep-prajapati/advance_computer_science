# Chapter 16: Quantum Programming with Qiskit

## Qiskit Ecosystem

Qiskit (IBM's quantum SDK) is the most widely used quantum programming framework:

```
qiskit              → Core: circuits, gates, transpilation
qiskit-aer          → Simulators (statevector, QASM, noise)
qiskit-ibm-runtime  → Access to IBM quantum hardware
qiskit-nature       → Quantum chemistry
qiskit-ml           → Quantum machine learning
qiskit-optimization → Combinatorial optimization
```

```bash
pip install qiskit qiskit-aer qiskit-ibm-runtime
pip install matplotlib pylatexenc  # visualization
```

---

## Building Circuits

```python
from qiskit import QuantumCircuit, QuantumRegister, ClassicalRegister

# Method 1: Simple
qc = QuantumCircuit(3, 2)  # 3 qubits, 2 classical bits

# Method 2: Named registers
qr = QuantumRegister(2, 'data')
ar = QuantumRegister(1, 'ancilla')
cr = ClassicalRegister(2, 'output')
qc = QuantumCircuit(qr, ar, cr)

# All gates
qc.h(qr[0])
qc.x(qr[1])
qc.y(qr[0])
qc.z(qr[1])
qc.s(qr[0])       # S gate = Z^(1/2)
qc.sdg(qr[0])     # S† gate
qc.t(qr[0])       # T gate = Z^(1/4)
qc.tdg(qr[0])     # T† gate
qc.rx(0.5, qr[0])
qc.ry(1.2, qr[1])
qc.rz(0.3, qr[0])
qc.p(0.7, qr[0])  # phase gate
qc.u(0.5, 1.2, 0.3, qr[0])  # general U gate

# Multi-qubit gates
qc.cx(qr[0], qr[1])      # CNOT
qc.cz(qr[0], qr[1])      # CZ
qc.cy(qr[0], qr[1])      # CY
qc.swap(qr[0], qr[1])    # SWAP
qc.ccx(qr[0], qr[1], ar[0])  # Toffoli
qc.cswap(ar[0], qr[0], qr[1])  # Fredkin

# Controlled versions
qc.ch(qr[0], qr[1])      # Controlled H
qc.crx(0.5, qr[0], qr[1])

# Measurement
qc.measure(qr, cr)         # all at once
qc.measure(qr[0], cr[0])  # individual

# Reset
qc.reset(qr[0])   # reset to |0⟩ (mid-circuit)

# Barrier (visual separator, no effect on computation)
qc.barrier()

# Draw
print(qc.draw())           # text
qc.draw('mpl')             # matplotlib figure
```

---

## Simulating Circuits

```python
from qiskit_aer import AerSimulator, StatevectorSimulator
from qiskit.quantum_info import Statevector

# 1. Statevector simulation (exact, exponential memory)
qc = QuantumCircuit(3)
qc.h(0)
qc.cx(0, 1)
qc.cx(0, 2)  # GHZ state

sv = Statevector(qc)
print(sv)                           # full state vector
print(sv.probabilities_dict())      # {'000': 0.5, '111': 0.5}
sv.draw('bloch')                    # Bloch sphere (single qubits only)
sv.draw('qsphere')                  # Q-sphere visualization

# 2. Sampling simulation (QASM)
qc_m = qc.copy()
qc_m.measure_all()
sim = AerSimulator()
job = sim.run(qc_m, shots=8192)
counts = job.result().get_counts()
print(counts)  # {'000': 4103, '111': 4089}

# 3. Density matrix simulation
from qiskit_aer import DensityMatrixSimulator
from qiskit.quantum_info import DensityMatrix
dm = DensityMatrix(qc)
print(dm.purity())  # 1.0 = pure state

# 4. Noisy simulation
from qiskit_aer.noise import NoiseModel, depolarizing_error
noise_model = NoiseModel()
noise_model.add_all_qubit_quantum_error(depolarizing_error(0.01, 1), ['h', 'rx'])
noise_model.add_all_qubit_quantum_error(depolarizing_error(0.05, 2), ['cx'])
noisy_sim = AerSimulator(noise_model=noise_model)
noisy_counts = noisy_sim.run(qc_m, shots=8192).result().get_counts()
print(noisy_counts)  # Imperfect, may see '001', '010', etc.
```

---

## Transpilation

Before running on hardware, circuits are transpiled to the target's native gate set and connectivity:

```python
from qiskit.compiler import transpile
from qiskit_aer import AerSimulator

# Simulate a backend with limited connectivity
from qiskit.providers.fake_provider import FakeManilaV2
backend = FakeManilaV2()  # simulated 5-qubit IBM backend

qc = QuantumCircuit(3)
qc.h(0)
qc.cx(0, 2)  # non-adjacent qubits
qc.cx(1, 2)

# Transpile for backend
t_qc = transpile(qc, backend=backend, optimization_level=3)
print(f"Original: {qc.count_ops()}")
print(f"Transpiled: {t_qc.count_ops()}")
print(f"Added SWAP gates: {t_qc.count_ops().get('swap', 0)}")
```

---

## Running on IBM Quantum Hardware

```python
from qiskit_ibm_runtime import QiskitRuntimeService, SamplerV2 as Sampler

# Setup (one-time): save credentials
QiskitRuntimeService.save_account(
    channel="ibm_quantum",
    token="YOUR_IBM_API_TOKEN",  # get from quantum.ibm.com
    set_as_default=True,
    overwrite=True
)

# Load account
service = QiskitRuntimeService()

# List available backends
backends = service.backends()
for b in backends:
    print(f"{b.name}: {b.num_qubits} qubits, pending: {b.status().pending_jobs}")

# Select a backend
backend = service.least_busy(min_num_qubits=5)

# Build circuit
qc = QuantumCircuit(2, 2)
qc.h(0)
qc.cx(0, 1)
qc.measure_all()

# Transpile
from qiskit.compiler import transpile
t_qc = transpile(qc, backend=backend, optimization_level=3)

# Run with Sampler (returns probability distribution)
sampler = Sampler(backend)
job = sampler.run([t_qc], shots=1024)
result = job.result()
counts = result[0].data.meas.get_counts()
print(f"Results from {backend.name}:", counts)
```

---

## Quantum Information Tools

```python
from qiskit.quantum_info import (
    Statevector, DensityMatrix, Operator,
    state_fidelity, average_gate_fidelity,
    entropy, partial_trace, random_statevector,
    SparsePauliOp, Pauli
)

# State fidelity
psi = Statevector.from_label('+')   # |+⟩
phi = Statevector([1/2**0.5, 0])    # |0⟩ + 0|1⟩ — wait, this is |0⟩/√... 
phi = Statevector([1, 0])           # |0⟩
print(state_fidelity(psi, phi))     # 0.5

# Operator fidelity
U = Operator.from_label('H')   # Hadamard
V = Operator.from_label('X')   # X gate
print(average_gate_fidelity(U, V))

# Entanglement entropy
bell = Statevector([1/2**0.5, 0, 0, 1/2**0.5])
rho = DensityMatrix(bell)
rho_A = partial_trace(rho, [1])
print(f"Entanglement entropy: {entropy(rho_A, base=2):.4f}")  # 1.0 ebit

# Pauli operator expectation
qc = QuantumCircuit(2)
qc.h(0); qc.cx(0, 1)
sv = Statevector(qc)
ZZ = SparsePauliOp('ZZ')
print(f"⟨ZZ⟩ = {sv.expectation_value(ZZ).real:.4f}")  # 1.0 for Bell state
```

---

## Quantum Circuit Library

Qiskit includes many pre-built circuits:

```python
from qiskit.circuit.library import (
    QFT,                    # Quantum Fourier Transform
    PhaseEstimation,        # QPE circuit
    GroverOperator,         # Grover's diffuser
    TwoLocal,               # Hardware-efficient ansatz
    EfficientSU2,           # Another VQA ansatz
    PauliEvolutionGate,     # exp(-iHt) for Pauli H
    ZZFeatureMap,           # Quantum kernel for ML
    RealAmplitudes,         # Parameterized circuit
)
import numpy as np

# QFT
qft = QFT(4)
print(qft.draw())

# Grover operator
from qiskit.circuit.library import GroverOperator
oracle = QuantumCircuit(3)
oracle.ccz(0, 1, 2)  # marks |111⟩
grover_op = GroverOperator(oracle)
print(grover_op.draw())

# Ansatz for VQE
ansatz = EfficientSU2(num_qubits=4, reps=2, entanglement='linear')
print(f"Parameters: {len(ansatz.parameters)}")
print(ansatz.draw())

# Bind parameters
params = {p: np.random.uniform(0, 2*np.pi) for p in ansatz.parameters}
bound = ansatz.assign_parameters(params)
```

---

## Advanced: Custom Gate Definitions

```python
from qiskit.circuit import Gate
from qiskit import QuantumCircuit
import numpy as np

# Define a custom gate
class iSWAP(Gate):
    def __init__(self):
        super().__init__('iswap', 2, [])
    
    def __array__(self, dtype=None):
        return np.array([[1, 0, 0, 0],
                         [0, 0, 1j, 0],
                         [0, 1j, 0, 0],
                         [0, 0, 0, 1]], dtype=dtype)

# Use in circuit
qc = QuantumCircuit(2)
qc.append(iSWAP(), [0, 1])
print(qc.draw())

# Decompose to native gates for simulation
from qiskit.quantum_info import Operator
print(Operator(qc))
```

---

## Primitives API (Modern Qiskit)

```python
from qiskit_ibm_runtime import SamplerV2, EstimatorV2

# Sampler: get measurement statistics
sampler = SamplerV2(backend)
pub = (circuit, )  # PUB = Primitive Unified Bloc
job = sampler.run([pub])
result = job.result()
counts = result[0].data.meas.get_counts()

# Estimator: compute expectation values
from qiskit.quantum_info import SparsePauliOp
estimator = EstimatorV2(backend)
observable = SparsePauliOp('ZZ')
pub = (circuit, observable)
job = estimator.run([pub])
exp_val = job.result()[0].data.evs
print(f"⟨ZZ⟩ = {exp_val}")
```

---

## Complete Example: Variational Circuit on Real Hardware

```python
from qiskit import QuantumCircuit
from qiskit.quantum_info import SparsePauliOp
from qiskit_ibm_runtime import QiskitRuntimeService, EstimatorV2
from qiskit.compiler import transpile
from scipy.optimize import minimize
import numpy as np

# Hamiltonian: Z⊗Z + X⊗I
H = SparsePauliOp.from_list([('ZZ', 1.0), ('XI', 0.5)])

def parameterized_circuit(theta):
    qc = QuantumCircuit(2)
    qc.ry(theta[0], 0)
    qc.ry(theta[1], 1)
    qc.cx(0, 1)
    qc.ry(theta[2], 0)
    return qc

# Load backend
service = QiskitRuntimeService()
backend = service.least_busy(min_num_qubits=2)
estimator = EstimatorV2(backend)

# Cost function (calls quantum hardware)
def cost(theta):
    qc = parameterized_circuit(theta)
    t_qc = transpile(qc, backend=backend)
    pub = (t_qc, H)
    result = estimator.run([pub]).result()
    return result[0].data.evs

# Optimize
theta0 = np.random.uniform(0, 2*np.pi, 3)
result = minimize(cost, theta0, method='COBYLA', options={'maxiter': 100})
print(f"Ground state energy: {result.fun:.4f}")
```

---

## Qiskit Cheat Sheet

```python
# Create         qc = QuantumCircuit(n, m)
# Gates          qc.h/x/y/z/s/t/cx/cz/swap/ccx(qubit)
# Rotation       qc.rx/ry/rz(angle, qubit)
# Measure        qc.measure(qubits, clbits)
# Draw           qc.draw() or qc.draw('mpl')
# Depth          qc.depth()
# Gate count     qc.count_ops()
# Statevector    Statevector(qc)
# Density matrix DensityMatrix(qc)
# Simulate       AerSimulator().run(qc, shots=1024).result().get_counts()
# Transpile      transpile(qc, backend, optimization_level=3)
# Operator       Operator(qc)  or SparsePauliOp('ZZ')
# Expectation    sv.expectation_value(operator)
# Fidelity       state_fidelity(sv1, sv2)
# Compose        qc1.compose(qc2, inplace=True)
# Inverse        qc.inverse()
# Repeat         qc.repeat(k)
# Control        qc.control(num_ctrl_qubits)
```
