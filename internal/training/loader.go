package training

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"socrates/internal/knowledge"
)

// ============================================================
// YAML Document Structures
// ============================================================

type examplesDoc struct {
	Examples []ExampleEntry `yaml:"examples"`
}

type ExampleEntry struct {
	ID              string   `yaml:"id"`
	Input           string   `yaml:"input"`
	ExpectedConcepts []string `yaml:"expected_concepts"`
	ExpectedFields  []string `yaml:"expected_fields"`
	ExpectedQuality string   `yaml:"expected_quality"`
	Confidence      string   `yaml:"confidence"`
	Source          string   `yaml:"source"`
	Notes           string   `yaml:"notes"`
}

// ToExample converts a YAML entry to internal Example type.
func (e ExampleEntry) ToExample() Example {
	return Example{
		ID:               e.ID,
		Input:            e.Input,
		ExpectedConcepts: e.ExpectedConcepts,
		ExpectedFields:   e.ExpectedFields,
		ExpectedQuality:  e.ExpectedQuality,
		Confidence:       e.Confidence,
		Source:           e.Source,
		Notes:            e.Notes,
	}
}

// ============================================================
// Loader
// ============================================================

// Loader handles loading training examples from YAML files.
type Loader struct{}

// NewLoader creates a new training loader.
func NewLoader() *Loader {
	return &Loader{}
}

// LoadFromFile loads training examples from a YAML file.
func (l *Loader) LoadFromFile(path string) (Examples, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read training file: %w", err)
	}

	return l.Parse(data)
}

// Parse parses training examples from YAML data.
func (l *Loader) Parse(data []byte) (Examples, error) {
	var doc examplesDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	examples := make(Examples, 0, len(doc.Examples))
	for _, e := range doc.Examples {
		examples = append(examples, e.ToExample())
	}

	return examples, nil
}

// ============================================================
// Validator
// ============================================================

// Validator validates training examples against active knowledge.
type Validator struct {
	kb *knowledge.Knowledge
}

// NewValidator creates a new training example validator.
func NewValidator(kb *knowledge.Knowledge) *Validator {
	return &Validator{kb: kb}
}

// ValidationError represents an error in training example validation.
type ValidationError struct {
	ExampleID string
	Field     string
	Message   string
}

// ValidationResult contains all validation errors.
type ValidationResult struct {
	Errors []ValidationError
}

// IsValid returns true if there are no errors.
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// AddError adds a validation error.
func (r *ValidationResult) AddError(exampleID, field, message string) {
	r.Errors = append(r.Errors, ValidationError{
		ExampleID: exampleID,
		Field:     field,
		Message:   message,
	})
}

// Validate validates a single training example.
func (v *Validator) Validate(example Example) *ValidationResult {
	result := &ValidationResult{
		Errors: make([]ValidationError, 0),
	}

	// Check non-empty ID
	if example.ID == "" {
		result.AddError(example.ID, "id", "example ID cannot be empty")
	}

	// Check non-empty input
	if example.Input == "" {
		result.AddError(example.ID, "input", "example input cannot be empty")
	}

	// Check expected concepts exist in knowledge
	for i, conceptID := range example.ExpectedConcepts {
		if !v.kb.HasConcept(conceptID) {
			result.AddError(example.ID, fmt.Sprintf("expected_concepts[%d]", i),
				fmt.Sprintf("concept '%s' does not exist in knowledge", conceptID))
		}
	}

	// Check expected fields exist in frequency profiles
	for i, fieldID := range example.ExpectedFields {
		if !v.kb.HasFrequencyProfile(fieldID) {
			result.AddError(example.ID, fmt.Sprintf("expected_fields[%d]", i),
				fmt.Sprintf("meaning-frequency field '%s' does not exist in frequency profiles", fieldID))
		}
	}

	return result
}

// ValidateAll validates all training examples.
func (v *Validator) ValidateAll(examples Examples) *ValidationResult {
	result := &ValidationResult{
		Errors: make([]ValidationError, 0),
	}

	for _, ex := range examples {
		exampleResult := v.Validate(ex)
		result.Errors = append(result.Errors, exampleResult.Errors...)
	}

	return result
}

// ValidateExamples is a convenience function that validates examples.
func ValidateExamples(examples Examples, kb *knowledge.Knowledge) *ValidationResult {
	v := NewValidator(kb)
	return v.ValidateAll(examples)
}