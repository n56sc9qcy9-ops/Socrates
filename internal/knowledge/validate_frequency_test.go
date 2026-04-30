package knowledge

import (
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

	// Should fail because concepts list is empty (but not because concepts don't exist)
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

	// The validator uses `continue` after empty ID, so only one error may be reported
	// This is expected behavior - we just need at least one error
	if len(result.Errors) == 0 {
		t.Error("Should collect at least one error for invalid profile")
	}
}
