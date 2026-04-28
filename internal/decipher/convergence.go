package decipher

import (
	"socrates/internal/knowledge"
)

// =============================================================================
// Phase E: Convergence Result Logic
// =============================================================================

// ConvergenceResult represents the result of convergence computation.
type ConvergenceResult struct {
	ActivatedConcepts []ActivatedConcept
	CoActivationScore float64
	TopConcepts       []ActivatedConcept
	RelationPaths     []string
}

// DetectConvergence computes passage-level convergence using generic concept activation.
// NO semantic bucket booleans - convergence emerges from graph structure.
func DetectConvergence(passageSignals []PassageSignal, directConcepts []string, kb *knowledge.Knowledge) ConvergenceResult {
	// Step 1: Compute activated concepts from fuzzy matches
	activated := ComputeActivatedConcepts(passageSignals, directConcepts, kb)

	// Step 2: Compute co-activation score from relation structure
	coActivation := ComputeCoActivation(activated, kb)

	// Step 3: Find top concepts by strength
	topConcepts := findTopConcepts(activated, 3)

	// Step 4: Identify relation paths between top concepts using knowledge base
	conceptIDs := make([]string, len(topConcepts))
	for i, tc := range topConcepts {
		conceptIDs[i] = tc.Concept
	}
	paths := getRelationPathsForConcepts(conceptIDs, kb)

	return ConvergenceResult{
		ActivatedConcepts: activated,
		CoActivationScore: coActivation,
		TopConcepts:       topConcepts,
		RelationPaths:     paths,
	}
}

// findTopConcepts returns the top N concepts by activation strength.
func findTopConcepts(activated []ActivatedConcept, n int) []ActivatedConcept {
	// Sort by strength descending
	sorted := make([]ActivatedConcept, len(activated))
	copy(sorted, activated)

	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Strength > sorted[i].Strength {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	if len(sorted) > n {
		return sorted[:n]
	}
	return sorted
}
