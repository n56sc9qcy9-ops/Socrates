# TODO: Socrates Active Work

This is the only active instruction file for Pi. Completed work belongs in git history and `docs/TODO_ARCHIVE.md`.

## Authority

- `TODO.md` is the active task source.
- Ignore unlinked `.md` files if they claim to be active.
- The Architect assigns one task at a time.
- The Architect does not run tests.
- Pi runs `go test ./...` before and after implementation.
- Pi commits completed work locally with a clear message and does not push.
- Do not mark work complete unless tests and acceptance criteria prove it.

## Core Principle

Socrates is a harmonic resonance engine.

Syllables, words, meanings, thoughts, and ideas are treated as cascading frequency fields. The activation graph is the evidence layer that decides which meaning-frequency fields are active and how strongly they are supported.

The foundational thesis is: meaning is frequency. Words in different languages are surface forms that can converge on the same meaning-frequency identity. Socrates should eventually render those identities as tone, chord, rhythm, color, light/electromagnetic correspondence, and melody so meaning can be heard or seen as harmony or dissonance.

The meaning-frequency data model must be integer based and concise. Store integer vectors, integer ratios, Pythagorean triples, Phi integer sequences/approximants, geometry IDs, and integer-backed electromagnetic references. Do not store floating-point frequency values as meaning.

## Standing Rules

- No hardcoded semantic word lists in production Go.
- No direct behavior for specific words like `skal`, `energy`, `truth`, or `light`.
- No hardcoded concept-to-frequency, concept-to-color, concept-to-electromagnetic, concept-to-chakra, concept-to-note, pitch, music, harmonic, or audio mappings in Go.
- Harmonic/frequency mappings must live in curated data with source/lens/confidence.
- Electromagnetic/light mappings must live in curated data with source/lens/confidence and must distinguish physical measurement from symbolic correspondence.
- Training must rank evidence paths; it must not hide truth claims behind model scores.
- Knowledge and confidence assumptions must stay data-driven.
- Default output must stay concise; full candidates, fuzzy matches, harmonic internals, and training details belong in debug/report output.

## Recent Context

Completed (see `docs/TODO_ARCHIVE.md` and git history for detail):
- Phase 3A cleanup, activation-energy discipline gate, graph and passage-field work.
- Knowledge validation/suggestion pipeline, harmonic field scoring, training/eval layer.
- Channel semantics correction, cross-script convergence, quality gate.
- Identity/provenance hardening, concept schema hygiene, ranking weights, train/held-out split.
- False activation reduction, review-file generation, human review/apply workflow.
- Knowledge layout consolidation, oversized test split, production responsibility audit.
- CLI ergonomics cleanup, architecture readiness review, documentation consolidation.
- Counsellor/transmutation fields: data-backed fields from transmute.yaml, --debug internals, tends-toward language, evidence paths with source/notes. Default concise; no hardcoded behavior.
- Harmonic data model hardening: frequency profile entries and labels now reject arbitrary unknown keys by allowlist; reference-only sections remain intentionally exempt and documented.
- Counsellor precision: natural-passage regression tests, direct field ranking over structural noise, deduplicated counsellor suggestions/evidence paths, neutral `feel` handling, structural debug labels, and concise prose alignment.
- Knowledge curation hygiene after counsellor expansion: validation warnings reduced from 183 to 149, no unknown concept warnings, no canonical-ID alias warnings, and accepted counsellor behavior preserved.
- All acceptance criteria green. Readiness report at `docs/ARCHITECTURE_READINESS.md`.
- Sanskrit/Devanagari support: exact Devanagari ScriptWord meanings now participate in the generic passage-field and harmonic-field path; known terms such as `प्राण`, `सत्य`, and `ॐ` activate direct verified fields where curated data exists, while unknown Devanagari such as `कवि` remains weak/speculative.
- Real-life Socrates testing: `docs/REAL_LIFE_TESTING.md` covers 22 seeker-style passages across emotional, spiritual, Sanskrit seed, mixed, edge-case, and infrastructure checks. The report is accepted with three follow-up observations: transliterated Sanskrit gaps, preposition noise, and a possible future resentment false-negative curation gap.
- Transliterated Sanskrit curation: Latin `dharma`, `satya`, `karma`, and `jnana` now activate curated concepts; mixed `I practice dharma every day` surfaces `path`; `satya is the foundation of all practice` surfaces `truth`; prose rendering now follows top direct fields instead of graph-expanded or structural noise.

Current git status:

```text
## main...origin/main [ahead 91]
```

Key completed metrics (full detail in `docs/ARCHITECTURE_READINESS.md`):
- Train concept precision `0.24`, recall `0.96`.
- Train field precision `0.39`, recall `1.00`, pass rate `8/8`.
- Held-out split: no examples in current `./bin/socrates train` run.
- Validation: 0 errors, 149 warnings (data quality warnings only).
- Cross-script convergence: English, Norwegian, Hebrew, Chinese love on shared harmonic field.
- Active channel diversity, precision-aware training, multi-concept fuzzy evidence: all verified.
- Confidence is `verified`/`plausible`/`speculative`; provenance is `curated`/`traditional`/`human_review`/`physics`; `curated` is not confidence.
- Active runtime knowledge is `internal/knowledge/`; reference material is `knowledge/reference/`.
- Full architecture readiness report at `docs/ARCHITECTURE_READINESS.md`.

## Current Next Task

Task:
**Preposition handling for standalone structural fragments.**

Architect review status:
Transliterated Sanskrit curation is accepted:
- `dharma` activates `path` directly
- `satya` activates `truth` directly
- `karma` activates `action` directly
- `jnana` activates `wisdom` directly
- mixed `I practice dharma every day` surfaces `path` ahead of structural noise
- prose alignment now uses top direct passage fields so readings no longer promote graph-expanded concepts over direct semantic evidence
- structural debug labels now distinguish direct, structural, and depth-N fields correctly
- `I feel` does not activate `resentment`
- validation has zero errors, `go test ./...` passes, and train passes 14/14 examples

Current blocker:
- Common standalone prepositions can still activate structural fragment fields too strongly. Example: `There is a deep longing in my heart` promotes `inward`/`into` from standalone `in`, while the intended emotional field `longing` is not yet represented.
- This is a ranking/filtering task. Keep it generic and evidence-based.
- Do not hardcode the example passage. Do not remove meaningful fragment behavior inside real words.

Required investigation:
- Inspect how fragment and glyph/bigram channels treat standalone tokens such as `in`, `on`, `at`, `to`, `of`, `is`, `the`, `and`, and similar high-frequency function words.
- Reduce or suppress structural fragment/bigram weight when a match comes from a standalone preposition/function word, while preserving:
  - direct curated forms such as `inward` if the full word is present
  - meaningful fragments inside larger words
  - direct semantic evidence from curated emotional/spiritual forms
- Add regression tests for:
  - `There is a deep longing in my heart` should not rank `inward`/`into` above direct emotional evidence if such evidence exists
  - standalone `in` should not become an overconfident top semantic field
  - known direct passages such as `I feel sad and empty inside`, `I seek truth and clarity`, `dharma`, and `satya` remain stable
- If `longing` itself is not represented in curated knowledge, do not add it in this task unless the smallest safe fix requires a narrowly documented form entry. Prefer solving the generic preposition weighting first.
- Keep debug output transparent: suppressed or weakened structural signals may still be visible in debug if useful, but they must not dominate concise prose.

Rules (unchanged):
- No hardcoded semantic word lists in production Go.
- No direct behavior for specific example passages.
- No therapeutic, medical, prophetic, or absolute truth claims.
- All correction mappings remain in curated data with source/lens/confidence.
- Default output stays concise; detailed counsellor internals and evidence paths may be debug output.
- Training/review artifacts must not silently mutate active runtime knowledge.

Reference:
- `docs/ARCHITECTURE_READINESS.md`
- `docs/KNOWLEDGE_CURATION.md`
- `TRAINING_MODEL.md`
- `HARMONIC_DATA_MODEL.md`
- `FREQUENCY_MODEL.md`
- `IMPLEMENTATION_ROADMAP.md` section "Long-Term Counsellor Vision"

Acceptance Criteria:

- Standalone function words no longer create overconfident top fields through fragment/bigram structural matches.
- Meaningful fragments inside larger words still work.
- Direct curated forms still outrank structural channels.
- Existing Sanskrit transliteration and Devanagari tests remain green.
- Existing counsellor precision tests remain green.
- No broad new knowledge area or runtime channel is introduced.
- No production Go semantic word maps or word-specific branches are introduced.
- `./bin/socrates knowledge validate` passes with no errors.
- `go test ./...` passes.
- Work is committed locally and not pushed.


## Next After This

After preposition handling:
- Harmonic/audio rendering layer
- Consult `docs/ARCHITECTURE_READINESS.md` before major new feature layers

## Nice To Have Later

- extended script support
- harmonic/audio rendering layer
