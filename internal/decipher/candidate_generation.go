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
// Bounded: speculative candidates (edit variants, n-grams, prefixes, suffixes)
// are capped for long inputs to prevent unbounded work.
// Returns candidates and total count of discarded candidates due to bounds.
func GenerateCandidateForms(input string, bounds ...CandidateBounds) ([]CandidateForm, int) {
	candidates := make([]CandidateForm, 0)
	seenForms := make(map[string]CandidateForm) // Track best candidate per form
	script := DetectScript(input)
	discarded := 0

	// Only process Latin script for neighbor generation
	if script != ScriptLatin {
		return candidates, 0
	}

	// Determine bounds with safe defaults
	maxCandidates := 50
	maxSpeculative := 30
	if len(bounds) > 0 && bounds[0].MaxCandidates > 0 {
		maxCandidates = bounds[0].MaxCandidates
	}
	if len(bounds) > 0 && bounds[0].MaxSpeculativeCandidates > 0 {
		maxSpeculative = bounds[0].MaxSpeculativeCandidates
	} else if len(bounds) > 0 && bounds[0].MaxCandidates > 0 {
		// If only MaxCandidates is set, use it as the speculative cap
		maxSpeculative = bounds[0].MaxCandidates
	}

	normalized := strings.ToLower(strings.TrimSpace(input))
	inputLen := len(normalized)

	// Phase 3 form generation methods:
	// CORE: High-confidence candidates always included first

	// 1. Normalized form (always first, highest confidence)
	addCandidate(&candidates, seenForms, CandidateForm{
		Form:       normalized,
		Method:     "normalized",
		Distance:   0,
		Confidence: "verified",
	})

	// 2. Consonant skeleton (always included)
	skeleton := consonantSkeleton(normalized)
	addCandidate(&candidates, seenForms, CandidateForm{
		Form:       skeleton,
		Method:     "consonant_skeleton",
		Distance:   0,
		Confidence: "verified",
	})

	// 3. Vowel skeleton (included if not same as normalized or skeleton)
	vowelSkeleton := vowelSkeleton(normalized)
	addCandidate(&candidates, seenForms, CandidateForm{
		Form:       vowelSkeleton,
		Method:     "vowel_skeleton",
		Distance:   0,
		Confidence: "plausible",
	})

	// 4. Phonetic variants (always included - core plausible form)
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

	// Count speculative candidates generated so far
	speculativeCount := 0

	// 6. Doubled-letter variants (insert doubled consonants)
	// These are speculative - they grow with input length
	doubledVariants := generateDoubledVariants(normalized)
	maxDoubled := 3
	for i, c := range doubledVariants {
		if speculativeCount >= maxSpeculative {
			discarded += len(doubledVariants) - i
			break
		}
		if i >= maxDoubled {
			discarded += len(doubledVariants) - i
			break
		}
		addCandidate(&candidates, seenForms, c)
		speculativeCount++
	}

	// 7. Common sound substitutions
	// These are plausible but can multiply; cap them
	soundSubs := generateSoundSubstitutions(normalized)
	maxSoundSubs := 4
	for i, c := range soundSubs {
		if speculativeCount >= maxSpeculative {
			discarded += len(soundSubs) - i
			break
		}
		if i >= maxSoundSubs {
			discarded += len(soundSubs) - i
			break
		}
		addCandidate(&candidates, seenForms, c)
		speculativeCount++
	}

	// SPECULATIVE: Bounded variants for long inputs
	// For short inputs (< 8 chars), generate all edit variants
	// For long inputs (>= 8 chars), cap edit variants based on speculative budget

	// 8. One-letter edit variants (bounded for long inputs)
	// These are the most expensive speculative source
	if inputLen < 8 {
		// Short inputs: generate all edit variants but cap at speculative budget
		for _, c := range generateEditVariants(normalized, 1) {
			if speculativeCount >= maxSpeculative {
				discarded++
				continue
			}
			addCandidate(&candidates, seenForms, c)
			speculativeCount++
		}
	} else {
		// Long inputs: cap edit variants strictly
		// Calculate what we would generate vs what we can afford
		totalDeletions := len(normalized)
		totalInsertions := len(normalized) * 26
		totalSubstitutions := len(normalized) * 25

		// Budget remaining for edit variants
		editBudget := maxSpeculative - speculativeCount
		if editBudget < 0 {
			editBudget = 0
		}

		// Cap deletions at editBudget (but also limit per-type)
		maxDeletions := min2(editBudget/2, inputLen-1)
		if maxDeletions > 10 {
			maxDeletions = 10
		}
		for i := 0; i < len(normalized) && i < maxDeletions; i++ {
			variant := normalized[:i] + normalized[i+1:]
			if len(variant) > 0 {
				addCandidate(&candidates, seenForms, CandidateForm{
					Form:       variant,
					Method:     "deletion",
					Distance:   1,
					Confidence: "speculative",
				})
				speculativeCount++
			}
		}

		// Track discarded edit variants
		discarded += totalDeletions - min2(maxDeletions, totalDeletions)
		discarded += totalInsertions // All insertions discarded for long inputs
		discarded += totalSubstitutions // All substitutions discarded for long inputs
	}

	// 9. N-grams (bounded for long inputs)
	// Cap n-grams: bigrams and trigrams grow quadratically with input length
	maxNgrams := 0
	if inputLen < 10 {
		// Short inputs: generate all but cap total n-grams
		maxNgrams = 10
	} else {
		// For long inputs, cap n-grams strictly
		maxNgrams = 6
	}

	// Calculate n-grams we would generate
	maxPossibleNgrams := 0
	if inputLen >= 2 {
		maxPossibleNgrams += inputLen - 1 // bigrams
	}
	if inputLen >= 3 {
		maxPossibleNgrams += inputLen - 2 // trigrams
	}

	count := 0
	// Only generate bigrams for bounded inputs
	for i := 0; i < len(normalized)-1 && count < maxNgrams; i++ {
		addCandidate(&candidates, seenForms, CandidateForm{
			Form:       normalized[i:i+2],
			Method:     "bigram",
			Distance:   0.5,
			Confidence: "plausible",
		})
		count++
	}
	// Count discarded n-grams
	if maxPossibleNgrams > count {
		discarded += maxPossibleNgrams - count
	}

	// 10. Prefix fragments (bounded for long inputs)
	maxPrefixLen := 4
	if inputLen > 12 {
		maxPrefixLen = 3
	}
	// Calculate prefixes we would generate
	originalPrefixes := max2(0, inputLen-2)
	cappedPrefixes := min2(originalPrefixes, maxPrefixLen-1)
	for l := 2; l <= maxPrefixLen && l <= inputLen-1 && (l-1) < maxPrefixLen; l++ {
		addCandidate(&candidates, seenForms, CandidateForm{
			Form:       normalized[:l],
			Method:     "prefix_" + string(rune('0'+l)),
			Distance:   0.3,
			Confidence: "plausible",
		})
	}
	// Track discarded prefixes
	if originalPrefixes > cappedPrefixes {
		discarded += originalPrefixes - cappedPrefixes
	}

	// 11. Suffix fragments (bounded for long inputs)
	maxSuffixLen := 4
	if inputLen > 12 {
		maxSuffixLen = 3
	}
	// Calculate suffixes we would generate
	originalSuffixes := max2(0, inputLen-2)
	cappedSuffixes := min2(originalSuffixes, maxSuffixLen-1)
	for l := 2; l <= maxSuffixLen && l <= inputLen-1 && (l-1) < maxSuffixLen; l++ {
		addCandidate(&candidates, seenForms, CandidateForm{
			Form:       normalized[len(normalized)-l:],
			Method:     "suffix_" + string(rune('0'+l)),
			Distance:   0.3,
			Confidence: "plausible",
		})
	}
	// Track discarded suffixes
	if originalSuffixes > cappedSuffixes {
		discarded += originalSuffixes - cappedSuffixes
	}

	// Apply absolute cap if we somehow exceeded maxCandidates
	if len(candidates) > maxCandidates {
		discarded += len(candidates) - maxCandidates
		candidates = candidates[:maxCandidates]
	}

	return candidates, discarded
}

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max2(a, b int) int {
	if a > b {
		return a
	}
	return b
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
