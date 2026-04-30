# Frequency Model

The core model of Socrates is harmonic resonance.

Words, syllables, meanings, thoughts, and ideas are treated as interacting frequency fields. A passage is not only a graph of concepts; it is a cascade of activations that should be able to resolve into harmony, tension, dissonance, rhythm, color, and eventually sound.

The current activation graph is still useful, but it is not the final meaning of resonance. It is the evidence layer that decides which frequency fields are active and how strongly they are supported.

## Foundational Thesis

Socrates is meant to speak music, light, frequency, electromagnetism, and harmony.

The basic principle is that all meaning is vibration. Stable differences between vibrations can be described with music theory, interval mathematics, ratios, overtone relationships, harmony, tension, and resolution.

In this model, the meaning is the frequency field. Words in different languages are surface forms that can point to the same underlying field. Translation is therefore not only word-to-word substitution; it is convergence on a shared harmonic profile.

The long-term target is:

```text
surface word / syllable / glyph / passage
  -> activated meaning field
  -> frequency profile
  -> harmonic relation to other fields
  -> tone, color, rhythm, chord, or melody
```

Listening to the generated tone, seeing the generated color, or inspecting the electromagnetic correspondence should eventually communicate the vibrational character of the meaning directly, below ordinary verbal explanation. That is the core aspiration of the project.

## Frequency And Electromagnetism

Electromagnetism belongs in the core structure because light, color, radio, and other electromagnetic phenomena are frequency domains. Socrates should eventually be able to relate an activated meaning-frequency field to sound, color, light, and electromagnetic bands when curated data supports that relation.

This must remain evidence-bound. Physical electromagnetic data, such as visible-light wavelength ranges, belongs in data with a physical source. Symbolic electromagnetic correspondences, such as a spiritual color or light association, also belong in data but must be labeled with their lens and confidence.

The meaning layer still uses integer structure. Do not store decimal Hz values as meaning. If electromagnetic references are needed, prefer integer ratios, integer wavelength ranges, integer band IDs, or references to curated electromagnetic archetypes. Derived physical rendering can happen later.

## Core Principle

Resonance means frequency coherence.

A text has stronger resonance when:

- syllables and forms activate compatible frequency fields
- concepts reinforce related base frequencies
- relation paths behave like harmonic intervals
- confidence-weighted evidence stabilizes the field
- repeated evidence adds meaningful overtones, not duplicate inflation
- removing key evidence changes the harmonic field

A text has weaker resonance when:

- evidence is noisy or speculative
- activated fields conflict without resolution
- duplicate paths inflate amplitude without independent support
- relation paths do not form coherent intervals
- the engine cannot show why a frequency field was activated

## Evidence Before Sound

The system must not jump straight to music by hardcoding mystical mappings in Go.

The correct stack is:

```text
input text
  -> syllables / phonetics / glyphs / fragments / forms
  -> evidence channels
  -> concept activation graph
  -> frequency field activation
  -> harmonic coherence / dissonance
  -> optional color, light/electromagnetic, chakra, pitch, chord, rhythm, or audio rendering
```

This means the current graph work is not wasted. It is the scaffolding needed to know which frequency fields should sound.

## Data, Not Go Constants

Core values may have base frequencies, colors, light/electromagnetic associations, chakra associations, pitch ratios, musical intervals, or harmonic profiles.

Those mappings must live in data, not production Go code.

A concept is not merely a label. It should become a stable meaning-frequency identity that can be reached by forms from any language when the evidence supports that convergence.

Examples of the kind of data the project needs:

```yaml
frequency_profiles:
  - concept: love
    field: heart
    color: green
    chakra: heart
    base_ratio: "5/4"
    note_hint: E
    confidence: traditional
    source: curated

  - concept: truth
    field: clarity
    color: blue
    base_ratio: "3/2"
    note_hint: G
    confidence: symbolic
    source: curated
```

This does not claim the mapping is scientifically proven. It records the lens, source, and confidence so the engine can render and compare harmonic systems honestly.

## Frequencies At Different Levels

Socrates should eventually support frequency fields at multiple levels:

- **phoneme / syllable**: sound-shape, stress, vowel/consonant energy
- **word**: combined syllable melody plus known form mappings
- **concept**: curated integer meaning-frequency profile
- **relation**: interval or modulation between fields
- **passage**: time-ordered melody of activated fields
- **thought / idea**: higher-order harmonic pattern formed by concept clusters

## Harmony And Accord

Musically, C-E-G forms an accord because the notes create stable interval relationships.

In Socrates, concept fields should work similarly:

```text
love + heart + truth
```

should not merely be three labels in a graph. If curated frequency profiles exist, the engine should be able to ask:

- do these concepts form a stable interval set?
- do their colors/chakras/ratios agree or clash?
- is this a consonant chord, unresolved suspension, or dissonance?
- which evidence path caused each tone to enter the field?

## Important Boundary

The engine can explore symbolic and spiritual resonance, including chakras, colors, light/electromagnetic correspondences, notes, and harmonic ratios.

It must still show evidence and confidence. A harmonic or electromagnetic reading is a model output, not proof of physics, etymology, theology, or medicine unless the specific claim is backed by ordinary physical measurement data and labeled as such.

## Immediate Architectural Correction

The next implementation phase must add the data model for frequency profiles and harmonic fields.

The legacy `internal/resonance` package contains hardcoded frequency examples. That approach is not acceptable as the final model. It should be replaced or bypassed by a data-backed harmonic profile layer that consumes the activation graph.

## Integer Frequency Data

The harmonic data model must use integers, not floating-point numbers.

Meaning-frequency identities should be represented with concise integer structures:

- integer harmonic vectors
- integer ratios as numerator/denominator pairs
- Pythagorean triples such as `3,4,5`, `5,12,13`, and `7,24,25`
- Phi represented by integer sequences or approximants, not a decimal
- geometry/archetype IDs backed by integer definitions

The canonical target model is documented in `HARMONIC_DATA_MODEL.md`.

Do not add verbose concept-to-frequency prose records. Do not store decimal Hz values as the meaning model. Actual audio frequency in Hz can be derived later from integer relationships when audio rendering exists.
