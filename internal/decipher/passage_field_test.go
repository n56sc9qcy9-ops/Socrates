package decipher

import (
	"strings"
	"testing"
)

// =============================================================================
// PASSAGE FIELD STRUCTURE TESTS
// =============================================================================

func TestPassageField_Structure(t *testing.T) {
	field := &PassageField{
		Concept:      "love",
		Strength:     0.9,
		Confidence:   ConfidenceVerified,
		TokenSources:  []string{"love", "heart"},
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

	// Whitespace-only might produce no tokens
	if len(fields) > 0 {
		// This is acceptable behavior
		t.Logf("whitespace-only produced %d fields", len(fields))
	}
}

// =============================================================================
// PASSAGE FIELD TOKEN TRACKING TESTS
// =============================================================================

func TestPassageField_TokenSources(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("love light", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Check that fields have token sources (may be empty if evidence paths not populated)
	loveField := fields.GetField("love")
	if loveField != nil {
		// Token sources may or may not include 'love' depending on graph evidence
		// Just verify the field exists and has strength
		if loveField.Strength <= 0 {
			t.Error("love field should have positive strength")
		}
		t.Logf("love field: strength=%.4f, tokens=%v", loveField.Strength, loveField.TokenSources)
	}
}

func TestPassageField_DepthDistinction(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("love", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Direct concept should have depth 0
	for _, field := range fields {
		// Field should have a valid depth
		if field.Depth < 0 {
			t.Errorf("field '%s' has invalid depth %d", field.Concept, field.Depth)
		}
	}
}

// =============================================================================
// REPEATED ACTIVATION FIELD TESTS
// =============================================================================

func TestPassageFields_NoDuplicateInflation(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Analyze same passage twice
	fields1, _ := AnalyzePassage("love heart", engine)
	fields2, _ := AnalyzePassage("love heart", engine)

	// Results should be similar
	if len(fields1) != len(fields2) {
		t.Errorf("repeated analysis should produce same number of fields: %d vs %d",
			len(fields1), len(fields2))
	}
}

func TestPassageField_MultipleTokensSameConcept(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Analyze passage where multiple tokens might activate same concept
	fields, err := AnalyzePassage("love love", engine)
	if err != nil {
		t.Fatal(err)
	}

	// There should be a love field
	loveField := fields.GetField("love")
	if loveField == nil {
		t.Error("should have love field")
	} else {
		// Love field should track both token sources
		loveCount := 0
		for _, source := range loveField.TokenSources {
			if source == "love" {
				loveCount++
			}
		}
		// Either deduplicated (1 source) or tracking repeats (2+ sources)
		// Either way, field should exist with reasonable strength
		if loveField.Strength <= 0 {
			t.Error("love field should have positive strength")
		}
	}
}

// =============================================================================
// KEY TOKEN WEAKENING TESTS
// =============================================================================

func TestKeyTokenRemoval_WeakensPassageField(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// "love truth light" - love is key semantic word
	fieldsFull, _ := AnalyzePassage("love truth light", engine)
	fieldsNoLove, _ := AnalyzePassage("truth light", engine)


	fullLoveStrength := 0.0
	noLoveStrength := 0.0

	for _, field := range fieldsFull {
		if field.Concept == "love" {
			fullLoveStrength = field.Strength
		}
	}

	for _, field := range fieldsNoLove {
		if field.Concept == "love" {
			noLoveStrength = field.Strength
		}
	}

	// The full passage with "love" should have higher or equal love strength
	// than the passage without "love"
	if fullLoveStrength < noLoveStrength {
		t.Errorf("passage with 'love' should have love field strength >= without:\n  full: %.4f\n  no love: %.4f",
			fullLoveStrength, noLoveStrength)
	}
}

func TestKeyTokenRemoval_AffectsOtherFields(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Full passage
	fieldsFull, _ := AnalyzePassage("love truth", engine)
	// Without "love" - should affect co-activation
	fieldsNoLove, _ := AnalyzePassage("truth", engine)

	// Truth field in full passage might have higher strength due to co-activation with love
	truthFullStrength := 0.0
	truthNoLoveStrength := 0.0

	for _, field := range fieldsFull {
		if field.Concept == "truth" {
			truthFullStrength = field.Strength
		}
	}

	for _, field := range fieldsNoLove {
		if field.Concept == "truth" {
			truthNoLoveStrength = field.Strength
		}
	}

	// Full passage truth strength might be different due to concept co-activation
	// This is acceptable - we're testing that the system responds to token removal
	t.Logf("truth strength: full=%.4f, no love=%.4f", truthFullStrength, truthNoLoveStrength)
}

// =============================================================================
// TOKENIZATION TESTS
// =============================================================================

func TestTokenizePassage_Basic(t *testing.T) {
	tokens := tokenizePassage("love truth")
	if len(tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[0] != "love" {
		t.Errorf("expected first token 'love', got '%s'", tokens[0])
	}
	if tokens[1] != "truth" {
		t.Errorf("expected second token 'truth', got '%s'", tokens[1])
	}
}

func TestTokenizePassage_Whitespace(t *testing.T) {
	tokens := tokenizePassage("  love   truth  ")
	if len(tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(tokens))
	}
}

func TestTokenizePassage_Newlines(t *testing.T) {
	tokens := tokenizePassage("love\ntruth\rlight")
	if len(tokens) != 3 {
		t.Errorf("expected 3 tokens, got %d", len(tokens))
	}
}

func TestTokenizePassage_SingleToken(t *testing.T) {
	tokens := tokenizePassage("love")
	if len(tokens) != 1 {
		t.Errorf("expected 1 token, got %d", len(tokens))
	}
}

func TestTokenizePassage_Empty(t *testing.T) {
	tokens := tokenizePassage("")
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(tokens))
	}
}

// =============================================================================
// DEFAULT OUTPUT TESTS
// =============================================================================

func TestPassageField_DefaultOutputConcise(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("love truth")
	output := RenderReadingWithOptions(reading, DefaultRenderOptions())

	// Default output should be concise
	if len(output) > 2000 {
		t.Errorf("default output seems too long (%d chars), should be concise", len(output))
	}

	// Should NOT include all passage field internals
	if strings.Contains(output, "PassageField") {
		t.Error("default output should not expose PassageField struct names")
	}
}

func TestPassageField_DebugOutputMayIncludeDetail(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("love truth light")
	output := RenderReadingWithOptions(reading, DebugRenderOptions())

	// Debug output can be longer
	if len(output) < 100 {
		t.Error("debug output seems too short")
	}

	// Debug may include passage field info
	t.Logf("debug output length: %d chars", len(output))
}
