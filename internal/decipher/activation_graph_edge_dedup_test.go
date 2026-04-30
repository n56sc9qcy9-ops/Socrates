package decipher

import (
	"testing"
)

// EDGE DEDUPLICATION TESTS
// =============================================================================

func TestAddEdge_DeduplicatesByIdentity(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.8, ConfidenceVerified)

	// Add edge first time
	e1 := g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)

	// Try to add duplicate edge (same from, to, relationType)
	e2 := g.AddEdge("a", "b", "related", 0.9, ConfidencePlausible) // Different weight

	// Should return existing edge, not create new one
	if e1 != e2 {
		t.Error("duplicate edge should return existing edge")
	}

	// Graph should only have one edge
	if len(g.Edges) != 1 {
		t.Errorf("graph should have 1 edge, got %d", len(g.Edges))
	}

	// Edge should retain original weight
	if e1.Weight != 0.8 {
		t.Errorf("edge should retain original weight 0.8, got %.2f", e1.Weight)
	}
}

func TestAddEdge_DifferentRelationTypesNotDuplicates(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.8, ConfidenceVerified)

	// Add edge with "related"
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)

	// Add edge with "similar" - different relation type, should be kept
	g.AddEdge("a", "b", "similar", 0.7, ConfidencePlausible)

	// Should have 2 edges
	if len(g.Edges) != 2 {
		t.Errorf("graph should have 2 edges (different relation types), got %d", len(g.Edges))
	}
}

func TestAddEdge_DifferentToNotDuplicate(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.8, ConfidenceVerified)
	g.AddNode("c", 0.6, ConfidenceVerified)

	// Add edge a -> b
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)

	// Add edge a -> c - different target, should be kept
	g.AddEdge("a", "c", "related", 0.7, ConfidencePlausible)

	// Should have 2 edges
	if len(g.Edges) != 2 {
		t.Errorf("graph should have 2 edges (different targets), got %d", len(g.Edges))
	}
}

// =============================================================================
