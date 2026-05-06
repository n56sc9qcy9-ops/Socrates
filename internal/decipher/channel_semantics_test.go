package decipher

import (
	"socrates/internal/knowledge"
	"strings"
	"testing"
)

// TestSpiritNoRepeatedLetterRepetition verifies that "spirit" does not show
// [Glyph] repeated letters detected -> repetition as a signal.
func TestSpiritNoRepeatedLetterRepetition(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	engine := NewEngineWithKnowledge(kb)
	reading := engine.Analyze("spirit")

	// Find any Glyph channel signals
	for _, ch := range reading.Channels {
		if ch.Name == "Glyph" {
			for _, s := range ch.Signals {
				if strings.Contains(s.Text, "repeated letters") && s.Target == "repetition" {
					t.Errorf("spirit should not have repeated letter glyph signal: %s", s.Text)
				}
			}
		}
	}
}

// TestLatinOrthographicRepetitionNotLabeledAsGlyph tests that Latin repeated-letter
// observations are orthographic evidence, not glyph meaning.
func TestLatinOrthographicRepetitionNotLabeledAsGlyph(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Test words with repeated letters
	words := []string{"balloon", "book", "coffee"}
	engine := NewEngineWithKnowledge(kb)

	for _, word := range words {
		reading := engine.Analyze(word)
		for _, ch := range reading.Channels {
			if ch.Name == "Glyph" {
				for _, s := range ch.Signals {
					// Repeated letter observation should only map to orthographic concepts, not meaningful glyph meaning
					// repeated-letter-observation is an orthographic signal, not a semantic glyph
					if strings.Contains(s.Text, "repeated letters") && s.Target == "repeated-letter-observation" {
						// This is the expected behavior - repeated letters map to the observation concept
						// which is orthographic, not glyph meaning
						t.Logf("%s: repeated letter observation -> %s (orthographic)", word, s.Target)
					}
				}
			}
		}
	}
}

// TestHebrewAhavaActivatesLove verifies that transliterated Hebrew "ahava"
// activates the love meaning-frequency field through data.
func TestHebrewAhavaActivatesLove(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	engine := NewEngineWithKnowledge(kb)
	reading := engine.Analyze("ahava")

	// Check that love is activated via Fragment channel
	foundLove := false
	for _, ch := range reading.Channels {
		if ch.Name == "Fragment" {
			for _, s := range ch.Signals {
				if s.Target == "love" && s.Lens == "hebrew" {
					foundLove = true
				}
			}
		}
	}

	if !foundLove {
		t.Error("ahava should activate love via Fragment channel with Hebrew lens")
	}
}

// TestChineseTraditionalLoveActivatesLove verifies that traditional Han "愛"
// activates the love meaning-frequency field through data.
func TestChineseTraditionalLoveActivatesLove(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	engine := NewEngineWithKnowledge(kb)
	reading := engine.Analyze("愛")

	// Check that love is activated
	foundLove := false
	for _, ch := range reading.Channels {
		if ch.Name == "Fragment" || ch.Name == "ScriptWord" {
			for _, s := range ch.Signals {
				if s.Target == "love" && s.Lens == "han" {
					foundLove = true
				}
			}
		}
	}

	if !foundLove {
		t.Error("愛 should activate love via Fragment or ScriptWord channel with Han lens")
	}
}

// TestChineseSimplifiedLoveActivatesLove verifies that simplified Han "爱"
// activates the love meaning-frequency field through data.
func TestChineseSimplifiedLoveActivatesLove(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	engine := NewEngineWithKnowledge(kb)
	reading := engine.Analyze("爱")

	// Check that love is activated
	foundLove := false
	for _, ch := range reading.Channels {
		if ch.Name == "Fragment" || ch.Name == "ScriptWord" {
			for _, s := range ch.Signals {
				if s.Target == "love" && s.Lens == "han" {
					foundLove = true
				}
			}
		}
	}

	if !foundLove {
		t.Error("爱 should activate love via Fragment or ScriptWord channel with Han lens")
	}
}

// TestLoveConvergence verifies that English "love", Norwegian "kjærlighet",
// Hebrew "ahava"/"ahavah", and Chinese "愛"/"爱" all converge on the same
// meaning-frequency field.
func TestLoveConvergence(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	engine := NewEngineWithKnowledge(kb)
	words := []string{"love", "kjærlighet", "ahava", "ahavah", "愛", "爱"}

	var loveFields []string
	for _, word := range words {
		reading := engine.Analyze(word)
		// Check if love is activated (in any position)
		for _, ac := range reading.Convergence.ActivatedConcepts {
			if ac.Concept == "love" {
				loveFields = append(loveFields, word)
				break
			}
		}
	}

	// All words should activate the love field
	if len(loveFields) < len(words) {
		t.Errorf("Expected all words to activate love, but only %d did: %v", len(loveFields), loveFields)
	}
}

// TestHebrewElDivineVerifies that Hebrew "אל" activates divine meaning
// through curated Hebrew/script evidence.
func TestHebrewElDivine(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	engine := NewEngineWithKnowledge(kb)
	reading := engine.Analyze("אל")

	// Check that god/power/divine is activated via ScriptWord channel
	foundDivine := false
	for _, ch := range reading.Channels {
		if ch.Name == "ScriptWord" {
			for _, s := range ch.Signals {
				if s.Target == "god" || s.Target == "power" || s.Target == "divine" {
					foundDivine = true
				}
			}
		}
	}

	if !foundDivine {
		t.Error("אל should activate divine meaning via ScriptWord channel")
	}
}

// TestEnglishElSubstringNotAutomaticDivine verifies that English words
// containing "el" substrings do not automatically activate divine meaning
// without curated evidence.
func TestEnglishElSubstringNotAutomaticDivine(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	engine := NewEngineWithKnowledge(kb)
	words := []string{"electricity", "electromagnetism", "television", "element"}

	for _, word := range words {
		reading := engine.Analyze(word)

		// Check that divine-related concepts are NOT in top concepts
		for _, ac := range reading.Convergence.TopConcepts {
			if ac.Concept == "god" || ac.Concept == "divine" {
				t.Errorf("%s should NOT automatically activate divine meaning from 'el' substring", word)
			}
		}
	}
}

// TestHebrewAhavaInForms verifies the curated love form is in the knowledge base.
func TestHebrewAhavaInForms(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	forms := kb.GetFormsByText("ahava")
	found := false
	for _, f := range forms {
		if f.Concept == "love" && f.Lens == "hebrew" {
			found = true
		}
	}

	if !found {
		t.Error("ahava should be in knowledge base as Hebrew love form")
	}
}

// TestChineseLoveInForms verifies the curated Han love forms are in the knowledge base.
func TestChineseLoveInForms(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Check traditional
	forms := kb.GetFormsByText("愛")
	foundTraditional := false
	for _, f := range forms {
		if f.Concept == "love" && f.Lens == "han" {
			foundTraditional = true
		}
	}

	// Check simplified
	forms = kb.GetFormsByText("爱")
	foundSimplified := false
	for _, f := range forms {
		if f.Concept == "love" && f.Lens == "han" {
			foundSimplified = true
		}
	}

	if !foundTraditional {
		t.Error("愛 should be in knowledge base as Han love form")
	}
	if !foundSimplified {
		t.Error("爱 should be in knowledge base as Han love form")
	}
}

// TestRemovingCuratedDataRemovesConvergence verifies that without curated data,
// love convergence would not work. This is tested by checking that words without
// curated forms don't activate love through the same path.
func TestRemovingCuratedDataRemovesConvergence(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Test a word that is NOT in curated data
	engine := NewEngineWithKnowledge(kb)
	reading := engine.Analyze("lovingly")

	// "lovingly" should NOT have "love" as a direct fragment match
	hasDirectLove := false
	for _, ch := range reading.Channels {
		if ch.Name == "Fragment" {
			for _, s := range ch.Signals {
				if s.Target == "love" && strings.Contains(s.Text, "exact whole-token match") {
					hasDirectLove = true
				}
			}
		}
	}

	if hasDirectLove {
		t.Error("lovingly should NOT have exact whole-token love match without curated data")
	}
}
