package decipher

import (
	"testing"
)

// GRAPH PROPAGATION BEHAVIOR TESTS
// =============================================================================

func TestGraphPropagate_MultipleSources(t *testing.T) {
	g := NewActivationGraph()
	g.AddNode("x", 1.0, ConfidenceVerified)
	g.AddNode("y", 0.8, ConfidenceVerified)
	g.AddNode("z", 0.5, ConfidencePlausible)
	g.AddEdge("x", "z", "synonym", 0.9, ConfidenceVerified)
	g.AddEdge("y", "z", "related", 0.7, ConfidencePlausible)
	g.PropagateActivation()

	// Z should receive activation from both sources
	zStrength := g.GetNode("z").Strength
	if zStrength <= 0.5 {
		t.Errorf("Z should have received propagated activation: %.4f", zStrength)
	}
}

// =============================================================================
