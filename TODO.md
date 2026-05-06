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
## main...origin/main [ahead 66]
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
**Harmonic data model hardening — COMPLETED this session.**

Architect review status:
Not accepted yet. Pi addressed label ranges and improved documentation/tests, but one core blocker remains unresolved.

Current blocker:
- Unknown fields in `internal/knowledge/frequencies.yaml` are still silently ignored by `yaml.Unmarshal`. Documentation now states this is a known risk, but that is not sufficient for the meaning-frequency substrate.
- Active frequency YAML must reject unknown fields, at least for `frequencies.yaml`. Use strict YAML decoding for this file or add an explicit unknown-field scanner before typed conversion.
- Add a real loader-level regression proving that `frequency_hz: 528.0`, `pitch: 432.0`, `color_rgb: [1.0, 0.2, 0.3]`, `em_band_id: em.visible.green`, or any other unrecognized harmonic identity field under a frequency profile causes loading or validation to fail.
- Keep current valid `frequencies.yaml` loading successfully. If strict decoding affects other knowledge files, limit the change to frequency loading first.
- Update `HARMONIC_DATA_MODEL.md` so it no longer says unknown active frequency YAML fields are accepted as a standing risk. The accepted rule is: unknown active frequency identity fields are rejected until explicitly modeled.
- Commit the correction locally and report clean status, tests, validation, and the strict-decoding decision.

Rules (unchanged, enforced by this hardening):
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
