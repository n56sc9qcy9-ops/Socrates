package decipher

import (
	"testing"
	"fmt"
	"strings"

	"socrates/internal/knowledge"
)


// BOUNDED CANDIDATE GENERATION TESTS
// =============================================================================

func TestGenerateCandidateForms_Bounded(t *testing.T) {
	// Test that long inputs produce bounded candidates
	longInput := "abcdefghijklmnop" // 16 chars
	candidates, discarded := GenerateCandidateForms(longInput, DefaultCandidateBounds())

	// Should not exceed MaxCandidates (50)
	if len(candidates) > 50 {
		t.Errorf("long input should produce at most 50 candidates, got %d", len(candidates))
	}

	// For long input, some candidates should be discarded
	// A 16-char input would generate ~16 deletions + 416 insertions + ~400 substitutions = ~832 edit variants
	// But we cap at 10
	if discarded == 0 {
		t.Error("long input should have discarded candidates (many edit variants should be skipped)")
	}
}

func TestGenerateCandidateForms_CoreRetained(t *testing.T) {
	// Test that high-confidence core candidates survive bounds
	candidates, _ := GenerateCandidateForms("energy", DefaultCandidateBounds())

	// Should always have normalized form
	foundNormalized := false
	// Should have consonant skeleton (energy -> nrg, not "nrgy")
	foundSkeleton := false
	// Should have phonetic
	foundPhonetic := false

	for _, cand := range candidates {
		if cand.Form == "energy" && cand.Method == "normalized" {
			foundNormalized = true
		}
		// Consonant skeleton of 'energy' removes vowels: e-n-e-r-g-y -> n-r-g
		if cand.Form == "nrg" && cand.Method == "consonant_skeleton" {
			foundSkeleton = true
		}
		if cand.Method == "phonetic" {
			foundPhonetic = true
		}
	}

	if !foundNormalized {
		t.Error("normalized 'energy' should be retained regardless of bounds")
	}
	if !foundSkeleton {
		t.Error("consonant skeleton 'nrg' should be retained regardless of bounds")
	}
	if !foundPhonetic {
		t.Error("phonetic variant should be retained regardless of bounds")
	}
}

func TestFuzzyMatchEvidence_Bounded(t *testing.T) {
	// Create many candidates to stress-test bounds
	candidates := make([]CandidateForm, 100)
	for i := range candidates {
		candidates[i] = CandidateForm{
			Form:       "cand" + fmt.Sprintf("%d", i),
			Method:     "normalized",
			Distance:   0,
			Confidence: "verified",
		}
	}

	anchors := make([]AnchorConcept, 20)
	for i := range anchors {
		anchors[i] = AnchorConcept{
			Form:       "anchor" + fmt.Sprintf("%d", i),
			Concept:    "concept",
			Confidence: "verified",
			Weight:     0.5,
		}
	}

	matches, discarded := FuzzyMatchEvidence(candidates, anchors, DefaultFuzzyBounds())

	// Should not exceed MaxMatches (20)
	if len(matches) > 20 {
		t.Errorf("should produce at most 20 matches, got %d", len(matches))
	}

	// With 100 candidates x 20 anchors = 2000 comparisons, and MaxComparisons = 500,
	// we should have discarded comparisons - assert this
	if discarded < 1000 {
		t.Errorf("should discard at least 1000 comparisons when 100 candidates x 20 anchors would exceed 500 cap, got %d", discarded)
	}
}

func TestFuzzyMatchEvidence_DiscardedCount(t *testing.T) {
	// Test that discarded comparison count is populated
	var kb *knowledge.Knowledge = testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	// Create a long input with many candidates
	longInput := "abcdefghijklmnopqrstuvwxyz"
	candidates, _ := GenerateCandidateForms(longInput, DefaultCandidateBounds())
	anchors := GetAllAnchors(kb)

	_, discarded := FuzzyMatchEvidence(candidates, anchors, DefaultFuzzyBounds())

	// Discarded count should be non-negative
	if discarded < 0 {
		t.Errorf("discarded comparisons should be >= 0, got %d", discarded)
	}

	// For long random input, we expect significant discarding
	// (many candidates will be compared but capped)
	if len(candidates) > 30 && len(anchors) > 50 && discarded == 0 {
		t.Error("long input with many candidates should have discarded comparisons when cap is reached")
	}
}

func TestGenerateForms_Energy(t *testing.T) {
	forms := GenerateForms("energy")

	if forms.Normalized != "energy" {
		t.Errorf("expected 'energy', got '%s'", forms.Normalized)
	}

	// Should generate multiple fragment paths
	if len(forms.Fragments) < 2 {
		t.Error("energy should generate multiple fragment paths")
	}

	// Should have phonetic keys
	if len(forms.PhoneticKeys) == 0 {
		t.Error("energy should generate phonetic keys")
	}
}

func TestGenerateForms_UnknownLatin(t *testing.T) {
	forms := GenerateForms("xyzqrs")

	if forms.Normalized != "xyzqrs" {
		t.Errorf("expected 'xyzqrs', got '%s'", forms.Normalized)
	}

	// Should still produce forms
	if len(forms.PhoneticKeys) == 0 {
		t.Error("unknown word should still produce phonetic forms")
	}

	// Should still have glyph analysis
	channels := RunAllChannels(forms, testKB())
	if len(channels) == 0 {
		t.Error("should produce channel results")
	}
}

func TestGenerateForms_NonLatin(t *testing.T) {
	// Hebrew
	formsHe := GenerateForms("רוח")
	if formsHe.Script != ScriptHebrew {
		t.Errorf("Hebrew input should have ScriptHebrew, got %s", formsHe.Script)
	}

	// Devanagari
	formsDev := GenerateForms("आत्मन्")
	if formsDev.Script != ScriptDevanagari {
		t.Errorf("Devanagari input should have ScriptDevanagari, got %s", formsDev.Script)
	}

	// Han
	formsHan := GenerateForms("道")
	if formsHan.Script != ScriptHan {
		t.Errorf("Han input should have ScriptHan, got %s", formsHan.Script)
	}
}

func TestEngine_Analyze_Inspired(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	reading := engine.Analyze("inspired")

	if reading.Input != "inspired" {
		t.Errorf("expected input 'inspired', got '%s'", reading.Input)
	}

	// Should have multiple channels
	if len(reading.Channels) < 3 {
		t.Errorf("should have multiple channels, got %d", len(reading.Channels))
	}

	// Should produce signals
	totalSignals := 0
	for _, ch := range reading.Channels {
		totalSignals += len(ch.Signals)
	}
	if totalSignals == 0 {
		t.Error("should produce at least some signals")
	}

	// Should have a reading
	if reading.ConciseReading == "" {
		t.Error("should produce a concise reading")
	}
}

func TestEngine_Analyze_Energy(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	reading := engine.Analyze("energy")

	if reading.Input != "energy" {
		t.Errorf("expected input 'energy', got '%s'", reading.Input)
	}

	// Should have multiple fragment paths (not hardcoded)
	if len(reading.Forms.Fragments) < 2 {
		t.Error("energy should generate multiple fragment paths")
	}

	// Should have a score
	if reading.Score.Overall < 0 || reading.Score.Overall > 1 {
		t.Errorf("score should be in [0, 1], got %f", reading.Score.Overall)
	}
}

func TestEngine_Analyze_Ruach(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	reading := engine.Analyze("רוח")

	if reading.Forms.Script != ScriptHebrew {
		t.Errorf("should detect Hebrew, got %s", reading.Forms.Script)
	}

	// Should have Hebrew-related signals
	foundBreath := false
	for _, ch := range reading.Channels {
		for _, sig := range ch.Signals {
			if sig.Target == "breath" || sig.Target == "spirit" {
				foundBreath = true
				break
			}
		}
	}
	if !foundBreath {
		t.Error("רוח should produce breath/spirit signals")
	}
}

func TestEngine_Analyze_Praana(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	// Test with Devanagari प्राण (prāṇa)
	reading := engine.Analyze("प्राण")

	if reading.Forms.Script != ScriptDevanagari {
		t.Errorf("should detect Devanagari, got %s", reading.Forms.Script)
	}

	// Should have breath/life-force signals from exact ScriptWord match
	foundBreath := false
	for _, ch := range reading.Channels {
		for _, sig := range ch.Signals {
			if sig.Target == "breath" || sig.Target == "life-force" {
				foundBreath = true
				break
			}
		}
	}
	if !foundBreath {
		t.Error("प्राण should produce breath/life-force signals via ScriptWord channel")
	}
}

func TestEngine_Analyze_QiTraditional(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	// Test with traditional Chinese 氣
	reading := engine.Analyze("氣")

	if reading.Forms.Script != ScriptHan {
		t.Errorf("should detect Han, got %s", reading.Forms.Script)
	}

	// Should have breath/qi signals from exact ScriptWord match
	foundBreath := false
	for _, ch := range reading.Channels {
		for _, sig := range ch.Signals {
			if sig.Target == "breath" || sig.Target == "qi" {
				foundBreath = true
				break
			}
		}
	}
	if !foundBreath {
		t.Error("氣 should produce breath/qi signals via ScriptWord channel")
	}
}

func TestEngine_Analyze_QiSimplified(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	// Test with simplified Chinese 气
	reading := engine.Analyze("气")

	if reading.Forms.Script != ScriptHan {
		t.Errorf("should detect Han, got %s", reading.Forms.Script)
	}

	// Should have breath/qi signals from exact ScriptWord match
	foundBreath := false
	for _, ch := range reading.Channels {
		for _, sig := range ch.Signals {
			if sig.Target == "breath" || sig.Target == "qi" {
				foundBreath = true
				break
			}
		}
	}
	if !foundBreath {
		t.Error("气 should produce breath/qi signals via ScriptWord channel")
	}
}

func TestEngine_Analyze_Skal(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	reading := engine.Analyze("skal")

	// Should produce some form of resonance (candidates, fuzzy matches, etc.)
	// We don't check for specific hardcoded targets anymore
	// Instead verify the generic evidence pipeline works

	if len(reading.Candidates) == 0 {
		t.Error("skal should generate candidate forms")
	}

	// Verify candidates include the normalized form and variants
	foundNormalized := false
	for _, cand := range reading.Candidates {
		if cand.Form == "skal" {
			foundNormalized = true
			break
		}
	}
	if !foundNormalized {
		t.Error("skal should produce normalized 'skal' candidate")
	}
}

func TestEngine_Analyze_Skall(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	reading := engine.Analyze("skall")

	// Should produce candidate forms for fuzzy matching
	if len(reading.Candidates) == 0 {
		t.Error("skall should generate candidate forms")
	}

	// Verify generic evidence is produced (not hardcoded semantic targets)
	// The system should still produce signals through generic channels
	foundAnySignal := false
	for _, ch := range reading.Channels {
		if len(ch.Signals) > 0 {
			foundAnySignal = true
			break
		}
	}
	if !foundAnySignal {
		t.Error("skall should produce some generic signals through channels")
	}
}

func TestEngine_Analyze_Shell(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	reading := engine.Analyze("shell")

	// Should have shell/scale signals from fragment seed
	foundShell := false
	for _, ch := range reading.Channels {
		for _, sig := range ch.Signals {
			if sig.Target == "shell" || sig.Target == "scale" {
				foundShell = true
				break
			}
		}
	}
	if !foundShell {
		t.Error("shell should produce shell/scale signals via Fragment channel")
	}
}

func TestEngine_Analyze_Prana(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	reading := engine.Analyze("prana")

	// Should have breath/life signals
	foundBreath := false
	for _, ch := range reading.Channels {
		for _, sig := range ch.Signals {
			if sig.Target == "breath" || sig.Target == "life" {
				foundBreath = true
				break
			}
		}
	}
	if !foundBreath {
		t.Error("prana should produce breath/life signals")
	}
}

func TestEngine_Analyze_UnknownHebrew(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	// Test with an unknown Hebrew word (not רוח)
	reading := engine.Analyze("מלך") // melek (king) - not ruach

	if reading.Forms.Script != ScriptHebrew {
		t.Errorf("should detect Hebrew, got %s", reading.Forms.Script)
	}

	// Unknown Hebrew should NOT produce breath/spirit signals
	for _, ch := range reading.Channels {
		for _, sig := range ch.Signals {
			if sig.Target == "breath" || sig.Target == "spirit" {
				t.Errorf("unknown Hebrew מלך should NOT produce breath/spirit, but got: %s -> %s", sig.Text, sig.Target)
			}
		}
	}
}

func TestEngine_Analyze_UnknownDevanagari(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	// Test with an unknown Devanagari word (not प्राण)
	reading := engine.Analyze("कवि") // kavi (poet) - not prana

	if reading.Forms.Script != ScriptDevanagari {
		t.Errorf("should detect Devanagari, got %s", reading.Forms.Script)
	}

	// Unknown Devanagari should NOT produce prana/life-force signals
	for _, ch := range reading.Channels {
		for _, sig := range ch.Signals {
			if sig.Target == "prana" || sig.Target == "life-force" || sig.Target == "life" {
				t.Errorf("unknown Devanagari कवि should NOT produce prana/life-force, but got: %s -> %s", sig.Text, sig.Target)
			}
		}
	}
}

func TestEngine_Analyze_UnknownHan(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	// Test with an unknown Han character (not 道)
	reading := engine.Analyze("山") // shan (mountain)

	if reading.Forms.Script != ScriptHan {
		t.Errorf("should detect Han, got %s", reading.Forms.Script)
	}

	// Unknown Han should NOT produce path/dao signals
	for _, ch := range reading.Channels {
		for _, sig := range ch.Signals {
			if sig.Target == "path" || sig.Target == "dao" || sig.Target == "way" {
				t.Errorf("unknown Han 山 should NOT produce path/dao, but got: %s -> %s", sig.Text, sig.Target)
			}
		}
	}
}

func TestEngine_Analyze_Dao(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	reading := engine.Analyze("道")

	if reading.Forms.Script != ScriptHan {
		t.Errorf("道 should be Han script, got %s", reading.Forms.Script)
	}

	// Should have path/way signals
	foundPath := false
	for _, ch := range reading.Channels {
		for _, sig := range ch.Signals {
			if sig.Target == "path" || sig.Target == "way" {
				foundPath = true
				break
			}
		}
	}
	if !foundPath {
		t.Error("道 should produce path/way signals")
	}
}

func TestEngine_MultiChannelConvergence(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	reading := engine.Analyze("inspired")

	// Check that multiple channels can point to similar targets
	// If we have convergence, it should score higher than isolated signals
	hasConverging := len(reading.ConvergingPatterns) > 0

	// If there are weak signals only, the overall should be lower
	if !hasConverging && len(reading.WeakSignals) > 0 {
		// Weak signals only case - this is valid behavior
		t.Log("Only weak signals detected - valid for this input")
	}

	// Score should reflect the pattern strength
	if reading.Score.Overall < 0 || reading.Score.Overall > 1 {
		t.Errorf("score should be in [0, 1], got %f", reading.Score.Overall)
	}
}

func TestEngine_SpeculativeLabeling(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	reading := engine.Analyze("xyzqrs")

	// Unknown word - signals should be labeled appropriately
	// All signals should have a confidence category
	for _, ch := range reading.Channels {
		for _, sig := range ch.Signals {
			if sig.Confidence != ConfidenceVerified &&
				sig.Confidence != ConfidencePlausible &&
				sig.Confidence != ConfidenceSpeculative {
				t.Errorf("signal should have valid confidence, got '%s'", sig.Confidence)
			}
		}
	}
}

func TestEngine_NoHardcodedReadings(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Test that the system doesn't just return pre-written readings
	// It should generate different patterns for different inputs
	inspired := engine.Analyze("inspired")
	energy := engine.Analyze("energy")

	// The readings should be different
	if inspired.ConciseReading == energy.ConciseReading {
		t.Error("different words should produce different readings")
	}

	// The fragment paths should be different
	if len(inspired.Forms.Fragments) > 0 && len(energy.Forms.Fragments) > 0 {
		// Compare the first fragment paths
		insParts := strings.Join(inspired.Forms.Fragments[0].Parts, "|")
		enParts := strings.Join(energy.Forms.Fragments[0].Parts, "|")
		if insParts == enParts {
			// This could happen by chance but is unlikely
			t.Log("Fragment paths happened to match")
		}
	}
}

func TestFragmentSeeding(t *testing.T) {
	// Test that fragment seeds are reusable, not final readings
	// Uses knowledge-based lookup
	var kb *knowledge.Knowledge = testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	seeds := knowledgeBasedFragmentLookup("in", kb)
	if len(seeds) == 0 {
		t.Error("should find fragment seed for 'in' via knowledge base")
	}

	// The seed should have multiple lenses
	if len(seeds[0].Lenses) < 1 {
		t.Error("fragment seed should have at least one lens")
	}
}

func TestPrimitivesExist(t *testing.T) {
	// Test that primitive concepts exist in knowledge base
	var kb *knowledge.Knowledge = testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	prims := knowledgeBasedPrimitiveLookup("breath", kb)
	if len(prims) == 0 {
		t.Error("should find primitive 'breath' via knowledge base")
	}

	// Check for expected primitives
	expectedPrims := []string{"breath", "spirit", "word", "life", "light", "truth"}
	for _, ep := range expectedPrims {
		prims := knowledgeBasedPrimitiveLookup(ep, kb)
		if len(prims) == 0 {
			t.Errorf("missing expected primitive via knowledge base: %s", ep)
		}
	}
}

func TestRenderReading(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	reading := engine.Analyze("truth")
	output := RenderReadingWithOptions(reading, DebugRenderOptions())

	// Debug output should contain key sections
	if !strings.Contains(output, "Input:") {
		t.Error("rendered output should contain 'Input:'")
	}
	if !strings.Contains(output, "Generated Forms:") {
		t.Error("rendered output should contain 'Generated Forms:'")
	}
	if !strings.Contains(output, "Resonance Channels:") {
		t.Error("rendered output should contain 'Resonance Channels:'")
	}
	if !strings.Contains(output, "Resonance Score:") {
		t.Error("rendered output should contain 'Resonance Score:'")
	}
	if !strings.Contains(output, "Reading:") {
		t.Error("rendered output should contain 'Reading:'")
	}
}

func TestDetectScript(t *testing.T) {
	tests := []struct {
		input    string
		expected ScriptType
	}{
		{"hello", ScriptLatin},
		{"רוח", ScriptHebrew},
		{"आत्मन्", ScriptDevanagari},
		{"道", ScriptHan},
		{"مرحبا", ScriptArabic},
		{"Γεια", ScriptGreek},
	}

	for _, tc := range tests {
		script := DetectScript(tc.input)
		if script != tc.expected {
			t.Errorf("DetectScript(%q) = %s, want %s", tc.input, script, tc.expected)
		}
	}
}

func TestConsonantSkeleton(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"energy", "nrg"},
		{"inspired", "nsprd"},
		{"truth", "trth"},
		{"spirit", "sprt"},
	}

	for _, tc := range tests {
		result := consonantSkeleton(tc.input)
		if result != tc.expected {
			t.Errorf("consonantSkeleton(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

func TestPhoneticVariants(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		{"philosophy", "filo"},
		{"psychology", "psi"},
		{"energy", "energi"},
	}

	for _, tc := range tests {
		variants := phoneticVariants(tc.input)
		found := false
		for _, v := range variants {
			if strings.Contains(v, tc.contains) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("phoneticVariants(%q) should contain %q", tc.input, tc.contains)
		}
	}
}

// =============================================================================
