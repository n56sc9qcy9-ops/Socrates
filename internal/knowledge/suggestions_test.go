package knowledge

import (
	"testing"
)

// =============================================================================
// SUGGESTION TESTS
// =============================================================================

func TestSuggestionReviewer_AddAndGet(t *testing.T) {
	reviewer := NewSuggestionReviewer()
	
	s := Suggestion{
		Type:      "form",
		ProposedForm: "test",
		Source:    "test",
		Evidence:  "test evidence",
		Confidence: "speculative",
	}
	
	reviewer.AddSuggestion(s)
	
	suggestions := reviewer.GetSuggestions()
	if len(suggestions) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(suggestions))
	}
}

func TestSuggestionReviewer_GetUnreviewed(t *testing.T) {
	reviewer := NewSuggestionReviewer()
	
	reviewer.AddSuggestion(Suggestion{Type: "form", Reviewed: false})
	reviewer.AddSuggestion(Suggestion{Type: "form", Reviewed: true})
	reviewer.AddSuggestion(Suggestion{Type: "form", Reviewed: false})
	
	unreviewed := reviewer.GetUnreviewedSuggestions()
	if len(unreviewed) != 2 {
		t.Errorf("expected 2 unreviewed suggestions, got %d", len(unreviewed))
	}
}

func TestSuggestionReviewer_MarkReviewed(t *testing.T) {
	reviewer := NewSuggestionReviewer()
	
	reviewer.AddSuggestion(Suggestion{Type: "form", Reviewed: false})
	
	reviewer.MarkReviewed(0, true, "test-reviewer")
	
	suggestions := reviewer.GetSuggestions()
	if !suggestions[0].Reviewed {
		t.Error("suggestion should be marked as reviewed")
	}
	if !suggestions[0].Accepted {
		t.Error("suggestion should be marked as accepted")
	}
	if suggestions[0].ReviewedBy != "test-reviewer" {
		t.Errorf("expected reviewer 'test-reviewer', got '%s'", suggestions[0].ReviewedBy)
	}
}

func TestSuggestionForUnknownInput_KnownForm(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}
	
	// "skall" should be a known form
	suggestion := SuggestionForUnknownInput("skall", kb)
	if suggestion != nil {
		t.Error("known form should not generate a suggestion")
	}
}

func TestSuggestionForWeaklyMatchedInput(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}
	
	// Test with no matches - should suggest
	matches := kb.GetFormsByText("completely_unknown_xyz")
	suggestions := SuggestionForWeaklyMatchedInput("completely_unknown_xyz", matches, kb)
	
	if len(suggestions) == 0 {
		t.Error("no matches should generate a suggestion")
	}
	
	if suggestions[0].Type != "form" {
		t.Errorf("expected type 'form', got '%s'", suggestions[0].Type)
	}
}

func TestSuggestionForWeaklyMatchedInput_StrongMatch(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}
	
	// Test with strong matches - should not suggest
	matches := kb.GetFormsByText("skall")
	suggestions := SuggestionForWeaklyMatchedInput("skall", matches, kb)
	
	if len(suggestions) > 0 {
		t.Error("strong matches should not generate a suggestion")
	}
}

func TestFormatSuggestionsForReview(t *testing.T) {
	suggestions := []Suggestion{
		{
			Type:          "form",
			ProposedForm:  "testform",
			Source:        "test",
			Evidence:      "test evidence",
			Confidence:    "speculative",
			Weight:        0.5,
			Rationale:     "test rationale",
			Reviewed:      false,
		},
	}
	
	output := FormatSuggestionsForReview(suggestions)
	
	if len(output) == 0 {
		t.Error("format should produce non-empty output")
	}
	
	// Check that expected content is present
	expected := []string{
		"form",
		"testform",
		"test",
		"test evidence",
		"speculative",
		"test rationale",
	}
	
	for _, exp := range expected {
		if !contains(output, exp) {
			t.Errorf("output should contain '%s'", exp)
		}
	}
}

func TestFormatSuggestionsForReview_Empty(t *testing.T) {
	output := FormatSuggestionsForReview([]Suggestion{})
	if output != "No suggestions to review.\n" {
		t.Error("empty suggestions should produce specific message")
	}
}

func TestFormatSuggestionsForReview_Accepted(t *testing.T) {
	suggestions := []Suggestion{
		{
			Type:        "form",
			ProposedForm: "test",
			Reviewed:    true,
			Accepted:    true,
			ReviewedBy: "curator",
		},
	}
	
	output := FormatSuggestionsForReview(suggestions)
	if !contains(output, "ACCEPTED") {
		t.Error("output should contain ACCEPTED status")
	}
	if !contains(output, "curator") {
		t.Error("output should contain reviewer name")
	}
}

func TestFormatSuggestionsForReview_Rejected(t *testing.T) {
	suggestions := []Suggestion{
		{
			Type:        "form",
			ProposedForm: "test",
			Reviewed:    true,
			Accepted:    false,
			ReviewedBy: "curator",
		},
	}
	
	output := FormatSuggestionsForReview(suggestions)
	if !contains(output, "REJECTED") {
		t.Error("output should contain REJECTED status")
	}
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// =============================================================================
// MINIMUM REVIEW RECORD TESTS
// =============================================================================

func TestMinimumReviewRecord_Concept(t *testing.T) {
	rec := MinimumReviewRecord{
		Type:          "concept",
		ConceptID:     "new-concept",
		ConceptName:   "New Concept",
		Aliases:       []string{"alias1", "alias2"},
		EvidenceSource: "etymology source",
		ConfidenceLevel: "verified",
		Weight:        1.0,
		Rationale:    "founded in etymology",
		Notes:         "additional notes",
		Accepted:      true,
		ReviewedBy:    "curator",
	}
	
	if rec.Type != "concept" {
		t.Errorf("expected type 'concept', got '%s'", rec.Type)
	}
	if rec.ConceptID != "new-concept" {
		t.Errorf("expected concept id 'new-concept', got '%s'", rec.ConceptID)
	}
}

func TestMinimumReviewRecord_Relation(t *testing.T) {
	rec := MinimumReviewRecord{
		Type:          "relation",
		From:          "concept-a",
		To:            "concept-b",
		RelType:       "related",
		RelWeight:     0.8,
		EvidenceSource: "textual analysis",
		ConfidenceLevel: "plausible",
		Accepted:      true,
		ReviewedBy:    "curator",
	}
	
	if rec.Type != "relation" {
		t.Errorf("expected type 'relation', got '%s'", rec.Type)
	}
	if rec.From != "concept-a" || rec.To != "concept-b" {
		t.Error("relation endpoints mismatch")
	}
}
