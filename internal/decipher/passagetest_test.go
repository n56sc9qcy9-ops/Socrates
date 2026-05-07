package decipher

import (
	"testing"
	"socrates/internal/knowledge"
)

// TestSanskritDevanagariPassageFields verifies Devanagari ScriptWord concepts
// flow through passage/harmonic field activation per current task acceptance
// criteria. All assertions check the generic activation path, not hardcoded
// behavior. This test STRENGTHENS the existing Sanskrit suite by requiring
// evidence to reach passage fields (not just channels).
func TestSanskritDevanagariPassageFields(t *testing.T) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		t.Fatal(err)
	}
	engine := NewEngineWithKnowledge(kb)

	// -----------------------------------------------------------------
	// Criterion 1: प्राण has breath and life-force as passage fields.
	// -----------------------------------------------------------------
	t.Run("PranaBreathLifeForce", func(t *testing.T) {
		reading := engine.Analyze("प्राण")
		fields := reading.PassageFields

		// breath must appear as a passage field (not only channel evidence)
		foundBreath := false
		for _, pf := range fields {
			if pf.Concept == "breath" {
				foundBreath = true
				if !pf.IsDirectEvidence {
					t.Errorf("breath should be direct evidence (got direct=%v)", pf.IsDirectEvidence)
				}
				if pf.Depth != 0 {
					t.Errorf("breath should be depth 0 (got %d)", pf.Depth)
				}
				if pf.Confidence != ConfidenceVerified {
					t.Errorf("breath should be verified confidence (got %s)", pf.Confidence)
				}
			}
		}
		if !foundBreath {
			t.Errorf("प्राण should have 'breath' in passage fields; got %d fields: %v",
				len(fields), conceptNames(fields))
		}

		// life-force must appear as a passage field
		foundLifeForce := false
		for _, pf := range fields {
			if pf.Concept == "life-force" {
				foundLifeForce = true
				if !pf.IsDirectEvidence {
					t.Errorf("life-force should be direct evidence")
				}
			}
		}
		if !foundLifeForce {
			t.Errorf("प्राण should have 'life-force' in passage fields; got %d fields: %v",
				len(fields), conceptNames(fields))
		}
	})

	// -----------------------------------------------------------------
	// Criterion 2: सत्य has truth as a passage field.
	// -----------------------------------------------------------------
	t.Run("SatyaTruth", func(t *testing.T) {
		reading := engine.Analyze("सत्य")
		fields := reading.PassageFields

		foundTruth := false
		for _, pf := range fields {
			if pf.Concept == "truth" {
				foundTruth = true
				if pf.Depth != 0 {
					t.Errorf("truth should be depth 0 (got %d)", pf.Depth)
				}
				if pf.Confidence != ConfidenceVerified {
					t.Errorf("truth should be verified confidence (got %s)", pf.Confidence)
				}
			}
		}
		if !foundTruth {
			t.Errorf("सत्य should have 'truth' in passage fields; got %d fields: %v",
				len(fields), conceptNames(fields))
		}
	})

	// -----------------------------------------------------------------
	// Criterion 3: Known Sanskrit with frequency profile produces harmonic field.
	// breath and prana have curated frequency profiles; ॐ does too.
	// -----------------------------------------------------------------
	t.Run("PranaHarmonicField", func(t *testing.T) {
		reading := engine.Analyze("प्राण")
		if reading.HarmonicField == nil {
			t.Error("प्राण should produce a harmonic field (breath has frequency profiles)")
		} else {
			if len(reading.HarmonicField.Tones) == 0 {
				t.Error("प्राण harmonic field should have at least one tone")
			}
			// Coherence should be positive
			if reading.HarmonicField.Coherence <= 0 {
				t.Errorf("प्राण harmonic field coherence should be >0 (got %.2f)",
					reading.HarmonicField.Coherence)
			}
			t.Logf("प्राण harmonic: %d tones, coherence=%.2f",
				len(reading.HarmonicField.Tones), reading.HarmonicField.Coherence)
		}
	})

	t.Run("OmHarmonicField", func(t *testing.T) {
		reading := engine.Analyze("ॐ")
		if reading.HarmonicField == nil {
			t.Error("ॐ should produce a harmonic field (sacred-sound/totality have profiles)")
		} else {
			t.Logf("ॐ harmonic: %d tones, coherence=%.2f",
				len(reading.HarmonicField.Tones), reading.HarmonicField.Coherence)
		}
	})

	// -----------------------------------------------------------------
	// Criterion 4: Unknown कवि remains weak/speculative.
	// -----------------------------------------------------------------
	t.Run("UnknownKaviWeak", func(t *testing.T) {
		reading := engine.Analyze("कवि")

		// Must NOT produce a harmonic field
		if reading.HarmonicField != nil {
			t.Errorf("unknown कवि should not produce a harmonic field")
		}

		// No passage field should have verified confidence
		for _, pf := range reading.PassageFields {
			if pf.Confidence == ConfidenceVerified {
				t.Errorf("unknown कवि should not have verified-confidence fields, but got %s (conf=%s)",
					pf.Concept, pf.Confidence)
			}
		}
		t.Logf("कवि passage fields: %d, confs: %v",
			len(reading.PassageFields), fieldConfs(reading.PassageFields))
	})
}

func conceptNames(fields PassageFields) []string {
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.Concept
	}
	return names
}

func fieldConfs(fields PassageFields) []string {
	confs := make([]string, len(fields))
	for i, f := range fields {
		confs[i] = f.Confidence
	}
	return confs
}
