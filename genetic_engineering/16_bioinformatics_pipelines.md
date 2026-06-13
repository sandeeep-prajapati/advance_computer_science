# Chapter 16: Bioinformatics Pipelines

## Why Pipelines Matter

Bioinformatics analysis is rarely one script. It is a graph of tools, files, metadata, containers, and compute environments.

```
FASTQ -> QC -> trim -> align -> sort/index -> call variants -> annotate -> report
```

A good pipeline is:

```
reproducible
restartable
parallelizable
portable
versioned
auditable
```

---

## Workflow Engines

```python
engines = {
    "Nextflow": "Popular for cloud/HPC bioinformatics; DSL2 modules",
    "Snakemake": "Pythonic rules; great for academic pipelines",
    "WDL/Cromwell": "Common in genomics consortia and Terra",
    "CWL": "Portable workflow specification",
}
```

Use a workflow engine when file dependencies and compute steps become hard to reason about manually.

---

## Minimal Snakemake Example

```python
SAMPLES = ["sample1", "sample2"]

rule all:
    input:
        expand("results/{sample}.bam.bai", sample=SAMPLES)

rule align:
    input:
        fastq="reads/{sample}.fastq.gz",
        ref="reference/genome.fa"
    output:
        bam="results/{sample}.bam"
    shell:
        "bwa mem {input.ref} {input.fastq} | samtools sort -o {output.bam}"

rule index:
    input:
        bam="results/{sample}.bam"
    output:
        bai="results/{sample}.bam.bai"
    shell:
        "samtools index {input.bam}"
```

Run:

```bash
snakemake --cores 8 --use-conda
```

---

## Minimal Nextflow Example

```groovy
nextflow.enable.dsl=2

params.reads = "reads/*.fastq.gz"
params.ref = "reference/genome.fa"

process FASTQC {
    publishDir "results/fastqc", mode: "copy"

    input:
    path read

    output:
    path "*_fastqc.zip"

    script:
    """
    fastqc ${read}
    """
}

workflow {
    reads_ch = Channel.fromPath(params.reads)
    FASTQC(reads_ch)
}
```

Run:

```bash
nextflow run main.nf -profile docker
```

---

## Containers and Environments

```dockerfile
FROM mambaorg/micromamba:1.5.8

COPY environment.yml /tmp/environment.yml
RUN micromamba install -y -n base -f /tmp/environment.yml && micromamba clean -a -y
```

Best practice:

```
pin tool versions
record reference genome build
store checksums for key inputs
separate raw data from derived outputs
make reports machine-readable and human-readable
```

---

## Metadata and Sample Sheets

```csv
sample_id,fastq_1,fastq_2,condition,batch
S1,reads/S1_R1.fastq.gz,reads/S1_R2.fastq.gz,control,b1
S2,reads/S2_R1.fastq.gz,reads/S2_R2.fastq.gz,treatment,b1
```

Validate metadata before launching expensive jobs:

```python
import pandas as pd
from pathlib import Path

def validate_samplesheet(path: str):
    df = pd.read_csv(path)
    required = {"sample_id", "fastq_1", "condition"}
    missing = required - set(df.columns)
    if missing:
        raise ValueError(f"Missing columns: {missing}")
    for fastq in df["fastq_1"]:
        if not Path(fastq).exists():
            raise FileNotFoundError(fastq)
    return df
```

---

## Cloud and HPC

```
HPC:   SLURM/PBS, shared filesystem, module system
Cloud: object storage, containers, autoscaling, spot instances
```

Key engineering concerns:

```
data locality
retry behavior
cost limits
logging
provenance
intermediate file cleanup
security and access control
```

---

## Key Takeaways

- Pipelines are software products, not throwaway scripts.
- Reproducibility depends on versions, references, metadata, and containers.
- Workflow engines turn file dependencies into an executable DAG.
- Start simple, but design so failed jobs can resume without rerunning everything.

