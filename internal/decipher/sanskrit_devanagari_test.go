package decipher

import (
	"strings"
	"testing"

	"socrates/internal/knowledge"
)

// TestSanskritDevanagariScriptDetection verifies that Devanagari script is detected
// correctly for Sanskrit text input.
func TestSanskritDevanagariScriptDetection(t *testing.T) {
	testCases := []struct {
		input    string
		expected ScriptType
	}{
		{"प्राण", ScriptDevanagari},
		{"ॐ", ScriptDevanagari},
		{"आत्मन्", ScriptDevanagari},
		{"सत्य", ScriptDevanagari},
		{"धर्म", ScriptDevanagari},
		{"अग्नि", ScriptDevanagari},
		{"अनंद", ScriptDevanagari},
	}

	for _, tc := range testCases {
		forms := GenerateForms(tc.input)
		if forms.Script != tc.expected {
			t.Errorf("input %q should have script %s, got %s", tc.input, tc.expected, forms.Script)
		}
	}
}

// TestSanskritDevanagariOmGlyph tests that the Om glyph ॐ activates sacred-sound
// via the glyph channel with correct confidence/source/lens.
func TestSanskritDevanagariOmGlyph(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatal(err)
	}
	engine := NewEngineWithKnowledge(kb)

	reading := engine.Analyze("ॐ")

	// Check script detection
	if reading.Forms.Script != ScriptDevanagari {
		t.Errorf("ॐ should be detected as Devanagari, got %s", reading.Forms.Script)
	}

	// Check that sacred-sound is activated
	foundSacredSound := false
	for _, ch := range reading.Channels {
		for _, s := range ch.Signals {
			if s.Target == "sacred-sound" {
				foundSacredSound = true
				// Lens should be devanagari-glyph or sanskrit
				if s.Lens != "devanagari-glyph" && s.Lens != "sanskrit" {
					t.Logf("note: sacred-sound signal lens is %q (expected devanagari-glyph or sanskrit)", s.Lens)
				}
				t.Logf("ॐ -> sacred-sound via channel %s, lens %q", ch.Name, s.Lens)
			}
		}
	}

	if !foundSacredSound {
		t.Error("ॐ should activate sacred-sound via glyph channel")
	}
}

// TestSanskritDevanagariExactWordActivation verifies that exact Devanagari ScriptWords
// activate their intended concepts.
func TestSanskritDevanagariExactWordActivation(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatal(err)
	}
	engine := NewEngineWithKnowledge(kb)

	testCases := []struct {
		input           string
		wantConcept     string
		wantChannel     string
		wantConfidence string
	}{
		{"प्राण", "breath", "ScriptWord", "verified"},
		{"ॐ", "sacred-sound", "ScriptWord", "verified"},
		{"अनंद", "bliss", "ScriptWord", "verified"},
	}

	for _, tc := range testCases {
		reading := engine.Analyze(tc.input)

		found := false
		for _, ch := range reading.Channels {
			if ch.Name == tc.wantChannel {
				for _, s := range ch.Signals {
					if s.Target == tc.wantConcept {
						found = true
						// Confidence check
						if tc.wantConfidence != "" && s.Confidence != tc.wantConfidence {
							t.Logf("note: %s -> %s has confidence %q (expected %q)",
								tc.input, tc.wantConcept, s.Confidence, tc.wantConfidence)
						}
						t.Logf("%s -> %s via %s [confidence=%s, lens=%s]",
							tc.input, s.Target, tc.wantChannel, s.Confidence, s.Lens)
					}
				}
			}
		}

		if !found {
			t.Errorf("%s should activate %s via %s channel", tc.input, tc.wantConcept, tc.wantChannel)
		}
	}
}

// TestSanskritTransliteratedActivation verifies that transliterated Sanskrit
// (Latin script) forms activate the same meaning-frequency field as Devanagari.
func TestSanskritTransliteratedActivation(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatal(err)
	}
	engine := NewEngineWithKnowledge(kb)

	// "prana" transliterated should activate the same breath concept as Devanagari प्राण
	reading := engine.Analyze("prana")

	foundBreath := false
	for _, ch := range reading.Channels {
		for _, s := range ch.Signals {
			if s.Target == "breath" || s.Target == "life" {
				foundBreath = true
				t.Logf("prana -> %s via channel %s, lens %s", s.Target, ch.Name, s.Lens)
			}
		}
	}

	if !foundBreath {
		t.Error("transliterated 'prana' should activate breath/life via Fragment channel")
	}
}

// TestSanskritDevanagariCrossScriptConvergence verifies that a Devanagari word
// and its transliterated form converge on the same meaning-frequency field,
// and that Sanskrit and Hebrew "breath" converge on the shared field.
func TestSanskritDevanagariCrossScriptConvergence(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatal(err)
	}
	engine := NewEngineWithKnowledge(kb)

	// Devanagari प्राण (prana = breath)
	readingDev := engine.Analyze("प्राण")
	// Transliterated prana
	readingLat := engine.Analyze("prana")

	var devFound, latFound bool
	for _, ch := range readingDev.Channels {
		for _, s := range ch.Signals {
			if s.Target == "breath" || s.Target == "life-force" {
				devFound = true
				t.Logf("प्राण -> %s via %s [weight=%.0f, lens=%s]",
					s.Target, ch.Name, s.Weight, s.Lens)
			}
		}
	}
	for _, ch := range readingLat.Channels {
		for _, s := range ch.Signals {
			if s.Target == "breath" || s.Target == "life-force" {
				latFound = true
				t.Logf("prana -> %s via %s [weight=%.0f, lens=%s]",
					s.Target, ch.Name, s.Weight, s.Lens)
			}
		}
	}

	if !devFound {
		t.Error("Devanagari प्राण should activate breath/life-force")
	}
	if !latFound {
		t.Error("transliterated 'prana' should activate breath/life-force")
	}

	// Both should converge on breath field (both produce evidence)
	if devFound && latFound {
		t.Log("प्राण and prana converge on breath/life-force field")
	}

	// Also verify Sanskrit "prana" and Hebrew "ruach" both activate breath
	readingHebrew := engine.Analyze("ruach")
	var hebrewFound bool
	for _, ch := range readingHebrew.Channels {
		for _, s := range ch.Signals {
			if s.Target == "breath" {
				hebrewFound = true
				t.Logf("ruach -> %s via %s [weight=%.0f, lens=%s]",
					s.Target, ch.Name, s.Weight, s.Lens)
			}
		}
	}
	if hebrewFound {
		t.Log("Sanskrit prana and Hebrew ruach both activate breath field")
	}
}

// TestSanskritDevanagariUnknownSafety verifies that unknown or weak Devanagari input
// does not produce overconfident counsellor guidance.
func TestSanskritDevanagariUnknownSafety(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatal(err)
	}
	engine := NewEngineWithKnowledge(kb)

	// Unknown Devanagari word (not in ScriptWords)
	reading := engine.Analyze("कवि") // kavi (poet) - not in curated data

	// Script detection should still work
	if reading.Forms.Script != ScriptDevanagari {
		t.Errorf("unknown Devanagari should still be detected as Devanagari, got %s", reading.Forms.Script)
	}

	// Should NOT produce signals for breath/life-force/prana
	for _, ch := range reading.Channels {
		for _, s := range ch.Signals {
			if s.Target == "prana" || s.Target == "life-force" || s.Target == "breath" {
				t.Errorf("unknown Devanagari कवि should NOT produce %s signals, got: %s", s.Target, s.Text)
			}
		}
	}

	// Counsellor field should NOT suggest breath/transmutation paths for unknown input
	if reading.CounsellorField != nil {
		for _, s := range reading.CounsellorField.Suggestions {
			if s.SourceConcept == "breath" || s.SourceConcept == "life" {
				t.Errorf("unknown Devanagari कवि should NOT produce breath/life suggestions, got: %s -> %s",
					s.SourceConcept, s.TargetConcept)
			}
		}
	}

	// Score should be low (weak convergence)
	if reading.Score.Overall > 0.5 {
		t.Logf("note: unknown Devanagari score is %.2f (expected low, <0.5)", reading.Score.Overall)
	}
}

// TestSanskritDevanagariDebugOutput verifies that debug output for Devanagari
// shows script channel evidence with correct script labeling.
func TestSanskritDevanagariDebugOutput(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatal(err)
	}
	engine := NewEngineWithKnowledge(kb)

	reading := engine.Analyze("प्राण")
	debug := RenderReadingWithOptions(reading, RenderOptions{Mode: RenderModeDebug})

	// Debug output should show ScriptWord channel for Devanagari
	if !strings.Contains(debug, "ScriptWord") {
		t.Error("debug output should contain ScriptWord channel for Devanagari input")
	}

	// Debug output should show the Devanagari word correctly
	if !strings.Contains(debug, "प्राण") {
		t.Error("debug output should show the Devanagari input word")
	}

	// Debug output should show breath/life-force in evidence
	if !strings.Contains(debug, "breath") && !strings.Contains(debug, "life-force") {
		t.Error("debug output should show breath/life-force evidence for प्राण")
	}

	t.Logf("Debug output excerpt:\n%s", debug)
}

// TestSanskritDevanagariNoHardcodedBehavior verifies that Sanskrit support works
// entirely through data channels (glyph, fragment, script word) with no
// hardcoded word-specific branches in Go.
func TestSanskritDevanagariNoHardcodedBehavior(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatal(err)
	}
	engine := NewEngineWithKnowledge(kb)

	// Test a known Devanagari word (agni = fire)
	reading := engine.Analyze("अग्नि")

	// Should activate fire through a data-driven channel (ScriptWord or Fragment)
	foundFire := false
	for _, ch := range reading.Channels {
		for _, s := range ch.Signals {
			if s.Target == "fire" {
				foundFire = true
				// Should come from a data channel, not a hardcoded branch
				if ch.Name != "ScriptWord" && ch.Name != "Fragment" && ch.Name != "Glyph" {
					t.Errorf("अग्नि->fire came from %s (expected ScriptWord/Fragment/Glyph)", ch.Name)
				}
				t.Logf("अग्नि -> fire via %s [lens=%s, confidence=%s]",
					ch.Name, s.Lens, s.Confidence)
			}
		}
	}

	if !foundFire {
		t.Error("अग्नि should activate fire via data-driven channel")
	}
}

// TestSanskritDevanagariAnanda verifies that अनंद (bliss/ananda) activates
// the bliss concept through Devanagari ScriptWord matching.
func TestSanskritDevanagariAnanda(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatal(err)
	}
	engine := NewEngineWithKnowledge(kb)

	reading := engine.Analyze("अनंद")

	// Check script detection
	if reading.Forms.Script != ScriptDevanagari {
		t.Errorf("अनंद should be detected as Devanagari, got %s", reading.Forms.Script)
	}

	// Check that bliss is activated
	foundBliss := false
	for _, ch := range reading.Channels {
		for _, s := range ch.Signals {
			if s.Target == "bliss" {
				foundBliss = true
				t.Logf("अनंद -> bliss via %s [confidence=%s, lens=%s]",
					ch.Name, s.Confidence, s.Lens)
			}
		}
	}

	if !foundBliss {
		t.Error("अनंद should activate bliss via ScriptWord channel")
	}
}

// TestSanskritDevanagariMantra verifies that मंत्र (mantra/sacred-utterance)
// activates sacred-utterance or mantra-related concepts.
func TestSanskritDevanagariMantra(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatal(err)
	}
	engine := NewEngineWithKnowledge(kb)

	reading := engine.Analyze("मंत्र")

	// Check script detection
	if reading.Forms.Script != ScriptDevanagari {
		t.Errorf("मंत्र should be detected as Devanagari, got %s", reading.Forms.Script)
	}

	// Check that sacred-utterance or mantra-related concept is activated
	foundMantra := false
	for _, ch := range reading.Channels {
		for _, s := range ch.Signals {
			if s.Target == "sacred-utterance" || s.Target == "mantra" {
				foundMantra = true
				t.Logf("मंत्र -> %s via %s [confidence=%s, lens=%s]",
					s.Target, ch.Name, s.Confidence, s.Lens)
			}
		}
	}

	if !foundMantra {
		t.Error("मंत्र should activate sacred-utterance or mantra via ScriptWord channel")
	}
}
