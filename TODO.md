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
- Phase 5/6 graph and passage-field work.
- First knowledge validation/suggestion pipeline.

Current git status reported after Pi's last task:

```text
## main...origin/main [ahead 15]
```

Architect review of knowledge pipeline:

- Good: validation and suggestion packages exist, curation doc exists, and trusted knowledge is still embedded from `internal/knowledge/*.yaml`.
- Must fix next: CLI docs say `socrates knowledge validate`, but implementation appears to use a `--validate` flag internally. The documented command must actually work.
- Must fix next: alias ambiguity is documented but not fully validated across concepts.
- Must fix next: form/glyph targets for missing concepts are warnings in places where the task expected strict policy. Decide strict vs explicit external/speculative policy and test it.
- Must fix next: suggestions exist as package functions, but no CLI path exposes them; docs imply `decipher unknown_word > suggestions.txt`, which does not actually produce suggestion records.
- Should fix next: `ftos` is a lossy custom formatter and should use standard formatting.
- Should fix next: `TestValidateKnowledge_ValidKnowledge` should assert embedded knowledge has no validation errors, not only that validation runs.

## Current Next Task

Task:
Harden the knowledge pipeline UX and strictness, then add the first small curated knowledge pack through it.

Context:
The pipeline exists, but the command-line behavior, validator strictness, suggestion exposure, and tests need to match the documented curation workflow. After hardening, add one deliberately small curated pack using the pipeline rather than ad hoc expansion.

Likely files:

- `cmd/socrates/main.go`
- `internal/knowledge/validate.go`
- `internal/knowledge/validate_test.go`
- `internal/knowledge/suggestions.go`
- `internal/knowledge/suggestions_test.go`
- `internal/knowledge/*.yaml`
- `docs/KNOWLEDGE_CURATION.md`
- `TODO.md`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Fix the CLI so the documented command works exactly:
   - `socrates knowledge validate`
   - `socrates knowledge validate --dir <path>`
4. Add tests or a small command parser helper test for the knowledge CLI behavior if practical.
5. Add a suggestion CLI path, for example:
   - `socrates knowledge suggest <word-or-phrase>`
   - It should print review records only; it must not mutate trusted YAML.
6. Update docs so they do not claim normal `decipher` output emits suggestions unless that is actually implemented.
7. Make alias ambiguity validation real. Duplicate aliases across different concept IDs must be rejected or explicitly reported according to policy.
8. Make target-reference validation strict and documented:
   - trusted forms, script words, relations, and glyph patterns should target existing concepts
   - if external/speculative targets are allowed, represent that policy explicitly and test it
9. Replace custom `ftos` with standard formatting.
10. Strengthen validator tests so embedded runtime knowledge must have zero validation errors.
11. Add tests proving suggestions remain separate from trusted embedded knowledge and from YAML mutation.
12. Add the first small curated knowledge pack through the validated path. Keep it narrow: 2-4 concepts/forms/relations max, with review notes in docs or comments. Prefer concepts that improve passage-field examples without adding canned readings.
13. Re-run validation after adding the pack.
14. Keep all decipher/render/passage/graph tests passing.
15. Do not import a large dictionary.
16. Do not auto-learn from user input.
17. Do not implement harmonic/audio rendering or concept-to-frequency mappings.
18. Run targeted tests while working, then `go test ./...`.
19. Commit the completed pipeline-hardening/curated-pack work locally with a clear message. Do not push.
20. Report final `git status --short --branch`, tests run, files changed, exact CLI examples, and what curated knowledge was added.

Pi work requirement:
Spend at least 60 focused minutes. If the CLI fixes are quick, use the remaining time on strict validation tests and the first curated pack.

Acceptance Criteria:

- Documented validation CLI works as written.
- Suggestion CLI exists and does not mutate trusted knowledge.
- Alias ambiguity is validated.
- Missing target policy is strict, explicit, and tested.
- Embedded runtime knowledge validates with zero errors.
- Curation docs match real CLI behavior.
- One small curated knowledge pack is added through the validated workflow.
- Existing tests still pass.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

If this passes, the next task should be Phase 7: evidence-first CLI polish using the larger, validated knowledge base.
