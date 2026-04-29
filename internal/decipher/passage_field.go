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
	TokenSources  []string
	Depth         int
	RelationPaths []string
}

// PassageFields is a collection of PassageField instances.
type PassageFields []*PassageField

// Merge merges another PassageFields collection into this one.
func (pf PassageFields) Merge(other PassageFields) PassageFields {
	fieldMap := make(map[string]*PassageField)

	for _, field := range pf {
		fieldMap[field.Concept] = field
	}

	for _, otherField := range other {
		if existing, exists := fieldMap[otherField.Concept]; exists {
			existing.Strength += otherField.Strength
			if otherField.Strength > existing.Strength {
				existing.Confidence = otherField.Confidence
			}
			existing.TokenSources = mergeStringSlices(existing.TokenSources, otherField.TokenSources)
			existing.RelationPaths = mergeStringSlices(existing.RelationPaths, otherField.RelationPaths)
		} else {
			fieldMap[otherField.Concept] = otherField
		}
	}

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
func BuildPassageFieldsFromGraph(graph *ActivationGraph) PassageFields {
	fields := make(PassageFields, 0, len(graph.Nodes))

	for _, node := range graph.Nodes {
		field := &PassageField{
			Concept:       node.Concept,
			Strength:      node.Strength,
			Confidence:    node.Confidence,
			TokenSources:  []string{},
			Depth:         node.Depth,
			RelationPaths: []string{},
		}

		tokenSet := make(map[string]bool)
		for _, evidence := range node.Evidence {
			if evidence.SourceToken != "" {
				tokenSet[evidence.SourceToken] = true
			}
		}
		for token := range tokenSet {
			field.TokenSources = append(field.TokenSources, token)
		}

		fields = append(fields, field)
	}

	return fields
}

// AnalyzePassage analyzes a multi-token passage through existing channels.
func AnalyzePassage(passage string, engine *Engine) (PassageFields, error) {
	tokens := tokenizePassage(passage)
	if len(tokens) == 0 {
		return nil, nil
	}

	var allChannels []ChannelResult
	var allPassageSignals []PassageSignal
	conceptExpansions := make(map[string][]knowledge.DecipherConceptRelation)

	for _, token := range tokens {
		reading := engine.Analyze(token)

		allChannels = append(allChannels, reading.Channels...)

		for _, sig := range reading.Convergence.ActivatedConcepts {
			for _, source := range sig.Sources {
				allPassageSignals = append(allPassageSignals, PassageSignal{
					Token:      source,
					Concept:    sig.Concept,
					Weight:     sig.Strength,
					Confidence: sig.Confidence,
					MatchForm:  token,
					MatchScore: 1.0,
				})
			}
		}

		for concept, expansions := range reading.ConceptExpansions {
			conceptExpansions[concept] = expansions
		}
	}

	graph := BuildGraphFromEvidence(allChannels, allPassageSignals, nil, conceptExpansions, engine.Knowledge)
	graph.PropagateActivation()

	passageFields := BuildPassageFieldsFromGraph(graph)

	return passageFields, nil
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

// containsString checks if a slice contains a string.
func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
