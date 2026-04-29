# TODO: Socrates Active Work

This is the only active instruction file for Pi. Keep it short. Completed work belongs in git history and `docs/TODO_ARCHIVE.md`, not here.

## Authority

- `TODO.md` is the active task source.
- Ignore unlinked `.md` files if they claim to be active.
- The Architect assigns one narrow task at a time.
- The Architect does not run tests.
- Pi runs `go test ./...` before and after implementation.
- Pi commits completed work locally with a clear message and does not push.
- Do not mark work complete unless tests and acceptance criteria prove it.

## Standing Rules

- No hardcoded semantic word lists in production Go.
- No direct behavior for specific words like `skal`, `energy`, `truth`, or `light`.
- No concept-to-frequency, pitch, music, harmonic, or audio mappings in Go.
- No harmonic/audio rendering yet.
- Knowledge and confidence assumptions must stay data-driven.
- Default output must stay concise; full candidates and fuzzy matches belong in debug output.

## Recent Context

Completed and committed locally:

- Phase 3A cleanup.
- Deduplication and bounded candidate/fuzzy work.
- Default/debug render split with `--debug`.

Current git status reported after Pi's last task:

```text
## main...origin/main [ahead 3]
```

## Current Next Task

Task:
Add the missing key-token removal regression for activation-energy discipline, and make the smallest fix needed if it fails.

Context:
The Activation-Energy Discipline Gate is mostly complete. The remaining acceptance gap is proving that removing an important evidence token weakens the relevant activation field. This protects against fake resonance where output looks evidence-based but scores/convergence do not actually depend on the evidence.

Likely files:

- `internal/decipher/engine_test.go`
- `internal/decipher/activation.go`
- `internal/decipher/convergence.go`
- `internal/decipher/scoring.go`
- `internal/decipher/types.go`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Inspect existing passage/convergence/activation tests before adding new behavior.
4. Add a regression test comparing a passage/input with a key evidence token against the same passage/input with that token removed.
5. Assert that the relevant activation/convergence/score weakens when the key token is removed.
6. Prefer a data-backed fixture or existing YAML-backed concepts over hardcoded production behavior.
7. If the test fails, make the smallest engine/scoring/convergence fix needed.
8. Do not broaden this into Phase 5 graph redesign.
9. Run targeted tests, then `go test ./...`.
10. Commit the completed work locally with a clear message. Do not push.
11. Report final `git status --short --branch`, tests run, and files changed.

Pi work requirement:
Spend at least 10 focused minutes. If the first test passes immediately, use the remaining time to check that the test is not tautological and would fail if the key-token dependency were broken.

Acceptance Criteria:

- A test proves removing a key evidence token weakens the relevant activation field.
- The test is data-driven and does not require production word-specific branches.
- Existing bounded-work and render-mode tests still pass.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

If the key-token regression passes, the Activation-Energy Discipline Gate can be treated as complete and the next task should move to Phase 5: Activation Graph hardening.
