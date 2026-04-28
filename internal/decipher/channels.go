package decipher

import (
	"socrates/internal/knowledge"
)

// channels.go holds the deprecated monolithic channel orchestration.
// Logic has been moved into focused files:
// - channel_glyph.go: runGlyphChannel, analyzeLatinGlyphs, etc.
// - channel_sound.go: runSoundChannel
// - channel_script_word.go: runScriptWordChannel
// - channel_fragment.go: runFragmentChannel, runWholeTokenMatching, FragmentLens.Weight
// - channel_cross_language.go: runCrossLanguageChannel
// - channel_symbolic.go: runSymbolicChannel
// - channel_scoring.go: calculateChannelScore, itoa

// RunAllChannels runs all resonance channels and returns results.
// Takes explicit knowledge parameter instead of global state.
func RunAllChannels(forms Forms, kb *knowledge.Knowledge) []ChannelResult {
	results := make([]ChannelResult, 0)

	// Glyph channel
	glyphResult := runGlyphChannel(forms, kb)
	results = append(results, glyphResult)

	// Sound channel
	soundResult := runSoundChannel(forms)
	results = append(results, soundResult)

	// ScriptWord channel (exact whole-word matching for non-Latin scripts)
	scriptWordResult := runScriptWordChannel(forms, kb)
	results = append(results, scriptWordResult)

	// Fragment channel (uses knowledge data when available)
	fragmentResult := runFragmentChannel(forms, kb)
	results = append(results, fragmentResult)

	// Cross-language channel
	crossResult := runCrossLanguageChannel(forms)
	results = append(results, crossResult)

	// Symbolic channel
	symbolicResult := runSymbolicChannel(forms, kb)
	results = append(results, symbolicResult)

	return results
}
