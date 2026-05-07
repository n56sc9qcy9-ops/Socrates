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
	Depth         int            // 0 = direct, 1+ = graph-expanded
	TokenSources  []string       // Original tokens that activated this field
	EvidencePaths []EvidencePath // Evidence paths explaining this field
	RelationPaths []string       // Relation paths through the activation graph
	EvidenceCount int            // Number of distinct evidence paths supporting this field
	// IsDirectEvidence is true if this field has any direct form/glyph/script evidence.
	// False means evidence comes only from symbolic neighbor expansion or other indirect sources.
	// Used by counsellor to distinguish direct evidence from propagated suggestions.
	IsDirectEvidence bool
}

// Gate returns true if this passage field passes a harmonic field gate.
// Requires minimum strength and minimum evidence count.
// Very strong verified concepts (strength >= 0.5) bypass evidence count requirements.
func (f *PassageField) Gate(gate HarmonicFieldGate) bool {
	// Check strength threshold
	if f.Strength < gate.MinConceptStrength {
		return false
	}

	// For very strong verified concepts, allow even with zero/low evidence count
	if f.Strength >= 0.5 && f.Confidence == ConfidenceVerified {
		return true
	}

	// Check evidence count
	if f.EvidenceCount < gate.MinEvidenceCount {
		return false
	}

	// Check confidence - verified and plausible pass, speculative fails
	if f.Confidence == ConfidenceVerified || f.Confidence == ConfidencePlausible {
		return true
	}
	return false
}

// GateStrength returns the gated strength and whether this field passes the gate.
func (f *PassageField) GateStrength(gate HarmonicFieldGate) (float64, bool) {
	if !f.Gate(gate) {
		return 0, false
	}

	// Apply speculative discount
	if f.Confidence == ConfidenceSpeculative {
		return f.Strength * gate.SpeculativeDiscount, true
	}

	return f.Strength, true
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
			// Merge: combine strength, keep higher confidence
			existing.Strength += otherField.Strength
			existing.Confidence = keepHigherConfidence(existing.Confidence, otherField.Confidence)
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

// TopFields returns the top N fields, prioritizing direct concept evidence.
// Fields with direct evidence (exact matches, curated fragments) rank above
// propagated/graph-expanded fields, regardless of accumulated strength.
func (pf PassageFields) TopFields(n int) PassageFields {
	if len(pf) <= n {
		return pf
	}

	sorted := make(PassageFields, len(pf))
	copy(sorted, pf)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			// Primary sort: direct evidence fields rank above propagated
			if !sorted[i].IsDirectEvidence && sorted[j].IsDirectEvidence {
				sorted[i], sorted[j] = sorted[j], sorted[i]
				continue
			}
			if sorted[i].IsDirectEvidence == sorted[j].IsDirectEvidence {
				// Secondary: within same evidence type, sort by strength (descending)
				if sorted[j].Strength > sorted[i].Strength {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
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
		// Count unique evidence paths by source type for gating
		evidCount := len(node.Evidence)

		// Check if any evidence is direct
		isDirectEvidence := false
		for _, ev := range node.Evidence {
			if ev.IsDirect {
				isDirectEvidence = true
				break
			}
		}

		field := &PassageField{
			Concept:          node.Concept,
			Strength:         node.Strength,
			Confidence:       node.Confidence,
			Depth:            node.Depth,
			TokenSources:     []string{},
			EvidencePaths:    make([]EvidencePath, 0),
			RelationPaths:    []string{},
			EvidenceCount:    evidCount,
			IsDirectEvidence: isDirectEvidence,
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

		// Analyze passage tokens for this single token.
		// For non-Latin scripts, analyze the whole string for whole-word lookup.
		var passageTokens []string
		if forms.Script != ScriptLatin {
			passageTokens = []string{forms.Original}
		} else {
			passageTokens = forms.Tokens
		}
		signals := AnalyzePassageTokens(passageTokens, kb)

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

// confidencePriority returns numeric priority for confidence levels.
// Higher number = higher priority.
func confidencePriority(conf string) int {
	switch conf {
	case ConfidenceVerified:
		return 3
	case ConfidencePlausible:
		return 2
	case ConfidenceSpeculative:
		return 1
	default:
		return 0
	}
}

// keepHigherConfidence returns the higher of two confidences.
// This is the correct function for merging: keep the better confidence.
func keepHigherConfidence(a, b string) string {
	if confidencePriority(a) > confidencePriority(b) {
		return a
	}
	return b
}
