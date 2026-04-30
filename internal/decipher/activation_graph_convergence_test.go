package decipher

import (
	"testing"

	"socrates/internal/knowledge"
)

// GRAPH-TO-CONVERGENCE BRIDGE TESTS
// =============================================================================

func TestToConvergenceResult(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	// Create test data
	channels := []ChannelResult{
		{
			Name: "test",
			Signals: []Signal{
				{Text: "test", Target: "test", Channel: "test", Confidence: ConfidenceVerified, Weight: 0.8},
			},
		},
	}

	passageSignals := []PassageSignal{
		{Token: "test", Concept: "test", Weight: 0.7, Confidence: ConfidencePlausible, MatchForm: "test"},
	}

	fuzzyMatches := []MatchEvidence{}
	conceptExpansions := map[string][]knowledge.DecipherConceptRelation{}

	graph := BuildGraphFromEvidence(channels, passageSignals, fuzzyMatches, conceptExpansions, kb)
	graph.PropagateActivation()

	result := graph.ToConvergenceResult()

	if len(result.ActivatedConcepts) == 0 {
		t.Error("convergence result should have activated concepts")
	}

	if result.CoActivationScore < 0.0 || result.CoActivationScore > 1.0 {
		t.Errorf("co-activation score should be in [0, 1], got %f", result.CoActivationScore)
	}
}

func TestGraphStats(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("a", 0.8, ConfidenceVerified)
	g.AddNode("b", 0.7, ConfidenceVerified)
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)

	stats := g.GraphStats()

	if stats["node_count"] != 2 {
		t.Errorf("node_count should be 2, got %v", stats["node_count"])
	}

	if stats["edge_count"] != 1 {
		t.Errorf("edge_count should be 1, got %v", stats["edge_count"])
	}
}

// =============================================================================
