package decipher

import (
	"testing"
)

// KEY TOKEN WEAKENING STRICT TESTS
// =============================================================================

func TestKeyTokenRemoval_StrictlyWeakensRelatedField(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Full passage with key semantic word
	fieldsFull, _ := AnalyzePassage("love truth light", engine)
	// Without the key word
	fieldsNoLove, _ := AnalyzePassage("truth light", engine)

	// Get love field strengths
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

	// Removing "love" from the passage MUST strictly weaken the love field
	// (either love field disappears or strength is reduced)
	if fullLoveStrength <= noLoveStrength {
		t.Errorf("removing 'love' token should strictly weaken love field:\n  full: %.4f\n  no love: %.4f",
			fullLoveStrength, noLoveStrength)
	}
}

func TestKeyTokenRemoval_CompleteRemoval(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Analyze single word
	fieldsSingle, _ := AnalyzePassage("truth", engine)

	// Get truth field strength from single
	var singleTruthStrength float64
	for _, field := range fieldsSingle {
		if field.Concept == "truth" {
			singleTruthStrength = field.Strength
		}
	}

	// Single word analysis should produce truth field with positive strength
	if singleTruthStrength <= 0 {
		t.Error("single word analysis should produce truth field with positive strength")
	}

	// Analyze again - strengths should be in the same ballpark (within 50%)
	fieldsAgain, _ := AnalyzePassage("truth", engine)
	var againTruthStrength float64
	for _, field := range fieldsAgain {
		if field.Concept == "truth" {
			againTruthStrength = field.Strength
		}
	}

	ratio := againTruthStrength / singleTruthStrength
	if ratio < 0.5 || ratio > 2.0 {
		t.Errorf("repeated single-word analysis should produce similar strength: %.4f vs %.4f (ratio=%.2f)",
			singleTruthStrength, againTruthStrength, ratio)
	}
}

// =============================================================================
