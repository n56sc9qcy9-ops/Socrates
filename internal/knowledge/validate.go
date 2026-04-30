package knowledge

// ValidationResult contains all validation errors and warnings.
type ValidationResult struct {
	Errors   []ValidationError   // Must-fix errors
	Warnings []ValidationWarning  // Advisory warnings
}

// ValidationError represents a fatal validation error.
type ValidationError struct {
	Field   string // e.g., "concepts[2].id"
	Message string // e.g., "duplicate concept id: love"
}

// ValidationWarning represents a non-fatal warning.
type ValidationWarning struct {
	Field   string
	Message string
}

// IsValid returns true if there are no errors.
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// AddError adds a fatal error.
func (r *ValidationResult) AddError(field, message string) {
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
}

// AddWarning adds an advisory warning.
func (r *ValidationResult) AddWarning(field, message string) {
	r.Warnings = append(r.Warnings, ValidationWarning{Field: field, Message: message})
}

// ValidArchetypeIDs contains all known archetype IDs for validation.
var ValidArchetypeIDs = map[string]bool{
	"pythagorean_triple_1": true,
	"pythagorean_triple_2": true,
	"pythagorean_triple_3": true,
	"pythagorean_triple_4": true,
	"pythagorean_triple_5": true,
	"pythagorean_triple_6": true,
	"pythagorean_triple_7": true,
	"phi_approximant_1":    true,
	"phi_approximant_2":    true,
	"phi_approximant_3":    true,
	"metatron_node_1":      true,
	"metatron_node_2":      true,
	"metatron_node_3":      true,
}

// ValidateKnowledge runs all validation checks on the knowledge base.
// This is the main entry point for knowledge validation.
func ValidateKnowledge(kb *Knowledge) *ValidationResult {
	result := &ValidationResult{
		Errors:   make([]ValidationError, 0),
		Warnings: make([]ValidationWarning, 0),
	}

	// Build concept set for reference
	conceptIDs := make(map[string]bool)
	for _, c := range kb.Concepts {
		conceptIDs[c.ID] = true
	}

	// Validate concepts
	validateConcepts(kb.Concepts, result)

	// Validate forms - check targets exist
	validateForms(kb.Forms, conceptIDs, result)

	// Validate script words - check meanings exist
	validateScriptWords(kb.ScriptWords, conceptIDs, result)

	// Validate relations - check endpoints exist
	validateRelations(kb.Relations, conceptIDs, result)

	// Validate glyph patterns - check targets exist
	validateGlyphPatterns(kb.GlyphPatterns, conceptIDs, result)

	// Validate frequency profiles - check all constraints
	validateFrequencyProfiles(kb.FrequencyProfiles, conceptIDs, result)

	return result
}

// validateConcepts validates concept definitions.
func validateConcepts(concepts []Concept, result *ValidationResult) {
	seenIDs := make(map[string]int) // id -> first index

	for i, c := range concepts {
		field := conceptField(i)

		// Check ID is non-empty
		if c.ID == "" {
			result.AddError(field+".id", "concept id cannot be empty")
			continue
		}

		// Check for duplicate IDs
		if firstIdx, exists := seenIDs[c.ID]; exists {
			result.AddError(field+".id", "duplicate concept id: "+c.ID+" (first seen at "+conceptField(firstIdx)+")")
		} else {
			seenIDs[c.ID] = i
		}

		// Check name is non-empty
		if c.Name == "" {
			result.AddError(field+".name", "concept name cannot be empty")
		}

		// Check aliases don't create ambiguity
		for j, alias := range c.Aliases {
			if alias == "" {
				result.AddWarning(field+".aliases["+itoa(j)+"]", "empty alias")
			}
		}
	}
}

// validateForms validates form/fragment definitions.
func validateForms(forms []Form, conceptIDs map[string]bool, result *ValidationResult) {
	seenForms := make(map[string]int) // form string -> first index

	for i, f := range forms {
		field := formField(i)

		// Check form is non-empty
		if f.Form == "" {
			result.AddError(field+".form", "form cannot be empty")
			continue
		}

		// Check for duplicate forms (unless explicitly allowed)
		if firstIdx, exists := seenForms[f.Form]; exists {
			result.AddWarning(field+".form", "duplicate form: "+f.Form+" (first at "+formField(firstIdx)+")")
		} else {
			seenForms[f.Form] = i
		}

		// Check target concept exists or is marked external/speculative
		if f.Concept != "" && f.Concept != "external" && f.Concept != "speculative" {
			if !conceptIDs[f.Concept] {
				result.AddWarning(field+".concept", "form target concept '"+f.Concept+"' does not exist; mark as 'external' or 'speculative' if unknown")
			}
		}

		// Validate confidence
		if !isValidConfidence(f.Confidence) {
			result.AddError(field+".confidence", "invalid confidence: "+f.Confidence)
		}

		// Validate weight
		if f.Weight <= 0 || f.Weight > 1 {
			result.AddError(field+".weight", "weight must be in (0, 1]; got "+ftos(f.Weight))
		}
	}
}

// validateScriptWords validates script word definitions.
func validateScriptWords(words []ScriptWord, conceptIDs map[string]bool, result *ValidationResult) {
	for i, w := range words {
		field := scriptWordField(i)

		// Check script and word are non-empty
		if w.Script == "" {
			result.AddError(field+".script", "script cannot be empty")
		}
		if w.Word == "" {
			result.AddError(field+".word", "word cannot be empty")
		}

		// Check meanings reference existing concepts
		for j, meaning := range w.Meanings {
			if meaning != "" && meaning != "external" {
				if !conceptIDs[meaning] {
					result.AddError(field+".meanings["+itoa(j)+"]", "meaning '"+meaning+"' does not exist as a concept")
				}
			}
		}

		// Validate confidence
		if !isValidConfidence(w.Confidence) {
			result.AddError(field+".confidence", "invalid confidence: "+w.Confidence)
		}

		// Validate weight
		if w.Weight <= 0 || w.Weight > 1 {
			result.AddError(field+".weight", "weight must be in (0, 1]; got "+ftos(w.Weight))
		}
	}
}

// validateRelations validates relation definitions.
func validateRelations(relations []Relation, conceptIDs map[string]bool, result *ValidationResult) {
	for i, r := range relations {
		field := relationField(i)

		// Check from concept exists
		if r.From == "" {
			result.AddError(field+".from", "relation 'from' concept cannot be empty")
		} else if !conceptIDs[r.From] {
			result.AddError(field+".from", "relation 'from' concept '"+r.From+"' does not exist")
		}

		// Check to concept exists
		if r.To == "" {
			result.AddError(field+".to", "relation 'to' concept cannot be empty")
		} else if !conceptIDs[r.To] {
			result.AddError(field+".to", "relation 'to' concept '"+r.To+"' does not exist")
		}

		// Check relation type is non-empty and valid identifier
		if r.Type == "" {
			result.AddError(field+".type", "relation type cannot be empty")
		} else if r.Type != "synonym" && r.Type != "related" && r.Type != "opposite" && r.Type != "part-of" && r.Type != "kind-of" {
			// Allow extended relation types but warn about them
			result.AddWarning(field+".type", "non-standard relation type: "+r.Type+"; consider using synonym, related, opposite, part-of, or kind-of")
		}

		// Validate weight
		if r.Weight <= 0 || r.Weight > 1 {
			result.AddError(field+".weight", "weight must be in (0, 1]; got "+ftos(r.Weight))
		}
	}
}

// validateGlyphPatterns validates glyph pattern definitions.
func validateGlyphPatterns(patterns []GlyphPattern, conceptIDs map[string]bool, result *ValidationResult) {
	for i, p := range patterns {
		field := glyphPatternField(i)

		// Check target concept exists or is marked external
		if p.Concept != "" && p.Concept != "external" {
			if !conceptIDs[p.Concept] {
				result.AddWarning(field+".concept", "glyph target concept '"+p.Concept+"' does not exist; consider marking as 'external'")
			}
		}

		// Validate confidence if present (empty is allowed, defaults to speculative)
		if p.Confidence != "" && !isValidConfidence(p.Confidence) {
			result.AddError(field+".confidence", "invalid confidence: '"+p.Confidence+"'; must be verified, plausible, or speculative")
		}
	}
}

// validateFrequencyProfiles validates frequency profile definitions.
func validateFrequencyProfiles(profiles []FrequencyProfile, conceptIDs map[string]bool, result *ValidationResult) {
	seenIDs := make(map[string]int) // meaning_frequency_id -> first index

	for i, fp := range profiles {
		field := frequencyProfileField(i)

		// Check meaning_frequency_id is non-empty
		if fp.MeaningFrequencyID == "" {
			result.AddError(field+".meaning_frequency_id", "meaning_frequency_id cannot be empty")
			continue
		}

		// Check for duplicate IDs
		if firstIdx, exists := seenIDs[fp.MeaningFrequencyID]; exists {
			result.AddError(field+".meaning_frequency_id", "duplicate meaning_frequency_id: "+fp.MeaningFrequencyID+" (first at "+frequencyProfileField(firstIdx)+")")
		} else {
			seenIDs[fp.MeaningFrequencyID] = i
		}

		// Check concepts list is non-empty
		if len(fp.Concepts) == 0 {
			result.AddError(field+".concepts", "concepts list cannot be empty")
		}

		// Check all concepts exist
		for j, concept := range fp.Concepts {
			if concept != "" && !conceptIDs[concept] {
				result.AddError(field+".concepts["+itoa(j)+"]", "concept '"+concept+"' does not exist")
			}
		}

		// Check vector has exactly 3 integers
		if len(fp.Vector) != 3 {
			result.AddError(field+".vector", "vector must have exactly 3 integers; got "+itoa(len(fp.Vector)))
		}

		// Check ratio has exactly 2 integers (numerator, denominator)
		if len(fp.Ratio) != 2 {
			result.AddError(field+".ratio", "ratio must have exactly 2 integers [num, denom]; got "+itoa(len(fp.Ratio)))
		}

		// Check denominator is not zero
		if len(fp.Ratio) == 2 && fp.Ratio[1] == 0 {
			result.AddError(field+".ratio[1]", "ratio denominator cannot be zero")
		}

		// Check archetype is valid if present
		if fp.Archetype != "" && !ValidArchetypeIDs[fp.Archetype] {
			result.AddError(field+".archetype", "unknown archetype: "+fp.Archetype+"; valid IDs: pythagorean_triple_1-7, phi_approximant_1-3, metatron_node_1-3")
		}

		// Validate confidence
		if !isValidConfidence(fp.Confidence) {
			result.AddError(field+".confidence", "invalid confidence: "+fp.Confidence)
		}

		// Validate weight is integer in 0-100 range
		if fp.Weight < 0 || fp.Weight > 100 {
			result.AddError(field+".weight", "weight must be integer in range [0, 100]; got "+itoa(fp.Weight))
		}
	}
}

func frequencyProfileField(i int) string { return "frequency_profiles[" + itoa(i) + "]" }

// isValidConfidence checks if a confidence string is valid.
func isValidConfidence(conf string) bool {
	switch conf {
	case "verified", "plausible", "speculative":
		return true
	default:
		return false
	}
}

// Field helper functions for error messages.
func conceptField(i int) string   { return "concepts[" + itoa(i) + "]" }
func formField(i int) string      { return "forms[" + itoa(i) + "]" }
func scriptWordField(i int) string { return "script_words[" + itoa(i) + "]" }
func relationField(i int) string  { return "relations[" + itoa(i) + "]" }
func glyphPatternField(i int) string { return "glyph_patterns[" + itoa(i) + "]" }

// String conversion helpers.
func itoa(i int) string {
	if i < 10 {
		return string(rune('0' + i))
	}
	// Simple implementation for small numbers
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return s
}

func ftos(f float64) string {
	// Simple float-to-string for common cases
	if f == 0 {
		return "0"
	}
	if f == 1 {
		return "1"
	}
	// For other values, use a simple representation
	return string(rune('0'+int(f*10)%10)) + "." + string(rune('0'+int(f*100)%10))
}
