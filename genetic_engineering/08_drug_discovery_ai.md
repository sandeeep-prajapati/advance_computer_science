# Chapter 8: AI-Driven Drug Discovery

## The Drug Discovery Pipeline

```
Target ID → Hit Finding → Lead Optimization → Preclinical → Clinical Trials
  (gene)   (1M+ compounds) (ADMET, potency)   (animal tests)  (human)

Timeline: ~12-15 years, $2.6 billion average
AI impact: compress hit finding + lead opt to months instead of years

Your CS Skills Apply:
  Virtual screening  = search/ranking algorithms
  QSAR              = supervised ML regression
  Generative design = VAE / diffusion models
  Docking           = optimization / physics simulation
  ADMET prediction  = neural network classifiers
```

---

## Molecular Representation

Before training models, represent molecules as data structures:

```python
from rdkit import Chem
from rdkit.Chem import AllChem, Descriptors, Draw
import numpy as np

def smiles_to_fingerprint(smiles: str, radius: int = 2,
                           nbits: int = 2048) -> np.ndarray:
    """Morgan fingerprint (= circular fingerprint = ECFP4)."""
    mol = Chem.MolFromSmiles(smiles)
    if mol is None:
        return None
    fp = AllChem.GetMorganFingerprintAsBitVect(mol, radius, nBits=nbits)
    return np.array(fp)

def smiles_to_descriptors(smiles: str) -> dict:
    """2D physicochemical descriptors (200+ features)."""
    mol = Chem.MolFromSmiles(smiles)
    if mol is None:
        return {}
    return {
        'MW': Descriptors.MolWt(mol),
        'LogP': Descriptors.MolLogP(mol),           # lipophilicity
        'TPSA': Descriptors.TPSA(mol),              # polar surface area
        'HBD': Descriptors.NumHDonors(mol),          # H-bond donors
        'HBA': Descriptors.NumHAcceptors(mol),       # H-bond acceptors
        'RotBonds': Descriptors.NumRotatableBonds(mol),
        'AromaticRings': Descriptors.NumAromaticRings(mol),
        'QED': Chem.QED.qed(mol),                   # drug-likeness 0-1
    }

def lipinski_rule_of_five(smiles: str) -> dict:
    """Check Lipinski's Ro5 for oral drug-likeness."""
    desc = smiles_to_descriptors(smiles)
    violations = sum([
        desc.get('MW', 999) > 500,
        desc.get('LogP', 99) > 5,
        desc.get('HBD', 99) > 5,
        desc.get('HBA', 99) > 10,
    ])
    return {**desc, 'ro5_violations': violations, 'drug_like': violations <= 1}

# Graph representation for GNNs
def smiles_to_graph(smiles: str) -> dict:
    """Convert SMILES to node/edge graph for GNN."""
    mol = Chem.MolFromSmiles(smiles)
    if mol is None:
        return None
    
    # Node features: atom type, charge, hybridization
    atom_features = []
    for atom in mol.GetAtoms():
        features = [
            atom.GetAtomicNum(),
            atom.GetFormalCharge(),
            int(atom.GetHybridization()),
            atom.GetNumImplicitHs(),
            int(atom.GetIsAromatic()),
            atom.GetDegree(),
        ]
        atom_features.append(features)
    
    # Edges: bonds
    edges = []
    edge_features = []
    for bond in mol.GetBonds():
        i, j = bond.GetBeginAtomIdx(), bond.GetEndAtomIdx()
        bond_type = bond.GetBondTypeAsDouble()
        edges.extend([[i, j], [j, i]])
        edge_features.extend([[bond_type], [bond_type]])
    
    return {
        'nodes': np.array(atom_features, dtype=np.float32),
        'edges': np.array(edges, dtype=np.int64).T if edges else np.zeros((2, 0), dtype=np.int64),
        'edge_features': np.array(edge_features, dtype=np.float32),
    }
```

---

## QSAR: ML for Property Prediction

Quantitative Structure-Activity Relationship:

```python
import pandas as pd
import numpy as np
from sklearn.ensemble import RandomForestRegressor, GradientBoostingClassifier
from sklearn.model_selection import cross_val_score
from sklearn.metrics import roc_auc_score

def train_activity_predictor(smiles_list: list, activities: list,
                              model_type: str = 'regression') -> dict:
    """
    Train a QSAR model to predict biological activity from structure.
    
    activities: IC50 values (regression) or active/inactive (classification)
    """
    # Generate fingerprints
    X = np.array([smiles_to_fingerprint(smi) for smi in smiles_list
                  if smiles_to_fingerprint(smi) is not None])
    y = np.array(activities)
    
    if model_type == 'regression':
        # Predict pIC50 (log-transformed IC50)
        model = RandomForestRegressor(n_estimators=200, n_jobs=-1, random_state=42)
        scores = cross_val_score(model, X, y, cv=5, scoring='r2')
        model.fit(X, y)
        print(f"Cross-val R²: {scores.mean():.3f} ± {scores.std():.3f}")
    else:
        model = GradientBoostingClassifier(n_estimators=200)
        scores = cross_val_score(model, X, y, cv=5, scoring='roc_auc')
        model.fit(X, y)
        print(f"Cross-val AUC: {scores.mean():.3f} ± {scores.std():.3f}")
    
    return {'model': model, 'cv_scores': scores}

def applicability_domain(train_fps: np.ndarray,
                          test_fp: np.ndarray,
                          threshold: float = 0.65) -> bool:
    """
    Check if a test compound is within the applicability domain.
    Tanimoto similarity to nearest training neighbor > threshold.
    """
    from rdkit import DataStructs
    max_sim = 0.0
    for train_fp in train_fps:
        sim = np.sum(test_fp & train_fp) / np.sum(test_fp | train_fp)  # Tanimoto
        max_sim = max(max_sim, sim)
    return max_sim >= threshold
```

---

## GNN for Molecular Property Prediction

Graph Neural Networks naturally handle molecular graphs:

```python
import torch
import torch.nn as nn
from torch_geometric.nn import GCNConv, global_mean_pool
from torch_geometric.data import Data, DataLoader

class MolecularGNN(nn.Module):
    """GNN for molecular property prediction."""
    
    def __init__(self, node_features: int = 6, hidden_dim: int = 128,
                 output_dim: int = 1):
        super().__init__()
        self.conv1 = GCNConv(node_features, hidden_dim)
        self.conv2 = GCNConv(hidden_dim, hidden_dim)
        self.conv3 = GCNConv(hidden_dim, hidden_dim)
        self.bn1 = nn.BatchNorm1d(hidden_dim)
        self.bn2 = nn.BatchNorm1d(hidden_dim)
        self.fc = nn.Sequential(
            nn.Linear(hidden_dim, 64),
            nn.ReLU(),
            nn.Dropout(0.3),
            nn.Linear(64, output_dim)
        )
        self.relu = nn.ReLU()
    
    def forward(self, data: Data) -> torch.Tensor:
        x, edge_index, batch = data.x, data.edge_index, data.batch
        x = self.relu(self.bn1(self.conv1(x, edge_index)))
        x = self.relu(self.bn2(self.conv2(x, edge_index)))
        x = self.relu(self.conv3(x, edge_index))
        x = global_mean_pool(x, batch)   # graph-level representation
        return self.fc(x).squeeze(-1)

def smiles_to_pyg_data(smiles: str, label: float = None) -> Data:
    """Convert SMILES to PyG Data object."""
    graph = smiles_to_graph(smiles)
    if graph is None:
        return None
    data = Data(
        x=torch.tensor(graph['nodes'], dtype=torch.float),
        edge_index=torch.tensor(graph['edges'], dtype=torch.long),
    )
    if label is not None:
        data.y = torch.tensor([label], dtype=torch.float)
    return data
```

---

## Generative Models for Drug Design

Generate novel drug-like molecules:

```python
import torch
import torch.nn as nn

class MolecularVAE(nn.Module):
    """
    Variational Autoencoder for molecular generation.
    Encode SMILES → latent space → decode back to SMILES.
    
    Latent space is continuous → interpolate between molecules.
    Key paper: Gómez-Bombarelli et al. 2018
    """
    
    def __init__(self, vocab_size: int, max_length: int = 120,
                 latent_dim: int = 256):
        super().__init__()
        self.latent_dim = latent_dim
        self.max_length = max_length
        
        # Encoder: SMILES sequence → latent vector
        self.encoder_rnn = nn.GRU(vocab_size, 512, num_layers=2,
                                   batch_first=True, bidirectional=True)
        self.fc_mu = nn.Linear(1024, latent_dim)
        self.fc_logvar = nn.Linear(1024, latent_dim)
        
        # Decoder: latent vector → SMILES sequence
        self.fc_z = nn.Linear(latent_dim, 512)
        self.decoder_rnn = nn.GRU(vocab_size + latent_dim, 512, num_layers=2,
                                   batch_first=True)
        self.output_fc = nn.Linear(512, vocab_size)
    
    def encode(self, x):
        _, h = self.encoder_rnn(x)
        h = h[-2:].transpose(0, 1).reshape(x.size(0), -1)  # bidirectional
        return self.fc_mu(h), self.fc_logvar(h)
    
    def decode(self, z, target=None, max_len=None):
        z_expanded = self.fc_z(z).unsqueeze(1)
        # Teacher forcing during training; autoregressive during inference
        if target is not None:
            z_rep = z.unsqueeze(1).expand(-1, target.size(1), -1)
            inp = torch.cat([target, z_rep], dim=-1)
            out, _ = self.decoder_rnn(inp)
        else:
            # Generate step by step
            out_tokens = []
            inp_t = torch.zeros(z.size(0), 1, self.output_fc.out_features).to(z.device)
            h = None
            for _ in range(max_len or self.max_length):
                z_t = z.unsqueeze(1)
                inp = torch.cat([inp_t, z_t], dim=-1)
                out_t, h = self.decoder_rnn(inp, h)
                logits = self.output_fc(out_t)
                token = logits.argmax(dim=-1)
                out_tokens.append(token)
                inp_t = torch.zeros_like(inp_t).scatter_(
                    -1, token.unsqueeze(-1), 1.0)
            out = torch.cat(out_tokens, dim=1)
        return self.output_fc(out)
    
    def forward(self, x):
        mu, logvar = self.encode(x)
        z = mu + torch.randn_like(mu) * torch.exp(0.5 * logvar)
        return self.decode(z, x), mu, logvar

# Bayesian optimization for property-guided generation
def optimize_molecule_bayesian(model: MolecularVAE,
                                property_predictor,
                                seed_smiles: str,
                                n_iterations: int = 100) -> list:
    """
    Optimize molecule toward desired property using Bayesian optimization.
    1. Encode seed to latent z
    2. BO explores latent space
    3. Decode candidates, predict property
    4. Return best molecules
    """
    from botorch.models import SingleTaskGP
    from botorch.acquisition import ExpectedImprovement
    from botorch.optim import optimize_acqf
    import torch
    
    # Encode seed molecule
    # (assumes model.encode is callable on preprocessed input)
    best_candidates = []
    # Simplified: in practice, run full BO loop here
    return best_candidates
```

---

## ADMET Prediction

Predict pharmacokinetic properties before synthesis:

```python
def predict_admet(smiles: str) -> dict:
    """
    Predict ADMET properties using rule-based + ML models.
    ADMET: Absorption, Distribution, Metabolism, Excretion, Toxicity
    """
    desc = smiles_to_descriptors(smiles)
    fp = smiles_to_fingerprint(smiles)
    
    predictions = {}
    
    # Absorption: HIA (Human Intestinal Absorption)
    # Rule: LogP 0-3, MW < 500, TPSA < 140
    predictions['oral_absorption'] = 'likely' if (
        0 <= desc.get('LogP', 99) <= 3 and
        desc.get('MW', 999) < 500 and
        desc.get('TPSA', 999) < 140
    ) else 'poor'
    
    # BBB (Blood-Brain Barrier permeability)
    # Rule: MW < 450, TPSA < 90, LogP 1-3
    predictions['bbb_penetration'] = 'likely' if (
        desc.get('MW', 999) < 450 and
        desc.get('TPSA', 999) < 90 and
        1 <= desc.get('LogP', -99) <= 3
    ) else 'unlikely'
    
    # hERG toxicity (cardiac arrhythmia risk)
    # Rule: basic N, MW > 400, LogP > 3 → higher risk
    from rdkit import Chem
    mol = Chem.MolFromSmiles(smiles)
    basic_n = sum(1 for atom in mol.GetAtoms()
                  if atom.GetAtomicNum() == 7 and atom.GetFormalCharge() == 0)
    predictions['herg_risk'] = 'high' if (
        basic_n > 0 and desc.get('MW', 0) > 400 and desc.get('LogP', 0) > 3
    ) else 'low'
    
    # Ames mutagenicity (toxicity alert)
    # Real: use ML model trained on Ames test data
    predictions['ames_mutagenic'] = 'unknown (use ML model)'
    
    return predictions

# Real ADMET tools: DeepPurpose, ADMET-AI, pkCSM (web API)
def admet_ai_predict(smiles_list: list) -> pd.DataFrame:
    """Use ADMET-AI package for comprehensive predictions."""
    try:
        from admet_ai import ADMETModel
        model = ADMETModel()
        return model.predict(smiles=smiles_list)
    except ImportError:
        print("Install: pip install admet-ai")
        return pd.DataFrame()
```

---

## Virtual Screening Pipeline

Screen millions of compounds against a target:

```python
import subprocess
from pathlib import Path

def virtual_screening_pipeline(target_pdb: str, ligand_library: str,
                                active_site_center: tuple,
                                output_dir: str,
                                top_n: int = 100) -> pd.DataFrame:
    """
    Virtual screening pipeline:
    1. Prepare protein (remove water, add hydrogens)
    2. Prepare ligands (convert to PDBQT)
    3. Dock with AutoDock Vina
    4. Score and rank
    """
    import os
    os.makedirs(output_dir, exist_ok=True)
    
    # Step 1: Prepare receptor
    subprocess.run([
        "prepare_receptor",
        "-r", target_pdb,
        "-o", f"{output_dir}/receptor.pdbqt",
        "-A", "hydrogens"
    ], check=True)
    
    # Step 2: Dock each ligand
    results = []
    for lig_file in Path(ligand_library).glob("*.mol2"):
        # Convert to PDBQT
        pdbqt_file = f"{output_dir}/{lig_file.stem}.pdbqt"
        subprocess.run([
            "prepare_ligand",
            "-l", str(lig_file),
            "-o", pdbqt_file
        ], capture_output=True)
        
        # Dock
        out_file = f"{output_dir}/{lig_file.stem}_out.pdbqt"
        result = subprocess.run([
            "vina",
            "--receptor", f"{output_dir}/receptor.pdbqt",
            "--ligand", pdbqt_file,
            "--center_x", str(active_site_center[0]),
            "--center_y", str(active_site_center[1]),
            "--center_z", str(active_site_center[2]),
            "--size_x", "20", "--size_y", "20", "--size_z", "20",
            "--out", out_file,
            "--exhaustiveness", "8",
        ], capture_output=True, text=True)
        
        # Parse best score
        import re
        scores = re.findall(r'\s+1\s+([-\d.]+)', result.stdout)
        if scores:
            results.append({'ligand': lig_file.stem, 'score': float(scores[0])})
    
    df = pd.DataFrame(results).sort_values('score')
    return df.head(top_n)
```

---

## Key Frameworks and Tools

```
Molecular ML:
  RDKit            — molecular manipulation (open source, essential)
  DeepChem         — ML for drug discovery (sklearn-like API)
  PyTorch Geometric — GNN for molecules
  DGL-LifeSci      — GNN drug discovery models

Generative:
  REINVENT         — RL-based molecule generation (AstraZeneca, open source)
  GuacaMol         — benchmark for molecular generation
  MolGPT           — GPT for molecule SMILES generation
  DiffSBDD         — diffusion model for structure-based drug design

Docking:
  AutoDock Vina    — fast docking (open source)
  GNINA            — DL-enhanced docking (CNNaffinity scoring)
  DiffDock         — diffusion model for blind docking

Property Prediction:
  ADMET-AI         — comprehensive ADMET predictions
  DeepPurpose      — DTI prediction (drug-target interaction)
  ChemProp         — message-passing GNN (MIT, excellent performance)

Databases:
  ChEMBL           — 2M+ bioactive compounds + activities
  PubChem          — 100M+ compounds
  ZINC15           — 750M+ purchasable molecules
  BindingDB        — protein-ligand binding data
```
