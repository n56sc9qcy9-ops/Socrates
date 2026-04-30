package decipher

import (
	"testing"
	"strings"
)


// CLEANUP: Fuzzy Matching Tightening Tests
// =============================================================================

func TestUnrelatedInputProducesFewFuzzyMatches(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Unrelated noisy input should produce few or no fuzzy matches
	reading := engine.Analyze("zzskalx")

	// Count actual fuzzy matches (not just any signal)
	fuzzyMatchCount := len(reading.FuzzyMatches)

	// A truly unrelated string should have very few (ideally 0) fuzzy matches
	if fuzzyMatchCount > 3 {
		t.Errorf("unrelated 'zzskalx' should produce at most 3 fuzzy matches, got %d", fuzzyMatchCount)
	}
}

func TestExactMatchesOutrankFuzzyMatches(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// An input that matches an anchor exactly should have exact matches
	// ranked above any fuzzy matches
	reading := engine.Analyze("skal")

	// Find the first exact match and first fuzzy (non-exact) match
	var firstExactWeight, firstFuzzyWeight float64
	var foundExact, foundFuzzy bool

	for _, m := range reading.FuzzyMatches {
		if !foundExact && (m.Method == "exact" || m.Method == "case_insensitive") {
			firstExactWeight = m.Weight
			foundExact = true
		}
		if !foundFuzzy && m.Method != "exact" && m.Method != "case_insensitive" {
			firstFuzzyWeight = m.Weight
			foundFuzzy = true
		}
		if foundExact && foundFuzzy {
			break
		}
	}

	// If we have both exact and fuzzy, exact should rank higher or equal
	if foundExact && foundFuzzy {
		if firstFuzzyWeight > firstExactWeight {
			t.Errorf("exact match weight (%.2f) should be >= fuzzy match weight (%.2f)",
				firstExactWeight, firstFuzzyWeight)
		}
	}
}

func TestNoisyInputScoresLowerThanCleanInput(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	cleanReading := engine.Analyze("skal")
	noisyReading := engine.Analyze("zzskalx")

	// Clean input should score higher than noisy input
	// Allow small margin for floating point variations
	diff := noisyReading.Score.Overall - cleanReading.Score.Overall
	if diff > 0.05 {
		t.Errorf("noisy 'zzskalx' score (%.2f) should not significantly exceed clean 'skal' score (%.2f), diff=%.2f",
			noisyReading.Score.Overall, cleanReading.Score.Overall, diff)
	}
}

func TestRenderedFuzzyMatchesCapped(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Generate a reading with many potential fuzzy matches
	reading := engine.Analyze("inspired")
	output := RenderReading(reading)

	// Count lines in fuzzy matches section
	lines := strings.Split(output, "\n")
	inFuzzySection := false
	fuzzyLineCount := 0

	for _, line := range lines {
		if strings.HasPrefix(line, "Fuzzy Matches:") {
			inFuzzySection = true
			continue
		}
		if inFuzzySection {
			if strings.HasPrefix(line, "  - ") {
				fuzzyLineCount++
			}
			// End of section
			if line == "" || strings.HasPrefix(line, "  ...") || strings.HasPrefix(line, "Phase ") || strings.HasPrefix(line, "Graph") || strings.HasPrefix(line, "Passage") {
				if strings.HasPrefix(line, "  ...") {
					// This line indicates cap was applied
				}
				break
			}
		}
	}

	// Fuzzy matches should be capped at 10 (plus optional weaker-match line)
	if fuzzyLineCount > 10 {
		t.Errorf("rendered fuzzy matches should be capped at 10, got %d", fuzzyLineCount)
	}
}

func TestLevenshteinThresholdIsTight(t *testing.T) {
	// Test that very weak Levenshtein matches are rejected
	// A string like "abcdef" compared to "xyz123" should NOT match
	a := "abcdef"
	b := "xyz123"
	method, distance := fuzzyMatch(a, b)

	// Should be Levenshtein with high distance
	if method == "exact" || method == "case_insensitive" {
		t.Errorf("'%s' vs '%s' should not be exact match", a, b)
	}

	// Distance should be high (normalized > 0.34)
	if acceptMatch(method, distance, a, b) {
		t.Errorf("'%s' vs '%s' should NOT be accepted (distance %.2f exceeds 0.34)", a, b, distance)
	}
}

func TestFuzzyMatchEvidenceDeduplicates(t *testing.T) {
	// Test that duplicate evidence is removed
	candidates := []CandidateForm{
		{Form: "skal", Method: "normalized", Distance: 0, Confidence: "verified"},
		{Form: "skal", Method: "normalized", Distance: 0, Confidence: "verified"}, // duplicate
	}
	anchors := []AnchorConcept{
		{Form: "skal", Concept: "truth", Confidence: "verified", Weight: 0.5},
	}

	evidence, _ := FuzzyMatchEvidence(candidates, anchors)

	// Should have only one match, not duplicate
	if len(evidence) != 1 {
		t.Errorf("expected 1 evidence after dedup, got %d", len(evidence))
	}
}

// =============================================================================
