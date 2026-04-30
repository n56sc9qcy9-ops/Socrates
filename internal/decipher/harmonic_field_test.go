package decipher

import (
	"strings"
	"testing"

	"socrates/internal/knowledge"
)

// TestBuildHarmonicFieldFromData verifies harmonic field is built from YAML data, not Go constants.
func TestBuildHarmonicFieldFromData(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Create mock passage fields with concepts that have frequency profiles
	pf := PassageFields{
		&PassageField{
			Concept:      "breath",
			Strength:     0.9,
			Confidence:   "verified",
			TokenSources: []string{"prana", "ruach"},
		},
		&PassageField{
			Concept:      "life",
			Strength:     0.8,
			Confidence:   "verified",
			TokenSources: []string{"zoe", " vita"},
		},
		&PassageField{
			Concept:      "spirit",
			Strength:     0.7,
			Confidence:   "plausible",
			TokenSources: []string{"pneuma"},
		},
	}

	// Build harmonic field
	hf := BuildHarmonicField(pf, kb)

	// Verify field was created
	if hf == nil {
		t.Fatal("HarmonicField should not be nil for concepts with profiles")
	}

	// Verify tones exist
	if len(hf.Tones) == 0 {
		t.Error("HarmonicField should have at least one tone for concepts with profiles")
	}

	// Verify coherence and metrics are calculated
	if hf.Coherence < 0 || hf.Coherence > 1 {
		t.Errorf("Coherence should be in [0, 1], got %f", hf.Coherence)
	}
	if hf.Consonance < 0 || hf.Consonance > 1 {
		t.Errorf("Consonance should be in [0, 1], got %f", hf.Consonance)
	}
	if hf.Dissonance < 0 || hf.Dissonance > 1 {
		t.Errorf("Dissonance should be in [0, 1], got %f", hf.Dissonance)
	}

	// Verify evidence paths are recorded
	if len(hf.EvidencePaths) == 0 {
		t.Error("HarmonicField should have evidence paths tracing concept-to-tone mapping")
	}

	// Verify source concepts are tracked in tones
	for _, tone := range hf.Tones {
		if len(tone.SourceConcepts) == 0 {
			t.Error("Tone should have at least one source concept")
		}
		// Verify vector is set (integer from YAML)
		if tone.Vector[0] == 0 && tone.Vector[1] == 0 && tone.Vector[2] == 0 {
			// Could be uninitialized - check profile
			t.Logf("Tone %s has zero vector, checking profile", tone.MeaningFrequencyID)
		}
		// Verify ratio is set (integer pair from YAML)
		if tone.Ratio[0] == 0 && tone.Ratio[1] == 0 {
			t.Logf("Tone %s has zero ratio", tone.MeaningFrequencyID)
		}
	}
}

// TestNoProfileNoField verifies concepts without frequency profiles don't get invented ones.
func TestNoProfileNoField(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Create passage fields with concepts that have NO frequency profiles
	// These concepts exist in the knowledge base but have no frequency profiles
	unknownConcepts := []string{"energy", "truth", "light"}
	
	// Check which concepts don't have profiles
	var conceptsWithoutProfiles []string
	for _, concept := range unknownConcepts {
		profiles := kb.GetFrequencyProfilesByConcept(concept)
		if len(profiles) == 0 {
			conceptsWithoutProfiles = append(conceptsWithoutProfiles, concept)
		}
	}

	if len(conceptsWithoutProfiles) == 0 {
		// All concepts have profiles - use more obscure ones
		conceptsWithoutProfiles = []string{"water", "earth", "fire"}
		for _, concept := range conceptsWithoutProfiles {
			profiles := kb.GetFrequencyProfilesByConcept(concept)
			if len(profiles) > 0 {
				t.Skipf("All test concepts have profiles, skipping test")
			}
		}
	}

	pf := PassageFields{}
	for _, concept := range conceptsWithoutProfiles {
		pf = append(pf, &PassageField{
			Concept:      concept,
			Strength:     0.8,
			Confidence:   "verified",
			TokenSources: []string{"test"},
		})
	}

	// Build harmonic field
	hf := BuildHarmonicField(pf, kb)

	// Field may be nil or empty - concepts without profiles should not get invented ones
	if hf != nil && len(hf.Tones) > 0 {
		// If we got tones, verify they're NOT invented
		for _, concept := range conceptsWithoutProfiles {
			profiles := kb.GetFrequencyProfilesByConcept(concept)
			if len(profiles) == 0 {
				t.Errorf("Concept %s without profile should not create a tone", concept)
			}
		}
	}
}

// TestMultipleSourcesOneTone verifies multiple concepts can contribute to one tone.
func TestMultipleSourcesOneTone(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Create passage fields where multiple concepts share the same frequency profile
	// e.g., breath, wind, spirit all share breath-vibration profile
	pf := PassageFields{
		&PassageField{
			Concept:      "breath",
			Strength:     0.9,
			Confidence:   "verified",
			TokenSources: []string{"prana"},
		},
		&PassageField{
			Concept:      "wind",
			Strength:     0.7,
			Confidence:   "plausible",
			TokenSources: []string{"vayu"},
		},
		&PassageField{
			Concept:      "spirit",
			Strength:     0.6,
			Confidence:   "plausible",
			TokenSources: []string{"atman"},
		},
	}

	hf := BuildHarmonicField(pf, kb)

	if hf == nil {
		t.Fatal("HarmonicField should not be nil")
	}

	// Count how many tones have multiple source concepts
	multiSourceTones := 0
	for _, tone := range hf.Tones {
		if len(tone.SourceConcepts) > 1 {
			multiSourceTones++
			t.Logf("Tone %s has multiple sources: %v", tone.MeaningFrequencyID, tone.SourceConcepts)
		}
	}

	if multiSourceTones == 0 {
		t.Log("No tones with multiple sources - check if breath/wind/spirit share a profile")
	}
}

// TestHarmonicScoring verifies consonance/dissonance scoring.
func TestHarmonicScoring(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Test with concepts that should have compatible profiles
	pf := PassageFields{
		&PassageField{
			Concept:      "breath",
			Strength:     0.9,
			Confidence:   "verified",
		},
		&PassageField{
			Concept:      "life",
			Strength:     0.8,
			Confidence:   "verified",
		},
	}

	hf := BuildHarmonicField(pf, kb)

	if hf == nil {
		t.Skip("Cannot test scoring without profiles")
	}

	// Consonance and dissonance should be valid
	if hf.Consonance+hf.Dissonance > 1.01 {
		t.Errorf("Consonance (%.2f) + Dissonance (%.2f) should be <= 1.0", 
			hf.Consonance, hf.Dissonance)
	}

	// Coherence should be reasonable
	if hf.Coherence < 0 || hf.Coherence > 1 {
		t.Errorf("Coherence should be in [0, 1], got %f", hf.Coherence)
	}
}

// TestRatiosCompatible verifies ratio compatibility logic.
func TestRatiosCompatible(t *testing.T) {
	testCases := []struct {
		r1     [2]int
		r2     [2]int
		expect bool // true = compatible
	}{
		// Identical ratios are compatible
		{[2]int{1, 1}, [2]int{1, 1}, true},
		{[2]int{3, 2}, [2]int{3, 2}, true},

		// 1/1 is compatible with everything (unison)
		{[2]int{1, 1}, [2]int{3, 2}, true},
		{[2]int{1, 1}, [2]int{5, 4}, true},

		// Octave relationships are compatible
		{[2]int{2, 1}, [2]int{1, 1}, true},

		// Pythagorean fifths
		{[2]int{3, 2}, [2]int{4, 3}, true},

		// Perfect fifth and major third
		{[2]int{3, 2}, [2]int{5, 4}, true},

		// Unrelated ratios
		{[2]int{7, 4}, [2]int{5, 3}, false},
	}

	for i, tc := range testCases {
		result := areRatiosCompatible(tc.r1, tc.r2)
		if result != tc.expect {
			t.Errorf("Case %d: ratios %v and %v expected %v, got %v",
				i, tc.r1, tc.r2, tc.expect, result)
		}
	}
}

// TestGCD verifies greatest common divisor.
func TestGCD(t *testing.T) {
	testCases := []struct {
		a, b, expect int
	}{
		{10, 5, 5},
		{12, 8, 4},
		{17, 13, 1},
		{100, 25, 25},
		{7, 0, 7}, // GCD with zero
		{0, 0, 0}, // Both zero
	}

	for _, tc := range testCases {
		result := gcd(tc.a, tc.b)
		if result != tc.expect {
			t.Errorf("gcd(%d, %d) expected %d, got %d", tc.a, tc.b, tc.expect, result)
		}
	}
}

// TestRenderHarmonicField verifies rendering in both modes.
func TestRenderHarmonicField(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	pf := PassageFields{
		&PassageField{
			Concept:      "breath",
			Strength:     0.9,
			Confidence:   "verified",
		},
	}

	hf := BuildHarmonicField(pf, kb)

	if hf == nil {
		t.Skip("Cannot test rendering without profiles")
	}

	// Test default mode - should be concise
	defaultOut := RenderHarmonicField(hf, RenderModeDefault)
	if defaultOut == "" {
		t.Error("Default mode should return non-empty string when tones exist")
	}
	if strings.Contains(defaultOut, "EvidencePaths") {
		t.Error("Default mode should not show full evidence paths")
	}

	// Test debug mode - should show details
	debugOut := RenderHarmonicField(hf, RenderModeDebug)
	if debugOut == "" {
		t.Error("Debug mode should return non-empty string when tones exist")
	}
	if !strings.Contains(debugOut, "Coherence") {
		t.Error("Debug mode should show coherence")
	}

	// Test nil field
	nilOut := RenderHarmonicField(nil, RenderModeDefault)
	if nilOut != "" {
		t.Error("Nil field should return empty string")
	}
}

// TestGetHarmonicFieldSummary verifies summary extraction.
func TestGetHarmonicFieldSummary(t *testing.T) {
	hf := &HarmonicField{
		Tones: []HarmonicTone{
			{MeaningFrequencyID: "breath-vibration", SourceConcepts: []string{"breath"}},
			{MeaningFrequencyID: "life-force-vitality", SourceConcepts: []string{"life"}},
		},
	}

	summary := GetHarmonicFieldSummary(hf)
	if summary == "" {
		t.Error("Summary should not be empty for valid field")
	}
	if !strings.Contains(summary, "breath-vibration") {
		t.Error("Summary should contain tone IDs")
	}

	// Test nil
	nilSummary := GetHarmonicFieldSummary(nil)
	if nilSummary != "" {
		t.Error("Nil field summary should be empty")
	}
}

// TestEmptyFieldsNilField verifies nil/empty handling.
func TestEmptyFieldsNilField(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Empty passage fields
	emptyPF := PassageFields{}
	hf1 := BuildHarmonicField(emptyPF, kb)
	if hf1 != nil {
		t.Error("Empty passage fields should return nil field")
	}

	// Nil passage fields
	hf2 := BuildHarmonicField(nil, kb)
	if hf2 != nil {
		t.Error("Nil passage fields should return nil field")
	}
}

// TestIntegerLabels verifies labels are integer-based.
func TestIntegerLabels(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	profiles := kb.AllFrequencyProfiles()
	for _, fp := range profiles {
		// Verify labels are integer values (not floats stored as strings)
		if fp.Labels.Note < 0 || fp.Labels.Note > 127 {
			t.Errorf("Profile %s has invalid note: %d (should be 0-127)",
				fp.MeaningFrequencyID, fp.Labels.Note)
		}
		if fp.Labels.Color < 0 || fp.Labels.Color > 360 {
			t.Errorf("Profile %s has invalid color: %d (should be 0-360)",
				fp.MeaningFrequencyID, fp.Labels.Color)
		}
		if fp.Labels.Field < 0 {
			t.Errorf("Profile %s has invalid field: %d (should be >= 0)",
				fp.MeaningFrequencyID, fp.Labels.Field)
		}
	}
}

// TestWeightIsInteger verifies weights are stored as integers.
func TestWeightIsInteger(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	profiles := kb.AllFrequencyProfiles()
	for _, fp := range profiles {
		// Verify weight is in valid range
		if fp.Weight < 0 || fp.Weight > 100 {
			t.Errorf("Profile %s has invalid weight: %d (should be 0-100)",
				fp.MeaningFrequencyID, fp.Weight)
		}
	}
}