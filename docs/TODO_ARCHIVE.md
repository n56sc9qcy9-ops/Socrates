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
