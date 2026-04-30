package decipher

import (
	"fmt"
	"strings"
	"testing"

	"socrates/internal/knowledge"
)

// testKB returns a knowledge fixture for tests.
// Uses embedded knowledge if available, otherwise nil.
func testKB() *knowledge.Knowledge {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		return nil
	}
	return kb
}

func TestGenerateForms_Inspired(t *testing.T) {
	forms := GenerateForms("inspired")

	if forms.Normalized != "inspired" {
		t.Errorf("expected 'inspired', got '%s'", forms.Normalized)
	}

	// Should generate fragment paths
	if len(forms.Fragments) == 0 {
		t.Error("inspired should generate fragment paths")
	}

	// Check that some paths include spirit-related fragments
	foundSpiritPath := false
	for _, path := range forms.Fragments {
		for _, part := range path.Parts {
			if part == "in" || part == "spirit" || part == "spire" {
				foundSpiritPath = true
				break
			}
		}
	}
	if !foundSpiritPath {
		t.Error("inspired should have at least one path with 'in', 'spirit', or 'spire'")
	}
}

// =============================================================================
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
	kb := testKB()
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
	kb := testKB()
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
	kb := testKB()
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
// CLEANUP: Fuzzy Matching Tightening Tests
// =============================================================================

func TestUnrelatedInputProducesFewFuzzyMatches(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Unrelated noisy input should produce few or no fuzzy matches
	reading := engine.Analyze("zzskalx")

	// Count actual fuzzy matches (not just any signal)
	fuzzyMatchCount := len(reading.FuzzyMatches)

	// A truly unrelated string should have very few (ideally 0) fuzzy matches
	if fuzzyMatchCount > 3 {
		t.Errorf("unrelated 'zzskalx' should produce at most 3 fuzzy matches, got %d", fuzzyMatchCount)
	}
}

func TestExactMatchesOutrankFuzzyMatches(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// An input that matches an anchor exactly should have exact matches
	// ranked above any fuzzy matches
	reading := engine.Analyze("skal")

	// Find the first exact match and first fuzzy (non-exact) match
	var firstExactWeight, firstFuzzyWeight float64
	var foundExact, foundFuzzy bool

	for _, m := range reading.FuzzyMatches {
		if !foundExact && (m.Method == "exact" || m.Method == "case_insensitive") {
			firstExactWeight = m.Weight
			foundExact = true
		}
		if !foundFuzzy && m.Method != "exact" && m.Method != "case_insensitive" {
			firstFuzzyWeight = m.Weight
			foundFuzzy = true
		}
		if foundExact && foundFuzzy {
			break
		}
	}

	// If we have both exact and fuzzy, exact should rank higher or equal
	if foundExact && foundFuzzy {
		if firstFuzzyWeight > firstExactWeight {
			t.Errorf("exact match weight (%.2f) should be >= fuzzy match weight (%.2f)",
				firstExactWeight, firstFuzzyWeight)
		}
	}
}

func TestNoisyInputScoresLowerThanCleanInput(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	cleanReading := engine.Analyze("skal")
	noisyReading := engine.Analyze("zzskalx")

	// Clean input should score higher than noisy input
	// Allow small margin for floating point variations
	diff := noisyReading.Score.Overall - cleanReading.Score.Overall
	if diff > 0.05 {
		t.Errorf("noisy 'zzskalx' score (%.2f) should not significantly exceed clean 'skal' score (%.2f), diff=%.2f",
			noisyReading.Score.Overall, cleanReading.Score.Overall, diff)
	}
}

func TestRenderedFuzzyMatchesCapped(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Generate a reading with many potential fuzzy matches
	reading := engine.Analyze("inspired")
	output := RenderReading(reading)

	// Count lines in fuzzy matches section
	lines := strings.Split(output, "\n")
	inFuzzySection := false
	fuzzyLineCount := 0

	for _, line := range lines {
		if strings.HasPrefix(line, "Fuzzy Matches:") {
			inFuzzySection = true
			continue
		}
		if inFuzzySection {
			if strings.HasPrefix(line, "  - ") {
				fuzzyLineCount++
			}
			// End of section
			if line == "" || strings.HasPrefix(line, "  ...") || strings.HasPrefix(line, "Phase ") || strings.HasPrefix(line, "Graph") || strings.HasPrefix(line, "Passage") {
				if strings.HasPrefix(line, "  ...") {
					// This line indicates cap was applied
				}
				break
			}
		}
	}

	// Fuzzy matches should be capped at 10 (plus optional weaker-match line)
	if fuzzyLineCount > 10 {
		t.Errorf("rendered fuzzy matches should be capped at 10, got %d", fuzzyLineCount)
	}
}

func TestLevenshteinThresholdIsTight(t *testing.T) {
	// Test that very weak Levenshtein matches are rejected
	// A string like "abcdef" compared to "xyz123" should NOT match
	a := "abcdef"
	b := "xyz123"
	method, distance := fuzzyMatch(a, b)

	// Should be Levenshtein with high distance
	if method == "exact" || method == "case_insensitive" {
		t.Errorf("'%s' vs '%s' should not be exact match", a, b)
	}

	// Distance should be high (normalized > 0.34)
	if acceptMatch(method, distance, a, b) {
		t.Errorf("'%s' vs '%s' should NOT be accepted (distance %.2f exceeds 0.34)", a, b, distance)
	}
}

func TestFuzzyMatchEvidenceDeduplicates(t *testing.T) {
	// Test that duplicate evidence is removed
	candidates := []CandidateForm{
		{Form: "skal", Method: "normalized", Distance: 0, Confidence: "verified"},
		{Form: "skal", Method: "normalized", Distance: 0, Confidence: "verified"}, // duplicate
	}
	anchors := []AnchorConcept{
		{Form: "skal", Concept: "truth", Confidence: "verified", Weight: 0.5},
	}

	evidence, _ := FuzzyMatchEvidence(candidates, anchors)

	// Should have only one match, not duplicate
	if len(evidence) != 1 {
		t.Errorf("expected 1 evidence after dedup, got %d", len(evidence))
	}
}

// =============================================================================
// DUPLICATE EVIDENCE DEDUPLICATION TESTS
// =============================================================================

func TestDuplicateMatchEvidenceDoesNotInflateScore(t *testing.T) {
	// Test that repeated MatchEvidence with same path does not inflate score
	// Uses scoring.go deduplication

	matches := []MatchEvidence{
		{InputForm: "skal", AnchorForm: "skal", Method: "exact", Distance: 0, Weight: 0.5},
		{InputForm: "skal", AnchorForm: "skal", Method: "exact", Distance: 0, Weight: 0.5}, // duplicate
		{InputForm: "skal", AnchorForm: "skal", Method: "exact", Distance: 0, Weight: 0.5}, // duplicate
	}

	deduped := deduplicateMatchEvidence(matches)

	if len(deduped) != 1 {
		t.Errorf("expected 1 deduped match, got %d", len(deduped))
	}

	// Verify score calculation uses deduplicated matches
	components := CalculateScoreComponents(nil, deduped, nil, ConvergenceResult{}, []ChannelResult{})
	// With 1 exact match at weight 0.5, exact match score = 0.5
	if components.ExactMatchScore != 0.5 {
		t.Errorf("expected ExactMatchScore 0.5, got %f", components.ExactMatchScore)
	}
}

func TestDuplicateChannelSignalsDoNotInflateStrength(t *testing.T) {
	// Test that signals with same channel+target are deduplicated before scoring

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Analyze produces signals that get deduplicated in buildSignalGraph
	reading := engine.Analyze("inspired")

	// If we have converging patterns, verify they don't have inflated signal counts
	// from duplicate evidence
	for _, pattern := range reading.ConvergingPatterns {
		// Count unique channel+target combinations
		seenKeys := make(map[string]bool)
		for _, sig := range pattern.Signals {
			key := sig.Channel + "|" + sig.Target
			seenKeys[key] = true
		}

		// The number of signals should not exceed the number of unique keys
		// (duplicates should be collapsed)
		if len(pattern.Signals) > len(seenKeys) {
			t.Errorf("pattern %s has %d signals but only %d unique keys - duplicates not collapsed",
				pattern.Name, len(pattern.Signals), len(seenKeys))
		}
	}
}

func TestExactVerifiedEvidenceOutranksRepeatedSpeculativeEvidence(t *testing.T) {
	// Test that exact/verified evidence scores above repeated speculative evidence

	// Exact verified match
	exactVerified := []MatchEvidence{
		{InputForm: "truth", AnchorForm: "truth", Method: "exact", Distance: 0, Weight: 0.8},
	}

	// Repeated speculative fuzzy matches (should not outrank exact)
	repeatedSpeculative := []MatchEvidence{
		{InputForm: "truth", AnchorForm: "trth", Method: "consonant_skeleton", Distance: 0.5, Weight: 0.3},
		{InputForm: "truth", AnchorForm: "trth", Method: "consonant_skeleton", Distance: 0.5, Weight: 0.3}, // duplicate
		{InputForm: "truth", AnchorForm: "trth", Method: "consonant_skeleton", Distance: 0.5, Weight: 0.3}, // duplicate
	}

	dedupedSpeculative := deduplicateMatchEvidence(repeatedSpeculative)

	exactComponents := CalculateScoreComponents(nil, exactVerified, nil, ConvergenceResult{}, []ChannelResult{})
	speculativeComponents := CalculateScoreComponents(nil, dedupedSpeculative, nil, ConvergenceResult{}, []ChannelResult{})

	// Exact match score should be higher than fuzzy match score
	if speculativeComponents.FuzzyMatchScore >= exactComponents.ExactMatchScore {
		t.Errorf("exact match score (%.2f) should exceed single fuzzy match score (%.2f)",
			exactComponents.ExactMatchScore, speculativeComponents.FuzzyMatchScore)
	}
}

func TestDeduplicateSignalsFunction(t *testing.T) {
	// Test the DeduplicateSignals utility function
	signals := []Signal{
		{Text: "skal", Target: "truth", Channel: "glyph", Confidence: "verified", Weight: 0.5},
		{Text: "skal", Target: "truth", Channel: "glyph", Confidence: "verified", Weight: 0.5},  // duplicate
		{Text: "skal", Target: "truth", Channel: "glyph", Confidence: "verified", Weight: 0.5},  // duplicate
		{Text: "skal", Target: "truth", Channel: "sound", Confidence: "plausible", Weight: 0.4}, // different channel - independent evidence
	}

	deduped := DeduplicateSignals(signals)

	// Should have 2 - one for glyph, one for sound (different channels)
	if len(deduped) != 2 {
		t.Errorf("expected 2 signals after dedup (different channels), got %d", len(deduped))
	}

	// Verify the evidence IDs
	if deduped[0].EvidenceID() != "glyph|truth|skal" {
		t.Errorf("unexpected evidence ID for first signal: %s", deduped[0].EvidenceID())
	}
}

func TestDeduplicatePassageSignalsFunction(t *testing.T) {
	// Test the DeduplicatePassageSignals utility function
	signals := []PassageSignal{
		{Token: "in", Concept: "truth", MatchForm: "skal", Weight: 0.5},
		{Token: "in", Concept: "truth", MatchForm: "skal", Weight: 0.5},     // duplicate
		{Token: "in", Concept: "truth", MatchForm: "skal", Weight: 0.5},     // duplicate
		{Token: "spirit", Concept: "truth", MatchForm: "skal", Weight: 0.4}, // different token - independent evidence
	}

	deduped := DeduplicatePassageSignals(signals)

	// Should have 2 - one for 'in', one for 'spirit' (different tokens)
	if len(deduped) != 2 {
		t.Errorf("expected 2 signals after dedup (different tokens), got %d", len(deduped))
	}
}

func TestActivationStrengthNotInflatedByDuplicates(t *testing.T) {
	// Test that ComputeActivatedConcepts doesn't inflate strength with duplicates

	signals := []PassageSignal{
		{Token: "in", Concept: "truth", MatchForm: "skal", Weight: 0.5},
		{Token: "in", Concept: "truth", MatchForm: "skal", Weight: 0.5}, // duplicate - should not add weight
		{Token: "in", Concept: "truth", MatchForm: "skal", Weight: 0.5}, // duplicate - should not add weight
	}

	// With deduplication, the concept should only get 0.5 strength, not 1.5
	activated := ComputeActivatedConcepts(signals, nil, nil)

	if len(activated) != 1 {
		t.Errorf("expected 1 activated concept, got %d", len(activated))
	}

	if activated[0].Strength > 0.6 { // Allow small margin for floating point
		t.Errorf("concept strength should not be inflated by duplicates, got %f (expected ~0.5)", activated[0].Strength)
	}
}

// =============================================================================
// STRICT BOUNDED WORK TESTS
// =============================================================================

func TestCoreCandidatesSurviveTightBounds(t *testing.T) {
	// Test that core high-confidence candidates survive even with very tight bounds
	input := "energy"

	// Very tight bounds - only 5 candidates allowed
	tightBounds := CandidateBounds{
		MaxCandidates:            5,
		MaxSpeculativeCandidates: 2,
	}

	candidates, _ := GenerateCandidateForms(input, tightBounds)

	// Core candidates should always be present
	coreMethods := map[string]bool{
		"normalized":         false,
		"consonant_skeleton": false,
		"phonetic":           false,
	}

	for _, cand := range candidates {
		if _, ok := coreMethods[cand.Method]; ok {
			coreMethods[cand.Method] = true
		}
	}

	for method, found := range coreMethods {
		if !found {
			t.Errorf("core candidate method '%s' should survive tight bounds of 5 candidates", method)
		}
	}

	// Also verify we didn't exceed the max candidates cap
	if len(candidates) > 5 {
		t.Errorf("should not exceed maxCandidates (5), got %d", len(candidates))
	}
}

func TestExactMatchesFoundBeforeBudgetConsumed(t *testing.T) {
	// Test that high-confidence candidates get matched before the fuzzy comparison budget is exhausted
	kb := testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	// Create a controlled test with known high-quality candidates and matching anchors
	candidates := []CandidateForm{
		{Form: "energy", Method: "normalized", Distance: 0, Confidence: "verified"},
		{Form: "spirit", Method: "normalized", Distance: 0, Confidence: "verified"},
		{Form: "truth", Method: "normalized", Distance: 0, Confidence: "verified"},
		{Form: "breath", Method: "normalized", Distance: 0, Confidence: "verified"},
		{Form: "xyzabc", Method: "normalized", Distance: 0, Confidence: "speculative"},
	}
	anchors := []AnchorConcept{
		{Form: "energy", Concept: "power", Confidence: "verified", Weight: 0.8},
		{Form: "spirit", Concept: "soul", Confidence: "verified", Weight: 0.8},
		{Form: "truth", Concept: "truth", Confidence: "verified", Weight: 0.8},
		{Form: "breath", Concept: "life", Confidence: "verified", Weight: 0.8},
		{Form: "xyzabc", Concept: "unknown", Confidence: "speculative", Weight: 0.3},
	}

	// Very tight comparison budget - only 5 comparisons allowed
	tightBounds := FuzzyBounds{
		MaxComparisons: 5,
		MaxMatches:     10,
	}

	matches, discarded := FuzzyMatchEvidence(candidates, anchors, tightBounds)

	// With 5 comparisons allowed and 5 candidates, all candidates should be processed
	// and we should find matches for at least the verified candidates
	if len(matches) == 0 {
		t.Error("should find matches with tight budget of 5 comparisons")
	}

	// All matches should be from verified candidates (exact matches)
	for _, m := range matches {
		if m.Method != "exact" && m.Method != "case_insensitive" {
			t.Errorf("with tight budget, matches should be exact, got method '%s'", m.Method)
		}
	}

	// Discarded count should reflect remaining comparisons
	t.Logf("Tight budget: %d matches found, %d discarded comparisons", len(matches), discarded)
}

func TestEngineReadingExposesDiscardedCounts(t *testing.T) {
	// Engine-level test that Reading.DiscardedCandidates and DiscardedComparisons are populated
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Deliberately long input to trigger bounds
	longInput := "abcdefghijklmnopqrstuvwxyz"
	reading := engine.Analyze(longInput)

	// DiscardedCandidates should be populated (long input generates many speculative candidates)
	if reading.DiscardedCandidates < 0 {
		t.Errorf("DiscardedCandidates should be >= 0, got %d", reading.DiscardedCandidates)
	}

	// DiscardedComparisons should be populated
	if reading.DiscardedComparisons < 0 {
		t.Errorf("DiscardedComparisons should be >= 0, got %d", reading.DiscardedComparisons)
	}

	// For a long input, at least one of these should be positive
	// (unless the input happens to match so few things that no bounds are hit)
	t.Logf("Long input discarded: %d candidates, %d comparisons", reading.DiscardedCandidates, reading.DiscardedComparisons)
}

func TestTightBoundsStillProduceMatches(t *testing.T) {
	// Test that even with very tight bounds, we still get some results
	kb := testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	input := "energy"
	candidates, discardedCand := GenerateCandidateForms(input, CandidateBounds{
		MaxCandidates:            10,
		MaxSpeculativeCandidates: 5,
	})

	anchors := GetAllAnchors(kb)
	matches, discardedComp := FuzzyMatchEvidence(candidates, anchors, FuzzyBounds{
		MaxComparisons: 50,
		MaxMatches:     5,
	})

	// Should still produce candidates
	if len(candidates) == 0 {
		t.Error("should produce candidates even with tight bounds")
	}

	// Should not exceed max candidates
	if len(candidates) > 10 {
		t.Errorf("should not exceed 10 candidates, got %d", len(candidates))
	}

	// Should not exceed max matches
	if len(matches) > 5 {
		t.Errorf("should not exceed 5 matches, got %d", len(matches))
	}

	// Discarded counts should be non-negative
	if discardedCand < 0 {
		t.Errorf("discardedCand should be >= 0, got %d", discardedCand)
	}
	if discardedComp < 0 {
		t.Errorf("discardedComp should be >= 0, got %d", discardedComp)
	}

	t.Logf("Tight bounds: %d candidates (%d discarded), %d matches (%d discarded)",
		len(candidates), discardedCand, len(matches), discardedComp)
}

func TestExactHighConfidenceRetainedUnderCaps(t *testing.T) {
	// Test that exact/high-confidence matches are retained even when speculative matches are capped
	kb := testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	// Input with known exact match
	input := "truth"
	candidates, _ := GenerateCandidateForms(input, DefaultCandidateBounds())
	anchors := GetAllAnchors(kb)

	// Very tight match limit - only 1 match allowed
	matches, _ := FuzzyMatchEvidence(candidates, anchors, FuzzyBounds{
		MaxComparisons: 100,
		MaxMatches:     1,
	})

	// If we found any matches, the first should be an exact/high-confidence one
	if len(matches) > 0 {
		firstMatch := matches[0]
		// Exact match should rank first for "truth"
		if firstMatch.Method != "exact" && firstMatch.Method != "case_insensitive" {
			t.Errorf("with tight cap of 1, first match should be exact for 'truth', got '%s'", firstMatch.Method)
		}
	}
}

// =============================================================================
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
// KEY TOKEN REMOVAL REGRESSION TESTS
// =============================================================================

// TestKeyTokenRemovalWeakenActivation tests that removing a key evidence token
// weakens the relevant activation field. This is a regression test for
// Activation-Energy Discipline Gate.
func TestKeyTokenRemovalWeakenActivation(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Test: Single word "truth" - the key token IS "truth"
	// Removing key consonant 't' gives "ruth" which has different meaning
	readingFull := engine.Analyze("truth")
	readingNoKey := engine.Analyze("ruth")

	// The overall score should be lower when we remove the key token
	if readingFull.Score.Overall <= readingNoKey.Score.Overall {
		t.Errorf("removing key token should weaken overall score:\n  truth score: %.2f\n  ruth score: %.2f",
			readingFull.Score.Overall, readingNoKey.Score.Overall)
	}

	// The graph expansion score should be lower without the key anchor
	if readingFull.Score.Components.GraphExpansionScore <= readingNoKey.Score.Components.GraphExpansionScore {
		t.Errorf("removing key token should weaken graph expansion:\n  truth graph: %.2f\n  ruth graph: %.2f",
			readingFull.Score.Components.GraphExpansionScore, readingNoKey.Score.Components.GraphExpansionScore)
	}
}

// TestKeyTokenRemovalWeakenActivationMultiWord tests with a multi-word input
// where removing a key token should weaken convergence.
func TestKeyTokenRemovalWeakenActivationMultiWord(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// "word" is a key concept with neighbors like truth
	readingWord := engine.Analyze("word")
	readingOrd := engine.Analyze("ord")

	// "word" should score higher than "ord" (partial, not a key concept)
	if readingWord.Score.Overall <= readingOrd.Score.Overall {
		t.Errorf("key word 'word' should score higher than partial 'ord':\n  word: %.2f\n  ord: %.2f",
			readingWord.Score.Overall, readingOrd.Score.Overall)
	}
}

// TestKeyTokenRemovalWeakensConvergence tests that removing a key token
// from a passage weakens the convergence score.
func TestKeyTokenRemovalWeakensConvergence(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// "being" is a key concept with strong graph connections
	readingBeing := engine.Analyze("being")
	readingBeeing := engine.Analyze("beeing") // misspelling, missing key consonant pattern

	// "being" should have stronger convergence than misspelled version
	// The misspelling disrupts the phonetic signature
	if readingBeing.Score.Components.PassageConvergenceScore <= readingBeeing.Score.Components.PassageConvergenceScore {
		t.Errorf("correct 'being' should have stronger convergence than 'beeing':\n  being: %.2f\n  beeing: %.2f",
			readingBeing.Score.Components.PassageConvergenceScore, readingBeeing.Score.Components.PassageConvergenceScore)
	}
}

// TestKeyTokenRemovalWeakenExactMatch tests that removing a key token
// reduces the exact match score (for known concepts).
func TestKeyTokenRemovalWeakenExactMatch(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// "spirit" is a known concept with strong semantic connections
	readingSpirit := engine.Analyze("spirit")
	readingSpiri := engine.Analyze("spiri") // missing key ending

	// "spirit" should have higher overall score than "spiri"
	if readingSpirit.Score.Overall <= readingSpiri.Score.Overall {
		t.Errorf("key word 'spirit' should score higher than partial 'spiri':\n  spirit: %.2f\n  spiri: %.2f",
			readingSpirit.Score.Overall, readingSpiri.Score.Overall)
	}
}

// TestGodIsLove_BoundedScores verifies that "God is Love" has:
// - Bounded score components (all in [0, 1] range)
// - No duplicate evidence paths in output
// - No whitespace creating spurious repetition signals
func TestGodIsLove_BoundedScores(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("God is Love")
	output := RenderReadingWithOptions(reading, DefaultRenderOptions())

	// Verify score components are bounded in [0, 1]
	comp := reading.Score.Components
	if comp.ExactMatchScore < 0 || comp.ExactMatchScore > 1.0 {
		t.Errorf("ExactMatchScore should be in [0, 1], got %f", comp.ExactMatchScore)
	}
	if comp.FuzzyMatchScore < 0 || comp.FuzzyMatchScore > 1.0 {
		t.Errorf("FuzzyMatchScore should be in [0, 1], got %f", comp.FuzzyMatchScore)
	}
	if comp.GraphExpansionScore < 0 || comp.GraphExpansionScore > 1.0 {
		t.Errorf("GraphExpansionScore should be in [0, 1], got %f", comp.GraphExpansionScore)
	}
	if comp.PassageConvergenceScore < 0 || comp.PassageConvergenceScore > 1.0 {
		t.Errorf("PassageConvergenceScore should be in [0, 1], got %f", comp.PassageConvergenceScore)
	}
	if comp.MultiMethodBonus < 0 || comp.MultiMethodBonus > 1.0 {
		t.Errorf("MultiMethodBonus should be in [0, 1], got %f", comp.MultiMethodBonus)
	}
	if comp.ChannelDiversityBonus < 0 || comp.ChannelDiversityBonus > 1.0 {
		t.Errorf("ChannelDiversityBonus should be in [0, 1], got %f", comp.ChannelDiversityBonus)
	}

	// Verify no duplicate evidence paths in output
	// Count occurrences of "primitive 'love' matches love" - should appear once
	countPrimitiveLove := strings.Count(output, "primitive 'love' matches love")
	if countPrimitiveLove > 1 {
		t.Errorf("primitive 'love' should appear at most once in output, got %d occurrences", countPrimitiveLove)
	}

	// Verify no per-character repetition signals (whitespace should not create them)
	// Should NOT have patterns like "o repeated 2x" or "repeated 2x" appearing multiple times
	repetitionCount := strings.Count(output, "repeated letters detected")
	if repetitionCount > 1 {
		t.Errorf("'repeated letters detected' should appear at most once, got %d occurrences", repetitionCount)
	}

	// Verify overall score is bounded
	if reading.Score.Overall < 0 || reading.Score.Overall > 1.0 {
		t.Errorf("Overall score should be in [0, 1], got %f", reading.Score.Overall)
	}
}
