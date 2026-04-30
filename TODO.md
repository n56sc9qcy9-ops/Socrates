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

Current git status before this handoff:

```text
## main...origin/main [ahead 36]
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
- Field precision improved without recall collapse:
  - train field precision `0.23 -> 0.39`, recall `1.00 -> 1.00`
  - held-out field precision `0.16 -> 0.27`, recall `0.62 -> 0.62`
  - train pass rate `0/8 -> 8/8`
  - held-out pass rate `0/8 -> 5/8`
- The next step is not counsellor prose and not new symbolic feature work. It is repository shape and knowledge layout cleanup.

## Current Next Task

Task:
Consolidate knowledge layout and clarify runtime/reference/review data boundaries.

Context:
The evidence pipeline, harmonic core, training evaluation, validation, channel semantics, identity/provenance model, concept schema guardrails, ranking weights, held-out evaluation, false-activation gating, review-file generation, and review/apply workflow are now in place.

The repository now has two knowledge-looking locations:

```text
knowledge/
  concepts.yaml
  forms.yaml
  relations.yaml

internal/knowledge/
  concepts.yaml
  forms.yaml
  relations.yaml
  glyphs.yaml
  frequencies.yaml
  loader.go
  ...
```

This is confusing. Future agents and humans may edit the wrong YAML. Socrates needs one clearly documented active runtime knowledge location and separate places for review, reference/source, and archive material.

The goal:

```text
active runtime knowledge: one authoritative location
review suggestions: separate from active runtime knowledge
reference/source material: separate from active runtime knowledge
archive: separate from active runtime knowledge
```

Reference:

- `ARCHITECTURE.md`
- `README.md`
- `docs/KNOWLEDGE_CURATION.md`
- `IMPLEMENTATION_ROADMAP.md`

Likely files:

- `README.md`
- `ARCHITECTURE.md`
- `docs/KNOWLEDGE_CURATION.md`
- `IMPLEMENTATION_ROADMAP.md`
- `internal/knowledge/loader.go`
- `internal/knowledge/*.yaml`
- `knowledge/*.yaml`
- `review/*`
- `training/*`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Run the documented validation command and record the result.
4. Decide and document the authoritative runtime knowledge location.
   - Preferred long-term shape: root `knowledge/` contains structured subfolders such as `runtime/`, `reference/`, `review/`, and `archive/`.
   - If Go embedding makes immediate root-runtime migration too large, keep `internal/knowledge/*.yaml` as active runtime for now, but move or rename root `knowledge/*.yaml` so it cannot be mistaken for runtime.
   - The final state of this task must have no ambiguous duplicate runtime-looking YAML.
5. Make the chosen active runtime knowledge location explicit in:
   - `README.md`
   - `docs/KNOWLEDGE_CURATION.md`
   - relevant architecture docs if needed
6. Ensure loader/validation behavior matches the documented active runtime location.
   - `knowledge validate` should validate active embedded/runtime knowledge by default.
   - `knowledge validate --dir <path>` should remain available for explicit alternate directories.
7. Separate non-runtime material.
   - Review suggestions belong under `review/`.
   - Training examples and weights belong under `training/`.
   - Reference/source/seed/archive material must not sit beside runtime YAML unless clearly named as non-runtime.
8. Do not silently delete useful source material.
   - If root `knowledge/*.yaml` is obsolete, move it to a clearly named archive/reference location or document why it is removed.
   - Preserve git history naturally through the move.
9. Add guardrails so future agents do not edit the wrong YAML.
   - Documentation must say where active knowledge lives.
   - If practical, add a lightweight test or validation check proving the active knowledge files exist where documented.
10. Do not change knowledge semantics in this task except what is necessary for layout.
11. Do not add new active concepts, forms, relations, examples, AI, embeddings, counsellor/transmutation fields, UI, or audio in this task.
12. Preserve completed guardrails:
   - active YAML validates
   - `knowledge validate` works
   - channel diversity counts active evidence only
   - training evaluation is precision-aware
   - multi-concept fuzzy evidence remains complete
   - cross-script love and Hebrew `El` boundary behavior still works
   - concept schema hygiene warnings remain visible and bounded
13. Keep implementation light. Prefer moving/renaming data and updating docs over a loader rewrite unless the loader rewrite is clearly simpler and safer.
14. Add or update tests proving:
   - documented active knowledge location is valid
   - default embedded validation still validates active runtime knowledge
   - explicit `--dir` validation still works
   - review/training/reference material is not accidentally treated as active runtime knowledge
15. Run targeted tests while working, then `go test ./...` or `make test`.
16. Commit the completed repository/knowledge layout cleanup locally with a clear message. Do not push.
17. Report final `git status --short --branch`, validation result, tests run, files changed/moved, and the final runtime/reference/review/archive layout.

Acceptance Criteria:

- There is exactly one documented active runtime knowledge location.
- Root `knowledge/` and `internal/knowledge/` are no longer ambiguous duplicates.
- Review, training, reference/source, and archive data are clearly separated from active runtime knowledge.
- Loader and validation behavior match the documentation.
- Useful source/reference material is preserved or intentionally removed with rationale.
- Existing harmonic core, training evaluation, quality gate behavior, identity/provenance hardening, and cross-script convergence still pass.
- Training does not mutate active knowledge.
- No black-box truth model is introduced.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

After repository shape and knowledge layout are clean, split oversized tests by behavior. Only after maintainability cleanup should Socrates add counsellor/transmutation fields.

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
