package training

import (
	"fmt"

	"socrates/internal/decipher"
)

// =============================================================================
// Evaluation Report Types (with ranking weights support)
// =============================================================================

// EvaluationReport is the complete evaluation report with train/held-out splits and weight info.
type EvaluationReport struct {
	// Weight info
	WeightSet string // "default" or path to weights file
	Weights   decipher.RankingWeights

	// Train split results
	Train EvaluationSplit

	// Held-out split results
	HeldOut EvaluationSplit

	// Overall combined metrics
	TotalExamples      int
	TotalPassed        int
	TotalFailed        int
	AvgConceptPrec     float64
	AvgConceptRec      float64
	AvgFieldPrec       float64
	AvgFieldRec        float64

	// Weight impact notes
	WeightNotes []string
}

// EvaluationSplit contains metrics for a single split (train or held-out).
type EvaluationSplit struct {
	Examples      int
	Passed        int
	Failed        int
	ConceptPrec   float64
	ConceptRec    float64
	FieldPrec     float64
	FieldRec      float64
	MissedFields  []string
	FalseFields   []string
	Detailed      []ExampleResult
}

// FeatureContribution holds contribution info for a feature in debug mode.
type FeatureContribution struct {
	Feature   string
	Value     float64
	Weight    float64
	Contrib   float64
	Source    string
}

// FeatureReport holds feature contributions for a single activation.
type FeatureReport struct {
	ConceptID  string
	FieldID    string
	IsExpected bool // True if this was in expected concepts/fields
	Features   []FeatureContribution
	TotalScore float64
	Passed     bool
}

// NewEvaluationReport creates a new evaluation report with the given weights.
func NewEvaluationReport(weightSet string, weights decipher.RankingWeights) *EvaluationReport {
	return &EvaluationReport{
		WeightSet: weightSet,
		Weights:   weights,
		Train:     EvaluationSplit{Detailed: make([]ExampleResult, 0)},
		HeldOut:   EvaluationSplit{Detailed: make([]ExampleResult, 0)},
	}
}

// AddTrainResult adds a train split example result.
func (r *EvaluationReport) AddTrainResult(result ExampleResult) {
	r.Train.Detailed = append(r.Train.Detailed, result)
}

// AddHeldOutResult adds a held-out split example result.
func (r *EvaluationReport) AddHeldOutResult(result ExampleResult) {
	r.HeldOut.Detailed = append(r.HeldOut.Detailed, result)
}

// Finalize computes aggregate metrics after all results are added.
func (r *EvaluationReport) Finalize() {
	// Train split aggregates
	r.Train.Examples = len(r.Train.Detailed)
	r.Train.Passed = countPassed(r.Train.Detailed)
	r.Train.Failed = r.Train.Examples - r.Train.Passed
	r.Train.ConceptPrec, r.Train.ConceptRec, r.Train.FieldPrec, r.Train.FieldRec = computeAverages(r.Train.Detailed)
	r.Train.MissedFields = collectMissed(r.Train.Detailed)
	r.Train.FalseFields = collectFalse(r.Train.Detailed)

	// Held-out split aggregates
	r.HeldOut.Examples = len(r.HeldOut.Detailed)
	r.HeldOut.Passed = countPassed(r.HeldOut.Detailed)
	r.HeldOut.Failed = r.HeldOut.Examples - r.HeldOut.Passed
	r.HeldOut.ConceptPrec, r.HeldOut.ConceptRec, r.HeldOut.FieldPrec, r.HeldOut.FieldRec = computeAverages(r.HeldOut.Detailed)
	r.HeldOut.MissedFields = collectMissed(r.HeldOut.Detailed)
	r.HeldOut.FalseFields = collectFalse(r.HeldOut.Detailed)

	// Overall combined
	r.TotalExamples = r.Train.Examples + r.HeldOut.Examples
	r.TotalPassed = r.Train.Passed + r.HeldOut.Passed
	r.TotalFailed = r.Train.Failed + r.HeldOut.Failed

	allDetailed := append(r.Train.Detailed, r.HeldOut.Detailed...)
	r.AvgConceptPrec, r.AvgConceptRec, r.AvgFieldPrec, r.AvgFieldRec = computeAverages(allDetailed)
}

func countPassed(results []ExampleResult) int {
	count := 0
	for _, r := range results {
		if r.Passed {
			count++
		}
	}
	return count
}

func computeAverages(results []ExampleResult) (cPrec, cRec, fPrec, fRec float64) {
	if len(results) == 0 {
		return 0, 0, 0, 0
	}
	for _, r := range results {
		cPrec += r.ConceptPrecision
		cRec += r.ConceptRecall
		fPrec += r.FieldPrecision
		fRec += r.FieldRecall
	}
	n := float64(len(results))
	return cPrec / n, cRec / n, fPrec / n, fRec / n
}

func collectMissed(results []ExampleResult) []string {
	missed := make(map[string]bool)
	for _, r := range results {
		for _, f := range r.MissedExpectedFields {
			missed[f] = true
		}
	}
	return mapKeys(missed)
}

func collectFalse(results []ExampleResult) []string {
	falseFields := make(map[string]bool)
	for _, r := range results {
		for _, f := range r.FalseActivatedFields {
			falseFields[f] = true
		}
	}
	return mapKeys(falseFields)
}

func mapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// FormatReport formats the evaluation report for display.
func (r *EvaluationReport) FormatReport() string {
	s := "\n" + repeat("=", 60) + "\n"
	s += "SOCRATES TRAINING EVALUATION REPORT\n"
	s += repeat("=", 60) + "\n\n"

	// Weight info
	s += "Weight Set: " + r.WeightSet + "\n"
	s += fmt.Sprintf("Total Examples: %d\n", r.TotalExamples)
	s += "\n"

	// Train split
	s += "TRAIN SPLIT\n"
	s += repeat("-", 40) + "\n"
	s += formatSplitSummary(&r.Train)
	s += "\n"

	// Held-out split
	s += "HELD-OUT SPLIT\n"
	s += repeat("-", 40) + "\n"
	s += formatSplitSummary(&r.HeldOut)
	s += "\n"

	// Combined
	s += "COMBINED METRICS\n"
	s += repeat("-", 40) + "\n"
	s += formatMetrics(r.AvgConceptPrec, r.AvgConceptRec, r.AvgFieldPrec, r.AvgFieldRec, r.TotalPassed, r.TotalFailed)

	return s
}

// formatSplitSummary formats a single split summary.
func formatSplitSummary(split *EvaluationSplit) string {
	if split.Examples == 0 {
		return "  (no examples)\n\n"
	}
	s := ""
	s += fmt.Sprintf("  Examples: %d\n", split.Examples)
	s += fmt.Sprintf("  Passed: %d\n", split.Passed)
	s += fmt.Sprintf("  Failed: %d\n", split.Failed)
	s += "\n"
	s += formatMetrics(split.ConceptPrec, split.ConceptRec, split.FieldPrec, split.FieldRec, split.Passed, split.Failed)
	s += "\n"
	return s
}

// formatMetrics formats precision/recall metrics.
func formatMetrics(cp, cr, fp, fr float64, passed, failed int) string {
	s := fmt.Sprintf("  Concept Precision: %.2f\n", cp)
	s += fmt.Sprintf("  Concept Recall: %.2f\n", cr)
	s += fmt.Sprintf("  Field Precision: %.2f\n", fp)
	s += fmt.Sprintf("  Field Recall: %.2f\n", fr)
	return s
}

func repeat(c string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += c
	}
	return result
}

// FormatDetailedResults formats detailed results for a split.
// This method exists on EvaluationSplit for compatibility with existing code.
func (split *EvaluationSplit) FormatDetailedResults(examples []Example) string {
	output := "Detailed Results:\n"
	for i, ex := range split.Detailed {
		status := "PASS"
		if !ex.Passed {
			status = "FAIL"
		}
		exID := ex.ExampleID
		if exID == "" {
			exID = fmt.Sprintf("example[%d]", i)
		}
		output += fmt.Sprintf("  [%s] %s: %s\n", status, exID, ex.Input)
		output += fmt.Sprintf("    Concepts: %d hits, %d misses, %d false pos (prec=%.2f, rec=%.2f)\n",
			ex.ConceptHits, ex.ConceptMisses, ex.ConceptFalsePos,
			ex.ConceptPrecision, ex.ConceptRecall)
		output += fmt.Sprintf("    Fields: %d hits, %d misses, %d false pos (prec=%.2f, rec=%.2f)\n",
			ex.FieldHits, ex.FieldMisses, ex.FieldFalsePos,
			ex.FieldPrecision, ex.FieldRecall)
		if ex.FailedReason != "" {
			output += fmt.Sprintf("    Failed: %s\n", ex.FailedReason)
		}
	}
	return output
}