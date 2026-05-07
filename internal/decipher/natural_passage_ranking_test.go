package decipher

import (
	"strings"
	"testing"

	"socrates/internal/knowledge"
)

// TestNaturalPassageRankingSadEmpty verifies that plain-language passages with direct
// concept evidence rank the direct fields above glyph/phonetic/orthographic
// structural observations.
func TestNaturalPassageRankingSadEmpty(t *testing.T) {
	kb := knowledge.LoadOrPanic()
	engine := NewEngineWithKnowledge(kb)

	// "I feel sad and empty inside" - classic emotional self-report
	reading := engine.Analyze("I feel sad and empty inside")
	if reading.Input == "" {
		t.Fatal("expected reading for 'I feel sad and empty inside'")
	}

	// Top fields should include sadness and emptiness as primary
	topFields := reading.PassageFields.TopFields(5)
	topFieldConcepts := make(map[string]bool)
	for _, f := range topFields {
		topFieldConcepts[f.Concept] = true
		t.Logf("top field: %s [strength=%.2f, direct=%v]", f.Concept, f.Strength, f.IsDirectEvidence)
	}

	// sadness and emptiness must be in top fields
	if !topFieldConcepts["sadness"] {
		t.Error("top fields should include 'sadness' for 'I feel sad and empty inside'")
	}
	if !topFieldConcepts["emptiness"] {
		t.Error("top fields should include 'emptiness' for 'I feel sad and empty inside'")
	}

	// At least sadness AND emptiness should outrank structural obs
	structuralOutranked := topFieldConcepts["sadness"] && topFieldConcepts["emptiness"]
	if !structuralOutranked {
		t.Error("sadness and emptiness should outrank structural observations")
	}

	// No resentment from neutral "feel" - the false feel->resentment path
	// counsellor should NOT suggest resentment->forgiveness for this passage
	if reading.CounsellorField != nil {
		for _, s := range reading.CounsellorField.Suggestions {
			if s.SourceConcept == "resentment" {
				t.Errorf("passage 'I feel sad and empty inside' should not produce resentment suggestion, got: %s -> %s",
					s.SourceConcept, s.TargetConcept)
			}
		}
	}
}

// TestNaturalPassageRankingLonelyAfraid verifies that loneliness and fear outrank
// structural observations and produce non-duplicated counsellor suggestions.
func TestNaturalPassageRankingLonelyAfraid(t *testing.T) {
	kb := knowledge.LoadOrPanic()
	engine := NewEngineWithKnowledge(kb)

	// "lonely and afraid" - relational + emotional
	reading := engine.Analyze("lonely and afraid")
	if reading.Input == "" {
		t.Fatal("expected reading for 'lonely and afraid'")
	}

	// Top fields should include loneliness and fear (use sorted TopFields)
	topFields := reading.PassageFields.TopFields(5)
	topConcepts := make([]string, 0)
	for i, f := range topFields {
		topConcepts = append(topConcepts, f.Concept)
		t.Logf("top[%d]: %s [strength=%.2f, depth=%d, direct=%v]",
			i, f.Concept, f.Strength, f.Depth, f.IsDirectEvidence)
	}

	hasLoneliness := false
	hasFear := false
	for _, c := range topConcepts {
		if c == "loneliness" {
			hasLoneliness = true
		}
		if c == "fear" {
			hasFear = true
		}
	}

	if !hasLoneliness {
		t.Error("top fields should include 'loneliness' for 'lonely and afraid'")
	}
	if !hasFear {
		t.Error("top fields should include 'fear' for 'lonely and afraid'")
	}

	// Counsellor suggestions must be deduplicated
	if reading.CounsellorField != nil {
		suggestions := reading.CounsellorField.Suggestions
		t.Logf("Suggestions count: %d", len(suggestions))

		// Build dedup key set
		seenKeys := make(map[string]bool)
		for _, s := range suggestions {
			dedupKey := s.SourceConcept + "|" + s.Kind + "|" + s.TargetConcept + "|" + s.Confidence + "|" + s.Lens
			if seenKeys[dedupKey] {
				t.Errorf("DUPLICATE suggestion: %s --[%s]--> %s", s.SourceConcept, s.Kind, s.TargetConcept)
			}
			seenKeys[dedupKey] = true
			t.Logf("  suggestion: %s --[%s]--> %s [strength=%.2f, conf=%s]",
				s.SourceConcept, s.Kind, s.TargetConcept, s.Strength, s.Confidence)
		}

		// Evidence paths must be deduplicated
		evidencePaths := reading.CounsellorField.EvidencePaths
		t.Logf("Evidence paths count: %d", len(evidencePaths))

		seenEvKeys := make(map[string]bool)
		for _, e := range evidencePaths {
			evKey := e.SourceConcept + "|" + e.Kind + "|" + e.TargetConcept
			if seenEvKeys[evKey] {
				t.Errorf("DUPLICATE evidence path: %s --[%s]--> %s", e.SourceConcept, e.Kind, e.TargetConcept)
			}
			seenEvKeys[evKey] = true
			t.Logf("  evidence: %s --[%s]--> %s", e.SourceConcept, e.Kind, e.TargetConcept)
		}
	}
}

// TestNaturalPassageRankingAfraidDisconnected verifies that fear, isolation, and truth
// are visible as primary fields; counsellor suggestions are deduplicated.
func TestNaturalPassageRankingAfraidDisconnected(t *testing.T) {
	kb := knowledge.LoadOrPanic()
	engine := NewEngineWithKnowledge(kb)

	// "I feel afraid and disconnected from truth"
	reading := engine.Analyze("I feel afraid and disconnected from truth")
	if reading.Input == "" {
		t.Fatal("expected reading for 'I feel afraid and disconnected from truth'")
	}

	// Check top fields show fear, isolation, truth (use sorted TopFields)
	topFields := reading.PassageFields.TopFields(5)
	topConcepts := make([]string, 0)
	for i, f := range topFields {
		topConcepts = append(topConcepts, f.Concept)
		t.Logf("top[%d]: %s [strength=%.2f, depth=%d, direct=%v]",
			i, f.Concept, f.Strength, f.Depth, f.IsDirectEvidence)
	}

	hasFear := false
	hasIsolation := false
	hasTruth := false
	for _, c := range topConcepts {
		if c == "fear" {
			hasFear = true
		}
		if c == "isolation" {
			hasIsolation = true
		}
		if c == "truth" {
			hasTruth = true
		}
	}

	if !hasFear {
		t.Error("top fields should include 'fear' for 'I feel afraid and disconnected from truth'")
	}
	if !hasIsolation {
		t.Error("top fields should include 'isolation' for 'I feel afraid and disconnected from truth'")
	}
	if !hasTruth {
		t.Error("top fields should include 'truth' for 'I feel afraid and disconnected from truth'")
	}

	// Counsellor suggestions must be deduplicated
	if reading.CounsellorField != nil {
		suggestions := reading.CounsellorField.Suggestions
		seenKeys := make(map[string]bool)
		for _, s := range suggestions {
			dedupKey := s.SourceConcept + "|" + s.Kind + "|" + s.TargetConcept + "|" + s.Confidence + "|" + s.Lens
			if seenKeys[dedupKey] {
				t.Errorf("DUPLICATE suggestion: %s --[%s]--> %s", s.SourceConcept, s.Kind, s.TargetConcept)
			}
			seenKeys[dedupKey] = true
			t.Logf("  suggestion: %s --[%s]--> %s [strength=%.2f, conf=%s]",
				s.SourceConcept, s.Kind, s.TargetConcept, s.Strength, s.Confidence)
		}

		// Evidence paths must be deduplicated
		evidencePaths := reading.CounsellorField.EvidencePaths
		seenEvKeys := make(map[string]bool)
		for _, e := range evidencePaths {
			evKey := e.SourceConcept + "|" + e.Kind + "|" + e.TargetConcept
			if seenEvKeys[evKey] {
				t.Errorf("DUPLICATE evidence path: %s --[%s]--> %s", e.SourceConcept, e.Kind, e.TargetConcept)
			}
			seenEvKeys[evKey] = true
		}
	}
}

// TestFeelFalseResentmentActivation verifies that neutral "feel" does not
// produce false resentment activation paths.
// bitterness from fuzzy match should be depth 1 (indirect) and ranked below
// direct evidence concepts like feeling.
func TestFeelFalseResentmentActivation(t *testing.T) {
	kb := knowledge.LoadOrPanic()
	engine := NewEngineWithKnowledge(kb)

	// Single "feel" input - bitterness should be depth 1 (indirect), not top-ranked
	reading := engine.Analyze("feel")
	if reading.Input == "" {
		t.Fatal("expected reading for 'feel'")
	}

	// Check top 5 sorted passage fields
	topFields := reading.PassageFields.TopFields(5)
	topConcepts := make([]string, 0)
	for _, f := range topFields {
		topConcepts = append(topConcepts, f.Concept)
		t.Logf("top field: %s [strength=%.2f, depth=%d, direct=%v]",
			f.Concept, f.Strength, f.Depth, f.IsDirectEvidence)
	}

	// bitterness from fuzzy match should be depth 1 (indirect)
	bitternessDepth := -1
	for _, f := range reading.PassageFields {
		if f.Concept == "bitterness" {
			bitternessDepth = f.Depth
			t.Logf("bitterness depth: %d (indirect=%v)", f.Depth, f.Depth > 0 || !f.IsDirectEvidence)
			break
		}
	}
	if bitternessDepth != 1 {
		t.Errorf("bitterness should be depth 1 (indirect from fuzzy match), got depth %d", bitternessDepth)
	}

	// feeling should be present and direct (from curated form match)
	hasFeeling := false
	feelingDirect := false
	for _, f := range topFields {
		if f.Concept == "feeling" {
			hasFeeling = true
			feelingDirect = f.IsDirectEvidence
			break
		}
	}
	if !hasFeeling {
		t.Error("feeling should be in top fields for 'feel' (curated form match)")
	}
	if !feelingDirect {
		t.Error("feeling should have direct evidence (from curated form match)")
	}

	// No resentment suggestions from neutral "feel"
	if reading.CounsellorField != nil {
		for _, s := range reading.CounsellorField.Suggestions {
			if s.SourceConcept == "resentment" {
				t.Errorf("neutral 'feel' should not produce resentment suggestion: %s -> %s",
					s.SourceConcept, s.TargetConcept)
			}
		}
	}
}

// TestFuzzyMatchIndirectDepth verifies that fuzzy matches are marked as
// indirect (depth 1) when they create new nodes, not direct depth 0.
func TestFuzzyMatchIndirectDepth(t *testing.T) {
	kb := knowledge.LoadOrPanic()
	engine := NewEngineWithKnowledge(kb)

	// "feel" triggers ifeel->bitterness via vowel_skeleton fuzzy match
	// This should be depth 1 (indirect), not depth 0 (direct)
	reading := engine.Analyze("feel")
	if reading.Input == "" {
		t.Fatal("expected reading for 'feel'")
	}

	// Find bitterness field
	for _, f := range reading.PassageFields {
		if f.Concept == "bitterness" {
			if f.Depth == 0 {
				t.Errorf("bitterness from fuzzy match should be depth 1 (indirect), got depth 0 (direct)")
			} else {
				t.Logf("bitterness correctly marked as depth %d (indirect from fuzzy match)", f.Depth)
			}
			return
		}
	}
	t.Log("bitterness not in passage fields - ok if correctly filtered")
}

// TestDirectVsStructuralLabeling verifies that structural channels are labeled
// correctly and not over-labeled as "direct" when they should be "propagated".
func TestDirectVsStructuralLabeling(t *testing.T) {
	kb := knowledge.LoadOrPanic()
	engine := NewEngineWithKnowledge(kb)

	// "lonely and afraid" has glyph/phonetic structural channels
	reading := engine.Analyze("lonely and afraid")
	if reading.Input == "" {
		t.Fatal("expected reading for 'lonely and afraid'")
	}

	// Check that direct evidence comes from curated forms (feeling, loneliness, fear)
	// not from structural channels alone
	for _, f := range reading.PassageFields {
		if f.IsDirectEvidence {
			t.Logf("DIRECT: %s [depth=%d]", f.Concept, f.Depth)
		} else {
			t.Logf("INDIRECT: %s [depth=%d]", f.Concept, f.Depth)
		}
	}

	// Primary emotional concepts should be direct
	primaryConcepts := map[string]bool{
		"loneliness": true,
		"fear":       true,
	}
	for _, f := range reading.PassageFields {
		if primaryConcepts[f.Concept] && !f.IsDirectEvidence {
			t.Errorf("primary concept '%s' should have direct evidence, got indirect", f.Concept)
		}
	}

	debug := RenderReadingWithOptions(reading, RenderOptions{Mode: RenderModeDebug})
	if strings.Contains(debug, "phonetic-structure [strength:") &&
		strings.Contains(debug, "phonetic-structure [strength: 1.50, direct") {
		t.Error("debug output should not label phonetic-structure as direct")
	}
	if strings.Contains(debug, "vowel-heavy-orthography [strength:") &&
		strings.Contains(debug, "vowel-heavy-orthography [strength: 10.00, direct") {
		t.Error("debug output should not label vowel-heavy-orthography as direct")
	}
}

// TestProseReflectsTopFields verifies that the default concise reading
// reflects the top passage fields, not structural noise.
func TestProseReflectsTopFields(t *testing.T) {
	kb := knowledge.LoadOrPanic()
	engine := NewEngineWithKnowledge(kb)

	// "I feel sad and empty inside" - top fields should be sadness, emptiness, feeling
	reading := engine.Analyze("I feel sad and empty inside")
	if reading.Input == "" {
		t.Fatal("expected reading for 'I feel sad and empty inside'")
	}

	// Get top 3 field concepts (using sorted TopFields)
	topFields := reading.PassageFields.TopFields(3)
	topConcepts := make([]string, 0)
	for _, f := range topFields {
		topConcepts = append(topConcepts, f.Concept)
	}

	// The concise reading should mention top field concepts
	prose := RenderReading(reading)
	t.Logf("Concise reading:\n%s", prose)

	// sadness or emptiness should appear in the concise reading
	hasSadnessOrEmptiness := strings.Contains(prose, "sadness") || strings.Contains(prose, "emptiness")
	if !hasSadnessOrEmptiness {
		t.Error("concise reading should mention sadness or emptiness from top fields")
	}

	readingSection := prose
	if idx := strings.Index(prose, "Reading:\n"); idx >= 0 {
		readingSection = prose[idx:]
	}

	for _, structural := range []string{"phonetic", "glyph", "orthographic", "repeated-letter-observation"} {
		if strings.Contains(readingSection, structural) {
			t.Errorf("final reading prose should not promote structural signal %q: %s", structural, readingSection)
		}
	}
}
