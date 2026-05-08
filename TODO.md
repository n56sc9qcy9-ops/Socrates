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
- Harmonic rendering foundation: `HarmonicField` now converts to deterministic render descriptors with tone vectors, integer ratios, integer note/color/field labels, relationship descriptors, and concept-to-tone source trace. This is textual/structured only; no audio playback was added.

Current git status:

```text
## main...origin/main [ahead 98]
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
**Expose harmonic render descriptors in CLI/debug output.**

Architect review status:
Harmonic rendering foundation is accepted:
- `internal/decipher/harmonic_render.go` adds `RenderDescriptor`, tone descriptors, relationship descriptors, and source trace entries
- descriptors are deterministic and derived from existing `HarmonicField` tone/profile data
- vectors and ratios remain integer data
- labels use integer note/color/field IDs already present in frequency profiles
- output is textual/structured only; no audio playback was added
- tests cover descriptor conversion, nil handling, string rendering, relationship descriptions, source trace deduplication, integer vector/ratio shape, and archetype names
- Pi reports all tests pass

Current blocker:
- The render descriptor exists as a code-level foundation, but the CLI does not yet expose it in default or debug output.
- A user still sees the older `HarmonicField(...)` summary unless they are reading code/tests.
- The next task should make render descriptors visible and useful without turning them into audio playback yet.

Required investigation:
- Inspect `render.go`, `engine.go`, and current harmonic field rendering.
- Wire the existing `RenderDescriptor` into CLI rendering:
  - default output should remain concise, for example one extra `Render:` line when a harmonic field exists
  - debug output may show the full `RenderDescriptorToString` block
  - unknown/weak inputs should not show authoritative render descriptors
- Keep the distinction clear:
  - `Harmonic:` is the computed field
  - `Render:` is a deterministic textual representation of that field
  - this is not audio playback and not a physical frequency claim
- Do not add concept-specific audio, color, pitch, or frequency branches in Go.
- Do not add floating-point meaning-frequency storage.
- Add tests proving:
  - default output for known harmonic inputs includes a concise render summary
  - debug output includes vector, ratio, integer labels, relationships, and source trace
  - unknown/weak inputs do not show render output
  - existing concise output remains readable and not bloated
  - descriptor tests from the foundation remain green

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

- Known harmonic inputs expose a concise render summary in default CLI output.
- Debug output exposes full deterministic render descriptor internals.
- Render output is derived only from existing integer harmonic field data and curated frequency profiles.
- Weak/unknown inputs do not show authoritative render output.
- No concept-to-audio/color/frequency hardcoding is introduced in production Go.
- No floating-point meaning-frequency storage is introduced.
- No broad new knowledge area or runtime channel is introduced.
- `./bin/socrates knowledge validate` passes with no errors.
- `go test ./...` passes.
- Work is committed locally and not pushed.


## Next After This

After CLI render exposure:
- Private app/UI workflow for personal journaling
- Saved reflection sessions and field history
- Consult `docs/ARCHITECTURE_READINESS.md` before major new feature layers

## Nice To Have Later

- extended script support
- harmonic/audio rendering layer
