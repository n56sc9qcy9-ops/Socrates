package decipher

import (
	"testing"
)

// BUILD PASSAGE FIELDS FROM GRAPH TESTS
// =============================================================================

func TestBuildPassageFieldsFromGraph_Empty(t *testing.T) {
	graph := NewActivationGraph()
	fields := BuildPassageFieldsFromGraph(graph)

	if len(fields) != 0 {
		t.Errorf("empty graph should produce empty fields, got %d", len(fields))
	}
}

func TestBuildPassageFieldsFromGraph_WithNodes(t *testing.T) {
	graph := NewActivationGraph()
	graph.AddNode("love", 1.0, ConfidenceVerified)
	graph.AddNode("light", 0.8, ConfidencePlausible)
	graph.AddEdge("love", "light", "related", 0.7, ConfidencePlausible)
	graph.PropagateActivation()

	fields := BuildPassageFieldsFromGraph(graph)

	if len(fields) != 2 {
		t.Errorf("graph with 2 nodes should produce 2 fields, got %d", len(fields))
	}

	// Both fields should have valid structure
	for _, field := range fields {
		if field.Concept == "" {
			t.Error("field should have a concept")
		}
		if field.Strength <= 0 {
			t.Error("field should have positive strength")
		}
	}
}

func TestBuildPassageFieldsFromGraph_TokenSourcesFromEvidence(t *testing.T) {
	graph := NewActivationGraph()
	node := graph.AddNode("truth", 1.0, ConfidenceVerified)
	node.Evidence = append(node.Evidence, EvidencePath{
		SourceToken: "truth",
		SourceType:  "passage_signal",
		Confidence:  ConfidenceVerified,
		Weight:      1.0,
	})

	fields := BuildPassageFieldsFromGraph(graph)

	truthField := fields.GetField("truth")
	if truthField == nil {
		t.Fatal("should have truth field")
	}

	if len(truthField.TokenSources) == 0 {
		t.Error("truth field should have 'truth' as token source from evidence")
	}

	foundTruth := false
	for _, src := range truthField.TokenSources {
		if src == "truth" {
			foundTruth = true
			break
		}
	}
	if !foundTruth {
		t.Error("truth field should have 'truth' as token source")
	}
}

// =============================================================================
