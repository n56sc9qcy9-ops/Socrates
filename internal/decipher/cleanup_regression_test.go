package decipher

import (
	"os"
	"strings"
	"testing"

	"socrates/internal/knowledge"
)

// TestNoSemanticTargetsInProduction verifies no hardcoded semantic targets exist in production code.
// This is a regression test to ensure semantic meanings are loaded from data, not Go code.
func TestNoSemanticTargetsInProduction(t *testing.T) {
	// These semantic targets must NOT appear in production Go files
	forbiddenTargets := []string{
		`Target: "breath"`,
		`Target: "ground"`,
		`Target: "spirit"`,
		`Target: "life"`,
		`Target: "inward"`,
		`Target: "being"`,
		`Target: "core"`,
		`Target: "sound"`,
	}

	// Production files to check
	prodFiles := []string{
		"activation.go",
		"candidate_generation.go",
		"channel_cross_language.go",
		"channel_fragment.go",
		"channel_glyph.go",
		"channel_scoring.go",
		"channel_script_word.go",
		"channel_sound.go",
		"channel_symbolic.go",
		"channels.go",
		"convergence.go",
		"discovery.go",
		"engine.go",
		"forms.go",
		"glyphs.go",
		"knowledge_bridge.go",
		"render.go",
		"scoring.go",
		"similarity.go",
		"types.go",
	}

	for _, file := range prodFiles {
		path := "internal/decipher/" + file
		content, err := os.ReadFile(path)
		if err != nil {
			t.Skipf("cannot read %s: %v", file, err)
		}
		contentStr := string(content)

		for _, target := range forbiddenTargets {
			if strings.Contains(contentStr, target) {
				t.Errorf("production file %s contains forbidden semantic target: %s", file, target)
			}
		}
	}
}

// TestNoFallbackGlyphPatterns verifies fallbackGetGlyphPatterns is removed from production.
func TestNoFallbackGlyphPatterns(t *testing.T) {
	path := "internal/decipher/glyphs.go"
	content, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("cannot read glyphs.go: %v", err)
	}

	if strings.Contains(string(content), "fallbackGetGlyphPatterns") {
		t.Error("glyphs.go still contains fallbackGetGlyphPatterns function")
	}
}

// TestGlyphDataProducesSignals verifies glyph data from YAML produces signals.
func TestGlyphDataProducesSignals(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Skip("no embedded knowledge available")
	}

	// Test Latin bigram pattern produces signal
	patterns := knowledgeBasedGlyphLookup(ScriptLatin, kb)

	foundInward := false
	foundAction := false
	for _, p := range patterns {
		if p.Pattern == "in" && p.Concept == "inward" {
			foundInward = true
		}
		if p.Pattern == "ly" && p.Concept == "action" {
			foundAction = true
		}
	}

	if !foundInward {
		t.Error("glyph data should include 'in' pattern with concept 'inward'")
	}
	if !foundAction {
		t.Error("glyph data should include 'ly' pattern with concept 'action'")
	}

	// Test Hebrew rune pattern
	hebrewPatterns := knowledgeBasedGlyphLookup(ScriptHebrew, kb)
	foundAleph := false
	for _, p := range hebrewPatterns {
		if p.Rune == 0x05D0 {
			foundAleph = true
			break
		}
	}
	if !foundAleph {
		t.Error("glyph data should include Hebrew Aleph (0x05D0)")
	}
}

// TestVowelStructureUsesGlyphData verifies vowel/consonant ratio uses glyph patterns.
func TestVowelStructureUsesGlyphData(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Skip("no embedded knowledge available")
	}

	// Verify heavy/light patterns exist in glyph data
	patterns := knowledgeBasedGlyphLookup(ScriptLatin, kb)

	foundHeavy := false
	foundLight := false
	for _, p := range patterns {
		if p.Pattern == "heavy" {
			foundHeavy = true
		}
		if p.Pattern == "light" {
			foundLight = true
		}
	}

	if !foundHeavy {
		t.Error("glyph data should include 'heavy' pattern for vowel-heavy structure")
	}
	if !foundLight {
		t.Error("glyph data should include 'light' pattern for consonant-heavy structure")
	}
}

// TestNilKnowledgeReturnsEmptyGlyphSlice verifies nil KB returns empty slice.
func TestNilKnowledgeReturnsEmptyGlyphSlice(t *testing.T) {
	patterns := knowledgeBasedGlyphLookup(ScriptLatin, nil)
	if patterns != nil && len(patterns) > 0 {
		t.Error("glyph lookup with nil KB should return nil or empty slice")
	}
}

// TestSoundChannelUsesNeutralTargets verifies sound channel uses non-semantic targets.
func TestSoundChannelUsesNeutralTargets(t *testing.T) {
	forms := GenerateForms("test")
	result := runSoundChannel(forms)

	// Sound channel should use neutral targets, not semantic ones
	for _, sig := range result.Signals {
		switch sig.Target {
		case "core", "sound", "spirit", "life", "breath", "ground", "inward", "being":
			t.Errorf("sound channel signal has semantic target %q", sig.Target)
		}
	}
}
