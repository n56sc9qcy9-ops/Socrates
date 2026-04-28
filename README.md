# Socrates

Socrates is an experimental language-resonance engine.

The project explores whether words, sounds, glyphs, scripts, fragments, and passages can be analyzed as activation fields. It is not a dictionary, journaling app, chatbot, or mystical text generator.

The base model is coherent activation energy first: a text resonates when independent evidence channels reinforce a shared activation field.

The engine should show evidence first: activated concepts, graph paths, convergence scores, and uncertainty. Generated forms, candidate lists, and full similarity matches should be available when debugging.

## Core Principle

The code must not know what specific words mean.

Go code may:

- generate candidate forms
- compute similarity
- load knowledge data
- activate concepts
- propagate activation through relations
- score convergence
- render evidence

Go code must not:

- contain semantic keyword maps
- branch on specific words for meaning
- make examples pass through hardcoded behavior
- claim symbolic readings are proven facts

All linguistic, symbolic, spiritual, harmonic, musical, and cross-language knowledge belongs in data.

## Resonance

Resonance means coherent activation energy.

A word or passage has stronger resonance when form, sound, glyph, fragment, passage, and relation evidence converge around related concepts. It has weaker resonance when matches are noisy, speculative, duplicated, or disconnected.

The project may later render a coherent activation field as harmonics or music, but audio is not the base engine. Harmonic output should consume evidence from the activation graph; it must not replace the evidence model or hardcode concept-to-pitch mappings in Go.

## Core Flow

```text
input text
  -> generated forms
  -> similarity matches against known data forms
  -> concept activation
  -> graph propagation
  -> coherent activation energy
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
knowledge/lenses.yaml
knowledge/confidence.yaml
```

SQLite can be considered later if YAML becomes too limited.

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
