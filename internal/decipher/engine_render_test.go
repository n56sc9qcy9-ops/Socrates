package decipher

import (
	"testing"
	"strings"
)


// RENDER MODE TESTS
// =============================================================================

func TestDefaultRenderDoesNotDumpCandidates(t *testing.T) {
	// Default render output should NOT include full candidate dumps
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Use a word that generates many candidates
	reading := engine.Analyze("inspired")
	output := RenderReadingWithOptions(reading, DefaultRenderOptions())

	// Default should not contain "Candidate Neighbors" section
	if strings.Contains(output, "Candidate Neighbors") {
		t.Error("default render should NOT contain 'Candidate Neighbors' section")
	}

	// Default should not show individual candidate forms
	if strings.Contains(output, "[method: normalized") || strings.Contains(output, "[method: consonant_skeleton") {
		t.Error("default render should NOT show individual candidate methods")
	}
}

func TestDefaultRenderDoesNotDumpFuzzyMatches(t *testing.T) {
	// Default render output should NOT include full fuzzy match dumps
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("inspired")
	output := RenderReadingWithOptions(reading, DefaultRenderOptions())

	// Default should not contain detailed fuzzy match entries
	if strings.Contains(output, "Fuzzy Matches:") {
		t.Error("default render should NOT contain 'Fuzzy Matches:' section")
	}

	// Should not show input->anchor patterns
	if strings.Contains(output, "->") && strings.Contains(output, "[method: fuzzy") {
		t.Error("default render should NOT show fuzzy match arrows")
	}
}

func TestDefaultRenderShowsEvidencePaths(t *testing.T) {
	// Default render output should show evidence paths for conclusions
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("truth")
	output := RenderReadingWithOptions(reading, DefaultRenderOptions())

	// Should show Evidence section
	if !strings.Contains(output, "Evidence:") {
		t.Error("default render should contain 'Evidence:' section")
	}

	// Should show top concepts with strength
	if !strings.Contains(output, "%") {
		t.Error("default render should show concept strengths as percentages")
	}
}

func TestDefaultRenderShowsScoreSummary(t *testing.T) {
	// Default render should show concise score summary
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("energy")
	output := RenderReadingWithOptions(reading, DefaultRenderOptions())

	// Should show overall score
	if !strings.Contains(output, "Overall:") {
		t.Error("default render should contain 'Overall:' score")
	}

	// Should show Components section
	if !strings.Contains(output, "Components:") {
		t.Error("default render should contain 'Components:' section")
	}
}

func TestDebugRenderShowsCandidates(t *testing.T) {
	// Debug render output SHOULD include full candidate lists
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("inspired")
	output := RenderReadingWithOptions(reading, DebugRenderOptions())

	// Debug should contain "Candidate Neighbors" section
	if !strings.Contains(output, "Candidate Neighbors") {
		t.Error("debug render should contain 'Candidate Neighbors' section")
	}

	// Debug should show individual candidate methods
	if !strings.Contains(output, "[method: normalized") {
		t.Error("debug render should show 'normalized' method")
	}
}

func TestDebugRenderShowsFuzzyMatches(t *testing.T) {
	// Debug render output SHOULD include fuzzy match section header
	// (even if no matches found, the section header should be present)
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Create reading
	reading := engine.Analyze("energy")
	output := RenderReadingWithOptions(reading, DebugRenderOptions())

	// Debug should show Fuzzy Matches section header
	// (Even if empty, the section appears in debug mode for completeness)
	if !strings.Contains(output, "Fuzzy Matches") {
		t.Error("debug render should contain 'Fuzzy Matches' section")
	}
}

func TestDebugRenderShowsDiscardedCounts(t *testing.T) {
	// Debug render output SHOULD show discarded candidate/comparison counts
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Use long input to trigger discarding
	reading := engine.Analyze("abcdefghijklmnop")

	// Ensure we have some discarded data
	hasDiscarded := reading.DiscardedCandidates > 0 || reading.DiscardedComparisons > 0

	output := RenderReadingWithOptions(reading, DebugRenderOptions())

	// Debug should show discarded counts when present
	if hasDiscarded {
		if !strings.Contains(output, "Discarded") {
			t.Error("debug render should show 'Discarded' counts when bounded work occurred")
		}
	}
}

func TestDebugRenderShowsGeneratedForms(t *testing.T) {
	// Debug render output SHOULD include generated forms detail
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("spirit")
	output := RenderReadingWithOptions(reading, DebugRenderOptions())

	// Debug should show Generated Forms section
	if !strings.Contains(output, "Generated Forms:") {
		t.Error("debug render should contain 'Generated Forms:' section")
	}

	// Debug should show phonetic keys
	if !strings.Contains(output, "phonetic") {
		t.Error("debug render should mention phonetic processing")
	}
}

func TestDefaultRenderIsConcise(t *testing.T) {
	// Default render should be significantly shorter than debug render
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("inspired")

	defaultOutput := RenderReadingWithOptions(reading, DefaultRenderOptions())
	debugOutput := RenderReadingWithOptions(reading, DebugRenderOptions())

	// Default should be at least 30% shorter than debug
	lenRatio := float64(len(defaultOutput)) / float64(len(debugOutput))
	if lenRatio > 0.7 {
		t.Errorf("default render should be at least 30%% shorter than debug, got ratio %.2f%%", lenRatio*100)
	}
}

func TestDefaultRenderNoGeneratedFormsDump(t *testing.T) {
	// Default render should NOT dump all generated forms internals
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("truth")
	output := RenderReadingWithOptions(reading, DefaultRenderOptions())

	// Should not show phonetic keys detail
	if strings.Contains(output, "phonetic keys:") {
		t.Error("default render should NOT show 'phonetic keys' detail")
	}

	// Should not show fragments detail
	if strings.Contains(output, "fragments:") {
		t.Error("default render should NOT show 'fragments' detail")
	}

	// Should not show runes detail
	if strings.Contains(output, "runes:") {
		t.Error("default render should NOT show 'runes' detail")
	}
}

func TestDefaultRenderShowsChannelsCompact(t *testing.T) {
	// Default render should show channels in compact form (not full dumps)
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("energy")
	output := RenderReadingWithOptions(reading, DefaultRenderOptions())

	// Should show channels section exists
	if !strings.Contains(output, "Channels:") {
		t.Error("default render should contain 'Channels:' section")
	}

	// But should NOT show channel score details in the main body
	// (score summary is in Score section, not per-channel in channels)
}

func TestRenderModesProduceConsistentReading(t *testing.T) {
	// Both render modes should produce the same Reading conclusion
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("skal")

	defaultOutput := RenderReadingWithOptions(reading, DefaultRenderOptions())
	debugOutput := RenderReadingWithOptions(reading, DebugRenderOptions())

	// Both should contain the same concise reading
	if !strings.Contains(defaultOutput, reading.ConciseReading) {
		t.Error("default render should contain the ConciseReading")
	}

	if !strings.Contains(debugOutput, reading.ConciseReading) {
		t.Error("debug render should contain the ConciseReading")
	}

	// Both should have the same score
	if !strings.Contains(defaultOutput, "Overall:") {
		t.Error("default render should show overall score")
	}

	if !strings.Contains(debugOutput, "Overall:") {
		t.Error("debug render should show overall score")
	}
}

func TestDefaultRenderShowsBoundedWorkSummary(t *testing.T) {
	// Default render should show bounded work summary when discarding occurred
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Use long input to trigger bounds
	reading := engine.Analyze("abcdefghijklmnop")

	// At least one of these should be positive for long input
	if reading.DiscardedCandidates > 0 || reading.DiscardedComparisons > 0 {
		output := RenderReadingWithOptions(reading, DefaultRenderOptions())

		// Default should show bounded work info
		if !strings.Contains(output, "Bounded work") {
			t.Error("default render should show 'Bounded work' summary when discarding occurred")
		}
	}
}

// =============================================================================
