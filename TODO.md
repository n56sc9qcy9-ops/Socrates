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

Current git status before this handoff:

```text
## main...origin/main [ahead 30]
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
- The new metrics are honest and useful: train/held-out examples currently pass `0` because precision is too low.
- Current combined metrics: concept precision `0.18`, concept recall `0.79`, field precision `0.19`, field recall `0.81`.
- The next step is not counsellor prose and not new symbolic feature work. It is false-activation reduction.

## Current Next Task

Task:
Reduce false activations through evidence gating and ranking calibration.

Context:
The evidence pipeline, harmonic core, training evaluation, validation, channel semantics, identity/provenance model, concept schema guardrails, ranking weights, and held-out evaluation are now in place.

The held-out evaluation is exposing the real architectural problem: Socrates finds many expected fields, but activates too many extra fields.

```text
Train:   field precision 0.23, field recall 1.00
HeldOut: field precision 0.16, field recall 0.62
Combined field precision 0.19, field recall 0.81
```

The next task is to reduce false activations without destroying recall. Focus on evidence gates, thresholds, and ranking calibration. Do not add new concepts or broaden the knowledge base to make examples pass.

Reference:

- `TRAINING_MODEL.md`
- `IMPLEMENTATION_ROADMAP.md`
- `ARCHITECTURE.md`

Likely files:

- `internal/decipher/scoring.go`
- `internal/decipher/types.go`
- `internal/decipher/*`
- `internal/training/*`
- `training/examples.yaml`
- `training/weights.yaml`
- `training/heldout.yaml`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Run the documented validation command and record the result.
4. Capture baseline train and held-out metrics before changing behavior.
5. Produce a concise false-activation analysis:
   - top false activated concepts/fields by frequency
   - channels or evidence paths most responsible for false activations
   - whether false activations come from fuzzy matches, graph propagation, weak forms, harmonic profile fan-out, or duplicate evidence
6. Add evidence gates or ranking calibration to reduce false activations.
   - Prefer gating weak fuzzy evidence and speculative/structural channels before lowering strong exact evidence.
   - Prefer limiting graph propagation fan-out/depth before removing valid relations.
   - Prefer making weak evidence require independent support before it can activate harmonic fields.
   - Preserve cross-script exact/script evidence.
7. Tune only through explicit ranking weights or small bounded scoring changes.
   - Do not hardcode behavior for specific words.
   - Do not add new active knowledge to make metrics look better.
   - Do not mutate training examples to hide failures.
8. Establish an improvement target.
   - Primary: increase held-out field precision.
   - Secondary: increase combined field precision.
   - Guardrail: do not collapse field recall below a documented acceptable floor.
   - If no safe improvement is possible, report why with evidence-path analysis.
9. Keep train and held-out metrics separate.
10. Keep examples separate from active knowledge.
11. Do not let training mutate active knowledge.
12. Do not write suggested weights back into active YAML in this task.
13. Do not add AI, embeddings, counsellor/transmutation fields, UI, or audio in this task.
14. Preserve completed guardrails:
   - active YAML validates
   - `knowledge validate` works
   - channel diversity counts active evidence only
   - training evaluation is precision-aware
   - multi-concept fuzzy evidence remains complete
   - cross-script love and Hebrew `El` boundary behavior still works
   - concept schema hygiene warnings remain visible and bounded
15. Keep the implementation light. Prefer explicit thresholds, gates, and reports over a general ML framework.
16. Add or update tests proving:
   - weak fuzzy/speculative evidence cannot dominate exact/verified evidence
   - unsupported weak evidence does not activate harmonic fields alone
   - graph propagation remains bounded and does not fan out into unrelated fields
   - train and held-out metrics are reported separately
   - precision-aware failure still fails noisy runs
   - cross-script love and Hebrew `El` boundary behavior still works
17. Run targeted tests while working, then `go test ./...` or `make test`.
18. Commit the completed false-activation reduction work locally with a clear message. Do not push.
19. Report final `git status --short --branch`, validation result, tests run, files changed, baseline metrics, final metrics, and a concise explanation of what false activations were reduced.

Acceptance Criteria:

- False activations are analyzed by evidence path/channel, not guessed.
- Held-out field precision improves or Pi provides a defensible no-safe-improvement explanation.
- Recall loss, if any, is measured and justified.
- Weak fuzzy/speculative evidence is gated more strictly than exact/verified evidence.
- Train and held-out metrics remain separate and visible.
- Existing harmonic core, training evaluation, quality gate behavior, identity/provenance hardening, and cross-script convergence still pass.
- Training does not mutate active knowledge.
- No black-box truth model is introduced.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

After false activations are reduced, add review-file generation for proposed weight changes. Only after that should Socrates add counsellor/transmutation fields.

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
