package training

import "fmt"

// ============================================================
// Training Example Types
// ============================================================

// Example represents a curated training example for supervised ranking.
// Examples are stored separately from active knowledge - they verify the engine,
// they are not automatically added as truth.
type Example struct {
	// ID is a unique identifier for this example.
	ID string

	// Input is the text to analyze.
	Input string

	// ExpectedConcepts are the concept IDs expected to activate.
	// These must exist in the active knowledge base.
	ExpectedConcepts []string

	// ExpectedFields are the meaning-frequency IDs expected to be active.
	// These must exist in the frequency profiles.
	ExpectedFields []string

	// ExpectedQuality is the expected harmonic quality.
	// Optional: if set, evaluator checks agreement.
	ExpectedQuality string

	// Confidence indicates the example's provenance and reliability.
	// e.g., "curated", "human_review", "verified"
	Confidence string

	// Source describes where this example came from.
	Source string

	// Notes contains optional notes about this example.
	Notes string
}

// Examples is a collection of training examples.
type Examples []Example

// ============================================================
// Evaluation Types
// ============================================================

// ExampleResult is the evaluation result for a single example.
type ExampleResult struct {
	ExampleID    string
	Input        string
	Passed       bool
	FailedReason string // Why the example failed, if any

	// Concept metrics
	ConceptPrecision float64
	ConceptRecall    float64
	ConceptHits      int // expected concepts that were found
	ConceptMisses    int // expected concepts not found
	ConceptFalsePos  int // found concepts not expected

	// Meaning-frequency field metrics
	FieldPrecision float64
	FieldRecall    float64
	FieldHits      int
	FieldMisses    int
	FieldFalsePos  int

	// Quality agreement
	QualityAgreed   bool
	ExpectedQuality string
	ActualQuality   string

	// Details for debugging
	ActivatedConcepts    []string
	ActivatedFields      []string
	MissedExpectedFields []string
	FalseActivatedFields []string
	EvidencePathsValid   bool // matched fields have evidence paths
}

// EvaluationResult is the complete evaluation report.
type EvaluationResult struct {
	TotalExamples  int
	PassedExamples int
	FailedExamples int

	// Aggregate metrics
	AvgConceptPrecision float64
	AvgConceptRecall    float64
	AvgFieldPrecision   float64
	AvgFieldRecall      float64

	// Precision warning threshold
	// If AvgFieldPrecision < MinAcceptablePrecision, the evaluation has excessive false activations
	MinAcceptablePrecision float64

	// Quality agreement
	QualityAgreements    int
	QualityDisagreements int

	// Detailed results per example
	ExampleResults []ExampleResult

	// Summary of missed and false activations
	AllMissedFields []string
	AllFalseFields  []string
}

// Metrics returns a human-readable summary of evaluation metrics.
func (r *EvaluationResult) Metrics() string {
	return formatEvaluationMetrics(r)
}

// formatEvaluationMetrics formats evaluation metrics as a string.
func formatEvaluationMetrics(r *EvaluationResult) string {
	s := "Evaluation Metrics:\n"
	s += "==================\n"
	s += fmt.Sprintf("Total examples: %d\n", r.TotalExamples)
	s += fmt.Sprintf("Passed: %d\n", r.PassedExamples)
	s += fmt.Sprintf("Failed: %d\n", r.FailedExamples)
	s += "\n"
	s += fmt.Sprintf("Avg concept precision: %.2f\n", r.AvgConceptPrecision)
	s += fmt.Sprintf("Avg concept recall: %.2f\n", r.AvgConceptRecall)
	s += fmt.Sprintf("Avg field precision: %.2f\n", r.AvgFieldPrecision)
	s += fmt.Sprintf("Avg field recall: %.2f\n", r.AvgFieldRecall)
	s += "\n"
	s += fmt.Sprintf("Quality agreements: %d\n", r.QualityAgreements)
	s += fmt.Sprintf("Quality disagreements: %d\n", r.QualityDisagreements)
	return s
}
