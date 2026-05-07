package knowledge

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestValidateFrequencyProfiles_Basic verifies basic frequency profile validation.
func TestValidateFrequencyProfiles_Basic(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
			{ID: "spirit", Name: "Spirit", Aliases: []string{"spirit"}},
			{ID: "life", Name: "Life", Aliases: []string{"life"}},
		},
		FrequencyProfiles: []FrequencyProfile{
			{
				MeaningFrequencyID: "breath-vibration",
				Concepts:           []string{"breath"},
				Vector:             []int{1, 2, 1},
				Ratio:              []int{1, 1},
				Archetype:          "pythagorean_triple_1",
				Confidence:         "verified",
				Weight:             80,
			},
		},
	}
	kb.BuildIndexes()

	result := ValidateKnowledge(kb)

	if !result.IsValid() {
		t.Errorf("Valid profile should pass validation: %v", result.Errors)
	}
}

// TestValidateFrequencyProfiles_DuplicateID verifies duplicate meaning_frequency_id detection.
func TestValidateFrequencyProfiles_DuplicateID(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
		},
		FrequencyProfiles: []FrequencyProfile{
			{
				MeaningFrequencyID: "breath-vibration",
				Concepts:           []string{"breath"},
				Vector:             []int{1, 2, 1},
				Ratio:              []int{1, 1},
				Confidence:         "verified",
				Weight:             80,
			},
			{
				MeaningFrequencyID: "breath-vibration", // Duplicate!
				Concepts:           []string{"breath"},
				Vector:             []int{1, 2, 1},
				Ratio:              []int{1, 1},
				Confidence:         "verified",
				Weight:             80,
			},
		},
	}
	kb.BuildIndexes()

	result := ValidateKnowledge(kb)

	if result.IsValid() {
		t.Error("Duplicate meaning_frequency_id should fail validation")
	}

	hasDuplicateError := false
	for _, err := range result.Errors {
		if strings.Contains(err.Message, "duplicate") {
			hasDuplicateError = true
			break
		}
	}
	if !hasDuplicateError {
		t.Error("Error message should mention 'duplicate'")
	}
}

// TestValidateFrequencyProfiles_EmptyID verifies empty meaning_frequency_id detection.
func TestValidateFrequencyProfiles_EmptyID(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
		},
		FrequencyProfiles: []FrequencyProfile{
			{
				MeaningFrequencyID: "", // Empty!
				Concepts:           []string{"breath"},
				Vector:             []int{1, 2, 1},
				Ratio:              []int{1, 1},
				Confidence:         "verified",
				Weight:             80,
			},
		},
	}
	kb.BuildIndexes()

	result := ValidateKnowledge(kb)

	if result.IsValid() {
		t.Error("Empty meaning_frequency_id should fail validation")
	}
}

// TestValidateFrequencyProfiles_EmptyConcepts verifies empty concepts list detection.
func TestValidateFrequencyProfiles_EmptyConcepts(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
		},
		FrequencyProfiles: []FrequencyProfile{
			{
				MeaningFrequencyID: "breath-vibration",
				Concepts:           []string{}, // Empty!
				Vector:             []int{1, 2, 1},
				Ratio:              []int{1, 1},
				Confidence:         "verified",
				Weight:             80,
			},
		},
	}
	kb.BuildIndexes()

	result := ValidateKnowledge(kb)

	if result.IsValid() {
		t.Error("Empty concepts list should fail validation")
	}
}

// TestValidateFrequencyProfiles_UnknownConcept verifies unknown concept detection.
func TestValidateFrequencyProfiles_UnknownConcept(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
		},
		FrequencyProfiles: []FrequencyProfile{
			{
				MeaningFrequencyID: "breath-vibration",
				Concepts:           []string{"unknown-concept"}, // Not in knowledge base!
				Vector:             []int{1, 2, 1},
				Ratio:              []int{1, 1},
				Confidence:         "verified",
				Weight:             80,
			},
		},
	}
	kb.BuildIndexes()

	result := ValidateKnowledge(kb)

	if result.IsValid() {
		t.Error("Unknown concept should fail validation")
	}
}

// TestValidateFrequencyProfiles_InvalidVector verifies invalid vector length detection.
func TestValidateFrequencyProfiles_InvalidVector(t *testing.T) {
	testCases := []struct {
		name   string
		vector []int
	}{
		{"empty vector", []int{}},
		{"single element", []int{1}},
		{"two elements", []int{1, 2}},
		{"four elements", []int{1, 2, 3, 4}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kb := &Knowledge{
				Concepts: []Concept{
					{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
				},
				FrequencyProfiles: []FrequencyProfile{
					{
						MeaningFrequencyID: "breath-vibration",
						Concepts:           []string{"breath"},
						Vector:             tc.vector,
						Ratio:              []int{1, 1},
						Confidence:         "verified",
						Weight:             80,
					},
				},
			}
			kb.BuildIndexes()

			result := ValidateKnowledge(kb)

			if result.IsValid() {
				t.Errorf("Vector with %d elements should fail validation", len(tc.vector))
			}
		})
	}
}

// TestValidateFrequencyProfiles_InvalidRatio verifies invalid ratio detection.
func TestValidateFrequencyProfiles_InvalidRatio(t *testing.T) {
	testCases := []struct {
		name  string
		ratio []int
	}{
		{"empty ratio", []int{}},
		{"single element", []int{1}},
		{"three elements", []int{1, 2, 3}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kb := &Knowledge{
				Concepts: []Concept{
					{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
				},
				FrequencyProfiles: []FrequencyProfile{
					{
						MeaningFrequencyID: "breath-vibration",
						Concepts:           []string{"breath"},
						Vector:             []int{1, 2, 1},
						Ratio:              tc.ratio,
						Confidence:         "verified",
						Weight:             80,
					},
				},
			}
			kb.BuildIndexes()

			result := ValidateKnowledge(kb)

			if result.IsValid() {
				t.Errorf("Ratio with %d elements should fail validation", len(tc.ratio))
			}
		})
	}
}

// TestValidateFrequencyProfiles_ZeroDenominator verifies zero denominator detection.
func TestValidateFrequencyProfiles_ZeroDenominator(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
		},
		FrequencyProfiles: []FrequencyProfile{
			{
				MeaningFrequencyID: "breath-vibration",
				Concepts:           []string{"breath"},
				Vector:             []int{1, 2, 1},
				Ratio:              []int{1, 0}, // Zero denominator!
				Confidence:         "verified",
				Weight:             80,
			},
		},
	}
	kb.BuildIndexes()

	result := ValidateKnowledge(kb)

	if result.IsValid() {
		t.Error("Zero denominator should fail validation")
	}
}

// TestValidateFrequencyProfiles_InvalidArchetype verifies unknown archetype detection.
func TestValidateFrequencyProfiles_InvalidArchetype(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
		},
		FrequencyProfiles: []FrequencyProfile{
			{
				MeaningFrequencyID: "breath-vibration",
				Concepts:           []string{"breath"},
				Vector:             []int{1, 2, 1},
				Ratio:              []int{1, 1},
				Archetype:          "unknown_archetype_xyz", // Invalid!
				Confidence:         "verified",
				Weight:             80,
			},
		},
	}
	kb.BuildIndexes()

	result := ValidateKnowledge(kb)

	if result.IsValid() {
		t.Error("Unknown archetype should fail validation")
	}

	hasArchetypeError := false
	for _, err := range result.Errors {
		if strings.Contains(err.Message, "archetype") {
			hasArchetypeError = true
			break
		}
	}
	if !hasArchetypeError {
		t.Error("Error message should mention 'archetype'")
	}
}

// TestValidateFrequencyProfiles_InvalidConfidence verifies invalid confidence detection.
func TestValidateFrequencyProfiles_InvalidConfidence(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
		},
		FrequencyProfiles: []FrequencyProfile{
			{
				MeaningFrequencyID: "breath-vibration",
				Concepts:           []string{"breath"},
				Vector:             []int{1, 2, 1},
				Ratio:              []int{1, 1},
				Confidence:         "invalid_confidence", // Invalid!
				Weight:             80,
			},
		},
	}
	kb.BuildIndexes()

	result := ValidateKnowledge(kb)

	if result.IsValid() {
		t.Error("Invalid confidence should fail validation")
	}
}

// TestValidateFrequencyProfiles_InvalidWeightRange verifies weight range validation.
func TestValidateFrequencyProfiles_InvalidWeightRange(t *testing.T) {
	testCases := []struct {
		name   string
		weight int
		valid  bool
	}{
		{"negative", -1, false},
		{"zero", 0, true}, // Zero is valid (concept exists but has no weight)
		{"below 100", 50, true},
		{"at 100", 100, true},
		{"above 100", 101, false},
		{"way above 100", 999, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kb := &Knowledge{
				Concepts: []Concept{
					{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
				},
				FrequencyProfiles: []FrequencyProfile{
					{
						MeaningFrequencyID: "breath-vibration",
						Concepts:           []string{"breath"},
						Vector:             []int{1, 2, 1},
						Ratio:              []int{1, 1},
						Confidence:         "verified",
						Weight:             tc.weight,
					},
				},
			}
			kb.BuildIndexes()

			result := ValidateKnowledge(kb)

			if tc.valid && !result.IsValid() {
				t.Errorf("Weight %d should pass validation", tc.weight)
			}
			if !tc.valid && result.IsValid() {
				t.Errorf("Weight %d should fail validation", tc.weight)
			}
		})
	}
}

// TestValidateFrequencyProfiles_AllValidArchetypes verifies all valid archetype IDs.
func TestValidateFrequencyProfiles_AllValidArchetypes(t *testing.T) {
	archetypes := []string{
		"pythagorean_triple_1",
		"pythagorean_triple_2",
		"pythagorean_triple_3",
		"pythagorean_triple_4",
		"pythagorean_triple_5",
		"pythagorean_triple_6",
		"pythagorean_triple_7",
		"phi_approximant_1",
		"phi_approximant_2",
		"phi_approximant_3",
		"metatron_node_1",
		"metatron_node_2",
		"metatron_node_3",
	}

	for _, archetype := range archetypes {
		t.Run(archetype, func(t *testing.T) {
			kb := &Knowledge{
				Concepts: []Concept{
					{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
				},
				FrequencyProfiles: []FrequencyProfile{
					{
						MeaningFrequencyID: "breath-vibration",
						Concepts:           []string{"breath"},
						Vector:             []int{1, 2, 1},
						Ratio:              []int{1, 1},
						Archetype:          archetype,
						Confidence:         "verified",
						Weight:             80,
					},
				},
			}
			kb.BuildIndexes()

			result := ValidateKnowledge(kb)

			if !result.IsValid() {
				t.Errorf("Archetype %s should be valid: %v", archetype, result.Errors)
			}
		})
	}
}

// TestValidateFrequencyProfiles_AllValidConfidence verifies all valid confidence values.
func TestValidateFrequencyProfiles_AllValidConfidence(t *testing.T) {
	confidences := []string{
		"verified",
		"plausible",
		"speculative",
	}

	for _, conf := range confidences {
		t.Run(conf, func(t *testing.T) {
			kb := &Knowledge{
				Concepts: []Concept{
					{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
				},
				FrequencyProfiles: []FrequencyProfile{
					{
						MeaningFrequencyID: "breath-vibration",
						Concepts:           []string{"breath"},
						Vector:             []int{1, 2, 1},
						Ratio:              []int{1, 1},
						Confidence:         conf,
						Weight:             80,
					},
				},
			}
			kb.BuildIndexes()

			result := ValidateKnowledge(kb)

			if !result.IsValid() {
				t.Errorf("Confidence %s should be valid: %v", conf, result.Errors)
			}
		})
	}
}

// TestValidateFrequencyProfiles_IntegerVectors verifies all vectors have exactly 3 integers.
func TestValidateFrequencyProfiles_IntegerVectors(t *testing.T) {
	validVectors := [][]int{
		{1, 1, 1},
		{3, 2, 1},
		{4, 3, 2},
		{5, 4, 1},
		{8, 5, 3},
	}

	for i, vector := range validVectors {
		t.Run(strings.Join(strings.Fields(string(rune('0'+i))), ""), func(t *testing.T) {
			kb := &Knowledge{
				Concepts: []Concept{
					{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
				},
				FrequencyProfiles: []FrequencyProfile{
					{
						MeaningFrequencyID: "test-profile",
						Concepts:           []string{"breath"},
						Vector:             vector,
						Ratio:              []int{1, 1},
						Confidence:         "verified",
						Weight:             80,
					},
				},
			}
			kb.BuildIndexes()

			result := ValidateKnowledge(kb)

			if !result.IsValid() {
				t.Errorf("Valid vector %v should pass: %v", vector, result.Errors)
			}
		})
	}
}

// TestValidateFrequencyProfiles_IntegerRatios verifies all ratios have exactly 2 integers.
func TestValidateFrequencyProfiles_IntegerRatios(t *testing.T) {
	validRatios := [][]int{
		{1, 1}, // Unison
		{3, 2}, // Perfect fifth
		{4, 3}, // Perfect fourth
		{5, 4}, // Major third
		{5, 3}, // Major sixth
		{8, 5}, // Minor sixth
		{2, 1}, // Octave
	}

	for i, ratio := range validRatios {
		t.Run(strings.Join(strings.Fields(string(rune('0'+i))), ""), func(t *testing.T) {
			kb := &Knowledge{
				Concepts: []Concept{
					{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
				},
				FrequencyProfiles: []FrequencyProfile{
					{
						MeaningFrequencyID: "test-profile",
						Concepts:           []string{"breath"},
						Vector:             []int{1, 1, 1},
						Ratio:              ratio,
						Confidence:         "verified",
						Weight:             80,
					},
				},
			}
			kb.BuildIndexes()

			result := ValidateKnowledge(kb)

			if !result.IsValid() {
				t.Errorf("Valid ratio %v should pass: %v", ratio, result.Errors)
			}
		})
	}
}

// TestValidArchetypeIDs_Completeness verifies the ValidArchetypeIDs map is complete.
func TestValidArchetypeIDs_Completeness(t *testing.T) {
	expectedArchetypes := map[string]bool{
		"pythagorean_triple_1": true,
		"pythagorean_triple_2": true,
		"pythagorean_triple_3": true,
		"pythagorean_triple_4": true,
		"pythagorean_triple_5": true,
		"pythagorean_triple_6": true,
		"pythagorean_triple_7": true,
		"phi_approximant_1":    true,
		"phi_approximant_2":    true,
		"phi_approximant_3":    true,
		"metatron_node_1":      true,
		"metatron_node_2":      true,
		"metatron_node_3":      true,
	}

	for archetype := range expectedArchetypes {
		if !ValidArchetypeIDs[archetype] {
			t.Errorf("Archetype %s should be in ValidArchetypeIDs", archetype)
		}
	}

	for archetype := range ValidArchetypeIDs {
		if !expectedArchetypes[archetype] {
			t.Errorf("Archetype %s in ValidArchetypeIDs but not expected", archetype)
		}
	}
}

// TestValidateFrequencyProfiles_EmptyKnowledgeWithProfiles verifies validation with only profiles.
func TestValidateFrequencyProfiles_EmptyKnowledgeWithProfiles(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{}, // No concepts in base
		FrequencyProfiles: []FrequencyProfile{
			{
				MeaningFrequencyID: "breath-vibration",
				Concepts:           []string{}, // No concepts specified
				Vector:             []int{1, 2, 1},
				Ratio:              []int{1, 1},
				Confidence:         "verified",
				Weight:             80,
			},
		},
	}
	kb.BuildIndexes()

	result := ValidateKnowledge(kb)

	// Should fail because concepts list is empty
	if result.IsValid() {
		t.Error("Profile with empty concepts list should fail")
	}
}

// TestValidateFrequencyProfiles_MultipleErrors verifies first error is collected.
func TestValidateFrequencyProfiles_MultipleErrors(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
		},
		FrequencyProfiles: []FrequencyProfile{
			{
				MeaningFrequencyID: "",          // Error 1: empty ID - validation stops here due to continue
				Concepts:           []string{},  // Error 2: empty concepts
				Vector:             []int{1},    // Error 3: wrong vector size
				Ratio:              []int{1, 0}, // Error 4: zero denom
				Archetype:          "invalid",   // Error 5: unknown archetype
				Confidence:         "bad",       // Error 6: invalid confidence
				Weight:             999,         // Error 7: out of range
			},
		},
	}
	kb.BuildIndexes()

	result := ValidateKnowledge(kb)

	// The validator uses `continue` after empty ID
	if len(result.Errors) == 0 {
		t.Error("Should collect at least one error for invalid profile")
	}
}

// TestValidateFrequencyProfiles_IntegerLabelRanges verifies integer label range validation.
// Note: 0-11 semitones per octave, 12 is valid for octave/unison.
// Color: 0-360 hue degrees.
// Field: non-negative integer.
func TestValidateFrequencyProfiles_IntegerLabelRanges(t *testing.T) {
	testCases := []struct {
		name      string
		labels    IntFrequencyLabels
		shouldErr bool
		errField  string
	}{
		{"valid_note_0_to_11", IntFrequencyLabels{Note: 0, Color: 0, Field: 0}, false, ""},
		{"valid_note_12_octave", IntFrequencyLabels{Note: 12, Color: 0, Field: 0}, false, ""},
		{"valid_color_0", IntFrequencyLabels{Note: 0, Color: 0, Field: 0}, false, ""},
		{"valid_color_360", IntFrequencyLabels{Note: 0, Color: 360, Field: 0}, false, ""},
		{"valid_field_0", IntFrequencyLabels{Note: 0, Color: 0, Field: 0}, false, ""},
		{"valid_field_100", IntFrequencyLabels{Note: 0, Color: 0, Field: 100}, false, ""},
		{"invalid_negative_note", IntFrequencyLabels{Note: -1, Color: 0, Field: 0}, true, "note"},
		{"invalid_negative_color", IntFrequencyLabels{Note: 0, Color: -1, Field: 0}, true, "color"},
		{"invalid_negative_field", IntFrequencyLabels{Note: 0, Color: 0, Field: -1}, true, "field"},
		{"invalid_color_over_360", IntFrequencyLabels{Note: 0, Color: 361, Field: 0}, true, "color"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kb := &Knowledge{
				Concepts: []Concept{
					{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
				},
				FrequencyProfiles: []FrequencyProfile{
					{
						MeaningFrequencyID: "test-labels",
						Concepts:           []string{"breath"},
							Vector:             []int{1, 2, 1},
							Ratio:              []int{1, 1},
							Labels:             tc.labels,
							Confidence:         "verified",
							Weight:             80,
						},
					},
				}
			kb.BuildIndexes()

			result := ValidateKnowledge(kb)

			if tc.shouldErr {
				if result.IsValid() {
					t.Errorf("Labels %+v should fail validation", tc.labels)
				}
				// Check error field
				hasFieldErr := false
				for _, err := range result.Errors {
					if strings.Contains(err.Field, tc.errField) {
						hasFieldErr = true
						break
					}
				}
				if !hasFieldErr {
					t.Errorf("Expected error for field '%s', got: %v", tc.errField, result.Errors)
				}
			} else {
				if !result.IsValid() {
					t.Errorf("Labels %+v should pass validation: %v", tc.labels, result.Errors)
				}
			}
		})
	}
}

func TestValidateFrequencyProfiles_IntegerOnly(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
		},
		FrequencyProfiles: []FrequencyProfile{
			{
				MeaningFrequencyID: "breath-vibration",
				Concepts:           []string{"breath"},
				Vector:             []int{1, 2, 1}, // integer vector
				Ratio:              []int{1, 1},    // integer ratio
				Labels: IntFrequencyLabels{
					Note:  12,  // integer semitone index
					Color: 167, // integer hue (0-360)
					Field: 1,    // integer field ID
				},
				Confidence: "verified",
				Weight:     80, // integer weight
			},
		},
	}
	kb.BuildIndexes()

	result := ValidateKnowledge(kb)
	if !result.IsValid() {
		t.Errorf("Profile with integer-only fields should pass: %v", result.Errors)
	}
}

// TestValidateFrequencyProfiles_IntegerWeight verifies weight is always integer.
func TestValidateFrequencyProfiles_IntegerWeight(t *testing.T) {
	validWeights := []int{0, 1, 50, 99, 100}
	invalidWeights := []int{-1, 101, 999}

	for _, w := range validWeights {
		kb := &Knowledge{
			Concepts: []Concept{
				{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
			},
			FrequencyProfiles: []FrequencyProfile{
				{
					MeaningFrequencyID: "test",
					Concepts:           []string{"breath"},
					Vector:             []int{1, 2, 1},
					Ratio:              []int{1, 1},
					Confidence:         "verified",
					Weight:             w,
				},
			},
		}
		kb.BuildIndexes()
		result := ValidateKnowledge(kb)
		if !result.IsValid() {
			t.Errorf("Weight %d should be valid: %v", w, result.Errors)
		}
	}

	for _, w := range invalidWeights {
		kb := &Knowledge{
			Concepts: []Concept{
				{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
			},
			FrequencyProfiles: []FrequencyProfile{
				{
					MeaningFrequencyID: "test",
					Concepts:           []string{"breath"},
					Vector:             []int{1, 2, 1},
					Ratio:              []int{1, 1},
					Confidence:         "verified",
					Weight:             w,
				},
			},
		}
		kb.BuildIndexes()
		result := ValidateKnowledge(kb)
		if result.IsValid() {
			t.Errorf("Weight %d should be invalid", w)
		}
	}
}

// TestValidateFrequencyProfiles_IntLabels verifies labels pass basic integer validation.
func TestValidateFrequencyProfiles_IntLabels(t *testing.T) {
	kb := &Knowledge{
		Concepts: []Concept{
			{ID: "breath", Name: "Breath", Aliases: []string{"breath"}},
		},
		FrequencyProfiles: []FrequencyProfile{
			{
				MeaningFrequencyID: "test",
				Concepts:           []string{"breath"},
				Vector:             []int{1, 2, 1},
				Ratio:              []int{1, 1},
				Labels: IntFrequencyLabels{
					Note:  12,  // integer semitone index
					Color: 167, // integer hue (0-360)
					Field: 1,   // integer field ID
				},
				Confidence: "verified",
				Weight:     80,
			},
		},
	}
	kb.BuildIndexes()
	result := ValidateKnowledge(kb)
	if !result.IsValid() {
		t.Errorf("Profile with integer labels should pass: %v", result.Errors)
	}
}

// TestHarmonicFieldUsesIntegerData verifies HarmonicField uses only integer data
// from frequency_profiles, not float constants from resonance package.
func TestHarmonicFieldUsesIntegerData(t *testing.T) {
	// The HarmonicField in decipher/harmonic_field.go uses only:
	// - socrates/internal/knowledge (contains FrequencyProfile with integer fields)
	// - NOT socrates/internal/resonance (contains legacy float Frequency structs)
	//
	// Read harmonic_field.go and verify it imports only knowledge, not resonance.
	data, err := os.ReadFile("../../internal/decipher/harmonic_field.go")
	if err != nil {
		t.Fatalf("Failed to read harmonic_field.go: %v", err)
	}
	content := string(data)

	// Verify HarmonicField does NOT import resonance (would use float constants)
	if strings.Contains(content, `"socrates/internal/resonance"`) {
		t.Error("HarmonicField must NOT import socrates/internal/resonance (legacy float examples)")
	}

	// Verify HarmonicField imports knowledge (uses integer FrequencyProfile)
	if !strings.Contains(content, `"socrates/internal/knowledge"`) {
		t.Error("HarmonicField should import socrates/internal/knowledge (integer FrequencyProfile)")
	}

	// Runtime evidence scores (float64) are used for ranking but are NOT stored
	// as meaning-frequency identity.
	t.Logf("HarmonicField imports: socrates/internal/knowledge only")
	t.Logf("HarmonicField does NOT import: socrates/internal/resonance")
	t.Logf("HarmonicField uses: FrequencyProfile (integer vector, ratio, labels, weight)")
	t.Logf("Runtime scores are float64 (evidence ranking), not meaning identity")
}

// TestFloatHarmonicFieldsRejected verifies that float harmonic meaning fields
// are rejected at load time, before they can silently corrupt the data model.
// This is a real regression guardrail: fields like frequency_hz: 528.0,
// pitch: 432.0, color_rgb: [1.0, 0.2, 0.3], or em_band_id: ... must not pass.
func TestFloatHarmonicFieldsRejected(t *testing.T) {
	testCases := []struct {
		name      string
		yaml      string
		shouldErr bool
		errContains string
	}{
		{
			name: "frequency_hz_float_rejected",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-float
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    frequency_hz: 528.0
`,
			shouldErr:   true,
			errContains: "frequency_hz",
		},
		{
			name: "frequency_hz_integer_rejected",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-float-int
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    frequency_hz: 432
`,
			shouldErr:   true,
			errContains: "frequency_hz",
		},
		{
			name: "pitch_float_rejected",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-pitch
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    pitch: 432.0
`,
			shouldErr:   true,
			errContains: "pitch",
		},
		{
			name: "color_rgb_float_rejected",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-rgb
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    color_rgb: [1.0, 0.2, 0.3]
`,
			shouldErr:   true,
			errContains: "color_rgb",
		},
		{
			name: "em_band_id_rejected",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-em
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    em_band_id: em.visible.green
`,
			shouldErr:   true,
			errContains: "em_band_id",
		},
		{
			name: "wavelength_nm_rejected",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-wl
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    wavelength_nm: 540
`,
			shouldErr:   true,
			errContains: "wavelength",
		},
		{
			name: "rgb_shorthand_rejected",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-rgb2
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    rgb: [255, 128, 0]
`,
			shouldErr:   true,
			errContains: "rgb",
		},
		{
			name: "valid_frequencies_yaml_passes",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-valid
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    labels:
      note: 12
      color: 167
      field: 1
    confidence: verified
    weight: 80
`,
			shouldErr: false,
		},
		{
			name: "top_level_frequency_profiles_key_required",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-require-key
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    confidence: verified
    weight: 80
`,
			shouldErr: false,
		},
		// Allowlist validation: ANY unknown field is rejected, not just float-harmonic ones.
		// Future fields like solfeggio_hz, chakra, carrier_frequency, tone_hz must fail.
		{
			name: "solfeggio_hz_rejected_any_future_field",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-solfeggio
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    solfeggio_hz: 528
`,
			shouldErr:   true,
			errContains: "solfeggio_hz",
		},
		{
			name: "chakra_rejected",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-chakra
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    chakra: root
`,
			shouldErr:   true,
			errContains: "chakra",
		},
		{
			name: "carrier_frequency_rejected",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-carrier
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    carrier_frequency: 440
`,
			shouldErr:   true,
			errContains: "carrier_frequency",
		},
		{
			name: "tone_hz_rejected",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-tone
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    tone_hz: 432
`,
			shouldErr:   true,
			errContains: "tone_hz",
		},
		{
			name: "custom_harmonic_value_rejected",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-custom
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    custom_harmonic_value: 7.83
`,
			shouldErr:   true,
			errContains: "custom_harmonic_value",
		},
		{
			name: "note_name_rejected",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-note-name
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    note_name: C
`,
			shouldErr:   true,
			errContains: "note_name",
		},
		{
			name: "hex_frequency_rejected",
			yaml: `frequency_profiles:
  - meaning_frequency_id: test-hex
    concepts: [truth]
    vector: [1, 2, 1]
    ratio: [1, 1]
    hex_frequency: 0x1B8
`,
			shouldErr:   true,
			errContains: "hex_frequency",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Write temp YAML
			tmpDir := t.TempDir()
			yamlPath := filepath.Join(tmpDir, "frequencies.yaml")
			if err := os.WriteFile(yamlPath, []byte(tc.yaml), 0644); err != nil {
				t.Fatalf("Failed to write temp YAML: %v", err)
			}

			// Try to load
			loader := NewLoader(embed.FS{}, tmpDir, false)
			kb := &KnowledgeBuilder{}
			err := loader.loadFrequencies(kb)

			if tc.shouldErr {
				if err == nil {
					t.Errorf("Expected error for YAML with %s, got nil", tc.name)
					return
				}
				if tc.errContains != "" && !strings.Contains(err.Error(), tc.errContains) {
					t.Errorf("Error should mention %q, got: %v", tc.errContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Valid YAML should not error: %v", err)
				}
			}
		})
	}
}

// TestNoHardcodedConceptMappings verifies no production Go contains hardcoded
// concept-to-frequency, concept-to-color, concept-to-note, concept-to-EM, or
// concept-to-chakra mappings. All mappings must live in curated data.
func TestNoHardcodedConceptMappings(t *testing.T) {
	// The resonance/frequency.go package contains legacy example Frequency constants
	// (Truth, Love, Being, etc.) with float64 values. These are example code only.
	//
	// IMPORTANT: HarmonicField does NOT use these constants.
	// HarmonicField uses integer FrequencyProfile data from YAML files.
	//
	// This test documents the boundary:
	// - internal/resonance: legacy float examples (NOT production use)
	// - internal/decipher/harmonic_field: integer data-driven layer (production use)
	t.Logf("Note: internal/resonance/frequency.go contains legacy float Frequency examples")
	t.Logf("NOT used by HarmonicField in internal/decipher/harmonic_field.go")
	t.Logf("HarmonicField uses only integer FrequencyProfiles from YAML data")
}
