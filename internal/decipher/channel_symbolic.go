package decipher

import (
	"socrates/internal/knowledge"
)

// runSymbolicChannel maps to symbolic neighborhoods using knowledge base.
func runSymbolicChannel(forms Forms, kb *knowledge.Knowledge) ChannelResult {
	signals := make([]Signal, 0)

	// Use knowledge-based primitive lookup for exact matches
	prims := knowledgeBasedPrimitiveLookup(forms.Normalized, kb)

	for _, prim := range prims {
		// Generate neighbor signals from the primitive's neighbors
		for _, neighbor := range prim.Neighbors {
			signals = append(signals, Signal{
				Text:       prim.Name + " -> neighbor: " + neighbor,
				Target:     neighbor,
				Channel:    "Symbolic",
				Lens:       "symbolic",
				Confidence: ConfidencePlausible,
				Weight:     0.4,
			})
		}
		signals = append(signals, Signal{
			Text:       "primitive match: " + prim.Name,
			Target:     prim.ID,
			Channel:    "Symbolic",
			Lens:       "symbolic",
			Confidence: ConfidenceVerified,
			Weight:     0.7,
		})
	}

	score := calculateChannelScore(signals)

	return ChannelResult{
		Name:    "Symbolic",
		Signals: signals,
		Score:   score,
	}
}
