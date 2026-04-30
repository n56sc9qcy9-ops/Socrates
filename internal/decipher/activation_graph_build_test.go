package decipher

import (
	"testing"
)

// ACTIVATION GRAPH TESTS
// =============================================================================

func TestNewActivationGraph(t *testing.T) {
	g := NewActivationGraph()

	if g == nil {
		t.Fatal("NewActivationGraph should return a non-nil graph")
	}

	if g.NodeCount() != 0 {
		t.Errorf("new graph should have 0 nodes, got %d", g.NodeCount())
	}

	if g.EdgeCount() != 0 {
		t.Errorf("new graph should have 0 edges, got %d", g.EdgeCount())
	}

	if g.MaxDepth != 2 {
		t.Errorf("default MaxDepth should be 2, got %d", g.MaxDepth)
	}

	if g.DecayFactor != 0.5 {
		t.Errorf("default DecayFactor should be 0.5, got %f", g.DecayFactor)
	}
}

func TestActivationGraph_AddNode(t *testing.T) {
	g := NewActivationGraph()

	// Add first node
	node1 := g.AddNode("spirit", 0.8, ConfidenceVerified)
	if node1 == nil {
		t.Fatal("AddNode should return non-nil node")
	}

	if node1.Concept != "spirit" {
		t.Errorf("node concept should be 'spirit', got '%s'", node1.Concept)
	}

	if node1.Strength != 0.8 {
		t.Errorf("node strength should be 0.8, got %f", node1.Strength)
	}

	if node1.BaseStrength != 0.8 {
		t.Errorf("node base strength should be 0.8, got %f", node1.BaseStrength)
	}

	if node1.Confidence != ConfidenceVerified {
		t.Errorf("node confidence should be '%s', got '%s'", ConfidenceVerified, node1.Confidence)
	}

	if node1.Depth != 0 {
		t.Errorf("new node depth should be 0, got %d", node1.Depth)
	}

	// Add second node
	node2 := g.AddNode("breath", 0.5, ConfidencePlausible)
	if node2.Concept != "breath" {
		t.Errorf("node concept should be 'breath', got '%s'", node2.Concept)
	}

	if g.NodeCount() != 2 {
		t.Errorf("graph should have 2 nodes, got %d", g.NodeCount())
	}

	// Adding same concept should merge, not duplicate
	node1Again := g.AddNode("spirit", 0.9, ConfidencePlausible)
	if node1Again != node1 {
		t.Error("AddNode for existing concept should return existing node")
	}

	if g.NodeCount() != 2 {
		t.Errorf("adding existing concept should not increase count, got %d", g.NodeCount())
	}

	// Merged strength should be strongest base
	if node1.BaseStrength != 0.9 {
		t.Errorf("merged base strength should be 0.9, got %f", node1.BaseStrength)
	}
}

func TestActivationGraph_AddEdge(t *testing.T) {
	g := NewActivationGraph()

	// Add nodes first
	g.AddNode("spirit", 0.8, ConfidenceVerified)
	g.AddNode("breath", 0.5, ConfidencePlausible)

	// Add edge between them
	edge := g.AddEdge("spirit", "breath", "related_to", 0.7, ConfidencePlausible)

	if edge == nil {
		t.Fatal("AddEdge should return non-nil edge")
	}

	if edge.From != "spirit" {
		t.Errorf("edge From should be 'spirit', got '%s'", edge.From)
	}

	if edge.To != "breath" {
		t.Errorf("edge To should be 'breath', got '%s'", edge.To)
	}

	if edge.RelationType != "related_to" {
		t.Errorf("edge relation type should be 'related_to', got '%s'", edge.RelationType)
	}

	if edge.Weight != 0.7 {
		t.Errorf("edge weight should be 0.7, got %f", edge.Weight)
	}

	if g.EdgeCount() != 1 {
		t.Errorf("graph should have 1 edge, got %d", g.EdgeCount())
	}
}

func TestActivationGraph_GetEdges(t *testing.T) {
	g := NewActivationGraph()

	// Add nodes
	g.AddNode("a", 0.5, ConfidencePlausible)
	g.AddNode("b", 0.5, ConfidencePlausible)
	g.AddNode("c", 0.5, ConfidencePlausible)

	// Add edges
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("a", "c", "similar", 0.6, ConfidencePlausible)
	g.AddEdge("b", "c", "related", 0.7, ConfidencePlausible)

	// Get edges from 'a'
	edgesFromA := g.GetEdgesFrom("a")
	if len(edgesFromA) != 2 {
		t.Errorf("should have 2 edges from 'a', got %d", len(edgesFromA))
	}

	// Get edges to 'c'
	edgesToC := g.GetEdgesTo("c")
	if len(edgesToC) != 2 {
		t.Errorf("should have 2 edges to 'c', got %d", len(edgesToC))
	}

	// Get edges from non-existent node
	edgesFromX := g.GetEdgesFrom("x")
	if len(edgesFromX) != 0 {
		t.Errorf("should have 0 edges from 'x', got %d", len(edgesFromX))
	}
}

// =============================================================================
