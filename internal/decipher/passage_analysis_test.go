package decipher

import (
	"testing"
)

// PASSAGE FIELD ANALYSIS TESTS
// =============================================================================

func TestAnalyzePassage_MultiToken(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("God is Love", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Should have passage fields
	if len(fields) == 0 {
		t.Error("should produce passage fields for multi-token input")
	}

	// Verify 'love' field exists and is strong
	loveField := fields.GetField("love")
	if loveField == nil {
		t.Error("should have 'love' field")
	} else if loveField.Strength <= 0 {
		t.Error("love field should have positive strength")
	}
}

func TestAnalyzePassage_SingleToken(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("love", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Single token should still produce fields
	if len(fields) == 0 {
		t.Error("should produce passage fields for single token")
	}
}

func TestAnalyzePassage_Empty(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("", engine)
	if err != nil {
		t.Fatal(err)
	}

	if fields != nil {
		t.Error("empty passage should produce nil fields")
	}
}

func TestAnalyzePassage_WhitespaceOnly(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("   \n\t  ", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Whitespace-only should produce no fields
	if len(fields) > 0 {
		t.Errorf("whitespace-only should produce no fields, got %d", len(fields))
	}
}

// =============================================================================
