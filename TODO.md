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
- Phase 5 activation graph hardening.
- Phase 6 passage fields integrated into `Reading`.

Current git status reported after Pi's last task:

```text
## main...origin/main [ahead 11]
```

Architect review notes:

- Runtime knowledge is loaded from embedded `internal/knowledge/*.yaml`, not automatically from root `knowledge/`.
- More concepts must enter through a curated, validated data pipeline. The app must not silently learn trusted concepts from user input.
- Preflight bug to fix: `isHigherConfidence()` / merge usage in `passage_field.go` can downgrade confidence during merge because the argument semantics are inverted.

## Current Next Task

Task:
Build a safe knowledge-growth pipeline: validate, import, and suggest knowledge without auto-promoting untrusted input.

Context:
The engine is now graph-backed and passage-aware, but the knowledge base is still hand-edited YAML. That is acceptable for early development, but the project needs a disciplined path for adding concepts, forms, relations, script words, and glyph patterns. The pipeline must keep knowledge auditable, data-driven, and test-protected.

Likely files:

- `internal/knowledge/knowledge.go`
- `internal/knowledge/loader.go`
- `internal/knowledge/knowledge_test.go`
- `internal/knowledge/*.yaml`
- `knowledge/*.yaml`
- `cmd/socrates/main.go`
- `internal/decipher/passage_field.go`
- new package/file if useful: `internal/knowledge/validate.go`
- new package/file if useful: `internal/knowledge/suggestions.go`
- new docs if useful: `docs/KNOWLEDGE_CURATION.md`

Instructions:

1. Start with `git status --short --branch` and record it.
2. Run `go test ./...` before changing code and record the baseline result.
3. Fix the `isHigherConfidence` merge bug in `passage_field.go` first and add a regression test proving verified confidence is not downgraded by merging plausible/speculative evidence.
4. Decide and document the authoritative runtime data directory. Current runtime uses embedded `internal/knowledge/*.yaml`; root `knowledge/*.yaml` must either be documented as non-runtime source material or kept in sync by tests/tooling.
5. Add a validation layer for knowledge data. It should check at least:
   - concept IDs are unique and non-empty
   - concept aliases do not ambiguously point to multiple concepts unless explicitly allowed
   - every form/script word/glyph target concept exists, or is explicitly marked external/speculative by policy
   - every relation endpoint exists
   - weights are in `(0, 1]`
   - confidence is one of `verified`, `plausible`, `speculative`
   - duplicate forms/relations are either rejected or deterministic
6. Add CLI support for validation, likely `socrates knowledge validate`, without redesigning the whole CLI.
7. Add an import path for curated knowledge packs, likely `socrates knowledge validate --dir <path>` first. Do not auto-write into embedded YAML in this task unless the write path is explicit and safe.
8. Add a suggestion mechanism for unknown or weakly matched inputs. Suggestions may be written to a separate review file or printed as structured output, but they must not become trusted concepts automatically.
9. Define the minimal review record for future additions: proposed concept/form/relation, evidence source, confidence, weight, rationale/notes, and whether it was accepted.
10. Add tests for validation failures: missing relation endpoint, form target missing concept, invalid confidence, invalid weight, duplicate concept ID, ambiguous alias.
11. Add tests proving suggestions are separate from trusted embedded knowledge.
12. Keep all decipher/render/passage/graph tests passing.
13. Do not import a large external dictionary in this task. Build the pipeline first.
14. Do not add auto-learning from user input into trusted YAML.
15. Do not implement harmonic/audio rendering or concept-to-frequency mappings.
16. Run targeted tests while working, then `go test ./...`.
17. Commit the completed knowledge-pipeline work locally with a clear message. Do not push.
18. Report final `git status --short --branch`, tests run, files changed, and how a human curator would add a new concept after this task.

Pi work requirement:
Spend at least 60 focused minutes. If validation is finished early, use the remaining time to add stronger failure tests and a concise curation doc.

Acceptance Criteria:

- Confidence merge bug is fixed and tested.
- The authoritative runtime knowledge location is documented.
- Knowledge validation exists and is covered by failure tests.
- CLI can validate embedded or directory-based knowledge.
- Suggestions are separated from trusted knowledge.
- New concepts have a documented human-curation path.
- Existing tests still pass.
- `go test ./...` passes.
- Work is committed locally and not pushed.

## Next After This

If this passes, the next task should add the first curated knowledge pack using the new pipeline, not by ad hoc YAML edits.
