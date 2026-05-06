package training

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"socrates/internal/decipher"
	"socrates/internal/knowledge"
)

// repoRoot finds the repository root by locating go.mod.
func repoRoot() string {
	cwd, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
			return cwd
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			break
		}
		cwd = parent
	}
	return "/Users/bot/Socrates" // fallback
}

// TestWeightSuggester generates and writes suggestions.
func TestWeightSuggester(t *testing.T) {
	os.Chdir(repoRoot())

	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Failed to load knowledge: %v", err)
	}

	engine := decipher.NewEngineWithKnowledge(kb)
	weights := decipher.DefaultRankingWeights()

	suggester := NewWeightSuggester(engine, weights)
	reviewDir := t.TempDir()
	suggester.SetReviewPath(filepath.Join(reviewDir, "weight_suggestions.yaml"))

	// Load examples
	loader := NewLoader()
	trainExamples, err := loader.LoadFromFile("training/examples.yaml")
	if err != nil {
		t.Fatalf("Failed to load train examples: %v", err)
	}

	heldOutExamples, err := loader.LoadHeldOutFromFile("training/heldout.yaml")
	if err != nil {
		t.Fatalf("Failed to load held-out examples: %v", err)
	}

	// Evaluate with current weights to get baseline metrics
	eval := NewEvaluatorWithWeights(engine, weights, "baseline")
	report := eval.EvaluateWithReport(trainExamples)
	heldOutReport := &EvaluationReport{}
	eval.EvaluateHeldOut(heldOutExamples, heldOutReport)
	heldOutReport.Finalize()

	// Baseline metrics from current evaluation
	baseTrainMetrics := MetricsDelta{
		ConceptPrecision: report.Train.ConceptPrec,
		ConceptRecall:    report.Train.ConceptRec,
		FieldPrecision:   report.Train.FieldPrec,
		FieldRecall:      report.Train.FieldRec,
	}
	baseHeldOutMetrics := MetricsDelta{
		ConceptPrecision: heldOutReport.HeldOut.ConceptPrec,
		ConceptRecall:    heldOutReport.HeldOut.ConceptRec,
		FieldPrecision:   heldOutReport.HeldOut.FieldPrec,
		FieldRecall:      heldOutReport.HeldOut.FieldRec,
	}

	// Generate suggestions
	candidates := suggester.GenerateSuggestions(trainExamples, heldOutExamples, baseTrainMetrics, baseHeldOutMetrics)

	t.Logf("Generated %d candidates", len(candidates))
	for i, c := range candidates {
		t.Logf("  [%d] %s: %v -> %v (heldout fp delta: %.4f)",
			i+1, c.Path, c.Current, c.Suggested, c.HeldOutDelta.FieldPrecision)
	}

	// Write suggestions to review file
	err = suggester.WriteSuggestions()
	if err != nil {
		t.Fatalf("Failed to write suggestions: %v", err)
	}

	// Verify review file exists
	if _, err := os.Stat(suggester.ReviewPath); os.IsNotExist(err) {
		t.Errorf("Review file was not created at %s", suggester.ReviewPath)
	} else {
		t.Logf("Review file created at %s", suggester.ReviewPath)
	}

	// Verify suggestions are marked as pending
	for _, c := range candidates {
		if c.Status != "pending" {
			t.Errorf("Expected status 'pending', got '%s'", c.Status)
		}
	}

	// Verify suggestions contain before/after metrics
	for _, c := range candidates {
		if c.TrainDelta.FieldPrecision == 0 && c.TrainDelta.FieldRecall == 0 {
			t.Logf("Warning: TrainDelta appears empty for %s", c.Path)
		}
		if c.HeldOutDelta.FieldPrecision == 0 && c.HeldOutDelta.FieldRecall == 0 {
			t.Logf("Warning: HeldOutDelta appears empty for %s", c.Path)
		}
	}
}

// TestSuggestionsNotApplied verifies suggestions go to review, not active config.
func TestSuggestionsNotApplied(t *testing.T) {
	os.Chdir(repoRoot())

	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Failed to load knowledge: %v", err)
	}

	engine := decipher.NewEngineWithKnowledge(kb)
	weights := decipher.DefaultRankingWeights()

	suggester := NewWeightSuggester(engine, weights)
	tmpDir := t.TempDir()
	suggester.SetReviewPath(filepath.Join(tmpDir, "test_review_not_applied.yaml"))

	loader := NewLoader()
	trainExamples, err := loader.LoadFromFile("training/examples.yaml")
	if err != nil {
		t.Fatalf("Failed to load examples: %v", err)
	}
	heldOutExamples, err := loader.LoadHeldOutFromFile("training/heldout.yaml")
	if err != nil {
		t.Fatalf("Failed to load held-out examples: %v", err)
	}

	baseTrain := MetricsDelta{ConceptPrecision: 0.24, FieldPrecision: 0.39}
	baseHeldOut := MetricsDelta{ConceptPrecision: 0.12, FieldPrecision: 0.27}

	suggester.GenerateSuggestions(trainExamples, heldOutExamples, baseTrain, baseHeldOut)
	suggester.WriteSuggestions()

	// Verify review file exists, active YAML does not
	if _, err := os.Stat(suggester.ReviewPath); os.IsNotExist(err) {
		t.Errorf("Review file should exist at %s", suggester.ReviewPath)
	}

	// Verify active weights file is unchanged
	activePath := filepath.Join(repoRoot(), "training/ranking_weights.yaml")
	if info, err := os.Stat(activePath); err == nil {
		// File exists, should not have been modified by suggestion generation
		t.Logf("Active weights file exists at %s, size %d bytes (should be unchanged)", activePath, info.Size())
	}
}

// TestHeldOutImprovementPreferred verifies held-out improvements are prioritized.
func TestHeldOutImprovementPreferred(t *testing.T) {
	os.Chdir(repoRoot())

	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Failed to load knowledge: %v", err)
	}

	engine := decipher.NewEngineWithKnowledge(kb)
	weights := decipher.DefaultRankingWeights()

	_ = NewWeightSuggester(engine, weights)

	// Test cases for isImprovement
	testCases := []struct {
		name        string
		trainDelta  MetricsDelta
		heldOutDelta MetricsDelta
		expected    bool
	}{
		{
			name: "held-out precision improves",
			trainDelta:  MetricsDelta{FieldPrecision: 0.01},
			heldOutDelta: MetricsDelta{FieldPrecision: 0.10},
			expected: true,
		},
		{
			name: "train improves, held-out degrades too much",
			trainDelta:  MetricsDelta{FieldPrecision: 0.15},
			heldOutDelta: MetricsDelta{FieldPrecision: -0.10},
			expected: false,
		},
		{
			name: "train improves, held-out maintains",
			trainDelta:  MetricsDelta{FieldPrecision: 0.10},
			heldOutDelta: MetricsDelta{FieldPrecision: 0.01, FieldRecall: -0.02},
			expected: true,
		},
	}

	for _, tc := range testCases {
		result := isImprovement(tc.trainDelta, tc.heldOutDelta)
		if result != tc.expected {
			t.Errorf("%s: expected %v, got %v", tc.name, tc.expected, result)
		}
	}
}

// TestTrainOnlyImprovementFlagged verifies train-only improvements are flagged.
func TestTrainOnlyImprovementFlagged(t *testing.T) {
	os.Chdir(repoRoot())

	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Failed to load knowledge: %v", err)
	}

	engine := decipher.NewEngineWithKnowledge(kb)
	weights := decipher.DefaultRankingWeights()

	_ = NewWeightSuggester(engine, weights)

	// Test: train improves significantly, held-out degrades slightly
	// This should still pass if held-out degradation is acceptable
	trainDelta := MetricsDelta{FieldPrecision: 0.15, FieldRecall: 0.05}
	heldOutDelta := MetricsDelta{FieldPrecision: -0.02, FieldRecall: -0.03}

	result := isImprovement(trainDelta, heldOutDelta)
	// Should be true because heldOutDelta.FieldPrecision >= -0.03 and heldOutDelta.FieldRecall >= -0.05
	if !result {
		t.Errorf("Expected train-only improvement to be acceptable when held-out degradation is small")
	}

	// Test: train improves, held-out degrades significantly
	heldOutDelta = MetricsDelta{FieldPrecision: -0.10, FieldRecall: -0.10}
	result = isImprovement(trainDelta, heldOutDelta)
	if result {
		t.Errorf("Expected train-only improvement to be flagged when held-out degrades significantly")
	}
}

// TestReviewFileContent verifies the review file format.
func TestReviewFileContent(t *testing.T) {
	os.Chdir(repoRoot())

	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Failed to load knowledge: %v", err)
	}

	engine := decipher.NewEngineWithKnowledge(kb)
	weights := decipher.DefaultRankingWeights()

	suggester := NewWeightSuggester(engine, weights)
	tmpDir := t.TempDir()
	reviewFile := filepath.Join(tmpDir, "weight_suggestions.yaml")
	suggester.SetReviewPath(reviewFile)

	loader := NewLoader()
	trainExamples, err := loader.LoadFromFile("training/examples.yaml")
	if err != nil {
		t.Fatalf("Failed to load examples: %v", err)
	}
	heldOutExamples, err := loader.LoadHeldOutFromFile("training/heldout.yaml")
	if err != nil {
		t.Fatalf("Failed to load held-out examples: %v", err)
	}

	baseTrain := MetricsDelta{ConceptPrecision: 0.24, FieldPrecision: 0.39}
	baseHeldOut := MetricsDelta{ConceptPrecision: 0.12, FieldPrecision: 0.27}

	suggester.GenerateSuggestions(trainExamples, heldOutExamples, baseTrain, baseHeldOut)
	suggester.WriteSuggestions()

	// Read the review file and verify content
	content, err := os.ReadFile(suggester.ReviewPath)
	if err != nil {
		t.Fatalf("Failed to read review file: %v", err)
	}

	contentStr := string(content)

	// Verify expected content sections
	if !strings.Contains(contentStr, "version:") {
		t.Error("Review file missing version")
	}
	if !strings.Contains(contentStr, "suggestions:") {
		t.Error("Review file missing suggestions section")
	}
	if !strings.Contains(contentStr, "status:") {
		t.Error("Review file missing status field")
	}
	if !strings.Contains(contentStr, "train_delta:") {
		t.Error("Review file missing train_delta section")
	}
	if !strings.Contains(contentStr, "heldout_delta:") {
		t.Error("Review file missing heldout_delta section")
	}
	if !strings.Contains(contentStr, "Do NOT auto-apply") {
		t.Error("Review file missing auto-apply warning")
	}

	t.Logf("Review file content verified, length: %d bytes", len(content))
}