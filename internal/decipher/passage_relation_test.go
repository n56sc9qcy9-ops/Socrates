package decipher

import (
	"strings"
	"testing"
)

// RELATION PATH TESTS
// =============================================================================

func TestPassageField_RelationPathsExist(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("love light", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Some fields should have relation paths when concepts are connected
	foundRelationPaths := false
	for _, field := range fields {
		if len(field.RelationPaths) > 0 {
			foundRelationPaths = true
			// Verify path format: "from ->(type) to"
			for _, path := range field.RelationPaths {
				if !strings.Contains(path, "->") {
					t.Errorf("relation path should contain '->': %s", path)
				}
			}
		}
	}

	if !foundRelationPaths {
		t.Error("at least some fields should have relation paths when concepts are connected")
	}
}

func TestPassageField_RelationPathsFromGraph(t *testing.T) {
	// Create a graph with explicit edges
	graph := NewActivationGraph()
	graph.AddNode("light", 1.0, ConfidenceVerified)
	graph.AddNode("sun", 0.8, ConfidencePlausible)
	graph.AddNode("truth", 0.6, ConfidenceSpeculative)
	graph.AddEdge("light", "sun", "synonym", 0.9, ConfidenceVerified)
	graph.AddEdge("sun", "truth", "related", 0.7, ConfidencePlausible)
	graph.PropagateActivation()

	fields := BuildPassageFieldsFromGraph(graph)

	// Light should have relation path to sun
	lightField := fields.GetField("light")
	if lightField == nil {
		t.Fatal("should have light field")
	}

	if len(lightField.RelationPaths) == 0 {
		t.Error("light field should have relation paths from graph edges")
	}
}

// =============================================================================
