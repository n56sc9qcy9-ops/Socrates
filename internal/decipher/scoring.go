package decipher

import (
	"socrates/internal/knowledge"
)

// =============================================================================
// Phase G: Scoring
// =============================================================================

// CalculateScoreComponents computes detailed score components.
// Uses ConvergenceResult for generic concept activation instead of semantic buckets.
func CalculateScoreComponents(
	candidates []CandidateForm,
	matches []MatchEvidence,
	expansions map[string][]knowledge.DecipherConceptRelation,
	convergence ConvergenceResult,
	channelCount int,
) ScoreComponents {
	components := ScoreComponents{}

	// Exact match score: high weight for exact matches
	var exactWeight, fuzzyWeight float64

	for _, m := range matches {
		if m.Distance == 0 {
			exactWeight += m.Weight
		} else {
			fuzzyWeight += m.Weight
		}
	}

	if len(matches) > 0 {
		components.ExactMatchScore = exactWeight / float64(len(matches))
		components.FuzzyMatchScore = fuzzyWeight / float64(len(matches)) * 0.8 // Fuzzy is weaker
	}

	// Graph expansion score
	var expansionWeight float64
	for _, exps := range expansions {
		for _, e := range exps {
			expansionWeight += e.Weight
		}
	}
	if len(expansions) > 0 {
		components.GraphExpansionScore = expansionWeight / float64(len(expansions)) * 0.7
	}

	// Passage convergence score - derived from generic activation, not semantic buckets
	if len(convergence.ActivatedConcepts) > 0 {
		// Combine co-activation score with number of activated concepts
		conceptCount := float64(len(convergence.ActivatedConcepts))
		components.PassageConvergenceScore = convergence.CoActivationScore*0.6 + (conceptCount/10.0)*0.4
		// Normalize
		if components.PassageConvergenceScore > 1.0 {
			components.PassageConvergenceScore = 1.0
		}
	}

	// Multi-method bonus: multiple methods agreeing
	methodSet := make(map[string]bool)
	for _, m := range matches {
		methodSet[m.Method] = true
	}
	if len(methodSet) >= 3 {
		components.MultiMethodBonus = 0.15
	} else if len(methodSet) >= 2 {
		components.MultiMethodBonus = 0.1
	}

	// Channel diversity bonus
	if channelCount >= 4 {
		components.ChannelDiversityBonus = 0.2
	} else if channelCount >= 3 {
		components.ChannelDiversityBonus = 0.1
	}

	return components
}

// CalculateFinalScore computes the final combined score.
func CalculateFinalScore(components ScoreComponents) float64 {
	score := components.ExactMatchScore*0.25 +
		components.FuzzyMatchScore*0.20 +
		components.GraphExpansionScore*0.15 +
		components.PassageConvergenceScore*0.25 +
		components.MultiMethodBonus +
		components.ChannelDiversityBonus

	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}
	return score
}
