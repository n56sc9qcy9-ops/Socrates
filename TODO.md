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
- All acceptance criteria green. Readiness report at `docs/ARCHITECTURE_READINESS.md`.

Current git status:

```text
## main...origin/main [ahead 76]
```

Key completed metrics (full detail in `docs/ARCHITECTURE_READINESS.md`):
- Train field precision `0.39`, recall `1.00`, pass rate `8/8`.
- Held-out field precision `0.27`, recall `0.62`, pass rate `5/8`.
- Validation: 0 errors, 183 warnings (data quality warnings only).
- Cross-script convergence: English, Norwegian, Hebrew, Chinese love on shared harmonic field.
- Active channel diversity, precision-aware training, multi-concept fuzzy evidence: all verified.
- Confidence is `verified`/`plausible`/`speculative`; provenance is `curated`/`traditional`/`human_review`/`physics`; `curated` is not confidence.
- Active runtime knowledge is `internal/knowledge/`; reference material is `knowledge/reference/`.
- Full architecture readiness report at `docs/ARCHITECTURE_READINESS.md`.

## Current Next Task

Task:
**Knowledge curation hygiene after counsellor expansion.**

Architect review status:
Harmonic data model hardening is accepted.

Counsellor precision is accepted after the follow-up fixes:
- focused natural-passage regression tests exist
- direct emotional/spiritual fields outrank structural observations in reviewed passages
- counsellor suggestions and evidence paths deduplicate
- neutral `feel` activates `feeling`, not `resentment`
- structural fields render as structural in debug output
- final prose no longer promotes structural-only signals

Current blocker:
- The counsellor vocabulary expansion increased validation warnings from 135 to 183. There are still no validation errors, but new curation debt should be reduced before the next feature layer.
- Many warnings are duplicate alias/canonical-ID overlaps introduced or exposed by the expanded emotional/spiritual vocabulary. These should be resolved by using relations/neighbors instead of aliases where the target is already a concept.
- Keep the accepted counsellor behavior intact while cleaning data. Do not delete direct forms needed by the new natural-passage tests.

Required investigation:
- Run `./bin/socrates knowledge validate` and group the 183 warnings by cause.
- Prefer data-only fixes in `internal/knowledge/concepts.yaml` and `internal/knowledge/forms.yaml`.
- Remove aliases that duplicate another concept's canonical ID when a relation/neighbor already expresses the connection.
- Keep accepted forms such as `feel`, `sad`, `empty`, `lonely`, `afraid`, `disconnected`, and `truth` working in the natural-passage regression tests.
- Do not hide warnings by weakening validation.
- Rebuild before validation because active runtime knowledge is embedded.

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

- Validation warning count is materially reduced from 183 without creating new validation categories.
- No `unknown concepts` warnings.
- No active form targets a missing concept.
- Natural-passage counsellor precision tests still pass.
- Existing harmonic core, training evaluation, review/apply workflow, validation, and cross-script convergence still pass.
- `./bin/socrates knowledge validate` passes with no errors.
- `go test ./...` passes.
- Work is committed locally and not pushed.


## Next After This

After knowledge curation hygiene:
- Extended script support (Arabic, Thai, etc.)
- Harmonic/audio rendering layer
- Consult `docs/ARCHITECTURE_READINESS.md` before major new feature layers

## Nice To Have Later

- extended script support (Arabic, Thai, etc.)
- harmonic/audio rendering layer
