package knowledge

import (
	"testing"
)

func TestLoadFromEmbed(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Check concepts loaded
	if len(kb.Concepts) == 0 {
		t.Error("no concepts loaded")
	}

	// Check forms loaded
	if len(kb.Forms) == 0 {
		t.Error("no forms loaded")
	}

	// Check relations loaded
	if len(kb.Relations) == 0 {
		t.Error("no relations loaded")
	}

	// Check script words loaded
	if len(kb.ScriptWords) == 0 {
		t.Error("no script words loaded")
	}
}

func TestKnowledgeIndexes(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Test concept by ID
	concept, ok := kb.GetConceptByID("breath")
	if !ok {
		t.Error("breath concept not found")
	}
	if concept.Name != "breath" {
		t.Errorf("expected name 'breath', got %q", concept.Name)
	}

	// Test concept by alias
	concept, ok = kb.GetConceptByAlias("prana")
	if !ok {
		t.Error("prana alias not found")
	}

	// Test forms by text
	forms := kb.GetFormsByText("skall")
	if len(forms) == 0 {
		t.Error("no forms for 'skall'")
	}

	// Test relations from
	relations := kb.GetRelationsFrom("breath")
	if len(relations) == 0 {
		t.Error("no relations from breath")
	}

	// Test neighbor concepts
	neighbors := kb.GetNeighborConcepts("life")
	if len(neighbors) == 0 {
		t.Error("no neighbors for life concept")
	}
}

func TestScriptWordsByScript(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Test Hebrew script words
	hebrewWords := kb.GetScriptWordsByScript("hebrew")
	if len(hebrewWords) == 0 {
		t.Error("no Hebrew script words found")
	}

	// Verify Hebrew script word structure
	for _, w := range hebrewWords {
		if w.Script != "hebrew" {
			t.Errorf("expected script 'hebrew', got %q", w.Script)
		}
		if len(w.Runes) == 0 {
			t.Error("Hebrew word has no runes")
		}
	}
}

func TestFormsByConcept(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Test forms by concept
	forms := kb.GetFormsByConcept("breath")
	if len(forms) == 0 {
		t.Error("no forms for breath concept")
	}

	// Verify at least one form has high weight
	var hasHighWeight bool
	for _, f := range forms {
		if f.Weight >= 0.6 {
			hasHighWeight = true
			break
		}
	}
	if !hasHighWeight {
		t.Error("no forms with high weight for breath concept")
	}
}

func TestRelationsByType(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Check that relations have proper structure
	for _, r := range kb.Relations {
		if r.From == "" || r.To == "" {
			t.Errorf("relation missing from/to: %+v", r)
		}
		if r.Weight <= 0 || r.Weight > 1 {
			t.Errorf("relation %s->%s has invalid weight: %f", r.From, r.To, r.Weight)
		}
	}
}

func TestConfidenceValues(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	validConfidences := map[string]bool{
		"verified":   true,
		"plausible": true,
		"speculative": true,
	}

	// Check forms have valid confidence
	for _, f := range kb.Forms {
		if !validConfidences[f.Confidence] {
			t.Errorf("form %s has invalid confidence: %q", f.Form, f.Confidence)
		}
	}

	// Check script words have valid confidence
	for _, w := range kb.ScriptWords {
		if !validConfidences[w.Confidence] {
			t.Errorf("script word %s has invalid confidence: %q", w.Word, w.Confidence)
		}
	}
}

func TestConceptNeighborsExist(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Check that all neighbor references point to existing concepts
	for _, c := range kb.Concepts {
		for _, neighborID := range c.Neighbors {
			if _, ok := kb.conceptsByID[neighborID]; !ok {
				t.Errorf("concept %s has non-existent neighbor: %s", c.ID, neighborID)
			}
		}
	}
}

func TestSkalInDataNotCode(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Verify 'skall' (Norwegian) form exists in data
	forms := kb.GetFormsByText("skall")
	if len(forms) == 0 {
		t.Error("'skall' form not found in data - it should be data, not code")
	}

	// Verify 'shell' form exists in data
	forms = kb.GetFormsByText("shell")
	if len(forms) == 0 {
		t.Error("'shell' form not found in data")
	}

	// Verify 'skull' form exists (skall -> skull via norwegian lens)
	forms = kb.GetFormsByConcept("skull")
	if len(forms) == 0 {
		t.Error("'skull' concept forms not found - relates to skall")
	}
}

func TestConceptRelationEndpoints(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Verify all relation endpoints exist
	conceptIDs := make(map[string]bool)
	for _, c := range kb.Concepts {
		conceptIDs[c.ID] = true
	}

	for _, r := range kb.Relations {
		if !conceptIDs[r.From] {
			t.Errorf("relation %s->%s: 'from' concept %q does not exist", r.From, r.To, r.From)
		}
		if !conceptIDs[r.To] {
			t.Errorf("relation %s->%s: 'to' concept %q does not exist", r.From, r.To, r.To)
		}
	}
}

func TestAllKnowledgeAccessible(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// Verify AllConcepts returns all concepts
	all := kb.AllConcepts()
	if len(all) != len(kb.Concepts) {
		t.Errorf("AllConcepts returned %d, expected %d", len(all), len(kb.Concepts))
	}

	// Verify AllForms returns all forms
	allForms := kb.AllForms()
	if len(allForms) != len(kb.Forms) {
		t.Errorf("AllForms returned %d, expected %d", len(allForms), len(kb.Forms))
	}

	// Verify AllRelations returns all relations
	allRelations := kb.AllRelations()
	if len(allRelations) != len(kb.Relations) {
		t.Errorf("AllRelations returned %d, expected %d", len(allRelations), len(kb.Relations))
	}
}