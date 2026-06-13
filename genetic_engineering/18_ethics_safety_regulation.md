# Chapter 18: Ethics, Safety, and Regulation

## Why Safety Is Part of the Engineering Problem

Genetic engineering can produce enormous benefit: therapies, diagnostics, crops, industrial enzymes, and environmental tools. It can also create harm if misused, poorly validated, or deployed without consent and oversight.

For a CS/AI practitioner, safety is not an appendix. It is part of the system requirements.

---

## Core Ethical Questions

```python
questions = {
    "consent": "Who agreed to the collection, analysis, editing, or sharing?",
    "equity": "Who benefits, who pays, and who is excluded?",
    "privacy": "Can genomic data identify people or relatives?",
    "reversibility": "Can an intervention be stopped or recalled?",
    "germline": "Could changes affect future generations?",
    "dual_use": "Could the same knowledge enable harm?",
}
```

---

## Somatic vs Germline Editing

```python
editing_context = {
    "somatic": "Affects treated individual only; common focus of gene therapy",
    "germline": "Affects embryos/gametes and can be inherited; highly restricted",
}
```

Most accepted therapeutic editing targets somatic cells for serious disease. Heritable human genome editing raises much harder scientific, ethical, and governance problems.

---

## Biosafety and Biosecurity

```python
safety_layers = [
    "risk assessment before work begins",
    "appropriate containment and training",
    "validated protocols and documentation",
    "access controls for sensitive materials/data",
    "incident reporting",
    "independent review for dual-use concerns",
]
```

Dual-use research of concern means legitimate research that could reasonably be misapplied to cause harm. When in doubt, seek institutional review and follow local law and policy.

---

## Human Genomic Privacy

Genomic data is hard to anonymize because it is identifying and familial.

```python
privacy_practices = {
    "minimize": "Collect only data needed for the purpose",
    "separate": "Store identifiers separately from genomic data",
    "encrypt": "Protect data at rest and in transit",
    "audit": "Log access to sensitive datasets",
    "consent": "Respect allowed uses and withdrawal rules",
    "aggregate": "Share summary statistics when individual data is unnecessary",
}
```

AI models trained on genomic or clinical data should be checked for leakage and membership inference risk.

---

## Regulation Overview

Regulatory details vary by country and product type, but common categories include:

```
clinical diagnostics
gene therapies
cell therapies
drugs and biologics
genetically modified crops
laboratory biosafety
human-subjects research
data protection
```

For medical products, evidence must show quality, safety, and efficacy. For research, oversight often includes institutional review boards, biosafety committees, and data access committees.

---

## Responsible AI in Biotech

```python
responsible_ai_checks = [
    "document training data provenance",
    "evaluate across populations and contexts",
    "avoid clinical claims beyond evidence",
    "track uncertainty and failure modes",
    "prevent automated misuse at scale",
    "keep human expert review in high-stakes decisions",
]
```

Useful artifacts:

```
model card
data sheet
risk assessment
validation report
audit log
deployment monitoring plan
```

---

## Key Takeaways

- In biotech, "can we?" and "should we?" must be answered together.
- Human genomic data deserves strong privacy assumptions.
- Medical and environmental deployment requires evidence, oversight, and monitoring.
- Responsible practice protects both people and the credibility of the science.

