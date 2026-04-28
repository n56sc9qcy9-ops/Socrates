# TODO: Socrates Active Work

This is the only active checklist. Completed phases live in `docs/TODO_ARCHIVE.md`.

Markers:

```text
[ ] not started
[/] in progress
[x] complete
[!] blocked / needs review
```

## North Star

Build a real language-resonance engine.

The base model is coherent activation energy first.

The code may generate forms, compute similarity, load knowledge data, activate concepts, propagate through a graph, score convergence, and render evidence. Linguistic, spiritual, symbolic, harmonic, musical, and cross-language knowledge must live in data, not Go logic.

Resonance is not evidence volume. Resonance is coherent activation: independent, confidence-weighted evidence channels reinforcing related concepts through visible paths.

## Non-Negotiable Rules

- [!] No hardcoded semantic word lists in production Go. Current violations include cross-language maps and semantic alias normalization in Go.
- [x] No direct `skal` behavior in production Go.
- [x] No modal/emptiness/contrast marker maps in production Go.
- [x] No silent fallback to legacy hardcoded knowledge.
- [!] Tests must guard against fake progress, shortcut semantics, duplicate score inflation, and verbosity regressions.
- [x] No hardcoded glyph/script symbolic meanings in channel code.
- [x] No global mutable knowledge dependency.
- [ ] No mixed-purpose files over 500 lines without explicit approval.
- [ ] Every output must show evidence paths, not only conclusions.
- [ ] Default output must be concise; full candidates and fuzzy matches must be debug output.
- [ ] No score inflation from duplicated fragments, repeated weak matches, or debug noise.
- [ ] No concept-to-frequency, concept-to-pitch, or harmonic mappings in Go.

## Current Priority

1. Finish remaining Phase 3A cleanup blockers.
2. Activation-Energy Discipline Gate: scoring, deduplication, and verbosity.
3. Phase 5: Activation Graph hardening.
4. Phase 6: Passage Field Analysis.
5. Phase 7: Evidence-First CLI Output modes.
6. Phase 8: Regression Tests.
7. Phase 9: Final Cleanup.
8. Future: Harmonic Rendering Layer.

Do not add exact word meanings to Go code. Add or change knowledge only in YAML.

Do not build audio yet. Harmonic rendering comes after the activation-energy core is bounded, deduplicated, data-driven, and evidence-first.

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
- [ ] Remove or rename legacy architecture types such as `Primitive`, `FragmentSeed`, and `FragmentLens` if they are no longer part of the real design.
- [ ] Remove root build/artifact files such as `socrates` and accidental empty files such as `skal`, after confirming they are not source files.
- [ ] Add regression tests that fail if semantic mappings return to production channel code.
- [ ] Add file-size/code-shape tests or checks for oversized mixed-purpose files.
- [ ] Move hardcoded cross-language echo/root maps from Go into data.
- [ ] Move semantic alias normalization out of `engine.go` into data or remove it.

Acceptance:

- [x] No production Go file contains hardcoded symbolic mappings such as `Target: "breath"` inside glyph/script rule code.
- [x] `internal/decipher/discovery.go` is split into focused files and no longer acts as the main dumping ground.
- [x] `internal/decipher/channels.go` is split into focused channel files.
- [x] Engine analysis can run with an explicit knowledge fixture without touching global state.
- [ ] No production Go file contains hardcoded cross-language or semantic alias maps.
- [ ] Root artifact files are removed or explicitly justified.
- [ ] `go test ./...` passes.
- [ ] Phase completion is recorded only after the above checks pass.

Known current blockers checked on 2026-04-28:

- `internal/decipher/channel_cross_language.go` still contains hardcoded cross-language semantic maps.
- `internal/decipher/engine.go` still contains semantic alias normalization.
- Root artifact files `skal` and `socrates` still exist in the working tree.
- Do not report Phase 3A complete until those are resolved and `go test ./...` passes.

## Activation-Energy Discipline Gate

Stop new feature work until scoring, deduplication, and output verbosity reflect coherent activation energy.

- [ ] Define score semantics for coherent activation energy.
- [ ] Deduplicate evidence by meaningful path, not only by rendered text.
- [ ] Prevent repeated fragments from increasing strength without independent support.
- [ ] Prevent weak fuzzy matches from dominating scores.
- [ ] Make exact/verified multi-channel evidence score above speculative single-channel evidence.
- [ ] Make candidate generation bounded for long input.
- [ ] Cap fuzzy matching work and expose discarded counts.
- [ ] Split default output from debug output.
- [ ] Keep default CLI output concise.
- [ ] Add tests for duplicate evidence, weak fuzzy noise, and key-token removal.

Acceptance:

- [ ] Duplicated evidence does not increase resonance after the first meaningful occurrence.
- [ ] Weak fuzzy matches do not dominate the score.
- [ ] Exact/verified multi-channel evidence scores higher than speculative single-channel evidence.
- [ ] Removing a key evidence token weakens the relevant activation field.
- [ ] Default output does not dump candidate or match internals.
- [ ] Debug output can still show generated forms and match details.
- [ ] `go test ./...` passes.

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
