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
- All acceptance criteria green. Readiness report at `docs/ARCHITECTURE_READINESS.md`.

Current git status:

```text
## main...origin/main [ahead 51]
```

Key completed metrics (full detail in `docs/ARCHITECTURE_READINESS.md`):
- Train field precision `0.39`, recall `1.00`, pass rate `8/8`.
- Held-out field precision `0.27`, recall `0.62`, pass rate `5/8`.
- Validation: 0 errors, 86 warnings (seed labels only).
- Cross-script convergence: English, Norwegian, Hebrew, Chinese love on shared harmonic field.
- Active channel diversity, precision-aware training, multi-concept fuzzy evidence: all verified.
- Confidence is `verified`/`plausible`/`speculative`; provenance is `curated`/`traditional`/`human_review`/`physics`; `curated` is not confidence.
- Active runtime knowledge is `internal/knowledge/`; reference material is `knowledge/reference/`.
- Full architecture readiness report at `docs/ARCHITECTURE_READINESS.md`.

## Current Next Task

Task:
Calibrate counsellor/transmutation fields so they remain humble, evidence-first, and not overdiagnostic.

Context:
The first counsellor/transmutation layer is implemented and committed. It loads `internal/knowledge/transmute.yaml`, validates transmutation relations, builds a `CounsellorField`, and renders default/debug output with evidence paths. The implementation is directionally correct, but the first output sample shows an architectural risk:

- A single explicit source field can activate neighbor fields, and those neighbor fields can produce multiple correction suggestions.
- The phrase "if field remains X, tends toward Y" can sound like deterministic moral diagnosis.
- The default output can become a list of corrections before the user has enough evidence context.

The goal of this pass is not to add more counsellor content. The goal is to make the counsellor layer restrained, transparent, and aligned with human sovereignty.

Rules:
- No hardcoded word-specific or concept-specific behaviour in production Go.
- All mappings must remain data-driven with source/lens/confidence.
- Training must not mutate active knowledge.
- No black-box truth model.
- No fixed-future prediction.
- Preserve human sovereignty. Counsel should suggest possible correction fields, not command the user.
- No medical, legal, or clinical claims.
- No claim that symbolic readings are proven divine facts.
- Default output must remain concise and evidence-first; detailed counsellor internals belong in debug/report output.
- Do not add new broad spiritual seed packs during this pass.
- Do not lower validation discipline to make warnings disappear.

Reference:

- `docs/ARCHITECTURE_READINESS.md`
- `docs/KNOWLEDGE_CURATION.md`
- `TRAINING_MODEL.md`
- `HARMONIC_DATA_MODEL.md`
- `FREQUENCY_MODEL.md`

Likely files:

- `internal/knowledge/transmute.yaml`
- `internal/knowledge/validate.go`
- `internal/decipher/types.go`
- `internal/decipher/render.go`
- `internal/decipher/counsellor.go`
- counsellor tests under `internal/decipher/`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Run `./bin/socrates knowledge validate` or the documented equivalent and record the result.
4. Review the current counsellor source selection.
   - Distinguish direct evidence from propagated/neighbor evidence in the source field.
   - Neighbor-propagated fields must not produce equal-strength counsel by default.
   - If a suggestion is driven only by propagated evidence, mark it weaker or keep it debug-only unless there is independent supporting evidence.
5. Review default rendering.
   - Replace deterministic-sounding phrasing with a softer evidence phrase.
   - Preferred shape: "possible correction field: source -> target" or "field X may be softened through Y".
   - Avoid "if field remains" in default output unless the evidence clearly supports an inertia statement.
   - Keep relation kind and confidence out of default unless needed for clarity.
6. Add a small cap or ranking rule for default counsellor suggestions.
   - Default output should show the highest-signal few suggestions, not every graph-adjacent correction.
   - Debug output may show all candidates and why they were included or suppressed.
7. Add suppression reasons to debug/report output if useful.
   - Example reasons: weak source, propagated-only source, speculative relation, below display threshold.
   - Keep the implementation light; do not build a large policy framework.
8. Review seed data count and quality.
   - Do not add new entries unless a test genuinely needs one.
   - If any current transmutation relation is too generic or unsupported, lower confidence/weight or leave it debug-only through gating.
9. Preserve completed guardrails:
   - active YAML validates
   - `knowledge validate` works
   - channel diversity counts active evidence only
   - training evaluation is precision-aware
   - multi-concept fuzzy evidence remains complete
   - cross-script love and Hebrew `El` boundary behavior still works
   - concept schema hygiene warnings remain visible and bounded
   - review/apply workflow does not mutate active knowledge without explicit accepted suggestions
10. Add tests proving:
   - direct evidence can produce a concise counsellor suggestion
   - propagated-only neighbor evidence does not flood default output
   - default output is not deterministic or commanding
   - debug output shows why suggestions were included/suppressed
   - no production Go hardcodes particular concepts such as forgiveness, humility, resentment, fear, or love as special cases
11. Run targeted tests while working, then `go test ./...` or `make test`.
12. Commit the completed calibration work locally with a clear message. Do not push.
13. Report final `git status --short --branch`, validation result, tests run, files changed, example default/debug outputs, and any unresolved warnings.

Acceptance Criteria:

- Counsellor/transmutation suggestions are data-backed.
- No hardcoded word-specific or concept-specific counsel exists in Go.
- Suggestions preserve evidence paths, source strength, confidence, and provenance.
- Default output is concise, humble, and does not sound like a fixed diagnosis.
- Default output does not flood the user with correction fields from neighbor expansion.
- Debug output exposes transmutation internals, including included/suppressed reasoning where implemented.
- No fixed-future, medical, legal, or divine-certainty claims are introduced.
- Existing harmonic core, training evaluation, review/apply workflow, validation, and cross-script convergence still pass.
- `go test ./...` passes.
- Work is committed locally and not pushed.


## Next After This

After counsellor/transmutation fields:
- Integer-only harmonic data model coverage
- Extended script support (Arabic, Thai, etc.)
- Harmonic/audio rendering layer
- Consult `docs/ARCHITECTURE_READINESS.md` before major new feature layers

## Nice To Have Later

- integer-only harmonic data model coverage
- extended script support (Arabic, Thai, etc.)
- harmonic/audio rendering layer


