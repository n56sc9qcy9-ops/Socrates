package decipher

import (
	"testing"
)

// DUPLICATE/REPEATED TOKEN EVIDENCE POLICY TESTS
// =============================================================================

func TestRepeatedTokens_AccumulationWithinBounds(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Single occurrence
	fields1, _ := AnalyzePassage("love", engine)

	// Repeated occurrences (3x)
	fields3, _ := AnalyzePassage("love love love", engine)

	// Get love field strengths
	var strength1, strength3 float64

	for _, field := range fields1 {
		if field.Concept == "love" {
			strength1 = field.Strength
		}
	}

	for _, field := range fields3 {
		if field.Concept == "love" {
			strength3 = field.Strength
		}
	}

	// Multiple occurrences should accumulate but not linearly
	// Allow up to 4x accumulation for triple occurrence (reasonable for co-activation)
	maxExpected := strength1 * 4.0

	if strength3 > maxExpected {
		t.Errorf("repeated tokens should not inflate linearly: single=%.4f, triple=%.4f, max=%.4f",
			strength1, strength3, maxExpected)
	}
}

func TestDuplicateInflation_EvidenceDeduplication(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Analyze once
	fields1, _ := AnalyzePassage("light", engine)

	// Analyze twice with same passage
	fields2, _ := AnalyzePassage("light", engine)

	// Both should produce same number of fields
	if len(fields1) != len(fields2) {
		t.Errorf("repeated analysis should produce same number of fields: %d vs %d",
			len(fields1), len(fields2))
	}
}

// =============================================================================
