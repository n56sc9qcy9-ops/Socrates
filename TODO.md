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
- All acceptance criteria green. Readiness report at `docs/ARCHITECTURE_READINESS.md`.

Current git status:

```text
## main...origin/main [ahead 72]
```

Key completed metrics (full detail in `docs/ARCHITECTURE_READINESS.md`):
- Train field precision `0.39`, recall `1.00`, pass rate `8/8`.
- Held-out field precision `0.27`, recall `0.62`, pass rate `5/8`.
- Validation: 0 errors, 135 warnings (data quality warnings only).
- Cross-script convergence: English, Norwegian, Hebrew, Chinese love on shared harmonic field.
- Active channel diversity, precision-aware training, multi-concept fuzzy evidence: all verified.
- Confidence is `verified`/`plausible`/`speculative`; provenance is `curated`/`traditional`/`human_review`/`physics`; `curated` is not confidence.
- Active runtime knowledge is `internal/knowledge/`; reference material is `knowledge/reference/`.
- Full architecture readiness report at `docs/ARCHITECTURE_READINESS.md`.

## Current Next Task

Task:
**Counsellor precision — direct field evidence must outrank structural noise.**

Architect review status:
Harmonic data model hardening is accepted.

Counsellor precision commit `e567fc5` is **not accepted yet**. It improves top-field ordering for the reviewed examples, but it does not satisfy the current task acceptance criteria.

Current blocker:
- No regression tests were added or updated in `e567fc5`, despite the acceptance criterion requiring tests for natural passage field ranking and counsellor suggestion precision.
- Counsellor suggestions can duplicate the same source-target-kind path. Example: `./bin/socrates "I feel afraid and disconnected from truth" --debug` currently shows `isolation --[transmutes_to]--> connection` twice in both suggestions and evidence paths.
- Neutral self-reporting language can create a false direct counsellor source. Example: `feel` currently appears as `resentment` via a weak match to `bitterness`, producing a resentment -> forgiveness suggestion for `I feel afraid...`. This is not acceptable as direct emotional evidence.
- Default prose still does not consistently reflect the improved top-field ranking. Example: `i feel sad and empty inside` shows top fields `sadness`, `emptiness`, `resentment`, but the final reading says `Activated concepts: emptiness, inward, action`.
- Structural channels have been marked `IsDirect:false` at the signal level, but debug rendering still labels fields such as phonetic/glyph observations as `direct`. Clarify whether this is only display wording or a true evidence-classification leak, and fix the misleading output if needed.

Required investigation:
- Add focused regression tests before further broad data expansion. Tests should lock the intended behavior for at least:
  - `i feel sad and empty inside`: top passage fields include `sadness` and `emptiness` above structural observations, and no neutral `feel` -> `resentment` counsellor source appears.
  - `lonely and afraid`: `loneliness` and `fear` outrank structural observations and produce non-duplicated counsellor suggestions.
  - `I feel afraid and disconnected from truth`: `fear`, `isolation`, and `truth` are visible as primary fields; counsellor suggestions are deduplicated.
- Fix counsellor suggestion/evidence-path deduplication by semantic identity such as `source|kind|target|confidence|lens`.
- Fix or remove the `feel` false activation path. Prefer adjusting fuzzy threshold/gating for passage signals or adding a neutral `feeling` concept only if it does not become a transmutation source. Do not map generic `feel` to a contracted field.
- Align default concise reading with `PassageFields.TopFields` or otherwise ensure the final prose does not re-promote structural noise after top-field sorting.
- Keep the current direct-first behavior if tests prove it is the right generic ranking rule, but avoid relying on untested sorting alone.

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

- Plain-language passages with direct concept evidence rank the direct fields above glyph/phonetic/orthographic structural observations.
- Counsellor source fields activate only from direct or clearly evidence-supported passage fields; weak graph-only noise must not produce confident guidance.
- Counsellor output shows source field, suggested correction field, relation kind, confidence, and evidence path in default or debug output.
- Existing `transmute.yaml` mappings remain data-backed; no source-to-target correction mapping is hardcoded in Go.
- New or updated tests cover natural passage field ranking and counsellor suggestion precision without relying on production word-specific branches.
- Existing harmonic core, training evaluation, review/apply workflow, validation, and cross-script convergence still pass.
- `./bin/socrates knowledge validate` passes with no errors.
- `go test ./...` passes.
- Work is committed locally and not pushed.


## Next After This

After counsellor precision:
- Extended script support (Arabic, Thai, etc.)
- Harmonic/audio rendering layer
- Consult `docs/ARCHITECTURE_READINESS.md` before major new feature layers

## Nice To Have Later

- extended script support (Arabic, Thai, etc.)
- harmonic/audio rendering layer
