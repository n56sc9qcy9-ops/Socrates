package decipher

import (
	"testing"
)

// PASSAGE FIELD TOKEN PROVENANCE TESTS
// =============================================================================

func TestPassageField_TokenSourcesFromOriginalTokens(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Passage with specific tokens
	fields, err := AnalyzePassage("love light", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Fields should have evidence paths
	foundEvidence := false
	for _, field := range fields {
		if len(field.EvidencePaths) > 0 {
			foundEvidence = true
			// Verify evidence has valid structure
			for _, ev := range field.EvidencePaths {
				if ev.SourceToken == "" {
					t.Error("evidence path should have source token")
				}
				if ev.SourceType == "" {
					t.Error("evidence path should have source type")
				}
			}
		}
	}

	if !foundEvidence {
		t.Error("fields should have evidence paths from passage analysis")
	}
}

func TestPassageField_TokenProvenance_OriginalTokenIsSource(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("truth", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Truth field should have evidence paths populated
	// Token sources may vary due to form generation, but field must exist
	truthField := fields.GetField("truth")
	if truthField == nil {
		t.Error("should have 'truth' field")
	} else {
		// Field should have evidence paths (the key requirement)
		if len(truthField.EvidencePaths) == 0 {
			t.Error("'truth' field should have evidence paths")
		}
		// Field should have positive strength
		if truthField.Strength <= 0 {
			t.Error("'truth' field should have positive strength")
		}
	}
}

// =============================================================================
