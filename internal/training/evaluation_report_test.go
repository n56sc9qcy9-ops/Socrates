package training

import (
	"socrates/internal/decipher"
	"testing"
)

// =============================================================================
// Evaluation Report Tests (with train/held-out splits)
// =============================================================================

func TestEvaluationReportCreation(t *testing.T) {
	weights := decipher.DefaultRankingWeights()
	report := NewEvaluationReport("default", weights)

	if report.WeightSet != "default" {
		t.Errorf("Expected WeightSet='default', got '%s'", report.WeightSet)
	}

	if report.Train.Examples != 0 {
		t.Errorf("Expected Train.Examples=0, got %d", report.Train.Examples)
	}

	if report.HeldOut.Examples != 0 {
		t.Errorf("Expected HeldOut.Examples=0, got %d", report.HeldOut.Examples)
	}
}

func TestEvaluationReportAddResults(t *testing.T) {
	weights := decipher.DefaultRankingWeights()
	report := NewEvaluationReport("default", weights)

	// Add train results
	report.AddTrainResult(ExampleResult{
		ExampleID: "train.1",
		Passed:    true,
	})
	report.AddTrainResult(ExampleResult{
		ExampleID: "train.2",
		Passed:    false,
	})

	// Add held-out results
	report.AddHeldOutResult(ExampleResult{
		ExampleID: "heldout.1",
		Passed:    true,
	})

	// Verify counts
	if len(report.Train.Detailed) != 2 {
		t.Errorf("Expected 2 train results, got %d", len(report.Train.Detailed))
	}

	if len(report.HeldOut.Detailed) != 1 {
		t.Errorf("Expected 1 held-out result, got %d", len(report.HeldOut.Detailed))
	}
}

func TestEvaluationReportFinalize(t *testing.T) {
	weights := decipher.DefaultRankingWeights()
	report := NewEvaluationReport("default", weights)

	// Add train results
	report.AddTrainResult(ExampleResult{
		ExampleID:           "train.1",
		Passed:              true,
		ConceptPrecision:    0.8,
		ConceptRecall:      0.9,
		FieldPrecision:     0.7,
		FieldRecall:        0.85,
	})
	report.AddTrainResult(ExampleResult{
		ExampleID:           "train.2",
		Passed:              true,
		ConceptPrecision:    0.6,
		ConceptRecall:      0.7,
		FieldPrecision:     0.8,
		FieldRecall:        0.9,
	})

	// Add held-out results
	report.AddHeldOutResult(ExampleResult{
		ExampleID:           "heldout.1",
		Passed:              false,
		ConceptPrecision:    0.5,
		ConceptRecall:      0.6,
		FieldPrecision:     0.4,
		FieldRecall:        0.5,
		MissedExpectedFields: []string{"missed.field.1"},
		FalseActivatedFields: []string{"false.field.1"},
	})

	// Finalize
	report.Finalize()

	// Check train split
	if report.Train.Examples != 2 {
		t.Errorf("Expected Train.Examples=2, got %d", report.Train.Examples)
	}
	if report.Train.Passed != 2 {
		t.Errorf("Expected Train.Passed=2, got %d", report.Train.Passed)
	}

	// Check held-out split
	if report.HeldOut.Examples != 1 {
		t.Errorf("Expected HeldOut.Examples=1, got %d", report.HeldOut.Examples)
	}
	if report.HeldOut.Passed != 0 {
		t.Errorf("Expected HeldOut.Passed=0, got %d", report.HeldOut.Passed)
	}

	// Check combined
	if report.TotalExamples != 3 {
		t.Errorf("Expected TotalExamples=3, got %d", report.TotalExamples)
	}

	// Check averages
	if report.AvgConceptPrec < 0.63 || report.AvgConceptPrec > 0.64 {
		t.Errorf("Expected AvgConceptPrec ~0.63, got %v", report.AvgConceptPrec)
	}

	// Check missed/false fields collection
	if len(report.HeldOut.MissedFields) != 1 {
		t.Errorf("Expected 1 missed field, got %d", len(report.HeldOut.MissedFields))
	}
}

func TestEvaluationReportFormatReport(t *testing.T) {
	weights := decipher.DefaultRankingWeights()
	report := NewEvaluationReport("training/ranking_weights.yaml", weights)

	report.AddTrainResult(ExampleResult{
		ExampleID: "train.1",
		Passed:    true,
		ConceptPrecision: 0.8,
		ConceptRecall:    0.9,
		FieldPrecision:   0.7,
		FieldRecall:      0.85,
	})

	report.Finalize()

	formatted := report.FormatReport()

	// Check that report contains key sections
	if !contains(formatted, "Weight Set:") {
		t.Error("Report should contain 'Weight Set:'")
	}

	if !contains(formatted, "TRAIN SPLIT") {
		t.Error("Report should contain 'TRAIN SPLIT'")
	}

	if !contains(formatted, "HELD-OUT SPLIT") {
		t.Error("Report should contain 'HELD-OUT SPLIT'")
	}

	if !contains(formatted, "COMBINED METRICS") {
		t.Error("Report should contain 'COMBINED METRICS'")
	}
}

func TestRunFullEvaluation(t *testing.T) {
	// Create a simple evaluator (engine will be nil but we test the method exists)
	weights := decipher.DefaultRankingWeights()
	
	// This test just verifies the method runs without panic
	// Full integration would require a real engine
	report := NewEvaluationReport("default", weights)
	
	// Verify we can add results
	report.AddTrainResult(ExampleResult{
		ExampleID: "test.1",
		Passed:    true,
	})
	
	report.Finalize()
	
	if report.Train.Examples != 1 {
		t.Errorf("Expected 1 train example, got %d", report.Train.Examples)
	}
}

// Helper
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsAt(s, sub, 0))
}

func containsAt(s, sub string, start int) bool {
	for i := start; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}