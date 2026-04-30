# TODO Archive

This file is historical context only. It is not an active instruction file.

The active instruction file is `TODO.md`.

This file preserves earlier planning context for Socrates.

The active instruction for the coding agent is `TODO.md`.

## Earlier Direction

The project was first reframed as a grounded reflective-intelligence prototype. That direction was useful for removing inflated language, but it did not match the core product idea.

The core idea is now language-resonance deciphering:

- words as sounds
- words as glyphs
- words as fragments
- words as cross-language echoes
- words as symbolic carriers
- resonance as converging signals across independent channels

## Important Correction

The system must not become:

- a normal journaling app
- a dictionary lookup tool
- a table of hardcoded final readings
- a mystical text generator

The useful prototype is a pattern engine that generates candidate readings and scores convergence.

## Previous Examples To Keep In Mind

- `inspired` may resonate through `in`, `spirit`, `spire`, breath, and inwardness.
- `energy` may generate possible paths such as `en | er | gi`, but must not hardcode "energy means one is life".
- Hebrew is important: `רוח / ruach`, `דבר / davar`, `אור / or`, `אמת / emet`.
- Sanskrit is important: `ॐ / om`, `प्राण / prana`, `आत्मन् / atman`, `सत् / sat`.
- Chinese is important: `氣 / 气 / qi`, `道 / dao`, `心 / xin`, `真 / zhen`.

These examples should guide tests and seed data, not become final canned interpretations.

## Completed Work: 2026-04-28

The active TODO was shortened after Phase 2B. This section preserves the completed implementation history.

### Phase 1 Correction: Remove Semantic Bucket Branching

Completed:

- Removed `hasModal`, `hasContrast`, and `hasEmptiness` from production code.
- Removed passage APIs that returned semantic bucket booleans.
- Removed hardcoded convergence names such as `hollow_obligation`, `possible_hollow_obligation`, `obligation_present`, and `emptiness_present`.
- Replaced semantic bucket branching with generic concept activation and co-activation scoring.
- Added tests that fail if the semantic bucket fields or hardcoded convergence names return.

Acceptance recorded:

- Production Go searches for semantic bucket fields returned no matches.
- Passage convergence output is derived from generic activated concept fields.

### Phase 1: Remove Shortcut Semantics From Code

Completed:

- Removed direct `skal` fragment behavior from Go code.
- Removed modal, emptiness, and contrast marker maps from Go code.
- Removed passage logic based on `hasModal && hasContrast && hasEmptiness`.
- Replaced shortcut tests with evidence-oriented tests.

Acceptance recorded:

- `skal` appears only in tests/docs/examples, not in semantic implementation.
- No production marker maps remain.

### Phase 2: Data-Driven Knowledge Layer

Completed:

- Created YAML knowledge data for concepts, forms, and relations.
- Implemented a loader package.
- Added confidence validation.
- Added relation endpoint validation.
- Built in-memory indexes for forms, concepts, script words, and relations.
- Moved known forms and relations into data instead of Go code.

Acceptance recorded:

- Adding or removing a known form does not require editing Go source.
- Tests can load fixture knowledge.

### Phase 2 Correction: Remove Active Hardcoded Knowledge Fallbacks

Completed:

- Removed hardcoded modal, shell, emptiness, breath, and concept graph shortcuts from production code.
- Replaced anchor retrieval with knowledge-backed data.
- Replaced relation expansion with knowledge-backed relations.
- Removed legacy fallback from knowledge bridge lookups.
- Added tests proving form and relation removal from fixture data changes activation/path behavior.

### Phase 2B: Remove Remaining Legacy Knowledge Paths

Completed:

- Deleted the hardcoded lexicon source that held `Primitives`, `FragmentSeeds`, and `ScriptWords`.
- Removed or rewrote `LookupFragment` and `LookupPrimitive` so they cannot read hardcoded Go knowledge.
- Updated split confidence scoring to use knowledge-backed lookups.
- Changed engine construction so knowledge loading errors fail clearly instead of falling back to an empty legacy engine.
- Updated explicit fixture engine construction so the provided knowledge base is the active dependency used by analysis.
- Documented the authoritative YAML location.
- Added tests proving explicit fixtures are used by `Analyze`.
- Added tests proving removing a form from active fixture data removes that activation.
- Added tests proving forbidden hardcoded knowledge patterns do not exist in production Go.

Verification recorded by Pi:

```text
rg -n "var FragmentSeeds|var ScriptWords|var Primitives|LookupFragment\(|LookupPrimitive\(" internal/decipher --type go | grep -v "_test.go"
PASS: no matches found in production code

rg -n "Fall back to legacy|legacy behavior" internal/decipher internal/knowledge --type go | grep -v "_test.go"
PASS: no legacy fallback comments found

go test ./...
PASS
```

Verification rechecked by Codex:

- Forbidden hardcoded knowledge names now appear only in the regression test guard.
- No legacy fallback comment remains in production Go.

## Completed Work: 2026-04-29

The active TODO was shortened again after the Activation-Energy Discipline Gate was mostly completed. Detailed implementation history is in local git commits.

### Phase 3A Architecture Cleanup

Completed:

- Split discovery and channel files by responsibility.
- Moved hardcoded glyph/script/cross-language mappings into YAML data or removed the hardcoded path.
- Removed global mutable knowledge dependency.
- Removed root artifacts such as `socrates` and `skal`.
- Added regression coverage against hardcoded semantic/cross-language maps returning.

Recorded verification:

- `go test ./...` passed.
- Local commit: `3576aaf harden bounded decipher work`.

### Bounded Candidate And Fuzzy Work

Completed:

- Bounded candidate generation for long inputs and tight custom caps.
- Preserved core normalized/skeleton/phonetic candidates before speculative candidates.
- Bounded fuzzy comparisons with deterministic ordering.
- Exposed discarded candidate and comparison counts on structured readings.
- Added tests for bounds, discarded counts, and exact/high-confidence priority.

Recorded verification from Pi:

- 53 decipher tests, 9 knowledge tests, and 5 resonance tests passed.

### Default/Debug Output Split

Completed:

- Added render modes/options.
- Kept default output concise and evidence-first.
- Moved full generated forms, candidate lists, fuzzy matches, and detailed internals to debug output.
- Added `--debug` CLI support.
- Added render-mode tests for default concision, debug details, evidence paths, and bounded-work summaries.

Recorded verification from Pi:

```text
ok  socrates/internal/decipher   0.600s
ok  socrates/internal/knowledge  (cached)
ok  socrates/internal/resonance  (cached)
```

Local commit:

- `f0e129e Split default output from debug output`

### Key-Token Removal Regression

Completed:

- Added regression tests proving selected YAML-backed evidence forms weaken when key evidence is removed or perturbed.
- Covered overall score, graph expansion score, passage convergence score, and exact-match weakening.
- Left production code unchanged; no word-specific production branches were added.
- Marked the Activation-Energy Discipline Gate complete.

Recorded verification from Pi:

```text
=== RUN   TestKeyTokenRemovalWeakenActivation
--- PASS
=== RUN   TestKeyTokenRemovalWeakenActivationMultiWord
--- PASS
=== RUN   TestKeyTokenRemovalWeakensConvergence
--- PASS
=== RUN   TestKeyTokenRemovalWeakenExactMatch
--- PASS
PASS
```

Local commit:

- `e4ede94 Add key-token removal regression tests for Activation-Energy Discipline Gate`

### Phase 5 First Activation Graph Slice

Completed:

- Added a first-class `ActivationGraph` with nodes, edges, evidence paths, bounded propagation, and cycle-prevention tests.
- Routed convergence through the graph internally.
- Bounded graph expansion score.
- Deduplicated repeated primitive fragment signals.
- Added a `"God is Love"` regression around bounded score components and duplicate default evidence.

Recorded verification from Pi:

```text
?     socrates/cmd/socrates       [no test files]
ok    socrates/internal/decipher  (cached)
ok    socrates/internal/knowledge (cached)
ok    socrates/internal/resonance (cached)
```

Architect review notes for next task:

- Whitespace can still create repetition evidence unless glyph repetition ignores non-letters.
- Propagation uses mutable visited state and should be made path-local/deterministic.
- Graph-derived nodes need to be clearly distinguished from direct nodes.
- Duplicate edges/evidence need stronger graph-level deduplication.

Local commit:

- `79963d4 Phase 5: Activation Graph layer with bounded scores and deduplication`

### Phase 5/6 Graph Hardening And First Passage Field Slice

Completed:

- Replaced global graph propagation `Visited` state with path-local visited tracking.
- Distinguished direct and graph-expanded nodes by depth.
- Deduplicated graph edges by `(from, to, relationType)`.
- Fixed glyph repetition so whitespace and punctuation are not counted as repeated-letter evidence.
- Added a first `PassageField` API and passage tokenization/analysis tests.

Recorded verification from Pi:

```text
go test ./...
PASS
```

Architect review notes for next task:

- Passage fields exist but are not yet returned on `Reading`.
- Original passage-token provenance needs to be made stricter.
- `PassageField.RelationPaths` exists but is not populated.
- Some tests are permissive and should become strict regressions.
- Graph propagation is currently documented as single-pass because repeated propagation accumulates strength; the contract should be guarded or made idempotent.

Local commit:

- `67031be Phase 5/6: graph hardening and passage field analysis`

### Passage Field Integration Into Reading

Completed:

- Added `PassageFields` to structured `Reading`.
- Routed multi-token analysis through `AnalyzePassageFromTokens` to avoid circular `engine.Analyze` calls.
- Added passage-field summaries to default render output and full passage-field details to debug output.
- Added tests for passage field structure, evidence paths, relation paths, token provenance, repeated-token bounds, and render behavior.

Recorded verification from Pi:

```text
go test ./...
PASS
```

Architect review notes for next task:

- Runtime knowledge is embedded from `internal/knowledge/*.yaml`.
- The project needs a curated knowledge-growth pipeline rather than ad hoc YAML edits or automatic learning.
- `passage_field.go` confidence merging should be checked because the helper semantics can downgrade confidence if used with inverted arguments.

Local commit:

- `c9c0c0e integrate passage fields into readings`

### First Knowledge Growth Pipeline

Completed:

- Added knowledge validation package with errors/warnings.
- Added suggestion/review data structures.
- Added curation documentation.
- Added initial CLI support for knowledge validation.
- Fixed passage-field confidence merge bug per Pi's report.

Recorded verification from Pi:

```text
go test ./...
PASS
```

Architect review notes for next task:

- CLI usage and implementation need alignment: docs say `socrates knowledge validate`, but implementation uses a validation flag internally.
- Alias ambiguity validation needs to check aliases across concepts, not only empty aliases.
- Target-reference policy needs to be stricter or explicitly external/speculative.
- Suggestions exist in code but are not exposed by CLI.
- Validation should assert embedded runtime knowledge has zero errors.

Local commits:

- `d4a69da Build knowledge-growth pipeline: validate, import, and suggest`
- `c32af07 Add KNOWLEDGE_CURATION.md guide for human curators`

### Harmonic Core Architecture Correction

Architect correction:

- Re-established harmonic resonance as the core model, not a future rendering layer.
- Clarified that syllables, words, meanings, thoughts, and ideas are cascading frequency fields.
- Clarified the foundational thesis that meaning is frequency, and words across languages are surface forms pointing to language-neutral meaning-frequency identities.
- Defined the activation graph as the evidence layer that selects and weights frequency fields.
- Added a required next task for YAML-backed frequency profiles and harmonic field scoring.
- Marked the hardcoded legacy `internal/resonance` examples as the wrong final shape.

Implementation direction:

- Concept-to-frequency, color, chakra, note, pitch, interval, and harmonic mappings must live in curated data.
- Go may load, validate, combine, and score frequency profiles, but must not encode those mappings as constants or word branches.
- A `Reading` should expose harmonic field structure only when active concepts have curated profiles.
