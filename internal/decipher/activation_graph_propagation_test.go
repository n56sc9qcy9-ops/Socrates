package decipher

import (
	"testing"

	"socrates/internal/knowledge"
)

// GRAPH BUILDING TESTS
// =============================================================================

func TestBuildGraphFromEvidence_Basic(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	// Create minimal test data
	channels := []ChannelResult{
		{
			Name: "glyph",
			Signals: []Signal{
				{Text: "spirit", Target: "spirit", Channel: "glyph", Confidence: ConfidenceVerified, Weight: 0.8},
			},
		},
	}

	passageSignals := []PassageSignal{
		{Token: "in", Concept: "spirit", Weight: 0.7, Confidence: ConfidencePlausible, MatchForm: "in"},
	}

	fuzzyMatches := []MatchEvidence{
		{InputForm: "inspired", AnchorForm: "spirit", Method: "phonetic", Distance: 0.2, Weight: 0.8},
	}

	conceptExpansions := map[string][]knowledge.DecipherConceptRelation{
		"spirit": {
			{To: "breath", Type: "related_to", Weight: 0.6, Confidence: ConfidencePlausible},
		},
	}

	graph := BuildGraphFromEvidence(channels, passageSignals, fuzzyMatches, conceptExpansions, kb)

	if graph == nil {
		t.Fatal("BuildGraphFromEvidence should return non-nil graph")
	}

	// Should have spirit node from direct channel
	spiritNode := graph.GetNode("spirit")
	if spiritNode == nil {
		t.Fatal("spirit node should exist")
	}

	// Should have breath node from graph expansion
	breathNode := graph.GetNode("breath")
	if breathNode == nil {
		t.Fatal("breath node should exist from graph expansion")
	}

	// Should have edge from spirit to breath
	edges := graph.GetEdgesFrom("spirit")
	if len(edges) == 0 {
		t.Error("should have edge from spirit")
	}
}

func TestBuildGraphFromEvidence_Deduplication(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	// Create data with duplicate evidence from SAME channel (same text, same channel)
	// These should be deduplicated
	channels := []ChannelResult{
		{
			Name: "channel1",
			Signals: []Signal{
				{Text: "test", Target: "test", Channel: "channel1", Confidence: ConfidenceVerified, Weight: 0.5},
				{Text: "test", Target: "test", Channel: "channel1", Confidence: ConfidencePlausible, Weight: 0.6}, // duplicate - same channel + same text
			},
		},
	}

	passageSignals := []PassageSignal{}
	fuzzyMatches := []MatchEvidence{}
	conceptExpansions := map[string][]knowledge.DecipherConceptRelation{}

	graph := BuildGraphFromEvidence(channels, passageSignals, fuzzyMatches, conceptExpansions, kb)

	testNode := graph.GetNode("test")
	if testNode == nil {
		t.Fatal("test node should exist")
	}

	// Same channel + same text = same key = 1 evidence entry (deduplicated)
	if len(testNode.Evidence) != 1 {
		t.Errorf("should have 1 evidence entry from same channel duplicates, got %d", len(testNode.Evidence))
	}
}

// =============================================================================
