# Chapter 5: Protein Structure & AlphaFold

## The Protein Folding Problem

A protein's function is determined by its 3D structure, and its structure is determined by its amino acid sequence. Predicting structure from sequence is the "holy grail" of computational biology.

```
Sequence: MGSSHHHHHHSSGLVPRGSHMRGPNPTAASLEASAGPF...  (1D string)
     ↓ (folding)
Structure: complex 3D shape with helices, sheets, loops
     ↓
Function: enzyme, receptor, antibody, structural protein, etc.

Levinthal's Paradox:
  If a protein tries all possible conformations at 10^13/sec:
  100-aa protein: 10^300 conformations × 10^-13 sec = > age of universe
  Yet proteins fold in microseconds to milliseconds!
  → There must be a folding pathway, not random search
```

---

## Protein Structure Levels

```python
# Primary: amino acid sequence (string)
sequence = "MVLSPADKTNVKAAWGKVGAHAGEYGAEALERMFLSFPTTKTYFPHFDLSHGSAQVKGHGKKVADALTNAVAHVDDMPNALSALSDLHAHKLRVDPVNFKLLSHCLLVTLAAHLPAEFTPAVHASLDKFLASVSTVLTSKYR"
# ^ Hemoglobin alpha subunit

# Secondary: local structure patterns
secondary_elements = {
    'helix': 'LSPADKTNVKAAWGK',    # alpha-helix (spiral backbone)
    'sheet': 'VGAHAGEYGAEAL',       # beta-sheet (flat, hydrogen-bonded strands)
    'loop':  'ERMFL',               # loop/turn (irregular connections)
}

# Tertiary: complete 3D fold (most important for function)
# Represented as atomic coordinates (PDB format)

# Quaternary: multi-chain assembly
# Hemoglobin = 2 alpha + 2 beta chains = tetramer
```

---

## PDB Format: Protein Structure Data

```python
from Bio.PDB import PDBParser, PDBIO, Select
import numpy as np

def load_pdb(pdb_file: str):
    """Load protein structure from PDB file."""
    parser = PDBParser(QUIET=True)
    structure = parser.get_structure("protein", pdb_file)
    return structure

def extract_ca_coordinates(structure) -> np.ndarray:
    """Extract Cα (alpha carbon) coordinates — backbone representation."""
    coords = []
    for model in structure:
        for chain in model:
            for residue in chain:
                if 'CA' in residue:
                    coords.append(residue['CA'].get_vector().get_array())
    return np.array(coords)

def calculate_rmsd(coords1: np.ndarray, coords2: np.ndarray) -> float:
    """Root Mean Square Deviation — structural similarity metric."""
    if len(coords1) != len(coords2):
        raise ValueError("Structures must have same number of atoms")
    diff = coords1 - coords2
    return np.sqrt(np.mean(np.sum(diff**2, axis=1)))

def get_secondary_structure(structure, dssp_path: str = 'dssp') -> dict:
    """Calculate secondary structure using DSSP algorithm."""
    from Bio.PDB.DSSP import DSSP
    model = list(structure)[0]
    dssp = DSSP(model, list(structure.get_models())[0], dssp=dssp_path)
    
    ss_dict = {}
    for key in dssp.property_keys:
        ss_dict[key] = dssp[key]
    return ss_dict
```

---

## AlphaFold2: The Revolution

AlphaFold2 (DeepMind, 2020) achieves near-experimental accuracy for protein structure prediction:

```
Key innovations:
1. Multiple Sequence Alignment (MSA) as input — evolutionary information
2. Evoformer: transformer that processes MSA + pair representations together
3. Structure module: directly outputs 3D coordinates
4. Self-distillation training trick
5. End-to-end differentiable (gradient flows through geometry)

Performance: CASP14 — TM-score > 0.9 for most domains
               Effectively SOLVED the protein structure prediction problem!

Input:  amino acid sequence
Output: 3D atomic coordinates + confidence (pLDDT score)
```

```python
# Using AlphaFold2 via ColabFold (easiest access)
# https://colab.research.google.com/github/sokrypton/ColabFold/blob/main/AlphaFold2.ipynb

# Or via ESMFold (Meta) — faster, no MSA needed
import torch
from transformers import EsmForProteinFolding, EsmTokenizer

def predict_with_esmfold(sequence: str) -> dict:
    """Predict protein structure using ESMFold (no MSA needed)."""
    model_name = "facebook/esmfold_v1"
    tokenizer = EsmTokenizer.from_pretrained(model_name)
    model = EsmForProteinFolding.from_pretrained(model_name)
    model.eval()
    
    tokenized = tokenizer(sequence, return_tensors="pt", add_special_tokens=False)
    
    with torch.no_grad():
        output = model(**tokenized)
    
    pdb_string = model.output_to_pdb(output)[0]
    
    return {
        'pdb': pdb_string,
        'plddt': output.plddt.squeeze().cpu().numpy(),  # per-residue confidence
        'mean_plddt': output.plddt.mean().item(),
    }

# Interpretation of pLDDT confidence scores
def interpret_plddt(plddt_scores: np.ndarray) -> dict:
    """Interpret per-residue confidence scores."""
    return {
        'very_high': np.sum(plddt_scores >= 90) / len(plddt_scores),  # backbone correct
        'high': np.sum((plddt_scores >= 70) & (plddt_scores < 90)) / len(plddt_scores),
        'low': np.sum((plddt_scores >= 50) & (plddt_scores < 70)) / len(plddt_scores),  # loops
        'very_low': np.sum(plddt_scores < 50) / len(plddt_scores),  # disordered regions
    }
```

---

## Protein Structure Databases

```python
# PDB (Protein Data Bank): 230,000+ experimental structures
# AlphaFold DB: 200 million+ predicted structures (entire UniProt!)

def fetch_alphafold_structure(uniprot_id: str) -> str:
    """Download AlphaFold predicted structure."""
    import requests
    url = f"https://alphafold.ebi.ac.uk/files/AF-{uniprot_id}-F1-model_v4.pdb"
    response = requests.get(url)
    if response.status_code == 200:
        return response.text
    raise ValueError(f"No AlphaFold structure for {uniprot_id}")

def fetch_pdb_structure(pdb_id: str) -> str:
    """Download experimental structure from PDB."""
    import requests
    url = f"https://files.rcsb.org/download/{pdb_id}.pdb"
    return requests.get(url).text

# Search by sequence
def search_pdb_by_sequence(sequence: str, identity_cutoff: float = 0.3) -> list:
    """Find PDB structures with similar sequences."""
    import requests
    url = "https://search.rcsb.org/rcsbsearch/v2/query"
    query = {
        "query": {
            "type": "terminal",
            "service": "sequence",
            "parameters": {
                "evalue_cutoff": 1,
                "identity_cutoff": identity_cutoff,
                "sequence_type": "protein",
                "value": sequence
            }
        },
        "return_type": "entry"
    }
    response = requests.post(url, json=query)
    return [r['identifier'] for r in response.json().get('result_set', [])]
```

---

## Structural Analysis with AI

```python
import torch
from torch_geometric.data import Data
from torch_geometric.nn import GCNConv

class ProteinGNN(torch.nn.Module):
    """Graph Neural Network for protein property prediction."""
    
    def __init__(self, node_features: int, hidden_dim: int = 128, output_dim: int = 1):
        super().__init__()
        self.conv1 = GCNConv(node_features, hidden_dim)
        self.conv2 = GCNConv(hidden_dim, hidden_dim)
        self.conv3 = GCNConv(hidden_dim, hidden_dim)
        self.pool = torch.nn.AdaptiveAvgPool1d(1)
        self.fc = torch.nn.Linear(hidden_dim, output_dim)
        self.relu = torch.nn.ReLU()
    
    def forward(self, data: Data) -> torch.Tensor:
        x, edge_index = data.x, data.edge_index
        x = self.relu(self.conv1(x, edge_index))
        x = self.relu(self.conv2(x, edge_index))
        x = self.relu(self.conv3(x, edge_index))
        
        # Global pooling
        x = torch.mean(x, dim=0, keepdim=True)
        return self.fc(x)

def protein_to_graph(coords: np.ndarray, sequence: str,
                     distance_threshold: float = 8.0) -> Data:
    """Convert protein structure to graph for GNN."""
    n = len(coords)
    
    # Node features: amino acid properties
    aa_features = encode_amino_acids(sequence)  # 20-dim one-hot + physicochemical
    
    # Edge: connect residues within distance threshold
    edges = []
    for i in range(n):
        for j in range(i+1, n):
            dist = np.linalg.norm(coords[i] - coords[j])
            if dist < distance_threshold:
                edges.append([i, j])
                edges.append([j, i])
    
    edge_index = torch.tensor(edges, dtype=torch.long).T if edges else torch.zeros((2, 0), dtype=torch.long)
    
    return Data(
        x=torch.tensor(aa_features, dtype=torch.float),
        edge_index=edge_index,
        pos=torch.tensor(coords, dtype=torch.float)
    )
```

---

## Protein Language Models

ESM-2 (Meta) is a transformer trained on protein sequences — the BERT of proteins:

```python
# ESM-2 for protein embeddings
from transformers import EsmModel, EsmTokenizer
import torch

def get_protein_embeddings(sequences: list, model_name: str = "facebook/esm2_t6_8M_UR50D") -> np.ndarray:
    """Extract protein embeddings using ESM-2."""
    tokenizer = EsmTokenizer.from_pretrained(model_name)
    model = EsmModel.from_pretrained(model_name)
    model.eval()
    
    embeddings = []
    for seq in sequences:
        inputs = tokenizer(seq, return_tensors="pt")
        with torch.no_grad():
            outputs = model(**inputs)
        # Mean pooling over residues
        embedding = outputs.last_hidden_state.mean(dim=1).squeeze().numpy()
        embeddings.append(embedding)
    
    return np.array(embeddings)

# Use for:
# 1. Protein similarity search (cosine similarity in embedding space)
# 2. Function prediction (train classifier on embeddings)
# 3. Variant effect prediction
# 4. Clustering protein families

def predict_variant_effect(wt_seq: str, variant_seq: str) -> float:
    """Predict if variant is deleterious using ESM-2."""
    embeddings = get_protein_embeddings([wt_seq, variant_seq])
    # Cosine distance between WT and variant embeddings
    from sklearn.metrics.pairwise import cosine_similarity
    sim = cosine_similarity([embeddings[0]], [embeddings[1]])[0][0]
    return 1 - sim  # higher = more different = potentially more deleterious
```

---

## Drug-Target Binding Prediction

```python
# Molecular docking: predict how a drug binds to a protein
# Tools: AutoDock Vina, GLIDE, Rosetta

def vina_dock(protein_pdb: str, ligand_mol: str, box_center: tuple,
              box_size: tuple = (20, 20, 20)) -> dict:
    """Run molecular docking with AutoDock Vina."""
    import subprocess
    
    cmd = [
        "vina",
        "--receptor", protein_pdb,
        "--ligand", ligand_mol,
        "--center_x", str(box_center[0]),
        "--center_y", str(box_center[1]),
        "--center_z", str(box_center[2]),
        "--size_x", str(box_size[0]),
        "--size_y", str(box_size[1]),
        "--size_z", str(box_size[2]),
        "--out", "docked_out.pdbqt",
    ]
    result = subprocess.run(cmd, capture_output=True, text=True)
    
    # Parse best docking score from output
    import re
    scores = re.findall(r'\s+1\s+([-\d.]+)', result.stdout)
    best_score = float(scores[0]) if scores else None
    
    return {'binding_affinity_kcal_mol': best_score, 'output': result.stdout}

# Deep learning approach: EquiBind, DiffDock
# Much faster than traditional docking; uses GNN + diffusion models
```

---

## Key Structural Bioinformatics Tools

```
Structure prediction:   AlphaFold2, ESMFold, RoseTTAFold (diffusion-based)
Protein design:         RFdiffusion (backbone), ProteinMPNN (sequence)
Docking:                AutoDock Vina, Glide, DiffDock
Visualization:          PyMOL, ChimeraX, UCSF Chimera
Analysis:               Bio.PDB, MDAnalysis, ProDy
Databases:              PDB, AlphaFold DB, CATH, SCOP
MD Simulation:          GROMACS, AMBER, OpenMM, NAMD
```
