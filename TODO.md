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

Current git status before this handoff:

```text
## main...origin/main [ahead 34]
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
- Field precision improved without recall collapse:
  - train field precision `0.23 -> 0.39`, recall `1.00 -> 1.00`
  - held-out field precision `0.16 -> 0.27`, recall `0.62 -> 0.62`
  - train pass rate `0/8 -> 8/8`
  - held-out pass rate `0/8 -> 5/8`
- The next step is not counsellor prose and not new symbolic feature work. It is human review/apply for generated suggestions.

## Current Next Task

Task:
Add human review/apply workflow for proposed weight changes.

Context:
The evidence pipeline, harmonic core, training evaluation, validation, channel semantics, identity/provenance model, concept schema guardrails, ranking weights, held-out evaluation, false-activation gating, and review-file generation are now in place.

Socrates can now generate pending weight suggestions in `review/weight_suggestions.yaml`. The next step is controlled human review: humans should be able to mark suggestions accepted or rejected, and only accepted suggestions should be applied through an explicit command.

The goal is:

```text
pending suggestions
  -> human marks accepted/rejected
  -> explicit apply command reads accepted suggestions
  -> updates target config only after validation
  -> leaves an audit trail
```

Reference:

- `TRAINING_MODEL.md`
- `IMPLEMENTATION_ROADMAP.md`
- `ARCHITECTURE.md`

Likely files:

- `internal/decipher/scoring.go`
- `internal/decipher/types.go`
- `internal/decipher/*`
- `internal/training/*`
- `training/weights.yaml`
- `review/weight_suggestions.yaml`
- new file if useful: `internal/training/weight_review.go`
- new file if useful: `review/applied_weight_changes.yaml`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Run the documented validation command and record the result.
4. Define the review/apply workflow for weight suggestions.
   - Suggestions must have statuses: `pending`, `accepted`, `rejected`.
   - Only `accepted` suggestions may be applied.
   - `pending` and `rejected` suggestions must never change runtime config.
5. Add a parser/validator for review files.
   - Invalid statuses fail.
   - Missing target path/current/suggested values fail.
   - Target paths must be restricted to known ranking weight fields.
6. Add an explicit apply command if the scope stays small.
   - Example: `socrates train apply-weight-suggestions`
   - It must require an explicit review file path or use the documented default review path.
   - It must update only `training/weights.yaml` or the project’s chosen ranking weight config, not active knowledge YAML.
7. Before applying, verify the current value in the target config still matches the suggestion's recorded current value.
   - If it does not match, refuse to apply and report drift.
8. After applying, run or require validation of the updated weight config.
9. Write an audit record for applied suggestions.
   - Include timestamp if available, suggestion path, old value, new value, metrics delta, and status.
   - Keep rejected suggestions in the review file or audit trail.
10. Do not auto-accept suggestions.
11. Do not generate new suggestions in this task unless needed for tests.
12. Do not mutate active knowledge.
13. Do not add new active concepts, forms, relations, examples, AI, embeddings, counsellor/transmutation fields, UI, or audio in this task.
14. Keep train and held-out metrics separate in any reporting.
15. Preserve completed guardrails:
   - active YAML validates
   - `knowledge validate` works
   - channel diversity counts active evidence only
   - training evaluation is precision-aware
   - multi-concept fuzzy evidence remains complete
   - cross-script love and Hebrew `El` boundary behavior still works
   - concept schema hygiene warnings remain visible and bounded
16. Keep the implementation light. Prefer explicit YAML update logic over a broad config-management framework.
17. Add or update tests proving:
   - pending suggestions do not apply
   - rejected suggestions do not apply
   - accepted suggestions apply only to allowed weight fields
   - current-value drift prevents apply
   - audit record is written
   - active knowledge YAML is not mutated
18. Run targeted tests while working, then `go test ./...` or `make test`.
19. Commit the completed review/apply workflow locally with a clear message. Do not push.
20. Report final `git status --short --branch`, validation result, tests run, files changed, and an example accepted/rejected/apply flow.

Acceptance Criteria:

- Pending/rejected suggestions cannot mutate config.
- Accepted suggestions can be applied only through an explicit command.
- Apply checks current-value drift before writing.
- Applied changes leave an audit trail.
- Active knowledge YAML is not mutated.
- Train and held-out metrics remain separate and visible.
- Existing harmonic core, training evaluation, quality gate behavior, identity/provenance hardening, and cross-script convergence still pass.
- Training does not mutate active knowledge.
- No black-box truth model is introduced.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

After human review/apply for weight suggestions is stable, do repository shape and maintainability cleanup before counsellor/transmutation fields.

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
