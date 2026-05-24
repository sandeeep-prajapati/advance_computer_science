# Chapter 19: Quantum Hardware

## Physical Implementations of Qubits

Different technologies encode quantum information in different physical systems:

| Technology | Company | Qubit type | Gate fidelity | T₂ coherence |
|---|---|---|---|---|
| Superconducting | IBM, Google | Energy levels of Josephson junction | 99.5% (1Q), 99% (2Q) | ~100 μs |
| Trapped ion | IonQ, Quantinuum | Electronic states of ions | 99.9% (1Q), 99.5% (2Q) | >10 minutes |
| Photonic | PsiQuantum, Xanadu | Photon polarization/path | Variable | Infinite (no decoherence) |
| Neutral atom | QuEra, Pasqal | Rydberg states of neutral atoms | 99.5% (1Q), 99% (2Q) | ~5 s |
| Silicon spin | Intel | Electron spin in silicon | ~99% | ms |
| Topological | Microsoft | Non-Abelian anyons | Theoretical | Inherently protected |

---

## Superconducting Qubits (IBM, Google)

The most advanced commercial technology:

```
Physics:
  Josephson junction: superconductor-insulator-superconductor
  At ~15 mK (colder than outer space), becomes quantum harmonic oscillator
  Anharmonicity makes it a 2-level system (qubit)
  
  Qubit types: Transmon (most common), Flux qubit, Charge qubit
  
Gate operations:
  Single-qubit: microwave pulses (~10-50 ns)
  Two-qubit: cross-resonance, CZ via flux tuning (~100-500 ns)
  
Cooling requirements: dilution refrigerator (~15 mK)
Scale: 50-1000+ physical qubits on chip

IBM Quantum roadmap:
  2021: 127 qubits (Eagle)
  2022: 433 qubits (Osprey)
  2023: 1121 qubits (Condor)
  2025: 100,000+ qubits (modular systems)
  
Connectivity:
  Heavy-hex topology: each qubit connects to 2-3 others
  Not all-to-all → SWAP gates needed for non-adjacent operations
```

```
IBM Heavy-Hex Topology (Eagle, 127 qubits):
  ○─○─○─○─○─○─○─○─○─○─○─○
  │   │   │   │   │   │  
  ○   ○   ○   ○   ○   ○  
  │   │   │   │   │   │  
  ○─○─○─○─○─○─○─○─○─○─○─○
```

---

## Trapped Ion Qubits (IonQ, Quantinuum)

```
Physics:
  Individual ions (e.g., Yb⁺, Ca⁺) trapped by electromagnetic fields
  Qubit: two electronic energy levels of the ion
  Gates: laser pulses (optical, ~1-100 μs) or microwave pulses
  
Advantages:
  Very long coherence times (>10 minutes vs ~100 μs for superconducting)
  High gate fidelity (>99.9% single-qubit, >99.5% two-qubit)
  All-to-all connectivity (any ion can interact with any other)
  No need for dilution refrigerator (room temperature trap)
  
Disadvantages:
  Slow gate speeds (μs vs ns for superconducting)
  Difficult to scale to thousands of qubits
  Photonic interconnects needed for modular systems
  
Key companies:
  IonQ: Ytterbium ions, publicly traded, AWS/Azure integration
  Quantinuum (Honeywell): Barium ions, H-series, highest fidelity
  Oxford Ionics: Erbium, chip-scale
  
Scale: 20-30 "algorithmic qubits" currently
```

---

## Photonic Quantum Computing

```
Physics:
  Qubits: photon polarization (H/V) or dual-rail (path)
  Gates: beam splitters, phase shifters, photon detectors
  
Types:
  Gaussian Boson Sampling: Xanadu's Borealis
  Linear optical QC (LOQC): probabilistic gates
  Fusion-based QC: PsiQuantum's approach
  
Advantages:
  Room temperature operation (no cooling)
  Natural for quantum communication/networks
  Photons are perfect flying qubits
  
Disadvantages:
  Photon loss (each optical element absorbs some photons)
  Probabilistic gates (entangling gates succeed with probability <1)
  Measurement is destructive
  
PsiQuantum approach:
  Fault-tolerant silicon photonics
  Manufacturing at GlobalFoundries (standard semiconductor fab)
  Target: 1 million+ photonic qubits
```

---

## Neutral Atom Qubits (QuEra, Pasqal)

```
Physics:
  Neutral atoms (Rb, Cs) held in optical tweezers (focused laser beams)
  Ground/excited states or hyperfine levels as qubit states
  Long-range interactions via Rydberg states (highly excited atoms)
  
Advantages:
  Large-scale connectivity via Rydberg interactions
  Reconfigurable: move atoms to create different topologies
  Long coherence times
  Hundreds to thousands of qubits
  
Disadvantages:
  Moderate gate fidelities (~99.5%)
  Atom loading is probabilistic
  Vacuum system required
  
QuEra:
  Aquila system: 256 neutral atom qubits (publicly available via AWS Braket)
  Programmable Hamiltonian simulation
  
Notable result: Harvard/QuEra (2023) demonstrated 48 logical qubits via error correction
  First demonstration of fault-tolerant algorithms at this scale
```

---

## Hardware Metrics

### Gate fidelity:
```
Average Gate Fidelity F = (1 - ε)
  ε = error per gate

Single-qubit fidelity: typically 99.5–99.9%
Two-qubit fidelity: typically 99.0–99.7%
State preparation and measurement (SPAM) error: ~1-5%

Error per gate needed for fault-tolerance: < ~1%
Current superconducting: ~0.1-1% (at the threshold)
Current trapped ion: ~0.01-0.1% (better!)
```

### Quantum Volume:
```
IBM's benchmark for processor capability:
  QV = max m such that m-qubit random circuit has heavy output > 2/3

Factors: qubit count, connectivity, gate fidelity, measurement error

QV progression:
  2017: QV=4
  2019: QV=16
  2020: QV=64
  2022: QV=512
  2023: QV=8192 (IBM Torino)

Higher QV = more complex circuits can run successfully
```

### Circuit Layer Operations per Second (CLOPS):
```
Speed metric: how fast can a backend run parametric circuit layers?
  Includes compilation, execution, and result retrieval

IBM Eagle: ~5,000 CLOPS
Important for VQA optimization (many circuit evaluations needed)
```

---

## Quantum Cloud Access

```
IBM Quantum:
  Free tier: Falcon (5Q), Heron (133Q), Eagle (127Q)
  Access via qiskit-ibm-runtime
  Jobs queue (wait time varies)

Google Quantum AI:
  Research collaboration required
  Access via cirq and Google Cloud

AWS Braket:
  Multiple providers: IonQ, Rigetti, OQC, QuEra, IQM
  Pay per shot/task
  
Azure Quantum:
  IonQ, Quantinuum, Rigetti
  Quantum credits for researchers
  
IonQ:
  Harmony, Aria, Forte systems
  Available on AWS/Azure/GCP

Quantinuum:
  H1-1, H2 systems (highest fidelity available)
  Commercial and research access
```

```python
# Connect to IBM Quantum (free tier)
from qiskit_ibm_runtime import QiskitRuntimeService

service = QiskitRuntimeService(channel="ibm_quantum", token="YOUR_TOKEN")
backends = service.backends(operational=True, min_num_qubits=5)
for b in backends[:5]:
    print(f"{b.name}: {b.num_qubits}Q, queue={b.status().pending_jobs}")

# Connect to AWS Braket
# pip install amazon-braket-sdk
from braket.aws import AwsDevice
from braket.circuits import Circuit

device = AwsDevice("arn:aws:braket:::device/quantum-simulator/amazon/sv1")
circuit = Circuit().h(0).cnot(0, 1)
task = device.run(circuit, shots=1000)
result = task.result()
print(result.measurement_counts)
```

---

## Noise Characterization

Understanding hardware noise is critical for algorithm design:

```python
from qiskit_ibm_runtime import QiskitRuntimeService
from qiskit.providers import BackendV2

# Get backend properties
service = QiskitRuntimeService()
backend = service.backend('ibm_brisbane')

# View coupling map
print(backend.coupling_map)

# View error rates (via calibration data)
from qiskit.transpiler import CouplingMap
cm = backend.coupling_map

# Randomized benchmarking (estimate gate errors)
# (requires qiskit-experiments)
from qiskit_experiments.library import StandardRB

qubits = [0]
rb_exp = StandardRB(qubits, lengths=range(1, 200, 20), backend=backend)
rb_data = rb_exp.run(backend).block_for_results()
result_data = rb_data.analysis_results("EPC")
print(f"1Q gate error: {result_data.value:.4f}")
```

---

## Roadmap to Fault-Tolerant QC

```
NISQ era (now, 2024):
  Physical qubits: 50–1000
  Error rate: 0.1–1%
  No error correction
  Useful for: NISQ algorithms (VQE, QAOA), quantum simulation demo

Early FTQC era (2027–2030?):
  Physical qubits: 10,000–100,000
  Logical qubits: 10–100 (with surface codes)
  Surface code threshold crossed
  Useful for: Small-scale error-corrected computations

Full FTQC era (2030s+):
  Physical qubits: millions
  Logical qubits: thousands
  Useful for: Shor's algorithm, large-scale quantum simulation, quantum chemistry

Key milestones to watch:
  1. Sustained error below surface code threshold (~1%)
  2. First logical qubit outperforming physical qubits
  3. Fault-tolerant T gate via magic state distillation at scale
  4. Break-even: fault-tolerant computation faster than classical
```

---

## Comparing Hardware Technologies

```
For algorithm developers:
  Short-term advantage: Trapped ion (high fidelity, all-to-all)
  Long-term scale: Superconducting (faster fabrication scaling)
  New paradigm: Neutral atoms (reconfigurable, large-scale)
  
For networking/communication: Photonic
For topological protection: Microsoft's topological (not yet realized)

"Horses for courses":
  Optimization (short depth): Superconducting (speed)
  Deep circuits (high fidelity): Trapped ion
  Quantum simulation: Neutral atoms (programmable Hamiltonians)
  Communication: Photonic
```
