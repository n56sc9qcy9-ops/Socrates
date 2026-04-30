package decipher

import (
	"socrates/internal/knowledge"
)

// =============================================================================
// Phase G: Scoring (with deduplication)
// =============================================================================

// CalculateScoreComponents computes detailed score components.
// Uses ConvergenceResult for generic concept activation instead of semantic buckets.
// IMPORTANT: Deduplicates evidence by meaningful path before scoring to prevent inflation.
func CalculateScoreComponents(
	candidates []CandidateForm,
	matches []MatchEvidence,
	expansions map[string][]knowledge.DecipherConceptRelation,
	convergence ConvergenceResult,
	channels []ChannelResult,
) ScoreComponents {
	components := ScoreComponents{}

	// Deduplicate matches by evidence ID - keep first (best) occurrence
	dedupedMatches := deduplicateMatchEvidence(matches)

	// Exact match score: high weight for exact matches
	var exactWeight, fuzzyWeight float64

	for _, m := range dedupedMatches {
		if m.Distance == 0 {
			exactWeight += m.Weight
		} else {
			fuzzyWeight += m.Weight
		}
	}

	if len(dedupedMatches) > 0 {
		components.ExactMatchScore = exactWeight / float64(len(dedupedMatches))
		components.FuzzyMatchScore = fuzzyWeight / float64(len(dedupedMatches)) * 0.8 // Fuzzy is weaker
	}

	// Graph expansion score - normalize by count AND cap at 1.0
	var expansionWeight float64
	var expansionCount int
	for _, exps := range expansions {
		for _, e := range exps {
			expansionWeight += e.Weight
			expansionCount++
		}
	}
	if expansionCount > 0 {
		components.GraphExpansionScore = expansionWeight / float64(expansionCount) * 0.7
		if components.GraphExpansionScore > 1.0 {
			components.GraphExpansionScore = 1.0
		}
	}

	// Passage convergence score - derived from generic activation, not semantic buckets
	// Deduplicate activated concepts first
	dedupedConcepts := deduplicateActivatedConcepts(convergence.ActivatedConcepts)
	if len(dedupedConcepts) > 0 {
		// Combine co-activation score with number of activated concepts
		conceptCount := float64(len(dedupedConcepts))
		components.PassageConvergenceScore = convergence.CoActivationScore*0.6 + (conceptCount/10.0)*0.4
		// Normalize
		if components.PassageConvergenceScore > 1.0 {
			components.PassageConvergenceScore = 1.0
		}
	}

	// Multi-method bonus: deduplicated methods agreeing
	methodSet := make(map[string]bool)
	for _, m := range dedupedMatches {
		methodSet[m.Method] = true
	}
	if len(methodSet) >= 3 {
		components.MultiMethodBonus = 0.15
	} else if len(methodSet) >= 2 {
		components.MultiMethodBonus = 0.1
	}

	// Channel diversity bonus - only count channels with ACTIVE meaningful signals
	activeChannelCount := countActiveChannels(channels)
	if activeChannelCount >= 4 {
		components.ChannelDiversityBonus = 0.2
	} else if activeChannelCount >= 3 {
		components.ChannelDiversityBonus = 0.1
	}

	return components
}

// countActiveChannels counts only channels with meaningful signals.
// A channel is active only if it has at least one signal with a non-empty target
// or non-trivial evidence (not just "no match" or structural observations).
func countActiveChannels(channels []ChannelResult) int {
	activeCount := 0
	for _, ch := range channels {
		if len(ch.Signals) == 0 {
			// Empty channel - not active
			continue
		}
		// Check if channel has meaningful signals
		hasMeaningful := false
		for _, sig := range ch.Signals {
			// Skip no-match signals and purely structural observations
			if sig.Target == "" || sig.Target == "no-match" || sig.Target == "phonetic-observation" {
				continue
			}
			// Skip orthographic-only signals with low weight
			if sig.Lens == "orthographic" && sig.Weight < 0.3 {
				continue
			}
			hasMeaningful = true
			break
		}
		if hasMeaningful {
			activeCount++
		}
	}
	return activeCount
}

// deduplicateMatchEvidence removes duplicate MatchEvidence entries by EvidenceID.
// Keeps the first (best) occurrence of each evidence path.
func deduplicateMatchEvidence(matches []MatchEvidence) []MatchEvidence {
	if len(matches) == 0 {
		return matches
	}
	seen := make(map[string]bool)
	result := make([]MatchEvidence, 0, len(matches))
	for _, m := range matches {
		id := m.EvidenceID()
		if !seen[id] {
			seen[id] = true
			result = append(result, m)
		}
	}
	return result
}

// deduplicateActivatedConcepts removes duplicate ActivatedConcept entries by Concept.
// Keeps the strongest occurrence of each concept.
func deduplicateActivatedConcepts(concepts []ActivatedConcept) []ActivatedConcept {
	if len(concepts) == 0 {
		return concepts
	}
	best := make(map[string]ActivatedConcept)
	for _, c := range concepts {
		if existing, ok := best[c.Concept]; ok {
			// Keep strongest
			if c.Strength > existing.Strength {
				best[c.Concept] = c
			}
		} else {
			best[c.Concept] = c
		}
	}
	result := make([]ActivatedConcept, 0, len(best))
	for _, c := range best {
		result = append(result, c)
	}
	return result
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
