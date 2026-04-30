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

Current git status before this handoff:

```text
## main...origin/main [ahead 26]
```

Architect correction:

- Socrates can now evaluate curated examples and can converge English, Norwegian, Hebrew, and Chinese love forms on a shared harmonic field.
- Pi completed the quality gate, including active-channel diversity, positional `knowledge validate`, precision-aware training evaluation, and multi-concept fuzzy evidence handling.
- Pi completed identity/provenance hardening: confidence is `verified`, `plausible`, `speculative`; provenance/source is `curated`, `traditional`, `human_review`, `physics`; `curated` is not confidence.
- Validation warning debt was reduced to zero.
- The next step is not counsellor prose, not new symbolic feature work, and not ranking-weight tuning yet.
- The remaining schema risk is concept hygiene: `id`, `name`, `aliases`, and broad spiritual equivalences must not collapse distinct fields.

## Current Next Task

Task:
Audit concept schema hygiene and guard broad spiritual identities.

Context:
Pi completed identity/provenance hardening in commit `3ef62c0`. Validation now reaches zero warnings. That is good, but clean validation does not prove that the concept model is spiritually or semantically clean.

Observed issue:

```yaml
- id: source
  name: source
  aliases: [source, origin, beginning, genesis, god, power, divine, totality]
```

This blurs canonical identity, display label, aliases, and neighboring concepts. It also risks collapsing distinct spiritual fields such as `source`, `god`, `divine`, and `power`.

Before ranking-weight tuning, audit the concept schema itself:

- Concept IDs are canonical identities.
- Names are human display labels.
- Aliases are alternate surface labels, not a place to store neighboring concepts.
- Broad spiritual terms may resonate through relations without being flattened into one alias bucket.

Reference:

- `docs/KNOWLEDGE_CURATION.md`
- `HARMONIC_DATA_MODEL.md`
- `ARCHITECTURE.md`

Likely files:

- `internal/knowledge/concepts.yaml`
- `internal/knowledge/validate.go`
- `internal/knowledge/validate_test.go`
- `internal/knowledge/knowledge.go`
- `docs/KNOWLEDGE_CURATION.md`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Run the documented validation command and record the result.
4. Audit `concepts.yaml` for concept hygiene issues:
   - `name` exactly equals `id`
   - alias equals own `id`
   - alias equals another concept's `id`
   - duplicate aliases across concepts
   - broad spiritual terms stored as aliases where relations would be more accurate
5. Add validator warnings or failures for the hygiene issues above.
   - Do not make seed-data convenience silently acceptable.
   - If `name == id` is temporarily accepted, it must be visible as a warning category until curated.
6. Clean the most important concept entries, especially broad spiritual fields:
   - `source`
   - `god`
   - `divine`
   - `power`
   - `one`
   - `light`
   - `truth`
   - `love`
7. Prefer relations over aliases when concepts are distinct but resonant.
   - Example: `source` may relate to `god`, `divine`, or `power`, but they should not automatically be identical aliases unless deliberately documented.
8. Do not do a broad concept-pack expansion in this task.
9. Preserve the completed quality-gate behavior:
   - active YAML validates
   - `knowledge validate` works
   - channel diversity counts active evidence only
   - training evaluation is precision-aware
   - multi-concept fuzzy evidence remains complete
   - cross-script love and Hebrew `El` boundary behavior still works
10. Keep the implementation light. Prefer validator guardrails and focused curation over large refactors.
11. Add or update tests proving:
   - alias cannot equal own concept ID without warning/failure
   - alias cannot equal another concept ID without warning/failure
   - duplicate aliases are detected
   - broad spiritual identities are represented through relations when distinct
   - validation remains compact and actionable
12. Run targeted tests while working, then `go test ./...` or `make test`.
13. Commit the completed concept schema hygiene work locally with a clear message. Do not push.
14. Report final `git status --short --branch`, validation result, concept hygiene warnings before/after, tests run, and files changed.

Acceptance Criteria:

- Concept `id`, `name`, `aliases`, and relations have distinct meanings.
- Own-ID aliases are detected and removed or explicitly warned.
- Cross-concept alias collisions are detected.
- Broad spiritual fields are not silently flattened into aliases.
- Validation remains compact and actionable.
- Existing harmonic core, training evaluation, quality gate behavior, and cross-script convergence still pass.
- Training does not mutate active knowledge.
- No black-box truth model is introduced.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

After concept schema hygiene is clean, return to configurable ranking weights and a small held-out evaluation split. Only after that should Socrates add counsellor/transmutation fields.

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
