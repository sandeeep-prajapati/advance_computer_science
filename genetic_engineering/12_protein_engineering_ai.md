# Chapter 12: Protein Engineering with AI

## From Protein Analysis to Protein Design

Protein engineering asks: "Can we change or create a protein with a desired function?"

```
Analyze: sequence -> structure -> function
Design:  function -> structure -> sequence
```

Traditional protein engineering used directed evolution: create variants, screen them, keep winners, repeat. AI adds better search, better priors, and better scoring.

---

## Design Problem Types

```python
design_tasks = {
    "stability": "Make a protein fold more reliably",
    "binding": "Improve affinity or specificity for a target",
    "enzyme": "Improve catalytic activity or substrate preference",
    "de_novo": "Create a new fold or binder from scratch",
    "humanization": "Reduce immunogenicity while preserving function",
}
```

Search space is enormous:

```python
def protein_search_space(length: int) -> int:
    return 20 ** length

print(protein_search_space(100))  # not searchable by brute force
```

---

## Directed Evolution as Optimization

```python
import random

AMINO_ACIDS = "ACDEFGHIKLMNPQRSTVWY"

def mutate_sequence(seq: str, mutation_rate: float = 0.02) -> str:
    chars = list(seq)
    for i, aa in enumerate(chars):
        if random.random() < mutation_rate:
            chars[i] = random.choice(AMINO_ACIDS.replace(aa, ""))
    return "".join(chars)

def directed_evolution_round(parent: str, score_fn, library_size: int = 1000) -> str:
    library = [mutate_sequence(parent) for _ in range(library_size)]
    return max(library, key=score_fn)
```

The practical bottleneck is usually the assay. ML is useful because it learns from small expensive experiments and chooses the next variants to test.

---

## Protein Language Models

Protein LMs learn constraints from millions of natural sequences.

```python
model_uses = {
    "embeddings": "Represent sequences for classifiers or clustering",
    "masked_prediction": "Suggest plausible amino acids at positions",
    "variant_effect": "Estimate whether a mutation is tolerated",
    "generation": "Sample new sequences with natural-like statistics",
}
```

Pseudo-workflow:

```python
def rank_variants(parent: str, candidate_variants: list[str], embedding_model, predictor):
    embeddings = [embedding_model(seq) for seq in candidate_variants]
    scores = predictor.predict(embeddings)
    return sorted(zip(candidate_variants, scores), key=lambda x: x[1], reverse=True)
```

---

## Structure-Conditioned Design

Given a target backbone, sequence design finds amino acids that stabilize it.

```
Backbone coordinates -> sequence design model -> candidate sequences
candidate sequences  -> structure prediction -> filter by confidence
filtered designs     -> docking/functional scoring -> experimental shortlist
```

Tools to know:

```
ProteinMPNN       sequence design for fixed backbones
RFdiffusion       generate protein backbones and binders
AlphaFold/ColabFold structure prediction and confidence
ESM-IF            inverse folding from structure to sequence
Rosetta           physics/statistical protein modeling
FoldX             mutation stability estimation
```

---

## Active Learning Loop

```python
import numpy as np

def acquisition_score(pred_mean: np.ndarray, pred_std: np.ndarray, beta: float = 2.0):
    """Upper confidence bound: exploit high mean, explore high uncertainty."""
    return pred_mean + beta * pred_std

def select_next_batch(candidates, model, batch_size: int = 96):
    pred_mean, pred_std = model.predict_with_uncertainty(candidates)
    scores = acquisition_score(pred_mean, pred_std)
    order = np.argsort(scores)[::-1]
    return [candidates[i] for i in order[:batch_size]]
```

Active learning works well when each experimental batch is costly and you can only test a tiny subset of the design space.

---

## Design Filters

```python
filters = {
    "identity_to_parent": "Avoid drifting too far unless de novo design",
    "predicted_stability": "Discard destabilized variants",
    "aggregation": "Remove sticky hydrophobic patches",
    "motifs": "Preserve catalytic or binding residues",
    "developability": "Solubility, expression, manufacturability",
    "off_targets": "Avoid undesired binding or immune risk",
}
```

AI does not remove the need for wet-lab validation. It makes the shortlist better.

---

## Key Takeaways

- Protein engineering is search over a huge constrained sequence space.
- Protein LMs provide powerful priors; structure models provide geometric validation.
- Active learning connects ML with real assays.
- Final success is measured experimentally, not by model confidence alone.

