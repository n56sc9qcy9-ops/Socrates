package decipher

import (
	"strings"
	"testing"
)

// DEFAULT VS DEBUG RENDER TESTS
// =============================================================================

func TestPassageField_DefaultOutputConcise(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("love truth")
	output := RenderReadingWithOptions(reading, DefaultRenderOptions())

	// Default output should be concise (under 2000 chars)
	if len(output) > 2000 {
		t.Errorf("default output seems too long (%d chars), should be concise", len(output))
	}

	// Default output should show top fields summary
	if !strings.Contains(output, "Top fields") && !strings.Contains(output, "field") {
		t.Error("default output should show passage field summary")
	}
}

func TestPassageField_DefaultOutputNoCandidates(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("love truth")
	output := RenderReadingWithOptions(reading, DefaultRenderOptions())

	// Default output should NOT include all candidates
	if strings.Contains(output, "Candidate Neighbors") {
		t.Error("default output should not include candidate neighbor list")
	}
}

func TestPassageField_DebugOutputFullDetails(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("love truth light")
	output := RenderReadingWithOptions(reading, DebugRenderOptions())

	// Debug output should include passage fields
	if !strings.Contains(output, "Passage Fields") {
		t.Error("debug output should include passage fields section")
	}

	// Debug should be longer than default
	defaultOutput := RenderReadingWithOptions(reading, DefaultRenderOptions())
	if len(output) <= len(defaultOutput) {
		t.Error("debug output should be longer than default output")
	}
}

func TestPassageField_DebugOutputEvidencePaths(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("truth")
	output := RenderReadingWithOptions(reading, DebugRenderOptions())

	// Debug output should show evidence
	if !strings.Contains(output, "evidence") {
		t.Error("debug output should include evidence paths")
	}
}

// =============================================================================
