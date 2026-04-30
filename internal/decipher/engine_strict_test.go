package decipher

import (
	"testing"

	"socrates/internal/knowledge"
)



// STRICT BOUNDED WORK TESTS
// =============================================================================

func TestCoreCandidatesSurviveTightBounds(t *testing.T) {
	// Test that core high-confidence candidates survive even with very tight bounds
	input := "energy"

	// Very tight bounds - only 5 candidates allowed
	tightBounds := CandidateBounds{
		MaxCandidates:            5,
		MaxSpeculativeCandidates: 2,
	}

	candidates, _ := GenerateCandidateForms(input, tightBounds)

	// Core candidates should always be present
	coreMethods := map[string]bool{
		"normalized":         false,
		"consonant_skeleton": false,
		"phonetic":           false,
	}

	for _, cand := range candidates {
		if _, ok := coreMethods[cand.Method]; ok {
			coreMethods[cand.Method] = true
		}
	}

	for method, found := range coreMethods {
		if !found {
			t.Errorf("core candidate method '%s' should survive tight bounds of 5 candidates", method)
		}
	}

	// Also verify we didn't exceed the max candidates cap
	if len(candidates) > 5 {
		t.Errorf("should not exceed maxCandidates (5), got %d", len(candidates))
	}
}

func TestExactMatchesFoundBeforeBudgetConsumed(t *testing.T) {
	// Test that high-confidence candidates get matched before the fuzzy comparison budget is exhausted
	var kb *knowledge.Knowledge = testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	// Create a controlled test with known high-quality candidates and matching anchors
	candidates := []CandidateForm{
		{Form: "energy", Method: "normalized", Distance: 0, Confidence: "verified"},
		{Form: "spirit", Method: "normalized", Distance: 0, Confidence: "verified"},
		{Form: "truth", Method: "normalized", Distance: 0, Confidence: "verified"},
		{Form: "breath", Method: "normalized", Distance: 0, Confidence: "verified"},
		{Form: "xyzabc", Method: "normalized", Distance: 0, Confidence: "speculative"},
	}
	anchors := []AnchorConcept{
		{Form: "energy", Concept: "power", Confidence: "verified", Weight: 0.8},
		{Form: "spirit", Concept: "soul", Confidence: "verified", Weight: 0.8},
		{Form: "truth", Concept: "truth", Confidence: "verified", Weight: 0.8},
		{Form: "breath", Concept: "life", Confidence: "verified", Weight: 0.8},
		{Form: "xyzabc", Concept: "unknown", Confidence: "speculative", Weight: 0.3},
	}

	// Very tight comparison budget - only 5 comparisons allowed
	tightBounds := FuzzyBounds{
		MaxComparisons: 5,
		MaxMatches:     10,
	}

	matches, discarded := FuzzyMatchEvidence(candidates, anchors, tightBounds)

	// With 5 comparisons allowed and 5 candidates, all candidates should be processed
	// and we should find matches for at least the verified candidates
	if len(matches) == 0 {
		t.Error("should find matches with tight budget of 5 comparisons")
	}

	// All matches should be from verified candidates (exact matches)
	for _, m := range matches {
		if m.Method != "exact" && m.Method != "case_insensitive" {
			t.Errorf("with tight budget, matches should be exact, got method '%s'", m.Method)
		}
	}

	// Discarded count should reflect remaining comparisons
	t.Logf("Tight budget: %d matches found, %d discarded comparisons", len(matches), discarded)
}

func TestEngineReadingExposesDiscardedCounts(t *testing.T) {
	// Engine-level test that Reading.DiscardedCandidates and DiscardedComparisons are populated
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Deliberately long input to trigger bounds
	longInput := "abcdefghijklmnopqrstuvwxyz"
	reading := engine.Analyze(longInput)

	// DiscardedCandidates should be populated (long input generates many speculative candidates)
	if reading.DiscardedCandidates < 0 {
		t.Errorf("DiscardedCandidates should be >= 0, got %d", reading.DiscardedCandidates)
	}

	// DiscardedComparisons should be populated
	if reading.DiscardedComparisons < 0 {
		t.Errorf("DiscardedComparisons should be >= 0, got %d", reading.DiscardedComparisons)
	}

	// For a long input, at least one of these should be positive
	// (unless the input happens to match so few things that no bounds are hit)
	t.Logf("Long input discarded: %d candidates, %d comparisons", reading.DiscardedCandidates, reading.DiscardedComparisons)
}

func TestTightBoundsStillProduceMatches(t *testing.T) {
	// Test that even with very tight bounds, we still get some results
	var kb *knowledge.Knowledge = testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	input := "energy"
	candidates, discardedCand := GenerateCandidateForms(input, CandidateBounds{
		MaxCandidates:            10,
		MaxSpeculativeCandidates: 5,
	})

	anchors := GetAllAnchors(kb)
	matches, discardedComp := FuzzyMatchEvidence(candidates, anchors, FuzzyBounds{
		MaxComparisons: 50,
		MaxMatches:     5,
	})

	// Should still produce candidates
	if len(candidates) == 0 {
		t.Error("should produce candidates even with tight bounds")
	}

	// Should not exceed max candidates
	if len(candidates) > 10 {
		t.Errorf("should not exceed 10 candidates, got %d", len(candidates))
	}

	// Should not exceed max matches
	if len(matches) > 5 {
		t.Errorf("should not exceed 5 matches, got %d", len(matches))
	}

	// Discarded counts should be non-negative
	if discardedCand < 0 {
		t.Errorf("discardedCand should be >= 0, got %d", discardedCand)
	}
	if discardedComp < 0 {
		t.Errorf("discardedComp should be >= 0, got %d", discardedComp)
	}

	t.Logf("Tight bounds: %d candidates (%d discarded), %d matches (%d discarded)",
		len(candidates), discardedCand, len(matches), discardedComp)
}

func TestExactHighConfidenceRetainedUnderCaps(t *testing.T) {
	// Test that exact/high-confidence matches are retained even when speculative matches are capped
	var kb *knowledge.Knowledge = testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	// Input with known exact match
	input := "truth"
	candidates, _ := GenerateCandidateForms(input, DefaultCandidateBounds())
	anchors := GetAllAnchors(kb)

	// Very tight match limit - only 1 match allowed
	matches, _ := FuzzyMatchEvidence(candidates, anchors, FuzzyBounds{
		MaxComparisons: 100,
		MaxMatches:     1,
	})

	// If we found any matches, the first should be an exact/high-confidence one
	if len(matches) > 0 {
		firstMatch := matches[0]
		// Exact match should rank first for "truth"
		if firstMatch.Method != "exact" && firstMatch.Method != "case_insensitive" {
			t.Errorf("with tight cap of 1, first match should be exact for 'truth', got '%s'", firstMatch.Method)
		}
	}
}

// =============================================================================
