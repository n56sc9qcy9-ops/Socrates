package decipher

import (
	"socrates/internal/knowledge"
)

// ConceptRelation represents a relation between concepts in the concept graph.
type ConceptRelation struct {
	From       string
	To         string
	Type       string
	Confidence string
	Weight     float64
}

// =============================================================================
// Phase D: Data-Driven Concept Graph Expansion
// =============================================================================

// ExpandConcept expands a concept to its related concepts using the knowledge graph.
// This replaces the hardcoded ConceptGraph.
func ExpandConcept(concept string, minWeight float64, kb *knowledge.Knowledge) []knowledge.DecipherConceptRelation {
	if kb == nil {
		return nil
	}

	relations := kb.GetConceptRelationsAsDecipher(concept)
	result := make([]knowledge.DecipherConceptRelation, 0)

	for _, rel := range relations {
		if rel.Weight >= minWeight {
			result = append(result, rel)
		}
	}

	return result
}

// ExpandConcepts expands multiple concepts using the knowledge graph.
// This replaces the hardcoded ConceptGraph expansion.
func ExpandConcepts(concepts []string, minWeight float64, kb *knowledge.Knowledge) map[string][]knowledge.DecipherConceptRelation {
	if kb == nil {
		return nil
	}

	return kb.ExpandConcepts(concepts, minWeight)
}

// getRelationPath checks if there's a relation between two concepts using the knowledge base.
func getRelationPath(from, to string, kb *knowledge.Knowledge) bool {
	if kb == nil {
		return false
	}

	relations := kb.GetConceptRelations(from)
	for _, rel := range relations {
		if rel.To == to || rel.From == to {
			return true
		}
	}

	// Also check reverse
	relations = kb.GetConceptRelations(to)
	for _, rel := range relations {
		if rel.To == from || rel.From == from {
			return true
		}
	}

	return false
}

// getRelationPathsForConcepts returns relation paths between multiple concepts.
func getRelationPathsForConcepts(concepts []string, kb *knowledge.Knowledge) []string {
	if kb == nil {
		return nil
	}

	paths := make([]string, 0)
	seen := make(map[string]bool)

	// Check each pair
	for i := 0; i < len(concepts); i++ {
		for j := i + 1; j < len(concepts); j++ {
			conceptA := concepts[i]
			conceptB := concepts[j]

			rels := kb.GetConceptRelations(conceptA)
			for _, rel := range rels {
				path := ""
				if rel.To == conceptB {
					path = conceptA + " --[" + rel.Type + "]--> " + conceptB
				} else if rel.From == conceptB {
					path = conceptB + " --[" + rel.Type + "]--> " + conceptA
				}

				if path != "" && !seen[path] {
					paths = append(paths, path)
					seen[path] = true
				}
			}
		}
	}

	return paths
}

// =============================================================================
// Phase E: Passage-Level Convergence via Generic Concept Activation
// =============================================================================

// PassageSignal represents a signal extracted from passage-level analysis.
type PassageSignal struct {
	Token      string
	Concept    string
	Weight     float64
	Confidence string
	MatchForm  string
	MatchScore float64
	// IsStandaloneToken is true if this signal comes from a standalone
	// preposition/function word. Used to suppress structural noise while
	// preserving meaningful direct matches from real words.
	IsStandaloneToken bool
}

// EvidenceID returns a deterministic identity for this signal.
// Used for deduplication before scoring.
func (ps PassageSignal) EvidenceID() string {
	// Include token, concept, and match form for meaningful identity
	return ps.Token + "|" + ps.Concept + "|" + ps.MatchForm
}

// ActivatedConcept represents a concept activated through form matching.
type ActivatedConcept struct {
	Concept           string
	Strength         float64
	Sources          []string // tokens that activated this concept
	Confidence       string
	IsDirectEvidence bool     // true if activation has direct (not graph-propagated) evidence
}

// AnalyzePassageTokens analyzes tokens for passage-level signals.
// Uses direct form matches from knowledge when available.
// Fuzzy matching is only used as a fallback for tokens without curated forms.
// Returns activated concepts - NO semantic bucket booleans.
func AnalyzePassageTokens(tokens []string, kb *knowledge.Knowledge) []PassageSignal {
	signals := make([]PassageSignal, 0)

	// Get all known anchors for fuzzy fallback
	anchors := GetAllAnchors(kb)

	for _, token := range tokens {
		// Track if this token is a standalone preposition/function word.
		isStandalone := IsPrepositionOrFunctionWord(token)

		// First: check for direct form match in knowledge base
		// If a token has ANY curated form, use it and skip fuzzy matching
		// to prevent false positives from similar but unrelated words.
		// Standalone tokens still get processed but flagged for downstream suppression.
		matchedToken := ""
		directForms := kb.GetFormsByText(token)
		for _, form := range directForms {
			// Use any curated form match (confidence affects signal level, not existence)
			// Skip speculative forms to avoid noise
			if form.Confidence != ConfidenceSpeculative {
				signals = append(signals, PassageSignal{
					Token:              token,
					Concept:            form.Concept,
					Weight:             float64(form.Weight),
					Confidence:         form.Confidence,
					MatchForm:          token,
					MatchScore:         1.0,
					IsStandaloneToken: isStandalone,
				})
				matchedToken = token
				break
			}
		}
		// If we found a direct form, skip fuzzy matching for this token
		if matchedToken == token {
			continue
		}

		// No direct match found - try fuzzy anchor matching
		candidates, _ := GenerateCandidateForms(token)
		matches, _ := FuzzyMatchEvidence(candidates, anchors)
		for _, match := range matches {
			if match.Distance < 1.0 { // Only strong matches
				// Find the concept from the matched anchor
				for _, anchor := range anchors {
					if anchor.Form == match.AnchorForm {
						signals = append(signals, PassageSignal{
							Token:              token,
							Concept:            anchor.Concept,
							Weight:             match.Weight,
							Confidence:         computeSignalConfidence(anchor.Confidence, match.Weight),
							MatchForm:          match.AnchorForm,
							MatchScore:         match.Weight,
							IsStandaloneToken: isStandalone,
						})
						break
					}
				}
			}
		}
	}

	return signals
}

// computeSignalConfidence determines confidence from anchor confidence and match weight.
// Lower weights reduce confidence unless anchor is verified.
func computeSignalConfidence(anchorConf string, weight float64) string {
	if anchorConf == ConfidenceVerified {
		return ConfidenceVerified
	}
	// Reduce confidence based on weight
	if weight > 0.7 {
		return ConfidencePlausible
	}
	return ConfidenceSpeculative
}

// ComputeActivatedConcepts converts passage signals to activated concepts.
// This groups signals by concept and sums activation strength.
// IMPORTANT: Deduplicates by evidence ID to prevent duplicate evidence inflation.
func ComputeActivatedConcepts(signals []PassageSignal, directConcepts []string, kb *knowledge.Knowledge) []ActivatedConcept {
	// Deduplicate signals first to prevent duplicate evidence inflation
	seenIDs := make(map[string]bool)
	dedupedSignals := make([]PassageSignal, 0, len(signals))
	for _, sig := range signals {
		id := sig.EvidenceID()
		if !seenIDs[id] {
			seenIDs[id] = true
			dedupedSignals = append(dedupedSignals, sig)
		}
	}

	activated := make(map[string]*ActivatedConcept)

	// Process fuzzy match signals (deduplicated)
	for _, sig := range dedupedSignals {
		if existing, ok := activated[sig.Concept]; ok {
			existing.Strength += sig.Weight
			existing.Sources = append(existing.Sources, sig.Token)
			// Upgrade confidence if higher
			if sig.Confidence == ConfidenceVerified {
				existing.Confidence = ConfidenceVerified
			}
		} else {
			activated[sig.Concept] = &ActivatedConcept{
				Concept:           sig.Concept,
				Strength:         sig.Weight,
				Sources:          []string{sig.Token},
				Confidence:       sig.Confidence,
				IsDirectEvidence: true, // passage signals have direct evidence
			}
		}
	}

	// Process direct concepts from channels
	for _, c := range directConcepts {
		if existing, ok := activated[c]; ok {
			existing.Strength += 0.5 // Direct concepts add baseline activation
		} else {
			activated[c] = &ActivatedConcept{
				Concept:    c,
				Strength:   0.5,
				Sources:    []string{"direct"},
				Confidence: ConfidencePlausible,
			}
		}
	}

	// Expand concepts through graph relations
	conceptList := make([]string, 0, len(activated))
	for c := range activated {
		conceptList = append(conceptList, c)
	}

	// Get expansions and add related concepts
	expansions := ExpandConcepts(conceptList, 0.3, kb)
	for _, relatedRels := range expansions {
		for _, rel := range relatedRels {
			if existing, ok := activated[rel.To]; ok {
				existing.Strength += rel.Weight * 0.5 // Related concepts get partial activation
			} else {
				activated[rel.To] = &ActivatedConcept{
					Concept:    rel.To,
					Strength:   rel.Weight * 0.5,
					Sources:    []string{"graph_expansion"},
					Confidence: ConfidencePlausible,
				IsDirectEvidence: false, // graph expansion is indirect
				}
			}
		}
	}

	// Convert to slice
	result := make([]ActivatedConcept, 0, len(activated))
	for _, ac := range activated {
		result = append(result, *ac)
	}

	return result
}

// ComputeCoActivation finds concepts that co-activate through relation paths.
// This is the generic way to detect convergence - through shared relation structure.
func ComputeCoActivation(activated []ActivatedConcept, kb *knowledge.Knowledge) float64 {
	if len(activated) < 2 {
		return 0.0
	}

	// Count concepts that share relation edges
	conceptSet := make(map[string]bool)
	for _, ac := range activated {
		conceptSet[ac.Concept] = true
	}

	var sharedRelations float64
	var totalPairs float64

	// Check each pair of activated concepts for shared relations
	for i := 0; i < len(activated); i++ {
		for j := i + 1; j < len(activated); j++ {
			totalPairs++

			// Check if there's a relation path between these concepts
			conceptA := activated[i].Concept
			conceptB := activated[j].Concept

			// Look for direct or indirect relations using knowledge base
			if getRelationPath(conceptA, conceptB, kb) {
				sharedRelations++
			}
		}
	}

	if totalPairs == 0 {
		return 0.0
	}

	return sharedRelations / totalPairs
}
