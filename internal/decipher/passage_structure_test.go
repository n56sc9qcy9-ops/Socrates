package decipher

import (
	"testing"
)

// PASSAGE FIELD STRUCTURE TESTS
// =============================================================================

func TestPassageField_Structure(t *testing.T) {
	field := &PassageField{
		Concept:      "love",
		Strength:     0.9,
		Confidence:   ConfidenceVerified,
		TokenSources: []string{"love", "heart"},
		Depth:        0,
	}

	if field.Concept != "love" {
		t.Errorf("expected concept 'love', got '%s'", field.Concept)
	}
	if field.Strength != 0.9 {
		t.Errorf("expected strength 0.9, got %f", field.Strength)
	}
	if field.Depth != 0 {
		t.Errorf("expected depth 0 (direct), got %d", field.Depth)
	}
}

func TestPassageFields_TopFields(t *testing.T) {
	fields := PassageFields{
		{Concept: "low", Strength: 0.2},
		{Concept: "high", Strength: 0.9},
		{Concept: "medium", Strength: 0.5},
	}

	top2 := fields.TopFields(2)

	if len(top2) != 2 {
		t.Errorf("TopFields(2) should return 2 fields, got %d", len(top2))
	}
	if top2[0].Concept != "high" {
		t.Errorf("top field should be 'high', got '%s'", top2[0].Concept)
	}
	if top2[1].Concept != "medium" {
		t.Errorf("second field should be 'medium', got '%s'", top2[1].Concept)
	}
}

func TestPassageFields_GetField(t *testing.T) {
	fields := PassageFields{
		{Concept: "love", Strength: 0.9},
		{Concept: "light", Strength: 0.7},
	}

	loveField := fields.GetField("love")
	if loveField == nil {
		t.Fatal("should find 'love' field")
	}
	if loveField.Strength != 0.9 {
		t.Errorf("expected strength 0.9, got %f", loveField.Strength)
	}

	missingField := fields.GetField("nonexistent")
	if missingField != nil {
		t.Error("should not find 'nonexistent' field")
	}
}

func TestPassageFields_Merge(t *testing.T) {
	fields1 := PassageFields{
		{Concept: "love", Strength: 0.5, TokenSources: []string{"love"}},
		{Concept: "light", Strength: 0.3, TokenSources: []string{"light"}},
	}

	fields2 := PassageFields{
		{Concept: "love", Strength: 0.4, TokenSources: []string{"heart"}},
		{Concept: "truth", Strength: 0.6, TokenSources: []string{"truth"}},
	}

	merged := fields1.Merge(fields2)

	// Should have 3 unique concepts
	if len(merged) != 3 {
		t.Errorf("merged should have 3 fields, got %d", len(merged))
	}

	// Love should have combined strength
	loveField := merged.GetField("love")
	if loveField == nil {
		t.Fatal("should have love field")
	}
	if loveField.Strength != 0.9 {
		t.Errorf("love should have combined strength 0.9, got %f", loveField.Strength)
	}
	// Love should have both token sources
	if len(loveField.TokenSources) != 2 {
		t.Errorf("love should have 2 token sources, got %d", len(loveField.TokenSources))
	}
}

// =============================================================================
