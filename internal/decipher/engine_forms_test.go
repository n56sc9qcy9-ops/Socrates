package decipher

import (
	"testing"
)


func TestGenerateForms_Inspired(t *testing.T) {
	forms := GenerateForms("inspired")

	if forms.Normalized != "inspired" {
		t.Errorf("expected 'inspired', got '%s'", forms.Normalized)
	}

	// Should generate fragment paths
	if len(forms.Fragments) == 0 {
		t.Error("inspired should generate fragment paths")
	}

	// Check that some paths include spirit-related fragments
	foundSpiritPath := false
	for _, path := range forms.Fragments {
		for _, part := range path.Parts {
			if part == "in" || part == "spirit" || part == "spire" {
				foundSpiritPath = true
				break
			}
		}
	}
	if !foundSpiritPath {
		t.Error("inspired should have at least one path with 'in', 'spirit', or 'spire'")
	}
}

// =============================================================================
