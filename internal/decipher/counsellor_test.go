package decipher

import (
	"os"
	"strings"
	"testing"

	"socrates/internal/knowledge"
)

func TestBuildCounsellorField(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Create passage fields with a concept that has transmutation relations
	passageFields := PassageFields{
		{Concept: "resentment", Strength: 0.8, IsDirectEvidence: true, TokenSources: []string{"resentment"}},
		{Concept: "love", Strength: 0.7, TokenSources: []string{"love"}},
	}

	field := BuildCounsellorField(passageFields, kb)

	// Should find resentment source
	if field == nil {
		t.Fatal("expected counsellor field, got nil")
	}

	if len(field.SourceFields) == 0 {
		t.Fatal("expected at least one source field")
	}

	if field.SourceFields[0].Concept != "resentment" {
		t.Errorf("expected source 'resentment', got '%s'", field.SourceFields[0].Concept)
	}

	// Should have suggestions from transmute.yaml
	if len(field.Suggestions) == 0 {
		t.Errorf("expected at least one suggestion for resentment")
	}

	// First suggestion should be resentment -> forgiveness
	if len(field.Suggestions) > 0 && field.Suggestions[0].SourceConcept != "resentment" {
		t.Errorf("expected source 'resentment', got '%s'", field.Suggestions[0].SourceConcept)
	}

	// Should have evidence paths
	if len(field.EvidencePaths) == 0 {
		t.Errorf("expected at least one evidence path")
	}
}

func TestBuildCounsellorFieldNoMatch(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Create passage fields with no concepts that have transmutation relations
	passageFields := PassageFields{
		{Concept: "one", Strength: 0.9, TokenSources: []string{"one"}},
		{Concept: "being", Strength: 0.8, TokenSources: []string{"being"}},
	}

	field := BuildCounsellorField(passageFields, kb)

	// Should be nil when no transmutation matches
	if field != nil {
		t.Errorf("expected nil for non-transmutation concepts, got %v", field)
	}
}

func TestBuildCounsellorFieldWeakStrength(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Create passage fields with strength below threshold
	passageFields := PassageFields{
		{Concept: "resentment", Strength: 0.1, TokenSources: []string{"resentment"}},
	}

	field := BuildCounsellorField(passageFields, kb)

	// Should be nil when strength is below threshold
	if field != nil {
		t.Errorf("expected nil for weak strength, got %v", field)
	}
}

func TestCounsellorFieldKinds(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// resentment: transmutes_to
	passageFields := PassageFields{
		{Concept: "resentment", Strength: 0.8, IsDirectEvidence: true, TokenSources: []string{"resentment"}},
	}
	field := BuildCounsellorField(passageFields, kb)
	if field == nil || len(field.Suggestions) == 0 {
		t.Fatal("expected counsellor field for resentment")
	}
	if field.Suggestions[0].Kind != "transmutes_to" {
		t.Errorf("expected kind 'transmutes_to', got '%s'", field.Suggestions[0].Kind)
	}

	// guilt: softens_through
	passageFields = PassageFields{
		{Concept: "guilt", Strength: 0.8, TokenSources: []string{"guilt"}},
	}
	field = BuildCounsellorField(passageFields, kb)
	if field == nil || len(field.Suggestions) == 0 {
		t.Fatal("expected counsellor field for guilt")
	}
	if field.Suggestions[0].Kind != "softens_through" {
		t.Errorf("expected kind 'softens_through', got '%s'", field.Suggestions[0].Kind)
	}

	// pride: corrects_through
	passageFields = PassageFields{
		{Concept: "pride", Strength: 0.8, TokenSources: []string{"pride"}},
	}
	field = BuildCounsellorField(passageFields, kb)
	if field == nil || len(field.Suggestions) == 0 {
		t.Fatal("expected counsellor field for pride")
	}
	if field.Suggestions[0].Kind != "corrects_through" {
		t.Errorf("expected kind 'corrects_through', got '%s'", field.Suggestions[0].Kind)
	}

	// fear: releases_into
	passageFields = PassageFields{
		{Concept: "fear", Strength: 0.8, TokenSources: []string{"fear"}},
	}
	field = BuildCounsellorField(passageFields, kb)
	if field == nil || len(field.Suggestions) == 0 {
		t.Fatal("expected counsellor field for fear")
	}
	if field.Suggestions[0].Kind != "releases_into" {
		t.Errorf("expected kind 'releases_into', got '%s'", field.Suggestions[0].Kind)
	}

	// shame: grounds_in
	passageFields = PassageFields{
		{Concept: "shame", Strength: 0.8, TokenSources: []string{"shame"}},
	}
	field = BuildCounsellorField(passageFields, kb)
	if field == nil || len(field.Suggestions) == 0 {
		t.Fatal("expected counsellor field for shame")
	}
	if field.Suggestions[0].Kind != "grounds_in" {
		t.Errorf("expected kind 'grounds_in', got '%s'", field.Suggestions[0].Kind)
	}
}

func TestCounsellorFieldMultipleSuggestions(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Check that each transmutation source maps to its correction target
	tests := []struct {
		source      string
		target      string
		minStrength float64
	}{
		{"resentment", "forgiveness", 0.5},
		{"pride", "humility", 0.5},
		{"fear", "trust", 0.5},
		{"guilt", "mercy", 0.4},
		{"shame", "self_acceptance", 0.4},
		{"attachment", "letting_go", 0.4},
		{"control", "surrender", 0.4},
		{"avoidance", "presence", 0.3},
		{"judgment", "mercy", 0.4},
		{"isolation", "connection", 0.5},
		{"grasping", "openness", 0.3},
		{"confusion", "clarity", 0.3},
		{"rigidity", "flexibility", 0.3},
		{"conflict", "harmony", 0.3},
		{"stagnation", "growth", 0.3},
	}

	for _, tc := range tests {
		passageFields := PassageFields{
			{Concept: tc.source, Strength: tc.minStrength + 0.1, TokenSources: []string{tc.source}},
		}
		field := BuildCounsellorField(passageFields, kb)
		if field == nil {
			t.Errorf("expected counsellor field for %s, got nil", tc.source)
			continue
		}
		if len(field.Suggestions) == 0 {
			t.Errorf("expected at least one suggestion for %s", tc.source)
			continue
		}
		found := false
		for _, s := range field.Suggestions {
			if s.TargetConcept == tc.target {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected suggestion %s -> %s, suggestions: %v", tc.source, tc.target, field.Suggestions)
		}
	}
}

func TestTransmutationConfidence(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	passageFields := PassageFields{
		{Concept: "resentment", Strength: 0.8, IsDirectEvidence: true, TokenSources: []string{"resentment"}},
	}
	field := BuildCounsellorField(passageFields, kb)
	if field == nil || len(field.Suggestions) == 0 {
		t.Fatal("expected counsellor field")
	}

	// resentment -> forgiveness should be "plausible"
	if field.Suggestions[0].Confidence != "plausible" {
		t.Errorf("expected confidence 'plausible', got '%s'", field.Suggestions[0].Confidence)
	}
}

func TestCounsellorDirectEvidenceProducesSuggestion(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Create passage field with direct concept (depth 0)
	passageFields := PassageFields{
		{Concept: "resentment", Strength: 0.8, Depth: 0, IsDirectEvidence: true, TokenSources: []string{"resentment"}},
	}

	field := BuildCounsellorField(passageFields, kb)

	if field == nil {
		t.Fatal("expected counsellor field for direct evidence")
	}

	if len(field.Suggestions) == 0 {
		t.Fatal("expected suggestion for direct resentment")
	}

	// Direct evidence should have normal confidence (plausible for resentment)
	if field.Suggestions[0].Confidence != "plausible" {
		t.Errorf("expected plausible confidence for direct evidence, got %s", field.Suggestions[0].Confidence)
	}
}

func TestCounsellorPropagatedNeighborEvidenceDowngraded(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Create passage field with propagated concept (depth 1) - simulating graph expansion
	// In real usage, depth 1 means it came from a neighbor relation, not direct evidence
	passageFields := PassageFields{
		{Concept: "fear", Strength: 0.4, Depth: 1, IsDirectEvidence: false, TokenSources: []string{"resentment -> neighbor: fear"}},
	}

	field := BuildCounsellorField(passageFields, kb)

	if field == nil {
		t.Fatal("expected counsellor field for propagated evidence")
	}


	if len(field.Suggestions) == 0 {
		t.Fatal("expected suggestion for propagated fear")
	}

	// Propagated evidence (depth 1) should have downgraded confidence
	// fear has confidence "plausible" in transmute.yaml, so depth 1 should downgrade to "speculative"
	if field.Suggestions[0].Confidence != "speculative" {
		t.Errorf("expected speculative confidence for depth 1 propagated evidence, got %s", field.Suggestions[0].Confidence)
	}
}

func TestCounsellorDefaultOutputNotCommanding(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Create reading with counsellor field
	passageFields := PassageFields{
		{Concept: "resentment", Strength: 0.8, Depth: 0, IsDirectEvidence: true, TokenSources: []string{"resentment"}},
	}
	field := BuildCounsellorField(passageFields, kb)

	reading := Reading{
		Input:           "resentment",
		CounsellorField: field,
	}

	output := RenderReading(reading)

	// Default output should not contain commanding language
	commandingPhrases := []string{
		"if field remains",
		"you should",
		"must",
		"always",
		"never",
	}
	for _, phrase := range commandingPhrases {
		if strings.Contains(output, phrase) {
			t.Errorf("default output contains commanding phrase: %s", phrase)
		}
	}

	// Should contain humble phrasing
	if !strings.Contains(output, "may be") && !strings.Contains(output, "possible") {
		t.Errorf("default output should contain humble phrasing like 'may be' or 'possible'")
	}
}

func TestCounsellorDefaultOutputCapped(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Create many passage fields that would produce many suggestions
	passageFields := PassageFields{
		{Concept: "resentment", Strength: 0.8, Depth: 0, IsDirectEvidence: true, TokenSources: []string{"resentment"}},
		{Concept: "guilt", Strength: 0.8, Depth: 0, TokenSources: []string{"guilt"}},
		{Concept: "judgment", Strength: 0.8, Depth: 0, TokenSources: []string{"judgment"}},
		{Concept: "fear", Strength: 0.8, Depth: 0, TokenSources: []string{"fear"}},
		{Concept: "shame", Strength: 0.8, Depth: 0, TokenSources: []string{"shame"}},
		{Concept: "pride", Strength: 0.8, Depth: 0, TokenSources: []string{"pride"}},
		{Concept: "attachment", Strength: 0.8, Depth: 0, TokenSources: []string{"attachment"}},
		{Concept: "control", Strength: 0.8, Depth: 0, TokenSources: []string{"control"}},
	}

	field := BuildCounsellorField(passageFields, kb)

	reading := Reading{
		Input:           "many tensions",
		CounsellorField: field,
	}

	output := RenderReading(reading)

	// Count suggestions shown - each line has exactly one suggestion
	// Output format: "  - possible field 'X' may be softened through 'Y'"
	// We count the "may be softened" pattern which appears once per suggestion
	suggestionCount := strings.Count(output, "may be softened")

	// Default output should be capped at 3 suggestions
	if suggestionCount > 3 {
		t.Errorf("default output should be capped at 3 suggestions, found %d", suggestionCount)
	}

	// Should indicate more exist in debug mode
	if !strings.Contains(output, "+") && field != nil && len(field.Suggestions) > 3 {
		t.Errorf("default output should indicate when suggestions are truncated")
	}
}

func TestCounsellorDebugOutputShowsReasoning(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Create reading with mixed direct and propagated concepts
	passageFields := PassageFields{
		{Concept: "resentment", Strength: 0.8, Depth: 0, TokenSources: []string{"direct"}},
		{Concept: "guilt", Strength: 0.6, Depth: 1, TokenSources: []string{"propagated"}},
	}
	field := BuildCounsellorField(passageFields, kb)

	reading := Reading{
		Input:           "mixed evidence",
		CounsellorField: field,
	}

	output := RenderReadingWithOptions(reading, DebugRenderOptions())

	// Debug output should show all suggestions
	if !strings.Contains(output, "Source fields:") {
		t.Errorf("debug output should show source fields")
	}

	// Debug output should show evidence paths
	if !strings.Contains(output, "Evidence paths:") {
		t.Errorf("debug output should show evidence paths")
	}

	// Debug output should show strength values
	if !strings.Contains(output, "strength:") {
		t.Errorf("debug output should show strength values")
	}
}

func TestNoHardcodedConceptBehavior(t *testing.T) {
	kb := knowledge.LoadOrPanic()

	// Get all transmutation source concepts from data
	dataSources := make(map[string]bool)
	for _, tr := range kb.AllTransmutations() {
		dataSources[tr.From] = true
	}

	// Verify specific concepts exist in data (not hardcoded in Go)
	specialConcepts := []string{"resentment", "forgiveness", "humility", "fear", "trust", "love"}
	for _, concept := range specialConcepts {
		// If the concept is a transmutation target, it should be in the data sources
		found := false
		for _, tr := range kb.AllTransmutations() {
			if tr.From == concept || tr.To == concept {
				found = true
				break
			}
		}
		if !found && concept != "love" {
			// love might not have a transmutation - that's ok
			t.Logf("concept %s not found as transmutation target", concept)
		}
	}

	// Verify counsellor.go doesn't contain hardcoded concept checks
	counsellorCode, err := os.ReadFile("counsellor.go")
	if err != nil {
		t.Fatal("could not read counsellor.go")
	}

	hardcodedConcepts := []string{"resentment", "forgiveness", "humility", "fear", "trust", "love"}
	for _, concept := range hardcodedConcepts {
		// These concepts appearing in string literals in counsellor.go would indicate hardcoding
		// We check they don't appear in case-sensitive contexts that would indicate special handling
		pattern := "\"" + concept + "\""
		if strings.Contains(string(counsellorCode), pattern) {
			t.Errorf("counsellor.go contains hardcoded concept reference: %s", concept)
		}
	}
}

// TestEngineCounsellorNeighborVsDirectEvidence proves that:
// 1. Form/glyph/script evidence (IsDirect=true) renders with normal confidence
// 2. Symbolic neighbor expansion (IsDirect=false) is downgraded
// 3. Direct suggestions appear first in ranked output (before propagated)
// This is an end-to-end regression test that runs through the real engine.
func TestEngineCounsellorNeighborVsDirectEvidence(t *testing.T) {
	kb := knowledge.LoadOrPanic()
	engine := NewEngineWithKnowledge(kb)

	// Analyze "resentment" - this triggers:
	// - Direct: primitive match for resentment (IsDirect=true, weight 0.7)
	// - Indirect: neighbor expansion to fear/guilt/judgment (IsDirect=false, weight 0.4)
	reading := engine.Analyze("resentment")
if reading.Input == "" {
		t.Fatal("expected reading")
	}

	counsellor := reading.CounsellorField
	if counsellor == nil {
		t.Fatal("expected counsellor field")
	}

	// Check source evidence classification
	directCount := 0
	propagatedCount := 0
	for _, src := range counsellor.SourceFields {
		if src.IsDirectEvidence {
			directCount++
			t.Logf("DIRECT: %s (strength=%.2f, depth=%d)", src.Concept, src.Strength, src.Depth)
		} else {
			propagatedCount++
			t.Logf("PROPAGATED: %s (strength=%.2f, depth=%d)", src.Concept, src.Strength, src.Depth)
		}
	}

	// resentment itself should be direct (from primitive match)
	// fear/guilt/judgment are from neighbor expansion, should be propagated
	if directCount == 0 {
		t.Error("expected at least one direct evidence source (resentment from primitive match)")
	}
	if propagatedCount == 0 {
		t.Error("expected at least one propagated/neighbor source (fear/guilt/judgment from neighbor)")
	}

	// Verify suggestions are ranked: direct evidence first, then by confidence/strength
	// resentment -> forgiveness should be first (direct evidence, plausible)
	if len(counsellor.Suggestions) == 0 {
		t.Fatal("expected suggestions")
	}
	firstSuggestion := counsellor.Suggestions[0]
	if firstSuggestion.SourceConcept != "resentment" {
		t.Errorf("expected first suggestion to be from resentment (direct), got %s",
			firstSuggestion.SourceConcept)
	}
	if firstSuggestion.Confidence != "plausible" {
		t.Errorf("expected first suggestion to have plausible confidence, got %s",
			firstSuggestion.Confidence)
	}
	t.Logf("First suggestion: %s -> %s [confidence=%s] ✓", firstSuggestion.SourceConcept,
		firstSuggestion.TargetConcept, firstSuggestion.Confidence)

	// All remaining suggestions should be from propagated sources (speculative)
	for i := 1; i < len(counsellor.Suggestions); i++ {
		sugg := counsellor.Suggestions[i]
		t.Logf("Ranked suggestion %d: %s -> %s [confidence=%s]", i, sugg.SourceConcept,
			sugg.TargetConcept, sugg.Confidence)
	}

	// Verify suggestions have appropriate confidence based on evidence type
	for _, sugg := range counsellor.Suggestions {
		t.Logf("Suggestion: %s -> %s [confidence=%s]", sugg.SourceConcept, sugg.TargetConcept, sugg.Confidence)
		
		// Find the source for this suggestion
		var isDirectSrc bool
		for _, src := range counsellor.SourceFields {
			if src.Concept == sugg.SourceConcept {
				isDirectSrc = src.IsDirectEvidence
				break
			}
		}
		
		if isDirectSrc {
			// Direct evidence: normal confidence (plausible for resentment, verified for some)
			if sugg.Confidence == "speculative" {
				t.Errorf("direct evidence source %s should not have speculative confidence, got %s",
					sugg.SourceConcept, sugg.Confidence)
			}
		} else {
			// Propagated/neighbor evidence: should be downgraded
			if sugg.Confidence != "speculative" {
				t.Errorf("propagated/neighbor source %s should have speculative confidence, got %s",
					sugg.SourceConcept, sugg.Confidence)
			}
		}
	}
}
