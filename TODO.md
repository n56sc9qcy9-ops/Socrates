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

## Corrected Core Principle

Socrates is a harmonic resonance engine.

Syllables, words, meanings, thoughts, and ideas are treated as cascading frequency fields. The activation graph is not the final definition of resonance; it is the evidence layer that decides which frequency fields are active and how strongly they are supported.

The project must not lose this principle again.

## Standing Rules

- No hardcoded semantic word lists in production Go.
- No direct behavior for specific words like `skal`, `energy`, `truth`, or `light`.
- No hardcoded concept-to-frequency, concept-to-color, concept-to-chakra, concept-to-note, pitch, music, harmonic, or audio mappings in Go.
- Harmonic/frequency mappings must live in curated data with source/lens/confidence.
- Knowledge and confidence assumptions must stay data-driven.
- Default output must stay concise; full candidates, fuzzy matches, and harmonic internals belong in debug output.

## Recent Context

Completed and committed locally:

- Phase 3A cleanup.
- Activation-Energy Discipline Gate.
- Phase 5/6 graph and passage-field work.
- First knowledge validation/suggestion pipeline.

Current git status before this handoff:

```text
## main...origin/main [ahead 16]
```

Architect correction:

- Existing docs had drifted: they described harmonics as a later rendering layer rather than the core principle.
- The graph work should remain, but only as the evidence scaffold for frequency-field activation.
- Legacy `internal/resonance` hardcodes sample frequencies in Go. That is the wrong final shape and should be replaced or bypassed by data-backed harmonic profiles.

## Current Next Task

Task:
Correct the architecture by adding data-backed frequency profiles and harmonic field scoring.

Context:
Before adding more ordinary concepts, Socrates needs the missing harmonic core. The engine should be able to say: these concepts are active, these active concepts have curated frequency profiles, these profiles form a coherent or dissonant field, and this is the evidence path for each tone in the field.

Likely files:

- `internal/knowledge/loader.go`
- `internal/knowledge/knowledge.go`
- `internal/knowledge/validate.go`
- `internal/knowledge/validate_test.go`
- `internal/knowledge/*.yaml`
- `internal/decipher/types.go`
- `internal/decipher/engine.go`
- `internal/decipher/render.go`
- new file if useful: `internal/knowledge/frequencies.yaml`
- new file if useful: `internal/decipher/harmonic_field.go`
- new test file if useful: `internal/decipher/harmonic_field_test.go`
- legacy review: `internal/resonance/*`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Preserve the knowledge-pipeline hardening task if it is already in progress, but prioritize the harmonic-core correction.
4. Add a YAML-backed frequency/harmonic profile model. It should support at least:
   - concept ID
   - base frequency or ratio representation
   - optional note/color/chakra/field labels
   - confidence
   - lens/system/source
   - weight
5. Load and validate frequency profiles through the knowledge layer.
6. Validation must reject profiles for unknown concepts, invalid confidence, invalid weight, and malformed frequency/ratio values.
7. Do not hardcode concept-to-frequency mappings in Go. Tests may use fixtures.
8. Build a harmonic field from active `Reading.PassageFields` / activation graph concepts.
9. Add harmonic scoring:
   - consonance/coherence for compatible ratios or shared systems
   - dissonance/tension for incompatible or weakly supported profiles
   - confidence-weighted strength
10. Preserve evidence paths from activated concepts into harmonic field tones.
11. Add structured result data to `Reading`, e.g. `HarmonicField`, without breaking existing render behavior.
12. Default render should summarize harmonic field only when profiles exist. Debug render should show profile details, ratios/frequencies, confidence, source/lens, and evidence paths.
13. Add tests proving `love + heart + truth` or a similarly data-backed cluster creates a harmonic field from data, not Go constants.
14. Add tests proving a concept without a frequency profile does not invent one.
15. Add tests proving malformed frequency profile YAML fails validation.
16. Decide what to do with legacy `internal/resonance`: document as legacy, bypass it, or refactor it to consume data-backed profiles. Do not expand its hardcoded constants.
17. Keep all existing decipher/render/passage/graph/knowledge tests passing.
18. Do not implement audio synthesis yet. The target is harmonic data and scoring, not sound output.
19. Run targeted tests while working, then `go test ./...`.
20. Commit the completed harmonic-core work locally with a clear message. Do not push.
21. Report final `git status --short --branch`, tests run, files changed, and an example harmonic field output.

Pi work requirement:
Spend at least 90 focused minutes. This is a core correction, not a small feature.

Acceptance Criteria:

- Frequency/harmonic profiles are represented in YAML data.
- Profiles load through the knowledge layer.
- Profiles validate strictly.
- No production Go file contains concept-to-frequency/color/chakra/note mappings.
- `Reading` can expose a harmonic field derived from active concepts.
- Harmonic field scoring is data-backed and confidence-weighted.
- Default/debug rendering handles harmonic fields appropriately.
- Legacy hardcoded `internal/resonance` is not expanded.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

If this passes, resume knowledge-pipeline hardening and then add curated frequency profiles/concept packs through the validated workflow.
