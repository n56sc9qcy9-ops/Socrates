package decipher

import (
	"sort"

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
// Priority: direct evidence > by strength descending.
// This ensures the prose reading reflects verified field activations, not just signal strength.
func findTopConcepts(activated []ActivatedConcept, n int) []ActivatedConcept {
	if len(activated) <= n {
		return activated
	}

	sorted := make([]ActivatedConcept, len(activated))
	copy(sorted, activated)

	// Sort: direct evidence first, then by strength descending
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].IsDirectEvidence != sorted[j].IsDirectEvidence {
			return sorted[i].IsDirectEvidence // true before false
		}
		return sorted[i].Strength > sorted[j].Strength
	})

	if len(sorted) > n {
		return sorted[:n]
	}
	return sorted
}

// findTopConceptsWithFields returns the top N concepts, preferring concepts
// that appear in the top passage fields when available.
// Priority: direct evidence > in-field > by strength.
func findTopConceptsWithFields(activated []ActivatedConcept, topFields []string, n int) []ActivatedConcept {
	if len(activated) <= n {
		return activated
	}

	// Build a set of top field concept IDs for fast lookup
	fieldSet := make(map[string]bool)
	for _, f := range topFields {
		fieldSet[f] = true
	}

	// Sort: direct evidence first, then in-field, then by strength descending
	sorted := make([]ActivatedConcept, len(activated))
	copy(sorted, activated)

	sort.Slice(sorted, func(i, j int) bool {
		// 1. Direct evidence first
		if sorted[i].IsDirectEvidence != sorted[j].IsDirectEvidence {
			return sorted[i].IsDirectEvidence
		}
		// 2. In top fields (by strength)
		inFieldI := fieldSet[sorted[i].Concept]
		inFieldJ := fieldSet[sorted[j].Concept]
		if inFieldI != inFieldJ {
			return inFieldI
		}
		// 3. By strength descending
		return sorted[i].Strength > sorted[j].Strength
	})

	if len(sorted) > n {
		return sorted[:n]
	}
	return sorted
}
