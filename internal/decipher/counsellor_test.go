package decipher

import (
	"testing"

	"socrates/internal/knowledge"
)

func TestBuildCounsellorField(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Create passage fields with a concept that has transmutation relations
	passageFields := PassageFields{
		{Concept: "resentment", Strength: 0.8, TokenSources: []string{"resentment"}},
		{Concept: "love", Strength: 0.7, TokenSources: []string{"love"}},
	}

	field := BuildCounsellorField(passageFields, kb)

	// Should find resentment source
	if field == nil {
		t.Fatal("expected counsellor field, got nil")
	}

	if len(field.SourceFields) == 0 {
		t.Fatal("expected at least one source field")
	}

	if field.SourceFields[0].Concept != "resentment" {
		t.Errorf("expected source 'resentment', got '%s'", field.SourceFields[0].Concept)
	}

	// Should have suggestions from transmute.yaml
	if len(field.Suggestions) == 0 {
		t.Errorf("expected at least one suggestion for resentment")
	}

	// First suggestion should be resentment -> forgiveness
	if len(field.Suggestions) > 0 && field.Suggestions[0].SourceConcept != "resentment" {
		t.Errorf("expected source 'resentment', got '%s'", field.Suggestions[0].SourceConcept)
	}

	// Should have evidence paths
	if len(field.EvidencePaths) == 0 {
		t.Errorf("expected at least one evidence path")
	}
}

func TestBuildCounsellorFieldNoMatch(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Create passage fields with no concepts that have transmutation relations
	passageFields := PassageFields{
		{Concept: "one", Strength: 0.9, TokenSources: []string{"one"}},
		{Concept: "being", Strength: 0.8, TokenSources: []string{"being"}},
	}

	field := BuildCounsellorField(passageFields, kb)

	// Should be nil when no transmutation matches
	if field != nil {
		t.Errorf("expected nil for non-transmutation concepts, got %v", field)
	}
}

func TestBuildCounsellorFieldWeakStrength(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Create passage fields with strength below threshold
	passageFields := PassageFields{
		{Concept: "resentment", Strength: 0.1, TokenSources: []string{"resentment"}},
	}

	field := BuildCounsellorField(passageFields, kb)

	// Should be nil when strength is below threshold
	if field != nil {
		t.Errorf("expected nil for weak strength, got %v", field)
	}
}

func TestCounsellorFieldKinds(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// resentment: transmutes_to
	passageFields := PassageFields{
		{Concept: "resentment", Strength: 0.8, TokenSources: []string{"resentment"}},
	}
	field := BuildCounsellorField(passageFields, kb)
	if field == nil || len(field.Suggestions) == 0 {
		t.Fatal("expected counsellor field for resentment")
	}
	if field.Suggestions[0].Kind != "transmutes_to" {
		t.Errorf("expected kind 'transmutes_to', got '%s'", field.Suggestions[0].Kind)
	}

	// guilt: softens_through
	passageFields = PassageFields{
		{Concept: "guilt", Strength: 0.8, TokenSources: []string{"guilt"}},
	}
	field = BuildCounsellorField(passageFields, kb)
	if field == nil || len(field.Suggestions) == 0 {
		t.Fatal("expected counsellor field for guilt")
	}
	if field.Suggestions[0].Kind != "softens_through" {
		t.Errorf("expected kind 'softens_through', got '%s'", field.Suggestions[0].Kind)
	}

	// pride: corrects_through
	passageFields = PassageFields{
		{Concept: "pride", Strength: 0.8, TokenSources: []string{"pride"}},
	}
	field = BuildCounsellorField(passageFields, kb)
	if field == nil || len(field.Suggestions) == 0 {
		t.Fatal("expected counsellor field for pride")
	}
	if field.Suggestions[0].Kind != "corrects_through" {
		t.Errorf("expected kind 'corrects_through', got '%s'", field.Suggestions[0].Kind)
	}

	// fear: releases_into
	passageFields = PassageFields{
		{Concept: "fear", Strength: 0.8, TokenSources: []string{"fear"}},
	}
	field = BuildCounsellorField(passageFields, kb)
	if field == nil || len(field.Suggestions) == 0 {
		t.Fatal("expected counsellor field for fear")
	}
	if field.Suggestions[0].Kind != "releases_into" {
		t.Errorf("expected kind 'releases_into', got '%s'", field.Suggestions[0].Kind)
	}

	// shame: grounds_in
	passageFields = PassageFields{
		{Concept: "shame", Strength: 0.8, TokenSources: []string{"shame"}},
	}
	field = BuildCounsellorField(passageFields, kb)
	if field == nil || len(field.Suggestions) == 0 {
		t.Fatal("expected counsellor field for shame")
	}
	if field.Suggestions[0].Kind != "grounds_in" {
		t.Errorf("expected kind 'grounds_in', got '%s'", field.Suggestions[0].Kind)
	}
}

func TestCounsellorFieldMultipleSuggestions(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Check that each transmutation source maps to its correction target
	tests := []struct {
		source      string
		target      string
		minStrength float64
	}{
		{"resentment", "forgiveness", 0.5},
		{"pride", "humility", 0.5},
		{"fear", "trust", 0.5},
		{"guilt", "mercy", 0.4},
		{"shame", "self_acceptance", 0.4},
		{"attachment", "letting_go", 0.4},
		{"control", "surrender", 0.4},
		{"avoidance", "presence", 0.3},
		{"judgment", "mercy", 0.4},
		{"isolation", "connection", 0.5},
		{"grasping", "openness", 0.3},
		{"confusion", "clarity", 0.3},
		{"rigidity", "flexibility", 0.3},
		{"conflict", "harmony", 0.3},
		{"stagnation", "growth", 0.3},
	}

	for _, tc := range tests {
		passageFields := PassageFields{
			{Concept: tc.source, Strength: tc.minStrength + 0.1, TokenSources: []string{tc.source}},
		}
		field := BuildCounsellorField(passageFields, kb)
		if field == nil {
			t.Errorf("expected counsellor field for %s, got nil", tc.source)
			continue
		}
		if len(field.Suggestions) == 0 {
			t.Errorf("expected at least one suggestion for %s", tc.source)
			continue
		}
		found := false
		for _, s := range field.Suggestions {
			if s.TargetConcept == tc.target {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected suggestion %s -> %s, suggestions: %v", tc.source, tc.target, field.Suggestions)
		}
	}
}

func TestTransmutationConfidence(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	passageFields := PassageFields{
		{Concept: "resentment", Strength: 0.8, TokenSources: []string{"resentment"}},
	}
	field := BuildCounsellorField(passageFields, kb)
	if field == nil || len(field.Suggestions) == 0 {
		t.Fatal("expected counsellor field")
	}

	// resentment -> forgiveness should be "plausible"
	if field.Suggestions[0].Confidence != "plausible" {
		t.Errorf("expected confidence 'plausible', got '%s'", field.Suggestions[0].Confidence)
	}
}
