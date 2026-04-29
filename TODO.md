# TODO: Socrates Active Work

This is the only active instruction and checklist file for the Architect and Coder. Completed phases live in `docs/TODO_ARCHIVE.md`.

## Instruction Authority

The Architect must read this file first and give the Coder one narrow task from the current next task only.

- `TODO.md` is the only active authority.
- A separate task file is active only when this file explicitly links to it.
- Unlinked `.md` files are background, archive, or historical context.
- If another `.md` file claims to be active but is not linked here, ignore that claim and follow `TODO.md`.
- Do not mark checklist items complete unless the listed acceptance criteria are actually met.
- Do not skip to later phases while the current next task has unmet acceptance criteria.
- The Architect does not run tests.
- The Coder/Pi runs the baseline tests before changing code, then runs tests again after changes.

Markers:

```text
[ ] not started
[/] in progress
[x] complete
[!] blocked / needs review
```

## Current Next Task

Task:
Split default output from debug output while preserving evidence-first structured data.

Context:
Phase 3A cleanup, deduplication, and bounded-work hardening are complete. Pi reported `go test ./...` passing after the bounded candidate/fuzzy work on 2026-04-29, with 53 decipher tests, 9 knowledge tests, and 5 resonance tests. The next activation-energy risk is output discipline: default output should show concise evidence paths and conclusions, while full generated candidates, fuzzy matches, and internal bounded-work details should move behind explicit debug output.

Files:
- `internal/decipher/types.go`
- `internal/decipher/engine.go`
- `internal/decipher/render.go`
- `cmd/socrates/main.go`
- `internal/decipher/engine_test.go`

Instructions:
1. Start with `git status --short --branch` and record whether the worktree is dirty, ahead, or behind. Do not pull/rebase over dirty local work.
2. Run `go test ./...` before changing code and record the exact baseline result in the final report.
3. Inspect current CLI and render output before editing. Identify exactly where candidates, fuzzy matches, generated forms, discarded counts, channels, and evidence paths are printed.
4. Define an explicit output mode shape. Prefer a small typed option such as `RenderModeDefault` / `RenderModeDebug` or a `RenderOptions` struct over ad hoc booleans.
5. Keep default output concise: input, top convergence/evidence paths, score summary, warnings, and bounded-work discard counts are acceptable.
6. Move full candidate lists, full fuzzy match lists, verbose generated forms, and internal match details to debug output only.
7. Preserve structured `Reading` fields. This task changes rendering/CLI output discipline, not the engine's ability to return data.
8. Add or update CLI flags for debug output only if the CLI currently has a natural place for it. Do not create a broad CLI redesign.
9. Add tests proving default rendering does not dump candidate or fuzzy-match internals.
10. Add tests proving debug rendering can still show generated forms, candidates, fuzzy matches, and discarded counts.
11. Add tests proving every default conclusion shown has an evidence path or channel signal reference, not only a bare concept name.
12. Keep existing public behavior stable where practical. Do not add word-specific branches or hardcoded semantic exceptions to make output tests pass.
13. Run targeted render/CLI tests after implementation, then run `go test ./...` as the final verification.
14. After final tests pass, update only the checklist items whose acceptance criteria are actually proven.
15. Update git at the end: commit the completed output-discipline changes with a clear message, or state exactly why committing was not possible. Do not push. Leave `git status --short --branch` in the final report.

Pi work requirement:
Spend at least 10 focused minutes on this task before reporting back. If the first code change passes quickly, use the remaining time for review, edge-case tests, and checking that default output is concise without hiding debug data.

Constraints:
- No hardcoded semantic word lists in production Go.
- No direct behavior for specific words like `skal`, `energy`, `truth`, or `light`.
- No concept-to-frequency, pitch, music, harmonic, or audio mappings.
- Do not implement harmonic/audio rendering.
- Do not add production Go branches for specific words.
- Knowledge and confidence assumptions must stay data-driven.
- Do not change candidate-generation or fuzzy-matching semantics unless required by a rendering test failure.
- Do not remove the existing candidate-generation methods.
- Do not make fuzzy matching depend on global mutable state.

Tests:
- Default render output does not include full candidate dumps.
- Default render output does not include full fuzzy-match dumps.
- Debug render output includes generated forms, candidates, fuzzy matches, and discarded counts.
- Default output still shows evidence paths or channel signal references for displayed conclusions.
- Existing bounded-work, deduplication, and fuzzy-match tests still pass.

Acceptance Criteria:
- Default CLI/render output is concise and does not dump candidate or match internals.
- Debug CLI/render output can still show generated forms, candidates, fuzzy matches, and discarded counts.
- Conclusions visible in default output have evidence paths or channel signal references.
- The relevant Activation-Energy Discipline Gate checklist items below are updated only if acceptance is met.
- `go test ./...` passes after changes.

Do Not:
- Do not implement harmonic/audio rendering.
- Do not add new semantic alias maps in Go.
- Do not mark the full Activation-Energy Discipline Gate complete.
- Do not leave the repository in an unreported dirty/ahead/behind state.

## North Star

Build a real language-resonance engine.

The base model is coherent activation energy first.

The code may generate forms, compute similarity, load knowledge data, activate concepts, propagate through a graph, score convergence, and render evidence. Linguistic, spiritual, symbolic, harmonic, musical, and cross-language knowledge must live in data, not Go logic.

Resonance is not evidence volume. Resonance is coherent activation: independent, confidence-weighted evidence channels reinforcing related concepts through visible paths.

## Non-Negotiable Rules

- [x] No hardcoded semantic word lists in production Go. Cross-language maps and semantic alias normalization removed.
- [x] No direct `skal` behavior in production Go.
- [x] No modal/emptiness/contrast marker maps in production Go.
- [x] No silent fallback to legacy hardcoded knowledge.
- [!] Tests must guard against fake progress, shortcut semantics, duplicate score inflation, and verbosity regressions.
- [x] No hardcoded glyph/script symbolic meanings in channel code.
- [x] No global mutable knowledge dependency.
- [ ] No mixed-purpose files over 500 lines without explicit approval.
- [ ] Every output must show evidence paths, not only conclusions.
- [ ] Default output must be concise; full candidates and fuzzy matches must be debug output.
- [x] No score inflation from duplicated fragments, repeated weak matches, or debug noise.
- [!] No concept-to-frequency, concept-to-pitch, or harmonic mappings in Go. Legacy `internal/resonance` frequency code still exists; do not expand it.

## Current Priority

1. [x] Finish remaining Phase 3A cleanup blockers. **COMPLETED 2026-04-28**
2. Activation-Energy Discipline Gate: scoring, deduplication, and verbosity.
3. Phase 5: Activation Graph hardening.
4. Phase 6: Passage Field Analysis.
5. Phase 7: Evidence-First CLI Output modes.
6. Phase 8: Regression Tests.
7. Phase 9: Final Cleanup.
8. Future: Harmonic Rendering Layer.

Do not add exact word meanings to Go code. Add or change knowledge only in YAML.

Do not build audio yet. Harmonic rendering comes after the activation-energy core is bounded, deduplicated, data-driven, and evidence-first.

The Architect should assign only one narrow coder task at a time. The current task is Activation-Energy Discipline Gate, not audio.

## Phase 3A: Architecture Cleanup Gate

Stop feature work until this is complete. This phase is a no-behavior-change cleanup except where hardcoded semantic knowledge is moved from Go into YAML.

- [x] Split `internal/decipher/discovery.go` by responsibility:
  - candidate generation
  - similarity
  - activation
  - convergence
  - scoring
- [x] Split `internal/decipher/channels.go` by channel.
- [x] Move Latin glyph patterns from Go code into data.
- [x] Move Hebrew glyph associations from Go code into data.
- [x] Move Han character associations/readings from Go code into data.
- [x] Remove global mutable `KnowledgeBase`.
- [x] Pass knowledge through `Engine` or explicit analyzer structs.
- [!] Remove or rename legacy architecture types such as `Primitive`, `FragmentSeed`, and `FragmentLens` if they are no longer part of the real design. Deferred cleanup; not the current gate.
- [x] Remove root build/artifact files such as `socrates` and accidental empty files such as `skal`, after confirming they are not source files.
- [x] Add regression tests that fail if semantic mappings return to production channel code.
- [!] Add file-size/code-shape tests or checks for oversized mixed-purpose files. Deferred cleanup; not the current gate.
- [x] Move hardcoded cross-language echo/root maps from Go into data.
- [x] Move semantic alias normalization out of `engine.go` into data or remove it.

Acceptance:

- [x] No production Go file contains hardcoded symbolic mappings such as `Target: "breath"` inside glyph/script rule code.
- [x] `internal/decipher/discovery.go` is split into focused files and no longer acts as the main dumping ground.
- [x] `internal/decipher/channels.go` is split into focused channel files.
- [x] Engine analysis can run with an explicit knowledge fixture without touching global state.
- [x] No production Go file contains hardcoded cross-language or semantic alias maps.
- [x] Root artifact files are removed or explicitly justified.
- [x] `go test ./...` passes.
- [x] Phase completion is recorded only after the above checks pass.

Phase 3A completed on 2026-04-28:

- Removed `internal/decipher/channel_cross_language.go` (cross-language maps moved to YAML forms)
- Simplified `normalizeTarget()` in `engine.go` (semantic aliasing handled by knowledge layer)
- Removed root artifact files `skal` and `socrates`
- Updated regression tests to reflect file removal
- Added `TestCrossLanguageFormsViaYAML` and `TestNoCrossLanguageChannel` tests
- All tests pass (`go test ./...` passes)
- Cross-language forms (chi, ruach, prana, om, logos, dao) are now data-driven via forms.yaml

## Activation-Energy Discipline Gate

Stop new feature work until scoring, deduplication, and output verbosity reflect coherent activation energy.

- [x] Define score semantics for coherent activation energy. Deduplication and bounded work are implemented; output separation implemented 2026-04-29.
- [x] Deduplicate evidence by meaningful path, not only by rendered text.
- [x] Prevent repeated fragments from increasing strength without independent support.
- [x] Prevent weak fuzzy matches from dominating scores.
- [x] Make exact/verified multi-channel evidence score above speculative single-channel evidence.
- [x] Make candidate generation bounded for long input.
- [x] Cap fuzzy matching work and expose discarded counts.
- [x] Split default output from debug output. Completed 2026-04-29.
- [x] Keep default CLI output concise. Default is ~50 lines, debug is ~170 lines.
- [/] Add tests for duplicate evidence, weak fuzzy noise, and key-token removal. Duplicate/weak tests are present; key-token removal remains.

Acceptance:

- [x] Duplicated evidence does not increase resonance after the first meaningful occurrence.
- [x] Weak fuzzy matches do not dominate the score.
- [x] Exact/verified multi-channel evidence scores higher than speculative single-channel evidence.
- [ ] Removing a key evidence token weakens the relevant activation field.
- [x] Default output does not dump candidate or match internals.
- [x] Debug output can still show generated forms and match details.
- [x] `go test ./...` passes for the completed Phase 3A and deduplication pass.
- [x] `go test ./...` passes after the strict bounded-work verification pass. Pi reported passing tests on 2026-04-29: 53 decipher, 9 knowledge, 5 resonance.

## Phase 3: Generic Form Generation

- [!] Implemented prematurely before Phase 3A passed. Re-verify after architecture cleanup.
- [ ] Generate n-grams.
- [ ] Generate prefix/suffix fragments.
- [ ] Deduplicate generated forms.
- [ ] Attach method and distance to every candidate.
- [ ] Make candidate output visible in debug CLI output.

Acceptance:

- [ ] Form generation works without knowing any word meaning.
- [ ] Candidate output explains how each form was produced.

## Phase 4: Similarity Engine

- [ ] Levenshtein distance.
- [ ] Normalized edit distance.
- [ ] N-gram similarity.
- [ ] Prefix/suffix overlap.
- [ ] Consonant skeleton similarity.
- [ ] Vowel skeleton similarity.
- [ ] Phonetic similarity.
- [ ] Weighted combined similarity.
- [ ] Configurable thresholds.
- [ ] Deduplicate matches by target form and best evidence.

Acceptance:

- [ ] `skal` can find nearby forms without a direct code path.
- [ ] Noisy input like `zzskalx` scores weaker than `skal`.
- [ ] Exact matches score higher than fuzzy matches.

## Phase 5: Activation Graph

Data flow:

```text
input candidates
  -> matched knowledge forms
  -> activated concepts
  -> related concepts
  -> convergence fields
```

- [ ] Define `Activation`.
- [ ] Define `ActivationEdge`.
- [ ] Define `ActivationGraph`.
- [ ] Activate concepts from form matches.
- [ ] Propagate activation through relations.
- [ ] Decay activation by relation weight.
- [ ] Track evidence paths.
- [ ] Prevent infinite loops.
- [ ] Rank activated concepts.

Acceptance:

- [ ] No special-case passage pattern logic is needed.
- [ ] Convergence emerges from graph activation strength.

## Phase 6: Passage Field Analysis

- [ ] Tokenize passage.
- [ ] Generate candidates for each token.
- [ ] Match each token against known forms.
- [ ] Merge token activations into one graph.
- [ ] Track which token produced each activation.
- [ ] Detect repeated activation fields.
- [ ] Detect concept co-activation through graph structure.
- [ ] Score passage-level convergence.

Acceptance:

- [ ] A passage can activate multiple fields at once.
- [ ] Removing one important token weakens the related field.
- [ ] Output shows evidence paths for the field.

## Phase 7: Evidence-First CLI Output

- [ ] Show concise resonance summary by default.
- [ ] Show activated concepts.
- [ ] Show graph propagation paths.
- [ ] Show convergence fields.
- [ ] Show weak signals.
- [ ] Show score components.
- [ ] Show warnings for speculative paths.
- [ ] Keep final prose short.
- [ ] Move generated forms and full similarity matches behind debug output.

Acceptance:

- [ ] User can see why a resonance reading appeared.
- [ ] Output does not hide behind poetic language.
- [ ] Default output is not dominated by candidate or fuzzy-match dumps.

## Phase 8: Regression Tests

- [ ] Test no semantic marker maps exist in production code.
- [ ] Test no hardcoded knowledge registries exist in production code.
- [ ] Test `skal` has no direct Go semantic branch.
- [ ] Test removing a form from YAML removes that activation.
- [ ] Test removing a relation from YAML removes that graph path.
- [ ] Test exact match beats fuzzy match.
- [ ] Test noisy input scores lower.
- [ ] Test passage convergence changes when a key token is removed.
- [ ] Test evidence paths are present for top concepts.
- [ ] Test all loaded knowledge has confidence labels.

Acceptance:

- [ ] Tests protect architecture, not example outputs.

## Phase 9: Final Cleanup

- [ ] Remove obsolete direct channel behavior that duplicates activation graph behavior.
- [ ] Keep only generic channels and generic graph logic.
- [ ] Update docs to match actual implementation.
- [ ] Remove stale examples that imply canned interpretations.

## Final Report Format

When stopping work, report:

- checklist items completed
- files changed
- tests run
- remaining shortcut semantics found, if any
- example CLI outputs
- known limitations
