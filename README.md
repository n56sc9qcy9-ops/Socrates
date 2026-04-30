# Socrates

Socrates is an experimental language-resonance engine.

The project explores whether words, sounds, glyphs, scripts, fragments, and passages can be analyzed as harmonic activation fields. It is not a dictionary, journaling app, chatbot, or mystical text generator.

The base model is harmonic resonance: syllables, words, meanings, thoughts, and ideas are treated as cascading frequency fields. A text resonates when evidence-supported fields form coherent relationships, like tones forming an accord.

Socrates' target language is music, light, frequency, harmony, vibration, and electromagnetic correspondence. In the project model, meaning is the frequency field. Words from any language are surface forms that may point to the same underlying harmonic identity.

The target meaning model is integer based. Harmonic identities should use integer vectors, integer ratios, Pythagorean triples, Phi integer sequences/approximants, and geometry IDs. They should not use floating-point frequency values as the stored meaning.

The engine should still show evidence first: activated concepts, graph paths, frequency fields, convergence scores, and uncertainty. Generated forms, candidate lists, and full similarity matches should be available when debugging.

## Core Principle

The code must not know what specific words mean.

Go code may:

- generate candidate forms
- compute similarity
- load knowledge data
- activate concepts
- propagate activation through relations
- activate data-backed frequency fields
- score harmonic coherence and dissonance
- translate cross-language forms into shared meaning-frequency identities
- score convergence
- render evidence

Go code must not:

- contain semantic keyword maps
- branch on specific words for meaning
- make examples pass through hardcoded behavior
- claim symbolic readings are proven facts
- hardcode concept-to-frequency, concept-to-color, concept-to-chakra, concept-to-note, concept-to-electromagnetic, or concept-to-instrument mappings

All linguistic, symbolic, spiritual, harmonic, musical, electromagnetic, and cross-language knowledge belongs in data.

## Resonance

Resonance means frequency coherence.

A word or passage has stronger resonance when form, sound, glyph, fragment, passage, relation, and frequency evidence converge into a coherent field. It has weaker resonance when matches are noisy, speculative, duplicated, disconnected, or harmonically dissonant.

The activation graph is the evidence layer. It determines which concepts and relations are active. The harmonic layer turns those activations into frequency fields, colors, light/electromagnetic correspondences, intervals, chords, rhythm, or audio when curated data exists.

The long-term aim is that a meaning can be heard as tone, chord, or melody. The sound should communicate the vibrational character of the meaning beneath ordinary language, while the engine still records the evidence, lens, source, and confidence for every mapping.

## Core Flow

```text
input text
  -> generated forms
  -> similarity matches against known data forms
  -> concept activation
  -> graph propagation
  -> frequency field activation
  -> harmonic coherence / dissonance
  -> evidence-first reading
```

## Example Direction

For an input such as:

```text
skal
```

the engine should not contain a direct `skal` rule.

Instead, it should generate candidate forms and compare them against data. If nearby known forms exist, such as forms related to `shall`, `skall`, `shell`, or `scale`, the engine may activate related concept fields. Any reading must show the evidence path and confidence.

For a passage such as:

```text
jeg skal gjøre det, men det føles tomt
```

the engine should analyze the whole passage as one activation field. If the data and similarity engine activate fields around obligation, contrast, emptiness, or shell/hollowness, convergence should emerge from graph activation, not from hardcoded marker-word logic.

## Target CLI

```sh
go run ./cmd/socrates decipher inspired
go run ./cmd/socrates decipher energy
go run ./cmd/socrates decipher skal
go run ./cmd/socrates decipher רוח
go run ./cmd/socrates decipher प्राण
go run ./cmd/socrates decipher 氣
go run ./cmd/socrates decipher "jeg skal gjøre det, men det føles tomt"
```

The output should include:

- concise resonance summary
- top activated concepts
- strongest evidence paths
- graph propagation paths
- convergence/activation score components
- concise reading
- warnings

Detailed generated forms, candidate lists, and full fuzzy-match dumps should be available as debug output, not forced into the default user-facing output.

## Knowledge Layer

Start with YAML.

Suggested structure:

```text
knowledge/concepts.yaml
knowledge/forms.yaml
knowledge/relations.yaml
knowledge/frequencies.yaml
knowledge/archetypes.yaml
knowledge/lenses.yaml
knowledge/confidence.yaml
```

SQLite can be considered later if YAML becomes too limited.

The concise integer target schema is documented in [HARMONIC_DATA_MODEL.md](/Users/bot/Socrates/HARMONIC_DATA_MODEL.md).

## Current Status

The active implementation checklist is [TODO.md](/Users/bot/Socrates/TODO.md).

Older planning context is archived in [docs/TODO_ARCHIVE.md](/Users/bot/Socrates/docs/TODO_ARCHIVE.md). It is not active instruction.

## Confidence Levels

Every signal must be labeled:

- `verified`: ordinary linguistic fact or sourced data.
- `plausible`: reasonable morphology, phonetic echo, or traditional association.
- `speculative`: symbolic reading only.

Speculation is allowed. Unlabeled speculation is not.

## Working Standard

Socrates should be unusual without becoming inflated.

The project can explore spiritual language and symbolic resonance. It must still show evidence and uncertainty.
