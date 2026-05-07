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

Completed (see `docs/TODO_ARCHIVE.md` and git history for detail):
- Phase 3A cleanup, activation-energy discipline gate, graph and passage-field work.
- Knowledge validation/suggestion pipeline, harmonic field scoring, training/eval layer.
- Channel semantics correction, cross-script convergence, quality gate.
- Identity/provenance hardening, concept schema hygiene, ranking weights, train/held-out split.
- False activation reduction, review-file generation, human review/apply workflow.
- Knowledge layout consolidation, oversized test split, production responsibility audit.
- CLI ergonomics cleanup, architecture readiness review, documentation consolidation.
- Counsellor/transmutation fields: data-backed fields from transmute.yaml, --debug internals, tends-toward language, evidence paths with source/notes. Default concise; no hardcoded behavior.
- Harmonic data model hardening: frequency profile entries and labels now reject arbitrary unknown keys by allowlist; reference-only sections remain intentionally exempt and documented.
- Counsellor precision: natural-passage regression tests, direct field ranking over structural noise, deduplicated counsellor suggestions/evidence paths, neutral `feel` handling, structural debug labels, and concise prose alignment.
- Knowledge curation hygiene after counsellor expansion: validation warnings reduced from 183 to 149, no unknown concept warnings, no canonical-ID alias warnings, and accepted counsellor behavior preserved.
- All acceptance criteria green. Readiness report at `docs/ARCHITECTURE_READINESS.md`.
- Sanskrit/Devanagari support: exact Devanagari ScriptWord meanings now participate in the generic passage-field and harmonic-field path; known terms such as `प्राण`, `सत्य`, and `ॐ` activate direct verified fields where curated data exists, while unknown Devanagari such as `कवि` remains weak/speculative.
- Real-life Socrates testing: `docs/REAL_LIFE_TESTING.md` covers 22 seeker-style passages across emotional, spiritual, Sanskrit seed, mixed, edge-case, and infrastructure checks. The report is accepted with three follow-up observations: transliterated Sanskrit gaps, preposition noise, and a possible future resentment false-negative curation gap.

Current git status:

```text
## main...origin/main [ahead 87]
```

Key completed metrics (full detail in `docs/ARCHITECTURE_READINESS.md`):
- Train concept precision `0.24`, recall `0.96`.
- Train field precision `0.39`, recall `1.00`, pass rate `8/8`.
- Held-out split: no examples in current `./bin/socrates train` run.
- Validation: 0 errors, 149 warnings (data quality warnings only).
- Cross-script convergence: English, Norwegian, Hebrew, Chinese love on shared harmonic field.
- Active channel diversity, precision-aware training, multi-concept fuzzy evidence: all verified.
- Confidence is `verified`/`plausible`/`speculative`; provenance is `curated`/`traditional`/`human_review`/`physics`; `curated` is not confidence.
- Active runtime knowledge is `internal/knowledge/`; reference material is `knowledge/reference/`.
- Full architecture readiness report at `docs/ARCHITECTURE_READINESS.md`.

## Current Next Task

Task:
**Curate transliterated Sanskrit forms.**

Architect review status:
Real-life testing is accepted:
- `docs/REAL_LIFE_TESTING.md` exists and is committed
- 22 passages were tested across plain emotional, spiritual alignment, Sanskrit seeds, mixed Sanskrit/English, edge cases, and infrastructure checks
- direct emotional/spiritual evidence generally outranks structural noise
- known Sanskrit seed terms remain first-class through passage/harmonic fields
- unknown `कवि` remains weak/speculative
- tests, validation, and training remained green

Current blocker:
- Mixed English/Sanskrit passages fail for common Latin transliterations that are not yet in curated data. `prana` and `agni` work because they already exist in `forms.yaml`; `dharma`, `satya`, `karma`, and `jnana` do not.
- This is a knowledge curation task, not a code architecture task.
- Do not add new parser behavior, new fuzzy rules, new runtime channels, audio rendering, or broad new knowledge during this task.

Required investigation:
- Add curated Latin transliteration entries in `internal/knowledge/forms.yaml` for:
  - `dharma` -> `path`
  - `satya` -> `truth`
  - `karma` -> an existing appropriate concept only if a defensible existing concept exists; otherwise document why it is deferred
  - `jnana` -> an existing appropriate concept only if a defensible existing concept exists; otherwise document why it is deferred
- Use existing data shape: source, lens, confidence, and weight must be explicit.
- Prefer conservative confidence. If the mapping is traditional but broad, use `plausible`; use `verified` only when the current knowledge conventions clearly justify it.
- Do not create new concepts unless absolutely necessary. If a needed concept is absent, defer that form with a note instead of expanding the ontology casually.
- Add focused tests proving mixed English/Sanskrit inputs activate the curated transliterated form directly:
  - `I practice dharma every day` should surface `path`
  - `satya is the foundation of all practice` should surface `truth`
  - any added `karma`/`jnana` mapping must have a test
- Preserve existing Devanagari behavior for `धर्म`, `सत्य`, `प्राण`, and unknown `कवि`.
- Update `docs/REAL_LIFE_TESTING.md` only if the observed mixed-passage limitation changes materially.

Rules (unchanged):
- No hardcoded semantic word lists in production Go.
- No direct behavior for specific example passages.
- No therapeutic, medical, prophetic, or absolute truth claims.
- All correction mappings remain in curated data with source/lens/confidence.
- Default output stays concise; detailed counsellor internals and evidence paths may be debug output.
- Training/review artifacts must not silently mutate active runtime knowledge.

Reference:
- `docs/ARCHITECTURE_READINESS.md`
- `docs/KNOWLEDGE_CURATION.md`
- `TRAINING_MODEL.md`
- `HARMONIC_DATA_MODEL.md`
- `FREQUENCY_MODEL.md`
- `IMPLEMENTATION_ROADMAP.md` section "Long-Term Counsellor Vision"

Acceptance Criteria:

- Curated transliterated Sanskrit form entries are added only where there is a defensible existing target concept.
- `dharma` activates `path` directly in mixed English/Sanskrit passage analysis.
- `satya` activates `truth` directly in mixed English/Sanskrit passage analysis.
- `karma` and `jnana` are either curated with tests or explicitly deferred with a short reason.
- Devanagari Sanskrit seed behavior remains unchanged.
- Unknown Devanagari input remains weak/speculative.
- No broad new knowledge area or runtime channel is introduced.
- No production Go semantic word maps or word-specific branches are introduced.
- `./bin/socrates knowledge validate` passes with no errors.
- `go test ./...` passes.
- Work is committed locally and not pushed.


## Next After This

After transliterated Sanskrit curation:
- Preposition handling for standalone fragments such as `in`
- Harmonic/audio rendering layer
- Consult `docs/ARCHITECTURE_READINESS.md` before major new feature layers

## Nice To Have Later

- extended script support
- harmonic/audio rendering layer
