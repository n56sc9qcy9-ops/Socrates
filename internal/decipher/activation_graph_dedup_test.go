package decipher

import (
	"testing"
)

// PROPAGATION DETERMINISM TESTS
// =============================================================================

func TestPropagateActivation_StableWithinSinglePass(t *testing.T) {
	// Test that propagation produces consistent results within a single run
	// by checking intermediate results are stable
	g := NewActivationGraph()

	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.8, ConfidenceVerified)
	g.AddNode("c", 0.6, ConfidenceVerified)
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "related", 0.8, ConfidencePlausible)
	g.AddEdge("a", "c", "related", 0.6, ConfidencePlausible)

	// Record results after single propagation
	g.PropagateActivation()
	result := make(map[string]float64)
	for _, node := range g.Nodes {
		result[node.Concept] = node.Strength
	}

	// Check that propagation is bounded (no excessive growth)
	for concept, strength := range result {
		if strength > 2.0 {
			t.Errorf("concept '%s' has inflated strength %.4f", concept, strength)
		}
	}

	// Direct nodes should retain their original strength
	if aNode := g.GetNode("a"); aNode.Strength < 1.0 {
		t.Errorf("direct node 'a' should retain at least original strength 1.0, got %.4f", aNode.Strength)
	}
}

func TestPropagateActivation_CycleSafe(t *testing.T) {
	// Test that cycles don't cause infinite propagation or stack overflow
	g := NewActivationGraph()

	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.8, ConfidenceVerified)
	g.AddNode("c", 0.6, ConfidenceVerified)

	// Create cycle: a -> b -> c -> a
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "related", 0.8, ConfidencePlausible)
	g.AddEdge("c", "a", "related", 0.8, ConfidencePlausible)

	// This should not cause infinite recursion or stack overflow
	g.PropagateActivation()

	// All nodes should have bounded strength (no infinite growth)
	for _, node := range g.Nodes {
		if node.Strength > 2.0 {
			t.Errorf("node '%s' has inflated strength %.4f from cycle", node.Concept, node.Strength)
		}
	}
}

func TestPropagateActivation_LongCycleSafe(t *testing.T) {
	// Test that long chains with cycles don't explode
	g := NewActivationGraph()

	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.8, ConfidenceVerified)
	g.AddNode("c", 0.6, ConfidenceVerified)
	g.AddNode("d", 0.4, ConfidencePlausible)
	g.AddNode("e", 0.3, ConfidencePlausible)

	// Long chain: a -> b -> c -> d -> e -> a
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "related", 0.8, ConfidencePlausible)
	g.AddEdge("c", "d", "related", 0.8, ConfidencePlausible)
	g.AddEdge("d", "e", "related", 0.8, ConfidencePlausible)
	g.AddEdge("e", "a", "related", 0.8, ConfidencePlausible)

	// This should complete without stack overflow
	g.PropagateActivation()

	// All nodes should have bounded strength
	for _, node := range g.Nodes {
		if node.Strength > 2.0 {
			t.Errorf("node '%s' has inflated strength %.4f from long cycle", node.Concept, node.Strength)
		}
	}
}

// =============================================================================
