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
## main...origin/main [ahead 50]
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
Add data-backed counsellor/transmutation fields.

Context:
Architecture readiness is complete. The engine is stable, documented, and ready for the next layer.

Counsellor/transmutation fields extend evidence-first reading with symbolic guidance capabilities. This is not fortune telling and not an oracle layer. Socrates should read the currently activated field, identify likely inertia if the field remains unchanged, and show data-backed correction/transmutation fields that may alter the path.

Target flow:

```text
input passage
  -> evidence-backed concept activation
  -> harmonic field
  -> detected tension / contraction / unresolved pattern
  -> data-backed transmutation relation
  -> possible correction field
  -> evidence-first counsel summary
```

Examples of future data-backed relations:

```text
resentment -> forgiveness
pride -> humility
fear -> trust
guilt -> mercy / atonement
past fixation -> presence / holy instant
separation -> loving consciousness
control -> surrender
judgment -> mercy
```

These are examples, not hardcoded Go behavior. If the active data does not support a path, Socrates must not invent it.

Rules:
- No hardcoded word-specific behaviour in Go
- All mappings must be data-driven with source/lens/confidence
- Training must not mutate active knowledge
- No black-box truth model
- No fixed-future prediction. Use language such as "if unchanged, this field tends toward..." rather than "this will happen."
- Preserve human sovereignty. Counsel should suggest possible correction fields, not command the user.
- No medical, legal, or clinical claims.
- No claim that symbolic readings are proven divine facts.
- Default output must remain concise and evidence-first; detailed counsellor internals belong in debug/report output.

Reference:

- `docs/ARCHITECTURE_READINESS.md`
- `docs/KNOWLEDGE_CURATION.md`
- `TRAINING_MODEL.md`
- `HARMONIC_DATA_MODEL.md`
- `FREQUENCY_MODEL.md`

Likely files:

- `internal/knowledge/*.yaml`
- `internal/knowledge/knowledge.go`
- `internal/knowledge/loader.go`
- `internal/knowledge/validate.go`
- `internal/decipher/types.go`
- `internal/decipher/engine.go`
- `internal/decipher/render.go`
- new file if useful: `internal/decipher/counsellor_field.go`
- new tests under `internal/decipher/`
- training examples if needed: `training/examples.yaml`, `training/heldout.yaml`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Run `./bin/socrates knowledge validate` or the documented equivalent and record the result.
4. Design the smallest data model for transmutation/counsellor relations.
   - Prefer extending existing relation/profile data if clean.
   - Add a new YAML file only if existing schemas become unclear.
   - Required fields: source field/concept, target correction field/concept, relation kind, confidence, source/provenance, lens, weight.
   - Suggested relation kinds: `transmutes_to`, `softens_through`, `corrects_through`, `releases_into`, `grounds_in`.
5. Validation must reject:
   - unknown source/target concept or field IDs
   - invalid confidence
   - invalid provenance/source
   - invalid weights
   - unsupported relation kinds
   - floating-point harmonic meaning data if any harmonic profile is added
6. Add a structured result to `Reading`, such as `CounsellorField` or `TransmutationField`.
   - It should include active tension/source concepts, suggested correction fields, relation paths, confidence, weights, and evidence paths.
   - It must be absent/empty when no data-backed transmutation path exists.
7. Build counsellor/transmutation suggestions from active passage fields and harmonic fields.
   - Do not inspect raw words directly for meaning.
   - Use activated concepts/fields and curated transmutation relations.
   - Require sufficient evidence strength before suggesting a correction field.
   - Prefer verified/plausible evidence over speculative evidence.
8. Rendering:
   - Default render: concise, evidence-first, only when counsellor fields exist.
   - Debug render: show relation kind, source/lens/confidence, weights, and evidence paths.
   - Avoid moralizing language.
   - Avoid "you must" phrasing.
9. Add minimal curated seed data for the first counsellor paths.
   - Keep it small and high-signal.
   - Use source/provenance such as `curated` or `traditional`.
   - Confidence should be `plausible` or `speculative` unless ordinary evidence supports `verified`.
   - Do not add a broad spiritual concept pack.
10. Add training/held-out examples only if needed to verify the feature.
   - Keep training separate from active knowledge.
   - Do not mutate active knowledge from training.
11. Preserve completed guardrails:
   - active YAML validates
   - `knowledge validate` works
   - channel diversity counts active evidence only
   - training evaluation is precision-aware
   - multi-concept fuzzy evidence remains complete
   - cross-script love and Hebrew `El` boundary behavior still works
   - concept schema hygiene warnings remain visible and bounded
   - review/apply workflow does not mutate active knowledge without explicit accepted suggestions
12. Add tests proving:
   - counsellor field appears only from data-backed transmutation relations
   - removing the transmutation data removes the counsel
   - weak/speculative evidence alone does not produce strong counsel
   - relation evidence paths are preserved
   - default render stays concise
   - debug render exposes source/lens/confidence/evidence paths
   - no production Go hardcodes examples like forgiveness, humility, ego, or love as special cases
13. Run targeted tests while working, then `go test ./...` or `make test`.
14. Commit the completed counsellor/transmutation field work locally with a clear message. Do not push.
15. Report final `git status --short --branch`, validation result, tests run, files changed, example default/debug outputs, and any added seed paths.

Acceptance Criteria:

- Counsellor/transmutation suggestions are data-backed.
- No hardcoded word-specific or concept-specific counsel exists in Go.
- Suggestions preserve evidence paths and confidence/provenance.
- Default output remains concise and evidence-first.
- Debug output exposes transmutation internals.
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

- counsellor/transmutation fields (current task)
- integer-only harmonic data model coverage
- extended script support (Arabic, Thai, etc.)
- harmonic/audio rendering layer



