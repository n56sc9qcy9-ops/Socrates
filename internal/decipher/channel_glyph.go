package decipher

import (
	"strings"

	"socrates/internal/knowledge"
)

// runGlyphChannel analyzes letter/shape patterns.
func runGlyphChannel(forms Forms, kb *knowledge.Knowledge) ChannelResult {
	signals := make([]Signal, 0)

	input := forms.Normalized
	script := forms.Script

	switch script {
	case ScriptLatin:
		signals = append(signals, analyzeLatinGlyphs(input, kb)...)
	case ScriptHebrew:
		signals = append(signals, analyzeHebrewGlyphs(forms.Runes, kb)...)
	case ScriptDevanagari:
		signals = append(signals, analyzeDevanagariGlyphs(forms.Runes, kb)...)
	case ScriptHan:
		signals = append(signals, analyzeHanGlyphs(forms.Runes, kb)...)
	default:
		signals = append(signals, analyzeLatinGlyphs(input, kb)...)
	}

	// Tag all signals with the Glyph channel
	for i := range signals {
		signals[i].Channel = "Glyph"
	}

	score := calculateChannelScore(signals)

	return ChannelResult{
		Name:    "Glyph",
		Signals: signals,
		Score:   score,
	}
}

// analyzeLatinGlyphs analyzes Latin letter patterns.
// Uses knowledge-based glyph lookup when available.
func analyzeLatinGlyphs(s string, kb *knowledge.Knowledge) []Signal {
	signals := make([]Signal, 0)

	if len(s) == 0 {
		return signals
	}

	// Character frequency
	charFreq := make(map[rune]int)
	for _, r := range s {
		charFreq[r]++
	}

	// Repeated letters
	for r, count := range charFreq {
		if count > 1 {
			signals = append(signals, Signal{
				Text:       string(r) + " repeated " + itoa(count) + "x",
				Target:     "repetition",
				Channel:    "Glyph",
				Lens:       "glyph",
				Confidence: ConfidencePlausible,
				Weight:     0.3,
			})
		}
	}

	// N-grams (bigrams and trigrams) - use knowledge base lookup
	if len(s) >= 2 {
		patterns := knowledgeBasedGlyphLookup(ScriptLatin, kb)
		patternsMap := make(map[string]GlyphPatternSpec)
		for _, p := range patterns {
			patternsMap[p.Pattern] = p
		}

		for i := 0; i < len(s)-1; i++ {
			bigram := s[i : i+2]
			if p, ok := patternsMap[bigram]; ok {
				signals = append(signals, Signal{
					Text:       "bigram pattern: " + bigram,
					Target:     p.Concept,
					Channel:    "Glyph",
					Lens:       "glyph",
					Confidence: p.Confidence,
					Weight:     p.Weight,
				})
			}
		}
	}

	// Prefix/suffix analysis - use knowledge base lookup
	if len(s) >= 3 {
		patterns := knowledgeBasedGlyphLookup(ScriptLatin, kb)
		for _, p := range patterns {
			if len(p.Pattern) >= 3 && len(p.Pattern) <= len(s) {
				prefix := s[:len(p.Pattern)]
				if prefix == p.Pattern {
					signals = append(signals, Signal{
						Text:       "prefix: " + prefix,
						Target:     p.Concept,
						Channel:    "Glyph",
						Lens:       "glyph",
						Confidence: p.Confidence,
						Weight:     p.Weight,
					})
				}
			}
		}
	}

	if len(s) >= 2 {
		patterns := knowledgeBasedGlyphLookup(ScriptLatin, kb)
		for _, p := range patterns {
			if len(p.Pattern) == 2 {
				suffix := s[len(s)-2:]
				if suffix == p.Pattern {
					signals = append(signals, Signal{
						Text:       "suffix: " + suffix,
						Target:     p.Concept,
						Channel:    "Glyph",
						Lens:       "glyph",
						Confidence: p.Confidence,
						Weight:     p.Weight,
					})
				}
			}
		}
	}

	// Vowel/consonant structure - use knowledge base patterns
	vowels := "aeiouAEIOU"
	vowelCount := 0
	consCount := 0
	for _, r := range s {
		if strings.ContainsRune(vowels, r) {
			vowelCount++
		} else if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
			consCount++
		}
	}
	if consCount > 0 {
		ratio := float64(vowelCount) / float64(consCount)
		patterns := knowledgeBasedGlyphLookup(ScriptLatin, kb)
		patternsMap := make(map[string]GlyphPatternSpec)
		for _, p := range patterns {
			patternsMap[p.Pattern] = p
		}
		if ratio > 0.5 {
			if p, ok := patternsMap["heavy"]; ok {
				signals = append(signals, Signal{
					Text:       "vowel-heavy structure",
					Target:     p.Concept,
					Channel:    "Glyph",
					Lens:       "glyph",
					Confidence: p.Confidence,
					Weight:     p.Weight,
				})
			}
		} else {
			if p, ok := patternsMap["light"]; ok {
				signals = append(signals, Signal{
					Text:       "consonant-heavy structure",
					Target:     p.Concept,
					Channel:    "Glyph",
					Lens:       "glyph",
					Confidence: p.Confidence,
					Weight:     p.Weight,
				})
			}
		}
	}

	return signals
}

// analyzeHebrewGlyphs analyzes Hebrew letter patterns.
// Uses knowledge base when available for rune-to-concept mapping.
// Returns letter-shape associations, NOT automatic breath/spirit for unknown words.
func analyzeHebrewGlyphs(runes []rune, kb *knowledge.Knowledge) []Signal {
	signals := make([]Signal, 0)

	if len(runes) == 0 {
		return signals
	}

	// Get patterns from knowledge base
	patterns := knowledgeBasedGlyphLookup(ScriptHebrew, kb)
	runeToPattern := make(map[uint32]GlyphPatternSpec)
	for _, p := range patterns {
		runeToPattern[p.Rune] = p
	}

	for _, r := range runes {
		if p, ok := runeToPattern[uint32(r)]; ok {
			signals = append(signals, Signal{
				Text:       "letter: " + p.Pattern,
				Target:     p.Concept,
				Channel:    "Glyph",
				Lens:       "hebrew-glyph",
				Confidence: p.Confidence,
				Weight:     p.Weight,
			})
		} else {
			// Unknown Hebrew letter
			signals = append(signals, Signal{
				Text:       "Hebrew letter",
				Target:     "unknown",
				Channel:    "Glyph",
				Lens:       "hebrew-glyph",
				Confidence: ConfidenceSpeculative,
				Weight:     0.1,
			})
		}
	}

	return signals
}

// analyzeDevanagariGlyphs analyzes Devanagari letter patterns.
// Uses knowledge base when available.
// Returns script/glyph signals, NOT automatic spiritual meaning.
func analyzeDevanagariGlyphs(runes []rune, kb *knowledge.Knowledge) []Signal {
	signals := make([]Signal, 0)

	// Get patterns from knowledge base
	patterns := knowledgeBasedGlyphLookup(ScriptDevanagari, kb)
	runeToPattern := make(map[uint32]GlyphPatternSpec)
	for _, p := range patterns {
		runeToPattern[p.Rune] = p
	}

	foundSpecific := false
	for _, r := range runes {
		if p, ok := runeToPattern[uint32(r)]; ok {
			foundSpecific = true
			signals = append(signals, Signal{
				Text:       "glyph: " + p.Pattern,
				Target:     p.Concept,
				Channel:    "Glyph",
				Lens:       "devanagari-glyph",
				Confidence: p.Confidence,
				Weight:     p.Weight,
			})
		}
	}

	// If no specific patterns found, mark as Devanagari script
	if !foundSpecific {
		signals = append(signals, Signal{
			Text:       "Devanagari script",
			Target:     "unknown",
			Channel:    "Glyph",
			Lens:       "devanagari-glyph",
			Confidence: ConfidenceSpeculative,
			Weight:     0.2,
		})
	}

	return signals
}

// analyzeHanGlyphs analyzes Chinese character patterns.
// Uses knowledge base when available for rune-to-meaning mapping.
func analyzeHanGlyphs(runes []rune, kb *knowledge.Knowledge) []Signal {
	signals := make([]Signal, 0)

	// Get patterns from knowledge base
	patterns := knowledgeBasedGlyphLookup(ScriptHan, kb)
	runeToPattern := make(map[uint32]GlyphPatternSpec)
	for _, p := range patterns {
		runeToPattern[p.Rune] = p
	}

	for _, r := range runes {
		if p, ok := runeToPattern[uint32(r)]; ok {
			signals = append(signals, Signal{
				Text:       "Han character",
				Target:     p.Concept,
				Channel:    "Glyph",
				Lens:       "han-glyph",
				Confidence: p.Confidence,
				Weight:     p.Weight,
			})
		} else {
			// Unknown Han character
			signals = append(signals, Signal{
				Text:       "Han character",
				Target:     "unknown",
				Channel:    "Glyph",
				Lens:       "han-glyph",
				Confidence: ConfidenceSpeculative,
				Weight:     0.2,
			})
		}
	}

	return signals
}
