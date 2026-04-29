# TODO: Socrates Active Work

This is the only active instruction file for Pi. Keep it short enough to use, but make the current task complete enough that Pi can work without babysitting. Completed work belongs in git history and `docs/TODO_ARCHIVE.md`.

## Authority

- `TODO.md` is the active task source.
- Ignore unlinked `.md` files if they claim to be active.
- The Architect assigns one task at a time.
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
- Activation-Energy Discipline Gate.
- Phase 5 first activation graph slice.
- Phase 5/6 graph hardening and first passage-field API.

Current git status reported after Pi's last task:

```text
## main...origin/main [ahead 10]
```

Architect review of Phase 5/6 slice:

- Good: graph propagation now uses path-local visited tracking, derived nodes can be depth `1`, edge identity is deduped, glyph repetition ignores non-letters, and `PassageField` exists.
- Gap: `PassageField` is not part of `Reading`, so the main engine/render path does not expose graph-backed passage fields.
- Gap: token provenance can be indirect or noisy because `AnalyzePassage` calls `engine.Analyze(token)` and then reuses convergence sources instead of preserving the original passage token as the primary source.
- Gap: `PassageField.RelationPaths` is defined but not populated.
- Gap: some passage tests are permissive (`t.Log`, `>=`, or "may or may not") and would not catch weak passage-field behavior.
- Gap: repeated `PropagateActivation()` on the same graph accumulates strength; Pi documented this as single-pass, but tests/API should guard or make the single-pass contract explicit.

## Current Next Task

Task:
Make passage fields production-facing, evidence-strict, and renderable without bloating default output.

Context:
The first passage-field API exists, but it is not yet an engine result. This task should finish Phase 6 enough that `Engine.Analyze` returns passage fields in `Reading`, debug output can inspect them, default output summarizes only top fields, and tests prove fields depend on actual passage tokens and graph relation paths.

Likely files:

- `internal/decipher/types.go`
- `internal/decipher/engine.go`
- `internal/decipher/render.go`
- `internal/decipher/passage_field.go`
- `internal/decipher/passage_field_test.go`
- `internal/decipher/activation_graph.go`
- `internal/decipher/activation_graph_test.go`
- `internal/decipher/engine_test.go`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Add passage fields to structured `Reading`, likely `PassageFields PassageFields`, without removing existing compatibility fields.
4. Route `Engine.Analyze` so phrase inputs populate `Reading.PassageFields` through the explicit activation graph.
5. Preserve real passage-token provenance. The original token from the passage must appear as the source for direct token activations.
6. Populate `PassageField.RelationPaths` from graph-backed relation paths relevant to that field.
7. Make `PassageField` carry enough evidence detail to defend each top field. Add an evidence-path field if needed rather than relying only on token strings.
8. Make repeated `PropagateActivation()` behavior explicit: either make propagation idempotent by resetting to base strengths before each run, or enforce a single-pass API with a test that prevents accidental repeated use.
9. Strengthen permissive tests. Replace `t.Log` and weak `>=` assertions with assertions that would fail if token provenance, field weakening, or relation paths broke.
10. Add a regression proving removing a key token from a passage strictly weakens the relevant field, not just "does not increase".
11. Add a regression proving repeated occurrences of the same token do not inflate a field beyond the meaningful evidence policy.
12. Add a regression proving relation paths appear on fields when connected concepts co-activate.
13. Default render should show only top passage fields and short evidence paths. It must remain concise.
14. Debug render should show full passage-field details, token sources, relation paths, and evidence paths.
15. Keep existing graph, bounded-work, render-mode, key-token, and architecture tests passing.
16. Do not implement harmonic/audio rendering or concept-to-frequency mappings.
17. Run targeted tests while working, then `go test ./...`.
18. Commit the completed passage-field integration locally with a clear message. Do not push.
19. Report final `git status --short --branch`, tests run, files changed, and limitations left intentionally out of scope.

Pi work requirement:
Spend at least 60 focused minutes. If the implementation passes quickly, use the remaining time to remove weak tests and add stricter passage-field regressions.

Acceptance Criteria:

- `Reading` exposes passage fields as structured data.
- `Engine.Analyze` populates passage fields for multi-token input.
- Passage fields preserve original token sources.
- Passage fields include relation paths or evidence paths sufficient to explain top fields.
- Removing a key passage token strictly weakens the related field.
- Duplicate/repeated token evidence follows a tested, non-inflating policy.
- Default output remains concise and evidence-first.
- Debug output exposes full passage-field details.
- Repeated graph propagation behavior is idempotent or explicitly guarded by tests.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

If this passes, the next task should be Phase 7: evidence-first CLI polish and examples using graph-backed passage fields.
