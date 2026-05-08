package decipher

import (
	"strings"
)

// RenderDescriptor contains deterministic textual render instructions
// for the harmonic field. All values come from integer profile data.
// No audio playback, no floating-point meaning storage.
type RenderDescriptor struct {
	// Field provides context about the overall harmonic field.
	Field FieldDescriptor

	// Tones are the individual tone render instructions.
	Tones []ToneDescriptor

	// Relationships describe harmonic interactions between tones.
	Relationships []RelationshipDescriptor

	// SourceTrace documents which concepts contributed to which tones.
	SourceTrace []SourceTraceEntry
}

// FieldDescriptor provides field-level context.
type FieldDescriptor struct {
	// ToneCount is the number of active tones.
	ToneCount int

	// Coherence is the overall harmonic coherence (0-1).
	Coherence float64

	// ConsonanceScore indicates consonant relationships.
	ConsonanceScore float64

	// DissonanceScore indicates dissonant relationships.
	DissonanceScore float64

	// DominantQuality is "consonant", "dissonant", or "neutral".
	DominantQuality string

	// ArchetypeSummary describes the geometric/spiritual archetype structure.
	ArchetypeSummary string
}

// ToneDescriptor describes a single tone for rendering.
type ToneDescriptor struct {
	// ID is the stable meaning-frequency identity.
	ID string

	// Vector is the 3D tone vector as integers [Tone1, Tone2, Tone3].
	Vector [3]int

	// Ratio is the frequency ratio as integers [numerator, denominator].
	Ratio [2]int

	// Archetype is the archetype ID (e.g., pythagorean_triple_1, phi_approximant_1).
	Archetype string

	// ArchetypeName is the human-readable archetype name.
	ArchetypeName string

	// Labels provides semantic render hints as integers.
	Labels LabelDescriptor

	// Strength indicates the combined strength from contributing concepts.
	Strength float64

	// Confidence is the confidence level of the source profile.
	Confidence string
}

// LabelDescriptor provides semantic labels as integers for rendering.
type LabelDescriptor struct {
	// Note is the integer note index (0-11 semitones, or extended range).
	Note int

	// Color is the integer hue (0-360 degrees).
	Color int

	// Field is the integer field ID.
	Field int
}

// RelationshipDescriptor describes a harmonic relationship between tones.
type RelationshipDescriptor struct {
	// From is the source tone ID.
	From string

	// To is the target tone ID.
	To string

	// Type is "consonant" or "dissonant".
	Type string

	// Description is a textual description of the relationship.
	Description string
}

// SourceTraceEntry documents concept-to-tone source mapping.
type SourceTraceEntry struct {
	// Concept is the source concept ID.
	Concept string

	// ToneID is the meaning-frequency ID of the tone.
	ToneID string

	// Weight is the contribution weight (0-1).
	Weight float64
}

// ToRenderDescriptor converts a HarmonicField to a RenderDescriptor.
// All values are deterministic from integer profile data.
func (hf *HarmonicField) ToRenderDescriptor(archetypes map[string]string) *RenderDescriptor {
	if hf == nil || len(hf.Tones) == 0 {
		return nil
	}

	// Build field descriptor
	field := FieldDescriptor{
		ToneCount:        len(hf.Tones),
		Coherence:        hf.Coherence,
		ConsonanceScore:  hf.Consonance,
		DissonanceScore:  hf.Dissonance,
		DominantQuality:  dominantQuality(hf.Consonance, hf.Dissonance),
		ArchetypeSummary: archetypeSummary(hf.Tones),
	}

	// Build tone descriptors
	toneDescs := make([]ToneDescriptor, len(hf.Tones))
	archetypeNames := loadArchetypeNames()
	for i, tone := range hf.Tones {
		archName := archetypeNames[tone.Archetype]
		if archName == "" {
			archName = tone.Archetype
		}

		toneDescs[i] = ToneDescriptor{
			ID:            tone.MeaningFrequencyID,
			Vector:        tone.Vector,
			Ratio:         tone.Ratio,
			Archetype:     tone.Archetype,
			ArchetypeName: archName,
			Labels: LabelDescriptor{
				Note:  tone.Labels.Note,
				Color: tone.Labels.Color,
				Field: tone.Labels.Field,
			},
			Strength:   tone.Strength,
			Confidence: tone.Confidence,
		}
	}

	// Build relationship descriptors
	relDescs := buildRelationshipDescriptors(hf.Tones)

	// Build source trace from evidence paths
	var sourceTrace []SourceTraceEntry
	seen := make(map[string]bool) // dedup by concept+tone
	for _, ep := range hf.EvidencePaths {
		key := ep.SourceConcept + "|" + ep.MeaningFrequencyID
		if seen[key] {
			continue
		}
		seen[key] = true
		sourceTrace = append(sourceTrace, SourceTraceEntry{
			Concept: ep.SourceConcept,
			ToneID:  ep.MeaningFrequencyID,
			Weight:  ep.SourceConceptWeight * float64(ep.ProfileWeight) / 100.0,
		})
	}

	return &RenderDescriptor{
		Field:        field,
		Tones:        toneDescs,
		Relationships: relDescs,
		SourceTrace:  sourceTrace,
	}
}

// dominantQuality determines the dominant harmonic quality.
func dominantQuality(consonance, dissonance float64) string {
	if consonance > dissonance {
		return "consonant"
	}
	if dissonance > 0.3 {
		return "dissonant"
	}
	return "neutral"
}

// archetypeSummary generates a textual summary of the archetype structure.
func archetypeSummary(tones []HarmonicTone) string {
	archetypes := make([]string, 0, len(tones))
	for _, t := range tones {
		if t.Archetype != "" {
			archetypes = append(archetypes, t.Archetype)
		}
	}

	if len(archetypes) == 0 {
		return "undetermined"
	}

	// Count archetype occurrences
	count := make(map[string]int)
	for _, a := range archetypes {
		count[a]++
	}

	// Build summary
	var summary []string
	for arch, c := range count {
		if c > 1 {
			summary = append(summary, arch+" (×"+intToStr(c)+")")
		} else {
			summary = append(summary, arch)
		}
	}

	return strJoin(summary, ", ")
}

// loadArchetypeNames returns the archetype ID to name mapping.
func loadArchetypeNames() map[string]string {
	return map[string]string{
		"pythagorean_triple_1": "Monad",
		"pythagorean_triple_2": "Triangular",
		"pythagorean_triple_3": "Triangular",
		"pythagorean_triple_4": "Triangular",
		"pythagorean_triple_5": "Primitive",
		"pythagorean_triple_6": "Primitive",
		"phi_approximant_1":    "Golden Ratio",
		"phi_approximant_2":    "Extended Golden",
		"metatron_node_1":      "Metatron Node",
	}
}

// buildRelationshipDescriptors creates relationship descriptions between tones.
func buildRelationshipDescriptors(tones []HarmonicTone) []RelationshipDescriptor {
	var rels []RelationshipDescriptor

	for i := 0; i < len(tones); i++ {
		for j := i + 1; j < len(tones); j++ {
			compatible := areRatiosCompatible(tones[i].Ratio, tones[j].Ratio)
			relType := "dissonant"
			description := "dissonant interval"

			if compatible {
				relType = "consonant"
				description = describeRelationship(tones[i].Ratio, tones[j].Ratio)
			}

			rels = append(rels, RelationshipDescriptor{
				From:        tones[i].MeaningFrequencyID,
				To:          tones[j].MeaningFrequencyID,
				Type:        relType,
				Description: description,
			})
		}
	}

	return rels
}

// describeRelationship provides textual description of the interval relationship.
func describeRelationship(r1, r2 [2]int) string {
	// Check for unison
	if r1[0] == r1[1] && r2[0] == r2[1] && r1[0] == r2[0] {
		return "unison"
	}

	// Check for octave relationships
	if r1[0] == 2*r2[0] && r1[1] == r2[1] {
		return "octave (descending)"
	}
	if r2[0] == 2*r1[0] && r2[1] == r1[1] {
		return "octave (ascending)"
	}

	// Check for perfect fifth (3:2)
	if (r1[0] == 3 && r1[1] == 2 && r2[0] == 1 && r2[1] == 1) ||
		(r1[0] == 1 && r1[1] == 1 && r2[0] == 3 && r2[1] == 2) {
		return "perfect fifth"
	}
	if (r1[0] == 3 && r1[1] == 2 && r2[0] == 4 && r2[1] == 3) ||
		(r1[0] == 4 && r1[1] == 3 && r2[0] == 3 && r2[1] == 2) {
		return "perfect fifth"
	}

	// Check for perfect fourth (4:3)
	if (r1[0] == 4 && r1[1] == 3 && r2[0] == 1 && r2[1] == 1) ||
		(r1[0] == 1 && r1[1] == 1 && r2[0] == 4 && r2[1] == 3) {
		return "perfect fourth"
	}

	// Check for major third (5:4)
	if (r1[0] == 5 && r1[1] == 4 && r2[0] == 1 && r2[1] == 1) ||
		(r1[0] == 1 && r1[1] == 1 && r2[0] == 5 && r2[1] == 4) {
		return "major third"
	}

	// Check for minor third (6:5)
	if (r1[0] == 6 && r1[1] == 5 && r2[0] == 1 && r2[1] == 1) ||
		(r1[0] == 1 && r1[1] == 1 && r2[0] == 6 && r2[1] == 5) {
		return "minor third"
	}

	// Check for common factors (harmonic series)
	gcdNum := gcd(r1[0], r2[0])
	if gcdNum > 1 {
		return "harmonic series (shared factor " + intToStr(gcdNum) + ")"
	}

	// Generic consonant description
	return "consonant interval"
}

// RenderDescriptorToString converts a RenderDescriptor to a human-readable string.
func RenderDescriptorToString(rd *RenderDescriptor) string {
	if rd == nil {
		return ""
	}

	var sb strings.Builder

	// Field summary
	sb.WriteString("HarmonicField[\n")
	sb.WriteString("  field: tone_count=" + intToStr(rd.Field.ToneCount))
	sb.WriteString(", coherence=" + floatToStr(rd.Field.Coherence))
	sb.WriteString(", quality=" + rd.Field.DominantQuality)
	sb.WriteString(", archetypes=[" + rd.Field.ArchetypeSummary + "]\n")

	// Tones
	sb.WriteString("  tones:\n")
	for _, t := range rd.Tones {
		sb.WriteString("    " + t.ID + "\n")
		sb.WriteString("      vector=" + vectorToStr(t.Vector))
		sb.WriteString(", ratio=" + ratioToStr(t.Ratio))
		sb.WriteString(", archetype=" + t.ArchetypeName + "\n")
		sb.WriteString("      labels: note=" + intToStr(t.Labels.Note))
		sb.WriteString(", color=" + intToStr(t.Labels.Color))
		sb.WriteString(", field=" + intToStr(t.Labels.Field) + "\n")
		sb.WriteString("      strength=" + floatToStr(t.Strength) + ", confidence=" + t.Confidence + "\n")
	}

	// Relationships (only if more than one tone)
	if len(rd.Relationships) > 0 {
		sb.WriteString("  relationships:\n")
		for _, r := range rd.Relationships {
			sb.WriteString("    " + r.From + " --[" + r.Type + "]--> " + r.To)
			sb.WriteString(": " + r.Description + "\n")
		}
	}

	// Source trace
	if len(rd.SourceTrace) > 0 {
		sb.WriteString("  sources:\n")
		for _, s := range rd.SourceTrace {
			sb.WriteString("    " + s.Concept + " -> " + s.ToneID)
			sb.WriteString(" (weight=" + floatToStr(s.Weight) + ")\n")
		}
	}

	sb.WriteString("]\n")

	return sb.String()
}

// vectorToStr formats a 3D integer vector.
func vectorToStr(v [3]int) string {
	return "[" + intToStr(v[0]) + "," + intToStr(v[1]) + "," + intToStr(v[2]) + "]"
}

// ratioToStr formats a frequency ratio.
func ratioToStr(r [2]int) string {
	return intToStr(r[0]) + "/" + intToStr(r[1])
}
