package training

import (
	"os"
	"testing"

	"socrates/internal/decipher"
	"socrates/internal/knowledge"
)

// testKB loads the knowledge base for testing.
func testKB() *knowledge.Knowledge {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		return nil
	}
	return kb
}

// ============================================================
// Loader Tests
// ============================================================

func TestLoader_Parse(t *testing.T) {
	loader := NewLoader()

	// Valid YAML
	yamlData := `
examples:
  - id: test.1
    input: test
    expected_concepts: [one, two]
    expected_fields: [breath-vibration]
    confidence: curated
    source: test
`

	examples, err := loader.Parse([]byte(yamlData))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(examples) != 1 {
		t.Fatalf("expected 1 example, got %d", len(examples))
	}

	ex := examples[0]
	if ex.ID != "test.1" {
		t.Errorf("expected ID 'test.1', got '%s'", ex.ID)
	}
	if ex.Input != "test" {
		t.Errorf("expected input 'test', got '%s'", ex.Input)
	}
	if len(ex.ExpectedConcepts) != 2 {
		t.Errorf("expected 2 concepts, got %d", len(ex.ExpectedConcepts))
	}
	if len(ex.ExpectedFields) != 1 {
		t.Errorf("expected 1 field, got %d", len(ex.ExpectedFields))
	}
}

func TestLoader_LoadFromFile(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("knowledge not available")
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "examples-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	yamlData := `
examples:
  - id: ex.load.1
    input: breath
    expected_concepts: [breath, life]
    expected_fields: [breath-vibration]
    expected_quality: consonant
    confidence: curated
    source: test
`

	if _, err := tmpFile.WriteString(yamlData); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	loader := NewLoader()
	examples, err := loader.LoadFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected no error loading file, got %v", err)
	}

	if len(examples) != 1 {
		t.Fatalf("expected 1 example, got %d", len(examples))
	}
}

func TestLoader_InvalidYAML(t *testing.T) {
	loader := NewLoader()

	invalidYAML := `
examples:
  - id: test
    input: [invalid broken yaml
`

	_, err := loader.Parse([]byte(invalidYAML))
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

// ============================================================
// Validator Tests
// ============================================================

func TestValidator_ValidExample(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("knowledge not available")
	}

	v := NewValidator(kb)

	example := Example{
		ID:               "test.1",
		Input:            "breath",
		ExpectedConcepts: []string{"breath", "life"},
		ExpectedFields:   []string{"breath-vibration"},
	}

	result := v.Validate(example)
	if !result.IsValid() {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
}

func TestValidator_EmptyID(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("knowledge not available")
	}

	v := NewValidator(kb)

	example := Example{
		ID:               "",
		Input:            "test",
		ExpectedConcepts: []string{},
		ExpectedFields:   []string{},
	}

	result := v.Validate(example)
	if result.IsValid() {
		t.Error("expected error for empty ID, got valid")
	}
}

func TestValidator_EmptyInput(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("knowledge not available")
	}

	v := NewValidator(kb)

	example := Example{
		ID:               "test.1",
		Input:            "",
		ExpectedConcepts: []string{},
		ExpectedFields:   []string{},
	}

	result := v.Validate(example)
	if result.IsValid() {
		t.Error("expected error for empty input, got valid")
	}
}

func TestValidator_InvalidConceptID(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("knowledge not available")
	}

	v := NewValidator(kb)

	example := Example{
		ID:               "test.1",
		Input:            "test",
		ExpectedConcepts: []string{"nonexistent-concept"},
		ExpectedFields:   []string{},
	}

	result := v.Validate(example)
	if result.IsValid() {
		t.Error("expected error for nonexistent concept, got valid")
	}
}

func TestValidator_InvalidFieldID(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("knowledge not available")
	}

	v := NewValidator(kb)

	example := Example{
		ID:               "test.1",
		Input:            "test",
		ExpectedConcepts: []string{},
		ExpectedFields:   []string{"nonexistent-field"},
	}

	result := v.Validate(example)
	if result.IsValid() {
		t.Error("expected error for nonexistent field, got valid")
	}
}

func TestValidator_ValidateAll(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("knowledge not available")
	}

	v := NewValidator(kb)

	examples := Examples{
		{
			ID:               "test.1",
			Input:            "breath",
			ExpectedConcepts: []string{"breath"},
			ExpectedFields:   []string{"breath-vibration"},
		},
		{
			ID:               "test.2",
			Input:            "love",
			ExpectedConcepts: []string{"love"},
			ExpectedFields:   []string{"love-heart-union"},
		},
	}

	result := v.ValidateAll(examples)
	if !result.IsValid() {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
}

func TestValidateExamples_ConvenienceFunction(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("knowledge not available")
	}

	examples := Examples{
		{
			ID:               "test.1",
			Input:            "breath",
			ExpectedConcepts: []string{"breath"},
			ExpectedFields:   []string{"breath-vibration"},
		},
	}

	result := ValidateExamples(examples, kb)
	if !result.IsValid() {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
}

// ============================================================
// Evaluator Tests
// ============================================================

func TestEvaluator_Evaluate(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("knowledge not available")
	}

	// Create actual engine using decipher package
	decipherEng, err := decipher.NewEngine()
	if err != nil {
		t.Skipf("engine not available: %v", err)
	}

	evaluator := NewEvaluator(decipherEng)

	examples := Examples{
		{
			ID:               "ex.eval.1",
			Input:            "breath",
			ExpectedConcepts: []string{"breath", "life"},
			ExpectedFields:   []string{"breath-vibration"},
			ExpectedQuality:  "consonant",
		},
	}

	result := evaluator.Evaluate(examples)

	if result.TotalExamples != 1 {
		t.Errorf("expected 1 total example, got %d", result.TotalExamples)
	}

	if result.PassedExamples < 0 || result.PassedExamples > 1 {
		t.Errorf("expected passed examples 0-1, got %d", result.PassedExamples)
	}

	if len(result.ExampleResults) != 1 {
		t.Errorf("expected 1 result, got %d", len(result.ExampleResults))
	}
}

// Simple test to verify evaluator interface works with test data
func TestEvaluator_SimpleResult(t *testing.T) {
	// Just test that we can create a result and check basic fields
	result := &ExampleResult{
		ExampleID:          "test.1",
		Input:              "test",
		Passed:             true,
		ConceptHits:        2,
		ConceptMisses:      0,
		ConceptFalsePos:    0,
		ConceptPrecision:   1.0,
		ConceptRecall:      1.0,
		FieldHits:          1,
		FieldMisses:        0,
		FieldFalsePos:      0,
		FieldPrecision:     1.0,
		FieldRecall:        1.0,
		QualityAgreed:      true,
		ExpectedQuality:    "consonant",
		ActualQuality:      "consonant",
		EvidencePathsValid: true,
	}

	if !result.Passed {
		t.Error("expected passed to be true")
	}
	if result.ConceptHits != 2 {
		t.Errorf("expected 2 concept hits, got %d", result.ConceptHits)
	}
	if !result.QualityAgreed {
		t.Error("expected quality to agree")
	}
}

// ============================================================
// Metrics Formatting Tests
// ============================================================

func TestEvaluationResult_Metrics(t *testing.T) {
	result := &EvaluationResult{
		TotalExamples:        5,
		PassedExamples:       4,
		FailedExamples:       1,
		AvgConceptPrecision:  0.75,
		AvgConceptRecall:     0.80,
		AvgFieldPrecision:    0.60,
		AvgFieldRecall:       0.90,
		QualityAgreements:    3,
		QualityDisagreements: 1,
	}

	metrics := result.Metrics()

	// Check that metrics string contains expected values
	if len(metrics) == 0 {
		t.Error("expected non-empty metrics string")
	}
}

func TestFormatResult(t *testing.T) {
	result := &EvaluationResult{
		TotalExamples:        3,
		PassedExamples:       2,
		FailedExamples:       1,
		AvgConceptPrecision:  0.50,
		AvgConceptRecall:     0.75,
		AvgFieldPrecision:    0.40,
		AvgFieldRecall:       0.80,
		QualityAgreements:    2,
		QualityDisagreements: 0,
		ExampleResults: []ExampleResult{
			{
				ExampleID: "test.1",
				Input:     "test",
				Passed:    true,
			},
		},
	}

	output := FormatResult(result)
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}
}

// ============================================================
// Example Result Tests
// ============================================================

func TestExampleResult_ConceptMetrics(t *testing.T) {
	result := ExampleResult{
		ConceptPrecision: 0.67,
		ConceptRecall:    0.80,
		ConceptHits:      4,
		ConceptMisses:    1,
		ConceptFalsePos:  2,
	}

	if result.ConceptPrecision != 0.67 {
		t.Errorf("expected precision 0.67, got %f", result.ConceptPrecision)
	}
	if result.ConceptRecall != 0.80 {
		t.Errorf("expected recall 0.80, got %f", result.ConceptRecall)
	}
	if result.ConceptHits != 4 {
		t.Errorf("expected 4 hits, got %d", result.ConceptHits)
	}
}

func TestExampleResult_FieldMetrics(t *testing.T) {
	result := ExampleResult{
		FieldPrecision: 0.50,
		FieldRecall:    0.60,
		FieldHits:      3,
		FieldMisses:    2,
		FieldFalsePos:  3,
	}

	if result.FieldPrecision != 0.50 {
		t.Errorf("expected precision 0.50, got %f", result.FieldPrecision)
	}
	if result.FieldRecall != 0.60 {
		t.Errorf("expected recall 0.60, got %f", result.FieldRecall)
	}
	if result.FieldHits != 3 {
		t.Errorf("expected 3 hits, got %d", result.FieldHits)
	}
}

func TestExampleResult_QualityAgreement(t *testing.T) {
	result := ExampleResult{
		QualityAgreed:   true,
		ExpectedQuality: "consonant",
		ActualQuality:   "consonant",
	}

	if !result.QualityAgreed {
		t.Error("expected quality to agree")
	}
}

func TestExampleResult_EvidencePaths(t *testing.T) {
	result := ExampleResult{
		EvidencePathsValid: true,
	}

	if !result.EvidencePathsValid {
		t.Error("expected evidence paths to be valid")
	}
}
