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
Audit and harden integer-only harmonic data model coverage.

Completed and accepted:
- Counsellor/transmutation calibration is complete.
- Direct evidence is ranked before propagated/speculative suggestions.
- Neighbor-derived evidence remains visibly propagated in debug output.
- Transmutation weights are normalized from integer percentages before scoring.
- Default counsel is concise and humble.
- Work committed locally through `ed68780`.

Context:
The counsellor layer is now restrained enough to move back to the harmonic core. Before Socrates grows toward audio, light, electromagnetic correspondence, or richer guidance, the meaning-frequency substrate must be audited and hardened.

The core rule remains: meaning-frequency identity must be integer based and concise. Decimal runtime scores are acceptable for evidence strength and ranking, but the stored harmonic meaning model must not become floating-point frequency numerology.

This task is an architecture hardening pass, not a feature expansion. Do not add audio synthesis, color rendering, chakras, UI, or a broad concept pack.

Rules:
- Do not store floating-point frequency values as meaning.
- Integer harmonic fields may include integer vectors, integer ratios, Pythagorean triples, Phi integer sequences/approximants, geometry IDs, and integer-backed electromagnetic references.
- Physical electromagnetic references must be clearly separated from symbolic correspondences.
- No hardcoded concept-to-frequency, concept-to-color, concept-to-note, pitch, chakra, electromagnetic, or music mappings in production Go.
- All harmonic mappings must live in curated data with source/lens/confidence.
- Runtime evidence scores may remain decimal because they are ranking signals, not meaning identity.
- Keep the model compact. Prefer a small schema with clear validation over many optional fields.
- Do not mutate active knowledge from training or review artifacts.

Reference:

- `docs/ARCHITECTURE_READINESS.md`
- `docs/KNOWLEDGE_CURATION.md`
- `TRAINING_MODEL.md`
- `HARMONIC_DATA_MODEL.md`
- `FREQUENCY_MODEL.md`

Likely files:

- `internal/knowledge/frequencies.yaml`
- `internal/knowledge/knowledge.go`
- `internal/knowledge/loader.go`
- `internal/knowledge/validate.go`
- `internal/decipher/harmonic_field.go`
- `internal/decipher/harmonic_field_test.go`
- `HARMONIC_DATA_MODEL.md`
- `FREQUENCY_MODEL.md`
- `docs/KNOWLEDGE_CURATION.md`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Run `./bin/socrates knowledge validate` or the documented equivalent and record the result.
4. Audit the active harmonic/frequency schema and data.
   - Identify every stored harmonic identity field.
   - Confirm which fields are integer identity data and which fields are runtime scores.
   - Confirm no meaning-frequency identity depends on a float.
5. Harden validation.
   - Reject floating-point harmonic meaning fields where integer identity is required.
   - Reject malformed integer ratios, vectors, triples, geometry IDs, or electromagnetic references.
   - Reject ambiguous physical/symbolic electromagnetic entries.
   - Keep accepted warning debt visible; do not silence unrelated warnings.
6. Harden loader/types only where necessary.
   - If fields are currently loose strings, make validation stricter before adding abstractions.
   - Do not create a large generalized physics model.
   - Keep backward compatibility with current valid YAML where possible.
7. Add focused tests proving:
   - active `frequencies.yaml` validates
   - float harmonic identity data is rejected
   - integer ratios/vectors/triples are accepted
   - invalid ratios/vectors/triples are rejected
   - physical electromagnetic references require a clear physical/source lens
   - symbolic electromagnetic correspondences remain distinguishable from physical claims
   - no production Go contains hardcoded concept-to-frequency/color/note/EM/chakra mappings
8. Review documentation.
   - Update `HARMONIC_DATA_MODEL.md` and/or `FREQUENCY_MODEL.md` if the current schema rules are unclear.
   - Document the boundary between integer meaning identity and decimal runtime evidence scoring.
   - Document the boundary between physical electromagnetic measurement and symbolic correspondence.
9. Preserve completed guardrails:
   - active YAML validates
   - `knowledge validate` works
   - channel diversity counts active evidence only
   - training evaluation is precision-aware
   - multi-concept fuzzy evidence remains complete
   - cross-script love and Hebrew `El` boundary behavior still works
   - concept schema hygiene warnings remain visible and bounded
   - review/apply workflow does not mutate active knowledge without explicit accepted suggestions
   - counsellor/transmutation direct-vs-propagated ranking remains correct
10. Run targeted tests while working, then `go test ./...` or `make test`.
11. Commit the completed harmonic data model hardening work locally with a clear message. Do not push.
12. Report final `git status --short --branch`, validation result, tests run, files changed, schema decisions, and any unresolved warnings.

Acceptance Criteria:

- Harmonic meaning identity remains integer-only.
- Decimal values are allowed only for runtime evidence/scoring, not stored meaning identity.
- Physical electromagnetic references are clearly separated from symbolic correspondence.
- Validation rejects malformed harmonic identity data.
- No concept-to-frequency/color/note/EM/chakra mapping is hardcoded in production Go.
- Documentation clearly explains the schema boundary.
- Existing harmonic core, training evaluation, review/apply workflow, validation, and cross-script convergence still pass.
- `go test ./...` passes.
- Work is committed locally and not pushed.


## Next After This

After harmonic data model hardening:
- Extended script support (Arabic, Thai, etc.)
- Harmonic/audio rendering layer
- Consult `docs/ARCHITECTURE_READINESS.md` before major new feature layers

## Nice To Have Later

- extended script support (Arabic, Thai, etc.)
- harmonic/audio rendering layer
