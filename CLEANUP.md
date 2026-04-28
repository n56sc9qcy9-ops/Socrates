# CLEANUP: Active Pi Task

This is the only active instruction file for Pi.

Do not edit `TODO.md` for this task.

Task:
Split `internal/decipher/channels.go` into focused channel files without changing behavior.

Context:
Phase 3A is still blocked because `channels.go` is about 668 lines and mixes orchestration, glyph analysis, sound analysis, script-word matching, fragment matching, cross-language matching, symbolic matching, and scoring helpers. This task is mechanical architecture cleanup only. Move code by whole function/type blocks so later cleanup can target individual channels safely.

Files:
- `internal/decipher/channels.go`
- `internal/decipher/channel_glyph.go`
- `internal/decipher/channel_sound.go`
- `internal/decipher/channel_script_word.go`
- `internal/decipher/channel_fragment.go`
- `internal/decipher/channel_cross_language.go`
- `internal/decipher/channel_symbolic.go`
- `internal/decipher/channel_scoring.go`
- existing tests under `internal/decipher/*_test.go`

Instructions:
1. Keep `RunAllChannels` in `internal/decipher/channels.go`.
2. Create `internal/decipher/channel_glyph.go` and move:
   - `runGlyphChannel`
   - `analyzeLatinGlyphs`
   - `analyzeHebrewGlyphs`
   - `analyzeDevanagariGlyphs`
   - `analyzeHanGlyphs`
3. Create `internal/decipher/channel_sound.go` and move:
   - `runSoundChannel`
4. Create `internal/decipher/channel_script_word.go` and move:
   - `runScriptWordChannel`
5. Create `internal/decipher/channel_fragment.go` and move:
   - `runFragmentChannel`
   - `runWholeTokenMatching`
   - `FragmentLens.Weight`
6. Create `internal/decipher/channel_cross_language.go` and move:
   - `runCrossLanguageChannel`
7. Create `internal/decipher/channel_symbolic.go` and move:
   - `runSymbolicChannel`
8. Create `internal/decipher/channel_scoring.go` and move:
   - `calculateChannelScore`
   - `itoa`
9. Clean imports in every file.
   - Glyph and sound files likely need `strings`.
   - Files that accept `*knowledge.Knowledge` need `socrates/internal/knowledge`.
10. Run formatting:

```bash
gofmt -w internal/decipher/*.go
```

11. Run tests:

```bash
go test ./...
```

Constraints:
- Do not change behavior.
- Do not change public function names.
- Do not change exported type names.
- Do not edit `TODO.md`.
- Do not refactor algorithms while moving code.
- Do not move hardcoded semantic maps into YAML in this task; leave that for a separate task after the split.
- Do not remove root artifacts.
- Do not report Phase 3A complete.

Tests:
- Existing tests must continue to pass.
- If any architecture test has a production-file list, update it to include the new `channel_*.go` files and remove `channels.go` only if the checked code no longer lives there.
- Do not add broad new behavior tests for channel semantics in this task.

Acceptance Criteria:
- `channels.go` is reduced to orchestration only, primarily `RunAllChannels` plus package/import/comment as needed.
- Each channel implementation lives in its focused `channel_*.go` file.
- No duplicate function/type declarations remain.
- No production behavior changes are made.
- Architecture regression tests still scan the current production file set.
- `go test ./...` passes.

Do Not:
- Do not edit `TODO.md`.
- Do not remove `skal` or `socrates`.
- Do not migrate hardcoded cross-language data in this task.
- Do not combine this with `engine.go` cleanup.
- Do not mark Phase 3A complete.
