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
- Direct-evidence/rendering support fix: `truth` single-token input now renders as direct evidence, Evidence and Top Fields display are aligned, graph propagation uses source/base strength to avoid repeated-token inflation, and repeated-token accumulation is regression-tested.
- Preposition/function-word handling: standalone-only structural signals are demoted from Top Fields and prose themes; `in` no longer presents `inward`/`into`/`one` as the concise reading; `There is a deep longing in my heart` now leads with emotional fields instead of preposition-derived fields.

Current git status:

```text
## main...origin/main [ahead 96]
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
**Harmonic/audio rendering foundation.**

Architect review status:
Preposition/function-word handling is accepted:
- standalone function-word signals are weighted down and marked so they can be filtered from top user-facing fields
- standalone-only patterns are excluded from prose theme identification
- `./bin/socrates "in"` no longer presents `inward`, `into`, or `one` in the concise reading
- `./bin/socrates "There is a deep longing in my heart"` leads with `pain`, `bitterness`, and `desire`, not `inward`/`into`
- `./bin/socrates "truth"` still shows `truth` as the main theme
- `./bin/socrates "I feel sad and empty inside"` still leads with `feeling`, `sadness`, and `emptiness`
- Sanskrit transliteration and Devanagari direct fields remain stable in representative checks
- Pi reports `go test ./...` passes and knowledge validation has zero errors

Current blocker:
- Socrates can compute integer harmonic fields, but it cannot yet render them into any user-facing tone, interval, chord, color swatch, or simple melody representation.
- The next layer should expose harmonic field data without inventing floating-point meaning claims or hardcoded concept-to-audio behavior.

Required investigation:
- Inspect current harmonic output structures and renderer paths before adding behavior.
- Add a first rendering layer that converts existing `HarmonicField` data into a deterministic, inspectable representation such as:
  - tone vector rows
  - ratio labels
  - interval/chord descriptors
  - integer note IDs already present in `frequency_profiles`
  - integer color/field labels already present in `frequency_profiles`
- Keep this as textual or structured output first. Do not require actual audio playback yet unless the existing CLI already has a safe place for it.
- Do not add concept-specific audio mappings in Go. Rendering must derive from curated frequency profiles and integer labels.
- Do not introduce floating-point frequency values as meaning. If a pitch preview is added later, it must be clearly a rendering choice, not the stored meaning.
- Add tests proving:
  - known harmonic inputs such as `truth`, `प्राण`, `dharma`, and `love` produce stable render descriptors
  - unknown/weak inputs do not produce authoritative renderings
  - rendering output is deterministic
  - existing harmonic field calculation remains unchanged
- Default output may show a concise harmonic summary; debug output may show full render internals.

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

- A deterministic harmonic render descriptor exists for harmonic fields.
- Render descriptors are derived only from existing integer harmonic field data and curated frequency profiles.
- Known harmonic inputs produce stable descriptors in tests.
- Weak/unknown inputs do not produce authoritative render descriptors.
- No concept-to-audio/color/frequency hardcoding is introduced in production Go.
- No floating-point meaning-frequency storage is introduced.
- No broad new knowledge area or runtime channel is introduced.
- `./bin/socrates knowledge validate` passes with no errors.
- `go test ./...` passes.
- Work is committed locally and not pushed.


## Next After This

After harmonic rendering foundation:
- Private app/UI workflow for personal journaling
- Saved reflection sessions and field history
- Consult `docs/ARCHITECTURE_READINESS.md` before major new feature layers

## Nice To Have Later

- extended script support
- harmonic/audio rendering layer
