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

Completed and committed locally:

- Phase 3A cleanup.
- Activation-Energy Discipline Gate.
- Phase 5/6 graph and passage-field work.
- First knowledge validation/suggestion pipeline.
- Data-backed harmonic frequency profiles and harmonic field scoring.
- Training/evaluation layer for supervised ranking over evidence paths.
- Channel semantics correction and cross-script meaning-frequency convergence.
- Architecture quality gate for validation, scoring, evaluation, and evidence-path integrity.
- Knowledge identity/provenance hardening; validation warnings reduced to zero.
- Concept schema hygiene; alias collisions guarded and broad spiritual fields separated.
- Configurable ranking weights and train/held-out evaluation split.
- False activation reduction through harmonic evidence gating and ranking calibration.
- Review-file generation for proposed ranking weight changes.
- Human review/apply workflow for proposed ranking weight changes.
- Knowledge layout consolidation; active runtime in `internal/knowledge/`, reference in `knowledge/reference/`.
- Oversized decipher tests split by behavior.
- Production responsibility audit completed; no safe extraction needed.
- CLI ergonomics cleanup; root natural input now defaults to decipher and help/build docs are clearer.

Current git status before this handoff:

```text
## main...origin/main [ahead 43]
?? internal/training/review/
```

Architect correction:

- Socrates can now evaluate curated examples and can converge English, Norwegian, Hebrew, and Chinese love forms on a shared harmonic field.
- Pi completed the quality gate, including active-channel diversity, positional `knowledge validate`, precision-aware training evaluation, and multi-concept fuzzy evidence handling.
- Pi completed identity/provenance hardening: confidence is `verified`, `plausible`, `speculative`; provenance/source is `curated`, `traditional`, `human_review`, `physics`; `curated` is not confidence.
- Validation warning debt was reduced to zero.
- Pi completed concept schema hygiene in commit `0a4661d`.
- Remaining validation warnings are now a single explicit seed-label category: `name == id` for uncurated display labels.
- That seed-label warning category is accepted as visible curation debt for now; it must not be hidden, but it no longer blocks ranking work.
- Pi completed ranking weights and held-out evaluation in commit `b66e47a`.
- Pi completed false activation reduction in commit `0030020`.
- Pi completed review-gated weight suggestion generation in commit `4b50136`.
- Pi completed human review/apply workflow for proposed weight changes in commit `c420013`.
- Pi completed knowledge layout consolidation in commit `916c504`.
- Pi completed behavior-based test split in commit `66a731a`.
- Pi completed production responsibility audit. `engine.go`, `activation_graph.go`, and `candidate_generation.go` have clear responsibilities and no low-risk extraction seam was found.
- Pi completed CLI ergonomics in commit `007922a`.
- Active runtime knowledge is documented as `internal/knowledge/`.
- Root `knowledge/` now holds non-runtime reference material under `knowledge/reference/`.
- Field precision improved without recall collapse:
  - train field precision `0.23 -> 0.39`, recall `1.00 -> 1.00`
  - held-out field precision `0.16 -> 0.27`, recall `0.62 -> 0.62`
  - train pass rate `0/8 -> 8/8`
  - held-out pass rate `0/8 -> 5/8`
- The next step is not counsellor prose and not new symbolic feature work. It is a final architecture readiness and documentation consolidation pass.

## Current Next Task

Task:
Perform architecture readiness review and documentation consolidation.

Context:
The core evidence pipeline, harmonic field layer, training/evaluation loop, review-gated weight workflow, knowledge layout, maintainability split, production responsibility audit, and CLI ergonomics are now in place.

Before Socrates moves into counsellor/transmutation fields, the project needs a long-horizon readiness pass:

- make sure docs match the code
- make sure TODO does not carry stale completed addenda as active guidance
- make sure all guardrails are documented in the right permanent files
- make sure generated artifacts are not sitting in source directories
- make sure the current CLI, validation, training, held-out evaluation, suggestion generation, and apply workflow all work as documented
- produce a concise readiness report with blockers and recommended next task

Reference:

- `README.md`
- `ARCHITECTURE.md`
- `FREQUENCY_MODEL.md`
- `HARMONIC_DATA_MODEL.md`
- `TRAINING_MODEL.md`
- `IMPLEMENTATION_ROADMAP.md`
- `docs/KNOWLEDGE_CURATION.md`
- `CONTRIBUTING.md`
- `TODO.md`

Likely files:

- `README.md`
- `ARCHITECTURE.md`
- `FREQUENCY_MODEL.md`
- `HARMONIC_DATA_MODEL.md`
- `TRAINING_MODEL.md`
- `IMPLEMENTATION_ROADMAP.md`
- `docs/KNOWLEDGE_CURATION.md`
- `CONTRIBUTING.md`
- `TODO.md`
- new file if useful: `docs/ARCHITECTURE_READINESS.md`
- generated/test artifact path: `internal/training/review/`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Run and record the documented operational commands:
   - `make test` or `go test ./...`
   - `make create`
   - `./bin/socrates help`
   - `./bin/socrates love`
   - `./bin/socrates knowledge validate`
   - `./bin/socrates train --heldout training/heldout.yaml`
   - `./bin/socrates suggest-weights`
   - `./bin/socrates apply-weights` only in a safe test flow or with a dry-run/no-op path if available; do not accidentally apply pending real suggestions
4. Resolve the untracked `internal/training/review/` artifact.
   - If it is generated test output, remove it or move tests to use `t.TempDir()`.
   - If it is an intended fixture, move it to a proper testdata path and document why.
   - Do not leave review/generated YAML under `internal/training/review/`.
5. Audit all major docs for drift against current behavior.
   - The active runtime knowledge location must be `internal/knowledge/`.
   - Reference knowledge must be described as `knowledge/reference/`.
   - Review files must be described as review data, not active runtime knowledge.
   - Training examples/weights/held-out data must be separate from active knowledge.
   - `curated` must be provenance/source, not confidence.
   - Confidence must remain `verified`, `plausible`, `speculative`.
   - Root natural-language CLI input must be documented.
   - Build path `bin/socrates` must be documented.
6. Move durable architectural guidance out of transient TODO addenda if appropriate.
   - Concept schema guardrails belong in curation/architecture docs.
   - Repository shape/maintainability guidance belongs in architecture/contributing docs.
   - TODO should remain the current task source, not a permanent archive of completed addenda.
7. Create or update a readiness report.
   - Suggested file: `docs/ARCHITECTURE_READINESS.md`.
   - Include current architecture status, verified commands, known metrics, accepted warning debt, remaining risks, and blockers before counsellor/transmutation fields.
   - Include the recommended next task.
8. Do not start counsellor/transmutation implementation in this task.
9. Do not add new active knowledge, concepts, forms, relations, AI, embeddings, UI, or audio.
10. Do not change engine behavior, scoring, ranking, validation, training thresholds, or YAML semantics except to fix a documented mismatch discovered during readiness review.
11. Preserve completed guardrails:
   - active YAML validates
   - `knowledge validate` works
   - channel diversity counts active evidence only
   - training evaluation is precision-aware
   - multi-concept fuzzy evidence remains complete
   - cross-script love and Hebrew `El` boundary behavior still works
   - concept schema hygiene warnings remain visible and bounded
12. Add or update tests only if they protect a documented readiness/cleanup issue.
13. Run targeted tests while working, then `go test ./...` or `make test`.
14. Commit the completed readiness/documentation consolidation locally with a clear message. Do not push.
15. Report final `git status --short --branch`, commands run, tests run, files changed, readiness conclusion, and recommended next task.

Acceptance Criteria:

- Documentation matches current code and data layout.
- `TODO.md` is focused on active/future task flow, not stale completed addenda.
- Generated/test artifacts are not left in source directories.
- Readiness report exists and identifies blockers or confirms readiness for the next architecture phase.
- Operational commands listed above are verified or exceptions are documented.
- Engine behavior is unchanged.
- Existing harmonic core, training evaluation, quality gate behavior, identity/provenance hardening, and cross-script convergence still pass.
- Training does not mutate active knowledge.
- No black-box truth model is introduced.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

After architecture readiness is documented, decide whether to add counsellor/transmutation fields or handle any blockers identified in the readiness report.

## Nice To Have Later

- CLI ergonomics:
  - Let root text default to decipher, e.g. `socrates what is the purpose of my life?`.
  - Join multi-argument root input into one passage so quotes are optional for ordinary questions.
  - Keep `socrates decipher <text>` as the explicit form.
  - Make root help show subcommand flags clearly, including debug/training flags.
  - Consider hiding or eventually removing the typo alias `descifer` after compatibility is no longer useful.
  - Keep `socrates knowledge validate` working exactly as documented.
  - Make build/run paths obvious: `make create` currently builds `bin/socrates`, while some manual testing has used repo-root `./socrates`.

## Architect Addendum: Concept Schema Guardrails

The identity-hardening task must also address concept schema hygiene.

Problem observed in `concepts.yaml`:

```yaml
- id: source
  name: source
  aliases: [source, origin, beginning, genesis, god, power, divine, totality]
```

This blurs machine identity, display name, aliases, and neighboring concepts.

Required guardrails:

- `id` is the canonical machine identity.
- `name` is the human display label and should not merely be an uncurated copy of `id` except during explicitly temporary seed data.
- `aliases` are alternate surface labels that resolve to the canonical `id`; they should not include the concept's own `id`.
- `aliases` should not include another concept's canonical `id`.
- alias collisions across concepts must fail validation or be explicitly disambiguated.
- broad spiritual concepts such as `source`, `god`, `divine`, and `power` should not be collapsed into aliases unless the curation decision is deliberate and documented; usually they should be separate concepts connected by relations.

Validation should warn or fail for:

- alias equals own concept ID
- alias equals another concept ID
- duplicate alias across concepts
- `name` exactly equals `id` across many entries, indicating uncurated labels

This is not cosmetic. Poor identity hygiene will poison resonance, training, and future counsellor guidance.

## Architect Addendum: Repository Shape And Maintainability

This is not the active task while Pi is working. It is a future architecture debt item that must be handled before major new feature layers.

Problem:

- Some production files have too many responsibilities.
- Some test files are too large for future agents to reason about safely.
- Runtime YAML and reference YAML exist in multiple locations, which risks edits landing in the wrong knowledge set.

Required future guardrails:

- Exactly one knowledge location must be authoritative for runtime behavior.
- Non-runtime knowledge must be clearly labeled as reference, review, archive, seed, or source material.
- Docs must state which YAML is embedded and active.
- Large tests should be split by behavior, not mechanically by line count.
- Production files should be split only after behavior is stable, and only along real ownership boundaries.

Known maintainability targets:

- `internal/decipher/engine.go`: keep orchestration here; move settled scoring, passage, harmonic, or rendering-prep responsibilities out when stable.
- `internal/decipher/activation_graph.go`: separate graph construction, propagation, evidence paths, and deduplication when stable.
- `internal/decipher/candidate_generation.go`: separate normalization, phonetic, n-gram, edit-variant, and skeleton generation when stable.
- `internal/decipher/engine_test.go`: split into render, scoring, bounds, regression, and integration suites.
- `internal/decipher/activation_graph_test.go`: split by graph construction, propagation, deduplication, and path behavior.
- `internal/decipher/passage_field_test.go`: split by merge, convergence, and regression behavior.

Knowledge layout target:

```text
active runtime knowledge: one clearly documented location
review suggestions: separate from active runtime knowledge
reference/source material: separate from active runtime knowledge
archive: separate from active runtime knowledge
```

This cleanup must not change resonance behavior by accident. It should be done as a dedicated architecture-maintenance task with before/after tests.
