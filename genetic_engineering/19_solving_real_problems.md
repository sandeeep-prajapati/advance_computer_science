# Chapter 19: Solving Real Problems

## Case-Study Thinking

Real computational biology work starts with a biological question, then maps it to data, models, validation, and decision-making.

```
question -> data -> analysis -> candidate mechanism -> validation -> decision
```

This chapter shows how to frame problems without getting lost in tool names.

---

## Case Study: Emerging Pathogen Genomics

Problem:

```
Track how a pathogen changes over time and geography.
```

Data:

```
viral/bacterial genomes
sample collection dates
locations
clinical metadata when available
```

Computational workflow:

```python
workflow = [
    "quality-control sequences",
    "align to reference",
    "call variants",
    "assign lineages or clades",
    "build phylogenetic tree",
    "monitor mutations in important genes",
    "share dashboards and reports",
]
```

Core CS concepts:

```
streaming data pipelines
versioned reference datasets
tree algorithms
anomaly detection
visualization
```

---

## Case Study: Sickle Cell Disease

Problem:

```
A mutation in the HBB gene changes hemoglobin behavior and causes sickled red blood cells.
```

Computational questions:

```python
questions = [
    "Which regulatory elements control fetal hemoglobin?",
    "Which edit strategy changes expression safely?",
    "How do we measure editing outcomes?",
    "Which off-targets need careful monitoring?",
    "How durable is the therapeutic response?",
]
```

Data layers:

```
genotype
CRISPR edit outcomes
RNA/protein expression
clinical blood measurements
off-target assessment
long-term follow-up
```

---

## Case Study: CAR-T Cell Therapy

Problem:

```
Engineer immune cells to recognize and attack cancer cells.
```

Computational roles:

```python
roles = {
    "target_selection": "Find antigens enriched on tumor and absent from critical normal tissue",
    "single_cell": "Profile tumor microenvironment and T cell states",
    "sequence_analysis": "Track engineered receptor constructs",
    "manufacturing_qc": "Monitor cell product quality",
    "clinical_modeling": "Study response, toxicity, and relapse",
}
```

Risk thinking:

```
on-target off-tumor effects
cytokine release syndrome
antigen escape
T cell exhaustion
manufacturing variability
```

---

## Case Study: Insulin Production

Problem:

```
Produce human insulin using engineered cells.
```

Computational biology view:

```python
steps = [
    "choose protein sequence",
    "optimize coding sequence for host expression",
    "design expression construct",
    "verify sequence integrity",
    "model expression burden",
    "optimize purification and quality analytics",
]
```

CS/AI contribution:

```
codon optimization
protein stability prediction
process analytics
quality-control automation
batch anomaly detection
```

---

## A General Problem Template

```markdown
## Problem
What biological or clinical decision must improve?

## Data
What data exists? What is missing? What are the biases?

## Model/Analysis
What computational method maps data to evidence?

## Validation
What independent data or experiment would disprove us?

## Decision
What action changes if the result is positive?

## Risk
What failure mode would cause harm or wasted effort?
```

---

## Key Takeaways

- Start from the decision, not the tool.
- Good biological projects combine multiple imperfect evidence sources.
- Validation is the bridge from interesting analysis to useful work.
- The best computational biologists can translate between code, data, experiments, and clinical or industrial constraints.

