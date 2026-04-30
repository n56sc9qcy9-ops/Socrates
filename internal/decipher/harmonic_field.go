package decipher

import (
	"socrates/internal/knowledge"
)

// BuildHarmonicField constructs a harmonic field from active passage fields.
// The field is built from data-backed frequency profiles, not from hardcoded Go constants.
func BuildHarmonicField(pf PassageFields, kb *knowledge.Knowledge) *HarmonicField {
	if pf == nil || len(pf) == 0 {
		return nil
	}

	// Collect all activated concepts and their strengths
	activatedConcepts := make(map[string]float64)
	for _, field := range pf {
		activatedConcepts[field.Concept] = field.Strength
	}

	if len(activatedConcepts) == 0 {
		return nil
	}

	// Look up frequency profiles for all activated concepts
	profileTones := make(map[string]*HarmonicTone) // key by meaning_frequency_id
	evidencePaths := make([]HarmonicEvidence, 0)

	for conceptID, conceptStrength := range activatedConcepts {
		profiles := kb.GetFrequencyProfilesByConcept(conceptID)

		for _, fp := range profiles {
			// Create or update tone for this meaning_frequency_id
			tone, exists := profileTones[fp.MeaningFrequencyID]
			if !exists {
				// New tone from this profile
				tone = &HarmonicTone{
					MeaningFrequencyID: fp.MeaningFrequencyID,
					SourceConcepts:     make([]string, 0),
					Vector:             [3]int{0, 0, 0},
					Ratio:              [2]int{1, 1},
					Archetype:          fp.Archetype,
					Labels:             fp.Labels,
					Strength:           0,
					Confidence:         fp.Confidence,
				}

				// Initialize vector and ratio from profile
				if len(fp.Vector) >= 3 {
					tone.Vector = [3]int{fp.Vector[0], fp.Vector[1], fp.Vector[2]}
				}
				if len(fp.Ratio) >= 2 {
					tone.Ratio = [2]int{fp.Ratio[0], fp.Ratio[1]}
				}

				profileTones[fp.MeaningFrequencyID] = tone
			}

			// Add this concept as a source
			tone.SourceConcepts = append(tone.SourceConcepts, conceptID)

			// Calculate contribution strength
			// Use profile weight (0-100) and concept strength (0-1) for combined strength
			profileStrength := float64(fp.Weight) / 100.0 // normalize to 0-1
			contribution := conceptStrength * profileStrength

			// Accumulate strength (weighted average for multiple sources)
			if len(tone.SourceConcepts) == 1 {
				tone.Strength = contribution
			} else {
				// Running average
				n := float64(len(tone.SourceConcepts))
				oldAvg := tone.Strength * (n - 1) / n
				tone.Strength = oldAvg + contribution/n
			}

			// Record evidence path
			evidencePaths = append(evidencePaths, HarmonicEvidence{
				SourceConcept:       conceptID,
				SourceConceptWeight: conceptStrength,
				MeaningFrequencyID:  fp.MeaningFrequencyID,
				ProfileWeight:       fp.Weight,
				Confidence:          fp.Confidence,
			})
		}
	}

	if len(profileTones) == 0 {
		return nil
	}

	// Collect tones
	tones := make([]HarmonicTone, 0, len(profileTones))
	for _, tone := range profileTones {
		tones = append(tones, *tone)
	}

	// Calculate coherence, consonance, and dissonance
	coherence, consonance, dissonance := calculateHarmonicMetrics(tones)

	return &HarmonicField{
		Tones:        tones,
		Coherence:    coherence,
		Consonance:   consonance,
		Dissonance:   dissonance,
		EvidencePaths: evidencePaths,
	}
}

// calculateHarmonicMetrics computes coherence, consonance, and dissonance scores.
func calculateHarmonicMetrics(tones []HarmonicTone) (coherence, consonance, dissonance float64) {
	if len(tones) < 2 {
		// Single tone or no tones - no tension
		if len(tones) == 1 {
			return 1.0, 1.0, 0.0
		}
		return 0.0, 0.0, 0.0
	}

	var compatibleCount, incompatibleCount int
	var totalWeight float64

	for i := 0; i < len(tones); i++ {
		for j := i + 1; j < len(tones); j++ {
			t1 := tones[i]
			t2 := tones[j]

			// Check if ratios are compatible (share common factors or form coherent intervals)
			compatible := areRatiosCompatible(t1.Ratio, t2.Ratio)

			// Weight by average strength
			weight := (t1.Strength + t2.Strength) / 2.0
			totalWeight += weight

			if compatible {
				compatibleCount++
			} else {
				incompatibleCount++
			}
		}
	}

	totalPairs := compatibleCount + incompatibleCount
	if totalPairs == 0 {
		return 1.0, 1.0, 0.0
	}

	// Consonance = compatible pairs / total pairs
	consonance = float64(compatibleCount) / float64(totalPairs)

	// Dissonance = incompatible pairs / total pairs
	dissonance = float64(incompatibleCount) / float64(totalPairs)

	// Coherence = weighted average considering tone strengths
	// High strength tones contribute more to overall coherence
	var weightedCoherence float64
	for _, tone := range tones {
		// Each tone contributes to coherence based on its strength
		weightedCoherence += tone.Strength
	}
	coherence = weightedCoherence / float64(len(tones))

	return coherence, consonance, dissonance
}

// areRatiosCompatible checks if two frequency ratios are harmonically compatible.
// Ratios are compatible if they share a common factor or if their combination
// forms a recognized harmonic interval.
func areRatiosCompatible(r1, r2 [2]int) bool {
	// If denominators are zero, not compatible
	if r1[1] == 0 || r2[1] == 0 {
		return false
	}

	// Check if ratios are identical
	if r1 == r2 {
		return true
	}

	// Calculate cross product for shared factors
	// For ratios a/b and c/d, compatibility is based on whether a*d and b*c share factors

	// Check for simple octave relationships
	// 1/1 is compatible with everything
	if (r1[0] == 1 && r1[1] == 1) || (r2[0] == 1 && r2[1] == 1) {
		return true
	}

	// Check if ratios share a common factor (reducing to simplest form)
	// Use GCD to normalize
	gcd1 := gcd(r1[0], r1[1])
	gcd2 := gcd(r2[0], r2[1])

	n1 := r1[0] / gcd1
	d1 := r1[1] / gcd1
	n2 := r2[0] / gcd2
	d2 := r2[1] / gcd2

	// If normalized forms match, compatible
	if n1 == n2 && d1 == d2 {
		return true
	}

	// Check if one is an octave of the other
	// For ratios a/b and c/d: if a*d = 2*b*c, they differ by an octave
	// Or: a*d = 4*b*c for double octave, etc.
	if n1*d2 == 2*n2*d1 || n2*d1 == 2*n1*d2 {
		return true
	}

	// Check for common factors in numerators (fifth relationships)
	// E.g., 3/2 (fifth) and 4/3 (fourth) are compatible
	gcdNum := gcd(n1, n2)
	if gcdNum > 1 {
		return true // Share a common factor in numerators
	}

	// Check for common factors in denominators
	gcdDen := gcd(d1, d2)
	if gcdDen > 1 {
		return true // Share a common factor in denominators
	}

	// Pythagorean compatibility: check if ratios come from the same harmonic series
	// Ratios are compatible if their cross-multiplication produces a Pythagorean triple
	// Check: a*d and b*c share a factor with their sum or difference
	prod1 := n1 * d2
	prod2 := n2 * d1
	gcdProd := gcd(prod1, prod2)

	// If the cross products share significant common factor, compatible
	if gcdProd > 1 {
		return true
	}

	// Fifths are always compatible (3/2 is the base Pythagorean ratio)
	// Check for perfect fifth (3/2) or its octave variants
	fifthRatios := [][2]int{
		{3, 2}, {6, 4}, {9, 6}, {12, 8}, // multiples of 3/2
	}
	for _, fr := range fifthRatios {
		if (n1 == fr[0] && d1 == fr[1]) || (n2 == fr[0] && d2 == fr[1]) {
			return true
		}
	}

	// Thirds (5/4) are compatible with fifths (3/2)
	thirdRatios := [][2]int{
		{5, 4}, {10, 8}, {15, 12}, // multiples of 5/4
	}
	for _, tr := range thirdRatios {
		if (n1 == tr[0] && d1 == tr[1]) || (n2 == tr[0] && d2 == tr[1]) {
			return true
		}
	}

	// Not compatible - these are dissonant ratios
	return false
}

// gcd computes the greatest common divisor using Euclidean algorithm.
func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

// RenderHarmonicField renders a harmonic field to a string.
// In default mode, shows summary only when profiles exist.
// In debug mode, shows full profile details.
func RenderHarmonicField(hf *HarmonicField, mode RenderMode) string {
	if hf == nil {
		return ""
	}

	if mode == RenderModeDefault {
		return hfRenderDefault(hf)
	}
	return hfRenderDebug(hf)
}

func hfRenderDefault(hf *HarmonicField) string {
	if len(hf.Tones) == 0 {
		return ""
	}

	// Show tone count and coherence summary
	s := "HarmonicField(" + intToStr(len(hf.Tones)) + " tones"
	s += ", coherence: " + floatToStr(hf.Coherence)
	if hf.Consonance > hf.Dissonance {
		s += ", consonant"
	} else if hf.Dissonance > 0.3 {
		s += ", dissonant"
	}
	s += ")"
	return s
}

func hfRenderDebug(hf *HarmonicField) string {
	if len(hf.Tones) == 0 {
		return "HarmonicField: no tones"
	}

	s := "HarmonicField:\n"
	s += "  Coherence: " + floatToStr(hf.Coherence) + "\n"
	s += "  Consonance: " + floatToStr(hf.Consonance) + "\n"
	s += "  Dissonance: " + floatToStr(hf.Dissonance) + "\n"
	s += "  Tones:\n"

	for _, tone := range hf.Tones {
		s += "    - " + tone.MeaningFrequencyID + "\n"
		s += "      Vector: [" + intToStr(tone.Vector[0]) + "," + intToStr(tone.Vector[1]) + "," + intToStr(tone.Vector[2]) + "]\n"
		s += "      Ratio: " + intToStr(tone.Ratio[0]) + "/" + intToStr(tone.Ratio[1]) + "\n"
		s += "      Archetype: " + tone.Archetype + "\n"
		s += "      Labels: note=" + intToStr(tone.Labels.Note) + ", color=" + intToStr(tone.Labels.Color) + ", field=" + intToStr(tone.Labels.Field) + "\n"
		s += "      Strength: " + floatToStr(tone.Strength) + "\n"
		s += "      Confidence: " + tone.Confidence + "\n"
		s += "      Sources: " + strJoin(tone.SourceConcepts, ", ") + "\n"
	}

	s += "  EvidencePaths (" + intToStr(len(hf.EvidencePaths)) + "):\n"
	for _, ep := range hf.EvidencePaths {
		s += "    - " + ep.SourceConcept + " → " + ep.MeaningFrequencyID
		s += " (concept_w=" + floatToStr(ep.SourceConceptWeight) + ", profile_w=" + intToStr(ep.ProfileWeight) + ")\n"
	}

	return s
}

func intToStr(i int) string {
	if i == 0 {
		return "0"
	}
	neg := ""
	if i < 0 {
		neg = "-"
		i = -i
	}
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return neg + s
}

func floatToStr(f float64) string {
	// Simple formatting to 2 decimal places
	i := int(f * 100)
	return intToStr(i/100) + "." + intToStr((i%100)/10) + intToStr(i%10)
}

func strJoin(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	s := strs[0]
	for i := 1; i < len(strs); i++ {
		s += sep + strs[i]
	}
	return s
}

// GetHarmonicFieldSummary returns a concise summary of the harmonic field.
func GetHarmonicFieldSummary(hf *HarmonicField) string {
	if hf == nil || len(hf.Tones) == 0 {
		return ""
	}

	toneNames := make([]string, 0, len(hf.Tones))
	for _, t := range hf.Tones {
		toneNames = append(toneNames, t.MeaningFrequencyID)
	}

	return "[" + strJoin(toneNames, ", ") + "]"
}