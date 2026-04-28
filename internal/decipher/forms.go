package decipher

import (
	"sort"
	"strings"
)

// GenerateForms creates the normalized forms from input.
func GenerateForms(input string) Forms {
	runes := []rune(input)
	script := DetectScript(input)
	normalized := normalizeForForms(input)
	tokens := tokenize(normalized)

	// Generate phonetic keys
	phoneticKeys := generatePhoneticKeys(normalized)

	// Generate fragment paths
	fragments := generateFragmentPaths(normalized)

	return Forms{
		Normalized:   normalized,
		Script:       script,
		Tokens:       tokens,
		Runes:        runes,
		PhoneticKeys: phoneticKeys,
		Fragments:    fragments,
	}
}

// normalizeForForms lowercases and removes diacritics for Latin scripts.
func normalizeForForms(s string) string {
	script := DetectScript(s)
	if script == ScriptHebrew || script == ScriptDevanagari || script == ScriptHan {
		return s // Keep original for non-Latin scripts
	}
	result := strings.ToLower(s)
	// Remove diacritics for Latin
	var clean strings.Builder
	for _, r := range result {
		if r >= 0x0300 && r <= 0x036F {
			continue // skip combining diacritics
		}
		clean.WriteRune(r)
	}
	return clean.String()
}

// tokenize splits input into tokens (words for Latin, characters for others).
func tokenize(s string) []string {
	script := DetectScript(s)
	if script == ScriptLatin {
		return strings.Fields(s)
	}
	// For non-Latin scripts, treat each character as a token
	tokens := make([]string, 0, len(s))
	for _, r := range s {
		tokens = append(tokens, string(r))
	}
	return tokens
}

// generatePhoneticKeys creates sound-based representations.
func generatePhoneticKeys(s string) []string {
	script := DetectScript(s)
	if script != ScriptLatin {
		return []string{} // Phonetic keys only for Latin script
	}

	keys := make([]string, 0)

	// Normalized form
	keys = append(keys, s)

	// Consonant skeleton
	skeleton := consonantSkeleton(s)
	keys = append(keys, "skeleton:"+skeleton)

	// Vowel skeleton
	vowelSkeleton := vowelSkeleton(s)
	keys = append(keys, "vowels:"+vowelSkeleton)

	// Phonetic variants
	keys = append(keys, phoneticVariants(s)...)

	return keys
}

// consonantSkeleton removes vowels, keeping consonant core.
func consonantSkeleton(s string) string {
	vowels := "aeiouyAEIOUY"
	var skeleton strings.Builder
	for _, r := range s {
		if !strings.ContainsRune(vowels, r) {
			skeleton.WriteRune(r)
		}
	}
	return skeleton.String()
}

// vowelSkeleton keeps only vowel sounds.
func vowelSkeleton(s string) string {
	vowels := "aeiouAEIOU"
	var skeleton strings.Builder
	for _, r := range s {
		if strings.ContainsRune(vowels, r) {
			skeleton.WriteRune(r)
		}
	}
	return skeleton.String()
}

// phoneticVariants applies simple phonetic substitutions.
func phoneticVariants(s string) []string {
	variants := make([]string, 0)

	// Apply substitutions to get phonetic keys
	variant := s

	// ph -> f (includes PH)
	variant = strings.ReplaceAll(strings.ReplaceAll(variant, "Ph", "F"), "ph", "f")

	// ch -> k (and keep as ch)
	variant = strings.ReplaceAll(strings.ReplaceAll(variant, "Ch", "K"), "ch", "k")

	// c -> k or s (prefer k) - includes C
	variant = strings.ReplaceAll(strings.ReplaceAll(variant, "C", "K"), "c", "k")

	// q -> k - includes Q
	variant = strings.ReplaceAll(strings.ReplaceAll(variant, "Q", "K"), "q", "k")

	// y -> i - includes Y
	variant = strings.ReplaceAll(strings.ReplaceAll(variant, "Y", "I"), "y", "i")

	// v -> w - includes V
	variant = strings.ReplaceAll(strings.ReplaceAll(variant, "V", "W"), "v", "w")

	// x -> ks - includes X
	variant = strings.ReplaceAll(strings.ReplaceAll(variant, "X", "KS"), "x", "ks")

	// Remove silent e
	variant = strings.TrimSuffix(variant, "e")

	// Collapse double letters
	variant = collapseDoubles(variant)

	variants = append(variants, "phonetic:"+variant)

	return variants
}

// collapseDoubles removes consecutive duplicate letters.
func collapseDoubles(s string) string {
	if len(s) == 0 {
		return s
	}
	var result strings.Builder
	last := rune(0)
	for _, r := range s {
		if r != last {
			result.WriteRune(r)
			last = r
		}
	}
	return result.String()
}

// generateFragmentPaths generates all possible splits of the input.
// This is the key algorithm - it generates candidates, NOT hardcoded paths.
func generateFragmentPaths(s string) []FragmentPath {
	script := DetectScript(s)
	if script != ScriptLatin {
		return []FragmentPath{} // Fragment generation only for Latin for now
	}

	if len(s) == 0 {
		return []FragmentPath{}
	}

	paths := make([]FragmentPath, 0)

	// Generate all reasonable splits up to maxParts fragments
	maxParts := 5
	minPartLen := 1
	maxPartLen := 6

	// Generate all possible split points
	splitPoints := generateSplitPoints(len(s), maxParts)

	for _, points := range splitPoints {
		parts := extractParts(s, points)
		if isValidSplit(parts, minPartLen, maxPartLen) {
			confidence := calculateSplitConfidence(parts, s)
			paths = append(paths, FragmentPath{
				Parts:      parts,
				Method:     "generated",
				Confidence: confidence,
			})
		}
	}

	// Sort by confidence (descending)
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].Confidence > paths[j].Confidence
	})

	// Limit to top paths
	if len(paths) > 20 {
		paths = paths[:20]
	}

	return paths
}

// generateSplitPoints generates all combinations of split points.
func generateSplitPoints(length, maxParts int) [][]int {
	results := make([][]int, 0)

	// Generate split points for 2 to maxParts parts
	for numParts := 2; numParts <= maxParts && numParts <= length; numParts++ {
		// Use recursive generation for split points
		points := make([]int, 0)
		generateSplitPointsRecursive(length, 0, numParts, points, &results)
	}

	return results
}

// generateSplitPointsRecursive recursively generates split points.
func generateSplitPointsRecursive(length, start, remaining int, current []int, results *[][]int) {
	if remaining == 1 {
		// Last segment goes to end
		*results = append(*results, append([]int{}, current...))
		return
	}

	for pos := start + 1; pos < length && remaining > 1; pos++ {
		current = append(current, pos)
		generateSplitPointsRecursive(length, pos, remaining-1, current, results)
		current = current[:len(current)-1]
	}
}

// extractParts splits string at given points.
func extractParts(s string, points []int) []string {
	parts := make([]string, 0)
	prev := 0
	for _, p := range points {
		parts = append(parts, s[prev:p])
		prev = p
	}
	parts = append(parts, s[prev:])
	return parts
}

// isValidSplit checks if all parts meet size constraints.
func isValidSplit(parts []string, minLen, maxLen int) bool {
	for _, p := range parts {
		if len(p) < minLen || len(p) > maxLen {
			return false
		}
	}
	return true
}

// calculateSplitConfidence scores a split based on fragment length and structure.
// Knowledge-based matching boost is handled at analysis time by the Engine.
func calculateSplitConfidence(parts []string, original string) float64 {
	// Base confidence
	conf := 0.5

	// Bonus for valid-length parts
	matchCount := 0
	for _, part := range parts {
		if len(part) >= 3 && len(part) <= 6 {
			matchCount++
		}
	}

	if len(parts) > 0 {
		matchRatio := float64(matchCount) / float64(len(parts))
		conf += matchRatio * 0.3
	}

	// Penalize very many fragments
	if len(parts) > 4 {
		conf -= 0.1
	}

	// Prefer shorter splits
	if len(parts) == 2 {
		conf += 0.1
	}

	// Clamp to [0, 1]
	if conf < 0 {
		conf = 0
	}
	if conf > 1 {
		conf = 1
	}

	return conf
}
