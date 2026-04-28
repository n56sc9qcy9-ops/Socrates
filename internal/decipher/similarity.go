package decipher

import (
	"math"
	"sort"
	"strings"

	"socrates/internal/knowledge"
)

// MatchEvidence represents evidence from fuzzy matching against anchors.
type MatchEvidence struct {
	InputForm  string
	AnchorForm string
	Method     string
	Distance   float64
	Weight     float64
}

// AnchorConcept is a minimal anchor concept for fuzzy matching.
type AnchorConcept struct {
	Form       string
	Concept    string
	Confidence string
	Weight     float64
}

// =============================================================================
// Phase B: Fuzzy Anchor Matching
// =============================================================================

// FuzzyMatchEvidence finds fuzzy matches between candidate forms and anchors.
// Uses method-aware acceptance logic to prevent weak matches.
// Key tightening: accepts only forms >= 4 chars for most methods (fragments need >= 3).
func FuzzyMatchEvidence(candidates []CandidateForm, anchors []AnchorConcept) []MatchEvidence {
	evidence := make([]MatchEvidence, 0)
	seen := make(map[string]bool) // Dedupe by inputForm:anchorForm

	for _, cand := range candidates {
		// CRITICAL: Skip very short candidates - they cause spurious matches
		// Only accept candidates >= 4 chars for fuzzy matching (fragments >= 3 are too short)
		if len(cand.Form) < 4 {
			continue
		}

		for _, anchor := range anchors {
			// Skip anchors that are too short (need at least 3 chars for meaningful match)
			if len(anchor.Form) < 3 {
				continue
			}

			method, distance := fuzzyMatch(cand.Form, anchor.Form)
			if acceptMatch(method, distance, cand.Form, anchor.Form) {
				// Deduplicate
				key := cand.Form + ":" + anchor.Form
				if seen[key] {
					continue
				}
				seen[key] = true

				weight := calculateMatchWeight(cand, anchor, method, distance)
				evidence = append(evidence, MatchEvidence{
					InputForm:  cand.Form,
					AnchorForm: anchor.Form,
					Method:     method,
					Distance:   distance,
					Weight:     weight,
				})
			}
		}
	}

	// Sort by weight descending, with exact methods ranking above fuzzy
	sort.Slice(evidence, func(i, j int) bool {
		return compareEvidence(evidence[i], evidence[j]) > 0
	})

	return evidence
}

// acceptMatch determines if a fuzzy match should be accepted based on method and distance.
// Method-aware rules:
// - exact match: accept (no length restriction since exact match is definitive)
// - case-insensitive exact match: accept
// - consonant skeleton: accept only if skeleton length >= 3 AND forms >= 4 chars
// - vowel skeleton: accept only if vowel skeleton length >= 3 and original forms >= 5 chars
// - phonetic: accept only if phonetic key length >= 4 AND forms >= 5 chars
// - Levenshtein: accept only if normalized distance <= 0.15 AND both forms >= 5 chars
func acceptMatch(method string, distance float64, inputForm, anchorForm string) bool {
	switch method {
	case "exact", "case_insensitive":
		return true
	case "consonant_skeleton":
		// Accept only if skeleton is at least 3 chars AND original forms >= 4 chars
		skelInput := consonantSkeleton(inputForm)
		return len(skelInput) >= 3 && len(inputForm) >= 4 && len(anchorForm) >= 4
	case "vowel_skeleton":
		// Accept only if vowel skeletons are at least 3 chars and originals >= 5
		vowInput := vowelSkeleton(inputForm)
		return len(vowInput) >= 3 && len(inputForm) >= 5 && len(anchorForm) >= 5
	case "phonetic":
		// Accept only if phonetic keys are at least 4 chars AND forms >= 5 chars
		phonInput := generatePhonetic(inputForm)
		return len(phonInput) >= 4 && len(inputForm) >= 5 && len(anchorForm) >= 5
	case "levenshtein":
		// Accept only if normalized distance is <= 0.15 AND both forms >= 5 chars
		if len(inputForm) < 5 || len(anchorForm) < 5 {
			return false
		}
		return distance <= 0.15
	default:
		return false
	}
}

// compareEvidence compares two MatchEvidence for sorting.
// Returns > 0 if a should come before b.
// Exact matches rank above fuzzy matches when weights are close.
func compareEvidence(a, b MatchEvidence) int {
	// Exact methods rank above all fuzzy methods
	aExact := a.Method == "exact" || a.Method == "case_insensitive"
	bExact := b.Method == "exact" || b.Method == "case_insensitive"

	if aExact && !bExact {
		return 1 // a comes first
	}
	if !aExact && bExact {
		return -1 // b comes first
	}

	// Otherwise, sort by weight descending
	if a.Weight != b.Weight {
		if a.Weight > b.Weight {
			return 1
		}
		return -1
	}

	// Tiebreaker: shorter distance first
	if a.Distance != b.Distance {
		if a.Distance < b.Distance {
			return 1
		}
		return -1
	}

	return 0
}

// fuzzyMatch compares two forms and returns the method and distance.
// Returns method name and Levenshtein distance.
func fuzzyMatch(a, b string) (string, float64) {
	// Exact match
	if a == b {
		return "exact", 0
	}

	// Normalized exact match (lowercase)
	if strings.ToLower(a) == strings.ToLower(b) {
		return "case_insensitive", 0
	}

	// Consonant skeleton match
	skelA := consonantSkeleton(a)
	skelB := consonantSkeleton(b)
	if skelA == skelB && len(skelA) > 0 {
		return "consonant_skeleton", 0.5
	}

	// Vowel skeleton match
	vowA := vowelSkeleton(a)
	vowB := vowelSkeleton(b)
	if vowA == vowB && len(vowA) > 0 {
		return "vowel_skeleton", 0.5
	}

	// Phonetic key match
	phonA := generatePhonetic(a)
	phonB := generatePhonetic(b)
	if phonA == phonB {
		return "phonetic", 0.3
	}

	// Levenshtein distance
	var dist float64 = float64(LevenshteinDistance(a, b))
	maxLen := math.Max(float64(len(a)), float64(len(b)))
	if maxLen > 0 {
		dist = dist / maxLen
	}

	return "levenshtein", dist
}

// LevenshteinDistance calculates the edit distance between two strings.
func LevenshteinDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	// Create matrix
	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
	}

	// Initialize first row and column
	for i := 0; i <= len(a); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(b); j++ {
		matrix[0][j] = j
	}

	// Fill in the rest
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			matrix[i][j] = minInt(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(a)][len(b)]
}

func minInt(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// calculateMatchWeight computes the weight of a match.
// Levenshtein matches fall off sharply - weak matches retain minimal weight.
func calculateMatchWeight(cand CandidateForm, anchor AnchorConcept, method string, distance float64) float64 {
	// Base weight from anchor
	weight := anchor.Weight

	// Adjust by candidate confidence
	switch cand.Confidence {
	case "verified":
		weight *= 1.0
	case "plausible":
		weight *= 0.8
	case "speculative":
		weight *= 0.5
	}

	// Adjust by match method
	switch method {
	case "exact":
		weight *= 1.0
	case "case_insensitive":
		weight *= 0.95
	case "consonant_skeleton":
		weight *= 0.8
	case "vowel_skeleton":
		weight *= 0.7
	case "phonetic":
		weight *= 0.75
	case "levenshtein":
		// Levenshtein falls off SHARPLY:
		// At distance 0.34 (max threshold), weight is only 15% of base
		// At distance 0.2, weight is ~50% of base
		// At distance 0.1, weight is ~75% of base
		weight *= (1.0 - distance*2.5) // Steep falloff
		if weight < 0.1 {              // Minimum floor for very weak matches
			weight = 0.1
		}
	}

	return weight
}

// =============================================================================
// Phase C: Data-Driven Anchor Knowledge
// =============================================================================

// GetAllAnchors returns all anchors for fuzzy matching from the knowledge base.
// This replaces hardcoded ModalAnchors, ShellAnchors, EmptyAnchors, BreathAnchors.
func GetAllAnchors(kb *knowledge.Knowledge) []AnchorConcept {
	if kb == nil {
		return nil
	}

	kbAnchors := kb.GetAllFormsAsAnchors()
	anchors := make([]AnchorConcept, len(kbAnchors))
	for i, a := range kbAnchors {
		anchors[i] = AnchorConcept{
			Form:       a.Form,
			Concept:    a.Concept,
			Confidence: a.Confidence,
			Weight:     a.Weight,
		}
	}
	return anchors
}
