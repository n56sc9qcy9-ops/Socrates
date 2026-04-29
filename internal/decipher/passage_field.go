package decipher

import (
	"socrates/internal/knowledge"
)

// PassageField represents a concept field activated by passage tokens.
// Groups activation by concept with evidence from multiple token sources.
type PassageField struct {
	Concept       string
	Strength      float64
	Confidence    string
	Depth         int             // 0 = direct, 1+ = graph-expanded
	TokenSources  []string       // Original tokens that activated this field
	EvidencePaths []EvidencePath // Evidence paths explaining this field
	RelationPaths []string       // Relation paths through the activation graph
}

// PassageFields is a collection of PassageField instances.
type PassageFields []*PassageField

// Merge merges another PassageFields collection into this one.
// Duplicate fields (same concept) are merged, not duplicated.
// Evidence paths are deduplicated by EvidenceID().
func (pf PassageFields) Merge(other PassageFields) PassageFields {
	fieldMap := make(map[string]*PassageField)

	// Add existing fields
	for _, field := range pf {
		fieldMap[field.Concept] = field
	}

	// Merge other fields
	for _, otherField := range other {
		if existing, exists := fieldMap[otherField.Concept]; exists {
			// Merge: combine strength, keep highest confidence
			existing.Strength += otherField.Strength
			if isHigherConfidence(otherField.Confidence, existing.Confidence) {
				existing.Confidence = otherField.Confidence
			}
			// Merge token sources without duplicates
			existing.TokenSources = mergeStringSlices(existing.TokenSources, otherField.TokenSources)
			// Merge evidence paths without duplicates
			existing.EvidencePaths = mergeEvidencePaths(existing.EvidencePaths, otherField.EvidencePaths)
			// Merge relation paths without duplicates
			existing.RelationPaths = mergeStringSlices(existing.RelationPaths, otherField.RelationPaths)
		} else {
			fieldMap[otherField.Concept] = otherField
		}
	}

	// Convert map back to slice
	result := make(PassageFields, 0, len(fieldMap))
	for _, field := range fieldMap {
		result = append(result, field)
	}
	return result
}

// GetField returns a field by concept, or nil if not found.
func (pf PassageFields) GetField(concept string) *PassageField {
	for _, field := range pf {
		if field.Concept == concept {
			return field
		}
	}
	return nil
}

// TopFields returns the top N fields by strength.
func (pf PassageFields) TopFields(n int) PassageFields {
	if len(pf) <= n {
		return pf
	}

	sorted := make(PassageFields, len(pf))
	copy(sorted, pf)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Strength > sorted[i].Strength {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted[:n]
}

// BuildPassageFieldsFromGraph constructs PassageFields from an activation graph.
// Populates token sources and evidence paths from graph nodes.
func BuildPassageFieldsFromGraph(graph *ActivationGraph) PassageFields {
	fields := make(PassageFields, 0, len(graph.Nodes))

	for _, node := range graph.Nodes {
		field := &PassageField{
			Concept:       node.Concept,
			Strength:      node.Strength,
			Confidence:    node.Confidence,
			Depth:         node.Depth,
			TokenSources:  []string{},
			EvidencePaths: make([]EvidencePath, 0),
			RelationPaths: []string{},
		}

		// Collect token sources from evidence
		tokenSet := make(map[string]bool)
		for _, evidence := range node.Evidence {
			if evidence.SourceToken != "" {
				tokenSet[evidence.SourceToken] = true
			}
		}
		for token := range tokenSet {
			field.TokenSources = append(field.TokenSources, token)
		}

		// Copy evidence paths from graph node
		field.EvidencePaths = append(field.EvidencePaths, node.Evidence...)

		// Collect relation paths from edges
		for _, edge := range graph.Edges {
			if edge.From == node.Concept || edge.To == node.Concept {
				path := formatRelationPath(edge)
				if !containsString(field.RelationPaths, path) {
					field.RelationPaths = append(field.RelationPaths, path)
				}
			}
		}

		fields = append(fields, field)
	}

	return fields
}

// formatRelationPath formats an edge as a human-readable relation path.
func formatRelationPath(edge *ActivationEdge) string {
	return edge.From + " ->(" + edge.RelationType + ") " + edge.To
}

// AnalyzePassageFromTokens analyzes a list of tokens through the activation graph.
// This is the internal version that avoids circular calls.
// For external use, call Engine.Analyze with a multi-token passage instead.
func AnalyzePassageFromTokens(tokens []string, kb *knowledge.Knowledge) PassageFields {
	if len(tokens) == 0 {
		return nil
	}

	var allChannels []ChannelResult
	var allPassageSignals []PassageSignal
	conceptExpansions := make(map[string][]knowledge.DecipherConceptRelation)

	for _, token := range tokens {
		// Generate forms for single token
		forms := GenerateForms(token)

		// Run channels for this token
		channels := RunAllChannels(forms, kb)
		allChannels = append(allChannels, channels...)

		// Analyze passage tokens for this single token
		signals := AnalyzePassageTokens(forms.Tokens, kb)

		// Mark signals with original token as source
		for _, sig := range signals {
			allPassageSignals = append(allPassageSignals, PassageSignal{
				Token:      token, // Original token from passage
				Concept:    sig.Concept,
				Weight:     sig.Weight,
				Confidence: sig.Confidence,
				MatchForm:  token,
				MatchScore: 1.0,
			})
		}

		// Expand concepts for this token
		directConcepts := extractDirectConcepts(channels)
		expansions := ExpandConcepts(directConcepts, 0.4, kb)
		for concept, exps := range expansions {
			conceptExpansions[concept] = exps
		}
	}

	// Build activation graph from collected data
	graph := BuildGraphFromEvidence(allChannels, allPassageSignals, nil, conceptExpansions, kb)
	graph.PropagateActivation()

	// Build passage fields from graph
	passageFields := BuildPassageFieldsFromGraph(graph)

	return passageFields
}

// AnalyzePassage analyzes a multi-token passage through existing channels.
// WARNING: This function calls engine.Analyze internally which may cause issues
// with multi-token inputs. For internal use, prefer AnalyzePassageFromTokens.
func AnalyzePassage(passage string, engine *Engine) (PassageFields, error) {
	tokens := tokenizePassage(passage)
	if len(tokens) == 0 {
		return nil, nil
	}

	// Use the internal function to avoid circular recursion
	return AnalyzePassageFromTokens(tokens, engine.Knowledge), nil
}

// tokenizePassage splits a passage into tokens.
func tokenizePassage(passage string) []string {
	tokens := make([]string, 0)
	current := ""
	for _, r := range passage {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		tokens = append(tokens, current)
	}
	return tokens
}

// mergeStringSlices merges two string slices without duplicates.
func mergeStringSlices(a, b []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(a)+len(b))
	for _, s := range a {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	for _, s := range b {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

// mergeEvidencePaths merges evidence paths without duplicates.
func mergeEvidencePaths(a, b []EvidencePath) []EvidencePath {
	seen := make(map[string]bool)
	result := make([]EvidencePath, 0, len(a)+len(b))
	for _, e := range a {
		id := e.EvidenceID()
		if !seen[id] {
			seen[id] = true
			result = append(result, e)
		}
	}
	for _, e := range b {
		id := e.EvidenceID()
		if !seen[id] {
			seen[id] = true
			result = append(result, e)
		}
	}
	return result
}

// containsString checks if a slice contains a string.
func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// isHigherConfidence returns true if b is higher priority than a.
func isHigherConfidence(a, b string) bool {
	if b == ConfidenceVerified && a != ConfidenceVerified {
		return true
	}
	if b == ConfidencePlausible && a == ConfidenceSpeculative {
		return true
	}
	return false
}
