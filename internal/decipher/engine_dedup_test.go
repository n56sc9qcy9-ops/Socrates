package decipher

import (
	"testing"
)


// DUPLICATE EVIDENCE DEDUPLICATION TESTS
// =============================================================================

func TestDuplicateMatchEvidenceDoesNotInflateScore(t *testing.T) {
	// Test that repeated MatchEvidence with same path does not inflate score
	// Uses scoring.go deduplication

	matches := []MatchEvidence{
		{InputForm: "skal", AnchorForm: "skal", Method: "exact", Distance: 0, Weight: 0.5},
		{InputForm: "skal", AnchorForm: "skal", Method: "exact", Distance: 0, Weight: 0.5}, // duplicate
		{InputForm: "skal", AnchorForm: "skal", Method: "exact", Distance: 0, Weight: 0.5}, // duplicate
	}

	deduped := deduplicateMatchEvidence(matches)

	if len(deduped) != 1 {
		t.Errorf("expected 1 deduped match, got %d", len(deduped))
	}

	// Verify score calculation uses deduplicated matches
	components := CalculateScoreComponents(nil, deduped, nil, ConvergenceResult{}, []ChannelResult{})
	// With 1 exact match at weight 0.5, exact match score = 0.5
	if components.ExactMatchScore != 0.5 {
		t.Errorf("expected ExactMatchScore 0.5, got %f", components.ExactMatchScore)
	}
}

func TestDuplicateChannelSignalsDoNotInflateStrength(t *testing.T) {
	// Test that signals with same channel+target are deduplicated before scoring

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Analyze produces signals that get deduplicated in buildSignalGraph
	reading := engine.Analyze("inspired")

	// If we have converging patterns, verify they don't have inflated signal counts
	// from duplicate evidence
	for _, pattern := range reading.ConvergingPatterns {
		// Count unique channel+target combinations
		seenKeys := make(map[string]bool)
		for _, sig := range pattern.Signals {
			key := sig.Channel + "|" + sig.Target
			seenKeys[key] = true
		}

		// The number of signals should not exceed the number of unique keys
		// (duplicates should be collapsed)
		if len(pattern.Signals) > len(seenKeys) {
			t.Errorf("pattern %s has %d signals but only %d unique keys - duplicates not collapsed",
				pattern.Name, len(pattern.Signals), len(seenKeys))
		}
	}
}

func TestExactVerifiedEvidenceOutranksRepeatedSpeculativeEvidence(t *testing.T) {
	// Test that exact/verified evidence scores above repeated speculative evidence

	// Exact verified match
	exactVerified := []MatchEvidence{
		{InputForm: "truth", AnchorForm: "truth", Method: "exact", Distance: 0, Weight: 0.8},
	}

	// Repeated speculative fuzzy matches (should not outrank exact)
	repeatedSpeculative := []MatchEvidence{
		{InputForm: "truth", AnchorForm: "trth", Method: "consonant_skeleton", Distance: 0.5, Weight: 0.3},
		{InputForm: "truth", AnchorForm: "trth", Method: "consonant_skeleton", Distance: 0.5, Weight: 0.3}, // duplicate
		{InputForm: "truth", AnchorForm: "trth", Method: "consonant_skeleton", Distance: 0.5, Weight: 0.3}, // duplicate
	}

	dedupedSpeculative := deduplicateMatchEvidence(repeatedSpeculative)

	exactComponents := CalculateScoreComponents(nil, exactVerified, nil, ConvergenceResult{}, []ChannelResult{})
	speculativeComponents := CalculateScoreComponents(nil, dedupedSpeculative, nil, ConvergenceResult{}, []ChannelResult{})

	// Exact match score should be higher than fuzzy match score
	if speculativeComponents.FuzzyMatchScore >= exactComponents.ExactMatchScore {
		t.Errorf("exact match score (%.2f) should exceed single fuzzy match score (%.2f)",
			exactComponents.ExactMatchScore, speculativeComponents.FuzzyMatchScore)
	}
}

func TestDeduplicateSignalsFunction(t *testing.T) {
	// Test the DeduplicateSignals utility function
	signals := []Signal{
		{Text: "skal", Target: "truth", Channel: "glyph", Confidence: "verified", Weight: 0.5},
		{Text: "skal", Target: "truth", Channel: "glyph", Confidence: "verified", Weight: 0.5},  // duplicate
		{Text: "skal", Target: "truth", Channel: "glyph", Confidence: "verified", Weight: 0.5},  // duplicate
		{Text: "skal", Target: "truth", Channel: "sound", Confidence: "plausible", Weight: 0.4}, // different channel - independent evidence
	}

	deduped := DeduplicateSignals(signals)

	// Should have 2 - one for glyph, one for sound (different channels)
	if len(deduped) != 2 {
		t.Errorf("expected 2 signals after dedup (different channels), got %d", len(deduped))
	}

	// Verify the evidence IDs
	if deduped[0].EvidenceID() != "glyph|truth|skal" {
		t.Errorf("unexpected evidence ID for first signal: %s", deduped[0].EvidenceID())
	}
}

func TestDeduplicatePassageSignalsFunction(t *testing.T) {
	// Test the DeduplicatePassageSignals utility function
	signals := []PassageSignal{
		{Token: "in", Concept: "truth", MatchForm: "skal", Weight: 0.5},
		{Token: "in", Concept: "truth", MatchForm: "skal", Weight: 0.5},     // duplicate
		{Token: "in", Concept: "truth", MatchForm: "skal", Weight: 0.5},     // duplicate
		{Token: "spirit", Concept: "truth", MatchForm: "skal", Weight: 0.4}, // different token - independent evidence
	}

	deduped := DeduplicatePassageSignals(signals)

	// Should have 2 - one for 'in', one for 'spirit' (different tokens)
	if len(deduped) != 2 {
		t.Errorf("expected 2 signals after dedup (different tokens), got %d", len(deduped))
	}
}

func TestActivationStrengthNotInflatedByDuplicates(t *testing.T) {
	// Test that ComputeActivatedConcepts doesn't inflate strength with duplicates

	signals := []PassageSignal{
		{Token: "in", Concept: "truth", MatchForm: "skal", Weight: 0.5},
		{Token: "in", Concept: "truth", MatchForm: "skal", Weight: 0.5}, // duplicate - should not add weight
		{Token: "in", Concept: "truth", MatchForm: "skal", Weight: 0.5}, // duplicate - should not add weight
	}

	// With deduplication, the concept should only get 0.5 strength, not 1.5
	activated := ComputeActivatedConcepts(signals, nil, nil)

	if len(activated) != 1 {
		t.Errorf("expected 1 activated concept, got %d", len(activated))
	}

	if activated[0].Strength > 0.6 { // Allow small margin for floating point
		t.Errorf("concept strength should not be inflated by duplicates, got %f (expected ~0.5)", activated[0].Strength)
	}
}

// =============================================================================
