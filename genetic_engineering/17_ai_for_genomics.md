# Chapter 17: AI for Genomics

## Genomes as Model Input

Genomics is a natural fit for ML:

```
DNA sequence        -> language modeling
chromatin tracks    -> computer vision / signal processing
gene networks       -> graph learning
variant effects     -> supervised and self-supervised prediction
single-cell states   -> representation learning
```

But biological ML has hard constraints: data leakage, confounding, cell-type specificity, and limited causal labels.

---

## One-Hot Encoding DNA

```python
import numpy as np

BASE_TO_IDX = {"A": 0, "C": 1, "G": 2, "T": 3}

def one_hot_dna(seq: str) -> np.ndarray:
    arr = np.zeros((len(seq), 4), dtype=np.float32)
    for i, base in enumerate(seq.upper()):
        if base in BASE_TO_IDX:
            arr[i, BASE_TO_IDX[base]] = 1.0
    return arr
```

For many regulatory models, the input is a fixed-length DNA window around a candidate enhancer, promoter, or variant.

---

## CNN for Regulatory Prediction

```python
import torch
import torch.nn as nn

class DNAConvNet(nn.Module):
    def __init__(self, num_outputs: int):
        super().__init__()
        self.net = nn.Sequential(
            nn.Conv1d(4, 128, kernel_size=15, padding=7),
            nn.ReLU(),
            nn.MaxPool1d(4),
            nn.Conv1d(128, 256, kernel_size=7, padding=3),
            nn.ReLU(),
            nn.AdaptiveMaxPool1d(1),
            nn.Flatten(),
            nn.Linear(256, num_outputs),
        )

    def forward(self, x):
        # x: batch x length x 4
        x = x.transpose(1, 2)
        return self.net(x)
```

The first-layer filters often learn motif-like patterns.

---

## Transformers for DNA

DNA transformers use tokenization strategies:

```
k-mers:       ACGTAC -> ACG, CGT, GTA, TAC
BPE/unigram: learned subword-like units
byte/base:    one token per base
```

Tasks:

```python
tasks = {
    "masked_language_modeling": "Learn genomic grammar without labels",
    "variant_effect_prediction": "Score reference vs alternate sequence",
    "enhancer_prediction": "Classify regulatory activity",
    "chromatin_prediction": "Predict tracks across cell types",
}
```

---

## Variant Effect Prediction

Score the reference and alternate allele in the same sequence context:

```python
def variant_delta(model, ref_seq: str, alt_seq: str) -> float:
    ref_x = torch.tensor(one_hot_dna(ref_seq)).unsqueeze(0)
    alt_x = torch.tensor(one_hot_dna(alt_seq)).unsqueeze(0)
    with torch.no_grad():
        ref_score = model(ref_x).sigmoid().item()
        alt_score = model(alt_x).sigmoid().item()
    return alt_score - ref_score
```

This is useful for noncoding variants where the effect is regulatory rather than protein-altering.

---

## Graph ML for Biology

```python
graph_types = {
    "gene_regulatory_network": "genes and transcription factors",
    "protein_interaction_network": "physical or functional protein links",
    "metabolic_network": "reactions and metabolites",
    "cell_cell_network": "ligand-receptor communication",
}
```

Graph models help when relationships matter more than individual features.

---

## Common Failure Modes

```python
failure_modes = [
    "train/test leakage through nearby genomic regions",
    "batch effects masquerading as biology",
    "ancestry confounding in human genetics",
    "evaluating on random splits instead of chromosome or donor splits",
    "overclaiming causal interpretation from predictive models",
    "using reference genome artifacts as features",
]
```

For genomic sequence models, chromosome-level splits are often safer than random splits.

---

## Key Takeaways

- DNA can be modeled as sequence, signal, graph, or multimodal data.
- Self-supervised models are powerful, but task-specific evaluation matters.
- Strong validation design is more important than a fancy architecture.
- Predictive biology is not automatically causal biology.

