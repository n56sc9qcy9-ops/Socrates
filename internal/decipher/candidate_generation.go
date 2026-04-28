package decipher

import (
	"strings"
)

// CandidateForm represents a generated candidate form from the input.
type CandidateForm struct {
	Form       string
	Method     string
	Distance   float64
	Confidence string
}

// =============================================================================
// Phase 3: Generic Form Generation
// =============================================================================

// GenerateCandidateForms creates all candidate forms from the input.
// This is the core neighbor generation algorithm - not hardcoded.
// Deduplicates forms by keeping the best match for each unique form.
func GenerateCandidateForms(input string) []CandidateForm {
	candidates := make([]CandidateForm, 0)
	seenForms := make(map[string]CandidateForm) // Track best candidate per form
	script := DetectScript(input)

	// Only process Latin script for neighbor generation
	if script != ScriptLatin {
		return candidates
	}

	normalized := strings.ToLower(strings.TrimSpace(input))

	// Phase 3 form generation methods:
	// 1. Normalized form
	addCandidate(&candidates, seenForms, CandidateForm{
		Form:       normalized,
		Method:     "normalized",
		Distance:   0,
		Confidence: "verified",
	})

	// 2. Consonant skeleton
	skeleton := consonantSkeleton(normalized)
	addCandidate(&candidates, seenForms, CandidateForm{
		Form:       skeleton,
		Method:     "consonant_skeleton",
		Distance:   0,
		Confidence: "verified",
	})

	// 3. Vowel skeleton
	vowelSkeleton := vowelSkeleton(normalized)
	addCandidate(&candidates, seenForms, CandidateForm{
		Form:       vowelSkeleton,
		Method:     "vowel_skeleton",
		Distance:   0,
		Confidence: "plausible",
	})

	// 4. Phonetic variants
	phonetic := generatePhonetic(normalized)
	addCandidate(&candidates, seenForms, CandidateForm{
		Form:       phonetic,
		Method:     "phonetic",
		Distance:   0,
		Confidence: "plausible",
	})

	// 5. De-doubled variant
	dedoubled := collapseDoubles(normalized)
	if dedoubled != normalized {
		addCandidate(&candidates, seenForms, CandidateForm{
			Form:       dedoubled,
			Method:     "dedoubled",
			Distance:   0,
			Confidence: "plausible",
		})
	}

	// 6. Doubled-letter variants (insert doubled consonants)
	for _, c := range generateDoubledVariants(normalized) {
		addCandidate(&candidates, seenForms, c)
	}

	// 7. Common sound substitutions
	for _, c := range generateSoundSubstitutions(normalized) {
		addCandidate(&candidates, seenForms, c)
	}

	// 8. One-letter edit variants
	for _, c := range generateEditVariants(normalized, 1) {
		addCandidate(&candidates, seenForms, c)
	}

	// Phase 3: N-grams (bigrams and trigrams)
	for _, c := range generateNgrams(normalized) {
		addCandidate(&candidates, seenForms, c)
	}

	// Phase 3: Prefix fragments
	for _, c := range generatePrefixFragments(normalized) {
		addCandidate(&candidates, seenForms, c)
	}

	// Phase 3: Suffix fragments
	for _, c := range generateSuffixFragments(normalized) {
		addCandidate(&candidates, seenForms, c)
	}

	return candidates
}

// addCandidate adds a candidate, keeping the best (lowest distance/highest confidence) for duplicates.
func addCandidate(candidates *[]CandidateForm, seenForms map[string]CandidateForm, c CandidateForm) {
	existing, exists := seenForms[c.Form]
	if !exists {
		// New form, add it
		seenForms[c.Form] = c
		*candidates = append(*candidates, c)
	} else if c.Distance < existing.Distance {
		// Better distance found, replace
		seenForms[c.Form] = c
		// Update in candidates slice
		for i := range *candidates {
			if (*candidates)[i].Form == c.Form {
				(*candidates)[i] = c
				break
			}
		}
	}
}

// generatePhonetic creates a phonetic representation.
func generatePhonetic(s string) string {
	result := s

	// Common substitutions
	result = strings.ReplaceAll(result, "ph", "f")
	result = strings.ReplaceAll(result, "ch", "k")
	result = strings.ReplaceAll(result, "c", "k")
	result = strings.ReplaceAll(result, "q", "k")
	result = strings.ReplaceAll(result, "y", "i")
	result = strings.ReplaceAll(result, "v", "w")
	result = strings.ReplaceAll(result, "x", "ks")

	// Remove silent e
	if strings.HasSuffix(result, "e") && len(result) > 2 {
		result = result[:len(result)-1]
	}

	// Collapse doubles
	result = collapseDoubles(result)

	return result
}

// generateDoubledVariants creates variants with doubled consonants.
func generateDoubledVariants(s string) []CandidateForm {
	variants := make([]CandidateForm, 0)
	vowels := "aeiouy"

	// Insert doubled consonant after vowels
	for i, r := range s {
		if strings.ContainsRune(vowels, r) && i+1 < len(s) {
			c := string(s[i+1])
			variant := s[:i+1] + c + c + s[i+2:]
			variants = append(variants, CandidateForm{
				Form:       variant,
				Method:     "doubled_consonant",
				Distance:   1,
				Confidence: "speculative",
			})
		}
	}

	return variants
}

// generateSoundSubstitutions creates variants using common sound substitutions.
func generateSoundSubstitutions(s string) []CandidateForm {
	variants := make([]CandidateForm, 0)

	// k/c/q substitution
	if strings.ContainsAny(s, "ckq") {
		for _, char := range []string{"k", "c", "q"} {
			variant := s
			for _, old := range []string{"c", "k", "q"} {
				if old != char {
					variant = strings.ReplaceAll(variant, old, char)
				}
			}
			if variant != s {
				variants = append(variants, CandidateForm{
					Form:       variant,
					Method:     "sound_kckq",
					Distance:   0.5,
					Confidence: "plausible",
				})
			}
		}
	}

	// s/sh substitution
	if strings.Contains(s, "sh") {
		variants = append(variants, CandidateForm{
			Form:       strings.ReplaceAll(s, "sh", "s"),
			Method:     "sound_sh_s",
			Distance:   0.5,
			Confidence: "plausible",
		})
	}

	// g/j/y substitution
	if strings.ContainsAny(s, "gjy") {
		variants = append(variants, CandidateForm{
			Form:       strings.ReplaceAll(strings.ReplaceAll(s, "g", "j"), "y", "j"),
			Method:     "sound_gjy",
			Distance:   0.5,
			Confidence: "speculative",
		})
	}

	// v/w substitution
	if strings.ContainsAny(s, "vw") {
		replacer := "w"
		if strings.Contains(s, "w") {
			replacer = "v"
		}
		variant := strings.ReplaceAll(s, "v", replacer)
		variant = strings.ReplaceAll(variant, "w", replacer)
		if variant != s {
			variants = append(variants, CandidateForm{
				Form:       variant,
				Method:     "sound_vw",
				Distance:   0.5,
				Confidence: "plausible",
			})
		}
	}

	// i/y substitution
	if strings.ContainsAny(s, "iy") && len(s) > 1 {
		variants = append(variants, CandidateForm{
			Form:       strings.ReplaceAll(s, "y", "i"),
			Method:     "sound_yi",
			Distance:   0.5,
			Confidence: "plausible",
		})
		variants = append(variants, CandidateForm{
			Form:       strings.ReplaceAll(s, "i", "y"),
			Method:     "sound_iy",
			Distance:   0.5,
			Confidence: "plausible",
		})
	}

	return variants
}

// generateEditVariants creates variants with one-letter edits.
func generateEditVariants(s string, editDistance int) []CandidateForm {
	variants := make([]CandidateForm, 0)
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	// Deletion variants
	for i := 0; i < len(s); i++ {
		variant := s[:i] + s[i+1:]
		if len(variant) > 0 {
			variants = append(variants, CandidateForm{
				Form:       variant,
				Method:     "deletion",
				Distance:   1,
				Confidence: "speculative",
			})
		}
	}

	// Insertion variants
	for i := 0; i <= len(s); i++ {
		for _, c := range alphabet {
			variant := s[:i] + string(c) + s[i:]
			variants = append(variants, CandidateForm{
				Form:       variant,
				Method:     "insertion",
				Distance:   1,
				Confidence: "speculative",
			})
		}
	}

	// Substitution variants
	for i := 0; i < len(s); i++ {
		sRune := rune(s[i])
		for _, c := range alphabet {
			if sRune != c {
				variant := s[:i] + string(c) + s[i+1:]
				variants = append(variants, CandidateForm{
					Form:       variant,
					Method:     "substitution",
					Distance:   1,
					Confidence: "speculative",
				})
			}
		}
	}

	return variants
}

// generateNgrams creates bigrams and trigrams from the input.
func generateNgrams(s string) []CandidateForm {
	variants := make([]CandidateForm, 0)

	if len(s) < 2 {
		return variants
	}

	// Bigrams
	for i := 0; i < len(s)-1; i++ {
		variants = append(variants, CandidateForm{
			Form:       s[i : i+2],
			Method:     "bigram",
			Distance:   0.5,
			Confidence: "plausible",
		})
	}

	// Trigrams
	for i := 0; i < len(s)-2; i++ {
		variants = append(variants, CandidateForm{
			Form:       s[i : i+3],
			Method:     "trigram",
			Distance:   0.4,
			Confidence: "plausible",
		})
	}

	return variants
}

// generatePrefixFragments creates prefix fragments of varying lengths.
func generatePrefixFragments(s string) []CandidateForm {
	variants := make([]CandidateForm, 0)

	// Generate prefixes of lengths 2 to min(5, len(s)-1)
	maxLen := len(s) - 1
	if maxLen > 5 {
		maxLen = 5
	}
	if maxLen < 2 {
		return variants
	}

	for l := 2; l <= maxLen; l++ {
		prefix := s[:l]
		variants = append(variants, CandidateForm{
			Form:       prefix,
			Method:     "prefix_" + string(rune('0'+l)),
			Distance:   0.3,
			Confidence: "plausible",
		})
	}

	return variants
}

// generateSuffixFragments creates suffix fragments of varying lengths.
func generateSuffixFragments(s string) []CandidateForm {
	variants := make([]CandidateForm, 0)

	// Generate suffixes of lengths 2 to min(5, len(s)-1)
	maxLen := len(s) - 1
	if maxLen > 5 {
		maxLen = 5
	}
	if maxLen < 2 {
		return variants
	}

	for l := 2; l <= maxLen; l++ {
		suffix := s[len(s)-l:]
		variants = append(variants, CandidateForm{
			Form:       suffix,
			Method:     "suffix_" + string(rune('0'+l)),
			Distance:   0.3,
			Confidence: "plausible",
		})
	}

	return variants
}
