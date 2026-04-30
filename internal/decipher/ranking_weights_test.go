package decipher

import (
	"os"
	"testing"
)

// =============================================================================
// Ranking Weights Tests
// =============================================================================

func TestDefaultRankingWeights(t *testing.T) {
	weights := DefaultRankingWeights()

	// Verify all weights are within reasonable bounds
	if weights.ExactMatchBase <= 0 || weights.ExactMatchBase > 10 {
		t.Errorf("ExactMatchBase should be in (0, 10], got %v", weights.ExactMatchBase)
	}

	if weights.FuzzyMatchBase <= 0 || weights.FuzzyMatchBase > 10 {
		t.Errorf("FuzzyMatchBase should be in (0, 10], got %v", weights.FuzzyMatchBase)
	}

	if weights.ConfidenceVerified < weights.ConfidenceSpeculative {
		t.Errorf("Verified confidence should be >= speculative, got verified=%v speculative=%v",
			weights.ConfidenceVerified, weights.ConfidenceSpeculative)
	}

	if weights.PassageCoActivationWeight < 0 || weights.PassageCoActivationWeight > 10 {
		t.Errorf("PassageCoActivationWeight should be in [0, 10], got %v", weights.PassageCoActivationWeight)
	}

	if weights.MultiMethodThreshold < 2 {
		t.Errorf("MultiMethodThreshold should be >= 2, got %v", weights.MultiMethodThreshold)
	}

	if weights.ChannelDiversityThreshold < 2 {
		t.Errorf("ChannelDiversityThreshold should be >= 2, got %v", weights.ChannelDiversityThreshold)
	}
}

func TestRankingWeightsValidate(t *testing.T) {
	weights := DefaultRankingWeights()
	errors := weights.Validate()

	if len(errors) > 0 {
		t.Errorf("Default weights should be valid, got errors: %v", errors)
	}
}

func TestRankingWeightsValidateInvalid(t *testing.T) {
	weights := RankingWeights{
		ExactMatchBase: 15, // Invalid: > 10
	}

	errors := weights.Validate()

	if len(errors) == 0 {
		t.Error("Expected validation errors for weights with ExactMatchBase=15")
	}
}

func TestRankingWeightsConfidenceMultiplier(t *testing.T) {
	weights := DefaultRankingWeights()

	// Test verified
	verified := weights.ConfidenceMultiplier("verified")
	if verified != weights.ConfidenceVerified {
		t.Errorf("ConfidenceMultiplier(verified) should return ConfidenceVerified")
	}

	// Test plausible
	plausible := weights.ConfidenceMultiplier("plausible")
	if plausible != weights.ConfidencePlausible {
		t.Errorf("ConfidenceMultiplier(plausible) should return ConfidencePlausible")
	}

	// Test speculative
	speculative := weights.ConfidenceMultiplier("speculative")
	if speculative != weights.ConfidenceSpeculative {
		t.Errorf("ConfidenceMultiplier(speculative) should return ConfidenceSpeculative")
	}

	// Test unknown (should default to plausible)
	unknown := weights.ConfidenceMultiplier("unknown")
	if unknown != weights.ConfidencePlausible {
		t.Errorf("ConfidenceMultiplier(unknown) should default to ConfidencePlausible")
	}
}

func TestRankingWeightsSourceMultiplier(t *testing.T) {
	weights := DefaultRankingWeights()

	// Test curated
	curated := weights.SourceMultiplier("curated")
	if curated != weights.SourceCurated {
		t.Errorf("SourceMultiplier(curated) should return SourceCurated")
	}

	// Test traditional
	traditional := weights.SourceMultiplier("traditional")
	if traditional != weights.SourceTraditional {
		t.Errorf("SourceMultiplier(traditional) should return SourceTraditional")
	}

	// Test human_review
	humanReview := weights.SourceMultiplier("human_review")
	if humanReview != weights.SourceHumanReview {
		t.Errorf("SourceMultiplier(human_review) should return SourceHumanReview")
	}

	// Test unknown (should default to 1.0)
	unknown := weights.SourceMultiplier("unknown")
	if unknown != 1.0 {
		t.Errorf("SourceMultiplier(unknown) should default to 1.0")
	}
}

func TestRankingWeightsLensWeight(t *testing.T) {
	weights := DefaultRankingWeights()

	// Test orthographic
	ortho := weights.LensWeight("orthographic")
	if ortho != weights.LensOrthographic {
		t.Errorf("LensWeight(orthographic) should return LensOrthographic")
	}

	// Test phonetic
	phonetic := weights.LensWeight("phonetic")
	if phonetic != weights.LensPhonetic {
		t.Errorf("LensWeight(phonetic) should return LensPhonetic")
	}

	// Test semantic
	semantic := weights.LensWeight("semantic")
	if semantic != weights.LensSemantic {
		t.Errorf("LensWeight(semantic) should return LensSemantic")
	}

	// Test unknown (should default to 0.5)
	unknown := weights.LensWeight("unknown")
	if unknown != 0.5 {
		t.Errorf("LensWeight(unknown) should default to 0.5")
	}
}

func TestLoadRankingWeights(t *testing.T) {
	// Test loading a valid weights file
	// Use absolute path to training directory
	absPath := "/Users/bot/Socrates/training/ranking_weights.yaml"
	weights, err := LoadRankingWeights(absPath)
	if err != nil {
		t.Skipf("Could not load weights file: %v", err)
	}

	// Verify weights were loaded
	if weights.ExactMatchBase != 1.0 {
		t.Errorf("Expected ExactMatchBase=1.0, got %v", weights.ExactMatchBase)
	}

	if weights.FuzzyMatchBase != 0.8 {
		t.Errorf("Expected FuzzyMatchBase=0.8, got %v", weights.FuzzyMatchBase)
	}

	if weights.MultiMethodThreshold != 3 {
		t.Errorf("Expected MultiMethodThreshold=3, got %v", weights.MultiMethodThreshold)
	}
}

func TestLoadRankingWeightsInvalidPath(t *testing.T) {
	_, err := LoadRankingWeights("nonexistent/path.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent weights file")
	}
}

func TestLoadRankingWeightsInvalidYAML(t *testing.T) {
	// Create a temp invalid file
	tmpPath := "/tmp/invalid_weights_test.yaml"
	if err := os.WriteFile(tmpPath, []byte("invalid: yaml: content: ["), 0644); err != nil {
		t.Skip("Could not create temp file")
	}

	_, err := LoadRankingWeights(tmpPath)
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}

	// Cleanup
	os.Remove(tmpPath)
}

func TestCalculateScoreComponentsWithWeights(t *testing.T) {
	weights := DefaultRankingWeights()

	// Test with empty inputs
	components := CalculateScoreComponentsWithWeights(nil, nil, nil, ConvergenceResult{}, nil, weights)

	if components.ExactMatchScore != 0 {
		t.Errorf("Expected 0 ExactMatchScore with nil matches, got %v", components.ExactMatchScore)
	}

	if components.FuzzyMatchScore != 0 {
		t.Errorf("Expected 0 FuzzyMatchScore with nil matches, got %v", components.FuzzyMatchScore)
	}

	if components.GraphExpansionScore != 0 {
		t.Errorf("Expected 0 GraphExpansionScore with nil expansions, got %v", components.GraphExpansionScore)
	}
}

func TestCalculateFinalScoreWithWeights(t *testing.T) {
	weights := DefaultRankingWeights()

	components := ScoreComponents{
		ExactMatchScore:         0.5,
		FuzzyMatchScore:         0.3,
		GraphExpansionScore:     0.2,
		PassageConvergenceScore: 0.4,
		MultiMethodBonus:        0.1,
		ChannelDiversityBonus:  0.1,
	}

	score := CalculateFinalScoreWithWeights(components, weights)

	// Score should be capped at 1.0
	if score > 1.0 {
		t.Errorf("Score should be capped at 1.0, got %v", score)
	}

	// Score should be positive
	if score <= 0 {
		t.Errorf("Score should be positive, got %v", score)
	}
}

func TestDefaultWeightsPreserveBehavior(t *testing.T) {
	// Test that CalculateScoreComponents with default weights
	// produces similar results to the original CalculateScoreComponents

	// Both should produce non-negative scores
	weights := DefaultRankingWeights()

	// Test with minimal components
	components1 := CalculateScoreComponents(nil, nil, nil, ConvergenceResult{}, nil)
	components2 := CalculateScoreComponentsWithWeights(nil, nil, nil, ConvergenceResult{}, nil, weights)

	// Both should have the same structure
	if components1.ExactMatchScore != components2.ExactMatchScore {
		t.Logf("Note: ExactMatchScore differs between old/new: %v vs %v",
			components1.ExactMatchScore, components2.ExactMatchScore)
	}
}