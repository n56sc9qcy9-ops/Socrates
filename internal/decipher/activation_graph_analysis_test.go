package decipher

import (
	"testing"
)

// GRAPH ANALYSIS TESTS
// =============================================================================

func TestGetTopNodes(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("low", 0.2, ConfidencePlausible)
	g.AddNode("high", 0.9, ConfidenceVerified)
	g.AddNode("medium", 0.5, ConfidencePlausible)

	top2 := g.GetTopNodes(2)

	if len(top2) != 2 {
		t.Errorf("GetTopNodes(2) should return 2 nodes, got %d", len(top2))
	}

	if top2[0].Concept != "high" {
		t.Errorf("top node should be 'high', got '%s'", top2[0].Concept)
	}

	if top2[1].Concept != "medium" {
		t.Errorf("second node should be 'medium', got '%s'", top2[1].Concept)
	}
}

func TestGetCoActivationScore(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("a", 0.8, ConfidenceVerified)
	g.AddNode("b", 0.7, ConfidenceVerified)
	g.AddNode("c", 0.6, ConfidenceVerified)
	g.AddNode("d", 0.5, ConfidencePlausible)

	// Connect a-b and b-c, but not d
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "related", 0.7, ConfidencePlausible)

	score := g.GetCoActivationScore()

	// Among top 10 concepts, connected pairs / total pairs
	// a-b are connected, b-c are connected, a-c not directly connected
	// Score should be > 0 since some pairs are connected
	if score <= 0.0 {
		t.Error("co-activation score should be > 0 when concepts are connected")
	}
}

func TestGetRelationPaths(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("a", 0.8, ConfidenceVerified)
	g.AddNode("b", 0.7, ConfidenceVerified)
	g.AddNode("c", 0.6, ConfidenceVerified)

	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "similar", 0.7, ConfidencePlausible)

	paths := g.GetRelationPaths([]string{"a", "b", "c"})

	// Should include a -> b and b -> c paths
	if len(paths) == 0 {
		t.Error("should return relation paths between concepts")
	}
}

// =============================================================================
