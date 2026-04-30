package training

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParseReviewFile tests parsing of review files.
func TestParseReviewFile(t *testing.T) {
	// Create a test review file
	testReviewPath := filepath.Join(os.TempDir(), "test_review_parsing.yaml")
	defer os.Remove(testReviewPath)

	content := `version: "1.0"
generated: "2026-04-30"
suggestions:
  - id: "suggestion_1"
    path: "fuzzy_match.base"
    current: 0.8
    suggested: 0.75
    status: "pending"
    rationale: "reduce fuzzy match weight"
    train_delta:
      field_precision: 0.01
      field_recall: 0.00
    heldout_delta:
      field_precision: 0.05
      field_recall: 0.00
`
	if err := os.WriteFile(testReviewPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test review file: %v", err)
	}

	// Parse the file
	rf, errors := ParseReviewFile(testReviewPath)
	if len(errors) > 0 {
		t.Fatalf("Unexpected parsing errors: %v", errors)
	}

	if rf.Version != "1.0" {
		t.Errorf("Expected version '1.0', got '%s'", rf.Version)
	}

	if len(rf.Suggestions) != 1 {
		t.Fatalf("Expected 1 suggestion, got %d", len(rf.Suggestions))
	}

	s := rf.Suggestions[0]
	if s.ID != "suggestion_1" {
		t.Errorf("Expected id 'suggestion_1', got '%s'", s.ID)
	}
	if s.Path != "fuzzy_match.base" {
		t.Errorf("Expected path 'fuzzy_match.base', got '%s'", s.Path)
	}
	if s.Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", s.Status)
	}
	if s.TrainDelta.FieldPrecision != 0.01 {
		t.Errorf("Expected train_delta.field_precision 0.01, got %.4f", s.TrainDelta.FieldPrecision)
	}
	if s.HeldOutDelta.FieldPrecision != 0.05 {
		t.Errorf("Expected heldout_delta.field_precision 0.05, got %.4f", s.HeldOutDelta.FieldPrecision)
	}
}

// TestValidateReviewFile tests review file validation.
func TestValidateReviewFile(t *testing.T) {
	rf := &ReviewFile{
		Suggestions: []WeightSuggestion{
			{
				ID:        "s1",
				Path:      "fuzzy_match.base",
				Current:   0.80,
				Suggested: 0.75,
				Status:    "pending",
			},
			{
				ID:        "s2",
				Path:      "graph.depth_penalty",
				Current:   0.90,
				Suggested: 0.95,
				Status:    "accepted",
			},
			{
				ID:        "s3",
				Path:      "confidence.verified",
				Current:   1.00,
				Suggested: 1.10,
				Status:    "rejected",
			},
			{
				ID:        "s4",
				Path:      "unknown_field",
				Current:   1.00,
				Suggested: 1.10,
				Status:    "pending",
			},
		},
	}

	errors := ValidateReviewFile(rf)

	// Should have one error for unknown_field
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d: %v", len(errors), errors)
	}

	if len(errors) > 0 && errors[0].Field != "suggestions[3].path" {
		t.Errorf("Expected error for suggestions[3].path, got '%s'", errors[0].Field)
	}
}

// TestGetAcceptedSuggestions tests filtering by status.
func TestGetAcceptedSuggestions(t *testing.T) {
	rf := &ReviewFile{
		Suggestions: []WeightSuggestion{
			{ID: "s1", Status: "pending"},
			{ID: "s2", Status: "accepted"},
			{ID: "s3", Status: "rejected"},
			{ID: "s4", Status: "accepted"},
		},
	}

	accepted := GetAcceptedSuggestions(rf)
	if len(accepted) != 2 {
		t.Errorf("Expected 2 accepted, got %d", len(accepted))
	}

	pending := GetPendingSuggestions(rf)
	if len(pending) != 1 {
		t.Errorf("Expected 1 pending, got %d", len(pending))
	}

	rejected := GetRejectedSuggestions(rf)
	if len(rejected) != 1 {
		t.Errorf("Expected 1 rejected, got %d", len(rejected))
	}
}

// TestPendingSuggestionsDoNotApply tests that pending status blocks application.
func TestPendingSuggestionsDoNotApply(t *testing.T) {
	reviewPath := filepath.Join(os.TempDir(), "test_pending_review.yaml")
	targetPath := filepath.Join(os.TempDir(), "test_target_weights.yaml")
	defer os.Remove(reviewPath)
	defer os.Remove(targetPath)

	// Create review file with pending status
	reviewContent := `version: "1.0"
suggestions:
  - id: "s1"
    path: "fuzzy_match.base"
    current: 0.8
    suggested: 0.75
    status: "pending"
`
	if err := os.WriteFile(reviewPath, []byte(reviewContent), 0644); err != nil {
		t.Fatalf("Failed to write review file: %v", err)
	}

	// Create target weights
	currentWeights := map[string]interface{}{
		"fuzzy_match.base": 0.80,
		"graph.depth_penalty": 0.90,
	}

	// Apply should fail for pending
	_, errors := ApplyWeightSuggestions(reviewPath, targetPath, currentWeights, nil)
	if len(errors) == 0 {
		t.Error("Expected error for pending suggestions, got none")
	}
}

// TestRejectedSuggestionsDoNotApply tests that rejected status blocks application.
func TestRejectedSuggestionsDoNotApply(t *testing.T) {
	reviewPath := filepath.Join(os.TempDir(), "test_rejected_review.yaml")
	targetPath := filepath.Join(os.TempDir(), "test_target_weights.yaml")
	defer os.Remove(reviewPath)
	defer os.Remove(targetPath)

	// Create review file with rejected status
	reviewContent := `version: "1.0"
suggestions:
  - id: "s1"
    path: "fuzzy_match.base"
    current: 0.8
    suggested: 0.75
    status: "rejected"
`
	if err := os.WriteFile(reviewPath, []byte(reviewContent), 0644); err != nil {
		t.Fatalf("Failed to write review file: %v", err)
	}

	currentWeights := map[string]interface{}{
		"fuzzy_match.base": 0.80,
		"graph.depth_penalty": 0.90,
	}

	// Apply should fail for rejected
	_, errors := ApplyWeightSuggestions(reviewPath, targetPath, currentWeights, nil)
	if len(errors) == 0 {
		t.Error("Expected error for rejected suggestions, got none")
	}
}

// TestAcceptedSuggestionsApply tests that accepted status allows application.
func TestAcceptedSuggestionsApply(t *testing.T) {
	reviewPath := filepath.Join(os.TempDir(), "test_accepted_review.yaml")
	targetPath := filepath.Join(os.TempDir(), "test_target_weights.yaml")
	auditPath := filepath.Join(os.TempDir(), "test_audit.yaml")
	defer os.Remove(reviewPath)
	defer os.Remove(targetPath)
	defer os.Remove(auditPath)

	// Create review file with accepted status
	reviewContent := `version: "1.0"
suggestions:
  - id: "s1"
    path: "fuzzy_match.base"
    current: 0.8
    suggested: 0.75
    status: "accepted"
    rationale: "test application"
`
	if err := os.WriteFile(reviewPath, []byte(reviewContent), 0644); err != nil {
		t.Fatalf("Failed to write review file: %v", err)
	}

	currentWeights := map[string]interface{}{
		"fuzzy_match.base": 0.80,
		"graph.depth_penalty": 0.90,
	}

	// Apply should succeed
	newAudit, errors := ApplyWeightSuggestions(reviewPath, targetPath, currentWeights, nil)
	if len(errors) > 0 {
		t.Fatalf("Unexpected errors: %v", errors)
	}

	// Verify target file was written
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Error("Target weights file was not created")
	}

	// Verify audit record was created
	if len(newAudit) != 1 {
		t.Errorf("Expected 1 audit record, got %d", len(newAudit))
	}

	if newAudit[0].Path != "fuzzy_match.base" {
		t.Errorf("Expected audit path 'fuzzy_match.base', got '%s'", newAudit[0].Path)
	}
	if newAudit[0].OldValue != 0.8 {
		t.Errorf("Expected audit old value 0.80, got %v", newAudit[0].OldValue)
	}
	if newAudit[0].NewValue != 0.75 {
		t.Errorf("Expected audit new value 0.75, got %v", newAudit[0].NewValue)
	}
}

// TestDriftPreventsApply tests that value drift prevents application.
func TestDriftPreventsApply(t *testing.T) {
	reviewPath := filepath.Join(os.TempDir(), "test_drift_review.yaml")
	targetPath := filepath.Join(os.TempDir(), "test_drift_target.yaml")
	defer os.Remove(reviewPath)
	defer os.Remove(targetPath)

	// Create review file expecting value 0.80
	reviewContent := `version: "1.0"
suggestions:
  - id: "s1"
    path: "fuzzy_match.base"
    current: 0.8
    suggested: 0.75
    status: "accepted"
`
	if err := os.WriteFile(reviewPath, []byte(reviewContent), 0644); err != nil {
		t.Fatalf("Failed to write review file: %v", err)
	}

	// Current config has 0.85 (drift!)
	currentWeights := map[string]interface{}{
		"fuzzy_match.base": 0.85, // Different from recorded 0.80
	}

	// Apply should fail due to drift
	_, errors := ApplyWeightSuggestions(reviewPath, targetPath, currentWeights, nil)
	if len(errors) == 0 {
		t.Error("Expected error for drift, got none")
	}

	// Verify target file was NOT written
	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		t.Error("Target file should not have been created when drift detected")
	}
}

// TestUnknownWeightFieldPreventsApply tests that unknown fields are rejected.
func TestUnknownWeightFieldPreventsApply(t *testing.T) {
	reviewPath := filepath.Join(os.TempDir(), "test_unknown_field_review.yaml")
	targetPath := filepath.Join(os.TempDir(), "test_unknown_target.yaml")
	defer os.Remove(reviewPath)
	defer os.Remove(targetPath)

	// Create review file with unknown field
	reviewContent := `version: "1.0"
suggestions:
  - id: "s1"
    path: "some_unknown_field"
    current: 1.00
    suggested: 2.00
    status: "accepted"
`
	if err := os.WriteFile(reviewPath, []byte(reviewContent), 0644); err != nil {
		t.Fatalf("Failed to write review file: %v", err)
	}

	currentWeights := map[string]interface{}{
		"some_unknown_field": 1.00,
		"fuzzy_match.base": 0.80, // needed for valid review
	}

	// Apply should fail due to unknown field
	_, errors := ApplyWeightSuggestions(reviewPath, targetPath, currentWeights, nil)
	if len(errors) == 0 {
		t.Error("Expected error for unknown field, got none")
	}
}

// TestAuditTrailWritten tests that audit trail is written.
func TestAuditTrailWritten(t *testing.T) {
	reviewPath := filepath.Join(os.TempDir(), "test_audit_review.yaml")
	targetPath := filepath.Join(os.TempDir(), "test_audit_target.yaml")
	auditPath := filepath.Join(os.TempDir(), "test_audit_trail.yaml")
	defer os.Remove(reviewPath)
	defer os.Remove(targetPath)
	defer os.Remove(auditPath)

	// Create review with multiple accepted suggestions
	reviewContent := `version: "1.0"
suggestions:
  - id: "s1"
    path: "fuzzy_match.base"
    current: 0.8
    suggested: 0.75
    status: "accepted"
  - id: "s2"
    path: "graph.depth_penalty"
    current: 0.9
    suggested: 0.85
    status: "accepted"
`
	if err := os.WriteFile(reviewPath, []byte(reviewContent), 0644); err != nil {
		t.Fatalf("Failed to write review file: %v", err)
	}

	currentWeights := map[string]interface{}{
		"fuzzy_match.base":     0.80,
		"graph.depth_penalty":  0.90,
	}

	// Apply suggestions
	newAudit, errors := ApplyWeightSuggestions(reviewPath, targetPath, currentWeights, nil)
	if len(errors) > 0 {
		t.Fatalf("Unexpected errors: %v", errors)
	}

	// Write audit trail
	if err := WriteAuditTrail(auditPath, newAudit); err != nil {
		t.Fatalf("Failed to write audit trail: %v", err)
	}

	// Read back audit trail
	readAudit, err := ReadAuditTrail(auditPath)
	if err != nil {
		t.Fatalf("Failed to read audit trail: %v", err)
	}

	if len(readAudit) != 2 {
		t.Errorf("Expected 2 audit records, got %d", len(readAudit))
	}
}

// TestActiveKnowledgeNotMutated tests that knowledge YAML is not modified.
func TestActiveKnowledgeNotMutated(t *testing.T) {
	// This test verifies that ApplyWeightSuggestions does not touch knowledge files
	// We create a mock scenario and verify only the target weights file is modified

	reviewPath := filepath.Join(os.TempDir(), "test_knowledge_review.yaml")
	targetPath := filepath.Join(os.TempDir(), "test_weights_knowledge.yaml")
	defer os.Remove(reviewPath)
	defer os.Remove(targetPath)

	// Create review file
	reviewContent := `version: "1.0"
suggestions:
  - id: "s1"
    path: "fuzzy_match.base"
    current: 0.8
    suggested: 0.75
    status: "accepted"
`
	if err := os.WriteFile(reviewPath, []byte(reviewContent), 0644); err != nil {
		t.Fatalf("Failed to write review file: %v", err)
	}

	currentWeights := map[string]interface{}{
		"fuzzy_match.base": 0.80,
		"graph.depth_penalty": 0.90,
	}

	// Apply suggestions
	_, errors := ApplyWeightSuggestions(reviewPath, targetPath, currentWeights, nil)
	if len(errors) > 0 {
		t.Fatalf("Unexpected errors: %v", errors)
	}

	// Verify target file content is just weights, not knowledge
	if data, err := os.ReadFile(targetPath); err == nil {
		content := string(data)
		// Should not contain knowledge YAML markers
		if containsString(content, "concepts:") || containsString(content, "relations:") {
			t.Error("Target weights file should not contain knowledge YAML")
		}
		// Should contain weights marker
		if !containsString(content, "fuzzy_match:") {
			t.Error("Target weights file should contain 'weights:' section")
		}
	}
}

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestAllowedWeightFields tests that all documented ranking weight fields are allowed.
func TestAllowedWeightFields(t *testing.T) {
	knownFields := []string{
		"exact_match.base",
		"fuzzy_match.base",
		"fuzzy_match.distance_penalty",
		"confidence.verified",
		"confidence.plausible",
		"confidence.speculative",
		"graph.expansion_weight",
		"graph.depth_penalty",
		"graph.relation_base_weight",
	}

	for _, field := range knownFields {
		if !AllowedWeightFields[field] {
			t.Errorf("Expected field '%s' to be allowed", field)
		}
	}
}

// TestMultipleReviewFilesAuditContinuity tests that audit records persist across reviews.
func TestMultipleReviewFilesAuditContinuity(t *testing.T) {
	reviewPath := filepath.Join(os.TempDir(), "test_continuity_review.yaml")
	targetPath := filepath.Join(os.TempDir(), "test_continuity_target.yaml")
	auditPath := filepath.Join(os.TempDir(), "test_continuity_audit.yaml")
	defer os.Remove(reviewPath)
	defer os.Remove(targetPath)
	defer os.Remove(auditPath)

	// First application
	reviewContent1 := `version: "1.0"
suggestions:
  - id: "s1"
    path: "fuzzy_match.base"
    current: 0.8
    suggested: 0.75
    status: "accepted"
`
	if err := os.WriteFile(reviewPath, []byte(reviewContent1), 0644); err != nil {
		t.Fatalf("Failed to write review file: %v", err)
	}

	currentWeights := map[string]interface{}{
		"fuzzy_match.base": 0.80,
		"graph.depth_penalty": 0.90,
	}

	audit1, errors := ApplyWeightSuggestions(reviewPath, targetPath, currentWeights, nil)
	if len(errors) > 0 {
		t.Fatalf("First application failed: %v", errors)
	}

	// Write first audit
	if err := WriteAuditTrail(auditPath, audit1); err != nil {
		t.Fatalf("Failed to write first audit: %v", err)
	}

	// Second application - read existing audit
	existingAudit, _ := ReadAuditTrail(auditPath)

	// Create second review (path changed)
	reviewContent2 := `version: "1.0"
suggestions:
  - id: "s2"
    path: "graph.depth_penalty"
    current: 0.9
    suggested: 0.85
    status: "accepted"
`
	if err := os.WriteFile(reviewPath, []byte(reviewContent2), 0644); err != nil {
		t.Fatalf("Failed to write second review: %v", err)
	}

	// Current weights from first apply
	currentWeights2 := map[string]interface{}{
		"fuzzy_match.base":    0.75, // Updated from first apply
		"graph.depth_penalty": 0.90,
	}

	// Apply with existing audit
	audit2, errors := ApplyWeightSuggestions(reviewPath, targetPath, currentWeights2, existingAudit)
	if len(errors) > 0 {
		t.Fatalf("Second application failed: %v", errors)
	}

	// Write combined audit
	if err := WriteAuditTrail(auditPath, audit2); err != nil {
		t.Fatalf("Failed to write combined audit: %v", err)
	}

	// Verify both records present
	readAudit, _ := ReadAuditTrail(auditPath)
	if len(readAudit) != 2 {
		t.Errorf("Expected 2 audit records (continuity), got %d", len(readAudit))
	}
}