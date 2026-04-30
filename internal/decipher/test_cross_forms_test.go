package decipher

import (
	"socrates/internal/knowledge"
	"testing"
)

func TestCrossLanguageFormsViaYAML(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatal(err)
	}

	// These forms were in the removed channel_cross_language.go
	// They should now be found in forms.yaml
	forms := []string{"chi", "ruach", "prana", "om", "logos", "dao"}

	allForms := kb.AllForms()
	formMap := make(map[string]bool)
	for _, f := range allForms {
		formMap[f.Form] = true
	}

	for _, form := range forms {
		if !formMap[form] {
			t.Errorf("form %q not found in knowledge forms", form)
		}
	}
}

// TestNoCrossLanguageChannel verifies cross-language channel is removed
func TestNoCrossLanguageChannel(t *testing.T) {
	forms := GenerateForms("chi")

	// Run all channels - should not fail even without cross-language channel
	kb, _ := knowledge.LoadFromEmbed()
	results := RunAllChannels(forms, kb)

	// Verify no channel named "Cross-Language" exists
	for _, ch := range results {
		if ch.Name == "Cross-Language" {
			t.Error("Cross-Language channel should be removed")
		}
	}

	// Verify other channels still work
	if len(results) < 4 {
		t.Errorf("expected at least 4 channels, got %d", len(results))
	}
}
