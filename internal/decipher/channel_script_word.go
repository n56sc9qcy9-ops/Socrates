package decipher

import (
	"socrates/internal/knowledge"
)

// runScriptWordChannel performs exact whole-word matching for non-Latin scripts.
// This channel ONLY produces signals for EXACT matches of known ScriptWords.
// It does NOT infer meaning for unknown words.
// Uses knowledge base exclusively.
func runScriptWordChannel(forms Forms, kb *knowledge.Knowledge) ChannelResult {
	signals := make([]Signal, 0)

	// Only process non-Latin scripts
	script := forms.Script
	if script == ScriptLatin {
		return ChannelResult{Name: "ScriptWord", Signals: signals, Score: 0.0}
	}

	// Check for exact ScriptWord match from knowledge base
	matchedWords := knowledgeBasedScriptWordLookup(string(script), forms.Runes, kb)

	for _, seed := range matchedWords {
		for _, meaning := range seed.Meanings {
			signals = append(signals, Signal{
				Text:       "exact word match: " + seed.Word,
				Target:     meaning,
				Channel:    "ScriptWord",
				Lens:       seed.Script + "-word",
				Confidence: seed.Confidence,
				Weight:     seed.Weight,
			})
		}
	}

	score := calculateChannelScore(signals)
	return ChannelResult{Name: "ScriptWord", Signals: signals, Score: score}
}
