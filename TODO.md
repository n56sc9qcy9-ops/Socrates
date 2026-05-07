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
- Sanskrit/Devanagari support: exact Devanagari ScriptWord meanings now participate in the generic passage-field and harmonic-field path; known terms such as `प्राण`, `सत्य`, and `ॐ` activate direct verified fields where curated data exists, while unknown Devanagari such as `कवि` remains weak/speculative. Arabic and Thai were not added.

Current git status:

```text
## main...origin/main [ahead 85]
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
**Real-life Socrates testing after Sanskrit/Devanagari support.**

Architect review status:
Sanskrit/Devanagari support is accepted after the follow-up fix:
- exact ScriptWord meanings are indexed as Forms in both builder and rebuilt-index paths
- non-Latin whole-word inputs remain whole tokens for passage/harmonic activation
- `प्राण` activates direct verified fields for `prana`, `breath`, and `life-force` and produces a harmonic field
- `सत्य` activates direct verified `truth` and produces a harmonic field
- `कवि` remains weak/speculative and does not produce a harmonic field
- Hebrew/Han existing script behavior and all prior counsellor/harmonic checks still pass
- no Arabic or Thai runtime knowledge was added

Current blocker:
- None in the code path. The next risk is whether Socrates feels spiritually useful, grounded, and non-authoritarian in real usage.
- Do not add Arabic, Thai, new language families, audio rendering, or large new knowledge during this task.

Required investigation:
- Create a small real-life testing protocol/report. Prefer a concise doc such as `docs/REAL_LIFE_TESTING.md` unless a better existing doc already owns this.
- Test Socrates manually with a fixed set of seeker-style prompts, including:
  - plain emotional passages: loneliness, sadness, fear, emptiness, resentment, longing, peace
  - spiritual alignment passages: truth, love, light, surrender, trust, faith, humility
  - Sanskrit seeds: `प्राण`, `सत्य`, `ॐ`, `धर्म`, `अग्नि`, `मंत्र`
  - mixed passages, for example English emotional text containing one Sanskrit term
- For each case, record:
  - top fields and whether direct evidence outranks structural noise
  - whether the reading is useful without claiming absolute truth about the person
  - whether suggestions remain evidence-backed and confidence-labeled
  - whether harmonic fields appear only when curated profiles support them
  - whether unknown or weak inputs remain humble/speculative
- Treat failures as review artifacts first. Do not silently add runtime knowledge to make examples pass.
- If a serious failure appears, stop and document the smallest correction needed.
- If the testing report is green, recommend the next architecture choice: deepen real-life testing, start harmonic/audio rendering, or consider Arabic later.

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

- A real-life testing report/protocol exists and is committed.
- The report includes the exact prompts tested and representative CLI outputs or summaries.
- The report separates spiritual usefulness from evidence support; Socrates must not present itself as absolute truth about the person.
- Direct emotional/spiritual evidence continues to outrank structural noise in plain passages.
- Known Sanskrit terms remain first-class through passage/harmonic fields.
- Unknown or weak inputs remain weak/speculative.
- No Arabic, Thai, or broad new script family is added.
- No runtime knowledge is silently mutated as part of testing unless a documented blocker requires a minimal fix.
- `./bin/socrates knowledge validate` passes with no errors.
- `go test ./...` passes.
- Work is committed locally and not pushed.


## Next After This

After Sanskrit/Devanagari support:
- Real-life Socrates testing before adding another script family
- Arabic later if real-life testing supports adding it
- Harmonic/audio rendering layer
- Consult `docs/ARCHITECTURE_READINESS.md` before major new feature layers

## Nice To Have Later

- extended script support (Arabic, Thai, etc.)
- harmonic/audio rendering layer
