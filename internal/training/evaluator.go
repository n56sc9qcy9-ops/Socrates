package training

import (
	"fmt"

	"socrates/internal/decipher"
)

// ============================================================
// Evaluator
// ============================================================

// Evaluator evaluates training examples against engine analysis results.
type Evaluator struct {
	engine *decipher.Engine
	kb     *decipher.Engine
}

// NewEvaluator creates a new training evaluator.
func NewEvaluator(engine *decipher.Engine) *Evaluator {
	return &Evaluator{
		engine: engine,
	}
}

// Evaluate runs the engine on all examples and returns the evaluation report.
func (e *Evaluator) Evaluate(examples Examples) *EvaluationResult {
	result := &EvaluationResult{
		TotalExamples:    len(examples),
		ExampleResults:    make([]ExampleResult, 0, len(examples)),
		AllMissedFields:   make([]string, 0),
		AllFalseFields:   make([]string, 0),
	}

	var totalConceptPrecision, totalConceptRecall float64
	var totalFieldPrecision, totalFieldRecall float64
	var qualityAgreements, qualityDisagreements int

	for _, example := range examples {
		exampleResult := e.evaluateExample(example)
		result.ExampleResults = append(result.ExampleResults, exampleResult)

		if exampleResult.Passed {
			result.PassedExamples++
		} else {
			result.FailedExamples++
		}

		totalConceptPrecision += exampleResult.ConceptPrecision
		totalConceptRecall += exampleResult.ConceptRecall
		totalFieldPrecision += exampleResult.FieldPrecision
		totalFieldRecall += exampleResult.FieldRecall

		if example.ExpectedQuality != "" {
			if exampleResult.QualityAgreed {
				qualityAgreements++
			} else {
				qualityDisagreements++
			}
		}

		// Collect missed and false fields
		result.AllMissedFields = append(result.AllMissedFields, exampleResult.MissedExpectedFields...)
		result.AllFalseFields = append(result.AllFalseFields, exampleResult.FalseActivatedFields...)
	}

	// Calculate averages
	if len(examples) > 0 {
		result.AvgConceptPrecision = totalConceptPrecision / float64(len(examples))
		result.AvgConceptRecall = totalConceptRecall / float64(len(examples))
		result.AvgFieldPrecision = totalFieldPrecision / float64(len(examples))
		result.AvgFieldRecall = totalFieldRecall / float64(len(examples))
	}

	result.QualityAgreements = qualityAgreements
	result.QualityDisagreements = qualityDisagreements

	return result
}

// evaluateExample evaluates a single training example.
func (e *Evaluator) evaluateExample(example Example) ExampleResult {
	result := ExampleResult{
		ExampleID:  example.ID,
		Input:     example.Input,
		Passed:    true, // assume passed until proven otherwise
	}

	// Run engine analysis
	reading := e.engine.Analyze(example.Input)

	// Collect activated concepts
	activatedConcepts := collectActivatedConcepts(reading)
	result.ActivatedConcepts = activatedConcepts

	// Collect activated meaning-frequency fields
	activatedFields := collectActivatedFields(reading)
	result.ActivatedFields = activatedFields

	// Calculate concept metrics
	result.ConceptHits, result.ConceptMisses, result.ConceptFalsePos =
		calculateConceptMetrics(example.ExpectedConcepts, activatedConcepts)
	
	if len(example.ExpectedConcepts) > 0 {
		result.ConceptRecall = float64(result.ConceptHits) / float64(len(example.ExpectedConcepts))
	} else {
		result.ConceptRecall = 1.0
	}

	if len(activatedConcepts) > 0 {
		result.ConceptPrecision = float64(result.ConceptHits) / float64(len(activatedConcepts))
	} else if len(example.ExpectedConcepts) == 0 {
		result.ConceptPrecision = 1.0
	} else {
		result.ConceptPrecision = 0.0
	}

	// Calculate field metrics
	result.FieldHits, result.FieldMisses, result.FieldFalsePos =
		calculateFieldMetrics(example.ExpectedFields, activatedFields)
	result.MissedExpectedFields = getMissedFields(example.ExpectedFields, activatedFields)
	result.FalseActivatedFields = getFalseFields(example.ExpectedFields, activatedFields)

	if len(example.ExpectedFields) > 0 {
		result.FieldRecall = float64(result.FieldHits) / float64(len(example.ExpectedFields))
	} else {
		result.FieldRecall = 1.0
	}

	if len(activatedFields) > 0 {
		result.FieldPrecision = float64(result.FieldHits) / float64(len(activatedFields))
	} else if len(example.ExpectedFields) == 0 {
		result.FieldPrecision = 1.0
	} else {
		result.FieldPrecision = 0.0
	}

	// Check quality agreement
	if example.ExpectedQuality != "" {
		actualQuality := getHarmonicQuality(reading)
		result.ExpectedQuality = example.ExpectedQuality
		result.ActualQuality = actualQuality
		result.QualityAgreed = (actualQuality == example.ExpectedQuality)
	}

	// Check evidence paths for matched fields
	result.EvidencePathsValid = checkEvidencePaths(reading, example.ExpectedFields)

	// Determine if example passed
	// Must have high recall (found expected items)
	// Precision is informational but not required for passing in v1
	if result.ConceptRecall < 0.6 && len(example.ExpectedConcepts) > 0 {
		result.Passed = false
	}
	if result.FieldRecall < 0.6 && len(example.ExpectedFields) > 0 {
		result.Passed = false
	}

	return result
}

// collectActivatedConcepts extracts activated concept IDs from a reading.
func collectActivatedConcepts(reading decipher.Reading) []string {
	conceptSet := make(map[string]bool)

	// Collect from passage fields
	for _, field := range reading.PassageFields {
		if field.Concept != "" {
			conceptSet[field.Concept] = true
		}
	}

	// Collect from convergence - these are ActivatedConcept structs, not strings
	for _, concept := range reading.Convergence.ActivatedConcepts {
		conceptSet[concept.Concept] = true
	}

	// Collect from top concepts
	for _, tc := range reading.Convergence.TopConcepts {
		conceptSet[tc.Concept] = true
	}

	// Convert to slice
	concepts := make([]string, 0, len(conceptSet))
	for c := range conceptSet {
		concepts = append(concepts, c)
	}

	return concepts
}

// collectActivatedFields extracts activated meaning-frequency IDs from a reading.
func collectActivatedFields(reading decipher.Reading) []string {
	if reading.HarmonicField == nil {
		return nil
	}

	fieldSet := make(map[string]bool)
	for _, tone := range reading.HarmonicField.Tones {
		if tone.MeaningFrequencyID != "" {
			fieldSet[tone.MeaningFrequencyID] = true
		}
	}

	fields := make([]string, 0, len(fieldSet))
	for f := range fieldSet {
		fields = append(fields, f)
	}

	return fields
}

// calculateConceptMetrics calculates hits, misses, and false positives for concepts.
func calculateConceptMetrics(expected, activated []string) (hits, misses, falsePos int) {
	expectedSet := make(map[string]bool)
	for _, e := range expected {
		expectedSet[e] = true
	}

	activatedSet := make(map[string]bool)
	for _, a := range activated {
		activatedSet[a] = true
	}

	// Hits: expected concepts that were activated
	for _, e := range expected {
		if activatedSet[e] {
			hits++
		} else {
			misses++
		}
	}

	// False positives: activated concepts not expected
	for _, a := range activated {
		if !expectedSet[a] {
			falsePos++
		}
	}

	return hits, misses, falsePos
}

// calculateFieldMetrics calculates hits, misses, and false positives for fields.
func calculateFieldMetrics(expected, activated []string) (hits, misses, falsePos int) {
	expectedSet := make(map[string]bool)
	for _, e := range expected {
		expectedSet[e] = true
	}

	activatedSet := make(map[string]bool)
	for _, a := range activated {
		activatedSet[a] = true
	}

	// Hits: expected fields that were activated
	for _, e := range expected {
		if activatedSet[e] {
			hits++
		} else {
			misses++
		}
	}

	// False positives: activated fields not expected
	for _, a := range activated {
		if !expectedSet[a] {
			falsePos++
		}
	}

	return hits, misses, falsePos
}

// getMissedFields returns fields that were expected but not activated.
func getMissedFields(expected, activated []string) []string {
	activatedSet := make(map[string]bool)
	for _, a := range activated {
		activatedSet[a] = true
	}

	var missed []string
	for _, e := range expected {
		if !activatedSet[e] {
			missed = append(missed, e)
		}
	}

	return missed
}

// getFalseFields returns fields that were activated but not expected.
func getFalseFields(expected, activated []string) []string {
	expectedSet := make(map[string]bool)
	for _, e := range expected {
		expectedSet[e] = true
	}

	var falsePos []string
	for _, a := range activated {
		if !expectedSet[a] {
			falsePos = append(falsePos, a)
		}
	}

	return falsePos
}

// getHarmonicQuality determines the harmonic quality from a reading.
func getHarmonicQuality(reading decipher.Reading) string {
	if reading.HarmonicField == nil {
		return "unknown"
	}

	hf := reading.HarmonicField

	// Quality is determined by consonance vs dissonance
	if hf.Consonance > hf.Dissonance+0.2 {
		return "consonant"
	} else if hf.Dissonance > hf.Consonance+0.2 {
		return "dissonant"
	} else {
		// When consonance and dissonance are similar, it's unresolved tension
		return "unresolved_tension"
	}
}

// checkEvidencePaths verifies that matched fields have evidence paths.
func checkEvidencePaths(reading decipher.Reading, expectedFields []string) bool {
	if reading.HarmonicField == nil {
		return len(expectedFields) == 0
	}

	// Build set of expected fields
	expectedSet := make(map[string]bool)
	for _, f := range expectedFields {
		expectedSet[f] = true
	}

	if len(expectedSet) == 0 {
		return true
	}

	// Check that each expected field has at least one evidence path
	fieldHasEvidence := make(map[string]bool)
	for _, ep := range reading.HarmonicField.EvidencePaths {
		fieldHasEvidence[ep.MeaningFrequencyID] = true
	}

	for _, expected := range expectedFields {
		if !fieldHasEvidence[expected] {
			return false
		}
	}

	return true
}

// FormatResult formats an evaluation result for display.
func FormatResult(result *EvaluationResult) string {
	s := "\nTraining Evaluation Results\n"
	s += "============================\n\n"

	s += fmt.Sprintf("Total examples: %d\n", result.TotalExamples)
	s += fmt.Sprintf("Passed: %d\n", result.PassedExamples)
	s += fmt.Sprintf("Failed: %d\n", result.FailedExamples)
	s += "\n"

	s += "Concept Metrics:\n"
	s += fmt.Sprintf("  Avg Precision: %.2f\n", result.AvgConceptPrecision)
	s += fmt.Sprintf("  Avg Recall: %.2f\n", result.AvgConceptRecall)
	s += "\n"

	s += "Meaning-Frequency Field Metrics:\n"
	s += fmt.Sprintf("  Avg Precision: %.2f\n", result.AvgFieldPrecision)
	s += fmt.Sprintf("  Avg Recall: %.2f\n", result.AvgFieldRecall)
	s += "\n"

	if result.QualityAgreements+result.QualityDisagreements > 0 {
		s += fmt.Sprintf("Quality Agreement: %d/%d\n",
			result.QualityAgreements,
			result.QualityAgreements+result.QualityDisagreements)
		s += "\n"
	}

	if len(result.AllMissedFields) > 0 {
		s += "Missed Expected Fields:\n"
		for _, f := range uniqueStrings(result.AllMissedFields) {
			s += fmt.Sprintf("  - %s\n", f)
		}
		s += "\n"
	}

	if len(result.AllFalseFields) > 0 {
		s += "False Activated Fields:\n"
		for _, f := range uniqueStrings(result.AllFalseFields) {
			s += fmt.Sprintf("  - %s\n", f)
		}
		s += "\n"
	}

	return s
}

// ExampleDetailedResult holds detailed info for formatted output.
type ExampleDetailedResult struct {
	ExampleID      string
	Input          string
	Passed         bool
	ConceptHits    int
	ConceptMisses  int
	ConceptFalsePos int
	ConceptPrec    float64
	ConceptRec     float64
	FieldHits      int
	FieldMisses    int
	FieldFalsePos  int
	FieldPrec      float64
	FieldRec       float64
	ExpectedQuality string
	ActualQuality   string
	QualityAgreed   bool
	EvidenceValid   bool
}

// FormatDetailedResult formats detailed per-example results.
func FormatDetailedResult(result *EvaluationResult, expectedConcepts []string, expectedFields []string) string {
	s := "\nDetailed Example Results\n"
	s += "========================\n\n"

	for i, ex := range result.ExampleResults {
		status := "PASS"
		if !ex.Passed {
			status = "FAIL"
		}

		exID := ex.ExampleID
		if exID == "" {
			exID = fmt.Sprintf("example[%d]", i)
		}

		s += fmt.Sprintf("[%s] %s: %s\n", status, exID, ex.Input)
		
		if len(expectedConcepts) > 0 || len(ex.ActivatedConcepts) > 0 {
			s += fmt.Sprintf("  Concepts: %d hits, %d misses, %d false pos (prec=%.2f, rec=%.2f)\n",
				ex.ConceptHits, ex.ConceptMisses, ex.ConceptFalsePos,
				ex.ConceptPrecision, ex.ConceptRecall)
		}

		if len(expectedFields) > 0 || len(ex.ActivatedFields) > 0 {
			s += fmt.Sprintf("  Fields: %d hits, %d misses, %d false pos (prec=%.2f, rec=%.2f)\n",
				ex.FieldHits, ex.FieldMisses, ex.FieldFalsePos,
				ex.FieldPrecision, ex.FieldRecall)
		}

		if ex.ExpectedQuality != "" {
			agreed := "✗"
			if ex.QualityAgreed {
				agreed = "✓"
			}
			s += fmt.Sprintf("  Quality %s: expected=%s, actual=%s\n",
				agreed, ex.ExpectedQuality, ex.ActualQuality)
		}

		if !ex.EvidencePathsValid {
			s += "  Evidence paths: MISSING for matched fields\n"
		}

		s += "\n"
	}

	return s
}

// FormatDetailedResults is a convenience wrapper that needs examples for context.
func (r *EvaluationResult) FormatDetailedResults(examples []Example) string {
	if len(examples) != len(r.ExampleResults) {
		// Fallback to basic formatting
		s := "\nDetailed Example Results\n"
		s += "========================\n\n"
		for i, ex := range r.ExampleResults {
			status := "PASS"
			if !ex.Passed {
				status = "FAIL"
			}
			s += fmt.Sprintf("[%s] Example %d\n", status, i+1)
			s += fmt.Sprintf("  Input: %s\n", ex.Input)
			s += fmt.Sprintf("  Concepts: %d hits, %d misses, %d false pos\n",
				ex.ConceptHits, ex.ConceptMisses, ex.ConceptFalsePos)
			s += fmt.Sprintf("  Fields: %d hits, %d misses, %d false pos\n",
				ex.FieldHits, ex.FieldMisses, ex.FieldFalsePos)
			s += "\n"
		}
		return s
	}

	// Use detailed format with expected values from examples
	var expectedConcepts []string
	var expectedFields []string
	for _, ex := range examples {
		expectedConcepts = ex.ExpectedConcepts
		expectedFields = ex.ExpectedFields
	}

	return FormatDetailedResult(r, expectedConcepts, expectedFields)
}

// uniqueStrings returns unique strings from a slice.
func uniqueStrings(strs []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range strs {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}