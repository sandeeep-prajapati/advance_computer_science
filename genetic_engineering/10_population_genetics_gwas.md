# Chapter 10: Population Genetics and GWAS

## Genetic Variation at Population Scale

Population genetics studies how variants differ across individuals and groups. In CS terms, each genome is a massive versioned document, and variants are diffs against a reference.

```
Reference genome:  A C G T A C G T
Person 1 genome:   A C G T A C G T
Person 2 genome:   A T G T A C A T
                         ^       ^
                       SNPs / variants
```

Common variant types:

```python
variant_types = {
    "SNP": "Single base substitution",
    "indel": "Small insertion or deletion",
    "SV": "Structural variant: deletion, duplication, inversion, translocation",
    "CNV": "Copy number variant",
}
```

---

## Allele Frequency

Allele frequency is the proportion of chromosomes carrying a variant.

```python
def allele_frequency(genotypes: list[int]) -> float:
    """
    Genotypes encoded as alternate allele count:
    0 = ref/ref, 1 = ref/alt, 2 = alt/alt
    """
    observed_alt = sum(genotypes)
    total_alleles = 2 * len(genotypes)
    return observed_alt / total_alleles

print(allele_frequency([0, 1, 1, 2, 0, 0]))
```

Rare variants can have large effects. Common variants usually have small effects and are discovered through statistics over many people.

---

## Hardy-Weinberg Equilibrium

For a biallelic variant with allele frequency `p` for reference and `q` for alternate:

```
Expected genotypes:
ref/ref = p^2
ref/alt = 2pq
alt/alt = q^2
```

```python
def hwe_expected_counts(alt_frequency: float, n: int) -> dict[str, float]:
    q = alt_frequency
    p = 1 - q
    return {
        "ref_ref": p * p * n,
        "ref_alt": 2 * p * q * n,
        "alt_alt": q * q * n,
    }
```

Strong deviation can indicate selection, population structure, genotyping error, or non-random mating.

---

## GWAS: Genome-Wide Association Studies

GWAS asks: "Which variants are statistically associated with a trait?"

```python
import pandas as pd
import statsmodels.api as sm

def run_single_variant_gwas(df: pd.DataFrame, snp_col: str) -> dict:
    """
    df columns:
      phenotype: quantitative trait
      snp_col: genotype dosage 0/1/2
      age, sex, pc1, pc2: covariates
    """
    X = df[[snp_col, "age", "sex", "pc1", "pc2"]]
    X = sm.add_constant(X)
    y = df["phenotype"]
    model = sm.OLS(y, X).fit()
    return {
        "beta": model.params[snp_col],
        "pvalue": model.pvalues[snp_col],
    }
```

Why covariates matter:

```
Age/sex          obvious biological confounders
PC1/PC2/...      ancestry/population structure correction
batch            genotyping or sequencing batch effects
relatedness      family structure inflates associations
```

---

## Multiple Testing

GWAS tests millions of variants. The classic significance threshold is very strict:

```python
GENOME_WIDE_SIGNIFICANCE = 5e-8

def bonferroni_threshold(num_tests: int, alpha: float = 0.05) -> float:
    return alpha / num_tests
```

Useful visualizations:

```
Manhattan plot: -log10(p-value) across chromosomes
QQ plot: observed vs expected p-values
LocusZoom: regional association around one hit
```

---

## Polygenic Risk Scores

A polygenic score sums many small effects:

```python
import numpy as np

def polygenic_score(genotype_matrix: np.ndarray, effect_sizes: np.ndarray) -> np.ndarray:
    """
    genotype_matrix: samples x variants, dosage encoded 0/1/2
    effect_sizes: variant weights from GWAS
    """
    return genotype_matrix @ effect_sizes
```

PRS models are sensitive to ancestry mismatch and ascertainment bias. Treat them as probabilistic risk stratification, not deterministic prediction.

---

## Practical Toolchain

```
PLINK / PLINK2        GWAS, QC, pruning, PCA
bcftools              VCF processing
GCTA / BOLT-LMM       mixed models and heritability
SAIGE / REGENIE       large biobank-scale association
PRSice / LDpred       polygenic scores
Hail                  scalable genomic analysis in Python/Spark
```

---

## Key Takeaways

- GWAS finds associations, not automatic mechanisms.
- Population structure is the core bug to control.
- Effect sizes from one cohort may not transfer cleanly to another.
- The best results combine statistics, biology, functional annotation, and replication.

