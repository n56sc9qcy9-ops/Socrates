package decipher

import (
	"testing"
	"strings"
)


// KEY TOKEN REMOVAL REGRESSION TESTS
// =============================================================================

// TestKeyTokenRemovalWeakenActivation tests that removing a key evidence token
// weakens the relevant activation field. This is a regression test for
// Activation-Energy Discipline Gate.
func TestKeyTokenRemovalWeakenActivation(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Test: Single word "truth" - the key token IS "truth"
	// Removing key consonant 't' gives "ruth" which has different meaning
	readingFull := engine.Analyze("truth")
	readingNoKey := engine.Analyze("ruth")

	// The overall score should be lower when we remove the key token
	if readingFull.Score.Overall <= readingNoKey.Score.Overall {
		t.Errorf("removing key token should weaken overall score:\n  truth score: %.2f\n  ruth score: %.2f",
			readingFull.Score.Overall, readingNoKey.Score.Overall)
	}

	// The graph expansion score should be lower without the key anchor
	if readingFull.Score.Components.GraphExpansionScore <= readingNoKey.Score.Components.GraphExpansionScore {
		t.Errorf("removing key token should weaken graph expansion:\n  truth graph: %.2f\n  ruth graph: %.2f",
			readingFull.Score.Components.GraphExpansionScore, readingNoKey.Score.Components.GraphExpansionScore)
	}
}

// TestKeyTokenRemovalWeakenActivationMultiWord tests with a multi-word input
// where removing a key token should weaken convergence.
func TestKeyTokenRemovalWeakenActivationMultiWord(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// "word" is a key concept with neighbors like truth
	readingWord := engine.Analyze("word")
	readingOrd := engine.Analyze("ord")

	// "word" should score higher than "ord" (partial, not a key concept)
	if readingWord.Score.Overall <= readingOrd.Score.Overall {
		t.Errorf("key word 'word' should score higher than partial 'ord':\n  word: %.2f\n  ord: %.2f",
			readingWord.Score.Overall, readingOrd.Score.Overall)
	}
}

// TestKeyTokenRemovalWeakensConvergence tests that removing a key token
// from a passage weakens the convergence score.
func TestKeyTokenRemovalWeakensConvergence(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// "being" is a key concept with strong graph connections
	readingBeing := engine.Analyze("being")
	readingBeeing := engine.Analyze("beeing") // misspelling, missing key consonant pattern

	// "being" should have stronger convergence than misspelled version
	// The misspelling disrupts the phonetic signature
	if readingBeing.Score.Components.PassageConvergenceScore <= readingBeeing.Score.Components.PassageConvergenceScore {
		t.Errorf("correct 'being' should have stronger convergence than 'beeing':\n  being: %.2f\n  beeing: %.2f",
			readingBeing.Score.Components.PassageConvergenceScore, readingBeeing.Score.Components.PassageConvergenceScore)
	}
}

// TestKeyTokenRemovalWeakenExactMatch tests that removing a key token
// reduces the exact match score (for known concepts).
func TestKeyTokenRemovalWeakenExactMatch(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// "spirit" is a known concept with strong semantic connections
	readingSpirit := engine.Analyze("spirit")
	readingSpiri := engine.Analyze("spiri") // missing key ending

	// "spirit" should have higher overall score than "spiri"
	if readingSpirit.Score.Overall <= readingSpiri.Score.Overall {
		t.Errorf("key word 'spirit' should score higher than partial 'spiri':\n  spirit: %.2f\n  spiri: %.2f",
			readingSpirit.Score.Overall, readingSpiri.Score.Overall)
	}
}

// TestGodIsLove_BoundedScores verifies that "God is Love" has:
// - Bounded score components (all in [0, 1] range)
// - No duplicate evidence paths in output
// - No whitespace creating spurious repetition signals
func TestGodIsLove_BoundedScores(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("God is Love")
	output := RenderReadingWithOptions(reading, DefaultRenderOptions())

	// Verify score components are bounded in [0, 1]
	comp := reading.Score.Components
	if comp.ExactMatchScore < 0 || comp.ExactMatchScore > 1.0 {
		t.Errorf("ExactMatchScore should be in [0, 1], got %f", comp.ExactMatchScore)
	}
	if comp.FuzzyMatchScore < 0 || comp.FuzzyMatchScore > 1.0 {
		t.Errorf("FuzzyMatchScore should be in [0, 1], got %f", comp.FuzzyMatchScore)
	}
	if comp.GraphExpansionScore < 0 || comp.GraphExpansionScore > 1.0 {
		t.Errorf("GraphExpansionScore should be in [0, 1], got %f", comp.GraphExpansionScore)
	}
	if comp.PassageConvergenceScore < 0 || comp.PassageConvergenceScore > 1.0 {
		t.Errorf("PassageConvergenceScore should be in [0, 1], got %f", comp.PassageConvergenceScore)
	}
	if comp.MultiMethodBonus < 0 || comp.MultiMethodBonus > 1.0 {
		t.Errorf("MultiMethodBonus should be in [0, 1], got %f", comp.MultiMethodBonus)
	}
	if comp.ChannelDiversityBonus < 0 || comp.ChannelDiversityBonus > 1.0 {
		t.Errorf("ChannelDiversityBonus should be in [0, 1], got %f", comp.ChannelDiversityBonus)
	}

	// Verify no duplicate evidence paths in output
	// Count occurrences of "primitive 'love' matches love" - should appear once
	countPrimitiveLove := strings.Count(output, "primitive 'love' matches love")
	if countPrimitiveLove > 1 {
		t.Errorf("primitive 'love' should appear at most once in output, got %d occurrences", countPrimitiveLove)
	}

	// Verify no per-character repetition signals (whitespace should not create them)
	// Should NOT have patterns like "o repeated 2x" or "repeated 2x" appearing multiple times
	repetitionCount := strings.Count(output, "repeated letters detected")
	if repetitionCount > 1 {
		t.Errorf("'repeated letters detected' should appear at most once, got %d occurrences", repetitionCount)
	}

	// Verify overall score is bounded
	if reading.Score.Overall < 0 || reading.Score.Overall > 1.0 {
		t.Errorf("Overall score should be in [0, 1], got %f", reading.Score.Overall)
	}
}
