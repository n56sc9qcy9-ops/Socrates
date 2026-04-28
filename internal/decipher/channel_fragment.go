package decipher

import (
	"socrates/internal/knowledge"
)

// runFragmentChannel analyzes fragment paths.
// For Latin script: includes exact whole-token matching.
// For non-Latin scripts: fragment matching is not applicable (see runScriptWordChannel).
// Uses knowledge base exclusively for data-driven lookups.
func runFragmentChannel(forms Forms, kb *knowledge.Knowledge) ChannelResult {
	signals := make([]Signal, 0)

	// First: exact whole-token matching
	exactSignals := runWholeTokenMatching(forms, kb)
	signals = append(signals, exactSignals...)

	// Second: fragment path analysis (knowledge base only)
	for _, path := range forms.Fragments {
		for _, part := range path.Parts {
			// Lookup fragment seeds from knowledge base
			seeds := knowledgeBasedFragmentLookup(part, kb)
			for _, seed := range seeds {
				for _, lens := range seed.Lenses {
					signals = append(signals, Signal{
						Text:       "fragment '" + part + "' -> " + lens.Target,
						Target:     lens.Target,
						Channel:    "Fragment",
						Lens:       lens.Lens,
						Confidence: lens.Confidence,
						Weight:     lens.Weight(path.Confidence),
					})
				}
			}

			// Lookup primitives from knowledge base
			prims := knowledgeBasedPrimitiveLookup(part, kb)
			for _, prim := range prims {
				signals = append(signals, Signal{
					Text:       "primitive '" + part + "' matches " + prim.Name,
					Target:     prim.ID,
					Channel:    "Fragment",
					Lens:       "primitive",
					Confidence: ConfidenceVerified,
					Weight:     0.6,
				})
			}
		}
	}

	score := calculateChannelScore(signals)

	return ChannelResult{
		Name:    "Fragment",
		Signals: signals,
		Score:   score,
	}
}

// runWholeTokenMatching checks if the entire input matches a known word/seed.
// This provides exact whole-token matching for words like prana, ruach, logos, mantra.
// ONLY processes Latin script - non-Latin scripts use runScriptWordChannel.
func runWholeTokenMatching(forms Forms, kb *knowledge.Knowledge) []Signal {
	signals := make([]Signal, 0)

	// Only process Latin script - non-Latin scripts use runScriptWordChannel
	script := forms.Script
	if script != ScriptLatin {
		return signals
	}

	input := forms.Normalized

	// Check fragment seeds for exact whole-token match from knowledge base
	seeds := knowledgeBasedFragmentLookup(input, kb)
	for _, seed := range seeds {
		if seed.Fragment == input {
			for _, lens := range seed.Lenses {
				signals = append(signals, Signal{
					Text:       "exact whole-token match: " + input,
					Target:     lens.Target,
					Channel:    "Fragment",
					Lens:       lens.Lens,
					Confidence: lens.Confidence,
					Weight:     lens.BaseWeight,
				})
			}
		}
	}

	return signals
}

// Weight adjusts weight by a multiplier (used by fragment confidence).
func (l FragmentLens) Weight(multiplier float64) float64 {
	weight := l.BaseWeight * multiplier
	if weight > 1.0 {
		return 1.0
	}
	return weight
}
