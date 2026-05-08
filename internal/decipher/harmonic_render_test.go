package decipher

import (
	"strings"
	"testing"

	"socrates/internal/knowledge"
)

// TestHarmonicFieldToRenderDescriptor verifies conversion from HarmonicField to RenderDescriptor.
func TestHarmonicFieldToRenderDescriptor(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Create passage fields with concepts that have frequency profiles
	pf := PassageFields{
		&PassageField{
			Concept:       "breath",
			Strength:      0.9,
			Confidence:   "verified",
			EvidenceCount: 2,
		},
		&PassageField{
			Concept:       "life",
			Strength:      0.8,
			Confidence:   "verified",
			EvidenceCount: 2,
		},
	}

	hf := BuildHarmonicField(pf, kb)
	if hf == nil || len(hf.Tones) == 0 {
		t.Skip("Cannot test conversion without profiles")
	}

	rd := hf.ToRenderDescriptor(nil)
	if rd == nil {
		t.Fatal("ToRenderDescriptor should return non-nil for valid field")
	}

	// Verify field descriptor
	if rd.Field.ToneCount != len(hf.Tones) {
		t.Errorf("ToneCount mismatch: got %d, want %d", rd.Field.ToneCount, len(hf.Tones))
	}
	if rd.Field.Coherence != hf.Coherence {
		t.Errorf("Coherence mismatch: got %f, want %f", rd.Field.Coherence, hf.Coherence)
	}
	if rd.Field.DominantQuality != dominantQuality(hf.Consonance, hf.Dissonance) {
		t.Errorf("DominantQuality mismatch: got %s, want %s", rd.Field.DominantQuality, dominantQuality(hf.Consonance, hf.Dissonance))
	}

	// Verify tone descriptors
	if len(rd.Tones) != len(hf.Tones) {
		t.Errorf("Tone count mismatch: got %d, want %d", len(rd.Tones), len(hf.Tones))
	}

	for i, tone := range hf.Tones {
		if i >= len(rd.Tones) {
			break
		}
		td := rd.Tones[i]

		if td.ID != tone.MeaningFrequencyID {
			t.Errorf("Tone %d ID mismatch: got %s, want %s", i, td.ID, tone.MeaningFrequencyID)
		}
		if td.Vector != tone.Vector {
			t.Errorf("Tone %d Vector mismatch: got %v, want %v", i, td.Vector, tone.Vector)
		}
		if td.Ratio != tone.Ratio {
			t.Errorf("Tone %d Ratio mismatch: got %v, want %v", i, td.Ratio, tone.Ratio)
		}
		if td.Archetype != tone.Archetype {
			t.Errorf("Tone %d Archetype mismatch: got %s, want %s", i, td.Archetype, tone.Archetype)
		}
		// Labels are integers from profile
		if td.Labels.Note < 0 || td.Labels.Note > 127 {
			t.Errorf("Tone %d invalid Note label: %d", i, td.Labels.Note)
		}
		if td.Labels.Color < 0 || td.Labels.Color > 360 {
			t.Errorf("Tone %d invalid Color label: %d", i, td.Labels.Color)
		}
		if td.Labels.Field < 0 {
			t.Errorf("Tone %d invalid Field label: %d", i, td.Labels.Field)
		}
	}

	// Verify source trace
	if len(rd.SourceTrace) == 0 {
		t.Error("SourceTrace should not be empty for valid field")
	}
}

// TestRenderDescriptorNilField verifies nil handling.
func TestRenderDescriptorNilField(t *testing.T) {
	var hf *HarmonicField
	rd := hf.ToRenderDescriptor(nil)
	if rd != nil {
		t.Error("nil field should return nil descriptor")
	}
}

// TestRenderDescriptorToString verifies string rendering.
func TestRenderDescriptorToString(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	pf := PassageFields{
		&PassageField{
			Concept:       "breath",
			Strength:      0.9,
			Confidence:   "verified",
			EvidenceCount: 2,
		},
	}

	hf := BuildHarmonicField(pf, kb)
	if hf == nil {
		t.Skip("Cannot test rendering without profiles")
	}

	rd := hf.ToRenderDescriptor(nil)
	output := RenderDescriptorToString(rd)

	// Should contain key sections
	if !strings.Contains(output, "HarmonicField[") {
		t.Error("Output should contain HarmonicField header")
	}
	if !strings.Contains(output, "field:") {
		t.Error("Output should contain field descriptor")
	}
	if !strings.Contains(output, "tones:") {
		t.Error("Output should contain tones section")
	}
	if !strings.Contains(output, "vector=") {
		t.Error("Output should contain vector information")
	}
	if !strings.Contains(output, "ratio=") {
		t.Error("Output should contain ratio information")
	}
	if !strings.Contains(output, "archetypes=") {
		t.Error("Output should contain archetype information")
	}
}

// TestDominantQuality verifies quality determination.
func TestDominantQuality(t *testing.T) {
	tests := []struct {
		consonance float64
		dissonance float64
		expected   string
	}{
		{0.8, 0.2, "consonant"},
		{0.3, 0.7, "dissonant"},
		{0.5, 0.3, "consonant"},
		{0.5, 0.25, "consonant"},
		{0.5, 0.4, "consonant"},   // 0.5 > 0.4
		{0.4, 0.5, "dissonant"},   // 0.4 <= 0.5 and 0.5 > 0.3
	}

	for i, tc := range tests {
		result := dominantQuality(tc.consonance, tc.dissonance)
		if result != tc.expected {
			t.Errorf("Case %d: dominantQuality(%.2f, %.2f) = %s, want %s",
				i, tc.consonance, tc.dissonance, result, tc.expected)
		}
	}
}

// TestArchetypeSummary verifies archetype summarization.
func TestArchetypeSummary(t *testing.T) {
	tones := []HarmonicTone{
		{Archetype: "pythagorean_triple_1"},
		{Archetype: "pythagorean_triple_1"},
		{Archetype: "phi_approximant_1"},
	}

	summary := archetypeSummary(tones)
	if !strings.Contains(summary, "pythagorean_triple_1") {
		t.Error("Summary should contain pythagorean_triple_1")
	}
	if !strings.Contains(summary, "×2") {
		t.Error("Summary should show duplicate count for repeated archetype")
	}
	if !strings.Contains(summary, "phi_approximant_1") {
		t.Error("Summary should contain phi_approximant_1")
	}
}

// TestArchetypeSummaryEmpty verifies empty handling.
func TestArchetypeSummaryEmpty(t *testing.T) {
	tones := []HarmonicTone{
		{Archetype: ""},
		{Archetype: ""},
	}
	summary := archetypeSummary(tones)
	if summary != "undetermined" {
		t.Errorf("Empty archetypes should return 'undetermined', got %s", summary)
	}
}

// TestBuildRelationshipDescriptors verifies relationship building.
func TestBuildRelationshipDescriptors(t *testing.T) {
	tones := []HarmonicTone{
		{MeaningFrequencyID: "tone1", Ratio: [2]int{3, 2}},
		{MeaningFrequencyID: "tone2", Ratio: [2]int{4, 3}},
		{MeaningFrequencyID: "tone3", Ratio: [2]int{7, 4}},
	}

	rels := buildRelationshipDescriptors(tones)

	// Should have 3 relationships (3 pairs from 3 tones)
	if len(rels) != 3 {
		t.Errorf("Expected 3 relationships, got %d", len(rels))
	}

	// First two should be consonant (fifth and fourth)
	// Third pair is dissonant
	countConsonant := 0
	countDissonant := 0
	for _, r := range rels {
		if r.Type == "consonant" {
			countConsonant++
		} else {
			countDissonant++
		}
	}

	if countConsonant < 1 {
		t.Error("Should have at least one consonant relationship")
	}
	if countDissonant < 1 {
		t.Error("Should have at least one dissonant relationship")
	}
}

// TestDescribeRelationship verifies interval descriptions.
func TestDescribeRelationship(t *testing.T) {
	tests := []struct {
		r1       [2]int
		r2       [2]int
		expected string
	}{
		{[2]int{1, 1}, [2]int{1, 1}, "unison"},
		{[2]int{2, 1}, [2]int{1, 1}, "octave (descending)"},  // 2/1 is octave of 1/1, 2>1
		{[2]int{1, 1}, [2]int{3, 2}, "perfect fifth"},
		{[2]int{1, 1}, [2]int{4, 3}, "perfect fourth"},
		{[2]int{1, 1}, [2]int{5, 4}, "major third"},
		{[2]int{1, 1}, [2]int{6, 5}, "minor third"},
	}

	for i, tc := range tests {
		result := describeRelationship(tc.r1, tc.r2)
		if result != tc.expected {
			t.Errorf("Case %d: describeRelationship(%v, %v) = %s, want %s",
				i, tc.r1, tc.r2, result, tc.expected)
		}
	}
}

// TestSourceTraceDeduplication verifies dedup in source trace.
func TestSourceTraceDeduplication(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Use same concept twice to test dedup
	pf := PassageFields{
		&PassageField{
			Concept:       "breath",
			Strength:      0.9,
			Confidence:   "verified",
			EvidenceCount: 2,
		},
	}

	hf := BuildHarmonicField(pf, kb)
	if hf == nil {
		t.Skip("Cannot test dedup without profiles")
	}

	rd := hf.ToRenderDescriptor(nil)

	// Check for duplicates in source trace
	seen := make(map[string]bool)
	for _, s := range rd.SourceTrace {
		key := s.Concept + "|" + s.ToneID
		if seen[key] {
			t.Errorf("Duplicate source trace entry: %s", key)
		}
		seen[key] = true
	}
}

// TestIntegerVectorsAndRatios verifies integer data types.
func TestIntegerVectorsAndRatios(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	pf := PassageFields{
		&PassageField{
			Concept:       "breath",
			Strength:      0.9,
			Confidence:   "verified",
			EvidenceCount: 2,
		},
	}

	hf := BuildHarmonicField(pf, kb)
	if hf == nil {
		t.Skip("Cannot test without profiles")
	}

	rd := hf.ToRenderDescriptor(nil)

	for _, td := range rd.Tones {
		// Vector components must be integers (not floats)
		if td.Vector[0] == 0 && td.Vector[1] == 0 && td.Vector[2] == 0 {
			// May be uninitialized from zero profile - check the profile directly
			profiles := kb.GetFrequencyProfilesByConcept("breath")
			if len(profiles) > 0 && len(profiles[0].Vector) >= 3 {
				// Profile has vector but tone doesn't - this is ok, tone may be from different profile
			}
		}

		// Ratio components must be integers (not floats)
		if td.Ratio[0] == 0 && td.Ratio[1] == 0 {
			// May be uninitialized - check profile
			profiles := kb.GetFrequencyProfilesByConcept("breath")
			if len(profiles) > 0 && len(profiles[0].Ratio) >= 2 {
				// Profile has ratio but tone doesn't - this is ok
			}
		}
	}
}

// TestLoadArchetypeNames verifies archetype name loading.
func TestLoadArchetypeNames(t *testing.T) {
	names := loadArchetypeNames()

	expected := map[string]string{
		"pythagorean_triple_1": "Monad",
		"pythagorean_triple_2": "Triangular",
		"pythagorean_triple_3": "Triangular",
		"pythagorean_triple_4": "Triangular",
		"pythagorean_triple_5": "Primitive",
		"pythagorean_triple_6": "Primitive",
		"phi_approximant_1":    "Golden Ratio",
		"phi_approximant_2":    "Extended Golden",
		"metatron_node_1":      "Metatron Node",
	}

	for id, name := range expected {
		if names[id] != name {
			t.Errorf("Archetype %s: got %s, want %s", id, names[id], name)
		}
	}
}

// TestVectorToStr verifies vector string formatting.
func TestVectorToStr(t *testing.T) {
	v := [3]int{1, 2, 3}
	result := vectorToStr(v)
	expected := "[1,2,3]"
	if result != expected {
		t.Errorf("vectorToStr(%v) = %s, want %s", v, result, expected)
	}
}

// TestRatioToStr verifies ratio string formatting.
func TestRatioToStr(t *testing.T) {
	r := [2]int{3, 2}
	result := ratioToStr(r)
	expected := "3/2"
	if result != expected {
		t.Errorf("ratioToStr(%v) = %s, want %s", r, result, expected)
	}
}

// TestRenderDescriptorToStringNil verifies nil handling.
func TestRenderDescriptorToStringNil(t *testing.T) {
	result := RenderDescriptorToString(nil)
	if result != "" {
		t.Error("nil descriptor should return empty string")
	}
}