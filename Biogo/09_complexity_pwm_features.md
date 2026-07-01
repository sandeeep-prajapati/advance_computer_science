# Chapter 9: Complexity, PWM & Features

## Complexity Package

Sequence complexity measures how repetitive or information-rich a sequence is.

Low complexity examples:

```text
AAAAAAAAAAAA
ATATATATATAT
```

Higher complexity example:

```text
ACGTTGCAAGTC
```

Explore:

```bash
go doc ./complexity
rg "complexity|Entropy|Dust" complexity
rg "func Test" complexity
```

---

## PWM Package

PWM means Position Weight Matrix.

It is often used for motif-like scoring:

```text
position 1: A likely
position 2: C likely
position 3: G likely
position 4: T likely
```

Explore:

```bash
go doc ./pwm
rg "PWM|Matrix|Score" pwm
rg "func Test" pwm
```

---

## Feat Package

The `feat` package relates to biological features.

Features can represent regions such as:

- genes
- exons
- regulatory regions
- sequence intervals

Explore:

```bash
go doc ./feat
rg "type .*Feature|Start|End|Location" feat
rg "func Test" feat
```

---

## Contribution Caution

These packages may encode biological meaning.

Before changing behavior:

- understand the biological model
- inspect tests
- add examples
- avoid renaming exported APIs casually
- benchmark scoring or scanning changes

Documentation improvements are very useful here because the concepts are domain-heavy.
