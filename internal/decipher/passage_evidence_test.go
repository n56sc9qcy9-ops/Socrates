package decipher

import (
	"testing"
)

// EVIDENCE PATH TESTS
// =============================================================================

func TestPassageField_EvidencePathsPopulated(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("truth", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Truth field should have evidence paths
	truthField := fields.GetField("truth")
	if truthField == nil {
		t.Error("should have truth field")
	} else if len(truthField.EvidencePaths) == 0 {
		t.Error("truth field should have evidence paths from activation graph")
	}
}

func TestEvidencePath_Deduplication(t *testing.T) {
	// Create paths with same ID
	path1 := EvidencePath{
		SourceToken: "love",
		SourceType:  "passage_signal",
		SourceForm:  "love",
		MatchForm:   "love",
		Confidence:  ConfidenceVerified,
		Weight:      1.0,
	}

	path2 := EvidencePath{
		SourceToken: "love",
		SourceType:  "passage_signal",
		SourceForm:  "love",
		MatchForm:   "love",
		Confidence:  ConfidenceVerified,
		Weight:      1.0,
	}

	// Both should have the same ID
	id1 := path1.EvidenceID()
	id2 := path2.EvidenceID()

	if id1 != id2 {
		t.Errorf("same evidence should have same ID: %s vs %s", id1, id2)
	}

	// Merging should deduplicate
	merged := mergeEvidencePaths([]EvidencePath{path1}, []EvidencePath{path2})
	if len(merged) != 1 {
		t.Errorf("merged evidence paths should have 1 entry, got %d", len(merged))
	}
}

// =============================================================================
