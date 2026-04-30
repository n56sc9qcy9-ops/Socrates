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

Current git status before this handoff:

```text
## main...origin/main [ahead 22]
```

Architect correction:

- Socrates can now evaluate curated examples and can converge English, Norwegian, Hebrew, and Chinese love forms on a shared harmonic field.
- The next step is not counsellor prose and not new symbolic feature work.
- External review found several architecture guardrail issues that must be fixed before ranking-weight tuning.
- The priority is validity and evidence integrity: active YAML must validate, channels must not inflate scores, training must fail on excessive false activations, and multi-concept evidence paths must remain complete.

## Current Next Task

Task:
Architecture quality gate: fix validation, scoring, evaluation, and evidence-path integrity.

Context:
Pi completed the channel semantics and cross-script convergence task in commit `240c958`. Before adding ranking weights, counsellor fields, or more symbolic knowledge, the existing architecture must pass its own integrity gates.

Known review findings to address:

- P1: Active embedded knowledge must pass `knowledge validate`. Current YAML has references such as `living`, `vitality`, `prana`, or `qi` being treated as concept IDs when they may be aliases or undefined concepts.
- P2: Glyph pattern kinds must stay separated before matching. Bigrams, prefixes, suffixes, and structural observations must not collapse into one flat matcher that activates arbitrary substrings.
- P2: Channel diversity must be based on channels with active meaningful evidence, not on configured channels that emitted nothing.
- P2: Training evaluation must fail or clearly fail quality gates on excessive false activations, not only on low recall.
- P2: Fuzzy evidence for multi-mapped anchor forms must preserve all mapped concepts, not only the first form mapping.
- P3: The documented validation command must work as advertised: `socrates knowledge validate`.

Reference:

- `TRAINING_MODEL.md`
- `docs/KNOWLEDGE_CURATION.md`
- `HARMONIC_DATA_MODEL.md`

Likely files:

- `cmd/socrates/main.go`
- `internal/decipher/activation_graph.go`
- `internal/decipher/channel_glyph.go`
- `internal/decipher/glyphs.go`
- `internal/decipher/scoring.go`
- `internal/decipher/*_test.go`
- `internal/knowledge/*.yaml`
- `internal/knowledge/validate.go`
- `internal/knowledge/validate_test.go`
- `internal/training/*`
- `training/examples.yaml`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Run the documented validation command before changing code and record the result:
   - `./socrates knowledge validate`
   - If the binary is stale, also verify through the Go command path Pi normally uses.
4. Fix active embedded YAML so it passes validation without weakening validation rules.
   - Prefer adding or correcting compact concept IDs/forms/relations over teaching the validator to ignore bad data.
   - If a value is an alias, map it through a form or alias field; do not pretend it is a concept ID.
5. Fix the documented CLI so `socrates knowledge validate` works. `--validate` may remain as compatibility, but the documented positional command must validate.
6. Preserve glyph pattern kind boundaries.
   - Bigrams should match only as bigrams.
   - Prefixes should match only at the start.
   - Suffixes should match only at the end.
   - Structural/orthographic observations must not be flattened into semantic glyph patterns.
7. Fix channel diversity scoring so only channels with active meaningful signals can contribute to diversity.
   - Empty configured channels must not increase confidence.
   - Weak/unknown/no-op signals should not count as full independent evidence.
8. Fix training evaluation so excessive false activations can fail an example or fail an aggregate quality gate.
   - Recall-only passing is not acceptable.
   - Precision thresholds should be explicit and documented in code/tests.
   - Keep examples separate from active knowledge.
9. Fix fuzzy evidence path construction for anchor forms with multiple concept mappings.
   - If a matched anchor form maps to multiple concepts, preserve all valid mapped concepts and evidence paths.
   - Do not stop at the first matching form.
10. Do not add new symbolic concepts, counsellor/transmutation fields, ranking-weight systems, AI, embeddings, UI, or audio in this task.
11. Keep the implementation light. Prefer small focused fixes over new abstractions unless they remove real duplication or protect an architecture boundary.
12. Add or update tests proving:
   - active embedded knowledge validates successfully
   - `socrates knowledge validate` runs the validator
   - glyph bigram/prefix/suffix/structure patterns do not cross-match outside their kind
   - empty channels do not earn channel diversity
   - training/evaluation fails or flags excessive false activations
   - multi-mapped fuzzy anchor forms preserve all expected concepts
   - existing cross-script love and Hebrew `El` boundary behavior still works
13. Run targeted tests while working, then `go test ./...`.
14. Commit the completed quality-gate work locally with a clear message. Do not push.
15. Report final `git status --short --branch`, validation result, tests run, files changed, and before/after metrics for training precision/recall.

Acceptance Criteria:

- Active embedded YAML passes the project validator.
- The documented `socrates knowledge validate` command works.
- Glyph pattern matching preserves pattern kind boundaries.
- Channel diversity only rewards active meaningful evidence.
- Training evaluation cannot pass cleanly while producing excessive false activations.
- Multi-concept forms preserve complete fuzzy evidence paths.
- Existing harmonic core, training evaluation, and cross-script convergence still pass.
- Training does not mutate active knowledge.
- No black-box truth model is introduced.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

After this quality gate is clean, return to configurable ranking weights and a small held-out evaluation split. Only after that should Socrates add counsellor/transmutation fields.

## Nice To Have Later

- CLI ergonomics:
  - Let root text default to decipher, e.g. `socrates what is the purpose of my life?`.
  - Join multi-argument root input into one passage so quotes are optional for ordinary questions.
  - Keep `socrates decipher <text>` as the explicit form.
  - Make root help show subcommand flags clearly, including debug/training flags.
  - Consider hiding or eventually removing the typo alias `descifer` after compatibility is no longer useful.
  - Keep `socrates knowledge validate` working exactly as documented.
