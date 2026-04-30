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

Current git status before this handoff:

```text
## main...origin/main [ahead 24]
```

Architect correction:

- Socrates can now evaluate curated examples and can converge English, Norwegian, Hebrew, and Chinese love forms on a shared harmonic field.
- Pi completed the quality gate, including active-channel diversity, positional `knowledge validate`, precision-aware training evaluation, and multi-concept fuzzy evidence handling.
- Pi began identity/provenance hardening, but the reported validation output still shows `239 warning(s)`.
- A warning count that high means the knowledge layer is still noisy. It is not an acceptable resting state before ranking weights or counsellor work.
- The next step is not counsellor prose, not new symbolic feature work, and not ranking-weight tuning yet.
- Concept IDs, aliases, confidence, and provenance must be kept distinct before more knowledge is added.
- Aliases may help resolve references, but active data should still preserve canonical concept identity and evidence paths.
- `curated` sounds like provenance/source, not confidence, unless the project explicitly defines it as a confidence level with semantics distinct from `verified`, `plausible`, and `speculative`.

## Current Next Task

Task:
Finish knowledge identity hardening and reduce validation warning debt.

Context:
Pi completed the architecture quality gate in commit `b6051fb` and began identity/provenance hardening. The reported implementation added `Source` fields, structured concept-reference resolution, valid source values, and `BuildIndexes()`. That direction is good.

However, the reported validator output was:

```text
Validation passed with warnings
239 warning(s)
```

That is too noisy. A valid knowledge layer should not require humans to ignore hundreds of warnings. The next task is to finish this hardening by making validation output actionable and compact.

Before ranking-weight tuning, lock down the knowledge semantics:

- Concept IDs are canonical identities.
- Aliases are surface names or alternate labels that resolve to canonical concept IDs.
- Forms/script words/glyphs are evidence surfaces.
- Confidence describes truth/support level, such as `verified`, `plausible`, or `speculative`.
- Provenance/source/lens describes where the mapping came from, such as `curated`, `human_review`, `traditional`, or `physics`.

Reference:

- `docs/KNOWLEDGE_CURATION.md`
- `HARMONIC_DATA_MODEL.md`
- `TRAINING_MODEL.md`

Likely files:

- `internal/knowledge/*.yaml`
- `internal/knowledge/validate.go`
- `internal/knowledge/validate_test.go`
- `internal/knowledge/knowledge.go`
- `internal/knowledge/loader.go`
- `internal/decipher/*`
- `internal/training/*`
- `training/examples.yaml`
- `docs/KNOWLEDGE_CURATION.md`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Run the documented validation command and record it:
   - `make test` if using the new Makefile
   - `./bin/socrates knowledge validate` if `make create` was used
   - `./socrates knowledge validate` only if the repo-root binary has been rebuilt
4. Report the current warning categories and counts before changing behavior.
   - Group warnings by cause, such as unknown concept, alias resolution, duplicate form, missing source, invalid confidence, etc.
   - Do not paste hundreds of repeated warning lines into the final report.
5. Decide and document the confidence vocabulary.
   - Preferred: keep confidence to `verified`, `plausible`, `speculative`.
   - Treat `curated` as `source`, `provenance`, or `review_status`, not confidence.
   - If `curated` remains a confidence value, define exactly what it means and how it ranks against the others.
6. Decide and document alias-reference behavior.
   - Active runtime identities should resolve to canonical concept IDs.
   - If YAML references an alias, validation should either normalize it or report the canonical target clearly.
   - Evidence paths should show both the surface/alias and the canonical concept when useful.
7. Update validation so it distinguishes:
   - direct canonical concept ID reference
   - alias reference resolved to a canonical concept ID
   - invalid unknown reference
8. Reduce validation warning debt.
   - Unknown concept warnings should be fixed or intentionally marked external/speculative according to the project rules.
   - Alias-resolution warnings should not flood normal validation if the alias resolves cleanly and this behavior is accepted.
   - Duplicate form warnings should be deduplicated, justified, or made more specific so intentional multi-concept mappings are not treated as noise.
   - The target is zero warnings if practical. If not practical, the remaining warnings must be few, categorized, and intentionally accepted.
9. Ensure meaning-frequency profiles, relations, and training expectations use canonical concept IDs unless there is a deliberate alias-resolution test.
10. Keep cross-linguistic forms like `ruach`, `prana`, `qi`, `ahava`, `pneuma`, `psyche`, `dao`, `de`, `ren`, `xin`, `echad`, and `ananda` as forms/aliases that resolve cleanly to canonical concepts.
11. Do not add broad new concept packs in this task. Only adjust data needed to clarify identity/provenance semantics and warning debt.
12. Preserve the completed quality-gate behavior:
   - active YAML validates
   - `knowledge validate` works
   - channel diversity counts active evidence only
   - training evaluation is precision-aware
   - multi-concept fuzzy evidence remains complete
   - cross-script love and Hebrew `El` boundary behavior still works
13. Keep the implementation light. Prefer small schema/validation clarifications over new architecture.
14. Add or update tests proving:
   - `curated` is not silently treated as confidence unless explicitly documented and tested
   - aliases resolve to canonical concept IDs without hiding unknown references
   - validator error/warning messages distinguish alias resolution from invalid references
   - frequency profiles and relations keep canonical IDs
   - training examples validate against canonical IDs or explicit alias-resolution cases
   - validation warning count is zero or intentionally bounded and categorized
15. Run targeted tests while working, then `go test ./...` or `make test`.
16. Commit the completed identity/provenance hardening locally with a clear message. Do not push.
17. Report final `git status --short --branch`, validation result, warning category counts before/after, tests run, files changed, and the final confidence/provenance decision.

Acceptance Criteria:

- Confidence and provenance are not blurred.
- Alias references resolve to canonical concept IDs in a visible, test-protected way.
- Unknown concept references still fail validation.
- Validation output is compact and actionable; hundreds of warnings are not acceptable.
- Frequency profiles, relations, and training examples remain identity-clean.
- Existing harmonic core, training evaluation, quality gate behavior, and cross-script convergence still pass.
- Training does not mutate active knowledge.
- No black-box truth model is introduced.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

After identity/provenance semantics are clean, return to configurable ranking weights and a small held-out evaluation split. Only after that should Socrates add counsellor/transmutation fields.

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
