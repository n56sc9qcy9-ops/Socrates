package decipher

import (
	"testing"
)

// PROPAGATION TESTS
// =============================================================================

func TestPropagateActivation_Basic(t *testing.T) {
	g := NewActivationGraph()

	// A -> B -> C chain
	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.0, ConfidencePlausible)
	g.AddNode("c", 0.0, ConfidencePlausible)

	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "related", 0.6, ConfidencePlausible)

	// Propagate
	g.PropagateActivation()

	// Node 'a' should have strength 1.0 (direct)
	aNode := g.GetNode("a")
	if aNode.Strength != 1.0 {
		t.Errorf("node 'a' should have strength 1.0, got %f", aNode.Strength)
	}

	// Node 'b' should get propagated strength
	bNode := g.GetNode("b")
	if bNode.Strength <= 0.0 {
		t.Error("node 'b' should have propagated strength")
	}

	// Node 'c' should get propagated strength (with double decay)
	cNode := g.GetNode("c")
	if cNode.Depth == 0 || cNode.Depth > 2 {
		t.Errorf("node 'c' depth should be 1 or 2, got %d", cNode.Depth)
	}
}

func TestPropagateActivation_Decay(t *testing.T) {
	g := NewActivationGraph()
	g.MaxDepth = 3
	g.DecayFactor = 0.5

	// A -> B chain with edge weight 0.8
	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.0, ConfidencePlausible)
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)

	g.PropagateActivation()

	// Expected: 1.0 * 0.8 * 0.5 = 0.4
	bNode := g.GetNode("b")
	if bNode.Strength != 0.4 {
		t.Errorf("node 'b' should have strength 0.4 (1.0 * 0.8 * 0.5), got %f", bNode.Strength)
	}
}

func TestPropagateActivation_BoundedDepth(t *testing.T) {
	g := NewActivationGraph()
	g.MaxDepth = 1 // Very shallow
	g.DecayFactor = 0.5

	// A -> B -> C chain
	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.0, ConfidencePlausible)
	g.AddNode("c", 0.0, ConfidencePlausible)

	g.AddEdge("a", "b", "r1", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "r2", 0.8, ConfidencePlausible)

	g.PropagateActivation()

	// Node 'c' should NOT get propagation (depth would be 2, but max is 1)
	cNode := g.GetNode("c")
	if cNode.Depth != 0 || cNode.Strength != 0.0 {
		t.Errorf("node 'c' should not be activated (exceeds max depth 1), got depth=%d, strength=%f", cNode.Depth, cNode.Strength)
	}

	// Node 'b' should be activated
	bNode := g.GetNode("b")
	if bNode.Depth == 0 {
		t.Error("node 'b' should have depth 1")
	}
}

func TestPropagateActivation_CyclePrevention(t *testing.T) {
	g := NewActivationGraph()
	g.MaxDepth = 2      // Limited depth
	g.DecayFactor = 0.5 // Decay to prevent inflation

	// A -> B -> A cycle
	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.0, ConfidencePlausible)

	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "a", "related", 0.8, ConfidencePlausible)

	g.PropagateActivation()

	// Node 'a' should only have its original strength plus limited re-propagation
	// The Visited flag prevents infinite loops, but allows one re-entry per path
	aNode := g.GetNode("a")

	// Expected behavior: 'a' starts at 1.0, then:
	// - Propagate a->b: b gets 1.0 * 0.8 * 0.5 = 0.4
	// - Propagate b->a: a gets additional 0.4 * 0.8 * 0.5 = 0.16 (1 cycle allowed)
	// - Propagate a(b from step 2): would check Visited, prevent re-entry
	// Final: ~1.16 (1 cycle allowed, not infinite)
	if aNode.Strength > 1.2 {
		t.Errorf("node 'a' should not have inflated strength from cycles, got %f (expected ~1.16 with 1 allowed cycle)", aNode.Strength)
	}

	if aNode.Strength < 1.0 {
		t.Errorf("node 'a' should retain at least original strength, got %f", aNode.Strength)
	}
}

// =============================================================================
