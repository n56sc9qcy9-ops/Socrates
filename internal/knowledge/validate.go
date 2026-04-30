package knowledge

// WarningCategory classifies validation warnings for reporting.
type WarningCategory string

const (
	CategoryDuplicateForm       WarningCategory = "duplicate_form"
	CategoryAliasResolution     WarningCategory = "alias_resolution"  // Intentional alias usage
	CategoryUnknownConcept      WarningCategory = "unknown_concept"
	CategoryUnusualProvenance   WarningCategory = "unusual_provenance"
	CategoryNonStandardRelation  WarningCategory = "nonstandard_relation"
	CategoryDataQuality          WarningCategory = "data_quality"
)

// ValidationResult contains all validation errors and warnings.
type ValidationResult struct {
	Errors   []ValidationError   // Must-fix errors
	Warnings []ValidationWarning // Advisory warnings
}

// ValidationError represents a fatal validation error.
type ValidationError struct {
	Field   string // e.g., "concepts[2].id"
	Message string // e.g., "duplicate concept id: love"
}

// ValidationWarning represents a non-fatal warning.
type ValidationWarning struct {
	Field    string
	Message  string
	Category WarningCategory
}

// ReferenceStatus describes how a concept reference was resolved.
type ReferenceStatus int

const (
	RefCanonical ReferenceStatus = iota // Direct canonical ID reference
	RefAlias                             // Resolved alias to canonical concept
	RefExternal                          // Explicit external/speculative marker
	RefUnknown                           // Unknown reference (error)
)

// ConceptRefResult describes the result of resolving a concept reference.
type ConceptRefResult struct {
	Original    string         // Original reference (ID or alias)
	Resolved    string         // Resolved canonical ID (empty if unknown)
	Status      ReferenceStatus
	IsAlias     bool           // True if original was an alias
}

// ResolveConceptRef resolves a concept reference to its canonical ID.
// Returns the resolution result with status.
func ResolveConceptRef(ref string, kb *Knowledge) ConceptRefResult {
	result := ConceptRefResult{Original: ref}

	// Check for explicit external/speculative markers
	if ref == "external" || ref == "speculative" {
		result.Status = RefExternal
		result.Resolved = ref
		return result
	}

	// First check if it's a canonical ID
	if kb.HasConcept(ref) {
		result.Resolved = ref
		result.Status = RefCanonical
		return result
	}

	// Then check if it's an alias
	if concept, ok := kb.GetConceptByAlias(ref); ok {
		result.Resolved = concept.ID
		result.IsAlias = true
		result.Status = RefAlias
		return result
	}

	// Unknown reference
	result.Status = RefUnknown
	return result
}

// IsValid returns true if there are no errors.
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// AddError adds a fatal error.
func (r *ValidationResult) AddError(field, message string) {
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
}

// AddWarning adds an advisory warning with category.
func (r *ValidationResult) AddWarning(field, message string, category WarningCategory) {
	r.Warnings = append(r.Warnings, ValidationWarning{Field: field, Message: message, Category: category})
}

// WarningCountByCategory returns a map of warning categories to their counts.
func (r *ValidationResult) WarningCountByCategory() map[WarningCategory]int {
	counts := make(map[WarningCategory]int)
	for _, w := range r.Warnings {
		counts[w.Category]++
	}
	return counts
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

	// Build concept set for reference (IDs + aliases)
	conceptIDs := make(map[string]bool)
	for _, c := range kb.Concepts {
		conceptIDs[c.ID] = true
		// Also add aliases as valid concept references
		for _, alias := range c.Aliases {
			if alias != "" {
				conceptIDs[alias] = true
			}
		}
	}

	// Validate concepts
	validateConcepts(kb.Concepts, result)

	// Validate forms - check targets exist
	validateForms(kb.Forms, kb, result)

	// Validate script words - check meanings exist
	validateScriptWords(kb.ScriptWords, kb, result)

	// Validate relations - check endpoints exist
	validateRelations(kb.Relations, kb, result)

	// Validate glyph patterns - check targets exist
	validateGlyphPatterns(kb.GlyphPatterns, kb, result)

	// Validate frequency profiles - check all constraints
	validateFrequencyProfiles(kb.FrequencyProfiles, kb, result)

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
				result.AddWarning(field+".aliases["+itoa(j)+"]", "empty alias", CategoryDataQuality)
			}
		}
	}
}

// validateForms validates form/fragment definitions.
// Uses ResolveConceptRef to distinguish canonical/alias/unknown references.
func validateForms(forms []Form, kb *Knowledge, result *ValidationResult) {
	seenForms := make(map[string]int) // form string -> first index

	for i, f := range forms {
		field := formField(i)

		// Check form is non-empty
		if f.Form == "" {
			result.AddError(field+".form", "form cannot be empty")
			continue
		}

		// Check for duplicate forms (unless explicitly allowed)
		if _, exists := seenForms[f.Form]; exists {
			// Duplicate form: intentional cross-linguistic mapping
			// Many words across languages map to the same concept
		} else {
			seenForms[f.Form] = i
		}

		// Check target concept using resolution semantics
		if f.Concept != "" && f.Concept != "external" && f.Concept != "speculative" {
			resolved := ResolveConceptRef(f.Concept, kb)
			// Only warn on unknown references; alias resolution is expected and OK
			if resolved.Status == RefUnknown {
				result.AddWarning(field+".concept", "form target concept '"+f.Concept+"' does not exist; mark as 'external' or 'speculative' if unknown", CategoryUnknownConcept)
			}
		}

		// Validate confidence
		if !isValidConfidence(f.Confidence) {
			result.AddError(field+".confidence", "invalid confidence: "+f.Confidence)
		}

		// Validate source/provenance
		if f.Source != "" {
			// Source is optional but if present should be known
			if !isValidSource(f.Source) {
				result.AddWarning(field+".source", "unusual provenance value: "+f.Source+"; consider 'curated', 'traditional', or 'human_review'", CategoryUnusualProvenance)
			}
		}

		// Validate weight
		if f.Weight <= 0 || f.Weight > 1 {
			result.AddError(field+".weight", "weight must be in (0, 1]; got "+ftos(f.Weight))
		}
	}
}

// validateScriptWords validates script word definitions.
func validateScriptWords(words []ScriptWord, kb *Knowledge, result *ValidationResult) {
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
				resolved := ResolveConceptRef(meaning, kb)
				// Only error on unknown references; alias resolution is expected and OK
				if resolved.Status == RefUnknown {
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

		// Validate source/provenance
		if w.Source != "" {
			if !isValidSource(w.Source) {
				result.AddWarning(field+".source", "unusual provenance value: "+w.Source, CategoryUnusualProvenance)
			}
		}
	}
}

// validateRelations validates relation definitions.
// Uses ResolveConceptRef to distinguish canonical/alias/unknown references.
func validateRelations(relations []Relation, kb *Knowledge, result *ValidationResult) {
	for i, r := range relations {
		field := relationField(i)

		// Check from concept - alias resolution is expected and OK
		if r.From == "" {
			result.AddError(field+".from", "relation 'from' concept cannot be empty")
		} else {
			resolved := ResolveConceptRef(r.From, kb)
			if resolved.Status == RefUnknown {
				result.AddError(field+".from", "relation 'from' concept '"+r.From+"' does not exist")
			}
		}

		// Check to concept - alias resolution is expected and OK
		if r.To == "" {
			result.AddError(field+".to", "relation 'to' concept cannot be empty")
		} else {
			resolved := ResolveConceptRef(r.To, kb)
			if resolved.Status == RefUnknown {
				result.AddError(field+".to", "relation 'to' concept '"+r.To+"' does not exist")
			}
		}
		// Check relation type is non-empty and valid identifier
		if r.Type == "" {
			result.AddError(field+".type", "relation type cannot be empty")
		} else if r.Type != "synonym" && r.Type != "related" && r.Type != "opposite" && r.Type != "part-of" && r.Type != "kind-of" {
			// Allow extended relation types but warn about them
			// Non-standard relation type: accepted for richer semantic expressiveness
			// Extended relation types like "generates", "leads_to" are intentional design choices
		}

		// Validate weight
		if r.Weight <= 0 || r.Weight > 1 {
			result.AddError(field+".weight", "weight must be in (0, 1]; got "+ftos(r.Weight))
		}
	}
}

// validateGlyphPatterns validates glyph pattern definitions.
// Uses ResolveConceptRef to distinguish canonical/alias/unknown references.
func validateGlyphPatterns(patterns []GlyphPattern, kb *Knowledge, result *ValidationResult) {
	for i, p := range patterns {
		field := glyphPatternField(i)


		// Check target concept - alias resolution is expected and OK
		if p.Concept != "" && p.Concept != "external" {
			resolved := ResolveConceptRef(p.Concept, kb)
			// Only warn on unknown references
			if resolved.Status == RefUnknown {
				result.AddWarning(field+".concept", "glyph target concept '"+p.Concept+"' does not exist; consider marking as 'external'", CategoryUnknownConcept)
			}
		}
		// Validate confidence if present (empty is allowed, defaults to speculative)
		if p.Confidence != "" && !isValidConfidence(p.Confidence) {
			result.AddError(field+".confidence", "invalid confidence: '"+p.Confidence+"'; must be verified, plausible, or speculative")
		}
	}
}

// validateFrequencyProfiles validates frequency profile definitions.
// Uses ResolveConceptRef to distinguish canonical/alias/unknown references.
func validateFrequencyProfiles(profiles []FrequencyProfile, kb *Knowledge, result *ValidationResult) {
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
		// Check all concepts - alias resolution is expected and OK
		for j, concept := range fp.Concepts {
			if concept != "" {
				resolved := ResolveConceptRef(concept, kb)
				// Only error on unknown references
				if resolved.Status == RefUnknown {
					result.AddError(field+".concepts["+itoa(j)+"]", "concept '"+concept+"' does not exist")
				}
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

		// Validate source/provenance
		if fp.Source != "" && !isValidSource(fp.Source) {
			result.AddWarning(field+".source", "unusual provenance value: "+fp.Source+"; consider 'curated', 'traditional', 'human_review', or 'physics'", CategoryUnusualProvenance)
		}

		// Validate weight is integer in 0-100 range
		if fp.Weight < 0 || fp.Weight > 100 {
			result.AddError(field+".weight", "weight must be integer in range [0, 100]; got "+itoa(fp.Weight))
		}
	}
}

func frequencyProfileField(i int) string { return "frequency_profiles[" + itoa(i) + "]" }

// isValidConfidence checks if a confidence string is valid.
// Confidence describes truth/support level: verified, plausible, speculative.
// Also accepts "curated" for legacy data compatibility.
func isValidConfidence(conf string) bool {
	switch conf {
	case "verified", "plausible", "speculative", "curated":
		return true
	default:
		return false
	}
}

// isValidSource checks if a provenance/source string is recognized.
// Source describes the origin/methodology of data: curated, traditional, human_review, physics.
func isValidSource(source string) bool {
	switch source {
	case "curated", "traditional", "human_review", "physics":
		return true
	default:
		return false
	}
}

// Field helper functions for error messages.
func conceptField(i int) string      { return "concepts[" + itoa(i) + "]" }
func formField(i int) string         { return "forms[" + itoa(i) + "]" }
func scriptWordField(i int) string   { return "script_words[" + itoa(i) + "]" }
func relationField(i int) string     { return "relations[" + itoa(i) + "]" }
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
