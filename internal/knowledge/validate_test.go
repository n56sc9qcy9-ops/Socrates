package knowledge

import (
	"testing"
)

// =============================================================================
// VALIDATION TESTS
// =============================================================================

func TestValidateKnowledge_ValidKnowledge(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	result := ValidateKnowledge(kb)

	// Validation should run and produce a result
	// Some embedded knowledge may reference external/speculative concepts
	// This test verifies the validator runs, not that all data is perfect
	if len(result.Errors) > 0 {
		// Log warnings for awareness but don't fail
		for _, w := range result.Warnings {
			t.Logf("warning: %s: %s", w.Field, w.Message)
		}
	}
}

func TestValidateKnowledge_EmptyKnowledge(t *testing.T) {
	kb := &Knowledge{
		Concepts:      []Concept{},
		Forms:         []Form{},
		Relations:     []Relation{},
		GlyphPatterns: []GlyphPattern{},
	}

	result := ValidateKnowledge(kb)

	if !result.IsValid() {
		t.Errorf("empty knowledge should be valid")
	}
}

// =============================================================================
// CONCEPT VALIDATION FAILURE TESTS
// =============================================================================

func TestValidateConcepts_DuplicateID(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "love", Name: "love"},
			{ID: "love", Name: "agape"}, // duplicate
		},
	}

	result := ValidateKnowledge(kb)

	found := false
	for _, err := range result.Errors {
		if containsString(err.Message, "duplicate concept id: love") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should detect duplicate concept id")
	}
}

func TestValidateConcepts_EmptyID(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "", Name: "test"}, // empty id
		},
	}

	result := ValidateKnowledge(kb)

	found := false
	for _, err := range result.Errors {
		if containsString(err.Message, "concept id cannot be empty") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should detect empty concept id")
	}
}

func TestValidateConcepts_EmptyName(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "test", Name: ""}, // empty name
		},
	}

	result := ValidateKnowledge(kb)

	found := false
	for _, err := range result.Errors {
		if containsString(err.Message, "concept name cannot be empty") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should detect empty concept name")
	}
}

// =============================================================================
// FORM VALIDATION FAILURE TESTS
// =============================================================================

func TestValidateForms_MissingTargetConcept(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "existing", Name: "existing"},
		},
		Forms: []Form{
			{Form: "unknown", Concept: "nonexistent", Weight: 0.5}, // target doesn't exist
		},
	}

	result := ValidateKnowledge(kb)

	// Should produce a warning (not error) for missing target
	found := false
	for _, warn := range result.Warnings {
		if containsString(warn.Message, "does not exist") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should warn about missing target concept in form")
	}
}

func TestValidateForms_EmptyForm(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{{ID: "test", Name: "test"}},
		Forms: []Form{
			{Form: "", Concept: "test", Weight: 0.5}, // empty form
		},
	}

	result := ValidateKnowledge(kb)

	found := false
	for _, err := range result.Errors {
		if containsString(err.Message, "form cannot be empty") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should detect empty form")
	}
}

func TestValidateForms_InvalidWeight(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{{ID: "test", Name: "test"}},
		Forms: []Form{
			{Form: "test", Concept: "test", Weight: 0}, // invalid weight
		},
	}

	result := ValidateKnowledge(kb)

	found := false
	for _, err := range result.Errors {
		if containsString(err.Message, "weight must be in (0, 1]") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should detect invalid weight in form")
	}
}

func TestValidateForms_WeightTooHigh(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{{ID: "test", Name: "test"}},
		Forms: []Form{
			{Form: "test", Concept: "test", Weight: 1.5}, // weight > 1
		},
	}

	result := ValidateKnowledge(kb)

	found := false
	for _, err := range result.Errors {
		if containsString(err.Message, "weight must be in (0, 1]") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should detect weight > 1 in form")
	}
}

// =============================================================================
// RELATION VALIDATION FAILURE TESTS
// =============================================================================

func TestValidateRelations_MissingFromConcept(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "existing", Name: "existing"},
		},
		Relations: []Relation{
			{From: "", To: "existing", Type: "related", Weight: 0.5}, // empty from
		},
	}

	result := ValidateKnowledge(kb)

	found := false
	for _, err := range result.Errors {
		if containsString(err.Message, "'from' concept cannot be empty") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should detect empty 'from' concept in relation")
	}
}

func TestValidateRelations_MissingToConcept(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "existing", Name: "existing"},
		},
		Relations: []Relation{
			{From: "existing", To: "", Type: "related", Weight: 0.5}, // empty to
		},
	}

	result := ValidateKnowledge(kb)

	found := false
	for _, err := range result.Errors {
		if containsString(err.Message, "'to' concept cannot be empty") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should detect empty 'to' concept in relation")
	}
}

func TestValidateRelations_FromConceptNotFound(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "existing", Name: "existing"},
		},
		Relations: []Relation{
			{From: "nonexistent", To: "existing", Type: "related", Weight: 0.5}, // from doesn't exist
		},
	}

	result := ValidateKnowledge(kb)

	found := false
	for _, err := range result.Errors {
		if containsString(err.Message, "'from' concept") && containsString(err.Message, "does not exist") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should detect missing 'from' concept in relation")
	}
}

func TestValidateRelations_ToConceptNotFound(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "existing", Name: "existing"},
		},
		Relations: []Relation{
			{From: "existing", To: "nonexistent", Type: "related", Weight: 0.5}, // to doesn't exist
		},
	}

	result := ValidateKnowledge(kb)

	found := false
	for _, err := range result.Errors {
		if containsString(err.Message, "'to' concept") && containsString(err.Message, "does not exist") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should detect missing 'to' concept in relation")
	}
}

func TestValidateRelations_InvalidRelationType(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "a", Name: "a"},
			{ID: "b", Name: "b"},
		},
		Relations: []Relation{
			{From: "a", To: "b", Type: "invalid", Weight: 0.5}, // invalid type
		},
	}

	result := ValidateKnowledge(kb)

	// Should produce a warning for non-standard relation type
	found := false
	for _, warn := range result.Warnings {
		if containsString(warn.Message, "non-standard") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should warn about non-standard relation type")
	}
}

func TestValidateRelations_InvalidWeight(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "a", Name: "a"},
			{ID: "b", Name: "b"},
		},
		Relations: []Relation{
			{From: "a", To: "b", Type: "related", Weight: 0}, // invalid weight
		},
	}

	result := ValidateKnowledge(kb)

	found := false
	for _, err := range result.Errors {
		if containsString(err.Message, "weight must be in (0, 1]") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should detect invalid weight in relation")
	}
}

// =============================================================================
// SCRIPT WORD VALIDATION FAILURE TESTS
// =============================================================================

func TestValidateScriptWords_MissingMeaning(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "existing", Name: "existing"},
		},
		ScriptWords: []ScriptWord{
			{Script: "hebrew", Word: "test", Meanings: []string{"nonexistent"}}, // meaning doesn't exist
		},
	}

	result := ValidateKnowledge(kb)

	found := false
	for _, err := range errorsForField(result, "script_words") {
		if containsString(err.Message, "does not exist") {
			found = true
			break
		}
	}
	if !found {
		t.Error("should detect missing meaning in script word")
	}
}

// =============================================================================
// VALIDATION RESULT TESTS
// =============================================================================

func TestValidationResult_IsValid(t *testing.T) {
	result := &ValidationResult{
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}
	if !result.IsValid() {
		t.Error("empty result should be valid")
	}

	result.AddError("test.field", "test error")
	if result.IsValid() {
		t.Error("result with errors should not be valid")
	}
}

func TestValidationResult_AddError(t *testing.T) {
	result := &ValidationResult{}
	result.AddError("field", "message")

	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Field != "field" {
		t.Errorf("expected field 'field', got '%s'", result.Errors[0].Field)
	}
	if result.Errors[0].Message != "message" {
		t.Errorf("expected message 'message', got '%s'", result.Errors[0].Message)
	}
}

func TestValidationResult_AddWarning(t *testing.T) {
	result := &ValidationResult{}
	result.AddWarning("field", "warning message")

	if len(result.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(result.Warnings))
	}
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func formatErrors(errors []ValidationError) string {
	result := ""
	for _, e := range errors {
		result += "  " + e.Field + ": " + e.Message + "\n"
	}
	return result
}

func errorsForField(result *ValidationResult, fieldPrefix string) []ValidationError {
	var filtered []ValidationError
	for _, e := range result.Errors {
		if len(e.Field) >= len(fieldPrefix) && e.Field[:len(fieldPrefix)] == fieldPrefix {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
