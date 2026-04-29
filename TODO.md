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
- Deduplication and bounded candidate/fuzzy work.
- Default/debug render split with `--debug`.
- Key-token removal regression tests.

Current git status reported after Pi's last task:

```text
## main...origin/main [ahead 5]
```

Manual output check to protect:

```text
go run ./cmd/socrates decipher "God is Love"
```

Observed issues that this graph task must address or explicitly justify:

- `Score.Components.GraphExpansionScore` rendered as `2.45`; score components should have clear bounded semantics and should not exceed their documented range.
- Default output repeats identical evidence such as `primitive 'love' matches love -> love`; graph evidence paths should deduplicate meaningful paths before rendering/scoring.
- Glyph channel treats repeated spaces as `repetition`; whitespace should not create symbolic glyph evidence unless deliberately justified.
- Exact/fuzzy scores show `0.00` even though `love` is visibly matched through fragments; this may be valid by current channel separation, but Pi should verify and document the boundary.
- Bounded work discarded counts are very high for a short phrase; Pi should verify phrase candidate generation is bounded per token or explain why whole-phrase candidate generation is expected here.

## Current Next Task

Task:
Implement the first first-class Activation Graph layer and route convergence through it.

Context:
The Activation-Energy Discipline Gate is complete. The next larger task is Phase 5 hardening. The engine already has `ActivatedConcept`, `PassageSignal`, `ConceptExpansions`, and string relation paths, but there is no explicit graph object that owns nodes, edges, propagation, decay, evidence paths, and loop prevention. Build that graph layer without changing the project into a different product.

Likely files:

- `internal/decipher/types.go`
- `internal/decipher/activation.go`
- `internal/decipher/convergence.go`
- `internal/decipher/engine.go`
- `internal/decipher/scoring.go`
- `internal/decipher/engine_test.go`
- new file if useful: `internal/decipher/activation_graph.go`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Inspect `activation.go`, `convergence.go`, `scoring.go`, and `engine.go` before editing.
4. Define explicit graph types, using names close to `ActivationGraph`, `ActivationNode`, `ActivationEdge`, and `EvidencePath` unless the existing code strongly suggests better names.
5. Build graph nodes from direct channel concepts, passage signals, fuzzy matches, and graph-expanded concepts where those are already available.
6. Build graph edges from YAML-backed knowledge relations only. Do not add concept relation maps in Go.
7. Track evidence paths on nodes/edges: source token/form/channel, match form when available, relation type, confidence, and weight/strength.
8. Propagate activation through relations with explicit decay based on relation weight and bounded depth. Default depth should be conservative, likely 1 or 2.
9. Prevent infinite loops and duplicate inflation. Repeated evidence for the same meaningful path must not keep increasing strength.
10. Keep `Reading.Convergence` and existing render output working. Add graph data to structured reading only if needed and keep default output concise.
11. Route `DetectConvergence` through the new graph internally, even if compatibility fields such as `ActivatedConcepts`, `TopConcepts`, and `RelationPaths` remain.
12. Add tests for graph node creation from evidence, relation propagation with decay, evidence path preservation, cycle/loop prevention, and deduplication.
13. Add a regression test or golden-ish assertion around `God is Love` or an equivalent short phrase proving graph scores are bounded, repeated evidence paths are not duplicated, and whitespace does not create symbolic repetition evidence.
14. Keep the existing key-token, bounded-work, render-mode, and architecture regression tests passing.
15. Do not implement Phase 6 passage field analysis beyond what is needed to build the graph from existing passage signals.
16. Do not implement harmonic/audio rendering or concept-to-frequency mappings.
17. Run targeted tests while working, then `go test ./...`.
18. Commit the completed graph work locally with a clear message. Do not push.
19. Report final `git status --short --branch`, tests run, files changed, and any graph limitations left intentionally out of scope.

Pi work requirement:
Spend at least 30 focused minutes. This is larger than the previous regression task. If the first implementation passes quickly, use the remaining time to review graph invariants and add one more meaningful regression test rather than stopping early.

Acceptance Criteria:

- `ActivationGraph` or equivalent first-class graph structure exists.
- Graph nodes represent activated concepts with strength, confidence, and evidence paths.
- Graph edges represent YAML-backed concept relations with relation type, weight, and propagated strength.
- Propagation is bounded and cycle-safe.
- Duplicate evidence paths do not inflate activation.
- `DetectConvergence` uses the graph internally.
- Score components exposed in default output have bounded, defensible ranges.
- Default evidence output deduplicates repeated meaningful paths.
- Whitespace does not create symbolic repetition evidence.
- Existing public render behavior remains stable.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

If this graph slice passes, the next task should move to Phase 6: passage field analysis on top of the explicit activation graph.
