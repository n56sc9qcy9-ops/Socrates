package decipher

import (
	"testing"
)

// ANALYZE PASSAGE FROM TOKENS TESTS
// =============================================================================

func TestAnalyzePassageFromTokens_Direct(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	tokens := []string{"love", "truth"}
	fields := AnalyzePassageFromTokens(tokens, engine.Knowledge)

	if len(fields) == 0 {
		t.Error("should produce passage fields from tokens")
	}
}

func TestAnalyzePassageFromTokens_TokenProvenance(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	tokens := []string{"light", "truth"}
	fields := AnalyzePassageFromTokens(tokens, engine.Knowledge)

	// Both fields should exist with evidence paths
	lightField := fields.GetField("light")
	if lightField == nil {
		t.Error("should have light field")
	} else {
		// Light field should have evidence paths (the key requirement)
		if len(lightField.EvidencePaths) == 0 {
			t.Error("light field should have evidence paths")
		}
		if lightField.Strength <= 0 {
			t.Error("light field should have positive strength")
		}
	}

	// Truth field should also exist
	truthField := fields.GetField("truth")
	if truthField == nil {
		t.Error("should have truth field")
	} else {
		if len(truthField.EvidencePaths) == 0 {
			t.Error("truth field should have evidence paths")
		}
	}
}
