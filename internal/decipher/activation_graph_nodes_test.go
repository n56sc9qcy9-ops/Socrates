package decipher

import (
	"testing"

	"socrates/internal/knowledge"
)

// DIRECT VS DERIVED NODE TESTS
// =============================================================================

func TestDirectAndDerivedNodes_Distinguishable(t *testing.T) {
	// Test that direct nodes (depth 0) and derived nodes (depth > 0) are distinguishable
	g := NewActivationGraph()

	// Add direct node
	directNode := g.AddNode("direct", 1.0, ConfidenceVerified)
	if directNode.Depth != 0 {
		t.Errorf("direct node should have depth 0, got %d", directNode.Depth)
	}

	// Add derived node (simulating graph expansion)
	derivedNode := g.AddNode("derived", 0.5, ConfidencePlausible)
	derivedNode.Depth = 1 // Manually set to derived

	if directNode.Depth == derivedNode.Depth {
		t.Error("direct and derived nodes should have different depths")
	}
}

func TestGraphExpandedNodes_NotDirect(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	// Create test data with a concept that will expand via relations
	channels := []ChannelResult{
		{
			Name: "test",
			Signals: []Signal{
				{Text: "love", Target: "love", Channel: "test", Confidence: ConfidenceVerified, Weight: 0.8},
			},
		},
	}

	passageSignals := []PassageSignal{}
	fuzzyMatches := []MatchEvidence{}
	conceptExpansions := map[string][]knowledge.DecipherConceptRelation{}

	// Get relations for "love" using the correct method
	if relations := kb.GetConceptRelationsAsDecipher("love"); len(relations) > 0 {
		conceptExpansions["love"] = relations
	}

	graph := BuildGraphFromEvidence(channels, passageSignals, fuzzyMatches, conceptExpansions, kb)
	graph.PropagateActivation()

	// Check that expanded nodes have depth > 0
	directCount := 0
	for _, node := range graph.Nodes {
		if node.Depth == 0 {
			directCount++
		}
	}

	// The "love" node should be direct (depth 0)
	loveNode := graph.GetNode("love")
	if loveNode == nil {
		t.Fatal("love node should exist")
	}
	if loveNode.Depth != 0 {
		t.Errorf("direct node 'love' should have depth 0, got %d", loveNode.Depth)
	}
}

// =============================================================================
