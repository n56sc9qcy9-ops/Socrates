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

	// Determine which tokens are standalone function words.
	// Signals from these tokens will be flagged as IsStandaloneToken.
	standaloneTokens := make(map[string]bool)
	for _, token := range forms.Tokens {
		if IsPrepositionOrFunctionWord(token) {
			standaloneTokens[token] = true
		}
	}

	switch script {
	case ScriptLatin:
		latinSignals := analyzeLatinGlyphs(input, kb)
		// Tag signals that come from standalone tokens (single-token input only).
		// For multi-token input, we can't distinguish which token each pattern
		// came from, so we don't mark signals as standalone.
		if standaloneTokens[input] && len(forms.Tokens) == 1 {
			for i := range latinSignals {
				latinSignals[i].IsStandaloneToken = true
			}
		}
		signals = append(signals, latinSignals...)
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

	// Apply weight reduction for signals from standalone tokens
	// to prevent function words from creating overconfident structural fields.
	signals = applyStandaloneWeightReduction(signals)

	score := calculateChannelScore(signals)

	return ChannelResult{
		Name:    "Glyph",
		Signals: signals,
		Score:   score,
	}
}

// analyzeLatinGlyphs analyzes Latin letter patterns.
// Uses knowledge-based glyph lookup when available.
// Orthographic observations (repeated letters, vowel/consonant structure)
// are ONLY emitted if curated glyph data explicitly supports that pattern.
func analyzeLatinGlyphs(s string, kb *knowledge.Knowledge) []Signal {
	signals := make([]Signal, 0)

	if len(s) == 0 {
		return signals
	}

	// Get ALL patterns from knowledge base first
	patterns := knowledgeBasedGlyphLookup(ScriptLatin, kb)
	patternsMap := make(map[string]GlyphPatternSpec)
	for _, p := range patterns {
		patternsMap[p.Pattern] = p
	}

	// Build set of consecutive repeated-letter bigrams actually present in the word
	// Only patterns like "ii", "ee", "oo" that appear consecutively in the word qualify
	// NOT automatic: only emit if glyph data explicitly maps repeated letters to concepts
	presentRepeatedBigrams := make(map[string]bool)
	for i := 0; i < len(s)-1; i++ {
		bigram := s[i : i+2]
		if bigram[0] == bigram[1] {
			// This is a consecutive repeated-letter bigram (like "ii")
			presentRepeatedBigrams[bigram] = true
		}
	}

	// Check for curated repeated-letter patterns ONLY if they actually exist in the word
	for repeatedBigram := range presentRepeatedBigrams {
		if spec, ok := patternsMap[repeatedBigram]; ok {
			signals = append(signals, Signal{
				Text:       "repeated letters detected",
				Target:     spec.Concept,
				Channel:    "Glyph",
				Lens:       "glyph",
				Confidence: spec.Confidence,
				Weight:     spec.Weight,
				IsDirect:   false,
			})
		}
		// If no curated pattern exists for this repeated bigram, emit nothing
	}

	// N-grams (bigrams) - use knowledge base lookup only
	if len(s) >= 2 {
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
				IsDirect:   false,
				})
			}
		}
	}

	// Prefix/suffix analysis - use knowledge base lookup only
	if len(s) >= 3 {
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
				IsDirect:   false,
					})
				}
			}
		}
	}

	if len(s) >= 2 {
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
				IsDirect:   false,
					})
				}
			}
		}
	}

	// Vowel/consonant structure - ONLY emit if knowledge base has explicit structural patterns
	// NOT automatic: orthographic observations are NOT glyph meaning by default
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

	// Only emit structural signals if curated data explicitly maps them
	if consCount > 0 {
		ratio := float64(vowelCount) / float64(consCount)

		// Look for explicit vowel_structure patterns in knowledge base
		for _, p := range patterns {
			if p.Lens == "vowel-structure" || p.Lens == "glyph-structure" {
				if ratio > 0.5 && (p.Pattern == "heavy" || p.Pattern == "vowel-heavy") {
					signals = append(signals, Signal{
						Text:       "vowel-heavy structure",
						Target:     p.Concept,
						Channel:    "Glyph",
						Lens:       "glyph",
						Confidence: p.Confidence,
						Weight:     p.Weight,
				IsDirect:   false,
					})
					break
				}
				if ratio <= 0.5 && (p.Pattern == "light" || p.Pattern == "consonant-heavy") {
					signals = append(signals, Signal{
						Text:       "consonant-heavy structure",
						Target:     p.Concept,
						Channel:    "Glyph",
						Lens:       "glyph",
						Confidence: p.Confidence,
						Weight:     p.Weight,
				IsDirect:   false,
					})
					break
				}
			}
		}
		// If no curated structural patterns exist, emit nothing
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
				IsDirect:   false,
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
				IsDirect:   false,
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
				IsDirect:   false,
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

// isLetter returns true if r is an alphabetic letter (Latin or otherwise).
// Excludes whitespace, punctuation, and other non-letter characters.
func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
		(r >= 0xC0 && r <= 0x24F) || // Latin Extended
		(r >= 0x370 && r <= 0x3FF) || // Greek
		(r >= 0x400 && r <= 0x4FF) || // Cyrillic
		(r >= 0x900 && r <= 0x97F) || // Devanagari
		(r >= 0x4E00 && r <= 0x9FFF) || // CJK Unified Ideographs
		(r >= 0x3040 && r <= 0x309F) || // Hiragana
		(r >= 0x30A0 && r <= 0x30FF) // Katakana
}
