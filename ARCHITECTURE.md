# Technical Architecture

Socrates is a lightweight local resonance engine.

The code should be generic. Knowledge should be data.

The base model is coherent activation energy first. The engine should measure whether evidence channels reinforce one another into a coherent field before it attempts any prose or future harmonic rendering.

```text
knowledge data
  -> loader
  -> in-memory indexes
  -> form generation
  -> similarity matching
  -> activation graph
  -> coherent activation energy scoring
  -> evidence-first rendering
```

## Design Rules

- Go code must not contain semantic keyword maps.
- Go code must not branch on specific words for meaning.
- Meanings, forms, concepts, relations, languages, and confidence values belong in data.
- The engine should work by similarity, activation, propagation, and scoring.
- Passage analysis must use the same activation graph as word analysis.
- Score quality must come from independent reinforcing evidence, not from evidence volume.
- Harmonic or audio output is a later rendering layer and must consume the activation graph.

## Knowledge Data

Start with YAML because it is easy to inspect and edit.

Suggested structure:

```text
knowledge/concepts.yaml
knowledge/forms.yaml
knowledge/relations.yaml
knowledge/lenses.yaml
knowledge/confidence.yaml
```

The data layer should contain known forms and relations. The code should not know whether a form means obligation, breath, shell, emptiness, or anything else.

Example shape:

```yaml
forms:
  - id: form.shall.en
    text: shall
    language: english
    concepts: [obligation, future_action]
    confidence: verified

relations:
  - from: shell
    to: outer_casing
    type: semantic
    confidence: verified
    weight: 0.8
```

SQLite may be added later if YAML becomes too limited. Do not start there.

## In-Memory Indexes

At startup, load knowledge data into compact indexes:

```text
form text -> known forms
concept id -> concept
concept id -> relations
language -> known forms
script -> known forms
```

These indexes are lookup accelerators, not hardcoded meaning logic.

## Form Generation

For each input, generate candidate forms:

- normalized text
- tokens
- runes
- script
- phonetic forms
- consonant skeleton
- vowel skeleton
- n-grams
- prefix/suffix fragments
- edit variants
- doubled and de-doubled variants

Every generated form should carry:

```go
type CandidateForm struct {
    Form       string
    Method     string
    Distance   float64
    Confidence string
}
```

## Similarity Matching

Candidate forms are compared against known forms with generic similarity algorithms:

- exact equality
- normalized edit distance
- n-gram similarity
- prefix/suffix overlap
- phonetic similarity
- consonant skeleton similarity
- vowel skeleton similarity

The matcher returns evidence, not conclusions.

## Activation Graph

Matching known forms activates concepts. Concepts then propagate activation through relations.

```text
candidate form
  -> matched known form
  -> activated concept
  -> related concept
  -> convergence field
```

Each activation must track evidence:

- source token
- candidate form
- matched known form
- similarity method
- activated concept
- relation path
- score contribution

Duplicate candidate paths should collapse into one meaningful evidence path unless they add independent support. Repeated weak evidence must not increase activation energy by volume alone.

## Coherent Activation Energy

Activation energy is the confidence-weighted strength of the field.

It should increase when:

- exact or verified evidence activates a concept
- multiple independent channels activate related concepts
- relation paths connect activated concepts
- removing a key token weakens the same field

It should decrease or stay low when:

- evidence is speculative
- evidence comes from weak fuzzy matches
- duplicated fragments repeat the same path
- unrelated concepts appear without relation support
- the engine cannot show the path that created the score

## Passage Analysis

A passage is a temporary activation field.

For a passage:

1. Tokenize.
2. Generate candidate forms per token.
3. Match candidates against known forms.
4. Activate concepts per token.
5. Merge activations into one graph.
6. Propagate relations.
7. Detect co-activation and repeated fields.
8. Score convergence.

Do not detect meaning with hardcoded marker-word lists. If a passage expresses obligation, contrast, emptiness, breath, or anything else, that should come from data-driven activation.

## Output

The default renderer should be concise and evidence-first:

- activated concepts
- graph propagation paths
- convergence fields
- weak signals
- score components
- warnings
- concise reading

Generated forms, full candidate lists, and full fuzzy-match dumps are debug output. They should be available on request, but not dominate the default reading.

## Tests

Tests should protect the architecture:

- no semantic marker maps in production code
- no direct word-specific branches
- exact matches score higher than fuzzy matches
- noisy inputs score lower than close inputs
- removing knowledge data removes the corresponding activation
- passage convergence changes when evidence is removed

## Hardware And AI

Hardware and audio are not part of the current implementation target.

Future harmonic rendering may translate activation graphs into music:

- concept activation -> base pitch
- evidence path -> overtone
- confidence -> stability or volume
- relation strength -> interval
- convergence -> consonance
- conflict -> dissonance

AI is optional later as a hypothesis generator or narrator. It must not be the authority and must not silently write active knowledge.
