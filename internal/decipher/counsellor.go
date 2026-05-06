package decipher

import (
	"math"

	"socrates/internal/knowledge"
)

// counsellorEngine produces data-backed counsellor/transmutation fields.
// NOT fortune telling — curated relational guidance from loaded knowledge.
// Rules:
// - No hardcoded transmutation logic in Go
// - All mappings must be data-driven with source/lens/confidence
// - No fixed-future claims; use "if unchanged, tends toward..." language
// - No medical, legal, or clinical claims
// - Preserve human sovereignty; suggest, don't command

const (
	// counsellorMinStrength is the minimum activation strength for a concept
	// to be considered as a potential transmutation source.
	counsellorMinStrength = 0.15

	// counsellorMinCombinedStrength is the minimum combined strength (source × weight)
	// for a transmutation suggestion to appear in output.
	counsellorMinCombinedStrength = 0.10
)

// BuildCounsellorField builds a counsellor field from passage fields and knowledge.
// Returns nil when no data-backed transmutation path exists.
// Uses activated concepts/fields — does NOT inspect raw words directly for meaning.
func BuildCounsellorField(passageFields PassageFields, kb *knowledge.Knowledge) *CounsellorField {
	if kb == nil || len(kb.AllTransmutations()) == 0 {
		return nil
	}

	// Find activated concepts that match transmutation source fields
	sourceFields := findTransmutationSources(passageFields, kb)
	if len(sourceFields) == 0 {
		return nil
	}

	// Build suggestions from data-backed transmutation relations
	suggestions := buildSuggestions(sourceFields, kb)
	if len(suggestions) == 0 {
		return nil
	}

	// Build evidence paths
	evidencePaths := buildEvidencePaths(sourceFields, kb)

	return &CounsellorField{
		SourceFields:  sourceFields,
		Suggestions:   suggestions,
		EvidencePaths: evidencePaths,
	}
}

// findTransmutationSources finds activated passage fields that match transmutation source concepts.
func findTransmutationSources(passageFields PassageFields, kb *knowledge.Knowledge) []TransmutationSource {
	seen := make(map[string]bool)
	var sources []TransmutationSource

	// Get all transmutation source concepts from knowledge
	transmutationSources := make(map[string]bool)
	for _, t := range kb.AllTransmutations() {
		transmutationSources[t.From] = true
	}

	// Scan activated passage fields for matches
	for _, field := range passageFields {
		if field.Strength < counsellorMinStrength {
			continue
		}
		if !transmutationSources[field.Concept] {
			continue
		}
		if seen[field.Concept] {
			continue
		}
		seen[field.Concept] = true

		sources = append(sources, TransmutationSource{
			Concept:           field.Concept,
			Strength:          field.Strength,
			Depth:             field.Depth,
			EvidenceSources:   field.TokenSources,
			IsDirectEvidence:  field.IsDirectEvidence,
		})
	}

	return sources
}

// buildSuggestions creates transmutation suggestions from source fields and knowledge.
func buildSuggestions(sources []TransmutationSource, kb *knowledge.Knowledge) []TransmutationSuggestion {
	var suggestions []TransmutationSuggestion

	for _, src := range sources {
		transmutations := kb.GetTransmutationsFrom(src.Concept)
		for _, t := range transmutations {
			// Combined strength = source field strength × relation weight (normalized to 0-1)
			// Weight is stored as integer percentage (0-100), source strength is normalized (0-1)
			combinedStrength := src.Strength * t.Weight / 100.0
			if combinedStrength < counsellorMinCombinedStrength {
				continue
			}

			// Downgrade confidence for indirect evidence sources.
			// Direct form/glyph/script evidence uses normal confidence.
			// Symbolic neighbor expansion or graph-propagated sources get downgraded.
			suggestionConf := t.Confidence
			// Check both Depth (graph propagation) and IsDirectEvidence (form evidence)
			isIndirect := src.Depth > 0 || !src.IsDirectEvidence
			if isIndirect {
				if t.Confidence == ConfidenceVerified {
					suggestionConf = ConfidencePlausible
				} else if t.Confidence == ConfidencePlausible {
					suggestionConf = ConfidenceSpeculative
				}
				// ConfidenceSpeculative stays speculative
			}

			suggestions = append(suggestions, TransmutationSuggestion{
				SourceConcept:   src.Concept,
				TargetConcept:   t.To,
				Kind:            t.Kind,
				Confidence:      suggestionConf,
				Source:          t.Source,
				Lens:            t.Lens,
				Weight:          t.Weight,
				Strength:        math.Round(combinedStrength*100) / 100, // preserve precision
			})
		}
	}

	return suggestions
}

// buildEvidencePaths creates evidence traces from source fields to suggestions.
func buildEvidencePaths(sources []TransmutationSource, kb *knowledge.Knowledge) []TransmutationEvidence {
	var paths []TransmutationEvidence

	for _, src := range sources {
		transmutations := kb.GetTransmutationsFrom(src.Concept)
		for _, t := range transmutations {
			combinedStrength := src.Strength * t.Weight / 100.0
			if combinedStrength < counsellorMinCombinedStrength {
				continue
			}

			isDirect := src.Depth == 0
			paths = append(paths, TransmutationEvidence{
				SourceConcept:   src.Concept,
				TargetConcept:   t.To,
				Kind:            t.Kind,
				DataSourceFile:  "internal/knowledge/transmute.yaml",
				Notes:           t.Notes,
				IsDirectSource:  isDirect,
			})
		}
	}

	return paths
}
