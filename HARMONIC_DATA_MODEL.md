# Harmonic Data Model

This is the target data model for Socrates' meaning layer.

Meaning is stored as integer harmonic structure, not as prose and not as floating-point frequency values.

Words are surface forms. Concepts are activation handles. The deeper identity is a language-neutral meaning-frequency field.

## Rules

- Store harmonic values as integers only.
- Do not store floating-point frequencies.
- Do not hardcode concept-to-frequency, concept-to-color, or concept-to-electromagnetic mappings in Go.
- Keep data concise.
- Prefer integer vectors, integer ratios, integer triples, symbolic geometry IDs, and electromagnetic band IDs backed by integer definitions.
- Use words only as labels, evidence, or aliases. The meaning identity is the harmonic field.

## Core Shape

Use short YAML keys in data. Loader structs may use clearer Go names.

```yaml
fields:
  - id: mf.love
    v: [5, 8, 13, 21, 34, 55, 89]
    r: [[5, 4], [3, 2]]
    g: [phi, metatron]
    em: [em.visible.green]
    c: symbolic
    w: 80
    src: curated
```

Keys:

- `id`: language-neutral meaning-frequency ID.
- `v`: integer harmonic vector.
- `r`: integer ratios as `[numerator, denominator]`.
- `g`: linked geometry/archetype IDs.
- `em`: optional linked electromagnetic/light band IDs.
- `c`: confidence label.
- `w`: integer weight, `0..100`.
- `src`: source or curation lens.

## Forms To Meaning

Forms should point toward concepts or directly toward meaning-frequency IDs through data.

```yaml
forms:
  - t: love
    lang: en
    field: mf.love
    c: verified
    w: 90

  - t: amor
    lang: es
    field: mf.love
    c: verified
    w: 90
```

The same field can be reached from many languages. Translation is convergence on the same harmonic identity, not word substitution.

## Relations As Intervals

Relations between meaning fields should also be integer based.

```yaml
intervals:
  - a: mf.love
    b: mf.truth
    r: [3, 2]
    k: fifth
    c: symbolic
    w: 70
```

Keys:

- `a`, `b`: meaning-frequency field IDs.
- `r`: integer ratio.
- `k`: interval/relation kind.
- `c`: confidence.
- `w`: integer weight.

## Integer Archetypes

Sacred geometry and harmonic archetypes are stored as integer definitions.

```yaml
archetypes:
  - id: tri.3.4.5
    k: pythagorean
    n: [3, 4, 5]

  - id: tri.5.12.13
    k: pythagorean
    n: [5, 12, 13]

  - id: phi
    k: continued_fraction
    n: [1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89]

  - id: metatron
    k: graph
    n: [13, 78]
```

Phi is not stored as a float. It is represented by integer structure, such as Fibonacci approximants or another explicitly curated integer sequence.

Pythagorean triples are stored as integer triples:

```text
3,4,5
5,12,13
7,24,25
8,15,17
9,40,41
11,60,61
12,35,37
```

## Electromagnetic References

Electromagnetic data should be represented as integer-backed references, not as decimal meaning values.

```yaml
em_bands:
  - id: em.visible.green
    k: visible_light
    wl_nm: [495, 570]
    c: verified
    src: physics

  - id: em.radio.low
    k: radio
    hz: [30000, 300000]
    c: verified
    src: physics
```

Keys:

- `id`: electromagnetic/light band ID.
- `k`: band kind.
- `wl_nm`: optional integer wavelength range in nanometers.
- `hz`: optional integer frequency range in hertz.
- `c`: confidence label.
- `src`: physical source or curation lens.

Symbolic links from meaning fields to electromagnetic bands are allowed, but they must be labeled as symbolic, traditional, or otherwise lens-specific. Physical ranges do not prove the symbolic mapping.

## Harmonic Scoring

The engine should compare active fields by integer math:

- matching field IDs
- shared archetypes
- shared electromagnetic/light band references
- compatible integer ratios
- low common multiples
- stable Pythagorean or Phi-derived relationships
- confidence-weighted evidence paths

No float comparison is needed for the core model.

## Boundary

Audio rendering is later. The first goal is a concise integer harmonic field model that can say:

```text
these meanings are active
these active meanings map to integer fields
these fields form coherent or dissonant integer relationships
these evidence paths activated them
```
