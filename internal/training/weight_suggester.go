package training

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"socrates/internal/decipher"
)

// WeightCandidate represents a proposed change to a weight configuration.
type WeightCandidate struct {
	Path       string      // Dot-path to the weight, e.g. "exact_match.base"
	Current    interface{} // Current value
	Suggested  interface{} // Proposed value
	Delta      string      // Human-readable change description
	Rationale  string      // Why this change is suggested
	TrainDelta MetricsDelta // How metrics change on train set
	HeldOutDelta MetricsDelta // How metrics change on held-out set
	Status     string      // "pending", "accepted", "rejected"
}

// MetricsDelta represents changes in evaluation metrics.
type MetricsDelta struct {
	ConceptPrecision float64
	ConceptRecall   float64
	FieldPrecision  float64
	FieldRecall     float64
	PassRateDelta   float64 // Change in pass rate (e.g., +2/8)
}

// WeightSuggester generates weight change suggestions based on evaluation results.
type WeightSuggester struct {
	engine           *decipher.Engine
	baseWeights      decipher.RankingWeights
	candidates       []WeightCandidate
	ReviewPath       string // Public path for review files
}

// NewWeightSuggester creates a new weight suggester.
func NewWeightSuggester(engine *decipher.Engine, weights decipher.RankingWeights) *WeightSuggester {
	return &WeightSuggester{
		engine:      engine,
		baseWeights: weights,
		candidates:  make([]WeightCandidate, 0),
		ReviewPath:  "review/weight_suggestions.yaml",
	}
}

// SetReviewPath sets the path for writing review files.
func (s *WeightSuggester) SetReviewPath(path string) {
	s.ReviewPath = path
}

// GenerateSuggestions creates weight change candidates based on evaluation results.
func (s *WeightSuggester) GenerateSuggestions(
	trainExamples, heldOutExamples Examples,
	baseTrainMetrics, baseHeldOutMetrics MetricsDelta,
) []WeightCandidate {
	s.candidates = make([]WeightCandidate, 0)

	// Define candidate weight changes to explore
	// Each candidate is a delta from current weights
	candidateChanges := s.defineCandidates()

	evaluator := NewEvaluator(s.engine)

	_ = evaluator // unused but shows intent

	for _, change := range candidateChanges {
		// Create modified weights
		modifiedWeights := applyWeightChange(s.baseWeights, change)

		// Create evaluator with modified weights
		eval := NewEvaluatorWithWeights(s.engine, modifiedWeights, "candidate")

		// Evaluate on train and held-out
		trainReport := eval.EvaluateWithReport(trainExamples)
		heldOutReport := &EvaluationReport{}
		eval.EvaluateHeldOut(heldOutExamples, heldOutReport)
		heldOutReport.Finalize()

		// Calculate deltas
		trainDelta := calculateMetricsDelta(baseTrainMetrics, trainReport)
		heldOutDelta := calculateMetricsDelta(baseHeldOutMetrics, heldOutReport)

		// Check if this is a good suggestion
		// Must improve held-out OR maintain train while improving held-out
		// Must not degrade held-out field precision significantly
		if isImprovement(trainDelta, heldOutDelta) {
			candidate := WeightCandidate{
				Path:          change.Path,
				Current:       change.CurrentValue,
				Suggested:     change.SuggestedValue,
				Delta:         change.DeltaDescription,
				Rationale:     change.Rationale,
				TrainDelta:    trainDelta,
				HeldOutDelta:  heldOutDelta,
				Status:        "pending",
			}
			s.candidates = append(s.candidates, candidate)
		}
	}

	return s.candidates
}

// weightChange represents a single weight parameter change.
type weightChange struct {
	Path            string
	CurrentValue    interface{}
	SuggestedValue  interface{}
	DeltaDescription string
	Rationale       string
	ApplyFunc       func(decipher.RankingWeights) decipher.RankingWeights
}

// defineCandidates returns the set of weight changes to explore.
func (s *WeightSuggester) defineCandidates() []weightChange {
	return []weightChange{
		// Exact match base weight - increase for better precision
		{
			Path:            "exact_match_base",
			CurrentValue:    s.baseWeights.ExactMatchBase,
			SuggestedValue:  s.baseWeights.ExactMatchBase * 1.1,
			DeltaDescription: fmt.Sprintf("increase from %.2f to %.2f", s.baseWeights.ExactMatchBase, s.baseWeights.ExactMatchBase*1.1),
			Rationale:       "Boost exact matches to improve precision",
			ApplyFunc: func(w decipher.RankingWeights) decipher.RankingWeights {
				w.ExactMatchBase *= 1.1
				return w
			},
		},
		// Fuzzy match base weight - decrease to reduce false activations
		{
			Path:            "fuzzy_match_base",
			CurrentValue:    s.baseWeights.FuzzyMatchBase,
			SuggestedValue:  s.baseWeights.FuzzyMatchBase * 0.9,
			DeltaDescription: fmt.Sprintf("decrease from %.2f to %.2f", s.baseWeights.FuzzyMatchBase, s.baseWeights.FuzzyMatchBase*0.9),
			Rationale:       "Reduce fuzzy weight to reduce false activations",
			ApplyFunc: func(w decipher.RankingWeights) decipher.RankingWeights {
				w.FuzzyMatchBase *= 0.9
				return w
			},
		},
		// Graph expansion weight - decrease to limit propagation
		{
			Path:            "graph_expansion_weight",
			CurrentValue:    s.baseWeights.GraphExpansionWeight,
			SuggestedValue:  s.baseWeights.GraphExpansionWeight * 0.9,
			DeltaDescription: fmt.Sprintf("decrease from %.2f to %.2f", s.baseWeights.GraphExpansionWeight, s.baseWeights.GraphExpansionWeight*0.9),
			Rationale:       "Limit graph propagation to reduce false activations",
			ApplyFunc: func(w decipher.RankingWeights) decipher.RankingWeights {
				w.GraphExpansionWeight *= 0.9
				return w
			},
		},
		// Graph depth penalty - increase to penalize deep propagation
		{
			Path:            "graph_depth_penalty",
			CurrentValue:    s.baseWeights.GraphDepthPenalty,
			SuggestedValue:  s.baseWeights.GraphDepthPenalty * 1.2,
			DeltaDescription: fmt.Sprintf("increase from %.2f to %.2f", s.baseWeights.GraphDepthPenalty, s.baseWeights.GraphDepthPenalty*1.2),
			Rationale:       "Increase depth penalty to reduce distant concept activations",
			ApplyFunc: func(w decipher.RankingWeights) decipher.RankingWeights {
				w.GraphDepthPenalty *= 1.2
				return w
			},
		},
		// Harmonic profile weight - increase for field precision
		{
			Path:            "harmonic_profile_weight",
			CurrentValue:    s.baseWeights.HarmonicProfileWeight,
			SuggestedValue:  s.baseWeights.HarmonicProfileWeight * 1.1,
			DeltaDescription: fmt.Sprintf("increase from %.2f to %.2f", s.baseWeights.HarmonicProfileWeight, s.baseWeights.HarmonicProfileWeight*1.1),
			Rationale:       "Boost harmonic profile influence for better field precision",
			ApplyFunc: func(w decipher.RankingWeights) decipher.RankingWeights {
				w.HarmonicProfileWeight *= 1.1
				return w
			},
		},
		// Plausible confidence weight - decrease slightly
		{
			Path:            "confidence_plausible",
			CurrentValue:    s.baseWeights.ConfidencePlausible,
			SuggestedValue:  s.baseWeights.ConfidencePlausible * 0.95,
			DeltaDescription: fmt.Sprintf("decrease from %.2f to %.2f", s.baseWeights.ConfidencePlausible, s.baseWeights.ConfidencePlausible*0.95),
			Rationale:       "Slightly reduce plausible confidence to filter weak evidence",
			ApplyFunc: func(w decipher.RankingWeights) decipher.RankingWeights {
				w.ConfidencePlausible *= 0.95
				return w
			},
		},
		// Channel diversity threshold - increase to require more channels
		{
			Path:            "channel_diversity_threshold",
			CurrentValue:    s.baseWeights.ChannelDiversityThreshold,
			SuggestedValue:  s.baseWeights.ChannelDiversityThreshold + 1,
			DeltaDescription: fmt.Sprintf("increase from %d to %d", s.baseWeights.ChannelDiversityThreshold, s.baseWeights.ChannelDiversityThreshold+1),
			Rationale:       "Require more channel diversity to activate fields",
			ApplyFunc: func(w decipher.RankingWeights) decipher.RankingWeights {
				w.ChannelDiversityThreshold++
				return w
			},
		},
		// Multi-method bonus - increase to reward agreement
		{
			Path:            "multi_method_bonus",
			CurrentValue:    s.baseWeights.MultiMethodBonus,
			SuggestedValue:  s.baseWeights.MultiMethodBonus + 0.05,
			DeltaDescription: fmt.Sprintf("increase from %.2f to %.2f", s.baseWeights.MultiMethodBonus, s.baseWeights.MultiMethodBonus+0.05),
			Rationale:       "Reward stronger agreement between multiple methods",
			ApplyFunc: func(w decipher.RankingWeights) decipher.RankingWeights {
				w.MultiMethodBonus += 0.05
				return w
			},
		},
	}
}

// applyWeightChange applies a weight change to the base weights.
func applyWeightChange(base decipher.RankingWeights, change weightChange) decipher.RankingWeights {
	return change.ApplyFunc(base)
}

// calculateMetricsDelta computes the change in metrics between base and new evaluation.
func calculateMetricsDelta(base MetricsDelta, report *EvaluationReport) MetricsDelta {
	if report == nil {
		return MetricsDelta{}
	}

	total := report.Train.Passed + report.Train.Failed

	var basePassRate, newPassRate, passRateDelta float64
	if total > 0 {
		newPassRate = float64(report.Train.Passed) / float64(total)
		basePassRate = 0 // Assume base was evaluated separately
		passRateDelta = newPassRate - basePassRate
	}

	return MetricsDelta{
		ConceptPrecision: report.Train.ConceptPrec - base.ConceptPrecision,
		ConceptRecall:    report.Train.ConceptRec - base.ConceptRecall,
		FieldPrecision:   report.Train.FieldPrec - base.FieldPrecision,
		FieldRecall:      report.Train.FieldRec - base.FieldRecall,
		PassRateDelta:    passRateDelta,
	}
}

// isImprovement determines if a weight change is an improvement.
// It prefers held-out improvement over train-only improvement.
// It flags train-only improvements that harm held-out metrics.
// Returns false if there's no change (delta is 0).
func isImprovement(trainDelta, heldOutDelta MetricsDelta) bool {
	// No improvement if nothing changed
	if trainDelta.FieldPrecision == 0 && trainDelta.FieldRecall == 0 &&
		heldOutDelta.FieldPrecision == 0 && heldOutDelta.FieldRecall == 0 {
		return false
	}

	// Good if held-out field precision improves
	if heldOutDelta.FieldPrecision > 0 {
		return true
	}

	// Good if held-out field recall maintains (no significant drop) and precision doesn't degrade
	if heldOutDelta.FieldPrecision >= 0 &&
		heldOutDelta.FieldRecall >= -0.05 {
		return true
	}

	// Acceptable if train improves significantly and held-out doesn't degrade much
	if trainDelta.FieldPrecision > 0.05 &&
		heldOutDelta.FieldPrecision >= -0.05 &&
		heldOutDelta.FieldRecall >= -0.05 {
		return true
	}

	// Acceptable if held-out recall improves even if precision drops slightly
	if heldOutDelta.FieldRecall > 0.05 &&
		heldOutDelta.FieldPrecision >= -0.03 {
		return true
	}

	return false
}

// WriteSuggestions writes the generated suggestions to a review file.
func (s *WeightSuggester) WriteSuggestions() error {
	// Ensure review directory exists
	reviewDir := filepath.Dir(s.ReviewPath)
	if err := os.MkdirAll(reviewDir, 0755); err != nil {
		return fmt.Errorf("failed to create review directory: %w", err)
	}

	// Generate YAML content
	content := generateSuggestionsYAML(s.candidates)

	// Write to file
	if err := os.WriteFile(s.ReviewPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write suggestions file: %w", err)
	}

	return nil
}

// generateSuggestionsYAML creates YAML representation of suggestions.
func generateSuggestionsYAML(candidates []WeightCandidate) string {
	content := `# Weight Suggestions Review File
# Generated: ` + time.Now().Format(time.RFC3339) + `
#
# This file contains proposed weight changes.
# Each suggestion has:
#   - path: dot-path to the weight parameter
#   - current: current value
#   - suggested: proposed value
#   - rationale: why this change is suggested
#   - train_delta: metric changes on train set
#   - heldout_delta: metric changes on held-out set
#   - status: "pending", "accepted", or "rejected"
#
# Do NOT auto-apply these changes. Review and accept/reject manually.
#
# =============================================================================

version: "1.0"
generated: "` + time.Now().Format(time.RFC3339) + `"

suggestions:`

	if len(candidates) == 0 {
		content += `
  - message: "No suggestions generated. Current weights are acceptable."
    status: "none"
`
	} else {
		for i, c := range candidates {
			content += fmt.Sprintf(`
  - id: suggestion_%d
    path: "%s"
    current: %v
    suggested: %v
    delta: "%s"
    rationale: "%s"
    status: "%s"
    train_delta:
      concept_precision: %.4f
      concept_recall: %.4f
      field_precision: %.4f
      field_recall: %.4f
    heldout_delta:
      concept_precision: %.4f
      concept_recall: %.4f
      field_precision: %.4f
      field_recall: %.4f`,
				i+1,
				c.Path,
				formatValue(c.Current),
				formatValue(c.Suggested),
				c.Delta,
				c.Rationale,
				c.Status,
				c.TrainDelta.ConceptPrecision,
				c.TrainDelta.ConceptRecall,
				c.TrainDelta.FieldPrecision,
				c.TrainDelta.FieldRecall,
				c.HeldOutDelta.ConceptPrecision,
				c.HeldOutDelta.ConceptRecall,
				c.HeldOutDelta.FieldPrecision,
				c.HeldOutDelta.FieldRecall,
			)
		}
	}

	content += `

# =============================================================================
# Review Instructions
# =============================================================================
# 1. Examine each suggestion's heldout_delta (held-out metrics).
# 2. Accept suggestions that improve held-out metrics.
# 3. Reject suggestions that improve train but degrade held-out metrics.
# 4. Update status to "accepted" or "rejected" after review.
# 5. Do NOT delete this file after accepting - keep for audit trail.
`

	return content
}

// formatValue formats a value for YAML output.
func formatValue(v interface{}) string {
	switch val := v.(type) {
	case float64:
		return fmt.Sprintf("%.2f", val)
	case int:
		return fmt.Sprintf("%d", val)
	case string:
		return fmt.Sprintf(`"%s"`, val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// GetCandidates returns the generated candidates.
func (s *WeightSuggester) GetCandidates() []WeightCandidate {
	return s.candidates
}