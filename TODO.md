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
- Phase 5 first activation graph slice.

Current git status reported after Pi's last task:

```text
## main...origin/main [ahead 8]
```

Architect review of Phase 5 first slice:

- Good: first-class graph types exist, convergence is routed through graph, scores are bounded, and `"God is Love"` has regression coverage.
- Must fix next: whitespace can still create a repetition signal because glyph frequency counts spaces before emitting generic `repeated letters detected`.
- Must fix next: graph propagation uses global mutable `Visited` state and appears order-dependent; use path-local visited tracking instead.
- Must fix next: graph-expanded nodes are created with depth `0`, so derived nodes can look direct and inflate direct counts.
- Should fix next: duplicate edge creation/evidence is not strongly deduplicated.

## Current Next Task

Task:
Finish activation graph correctness hardening and implement the first Phase 6 passage-field analysis slice.

Context:
The Activation-Energy Discipline Gate is complete and Phase 5 has a first graph implementation. The next task is intentionally larger: harden the graph invariants that still look weak, then build the first Phase 6 passage-field layer on top of the explicit graph. The outcome should be a multi-token passage analysis where each token can contribute evidence, repeated activation fields are merged, and removing a key token weakens the relevant field.

Likely files:

- `internal/decipher/types.go`
- `internal/decipher/activation.go`
- `internal/decipher/activation_graph.go`
- `internal/decipher/convergence.go`
- `internal/decipher/engine.go`
- `internal/decipher/scoring.go`
- `internal/decipher/engine_test.go`
- `internal/decipher/channel_glyph.go`
- new file if useful: `internal/decipher/passage_field.go`
- new test file if useful: `internal/decipher/passage_field_test.go`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Re-run or inspect the `"God is Love"` regression expectations. Strengthen it so whitespace cannot create a repetition signal at all.
4. Fix Latin glyph repetition so only letters are considered for repeated-letter evidence. Spaces and punctuation must not create symbolic glyph evidence.
5. Replace global mutable `Visited` propagation with path-local cycle tracking. Propagation should be deterministic and independent of map iteration order.
6. Distinguish direct nodes from graph-expanded nodes. Derived nodes should have depth greater than 0 and should not count as direct evidence.
7. Deduplicate graph edges by meaningful identity: from concept, to concept, relation type, and evidence path.
8. Add tests proving propagation is deterministic across repeated runs and safe through cycles.
9. Add tests proving derived nodes are not counted as direct nodes.
10. Add tests proving duplicate edges/evidence do not inflate activation.
11. Implement a first `PassageField` or equivalent structure that groups passage-level activation by field/concept, token sources, evidence paths, strength, confidence, and relation paths.
12. Passage-field analysis must tokenize a passage, analyze each token through existing candidate/match/channel logic where practical, merge token activations into the graph, and preserve which token produced each activation.
13. Detect repeated activation fields across tokens without duplicate inflation.
14. Detect concept co-activation through graph structure using the explicit activation graph.
15. Add a passage-level regression where removing one important token weakens the relevant field. Use data-backed concepts and no production word-specific branches.
16. Keep default output concise. Add passage field detail to debug output if useful; default should show only top fields and evidence paths.
17. Keep all existing key-token, bounded-work, render-mode, graph, and architecture regression tests passing.
18. Do not implement harmonic/audio rendering or concept-to-frequency mappings.
19. Run targeted tests while working, then `go test ./...`.
20. Commit the completed graph/passage work locally with a clear message. Do not push.
21. Report final `git status --short --branch`, tests run, files changed, and any limitations left intentionally out of scope.

Pi work requirement:
Spend at least 60 focused minutes. This is deliberately about twice the previous task size. If the first implementation passes quickly, use the remaining time to harden invariants, add one more meaningful passage regression, and inspect default/debug output.

Acceptance Criteria:

- Whitespace and punctuation do not create repeated-letter glyph evidence.
- Graph propagation is bounded, cycle-safe, and deterministic.
- Direct nodes and graph-derived nodes are distinguishable by depth/source.
- Duplicate graph edges and evidence paths do not inflate activation.
- Passage-field structure exists and records field concept, strength, confidence, token sources, evidence paths, and relation paths.
- Passage analysis merges multi-token activations into the graph.
- Repeated activation fields across tokens are detected without duplicate inflation.
- Removing an important passage token weakens the related field.
- Default output remains concise; debug output can expose passage-field detail.
- Duplicate evidence paths do not inflate activation.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

If this passes, the next task should be Phase 7: evidence-first CLI polish on top of graph-backed passage fields.
