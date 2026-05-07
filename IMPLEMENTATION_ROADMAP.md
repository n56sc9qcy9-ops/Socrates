# Implementation Roadmap

This roadmap is aligned with `TODO.md`. The active checklist is always `TODO.md`; this file explains the sequence at a higher level.

The goal is the real resonance engine:

```text
knowledge data
  -> in-memory indexes
  -> generic form generation
  -> similarity matching
  -> activation graph
  -> coherent activation energy
  -> evidence-first reading
```

The code must not contain semantic keyword maps. Meanings, known forms, concepts, and relations belong in data.

The base model is coherent activation energy first. Harmonic or audio output is a later rendering layer; it should consume the activation graph after the core can produce bounded, deduplicated, evidence-first fields.

## Long-Term Goals Preserved

The long-term vision is still present, but it is intentionally sequenced after the core engine is trustworthy.

Preserved future goals:

- coherent activation energy as the base resonance model
- data-driven linguistic, symbolic, cross-language, harmonic, and musical knowledge
- integer-only meaning-frequency storage
- human curation of proposed knowledge changes
- optional AI assistance as a hypothesis generator only
- harmonic/audio rendering of activation fields as listenable music
- possible later storage, UI, light, or hardware interfaces

## Long-Term Counsellor Vision

Socrates' long-term aim is to help a person inspect a thought-field for resonance, dissonance, contraction, and possible correction paths. It should support a seeker in noticing where a thought pattern appears off balance, while never presenting symbolic resonance as absolute proof, medical advice, prophecy, or external authority.

In this vision, truth is not ingested as an unexamined assertion. Truth enters Socrates as curated, inspectable knowledge: concepts, forms, relations, harmonic profiles, confidence labels, sources, lenses, and tests. A reading is trustworthy only to the degree that it can show the evidence paths that activated each field.

The counsellor layer should eventually answer questions such as:

- Which fields are active in this passage or thought?
- Which active fields appear coherent, contracted, dissonant, unresolved, or weakly supported?
- Which correction fields are suggested by curated transmutation relations?
- What evidence, confidence, source, and lens produced each suggestion?
- Which parts of the reading are verified, plausible, speculative, or unknown?

The system must preserve human sovereignty. It may suggest "this activated field tends toward this correction path under this lens"; it must not say "this is the truth about you" or imply clinical, spiritual, or future certainty.

Before UI, audio, hardware, or broad spiritual guidance work, the core must reliably activate obvious direct fields in ordinary passages. For example, plain-language emotional or spiritual input should not be dominated by low-level glyph, phonetic, or fragment noise when direct concept evidence is present. Field correction is not ready until source fields, target fields, confidence, evidence paths, and uncertainty render clearly by default or debug mode.

Not active yet:

- literal claims that meaning is scientifically proven to be frequency
- audio synthesis before evidence paths are clean
- concept-to-frequency mappings in Go
- floating-point frequency values as the stored meaning model
- verbose prose-heavy harmonic records
- AI silently writing active knowledge

Older visionary language about truth, love, wisdom, source, and frequency should be treated as philosophical background unless converted into data-backed concepts, relations, tests, or documented future rendering requirements.

## Phase 1: Remove Shortcut Semantics

Goal: delete prototype shortcuts from production code.

Tasks:

- Remove hardcoded semantic maps such as modal, emptiness, and contrast marker maps.
- Remove direct word-specific semantic behavior from Go code.
- Remove tests that pass because a word is directly hardcoded.
- Keep only generic form generation, similarity, graph activation, scoring, and rendering logic.

Done when:

- searching production Go files finds no semantic marker maps
- examples are explained through evidence paths, not direct branches

## Phase 2: Data-Driven Knowledge Layer

Goal: move all linguistic and symbolic knowledge out of Go code.

Tasks:

- Create `knowledge/` files for concepts, forms, relations, lenses, and confidence labels.
- Implement a loader.
- Validate confidence values.
- Validate relation endpoints.
- Build in-memory indexes.
- Allow tests to load fixture knowledge.

Done when:

- adding, removing, or changing a known form requires editing data, not Go code
- invalid knowledge data fails with clear errors

## Phase 3: Generic Form Generation

Goal: generate candidate forms without knowing meanings.

Tasks:

- Normalize input.
- Tokenize text and passages.
- Generate phonetic forms.
- Generate consonant and vowel skeletons.
- Generate n-grams.
- Generate prefix/suffix fragments.
- Generate edit-distance candidates.
- Deduplicate candidates.
- Attach method and distance to every candidate.

Done when:

- unknown input still produces useful candidate forms
- generated candidates are visible in CLI output

## Phase 4: Similarity Engine

Goal: compare input candidates against known forms generically.

Tasks:

- Implement normalized edit distance.
- Implement n-gram similarity.
- Implement phonetic similarity.
- Implement skeleton similarity.
- Combine scores with configurable weights.
- Deduplicate matches by best evidence.

Done when:

- close forms score higher than noisy forms
- exact matches score higher than fuzzy matches
- no example requires a direct code path

## Phase 4A: Activation-Energy Discipline Gate

Goal: ensure scores measure coherent activation, not evidence volume.

This gate comes before additional feature work. The engine should not become more expressive until it stops rewarding duplicate paths, weak fuzzy noise, and verbose debug output.

Tasks:

- Define coherent activation energy in score terms.
- Deduplicate evidence by meaningful path.
- Prevent repeated fragments from inflating scores without independent support.
- Bound candidate generation for long input.
- Bound fuzzy matching work and expose discarded counts.
- Make exact/verified multi-channel evidence outrank speculative single-channel evidence.
- Split default output from debug output.
- Keep default CLI output concise.

Done when:

- duplicated evidence does not increase resonance after the first meaningful occurrence
- weak fuzzy matches do not dominate scores
- removing a key token weakens the relevant activation field
- user-facing output is concise by default
- debug output can still expose full candidates and match details

## Phase 5: Activation Graph

Goal: make coherent activation energy emerge from graph structure.

Tasks:

- Define activations and activation edges.
- Activate concepts from form matches.
- Propagate activation through relations.
- Decay activation by relation strength.
- Track evidence paths.
- Prevent graph loops.
- Rank activated concepts.

Done when:

- convergence emerges from graph structure
- output can show why each concept became active
- graph strength reflects independent support rather than duplicated evidence

## Phase 6: Passage Field Analysis

Goal: analyze an entire passage as one resonance field.

Tasks:

- Generate candidates for each token.
- Match each token against known data forms.
- Merge token activations into one passage graph.
- Track token evidence.
- Detect repeated concept fields through activation counts.
- Detect co-activated concepts through graph proximity.
- Score passage-level convergence.

Done when:

- removing a key token changes the convergence score
- passage readings are produced without hardcoded marker-word logic

## Phase 7: Evidence-First Output

Goal: make every reading inspectable without making default output noisy.

Tasks:

- Show concise activation summary.
- Show activated concepts.
- Show propagation paths.
- Show convergence fields.
- Show weak signals.
- Show score components.
- Keep final prose short and secondary.
- Add debug mode for generated forms, candidates, and full similarity matches.

Done when:

- users can see the evidence before the interpretation
- speculative paths are labeled
- candidate and fuzzy-match dumps are opt-in

## Phase 8: Review And Curation

Goal: allow knowledge to evolve without changing code.

Tasks:

- Export candidate matches and graph paths.
- Save proposed knowledge changes to review files.
- Save proposed weight changes to review files.
- Track proposed, accepted, and rejected data.
- Keep rejected data out of active analysis.
- Never silently write training or AI suggestions into active runtime YAML.

Done when:

- humans can curate resonance knowledge
- the engine remains deterministic without AI
- accepted changes have a review trail before they enter active knowledge

## Phase 9: Supervised Ranking And Verification

Goal: make Socrates trainable and measurable without turning it into a black-box authority.

Tasks:

- Add curated training examples outside active knowledge.
- Define expected concepts, meaning-frequency fields, and harmonic qualities.
- Extract evidence-path features from `Reading`.
- Evaluate current rankings against expected fields.
- Report precision, recall, missed fields, false activations, and evidence-path validity.
- Add configurable ranking weights for evidence channels.
- Tune weights only after evaluation reports are stable.
- Suggest improved form/relation/profile weights from training results.
- Write suggested weight changes to review files, not active knowledge.

Done when:

- `socrates train` can prove whether a knowledge/ranking change improved or degraded readings
- expected fields are ranked above false activations for curated examples
- every training result still shows the evidence paths that caused it
- proposed weight changes include before/after metrics and require human acceptance before active YAML changes

## Phase 10: Optional AI Assistance

Goal: use AI only as a hypothesis generator.

AI may suggest:

- candidate forms
- possible translations
- possible symbolic relations
- prose summaries

AI must not silently write active knowledge. Suggestions should be speculative until reviewed.

## Phase 11: Harmonic Rendering Layer

Goal: render coherent activation fields as listenable harmonics.

This phase must not start until the activation graph is clean, bounded, deduplicated, and evidence-first.

Tasks:

- Add the concise integer harmonic data model from `HARMONIC_DATA_MODEL.md`.
- Store meaning-frequency identities as integer vectors, integer ratios, and integer-backed archetypes.
- Represent Phi through integer sequences/approximants, not floats.
- Store Pythagorean triples and geometry systems as integer archetypes.
- Map active concepts to language-neutral meaning-frequency IDs.
- Map evidence paths to integer-backed overtones.
- Map confidence to stability, volume, or timbre.
- Map relation strength to integer interval ratios.
- Map convergence to consonance.
- Map conflict or weak evidence to dissonance.
- Render sentence progression as rhythm or time.

Done when:

- audio rendering consumes activation evidence instead of inventing meaning
- all concept-to-frequency and concept-to-pitch mappings live in data
- stored meaning-frequency values are integer-only
- the user can inspect the evidence before listening

## Phase 12: Optional Storage And Interface

SQLite, UI, light, and hardware are later concerns.

Only add them after the activation graph produces useful evidence.
