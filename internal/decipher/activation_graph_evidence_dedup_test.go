package decipher

import (
	"testing"
)

// EVIDENCE DEDUPLICATION TESTS
// =============================================================================

func TestActivationNode_AddEvidence_DeduplicatesByID(t *testing.T) {
	node := &ActivationNode{
		Concept:  "test",
		Evidence: make([]EvidencePath, 0),
	}

	// Add first evidence
	evidence1 := EvidencePath{
		SourceToken: "token1",
		SourceForm:  "form1",
		SourceType:  "direct",
		Confidence:  ConfidenceVerified,
		Weight:      0.8,
	}
	if !node.AddEvidence(evidence1) {
		t.Error("first evidence should be added")
	}

	// Add duplicate evidence (same ID)
	evidence1Dup := EvidencePath{
		SourceToken: "token1", // Same token
		SourceForm:  "form1",  // Same form
		SourceType:  "direct", // Same type
		Confidence:  ConfidenceVerified,
		Weight:      0.8,
	}
	if node.AddEvidence(evidence1Dup) {
		t.Error("duplicate evidence should not be added")
	}

	if len(node.Evidence) != 1 {
		t.Errorf("node should have 1 evidence, got %d", len(node.Evidence))
	}
}

func TestActivationNode_AddEvidence_DifferentIDNotDuplicates(t *testing.T) {
	node := &ActivationNode{
		Concept:  "test",
		Evidence: make([]EvidencePath, 0),
	}

	// Add evidence from different sources
	evidence1 := EvidencePath{
		SourceToken: "token1",
		SourceForm:  "form1",
		SourceType:  "direct",
		Confidence:  ConfidenceVerified,
		Weight:      0.8,
	}
	evidence2 := EvidencePath{
		SourceToken: "token2", // Different token
		SourceForm:  "form2",
		SourceType:  "passage", // Different type
		Confidence:  ConfidencePlausible,
		Weight:      0.6,
	}

	node.AddEvidence(evidence1)
	node.AddEvidence(evidence2)

	if len(node.Evidence) != 2 {
		t.Errorf("node should have 2 evidence from different sources, got %d", len(node.Evidence))
	}
}
