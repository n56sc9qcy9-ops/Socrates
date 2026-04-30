package decipher

import (
	"testing"
)

// NODE EVIDENCE TESTS
// =============================================================================

func TestActivationNode_AddEvidence(t *testing.T) {
	g := NewActivationGraph()
	node := g.AddNode("spirit", 0.8, ConfidenceVerified)

	// Add first evidence
	evidence1 := EvidencePath{
		SourceToken: "in",
		SourceForm:  "in",
		SourceType:  "direct_channel",
		Confidence:  ConfidenceVerified,
		Weight:      0.5,
	}

	if !node.AddEvidence(evidence1) {
		t.Error("first evidence should be added")
	}

	if len(node.Evidence) != 1 {
		t.Errorf("node should have 1 evidence, got %d", len(node.Evidence))
	}

	// Add duplicate evidence
	evidence1Duplicate := EvidencePath{
		SourceToken: "in",
		SourceForm:  "in",
		SourceType:  "direct_channel",
		Confidence:  ConfidencePlausible,
		Weight:      0.6,
	}

	if node.AddEvidence(evidence1Duplicate) {
		t.Error("duplicate evidence should not be added")
	}

	if len(node.Evidence) != 1 {
		t.Errorf("duplicate evidence should not increase count, got %d", len(node.Evidence))
	}

	// Add different evidence
	evidence2 := EvidencePath{
		SourceToken: "spirit",
		SourceForm:  "spirit",
		SourceType:  "fuzzy_match",
		MatchForm:   "spirit",
		Confidence:  ConfidencePlausible,
		Weight:      0.7,
	}

	if !node.AddEvidence(evidence2) {
		t.Error("different evidence should be added")
	}

	if len(node.Evidence) != 2 {
		t.Errorf("node should have 2 evidence, got %d", len(node.Evidence))
	}
}

func TestEvidencePath_EvidenceID(t *testing.T) {
	// Test evidence ID generation for deduplication
	fuzzyEv := EvidencePath{
		SourceToken: "inspired",
		SourceForm:  "inspired",
		SourceType:  "fuzzy_match",
		MatchForm:   "spirit",
	}

	fuzzyID := fuzzyEv.EvidenceID()
	if fuzzyID != "fuzzy|inspired|spirit" {
		t.Errorf("fuzzy evidence ID should be 'fuzzy|inspired|spirit', got '%s'", fuzzyID)
	}

	passageEv := EvidencePath{
		SourceToken: "in",
		SourceForm:  "in",
		SourceType:  "passage_signal",
		MatchForm:   "in",
	}

	passageID := passageEv.EvidenceID()
	if passageID != "passage|in|in|in" {
		t.Errorf("passage evidence ID should be 'passage|in|in|in', got '%s'", passageID)
	}

	graphEv := EvidencePath{
		SourceType:   "graph_expansion",
		RelationType: "related_to",
		RelationFrom: "spirit",
		RelationTo:   "breath",
	}

	graphID := graphEv.EvidenceID()
	if graphID != "graph|related_to|spirit|breath" {
		t.Errorf("graph evidence ID should be 'graph|related_to|spirit|breath', got '%s'", graphID)
	}
}

// =============================================================================
