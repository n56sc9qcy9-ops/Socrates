package decipher

import (
	"socrates/internal/knowledge"
)

// runSymbolicChannel maps to symbolic neighborhoods using knowledge base.
func runSymbolicChannel(forms Forms, kb *knowledge.Knowledge) ChannelResult {
	signals := make([]Signal, 0)

	// Determine if the input token is a standalone function/preposition word.
	// For multi-token input: do NOT mark as standalone since the normalized
	// form is the full passage, not a single token.
	// Only single-token input gets the standalone flag.
	isStandalone := len(forms.Tokens) == 1 && IsPrepositionOrFunctionWord(forms.Normalized)

	// Use knowledge-based primitive lookup for exact matches
	prims := knowledgeBasedPrimitiveLookup(forms.Normalized, kb)

	for _, prim := range prims {
		// Generate neighbor signals from the primitive's neighbors
		for _, neighbor := range prim.Neighbors {
			signals = append(signals, Signal{
				Text:              prim.Name + " -> neighbor: " + neighbor,
				Target:            neighbor,
				Channel:           "Symbolic",
				Lens:              "symbolic",
				Confidence:        ConfidencePlausible,
				Weight:            0.4,
				IsStandaloneToken: isStandalone,
			})
		}
		signals = append(signals, Signal{
			Text:              "primitive match: " + prim.Name,
			Target:            prim.ID,
			Channel:           "Symbolic",
			Lens:              "symbolic",
			Confidence:        ConfidenceVerified,
			Weight:            0.7,
			IsDirect:          true, // Direct primitive match from symbolic knowledge
			IsStandaloneToken: isStandalone,
		})
	}

	score := calculateChannelScore(signals)

	return ChannelResult{
		Name:    "Symbolic",
		Signals: signals,
		Score:   score,
	}
}
