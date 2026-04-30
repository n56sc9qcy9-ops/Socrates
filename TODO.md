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

Current git status before this handoff:

```text
## main...origin/main [ahead 20]
```

Architect correction:

- Socrates can now evaluate curated examples separately from active knowledge.
- The next step is not counsellor prose yet.
- The next step is ranking-weight hardening: make evidence-path scores configurable, measurable, and auditable.
- Weight changes must be judged by training evaluation, not by subjective output preference.

## Current Next Task

Task:
Add configurable ranking weights and harden training evaluation.

Context:
Socrates now has a committed training/evaluation layer. It can load curated examples, validate them against active knowledge, run `Engine.Analyze`, and report concept/field precision and recall.

The next step is to make the ranking model explicit. Current scoring weights still live inside engine logic. Move the weight assumptions into a configurable structure and make evaluation show whether a weight change improves or degrades the curated examples.

Reference:

- `TRAINING_MODEL.md`

Likely files:

- `cmd/socrates/main.go`
- `internal/decipher/types.go`
- `internal/decipher/scoring.go`
- `internal/decipher/*`
- `internal/training/*`
- `training/examples.yaml`
- new file if useful: `internal/decipher/ranking_weights.go`
- new file if useful: `internal/training/report.go`
- new file if useful: `training/weights.yaml`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Review the committed training implementation and record any correctness issues before changing behavior.
4. Add a `RankingWeights` or equivalent structure for evidence-path scoring assumptions.
5. Include weights for at least:
   - exact form match
   - fuzzy form match
   - phonetic/script/glyph evidence if currently scored
   - confidence labels
   - graph relation support
   - propagation depth penalty
   - passage co-activation
   - harmonic profile support
   - harmonic ratio/archetype compatibility
   - duplicate/noise penalty
   - dissonance penalty
6. Provide default weights that preserve current behavior as closely as practical.
7. Allow evaluation to run with default weights.
8. Add optional `--weights <path>` support if the implementation shape is clear.
9. Evaluation output must show which weight set was used.
10. Add a detailed/debug report that shows feature contributions for expected and false activations.
11. Keep training examples separate from active knowledge.
12. Do not let training mutate active knowledge.
13. Do not add AI or embeddings in this task.
14. Do not add automatic optimization yet unless it is isolated behind a separate explicit command and deterministic.
15. Fix any discovered training CLI mismatch or evaluator correctness issue while staying scoped.
16. Add tests proving:
   - default weights preserve baseline evaluation behavior
   - invalid weight files fail validation if `--weights` is added
   - evaluation reports the active weight set
   - expected activations expose feature/evidence contributions in debug/report mode
   - training still does not mutate active knowledge
17. Run targeted tests while working, then `go test ./...`.
18. Commit the completed ranking-weight work locally with a clear message. Do not push.
19. Report final `git status --short --branch`, tests run, files changed, and an example weighted evaluation output.

Acceptance Criteria:

- Ranking weights are explicit and inspectable.
- Default weights preserve current behavior as closely as practical.
- Training evaluation reports which weights were used.
- Debug/report output shows evidence or feature contributions for activations.
- Training does not mutate active knowledge.
- No black-box truth model is introduced.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

After ranking weights are explicit and measurable, add a small held-out evaluation split. Only after that should Socrates add counsellor/transmutation fields.
