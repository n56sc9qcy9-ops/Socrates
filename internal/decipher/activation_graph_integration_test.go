package decipher

import (
	"testing"
)

// INTEGRATION TESTS
// =============================================================================

func TestEngine_Analyze_RoutesThroughGraph(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("inspired")

	// Reading should be produced normally
	if reading.Input != "inspired" {
		t.Errorf("expected input 'inspired', got '%s'", reading.Input)
	}

	// Should have convergence result (populated through graph)
	if len(reading.Convergence.ActivatedConcepts) == 0 {
		t.Error("convergence should have activated concepts from graph")
	}

	// Should have co-activation score
	if reading.Convergence.CoActivationScore < 0.0 || reading.Convergence.CoActivationScore > 1.0 {
		t.Errorf("co-activation score should be in [0, 1], got %f", reading.Convergence.CoActivationScore)
	}

	// Should have a reading
	if reading.ConciseReading == "" {
		t.Error("should produce a concise reading")
	}
}

func TestActivationGraph_EdgeCase_EmptyGraph(t *testing.T) {
	g := NewActivationGraph()

	g.PropagateActivation()

	stats := g.GraphStats()
	if stats["node_count"] != 0 {
		t.Errorf("empty graph should have 0 nodes, got %v", stats["node_count"])
	}

	score := g.GetCoActivationScore()
	if score != 0.0 {
		t.Errorf("empty graph should have co-activation score 0, got %f", score)
	}

	topNodes := g.GetTopNodes(5)
	if len(topNodes) != 0 {
		t.Errorf("empty graph should have 0 top nodes, got %d", len(topNodes))
	}
}

func TestActivationGraph_EdgeCase_SingleNode(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("solo", 0.8, ConfidenceVerified)
	g.PropagateActivation()

	score := g.GetCoActivationScore()
	if score != 0.0 {
		t.Errorf("single node should have co-activation score 0, got %f", score)
	}

	topNodes := g.GetTopNodes(3)
	if len(topNodes) != 1 {
		t.Errorf("single node graph should have 1 top node, got %d", len(topNodes))
	}
}

// =============================================================================
