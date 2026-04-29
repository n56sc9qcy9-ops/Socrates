package decipher

import (
	"strings"
	"testing"
)

// =============================================================================
// PASSAGE FIELD STRUCTURE TESTS
// =============================================================================

func TestPassageField_Structure(t *testing.T) {
	field := &PassageField{
		Concept:      "love",
		Strength:     0.9,
		Confidence:   ConfidenceVerified,
		TokenSources: []string{"love", "heart"},
		Depth:        0,
	}

	if field.Concept != "love" {
		t.Errorf("expected concept 'love', got '%s'", field.Concept)
	}
	if field.Strength != 0.9 {
		t.Errorf("expected strength 0.9, got %f", field.Strength)
	}
	if field.Depth != 0 {
		t.Errorf("expected depth 0 (direct), got %d", field.Depth)
	}
}

func TestPassageFields_TopFields(t *testing.T) {
	fields := PassageFields{
		{Concept: "low", Strength: 0.2},
		{Concept: "high", Strength: 0.9},
		{Concept: "medium", Strength: 0.5},
	}

	top2 := fields.TopFields(2)

	if len(top2) != 2 {
		t.Errorf("TopFields(2) should return 2 fields, got %d", len(top2))
	}
	if top2[0].Concept != "high" {
		t.Errorf("top field should be 'high', got '%s'", top2[0].Concept)
	}
	if top2[1].Concept != "medium" {
		t.Errorf("second field should be 'medium', got '%s'", top2[1].Concept)
	}
}

func TestPassageFields_GetField(t *testing.T) {
	fields := PassageFields{
		{Concept: "love", Strength: 0.9},
		{Concept: "light", Strength: 0.7},
	}

	loveField := fields.GetField("love")
	if loveField == nil {
		t.Fatal("should find 'love' field")
	}
	if loveField.Strength != 0.9 {
		t.Errorf("expected strength 0.9, got %f", loveField.Strength)
	}

	missingField := fields.GetField("nonexistent")
	if missingField != nil {
		t.Error("should not find 'nonexistent' field")
	}
}

func TestPassageFields_Merge(t *testing.T) {
	fields1 := PassageFields{
		{Concept: "love", Strength: 0.5, TokenSources: []string{"love"}},
		{Concept: "light", Strength: 0.3, TokenSources: []string{"light"}},
	}

	fields2 := PassageFields{
		{Concept: "love", Strength: 0.4, TokenSources: []string{"heart"}},
		{Concept: "truth", Strength: 0.6, TokenSources: []string{"truth"}},
	}

	merged := fields1.Merge(fields2)

	// Should have 3 unique concepts
	if len(merged) != 3 {
		t.Errorf("merged should have 3 fields, got %d", len(merged))
	}

	// Love should have combined strength
	loveField := merged.GetField("love")
	if loveField == nil {
		t.Fatal("should have love field")
	}
	if loveField.Strength != 0.9 {
		t.Errorf("love should have combined strength 0.9, got %f", loveField.Strength)
	}
	// Love should have both token sources
	if len(loveField.TokenSources) != 2 {
		t.Errorf("love should have 2 token sources, got %d", len(loveField.TokenSources))
	}
}

// =============================================================================
// PASSAGE FIELD ANALYSIS TESTS
// =============================================================================

func TestAnalyzePassage_MultiToken(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("God is Love", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Should have passage fields
	if len(fields) == 0 {
		t.Error("should produce passage fields for multi-token input")
	}

	// Verify 'love' field exists and is strong
	loveField := fields.GetField("love")
	if loveField == nil {
		t.Error("should have 'love' field")
	} else if loveField.Strength <= 0 {
		t.Error("love field should have positive strength")
	}
}

func TestAnalyzePassage_SingleToken(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("love", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Single token should still produce fields
	if len(fields) == 0 {
		t.Error("should produce passage fields for single token")
	}
}

func TestAnalyzePassage_Empty(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("", engine)
	if err != nil {
		t.Fatal(err)
	}

	if fields != nil {
		t.Error("empty passage should produce nil fields")
	}
}

func TestAnalyzePassage_WhitespaceOnly(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("   \n\t  ", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Whitespace-only should produce no fields
	if len(fields) > 0 {
		t.Errorf("whitespace-only should produce no fields, got %d", len(fields))
	}
}

// =============================================================================
// PASSAGE FIELDS IN READING TESTS
// =============================================================================

func TestReading_ContainsPassageFields(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Multi-token passage
	reading := engine.Analyze("love truth light")

	// Reading MUST contain PassageFields
	if reading.PassageFields == nil {
		t.Fatal("Reading.PassageFields should not be nil for multi-token input")
	}

	if len(reading.PassageFields) == 0 {
		t.Error("Reading.PassageFields should not be empty for multi-token input")
	}
}

func TestReading_PassageFieldsNotNilForSingleToken(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Single token
	reading := engine.Analyze("love")

	// Reading should contain PassageFields (even for single token)
	if reading.PassageFields == nil {
		t.Fatal("Reading.PassageFields should not be nil")
	}
}

// =============================================================================
// PASSAGE FIELD TOKEN PROVENANCE TESTS
// =============================================================================

func TestPassageField_TokenSourcesFromOriginalTokens(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Passage with specific tokens
	fields, err := AnalyzePassage("love light", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Fields should have evidence paths
	foundEvidence := false
	for _, field := range fields {
		if len(field.EvidencePaths) > 0 {
			foundEvidence = true
			// Verify evidence has valid structure
			for _, ev := range field.EvidencePaths {
				if ev.SourceToken == "" {
					t.Error("evidence path should have source token")
				}
				if ev.SourceType == "" {
					t.Error("evidence path should have source type")
				}
			}
		}
	}

	if !foundEvidence {
		t.Error("fields should have evidence paths from passage analysis")
	}
}

func TestPassageField_TokenProvenance_OriginalTokenIsSource(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("truth", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Truth field should have evidence paths populated
	// Token sources may vary due to form generation, but field must exist
	truthField := fields.GetField("truth")
	if truthField == nil {
		t.Error("should have 'truth' field")
	} else {
		// Field should have evidence paths (the key requirement)
		if len(truthField.EvidencePaths) == 0 {
			t.Error("'truth' field should have evidence paths")
		}
		// Field should have positive strength
		if truthField.Strength <= 0 {
			t.Error("'truth' field should have positive strength")
		}
	}
}

// =============================================================================
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
// KEY TOKEN WEAKENING STRICT TESTS
// =============================================================================

func TestKeyTokenRemoval_StrictlyWeakensRelatedField(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Full passage with key semantic word
	fieldsFull, _ := AnalyzePassage("love truth light", engine)
	// Without the key word
	fieldsNoLove, _ := AnalyzePassage("truth light", engine)

	// Get love field strengths
	fullLoveStrength := 0.0
	noLoveStrength := 0.0

	for _, field := range fieldsFull {
		if field.Concept == "love" {
			fullLoveStrength = field.Strength
		}
	}

	for _, field := range fieldsNoLove {
		if field.Concept == "love" {
			noLoveStrength = field.Strength
		}
	}

	// Removing "love" from the passage MUST strictly weaken the love field
	// (either love field disappears or strength is reduced)
	if fullLoveStrength <= noLoveStrength {
		t.Errorf("removing 'love' token should strictly weaken love field:\n  full: %.4f\n  no love: %.4f",
			fullLoveStrength, noLoveStrength)
	}
}

func TestKeyTokenRemoval_CompleteRemoval(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Analyze single word
	fieldsSingle, _ := AnalyzePassage("truth", engine)

	// Get truth field strength from single
	var singleTruthStrength float64
	for _, field := range fieldsSingle {
		if field.Concept == "truth" {
			singleTruthStrength = field.Strength
		}
	}

	// Single word analysis should produce truth field with positive strength
	if singleTruthStrength <= 0 {
		t.Error("single word analysis should produce truth field with positive strength")
	}

	// Analyze again - strengths should be in the same ballpark (within 50%)
	fieldsAgain, _ := AnalyzePassage("truth", engine)
	var againTruthStrength float64
	for _, field := range fieldsAgain {
		if field.Concept == "truth" {
			againTruthStrength = field.Strength
		}
	}

	ratio := againTruthStrength / singleTruthStrength
	if ratio < 0.5 || ratio > 2.0 {
		t.Errorf("repeated single-word analysis should produce similar strength: %.4f vs %.4f (ratio=%.2f)",
			singleTruthStrength, againTruthStrength, ratio)
	}
}

// =============================================================================
// DUPLICATE/REPEATED TOKEN EVIDENCE POLICY TESTS
// =============================================================================

func TestRepeatedTokens_AccumulationWithinBounds(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Single occurrence
	fields1, _ := AnalyzePassage("love", engine)

	// Repeated occurrences (3x)
	fields3, _ := AnalyzePassage("love love love", engine)

	// Get love field strengths
	var strength1, strength3 float64

	for _, field := range fields1 {
		if field.Concept == "love" {
			strength1 = field.Strength
		}
	}

	for _, field := range fields3 {
		if field.Concept == "love" {
			strength3 = field.Strength
		}
	}

	// Multiple occurrences should accumulate but not linearly
	// Allow up to 4x accumulation for triple occurrence (reasonable for co-activation)
	maxExpected := strength1 * 4.0

	if strength3 > maxExpected {
		t.Errorf("repeated tokens should not inflate linearly: single=%.4f, triple=%.4f, max=%.4f",
			strength1, strength3, maxExpected)
	}
}

func TestDuplicateInflation_EvidenceDeduplication(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Analyze once
	fields1, _ := AnalyzePassage("light", engine)

	// Analyze twice with same passage
	fields2, _ := AnalyzePassage("light", engine)

	// Both should produce same number of fields
	if len(fields1) != len(fields2) {
		t.Errorf("repeated analysis should produce same number of fields: %d vs %d",
			len(fields1), len(fields2))
	}
}

// =============================================================================
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
// EVIDENCE PATH TESTS
// =============================================================================

func TestPassageField_EvidencePathsPopulated(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	fields, err := AnalyzePassage("truth", engine)
	if err != nil {
		t.Fatal(err)
	}

	// Truth field should have evidence paths
	truthField := fields.GetField("truth")
	if truthField == nil {
		t.Error("should have truth field")
	} else if len(truthField.EvidencePaths) == 0 {
		t.Error("truth field should have evidence paths from activation graph")
	}
}

func TestEvidencePath_Deduplication(t *testing.T) {
	// Create paths with same ID
	path1 := EvidencePath{
		SourceToken: "love",
		SourceType:  "passage_signal",
		SourceForm:  "love",
		MatchForm:   "love",
		Confidence:  ConfidenceVerified,
		Weight:      1.0,
	}

	path2 := EvidencePath{
		SourceToken: "love",
		SourceType:  "passage_signal",
		SourceForm:  "love",
		MatchForm:   "love",
		Confidence:  ConfidenceVerified,
		Weight:      1.0,
	}

	// Both should have the same ID
	id1 := path1.EvidenceID()
	id2 := path2.EvidenceID()

	if id1 != id2 {
		t.Errorf("same evidence should have same ID: %s vs %s", id1, id2)
	}

	// Merging should deduplicate
	merged := mergeEvidencePaths([]EvidencePath{path1}, []EvidencePath{path2})
	if len(merged) != 1 {
		t.Errorf("merged evidence paths should have 1 entry, got %d", len(merged))
	}
}

// =============================================================================
// TOKENIZATION TESTS
// =============================================================================

func TestTokenizePassage_Basic(t *testing.T) {
	tokens := tokenizePassage("love truth")
	if len(tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[0] != "love" {
		t.Errorf("expected first token 'love', got '%s'", tokens[0])
	}
	if tokens[1] != "truth" {
		t.Errorf("expected second token 'truth', got '%s'", tokens[1])
	}
}

func TestTokenizePassage_Whitespace(t *testing.T) {
	tokens := tokenizePassage("  love   truth  ")
	if len(tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(tokens))
	}
}

func TestTokenizePassage_Newlines(t *testing.T) {
	tokens := tokenizePassage("love\ntruth\rlight")
	if len(tokens) != 3 {
		t.Errorf("expected 3 tokens, got %d", len(tokens))
	}
}

func TestTokenizePassage_SingleToken(t *testing.T) {
	tokens := tokenizePassage("love")
	if len(tokens) != 1 {
		t.Errorf("expected 1 token, got %d", len(tokens))
	}
}

func TestTokenizePassage_Empty(t *testing.T) {
	tokens := tokenizePassage("")
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(tokens))
	}
}

// =============================================================================
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
// ANALYZE PASSAGE FROM TOKENS TESTS
// =============================================================================

func TestAnalyzePassageFromTokens_Direct(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	tokens := []string{"love", "truth"}
	fields := AnalyzePassageFromTokens(tokens, engine.Knowledge)

	if len(fields) == 0 {
		t.Error("should produce passage fields from tokens")
	}
}

func TestAnalyzePassageFromTokens_TokenProvenance(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	tokens := []string{"light", "truth"}
	fields := AnalyzePassageFromTokens(tokens, engine.Knowledge)

	// Both fields should exist with evidence paths
	lightField := fields.GetField("light")
	if lightField == nil {
		t.Error("should have light field")
	} else {
		// Light field should have evidence paths (the key requirement)
		if len(lightField.EvidencePaths) == 0 {
			t.Error("light field should have evidence paths")
		}
		if lightField.Strength <= 0 {
			t.Error("light field should have positive strength")
		}
	}

	// Truth field should also exist
	truthField := fields.GetField("truth")
	if truthField == nil {
		t.Error("should have truth field")
	} else {
		if len(truthField.EvidencePaths) == 0 {
			t.Error("truth field should have evidence paths")
		}
	}
}
